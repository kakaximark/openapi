package handler

import (
	"net/http"
	"openapi/internal/middleware"
	"openapi/internal/service"

	"github.com/gin-gonic/gin"
)

// UpdateAliasRequest 更新别名的请求体
type UpdateAliasRequest struct {
	VersionId   string `json:"versionid" binding:"required"`
	ServiceName string `json:"servicename" binding:"required"`
	AliasName   string `json:"aliasname" binding:"required"`
}

// UpdateAliasHandler godoc
// @Summary      Update alias
// @Description  Update service alias to point to a specific version
// @Tags         aliases
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Param        headers         header  middleware.RequestHeaders  true  "Request headers"
// @Param        request       body    UpdateAliasRequest  true  "Update alias request"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.Response  "Invalid request headers or body"
// @Failure      401  {object}  model.Response  "Unauthorized"
// @Failure      500  {object}  model.Response  "Server error"
// @Router       /api/v1/services/aliases [put]
func UpdateAliasHandler(c *gin.Context) {

	// 从上下文中获取已验证的请求头信息
	headers := middleware.GetHeadersFromContext(c)
	if headers == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get headers from context",
		})
		return
	}

	var req UpdateAliasRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	aliasInfo, err := service.UpdateAlias(headers.Env, headers.CountryCode, req.ServiceName, req.AliasName, req.VersionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update alias",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    aliasInfo,
	})
}
