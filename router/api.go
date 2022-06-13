package router

import (
	"baal/controller"

	"github.com/gin-gonic/gin"
)

// SetupAPI ...
func SetupAPI(r *gin.Engine, c *controller.Controllers) {
	api := r.Group("api")

	v1 := api.Group("v1")
	{
		v1.GET("/", c.Index.Status)
		v1.GET("/oauth", c.OAuth.LoginURL)
		v1.GET("/oauth/callback", c.OAuth.LoginCallBack)
		v1.POST("/oauth/token", c.OAuth.Token)
	}
}
