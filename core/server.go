package core

import (
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/gin-gonic/gin"
)

var Server *gin.Engine

func init() {
	Server = gin.New()
}

var Tail = "--来自sillyGirl，傻妞技术交流QQ群882314490，电报交流群https://t.me/kczz2021。"

func RunServer() {
	if sillyGirl.GetBool("enable_http_server", false) == false {
		return
	}

	gin.SetMode(gin.ReleaseMode)
	logs.Info("开启httpserver----0.0.0.0:" + sillyGirl.Get("port", "8080"))
	Server.Run("0.0.0.0:" + sillyGirl.Get("port", "8080"))
}
