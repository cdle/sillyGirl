package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/dop251/goja"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Send       func(goja.Value)                     `json:"send"`
	SendStatus func(int)                            `json:"sendStatus"`
	Json       func(...interface{})                 `json:"json"`
	Header     func(string, string)                 `json:"header"`
	Render     func(string, map[string]interface{}) `json:"render"`
	Redirect   func(...interface{})                 `json:"redirect"`
	Status     func(int) goja.Value                 `json:"status"`
}

type Request struct {
	Body        func() string       `json:"body"`
	Json        func() interface{}  `json:"json"`
	IP          func() string       `json:"ip"`
	OriginalUrl func() string       `json:"originalUrl"`
	Query       func(string) string `json:"query"`
	PostForm    func(string) string `json:"postForm"`
	Path        func() string       `json:"path"`
	Header      func(string) string `json:"header"`
	Method      func() string       `json:"method"`
}

func rpo(obj *goja.Object, father string, text string, vm *goja.Runtime) string {
	for _, key := range obj.Keys() {
		v := obj.Get(key).String()
		fkey := strings.TrimLeft(father+"."+key, ".")
		text = strings.ReplaceAll(text, "#"+fkey+"#", v)
		if v == `[object Object]` {
			text = rpo(obj.Get(key).ToObject(vm), fkey, text, vm)
		}
	}
	return text
}

func init() {
	_, err := os.Stat("/etc/sillyGirl/views/home/hello.html")
	if os.IsNotExist(err) {
		os.MkdirAll("/etc/sillyGirl/views/home", os.ModePerm)
		os.MkdirAll("/etc/sillyGirl/assets", os.ModePerm)
		os.WriteFile("/etc/sillyGirl/views/home/hello.html", []byte(
			`<!DOCTYPE html>
<html lang="en">

<head>
	<meta charset="UTF-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{ .title }}</title>
	<style>
		body {
			background-image: url("{{ .data.image }}");
		}
	</style>
</head>

<body>
	{{ .data.text }}
</body>

</html>`), os.ModePerm)
	}

	_, err = os.Stat("/etc/sillyGirl/express.js")
	var d = "`"
	if os.IsNotExist(err) {
		os.WriteFile("/etc/sillyGirl/views/home/hello.html", []byte(
			`<!DOCTYPE html>
<html lang="en">

<head>
	<meta charset="UTF-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{ .title }}</title>
	<style>
		body {
			background-image: url("{{ .data.image }}");
		}
	</style>
</head>

<body>
	{{ .data.text }}
</body>

</html>`), os.ModePerm)
		os.WriteFile("/etc/sillyGirl/express.js", []byte(
			`// 获取web服务实例
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
	res.send(`+d+`你好，${name}！`+d+`)
	// 三种类型日志输出
	logs.Info(`+d+`%s，访问了 ${req.path()} 接口`+d+`, name)
	logs.Warn(`+d+`%s，访问了 ${req.path()} 接口`+d+`, name)
	logs.Debug(`+d+`%s，访问了 ${req.path()} 接口`+d+`, name)
})

// 获取表单数据
app.post('/post', (req, res) => {
	var name = req.postForm("name")
	res.send(`+d+`你好，${name}！`+d+`)
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
	sillyGirl.bucketSet(bucket, keyname, `+d+`访问地址：${req.ip()} + \n日期时间：${(new Date()).toLocaleString()}`+d+`)
})

`), os.ModePerm)
	}
	Server.Static("/assets", "/etc/sillyGirl/assets")
	Server.LoadHTMLGlob("/etc/sillyGirl/views/**/*")
	Server.NoRoute(func(c *gin.Context) {
		var status = http.StatusOK
		var content = ""
		var isJson bool
		var method = strings.ToLower(c.Request.Method)
		var bodyData, _ = ioutil.ReadAll(c.Request.Body)
		var isRedirect bool
		vm := goja.New()
		script, err := os.ReadFile("/etc/sillyGirl/express.js")
		if err != nil {
			c.String(404, err.Error())
			return
		}
		vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
		vm.Set("Logger", Logger)
		vm.Set("SillyGirl", SillyGirl)
		Render := func(path string, obj map[string]interface{}) {
			c.HTML(http.StatusOK, path, obj)
		}
		var res *goja.Object
		res = vm.ToValue(&Response{
			Send: func(gv goja.Value) {
				gve := gv.Export()
				switch gve.(type) {
				case string:
					content += gve.(string)
				default:
					d, err := json.Marshal(gve)
					if err == nil {
						content += string(d)
						isJson = true
					} else {
						content += fmt.Sprint(gve)
					}
				}
			},
			SendStatus: func(st int) {
				status = st
			},
			Json: func(ps ...interface{}) {
				if len(ps) == 1 {
					d, err := json.Marshal(ps[0])
					if err == nil {
						content += string(d)
						isJson = true
					} else {
						content += fmt.Sprint(ps[0])
					}
				}
				isJson = true
			},
			Header: func(str, value string) {
				c.Header(str, value)
			},
			Render: Render,
			Redirect: func(is ...interface{}) {
				a := 302
				b := ""
				for _, i := range is {
					switch i.(type) {
					case string:
						b = i.(string)
					default:
						a = Int(i)
					}
				}
				c.Redirect(a, b)
				isRedirect = true
			},
			Status: func(i int) goja.Value {
				status = i
				return res
			},
		}).(*goja.Object)
		req := vm.ToValue(&Request{
			Body: func() string {
				return string(bodyData)
			},
			Json: func() interface{} {
				var i interface{}
				if json.Unmarshal(bodyData, &i) != nil {
					return nil
				}
				return i
			},
			IP:          c.ClientIP,
			OriginalUrl: c.Request.URL.String,
			Query:       c.Query,
			PostForm:    c.PostForm,
			Path: func() string {
				return c.Request.URL.Path
			},
			Header: c.GetHeader,
			Method: func() string {
				return c.Request.Method
			},
		}).(*goja.Object)
		handled := false
		vm.Set("Express",
			func(call goja.ConstructorCall) *goja.Object {
				for _, m := range []string{"get", "post", "delete", "put"} {
					mm := m
					call.This.Set(mm, func(relativePath string, handle func(*goja.Object, *goja.Object)) {
						if method == mm && relativePath == c.Request.URL.Path {
							handled = true
							handle(req, res)
						}
					})
				}
				return nil
			},
		)
		_, err = vm.RunString(string(script))
		if err != nil {
			c.String(http.StatusBadGateway, err.Error())
			return
		}
		if !handled {
			c.String(404, "page nono n ot found")
			return
		}
		if isRedirect {
			return
		}
		if isJson {
			c.Header("Content-Type", "application/json")
		}

		c.String(status, content)
	})
}

