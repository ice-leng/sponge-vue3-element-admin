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

var _ RoleMenuHandler = (*roleMenuHandler)(nil)

// RoleMenuHandler defining the handler interface
type RoleMenuHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type roleMenuHandler struct {
	logic logic.RoleMenuLogic
}

// NewRoleMenuHandler creating the handler interface
func NewRoleMenuHandler() RoleMenuHandler {
	return &roleMenuHandler{
		logic: logic.NewRoleMenuLogic(),
	}
}

// Create a record
// @Summary create roleMenu
// @Description submit information to create roleMenu
// @Tags admin/roleMenu
// @accept json
// @Produce json
// @Param data body CreateRoleMenuRequest true "roleMenu information"
// @Success 200 {object} CreateRoleMenuReply{}
// @Router /admin/v1/roleMenu [post]
// @Security BearerAuth
func (h *roleMenuHandler) Create(c *gin.Context) {
	form := &types.CreateRoleMenuRequest{}
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
// @Summary delete roleMenu
// @Description delete roleMenu by id
// @Tags admin/roleMenu
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} DeleteRoleMenuByIDReply{}
// @Router /admin/v1/roleMenu/{id} [delete]
// @Security BearerAuth
func (h *roleMenuHandler) DeleteByID(c *gin.Context) {
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
// @Summary update roleMenu
// @Description update roleMenu information by id
// @Tags admin/roleMenu
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body UpdateRoleMenuByIDRequest true "roleMenu information"
// @Success 200 {object} UpdateRoleMenuByIDReply{}
// @Router /admin/v1/roleMenu/{id} [put]
// @Security BearerAuth
func (h *roleMenuHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := handlerfunc.GetIdFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateRoleMenuByIDRequest{}
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
// @Summary get roleMenu detail
// @Description get roleMenu detail by id
// @Tags admin/roleMenu
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} GetRoleMenuByIDReply{}
// @Router /admin/v1/roleMenu/{id} [get]
// @Security BearerAuth
func (h *roleMenuHandler) GetByID(c *gin.Context) {
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
// @Summary list of roleMenus by query parameters
// @Description list of roleMenus by paging and conditions
// @Tags admin/roleMenu
// @accept json
// @Produce json
// @Param request query ListRoleMenusRequest true "query parameters"
// @Success 200 {object} ListRoleMenusReply{}
// @Router /admin/v1/roleMenu [get]
// @Security BearerAuth
func (h *roleMenuHandler) List(c *gin.Context) {
	request := &types.ListRoleMenusRequest{}
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
