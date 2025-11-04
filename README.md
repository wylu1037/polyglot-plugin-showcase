# Polyglot Plugin Showcase

A dynamic plugin system demonstration based on [hashicorp/go-plugin](https://github.com/hashicorp/go-plugin).

## 🎯 Project Goals

Demonstrate how to build a complete plugin-based application system using `go-plugin`:
- **Backend (host-server)**: Go + [Echo](https://github.com/labstack/echo) - Plugin host server
- **Frontend (host-web)**: React + [React Router](https://reactrouter.com) - Plugin management UI
- **Plugins**: Independent Go binaries communicating via gRPC/net-rpc

## 📁 Project Structure

```
polyglot-plugin-showcase/
├── host-server/       # Go backend server (Echo)
│   ├── main.go        # Entry point
│   ├── plugin/        # Plugin manager
│   └── api/           # RESTful API handlers
├── host-web/          # React frontend (React Router)
│   ├── src/
│   │   ├── routes/    # Route definitions
│   │   ├── components/# UI components
│   │   └── api/       # API client
│   └── package.json
├── plugins/           # Plugin implementations
│   ├── example-plugin/
│   └── another-plugin/
├── proto/             # gRPC protocol definitions (if using gRPC)
└── Makefile           # Build scripts
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- make (optional)

### Start Backend

```bash
cd host-server
go mod download
go run main.go
# Server runs at http://localhost:8080
```

### Start Frontend

```bash
cd host-web
npm install
npm run dev
# Frontend runs at http://localhost:5173
```

### Build Plugins

```bash
cd plugins/example-plugin
go build -o example-plugin
# Place the binary in host-server/plugins/ directory
```

## 🔌 How Plugins Work

```
┌─────────────────┐
│  React Frontend │  HTTP Request
│  (host-web)     ├──────────┐
└─────────────────┘          │
                             ▼
                    ┌────────────────┐
                    │  Echo Server   │
                    │  (host-server) │
                    └────────┬───────┘
                             │ go-plugin (gRPC/net-rpc)
                    ┌────────▼────────┐
                    │  Plugin Process │
                    │  (separate proc)│
                    └─────────────────┘
```

**Key Concepts**:
1. **Plugins are separate processes**: Each plugin runs in its own Go process
2. **Communication**: Inter-process communication via gRPC or net-rpc
3. **Lifecycle management**: host-server manages plugin start/stop/restart
4. **Protocol definition**: Plugin capabilities defined via interfaces (in `proto/` or code)

## 📋 Features

- [ ] Dynamic plugin loading/unloading
- [ ] Plugin health monitoring
- [ ] Plugin communication logging
- [ ] Web UI management interface
- [ ] Plugin hot-reload
- [ ] Plugin dependency management

## 🛠️ Tech Stack

### Backend
- **Web Framework**: [Echo](https://echo.labstack.com/) - High-performance Go web framework
- **Plugin System**: [go-plugin](https://github.com/hashicorp/go-plugin) - HashiCorp's plugin library
- **Communication Protocol**: gRPC / net-rpc

### Frontend
- **Framework**: React 18
- **Routing**: [React Router v7](https://reactrouter.com)
- **Build Tool**: Vite
- **UI Library**: (TBD - suggest shadcn/ui or Ant Design)

## 📖 Development Guide

### Creating a New Plugin

1. Create a new directory under `plugins/`
2. Implement the plugin interface (refer to `plugins/example-plugin/`)
3. Build as a standalone binary
4. Register with host-server via API or configuration file

### API Design

```
GET    /api/plugins          # List all plugins
POST   /api/plugins/load     # Load a plugin
DELETE /api/plugins/:id      # Unload a plugin
GET    /api/plugins/:id/info # Plugin details
POST   /api/plugins/:id/call # Call plugin method
```

## ⚠️ Important Notes

1. **Security**: Plugins run in separate processes, but still need to validate plugin sources
2. **Resource Management**: Ensure plugin processes are properly closed to avoid zombie processes
3. **Error Handling**: Plugin crashes should not affect the main server
4. **Version Compatibility**: Define clear plugin interface versions

## 🤔 Design Considerations

**Why go-plugin?**
- Process isolation: Plugin crashes don't affect main program
- Language agnostic: Theoretically supports plugins in any language (if protocol is implemented)
- Battle-tested: Used by Terraform, Vault, and other HashiCorp projects

**Architecture Trade-offs**:
- ✅ Strong isolation, high stability
- ❌ Process communication overhead, deployment complexity

## 📚 References

- [go-plugin Official Documentation](https://github.com/hashicorp/go-plugin)
- [Echo Framework Documentation](https://echo.labstack.com/docs)
- [React Router Documentation](https://reactrouter.com)
- [gRPC Go Quick Start](https://grpc.io/docs/languages/go/quickstart/)

## 📝 License

MIT

---

**Development Status**: 🚧 Initializing...
