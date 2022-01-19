package qinglong

import (
	"errors"
	"fmt"
	"strings"

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
	PtPin     string `json:"-"`
	Index     int    `json:"-"`
}

func GetEnv(ql *QingLong, id string) (*Env, error) {
	envs, err := GetEnvs(ql, "")
	if err != nil {
		return nil, err
	}
	for _, env := range envs {
		if env.ID == id {
			return &env, nil
		}
	}
	return nil, errors.New("无效ID。")
}

func GetEnvs(ql *QingLong, searchValue string) ([]Env, error) {
	er := EnvResponse{}
	if _, err := Req(ql, ENVS, &er, "?searchValue="+searchValue); err != nil {
		return nil, err
	}
	return er.Data, nil
}

func GetEnvss(ql *QingLong, searchValue string) ([]Env, error) {
	er := EnvResponse{}
	if _, err := Req(ql, ENVS, &er, "?searchValue="+searchValue); err != nil {
		return nil, err
	}
	return er.Data, nil
}

// func SetEnv(e Env) error {
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
// 			return Req(nil, PUT, ENVS, env)
// 		}
// 	}
// 	return AddEnv(e)
// }

func UdpEnv(ql *QingLong, env Env) error {
	env.Created = 0
	env.Timestamp = ""
	_, err := Req(ql, PUT, ENVS, env)
	return err
}

// func ModEnv(e Env) error {
// 	envs, err := GetEnvs("")
// 	if err != nil {
// 		return err
// 	}
// 	for _, env := range envs {
// 		if env.ID == e.ID {
// 			if e.Remarks != "" {
// 				env.Remarks = e.Remarks
// 			}
// 			if e.Value != "" {
// 				env.Value = e.Value
// 			}
// 			if e.Name != "" {
// 				env.Name = e.Name
// 			}
// 			env.Created = 0
// 			env.Timestamp = ""
// 			return Req(nil, PUT, ENVS, env)
// 		}
// 	}
// 	return errors.New("找不到环境变量")
// }

func AddEnv(ql *QingLong, es ...Env) error {
	nn := []Env{}
	for i := range es {
		n := es[i]
		n.Created = 0
		n.Timestamp = ""
		n.ID = ""
		nn = append(nn, n)
	}
	_, err := Req(ql, POST, ENVS, nn)
	return err
}

func RemEnv(ql *QingLong, es ...Env) error {
	v := []string{}
	for i := range es {
		v = append(v, fmt.Sprintf(`"%s"`, es[i].ID))
	}
	_, err := Req(ql, DELETE, ENVS, []byte(`[`+strings.Join(v, ",")+`]`))
	return err
}

func DisableEnv(ql *QingLong, es ...Env) error {
	v := []string{}
	for i := range es {
		v = append(v, fmt.Sprintf(`"%s"`, es[i].ID))
	}
	_, err := Req(ql, PUT, ENVS, "/disable", []byte(`[`+strings.Join(v, ",")+`]`))
	return err
}

func DisableCron(ql *QingLong, es ...Cron) error {
	v := []string{}
	if ql.IsSqlite() {
		for i := range es {
			v = append(v, fmt.Sprintf(`%s`, es[i].ID))
		}
	} else {
		for i := range es {
			v = append(v, fmt.Sprintf(`"%s"`, es[i].ID))
		}
	}
	_, err := Req(ql, PUT, CRONS, "/disable", []byte(`[`+strings.Join(v, ",")+`]`))
	return err
}

func EnableCron(ql *QingLong, es ...Cron) error {
	v := []string{}
	if ql.IsSqlite() {
		for i := range es {
			v = append(v, fmt.Sprintf(`%s`, es[i].ID))
		}
	} else {
		for i := range es {
			v = append(v, fmt.Sprintf(`"%s"`, es[i].ID))
		}
	}
	_, err := Req(ql, PUT, CRONS, "/enable", []byte(`[`+strings.Join(v, ",")+`]`))
	return err
}

func initEnv() {
	core.AddCommand("ql", []core.Function{
		{
			Rules: []string{`cookie status`},
			Admin: true,
			Handle: func(_ core.Sender) interface{} {
				return "等待更新。"
				// type Count struct {
				// 	Total        int
				// 	Disable      int
				// 	TodayCreate  int
				// 	TodayDisable int
				// 	TodayUpdate  int
				// }
				// envs, err := GetEnvs("")
				// if err != nil {
				// 	return err
				// }
				// today := time.Now()
				// var cookies = map[string]*Count{}
				// for _, env := range envs {
				// 	var c *Count
				// 	if _, ok := cookies[env.Name]; !ok {
				// 		cookies[env.Name] = &Count{}
				// 	}
				// 	c = cookies[env.Name]
				// 	c.Total++
				// 	if strings.Contains(env.Timestamp, fmt.Sprintf(`%s %s`, today.Month().String()[0:3], today.Format("02 2006"))) {
				// 		if env.Status != 0 {
				// 			c.TodayDisable++
				// 		} else {
				// 			c.TodayUpdate++
				// 		}
				// 	}
				// 	if env.Status != 0 {
				// 		c.Disable++
				// 	}
				// 	if time.Unix(env.Created, 0).Format("2006-01-02") == today.Format("2006-01-02") {
				// 		c.TodayCreate++
				// 	}
				// }
				// ss := []string{}
				// for name, c := range cookies {
				// 	ss = append(ss, fmt.Sprintf(`%s 今日新增%d，今日更新%d，今日失效%d，总数%d，有效%d，无效%d`, name, c.TodayCreate, c.TodayUpdate, c.TodayDisable, c.Total, c.Total-c.Disable, c.Disable))
				// }
				// return strings.Join(ss, "\n")
			},
		},
	})
}
