package wxgzh

import (
	"encoding/json"
	"fmt"

	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
	server "github.com/rixingyike/wechat"
)

var wxsv = core.NewBucket("wxsv")
var app *server.Server

func init() {
	if wxsv.Get("app_id") == "" {
		return
	}
	app = server.New(&server.WxConfig{
		AppId:          wxsv.Get("app_id"),
		Secret:         wxsv.Get("app_secret"),
		Token:          wxsv.Get("token"),
		EncodingAESKey: wxsv.Get("encoding_aes_key"),
		DateFormat:     "XML",
	})

	core.Pushs["wxsv"] = func(i1 interface{}, s1 string, _ interface{}, _ string) {
		app.SendText(fmt.Sprint(i1), s1)
	}
	core.OttoFuncs["GetWxmpAccessToken"] = app.GetAccessToken
	core.Server.Any("/wxsv/", func(c *gin.Context) {
		ctx := app.VerifyURL(c.Writer, c.Request)
		sender := &Sender{}

		if ctx.Msg.Event == "subscribe" {
			sender.admin = true
			ctx.Msg.Content = "wxsv subscribe"
		}

		if ctx.Msg.Event == "CLICK" {
			sender.admin = true
			ctx.Msg.Content = "wxsv " + ctx.Msg.EventKey
		}

		sender.tp = "wxsv"
		sender.Message = ctx.Msg.Content
		sender.uid = ctx.Msg.FromUserName
		sender.ctx = ctx
		core.Senders <- sender
	})

	core.AddCommand("", []core.Function{
		{
			Admin: true,
			Rules: []string{"init wxsv"},
			// Cron:  "1 1 * * *",
			Handle: func(_ core.Sender) interface{} {
				c := &core.Faker{
					Type:    "carry",
					Message: "wxsv init",
					Carry:   make(chan string),
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
				bt := server.Menu{}
				json.Unmarshal([]byte(f), &bt)
				if len(bt.Button) < 0 {
					return "没解析出菜单，" + f
				}
				app.AddMenu(&bt)
				return f
			},
		},
	})
}
