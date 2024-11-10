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

var _ MenuDao = (*menuDao)(nil)

// MenuDao defining the dao interface
type MenuDao interface {
	Create(ctx context.Context, table *model.Menu) error
	DeleteByID(ctx context.Context, id uint64) error
	DeleteByIDs(ctx context.Context, ids []uint64) error
	UpdateByID(ctx context.Context, table *model.Menu) error
	GetByID(ctx context.Context, id uint64) (*model.Menu, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Menu, int64, error)
	GetByParams(ctx context.Context, params *types.ListMenusRequest) ([]*model.Menu, int64, error)
	Routes(ctx context.Context, roleIds []uint64) ([]model.MenuItem, error)
	Options(ctx context.Context, request *types.OptionMenusRequest) ([]types.Options, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Menu) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Menu) error
}

type menuDao struct {
	db    *gorm.DB
	cache cache.MenuCache     // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewMenuDao creating the dao interface
func NewMenuDao(db *gorm.DB, xCache cache.MenuCache) MenuDao {
	if xCache == nil {
		return &menuDao{db: db}
	}
	return &menuDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *menuDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *menuDao) Create(ctx context.Context, table *model.Menu) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a record by id
func (d *menuDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Menu{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// DeleteByIDs delete records by batch id
func (d *menuDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.Menu{}).Error
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
func (d *menuDao) UpdateByID(ctx context.Context, table *model.Menu) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *menuDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Menu) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.ParentID != 0 {
		update["parent_id"] = table.ParentID
	}
	if table.Name != "" {
		update["name"] = table.Name
	}
	if table.Type != "" {
		update["type"] = table.Type
	}
	if table.Path != "" {
		update["path"] = table.Path
	}
	if table.Component != "" {
		update["component"] = table.Component
	}
	if table.Perm != "" {
		update["perm"] = table.Perm
	}
	if table.Sort != 0 {
		update["sort"] = table.Sort
	}
	if table.Visible != 0 {
		update["visible"] = table.Visible
	}
	if table.Icon != "" {
		update["icon"] = table.Icon
	}
	if table.Redirect != "" {
		update["redirect"] = table.Redirect
	}
	if table.AlwaysShow != 0 {
		update["always_show"] = table.AlwaysShow
	}
	if table.KeepAlive != 0 {
		update["keep_alive"] = table.KeepAlive
	}
	if table.Params != "" {
		update["params"] = table.Params
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *menuDao) GetByID(ctx context.Context, id uint64) (*model.Menu, error) {
	// no cache
	if d.cache == nil {
		record := &model.Menu{}
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
			table := &model.Menu{}
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
			err = d.cache.Set(ctx, id, table, cache.MenuExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Menu)
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
func (d *menuDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Menu, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Menu{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Menu{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

func (d *menuDao) GetByParams(ctx context.Context, request *types.ListMenusRequest) ([]*model.Menu, int64, error) {
	page := query.NewPage(request.Page-1, request.PageSize, request.Sort)

	db := d.db.WithContext(ctx).Model(&model.Menu{}).Order(page.Sort())
	if request.StartTime != "" && request.EndTime != "" {
		db = db.Where("created_at BETWEEN ? AND ?", request.StartTime, request.EndTime)
	}

	if request.ParentID != nil {
		db = db.Where("parent_id = ?", *request.ParentID)
	}

	if request.Keywords != "" {
		db = db.Where("name LIKE ?", "%"+request.Keywords+"%")
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

	records := []*model.Menu{}
	err := db.Find(&records).Error
	if err != nil {
		return nil, 0, err
	}
	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *menuDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Menu) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *menuDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Menu{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *menuDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Menu) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *menuDao) getListByPid(ctx context.Context, pid uint64, roleIds []uint64, onlyParent bool) ([]*model.Menu, error) {
	db := d.db.WithContext(ctx).
		Model(&model.Menu{}).
		Where("t_menu.parent_id = ?", pid).
		Order("t_menu.sort asc")
	if roleIds != nil && len(roleIds) > 0 {
		db = db.Joins("join t_role_menu on t_role_menu.menu_id = t_menu.id").
			Where("t_role_menu.role_id in ?", roleIds).
			Group("t_menu.id")
	}
	if onlyParent {
		db = db.Where("t_menu.type in ?", []string{"CATALOG", "MENU"})
	}
	var records []*model.Menu
	err := db.Find(&records).Error
	return records, err
}

func (d *menuDao) getChildren(ctx context.Context, pid uint64, roleIds []uint64) ([]model.Children, error) {
	var items []model.Children
	menus, err := d.getListByPid(ctx, pid, roleIds, false)
	if err != nil {
		return make([]model.Children, 0), err
	}
	for _, menu := range menus {
		params := menu.Params
		meta := model.ChildrenMeta{
			Title:      menu.Name,
			Icon:       menu.Icon,
			Hidden:     menu.Visible != 1,
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

func (d *menuDao) Routes(ctx context.Context, roleIds []uint64) ([]model.MenuItem, error) {
	tops, err := d.getListByPid(ctx, 0, roleIds, false)
	if err != nil {
		return nil, err
	}
	var items []model.MenuItem
	for _, top := range tops {
		params := top.Params
		meta := model.MenuMeta{
			Title:      top.Name,
			Icon:       top.Icon,
			Hidden:     top.Visible != 1,
			AlwaysShow: top.AlwaysShow == 1,
			Params:     &params,
		}
		children, _ := d.getChildren(ctx, top.ID, roleIds)
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

func (d *menuDao) childrenOption(ctx context.Context, pid uint64, request *types.OptionMenusRequest) ([]types.Options, error) {
	items := make([]types.Options, 0)
	menus, err := d.getListByPid(ctx, pid, nil, request.OnlyParent)
	if err != nil {
		return items, err
	}
	for _, menu := range menus {
		children, _ := d.childrenOption(ctx, menu.ID, request)
		item := types.Options{
			Label:    menu.Name,
			Value:    menu.ID,
			Children: children,
		}
		items = append(items, item)
	}
	return items, nil
}

func (d *menuDao) Options(ctx context.Context, request *types.OptionMenusRequest) ([]types.Options, error) {
	tops, err := d.getListByPid(ctx, 0, nil, request.OnlyParent)
	if err != nil {
		return make([]types.Options, 0), err
	}

	var items []types.Options
	for _, top := range tops {
		children, _ := d.childrenOption(ctx, top.ID, request)
		item := types.Options{
			Label:    top.Name,
			Value:    top.ID,
			Children: children,
		}
		items = append(items, item)
	}
	return items, nil
}
