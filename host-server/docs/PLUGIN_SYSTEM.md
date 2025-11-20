# 插件系统架构文档

> 基于 Terraform 插件设计理念的动态插件系统

## 概述

本插件系统采用了 Terraform 的动态化设计理念，实现了灵活、可扩展的插件架构。

### 核心特性

- ✅ **动态类型系统**：插件类型不再硬编码，支持任意自定义类型
- ✅ **命名空间隔离**：支持多来源插件，避免命名冲突
- ✅ **跨平台支持**：明确区分操作系统和架构
- ✅ **自动注册**：Registry 动态识别和注册新插件类型
- ✅ **握手验证**：启动插件进程验证其真实能力
- ✅ **结构化元数据**：强类型约束的插件元信息

---

## 目录结构

### 标准路径格式

```
bin/plugins/{namespace}/{type}/{name}/{version}/{os}_{arch}/plugin
```

#### 示例

```
bin/plugins/
├── official/                      # 官方命名空间
│   ├── data-processing/           # 数据处理类型
│   │   ├── anonymizer/            # 插件名称
│   │   │   └── v1.0.0/            # 版本
│   │   │       ├── darwin_arm64/  # macOS ARM64
│   │   │       │   └── plugin     # 二进制文件
│   │   │       └── linux_amd64/   # Linux AMD64
│   │   │           └── plugin
│   │   └── converter/
│   │       └── v1.2.0/
│   │           └── linux_amd64/
│   │               └── plugin
│   └── security/                  # 安全类型
│       └── encryptor/
│           └── v2.1.0/
│               └── darwin_arm64/
│                   └── plugin
└── third-party/                   # 第三方命名空间
    └── analytics/
        └── metrics/
            └── v0.5.0/
                └── linux_amd64/
                    └── plugin
```

### 路径组成部分

| 部分          | 说明                   | 示例                                           |
| ------------- | ---------------------- | ---------------------------------------------- |
| **namespace** | 命名空间，区分插件来源 | `official`, `third-party`, `mycompany`         |
| **type**      | 插件功能分类           | `data-processing`, `security`, `analytics`     |
| **name**      | 插件名称               | `anonymizer`, `encryptor`, `metrics`           |
| **version**   | 语义化版本             | `v1.0.0`, `v2.1.3-beta`                        |
| **os_arch**   | 操作系统\_架构         | `linux_amd64`, `darwin_arm64`, `windows_amd64` |
| **plugin**    | 可执行二进制文件       | `plugin` (必须可执行)                          |

---

## 插件类型系统

### 动态类型

不同于旧版的硬编码枚举，新系统支持任意自定义类型：

```go
// ❌ 旧版：硬编码类型
const (
    PluginTypeDesensitization = "desensitization"
    PluginTypeEncryption      = "encryption"
    // 添加新类型需要修改代码
)

// ✅ 新版：动态类型
type PluginType = string  // 任意字符串
```

### 推荐的类型分类

虽然类型不受限制，但推荐使用以下通用分类：

| 类型              | 说明     | 示例插件                            |
| ----------------- | -------- | ----------------------------------- |
| `data-processing` | 数据处理 | anonymizer, converter, transformer  |
| `security`        | 安全相关 | encryptor, validator, authenticator |
| `analytics`       | 分析统计 | metrics, logger, profiler           |
| `integration`     | 外部集成 | api-connector, database-sync        |
| `extension`       | 通用扩展 | custom-handler, webhook             |

---

## 插件元数据

### 结构定义

```go
type PluginMetadata struct {
    Name            string              `json:"name"`              // 插件名称
    Version         string              `json:"version"`           // 版本号
    Type            string              `json:"type"`              // 插件类型
    Description     string              `json:"description"`       // 描述
    Author          string              `json:"author"`            // 作者
    ProtocolVersion int32               `json:"protocol_version"`  // 协议版本
    Capabilities    []string            `json:"capabilities"`      // 能力列表
    Methods         []MethodMetadata    `json:"methods"`           // 方法定义
    Dependencies    []Dependency        `json:"dependencies"`      // 依赖项
}
```

### 示例元数据

```json
{
  "name": "anonymizer",
  "version": "v1.0.0",
  "type": "data-processing",
  "description": "数据匿名化插件，支持多种匿名化策略",
  "author": "Official Team",
  "protocol_version": 1,
  "capabilities": ["hash", "mask", "tokenize"],
  "methods": [
    {
      "name": "Anonymize",
      "description": "对数据进行匿名化处理",
      "parameters": {
        "data": {
          "type": "string",
          "description": "需要匿名化的数据",
          "required": true
        },
        "strategy": {
          "type": "string",
          "description": "匿名化策略",
          "required": false,
          "default": "hash"
        }
      },
      "returns": {
        "type": "string",
        "description": "匿名化后的数据"
      }
    }
  ],
  "dependencies": []
}
```

---

## API 使用指南

### 1. 安装插件

```bash
POST /api/plugins/install
Content-Type: application/json

{
  "downloadURL": "https://example.com/plugins/anonymizer_v1.0.0_linux_amd64",
  "namespace": "official",
  "name": "anonymizer",
  "version": "v1.0.0",
  "type": "data-processing",
  "os": "linux",
  "arch": "amd64",
  "description": "数据匿名化插件",
  "config": {
    "default_strategy": "hash"
  }
}
```

### 2. 列出插件（支持多维度过滤）

```bash
GET /api/plugins?namespace=official&type=data-processing&os=linux&arch=amd64&status=active
```

**查询参数：**

| 参数        | 说明     | 可选值                                                  |
| ----------- | -------- | ------------------------------------------------------- |
| `namespace` | 命名空间 | 任意字符串                                              |
| `type`      | 插件类型 | 任意字符串                                              |
| `os`        | 操作系统 | `linux`, `darwin`, `windows`                            |
| `arch`      | 架构     | `amd64`, `arm64`                                        |
| `status`    | 状态     | `active`, `inactive`, `disabled`, `error`, `installing` |

