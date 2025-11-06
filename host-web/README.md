# Host Web - 插件管理前端

基于 React + React Router v7 + Tailwind CSS + shadcn/ui 的插件管理系统前端界面。

## 技术栈

- **React 18** - UI 框架
- **React Router v7** - 路由管理
- **Tailwind CSS 4** - 样式框架 (Vite 插件)
- **shadcn/ui** - UI 组件库
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **Kubb** - 从 Swagger/OpenAPI 自动生成 API 客户端
- **React Query** - 数据获取和状态管理
- **Axios** - HTTP 客户端

## 快速开始

### 1. 安装依赖

```bash
pnpm install
```

### 2. 生成 API 客户端代码

确保后端服务器正在运行 (http://localhost:8080),然后生成 API 客户端:

```bash
pnpm generate:api
```

这将从后端的 Swagger 文档自动生成:
- TypeScript 类型定义
- Axios 客户端函数
- React Query hooks

### 3. 启动开发服务器

```bash
pnpm dev
```

访问 http://localhost:5173

## 项目结构

```
host-web/
├── src/
│   ├── api/
│   │   └── generated/        # Kubb 自动生成的代码
│   │       ├── types/        # TypeScript 类型
│   │       ├── clients/      # Axios 客户端
│   │       └── hooks/        # React Query hooks
│   ├── components/
│   │   └── ui/               # shadcn/ui 组件
│   ├── routes/               # 路由页面
│   │   ├── PluginList.tsx    # 插件列表
│   │   ├── PluginDetail.tsx  # 插件详情
│   │   └── PluginDiscover.tsx # 发现新插件
│   ├── lib/                  # 工具库
│   │   ├── utils.ts          # 工具函数
│   │   └── query-client.ts   # React Query 配置
│   ├── App.tsx               # 根组件
│   ├── router.tsx            # 路由配置
│   ├── main.tsx              # 入口文件
│   └── index.css             # 全局样式
├── kubb.config.ts            # Kubb 配置
├── components.json           # shadcn/ui 配置
├── vite.config.ts            # Vite 配置
└── package.json
```

## 可用脚本

- `pnpm dev` - 启动开发服务器
- `pnpm build` - 构建生产版本
- `pnpm preview` - 预览生产构建
- `pnpm generate:api` - 从 Swagger 生成 API 客户端代码

## 路由

- `/` - 插件列表页面
- `/plugins/:id` - 插件详情页面
- `/plugins/discover` - 发现新插件页面

## API 集成

前端通过 Kubb 自动生成的代码与后端 API 交互:

1. **自动生成**: 运行 `pnpm generate:api` 从后端 Swagger 文档生成客户端代码
2. **类型安全**: 所有 API 调用都有完整的 TypeScript 类型支持
3. **React Query**: 使用生成的 hooks 进行数据获取和缓存

示例:

```typescript
import { useGetPluginStoresQuery } from '@/api/generated/hooks'

function PluginList() {
  const { data, isLoading, error } = useGetPluginStoresQuery()
  
  // 使用数据...
}
```

## 环境变量

创建 `.env.local` 文件配置环境变量:

```
VITE_API_BASE_URL=http://localhost:8080
```

## 开发工作流

1. 后端更新 API 后,运行 `pnpm generate:api` 重新生成客户端代码
2. 前端直接使用生成的 React Query hooks
3. 类型自动同步,减少前后端对接错误

## 添加新的 UI 组件

使用 shadcn/ui CLI 添加组件:

```bash
pnpm dlx shadcn@latest add [component-name]
```

例如:
```bash
pnpm dlx shadcn@latest add dropdown-menu
```

## 构建生产版本

```bash
pnpm build
```

构建产物将输出到 `dist/` 目录。

