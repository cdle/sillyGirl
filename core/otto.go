package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/dop251/goja"
)

type JsReply string

func init() {

	files, err := ioutil.ReadDir(ExecPath + "/develop/replies")
	if err != nil {
		logs.Warn("打开文件夹%s错误，%v", ExecPath+"/develop/replies", err)
		return
	}
	var o = NewBucket("otto")
	get := func(call goja.FunctionCall) string {
		key := call.Argument(0).String()
		value := call.Argument(1).String()
		return o.Get(key, value)
	}
	set := func(call goja.FunctionCall) {
		key := call.Argument(0).String()
		value := call.Argument(1).String()
		o.Set(key, value)
	}
	push := func(call goja.Value) {
		imType := call.ToObject(nil).Get("imType").String()
		groupCode := call.ToObject(nil).Get("groupCode").ToInteger()
		userID := call.ToObject(nil).Get("userID").ToInteger()
		content := call.ToObject(nil).Get("content").String()
		if groupCode != 0 {
			if push, ok := GroupPushs[imType]; ok {
				push(int(groupCode), int(userID), content)
			}
		} else {
			if push, ok := Pushs[imType]; ok {
				push(int(userID), content)
			}
		}
	}
	request := func(call goja.Value) interface{} {
		url := ""
		dataType := ""
		method := "get"
		body := ""
		{
			url = call.ToObject(nil).Get("url").String()
		}
		{
			dataType = call.ToObject(nil).Get("dataType").String()
		}
		{
			v := call.ToObject(nil).Get("body").String()
			body = v
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
			return goja.Undefined()
		}
		if strings.Contains(dataType, "json") {
			// obj, err := goja.New().Object(fmt.Sprintf(`(%s)`, data))
			// if err != nil {
			// 	return goja.Undefined()
			// }
			// return obj
		}
		return data
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
		for _, res := range regexp.MustCompile(`\[rule:([^\[\]]+)\]`).FindAllStringSubmatch(data, -1) {
			rules = append(rules, strings.Trim(res[1], " "))
		}
		cron := ""
		if res := regexp.MustCompile(`\[cron:([^\[\]]+)\]`).FindStringSubmatch(data); len(res) != 0 {
			cron = strings.Trim(res[1], " ")
		}
		admin := false
		if res := regexp.MustCompile(`\[admin:([^\[\]]+)\]`).FindStringSubmatch(data); len(res) != 0 {
			admin = strings.Trim(res[1], " ") == "true"
		}
		if len(rules) == 0 {
			logs.Warn("回复：%s找不到规则", jr, err)
			continue
		}
		var handler = func(s Sender) interface{} {
			template := data
			for k, v := range s.GetMatch() {
				template = strings.Replace(template, fmt.Sprintf(`param(%d)`, k+1), fmt.Sprintf(`"%s"`, v), -1)
			}
			vm := goja.New()
			vm.Set("set", set)
			vm.Set("get", get)
			vm.Set("request", request)
			vm.Set("push", push)
			vm.Set("sendText", func(call goja.Value) {
				s.Reply(call.String())

			})
			vm.Set("sendImage", func(call goja.Value) {
				s.Reply(ImageUrl(call.String()))
			})
			rt, err := vm.RunString(template + `
""
`)
			if err != nil {
				return err
			}
			result := rt.String()
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
		AddCommand("", []Function{
			{
				Handle: handler,
				Rules:  rules,
				Cron:   cron,
				Admin:  admin,
			},
		})
	}
}
