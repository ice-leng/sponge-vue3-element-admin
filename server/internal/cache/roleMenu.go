package cache

import (
	"context"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/encoding"
	"github.com/zhufuyi/sponge/pkg/utils"

	"admin/internal/model"
)

const (
	// cache prefix key, must end with a colon
	roleMenuCachePrefixKey = "roleMenu:"
	// RoleMenuExpireTime expire time
	RoleMenuExpireTime = 5 * time.Minute
)

var _ RoleMenuCache = (*roleMenuCache)(nil)

// RoleMenuCache cache interface
type RoleMenuCache interface {
	Set(ctx context.Context, id uint64, data *model.RoleMenu, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.RoleMenu, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.RoleMenu, error)
	MultiSet(ctx context.Context, data []*model.RoleMenu, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// roleMenuCache define a cache struct
type roleMenuCache struct {
	cache cache.Cache
}

// NewRoleMenuCache new a cache
func NewRoleMenuCache(cacheType *model.CacheType) RoleMenuCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.RoleMenu{}
		})
		return &roleMenuCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.RoleMenu{}
		})
		return &roleMenuCache{cache: c}
	}

	return nil // no cache
}

// GetRoleMenuCacheKey cache key
func (c *roleMenuCache) GetRoleMenuCacheKey(id uint64) string {
	return roleMenuCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *roleMenuCache) Set(ctx context.Context, id uint64, data *model.RoleMenu, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetRoleMenuCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *roleMenuCache) Get(ctx context.Context, id uint64) (*model.RoleMenu, error) {
	var data *model.RoleMenu
	cacheKey := c.GetRoleMenuCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *roleMenuCache) MultiSet(ctx context.Context, data []*model.RoleMenu, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetRoleMenuCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *roleMenuCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.RoleMenu, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetRoleMenuCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.RoleMenu)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.RoleMenu)
	for _, id := range ids {
		val, ok := itemMap[c.GetRoleMenuCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *roleMenuCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetRoleMenuCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *roleMenuCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetRoleMenuCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
