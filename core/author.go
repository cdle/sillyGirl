package core

func init() {
	AddCommand("", []Function{
		{
			Admin: true,
			Rules: []string{"raw (吃早点)", "raw (午饭)", "raw (饿了么)", "raw (饿了)", "raw (外卖)", "raw (美团)"},
			Handle: func(s Sender) interface{} {
				return nil
			},
		},
	})
}
