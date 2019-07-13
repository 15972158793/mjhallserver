package gameserver

import (
	"fmt"
	"lib"
	"math"
	//"sort"
	"staticfunc"
	"time"
)

var HHDZ_BS []int = []int{0, 0, 2, 3, 4, 5, 10} //!　单张　9以下对子　9以上对子　顺子　金花　顺金　豹子

type Rec_HHDZ_Info struct {
	GameType int                   `json:"gametype"`
	Time     int64                 `json:"time"` //! 记录时间
	Info     []Son_Rec_HHDZ_Person `json:"info"`
}
type Son_Rec_HHDZ_Person struct {
	Uid    int64    `json:"uid"`
	Name   string   `json:"name"`
	Head   string   `json:"head"`
	Score  int      `json:"score"`
	Result [2][]int `json:"result"`
	Bets   [3]int   `json:"bets"`
}

type Game_GoldHHDZSeat struct {
	Person *Game_HHDZ_Person
	Robot  *lib.Robot
}

type Game_HHDZ struct {
	PersonMgr map[int64]*Game_HHDZ_Person  `json:"personmgr"`
	Bets      [3]map[*Game_HHDZ_Person]int `json:"bets"`
	Result    [2][]int                     `json:"result"`
	Time      int64                        `json:"time"`
	Seat      [6]Game_GoldHHDZSeat         `json:"seat"`
	Total     int                          `json:"total"` //!　这局下了多少
	Trend     []int                        `json:"trend"` //! 十位 0-黑赢 1-红赢    个位 牌型 1-单张 2-对子(9-A) 3-顺子 4-金花 5-顺金 6-豹子
	BetTime   int                          `json:"bettime"`
	room      *Room
	Robot     lib.ManyGameRobot //! 机器人结构
}

func NewGame_GoldHHDZ() *Game_HHDZ {
	game := new(Game_HHDZ)
	game.PersonMgr = make(map[int64]*Game_HHDZ_Person)
	for i := 0; i < 20; i++ {
		cardmgr := NewCard_LYC()
		game.Result[0] = cardmgr.Deal(3)
		game.Result[1] = cardmgr.Deal(3)
		trend := game.IsType(game.Result[0], game.Result[1])
		game.Trend = append(game.Trend, trend)
	}
	for i := 0; i < len(game.Bets); i++ {
		game.Bets[i] = make(map[*Game_HHDZ_Person]int)
	}

	return game
}

type Game_HHDZ_Person struct {
	Uid       int64  `json:"uid"`
	Gold      int    `json:"gold"`
	Total     int    `json:"total"`
	Win       int    `json:"win"`       //! 赢钱
	Cost      int    `json:"cost"`      //! 抽水
	Bets      int    `json:"bets"`      //! 单局下注
	BetInfo   [3]int `json:"betinfo"`   //!　下注详情
	BeBets    int    `json:"bebets"`    //! 上局下注
	BeBetInfo [3]int `json:"bebetinfo"` //!　上局下注详情
	Name      string `json:"name"`
	Head      string `json:"head"`
	Online    bool   `json:"online"`
	Round     int    `json:"round"` //! 未下注轮数
	Seat      int    `json:"seat"`
	IP        string `json:"ip"`
	Address   string `json:"address"`
	Sex       int    `json:"sex"`
}

type Msg_GameHHDZ_Total struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"`
}

type Son_GameHHDZ_Info struct {
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Total   int    `json:"total"` //! 总金币
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}
type Msg_GameHHDZ_List struct {
	Info []Son_GameHHDZ_Info `json:"info"`
}

