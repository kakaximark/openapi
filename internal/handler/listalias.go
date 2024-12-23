package handler

import (
	"net/http"
	"openapi/internal/logger"
	"openapi/internal/middleware"
	"openapi/internal/service"

	"github.com/gin-gonic/gin"
)

// ListAliasRequest 获取服务别名请求结构体
type ListAliasRequest struct {
	ServiceName string `json:"ServiceName" binding:"required"`
}

// ListAliasHandler godoc
// @Summary      List aliases
// @Description  Get list of aliases for a specific service
// @Tags         aliases
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Param        Env           header  string  true  "Environment (e.g. pre, prod)"
// @Param        Country-Code  header  string  true  "Country code (e.g. US, CN)"
// @Param        request       body    ListAliasRequest  true  "List alias request"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.Response  "Invalid request headers or body"
// @Failure      401  {object}  model.Response  "Unauthorized"
// @Failure      500  {object}  model.Response  "Server error"
// @Router       /api/v1/services/aliases [post]
func ListAliasHandler(c *gin.Context) {
	var req ListAliasRequest
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

	aliases, err := service.ListAlias(headers.Env, headers.CountryCode, req.ServiceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get aliases list",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    aliases,
	})
}
