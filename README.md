# Docker 可视化面板

Docker 容器管理的可视化 Web 面板，提供对容器、镜像、数据卷、网络和 Docker Compose 的完整管理能力。

## 项目结构

```
docker-panel/
├── internal/                      # Go 后端
│   ├── config/                    # 配置管理 (Viper + TOML)
│   ├── db/                        # SQLite 数据库初始化 (GORM)
│   ├── docker/                    # Docker 客户端接口 + 实现
│   ├── handler/                   # Gin HTTP 处理层
│   │   ├── register.go            # 路由注册 + 依赖注入
│   │   ├── middleware.go          # CORS / 日志 / 恢复 / JWT 认证中间件
│   │   ├── container_handler.go
│   │   ├── image_handler.go
│   │   ├── volume_handler.go
│   │   ├── network_handler.go
│   │   ├── compose_handler.go
│   │   └── user_handler.go
│   ├── models/                    # GORM 数据模型
│   ├── service/                   # 业务逻辑层（不含 Gin 引用）
│   │   ├── container_service.go
│   │   ├── image_service.go
│   │   ├── volume_service.go
│   │   ├── network_service.go
│   │   ├── compose_service.go
│   │   ├── user_service.go
│   │   ├── recent_containers.go
│   │   └── response.go            # 统一响应结构
│   └── utils/jwt.go               # JWT 工具
├── docker_cli_wrapper/            # Docker 客户端 mock (gomock)
│   └── mock/client_mock.go
├── web/
│   ├── frontend/                  # Vue 3 + TypeScript 前端
│   │   ├── src/
│   │   │   ├── views/             # 页面组件
│   │   │   ├── composables/       # 视图状态管理
│   │   │   ├── components/        # 通用 UI 组件
│   │   │   ├── services/          # API 调用层
│   │   │   │   └── modules/       # 按资源划分的 API 模块
│   │   │   ├── router/            # Vue Router 配置
│   │   │   ├── store/             # Pinia 状态 (认证/UI)
│   │   │   └── styles/            # CSS 样式
│   │   ├── index.html
│   │   ├── package.json
│   │   └── vite.config.js
│   └── dist/                      # 前端构建产物 (go:embed)
├── .github/workflows/             # CI/CD 流水线
├── main.go                        # 程序入口 + DI 组装
├── config.toml                    # 配置文件（首次运行自动生成）
└── docker-panel.db                # SQLite 数据库（首次运行自动创建）
```

## 技术栈

- **前端**: Vue 3 + TypeScript + Vite + Element Plus + Vue Router + Pinia
- **后端**: Go 1.25+ + Gin + GORM + SQLite
- **Docker**: 通过 Docker SDK (docker/docker/client) 直接与 Docker 守护进程通信
- **认证**: JWT（登录后 Bearer Token）
- **配置**: Viper (TOML 格式)
- **部署**: 单个可执行文件（`//go:embed` 内嵌前端静态资源）

## 功能特性

### 功能特性

**容器管理**
- 容器列表查看，状态、端口、挂载点信息展示
- 容器创建
- 容器详情查看
- 容器生命周期管理
- 容器删除
- 实时日志查看

**镜像管理**
- 镜像列表与详情
- 拉取 / 删除 镜像
- 运行镜像

**数据卷管理**
- 卷列表、创建、详情查看、删除

**网络管理**
- 网络列表、创建、详情查看、删除
- 容器接入 / 断开网络

**Docker Compose**
- Compose 项目列表与上传
- 项目状态查看、日志、启动 / 停止
- 删除项目

**用户认证**
- JWT 登录认证
- 首次运行自动创建管理员账号
- 前端路由守卫（未登录自动跳转登录页）

**仪表盘**
- 系统概览
- 活跃容器

## 快速开始

### 方式一：二进制文件运行

