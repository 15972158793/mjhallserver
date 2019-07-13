package gameserver

import (
	"fmt"
	"lib"
	"math"
	"staticfunc"
	"time"
)

var SM_TIME = 70
var ROBOTBETS []int = []int{100, 500, 1000}

type Rec_SaiMa_Info struct {
	GameType int                    `json:"gametype"`
	Time     int64                  `json:"time"`
	Info     []Son_Rec_SaiMa_Person `json:"info"`
}
type Son_Rec_SaiMa_Person struct {
	Uid    int64   `json:"uid"`    //! uid
	Name   string  `json:"name"`   //! 名字
	Head   string  `json:"head"`   //! 头像
	Score  int     `json:"score"`  //! 这局的收益
	Result [2]int  `json:"result"` //! 前两名
	PL     [15]int `json:"pl"`     //! 15个区的赔率
	Bets   [15]int `json:"bets"`   //! 自己在15个区分别下了多少
}

type Game_GoldSaiMa struct {
	PersonMgr map[int64]*Game_GoldSaiMa_Person   `json:"personmgr"`
	Bets      [15]map[*Game_GoldSaiMa_Person]int `json:"bets"`
	RoBotBet  [15]int                            `json:"robotbet"`
	PL        [15]int                            `json:"pl"`     //! 15个区的赔率
	Result    [2]int                             `json:"result"` //! 结果
	Time      int64                              `json:"time"`
	Total     int                                `json:"total"`    //! 本局一共下了多少
	Money     int                                `json:"money"`    //! 系统庄的钱
	Trend     [][3]int                           `json:"trend"`    //! 走势
	Active    [3]int                             `json:"active"`   //! 活跃的三匹马
	Starter   int                                `json:"starter"`  //! 开始活跃的马
	Sprinter  int                                `json:"sprinter"` //! 中途活跃的马
	Finisher  int                                `json:"finisher"` //! 最后活跃的马
	Horse     [6]int                             `json:"horse"`    //! 排名

	room *Room
}

func NewGame_GoldSaiMa() *Game_GoldSaiMa {
	game := new(Game_GoldSaiMa)
	game.PersonMgr = make(map[int64]*Game_GoldSaiMa_Person)
	for i := 0; i < len(game.Bets); i++ {
		game.Bets[i] = make(map[*Game_GoldSaiMa_Person]int)
	}
	for i := 0; i < len(game.PL); i++ {
		game.PL[i] = -1
	}
	for i := 0; i < 5; i++ {
		tmp := []int{1, 2, 3, 4, 5, 6}
		trend := [3]int{0, 0, 0}
		for i := 0; i < 2; i++ {
			index := lib.HF_GetRandom(len(tmp))
			trend[i] = tmp[index]
			copy(tmp[index:], tmp[index+1:])
			tmp = tmp[:len(tmp)-1]
		}
		tmp = []int{3, 4, 5, 8, 30, 20, 50, 80, 100, 125, 175}
		trend[2] = tmp[lib.HF_GetRandom(len(tmp))]

		game.Trend = append(game.Trend, trend)
	}

	game.Time = 0

	arr := []int{1, 2, 3, 4, 5, 6}
	//! 随机三匹活跃马
	for i := 0; i < len(game.Active); i++ {
		index := lib.HF_GetRandom(len(arr))
		game.Active[i] = arr[index]
		copy(arr[index:], arr[index+1:])
		arr = arr[:len(arr)-1]
	}
	//! 设置活跃马的词缀
	game.Starter = game.Active[0]
	game.Sprinter = game.Active[1]
	game.Finisher = game.Active[2]

	//! 设置每个区的赔率
	var activePL []int //! 获取三个活跃马的赔率区下标
	for i := 0; i < len(game.Active)-1; i++ {
		for j := i + 1; j < len(game.Active); j++ {
			activePL = append(activePL, game.GetPLPos(game.Active[j]*10+game.Active[i]))

		}
	}
	arr = []int{3, 4, 5}
	for i := 0; i < len(activePL); i++ { //! 优先设置三个活跃区的赔率
		index := lib.HF_GetRandom(len(arr))
		game.PL[activePL[i]] = arr[index]
		copy(arr[index:], arr[index+1:])
		arr = arr[:len(arr)-1]
	}

	arr = []int{1000, 60, 10, 500, 8, 30, 125, 175, 80, 100, 20, 50}
	for i := 0; i < len(game.PL); i++ { //! 设置其他区的赔率
		if game.PL[i] != -1 {
			continue
		}
		index := lib.HF_GetRandom(len(arr))
		game.PL[i] = arr[index]
		copy(arr[index:], arr[index+1:])
		arr = arr[:len(arr)-1]
	}

	return game
}

