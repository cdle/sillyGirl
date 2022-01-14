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
	Number   int    `json:"-"`
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
		core.AddCommand("", []core.Function{
			{
				Rules: []string{"青龙管理"},
				Admin: true,
				Handle: func(s core.Sender) interface{} {
					var ql *QingLong
					var ls []string
					nn := []*QingLong{}
					sss := qinglong.Get("QLS")
					json.Unmarshal([]byte(sss), &nn)
					t := ""
				hh:
					ls = []string{}
					for i := range nn {
						t := ""
						if nn[i].Default {
							t = "- 默认"
						}
						ls = append(ls, fmt.Sprintf("%d. %s %s", i+1, nn[i].Name, t))
					}
					s.Reply("请选择容器进行编辑：(-删除，0增加，q退出)\n" + strings.Join(ls, "\n"))
					r := s.Await(s, nil)
					is := r.(string)
					i := 0
					if is == "q" {
						goto stop
					}
					if is == "0" {
						ql = &QingLong{}
						nn = append(nn, ql)
					}
					i = core.Int(is)
					if i < 0 && i >= -len(QLS) {
						for j := range nn {
							if j == -i-1 {
								nn = append(nn[:j], nn[j+1:]...)
								break
							}
						}
						goto hh
					}
					if i > 0 && i <= len(QLS) {
						ql = nn[i-1]
					}
					if ql == nil {
						goto hh
					}
					if ql.Host == "" {
					oo:
						s.Reply("请输入青龙面板地址：")
						ql.Host = regexp.MustCompile(`^(https?://[\.\w]+:?\d*)`).FindString(s.Await(s, nil).(string))
						if ql.Host == "" {
							goto oo
						}
					}
					if ql.ClientID == "" {
						s.Reply("请输入ClientID：")
						ql.ClientID = s.Await(s, nil).(string)
					}
					if ql.ClientSecret == "" {
						s.Reply("请输入ClientSecret：")
						ql.ClientSecret = s.Await(s, nil).(string)
					}
					if ql.Name == "" {
						s.Reply("请输入备注：")
						ql.Name = s.Await(s, nil).(string)
					}
					for {
						if ql.Default {
							t = "取消默认"
						} else {
							t = "设置默认"
						}
						s.Reply(fmt.Sprintf("请选择要编辑的属性(q退出)：\n%s", strings.Join(
							[]string{
								fmt.Sprintf("1. 容器备注 - %s", ql.Name),
								fmt.Sprintf("2. 面板地址 - %s", ql.Host),
								fmt.Sprintf("3. ClientID - %s", ql.ClientID),
								fmt.Sprintf("4. ClientSecret - %s", ql.ClientSecret),
								fmt.Sprintf("5. %s", t),
							}, "\n")))
						switch s.Await(s, nil) {
						default:
							goto hh
						case "1":
							s.Reply("请输入备注：")
							ql.Name = s.Await(s, nil).(string)
						case "2":
						oo1:
							s.Reply("请输入青龙面板地址：")
							ql.Host = regexp.MustCompile(`^(https?://[\.\w]+:?\d*)`).FindString(s.Await(s, nil).(string))
							if ql.Host == "" {
								goto oo1
							}
						case "3":
							s.Reply("请输入ClientID：")
							ql.ClientID = s.Await(s, nil).(string)
						case "4":
							s.Reply("请输入ClientSecret：")
							ql.ClientSecret = s.Await(s, nil).(string)
						case "5":
							ql.Default = !ql.Default
						case "q":
							goto hh
						}
					}
				stop:
					s.Reply("是否保存修改？(Y/n)")
					if s.Await(s, func(s core.Sender) interface{} {
						return core.YesNo
					}) == core.Yes {
						QLS = nn
						d, _ := json.Marshal(nn)
						qinglong.Set("QLS", string(d))
						return "已保存修改。"
					}
					return "未作修改。"
				},
			},
		})
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

func (ql *QingLong) SetNumber(i int) {
	ql.Lock()
	defer ql.Unlock()
	ql.Number = i
}

func (ql *QingLong) GetNumber() int {
	ql.RLock()
	defer ql.RUnlock()
	return ql.Number
}

func (ql *QingLong) SetClientID(i string) {
	ql.Lock()
	defer ql.Unlock()
	ql.ClientID = i
}

func (ql *QingLong) GetClientID() string {
	ql.RLock()
	defer ql.RUnlock()
	return ql.ClientID
}

func (ql *QingLong) SetClientSecret(i string) {
	ql.Lock()
	defer ql.Unlock()
	ql.ClientSecret = i
}

