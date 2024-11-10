package handler

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/utils"

	"admin/internal/cache"
	"admin/internal/dao"
	"admin/internal/ecode"
	"admin/internal/model"
	"admin/internal/types"
)

var _ ConfigHandler = (*configHandler)(nil)

// ConfigHandler defining the handler interface
type ConfigHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type configHandler struct {
	iDao dao.ConfigDao
}

// NewConfigHandler creating the handler interface
func NewConfigHandler() ConfigHandler {
	return &configHandler{
		iDao: dao.NewConfigDao(
			model.GetDB(),
			cache.NewConfigCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create config
// @Description submit information to create config
// @Tags config
// @accept json
// @Produce json
// @Param data body types.CreateConfigRequest true "config information"
// @Success 200 {object} types.CreateConfigReply{}
// @Router /api/v1/config [post]
// @Security BearerAuth
func (h *configHandler) Create(c *gin.Context) {
	form := &types.CreateConfigRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	config := &model.Config{}
	err = copier.Copy(config, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateConfig)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, config)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": config.ID})
}

// DeleteByID delete a record by id
// @Summary delete config
// @Description delete config by id
// @Tags config
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteConfigByIDReply{}
// @Router /api/v1/config/{id} [delete]
// @Security BearerAuth
func (h *configHandler) DeleteByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.Error(c, ecode.InvalidParams)
		return
	}

	var ids []uint64
	for _, v := range strings.Split(idStr, ",") {
		ids = append(ids, utils.StrToUint64(v))
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteByIDs(ctx, ids)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", idStr), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update information by id
// @Summary update config
// @Description update config information by id
// @Tags config
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateConfigByIDRequest true "config information"
// @Success 200 {object} types.UpdateConfigByIDReply{}
// @Router /api/v1/config/{id} [put]
// @Security BearerAuth
func (h *configHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getConfigIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateConfigByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	config := &model.Config{}
	err = copier.Copy(config, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDConfig)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, config)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get config detail
// @Description get config detail by id
// @Tags config
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetConfigByIDReply{}
// @Router /api/v1/config/{id} [get]
// @Security BearerAuth
func (h *configHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getConfigIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	config, err := h.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.ConfigObjDetail{}
	err = copier.Copy(data, config)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDConfig)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, data)
}

// List of records by query parameters
// @Summary list of configs by query parameters
// @Description list of configs by paging and conditions
// @Tags config
// @accept json
// @Produce json
// @Param request query types.ListConfigsRequest true "query parameters"
// @Success 200 {object} types.ListConfigsReply{}
// @Router /api/v1/config [get]
// @Security BearerAuth
func (h *configHandler) List(c *gin.Context) {
	request := &types.ListConfigsRequest{}
	err := c.ShouldBindQuery(request)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	configs, total, err := h.iDao.GetByParams(ctx, request)
	if err != nil {
		logger.Error("GetByParams error", logger.Err(err), logger.Any("request", request), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertConfigs(configs)
	if err != nil {
		response.Error(c, ecode.ErrListConfig)
		return
	}

	response.Success(c, gin.H{
		"list":  data,
		"total": total,
	})
}

func getConfigIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertConfig(config *model.Config) (*types.ConfigObjDetail, error) {
	data := &types.ConfigObjDetail{}
	err := copier.Copy(data, config)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertConfigs(fromValues []*model.Config) ([]*types.ConfigObjDetail, error) {
	toValues := []*types.ConfigObjDetail{}
	for _, v := range fromValues {
		data, err := convertConfig(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
