package core

import (
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/cdle/sillyGirl/im"
	cron "github.com/robfig/cron/v3"
)

var c *cron.Cron

func init() {
	c = cron.New()
	c.Start()
}

type Function struct {
	Rules   []string
	FindAll bool
	Admin   bool
	Handle  func(s im.Sender) interface{}
	Cron    string
}

var pname = regexp.MustCompile(`/([^/\s]+)`).FindStringSubmatch(os.Args[0])[1]

var name = func() string {
	return sillyGirl.Get("name", "傻妞")
}

var functions = []Function{
	{
		Rules: []string{"^升级$"},
		Admin: true,
		Handle: func(s im.Sender) interface{} {
			if runtime.GOOS == "darwin" {
				return "沙雕。"
			}
			s.Reply(name() + "开始拉取代码。")
			rtn, err := exec.Command("sh", "-c", "cd "+ExecPath+" && git stash && git pull").Output()
			if err != nil {
				return name() + "拉取代失败：" + err.Error() + "。"
			}
			t := string(rtn)
			if !strings.Contains(t, "changed") {
				if strings.Contains(t, "Already") || strings.Contains(t, "已经是最新") {
					return name() + "已是最新版啦。"
				} else {
					return name() + "拉取代失败：" + t + "。"
				}
			} else {
				s.Reply(name() + "拉取代码成功。")
			}
			s.Reply(name() + "正在编译程序。")
			rtn, err = exec.Command("sh", "-c", "cd "+ExecPath+" && go build -o "+pname).Output()
			if err != nil {
				return name() + "编译失败：" + err.Error()
			} else {
				s.Reply(name() + "编译成功。")
			}
			s.Reply(name() + "重启程序。")
			Daemon()
			return nil
		},
	},
	{
		Rules: []string{"^重启$"},
		Admin: true,
		Handle: func(s im.Sender) interface{} {
			s.Reply(name() + "重启程序。")
			Daemon()
			return nil
		},
	},
}

var Senders chan im.Sender

func initToHandleMessage() {
	Senders = make(chan im.Sender)
	go func() {
		for {
			go handleMessage(<-Senders)
		}
	}()
}

func AddCommand(prefix string, cmds []Function) {
	for j := range cmds {
		for i := range cmds[j].Rules {
			if strings.Contains(cmds[j].Rules[i], "raw ") {
				cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], "raw ", "", -1)
				continue
			}
			if prefix != "" {
				cmds[j].Rules[i] = prefix + `\s+` + cmds[j].Rules[i]
			}
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], "(", `[(]`, -1)
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], ")", `[)]`, -1)
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], " ", `\s+`, -1)
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], "?", `(\S+)`, -1)
			cmds[j].Rules[i] = "^" + cmds[j].Rules[i] + "$"
		}
		functions = append(functions, cmds[j])
		if cmds[j].Cron != "" {
			c.AddFunc(cmds[j].Cron, func() {
				cmds[j].Handle(nil)
			})
		}
	}
}

func handleMessage(sender im.Sender) {
	for _, function := range functions {
		for _, rule := range function.Rules {
			var matched bool
			if function.FindAll {
				if res := regexp.MustCompile(rule).FindAllStringSubmatch(sender.GetContent(), -1); len(res) > 0 {
					tmp := [][]string{}
					for i := range res {
						tmp = append(tmp, res[i][1:])
					}
					sender.SetAllMatch(tmp)
					matched = true
				}
			} else {
				if res := regexp.MustCompile(rule).FindStringSubmatch(sender.GetContent()); len(res) > 0 {
					sender.SetMatch(res[1:])
					matched = true
				}
			}
			if matched {
				if function.Admin && !sender.IsAdmin() {
					sender.Reply("没有权限操作")
					return
				}
				rt := function.Handle(sender)
				if rt != nil {
					sender.Reply(rt)
				}
				return
			}
		}
	}
}

func FetchCookieValue(ps ...string) string {
	var key, cookies string
	if len(ps) == 2 {
		if len(ps[0]) > len(ps[1]) {
			key, cookies = ps[1], ps[0]
		} else {
			key, cookies = ps[0], ps[1]
		}
	}
	match := regexp.MustCompile(key + `=([^;]*);{0,1}`).FindStringSubmatch(cookies)
	if len(match) == 2 {
		return match[1]
	} else {
		return ""
	}
}
