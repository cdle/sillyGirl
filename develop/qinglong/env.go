package qinglong

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/buger/jsonparser"
)

type EnvResponse struct {
	Code int   `json:"code"`
	Data []Env `json:"data"`
}

type Env struct {
	Value   string `json:"value"`
	ID      string `json:"_id"`
	Status  int    `json:"status"`
	Name    string `json:"name"`
	Remarks string `json:"remarks"`
}

func GetEnv(searchValue string) ([]Env, error) {
	er := EnvResponse{}
	if err := req(PUT, ENVS, &er, "?searchValue="+searchValue); err != nil {
		return nil, err
	}
	return er.Data, nil
}

func SetEnv(e *Env) error {
	es, err := GetEnv(e.Name)
	if err != nil {
		return err
	}
	if len(es) == 0 {
		return AddEnv(e)
	}
	e.ID = es[0].ID
	return req(PUT, ENVS, e)
}

func AddEnv(e *Env) error {
	return req(POST, ENVS, []Env{*e})
}

func RemEnv(e *Env) error {
	es, err := GetEnv(e.Name)
	if err != nil {
		return err
	}
	if len(es) == 0 {
		return nil
	}
	td := []string{}
	for _, e := range es {
		td = append(td, e.ID)
	}
	return req(DELETE, ENVS, td)
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
				method = p.(string)
			default:
				apd = p.(string)
			}
		case []byte:
			body = p.([]byte)
		default:
			if strings.Contains(reflect.TypeOf(&p).String(), "*") {
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
		req = httplib.Delete(Config.Host + "/open/" + api)
	case DELETE:
		req = httplib.Delete(Config.Host + "/open/" + api)
	case PUT:
		req = httplib.Put(Config.Host + "/open/" + api)
	}
	if method != GET && len(body) > 0 {
		req.Body(body)
	}
	req.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	data, _ := json.Marshal(body)
	req.Body(data)
	data, err = req.Bytes()
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
