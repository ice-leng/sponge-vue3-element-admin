package routers

import (
	"admin/internal/handler"
	"admin/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware/auth"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		configRouter(group, handler.NewConfigHandler())
	})
}

func configRouter(group *gin.RouterGroup, h handler.ConfigHandler) {
	g := group.Group("/config")

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithVerify(fn))
	g.Use(auth.Auth(auth.WithExtraVerify(middlewares.VerifyToken)))

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("", h.Create)           // [post] /api/v1/config
	g.DELETE("/:id", h.DeleteByID) // [delete] /api/v1/config/:id
	g.PUT("/:id", h.UpdateByID)    // [put] /api/v1/config/:id
	g.GET("/:id", h.GetByID)       // [get] /api/v1/config/:id
	g.GET("", h.List)              // [get] /api/v1/config
	g.GET("/dict", h.Dict)         // [get] /api/v1/config/dict
}
