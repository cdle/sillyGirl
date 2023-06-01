package core

import (
	"encoding/json"

	"github.com/cdle/sillyGirl/utils"
)

var platforms = []string{}

func init() {
	nickname.Foreach(func(b1, b2 []byte) error {
		v := &Nickname{}
		err := json.Unmarshal(b2, v)
		if err == nil {
			platforms = append(platforms, v.Platform)
		}
		return nil
	})
	platforms = utils.Unique(platforms)
}

func getPltsArray() []string {
	return utils.Unique(platforms, GetAdapterBotPlts())
}

func getPltsLabel() []map[string]string {
	ms := []map[string]string{}
	for _, v := range getPltsArray() {
		ms = append(ms, map[string]string{"label": v, "value": v})
	}
	return ms
}
