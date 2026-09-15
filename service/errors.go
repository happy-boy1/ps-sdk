package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"ps-sdk/model"
)

// 服务层统一错误，所有适配器返回的错误都可被 errors.As 识别
type Error struct {
	Platform int16  // 平台 ID
	Code     string // 平台短码
	Op       string // 操作名，如 "列出电站"
	Err      error  // 底层原因
}

func (e *Error) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("[%s]%s 失败", e.Code, e.Op)
	}
	return fmt.Sprintf("[%s]%s 失败: %v", e.Code, e.Op, e.Err)
}

// Unwrap 保留底层错误链，便于 errors.Is / errors.As 继续向下匹配 SDK 错误
func (e *Error) Unwrap() error { return e.Err }

// wrapErr 包装适配器内部错误，已是 *Error 时不重复包装
func wrapErr(platform int16, code, op string, err error) error {
	if err == nil {
		return nil
	}
	var svcErr *Error
	if errors.As(err, &svcErr) {
		return err
	}
	return &Error{Platform: platform, Code: code, Op: op, Err: err}
}

// PartialError 部分成功：设备列表拿到了子集，同时存在不可用的设备族。
// 调用方应当写入已拿到的设备，并把 Warnings 作为非致命信息记录。
type PartialError struct {
	Devices  []model.PowerDevice
	Warnings []error
}

func (e *PartialError) Error() string {
	if len(e.Warnings) == 0 {
		return "部分设备不可用"
	}
	return fmt.Sprintf("部分设备不可用(%d 项): %v", len(e.Warnings), e.Warnings[0])
}

// Unwrap 暴露首条告警，便于 errors.Is / errors.As 继续匹配
func (e *PartialError) Unwrap() error {
	if len(e.Warnings) == 0 {
		return nil
	}
	return e.Warnings[0]
}

// ErrCredentialMissing 平台凭据缺失
var ErrCredentialMissing = errors.New("平台凭据缺失")

// ErrPlatformUnknown 平台标识未知
var ErrPlatformUnknown = errors.New("未知平台")

// Logf 服务层日志输出口，默认走标准库 log，可替换为自定义日志实现
var Logf = log.Printf

// Options 请求与同步参数
type Options struct {
	Timeout     time.Duration // 单次请求超时，默认 30s
	MaxAttempts int           // 含首调在内的最大尝试次数，0 表示用 SDK 默认值
	RetryDelay  time.Duration // 重试基础退避，0 表示用 SDK 默认值
	Debug       bool          // 打印请求日志
	DevMode     bool          // 打印完整请求/响应报文
	Language    string        // 部分平台的接口语言，SolarMan 使用，默认 zh
	Concurrency int           // 同步设备时的并发电站数，默认 1
}

// WithDefaults 补齐零值，保证适配器构造与同步时参数合法
func (o Options) WithDefaults() Options {
	if o.Timeout <= 0 {
		o.Timeout = 30 * time.Second
	}
	if o.Concurrency < 1 {
		o.Concurrency = 1
	}
	return o
}
