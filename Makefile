.PHONY: help server-dev server-build web-install web-dev web-build web-generate

help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# 后端命令
server-dev: ## 运行后端开发服务器
	cd host-server && go run cmd/server/main.go

server-build: ## 构建后端二进制文件
	cd host-server && go build -o bin/server cmd/server/main.go

swagger: ## 生成 Swagger/OpenAPI 文档
	cd host-server && swag init -g cmd/server/main.go -o docs

swagger-fmt: ## 格式化 Swagger 注释
	cd host-server && swag fmt

# 前端命令
web-install: ## 安装前端依赖
	cd host-web && pnpm install

web-dev: ## 运行前端开发服务器
	cd host-web && pnpm dev

web-build: ## 构建前端生产版本
	cd host-web && pnpm build

web-generate: ## 从 Swagger 生成 API 客户端代码
	cd host-web && pnpm generate:api

# 全栈命令
install: ## 安装所有依赖
	@echo "安装后端依赖..."
	cd host-server && go mod download
	@echo "安装前端依赖..."
	cd host-web && pnpm install

dev: ## 同时运行前后端开发服务器 (需要两个终端)
	@echo "请在两个终端分别运行:"
	@echo "  终端 1: make server-dev"
	@echo "  终端 2: make web-dev"

build: server-build web-build ## 构建前后端

# 插件命令
plugin-proto: ## 使用 buf 生成插件 protobuf 代码
	cd proto && buf generate

plugin-proto-lint: ## 检查 proto 文件
	cd proto && buf lint

plugin-proto-breaking: ## 检查 proto 文件的破坏性变更
	cd proto && buf breaking --against '.git#branch=main'

plugin-desensitization: ## 构建数据脱敏插件
	cd plugins/desensitization && go build -o ../../host-server/bin/plugins/desensitization

plugin-build: plugin-desensitization ## 构建所有插件

plugin-clean: ## 清理插件二进制文件
	rm -rf host-server/bin/plugins/*

plugin-test: ## 运行插件测试
	cd plugins/desensitization/impl && go test -v

plugin-example: ## 运行插件示例
	cd plugins/desensitization/example && go run main.go

