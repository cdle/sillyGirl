package qinglong

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
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
	Token        string `json:"-"`
	Error        error  `json:"-"`
	Default      bool   `json:"default"`
	sync.RWMutex
	idSqlite bool
	Name     string `json:"name"`
}

// var Config *QingLong
var qinglong = core.NewBucket("qinglong")
var QLS = []*QingLong{}

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
	go func() {
		if !qinglong.GetBool("enable_qinglong", true) {
			return
		}
		initConfig()
		initTask()
		initEnv()
		initqls()
		initCron()
	}()
}

func initqls() {
	s := qinglong.Get("QLS")
	json.Unmarshal([]byte(s), &QLS)
	if len(QLS) == 0 {
		Config := &QingLong{}
		Config.Host = regexp.MustCompile(`^(https?://[\.\w]+:?\d*)`).FindString(qinglong.Get("host"))
		Config.ClientID = qinglong.Get("client_id")
		Config.ClientSecret = qinglong.Get("client_secret")
		if Config.Host != "" {
			QLS = append(QLS, Config)
			d, _ := json.Marshal(QLS)
			qinglong.Set("QLS", string(d))
		}
	}
	for _, ql := range QLS {
		if ql.Name == "" {
			ql.Name = ql.Host
		}
		if ql.Host == "" {

		}
		_, err := ql.GetToken()
		if err == nil {
			logs.Info("青龙面板(%s)连接成功。", ql.Name)
		} else {
			logs.Warn("青龙面板(%s)连接错误，%v", ql.Name, err)
		}
	}
	logs.Info("青龙360安全卫士为您保驾护航，杜绝一切流氓脚本！")
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

func Req(p interface{}, ps ...interface{}) (*QingLong, error) {
	if len(QLS) == 0 {
		return nil, errors.New("未配置容器。")
	}
	var s core.Sender
	var ql *QingLong
	var qls []*QingLong
	switch p.(type) {
	case core.Sender:
		s = p.(core.Sender)
	case *QingLong:
		ql = p.(*QingLong)
	case []*QingLong:
		qls = p.([]*QingLong)
	}
	if qls != nil {
		for i := range qls {
			Req(qls[i], ps...)
		}
		return nil, nil
	}

	if ql == nil {
		if len(QLS) > 1 {
			if s != nil {

				ls := []string{}
				for i := range QLS {
					ls = append(ls, fmt.Sprintf("%d. %s", i+1, QLS[i].Name))
				}
				s.Reply("请选择容器：\n" + strings.Join(ls, "\n"))
				r := s.Await(s, func(s core.Sender) interface{} {
					return core.Range([]int{1, len(QLS)})
				}, time.Second*10)
				switch r {
				case nil:
				default:
					ql = QLS[r.(int)-1]
				}
			}
		} else {
			ql = QLS[0]
		}
	}

	if ql == nil {
		for i := range QLS {
			if QLS[i].Default {
				ql = QLS[i]
				if s != nil {
					s.Reply(fmt.Sprintf("已默认选择容器%s", ql.Name))
				}
				break
			}
		}
	}

	if ql == nil {
		return nil, errors.New("未选择容器。")
	}

	ql.RLock()
	defer ql.RUnlock()
	token, err := ql.GetToken()
	if err != nil {
		return nil, err
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
		if ql.idSqlite {
			s := string(body)
			for _, v := range regexp.MustCompile(`"_id":"(\d+)",`).FindAllStringSubmatch(s, -1) {
				s = strings.Replace(s, v[0], `"id":`+v[1]+`,`, -1)
			}
			body = []byte(s)
			// body = []byte(strings.ReplaceAll(string(body), `"_id"`, `"id"`))
		}
		req.Body(body)
	}
	// logs.Info(ql.idSqlite, string(body))
	data, err := req.Bytes()
	if err != nil {
		return nil, err
	}
	if strings.Contains(string(data), `"id"`) {
		s := string(data)
		for _, v := range regexp.MustCompile(`"id":(\d+),`).FindAllStringSubmatch(s, -1) {
			s = strings.Replace(s, v[0], `"_id":"`+v[1]+`",`, -1)
		}
		data = []byte(s)
		if !ql.idSqlite {
			go func() {
				ql.Lock()
				ql.idSqlite = true
				ql.Unlock()
			}()
		}
	}
	// logs.Info(ql.idSqlite, string(data))
	code, _ := jsonparser.GetInt(data, "code")
	if code != 200 {
		return nil, errors.New(string(data))
	}
	if toParse != nil {
		if err := json.Unmarshal(data, toParse); err != nil {
			return nil, err
		}
	}
	if get != nil {
		if *get, err = jsonparser.GetString(data, *get); err != nil {
			return nil, err
		}
	}
	if c != nil {
		c.Value, _ = jsonparser.GetString(data, strings.Split(c.Get, ".")...)
	}
	return ql, nil
}

func QinglongSC(s core.Sender) (error, []*QingLong) {
	if len(QLS) == 0 {
		return errors.New("未配置容器。"), nil
	}
	ls := []string{}
	for i := range QLS {
		ls = append(ls, fmt.Sprintf("%d. %s", i+1, QLS[i].Name))
	}
	s.Reply("请选择容器：\n" + strings.Join(ls, "\n"))
	r := s.Await(s, func(s core.Sender) interface{} {
		return core.Range([]int{1, len(QLS)})
	}, time.Second*10)
	switch r {
	case nil:
		s.Reply()
		return errors.New("你没有选择容器。"), []*QingLong{}
	default:
		index := r.(int) - 1
		if index != len(QLS) {
			return nil, []*QingLong{QLS[index]}
		} else {
			return nil, QLS
		}
	}
}
