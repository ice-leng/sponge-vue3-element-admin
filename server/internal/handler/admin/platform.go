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

var _ PlatformHandler = (*platformHandler)(nil)

// PlatformHandler defining the handler interface
type PlatformHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	Me(c *gin.Context)
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	ChangePassword(c *gin.Context)
	ResetPassword(c *gin.Context)
}

type platformHandler struct {
	logic logic.PlatformLogic
}

// NewPlatformHandler creating the handler interface
func NewPlatformHandler() PlatformHandler {
	return &platformHandler{
		logic: logic.NewPlatformLogic(),
	}
}

// Create a record
// @Summary create platform
// @Description submit information to create platform
// @Tags platform
// @accept json
// @Produce json
// @Param data body types.CreatePlatformRequest true "platform information"
// @Success 200 {object} types.CreatePlatformReply{}
// @Router /admin/v1/platform [post]
// @Security BearerAuth
func (h *platformHandler) Create(c *gin.Context) {
	form := &types.CreatePlatformRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		response.Error(c, ecode.InvalidParams.RewriteMsg(validator.GetValidatorErrorMsg(err)))
		return
	}

	ctx := middleware.WrapCtx(c)
	id, err := h.logic.Create(ctx, form)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": id})
}

// DeleteByID delete a record by id
// @Summary delete platform
// @Description delete platform by id
// @Tags platform
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeletePlatformByIDReply{}
// @Router /admin/v1/platform/{id} [delete]
// @Security BearerAuth
func (h *platformHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := handlerfunc.GetIdFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.logic.DeleteByID(ctx, id)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update information by id
// @Summary update platform
// @Description update platform information by id
// @Tags platform
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdatePlatformByIDRequest true "platform information"
// @Success 200 {object} types.UpdatePlatformByIDReply{}
// @Router /admin/v1/platform/{id} [put]
// @Security BearerAuth
func (h *platformHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := handlerfunc.GetIdFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdatePlatformByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		response.Error(c, ecode.InvalidParams.RewriteMsg(validator.GetValidatorErrorMsg(err)))
		return
	}
	form.ID = id

	ctx := middleware.WrapCtx(c)
	err = h.logic.UpdateByID(ctx, form)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get platform detail
// @Description get platform detail by id
// @Tags platform
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetPlatformByIDReply{}
// @Router /admin/v1/platform/{id} [get]
// @Security BearerAuth
func (h *platformHandler) GetByID(c *gin.Context) {
	_, id, isAbort := handlerfunc.GetIdFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	data, err := h.logic.GetByID(ctx, id)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, data)
}

// List of records by query parameters
// @Summary list of platforms by query parameters
// @Description list of platforms by paging and conditions
// @Tags platform
// @accept json
// @Produce json
// @Param request query types.ListPlatformsRequest true "query parameters"
// @Success 200 {object} types.ListPlatformsReply{}
// @Router /admin/v1/platform [get]
// @Security BearerAuth
func (h *platformHandler) List(c *gin.Context) {
	request := &types.ListPlatformsRequest{}
	err := c.ShouldBindQuery(request)
	if err != nil {
		response.Error(c, ecode.InvalidParams.RewriteMsg(validator.GetValidatorErrorMsg(err)))
		return
	}

	ctx := middleware.WrapCtx(c)
	data, total, err := h.logic.List(ctx, request)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("List error", logger.Err(err), logger.Any("request", request), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{
		"list":  data,
		"total": total,
	})
}

// Me of records
// @Summary current information
// @Description current information
// @Tags platform
// @accept json
// @Produce json
// @Success 200 {object} types.MeReply{}
// @Router /admin/v1/platform/me [get]
// @Security BearerAuth
func (h *platformHandler) Me(c *gin.Context) {
	ctx := middleware.WrapCtx(c)
	id := handlerfunc.GetCurrentUid(c)
	data, err := h.logic.Me(ctx, id)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("Me error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	response.Success(c, data)
}

// GetProfile get me information
// @Summary current information
// @Description current information
// @Tags platform
// @accept json
// @Produce json
// @Success 200 {object} types.ProfileReply{}
// @Router /admin/v1/platform/profile [get]
// @Security BearerAuth
func (h *platformHandler) GetProfile(c *gin.Context) {
	ctx := middleware.WrapCtx(c)
	id := handlerfunc.GetCurrentUid(c)
	data, err := h.logic.Profile(ctx, id)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("Profile error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	response.Success(c, data)
}

// UpdateProfile update information by self
// @Summary update platform
// @Description update platform information by self
// @Tags platform
// @accept json
// @Produce json
// @Param data body types.UpdatePlatformByIDRequest true "platform information"
// @Success 200 {object} common.Result{}
// @Router /admin/v1/platform/profile [put]
// @Security BearerAuth
func (h *platformHandler) UpdateProfile(c *gin.Context) {
	form := &types.UpdatePlatformByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		response.Error(c, ecode.InvalidParams.RewriteMsg(validator.GetValidatorErrorMsg(err)))
		return
	}
	ctx := middleware.WrapCtx(c)
	form.ID = handlerfunc.GetCurrentUid(c)
	err = h.logic.UpdateByID(ctx, form)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("UpdateProfile error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// ChangePassword change password by self
// @Summary change password by self
// @Description change password by self
// @Tags platform
// @accept json
// @Produce json
// @Param data body types.ChangePasswordRequest true "platform information"
// @Success 200 {object} common.Result{}
// @Router /admin/v1/platform/password [put]
// @Security BearerAuth
func (h *platformHandler) ChangePassword(c *gin.Context) {
	form := &types.ChangePasswordRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		response.Error(c, ecode.InvalidParams.RewriteMsg(validator.GetValidatorErrorMsg(err)))
		return
	}

	ctx := middleware.WrapCtx(c)
	form.ID = handlerfunc.GetCurrentUid(c)
	err = h.logic.ChangePassword(ctx, form)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("ChangePassword error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	response.Success(c)
}

// ResetPassword reset password by self
// @Summary reset password by self
// @Description reset password by self
// @Tags platform
// @accept json
// @Produce json
// @Param data body types.ResetPasswordRequest true "platform information"
// @Success 200 {object} common.Result{}
// @Router /admin/v1/platform/password/reset [put]
// @Security BearerAuth
func (h *platformHandler) ResetPassword(c *gin.Context) {
	request := &types.ResetPasswordRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		response.Error(c, ecode.InvalidParams.RewriteMsg(validator.GetValidatorErrorMsg(err)))
		return
	}

	ctx := middleware.WrapCtx(c)
	request.ID = handlerfunc.GetCurrentUid(c)
	err = h.logic.ResetPassword(ctx, request)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("ResetPassword error", logger.Err(err), logger.Any("form", request), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	response.Success(c)
}
