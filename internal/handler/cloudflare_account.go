package handler

import (
	"net/http"
	"strconv"

	"openapi/internal/logger"
	"openapi/internal/middleware"
	"openapi/internal/model"
	"openapi/internal/service"

	"github.com/gin-gonic/gin"
)

// ListCloudflareAccountsHandler godoc
// @Summary      List Cloudflare accounts
// @Description  Get list of all Cloudflare accounts
// @Tags         cloudflare-accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Server error"
// @Router       /api/v1/cloudflare/accounts [get]
func ListCloudflareAccountsHandler(c *gin.Context) {
	// 从上下文中获取已验证的请求头信息
	headers := middleware.GetHeadersFromContext(c)
	if headers == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get headers from context",
		})
		return
	}

	accounts, err := service.ListCloudflareAccounts()
	if err != nil {
		logger.Error("Failed to get accounts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get accounts",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    accounts,
	})
}

// CreateCloudflareAccountHandler godoc
// @Summary      Create Cloudflare account
// @Description  Create a new Cloudflare account
// @Tags         cloudflare-accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Param        account  body      model.CloudflareAccountInfo  true  "Account info"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}  "Invalid request body"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Router       /api/v1/cloudflare/accounts [post]
func CreateCloudflareAccountHandler(c *gin.Context) {
	// 从上下文中获取已验证的请求头信息
	headers := middleware.GetHeadersFromContext(c)
	if headers == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get headers from context",
		})
		return
	}

	var account model.CloudflareAccountInfo
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if err := service.CreateCloudflareAccount(&account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create account",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    account,
	})
}

// UpdateCloudflareAccountHandler godoc
// @Summary      Update Cloudflare account
// @Description  Update an existing Cloudflare account
// @Tags         cloudflare-accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Param        id       path      int                          true  "Account ID"
// @Param        account  body      model.CloudflareAccountInfo  true  "Account info"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}  "Invalid request body"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Router       /api/v1/cloudflare/accounts/{id} [put]
func UpdateCloudflareAccountHandler(c *gin.Context) {
	// 从上下文中获取已验证的请求头信息
	headers := middleware.GetHeadersFromContext(c)
	if headers == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get headers from context",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid account ID",
			"error":   err.Error(),
		})
		return
	}

	var account model.CloudflareAccountInfo
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if err := service.UpdateCloudflareAccount(uint(id), &account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update account",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
	})
}

// DeleteCloudflareAccountHandler godoc
// @Summary      Delete Cloudflare account
// @Description  Delete a Cloudflare account
// @Tags         cloudflare-accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}  "Invalid request body"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Router       /api/v1/cloudflare/accounts/{id} [delete]
func DeleteCloudflareAccountHandler(c *gin.Context) {
	// 从上下文中获取已验证的请求头信息
	headers := middleware.GetHeadersFromContext(c)
	if headers == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get headers from context",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid account ID",
			"error":   err.Error(),
		})
		return
	}

	if err := service.DeleteCloudflareAccount(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to delete account",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
	})
}
