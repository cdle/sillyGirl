package core

var Pushs = map[string]func(int, string){}

func Push(class string, uid int, content string) {
	if push, ok := Pushs[class]; ok {
		push(uid, content)
	}
}
