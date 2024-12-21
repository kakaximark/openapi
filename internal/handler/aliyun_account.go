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

// ListAliyunAccountsHandler godoc
// @Summary      List Aliyun accounts
// @Description  Get list of all Aliyun accounts
// @Tags         aliyun-accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Server error"
// @Router       /api/v1/aliyun/accounts [get]
func ListAliyunAccountsHandler(c *gin.Context) {
	// 从上下文中获取已验证的请求头信息
	headers := middleware.GetHeadersFromContext(c)
	if headers == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get headers from context",
		})
		return
	}

	accounts, err := service.ListAliyunAccounts()
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

// CreateAliyunAccountHandler godoc
// @Summary      Create Aliyun account
// @Description  Create a new Aliyun account
// @Tags         aliyun-accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Param        account  body      model.AliyunAccountInfo  true  "Account info"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}  "Invalid request body"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Router       /api/v1/aliyun/accounts [post]
func CreateAliyunAccountHandler(c *gin.Context) {
	// 从上下文中获取已验证的请求头信息
	headers := middleware.GetHeadersFromContext(c)
	if headers == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get headers from context",
		})
		return
	}

	var account model.AliyunAccountInfo
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if err := service.CreateAliyunAccount(&account); err != nil {
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

// UpdateAliyunAccountHandler godoc
// @Summary      Update Aliyun account
// @Description  Update an existing Aliyun account
// @Tags         aliyun-accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Param        id       path      int                     true  "Account ID"
// @Param        account  body      model.AliyunAccountInfo true  "Account info"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}  "Invalid request body"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Router       /api/v1/aliyun/accounts/{id} [put]
func UpdateAliyunAccountHandler(c *gin.Context) {
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

	var account model.AliyunAccountInfo
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if err := service.UpdateAliyunAccount(uint(id), &account); err != nil {
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

// DeleteAliyunAccountHandler godoc
// @Summary      Delete Aliyun account
// @Description  Delete an Aliyun account
// @Tags         aliyun-accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}  "Invalid request body"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Router       /api/v1/aliyun/accounts/{id} [delete]
func DeleteAliyunAccountHandler(c *gin.Context) {
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

	if err := service.DeleteAliyunAccount(uint(id)); err != nil {
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
