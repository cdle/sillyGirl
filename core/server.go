package core

import "github.com/gin-gonic/gin"

var Server *gin.Engine

func init() {
	Server = gin.New()
}

var Tail = "--来自sillyGirl，傻妞技术交流QQ群882314490，电报交流群https://t.me/kczz2021。"

func RunServer() {
	if sillyGirl.GetBool("enable_http_server", false) == false {
		return
	}
	Server.GET("/", func(c *gin.Context) {
		c.String(200, Tail)
	})
	gin.SetMode(gin.ReleaseMode)
	Server.Run("0.0.0.0:" + sillyGirl.Get("port", "8080"))
}
