package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/cdle/sillyGirl/core/storage"
	"github.com/cdle/sillyGirl/utils"
)

var DataHome = utils.GetDataHome()

func Init() {
	initLoc()
	sillyGirl = MakeBucket("sillyGirl")
	_, err := os.Stat(DataHome)
	if err != nil {
		os.MkdirAll(DataHome, os.ModePerm)
	}
	// utils.ReadYaml(utils.ExecPath+"/conf/", &Config, "https://raw.githubusercontent.com/cdle/sillyGirl/main/conf/demo_config.yaml")
	initToHandleMessage()
	sillyGirl.Set("compiled_at", compiled_at)
	console.Log("编译版本：%s", compiled_at)
	initWeb()
	initCarry()
	sillyGirl.Set("started_at", time.Now().Format("2006-01-02 15:04:05"))
	storage.Watch(sillyGirl, "compiled_at", func(old, new, key string) *storage.Final {
		if old != new {
			var transport *http.Transport
			instance, err := GetProxyTransport("https://raw.githubusercontent.com", "", nil)
			if err != nil {
				console.Error("升级代理错误：%s", err)
				return &storage.Final{
					Error: fmt.Errorf("升级代理错误：：%s", err),
				}
			}
			if instance != nil {
				defer instance.Close()
			}
			if instance != nil {
				transport = &http.Transport{
					Dial: func(string, string) (net.Conn, error) {
						return instance, nil
					},
					MaxIdleConns:          100,
					IdleConnTimeout:       90 * time.Second,
					TLSHandshakeTimeout:   10 * time.Second,
					ExpectContinueTimeout: 1 * time.Second,
				}
			}
			var client = &http.Client{}
			if transport != nil {
				client.Transport = transport
			}
			var body io.Reader
			var data []byte
			var latest_version = ""

			console.Debug("正在从 cdle/binary 获取版本号...")
			qurl := "https://raw.githubusercontent.com/cdle/binary/main/compile_time.go"
			req, _ := http.NewRequest("GET", qurl, strings.NewReader(""))
			resp, err := client.Do(req)
			if err != nil {
				console.Error("获取版本号错误：%s", err)
				// return &storage.Final{
				// 	Error: fmt.Errorf("貌似网络不太行啊：%s", err),
				// }
				goto PROXY
			}
			defer resp.Body.Close()
			data, _ = ioutil.ReadAll(resp.Body)
			latest_version = regexp.MustCompile(`\d{13}`).FindString(string(data))
			if latest_version <= compiled_at {
				console.Debug("当前版本 %s 已是最新，无需升级", compiled_at)
				return &storage.Final{
					Message: fmt.Sprintf("当前版本 %s 已是最新，无需升级", compiled_at),
				}
			}
			client = &http.Client{}
			if transport != nil {
				client.Transport = transport
			}
			console.Debug("正在从 cdle/binary 获取最新版本 %s 编译文件...", latest_version)
			qurl = "https://raw.githubusercontent.com/cdle/binary/master/sillyGirl_" + runtime.GOOS + "_" + runtime.GOARCH + "_" + latest_version
			if runtime.GOOS == "windows" {
				qurl += ".exe"
			}
			req, _ = http.NewRequest("GET", qurl, strings.NewReader(""))
			resp, err = client.Do(req)
			if err != nil {
				console.Error("获取最新编译文件错误：%s", err)
				// return &storage.Final{
				// 	Error: fmt.Errorf("升级时貌似网络不太行啊：%v", err),
				// }
				goto PROXY
			}
			defer resp.Body.Close()
			body = resp.Body
			goto CREATE

		PROXY:
			//使用免费代理下载
			console.Info("正在重新尝试下载...")
			qurl = "http://127.0.0.1:8765/api/download?version=" + compiled_at + "&goos=" + runtime.GOOS + "&goarch=" + runtime.GOARCH
			resp, err = http.Get(qurl)
			if err != nil {
				return &storage.Final{
					Error: fmt.Errorf("升级时貌似网络不太行啊"),
				}
			}
			defer resp.Body.Close()
			body = resp.Body
			switch resp.Header.Get("Result") {
			case "newest":
				return &storage.Final{
					Message: fmt.Sprintf("当前版本 %s 已是最新，无需升级", compiled_at),
				}
			case "fail":
				return &storage.Final{
					Error: fmt.Errorf("升级失败"),
				}
			case "ok":
			}

		CREATE:
			console.Debug("正在创建编译文件...")
			filename := utils.ExecPath + "/" + utils.ProcessName
			ready := ""
			if runtime.GOOS == "windows" {
				ready = strings.Replace(filename, ".exe", ".ready.exe", -1)
			} else {
				ready += ".ready"
			}
			f, err := os.OpenFile(ready, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				console.Error("创建编译文件错误：%v", err)
				return &storage.Final{
					Error: fmt.Errorf("创建编译文件错误：%v", err),
				}
			}
			defer f.Close()
			i, err := io.Copy(f, body)
			if i < 2646140 || err != nil {
				console.Error("创建编译文件错误：%v", i)
				return &storage.Final{
					Error: fmt.Errorf("创建编译文件错误：%v", i),
				}
			}
			if runtime.GOOS == "windows" {
				console.Log("正在准备重启...")
				go func() {
					time.Sleep(time.Second)
					utils.Daemon("ready")
				}()
				return nil
			} else {
				console.Debug("正在删除旧程序错误...")
				if err = os.RemoveAll(filename); err != nil {
					console.Error("删除旧程序错误：%v", err)
					return &storage.Final{
						Error: fmt.Errorf("删除旧程序错误：%v", err),
					}
				}
			}
			console.Debug("正在移动新程序错误...")
			if err = os.Rename(ready, filename); err != nil {
				console.Error("移动新程序错误：%v", err)
				return &storage.Final{
					Error: fmt.Errorf("移动新程序错误：%v", err),
				}
			}
			go func() {
				console.Debug("正在重启...")
				time.Sleep(time.Second)
				utils.Daemon()
			}()
			return &storage.Final{
				Message: "升级成功，即将重启！",
			}
		}
		return nil
	})
	storage.Watch(sillyGirl, "started_at", func(old, new, key string) *storage.Final {
		if old != new {
			go func() {
				time.Sleep(time.Second)
				utils.Daemon()
			}()
			return &storage.Final{
				Message: "马上重启！",
			}
		}
		return nil
	})

	api_key := sillyGirl.GetString("api_key")
	if api_key == "" {
		api_key := time.Now().UnixNano()
		sillyGirl.Set("api_key", api_key)
	}
	// if sillyGirl.GetString("uuid") == "" {
	sillyGirl.Set("uuid", utils.GenUUID())
	// }
	httplib.SetDefaultSetting(httplib.BeegoHTTPSettings{
		ConnectTimeout:   time.Second * 10,
		ReadWriteTimeout: time.Second * 10,
		UserAgent:        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36",
	})
	initPlugins()
	initReboot()
	initListenReply()
	// initPluginFile()
	initWebPluginList()
	go initPluginList()
	initPluginPublish()

}
