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

var _ PlatformHandler = (*platformHandler)(nil)

// PlatformHandler defining the handler interface
type PlatformHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type platformHandler struct {
	iDao dao.PlatformDao
}

// NewPlatformHandler creating the handler interface
func NewPlatformHandler() PlatformHandler {
	return &platformHandler{
		iDao: dao.NewPlatformDao(
			model.GetDB(),
			cache.NewPlatformCache(model.GetCacheType()),
		),
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
// @Router /api/v1/platform [post]
// @Security BearerAuth
func (h *platformHandler) Create(c *gin.Context) {
	form := &types.CreatePlatformRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	platform := &model.Platform{}
	err = copier.Copy(platform, form)
	if err != nil {
		response.Error(c, ecode.ErrCreatePlatform)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, platform)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": platform.ID})
}

// DeleteByID delete a record by id
// @Summary delete platform
// @Description delete platform by id
// @Tags platform
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeletePlatformByIDReply{}
// @Router /api/v1/platform/{id} [delete]
// @Security BearerAuth
func (h *platformHandler) DeleteByID(c *gin.Context) {
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
// @Summary update platform
// @Description update platform information by id
// @Tags platform
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdatePlatformByIDRequest true "platform information"
// @Success 200 {object} types.UpdatePlatformByIDReply{}
// @Router /api/v1/platform/{id} [put]
// @Security BearerAuth
func (h *platformHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getPlatformIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdatePlatformByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	platform := &model.Platform{}
	err = copier.Copy(platform, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDPlatform)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, platform)
	if err != nil {
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
// @Router /api/v1/platform/{id} [get]
// @Security BearerAuth
func (h *platformHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getPlatformIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	platform, err := h.iDao.GetByID(ctx, id)
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

	data := &types.PlatformObjDetail{}
	err = copier.Copy(data, platform)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDPlatform)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"platform": data})
}

// List of records by query parameters
// @Summary list of platforms by query parameters
// @Description list of platforms by paging and conditions
// @Tags platform
// @accept json
// @Produce json
// @Param request query types.ListPlatformsRequest true "query parameters"
// @Success 200 {object} types.ListPlatformsReply{}
// @Router /api/v1/platform/list [get]
// @Security BearerAuth
func (h *platformHandler) List(c *gin.Context) {
	request := &types.ListPlatformsRequest{}
	err := c.ShouldBindQuery(request)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	platforms, total, err := h.iDao.GetByParams(ctx, request)
	if err != nil {
		logger.Error("GetByParams error", logger.Err(err), logger.Any("request", request), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertPlatforms(platforms)
	if err != nil {
		response.Error(c, ecode.ErrListPlatform)
		return
	}

	response.Success(c, gin.H{
		"list":  data,
		"total": total,
	})
}

func getPlatformIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertPlatform(platform *model.Platform) (*types.PlatformObjDetail, error) {
	data := &types.PlatformObjDetail{}
	err := copier.Copy(data, platform)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertPlatforms(fromValues []*model.Platform) ([]*types.PlatformObjDetail, error) {
	toValues := []*types.PlatformObjDetail{}
	for _, v := range fromValues {
		data, err := convertPlatform(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
