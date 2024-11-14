package handler

import (
	"encoding/json"
	"errors"
	"github.com/huandu/xstrings"
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
	iDao dao.MenuDao
}

// NewMenuHandler creating the handler interface
func NewMenuHandler() MenuHandler {
	return &menuHandler{
		iDao: dao.NewMenuDao(
			model.GetDB(),
			cache.NewMenuCache(model.GetCacheType()),
		),
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
// @Router /api/v1/menu [post]
// @Security BearerAuth
func (h *menuHandler) Create(c *gin.Context) {
	form := &types.CreateMenuRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	menu := &model.Menu{}
	err = copier.Copy(menu, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateMenu)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, menu)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": menu.ID})
}

// DeleteByID delete a record by id
// @Summary delete menu
// @Description delete menu by id
// @Tags menu
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteMenuByIDReply{}
// @Router /api/v1/menu/{id} [delete]
// @Security BearerAuth
func (h *menuHandler) DeleteByID(c *gin.Context) {
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
// @Summary update menu
// @Description update menu information by id
// @Tags menu
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateMenuByIDRequest true "menu information"
// @Success 200 {object} types.UpdateMenuByIDReply{}
// @Router /api/v1/menu/{id} [put]
// @Security BearerAuth
func (h *menuHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getMenuIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateMenuByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	menu := &model.Menu{}
	err = copier.Copy(menu, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDMenu)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, menu)
	if err != nil {
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
// @Router /api/v1/menu/{id} [get]
// @Security BearerAuth
func (h *menuHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getMenuIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	menu, err := h.iDao.GetByID(ctx, id)
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

	data := &types.MenuObjDetail{}
	err = copier.Copy(data, menu)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDMenu)
		return
	}
	data.RouteName = xstrings.FirstRuneToUpper(data.Path)
	// Note: if copier.Copy cannot assign a value to a field, add it here

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
// @Router /api/v1/menu [get]
// @Security BearerAuth
func (h *menuHandler) List(c *gin.Context) {
	request := &types.ListMenusRequest{}
	err := c.ShouldBindQuery(request)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	var pid uint64 = 0
	request.Sort = "id"
	if request.ParentID == nil {
		request.ParentID = &pid
	}
	menus, _, err := h.iDao.GetByParams(ctx, request)
	if err != nil {
		logger.Error("GetByParams error", logger.Err(err), logger.Any("request", request), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := h.convertMenus(c, request, menus)
	if err != nil {
		response.Error(c, ecode.ErrListMenu)
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
// @Success 200 {object} types.Result{}
// @Router /api/v1/menu/routes [get]
// @Security BearerAuth
func (h *menuHandler) Routes(c *gin.Context) {
	ctx := middleware.WrapCtx(c)
	roleId := c.GetString("roleId")
	var roleIds []uint64
	_ = json.Unmarshal([]byte(roleId), &roleIds)
	result, err := h.iDao.Routes(ctx, roleIds)
	if err != nil {
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
// @Success 200 {object} types.OptionsReply{}
// @Router /api/v1/menu/options [get]
// @Security BearerAuth
func (h *menuHandler) Options(c *gin.Context) {
	request := &types.OptionMenusRequest{}
	err := c.ShouldBindQuery(request)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	options, _ := h.iDao.Options(ctx, request)
	response.Success(c, options)
}

func getMenuIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func (h *menuHandler) convertMenu(c *gin.Context, request *types.ListMenusRequest, menu *model.Menu) (*types.MenuObjPage, error) {
	data := &types.MenuObjPage{}
	err := copier.Copy(data, menu)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.RouteName = menu.Name
	ctx := middleware.WrapCtx(c)

	request.ParentID = &menu.ID
	menus, _, err2 := h.iDao.GetByParams(ctx, request)
	if err2 != nil {
		return nil, err2
	}
	children, err3 := h.convertMenus(c, request, menus)
	if err3 != nil {
		return nil, err3
	}
	data.Children = children
	return data, nil
}

func (h *menuHandler) convertMenus(c *gin.Context, request *types.ListMenusRequest, fromValues []*model.Menu) ([]*types.MenuObjPage, error) {
	toValues := []*types.MenuObjPage{}
	for _, v := range fromValues {
		data, err := h.convertMenu(c, request, v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
