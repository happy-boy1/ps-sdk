package sungrow

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/imroc/req/v3"
)

type Credentials struct {
	AppID     string
	AppSecret string

	UserAccount  string
	UserPassword string

	BaseURL string
}

func (c Credentials) Validate() error {
	if strings.TrimSpace(c.AppID) == "" {
		return errors.New("[Sungrow]: AppID 是空值")
	}

	if strings.TrimSpace(c.AppSecret) == "" {
		return errors.New("[Sungrow]: AppSecret 是空值")
	}

	if strings.TrimSpace(c.UserAccount) == "" {
		return errors.New("[Sungrow]: UserAccount 是空值")
	}

	if strings.TrimSpace(c.UserPassword) == "" {
		return errors.New("[Sungrow]: UserPassword 是空值")
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
var ErrNotLoggedIn = errors.New("[Sungrow]: 未获取 Token 令牌")

// ErrMissingCredentials 凭证未填写。构造客户端不再校验凭证，首次请求时才返回本错误
var ErrMissingCredentials = errors.New("[Sungrow]: 未配置 AppID/AppSecret/UserAccount/UserPassword，请先在 Credentials 中填写")

// loginError 标记登录阶段失败，避免被重试逻辑再次提交密码
type loginError struct{ err error }

func (e *loginError) Error() string { return e.err.Error() }
func (e *loginError) Unwrap() error { return e.err }

// tokenCarrier 需要注入 appkey 与 token 的请求体。
// 业务接口的请求体都是「内嵌 Request 的结构体指针」，如 &PowerStationListRequest{}，
// 由 callOnce 统一注入，接口方法内部不必再关心 token
type tokenCarrier interface {
	setAuth(appKey, token string)
}

type SungrowSDK struct {
	credsMu sync.RWMutex // 保护 creds，允许运行期更换凭证
	creds   Credentials
	client  *req.Client

	mu     sync.RWMutex
	token  string
	expiry time.Time
	uid    int64

	loginMu sync.Mutex

	retryDelay  time.Duration
	maxAttempts int
	autoLogin   bool
	debugf      func(format string, args ...any)
}

// Option 客户端可选项
type Option func(*SungrowSDK)

// WithClient 自定义 HTTP 客户端
func WithClient(client *req.Client) Option {
	return func(sdk *SungrowSDK) {
		if client != nil {
			sdk.client = client
		}
	}
}

// WithTimeout 设置单次请求超时
func WithTimeout(d time.Duration) Option {
	return func(sdk *SungrowSDK) {
		if d > 0 {
			sdk.client.SetTimeout(d)
		}
	}
}

// WithAccessToken 复用已有 Token，跳过首次获取。
// 平台未提供 token 有效期字段，expiresIn 为 0 时按 TokenTTL 估算
func WithAccessToken(token string, expiresIn time.Duration) Option {
	return func(sdk *SungrowSDK) {
		if strings.TrimSpace(token) == "" {
			return
		}
		if expiresIn <= 0 {
			expiresIn = TokenTTL
		}
		sdk.mu.Lock()
		sdk.token = token
		sdk.expiry = time.Now().Add(expiresIn)
		sdk.mu.Unlock()
	}
}

// WithDebugf 打开调试日志
func WithDebugf(f func(format string, args ...any)) Option {
	return func(sdk *SungrowSDK) {
		sdk.debugf = f
	}
}

// WithDevMode 打印请求/响应明细
func WithDevMode(flag bool) Option {
	return func(sdk *SungrowSDK) {
		if flag {
			sdk.client = sdk.client.DevMode()
		}
	}
}

// WithMaxAttempts 设置鉴权失效时的最大尝试次数，默认 2
func WithMaxAttempts(n int) Option {
	return func(sdk *SungrowSDK) {
		if n > 0 {
			sdk.maxAttempts = n
		}
	}
}

// WithRetryBaseDelay 设置重试基础退避时长，默认 RetryBaseDelay
func WithRetryBaseDelay(d time.Duration) Option {
	return func(sdk *SungrowSDK) {
		if d > 0 {
			sdk.retryDelay = d
		}
	}
}

// WithDisableAutoLogin 关闭自动获取 Token，失效时直接返回错误
func WithDisableAutoLogin() Option {
	return func(sdk *SungrowSDK) {
		sdk.autoLogin = false
	}
}

// NewSungrowSDK 创建阳光云 OpenAPI 客户端。凭证允许留空，调用接口时才校验
func NewSungrowSDK(creds Credentials, opts ...Option) (*SungrowSDK, error) {
	sdk := &SungrowSDK{
		creds:       creds,
		client:      req.C().SetBaseURL(creds.baseURL()).SetTimeout(30 * time.Second),
		maxAttempts: MaxAttempts,
		retryDelay:  RetryBaseDelay,
		autoLogin:   true,
	}

	for _, opt := range opts {
		opt(sdk)
	}

	return sdk, nil
}

// Client 返回底层 HTTP 客户端
func (sdk *SungrowSDK) Client() *req.Client { return sdk.client }

// Credentials 返回凭证副本
func (sdk *SungrowSDK) Credentials() Credentials {
	sdk.credsMu.RLock()
	defer sdk.credsMu.RUnlock()
	return sdk.creds
}

// SetCredentials 运行期更换凭证，会同时清空已缓存的 token
func (sdk *SungrowSDK) SetCredentials(creds Credentials) {
	sdk.credsMu.Lock()
	sdk.creds = creds
	sdk.credsMu.Unlock()

	sdk.client.SetBaseURL(creds.baseURL())
	sdk.invalidateToken()
}

// Token 返回当前 access_token
func (sdk *SungrowSDK) Token() string {
	sdk.mu.RLock()
	defer sdk.mu.RUnlock()
	return sdk.token
}

// TokenExpiry 返回 access_token 失效时间
func (sdk *SungrowSDK) TokenExpiry() time.Time {
	sdk.mu.RLock()
	defer sdk.mu.RUnlock()
	return sdk.expiry
}

// UID 返回 Token 对应的用户 ID
func (sdk *SungrowSDK) UID() int64 {
	sdk.mu.RLock()
	defer sdk.mu.RUnlock()
	return sdk.uid
}

// 获取Token
type TokenRequest struct {
	Request

	UserAccount  string `json:"user_account"`
	UserPassword string `json:"user_password"`
}

type TokenResult struct {
	UserMasterOrgId   string `json:"user_master_org_id"`
	MobileTel         string `json:"mobile_tel"`
	UserName          string `json:"user_name"`
	Language          string `json:"language"`
	Token             string `json:"token"`
	ErrTimes          string `json:"err_times"`
	UserId            string `json:"user_id"`
	LoginState        string `json:"login_state"`
	DisableTime       string `json:"disable_time"`
	CountryName       string `json:"country_name"`
	UserAccount       string `json:"user_account"`
	UserMasterOrgName string `json:"user_master_org_name"`
	Email             string `json:"email"`
	CountryId         string `json:"country_id"`
}

// Login 获取 Token 并缓存。总是发起一次登录请求；
// 业务接口调用无需显式登录，token 为空或失效时会自动获取（见 ensureToken）
func (sdk *SungrowSDK) Login() (*TokenResult, error) {
	sdk.loginMu.Lock()
	defer sdk.loginMu.Unlock()
	return sdk.login()
}

func (sdk *SungrowSDK) login() (*TokenResult, error) {
	creds := sdk.Credentials()
	if err := creds.Validate(); err != nil {
		return nil, ErrMissingCredentials
	}

	body := TokenRequest{
		Request: Request{
			AppKey: creds.AppID,
		},
		UserAccount:  creds.UserAccount,
		UserPassword: creds.UserPassword,
	}

	var out TokenResult

	err := sdk.do(PathLogin, body, &out)
	if err != nil {
		return nil, err
	}

	if out.Token == "" {
		// 平台以 result_data.login_state 表达失败原因，token 为空即登录未成功
		if state, ok := ParseLoginState(out.LoginState); ok && state != LoginStateOK {
			return nil, fmt.Errorf("[Sungrow] %s 登录失败: login_state=%s (%s)",
				PathLogin, strings.TrimSpace(out.LoginState), state)
		}
		return nil, fmt.Errorf("[Sungrow] %s: 调用成功，但是token为空值", PathLogin)
	}

	sdk.mu.Lock()
	sdk.token = out.Token
	sdk.expiry = time.Now().Add(TokenTTL)
	if uid, err := strconv.ParseInt(strings.TrimSpace(out.UserId), 10, 64); err == nil {
		sdk.uid = uid
	}
	sdk.mu.Unlock()

	sdk.logf("[Sungrow] 登录成功，token 有效期至 %s", sdk.TokenExpiry().Format(time.RFC3339))

	return &out, nil
}

func (sdk *SungrowSDK) Logout() {
	sdk.mu.Lock()
	sdk.token = ""
	sdk.expiry = time.Time{}
	sdk.uid = 0
	sdk.mu.Unlock()
}

func (sdk *SungrowSDK) EnsureLogin() error {
	_, err := sdk.ensureToken()
	return err
}

// cachedToken 返回尚未失效的缓存 token，失效提前量 TokenRefreshAhead
func (sdk *SungrowSDK) cachedToken() (string, bool) {
	sdk.mu.RLock()
	defer sdk.mu.RUnlock()

	if sdk.token == "" || !time.Now().Before(sdk.expiry.Add(-TokenRefreshAhead)) {
		return "", false
	}
	return sdk.token, true
}

// ensureToken 返回可用 token：命中缓存直接复用，否则获取并缓存后再返回
func (sdk *SungrowSDK) ensureToken() (string, error) {
	if token, ok := sdk.cachedToken(); ok {
		return token, nil
	}

	if !sdk.autoLogin {
		// 关闭自动登录时沿用已有 token，由服务端判定是否失效
		if token := sdk.Token(); token != "" {
			return token, nil
		}

		return "", ErrNotLoggedIn
	}

	return sdk.loginShared()
}

// loginShared 串行化获取 Token：取到锁后二次确认缓存，
// 避免并发请求重复登录（同账号重复登录可能使先前 token 失效）
func (sdk *SungrowSDK) loginShared() (string, error) {
	sdk.loginMu.Lock()
	defer sdk.loginMu.Unlock()

	if token, ok := sdk.cachedToken(); ok {
		return token, nil
	}

	if _, err := sdk.login(); err != nil {
		return "", err
	}

	return sdk.Token(), nil
}

func (sdk *SungrowSDK) invalidateToken() {
	sdk.mu.Lock()
	sdk.token = ""
	sdk.expiry = time.Time{}
	sdk.mu.Unlock()
}

func (sdk *SungrowSDK) decode(path string, resp *req.Response) (*Envelope, error) {
	raw := resp.Bytes()
	sdk.logf("[Sungrow] %s <= %d", path, resp.StatusCode)

	var wire Envelope

	if err := json.Unmarshal(raw, &wire); err != nil {
		return nil, &BadResponseError{Path: path, StatusCode: resp.StatusCode, Body: string(raw)}
	}

	if wire.ResultCode == ResultCodeSuccess && wire.ResultMsg == "" && wire.ResultData == nil {
		return nil, &BadResponseError{Path: path, StatusCode: resp.StatusCode, Body: string(raw)}
	}

	return &wire, nil
}

func (sdk *SungrowSDK) backoff(attempt int, lastErr error) time.Duration {
	var apiErr *APIError
	if errors.As(lastErr, &apiErr) {
		switch apiErr.ResultCode {
		case ResultCodeCallTooFrequently:
			return time.Minute
		case ResultCodeErTokenLoginInvalid:
			// token 失效已丢弃缓存并重新获取，无需退避
			return 0
		}
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

// callOnce 发送单次请求。
// 除登录接口外，请求体（内嵌 Request 的结构体指针）会在发送前统一注入 appkey 与 token；
// token 为空或已失效时先获取并缓存，再发起本次调用
func (sdk *SungrowSDK) callOnce(path string, body any) (*Envelope, error) {
	creds := sdk.Credentials()
	if err := creds.Validate(); err != nil {
		return nil, &loginError{err: ErrMissingCredentials}
	}

	if body != nil && path != PathLogin {
		carrier, ok := body.(tokenCarrier)
		if !ok {
			return nil, fmt.Errorf("[Sungrow] %s 请求体需为内嵌 Request 的结构体指针（如 &Req{}），当前类型 %T", path, body)
		}

		token, err := sdk.ensureToken()
		if err != nil {
			// 标记为登录失败：连续输错密码会锁定账户，不重试
			return nil, &loginError{err: err}
		}

		carrier.setAuth(creds.AppID, token)
	}

	r := sdk.client.R().
		SetHeader(HeaderXAccessKey, creds.AppSecret).
		SetHeader("Content-Type", "application/json").
		SetHeader(HeaderSysCode, "901")

	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("[Sungrow] %s 请求体序列化失败: %w", path, err)
		}
		r = r.SetBodyBytes(payload)
	}

	resp, err := r.Post(path)
	if err != nil {
		return nil, fmt.Errorf("[Sungrow] %s 请求失败: %w", path, err)
	}

	return sdk.decode(path, resp)
}

func (sdk *SungrowSDK) call(path string, body any) (*Envelope, error) {
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
			if errors.As(err, &bad) || errors.Is(err, ErrNotLoggedIn) {
				return nil, err
			}
			lastErr = err
			continue
		}
		if env.OK() {
			return env, nil
		}

		apiErr := newAPIError(path, env)
		switch {
		case apiErr.ResultCode == ResultCodeErTokenLoginInvalid && sdk.autoLogin:
			// token 失效：丢弃缓存，下次尝试会重新获取 token 再调用
			sdk.invalidateToken()
			sdk.logf("[Sungrow] %s token 失效(%s)，重新获取后重试", path, apiErr.ResultCode)
			lastErr = apiErr
		case apiErr.Retryable():
			lastErr = apiErr
		default:
			return env, nil
		}
	}
	return nil, lastErr
}

func (sdk *SungrowSDK) do(path string, body, out any) error {
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

	return decodeData(env.ResultData, out)
}

func (sdk *SungrowSDK) logf(format string, args ...any) {
	if sdk.debugf != nil {
		sdk.debugf(format, args...)
	}
}

// newAPIError 仅在响应失败时返回非 nil
func newAPIError(path string, env *Envelope) *APIError {
	if env.OK() {
		return nil
	}
	return &APIError{Path: path, ResultCode: env.ResultCode, Message: env.ResultMsg}
}

func decodeData(data json.RawMessage, out any) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("[Sungrow] data 解析失败: %w", err)
	}
	return nil
}
