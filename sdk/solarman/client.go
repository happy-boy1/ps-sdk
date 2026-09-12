// Package solarman 实现 SolarMan（小麦智电 / 小麦商家版）开放平台 OpenAPI。
//
// 平台使用 OAuth2 风格的 Bearer Token：先调 /account/v1.0/token 换取 access_token，
// 之后每个请求在 Authorization 头带上 "bearer {access_token}"。
// 所有响应把业务字段与 code/msg/success/requestId 平铺在同一层。
package solarman

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/imroc/req/v3"
)

// Credentials 接入凭证，在 SolarMan 开放平台创建应用后获取
type Credentials struct {
	AppID     string // 应用 APPID（作为 query 参数下发）
	AppSecret string // 应用密钥

	// 登录身份三选一：Email / Mobile(+CountryCode) / UserName
	Email       string
	Mobile      string
	CountryCode string // 手机号登录时必填，如 86
	UserName    string

	Password       string // 明文密码，SDK 内部做 SHA256 小写
	PasswordSHA256 string // 已算好的 SHA256 小写密文，填写后优先使用

	OrgID   int64  // 商家 ID，不填为 C 端（小麦智电）Token，填了为商家版 Token
	BaseURL string // 数据中心地址，缺省中国区，国际区用 BaseURLGlobal
}

// Validate 校验必填项
func (c Credentials) Validate() error {
	if strings.TrimSpace(c.AppID) == "" {
		return errors.New("[Solarman]: AppID 是空值")
	}
	if strings.TrimSpace(c.AppSecret) == "" {
		return errors.New("[Solarman]: AppSecret 是空值")
	}
	if err := c.validateIdentity(); err != nil {
		return err
	}
	if strings.TrimSpace(c.Password) == "" && strings.TrimSpace(c.PasswordSHA256) == "" {
		return errors.New("[Solarman]: Password 与 PasswordSHA256 不能同时为空")
	}
	return nil
}

func (c Credentials) validateIdentity() error {
	n := 0
	if strings.TrimSpace(c.Email) != "" {
		n++
	}
	if strings.TrimSpace(c.Mobile) != "" {
		n++
		if strings.TrimSpace(c.CountryCode) == "" {
			return errors.New("[Solarman]: 使用手机号登录时 CountryCode 必填")
		}
	}
	if strings.TrimSpace(c.UserName) != "" {
		n++
	}
	switch n {
	case 1:
		return nil
	case 0:
		return errors.New("[Solarman]: Email / Mobile / UserName 必须填写一个")
	default:
		return errors.New("[Solarman]: Email / Mobile / UserName 只能填写一个")
	}
}

// passwordSHA256 返回 SHA256 小写密文
func (c Credentials) passwordSHA256() string {
	if h := strings.TrimSpace(c.PasswordSHA256); h != "" {
		return strings.ToLower(h)
	}
	return SHA256Hex(c.Password)
}

func (c Credentials) baseURL() string {
	if strings.TrimSpace(c.BaseURL) == "" {
		return DefaultBaseURL
	}
	return strings.TrimRight(c.BaseURL, "/")
}

// SHA256Hex 计算 SHA256 小写十六进制摘要，用于 password 字段
func SHA256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// ErrNotLoggedIn 未获取 Token 且未开启自动获取
var ErrNotLoggedIn = errors.New("[Solarman]: 尚未获取 access_token")

// SolarmanSDK SolarMan 开放平台客户端，并发安全
type SolarmanSDK struct {
	creds  Credentials
	client *req.Client

	mu           sync.RWMutex // 保护 token
	token        string
	refreshToken string
	expiry       time.Time
	uid          int64

	loginMu sync.Mutex // 串行化获取 Token

	language     string
	appIDInQuery bool
	maxAttempts  int
	autoLogin    bool
	debugf       func(format string, args ...any)
}

// Option 客户端可选项
type Option func(*SolarmanSDK)

