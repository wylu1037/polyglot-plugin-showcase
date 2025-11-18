# Converter Plugin 使用指南

## 概述

Converter 插件是一个数据格式转换工具,专门用于将从数据库查询得到的 JSON 格式数据转换为其他常用格式,包括 CSV、TXT 和 HTML。

## 功能特性

### 1. JSON 转 CSV
- 支持单个对象或对象数组
- 自动提取表头(对象的键名)
- 支持自定义分隔符
- 优雅处理缺失字段

**使用场景:**
- 导出数据库查询结果到 Excel
- 生成数据报表
- 数据迁移

### 2. JSON 转 TXT
提供两种格式:
- **key-value**: 人类可读的键值对格式,支持嵌套结构
- **json-pretty**: 格式化的 JSON,带缩进

**使用场景:**
- 日志文件生成
- 配置文件导出
- 数据审查

### 3. JSON 转 HTML
- 生成 HTML 表格
- 可选的现代化 CSS 样式
- 支持生成完整 HTML 页面或仅表格片段
- 自动 HTML 转义,确保安全性

**使用场景:**
- 生成数据报告
- 创建数据展示页面
- 邮件内容生成

## 快速开始

### 构建插件

```bash
cd plugins/converter
go build -o ../../host-server/bin/plugins/converter/converter_v1.0.0 .
```

### 运行示例

```bash
cd plugins/converter/example
go run main.go
```

## API 使用示例

### 1. 转换为 CSV

**基本用法:**
```go
converter := &impl.ConverterImpl{}
jsonData := `[
    {"id": "1", "name": "John", "age": "30"},
    {"id": "2", "name": "Jane", "age": "25"}
]`

result, err := converter.ConvertToCSV(jsonData, map[string]string{})
```

**输出:**
```csv
age,id,name
30,1,John
25,2,Jane
```

**自定义分隔符:**
```go
result, err := converter.ConvertToCSV(jsonData, map[string]string{
    "delimiter": ";",
})
```

### 2. 转换为 TXT

**Key-Value 格式:**
```go
jsonData := `{
    "user": {
        "name": "John",
        "age": "30"
    },
    "active": true
}`

result, err := converter.ConvertToTXT(jsonData, map[string]string{
    "format": "key-value",
})
```

**输出:**
```
active: true
user:
  age: 30
  name: John
```

**Pretty JSON 格式:**
```go
result, err := converter.ConvertToTXT(jsonData, map[string]string{
    "format": "json-pretty",
})
```

### 3. 转换为 HTML

**基本表格:**
```go
jsonData := `[{"name": "John", "age": "30"}]`

result, err := converter.ConvertToHTML(jsonData, map[string]string{
    "styled": "true",
})
```

**完整 HTML 页面:**
```go
result, err := converter.ConvertToHTML(jsonData, map[string]string{
    "styled": "true",
    "full_page": "true",
})
```

## 参数说明

### ConvertToCSV 参数
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| data | string | 必填 | JSON 数据 |
| delimiter | string | "," | CSV 分隔符 |

### ConvertToTXT 参数
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| data | string | 必填 | JSON 数据 |
| format | string | "key-value" | 格式类型: "key-value" 或 "json-pretty" |

### ConvertToHTML 参数
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| data | string | 必填 | JSON 数据 |
| styled | string | "true" | 是否应用 CSS 样式 |
| full_page | string | "false" | 是否生成完整 HTML 页面 |

## 数据库集成示例

假设你从数据库查询得到 JSON 格式的数据:

```go
// 从数据库查询
rows, err := db.Query("SELECT id, name, email, age FROM users")
// ... 处理查询结果,转换为 JSON

jsonData := `[
    {"id": "1", "name": "John", "email": "john@example.com", "age": "30"},
    {"id": "2", "name": "Jane", "email": "jane@example.com", "age": "25"}
]`

// 使用 converter 插件转换
converter := &impl.ConverterImpl{}

// 导出为 CSV 供 Excel 使用
csv, _ := converter.ConvertToCSV(jsonData, map[string]string{})
os.WriteFile("users.csv", []byte(csv), 0644)

// 生成 HTML 报告
html, _ := converter.ConvertToHTML(jsonData, map[string]string{
    "styled": "true",
    "full_page": "true",
})
os.WriteFile("users_report.html", []byte(html), 0644)

// 生成文本日志
txt, _ := converter.ConvertToTXT(jsonData, map[string]string{
    "format": "key-value",
})
os.WriteFile("users.txt", []byte(txt), 0644)
```

## 错误处理

插件提供详细的错误信息:

```go
result, err := converter.ConvertToCSV("", map[string]string{})
if err != nil {
    // 错误: "JSON data cannot be empty"
    log.Fatal(err)
}

result, err = converter.ConvertToCSV("{invalid}", map[string]string{})
if err != nil {
    // 错误: "invalid JSON data: ..."
    log.Fatal(err)
}
```

## 测试

运行单元测试:

```bash
cd plugins/converter
go test ./impl/... -v
```

所有测试都应该通过,覆盖以下场景:
- 单个对象转换
- 对象数组转换
- 自定义选项
- 嵌套结构
- 错误情况

## 性能考虑

- 插件使用流式处理,内存占用低
- 对于大型数据集,建议分批处理
- HTML 生成包含内联样式,适合中小型数据集

## 插件信息

- **名称**: converter
- **版本**: 1.0.0
- **类型**: converter
- **输入格式**: JSON
- **输出格式**: CSV, TXT, HTML
- **协议版本**: 与 host-server 兼容

## 常见问题

**Q: 支持哪些 JSON 数据结构?**
A: 
- CSV/HTML: 支持对象或对象数组
- TXT: 支持任何有效的 JSON 结构

**Q: 如何处理嵌套的 JSON 对象?**
A: 
- CSV/HTML: 嵌套对象会被转换为字符串表示
- TXT (key-value): 支持嵌套结构,会自动缩进显示

**Q: 生成的 HTML 是否安全?**
A: 是的,所有 HTML 输出都经过转义处理,防止 XSS 攻击。

**Q: 可以自定义 HTML 样式吗?**
A: 当前版本使用内置样式。如需自定义,可以设置 `styled: "false"` 然后添加自己的 CSS。

## 下一步

- 查看 `example/main.go` 了解更多使用示例
- 阅读 `README.md` 了解架构细节
- 查看测试文件 `impl/converter_test.go` 了解更多用例

