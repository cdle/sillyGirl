package lottery

import (
	"fmt"
	"time"

	"github.com/cdle/sillyGirl/core"
)

var lottery = core.NewBucket("lottery")

type People struct {
	ID        int
	CreatedAt time.Time //参与时间
	Username  string    //用户名
	UserID    string    //用户ID
}

type Lottery struct {
	ID         int
	Name       string    //抽奖名称
	ImType     string    //群组类型
	GroupCode  string    //群组编号
	Prizes     []string  //奖品列表
	OpenMethod string    //开奖方式
	OpenTime   time.Time //开奖时间
	OpenNumber int       //开奖人数
	UserNumber int       //参数人数
	LuckyDogs  []string  //中奖者
	Keyword    string    //关键词
}

var cancel = "/cancel"
var e按时间自动开奖 = "按时间自动开奖"
var e按人数自动开奖 = "按人数自动开奖"
var e确定 = "确定"
var e取消 = "取消"

var help = `
· 如果你是活动参与者，中奖后需到我这里领取奖品。
· 如果你是群组管理员，邀请我进你所管理的群，身份为管理员，并给与删除消息和置顶消息的权限。通过 /create 命令可以在群里发起抽奖活动。
你也可以使用下面的命令来控制我：x w

参与者
/list - 已参与的活动
/wait - 待开奖的活动
/winlist - 领取奖品

发起者
/create - 在群组中使用此命令来创建一个抽奖活动
/released - 查询已发布的抽奖活动
/edit - 命令后面加上活动 ID 可以修改已发布的活动 ( ID 通过 /released 命令查询)
/close - 命令后面加上活动 ID 可以关闭正在进行中的活动 ( ID 通过 /released 命令查询)

其他命令
/cancel - 取消当前会话 ( 例如：取消当前正在创建的抽奖活动 )
`

func init() {
	core.AddCommand("", []core.Function{
		{
			Rules: []string{`raw [\s\S]+`},
			Handle: func(s core.Sender) interface{} {
				s.Continue()
				return nil
			},
		},
		{
			Rules: []string{`raw ^抽奖$`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				s.Reply(help)
				var stop = false
				var c = func(s string) bool {
					if s == cancel {
						stop = true
					}
					return stop
				}
				for {
					if stop == true {
						break
					}
					s.Await(s, func(s core.Sender) interface{} {
						switch s.GetContent() {
						//参与者
						case "/list":
							s.Reply("暂未实现。")
						case "/wait":
							s.Reply("暂未实现。")
						case "/winlist":
							s.Reply("暂未实现。")
						//发起者
						case "/create":
							Create(s, c)
						case "/released":
							s.Reply("暂未实现。")
						case "/edit":
							s.Reply("暂未实现。")
						case "/close":
							s.Reply("暂未实现。")

						//其他命令
						case cancel:
							stop = true
						case "/help":
							s.Reply(help)
						default:
							s.Reply("不支持的指令。")
						}
						return nil
					}, `/[a-z]+`)
				}
				return "已退出设置"
			},
		},
	})
}

