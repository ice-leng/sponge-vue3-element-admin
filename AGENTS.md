# AGENTS.md - sponge-vue3-element-admin

> **面向 Agent 的项目规则与执行指南**  
> 面向人类的说明请看 `README.md`；本文件仅供 Agent 执行参考。

---

## 1. 规则分级

| 等级 | 含义 | 违规后果 |
|------|------|----------|
| **必选** | 必须满足，未满足视为任务未完成 | 直接停止，反馈失败点 |
| **建议** | 无冲突时默认遵循；冲突需说明原因 | 记录理由，不阻断 |

**规则治理**：
- **必选**：每新增 1 条，必须删除/合并 ≥1 条旧规则，保持总量不膨胀
- **建议**：每两周审计一次，清理重复/不可验证/低频条目

---

## 2. 项目结构（单体双端）

```
sponge-vue3-element-admin/
├── server/                 # Go 后端（Sponge + Gin + GORM）
│   ├── cmd/admin/          # 服务入口与初始化
│   ├── internal/
│   │   ├── handler/        # HTTP 绑定/解析 → 调用 logic → 统一响应
│   │   ├── logic/          # 业务编排、错误转换、唯一索引预查重
│   │   ├── dao/            # 数据访问细节、唯一键冲突转业务错误
│   │   ├── model/          # GORM 模型（dao 专用，禁止跨层直连）
│   │   ├── types/          # DTO/请求/响应定义
│   │   ├── routers/        # 路由注册
│   │   ├── server/         # HTTP/gRPC 启动逻辑
│   │   ├── middleware/     # 统一中间件
│   │   ├── cache/          # 缓存封装
│   │   ├── config/         # 配置加载
│   │   ├── constant/       # 常量/枚举
│   │   ├── database/       # DB 初始化
│   │   └── ecode/          # 业务错误码
│   ├── configs/            # 配置模板（admin.yml）
│   ├── docs/               # Swagger/OpenAPI 产物
│   ├── deployments/        # 部署脚本
│   ├── scripts/            # 构建脚本
│   ├── pkg/                # 公共库
│   ├── go.mod / go.sum
│   └── Makefile
│
└── web/                    # Vue3 前端
    ├── src/
    │   ├── api/            # API 请求封装
    │   ├── components/     # 通用组件
    │   ├── views/          # 页面视图
    │   ├── router/         # 路由配置
    │   ├── store/          # Pinia 状态
    │   ├── types/          # TS 类型定义
    │   ├── utils/          # 工具函数
    │   ├── hooks/          # 组合式函数
    │   ├── layout/         # 布局组件
    │   ├── styles/         # 全局样式
    │   ├── plugins/        # 插件注册
    │   ├── directive/      # 自定义指令
    │   ├── enums/          # 枚举常量
    │   └── lang/           # 国际化
    ├── package.json
    ├── vite.config.ts
    ├── tsconfig.json
    └── .eslintrc / .prettierrc / stylelint.config
```

---

## 3. 常用命令

### 后端
| 命令 | 说明 |
|------|------|
| `make run` | 构建并运行（可选 `Config=configs/dev.yml`） |
| `make docs` | 生成 Swagger 文档 |
| `make test` | `go test -short ./...` |
| `make cover` | 覆盖率报告 |
| `make ci-lint` | `gofmt` + `golangci-lint` |
| `make build` | 构建 Linux amd64 二进制（`cmd/admin`） |

### 前端
| 命令 | 说明 |
|------|------|
| `pnpm dev` | 开发服务器（Vite） |
| `pnpm build` | 生产构建（`vue-tsc --noEmit && vite build`） |
| `pnpm preview` | 预览构建产物 |
| `pnpm type-check` | `vue-tsc --noEmit` |
| `pnpm lint:eslint` | ESLint + 自动修复 |
| `pnpm lint:prettier` | Prettier 格式化 |
| `pnpm lint:stylelint` | Stylelint 修复 |

---

## 4. 后端开发规则（Go + Sponge + Gin）

### 4.1 必选
- **仅修改任务直接相关文件**；发现无关改动/脏变更先暂停并确认。
- **严格分层**：`handler → logic → dao`，**禁止** `handler` 直连 `dao`/`model`，**禁止** `logic` 直连 `model`。
- **handler 职责**：绑定/校验参数 → `middleware.WrapCtx(c)` 取 `context.Context` → 调用 `logic` → `response.Success/Error/Output` 统一响应。
- **logic 职责**：业务编排、事务控制、错误语、唯一索引预查重（Create/Update 均需）、错误码转换。
- **dao 职责**：SQL 构建、分页、软删、唯一键冲突（MySQL 1062）转业务错误。
- **业务语义不清**（字段含义、状态流转、兼容策略、脏数据处理）必须先确认，禁止臆断。

