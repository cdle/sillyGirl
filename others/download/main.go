package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var version = ""

func init() {
	resp, err := http.Get("https://raw.githubusercontent.com/cdle/binary/main/compile_time.go")
	if err != nil {
		panic(fmt.Errorf("获取版本号错误：%s", err))
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("获取版本号错误：%s", err))
	}
	version = regexp.MustCompile(`\d{13}`).FindString(string(data))

	if version == "" {
		panic(fmt.Errorf("获取版本号错误：%v", version))
	}
	fmt.Println("获取版本号：", version)
}

var Result = "Result"

type Class struct {
	Os   string
	Arch string
	Data []byte
	sync.RWMutex
}

var classes = []*Class{
	{
		Os:   "linux",
		Arch: "amd64",
	},
	{
		Os:   "linux",
		Arch: "arm64",
	},
	{
		Os:   "windows",
		Arch: "amd64",
	},
}

func main() {
	router := gin.Default()
	router.GET("/api/download", func(c *gin.Context) {
		if version <= c.Query("version") {
			c.Header(Result, "newest")
			c.String(200, "版本已是最新，无需升级！")
			return
		}
		var goos = c.Query("goos")
		var goarch = c.Query("goarch")
		for _, class := range classes {
			if goos == class.Os && goarch == class.Arch {
				func() {
					class.Lock()
					defer class.Unlock()
					if class.Data == nil {
						fmt.Println("没有就下载", version)
						qurl := "https://raw.githubusercontent.com/cdle/binary/master/sillyGirl_" + goos + "_" + goarch + "_" + version
						if goos == "windows" {
							qurl += ".exe"
						}
						fmt.Println(qurl)
						req, _ := http.NewRequest("GET", qurl, strings.NewReader(""))
						resp, err := (&http.Client{}).Do(req)
						if err != nil {
							c.Header(Result, "fail")
							c.String(200, err.Error())
							return
						}
						if resp.StatusCode != 200 {
							c.Header(Result, "fail")
							c.String(200, "下载失败")
							return
						}
						data, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							c.Header(Result, "fail")
							c.String(200, err.Error())
							return
						}
						class.Data = data
					}
				}()
				c.Header(Result, "ok")
				c.Data(200, "application/octet-stream", class.Data)
			}
		}
		c.Header(Result, "fail")
		c.String(200, "")
	})
	//http://127.0.0.1:8765/api/version?version=
	router.GET("/api/version", func(c *gin.Context) {
		v := c.Query("version")
		if v != "" {
			for _, class := range classes {
				class.Lock()
				class.Data = nil
				class.Unlock()
			}
			version = v
			fmt.Println("更新版本号：", version)
		}
		c.String(200, version)
	})

	router.Run(":8765")
}
