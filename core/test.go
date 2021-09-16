package core

import (
	"errors"
	"strings"

	"github.com/cdle/sillyGirl/im"
)

func initSys() {
	AddCommand("", []Function{
		{
			Rules: []string{"raw ^name$"},
			Handle: func(_ im.Sender) interface{} {
				return name()
			},
		},
		{
			Rules: []string{"raw ^升级$"},
			Admin: true,
			Handle: func(s im.Sender) interface{} {
				s.Reply(name() + "开始拉取代码。")
				need, err := GitPull("")
				if err != nil {
					return err
				}
				if !need {
					return name() + "已是最新版。"
				}
				s.Reply(name() + "开始拉取成功。")
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
			Handle: func(_ im.Sender) interface{} {
				Daemon()
				return nil
			},
		},
		{
			Rules: []string{"raw ^命令$"},
			Handle: func(_ im.Sender) interface{} {
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
