package test

import (
	"baal/controllers"
	"baal/routers"

	"github.com/gin-gonic/gin"
)

// MockRouter will setup & create router
func MockRouter() *gin.Engine {
	return routers.Setup(&controllers.Controllers{})
}
