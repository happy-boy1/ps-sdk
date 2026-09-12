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

// Credentials API 账户凭证
type Credentials struct {
	UserName   string // API 账户名
	SystemCode string // API 账户密码
	BaseURL    string // 管理系统域名，缺省用 DefaultBaseURL
}

// Validate 校验必填项
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

// ErrNotLoggedIn 未登录且未开启自动登录
var ErrNotLoggedIn = errors.New("[HuaWei]: 尚未登录")

// loginError 标记登录阶段失败，避免被重试逻辑再次提交密码
type loginError struct{ err error }

func (e *loginError) Error() string { return e.err.Error() }
func (e *loginError) Unwrap() error { return e.err }

// FusionSolarSDK FusionSolar 北向 OpenAPI 客户端，并发安全
type FusionSolarSDK struct {
	creds  Credentials
	client *req.Client

	mu     sync.RWMutex // 保护 token/expiry
	token  string
	expiry time.Time

	loginMu sync.Mutex // 串行化登录，避免重复登录使旧 token 失效

	debugf      func(format string, args ...any)
	maxAttempts int
	retryDelay  time.Duration
	autoLogin   bool
}

// Option 客户端可选项
type Option func(*FusionSolarSDK)

// WithClient 自定义 HTTP 客户端
func WithClient(client *req.Client) Option {
	return func(sdk *FusionSolarSDK) {
		if client != nil {
			sdk.client = client
		}
	}
}

// WithTimeout 设置单次请求超时
func WithTimeout(d time.Duration) Option {
	return func(sdk *FusionSolarSDK) {
		if d > 0 {
			sdk.client.SetTimeout(d)
		}
	}
}

// WithToken 复用已有 XSRF-TOKEN，避免重复登录
func WithToken(token string) Option {
	return func(sdk *FusionSolarSDK) {
		if token != "" {
			sdk.token = token
			sdk.expiry = time.Now().Add(TokenTTL)
		}
	}
}

// WithDebugf 打开调试日志
func WithDebugf(f func(format string, args ...any)) Option {
	return func(sdk *FusionSolarSDK) {
		sdk.debugf = f
	}
}

// WithDevMode 打印请求/响应明细
func WithDevMode(flag bool) Option {
	return func(sdk *FusionSolarSDK) {
		if flag {
			sdk.client = sdk.client.DevMode()
		}
	}
}

// WithMaxAttempts 设置含首调在内的最大尝试次数
func WithMaxAttempts(n int) Option {
	return func(sdk *FusionSolarSDK) {
		if n > 0 {
			sdk.maxAttempts = n
		}
	}
}

// WithRetryBaseDelay 设置重试基础退避时长
func WithRetryBaseDelay(d time.Duration) Option {
	return func(sdk *FusionSolarSDK) {
		if d > 0 {
			sdk.retryDelay = d
		}
	}
}

// WithDisableAutoLogin 关闭自动登录，token 失效时直接返回 305 错误
func WithDisableAutoLogin() Option {
	return func(sdk *FusionSolarSDK) {
		sdk.autoLogin = false
	}
}

// NewFusionSolarSDK 创建 FusionSolar 北向接口客户端
func NewFusionSolarSDK(creds Credentials, opts ...Option) (*FusionSolarSDK, error) {
	if err := creds.Validate(); err != nil {
		return nil, err
	}

	sdk := &FusionSolarSDK{
		creds:       creds,
		client:      req.C().SetBaseURL(creds.baseURL()).SetTimeout(DefaultTimeout),
		maxAttempts: MaxAttempts,
		retryDelay:  RetryBaseDelay,
		autoLogin:   true,
	}

	if sdk.client == nil {
		return nil, errors.New("[HuaWei]: HTTP 客户端为空")
	}

	for _, opt := range opts {
		opt(sdk)
	}

	return sdk, nil
}

// Client 返回底层 HTTP 客户端
func (sdk *FusionSolarSDK) Client() *req.Client { return sdk.client }

// Credentials 返回凭证副本
func (sdk *FusionSolarSDK) Credentials() Credentials { return sdk.creds }

// Token 返回当前 XSRF-TOKEN
func (sdk *FusionSolarSDK) Token() string {
	sdk.mu.RLock()
	defer sdk.mu.RUnlock()
	return sdk.token
}

// TokenExpiry 返回当前 token 失效时间
func (sdk *FusionSolarSDK) TokenExpiry() time.Time {
	sdk.mu.RLock()
	defer sdk.mu.RUnlock()
	return sdk.expiry
}