type Game_GoldSaiMa_Person struct {
	Uid       int64   `json:"uid"`
	Gold      int     `json:"gold"`
	Total     int     `json:"total"`
	CurBets   int     `json:"curbets"`
	Win       int     `json:"win"`       //! 本局赢了多少
	Cost      int     `json:"cost"`      //! 抽水
	Bets      int     `json:"bets"`      //! 本局下了多少
	BetInfo   [15]int `json:"betinfo"`   //! 每个区的下注
	BeBets    int     `json:"bebets"`    //! 上局下注
	BeBetInfo [15]int `json:"bebetinfo"` //! 上局每个区的下注
	Round     int     `json:"round"`
	Name      string  `json:"name"`
	Head      string  `json:"head"`
	Online    bool    `json:"online"`
	IP        string  `json:"ip"`
	Address   string  `json:"address"`
	Sex       int     `json:"int"`
}

type Msg_GameGoldSaiMa_Total struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"`
}

type Msg_GameGoldSaiMa_Info struct {
	Begin     bool     `json:"begin"`
	Time      int64    `json:"time"`
	PL        [15]int  `json:"pl"`        //! 15个区的赔率
	Active    [3]int   `json:"active"`    //! 活跃的三匹马
	BetInfo   [15]int  `json:"betinfo"`   //! 自己15个区的下注
	GameTatol [15]int  `json:"gametotal"` //! 15个区下的总注
	Total     int      `json:"total"`     //! 自己的钱
	Trend     [][3]int `json:"trend"`     //! 走势
}

type Msg_GameGoldSaiMa_End struct {
	Uid    int64                      `json:"uid"`
	Total  int                        `json:"total"` //! 总金币
	Win    int                        `json:"win"`   //! 赢了多少
	WinPL  int                        `json:"winpl"` //! 获胜区的赔率
	Rating []Son_GameGoldSaiMa_Person `json:"rating"`
	PL     [15]int                    `json:"pl"`     //! 15个区的赔率
	Active [3]int                     `json:"active"` //! 活跃的三匹马
	Horse  [6]int                     `json:"horse"`  //! 排名
	Rand   int                        `json:"rand"`
}

type Son_GameGoldSaiMa_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Win   int    `json:"win"`   //! 赢了多少
	Total int    `json:"total"` //! 自己的钱
}

type Msg_GameGoldSaiMa_Bets struct {
	Uid       int64   `json:"uid"`
	Index     int     `json:"index"`     //! 投注的下标
	Gold      int     `json:"gold"`      //! 投注金额
	Total     int     `json:"total"`     //! 玩家金币总额
	GameTotal [15]int `json:"gametotal"` //! 游戏15个区域的下注
}
type Msg_GameGoldSaiMa_GoOn struct {
	Uid       int64   `json:"uid"`
	Gold      [15]int `json:"gold"`      //! 玩家每个区本次的下注
	Total     int     `json:"total"`     //! 玩家金币总额
	GameTotal [15]int `json:"gametotal"` //! 游戏15个区域的下注
}

//! 获取两匹马的组合在赔率数组中的下标
func (self *Game_GoldSaiMa) GetPLPos(index int) int {
	if index/10 > index%10 {
		index = index%10*10 + index/10
	}

	tmp := []int{16, 15, 14, 13, 12, 26, 25, 24, 23, 36, 35, 34, 46, 45, 56}
	for i := 0; i < len(tmp); i++ {
		if index == tmp[i] {
			return i
		}
	}
	return -1
}

//! 得到这个位置下了多少钱 true加上机器人下注 false不加机器人下注
func (self *Game_GoldSaiMa) GetMoneyPos(index int, all bool) int {
	total := 0
	for _, value := range self.Bets[index] {
		total += value
		if all {
			total += self.RoBotBet[index]
		}
	}
	return total
}

