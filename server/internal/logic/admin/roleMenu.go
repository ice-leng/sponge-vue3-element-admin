package admin

import (
	"admin/internal/cache"
	"admin/internal/dao"
	"admin/internal/database"
	"admin/internal/ecode"
	"admin/internal/model"
	types "admin/internal/types/admin"
	"context"
	"errors"

	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

type RoleMenuLogic interface {
	Create(ctx context.Context, request *types.CreateRoleMenuRequest) (uint64, error)
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, request *types.UpdateRoleMenuByIDRequest) error
	GetByID(ctx context.Context, id uint64) (*types.RoleMenuObjDetail, error)
	List(ctx context.Context, request *types.ListRoleMenusRequest) ([]*types.RoleMenuObjDetail, int64, error)
}

type roleMenuLogic struct {
	iDao dao.RoleMenuDao
}

func NewRoleMenuLogic() RoleMenuLogic {
	return NewRoleMenuLogicByDAO(
		dao.NewRoleMenuDao(
			database.GetDB(),
			cache.NewRoleMenuCache(database.GetCacheType()),
		),
	)
}

func NewRoleMenuLogicByDAO(iDao dao.RoleMenuDao) RoleMenuLogic {
	return &roleMenuLogic{
		iDao: iDao,
	}
}

func (l roleMenuLogic) Create(ctx context.Context, request *types.CreateRoleMenuRequest) (uint64, error) {
	roleMenu := &model.RoleMenu{}
	err := copier.Copy(roleMenu, request)
	if err != nil {
		return 0, ecode.ErrCreateRoleMenu.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	err = l.iDao.Create(ctx, roleMenu)
	return roleMenu.ID, err
}

func (l roleMenuLogic) DeleteByID(ctx context.Context, id uint64) error {
	return l.iDao.DeleteByID(ctx, id)
}

func (l roleMenuLogic) UpdateByID(ctx context.Context, request *types.UpdateRoleMenuByIDRequest) error {
	roleMenu := &model.RoleMenu{}
	err := copier.Copy(roleMenu, request)
	if err != nil {
		return ecode.ErrUpdateByIDRoleMenu.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return l.iDao.UpdateByID(ctx, roleMenu)
}

func (l roleMenuLogic) GetByID(ctx context.Context, id uint64) (*types.RoleMenuObjDetail, error) {
	roleMenu, err := l.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return nil, ecode.NotFound.Err()
		}
		return nil, err
	}

	data, err := convertRoleMenu(roleMenu)
	if err != nil {
		return nil, ecode.ErrGetByIDRoleMenu.Err()
	}

	return data, nil
}

func (l roleMenuLogic) List(ctx context.Context, request *types.ListRoleMenusRequest) ([]*types.RoleMenuObjDetail, int64, error) {
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
	if request.RoleId != nil {
		params.Columns = append(params.Columns, query.Column{
			Name:  "role_id",
			Exp:   "=",
			Value: *request.RoleId,
		})
	}

	roleMenus, total, err := l.iDao.GetByColumns(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	data, err := convertRoleMenus(roleMenus)
	if err != nil {
		return nil, 0, ecode.ErrListRoleMenu.Err()
	}

	return data, total, nil
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
