package huawei

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/imroc/req/v3"
)

/*
**************************************************
定义调用华为 FusionSolar API的常量
**************************************************
*/

type Credentials struct {
	UserName   string
	SystemCode string
	BaseURL    string
}

func (c Credentials) Validate() error {
	if strings.TrimSpace(c.UserName) == "" {
		return errors.New("[HuaWei]: UserName 是空值")
	}

	if strings.TrimSpace(c.SystemCode) == "" {
		return errors.New("[HuaWei]: SystemCode 是空值")
	}

	return nil
}

func (c Credentials) baseURL() string {
	if strings.TrimSpace(c.BaseURL) == "" {
		return DefaultBaseURL
	}

	return strings.TrimRight(c.BaseURL, "/")
}

type Envelope struct {
	Success  bool            `json:"success"`
	FailCode FailCode        `json:"failCode"`
	Params   json.RawMessage `json:"params"`
	Message  string          `json:"message"`
	Data     json.RawMessage `json:"data"`
}

func (e *Envelope) OK() bool { return e.Success && e.FailCode == FailCodeOK }

type CommonParams struct {
	CurrentTime  int64  `json:"currentTime"`
	StationCodes string `json:"stationCodes,omitempty"`
}

type APIError struct {
	FailCode FailCode
	Message  string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[HuaWei]业务错误: failCode=%d message=%s hint=%s",
		e.FailCode, e.Message, e.FailCode.hint())
}

func (e *APIError) Retryable() bool {
	switch e.FailCode {
	case FailCodeRateLimited, FailCodeSystemBusy, FailCodeServerBusy, FailCodeServerError:
		return true
	default:
		return false
	}
}

func (c FailCode) hint() string {
	switch c {
	case FailCodeNotLogin:
		return "token 过期或未登录，需重新登录"
	case FailCodeNoPermission:
		return "该 API 账户没有此接口权限"
	case FailCodeRateLimited:
		return "超过单用户限流，请降低调用频率"
	case FailCodeSystemBusy:
		return "系统级限流，需间隔 1 分钟以上重试"
	case FailCodeLoginFailed:
		return "登录失败：账号或密码错误 / 参数名错误 / 账户被锁定 / 密码过期 / 在线会话数达上限"
	case FailCodeStationLimitExceeded:
		return "一次最多支持查询 100 个电站，请分批"
	case FailCodeDailyQuotaExceeded:
		return "单 API 账户单日调用超过最大次数"
	case FailCodeAccountDisabled:
		return "API 账户已禁用"
	case FailCodeAccountExpired:
		return "API 账户已超期"
	case FailCodeTimeInvalid:
		return "开始时间不能大于等于结束时间"
	case FailCodeTimeNegative:
		return "时间参数中存在负值"
	case FailCodeStartTimeFuture:
		return "只输入开始时间时，开始时间不能大于等于当前时间"
	default:
		return "详见华为北向API文档错误码列表"
	}
}

type FusionSolarSDK struct {
	creds  Credentials
	client *req.Client

	mu        sync.Mutex
	token     string
	expiry    time.Time
	loggingIn bool
	loginDone chan struct{}

	debugf func(format string, args ...any)
}

type Option func(*FusionSolarSDK)

func WithClient(client *req.Client) Option {
	return func(sdk *FusionSolarSDK) {
		sdk.client = client
	}
}

func WithTimeout(d time.Duration) Option {
	return func(sdk *FusionSolarSDK) {
		if d > 0 {
			sdk.client.SetTimeout(d)
		}
	}
}

func WithToken(token string) Option {
	return func(sdk *FusionSolarSDK) {
		if token != "" {
			sdk.token = token
			sdk.expiry = time.Now().Add(TokenTTL)
		}
	}
}

func WithDebugF(f func(formath string, args ...any)) Option {
	return func(sdk *FusionSolarSDK) {
		sdk.debugf = f
	}
}

func NewFunsionSolarSDK(creds Credentials, opts ...Option) (*FusionSolarSDK, error) {
	if err := creds.Validate(); err != nil {
		return nil, err
	}
	client := req.C()

	sdk := &FusionSolarSDK{
		creds:  creds,
		client: client,
	}

	for _, opt := range opts {
		opt(sdk)
	}

	return sdk, nil
}
