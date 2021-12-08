// 获取web服务实例
var app = Express();
// 获取日志实例
var logs = Logger();
// 获取傻妞实例
var sillyGirl = SillyGirl();

// 首页
app.get("/", (req, res) => {
     // 渲染模版
     res.render(
          "hello.html",// 模版文件目录 /etc/sillyGirl/views
          {
               title: "世界，你好。", data: {
                    text: "Hello world!",
                    image: "assets/test.jpeg",// 静态文件目录 /etc/sillyGirl/assets
               }
          }
     )

     // 页面提示404
     // res.status(404).send("页面找不到了")

     // 跳转指定网页
     // res.redirect("https://github.com/cdle/sillyGirl")
})

// 响应普通文本
app.get('/text', (req, res) => {
     res.send('这是一段普通的文字。')
})

// 获取请求的json数据，响应json数据
app.post('/json', (req, res) => {
     var data = req.json()
     res.json(data)
})

// 获取url中的参数
app.get('/query', (req, res) => {
     var name = req.query("name")
     res.send(`你好，${name}！`)
     // 三种类型日志输出
     logs.Info(`%s，访问了 ${req.path()} 接口`, name)
     logs.Warn(`%s，访问了 ${req.path()} 接口`, name)
     logs.Debug(`%s，访问了 ${req.path()} 接口`, name)
})

// 获取表单数据
app.post('/post', (req, res) => {
     var name = req.postForm("name")
     res.send(`你好，${name}！`)
})

// 推送私聊消息
app.get('/sendPrivateMsg', (req, res) => {
     sillyGirl.push({
          imType: "tg",
          userID: "1837585653",
          content: "你的大香蕉成熟了，请快到app领取。"
     })
})

// 推送群聊消息
app.post('/sendGroupMsg', (req, res) => {
     sillyGirl.push({
          imType: "tg",
          groupCode: -1001583071436,
          content: "该喝开水啦。"
     })
})

// 数据存储
app.get('/lastTime', (req, res) => {
     var bucket = "test"
     var keyname = "lastTime"
     var lastTime = sillyGirl.bucketGet(bucket, keyname)
     res.send(lastTime)
     sillyGirl.bucketSet(bucket, keyname, `访问地址：${req.ip()} + \n日期时间：${(new Date()).toLocaleString()}`)
})

