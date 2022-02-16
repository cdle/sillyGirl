package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/client/httplib"
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
	SetCookie  func(string, string)                 `json:"setCookie"`
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
	Cookie      func(string) string `json:"cookie"`
}

type SillyGirlWeb struct {
	BucketGet  func(bucket, key string) interface{} `json:"bucketGet"`
	BucketSet  func(bucket, key, value string)      `json:"bucketSet"`
	BucketKeys func(bucket string) []string         `json:"bucketKeys"`
	Push       func(obj map[string]interface{})     `json:"push"`
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

var Handle = make(map[string]func(c *gin.Context))

func init() {

	_, err := os.Stat(dataHome + "/views/home/hello.html")
	if os.IsNotExist(err) {
		os.MkdirAll(dataHome+"/views/home", os.ModePerm)
		os.MkdirAll(dataHome+"/assets", os.ModePerm)
		os.WriteFile(dataHome+"/views/home/hello.html", []byte(
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

	_, err = os.Stat(dataHome + "/express.js")
	var d = "`"
	if os.IsNotExist(err) {
		os.WriteFile(dataHome+"/views/home/hello.html", []byte(
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
		os.WriteFile(dataHome+"/express.js", []byte(
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
	Server.Static("/assets", dataHome+"/assets")
	Server.LoadHTMLGlob(dataHome + "/views/**/*")

	Handle["default"] = func(c *gin.Context) {
		script, err := os.ReadFile(dataHome + "/express.js")
		if err != nil {
			c.String(404, err.Error())
			return
		}
		vm, req := newVm(c)
		var status = http.StatusOK
		var content = ""
		var isJson bool
		var isRedirect bool
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
			SetCookie: func(name, value string) {
				c.SetCookie(name, value, 1000*60, "/", "", false, true)
			},
		}).(*goja.Object)
		var method = strings.ToLower(c.Request.Method)
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
	}
	Server.NoRoute(func(c *gin.Context) {
		patchPostForm(c)
		p := c.Request.URL.Path
		var f func(c *gin.Context) = nil
		if len(p) > 1 {
			split := strings.Split(p, "/")
			f = Handle[split[1]]
		}
		if f == nil {
			f = Handle["default"]
		}
		f(c)
	})
	initWebPlugin()
}

func initWebPlugin() {
	//请求接口插件化为目录:
	//pluginRoot
	// - dir1 //目录名称做为请求路径
	// -- static //静态文件目录
	// -- *.js //接口本体,请求路径为/dir/js名称,只支持2级
	// - dir2 //目录2
	rootPath := ExecPath + "/plugin/web"
	rootFiles, err := ioutil.ReadDir(rootPath)
	if err != nil {
		os.MkdirAll(rootPath, os.ModePerm)
		return
	}
	for _, base := range rootFiles {
		if !base.IsDir() {
			continue
		}
		if ok, _ := regexp.MatchString("[A-z0-9]+", base.Name()); !ok {
			continue
		}
		pluginPath := path.Join(rootPath, base.Name())
		files, _ := ioutil.ReadDir(pluginPath)
		var plugin []string
		for _, v := range files {
			if v.IsDir() {
				continue
			}
			if ok, _ := regexp.MatchString("[A-z0-9]+\\.js", v.Name()); !ok {
				continue
			}
			plugin = append(plugin, strings.TrimPrefix(v.Name(), ".js"))
		}
		if len(plugin) > 0 {
			_, err := ioutil.ReadDir(pluginPath + "/static")
			if err == nil {
				Server.Static("/"+base.Name()+"/static", pluginPath+"/static")
			}
			Handle[base.Name()] = func(c *gin.Context) {
				p := strings.Split(c.Request.URL.Path, "/")
				if len(p) > 3 || (len(p) == 3 && "" != p[2]) {
					file, e := os.ReadFile(pluginPath + "/" + p[2] + ".js")
					if e != nil {
						c.String(404, "plugin not find")
						return
					}
					vm, _ := newVm(c)
					var status = 200
					var isOver bool
					var isJson bool
					var content = ""
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
						Render: func(path string, obj map[string]interface{}) {
							isOver = true
							c.HTML(http.StatusOK, path, obj)
						},
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
							isOver = true
						},
						Status: func(i int) goja.Value {
							status = i
							return res
						},
						SetCookie: func(name, value string) {
							c.SetCookie(name, value, 1000*60, "/", "", false, true)
						},
					}).(*goja.Object)
					vm.Set("__response", res)
					importedJs := make(map[string]struct{})
					importedJs[pluginPath+"/"+p[2]+".js"] = struct{}{}
					vm.Set("importJs", func(file string) error {
						js, e := ReadJs(file, pluginPath+"/", importedJs)
						if e != nil {
							return e
						}
						vm.RunString(string(js))
						return nil
					})
					vm.Set("importDir", func(dir string) error {
						return importDir(dir, pluginPath, importedJs, vm)
					})
					_, err = vm.RunString(string(file))
					if err != nil {
						c.String(http.StatusBadGateway, err.Error())
						return
					}
					if isOver {
						return
					}
					if isJson {
						c.Header("Content-Type", "application/json")
					}

					c.String(status, content)
				} else {
					_, err := os.Stat(pluginPath + "/static/index.html")
					if err != nil {
						c.String(404, "plugin not find")
					} else {
						c.Redirect(302, "/"+p[1]+"/static")
					}
				}
			}
		}
	}
}

func newVm(c *gin.Context) (*goja.Runtime, *goja.Object) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	vm.Set("Logger", Logger)
	vm.Set("console", console)
	vm.Set("SillyGirl", SillyGirl)
	s := vm.ToValue(&SillyGirlWeb{
		BucketGet: func(bucket, key string) interface{} {
			return Bucket(bucket).Get(key)
		},
		BucketSet: func(bucket, key, value string) {
			Bucket(bucket).Set(key, value)
		},
		BucketKeys: func(bucket string) []string {
			ss := []string{}
			Bucket(bucket).Foreach(func(k, _ []byte) error {
				ss = append(ss, string(k))
				return nil
			})
			return ss
		},
		Push: func(obj map[string]interface{}) {
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
		},
	}).(*goja.Object)
	vm.Set("sillyGirl", s)
	vm.Set("Request", newrequest)
	vm.Set("request", request)
	vm.Set("fetch", request)
	vm.Set("require", require)
	var bodyData, _ = ioutil.ReadAll(c.Request.Body)
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
		PostForm: func(s string) string {
			return c.PostForm(s)
		},
		Path: func() string {
			return c.Request.URL.Path
		},
		Header: c.GetHeader,
		Method: func() string {
			return c.Request.Method
		},
		Cookie: func(s string) string {
			var cookie, _ = c.Cookie(s)
			return cookie
		},
	}).(*goja.Object)
	vm.Set("__request", req)
	return vm, req
}

func patchPostForm(c *gin.Context) {
	if c.Request.Method == "POST" {
		c.Request.ParseForm()
	}
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

func newrequest() interface{} {
	return request
}

func require(str string) interface{} {
	switch str {
	case "request":
		return request
	}
	return nil
}

func request(wt interface{}, handles ...func(error, map[string]interface{}, interface{}) interface{}) interface{} {
	var method = "get"
	var url = ""
	var req *httplib.BeegoHTTPRequest
	var headers map[string]interface{}
	var formData map[string]interface{}
	var isJson bool
	var isJsonBody bool
	var body string
	var location bool
	var useproxy bool
	switch wt.(type) {
	case string:
		url = wt.(string)
	default:
		props := wt.(map[string]interface{})
		for i := range props {
			switch i {
			case "headers":
				headers = props["headers"].(map[string]interface{})
			case "method":
				method = strings.ToLower(props["method"].(string))
			case "url":
				url = props["url"].(string)
			case "json":
				isJson = props["json"].(bool)
			case "dataType":
				switch props["dataType"].(string) {
				case "json":
					isJson = true
				case "location":
					location = true
				}
			case "body":
				if v, ok := props["body"].(string); !ok {
					d, _ := json.Marshal(props["body"])
					body = string(d)
					isJsonBody = true
				} else {
					body = v
				}
			case "formData":
				formData = props["formData"].(map[string]interface{})
			case "useproxy":
				useproxy = props["useproxy"].(bool)
			}
		}
	}
	switch strings.ToLower(method) {
	case "post":
		req = httplib.Post(url)
	case "put":
		req = httplib.Put(url)
	case "delete":
		req = httplib.Delete(url)
	default:
		req = httplib.Get(url)
	}
	if isJsonBody {
		req.Header("Content-Type", "application/json")
	}
	//自定义header优先级最高
	for i := range headers {
		req.Header(i, fmt.Sprint(headers[i]))
	}
	for i := range formData {
		req.Param(i, fmt.Sprint(formData[i]))
	}
	if body != "" {
		req.Body(body)
	}
	if location {
		req.SetCheckRedirect(func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		})
		rsp, err := req.Response()
		if err == nil && (rsp.StatusCode == 301 || rsp.StatusCode == 302) {
			return rsp.Header.Get("Location")
		} else
		//非重定向,允许用户自定义判断
		if len(handles) == 0 {
			return err
		}
	}
	if useproxy && Transport != nil {
		req.SetTransport(Transport)
	}
	rsp, err := req.Response()
	rspObj := map[string]interface{}{}
	var bd interface{}
	if err == nil {
		rspObj["status"] = rsp.StatusCode
		rspObj["statusCode"] = rsp.StatusCode
		data, _ := ioutil.ReadAll(rsp.Body)
		if isJson {
			var v interface{}
			json.Unmarshal(data, &v)
			bd = v
		} else {
			bd = string(data)
		}
		rspObj["body"] = bd
		rspObj["header"] = rsp.Header
	}
	if len(handles) > 0 {
		return handles[0](err, rspObj, bd)
	} else {
		return bd
	}
}

var console = map[string]func(...interface{}){
	"info": func(v ...interface{}) {
		if len(v) == 0 {
			return
		}
		if len(v) == 1 {
			logs.Info(v[0])
			return
		}
		logs.Info(v[0], v[1:]...)
	},
	"debug": func(v ...interface{}) {
		if len(v) == 0 {
			return
		}
		if len(v) == 1 {
			logs.Debug(v[0])
			return
		}
		logs.Debug(v[0], v[1:]...)
	},
	"warn": func(v ...interface{}) {
		if len(v) == 0 {
			return
		}
		if len(v) == 1 {
			logs.Warn(v[0])
			return
		}
		logs.Warn(v[0], v[1:]...)
	},
	"error": func(v ...interface{}) {
		if len(v) == 0 {
			return
		}
		if len(v) == 1 {
			logs.Error(v[0])
			return
		}
		logs.Error(v[0], v[1:]...)
	},
	"log": func(v ...interface{}) {
		if len(v) == 0 {
			return
		}
		if len(v) == 1 {
			logs.Info(v[0])
			return
		}
		logs.Info(v[0], v[1:]...)
	},
}
