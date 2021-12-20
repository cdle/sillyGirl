package qinglong

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/buger/jsonparser"
	"github.com/cdle/sillyGirl/core"
)

type QingLong struct {
	Host         string `json:"host"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Token        string
}

var Config *QingLong
var qinglong = core.NewBucket("qinglong")

var expiration int64
var GET = "GET"
var PUT = "PUT"
var POST = "POST"
var DELETE = "DELETE"
var ENVS = "envs"
var CRONS = "crons"
var CONFIG = "configs"

type Carrier struct {
	Get   string
	Value string
}

func init() {
	if !qinglong.GetBool("enable_qinglong", true) {
		return
	}
	initConfig()
	initTask()
	initEnv()
	Config = &QingLong{}
	Config.Host = qinglong.Get("host", "http://127.0.0.1:5700")
	Config.ClientID = qinglong.Get("client_id")
	Config.ClientSecret = qinglong.Get("client_secret")
	if Config.Host == "" {
		return
	}
	if v := regexp.MustCompile(`^(https?://[\.\w]+:?\d*)`).FindStringSubmatch(Config.Host); len(v) == 2 {
		Config.Host = v[1]
	}
	_, err := Config.GetToken()
	if err == nil {
		logs.Info("青龙面板连接成功。")
	}
	initCron()

}

func (ql *QingLong) GetToken() (string, error) {
	if ql.Token != "" && expiration > time.Now().Unix() {
		return ql.Token, nil
	}
	req := httplib.Get(fmt.Sprintf("%s/open/auth/token?client_id=%s&client_secret=%s", ql.Host, ql.ClientID, ql.ClientSecret))
	data, err := req.Bytes()
	if err != nil {
		msg := fmt.Sprintf("青龙连接失败：%v", err)
		logs.Warn(msg)
		return "", errors.New(msg)
	}
	code, _ := jsonparser.GetInt(data, "code")
	if code != 200 {
		msg := fmt.Sprintf("青龙登录失败：%v", string(data))
		logs.Warn(msg)
		return "", errors.New(msg)
	}
	ql.Token, _ = jsonparser.GetString(data, "data", "token")
	expiration, _ = jsonparser.GetInt(data, "data", "expiration")
	return ql.Token, nil
}

func (ql *QingLong) Req(ps ...interface{}) error {
	if ql.Host == "" {
		return nil
	}
	token, err := ql.GetToken()
	if err != nil {
		return err
	}
	method := GET
	body := []byte{}
	api := ENVS
	apd := ""
	var get *string
	var c *Carrier
	var toParse interface{}
	for _, p := range ps {
		switch p.(type) {
		case string:
			switch p.(string) {
			case GET, POST, DELETE, PUT:
				method = p.(string)
			case ENVS, CRONS, CONFIG:
				api = p.(string)
			default:
				apd = p.(string)
			}
		case []byte:
			body = p.([]byte)
		case *Carrier:
			c = p.(*Carrier)
		case *string:
			get = p.(*string)
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
	api = strings.Trim(api, " ")
	switch method {
	case GET:
		req = httplib.Get(ql.Host + "/open/" + api)
	case POST:
		req = httplib.Post(ql.Host + "/open/" + api)
	case DELETE:
		req = httplib.Delete(ql.Host + "/open/" + api)
	case PUT:
		req = httplib.Put(ql.Host + "/open/" + api)
	}
	req.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header("Content-Type", "application/json;charset=UTF-8")
	req.SetTimeout(time.Second*5, time.Second*5)
	if method != GET {
		req.Body(body)
	}
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
			return err
		}
	}
	if get != nil {
		if *get, err = jsonparser.GetString(data, *get); err != nil {
			return err
		}
	}
	if c != nil {
		c.Value, _ = jsonparser.GetString(data, strings.Split(c.Get, ".")...)
	}
	return nil
}
