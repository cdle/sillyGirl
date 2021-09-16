package core

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
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
	logs.Info(sillyGirl.Get("name", "傻妞") + "以静默形式运行")
	os.Exit(0)
}

func GitPull(filename string) (bool, error) {
	if runtime.GOOS == "darwin" {
		return false, errors.New("骂你一句沙雕。")
	}
	rtn, err := exec.Command("sh", "-c", "cd "+ExecPath+filename+" && git stash && git pull").Output()
	if err != nil {
		return false, errors.New("拉取代失败：" + err.Error() + "。")
	}
	t := string(rtn)
	if !strings.Contains(t, "changed") {
		if strings.Contains(t, "Already") || strings.Contains(t, "已经是最新") {
			return false, nil
		} else {
			return false, errors.New("拉取代失败：" + t + "。")
		}
	}
	return true, nil
}

func CompileCode() error {
	_, err := exec.Command("sh", "-c", "cd "+ExecPath+" && go build -o "+pname).Output()
	if err != nil {
		return errors.New("编译失败：" + err.Error() + "。")
	}
	return nil
}
