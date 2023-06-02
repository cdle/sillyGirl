package core

import (
	"context"
	"embed"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/cdle/sillyGirl/core/logs"
	"github.com/cdle/sillyGirl/core/storage"
	"github.com/cdle/sillyGirl/utils"
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

var Server = gin.New()

func init() {
	gin.SetMode(gin.ReleaseMode)
	// Server.Use(gin.Recovery())
	Server.Use(Cors())
	Server.Use(gzip.Gzip(gzip.DefaultCompression))
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
		var req = &Request{
			c:    c,
			uuid: uuid,
		}
		var res = &Response{
			c: c,
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
					// fmt.Println(c.Request.URL.Path, key)
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

		c.String(404, "页面被喵咪劫走了。") //
		//开启代理模式

		// handleHTTP(c.Writer, c.Request)
	})

	port := sillyGirl.GetString("port", "8080")
	srvs := []*http.Server{{
		Addr:    ":" + port,
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
			logs.Info("Http服务(%s)重新运行。", port)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logs.Error("Http服务(%s)运行失败：%s", port, err.Error())
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
				logs.Info("Http服务(%s)关闭。", old)
			}
			srvs = srvs[1:]
		}
		return &storage.Final{
			Now: new,
		}
	})
	go func() {
		logs.Info("Http服务(%s)开始运行。", port)
		if err := srvs[0].ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.Error("Http服务运行失败：%s。", err.Error())
		}
	}()
}

// var httpProxys = MakeBucket("httpProxys")

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	// console.Info("%s", utils.JsonMarshal(req.Header))
	// u, err := url.Parse("addr")
	// if err != nil {
	// 	core.Logs.Warn("can't connect to the http proxy:", err)
	// 	return
	// }
	// Transport = &http.Transport{Proxy: http.ProxyURL(u)}
	// resp, err := Transport.RoundTrip(req)
	// fmt.Println("====")
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
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