func (self *Game_GoldSaiMa) getinfo(uid int64) *Msg_GameGoldSaiMa_Info {
	var msg Msg_GameGoldSaiMa_Info
	msg.Begin = self.room.Begin
	if self.Time == 0 {
		msg.Time = 0
	} else {
		msg.Time = self.Time - time.Now().Unix()
	}

	msg.Active = self.Active
	msg.PL = self.PL
	for i := 0; i < len(self.Bets); i++ {
		msg.GameTatol[i] = self.GetMoneyPos(i, true)
	}
	person := self.GetPerson(uid)
	if person == nil {
		return nil
	}
	msg.BetInfo = person.BetInfo
	msg.Total = person.Total
	msg.Trend = self.Trend
	return &msg
}

//! 同步金币
func (self *Game_GoldSaiMa_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

func (self *Game_GoldSaiMa) GetPerson(uid int64) *Game_GoldSaiMa_Person {
	return self.PersonMgr[uid]
}

//! 同步总分
func (self *Game_GoldSaiMa) SendTotal(uid int64, total int) {
	var msg Msg_GameGoldSaiMa_Total
	msg.Uid = uid
	msg.Total = total
	self.room.SendMsg(uid, "gamegoldtotal", &msg)
}

//! 设置时间
func (self *Game_GoldSaiMa) SetTime(t int) {
	if t == 0 {
		self.Time = 0
	} else {
		self.Time = time.Now().Unix() + int64(t)
	}

	var msg Msg_SetTime
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	self.room.broadCastMsg("gametime", &msg)
}

//! 是否下注了
func (self *Game_GoldSaiMa) OnIsBets(uid int64) bool {
	value, ok := self.PersonMgr[uid]
	if ok {
		return value.Bets > 0
	}
	return false
}

//! 下注
func (self *Game_GoldSaiMa) GameBets(uid int64, index int, gold int) {
	if uid == 0 {
		self.RoBotBet[index] += gold
		var msg Msg_GameGoldSaiMa_Bets
		msg.Uid = uid
		msg.Index = index
		msg.Gold = gold
		msg.Total = 0
		for i := 0; i < 15; i++ {
			msg.GameTotal[i] = self.GetMoneyPos(i, true)
		}
		self.room.broadCastMsg("gamegoldsaimabets", &msg)
		return
	}

	if index < 0 || index > 14 {
		return
	}

	if gold <= 0 {
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() >= int64(SM_TIME-2) {
		self.room.SendErr(uid, "正在开奖，请稍后下注。")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if person.Total < lib.GetManyMgr().GetProperty(self.room.Type).MinBet {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("%d金币才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("%d金币才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("%d万金币才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/10000))
		}
		return
	}

	if person.Total < gold {
		self.room.SendErr(uid, "您的金币不足，请前往充值。")
		return
	}

	if person.Bets+gold > lib.GetManyMgr().GetProperty(self.room.Type).MaxBet {
		self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d。", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet))
		return
	}

	person.CurBets += gold
	person.Bets += gold
	person.Total -= gold
	person.BetInfo[index] += gold
	person.Round = 0
	self.Total += gold
	self.Bets[index][person] += gold

	var msg Msg_GameGoldSaiMa_Bets
	msg.Uid = uid
	msg.Index = index
	msg.Gold = gold
	msg.Total = person.Total
	for i := 0; i < 15; i++ {
		msg.GameTotal[i] = self.GetMoneyPos(i, true)
	}
	self.room.broadCastMsg("gamegoldsaimabets", &msg)

}

//! 续押
func (self *Game_GoldSaiMa) GameGoOn(uid int64) {
	if uid == 0 {
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() >= int64(SM_TIME-2) {
		self.room.SendErr(uid, "正在开奖,请稍后下注。")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if person.Total < lib.GetManyMgr().GetProperty(self.room.Type).MinBet {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("%d金币才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("%d金币才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("%d万金币才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/10000))
		}
		return
	}

	if person.Total < person.BeBets {
		self.room.SendErr(uid, "您的金币不足，请前往充值。")
		return
	}

	if person.Bets+person.BeBets > lib.GetManyMgr().GetProperty(self.room.Type).MaxBet {
		self.room.SendErr(uid, fmt.Sprintln("单局下注不能超过%d。", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet))
		return
	}

	person.CurBets += person.BeBets
	person.Bets += person.BeBets
	person.Total -= person.BeBets
	self.Total += person.BeBets
	for i := 0; i < len(person.BeBetInfo); i++ {
		person.BetInfo[i] += person.BeBetInfo[i]
		self.Bets[i][person] += person.BeBetInfo[i]
	}
	person.Round = 0

	var msg Msg_GameGoldSaiMa_GoOn
	msg.Uid = uid
	msg.Gold = person.BeBetInfo
	msg.Total = person.Total
	for i := 0; i < 15; i++ {
		msg.GameTotal[i] = self.GetMoneyPos(i, true)
	}
	self.room.broadCastMsg("gamegoldsaimagoon", &msg)

}

func (self *Game_GoldSaiMa) OnBegin() {
	self.room.Begin = true
	//! ------------------------------- 确定一二名
	winlst := make([]int, 0)
	lostlst := make([]int, 0)
	mustwinlst := make([]int, 0)  //!　一定会赢（当前奖池已经跌过最低线）
	mustlostlst := make([]int, 0) //! 一定会输（当前奖池已经超过最高线）
	ret := -1
	for i := 0; i < len(self.Active); i++ {
		for j := 1; j < 7; j++ {
			if j == self.Active[i] {
				continue
			}
			index := self.GetPLPos(self.Active[i]*10 + j)
			win := self.Total - self.GetMoneyPos(index, false)*self.PL[index]
			if GetServer().SaiMaMoney[self.room.Type%130000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && GetServer().SaiMaMoney[self.room.Type%130000]+int64(win) <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
				if win > 0 {
					if self.Active[i] == self.Finisher { //! 终结者多加一次
						winlst = append(winlst, self.Active[i]*10+j)
					}

					if j == self.Finisher || j == self.Starter || j == self.Sprinter { //！两匹马都是活跃马多加一次
						winlst = append(winlst, self.Active[i]*10+j)
					}

					winlst = append(winlst, self.Active[i]*10+j)
				} else {
					if self.Active[i] == self.Finisher { //! 终结者多加一次
						lostlst = append(lostlst, self.Active[i]*10+j)
					}

					if j == self.Finisher || j == self.Starter || j == self.Sprinter { //！两匹马都是活跃马多加一次
						lostlst = append(lostlst, self.Active[i]*10+j)
					}

					lostlst = append(lostlst, self.Active[i]*10+j)
				}
			}
			if win > 0 {
				if self.Active[i] == self.Finisher { //! 终结者多加一次
					mustwinlst = append(mustwinlst, self.Active[i]*10+j)
				}

				if j == self.Finisher || j == self.Starter || j == self.Sprinter { //！两匹马都是活跃马多加一次
					mustwinlst = append(mustwinlst, self.Active[i]*10+j)
				}

				mustwinlst = append(mustwinlst, self.Active[i]*10+j)
			} else {
				if self.Active[i] == self.Finisher { //! 终结者多加一次
					mustlostlst = append(mustlostlst, self.Active[i]*10+j)
				}

				if j == self.Finisher || j == self.Starter || j == self.Sprinter { //！两匹马都是活跃马多加一次
					mustlostlst = append(mustlostlst, self.Active[i]*10+j)
				}

				mustlostlst = append(mustlostlst, self.Active[i]*10+j)
			}
		}
	}
	willwin := true
	if GetServer().SaiMaMoney[self.room.Type%130000] < lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "--------会赢1")
		willwin = true
	} else if GetServer().SaiMaMoney[self.room.Type%130000] > lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "--------会输1")
		willwin = false
	} else {
		sum := lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax - lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin
		s1 := GetServer().SaiMaMoney[self.room.Type%130000] - lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin
		pro := 1000000 * s1 / int64(sum)
		if lib.HF_GetRandom(1000000) < int(pro) {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "--------会输2")
			willwin = false
		} else {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "--------会赢2")
			willwin = true
		}
	}

	//willwin = false

	if willwin {
		if GetServer().SaiMaMoney[self.room.Type%130000] < lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && len(mustwinlst) != 0 {
			ret = mustwinlst[lib.HF_GetRandom(len(mustwinlst))]
		} else if len(winlst) != 0 {
			ret = winlst[lib.HF_GetRandom(len(winlst))]
		}
	} else {
		if GetServer().SaiMaMoney[self.room.Type%130000] > lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax && len(mustlostlst) != 0 {
			ret = mustlostlst[lib.HF_GetRandom(len(mustlostlst))]
		} else if len(lostlst) != 0 {
			ret = lostlst[lib.HF_GetRandom(len(lostlst))]
		}
	}

	if ret == -1 { //! 没有合适的方案,随机开
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "没有合适的方案,随机开")
		index := self.Active[lib.HF_GetRandom(len(self.Active))]
		for true {
			tmp := lib.HF_GetRandom(6) + 1
			if tmp == index {
				continue
			}
			ret = index + tmp*10
			break
		}
	}

	if lib.HF_GetRandom(9) < 5 { //! 随机第一名和第二名
		self.Result[0] = ret % 10
		self.Result[1] = ret / 10
	} else {
		self.Result[0] = ret / 10
		self.Result[1] = ret % 10
	}
	self.Horse[0] = self.Result[0]
	self.Horse[1] = self.Result[1]
	//! 随机3到6名
	arr := []int{1, 2, 3, 4, 5, 6}
	for i := 0; i < len(self.Result); i++ {
		for j := 0; j < len(arr); j++ {
			if arr[j] == self.Result[i] {
				copy(arr[j:], arr[j+1:])
				arr = arr[:len(arr)-1]
				break
			}
		}
	}

	for i := 2; i < len(self.Horse); i++ {
		index := lib.HF_GetRandom(len(arr))
		self.Horse[i] = arr[index]
		copy(arr[index:], arr[index+1:])
		arr = arr[:len(arr)-1]
	}

	self.OnEnd()

}