func Create(s core.Sender, c func(string) bool) {
	cancal := false
	l := &Lottery{}
	show := ""
	s.Reply("请设置奖品名称：")
	s.Await(s, func(s core.Sender) interface{} {
		rt := s.GetContent()
		if c(rt) {
			cancal = true
			return nil
		}
		l.Name = rt
		show += fmt.Sprintf("奖品名称：%s", l.Name)
		return show
	})
	if cancal {
		return
	}

	var prizeNumber = 0
	for {
		s.Reply("请设置奖品数量：")
		s.Await(s, func(s core.Sender) interface{} {
			rt := s.GetContent()
			if c(rt) {
				cancal = true
			}
			prizeNumber = core.Int(rt)
			return nil
		})

		if prizeNumber != 0 {
			show += fmt.Sprintf("奖品数量：%d\n", prizeNumber)
			s.Reply(show)
			break
		}
		if cancal {
			return
		}
	}

	s.Reply(`请设置奖品内容 ( 1. 可以直接填写 APP 兑换码、支付宝口令红包等奖品让机器人自动发奖；也可留下你的联系方式，让中奖者主动联系你领奖。2. 有多少奖品数就回复多少次。 )：`)
	for i := 0; i < prizeNumber; i++ {
		s.Await(s, func(s core.Sender) interface{} {
			rt := s.GetContent()
			if c(rt) {
				cancal = true
				return nil
			}
			l.Prizes = append(l.Prizes, rt)
			if i == prizeNumber-1 {
				show += "奖品列表：\n"
				for j := range l.Prizes {
					show += fmt.Sprintf("%d. %s\n", j+1, l.Prizes[j])
				}
				return show
			}
			return fmt.Sprintf("继续设置下一个奖品内容：")
		})
		if cancal {
			return
		}
	}
	var choose = 0
	var tip = "请选择开奖方式：\n1. 按时间自动开奖\n2. 按人数自动开奖"
	for {
		s.Reply(tip)
		s.Await(s, func(s core.Sender) interface{} {
			rt := s.GetContent()
			if c(rt) {
				cancal = true
				return nil
			}
			return nil
		})
		if cancal {
			return
		}
		if choose == 1 || choose == 2 {
			break
		}
	}
	if choose == 1 {
		show += fmt.Sprintf("%s：\n", e按时间自动开奖)
		s.Reply(show)
		var tip = "请设置开奖时间 ( 格式：年-月-日 时:分 ) ："
		var rt = ""
		for {
			s.Reply(tip)
			s.Await(s, func(s core.Sender) interface{} {
				rt = s.GetContent()
				if c(rt) {
					cancal = true
					return nil
				}
				return nil
			})
			if cancal {
				return
			}
			if openTime, err := time.ParseInLocation("2006-01-02 15:04", rt, time.Local); err == nil {
				l.OpenTime = openTime
				break
			}
		}
		show += l.OpenTime.Format("开奖时间：2006-01-02 15:04\n")
		s.Reply(show)
	}
	if choose == 2 {
		show += fmt.Sprintf("%s：\n", e按时间自动开奖)
		s.Reply(show)
		var tip = "请设置开奖人数 ："
		var rt = ""
		for {
			s.Reply(tip)
			s.Await(s, func(s core.Sender) interface{} {
				rt = s.GetContent()
				if c(rt) {
					cancal = true
					return nil
				}
				return nil
			})
			if cancal {
				return
			}
			if l.OpenNumber = core.Int(rt); l.OpenNumber != 0 {
				break
			}
		}
		show += fmt.Sprintf("开奖人数：%d", l.OpenNumber)
		s.Reply(show)
	}
	s.Reply("请设置参与关键词：")
	s.Await(s, func(s core.Sender) interface{} {
		rt := s.GetContent()
		if c(rt) {
			cancal = true
			return nil
		}
		l.Keyword = rt
		show += fmt.Sprintf("关键词：%s", l.Keyword)
		return show
	})
	if cancal {
		return
	}
	tip = "已全部设置完成，是否发布？(确定/取消)"
	var rt = ""
	for {
		s.Reply(tip)
		s.Await(s, func(s core.Sender) interface{} {
			rt = s.GetContent()
			if c(rt) {
				cancal = true
				return nil
			}
			return nil
		})
		if cancal {
			return
		}
		if rt == e确定 || rt == e取消 {
			break
		}
	}
	if rt == e取消 {
		s.Reply(fmt.Sprintf("已取消 %s 抽奖活动发布", l.Name))
	}
	if rt == e确定 {
		lottery.Create(l)
		s.Reply(fmt.Sprintf("%s 抽奖活动已发布\n参与关键词：%s", l.Name, l.Keyword))
	}
	c(cancel)
}
