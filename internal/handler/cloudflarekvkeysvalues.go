package handler

import (
	"net/http"

	"openapi/internal/logger"
	"openapi/internal/middleware"
	"openapi/internal/service"

	"github.com/gin-gonic/gin"
)

type GetKVKeyValuesRequest struct {
	NameSpaceId string `json:"namespaceid" binding:"required"`
	KeyName     string `json:"keyname" binding:"required"`
}

type UpdateKVKeyValuesRequest struct {
	NameSpaceId string `json:"namespaceid" binding:"required"`
	KeyName     string `json:"keyname" binding:"required"`
	KeyValue    string `json:"keyvalue" binding:"required"`
}

// GetKVKeyValuesHandler godoc
// @Summary      Get KV namespace key values
// @Description  Get values for keys in a Cloudflare KV namespace
// @Tags         cloudflare-kv
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        headers         header  middleware.RequestHeaders  true  "Request headers"
// @Param        request       body    GetKVKeyValuesRequest  true  "Get KV key values request"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.Response  "Invalid request headers or body"
// @Failure      401  {object}  model.Response  "Unauthorized"
// @Failure      500  {object}  model.Response  "Server error"
// @Router       /api/v1/cloudflare/kv/namespaces/keys/values [post]
func GetKVKeyValuesHandler(c *gin.Context) {
	// 从上下文中获取已验证的请求头信息
	headers := middleware.GetHeadersFromContext(c)
	if headers == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get headers from context",
		})
		return
	}

	var req GetKVKeyValuesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}
	logger.Info("Getting KV key value for CountryCode: %s, Env: %s, Namespace: %s, KeyName: %s",
		headers.CountryCode, headers.Env, req.NameSpaceId, req.KeyName)

	resp, err := service.GetKVKeyValues(headers.CountryCode, headers.Env, req.NameSpaceId, req.KeyName)
	if err != nil {
		handleCloudflareError(c, "get KV key value", err)
		return
	}

	c.String(http.StatusOK, resp.RawData)
}

// UpdateKVKeyValuesHandler godoc
// @Summary      Update KV namespace key values
// @Description  Update values for keys in a Cloudflare KV namespace
// @Tags         cloudflare-kv
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        headers         header  middleware.RequestHeaders  true  "Request headers"
// @Param        request       body    UpdateKVKeyValuesRequest  true  "Update KV key values request"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.Response  "Invalid request headers or body"
// @Failure      401  {object}  model.Response  "Unauthorized"
// @Failure      500  {object}  model.Response  "Server error"
// @Router       /api/v1/cloudflare/kv/namespaces/keys/values [put]
func UpdateKVKeyValuesHandler(c *gin.Context) {
	var req UpdateKVKeyValuesRequest
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

	logger.Info("Updating KV key value for CountryCode: %s, Env: %s, NamespaceId: %s, KeyName: %s, KeyValue: %s",
		headers.CountryCode, headers.Env, req.NameSpaceId, req.KeyName, req.KeyValue)

	resp, err := service.UpdateKVKeyValues(headers.CountryCode, headers.Env, req.NameSpaceId, req.KeyName, req.KeyValue)
	if err != nil {
		handleCloudflareError(c, "update KV key value", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    resp,
	})
}
