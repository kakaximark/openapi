package service

import (
	"fmt"
	"openapi/internal/db"
	"openapi/internal/model"
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	fc_open20210406 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

var (
	fcClient    *fc_open20210406.Client
	clientMutex sync.RWMutex
)

// ClientConfig 客户端配置
type ClientConfig struct {
	AliyunConfig      *model.AliyunAccountInfo
	EnvironmentConfig *model.EnvironmentConfig
}

// LoadZoneInfoConfig 加载区域信息配置
func LoadZoneInfoConfig() ([]*model.EnvironmentConfig, error) {
	var envConfigs []*model.EnvironmentConfig

	// 查询所有有效的配置
	if err := db.DB.Where("environment IS NOT NULL AND country_code IS NOT NULL").
		Order("created_at DESC"). // 按创建时间倒序
		Find(&envConfigs).Error; err != nil {
		return nil, fmt.Errorf("failed to load zone info configs: %v", err)
	}

	// 如果没有找到任何配置
	if len(envConfigs) == 0 {
		return nil, fmt.Errorf("no valid zone info configs found")
	}

	return envConfigs, nil
}

// LoadClientConfig 根据环境和区域加载配置
func LoadClientConfig(env, countryCode string) (*ClientConfig, error) {
	// 加载阿里云配置
	var aliyunConfig model.AliyunAccountInfo
	if err := db.DB.Where("environment = ? AND is_active = ?", env, true).First(&aliyunConfig).Error; err != nil {
		return nil, fmt.Errorf("failed to load aliyun config: %v", err)
	}

	// 加载环境配置
	var envConfig model.EnvironmentConfig
	if err := db.DB.Where("environment = ? AND country_code = ?", env, countryCode).First(&envConfig).Error; err != nil {
		return nil, fmt.Errorf("failed to load environment config: %v", err)
	}

	return &ClientConfig{
		AliyunConfig:      &aliyunConfig,
		EnvironmentConfig: &envConfig,
	}, nil
}

// initFCClient 初始化 FC 客户端
func initFCClient(env, region string) error {
	// 加载配置
	config, err := LoadClientConfig(env, region)
	if err != nil {
		return err
	}

	apiConfig := &openapi.Config{
		AccessKeyId:     tea.String(config.AliyunConfig.AccessKeyID),
		AccessKeySecret: tea.String(config.AliyunConfig.AccessKeySecret),
	}

	endpoint := fmt.Sprintf("%s.%s.fc.aliyuncs.com",
		config.AliyunConfig.AccountID,
		config.EnvironmentConfig.Region)
	apiConfig.Endpoint = tea.String(endpoint)

	client, err := fc_open20210406.NewClient(apiConfig)
	if err != nil {
		return err
	}

	clientMutex.Lock()
	fcClient = client
	clientMutex.Unlock()

	return nil
}

// GetFCClient 获取 FC 客户端（带缓存）
func GetFCClient(env, region string) (*fc_open20210406.Client, error) {
	clientMutex.RLock()
	if fcClient != nil {
		defer clientMutex.RUnlock()
		return fcClient, nil
	}
	clientMutex.RUnlock()

	if err := initFCClient(env, region); err != nil {
		return nil, err
	}

	clientMutex.RLock()
	defer clientMutex.RUnlock()
	return fcClient, nil
}

// init 注册配置变更回调
func init() {
	RegisterConfigCallback(func() {
		clientMutex.Lock()
		fcClient = nil
		clientMutex.Unlock()
	})
}
