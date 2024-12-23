package service

import (
	"openapi/internal/db"
	"openapi/internal/model"
)

// CreateCloudflareAccount 创建Cloudflare账号
func CreateCloudflareAccount(account *model.CloudflareAccountInfo) error {
	return db.DB.Create(account).Error
}

// GetCloudflareAccount 获取单个账号信息
func GetCloudflareAccount(id uint) (*model.CloudflareAccountInfo, error) {
	var account model.CloudflareAccountInfo
	err := db.DB.First(&account, id).Error
	return &account, err
}

// ListCloudflareAccounts 获取账号列表
func ListCloudflareAccounts() ([]model.CloudflareAccountInfo, error) {
	var accounts []model.CloudflareAccountInfo
	err := db.DB.Find(&accounts).Error
	return accounts, err
}

// UpdateCloudflareAccount 更新账号信息
func UpdateCloudflareAccount(id uint, account *model.CloudflareAccountInfo) error {
	return db.DB.Model(&model.CloudflareAccountInfo{}).Where("id = ?", id).Updates(account).Error
}

// DeleteCloudflareAccount 删除账号
func DeleteCloudflareAccount(id uint) error {
	return db.DB.Delete(&model.CloudflareAccountInfo{}, id).Error
}
