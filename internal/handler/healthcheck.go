package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheckHandler godoc
// @Summary      Health check endpoint
// @Description  Returns OK if the service is running
// @Tags         system
// @Produce      plain
// @Success      200  {string}  string  "ok"
// @Router       /healthcheck [get]
func HealthCheckHandler(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
