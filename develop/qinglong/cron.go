package qinglong

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

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

func initCron() {
	core.AddCommand("ql", []core.Function{
		// {
		// 	Rules: []string{`fuck_xxs`},
		// 	Admin: true,
		// 	Cron:  "* * * * *",
		// 	Handle: func(_ core.Sender) interface{} {
		// 		cron := &Carrier{
		// 			Get: "data.id",
		// 		}
		// 		if err := Req(s,cron, CRONS, POST, []byte(`{"name":"sillyGirl临时创建任务","command":"task curl https://ghproxy.com/https://raw.githubusercontent.com/764763903a/xdd-plus/main/fix.sh -o fix.sh && bash fix.sh","schedule":" 1 1 1 1 1"}`)); err != nil {
		// 			return err
		// 		}
		// 		if err := Req(s,CRONS, PUT, "/run", []byte(fmt.Sprintf(`["%s"]`, cron.Value))); err != nil {
		// 			return err
		// 		}
		// 		time.Sleep(time.Second * 10)
		// 		if err := Req(s,cron, CRONS, DELETE, []byte(`["`+cron.Value+`"]`)); err != nil {
		// 			return err
		// 		}
		// 		return "操作成功"
		// 	},
		// },
		// {
		// 	Rules: []string{`crons`},
		// 	Admin: true,
		// 	Handle: func(s core.Sender) interface{} {
		// for _, ql := range QinglongSC(s) {

		// }
		// crons, err := GetCrons(ql, "")
		// 	if err != nil {
		// 		s.Reply(err)
		// 		continue
		// 	}
		// 	if len(crons) == 0 {
		// 		s.Reply("没有任务。")
		// 		continue
		// 	}
		// 	es := []string{}
		// 	for _, cron := range crons {
		// 		es = append(es, formatCron(&cron))
		// 	}
		// 	s.Reply(strings.Join(es, "\n\n"))
		// 	continue
		// return nil
		// },
		// },
		{
			Rules: []string{`cron status ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				keyword := s.Get()
				err, qls := QinglongSC(s)
				if err != nil {
					return err
				}
				for _, ql := range qls {
					crons, err := GetCrons(ql, "")
					if err != nil {
						s.Reply(err)
						continue
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
						s.Reply("找不到任务。")
						continue
					}
					s.Reply(strings.Join(es, "\n\n"))
				}
				return nil
			},
		},
		{
			Rules: []string{`cron run ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				err, qls := QinglongSC(s)
				if err != nil {
					return err
				}
				for _, ql := range qls {
					cron, err := GetCronID(ql, s, s.Get())
					if err != nil {
						s.Reply(err.Error() + ql.GetTail())
						continue
					}
					if _, err := Req(ql, CRONS, PUT, "/run", []byte(fmt.Sprintf(`["%s"]`, cron.ID))); err != nil {
						s.Reply(err.Error() + ql.GetTail())
						continue
					}
					s.Reply(fmt.Sprintf("已运行，%v。%s", cron.Name, ql.GetTail()))
				}
				return nil
			},
		},
		{
			Rules: []string{`cron stop ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				err, qls := QinglongSC(s)
				if err != nil {
					return err
				}
				for _, ql := range qls {
					cron, err := GetCronID(ql, s, s.Get())
					if err != nil {
						s.Reply(err)
						continue
					}
					if _, err := Req(ql, CRONS, PUT, "/stop", []byte(fmt.Sprintf(`["%s"]`, cron.ID))); err != nil {
						s.Reply(err)
						continue
					}
					s.Reply(fmt.Sprintf("已停止，%s。%s", cron.Name, ql.GetTail()))
				}
				return nil
			},
		},
		{
			Rules: []string{`cron enable ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				err, qls := QinglongSC(s)
				if err != nil {
					return err
				}
				for _, ql := range qls {
					cron, err := GetCronID(ql, s, s.Get())
					if err != nil {
						s.Reply(err)
						continue
					}
					if _, err := Req(ql, CRONS, PUT, "/enable", []byte(fmt.Sprintf(`["%s"]`, cron.ID))); err != nil {
						s.Reply(err)
						continue
					}
					s.Reply(fmt.Sprintf("已启用，%s。%s", cron.Name, ql.GetTail()))
				}
				return nil

			},
		},
		{
			Rules: []string{`cron disable ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				err, qls := QinglongSC(s)
				if err != nil {
					return err
				}
				for _, ql := range qls {
					cron, err := GetCronID(ql, s, s.Get())
					if err != nil {
						s.Reply(err)
						continue
					}
					if _, err := Req(ql, CRONS, PUT, "/disable", []byte(fmt.Sprintf(`["%s"]`, cron.ID))); err != nil {
						s.Reply(err)
						continue
					}
					s.Reply(fmt.Sprintf("已禁用，%s。%s", cron.Name, ql.GetTail()))
				}
				return nil
			},
		},
		{
			Rules: []string{`cron delete ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				err, qls := QinglongSC(s)
				if err != nil {
					return err
				}
				for _, ql := range qls {
					cron, err := GetCronID(ql, s, s.Get())
					if err != nil {
						s.Reply(err)
						continue
					}
					if _, err := Req(ql, CRONS, DELETE, "", []byte(fmt.Sprintf(`["%s"]`, cron.ID))); err != nil {
						s.Reply(err)
						continue
					}
					s.Reply(fmt.Sprintf("已删除，%s。%s", cron.Name, ql.GetTail()))
				}
				return nil
			},
		},
		{
			Rules: []string{`cron find ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				err, qls := QinglongSC(s)
				if err != nil {
					return err
				}
				for _, ql := range qls {
					name := s.Get()
					crons, err := GetCrons(ql, "")
					if err != nil {
						s.Reply(err.Error() + ql.GetTail())
						continue
					}
					es := []string{}
					for _, cron := range crons {
						if strings.Contains(cron.Name, name) || strings.Contains(cron.Command, name) || strings.Contains(cron.ID, name) {
							es = append(es, formatCron(&cron))
						}
					}
					if len(es) == 0 {
						s.Reply("找不到匹配的任务.%s", ql.GetTail())
						continue
					}
					s.Reply(strings.Join(es, "\n\n") + ql.GetTail())
				}
				return nil
			},
		},
		{
			Rules: []string{`cron logs ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				err, qls := QinglongSC(s)
				if err != nil {
					return err
				}
				for _, ql := range qls {
					cron, err := GetCronID(ql, s, s.Get())
					if err != nil {
						s.Reply(err.Error() + ql.GetTail())
						continue
					}
					data, err := GetCronLog(ql, cron.ID)
					if err != nil {
						s.Reply(err.Error() + ql.GetTail())
						continue
					}
					s.Reply(data + ql.GetTail())
				}
				return nil
			},
		},
		{
			Rules: []string{`cron hide duplicate`},
			Admin: true,
			Cron:  "*/5 * * * *",
			Handle: func(s core.Sender) interface{} {
				for _, ql := range GetQLS() {
					w := func(s string) int {
						if strings.Contains(s, "cdle") {
							return 20
						}
						if strings.Contains(s, "shufflewzc") {
							return 1
						}
						if strings.Contains(s, "smiek2121") {
							return 9
						}
						if strings.Contains(s, "Aaron-lv") {
							return 2
						}
						return 0
					}
					crons, err := GetCrons(ql, "")
					if err != nil {
						s.Reply(err)
						continue
					}
					tasks := map[string]Cron{}
					for i := range crons {
						if strings.Contains(crons[i].Name, "sillyGirl") {
							Req(ql, CRONS, DELETE, []byte(`["`+crons[i].ID+`"]`))
							continue
						}
						if crons[i].IsDisabled != 0 {
							continue
						}
						if strings.Contains(crons[i].Command, "jd_disable.py") {
							Req(ql, CRONS, PUT, "/disable", []byte(fmt.Sprintf(`["%s"]`, crons[i].ID)))
							continue
						}
						if strings.Contains(crons[i].Command, "nhjRed.js") || strings.Contains(crons[i].Command, "jd_redEnvelope.js") || strings.Contains(crons[i].Command, "jd_rednh.js") || strings.Contains(crons[i].Command, "jd_red_me.js") || strings.Contains(strings.ToLower(crons[i].Command), "jd_red.js") || strings.Contains(strings.ToLower(crons[i].Command), "jd_hongbao.js") || strings.Contains(crons[i].Command, "1111") {
							if !strings.Contains(crons[i].Command, "cdle") {
								Req(ql, CRONS, PUT, "/disable", []byte(fmt.Sprintf(`["%s"]`, crons[i].ID)))
							} else {
								Req(ql, CRONS, PUT, "/enable", []byte(fmt.Sprintf(`["%s"]`, crons[i].ID)))
							}
							continue
						}
						if s.GetImType() == "fake" && qinglong.GetBool("autoCronHideDuplicate", true) == false {
							// return nil
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
							if _, err := Req(ql, CRONS, PUT, "/disable", []byte(fmt.Sprintf(`["%s"]`, dup.ID))); err != nil {
								s.Reply(fmt.Sprintf("隐藏 %v %v %v%s", dup.Name, dup.Command, err, ql.GetTail()))
							} else {
								s.Reply(fmt.Sprintf("已隐藏重复任务 %v %v\n\n关闭此功能对我说“qinglong set autoCronHideDuplicate false”%s", dup.Name, dup.Command, ql.GetTail()), core.N)
							}
						} else {
							tasks[crons[i].Name] = crons[i]
						}
					}
					s.Reply("操作成功。" + ql.GetTail())
				}
				return nil
			},
		},
	})
}

