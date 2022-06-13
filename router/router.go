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
	r := gin.New()
	r.Use(
		middleware.CorsMiddleware(),
		gin.Logger(),
		gin.Recovery(),
	)

	if config.Global.IsDev() {
		gin.SetMode(gin.DebugMode)
		SetupSwagger(r)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	SetupBasicAPI(r, c)
	SetupAPI(r, c)

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
