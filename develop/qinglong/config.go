package qinglong

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cdle/sillyGirl/core"
)

func initConfig() {
	core.AddCommand("ql", []core.Function{
		{
			Rules: []string{`config`},
			Admin: true,
			Handle: func(_ core.Sender) interface{} {
				config, err := GetConfig()
				if err != nil {
					return err
				}
				return config
			},
		},
	})
	core.AddCommand("ql", []core.Function{
		{
			Rules: []string{`envs`},
			Admin: true,
			Handle: func(_ core.Sender) interface{} {
				envs, err := GetConfigEnvs("")
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
			Handle: func(s core.Sender) interface{} {
				name := s.Get()
				envs, err := GetConfigEnvs(name)
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
			Handle: func(s core.Sender) interface{} {
				m := s.Get()
				envs, err := GetConfigEnvs(m)
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
			Handle: func(s core.Sender) interface{} {
				err := SetConfigEnv(Env{
					Name:   s.Get(0),
					Value:  s.Get(1),
					Status: 3,
				})
				if err != nil {
					return err
				}
				return fmt.Sprintf("操作成功")
			},
		},
		// {
		// 	Rules: []string{`env delete ?`},
		// 	Admin: true,
		// 	Handle: func(s core.Sender) interface{} {
		// 		if err := DelEnv(&Env{ID: s.Get()}); err != nil {
		// 			return err
		// 		}
		// 		return "操作成功"
		// 	},
		// },
		{
			Rules: []string{`env remark ? ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				if err := SetConfigEnv(Env{Name: s.Get(0), Remarks: s.Get(1), Status: 3}); err != nil {
					return err
				}
				return "操作成功"
			},
		},
		{
			Rules: []string{`env disable ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				if err := SetConfigEnv(Env{Name: s.Get(), Status: 1}); err != nil {
					return err
				}
				return "操作成功"
			},
		},
		{
			Rules: []string{`env enable ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				if err := SetConfigEnv(Env{Name: s.Get()}); err != nil {
					return err
				}
				return "操作成功"
			},
		},
	})
}

func GetConfig() (string, error) {
	config := "data"
	if err := Config.Req(CONFIG, &config, "/config.sh"); err != nil {
		return "", err
	}
	return config, nil
}

func SvaeConfig(content string) error {
	if err := Config.Req(POST, CONFIG, map[string]interface{}{
		"name":    "config.sh",
		"content": content,
	}, "/save"); err != nil {
		return err
	}
	return nil
}

func GetConfigEnvs(searchValue string) ([]Env, error) {
	envs := []Env{}
	content, err := GetConfig()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if searchValue != "" && !strings.Contains(line, searchValue) {
			continue
		}
		for i, pattern := range []string{`^\s*export\s+([^'"=\s]+)=[ '"]?(.*?)['"]?$`, `^\s*#[#\s]*export\s+([^'"=\s]+)=[ '"]?(.*?)['"]?$`, `^\s*([^'"=\s]+)=[ '"]?(.*?)['"]?$`, `^\s*#[#\s]*([^'"=\s]+)=[ '"]?(.*?)['"]?$`} {
			if v := regexp.MustCompile(pattern).FindStringSubmatch(line); len(v) > 0 {
				e := Env{}
				if i == 1 || i == 3 {
					e.Status = 1
				}
				e.Name = v[1]
				e.Value = v[2]
				envs = append(envs, e)
				break
			}
		}
	}
	return envs, nil
}

func SetConfigEnv(envs ...Env) error {
	config, err := GetConfig()
	if err != nil {
		return err
	}
	lines := strings.Split(config, "\n")
	for _, env := range envs {
		if env.Name == "" {
			continue
		}
		set := false
		for j, line := range lines {
			for i, pattern := range []string{`^\s*export\s+([^'"=\s]+)=[ '"]?(.*?)['"]?$`, `^\s*#[#\s]*export\s+([^'"=\s]+)=[ '"]?(.*?)['"]?$`, `^\s*([^'"=\s]+)=[ '"]?(.*?)['"]?$`, `^\s*#[#\s]*([^'"=\s]+)=[ '"]?(.*?)['"]?$`} {
				if v := regexp.MustCompile(pattern).FindStringSubmatch(line); len(v) > 0 {
					e := Env{}
					if i == 1 || i == 3 {
						e.Status = 1
					}
					e.Name = v[1]
					e.Value = v[2]
					if env.Name != e.Name {
						continue
					}
					if env.Value != "" && env.Value != e.Value {
						e.Value = env.Value
					}
					if env.Status < 2 {
						e.Status = env.Status
					}
					h := ""
					if e.Status == 1 {
						h = "## "
					}
					if i <= 1 {
						h += "export "
					}
					lines[j] = h + fmt.Sprintf("%s=\"%s\"", e.Name, e.Value)
					set = true
					break
				}
			}
		}
		if !set {
			lines = append(lines, fmt.Sprintf("export %s=\"%s\"", env.Name, env.Value))
		}
	}
	return SvaeConfig(strings.Join(lines, "\n"))
}

func formatEnv(env *Env) string {
	status := "已启用"
	if env.Status != 0 {
		status = "已禁用"
	}
	return fmt.Sprintf("名称：%v\n状态：%v\n值：%v", env.Name, status, env.Value)
}
