package qinglong

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/cdle/sillyGirl/core"
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
			Handle: func(_ core.Sender) interface{} {
				crons, err := GetCrons("")
				if err != nil {
					return err
				}
				if len(crons) == 0 {
					return "没有任务。"
				}
				es := []string{}
				for _, cron := range crons {
					es = append(es, formatCron(&cron))
				}
				return strings.Join(es, "\n\n")
			},
		},
		{
			Rules: []string{`cron status ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				keyword := s.Get()
				crons, err := GetCrons("")
				if err != nil {
					return err
				}
				es := []string{}
				for _, cron := range crons {
					if cron.ID == keyword {
						es = append(es, formatCron(&cron))
						break
					}
					if regexp.MustCompile(keyword+"$").FindString(cron.Command) != "" {
						es = append(es, formatCron(&cron))
					}
				}
				if len(es) == 0 {
					return "找不到任务。"
				}
				return strings.Join(es, "\n\n")
			},
		},
		{
			Rules: []string{`cron run ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				cron, err := GetCronID(s.Get())
				if err != nil {
					return err
				}
				if err := Req(CRONS, PUT, "/run", []byte(fmt.Sprintf(`["%s"]`, cron.ID))); err != nil {
					return err
				}
				return fmt.Sprintf("已运行 %v", cron.Name)
			},
		},
		{
			Rules: []string{`cron stop ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				cron, err := GetCronID(s.Get())
				if err != nil {
					return err
				}
				if err := Req(CRONS, PUT, "/stop", []byte(fmt.Sprintf(`["%s"]`, cron.ID))); err != nil {
					return err
				}
				return "操作成功"
			},
		},
		{
			Rules: []string{`cron enable ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				cron, err := GetCronID(s.Get())
				if err != nil {
					return err
				}
				if err := Req(CRONS, PUT, "/enable", []byte(fmt.Sprintf(`["%s"]`, cron.ID))); err != nil {
					return err
				}
				return "操作成功"
			},
		},
		{
			Rules: []string{`cron disable ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				cron, err := GetCronID(s.Get())
				if err != nil {
					return err
				}
				if err := Req(CRONS, PUT, "/disable", []byte(fmt.Sprintf(`["%s"]`, cron.ID))); err != nil {
					return err
				}
				return "操作成功"
			},
		},
		{
			Rules: []string{`cron find ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
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
			Handle: func(s core.Sender) interface{} {
				cron, err := GetCronID(s.Get())
				if err != nil {
					return err
				}
				data, err := GetCronLog(cron.ID)
				if err != nil {
					return err
				}
				return data
			},
		},
		{
			Rules: []string{`update`},
			Admin: true,
			Handle: func(_ core.Sender) interface{} {
				cron, err := GetCronID("更新面板")
				if err != nil {
					return err
				}
				if err := Req(CRONS, PUT, "/run", []byte(fmt.Sprintf(`["%s"]`, cron.ID))); err != nil {
					return err
				}
				return "操作成功"
			},
		},
		{
			Rules: []string{`update logs`},
			Admin: true,
			Handle: func(_ core.Sender) interface{} {
				cron, err := GetCronID("更新面板")
				if err != nil {
					return err
				}
				data, err := GetCronLog(cron.ID)
				if err != nil {
					return err
				}
				return data
			},
		},
		{
			Rules: []string{`cron hide duplicate`},
			Admin: true,
			Cron:  "*/5 * * * *",
			Handle: func(_ core.Sender) interface{} {
				w := func(s string) int {
					if strings.Contains(s, "shufflewzc") {
						return 1
					}
					if strings.Contains(s, "novpx") {
						return 10
					}
					if strings.Contains(s, "smiek2221") {
						return 9
					}
					if strings.Contains(s, "Aaron-lv") {
						return 8
					}
					return 0
				}
				crons, err := GetCrons("")
				if err != nil {
					return err
				}
				tasks := map[string]Cron{}
				for i := range crons {
					if crons[i].IsDisabled != 0 {
						continue
					}
					if task, ok := tasks[crons[i].Name]; ok {
						var dup Cron
						if w(task.Command) > w(crons[i].Command) {
							dup = crons[i]
						} else {
							dup = task
							tasks[crons[i].Name] = crons[i]
						}
						if err := Req(CRONS, PUT, "/disable", []byte(fmt.Sprintf(`["%s"]`, dup.ID))); err != nil {
							// s.Reply(fmt.Sprintf("隐藏 %v %v %v", dup.Name, dup.Command, err))
						} else {
							// s.Reply(fmt.Sprintf("已隐藏重复任务 %v %v", dup.Name, dup.Command))
						}
					} else {
						tasks[crons[i].Name] = crons[i]
					}
				}
				return "操作成功"
			},
		},
	})
}

func GetCrons(searchValue string) ([]Cron, error) {
	er := CronResponse{}
	if err := Req(CRONS, &er, "?searchValue="+searchValue); err != nil {
		return nil, err
	}
	return er.Data, nil
}

func GetCronLog(id string) (string, error) {
	c := &Carrier{
		Get: "data",
	}
	if err := Req(CRONS, "/"+id+"/log", c); err != nil {
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

func GetCronID(keyword string) (*Cron, error) {
	crons, err := GetCrons("")
	if err != nil {
		return nil, err
	}
	cs := []Cron{}
	for _, cron := range crons {
		if cron.IsDisabled != 0 {
			continue
		}
		if cron.ID == keyword {
			cs = append(cs, cron)
			break
		}
		if strings.Contains(cron.Name, keyword) {
			cs = append(cs, cron)
		}
		if regexp.MustCompile(keyword+"$").FindString(cron.Command) != "" {
			cs = append(cs, cron)
		}
	}
	if len(cs) == 0 {
		return nil, errors.New("找不到任务。")
	}
	if len(cs) != 1 {
		return nil, errors.New("搜索到多个任务，请再具体一些！")
	}
	return &cs[0], nil
}
