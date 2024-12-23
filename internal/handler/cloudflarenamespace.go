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
// @Param        headers         header  middleware.RequestHeaders  true  "Request headers"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.Response  "Invalid request headers"
// @Failure      401  {object}  model.Response  "Unauthorized"
// @Failure      500  {object}  model.Response  "Server error"
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

	logger.Info("Getting KV namespaces for CountryCode: %s, Env: %s", headers.CountryCode, headers.Env)
	resp, err := service.GetKVNamespaces(headers.CountryCode, headers.Env)
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
