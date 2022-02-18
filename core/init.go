package core

import (
	"bufio"
	"os"
	"regexp"
	"time"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/cdle/sillyGirl/utils"
)

var Duration time.Duration

var DataHome = utils.GetDataHome()

func Init() {
	sillyGirl = MakeBucket("sillyGirl")
	_, err := os.Stat(DataHome)
	if err != nil {
		os.MkdirAll(DataHome, os.ModePerm)
	}
	for _, arg := range os.Args {
		if arg == "-d" {
			utils.Daemon()
		}
	}
	ReadYaml(utils.ExecPath+"/conf/", &Config, "https://raw.githubusercontent.com/cdle/sillyGirl/main/conf/demo_config.yaml")
	InitReplies()
	initToHandleMessage()
	file, err := os.Open(DataHome + "/sets.conf")
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if v := regexp.MustCompile(`^\s*set\s+(\S+)\s+(\S+)\s+(\S+.*)`).FindStringSubmatch(line); len(v) > 0 {
				b := MakeBucket(v[1])
				if b.GetString(v[2]) != v[3] {
					b.Set(v[2], v[3])
				}
			}
		}
		file.Close()
	}
	initSys()
	Duration = time.Duration(sillyGirl.GetInt("duration", 5)) * time.Second
	sillyGirl.Set("started_at", time.Now().Format("2006-01-02 15:04:05"))
	api_key := sillyGirl.GetString("api_key")
	if api_key == "" {
		api_key := time.Now().UnixNano()
		sillyGirl.Set("api_key", api_key)
	}
	if sillyGirl.GetString("uuid") == "" {
		sillyGirl.Set("uuid", utils.GenUUID())
	}
	httplib.SetDefaultSetting(httplib.BeegoHTTPSettings{
		ConnectTimeout:   time.Second * 10,
		ReadWriteTimeout: time.Second * 10,
	})
	initGoja()
	initReboot()
}
