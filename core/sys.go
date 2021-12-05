package core

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

var BeforeStop = []func(){}

var getPidf = func() string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s\\%s", ExecPath, "sillyGirl.pid")
	}
	return "/var/run/sillyGirl.pid"
}

var pidf = getPidf()

func Daemon() {
	for _, bs := range BeforeStop {
		bs()
	}
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
	err = os.WriteFile(pidf, []byte(fmt.Sprintf("%d", proc.Process.Pid)), 0o644)
	if err != nil {
		logs.Warn(err)
	}
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
	app := "sh"
	param := "-c"
	if runtime.GOOS == "windows" {
		app = "cmd"
		param = "/c"
	}
	cmd := exec.Command(app, param, "cd "+ExecPath+" && go build -o "+pname)
	_, err := cmd.Output()
	if err != nil {
		return errors.New("编译失败：" + err.Error() + "。")
	}
	sillyGirl.Set("compiled_at", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

func Download() error {
	url := "https://github.com/cdle/sillyGirl/releases/download/main/sillyGirl_linux_"
	if sillyGirl.GetBool("downlod_use_ghproxy", false) { //
		url = "https://mirror.ghproxy.com/" + url
	}
	url += runtime.GOARCH
	cmd := exec.Command("sh", "-c", "cd "+ExecPath+" && wget "+url+" -O temp && mv temp "+pname+"  && chmod 777 "+pname)
	_, err := cmd.Output()
	if err != nil {
		return errors.New("失败：" + err.Error() + "。")
	}
	// sillyGirl.Set("compiled_at", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

func killp() {
	// data, _ := os.ReadFile(pidf)
	// pid := Int(string(data))
	// if pid > 0 {
	// 	syscall.Kill(-pid, syscall.SIGKILL)
	// }
	pids, err := ppid()
	if err == nil {
		if len(pids) == 0 {
			return
		} else {
			exec.Command("sh", "-c", "kill -9 "+strings.Join(pids, " ")).Output()
		}
	} else {
		return
	}
}

func ppid() ([]string, error) {
	pid := fmt.Sprint(os.Getpid())
	pids := []string{}
	rtn, err := exec.Command("sh", "-c", "pidof "+pname).Output()
	if err != nil {
		return pids, err
	}
	re := regexp.MustCompile(`[\d]+`)
	for _, v := range re.FindAll(rtn, -1) {
		if string(v) != pid {
			pids = append(pids, string(v))
		}
	}
	return pids, nil
}
