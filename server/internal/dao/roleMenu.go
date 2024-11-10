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

var _ RoleMenuDao = (*roleMenuDao)(nil)

// RoleMenuDao defining the dao interface
type RoleMenuDao interface {
	Create(ctx context.Context, table *model.RoleMenu) error
	DeleteByID(ctx context.Context, id uint64) error
	DeleteByIDs(ctx context.Context, ids []uint64) error
	UpdateByID(ctx context.Context, table *model.RoleMenu) error
	GetByID(ctx context.Context, id uint64) (*model.RoleMenu, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.RoleMenu, int64, error)
	GetByParams(ctx context.Context, params *types.ListRoleMenusRequest) ([]*model.RoleMenu, int64, error)
	UpdateByRoleIds(ctx context.Context, roleId uint64, menuIds []uint64) error

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.RoleMenu) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.RoleMenu) error
}

type roleMenuDao struct {
	db    *gorm.DB
	cache cache.RoleMenuCache // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewRoleMenuDao creating the dao interface
func NewRoleMenuDao(db *gorm.DB, xCache cache.RoleMenuCache) RoleMenuDao {
	if xCache == nil {
		return &roleMenuDao{db: db}
	}
	return &roleMenuDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *roleMenuDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *roleMenuDao) Create(ctx context.Context, table *model.RoleMenu) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a record by id
func (d *roleMenuDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.RoleMenu{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// DeleteByIDs delete records by batch id
func (d *roleMenuDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.RoleMenu{}).Error
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
func (d *roleMenuDao) UpdateByID(ctx context.Context, table *model.RoleMenu) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *roleMenuDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.RoleMenu) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.RoleID != 0 {
		update["role_id"] = table.RoleID
	}
	if table.MenuID != 0 {
		update["menu_id"] = table.MenuID
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *roleMenuDao) GetByID(ctx context.Context, id uint64) (*model.RoleMenu, error) {
	// no cache
	if d.cache == nil {
		record := &model.RoleMenu{}
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
			table := &model.RoleMenu{}
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
			err = d.cache.Set(ctx, id, table, cache.RoleMenuExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.RoleMenu)
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
func (d *roleMenuDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.RoleMenu, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.RoleMenu{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.RoleMenu{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

func (d *roleMenuDao) GetByParams(ctx context.Context, request *types.ListRoleMenusRequest) ([]*model.RoleMenu, int64, error) {
	page := query.NewPage(request.Page-1, request.PageSize, request.Sort)

	db := d.db.WithContext(ctx).Model(&model.RoleMenu{}).Order(page.Sort())
	if request.StartTime != "" && request.EndTime != "" {
		db = db.Where("created_at BETWEEN ? AND ?", request.StartTime, request.EndTime)
	}

	if request.RoleId != nil {
		db = db.Where("role_id = ?", request.RoleId)
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

	records := []*model.RoleMenu{}
	err := db.Find(&records).Error
	if err != nil {
		return nil, 0, err
	}
	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *roleMenuDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.RoleMenu) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *roleMenuDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Unscoped().Delete(&model.RoleMenu{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *roleMenuDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.RoleMenu) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *roleMenuDao) UpdateByRoleIds(ctx context.Context, roleId uint64, menuIds []uint64) error {
	tx := d.db.Begin()

	if err := d.DeleteByTx(ctx, tx, roleId); err != nil {
		tx.Rollback()
		return err
	}

	var items []model.RoleMenu
	for _, menuId := range menuIds {
		items = append(items, model.RoleMenu{
			RoleID: roleId,
			MenuID: menuId,
		})
	}
	if err := tx.WithContext(ctx).Create(items).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
