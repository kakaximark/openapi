package handler

import (
	"net/http"

	"openapi/internal/logger"
	"openapi/internal/middleware"
	"openapi/internal/service"

	"github.com/gin-gonic/gin"
)

// GetBucketRequest 获取bucket信息的请求结构
type GetBucketRequest struct {
	BucketName string `json:"bucketname" binding:"required" example:"bucket-name"`
}

// DeleteDirectoryRequest 删除目录的请求结构
type DeleteDirectoryRequest struct {
	DirPath    string `json:"dirpath" binding:"required" example:"path/to/directory"`
	BucketName string `json:"bucketname" binding:"required" example:"bucket-name"`
}

type CopyDirectoryRequest struct {
	SourceDir  string `json:"sourcedir" binding:"required" example:"path/to/source/directory"`
	TargetDir  string `json:"targetdir" binding:"required" example:"path/to/target/directory"`
	BucketName string `json:"bucketname" binding:"required" example:"bucket-name"`
}

// GetBucketHandler godoc
// @Summary      Get bucket info
// @Description  Get information about a Cloudflare R2 bucket
// @Tags         cloudflare
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Param        headers         header  middleware.RequestHeaders  true  "Request headers"
// @Param        request         body    GetBucketRequest  true  "Get bucket request"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.Response  "Invalid request headers"
// @Failure      401  {object}  model.Response  "Unauthorized"
// @Failure      500  {object}  model.Response  "Server error"
// @Router       /api/v1/cloudflare/bucketinfo [post]
func GetBucketHandler(c *gin.Context) {
	var req GetBucketRequest
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

	logger.Info("BucketName: %s, Country-Code: %s, Env: %s", req.BucketName, headers.CountryCode, headers.Env)

	resp, err := service.GetBucketInfo(headers.CountryCode, headers.Env, req.BucketName)
	if err != nil {
		logger.Error("Failed to get bucket info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get bucket info",
			"error":   err.Error(),
		})
		return
	}

	logger.Info("Successfully retrieved bucket info")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    resp,
	})
}

// CopyDirectoryHandler godoc
// @Summary      Copy directory
// @Description  Copy a directory in Cloudflare R2 bucket
// @Tags         cloudflare
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        headers         header  middleware.RequestHeaders  true  "Request headers"
// @Param        request         body    CopyDirectoryRequest  true  "Copy directory request"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.Response  "Invalid request headers or body"
// @Failure      401  {object}  model.Response  "Unauthorized"
// @Failure      500  {object}  model.Response  "Server error"
// @Router       /api/v1/cloudflare/bucketinfo/copy [post]
func CopyDirectoryHandler(c *gin.Context) {
	var req CopyDirectoryRequest
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

	logger.Info("Copying directory in bucket %s from %s to %s",
		req.BucketName, req.SourceDir, req.TargetDir)

	// 调用服务层复制目录
	err := service.CopyDirectory(headers.CountryCode, headers.Env, req.BucketName, req.SourceDir, req.TargetDir)
	if err != nil {
		logger.Error("Failed to copy directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to copy directory",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Successfully copied directory",
	})
}

// DeleteDirectoryHandler godoc
// @Summary      Delete directory
// @Description  Delete a directory in Cloudflare R2 bucket
// @Tags         cloudflare
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        headers         header  middleware.RequestHeaders  true  "Request headers"
// @Param        request         body    DeleteDirectoryRequest  true  "Delete directory request"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.Response  "Invalid request headers or body"
// @Failure      401  {object}  model.Response  "Unauthorized"
// @Failure      500  {object}  model.Response  "Server error"
// @Router       /api/v1/cloudflare/bucketinfo [delete]
func DeleteDirectoryHandler(c *gin.Context) {
	var req DeleteDirectoryRequest
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

	logger.Info("Deleting directory %s in bucket %s", req.DirPath, req.BucketName)

	// 调用服务层删除目录
	err := service.DeleteDirectory(headers.CountryCode, headers.Env, req.BucketName, req.DirPath)
	if err != nil {
		logger.Error("Failed to delete directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to delete directory",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Successfully deleted directory",
	})
}
