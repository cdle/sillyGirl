package core

import (
	"errors"

	"github.com/cdle/sillyGirl/im"
)

func init() {
	AddCommand("", []Function{
		{
			Rules: []string{"set ? ?"},
			Handle: func(s im.Sender) interface{} {
				sillyGirl.Set(s.Get(0), s.Get(1))
				return "设置成功"
			},
		},
		{
			Rules: []string{"delete ?"},
			Handle: func(s im.Sender) interface{} {
				sillyGirl.Set(s.Get(), "")
				return "删除成功"
			},
		},
		{
			Rules: []string{"get ?"},
			Handle: func(s im.Sender) interface{} {
				v := sillyGirl.Get(s.Get())
				if v == "" {
					return errors.New("空值")
				}
				return v
			},
		},
	})
}
