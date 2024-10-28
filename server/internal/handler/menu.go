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

var _ MenuHandler = (*menuHandler)(nil)

// MenuHandler defining the handler interface
type MenuHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
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
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"menu": data})
}

// List of records by query parameters
// @Summary list of menus by query parameters
// @Description list of menus by paging and conditions
// @Tags menu
// @accept json
// @Produce json
// @Param request query types.ListMenusRequest true "query parameters"
// @Success 200 {object} types.ListMenusReply{}
// @Router /api/v1/menu/list [get]
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
	menus, total, err := h.iDao.GetByParams(ctx, request)
	if err != nil {
		logger.Error("GetByParams error", logger.Err(err), logger.Any("request", request), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertMenus(menus)
	if err != nil {
		response.Error(c, ecode.ErrListMenu)
		return
	}

	response.Success(c, gin.H{
		"list":  data,
		"total": total,
	})
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

func convertMenu(menu *model.Menu) (*types.MenuObjDetail, error) {
	data := &types.MenuObjDetail{}
	err := copier.Copy(data, menu)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertMenus(fromValues []*model.Menu) ([]*types.MenuObjDetail, error) {
	toValues := []*types.MenuObjDetail{}
	for _, v := range fromValues {
		data, err := convertMenu(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}