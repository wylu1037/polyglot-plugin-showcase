.PHONY: help install dev build generate

help: ## Show help information
	@echo "📋 Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Installation
install: ## Install all dependencies
	@echo "📦 Installing dependencies..."
	@echo "  → Backend dependencies"
	@cd host-server && go mod download
	@echo "  → Frontend dependencies"
	@cd host-web && pnpm install
	@echo "✅ Installation complete"

# Development
dev: ## Start development servers (requires two terminals)
	@echo "🚀 Starting development environment"
	@echo "  Please run in two terminals:"
	@echo "  Terminal 1: cd host-server && go run cmd/server/main.go"
	@echo "  Terminal 2: cd host-web && pnpm dev"

server-dev: ## Start backend server only
	@echo "🔧 Starting backend server..."
	@cd host-server && go run cmd/server/main.go

web-dev: ## Start frontend server only
	@echo "🎨 Starting frontend server..."
	@cd host-web && pnpm dev

build: ## Build backend, frontend and all plugins
	@echo "📦 Building..."
	@echo "  → Building backend"
	@cd host-server && go build -o bin/server cmd/server/main.go
	@echo "  → Building frontend"
	@cd host-web && pnpm build
	@echo "  → Building plugins"
	@cd plugins/desensitization && go build -o ../../host-server/bin/plugins/builtin/data-processing/desensitization/v1.0.0/darwin_arm64/plugin .
	@cd plugins/dpanonymizer && go build -o ../../host-server/bin/plugins/builtin/data-processing/dpanonymizer/v1.0.0/darwin_arm64/plugin .
	@cd plugins/converter && go build -o ../../host-server/bin/plugins/builtin/data-processing/converter/v1.0.0/darwin_arm64/plugin .
	@echo "🏁 Build complete"

plugin-build: ## Build all plugins only
	@echo "🔌 Building plugins..."
	@mkdir -p host-server/bin/plugins/builtin/data-processing/desensitization/v1.0.0/darwin_arm64
	@mkdir -p host-server/bin/plugins/builtin/data-processing/dpanonymizer/v1.0.0/darwin_arm64
	@mkdir -p host-server/bin/plugins/builtin/data-processing/converter/v1.0.0/darwin_arm64
	@cd plugins/desensitization && go build -o ../../host-server/bin/plugins/builtin/data-processing/desensitization/v1.0.0/darwin_arm64/plugin .
	@cd plugins/dpanonymizer && go build -o ../../host-server/bin/plugins/builtin/data-processing/dpanonymizer/v1.0.0/darwin_arm64/plugin .
	@cd plugins/converter && go build -o ../../host-server/bin/plugins/builtin/data-processing/converter/v1.0.0/darwin_arm64/plugin .
	@echo "✅ Plugins built successfully"

generate: ## Generate all code (API docs + Frontend client + Plugin protocol)
	@echo "⚙️  Generating code..."
	@echo "  → Generating Swagger docs"
	@cd host-server && swag init -g cmd/server/main.go -o docs
	@echo "  → Generating frontend API client"
	@cd host-web && pnpm generate:api
	@echo "  → Generating plugin protobuf code"
	@cd proto && buf generate
	@echo "🏁 Code generation complete"

swagger: ## Generate Swagger documentation only
	@echo "📖 Generating Swagger docs..."
	@cd host-server && swag init -g cmd/server/main.go -o docs
	@echo "✅ Swagger docs generated"

