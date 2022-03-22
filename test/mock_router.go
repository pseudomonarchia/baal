package test

import (
	"baal/controller"
	"baal/router"

	"github.com/gin-gonic/gin"
)

// MockRouter will setup & create router
func MockRouter() *gin.Engine {
	return router.New(&controller.Controllers{})
}
