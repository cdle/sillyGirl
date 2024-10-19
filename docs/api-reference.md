# API 参考

本文档提供 SillyGirl 暴露的所有编程接口的详细参考，包括 REST API、gRPC API 和 JavaScript 插件 API。

## 目录

- [REST API](#rest-api)
  - [认证](#认证)
  - [Bucket API](#bucket-api)
  - [Plugin API](#plugin-api)
  - [File API](#file-api)
- [WebSocket](#websocket)
- [gRPC API](#grpc-api)
  - [服务定义](#服务定义)
  - [消息类型](#消息类型)
- [JavaScript API](#javascript-api)
  - [全局函数](#全局函数)
  - [Sender 接口](#sender-接口)
  - [Bucket 接口](#bucket-接口)
  - [Cron 接口](#cron-接口)
  - [Request / Response 接口](#request--response-接口)

## REST API

Base URL: `http://host:port/api`

### 认证

部分 API 需要认证，通过请求头传递：

```
Authorization: Bearer <token>
```

Token 获取方式：
- 首次访问 Admin 面板时自动生成
- 存储在 Bucket `sillyGirl` 的 `api_key` 字段中

### Bucket API

#### 获取键值

```
GET /api/bucket/:name/:key
```

**参数**：
- `name` - Bucket 名称
- `key` - 键名

**响应**：
```json
{
  "value": "string"
}
```

#### 设置键值

```
POST /api/bucket/:name/:key
Content-Type: application/json

{"value": "new_value"}
```

**响应**：
```json
{
  "changed": true,
  "message": ""
}
```

#### 删除键

```
DELETE /api/bucket/:name/:key
```

#### 获取所有键名

```
GET /api/bucket/:name/keys
```

**响应**：
```json
{
  "keys": ["key1", "key2"]
}
```

### Plugin API

#### 下载插件

```
GET /api/plugins/download?uuid=<uuid>
```

下载指定 UUID 的公开插件代码。

**响应**：插件 JavaScript 源代码（text/plain）或 ZIP 文件（application/zip）

```
GET /api/plugins/download/:uuid
```

同上，使用路径参数。

### File API

#### Base64 解码

```
GET /api/decode/:random
```

用于 Admin 面板的资源解码。

#### 文件查找

```
GET /api/file/:filename
```

## WebSocket

### Web 聊天

```
GET /api/web_chat?rid=<room_id>&ctt=<content>
```

这是 SillyGirl 的 WebSocket 替代实现，基于 HTTP 长轮询：

- `rid` - 房间 ID / 用户标识
- `ctt` - 要发送的消息内容（可选，为空时仅接收）

**请求示例**：

```bash
# 发送消息并接收回复
curl "http://localhost:8080/api/web_chat?rid=user123&ctt=你好"

# 仅接收消息（长轮询，最多等待4秒）
curl "http://localhost:8080/api/web_chat?rid=user123"
```

**响应格式**：

```json
[
  {
    "t": "chat",
    "c": "Hello World!",
    "m": []
  }
]
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `t` | string | 消息类型：`chat`/`info`/`warn`/`error`/`debug` |
| `c` | string | 消息内容 |
| `m` | string[] | 图片 URL 列表 |

### 管理面板实时通信

Admin 面板通过 WebSocket 与后端实时同步日志、插件状态和配置变更。连接由前端在加载时自动建立。

## gRPC API

### 服务定义

**Proto 文件**：`proto3/srpc.proto`

**Go 包**：`github.com/cdle/sillyGirl/proto3/srpc`

**服务名**：`SillyGirlService`

### 消息类型

#### Empty
```protobuf
message Empty {}
```

#### Default
```protobuf
message Default {
  string value = 1;
}
```

#### BucketSetRequest
```protobuf
message BucketSetRequest {
  string name = 1;
  string key = 2;
  string value = 3;
}
```

#### BucketSetResponse
```protobuf
message BucketSetResponse {
  bool changed = 1;
  string message = 2;
}
```

#### BucketKeyRequest
```protobuf
message BucketKeyRequest {
  string name = 1;
  string key = 2;
}
```

#### BucketRequest
```protobuf
message BucketRequest {
  string name = 1;
}
```

#### BucketKeysResponse
```protobuf
message BucketKeysResponse {
  repeated string keys = 1;
}
```

#### LenResponse
```protobuf
message LenResponse {
  int32 length = 1;
}
```

#### BoolResponse
```protobuf
message BoolResponse {
  bool value = 1;
}
```

#### BucketsResponse
```protobuf
message BucketsResponse {
  repeated string buckets = 1;
}
```

#### SenderRequest
```protobuf
message SenderRequest {
  string uuid = 1;
}
```

#### ReplyRequest
```protobuf
message ReplyRequest {
  string uuid = 1;
  string content = 2;
}
```

#### SenderContentRequest
```protobuf
message SenderContentRequest {
  string uuid = 1;
  string content = 2;
}
```

#### SenderListenRequest
```protobuf
message SenderListenRequest {
  string uuid = 1;
  repeated string rules = 2;
  int32 timeout = 3;
  bool listen_group = 4;
  bool listen_private = 5;
  bool require_admin = 6;
  repeated string allow_platforms = 7;
  repeated string prohibit_platforms = 8;
  repeated string allow_users = 9;
  repeated string prohibit_users = 10;
  repeated string allow_groups = 11;
  repeated string prohibit_groups = 12;
  bool persistent = 13;
  string value = 14;
  string plugin_id = 15;
}
```

#### SenderListenResponse
```protobuf
message SenderListenResponse {
  string echo = 1;
  string uuid = 2;
}
```

#### AdapterRegistRequest
```protobuf
message AdapterRegistRequest {
  string platform = 1;
  string bot_id = 2;
}
```

#### AdapterRequest
```protobuf
message AdapterRequest {
  string platform = 1;
  string bot_id = 2;
  string value = 3;
}
```

#### ConsoleRequest
```protobuf
message ConsoleRequest {
  string type = 1;
  string content = 2;
  string plugin_id = 3;
}
```

### RPC 方法列表

| 方法 | 请求 | 响应 | 说明 |
|------|------|------|------|
| `BucketGet` | `BucketKeyRequest` | `Default` | 获取 Bucket 键值 |
| `BucketSet` | `BucketSetRequest` | `BucketSetResponse` | 设置 Bucket 键值 |
| `BucketDelete` | `BucketRequest` | `Empty` | 删除 Bucket |
| `BucketKeys` | `BucketRequest` | `BucketKeysResponse` | 获取所有键名 |
| `BucketLen` | `BucketRequest` | `LenResponse` | 获取键数量 |
| `BucketGetAll` | `BucketRequest` | `Default` | 获取所有键值（JSON） |
| `BucketBuckets` | `Empty` | `BucketsResponse` | 获取所有 Bucket 名 |
| `BucketWatch` | `stream BucketWatchRequest` | `stream BucketWatchResponse` | 流式监听变更 |
| `SenderGetUserId` | `SenderRequest` | `Default` | 获取用户ID |
| `SenderGetUserName` | `SenderRequest` | `Default` | 获取用户名 |
| `SenderGetChatId` | `SenderRequest` | `Default` | 获取群聊ID |
| `SenderGetChatName` | `SenderRequest` | `Default` | 获取群聊名 |
| `SenderGetMessageId` | `SenderRequest` | `Default` | 获取消息ID |
| `SenderIsAdmin` | `SenderRequest` | `BoolResponse` | 是否管理员 |
| `SenderGetPlatform` | `SenderRequest` | `Default` | 获取平台 |
| `SenderGetBotId` | `SenderRequest` | `Default` | 获取机器人ID |
| `SenderGetContent` | `SenderRequest` | `Default` | 获取消息内容 |
| `SenderSetContent` | `SenderContentRequest` | `Empty` | 设置消息内容 |
| `SenderContinue` | `SenderRequest` | `Empty` | 继续匹配 |
| `SenderListen` | `stream SenderListenRequest` | `stream SenderListenResponse` | 流式消息监听 |
| `SenderEvent` | `SenderRequest` | `Default` | 获取事件数据 |
| `SenderReply` | `ReplyRequest` | `Default` | 发送回复 |
| `SenderParam` | `ReplyRequest` | `Default` | 获取参数 |
| `SenderAction` | `ReplyRequest` | `Default` | 执行动作 |
| `SenderDestroy` | `ReplyRequest` | `Empty` | 销毁 Sender |
| `AdapterRegist` | `stream AdapterRegistRequest` | `stream Default` | 注册适配器 |
| `AdapterReceive` | `AdapterRequest` | `Empty` | 接收消息 |
| `AdapterPush` | `AdapterRequest` | `Default` | 推送消息 |
| `AdapterDestroy` | `AdapterRequest` | `Empty` | 销毁适配器 |
| `AdapterSender` | `AdapterRequest` | `Default` | 获取 Sender |
| `Console` | `ConsoleRequest` | `Empty` | 控制台日志 |

### Python 客户端示例

```python
import grpc
from proto3 import srpc_pb2, srpc_pb2_grpc

channel = grpc.insecure_channel('localhost:8080')
stub = srpc_pb2_grpc.SillyGirlServiceStub(channel)

# 获取 Bucket 值
req = srpc_pb2.BucketKeyRequest(name="sillyGirl", key="port")
resp = stub.BucketGet(req)
print(f"Port: {resp.value}")

# 设置 Bucket 值
req = srpc_pb2.BucketSetRequest(name="app", key="test", value="hello")
resp = stub.BucketSet(req)
print(f"Changed: {resp.changed}")
```

## JavaScript API

### 全局函数

```js
Bucket(name: string): Bucket
Cron(): Cron
Express(): Express
sleep(ms: number): void
md5(str: string): string
uuid(): string
running(): boolean
```

### Sender 接口

`Sender` 对象通过全局变量 `s` 或 `sender` 访问。

#### 用户信息

```js
s.getUserId(): string
s.getUserName(): string
s.getChatId(): string
s.getChatName(): string
s.getMessageId(): string
s.getPlatform(): string
s.getBotId(): string
s.isAdmin(): boolean
s.getLevel(): number
s.setLevel(level: number): void
```

#### 内容操作

```js
s.getContent(): string
s.setContent(content: string): void
s.continue(): void
```

#### 回复与消息

```js
s.reply(...texts: string[]): { message_id: string, error: string }
s.recallMessage(messageId: string | string[] | string[][]): void
```

#### 参数捕获

```js
s.param(name: string): string
s.param(index: number): string
s.get(index: number): string
s.getAllMatch(): string[][]
```

#### 群管功能

```js
s.kick(userId: string): string | null
s.unkick(userId: string): string | null
s.ban(userId: string, duration: number): string | null
s.unban(userId: string): string | null
```

#### 监听

```js
s.listen(options: ListenOptions): Sender | undefined
```

`ListenOptions` 结构：

```js
{
  rules: string[],           // 匹配规则数组
  timeout?: number,          // 超时毫秒
  handle?: Function,         // 回调函数
  private?: boolean,         // 允许私聊
  group?: boolean,           // 允许群聊
  require_admin?: boolean,   // 需要管理员
  allow_platforms?: string[],
  prohibit_platforms?: string[],
  allow_users?: string[],
  allow_groups?: string[],
  prohibit_users?: string[],
  prohibit_groups?: string[],
  user_id?: string,
  chat_id?: string,
  platform?: string,
}
```

#### 其他

```js
s.holdOn(text?: string): string
s.action(options: object): { result: any, error: any }
s.doAction(options: object): { result: any, error: any }
s.getVar(key: string): any
s.setVar(key: string, value: any): void
s.setVars(kvs: object): void
s.getVars(): object
s.getReplyUserID(): number
s.isReply(): boolean
```

### Bucket 接口

```js
interface Bucket {
  get(key: string, defaultValue?: any): any
  set(key: string, value: any): Error | null
  set2(key: string, value: any): Error | null
  delete(key: string): Error | null
  keys(): string[]
  watch(key: string, callback: (old: any, new_: any, key: string) => void): void
  getAll(): Record<string, any>
  empty(): Error | undefined
  len(): number
  buckets(): string[]
}
```

### Cron 接口

```js
interface Cron {
  add(crontab: string, callback: Function): { id: number, error: string }
  remove(id: number): void
}
```

### Request / Response 接口

仅在 `@web true` 插件的 HTTP 回调中使用。

#### Request

```js
req.url(): string
req.method(): string
req.header(key: string): string
req.body(): any
req.query(key: string): string
req.path(): string
req.param(key: string): string
```

#### Response

```js
res.send(text: string): void
res.json(obj: object): void
res.redirect(url: string): void
res.status(code: number): void
```
