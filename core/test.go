package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
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
			if tp == "fake" { //&& sillyGirl.GetBool("update_notify", false) == true { //
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
				if s.GetImType() == "fake" && !sillyGirl.GetBool("auto_update", true) {
					return nil
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
					return err
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
					return err
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
				for _, f := range functions {
					// f := f
					// for i := range f.Rules {
					// 	f.Rules[i] = strings.Trim(f.Rules[i], "^$")
					// 	f.Rules[i] = strings.Replace(f.Rules[i], `\s+`, " ", -1)
					// 	f.Rules[i] = strings.Replace(f.Rules[i], `(\S+)`, "?", -1)
					// 	f.Rules[i] = strings.Replace(f.Rules[i], `[(]`, "(", -1)
					// 	f.Rules[i] = strings.Replace(f.Rules[i], `[)]`, ")", -1)
					// }
					ss = append(ss, strings.Join(f.Rules, " "))
				}
				return strings.Join(ss, "\n")
			},
		},
		{
			Admin: true,
			Rules: []string{"set ? ? ?", "delete ? ?"},
			Handle: func(s Sender) interface{} {
				b := Bucket(s.Get(0))
				if !IsBucket(b) {
					return errors.New("不存在的存储桶")
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
			Rules: []string{"get ? ?"},
			Handle: func(s Sender) interface{} {
				s.Disappear()
				b := Bucket(s.Get(0))
				if !IsBucket(b) {
					return errors.New("不存在的存储桶")
				}
				v := b.Get(s.Get(1))
				if v == "" {
					return errors.New("无值")
				}
				return v
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
				return fmt.Sprintf("%d", s.GetChatID())
			},
		},
		{
			Rules: []string{"raw ^compiled_at$"},
			Handle: func(s Sender) interface{} {
				return sillyGirl.Get("compiled_at")
			},
		},
		{
			Rules: []string{"raw ^started_at$"},
			Handle: func(s Sender) interface{} {
				return sillyGirl.Get("started_at")
			},
		},
		{
			Rules: []string{"^守护傻妞"},
			Handle: func(s Sender) interface{} {
				service := `
[Service]
Type=forking
ExecStart=` + ExecPath + "/" + pname + ` -d
PIDFile=/var/run/` + pname + `.pid
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
				id := fmt.Sprintf("%v", s.GetUserID())
				data, err := httplib.Get("http://hm.suol.cc/API/cyjl.php?id=" + id + "&msg=开始成语接龙").String()
				if err != nil {
					s.Reply(err)
				}
				s.Reply(data)
				stop := false
				if strings.Contains(data, "你赢") {
					stop = true
				}
				for {
					if stop == true {
						break
					}
					s.Await(s, func(s2 Sender) interface{} {
						ct := s2.GetContent()
						if ct == "退出接龙" {
							stop = true
							return "不要走决战到天亮，啊哦～"
						}
						if strings.Contains(ct, "认输") {
							stop = true
							return "菜*，见一次虐一次！"
						}
						cy := regexp.MustCompile("^[一-龥]+$").FindString(ct)
						if cy == "" {
							s2.Disappear(time.Millisecond * 500)
							return "请认真接龙，一站到底！"
						}
						data, err := httplib.Get("http://hm.suol.cc/API/cyjl.php?id=" + id + "&msg=我接" + cy).String()
						if err != nil {
							s2.Reply(err)
							return nil
						}
						if strings.Contains(data, "file_get_contents") {
							ss := strings.Split(data, "\n")
							return ss[len(ss)-1]
						}
						if strings.Contains(data, "你赢") {
							stop = true
						} else if strings.Contains(data, "恭喜") {

						} else {
							data += "玩不过就认输呗。"
						}
						return data
					}, `[\s\S]*`, time.Duration(time.Second*300))

				}
				return nil
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
