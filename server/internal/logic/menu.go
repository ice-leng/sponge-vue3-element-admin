package logic

import (
	"admin/internal/cache"
	"admin/internal/dao"
	"admin/internal/database"
	"admin/internal/ecode"
	"admin/internal/model"
	"admin/internal/types"
	"context"
	"errors"

	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
	"github.com/huandu/xstrings"
)

type MenuLogic interface {
	Create(ctx context.Context, request *types.CreateMenuRequest) (uint64, error)
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, request *types.UpdateMenuByIDRequest) error
	GetByID(ctx context.Context, id uint64) (*types.MenuObjDetail, error)
	List(ctx context.Context, request *types.ListMenusRequest) ([]*types.MenuObjPage, error)
	Routes(ctx context.Context, roleIds types.LocalIntArray) ([]model.MenuItem, error)
	Options(ctx context.Context, request *types.OptionMenusRequest) ([]types.Options, error)
}

type menuLogic struct {
	iDao dao.MenuDao
}

func NewMenuLogic() MenuLogic {
	return NewMenuLogicByDAO(
		dao.NewMenuDao(
			database.GetDB(),
			cache.NewMenuCache(database.GetCacheType()),
		),
	)
}

func NewMenuLogicByDAO(iDao dao.MenuDao) MenuLogic {
	return &menuLogic{
		iDao: iDao,
	}
}

