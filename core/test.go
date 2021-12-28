package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/beego/beego/v2/adapter/httplib"
)

func init() {
	go func() {
		v := sillyGirl.Get("rebootInfo")
		defer sillyGirl.Set("rebootInfo", "")
		if v != "" {
			vv := strings.Split(v, " ")
			tp, cd, ud := vv[0], Int(vv[1]), vv[2]
			if tp == "fake" { //&& sillyGirl.GetBool("update_notify", false) == true {
				// time.Sleep(time.Second * 10)
				// NotifyMasters("自动更新完成。")
				return
			}
			msg := "重启完成。"
			for i := 0; i < 10; i++ {
				if cd == 0 {
					if push, ok := Pushs[tp]; ok {
						push(ud, msg, nil, "")
						break
					}
				} else {
					if push, ok := GroupPushs[tp]; ok {
						push(cd, ud, msg, "")
						break
					}
				}
				time.Sleep(time.Second)
			}
		}
	}()
}

func initSys() {
	AddCommand("", []Function{
		// {//
		// 	Rules: []string{"unintsall sillyGirl"},
		// 	Admin: true,
		// 	Handle: func(s Sender) interface{} {
		// 		return ""
		// 	},
		// },
		//
		{
			Rules: []string{"raw ^name$"},
			Handle: func(s Sender) interface{} {
				s.Disappear()
				return name()
			},
		},
		{
			Rules: []string{"reply ? ?"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				a := s.Get(1)
				if a == "nil" {
					a = ""
				}
				Bucket(fmt.Sprintf("reply%s%d", s.GetImType(), s.GetChatID())).Set(s.Get(0), a)
				return "设置成功。"
			},
		},
		{
			Rules: []string{"replies"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				rt := ""
				Bucket(fmt.Sprintf("reply%s%d", s.GetImType(), s.GetChatID())).Foreach(func(k, v []byte) error {
					rt += fmt.Sprintf("%s === %s\n", k, v)
					return nil
				})
				return strings.Trim(rt, "\n")
			},
		},
		{
			Rules: []string{"raw ^卸载$"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				if runtime.GOOS == "windows" {
					return "windows系统不支持此命令"
				}
				s.Reply("您真的要卸载" + name() + "吗？(5秒后默认卸载，Y/n)")
				switch s.Await(s, func(s Sender) interface{} {
					return YesNo
				}, time.Second*5) {
				case No:
					return name() + "将继续为您服务！"
				}
				s.Reply("是否删除用户数据？(5秒后默认删除，Y/n)")
				clear := true
				switch s.Await(s, func(s Sender) interface{} {
					return YesNo
				}, time.Second*5) {
				case No:
					clear = false
					return name() + "将继续为您服务！"
				}
				s.Reply("进入冷静期，给你5秒时间思考，输入任意字符取消卸载：")
				if s.Await(s, nil, time.Second*5) != nil {
					return name() + "将继续为您服务！"
				}
				s.Reply("你终究还是下得了狠心，不过那又怎样？")
				time.Sleep(time.Second * 2)
				s.Reply("请在5秒内输入“我是🐶”完成卸载：")
				rt := s.Await(s, nil, time.Second*5)
				switch rt.(type) {
				case nil:
					return "你的打字速度不够快啊，请重新卸载～"
				case string:
					if rt.(string) != "我是🐶" {
						return "输入错误，请重新卸载～"
					}
				}
				if !sillyGirl.GetBool("forbid_uninstall") {
					if clear {
						os.RemoveAll(dataHome)
					}
					os.RemoveAll(ExecPath)
					os.RemoveAll("/usr/lib/systemd/system/sillyGirl.service")
				}
				s.Reply("卸载完成，下次重启你就再也见不到我了。")
				time.Sleep(time.Second)
				s.Reply("是否立即重启？")
				s.Reply("正在重启...")
				time.Sleep(time.Second)
				os.Exit(0)
				return nil
			},
		},
		{
			Rules: []string{"raw ^升级$"},
			// Cron:  "*/1 * * * *",
			Admin: true,
			Handle: func(s Sender) interface{} {
				if runtime.GOOS == "windows" {
					return "windows系统不支持此命令"
				}

				if s.GetImType() == "fake" && !sillyGirl.GetBool("auto_update", true) && compiled_at == "" {
					return nil
				}

				if compiled_at != "" {
					str := ""
					for i, prefix := range []string{"https://ghproxy.com/", ""} {
						if str == "" && s.GetImType() != "fake" {
							if v, ok := OttoFuncs["version"]; ok {
								if rt := v(""); rt != "" {
									str = regexp.MustCompile(`\d{13}`).FindString(rt)
								}
							}
						}
						if str == "" {
							data, _ := httplib.Get(prefix + "https://raw.githubusercontent.com/cdle/binary/master/compile_time.go").String()
							rt := regexp.MustCompile(`\d{13}`).FindString(data)
							if strings.Contains(data, "package") {
								str = rt
							}
						}
						if str != "" {
							if s.GetImType() == "fake" {
								ver := sillyGirl.Get("compiled_at")
								if str > ver && ver > compiled_at {
									return nil
								}
								if ver < str && str > compiled_at {
									sillyGirl.Set("compiled_at", str)
									NotifyMasters(fmt.Sprintf("检测到更新版本(%s)。", str))
								}
								return nil
							} else {
								s.Reply(fmt.Sprintf("检测到最新版本(%s)。", str))
							}
							if str > compiled_at {
								if i == 0 {
									s.Reply("正在从ghproxy.com下载更新...")
								} else {
									s.Reply("尝试从github.com下载更新...")
								}
								req := httplib.Get(prefix + "https://raw.githubusercontent.com/cdle/binary/master/sillyGirl_linux_" + runtime.GOARCH + "_" + str)
								if i == 1 && Transport != nil {
									req.SetTransport(Transport)
								}
								req.SetTimeout(time.Minute*5, time.Minute*5)
								data, err := req.Bytes()
								if err != nil {
									// return "下载程序错误：" + err.Error()
									continue
								}
								if len(data) < 2646147 {
									// return "下载失败。"
									continue
								}
								filename := ExecPath + "/" + pname
								if err = os.RemoveAll(filename); err != nil {
									return "删除旧程序错误：" + err.Error()
								}

								if f, err := os.OpenFile(filename, syscall.O_CREAT, 0777); err != nil {
									return "创建程序错误：" + err.Error()
								} else {
									_, err := f.Write(data)
									f.Close()
									if err != nil {
										des := err.Error()
										if err = os.WriteFile(filename, data, 777); err != nil {
											return "写入程序错误：" + des + "\n" + err.Error()
										}
									}
								}
								s.Reply("更新完成，重启生效，是否立即重启？(Y/n，3秒后自动确认。)")
								if s.Await(s, func(s Sender) interface{} {
									return YesNo
								}, time.Second*3) == No {
									return "好的，下次重启生效。。"
								}
								go func() {
									time.Sleep(time.Second)
									Daemon()
								}()
								sillyGirl.Set("rebootInfo", fmt.Sprintf("%v %v %v", s.GetImType(), s.GetChatID(), s.GetUserID()))
								return "正在重启。"
							} else {
								return fmt.Sprintf("当前版本(%s)最新，无需升级。", compiled_at)
							}
						} else {
							continue
						}
					}
					return `无法升级，你网不好。建议您手动于linux执行一键升级命令： s=sillyGirl;a=arm64;if [[ $(uname -a | grep "x86_64") != "" ]];then a=amd64;fi ;if [ ! -d $s ];then mkdir $s;fi ;cd $s;wget https://mirror.ghproxy.com/https://github.com/cdle/${s}/releases/download/main/${s}_linux_$a -O $s && chmod 777 $s;pkill -9 $s;$(pwd)/$s`
				}

				s.Reply("开始检查核心更新...", E)
				update := false
				record := func(b bool) {
					if !update && b {
						update = true
					}
				}
				need, err := GitPull("")
				if err != nil {
					return "请使用以下命令手动升级：\n cd " + ExecPath + " && git stash && git pull && go build && ./" + pname
				}
				if !need {
					s.Reply("核心功能已是最新。", E)
				} else {
					record(need)
					s.Reply("核心功能发现更新。", E)
				}
				files, _ := ioutil.ReadDir(ExecPath + "/develop")
				for _, f := range files {
					if f.IsDir() && f.Name() != "replies" {
						if f.Name() == "qinglong" {
							continue
						}
						if strings.HasPrefix(f.Name(), "_") {
							continue
						}
						s.Reply("检查扩展"+f.Name()+"更新...", E)
						need, err := GitPull("/develop/" + f.Name())
						if err != nil {
							s.Reply("扩展"+f.Name()+"更新错误"+err.Error()+"。", E)
						}
						if !need {
							s.Reply("扩展"+f.Name()+"已是最新。", E)
						} else {
							record(need)
							s.Reply("扩展"+f.Name()+"发现更新。", E)
						}
					}
				}
				if !update {
					s.Reply("没有更新。", E)
					return nil
				}
				s.Reply("正在编译程序...", E)
				if err := CompileCode(); err != nil {
					return "请使用以下命令手动编译：\n cd " + ExecPath + " && go build && ./" + pname
				}
				s.Reply("编译程序完毕。", E)
				sillyGirl.Set("rebootInfo", fmt.Sprintf("%v %v %v", s.GetImType(), s.GetChatID(), s.GetUserID()))
				s.Reply("更新完成，即将重启！", E)
				go func() {
					time.Sleep(time.Second)
					Daemon()
				}()
				return nil
			},
		},
		{
			Rules: []string{"raw ^编译$"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				if compiled_at != "" {
					return "编译个🐔8。"
				}
				s.Reply("正在编译程序...", E)
				if err := CompileCode(); err != nil {
					return err
				}
				s.Reply("编译程序完毕。", E)
				return nil
			},
		},
		{
			Rules: []string{"raw ^重启$"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				s.Disappear()
				sillyGirl.Set("rebootInfo", fmt.Sprintf("%v %v %v", s.GetImType(), s.GetChatID(), s.GetUserID()))
				s.Reply("即将重启！", E)
				Daemon()
				return nil
			},
		},
		{
			Rules: []string{"raw ^status$"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				return fmt.Sprintf("总计：%d，已处理：%d，运行中：%d", total, finished, total-finished)
			},
		},
		{
			Rules: []string{"raw ^命令$"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				s.Disappear()
				ss := []string{}
				ruless := [][]string{}
				for _, f := range functions {
					if len(f.Rules) > 0 {
						rules := []string{}
						for i := range f.Rules {
							rules = append(rules, f.Rules[i])
						}
						ruless = append(ruless, rules)
					}
				}
				for j := range ruless {
					for i := range ruless[j] {
						ruless[j][i] = strings.Trim(ruless[j][i], "^$")
						ruless[j][i] = strings.Replace(ruless[j][i], `(\S+)`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `(\S*)`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `(.+)`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `(.*)`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `\s+`, " ", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `\s*`, " ", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `.+`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `.*`, "?", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `[(]`, "(", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `[)]`, ")", -1)
						ruless[j][i] = strings.Replace(ruless[j][i], `([\s\S]+)`, "?", -1)
					}
					ss = append(ss, strings.Join(ruless[j], "\n"))
				}

				return strings.Join(ss, "\n")
			},
		},
		{
			Admin: true,
			Rules: []string{"set ? ? ?", "delete ? ?", "? set ? ?", "? delete ?", "set ? ?", "? set ?"},
			Handle: func(s Sender) interface{} {
				name := s.Get(0)
				if name == "silly" {
					name = "sillyGirl"
				}
				b := Bucket(name)
				if !IsBucket(b) && !strings.HasPrefix(name, "tgc_") {
					s.Continue()
					return nil
				}
				old := b.Get(s.Get(1))
				b.Set(s.Get(1), s.Get(2))
				go func() {
					s.Await(s, func(_ Sender) interface{} {
						b.Set(s.Get(1), old)
						return "已撤回。"
					}, "^撤回$", time.Second*60)
				}()
				return "操作成功，在60s内可\"撤回\"。"
			},
		},
		{
			Admin: true,
			Rules: []string{"get ? ?", "? get ?"},
			Handle: func(s Sender) interface{} {
				name := s.Get(0)
				if name == "silly" {
					name = "sillyGirl"
				}
				b := Bucket(name)
				if !IsBucket(b) {
					s.Continue()
					return nil
				}
				s.Disappear()
				v := b.Get(s.Get(1))
				if v == "" {
					return errors.New("无值")
				}
				return v
			},
		},
		{
			Admin: true,
			Rules: []string{"list ?"},
			Handle: func(s Sender) interface{} {
				name := s.Get(0)
				if name == "silly" {
					name = "sillyGirl"
				}
				if s.GetChatID() != 0 && name != "reply" {
					return "请私聊我。"
				} //fanlivip
				if name != "fanlivip" && name != "otto" && name != "reply" && name != "wxsv" && name != "sillyGirl" && name != "qinglong" && name != "wx" && name != "wxmp" && name != "tg" && name != "qq" && !strings.HasPrefix(name, "tgc_") {
					s.Continue()
					return nil
				}
				if s.GetChatID() != 0 {
					s.Disappear()
				}
				b := Bucket(name)
				// if !IsBucket(b) {
				// s.Continue()
				// return nil
				// }
				rt := ""
				b.Foreach(func(k, v []byte) error {
					rt += fmt.Sprintf("%s === %s\n", k, v)
					return nil
				})
				return strings.Trim(rt, "\n")
			},
		},
		{
			Admin: true,
			Rules: []string{"send ? ? ?"},
			Handle: func(s Sender) interface{} {
				if push, ok := Pushs[s.Get(0)]; ok {
					push(s.Get(1), s.Get(2), nil, "")
				}
				return "发送成功呢"
			},
		},
		{
			Rules: []string{"raw ^myuid$"},
			Handle: func(s Sender) interface{} {
				return fmt.Sprint(s.GetUserID())
			},
		},
		{
			Rules: []string{"raw ^groupCode$"},
			Handle: func(s Sender) interface{} {
				return fmt.Sprint(s.GetChatID())
			},
		},
		{
			Rules: []string{"raw ^compiled_at$"},
			Handle: func(s Sender) interface{} {
				return sillyGirl.Get("compiled_at")
			},
		},
		{
			Rules: []string{"notify ?"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				NotifyMasters(s.Get())
				return "通知成功。"
			},
		},
		{
			Rules: []string{"raw ^started_at$"},
			Handle: func(s Sender) interface{} {
				return sillyGirl.Get("started_at")
			},
		},
		{
			Rules: []string{"守护傻妞"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				if runtime.GOOS == "windows" {
					return "windows系统不支持此命令"
				}
				service := `
[Unit]
Description=silly silly girl bot
After=network.target mysql.service mariadb.service mysqld.service
[Service]
Type=forking
ExecStart=` + ExecPath + "/" + pname + ` -d
PIDFile=/var/run/sillyGirl.pid
Restart=always
User=root
Group=root
				
[Install]
WantedBy=multi-user.target
Alias=sillyGirl.service`
				data, err := exec.Command("sh", "-c", "type systemctl").Output()
				if err != nil {
					s.Reply(err)
					return nil
				}

				if !strings.Contains(string(data), "bin") {
					s.Reply(data)
					return nil
				}
				os.WriteFile("/usr/lib/systemd/system/sillyGirl.service", []byte(service), 0o644)
				exec.Command("systemctl", "disable", string(sillyGirl)).Output()
				exec.Command("systemctl", "enable", string(sillyGirl)).Output()
				return "电脑重启后生效。"
			},
		},
		// {
		// 	Rules: []string{"raw .*pornhub.*"},
		// 	Handle: func(s Sender) interface{} {
		// 		s.Reply("你已涉黄永久禁言。")
		// 		for {
		// 			s.Await(s, func(s2 Sender, _ error) interface{} {
		// 				s2.Disappear(time.Millisecond * 50)
		// 				return "你已被禁言。"
		// 			}, `[\s\S]*`, time.Duration(time.Second*300))
		// 		}
		// 	},
		// },
		{
			Rules: []string{"raw ^成语接龙$"},
			Handle: func(s Sender) interface{} {
				if sillyGirl.GetBool("disable_成语接龙", false) {
					s.Continue()
					return nil
				}
				begin := ""
				fword := func(cy string) string {
					begin = strings.Replace(regexp.MustCompile(`([一-龥])】`).FindString(cy), "】", "", -1)
					return begin
				}
				id := fmt.Sprintf("%v", s.GetUserID())
			start:
				data, err := httplib.Get("http://hm.suol.cc/API/cyjl.php?id=" + id + "&msg=开始成语接龙").String()
				if err != nil {
					s.Reply(err)
				}
				s.Reply(data)
				fword(data)
				stop := false
				win := false
				if strings.Contains(data, "你赢") {
					stop = true
					win = true
				}
				if strings.Contains(data, "我赢") {
					stop = true
				}
				if !stop {
					s.Await(s, func(s2 Sender) interface{} {
						ct := s2.GetContent()

						me := s2.GetUserID() == s.GetUserID()
						if strings.Contains(ct, "小爱提示") || ct == "q" {
							s2.SetContent(fmt.Sprintf("小爱%s字开头的成语有哪些？", begin))
							s2.Continue()
							return Again
						}
						if strings.Contains(ct, "认输") {
							if me || s2.IsAdmin() {
								stop = true
								return nil
							} else {
								return GoAgain("你认输有个屁用。")
							}
						}
						if regexp.MustCompile("^"+begin).FindString(ct) == "" || strings.Contains(ct, "接龙") {
							if me {
								return GoAgain(fmt.Sprintf("现在是接【%s】开头的成语哦。", begin))
							} else {
								if ct == "成语接龙" {
									return GoAgain(fmt.Sprintf("现在是接【%s】开头的成语哦。", begin))
								}
								s2.Continue()
								return Again
							}
						}
						cy := regexp.MustCompile("^[一-龥]+$").FindString(ct)
						if cy == "" {
							s2.Disappear(time.Millisecond * 500)
							return GoAgain("请认真接龙，一站到底！")
						}
						data, err := httplib.Get("http://hm.suol.cc/API/cyjl.php?id=" + id + "&msg=我接" + cy).String()
						if err != nil {
							s2.Reply(err)
							return Again
						}
						if strings.Contains(data, "file_get_contents") {
							ss := strings.Split(data, "\n")
							return GoAgain(ss[len(ss)-1])
						}
						if strings.Contains(data, "你赢") {
							stop = true
							win = true
							if !me {
								defer s.Reply("反正不是你赢，嘿嘿。")
							}
						} else if strings.Contains(data, "我赢") {
							stop = true
							win = false
						} else if strings.Contains(data, "恭喜") {
							fword(data)
							if !me {
								data += "\n你很可拷，观棋不语真君子懂不懂啊。"
							}
						} else {
							if me {
								data += "\n玩不过就认输呗。"
							} else {
								data += "\n你以为你会，结果出丑了吧。"
							}
						}
						if !stop {
							return GoAgain(data)
						}
						return data
					}, ForGroup)
				}
				time.Sleep(time.Microsecond * 100)
				s.Reply("还玩吗？[Y/n]")
				if s.Await(s, func(s2 Sender) interface{} {
					return YesNo
				}, time.Second*6) == Yes {
					goto start
				}
				if !win {
					s.Reply("菜*，见一次虐一次！")
				} else {
					s.Reply("大爷下次再来玩啊～")
				}
				return nil
			},
		},
		{
			Rules: []string{"^machineId$"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				return fmt.Sprintf("你的机器码：%s", OttoFuncs["machineId"](""))
			},
		},
		{
			Rules: []string{"^time$"},
			Handle: func(s Sender) interface{} {
				return OttoFuncs["timeFormat"]("2006-01-02 15:04:05")
			},
		},
	})
}

func IsBucket(b Bucket) bool {
	for i := range Buckets {
		if Buckets[i] == b {
			return true
		}
	}
	return false
}
