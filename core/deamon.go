package core

import (
	"os"
	"os/exec"
	"strings"

	"github.com/astaxie/beego/logs"
)

func Daemon() {
	args := os.Args[1:]
	execArgs := make([]string, 0)
	l := len(args)
	for i := 0; i < l; i++ {
		if strings.Index(args[i], "-d") == 0 {
			continue
		}

		execArgs = append(execArgs, args[i])
	}
	proc := exec.Command(os.Args[0], execArgs...)
	err := proc.Start()
	if err != nil {
		panic(err)
	}
	logs.Info("傻妞以静默形式运行")
	os.Exit(0)
}
