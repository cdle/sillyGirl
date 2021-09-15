package jd_asset

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/buger/jsonparser"
	"github.com/cdle/sillyGirl/core"
	"github.com/cdle/sillyGirl/develop/qinglong"
	"github.com/cdle/sillyGirl/im"
)

type JdCookie struct {
	ID        int
	PtKey     string
	PtPin     string
	WsKey     string
	Note      string
	Nickname  string
	BeanNum   string
	UserLevel string
	LevelName string
}

var ua = `Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_2 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8H7 Safari/6533.18.5 UCBrowser/13.4.2.1122`

func init() {
	core.AddCommand("jd", []core.Function{
		{
			Rules: []string{`asset ?`},
			Admin: true,
			Handle: func(s im.Sender) interface{} {
				a := s.Get()
				envs, err := qinglong.GetEnvs("JD_COOKIE")
				if err != nil {
					return err
				}
				if len(envs) == 0 {
					return "青龙没有京东账号。"
				}
				cks := []JdCookie{}
				for _, env := range envs {
					pt_key := FetchJdCookieValue("pt_key", env.Value)
					pt_pin := FetchJdCookieValue("pt_pin", env.Value)
					if pt_key != "" && pt_pin != "" {
						cks = append(cks, JdCookie{
							PtKey: pt_key,
							PtPin: pt_pin,
							Note:  env.Remarks,
						})
					}
				}
				cks = LimitJdCookie(cks, a)
				if len(cks) == 0 {
					return "没有匹配的京东账号。"
				}
				for _, ck := range cks {
					go s.Reply(ck.QueryAsset())
				}
				return nil
			},
		},
		{
			Rules: []string{`today bean(?)`},
			Admin: true,
			Handle: func(s im.Sender) interface{} {
				a := s.Get()
				envs, err := qinglong.GetEnvs("JD_COOKIE")
				if err != nil {
					return err
				}
				if len(envs) == 0 {
					return "青龙没有京东账号。"
				}
				cks := []JdCookie{}
				for _, env := range envs {
					pt_key := FetchJdCookieValue("pt_key", env.Value)
					pt_pin := FetchJdCookieValue("pt_pin", env.Value)
					if pt_key != "" && pt_pin != "" {
						cks = append(cks, JdCookie{
							PtKey: pt_key,
							PtPin: pt_pin,
							Note:  env.Remarks,
						})
					}
				}
				cks = LimitJdCookie(cks, a)
				if len(cks) == 0 {
					return "没有匹配的京东账号。"
				}
				var beans []chan int
				for _, ck := range cks {
					var bean = make(chan int)
					go GetTodayBean(&ck, bean)
					beans = append(beans, bean)
				}
				all := 0
				for i := range beans {
					all += <-beans[i]
				}
				return fmt.Sprintf("今日收入%d京豆。", all)
			},
		},
		{
			Rules: []string{`yestoday bean(?)`},
			Admin: true,
			Handle: func(s im.Sender) interface{} {
				a := s.Get()
				envs, err := qinglong.GetEnvs("JD_COOKIE")
				if err != nil {
					return err
				}
				if len(envs) == 0 {
					return "青龙没有京东账号。"
				}
				cks := []JdCookie{}
				for _, env := range envs {
					pt_key := FetchJdCookieValue("pt_key", env.Value)
					pt_pin := FetchJdCookieValue("pt_pin", env.Value)
					if pt_key != "" && pt_pin != "" {
						cks = append(cks, JdCookie{
							PtKey: pt_key,
							PtPin: pt_pin,
							Note:  env.Remarks,
						})
					}
				}
				cks = LimitJdCookie(cks, a)
				if len(cks) == 0 {
					return "没有匹配的京东账号。"
				}
				var beans []chan int
				for _, ck := range cks {
					var bean = make(chan int)
					go GetYestodayBean(&ck, bean)
					beans = append(beans, bean)
				}
				all := 0
				for i := range beans {
					all += <-beans[i]
				}
				return fmt.Sprintf("昨日收入%d京豆。", all)
			},
		},
		{
			Rules: []string{`bean(?)`},
			Admin: true,
			Handle: func(s im.Sender) interface{} {
				a := s.Get()
				envs, err := qinglong.GetEnvs("JD_COOKIE")
				if err != nil {
					return err
				}
				if len(envs) == 0 {
					return "青龙没有京东账号。"
				}
				cks := []JdCookie{}
				for _, env := range envs {
					pt_key := FetchJdCookieValue("pt_key", env.Value)
					pt_pin := FetchJdCookieValue("pt_pin", env.Value)
					if pt_key != "" && pt_pin != "" {
						cks = append(cks, JdCookie{
							PtKey: pt_key,
							PtPin: pt_pin,
							Note:  env.Remarks,
						})
					}
				}
				cks = LimitJdCookie(cks, a)
				if len(cks) == 0 {
					return "没有匹配的京东账号。"
				}
				all := 0
				for _, ck := range cks {
					ck.Available()
					all += Int(ck.BeanNum)
				}
				return fmt.Sprintf("总计%d京豆。", all)
			},
		},
	})

}

func LimitJdCookie(cks []JdCookie, a string) []JdCookie {
	ncks := []JdCookie{}
	if s := strings.Split(a, "-"); len(s) == 2 {
		for i := range cks {
			if i+1 >= Int(s[0]) && i+1 <= Int(s[1]) {
				ncks = append(ncks, cks[i])
			}
		}
	} else if x := regexp.MustCompile(`^[\s\d,]+$`).FindString(a); x != "" {
		xx := regexp.MustCompile(`(\d+)`).FindAllStringSubmatch(a, -1)
		for i := range cks {
			for _, x := range xx {
				if fmt.Sprint(i+1) == x[1] {
					ncks = append(ncks, cks[i])
				}
			}

		}
	} else if a != "" {
		a = strings.Replace(a, " ", "", -1)
		for i := range cks {
			if strings.Contains(cks[i].Note, a) || strings.Contains(cks[i].Nickname, a) || strings.Contains(cks[i].PtPin, a) {
				ncks = append(ncks, cks[i])
			}
		}
	}
	return ncks
}

