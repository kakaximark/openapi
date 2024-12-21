package handler

import (
	"net/http"

	"openapi/internal/logger"
	"openapi/internal/middleware"
	"openapi/internal/model"
	"openapi/internal/service"

	"github.com/gin-gonic/gin"
)

// LoginHandler godoc
// @Summary      User login
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body  model.LoginRequest  true  "Login credentials"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}  "Invalid request body"
// @Failure      401  {object}  map[string]interface{}  "Invalid credentials"
// @Failure      500  {object}  map[string]interface{}  "Server error"
// @Router       /api/v1/auth/login [post]
func LoginHandler(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid login request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// 验证用户
	user, err := service.ValidateUser(req.Username, req.Password)
	if err != nil {
		logger.Error("Login failed for user %s: %v", req.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Invalid username or password",
			"error":   err.Error(),
		})
		return
	}

	// 生成新token
	token, err := middleware.GenerateToken(user.Username, user.ID)
	if err != nil {
		logger.Error("Failed to generate token for user %s: %v", user.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
		return
	}

	// 更新用户会话
	if err := service.UpdateUserToken(user.ID, user.Username, token, c.ClientIP(), c.GetHeader("User-Agent")); err != nil {
		logger.Error("Failed to update user token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update user token",
			"error":   err.Error(),
		})
		return
	}

	logger.Info("User %s logged in successfully", user.Username)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data": gin.H{
			"token": token,
			// "username": user.Username,
			// "is_admin": user.IsAdmin,
		},
	})
}
