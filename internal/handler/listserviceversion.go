package handler

import (
	"net/http"
	"openapi/internal/logger"
	"openapi/internal/middleware"
	"openapi/internal/service"

	"github.com/gin-gonic/gin"
)

// ListServiceVersionRequest 获取服务版本请求结构体
type ListServiceVersionRequest struct {
	ServiceName string `json:"ServiceName" binding:"required"`
}

// ListServiceVersionHandler godoc
// @Summary      List service versions
// @Description  Get list of versions for a specific service
// @Tags         services
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Param        Env           header  string  true  "Environment (e.g. pre, prod)"
// @Param        Country-Code  header  string  true  "Country code (e.g. US, CN)"
// @Param        request       body    ListServiceVersionRequest  true  "List service version request"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.Response  "Invalid request headers or body"
// @Failure      401  {object}  model.Response  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Server error"
// @Router       /api/v1/services/versions [post]
func ListServiceVersionHandler(c *gin.Context) {
	// 从上下文中获取已验证的请求头信息
	headers := middleware.GetHeadersFromContext(c)
	if headers == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get headers from context",
		})
		return
	}

	var req ListServiceVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	settings, err := service.ListServiceVersion(headers.Env, headers.CountryCode, req.ServiceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get services list",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    settings,
	})
}
