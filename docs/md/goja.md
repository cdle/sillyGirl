# 傻妞

一个不太有用的机器人，不生产消息，只搬运消息。

## 特性

- 简单易用的消息搬运功能。
- 简单强大的自定义回复功能。
- 完整支持 ECMAScript 5.1 的插件系统，基于 [otto](https://github.com/robertkrimen/otto)。
- 支持通过内置的阉割版 `Express` / `request` ，接入互联网。
- 内置 `Cron` ，轻松实现定时任务。
- 持久化的 `Bucket` 存储模块。
- 支持同时接入多个平台多个机器人，自己开发。

## 快速上手

### 安装

在 [releases](https://github.com/cdle/sillyGirl/releases) 中找到合适自己系统版本的程序运行带 `-t` 可以开启终端机器人，直接与程序进行交互。

```shell
./sillyplus -t
2023/05/24 14:12:01.859 [I]  默认使用boltdb进行数据存储。
2023/05/24 14:12:01.950 [I]  Http服务已运行(8080)。
```

### 开发第一个插件

```js
/**
 * @title HelleWorld
 * @rule raw ^你好$
 */

s.reply("Helle World!");
```

怼着程序输入 `你好` ，就可以看到机器人回复的 `Helle World!` 了

```sh
你好
2023/05/24 14:15:48.350 [I]  匹配到规则：^你好$
Helle World!
```

插件注释 `@rule raw ^你好$` 中的正则表达式被消息匹配时插件脚本就会被触发。

### 添加和销毁定时任务

插件注释 `@service true` 时是作为傻妞系统服务，伴随机器人一起启动。

```js
/**
 * @title 定时任务
 * @service true
 */

const task = Cron();
let taskId = 0;
let times = 5;
const { id } = task.add("*/5 * * * * *", () => {
  // 同样支持分钟级任务，如：*/5 * * * *
  times--;
  console.log(
    `每5秒执行一次任务，${
      times ? `${times}次后结束任务` : "这是最后一次任务"
    }。`
  );
  if (times == 0) {
    task.remove(taskId); //移除任务
  }
});
taskId = id;
```

程序输出：

```sh
2023/05/27 19:57:00.000 [I]  每5秒执行一次任务，4次后结束任务。
2023/05/27 19:57:05.001 [I]  每5秒执行一次任务，3次后结束任务。
2023/05/27 19:57:10.001 [I]  每5秒执行一次任务，2次后结束任务。
2023/05/27 19:57:15.001 [I]  每5秒执行一次任务，1次后结束任务。
2023/05/27 19:57:20.000 [I]  每5秒执行一次任务，这是最后一次任务。
```

### 与用户交互

```js
/**
 * @title 用户交互插件
 * @rule 猜拳
 */

s.reply("你先出，请在10秒内出拳！");
ns = s.listen({
  rules: ["[出拳:剪刀,石头,布]"], // []中出拳是参数名，剪刀,石头,布是参数可能值
  timeout: 10000, // 超时设置
  handle: (s) => {
    let choose = s.param("出拳");
    s.reply(
      `我出${
        choose == "石头" ? "剪刀" : choose == "布" ? "剪刀" : "石头"
      }，我赢了。`
    );
  },
});
if (!ns) {
  s.reply("你没出拳，算我赢了！");
}
```

### 开发 HTTP 接口

```js
/**
 * @title 第一个web服务
 * @service true
 */

const app = Express(); //导入HTTP服务，傻妞默认开启，端口8080
app.get("/helloWorld", function (req, res) {
  res.send("Hello world!");
});
```

打开浏览器访问 `http://127.0.0.1:8080/helloWorld` ，当然地址根据实际情况，理论上可以看到接口返回的 `Hello world!` 。

### 实现一个 HTTP 请求

```js
/**
 * @title 实现一个HTTP 请求
 * @service true
 */

let api = "/testRequest"; //接口地址

//第一步，实现一个原样返回请求数据的接口
const app = Express();
app.post(api, (req, res) => res.json(req.json()));

//第二步，请求第一步实现的接口
const port = Bucket("app").port ?? "8080"; // 获取http服务端口
const url = `http://127.0.0.1:${port}${api}`;
const { body } = request({
  url,
  method: "POST",
  body: { value: "test" },
  json: true,
});
console.log(`value is ${body.value}`);
```

### 持久化存储

```js
/**
 * @title 持久化存储
 * @rule raw ^我是谁$
 * @rule 我是[姓名]
 */

const user = Bucket("user"); //初始化存储桶user
let name = s.param("姓名");

if (user.name == "") {
  s.reply(`我不知道你是谁！`);
} else if (name == "谁") {
  s.reply(`你是${user.name}`);
} else {
  user.name = name;
  s.reply(`好的，你的姓名更新为${user.name}`);
}
```

插件实现了记名字的功能，其中`姓名`是方括号里匹配到的值，本质还是正则匹配到的。

```
我是谁
2023/05/24 15:43:40.121 [I]  匹配到规则：^我是谁$
我不知道你是谁！
我是小千
2023/05/24 15:43:49.735 [I]  匹配到规则：^我是([\s\S]+)$
好的，你的姓名更新为小千
我是谁
2023/05/24 15:43:53.727 [I]  匹配到规则：^我是谁$
你是小千
```

有了 `Bucket` 才有了傻妞从不认识小千到认识小千的过程。

### 管理员

```js
const masters = Bucket("qq")["masters"];
```

`masters` 是管理员账号通过"&"拼接起来的，系统默认依此判断用户是否是管理员。

### 群组消息

默认不监听不回复任何群组，监听口令 `listen` 和 `unlisten`，回复口令 `reply` 和 `noreply`，需要管理员在对应群组发送口令。

## 深入了解

### 插件注释

| 字段          | 举例                               | 用法                                               |
| ------------- | ---------------------------------- | -------------------------------------------------- |
| `title`       | HelloWorld                         | 插件标题                                           |
| `rule`        | raw `^我是([\s\S]+)$`              | 可写多行，取括号内参数 `s.param(1)` ，多个参数类推 |
| `priority`    | `1`                                | 插件优先级，越高则优先处理                         |
| `service`     | `true`                             | 插件后台任务执行脚本，避免重复运行                 |
| `disable`     | `true`                             | 禁用脚本                                           |
| `form`        | `{title: "姓名", key:"user.name"}` | 插件表，key 值对应 `存储桶.键名`                   |
| `public`      | `true`                             | 公开插件                                           |
| `create_at`   | 2023-05-24 15:14:53                | 插件创建时间                                       |
| `description` | 本插件用于每天向女友问好           | 插件描述                                           |
| `author`      | `cdle`                             | 插件作者                                           |
| `version`     | `v1.0.0`                           | 插件版本                                           |
| `icon`        | url 省略...                        | 给插件增加图标                                     |

### Sender

傻妞搬运的核心对象，在插件中为全局变量 s or sender。

```ts
interface Sender {
  getUserId(): string; //获取用户ID
  getUserName(): string; //获取用户昵称
  getChatId(): string; //获取群聊ID
  getChatName(): string; //获取群聊名称
  getMessageId(): string; //获取消息ID
  getContent(): string; //获取消息内容
  continue(): void; //使消息继续往下匹配正则，消息正常第一次被匹配就会停止继续匹配
  setContent(content: string): void; //修改接收到的消息内容，可配合`continue`被其他规则匹配
  param(index: string | number): string; //获取`rule`匹配参数，可取[]内参数，?型参数从1开始取，例 `@rule 回复 ?` 对应 `s.param(1)`
  holdOn(content: string): string; //持续监听
  listen({
    rules: string[]; //匹配规则
    timeout: number; //超时，单位毫秒
    handle: (s: Sender): string;//如果匹配成功，则进入消息处理逻辑。如果将 holdOn(content) 的结果作为返回值，会继续监听
    listen_private: boolean; //监听用户群内消息时，同时监听用户消息
    listen_group: boolean; //监听用户群内消息时，同时监听群员消息
    allow_platforms: string[]; //平台白名单
    prohibit_platforms: string[]; //平台黑名单
    allow_groups: string[]; //群聊白名单
    prohibit_groups: string[]; //群聊黑名单
    allow_users: string[]; //用户白名单
    prohibit_users: string[]; //群聊白名单
  }): Sender; //超时，返回undefined
  isAdmin(): boolean; //判断消息是否来自管理员
  getPlatform(): string; //获取消息平台
  getBotId(): string; //获取机器人ID
  reply(content: string): {message_id: string, error: string}; //回复消息，媒体消息推荐使用CQ码实现，返回消息ID
  recallMessage(meesageId: string | string[] | number): {error: string}; //撤回消息，number类型时为延时毫秒
  kick(user_id: string): {error: string}; //移出群聊
  unkick(user_id: string): {error: string}; //取消移出群聊
  ban(user_id: string, duration: number): {error: string}; //禁言，并指定时长
  unban(user_id: string): {error: string};  //取消禁言
}
```

### Express `Request` / `Response`

只能说是够用，有需求可联系作者。插件中通过 `Express()` 返回一个对象。

```ts
interface Request {
  body(): string; //获取请求体
  json(): any; //将请求体解析为JSON
  ip(): string; //获取客户端IP地址
  originalUrl(): string; //获取原始请求URL
  query(param: string): string; //获取查询参数
  param(i: number): string; //根据索引获取路径参数
  querys(): Record<string, string[]>; //获取所有查询参数
  postForm(s: string): string; //获取表单数据
  postForms(): Record<string, string[]>; //获取所有表单数据
  path(): string; //获取请求路径
  header(s: string): string; //获取请求头
  get(s: string): string; //获取请求头
  headers(): Record<string, string[]>; //获取所有请求头
  method(): string; //获取请求方法
  cookie(s: string): string; //获取 cookie
  cookies(): Record<string, string>; //获取 cookies
  continue(): void; //继续匹配其他路由
  setSession(k: string, v: string): string; //设置会话值
  getSession(k: string): string; //获取会话值
  getSessionId(): string; //获取会话ID
  destroySession(): string; //销毁会话
  logined(): boolean; //是否面板登录状态
}

interface Response {
  send(body: any): Response; //发送响应体
  sendStatus(status: number): Response; //发送状态码
  json(...ps: any[]): Response; //发送JSON响应
  header(str: string, value: string): Response; //设置响应头
  set(str: string, value: string): void; //设置响应头
  render(view: string, params: Record<string, any>): Response; //渲染视图
  redirect(...is: any[]): void; //重定向到URL
  status(i: number, ...s: string[]): Response; //设置状态码和文本
  setCookie(name: string, value: string, ...i: any[]): Response; //设置 Cookie
  stop(): void; //代码片段停止
}
```

### request

由 `net/http` 封装而成，如有更多需求可以联系作者。

```ts
function request(options: {
  url: string; //请求地址
  method: string; //请求方法
  headers: { [key: string]: string }; //请求头
  json: boolean; // 返回json对象，等价于 responseType: "json"
  timeout: number; //超时参数，单位毫秒
  form: { [key: string]: any }; //formData表单数据，优先于下面的body
  body: any; // 请求体，支持字符串、二进制，对象自动转json字符串和添加相应请求头
  allow_redirects: boolean; // 是否允许重定向，默认允许
  proxy: {};
}): {
  status: number; // 状态码，同statusCode
  headers: { [key: string]: string };
  body: any;
};
```

### Adapter

```ts
interface Message{
  message_id: string; // 消息ID
  user_id: string;    // 用户ID
  chat_id: string;    // 聊天ID
  content: string;    // 聊天内容
  user_name: string;  // 用户名
  chat_name: string;  // 群组名
}

class Adapter(botplt: string, botid: string) {
  isAdapter(botid: string): boolean; //判断id是否为机器人
  push(message: Message): string; //推送消息，无视禁言设置
  getReplyMessage(): Promise<message: Message>; //获取一条回复消息，实际发送成功后，如果有id，请设置 message.message_id
  setReplyHandler(func: (message: Message): string): void; //设置回复事件处理方法，方法中返回消息ID，不推荐使用。
  receive(message: Message): Sender; //接收一个消息，并返回一个Sender对象
  setRecallMessage(func: (i: string | string[]) => boolean): void;//设置撤回消息函数。
  setGroupKick(func: (user_id: string, chat_id: string, reject_add_request: boolean) => void): boolean; //设置群聊成员移除函数，reject_add_request指5是否继续接受请求
  setGroupBan(func: (user_id: string, chat_id: string, duration: number) => void): boolean;//设置群聊成员禁言函数
  setGroupUnban(func: (user_id: string, chat_id: string) => void): boolean;//设置群聊成员解除禁言函数
  setIsAdmin(func: (user_id: string) => boolean): void; //设置用户是否是成员函数，默认自动实现
  destroy(): void;//销毁机器人
}
function getAdapter(platform: string, bot_id string): {Adapter: string, error: string}; //获取一个机器人

function getAdapterBotsID(bot_id string): Adapter[]; //获取一个平台的所有机器人

function getAdapterBotPlts(platform: string): string[]; //所有机器人平台
```

### Bucket

例：通过 `Bucket("app")` 初始化一个 app 存储痛

```ts
interface Bucket(name: string) {
  get(key: string, defaultValue: any): any; // 取值
  set(key: string, value: any): Error | null; // 设值
  watch(key: string, event: (old: any, new_: any, key: string) => void); // 设置监听器，key 值为 * 时将监听整个桶的存储事件
  getAll(): []; // 获取全部值
  delete(key: string): Error | null; // 删值
  empty(): Error | undefined; // 清空桶
  keys(): string[]; // 获取所有键名
  len(): number | undefined; // 获取数据数目
  buckets(): string[]; // 获取所有存在的桶名
  _name(): string; // 获取当前桶名
}
```

### Cron

可以通过`let task = Cron()`返回的对象来添加定时任务 `const {id, error} = task.add("* * * * *", ()=>{})`

```ts
interface Cron {
  add(crontab: string, ()=>void): {id: number, error: string}//添加定时任务 crontab同时支持秒级和分钟级
  remove(id: number): void//移除定时任务
}
```

### 插件表单

可以使用注释 `@form {title: "标题", key: "test.title"}` 添加表单元素。当如也可以直接在插件代码中添加，如下。

```js
// 单个表单元素
Form({
  title: "姓名",
  key: "test.name",
});
// 多个表单元素
Form([
  {
    title: "姓名",
    key: "test.name",
  },
  {
    title: "性别",
    key: "test.sex",
  },
]);
// 使用schema-form
Form([
  {
    title: "创建时间",
    key: "test.createName",
    dataIndex: "test.createName",
    valueType: "date",
  },
  {
    title: "创建时间",
    key: "test.createName",
    dataIndex: "test.createName",
    valueType: "date",
  },
  {
    title: "分组",
    valueType: "group",
    columns: [
      {
        title: "状态",
        dataIndex: "test.groupState",
        valueType: "select",
        width: "xs",
        valueEnum: {
          all: { text: "全部", status: "Default" },
          open: {
            text: "未解决",
            status: "Error",
          },
          closed: {
            text: "已解决",
            status: "Success",
            disabled: true,
          },
          processing: {
            text: "解决中",
            status: "Processing",
          },
        },
      },
      {
        title: "标题",
        width: "md",
        dataIndex: "test.groupTitle",
        formItemProps: {
          rules: [
            {
              required: true,
              message: "此项为必填项",
            },
          ],
        },
      },
    ],
  },
]);
```

### 其他

```ts
sleep(millsec: number): void; //等待
md5(string): string; //加密
running(): boolean; //服务是否运行
genUuid(): string; //生成uuid
```

### 项目赞助

打开微信扫一扫，深入了解作者~
![](https://raw.githubusercontent.com/cdle/sillyGirl/main/appreciate.jpg)
