package qinglong

import (
	"errors"
)

type EnvResponse struct {
	Code int   `json:"code"`
	Data []Env `json:"data"`
}

type Env struct {
	Value   string `json:"value,omitempty"`
	ID      string `json:"_id,omitempty"`
	Status  int    `json:"status,omitempty"`
	Name    string `json:"name,omitempty"`
	Remarks string `json:"remarks,omitempty"`
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
			return Config.Req(PUT, ENVS, env)
		}
	}
	return errors.New("找不到环境变量")
}

func AddEnv(e Env) error {
	return Config.Req(POST, ENVS, []Env{e})
}

func RemEnv(e *Env) error {
	return Config.Req(DELETE, ENVS, []byte(`["`+e.ID+`"]`))
}
