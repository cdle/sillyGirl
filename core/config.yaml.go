package core

import (
	"io/ioutil"
	"os"

	"github.com/astaxie/beego/logs"
	"gopkg.in/yaml.v2"
)

type Yaml struct {
	Replies []Reply
}

var Config Yaml

func ReadYaml(confDir string, conf interface{}, _ string) {
	path := confDir + "config.yaml"
	if _, err := os.Stat(confDir); err != nil {
		os.MkdirAll(confDir, os.ModePerm)
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return
	}
	s, _ := ioutil.ReadAll(f)
	if len(s) == 0 {
		return
	}
	f.Close()
	content, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Warn("解析配置文件%s读取错误: %v", path, err)
		return
	}
	if yaml.Unmarshal(content, conf) != nil {
		logs.Warn("解析配置文件%s出错: %v", path, err)
		return
	}
	logs.Info("解析配置文件%s", path)
}
