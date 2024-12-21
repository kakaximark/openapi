package service

import (
	"openapi/internal/db"
	"openapi/internal/model"
)

// CreateAliyunAccount 创建阿里云账号
func CreateAliyunAccount(account *model.AliyunAccountInfo) error {
	return db.DB.Create(account).Error
}

// GetAliyunAccount 获取单个账号信息
func GetAliyunAccount(id uint) (*model.AliyunAccountInfo, error) {
	var account model.AliyunAccountInfo
	err := db.DB.First(&account, id).Error
	return &account, err
}

// ListAliyunAccounts 获取账号列表
func ListAliyunAccounts() ([]model.AliyunAccountInfo, error) {
	var accounts []model.AliyunAccountInfo
	err := db.DB.Find(&accounts).Error
	return accounts, err
}

// UpdateAliyunAccount 更新账号信息
func UpdateAliyunAccount(id uint, account *model.AliyunAccountInfo) error {
	return db.DB.Model(&model.AliyunAccountInfo{}).Where("id = ?", id).Updates(account).Error
}

// DeleteAliyunAccount 删除账号
func DeleteAliyunAccount(id uint) error {
	return db.DB.Delete(&model.AliyunAccountInfo{}, id).Error
}
