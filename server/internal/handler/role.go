package handler

import (
	"encoding/json"
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
	iDao         dao.RoleDao
	iRoleMenuDao dao.RoleMenuDao
}

// NewRoleHandler creating the handler interface
func NewRoleHandler() RoleHandler {
	return &roleHandler{
		iDao: dao.NewRoleDao(
			model.GetDB(),
			cache.NewRoleCache(model.GetCacheType()),
		),
		iRoleMenuDao: dao.NewRoleMenuDao(
			model.GetDB(),
			cache.NewRoleMenuCache(model.GetCacheType()),
		),
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
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	role := &model.Role{}
	err = copier.Copy(role, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateRole)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, role)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": role.ID})
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
	_, id, isAbort := getRoleIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateRoleByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	role := &model.Role{}
	err = copier.Copy(role, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDRole)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, role)
	if err != nil {
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
	_, id, isAbort := getRoleIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	role, err := h.iDao.GetByID(ctx, id)
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

	data := &types.RoleObjDetail{}
	err = copier.Copy(data, role)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDRole)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

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
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	roles, total, err := h.iDao.GetByParams(ctx, request)
	if err != nil {
		logger.Error("GetByParams error", logger.Err(err), logger.Any("request", request), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertRoles(roles)
	if err != nil {
		response.Error(c, ecode.ErrListRole)
		return
	}

	response.Success(c, gin.H{
		"list":  data,
		"total": total,
	})
}

func getRoleIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertRole(role *model.Role) (*types.RoleObjDetail, error) {
	data := &types.RoleObjDetail{}
	err := copier.Copy(data, role)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertRoles(fromValues []*model.Role) ([]*types.RoleObjDetail, error) {
	toValues := []*types.RoleObjDetail{}
	for _, v := range fromValues {
		data, err := convertRole(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
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
	status := 1
	params := types.ListRolesRequest{
		Sort:   "sort",
		Status: &status,
	}
	roles, _, _ := h.iDao.GetByParams(ctx, &params)
	var options []types.Options
	for _, role := range roles {
		options = append(options, types.Options{
			Value: role.ID,
			Label: role.Name,
		})
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
	_, id, isAbort := getRoleIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}
	var menuIds []uint64
	params := types.ListRoleMenusRequest{
		RoleId: &id,
	}
	ctx := middleware.WrapCtx(c)
	roleMenus, _, _ := h.iRoleMenuDao.GetByParams(ctx, &params)

	for _, roleMenu := range roleMenus {
		menuIds = append(menuIds, roleMenu.MenuID)
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
	_, id, isAbort := getRoleIDFromPath(c)
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
	err = h.iRoleMenuDao.UpdateByRoleIds(ctx, id, menuIds)
	if err != nil {
		logger.Error("UpdateByRoleId error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}