func (l menuLogic) Create(ctx context.Context, request *types.CreateMenuRequest) (uint64, error) {
	menu := &model.Menu{}
	err := copier.Copy(menu, request)
	if err != nil {
		return 0, ecode.ErrCreateMenu.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	err = l.iDao.Create(ctx, menu)
	return menu.ID, err
}

func (l menuLogic) DeleteByID(ctx context.Context, id uint64) error {
	return l.iDao.DeleteByID(ctx, id)
}

func (l menuLogic) UpdateByID(ctx context.Context, request *types.UpdateMenuByIDRequest) error {
	menu := &model.Menu{}
	err := copier.Copy(menu, request)
	if err != nil {
		return ecode.ErrUpdateByIDMenu.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return l.iDao.UpdateByID(ctx, menu)
}

func (l menuLogic) GetByID(ctx context.Context, id uint64) (*types.MenuObjDetail, error) {
	menu, err := l.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return nil, ecode.NotFound.Err()
		}
		return nil, err
	}

	data := &types.MenuObjDetail{}
	err = copier.Copy(data, menu)
	if err != nil {
		return nil, ecode.ErrGetByIDMenu.Err()
	}
	data.RouteName = xstrings.FirstRuneToUpper(data.Path)
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func (l menuLogic) List(ctx context.Context, request *types.ListMenusRequest) ([]*types.MenuObjPage, error) {
	params := &query.Params{
		Page:    request.Page - 1,
		Limit:   request.PageSize,
		Sort:    request.Sort,
		Columns: []query.Column{},
	}
	if request.StartTime != "" && request.EndTime != "" {
		params.Columns = append(params.Columns, query.Column{
			Name:  "created_at",
			Exp:   ">=",
			Value: request.StartTime,
		})
		params.Columns = append(params.Columns, query.Column{
			Name:  "created_at",
			Exp:   "<",
			Value: request.EndTime + " 23:59:59",
		})
	}
	if request.ParentID != nil {
		params.Columns = append(params.Columns, query.Column{
			Name:  "parent_id",
			Exp:   "=",
			Value: *request.ParentID,
		})
	}
	if request.Keywords != "" {
		params.Columns = append(params.Columns, query.Column{
			Name:  "name",
			Exp:   "like",
			Value: "%" + request.Keywords + "%",
		})
	}

	menus, _, err := l.iDao.GetByColumns(ctx, params)
	if err != nil {
		return nil, err
	}

	data, err := l.convertMenus(ctx, request, menus)
	if err != nil {
		return nil, ecode.ErrListMenu.Err()
	}

	return data, nil
}

func (l menuLogic) Routes(ctx context.Context, roleIds types.LocalIntArray) ([]model.MenuItem, error) {
	tops, err := l.iDao.GetListByPid(ctx, 0, roleIds, false)
	if err != nil {
		return nil, err
	}
	var items []model.MenuItem
	for _, top := range tops {
		params := top.Params
		meta := model.MenuMeta{
			Title:      top.Name,
			Icon:       top.Icon,
			Hidden:     *top.Visible != 1,
			AlwaysShow: top.AlwaysShow == 1,
			Params:     &params,
		}
		children, _ := l.getChildren(ctx, top.ID, roleIds)
		item := model.MenuItem{
			Path:      top.Path,
			Name:      top.Path,
			Component: top.Component,
			Redirect:  top.Redirect,
			Meta:      meta,
			Children:  children,
		}
		items = append(items, item)
	}
	return items, nil
}

func (l menuLogic) Options(ctx context.Context, request *types.OptionMenusRequest) ([]types.Options, error) {
	tops, err := l.iDao.GetListByPid(ctx, 0, nil, request.OnlyParent)
	if err != nil {
		return make([]types.Options, 0), err
	}

	var items []types.Options
	for _, top := range tops {
		children, _ := l.childrenOption(ctx, top.ID, request)
		item := types.Options{
			Label:    top.Name,
			Value:    top.ID,
			Children: children,
		}
		items = append(items, item)
	}
	return items, nil
}

func (l *menuLogic) convertMenu(ctx context.Context, request *types.ListMenusRequest, menu *model.Menu) (*types.MenuObjPage, error) {
	data := &types.MenuObjPage{}
	err := copier.Copy(data, menu)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.RouteName = menu.Name

	params := &query.Params{
		Page:  request.Page - 1,
		Limit: request.PageSize,
		Sort:  request.Sort,
		Columns: []query.Column{
			{
				Name:  "parent_id",
				Exp:   "=",
				Value: menu.ID,
			},
		},
	}
	if request.StartTime != "" && request.EndTime != "" {
		params.Columns = append(params.Columns, query.Column{
			Name:  "created_at",
			Exp:   ">=",
			Value: request.StartTime,
		})
		params.Columns = append(params.Columns, query.Column{
			Name:  "created_at",
			Exp:   "<",
			Value: request.EndTime + " 23:59:59",
		})
	}
	if request.Keywords != "" {
		params.Columns = append(params.Columns, query.Column{
			Name:  "name",
			Exp:   "like",
			Value: "%" + request.Keywords + "%",
		})
	}

	menus, _, err2 := l.iDao.GetByColumns(ctx, params)
	if err2 != nil {
		return nil, err2
	}
	children, err3 := l.convertMenus(ctx, request, menus)
	if err3 != nil {
		return nil, err3
	}
	data.Children = children
	return data, nil
}

func (l *menuLogic) convertMenus(ctx context.Context, request *types.ListMenusRequest, fromValues []*model.Menu) ([]*types.MenuObjPage, error) {
	toValues := []*types.MenuObjPage{}
	for _, v := range fromValues {
		data, err := l.convertMenu(ctx, request, v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}

func (l menuLogic) getChildren(ctx context.Context, pid uint64, roleIds []uint64) ([]model.Children, error) {
	items := make([]model.Children, 0)
	menus, err := l.iDao.GetListByPid(ctx, pid, roleIds, false)
	if err != nil {
		return make([]model.Children, 0), err
	}
	for _, menu := range menus {
		params := menu.Params
		meta := model.ChildrenMeta{
			Title:      menu.Name,
			Icon:       menu.Icon,
			Hidden:     *menu.Visible != 1,
			KeepAlive:  menu.KeepAlive == 1,
			AlwaysShow: menu.AlwaysShow == 1,
			Params:     &params,
		}
		item := model.Children{
			Path:      menu.Path,
			Name:      menu.Path,
			Component: menu.Component,
			Meta:      meta,
		}
		items = append(items, item)
	}
	return items, nil
}

func (l menuLogic) childrenOption(ctx context.Context, pid uint64, request *types.OptionMenusRequest) ([]types.Options, error) {
	items := make([]types.Options, 0)
	menus, err := l.iDao.GetListByPid(ctx, pid, nil, request.OnlyParent)
	if err != nil {
		return items, err
	}
	for _, menu := range menus {
		children, _ := l.childrenOption(ctx, menu.ID, request)
		item := types.Options{
			Label:    menu.Name,
			Value:    menu.ID,
			Children: children,
		}
		items = append(items, item)
	}
	return items, nil
}
