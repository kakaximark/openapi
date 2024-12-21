package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"openapi/internal/constants"
)

// GetKVNamespaces 获取 KV 命名空间列表
func GetKVNamespaces(accountID, authToken string) (*constants.CFResponse[constants.KVNamespace], error) {
	url := fmt.Sprintf("%s/accounts/%s/storage/kv/namespaces", constants.CFBaseURL, accountID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 添加认证头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var cfResp constants.CFResponse[constants.KVNamespace]
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &cfResp, nil
}
