package core

import (
	"fmt"
	"regexp"

	"github.com/cdle/sillyGirl/im"
	"github.com/cdle/sillyGirl/im/tg"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Function struct {
	Rules   []string
	FindAll bool
	Admin   bool
	Handle  func(s im.Sender) bool
}

var functions = []Function{
	{
		Rules: []string{"^傻妞 (.*)$", "^傻妞$"},
		Handle: func(s im.Sender) bool {
			m := s.Get()
			if m != "" {
				s.Reply(fmt.Sprintf("哎呀，傻妞不懂%s是什么意思啦。", m))
			} else {
				s.Reply("请说，我在。")
			}
			return true
		},
	},
}

var Senders chan im.Sender

func initToHandleMessage() {
	for _, im := range Config.Im {
		switch im.Type {
		case "tg":
			tg.Handler = func(message *tb.Message) {
				Senders <- &tg.Sender{
					Message: message,
				}
			}
			go tg.RunBot(&im)
		case "qq":

		}
	}
	Senders = make(chan im.Sender)
	go func() {
		for {
			go handleMessage(<-Senders)
		}
	}()
}

func AddCommand(cmd *Function) {
	functions = append(functions, *cmd)
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
			if matched && function.Handle(sender) {
				return
			}
		}
	}
}
