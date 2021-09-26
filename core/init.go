package core

import (
	"bufio"
	"os"
	"regexp"
	"time"
)

var Duration time.Duration

func init() {
	killp()
	for _, arg := range os.Args {
		if arg == "-d" {
			initStore()
			Daemon()
		}
	}
	_, err := os.Stat("/etc/sillyGirl/")
	if err != nil {
		os.MkdirAll("/etc/sillyGirl/", os.ModePerm)
	}
	initStore()
	ReadYaml(ExecPath+"/conf/", &Config, "https://raw.githubusercontent.com/cdle/sillyGirl/main/conf/demo_config.yaml")
	InitReplies()
	initToHandleMessage()
	file, err := os.Open("/etc/sillyGirl/sets.conf")
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if v := regexp.MustCompile(`^\s*set\s+(\S+)\s+(\S+)\s+(\S+)`).FindStringSubmatch(line); len(v) > 0 {
				b := Bucket(v[1])
				if b.Get(v[2]) != v[3] {
					b.Set(v[2], v[3])
				}
			}
		}
		file.Close()
	}
	initSys()
	Duration = time.Duration(sillyGirl.GetInt("duration", 5)) * time.Second
}
