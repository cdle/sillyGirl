package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/buger/jsonparser"
)

func init() {
	AddCommand("", []Function{
		{
			Rules: []string{
				"^小爱(.*)$",
			},
			Handle: func(s Sender) interface{} {
				api := sillyGirl.Get("小爱同学")
				if api == "" {
					return "未设置小同学api"
				}
				reply := func(str string) string {
					str, _ = httplib.Get(fmt.Sprintf(api, str)).String()
					if gjson := sillyGirl.Get("小爱同学gjson"); gjson != "" {
						str, _ = jsonparser.GetString([]byte(str), strings.Split(gjson, ".")...)
					}
					if str == "" {
						str = "小爱暂时无法答复您。"
					}
					return str
				}
				msg := s.Get()
				msg = strings.Trim(msg, " ")
				if strings.Contains(msg, "对话模式") {
					stop := false
					s.Reply(reply("小爱"))
					for {
						if stop {
							return nil
						}
						s.Await(s, func(s2 Sender) interface{} {
							msg := s2.GetContent()
							msg = strings.Trim(msg, " ")
							if strings.Contains(msg, "闭嘴") {
								stop = true
							}
							return reply(msg)
						}, `[\s\S]*`, time.Duration(time.Second*5000))
					}
				}
				if msg == "" || msg == "同学" {
					msg = "小爱"
					go func() {
						s.Reply(reply(s.Await(s, nil).(string)))
					}()
				}
				return reply(msg)
			},
		},
	})
}
