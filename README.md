Gin 企业项目模板。HTTP 服务使用 Gin 和 Huma，配置来自环境变量，日志使用 zap。

## 环境准备

需要 Go（版本见 `go.mod`）、Node.js、pnpm、PostgreSQL、Make。Air 仅用于 `make dev` 热重载；安装后需确保 `$(go env GOPATH)/bin` 在 `PATH` 中。

```bash
go mod download
cd web/frontend && pnpm install && cd ../..
go install github.com/air-verse/air@latest
```

创建并填写 `.env` 中的 PostgreSQL 连接配置后才能启动后端。开发环境的客户端和 JWT 配置有默认值；生产环境必须自行设置安全的认证密钥。首次运行或数据库结构有更新时，先按顺序应用 `migrations/` 下的 SQL 文件；服务启动时不会自动迁移数据库。

## 本地启动

在两个终端分别启动服务：

```bash
make dev          # Gin 后端，Air 热重载，默认 http://127.0.0.1:8080
make dev client   # Vue 前端，Vite，默认 http://127.0.0.1:5173
```

Vite 会将 `/api` 请求代理到本机 Gin。只开发后端时运行 `make dev`；只开发前端时运行 `make dev client`。

直接启动 Go 服务（不启用 Air）：

```bash
make run
```

## 前端构建

前端源码位于 `web/frontend/`。构建产物输出到 `web/dist/` 并由 Go 嵌入二进制；发布前先构建前端：

```bash
cd web/frontend
pnpm install
pnpm build
cd ../..
make build
```

`node_modules` 不纳入版本控制。

## 配置

环境变量以 `APP_` 开头，下划线对应配置层级。例如 `APP_SERVER_PORT` 对应 `server.port`。

| 变量 | 说明 |
| --- | --- |
| `APP_APP_ENV` | `development` 输出彩色控制台日志，`production` 输出 JSON |
| `APP_SERVER_HOST` | 监听地址 |
| `APP_SERVER_PORT` | 监听端口 |
| `APP_POSTGRES_HOST` | PostgreSQL 地址 |
| `APP_POSTGRES_PORT` | PostgreSQL 端口，默认 5432 |
| `APP_POSTGRES_USER` | PostgreSQL 用户 |
| `APP_POSTGRES_PASSWORD` | PostgreSQL 密码 |
| `APP_POSTGRES_DATABASE` | PostgreSQL 数据库名 |
| `APP_POSTGRES_SSLMODE` | SSL 模式，默认 `require` |
| `APP_AUTH_CLIENT_ID` | 签发令牌的客户端 ID。开发环境默认为 `gin-enterprise` |
| `APP_AUTH_CLIENT_SECRET` | 客户端密钥。开发环境有默认值，生产环境必填 |
| `APP_AUTH_JWT_SECRET` | JWT 签名密钥，至少 32 字节。生产环境必填 |
| `APP_AUTH_ISSUER` | 令牌签发方，默认 `gin-enterprise` |
| `APP_AUTH_JWT_TTL` | 访问令牌有效期，单位秒，默认 3600 |

## 构建与检查

```bash
make build       # 编译当前系统的二进制
make build-all   # 并行编译 linux、macOS、windows 的常用架构
make fmt         # 用 go fmt 格式化全部 Go 源码
make vet         # 用 go vet 检查常见错误
make lint        # 用 .tools/bin/golangci-lint 按 .golangci.yml 检查
make tidy        # 整理 go.mod 和 go.sum，补上缺失依赖并去掉未使用的依赖
make sqlc        # 根据 SQL 查询和 schema 重新生成数据库访问代码（需安装 sqlc）
make clean       # 删除 bin
```

`make lint` 使用 `.tools/bin/golangci-lint`。可将 `golangci-lint` 安装到 PATH，或安装到项目本地目录：

```bash
mkdir -p .tools/bin
GOBIN="$PWD/.tools/bin" go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```

前端检查在 `web/frontend/` 执行：

```bash
pnpm type-check
pnpm lint
```

## 列表查询参数

OpenAPI JSON：`http://127.0.0.1:8080/openapi/json`，可作为 Orval 输入；服务重启后包含最新参数定义。

