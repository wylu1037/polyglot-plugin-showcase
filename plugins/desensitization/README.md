# 数据脱敏插件 (Desensitization Plugin)

基于 [hashicorp/go-plugin](https://github.com/hashicorp/go-plugin) 实现的数据脱敏插件,支持多种常见数据类型的脱敏处理。

## 📋 功能特性

该插件提供以下数据脱敏功能:

| 功能 | 说明 | 示例 |
|------|------|------|
| **姓名脱敏** | 保留第一个字符,其余用 `*` 替换 | `张三` → `张*` |
| **手机号脱敏** | 保留前3位和后4位,中间4位用 `*` 替换 | `13812345678` → `138****5678` |
| **身份证号脱敏** | 保留前2位和后2位,其余用 `*` 替换 | `110101199001011234` → `11**************34` |
| **邮箱脱敏** | 保留用户名首字符和完整域名,其余用 `*` 替换 | `user@example.com` → `u***@example.com` |
| **银行卡号脱敏** | 保留前6位和后3位,其余用 `*` 替换 | `6222021234567890123` → `622202**********123` |
| **地址脱敏** | 保留前1/3内容,其余用 `*` 替换 | `北京市朝阳区某某街道123号` → `北京市朝阳区********` |

## 🏗️ 架构设计

```
┌─────────────────┐
│   Host Process  │
│  (Your App)     │
└────────┬────────┘
         │ gRPC over Unix Socket
         │
┌────────▼────────┐
│  Plugin Process │
│  (Desensitizer) │
└─────────────────┘
```

### 核心组件

1. **接口定义** (`proto/shared/interface.go`)
   - 定义了 `Desensitizer` 接口
   - 所有插件必须实现此接口

2. **gRPC 协议** (`proto/shared/desensitizer.proto`)
   - 定义了 gRPC 服务和消息格式
   - 支持跨语言插件开发

3. **插件实现** (`plugins/desensitization/impl/`)
   - 具体的脱敏算法实现
   - 参考 [WGrape/golib](https://github.com/WGrape/golib/blob/main/desensitization/desensitization.go)

4. **插件主程序** (`plugins/desensitization/main.go`)
   - 插件入口,启动 gRPC 服务器

## 🚀 快速开始

### 构建插件

```bash
# 方式1: 使用 Makefile
make plugin-desensitization

# 方式2: 手动构建
cd plugins/desensitization
go build -o ../../host-server/bin/plugins/desensitization
```

构建成功后,插件二进制文件位于: `host-server/bin/plugins/desensitization`

### 运行示例

```bash
cd plugins/desensitization/example
go run main.go
```

### 在你的应用中使用

```go
package main

import (
    "fmt"
    "log"
    "os/exec"

    "github.com/hashicorp/go-plugin"
    "github.com/wylu1037/polyglot-plugin-showcase/proto/desensitization"
)

func main() {
    // 1. 创建插件客户端
    client := plugin.NewClient(&plugin.ClientConfig{
        HandshakeConfig: desensitization.Handshake,
        Plugins:         desensitization.PluginMap,
        Cmd:             exec.Command("./host-server/bin/plugins/desensitization"),
        AllowedProtocols: []plugin.Protocol{
            plugin.ProtocolGRPC,
        },
    })
    defer client.Kill()

    // 2. 连接到插件
    rpcClient, err := client.Client()
    if err != nil {
        log.Fatal(err)
    }

    // 3. 获取插件实例
    raw, err := rpcClient.Dispense("desensitizer")
    if err != nil {
        log.Fatal(err)
    }

    // 4. 使用插件(就像调用本地接口一样!)
    desensitizer := raw.(desensitization.Desensitizer)
    
    result, err := desensitizer.DesensitizeName("张三")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(result) // 输出: 张*
}
```

## 📖 API 文档

### Desensitizer 接口

```go
type Desensitizer interface {
    // DesensitizeName 姓名脱敏
    DesensitizeName(name string) (string, error)

    // DesensitizeTelNo 手机号脱敏
    DesensitizeTelNo(telNo string) (string, error)

    // DesensitizeIDNumber 身份证号脱敏
    DesensitizeIDNumber(idNumber string) (string, error)

    // DesensitizeEmail 邮箱脱敏
    DesensitizeEmail(email string) (string, error)

    // DesensitizeBankCard 银行卡号脱敏
    DesensitizeBankCard(cardNumber string) (string, error)

    // DesensitizeAddress 地址脱敏
    DesensitizeAddress(address string) (string, error)
}
```

### 错误处理

所有方法在以下情况会返回错误:

- 输入为空字符串
- 输入格式不符合要求(如手机号不是11位)
- 其他验证失败

## 🔧 开发指南

### 添加新的脱敏方法

1. 在 `proto/shared/desensitizer.proto` 中添加 RPC 方法定义
2. 使用 buf 重新生成 protobuf 代码: `make plugin-proto`
3. 在 `proto/shared/interface.go` 中添加接口方法
4. 在 `proto/shared/grpc.go` 中实现 gRPC 客户端和服务端方法
5. 在 `plugins/desensitization/impl/desensitizer.go` 中实现具体逻辑
6. 编写测试用例
7. 重新构建插件: `make plugin-desensitization`

### 测试插件

```bash
# 运行示例程序
cd plugins/desensitization/example
go run main.go

# 或者编写单元测试
cd plugins/desensitization/impl
go test -v
```

### 调试插件

插件使用 `go-plugin` 的日志系统,默认会输出调试信息:

```bash
# 设置日志级别
export PLUGIN_LOG_LEVEL=DEBUG
go run example/main.go
```

## 🎯 设计考虑

### 为什么使用 go-plugin?

1. **进程隔离**: 插件崩溃不会影响主程序
2. **语言无关**: 理论上支持任何语言编写插件(只要实现 gRPC 协议)
3. **成熟稳定**: 被 Terraform、Vault 等 HashiCorp 产品广泛使用
4. **易于部署**: 插件就是独立的二进制文件

### 性能考虑

- 使用 gRPC 进行通信,性能优于传统 net/rpc
- Unix Socket 通信,本地网络开销很小
- 适合中等频率调用场景
- 如需极高性能,考虑使用共享库或内嵌实现

### 安全考虑

- 插件运行在独立进程中,有一定隔离性
- 建议验证插件二进制文件的签名
- 可以配置 TLS 加密通信(本示例未启用)
- 限制插件的文件系统和网络访问权限

## 📚 参考资料

- [hashicorp/go-plugin 官方文档](https://github.com/hashicorp/go-plugin)
- [gRPC Go 快速开始](https://grpc.io/docs/languages/go/quickstart/)
- [WGrape/golib 脱敏实现](https://github.com/WGrape/golib/blob/main/desensitization/desensitization.go)

## 📝 许可证

MIT License

---

**开发状态**: ✅ 已完成并测试通过

