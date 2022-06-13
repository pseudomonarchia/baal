package router

import (
	"baal/controller"

	"github.com/gin-gonic/gin"
)

// SetupBasicAPI ...
func SetupBasicAPI(r *gin.Engine, c *controller.Controllers) {
	r.GET("/health", c.Health.Check)
}
