##### sillyGirl

##### 傻妞机器人名

set sillyGirl name 傻妞

##### 傻妞http服务端口

set sillyGirl port 8080

##### 傻妞消息撤回等待时间，单位秒

set sillyGirl duration 5

##### 傻妞自动升级是否通知

set sillyGirl update_notify false

#### 是否开启傻妞自动更新

set sillyGirl auto_update true

##### 傻妞内置赞赏码

set sillyGirl appreciate https://gitee.com/aiancandle/sillyGirl/raw/main/appreciate.jpg

#### 是否启动http服务

set sillyGirl enable_http_server false

##### 设置青龙openapi的client_id参数

set qinglong client_id ?

##### 设置青龙openapi的client_secret参数

set qinglong client_secret ?

##### 青龙是否开启自动隐藏重复任务功能

set qinglong autoCronHideDuplicate true

##### 设置青龙面板地址

set qinglong host http://127.0.0.1:5700

##### 设置主qq账号

set qq default_bot 10000

##### 指定要监听的qq群

set qq onGroups g1&g2&g3...

##### 设置qq管理员

set qq masters q1&q2&q3...

##### 设置接受通知的qq账号

set qq notifier q1&q2&q3...

##### 设置telegram机器人token

set tg token ?

##### 设置telegram机器人代理

set tg http_proxy ?

##### 设置telegram机器人管理员

set tg masters t1&t2&t3...

##### 设置接受通知的telegram账号

set tg notifier t1&t2&t3...

##### 设置微信公众平台app_id

set wxmp app_id ?

##### 设置微信公众平台app_secret

set wxmp app_secret ?

##### 设置微信公众平台token

set wxmp token ?

##### 设置微信公众平台encoding_aes_key

set wxmp encoding_aes_key ?

##### 设置微信公众平台管理员

set wxmp masters w1&w2&w3...

##### 设置公众号关注事件回复

set wxmp subscribe_reply 感谢关注！

##### 设置公众号默认回复

set wxmp default_reply 无法回复该消息

##### 傻妞内置微信插件，依赖于[可爱猫](https://www.keaimao.com/)和[http-sdk](https://www.vwzx.com/keaimao-http-sdk)

##### 傻妞远程处理接口 /wx/receive

##### 设置插件调用地址，确保傻妞可以访问可爱猫端口

#set wx api_url ?

##### 设置图片转发模式，否则可能会出现此图片来自xx未经允许不得使用的提示

#set wx relay_mode true

##### 设置指定转发地址，格式为 https://域名/relay?url=%s，不知道不用填

#set wx relaier ?

##### 设置傻妞是否动态网络地址，适用于傻妞家庭宽带而可爱猫在云服务器的情况下

set wx sillyGirl_dynamic_ip true

##### 设置可爱猫是否动态网络地址，适用于可爱猫家庭宽带而傻妞在云服务器的情况下

#set wx keaimao_dynamic_ip true

##### 设置可爱猫端口

#set wx keaimao_port ?

![Image text](https://raw.githubusercontent.com/cdle/sillyGirl/main/appreciate.jpg)
