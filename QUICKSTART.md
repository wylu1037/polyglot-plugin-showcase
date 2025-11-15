# 快速入门指南

## 前置条件

- Go 1.21+
- PostgreSQL 数据库
- Make (可选)

## 步骤 1: 配置数据库

1. 创建数据库：

```bash
createdb polyglot_plugin
```

2. 更新配置文件 `host-server/config.yaml`:

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  database: polyglot_plugin
  ssl_mode: disable
```

## 步骤 2: 构建插件

```bash
# 构建所有插件
make plugin-build

# 或者单独构建 desensitization 插件
make plugin-desensitization
```

这会将插件二进制文件构建到 `host-server/bin/plugins/` 目录。

## 步骤 3: 启动服务器

```bash
# 方式 1: 使用 Make
make server-dev

# 方式 2: 直接运行
cd host-server
go run cmd/server/main.go
```

服务器将在 `http://localhost:8080` 启动。

## 步骤 4: 测试插件系统

### 4.1 安装插件（本地文件）

由于我们已经构建了插件，可以手动将其添加到数据库：

```bash
# 使用 psql 或任何数据库客户端
psql -d polyglot_plugin -c "
INSERT INTO plugins (name, version, type, description, status, binary_path, protocol, protocol_version, checksum_type)
VALUES ('desensitization', 'v1.0.0', 'desensitization', 'Data desensitization plugin', 'inactive', './bin/plugins/desensitization', 'grpc', 1, 'sha256')
RETURNING id;
"
```

或者通过 API 安装（如果插件托管在 Maven 上）：

```bash
curl -X POST http://localhost:8080/api/plugins/install \
  -H "Content-Type: application/json" \
  -d '{
    "download_url": "https://your-maven-repo.com/plugins/desensitization-v1.0.0",
    "name": "desensitization",
    "version": "v1.0.0",
    "type": "desensitization",
    "description": "Data desensitization plugin"
  }'
```

### 4.2 列出所有插件

```bash
curl http://localhost:8080/api/plugins
```

### 4.3 激活插件

```bash
# 假设插件 ID 为 1
curl -X POST http://localhost:8080/api/plugins/1/activate
```

### 4.4 调用插件方法

```bash
# 脱敏姓名
curl -X POST http://localhost:8080/api/plugins/1/call \
  -H "Content-Type: application/json" \
  -d '{
    "method": "DesensitizeName",
    "params": {
      "data": "张三"
    }
  }'

# 预期输出: {"success":true,"data":"张**"}

# 脱敏手机号
curl -X POST http://localhost:8080/api/plugins/1/call \
  -H "Content-Type: application/json" \
  -d '{
    "method": "DesensitizeTelNo",
    "params": {
      "data": "13812345678"
    }
  }'

# 预期输出: {"success":true,"data":"138****5678"}

# 脱敏邮箱
curl -X POST http://localhost:8080/api/plugins/1/call \
  -H "Content-Type: application/json" \
  -d '{
    "method": "DesensitizeEmail",
    "params": {
      "data": "user@example.com"
    }
  }'

# 预期输出: {"success":true,"data":"u***@example.com"}
```

### 4.5 停用插件

```bash
curl -X POST http://localhost:8080/api/plugins/1/deactivate
```

### 4.6 卸载插件

```bash
curl -X DELETE http://localhost:8080/api/plugins/1
```

## 步骤 5: 开发自己的插件

参考 `PLUGIN_SYSTEM.md` 中的插件开发指南。

## 常见问题

### Q: 插件激活失败

**A**: 检查以下几点：
1. 插件二进制文件是否存在且可执行
2. 插件路径是否正确（相对于服务器工作目录）
3. 查看服务器日志中的详细错误信息

### Q: 数据库连接失败

**A**: 确保：
1. PostgreSQL 服务正在运行
2. 数据库已创建
3. config.yaml 中的数据库配置正确

### Q: 插件调用返回错误

**A**: 确认：
1. 插件状态为 `active`
2. 方法名拼写正确（区分大小写）
3. 参数格式正确

## 下一步

- 阅读 `PLUGIN_SYSTEM.md` 了解完整的系统文档
- 查看 `plugins/desensitization` 了解插件实现示例
- 开发自己的插件类型

## 项目结构

```
polyglot-plugin-showcase/
├── host-server/              # 主服务器
│   ├── cmd/server/          # 入口点
│   ├── app/
│   │   ├── database/        # 数据库层
│   │   └── modules/plugins/ # 插件模块
│   ├── internal/
│   │   ├── bootstrap/       # 应用启动
│   │   └── plugin/          # 插件管理核心
│   ├── config/              # 配置
│   └── bin/plugins/         # 插件二进制存储
├── plugins/                 # 插件实现
│   └── desensitization/     # 脱敏插件示例
├── proto/                   # Protocol Buffers 定义
│   └── desensitization/     # 脱敏插件接口
└── Makefile                 # 构建脚本
```

## 有用的命令

```bash
# 构建服务器
make server-build

# 运行服务器（开发模式）
make server-dev

# 构建所有插件
make plugin-build

# 清理插件二进制
make plugin-clean

# 生成 proto 代码
make plugin-proto

# 运行插件测试
make plugin-test
```

