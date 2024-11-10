package routers

import (
	"admin/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"

	"admin/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		platformRouter(group, handler.NewPlatformHandler())
	})
}

func platformRouter(group *gin.RouterGroup, h handler.PlatformHandler) {
	g := group.Group("/platform", middleware.Auth(middleware.WithVerify(middlewares.VerifyToken), middleware.WithSwitchHTTPCode()))

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("", h.Create)                      // [post] /api/v1/platform
	g.DELETE("/:id", h.DeleteByID)            // [delete] /api/v1/platform/:id
	g.PUT("/:id", h.UpdateByID)               // [put] /api/v1/platform/:id
	g.GET("/:id", h.GetByID)                  // [get] /api/v1/platform/:id
	g.GET("", h.List)                         // [get] /api/v1/platform
	g.GET("/me", h.Me)                        // [get] /api/v1/platform/me
	g.GET("/profile", h.GetProfile)           // [get] /api/v1/platform/profile
	g.PUT("/profile", h.UpdateProfile)        // [put] /api/v1/platform/profile
	g.PUT("/password", h.ChangePassword)      // [put] /api/v1/platform/password
	g.PUT("/password/reset", h.ResetPassword) // [put] /api/v1/platform/password/password/reset
}
