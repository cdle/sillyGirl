package main

import (
	"bufio"
	"os"

	"github.com/beego/beego/v2/core/logs"
	"github.com/cdle/sillyGirl/core"
)

func main() {
	go core.RunServer()
	logs.Info("傻妞用不了了？关注频道 https://t.me/kczz2021 获取最新消息。")
	reader := bufio.NewReader(os.Stdin)
	core.Init123()
	for {
		data, _, _ := reader.ReadLine()
		core.Senders <- &core.Faker{
			Type:    "terminal",
			Message: string(data),
		}
	}
}