// LoginBody 登录请求体
type LoginBody struct {
	UserName   string `json:"userName"`
	SystemCode string `json:"systemCode"`
}

// LogoutBody 注销请求体
type LogoutBody struct {
	XSRFToken string `json:"xsrfToken"`
}

// Login 登录并缓存 XSRF-TOKEN。token 有效期 30 分钟，期间持续调用会自动续期；
// 同一账户仅允许一个在线会话，重复登录会使先前 token 失效
func (sdk *FusionSolarSDK) Login() error {
	sdk.loginMu.Lock()
	defer sdk.loginMu.Unlock()
	return sdk.login()
}

func (sdk *FusionSolarSDK) login() error {
	body := LoginBody{UserName: sdk.creds.UserName, SystemCode: sdk.creds.SystemCode}

	resp, err := sdk.client.R().SetBody(body).SetContentType("application/json").Post(PathLogin)
	if err != nil {
		return fmt.Errorf("[HuaWei]%s 请求失败: %w", PathLogin, err)
	}

	env, err := sdk.decode(PathLogin, resp)
	if err != nil {
		return err
	}
	if !env.OK() {
		return newAPIError(PathLogin, env)
	}

	token := pickToken(resp)
	if token == "" {
		return errors.New("[HuaWei]: 登录响应未返回 XSRF-TOKEN")
	}

	sdk.mu.Lock()
	sdk.token = token
	sdk.expiry = time.Now().Add(TokenTTL)
	sdk.mu.Unlock()

	sdk.logf("[HuaWei] 登录成功，token 有效期至 %s", sdk.TokenExpiry().Format(time.RFC3339))
	return nil
}

// Logout 注销当前 token，建议非必要不调用
func (sdk *FusionSolarSDK) Logout() error {
	token := sdk.Token()
	if token == "" {
		return ErrNotLoggedIn
	}

	resp, err := sdk.client.R().
		SetBody(LogoutBody{XSRFToken: token}).
		SetContentType("application/json").
		Post(PathLogout)
	sdk.invalidateToken()
	if err != nil {
		return fmt.Errorf("[HuaWei]%s 请求失败: %w", PathLogout, err)
	}

	env, err := sdk.decode(PathLogout, resp)
	if err != nil {
		return err
	}
	if !env.OK() {
		return newAPIError(PathLogout, env)
	}
	return nil
}

// EnsureLogin 未登录或 token 即将过期时登录
func (sdk *FusionSolarSDK) EnsureLogin() error {
	_, err := sdk.ensureToken()
	return err
}

func (sdk *FusionSolarSDK) ensureToken() (string, error) {
	sdk.mu.RLock()
	token, expiry := sdk.token, sdk.expiry
	sdk.mu.RUnlock()

	if token != "" && time.Now().Before(expiry.Add(-TokenRefreshAhead)) {
		return token, nil
	}
	if !sdk.autoLogin {
		if token != "" {
			return token, nil
		}
		return "", ErrNotLoggedIn
	}
	if err := sdk.Login(); err != nil {
		return "", err
	}
	return sdk.Token(), nil
}

func (sdk *FusionSolarSDK) invalidateToken() {
	sdk.mu.Lock()
	sdk.token = ""
	sdk.expiry = time.Time{}
	sdk.mu.Unlock()
}

// call 发送业务请求，处理 305 重登录与限流退避重试。
// 非重试类业务错误仍返回信封（供 doTask 检查 failCode），由 do/doTask 统一转换为 *APIError
func (sdk *FusionSolarSDK) call(path string, body any) (*Envelope, error) {
	attempts := sdk.maxAttempts
	if attempts < 1 {
		attempts = 1
	}

	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		if attempt > 0 {
			time.Sleep(sdk.backoff(attempt, lastErr))
		}

		env, err := sdk.callOnce(path, body)
		if err != nil {
			// 登录失败不可重试：连续输错密码会锁定账户
			var loginErr *loginError
			if errors.As(err, &loginErr) {
				return nil, loginErr.err
			}
			// 网关异常与不可解析响应为确定性失败，不重试
			var bad *BadResponseError
			var exc *Exception
			if errors.As(err, &bad) || errors.As(err, &exc) || errors.Is(err, ErrNotLoggedIn) {
				return nil, err
			}
			lastErr = err
			continue
		}
		if env.OK() || env.FailCode == FailCodeTaskPartial {
			return env, nil
		}

		apiErr := newAPIError(path, env)
		switch {
		case apiErr.FailCode == FailCodeNotLogin && sdk.autoLogin:
			sdk.invalidateToken() // token 失效，重新登录后重试
			lastErr = apiErr
		case apiErr.Retryable():
			lastErr = apiErr
		default:
			return env, nil
		}
	}
	return nil, lastErr
}

