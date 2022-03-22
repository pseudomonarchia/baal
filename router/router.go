package router

import (
	"baal/config"
	"baal/controller"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Router represents a global router struct
type Router struct {
	route *gin.Engine
	conf  *config.GlobalConf
}

var _ = (*Router)(nil)

// Module is used for `fx.provider` to inject dependencies
var Module fx.Option = fx.Options(fx.Provide(registration))

// New will register all routers
func New(c *controller.Controllers) *gin.Engine {
	r := gin.Default()
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")

	api := r.Group("api")
	v1 := api.Group("v1")

	v1.GET("/", c.Index.Status)

	return r
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

func registration(c *controller.Controllers, conf *config.GlobalConf) *Router {
	if conf.IsDev() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := New(c)
	return &Router{r, conf}
}