// WithClient 自定义 HTTP 客户端
func WithClient(client *req.Client) Option {
	return func(sdk *SolarmanSDK) {
		if client != nil {
			sdk.client = client
		}
	}
}

// WithTimeout 设置单次请求超时
func WithTimeout(d time.Duration) Option {
	return func(sdk *SolarmanSDK) {
		if d > 0 {
			sdk.client.SetTimeout(d)
		}
	}
}

// WithLanguage 设置 query 参数 language，默认 zh
func WithLanguage(lang string) Option {
	return func(sdk *SolarmanSDK) {
		if strings.TrimSpace(lang) != "" {
			sdk.language = lang
		}
	}
}

// WithAccessToken 复用已有 Token，跳过首次获取
func WithAccessToken(token, refreshToken string, expiresIn time.Duration) Option {
	return func(sdk *SolarmanSDK) {
		if token != "" {
			sdk.token = token
			sdk.refreshToken = refreshToken
			sdk.expiry = time.Now().Add(expiresIn)
		}
	}
}

// WithDebugf 打开调试日志
func WithDebugf(f func(format string, args ...any)) Option {
	return func(sdk *SolarmanSDK) {
		sdk.debugf = f
	}
}

// WithDevMode 打印请求/响应明细
func WithDevMode(flag bool) Option {
	return func(sdk *SolarmanSDK) {
		if flag {
			sdk.client = sdk.client.DevMode()
		}
	}
}

// WithAppIDInQuery 是否在非 Token 接口的 query 中也带上 appId，默认 false
func WithAppIDInQuery(flag bool) Option {
	return func(sdk *SolarmanSDK) {
		sdk.appIDInQuery = flag
	}
}

// WithMaxAttempts 设置鉴权失效时的最大尝试次数，默认 2
func WithMaxAttempts(n int) Option {
	return func(sdk *SolarmanSDK) {
		if n > 0 {
			sdk.maxAttempts = n
		}
	}
}

// WithDisableAutoLogin 关闭自动获取 Token，失效时直接返回错误
func WithDisableAutoLogin() Option {
	return func(sdk *SolarmanSDK) {
		sdk.autoLogin = false
	}
}

// NewSolarmanSDK 创建 SolarMan 客户端
func NewSolarmanSDK(creds Credentials, opts ...Option) (*SolarmanSDK, error) {
	sdk := &SolarmanSDK{
		creds:       creds,
		client:      req.C().SetBaseURL(creds.baseURL()).SetTimeout(30 * time.Second),
		language:    LangZh,
		maxAttempts: 2,
		autoLogin:   true,
	}

	if sdk.client == nil {
		return nil, errors.New("[SolarMan]: HTTP客户端为空")
	}

	for _, opt := range opts {
		opt(sdk)
	}
	return sdk, nil
}

// Client 返回底层 HTTP 客户端
func (sdk *SolarmanSDK) Client() *req.Client { return sdk.client }

// Credentials 返回凭证副本
func (sdk *SolarmanSDK) Credentials() Credentials { return sdk.creds }

// SetCredentials 后续补充凭证
func (sdk *SolarmanSDK) SetCredentials(creds Credentials) {
	sdk.creds = creds
	sdk.client.SetBaseURL(creds.baseURL())
}

// Token 返回当前 access_token
func (sdk *SolarmanSDK) Token() string {
	sdk.mu.RLock()
	defer sdk.mu.RUnlock()
	return sdk.token
}

// RefreshToken 返回 refresh_token
func (sdk *SolarmanSDK) RefreshToken() string {
	sdk.mu.RLock()
	defer sdk.mu.RUnlock()
	return sdk.refreshToken
}

// TokenExpiry 返回 access_token 失效时间
func (sdk *SolarmanSDK) TokenExpiry() time.Time {
	sdk.mu.RLock()
	defer sdk.mu.RUnlock()
	return sdk.expiry
}

// UID 返回 Token 对应的用户 ID
func (sdk *SolarmanSDK) UID() int64 {
	sdk.mu.RLock()
	defer sdk.mu.RUnlock()
	return sdk.uid
}