type Asset struct {
	Nickname string
	Bean     struct {
		Total       int
		TodayIn     int
		TodayOut    int
		YestodayIn  int
		YestodayOut int
		ToExpire    []int
	}
	RedPacket struct {
		Total      float64
		ToExpire   float64
		ToExpireJd float64
		ToExpireJx float64
		ToExpireJs float64
		ToExpireJk float64
		Jd         float64
		Jx         float64
		Js         float64
		Jk         float64
	}
	Other struct {
		JsCoin   float64
		NcStatus float64
		McStatus float64
	}
}

var Int = func(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

var Float64 = func(s string) float64 {
	i, _ := strconv.ParseFloat(s, 64)
	return i
}

func (ck *JdCookie) QueryAsset() string {
	msgs := []string{}
	if ck.Note != "" {
		msgs = append(msgs, fmt.Sprintf("账号备注：%s", ck.Note))
	}
	asset := Asset{}
	if ck.Available() {
		msgs = append(msgs, fmt.Sprintf("用户等级：%v", ck.UserLevel))
		msgs = append(msgs, fmt.Sprintf("等级名称：%v", ck.LevelName))
		cookie := fmt.Sprintf("pt_key=%s;pt_pin=%s;", ck.PtKey, ck.PtPin)
		var rpc = make(chan []RedList)
		var fruit = make(chan string)
		var pet = make(chan string)
		var gold = make(chan int64)
		var egg = make(chan int64)
		var tyt = make(chan string)
		var mmc = make(chan int64)
		var zjb = make(chan int64)
		go redPacket(cookie, rpc)
		go initFarm(cookie, fruit)
		go initPetTown(cookie, pet)
		go jsGold(cookie, gold)
		go jxncEgg(cookie, egg)
		go tytCoupon(cookie, tyt)
		go mmCoin(cookie, mmc)
		go jdzz(cookie, zjb)
		today := time.Now().Local().Format("2006-01-02")
		yestoday := time.Now().Local().Add(-time.Hour * 24).Format("2006-01-02")
		page := 1
		end := false
		for {
			if end {
				msgs = append(msgs, []string{
					fmt.Sprintf("昨日收入：%d京豆", asset.Bean.YestodayIn),
					fmt.Sprintf("昨日支出：%d京豆", asset.Bean.YestodayOut),
					fmt.Sprintf("今日收入：%d京豆", asset.Bean.TodayIn),
					fmt.Sprintf("今日支出：%d京豆", asset.Bean.TodayOut),
				}...)
				break
			}
			bds := getJingBeanBalanceDetail(page, cookie)
			if bds == nil {
				end = true
				msgs = append(msgs, "京豆数据异常")
				break
			}
			for _, bd := range bds {
				amount := Int(bd.Amount)
				if strings.Contains(bd.Date, today) {
					if amount > 0 {
						asset.Bean.TodayIn += amount
					} else {
						asset.Bean.TodayOut += -amount
					}
				} else if strings.Contains(bd.Date, yestoday) {
					if amount > 0 {
						asset.Bean.YestodayIn += amount
					} else {
						asset.Bean.YestodayOut += -amount
					}
				} else {
					end = true
					break
				}
			}
			page++
		}
		msgs = append(msgs, fmt.Sprintf("当前京豆：%v京豆", ck.BeanNum))
		ysd := int(time.Now().Add(24 * time.Hour).Unix())
		if rps := <-rpc; len(rps) != 0 {
			for _, rp := range rps {
				b := Float64(rp.Balance)
				asset.RedPacket.Total += b
				if strings.Contains(rp.ActivityName, "京喜") || strings.Contains(rp.OrgLimitStr, "京喜") {
					asset.RedPacket.Jx += b
					if ysd >= rp.EndTime {
						asset.RedPacket.ToExpireJx += b
						asset.RedPacket.ToExpire += b
					}
				} else if strings.Contains(rp.ActivityName, "极速版") {
					asset.RedPacket.Js += b
					if ysd >= rp.EndTime {
						asset.RedPacket.ToExpireJs += b
						asset.RedPacket.ToExpire += b
					}

				} else if strings.Contains(rp.ActivityName, "京东健康") {
					asset.RedPacket.Jk += b
					if ysd >= rp.EndTime {
						asset.RedPacket.ToExpireJk += b
						asset.RedPacket.ToExpire += b
					}
				} else {
					asset.RedPacket.Jd += b
					if ysd >= rp.EndTime {
						asset.RedPacket.ToExpireJd += b
						asset.RedPacket.ToExpire += b
					}
				}
			}
			e := func(m float64) string {
				if m > 0 {
					return fmt.Sprintf(`(今日过期%.2f)`, m)
				}
				return ""
			}
			msgs = append(msgs, []string{
				fmt.Sprintf("所有红包：%.2f%s元🧧", asset.RedPacket.Total, e(asset.RedPacket.ToExpire)),
				fmt.Sprintf("京喜红包：%.2f%s元", asset.RedPacket.Jx, e(asset.RedPacket.ToExpireJx)),
				fmt.Sprintf("极速红包：%.2f%s元", asset.RedPacket.Js, e(asset.RedPacket.ToExpireJs)),
				fmt.Sprintf("健康红包：%.2f%s元", asset.RedPacket.Jk, e(asset.RedPacket.ToExpireJk)),
				fmt.Sprintf("京东红包：%.2f%s元", asset.RedPacket.Jd, e(asset.RedPacket.ToExpireJd)),
			}...)
		} else {
			msgs = append(msgs, "暂无红包数据🧧")
		}
		msgs = append(msgs, fmt.Sprintf("东东农场：%s", <-fruit))
		msgs = append(msgs, fmt.Sprintf("东东萌宠：%s", <-pet))
		gn := <-gold
		msgs = append(msgs, fmt.Sprintf("极速金币：%d(≈%.2f元)💰", gn, float64(gn)/10000))
		zjbn := <-zjb
		if zjbn != 0 {
			msgs = append(msgs, fmt.Sprintf("京东赚赚：%d金币(≈%.2f元)💰", zjbn, float64(zjbn)/10000))
		} else {
			msgs = append(msgs, fmt.Sprintf("京东赚赚：暂无数据"))
		}
		mmcCoin := <-mmc
		if mmcCoin != 0 {
			msgs = append(msgs, fmt.Sprintf("京东秒杀：%d秒秒币(≈%.2f元)💰", mmcCoin, float64(mmcCoin)/1000))
		} else {
			msgs = append(msgs, fmt.Sprintf("京东秒杀：暂无数据"))
		}
		msgs = append(msgs, fmt.Sprintf("推一推券：%s", <-tyt))
		msgs = append(msgs, fmt.Sprintf("惊喜牧场：%d枚鸡蛋🥚", <-egg))
		if ck.Note != "" {
			msgs = append([]string{
				fmt.Sprintf("账号备注：%s", ck.Note),
			}, msgs...)
		}
		msgs = append([]string{
			fmt.Sprintf("账号昵称：%s", ck.Nickname),
		}, msgs...)
	} else {
		ck.PtPin, _ = url.QueryUnescape(ck.PtPin)
		msgs = append(msgs, fmt.Sprintf("京东账号：%s", ck.PtPin))
		if ck.Note != "" {
			msgs = append(msgs, fmt.Sprintf("账号备注：%s", ck.Note))
		}
		msgs = append(msgs, []string{
			"提醒：该账号已过期，请重新登录",
		}...)
	}
	ck.PtPin, _ = url.QueryUnescape(ck.PtPin)
	return strings.Join(msgs, "\n")
}

type BeanDetail struct {
	Date         string `json:"date"`
	Amount       string `json:"amount"`
	EventMassage string `json:"eventMassage"`
}

func getJingBeanBalanceDetail(page int, cookie string) []BeanDetail {
	type AutoGenerated struct {
		Code       string       `json:"code"`
		DetailList []BeanDetail `json:"detailList"`
	}
	a := AutoGenerated{}
	req := httplib.Post(`https://api.m.jd.com/client.action?functionId=getJingBeanBalanceDetail`)
	req.Header("User-Agent", ua)
	req.Header("Host", "api.m.jd.com")
	req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.Header("Cookie", cookie)
	req.Body(fmt.Sprintf(`body={"pageSize": "20", "page": "%d"}&appid=ld`, page))
	data, err := req.Bytes()
	if err != nil {
		return nil
	}
	json.Unmarshal(data, &a)
	return a.DetailList
}

type RedList struct {
	ActivityName string `json:"activityName"`
	Balance      string `json:"balance"`
	BeginTime    int    `json:"beginTime"`
	DelayRemark  string `json:"delayRemark"`
	Discount     string `json:"discount"`
	EndTime      int    `json:"endTime"`
	HbID         string `json:"hbId"`
	HbState      int    `json:"hbState"`
	IsDelay      bool   `json:"isDelay"`
	OrgLimitStr  string `json:"orgLimitStr"`
}

func redPacket(cookie string, rpc chan []RedList) {
	type UseRedInfo struct {
		Count   int       `json:"count"`
		RedList []RedList `json:"redList"`
	}
	type Data struct {
		AvaiCount      int        `json:"avaiCount"`
		Balance        string     `json:"balance"`
		CountdownTime  string     `json:"countdownTime"`
		ExpiredBalance string     `json:"expiredBalance"`
		ServerCurrTime int        `json:"serverCurrTime"`
		UseRedInfo     UseRedInfo `json:"useRedInfo"`
	}
	type AutoGenerated struct {
		Data    Data   `json:"data"`
		Errcode int    `json:"errcode"`
		Msg     string `json:"msg"`
	}
	a := AutoGenerated{}
	req := httplib.Get(`https://m.jingxi.com/user/info/QueryUserRedEnvelopesV2?type=1&orgFlag=JD_PinGou_New&page=1&cashRedType=1&redBalanceFlag=1&channel=1&_=` + fmt.Sprint(time.Now().Unix()) + `&sceneval=2&g_login_type=1&g_ty=ls`)
	req.Header("User-Agent", ua)
	req.Header("Host", "m.jingxi.com")
	req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	req.Header("Accept-Language", "zh-cn")
	req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Referer", "https://st.jingxi.com/my/redpacket.shtml?newPg=App")
	req.Header("Cookie", cookie)
	data, _ := req.Bytes()
	json.Unmarshal(data, &a)
	rpc <- a.Data.UseRedInfo.RedList
}

func initFarm(cookie string, state chan string) {
	type RightUpResouces struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type TurntableInit struct {
		TimeState int `json:"timeState"`
	}
	type MengchongResouce struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type GUIDPopupTask struct {
		GUIDPopupTask string `json:"guidPopupTask"`
	}
	type IosConfigResouces struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type TodayGotWaterGoalTask struct {
		CanPop bool `json:"canPop"`
	}
	type LeftUpResouces struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type RightDownResouces struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type FarmUserPro struct {
		TotalEnergy     int    `json:"totalEnergy"`
		TreeState       int    `json:"treeState"`
		CreateTime      int64  `json:"createTime"`
		TreeEnergy      int    `json:"treeEnergy"`
		TreeTotalEnergy int    `json:"treeTotalEnergy"`
		ShareCode       string `json:"shareCode"`
		WinTimes        int    `json:"winTimes"`
		NickName        string `json:"nickName"`
		CouponKey       string `json:"couponKey"`
		CouponID        string `json:"couponId"`
		CouponEndTime   int64  `json:"couponEndTime"`
		Type            string `json:"type"`
		SimpleName      string `json:"simpleName"`
		Name            string `json:"name"`
		GoodsImage      string `json:"goodsImage"`
		SkuID           string `json:"skuId"`
		LastLoginDate   int64  `json:"lastLoginDate"`
		NewOldState     int    `json:"newOldState"`
		OldMarkComplete int    `json:"oldMarkComplete"`
		CommonState     int    `json:"commonState"`
		PrizeLevel      int    `json:"prizeLevel"`
	}
	type LeftDownResouces struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type LoadFriend struct {
		Code            string      `json:"code"`
		StatisticsTimes interface{} `json:"statisticsTimes"`
		SysTime         int64       `json:"sysTime"`
		Message         interface{} `json:"message"`
		FirstAddUser    bool        `json:"firstAddUser"`
	}
	type AutoGenerated struct {
		Code                  string                `json:"code"`
		RightUpResouces       RightUpResouces       `json:"rightUpResouces"`
		TurntableInit         TurntableInit         `json:"turntableInit"`
		IosShieldConfig       interface{}           `json:"iosShieldConfig"`
		MengchongResouce      MengchongResouce      `json:"mengchongResouce"`
		ClockInGotWater       bool                  `json:"clockInGotWater"`
		GUIDPopupTask         GUIDPopupTask         `json:"guidPopupTask"`
		ToFruitEnergy         int                   `json:"toFruitEnergy"`
		StatisticsTimes       interface{}           `json:"statisticsTimes"`
		SysTime               int64                 `json:"sysTime"`
		CanHongbaoContineUse  bool                  `json:"canHongbaoContineUse"`
		ToFlowTimes           int                   `json:"toFlowTimes"`
		IosConfigResouces     IosConfigResouces     `json:"iosConfigResouces"`
		TodayGotWaterGoalTask TodayGotWaterGoalTask `json:"todayGotWaterGoalTask"`
		LeftUpResouces        LeftUpResouces        `json:"leftUpResouces"`
		MinSupportAPPVersion  string                `json:"minSupportAPPVersion"`
		LowFreqStatus         int                   `json:"lowFreqStatus"`
		FunCollectionHasLimit bool                  `json:"funCollectionHasLimit"`
		Message               interface{}           `json:"message"`
		TreeState             int                   `json:"treeState"`
		RightDownResouces     RightDownResouces     `json:"rightDownResouces"`
		IconFirstPurchaseInit bool                  `json:"iconFirstPurchaseInit"`
		ToFlowEnergy          int                   `json:"toFlowEnergy"`
		FarmUserPro           FarmUserPro           `json:"farmUserPro"`
		RetainPopupLimit      int                   `json:"retainPopupLimit"`
		ToBeginEnergy         int                   `json:"toBeginEnergy"`
		LeftDownResouces      LeftDownResouces      `json:"leftDownResouces"`
		EnableSign            bool                  `json:"enableSign"`
		LoadFriend            LoadFriend            `json:"loadFriend"`
		HadCompleteXgTask     bool                  `json:"hadCompleteXgTask"`
		OldUserIntervalTimes  []int                 `json:"oldUserIntervalTimes"`
		ToFruitTimes          int                   `json:"toFruitTimes"`
		OldUserSendWater      []string              `json:"oldUserSendWater"`
	}
	a := AutoGenerated{}
	req := httplib.Post(`https://api.m.jd.com/client.action?functionId=initForFarm`)
	req.Header("accept", "*/*")
	req.Header("accept-encoding", "gzip, deflate, br")
	req.Header("accept-language", "zh-CN,zh;q=0.9")
	req.Header("cache-control", "no-cache")
	req.Header("cookie", cookie)
	req.Header("origin", "https://home.m.jd.com")
	req.Header("pragma", "no-cache")
	req.Header("referer", "https://home.m.jd.com/myJd/newhome.action")
	req.Header("sec-fetch-dest", "empty")
	req.Header("sec-fetch-mode", "cors")
	req.Header("sec-fetch-site", "same-site")
	req.Header("User-Agent", ua)
	req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.Body(`body={"version":4}&appid=wh5&clientVersion=9.1.0`)
	data, _ := req.Bytes()
	json.Unmarshal(data, &a)

	rt := a.FarmUserPro.Name
	if rt == "" {
		rt = "数据异常"
	} else {
		if a.TreeState == 2 || a.TreeState == 3 {
			rt += "已可领取⏰"
		} else if a.TreeState == 1 {
			rt += fmt.Sprintf("种植中，进度%.2f%%🍒", 100*float64(a.FarmUserPro.TreeEnergy)/float64(a.FarmUserPro.TreeTotalEnergy))
		} else if a.TreeState == 0 {
			rt = "您忘了种植新的水果⏰"
		}
	}
	state <- rt
}

func initPetTown(cookie string, state chan string) {
	type ResourceList struct {
		AdvertID string `json:"advertId"`
		ImageURL string `json:"imageUrl"`
		Link     string `json:"link"`
		ShopID   string `json:"shopId"`
	}
	type PetPlaceInfoList struct {
		Place  int `json:"place"`
		Energy int `json:"energy"`
	}
	type PetInfo struct {
		AdvertID     string `json:"advertId"`
		NickName     string `json:"nickName"`
		IconURL      string `json:"iconUrl"`
		ClickIconURL string `json:"clickIconUrl"`
		FeedGifURL   string `json:"feedGifUrl"`
		HomePetImage string `json:"homePetImage"`
		CrossBallURL string `json:"crossBallUrl"`
		RunURL       string `json:"runUrl"`
		TickleURL    string `json:"tickleUrl"`
	}
	type GoodsInfo struct {
		GoodsName        string `json:"goodsName"`
		GoodsURL         string `json:"goodsUrl"`
		GoodsID          string `json:"goodsId"`
		ExchangeMedalNum int    `json:"exchangeMedalNum"`
		ActivityID       string `json:"activityId"`
		ActivityIds      string `json:"activityIds"`
	}
	type Result struct {
		ShareCode              string             `json:"shareCode"`
		HisHbFlag              bool               `json:"hisHbFlag"`
		MasterHelpPeoples      []interface{}      `json:"masterHelpPeoples"`
		HelpSwitchOn           bool               `json:"helpSwitchOn"`
		UserStatus             int                `json:"userStatus"`
		TotalEnergy            int                `json:"totalEnergy"`
		MasterInvitePeoples    []interface{}      `json:"masterInvitePeoples"`
		ShareTo                string             `json:"shareTo"`
		PetSportStatus         int                `json:"petSportStatus"`
		UserImage              string             `json:"userImage"`
		MasterHelpReward       int                `json:"masterHelpReward"`
		ShowHongBaoExchangePop bool               `json:"showHongBaoExchangePop"`
		ShowNeedCollectPop     bool               `json:"showNeedCollectPop"`
		PetSportReward         string             `json:"petSportReward"`
		NewhandBubble          bool               `json:"newhandBubble"`
		ResourceList           []ResourceList     `json:"resourceList"`
		ProjectBubble          bool               `json:"projectBubble"`
		MasterInvitePop        bool               `json:"masterInvitePop"`
		MasterInviteReward     int                `json:"masterInviteReward"`
		MedalNum               int                `json:"medalNum"`
		MasterHelpPop          bool               `json:"masterHelpPop"`
		MeetDays               int                `json:"meetDays"`
		PetPlaceInfoList       []PetPlaceInfoList `json:"petPlaceInfoList"`
		MedalPercent           float64            `json:"medalPercent"`
		CharitableSwitchOn     bool               `json:"charitableSwitchOn"`
		PetInfo                PetInfo            `json:"petInfo"`
		NeedCollectEnergy      int                `json:"needCollectEnergy"`
		FoodAmount             int                `json:"foodAmount"`
		InviteCode             string             `json:"inviteCode"`
		RulesURL               string             `json:"rulesUrl"`
		PetStatus              int                `json:"petStatus"`
		GoodsInfo              GoodsInfo          `json:"goodsInfo"`
	}
	type AutoGenerated struct {
		Code       string `json:"code"`
		ResultCode string `json:"resultCode"`
		Message    string `json:"message"`
		Result     Result `json:"result"`
	}
	a := AutoGenerated{}
	req := httplib.Post(`https://api.m.jd.com/client.action?functionId=initPetTown`)
	req.Header("Host", "api.m.jd.com")
	req.Header("User-Agent", ua)
	req.Header("cookie", cookie)
	req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.Body(`body={}&appid=wh5&loginWQBiz=pet-town&clientVersion=9.0.4`)
	data, _ := req.Bytes()
	json.Unmarshal(data, &a)
	rt := ""
	if a.Code == "0" && a.ResultCode == "0" && a.Message == "success" {
		if a.Result.UserStatus == 0 {
			rt = "请手动开启活动⏰"
		} else if a.Result.GoodsInfo.GoodsName == "" {
			rt = "你忘了选购新的商品⏰"
		} else if a.Result.PetStatus == 5 {
			rt = a.Result.GoodsInfo.GoodsName + "已可领取⏰"
		} else if a.Result.PetStatus == 6 {
			rt = a.Result.GoodsInfo.GoodsName + "未继续领养新的物品⏰"
		} else {
			rt = a.Result.GoodsInfo.GoodsName + fmt.Sprintf("领养中，进度%.2f%%，勋章%d/%d🐶", a.Result.MedalPercent, a.Result.MedalNum, a.Result.GoodsInfo.ExchangeMedalNum)
		}
	} else {
		rt = "数据异常"
	}
	state <- rt
}

func jsGold(cookie string, state chan int64) { //

	type BalanceVO struct {
		CashBalance       string `json:"cashBalance"`
		EstimatedAmount   string `json:"estimatedAmount"`
		ExchangeGold      string `json:"exchangeGold"`
		FormatGoldBalance string `json:"formatGoldBalance"`
		GoldBalance       int    `json:"goldBalance"`
	}
	type Gears struct {
		Amount         string `json:"amount"`
		ExchangeAmount string `json:"exchangeAmount"`
		Order          int    `json:"order"`
		Status         int    `json:"status"`
		Type           int    `json:"type"`
	}
	type Data struct {
		Advertise      string    `json:"advertise"`
		BalanceVO      BalanceVO `json:"balanceVO"`
		Gears          []Gears   `json:"gears"`
		IsGetCoupon    bool      `json:"isGetCoupon"`
		IsGetCouponEid bool      `json:"isGetCouponEid"`
		IsLogin        bool      `json:"isLogin"`
		NewPeople      bool      `json:"newPeople"`
	}
	type AutoGenerated struct {
		Code      int    `json:"code"`
		Data      Data   `json:"data"`
		IsSuccess bool   `json:"isSuccess"`
		Message   string `json:"message"`
		RequestID string `json:"requestId"`
	}
	a := AutoGenerated{}
	req := httplib.Post(`https://api.m.jd.com?functionId=MyAssetsService.execute&appid=market-task-h5`)
	req.Header("Accept", "application/json, text/plain, */*")
	req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Cookie", cookie)
	req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.Header("Origin", "https://gold.jd.com")
	req.Header("Host", "api.m.jd.com")
	req.Header("Connection", "keep-alive")
	req.Header("User-Agent", ua)
	req.Header("Referer", "https://gold.jd.com/")
	req.Body(`functionId=MyAssetsService.execute&body={"method":"goldShopPage","data":{"channel":1}}&_t=` + fmt.Sprint(time.Now().Unix()) + `&appid=market-task-h5;`)
	data, _ := req.Bytes()
	json.Unmarshal(data, &a)
	state <- int64(a.Data.BalanceVO.GoldBalance)
}

func jxncEgg(cookie string, state chan int64) {
	req := httplib.Get("https://m.jingxi.com/jxmc/queryservice/GetHomePageInfo?channel=7&sceneid=1001&activeid=null&activekey=null&isgift=1&isquerypicksite=1&_stk=activeid%2Cactivekey%2Cchannel%2Cisgift%2Cisquerypicksite%2Csceneid&_ste=1&h5st=20210818211830955%3B4408816258824161%3B10028%3Btk01w8db21b2130ny2eg0siAPpNQgBqjGzYfuG6IP7Z%2BAOB40BiqLQ%2Blglfi540AB%2FaQrTduHbnk61ngEeKn813gFeRD%3Bd9a0b833bf99a29ed726cbffa07ba955cc27d1ff7d2d55552878fc18fc667929&_=1629292710957&sceneval=2&g_login_type=1&g_ty=ls")
	req.Header("User-Agent", ua)
	req.Header("Host", "m.jingxi.com")
	req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	req.Header("Accept-Language", "zh-cn")
	req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Referer", "https://st.jingxi.com/pingou/jxmc/index.html?nativeConfig=%7B%22immersion%22%3A1%2C%22toColor%22%3A%22%23e62e0f%22%7D&;__mcwvt=sjcp&ptag=7155.9.95")
	req.Header("Cookie", cookie)
	data, _ := req.Bytes()

	egg, _ := jsonparser.GetInt(data, "data", "eggcnt")
	state <- egg
}

func tytCoupon(cookie string, state chan string) {

	type DiscountInfo struct {
		High string        `json:"high"`
		Info []interface{} `json:"info"`
	}
	type ExtInfo struct {
		Num5              string `json:"5"`
		Num12             string `json:"12"`
		Num16             string `json:"16"`
		Num21             string `json:"21"`
		Num52             string `json:"52"`
		Num54             string `json:"54"`
		Num74             string `json:"74"`
		BusinessLabel     string `json:"business_label"`
		LimitOrganization string `json:"limit_organization"`
		UserLabel         string `json:"user_label"`
	}
	type Useable struct {
		AreaDesc         string        `json:"areaDesc"`
		AreaType         int           `json:"areaType"`
		Batchid          string        `json:"batchid"`
		BeanNumForPerson int           `json:"beanNumForPerson"`
		BeanNumForPlat   int           `json:"beanNumForPlat"`
		BeginTime        string        `json:"beginTime"`
		CanBeSell        bool          `json:"canBeSell"`
		CanBeShare       bool          `json:"canBeShare"`
		CompleteTime     string        `json:"completeTime"`
		CouponKind       int           `json:"couponKind"`
		CouponStyle      int           `json:"couponStyle"`
		CouponTitle      string        `json:"couponTitle"`
		Couponid         string        `json:"couponid"`
		Coupontype       int           `json:"coupontype"`
		CreateTime       string        `json:"createTime"`
		Discount         string        `json:"discount"`
		DiscountInfo     DiscountInfo  `json:"discountInfo"`
		EndTime          string        `json:"endTime"`
		ExpireType       int           `json:"expireType"`
		ExtInfo          ExtInfo       `json:"extInfo"`
		HourCoupon       int           `json:"hourCoupon"`
		IsOverlay        int           `json:"isOverlay"`
		LimitStr         string        `json:"limitStr"`
		LinkStr          string        `json:"linkStr"`
		OperateTime      string        `json:"operateTime"`
		OrderID          string        `json:"orderId"`
		OverlayDesc      string        `json:"overlayDesc"`
		PassKey          string        `json:"passKey"`
		Pin              string        `json:"pin"`
		PlatFormInfo     string        `json:"platFormInfo"`
		Platform         int           `json:"platform"`
		PlatformDetails  []interface{} `json:"platformDetails"`
		PwdKey           string        `json:"pwdKey"`
		Quota            string        `json:"quota"`
		SellID           string        `json:"sellId"`
		ShareID          string        `json:"shareId"`
		ShopID           string        `json:"shopId"`
		ShopName         string        `json:"shopName"`
		State            int           `json:"state"`
		UseTime          string        `json:"useTime"`
		VenderID         string        `json:"venderId"`
	}
	type Coupon struct {
		Curtimestamp           int       `json:"curtimestamp"`
		ExpiredCount           int       `json:"expired_count"`
		IsHideBaiTiaoInJxWxapp int       `json:"isHideBaiTiaoInJxWxapp"`
		IsHideMailInWxapp      int       `json:"isHideMailInWxapp"`
		Useable                []Useable `json:"useable"`
		UseableCount           int       `json:"useable_count"`
		UsedCount              int       `json:"used_count"`
	}
	type AutoGenerated struct {
		Coupon    Coupon `json:"coupon"`
		ErrMsg    string `json:"errMsg"`
		ErrorCode int    `json:"errorCode"`
		HasNext   int    `json:"hasNext"`
		Jdpin     string `json:"jdpin"`
		State     int    `json:"state"`
		Uin       string `json:"uin"`
	}
	a := AutoGenerated{}
	req := httplib.Get(`https://m.jingxi.com/activeapi/queryjdcouponlistwithfinance?state=1&wxadd=1&filterswitch=1&_=1629296270692&sceneval=2&g_login_type=1&callback=jsonpCBKB&g_ty=ls`)
	req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Cookie", cookie)
	req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.Header("Host", "m.jingxi.com")

	req.Header("User-Agent", ua)
	req.Header("Referer", "https://st.jingxi.com/my/coupon/jx.shtml?sceneval=2&ptag=7155.1.18")
	data, _ := req.Bytes()
	res := regexp.MustCompile(`jsonpCBKB[(](.*)\s+[)];}catch`).FindSubmatch(data)
	rt := "暂无数据"
	if len(res) > 0 {
		json.Unmarshal(res[1], &a)
		num := 0
		toexp := 0
		tm := int(time.Now().Unix() * 1000)
		for _, cp := range a.Coupon.Useable {
			if strings.Contains(cp.CouponTitle, "推推5.01") {
				num++
				if Int(cp.EndTime) < tm {
					toexp++
				}
			}
		}
		if num == 0 {
			rt = "无优惠券"
		} else {
			rt = fmt.Sprintf("%d张5元优惠券", num)
			if toexp > 0 {
				rt += fmt.Sprintf("(今天将过期%d张)⏰", toexp)
			} else {
				rt += "🎰"
			}
		}
	}
	state <- rt
}

func mmCoin(cookie string, state chan int64) {
	req := httplib.Post(`https://api.m.jd.com/client.action`)
	req.Header("Host", "api.m.jd.com")
	req.Header("Accept", "application/json, text/plain, */*")
	req.Header("Origin", "https://h5.m.jd.com")

	req.Header("User-Agent", ua)
	req.Header("cookie", cookie)
	req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.Body(`uuid=3245ad3d16ab2153c69f9ca91cd2e931b06a3bb8&clientVersion=10.1.0&client=wh5&osVersion=&area=&networkType=wifi&functionId=homePageV2&body=%7B%7D&appid=SecKill2020`)
	data, _ := req.Bytes()
	mmc, _ := jsonparser.GetInt(data, "result", "assignment", "assignmentPoints")
	state <- mmc
}

func jdzz(cookie string, state chan int64) { //
	req := httplib.Get(`https://api.m.jd.com/client.action?functionId=interactTaskIndex&body={}&client=wh5&clientVersion=9.1.0`)
	req.Header("Host", "api.m.jd.com")
	req.Header("Accept-Language", "zh-cn")
	req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Referer", "http://wq.jd.com/wxapp/pages/hd-interaction/index/index")
	req.Header("User-Agent", ua)
	req.Header("cookie", cookie)
	req.Header("Content-Type", "application/json")
	data, _ := req.Bytes()
	mmc, _ := jsonparser.GetString(data, "data", "totalNum")
	state <- int64(Int(mmc))
}

func (ck *JdCookie) Available() bool {
	cookie := "pt_key=" + ck.PtKey + ";pt_pin=" + ck.PtPin + ";"
	if ck == nil {
		return true
	}
	req := httplib.Get("https://me-api.jd.com/user_new/info/GetJDUserInfoUnion")
	req.Header("Cookie", cookie)
	req.Header("Accept", "*/*")
	req.Header("Accept-Language", "zh-cn,")
	req.Header("Connection", "keep-alive,")
	req.Header("Referer", "https://home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&")
	req.Header("Host", "me-api.jd.com")
	req.Header("User-Agent", "jdapp;iPhone;9.4.4;14.3;network/4g;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1")
	data, err := req.Bytes()
	if err != nil {
		return true
	}
	ui := &UserInfoResult{}
	if nil != json.Unmarshal(data, ui) {
		return true
	}
	switch ui.Retcode {
	case "1001": //ck.BeanNum
		if ui.Msg == "not login" {
			return false
		}
	case "0":
		if url.QueryEscape(ui.Data.UserInfo.BaseInfo.CurPin) != ck.PtPin {
			return av2(cookie)
		}
		if ui.Data.UserInfo.BaseInfo.Nickname != ck.Nickname || ui.Data.AssetInfo.BeanNum != ck.BeanNum || ui.Data.UserInfo.BaseInfo.UserLevel != ck.UserLevel || ui.Data.UserInfo.BaseInfo.LevelName != ck.LevelName {
			ck.UserLevel = ui.Data.UserInfo.BaseInfo.UserLevel
			ck.LevelName = ui.Data.UserInfo.BaseInfo.LevelName
			ck.Nickname = ui.Data.UserInfo.BaseInfo.Nickname
			ck.BeanNum = ui.Data.AssetInfo.BeanNum
		}
		return true
	}
	return av2(cookie)
}

func av2(cookie string) bool {
	req := httplib.Get(`https://m.jingxi.com/user/info/GetJDUserBaseInfo?_=1629334995401&sceneval=2&g_login_type=1&g_ty=ls`)
	req.Header("User-Agent", ua)
	req.Header("Host", "m.jingxi.com")
	req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	req.Header("Accept-Language", "zh-cn")
	req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Referer", "https://st.jingxi.com/my/userinfo.html?&ptag=7205.12.4")
	req.Header("Cookie", cookie)
	data, err := req.String()
	if err != nil {
		return true
	}
	return !strings.Contains(data, "login")
}

type UserInfoResult struct {
	Data struct {
		JdVvipCocoonInfo struct {
			JdVvipCocoon struct {
				DisplayType   int    `json:"displayType"`
				HitTypeList   []int  `json:"hitTypeList"`
				Link          string `json:"link"`
				Price         string `json:"price"`
				Qualification int    `json:"qualification"`
				SellingPoints string `json:"sellingPoints"`
			} `json:"JdVvipCocoon"`
			JdVvipCocoonStatus string `json:"JdVvipCocoonStatus"`
		} `json:"JdVvipCocoonInfo"`
		JdVvipInfo struct {
			JdVvipStatus string `json:"jdVvipStatus"`
		} `json:"JdVvipInfo"`
		AssetInfo struct {
			AccountBalance string `json:"accountBalance"`
			BaitiaoInfo    struct {
				AvailableLimit     string `json:"availableLimit"`
				BaiTiaoStatus      string `json:"baiTiaoStatus"`
				Bill               string `json:"bill"`
				BillOverStatus     string `json:"billOverStatus"`
				Outstanding7Amount string `json:"outstanding7Amount"`
				OverDueAmount      string `json:"overDueAmount"`
				OverDueCount       string `json:"overDueCount"`
				UnpaidForAll       string `json:"unpaidForAll"`
				UnpaidForMonth     string `json:"unpaidForMonth"`
			} `json:"baitiaoInfo"`
			BeanNum    string `json:"beanNum"`
			CouponNum  string `json:"couponNum"`
			CouponRed  string `json:"couponRed"`
			RedBalance string `json:"redBalance"`
		} `json:"assetInfo"`
		FavInfo struct {
			FavDpNum    string `json:"favDpNum"`
			FavGoodsNum string `json:"favGoodsNum"`
			FavShopNum  string `json:"favShopNum"`
			FootNum     string `json:"footNum"`
			IsGoodsRed  string `json:"isGoodsRed"`
			IsShopRed   string `json:"isShopRed"`
		} `json:"favInfo"`
		GrowHelperCoupon struct {
			AddDays     int     `json:"addDays"`
			BatchID     int     `json:"batchId"`
			CouponKind  int     `json:"couponKind"`
			CouponModel int     `json:"couponModel"`
			CouponStyle int     `json:"couponStyle"`
			CouponType  int     `json:"couponType"`
			Discount    float64 `json:"discount"`
			LimitType   int     `json:"limitType"`
			MsgType     int     `json:"msgType"`
			Quota       float64 `json:"quota"`
			RoleID      int     `json:"roleId"`
			State       int     `json:"state"`
			Status      int     `json:"status"`
		} `json:"growHelperCoupon"`
		KplInfo struct {
			KplInfoStatus string `json:"kplInfoStatus"`
			Mopenbp17     string `json:"mopenbp17"`
			Mopenbp22     string `json:"mopenbp22"`
		} `json:"kplInfo"`
		OrderInfo struct {
			CommentCount     string        `json:"commentCount"`
			Logistics        []interface{} `json:"logistics"`
			OrderCountStatus string        `json:"orderCountStatus"`
			ReceiveCount     string        `json:"receiveCount"`
			WaitPayCount     string        `json:"waitPayCount"`
		} `json:"orderInfo"`
		PlusPromotion struct {
			Status int `json:"status"`
		} `json:"plusPromotion"`
		UserInfo struct {
			BaseInfo struct {
				AccountType    string `json:"accountType"`
				BaseInfoStatus string `json:"baseInfoStatus"`
				CurPin         string `json:"curPin"`
				DefinePin      string `json:"definePin"`
				HeadImageURL   string `json:"headImageUrl"`
				LevelName      string `json:"levelName"`
				Nickname       string `json:"nickname"`
				Pinlist        string `json:"pinlist"`
				UserLevel      string `json:"userLevel"`
			} `json:"baseInfo"`
			IsHideNavi     string `json:"isHideNavi"`
			IsHomeWhite    string `json:"isHomeWhite"`
			IsJTH          string `json:"isJTH"`
			IsKaiPu        string `json:"isKaiPu"`
			IsPlusVip      string `json:"isPlusVip"`
			IsQQFans       string `json:"isQQFans"`
			IsRealNameAuth string `json:"isRealNameAuth"`
			IsWxFans       string `json:"isWxFans"`
			Jvalue         string `json:"jvalue"`
			OrderFlag      string `json:"orderFlag"`
			PlusInfo       struct {
			} `json:"plusInfo"`
			XbScore string `json:"xbScore"`
		} `json:"userInfo"`
		UserLifeCycle struct {
			IdentityID      string `json:"identityId"`
			LifeCycleStatus string `json:"lifeCycleStatus"`
			TrackID         string `json:"trackId"`
		} `json:"userLifeCycle"`
	} `json:"data"`
	Msg       string `json:"msg"`
	Retcode   string `json:"retcode"`
	Timestamp int64  `json:"timestamp"`
}

func FetchJdCookieValue(ps ...string) string {
	var key, cookies string
	if len(ps) == 2 {
		if len(ps[0]) > len(ps[1]) {
			key, cookies = ps[1], ps[0]
		} else {
			key, cookies = ps[0], ps[1]
		}
	}
	match := regexp.MustCompile(key + `=([^;]*);{0,1}`).FindStringSubmatch(cookies)
	if len(match) == 2 {
		return match[1]
	} else {
		return ""
	}
}

func GetTodayBean(ck *JdCookie, state chan int) {
	cookie := fmt.Sprintf("pt_key=%s;pt_pin=%s;", ck.PtKey, ck.PtPin)
	today := time.Now().Local().Format("2006-01-02")
	page := 1
	end := false
	in := 0
	defer func() {
		state <- in
	}()
	for {
		if end {
			return
		}
		bds := getJingBeanBalanceDetail(page, cookie)
		if bds == nil {
			break
		}
		for _, bd := range bds {
			amount := Int(bd.Amount)
			if strings.Contains(bd.Date, today) {
				if amount > 0 {
					in += amount
				} else {

				}
			} else {
				end = true
				break
			}
		}
		page++
	}
	return
}

func GetYestodayBean(ck *JdCookie, state chan int) {
	cookie := fmt.Sprintf("pt_key=%s;pt_pin=%s;", ck.PtKey, ck.PtPin)
	today := time.Now().Local().Format("2006-01-02")
	yestoday := time.Now().Local().Add(-time.Hour * 24).Format("2006-01-02")
	page := 1
	end := false
	in := 0
	defer func() {
		state <- in
	}()
	for {
		if end {
			return
		}
		bds := getJingBeanBalanceDetail(page, cookie)
		if bds == nil {
			break
		}
		for _, bd := range bds {
			amount := Int(bd.Amount)
			if strings.Contains(bd.Date, yestoday) {
				if amount > 0 {
					in += amount
				} else {

				}
			} else if strings.Contains(bd.Date, today) {

			} else {
				end = true
				break
			}
		}
		page++
	}
	return
}
