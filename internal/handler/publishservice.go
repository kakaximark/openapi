package handler

import (
	"net/http"
	"openapi/internal/middleware"
	"openapi/internal/service"

	"openapi/internal/logger"

	"github.com/gin-gonic/gin"
)

// PublicServiceRequest 发布服务的请求体
type PublicServiceRequest struct {
	Description string `json:"description" binding:"required"`
	ServiceName string `json:"servicename" binding:"required"`
}

// PublicServiceHandler godoc
// @Summary      Publish service
// @Description  Publish a service version
// @Tags         services
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Param        headers         header  middleware.RequestHeaders  true  "Request headers"
// @Param        request       body    PublicServiceRequest  true  "Publish service request"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.Response  "Invalid request headers or body"
// @Failure      401  {object}  model.Response  "Unauthorized"
// @Failure      500  {object}  model.Response  "Server error"
// @Router       /api/v1/services/publish [post]
func PublicServiceHandler(c *gin.Context) {

	// 从上下文中获取已验证的请求头信息
	headers := middleware.GetHeadersFromContext(c)
	if headers == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get headers from context",
		})
		return
	}

	var req PublicServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	logger.Info("env: %s, countryCode: %s, serviceName: %s, description: %s", headers.Env, headers.CountryCode, req.ServiceName, req.Description)
	settings, err := service.PublicService(headers.Env, headers.CountryCode, req.ServiceName, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to publish service",
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
