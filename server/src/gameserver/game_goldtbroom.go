package gameserver

import (
	"fmt"
	"lib"
	"math"
	"staticfunc"
	"time"
)

var TB_TIME = 30

type Game_GoldTBRoom struct {
	PersonMgr map[int64]*Game_GoldTB_Person  `json:"personmgr"`
	Bets      [2]map[*Game_GoldTB_Person]int `json:"bets"`
	Kill      [2]*Son_GameGoldTB_Info        `json:"kill"` //! 0-杀单玩家 1-杀双玩家
	Result    [2]int                         `json:"result"`
	Lucky     [2]int                         `json:"lucky"` //! 幸运数
	Dealer    *Game_GoldTB_Person            `json:"dealer"`
	Round     int                            `json:"round"`   //! 连庄轮数
	DownUid   int64                          `json:"downuid"` //! 下庄的人
	Time      int64                          `json:"time"`
	LstDeal   []*Game_GoldTB_Person          `json:"lstdeal"` //! 上庄列表
	Seat      [6]*Game_GoldTB_Person         `json:"seat"`
	Total     int                            `json:"total"` //!这局一共下了多少钱
	Trend     [][2]int                       `json:"trend"` //!　走势

	room *Room
}

func NewGame_GoldTBRoom() *Game_GoldTBRoom {
	game := new(Game_GoldTBRoom)
	game.PersonMgr = make(map[int64]*Game_GoldTB_Person)
	game.Lucky = [2]int{lib.HF_GetRandom(6) + 1, lib.HF_GetRandom(6) + 1}
	for i := 0; i < 2; i++ {
		game.Kill[i] = nil
	}
	for i := 0; i < len(game.Bets); i++ {
		game.Bets[i] = make(map[*Game_GoldTB_Person]int)
	}
	return game
}

func (self *Game_GoldTBRoom) GetPerson(uid int64) *Game_GoldTB_Person {
	return self.PersonMgr[uid]
}

func (self *Game_GoldTBRoom) GameSeat(uid int64, index int) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if index < 0 || index > 7 {
		return
	}

	if self.Dealer == person {
		self.room.SendErr(uid, "庄家无法坐下")
		return
	}

	for i := 0; i < len(self.Seat); i++ {
		if self.Seat[i] == person {
			return
		}
	}

	if self.Seat[index] != nil {
		if person.Total <= self.Seat[index].Total {
			self.room.SendErr(uid, "该位置已经有人坐了")
			return
		}
		self.Seat[index].Seat = -1
	}

	self.Seat[index] = person
	person.Seat = index

	var msg Msg_GameGoldTB_UpdSeat
	msg.Uid = uid
	msg.Index = index
	msg.Head = person.Head
	msg.Name = person.Name
	msg.Total = person.Total
	msg.IP = person.IP
	msg.Address = person.Address
	msg.Sex = person.Sex
	self.room.broadCastMsg("gamegoldtbseat", &msg)
}

//! 同步总分
func (self *Game_GoldTBRoom) SendTotal(uid int64, total int) {
	var msg Msg_GameGoldTB_Total
	msg.Uid = uid
	msg.Total = total

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if person.Seat < 0 {
		self.room.SendMsg(uid, "gamegoldtotal", &msg)
	} else {
		self.room.broadCastMsg("gamegoldtotal", &msg)
	}
}

//! 设置时间
func (self *Game_GoldTBRoom) SetTime(t int) {
	if t == 0 {
		self.Time = 0
	} else {
		self.Time = time.Now().Unix() + int64(t)
	}
	var msg Msg_SetTime
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	self.room.broadCastMsg("gametime", &msg)
}

//! 得到这个位置下了多少钱
func (self *Game_GoldTBRoom) GetMoneyPos(index int) int {
	total := 0
	for _, value := range self.Bets[index] {
		total += value
	}
	return total
}

//! 是否是幸运数
func (self *Game_GoldTBRoom) IsLucky() bool {
	return (self.Result[0] == self.Lucky[0] && self.Result[1] == self.Lucky[1]) || (self.Result[0] == self.Lucky[1] && self.Result[1] == self.Lucky[0])
}

