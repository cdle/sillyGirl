package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/robertkrimen/otto"
)

type JsReply string

func init() {
	go func() {
		time.Sleep(time.Second)
		init123()
	}()
}

func init123() {
	files, err := ioutil.ReadDir(ExecPath + "/develop/replies")
	if err != nil {
		logs.Warn("打开文件夹%s错误，%v", ExecPath+"/develop/replies", err)
		return
	}
	var o = NewBucket("otto")
	get := func(call otto.FunctionCall) (result otto.Value) {
		key := call.Argument(0).String()
		value := call.Argument(1).String()
		result, _ = otto.ToValue(o.Get(key, value))
		return
	}
	set := func(call otto.FunctionCall) interface{} {
		key := call.Argument(0).String()
		value := call.Argument(1).String()
		o.Set(key, value)
		return otto.Value{}
	}
	push := func(call otto.Value) interface{} {
		imType, _ := call.Object().Get("imType")
		groupCode, _ := call.Object().Get("groupCode")
		userID, _ := call.Object().Get("userID")
		content, _ := call.Object().Get("content")
		gid, _ := groupCode.ToInteger()
		if gid != 0 {
			if push, ok := GroupPushs[imType.String()]; ok {
				uid, _ := userID.ToInteger()
				push(int(gid), int(uid), content.String())
			}
		} else {
			if push, ok := Pushs[imType.String()]; ok {
				uid, _ := userID.ToInteger()
				push(int(uid), content.String())
			}
		}
		return otto.Value{}
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
			return otto.Value{}
		}
		if strings.Contains(dataType, "json") {
			obj, err := otto.New().Object(fmt.Sprintf(`(%s)`, data))
			if err != nil {
				return otto.Value{}
			}
			return obj
		}
		result, err := otto.ToValue(data)
		if err != nil {
			return otto.Value{}
		}
		return result
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
		for _, res := range regexp.MustCompile(`\[rule:(.+)\]`).FindAllStringSubmatch(data, -1) {
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
			template = strings.Replace(template, "ImType()", fmt.Sprintf(`"%s"`, s.GetImType()), -1)
			template = strings.Replace(template, "GetChatID()", fmt.Sprint(s.GetChatID()), -1)
			param := func(call otto.Value) otto.Value {
				i, _ := call.ToInteger()
				v, _ := otto.ToValue(s.Get(int(i - 1)))
				return v
			}
			vm := otto.New()
			vm.Set("Delete", func() {
				s.Delete()
			})
			vm.Set("Continue", func() {
				s.Continue()
			})
			vm.Set("GetUsername", func() otto.Value {
				v, _ := otto.ToValue(s.GetUsername())
				return v
			})
			vm.Set("GetUserID", func() otto.Value {
				v, _ := otto.ToValue(s.GetUserID())
				return v
			})
			vm.Set("set", set)
			vm.Set("param", param)
			vm.Set("get", get)
			vm.Set("request", request)
			vm.Set("push", push)
			vm.Set("sendText", func(call otto.Value) interface{} {
				s.Reply(call.String())
				return otto.Value{}
			})
			vm.Set("sendImage", func(call otto.Value) interface{} {
				s.Reply(ImageUrl(call.String()))
				return otto.Value{}
			})
			rt, err := vm.Run(template + `
""
`)
			if err != nil {
				return err
			}
			result := rt.String()
			for _, v := range regexp.MustCompile(`\[image:\s*([^\s\[\]]+)\s*\]`).FindAllStringSubmatch(result, -1) {
				s.Reply(ImageUrl(v[1]))
				result = strings.Replace(result, fmt.Sprintf(`[image:%s]\n`, v[1]), "", -1)
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
