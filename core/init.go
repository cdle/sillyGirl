package core

import (
	"os"
	"time"

	"github.com/beego/beego/v2/adapter/httplib"
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

	sillyGirl.Set("started_at", time.Now().Format("2006-01-02 15:04:05"))
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
	if compiled_at != "" {
		console.Log("编译时间戳，", compiled_at)
	}
}
