package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/cdle/sillyGirl/utils"
	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
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
	GetStatus  func() int                           `json:"getStatus"`
	IsComplete func() bool                          `json:"isComplete"`
	SetCookie  func(string, string, ...interface{}) `json:"setCookie"`
}

type Request struct {
	Body        func() string              `json:"body"`
	Json        func() interface{}         `json:"json"`
	IP          func() string              `json:"ip"`
	OriginalUrl func() string              `json:"originalUrl"`
	Query       func(string) string        `json:"query"`
	Querys      func() map[string][]string `json:"querys"`
	PostForm    func(string) string        `json:"postForm"`
	PostForms   func() map[string][]string `json:"postForms"`
	Path        func() string              `json:"path"`
	Header      func(string) string        `json:"header"`
	Headers     func() map[string][]string `json:"headers"`
	Method      func() string              `json:"method"`
	Cookie      func(string) string        `json:"cookie"`
}

type SillyGirlJs struct {
	BucketGet  func(bucket, key string) string                 `json:"bucketGet"`
	BucketSet  func(bucket, key, value string)                 `json:"bucketSet"`
	BucketKeys func(bucket string) []string                    `json:"bucketKeys"`
	Push       func(obj map[string]interface{})                `json:"push"`
	Session    func(wt interface{}) func(...int) SessionResult `json:"session"`
	Call       func(key string) interface{}                    `json:"call"`
}

type BucketJs struct {
	Get     func(bucket, key string) string `json:"get"`
	Set     func(bucket, key, value string) `json:"set"`
	Keys    func(bucket string) []string    `json:"keys"`
	Size    func(bucket string) int64       `json:"size"`
	Buckets func() []string                 `json:"buckets"`
	Empty   func(bucket string) bool        `json:"empty"`
}
type SessionResult struct {
	HasNext bool   `json:"hasNext"`
	Message string `json:"message"`
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

var BucketJsImpl = &BucketJs{
	Get: func(bucket, key string) string {
		return MakeBucket(bucket).GetString(key)
	},
	Set: func(bucket, key, value string) {
		bk := MakeBucket(bucket)
		bk.Set(key, value)
		if value == "" {
			empty, e := bk.Empty()
			if e != nil {
				return
			}
			if empty {
				bk.Delete()
			}
		}
	},
	Keys: func(bucket string) []string {
		ss := []string{}
		MakeBucket(bucket).Foreach(func(k, _ []byte) error {
			ss = append(ss, string(k))
			return nil
		})
		return ss
	},
	Size: func(bucket string) int64 {
		size, _ := MakeBucket(bucket).Size()
		return size
	},
	Empty: func(bucket string) bool {
		empty, _ := MakeBucket(bucket).Empty()
		return empty
	},
	Buckets: func() []string {
		ss := []string{}
		b, e := MakeBucket("").Buckets()
		if e != nil {
			return ss
		}
		for _, data := range b {
			ss = append(ss, string(data))
		}
		return ss
	},
}

var Handle = make(map[string]func(c *gin.Context))

var Server *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	Server = gin.New()
	initWeb()
}

