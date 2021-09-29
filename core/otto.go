package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/robertkrimen/otto"
)

type JsReply string

func init() {
	files, err := ioutil.ReadDir(ExecPath + "/develop/replies")
	if err != nil {
		logs.Warn("打开文件夹%s错误，%v", ExecPath+"/develop/replies", err)
		return
	}
	for _, v := range files {
		if v.IsDir() {
			continue
		}
		jr := string(ExecPath + "/develop/replies/" + v.Name())
		data := ""
		if strings.Contains(jr, "http") {
			data, err = httplib.Get(jr).String()
			if err != nil {
				logs.Warn("回复：%s获取失败%v", jr, err)
				continue
			}
		} else {
			f, err := os.Open(jr)
			if err != nil {
				logs.Warn("回复：%s打开失败%v", jr, err)
				continue
			}
			v, _ := ioutil.ReadAll(f)
			data = string(v)
		}
		rules := []string{}
		for _, v := range regexp.MustCompile(`\[rule:\s*([^\s\[\]]+)\s*\]`).FindAllStringSubmatch(data, -1) {
			rules = append(rules, v[1])
		}
		if len(rules) == 0 {
			logs.Warn("回复：%s找不到规则", jr, err)
			continue
		}
		request := func(call otto.Value) interface{} {
			url := ""
			dataType := ""
			method := "get"
			body := ""
			{
				v, _ := call.Object().Get("url")
				url = v.String()
			}
			{
				v, _ := call.Object().Get("dataType")
				dataType = v.String()
			}
			{
				v, _ := call.Object().Get("body")
				body = v.String()
			}
			var req *httplib.BeegoHTTPRequest
			switch strings.ToLower(method) {
			case "delete":
				req = httplib.Delete(url)
			case "post":
				req = httplib.Post(url)
			case "put":
				req = httplib.Put(url)
			default:
				req = httplib.Get(url)
			}
			if body != "" {
				req.Body(body)
			}
			data, err := req.String()
			if err != nil {
				return nil
			}

			if strings.Contains(dataType, "json") {
				obj, _ := otto.New().Object(fmt.Sprintf(`(%s)`, data))
				return obj
			}
			result, _ := otto.ToValue(data)
			return result
		}
		var handler = func(s Sender) interface{} {
			template := data
			for k, v := range s.GetMatch() {
				template = strings.Replace(template, fmt.Sprintf(`param(%d)`, k+1), fmt.Sprintf(`"%s"`, v), -1)
			}
			vm := otto.New()
			vm.Set("request", request)
			vm.Set("sendText", func(call otto.Value) interface{} {
				s.Reply(call.String())
				return nil
			})
			vm.Set("sendImage", func(call otto.Value) interface{} {
				s.Reply(ImageUrl(call.String()))
				return nil
			})
			rt, err := vm.Run(template + `
""
`)
			result := rt.String()
			if err != nil {
				return err
			}
			for _, v := range regexp.MustCompile(`\[image:\s*([^\s\[\]]+)\s*\]`).FindAllStringSubmatch(result, -1) {
				s.Reply(ImageUrl(v[1]))
				result = strings.Replace(result, fmt.Sprintf(`[image:%s]`, v[1]), "", -1)
			}
			if result == "" {
				return nil
			}
			return result
		}
		logs.Warn("回复：%s添加成功", jr)
		functions = append(functions, Function{
			Handle: handler,
			Rules:  rules,
		})
	}
}
