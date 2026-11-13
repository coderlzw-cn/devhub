APP := server
CMD := ./cmd/server
OUT := bin
PKG := gin-enterprise/internal/buildinfo

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w \
	-X $(PKG).Version=$(VERSION) \
	-X $(PKG).Commit=$(COMMIT) \
	-X $(PKG).Time=$(BUILD_TIME)

GOLANGCI := .tools/bin/golangci-lint

# -trimpath 去掉二进制里的本机源码路径。
# -v 按包打印编译进度。Go 没有百分比进度。
# -s -w 在 LDFLAGS 里，用于去掉符号表和调试信息。
BUILD_FLAGS := -trimpath -v -ldflags "$(LDFLAGS)"

.PHONY: build build-all run dev client fmt vet lint tidy sqlc clean

# 编译当前系统的二进制。
build:
	CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(OUT)/$(APP) $(CMD)

# 并行编译 linux、macOS、windows 的常用架构。任一平台失败则整体失败。
build-all:
	@mkdir -p $(OUT)
	@set -e; \
	pids=""; \
	for spec in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64; do \
		os=$${spec%/*}; \
		arch=$${spec#*/}; \
		ext=""; \
		if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		echo "开始编译 $$os/$$arch"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build $(BUILD_FLAGS) -o $(OUT)/$(APP)-$$os-$$arch$$ext $(CMD) & \
		pids="$$pids $$!"; \
	done; \
	status=0; \
	for pid in $$pids; do \
		wait $$pid || status=1; \
	done; \
	exit $$status

run:
	go run -ldflags "$(LDFLAGS)" $(CMD)

# 启动后端 Air 热重载；client 模式只启动前端 Vite。
dev:
	@if [ "$(filter client,$(MAKECMDGOALS))" ]; then \
		cd web/frontend && pnpm dev; \
	else \
		air; \
	fi

# make dev client 的纯前端模式标记。
client:

# 用 go fmt 格式化全部 Go 源码。
fmt:
	go fmt ./...

# 用 go vet 检查常见错误。
vet:
	go vet ./...

# 用 .tools/bin/golangci-lint 按 .golangci.yml 检查。
lint:
	$(GOLANGCI) run ./...

# 整理 go.mod 和 go.sum，补上缺失依赖并去掉未使用的依赖。
tidy:
	go mod tidy

# 从 SQL 查询和数据库 schema 生成类型安全的 Go 代码。
sqlc:
	sqlc generate

# 删除编译产物 bin。
clean:
	rm -rf $(OUT)