// ---------------------------------------------------------------------------
// 2.1 获取 Token
// ---------------------------------------------------------------------------

// TokenRequest 2.1 获取 Token 请求体
type TokenRequest struct {
	AppSecret   string `json:"appSecret"`             // 应用密钥
	CountryCode string `json:"countryCode,omitempty"` // 国家代码，手机号登录时必填
	Email       string `json:"email,omitempty"`       // 邮箱登录
	Mobile      string `json:"mobile,omitempty"`      // 手机号登录
	OrgID       *Int64 `json:"orgId,omitempty"`       // 商家 ID，不传为 C 端 Token
	Password    string `json:"password"`              // 密码的 SHA256 小写密文
	Username    string `json:"username,omitempty"`    // 用户名登录
}

// TokenResult 2.1 获取 Token 返回
type TokenResult struct {
	Response
	AccessToken  string `json:"access_token"`  // 访问令牌
	TokenType    string `json:"token_type"`    // 固定 bearer
	RefreshToken string `json:"refresh_token"` // 更新令牌
	ExpiresIn    Int64  `json:"expires_in"`    // 有效期，单位秒（约 60 天）
	Scope        string `json:"scope"`         // 授权范围
	UID          Int64  `json:"uid"`           // 用户 ID
}

// AcquireToken 调用 2.1 获取 Token 并缓存。
// req 中留空的字段会用凭证里的值补齐，因此 AcquireToken(TokenRequest{}) 等价于 Login()
func (sdk *SolarmanSDK) AcquireToken(req TokenRequest) (*TokenResult, error) {
	if strings.TrimSpace(req.AppSecret) == "" {
		req.AppSecret = sdk.creds.AppSecret
	}
	if strings.TrimSpace(req.Password) == "" {
		req.Password = sdk.creds.passwordSHA256()
	}
	if strings.TrimSpace(req.CountryCode) == "" {
		req.CountryCode = sdk.creds.CountryCode
	}
	if strings.TrimSpace(req.Email) == "" {
		req.Email = sdk.creds.Email
	}
	if strings.TrimSpace(req.Mobile) == "" {
		req.Mobile = sdk.creds.Mobile
	}
	if strings.TrimSpace(req.Username) == "" {
		req.Username = sdk.creds.UserName
	}
	if req.OrgID == nil && sdk.creds.OrgID > 0 {
		req.OrgID = Int64Ptr(sdk.creds.OrgID)
	}

	var out TokenResult
	if err := sdk.do(PathToken, req, &out, true); err != nil {
		return nil, err
	}
	if out.AccessToken == "" {
		return nil, errors.New("[Solarman]: 获取 Token 成功但未返回 access_token")
	}

	sdk.mu.Lock()
	sdk.token = out.AccessToken
	sdk.refreshToken = out.RefreshToken
	sdk.uid = out.UID.Int()
	ttl := time.Duration(out.ExpiresIn.Int()) * time.Second
	if ttl <= 0 {
		ttl = 60 * 24 * time.Hour // 文档：默认约 60 天
	}
	sdk.expiry = time.Now().Add(ttl)
	sdk.mu.Unlock()

	sdk.logf("[Solarman] 获取 Token 成功，有效期至 %s", sdk.TokenExpiry().Format(time.RFC3339))
	return &out, nil
}

// Login 按凭证里的登录身份获取 Token
func (sdk *SolarmanSDK) Login() error {
	req := TokenRequest{
		CountryCode: sdk.creds.CountryCode,
		Email:       sdk.creds.Email,
		Mobile:      sdk.creds.Mobile,
		Username:    sdk.creds.UserName,
	}
	if sdk.creds.OrgID > 0 {
		req.OrgID = Int64Ptr(sdk.creds.OrgID)
	}
	_, err := sdk.AcquireToken(req)
	return err
}

