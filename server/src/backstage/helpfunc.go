package backstage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lib"
	"log"
	//"net"
	//"net/http"
	"staticfunc"
	"time"
)

//! 一个玩家的结构
type HT_Player struct {
	Code         int     `json:"code"`
	Msg          string  `json:"msg"`
	Uid          int     `json:"uid"`       //! uid
	NickName     string  `json:"nickname"`  //! 昵称
	Head         string  `json:"head"`      //! 头像
	Level        int     `json:"level"`     //! 等级
	Deepin       int     `json:"deepin"`    //! 深度
	Card         int     `json:"card"`      //! 房卡数
	Gold         int     `json:"gold"`      //! 金币数
	Score        float64 `json:"score"`     //! 可兑换推广额
	T_Score      int     `json:"t_score"`   //! 已兑换推广额
	AgentCard    int     `json:"agentcard"` //! 代理房卡
	Rating       string  `json:"rating"`    //! 返利
	Top_Group    string  `json:"top_group"`
	PassWord     string  `json:"password"`
	Union_id     string  `json:"union_id"`
	ChildNum     int     `json:"childnum"` //! 下级人数
	Score1       float64 `json:"score1"`
	Score2       float64 `json:"score2"`
	CreateTime   string  `json:"createtime"`
	Parent       int     `json:"parent"`
	Name         string  `json:"name"`
	Alipay       string  `json:"alipay"`
	AliName      string  `json:"aliname"`
	Bankcard     string  `json:"bankcard"`
	Bankname     string  `json:"bankname"`
	Phone        string  `json:"phone"`
	AllCost      int64   `json:"allcost"`
	Commission   int     `json:"commission"`   //! 总佣金
	T_Commission int     `json:"t_commission"` //! 已提佣金
	Bills1       int64   `json:"bills1"`       //! 直属流水
	Bills2       int64   `json:"bills2"`       //! 下级流水
	AllBills     int64   `json:"allbills"`     //! 总流水
	DayBills     int64   `json:"daybills"`     //! 日流水
	WeekBills    int64   `json:"weekbills"`    //! 周流水
	MonthBills   int64   `json:"monthbills"`   //! 月流水
	Commission1  int     `json:"commission1"`  //! 直属佣金
	Commission2  int     `json:"commission2"`  //! 推广员佣金
	TodayAdd     int     `json:"todayadd"`     //! 今日新增
	MonthAdd     int     `json:"monthadd"`     //! 本月新增
	WeekAdd      int     `json:"weekadd"`      //! 本周新增
	AreaScale    int     `json:"areascale"`    //! 区域比例
	AreaScore    float64 `json:"areascore"`    //! 区域收益
	AreaTScore   int     `json:"areatscore"`   //! 已提现收益
	AgentMode    int     `json:"agentmode"`    //! 代理模式 0抽水模式 1流水模式
	MoneyMode    int     `json:"moneymode"`    //!
}

//! 一个玩家的基础结构
type HT_BasePlayer struct {
	Code       int     `json:"code"`
	Msg        string  `json:"msg"`
	Uid        int     `json:"uid"`       //! uid
	NickName   string  `json:"nickname"`  //! 昵称
	Head       string  `json:"head"`      //! 头像
	Level      int     `json:"level"`     //! 等级
	Deepin     int     `json:"deepin"`    //! 深度
	AgentCard  int     `json:"agentcard"` //! 代理房卡
	Score1     float64 `json:"score1"`
	Score2     float64 `json:"score2"`
	CreateTime string  `json:"createtime"`
}

//! 返回码
type HT_Code struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

//! 编码错误码
func HF_Code(code int, msg string) []byte {
	var node HT_Code
	node.Code = code
	node.Msg = msg
	log.Println(node)
	return lib.HF_JtoB(&node)
}

//! 查找一个玩家所有的子代理
func HF_FindSon(uid int, lst *[]int, deep int) {
	deep++
	res := GetServer().GetChildren(uid)
	for i := 0; i < len(res.Child); i++ {
		*lst = append(*lst, res.Child[i])
		HF_FindSon(res.Child[i], lst, deep)
	}
}

//! 查找一个玩家下面第N层代理
func HF_FindSonFromDeep(uid int, lst *[]int, deep int, destdeep int) {
	deep++
	res := GetServer().GetChildren(uid)
	for i := 0; i < len(res.Child); i++ {
		if deep == destdeep {
			*lst = append(*lst, res.Child[i])
		} else if deep < destdeep {
			HF_FindSonFromDeep(res.Child[i], lst, deep, destdeep)
		}
	}
}

//! 请求数据
func HF_GetPlayer(uid int, timeout time.Duration) (string, string, int, int, bool) {
	str := fmt.Sprintf("http://%s:%d/findplayer?uid=%d", GetServer().Con.GameIP, GetServer().Con.LoginPort, uid)
	res, err := lib.HF_Get(str, timeout)
	if err == nil {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err == nil && string(body) != "false" {
			var msg staticfunc.Msg_PersonInfo
			err := json.Unmarshal(body, &msg)
			if err == nil {
				return msg.Nickname, msg.Headurl, msg.Repertory[0].Repertory_count, msg.Repertory[1].Repertory_count, true
			}
		}
	}

	return "", "", 0, 0, false
}
