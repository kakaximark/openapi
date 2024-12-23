package config

import (
	"encoding/json"
	"os"
)

// MysqlConfig MySQL配置
type MysqlConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// Config 全局配置
type Config struct {
	Mysql MysqlConfig `json:"mysql"`
}

var GlobalConfig Config

// LoadConfig 加载配置
func LoadConfig(configPath string) error {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, &GlobalConfig)
}
