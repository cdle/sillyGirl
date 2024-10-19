# 快速开始

本文档将引导你完成 SillyGirl 的下载、安装、配置和第一个插件开发。

## 目录

- [系统要求](#系统要求)
- [安装方式](#安装方式)
  - [二进制安装（推荐）](#二进制安装推荐)
  - [源码编译](#源码编译)
  - [Docker 部署](#docker-部署)
- [首次运行](#首次运行)
- [配置说明](#配置说明)
- [终端模式详解](#终端模式详解)
- [常见问题](#常见问题)

## 系统要求

| 项目 | 要求 |
|------|------|
| 操作系统 | Linux / macOS / Windows |
| Go 版本 | 1.18+（仅源码编译需要）|
| 内存 | ≥ 64MB |
| 磁盘 | ≥ 50MB |
| 网络 | 可访问互联网（插件市场、自动升级）|

## 安装方式

### 二进制安装（推荐）

从 [GitHub Releases](../../releases) 下载对应平台的预编译二进制：

```bash
# 下载（以 Linux amd64 为例）
wget https://github.com/cdle/sillyGirl/releases/download/v2.0.0/sillyGirl_linux_amd64
mv sillyGirl_linux_amd64 sillyGirl
chmod +x sillyGirl

# 运行（终端模式）
./sillyGirl -t
```

其他平台文件名对照：

| 平台 | 文件名 |
|------|--------|
| Linux AMD64 | `sillyGirl_linux_amd64` |
| Linux ARM64 | `sillyGirl_linux_arm64` |
| macOS AMD64 | `sillyGirl_darwin_amd64` |
| macOS ARM64 | `sillyGirl_darwin_arm64` |
| Windows AMD64 | `sillyGirl_windows_amd64.exe` |

### 源码编译

```bash
# 克隆仓库
git clone https://github.com/cdle/sillyGirl.git
cd sillyGirl

# 编译
go build -ldflags "-s -w" -o sillyGirl

# 运行
./sillyGirl -t
```

编译参数说明：
- `-ldflags "-s -w"`：去除符号表和调试信息，减小体积
- 如需交叉编译，设置 `GOOS` 和 `GOARCH` 环境变量

### Docker 部署

```bash
# 构建镜像
docker build -t sillygirl:latest .

# 运行容器（挂载数据目录持久化存储）
docker run -d \
  --name sillygirl \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  --restart unless-stopped \
  sillygirl:latest

# 查看日志
docker logs -f sillygirl
```

## 首次运行

启动程序后，你将看到类似如下的日志输出：

```
2023/06/01 08:26:40 [I] 默认使用 boltdb 进行数据存储。
2023/06/01 08:26:40 [I] Http 服务已运行(8080)。
2023/06/01 08:26:40 [I] 管理员面板:
2023/06/01 08:26:40 [I]   > 本机: http://localhost:8080/admin
2023/06/01 08:26:40 [I]   > 局域网: http://192.168.1.100:8080/admin
```

### 访问管理面板

打开浏览器访问 `http://localhost:8080/admin`，首次访问会自动生成管理员 Token。

### 终端模式交互

使用 `-t` 参数启动时，程序会创建一个虚拟的终端机器人。你可以直接在命令行输入消息与之交互：

```
> 你好
Hello World!
> 猜拳
你先出，请在10秒内出拳！
> 石头
我出布，我赢了！
```

终端模式下所有插件规则均可正常触发，是开发和调试插件的最佳环境。

## 配置说明

SillyGirl 的配置通过 **Bucket 存储系统** 管理，而非传统的配置文件。你可以通过以下方式修改配置：

### 1. 管理面板

在 Admin 面板的"配置"页面中直接修改键值对。

### 2. 插件代码

```js
const app = Bucket("app");
app.set("port", 9090);  // 修改 HTTP 端口
```

### 3. 环境变量

部分配置支持通过环境变量注入：

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `SILLYGIRL_PORT` | HTTP 服务端口 | 8080 |
| `SILLYGIRL_REDIS_ADDR` | Redis 地址 | "" |
| `SILLYGIRL_REDIS_PASSWORD` | Redis 密码 | "" |

### 常用配置项

| Bucket | Key | 说明 | 默认值 |
|--------|-----|------|--------|
| `sillyGirl` | `port` | HTTP 服务端口 | 8080 |
| `sillyGirl` | `storage` | 存储后端 (`boltdb`/`redis`) | `boltdb` |
| `sillyGirl` | `redis_addr` | Redis 地址 | "" |
| `sillyGirl` | `redis_password` | Redis 密码 | "" |
| `sillyGirl` | `debug` | 调试模式 | `false` |
| `sillyGirl` | `api_key` | Admin API 密钥 | 自动生成 |
| `sillyGirl` | `listen_admin` | 管理员指令监听 | `true` |
| `sillyGirl` | `recall` | 自动撤回正则（&分隔）| "" |

### 存储后端切换

**BoltDB（默认）**：
- 嵌入式，无需额外依赖
- 数据存储在 `./sillyGirl.db`
- 适合单机部署

**Redis**：
```bash
# 方式一：管理面板修改
# storage = redis
# redis_addr = 127.0.0.1:6379
# redis_password = your_password

# 方式二：插件代码
Bucket("sillyGirl").set("storage", "redis");
Bucket("sillyGirl").set("redis_addr", "127.0.0.1:6379");
```

修改后重启生效。

## 终端模式详解

`./sillyGirl -t` 启动的终端模式是一个特殊的 QQ/Web 适配器模拟器，具有以下特性：

- **全权限用户**：终端输入者默认拥有管理员权限
- **即时反馈**：无需网络往返，插件响应延迟极低
- **日志同步**：插件的 `console.log` 会同时输出到终端和日志系统
- **插件热重载**：在 Admin 面板修改插件后，终端模式立即生效

## 常见问题

**Q: 启动后无法访问 Admin 面板？**
A: 检查防火墙是否放行对应端口。如果绑定的是 `localhost`，请使用 `127.0.0.1` 访问；如需外网访问，确保 `port` 配置正确且防火墙放行。

**Q: 如何重置管理员密码 / API Key？**
A: 删除 Bucket `sillyGirl` 中的 `api_key` 项，重启后自动生成新的 Key。

**Q: 插件安装后没有反应？**
A: 检查插件注释中的 `@rule` 是否匹配你的输入；检查插件是否被禁用（`@disable`）；查看日志确认是否有报错。

**Q: 如何备份数据？**
A: 如果使用 BoltDB，直接备份 `sillyGirl.db` 文件即可；如果使用 Redis，使用 `redis-cli SAVE` 或 Redis 持久化机制备份。

**Q: 如何关闭自动升级？**
A: 目前自动升级通过检测 `compiled_at` 变量触发。如需关闭，可在编译时设置固定版本号，或断网运行。
