package dao

import (
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

var _ PlatformDao = (*platformDao)(nil)

// PlatformDao defining the dao interface
type PlatformDao interface {
	Create(ctx context.Context, table *model.Platform) error
	DeleteByID(ctx context.Context, id uint64) error
	DeleteByIDs(ctx context.Context, ids []uint64) error
	UpdateByID(ctx context.Context, table *model.Platform) error
	GetByID(ctx context.Context, id uint64) (*model.Platform, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Platform, int64, error)
	GetByParams(ctx context.Context, params *types.ListPlatformsRequest) ([]*model.Platform, int64, error)
	GetByUsername(ctx context.Context, username string) (*model.Platform, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Platform) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Platform) error
}

type platformDao struct {
	db    *gorm.DB
	cache cache.PlatformCache // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewPlatformDao creating the dao interface
func NewPlatformDao(db *gorm.DB, xCache cache.PlatformCache) PlatformDao {
	if xCache == nil {
		return &platformDao{db: db}
	}
	return &platformDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *platformDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *platformDao) Create(ctx context.Context, table *model.Platform) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a record by id
func (d *platformDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Platform{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// DeleteByIDs delete records by batch id
func (d *platformDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.Platform{}).Error
	if err != nil {
		return err
	}

	// delete cache
	for _, id := range ids {
		_ = d.deleteCache(ctx, id)
	}

	return nil
}

// UpdateByID update a record by id
func (d *platformDao) UpdateByID(ctx context.Context, table *model.Platform) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *platformDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Platform) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.Username != "" {
		update["username"] = table.Username
	}
	if table.Password != "" {
		update["password"] = table.Password
	}
	if table.Avatar != "" {
		update["avatar"] = table.Avatar
	}
	if table.RoleID != nil {
		update["role_id"] = table.RoleID
	}
	if table.Status != nil {
		update["status"] = table.Status
	}
	if table.LastTime != nil {
		update["last_time"] = table.LastTime
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *platformDao) GetByID(ctx context.Context, id uint64) (*model.Platform, error) {
	// no cache
	if d.cache == nil {
		record := &model.Platform{}
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
			table := &model.Platform{}
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
			err = d.cache.Set(ctx, id, table, cache.PlatformExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Platform)
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

func (d *platformDao) GetByUsername(ctx context.Context, username string) (*model.Platform, error) {
	record := &model.Platform{}
	err := d.db.WithContext(ctx).Where("username = ?", username).First(record).Error
	return record, err
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
func (d *platformDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Platform, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Platform{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Platform{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

func (d *platformDao) GetByParams(ctx context.Context, request *types.ListPlatformsRequest) ([]*model.Platform, int64, error) {
	page := query.NewPage(request.Page-1, request.PageSize, request.Sort)

	db := d.db.WithContext(ctx).Model(&model.Platform{}).Order(page.Sort())
	if request.StartTime != "" && request.EndTime != "" {
		db = db.Where("created_at BETWEEN ? AND ?", request.StartTime, request.EndTime+" 23:59:59")
	}

	if request.Username != "" {
		db = db.Where("username like ?", "%"+request.Username+"%")
	}
	if request.Status != nil {
		db = db.Where("status = ?", request.Status)
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

	records := []*model.Platform{}
	err := db.Find(&records).Error
	if err != nil {
		return nil, 0, err
	}
	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *platformDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Platform) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *platformDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Platform{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *platformDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Platform) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
