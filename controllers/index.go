package controllers

import (
	"baal/libs/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index Used to store the global controller struct
type Index struct {
	log *logger.Logger
}

// Status Returns whether the server is alive when the route is accessed
func (i *Index) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"services": "alive"})
}
