package routers

import (
	"admin/internal/handler"
	"admin/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		roleMenuRouter(group, handler.NewRoleMenuHandler())
	})
}

func roleMenuRouter(group *gin.RouterGroup, h handler.RoleMenuHandler) {
	g := group.Group("/roleMenu")

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithVerify(fn))
	g.Use(middleware.Auth(middleware.WithExtraVerify(middlewares.VerifyToken), middleware.WithSignKey([]byte(middlewares.JwtSignKey))))

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("", h.Create)           // [post] /api/v1/roleMenu
	g.DELETE("/:id", h.DeleteByID) // [delete] /api/v1/roleMenu/:id
	g.PUT("/:id", h.UpdateByID)    // [put] /api/v1/roleMenu/:id
	g.GET("/:id", h.GetByID)       // [get] /api/v1/roleMenu/:id
	g.GET("", h.List)              // [get] /api/v1/roleMenu
}
