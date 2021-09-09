package qinglong

import (
	"errors"
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
	core.AddCommand("ql", []core.Function{
		{
			Rules: []string{`envs`},
			Admin: true,
			Handle: func(_ im.Sender) interface{} {
				envs, err := GetEnvs("")
				if err != nil {
					return err
				}
				if len(envs) == 0 {
					return "未设置任何环境变量"
				}
				es := []string{}
				for _, env := range envs {
					es = append(es, formatEnv(&env))
				}
				return strings.Join(es, "\n\n")
			},
		},
		{
			Rules: []string{`env get ?`},
			Admin: true,
			Handle: func(s im.Sender) interface{} {
				name := s.Get()
				envs, err := GetEnvs(name)
				if err != nil {
					return err
				}
				if len(envs) == 0 {
					return "未设置该环境变量"
				}
				es := []string{}
				for _, env := range envs {
					if env.Name == name {
						es = append(es, formatEnv(&env))
					}
				}
				return strings.Join(es, "\n\n")
			},
		},
		{
			Rules: []string{`env find ?`},
			Admin: true,
			Handle: func(s im.Sender) interface{} {
				m := s.Get()
				envs, err := GetEnvs(m)
				if err != nil {
					return err
				}
				if len(envs) == 0 {
					return "找不到环境变量"
				}
				es := []string{}
				for _, env := range envs {
					es = append(es, formatEnv(&env))
				}
				return strings.Join(es, "\n\n")
			},
		},
		{
			Rules: []string{`env set ? ?`},
			Admin: true,
			Handle: func(s im.Sender) interface{} {
				err := SetEnv(Env{
					Name:  s.Get(0),
					Value: s.Get(1),
				})
				if err != nil {
					return err
				}
				return fmt.Sprintf("操作成功")
			},
		},
		{
			Rules: []string{`env del ?`},
			Admin: true,
			Handle: func(s im.Sender) interface{} {
				if err := RemEnv(&Env{ID: s.Get()}); err != nil {
					return err
				}
				return "操作成功"
			},
		},
		{
			Rules: []string{`env remark ? ?`},
			Admin: true,
			Handle: func(s im.Sender) interface{} {
				if err := ModEnv(&Env{ID: s.Get(0), Remarks: s.Get(1)}); err != nil {
					return err
				}
				return "操作成功"
			},
		},
	})
}

func GetEnv(name string) (*Env, error) {
	envs, err := GetEnvs("")
	if err != nil {
		return nil, err
	}
	for _, env := range envs {
		if env.Name == name {
			return &env, nil
		}
	}
	return nil, nil
}

// func GetEnvs(searchValue string) ([]Env, error) {
// 	er := EnvResponse{}
// 	if err := req(ENVS, &er, "?searchValue="+searchValue); err != nil {
// 		return nil, err
// 	}
// 	return er.Data, nil
// }

// func SetEnv(e *Env) error {
// 	if e.Name == "JD_COOKIE" {
// 		return errors.New("不支持的操作")
// 	}
// 	envs, err := GetEnvs("")
// 	if err != nil {
// 		return err
// 	}
// 	for _, env := range envs {
// 		if env.Name == e.Name {
// 			if e.Remarks != "" {
// 				env.Remarks = e.Remarks
// 			}
// 			if e.Value != "" {
// 				env.Value = e.Value
// 			}
// 			if e.Name != "" {
// 				env.Name = e.Name
// 			}
// 			env.Timestamp = ""
// 			return req(PUT, ENVS, env)
// 		}
// 	}
// 	return AddEnv(e)
// }

func ModEnv(e *Env) error {
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
			env.Timestamp = ""
			return req(PUT, ENVS, env)
		}
	}
	return errors.New("找不到环境变量")
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
	return fmt.Sprintf("名称：%v\n备注：%v\n状态：%v\n值：%v", env.Name, env.Remarks, status, env.Value)
}
