package handler

import (
	"net/http"

	"openapi/internal/logger"
	"openapi/internal/middleware"
	"openapi/internal/service"

	"github.com/gin-gonic/gin"
)

// GetKVNamespacesHandler godoc
// @Summary      Get KV namespaces
// @Description  Get list of KV namespaces for a Cloudflare account
// @Tags         cloudflare-kv
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Param        AccountId     path    string  true  "Account ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}  "Invalid request headers"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Server error"
// @Router       /api/v1/cloudflare/{accountId}/kv/namespaces [get]
func GetKVNamespacesHandler(c *gin.Context) {
	// 从上下文中获取已验证的请求头信息
	headers := middleware.GetHeadersFromContext(c)
	if headers == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get headers from context",
		})
		return
	}

	// 从查询参数获取账号 ID
	accountID := c.Param("accountId")
	if accountID == "" {
		logger.Error("accountId is required")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "accountId is required",
		})
		return
	}

	logger.Info("Getting KV namespaces for account: %s", accountID)
	resp, err := service.GetKVNamespaces(accountID, headers.Authorization)
	if err != nil {
		logger.Error("Failed to get KV namespaces: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get KV namespaces",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    resp,
	})
}
