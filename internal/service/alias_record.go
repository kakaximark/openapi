package service

import (
	"openapi/internal/db"
)

// CreateAliasRecord 创建别名记录
func CreateAliasRecord(serviceName, aliasName, versionID, description string) error {
	record := &db.AliasRecord{
		ServiceName: serviceName,
		AliasName:   aliasName,
		VersionID:   versionID,
		Description: description,
	}
	return db.DB.Create(record).Error
}

// GetAliasRecords 获取别名记录列表
func GetAliasRecords(serviceName string) ([]db.AliasRecord, error) {
	var records []db.AliasRecord
	query := db.DB
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}
	err := query.Find(&records).Error
	return records, err
}
