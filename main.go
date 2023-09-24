package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	// _ "github.com/cdle/sillyGirl/adapters/qq"
	"github.com/cdle/sillyGirl/adapters/web"
	"github.com/cdle/sillyGirl/core"
	"github.com/cdle/sillyGirl/core/common"

	"github.com/cdle/sillyGirl/utils"
)

var app = core.MakeBucket("app")
var ip = ""

func main() {
	ip = app.GetString("ip")
	go func() {
		var err error
		ip, err = utils.GetPublicIP()
		if err == nil {
			app.Set("ip", ip)
		}
	}()
	core.Init()
	if app.GetBool("anti_kasi") {
		go utils.MonitorGoroutine()
	}
	d := false
	for _, arg := range os.Args {
		if arg == "-d" {
			d = true
		}
	}
	go func() { //弹出浏览器
		if runtime.GOOS != "windows" {
			return
		}
		time.Sleep(time.Second * 3)
		if web.GetUserNumber() == 0 {
			app := core.MakeBucket("app")
			port := app.GetInt("port", 8080)
			url := fmt.Sprintf("http://localhost:%d/admin", port)
			cmd := exec.Command("cmd", "/c", "start", url)
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(stdout))
		}
	}()
	if !d {
		t := false
		for _, arg := range os.Args {
			if arg == "-t" {
				t = true
			}
		}
		if t {
			// core.Logs.Info("Terminal机器人已连接")

			a := &core.Factory{}
			a.Init("terminal", "default", nil)
			i := 1
			a.SetIsAdmin(func(s string) bool {
				return true
			})
			a.SetReplyHandler(func(m map[string]interface{}) string {
				i++
				fmt.Printf("\x1b[%dm%s \x1b[0m\n", 31, m[core.CONETNT])
				return fmt.Sprint(i)
			})
			// a.SetActionHandler(func(m map[string]interface{}) string {
			// 	fmt.Println(`do action: ` + string(utils.JsonMarshal(m)))
			// 	return ""
			// })
			reader := bufio.NewReader(os.Stdin)
			for {
				input, err := reader.ReadString('\n')
				if err != nil {
					core.Logs.Error(err)
					a.Destroy()
					break
				}
				input = strings.TrimSpace(input)

				s := &core.CustomSender{
					F: a,
				}
				i++
				s.SetFsps(&common.FakerSenderParams{
					Content:   input,
					MessageID: fmt.Sprint(i),
				})
				core.Messages <- s
				// 在这里可以根据输入执行相应的逻辑
			}
		} else {
			// core.Logs.Info("Terminal机器人不可用，运行带-t参数即可启用")
		}
	}
	select {}
}

//  git add . && git commit -m "x" && git push
