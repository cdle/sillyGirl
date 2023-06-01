package core

import (
	"fmt"

	"github.com/cdle/sillyGirl/core/common"
	"github.com/cdle/sillyGirl/utils"
)

var IsCdle = false

func init() {
	if sillyGirl.GetString("is_cdle") == "silly8023" {
		IsCdle = true
	}
	if !IsCdle {
		return
	}
	AddCommand([]*common.Function{ //认证订阅
		{
			Admin: true,
			Rules: []string{"identify sublink [地址] [组织]"},
			Handle: func(s common.Sender) interface{} {
				address := s.Get(0)
				organization := s.Get(1)
				// machine_id := s.Get(1)
				if err := CheckPluginAddress(address); err != nil {
					return err
				}
				str, err := EncryptByAes(utils.JsonMarshal(common.PluginPublisher{
					Address:      address,
					Organization: organization,
					Identified:   true,
					// MachineID:    machine_id,
				}))
				if err != nil {
					return err
				}
				sublink := fmt.Sprintf("link://%s", str)
				return sublink
			},
		},
	})

	// 原创记录
	// var pcr = MakeBucket("pcr")
	// GinApi(POST, "/api/plugins/record", func(c *gin.Context) {
	// 	data, _ := ioutil.ReadAll(c.Request.Body)
	// 	str, _ := EncryptByAes(data)
	// 	v := PluginCreateRecord{}
	// 	if json.Unmarshal([]byte(str), &v) == nil {
	// 		o := v
	// 		pcr.First(o)
	// 		if o.MachineID != "" {
	// 			pcr.Create(v)
	// 		}
	// 	}
	// })
}

// type PluginCreateRecord struct {
// 	ID        string `json:"id"`
// 	Unix      int64  `json:"unix"`
// 	MachineID string `json:"machine_id"`
// 	IP        string `json:"ip"`
// 	Title     string `json:"title"`
// }
