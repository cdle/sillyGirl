# 架构设计文档

本文档深入解析 SillyGirl 的内部架构、核心模块设计与数据流。

## 目录

- [整体设计哲学](#整体设计哲学)
- [模块架构图](#模块架构图)
- [核心模块详解](#核心模块详解)
  - [Adapter Manager](#adapter-manager)
  - [Message Router](#message-router)
  - [Plugin Engine](#plugin-engine)
  - [Bucket Storage](#bucket-storage)
  - [Web Server](#web-server)
  - [gRPC Services](#grpc-services)
  - [Cron Scheduler](#cron-scheduler)
- [数据流](#数据流)
- [插件生命周期](#插件生命周期)
- [存储抽象设计](#存储抽象设计)
- [安全与隔离](#安全与隔离)

## 整体设计哲学

SillyGirl 的设计围绕三个核心原则：

1. **插件优先**：所有业务逻辑都通过 JavaScript 插件实现，核心框架仅提供基础设施
2. **平台无关**：通过 Adapter 抽象屏蔽底层平台差异，同一插件可同时服务于 QQ、Web 等多个平台
3. **运行时动态**：支持插件热重载、配置热更新、存储后端热切换，最大化运维灵活性

## 模块架构图

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              Client Layer                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐ │
│  │  QQ Client  │  │ Web Browser │  │  gRPC Client│  │  Pagermaid Bot  │ │
│  │ (CQHTTP/WS) │  │  (Admin UI) │  │  (Python/JS)│  │  (Python Bridge)│ │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └────────┬────────┘ │
└─────────┼────────────────┼────────────────┼──────────────────┼──────────┘
          │                │                │                  │
          └────────────────┴────────────────┴──────────────────┘
                                    │
                    ┌───────────────▼────────────────┐
                    │        Transport Layer          │
                    │    (HTTP / WebSocket / TCP)     │
                    └───────────────┬────────────────┘
                                    │
┌───────────────────────────────────▼─────────────────────────────────────┐
│                            Adapter Layer                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    Adapter Manager (core/adapter.go)               │   │
│  │  • Factory：适配器工厂，管理每个平台实例的生命周期                   │   │
│  │  • Bots Map：[platform, botid] → Factory 的线程安全映射            │   │
│  │  • Sender：统一的消息发送接口，封装平台差异                         │   │
│  │  • CustomSender：标准 Sender 实现，供所有适配器复用                │   │
│  └─────────────────────────────────────────────────────────────────┘   │
└───────────────────────────────────┬─────────────────────────────────────┘
                                    │ Messages chan
                    ┌───────────────▼────────────────┐
                    │         Message Router          │
                    │      (core/function.go)         │
                    │  • Group Filter (listen/reply)  │
                    │  • User Blocklist               │
                    │  • Admin Commands               │
                    └───────────────┬────────────────┘
                                    │
┌───────────────────────────────────▼─────────────────────────────────────┐
│                           Core Engine Layer                              │
│                                                                          │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────────┐  │
│  │  Plugin Engine  │  │  Bucket Storage │  │    Cron Scheduler       │  │
│  │ (plugin_core.go)│  │  (bucket.go)    │  │    (cron.go)            │  │
│  │  • Goja VM Pool │  │  • BoltDB       │  │  • robfig/cron/v3       │  │
│  │  • Rule Compile │  │  • Redis        │  │  • Per-platform cron    │  │
│  │  • Priority Q   │  │  • MongoDB      │  │  • Lifecycle mgmt       │  │
│  │  • Hot Reload   │  │  • Watch/Notify │  │                         │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────────────┘  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────────┐  │
│  │   Web Server    │  │  gRPC Services  │  │    Log Framework        │  │
│  │    (web.go)     │  │ (grpc_*.go)     │  │    (core/logs/)         │  │
│  │  • Gin Router   │  │  • srpc.proto   │  │  • Multi-backend        │  │
│  │  • Static Embed │  │  • Go/Python/JS │  │  • Level/Format/Rotate  │  │
│  │  • WS Handler   │  │    clients      │  │                         │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────────────┘  │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │              Node.js Compatibility Layer (node_*.go)              │   │
│  │  • request / crypto / os / regexp / xml / buffer / strings        │   │
│  └─────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
```

## 核心模块详解

### Adapter Manager

**文件**：`core/adapter.go`

Adapter Manager 负责管理所有平台机器人实例的生命周期。

**核心数据结构**：

```go
type Factory struct {
    botid      string              // 机器人标识
    botplt     string              // 平台类型 (qq/web/...)
    uuid       string              // 唯一标识
    msgChan    chan MsgChan        // 消息接收通道
    reply      func(map[string]interface{}) string  // 回复处理函数
    isAdmin    func(string) bool   // 管理员判断函数
    vm         *goja.Runtime       // 关联的 JS 运行时
    ctx        context.Context     // 生命周期上下文
    cancel     context.CancelFunc  // 取消函数
    destroid   bool                // 是否已销毁
}

var Bots = map[Bot]*Factory{}  // 全局适配器映射，Bot = [2]string{platform, botid}
```

**关键行为**：
- **初始化** (`Init`)：注册到全局 `Bots` 映射，初始化消息通道
- **冲突处理**：同一 `[platform, botid]` 已存在时，自动销毁旧实例
- **健康检查**：连续错误超过 5 次自动销毁
- **销毁** (`Destroy`)：清理资源、关闭通道、从映射移除

### Message Router

**文件**：`core/function.go`

Message Router 是消息流转的核心枢纽，负责过滤、匹配和分发。

**消息处理流程**：

1. **接收**：从 `Messages` channel 接收 `Sender` 对象
2. **过滤层**：
   - 群组监听检查 (`listenOnGroups`)：非监听群组的消息默认忽略（管理员除外）
   - 用户屏蔽检查 (`noListenUsers`)：屏蔽用户的消息直接丢弃
   - 群管指令处理：管理员发送 `listen`/`unlisten`/`reply`/`noreply` 等指令时直接处理
3. **昵称记录**：自动记录用户和群组的昵称映射
4. **分发**：通过 `go HandleMessage(s)` 异步处理

**HandleMessage 内部**：

1. **等待器匹配** (`waits`)：优先检查是否有插件通过 `s.listen()` 等待此消息
2. **自动回复** (`replies`)：检查关键词自动回复规则
3. **插件匹配** (`Functions`)：按优先级遍历所有插件，正则匹配 `rule`
4. **消息撤回** (`recall`)：匹配配置的正则则自动撤回消息

### Plugin Engine

**文件**：`core/plugin_core.go`, `core/plugin_impl.go`

Plugin Engine 是 SillyGirl 最核心的创新点，实现了在 Go 程序中运行 JavaScript 插件的完整机制。

**插件加载流程**：

1. **扫描存储**：从 `plugins` Bucket 读取所有已保存的插件代码
2. **解析元数据**：通过正则提取注释中的 `@title`、`@rule` 等字段
3. **编译脚本**：使用 `goja.Compile()` 预编译为 `goja.Program`
4. **注册命令**：调用 `AddCommand()` 将插件加入全局 `Functions` 数组
5. **规则格式化**：`fmtRule()` 将声明式规则转换为标准正则表达式

**执行模型**：

```
消息到达
  → 匹配 Function
  → 创建新 eventloop.NewEventLoop()
  → 在 loop 中运行 goja.Program
  → 注入全局对象 (s, Bucket, Cron, ...)
  → 执行用户代码
  → 回收 VM
```

每个插件执行都在独立的 EventLoop 中，通过 `recover()` 捕获 panic，确保单个插件崩溃不影响系统。

**JS 全局对象注入**（`SetPluginMethod`）：

- `s` / `sender` → `SenderJsIplm`（封装 Go 的 Sender 接口）
- `Bucket(name)` → 返回存储桶对象
- `Cron()` → 返回定时任务管理器
- `initAdapter()` / `InitAdapter()` → 创建新适配器
- `running()` → 查询插件运行状态
- `uuid()` / `genUUID()` → 生成 UUID

### Bucket Storage

**文件**：`core/bucket.go`, `core/storage/`

Bucket 是 SillyGirl 的统一存储抽象，设计灵感来自 AWS S3 的 Bucket 概念。

**接口定义**：

```go
type Bucket interface {
    Set(interface{}, interface{}) (string, bool, error)
    Set2(interface{}, interface{}) (string, bool, error)
    Copy(string) Bucket
    GetString(...interface{}) string
    GetBytes(string) []byte
    GetInt(string, ...int) int
    GetBool(string, ...bool) bool
    Foreach(func([]byte, []byte) error)
    Keys() ([]string, error)
    // ...
}
```

**类型系统**：

Bucket 通过前缀标记实现类型透明存储：

| 前缀 | 类型 | 示例 |
|------|------|------|
| 无 | string | `hello` |
| `d:` | int | `d:100` |
| `f:` | float | `f:3.14` |
| `b:` | bool | `b:true` |
| `o:` | object (JSON) | `o:{"a":1}` |

**监听机制**：

```go
storage.Watch(bucket, "key", func(old, new, key string) *Final {
    // 值变更时触发
    return &Final{Message: "配置已更新"}
})
```

监听是全局的，所有 `Watch` 注册器保存在 `Listens` 切片中，每次 `Set` 操作后遍历触发。

### Web Server

**文件**：`core/web.go`

基于 Gin 的 HTTP 服务，承担多个职责：

1. **Admin 面板**：嵌入的 React 静态资源（`//go:embed admin/*`）
2. **REST API**：`/api/plugins/download`、文件服务等
3. **WebSocket**：`/api/web_chat` 长轮询实现实时聊天
4. **插件 HTTP 路由**：插件通过 `@web` 或 `Express()` 注册的路由
5. **动态端口**：支持运行时通过 Bucket 修改监听端口，平滑重启

**NoRoute 处理逻辑**：

```
请求到达
  → 匹配 /admin/* 静态资源
  → 匹配 GinApi 注册的固定路由
  → 匹配 WebSocket 升级
  → 匹配插件 @web HTTP 路由
  → 匹配动态 httpListens（Express 注册）
  → 返回 404
```

### gRPC Services

**文件**：`proto3/srpc.proto`, `core/grpc_*.go`

gRPC 服务使 SillyGirl 可以被其他语言编写的客户端调用。

**服务列表**：

| 服务 | 文件 | 说明 |
|------|------|------|
| `SillyGirlService` | `grpc_sender.go` | 主服务，包含所有 RPC |
| Bucket 操作 | `grpc_bucket.go` | 存储读写、监听 |
| Adapter 管理 | `grpc_adapter.go` | 适配器注册、接收、推送 |
| Plugin 管理 | `grpc_plugins.go` | 插件安装、卸载、列表 |
| Asset 管理 | `grpc_asset.go` | 静态资源 |
| Queue 服务 | `grpc_queue.go` | 队列操作 |
| Runtime 服务 | `grpc_runtime.go` | 运行时控制 |

**多语言支持**：

- Go：`proto3/srpc/srpc.pb.go`
- Python：`proto3/srpc_pb2.py`, `proto3/srpc_pb2_grpc.py`
- JavaScript/TypeScript：`proto3/srpc.js`, `proto3/srpc.ts`
- Node.js：`proto3/sillygirl.js`

### Cron Scheduler

**文件**：`core/cron.go`

基于 `robfig/cron/v3` 的定时任务调度，支持平台隔离：

```go
type Function struct {
    Cron map[string]string  // platform -> cron expression
    CronIds []int           // 已注册的 cron entry IDs
}
```

每个 `@cron` 规则会在 `AddCommand()` 时被解析并注册到全局 `CRON` 调度器。当定时触发时，会构造一个虚拟的 `Sender` 对象（平台类型为 `"cron"`）传入插件处理函数。

## 数据流

### 消息接收与回复

```
┌─────────┐    HTTP/WS     ┌──────────┐    Messages chan    ┌─────────────┐
│ QQ/Web  │───────────────→│ Adapter  │────────────────────→│ Message     │
│ Platform│                │ Factory  │                     │ Router      │
│         │←───────────────│          │                     │             │
└─────────┘   Reply()      └──────────┘                     └──────┬──────┘
                                                                   │
                                                                   ▼
┌─────────┐    msgChan     ┌──────────┐    rule match       ┌─────────────┐
│ Adapter │←───────────────│ Plugin   │←───────────────────│ HandleMessage│
│ (await) │                │ Engine   │                     │             │
└─────────┘                └──────────┘                     └─────────────┘
```

### 插件执行时序

```
Message Received
  → HandleMessage()
    → Waiters Check (s.listen)
    → Auto Reply Check
    → for _, function := range Functions
      → regexp.Match(rule, content)
      → if matched:
        → eventloop.NewEventLoop()
        → loop.Run(func(vm) {
            → SetPluginMethod(vm)
            → vm.RunProgram(prg)
          })
        → sender.Finish()
        → if sender.IsAtLast() → sender.Reply(accumulated)
```

## 插件生命周期

```
[插件代码写入 storage]
         │
         ▼
[storage.Watch 触发]
         │
         ▼
[initPlugin()] ──→ [pluginParse()] ──→ [goja.Compile()]
         │
         ▼
[AddCommand()] ──→ [fmtRule()] ──→ [cron.AddFunc()]
         │
         ▼
[加入 Functions 数组]
         │
    ┌────┴────┐
    ▼         ▼
[消息匹配]  [定时触发]
    │         │
    ▼         ▼
[f.Handle()] (Goja VM 执行)
    │
    ▼
[插件输出 / 错误捕获]
```

**卸载流程**：

```
plugins.Set(uuid, "") 或 plugins.Set(uuid, "uninstall")
  → storage.Watch 触发
  → 查找 Functions 中对应 UUID
  → 移除 cron jobs
  → 移除 HTTP listeners
  → 取消 waits
  → 从 Functions 切片删除
  → 销毁关联 adapters
```

## 存储抽象设计

SillyGirl 的存储层采用三层架构：

```
┌─────────────────────────────────────────┐
│           Application Layer              │
│         (plugin JS / core logic)         │
└─────────────────┬───────────────────────┘
                  │ Bucket 接口
┌─────────────────▼───────────────────────┐
│         Storage Abstraction              │
│    (core/storage/main.go Bucket iface)   │
└─────────────────┬───────────────────────┘
                  │
      ┌───────────┼───────────┐
      ▼           ▼           ▼
┌─────────┐ ┌─────────┐ ┌─────────┐
│ BoltDB  │ │  Redis  │ │ MongoDB │
│(boltdb) │ │ (redis) │ │(mongodb)│
└─────────┘ └─────────┘ └─────────┘
```

**选择策略**：
- 默认使用 BoltDB，零配置启动
- 通过 `sillyGirl.storage` 配置切换后端
- 切换前会自动测试 Redis 连通性
- 数据不自动迁移，切换后端后原数据不可见

## 安全与隔离

### 插件隔离

- **VM 隔离**：每个插件运行在独立的 Goja Runtime 中
- **panic 捕获**：`recover()` 包裹所有插件入口，单个插件崩溃不影响其他插件
- **资源限制**：通过 EventLoop 控制执行，避免死循环阻塞主线程

### 权限控制

- **管理员系统**：每个平台独立维护 `masters` 列表，以 `&` 分隔
- **指令保护**：`listen`/`unlisten` 等群管指令仅管理员可用
- **插件级权限**：`@admin true` 的插件仅管理员可触发

### 网络安全

- **CORS**：默认开启跨域支持，方便前端开发
- **API Key**：Admin 面板和 REST API 使用 Token 认证
- **UUID Cookie**：Web 用户通过 Cookie 标识身份
