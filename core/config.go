package core

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/cdle/sillyGirl/im"
	"gopkg.in/yaml.v2"
)

type Yaml struct {
	Im      []im.Config
	Replies []Reply
}

var ExecPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

var Config Yaml

func init() {
	confDir := ExecPath + "/conf"
	if _, err := os.Stat(confDir); err != nil {
		os.MkdirAll(confDir, os.ModePerm)
	}
	for _, name := range []string{"config.yaml"} {
		f, err := os.OpenFile(ExecPath+"/conf/"+name, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			logs.Warn(err)
		}
		s, _ := ioutil.ReadAll(f)
		if len(s) == 0 {
			logs.Info("下载配置%s", name)
			r, err := httplib.Get("https://ghproxy.com/" + "https://raw.githubusercontent.com/cdle/sillyGirl/main/conf/demo_" + name).Response()
			if err == nil {
				io.Copy(f, r.Body)
			}
		}
		f.Close()
	}
	content, err := ioutil.ReadFile(ExecPath + "/conf/config.yaml")
	if err != nil {
		logs.Warn("解析config.yaml读取错误: %v", err)
	}
	if yaml.Unmarshal(content, &Config) != nil {
		logs.Warn("解析config.yaml出错: %v", err)
	}

	InitReplies()
	initToHandleMessage()
}