func (self *Game_GoldSaiMa) OnEnd() {
	self.room.Begin = false
	//! -------------------------------- 结算
	pos := self.GetPLPos(self.Result[0]*10 + self.Result[1])
	posPL := self.PL[pos]

	result := [3]int{self.Result[0], self.Result[1], self.PL[pos]}
	trend := [][3]int{result}
	trend = append(trend, self.Trend...)
	if len(trend) > 5 {
		trend = trend[0:5]
	}
	self.Trend = trend

	dealwin := self.Total - self.GetMoneyPos(pos, false)*self.PL[pos]
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "------- dealwin : ", dealwin, " self.total : ", self.Total, " self.GetMoneyPos(pos)*self.PL[pos] : ", self.GetMoneyPos(pos, false)*self.PL[pos])
	if dealwin != 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "----- sqlbzwlog")
		GetServer().SqlBZWLog(&SQL_BZWLog{1, dealwin, time.Now().Unix(), self.room.Type})
	}
	if dealwin > 0 {
		cost := int(math.Ceil(float64(dealwin) * float64(lib.GetManyMgr().GetProperty(self.room.Type).DealCost) / 100.0))
		dealwin -= cost
	}
	if dealwin != 0 {
		GetServer().SetSaiMaMoney(self.room.Type%130000, GetServer().SaiMaMoney[self.room.Type%130000]+int64(dealwin))
	}

	for i := 0; i < len(self.Bets); i++ {
		if pos == i {
			for key, value := range self.Bets[i] {
				winmoney := value * self.PL[pos]
				key.Win += winmoney
				key.Cost = int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
			}
		}
	}

	var winner []Son_GameGoldSaiMa_Person
	for _, value := range self.PersonMgr {
		if value.Win > 0 {
			if value.Win > 0 && value.Cost != 0 {
				value.Win -= value.Cost
				GetServer().SqlAgentGoldLog(value.Uid, value.Cost, self.room.Type)
				GetServer().SqlAgentBillsLog(value.Uid, value.Cost/2, self.room.Type)
			}
			value.Total += value.Win

			winner = append(winner, Son_GameGoldSaiMa_Person{value.Uid, value.Name, value.Head, value.Win, value.Total})
		} else if value.Win-value.Bets < 0 {
			cost := int(math.Ceil(float64(value.Win-value.Bets) * float64(lib.GetManyMgr().GetProperty(self.room.Type).Cost) / 200.0))
			GetServer().SqlAgentBillsLog(value.Uid, cost, self.room.Type)
		}
		//! 插入战绩
		if value.Bets > 0 {
			var record Rec_SaiMa_Info
			record.Time = time.Now().Unix()
			record.GameType = self.room.Type
			var rec Son_Rec_SaiMa_Person
			rec.Uid = value.Uid
			rec.Name = value.Name
			rec.Head = value.Head
			rec.Score = value.Win - value.Bets
			rec.Bets = value.BetInfo
			rec.PL = self.PL
			rec.Result = self.Result
			record.Info = append(record.Info, rec)
			GetServer().InsertRecord(self.room.Type, value.Uid, lib.HF_JtoA(&record), rec.Score)
		}
	}

	//! 排行榜
	for j := 0; j < len(winner)-1; j++ {
		for i := 0; i < len(winner)-1-j; i++ {
			if winner[i].Win < winner[i+1].Win {
				tmp := winner[i]
				winner[i] = winner[i+1]
				winner[i+1] = tmp
			}
		}
	}

	//! -------------------------------- 初始化
	for i := 0; i < len(self.Active); i++ { //! 活跃马
		self.Active[i] = -1
	}
	for i := 0; i < len(self.PL); i++ { //! 赔率
		self.PL[i] = -1
	}
	for i := 0; i < len(self.Bets); i++ { //! 下注
		self.Bets[i] = make(map[*Game_GoldSaiMa_Person]int)
	}
	self.Total = 0
	self.Starter = -1
	self.Sprinter = -1
	self.Finisher = -1
	for i := 0; i < len(self.RoBotBet); i++ {
		self.RoBotBet[i] = 0
	}

	//! -------------------------------- 设置下一局的 赔率 ，活跃马
	arr := []int{1, 2, 3, 4, 5, 6}
	//! 随机三匹活跃马
	for i := 0; i < len(self.Active); i++ {
		index := lib.HF_GetRandom(len(arr))
		self.Active[i] = arr[index]
		copy(arr[index:], arr[index+1:])
		arr = arr[:len(arr)-1]
	}
	//! 设置活跃马的词缀
	self.Starter = self.Active[0]
	self.Sprinter = self.Active[1]
	self.Finisher = self.Active[2]

	//! 设置每个区的赔率
	var activePL []int //! 获取三个活跃马的赔率区下标
	for i := 0; i < len(self.Active)-1; i++ {
		for j := i + 1; j < len(self.Active); j++ {
			activePL = append(activePL, self.GetPLPos(self.Active[j]*10+self.Active[i]))

		}
	}
	arr = []int{3, 4, 5}
	for i := 0; i < len(activePL); i++ { //! 优先设置三个活跃区的赔率
		index := lib.HF_GetRandom(len(arr))
		self.PL[activePL[i]] = arr[index]
		copy(arr[index:], arr[index+1:])
		arr = arr[:len(arr)-1]
	}

	arr = []int{1000, 60, 10, 500, 8, 30, 125, 175, 80, 100, 20, 50}
	for i := 0; i < len(self.PL); i++ { //! 设置其他区的赔率
		if self.PL[i] != -1 {
			continue
		}
		index := lib.HF_GetRandom(len(arr))
		self.PL[i] = arr[index]
		copy(arr[index:], arr[index+1:])
		arr = arr[:len(arr)-1]
	}

	var rating []Son_GameGoldSaiMa_Person
	for i := 0; i < len(winner); i++ {
		var son Son_GameGoldSaiMa_Person
		son.Name = winner[i].Name
		son.Head = winner[i].Head
		son.Total = winner[i].Total
		son.Uid = winner[i].Uid
		son.Win = winner[i].Win
		rating = append(rating, son)
		if i >= 2 {
			break
		}
	}
	for _, value := range self.PersonMgr {
		var msg Msg_GameGoldSaiMa_End
		msg.Uid = value.Uid
		msg.Win = value.Win
		msg.WinPL = posPL
		msg.Total = value.Total
		msg.Rating = rating
		msg.PL = self.PL
		msg.Active = self.Active
		msg.Horse = self.Horse
		msg.Rand = lib.HF_GetRandom(10000)
		self.room.SendMsg(value.Uid, "gamegoldsaimaend", &msg)
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "---- win : ", value.Win, " cost : ", value.Cost)
	}

	for i := 0; i < len(self.Horse); i++ {
		self.Horse[i] = -1
	}

	//! 清理玩家
	for key, value := range self.PersonMgr {
		if value.Online {
			if value.Round >= 5 {
				self.room.KickViewByUid(value.Uid, 96)
			} else {
				value.BeBets = value.Bets
				value.BeBetInfo = value.BetInfo
				value.Win = 0
				value.Cost = 0
				value.Bets = 0
				for j := 0; j < len(value.BetInfo); j++ {
					value.BetInfo[j] = 0
				}
				continue
			}
		}
		delete(self.PersonMgr, key)
	}

	self.SetTime(SM_TIME)
}

