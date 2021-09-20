package wxgzh

import (
	"fmt"

	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

var wxmp = core.NewBucket("wxmp")

func init() {
	core.Server.Any("/wx", func(c *gin.Context) {
		wc := wechat.NewWechat()
		memory := cache.NewMemory()
		cfg := &offConfig.Config{
			AppID:          wxmp.Get("app_id"),
			AppSecret:      wxmp.Get("app_secret"),
			Token:          wxmp.Get("token"),
			EncodingAESKey: wxmp.Get("encoding_aes_key"),
			Cache:          memory,
		}
		officialAccount := wc.GetOfficialAccount(cfg)
		server := officialAccount.GetServer(c.Request, c.Writer)
		server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
			//TODO
			//回复消息：演示回复用户发送的消息
			fmt.Println(msg.Content)
			text := message.NewText(msg.Content)
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
		})
		//处理消息接收以及回复
		err := server.Serve()
		if err != nil {
			fmt.Println(err)
			return
		}
		//发送回复的消息
		server.Send()
	})
}
