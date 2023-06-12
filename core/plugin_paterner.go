package core

import "regexp"

func getPaterner(uuid, path string) {
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
	if address != "" {
		for _, f := range ls {
			if f.Address == address && (f.Title == path || f.UUID == path) {
				data = plugins.GetBytes(f.UUID)
				if data != nil {
					return
				}
			}
		}
		if data == nil {
			for _, l := range ls {
				if l.Address == address && (l.Title == path || l.UUID == path) {
					data = fetchScript(l.Address, l.UUID)
					if data == nil {
						console.Warn("无法从订阅源获取 %s 的协作脚本 %s ", GetScriptNameByUUID(uuid), path)
					} else {
						console.Log("已从订阅源获取 %s 的协作脚本 %s", GetScriptNameByUUID(uuid), path)
						plugins.Set(l.UUID, string(data))
					}
					return
				}
			}
		}
	}
	if data == nil {
		console.Warn("找不到 %s 的协作脚本 %s", GetScriptNameByUUID(uuid), path)
	}
}
