# 🔌 Polyglot Plugin Showcase

A production-ready dynamic plugin system demonstration based on [HashiCorp go-plugin](https://github.com/hashicorp/go-plugin), featuring a complete plugin management platform with RESTful API, modern web UI, and interactive API documentation.

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue.svg)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18-blue.svg)](https://reactjs.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## ✨ Features

- 🚀 **Dynamic Plugin Management** - Load, unload, activate, and deactivate plugins at runtime
- 🔒 **Process Isolation** - Each plugin runs in a separate process for maximum stability
- 📡 **gRPC Communication** - High-performance inter-process communication
- 🌐 **RESTful API** - Complete plugin management API with Echo framework
- 💻 **Modern Web UI** - React-based management interface with shadcn/ui
- 📚 **Interactive API Docs** - Beautiful API documentation powered by Scalar
- 🗄️ **Database Persistence** - Plugin metadata stored in PostgreSQL
- 🔄 **Auto-reload** - Automatically load active plugins on server startup
- 🛡️ **Error Handling** - Comprehensive error handling with structured responses
- 🧪 **Example Plugins** - Data desensitization and differential privacy anonymization plugins

## 🎯 Project Goals

Demonstrate how to build a complete, production-ready plugin-based application system:

- **Backend (host-server)**: Go + Echo + PostgreSQL - Plugin host server with RESTful API
- **Frontend (host-web)**: React + React Router + shadcn/ui - Plugin management UI
- **Plugins**: Independent Go binaries communicating via gRPC
- **Documentation**: Interactive API documentation with Scalar

## 📁 Project Structure

```
polyglot-plugin-showcase/
├── host-server/              # Backend server (Go + Echo)
│   ├── cmd/server/           # Application entry point
│   ├── app/
│   │   ├── common/           # Common utilities (errors, responses)
│   │   ├── database/         # Database connection and migrations
│   │   ├── modules/plugins/  # Plugin management module
│   │   │   ├── controller/   # HTTP handlers
│   │   │   ├── service/      # Business logic
│   │   │   └── repository/   # Data access layer
│   │   └── router/           # Route definitions
│   ├── internal/
│   │   ├── bootstrap/        # Application bootstrap (server, docs)
│   │   └── plugin/           # Plugin manager and registry
│   ├── config/               # Configuration management
│   ├── docs/                 # Auto-generated API documentation
│   └── bin/plugins/          # Compiled plugin binaries
│
├── host-web/                 # Frontend (React + Vite)
│   ├── src/
│   │   ├── routes/           # Page components
│   │   ├── components/ui/    # UI components (shadcn/ui)
│   │   ├── lib/              # Utilities
│   │   └── router.tsx        # Route configuration
│   └── package.json
│
├── plugins/                  # Plugin implementations
│   ├── desensitization/      # Data desensitization plugin
│   │   ├── main.go           # Plugin entry point
│   │   ├── adapter/          # Plugin adapter (implements common interface)
│   │   ├── impl/             # Business logic implementation
│   │   └── example/          # Standalone example
│   └── dpanonymizer/         # Differential privacy anonymization plugin
│       ├── main.go           # Plugin entry point
│       ├── adapter/          # Plugin adapter
│       ├── impl/             # DP algorithms implementation
│       └── example/          # Standalone example
│
├── proto/                    # Protocol definitions
│   ├── common/               # Common plugin interface (gRPC)
│   │   ├── plugin.proto      # Protocol definition
│   │   ├── grpc.go           # go-plugin integration
│   │   └── interface.go      # Go interface
│   └── desensitization/      # Plugin-specific protocols
│
└── Makefile                  # Build automation
```

## 🚀 Quick Start

### Prerequisites

- **Go** 1.25+ - [Download](https://golang.org/dl/)
- **Node.js** 18+ - [Download](https://nodejs.org/)
- **PostgreSQL** 14+ - [Download](https://www.postgresql.org/download/)
- **pnpm** (optional) - `npm install -g pnpm`
- **buf** (optional, for proto generation) - [Install](https://buf.build/docs/installation)

### 1. Database Setup

```bash
# Create database
createdb polyglot_plugin

# Or using psql
psql -U postgres
CREATE DATABASE polyglot_plugin;
```

### 2. Start Backend

```bash
# Clone repository
git clone https://github.com/wylu1037/polyglot-plugin-showcase.git
cd polyglot-plugin-showcase

# Configure database
cd host-server
cp config.example.yaml config.yaml
# Edit config.yaml with your database credentials

# Install dependencies
go mod download

# Run server
make server-dev
# Or: go run cmd/server/main.go

# Server runs at http://localhost:8080
# API Docs at http://localhost:8080/docs
```

### 3. Start Frontend

```bash
cd host-web

# Install dependencies
pnpm install
# Or: npm install

# Start development server
pnpm dev
# Or: npm run dev

# Frontend runs at http://localhost:5173
```

### 4. Build and Install Plugins

```bash
# Build all plugins
make plugin-build

# Or build individually:
make plugin-desensitization    # Data desensitization plugin
make plugin-dpanonymizer       # Differential privacy plugin

# Run plugin examples:
make plugin-example-desensitization
make plugin-example-dpanonymizer

# Install via API (see API documentation)
```

## 🔌 How It Works

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     React Frontend (host-web)                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Plugin List  │  │Plugin Detail │  │Plugin Discover│      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└────────────────────────────┬────────────────────────────────┘
                             │ HTTP/REST API
                             ▼
┌─────────────────────────────────────────────────────────────┐
│              Echo Server (host-server)                       │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Controller → Service → Repository → Database        │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Plugin Manager (Registry + Lifecycle Management)    │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────────┬────────────────────────────────┘
                             │ go-plugin (gRPC)
                ┌────────────┼────────────┐
                ▼            ▼            ▼
        ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
        │  Plugin A   │ │  Plugin B   │ │  Plugin C   │
        │  (Process)  │ │  (Process)  │ │  (Process)  │
        └─────────────┘ └─────────────┘ └─────────────┘
```

### Key Concepts

1. **Process Isolation**: Each plugin runs in its own Go process
2. **gRPC Communication**: High-performance inter-process communication
3. **Plugin Registry**: Centralized plugin management and discovery
4. **Lifecycle Management**: Install → Activate → Execute → Deactivate → Uninstall
5. **Common Interface**: All plugins implement the same `PluginInterface`

### Plugin Lifecycle

```
┌─────────────┐
│  Inactive   │ ← Initial state
└──────┬──────┘
       │ Install
       ▼
┌─────────────┐
│  Installed  │ ← Binary downloaded, metadata stored
└──────┬──────┘
       │ Activate
       ▼
┌─────────────┐
│   Active    │ ← Process started, ready to serve
└──────┬──────┘
       │ Deactivate
       ▼
┌─────────────┐
│  Installed  │
└──────┬──────┘
       │ Uninstall
       ▼
┌─────────────┐
│   Removed   │
└─────────────┘
```

## 📖 API Documentation

### Interactive Documentation

Visit **http://localhost:8080/docs** for interactive API documentation powered by Scalar.

### Key Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/plugins/install` | Install a new plugin |
| `GET` | `/api/plugins` | List all plugins |
| `GET` | `/api/plugins/{id}` | Get plugin details |
| `POST` | `/api/plugins/{id}/activate` | Activate a plugin |
| `POST` | `/api/plugins/{id}/deactivate` | Deactivate a plugin |
| `DELETE` | `/api/plugins/{id}` | Uninstall a plugin |
| `POST` | `/api/plugins/{id}/call` | Execute plugin method |

### Example: Install Plugin

```bash
curl -X POST http://localhost:8080/api/plugins/install \
  -H "Content-Type: application/json" \
  -d '{
    "name": "desensitization",
    "version": "1.0.0",
    "type": "grpc",
    "download_url": "https://example.com/plugins/desensitization_v1.0.0",
    "description": "Data desensitization plugin"
  }'
```

### Example: Call Plugin Method

```bash
curl -X POST http://localhost:8080/api/plugins/1/call \
  -H "Content-Type: application/json" \
  -d '{
    "method": "DesensitizeName",
    "params": {
      "data": "张三"
    }
  }'
```

## 🛠️ Development

### Project Commands

```bash
# Backend
make server-dev          # Start development server
make server-build        # Build production binary
make swagger             # Generate API documentation

# Frontend
make web-dev             # Start frontend dev server
make web-build           # Build production bundle
make web-generate        # Generate API client from OpenAPI

# Plugins
make plugin-build                    # Build all plugins
make plugin-desensitization          # Build desensitization plugin
make plugin-dpanonymizer             # Build differential privacy plugin
make plugin-proto                    # Generate plugin protobuf code
make plugin-test                     # Run all plugin tests
make plugin-example-desensitization  # Run desensitization example
make plugin-example-dpanonymizer     # Run differential privacy example

# Full stack
make install             # Install all dependencies
make build               # Build everything
```

### Creating a New Plugin

1. **Create plugin directory**:
```bash
mkdir -p plugins/my-plugin/{impl,adapter}
cd plugins/my-plugin
```

2. **Implement the plugin interface**:
```go
// adapter/adapter.go
package adapter

import "github.com/wylu1037/polyglot-plugin-showcase/proto/common"

type MyPluginAdapter struct{}

func NewMyPluginAdapter() *MyPluginAdapter {
    return &MyPluginAdapter{}
}

func (a *MyPluginAdapter) GetMetadata() (*common.MetadataResponse, error) {
    return &common.MetadataResponse{
        Name:            "my-plugin",
        Version:         "1.0.0",
        Description:     "My awesome plugin",
        Methods:         []string{"MyMethod"},
        ProtocolVersion: 1,
    }, nil
}

func (a *MyPluginAdapter) Execute(method string, params map[string]string) (*common.ExecuteResponse, error) {
    switch method {
    case "MyMethod":
        // Your logic here
        result := "Hello, " + params["name"]
        return &common.ExecuteResponse{
            Success: true,
            Result:  &result,
        }, nil
    default:
        errMsg := "unknown method: " + method
        return &common.ExecuteResponse{
            Success: false,
            Error:   &errMsg,
        }, nil
    }
}
```

3. **Create main.go**:
```go
// main.go
package main

import (
    "github.com/hashicorp/go-plugin"
    "github.com/wylu1037/polyglot-plugin-showcase/plugins/my-plugin/adapter"
    "github.com/wylu1037/polyglot-plugin-showcase/proto/common"
)

func main() {
    plugin.Serve(&plugin.ServeConfig{
        HandshakeConfig: common.Handshake,
        Plugins: map[string]plugin.Plugin{
            "my-plugin": &common.PluginGRPCPlugin{
                Impl: adapter.NewMyPluginAdapter(),
            },
        },
        GRPCServer: plugin.DefaultGRPCServer,
    })
}
```

4. **Build and install**:
```bash
go build -o ../../host-server/bin/plugins/my-plugin/my-plugin_v1.0.0
```

### Updating API Documentation

After modifying controller comments:

```bash
cd host-server
swag init -g cmd/server/main.go -o docs
# Or: make swagger
```

## 🧪 Testing

### Run Tests

```bash
# Backend tests
cd host-server
go test ./...

# Plugin tests
cd plugins/desensitization
go test ./impl/...

# Frontend tests
cd host-web
pnpm test
```

### Test Plugin Standalone

```bash
cd plugins/desensitization/example
go run main.go
```

## 🛡️ Security Considerations

1. **Plugin Verification**: Always verify plugin checksums before installation
2. **Process Isolation**: Plugins run in separate processes, limiting blast radius
3. **Resource Limits**: Consider implementing resource limits for plugin processes
4. **Input Validation**: Validate all plugin inputs and outputs
5. **Authentication**: Add authentication/authorization for plugin management APIs

## 🎨 Tech Stack

### Backend
- **Framework**: [Echo](https://echo.labstack.com/) - High-performance Go web framework
- **Plugin System**: [go-plugin](https://github.com/hashicorp/go-plugin) - HashiCorp's plugin library
- **Database**: PostgreSQL + [GORM](https://gorm.io/)
- **DI Container**: [Uber Fx](https://uber-go.github.io/fx/)
- **API Docs**: [Scalar](https://scalar.com/) + [Swaggo](https://github.com/swaggo/swag)
- **Protocol**: gRPC + Protocol Buffers

### Frontend
- **Framework**: React 18
- **Routing**: [React Router v7](https://reactrouter.com)
- **Build Tool**: Vite
- **UI Library**: [shadcn/ui](https://ui.shadcn.com/)
- **Styling**: Tailwind CSS
- **API Client**: [TanStack Query](https://tanstack.com/query)
- **Code Generation**: [Kubb](https://kubb.dev/)

### DevOps
- **Build**: Make, Go build, Vite
- **Proto**: [Buf](https://buf.build/)
- **Testing**: Go testing, Vitest

## 📚 References

- [go-plugin Official Documentation](https://github.com/hashicorp/go-plugin)
- [Echo Framework Documentation](https://echo.labstack.com/docs)
- [React Router Documentation](https://reactrouter.com)
- [gRPC Go Quick Start](https://grpc.io/docs/languages/go/quickstart/)
- [Scalar API Documentation](https://github.com/scalar/scalar)
- [Protocol Buffers](https://protobuf.dev/)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details

## 🙏 Acknowledgments

- [HashiCorp](https://www.hashicorp.com/) for the excellent go-plugin library
- [Echo](https://echo.labstack.com/) for the high-performance web framework
- [Scalar](https://scalar.com/) for the beautiful API documentation

---

**Development Status**: ✅ Production Ready

**Last Updated**: 2025-11-18

## 🔐 Featured Plugins

### 1. Data Desensitization Plugin
Provides various data masking methods for sensitive information:
- Name desensitization
- Phone number masking
- ID card number masking
- Email address masking
- Bank card number masking
- Address masking

### 2. Differential Privacy Anonymization Plugin
Implements Google's Differential Privacy library for privacy-preserving data analysis:
- **Noise Addition**: Laplace and Gaussian noise mechanisms
- **Aggregations**: Differentially private count, sum, mean, and variance
- **Privacy Guarantees**: Configurable ε (epsilon) and δ (delta) parameters
- **Use Cases**: Statistical reporting, data analytics, machine learning

See [plugins/dpanonymizer/README.md](plugins/dpanonymizer/README.md) for detailed documentation.
