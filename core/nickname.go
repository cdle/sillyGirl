package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Nickname struct {
	ID       string   `json:"i"`
	Group    bool     `json:"g"`
	Unix     int      `json:"u"`
	Value    string   `json:"v"`
	Platform string   `json:"p"`
	BotsID   []string `json:"bs"`
}

var nickname = MakeBucket("nickname")

type NicklabeL struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	Platform string `json:"platform"`
	ChatName string `json:"chat_name"`
}

func CreateNickName(nick *Nickname) {
	nick.Unix = int(time.Now().Unix())
	nickname.Create(nick)
}

func init() {
	GinApi(GET, "/api/nickname/labels", RequireAuth, func(ctx *gin.Context) {
		group := true
		keyword := ctx.Query("gkeyword")
		if keyword == "" {
			keyword = ctx.Query("ukeyword")
			if keyword != "" {
				group = false
			}
		}
		if keyword == "" && strings.Contains(ctx.Request.URL.String(), "ukeyword") {
			group = false
		}
		platform := ctx.Query("platform")
		data := []NicklabeL{}
		data2 := []NicklabeL{}
		// if keyword != "" {
		// full := false
		nickname.Foreach(func(b1, b2 []byte) error {
			v := &Nickname{}
			code := string(b1)
			err := json.Unmarshal(b2, v)
			if err == nil {
				if v.Group != group {
					return nil
				}
				if platform != "" && v.Platform != platform {
					return nil
				}
				nl := NicklabeL{
					ChatName: v.Value,
					Value:    code,
					Platform: v.Platform,
				}
				if !group {
					nl.Label = fmt.Sprintf("%s(%s)", v.Value, code)
				} else {
					nl.Label = fmt.Sprintf("%s %s@%s", v.Value, code, v.Platform)
				}
				if strings.HasPrefix(code, keyword) || strings.Contains(v.Value, keyword) {
					data = append(data, nl)
					// if code == keyword {
					// 	full = true
					// }
				}
				data2 = append(data2, nl)
			}
			return nil
		})
		// if !full {
		// 	data = append([]NicklabeL{{
		// 		Label: keyword,
		// 		Value: keyword,
		// 	}}, data...)
		// } else
		if len(data) == 0 {
			data = append([]NicklabeL{{
				Label: keyword,
				Value: keyword,
			}}, data2...)
		}
		// }
		ctx.JSON(200, map[string]interface{}{
			"success": true,
			"data":    data,
		})
	})
}
