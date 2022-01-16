package core

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

var OttoFuncs = map[string]interface{}{
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

type Strings struct {
}

func (sender *Strings) Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func (sender *Strings) Replace(s string, old string, new string, n int) string {
	return strings.Replace(s, old, new, n)
}

func (sender *Strings) ReplaceAll(s string, old string, new string) string {
	return strings.ReplaceAll(s, old, new)
}

type Fmt struct {
}

func (sender *Fmt) Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func (sender *Fmt) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func (sender *Fmt) Println(a ...interface{}) (int, error) {
	return fmt.Println(a...)
}

func (sender *Fmt) Print(a ...interface{}) (int, error) {
	return fmt.Print(a...)
}

type JsSender struct {
	Sender Sender
	again  string
}

func (sender *JsSender) Continue() {
	sender.Sender.Continue()
}

func (sender *JsSender) GetUserID() string {
	return sender.Sender.GetUserID()
}

func (sender *JsSender) SetContent(s string) {
	sender.Sender.SetContent(s)
}

func (sender *JsSender) GetContent() string {
	return sender.Sender.GetContent()
}
func (sender *JsSender) GetImType() string {
	return sender.Sender.GetImType()
}
func (sender *JsSender) RecallMessage(p ...interface{}) {
	sender.Sender.RecallMessage(p...)
}
func (sender *JsSender) GetUsername() string {
	return sender.Sender.GetUsername()
}
func (sender *JsSender) GetMessageID() string {
	return sender.Sender.GetMessageID()
}

func (sender *JsSender) GetGroupCode() int {
	return sender.Sender.GetChatID()
}
func (sender *JsSender) IsAdmin() bool {
	return sender.Sender.IsAdmin()
}
func (sender *JsSender) Reply(text string) []string {
	if text == "" {
		return []string{}
	}
	i, _ := sender.Sender.Reply(text)
	return i
}
func (sender *JsSender) Await(timeout int, fg bool, hd func(*JsSender) string) *JsSender {
	options := []interface{}{}
	if timeout != 0 {
		options = append(options, time.Duration(timeout)*time.Millisecond)
	}
	if fg {
		options = append(options, ForGroup)
	}
	var newJsSender *JsSender
	sender.Sender.Await(sender.Sender, func(sender Sender) interface{} {
		newJsSender = &JsSender{}
		newJsSender.Sender = sender
		if hd != nil {
			var rt = hd(newJsSender)
			if strings.HasPrefix(rt, "go_again_") {
				rt = strings.Replace(rt, "go_again_", "", 1)
				return GoAgain(rt)
			} else {
				if rt == "" {
					return nil
				}
				return rt
			}
		}
		return nil
	}, options...)
	return newJsSender
}

func Init123() {
	files, err := ioutil.ReadDir(ExecPath + "/develop/replies")
	if err != nil {
		os.MkdirAll(ExecPath+"/develop/replies", os.ModePerm)

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
		basePath := ExecPath + "/develop/replies/"
		jr := basePath + v.Name()
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
		//取前1000个字符应该够了,否则有些js长度太长影响正则判断性能
		l := 1000
		if len([]rune(data)) < 1000 {
			l = len([]rune(data))
		}
		data = string([]rune(data)[:l])
		rules := []string{}
		for _, res := range regexp.MustCompile(`\[rule:(.+)]`).FindAllStringSubmatch(data, -1) {
			rules = append(rules, strings.Trim(res[1], " "))
		}
		cron := ""
		if res := regexp.MustCompile(`\[cron:([^\[\]]+)]`).FindStringSubmatch(data); len(res) != 0 {
			cron = strings.Trim(res[1], " ")
		}
		admin := false
		if res := regexp.MustCompile(`\[admin:([^\[\]]+)]`).FindStringSubmatch(data); len(res) != 0 {
			admin = strings.Trim(res[1], " ") == "true"
		}
		disable := false
		if res := regexp.MustCompile(`\[disable:([^\[\]]+)]`).FindStringSubmatch(data); len(res) != 0 {
			disable = strings.Trim(res[1], " ") == "true"
		}
		priority := 0
		if res := regexp.MustCompile(`\[priority:([^\[\]]+)]`).FindStringSubmatch(data); len(res) != 0 {
			priority = Int(strings.Trim(res[1], " "))
		}
		server := ""
		if res := regexp.MustCompile(`\[server:([^\[\]]+)]`).FindStringSubmatch(data); len(res) != 0 {
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
			vm.Set("call", func(key string) interface{} {
				if f, ok := OttoFuncs[key]; ok {
					return f
				}
				return nil
			})
			vm.Set("require", require)
			vm.Set("Request", newrequest)
			vm.Set("request", request)
			vm.Set("cancall", func(key string) interface{} {
				_, ok := OttoFuncs[key]
				return ok
			})
			vm.Set("Delete", s.Delete)
			vm.Set("timeFmt", func(str string) string {
				return time.Now().Format(str)
			})
			vm.Set("GetChatID", s.GetChatID)
			vm.Set("GoAgain", func(s string) string {
				return "go_again_" + s
			})
			vm.Set("GetImType", s.GetImType)
			vm.Set("ImType", s.GetImType)
			vm.Set("Continue", s.Continue)
			vm.Set("GetUsername", s.GetUsername)
			vm.Set("GetChatname", s.GetChatname)
			vm.Set("GetMessageID", s.GetMessageID)
			vm.Set("RecallMessage", s.RecallMessage)
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
			vm.Set("sendText", func(text string) []string {
				if text == "" {
					return []string{}
				}
				i, _ := s.Reply(text)
				return i
			})
			vm.Set("Logger", Logger)
			vm.Set("console", console)
			vm.Set("SillyGirl", SillyGirl)
			vm.Set("image", func(url string) interface{} {
				return `[CQ:image,file=` + url + `]`
			})

			vm.Set("sendImage", func(url string) []string {
				if url == "" {
					return nil
				}
				i, _ := s.Reply(ImageUrl(url))
				return i
			})
			vm.Set("sendVideo", func(url string) []string {
				if url == "" {
					return nil
				}
				i, _ := s.Reply(VideoUrl(url))
				return i
			})

			vm.Set("Sender", &JsSender{
				Sender: s,
			})

			vm.Set("fmt", &Fmt{})
			vm.Set("strings", &Strings{})
			vm.Set("nil", nil)

			// vm.Set("fmt", fmt)
			importedJs := make(map[string]struct{})
			importedJs[jr[len(basePath):]] = struct{}{}
			//2个或者2个以上"/"
			regexp1, _ := regexp.Compile("/{2,}")
			importJs := func(file string) error {
				if file == "" {
					return errors.New("路径不能为空")
				}
				if strings.Contains(file, "..") {
					return errors.New("不能使用父路径")
				}
				file = strings.Replace(file, "./", "", -1)
				file = regexp1.ReplaceAllString(file, "/")
				if !strings.HasSuffix(file, ".js") {
					file = file + ".js"
				}
				if _, ok := importedJs[file]; ok {
					return nil
				}
				importedJs[file] = struct{}{}
				filePath := basePath + file
				f, err := os.Open(filePath)
				if err != nil {
					return err
				}
				v, _ := ioutil.ReadAll(f)
				vm.RunString(string(v))
				return nil
			}
			vm.Set("importJs", importJs)
			vm.Set("importDir", func(dir string) error {
				if dir == "" {
					return errors.New("路径不能为空")
				}
				if strings.Contains(dir, "..") {
					return errors.New("不能使用父路径")
				}
				dir = strings.Replace(dir, "./", "", -1)
				dir = regexp1.ReplaceAllString(dir, "/")
				//统一处理为没有前后"/"
				dir = strings.TrimPrefix(dir, "/")
				dir = strings.TrimSuffix(dir, "/")
				files, err := ioutil.ReadDir(basePath + dir)
				if err != nil {
					return err
				}
				var firstErr error = nil
				for _, v := range files {
					if v.IsDir() {
						continue
					}
					if !strings.Contains(v.Name(), ".js") {
						continue
					}
					if err := importJs(dir + "/" + v.Name()); err != nil {
						if firstErr == nil {
							firstErr = err
						}
					}
				}
				return firstErr
			})
			_, err = vm.RunString(template)
			if err != nil {
				if strings.Contains(err.Error(), "window") {
					return nil
				}
				return err
			}
			return nil
			// result := rt.String()
			// for _, v := range regexp.MustCompile(`\[image:\s*([^\s\[\]]+)\s*]`).FindAllStringSubmatch(result, -1) {
			// 	s.Reply(ImageUrl(v[1]))
			// 	result = strings.Replace(result, fmt.Sprintf(`[image:%s]\n`, v[1]), "", -1)
			// }
			// if result == "" || result == "undefined" {
			// 	return nil
			// }
			// return result
		}

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

func ToImage(url string) string {
	return `[CQ:image,file=` + url + `]`
}
