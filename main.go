package main

import (
	"bufio"
	"os"
	"strings"
	"time"

	_ "github.com/cdle/sillyGirl/adapters/qq"
	_ "github.com/cdle/sillyGirl/adapters/web"
	"github.com/cdle/sillyGirl/core"

	"github.com/cdle/sillyGirl/utils"
)

var sillyGirl = core.MakeBucket("sillyGirl")

func main() {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = loc
	core.Init()
	if sillyGirl.GetBool("anti_kasi") {
		go utils.MonitorGoroutine()
	}
	d := false
	for _, arg := range os.Args {
		if arg == "-d" {
			d = true
		}
		if arg == "-r" { //准备程序->原程序
			rfix := ".ready.exe"
			if strings.Contains(os.Args[0], rfix) {
				err := utils.CopyFile(utils.ProcessName, strings.Replace(utils.ProcessName, rfix, ".exe", -1))
				if err == nil {
					utils.Daemon("reset")
				}
			} else {
				os.Remove(strings.ReplaceAll(os.Args[0], ".exe", rfix))
			}
			continue
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
			core.Logs.Info("Terminal机器人已连接")
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				data := scanner.Text()
				f := &core.Faker{
					Type:    "terminal",
					Message: string(data),
					Admin:   true,
				}
				core.Messages <- f
			}
			core.Logs.Info("Terminal机器人异常,请检查运行环境设置,如果是docker环境,请附加-it参数")
		} else {
			// core.Logs.Info("Terminal机器人不可用，运行带-t参数即可启用")
		}
	}
	select {}
}

//
