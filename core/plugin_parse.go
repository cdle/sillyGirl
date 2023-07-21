package core

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cdle/sillyGirl/core/common"
	"github.com/cdle/sillyGirl/core/storage"
	"github.com/cdle/sillyGirl/utils"
)

func pluginParse(script string, uuid string) (*common.Function, []func()) {
	var cbs = []func(){}
	var rules []string
	var imType *common.Filter
	var userId *common.Filter
	var groupId *common.Filter
	var cron = map[string]string{}
	var admin bool
	var disable bool
	var priority int
	var title string
	var public bool
	var description string
	var icon string
	var version string = "v1.0.0"
	var author string
	var create_at string
	var module bool
	var encrypt bool
	var onStart bool
	var origin = "自定义"
	var https = []*common.Http{}
	var message *common.Reply
	var FindAll bool
	var hasForm bool
	var carry bool
	var classes = []string{}
	ks := map[string]bool{}
	ress := regexp.MustCompile(
		`\*\s?@([\d\w+-]+)\s+([^\n]+?)\n`,
	).FindAllStringSubmatch(script, -1)
	for _, res := range ress {
		switch res[1] {
		case "rule", "match", "regex", "pattern":
			rule := strings.TrimSpace(res[2])
			rule = parseReply3(rule, func(s1, s2 string) {
				k := s1 + "." + s2
				if _, ok := ks[k]; !ok {
					cbs = append(cbs, func() {
						storage.Watch(MakeBucket(s1), s2, func(old, new, key string) *storage.Final {
							return &storage.Final{
								EndFunc: func() {
									plugins.Set(uuid, "reload")
								},
							}
						}, uuid)
					})
					ks[k] = true
				}
			})
			_rs := []string{}
		FR:
			ress := regexp.MustCompile(`\[([^\s\[\]]+)\]`).FindAllStringSubmatch(rule, -1)
			if len(ress) != 0 {
				res := ress[len(ress)-1]
				var inner = res[1]
				slice := strings.SplitN(inner, ":", 2)
				name := slice[0]
				ps := ""
				if len(slice) == 2 {
					ps = slice[1]
				}
				if strings.HasSuffix(name, "?") {
					name = strings.TrimRight(name, "?")
					rep := ""
					if ps == "" {
						rep = fmt.Sprintf("[%s]", name)
					} else {
						rep = fmt.Sprintf("[%s:%s]", name, ps)
					}
					for l := range _rs {
						_rs[l] = strings.Replace(_rs[l], res[0], rep, 1)
					}
					rule1 := strings.Replace(rule, res[0], rep, 1)
					if len(_rs) == 0 {
						_rs = append(_rs, rule1)
					}
					rule = strings.Replace(rule, res[0], "", 1)
					rule = regexp.MustCompile("\x20{2,}").ReplaceAllString(rule, " ")
					rule = strings.TrimSpace(rule)
					_rs = append(_rs, rule)
					goto FR
				}
			}
			if len(_rs) != 0 {
				rules = append(rules, _rs...)
			} else {
				rules = append(rules, rule)
			}
		case "class":
			classes = append(classes, regexp.MustCompile(`[\S]+`).FindAllString(res[2], -1)...)
			classes = utils.Unique(classes)
		case "platform", "imType", "platform+", "imType+":
			var item []string
			for _, i := range regexp.MustCompile(`[\d\w-]+`).FindAllString(res[2], -1) {
				item = append(item, strings.TrimSpace(i))
			}
			imType = &common.Filter{
				BlackMode: false,
				Items:     item,
			}
		case "platform-", "imType-":
			var item []string
			for _, i := range regexp.MustCompile(`[\d\w-]+`).FindAllString(res[2], -1) {
				item = append(item, strings.TrimSpace(i))
			}
			imType = &common.Filter{
				BlackMode: true,
				Items:     item,
			}
		case "userId", "userID", "uid", "userId+", "userID+", "uid+":
			var item []string
			for _, i := range regexp.MustCompile(`[\d\w-]+`).FindAllString(res[2], -1) {
				item = append(item, strings.TrimSpace(i))
			}
			userId = &common.Filter{
				BlackMode: false,
				Items:     item,
			}
		case "userId-", "userID-", "uid-":
			var item []string
			for _, i := range regexp.MustCompile(`[\d\w-]+`).FindAllString(res[2], -1) {
				item = append(item, strings.TrimSpace(i))
			}
			userId = &common.Filter{
				BlackMode: true,
				Items:     item,
			}
		case "groupId", "groupID", "groupCode", "chat_id", "chat_id+", "chatId", "chatID", "gid", "groupId+", "groupID+", "groupCode+", "chatId+", "chatID+", "gid+":
			var item []string
			for _, i := range regexp.MustCompile(`[\d\w-]+`).FindAllString(res[2], -1) {
				item = append(item, strings.TrimSpace(i))
			}
			groupId = &common.Filter{
				BlackMode: false,
				Items:     item,
			}
		case "groupId-", "groupID-", "groupCode-", "chatId-", "chat_id-", "chatID-", "gid-":
			var item []string
			for _, i := range regexp.MustCompile(`[\d\w-]+`).FindAllString(res[2], -1) {
				item = append(item, strings.TrimSpace(i))
			}
			groupId = &common.Filter{
				BlackMode: true,
				Items:     item,
			}

		case "admin":
			admin = strings.TrimSpace(res[2]) == "true"
		case "disable":
			disable = strings.TrimSpace(res[2]) == "true"
		case "findall":
			FindAll = strings.TrimSpace(res[2]) == "true"
		case "priority":
			priority = utils.Int(strings.TrimSpace(res[2]))
		case "title", "name", "show":
			title = strings.TrimSpace(res[2])
		case "public":
			public = strings.TrimSpace(res[2]) == "true"
		case "description":
			description = strings.TrimSpace(res[2])
		case "icon":
			icon = strings.TrimSpace(res[2])
		case "version":
			version = strings.TrimSpace(res[2])
		case "author":
			author = strings.TrimSpace(res[2])
		case "http":
			ss := regexp.MustCompile(`[\S]+`).FindAllString(strings.TrimSpace(res[2]), -1)
			if len(ss) == 2 {
				https = append(https, &common.Http{
					Path:   ss[1],
					Method: strings.ToUpper(ss[0]),
				})
			} else {
				console.Warn("http param is not 2")
			}
		case "message":
			ss := regexp.MustCompile(`[\S]+`).FindAllString(strings.TrimSpace(res[2]), -1)
			if len(ss) > 1 {
				if len(ss) == 2 && ss[1] == "*" {
					message = &common.Reply{
						Platform: ss[0],
						BotsID:   []string{},
					}
				} else {
					message = &common.Reply{
						Platform: ss[0],
						BotsID:   ss[1:],
					}
				}

			} else {
				console.Warn("message param is 0")
			}
		case "create_at":
			create_at = strings.TrimSpace(res[2])
		case "origin":
			origin = strings.TrimSpace(res[2])
		case "module":
			module = strings.TrimSpace(res[2]) == "true"
		case "carry":
			carry = strings.TrimSpace(res[2]) == "true"
		case "encrypt":
			encrypt = strings.TrimSpace(res[2]) == "true"
		case "on_start":
			onStart = strings.TrimSpace(res[2]) == "true"
		case "form":
			hasForm = true
		case "paterner":
			// paterner := strings.TrimSpace(res[2])
			// go func() {
			// 	time.Sleep(time.Second * 2)
			// 	getPaterner(uuid, strings.TrimSpace(paterner))
			// }()
		default:
			cron_ := strings.TrimSpace(res[2])
			cron_ = strings.ReplaceAll(cron_, `\/`, "/")
			if strings.HasPrefix(res[1], "cron") {
				cron[res[1]] = cron_
			}
		}
	}
	if !hasForm {
		hasForm = strings.Contains(script, "Form(")
	}
	return &common.Function{
		Rules:       rules,
		ImType:      imType,
		UserId:      userId,
		GroupId:     groupId,
		Cron:        cron,
		Admin:       admin,
		Priority:    priority,
		Disable:     disable,
		UUID:        uuid,
		Title:       title,
		Public:      public,
		Description: description,
		Icon:        icon,
		Version:     version,
		Author:      author,
		CreateAt:    create_at,
		Module:      module,
		Encrypt:     encrypt,
		OnStart:     onStart,
		Origin:      origin,
		Running:     onStart,
		Reply:       message,
		Https:       https,
		FindAll:     FindAll,
		HasForm:     hasForm,
		Carry:       carry,
		Classes:     classes,
	}, cbs
}
