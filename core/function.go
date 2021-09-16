package core

import (
	"fmt"
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
	Regex   bool
	Cron    string
}

var pname = regexp.MustCompile(`/([^/\s]+)`).FindStringSubmatch(os.Args[0])[1]

var functions = []Function{
	{
		Rules: []string{"^傻妞 (.*)$", "^傻妞$"},
		Handle: func(s im.Sender) interface{} {
			m := s.Get()
			if m != "" {
				s.Reply(fmt.Sprintf("哎呀，傻妞不懂%s是什么意思啦。", m))
			} else {
				s.Reply("请说，我在。")
			}
			return nil
		},
	},
	{
		Rules: []string{"^升级$"},
		Admin: true,
		Handle: func(s im.Sender) interface{} {
			if runtime.GOOS == "darwin" {
				return "沙雕。"
			}
			s.Reply("傻妞开始拉取代码。")
			rtn, err := exec.Command("sh", "-c", "cd "+ExecPath+" && git stash && git pull").Output()
			if err != nil {
				return "傻妞拉取代失败：" + err.Error() + "。"
			}
			t := string(rtn)
			if !strings.Contains(t, "changed") {
				if strings.Contains(t, "Already") || strings.Contains(t, "已经是最新") {
					return "傻妞已是最新版啦。"
				} else {
					return "傻妞拉取代失败：" + t + "。"
				}
			} else {
				s.Reply("傻妞拉取代码成功。")
			}
			s.Reply("傻妞正在编译程序。")
			rtn, err = exec.Command("sh", "-c", "cd "+ExecPath+" && go build -o "+pname).Output()
			if err != nil {
				return "傻妞编译失败：" + err.Error()
			} else {
				s.Reply("傻妞编译成功。")
			}
			s.Reply("傻妞重启程序。")
			Daemon()
			return nil
		},
	},
	{
		Rules: []string{"^重启$"},
		Admin: true,
		Handle: func(s im.Sender) interface{} {
			s.Reply("傻妞重启程序。")
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
	for _, cmd := range cmds {
		if !cmd.Regex {
			for i := range cmd.Rules {
				if prefix != "" {
					cmd.Rules[i] = prefix + `\s+` + cmd.Rules[i]
				}
				cmd.Rules[i] = strings.Replace(cmd.Rules[i], "(", `[(]`, -1)
				cmd.Rules[i] = strings.Replace(cmd.Rules[i], ")", `[)]`, -1)
				cmd.Rules[i] = strings.Replace(cmd.Rules[i], " ", `\s+`, -1)
				cmd.Rules[i] = strings.Replace(cmd.Rules[i], "?", `(\S+)`, -1)
				cmd.Rules[i] = "^" + cmd.Rules[i] + "$"
			}
		}
		functions = append(functions, cmd)
		if cmd.Cron != "" {
			c.AddFunc(cmd.Cron, func() {
				cmd.Handle(nil)
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