从 [GitHub Releases](https://github.com/123508/docker-panel/releases) 页面下载对应平台的压缩包，解压后直接运行。

**Linux x86_64：**

```bash
curl -LO https://github.com/123508/docker-panel/releases/latest/download/docker-panel-linux-amd64.tar.gz
tar -xzf docker-panel-linux-amd64.tar.gz
./docker-panel
```

**Linux ARM64：**

```bash
curl -LO https://github.com/123508/docker-panel/releases/latest/download/docker-panel-linux-arm64.tar.gz
tar -xzf docker-panel-linux-arm64.tar.gz
./docker-panel
```

**macOS（Intel）：**

```bash
curl -LO https://github.com/123508/docker-panel/releases/latest/download/docker-panel-darwin-amd64.tar.gz
tar -xzf docker-panel-darwin-amd64.tar.gz
./docker-panel
```

**macOS（Apple Silicon）：**

```bash
curl -LO https://github.com/123508/docker-panel/releases/latest/download/docker-panel-darwin-arm64.tar.gz
tar -xzf docker-panel-darwin-arm64.tar.gz
./docker-panel
```

**Windows：**

下载 `docker-panel-windows-amd64.zip`，解压后运行 `docker-panel-windows-amd64.exe`。

### 方式二：Docker 运行

**Linux / macOS：**

```bash
docker run -d \
  --name docker-panel \
  -p 8080:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v docker-panel-data:/app \
  docker-panel:1.0.0
```

**Windows PowerShell：**

```powershell
docker run -d `
  --name docker-panel `
  -p 8080:8080 `
  -v /var/run/docker.sock:/var/run/docker.sock `
  -v docker-panel-data:/app `
  docker-panel:1.0.0
```

**Windows CMD：**

```cmd
docker run -d ^
  --name docker-panel ^
  -p 8080:8080 ^
  -v /var/run/docker.sock:/var/run/docker.sock ^
  -v docker-panel-data:/app ^
  docker-panel:1.0.0
```

### 参数说明

- `/var/run/docker.sock` — 挂载宿主机 Docker Unix socket，使容器内可访问宿主机 Docker 守护进程
- `docker-panel-data:/app` — 命名卷持久化 SQLite 数据库，容器删除后数据不丢失
- `-p 8080:8080` — 端口映射，左侧宿主机端口可按需修改（如 `-p 8081:8080`）

### 访问应用

首次运行时，程序会自动生成 `config.toml` 配置文件和 `docker-panel.db` 数据库，并创建管理员账号。

```
http://localhost:8080
```

默认管理员账号：
- 用户名: `admin`
- 密码：首次运行自动随机生成，保存在 `config.toml` 的 `admin_password` 字段中

## 配置文件

程序首次运行时会在当前目录自动生成 `config.toml`（密码随机生成）：

```toml
[user]
admin_username = "admin"
admin_password = "your_password"    # 首次运行自动随机生成

[server]
bind_ip = "0.0.0.0"
bind_port = "8080"
debug = false

[jwt]
secret = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"    # 首次运行自动随机生成
```

### 配置说明

- `user.admin_username` — 管理员账号
- `user.admin_password` — 管理员密码（首次运行自动随机生成；如需重置，直接修改此字段后重启即可）
- `server.bind_ip` — 服务器绑定 IP
- `server.bind_port` — 服务器端口
- `server.debug` — 调试模式开关
- `jwt.secret` — JWT 签名密钥（首次运行自动随机生成 64 位十六进制字符串）

## 开发指南

### 1. 安装依赖

```bash
# 安装前端依赖
cd web/frontend
npm install
cd ../..
```

### 2. 开发模式

**方式一：前后端分离开发**

终端1 - 启动前端开发服务器（端口 5173，自动代理 `/api` 到后端）：
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

调试模式（查看路由表、Docker 版本等信息）：
```bash
go run main.go -mode debug
```

### 3. 构建发布版本

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

构建后的单文件已包含所有前端资源，可直接分发部署。

### 4. Docker 镜像构建

```bash
docker build -t docker-panel:1.0.0 .
```

## API 概览

所有 `/api/v1/*` 路由需要 `Authorization: Bearer <token>` 头。

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/health` | GET | 健康检查 |
| `/api/login` | POST | 登录获取 JWT |
| `/api/v1/containers` | GET | 容器列表 |
| `/api/v1/containers` | POST | 创建容器 |
| `/api/v1/containers/recent` | GET | 最近操作容器（LRU） |
| `/api/v1/containers/:id` | GET/DELETE | 详情 / 删除 |
| `/api/v1/containers/:id/start` | POST | 启动 |
| `/api/v1/containers/:id/stop` | POST | 停止 |
| `/api/v1/containers/:id/restart` | POST | 重启 |
| `/api/v1/containers/:id/kill` | POST | 强制停止 |
| `/api/v1/containers/:id/pause` | POST | 暂停 |
| `/api/v1/containers/:id/unpause` | POST | 恢复 |
| `/api/v1/containers/:id/rename` | POST | 重命名 |
| `/api/v1/containers/:id/logs` | GET | 日志 |
| `/api/v1/containers/:id/logs/ws` | GET | 日志 WebSocket |
| `/api/v1/containers/:id/exec` | POST | 执行命令 |
| `/api/v1/containers/:id/terminal` | GET | Web 终端 |
| `/api/v1/containers/:id/top` | GET | 进程列表 |
| `/api/v1/containers/:id/ports` | GET | 端口映射 |
| `/api/v1/containers/:id/mounts` | GET | 挂载信息 |
| `/api/v1/containers/:id/export` | GET | 导出容器 |
| `/api/v1/images` | GET | 镜像列表 |
| `/api/v1/images/:id` | GET/DELETE | 详情 / 删除 |
| `/api/v1/images/pull` | POST | 拉取镜像 |
| `/api/v1/images/tag` | POST | 标记镜像 |
| `/api/v1/images/push` | POST | 推送镜像 |
| `/api/v1/images/:id/save` | GET | 导出镜像文件 |
| `/api/v1/images/load` | POST | 加载镜像 |
| `/api/v1/images/import` | POST | 导入镜像 |
| `/api/v1/volumes` | GET/POST | 列表 / 创建 |
| `/api/v1/volumes/:name` | GET/DELETE | 详情 / 删除 |
| `/api/v1/networks` | GET/POST | 列表 / 创建 |
| `/api/v1/networks/:id` | GET/DELETE | 详情 / 删除 |
| `/api/v1/networks/:id/connect` | POST | 接入容器 |
| `/api/v1/networks/:id/disconnect` | POST | 断开容器 |
| `/api/v1/compose/projects` | GET/POST | 列表 / 上传 |
| `/api/v1/compose/projects/:name/ps` | GET | 项目状态 |
| `/api/v1/compose/projects/:name/logs` | GET | 项目日志 |
| `/api/v1/compose/projects/:name/up` | POST | 启动项目 |
| `/api/v1/compose/projects/:name/stop` | POST | 停止项目 |
| `/api/v1/compose/projects/:name/restart` | POST | 重启项目 |
| `/api/v1/compose/projects/:name/scale` | POST | 弹性伸缩 |
| `/api/v1/compose/projects/:name` | DELETE | 删除项目 |

## 架构

后端采用分层设计：

```
main.go (DI 依赖注入)
  └── handler.Dependencies → handler.NewRouter
       ├── ContainerHandler → ContainerService → docker.DockerClientInterface
       ├── ImageHandler → ImageService           (同层级)
       ├── VolumeHandler → VolumeService          (同层级)
       ├── NetworkHandler → NetworkService        (同层级)
       ├── ComposeHandler → ComposeService        (同层级)
       ├── UserHandler → UserService (SQLite via GORM, internal/db)
       └── RecentContainers (跨请求共享的进程级 LRU)
```

- **Handler 层** (`internal/handler/`): 仅做 HTTP ↔ Service 的翻译，无业务逻辑。所有响应通过 `respondJSON` 统一返回 `{ code, msg, data }` 格式。
- **Service 层** (`internal/service/`): 纯业务逻辑，不引用 Gin。通过 `docker.DockerClientInterface` 与 Docker 交互，该接口由 `docker_cli_wrapper/mock` 提供 gomock 测试替身。
- **Docker 客户端** (`internal/docker/client.go`): 定义 `DockerClientInterface` 接口，连接 Docker 守护进程，所有 service 通过该接口间接调用 Docker API。

## 测试

```bash
go test -v ./...                # 全部测试
go test -v -run ^TestName$ ./internal/handler  # 单个 handler 测试
```

- Service 测试通过 `gomock.Controller` + `docker_cli_wrapper/mock` 模拟 Docker 客户端
- Handler 测试构建完整 Gin 路由，验证 HTTP 状态码和响应结构
