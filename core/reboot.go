package core

import (
	"os"
	"strings"
	"time"

	"github.com/cdle/sillyGirl/utils"
)

func initReboot() {
	go func() {
		data, _ := os.ReadFile(utils.ExecPath + "/rebootInfo")
		if v := string(data); v != "" {
			defer os.RemoveAll(utils.ExecPath + "/rebootInfo")
			vv := strings.Split(v, " ")
			tp, cd, ud := vv[0], vv[1], vv[2]
			if tp == "fake" {
				return
			}
			msg := "重启完成。"
			for i := 0; i < 10; i++ {
				dapter, _ := GetAdapter(tp, "")
				if dapter != nil {
					break
				}
				dapter.Push(Message{
					USER_ID: ud,
					CHAT_ID: cd,
					CONETNT: msg,
				})
				time.Sleep(time.Second)
			}
		}
	}()
}
