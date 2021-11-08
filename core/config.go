package core

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
	"gopkg.in/yaml.v2"
)

type Yaml struct {
	Replies []Reply
}

var ExecPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

var Config Yaml

func ReadYaml(confDir string, conf interface{}, _ string) {
	path := confDir + "config.yaml"
	if _, err := os.Stat(confDir); err != nil {
		os.MkdirAll(confDir, os.ModePerm)
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		// logs.Warn(err)
		return
	}
	s, _ := ioutil.ReadAll(f)
	if len(s) == 0 {
		// logs.Info("下载配置%s", url)
		// r, err := httplib.Get("https://ghproxy.com/" + url).Response()//
		// if err == nil {
		// 	io.Copy(f, r.Body)
		// }
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
