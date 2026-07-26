package handler

import (
	"admin/internal/cache"
	"admin/internal/dao"
	"admin/internal/database"
	"admin/internal/logic"
	"admin/internal/model"
	"admin/internal/types"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gocrypto"
	"github.com/go-dev-frame/sponge/pkg/gotest"
	"github.com/go-dev-frame/sponge/pkg/httpcli"
	"github.com/go-dev-frame/sponge/pkg/utils"
	"github.com/jinzhu/copier"
)

// withAuth wraps a handler to inject auth context (simulates JWT middleware)
func withAuth(uid uint64, fn gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("id", uid)
		fn(c)
	}
}

func newPlatformHandler() *gotest.Handler {
	testData := &model.Platform{}
	testData.ID = 1
	// you can set the other fields of testData here, such as:
	//testData.CreatedAt = time.Now()
	//testData.UpdatedAt = testData.CreatedAt

	// init mock cache for platform
	c := gotest.NewCache(map[string]interface{}{utils.Uint64ToStr(testData.ID): testData})
	c.ICache = cache.NewPlatformCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})

	// init mock dao
	d := gotest.NewDao(c, testData)
	d.IDao = dao.NewPlatformDao(d.DB, c.ICache.(cache.PlatformCache))

	// init mock role dao (separate cache)
	cRole := gotest.NewCache(map[string]interface{}{})
	cRole.ICache = cache.NewRoleCache(&database.CacheType{
		CType: "redis",
		Rdb:   cRole.RedisClient,
	})
	dRole := gotest.NewDao(cRole, &model.Role{})
	dRole.IDao = dao.NewRoleDao(dRole.DB, cRole.ICache.(cache.RoleCache))

	// init mock config dao (separate cache)
	cConfig := gotest.NewCache(map[string]interface{}{})
	cConfig.ICache = cache.NewConfigCache(&database.CacheType{
		CType: "redis",
		Rdb:   cConfig.RedisClient,
	})
	dConfig := gotest.NewDao(cConfig, &model.Config{})
	dConfig.IDao = dao.NewConfigDao(dConfig.DB, cConfig.ICache.(cache.ConfigCache))

	// init mock handler
	h := gotest.NewHandler(d, testData)
	h.IHandler = &platformHandler{
		logic: logic.NewPlatformLogicByDAO(
			d.IDao.(dao.PlatformDao),
			dRole.IDao.(dao.RoleDao),
			logic.NewConfigLogicByDAO(dConfig.IDao.(dao.ConfigDao), nil),
		),
	}
	iHandler := h.IHandler.(PlatformHandler)

	testFns := []gotest.RouterInfo{
		{
			FuncName:    "Create",
			Method:      http.MethodPost,
			Path:        "/platform",
			HandlerFunc: iHandler.Create,
		},
		{
			FuncName:    "DeleteByID",
			Method:      http.MethodDelete,
			Path:        "/platform/:id",
			HandlerFunc: iHandler.DeleteByID,
		},
		{
			FuncName:    "UpdateByID",
			Method:      http.MethodPut,
			Path:        "/platform/:id",
			HandlerFunc: iHandler.UpdateByID,
		},
		{
			FuncName:    "GetByID",
			Method:      http.MethodGet,
			Path:        "/platform/:id",
			HandlerFunc: iHandler.GetByID,
		},
		{
			FuncName:    "List",
			Method:      http.MethodGet,
			Path:        "/platform",
			HandlerFunc: iHandler.List,
		},
		{
			FuncName:    "Me",
			Method:      http.MethodGet,
			Path:        "/platform/me",
			HandlerFunc: withAuth(testData.ID, iHandler.Me),
		},
		{
			FuncName:    "GetProfile",
			Method:      http.MethodGet,
			Path:        "/platform/profile",
			HandlerFunc: withAuth(testData.ID, iHandler.GetProfile),
		},
		{
			FuncName:    "UpdateProfile",
			Method:      http.MethodPut,
			Path:        "/platform/profile",
			HandlerFunc: withAuth(testData.ID, iHandler.UpdateProfile),
		},
		{
			FuncName:    "ChangePassword",
			Method:      http.MethodPut,
			Path:        "/platform/password",
			HandlerFunc: withAuth(testData.ID, iHandler.ChangePassword),
		},
		{
			FuncName:    "ResetPassword",
			Method:      http.MethodPut,
			Path:        "/platform/password/reset",
			HandlerFunc: withAuth(testData.ID, iHandler.ResetPassword),
		},
	}

	h.GoRunHTTPServer(testFns)

	time.Sleep(time.Millisecond * 200)
	return h
}

func Test_platformHandler_Create(t *testing.T) {
	h := newPlatformHandler()
	defer h.Close()
	testData := &types.CreatePlatformRequest{}
	_ = copier.Copy(testData, h.TestData.(*model.Platform))

	h.MockDao.SQLMock.ExpectBegin()
	args := h.MockDao.GetAnyArgs(h.TestData)
	h.MockDao.SQLMock.ExpectExec("INSERT INTO .*").
		WithArgs(args[:len(args)-1]...). // adjusted for the amount of test data
		WillReturnResult(sqlmock.NewResult(1, 1))
	h.MockDao.SQLMock.ExpectCommit()

	result := &httpcli.StdResult{}
	err := httpcli.Post(result, h.GetRequestURL("Create"), testData)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", result)
}

