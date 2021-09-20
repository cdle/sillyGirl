package core

import "github.com/gin-gonic/gin"

var Server *gin.Engine

func init() {
	Server = gin.New()
}

func RunServer() {
	Server.Run(":" + sillyGirl.Get("port", "8080"))
}
