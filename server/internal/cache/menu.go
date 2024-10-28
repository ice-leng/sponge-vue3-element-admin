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
	menuCachePrefixKey = "menu:"
	// MenuExpireTime expire time
	MenuExpireTime = 5 * time.Minute
)

var _ MenuCache = (*menuCache)(nil)

// MenuCache cache interface
type MenuCache interface {
	Set(ctx context.Context, id uint64, data *model.Menu, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Menu, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Menu, error)
	MultiSet(ctx context.Context, data []*model.Menu, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// menuCache define a cache struct
type menuCache struct {
	cache cache.Cache
}

// NewMenuCache new a cache
func NewMenuCache(cacheType *model.CacheType) MenuCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Menu{}
		})
		return &menuCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Menu{}
		})
		return &menuCache{cache: c}
	}

	return nil // no cache
}

// GetMenuCacheKey cache key
func (c *menuCache) GetMenuCacheKey(id uint64) string {
	return menuCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *menuCache) Set(ctx context.Context, id uint64, data *model.Menu, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetMenuCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *menuCache) Get(ctx context.Context, id uint64) (*model.Menu, error) {
	var data *model.Menu
	cacheKey := c.GetMenuCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *menuCache) MultiSet(ctx context.Context, data []*model.Menu, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetMenuCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *menuCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Menu, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetMenuCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Menu)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Menu)
	for _, id := range ids {
		val, ok := itemMap[c.GetMenuCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *menuCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetMenuCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *menuCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetMenuCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
