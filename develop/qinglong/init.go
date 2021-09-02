package qinglong

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
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
	core.ReadYaml(core.ExecPath+"/develop/qinglong/conf/", &Config, "https://raw.githubusercontent.com/cdle/sillyGirl/main/develop/qinglong/conf/demo_config.yaml")
	_, err := getToken()
	if err == nil {
		logs.Info("青龙已连接")
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

func req(ps ...interface{}) error {
	token, err := getToken()
	if err != nil {
		return err
	}
	method := GET
	body := []byte{}
	api := ENVS
	apd := ""
	var toParse interface{}
	for _, p := range ps {
		switch p.(type) {
		case string:
			switch p.(string) {
			case GET, POST, DELETE, PUT:
				method = p.(string)
			case ENVS:
				api = p.(string)
			default:
				apd = p.(string)
			}
		case []byte:
			body = p.([]byte)
		default:
			if strings.Contains(reflect.TypeOf(p).String(), "*") {
				toParse = p
			} else {
				body, _ = json.Marshal(p)
			}
		}
	}
	var req *httplib.BeegoHTTPRequest
	api += apd
	switch method {
	case GET:
		req = httplib.Get(Config.Host + "/open/" + api)
	case POST:
		req = httplib.Post(Config.Host + "/open/" + api)
	case DELETE:
		req = httplib.Delete(Config.Host + "/open/" + api)
	case PUT:
		req = httplib.Put(Config.Host + "/open/" + api)
	}
	req.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header("Content-Type", "application/json;charset=UTF-8")
	if method != GET {
		req.Body(body)
	}
	fmt.Println(Config.Host + "/open/" + api)
	fmt.Println(method)
	fmt.Println(string(body))
	req.Body(body)
	data, err := req.Bytes()
	if err != nil {
		return err
	}
	code, _ := jsonparser.GetInt(data, "code")
	if code != 200 {
		return errors.New(string(data))
	}
	if toParse != nil {
		if err := json.Unmarshal(data, toParse); err != nil {
			return errors.New(fmt.Sprintf("解析错误：%v,%v", err, string(data)))
		}
	}
	return nil
}
