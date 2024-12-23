package handler

import (
	"net/http"

	"openapi/internal/service"

	"github.com/gin-gonic/gin"
)

// LogoutHandler godoc
// @Summary      User logout
// @Description  Invalidate user's current session
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Authorization"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Server error"
// @Router       /api/v1/auth/logout [post]
func LogoutHandler(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if len(token) > 7 { // 去掉 "Bearer "
		token = token[7:]
	}

	if err := service.LogoutSession(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to logout",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Successfully logged out",
	})
}
