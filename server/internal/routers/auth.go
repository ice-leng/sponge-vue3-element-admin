package routers

import (
	"admin/internal/handler"
	"admin/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		authRouter(group, handler.NewAuthHandler())
	})
}

func authRouter(group *gin.RouterGroup, h handler.AuthHandler) {
	g := group.Group("/auth")

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/login", h.Login)                                                                                                                 // [post] /api/v1/auth/login
	g.GET("/captcha", h.Captcha)                                                                                                              // [get] /api/v1/auth/captcha
	g.POST("/refreshToken", middleware.Auth(middleware.WithVerify(middlewares.VerifyToken), middleware.WithSwitchHTTPCode()), h.RefreshToken) // [post] /api/v1/auth/refreshToken
	g.DELETE("/logout", middleware.Auth(middleware.WithVerify(middlewares.VerifyToken), middleware.WithSwitchHTTPCode()), h.Logout)           // [delete] /api/v1/auth/logout
}
