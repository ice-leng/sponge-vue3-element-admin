package admin

import (
	"admin/internal/ecode"
	logic "admin/internal/logic/admin"
	types "admin/internal/types/admin"
	"admin/pkg/gin/handlerfunc"
	"admin/pkg/gin/validator"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
)

var _ AuthHandler = (*authHandler)(nil)

// AuthHandler defining the handler interface
type AuthHandler interface {
	Login(c *gin.Context)
	Captcha(c *gin.Context)
	Logout(c *gin.Context)
}

type authHandler struct {
	logic logic.AuthLogic
}

// NewAuthHandler creating the handler interface
func NewAuthHandler() AuthHandler {
	return &authHandler{
		logic: logic.NewAuthLogic(),
	}
}

// Login
// @Summary with username and password
// @Description with username and password
// @Tags admin/auth
// @Accept json
// @Produce json
// @Param data body types.LoginRequest true "login information"
// @Success 200 {object} types.LoginReply{}
// @Router /api/v1/auth/login [post]
func (h *authHandler) Login(c *gin.Context) {
	request := &types.LoginRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		response.Error(c, ecode.InvalidParams.RewriteMsg(validator.GetValidatorErrorMsg(err)))
		return
	}

	ctx := middleware.WrapCtx(c)
	data, err := h.logic.Login(ctx, request)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("Login error", logger.Err(err), logger.Any("request", request), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, data)
}

// Logout of logout
// @Summary logout
// @Description logout
// @Tags admin/auth
// @accept json
// @Produce json
// @Success 200 {object} common.Result{}
// @Router /admin/v1/auth/logout [delete]
// @Security BearerAuth
func (h *authHandler) Logout(c *gin.Context) {
	ctx := middleware.WrapCtx(c)
	err := h.logic.Logout(ctx)
	if err != nil {
		logger.Error("Logout error", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	response.Success(c)
}

// Captcha get a captcha
// @Summary get a captcha
// @Description get a captcha
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} types.CaptchaReply{}
// @Router /admin/v1/auth/captcha [get]
func (h *authHandler) Captcha(c *gin.Context) {
	ctx := middleware.WrapCtx(c)
	data, err := h.logic.Captcha(ctx)
	if err != nil {
		logger.Error("Captcha error", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	response.Success(c, data)
}
