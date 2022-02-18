package core

import (
	"strings"
)

var Pushs = map[string]func(interface{}, string, interface{}, string){}
var GroupPushs = map[string]func(interface{}, interface{}, string, string){}

func NotifyMasters(content string) {
	content = strings.Trim(content, " ")
	if sillyGirl.GetBool("ignore_notify", false) == true {
		return
	}
	for _, class := range []string{"tg", "qq", "wx"} {
		notify := MakeBucket(class).GetString("notifiers")
		if notify == "" {
			notify = MakeBucket(class).GetString("masters")
		}
		for _, v := range strings.Split(notify, "&") {
			if push, ok := Pushs[class]; ok {
				push(v, content, nil, "")
			}
		}
	}
}
