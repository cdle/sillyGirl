package core

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	cron "github.com/robfig/cron/v3"
)

var c *cron.Cron

func init() {
	if runtime.GOOS != "windows" {
		pname = regexp.MustCompile(`/([^/\s]+)$`).FindStringSubmatch(os.Args[0])[1]
	}
	c = cron.New()
	c.Start()
}

type Function struct {
	Rules   []string
	FindAll bool
	Admin   bool
	Handle  func(s Sender) interface{}
	Cron    string
}

var pname = ""

var name = func() string {
	return sillyGirl.Get("name", "傻妞")
}

var functions = []Function{}

var Senders chan Sender

func initToHandleMessage() {
	Senders = make(chan Sender)
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
			cmds[j].Rules[i] = strings.ReplaceAll(cmds[j].Rules[i], `\r\a\w`, "raw")
			if strings.Contains(cmds[j].Rules[i], "$") {
				continue
			}
			if prefix != "" {
				cmds[j].Rules[i] = prefix + `\s+` + cmds[j].Rules[i]
			}
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], "(", `[(]`, -1)
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], ")", `[)]`, -1)
			cmds[j].Rules[i] = regexp.MustCompile(`\?$`).ReplaceAllString(cmds[j].Rules[i], `(.+)`)
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], " ", `\s+`, -1)
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], "?", `(\S+)`, -1)
			cmds[j].Rules[i] = "^" + cmds[j].Rules[i] + "$"
		}
		functions = append(functions, cmds[j])
		if cmds[j].Cron != "" {
			cmd := cmds[j]
			if _, err := c.AddFunc(cmds[j].Cron, func() {
				cmd.Handle(&Faker{})
			}); err != nil {
				// logs.Warn("任务%v添加失败%v", cmds[j].Rules[0], err)
			} else {
				// logs.Warn("任务%v添加成功", cmds[j].Rules[0])
			}
		}
	}
}

func handleMessage(sender Sender) {
	defer sender.Finish()
	recall := sillyGirl.Get("recall")
	if recall != "" {
		recalled := false
		for _, v := range strings.Split(recall, "&") {
			reg, err := regexp.Compile(v)
			if err == nil {
				if reg.FindString(sender.GetContent()) != "" {
					if !sender.IsAdmin() && sender.GetImType() != "wx" {
						sender.Delete()
						sender.Reply("本妞清除了不好的消息～", time.Duration(time.Second))
						recalled = true
						break
					}
				}
			}
		}
		if recalled == true {
			return
		}
	}

	defer func() {
		logs.Info("%v ==> %v", sender.GetContent(), "finished")
	}()
	u, g, i := fmt.Sprint(sender.GetUserID()), fmt.Sprint(sender.GetChatID()), fmt.Sprint(sender.GetImType())
	con := true
	mtd := false
	waits.Range(func(k, v interface{}) bool {
		c := v.(*Carry)
		vs, _ := url.ParseQuery(k.(string))
		userID := vs.Get("u")
		chatID := vs.Get("c")
		imType := vs.Get("i")
		forGroup := vs.Get("f")
		if imType != i {
			return true
		}
		if chatID != g {
			return true
		}
		if userID != u && forGroup == "" {
			return true
		}
		if m := regexp.MustCompile(c.Pattern).FindString(sender.GetContent()); m != "" {
			mtd = true
			c.Chan <- sender
			sender.Reply(<-c.Result)
			if !sender.IsContinue() {
				con = false
				return false
			}
		}
		return true
	})
	if mtd && !con {
		return
	}
	// if v, ok := waits.Load(key); ok {
	// 	c := v.(*Carry)
	// 	if m := regexp.MustCompile(c.Pattern).FindString(sender.GetContent()); m != "" {
	// 		c.Chan <- sender
	// 		sender.Reply(<-c.Result)
	// 		return
	// 	}
	// }
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
				logs.Info("%v ==> %v", sender.GetContent(), rule)
				if function.Admin && !sender.IsAdmin() {
					sender.Delete()
					sender.Disappear()
					// if sender.GetImType() != "wx" && sender.GetImType() != "qq" {
					sender.Reply("再捣乱我就报警啦～")
					// }
					return
				}
				rt := function.Handle(sender)
				if rt != nil {
					sender.Reply(rt)
				}
				if sender.IsContinue() {
					goto goon
				}
				return
			}
		}
	goon:
	}
	reply.Foreach(func(k, v []byte) error {
		if string(v) == "" {
			return nil
		}
		reg, err := regexp.Compile(string(k))
		if err == nil {
			if reg.FindString(sender.GetContent()) != "" {
				sender.Reply(string(v))
			}
		}
		return nil
	})

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
