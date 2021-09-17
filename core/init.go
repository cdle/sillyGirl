package core

import (
	"os"
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
	initSys()
	Duration = time.Duration(sillyGirl.GetInt("duration", 5)) * time.Second
}
