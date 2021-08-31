# 配置
初次启动创建3个文件夹，分别为 `language` `plugins`
`/etc/sillyplus` 会自动生成一些启动所需的配置文件，已进行详细注释，根据自己情况来填写；

`language` 计算机语言支持依赖，请不要动它；

`plugins` 插件目录，从插件市场安装的插件都在这里；

`/etc/sillyplus` 为系统数据库存放目录；

# 面板

运行时会提示面板，本人蹩脚前端，将就用吧；

# 适配器

运行傻妞`-t`启用自带`terminal`适配器，其他适配器请前往插件市场，目前已支持`QQ`、`QQ频道`、`微信`、`公众号`、`飞书`、`钉钉`、`Telegram`、`Pagermaid`等，从`机器人`分类安装；


# 管理员命令没反应？群聊不回复群友？
参见 [**常见问题Q&A**](/help/Q&A.md)

# 常用指令

首次安装时将自动从插件市场安装，如果不需要可以禁用

```js
//获取数据库数据
get 表 key
//例如获取管理员
get qq masters
// 设置数据库
set 表 key value
//例如设置管理员
set qq masters 123456789
//获取时间
time
//启动时间
started_at
//获取机器码
machine_id
//获取版本
compiled_at
// 获取群id
chat_id
//获取个人id
myuid
//监听群消息 （默认屏蔽所有群）
listen
//屏蔽群消息
unlisten
//不回复该群
unreply
//回复该群
reply

```

# 其他命令
其他命令要视插件情况而定,具体问队友插件作者

# 关于报错!

## Error: Cannot find module 'xxxxx'
统一为缺少npm模块,通过管理员对机器人发送 npm i xxxx 命令安装模块后重启即可解决

## Error: Cannot find module './xxxxx'
统一为缺少自定义模块,谁写的插件找谁要这些模块,一般对应的插件仓库都有的,是你没装好!