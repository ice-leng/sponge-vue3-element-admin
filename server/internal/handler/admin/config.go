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

var _ ConfigHandler = (*configHandler)(nil)

// ConfigHandler defining the handler interface
type ConfigHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	Dict(c *gin.Context)
}

type configHandler struct {
	logic logic.ConfigLogic
}

// NewConfigHandler creating the handler interface
func NewConfigHandler() ConfigHandler {
	return &configHandler{
		logic: logic.NewConfigLogic(),
	}
}

// Create a record
// @Summary create config
// @Description submit information to create config
// @Tags admin/config
// @accept json
// @Produce json
// @Param data body CreateConfigRequest true "config information"
// @Success 200 {object} CreateConfigReply{}
// @Router /admin/v1/config [post]
// @Security BearerAuth
func (h *configHandler) Create(c *gin.Context) {
	form := &types.CreateConfigRequest{}
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
// @Summary delete config
// @Description delete config by id
// @Tags admin/config
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} DeleteConfigByIDReply{}
// @Router /admin/v1/config/{id} [delete]
// @Security BearerAuth
func (h *configHandler) DeleteByID(c *gin.Context) {
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
// @Summary update config
// @Description update config information by id
// @Tags admin/config
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body UpdateConfigByIDRequest true "config information"
// @Success 200 {object} UpdateConfigByIDReply{}
// @Router /admin/v1/config/{id} [put]
// @Security BearerAuth
func (h *configHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := handlerfunc.GetIdFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateConfigByIDRequest{}
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
// @Summary get config detail
// @Description get config detail by id
// @Tags admin/config
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} GetConfigByIDReply{}
// @Router /admin/v1/config/{id} [get]
// @Security BearerAuth
func (h *configHandler) GetByID(c *gin.Context) {
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
// @Summary list of configs by query parameters
// @Description list of configs by paging and conditions
// @Tags admin/config
// @accept json
// @Produce json
// @Param request query ListConfigsRequest true "query parameters"
// @Success 200 {object} ListConfigsReply{}
// @Router /admin/v1/config [get]
// @Security BearerAuth
func (h *configHandler) List(c *gin.Context) {
	request := &types.ListConfigsRequest{}
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

// Dict 字典
// @Summary get dict
// @Description get dict
// @Tags admin/config
// @Accept json
// @Produce json
// @Success 200 {object} GetConfigByIDReply{}
// @Router /admin/v1/config/dict [get]
// @Security BearerAuth
func (h *configHandler) Dict(c *gin.Context) {
	ctx := middleware.WrapCtx(c)
	response.Success(c, h.logic.Dict(ctx))
}
