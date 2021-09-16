package core

import "os"

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
}
