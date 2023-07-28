package core

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/cdle/sillyGirl/core/common"
	"github.com/cdle/sillyGirl/core/storage"
	"github.com/cdle/sillyGirl/utils"
	"github.com/gin-gonic/gin"
)

var plugin_subcribe_addresses = sillyGirl.GetString("plugin_subcribe_addresses")
var plugin_subcribe_data = MakeBucket("plugin_subcribe_data")

type RequestPluginResult struct {
	Success bool               `json:"success"`
	Data    []*common.Function `json:"data"`
	Page    int                `json:"page"`
	Total   int                `json:"total"`
	Tab1    int                `json:"tab1"`
	Tab2    int                `json:"tab2"`
	Tab3    int                `json:"tab3"`
	Tab     string             `json:"tab"`
	Time    time.Time          `json:"time"`
	Classes map[string]int     `json:"classes"`
	Origins map[string]string  `json:"origins"`
}

var plugin_list = []*common.Function{}

var cdle_sublink = `
//傻妞官方
link://T4EywWN46ztYBhHNdOl6Tpz8QQsCZGj8JvdRJ5QKatJm0P+mI/G3ruO7AC04guqqKKa29VOvTGR7ATUJGYayRBpG2RFq+6ZPK3vcu6KCDGvRE3S43Gj42EXfvs04M6s4
//大灰机
link://T4EywWN46ztYBhHNdOl6ThvHulsS6Fo5vRI+WFDJEtMBltBIFj2gLoSSIFXLSRmeAwYxIkikr+TZUzuTr2QYZ7edh12jsIgAv3s0FR3pqace1TX5/6r2rcc52HlAkCPU
`

// 搬运中心
// link://T4EywWN46ztYBhHNdOl6Tpz8QQsCZGj8JvdRJ5QKatJYds3a/BticqD0hzidGsOysEx/RK/nKppChxMLb6QGczhWjGC/M2ETxWb+Jl+6q/x+LP4gy+ibeAEzatOYwdZMckI8nN/R6mY/HW2dyBtp0qH1ldICn6Wl+9YowLvvpLU=
// 木子李
// link://T4EywWN46ztYBhHNdOl6Tpz8QQsCZGj8JvdRJ5QKatL/GWakSkUWVNTd/jJS4YaqGXqvoJOxtEwVxbfBpmsMdTpKFr7K/+9MW/CJFpFsLFGM3yRxh2z8fVsDZUV6GoXei5QhOviIvo5ys7N5b6MRiEmbVATiiTEovz3IBg8nObQ=
func initPluginList() {
	list := []*common.Function{}
	var carrys []chan []*common.Function
	for _, v := range regexp.MustCompile(`link://([^\s#]+)`).FindAllStringSubmatch(cdle_sublink+"\n"+plugin_subcribe_addresses+"\n", -1) {
		sublink := v[1]
		ppr := common.PluginPublisher{}
		var data []byte
		func() {
			defer func() {
				err := recover()
				if err != nil {
					console.Error("initPluginList：", err)
				}
			}()
			data, _ = DecryptByAes(sublink)
		}()

		if data == nil {
			continue
		}
		json.Unmarshal(data, &ppr)
		if ppr.Address != "" {
			carry := make(chan []*common.Function)
			carrys = append(carrys, carry)
			go func() {
				rr := RequestPluginResult{}
				data := plugin_subcribe_data.GetBytes(sublink)
				json.Unmarshal(data, &rr)
				if !rr.Success || rr.Time.Add(time.Second*3).Before(time.Now()) {
					address := ""
					if !strings.HasSuffix(ppr.Address, "list.json") {
						address = ppr.Address + "/api/plugins/list.json"
					} else {
						address = ppr.Address
					}
					req := httplib.Get(address)
					req.SetTimeout(time.Second*2, time.Second*2)
					data, _ := req.Bytes()
					json.Unmarshal(data, &rr)
					if rr.Success {
						rr.Time = time.Now()
						plugin_subcribe_data.Set(sublink, string(utils.JsonMarshal(rr)))
					}
				}
				for i := range rr.Data {
					rr.Data[i].Address = ppr.Address
					rr.Data[i].Organization = ppr.Organization
					rr.Data[i].Identified = ppr.Identified
				}
				n := len(rr.Data)
				flag := true
				for i := 0; i < n && flag; i++ {
					flag = false
					for j := 0; j < n-i-1; j++ {
						if rr.Data[j].CreateAt < rr.Data[j+1].CreateAt {
							rr.Data[j], rr.Data[j+1] = rr.Data[j+1], rr.Data[j]
							flag = true
						}
					}
				}
				carry <- rr.Data
			}()
		}
	}
	for _, carry := range carrys {
		list = append(list, <-carry...)
	}
	cyzl := "7642f5de-3300-11ed-8a79-52540066b468"
	plugin_list = list
	if sillyGirl.GetString("password") == "" && plugins.GetString(cyzl) == "" { //自动安装老版命令
		plugins.Set(cyzl, "install")
	}
	// if plugins.GetString("78b15932-334f-11ed-8b59-aaaa00117a5c") == "" { //自动安装比价文案
	// 	plugins.Set("78b15932-334f-11ed-8b59-aaaa00117a5c", "install")
	// }
}

var plugin_downloads = MakeBucket("plugin_downloads")