func (ql *QingLong) GetClientSecret() string {
	ql.RLock()
	defer ql.RUnlock()
	return ql.ClientSecret
}

func (ql *QingLong) SetHost(i string) {
	ql.Lock()
	defer ql.Unlock()
	ql.Host = i
}

func (ql *QingLong) GetHost() string {
	ql.RLock()
	defer ql.RUnlock()
	return ql.Host
}

func (ql *QingLong) SetName(i string) {
	ql.Lock()
	defer ql.Unlock()
	ql.Name = i
}

func (ql *QingLong) SetIsSqlite() {
	ql.Lock()
	defer ql.Unlock()
	ql.idSqlite = true
}

func (ql *QingLong) IsSqlite() bool {
	ql.RLock()
	defer ql.RUnlock()
	return ql.idSqlite
}

func (ql *QingLong) GetName() string {
	ql.RLock()
	defer ql.RUnlock()
	return ql.Name
}

func (ql *QingLong) SetToken(i string) {
	ql.Lock()
	defer ql.Unlock()
	ql.Name = i
}

func (ql *QingLong) GetToken() (string, error) {
	ql.RLock()
	defer ql.RUnlock()
	if ql.Token != "" && expiration > time.Now().Unix() {
		return ql.Token, nil
	}
	req := httplib.Get(fmt.Sprintf("%s/open/auth/token?client_id=%s&client_secret=%s", ql.Host, ql.ClientID, ql.ClientSecret))
	data, err := req.Bytes()
	if err != nil {
		msg := fmt.Sprintf("青龙连接失败：%v", err)
		// logs.Warn(msg)
		return "", errors.New(msg)
	}
	code, _ := jsonparser.GetInt(data, "code")
	if code != 200 {
		msg := fmt.Sprintf("青龙登录失败：%v", string(data))
		// logs.Warn(msg)
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
	if s != nil && !s.IsAdmin() { //普通用户自动分配
		for i := range QLS {
			if QLS[i].Default {
				ql = QLS[i]
				break
			}
		}
		if ql == nil {
			ql = QLS[0]
		}
	}
	if ql == nil {
		if len(QLS) > 1 {
			if s != nil {
				ls := []string{}
				for i := range QLS {
					ls = append(ls, fmt.Sprintf("%d. %s", i+1, QLS[i].Name))
				}
				ls = append(ls, fmt.Sprintf("%d. %s", len(QLS)+1, "所有容器"))
				s.Reply("请选择容器：\n" + strings.Join(ls, "\n"))
				r := s.Await(s, func(s core.Sender) interface{} {
					return core.Range([]int{1, len(QLS) + 1})
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
		req = httplib.Get(ql.GetHost() + "/open/" + api)
	case POST:
		req = httplib.Post(ql.GetHost() + "/open/" + api)
	case DELETE:
		req = httplib.Delete(ql.GetHost() + "/open/" + api)
	case PUT:
		req = httplib.Put(ql.GetHost() + "/open/" + api)
	}
	req.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header("Content-Type", "application/json;charset=UTF-8")
	req.SetTimeout(time.Second*5, time.Second*5)
	if method != GET {
		if ql.IsSqlite() {
			s := string(body)
			for _, v := range regexp.MustCompile(`"_id":"(\d+)",`).FindAllStringSubmatch(s, -1) {
				s = strings.Replace(s, v[0], `"id":`+v[1]+`,`, -1)
			}
			body = []byte(s)
		}
		req.Body(body)
	}
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
		if !ql.IsSqlite() {
			ql.SetIsSqlite()
		}
	}

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

func GetQinglongByClientID(s string) (error, *QingLong) {
	for i := range QLS {
		if QLS[i].ClientID == s {
			return nil, QLS[i]
		}
	}
	if len(QLS) == 0 {
		return errors.New("未配置容器。"), nil
	}
	var ql *QingLong
	min := 10000000
	for i := range QLS {
		if num := QLS[i].GetNumber(); num <= min {
			min = num
			ql = QLS[i]
		}
	}
	return errors.New("默认获取了一个容器。"), ql
}

func QinglongSC(s core.Sender) (error, []*QingLong) {
	if len(QLS) == 0 {
		return errors.New("未配置容器。"), nil
	}
	if len(QLS) == 1 {
		return nil, QLS
	}
	var ql *QingLong
	if s != nil && !s.IsAdmin() { //普通用户自动分配
		for i := range QLS {
			if QLS[i].Default {
				ql = QLS[i]
				break
			}
		}
		if ql == nil {
			ql = QLS[0]
		}
	}
	if ql != nil {
		return nil, []*QingLong{ql}
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
