package huawei

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Envelope 北向接口统一响应封装
type Envelope struct {
	Success  bool            `json:"success"`
	FailCode FailCode        `json:"failCode"`
	Message  string          `json:"message"`
	Params   json.RawMessage `json:"params"`
	Data     json.RawMessage `json:"data"`
}

// OK 是否完全成功
func (e *Envelope) OK() bool { return e.Success && e.FailCode == FailCodeOK }

// Exception 网关层异常（非业务错误码格式）
type Exception struct {
	ExceptionID   string `json:"exceptionId"`
	ExceptionType string `json:"exceptionType"`
	DescArgs      []any  `json:"descArgs"`
	ReasonArgs    []any  `json:"reasonArgs"`
	DetailArgs    []any  `json:"detailArgs"`
	AdviceArgs    []any  `json:"adviceArgs"`
}

func (e *Exception) Error() string {
	return fmt.Sprintf("[HuaWei]接口异常: %s(%s) reason=%v detail=%v",
		e.ExceptionID, e.ExceptionType, e.ReasonArgs, e.DetailArgs)
}

// APIError 业务错误（failCode 非 0）
type APIError struct {
	Path     string
	FailCode FailCode
	Message  string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[HuaWei]%s 业务错误: failCode=%d message=%s hint=%s",
		e.Path, e.FailCode, e.Message, e.FailCode.Hint())
}

// Retryable 是否值得重试
func (e *APIError) Retryable() bool { return e.FailCode.Retryable() }

// Retryable 是否值得重试
func (c FailCode) Retryable() bool {
	switch c {
	case FailCodeRateLimited, FailCodeSystemBusy, FailCodeServerBusy, FailCodeServerError:
		return true
	default:
		return false
	}
}

// Hint 错误码排查提示
func (c FailCode) Hint() string {
	switch c {
	case FailCodeOK:
		return "成功"
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
	case FailCodeTaskPartial:
		return "控制类任务部分成功"
	case FailCodeTaskFailed:
		return "控制类任务失败"
	default:
		return "详见华为北向接口文档错误码列表"
	}
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
	return fmt.Sprintf("[HuaWei]%s 响应解析失败: status=%d body=%s", e.Path, e.StatusCode, body)
}

// PageResult 分页查询结果
type PageResult[T any] struct {
	Total     int64 `json:"total"`
	PageCount int64 `json:"pageCount"`
	PageNo    int   `json:"pageNo"`
	PageSize  int   `json:"pageSize"`
	List      []T   `json:"list"`
}

// Float64 兼容 JSON 的数值、字符串与 null 三种形态。
// 华为部分接口把文档标注为 Double 的字段以字符串返回（如电站经纬度）
type Float64 float64

// Float 转回 float64
func (f Float64) Float() float64 { return float64(f) }

// Ptr 取指针，便于构造可选字段
func (f Float64) Ptr() *Float64 { return &f }

// UnmarshalJSON 兼容 1.5 / "1.5" / null
func (f *Float64) UnmarshalJSON(b []byte) error {
	s := strings.TrimSpace(string(b))
	if s == "" || s == "null" {
		*f = 0
		return nil
	}
	v, err := strconv.ParseFloat(strings.Trim(strings.TrimSpace(s), `"`), 64)
	if err != nil {
		return fmt.Errorf("[HuaWei]: 无法解析数值 %s", s)
	}
	*f = Float64(v)
	return nil
}

// MarshalJSON 始终输出数值
func (f Float64) MarshalJSON() ([]byte, error) { return json.Marshal(float64(f)) }

// Float64Ptr 返回 Float64 指针，用于可选数值字段
func Float64Ptr(v float64) *Float64 {
	f := Float64(v)
	return &f
}

// ItemMap 数据项 key-value 集合（dataItemMap / dataItems）
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
	switch s := v.(type) {
	case string:
		return s, true
	case json.Number:
		return s.String(), true
	default:
		return fmt.Sprint(v), true
	}
}

// Bool 取布尔值（兼容 0/1 与 "true"/"false"）
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
		if err != nil {
			return false, false
		}
		return p, true
	}
	if f, ok := toFloat64(v); ok {
		return f != 0, true
	}
	return false, false
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

// SplitCodes 拆分逗号分隔的电站编号/设备 SN
func SplitCodes(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// JoinCodes 拼接为接口要求的英文逗号分隔字符串，超过 max 直接截断
func JoinCodes(codes []string) string {
	if len(codes) > MaxBatch {
		codes = codes[:MaxBatch]
	}
	return strings.Join(codes, ",")
}

// JoinInts 拼接整数列表
func JoinInts[T ~int | ~int64](vals []T) string {
	parts := make([]string, 0, len(vals))
	for _, v := range vals {
		parts = append(parts, strconv.FormatInt(int64(v), 10))
	}
	return strings.Join(parts, ",")
}

// Millis 转毫秒时间戳
func Millis(t time.Time) int64 { return t.UnixMilli() }

// ParseTime 解析接口返回的带时区时间，如 2020-02-06T00:00:00+08:00
func ParseTime(s string) (time.Time, error) { return time.Parse(time.RFC3339, s) }

// Time 解析带时区时间，失败返回零值
type Time struct{ time.Time }

func (t *Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		t.Time = time.Time{}
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(t.Format(time.RFC3339))
}
