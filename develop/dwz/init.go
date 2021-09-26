package dwz

import (
	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
)

func init() {
	core.Server.GET("/http:add", func(c *gin.Context) {
		addr := "http" + c.Param("add")
		c.String(200, addr)
	})
}
