package model

import "gorm.io/gorm"

// CloudflareAccountInfo Cloudflare账号信息表
type CloudflareAccountInfo struct {
	gorm.Model
	SiteClient      string `gorm:"column:site_client;type:varchar(50);not null" json:"site_client"`
	AccountID       string `gorm:"column:account_id;type:varchar(100);not null" json:"account_id"`
	AccessKeyID     string `gorm:"column:access_key_id;type:varchar(100);not null" json:"access_key_id"`
	AccessKeySecret string `gorm:"column:access_key_secret;type:varchar(100);not null" json:"access_key_secret"`
	ApiToken        string `gorm:"column:api_token;type:varchar(100);not null" json:"api_token"`
	Environment     string `gorm:"column:environment;type:varchar(20);not null;default:'prod'" json:"environment"`
	CountryCode     string `gorm:"column:country_code;type:varchar(20);not null;default:'global'" json:"country_code"`
	IsActive        bool   `gorm:"column:is_active;default:true" json:"is_active"`
	Description     string `gorm:"column:description;type:varchar(255)" json:"description"`
}

// TableName 指定表名
func (CloudflareAccountInfo) TableName() string {
	return "cloudflare_account_info"
}
