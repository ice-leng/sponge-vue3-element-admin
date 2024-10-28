package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"admin/internal/model"
)

func newPlatformCache() *gotest.Cache {
	record1 := &model.Platform{}
	record1.ID = 1
	record2 := &model.Platform{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewPlatformCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_platformCache_Set(t *testing.T) {
	c := newPlatformCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Platform)
	err := c.ICache.(PlatformCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(PlatformCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_platformCache_Get(t *testing.T) {
	c := newPlatformCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Platform)
	err := c.ICache.(PlatformCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(PlatformCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(PlatformCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_platformCache_MultiGet(t *testing.T) {
	c := newPlatformCache()
	defer c.Close()

	var testData []*model.Platform
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Platform))
	}

	err := c.ICache.(PlatformCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(PlatformCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Platform))
	}
}

func Test_platformCache_MultiSet(t *testing.T) {
	c := newPlatformCache()
	defer c.Close()

	var testData []*model.Platform
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Platform))
	}

	err := c.ICache.(PlatformCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_platformCache_Del(t *testing.T) {
	c := newPlatformCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Platform)
	err := c.ICache.(PlatformCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_platformCache_SetCacheWithNotFound(t *testing.T) {
	c := newPlatformCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Platform)
	err := c.ICache.(PlatformCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewPlatformCache(t *testing.T) {
	c := NewPlatformCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewPlatformCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewPlatformCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
