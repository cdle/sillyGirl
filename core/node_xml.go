package core

import (
	"github.com/clbanning/mxj"
)

func (sender *Strings) ParseXml(str string) map[string]interface{} {
	m, err := mxj.NewMapXml([]byte(str))
	if err != nil {
		pluginConsole(sender.UUID).Error("xml解析错误：", err)
	}
	return m
}

func (sender *Strings) Xml(m map[string]interface{}) string {
	xmlStr, err := mxj.Map(m).Xml()
	if err != nil {
		pluginConsole(sender.UUID).Error("xml编码错误：", err)
	}
	return string(xmlStr)
}
