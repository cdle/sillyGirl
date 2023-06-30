package core

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cdle/sillyGirl/utils"
	"github.com/dop251/goja"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

type Reply struct {
	Index     int      `json:"index,omitempty"` // 排序
	ID        int      `json:"id"`              // ID，主键
	Nickname  string   `json:"nickname"`        // 类型 0全部 1用户 2群聊
	Number    string   `json:"number"`          // 号码 明确用户和群聊
	Priority  int      `json:"priority"`        // 决定 replies 排序，优先级越高排的越靠前
	Keyword   string   `json:"keyword"`         // 关键词，模糊查询
	Value     string   `json:"value"`           // 值，模糊查询
	CreatedAt int      `json:"created_at"`      // 创建时间
	Platforms []string `json:"platforms"`       // 平台
}

var replies []Reply //一切增删查改只需作用到这个变量
var repliesLock sync.RWMutex

func init() {
	REPLY.Foreach(func(b1, b2 []byte) error {
		repliesLock.Lock()
		defer repliesLock.Unlock()
		rp := Reply{}
		err := json.Unmarshal(b2, &rp)
		if err != nil {
			return nil
		}
		replies = append(replies, rp)
		sort.Slice(replies, func(i, j int) bool {
			return replies[i].Priority > replies[j].Priority
		})
		return nil
	})
	GinApi(GET, "/api/reply/list", func(ctx *gin.Context) {
		repliesLock.RLock()
		defer repliesLock.RUnlock()
		page, _ := strconv.Atoi(ctx.DefaultQuery("current", "1"))
		perPage, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))
		keyword := ctx.Query("keyword")
		value := ctx.Query("value")
		// class_ := ctx.Query("class")
		// class := utils.Int(class_)
		number := ctx.Query("number")
		// filter replies based on the query parameters
		filteredReplies := make([]Reply, 0, len(replies))
		for _, reply := range replies {
			if keyword != "" && !strings.Contains(reply.Keyword, keyword) {
				continue
			}
			if value != "" && !strings.Contains(reply.Value, value) {
				continue
			}
			// if class_ != "" && reply.Class != class {
			// 	continue
			// }
			if number != "" && reply.Number != number {
				continue
			}
			filteredReplies = append(filteredReplies, reply)
		}
		sort.Slice(filteredReplies, func(i, j int) bool {
			return filteredReplies[i].CreatedAt > filteredReplies[j].CreatedAt
		})
		// paginate the filtered replies
		start := (page - 1) * perPage
		end := start + perPage
		if end > len(filteredReplies) {
			end = len(filteredReplies)
		}
		paginatedReplies := filteredReplies[start:end]
		index := start + 1
		for i := range paginatedReplies {
			filteredReplies[i].Index = index
			index++
			if filteredReplies[i].Nickname == "" || len(filteredReplies[i].Platforms) == 0 {
				nk := Nickname{ID: filteredReplies[i].Number}
				nickname.First(&nk)
				if nk.Value != "" && filteredReplies[i].Nickname == "" {
					filteredReplies[i].Nickname = nk.Value
				}
				if nk.Platform != "" && len(filteredReplies[i].Platforms) == 0 {
					filteredReplies[i].Platforms = []string{nk.Platform}
				}
			}
		}
		ctx.JSON(200, map[string]interface{}{
			"success":   true,
			"data":      paginatedReplies,
			"page":      page,
			"total":     len(filteredReplies),
			"platforms": getPltsLabel(),
		})
	})

	GinApi(POST, "/api/reply", func(ctx *gin.Context) {
		repliesLock.Lock()
		defer repliesLock.Unlock()
		var reply Reply
		data, _ := ioutil.ReadAll(ctx.Request.Body)
		var v = map[string]interface{}{}
		if err := json.Unmarshal(data, &reply); err != nil {
			ctx.JSON(200, map[string]interface{}{
				"success":      false,
				"errorMessage": err.Error(),
			})
			return
		}
		json.Unmarshal(data, &v)
		has := func(str string) bool {
			_, ok := v[str]
			return ok
		}
		if reply.ID < 0 {
			reply.ID = 0
		}
		// find existing reply with the same ID
		var existingReply *Reply
		for i, r := range replies {
			if r.ID == reply.ID {
				existingReply = &replies[i]
				break
			}
		}
		if existingReply != nil {
			// update existing reply
			if has("nickname") {
				existingReply.Nickname = reply.Nickname
			}
			if has("number") {
				existingReply.Number = reply.Number
			}
			if has("keyword") {
				existingReply.Keyword = reply.Keyword
			}
			if has("value") {
				existingReply.Value = reply.Value
			}
			if has("priority") {
				existingReply.Priority = reply.Priority
			}
			if has("platforms") {
				existingReply.Platforms = reply.Platforms
			}
			reply = *existingReply
			err := REPLY.Create(&reply)
			if err != nil {
				ctx.JSON(200, map[string]interface{}{
					"success":      false,
					"errorMessage": err.Error(),
				})
				return
			}
		} else {
			reply.CreatedAt = int(time.Now().Unix())
			err := REPLY.Create(&reply)
			if err != nil {
				ctx.JSON(200, map[string]interface{}{
					"success":      false,
					"errorMessage": err.Error(),
				})
				return
			}
			replies = append(replies, reply)
		}
		sort.Slice(replies, func(i, j int) bool {
			return replies[i].Priority > replies[j].Priority
		})
		ctx.JSON(200, map[string]interface{}{
			"success": true,
		})
	})
	//删除功能
	GinApi(DELETE, "/api/reply", func(ctx *gin.Context) {
		repliesLock.Lock()
		defer repliesLock.Unlock()
		id := utils.Int(ctx.Query("id"))
		for i, r := range replies {
			if r.ID == id {
				REPLY.Set(r.ID, nil)
				replies = append(replies[:i], replies[i+1:]...)
				break
			}
		}
		ctx.JSON(200, map[string]interface{}{
			"success": true,
		})
	})
}

