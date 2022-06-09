package middleware

import (
	"baal/lib/header"
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
		header.ContentType,
		header.Upgrade,
		header.Origin,
		header.Host,
		header.Connection,
		header.Authorization,
		header.AcceptEncoding,
		header.AcceptLanguage,
		header.AccessControlAllowOrigin,
		header.AccessControlRequestMethod,
		header.AccessControlRequestHeaders,
		header.AccessControlAllowCredentials,
		header.AccessControlMaxAge,
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
		AllowCredentials:       true,
	}

	return cors.New(conf)
}
