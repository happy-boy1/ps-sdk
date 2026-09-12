package ginlong

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Envelope 锦浪统一响应封装
type Envelope struct {
	Success bool            `json:"success"`
	Code    Code            `json:"code"`
	Msg     string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}

// OK 是否成功
func (e *Envelope) OK() bool { return e.Success && e.Code == CodeOK }

// Code 返回码，文档为 String，兼容数值返回
type Code string

func (c *Code) UnmarshalJSON(b []byte) error {
	s := strings.Trim(strings.TrimSpace(string(b)), `"`)
	*c = Code(s)
	return nil
}

// APIError 业务错误
type APIError struct {
	Path string
	Code Code
	Msg  string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[Ginlong]%s 业务错误: code=%s msg=%s hint=%s",
		e.Path, e.Code, e.Msg, CodeHint(string(e.Code)))
}

// BadResponseError 响应不是预期的 JSON 结构
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
	return fmt.Sprintf("[Ginlong]%s 响应解析失败: status=%d body=%s", e.Path, e.StatusCode, body)
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
	v, ok := m[key]
	if !ok {
		return nil, false
	}
	switch t := v.(type) {
	case map[string]any:
		return ItemMap(t), true
	default:
		return nil, false
	}
}

// List 取子数组
func (m ItemMap) List(key string) ([]ItemMap, bool) {
	v, ok := m[key]
	if !ok {
		return nil, false
	}
	arr, ok := v.([]any)
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

// Num 兼容数值、字符串与 null 的数字，锦浪部分字段两种形态都会出现
type Num float64

// Float 转回 float64
func (n Num) Float() float64 { return float64(n) }

// Ptr 取指针，便于构造可选字段
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

// Int64 兼容 数值 / 字符串 / 空串 / null 的整数。
// 锦浪的 id、stationId、userId、collectorId、dataTimestamp 实际都以字符串返回
type Int64 int64

// Int 转回 int64
func (n Int64) Int() int64 { return int64(n) }

// Ptr 取指针，便于构造可选字段
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

// MarshalJSON 按官方请求示例输出字符串形式（{"id":"1308675217944611083"}）
func (n Int64) MarshalJSON() ([]byte, error) { return json.Marshal(strconv.FormatInt(int64(n), 10)) }

// Int64Ptr 返回 Int64 指针
func Int64Ptr(v int64) *Int64 {
	n := Int64(v)
	return &n
}

// Str 兼容 字符串 / 数值 / 布尔 / null 的字符串，用于文档类型标注不可靠的字段
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

// intEnum 兼容 数值 / 字符串 / null 的整数枚举基类
type intEnum int

func (e *intEnum) unmarshal(b []byte) int {
	var n Int64
	if err := n.UnmarshalJSON(b); err != nil {
		return 0
	}
	return int(n)
}

// Page 分页对象
type Page[T any] struct {
	Current int `json:"current"` // 当前页
	Pages   int `json:"pages"`   // 总页数
	Size    int `json:"size"`    // 每页条数
	Total   int `json:"total"`   // 总条数
	Records []T `json:"records"` // 数据
}

// UnmarshalJSON 逐条解析记录，并顺带为每条记录填充原始数据兜底
func (p *Page[T]) UnmarshalJSON(b []byte) error {
	var wire struct {
		Current int               `json:"current"`
		Pages   int               `json:"pages"`
		Size    int               `json:"size"`
		Total   int               `json:"total"`
		Records []json.RawMessage `json:"records"`
	}
	if err := json.Unmarshal(b, &wire); err != nil {
		return err
	}
	p.Current, p.Pages, p.Size, p.Total = wire.Current, wire.Pages, wire.Size, wire.Total
	p.Records = make([]T, 0, len(wire.Records))
	for _, raw := range wire.Records {
		var item T
		if err := json.Unmarshal(raw, &item); err != nil {
			return err
		}
		if s, ok := any(&item).(rawSetter); ok {
			var m ItemMap
			if json.Unmarshal(raw, &m) == nil {
				s.setRaw(m)
			}
		}
		p.Records = append(p.Records, item)
	}
	return nil
}

// Chart 图表数据，锦浪图表接口的 data 有两种形态：
// 对象数组（看 Items）或「指标名 → 并列数组」的对象（看 Raw，键名与 searchinfo 一致）
type Chart[T any] struct {
	rawHolder
	Items []T `json:"-"` // data 为对象数组时，每项一条记录
}

func (c *Chart[T]) UnmarshalJSON(b []byte) error {
	raw := bytes.TrimSpace(b)
	if len(raw) > 0 && raw[0] == '[' {
		return json.Unmarshal(raw, &c.Items)
	}
	return nil // 对象形态由 fillRaw 填入 Raw
}

// rawHolder 为结果结构体提供原始数据兜底，未建模的字段从 Raw 取
type rawHolder struct {
	Raw ItemMap `json:"-"` // 完整原始数据
}

func (r *rawHolder) setRaw(m ItemMap) { r.Raw = m }

// rawSetter 内部接口，用于填充原始数据
type rawSetter interface{ setRaw(ItemMap) }

// PageRequest 通用分页请求，可嵌入各列表请求
type PageRequest struct {
	PageNo   int `json:"pageNo"`   // 页码，从 1 开始
	PageSize int `json:"pageSize"` // 每页条数，最大 100
}

// Normalize 补默认值并限制上限
func (p PageRequest) Normalize() PageRequest {
	if p.PageNo <= 0 {
		p.PageNo = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = DefaultPageSize
	}
	if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}
	return p
}