type Msg_GameHHDZ_UpdSeat struct {
	Index   int    `json:"index"` //!　位置下标
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Total   int    `json:"total"`
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type Msg_GameHHDZ_Bets struct {
	Uid   int64 `json:"uid"`
	Index int   `json:"index"`
	Gold  int   `json:"gold"`
	Total int   `json:"total"`
}

type Msg_GameHHDZ_GoOn struct {
	Uid   int64  `json:"uid"`
	Gold  [3]int `json:"gold"`  //！ 三个区域分别下了多少钱
	Total int    `json:"total"` //! 总金币
}

type Msg_GameHHDZ_Info struct {
	Begin   bool                 `json:"begin"` //! 是否开始游戏
	Time    int64                `json:"time"`  //! 倒计时
	Seat    [6]Son_GameHHDZ_Info `json:"seat"`  //! 位置信息
	Bets    [3]int               `json:"bets"`  //! 三个下注区
	Total   int                  `json:"total"` //! 自己的钱
	Trend   []int                `json:"trend"` //! 走势
	Money   []int                `json:"money"`
	BetTime int                  `json:"bettime"`
}

type GameHHDZ_CanResult struct {
	Black []int `json:"black"`
	Red   []int `json:"red"`
}

//type GameGold_BigWin struct {
//	Uid  int64  `json:"uid"`
//	Name string `json:"name"`
//	Head string `json:"head"`
//	Win  int    `json:"win"`
//}

type Msg_GameHHDZ_Balance struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"` //! 总金币
	Win   int   `json:"win"`   //! 赢了多少
}

type Msg_GameHHDZ_End struct {
	Uid     int64    `json:"uid"`
	Name    string   `json:"name"`
	Head    string   `json:"head"`
	CT      []int    `json:"ct"`
	Trend   int      `json:"trend"`  //! 十位 0-黑赢 1-红赢    个位 牌型 1-单张 2-对子(9-A) 3-顺子 4-金花 5-顺金 6-豹子
	Result  [2][]int `json:"result"` //! 0-黑的牌 1-红的牌
	Money   []int    `json:"money"`
	BetTime int      `json:"bettime"`
}

func (self *Game_HHDZ) getInfo(uid int64, total int) *Msg_GameHHDZ_Info {
	var msg Msg_GameHHDZ_Info
	msg.Begin = self.room.Begin
	if self.Time == 0 {
		msg.Time = 0
	} else {
		msg.Time = self.Time - time.Now().Unix()
	}
	msg.Total = total
	msg.Trend = self.Trend
	msg.Money = lib.GetManyMoneyMgr().GetProperty(self.room.Type).Money
	msg.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime
	for i := 0; i < len(self.Bets); i++ {
		msg.Bets[i] = self.GetMoneyPos(i, true)
	}

	for i := 0; i < len(self.Seat); i++ {
		if self.Seat[i].Person != nil {
			msg.Seat[i].Uid = self.Seat[i].Person.Uid
			msg.Seat[i].Name = self.Seat[i].Person.Name
			msg.Seat[i].Head = self.Seat[i].Person.Head
			msg.Seat[i].Total = self.Seat[i].Person.Total
			msg.Seat[i].IP = self.Seat[i].Person.IP
			msg.Seat[i].Address = self.Seat[i].Person.Address
			msg.Seat[i].Sex = self.Seat[i].Person.Sex
		} else if self.Seat[i].Robot != nil {
			msg.Seat[i].Uid = self.Seat[i].Robot.Id
			msg.Seat[i].Name = self.Seat[i].Robot.Name
			msg.Seat[i].Head = self.Seat[i].Robot.Head
			msg.Seat[i].Total = self.Seat[i].Robot.GetMoney()
			msg.Seat[i].IP = self.Seat[i].Robot.IP
			msg.Seat[i].Address = self.Seat[i].Robot.Address
			msg.Seat[i].Sex = self.Seat[i].Robot.Sex
		}
	}
	return &msg

}

func (self *Game_HHDZ) GetPerson(uid int64) *Game_HHDZ_Person {
	return self.PersonMgr[uid]
}

func (self *Game_HHDZ) GetMoneyPos(index int, robot bool) int { //! 获取该位置一共下了多少
	total := 0
	for _, value := range self.Bets[index] {
		total += value
	}
	if robot {
		for _, value := range self.Robot.RobotsBet[index] {
			total += value
		}
	}
	return total
}