func (self *Game_GoldSaiMa) OnInit(room *Room) {
	self.room = room
	self.Money = 100000000
}

func (self *Game_GoldSaiMa) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldSaiMa) OnSendInfo(person *Person) {
	value, ok := self.PersonMgr[person.Uid]
	if ok {
		value.Online = true
		value.IP = person.ip
		value.Address = person.minfo.Address
		value.Sex = person.Sex
		value.CurBets = 0
		value.SynchroGold(person.Gold)
		person.SendMsg("gamegoldsaimainfo", self.getinfo(person.Uid))
		return
	}
	if self.Time == 0 {
		self.SetTime(15)
	}

	_person := new(Game_GoldSaiMa_Person)
	_person.Uid = person.Uid
	_person.Gold = person.Gold
	_person.Total = person.Gold
	_person.Name = person.Name
	_person.Head = person.Imgurl
	_person.IP = person.ip
	_person.Address = person.minfo.Address
	_person.Sex = person.Sex
	_person.Online = true
	self.PersonMgr[person.Uid] = _person
	person.SendMsg("gamegoldsaimainfo", self.getinfo(person.Uid))

}

func (self *Game_GoldSaiMa) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold":
		person := self.GetPerson(msg.V.(*staticfunc.Msg_SynchroGold).Uid)
		if person != nil {
			person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal(person.Uid, person.Total)
		}
	case "gamebzwbets":
		self.GameBets(msg.Uid, msg.V.(*Msg_GameGoldBZW_Bets).Index, msg.V.(*Msg_GameGoldBZW_Bets).Gold)
	case "gamebzwgoon":
		self.GameGoOn(msg.Uid)
	}
}

