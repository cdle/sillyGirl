package core

import (
	"crypto/md5"
	"encoding/json"
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
	"github.com/dop251/goja"
)

type JsReply string

var o = NewBucket("otto")

var OttoFuncs = map[string]func(string) string{
	"machineId": func(_ string) string {
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

func Init123() {
	files, err := ioutil.ReadDir(ExecPath + "/develop/replies")
	if err != nil {
		os.MkdirAll(ExecPath+"/develop/replies", os.ModePerm)
		// logs.Warn("打开文件夹%s错误，%v", "develop/replies", err)
		return
	}
	get := func(key string) string {
		return o.Get(key)
	}
	bucketGet := func(bucket, key string) string {
		return o.Get(key, Bucket(bucket).Get(key))
	}
	bucketSet := func(bucket, key, value string) {
		Bucket(bucket).Set(key, value)
	}
	bucketKeys := func(bucket string) []string {
		b := Bucket(bucket)
		if !IsBucket(b) {
			return []string{}
		}
		slice := []string{}
		b.Foreach(func(k, _ []byte) error {
			slice = append(slice, string(k))
			return nil
		})
		return slice
	}
	set := func(key, value string) {
		o.Set(key, value)
	}
	notifyMasters := func(content string) {
		NotifyMasters(content)
	}
	sleep := func(i int) {
		time.Sleep(time.Duration(i) * time.Millisecond)
	}
	push := func(obj *goja.Object) {
		imType := ""
		groupCode := 0
		userID := ""
		content := ""
		for _, key := range obj.Keys() {
			switch key {
			case "imType":
				imType = obj.Get(key).String()
			case "groupCode":
				groupCode = int(obj.Get(key).ToInteger())
			case "chatID":
				groupCode = int(obj.Get(key).ToInteger())
			case "userID":
				userID = obj.Get(key).String()
			case "content":
				content = obj.Get(key).String()
			}
		}
		gid := Int(groupCode)
		if gid != 0 {
			if push, ok := GroupPushs[imType]; ok {
				push(int(gid), userID, content, "")
			}
		} else {
			if push, ok := Pushs[imType]; ok {
				push(userID, content, nil, "")
			}
		}
		return
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
			param := func(i int) string {
				return s.Get(int(i - 1))
			}
			vm := goja.New()
			request := func(obj *goja.Object) interface{} {
				url := ""
				dataType := ""
				method := "get"
				body := ""
				var useProxy bool
				var headers *goja.Object
				var req *httplib.BeegoHTTPRequest

				for _, key := range obj.Keys() {
					switch strings.ToLower(key) {
					case "url":
						url = obj.Get(key).String()
					case "datatype":
						dataType = obj.Get(key).String()
					case "body":
						v := obj.Get(key).String()
						if v == `[object Object]` {
							d, _ := obj.Get(key).ToObject(vm).MarshalJSON()
							body = string(d)
						} else {
							body = obj.Get(key).String()
						}
					case "method":
						method = obj.Get(key).String()
					case "headers":
						headers = obj.Get(key).ToObject(vm)
					case "useproxy":
						if obj.Get(key).ToBoolean() {
							useProxy = !useProxy
						}
					}
				}
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
				if headers != nil {
					for _, key := range headers.Keys() {
						req.Header(key, headers.Get(key).String())
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
					return url
				}
				data, err := req.Bytes()
				if err != nil {
					return ""
				}
				if strings.Contains(dataType, "json") {
					var x = map[string]interface{}{}
					json.Unmarshal(data, &x)
					return x
				}
				return string(data)
			}
			vm.Set("call", func(key, value string) interface{} {
				if f, ok := OttoFuncs[key]; ok {
					return f(value)
				}
				return nil
			})

			vm.Set("cancall", func(key string) interface{} {
				_, ok := OttoFuncs[key]
				return ok
			})
			vm.Set("Delete", s.Delete)
			vm.Set("GetChatID", s.GetChatID)
			vm.Set("ImType", func() string {
				return s.GetImType()
			})
			vm.Set("Continue", s.Continue)
			vm.Set("GetUsername", s.GetUsername)
			vm.Set("GetChatname", s.GetChatname)
			vm.Set("Debug", func(str string) {
				logs.Debug(str)
			})
			vm.Set("GroupKick", func(uid string, reject_add_request bool) {
				s.GroupKick(uid, reject_add_request)
			})
			vm.Set("GroupBan", func(uid string, t int) {
				s.GroupBan(uid, t)
			})
			vm.Set("GetUserID", s.GetUserID)
			vm.Set("GetContent", s.GetContent)
			vm.Set("notifyMasters", notifyMasters)
			vm.Set("breakIn", func(str string) {
				s := s.Copy()
				s.SetContent(str)
				Senders <- s
			})
			vm.Set("input", func(vs ...interface{}) string {
				str := ""
				var i int64
				j := ""
				if len(vs) > 0 {
					i = Int64(vs[0])
				}
				if len(vs) > 1 {
					j = fmt.Sprint(vs[1])
				}
				options := []interface{}{}
				options = append(options, time.Duration(i)*time.Millisecond)
				if j != "" {
					options = append(options, ForGroup)
				}
				if rt := s.Await(s, nil, options...); rt != nil {
					str = rt.(string)
				}
				return str
			})

			vm.Set("sleep", sleep)
			vm.Set("isAdmin", s.IsAdmin)
			vm.Set("set", set)
			vm.Set("param", param)
			vm.Set("get", get)
			vm.Set("bucketGet", bucketGet)
			vm.Set("bucketSet", bucketSet)
			vm.Set("bucketKeys", bucketKeys)
			vm.Set("request", request)
			vm.Set("push", push)
			vm.Set("sendText", func(text string) {
				s.Reply(text)

			})
			vm.Set("image", func(url string) interface{} {
				return `[CQ:image,file=` + url + `]`
			})
			vm.Set("sendImage", func(url string) {
				s.Reply(ImageUrl(url))
			})
			vm.Set("sendVideo", func(url string) {
				if url == "" {
					return
				}
				s.Reply(VideoUrl(url))
			})
			rt, err := vm.RunString(template)
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
		// logs.Warn("回复：%s添加成功", jr)
		AddCommand("", []Function{
			{
				Handle:   handler,
				Rules:    rules,
				Cron:     cron,
				Admin:    admin,
				Priority: priority,
				Disable:  disable,
			},
		})
	}
}