### 4.2 建议
- 目录沿用既有结构：`internal/types`、`internal/routers`、`internal/server`。
- DTO ↔ `model` 优先 `copier.Copy`，补字段显式处理。
- 测试风格贴近现有 `gotest/sqlmock`。
- 行长 ≤200，优先复用现有中间件与启动逻辑。

---

## 5. 代码生成 Gate 流程（强制，顺序执行，任一步失败即停止）
### 5.1 前置校验（MySQL CLI）
```bash
SHOW TABLES LIKE '<表名>';
SHOW CREATE TABLE <表名>;
SHOW INDEX FROM <表名>;
```
- 任一失败或表不存在 → **禁止**执行 `sponge web http`。

### 5.2 执行生成
```bash
# ⚠️ 必须在 项目根目录下执行（即 /，不是 server/）
# 从 configs/admin.yml 取 database.mysql.dsn，去 ? 前缀作 <dsn主串>
# prefix 取表名前缀（如 t_goods -> t_）
sponge web http \
  --module-name=admin \
  --server-name=admin \
  --project-name=admin \
  --repo-addr= \
  --db-driver=mysql \
  --db-dsn="<dsn主串>;prefix=<prefix>" \
  --db-table=<表名> \
  --embed=true \
  --suited-mono-repo=false \
  --extended-api=false \
  --out=$(pwd)
```

### 5.2.1 关联表识别与文件清理（强制）
**关联表判定条件**（同时满足）：
1. 字段仅包含：`id` + `created_at`/`updated_at`/`deleted_at` + 2个外键字段
2. 唯一索引在外键组合上（如 `uk_project_domain (project_id, domain_id)`）
3. 无业务字段（非外键的 varchar/text/json 等）

**识别为关联表后**，代码生成完成只保留：
- `internal/model/<表名>.go` — GORM 模型
- `internal/dao/<表名>.go` — 数据访问层

**删除以下文件**：
- `internal/handler/<表名>.go` 及 `*_test.go`
- `internal/logic/<表名>.go`
- `internal/types/<表名>_types.go`
- `internal/routers/<表名>.go`
- `internal/ecode/<表名>_http.go`
- `internal/cache/<表名>.go` 及 `*_test.go`
- `web/src/api/<表名>.api.ts`
- `web/src/views/<表名>/`

**原因**：关联表是多对多关系表，业务逻辑由主表的 handler/logic 统一管理，不需要独立的 CRUD 接口。

### 5.3 生成后硬约束同步（`types` + `logic` + `dao`）
| 项 | 要求 |
|----|------|
| `binding` | `internal/types/*_types.go` **禁止**空 `binding:""` |
| 映射规则 | `NOT NULL` 无默认值 → `required`<br>`varchar/char` → `max=n`<br>枚举语义 → `oneof`<br>金额字符串 → `numeric`<br>更新接口 → `omitempty,...` |
| **唯一索引**（`Non_unique=0` 且 `Key_name != PRIMARY`） | **必须**在三处落地：<br>1. `logic.Create`：预查重，命中返回明确业务错误<br>2. `logic.Update`：预查重且排除当前 ID，命中返回明确业务错误<br>3. `dao`：MySQL 1062 统一转可识别业务错误 |
| 唯一索引语义不清（空值、大小写不敏感、联合唯一解释） | **必须先向用户确认** |

### 5.4 列表搜索规则（强制）
| 规则 | 说明 |
|------|------|
| 黑名单字段 | `password/token/secret/salt` 默认禁止搜索，**且不得在列表 DTO 返回**（确需开放先确认） |
| 时间字段 | `created_at/updated_at/deleted_at` 不做通用搜索 |
| 其余字段 | 至少提供一种搜索能力（精确/模糊/范围其一） |
| 搜索白名单 | 必须在 `types`、`logic`、`docs/swagger.yaml` **三处一致** |

### 5.5 测试与文档验收（强制）
- 受影响包测试通过（至少 `types/logic/dao`）
- 唯一索引相关至少包含：
  - `Create` 冲突测试
  - `Update` 冲突测试（排除自身）
