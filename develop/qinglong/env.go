package qinglong

import (
	"fmt"
	"strings"

	"github.com/cdle/sillyGirl/core"
	"github.com/cdle/sillyGirl/im"
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
	Timestamp string `json:"timestamp,omitempty"`
	Remarks   string `json:"remarks,omitempty"`
}

func init() {
	core.AddCommand([]core.Function{
		{
			Rules: []string{`^env\s+get\s+([\S]*)$`},
			Handle: func(s im.Sender) interface{} {
				m := s.Get()
				env, err := GetEnv(m)
				if err != nil {
					return err
				}
				if env == nil {
					return "未设置该环境变量"
				}
				if env != nil {
					return formatEnv(env)
				}
				return nil
			},
		},
		{
			Rules: []string{`^env\s+find\s+([\S]*)$`},
			Handle: func(s im.Sender) interface{} {
				m := s.Get()
				envs, err := GetEnvs(m)
				if err != nil {
					return err
				}
				if len(envs) == 0 {
					return "未设置该环境变量"
				}
				es := []string{}
				for _, env := range envs {
					es = append(es, formatEnv(&env))
				}
				return strings.Join(es, "\n\n")
			},
		},
		{
			Rules: []string{`^export\s+([^'"=]+)=['"]?([^=]+?)['"]?$`, `^env\s+set\s+([^'"=]+)=['"]?([^=]+?)['"]?$`},
			Handle: func(s im.Sender) interface{} {
				e := &Env{
					Name:  s.Get(0),
					Value: s.Get(1),
				}
				err := SetEnv(e)
				if err != nil {
					return err
				}
				return fmt.Sprintf("操作成功")
			},
		},
		{
			Rules: []string{`^env\s+del\s+([\S]*)$`},
			Handle: func(s im.Sender) interface{} {
				if err := RemEnv(&Env{ID: s.Get()}); err != nil {
					return err
				}
				return "操作成功"
			},
		},
	})
}

func GetEnv(searchValue string) (*Env, error) {
	envs, err := GetEnvs(searchValue)
	if err != nil {
		return nil, err
	}

	if len(envs) == 0 {
		return nil, nil
	}
	return &envs[0], nil
}

func GetEnvs(searchValue string) ([]Env, error) {
	er := EnvResponse{}
	if err := req(ENVS, &er, "?searchValue="+searchValue); err != nil {
		return nil, err
	}
	return er.Data, nil
}

func SetEnv(e *Env) error {
	es, err := GetEnvs(e.Name)
	if err != nil {
		return err
	}
	if len(es) == 0 {
		return AddEnv(e)
	}
	e.ID = es[0].ID
	return req(PUT, ENVS, *e)
}

func AddEnv(e *Env) error {
	return req(POST, ENVS, []Env{*e})
}

func RemEnv(e *Env) error {
	return req(DELETE, ENVS, []byte(`["`+e.ID+`"]`))
}

func formatEnv(env *Env) string {
	status := "已启用"
	if env.Status != 0 {
		status = "已禁用"
	}
	if env.Remarks == "" {
		env.Remarks = "无"
	}
	return fmt.Sprintf("名称：%v\n编号：%v\n备注：%v\n状态：%v\n时间：%v\n值：%v", env.Name, env.ID, env.Remarks, status, env.Timestamp, env.Value)
}
