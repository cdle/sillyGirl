package core

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/beego/beego/v2/core/logs"
	"github.com/cdle/sillyGirl/utils"
)

var total uint64 = 0
var finished uint64 = 0
var contents sync.Map

type Function struct {
	Rules    []string
	ImType   *Filter
	UserId   *Filter
	GroupId  *Filter
	FindAll  bool
	Admin    bool
	Handle   func(s Sender) interface{}
	Cron     string
	Show     string
	Priority int
	Disable  bool
	Hash     string
}
type Filter struct {
	BlackMode bool
	Items     []string
}

var reply Bucket

var name = func() string {
	return sillyGirl.GetString("name", "傻妞")
}
var Functions = []Function{}

var Senders chan Sender

func initToHandleMessage() {
	reply = MakeBucket("reply")
	Senders = make(chan Sender)
	go func() {
		for {
			s := <-Senders
			if s.GetImType() != "terminal" {
				logs.Info("接收到消息：%s", s.GetContent())
			}
			go HandleMessage(s)
		}
	}()
}

func AddCommand(prefix string, cmds []Function) {
	for j := range cmds {
		if cmds[j].Disable {
			continue
		}
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
			cmds[j].Rules[i] = regexp.MustCompile(`\?$`).ReplaceAllString(cmds[j].Rules[i], `([\s\S]+)`)
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], " ", `\s+`, -1)
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], "?", `(\S+)`, -1)
			cmds[j].Rules[i] = "^" + cmds[j].Rules[i] + "$"
		}
		{
			lf := len(Functions)
			for i := range Functions {
				f := lf - i - 1
				if Functions[f].Priority > cmds[j].Priority {
					Functions = append(Functions[:f+1], append([]Function{cmds[j]}, Functions[f+1:]...)...)
					break
				}
			}
			if len(Functions) == lf {
				if lf > 0 {
					if Functions[0].Priority < cmds[j].Priority && Functions[lf-1].Priority < cmds[j].Priority {
						Functions = append([]Function{cmds[j]}, Functions...)
					} else {
						Functions = append(Functions, cmds[j])
					}
				} else {
					Functions = append(Functions, cmds[j])
				}
			}
		}

		if cmds[j].Cron != "" {
			cmd := cmds[j]
			if _, err := C.AddFunc(cmds[j].Cron, func() {
				cmd.Handle(&Faker{})
			}); err != nil {

			} else {

			}
		}
	}
}

func HandleMessage(sender Sender) {
	num := atomic.AddUint64(&total, 1)
	defer atomic.AddUint64(&finished, 1)
	ct := sender.GetContent()
	contents.Store(num, ct)
	defer func() {
		contents.Delete(num)
	}()
	content := utils.TrimHiddenCharacter(ct)
	defer func() {
		sender.Finish()
		if sender.IsAtLast() {
			s := sender.MessagesToSend()
			if s != "" {
				sender.Reply(s)
			}
		}
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
		if chatID != g && (forGroup != "me" || g != "0") {
			return true
		}
		if userID != u && (forGroup == "" || forGroup == "me") {
			return true
		}
		if m := regexp.MustCompile(c.Pattern).FindString(content); m != "" {
			r := false
			mtd = true
			if f, ok := c.Sender.(*Faker); ok && f.Carry != nil {
				if s1, o := sender.(*Faker); o && s1.Carry != nil {
					f.Carry = s1.Carry
					r = true
				}
			}
			c.Chan <- sender
			sender.Reply(<-c.Result)
			if r {
				sender.(*Faker).Carry = nil
			}
			if !sender.IsContinue() {
				con = false
				return false
			}
			content = utils.TrimHiddenCharacter(sender.GetContent())
		}
		return true
	})
	if mtd && !con {
		return
	}
	replied := false
	MakeBucket(fmt.Sprintf("reply%s%d", sender.GetImType(), sender.GetChatID())).Foreach(func(k, v []byte) error {
		if string(v) == "" {
			return nil
		}
		reg, err := regexp.Compile(string(k))
		if err == nil {
			if reg.FindString(content) != "" {
				replied = true
				r := string(v)
				if strings.Contains(r, "$") {
					sender.Reply(reg.ReplaceAllString(content, r))
				} else {
					sender.Reply(r)
				}
			}
		}
		return nil
	})

	if !replied {
		reply.Foreach(func(k, v []byte) error {
			if string(v) == "" {
				return nil
			}
			reg, err := regexp.Compile(string(k))
			if err == nil {
				if reg.FindString(content) != "" {
					replied = true
					r := string(v)
					if strings.Contains(r, "$") {
						sender.Reply(reg.ReplaceAllString(content, r))
					} else {
						sender.Reply(r)
					}
				}
			}
			return nil
		})
	}

	for _, function := range Functions {
		if black(function.ImType, sender.GetImType()) || black(function.UserId, sender.GetUserID()) || black(function.GroupId, fmt.Sprint(sender.GetChatID())) {
			continue
		}
		for _, rule := range function.Rules {
			var matched bool
			if function.FindAll {
				if res := regexp.MustCompile(rule).FindAllStringSubmatch(content, -1); len(res) > 0 {
					tmp := [][]string{}
					for i := range res {
						tmp = append(tmp, res[i][1:])
					}
					logs.Info("匹配到规则：%s", rule)
					sender.SetAllMatch(tmp)
					matched = true
				}
			} else {
				if res := regexp.MustCompile(rule).FindStringSubmatch(content); len(res) > 0 {
					logs.Info("匹配到规则：%s", rule)
					sender.SetMatch(res[1:])
					matched = true
				}
			}
			if matched {
				if function.Admin && !sender.IsAdmin() {
					sender.Delete()
					sender.Disappear()
					return
				}
				rt := function.Handle(sender)
				if rt != nil {
					sender.Reply(rt)
				}
				if sender.IsContinue() {
					sender.ClearContinue()
					content = utils.TrimHiddenCharacter(sender.GetContent())
					logs.Info("继续去处理：%s", content)
					goto next
				}
				return
			}
		}
	next:
	}

	recall := sillyGirl.GetString("recall")
	if recall != "" {
		recalled := false
		for _, v := range strings.Split(recall, "&") {
			reg, err := regexp.Compile(v)
			if err == nil {
				if reg.FindString(content) != "" {
					if !sender.IsAdmin() && sender.GetImType() != "wx" {
						sender.Delete()
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
}
func black(filter *Filter, str string) bool {
	if filter != nil {
		if filter.BlackMode {
			if utils.Contains(filter.Items, str) {
				return true
			}
		} else {
			if !utils.Contains(filter.Items, str) {
				return true
			}
		}
	}
	return false
}
