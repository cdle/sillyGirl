package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/cdle/sillyGirl/core/storage"
	"github.com/cdle/sillyGirl/utils"
	"github.com/gin-gonic/gin"
)

var authBucket = MakeBucket("auths")
var auths = []*Auth{}

func init() {
	storage.Watch(sillyGirl, "name", func(old, new, key string) *storage.Final {
		if old == new {
			return &storage.Final{
				Error: errors.New("unchanged"),
			}
		}
		return nil
	})
	authBucket.Foreach(func(b1, b2 []byte) error {
		auth := &Auth{}
		if json.Unmarshal(b2, auth) == nil {
			if math.Abs(float64(int(time.Now().Unix())-auth.CreatedAt)) < 86400 {
				auths = append(auths, auth)
			}
		}
		return nil
	})
	var password = sillyGirl.GetString("password")
	var name = sillyGirl.GetString("name", "傻妞")
	if password == "" {
		password = utils.GenUUID()
		console.Info("可视化面板临时账号密码：%s %s", name, password)
	}
	storage.Watch(sillyGirl, "password", func(old, new, key string) *storage.Final {
		password, _ = EncryptByAes([]byte(new))
		return &storage.Final{
			Now: password,
		}
	})
	storage.Watch(sillyGirl, "name", func(old, new, key string) *storage.Final {
		name = new
		return nil
	})
	///可视化部分
	GinApi(POST, "/api/login/account", func(ctx *gin.Context) {
		var auth = struct {
			Password string `json:"password"`
			Username string `json:"username"`
		}{}
		json.NewDecoder(ctx.Request.Body).Decode(&auth)
		epassword, _ := EncryptByAes([]byte(auth.Password))
		if (auth.Password == password || epassword == password) && auth.Username == name {
			token := utils.GenUUID()
			auth := &Auth{
				IP:        ctx.ClientIP(),
				UserAgent: ctx.Request.UserAgent(),
				Token:     token,
				CreatedAt: int(time.Now().Unix()),
			}
			authBucket.Create(auth)
			auths = append(auths, auth)
			console.Log("登录成功，当前有效令牌数%d，总数%d", len(ValidAuths()), len(auths))
			ctx.SetCookie("token", token, 86400, "/", "", false, true)
			ctx.JSON(200, map[string]interface{}{
				"status":           "ok",
				"type":             "account",
				"currentAuthority": "admin",
			})
		} else {
			ctx.JSON(200, map[string]interface{}{
				"status":           "error",
				"type":             "account",
				"currentAuthority": "guest",
			})
		}
	})
	GinApi(POST, "/api/login/outLogin", DestroyAuth, func(ctx *gin.Context) {
		sillyGirl.Set("web_token", "")
		ctx.JSON(200, map[string]interface{}{
			"success": true,
		})
	})
	pluginNextUuid := sillyGirl.GetString("pluginNextUuid")
	if pluginNextUuid == "" {
		pluginNextUuid = utils.GenUUID()
		sillyGirl.Set("pluginNextUuid", pluginNextUuid)
	}
	GinApi(GET, "/api/currentUser", RequireAuth, func(ctx *gin.Context) {
		rs := []Route{}
		for _, f := range Functions {
			if f.UUID == pluginNextUuid {
				pluginNextUuid = utils.GenUUID()
				sillyGirl.Set("pluginNextUuid", pluginNextUuid)
			}
			if f.UUID != "" {
				name := f.Title
				if name == "" {
					name = "无名脚本"
				}
				if f.Module {
					name = name + " 🔧"
				}
				if f.OnStart {
					name = name + " 💫"
				}
				if f.Encrypt {
					name = name + " 🔒"
				}
				if f.Public {
					name = name + " 👑"
				}
				rs = append(rs, Route{
					Path:      fmt.Sprintf(`/script/%s`, f.UUID),
					Name:      name,
					Component: "./Script",
					CreateAt:  f.CreateAt,
				})
			}
		}
		rrs := rs
		n := len(rrs)
		flag := true
		for i := 0; i < n && flag; i++ {
			flag = false
			for j := 0; j < n-i-1; j++ {
				if rrs[j].CreateAt < rrs[j+1].CreateAt {
					rrs[j], rrs[j+1] = rrs[j+1], rrs[j]
					flag = true
				}
			}
		}
		rrs = append(rrs, Route{
			Path:      fmt.Sprintf(`/script/%s`, pluginNextUuid),
			Name:      "+新增脚本",
			Component: "./Script",
		})
		ctx.JSON(200, map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"name":    sillyGirl.GetString("name"),
				"avatar":  "https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png",
				"plugins": rrs,
			},
		})
	})
}

func DestroyAuth(c *gin.Context) {
	token, _ := c.Cookie("token")
	auth, _ := CheckAuth(token)
	if auth != nil {
		auth.ExpiredAt = int(time.Now().Unix())
		authBucket.Create(auth)
	}
}

var tempAuth sync.Map

func getTempAuth() string {
	uuid := utils.GenUUID()
	tempAuth.Store(uuid, time.Now().Unix())
	return uuid
}

func checkTempAuth(uuid string) bool {
	unix, ok := tempAuth.LoadAndDelete(uuid)
	if !ok {
		return false
	}
	if time.Now().Unix()-unix.(int64) > 1 {
		return false
	}
	return true
}

func RequireAuth(c *gin.Context) {
	token, _ := c.Cookie("token")
	_, err := CheckAuth(token)
	if err != nil && !checkTempAuth(token) {
		c.JSON(401, map[string]interface{}{
			"data": map[string]interface{}{
				"isLogin": false,
			},
			"errorCode":    "401",
			"errorMessage": err.Error(),
			"success":      true,
			"showType":     9,
		})
		panic(err)
	}
}

func CheckAuth(token string) (*Auth, error) {
	var errorMessage = "请先登录！"
	if token != "" {
		auths := auths
		for i := range auths {
			if auths[i].Token == token && auths[i].ExpiredAt == 0 {
				if math.Abs(float64(int(time.Now().Unix())-auths[i].CreatedAt)) > 86400 {
					auths[i].ExpiredAt = int(time.Now().Unix())
					authBucket.Create(auths[i])
					errorMessage = "授权已过期！"
				} else {
					return auths[i], nil
				}
			} else {
				errorMessage = "非法访问！"
			}
		}
	}
	return nil, errors.New(errorMessage)
}

func ValidAuths() []*Auth {
	tmp := []*Auth{}
	for _, auth := range auths {
		if auth.ExpiredAt == 0 {
			tmp = append(tmp, auth)
		}

	}
	return tmp
}
