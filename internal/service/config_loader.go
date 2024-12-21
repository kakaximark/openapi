package service

import (
	"openapi/internal/db"
	"openapi/internal/model"
)

// LoadAliyunConfig 从数据库加载激活的阿里云配置
func LoadAliyunConfig() (*model.AliyunAccountInfo, error) {
	var config model.AliyunAccountInfo
	err := db.DB.Where("is_active = ?", true).First(&config).Error
	return &config, err
}

// LoadEnvironmentConfig 从数据库加载激活的环境配置
func LoadEnvironmentConfig() (*model.EnvironmentConfig, error) {
	var config model.EnvironmentConfig
	err := db.DB.Where("is_active = ?", true).First(&config).Error
	return &config, err
}

// LoadCloudflareConfig 从数据库加载激活的Cloudflare配置
func LoadCloudflareConfig() (*model.CloudflareAccountInfo, error) {
	var config model.CloudflareAccountInfo
	err := db.DB.Where("is_active = ?", true).First(&config).Error
	return &config, err
}
