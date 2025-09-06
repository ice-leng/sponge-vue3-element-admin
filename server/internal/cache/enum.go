package cache

import (
	"admin/internal/config"
	"admin/internal/pkg/util"
	"admin/internal/types"
	"context"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"time"

	"github.com/go-dev-frame/sponge/pkg/cache"
	"github.com/go-dev-frame/sponge/pkg/encoding"
)

const (
	// cache prefix key, must end with a colon
	enumCachePrefixKey = "enum:"
	// all
	enumCacheKey = "all"
)

var _ EnumCache = (*enumCache)(nil)

// EnumCache cache interface
type EnumCache interface {
	Set(ctx context.Context, key string, data []*types.Options) error
	Get(ctx context.Context, key string) ([]*types.Options, error)
	MultiGet(ctx context.Context, keys []string) (map[string][]*types.Options, error)
	MultiSet(ctx context.Context, data map[string][]*types.Options) error

	GetAll(ctx context.Context) map[string][]*types.Options
	GetLabel(ctx context.Context, key string, value interface{}) string
}

// enumCache define a cache struct
type enumCache struct {
	cache cache.Cache
}

// NewEnumCache new a cache
func NewEnumCache() EnumCache {
	c := cache.NewMemoryCache("", encoding.JSONEncoding{}, func() interface{} {
		return &[]types.Options{}
	})

	cc := &enumCache{cache: c}
	_ = cc.load()
	return cc
}

func (c *enumCache) load() error {
	options := c.getOptions()
	err := c.setAll(options)
	if err != nil {
		return err
	}
	return c.MultiSet(context.Background(), options)
}

// GetEnumCacheKey cache key
func (c *enumCache) GetEnumCacheKey(key string) string {
	return enumCachePrefixKey + key
}

// Set write to cache
func (c *enumCache) Set(ctx context.Context, key string, data []*types.Options) error {
	if data == nil || key == "" {
		return nil
	}
	cacheKey := c.GetEnumCacheKey(key)
	err := c.cache.Set(ctx, cacheKey, data, 0*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// GetOptions cache key
func (c *enumCache) getOptions() map[string][]*types.Options {
	if config.Get().App.Env == "prod" {
		// 尝试从当前工作目录获取缓存路径
		wd, _ := os.Getwd()
		filePath := filepath.Join(wd, "enum.json")
		return util.EnumChangeDictByFile(filePath)
	}
	_, filename, _, _ := runtime.Caller(0)
	root := path.Dir(path.Dir(filename))
	enumDir := filepath.Join(root, "constant", "enum")
	return util.EnumChangeDict(enumDir)
}

func (c *enumCache) setAll(options map[string][]*types.Options) error {
	cacheKey := c.GetEnumCacheKey(enumCacheKey)
	err := c.cache.Set(context.Background(), cacheKey, &options, 0*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *enumCache) Get(ctx context.Context, key string) ([]*types.Options, error) {
	var data []*types.Options
	cacheKey := c.GetEnumCacheKey(key)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetAll cache value
func (c *enumCache) GetAll(ctx context.Context) map[string][]*types.Options {
	var data map[string][]*types.Options
	cacheKey := c.GetEnumCacheKey(enumCacheKey)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return make(map[string][]*types.Options)
	}
	return data
}

// MultiSet multiple set cache
func (c *enumCache) MultiSet(ctx context.Context, data map[string][]*types.Options) error {
	valMap := make(map[string]interface{})
	for k, v := range data {
		cacheKey := c.GetEnumCacheKey(k)
		valMap[cacheKey] = &v
	}

	err := c.cache.MultiSet(ctx, valMap, 0*time.Second)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *enumCache) MultiGet(ctx context.Context, ks []string) (map[string][]*types.Options, error) {
	var keys []string
	for _, v := range ks {
		cacheKey := c.GetEnumCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string][]*types.Options)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[string][]*types.Options)
	for _, id := range ks {
		val, ok := itemMap[c.GetEnumCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

func (c *enumCache) GetLabel(ctx context.Context, key string, value interface{}) string {
	options, err := c.Get(ctx, key)
	if err != nil {
		return ""
	}

	for _, v := range options {
		ov := reflect.ValueOf(v.Value).Convert(reflect.TypeOf(value)).Interface()
		if ov == value {
			return v.Label
		}
	}

	return ""
}
