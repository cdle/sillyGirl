package core

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
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
