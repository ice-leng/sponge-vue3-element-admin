package admin

import (
	"admin/internal/cache"
	"admin/internal/dao"
	"admin/internal/database"
	"admin/internal/ecode"
	"admin/internal/model"
	types "admin/internal/types/admin"
	common "admin/internal/types/common"
	"context"
	"errors"

	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

type RoleLogic interface {
	Create(ctx context.Context, request *types.CreateRoleRequest) (uint64, error)
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, request *types.UpdateRoleByIDRequest) error
	GetByID(ctx context.Context, id uint64) (*types.RoleObjDetail, error)
	List(ctx context.Context, request *types.ListRolesRequest) ([]*types.RoleObjDetail, int64, error)
	Options(ctx context.Context) ([]common.Options, error)
	MenuIds(ctx context.Context, id uint64) ([]uint64, error)
	Menus(ctx context.Context, id uint64, menuIds []uint64) error
}

type roleLogic struct {
	iDao         dao.RoleDao
	iRoleMenuDao dao.RoleMenuDao
}

func NewRoleLogic() RoleLogic {
	return NewRoleLogicByDAO(
		dao.NewRoleDao(
			database.GetDB(),
			cache.NewRoleCache(database.GetCacheType()),
		),
		dao.NewRoleMenuDao(
			database.GetDB(),
			cache.NewRoleMenuCache(database.GetCacheType()),
		),
	)
}

func NewRoleLogicByDAO(iDao dao.RoleDao, iRoleMenuDao dao.RoleMenuDao) RoleLogic {
	return &roleLogic{
		iDao:         iDao,
		iRoleMenuDao: iRoleMenuDao,
	}
}

func (l roleLogic) Create(ctx context.Context, request *types.CreateRoleRequest) (uint64, error) {
	role := &model.Role{}
	err := copier.Copy(role, request)
	if err != nil {
		return 0, ecode.ErrCreateRole.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	err = l.iDao.Create(ctx, role)
	return role.ID, err
}

func (l roleLogic) DeleteByID(ctx context.Context, id uint64) error {
	return l.iDao.DeleteByID(ctx, id)
}

func (l roleLogic) UpdateByID(ctx context.Context, request *types.UpdateRoleByIDRequest) error {
	role := &model.Role{}
	err := copier.Copy(role, request)
	if err != nil {
		return ecode.ErrUpdateByIDRole.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return l.iDao.UpdateByID(ctx, role)
}

func (l roleLogic) GetByID(ctx context.Context, id uint64) (*types.RoleObjDetail, error) {
	role, err := l.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return nil, ecode.NotFound.Err()
		}
		return nil, err
	}

	data, err := convertRole(role)
	if err != nil {
		return nil, ecode.ErrGetByIDRole.Err()
	}

	return data, nil
}

func (l roleLogic) List(ctx context.Context, request *types.ListRolesRequest) ([]*types.RoleObjDetail, int64, error) {
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
	if request.Status != nil {
		params.Columns = append(params.Columns, query.Column{
			Name:  "status",
			Exp:   "=",
			Value: *request.Status,
		})
	}

	roles, total, err := l.iDao.GetByColumns(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	data, err := convertRoles(roles)
	if err != nil {
		return nil, 0, ecode.ErrListRole.Err()
	}

	return data, total, nil
}

func (l roleLogic) Options(ctx context.Context) ([]common.Options, error) {
	params := &query.Params{
		Page:  0,
		Limit: 1000,
		Sort:  "sort",
		Columns: []query.Column{
			{
				Name:  "status",
				Exp:   "=",
				Value: 1,
			},
		},
	}

	roles, _, _ := l.iDao.GetByColumns(ctx, params)
	var options []common.Options
	for _, role := range roles {
		options = append(options, common.Options{
			Value: role.ID,
			Label: role.Name,
		})
	}

	return options, nil
}

func (l roleLogic) MenuIds(ctx context.Context, id uint64) ([]uint64, error) {
	var menuIds []uint64
	params := &query.Params{
		Page:  0,
		Limit: 1000,
		Sort:  "id",
		Columns: []query.Column{
			{
				Name:  "role_id",
				Exp:   "=",
				Value: id,
			},
		},
	}
	roleMenus, _, _ := l.iRoleMenuDao.GetByColumns(ctx, params)
	for _, roleMenu := range roleMenus {
		menuIds = append(menuIds, roleMenu.MenuID)
	}

	return menuIds, nil
}

func (l roleLogic) Menus(ctx context.Context, id uint64, menuIds []uint64) error {
	return l.iRoleMenuDao.UpdateByRoleIds(ctx, id, menuIds)
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
