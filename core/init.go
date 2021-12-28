package core

import (
	"bufio"
	"os"
	"regexp"
	"runtime"
	"time"

	"github.com/beego/beego/v2/adapter/httplib"
)

var Duration time.Duration

var dataHome = ""

func init() {
	killp()
	if runtime.GOOS == "windows" {
		_, err := os.Stat(`C:\ProgramData\sillyGirl`)
		if err != nil {
			os.MkdirAll(`C:\ProgramData\sillyGirl`, os.ModePerm)
		}
		dataHome = `C:\ProgramData\sillyGirl`
	} else {
		_, err := os.Stat("/etc/sillyGirl")
		if err != nil {
			os.MkdirAll("/etc/sillyGirl", os.ModePerm)
		}
		dataHome = `/etc/sillyGirl`
	}

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

	file, err := os.Open(dataHome + "/sets.conf")

	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if v := regexp.MustCompile(`^\s*set\s+(\S+)\s+(\S+)\s+(\S+.*)`).FindStringSubmatch(line); len(v) > 0 {
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
	sillyGirl.Set("started_at", time.Now().Format("2006-01-02 15:04:05"))
	api_key := sillyGirl.Get("api_key")
	if api_key == "" {
		api_key := time.Now().UnixNano()
		sillyGirl.Set("api_key", api_key)
	}
	if sillyGirl.Get("uuid") == "" {
		sillyGirl.Set("uuid", GetUUID())
	}
	httplib.SetDefaultSetting(httplib.BeegoHTTPSettings{
		ConnectTimeout:   time.Second * 10,
		ReadWriteTimeout: time.Second * 10,
	})
}
