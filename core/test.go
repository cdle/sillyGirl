package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

func init() {
	go func() {
		v := sillyGirl.Get("rebootInfo")
		if v != "" {
			vv := strings.Split(v, " ")
			tp, cd, ud := vv[0], Int(vv[1]), Int(vv[2])
			if tp == "fake" && sillyGirl.GetBool("update_notify", false) == true { //
				time.Sleep(time.Second * 10)
				NotifyMasters("自动更新完成。")
				return
			}
			msg := "重启完成。"
			for i := 0; i < 10; i++ {
				if cd == 0 {
					if push, ok := Pushs[tp]; ok {
						push(ud, msg)
						break
					}
				} else {
					if push, ok := GroupPushs[tp]; ok {
						push(cd, ud, msg)
						break
					}
				}
				time.Sleep(time.Second)
			}
			sillyGirl.Set("rebootInfo", "")
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
			Rules: []string{"raw ^升级$"},
			Cron:  "*/1 * * * *",
			Admin: true,
			Handle: func(s Sender) interface{} {
				s.Reply("开始检查核心更新...", E)
				update := false
				record := func(b bool) {
					if !update && b {
						update = true
					}
				}
				need, err := GitPull("")
				if err != nil {
					return err
				}
				if !need {
					s.Reply("核心功能已是最新。", E)
				} else {
					record(need)
					s.Reply("核心功能发现更新。", E)
				}
				files, _ := ioutil.ReadDir(ExecPath + "/develop")
				for _, f := range files {
					if f.IsDir() {
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
					return err
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
			Rules: []string{"raw ^命令$"},
			Handle: func(s Sender) interface{} {
				s.Disappear()
				ss := []string{}
				for _, f := range functions {
					ss = append(ss, strings.Join(f.Rules, " "))
				}
				return strings.Join(ss, "\n")
			},
		},
		{
			Admin: true,
			Rules: []string{"set ? ? ?"},
			Handle: func(s Sender) interface{} {
				s.Disappear()
				b := Bucket(s.Get(0))
				if !IsBucket(b) {
					return errors.New("不存在的存储桶")
				}
				b.Set(s.Get(1), s.Get(2))
				return "设置成功"
			},
		},
		{
			Admin: true,
			Rules: []string{"delete ? ?"},
			Handle: func(s Sender) interface{} {
				s.Disappear()
				b := Bucket(s.Get(0))
				if !IsBucket(b) {
					return errors.New("不存在的存储桶")
				}
				b.Set(s.Get(1), "")
				return "删除成功"
			},
		},
		{
			Admin: true,
			Rules: []string{"get ? ?"},
			Handle: func(s Sender) interface{} {
				s.Disappear()
				b := Bucket(s.Get(0))
				if !IsBucket(b) {
					return errors.New("不存在的存储桶")
				}
				v := b.Get(s.Get(1))
				if v == "" {
					return errors.New("空值")
				}
				return v
			},
		},
		{
			Admin: true,
			Rules: []string{"send ? ? ?"},
			Handle: func(s Sender) interface{} {
				Push(s.Get(0), Int(s.Get(1)), s.Get(2))
				return "发送成功呢"
			},
		},
		{
			Rules: []string{"raw ^myuid$"},
			Handle: func(s Sender) interface{} {
				return fmt.Sprint(s.GetUserID())
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
