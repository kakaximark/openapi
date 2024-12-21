package constants

import (
	"encoding/json"
)

// Cloudflare API 相关常量
const (
	// CFBaseURL Cloudflare API 基础 URL
	CFBaseURL = "https://api.cloudflare.com/client/v4/accounts"
)

// KVNamespace 表示 KV 命名空间信息
type KVNamespace struct {
	ID                  string `json:"id"`
	Title               string `json:"title"`
	SupportsUrlEncoding bool   `json:"supports_url_encoding"`
}

// KVKeys 表示 KV keys信息
type KVKeys struct {
	Name string `json:"name"`
}

// KVKeysValues 表示 KV keys值信息
type KVKeysValues struct {
	Value string
}

// UpdateKVKeysValues 表示更新 KV keys值
type UpdateKVKeysValues struct {
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Success  bool          `json:"success"`
	Result   []interface{} `json:"result"`
}

// CFResponse Cloudflare API 通用响应结构
type CFResponse[T any] struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Result   []T           `json:"result"`
}

// CFRawResponse 处理非 JSON 格式的响应
type CFRawResponse struct {
	RawData string
}

// PagesProject 表示 Pages 项目信息
type PagesProject struct {
	ID                 string            `json:"id"`
	Name               string            `json:"name"`
	Subdomain          string            `json:"subdomain"`
	CreatedOn          string            `json:"created_on"`
	ModifiedOn         string            `json:"modified_on"`
	ProductionBranch   string            `json:"production_branch"`
	Deployment_configs DeploymentConfigs `json:"deployment_configs"`
}

// DeploymentConfigs 部署配置
type DeploymentConfigs struct {
	Production json.RawMessage `json:"production"`
}
