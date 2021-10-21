package qinglong

import (
	"fmt"
	"strings"
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
				for {
					time.Sleep(time.Microsecond * 300)
					data, _ := GetCronLog(cron.Value)
					if strings.Contains(data, "执行结束...") {
						s.Reply(data)
						break
					}
				}
				if err := Config.Req(CRONS, DELETE, []byte(`["`+cron.Value+`"]`)); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Rules: []string{`task ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				cron := &Carrier{
					Get: "data._id",
				}
				if err := Config.Req(cron, CRONS, POST, []byte(`{"name":"sillyGirl临时创建任务","command":"task `+s.Get()+`","schedule":" 1 1 1 1 1"}`)); err != nil {
					return err
				}
				if err := Config.Req(CRONS, PUT, "/run", []byte(fmt.Sprintf(`["%s"]`, cron.Value))); err != nil {
					return err
				}
				for {
					time.Sleep(time.Second)
					data, _ := GetCronLog(cron.Value)
					if strings.Contains(data, "执行结束...") {
						s.Reply(data)
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
			Rules: []string{`repo ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				cron := &Carrier{
					Get: "data._id",
				}
				if err := Config.Req(cron, CRONS, POST, []byte(`{"name":"sillyGirl临时创建任务","command":"task `+s.Get()+`","schedule":" 1 1 1 1 1"}`)); err != nil {
					return err
				}
				if err := Config.Req(CRONS, PUT, "/run", []byte(fmt.Sprintf(`["%s"]`, cron.Value))); err != nil {
					return err
				}
				for {
					time.Sleep(time.Second)
					data, _ := GetCronLog(cron.Value)
					if strings.Contains(data, "执行结束...") {
						s.Reply(data)
						break
					}
				}
				if err := Config.Req(cron, CRONS, DELETE, []byte(`["`+cron.Value+`"]`)); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
