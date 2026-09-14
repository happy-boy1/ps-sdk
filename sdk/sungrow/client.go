package sungrow

import (
	"encoding/json"
	"errors"
	"fmt"
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

var ErrNotLoggedIn = errors.New("[Sungrow]: 未获取 Token 令牌")

type loginError struct{ err error }

func (e *loginError) Error() string { return e.err.Error() }
func (e *loginError) Unwrap() error { return e.err }

type SungrowSDK struct {
	creds  Credentials
	client *req.Client

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

// WithAccessToken 复用已有 Token，跳过首次获取
func WithAccessToken(token, refreshToken string, expiresIn time.Duration) Option {
	return func(sdk *SungrowSDK) {
		if token != "" {
			sdk.token = token
			sdk.expiry = time.Now().Add(expiresIn)
		}
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

// WithDisableAutoLogin 关闭自动获取 Token，失效时直接返回错误
func WithDisableAutoLogin() Option {
	return func(sdk *SungrowSDK) {
		sdk.autoLogin = false
	}
}

func NewSungrowSDK(creds Credentials, opts ...Option) (*SungrowSDK, error) {
	sdk := &SungrowSDK{
		creds:       creds,
		client:      req.C().SetBaseURL(creds.baseURL()).SetTimeout(30 * time.Second),
		maxAttempts: MaxAttempts,
		retryDelay:  RetryBaseDelay,
		autoLogin:   true,
	}

	if sdk.client == nil {
		return nil, errors.New("[Sungrow]: HTTP客户端为空")
	}

	for _, opt := range opts {
		opt(sdk)
	}

	return sdk, nil
}

// Client 返回底层 HTTP 客户端
func (sdk *SungrowSDK) Client() *req.Client { return sdk.client }

// Credentials 返回凭证副本
func (sdk *SungrowSDK) Credentials() Credentials { return sdk.creds }

// SetCredentials 后续补充凭证
func (sdk *SungrowSDK) SetCredentials(creds Credentials) {
	sdk.creds = creds
	sdk.client.SetBaseURL(creds.baseURL())
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

func (sdk *SungrowSDK) Login() (*TokenResult, error) {
	sdk.loginMu.Lock()
	defer sdk.loginMu.Unlock()
	return sdk.login()
}

func (sdk *SungrowSDK) login() (*TokenResult, error) {
	body := TokenRequest{
		Request: Request{
			AppKey: sdk.creds.AppID,
		},
		UserAccount:  sdk.creds.UserAccount,
		UserPassword: sdk.creds.UserPassword,
	}

	var out TokenResult

	err := sdk.do(PathLogin, body, &out)
	if err != nil {
		return nil, err
	}

	if out.Token == "" {
		return nil, fmt.Errorf("[Sungrow] %s: 调用成功，但是token为空值", PathLogin)
	}

	sdk.mu.Lock()
	sdk.token = out.Token
	sdk.expiry = time.Now().Add(TokenTTL)
	sdk.mu.Unlock()

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

func (sdk *SungrowSDK) ensureToken() (string, error) {
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

	if _, err := sdk.Login(); err != nil {
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
	sdk.logf("[Sungrow %s <= %d]", path, resp.StatusCode)

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
	if errors.As(lastErr, &apiErr) && apiErr.ResultCode == ResultCodeCallTooFrequently {
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

func (sdk *SungrowSDK) callOnce(path string, body any) (*Envelope, error) {

	r := sdk.client.R().
		SetHeader(HeaderXAccessKey, sdk.creds.AppSecret).
		SetHeader("Content-Type", "application/json").
		SetHeader("sys_code", "901")

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
			// sdk.invalidateToken() // token 失效，重新登录后重试
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
		return fmt.Errorf("[HuaWei] data 解析失败: %w", err)
	}
	return nil
}
