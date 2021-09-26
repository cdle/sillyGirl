package dwz

import (
	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
)

func init() {
	core.Server.GET("/dwz", func(c *gin.Context) {
		addr := c.Param("addr")
		c.String(200, addr)
	})
}
