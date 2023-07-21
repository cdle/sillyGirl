package core

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/cdle/sillyGirl/core/common"
	"github.com/cdle/sillyGirl/core/storage"
	"github.com/cdle/sillyGirl/utils"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

var tasks = MakeBucket("tasks")

type TasksResult struct {
	Success bool      `json:"success"`
	Data    []*Tasks  `json:"data"`
	Page    int       `json:"page"`
	Total   int       `json:"total"`
	Time    time.Time `json:"time"`
}

type Sender struct {
	ChatID   string `json:"chat_id"`
	UserID   string `json:"user_id"`
	Platfrom string `json:"platform"`
	BotID    string `json:"bot_id"`
}

type Tasks struct {
	Index     int           `json:"id"`       //编号 顺序编号
	ID        string        `json:"task_id"`  //任务ID
	Title     string        `json:"title"`    //任务名
	Schedule  string        `json:"schedule"` //计划时间
	Senders   []Sender      `json:"senders"`  //发送人
	Command   string        `json:"command"`  //消息指令
	Scripts   []string      `json:"scripts"`  //触发脚本
	CronID    int           `json:"cron_id"`
	CreatedAt int           `json:"created_at"` //创建时间戳(秒)转换成日期
	Remark    string        `json:"remark"`
	Enable    bool          `json:"enable"`
	Handle    func()        `json:"-"`
	Icons     []interface{} `json:"icons"`
}

var pts = []*Tasks{}

func RegistTasks(pt *Tasks) {
	pt.Handle = func() {
		content := pt.Command
		for _, meta := range pt.Senders {
			adapter, _ := GetAdapter(meta.Platfrom, meta.BotID)
			if adapter != nil {
				sender := adapter.Sender2(nil)
				sender.SetFsps(&common.FakerSenderParams{
					Content: content,
					ChatID:  meta.ChatID,
					UserID:  meta.UserID,
				})
				for _, script := range pt.Scripts {
					for _, function := range Functions {
						if function.UUID == script {
							for i := range function.Rules {
								reg, err := regexp.Compile(function.Rules[i])
								if err == nil {
									if res := reg.FindStringSubmatch(content); len(res) > 0 {
										sender.SetMatch(res[1:])
										sender.SetParams(function.Params[i])
									}
								}
							}
							function.Handle(sender, nil)
							break
						}
					}
				}
			}
		}
	}
	cid, _ := CRON.AddFunc(pt.Schedule, pt.Handle)
	pt.CronID = int(cid)
	console.Debug("已添加计划任务：%s(%v)", pt.Title, pt.CronID)
}