func (self *Game_HHDZ_Person) SynchroGold(gold int) { //! 同步金币
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

func (self *Game_HHDZ) SendTotal(uid int64, total int) { //! 同步总分
	var msg Msg_GameHHDZ_Total
	msg.Uid = uid
	msg.Total = total

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if person.Seat < 0 {
		self.room.SendMsg(uid, "gamehhdztotal", &msg)
	} else {
		self.room.broadCastMsg("gamehhdztotal", &msg)
	}
}

func (self *Game_HHDZ) SetTime(t int) { //! 设置时间
	if t == 0 {
		self.Time = 0
	} else {
		self.Time = time.Now().Unix() + int64(t)
	}

	var msg Msg_SetTime
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	self.room.broadCastMsg("gametime", &msg)
}

func (self *Game_HHDZ) GamePlayerList(uid int64) { //! 获取无座玩家列表
	var msg Msg_GameHHDZ_List
	tmp := make(map[int64]Son_GameHHDZ_Info)
	for _, value := range self.PersonMgr {
		if value.Seat >= 0 {
			continue
		}

		var node Son_GameHHDZ_Info
		node.Uid = value.Uid
		node.Name = value.Name
		node.Total = value.Total
		node.Head = value.Head
		tmp[node.Uid] = node
	}
	for i := 0; i < len(self.Robot.Robots); i++ {
		if self.Robot.Robots[i].GetSeat() >= 0 {
			continue
		}

		var node Son_GameHHDZ_Info
		node.Uid = self.Robot.Robots[i].Id
		node.Name = self.Robot.Robots[i].Name
		node.Total = self.Robot.Robots[i].GetMoney()
		node.Head = self.Robot.Robots[i].Head
		tmp[node.Uid] = node
	}
	for _, value := range tmp {
		msg.Info = append(msg.Info, value)
	}
	self.room.SendMsg(uid, "gameplayerlist", &msg)
}

func (self *Game_HHDZ) GameSeat(uid int64, index int) { //! 坐下
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if index < 0 || index > 5 {
		return
	}

	if person.Total < lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("金币必须大于%d才能坐下", lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("金币必须大于%d才能坐下", lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("金币必须大于%d才能坐下", lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney/10000))
		}
		return
	}

	for i := 0; i < len(self.Seat); i++ {
		if self.Seat[i].Person == person {
			return
		}
	}

	if self.Seat[index].Person != nil {
		if person.Total <= self.Seat[index].Person.Total {
			self.room.SendErr(uid, "该位置已经有人坐了")
			return
		}
		//! 把原来在这个位置上的人挤下去
		self.Seat[index].Person.Seat = -1
	} else if self.Seat[index].Robot != nil {
		if person.Total <= self.Seat[index].Robot.GetMoney() {
			self.room.SendErr(uid, "该位置已经有人坐了")
			return
		}
		//! 把原来在这个位置上的人挤下去
		self.Seat[index].Robot.SetSeat(-1)
	}

	self.Seat[index].Person = person
	self.Seat[index].Robot = nil
	person.Seat = index

	var msg Msg_GameHHDZ_UpdSeat
	msg.Uid = uid
	msg.Index = index
	msg.Head = person.Head
	msg.Name = person.Name
	msg.Total = person.Total
	msg.IP = person.IP
	msg.Address = person.Address
	msg.Sex = person.Sex
	self.room.broadCastMsg("gamehhdzseat", &msg)
}

//! 机器人坐下
func (self *Game_HHDZ) RobotSeat(index int, robot *lib.Robot) {
	if index < 0 || index > 5 {
		return
	}

	if robot.GetMoney() < lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney {
		return
	}

	for i := 0; i < len(self.Seat); i++ {
		if self.Seat[i].Robot == robot {
			return
		}
	}

	if self.Seat[index].Person != nil || self.Seat[index].Robot != nil {
		return
	}

	self.Seat[index].Person = nil
	self.Seat[index].Robot = robot
	robot.SetSeat(index)

	var msg Msg_GameHHDZ_UpdSeat
	msg.Uid = robot.Id
	msg.Index = index
	msg.Head = robot.Head
	msg.Name = robot.Name
	msg.Total = robot.GetMoney()
	msg.IP = robot.IP
	msg.Address = robot.Address
	msg.Sex = robot.Sex
	self.room.broadCastMsg("gamehhdzseat", &msg)
}

