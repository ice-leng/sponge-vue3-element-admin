package dao

import (
	"admin/internal/pkg/util"
	"admin/internal/types"
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	cacheBase "github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/ggorm/query"
	"github.com/zhufuyi/sponge/pkg/utils"

	"admin/internal/cache"
	"admin/internal/model"
)

var _ ConfigDao = (*configDao)(nil)

// ConfigDao defining the dao interface
type ConfigDao interface {
	Create(ctx context.Context, table *model.Config) error
	DeleteByID(ctx context.Context, id uint64) error
	DeleteByIDs(ctx context.Context, ids []uint64) error
	UpdateByID(ctx context.Context, table *model.Config) error
	GetByID(ctx context.Context, id uint64) (*model.Config, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Config, int64, error)
	GetByParams(ctx context.Context, params *types.ListConfigsRequest) ([]*model.Config, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Config) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Config) error

	MakePathByConfig(ctx context.Context, path, key string) string
}

type configDao struct {
	db    *gorm.DB
	cache cache.ConfigCache   // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewConfigDao creating the dao interface
func NewConfigDao(db *gorm.DB, xCache cache.ConfigCache) ConfigDao {
	if xCache == nil {
		return &configDao{db: db}
	}
	return &configDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *configDao) deleteCache(ctx context.Context, table *model.Config) error {
	if d.cache != nil {
		err := d.cache.Del(ctx, table.ID)
		if err != nil {
			return err
		}
		return d.cache.DelByKey(ctx, table.Key)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *configDao) Create(ctx context.Context, table *model.Config) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a record by id
func (d *configDao) DeleteByID(ctx context.Context, id uint64) error {
	table, err := d.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Config{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, table)

	return nil
}

// DeleteByIDs delete records by batch id
func (d *configDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	params := &query.Params{
		Columns: []query.Column{
			{
				Name:  "id",
				Exp:   "in",
				Value: ids,
			},
		},
	}
	tables, _, _ := d.GetByColumns(ctx, params)

	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.Config{}).Error
	if err != nil {
		return err
	}

	// delete cache
	for _, table := range tables {
		_ = d.deleteCache(ctx, table)
	}

	return nil
}

// UpdateByID update a record by id
func (d *configDao) UpdateByID(ctx context.Context, table *model.Config) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table)

	return err
}

func (d *configDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Config) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.Name != "" {
		update["name"] = table.Name
	}
	if table.Description != "" {
		update["description"] = table.Description
	}
	if table.Key != "" {
		update["key"] = table.Key
	}
	if table.Value != "" {
		update["value"] = table.Value
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *configDao) GetByID(ctx context.Context, id uint64) (*model.Config, error) {
	// no cache
	if d.cache == nil {
		record := &model.Config{}
		err := d.db.WithContext(ctx).Where("id = ?", id).First(record).Error
		return record, err
	}

	// get from cache or database
	record, err := d.cache.Get(ctx, id)
	if err == nil {
		return record, nil
	}

	if errors.Is(err, model.ErrCacheNotFound) {
		// for the same id, prevent high concurrent simultaneous access to database
		val, err, _ := d.sfg.Do(utils.Uint64ToStr(id), func() (interface{}, error) { //nolint
			table := &model.Config{}
			err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
			if err != nil {
				// if data is empty, set not found cache to prevent cache penetration, default expiration time 10 minutes
				if errors.Is(err, model.ErrRecordNotFound) {
					err = d.cache.SetCacheWithNotFound(ctx, id)
					if err != nil {
						return nil, err
					}
					return nil, model.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			err = d.cache.Set(ctx, id, table, cache.ConfigExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Config)
		if !ok {
			return nil, model.ErrRecordNotFound
		}
		return table, nil
	} else if errors.Is(err, cacheBase.ErrPlaceholder) {
		return nil, model.ErrRecordNotFound
	}

	// fail fast, if cache error return, don't request to db
	return nil, err
}

// GetByKey a record by key
func (d *configDao) GetByKey(ctx context.Context, key string) (*model.Config, error) {
	// no cache
	if d.cache == nil {
		record := &model.Config{}
		err := d.db.WithContext(ctx).Where("`key` = ?", key).First(record).Error
		return record, err
	}

	// get from cache or database
	record, err := d.cache.GetByKey(ctx, key)
	if err == nil {
		return record, nil
	}

	if errors.Is(err, model.ErrCacheNotFound) {
		// for the same id, prevent high concurrent simultaneous access to database
		val, err, _ := d.sfg.Do(key, func() (interface{}, error) { //nolint
			table := &model.Config{}
			err = d.db.WithContext(ctx).Where("`key` = ?", key).First(table).Error
			if err != nil {
				// if data is empty, set not found cache to prevent cache penetration, default expiration time 10 minutes
				if errors.Is(err, model.ErrRecordNotFound) {
					err = d.cache.SetCacheByKeyWithNotFound(ctx, key)
					if err != nil {
						return nil, err
					}
					return nil, model.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			err = d.cache.SetByKey(ctx, key, table, cache.ConfigExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, key=%d", err, key)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Config)
		if !ok {
			return nil, model.ErrRecordNotFound
		}
		return table, nil
	} else if errors.Is(err, cacheBase.ErrPlaceholder) {
		return nil, model.ErrRecordNotFound
	}

	// fail fast, if cache error return, don't request to db
	return nil, err
}

// GetByColumns get paging records by column information,
// Note: query performance degrades when table rows are very large because of the use of offset.
//
// params includes paging parameters and query parameters
// paging parameters (required):
//
//	page: page number, starting from 0
//	limit: lines per page
//	sort: sort fields, default is id backwards, you can add - sign before the field to indicate reverse order, no - sign to indicate ascending order, multiple fields separated by comma
//
// query parameters (not required):
//
//	name: column name
//	exp: expressions, which default is "=",  support =, !=, >, >=, <, <=, like, in
//	value: column value, if exp=in, multiple values are separated by commas
//	logic: logical type, defaults to and when value is null, only &(and), ||(or)
//
// example: search for a male over 20 years of age
//
//	params = &query.Params{
//	    Page: 0,
//	    Limit: 20,
//	    Columns: []query.Column{
//		{
//			Name:    "age",
//			Exp: ">",
//			Value:   20,
//		},
//		{
//			Name:  "gender",
//			Value: "male",
//		},
//	}
func (d *configDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Config, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Config{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Config{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

func (d *configDao) GetByParams(ctx context.Context, request *types.ListConfigsRequest) ([]*model.Config, int64, error) {
	page := query.NewPage(request.Page-1, request.PageSize, request.Sort)

	db := d.db.WithContext(ctx).Model(&model.Config{}).Order(page.Sort())
	if request.StartTime != "" && request.EndTime != "" {
		db = db.Where("created_at BETWEEN ? AND ?", request.StartTime, request.EndTime)
	}

	var total int64 = 0
	if request.Sort != "ignore count" { // determine if count is required
		err := db.Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	if request.PageSize > 0 {
		db = db.Limit(page.Limit()).Offset(page.Page() * page.Limit())
	}

	records := []*model.Config{}
	err := db.Find(&records).Error
	if err != nil {
		return nil, 0, err
	}
	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *configDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Config) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *configDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	table, err := d.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Config{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, table)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *configDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Config) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table)

	return err
}

func (d *configDao) MakePathByConfig(ctx context.Context, path, key string) string {
	host := ""
	config, _ := d.GetByKey(ctx, key)
	if config != nil {
		host = config.Value
	}
	return util.ImageMakePath(path, host)
}