func Test_platformHandler_DeleteByID(t *testing.T) {
	h := newPlatformHandler()
	defer h.Close()

	h.MockDao.SQLMock.ExpectBegin()
	h.MockDao.SQLMock.ExpectExec("UPDATE .*deleted_at").
		WithArgs(sqlmock.AnyArg(), h.TestData.(*model.Platform).ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	h.MockDao.SQLMock.ExpectCommit()

	result := &httpcli.StdResult{}
	err := httpcli.Delete(result, h.GetRequestURL("DeleteByID", h.TestData.(*model.Platform).ID))
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}
}

func Test_platformHandler_UpdateByID(t *testing.T) {
	h := newPlatformHandler()
	defer h.Close()
	testData := &types.UpdatePlatformByIDRequest{}
	_ = copier.Copy(testData, h.TestData.(*model.Platform))

	h.MockDao.SQLMock.ExpectBegin()
	h.MockDao.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(sqlmock.AnyArg(), h.TestData.(*model.Platform).ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	h.MockDao.SQLMock.ExpectCommit()

	result := &httpcli.StdResult{}
	err := httpcli.Put(result, h.GetRequestURL("UpdateByID", h.TestData.(*model.Platform).ID), testData)
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}
}

func Test_platformHandler_GetByID(t *testing.T) {
	h := newPlatformHandler()
	defer h.Close()

	h.MockDao.SQLMock.ExpectQuery("SELECT .*").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(h.TestData.(*model.Platform).ID))

	result := &httpcli.StdResult{}
	err := httpcli.Get(result, h.GetRequestURL("GetByID", h.TestData.(*model.Platform).ID))
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}
}

func Test_platformHandler_List(t *testing.T) {
	h := newPlatformHandler()
	defer h.Close()

	h.MockDao.SQLMock.ExpectQuery("SELECT .*").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(h.TestData.(*model.Platform).ID))

	result := &httpcli.StdResult{}
	params := httpcli.KV{"page": 1, "pageSize": 10, "sort": "ignore count"}
	err := httpcli.Get(result, h.GetRequestURL("List"), httpcli.WithParams(params))
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}
}

func Test_platformHandler_Me(t *testing.T) {
	h := newPlatformHandler()
	defer h.Close()
	testData := h.TestData.(*model.Platform)

	// column names and corresponding data
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(testData.ID)

	h.MockDao.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID, 1).
		WillReturnRows(rows)

	result := &httpcli.StdResult{}
	err := httpcli.Get(result, h.GetRequestURL("Me"))
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}
}

func Test_platformHandler_GetProfile(t *testing.T) {
	h := newPlatformHandler()
	defer h.Close()
	testData := h.TestData.(*model.Platform)

	// column names and corresponding data
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(testData.ID)

	h.MockDao.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID, 1).
		WillReturnRows(rows)

	result := &httpcli.StdResult{}
	err := httpcli.Get(result, h.GetRequestURL("GetProfile"))
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}
}

func Test_platformHandler_UpdateProfile(t *testing.T) {
	h := newPlatformHandler()
	defer h.Close()
	testData := &types.UpdatePlatformByIDRequest{}
	_ = copier.Copy(testData, h.TestData.(*model.Platform))

	h.MockDao.SQLMock.ExpectBegin()
	h.MockDao.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(sqlmock.AnyArg(), h.TestData.(*model.Platform).ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	h.MockDao.SQLMock.ExpectCommit()

	result := &httpcli.StdResult{}
	err := httpcli.Put(result, h.GetRequestURL("UpdateProfile"), testData)
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}
}

func Test_platformHandler_ChangePassword(t *testing.T) {
	h := newPlatformHandler()
	defer h.Close()
	testData := &types.ChangePasswordRequest{
		OldPassword:     "old123456",
		NewPassword:     "new123456",
		ConfirmPassword: "new123456",
	}
	testData.ID = h.TestData.(*model.Platform).ID

	// ChangePassword calls GetByID first - must include password hash for verification
	oldPasswordHash, _ := gocrypto.HashAndSaltPassword("old123456")
	h.MockDao.SQLMock.ExpectQuery("SELECT .*").
		WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).AddRow(h.TestData.(*model.Platform).ID, oldPasswordHash))
	h.MockDao.SQLMock.ExpectBegin()
	h.MockDao.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), h.TestData.(*model.Platform).ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	h.MockDao.SQLMock.ExpectCommit()

	result := &httpcli.StdResult{}
	err := httpcli.Put(result, h.GetRequestURL("ChangePassword"), testData)
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}
}

func Test_platformHandler_ResetPassword(t *testing.T) {
	h := newPlatformHandler()
	defer h.Close()
	testData := &types.ResetPasswordRequest{}
	_ = copier.Copy(testData, h.TestData.(*model.Platform))

	// ResetPassword calls GetByID first
	h.MockDao.SQLMock.ExpectQuery("SELECT .*").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(h.TestData.(*model.Platform).ID))
	h.MockDao.SQLMock.ExpectBegin()
	h.MockDao.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), h.TestData.(*model.Platform).ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	h.MockDao.SQLMock.ExpectCommit()

	result := &httpcli.StdResult{}
	err := httpcli.Put(result, h.GetRequestURL("ResetPassword"), testData)
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}
}

func TestNewPlatformHandler(t *testing.T) {
	defer func() {
		recover()
	}()
	_ = NewPlatformHandler()
}
