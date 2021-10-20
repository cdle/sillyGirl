package core

import "strings"

func init() {
	go func() {
		recall := sillyGirl.Get("recall")
		if recall != "" {
			rules := []string{}
			for _, v := range strings.Split(recall, "&") {
				rules = append(rules, "raw "+v)
			}
			AddCommand("", []Function{
				{
					Rules: rules,
					Handle: func(s Sender) interface{} {
						if !s.IsAdmin() {
							s.Delete()
						}
						return nil
					},
				},
			})
		}
	}()
}