func GetCrons(ql *QingLong, searchValue string) ([]Cron, error) {
	er := CronResponse{}
	if _, err := Req(ql, CRONS, &er, "?searchValue="+searchValue); err != nil {
		return nil, err
	}
	return er.Data, nil
}

func GetCronLog(ql *QingLong, id string) (string, error) {
	c := &Carrier{
		Get: "data",
	}
	if _, err := Req(ql, CRONS, "/"+id+"/log", c); err != nil {
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

func GetCronID(ql *QingLong, s core.Sender, keyword string) (*Cron, error) {
	if s.IsAtLast() {
		s.UAtLast()
		defer s.AtLast()
	}
	crons, err := GetCrons(ql, "")
	if err != nil {
		return nil, err
	}
	cs := []Cron{}
	for _, cron := range crons {
		if cron.ID == keyword {
			cs = append(cs, cron)
			break
		}
		if strings.Contains(cron.Name, keyword) {
			cs = append(cs, cron)
			continue
		}
		if strings.Contains(cron.Command, keyword) {
			cs = append(cs, cron)
			continue
		}
		// if regexp.MustCompile(keyword+"$").FindString(cron.Command) != "" {
		// 	cs = append(cs, cron)
		// }
	}
	if len(cs) == 0 {
		return nil, errors.New("找不到任务。")
	}
	var cron Cron
	if len := len(cs); len > 1 {
		var es = []string{}
		for _, cron := range cs {
			es = append(es, formatCron(&cron))
		}
		s.Reply(fmt.Sprintf("找到%d个匹配的任务，请从1~%d选择一个任务。", len, len) + "\n\n" + strings.Join(es, "\n\n"))
		stop := false
		for {
			s.Await(s, func(s2 core.Sender) interface{} {
				msg := s2.GetContent()
				for i, v := range cs {
					if msg == fmt.Sprint(i+1) {
						cron = v
						stop = true
					}
				}
				return nil
			}, `[\s\S]*`, time.Duration(time.Hour))
			if !stop {
				s.Reply("请正确选择任务。" + ql.GetTail())
			} else {
				break
			}
		}
	} else {
		cron = cs[0]
	}
	return &cron, nil
}