var REPLY = MakeBucket("reply")

// // 能处理字符：你好，我是${ user.name }
// func parseReply(str string) string {
// 	re := regexp.MustCompile(`\$\{\s*([^\s{}]+)\s*\}`)
// 	return re.ReplaceAllStringFunc(str, func(match string) string {
// 		bk := match[2 : len(match)-1]
// 		b_k := strings.Split(bk, ".")
// 		if len(b_k) != 3 {
// 			return fmt.Sprintf("${%s}", bk)
// 		}
// 		return MakeBucket(b_k[1]).GetString(b_k[2])
// 	})
// }

// 能处理字符：你好，我是${ user.name ?? 6 }
func parseReply2(str string) string {
	re := regexp.MustCompile(`\$\{\s*([^{}]+)\s*\}`)
	return re.ReplaceAllStringFunc(str, func(match string) string {
		script := match[2 : len(match)-1]
		script = regexp.MustCompile(`(\w+)\.(\w+)`).ReplaceAllStringFunc(script, func(match string) string {
			parts := strings.Split(match, ".")
			return fmt.Sprintf(`Bucket("%s")["%s"]`, parts[0], parts[1])
		})
		vm := goja.New()
		vm.Set("Bucket", func(name string) interface{} {
			return JsBucket(vm, name, "", false)
		})
		v, err := vm.RunString(script)
		if err == nil {
			return v.String()
		}
		return match
		// b_k := strings.Split(bk, ".")
		// if len(b_k) != 3 {
		// 	return fmt.Sprintf("${%s}", bk)
		// }
		// return MakeBucket(b_k[1]).GetString(b_k[2])
		// return ""
	})
}

func parseReply3(str string, f func(string, string)) string {
	ks := map[string]bool{}
	re := regexp.MustCompile(`\$\{\s*([^{}]+)\s*\}`)
	return re.ReplaceAllStringFunc(str, func(match string) string {
		script := match[2 : len(match)-1]
		script = regexp.MustCompile(`(\w+)\.(\w+)`).ReplaceAllStringFunc(script, func(match string) string {
			parts := strings.Split(match, ".")
			k := fmt.Sprintf(`%s.%s`, parts[0], parts[1])
			if _, ok := ks[k]; !ok {
				ks[k] = true
				f(parts[0], parts[1])
			}
			return fmt.Sprintf(`Bucket("%s")["%s"]`, parts[0], parts[1])
		})
		vm := goja.New()
		vm.Set("Bucket", func(name string) interface{} {
			return JsBucket(vm, name, "", false)
		})
		v, err := vm.RunString(script)
		if err == nil {
			return v.String()
		}
		return match
		// b_k := strings.Split(bk, ".")
		// if len(b_k) != 3 {
		// 	return fmt.Sprintf("${%s}", bk)
		// }
		// return MakeBucket(b_k[1]).GetString(b_k[2])
		// return ""
	})
}