func initWebPluginList() {
	storage.Watch(sillyGirl, "plugin_subcribe_addresses", func(old, new, key string) *storage.Final {
		plugin_subcribe_addresses = new
		return nil
	})
	GinApi(GET, "/api/plugins/list.json", func(ctx *gin.Context) {
		// ctx.QueryArray()
		origins := ctx.QueryArray("origin[]")
		current := utils.Int(ctx.Query("current"))
		pageSize := utils.Int(ctx.Query("pageSize"))
		activeKey := ctx.Query("activeKey")
		init := ctx.Query("init")
		keyword := ctx.Query("keyword")
		class := ctx.Query("class")
		mclass := ctx.Query("mclass")
		rr := RequestPluginResult{
			Success: true,
		}
		if pageSize == 0 {
			pageSize = 10
		}
		if class == "" {
			class = "全部"
		}
		rr.Page = current
		rr.Data = []*common.Function{}
		if current != 0 {
			if current == 1 && init != "false" {
				initPluginList()
			}
			var list []*common.Function
			if keyword == "" {
				if len(origins) == 0 {
					list = append(list, plugin_list...)

				} else {
					for _, f := range plugin_list {
						if Contains(origins, f.Organization) {
							list = append(list, f)
						}
					}
				}
			} else {
				if len(origins) == 0 {
					for _, f := range plugin_list {
						if strings.Contains(f.Title, keyword) || strings.Contains(f.Organization, keyword) {
							list = append(list, f)
						}
					}
				} else {
					for _, f := range plugin_list {
						if strings.Contains(f.Title, keyword) || strings.Contains(f.Organization, keyword) {
							if Contains(origins, f.Organization) {
								list = append(list, f)
							}
						}
					}
				}

			}
			rr.Total = len(list)
			tab1 := []*common.Function{}
			tab2 := []*common.Function{}
			tab3 := []*common.Function{}
			fc := []*common.Function{}
			fc = append(fc, Functions...)
			classes := map[string][]*common.Function{}
			classesNum := map[string]int{}
			for i := range list {
				if len(list[i].Classes) == 0 {
					class := "未分类"
					if _, ok := classes[class]; !ok {
						classes[class] = []*common.Function{}
					}
					classes[class] = append(classes[class], list[i])
				} else {
					for _, class := range list[i].Classes {
						class = strings.TrimRight(class, "类")
						if _, ok := classes[class]; !ok {
							classes[class] = []*common.Function{}
						}
						classes[class] = append(classes[class], list[i])
					}
				}
			}

			for class, fs := range classes {
				classesNum[class] = len(fs)
			}
			classesNum["全部"] = len(list)
			if class != "全部" {
				list, _ = classes[class]
			}
			rr.Classes = classesNum
			var origins = map[string]string{}
			for i := range list { //处理第二分类
				if list[i].Organization != "" {
					origins[list[i].Organization] = list[i].Organization
				}
				ded := false
				for j := range fc {
					if list[i].UUID == fc[j].UUID {
						if list[i].Version != fc[j].Version {
							tab3 = append(tab3, list[i])
						}
						ded = true
						break
					}
				}
				if ded {
					tab1 = append(tab1, list[i]) //已安装
				} else {
					tab2 = append(tab2, list[i])
				}
			}
			rr.Origins = origins
			if activeKey == "tab2" {
				list = tab2
				rr.Tab1 = len(tab1)
				rr.Tab2 = len(tab2)
				rr.Tab3 = len(tab3)
			} else if activeKey == "tab3" {
				list = tab3
				rr.Tab1 = len(tab1)
				rr.Tab2 = len(tab2)
				rr.Tab3 = len(tab3)
			} else {
				list = tab1
				rr.Tab1 = len(tab1)
				rr.Tab2 = len(tab2)
				rr.Tab3 = len(tab3)
			}
			tab := ""
			if mclass == "true" {
				if rr.Tab2 > rr.Tab1 {
					list = tab2
					tab = "tab2"
				} else {
					list = tab1
					tab = "tab1"
				}
			}
			rr.Tab = tab
			rr.Total = len(list)
			if len(list) == 0 {
				ctx.JSON(200, rr)
				return
			}
			if last := (rr.Total + pageSize - 1) / pageSize; current > last {
				current = last
			}
			begin := (current - 1) * pageSize
			end := (current) * pageSize
			if end > rr.Total {
				end = rr.Total
			}
			if begin > end {
				begin = end
			}
			rr.Data = append(rr.Data, list[begin:end]...)
			publics := []string{}
			for _, f := range Functions {
				if f.Public && f.UUID != "" {
					publics = append(publics, f.UUID)
				}
			}
			for i := range rr.Data {
				rr.Data[i].HasForm = false
				rr.Data[i].Running = false
				for j := range fc {
					if rr.Data[i].UUID == fc[j].UUID {
						rr.Data[i].Messages = GetPluginMessage(rr.Data[i].UUID)
						if rr.Data[i].Version != fc[j].Version {
							rr.Data[i].Status = 1
						} else {
							rr.Data[i].Status = 2
						}
						if Contains(publics, rr.Data[i].UUID) {
							rr.Data[i].Status = 6
						}
						if rr.Data[i].Icon == "" {
							rr.Data[i].Icon = "https://blog.example.com/huli.jpeg"
						}
						if fc[j].HasForm {
							rr.Data[i].HasForm = true
						}
						if fc[j].Running {
							rr.Data[i].Running = true
						}
						rr.Data[i].Debug = plugin_debug.GetString(rr.Data[i].UUID) == "b:true"
					}
				}
				rr.Data[i].Description = parseReply2(rr.Data[i].Description)
			}

			ctx.JSON(200, rr)
			return
		}

		ctx.JSON(200, GetPublicResponse())
	})
}

func GetPublicResponse() *RequestPluginResult {
	rr := &RequestPluginResult{
		Success: true,
	}
	fs := []*common.Function{}
	for _, f := range Functions {
		if f.Public {
			fs = append(fs, f)
			f.Downloads = plugin_downloads.GetInt(f.UUID)
		}
	}
	rr.Total = len(fs)
	rr.Data = fs
	rr.Page = 1
	return rr
}
