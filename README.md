# SillyGirl

> 一个可扩展、跨平台的开源机器人框架，内置强大的 JavaScript 插件系统与丰富的交互能力。

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-stable-green.svg)]()

## 目录

- [简介](#简介)
- [核心特性](#核心特性)
- [架构概览](#架构概览)
- [快速开始](#快速开始)
  - [二进制安装](#二进制安装)
  - [源码编译](#源码编译)
  - [Docker 部署](#docker-部署)
  - [第一个插件](#开发第一个插件)
- [项目结构](#项目结构)
- [技术栈](#技术栈)
- [文档](#文档)
- [版本历史](#版本历史)
- [许可](#许可)

## 简介

SillyGirl 是一个基于 **Go 语言** 开发的高性能开源机器人框架，其设计核心围绕强大的 **JavaScript 插件系统** 展开。框架内置了完整的 ECMAScript 5.1 运行时（基于 [Goja](https://github.com/dop251/goja) 引擎），允许开发者使用熟悉的 JavaScript 语法编写插件，并通过热重载机制实现功能的动态扩展，无需重启服务。

框架提供了丰富的内置能力：持久化键值存储、Cron 定时任务调度、HTTP/WebSocket 服务、gRPC 跨语言 RPC、Web Admin 管理面板、多平台机器人适配器等。开发者可以通过简单的 JavaScript 脚本快速构建具有复杂交互逻辑的机器人应用，并同时接入多个平台（QQ、Web、Pagermaid 等）的多个机器人实例，实现统一的业务逻辑与跨平台消息互通。

## 核心特性

### JavaScript 插件系统
- **完整 ES5.1 支持**：基于 Goja 引擎，支持闭包、原型链、正则表达式等标准语法
- **热重载机制**：插件文件变更后自动重新加载，开发调试零停机
- **丰富的元数据注解**：通过注释声明规则匹配、定时任务、HTTP 路由、权限控制等
- **Node.js 兼容层**：内置 `request`、`crypto`、`os` 等常用 Node API 的模拟实现
- **插件市场**：支持订阅远程插件源，一键安装、更新、卸载

### 多平台适配器架构
- **统一抽象接口**：所有平台通过标准化的 `Sender` 和 `Factory` 接口接入核心引擎
- **多实例管理**：同一平台可同时接入多个机器人账号，支持负载均衡与故障转移
- **内置适配器**：QQ（CQHTTP/OQ）、Web（内置聊天页）、Pagermaid（Python 桥接）
- **自定义适配器**：通过 gRPC 或 Go 接口自行开发新平台适配器

### 交互与存储能力
- **Bucket 持久化存储**：键值对存储抽象，支持 BoltDB（默认）、Redis、MongoDB 后端
- **存储变更监听**：支持 `watch` 机制，配置变更实时通知插件，实现热配置更新
- **Cron 定时任务**：基于 `robfig/cron`，支持秒级和分钟级表达式，多平台独立调度
- **消息监听与等待**：`s.listen()` 支持按规则捕获后续消息，实现对话式交互
- **群聊管理**：内置禁言、踢人、群组监听/屏蔽等群管能力

### 网络服务
- **HTTP 服务**：基于 Gin 框架，插件可通过注释声明 HTTP 路由，或运行时动态注册
- **WebSocket**：内置实时通信通道，Admin 面板与 Web 聊天均基于此
- **gRPC 服务**：提供跨语言调用的 RPC 接口（Bucket、Plugin、Adapter、Sender 等）
- **Admin 管理面板**：基于 React 的可视化界面，支持插件管理、存储浏览、日志查看、配置修改

### 运维与扩展
- **自动升级**：内置版本检测与二进制热更新机制
- **日志系统**：完整的分级日志框架，支持文件、控制台、ES、Slack、SMTP 等多种后端
- **容器化部署**：提供 Dockerfile，支持 Docker 一键部署
- **代理支持**：内置 HTTP/SOCKS5 代理传输层，支持翻墙与内网穿透场景

## 架构概览

```
┌─────────────────────────────────────────────────────────────┐
│                        Adapters                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────┐  │
│  │    QQ    │  │   Web    │  │Pagermaid │  │  Custom    │  │
│  │ (CQHTTP) │  │ (ChatUI) │  │ (Python) │  │  (gRPC)    │  │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └─────┬──────┘  │
└───────┼─────────────┼─────────────┼──────────────┼─────────┘
        │             │             │              │
        └─────────────┴─────────────┴──────────────┘
                              │
                    ┌─────────▼──────────┐
                    │   Message Router   │
                    │  (Listen/Reply/    │
                    │   Group Filter)    │
                    └─────────┬──────────┘
                              │
┌─────────────────────────────▼───────────────────────────────┐
│                      Core Engine                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   Plugin    │  │   Bucket    │  │   Adapter Manager   │ │
│  │   Engine    │  │   Storage   │  │   (Factory/Pool)    │ │
│  │  (Goja VM)  │  │(BoltDB/Redis│  │                     │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │    Cron     │  │   Web/Gin   │  │   gRPC Services     │ │
│  │  Scheduler  │  │   Server    │  │  (srpc.proto)       │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

**数据流**：
1. 适配器接收平台原始消息，构造 `Sender` 对象
2. 消息进入路由层，进行群组过滤、用户屏蔽、管理员指令处理
3. 插件引擎按优先级遍历所有插件，正则匹配 `rule` 规则
4. 匹配的插件在隔离的 Goja 运行时中执行，通过 `Sender` 接口回复
5. 适配器将回复发回原始平台

## 快速开始

### 二进制安装

从 [Releases](../../releases) 下载对应系统的可执行文件：

```bash
# Linux / macOS
chmod +x sillyGirl
./sillyGirl -t

# Windows
sillyGirl.exe -t
```

`-t` 参数开启终端机器人模式，启动后可直接在命令行与程序交互：

```
2023/06/01 08:26:40 [I] 默认使用 boltdb 进行数据存储。
2023/06/01 08:26:40 [I] Http 服务已运行(8080)。
```

访问 `http://localhost:8080/admin` 打开 Admin 管理面板。

### 源码编译

```bash
git clone https://github.com/cdle/sillyGirl.git
cd sillyGirl
go build -o sillyGirl
```

### Docker 部署

```bash
docker build -t sillygirl .
docker run -d -p 8080:8080 -v $(pwd)/data:/data sillygirl
```

### 开发第一个插件

创建 `hello.js`：

```js
/**
 * @title HelloWorld
 * @rule raw ^你好$
 */

s.reply("Hello World!");
```

在终端输入 `你好`，即可看到回复 `Hello World!`。

**进阶示例 — 猜拳游戏**：

```js
/**
 * @title 猜拳游戏
 * @rule 猜拳
 */

s.reply("你先出，请在10秒内出拳！");
const result = s.listen({
  rules: ["[出拳:剪刀,石头,布]"],
  timeout: 10000,
  handle: (s) => {
    const choose = s.param("出拳");
    const win = { "石头": "布", "剪刀": "石头", "布": "剪刀" };
    s.reply(`我出${win[choose]}，我赢了！`);
  },
});
if (!result) {
  s.reply("你没出拳，算我赢了！");
}
```

更多开发文档见 [docs/](docs/)。

## 项目结构

```
sillyGirl/
├── adapters/              # 平台适配器
│   ├── qq/               # QQ 机器人适配器
│   ├── web/              # Web 聊天适配器
│   └── pagermaid/        # Pagermaid 桥接适配器
├── core/                  # 核心框架
│   ├── admin/            # React 管理面板（编译产物，embed）
│   ├── common/           # 公共接口定义（Sender、Function）
│   ├── logs/             # 分级日志框架
│   ├── storage/          # 存储抽象与后端实现
│   ├── adapter.go        # 适配器工厂与消息收发
│   ├── bucket.go         # Bucket 键值存储
│   ├── function.go       # 消息路由与规则匹配
│   ├── init.go           # 系统初始化流程
│   ├── plugin_core.go    # 插件引擎（加载/卸载/热重载）
│   ├── plugin_impl.go    # JS API 实现（Sender、Cron、Bucket）
│   ├── web.go            # Gin Web 服务器与 Admin 面板
│   └── grpc_*.go         # gRPC 服务实现
├── proto3/               # Protobuf 定义与多语言生成代码
├── mongodb/              # MongoDB 存储后端
├── emoji/                # Emoji 数据处理
├── docs/                 # 项目文档
├── main.go               # 程序入口
├── go.mod                # Go 模块依赖
└── .dockerfile           # 容器构建配置
```

## 技术栈

| 层次 | 技术 | 说明 |
|------|------|------|
| 语言 | Go 1.18+ | 核心框架开发语言 |
| JS 运行时 | [Goja](https://github.com/dop251/goja) | ECMAScript 5.1，纯 Go 实现 |
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) | HTTP 服务与 REST API |
| 前端 | React / Ant Design Pro | Admin 管理面板 |
| 存储 | BoltDB / Redis / MongoDB | 键值对持久化 |
| 定时任务 | [robfig/cron/v3](https://github.com/robfig/cron) | Cron 表达式调度 |
| RPC | [gRPC](https://grpc.io) | 跨语言服务接口 |
| 消息协议 | CQHTTP / 自定义 | QQ 等平台的通信协议 |

## 文档

| 文档 | 说明 |
|------|------|
| [docs/quickstart.md](docs/quickstart.md) | 详细安装与配置指南 |
| [docs/plugin-dev.md](docs/plugin-dev.md) | 插件开发完整指南与 API 详解 |
| [docs/architecture.md](docs/architecture.md) | 架构设计与核心模块分析 |
| [docs/api-reference.md](docs/api-reference.md) | REST、gRPC 与 JavaScript API 参考 |
| [docs/deployment.md](docs/deployment.md) | 二进制、Docker 与反向代理部署 |

## 版本历史

本项目经历了两个主要阶段：

- **v1 (2021)** — 早期探索版本，基于直接函数调用的简单机器人框架
- **v2 (2023)** — 全面重构，引入 Goja JS 插件系统、Bucket 存储抽象、gRPC 服务、Admin 面板等现代架构

## 许可

[MIT](LICENSE)

---

*本项目不再活跃维护，但代码和文档保持开源状态，供社区参考和使用。*
