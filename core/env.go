package core

import "fmt"

var env = NewBucket("env")

func init() {
	// var xa = `https://jintia.jintias.cn/api/xatx.php?msg={{1}}`
	// replies := []string{
	// 	`你是，(.*)？=>你好，我是{{1}}。`,
	// 	`(小爱\S*)=>gjson(req(xa,1), text)`,
	// }
	// fmt.Println(replies)
	var template = `
var content = {{1}}
var data = request

	`
	fmt.Println(template)
}
