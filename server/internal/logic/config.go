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
	"strings"

	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

type ConfigLogic interface {
	Create(ctx context.Context, request *types.CreateConfigRequest) (uint64, error)
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, request *types.UpdateConfigByIDRequest) error
	GetByID(ctx context.Context, id uint64) (*types.ConfigObjDetail, error)
	List(ctx context.Context, request *types.ListConfigsRequest) ([]*types.ConfigObjDetail, int64, error)
	Dict(ctx context.Context) map[string]interface{}
	MakePathByConfig(ctx context.Context, path, key string) string
}

type configLogic struct {
	iDao  dao.ConfigDao
	cEnum cache.EnumCache
}

func NewConfigLogic() ConfigLogic {
	return NewConfigLogicByDAO(
		dao.NewConfigDao(
			database.GetDB(),
			cache.NewConfigCache(database.GetCacheType()),
		),
		cache.NewEnumCache(),
	)
}

func NewConfigLogicByDAO(iDao dao.ConfigDao, cEnum cache.EnumCache) ConfigLogic {
	return &configLogic{
		iDao:  iDao,
		cEnum: cEnum,
	}
}

func (l configLogic) Create(ctx context.Context, request *types.CreateConfigRequest) (uint64, error) {
	config := &model.Config{}
	err := copier.Copy(config, request)
	if err != nil {
		return 0, ecode.ErrCreateConfig.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	err = l.iDao.Create(ctx, config)
	return config.ID, err
}

func (l configLogic) DeleteByID(ctx context.Context, id uint64) error {
	return l.iDao.DeleteByID(ctx, id)
}

func (l configLogic) UpdateByID(ctx context.Context, request *types.UpdateConfigByIDRequest) error {
	config := &model.Config{}
	err := copier.Copy(config, request)
	if err != nil {
		return ecode.ErrUpdateByIDConfig.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return l.iDao.UpdateByID(ctx, config)
}

func (l configLogic) GetByID(ctx context.Context, id uint64) (*types.ConfigObjDetail, error) {
	config, err := l.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return nil, ecode.NotFound.Err()
		}
		return nil, err
	}

	data, err := convertConfig(config)
	if err != nil {
		return nil, ecode.ErrGetByIDConfig.Err()
	}

	return data, nil
}

func (l configLogic) List(ctx context.Context, request *types.ListConfigsRequest) ([]*types.ConfigObjDetail, int64, error) {
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
	if request.Name != "" {
		params.Columns = append(params.Columns, query.Column{
			Name:  "name",
			Exp:   "like",
			Value: "%" + request.Name + "%",
			Logic: "or:(",
		})
		params.Columns = append(params.Columns, query.Column{
			Name:  "`key`",
			Exp:   "like",
			Value: "%" + request.Name + "%",
			Logic: "and:)",
		})
	}
	configs, total, err := l.iDao.GetByColumns(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	data, err := convertConfigs(configs)
	if err != nil {
		return nil, 0, ecode.ErrListConfig.Err()
	}
	return data, total, nil
}

func (l configLogic) Dict(ctx context.Context) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range l.cEnum.GetAll(ctx) {
		result[k] = v
	}
	return result
}

func (l configLogic) MakePathByConfig(ctx context.Context, path, key string) string {
	if path == "" {
		return ""
	}
	if len(path) > 4 && (strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")) {
		return path
	}
	cfg, err := l.iDao.GetByKey(ctx, key)
	if err != nil || cfg == nil || cfg.Value == "" {
		return path
	}
	return strings.TrimRight(cfg.Value, "/") + "/" + strings.TrimLeft(path, "/")
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
