package qinglong

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
	es, err := GetEnvs(e.Name)
	if err != nil {
		return err
	}
	if len(es) == 0 {
		return nil
	}
	td := []string{}
	for _, e := range es {
		td = append(td, e.ID)
	}
	return req(DELETE, ENVS, td)
}
