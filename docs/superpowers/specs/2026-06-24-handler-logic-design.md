# Handler Logic Generation Design

## Goal

Generate logic-layer code for every handler module that still keeps business work in `server/internal/handler`, using the existing Platform implementation as the reference style.

The target modules are:

- `config`
- `menu`
- `role`
- `roleMenu`
- `dashboard`
- `upload`

`platform` is already the reference implementation and should not be structurally rewritten.

## Reference Pattern

Use the Platform code as the canonical pattern:

- Handler style: `server/internal/handler/platform.go`
- Logic style: `server/internal/logic/platform.go`

Each converted handler should:

- hold a `logic.XLogic` field
- construct it through `logic.NewXLogic()`
- keep Gin request binding, path/current-user extraction, response writing, and request logging in the handler
- call logic methods for DAO access, data conversion, file work, and business operations
- use `handlerfunc.IsErrcode(err)` before falling back to a logged internal server error
- use validator messages for bind errors where the Platform handler already does so

Each generated logic file should:

- define an `XLogic` interface
- define a private `xLogic` struct
- expose `NewXLogic()` for production wiring
- expose `NewXLogicByDAO(...)` or equivalent dependency injection constructor for tests
- keep conversion helpers and business helpers in the logic package
- return typed response data, ids, totals, and errors without writing HTTP responses

## Scope

In scope:

- Move CRUD business logic from `config.go`, `menu.go`, `role.go`, and `roleMenu.go` handlers into matching logic files.
- Move menu tree conversion, role menu permission updates, role option building, config enum lookup, and object conversion helpers into logic.
- Move dashboard statistics and echarts response construction into `DashboardLogic`.
- Move upload file naming, directory creation, local file write, and URL construction into `UploadLogic`.
- Keep route paths, request structs, response payload shapes, and DAO behavior unchanged.
- Update handler tests only as needed to inject logic while preserving SQL mock coverage.

Out of scope:

- No generic CRUD abstraction.
- No route, API, request, response, model, DAO, cache, or database schema redesign.
- No unrelated cleanup or behavior changes.
- No dependency additions.

## Module Design

### Config

`ConfigLogic` should provide:

- `Create(ctx, *types.CreateConfigRequest) (uint64, error)`
- `DeleteByID(ctx, []uint64) error`
- `UpdateByID(ctx, *types.UpdateConfigByIDRequest) error`
- `GetByID(ctx, uint64) (*types.ConfigObjDetail, error)`
- `List(ctx, *types.ListConfigsRequest) ([]*types.ConfigObjDetail, int64, error)`
- `Dict(ctx) map[string]interface{}` or the exact type returned by `cache.EnumCache.GetAll`

The handler keeps comma-separated id parsing for delete requests, then passes the parsed ids to logic.

### Menu

`MenuLogic` should provide:

- `Create(ctx, *types.CreateMenuRequest) (uint64, error)`
- `DeleteByID(ctx, []uint64) error`
- `UpdateByID(ctx, *types.UpdateMenuByIDRequest) error`
- `GetByID(ctx, uint64) (*types.MenuObjDetail, error)`
- `List(ctx, *types.ListMenusRequest) ([]*types.MenuObjPage, error)`
- `Routes(ctx, types.LocalIntArray) ([]model.MenuItem, error)`
- `Options(ctx, *types.OptionMenusRequest) ([]types.Options, error)`

Menu list behavior must stay the same:

- force `request.Sort = "id"`
- default `ParentID` to `0`
- recursively load children
- set `RouteName` in the same way as current handler code

### Role

`RoleLogic` should provide:

- `Create(ctx, *types.CreateRoleRequest) (uint64, error)`
- `DeleteByID(ctx, []uint64) error`
- `UpdateByID(ctx, *types.UpdateRoleByIDRequest) error`
- `GetByID(ctx, uint64) (*types.RoleObjDetail, error)`
- `List(ctx, *types.ListRolesRequest) ([]*types.RoleObjDetail, int64, error)`
- `Options(ctx) ([]types.Options, error)`
- `MenuIds(ctx, uint64) ([]uint64, error)`
- `Menus(ctx, uint64, []uint64) error`

The handler keeps raw JSON parsing for the menu id array, then passes the role id and menu ids to logic.

### RoleMenu

`RoleMenuLogic` should provide:

- `Create(ctx, *types.CreateRoleMenuRequest) (uint64, error)`
- `DeleteByID(ctx, []uint64) error`
- `UpdateByID(ctx, *types.UpdateRoleMenuByIDRequest) error`
- `GetByID(ctx, uint64) (*types.RoleMenuObjDetail, error)`
- `List(ctx, *types.ListRoleMenusRequest) ([]*types.RoleMenuObjDetail, int64, error)`

### Dashboard

`DashboardLogic` should provide:

- `Statistics(ctx) []types.DashboardStatisticsItem`
- `Echarts(ctx, *types.DashboardEchartsRequest) types.Echarts`

The current static response values must remain unchanged.

### Upload

`UploadLogic` should provide:

- `Local(ctx, fileHeader *multipart.FileHeader) (*types.UploadItem, error)`

The logic method should keep current behavior:

- read form field `file` in the handler
- derive a new filename from original base name plus current timestamp
- write to `uploads/YYYY-MM-DD`
- return name, path, and config-based URL

## Error Handling

Handlers should follow Platform error flow:

1. Bind or path parsing errors return `ecode.InvalidParams`.
2. Logic errors that are errcode values are returned directly through `handlerfunc.IsErrcode`.
3. Other errors are logged in the handler and returned as internal server errors.

Logic should:

- return existing module-specific `ecode.ErrCreate*`, `ecode.ErrUpdateByID*`, `ecode.ErrGetByID*`, and `ecode.ErrList*` for conversion failures
- preserve current not-found response behavior for non-Platform modules by returning `ecode.NotFound.Err()` when DAO returns `database.ErrRecordNotFound`
- return DAO and filesystem errors upward unchanged unless an existing handler already maps them to an errcode

## Testing And Verification

After implementation:

- run focused handler tests for the affected modules
- run logic tests if new tests are added
- run `go test -count=1 -short ./internal/...` from `server`
- run `gofmt` on changed Go files

Expected observable result:

- existing HTTP routes keep the same request and response behavior
- handler files no longer contain DAO/model/copier business code for the converted modules
- new logic files compile and own the moved business behavior