func (sdk *FusionSolarSDK) callOnce(path string, body any) (*Envelope, error) {
	token, err := sdk.ensureToken()
	if err != nil {
		return nil, &loginError{err: err}
	}

	r := sdk.client.R().
		SetHeader(HeaderToken, token).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json, */*")

	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("[HuaWei] %s 请求体序列化失败: %w", path, err)
		}
		r = r.SetBodyBytes(payload)
	}

	resp, err := r.Post(path)
	if err != nil {
		return nil, fmt.Errorf("[HuaWei] %s 请求失败: %w", path, err)
	}
	return sdk.decode(path, resp)
}

// decode 解析统一响应体
func (sdk *FusionSolarSDK) decode(path string, resp *req.Response) (*Envelope, error) {
	raw := resp.Bytes()
	sdk.logf("[HuaWei] %s <= %d", path, resp.StatusCode)

	var wire struct {
		Envelope
		Exception
	}
	if err := json.Unmarshal(raw, &wire); err != nil {
		return nil, &BadResponseError{Path: path, StatusCode: resp.StatusCode, Body: string(raw)}
	}
	if wire.ExceptionID != "" {
		return nil, &wire.Exception
	}
	if !wire.Success && wire.FailCode == FailCodeOK && wire.Message == "" && wire.Data == nil {
		return nil, &BadResponseError{Path: path, StatusCode: resp.StatusCode, Body: string(raw)}
	}
	return &wire.Envelope, nil
}

// do 发送请求并把 data 解码到 out，failCode 非 0 返回 *APIError
func (sdk *FusionSolarSDK) do(path string, body, out any) error {
	env, err := sdk.call(path, body)
	if err != nil {
		return err
	}
	if !env.OK() {
		return newAPIError(path, env)
	}
	if out == nil {
		return nil
	}
	return decodeData(env.Data, out)
}

// doTask 控制类异步任务：failCode=1 视为部分成功仍返回结果，=2 返回错误但结果已填充
func (sdk *FusionSolarSDK) doTask(path string, body, out any) error {
	env, err := sdk.call(path, body)
	if err != nil {
		return err
	}

	var decodeErr error
	if out != nil {
		decodeErr = decodeData(env.Data, out)
	}
	if env.FailCode == FailCodeOK || env.FailCode == FailCodeTaskPartial {
		return decodeErr
	}
	return newAPIError(path, env)
}

func decodeData(data json.RawMessage, out any) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("[HuaWei] data 解析失败: %w", err)
	}
	return nil
}

// newAPIError 仅在响应失败时返回非 nil
func newAPIError(path string, env *Envelope) *APIError {
	if env.OK() {
		return nil
	}
	return &APIError{Path: path, FailCode: env.FailCode, Message: env.Message}
}

// backoff 429 系统限流固定等待 1 分钟，其余指数退避
func (sdk *FusionSolarSDK) backoff(attempt int, lastErr error) time.Duration {
	var apiErr *APIError
	if errors.As(lastErr, &apiErr) && apiErr.FailCode == FailCodeSystemBusy {
		return time.Minute
	}
	base := sdk.retryDelay
	if base <= 0 {
		base = RetryBaseDelay
	}
	d := base << (attempt - 1)
	if d > 30*time.Second {
		d = 30 * time.Second
	}
	return d
}

func (sdk *FusionSolarSDK) logf(format string, args ...any) {
	if sdk.debugf != nil {
		sdk.debugf(format, args...)
	}
}

// pickToken 兼容新老版本的 token 获取方式（响应头 / Set-Cookie）
func pickToken(resp *req.Response) string {
	if v := resp.GetHeader(HeaderToken); v != "" {
		return v
	}
	for k, vs := range resp.Header {
		if strings.EqualFold(k, HeaderToken) && len(vs) > 0 {
			return vs[0]
		}
	}
	for _, c := range resp.Cookies() {
		if strings.EqualFold(c.Name, HeaderToken) {
			return c.Value
		}
	}
	return ""
}
