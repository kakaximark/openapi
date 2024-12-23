package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"openapi/internal/constants"
	"openapi/internal/logger"
)

// GetKVKeyValues 获取 KV key的值
func GetKVKeyValues(countryCode, env, namespaceId, keyName string) (*constants.CFRawResponse, error) {
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
	url := fmt.Sprintf("%s/%s/storage/kv/namespaces/%s/values/%s", constants.CFBaseURL, config.CloudflareConfig.AccountID, namespaceId, keyName)

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 添加认证头
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	logger.Info("Getting KV key value for countryCode: %s, env: %s, namespace: %s, key: %s", countryCode, env, namespaceId, keyName)

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

	// 尝试解析为 JSON 格式
	var jsonResp constants.CFResponse[constants.KVKeysValues]
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		// 如果不是 JSON 格式，直接返回原始字符串
		return &constants.CFRawResponse{
			RawData: string(body),
		}, nil
	}

	// 如果是 JSON 格式但解析失败
	if !jsonResp.Success {
		return nil, fmt.Errorf("API request failed: %v", jsonResp.Errors)
	}

	// 如果是 JSON 格式且成功解析，转换为原始响应
	return &constants.CFRawResponse{
		RawData: jsonResp.Result[0].Value,
	}, nil
}

// UpdateKVKeyValues 更新 KV key的值
func UpdateKVKeyValues(countryCode, env, namespaceId, keyName string, value string) (*constants.CFResponse[constants.UpdateKVKeysValues], error) {
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

	url := fmt.Sprintf("%s/%s/storage/kv/namespaces/%s/values/%s", constants.CFBaseURL, config.CloudflareConfig.AccountID, namespaceId, keyName)

	req, err := http.NewRequest("PUT", url, strings.NewReader(value))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "text/plain")

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

	var cfResp constants.CFResponse[constants.UpdateKVKeysValues]
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &cfResp, nil
}
