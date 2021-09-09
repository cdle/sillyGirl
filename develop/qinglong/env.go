package qinglong

import (
	"errors"
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

func GetEnvs(searchValue string) ([]Env, error) {
	er := EnvResponse{}
	if err := req(ENVS, &er, "?searchValue="+searchValue); err != nil {
		return nil, err
	}
	return er.Data, nil
}

func SetEnv(e *Env) error {
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
			env.Timestamp = ""
			return req(PUT, ENVS, env)
		}
	}
	return AddEnv(e)
}

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
