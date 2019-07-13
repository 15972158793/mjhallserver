package gameserver

import (
	"fmt"
	"lib"
	"math"
	"staticfunc"
	"time"
)

type Rec_TB_Info struct {
	GameType int                 `json:"gametype"`
	Time     int64               `json:"time"`
	Info     []Son_Rec_TB_Person `json:"info"`
}
type Son_Rec_TB_Person struct {
	Uid     int64    `json:"uid"`
	Name    string   `json:"name"`
	Head    string   `json:"head"`
	Score   int      `json:"score"`
	Result  [2]int   `json:"result"`
	Bets    [2]int   `json:"bets"`
	AllBets [2]int   `json:"allbets"`
	Dealer  bool     `json:"dealer"`
	Kill    [2]int64 `json:"kill"`
	Lucky   [2]int   `json:"lucky"`
}

type Game_GoldTBSeat struct {
	Person *Game_GoldTB_Person
	Robot  *lib.Robot
}

type Game_GoldTB struct {
	PersonMgr map[int64]*Game_GoldTB_Person  `json:"personmgr"`
	Bets      [2]map[*Game_GoldTB_Person]int `json:"bets"`
	Kill      [2]*Son_GameGoldTB_Info        `json:"kill"` //! 杀的玩家 0-单 1-双
	Result    [2]int                         `json:"result"`
	Lucky     [2]int                         `json:"lucky"` //! 幸运数
	Dealer    *Game_GoldTB_Person            `json:"dealer"`
	Round     int                            `json:"round"`   //! 连庄轮数
	DownUid   int64                          `json:"downuid"` //! 下庄的人
	Time      int64                          `json:"time"`
	LstDeal   []*Game_GoldTB_Person          `json:"lstdeal"` //! 上庄列表
	Seat      [8]Game_GoldTBSeat             `json:"seat"`    //! 8个位置
	Total     int                            `json:"total"`   //! 这局一共下了多少钱
	Money     int                            `json:"money"`   //! 系统的钱
	Trend     [][2]int                       `json:"trend"`   //! 走势
	Next      []int                          `json:"next"`
	Robot     lib.ManyGameRobot              //! 机器人结构
	BetTime   int                            `json:"bettime"`

	room *Room
}

func NewGame_GoldTB() *Game_GoldTB {
	game := new(Game_GoldTB)
	game.PersonMgr = make(map[int64]*Game_GoldTB_Person)
	game.Lucky = [2]int{lib.HF_GetRandom(6) + 1, lib.HF_GetRandom(6) + 1}
	for i := 0; i < len(game.Bets); i++ {
		game.Bets[i] = make(map[*Game_GoldTB_Person]int)
	}
	for i := 0; i < 20; i++ {
		game.Trend = append(game.Trend, [2]int{lib.HF_GetRandom(6) + 1, lib.HF_GetRandom(6) + 1})
	}
	for i := 0; i < 2; i++ {
		game.Kill[i] = nil
	}
	return game
}

type Game_GoldTB_Person struct {
	Uid       int64  `json:"uid"`
	Gold      int    `json:"gold"`      //! 进来时候的钱
	Total     int    `json:"total"`     //! 当前的钱
	Win       int    `json:"win"`       //! 本局赢了多少
	Cost      int    `json:"cost"`      //! 抽水
	Bets      int    `json:"bets"`      //! 本局下了多少
	BetInfo   [2]int `json:"betinfo"`   //！ 本局下注的具体情况
	BeBets    int    `json:"bebets`     //! 上一局下了多少
	BeBetInfo [2]int `json:"bebetinfo"` //！ 上一局下注具体情况
	Name      string `json:"name"`      //! 名字
	Head      string `json:"head"`      //! 头像
	Online    bool   `json:"online"`
	Round     int    `json:"round"` //!  不下注轮数
	Seat      int    `json:"seat"`  //! 0-7有座 -1无座 100庄家
	IP        string `json:"ip"`
	Address   string `json:"address"`
	Sex       int    `json:"sex"`
}

