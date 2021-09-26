package core

import (
	"bufio"
	"os"
	"regexp"
	"strings"
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
	initStore()
	ReadYaml(ExecPath+"/conf/", &Config, "https://raw.githubusercontent.com/cdle/sillyGirl/main/conf/demo_config.yaml")
	InitReplies()
	initToHandleMessage()
	file, err := os.Open("/etc/sillyGirl/sets.conf")
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if regexp.MustCompile(`^\s*set`).MatchString(line) {
				Senders <- &Faker{
					Message: strings.Trim(line, " "),
				}
			}
		}
		file.Close()
	}
	initSys()
	Duration = time.Duration(sillyGirl.GetInt("duration", 5)) * time.Second
}
