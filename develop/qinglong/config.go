package qinglong

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cdle/sillyGirl/core"
	"github.com/cdle/sillyGirl/im"
)

func init() {
	core.AddCommand("ql", []core.Function{
		{
			Rules: []string{`config`},
			Admin: true,
			Handle: func(_ im.Sender) interface{} {
				config, err := GetConfig()
				if err != nil {
					return err
				}
				return config
			},
		},
	})
}

func GetConfig() (string, error) {
	config := "data"
	if err := req(CONFIG, &config, "/config.sh"); err != nil {
		return "", err
	}
	return config, nil
}

func SvaeConfig(content string) error {
	if err := req(POST, CONFIG, map[string]interface{}{
		"name":    "config.sh",
		"content": content,
	}, "/save"); err != nil {
		return err
	}
	return nil
}

func GetEnvs(searchValue string) ([]Env, error) {
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

func SetEnv(env Env) error {
	config, err := GetConfig()
	if err != nil {
		return err
	}
	lines := strings.Split(config, "\n")
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
				if env.Value != e.Value {
					e.Value = env.Value
				}
				if env.Status != e.Status {
					e.Status = env.Status
				}
				h := ""
				if e.Status == 1 {
					h = "# "
				}
				if i <= 1 {
					h = "export "
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
	return SvaeConfig(strings.Join(lines, "\n"))
}
