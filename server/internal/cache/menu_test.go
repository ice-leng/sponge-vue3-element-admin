package cache

import (
	"admin/internal/database"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-dev-frame/sponge/pkg/gotest"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"admin/internal/model"
)

func newMenuCache() *gotest.Cache {
	record1 := &model.Menu{}
	record1.ID = 1
	record2 := &model.Menu{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewMenuCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_menuCache_Set(t *testing.T) {
	c := newMenuCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Menu)
	err := c.ICache.(MenuCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(MenuCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_menuCache_Get(t *testing.T) {
	c := newMenuCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Menu)
	err := c.ICache.(MenuCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(MenuCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(MenuCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_menuCache_MultiGet(t *testing.T) {
	c := newMenuCache()
	defer c.Close()

	var testData []*model.Menu
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Menu))
	}

	err := c.ICache.(MenuCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(MenuCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Menu))
	}
}

func Test_menuCache_MultiSet(t *testing.T) {
	c := newMenuCache()
	defer c.Close()

	var testData []*model.Menu
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Menu))
	}

	err := c.ICache.(MenuCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_menuCache_Del(t *testing.T) {
	c := newMenuCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Menu)
	err := c.ICache.(MenuCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_menuCache_SetCacheWithNotFound(t *testing.T) {
	c := newMenuCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Menu)
	err := c.ICache.(MenuCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(MenuCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewMenuCache(t *testing.T) {
	c := NewMenuCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewMenuCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewMenuCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
