package core

import (
	"fmt"
	"strings"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/dop251/goja_nodejs/require"
)

func mapFileSystemSourceLoader(uuid string) require.SourceLoader {
	return func(path string) ([]byte, error) {
		path = strings.ReplaceAll(path, `node_modules/`, "")
		var data []byte
		var address = ""
		ls := plugin_list
		for _, f := range ls {
			if f.UUID == uuid {
				address = f.Address
				break
			}
		}
		if address != "" {
			for _, f := range ls {
				if f.Address == address && f.Title == path {
					data = plugins.GetBytes(f.UUID)
				}
			}
			if data == nil {
				for _, l := range ls {
					if l.Address == address && l.Title == path {
						data = fetchScript(l.Address, l.UUID)
						if data == nil {
							return nil, fmt.Errorf("无法从订阅源获取%s模块", path)
						} else {
							console.Log("已从订阅源获取%s模块", path)
							plugins.Set(l.UUID, string(data))
						}
						break
					}
				}
			}
		}
		if len(data) == 0 {
			fs := Functions
			for _, f := range fs {
				if f.Title == path {
					data = plugins.GetBytes(f.UUID)
				}
			}
		}
		if data == nil {
			return nil, fmt.Errorf("缺少%s模块", path) //require.ModuleFileDoesNotExistError
		}
		su := &ScriptUtils{
			script: string(data),
		}
		if su.GetValue("encrypt") == "true" {
			data = []byte(DecryptPlugin(su.script))
		}
		data = []byte(halfDeEct(string(data)))
		return data, nil
	}
}

func fetchScript(address, uuid string) (data []byte) {
	var prefix = "?uuid=" + uuid
	if !strings.HasSuffix(address, "list.json") {
		address = address + "/api/plugins/download" + prefix
	} else {
		address = strings.ReplaceAll(address, "list.json", "download"+prefix)
	}
	data, _ = httplib.Get(address).Bytes()
	return
}
