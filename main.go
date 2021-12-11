package main

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/cdle/sillyGirl/core"
)

func main() {
	go core.RunServer()
	logs.Info("傻妞用不了了？关注频道 https://t.me/kczz2021 获取最新消息。")
	for {
		var input string
		fmt.Scanln(&input)
		fmt.Println(input)
		core.Senders <- &core.Faker{
			Type: "terminal",
		}
	}
}