//! 下注
func (self *Game_GoldTBRoom) GameBets(uid int64, index int, gold int) {
	if uid == 0 {
		return
	}

	if index != 0 && index != 1 {
		return
	}

	if gold <= 0 {
		return
	}

	if self.Time == 0 {
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() <= 5 {
		self.room.SendErr(uid, "现在是抢杀阶段，不能下注")
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() >= int64(TB_TIME-2) {
		self.room.SendErr(uid, "正在开奖，请稍后下注")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if self.Dealer == person {
		self.room.SendErr(uid, "庄家不用下注")
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
		self.room.SendErr(uid, "您的金币不足，请前往充值")
		return
	}

	if person.Bets+gold > lib.GetManyMgr().GetProperty(self.room.Type).MaxBet {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d。", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d。", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d万。", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet/10000))
		}
		return
	}
	dealtotal := 0
	if self.Dealer != nil {
		dealtotal = self.Dealer.Total
	}

	if self.GetMoneyPos(index)+gold > dealtotal {
		self.room.SendErr(uid, "庄家金币不足，该位置无法下注")
		return
	}

	person.Bets += gold
	person.Total -= gold
	person.BetInfo[index] += gold
	person.Round = 0
	self.Bets[index][person] += gold

	var msg Msg_GameGoldTB_Bets
	msg.Uid = uid
	msg.Index = index
	msg.Gold = gold
	msg.Total = person.Total
	msg.BetInfo = person.BetInfo
	for i := 0; i < len(self.Bets); i++ {
		msg.GameTotal[i] = self.GetMoneyPos(i)
	}
	self.room.broadCastMsg("gamegoldtbbets", &msg)

	/*
		player := GetPersonMgr().GetPerson(person.Uid)
		if player != nil {
			//	player.OnGaming(self.room.Type)
		}
	*/
	if dealtotal == self.GetMoneyPos(0) && dealtotal == self.GetMoneyPos(1) {
		self.SetTime(5)
	}

}

//! 续押
func (self *Game_GoldTBRoom) GameGoOn(uid int64) {
	if uid == 0 {
		return
	}

	if self.Time == 0 {
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() <= 5 {
		self.room.SendErr(uid, "现在是抢杀阶段，不能下注")
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() >= int64(TB_TIME-2) {
		self.room.SendErr(uid, "正在开奖，请稍后下注")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if self.Dealer == person {
		self.room.SendErr(uid, "庄家不用下注")
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
		self.room.SendErr(uid, "您的金币不足，请前往充值")
		return
	}

	if person.Bets+person.BeBets > lib.GetManyMgr().GetProperty(self.room.Type).MaxBet {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d。", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d。", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d万。", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet/10000))
		}
		return
	}

	dealtotal := 0
	if self.Dealer != nil {
		dealtotal = self.Dealer.Total
	}

	for i := 0; i < len(person.BeBetInfo); i++ {
		if self.GetMoneyPos(i)+person.BeBetInfo[i] > dealtotal {
			self.room.SendErr(uid, "庄家金币不足，该位置无法下注")
			return
		}
	}

	person.Bets += person.BeBets
	person.Total -= person.BeBets
	for i := 0; i < len(person.BeBetInfo); i++ {
		person.BetInfo[i] += person.BeBetInfo[i]
		self.Bets[i][person] += person.BeBetInfo[i]
	}
	person.Round = 0

	var msg Msg_GameGoldTB_GoOn
	msg.Uid = uid
	msg.Gold = person.BeBetInfo
	msg.Total = person.Total
	for i := 0; i < 2; i++ {
		msg.GameTotal[i] = self.GetMoneyPos(i)
	}
	self.room.broadCastMsg("gamegoldtbgoon", &msg)

	if dealtotal == self.GetMoneyPos(0) && dealtotal == self.GetMoneyPos(1) {
		self.SetTime(5)
	}
}

//! 上庄
func (self *Game_GoldTBRoom) GameUpDeal(uid int64) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if person.Total < lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("金币必须大于%d才能上庄", lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("金币必须大于%d才能上庄", lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("金币必须大于%d万才能上庄", lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney/10000))
		}
		return
	}

	if self.Dealer == person {
		self.DownUid = 0
	} else {
		for i := 0; i < len(self.LstDeal); i++ {
			if self.LstDeal[i] == person {
				self.room.SendErr(uid, "您已经在上庄列表中，请等待上庄")
				return
			}
		}
		if len(self.LstDeal) == 0 {
			self.Round = 0
		}
		self.LstDeal = append(self.LstDeal, person)
	}
	person.Round = 0

	if self.Dealer == nil {
		self.ChageDeal()
		if self.Dealer != nil {
			self.SetTime(22)
		}
	} else {
		var msg Msg_GameGoldTB_DealList
		msg.Type = 0
		msg.Info = make([]Son_GameGoldTB_Info, 0)
		for i := 0; i < len(self.LstDeal); i++ {
			msg.Info = append(msg.Info, Son_GameGoldTB_Info{self.LstDeal[i].Uid, self.LstDeal[i].Name, self.LstDeal[i].Head, self.LstDeal[i].Total, self.LstDeal[i].IP, self.LstDeal[i].Address, self.LstDeal[i].Sex})
		}
		self.room.SendMsg(uid, "gamegoldtbdeal", &msg)
	}

}

//! 下庄
func (self *Game_GoldTBRoom) GameReDeal(uid int64) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if self.Dealer == person {
		self.DownUid = uid
		self.room.SendErr(uid, "您已成功下庄，请等待本局结束")
	} else {
		for i := 0; i < len(self.LstDeal); i++ {
			if self.LstDeal[i] == person {
				copy(self.LstDeal[i:], self.LstDeal[i+1:])
				self.LstDeal = self.LstDeal[:len(self.LstDeal)-1]
				break
			}
		}
	}

	var msg Msg_GameGoldTB_DealList
	msg.Type = 1
	msg.Info = make([]Son_GameGoldTB_Info, 0)
	self.room.SendMsg(uid, "gametbdeal", &msg)
}

func (self *Game_GoldTBRoom) getInfo(uid int64, total int, mybets [2]int) *Msg_GameGoldTB_Info {
	var msg Msg_GameGoldTB_Info
	msg.Begin = self.room.Begin
	if self.Time == 0 {
		msg.Time = 0
	} else {
		msg.Time = self.Time - time.Now().Unix()
	}
	msg.Total = total
	msg.Trend = self.Trend
	msg.Lucky = self.Lucky
	msg.MyBets = mybets
	msg.Money = lib.GetManyMoneyMgr().GetProperty(self.room.Type).Money
	if self.Dealer != nil && self.Dealer.Uid == uid {
		msg.IsDeal = true
	} else {
		for i := 0; i < len(self.LstDeal); i++ {
			if self.LstDeal[i].Uid == uid {
				msg.IsDeal = true
				break
			}
		}
	}
	for i := 0; i < len(self.Kill); i++ {
		if self.Kill[i] != nil {
			msg.Kill[i].Uid = self.Kill[i].Uid
			msg.Kill[i].Name = self.Kill[i].Name
			msg.Kill[i].Head = self.Kill[i].Head
			msg.Kill[i].Total = self.Kill[i].Total
			msg.Kill[i].Address = self.Kill[i].Address
			msg.Kill[i].Sex = self.Kill[i].Sex
			msg.Kill[i].IP = self.Kill[i].IP
		}
	}
	for i := 0; i < len(self.Bets); i++ {
		msg.Bets[i] = self.GetMoneyPos(i)
	}
	for i := 0; i < 6; i++ {
		if self.Seat[i] != nil {
			msg.Seat[i].Uid = self.Seat[i].Uid
			msg.Seat[i].Name = self.Seat[i].Name
			msg.Seat[i].Head = self.Seat[i].Head
			msg.Seat[i].Total = self.Seat[i].Total
			msg.Seat[i].IP = self.Seat[i].IP
			msg.Seat[i].Address = self.Seat[i].Address
			msg.Seat[i].Sex = self.Seat[i].Sex
		}
	}
	if self.Dealer != nil {
		msg.Dealer.Uid = self.Dealer.Uid
		msg.Dealer.Name = self.Dealer.Name
		msg.Dealer.Head = self.Dealer.Head
		msg.Dealer.Total = self.Dealer.Total
		msg.Dealer.IP = self.Dealer.IP
		msg.Dealer.Address = self.Dealer.Address
		msg.Dealer.Sex = self.Dealer.Sex
	}
	return &msg
}

//! 换庄
func (self *Game_GoldTBRoom) ChageDeal() {
	if self.Dealer != nil {
		self.Dealer.Seat = -1
	}

	self.Dealer = nil
	for len(self.LstDeal) > 0 {
		if self.LstDeal[0].Total >= lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney {
			self.Dealer = self.LstDeal[0]
			self.LstDeal = self.LstDeal[1:]
			self.Dealer.Seat = 100
			for i := 0; i < len(self.Seat); i++ {
				if self.Seat[i] == self.Dealer {
					var msg Msg_GameGoldTB_UpdSeat
					msg.Index = i
					self.room.broadCastMsg("gamegoldtbseat", &msg)
					self.Seat[i] = nil
					break
				}
			}
			break
		} else {
			self.LstDeal = self.LstDeal[1:]
		}
	}
	self.DownUid = 0
	self.Round = 0
	self.Lucky = [2]int{lib.HF_GetRandom(6) + 1, lib.HF_GetRandom(6) + 1}

	var msg Msg_GameGoldTB_Deal
	if self.Dealer != nil {
		msg.Uid = self.Dealer.Uid
		msg.Name = self.Dealer.Name
		msg.Head = self.Dealer.Head
		msg.Lucky = self.Lucky
		msg.Total = self.Dealer.Total
		msg.IP = self.Dealer.IP
		msg.Address = self.Dealer.Address
		msg.Sex = self.Dealer.Sex
	}
	self.room.broadCastMsg("gamerob", &msg)
}

//! 申请无座玩家
func (self *Game_GoldTBRoom) GamePlayerList(uid int64) {
	has := false
	var msg Msg_GameGoldTB_List
	for _, value := range self.PersonMgr {
		if value.Seat >= 0 {
			continue
		}

		var node Son_GameGoldTB_Info
		node.Uid = value.Uid
		node.Name = value.Name
		node.Total = value.Total
		node.Head = value.Head
		msg.Info = append(msg.Info, node)
		if uid == value.Uid {
			has = true
		}
		if len(msg.Info) >= 100 {
			break
		}
	}
	if !has {
		person := self.GetPerson(uid)
		if person == nil {
			return
		}
		if person.Seat < 0 {
			msg.Info = append(msg.Info, Son_GameGoldTB_Info{person.Uid, person.Name, person.Head, person.Total, person.IP, person.Address, person.Sex})
		}
	}
	self.room.SendMsg(uid, "gameplayerlist", &msg)
}

//! 抢杀
func (self *Game_GoldTBRoom) GetKill(uid int64, index int) {
	if uid == 0 {
		return
	}

	if index != 0 && index != 1 {
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() > 5 {
		self.room.SendErr(uid, "现在不是抢杀时间。")
		return
	}
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if person.Total < lib.GetManyMgr().GetProperty(self.room.Type).MinBet {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("%d金币以上才能杀", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("%d金币以上才能杀", lib.GetManyMgr().GetProperty(self.room.Type).MinBet))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("%d万金币以上才能杀", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/10000))
		}
		return
	}

	if self.Dealer == person {
		self.room.SendErr(uid, "庄家不用抢杀")
		return
	}

	if person.BetInfo[index] > 0 {
		self.room.SendErr(uid, "不能杀已下注的区域")
		return
	}

	if (index == 0 && person.Total < self.GetMoneyPos(0)) || (index == 1 && person.Total < self.GetMoneyPos(1)) {
		self.room.SendErr(uid, "您的金币不足，请前往充值")
		return
	}

	//if self.Dealer != nil && (self.Kill[0] != nil || self.Kill[1] != nil) {
	//	return
	//}

	if self.Kill[index] != nil {
		self.room.SendErr(uid, "该位置已有人抢杀")
		return
	}

	var son Son_GameGoldTB_Info
	son.Uid = person.Uid
	son.Total = person.Total
	son.Sex = person.Sex
	son.Name = person.Name
	son.IP = person.IP
	son.Head = person.Head
	son.Address = person.Address

	self.Kill[index] = &son
	person.Round = 0

	var msg Msg_GameGoldTB_Kill
	msg.Uid = uid
	msg.Head = person.Head
	msg.Name = person.Name
	msg.Type = index
	self.room.broadCastMsg("gamegoldtbkill", &msg)

	if self.Kill[0] != nil && self.Kill[1] != nil {
		self.OnBegin()
		return
	}

	//if self.Dealer == nil { //! 系统庄

	//} else {
	//	if self.Kill[0] != nil || self.Kill[1] != nil {
	//		self.OnBegin()
	//		return
	//	}
	//}
}

