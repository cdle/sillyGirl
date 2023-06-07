package core

import (
	"fmt"
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
	sillyGirl = MakeBucket("sillyGirl")
	_, err := os.Stat(DataHome)
	if err != nil {
		os.MkdirAll(DataHome, os.ModePerm)
	}
	// utils.ReadYaml(utils.ExecPath+"/conf/", &Config, "https://raw.githubusercontent.com/cdle/sillyGirl/main/conf/demo_config.yaml")
	initToHandleMessage()
	sillyGirl.Set("compiled_at", compiled_at)
	console.Log("编译版本：%s", compiled_at)
	sillyGirl.Set("started_at", time.Now().Format("2006-01-02 15:04:05"))
	storage.Watch(sillyGirl, "compiled_at", func(old, new, key string) *storage.Final {
		if old != new {
			console.Debug("正在从 cdle/binary 获取版本号...")
			data, err := httplib.Get("https://raw.githubusercontent.com/cdle/binary/main/compile_time.go").Bytes()
			if err != nil {
				console.Error("获取版本号错误：%s", err)
				return &storage.Final{
					Error: fmt.Errorf("貌似网络不太行啊：%s", err),
				}
			}
			latest_version := regexp.MustCompile(`\d{13}`).FindString(string(data))
			if latest_version <= compiled_at {
				console.Debug("当前版本 %s 已是最新，无需升级", compiled_at)
				return &storage.Final{
					Message: fmt.Sprintf("当前版本 %s 已是最新，无需升级", compiled_at),
				}
			}
			console.Debug("正在从 cdle/binary 获取最新版本 %s 编译文件...", latest_version)
			qurl := "https://raw.githubusercontent.com/cdle/binary/master/sillyGirl_linux_" + runtime.GOARCH + "_" + latest_version
			if runtime.GOARCH == "windows" {
				qurl += ".exe"
			}
			req := httplib.Get(qurl)
			req.SetTimeout(time.Minute*5, time.Minute*5)
			data, err = req.Bytes()
			if err != nil {
				console.Error("获取最新编译文件错误：%s", err)
				return &storage.Final{
					Error: fmt.Errorf("升级时貌似网络不太行啊：%v", err),
				}
			}
			if len(data) < 2646140 {
				console.Error("获取最新编译文件错误：%v", len(data))
				return &storage.Final{
					Error: fmt.Errorf("升级时貌似网络不太行啊！%v", len(data)),
				}
			}
			console.Debug("正在创建编译文件...")
			filename := utils.ExecPath + "/" + utils.ProcessName
			ready := strings.Replace(filename, ".exe", ".ready.exe", -1)
			if f, err := os.OpenFile(ready, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777); err != nil {
				console.Error("创建编译文件错误：%v", err)
				return &storage.Final{
					Error: fmt.Errorf("创建编译文件错误：%v", err),
				}
			} else {
				_, err := f.Write(data)
				f.Close()
				if err != nil {
					des := err.Error()
					if err = os.WriteFile(ready, data, 0777); err != nil {
						console.Error("写入编译文件错误：%s || %s", des, err)
						return &storage.Final{
							Error: fmt.Errorf("写入编译文件错误：%s || %s", des, err),
						}
					}
				}
			}
			if runtime.GOOS == "window" {
				utils.Daemon("ready")
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
