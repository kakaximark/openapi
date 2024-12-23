package handler

import (
	"net/http"

	"openapi/internal/logger"
	"openapi/internal/middleware"
	"openapi/internal/service"

	"github.com/gin-gonic/gin"
)

// ListServiceHandler godoc
// @Summary      List FC services
// @Description  Get list of all Function Compute services
// @Tags         services
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Param        headers         header  middleware.RequestHeaders  true  "Request headers"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.Response  "Invalid request headers"
// @Failure      401  {object}  model.Response  "Unauthorized"
// @Failure      500  {object}  model.Response  "Server error"
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