func (self *Game_GoldTBRoom) IsType() int { //! 0-单 1-双
	if (self.Result[0]+self.Result[1])%2 == 0 {
		return 1
	}
	return 0
}

func (self *Game_GoldTBRoom) OnBegin() {
	if self.room.IsBye() {
		return
	}

	if self.Dealer == nil {
		self.SetTime(0)
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "------- onbegin  没有庄家")
		return
	}

	self.room.Begin = true

	self.Result[0] = lib.HF_GetRandom(6) + 1
	self.Result[1] = lib.HF_GetRandom(6) + 1

	self.OnEnd()
}

func (self *Game_GoldTBRoom) OnEnd() {
	self.room.Begin = false

	tmp := [][2]int{self.Result}
	tmp = append(tmp, self.Trend...)
	if len(tmp) > 20 {
		tmp = tmp[0:20]
	}
	self.Trend = tmp
	trend := self.IsType()

	var money [2]int = [2]int{self.GetMoneyPos(0), self.GetMoneyPos(1)}
	var result [2]int = [2]int{money[0], money[1]}

	if self.IsLucky() {
		for key, value := range self.Bets[trend] {
			key.Win = value
		}
		result[trend] -= money[trend]
	} else {
		for key, value := range self.Bets[trend] {
			winmoney := value * 2
			key.Win += winmoney
			key.Cost = int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		result[trend] -= 2 * money[trend]
	}

	if self.Kill[0] != nil { //! 杀单
		per := self.GetPerson(self.Kill[0].Uid)
		if per != nil {
			if trend == 0 { //! 开单，赔单的钱
				per.Win -= money[0]
				result[0] += money[0]
			} else if trend == 1 { //! 开双，吃单的钱
				per.Win += money[0]
				per.Cost += int(math.Ceil(float64(money[0]) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
				result[0] = 0
			}
		}
	}

	if self.Kill[1] != nil { //! 杀双
		per := self.GetPerson(self.Kill[1].Uid)
		if per != nil {
			if trend == 0 { //! 开单，吃双的钱
				per.Win += money[1]
				per.Cost += int(math.Ceil(float64(money[1]) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
				result[1] = 0
			} else if trend == 1 { //! 开双，赔双的钱
				per.Win -= money[1]
				result[1] += money[1]
			}
		}
	}
	var bigwin *Game_GoldTB_Person
	for _, value := range self.PersonMgr {
		if value.Win > 0 {
			value.Win -= value.Cost
			GetServer().SqlAgentGoldLog(value.Uid, value.Cost, self.room.Type)

			if bigwin == nil {
				bigwin = value
			} else if value.Win > bigwin.Win {
				bigwin = value
			}
		}
		value.Total += value.Win

		var msg Msg_GameGoldTB_balance
		msg.Uid = value.Uid
		msg.Total = value.Total
		msg.Win = value.Win
		find := false
		for j := 0; j < len(self.Seat); j++ {
			if self.Seat[j] == value {
				self.room.broadCastMsg("gamegoldtbbalance", &msg)
				find = true
				break
			}
		}
		if !find {
			self.room.SendMsg(value.Uid, "gamegoldtbbalance", &msg)
		}

		//! 插入战绩
		if value.Bets > 0 || (self.Kill[0] != nil && self.Kill[0].Uid == value.Uid) || (self.Kill[1] != nil && self.Kill[1].Uid == value.Uid) {
			var record Rec_TB_Info
			record.Time = time.Now().Unix()
			record.GameType = self.room.Type
			var rec Son_Rec_TB_Person
			rec.Uid = value.Uid
			rec.Name = value.Name
			rec.Head = value.Head
			rec.Score = value.Win - value.Bets
			rec.Result = self.Result
			rec.Bets = value.BetInfo
			rec.Lucky = self.Lucky
			rec.AllBets[0] = self.GetMoneyPos(0)
			rec.AllBets[1] = self.GetMoneyPos(1)

			for i := 0; i < len(self.Kill); i++ {
				if self.Kill[i] == nil {
					rec.Kill[i] = 0
				} else {
					rec.Kill[i] = self.Kill[i].Uid
				}
			}
			rec.Dealer = false
			record.Info = append(record.Info, rec)
			GetServer().InsertRecord(self.room.Type, value.Uid, lib.HF_JtoA(&record), rec.Score)
		}
	}

	dealwin := result[0] + result[1]
	if dealwin > 0 && self.Dealer != nil {
		cost := int(math.Ceil(float64(dealwin) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		dealwin -= cost
	}

	if self.Dealer != nil {
		var record Rec_TB_Info
		record.Time = time.Now().Unix()
		record.GameType = self.room.Type
		var rec Son_Rec_TB_Person
		rec.Uid = self.Dealer.Uid
		rec.Name = self.Dealer.Name
		rec.Head = self.Dealer.Head
		rec.Score = dealwin
		rec.Result = self.Result
		rec.Lucky = self.Lucky
		rec.AllBets[0] = self.GetMoneyPos(0)
		rec.AllBets[1] = self.GetMoneyPos(1)

		for i := 0; i < len(self.Kill); i++ {
			if self.Kill[i] == nil {
				rec.Kill[i] = 0
			} else {
				rec.Kill[i] = self.Kill[i].Uid
			}
		}
		rec.Dealer = true
		record.Info = append(record.Info, rec)
		GetServer().InsertRecord(self.room.Type, self.Dealer.Uid, lib.HF_JtoA(&record), rec.Score)
	}

	{
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "-----dealwin : ", dealwin)
		var msg Msg_GameGoldTB_balance
		if self.Dealer != nil {
			if self.Dealer.Total+dealwin > 0 {
				self.Dealer.Total += dealwin
			} else {
				self.Dealer.Total = 0
				dealwin = -self.Dealer.Total
			}
			msg.Uid = self.Dealer.Uid
			msg.Total = self.Dealer.Total
		}
		msg.Win = dealwin
		self.room.broadCastMsg("gamegoldtbbalance", &msg)
	}

	{
		var msg Msg_GameGoldTB_End
		msg.Result = self.Result
		if bigwin != nil {
			msg.Uid = bigwin.Uid
			msg.Name = bigwin.Name
			msg.Head = bigwin.Head
		}
		msg.Money = lib.GetManyMoneyMgr().GetProperty(self.room.Type).Money
		self.room.broadCastMsg("gamegoldtbend", &msg)
	}

	self.Total = 0
	for i := 0; i < len(self.Kill); i++ {
		self.Kill[i] = nil
	}

	//! 清理玩家
	for key, value := range self.PersonMgr {
		if value.Online {
			find := false
			for i := 0; i < len(self.LstDeal); i++ {
				if self.LstDeal[i] == value {
					find = true
					break
				}
			}

			if !find && value.Seat < 0 {
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

		//! 走的人在上庄列表上
		for j := 0; j < len(self.LstDeal); j++ {
			if self.LstDeal[j] == value {
				copy(self.LstDeal[j:], self.LstDeal[j+1:])
				self.LstDeal = self.LstDeal[:len(self.LstDeal)-1]
				break
			}
		}

		if self.Dealer == value {
			self.ChageDeal()
		}
		//! 走的人是位置上的人
		for j := 0; j < len(self.Seat); j++ {
			if self.Seat[j] == value {
				self.Seat[j] = nil
				var msg Msg_GameGoldTB_UpdSeat
				msg.Index = j
				self.room.broadCastMsg("gamegoldtbseat", &msg)
				break
			}
		}
		delete(self.PersonMgr, key)
	}

	for i := 0; i < len(self.Bets); i++ {
		self.Bets[i] = make(map[*Game_GoldTB_Person]int)
	}

	//! 判断庄家是否能继续连庄
	if self.Dealer != nil && (self.Dealer.Total < lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney || self.DownUid == self.Dealer.Uid) {
		self.ChageDeal()
	} else if self.Dealer == nil && len(self.LstDeal) > 0 {
		self.ChageDeal()
	} else if self.Dealer != nil {
		if self.Round >= 10 && len(self.LstDeal) > 0 {
			self.ChageDeal()
		} else {
			self.Round++
		}
	}

	if self.Dealer == nil {
		self.SetTime(0)
	} else {
		self.SetTime(TB_TIME)
	}

}
func (self *Game_GoldTBRoom) OnIsBets(uid int64) bool {
	value, ok := self.PersonMgr[uid]
	if ok {
		return value.Bets > 0
	}
	return false
}

func (self *Game_GoldTBRoom) OnInit(room *Room) {
	self.room = room
}

func (self *Game_GoldTBRoom) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldTBRoom) OnSendInfo(person *Person) {
	value, ok := self.PersonMgr[person.Uid]
	if ok {
		value.Online = true
		value.Round = 0
		value.IP = person.ip
		value.Address = person.minfo.Address
		value.Sex = person.Sex
		value.SynchroGold(person.Gold)
		person.SendMsg("gamegoldtbinfo", self.getInfo(value.Uid, value.Total, value.BetInfo))
		return
	}

	_person := new(Game_GoldTB_Person)
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
	person.SendMsg("gamegoldtbinfo", self.getInfo(person.Uid, person.Gold, [2]int{0, 0}))
}

func (self *Game_GoldTBRoom) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold": //! 同步金币
		person := self.GetPerson(msg.V.(*staticfunc.Msg_SynchroGold).Uid)
		if person != nil {
			person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal(person.Uid, person.Total)
		}
	case "gamebzwbets": //! 下注
		self.GameBets(msg.Uid, msg.V.(*Msg_GameGoldBZW_Bets).Index, msg.V.(*Msg_GameGoldBZW_Bets).Gold)
	case "gamebzwgoon": //! 续押
		self.GameGoOn(msg.Uid)
	case "gamerob": //! 上庄
		self.GameUpDeal(msg.Uid)
	case "gameredeal": //！ 下庄
		self.GameReDeal(msg.Uid)
	case "gamebzwseat":
		self.GameSeat(msg.Uid, msg.V.(*Msg_GameGoldBZW_Seat).Index)
	case "gameplayerlist":
		self.GamePlayerList(msg.Uid)
	case "gametbkill":
		self.GetKill(msg.Uid, msg.V.(*Msg_GameGoldBZW_Seat).Index)
	}
}

func (self *Game_GoldTBRoom) OnBye() {

}

func (self *Game_GoldTBRoom) OnExit(uid int64) {
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

func (self *Game_GoldTBRoom) OnIsDealer(uid int64) bool {
	if self.Dealer != nil && self.Dealer.Uid == uid {
		return true
	}
	return false
}

func (self *Game_GoldTBRoom) OnBalance() {
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

func (self *Game_GoldTBRoom) OnTime() {
	if self.Time == 0 {
		return
	}

	if time.Now().Unix() < self.Time {
		return
	}

	if !self.room.Begin {
		self.OnBegin()
	}
}
