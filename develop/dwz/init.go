package dwz

import (
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strings"

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
	return regexp.MustCompile(`https?://[\.\w]+:?\d*`).FindString(dwz.Get("address")) + "/" + dwz.Get("prefix", "d") + encode(int64(su.ID))
}

func getWz(id string) string {
	su := &ShortUrl{
		ID: int(decode(id)),
	}
	dwz.First(su)
	return fmt.Sprint(su.Url)
}

var chars string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func encode(num int64) string {
	bytes := []byte{}
	for num > 0 {
		bytes = append(bytes, chars[num%62])
		num = num / 62
	}
	reverse(bytes)
	return string(bytes)
}

func decode(str string) int64 {
	var num int64
	n := len(str)
	for i := 0; i < n; i++ {
		pos := strings.IndexByte(chars, str[i])
		num += int64(math.Pow(62, float64(n-i-1)) * float64(pos))
	}
	return num
}

func reverse(a []byte) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}
