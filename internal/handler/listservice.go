package handler

import (
	"net/http"

	"openapi/internal/logger"
	"openapi/internal/middleware"
	"openapi/internal/service"

	"github.com/gin-gonic/gin"
)

// ListServiceRequest 定义请求结构
type ListServiceRequest struct {
	Env         string `json:"env" binding:"required" example:"prod"`
	CountryCode string `json:"countryCode" binding:"required" example:"CN"`
}

// ListServiceHandler godoc
// @Summary      List FC services
// @Description  Get list of all Function Compute services
// @Tags         services
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Param        Env           header  string  true  "Environment (e.g. pre, prod)"
// @Param        Country-Code  header  string  true  "Country code (e.g. US, CN)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}  "Invalid request headers"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Server error"
// @Router       /api/v1/services [get]
func ListServiceHandler(c *gin.Context) {
	// 从上下文中获取已验证的请求头信息
	headers := middleware.GetHeadersFromContext(c)
	if headers == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get headers from context",
		})
		return
	}

	logger.Info("env: %s, countryCode: %s", headers.Env, headers.CountryCode)

	// 使用获取到的 header 参数调用服务
	services, err := service.ListService(headers.Env, headers.CountryCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get services",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    services,
	})
}
