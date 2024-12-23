package db

import (
	"fmt"
	"openapi/internal/config"
	"openapi/internal/model"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	cfg := config.GlobalConfig.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	// 自动迁移数据库结构，创建表
	err = DB.AutoMigrate(&AliasRecord{}, &model.User{}, &model.UserSession{}, &model.AliyunAccountInfo{}, &model.EnvironmentConfig{}, &model.CloudflareAccountInfo{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	// 检查是否需要创建管理员账户
	var adminCount int64
	err = DB.Model(&model.User{}).Where("is_admin = ?", true).Count(&adminCount).Error
	if err != nil {
		return fmt.Errorf("failed to check admin count: %v", err)
	}

	// 如果没有管理员账户，创建一个
	if adminCount == 0 {
		admin := &model.User{
			Username:           "admin",
			Password:           "784512",
			IsAdmin:            true,
			Status:             1,
			LoginAttempts:      0,
			LastLoginAttemptAt: time.Now(),
		}
		// 先加密密码
		if err := admin.EncryptPassword(); err != nil {
			return fmt.Errorf("failed to encrypt password: %v", err)
		}
		// 保存到数据库
		if err := DB.Create(admin).Error; err != nil {
			return fmt.Errorf("failed to create admin user: %v", err)
		}
	}

	return nil
}

// AliasRecord 别名记录模型，自动创建 alias_records 表（表名是结构体名称的蛇形复数形式）
type AliasRecord struct {
	gorm.Model         // 包含 ID, CreatedAt, UpdatedAt, DeletedAt 字段
	ServiceName string `gorm:"index"`
	AliasName   string `gorm:"index"`
	VersionID   string
	Description string
}
