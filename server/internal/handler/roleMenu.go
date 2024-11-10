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
	iDao dao.RoleMenuDao
}

// NewRoleMenuHandler creating the handler interface
func NewRoleMenuHandler() RoleMenuHandler {
	return &roleMenuHandler{
		iDao: dao.NewRoleMenuDao(
			model.GetDB(),
			cache.NewRoleMenuCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create roleMenu
// @Description submit information to create roleMenu
// @Tags roleMenu
// @accept json
// @Produce json
// @Param data body types.CreateRoleMenuRequest true "roleMenu information"
// @Success 200 {object} types.CreateRoleMenuReply{}
// @Router /api/v1/roleMenu [post]
// @Security BearerAuth
func (h *roleMenuHandler) Create(c *gin.Context) {
	form := &types.CreateRoleMenuRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	roleMenu := &model.RoleMenu{}
	err = copier.Copy(roleMenu, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateRoleMenu)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, roleMenu)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": roleMenu.ID})
}

// DeleteByID delete a record by id
// @Summary delete roleMenu
// @Description delete roleMenu by id
// @Tags roleMenu
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteRoleMenuByIDReply{}
// @Router /api/v1/roleMenu/{id} [delete]
// @Security BearerAuth
func (h *roleMenuHandler) DeleteByID(c *gin.Context) {
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
// @Summary update roleMenu
// @Description update roleMenu information by id
// @Tags roleMenu
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateRoleMenuByIDRequest true "roleMenu information"
// @Success 200 {object} types.UpdateRoleMenuByIDReply{}
// @Router /api/v1/roleMenu/{id} [put]
// @Security BearerAuth
func (h *roleMenuHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getRoleMenuIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateRoleMenuByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	roleMenu := &model.RoleMenu{}
	err = copier.Copy(roleMenu, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDRoleMenu)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, roleMenu)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get roleMenu detail
// @Description get roleMenu detail by id
// @Tags roleMenu
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetRoleMenuByIDReply{}
// @Router /api/v1/roleMenu/{id} [get]
// @Security BearerAuth
func (h *roleMenuHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getRoleMenuIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	roleMenu, err := h.iDao.GetByID(ctx, id)
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

	data := &types.RoleMenuObjDetail{}
	err = copier.Copy(data, roleMenu)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDRoleMenu)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, data)
}

// List of records by query parameters
// @Summary list of roleMenus by query parameters
// @Description list of roleMenus by paging and conditions
// @Tags roleMenu
// @accept json
// @Produce json
// @Param request query types.ListRoleMenusRequest true "query parameters"
// @Success 200 {object} types.ListRoleMenusReply{}
// @Router /api/v1/roleMenu [get]
// @Security BearerAuth
func (h *roleMenuHandler) List(c *gin.Context) {
	request := &types.ListRoleMenusRequest{}
	err := c.ShouldBindQuery(request)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	roleMenus, total, err := h.iDao.GetByParams(ctx, request)
	if err != nil {
		logger.Error("GetByParams error", logger.Err(err), logger.Any("request", request), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertRoleMenus(roleMenus)
	if err != nil {
		response.Error(c, ecode.ErrListRoleMenu)
		return
	}

	response.Success(c, gin.H{
		"list":  data,
		"total": total,
	})
}

func getRoleMenuIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertRoleMenu(roleMenu *model.RoleMenu) (*types.RoleMenuObjDetail, error) {
	data := &types.RoleMenuObjDetail{}
	err := copier.Copy(data, roleMenu)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertRoleMenus(fromValues []*model.RoleMenu) ([]*types.RoleMenuObjDetail, error) {
	toValues := []*types.RoleMenuObjDetail{}
	for _, v := range fromValues {
		data, err := convertRoleMenu(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
