package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	methods = []string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
		"OPTIONS",
	}
	headers = []string{
		"Content-Type",
		"Upgrade",
		"Origin",
		"Connection",
		"Accept-Encoding",
		"Accept-Language",
		"Host",
		"Authorization",
	}
)

// CorsMiddleware customize cors config
func CorsMiddleware() gin.HandlerFunc {
	conf := cors.Config{
		MaxAge:                 12 * time.Hour,
		AllowBrowserExtensions: true,
		AllowAllOrigins:        true,
		AllowMethods:           methods,
		AllowHeaders:           headers,
	}

	return cors.New(conf)
}
