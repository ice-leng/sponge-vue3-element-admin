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
	platformCachePrefixKey = "platform:"
	// PlatformExpireTime expire time
	PlatformExpireTime = 5 * time.Minute
)

var _ PlatformCache = (*platformCache)(nil)

// PlatformCache cache interface
type PlatformCache interface {
	Set(ctx context.Context, id uint64, data *model.Platform, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Platform, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Platform, error)
	MultiSet(ctx context.Context, data []*model.Platform, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// platformCache define a cache struct
type platformCache struct {
	cache cache.Cache
}

// NewPlatformCache new a cache
func NewPlatformCache(cacheType *model.CacheType) PlatformCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Platform{}
		})
		return &platformCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Platform{}
		})
		return &platformCache{cache: c}
	}

	return nil // no cache
}

// GetPlatformCacheKey cache key
func (c *platformCache) GetPlatformCacheKey(id uint64) string {
	return platformCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *platformCache) Set(ctx context.Context, id uint64, data *model.Platform, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetPlatformCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *platformCache) Get(ctx context.Context, id uint64) (*model.Platform, error) {
	var data *model.Platform
	cacheKey := c.GetPlatformCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *platformCache) MultiSet(ctx context.Context, data []*model.Platform, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetPlatformCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *platformCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Platform, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetPlatformCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Platform)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Platform)
	for _, id := range ids {
		val, ok := itemMap[c.GetPlatformCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *platformCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetPlatformCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *platformCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetPlatformCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
