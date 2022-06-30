package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/DeanThompson/ginpprof"
	"github.com/beego/beego/v2/core/logs"
	"github.com/cdle/sillyGirl/core"
	"github.com/cdle/sillyGirl/utils"
)

func main() {
	core.Init()
	ginpprof.Wrapper(core.Server)
	sillyGirl := core.MakeBucket("sillyGirl")
	go monitorGoroutine()
	port := sillyGirl.GetString("port", "8080")
	logs.Info("Http服务已运行(%s)。", sillyGirl.GetString("port", "8080"))
	go core.Server.Run("0.0.0.0:" + port)
	logs.Info("关注频道 https://t.me/kczz2021 获取最新消息。")
	d := false
	for _, arg := range os.Args {
		if arg == "-d" {
			d = true
		}
	}
	if !d {
		t := false
		for _, arg := range os.Args {
			if arg == "-t" {
				t = true
			}
		}
		if t {
			logs.Info("终端交互已启用。")
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				data := scanner.Text()
				f := &core.Faker{
					Type:    "terminal",
					Message: string(data),
					Carry:   make(chan string),
				}
				core.Senders <- f
				go func() {
					for v := range f.Listen() {
						fmt.Printf("\x1b[%dm%s \x1b[0m\n", 31, v)
					}
				}()
			}
			logs.Info("终端交互异常,请检查运行环境设置,如果是docker环境,请附加-it参数")
		} else {
			logs.Info("终端交互不可用，运行带-t参数即可启用。")
		}
	}

	select {}
}

func monitorGoroutine() {
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
				utils.Daemon()
			}
		}
	}
}
