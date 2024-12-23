package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"openapi/internal/constants"
	"openapi/internal/logger"
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

	// 读取响应后，先打印完整的响应内容
	// logger.Info("Raw response: %s", string(body))

	// 解析响应
	var cfResp constants.CFResponse[constants.PagesProject]
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// 打印每个项目的详细信息
	// for _, project := range cfResp.Result {
	// 	logger.Info("Project details - Name: %s, ID: %s", project.Name, project.ID)
	// 	if project.Deployment_configs.Production != nil {
	// 		logger.Info("Production config: %s", string(project.Deployment_configs.Production))
	// 	} else {
	// 		logger.Info("Production config is nil for project: %s", project.Name)
	// 	}
	// }

	return &cfResp, nil
}
