package routers

import (
	"admin/internal/handler"
	"admin/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		uploadRouter(group, handler.NewUploadHandler())
	})
}

func uploadRouter(group *gin.RouterGroup, h handler.UploadHandler) {
	g := group.Group("/upload", middleware.Auth(middleware.WithVerify(middlewares.VerifyToken), middleware.WithSwitchHTTPCode()))

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/local", h.Local) // [post] /api/v1/upload/local
}
