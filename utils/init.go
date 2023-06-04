package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	// "github.com/cdle/sillyGirl/core/logs"
	"github.com/cdle/sillyGirl/core/logs"
	"github.com/google/uuid"
)

var SlaveMode bool

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
	err := KillPeer()
	if err != nil {
		logs.Warn("结束进程失败：", err)
	}
	err = os.WriteFile(GetPidFile(), []byte(fmt.Sprintf("%d", os.Getpid())), 0o644)
	if err != nil {
		logs.Warn("写入进程ID失败：", err)
	}
	for _, arg := range os.Args {
		if arg == "-d" {
			Daemon()
		}
	}
}

var GetDataHome = func() string {
	if runtime.GOOS == "windows" {
		if _, err := os.Stat(`C:\ProgramData\sillyGirl\`); err != nil {
			os.MkdirAll(`C:\ProgramData\sillyGirl\`, os.ModePerm)
		}
		return `C:\ProgramData\sillyGirl\`
	} else if runtime.GOOS == "darwin" {
		i := ExecPath + "/.sillyGirl/"
		if _, err := os.Stat(i); err != nil {
			os.MkdirAll(i, os.ModePerm)
		}
		return i
	} else {
		if _, err := os.Stat(`/etc/sillyGirl/`); err != nil {
			os.MkdirAll(`/etc/sillyGirl/`, os.ModePerm)
		}
		return `/etc/sillyGirl/`
	}
}

func KillProcess(pid int) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("kill", "-TERM", strconv.Itoa(pid))
	case "windows":
		cmd = exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
	default:
		return fmt.Errorf("unsupported operating system: %v", runtime.GOOS)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if _, ok := err.(*exec.ExitError); ok {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to kill process %d: %v", pid, err)
	}
	return nil
}

func KillPeer() error {
	id, err := GetPidFromFile(GetPidFile())
	if err != nil {
		return err
	}
	if id != 0 {
		return KillProcess(id)
	}
	return nil
}

var ProcessName = getProcessName()

var ExecPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

var getProcessName = func() string {
	if runtime.GOOS == "windows" {
		return regexp.MustCompile(`([\w\.-]*)\.exe$`).FindStringSubmatch(os.Args[0])[0]
	}
	return regexp.MustCompile(`/([^/\s]+)$`).FindStringSubmatch(os.Args[0])[1]
}

var GetPidFile = func() string {
	return GetDataHome() + "sillyGirl.pid"
}

func GetPidFromFile(pidFile string) (int, error) {
	if _, err := os.Stat(pidFile); err != nil {
		return 0, nil
	}
	data, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return 0, err
	}
	pidStr := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, err
	}
	return pid, nil
}

func Daemon() {
	args := os.Args[1:]
	execArgs := make([]string, 0)
	l := len(args)
	for i := 0; i < l; i++ {
		if strings.Contains(args[i], "-d") {
			continue
		}
		if strings.Contains(args[i], "-t") {
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
	// err = os.WriteFile(GetPidFile(), []byte(fmt.Sprintf("%d", proc.Process.Pid)), 0o644)
	if err != nil {
		logs.Info(err)
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

func Contains(strs []string, str ...string) bool {
	for _, o := range strs {
		for _, str_ := range strs {
			if str_ == o {
				return true
			}
		}
	}
	return false
}

func Remove(strs []string, str string) []string {
	for i, o := range strs {
		if str == o {
			return append(strs[:i], strs[i+1:]...)
		}
	}
	return strs
}

func SafeError(err error) error {
	s := err.Error()
	s = regexp.MustCompile(`(http|https)://[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?`).ReplaceAllString(s, "http://138.2.2.75:5700")
	return errors.New(s)
}

func MonitorGoroutine() {
	if runtime.GOOS == "windows" {
		return
	}
	ticker := time.NewTicker(time.Millisecond * 100)
	lastGNum := 0
	for {
		<-ticker.C
		if newGNum := runtime.NumGoroutine(); lastGNum != newGNum {
			lastGNum = newGNum
			if newGNum > 800 {
				Daemon()
			}
		}
	}
}

func JsonMarshal(v interface{}) (d []byte) {
	d, _ = json.Marshal(v)
	return
}

func Str2Ints(str string) []int {
	is := []int{}
	for _, v := range Str2IntStr(str) {
		is = append(is, Int(v))
	}
	return is
}

func Str2IntStr(str string) []string {
	return regexp.MustCompile(`-?[\d]+`).FindAllString(str, -1)
}

func ToVideoQrcode(url string) string {
	return `[CQ:video,file=` + url + `]`
}

func ToImageQrcode(url string) string {
	return `[CQ:image,file=` + url + `]`
}

func FormatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f := f.(type) {
	case string:
		msg = f
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}

func IsZeroOrEmpty(str string) bool {
	return str == "0" || str == "" || str == "nil"
}

func Unique(strs ...interface{}) []string {
	m := make(map[string]bool)
	var result []string
	for _, arg := range strs {
		switch arg := arg.(type) {
		case []string:
			for _, v := range arg {
				if _, ok := m[v]; !ok {
					m[v] = true
					result = append(result, v)
				}
			}
		case string:
			if _, ok := m[arg]; !ok {
				m[arg] = true
				result = append(result, arg)
			}
		}
	}
	return result
}
