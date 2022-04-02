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
	"github.com/cdle/sillyGirl/utils"
)

func GitPull(filename string) (bool, error) {
	if runtime.GOOS == "darwin" {
		return false, errors.New("骂你一句沙雕。")
	}
	rtn, err := exec.Command("sh", "-c", "cd "+utils.ExecPath+filename+" && git stash && git pull").Output()
	if err != nil {
		return false, errors.New("拉取代失败：" + err.Error() + "。")
	}
	t := string(rtn)
	if !strings.Contains(t, "changed") {
		if strings.Contains(t, "Already") || strings.Contains(t, "已经是最新") {
			return false, nil
		} else {
			return false, errors.New("拉取代失败：" + t + "。")
		}
	}
	return true, nil
}

func CompileCode() error {
	app := "sh"
	param := "-c"
	if runtime.GOOS == "windows" {
		app = "cmd"
		param = "/c"
	}
	cmd := exec.Command(app, param, "cd "+utils.ExecPath+" && go build -o "+utils.ProcessName)
	_, err := cmd.Output()
	if err != nil {
		return errors.New("编译失败：" + err.Error() + "。")
	}
	sillyGirl.Set("compiled_at", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

func Download() error {
	url := "https://github.com/cdle/sillyGirl/releases/download/main/sillyGirl_linux_"
	if sillyGirl.GetBool("downlod_use_ghproxy", false) { //
		url = "https://mirror.ghproxy.com/" + url
	}
	url += runtime.GOARCH
	cmd := exec.Command("sh", "-c", "cd "+utils.ExecPath+" && wget "+url+" -O temp && mv temp "+utils.ProcessName+"  && chmod 777 "+utils.ProcessName)
	_, err := cmd.Output()
	if err != nil {
		return errors.New("失败：" + err.Error() + "。")
	}
	// sillyGirl.Set("compiled_at", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

func initReboot() {
	go func() {
		v := sillyGirl.GetString("rebootInfo")
		defer sillyGirl.Set("rebootInfo", "")
		if v != "" {
			vv := strings.Split(v, " ")
			tp, cd, ud := vv[0], utils.Int(vv[1]), vv[2]
			if tp == "fake" {
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
		{
			Rules: []string{"raw ^name$"},
			Handle: func(s Sender) interface{} {
				s.Disappear()
				return name()
			},
		},
		{
			Rules: []string{"reply empty all"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				b := MakeBucket(fmt.Sprintf("reply%s%d", s.GetImType(), s.GetChatID()))
				b.Foreach(func(k, v []byte) error {
					b.Set(string(k), "")
					return nil
				})
				return "清空成功。"
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
				MakeBucket(fmt.Sprintf("reply%s%d", s.GetImType(), s.GetChatID())).Set(s.Get(0), a)
				return "设置成功。"
			},
		},
		{
			Rules: []string{"replies"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				rt := ""
				MakeBucket(fmt.Sprintf("reply%s%d", s.GetImType(), s.GetChatID())).Foreach(func(k, v []byte) error {
					rt += fmt.Sprintf("%s === %s\n", k, v)
					return nil
				})
				return strings.Trim(rt, "\n")
			},
		},
		{
			Rules: []string{"升级 ?", "^升级$"},
			// Cron:  "*/1 * * * *",
			Admin: true,
			Handle: func(s Sender) interface{} {
				if runtime.GOOS == "windows" {
					return "windows系统不支持此命令"
				}
				if s.GetImType() == "fake" && !sillyGirl.GetBool("auto_update", true) && compiled_at == "" {
					return nil
				}
				var kz = s.Get(0)
				if compiled_at != "" {
					str := ""
					pxs := []string{}
					if p := sillyGirl.GetString("download_prefix"); p != "" {
						pxs = append(pxs, p)
					}
					pxs = append(pxs, "")
					pxs = append(pxs, "https://gitee.yanyuge.workers.dev/")
					pxs = append(pxs, "https://ghproxy.com/")
					for _, prefix := range pxs {
						if str == "" && s.GetImType() != "fake" {
							if v, ok := OttoFuncs["version"]; ok {
								if rt := v.(func(string) string)(""); rt != "" {
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
								ver := sillyGirl.GetString("compiled_at")
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
								s.Reply(fmt.Sprintf("正在从%s下载更新...", prefix))
								req := httplib.Get(prefix + "https://raw.githubusercontent.com/cdle/binary/master/sillyGirl_linux_" + runtime.GOARCH + "_" + str)
								if prefix == "" && Transport != nil {
									req.SetTransport(Transport)
								}
								req.SetTimeout(time.Minute*5, time.Minute*5)
								data, err := req.Bytes()
								if err != nil {
									continue
								}
								if len(data) < 2646147 {
									continue
								}
								filename := utils.ExecPath + "/" + utils.ProcessName
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
									utils.Daemon()
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
					return `无法升级，你网不好。建议您手动于linux执行一键升级命令： s=sillyGirl;a=arm64;if [[ $(uname -a | grep "x86_64") != "" ]];then a=amd64;fi ;cd ` + utils.ExecPath + `$s;wget https://github.com/cdle/${s}/releases/download/main/${s}_linux_$a -O $s && chmod 777 $s;pkill -9 $s;$(pwd)/$s -t`
				}

				s.Reply("开始检查核心更新...", E)
				update := false
				record := func(b bool) {
					if !update && b {
						update = true
					}
				}
				var need bool
				var err error
				if kz == "" || kz == "core" {
					need, err = GitPull("")
					if err != nil {
						return "请使用以下命令手动升级：\n cd " + utils.ExecPath + " && git stash && git pull && go build && ./" + utils.ProcessName
					}
					if !need {

					} else {
						record(need)
						s.Reply("核心功能发现更新。", E)
					}
				}

				files, _ := ioutil.ReadDir(utils.ExecPath + "/develop")
				for _, f := range files {
					if f.IsDir() && f.Name() != "replies" {
						if kz != "" && kz != f.Name() {
							continue
						}
						if strings.HasPrefix(f.Name(), "_") {
							continue
						}
						need, err := GitPull("/develop/" + f.Name())
						if err != nil {
							s.Reply("扩展"+f.Name()+"更新错误"+err.Error()+"。", E)
						}
						if !need {
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
					return "请使用以下命令手动编译：\n cd " + utils.ExecPath + " && go build && ./" + utils.ProcessName
				}
				s.Reply("编译程序完毕。", E)
				sillyGirl.Set("rebootInfo", fmt.Sprintf("%v %v %v", s.GetImType(), s.GetChatID(), s.GetUserID()))
				s.Reply("更新完成，即将重启！", E)
				go func() {
					time.Sleep(time.Second)
					utils.Daemon()
				}()
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
				utils.Daemon()
				return nil
			},
		},
		{
			Rules: []string{"raw ^status$"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				s.Disappear()
				ss := []string{}
				contents.Range(func(key, value interface{}) bool {
					ss = append(ss, fmt.Sprintf("%v. %v", key, value))
					return true
				}) //runtime.NumGoroutine()
				return fmt.Sprintf("总计：%d，已处理：%d，运行中：%d\n\n%s", total, finished, total-finished, strings.Join(ss, "\n"))
			},
		},
		{
			Rules: []string{"raw ^命令$"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				s.Disappear()
				ss := []string{}
				ruless := [][]string{}
				for _, f := range Functions {
					if len(f.Rules) > 0 {
						if f.Show != "" {
							ss = append(ss, fmt.Sprint(f.Show))
							continue
						}
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
			Admin:    true,
			Priority: 10000,
			Rules:    []string{"set ? ? ?", "delete ? ?", "? set ? ?", "? delete ?", "set ? ?", "? set ?"},
			Handle: func(s Sender) interface{} {
				name := s.Get(0)
				if name == "silly" {
					name = "sillyGirl"
				}
				b := MakeBucket(name)
				old := b.GetString(s.Get(1))
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
			Rules: []string{"empty ?", "empty ? ?", "? empty ?"},
			Handle: func(s Sender) interface{} {
				name := s.Get(0)
				filter := s.Get(1)
				if name == "silly" {
					name = "sillyGirl"
				}
				a := ""
				if filter != "" {
					a = "中包含" + filter
				}
				s.Reply("20秒内回复任意取消清空" + name + a + "的记录。")

				switch s.Await(s, nil, time.Second*20) {
				case nil:
				case "快":
				default:
					return "已取消。"
				}
				if filter == "" {
					// db.Update(func(t *bolt.Tx) error {
					// 	err := t.DeleteBucket([]byte(name))
					// 	if err != nil {
					// 		s.Reply(err)
					// 	}
					// 	return nil
					// })
					return fmt.Sprintf("已清空。")
				}
				b := MakeBucket(name)
				i := 0
				b.Foreach(func(k, v []byte) error {
					if filter == "" || strings.Contains(string(k), filter) || strings.Contains(string(v), filter) {
						b.Set(string(k), "")
						i++
					}
					return nil
				})
				return fmt.Sprintf("已清空%d个记录。", i)
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
				b := MakeBucket(name)
				s.Disappear()
				v := b.GetString(s.Get(1))
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
				}
				// if name != "fanlivip" && name != "otto" && name != "reply" && name != "wxsv" && name != "sillyGirl" && name != "qinglong" && name != "wx" && name != "wxmp" && name != "tg" && name != "qq" && !strings.HasPrefix(name, "tgc_") {
				// 	s.Continue()
				// 	return nil
				// }
				if s.GetChatID() != 0 {
					s.Disappear()
				}
				b := MakeBucket(name)
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
				return sillyGirl.GetString("compiled_at")
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
				return sillyGirl.GetString("started_at")
			},
		},
		{
			Rules: []string{"^machineId$"},
			Admin: true,
			Handle: func(s Sender) interface{} {
				return fmt.Sprintf("你的机器码：%s", OttoFuncs["machineId"].(func(string) string)(""))
			},
		},
		{
			Rules: []string{"^time$"},
			Handle: func(s Sender) interface{} {
				return OttoFuncs["timeFormat"].(func(string) string)("2006-01-02 15:04:05")
			},
		},
	})
	if !isReleaseVersion() {
		AddCommand("", []Function{
			{
				Rules: []string{"raw ^编译$"},
				Admin: true,
				Handle: func(s Sender) interface{} {
					s.Reply("正在编译程序...", E)
					if err := CompileCode(); err != nil {
						return err
					}
					s.Reply("编译程序完毕。", E)
					return nil
				},
			},
		})
	}
	if !inDocker() {
		return
	}
	AddCommand("", []Function{
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
						os.RemoveAll(DataHome)
					}
					os.RemoveAll(utils.ExecPath)
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
ExecStart=` + utils.ExecPath + "/" + utils.ProcessName + ` -d
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
				exec.Command("systemctl", "disable", sillyGirl.String()).Output()
				exec.Command("systemctl", "enable", sillyGirl.String()).Output()
				return "电脑重启后生效。"
			},
		},
	})
}
func inDocker() bool {
	info, e := os.Stat("/.dockerenv")
	return e != nil && info != nil && !info.IsDir() && info.Size() == 0
}

func isReleaseVersion() bool {
	return compiled_at != ""
}
