package core

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/cdle/sillyGirl/core/common"
	"github.com/cdle/sillyGirl/core/storage"
	"github.com/cdle/sillyGirl/utils"
	"github.com/goccy/go-json"
)

var plugin_path = "/etc/sillyGirl/public/"
var plugin_download_file = plugin_path + "download"

func CheckPluginAddress(address string) error {
	if !strings.HasSuffix(address, "list.json") {
		address += "/api/plugins/list.json"
	}
	data, _ := httplib.Get(address).Bytes()
	rr := RequestPluginResult{}
	json.Unmarshal(data, &rr)
	if !rr.Success {
		return errors.New("无效的地址")
	}
	return nil
}

func initPluginPublish() {
	storage.Watch(sillyGirl, "plugin_sublink", func(old, address, key string) *storage.Final {
		if strings.HasPrefix(address, "sub://") {
			return nil
		}
		if err := CheckPluginAddress(address); err != nil {
			return &storage.Final{
				Error: err,
			}
		}
		str, err := EncryptByAes(utils.JsonMarshal(common.PluginPublisher{
			Address: address,
			// MachineID: GetMachineID(),
		}))
		sublink := fmt.Sprintf("sub://%s", str)
		return &storage.Final{
			Now:     sublink,
			Message: sublink,
			Error:   err,
		}
	})

	os.MkdirAll(plugin_download_file, 0666)
	os.WriteFile(plugin_path+"list.json", utils.JsonMarshal(GetPublicResponse()), 0666)
	for _, f := range Functions {
		if f.UUID != "" && f.Public {
			os.WriteFile(fmt.Sprintf("%s/%s.js", plugin_download_file, f.UUID), []byte(publicScript(plugins.GetString(f.UUID))), 0666)
		}
	}
}

func publicScript(str string) string {
	su := &ScriptUtils{
		script: str,
	}
	if version := su.GetValue("version"); regexp.MustCompile(`v\d+\.\d+\.\d`).FindString(version) != version {
		su.SetValue("version", "v1.0.0")
	}
	if su.GetValue("author") == "" {
		su.SetValue("author", "佚名")
	}
	if su.GetValue("description") == "" {
		su.SetValue("description", "🐒这个人很懒什么都没有留下")
	}
	if su.GetValue("public") == "true" {
		su.SetValue("public", "false")
	}
	if su.GetValue("title") == "" {
		su.SetValue("title", "无名脚本")
	}
	if su.GetValue("message") != "" {
		su.DeleteValue("message")
	}
	create_at := su.GetValue("create_at")
	if _, err := time.Parse("2006-01-02 15:04:05", create_at); create_at == "" || err != nil {
		su.SetValue("create_at", time.Now().Format("2006-01-02 15:04:05"))
	}
	if su.GetValue("encrypt") == "true" {
		su.script = EncryptPlugin(su.script)
	}
	su.script = halfEct(su.script)
	return su.script
}
