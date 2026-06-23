# Handler Logic Generation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Move non-Platform handler business logic into `server/internal/logic` using the existing Platform handler and logic code as the reference implementation.

**Architecture:** Each affected handler will own a `logic.XLogic` dependency and keep only HTTP binding, path/current-user extraction, response writing, and request logging. Each new logic file will own DAO/cache dependencies, model/type conversion, static data construction, file work, and business operations through `NewXLogic()` plus injection constructors for tests.

**Tech Stack:** Go 1.24.2, Gin, Sponge DAO/cache/response/errcode packages, GORM, existing `go test` test suite.

---

## File Structure

- Create `server/internal/logic/config.go`: `ConfigLogic`, DAO/cache enum wiring, config CRUD, dict lookup, config conversion helpers.
- Create `server/internal/logic/menu.go`: `MenuLogic`, menu CRUD, recursive menu page conversion, routes/options forwarding.
- Create `server/internal/logic/role.go`: `RoleLogic`, role CRUD, role options, menu id lookup, role-menu update.
- Create `server/internal/logic/roleMenu.go`: `RoleMenuLogic`, role-menu CRUD and conversion helpers.
- Create `server/internal/logic/dashboard.go`: `DashboardLogic`, existing static statistics and echarts data.
- Create `server/internal/logic/upload.go`: `UploadLogic`, local upload filename/path/write/url logic.
- Modify `server/internal/handler/config.go`: replace DAO/cache fields and business code with Platform-style logic calls.
- Modify `server/internal/handler/menu.go`: replace DAO/model/copier business code with Platform-style logic calls.
- Modify `server/internal/handler/role.go`: replace DAO/model/copier/role-menu business code with Platform-style logic calls.
- Modify `server/internal/handler/roleMenu.go`: replace DAO/model/copier business code with Platform-style logic calls.
- Modify `server/internal/handler/dashboard.go`: call `logic.DashboardLogic`.
- Modify `server/internal/handler/upload.go`: call `logic.UploadLogic`.
- Modify handler tests as needed: use `logic.NewXLogicByDAO(...)` when constructing test handlers so SQL mocks still cover DAO calls.

## Task 1: Add Config Logic And Convert Config Handler

**Files:**
- Create: `server/internal/logic/config.go`
- Modify: `server/internal/handler/config.go`
- Modify: `server/internal/handler/config_test.go`

- [ ] **Step 1: Write the failing compile target**

Run:

```bash
cd server
go test -count=1 -short ./internal/handler -run 'Test_configHandler_(Create|DeleteByID|UpdateByID|GetByID|List)|TestNewConfigHandler'
```

Expected before implementation: current tests pass or compile against old handler. After changing handler to `logic.ConfigLogic` without `logic/config.go`, compilation fails with undefined `logic.ConfigLogic` or missing constructor. This confirms the test target covers the integration point.

- [ ] **Step 2: Create `server/internal/logic/config.go`**

Implement:

```go
package logic

import (
	"admin/internal/cache"
	"admin/internal/dao"
	"admin/internal/database"
	"admin/internal/ecode"
	"admin/internal/model"
	"admin/internal/types"
	"context"
	"errors"

	"github.com/go-dev-frame/sponge/pkg/copier"
)

type ConfigLogic interface {
	Create(ctx context.Context, request *types.CreateConfigRequest) (uint64, error)
	DeleteByID(ctx context.Context, ids []uint64) error
	UpdateByID(ctx context.Context, request *types.UpdateConfigByIDRequest) error
	GetByID(ctx context.Context, id uint64) (*types.ConfigObjDetail, error)
	List(ctx context.Context, request *types.ListConfigsRequest) ([]*types.ConfigObjDetail, int64, error)
	Dict(ctx context.Context) map[string]interface{}
}
```

Add `configLogic` with `iDao dao.ConfigDao` and `cEnum cache.EnumCache`, plus `NewConfigLogic()` and `NewConfigLogicByDAO(iDao dao.ConfigDao, cEnum cache.EnumCache) ConfigLogic`.

Move current handler bodies for create/update/get/list/dict into logic methods. Use `ecode.NotFound.Err()` when `GetByID` receives `database.ErrRecordNotFound`. Move `convertConfig` and `convertConfigs` into this file.

- [ ] **Step 3: Convert `server/internal/handler/config.go`**

Change the handler struct to:

```go
type configHandler struct {
	logic logic.ConfigLogic
}
```

Change `NewConfigHandler` to call `logic.NewConfigLogic()`. For each method, follow the Platform handler structure: bind request, build `ctx := middleware.WrapCtx(c)`, call logic, return errcode errors through `handlerfunc.IsErrcode(err)`, otherwise log and return internal server error.

Keep comma-separated delete id parsing in handler and pass `[]uint64` into `h.logic.DeleteByID(ctx, ids)`.

- [ ] **Step 4: Update `server/internal/handler/config_test.go` injection**

Construct the test handler as:

```go
h.IHandler = &configHandler{
	logic: logic.NewConfigLogicByDAO(d.IDao.(dao.ConfigDao), cache.NewEnumCache()),
}
```

Add the `admin/internal/logic` import and remove imports no longer used by the test.

- [ ] **Step 5: Run focused config tests**

Run:

```bash
cd server
go test -count=1 -short ./internal/handler -run 'Test_configHandler_(Create|DeleteByID|UpdateByID|GetByID|List)|TestNewConfigHandler'
```

Expected: PASS.

## Task 2: Add Menu Logic And Convert Menu Handler

**Files:**
- Create: `server/internal/logic/menu.go`
- Modify: `server/internal/handler/menu.go`
- Modify: `server/internal/handler/menu_test.go`

- [ ] **Step 1: Create `server/internal/logic/menu.go`**

Define:

```go
type MenuLogic interface {
	Create(ctx context.Context, request *types.CreateMenuRequest) (uint64, error)
	DeleteByID(ctx context.Context, ids []uint64) error
	UpdateByID(ctx context.Context, request *types.UpdateMenuByIDRequest) error
	GetByID(ctx context.Context, id uint64) (*types.MenuObjDetail, error)
	List(ctx context.Context, request *types.ListMenusRequest) ([]*types.MenuObjPage, error)
	Routes(ctx context.Context, roleIds types.LocalIntArray) ([]model.MenuItem, error)
	Options(ctx context.Context, request *types.OptionMenusRequest) ([]types.Options, error)
}
```

Add `menuLogic`, `NewMenuLogic()`, and `NewMenuLogicByDAO(iDao dao.MenuDao) MenuLogic`.

Move current create/update/get/list/routes/options code into logic. Move recursive `convertMenu` and `convertMenus` into logic and replace `*gin.Context` use with `context.Context`. Preserve `request.Sort = "id"`, default `ParentID = 0`, recursive child loading, and `RouteName` assignments.

- [ ] **Step 2: Convert `server/internal/handler/menu.go`**

Use:

```go
type menuHandler struct {
	logic logic.MenuLogic
}
```

Use `logic.NewMenuLogic()` in the constructor. Keep bind/path/delete id parsing in handler. Forward route role ids from Gin context:

```go
roleIds, _ := c.Get("roleId")
result, err := h.logic.Routes(ctx, roleIds.(types.LocalIntArray))
```

Handle errors using the Platform pattern.

- [ ] **Step 3: Update `server/internal/handler/menu_test.go` injection**

Construct:

```go
h.IHandler = &menuHandler{
	logic: logic.NewMenuLogicByDAO(d.IDao.(dao.MenuDao)),
}
```

Add the `admin/internal/logic` import and remove imports no longer used by the test.

- [ ] **Step 4: Run focused menu tests**

Run:

```bash
cd server
go test -count=1 -short ./internal/handler -run 'Test_menuHandler_(Create|DeleteByID|UpdateByID|GetByID|List)|TestNewMenuHandler'
```

Expected: PASS.

## Task 3: Add Role And RoleMenu Logic, Convert Their Handlers

**Files:**
- Create: `server/internal/logic/role.go`
- Create: `server/internal/logic/roleMenu.go`
- Modify: `server/internal/handler/role.go`
- Modify: `server/internal/handler/roleMenu.go`
- Modify: `server/internal/handler/role_test.go`
- Modify: `server/internal/handler/roleMenu_test.go`

- [ ] **Step 1: Create `server/internal/logic/roleMenu.go`**

Define `RoleMenuLogic` with CRUD/List methods matching the design spec. Add `roleMenuLogic`, `NewRoleMenuLogic()`, and `NewRoleMenuLogicByDAO(iDao dao.RoleMenuDao) RoleMenuLogic`. Move `convertRoleMenu` and `convertRoleMenus` into this file.

- [ ] **Step 2: Create `server/internal/logic/role.go`**

Define:

```go
type RoleLogic interface {
	Create(ctx context.Context, request *types.CreateRoleRequest) (uint64, error)
	DeleteByID(ctx context.Context, ids []uint64) error
	UpdateByID(ctx context.Context, request *types.UpdateRoleByIDRequest) error
	GetByID(ctx context.Context, id uint64) (*types.RoleObjDetail, error)
	List(ctx context.Context, request *types.ListRolesRequest) ([]*types.RoleObjDetail, int64, error)
	Options(ctx context.Context) ([]types.Options, error)
	MenuIds(ctx context.Context, id uint64) ([]uint64, error)
	Menus(ctx context.Context, id uint64, menuIds []uint64) error
}
```

Add `roleLogic`, `NewRoleLogic()`, and `NewRoleLogicByDAO(iDao dao.RoleDao, iRoleMenuDao dao.RoleMenuDao) RoleLogic`. Move role conversion helpers and role-menu operations into logic. Preserve current behavior that `Options` ignores DAO errors, unless an error is returned by the existing code path after refactor would be newly observable; prefer exact current response behavior.

- [ ] **Step 3: Convert `role.go` and `roleMenu.go` handlers**

