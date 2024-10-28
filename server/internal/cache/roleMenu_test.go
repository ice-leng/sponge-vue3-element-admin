package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"admin/internal/model"
)

func newRoleMenuCache() *gotest.Cache {
	record1 := &model.RoleMenu{}
	record1.ID = 1
	record2 := &model.RoleMenu{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewRoleMenuCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_roleMenuCache_Set(t *testing.T) {
	c := newRoleMenuCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.RoleMenu)
	err := c.ICache.(RoleMenuCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(RoleMenuCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_roleMenuCache_Get(t *testing.T) {
	c := newRoleMenuCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.RoleMenu)
	err := c.ICache.(RoleMenuCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(RoleMenuCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(RoleMenuCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_roleMenuCache_MultiGet(t *testing.T) {
	c := newRoleMenuCache()
	defer c.Close()

	var testData []*model.RoleMenu
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.RoleMenu))
	}

	err := c.ICache.(RoleMenuCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(RoleMenuCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.RoleMenu))
	}
}

func Test_roleMenuCache_MultiSet(t *testing.T) {
	c := newRoleMenuCache()
	defer c.Close()

	var testData []*model.RoleMenu
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.RoleMenu))
	}

	err := c.ICache.(RoleMenuCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_roleMenuCache_Del(t *testing.T) {
	c := newRoleMenuCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.RoleMenu)
	err := c.ICache.(RoleMenuCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_roleMenuCache_SetCacheWithNotFound(t *testing.T) {
	c := newRoleMenuCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.RoleMenu)
	err := c.ICache.(RoleMenuCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewRoleMenuCache(t *testing.T) {
	c := NewRoleMenuCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewRoleMenuCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewRoleMenuCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
