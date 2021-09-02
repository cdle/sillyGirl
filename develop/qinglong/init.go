package qinglong

import (
	"errors"
	"fmt"
	"time"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/buger/jsonparser"
	"github.com/cdle/sillyGirl/core"
)

type Yaml struct {
	Host         string
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

var Config Yaml

var token string
var expiration int64
var GET = "GET"
var PUT = "PUT"
var POST = "POST"
var DELETE = "DELETE"
var ENVS = "envs"

func init() {
	core.ReadYaml(core.ExecPath+"/develop/qinglong/conf/", &Config, "https://raw.githubusercontent.com/cdle/sillyGirl/main/develop/qinglong/conf/config.yaml")
	token, err := getToken()
	if err == nil {
		logs.Info("青龙可及%v", token)
	}

}

func getToken() (string, error) {
	if token != "" && expiration > time.Now().Unix() {
		return token, nil
	}
	req := httplib.Get(fmt.Sprintf("%s/open/auth/token?client_id=%s&client_secret=%s", Config.Host, Config.ClientID, Config.ClientSecret))
	data, err := req.Bytes()
	if err != nil {
		msg := fmt.Sprintf("青龙链接失败：%v", err)
		logs.Warn(msg)
		return "", errors.New(msg)
	}
	code, _ := jsonparser.GetInt(data, "code")
	if code != 200 {
		msg := fmt.Sprintf("青龙登录失败：%v", string(data))
		logs.Warn(msg)
		return "", errors.New(msg)
	}
	token, _ = jsonparser.GetString(data, "data", "token")
	expiration, _ = jsonparser.GetInt(data, "data", "expiration")
	return token, nil
}
