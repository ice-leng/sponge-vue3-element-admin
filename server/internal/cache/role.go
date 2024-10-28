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
	roleCachePrefixKey = "role:"
	// RoleExpireTime expire time
	RoleExpireTime = 5 * time.Minute
)

var _ RoleCache = (*roleCache)(nil)

// RoleCache cache interface
type RoleCache interface {
	Set(ctx context.Context, id uint64, data *model.Role, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Role, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Role, error)
	MultiSet(ctx context.Context, data []*model.Role, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// roleCache define a cache struct
type roleCache struct {
	cache cache.Cache
}

// NewRoleCache new a cache
func NewRoleCache(cacheType *model.CacheType) RoleCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Role{}
		})
		return &roleCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Role{}
		})
		return &roleCache{cache: c}
	}

	return nil // no cache
}

// GetRoleCacheKey cache key
func (c *roleCache) GetRoleCacheKey(id uint64) string {
	return roleCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *roleCache) Set(ctx context.Context, id uint64, data *model.Role, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetRoleCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *roleCache) Get(ctx context.Context, id uint64) (*model.Role, error) {
	var data *model.Role
	cacheKey := c.GetRoleCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *roleCache) MultiSet(ctx context.Context, data []*model.Role, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetRoleCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *roleCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Role, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetRoleCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Role)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Role)
	for _, id := range ids {
		val, ok := itemMap[c.GetRoleCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *roleCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetRoleCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *roleCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetRoleCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
