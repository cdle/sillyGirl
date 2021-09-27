package dwz

import (
	"fmt"
	"net/http"

	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
)

var dwz = core.NewBucket("dwz")

func init() {
	core.Server.GET("/dwz", func(c *gin.Context) {
		url := c.Query("url")
		dwz := getDwz(url)
		c.String(200, dwz)
	})
	core.Server.GET("/d:id", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, getWz(id))
	})
	core.AddCommand("", []core.Function{
		{
			Rules: []string{"dwz ?"},
			Handle: func(s core.Sender) interface{} {
				return getDwz(s.Get())
			},
		},
	})
}

type ShortUrl struct {
	ID  int
	Url string
}

func getDwz(url string) string {
	su := &ShortUrl{
		Url: url,
	}
	dwz.Create(su)
	return dwz.Get("address", "https://4co.cc"+"/d"+fmt.Sprint(su.ID))
}

func getWz(id string) string {
	su := &ShortUrl{
		ID: core.Int(id),
	}
	dwz.First(su)
	return fmt.Sprint(su.Url)
}