所有列表使用 `page`（默认 1）、`page_size`（默认 10，最大 100）和可选 `offset`（最小 0）。显式传入 `offset` 时优先使用该偏移，包括 `offset=0`，响应页码按实际偏移计算。

- `/api/categories`：`name`、`slug`、`description` 分别模糊匹配对应字段。
- `/api/sites`：`title`、`url`、`description` 分别模糊匹配对应字段；`category_id`、`user_id` 为 UUID 精确筛选；`is_pinned` 为可选布尔筛选，未传表示不限，`false` 表示明确筛选 false。

模糊查询为不区分大小写的包含匹配，多个筛选条件取交集；空文本不筛选，`total` 为筛选后的总数。`%`、`_` 按普通字符匹配。query 参数与数据库字段保持同名，使用 snake_case。

```http
GET /api/sites?title=github&url=github.com&is_pinned=false&page=1&page_size=20
GET /api/categories?name=前端&page=1&page_size=10
```

## 部署

Linux 不使用容器时，用 systemd。服务文件在 `deployments/systemd/gin-enterprise.service`。将二进制和 `.env` 放到 `/opt/gin-enterprise/`，再执行：

```bash
sudo cp deployments/systemd/gin-enterprise.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now gin-enterprise
```

数据库访问使用 sqlc 生成代码。新增数据库查询后执行 `sqlc generate`；新环境部署前需按发布流程应用 `migrations/` 下的 SQL 迁移，服务启动时不会自动改动数据库结构。

## 目录职责

```text
.
├── cmd/
│   └── server/              # 服务端程序入口，负责启动应用
├── internal/                # 本模块私有代码，外部项目不可导入
│   ├── app/                 # 应用初始化、依赖组装和进程生命周期管理
│   ├── buildinfo/           # 构建版本、提交号和构建时间
│   ├── config/              # 配置加载、校验和解析
│   ├── constant/            # 内部共用常量
│   ├── database/            # 数据库连接创建、配置和关闭
│   ├── db/sqlc/             # sqlc 生成的类型安全查询代码
│   ├── dto/
│   │   ├── request/         # HTTP 请求参数及请求体 DTO
│   │   └── response/        # HTTP 响应 DTO
│   ├── errors/              # 业务错误定义及错误转换约定
│   ├── handler/
│   │   └── http/
│   │       ├── middleware/  # 请求 ID、日志、超时、恢复等横切处理
│   │       ├── response/    # HTTP 响应体和统一错误响应格式
│   │       ├── router/      # 路由注册及接口与服务层绑定
│   │       └── websocket/   # WebSocket 连接处理、管理和消息分发
│   ├── logger/              # 日志初始化与配置封装
│   ├── model/               # 数据库实体和内部领域模型
│   ├── pkg/                 # 仅供 internal 内部使用的共享实现
│   ├── repository/          # 数据访问封装
│   │   └── queries/         # 手写 SQL，由 sqlc 生成 Go 查询代码
│   ├── scheduler/           # 定时任务注册、启动和停止
│   ├── service/             # 业务逻辑和用例编排
│   └── worker/              # 异步后台任务执行逻辑
├── pkg/                     # 可供其他项目复用的公共组件
│   ├── response/            # 通用响应组件
│   └── utils/               # 不依赖业务域的通用工具
├── api/
│   └── openapi/             # OpenAPI 接口规范文件
├── configs/                 # 可提交的配置模板及非敏感默认配置
├── deployments/             # 部署相关配置
│   ├── docker/              # Docker 镜像构建及容器运行配置
│   ├── k8s/                 # Kubernetes 部署清单
│   └── systemd/             # Linux systemd 服务配置
├── docs/                    # 架构、接口和运维文档
├── migrations/              # 数据库结构变更脚本
├── scripts/                 # 开发、构建、检查和运维辅助脚本
├── storage/                 # 本地运行时文件，不提交运行时数据
│   ├── logs/                # 本地运行日志
│   └── uploads/             # 本地上传文件
├── test/                    # 测试代码
│   ├── e2e/                 # 验证完整服务流程的端到端测试
│   └── integration/         # 验证多个组件协作的集成测试
```

空目录中的 `.gitkeep` 仅用于让 Git 保留目录；目录开始存放实际文件后可以移除对应占位文件。`internal` 下的代码只在本模块内使用；需要跨项目复用的组件放在 `pkg`。