Change handler structs to hold `logic.RoleLogic` and `logic.RoleMenuLogic`. Constructors call `logic.NewRoleLogic()` and `logic.NewRoleMenuLogic()`. Handlers keep request binding, raw JSON parsing for role menu ids, id parsing, responses, and logging. Business work moves to logic calls.

- [ ] **Step 4: Update role tests**

For `role_test.go`, construct:

```go
h.IHandler = &roleHandler{
	logic: logic.NewRoleLogicByDAO(d.IDao.(dao.RoleDao), dao.NewRoleMenuDao(d.DB, cache.NewRoleMenuCache(&database.CacheType{CType: "redis", Rdb: c.RedisClient}))),
}
```

For `roleMenu_test.go`, construct:

```go
h.IHandler = &roleMenuHandler{
	logic: logic.NewRoleMenuLogicByDAO(d.IDao.(dao.RoleMenuDao)),
}
```

Add imports and remove unused imports after compiling.

- [ ] **Step 5: Run focused role tests**

Run:

```bash
cd server
go test -count=1 -short ./internal/handler -run 'Test_roleHandler_(Create|DeleteByID|UpdateByID|GetByID|List)|TestNewRoleHandler|Test_roleMenuHandler_(Create|DeleteByID|UpdateByID|GetByID|List)|TestNewRoleMenuHandler'
```

Expected: PASS.

## Task 4: Add Dashboard And Upload Logic, Convert Handlers

**Files:**
- Create: `server/internal/logic/dashboard.go`
- Create: `server/internal/logic/upload.go`
- Modify: `server/internal/handler/dashboard.go`
- Modify: `server/internal/handler/upload.go`

- [ ] **Step 1: Create `server/internal/logic/dashboard.go`**

Define `DashboardLogic` with `Statistics(ctx context.Context) []types.DashboardStatisticsItem` and `Echarts(ctx context.Context, request *types.DashboardEchartsRequest) types.Echarts`. Move the exact static data from `dashboardHandler` into these methods.

- [ ] **Step 2: Convert dashboard handler**

Give `dashboardHandler` a `logic logic.DashboardLogic` field, construct it through `logic.NewDashboardLogic()`, and call logic methods after binding query params. Keep invalid query handling as currently implemented.

- [ ] **Step 3: Create `server/internal/logic/upload.go`**

Define `UploadLogic` with:

```go
Local(ctx context.Context, fileHeader *multipart.FileHeader) (*types.UploadItem, error)
```

Move existing upload behavior into logic: `fileHeader.Open()`, new filename with `gocrypto.Md5`, `uploads/YYYY-MM-DD`, `gofile.CreateDir`, `os.Create`, `io.Copy`, and config-based URL. Return `errcode.NewError(10001, err.Error())` for existing filesystem/open/copy failures to preserve current API errors.

- [ ] **Step 4: Convert upload handler**

Keep `c.Request.FormFile("file")` in the handler to preserve the request boundary. Pass the returned file header to `h.logic.Local(ctx, file)`. Use Platform-style error handling.

- [ ] **Step 5: Run package compile tests**

Run:

```bash
cd server
go test -count=1 -short ./internal/handler ./internal/logic
```

Expected: PASS.

## Task 5: Format, Search, And Full Verification

**Files:**
- All changed Go files.

- [ ] **Step 1: Run gofmt**

Run:

```bash
cd server
gofmt -w internal/logic/config.go internal/logic/menu.go internal/logic/role.go internal/logic/roleMenu.go internal/logic/dashboard.go internal/logic/upload.go internal/handler/config.go internal/handler/menu.go internal/handler/role.go internal/handler/roleMenu.go internal/handler/dashboard.go internal/handler/upload.go internal/handler/config_test.go internal/handler/menu_test.go internal/handler/role_test.go internal/handler/roleMenu_test.go
```

Expected: no output.

- [ ] **Step 2: Search for stale handler business imports**

Run:

```bash
cd server
rg -n '"admin/internal/(cache|dao|database|model)"|github.com/jinzhu/copier|pkg/copier' internal/handler/config.go internal/handler/menu.go internal/handler/role.go internal/handler/roleMenu.go internal/handler/dashboard.go internal/handler/upload.go
```

Expected: no matches except imports still required by handler-only concerns. Investigate and remove stale imports when matches are unused or indicate business code left in handler.

- [ ] **Step 3: Run affected handler tests**

Run:

```bash
cd server
go test -count=1 -short ./internal/handler
```

Expected: PASS.

- [ ] **Step 4: Run broader server tests**

Run:

```bash
cd server
go test -count=1 -short ./internal/...
```

Expected: PASS, or a clearly documented pre-existing/environment failure unrelated to the refactor.

- [ ] **Step 5: Review diff**

Run:

```bash
git diff --stat
git diff -- server/internal/logic server/internal/handler
```

Expected: changes are limited to new logic files, converted handlers, and necessary test injection updates. No generated noise, credentials, dependency changes, or unrelated formatting churn.

