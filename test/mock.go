package test

import (
	"baal/controller"
	"baal/router"
	"baal/service"
	"baal/service/mocks"

	"github.com/gin-gonic/gin"
)

// MockSrvRoute mock new Route
func MockSrvRoute(ss ...interface{}) *gin.Engine {
	s := loopInjectionServices(ss...)
	controllers := controller.New(s)
	return router.New(controllers).Route
}

func loopInjectionServices(ss ...interface{}) *service.Services {
	services := &service.Services{
		OAuth: &mocks.OAuthFace{},
		User:  &mocks.UserFace{},
	}

	for _, s := range ss {
		switch s.(type) {
		case *mocks.OAuthFace:
			services.OAuth = s.(*mocks.OAuthFace)
			break
		case *mocks.UserFace:
			services.User = s.(*mocks.UserFace)
			break
		default:
			break
		}
	}

	return services
}
