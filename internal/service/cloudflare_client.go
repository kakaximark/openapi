package service

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"openapi/internal/db"
	"openapi/internal/logger"
	"openapi/internal/model"
)

// CloudflareClient Cloudflare 客户端结构体
type CloudflareClient struct {
	httpClient *http.Client
	config     *CFClientConfig
	mutex      sync.RWMutex
}

// CFClientConfig 客户端配置
type CFClientConfig struct {
	CloudflareConfig *model.CloudflareAccountInfo
}

var (
	defaultClient *CloudflareClient
	once          sync.Once
)

// NewCloudflareClient 创建新的 Cloudflare 客户端
func NewCloudflareClient() *CloudflareClient {
	return &CloudflareClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetDefaultClient 获取默认客户端（单例模式）
func GetDefaultClient() *CloudflareClient {
	once.Do(func() {
		defaultClient = NewCloudflareClient()
	})
	return defaultClient
}

// LoadConfig 加载配置
func (c *CloudflareClient) LoadConfig(env, countryCode string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 加载 Cloudflare 配置
	var cloudflareConfig model.CloudflareAccountInfo
	if err := db.DB.Where("environment = ? AND is_active = ?", env, true).First(&cloudflareConfig).Error; err != nil {
		return fmt.Errorf("failed to load cloudflare config: %v", err)
	}

	c.config = &CFClientConfig{
		CloudflareConfig: &cloudflareConfig,
	}

	logger.Info("Loaded Cloudflare config for env: %s, country: %s", env, countryCode)
	return nil
}

// GetConfig 获取当前配置
func (c *CloudflareClient) GetConfig() *CFClientConfig {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.config
}

// Do 执行 HTTP 请求
func (c *CloudflareClient) Do(req *http.Request) (*http.Response, error) {
	// 添加通用请求头
	if c.config != nil && c.config.CloudflareConfig != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.CloudflareConfig.ApiToken))
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

// Reset 重置客户端配置
func (c *CloudflareClient) Reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.config = nil
}

// GetClient 获取底层 HTTP 客户端
func (c *CloudflareClient) GetClient() *http.Client {
	return c.httpClient
}
