package core

import (
	"errors"

	"github.com/cdle/sillyGirl/im"
)

func init() {
	AddCommand("", []Function{
		{
			Rules: []string{"set ? ? ?"},
			Handle: func(s im.Sender) interface{} {
				Bucket(s.Get(0)).Set(s.Get(1), s.Get(2))
				return "设置成功"
			},
		},
		{
			Rules: []string{"delete ? ?"},
			Handle: func(s im.Sender) interface{} {
				Bucket(s.Get(0)).Set(s.Get(1), "")
				return "删除成功"
			},
		},
		{
			Rules: []string{"get ? ? ?"},
			Handle: func(s im.Sender) interface{} {
				v := Bucket(s.Get(0)).Get(s.Get(1))
				if v == "" {
					return errors.New("空值")
				}
				return v
			},
		},
	})
}
