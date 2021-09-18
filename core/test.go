package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/cdle/sillyGirl/im"
)

func init() {
	go func() {
		v := sillyGirl.Get("rebootInfo")
		if v != "" {
			vv := strings.Split(v, " ")
			tp, cd, ud := vv[0], Int(vv[1]), Int(vv[2])
			msg := name() + "重启完成。"
			if cd == 0 {
				Push(tp, ud, msg)
			} else {
				for i := 0; i < 10; i++ {
					if push, ok := GroupPushs[tp]; ok {
						push(cd, ud, msg)
						break
					}
					time.Sleep(time.Second)
				}
			}
			sillyGirl.Set("rebootInfo", "")
		}
	}()
}

func initSys() {
	AddCommand("", []Function{
		{
			Rules: []string{"raw ^name$"},
			Handle: func(s im.Sender) interface{} {
				s.Disappear()
				return name()
			},
		},
		{
			Rules: []string{"raw ^升级$"},
			Admin: true,
			Handle: func(s im.Sender) interface{} {
				s.Disappear()
				s.Reply(name() + "开始检查核心功能。")
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
					s.Reply(name() + "核心功能已是最新。")
				} else {
					record(need)
					s.Reply(name() + "核心功能发现更新。")
				}
				files, _ := ioutil.ReadDir(ExecPath + "/develop")
				for _, f := range files {
					if f.IsDir() {
						need, err := GitPull("/develop/" + f.Name())
						if err != nil {
							s.Reply(name() + "扩展" + f.Name() + "更新错误" + err.Error() + "。")
						}
						if !need {
							s.Reply(name() + "扩展" + f.Name() + "已是最新。")
						} else {
							record(need)
							s.Reply(name() + "扩展" + f.Name() + "发现更新。")
						}
					}
				}
				if !update {
					return name() + "没有更新。"
				}
				s.Reply(name() + "正在编译程序。")
				if err := CompileCode(); err != nil {
					return err
				}
				s.Reply(name()+"编译程序完毕。", time.Duration(0))
				sillyGirl.Set("rebootInfo", fmt.Sprintf("%v %v %v", s.GetImType(), s.GetChatID(), s.GetUserID()))
				Daemon()
				return nil
			},
		},
		{
			Rules: []string{"raw ^重启$"},
			Admin: true,
			Handle: func(s im.Sender) interface{} {
				s.Disappear()
				sillyGirl.Set("rebootInfo", fmt.Sprintf("%v %v %v", s.GetImType(), s.GetChatID(), s.GetUserID()))
				Daemon()
				return nil
			},
		},
		{
			Rules: []string{"raw ^命令$"},
			Handle: func(s im.Sender) interface{} {
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
			Handle: func(s im.Sender) interface{} {
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
			Handle: func(s im.Sender) interface{} {
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
			Handle: func(s im.Sender) interface{} {
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
