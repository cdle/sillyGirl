package core

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/denisbrodbeck/machineid"
	"github.com/robertkrimen/otto"
)

type JsReply string

var o = NewBucket("otto")

func init() {
	go func() {
		time.Sleep(time.Second)
		// if o.GetBool("enable_price", true) {
		// 	os.MkdirAll("develop/replies", os.ModePerm)
		// 	if data, err := os.ReadFile("scripts/jd_price.js"); err == nil {
		// 		os.WriteFile("develop/replies/jd_price.js", data, os.ModePerm)
		// 	}
		// 	os.Remove("develop/replies/price.js")
		// } else {
		// 	os.Remove("develop/replies/jd_price.js")
		// }
		os.Remove("develop/replies/jd_price.js")
		init123()
	}()
}

var OttoFuncs = map[string]func(string) string{
	"machineId": func(_ string) string {
		// data, _ := os.ReadFile("/var/lib/dbus/machine-id")
		// id := regexp.MustCompile(`\w+`).FindString(string(data))
		// if id == "" {
		// 	data, _ = os.ReadFile("/etc/machine-id")
		// 	id = regexp.MustCompile(`\w+`).FindString(string(data))
		// }
		id, err := machineid.ProtectedID("sillyGirl")
		if err != nil {
			id = sillyGirl.Get("machineId")
			if id == "" {
				id = GetUUID()
				sillyGirl.Set("machineId", id)
			}
		}
		return id
	},
	"uuid": func(_ string) string {
		return GetUUID()
	},
	"md5": func(str string) string {
		w := md5.New()
		io.WriteString(w, str)
		md5str := fmt.Sprintf("%x", w.Sum(nil))
		return md5str
	},
	"timeFormat": func(str string) string {
		return time.Now().Format(str)
	},
}

