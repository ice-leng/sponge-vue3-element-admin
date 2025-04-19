package routers

import (
	"admin/internal/handler"
	"admin/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		menuRouter(group, handler.NewMenuHandler())
	})
}

func menuRouter(group *gin.RouterGroup, h handler.MenuHandler) {
	g := group.Group("/menu")

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithVerify(fn))
	g.Use(middleware.Auth(middleware.WithExtraVerify(middlewares.VerifyToken), middleware.WithSignKey([]byte(middlewares.JwtSignKey))))

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("", h.Create)           // [post] /api/v1/menu
	g.DELETE("/:id", h.DeleteByID) // [delete] /api/v1/menu/:id
	g.PUT("/:id", h.UpdateByID)    // [put] /api/v1/menu/:id
	g.GET("/:id", h.GetByID)       // [get] /api/v1/menu/:id
	g.GET("", h.List)              // [get] /api/v1/menu
	g.GET("/routes", h.Routes)     // [get] /api/v1/menu/routes
	g.GET("/options", h.Options)   // [get] /api/v1/menu/options
}
