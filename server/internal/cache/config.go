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
	configCachePrefixKey = "config:"
	// ConfigExpireTime expire time
	ConfigExpireTime = 5 * time.Minute
)

var _ ConfigCache = (*configCache)(nil)

// ConfigCache cache interface
type ConfigCache interface {
	Set(ctx context.Context, id uint64, data *model.Config, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Config, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Config, error)
	MultiSet(ctx context.Context, data []*model.Config, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error

	SetByKey(ctx context.Context, key string, data *model.Config, duration time.Duration) error
	GetByKey(ctx context.Context, key string) (*model.Config, error)
	DelByKey(ctx context.Context, key string) error
	SetCacheByKeyWithNotFound(ctx context.Context, key string) error
}

// configCache define a cache struct
type configCache struct {
	cache cache.Cache
}

// NewConfigCache new a cache
func NewConfigCache(cacheType *model.CacheType) ConfigCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Config{}
		})
		return &configCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Config{}
		})
		return &configCache{cache: c}
	}

	return nil // no cache
}

// GetConfigCacheKey cache key
func (c *configCache) GetConfigCacheKey(id uint64) string {
	return configCachePrefixKey + utils.Uint64ToStr(id)
}

func (c *configCache) GetConfigCacheKeyByKey(key string) string {
	return configCachePrefixKey + key
}

// Set write to cache
func (c *configCache) Set(ctx context.Context, id uint64, data *model.Config, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetConfigCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *configCache) Get(ctx context.Context, id uint64) (*model.Config, error) {
	var data *model.Config
	cacheKey := c.GetConfigCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *configCache) MultiSet(ctx context.Context, data []*model.Config, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetConfigCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *configCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Config, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetConfigCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Config)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Config)
	for _, id := range ids {
		val, ok := itemMap[c.GetConfigCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *configCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetConfigCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *configCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetConfigCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetByKey write to cache
func (c *configCache) SetByKey(ctx context.Context, key string, data *model.Config, duration time.Duration) error {
	if data == nil || key == "" {
		return nil
	}
	cacheKey := c.GetConfigCacheKeyByKey(key)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// GetByKey cache value
func (c *configCache) GetByKey(ctx context.Context, key string) (*model.Config, error) {
	var data *model.Config
	cacheKey := c.GetConfigCacheKeyByKey(key)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DelByKey delete cache
func (c *configCache) DelByKey(ctx context.Context, key string) error {
	cacheKey := c.GetConfigCacheKeyByKey(key)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

func (c *configCache) SetCacheByKeyWithNotFound(ctx context.Context, key string) error {
	cacheKey := c.GetConfigCacheKeyByKey(key)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
