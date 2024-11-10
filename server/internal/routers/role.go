package routers

import (
	"admin/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"

	"admin/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		roleRouter(group, handler.NewRoleHandler())
	})
}

func roleRouter(group *gin.RouterGroup, h handler.RoleHandler) {
	g := group.Group("/role", middleware.Auth(middleware.WithVerify(middlewares.VerifyToken), middleware.WithSwitchHTTPCode()))

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("", h.Create)             // [post] /api/v1/role
	g.DELETE("/:id", h.DeleteByID)   // [delete] /api/v1/role/:id
	g.PUT("/:id", h.UpdateByID)      // [put] /api/v1/role/:id
	g.GET("/:id", h.GetByID)         // [get] /api/v1/role/:id
	g.GET("", h.List)                // [get] /api/v1/role
	g.GET("/options", h.Options)     // [get] /api/v1/role/options
	g.GET("/:id/menuIds", h.MenuIds) // [get] /api/v1/role/:id/menuIds
	g.PUT("/:id/menus", h.Menus)     // [put] /api/v1/role/:id/menus
}
