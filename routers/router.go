package routers

import (
	"baal/configs"
	"baal/controllers"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Router represents a global router struct
type Router struct {
	route *gin.Engine
	conf  *configs.GlobalConf
}

// SetupRoutes will register all routers
func SetupRoutes(c *controllers.Controllers) *gin.Engine {
	r := gin.Default()
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")

	api := r.Group("api")
	v1 := api.Group("v1")

	v1.GET("/", c.Index.Status)

	return r
}

func registration(c *controllers.Controllers, conf *configs.GlobalConf) *Router {
	gin.SetMode(conf.MODE)
	r := SetupRoutes(c)

	return &Router{r, conf}
}

// Serve will return `http.server` and service `port`
func (r *Router) Serve() (*http.Server, string) {
	port := r.conf.PORT
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        r.route,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s, port
}

// Module is used for `fx.provider` to inject dependencies
var Module fx.Option = fx.Options(fx.Provide(registration))
