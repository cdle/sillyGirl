package jd_help

import (
	"regexp"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/cdle/sillyGirl/develop/qinglong"
)

func init() {
	crons, _ := qinglong.GetCrons("")
	for _, cron := range crons {
		if strings.Contains(cron.Command, "jd_get_share_code.js") && cron.IsDisabled == 0 {
			data, err := qinglong.GetCronLog(cron.ID)
			if err != nil {
				logs.Warn("助力码日志获取失败：%v", err)
				return
			}
			if data == "" {
				logs.Warn("助力码日志为空：%v", err)
				return
			}
			var codes = map[string][]string{
				"Fruit":        {},
				"Pet":          {},
				"Bean":         {},
				"JdFactory":    {},
				"DreamFactory": {},
				"Sgmh":         {},
				"Cash":         {},
			}
			for _, v := range regexp.MustCompile(`京东账号\d*（(.*)）(.*)】(\S*)`).FindAllStringSubmatch(data, -1) {
				if !strings.Contains(v[3], "种子") && !strings.Contains(v[3], "undefined") {
					// pt_pin := url.QueryEscape(v[1])
					for key, ss := range map[string][]string{
						"Fruit":        {"京东农场", "东东农场"},
						"Pet":          {"京东萌宠"},
						"Bean":         {"种豆得豆"},
						"JdFactory":    {"东东工厂"},
						"DreamFactory": {"京喜工厂"},
						"Jdzz":         {"京东赚赚"},
						"Sgmh":         {"闪购盲盒"},
						"Cash":         {"签到领现金"},
					} {
						for _, s := range ss {
							if strings.Contains(v[2], s) && v[3] != "" {
								codes[key] = append(codes[key], v[3])
							}
						}
					}
				}
			}
			var e = map[string]string{
				"Fruit":        "",
				"Pet":          "",
				"Bean":         "",
				"JdFactory":    "",
				"DreamFactory": "",
				"Sgmh":         "",
				"Cfd":          "",
				"Cash":         "",
			}
			for k := range codes {
				vv := codes[k]
				for i := range vv {
					vv[i] = strings.Replace(vv[i], `"`, `\"`, -1)

				}
				e[k] += strings.Join(vv, "@")
			}
			for k := range e {
				n := []string{}
				for i := 0; i < 20; i++ {
					n = append(n, e[k])
				}
				e[k] = strings.Join(n, "&")
			}
			var f = map[string]string{}
			for k := range e {
				switch k {
				case "Fruit":
					f["FRUITSHARECODES"] = e[k]
				case "Pet":
					f["PETSHARECODES"] = e[k]
				case "Bean":
					f["PLANT_BEAN_SHARECODES"] = e[k]
				case "JdFactory":
					f["DDFACTORY_SHARECODES"] = e[k]
				case "DreamFactory":
					f["DREAM_FACTORY_SHARE_CODES"] = e[k]
				case "Sgmh":
					f["JDSGMH_SHARECODES"] = e[k]
				case "Cash":
					f["JD_CASH_SHARECODES"] = e[k]
				}
			}
			envs := []qinglong.Env{}
			for i := range f {
				envs = append(envs, qinglong.Env{
					Name:  i,
					Value: f[i],
				})
			}
			qinglong.SetConfigEnv(envs...)
			return
		}
	}
}
