package main

import (
	"bufio"
	"os"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/cdle/sillyGirl/core"
)

func main() {
	core.Init123()
	sillyGirl := core.Bucket("sillyGirl")
	port := sillyGirl.Get("port", "8080")
	logs.Info("Http服务已运行(%s)。", sillyGirl.Get("port", "8080"))
	go core.Server.Run("0.0.0.0:" + port)
	logs.Info("关注频道 https://t.me/kczz2021 获取最新消息。")
	reader := bufio.NewReader(os.Stdin)
	var tm int64
	for {
		data, _, _ := reader.ReadLine()
		core.Senders <- &core.Faker{
			Type:    "terminal",
			Message: string(data),
		}
		nm := time.Now().UnixNano() / int64(time.Millisecond)
		if tm == 0 {
			tm = nm
		} else {
			if nm-tm < 10 {
				break
			}
		}
	} //
	select {}
}
