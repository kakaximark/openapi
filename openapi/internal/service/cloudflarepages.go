package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"openapi/internal/constants"
	"openapi/internal/logger"
	"strings"
)

// GetPagesProject 获取 Pages 项目列表
func GetPagesProject(countryCode, env string) (*constants.CFResponse[constants.PagesProject], error) {
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
	url := fmt.Sprintf("%s/%s/pages/projects",
		constants.CFBaseURL,
		config.CloudflareConfig.AccountID)

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	logger.Info("Getting Pages projects for account: %s", config.CloudflareConfig.AccountID)

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s",
			resp.StatusCode, string(body))
	}

	// 解析响应
	var cfResp constants.CFResponse[constants.PagesProject]
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &cfResp, nil
}

// GetKVNamespaces 获取 KV 命名空间列表
func GetKVNamespaces(countryCode, env string) (*constants.CFResponse[constants.KVNamespace], error) {
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
	url := fmt.Sprintf("%s/%s/storage/kv/namespaces",
		constants.CFBaseURL,
		config.CloudflareConfig.AccountID)

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	logger.Info("Getting KV namespaces for account: %s", config.CloudflareConfig.AccountID)

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

// CombinePageProjectKVNamespacesAndKeys 组合所有三种数据结构
func CombinePageProjectKVNamespacesAndKeys(pagesResp *constants.CFResponse[constants.PagesProject], kvResp *constants.CFResponse[constants.KVNamespace], countryCode, env string) error {
	var filteredProjects []constants.PagesProject

	// 遍历所有 Pages 项目
	for i := range pagesResp.Result {
		project := pagesResp.Result[i]
		namespaceID := project.Deployment_configs.Production.KV_namespaces.Kv.Namespace_ID

		// 遍历所有 KV 命名空间
		for _, namespace := range kvResp.Result {
			// 当 namespace_id 和 id 相同时，检查 KV Keys
			if namespace.ID == namespaceID && namespaceID != "" {
				// 获取该命名空间的 KV Keys
				keysResp, err := GetKVKeys(countryCode, env, namespaceID)
				if err != nil {
					return fmt.Errorf("failed to get KV keys for namespace %s: %v", namespaceID, err)
				}

				// 检查是否有包含 ProdVersion 的键
				var prodVersionKeys []constants.KVKeys
				for _, key := range keysResp.Result {
					if strings.Contains(key.Name, "ProdVersion") {
						prodVersionKeys = append(prodVersionKeys, key)
					}
				}

				// 如果找到包含 ProdVersion 的键
				if len(prodVersionKeys) > 0 {
					// 更新项目信息
					project.Deployment_configs.Production.KV_namespaces.Kv = constants.KVInfo{
						Namespace_ID:        namespaceID,
						Title:               namespace.Title,
						SupportsURLEncoding: namespace.SupportsUrlEncoding,
						Keys:                prodVersionKeys, // 保留所有包含 ProdVersion 的键
						HasProdVersion:      true,
					}

					// 添加到过滤后的项目列表
					filteredProjects = append(filteredProjects, project)
					logger.Info("Found project with ProdVersion keys - Project: %s, Namespace: %s, Keys count: %d",
						project.Name, namespace.Title, len(prodVersionKeys))
				}
				break
			}
		}
	}

	// 用过滤后的项目替换原始结果
	pagesResp.Result = filteredProjects
	return nil
}

// GetPagesProjectWithKVNamespacesAndKeys 获取完整的组合数据
func GetPagesProjectWithKVNamespacesAndKeys(countryCode, env string) (*constants.CFResponse[constants.PagesProject], error) {
	// 获取 Pages 项目数据
	pagesResp, err := GetPagesProject(countryCode, env)
	if err != nil {
		return nil, fmt.Errorf("failed to get pages project: %v", err)
	}

	// 获取 KV 命名空间数据
	kvResp, err := GetKVNamespaces(countryCode, env)
	if err != nil {
		return nil, fmt.Errorf("failed to get KV namespaces: %v", err)
	}

	// 组合所有数据
	if err := CombinePageProjectKVNamespacesAndKeys(pagesResp, kvResp, countryCode, env); err != nil {
		return nil, fmt.Errorf("failed to combine data: %v", err)
	}

	return pagesResp, nil
}