### 3. 调用插件

```bash
POST /api/plugins/123/call
Content-Type: application/json

{
  "method": "Anonymize",
  "params": {
    "data": "user@example.com",
    "strategy": "hash"
  }
}
```

---

## 开发指南

### 创建新插件

#### 1. 选择类型和命名空间

```bash
# 决定你的插件归属
NAMESPACE="mycompany"          # 命名空间
TYPE="data-processing"         # 类型
NAME="my-anonymizer"          # 名称
VERSION="v1.0.0"              # 版本
OS="linux"                    # 操作系统
ARCH="amd64"                  # 架构
```

#### 2. 创建目录结构

```bash
mkdir -p bin/plugins/${NAMESPACE}/${TYPE}/${NAME}/${VERSION}/${OS}_${ARCH}
```

#### 3. 实现插件接口

插件必须实现 gRPC 接口并返回正确的元数据：

```go
func (p *MyPlugin) GetMetadata(ctx context.Context, req *common.MetadataRequest) (*common.MetadataResponse, error) {
    return &common.MetadataResponse{
        Name:            "my-anonymizer",
        Version:         "v1.0.0",
        Description:     "Custom anonymization plugin",
        Methods:         []string{"Anonymize", "DeAnonymize"},
        Capabilities:    map[string]string{
            "hash": "SHA-256",
            "mask": "Partial masking",
        },
        ProtocolVersion: 1,
    }, nil
}
```

#### 4. 编译并放置

```bash
# 编译
go build -o plugin main.go

# 移动到正确位置
mv plugin bin/plugins/${NAMESPACE}/${TYPE}/${NAME}/${VERSION}/${OS}_${ARCH}/

# 设置可执行权限
chmod +x bin/plugins/${NAMESPACE}/${TYPE}/${NAME}/${VERSION}/${OS}_${ARCH}/plugin
```

### 跨平台编译

```bash
# macOS ARM64
GOOS=darwin GOARCH=arm64 go build -o plugin_darwin_arm64 main.go

# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o plugin_linux_amd64 main.go

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o plugin_windows_amd64.exe main.go
```

---

## 设计理念

### 学习自 Terraform

本系统借鉴了 Terraform 的以下设计理念：

1. **简化类型系统**

   - 只有一种技术实现（基于 go-plugin 的 gRPC 插件）
   - 通过元数据声明功能，而非代码硬编码

2. **进程隔离**

   - 每个插件独立进程
   - 崩溃不影响主系统

3. **协议优先**

   - 定义良好的 gRPC 协议
   - 支持多语言实现

4. **自描述能力**
   - 插件通过 `GetMetadata()` 声明自己的能力
   - 主系统通过元数据进行验证和调用

### 与旧版对比

| 特性     | 旧版                      | 新版                                                     |
| -------- | ------------------------- | -------------------------------------------------------- |
| 类型系统 | 硬编码枚举                | 动态字符串                                               |
| 路径结构 | `{type}/{name}_{version}` | `{namespace}/{type}/{name}/{version}/{os}_{arch}/plugin` |
| 跨平台   | 不支持                    | 原生支持                                                 |
| 命名空间 | 无                        | 多命名空间隔离                                           |
| 元数据   | 弱类型 JSON               | 结构化定义                                               |
| 握手验证 | 可选                      | 强制验证                                                 |

---

## 最佳实践

### 1. 命名规范

- **namespace**: 小写字母、数字、连字符，如 `my-company`
- **type**: 小写字母、连字符，如 `data-processing`
- **name**: 小写字母、数字、连字符，如 `my-plugin`
- **version**: 遵循语义化版本，如 `v1.2.3`

### 2. 版本管理

```
v1.0.0  -> 初始稳定版本
v1.0.1  -> bug修复
v1.1.0  -> 新增功能（向后兼容）
v2.0.0  -> 破坏性变更
```

### 3. 元数据完整性

确保 `GetMetadata()` 返回完整信息：

- ✅ 所有支持的方法
- ✅ 每个方法的参数和返回值
- ✅ 插件的能力和限制
- ✅ 依赖的其他插件或服务

### 4. 错误处理

插件应该优雅处理错误并返回有意义的错误信息：

```go
if err != nil {
    return nil, status.Error(codes.InvalidArgument,
        fmt.Sprintf("invalid data format: %v", err))
}
```

---

## 故障排查

### 插件未被发现

```bash
# 检查路径结构
ls -la bin/plugins/*/

# 确认文件可执行
ls -l bin/plugins/official/data-processing/anonymizer/v1.0.0/linux_amd64/plugin

# 查看警告日志
Warning: failed to parse plugin path ...
```

### 握手失败

```bash
# 测试插件能否启动
./bin/plugins/official/data-processing/anonymizer/v1.0.0/linux_amd64/plugin

# 检查协议版本
# 确保插件 ProtocolVersion 与主系统兼容
```

### 类型不匹配

动态类型系统会自动注册未知类型，不会因类型问题拒绝加载。

---

## 参考资料

- [Terraform Plugin Development](https://www.terraform.io/plugin)
- [HashiCorp go-plugin](https://github.com/hashicorp/go-plugin)
- [语义化版本规范](https://semver.org/)

---

## 更新日志

### v2.0.0 (2025-01-20)

- ✅ 插件类型动态化
- ✅ 新增命名空间支持
- ✅ 新增跨平台 OS/Arch 字段
- ✅ 结构化元数据定义
- ✅ Registry 自动注册机制
- ✅ 强制握手验证

### v1.0.0

- 初始版本（旧路径格式）
