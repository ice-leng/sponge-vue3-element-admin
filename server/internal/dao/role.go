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

var _ RoleDao = (*roleDao)(nil)

// RoleDao defining the dao interface
type RoleDao interface {
	Create(ctx context.Context, table *model.Role) error
	DeleteByID(ctx context.Context, id uint64) error
	DeleteByIDs(ctx context.Context, ids []uint64) error
	UpdateByID(ctx context.Context, table *model.Role) error
	GetByID(ctx context.Context, id uint64) (*model.Role, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.Role, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Role, int64, error)
	GetByParams(ctx context.Context, params *types.ListRolesRequest) ([]*model.Role, int64, error)
	GetPermissionsByIds(ctx context.Context, ids []uint64) ([]string, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Role) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Role) error
}

type roleDao struct {
	db    *gorm.DB
	cache cache.RoleCache     // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewRoleDao creating the dao interface
func NewRoleDao(db *gorm.DB, xCache cache.RoleCache) RoleDao {
	if xCache == nil {
		return &roleDao{db: db}
	}
	return &roleDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *roleDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *roleDao) Create(ctx context.Context, table *model.Role) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a record by id
func (d *roleDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Role{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// DeleteByIDs delete records by batch id
func (d *roleDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.Role{}).Error
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
func (d *roleDao) UpdateByID(ctx context.Context, table *model.Role) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *roleDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Role) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.Name != "" {
		update["name"] = table.Name
	}
	if table.Code != "" {
		update["code"] = table.Code
	}
	if table.Sort != 0 {
		update["sort"] = table.Sort
	}
	if table.Status != 0 {
		update["status"] = table.Status
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *roleDao) GetByID(ctx context.Context, id uint64) (*model.Role, error) {
	// no cache
	if d.cache == nil {
		record := &model.Role{}
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
			table := &model.Role{}
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
			err = d.cache.Set(ctx, id, table, cache.RoleExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Role)
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

// GetByIDs get records by batch id
func (d *roleDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.Role, error) {
	// no cache
	if d.cache == nil {
		var records []*model.Role
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		itemMap := make(map[uint64]*model.Role)
		for _, record := range records {
			itemMap[record.ID] = record
		}
		return itemMap, nil
	}

	// get form cache or database
	itemMap, err := d.cache.MultiGet(ctx, ids)
	if err != nil {
		return nil, err
	}

	var missedIDs []uint64
	for _, id := range ids {
		_, ok := itemMap[id]
		if !ok {
			missedIDs = append(missedIDs, id)
			continue
		}
	}

	// get missed data
	if len(missedIDs) > 0 {
		// find the id of an active placeholder, i.e. an id that does not exist in database
		var realMissedIDs []uint64
		for _, id := range missedIDs {
			_, err = d.cache.Get(ctx, id)
			if errors.Is(err, cacheBase.ErrPlaceholder) {
				continue
			}
			realMissedIDs = append(realMissedIDs, id)
		}

		if len(realMissedIDs) > 0 {
			var missedData []*model.Role
			err = d.db.WithContext(ctx).Where("id IN (?)", realMissedIDs).Find(&missedData).Error
			if err != nil {
				return nil, err
			}

			if len(missedData) > 0 {
				for _, data := range missedData {
					itemMap[data.ID] = data
				}
				err = d.cache.MultiSet(ctx, missedData, cache.RoleExpireTime)
				if err != nil {
					return nil, err
				}
			} else {
				for _, id := range realMissedIDs {
					_ = d.cache.SetCacheWithNotFound(ctx, id)
				}
			}
		}
	}

	return itemMap, nil
}

func (d *roleDao) GetPermissionsByIds(ctx context.Context, ids []uint64) ([]string, error) {
	var perms []string
	err := d.db.WithContext(ctx).
		Model(&model.Role{}).
		Joins("LEFT JOIN t_role_menu as role_menu ON role_menu.role_id = t_role.id").
		Joins("LEFT JOIN t_menu as menu ON menu.id = role_menu.menu_id").
		Where("menu.perm != ?", "").
		Where("t_role.id IN (?)", ids).
		Pluck("menu.perm", &perms).
		Error
	return perms, err
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
func (d *roleDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Role, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Role{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Role{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

func (d *roleDao) GetByParams(ctx context.Context, request *types.ListRolesRequest) ([]*model.Role, int64, error) {
	page := query.NewPage(request.Page-1, request.PageSize, request.Sort)

	db := d.db.WithContext(ctx).Model(&model.Role{}).Order(page.Sort())
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

	records := []*model.Role{}
	err := db.Find(&records).Error
	if err != nil {
		return nil, 0, err
	}
	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *roleDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Role) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *roleDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Role{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *roleDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Role) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
