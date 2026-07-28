package admin

import (
	handler "admin/internal/handler/admin"
	"admin/internal/middlewares"
	"admin/internal/routers"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware/auth"
)

func init() {
	routers.AdminV1RouterFns = append(routers.AdminV1RouterFns, func(group *gin.RouterGroup) {
		platformRouter(group, handler.NewPlatformHandler())
	})
}

func platformRouter(group *gin.RouterGroup, h handler.PlatformHandler) {
	g := group.Group("/platform")

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithVerify(fn))
	g.Use(auth.Auth(auth.WithExtraVerify(middlewares.VerifyToken)))
	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("", h.Create)                      // [post] /admin/v1/platform
	g.DELETE("/:id", h.DeleteByID)            // [delete] /admin/v1/platform/:id
	g.PUT("/:id", h.UpdateByID)               // [put] /admin/v1/platform/:id
	g.GET("/:id", h.GetByID)                  // [get] /admin/v1/platform/:id
	g.GET("", h.List)                         // [get] /admin/v1/platform
	g.GET("/me", h.Me)                        // [get] /admin/v1/platform/me
	g.GET("/profile", h.GetProfile)           // [get] /admin/v1/platform/profile
	g.PUT("/profile", h.UpdateProfile)        // [put] /admin/v1/platform/profile
	g.PUT("/password", h.ChangePassword)      // [put] /admin/v1/platform/password
	g.PUT("/password/reset", h.ResetPassword) // [put] /admin/v1/platform/password/password/reset
}
