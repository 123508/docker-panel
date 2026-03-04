# 配置文件示例

此文件展示 `config.toml` 的所有可配置项及其说明。

## 完整配置示例

```toml
[user]
# 管理员用户名
admin_username = "admin"

# 管理员密码
# 注意: 建议首次运行后立即修改默认密码
admin_password = "admin123"

[server]
# 服务器绑定IP
# "0.0.0.0" - 监听所有网卡
# "127.0.0.1" - 仅本地访问
bind_ip = "0.0.0.0"

# 服务器绑定端口
bind_port = "8080"
```

## 配置项说明

### [user] 用户配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| admin_username | string | "admin" | 管理员账号 |
| admin_password | string | "admin123" | 管理员密码 |

### [server] 服务器配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| bind_ip | string | "0.0.0.0" | 服务器绑定IP地址 |
| bind_port | string | "8080" | 服务器监听端口 |

## 使用说明

1. **首次运行**: 程序会自动在当前目录生成 `config.toml`
2. **修改配置**: 编辑 `config.toml` 文件后重启程序生效
3. **安全建议**: 
   - 首次运行后请立即修改默认密码
   - 如需公网访问，请确保使用强密码
   - 考虑配置防火墙规则

## 配置示例

### 仅本地访问
```toml
[server]
bind_ip = "127.0.0.1"
bind_port = "8080"
```

### 自定义端口
```toml
[server]
bind_ip = "0.0.0.0"
bind_port = "3000"
```

### 生产环境推荐配置
```toml
[user]
admin_username = "admin"
admin_password = "your-strong-password-here"

[server]
bind_ip = "0.0.0.0"
bind_port = "8080"
```
