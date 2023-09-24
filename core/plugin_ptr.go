package core

import "regexp"

func getPaterner(uuid, path string) (string, string, bool) {
	var ss = regexp.MustCompile(`\S+`).FindAllString(path, -1)
	path = ss[0]
	var data []byte
	var address = ""
	ls := plugin_list
	for _, f := range ls {
		if f.UUID == uuid {
			address = f.Address
			break
		}
	}
	for _, f := range ls { //同源
		if f.Address == address && (f.Title == path || f.UUID == path) {
			data = plugins.GetBytes(f.UUID)
			if data != nil {
				return string(data), f.UUID, true //本地取
			}
			data = fetchScript(f.Address, f.UUID)
			if data == nil {
				console.Warn("无法从订阅源获取 %s 的协作脚本 %s ", GetScriptNameByUUID(uuid), path)
				return "", f.UUID, false
			} else {
				console.Log("已从订阅源获取 %s 的协作脚本 %s", GetScriptNameByUUID(uuid), path)
				plugins.Set(f.UUID, string(data))
				return string(data), f.UUID, true
			}
		}
	}
	for _, f := range ls { //异源
		if f.Address != address && (f.Title == path || f.UUID == path) {
			data = plugins.GetBytes(f.UUID)
			if data != nil {
				return string(data), f.UUID, true //本地取
			}
			data = fetchScript(f.Address, f.UUID)
			if data == nil {
				console.Warn("无法从订阅源获取 %s 的协作脚本 %s ", GetScriptNameByUUID(uuid), path)
				return "", f.UUID, false
			} else {
				console.Log("已从订阅源获取 %s 的协作脚本 %s", GetScriptNameByUUID(uuid), path)
				plugins.Set(f.UUID, string(data))
				return string(data), f.UUID, true
			}
		}
	}
	return "", "", false
}