func initWeb() {
	_, err := os.Stat(DataHome + "/views/home/hello.html")
	if os.IsNotExist(err) {
		os.MkdirAll(DataHome+"/views/home", os.ModePerm)
		os.MkdirAll(DataHome+"/assets", os.ModePerm)
		os.WriteFile(DataHome+"/views/home/hello.html", []byte(
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

	_, err = os.Stat(DataHome + "/express.js")
	var d = "`"
	if os.IsNotExist(err) {
		os.WriteFile(DataHome+"/views/home/hello.html", []byte(
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
		os.WriteFile(DataHome+"/express.js", []byte(
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
	Server.Static("/assets", DataHome+"/assets")
	Server.LoadHTMLGlob(DataHome + "/views/**/*")

	Handle["default"] = func(c *gin.Context) {
		script, err := os.ReadFile(DataHome + "/express.js")
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
						a = utils.Int(i)
					}
				}
				c.Redirect(a, b)
				isRedirect = true
			},
			Status: func(i int) goja.Value {
				status = i
				return res
			},
			SetCookie: func(name, value string, i ...interface{}) {
				c.SetCookie(name, value, 1000*60, "/", "", false, true)
			},
			IsComplete: func() bool {
				return isRedirect || len(content) > 0
			},
			GetStatus: func() int {
				return status
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
	rootPath := utils.ExecPath + "/plugin/web"
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
		info, e1 := ioutil.ReadDir(pluginPath + "/static")
		if e1 == nil && info != nil && len(info) > 0 {
			Server.Static("/"+base.Name()+"/static", pluginPath+"/static")
		}
		files, _ := ioutil.ReadDir(pluginPath)
		hasPlugin := false
		for _, v := range files {
			if v.IsDir() {
				continue
			}
			if ok, _ := regexp.MatchString("[A-z0-9]+\\.js", v.Name()); !ok {
				continue
			}
			hasPlugin = true
			break
		}
		if hasPlugin {
			Handle[base.Name()] = func(c *gin.Context) {
				p := strings.Split(c.Request.URL.Path, "/")
				if len(p) > 3 || (len(p) == 3 && "" != p[2]) {
					file, e := os.ReadFile(pluginPath + "/" + p[2] + ".js")
					before, beforee := ioutil.ReadFile(pluginPath + "/$beforeRequest.js")
					if e != nil && beforee != nil {
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
								isJson = true
								d, err := json.Marshal(ps[0])
								if err == nil {
									content += string(d)
								} else {
									content += fmt.Sprint(ps[0])
								}
							}
						},
						Header: func(str, value string) {
							c.Header(str, value)
						},
						Render: func(path string, obj map[string]interface{}) {
							isOver = true
							c.HTML(http.StatusOK, path, obj)
						},
						Redirect: func(is ...interface{}) {
							isOver = true
							a := 302
							b := ""
							for _, i := range is {
								switch i.(type) {
								case string:
									b = i.(string)
								default:
									a = utils.Int(i)
								}
							}
							c.Redirect(a, b)
						},
						Status: func(i int) goja.Value {
							status = i
							return res
						},
						SetCookie: func(name, value string, i ...interface{}) {
							time := 60 * 60 * 24 * 30 //秒,30天
							path := p[2]
							l := len(i)
							if l > 2 {
								l = 2
							}
							for j := 0; j < l; j++ {
								switch i[j].(type) {
								case string:
									path = i[j].(string)
								case int:
									time = i[j].(int)
								case int64:
									time = int(i[j].(int64))
								default:
								}
							}
							c.SetCookie(name, value, time, path, "", false, true)
						},
						IsComplete: func() bool {
							return isOver || len(content) > 0
						},
						GetStatus: func() int {
							return status
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
						vm.RunScript(file, string(js))
						return nil
					})
					vm.Set("importDir", func(dir string) error {
						return importDir(dir, pluginPath+"/", importedJs, vm)
					})
					if beforee == nil {
						_, rune := vm.RunScript("$beforeRequest.js", string(before))
						if rune != nil {
							c.String(http.StatusBadGateway, rune.Error())
							return
						}
					}
					if !isOver && err == nil && len(content) == 0 && status == 200 {
						_, rune := vm.RunScript(p[2], string(file))
						if rune != nil {
							c.String(http.StatusBadGateway, rune.Error())
							return
						}
					}
					after, aftere := ioutil.ReadFile(pluginPath + "/$afterRequest.js")
					if aftere != nil {
						_, rune := vm.RunScript("$afterRequest.js", string(after))
						if rune != nil {
							c.String(http.StatusBadGateway, rune.Error())
							return
						}
					}
					if isOver {
						return
					}
					if isJson {
						c.Header("Content-Type", "application/json")
					}
					if status == 200 && len(content) == 0 {
						c.String(404, "plugin not message")
					} else {
						c.String(status, content)
					}
				} else {
					i, e := os.Stat(pluginPath + "/static/index.html")
					if e == nil && (i == nil || !i.IsDir()) {
						c.Redirect(302, "/"+p[1]+"/static")
					} else {
						c.String(404, "plugin not find")
					}
				}
			}
		}
	}
}

type myFieldNameMapper struct{}

func (tfm myFieldNameMapper) FieldName(_ reflect.Type, f reflect.StructField) string {
	tag := f.Tag.Get(`json`)
	if idx := strings.IndexByte(tag, ','); idx != -1 {
		tag = tag[:idx]
	}
	if parser.IsIdentifier(tag) {
		return tag
	}
	return uncapitalize(f.Name)
}

func uncapitalize(s string) string {
	return strings.ToLower(s[0:1]) + s[1:]
}

func (tfm myFieldNameMapper) MethodName(_ reflect.Type, m reflect.Method) string {
	return uncapitalize(m.Name)
}

func newVm(c *gin.Context) (*goja.Runtime, *goja.Object) {
	vm := goja.New()
	vm.SetFieldNameMapper(myFieldNameMapper{})
	vm.Set("Logger", Logger)
	vm.Set("console", console)
	s := NewSillyGirl(vm)
	vm.Set("SillyGirl", func() interface{} { return s })
	vm.Set("sillyGirl", s)
	vm.Set("Request", newrequest)
	vm.Set("request", request)
	vm.Set("bucket", BucketJsImpl)
	vm.Set("fetch", request)
	vm.Set("require", require)
	var bodyData, _ = ioutil.ReadAll(c.Request.Body)
	query := c.Request.URL.Query()
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
		Querys: func() map[string][]string {
			return query
		},
		PostForm: func(s string) string {
			return c.PostForm(s)
		},
		PostForms: func() map[string][]string {
			return c.Request.PostForm
		},
		Path: func() string {
			return c.Request.URL.Path
		},
		Header: c.GetHeader,
		Headers: func() map[string][]string {
			return c.Request.Header
		},
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

func NewSillyGirl(vm *goja.Runtime) *SillyGirlJs {
	dufaultUserId := fmt.Sprintf("carry_%d", rand.Int63())
	return &SillyGirlJs{
		BucketGet: func(bucket, key string) string {
			return BucketJsImpl.Get(bucket, key)
		},
		BucketSet: func(bucket, key, value string) {
			BucketJsImpl.Set(bucket, key, value)
		},
		BucketKeys: func(bucket string) []string {
			return BucketJsImpl.Keys(bucket)
		},
		Push: func(obj map[string]interface{}) {
			imType := obj["imType"].(string)
			groupCode := 0
			var userID interface{}
			if _, ok := obj["groupCode"]; ok {
				groupCode = utils.Int(obj["groupCode"])
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
		Session: func(info interface{}) func(...int) SessionResult {
			userId := dufaultUserId
			msg := ""
			imTpye := "carry"
			chatId := 0
			switch info.(type) {
			case string:
				msg = info.(string)
			default:
				props := info.(map[string]interface{})
				for i := range props {
					switch strings.ToLower(i) {
					case "imtype":
						imTpye = props[i].(string)
					case "msg":
						msg = props[i].(string)
					case "chatid":
						chatId = utils.Int(props[i])
					case "userid":
						userId = props[i].(string)
					}
				}
			}
			if msg == "" {
				return nil
			}
			c := &Faker{
				Type:    imTpye,
				Message: msg,
				Carry:   make(chan string),
				UserID:  userId,
				ChatID:  chatId,
			}
			Senders <- c
			var f = func(i ...int) SessionResult {
				timeOut := 1000 * 100
				if len(i) > 0 {
					timeOut = i[0]
				}
				select {
				case v, ok := <-c.Listen():
					return SessionResult{
						HasNext: ok,
						Message: v,
					}
				case <-time.After(time.Millisecond * time.Duration(timeOut)):
					return SessionResult{
						HasNext: false,
						Message: "已超时",
					}
				}
			}
			return f
		},
		Call: func(key string) interface{} {
			if f, ok := OttoFuncs[key]; ok {
				return f
			}
			return nil
		},
	}
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
	var timeout time.Duration = 0
	switch wt.(type) {
	case string:
		url = wt.(string)
	default:
		props := wt.(map[string]interface{})
		for i := range props {
			switch strings.ToLower(i) {
			case "timeout":
				timeout = time.Duration(utils.Int64(props[i]) * 1000 * 1000)
			case "headers":
				headers = props[i].(map[string]interface{})
			case "method":
				method = strings.ToLower(props[i].(string))
			case "url":
				url = props[i].(string)
			case "json":
				isJson = props[i].(bool)
			case "datatype":
				switch props[i].(type) {
				case string:
					switch strings.ToLower(props[i].(string)) {
					case "json":
						isJson = true
					case "location":
						location = true
					}
				}
			case "body":
				if v, ok := props[i].(string); !ok {
					d, _ := json.Marshal(props[i])
					body = string(d)
					isJsonBody = true
				} else {
					body = v
				}
			case "formdata":
				formData = props[i].(map[string]interface{})
			case "useproxy":
				useproxy = props[i].(bool)
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
	if timeout != 0 {
		req.SetTimeout(timeout, timeout)
	}
	if isJsonBody {
		req.Header("Content-Type", "application/json")
	}
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
		h := make(map[string][]string)
		for k := range rsp.Header {
			h[k] = rsp.Header[k]
		}
		rspObj["headers"] = h
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
