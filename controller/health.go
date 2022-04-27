package controller

import (
	"baal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health ...
type Health struct {
	Service *service.Services
}

// Check to confirm server health
func (h *Health) Check(c *gin.Context) {
	c.String(http.StatusOK, `"ok!"`)
}
