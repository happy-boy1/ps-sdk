package sungrow

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Request struct {
	AppKey      string          `json:"appkey"`
	Token       string          `json:"token"`
	Lang        string          `json:"lang,omitempty"`
	ApiKeyParam json.RawMessage `json:"api_key_param,omitempty"`
	TimeStamp   string          `json:"timestamp,omitempty"`
	Nonce       string          `json:"nonce,omitempty"`
}

// setAuth 由 SungrowSDK.callOnce 在发送前调用，注入 appkey 与当前有效 token。
// 调用方显式指定 appkey 时保留其值，token 始终以 SDK 缓存的值为准。
func (r *Request) setAuth(appKey, token string) {
	if strings.TrimSpace(r.AppKey) == "" {
		r.AppKey = appKey
	}
	r.Token = token
}

type Envelope struct {
	ReqSerialNum string          `json:"req_serial_num"`
	ResultCode   ResultCode      `json:"result_code"`
	ResultMsg    string          `json:"result_msg"`
	ResultData   json.RawMessage `json:"result_data"`
}

func (e *Envelope) OK() bool {
	return e.ResultCode == ResultCodeSuccess
}

// APIError
type APIError struct {
	Path       string
	ResultCode ResultCode
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[Sungrow] %s 业务错误: resultCode=%s messgae=%s hint=%s",
		e.Path, e.ResultCode, e.Message, e.ResultCode.Hint())
}

func (e *APIError) Retryable() bool { return e.ResultCode.Retryable() }

type BadResponseError struct {
	Path       string
	StatusCode int
	Body       string
}

func (e *BadResponseError) Error() string {
	body := e.Body
	if len(body) > 512 {
		body = body[:512] + "..."
	}
	return fmt.Sprintf("[Sungrow] %s 响应解析失败: status=%d body=%s",
		e.Path, e.StatusCode, e.Body)
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
