// Package tools 提供文件读写等通用小工具
package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ReadJSONFile 读取 JSON 文件并反序列化到 data
func ReadJSONFile(path string, data any) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, data)
}

// WriteJSONFile 以缩进格式写入 JSON 文件，目录不存在时自动创建
func WriteJSONFile(path string, data any) error {
	payload, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, payload, 0o644)
}
