package qinglong

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cdle/sillyGirl/core"
	"github.com/cdle/sillyGirl/im"
)

type CronResponse struct {
	Code int    `json:"code"`
	Data []Cron `json:"data"`
}
type Cron struct {
	Name       string      `json:"name"`
	Command    string      `json:"command"`
	Schedule   string      `json:"schedule"`
	Saved      bool        `json:"saved"`
	ID         string      `json:"_id"`
	Created    int64       `json:"created"`
	Status     int         `json:"status"`
	Timestamp  string      `json:"timestamp"`
	IsSystem   int         `json:"isSystem"`
	IsDisabled int         `json:"isDisabled"`
	LogPath    string      `json:"log_path"`
	Pid        interface{} `json:"pid"`
}

func init() {
	core.AddCommand("ql", []core.Function{
		{
			Rules: []string{`crons`},
			Admin: true,
			Handle: func(_ im.Sender) interface{} {
				crons, err := GetCrons("")
				if err != nil {
					return err
				}
				if len(crons) == 0 {
					return "没有任务"
				}
				es := []string{}
				for _, cron := range crons {
					es = append(es, formatCron(&cron))
				}
				return strings.Join(es, "\n\n")
			},
		},
		{
			Rules: []string{`cron get ?`},
			Admin: true,
			Handle: func(s im.Sender) interface{} {
				name := s.Get()
				crons, err := GetCrons("")
				if err != nil {
					return err
				}
				es := []string{}
				for _, cron := range crons {
					if cron.Name == name || regexp.MustCompile(name+"$").FindString(cron.Command) != "" {
						es = append(es, formatCron(&cron))
					}
				}
				if len(es) == 0 {
					return "找不到该任务"
				}
				return strings.Join(es, "\n\n")
			},
		},
		{
			Rules: []string{`cron find ?`},
			Admin: true,
			Handle: func(s im.Sender) interface{} {
				name := s.Get()
				crons, err := GetCrons("")
				if err != nil {
					return err
				}
				es := []string{}
				for _, cron := range crons {
					if strings.Contains(cron.Name, name) || strings.Contains(cron.Command, name) {
						es = append(es, formatCron(&cron))
					}
				}
				if len(es) == 0 {
					return "找不到匹配的任务"
				}
				return strings.Join(es, "\n\n")
			},
		},
		{
			Rules: []string{`cron logs ?`},
			Admin: true,
			Handle: func(s im.Sender) interface{} {
				data, err := GetCronLog(s.Get())
				if err != nil {
					return err
				}
				return data
			},
		},
	})
}

func GetCrons(searchValue string) ([]Cron, error) {
	er := CronResponse{}
	if err := req(CRONS, &er, "?searchValue="+searchValue); err != nil {
		return nil, err
	}
	return er.Data, nil
}

func GetCronLog(id string) (string, error) {
	c := &Carrier{
		Get: "data",
	}
	if err := req(CRONS, "/"+id+"/log", c); err != nil {
		return "", err
	}
	return c.Value, nil
}

func formatCron(cron *Cron) string {
	status := "空闲中"
	if cron.IsDisabled != 0 {
		status = "已禁用"
	}
	if cron.Pid != nil && fmt.Sprint(cron.Pid) != "" {
		status = "运行中"
	}
	return strings.Join([]string{
		fmt.Sprintf("任务名：%v", cron.Name),
		fmt.Sprintf("编号：%v", cron.ID),
		fmt.Sprintf("命令：%v", cron.Command),
		fmt.Sprintf("定时：%v", cron.Schedule),
		fmt.Sprintf("状态：%v", status),
	}, "\n")
}
