package core

import (
	"strings"

	"github.com/cdle/sillyGirl/utils"
	"github.com/gin-gonic/gin"
)

type Master struct {
	Platform string `json:"platform"`
	Nickname string `json:"nickname"`
	ID       string `json:"number"`
	Index    int    `json:"id"`
	Unix     int    `json:"unix"`
}

func init() {
	GinApi(GET, "/api/master/list", func(c *gin.Context) {
		plts := getPltsArray()
		ms := []Master{}
		i := 1
		for _, plt := range plts {
			v := MakeBucket(plt)
			masters := strings.Split(v.GetString("masters"), "&")
			for _, master := range masters {
				if master == "" {
					continue
				}
				nk := Nickname{ID: master}
				nickname.First(&nk)
				ms = append(ms, Master{
					Platform: plt,
					Nickname: nk.Value,
					ID:       master,
					Index:    i,
					Unix:     nk.Unix,
				})
				i++
			}
		}
		c.JSON(200, map[string]interface{}{
			"success":   true,
			"data":      ms,
			"platforms": getPltsLabel(),
		})
	})
	GinApi(POST, "/api/master", func(c *gin.Context) {
		m := Master{}
		c.BindJSON(&m)
		if m.ID == "" {
			c.JSON(200, map[string]interface{}{
				"success":      false,
				"errorMessage": "缺少号码字段",
			})
			return
		}
		if m.Platform == "" {
			nk := Nickname{ID: m.ID}
			nickname.First(&nk)
			if nk.Platform != "" {
				m.Platform = nk.Platform
			}
		}
		if m.Platform == "" {
			c.JSON(200, map[string]interface{}{
				"success":      false,
				"errorMessage": "缺少平台字段",
			})
			return
		}
		v := MakeBucket(m.Platform)
		masters := strings.Split(v.GetString("masters"), "&")
		v.Set("masters", strings.Join(utils.Unique(masters, m.ID), "&"))

		c.JSON(200, map[string]interface{}{
			"success": true,
		})
	})
	GinApi(DELETE, "/api/master", func(c *gin.Context) {
		m := Master{}
		c.BindJSON(&m)
		if m.ID == "" {
			c.JSON(200, map[string]interface{}{
				"success":      false,
				"errorMessage": "缺少账号字段",
			})
			return
		}
		if m.Platform == "" {
			nk := Nickname{ID: m.ID}
			nickname.First(&nk)
			if nk.Platform != "" {
				m.Platform = nk.Platform
			}
		}
		if m.Platform == "" {
			c.JSON(200, map[string]interface{}{
				"success":      false,
				"errorMessage": "缺少平台字段",
			})
			return
		}
		v := MakeBucket(m.Platform)
		masters := strings.Split(v.GetString("masters"), "&")
		v.Set("masters", strings.Join(utils.Remove(masters, m.ID), "&"))

		c.JSON(200, map[string]interface{}{
			"success": true,
		})
	})
}
