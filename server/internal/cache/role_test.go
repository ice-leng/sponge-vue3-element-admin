package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"admin/internal/model"
)

func newRoleCache() *gotest.Cache {
	record1 := &model.Role{}
	record1.ID = 1
	record2 := &model.Role{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewRoleCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_roleCache_Set(t *testing.T) {
	c := newRoleCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Role)
	err := c.ICache.(RoleCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(RoleCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_roleCache_Get(t *testing.T) {
	c := newRoleCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Role)
	err := c.ICache.(RoleCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(RoleCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(RoleCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_roleCache_MultiGet(t *testing.T) {
	c := newRoleCache()
	defer c.Close()

	var testData []*model.Role
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Role))
	}

	err := c.ICache.(RoleCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(RoleCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Role))
	}
}

func Test_roleCache_MultiSet(t *testing.T) {
	c := newRoleCache()
	defer c.Close()

	var testData []*model.Role
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Role))
	}

	err := c.ICache.(RoleCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_roleCache_Del(t *testing.T) {
	c := newRoleCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Role)
	err := c.ICache.(RoleCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_roleCache_SetCacheWithNotFound(t *testing.T) {
	c := newRoleCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Role)
	err := c.ICache.(RoleCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewRoleCache(t *testing.T) {
	c := NewRoleCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewRoleCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewRoleCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