func init123() {
	files, err := ioutil.ReadDir(ExecPath + "/develop/replies")
	if err != nil {
		os.MkdirAll(ExecPath+"/develop/replies", os.ModePerm)
		// logs.Warn("打开文件夹%s错误，%v", "develop/replies", err)
		return
	}

	get := func(call otto.FunctionCall) (result otto.Value) {
		key := call.Argument(0).String()
		value := call.Argument(1).String()
		result, _ = otto.ToValue(o.Get(key, value))
		return
	}
	bucketGet := func(bucket otto.Value, key otto.Value) (result otto.Value) {
		result, _ = otto.ToValue(o.Get(key, Bucket(bucket.String()).Get(key.String())))
		return
	}
	bucketSet := func(bucket otto.Value, key otto.Value, value otto.Value) (result otto.Value) {
		Bucket(bucket.String()).Set(key.String(), value.String())
		return otto.Value{}
	}
	bucketKeys := func(bucket otto.Value) (result otto.Value) {
		b := Bucket(bucket.String())
		if !IsBucket(b) {
			result, _ = otto.ToValue("")
			return
		}
		rt := ""
		b.Foreach(func(k, _ []byte) error {
			rt += fmt.Sprintf("%s;", k)
			return nil
		})
		result, _ = otto.ToValue(rt)
		return
	}
	set := func(key otto.Value, value otto.Value) interface{} {
		o.Set(key.String(), value.String())
		return otto.Value{}
	}
	sleep := func(value otto.Value) interface{} {
		i, _ := value.ToInteger()
		time.Sleep(time.Duration(i) * time.Millisecond)
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
				push(int(gid), userID, content.String(), "")
			}
		} else {
			if push, ok := Pushs[imType.String()]; ok {
				push(userID, content.String(), nil, "")
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
		{
			v, _ := call.Object().Get("method")
			method = v.String()
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
		{
			v, err := call.Object().Get("headers")
			if err == nil && v.IsObject() {
				headers := v.Object()
				for _, key := range headers.Keys() {
					v, _ := headers.Get(key)
					req.Header(key, v.String())
				}
			}
		}
		if body != "" {
			if body != "" && body != "undefined" {
				req.Body(body)
				req.Header("Content-Type", "application/json")
			}
			req.Body(body)
		}
		if dataType == "location" {
			req.SetCheckRedirect(func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			})
			rsp, err := req.Response()
			if err == nil && (rsp.StatusCode == 301 || rsp.StatusCode == 302) {
				url = rsp.Header.Get("Location")
			}
			result, err := otto.ToValue(url)
			if err != nil {
				return otto.Value{}
			}
			return result
		}
		{
			v, _ := call.Object().Get("useProxy")
			useProxy, _ := v.ToBoolean()
			if useProxy && Transport != nil {
				req.SetTransport(Transport)
			}
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
		if !strings.Contains(v.Name(), ".js") {
			continue
		}
		jr := ExecPath + "/develop/replies/" + v.Name()
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
		disable := false
		if res := regexp.MustCompile(`\[disable:([^\[\]]+)\]`).FindStringSubmatch(data); len(res) != 0 {
			admin = strings.Trim(res[1], " ") == "true"
		}
		priority := 0
		if res := regexp.MustCompile(`\[priority:([^\[\]]+)\]`).FindStringSubmatch(data); len(res) != 0 {
			priority = Int(strings.Trim(res[1], " "))
		}
		server := ""
		if res := regexp.MustCompile(`\[server:([^\[\]]+)\]`).FindStringSubmatch(data); len(res) != 0 {
			server = strings.TrimSpace(res[1])
		}
		if len(rules) == 0 && cron == "" && server == "" {
			logs.Warn("回复：%s无效文件", jr, err)
			continue
		}
		var handler = func(s Sender) interface{} {
			data, err := os.ReadFile(jr)
			if err != nil {
				return nil
			}
			template := string(data)
			template = strings.Replace(template, "ImType()", fmt.Sprintf(`"%s"`, s.GetImType()), -1)
			param := func(call otto.Value) otto.Value {
				i, _ := call.ToInteger()
				v, _ := otto.ToValue(s.Get(int(i - 1)))
				return v
			}
			vm := otto.New()
			vm.Set("call", func(name otto.Value, arg otto.Value) interface{} {
				key := name.String()
				value := arg.String()
				if f, ok := OttoFuncs[key]; ok {
					v, _ := otto.ToValue(f(value))
					return v
				}
				return otto.Value{}
			})
			vm.Set("cancall", func(name otto.Value) interface{} {
				key := name.String()
				if _, ok := OttoFuncs[key]; ok {
					return otto.TrueValue()
				}
				return otto.FalseValue()
			})
			vm.Set("Delete", func() {
				s.Delete()
			})
			vm.Set("GetChatID", func() otto.Value {
				v, _ := otto.ToValue(s.GetChatID())
				return v
			})
			vm.Set("Continue", func() {
				s.Continue()
			})
			vm.Set("GetUsername", func() otto.Value {
				v, _ := otto.ToValue(s.GetUsername())
				return v
			})
			vm.Set("GetChatname", func() otto.Value {
				v, _ := otto.ToValue(s.GetChatname())
				return v
			})
			vm.Set("Debug", func(str otto.Value) otto.Value {
				logs.Debug(str)
				return otto.Value{}
			})
			vm.Set("GroupKick", func(uid otto.Value, reject_add_request otto.Value) {
				f, _ := reject_add_request.ToBoolean()
				s.GroupKick(uid.String(), f)
			})
			vm.Set("GroupBan", func(uid otto.Value, duration otto.Value) {
				f, _ := duration.ToInteger()
				s.GroupBan(uid.String(), int(f))
			})
			vm.Set("GetUserID", func() otto.Value {
				v, _ := otto.ToValue(s.GetUserID())
				return v
			})
			vm.Set("GetContent", func() otto.Value {
				v, _ := otto.ToValue(s.GetContent())
				return v
			})
			vm.Set("breakIn", func(str otto.Value) otto.Value {
				s := s.Copy()
				s.SetContent(str.String())
				Senders <- s
				return otto.Value{}
			})
			vm.Set("input", func(vs ...otto.Value) interface{} {
				str := ""
				var i int64
				j := ""
				if len(vs) > 0 {
					i, _ = vs[0].ToInteger()
				}
				if len(vs) > 1 {
					j, _ = vs[1].ToString()
				}
				options := []interface{}{}
				options = append(options, time.Duration(i)*time.Millisecond)
				if j != "" {
					options = append(options, ForGroup)
				}
				if rt := s.Await(s, nil, options...); rt != nil {
					str = rt.(string)
				}
				v, _ := otto.ToValue(str)
				return v
			})

			vm.Set("sleep", sleep)
			vm.Set("isAdmin", func() interface{} {
				if s.IsAdmin() {
					return otto.TrueValue()
				}
				return otto.FalseValue()
			})
			vm.Set("set", set)
			vm.Set("param", param)
			vm.Set("get", get)
			vm.Set("bucketGet", bucketGet)
			vm.Set("bucketSet", bucketSet)
			vm.Set("bucketKeys", bucketKeys)
			vm.Set("request", request)
			vm.Set("push", push)
			vm.Set("sendText", func(call otto.Value) interface{} {
				s.Reply(call.String())
				return otto.Value{}
			})
			vm.Set("image", func(call otto.Value) interface{} {
				v, _ := otto.ToValue(`[CQ:image,file=` + call.String() + `]`)
				return v
			})
			vm.Set("sendImage", func(call otto.Value) interface{} {
				s.Reply(ImageUrl(call.String()))
				return otto.Value{}
			})
			vm.Set("sendVideo", func(call otto.Value) interface{} {
				url := call.String()
				if url == "" {
					return otto.Value{}
				}
				s.Reply(VideoUrl(url))
				return otto.Value{}
			})
			rt, err := vm.Run(template)
			if err != nil {
				return err
			}
			result := rt.String()
			for _, v := range regexp.MustCompile(`\[image:\s*([^\s\[\]]+)\s*\]`).FindAllStringSubmatch(result, -1) {
				s.Reply(ImageUrl(v[1]))
				result = strings.Replace(result, fmt.Sprintf(`[image:%s]\n`, v[1]), "", -1)
			}
			if result == "" || result == "undefined" {
				return nil
			}
			return result
		}
		logs.Warn("回复：%s添加成功", jr)
		AddCommand("", []Function{
			{
				Handle:   handler,
				Rules:    rules,
				Cron:     cron,
				Admin:    admin,
				Priority: priority,
				Disable:  disable,
				Server:   server,
			},
		})
	}
}
