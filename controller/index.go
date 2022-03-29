package controller

import (
	"baal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index ...
type Index struct {
	Service *service.Services
}

// Status Returns whether the server is alive when the route is accessed
func (i *Index) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"services": "alive"})
}
