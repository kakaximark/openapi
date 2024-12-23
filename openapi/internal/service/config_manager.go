package service

import (
	"fmt"
	"sync"

	"openapi/internal/model"
)

var (
	aliyunConfig     *model.AliyunAccountInfo
	cloudflareConfig *model.CloudflareAccountInfo
	configMutex      sync.RWMutex
	configCallbacks  []ConfigChangeCallback
)

// ConfigChangeCallback 配置变更回调函数类型
type ConfigChangeCallback func()

// RegisterConfigCallback 注册配置变更回调
func RegisterConfigCallback(callback ConfigChangeCallback) {
	configMutex.Lock()
	defer configMutex.Unlock()
	configCallbacks = append(configCallbacks, callback)
}

// notifyConfigChange 通知配置变更
func notifyConfigChange() {
	for _, callback := range configCallbacks {
		callback()
	}
}

// RefreshConfigs 刷新所有配置
func RefreshConfigs() error {
	configMutex.Lock()
	defer configMutex.Unlock()

	// 刷新阿里云配置
	aConfig, err := LoadAliyunConfig()
	if err != nil {
		return fmt.Errorf("failed to refresh aliyun config: %v", err)
	}
	aliyunConfig = aConfig

	// 刷新Cloudflare配置
	cConfig, err := LoadCloudflareConfig()
	if err != nil {
		return fmt.Errorf("failed to refresh cloudflare config: %v", err)
	}
	cloudflareConfig = cConfig

	// 通知配置变更
	notifyConfigChange()
	return nil
}

// ConfigError 配置错误类型
type ConfigError struct {
	Source  string
	Message string
	Err     error
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("%s: %s: %v", e.Source, e.Message, e.Err)
}

// NewConfigError 创建配置错误
func NewConfigError(source, message string, err error) error {
	return &ConfigError{
		Source:  source,
		Message: message,
		Err:     err,
	}
}

// GetAliyunConfig 获取阿里云配置（带缓存）
func GetAliyunConfig() (*model.AliyunAccountInfo, error) {
	configMutex.RLock()
	defer configMutex.RUnlock()

	if aliyunConfig == nil {
		return nil, NewConfigError("aliyun", "config not loaded", nil)
	}
	return aliyunConfig, nil
}

// GetCloudflareConfig 获取Cloudflare配置（带缓存）
func GetCloudflareConfig() (*model.CloudflareAccountInfo, error) {
	configMutex.RLock()
	defer configMutex.RUnlock()

	if cloudflareConfig == nil {
		return nil, NewConfigError("cloudflare", "config not loaded", nil)
	}
	return cloudflareConfig, nil
}

// InvalidateCache 使配置缓存失效
func InvalidateCache() {
	configMutex.Lock()
	defer configMutex.Unlock()

	aliyunConfig = nil
	cloudflareConfig = nil
}
