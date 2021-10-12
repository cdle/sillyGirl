package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/httplib"
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
					if str == "" {
						str = "暂时无法回复。"
					}
					return str
				}
				msg := s.Get()
				if strings.Contains(msg, "对话模式") {
					stop := false
					for {
						if stop {
							break
						}
						s.Await(s, func(s2 Sender) interface{} {
							msg := s2.Get()
							if msg == "闭嘴" {
								stop = true
							}
							return reply(msg)
						}, `[\s\S]*`, time.Duration(time.Second*300))
					}
				}
				if msg == "" {
					msg = "小爱"
				}
				return reply(msg)
			},
		},
	})
}
