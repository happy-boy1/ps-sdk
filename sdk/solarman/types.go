package solarman

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Response 所有接口的公共响应字段，业务字段与它平铺在同一层
type Response struct {
	Code      Code   `json:"code"`      // 信息码，成功时为 null
	Msg       string `json:"msg"`       // 信息描述
	Success   bool   `json:"success"`   // 是否成功
	RequestID string `json:"requestId"` // 请求标识
}

// OK 是否成功
func (r *Response) OK() bool { return r.Success && !r.Code.IsError() }

// IsError 信息码是否表示失败（成功时为 null 或 1000000）
func (c Code) IsError() bool {
	s := strings.TrimSpace(string(c))
	return s != "" && s != "null" && s != CodeOK
}

// envelope 一次响应的原始报文与已解析的公共字段
type envelope struct {
	raw json.RawMessage
	env Response
}

// Code 兼容 字符串 / 数值 / null 的信息码
func (c *Code) UnmarshalJSON(b []byte) error {
	s := strings.Trim(strings.TrimSpace(string(b)), `"`)
	if s == "null" {
		s = ""
	}
	*c = Code(s)
	return nil
}

// APIError 业务错误
type APIError struct {
	Path      string
	Code      Code
	Msg       string
	RequestID string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[Solarman]%s 业务错误: code=%s msg=%s hint=%s requestId=%s",
		e.Path, e.Code, e.Msg, CodeHint(string(e.Code)), e.RequestID)
}

// Retryable 是否值得重试
func (e *APIError) Retryable() bool {
	switch string(e.Code) {
	case CodeServiceUnavailable, CodeTooFrequently, CodeRPCException:
		return true
	default:
		return false
	}
}

// IsAuthError 是否鉴权失效（需重新获取 Token）
func (e *APIError) IsAuthError() bool {
	switch string(e.Code) {
	case CodeInvalidToken, CodeTokenNotFound, CodeAppIDNotFound, CodeInvalidAppID, CodeInvalidOrgID:
		return true
	default:
		return false
	}
}

// BadResponseError 响应不是预期的 JSON 结构
type BadResponseError struct {
	Path string
	Body string
}

func (e *BadResponseError) Error() string {
	body := e.Body
	if len(body) > 512 {
		body = body[:512] + "..."
	}
	return fmt.Sprintf("[Solarman]%s 响应解析失败: %s", e.Path, body)
}

// ItemMap 原始数据，用于读取未建模的字段
type ItemMap map[string]any

// Raw 取原始值
func (m ItemMap) Raw(key string) (any, bool) {
	v, ok := m[key]
	return v, ok
}

// Float64 取浮点值
func (m ItemMap) Float64(key string) (float64, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}
	return toFloat64(v)
}

// Int 取整数值
func (m ItemMap) Int(key string) (int, bool) {
	f, ok := m.Float64(key)
	return int(f), ok
}

// Int64 取长整数值
func (m ItemMap) Int64(key string) (int64, bool) {
	f, ok := m.Float64(key)
	return int64(f), ok
}

// String 取字符串值
func (m ItemMap) String(key string) (string, bool) {
	v, ok := m[key]
	if !ok || v == nil {
		return "", false
	}
	if s, ok := v.(string); ok {
		return s, true
	}
	return fmt.Sprint(v), true
}

// Bool 取布尔值
func (m ItemMap) Bool(key string) (bool, bool) {
	v, ok := m[key]
	if !ok || v == nil {
		return false, false
	}
	switch b := v.(type) {
	case bool:
		return b, true
	case string:
		p, err := strconv.ParseBool(strings.TrimSpace(b))
		return p, err == nil
	}
	if f, ok := toFloat64(v); ok {
		return f != 0, true
	}
	return false, false
}

// Sub 取子对象
func (m ItemMap) Sub(key string) (ItemMap, bool) {
	if v, ok := m[key].(map[string]any); ok {
		return ItemMap(v), true
	}
	return nil, false
}

