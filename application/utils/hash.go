package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// ComputeHash 计算hash值
func ComputeHash(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to serialize data: %w", err)
	}
	hash := sha256.Sum256(jsonData)
	return fmt.Sprintf("%x", hash), nil
}
