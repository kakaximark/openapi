package model

import (
	"gorm.io/gorm"
)

// AliyunAccountInfo 阿里云账号信息表
type AliyunAccountInfo struct {
	gorm.Model
	SiteClient      string `gorm:"column:site_client;type:varchar(50);not null" json:"site_client"`
	AccessKeyID     string `gorm:"column:access_key_id;type:varchar(100);not null" json:"access_key_id"`
	AccessKeySecret string `gorm:"column:access_key_secret;type:varchar(100);not null" json:"access_key_secret"`
	AccountID       string `gorm:"column:account_id;type:varchar(100);not null" json:"account_id"`
	MainAccountID   string `gorm:"column:main_account_id;type:varchar(100);not null" json:"main_account_id"`
	Environment     string `gorm:"column:environment;type:varchar(20);not null" json:"environment"`
	Region          string `gorm:"column:region;type:varchar(50);not null" json:"region"`
	IsActive        bool   `gorm:"column:is_active;default:true" json:"is_active"`
	Description     string `gorm:"column:description;type:varchar(255)" json:"description"`
}

type EnvironmentConfig struct {
	gorm.Model
	Environment string `gorm:"column:environment;type:varchar(20);not null" json:"environment"`
	CountryCode string `gorm:"column:country_code;type:varchar(50);not null" json:"country_code"`
	Region      string `gorm:"column:region;type:varchar(50);not null" json:"region"`
	Variables   string `gorm:"column:variables;type:json" json:"variables"`
	Permissions string `gorm:"column:permissions;type:json" json:"permissions"`
}

// TableName 指定表名
func (AliyunAccountInfo) TableName() string {
	return "aliyun_account_info"
}

func (EnvironmentConfig) TableName() string {
	return "environment_config"
}
