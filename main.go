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

	core.Init123()
	sillyGirl := core.Bucket("sillyGirl")
	port := sillyGirl.Get("port", "8080")
	logs.Info("http服务已运行(%s)。" + sillyGirl.Get("port", "8080"))
	go core.Server.Run("0.0.0.0:" + port)
	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		core.Senders <- &core.Faker{
			Type:    "terminal",
			Message: string(data),
		}
	}
}
