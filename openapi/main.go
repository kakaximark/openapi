package main

import (
	"fmt"
	"os"
	"path/filepath"

	"openapi/internal/config"
	"openapi/internal/db"
	"openapi/internal/logger"
	"openapi/internal/router"
	"openapi/internal/service"

	_ "openapi/docs" // 导入 swagger 文档
)

// @title        OpenAPI Service
// @version      1.0
// @description  OpenAPI service for managing Aliyun FC and Cloudflare R2.
// @host         localhost:8080
// @BasePath     /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 初始化日志
	if err := logger.InitLogger(); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	// 加载数据库配置
	configPath := filepath.Join("configs", "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("Configuration file not found at %s", configPath))
	}

	if err := config.LoadConfig(configPath); err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// 初始化数据库连接
	if err := db.InitDB(); err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	// 从数据库加载其他配置
	if err := service.RefreshConfigs(); err != nil {
		logger.Error("Failed to load configs: %v", err)
	} else {
		logger.Info("Successfully loaded all configs")
	}

	// 设置路由
	r := router.SetupRouter()

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