type Msg_GameGoldTB_Info struct {
	Begin   bool                   `json:"begin"`  //! 是否开始
	Time    int64                  `json:"time"`   //! 倒计时
	Seat    [8]Son_GameGoldTB_Info `json:"seat"`   //! 8个位置
	Bets    [2]int                 `json:"bets"`   //! 2个区域的下注
	Dealer  Son_GameGoldTB_Info    `json:"dealer"` //! 庄的信息
	Total   int                    `json:"total"`  //! 自己的钱
	Trend   [][2]int               `json:"trend"`  //! 历史记录
	Kill    [2]Son_GameGoldTB_Info `json:"kill"`   //! 杀的玩家 0-单 1-双
	Lucky   [2]int                 `json:"lucky"`  //! 幸运数
	IsDeal  bool                   `json:"isdeal"` //! 是否可以下庄
	MyBets  [2]int                 `json:"mybets"`
	Money   []int                  `json:"money"`
	BetTime int                    `json:"bettime"`
}
type Son_GameGoldTB_Info struct {
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Total   int    `json:"total"`
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"` //! -1是机器人
}

type Msg_GameGoldTB_Total struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"`
}

type Msg_GameGoldTB_UpdSeat struct {
	Uid     int64  `json:"uid"`
	Index   int    `json:"index"`
	Head    string `json:"head"`
	Name    string `json:"name"`
	Total   int    `json:"total"`
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type Msg_GameGoldTB_List struct {
	Info []Son_GameGoldTB_Info `json:"info"`
}

type Msg_GameGoldTB_Bets struct {
	Uid       int64  `json:"uid"`
	Index     int    `json:"index"`
	Gold      int    `json:"gold"`
	Total     int    `json:"total"`
	BetInfo   [2]int `json:"betinfo"`
	GameTotal [2]int `json:"gametotal"`
}

type Msg_GameGoldTB_GoOn struct {
	Uid       int64  `json:"uid"`
	Gold      [2]int `json:"gold"`
	Total     int    `json:"total"`
	GameTotal [2]int `json:"gametotal"`
}

type Msg_GameGoldTB_DealList struct {
	Type int                   `json:"type"`
	Info []Son_GameGoldTB_Info `json:"info"`
}

type Msg_GameGoldTB_Kill struct {
	Uid  int64  `json:"uid"`
	Name string `json:"name"`
	Head string `json:"head"`
	Type int    `json:"type"`
}
type Msg_GameGoldTB_Deal struct {
	Uid     int64  `json:"uid"`
	Head    string `json:"head"`
	Name    string `json:"name"`
	Lucky   [2]int `json:"lucky"` //! 幸运数
	Total   int    `json:"total"`
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type Msg_GameGoldTB_balance struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"`
	Win   int   `json:"win"`
}

type Msg_GameGoldTB_End struct {
	Result  [2]int `json:"result"`
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Money   []int  `json:"money"`
	BetTime int    `json:"bettime"`
}

func (self *Game_GoldTB) getInfo(uid int64, total int, mybets [2]int) *Msg_GameGoldTB_Info {
	var msg Msg_GameGoldTB_Info
	msg.Begin = self.room.Begin
	msg.Time = self.Time - time.Now().Unix()
	msg.Total = total
	msg.Trend = self.Trend
	msg.Lucky = self.Lucky
	msg.MyBets = mybets
	msg.Money = lib.GetManyMoneyMgr().GetProperty(self.room.Type).Money
	msg.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime
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
		msg.Bets[i] = self.GetMoneyPos(i, true)
	}
	for i := 0; i < 8; i++ {
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
	if self.Dealer != nil {
		msg.Dealer.Uid = self.Dealer.Uid
		msg.Dealer.Name = self.Dealer.Name
		msg.Dealer.Head = self.Dealer.Head
		msg.Dealer.Total = self.Dealer.Total
		msg.Dealer.IP = self.Dealer.IP
		msg.Dealer.Address = self.Dealer.Address
		msg.Dealer.Sex = self.Dealer.Sex
	} else {
		msg.Dealer.Total = self.Money
	}
	return &msg
}

//! 得到这个位置下了多少钱
func (self *Game_GoldTB) GetMoneyPos(index int, robot bool) int {
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

//! 得到这个位置机器人下了多少钱
//func (self *Game_GoldTB) GetMoneyPosByRobot(index int) int {
//	total := 0
//	for _, value := range self.Robot.RobotsBet[index] {
//		total += value
//	}
//	return total
//}

func (self *Game_GoldTB) GetPerson(uid int64) *Game_GoldTB_Person {
	return self.PersonMgr[uid]
}

//! 同步金币
func (self *Game_GoldTB_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

//! 同步总分
func (self *Game_GoldTB) SendTotal(uid int64, total int) {
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
func (self *Game_GoldTB) SetTime(t int) {
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
func (self *Game_GoldTB) OnIsBets(uid int64) bool {
	value, ok := self.PersonMgr[uid]
	if ok {
		if self.Kill[0] != nil && self.Kill[0].Uid == uid {
			return true
		}
		if self.Kill[1] != nil && self.Kill[1].Uid == uid {
			return true
		}
		return value.Bets > 0
	}
	return false
}

//! 坐下
func (self *Game_GoldTB) GameSeat(uid int64, index int) {
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
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "-----------zuoxia : ", lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney)
	if person.Total < lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("金币必须大于%d才能坐下", lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("金币必须大于%d才能坐下", lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("金币必须大于%d万才能坐下", lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney/10000))
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
		self.Seat[index].Person.Seat = -1
	} else if self.Seat[index].Robot != nil {
		if person.Total <= self.Seat[index].Robot.GetMoney() {
			self.room.SendErr(uid, "该位置已经有人坐了")
			return
		}
		self.Seat[index].Robot.SetSeat(-1)
	}

	self.Seat[index].Person = person
	self.Seat[index].Robot = nil
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

//! 机器人坐下
func (self *Game_GoldTB) RobotSeat(index int, robot *lib.Robot) {
	if index < 0 || index > 7 {
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

	var msg Msg_GameGoldTB_UpdSeat
	msg.Uid = robot.Id
	msg.Index = index
	msg.Head = robot.Head
	msg.Name = robot.Name
	msg.Total = robot.GetMoney()
	msg.IP = robot.IP
	msg.Address = robot.Address
	msg.Sex = robot.Sex
	self.room.broadCastMsg("gamegoldtbseat", &msg)
}

//! 申请无座玩家
func (self *Game_GoldTB) GamePlayerList(uid int64) {
	var msg Msg_GameGoldTB_List
	tmp := make(map[int64]Son_GameGoldTB_Info)
	for _, value := range self.PersonMgr {
		if value.Seat >= 0 {
			continue
		}

		var node Son_GameGoldTB_Info
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

		var node Son_GameGoldTB_Info
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

//! 下注
func (self *Game_GoldTB) GameBets(uid int64, index int, gold int) {
	if uid == 0 {
		return
	}

	if index != 0 && index != 1 {
		return
	}

	if gold <= 0 {
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() <= 5 {
		self.room.SendErr(uid, "现在是抢杀阶段，不能下注")
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() >= int64(self.BetTime-2) {
		self.room.SendErr(uid, "正在开奖，请稍后下注")
		return
	}

	if self.Dealer == nil && lib.GetManyMgr().GetProperty(self.room.Type).SysNoBets == 1 {
		self.room.SendErr(uid, "请等待玩家上庄")
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
	if self.Dealer == nil {
		dealtotal = self.Money
	} else {
		dealtotal = self.Dealer.Total + self.Robot.RobotTotal
	}

	if self.GetMoneyPos(index, true)+gold > dealtotal {
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
		msg.GameTotal[i] = self.GetMoneyPos(i, true)
	}
	self.room.broadCastMsg("gamegoldtbbets", &msg)

	if dealtotal == self.GetMoneyPos(0, true) && dealtotal == self.GetMoneyPos(1, true) {
		self.SetTime(5)
	}
}

//! 续押
func (self *Game_GoldTB) GameGoOn(uid int64) {
	if uid == 0 {
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() <= 5 {
		self.room.SendErr(uid, "现在是抢杀阶段，不能下注")
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() >= int64(self.BetTime-2) {
		self.room.SendErr(uid, "正在开奖，请稍后下注")
		return
	}

	if self.Dealer == nil && lib.GetManyMgr().GetProperty(self.room.Type).SysNoBets == 1 {
		self.room.SendErr(uid, "请等待玩家上庄")
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
	if self.Dealer == nil {
		dealtotal = self.Money
	} else {
		dealtotal = self.Dealer.Total + self.Robot.RobotTotal
	}

	for i := 0; i < len(person.BeBetInfo); i++ {
		if self.GetMoneyPos(i, true)+person.BeBetInfo[i] > dealtotal {
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
		msg.GameTotal[i] = self.GetMoneyPos(i, true)
	}
	self.room.broadCastMsg("gamegoldtbgoon", &msg)

	if dealtotal == self.GetMoneyPos(0, true) && dealtotal == self.GetMoneyPos(1, true) {
		self.SetTime(5)
	}
}

//! 上庄
func (self *Game_GoldTB) GameUpDeal(uid int64) {
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

	var msg Msg_GameGoldTB_DealList
	msg.Type = 0
	msg.Info = make([]Son_GameGoldTB_Info, 0)
	for i := 0; i < len(self.LstDeal); i++ {
		msg.Info = append(msg.Info, Son_GameGoldTB_Info{self.LstDeal[i].Uid, self.LstDeal[i].Name, self.LstDeal[i].Head, self.LstDeal[i].Total, self.LstDeal[i].IP, self.LstDeal[i].Address, self.LstDeal[i].Sex})
	}
	self.room.SendMsg(uid, "gamegoldtbdeal", &msg)
}

//! 下庄
func (self *Game_GoldTB) GameReDeal(uid int64) {
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

//! 抢杀
func (self *Game_GoldTB) GetKill(uid int64, index int) {
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

	if (index == 0 && person.Total < self.GetMoneyPos(0, true)) || (index == 1 && person.Total < self.GetMoneyPos(1, true)) {
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

//! 换庄
func (self *Game_GoldTB) ChageDeal() {
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
				if self.Seat[i].Person == self.Dealer {
					var msg Msg_GameGoldTB_UpdSeat
					msg.Index = i
					self.room.broadCastMsg("gamegoldtbseat", &msg)
					self.Seat[i].Person = nil
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
	} else {
		msg.Total = self.Money
		msg.Lucky = self.Lucky
	}
	self.room.broadCastMsg("gamerob", &msg)
}

func (self *Game_GoldTB) IsType() int { //! 0-单 1-双
	if (self.Result[0]+self.Result[1])%2 == 0 {
		return 1
	}
	return 0
}

func (self *Game_GoldTB) IsTypePlus() int { //! 0-单 1-双 2-幸运数字
	if self.IsLucky() {
		return 2
	}

	if (self.Result[0]+self.Result[1])%2 == 0 {
		return 1
	}
	return 0
}

//! 是否是幸运数
func (self *Game_GoldTB) IsLucky() bool {
	return (self.Result[0] == self.Lucky[0] && self.Result[1] == self.Lucky[1]) || (self.Result[0] == self.Lucky[1] && self.Result[1] == self.Lucky[0])
}

//! 判断幸运数字类型
func (self *Game_GoldTB) LuckyType() int {
	if (self.Lucky[0]+self.Lucky[1])%2 == 0 {
		return 1
	}
	return 0
}

//! 3种情况下机器人赢钱
func (self *Game_GoldTB) GetRobotWin(_type int /*0单,1双,2幸运数字*/) int {
	total := 0
	if _type == 2 { //! 开幸运数字
		if self.LuckyType() == 0 { //! 幸运数字是单
			if self.Dealer != nil || (self.Kill[1] != nil && self.Kill[1].Sex != -1) { //! 玩家庄或者杀双
				for _, value := range self.Robot.RobotsBet[1] {
					total -= value
				}
			}
		} else { //! 幸运数字是双
			if self.Dealer != nil || (self.Kill[0] != nil && self.Kill[0].Sex != -1) { //! 玩家庄或者杀单
				for _, value := range self.Robot.RobotsBet[0] {
					total -= value
				}
			}
		}
	} else if _type == 0 { //! 开单
		if self.Dealer != nil || (self.Kill[0] != nil && self.Kill[0].Sex != -1) {
			for _, value := range self.Robot.RobotsBet[0] {
				total += value
			}
		}
		if self.Dealer != nil || (self.Kill[1] != nil && self.Kill[1].Sex != -1) {
			for _, value := range self.Robot.RobotsBet[1] {
				total -= value
			}
		}
	} else if _type == 1 { //! 开双
		if self.Dealer != nil || (self.Kill[0] != nil && self.Kill[0].Sex != -1) {
			for _, value := range self.Robot.RobotsBet[0] {
				total -= value
			}
		}
		if self.Dealer != nil || (self.Kill[1] != nil && self.Kill[1].Sex != -1) {
			for _, value := range self.Robot.RobotsBet[1] {
				total += value
			}
		}
	}

	//! 机器人杀需要完善
	//if self.Kill[0] != nil && self.Kill[0].Sex == -1 { //! 机器人杀单
	//	if _type == 0 {
	//		total -= self.GetMoneyPos(0, false)
	//	} else if _type == 1 {
	//		total += self.GetMoneyPos(0, false)
	//	} else if _type == 2 {

	//	}
	//}
	//if self.Kill[1] != nil && self.Kill[1].Sex == -1 { //! 机器人杀双
	//	if _type == 0 {
	//		total += self.GetMoneyPos(1, false)
	//	} else {
	//		total -= self.GetMoneyPos(1, false)
	//	}
	//}
	return total
}

//! 3种情况下庄家赢钱
func (self *Game_GoldTB) GetResultWin(_type int, robot bool) int {
	dan := self.GetMoneyPos(0, robot)
	shuang := self.GetMoneyPos(1, robot)
	danresult := dan
	shuangresult := shuang
	if _type == 0 { //! 开单
		danresult -= 2 * dan
	} else if _type == 1 { //! 开双
		shuangresult -= 2 * shuang
	} else if _type == 2 { //! 幸运数字
		if self.LuckyType() == 0 {
			danresult -= dan
		} else {
			shuangresult -= shuang
		}
	}

	if self.Kill[0] != nil && (robot || self.Kill[0].Sex != -1) { //! 杀单
		if _type == 0 {
			danresult += dan
		} else {
			danresult = 0
		}
	}
	if self.Kill[1] != nil && (robot || self.Kill[1].Sex != -1) { //! 杀双
		if _type == 1 {
			shuangresult += shuang
		} else {
			shuangresult = 0
		}
	}
	return danresult + shuangresult
}

//! 庄家必赢或必输
func (self *Game_GoldTB) WinOrLost(iswin bool, robot bool) {
	lucky := false //! 幸运数字是否满足条件

	win := self.GetResultWin(2, robot)
	if (iswin && win >= 0) || (!iswin && win <= 0) { //! 满足条件
		lucky = true
	}

	lst := make([]int, 0)
	for i := 0; i <= 1; i++ {
		win = self.GetResultWin(i, robot)
		if (iswin && win >= 0) || (!iswin && win <= 0) { //! 满足条件
			lst = append(lst, i)
		}
	}

	if len(lst) == 0 && lucky { //! 只有幸运数字满足条件
		self.SetTrendPlus(2)
	} else {
		if lucky && lib.HF_GetRandom(180) < 10 {
			self.SetTrendPlus(2)
			return
		}
		if len(lst) > 0 {
			self.SetTrendPlus(lst[lib.HF_GetRandom(len(lst))])
		}
	}
}

//! 指定开什么类型
func (self *Game_GoldTB) SetTrend(trend int) {
	for i := 0; i < 10000000; i++ {
		if self.IsType() == trend {
			break
		}

		self.Result[0] = lib.HF_GetRandom(6) + 1
		self.Result[1] = lib.HF_GetRandom(6) + 1
	}
}

//! 指定开什么类型
func (self *Game_GoldTB) SetTrendPlus(_type int) {
	if _type == 2 {
		if lib.HF_GetRandom(100) < 50 {
			self.Result[0] = self.Lucky[0]
			self.Result[1] = self.Lucky[1]
		} else {
			self.Result[0] = self.Lucky[1]
			self.Result[1] = self.Lucky[0]
		}
	} else {
		for i := 0; i < 10000000; i++ {
			if self.IsType() == _type && !self.IsLucky() {
				break
			}

			self.Result[0] = lib.HF_GetRandom(6) + 1
			self.Result[1] = lib.HF_GetRandom(6) + 1
		}
	}
}

func (self *Game_GoldTB) OnBegin() {
	if self.room.IsBye() {
		return
	}
	self.room.Begin = true

	self.Result[0] = lib.HF_GetRandom(6) + 1
	self.Result[1] = lib.HF_GetRandom(6) + 1

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "next:", self.Next)
	if len(self.Next) == 1 {
		if self.Next[0] == 0 { //! 单
			self.SetTrend(0)
		} else { //! 双
			self.SetTrend(1)
		}

		self.Next = make([]int, 0)
		self.OnEnd()
		return
	} else if len(self.Next) == 2 {
		if self.Next[0] >= 1 && self.Next[0] <= 6 {
			self.Result[0] = self.Next[0]
		}
		if self.Next[1] >= 1 && self.Next[1] <= 6 {
			self.Result[1] = self.Next[1]
		}

		self.Next = make([]int, 0)
		self.OnEnd()
		return
	}

	if self.Dealer != nil { //! 玩家庄
		if self.Robot.RobotTotal == 0 { //! 没有机器人下注
			if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost == 102 { //! 平衡模式
				win := self.GetResultWin(self.IsTypePlus(), false)
				if GetServer().TBUserMoney[self.room.Type%140000]+int64(win) > lib.GetManyMgr().GetProperty(self.room.Type).PlayerMax { //! 玩家必须输
					self.WinOrLost(false, false)
				} else if GetServer().TBUserMoney[self.room.Type%140000]+int64(win) < lib.GetManyMgr().GetProperty(self.room.Type).PlayerMin { //! 玩家必须赢
					self.WinOrLost(true, false)
				}
			} else if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost == 101 { //! 随机模式

			} else if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost <= 100 { //! 概率模式
				if lib.HF_GetRandom(100) < lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost { //! 玩家赢
					self.WinOrLost(true, false)
				} else { //! 玩家输
					self.WinOrLost(false, false)
				}
			}
		} else {
			lucky := false
			win := self.GetRobotWin(2)
			if lib.GetRobotMgr().GetRobotWin(self.room.Type)+win >= 0 {
				lucky = true
			} else if win >= 0 {
				lucky = true
			}

			lst := make([]int, 0)
			for i := 0; i <= 1; i++ {
				win = self.GetRobotWin(i)
				if lib.GetRobotMgr().GetRobotWin(self.room.Type)+win >= 0 {
					lst = append(lst, i)
				} else if win >= 0 {
					lst = append(lst, i)
				}
			}

			if len(lst) == 0 && lucky { //! 只有幸运数字满足条件
				self.SetTrendPlus(2)
			} else {
				if lucky && lib.HF_GetRandom(180) < 10 {
					self.SetTrendPlus(2)
				} else if len(lst) > 0 {
					self.SetTrendPlus(lst[lib.HF_GetRandom(len(lst))])
				}
			}
		}
	} else { //! 系统庄
		var dealwin [3]int
		var robotwin [3]int
		lucky := false
		lst := make([]int, 0)
		for i := 0; i < 3; i++ {
			dealwin[i] = self.GetResultWin(i, false)
			robotwin[i] = self.GetRobotWin(i)
			if GetServer().TBMoney[self.room.Type%140000] < lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
				if dealwin[i] >= 0 && (robotwin[i] >= 0 || lib.GetRobotMgr().GetRobotWin(self.room.Type)+robotwin[i] >= 0) {
					if i == 2 {
						lucky = true
					} else {
						lst = append(lst, i)
					}
				}
			} else if GetServer().TBMoney[self.room.Type%140000] > lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
				if dealwin[i] <= 0 && (robotwin[i] >= 0 || lib.GetRobotMgr().GetRobotWin(self.room.Type)+robotwin[i] >= 0) {
					if i == 2 {
						lucky = true
					} else {
						lst = append(lst, i)
					}
				}
			} else if robotwin[i] >= 0 || lib.GetRobotMgr().GetRobotWin(self.room.Type)+robotwin[i] >= 0 {
				if i == 2 {
					lucky = true
				} else {
					lst = append(lst, i)
				}
			}
		}
		if len(lst) == 0 && lucky { //! 只有幸运数字满足条件
			self.SetTrendPlus(2)
		} else {
			if lucky && lib.HF_GetRandom(180) < 10 {
				self.SetTrendPlus(2)
			} else if len(lst) > 0 {
				self.SetTrendPlus(lst[lib.HF_GetRandom(len(lst))])
			}
		}
	}

	self.OnEnd()
}

func (self *Game_GoldTB) OnEnd() {
	self.room.Begin = false

	tmp := [][2]int{self.Result}
	tmp = append(tmp, self.Trend...)
	if len(tmp) > 20 {
		tmp = tmp[0:20]
	}
	self.Trend = tmp
	trend := self.IsType()

	var money [2]int = [2]int{self.GetMoneyPos(0, true), self.GetMoneyPos(1, true)}

	if self.IsLucky() {
		for key, value := range self.Bets[trend] {
			key.Win = value
		}
		for key, value := range self.Robot.RobotsBet[trend] {
			key.AddWin(value)
		}
	} else {
		for key, value := range self.Bets[trend] {
			winmoney := value * 2
			key.Win += winmoney
			key.Cost = int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[trend] {
			winmoney := value * 2
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
		}
	}

	if self.Kill[0] != nil && self.Kill[0].Sex != -1 { //! 杀单
		per := self.GetPerson(self.Kill[0].Uid)
		if per != nil {
			if trend == 0 { //! 开单，赔单的钱
				per.Win -= money[0]
			} else if trend == 1 { //! 开双，吃单的钱
				per.Win += money[0]
				per.Cost += int(math.Ceil(float64(money[0]) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
			}
		}
	}

	if self.Kill[1] != nil && self.Kill[1].Sex != -1 { //! 杀双
		per := self.GetPerson(self.Kill[1].Uid)
		if per != nil {
			if trend == 0 { //! 开单，吃双的钱
				per.Win += money[1]
				per.Cost += int(math.Ceil(float64(money[1]) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
			} else if trend == 1 { //! 开双，赔双的钱
				per.Win -= money[1]
			}
		}
	}
	var bigwin *GameGold_BigWin
	for _, value := range self.PersonMgr {
		if value.Win > 0 {
			value.Win -= value.Cost
			GetServer().SqlAgentGoldLog(value.Uid, value.Cost, self.room.Type)
			GetServer().SqlAgentBillsLog(value.Uid, value.Cost/2, self.room.Type)

			if bigwin == nil {
				bigwin = &GameGold_BigWin{value.Uid, value.Name, value.Head, value.Win}
			} else if value.Win > bigwin.Win {
				bigwin = &GameGold_BigWin{value.Uid, value.Name, value.Head, value.Win}
			}
		} else if value.Win-value.Bets < 0 {
			cost := int(math.Ceil(float64(value.Win-value.Bets) * float64(lib.GetManyMgr().GetProperty(self.room.Type).Cost) / 200.0))
			GetServer().SqlAgentBillsLog(value.Uid, cost, self.room.Type)
		}
		value.Total += value.Win

		var msg Msg_GameGoldTB_balance
		msg.Uid = value.Uid
		msg.Total = value.Total
		msg.Win = value.Win
		find := false
		for j := 0; j < len(self.Seat); j++ {
			if self.Seat[j].Person == value {
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
			rec.AllBets[0] = self.GetMoneyPos(0, true)
			rec.AllBets[1] = self.GetMoneyPos(1, true)

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

	for i := 0; i < len(self.Robot.Robots); i++ {
		if self.Robot.Robots[i].GetWin() > 0 {
			self.Robot.Robots[i].AddWin(-self.Robot.Robots[i].GetCost())
			self.Robot.Robots[i].AddMoney(self.Robot.Robots[i].GetWin())

			for j := 0; j < len(self.Seat); j++ {
				if self.Seat[j].Robot == self.Robot.Robots[i] {
					var msg Msg_GameGoldTB_balance
					msg.Uid = self.Robot.Robots[i].Id
					msg.Total = self.Robot.Robots[i].GetMoney()
					msg.Win = self.Robot.Robots[i].GetWin()
					self.room.broadCastMsg("gamegoldtbbalance", &msg)
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

	robotwin := self.GetRobotWin(self.IsTypePlus())
	if robotwin != 0 {
		lib.GetRobotMgr().AddRobotWin(self.room.Type, robotwin)
		GetServer().SqlBZWLog(&SQL_BZWLog{1, robotwin, time.Now().Unix(), self.room.Type + 10000000})
	}

	dealwin1 := self.GetResultWin(self.IsTypePlus(), true)  //! 带机器人结果
	dealwin2 := self.GetResultWin(self.IsTypePlus(), false) //! 不带机器人结果
	if self.Dealer == nil && dealwin2 != 0 {
		GetServer().SqlBZWLog(&SQL_BZWLog{1, dealwin2, time.Now().Unix(), self.room.Type})
	}
	if self.Dealer == nil && dealwin2 > 0 {
		cost := int(math.Ceil(float64(dealwin2) * float64(lib.GetManyMgr().GetProperty(self.room.Type).DealCost) / 100.0))
		dealwin2 -= cost
	} else if self.Dealer != nil {
		if dealwin1 > 0 {
			cost := int(math.Ceil(float64(dealwin1) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
			dealwin1 -= cost
			GetServer().SqlAgentGoldLog(self.Dealer.Uid, cost, self.room.Type)
			GetServer().SqlAgentBillsLog(self.Dealer.Uid, cost/2, self.room.Type)
		} else if dealwin1 < 0 && self.Dealer != nil {
			cost := int(math.Ceil(float64(dealwin1) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 200.0))
			GetServer().SqlAgentBillsLog(self.Dealer.Uid, cost, self.room.Type)
		}
	}

	lib.GetLogMgr().Output(lib.LOG_ERROR, "RobotTotal=", self.Robot.RobotTotal)
	lib.GetLogMgr().Output(lib.LOG_ERROR, "robotwin=", robotwin)
	lib.GetLogMgr().Output(lib.LOG_ERROR, "dealwin1=", dealwin1)
	lib.GetLogMgr().Output(lib.LOG_ERROR, "dealwin2=", dealwin2)

	if self.Dealer != nil {
		var record Rec_TB_Info
		record.Time = time.Now().Unix()
		record.GameType = self.room.Type
		var rec Son_Rec_TB_Person
		rec.Uid = self.Dealer.Uid
		rec.Name = self.Dealer.Name
		rec.Head = self.Dealer.Head
		rec.Score = dealwin1
		rec.Result = self.Result
		rec.Lucky = self.Lucky
		rec.AllBets[0] = self.GetMoneyPos(0, true)
		rec.AllBets[1] = self.GetMoneyPos(1, true)

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

	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 16
	self.SetTime(self.BetTime)
	{
		self.Money = lib.GetManyMgr().GetProperty(self.room.Type).DealInitMoney
		var msg Msg_GameGoldTB_balance
		if self.Dealer != nil {
			if self.Dealer.Total+dealwin1 > 0 {
				self.Dealer.Total += dealwin1
			} else {
				self.Dealer.Total = 0
				dealwin1 = -self.Dealer.Total
			}
			msg.Uid = self.Dealer.Uid
			msg.Total = self.Dealer.Total
			GetServer().SetTBUserMoney(self.room.Type%140000, GetServer().TBUserMoney[self.room.Type%140000]+int64(dealwin1))
		} else {
			GetServer().SetTBMoney(self.room.Type%140000, GetServer().TBMoney[self.room.Type%140000]+int64(dealwin2))
			if lib.GetManyMgr().GetProperty(self.room.Type).DealChange == 1 {
				self.Money += int(GetServer().TBMoney[self.room.Type%140000])
			}
			msg.Total = self.Money
		}
		msg.Win = dealwin1
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
		msg.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime
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
			if self.Seat[j].Person == value {
				self.Seat[j].Person = nil
				var msg Msg_GameGoldTB_UpdSeat
				msg.Index = j
				self.room.broadCastMsg("gamegoldtbseat", &msg)
				break
			}
		}
		delete(self.PersonMgr, key)
	}

	//! 返回机器人
	self.Robot.Init(2, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)

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

	//! 判断坐下的人是否能继续坐
	for i := 0; i < len(self.Seat); i++ {
		if self.Seat[i].Person == nil {
			continue
		}
		if self.Seat[i].Person.Total < lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney {
			self.Seat[i].Person.Seat = -1
			var msg Msg_GameGoldTB_UpdSeat
			msg.Index = i
			self.room.broadCastMsg("gamegoldtbseat", &msg)
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
			var msg Msg_GameGoldTB_UpdSeat
			msg.Index = i
			self.room.broadCastMsg("gamegoldtbseat", &msg)
			self.Seat[i].Robot = nil
		}
	}
}

func (self *Game_GoldTB) OnInit(room *Room) {
	self.room = room
	self.Money = lib.GetManyMgr().GetProperty(self.room.Type).DealInitMoney
	if lib.GetManyMgr().GetProperty(self.room.Type).DealChange == 1 {
		self.Money += int(GetServer().TBMoney[self.room.Type%140000])
	}
	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 16
	//! 载入机器人
	self.Robot.Init(2, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)
}

func (self *Game_GoldTB) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldTB) OnSendInfo(person *Person) {
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

func (self *Game_GoldTB) OnMsg(msg *RoomMsg) {
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
	case "gamesetnext":
		self.Next = msg.V.(*staticfunc.Msg_SetDealNext).Next
	}
}

func (self *Game_GoldTB) OnBye() {

}

func (self *Game_GoldTB) OnExit(uid int64) {
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

func (self *Game_GoldTB) OnIsDealer(uid int64) bool {
	if self.Dealer != nil && self.Dealer.Uid == uid {
		return true
	}
	return false
}

func (self *Game_GoldTB) OnBalance() {
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

func (self *Game_GoldTB) OnTime() {
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
		if self.Dealer == nil && lib.GetManyMgr().GetProperty(self.room.Type).SysNoBets == 1 {
			return
		}

		if self.Time-time.Now().Unix() >= int64(lib.GetManyMgr().GetProperty(self.room.Type).BetTime) || self.Time-time.Now().Unix() <= 5 {
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
			if self.Dealer != nil { //! 玩家庄判断是否能下
				if self.GetMoneyPos(index, true)+gold > self.Dealer.Total {
					self.Robot.GameBackBets(self.Robot.Robots[i], index, gold)
					continue
				}
			}
			var msg Msg_GameGoldTB_Bets
			msg.Uid = self.Robot.Robots[i].Id
			msg.Index = index
			msg.Gold = gold
			msg.Total = self.Robot.Robots[i].GetMoney()
			for i := 0; i < len(self.Bets); i++ {
				msg.GameTotal[i] = self.GetMoneyPos(i, true)
			}
			self.room.broadCastMsg("gamegoldtbbets", &msg)
		}
		return
	}
	if !self.room.Begin {
		self.OnBegin()
	}
}