func Logger(call goja.ConstructorCall) *goja.Object {
	call.This.Set("Info", func(v ...interface{}) {
		if len(v) == 0 {
			return
		}
		if len(v) == 1 {
			logs.Info(v[0])
			return
		}
		logs.Info(v[0], v[1:]...)
	})
	call.This.Set("Debug", func(v ...interface{}) {
		if len(v) == 0 {
			return
		}
		if len(v) == 1 {
			logs.Debug(v[0])
			return
		}
		logs.Debug(v[0], v[1:]...)
	})
	call.This.Set("Warn", func(v ...interface{}) {
		if len(v) == 0 {
			return
		}
		if len(v) == 1 {
			logs.Warn(v[0])
			return
		}
		logs.Warn(v[0], v[1:]...)
	})
	call.This.Set("Error", func(v ...interface{}) {
		if len(v) == 0 {
			return
		}
		if len(v) == 1 {
			logs.Error(v[0])
			return
		}
		logs.Error(v[0], v[1:]...)
	})
	return nil
}

func SillyGirl(call goja.ConstructorCall) *goja.Object {
	call.This.Set("bucketGet", func(bucket, key string) interface{} {
		return Bucket(bucket).Get(key)
	})
	call.This.Set("bucketSet", func(bucket, key, value string) {
		Bucket(bucket).Set(key, value)
	})
	call.This.Set("bucketKeys", func(bucket string) []string {
		ss := []string{}
		Bucket(bucket).Foreach(func(k, _ []byte) error {
			ss = append(ss, string(k))
			return nil
		})
		return ss
	})
	call.This.Set("push", func(obj map[string]interface{}) {
		imType := obj["imType"].(string)
		groupCode := 0
		var userID interface{}
		if _, ok := obj["groupCode"]; ok {
			groupCode = Int(obj["groupCode"])
		} else {
			userID = obj["userID"]
		}
		content := obj["content"].(string)
		if groupCode != 0 {
			if push, ok := GroupPushs[imType]; ok {
				push(groupCode, userID, content, "")
			}
		} else {
			if push, ok := Pushs[imType]; ok {
				push(userID, content, nil, "")
			}
		}
	})
	return nil
}
