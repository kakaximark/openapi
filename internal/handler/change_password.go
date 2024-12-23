package handler

import (
	"net/http"

	"openapi/internal/db"
	"openapi/internal/model"

	"github.com/gin-gonic/gin"
)

// ChangePasswordHandler 修改密码
func ChangePasswordHandler(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
		})
		return
	}

	// 获取当前用户
	userID := c.GetUint("user_id") // 从 JWT 中获取
	var user model.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "User not found",
		})
		return
	}

	// 验证旧密码
	if err := user.ValidatePassword(req.OldPassword); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Invalid old password",
		})
		return
	}

	// 更新新密码
	user.Password = req.NewPassword
	if err := user.EncryptPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to encrypt password",
		})
		return
	}

	// 保存到数据库
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Password updated successfully",
	})
}
