package admin

import (
	"admin/internal/ecode"
	logic "admin/internal/logic/admin"
	types "admin/internal/types/admin"
	"admin/internal/types/common"
	"admin/pkg/gin/handlerfunc"
	"admin/pkg/gin/validator"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
)

var _ MenuHandler = (*menuHandler)(nil)

// MenuHandler defining the handler interface
type MenuHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	Routes(c *gin.Context)
	Options(c *gin.Context)
}

type menuHandler struct {
	logic logic.MenuLogic
}

// NewMenuHandler creating the handler interface
func NewMenuHandler() MenuHandler {
	return &menuHandler{
		logic: logic.NewMenuLogic(),
	}
}

// Create a record
// @Summary create menu
// @Description submit information to create menu
// @Tags menu
// @accept json
// @Produce json
// @Param data body types.CreateMenuRequest true "menu information"
// @Success 200 {object} types.CreateMenuReply{}
// @Router /admin/v1/menu [post]
// @Security BearerAuth
func (h *menuHandler) Create(c *gin.Context) {
	form := &types.CreateMenuRequest{}
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
// @Summary delete menu
// @Description delete menu by id
// @Tags menu
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteMenuByIDReply{}
// @Router /admin/v1/menu/{id} [delete]
// @Security BearerAuth
func (h *menuHandler) DeleteByID(c *gin.Context) {
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
// @Summary update menu
// @Description update menu information by id
// @Tags menu
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateMenuByIDRequest true "menu information"
// @Success 200 {object} types.UpdateMenuByIDReply{}
// @Router /admin/v1/menu/{id} [put]
// @Security BearerAuth
func (h *menuHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := handlerfunc.GetIdFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateMenuByIDRequest{}
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
// @Summary get menu detail
// @Description get menu detail by id
// @Tags menu
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetMenuByIDReply{}
// @Router /admin/v1/menu/{id} [get]
// @Security BearerAuth
func (h *menuHandler) GetByID(c *gin.Context) {
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
// @Summary list of menus by query parameters
// @Description list of menus by paging and conditions
// @Tags menu
// @accept json
// @Produce json
// @Param request query types.ListMenusRequest true "query parameters"
// @Success 200 {object} types.ListMenusReply{}
// @Router /admin/v1/menu [get]
// @Security BearerAuth
func (h *menuHandler) List(c *gin.Context) {
	request := &types.ListMenusRequest{}
	err := c.ShouldBindQuery(request)
	if err != nil {
		response.Error(c, ecode.InvalidParams.RewriteMsg(validator.GetValidatorErrorMsg(err)))
		return
	}

	var pid uint64 = 0
	request.Sort = "id"
	if request.ParentID == nil {
		request.ParentID = &pid
	}

	ctx := middleware.WrapCtx(c)
	data, err := h.logic.List(ctx, request)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("List error", logger.Err(err), logger.Any("request", request), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, data)
}

// Routes of records routes
// @Summary list of routes
// @Description list routes
// @Tags menu
// @accept json
// @Produce json
// @Success 200 {object} common.Result{}
// @Router /admin/v1/menu/routes [get]
// @Security BearerAuth
func (h *menuHandler) Routes(c *gin.Context) {
	ctx := middleware.WrapCtx(c)
	roleIds, _ := c.Get("roleId")
	result, err := h.logic.Routes(ctx, roleIds.(common.LocalIntArray))
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("Routes error", logger.Err(err), logger.Any("roleIds", roleIds), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	response.Success(c, result)
}

// Options get role options
// @Summary get role options
// @Description get role options
// @Tags role
// @Accept json
// @Produce json
// @Param request query types.OptionMenusRequest true "query parameters"
// @Success 200 {object} common.OptionsReply{}
// @Router /admin/v1/menu/options [get]
// @Security BearerAuth
func (h *menuHandler) Options(c *gin.Context) {
	request := &types.OptionMenusRequest{}
	err := c.ShouldBindQuery(request)
	if err != nil {
		response.Error(c, ecode.InvalidParams.RewriteMsg(validator.GetValidatorErrorMsg(err)))
		return
	}

	ctx := middleware.WrapCtx(c)
	options, err := h.logic.Options(ctx, request)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("Options error", logger.Err(err), logger.Any("request", request), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	response.Success(c, options)
}
