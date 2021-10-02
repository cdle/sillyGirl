package core

import "github.com/gin-gonic/gin"

var Server *gin.Engine

func init() {
	Server = gin.New()
}

var Tail = "--来自sillyGirl，傻妞技术交流群654346133。"

func RunServer() {
	if sillyGirl.GetBool("enable_http_server", false) == false {
		return
	}
	Server.GET("/", func(c *gin.Context) {
		c.String(200, Tail)
	})
	Server.Run(":" + sillyGirl.Get("port", "8080"))
}