- 新增列表筛选能力至少补 1 条 `logic.List` 条件下推测试
- API/路由变更必须执行 `make docs` 并核对产物
- 提交前必须通过 `make test`（建议同时执行 `make ci-lint`）

---

## 6. 前端开发规则（Vue3 + Vite + TS + Element Plus + Pinia）

### 6.1 必选
- **目录结构**：严格遵循 `web/src/` 既有目录，严禁随意新增顶层目录。
- **API 请求**：统一走 `src/api/` 封装，**禁止**在组件/页面直接 `axios`/`fetch`。
- **类型定义**：后端 DTO 变更后，必须同步更新 `src/types/` 对应 TS 接口；API 请求/响应均需强类型。
- **路由懒加载**：所有路由组件必须使用 `defineAsyncComponent` 或 `() => import()` 懒加载。
- **状态管理**：全局状态用 Pinia（`src/store/`），组件内部状态用 `ref`/`reactive`。
- **代码规范**：提交前必须通过 `pnpm lint:eslint`、`pnpm lint:prettier`、`pnpm lint:stylelint`。

### 6.2 建议
- 组件/页面命名：PascalCase（`UserList.vue`）、目录 kebab-case（`user-list/`）。
- 优先复用 `src/components/` 通用组件与 `src/hooks/` 组合式函数。
- 国际化键统一放 `src/lang/`，禁止硬编码中文。
- 环境变量区分 `.env.development` / `.env.production`，**不提交**敏感信息。

---

## 7. 命名与格式统一规范

| 项 | 规范 |
|----|------|
| Go 包名 | 小写单词（`userdao`、`userlogic`） |
| Go 导出标识符 | 驼峰（`UserInfo`、`CreateUser`） |
| Go 格式化 | `gofmt` / `goimports`（CI 强制） |
| TS/JS 变量/函数 | camelCase |
| TS/JS 类型/组件 | PascalCase |
| TS/JS 常量/枚举 | UPPER_SNAKE_CASE |
| Vue 文件/组件 | PascalCase |
| Git 提交信息 | 简洁祈使句（`add user cache metrics`） |
| PR 描述 | 说明变更；有 API 变更时必须包含 `docs/` 更新 |

---

## 8. 配置与部署

- **配置文件**：`server/configs/*.yml`，**不提交敏感信息**（DSN、密钥、Token 等）。
- **部署参考**：`server/deployments/`（binary、Docker、K8S）。
- **前端构建产物**：`web/dist/` 交由 Nginx/静态服务托管。

---

## 9. 核心原则

1. **使用中文**
2. **不清楚不确定的逻辑，一定要问用户，不要自己猜。**
2. **最小改动**：只改当前任务直接相关文件，避免无关重构。
3. **可验证**：每条必选规则都要有对应的自动化检查或测试用例。

---

### 10 CodeGraph 使用规则（必选）

项目根目录存在 `.codegraph/` 时，代码查询**默认走 CodeGraph**，不回到 Read/Grep。

| 场景 | 工具 |
|------|------|
| 查符号调用链、影响范围、依赖关系 | `codegraph_explore` |
| 查单个符号的源码+上下游 | `codegraph_explore` 或 `codegraph_node` |
| 查调用方/被调用方 | `codegraph_callers` / `codegraph_callees` |

**例外**（允许直接 Read/Grep）：
- 无 `.codegraph/` 的项目或子项目
- CodeGraph 结果被截断，需要读取文件其余部分
- 刚新建/修改、尚未重新索引的内容
- 非源码文件：SQL 迁移、Swagger YAML、proto、二进制等

执行要求：
- `codegraph_explore` 已返回完整源码时，**直接读结果**，禁止再 Read 同一文件。
- 文件修改后需要 `codegraph sync` 再查。

---

## 11. 附录：常用文件路径速查

| 场景 | 路径 |
|------|------|
| 后端入口 | `server/cmd/admin/main.go` |
| 后端路由注册 | `server/internal/routers/router.go` |
| 后端统一响应 | `server/internal/handler/response.go` |
| 后端错误码 | `server/internal/ecode/code.go` |
| 后端中间件 | `server/internal/middleware/` |
| 后端配置模板 | `server/configs/admin.yml` |
| 前端入口 | `web/src/main.ts` |
| 前端路由 | `web/src/router/index.ts` |
| 前端 API 封装 | `web/src/api/` |
| 前端全局状态 | `web/src/store/` |
| 前端类型定义 | `web/src/types/` |
| 前端环境变量 | `web/.env.development` / `.env.production` |
