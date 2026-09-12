package ginlong

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/imroc/req/v3"
)

// Credentials 锦浪云 API 凭证，在锦浪云 WEB 端「服务 - API 管理」中获取
type Credentials struct {
	APIID     string // KeyID
	APISecret string // KeySecret
	BaseURL   string // 接口地址，缺省用 DefaultBaseURL
}

// Validate 校验必填项
func (c Credentials) Validate() error {
	if strings.TrimSpace(c.APIID) == "" {
		return errors.New("[Ginlong]: APIID 是空值")
	}
	if strings.TrimSpace(c.APISecret) == "" {
		return errors.New("[Ginlong]: APISecret 是空值")
	}
	return nil
}

func (c Credentials) baseURL() string {
	if strings.TrimSpace(c.BaseURL) == "" {
		return DefaultBaseURL
	}
	return strings.TrimRight(c.BaseURL, "/")
}

// ErrMissingCredentials 凭证未填写
var ErrMissingCredentials = errors.New("[Ginlong]: 未配置 APIID/APISecret，请先在 Credentials 中填写")

// Client 锦浪云 API 客户端，并发安全
type SolisSDK struct {
	creds  Credentials
	client *req.Client

	signContentType   string
	headerContentType string
	now               func() time.Time
	debugf            func(format string, args ...any)
}

// Option 客户端可选项
type Option func(*SolisSDK)

// WithClient 自定义 HTTP 客户端
func WithClient(client *req.Client) Option {
	return func(c *SolisSDK) {
		if client != nil {
			c.client = client
		}
	}
}

// WithTimeout 设置单次请求超时
func WithTimeout(d time.Duration) Option {
	return func(c *SolisSDK) {
		if d > 0 {
			c.client.SetTimeout(d)
		}
	}
}

// WithDebugf 打开调试日志
func WithDebugf(f func(format string, args ...any)) Option {
	return func(c *SolisSDK) {
		c.debugf = f
	}
}

// WithDevMode 打印请求/响应明细
func WithDevMode(flag bool) Option {
	return func(c *SolisSDK) {
		if flag {
			c.client = c.client.DevMode()
		}
	}
}

// WithContentType 覆盖请求头与签名使用的 Content-Type，默认 application/json
func WithContentType(ct string) Option {
	return func(c *SolisSDK) {
		if ct != "" {
			c.signContentType = ct
			c.headerContentType = ct
		}
	}
}

// WithClock 注入时钟，便于测试
func WithClock(now func() time.Time) Option {
	return func(c *SolisSDK) {
		if now != nil {
			c.now = now
		}
	}
}

// NewSolisSDK 创建锦浪云 API 客户端。凭证允许留空，调用接口时才校验
func NewSolisSDK(creds Credentials, opts ...Option) (*SolisSDK, error) {
	if err := creds.Validate(); err != nil {
		return nil, err
	}

	sdk := &SolisSDK{
		creds:             creds,
		client:            req.C().SetBaseURL(creds.baseURL()).SetTimeout(30 * time.Second),
		signContentType:   SignContentType,
		headerContentType: HeaderContentType,
		now:               time.Now,
	}

	if sdk.client == nil {
		return nil, errors.New("[GinLong]: HTTP客户端为空")
	}

	for _, opt := range opts {
		opt(sdk)
	}
	return sdk, nil
}

// SolisSDK 返回底层 HTTP 客户端
func (c *SolisSDK) Client() *req.Client { return c.client }

// Credentials 返回凭证副本
func (c *SolisSDK) Credentials() Credentials { return c.creds }

// SetCredentials 后续补充凭证
func (c *SolisSDK) SetCredentials(creds Credentials) { c.creds = creds }

// ContentMD5 计算 Content-MD5：base64(md5(body))
func ContentMD5(body []byte) string {
	sum := md5.Sum(body)
	return base64.StdEncoding.EncodeToString(sum[:])
}

// Sign 计算签名：base64(HmacSHA1(apiSecret, "POST\n"+md5+"\n"+contentType+"\n"+date+"\n"+resource))
func Sign(apiSecret, contentMD5, contentType, date, resource string) string {
	mac := hmac.New(sha1.New, []byte(apiSecret))
	mac.Write([]byte("POST\n" + contentMD5 + "\n" + contentType + "\n" + date + "\n" + resource))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// Authorization 拼装 Authorization 头：API {apiId}:{sign}
func Authorization(apiID, sign string) string {
	return "API " + apiID + ":" + sign
}

// call 发送一次业务请求
func (c *SolisSDK) call(path string, body any) (*Envelope, error) {
	if err := c.creds.Validate(); err != nil {
		return nil, ErrMissingCredentials
	}

	payload, err := marshalBody(body)
	if err != nil {
		return nil, fmt.Errorf("[Ginlong]%s 请求体序列化失败: %w", path, err)
	}

	contentMD5 := ContentMD5(payload)
	date := c.now().UTC().Format(TimeLayout)
	sign := Sign(c.creds.APISecret, contentMD5, c.signContentType, date, path)

	resp, err := c.client.R().
		SetHeader("Content-Type", c.headerContentType).
		SetHeader("Content-MD5", contentMD5).
		SetHeader("Date", date).
		SetHeader("Authorization", Authorization(strings.TrimSpace(c.creds.APIID), sign)).
		SetBodyBytes(payload).
		Post(path)
	if err != nil {
		return nil, fmt.Errorf("[Ginlong]%s 请求失败: %w", path, err)
	}

	raw := resp.Bytes()
	c.logf("[Ginlong] %s <= %d %s", path, resp.StatusCode, string(raw))

	var env Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, &BadResponseError{Path: path, StatusCode: resp.StatusCode, Body: string(raw)}
	}
	if env.Msg == "" && env.Code == "" && !env.Success && env.Data == nil {
		return nil, &BadResponseError{Path: path, StatusCode: resp.StatusCode, Body: string(raw)}
	}
	return &env, nil
}

// do 发送请求并把 data 解码到 out，失败码返回 *APIError
func (c *SolisSDK) do(path string, body, out any) error {
	env, err := c.call(path, body)
	if err != nil {
		return err
	}
	if !env.OK() {
		return &APIError{Path: path, Code: env.Code, Msg: env.Msg}
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(env.Data, out); err != nil {
		return fmt.Errorf("[Ginlong]%s data 解析失败: %w", path, err)
	}
	fillRaw(env.Data, out)
	return nil
}

// Call 调用任意接口并返回原始数据，用于读取未建模的接口或字段
func (c *SolisSDK) Call(path string, body any) (ItemMap, error) {
	var out ItemMap
	if err := c.do(path, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CallList 调用任意返回数组的接口并返回原始数据列表
func (c *SolisSDK) CallList(path string, body any) ([]ItemMap, error) {
	var out []ItemMap
	if err := c.do(path, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// fillRaw 为结果填充原始数据，支持三层：
// 1) out 实现 rawSetter → 直接填入 data
// 2) out 是切片 → 按元素填入
// 3) out 是结构体 → 按 json tag 递归下钻，让嵌套的 Page.Records 等也能拿到 Raw
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
			if !f.IsExported() { // 匿名的 rawHolder 走上面那条分支
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

func (c *SolisSDK) logf(format string, args ...any) {
	if c.debugf != nil {
		c.debugf(format, args...)
	}
}
