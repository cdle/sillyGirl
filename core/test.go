package core

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/cdle/sillyGirl/im"
)

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
					record(need)
					s.Reply(name() + "核心功能已是最新。")
				}
				files, _ := ioutil.ReadDir(ExecPath + "/develop")
				for _, f := range files {
					if f.IsDir() {
						need, err := GitPull("/develop/" + f.Name())
						if err != nil {
							s.Reply(name() + "扩展" + f.Name() + "更新错误" + err.Error() + "。")
						}
						if !need {
							record(need)
							s.Reply(name() + "扩展" + f.Name() + "已是最新。")
						}
					}
				}
				if !need {
					return name() + "没有更新。"
				}
				s.Reply(name() + "正在编译程序。")
				if err := CompileCode(); err != nil {
					return err
				}
				s.Reply(name() + "编译程序完毕。")
				Daemon()
				return nil
			},
		},
		{
			Rules: []string{"raw ^重启$"},
			Admin: true,
			Handle: func(s im.Sender) interface{} {
				s.Disappear()
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