func (self *Game_GoldSaiMa) OnBye() {

}

func (self *Game_GoldSaiMa) OnExit(uid int64) {
	value, ok := self.PersonMgr[uid]
	if ok {
		value.Online = false
		gold := value.Total - value.Gold
		if gold > 0 {
			GetRoomMgr().AddCard(value.Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
		} else if gold < 0 {
			GetRoomMgr().CostCard(value.Uid, staticfunc.TYPE_GOLD, -gold, self.room)
		}
		value.Gold = value.Total
	}
}

func (self *Game_GoldSaiMa) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_GoldSaiMa) OnBalance() {
	for _, value := range self.PersonMgr {
		value.Total += value.Bets

		gold := value.Total - value.Gold
		if gold > 0 {
			GetRoomMgr().AddCard(value.Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
		} else if gold < 0 {
			GetRoomMgr().CostCard(value.Uid, staticfunc.TYPE_GOLD, -gold, self.room)
		}
		value.Gold = value.Total
	}
}

func (self *Game_GoldSaiMa) RobotBet() {

}

func (self *Game_GoldSaiMa) OnTime() {
	if self.Time == 0 {
		return
	}

	if time.Now().Unix() < self.Time {
		//if self.Time-time.Now().Unix() <= 20 {
		//	active := make([]int, 0)
		//	notActive := make([]int, 0)
		//	for j := 0; j < len(self.PL); j++ {
		//		if self.PL[j] <= 10 {
		//			active = append(active, j)
		//		} else {
		//			notActive = append(notActive, j)
		//		}
		//	}
		//	for i := 0; i < lib.HF_GetRandom(2)+1; i++ {
		//		if lib.HF_GetRandom(100)+1 <= 80 { //!下活跃马
		//			self.GameBets(0, active[lib.HF_GetRandom(len(active))], ROBOTBETS[lib.HF_GetRandom(len(ROBOTBETS))])
		//		} else {
		//			self.GameBets(0, notActive[lib.HF_GetRandom(len(notActive))], ROBOTBETS[lib.HF_GetRandom(len(ROBOTBETS))])
		//		}
		//	}
		//}
		return
	}

	if !self.room.Begin {
		self.OnBegin()
	}
}
