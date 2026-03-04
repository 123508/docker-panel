# Docker 可视化面板

Docker 容器管理的可视化 Web 面板基础框架

## 项目结构

```
docker-panel/
├── internal/            # 内部包
│   ├── config/         # 配置管理 (Viper)
│   ├── api/
│   ├── handler/
│   ├── service/
│   └── docker/
├── web/
│   ├── frontend/        # Vue 前端源码
│   │   ├── src/
│   │   │   ├── App.vue
│   │   │   └── main.js
│   │   ├── index.html
│   │   ├── package.json
│   │   └── vite.config.js
│   └── dist/            # 前端编译产物（已嵌入到Go二进制）
├── main.go              # 主程序入口
├── config.toml          # 配置文件（首次运行自动生成）
└── go.mod
```

## 技术栈

- **前端**: Vue 3 + Vite
- **后端**: Go 1.23+
- **配置**: Viper (TOML 格式)
- **部署**: 单个可执行文件（embed 静态资源）

## 快速开始

### 1. 安装依赖

```bash
# 安装前端依赖
cd web/frontend
npm install
cd ../..
```

### 2. 开发模式

**方式一：前后端分离开发**

终端1 - 启动前端开发服务器：
```bash
cd web/frontend
npm run dev
```

终端2 - 启动后端：
```bash
go run main.go
```

**方式二：构建后整体运行**

```bash
# 构建前端
cd web/frontend
npm run build
cd ../..

# 运行后端（会自动加载编译后的前端）
go run main.go
```

### 3. 访问应用

首次运行时，程序会自动生成 `config.toml` 配置文件。

打开浏览器访问：
```
http://localhost:8080
```

默认管理员账号：
- 用户名: `admin`
- 密码: `admin123`

## 配置文件

程序首次运行时会在当前目录自动生成 `config.toml`：

```toml
[user]
admin_username = "admin"
admin_password = "admin123"

[server]
bind_ip = "0.0.0.0"
bind_port = "8080"
```

### 配置说明

- `user.admin_username` - 管理员账号
- `user.admin_password` - 管理员密码
- `server.bind_ip` - 服务器绑定IP（`0.0.0.0` 表示监听所有网卡）
- `server.bind_port` - 服务器端口

## 构建发布版本

```bash
# 1. 构建前端
cd web/frontend
npm run build
cd ../..

# 2. 构建 Go 二进制文件
go build -o docker-panel

# 3. 运行
./docker-panel
```

构建后的可执行文件已包含所有前端资源，可以直接部署。

## API 端点

当前提供的模拟 API：

- `GET /api/health` - 健康检查
- `GET /api/containers` - 获取容器列表（模拟数据）
- `GET /api/images` - 获取镜像列表（模拟数据）

## 当前状态

✅ 已完成：
- 基础项目结构
- Go 后端框架（embed 静态文件 + SPA 路由支持）
- Vue 单页面前端框架
- API 路由（模拟数据）
- 前后端集成
- 配置文件管理（Viper + TOML）
- 首次运行自动生成配置

⏳ 待实现：
- Docker API 集成
- 实际容器/镜像管理功能
- WebSocket 实时更新
- 认证授权
- 更多管理功能

## 开发说明

### 前端开发

前端源码位于 `web/frontend/`：
- 使用 Vite 作为构建工具
- 构建输出自动到 `web/dist/`
- 修改 `vite.config.js` 可调整构建配置

### 后端开发

- `main.go` - 主程序，包含路由和 API handler
- `internal/config/` - 配置管理模块
- `internal/` - 其他内部包
- 使用 `//go:embed` 嵌入前端编译产物

### 配置开发

配置文件使用 Viper 库管理：
- 配置模块位于 `internal/config/config.go`
- 使用 TOML 格式
- 首次运行自动生成默认配置
- 通过 `config.AppConfig` 访问配置

### 添加新 API

在 `main.go` 中添加新的 handler：

```go
func newHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    // ... 实现逻辑
}

// 在 main() 中注册
http.HandleFunc("/api/new", newHandler)
```

## 端口配置

默认端口：`8080`

修改 `config.toml` 文件中的 `server.bind_port` 来更改端口：

```toml
[server]
bind_ip = "0.0.0.0"
bind_port = "8090"  # 修改为其他端口
```
