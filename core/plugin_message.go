package core

import (
	"encoding/json"
	"time"

	"github.com/cdle/sillyGirl/utils"
)

type PMsg struct {
	Class   string `json:"class"`
	Unix    int    `json:"unix"`
	Content string `json:"content"`
}

var plugin_messages = MakeBucket("plugin_messages")

func WritePluginMessage(uuid string, class string, content string) {
	if uuid == "" {
		return
	}
	var data = plugin_messages.GetBytes(uuid)
	pmsgs := []PMsg{}
	json.Unmarshal(data, &pmsgs)
	s := &Strings{}
	ok := false
	new := PMsg{
		Class:   class,
		Unix:    int(time.Now().Unix()),
		Content: content,
	}
	for i := range pmsgs {
		if pmsgs[i].Class != class {
			continue
		}
		if content == "" {
			continue
		}
		if s.Similarity(pmsgs[i].Content, content) > 0.9 {
			pmsgs[i] = new
			ok = true
			break
		}
	}
	if !ok {
		pmsgs = append(pmsgs, new)
	}
	plugin_messages.Set(uuid, utils.JsonMarshal(pmsgs))
}

func GetPluginMessage(uuid string) []PMsg {
	var data = plugin_messages.GetBytes(uuid)
	pmsgs := []PMsg{}
	json.Unmarshal(data, &pmsgs)
	return pmsgs
}
