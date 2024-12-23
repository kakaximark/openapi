package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// HeaderConfig 定义需要验证的 header 配置
type HeaderConfig struct {
	Name     string
	Required bool
}

// RequestHeaders 存储请求头信息的结构体
type RequestHeaders struct {
	Authorization string
	Env           string
	CountryCode   string
}

// 预定义的 header 组合
var (
	// CommonHeaders 通用的 header 参数（Authorization, Env, Country-Code）
	CommonHeaders = []HeaderConfig{
		{Name: "Authorization", Required: true},
		{Name: "Env", Required: true},
		{Name: "Country-Code", Required: true},
	}

	// KVHeaders Cloudflare KV 操作需要的 header 参数
	KVHeaders = []HeaderConfig{
		{Name: "Authorization", Required: true},
		{Name: "Env", Required: true},
		{Name: "Country-Code", Required: true},
		{Name: "NamespaceId", Required: true},
		{Name: "KeyName", Required: true},
	}

	// CloudflareHeaders Cloudflare R2 操作需要的 header 参数
	CloudflareHeaders = []HeaderConfig{
		{Name: "Authorization", Required: true},
		{Name: "Env", Required: true},
		{Name: "Country-Code", Required: true},
		{Name: "BucketName", Required: true},
	}

	// AliasHeaders 别名操作需要的 header 参数
	AliasHeaders = []HeaderConfig{
		{Name: "Authorization", Required: true},
		{Name: "Env", Required: true},
		{Name: "Country-Code", Required: true},
		{Name: "AliasName", Required: true},
	}

	// PublishServiceHeaders 发布服务需要的 header 参数
	PublishServiceHeaders = []HeaderConfig{
		{Name: "Authorization", Required: true},
		{Name: "Env", Required: true},
		{Name: "Country-Code", Required: true},
		{Name: "ServiceName", Required: true},
	}
)

// ExtractToken 从 Authorization header 中提取 token
func ExtractToken(authorization string) string {
	return strings.TrimPrefix(authorization, "Bearer ")
}

// GetRequestHeaders 获取并返回请求头信息
func GetRequestHeaders(c *gin.Context) *RequestHeaders {
	return &RequestHeaders{
		Authorization: ExtractToken(c.GetHeader("Authorization")),
		Env:           c.GetHeader("Env"),
		CountryCode:   c.GetHeader("Country-Code"),
	}
}

// ValidateAndGetHeaders 验证并获取请求头信息
func ValidateAndGetHeaders(headers ...HeaderConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		missingHeaders := make([]string, 0)

		// 验证必需的 headers
		for _, header := range headers {
			if header.Required && c.GetHeader(header.Name) == "" {
				missingHeaders = append(missingHeaders, header.Name)
			}
		}

		if len(missingHeaders) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Missing required headers",
				"data": gin.H{
					"missing_headers": missingHeaders,
				},
			})
			c.Abort()
			return
		}

		// 获取请求头信息并存储到上下文中
		requestHeaders := GetRequestHeaders(c)
		c.Set("headers", requestHeaders)

		c.Next()
	}
}

// GetHeadersFromContext 从上下文中获取请求头信息
func GetHeadersFromContext(c *gin.Context) *RequestHeaders {
	headers, exists := c.Get("headers")
	if !exists {
		return nil
	}
	return headers.(*RequestHeaders)
}
