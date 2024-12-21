package handler

import (
	"net/http"
	"openapi/internal/logger"
	"openapi/internal/middleware"
	"openapi/internal/service"

	"github.com/gin-gonic/gin"
)

type GetKVKeysRequest struct {
	NamespaceId string `json:"namespaceId" binding:"required"`
}

// GetKVKeysHandler godoc
// @Summary      Get KV namespace keys
// @Description  Get list of keys in a Cloudflare KV namespace
// @Tags         cloudflare-kv
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Param        Env           header  string  true  "Environment (e.g. pre, prod)"
// @Param        Country-Code  header  string  true  "Country code (e.g. US, CN)"
// @Param        request       body    GetKVKeysRequest  true  "Get KV keys request"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.Response  "Invalid request headers or body"
// @Failure      401  {object}  model.Response  "Unauthorized"
// @Failure      500  {object}  model.Response  "Server error"
// @Router       /api/v1/cloudflare/kv/namespaces/keys [post]
func GetKVKeysHandler(c *gin.Context) {
	var req GetKVKeysRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}
	// 从上下文中获取已验证的请求头信息
	headers := middleware.GetHeadersFromContext(c)
	if headers == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get headers from context",
		})
		return
	}

	logger.Info("Getting KV keys for CountryCode: %s, Env: %s, Namespace: %s",
		headers.CountryCode, headers.Env, req.NamespaceId)

	resp, err := service.GetKVKeys(headers.CountryCode, headers.Env, req.NamespaceId)
	if err != nil {
		logger.Error("Failed to get KV keys: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get KV keys",
			"error":   err.Error(),
		})
		return
	}

	logger.Info("Successfully retrieved KV keys")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    resp,
	})
}
