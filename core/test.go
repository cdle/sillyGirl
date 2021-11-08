package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/beego/beego/v2/adapter/httplib"
)

func init() {
	go func() {
		v := sillyGirl.Get("rebootInfo")
		defer sillyGirl.Set("rebootInfo", "")
		if v != "" {
			vv := strings.Split(v, " ")
			tp, cd, ud := vv[0], Int(vv[1]), Int(vv[2])
			if tp == "fake" { //&& sillyGirl.GetBool("update_notify", false) == true {
				// time.Sleep(time.Second * 10)
				// NotifyMasters("自动更新完成。")
				return
			}
			msg := "重启完成。"
			for i := 0; i < 10; i++ {
				if cd == 0 {
					if push, ok := Pushs[tp]; ok {
						push(ud, msg)
						break
					}
				} else {
					if push, ok := GroupPushs[tp]; ok {
						push(cd, ud, msg)
						break
					}
				}
				time.Sleep(time.Second)
			}
		}
	}()
}

func initSys() {
	AddCommand("", []Function{
		// {//
		// 	Rules: []string{"unintsall sillyGirl"},
		// 	Admin: true,
		// 	Handle: func(s Sender) interface{} {
		// 		return ""
		// 	},
		// },
		//
		{
			Rules: []string{"raw ^name$"},
			Handle: func(s Sender) interface{} {
				s.Disappear()
				return name()
			},
		},
		{
			Rules: []string{"raw ^升级$"},
			Cron:  "*/1 * * * *",
			Admin: true,
			Handle: func(s Sender) interface{} {
				if runtime.GOOS == "windows" {
					return "windows系统不支持此命令"
				}
				if s.GetImType() == "fake" && !sillyGirl.GetBool("auto_update", true) {
					return nil
				}
				if s.GetImType() != "fake" {
					if compiled_at != "" {
						// prefix := "https://ghproxy.com/"
						//
						// prefix := sillyGirl.Get("download_prefix")
						for _, prefix := range []string{"https://ghproxy.com/", ""} {
							data, _ := httplib.Get(prefix + "https://raw.githubusercontent.com/cdle/binary/master/compile_time.go").String()
							if str := regexp.MustCompile(`\d+`).FindString(data); str != "" && strings.Contains(data, "package") {
								if str > compiled_at {
									s.Reply("正在下载更新...")
									data, err := httplib.Get(prefix + "https://raw.githubusercontent.com/cdle/binary/master/sillyGirl_linux_amd64_" + str).Bytes()
									if err != nil {
										return "下载程序错误：" + err.Error()
									}
									if len(data) < 2646147 {
										return "下载失败。"
									}
									filename := ExecPath + "/" + pname
									if err = os.RemoveAll(filename); err != nil {
										return "删除旧程序错误：" + err.Error()
									}
									if err = os.WriteFile(filename, data, 777); err != nil {
										return "写入程序错误：" + err.Error()
									}
									return "下载完成，请对我说\"重启\"。"
								} else {
									return fmt.Sprintf("当前版本(%s)最新，无需升级。", compiled_at)
								}
							} else {
								continue
							}
						}
						return "无法升级."
					}
				}

				s.Reply("开始检查核心更新...", E)
				update := false
				record := func(b bool) {
					if !update && b {
						update = true
					}
				}
				need, err := GitPull("")
				if err != nil {
					return "请使用以下命令手动升级：\n cd " + ExecPath + " && git stash && git pull && go build && ./" + pname
				}
				if !need {
					s.Reply("核心功能已是最新。", E)
				} else {
					record(need)
					s.Reply("核心功能发现更新。", E)
				}
				files, _ := ioutil.ReadDir(ExecPath + "/develop")
				for _, f := range files {
					if f.IsDir() && f.Name() != "replies" {
						if f.Name() == "qinglong" {
							continue
						}
						if strings.HasPrefix(f.Name(), "_") {
							continue
						}
						s.Reply("检查扩展"+f.Name()+"更新...", E)
						need, err := GitPull("/develop/" + f.Name())
						if err != nil {
							s.Reply("扩展"+f.Name()+"更新错误"+err.Error()+"。", E)
						}
						if !need {
							s.Reply("扩展"+f.Name()+"已是最新。", E)
						} else {
							record(need)
							s.Reply("扩展"+f.Name()+"发现更新。", E)
						}
					}
				}
				if !update {
					s.Reply("没有更新。", E)
					return nil
				}
				s.Reply("正在编译程序...", E)
				if err := CompileCode(); err != nil {
					return "请使用以下命令手动编译：\n cd " + ExecPath + " && go build && ./" + pname
				}
				s.Reply("编译程序完毕。", E)
				sillyGirl.Set("rebootInfo", fmt.Sprintf("%v %v %v", s.GetImType(), s.GetChatID(), s.GetUserID()))
				s.Reply("更新完成，即将重启！", E)
				go func() {
					time.Sleep(time.Second)
					Daemon()
				}()
				return nil
			},
		},
		{
			Rules: []string{"raw ^编译$"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				if sillyGirl.Get("compiled_at") == "" {
					return "编译个🐔8。"
				}
				s.Reply("正在编译程序...", E)
				if err := CompileCode(); err != nil {
					return err
				}
				s.Reply("编译程序完毕。", E)
				return nil
			},
		},
		{
			Rules: []string{"raw ^重启$"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				s.Disappear()
				sillyGirl.Set("rebootInfo", fmt.Sprintf("%v %v %v", s.GetImType(), s.GetChatID(), s.GetUserID()))
				s.Reply("即将重启！", E)
				Daemon()
				return nil
			},
		},
		{
			Rules: []string{"raw ^命令$"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				s.Disappear()
				ss := []string{}
				ruless := [][]string{}
				for _, f := range functions {
					if len(f.Rules) > 0 {
						rules := []string{}
						for i := range f.Rules {
							rules = append(rules, f.Rules[i])
						}
						ruless = append(ruless, rules)
					}
				}
				for j := range ruless {
					for i := range ruless[j] {
						ruless[j][i] = strings.Trim(ruless[j][i], "^$")
						ruless[j][i] = strings.Replace(ruless[j][i], `(\S+)`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `(\S*)`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `(.+)`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `(.*)`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `\s+`, " ", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `\s*`, " ", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `.+`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `.*`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `[(]`, "(", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `[)]`, ")", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `([\s\S]+)`, "?", -1)
					}
					ss = append(ss, strings.Join(ruless[j], "\n"))
				}

				return strings.Join(ss, "\n")
			},
		},
		{
			Admin: true,
			Rules: []string{"set ? ? ?", "delete ? ?", "? set ? ?", "? delete ?", "set ? ?", "? set ?"},
			Handle: func(s Sender) interface{} {
				name := s.Get(0)
				if name == "silly" {
					name = "sillyGirl"
				}
				b := Bucket(name)
				if !IsBucket(b) {
					s.Continue()
					return nil
				}
				old := b.Get(s.Get(1))
				b.Set(s.Get(1), s.Get(2))
				go func() {
					s.Await(s, func(_ Sender) interface{} {
						b.Set(s.Get(1), old)
						return "已撤回。"
					}, "^撤回$", time.Second*60)
				}()
				return "操作成功，在60s内可\"撤回\"。"
			},
		},
		{
			Admin: true,
			Rules: []string{"get ? ?", "? get ?"},
			Handle: func(s Sender) interface{} {

				name := s.Get(0)
				if name == "silly" {
					name = "sillyGirl"
				}
				b := Bucket(name)
				if !IsBucket(b) {
					s.Continue()
					return nil
				}
				s.Disappear()
				v := b.Get(s.Get(1))
				if v == "" {
					return errors.New("无值")
				}
				return v
			},
		},
		{
			Admin: true,
			Rules: []string{"list ?"},
			Handle: func(s Sender) interface{} {
				name := s.Get(0)
				if name == "silly" {
					name = "sillyGirl"
				}
				if s.GetChatID() != 0 && name != "reply" {
					return "请私聊我。"
				}
				if name != "otto" && name != "reply" && name != "sillyGirl" && name != "qinglong" {
					s.Continue()
					return nil
				}
				if s.GetChatID() != 0 {
					s.Disappear()
				}
				b := Bucket(name)
				if !IsBucket(b) {
					s.Continue()
					return nil
				}
				rt := ""
				b.Foreach(func(k, v []byte) error {
					rt += fmt.Sprintf("%s === %s\n", k, v)
					return nil
				})
				return strings.Trim(rt, "\n")
			},
		},
		{
			Admin: true,
			Rules: []string{"send ? ? ?"},
			Handle: func(s Sender) interface{} {
				Push(s.Get(0), Int(s.Get(1)), s.Get(2))
				return "发送成功呢"
			},
		},
		{
			Rules: []string{"raw ^myuid$"},
			Handle: func(s Sender) interface{} {
				return fmt.Sprint(s.GetUserID())
			},
		},
		{
			Rules: []string{"raw ^groupCode$"},
			Handle: func(s Sender) interface{} {
				return fmt.Sprint(s.GetChatID())
			},
		},
		{
			Rules: []string{"raw ^compiled_at$"},
			Handle: func(s Sender) interface{} {
				return sillyGirl.Get("compiled_at")
			},
		},
		{
			Rules: []string{"notify ?"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				NotifyMasters(s.Get())
				return "通知成功。"
			},
		},
		{
			Rules: []string{"raw ^started_at$"},
			Handle: func(s Sender) interface{} {
				return sillyGirl.Get("started_at")
			},
		},
		{
			Rules: []string{"守护傻妞"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				service := `
[Unit]
Description=silly silly girl bot
After=network.target mysql.service mariadb.service mysqld.service
[Service]
Type=forking
ExecStart=` + ExecPath + "/" + pname + ` -d
PIDFile=/var/run/sillyGirl.pid
Restart=always
User=root
Group=root
				
[Install]
WantedBy=multi-user.target
Alias=sillyGirl.service`
				data, err := exec.Command("sh", "-c", "type systemctl").Output()
				if err != nil {
					s.Reply(err)
					return nil
				}

				if !strings.Contains(string(data), "bin") {
					s.Reply(data)
					return nil
				}
				os.WriteFile("/usr/lib/systemd/system/sillyGirl.service", []byte(service), 0o644)
				exec.Command("systemctl", "disable", string(sillyGirl)).Output()
				exec.Command("systemctl", "enable", string(sillyGirl)).Output()
				return "电脑重启后生效。"
			},
		},
		// {
		// 	Rules: []string{"raw .*pornhub.*"},
		// 	Handle: func(s Sender) interface{} {
		// 		s.Reply("你已涉黄永久禁言。")
		// 		for {
		// 			s.Await(s, func(s2 Sender, _ error) interface{} {
		// 				s2.Disappear(time.Millisecond * 50)
		// 				return "你已被禁言。"
		// 			}, `[\s\S]*`, time.Duration(time.Second*300))
		// 		}
		// 	},
		// },
		{
			Rules: []string{"raw ^成语接龙$"},
			Handle: func(s Sender) interface{} {
				begin := ""
				fword := func(cy string) string {
					begin = strings.Replace(regexp.MustCompile(`([一-龥])】`).FindString(cy), "】", "", -1)
					return begin
				}
				id := fmt.Sprintf("%v", s.GetUserID())
			start:
				data, err := httplib.Get("http://hm.suol.cc/API/cyjl.php?id=" + id + "&msg=开始成语接龙").String()
				if err != nil {
					s.Reply(err)
				}
				s.Reply(data)
				fword(data)
				stop := false
				win := false
				if strings.Contains(data, "你赢") {
					stop = true
					win = true
				}
				if strings.Contains(data, "我赢") {
					stop = true
				}
				if !stop {
					s.Await(s, func(s2 Sender) interface{} {
						ct := s2.GetContent()
						me := s2.GetUserID() == s.GetUserID()
						if strings.Contains(ct, "小爱提示") || ct == "q" {
							s2.SetContent(fmt.Sprintf("小爱%s字开头的成语有哪些？", begin))
							s2.Continue()
							return Again
						}
						if strings.Contains(ct, "认输") {
							if me {
								stop = true
								return nil
							} else {
								return GoAgain("你认输有个屁用。")
							}
						}
						if regexp.MustCompile("^"+begin).FindString(ct) == "" || strings.Contains(ct, "接龙") {
							if me {
								return GoAgain(fmt.Sprintf("现在是接【%s】开头的成语哦。", begin))
							} else {
								s2.Continue()
								return Again
							}
						}
						cy := regexp.MustCompile("^[一-龥]+$").FindString(ct)
						if cy == "" {
							s2.Disappear(time.Millisecond * 500)
							return GoAgain("请认真接龙，一站到底！")
						}
						data, err := httplib.Get("http://hm.suol.cc/API/cyjl.php?id=" + id + "&msg=我接" + cy).String()
						if err != nil {
							s2.Reply(err)
							return Again
						}
						if strings.Contains(data, "file_get_contents") {
							ss := strings.Split(data, "\n")
							return GoAgain(ss[len(ss)-1])
						}
						if strings.Contains(data, "你赢") {
							stop = true
							win = true
							if !me {
								defer s.Reply("反正不是你赢，嘿嘿。")
							}
						} else if strings.Contains(data, "我赢") {
							stop = true
							win = false
						} else if strings.Contains(data, "恭喜") {
							fword(data)
							if !me {
								data += "\n你很可拷，观棋不语真君子懂不懂啊。"
							}
						} else {
							if me {
								data += "\n玩不过就认输呗。"
							} else {
								data += "\n你以为你会，结果出丑了吧。"
							}
						}
						if !stop {
							return GoAgain(data)
						}
						return data
					}, ForGroup)
				}
				time.Sleep(time.Microsecond * 100)
				s.Reply("还玩吗？[Y/n]")
				if s.Await(s, func(s2 Sender) interface{} {
					return YesNo
				}, ForGroup) == Yes {
					goto start
				}
				if !win {
					s.Reply("菜*，见一次虐一次！")
				} else {
					s.Reply("大爷下次再来玩啊～")
				}
				return nil
			},
		},
		{
			Rules: []string{"^machineId$"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				return OttoFuncs["machineId"]("")
			},
		},
		{
			Rules: []string{"^time$"},
			Handle: func(s Sender) interface{} {
				return OttoFuncs["timeFormat"]("2006-01-02 15:04:05")
			},
		},
	})
}

func IsBucket(b Bucket) bool {
	for i := range Buckets {
		if Buckets[i] == b {
			return true
		}
	}
	return false
}