// LoginWithOrg 切换到指定商家并重新获取 Token
func (sdk *SolarmanSDK) LoginWithOrg(orgID int64) error {
	sdk.creds.OrgID = orgID
	return sdk.Login()
}

// Logout 清空本地 Token（平台无注销接口）
func (sdk *SolarmanSDK) Logout() {
	sdk.mu.Lock()
	sdk.token = ""
	sdk.refreshToken = ""
	sdk.expiry = time.Time{}
	sdk.uid = 0
	sdk.mu.Unlock()
}

// EnsureToken 未获取 Token 或即将过期时重新获取
func (sdk *SolarmanSDK) EnsureToken() error {
	_, err := sdk.ensureToken()
	return err
}

func (sdk *SolarmanSDK) ensureToken() (string, error) {
	sdk.mu.RLock()
	token, expiry := sdk.token, sdk.expiry
	sdk.mu.RUnlock()

	if token != "" && time.Now().Before(expiry.Add(-5*time.Minute)) {
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

func (sdk *SolarmanSDK) invalidateToken() {
	sdk.mu.Lock()
	sdk.token = ""
	sdk.expiry = time.Time{}
	sdk.mu.Unlock()
}

// ---------------------------------------------------------------------------
// 请求内核
// ---------------------------------------------------------------------------

// do 发送请求并解析到 out，out 需内嵌 Response
func (sdk *SolarmanSDK) do(path string, body, out any, noAuth bool) error {
	raw, err := sdk.call(path, body, noAuth)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return &BadResponseError{Path: path, Body: string(raw)}
	}
	if r, ok := out.(responder); ok {
		if resp := r.base(); !resp.OK() {
			return &APIError{Path: path, Code: resp.Code, Msg: resp.Msg, RequestID: resp.RequestID}
		}
	}
	fillRaw(raw, out)
	return nil
}

// Call 调用任意接口并返回原始响应，用于读取未建模的字段
func (sdk *SolarmanSDK) Call(path string, body any) (ItemMap, error) {
	var out struct {
		Response
	}
	var raw json.RawMessage
	var err error
	if raw, err = sdk.call(path, body, false); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, &BadResponseError{Path: path, Body: string(raw)}
	}
	if !out.OK() {
		return nil, &APIError{Path: path, Code: out.Code, Msg: out.Msg, RequestID: out.RequestID}
	}
	var m ItemMap
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, &BadResponseError{Path: path, Body: string(raw)}
	}
	return m, nil
}

// doAppAuth 用 appId + appSecret（query 参数）鉴权，不带 bearer Token。
// 适用于文档中没有 authorization 头的接口：2.4 注册帐号、2.8 重置密码、5.1 生成验证码
func (sdk *SolarmanSDK) doAppAuth(path string, body, out any) error {
	if strings.TrimSpace(sdk.creds.AppID) == "" || strings.TrimSpace(sdk.creds.AppSecret) == "" {
		return errors.New("[Solarman]: 该接口需要 AppID 与 AppSecret，请先在 Credentials 中填写")
	}

	raw, err := sdk.callWithQuery(path, body, map[string]string{
		QueryAppID:     sdk.creds.AppID,
		QueryAppSecret: sdk.creds.AppSecret,
	}, true)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return &BadResponseError{Path: path, Body: string(raw)}
	}
	if r, ok := out.(responder); ok {
		if resp := r.base(); !resp.OK() {
			return &APIError{Path: path, Code: resp.Code, Msg: resp.Msg, RequestID: resp.RequestID}
		}
	}
	fillRaw(raw, out)
	return nil
}

// call 发送请求，鉴权失效时自动重新获取 Token 并重试
func (sdk *SolarmanSDK) call(path string, body any, noAuth bool) (json.RawMessage, error) {
	return sdk.callWithQuery(path, body, nil, noAuth)
}

