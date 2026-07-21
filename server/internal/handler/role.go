package handler

import (
	"admin/internal/ecode"
	"admin/internal/logic"
	"admin/internal/types"
	"admin/pkg/gin/handlerfunc"
	"admin/pkg/gin/validator"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
)

var _ RoleHandler = (*roleHandler)(nil)

// RoleHandler defining the handler interface
type RoleHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	Options(c *gin.Context)
	MenuIds(c *gin.Context)
	Menus(c *gin.Context)
}

type roleHandler struct {
	logic logic.RoleLogic
}

// NewRoleHandler creating the handler interface
func NewRoleHandler() RoleHandler {
	return &roleHandler{
		logic: logic.NewRoleLogic(),
	}
}

// Create a record
// @Summary create role
// @Description submit information to create role
// @Tags role
// @accept json
// @Produce json
// @Param data body types.CreateRoleRequest true "role information"
// @Success 200 {object} types.CreateRoleReply{}
// @Router /api/v1/role [post]
// @Security BearerAuth
func (h *roleHandler) Create(c *gin.Context) {
	form := &types.CreateRoleRequest{}
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
// @Summary delete role
// @Description delete role by id
// @Tags role
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteRoleByIDReply{}
// @Router /api/v1/role/{id} [delete]
// @Security BearerAuth
func (h *roleHandler) DeleteByID(c *gin.Context) {
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
// @Summary update role
// @Description update role information by id
// @Tags role
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateRoleByIDRequest true "role information"
// @Success 200 {object} types.UpdateRoleByIDReply{}
// @Router /api/v1/role/{id} [put]
// @Security BearerAuth
func (h *roleHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := handlerfunc.GetIdFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateRoleByIDRequest{}
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
// @Summary get role detail
// @Description get role detail by id
// @Tags role
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetRoleByIDReply{}
// @Router /api/v1/role/{id} [get]
// @Security BearerAuth
func (h *roleHandler) GetByID(c *gin.Context) {
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
// @Summary list of roles by query parameters
// @Description list of roles by paging and conditions
// @Tags role
// @accept json
// @Produce json
// @Param request query types.ListRolesRequest true "query parameters"
// @Success 200 {object} types.ListRolesReply{}
// @Router /api/v1/role [get]
// @Security BearerAuth
func (h *roleHandler) List(c *gin.Context) {
	request := &types.ListRolesRequest{}
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

// Options get role options
// @Summary get role options
// @Description get role options
// @Tags role
// @Accept json
// @Produce json
// @Success 200 {object} types.OptionsReply{}
// @Router /api/v1/role/options [get]
// @Security BearerAuth
func (h *roleHandler) Options(c *gin.Context) {
	ctx := middleware.WrapCtx(c)
	options, err := h.logic.Options(ctx)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("Options error", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	response.Success(c, options)
}

// MenuIds get role menuIds
// @Summary get role menuIds
// @Description get role menuIds
// @Tags role
// @Accept json
// @Produce json
// @Success 200 {object} types.Result{}
// @Router /api/v1/roles/{id}/menuIds [get]
// @Security BearerAuth
func (h *roleHandler) MenuIds(c *gin.Context) {
	_, id, isAbort := handlerfunc.GetIdFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	menuIds, err := h.logic.MenuIds(ctx, id)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("MenuIds error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, menuIds)
}

// Menus update permission
// @Summary update permission
// @Description update permission
// @Tags role
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.UpdateRoleByIDReply{}
// @Router /api/v1/role/{id}/menus [put]
// @Security BearerAuth
func (h *roleHandler) Menus(c *gin.Context) {
	_, id, isAbort := handlerfunc.GetIdFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	rawData, err := c.GetRawData()
	if err != nil {
		response.Error(c, ecode.InvalidParams)
		return
	}
	var menuIds []uint64
	err = json.Unmarshal(rawData, &menuIds)
	if err != nil {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err = h.logic.Menus(ctx, id, menuIds)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("Menus error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}
