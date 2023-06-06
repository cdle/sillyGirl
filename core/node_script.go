package core

func Script(uuid string) map[string]interface{} {
	su := &ScriptUtils{
		script: plugins.GetString(uuid),
	}
	var o = map[string]interface{}{
		"get": su.GetValue,
		"save": func() {
			plugins.Set(uuid, su.script)
		},
	}
	o["set"] = func(key, value string) map[string]interface{} {
		su.SetValue(key, value)
		return o
	}
	return o
}