// callWithQuery 发送请求，可附加 query 参数
func (sdk *SolarmanSDK) callWithQuery(path string, body any, extra map[string]string, noAuth bool) (json.RawMessage, error) {
	attempts := sdk.maxAttempts
	if attempts < 1 {
		attempts = 1
	}

	var lastErr error
	for i := 0; i < attempts; i++ {
		if i > 0 {
			sdk.invalidateToken()
		}

		token := ""
		if !noAuth {
			t, err := sdk.ensureToken()
			if err != nil {
				return nil, err
			}
			token = t
		}

		raw, err := sdk.post(path, body, token, noAuth, extra)
		if err != nil {
			return nil, err
		}

		var probe Response
		if err := json.Unmarshal(raw, &probe); err != nil {
			return nil, &BadResponseError{Path: path, Body: string(raw)}
		}
		if probe.OK() {
			return raw, nil
		}

		apiErr := &APIError{Path: path, Code: probe.Code, Msg: probe.Msg, RequestID: probe.RequestID}
		if apiErr.IsAuthError() && sdk.autoLogin && !noAuth && i < attempts-1 {
			lastErr = apiErr
			sdk.logf("[Solarman] %s 鉴权失效(%s)，重新获取 Token 后重试", path, probe.Code)
			continue
		}
		return raw, nil // 其余业务错误保留原始报文，由 do 统一转换
	}
	return nil, lastErr
}

func (sdk *SolarmanSDK) post(path string, body any, token string, noAuth bool, extra map[string]string) (json.RawMessage, error) {
	payload, err := marshalBody(body)
	if err != nil {
		return nil, fmt.Errorf("[Solarman]%s 请求体序列化失败: %w", path, err)
	}

	params := map[string]string{QueryLanguage: sdk.language}
	if noAuth || sdk.appIDInQuery {
		params[QueryAppID] = sdk.creds.AppID
	}
	for k, v := range extra {
		params[k] = v
	}

	r := sdk.client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParams(params).
		SetBodyBytes(payload)
	if !noAuth {
		r = r.SetHeader(HeaderAuthorization, BearerPrefix+token)
	}

	resp, err := r.Post(path)
	if err != nil {
		return nil, fmt.Errorf("[Solarman]%s 请求失败: %w", path, err)
	}

	raw := resp.Bytes()
	sdk.logf("[Solarman] %s <= %d", path, resp.StatusCode)
	return raw, nil
}

// marshalBody 请求体序列化，空请求体发送 {}
func marshalBody(body any) ([]byte, error) {
	if body == nil {
		return []byte("{}"), nil
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	if string(payload) == "null" {
		return []byte("{}"), nil
	}
	return payload, nil
}

// fillRaw 为结果填充原始响应，支持顶层结构体、切片与嵌套结构体
func fillRaw(data json.RawMessage, out any) {
	if len(data) == 0 || string(data) == "null" {
		return
	}
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return
	}
	rv = rv.Elem()

	switch rv.Kind() {
	case reflect.Slice:
		var raws []json.RawMessage
		if json.Unmarshal(data, &raws) != nil || rv.Len() != len(raws) {
			return
		}
		for i := 0; i < rv.Len(); i++ {
			fillRaw(raws[i], rv.Index(i).Addr().Interface())
		}

	case reflect.Struct:
		var fields map[string]json.RawMessage
		if json.Unmarshal(data, &fields) != nil {
			return
		}
		if s, ok := out.(rawSetter); ok {
			var m ItemMap
			if json.Unmarshal(data, &m) == nil {
				s.setRaw(m)
			}
		}
		t := rv.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !f.IsExported() { // 匿名的 Response / rawHolder 走上面那条分支
				continue
			}
			key := strings.Split(f.Tag.Get("json"), ",")[0]
			if key == "" || key == "-" {
				continue
			}
			sub, ok := fields[key]
			if !ok || !rv.Field(i).CanAddr() {
				continue
			}
			fillRaw(sub, rv.Field(i).Addr().Interface())
		}
	}
}

func (sdk *SolarmanSDK) logf(format string, args ...any) {
	if sdk.debugf != nil {
		sdk.debugf(format, args...)
	}
}
