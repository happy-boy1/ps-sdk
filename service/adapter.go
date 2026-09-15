package service

import (
	"strings"
	"sync"
	"time"

	"ps-sdk/model"
)

// Adapter 平台适配器。四个平台的 SDK 差异（鉴权方式、分页字段、设备类型编码）
// 全部收敛在实现内部，上层只依赖本接口。
type Adapter interface {
	// Platform 返回平台元信息
	Platform() Platform
	// ListPowerStations 拉取账号下全部电站，站内 ID 尚未分配（StationID 为 0）
	ListPowerStations() ([]model.PowerStation, error)
	// ListPowerDevices 拉取指定电站下的设备，stationID 为 model.PowerStation.StationID。
	// 返回 *PartialError 时表示只拿到部分设备族，调用方应写入 Devices 并记录 Warnings
	ListPowerDevices(stationID uint64) ([]model.PowerDevice, error)
}

// ConfigProvider 由配置包（pkg/config）实现，为服务层提供平台凭据与请求参数。
// 服务层只依赖本接口，不直接依赖具体配置实现。
type ConfigProvider interface {
	// Credential 返回平台在配置文件中登记的凭据，未配置返回 false
	Credential(code string) (Credential, bool)
	// Options 返回请求超时、重试、日志与同步并发等运行参数
	Options() Options
}

// provider 全局配置来源，由 InitWith 注入；为 nil 时只用数据库与环境变量
var provider ConfigProvider

// Credential 平台接入凭据。字段按平台取用：
//
//	solarman  AppID/AppSecret/Account/Email/Mobile/UserName/CountryCode/Password/OrgID
//	sungrow   AppID/AppSecret/Account/Password
//	ginlong   AppID/AppSecret
//	huawei    Account/Password
type Credential struct {
	PlatformID int16
	Code       string

	APIURL    string // 接口地址，空则用平台默认地址
	AppID     string // 应用 Key / APIID / 华为 API 账户名
	AppSecret string // 应用密钥 / APISecret

	Account     string // 登录身份：手机号 / 用户名
	Email       string // 邮箱登录，SolarMan 使用
	Mobile      string // 手机号登录，SolarMan 使用
	UserName    string // 用户名登录，SolarMan 使用
	CountryCode string // 手机号国家码，SolarMan 使用
	Password    string // 登录密码
	OrgID       int64  // 商家 ID，SolarMan 使用
}

// Account 返回首个非空的登录身份，顺序为 Account → Email → Mobile → UserName
func (c Credential) Identity() string {
	return firstNonEmpty(c.Account, c.Email, c.Mobile, c.UserName)
}

// adapterBuilder 平台适配器构造函数
type adapterBuilder func(Credential, Options) (Adapter, error)

// adapterBuilders 已注册的平台构造器，key 为平台短码。
// 新增平台时实现 Adapter 并在此登记即可，上层代码无需改动。
var adapterBuilders = map[string]adapterBuilder{}

// registerAdapter 登记平台适配器构造器
func registerAdapter(code string, builder adapterBuilder) {
	adapterBuilders[normalizeCode(code)] = builder
}

// SupportedPlatforms 返回已实现适配器的平台短码
func SupportedPlatforms() []string {
	out := make([]string, 0, len(adapterBuilders))
	for code := range adapterBuilders {
		out = append(out, code)
	}
	return out
}

// NewAdapter 按凭据构造平台适配器
func NewAdapter(cred Credential, opts Options) (Adapter, error) {
	code := normalizeCode(cred.Code)
	builder, ok := adapterBuilders[code]
	if !ok {
		return nil, wrapErr(cred.PlatformID, cred.Code, "构造适配器", ErrPlatformUnknown)
	}
	return builder(cred, opts.WithDefaults())
}

// platformCaller 各适配器共用的请求节流。
// 平台限流阈值各不相同（SolarMan 300 次/10 秒、锦浪 2 次/秒），且多数未公开，
// 因此采用「固定最小间隔 + 被限流后自适应退避」的策略：
// 并发只影响等待顺序，请求始终按间隔串行发出，被限流时退避会自动拉大间隔。
type platformCaller struct {
	mu      sync.Mutex
	next    time.Time     // 下一个可发起的时刻
	backoff time.Duration // 当前附加退避
}

// 节流参数
const (
	minRequestGap = 5 * time.Millisecond // 基础最小间隔
	backoffStep   = 100 * time.Millisecond
	backoffMax    = 5 * time.Second

	// deviceListAttempts 设备列表被限流时的最大尝试次数
	deviceListAttempts = 6
)

// wait 阻塞到下一个可用时隙
func (c *platformCaller) wait() {
	c.mu.Lock()
	now := time.Now()
	if c.next.Before(now) {
		c.next = now
	}
	at := c.next
	c.next = at.Add(minRequestGap).Add(c.backoff)
	c.mu.Unlock()

	if wait := time.Until(at); wait > 0 {
		time.Sleep(wait)
	}
}

// settle 根据本次调用结果调整退避：被限流则加倍，正常则逐步收敛
func (c *platformCaller) settle(rateLimited bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch {
	case rateLimited:
		if c.backoff == 0 {
			c.backoff = backoffStep
		} else if c.backoff < backoffMax {
			c.backoff *= 2
		}
		c.next = time.Now().Add(c.backoff)
	case c.backoff > 0:
		c.backoff /= 2
	}
}

// limitBackoff 提供带节流的调用入口，嵌入到各适配器中
type limitBackoff struct {
	caller platformCaller
}

// throttle 串行化并节流一次平台调用
func (l *limitBackoff) throttle(rateLimited func(error) bool, call func() error) error {
	l.caller.wait()
	err := call()
	l.caller.settle(rateLimited != nil && rateLimited(err))
	return err
}

// throttleRetry 在 throttle 之上对限流错误重试。
// 退避由 caller 自适应放大，因此重试间隔会随限流次数增长而变长。
func (l *limitBackoff) throttleRetry(rateLimited func(error) bool, attempts int, call func() error) error {
	if attempts < 1 {
		attempts = 1
	}

	var err error
	for attempt := 0; attempt < attempts; attempt++ {
		err = l.throttle(rateLimited, call)
		if err == nil {
			return nil
		}
		if rateLimited == nil || !rateLimited(err) {
			return err
		}
		Logf("[service] 平台限流，退避 %.0fms 后重试(%d/%d): %v",
			float64(l.currentBackoff().Milliseconds()), attempt+1, attempts, err)
	}
	return err
}

// currentBackoff 返回当前退避时长，用于日志
func (l *limitBackoff) currentBackoff() time.Duration {
	l.caller.mu.Lock()
	defer l.caller.mu.Unlock()
	return l.caller.backoff
}

// isRateLimitError 判定错误是否为平台限流。
// 通过错误信息匹配各平台的高频提示，避免与具体 SDK 错误类型耦合。
func isRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	message := err.Error()
	for _, keyword := range []string{"TooFrequently", "超频", "调用频繁", "频率", "rate limit"} {
		if strings.Contains(message, keyword) {
			return true
		}
	}
	return false
}
