package core

import (
	"errors"

	"github.com/cdle/sillyGirl/im"
)

func init() {
	AddCommand("", []Function{
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
