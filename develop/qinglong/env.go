package qinglong

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cdle/sillyGirl/core"
)

type EnvResponse struct {
	Code int   `json:"code"`
	Data []Env `json:"data"`
}

type Env struct {
	Value     string `json:"value,omitempty"`
	ID        string `json:"_id,omitempty"`
	Status    int    `json:"status,omitempty"`
	Name      string `json:"name,omitempty"`
	Remarks   string `json:"remarks,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Created   int64  `json:"created,omitempty"`
}

func GetEnv(id string) (*Env, error) {
	envs, err := GetEnvs("")
	if err != nil {
		return nil, err
	}
	for _, env := range envs {
		if env.ID == id {
			return &env, nil
		}
	}
	return nil, nil
}

func GetEnvs(searchValue string) ([]Env, error) {
	er := EnvResponse{}
	if err := Config.Req(ENVS, &er, "?searchValue="+searchValue); err != nil {
		return nil, err
	}
	return er.Data, nil
}

func GetEnvss(searchValue string) ([]Env, error) {
	er := EnvResponse{}
	if err := Config.Req(ENVS, &er, "?searchValue="+searchValue); err != nil {
		return nil, err
	}
	return er.Data, nil
}

func SetEnv(e Env) error {
	envs, err := GetEnvs("")
	if err != nil {
		return err
	}
	for _, env := range envs {
		if env.Name == e.Name {
			if e.Remarks != "" {
				env.Remarks = e.Remarks
			}
			if e.Value != "" {
				env.Value = e.Value
			}
			if e.Name != "" {
				env.Name = e.Name
			}
			return Config.Req(PUT, ENVS, env)
		}
	}
	return AddEnv(e)
}

func UdpEnv(env Env) error {
	env.Created = 0
	env.Timestamp = ""
	return Config.Req(PUT, ENVS, env)
}

func ModEnv(e Env) error {
	envs, err := GetEnvs("")
	if err != nil {
		return err
	}
	for _, env := range envs {
		if env.ID == e.ID {
			if e.Remarks != "" {
				env.Remarks = e.Remarks
			}
			if e.Value != "" {
				env.Value = e.Value
			}
			if e.Name != "" {
				env.Name = e.Name
			}
			env.Created = 0
			env.Timestamp = ""
			return Config.Req(PUT, ENVS, env)
		}
	}
	return errors.New("找不到环境变量")
}

func AddEnv(e Env) error {
	e.Created = 0
	e.Timestamp = ""
	return Config.Req(POST, ENVS, []Env{e})
}

func RemEnv(e *Env) error {
	return Config.Req(DELETE, ENVS, []byte(`["`+e.ID+`"]`))
}

func initEnv() {
	core.AddCommand("ql", []core.Function{
		{
			Rules: []string{`cookie status`},
			Admin: true,
			Handle: func(_ core.Sender) interface{} {
				type Count struct {
					Total        int
					Disable      int
					TodayCreate  int
					TodayDisable int
					TodayUpdate  int
				}
				envs, err := GetEnvs("")
				if err != nil {
					return err
				}
				today := time.Now()
				var cookies = map[string]*Count{}
				for _, env := range envs {
					var c *Count
					if _, ok := cookies[env.Name]; !ok {
						cookies[env.Name] = &Count{}
					}
					c = cookies[env.Name]
					c.Total++
					if strings.Contains(env.Timestamp, fmt.Sprintf(`%s %s`, today.Month().String()[0:3], today.Format("02 2006"))) {
						if env.Status != 0 {
							c.TodayDisable++
						} else {
							c.TodayUpdate++
						}
					}
					if env.Status != 0 {
						c.Disable++
					}
					if time.Unix(env.Created, 0).Format("2006-01-02") == today.Format("2006-01-02") {
						c.TodayCreate++
					}
				}
				ss := []string{}
				for name, c := range cookies {
					ss = append(ss, fmt.Sprintf(`%s 今日新增%d，今日更新%d，今日失效%d，总数%d，有效%d，无效%d`, name, c.TodayCreate, c.TodayUpdate, c.TodayDisable, c.Total, c.Total-c.Disable, c.Disable))
				}
				return strings.Join(ss, "\n")
			},
		},
	})
}
