package core

import (
	"strings"

	"github.com/beego/beego/v2/adapter/httplib"
)

var Pushs = map[string]func(interface{}, string, interface{}, string){}
var GroupPushs = map[string]func(interface{}, interface{}, string, string){}

// type Chat struct {
// 	Class  string
// 	ID     int
// 	UserID int
// }

// func (ct *Chat) Push(content interface{}) {
// 	switch content.(type) {
// 	case string:
// 		if push, ok := GroupPushs[ct.Class]; ok {
// 			push(ct.ID, ct.UserID, content.(string))
// 		}
// 	case error:
// 		if push, ok := GroupPushs[ct.Class]; ok {
// 			push(ct.ID, ct.UserID, content.(error).Error())
// 		}
// 	}
// }

func NotifyMasters(content string) {
	go func() {
		content = strings.Trim(content, " ")
		if sillyGirl.GetBool("ignore_notify", false) == true {
			return
		}
		if token := sillyGirl.Get("pushplus"); token != "" {
			httplib.Get("http://www.pushplus.plus/send?token=" + token + "&title=0101010&content=" + content + "&template=html")
		}
		for _, class := range []string{"tg", "qq", "wx"} {
			notify := Bucket(class).Get("notifiers")
			if notify == "" {
				notify = Bucket(class).Get("masters")
			}
			for _, v := range strings.Split(notify, "&") {
				if push, ok := Pushs[class]; ok {
					push(v, content, nil, "")
				}
			}
		}
	}()
}