func (self *Game_HHDZ) GameBets(uid int64, index int, gold int) { //! 下注
	if uid == 0 {
		return
	}

	if index < 0 || index > 2 {
		return
	}

	if gold <= 0 {
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() >= int64(self.BetTime-2) {
		self.room.SendErr(uid, "正在开奖，请稍后下注")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if person.Total < lib.GetManyMgr().GetProperty(self.room.Type).MinBet {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("%d金币以上才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("%d金币以上才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("%d金币以上才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/10000))
		}
		return
	}

	if person.Total < gold {
		self.room.SendErr(uid, "您的金币不足，请前往充值")
		return
	}

	if person.Bets+gold > lib.GetManyMgr().GetProperty(self.room.Type).MaxBet {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet/10000))
		}
		return
	}

	person.Bets += gold
	person.Total -= gold
	person.BetInfo[index] += gold
	person.Round = 0
	self.Total += gold
	self.Bets[index][person] += gold

	var msg Msg_GameHHDZ_Bets
	msg.Uid = uid
	msg.Index = index
	msg.Gold = gold
	msg.Total = person.Total
	self.room.broadCastMsg("gamehhdzbets", &msg)

}

func (self *Game_HHDZ) GameGoOn(uid int64) { //! 续压
	if uid == 0 {
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() >= int64(self.BetTime-2) {
		self.room.SendErr(uid, "正在开奖,请稍后下注")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if person.Total < lib.GetManyMgr().GetProperty(self.room.Type).MinBet {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("%d金币以上才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("%d金币以上才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("%d金币以上才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/10000))
		}
		return
	}

	if person.Total < person.BeBets {
		self.room.SendErr(uid, "您的金币不足，请前往充值")
		return
	}

	if person.Bets+person.BeBets > lib.GetManyMgr().GetProperty(self.room.Type).MaxBet {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet/10000))
		}
		return
	}

	person.Bets += person.BeBets
	person.Total -= person.BeBets
	self.Total += person.BeBets
	for i := 0; i < len(person.BeBetInfo); i++ {
		person.BetInfo[i] = person.BeBetInfo[i]
		self.Bets[i][person] += person.BeBetInfo[i]
	}
	person.Round = 0

	var msg Msg_GameHHDZ_GoOn
	msg.Uid = uid
	msg.Gold = person.BeBetInfo
	msg.Total = person.Total
	self.room.broadCastMsg("gamehhdzgoon", &msg)
}

func (self *Game_HHDZ) IsType(card []int, _card []int) int { //! 十位 0-黑胜 1-红胜  个位 0-单张 1-对子 2-可得分对子 3-顺子 4-金花 5-顺金 6-豹子
	result := ZjhCardCompare2(card, _card) * 10
	cardtype, _ := GetZjhType(card)
	cardtype2, _ := GetZjhType(_card)
	if cardtype > cardtype2 {
		if cardtype/100 == 1 {
			return result
		}
		if cardtype/100 == 2 && cardtype%100 < 9 {
			return result + 1
		}
		return result + cardtype/100
	} else {
		if cardtype2/100 == 1 {
			return result
		}
		if cardtype2/100 == 2 && cardtype2%100 < 9 {
			return result + 1
		}
		return result + cardtype2/100
	}
	return -1
}

func (self *Game_HHDZ) GetDealWin(card []int, _card []int) int { //! 庄家能赢多少
	lost := 0
	ct := self.IsType(card, _card)
	if ct/10 == 0 {
		lost += self.GetMoneyPos(0, false) * 2
	} else if ct/10 == 1 {
		lost += self.GetMoneyPos(1, false) * 2
	}
	lost += self.GetMoneyPos(2, false) * HHDZ_BS[ct%10]

	return self.Total - lost
}

func (self *Game_HHDZ) OnBegin() {
	if self.room.IsBye() {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "房间已解散")
		return
	}
	self.room.Begin = true

	//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "max : ", lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax, " min : ", lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin)

	lst := make([]GameHHDZ_CanResult, 0)
	winLst := make([]GameHHDZ_CanResult, 0)
	lostLst := make([]GameHHDZ_CanResult, 0)

	for i := 0; i < 100; i++ {
		card := NewCard_LYC()
		var msg GameHHDZ_CanResult
		msg.Black = card.Deal(3)
		msg.Red = card.Deal(3)

		dealwin := self.GetDealWin(msg.Black, msg.Red)
		if GetServer().HhdzSysMoney[self.room.Type%210000]+int64(dealwin) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && GetServer().HhdzSysMoney[self.room.Type%210000]+int64(dealwin) <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
			lst = append(lst, msg)
		}
		if dealwin >= 0 {
			winLst = append(winLst, msg)
		} else {
			lostLst = append(lostLst, msg)
		}
	}

	if len(lst) == 0 {
		if GetServer().HhdzSysMoney[self.room.Type%210000] >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax && len(lostLst) > 0 { //! 一定输
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------- 庄一定输", len(lostLst))
			index := lib.HF_GetRandom(len(lostLst))
			self.Result[0] = lostLst[index].Black
			self.Result[1] = lostLst[index].Red
		} else if GetServer().HhdzSysMoney[self.room.Type%210000] <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && len(winLst) > 0 { //! 一定赢
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------- 庄一定赢 ", len(winLst))
			index := lib.HF_GetRandom(len(winLst))
			self.Result[0] = winLst[index].Black
			self.Result[1] = winLst[index].Red
		} else {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------- 纯随机")
			card := NewCard_LYC()
			self.Result[0] = card.Deal(3)
			self.Result[1] = card.Deal(3)
		}
	} else {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "------随机lst")
		index := lib.HF_GetRandom(len(lst))
		self.Result[0] = lst[index].Black
		self.Result[1] = lst[index].Red
	}

	self.OnEnd()
}

