package core

import (
	"archive/zip"
	"bytes"
	"context"
	"embed"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/cdle/sillyGirl/core/logs"
	"github.com/cdle/sillyGirl/core/storage"
	"github.com/cdle/sillyGirl/utils"
	"github.com/dop251/goja"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

//go:embed admin/*
var static embed.FS

var Handle = make(map[string]func(c *gin.Context))

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

var Server *gin.Engine

func initWeb() {
	for _, arg := range os.Args { //处理升级
		if arg == "-r" { //准备程序->原程序
			rfix := ".ready.exe"
			ofix := ".exe"
			if strings.Contains(os.Args[0], rfix) {
				err := utils.CopyFile(utils.ProcessName, strings.Replace(utils.ProcessName, rfix, ofix, -1))
				if err == nil {
					utils.Daemon("reset")
				}
			} else {
				os.Remove(strings.ReplaceAll(os.Args[0], ofix, rfix))
			}
			continue
		}
	}
	gin.SetMode(gin.ReleaseMode)
	Server = gin.New()
	// Server.Use(gin.Recovery())
	Server.Use(Cors())
	Server.Use(gzip.Gzip(gzip.DefaultCompression))
	Server.GET("/api/file/:filename", FindFile)
	Server.GET("/api/decode/:random", Base642Binary)

	Server.GET("/api/plugins/download", func(c *gin.Context) {
		uuid := c.Query("uuid")
		for _, f := range Functions {
			if f.UUID == uuid && f.Public {
				plugin_downloads.Set(f.UUID, plugin_downloads.GetInt(f.UUID)+1)
				if f.Type == "goja" {
					c.String(200, publicScript(plugins.GetString(f.UUID)))
					return
				} else {
					v, ok := plugins_id.Load(f.UUID)
					if !ok {
						return
					}
					dir := filepath.Dir(v.(string))
					if _, err := os.Stat(dir); err != nil { //执行压缩
						return
					}
					ss := strings.Split(dir, "/")
					name := ss[len(ss)-1]
					buf := new(bytes.Buffer)
					w := zip.NewWriter(buf)
					// dir = strings.Replace(dir, utils.ExecPath+"", ".", 1)
					err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}

						if info.IsDir() && info.Name() == "node_modules" {
							return filepath.SkipDir
						}

						if info.IsDir() {
							return nil
						}

						// 将路径转换为相对路径
						relPath, err := filepath.Rel(dir, path)
						is_index := relPath == "main.js"

						relPath = name + "/" + relPath
						if err != nil {
							return err
						}

						file, err := os.Open(path)
						if err != nil {
							return err
						}
						defer file.Close()

						fh, err := zip.FileInfoHeader(info)
						if err != nil {
							return err
						}

						// 使用相对路径作为文件名
						fh.Name = relPath

						wr, err := w.CreateHeader(fh)
						if err != nil {
							return err
						}
						if is_index {
							var data []byte
							data, err = ioutil.ReadAll(file)
							if err != nil {
								return err
							}
							su := &ScriptUtils{
								script: string(data),
							}
							if su.GetValue("public") == "true" {
								su.SetValue("public", "false")
							}
							_, err = wr.Write([]byte(su.script))
						} else {
							_, err = io.Copy(wr, file)
						}
						return err
					})
					if err != nil {
						c.String(http.StatusInternalServerError, fmt.Sprintf("ZIP creation failed: %s", err))
						return
					}
					err = w.Close()
					if err != nil {
						c.String(http.StatusInternalServerError, fmt.Sprintf("ZIP creation failed: %s", err))
						return
					}
					c.Data(http.StatusOK, "application/zip", buf.Bytes())
					// zippath := utils.ExecPath + "/public/" + f.UUID + ".zip"
					// file, err := os.Open(zippath)
					// if err != nil {
					// 	return
					// }
					// defer file.Close()
					// c.Header("Content-Type", "application/zip")
					// io.Copy(c.Writer, file)
					return
				}
			}
		}
	})
	Server.GET("/api/plugins/download/:uuid", func(c *gin.Context) {
		uuid := c.Param("uuid")
		for _, f := range Functions {
			if f.UUID == uuid && f.Public {
				plugin_downloads.Set(f.UUID, plugin_downloads.GetInt(f.UUID)+1)
				c.String(200, publicScript(plugins.GetString(f.UUID)))
				return
			}
		}
	})
	Server.NoRoute(func(c *gin.Context) {
		if c.Request.URL.Path != "/api/web_chat" {
			logs.Debug(c.Request.URL.Path)
		}
		c.Status(200)
		if strings.HasPrefix(c.Request.URL.Path, "/admin") {
			if file, err := static.Open(strings.Trim(c.Request.URL.Path, "/")); err == nil {
				fs, _ := file.Stat()
				if !fs.IsDir() {
					defer file.Close()
					c.Header("cache-control", "max-age=864000")
					io.Copy(c.Writer, file)
					return
				} else {
					file.Close()
				}
			}
			data, err := static.ReadFile("admin/index.html")
			if err == nil {
				c.Header("Content-Type", "text/html; charset=utf-8")
				c.Writer.Write(data)
				return
			}
		}
		for _, req := range ss {
			if c.Request.URL.Path == req.Path && (req.Method == c.Request.Method || req.Method == "ANY") {
				req.Handle(c)
				return
			}
		}
		uuid, _ := c.Cookie("uuid")
		if uuid == "" {
			uuid = utils.GenUUID()
			c.SetCookie("uuid", uuid, 8640000, "/", "", false, true)
		}
		if IsWebSocketRequest(c.Request) {
			handleWebsocket(c)
			return
		}
		var req = &Request{
			c:    c,
			uuid: uuid,
		}
		var res = &Response{
			c: c,
		}
		for _, function := range Functions {
			if len(function.Https) != 0 {
				for _, http := range function.Https {
					path := http.Path
					method := http.Method
					matched := false
					if strings.HasPrefix(path, "^") {
						reg, err := regexp.Compile(path)
						if err != nil {
							console.Error(err)
							continue
						}
						req.ress = reg.FindAllStringSubmatch(c.Request.URL.Path, -1)
						matched = len(req.ress) != 0
					}
					if (matched || c.Request.URL.Path == path) && (c.Request.Method == method || "ANY" == method) {
						req.handled = true
						function.Handle(&Faker{
							Type: "http",
						}, func(vm *goja.Runtime) {
							vm.Set("res", res)
							vm.Set("req", req)
						})
						goto HELL
					}
				}
			}
		}
	HELL:

		if req.handled {
			if res.isRedirect {
				return
			}
			if res.isJson {
				c.Header("Content-Type", "application/json")
			}
			c.String(res.status, res.content)
			return
		}

		httpListensAny.Range(func(_, value any) bool {
			web := value.(*HttpListen)
			if web.Closed {
				return true
			}
			req.handled = true
			var end = make(chan bool, 1)
			web.Chan <- &RR{
				Req: req,
				Res: res,
				End: func() {
					if !web.Closed {
						end <- true
					}
				},
			}
			select {
			case <-end:
				// fmt.Println("end", req.handled)
			case <-time.After(time.Second * 10):
				logs.Warn("%s 响应超时", c.Request.URL.Path)
			}
			close(end)
			if req.handled {
				if res.isRedirect {
					return false
				}
				if res.isJson {
					c.Header("Content-Type", "application/json")
				}
				// fmt.Println(res.status, res.content)
				c.String(res.status, res.content)
				return false
			}
			return true
		})

		if req.handled {
			return
		}
		httpListens.Range(func(key, value any) bool {
			web := value.(*HttpListen)
			if web.Closed {
				return true
			}
			// fmt.Println(c.Request.URL.Path, key, web.Method, c.Request.Method)
			if web.Method == c.Request.Method {
				var matched = web.Path == c.Request.URL.Path
				if !matched && strings.HasPrefix(web.Path, "^") {
					reg, err := regexp.Compile(web.Path)
					if err != nil {
						console.Error(err)
						return true
					}
					req.ress = reg.FindAllStringSubmatch(c.Request.URL.Path, -1)
					matched = len(req.ress) != 0
				}
				if matched {
					httpListens.Delete(key)
					req.handled = true
					var end = make(chan bool, 1)
					web.Chan <- &RR{
						Req: req,
						Res: res,
						End: func() {
							if !web.Closed {
								end <- true
							}
						},
					}
					select {
					case <-end:
					case <-time.After(time.Second * 10):
						logs.Warn("%s 响应超时", c.Request.URL.Path)
					}
					web.Closed = true
					close(end)
					if !web.Closed {
						close(web.Chan)
					}
					if req.handled {
						if res.isRedirect {
							return false
						}
						if res.isJson {
							c.Header("Content-Type", "application/json")
						}
						// fmt.Println(res.status, res.content)
						c.String(res.status, res.content)
						return false
					}
				}
			}
			return true
		})
		if req.handled {
			return
		}
		for _, web := range webs {
			if web.handles == nil {
				continue
			}
			if web.method == c.Request.Method {
				var matched = web.path == c.Request.URL.Path
				if !matched && strings.HasPrefix(web.path, "^") {
					reg, err := regexp.Compile(web.path)
					if err != nil {
						console.Error(err)
						continue
					}
					req.ress = reg.FindAllStringSubmatch(c.Request.URL.Path, -1)
					matched = len(req.ress) != 0
				}
				if matched {
					req.handled = true
					func() {
						defer func() {
							if err := recover(); err != nil {
								if fmt.Sprint(err) != "stop" {
									console.Error(err)
								}
							}
						}()
						for _, handle := range web.handles {
							handle(req, res)
						}
					}()
					if req.handled {
						if res.isRedirect {
							return
						}
						if res.isJson {
							c.Header("Content-Type", "application/json")
						}
						// fmt.Println(res.status, res.content)
						c.String(res.status, res.content)
						return
					}
				}
			}
		}

		c.String(404, "页面被喵咪劫走了") //
		//开启代理模式

		// handleHTTP(c.Writer, c.Request)
	})

	port := sillyGirl.GetInt("port")
	if port == 0 {
		sillyGirl.Set("port", 8080)
		port = 8080
	}
	srvs := []*http.Server{{
		Addr:    ":" + fmt.Sprint(port),
		Handler: Server,
	}}

	storage.Watch(sillyGirl, "port", func(old, new, key string) *storage.Final {
		if new == "" {
			new = "8080"
		}
		if old == new {
			return nil
		}
		port := new
		console.Log("port", new)
		srv := &http.Server{
			Addr:    ":" + port,
			Handler: Server,
		}
		var ch = make(chan error, 1)
		srvs = append(srvs, srv)

		go func() {
			logs.Info("Http服务(%v)重新运行", port)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logs.Error("Http服务(%v)运行失败：%s", port, err.Error())
				ch <- err
			}
		}()
		select {
		case err := <-ch:
			srvs = srvs[:len(srvs)-1]
			return &storage.Final{
				Error: err,
			}
		case <-time.After(1 * time.Millisecond * 100):
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := srvs[0].Shutdown(ctx); err == nil {
				logs.Info("Http服务(%v)关闭", old)
			}
			srvs = srvs[1:]
		}
		return &storage.Final{
			Now: new,
		}
	})

	// logs.Info("Http服务(%s)开始运行", port)

	logs.Info("管理员面板:")
	logs.Info("  > 本机: http://localhost:%d/admin", port)
	local_ip := getLocalIP()
	logs.Info("  > 局域网: http://%v:%d/admin", local_ip, port)
	ip := sillyGirl.GetString("ip")
	if ip != "" {
		logs.Info("  > 广域网: http://%v:%d/admin", ip, port)
	}
	sillyGirl.Set("local_ip", local_ip)
	go func() {
		if err := srvs[0].ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.Error("Http服务运行失败：%s", err.Error())
		}
	}()
}

type Req struct {
	Method string
	Path   string
	Handle func(c *gin.Context)
}

var ss = []Req{}

type Auth struct {
	ID        int
	IP        string
	UserAgent string
	Token     string
	CreatedAt int
	ExpiredAt int
}

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
	ANY    = "ANY"
)

func GinApi(method string, path string, fs ...func(c *gin.Context)) {
	ss = append(ss, Req{
		Method: method,
		Path:   path,
		Handle: func(c *gin.Context) {
			defer func() {
				recover()
			}()
			for _, f := range fs {
				f(c)
			}
		},
	})
}

func getLocalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "127.0.0.1"
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip.To4() != nil {
				return ip.String()
			}
		}
	}
	return "127.0.0.1"
}
