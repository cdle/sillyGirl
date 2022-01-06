package wxgzh

import (
	"fmt"

	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
	server "github.com/rixingyike/wechat"
)

var wxsv = core.NewBucket("wxsv")

func init() {
	cfg := &server.WxConfig{
		AppId:          wxsv.Get("app_id"),
		Secret:         wxsv.Get("app_secret"),
		Token:          wxsv.Get("token"),
		EncodingAESKey: wxsv.Get("encoding_aes_key"),
		DateFormat:     "XML",
	}
	app := server.New(cfg)
	// app.AddMenu(&server.Menu{
	// 	Button: []server.Button{
	// 		{
	// 			Name: "购物功能",
	// 		},
	// 		{
	// 			Name: "好玩功能",
	// 		},
	// 		{
	// 			Name: "其他功能",
	// 		},
	// 	},
	// })
	core.Pushs["wxsv"] = func(i1 interface{}, s1 string, _ interface{}, _ string) {
		app.SendText(fmt.Sprint(i1), s1)
	}
	core.Server.Any("/wxsv/", func(c *gin.Context) {
		ctx := app.VerifyURL(c.Writer, c.Request)
		if ctx.Msg.Event == "subscribe" {
			ctx.NewText(wxsv.Get("subscribe_reply", "感谢关注！")).Reply()
			return
		}
		sender := &Sender{}
		sender.tp = "wxsv"
		sender.Message = ctx.Msg.Content
		sender.uid = ctx.Msg.FromUserName
		sender.ctx = ctx
		core.Senders <- sender
	})

	core.AddCommand("", []core.Function{
		{
			Admin: true,
			Rules: []string{"init wxsv menu"},
			Cron:  "1 1 * * *",
			Handle: func(_ core.Sender) interface{} {
				c := &core.Faker{
					Type:    "carry",
					Message: wxsv.Get("app_id"),
				}
				core.Senders <- c
				f := ""
				for {
					v, ok := <-c.Listen()
					if !ok {
						break
					}
					f = v
				}
				return f
			},
		},
	})
}