func (self *Game_HHDZ) OnEnd() {
	self.room.Begin = false

	trend := self.IsType(self.Result[0], self.Result[1])
	tmp := make([]int, 0)
	tmp = append(tmp, trend)
	tmp = append(tmp, self.Trend...)
	if len(tmp) > 20 {
		tmp = tmp[0:20]
	}
	self.Trend = tmp

	dealwin := 0

	for i := 0; i < len(self.Bets); i++ {
		if i != 2 {
			if trend/10 == i {
				for key, value := range self.Bets[i] {
					winMoney := value * 2
					dealwin -= winMoney
					key.Win += winMoney
					key.Cost += int(math.Ceil(float64(winMoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
				}

				for key, value := range self.Robot.RobotsBet[i] {
					winMoney := value * 2
					key.AddWin(winMoney)
					key.AddCost(int(math.Ceil(float64(winMoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
				}
			}
		} else {
			if trend%10 != 0 {
				for key, value := range self.Bets[i] {
					winMoney := value * HHDZ_BS[trend%10]
					if winMoney > 0 {
						dealwin -= winMoney
						key.Win += winMoney
						key.Cost += int(math.Ceil(float64(winMoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
					}
				}

				for key, value := range self.Robot.RobotsBet[i] {
					winMoney := value * HHDZ_BS[trend%10]
					if winMoney > 0 {
						key.AddWin(winMoney)
						key.AddCost(int(math.Ceil(float64(winMoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
					}
				}
			}
		}
	}

	dealwin += self.Total
	if dealwin != 0 {
		GetServer().SqlBZWLog(&SQL_BZWLog{1, dealwin, time.Now().Unix(), self.room.Type})
	}
	if dealwin > 0 {
		cost := int(math.Ceil(float64(dealwin) * float64(lib.GetManyMgr().GetProperty(self.room.Type).DealCost) / 100.0))
		dealwin -= cost
	}
	GetServer().SetHhdzMoney(self.room.Type%210000, GetServer().HhdzSysMoney[self.room.Type%210000]+int64(dealwin))

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "庄家收益 ： ", dealwin)

	var bigwin *GameGold_BigWin = nil
	for _, value := range self.PersonMgr {
		if value.Win > 0 {
			value.Win -= value.Cost
			GetServer().SqlAgentGoldLog(value.Uid, value.Cost, self.room.Type)
			GetServer().SqlAgentBillsLog(value.Uid, value.Cost/2, self.room.Type)
			value.Total += value.Win

			var msg Msg_GameHHDZ_Balance
			msg.Uid = value.Uid
			msg.Win = value.Win
			msg.Total = value.Total
			find := false
			for j := 0; j < len(self.Seat); j++ {
				if self.Seat[j].Person == value {
					self.room.broadCastMsg("gamehhdzbalance", &msg)
					find = true
					break
				}
			}
			if !find {
				self.room.SendMsg(value.Uid, "gamehhdzbalance", &msg)
			}

			if bigwin == nil {
				bigwin = &GameGold_BigWin{value.Uid, value.Name, value.Head, value.Win}
			} else if value.Win > bigwin.Win {
				bigwin = &GameGold_BigWin{value.Uid, value.Name, value.Head, value.Win}
			}
		} else if value.Win-value.Bets < 0 {
			cost := int(math.Ceil(float64(value.Win-value.Bets) * float64(lib.GetManyMgr().GetProperty(self.room.Type).Cost) / 200.0))
			GetServer().SqlAgentBillsLog(value.Uid, cost, self.room.Type)
		}

		if value.Bets > 0 {
			var record Rec_HHDZ_Info
			record.Time = time.Now().Unix()
			record.GameType = self.room.Type
			var rec Son_Rec_HHDZ_Person
			rec.Uid = value.Uid
			rec.Name = value.Name
			rec.Head = value.Head
			rec.Score = value.Win - value.Bets
			rec.Result = self.Result
			rec.Bets = value.BetInfo
			record.Info = append(record.Info, rec)
			GetServer().InsertRecord(self.room.Type, value.Uid, lib.HF_JtoA(&record), rec.Score)
		}
	}

	for i := 0; i < len(self.Robot.Robots); i++ {
		if self.Robot.Robots[i].GetWin() > 0 {
			self.Robot.Robots[i].AddWin(-self.Robot.Robots[i].GetCost())
			self.Robot.Robots[i].AddMoney(self.Robot.Robots[i].GetWin())

			for j := 0; j < len(self.Seat); j++ {
				if self.Seat[j].Robot == self.Robot.Robots[i] {
					var msg Msg_GameHHDZ_Balance
					msg.Uid = self.Robot.Robots[i].Id
					msg.Total = self.Robot.Robots[i].GetMoney()
					msg.Win = self.Robot.Robots[i].GetWin()
					self.room.broadCastMsg("gamehhdzbalance", &msg)
					break
				}
			}

			if bigwin == nil {
				bigwin = &GameGold_BigWin{self.Robot.Robots[i].Id, self.Robot.Robots[i].Name, self.Robot.Robots[i].Head, self.Robot.Robots[i].GetWin()}
			} else if self.Robot.Robots[i].GetWin() > bigwin.Win {
				bigwin = &GameGold_BigWin{self.Robot.Robots[i].Id, self.Robot.Robots[i].Name, self.Robot.Robots[i].Head, self.Robot.Robots[i].GetWin()}
			}
		}
	}

	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 12
	self.SetTime(self.BetTime)

	{ //! 大赢家
		var msg Msg_GameHHDZ_End
		msg.Result = self.Result
		msg.Trend = trend
		msg.Money = lib.GetManyMoneyMgr().GetProperty(self.room.Type).Money
		msg.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime
		msg.CT = make([]int, 2)
		ct, _ := GetZjhType(self.Result[0])
		msg.CT[0] = ct / 100
		ct, _ = GetZjhType(self.Result[1])
		msg.CT[1] = ct / 100
		if bigwin != nil {
			msg.Uid = bigwin.Uid
			msg.Name = bigwin.Name
			msg.Head = bigwin.Head
		}
		self.room.broadCastMsg("gamehhdzend", &msg)
	}
	self.Total = 0

	//! 清理玩家
	for key, value := range self.PersonMgr {
		if value.Online {
			if value.Seat < 0 {
				value.Round++
			}

			if value.Round >= 5 && GetPersonMgr().GetPerson(value.Uid) == nil {
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

		for j := 0; j < len(self.Seat); j++ {
			if self.Seat[j].Person == value {
				self.Seat[j].Person = nil
				var msg Msg_GameHHDZ_UpdSeat
				msg.Index = j
				self.room.broadCastMsg("gamehhdzseat", &msg)
				break
			}
		}

		delete(self.PersonMgr, key)
	}

	//! 载入机器人
	self.Robot.Init(3, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)

	for i := 0; i < len(self.Bets); i++ {
		self.Bets[i] = make(map[*Game_HHDZ_Person]int)
	}

	//! 坐下的人是否还能坐下
	for i := 0; i < len(self.Seat); i++ {
		if self.Seat[i].Person == nil {
			continue
		}

		if self.Seat[i].Person.Total < lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney {
			self.Seat[i].Person.Seat = -1
			var msg Msg_GameHHDZ_UpdSeat
			msg.Index = i
			self.room.broadCastMsg("gamehhdzseat", &msg)
			self.Seat[i].Person = nil
		}
	}
	for i := 0; i < len(self.Seat); i++ {
		if self.Seat[i].Robot == nil {
			continue
		}
		find := false
		for j := 0; j < len(self.Robot.Robots); j++ {
			if self.Robot.Robots[j] == self.Seat[i].Robot {
				find = true
				break
			}
		}
		if !find || self.Seat[i].Robot.GetSeat() != i || self.Seat[i].Robot.GetMoney() < lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney {
			self.Seat[i].Robot.SetSeat(-1)
			var msg Msg_GameHHDZ_UpdSeat
			msg.Index = i
			self.room.broadCastMsg("gamehhdzseat", &msg)
			self.Seat[i].Robot = nil
		}
	}
}

func (self *Game_HHDZ) OnInit(room *Room) {
	self.room = room
	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 12
	//! 载入机器人
	self.Robot.Init(3, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)
}

func (self *Game_HHDZ) OnRobot(robot *lib.Robot) {

}

func (self *Game_HHDZ) OnSendInfo(person *Person) {
	if self.Time == 0 {
		self.SetTime(lib.GetManyMgr().GetProperty(self.room.Type).BetTime)
	}

	value, ok := self.PersonMgr[person.Uid]
	if ok {
		value.Online = true
		value.Round = 0
		value.IP = person.ip
		value.Address = person.minfo.Address
		value.Sex = person.Sex
		value.SynchroGold(person.Gold)
		person.SendMsg("gamehhdzinfo", self.getInfo(person.Uid, value.Total))
		return
	}

	_person := new(Game_HHDZ_Person)
	_person.Uid = person.Uid
	_person.Seat = -1
	_person.Gold = person.Gold
	_person.Total = person.Gold
	_person.Name = person.Name
	_person.Head = person.Imgurl
	_person.IP = person.ip
	_person.Address = person.minfo.Address
	_person.Sex = person.Sex
	_person.Online = true
	self.PersonMgr[person.Uid] = _person
	person.SendMsg("gamehhdzinfo", self.getInfo(person.Uid, person.Gold))
}

func (self *Game_HHDZ) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold": //! 同步金币
		person := self.GetPerson(msg.V.(*staticfunc.Msg_SynchroGold).Uid)
		if person != nil {
			person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal(person.Uid, person.Total)
		}
	case "gamebzwbets": //! 下注
		self.GameBets(msg.Uid, msg.V.(*Msg_GameGoldBZW_Bets).Index, msg.V.(*Msg_GameGoldBZW_Bets).Gold)
	case "gamebzwgoon":
		self.GameGoOn(msg.Uid)
	case "gamebzwseat":
		self.GameSeat(msg.Uid, msg.V.(*Msg_GameGoldBZW_Seat).Index)
	case "gameplayerlist":
		self.GamePlayerList(msg.Uid)
	}
}

func (self *Game_HHDZ) OnBye() {

}

func (self *Game_HHDZ) OnExit(uid int64) {
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

func (self *Game_HHDZ) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_HHDZ) OnIsBets(uid int64) bool {
	value, ok := self.PersonMgr[uid]
	if ok {
		return value.Bets > 0
	}
	return false
}

func (self *Game_HHDZ) OnBalance() {
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

func (self *Game_HHDZ) OnTime() {
	if self.Time == 0 {
		return
	}

	for i := 0; i < len(self.Robot.Robots); i++ {
		if self.Robot.Robots[i].GetSeat() >= 0 {
			continue
		}
		if lib.HF_GetRandom(100) < 90 {
			continue
		}
		self.RobotSeat(lib.HF_GetRandom(8), self.Robot.Robots[i])
	}

	if time.Now().Unix() < self.Time {
		if self.Time-time.Now().Unix() >= int64(lib.GetManyMgr().GetProperty(self.room.Type).BetTime) {
			return
		}

		for i := 0; i < len(self.Robot.Robots); i++ {
			if lib.HF_GetRandom(100) >= 100-lib.GetRobotMgr().GetRobotSet(self.room.Type).BetRate {
				continue
			}

			index, gold, _ := self.Robot.GameBets(self.Robot.Robots[i])
			if gold == 0 {
				continue
			}
			var msg Msg_GameGoldBZW_Bets
			msg.Uid = self.Robot.Robots[i].Id
			msg.Index = index
			msg.Gold = gold
			msg.Total = self.Robot.Robots[i].GetMoney()
			self.room.broadCastMsg("gamehhdzbets", &msg)
		}
		return
	}

	if !self.room.Begin {
		self.OnBegin()
	}
}
