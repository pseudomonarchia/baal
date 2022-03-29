package router

import (
	"baal/config"
	"baal/controller"
	"baal/middleware"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Router represents a global router struct
type Router struct {
	Route      *gin.Engine
	Controller *controller.Controllers
}

// New will register all routers
func New(c *controller.Controllers) *Router {
	if config.Global.IsDev() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(
		middleware.CorsMiddleware(),
		gin.Logger(),
		gin.Recovery(),
	)

	r.StaticFile("/favicon.ico", "./assets/favicon.ico")

	api := r.Group("api")
	v1 := api.Group("v1")

	v1.GET("/", c.Index.Status)
	v1.GET("/login", c.OAuth.LoginURL)
	v1.GET("/login/callback", c.OAuth.LoginCallBack)

	return &Router{r, c}
}

// Serve will return `http.server` and service `port`
func (r *Router) Serve(port int) *http.Server {
	return &http.Server{
		Addr:           fmt.Sprintf(":%s", strconv.Itoa(port)),
		Handler:        r.Route,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
