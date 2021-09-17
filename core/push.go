package core

var Pushs = map[string]func(int, string){}
var GroupPushs = map[string]func(int, int, string){}

func Push(class string, uid int, content string) {
	if push, ok := Pushs[class]; ok {
		push(uid, content)
	}
}

type Chat struct {
	Class  string
	ID     int
	UserID int
}

func (ct *Chat) Push(content interface{}) {
	switch content.(type) {
	case string:
		if push, ok := GroupPushs[ct.Class]; ok {
			push(ct.ID, ct.UserID, content.(string))
		}
	case error:
		if push, ok := GroupPushs[ct.Class]; ok {
			push(ct.ID, ct.UserID, content.(error).Error())
		}
	}
}
