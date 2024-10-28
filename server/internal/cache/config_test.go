package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"admin/internal/model"
)

func newConfigCache() *gotest.Cache {
	record1 := &model.Config{}
	record1.ID = 1
	record2 := &model.Config{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewConfigCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_configCache_Set(t *testing.T) {
	c := newConfigCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Config)
	err := c.ICache.(ConfigCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(ConfigCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_configCache_Get(t *testing.T) {
	c := newConfigCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Config)
	err := c.ICache.(ConfigCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(ConfigCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(ConfigCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_configCache_MultiGet(t *testing.T) {
	c := newConfigCache()
	defer c.Close()

	var testData []*model.Config
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Config))
	}

	err := c.ICache.(ConfigCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(ConfigCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Config))
	}
}

func Test_configCache_MultiSet(t *testing.T) {
	c := newConfigCache()
	defer c.Close()

	var testData []*model.Config
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Config))
	}

	err := c.ICache.(ConfigCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_configCache_Del(t *testing.T) {
	c := newConfigCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Config)
	err := c.ICache.(ConfigCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_configCache_SetCacheWithNotFound(t *testing.T) {
	c := newConfigCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Config)
	err := c.ICache.(ConfigCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewConfigCache(t *testing.T) {
	c := NewConfigCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewConfigCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewConfigCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
