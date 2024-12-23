package handler

import (
	"net/http"
	"openapi/internal/logger"
	"openapi/internal/middleware"
	"openapi/internal/service"

	"github.com/gin-gonic/gin"
)

// GetPagesProjectHandlerInfo godoc
// @Summary      Get Pages projects
// @Description  Get list of Cloudflare Pages projects
// @Tags         cloudflare-pages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        headers         header  middleware.RequestHeaders  true  "Request headers"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.Response  "Invalid request headers"
// @Failure      401  {object}  model.Response  "Unauthorized"
// @Failure      500  {object}  model.Response  "Server error"
// @Router       /api/v1/cloudflare/pages/info [get]
func GetPagesProjectHandlerInfo(c *gin.Context) {
	// 从上下文中获取已验证的请求头信息
	headers := middleware.GetHeadersFromContext(c)
	if headers == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get headers from context",
		})
		return
	}

	logger.Info("Country-Code: %s, Env: %s", headers.CountryCode, headers.Env)
	resp, err := service.GetPagesProjectWithKVNamespacesAndKeys(headers.CountryCode, headers.Env)
	if err != nil {
		logger.Error("Failed to get pages project: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get pages project",
			"error":   err.Error(),
		})
		return
	}

	logger.Info("Successfully retrieved %d pages project", len(resp.Result))
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    resp,
	})
}