// List 取子数组
func (m ItemMap) List(key string) ([]ItemMap, bool) {
	arr, ok := m[key].([]any)
	if !ok {
		return nil, false
	}
	out := make([]ItemMap, 0, len(arr))
	for _, e := range arr {
		if mm, ok := e.(map[string]any); ok {
			out = append(out, ItemMap(mm))
		}
	}
	return out, true
}

func toFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(n), 64)
		return f, err == nil
	default:
		return 0, false
	}
}

// Num 兼容 数值 / 字符串 / null 的数字
type Num float64

// Float 转回 float64
func (n Num) Float() float64 { return float64(n) }

// Ptr 取指针
func (n Num) Ptr() *Num { return &n }

func (n *Num) UnmarshalJSON(b []byte) error {
	s := strings.Trim(strings.TrimSpace(string(b)), `"`)
	if s == "" || s == "null" {
		*n = 0
		return nil
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		*n = 0
		return nil
	}
	*n = Num(v)
	return nil
}

func (n Num) MarshalJSON() ([]byte, error) { return json.Marshal(float64(n)) }

// NumPtr 返回 Num 指针
func NumPtr(v float64) *Num {
	n := Num(v)
	return &n
}

// Int64 兼容 数值 / 字符串 / 空串 / null 的整数，平台的 id 与时间戳两种形态都会出现
type Int64 int64

// Int 转回 int64
func (n Int64) Int() int64 { return int64(n) }

// Ptr 取指针
func (n Int64) Ptr() *Int64 { return &n }

func (n *Int64) UnmarshalJSON(b []byte) error {
	s := strings.Trim(strings.TrimSpace(string(b)), `"`)
	if s == "" || s == "null" {
		*n = 0
		return nil
	}
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		*n = Int64(v)
		return nil
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		*n = Int64(f)
		return nil
	}
	*n = 0
	return nil
}

func (n Int64) MarshalJSON() ([]byte, error) { return json.Marshal(int64(n)) }

// Int64Ptr 返回 Int64 指针
func Int64Ptr(v int64) *Int64 {
	n := Int64(v)
	return &n
}

// Str 兼容 字符串 / 数值 / 布尔 / null 的字符串
type Str string

// String 转回 string
func (s Str) String() string { return string(s) }

func (s *Str) UnmarshalJSON(b []byte) error {
	raw := strings.TrimSpace(string(b))
	if raw == "" || raw == "null" {
		*s = ""
		return nil
	}
	var v string
	if err := json.Unmarshal([]byte(raw), &v); err == nil {
		*s = Str(v)
		return nil
	}
	*s = Str(strings.Trim(raw, `"`))
	return nil
}

func (s Str) MarshalJSON() ([]byte, error) { return json.Marshal(string(s)) }

// IntEnum 兼容 数值 / 字符串 / null 的整数枚举基类型
type IntEnum int

// Int 转回 int
func (e IntEnum) Int() int { return int(e) }

func (e *IntEnum) UnmarshalJSON(b []byte) error {
	var n Int64
	if err := n.UnmarshalJSON(b); err != nil {
		return err
	}
	*e = IntEnum(n)
	return nil
}

func (e IntEnum) MarshalJSON() ([]byte, error) { return json.Marshal(int(e)) }

// PageRequest 通用分页请求，可匿名嵌入各列表请求
type PageRequest struct {
	Page int `json:"page"`           // 第几页，从 1 开始
	Size int `json:"size,omitempty"` // 每页显示个数
}

// Normalize 补默认页码
func (p PageRequest) Normalize() PageRequest {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p
}

// rawHolder 为结果结构体提供原始数据兜底，未建模的字段从 Raw 取
type rawHolder struct {
	Raw ItemMap `json:"-"` // 完整原始响应
}

func (r *rawHolder) setRaw(m ItemMap) { r.Raw = m }

// rawSetter 内部接口，用于填充原始数据
type rawSetter interface{ setRaw(ItemMap) }
