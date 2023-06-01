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
	Time    time.Time          `json:"time"`
}

var plugin_list = []*common.Function{}

var cdle_sublink = "link://T4EywWN46ztYBhHNdOl6TjL4plwqQWRUoqr8w0KFmMqAdblZX3/xtrZARf3VKKQmH6iQNfyWvB2bqf6P1n/CMh1KLHLbTvUzh9zBQS2u9GeYwAp0APEZvQV1O6pb5g9V/dd6TLH54ssD92DAuMa1xw=="

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

func initWebPluginList() {
	storage.Watch(sillyGirl, "plugin_subcribe_addresses", func(old, new, key string) *storage.Final {
		plugin_subcribe_addresses = new
		return nil
	})
	GinApi(GET, "/api/plugins/download", func(c *gin.Context) {
		uuid := c.Query("uuid")
		for _, f := range Functions {
			if f.UUID == uuid && f.Public {
				c.String(200, publicScript(plugins.GetString(f.UUID)))
				return
			}
		}
	})
	GinApi(GET, "/api/plugins/list.json", func(ctx *gin.Context) {
		current := utils.Int(ctx.Query("current"))
		pageSize := utils.Int(ctx.Query("pageSize"))
		init := ctx.Query("init")
		title := ctx.Query("title")
		rr := RequestPluginResult{
			Success: true,
		}
		if pageSize == 0 {
			pageSize = 10
		}
		rr.Page = current
		if current != 0 {
			var list []*common.Function
			if title == "" {
				list = plugin_list
			} else {
				for _, f := range plugin_list {
					if strings.Contains(f.Title, title) || strings.Contains(f.Organization, title) {
						list = append(list, f)
					}
				}
			}
			if current == 1 && init != "false" {
				initPluginList()
			}
			rr.Total = len(list)
			begin := (current - 1) * pageSize
			end := (current) * pageSize
			if end > rr.Total {
				end = rr.Total
			}
			if begin > end {
				begin = end
			}
			rr.Data = append(rr.Data, list[begin:end]...)
			fc := []*common.Function{}
			fc = append(fc, Functions...)
			publics := []string{}
			for _, f := range Functions {
				if f.Public && f.UUID != "" {
					publics = append(publics, f.UUID)
				}
			}
			for i := range rr.Data {
				for j := range fc {
					if rr.Data[i].UUID == fc[j].UUID {
						if rr.Data[i].Version != fc[j].Version {
							rr.Data[i].Status = 1
						} else {
							rr.Data[i].Status = 2
						}
						if Contains(publics, rr.Data[i].UUID) {
							rr.Data[i].Status = 6
						}
						// if rr.Data[i].MachineID == machine_id {
						// 	rr.Data[i].MachineID = ""
						// }
					}
				}
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
		}
	}
	rr.Total = len(fs)
	rr.Data = fs
	rr.Page = 1
	return rr
}
