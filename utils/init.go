package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/google/uuid"
)

func GenUUID() string {
	u2, _ := uuid.NewUUID()
	return u2.String()
}

func Float64(str interface{}) float64 {
	f, _ := strconv.ParseFloat(fmt.Sprint(str), 64)
	return f
}

func TrimHiddenCharacter(originStr string) string {
	srcRunes := []rune(originStr)
	dstRunes := make([]rune, 0, len(srcRunes))
	for _, c := range srcRunes {
		if c >= 0 && c <= 31 && c != 10 {
			continue
		}
		if c == 127 {
			continue
		}

		dstRunes = append(dstRunes, c)
	}
	return strings.ReplaceAll(string(dstRunes), "￼", "")
}

func ForCQ(content string, callback func(key string, values map[string]string)) {

}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Itob(i uint64) []byte {
	return []byte(fmt.Sprint(i))
}

var Int = func(s interface{}) int {
	i, _ := strconv.Atoi(fmt.Sprint(s))
	return i
}

var Int64 = func(s interface{}) int64 {
	i, _ := strconv.Atoi(fmt.Sprint(s))
	return int64(i)
}

func init() {
	for _, arg := range os.Args {
		if arg == "-d" {
			Daemon()
		}
	}
	KillPeer()
}

var GetDataHome = func() string {
	if runtime.GOOS == "windows" {
		return `C:\ProgramData\sillyGirl`
	} else {
		return `/etc/sillyGirl`
	}
}

func KillPeer() {
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

var ProcessName = getProcessName()

func ppid() ([]string, error) {
	pid := fmt.Sprint(os.Getpid())
	pids := []string{}
	rtn, err := exec.Command("sh", "-c", "pidof "+ProcessName).Output()
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

var ExecPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

var getProcessName = func() string {
	if runtime.GOOS == "windows" {
		return regexp.MustCompile(`([\w\.-]*)\.exe$`).FindStringSubmatch(os.Args[0])[0]
	}
	return regexp.MustCompile(`/([^/\s]+)$`).FindStringSubmatch(os.Args[0])[1]
}

var GetPidFile = func() string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s\\%s", ExecPath, "sillyGirl.pid")
	}
	return "/var/run/sillyGirl.pid"
}

var Runnings = []func(){}

func Daemon() {
	for _, bs := range Runnings {
		bs()
	}
	args := os.Args[1:]
	execArgs := make([]string, 0)
	l := len(args)
	for i := 0; i < l; i++ {
		if strings.Index(args[i], "-d") != -1 {
			continue
		}
		if strings.Index(args[i], "-t") != -1 {
			continue
		}
		execArgs = append(execArgs, args[i])
	}
	proc := exec.Command(os.Args[0], execArgs...)
	err := proc.Start()
	if err != nil {
		panic(err)
	}
	logs.Info("程序以静默形式运行")
	err = os.WriteFile(GetPidFile(), []byte(fmt.Sprintf("%d", proc.Process.Pid)), 0o644)
	if err != nil {
		logs.Warn(err)
	}
	os.Exit(0)
}

func FetchCookieValue(ps ...string) string {
	var key, cookies string
	if len(ps) == 2 {
		if len(ps[0]) > len(ps[1]) {
			key, cookies = ps[1], ps[0]
		} else {
			key, cookies = ps[0], ps[1]
		}
	}
	match := regexp.MustCompile(key + `=([^;]*);{0,1}`).FindStringSubmatch(cookies)
	if len(match) == 2 {
		return strings.Trim(match[1], " ")
	} else {
		return ""
	}
}

func Contains(strs []string, str string) bool {
	for _, o := range strs {
		if str == o {
			return true
		}
	}
	return false
}

func SafeError(err error) error {
	s := err.Error()
	s = regexp.MustCompile(`(http|https)://[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?`).ReplaceAllString(s, "http://138.2.2.75:5700")
	return errors.New(s)
}
