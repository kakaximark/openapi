package handler

import (
	"net/http"

	"openapi/internal/logger"

	"github.com/gin-gonic/gin"
)

// validateCloudflareRequest 验证 Cloudflare 请求的通用参数
func validateCloudflareRequest(c *gin.Context) (CountryCode, Env, NamespaceId, KeyName string, ok bool) {
	// 获取路径参数
	CountryCode = c.GetHeader("CountryCode")
	Env = c.GetHeader("Env")
	NamespaceId = c.GetHeader("NamespaceId")
	KeyName = c.GetHeader("KeyName")
	if CountryCode == "" || Env == "" || NamespaceId == "" || KeyName == "" {
		logger.Error("CountryCode or Env or NamespaceId or KeyName is required")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "CountryCode or Env or NamespaceId or KeyName is required",
		})
		return
	}

	ok = true
	return
}

// handleCloudflareError 处理 Cloudflare 错误响应
func handleCloudflareError(c *gin.Context, operation string, err error) {
	logger.Error("Failed to %s: %v", operation, err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    500,
		"message": "Failed to " + operation,
		"error":   err.Error(),
	})
}
