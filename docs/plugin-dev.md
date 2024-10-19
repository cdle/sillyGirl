# 插件开发指南

SillyGirl 的插件系统基于 [Goja](https://github.com/dop251/goja) JavaScript 引擎，完整支持 ECMAScript 5.1。插件以 `.js` 文件形式存在，通过顶部注释声明元数据，由框架自动加载、匹配和执行。

## 目录

- [插件结构](#插件结构)
- [注释元数据](#注释元数据)
- [全局对象与 API](#全局对象与-api)
  - [sender (s)](#sender-s)
  - [Bucket(name)](#bucketname)
  - [Cron()](#cron)
  - [Express()](#express)
  - [其他全局函数](#其他全局函数)
- [规则匹配语法](#规则匹配语法)
- [消息监听与会话](#消息监听与会话)
- [HTTP 路由](#http-路由)
- [定时任务](#定时任务)
- [完整示例](#完整示例)
- [调试技巧](#调试技巧)

## 插件结构

一个最小插件由注释元数据和执行代码组成：

```js
/**
 * @title HelloWorld
 * @rule raw ^你好$
 */

s.reply("Hello World!");
```

插件文件可以包含多个 `@rule`，每条规则独立匹配：

```js
/**
 * @title 多功能助手
 * @rule raw ^你好$
 * @rule raw ^再见$
 * @rule 天气 [城市]
 */

const content = s.getContent();
if (content === "你好") {
  s.reply("Hello!");
} else if (content === "再见") {
  s.reply("Goodbye!");
} else {
  const city = s.param("城市");
  s.reply(`${city}今天天气晴朗！`);
}
```

## 注释元数据

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `title` | string | 是 | 插件标题，显示在管理面板和插件市场 |
| `rule` | string | 否 | 消息匹配规则，支持多行。详见[规则匹配语法](#规则匹配语法) |
| `cron` | string | 否 | 定时任务表达式，支持多行，格式为 `平台 cron表达式` |
| `priority` | number | 否 | 匹配优先级，数字越大越优先，默认 0 |
| `on_start` | boolean | 否 | `true` 时随系统启动执行一次，常用于初始化服务 |
| `module` | boolean | 否 | `true` 时表示为模块插件，不响应消息规则 |
| `web` | boolean | 否 | `true` 时表示为 Web 插件，可声明 HTTP 路由 |
| `version` | string | 否 | 版本号，如 `v1.0.0` |
| `author` | string | 否 | 作者名 |
| `description` | string | 否 | 插件描述 |
| `icon` | string | 否 | 插件图标 URL |
| `public` | boolean | 否 | `true` 时允许发布到插件市场 |
| `disable` | boolean | 否 | `true` 时禁用插件 |
| `admin` | boolean | 否 | `true` 时仅管理员可触发 |
| `platform` | string | 否 | 限制平台，如 `qq`、`web` |

### 元数据示例

```js
/**
 * @title 每日早报
 * @rule raw ^早报$
 * @cron qq 0 9 * * *
 * @cron web 0 9 * * *
 * @priority 10
 * @version v1.2.0
 * @author cdle
 * @description 每天早上9点推送新闻早报
 * @icon https://example.com/icon.png
 * @public true
 */
```

## 全局对象与 API

### sender (s)

当前消息的 Sender 对象，是插件中最核心的交互入口。

#### 用户信息

```js
s.getUserId()       // 获取用户ID（string）
s.getUserName()     // 获取用户昵称（string）
s.getChatId()       // 获取群聊ID，私聊时为空字符串（string）
s.getChatName()     // 获取群聊名称（string）
s.getMessageId()    // 获取消息ID（string）
s.getPlatform()     // 获取平台类型，如 "qq"、"web"（string）
s.getBotId()        // 获取当前机器人ID（string）
s.isAdmin()         // 判断用户是否为管理员（boolean）
```

#### 内容操作

```js
s.getContent()      // 获取消息原始内容（string）
s.setContent(text)  // 修改当前消息内容（影响后续插件匹配）
s.continue()        // 继续匹配后续规则（默认匹配成功即停止）
```

#### 回复消息

```js
s.reply("文本")     // 回复文本消息，返回 { message_id, error }
s.reply("文本1", "文本2")  // 多段回复
```

#### 参数捕获

```js
// 对于规则 "天气 [城市]"
s.param("城市")     // 通过名称获取捕获组
s.param(1)          // 通过索引获取第1个捕获组（从1开始）
s.get(1)            // param 的别名
s.getAllMatch()     // 获取所有匹配组（二维数组）
```

#### 群管功能

```js
s.kick(userId)          // 踢出群成员
s.unkick(userId)        // 取消踢出
s.ban(userId, duration) // 禁言，duration 为秒数
s.unban(userId)         // 解除禁言
s.recallMessage(messageId)  // 撤回消息
```

### Bucket(name)

持久化键值存储，数据自动持久化到 BoltDB 或 Redis。

```js
const bucket = Bucket("myapp");

// 基础读写
bucket.set("key", "value");
const value = bucket.get("key", "default_value");  // 支持默认值

// 类型自动转换
bucket.set("count", 100);        // 存储为数字
bucket.set("enabled", true);     // 存储为布尔值
bucket.set("data", { a: 1 });    // 存储为对象（自动 JSON 序列化）

// 监听变更
bucket.watch("key", (oldValue, newValue, key) => {
  console.log(`${key} changed: ${oldValue} -> ${newValue}`);
});

// 其他方法
bucket.keys();       // 获取所有键名（string[]）
bucket.getAll();     // 获取所有键值（object）
bucket.delete("key"); // 删除键
bucket.empty();      // 清空 bucket
bucket.len();        // 获取键数量（number）
```

**作用域说明**：每个 Bucket 是独立的命名空间，不同插件建议使用不同的 Bucket 名称，避免键冲突。

### Cron()

定时任务调度器，支持标准 Cron 表达式。

```js
const task = Cron();

// 添加任务（6字段：秒 分 时 日 月 周）
const { id, error } = task.add("*/5 * * * * *", () => {
  console.log("每5秒执行一次");
});

// 添加任务（5字段：分 时 日 月 周，秒自动补0）
const { id, error } = task.add("0 9 * * *", () => {
  console.log("每天早上9点执行");
});

// 移除任务
task.remove(id);
```

> 注意：插件中更推荐通过 `@cron` 注释声明定时任务，由框架统一管理生命周期。

### Express()

HTTP 服务路由注册（基于内置 Web 服务器）。

```js
const app = Express();

app.get("/hello", (req, res) => {
  res.send("Hello World!");
});

app.post("/api/data", (req, res) => {
  const body = req.body();
  res.json({ success: true, data: body });
});
```

> 注意：Web 插件需要声明 `@web true` 元数据，且框架会自动清理已注册的路由。

### 其他全局函数

```js
sleep(ms)           // 同步阻塞睡眠（毫秒）
md5(str)            // MD5 哈希
uuid()              // 生成 UUID
running()           // 判断当前插件是否仍在运行（boolean）
fmt.Sprintf(format, ...args)  // 格式化字符串
fmt.Printf(format, ...args)   // 格式化输出
time.Now()          // 获取当前时间对象
time.Sleep(ms)      // 睡眠（毫秒）
time.Unix(sec)      // Unix 时间戳转时间对象
time.Parse(str, layout, locale)  // 解析时间字符串
```

## 规则匹配语法

规则（`@rule`）支持多种语法，框架会自动转换为正则表达式进行匹配。

### 基础规则

```js
/**
 * @rule raw ^你好$      // 原始正则，完全匹配"你好"
 * @rule raw ^/help$     // 原始正则，匹配 "/help"
 * @rule 你好            // 自动转换为 ^你好$，完全匹配
 * @rule ^天气           // 以"天气"开头
 * @rule 帮助$           // 以"帮助"结尾
 */
```

### 参数捕获

```js
/**
 * @rule 天气 [城市]              // 匹配"天气 北京"，捕获"北京"
 * @rule [操作:登录,注册,退出]     // 匹配"登录"、"注册"或"退出"
 * @rule 查询 ?                   // ? 匹配任意非空白字符
 * @rule 搜索 *                   // * 匹配任意内容（包括空白）
 */
```

### 可选参数

```js
/**
 * @rule 天气 [城市?]             // ? 表示可选，匹配"天气"或"天气 北京"
 */
```

### 优先级与冲突

当多条规则可能同时匹配一条消息时，框架按以下顺序决定执行：

1. **监听规则**（`s.listen` 注册的）优先于普通规则
2. **高 priority** 优先于低 priority
3. 同一优先级下，先加载的插件优先

使用 `s.continue()` 可以让当前插件执行完毕后继续匹配后续规则：

```js
/**
 * @title 日志中间件
 * @rule *
 * @priority 999
 */

console.log("收到消息:", s.getContent());
s.continue();  // 继续让其他插件处理
```

## 消息监听与会话

`s.listen()` 是实现对话式交互的核心 API，支持等待用户后续输入并按规则匹配。

### 基础用法

```js
/**
 * @title 注册流程
 * @rule raw ^注册$
 */

s.reply("请输入你的用户名：");
const result = s.listen({
  rules: ["raw ^(.+)$"],  // 捕获任意输入
  timeout: 30000,          // 30秒超时
  handle: (s2) => {
    const username = s2.param(1);
    s2.reply(`用户名 "${username}" 注册成功！`);
  },
});

if (!result) {
  s.reply("注册超时，请重试。");
}
```

### 监听选项

```js
s.listen({
  rules: ["规则1", "规则2"],        // 匹配规则数组
  timeout: 10000,                   // 超时时间（毫秒）
  handle: (s2) => { ... },          // 匹配后的回调函数
  private: true,                    // 允许在私聊中触发
  group: true,                      // 允许在群聊中触发
  require_admin: false,             // 是否要求管理员权限
  allow_platforms: ["qq"],          // 限制平台
  prohibit_platforms: ["web"],      // 禁止平台
  allow_users: ["12345"],           // 仅允许指定用户
  allow_groups: ["67890"],          // 仅允许指定群组
  user_id: s.getUserId(),           // 仅监听当前用户
  chat_id: s.getChatId(),           // 仅监听当前群组
});
```

### 持久化监听

添加 `"persistent"` 参数可创建长期有效的监听器：

```js
/**
 * @title 关键词监控
 * @rule raw ^启动监控$
 */

s.listen({
  rules: ["raw ^(.+)$"],
  timeout: 0,           // 0 表示永不超时
  handle: (s2) => {
    console.log("监控到消息:", s2.getContent());
    // 返回空字符串或 undefined 会继续监听
  },
}, "persistent");

s.reply("监控已启动");
```

### HoldOn 与 GoAgain

```js
/**
 * @title 循环输入
 * @rule raw ^开始$
 */

function ask() {
  s.reply("请输入内容（输入'结束'停止）：");
  const r = s.listen({
    rules: ["raw ^(.+)$"],
    timeout: 30000,
    handle: (s2) => {
      if (s2.param(1) === "结束") {
        s2.reply("已结束");
        return;
      }
      s2.reply("你输入了：" + s2.param(1));
      return s2.holdOn("开始");  // 重新触发当前插件
    },
  });
}
ask();
```

## HTTP 路由

声明 `@web true` 的插件可以注册 HTTP 路由：

```js
/**
 * @title WebAPI 示例
 * @web true
 */

Express().get("/api/status", (req, res) => {
  res.json({
    status: "ok",
    time: new Date().toISOString(),
    plugin: "WebAPI 示例",
  });
});

Express().post("/api/echo", (req, res) => {
  const body = req.body();
  res.json({ echo: body });
});
```

### Request / Response 对象

```js
req.url()           // 请求 URL
req.method()        // 请求方法
req.header("key")   // 获取请求头
req.body()          // 获取请求体（已解析 JSON）
req.query("key")    // 获取 URL 参数
req.path()          // 请求路径
req.param("key")    // 获取路径参数

res.send("text")    // 发送文本响应
res.json(obj)       // 发送 JSON 响应
res.redirect(url)   // 重定向
res.status(code)    // 设置状态码
```

## 定时任务

通过 `@cron` 注释声明定时任务：

```js
/**
 * @title 每日提醒
 * @cron qq 0 9 * * *       // 每天早上9点（QQ平台）
 * @cron web 0 12 * * *     // 每天中午12点（Web平台）
 */

s.reply("该起床啦！今天是 " + time.Now().Format("2006-01-02"));
```

Cron 表达式格式（5字段或6字段）：

```
秒(可选) 分 时 日 月 周
```

| 字段 | 范围 | 特殊字符 |
|------|------|----------|
| 秒 | 0-59 | `, - * /` |
| 分 | 0-59 | `, - * /` |
| 时 | 0-23 | `, - * /` |
| 日 | 1-31 | `, - * /` |
| 月 | 1-12 | `, - * /` |
| 周 | 0-6（0=周日）| `, - * /` |

示例：

| 表达式 | 说明 |
|--------|------|
| `*/5 * * * *` | 每5分钟 |
| `0 */1 * * *` | 每小时 |
| `0 9 * * 1-5` | 工作日早上9点 |
| `0 0 1 * *` | 每月1号零点 |
| `*/30 * * * * *` | 每30秒 |

## 完整示例

### 示例 1：记忆名字

```js
/**
 * @title 记忆名字
 * @rule raw ^我是谁$
 * @rule 我是[姓名]
 * @version v1.0.0
 * @author cdle
 */

const user = Bucket("user_names");
const name = s.param("姓名");

if (!name) {
  const stored = user.get(s.getUserId());
  if (stored) {
    s.reply(`你是 ${stored}`);
  } else {
    s.reply("我还不知道你是谁，告诉我吧：我是[你的名字]");
  }
} else {
  user.set(s.getUserId(), name);
  s.reply(`好的，我记住你了，${name}！`);
}
```

### 示例 2：倒计时提醒

```js
/**
 * @title 倒计时
 * @rule 倒计时 [分钟:1,2,3,5,10] 分钟
 * @version v1.0.0
 */

const minutes = parseInt(s.param(1));
s.reply(`好的，${minutes}分钟后提醒你。`);

// 使用 Cron 不太适合一次性延时，这里用 sleep
// 注意：sleep 会阻塞，长时间任务建议用 cron 或其他方式
go(() => {
  sleep(minutes * 60 * 1000);
  s.reply("⏰ 时间到了！");
});
```

> 注：实际开发中长时间后台任务建议使用 `@cron` 或外部调度。

### 示例 3：简易 HTTP API

```js
/**
 * @title 天气查询 API
 * @web true
 * @version v1.0.0
 */

Express().get("/api/weather", (req, res) => {
  const city = req.query("city") || "北京";
  // 实际场景中这里应调用第三方天气 API
  res.json({
    city: city,
    weather: "晴",
    temperature: "25°C",
    updated_at: time.Now().Format("2006-01-02 15:04:05"),
  });
});
```

## 调试技巧

### 1. 使用终端模式

```bash
./sillyGirl -t
```

终端模式是开发插件的最快方式，修改插件后立即生效，无需重启。

### 2. 开启调试模式

在 Admin 面板将 `sillyGirl.debug` 设为 `true`，或在插件中：

```js
console.log("调试信息:", someVariable);
console.debug("详细调试:", detailedInfo);
console.error("错误:", err);
```

### 3. 查看插件状态

在 Admin 面板的"插件"页面，可以查看：
- 插件加载状态
- 规则列表
- 错误日志
- 性能统计

### 4. 安全执行

插件运行异常时框架会自动捕获 panic，不会影响其他插件或核心服务。你可以在 Admin 面板查看具体的错误堆栈。

### 5. 模块复用

将公共逻辑抽取为 `@module true` 插件，其他插件通过 `require` 或全局变量复用：

```js
/**
 * @title 工具模块
 * @module true
 */

// 定义全局工具函数
RegistFuncs["utils"] = {
  formatTime: (t) => t.Format("2006-01-02"),
  isWeekend: (t) => t.Weekday() === 0 || t.Weekday() === 6,
};
```
