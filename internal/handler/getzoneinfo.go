package handler

import (
	"net/http"
	"openapi/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetZoneInfoHandler godoc
// @Summary      Get zone information and validate token
// @Description  Validate token and get user information along with available zones
// @Tags         system
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string  true  "Bearer {token}"
// @Success      200  {object}  map[string]interface{} "Response with userid, username, environment, and country_code"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized or invalid token"
// @Failure      500  {object}  map[string]interface{}  "Server error"
// @Router       /api/v1/system/zones [get]
func GetZoneInfoHandler(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Authorization header is required",
		})
		return
	}

	if !strings.HasPrefix(token, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Invalid authorization format",
		})
		return
	}

	token = token[7:] // 去掉 "Bearer "

	userid, username, err := service.GetZoneInfo(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Invalid or expired token",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"userid":       userid,
			"username":     username,
			"environment":  "pre",
			"country_code": []string{"US", "BZ"},
		},
		"message": "Success",
	})
}
