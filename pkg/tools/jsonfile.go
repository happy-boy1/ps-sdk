package tools

import (
	"encoding/json"
	"os"
)

func ReadJsonFile(path string, data any) error { return nil }

func WriteJsonFile(path string, data any) error {
	dataBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(dataBytes)
	if err != nil {
		return err
	}
	return nil
}
