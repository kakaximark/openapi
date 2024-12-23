package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"openapi/internal/constants"
	"openapi/internal/logger"
)

// GetKVKeys 获取 KV keys列表
func GetKVKeys(countryCode, env, namespaceId string) (*constants.CFResponse[constants.KVKeys], error) {
	// 获取默认客户端
	client := GetDefaultClient()

	// 加载配置
	if err := client.LoadConfig(env, countryCode); err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	// 获取配置信息
	config := client.GetConfig()
	if config == nil || config.CloudflareConfig == nil {
		return nil, fmt.Errorf("cloudflare config not loaded")
	}

	// 构建请求 URL
	url := fmt.Sprintf("%s/%s/storage/kv/namespaces/%s/keys",
		constants.CFBaseURL,
		config.CloudflareConfig.AccountID,
		namespaceId)

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	logger.Info("Getting KV keys for account: %s, namespace: %s", config.CloudflareConfig.AccountID, namespaceId)

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

	var cfResp constants.CFResponse[constants.KVKeys]
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &cfResp, nil
}