func init() {
	tasks.Foreach(func(b1, b2 []byte) error {
		pt := Tasks{}
		err := json.Unmarshal(b2, &pt)
		if err != nil {
			return nil
		}
		RegistTasks(&pt)
		pts = append(pts, &pt)
		return nil
	})
	sort.Sort(byCreatedAt2(pts))
	for i := range pts {
		pts[i].Index = i + 1
	}
	storage.Watch(tasks, nil, func(old, new, key string) *storage.Final {
		console.Log("已更新计划任务")
		ocg := Tasks{}
		ncg := Tasks{}
		json.Unmarshal([]byte(old), &ocg)
		json.Unmarshal([]byte(new), &ncg)
		tmp := pts
		if old != "" {
			if new == "" { // 删除
				if ocg.ID != "" {
					for i, cg := range tmp {
						if cg.ID == ocg.ID {
							CRON.Remove(cron.EntryID(tmp[i].CronID))
							tmp = append(tmp[:i], tmp[i+1:]...)
							break
						}
					}
				} else {
					return nil
				}
			} else { // 修改
				if ocg.ID != "" {
					for i, cg := range tmp {
						if cg.ID == ocg.ID {
							CRON.Remove(cron.EntryID(tmp[i].CronID))
							tmp[i] = &ncg
							RegistTasks(&ncg)
							//todo 增
							break
						}
					}
				} else {
					return nil
				}
			}
		} else { //创建
			if ncg.ID != "" {
				tmp = append(tmp, &ncg)
				RegistTasks(&ncg)
				//todo 增
			} else {
				return nil
			}
		}
		sort.Sort(byCreatedAt2(pts))
		for i := range tmp {
			tmp[i].Index = i + 1
		}
		pts = tmp
		return nil
	})
	GinApi(GET, "/api/tasks", RequireAuth, func(ctx *gin.Context) {
		current := utils.Int(ctx.Query("current"))
		pageSize := utils.Int(ctx.Query("pageSize"))
		rr := TasksResult{
			Success: true,
		}
		pts := pts
		rr.Total = len(pts)
		if current == 0 {
			current = 1
		}
		if pageSize == 0 {
			pageSize = 20
		}
		begin := (current - 1) * pageSize
		end := (current) * pageSize
		if end > rr.Total {
			end = rr.Total
		}
		if begin > end {
			begin = end
		}
		rr.Data = pts[begin:end]
		for i := range rr.Data {
			rr.Data[i].Icons = []interface{}{}
			for _, script := range rr.Data[i].Scripts {
				for _, f := range Functions {
					if f.UUID == script {
						if f.Icon != "" {
							rr.Data[i].Icons = append(rr.Data[i].Icons, map[string]interface{}{
								"link":  f.Icon,
								"title": f.Title,
							})
							break
						}
					}
				}
			}
		}
		ctx.JSON(200, rr)
	})
	GinApi(POST, "/api/tasks", RequireAuth, func(ctx *gin.Context) {
		// 将请求的 JSON 数据解析为一个 map[string]interface{} 类型的变量
		var updateData map[string]interface{}
		err := ctx.BindJSON(&updateData)
		if err != nil {
			ctx.JSON(200, map[string]interface{}{
				"success":      false,
				"errorMessage": err.Error(),
			})
			return
		}
		v, ok := updateData["task_id"]
		if !ok {
			ctx.JSON(200, map[string]interface{}{
				"success":      false,
				"errorMessage": "任务ID不能为空",
			})
			return
		}
		task_id := v.(string)
		var tp = Tasks{
			ID: task_id,
		}
		tasks.First(&tp)
		for key, value := range updateData {
			switch key {
			case "title":
				if v, ok := value.(string); ok {
					tp.Title = v
				}
			case "remark":
				if v, ok := value.(string); ok {
					tp.Remark = v
				}
			case "schedule":
				if v, ok := value.(string); ok {
					tp.Schedule = v
					id, err := CRON.AddFunc(tp.Schedule, func() {})
					if err != nil {
						ctx.JSON(200, map[string]interface{}{
							"success":      false,
							"errorMessage": "Cron表达式错误：" + err.Error(),
						})
						return
					}
					CRON.Remove(id)
				}
			case "senders":
				ss := []Sender{}
				err := json.Unmarshal(utils.JsonMarshal(value), &ss)
				if err != nil {
					ctx.JSON(200, map[string]interface{}{
						"success":      false,
						"errorMessage": "Senders错误：" + err.Error(),
					})
					return
				}
				tp.Senders = ss
			case "command":
				if v, ok := value.(string); ok {
					tp.Command = v
				}
			case "scripts":
				if v, ok := value.([]interface{}); ok {
					tp.Scripts = toStringSlice(v)
				}
			case "enable":
				if v, ok := value.(bool); ok {
					tp.Enable = v
				}
			}
		}
		if tp.CreatedAt == 0 {
			tp.CreatedAt = int(time.Now().Unix())
		}
		tasks.Set(task_id, utils.JsonMarshal(tp))
		if err != nil {
			ctx.JSON(200, map[string]interface{}{
				"success":      false,
				"errorMessage": err.Error(),
			})
			return
		}
		ctx.JSON(200, map[string]interface{}{
			"success": true,
		})
	})
	GinApi(DELETE, "/api/tasks", RequireAuth, func(ctx *gin.Context) {
		pt := &Tasks{}
		err := ctx.BindJSON(pt)
		if err != nil {
			ctx.JSON(200, map[string]interface{}{
				"success":      false,
				"errorMessage": err.Error(),
			})
			return
		}
		if pt.ID == "" {
			ctx.JSON(200, map[string]interface{}{
				"success":      false,
				"errorMessage": "任务ID不为空",
			})
			return
		}
		tasks.Set(pt.ID, "")
		ctx.JSON(200, map[string]interface{}{
			"success": true,
		})
	})
	GinApi(GET, "/api/task/selects", RequireAuth, func(ctx *gin.Context) {
		var scripts = map[string]string{}
		var task_id = ctx.Query("task_id")
		var pts = pts
		var chat_ids = []string{}
		var user_ids = []string{}
		for _, pt := range pts {
			if pt.ID == task_id {
				for _, s := range pt.Senders {
					if s.ChatID != "" {
						chat_ids = append(chat_ids, s.ChatID)
					}
					if s.UserID != "" {
						user_ids = append(user_ids, s.UserID)
					}
				}
				break
			}
		}
		functions := Functions
		for _, function := range functions {
			if function.UUID != "" {
				scripts[function.UUID] = function.Title + function.Suffix
			}
		}
		var user_names = []NicklabeL{}
		var group_names = []NicklabeL{{
			Label: "私聊",
			Value: "",
		}}
		nickname.Foreach(func(b1, b2 []byte) error {
			v := &Nickname{}
			code := string(b1)
			err := json.Unmarshal(b2, v)
			if err == nil {
				if v.Group {
					if Contains(chat_ids, code) {
						group_names = append(group_names, NicklabeL{
							Label: fmt.Sprintf("%s(%s)", v.Value, code),
							Value: code,
						})
					}
				} else {
					if Contains(user_ids, code) {
						user_names = append(user_names, NicklabeL{
							Label: fmt.Sprintf("%s(%s)", v.Value, code),
							Value: code,
						})
					}

				}
			}
			return nil
		})
		platforms := map[string][]string{}
		for _, plt := range getPltsArray() {
			platforms[plt] = GetAdapterBotsID(plt)
		}
		ctx.JSON(200, map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"scripts":     scripts,
				"platforms":   platforms,
				"user_names":  user_names,
				"group_names": group_names,
			},
		})
	})
	GinApi(GET, "/api/tasks/run", RequireAuth, func(ctx *gin.Context) {
		var task_id = ctx.Query("task_id")
		for _, pt := range pts {
			if pt.ID == task_id {
				pt.Handle()
			}
		}
		ctx.JSON(200, map[string]interface{}{
			"success": true,
		})
	})

}

type byCreatedAt2 []*Tasks

func (s byCreatedAt2) Len() int {
	return len(s)
}

func (s byCreatedAt2) Less(i, j int) bool {
	return s[i].CreatedAt > s[j].CreatedAt
}

func (s byCreatedAt2) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
