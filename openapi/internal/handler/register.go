package handler

import (
	"errors"
	"net/http"
	"time"

	"openapi/internal/db"
	"openapi/internal/logger"
	"openapi/internal/model"

	"github.com/gin-gonic/gin"
)

// 添加注册请求结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// internal/handler/register.go
func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
		})
		return
	}

	// 验证用户名和密码长度
	if err := validateUserInput(req.Username, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	logger.Info("Attempting to register user: %s, %s", req.Username, req.Password)
	// 检查用户是否存在
	count := int64(0)
	if err := db.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		logger.Error("Database error when checking username: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Database error",
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Username already exists",
		})
		return
	}

	// 创建用户对象
	user := &model.User{
		Username:           req.Username,
		Password:           req.Password,
		LastLoginAt:        time.Now().Unix(),
		LoginAttempts:      0,
		Status:             1,
		IsAdmin:            false,
		LastLoginAttemptAt: time.Now(),
	}

	// 加密密码
	if err := user.EncryptPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to encrypt password",
		})
		return
	}

	// 保存到数据库
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "User created successfully",
	})
}

// 验证用户输入
func validateUserInput(username, password string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if len(password) < 6 || len(password) > 20 {
		return errors.New("password must be between 6 and 20 characters")
	}
	return nil
}
