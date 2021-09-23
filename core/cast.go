package core

import (
	"time"

	"github.com/beego/beego/v2/adapter/httplib"
)

var channels = []string{}

func init() {
	go func() {
		time.Sleep(time.Second * 20)
		for {
			for _, channel := range channels {
				msg, _ := httplib.Get(channel).String()
				if msg != "" && sillyGirl.Get(channel) != msg {
					NotifyMasters(msg)
					sillyGirl.Set(channel, msg)
				}
			}
			time.Sleep(time.Minute)
		}
	}()
}
