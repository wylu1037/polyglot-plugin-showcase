# Protocol Buffers 定义

本目录包含所有插件的 Protocol Buffers 定义和生成的代码。

## 📁 目录结构

```
proto/
├── buf.yaml              # Buf 配置文件
├── buf.gen.yaml          # Buf 代码生成配置
├── go.mod                # Go 模块定义
├── desensitization/      # 数据脱敏插件的 proto 定义
│   ├── desensitizer.proto       # protobuf 定义
│   ├── desensitizer.pb.go       # 生成的 Go 消息代码
│   ├── desensitizer_grpc.pb.go  # 生成的 gRPC 服务代码
│   ├── interface.go             # Go 接口定义
│   └── grpc.go                  # gRPC 客户端/服务端实现
├── another_plugin/       # 其他插件(示例)
│   └── ...
└── README.md             # 本文件
```

## 🔧 使用 Buf 管理 Proto

我们使用 [Buf](https://buf.build/) 来管理 Protocol Buffers 文件和代码生成。Buf 是一个现代化的 protobuf 工具,提供了更好的开发体验。

### 为什么使用 Buf?

- ✅ **更简单**: 不需要手动管理 protoc 和各种插件
- ✅ **更快**: 增量编译,只编译修改的文件
- ✅ **更安全**: 内置 lint 和 breaking change 检测
- ✅ **更现代**: 支持远程插件,无需本地安装
- ✅ **更一致**: 团队成员使用相同的配置和版本

### 安装 Buf

```bash
# macOS
brew install bufbuild/buf/buf

# Linux
curl -sSL "https://github.com/bufbuild/buf/releases/latest/download/buf-$(uname -s)-$(uname -m)" -o /usr/local/bin/buf
chmod +x /usr/local/bin/buf

# 或使用 Go 安装
go install github.com/bufbuild/buf/cmd/buf@latest
```

## 🚀 常用命令

### 生成代码

```bash
# 在 proto 目录下
buf generate

# 或在项目根目录使用 Makefile
make plugin-proto
```

### 检查 Proto 文件

```bash
# Lint 检查
buf lint

# 或使用 Makefile
make plugin-proto-lint
```

### 检查破坏性变更

```bash
# 对比当前分支和 main 分支
buf breaking --against '.git#branch=main'

# 或使用 Makefile
make plugin-proto-breaking
```

## 📝 配置文件说明

### buf.yaml

定义了 buf 模块和 lint/breaking 规则:

```yaml
version: v2
modules:
  - path: .
lint:
  use:
    - STANDARD
  except:
    # 对于内部插件系统,这些规则可以放宽
    - PACKAGE_VERSION_SUFFIX
    - SERVICE_SUFFIX
    - RPC_REQUEST_RESPONSE_UNIQUE
    - RPC_REQUEST_STANDARD_NAME
    - RPC_RESPONSE_STANDARD_NAME
breaking:
  use:
    - FILE
```

### buf.gen.yaml

定义了代码生成配置:

```yaml
version: v2
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: github.com/wylu1037/polyglot-plugin-showcase/proto
plugins:
  - remote: buf.build/protocolbuffers/go
    out: .
    opt:
      - paths=source_relative
  - remote: buf.build/grpc/go
    out: .
    opt:
      - paths=source_relative
```

## 🔌 添加新的插件接口

1. 在 `proto/` 目录下创建新的插件目录,如 `your_plugin/`
2. 在该目录下创建 `.proto` 文件
3. 定义服务和消息类型
4. 运行 `buf generate` 生成代码
5. 在同一目录下创建对应的 Go 接口和 gRPC 实现

### 示例

```protobuf
// proto/your_plugin/your_plugin.proto
syntax = "proto3";

package your_plugin;

option go_package = "github.com/wylu1037/polyglot-plugin-showcase/proto/your_plugin";

service YourPlugin {
  rpc DoSomething(YourRequest) returns (YourResponse);
}

message YourRequest {
  string data = 1;
}

message YourResponse {
  string result = 1;
}
```

### 目录结构

```
proto/
├── desensitization/      # 脱敏插件
├── your_plugin/          # 你的新插件
│   ├── your_plugin.proto
│   ├── your_plugin.pb.go
│   ├── your_plugin_grpc.pb.go
│   ├── interface.go
│   └── grpc.go
└── shared/               # 真正共享的类型(可选)
    └── common.proto
```

## 🛠️ 开发工作流

1. **创建插件目录**
   ```bash
   mkdir proto/your_plugin
   ```

2. **编写 proto 文件**
   ```bash
   vim proto/your_plugin/your_plugin.proto
   ```

3. **检查语法和风格**
   ```bash
   make plugin-proto-lint
   ```

4. **生成代码**
   ```bash
   make plugin-proto
   ```

5. **实现接口**
   - 创建 `your_plugin/interface.go` (Go 接口定义)
   - 创建 `your_plugin/grpc.go` (gRPC 客户端/服务端)
   - 在 `plugins/your_plugin/` 中实现具体逻辑

6. **测试**
   ```bash
   make plugin-test
   make plugin-example
   ```

## 📚 参考资料

- [Buf 官方文档](https://buf.build/docs)
- [Protocol Buffers 语言指南](https://protobuf.dev/programming-guides/proto3/)
- [gRPC Go 快速开始](https://grpc.io/docs/languages/go/quickstart/)

---

**提示**: 生成的 `.pb.go` 和 `_grpc.pb.go` 文件不应该手动编辑,它们会在运行 `buf generate` 时自动重新生成。

