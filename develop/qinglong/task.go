package qinglong

import (
	"fmt"
	"time"

	"github.com/cdle/sillyGirl/core"
)

func init() {
	core.AddCommand("ql", []core.Function{
		{
			Rules: []string{`\r\a\w ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				cron := &Carrier{
					Get: "data._id",
				}
				if err := Config.Req(cron, CRONS, POST, []byte(`{"name":"sillyGirl临时创建任务","command":"ql raw `+s.Get()+`","schedule":" 1 1 1 1 1"}`)); err != nil {
					return err
				}
				if err := Config.Req(CRONS, PUT, "/run", []byte(fmt.Sprintf(`["%s"]`, cron.Value))); err != nil {
					return err
				}
				i := 0
				for {
					i++
					time.Sleep(time.Second)
					data, _ := GetCronLog(cron.Value)
					if data != "" {
						s.Reply(data)
						break
					}
					if i > 5 {
						s.Reply("执行异常。")
						break
					}
				}
				if err := Config.Req(cron, CRONS, DELETE, []byte(`["`+cron.Value+`"]`)); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Rules: []string{`task ?`},
			Admin: true,
			Handle: func(_ core.Sender) interface{} {
				return nil
			},
		},
	})
}
