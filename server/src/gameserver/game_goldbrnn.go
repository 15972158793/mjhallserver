package gameserver

import (
	"fmt"
	"lib"
	"math"
	//	"sort"
	"staticfunc"
	"time"
)

type Rec_BrNN_Info struct {
	GameType int                   `json:"gametype"`
	Time     int64                 `json:"time"`
	Info     []Son_Rec_BrNN_Person `json:"info"`
}

type Son_Rec_BrNN_Person struct {
	Uid    int64    `json:"uid"`
	Name   string   `json:"name"`
	Head   string   `json:"head"`
	Score  int      `json:"score"`
	Result [5][]int `json:"result"`
	Bets   [4]int   `json:"bets"`
}

type Game_GoldBrNNSeat struct {
	Person *Game_GoldBrNN_Person
	Robot  *lib.Robot
}

func (self *Game_GoldBrNNSeat) GetTotal() int {
	if self.Person != nil {
		return self.Person.Total
	} else if self.Robot != nil {
		return self.Robot.GetMoney()
	}
	return 0
}

type Game_GoldBrNN struct {
	PersonMgr   map[int64]*Game_GoldBrNN_Person  `json:"personmgr"`
	Bets        [4]map[*Game_GoldBrNN_Person]int `json:"bets"`        //! 下注详情
	Result      [5][]int                         `json:"result"`      //! 结果
	Dealer      *Game_GoldBrNN_Person            `json:"dealer"`      //! 庄
	RobotDealer *lib.Robot                       `json:"robotdealer"` //! 机器人庄
	LstDeal     []Game_GoldBrNNSeat              `json:"lstdeal"`     //! 上庄列表
	Round       int                              `json:"round"`       //! 连庄轮数
	DownUid     int64                            `json:"downuid"`     //! 下庄uid
	Time        int64                            `json:"time"`        //! 倒计时
	Seat        [6]Game_GoldBrNNSeat             `json:"seat"`        //! 座位信息
	Total       int                              `json:"total"`       //! 本局总下注
	Trend       [4][]int                         `json:"trend"`       //! 走势
	Money       int                              `json:"money"`       //! 系统庄的钱
	Robot       lib.ManyGameRobot                //! 机器人结构
	BetTime     int                              `json:"bettime"`

	room *Room
}

func NewGame_GoldBrNN() *Game_GoldBrNN {
	game := new(Game_GoldBrNN)
	game.PersonMgr = make(map[int64]*Game_GoldBrNN_Person)
	for i := 0; i < len(game.Bets); i++ {
		game.Bets[i] = make(map[*Game_GoldBrNN_Person]int)
	}
	for i := 0; i < 10; i++ {
		game.Trend[0] = append(game.Trend[0], lib.HF_GetRandom(2))
		game.Trend[1] = append(game.Trend[1], lib.HF_GetRandom(2))
		game.Trend[2] = append(game.Trend[2], lib.HF_GetRandom(2))
		game.Trend[3] = append(game.Trend[3], lib.HF_GetRandom(2))
	}

	return game
}

type Game_GoldBrNN_Person struct {
	Uid       int64  `json:"uid"`
	Gold      int    `json:"gold"`      //! 进来时有多少金币
	Total     int    `json:"total"`     //! 当前的钱
	Win       int    `json:"win"`       //! 赢了多少
	Cost      int    `json:"cost"`      //! 手续费
	Bets      int    `json:"bets"`      //! 本局下了多少
	BetInfo   [4]int `json:"betinfo"`   //! 下注详情
	BeBets    int    `json:"bebets"`    //! 上把下了多少
	BeBetInfo [4]int `json:"bebetinfo"` //! 上把下注详情
	Name      string `json:"name"`      //! 名字
	Head      string `json:"head"`
	Online    bool   `json:"online"`
	Round     int    `json:"round"` //! 不下注轮数
	Seat      int    `json:"seat"`  //! 0-6有座位 -1无座 100庄家
	IP        string `json:"ip"`
	Address   string `json:"address"`
	Sex       int    `json:"sex"`
}

type Msg_GameGoldBrNN_Info struct {
	Begin   bool                     `json:"begin"`  //! 是否开始游戏
	Time    int64                    `json:"time"`   //! 倒计时
	Seat    [6]Son_GameGoldBrNN_Info `json:"seat"`   //! 座位信息
	Bets    [4]int                   `json:"bets"`   //! 4个区域下注
	Dealer  Son_GameGoldBrNN_Info    `json:"dealer"` //! 庄的信息
	Total   int                      `json:"total"`  //! 自己的钱
	Trend   [4][]int                 `json:"trend"`  //! 走势
	IsDeal  bool                     `json:"isdeal"` //! 是否可以下庄
	Result  [5][]int                 `json:"result"` //! 场上的牌
	Money   []int                    `json:"money"`
	BetTime int                      `json:"bettime"`
}

type Msg_GameGoldBrNN_UpdSeat struct {
	Index   int    `json:"index"` //! 座位的下标
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Total   int    `json:"total"` //　玩家金币
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type Msg_GameGoldBrNN_DealList struct {
	Type int                     `json:"type"` //! 0-上庄 1-下庄
	Info []Son_GameGoldBrNN_Info `json:"info"`
}

type Msg_GameGoldBrNN_List struct {
	Info []Son_GameGoldBrNN_Info `json:"info"`
}

type Son_GameGoldBrNN_Info struct {
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Total   int    `json:"total"` //! 金币数量
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type Msg_GameGoldBrNN_Bets struct {
	Uid   int64 `json:"uid"`
	Index int   `json:"index"` //! 下注位置
	Gold  int   `json:"gold"`  //! 下了多少钱
	Total int   `json:"total"` //! 还剩多少钱
}

type BrNN_Lst struct {
	Reuslt [5][]int `json:"result"`
}

type Msg_GameGoldBrNN_Balance struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"` //! 总金币
	Win   int   `json:"win"`   //!　赢了多少
}

type Msg_GameGoldBrnn_End struct {
	Uid     int64    `json:"uid"`    //! 大赢家uid
	Name    string   `json:"name"`   //! 大赢家名字
	Head    string   `json:"head"`   //! 头像
	Result  [5][]int `json:"result"` //!　这一局开的什么牌
	Kill    int      `json:"kill"`   //! 0-非通杀通赔 1-通杀 2-通赔
	CT      []int    `json:"ct"`
	Next    [5][]int `json:"next"` //! 下一局会开什么3张名牌2张暗牌
	Money   []int    `json:"money"`
	BetTime int      `json:"bettime"`
}

type Msg_GameGoldBrNN_Deal struct {
	Uid     int64  `json:"uid"`
	Head    string `json:"head"`
	Name    string `json:"name"`
	Total   int    `json:"total"`
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

func (self *Game_GoldBrNN) getInfo(uid int64, total int) *Msg_GameGoldBrNN_Info {
	var msg Msg_GameGoldBrNN_Info
	msg.Begin = self.room.Begin
	msg.Time = self.Time - time.Now().Unix()
	msg.Total = total
	msg.Trend = self.Trend
	msg.IsDeal = false
	msg.Money = lib.GetManyMoneyMgr().GetProperty(self.room.Type).Money
	msg.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime
	lib.HF_DeepCopy(&msg.Result, &self.Result)
	for i := 0; i < len(self.Result); i++ {
		msg.Result[i][3] = 0
		msg.Result[i][4] = 0
	}
	if self.Dealer != nil && self.Dealer.Uid == uid {
		msg.IsDeal = true
	} else {
		for i := 0; i < len(self.LstDeal); i++ {
			if self.LstDeal[i].Person != nil && self.LstDeal[i].Person.Uid == uid {
				msg.IsDeal = true
				break
			}
		}
	}
	for i := 0; i < len(self.Bets); i++ {
		msg.Bets[i] = self.GetMoneyPos(i, false)
	}
	if self.Dealer == nil { //! 系统庄的时候上面计算了机器人的下注
		for i := 0; i < len(self.Bets); i++ {
			msg.Bets[i] += self.GetMoneyPosByRobot(i)
		}
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
	if self.Dealer != nil {
		msg.Dealer.Uid = self.Dealer.Uid
		msg.Dealer.Name = self.Dealer.Name
		msg.Dealer.Head = self.Dealer.Head
		msg.Dealer.Total = self.Dealer.Total
		msg.Dealer.IP = self.Dealer.IP
		msg.Dealer.Address = self.Dealer.Address
		msg.Dealer.Sex = self.Dealer.Sex
	} else if self.RobotDealer != nil {
		msg.Dealer.Uid = self.RobotDealer.Id
		msg.Dealer.Name = self.RobotDealer.Name
		msg.Dealer.Head = self.RobotDealer.Head
		msg.Dealer.Total = self.RobotDealer.GetMoney()
		msg.Dealer.IP = self.RobotDealer.IP
		msg.Dealer.Address = self.RobotDealer.Address
		msg.Dealer.Sex = self.RobotDealer.Sex
	} else {
		msg.Dealer.Total = self.Money
	}

	return &msg
}

//! 发牌
func (self *Game_GoldBrNN) GameCard() {
	card := NewCard_NiuNiu(false)
	for i := 0; i < len(self.Result); i++ {
		self.Result[i] = card.Deal(5)
	}
}

//! 同步金币
func (self *Game_GoldBrNN_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

//! 得到这个位置下了多少钱
func (self *Game_GoldBrNN) GetMoneyPos(index int, robot bool) int {
	total := 0
	for _, value := range self.Bets[index] {
		total += value
	}
	if robot || self.Dealer != nil { //! 是玩家庄,判断机器人下注
		for _, value := range self.Robot.RobotsBet[index] {
			total += value
		}
	}
	return total
}

//! 得到这个位置机器人下了多少钱
func (self *Game_GoldBrNN) GetMoneyPosByRobot(index int) int {
	total := 0
	for _, value := range self.Robot.RobotsBet[index] {
		total += value
	}
	return total
}

func (self *Game_GoldBrNN) GetPerson(uid int64) *Game_GoldBrNN_Person {
	return self.PersonMgr[uid]
}

//! 同步总分
func (self *Game_GoldBrNN) SendTotal(uid int64, total int) {
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
func (self *Game_GoldBrNN) SetTime(t int) {
	if t == 0 {
		self.Time = 0
	} else {
		self.Time = time.Now().Unix() + int64(t)
	}
	var msg Msg_SetTime
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	self.room.broadCastMsg("gametime", &msg)
}

//! 坐下
func (self *Game_GoldBrNN) GameSeat(uid int64, index int) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if index < 0 || index >= len(self.Seat) {
		return
	}

	if self.Dealer == person {
		self.room.SendErr(uid, "庄家无法坐下")
		return
	}

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

	var msg Msg_GameGoldBrNN_UpdSeat
	msg.Uid = uid
	msg.Total = person.Total
	msg.Sex = person.Sex
	msg.Name = person.Name
	msg.IP = person.IP
	msg.Index = person.Seat
	msg.Head = person.Head
	msg.Address = person.Address
	self.room.broadCastMsg("gamegoldbrnnseat", &msg)
}

//! 机器人坐下
func (self *Game_GoldBrNN) RobotSeat(index int, robot *lib.Robot) {
	if index < 0 || index >= len(self.Seat) {
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

	var msg Msg_GameGoldBrNN_UpdSeat
	msg.Uid = robot.Id
	msg.Index = index
	msg.Head = robot.Head
	msg.Name = robot.Name
	msg.Total = robot.GetMoney()
	msg.IP = robot.IP
	msg.Address = robot.Address
	msg.Sex = robot.Sex
	self.room.broadCastMsg("gamegoldbrnnseat", &msg)
}

//! 上庄
func (self *Game_GoldBrNN) GameUpDeal(uid int64) {
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
			if self.LstDeal[i].Person == person {
				self.room.SendErr(uid, "您已经在上庄列表中，请等待上庄")
				return
			}
		}
		if len(self.LstDeal) == 0 {
			self.Round = 0
		}
		self.LstDeal = append(self.LstDeal, Game_GoldBrNNSeat{person, nil})
	}

	person.Round = 0

	var msg Msg_GameGoldBrNN_DealList
	msg.Type = 0
	msg.Info = make([]Son_GameGoldBrNN_Info, 0)
	for i := 0; i < len(self.LstDeal); i++ {
		if self.LstDeal[i].Person != nil {
			msg.Info = append(msg.Info, Son_GameGoldBrNN_Info{self.LstDeal[i].Person.Uid, self.LstDeal[i].Person.Name, self.LstDeal[i].Person.Head, self.LstDeal[i].Person.Total, self.LstDeal[i].Person.IP, self.LstDeal[i].Person.Address, self.LstDeal[i].Person.Sex})
		} else if self.LstDeal[i].Robot != nil {
			msg.Info = append(msg.Info, Son_GameGoldBrNN_Info{self.LstDeal[i].Robot.Id, self.LstDeal[i].Robot.Name, self.LstDeal[i].Robot.Head, self.LstDeal[i].Robot.GetMoney(), self.LstDeal[i].Robot.IP, self.LstDeal[i].Robot.Address, self.LstDeal[i].Robot.Sex})
		}
	}
	self.room.SendMsg(uid, "gamegoldbrnndeal", &msg)
}

//! 机器人上庄
func (self *Game_GoldBrNN) RobotUpDeal(robot *lib.Robot) {
	if robot.GetMoney() < lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney {
		return
	}

	if robot.GetSeat() == 100 {
		return
	}

	for i := 0; i < len(self.LstDeal); i++ {
		if self.LstDeal[i].Robot == robot {
			return
		}
	}

	if len(self.LstDeal) == 0 {
		self.Round = 0
	}
	self.LstDeal = append(self.LstDeal, Game_GoldBrNNSeat{nil, robot})
}

//! 下庄
func (self *Game_GoldBrNN) GameReDeal(uid int64) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if self.Dealer == person {
		self.DownUid = uid
		self.room.SendErr(uid, "您已成功下庄，请等待本局结束")
	} else {
		for i := 0; i < len(self.LstDeal); i++ {
			if self.LstDeal[i].Person == person {
				copy(self.LstDeal[i:], self.LstDeal[i+1:])
				self.LstDeal = self.LstDeal[:len(self.LstDeal)-1]
				break
			}
		}
	}

	var msg Msg_GameGoldBrNN_DealList
	msg.Type = 1
	msg.Info = make([]Son_GameGoldBrNN_Info, 0)
	self.room.SendMsg(uid, "gamegoldbrnndeal", &msg)

}

//! 申请无座玩家
func (self *Game_GoldBrNN) GamePlayerList(uid int64) {
	var msg Msg_GameGoldBrNN_List
	tmp := make(map[int64]Son_GameGoldBrNN_Info)
	for _, value := range self.PersonMgr {
		if value.Seat >= 0 {
			continue
		}

		var node Son_GameGoldBrNN_Info
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

		var node Son_GameGoldBrNN_Info
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

//! 获得庄家最多可能输多少
func (self *Game_GoldBrNN) GetMaxLost(robot bool) int {
	total := 0
	if self.room.Type%230000 == 0 { //! 低倍
		for i := 0; i < 4; i++ { //! 检查是否可能开五花牛，如果不能降到四炸
			b1 := 6
			for j := 0; j < len(self.Result[i]); j++ {
				if self.Result[i][j] == 0 {
					break
				}
				if self.Result[i][j]/10 <= 10 {
					b1 = 5
				}
			}

			if b1 == 5 { //! 检查是否可能开四炸
				for i := 0; i < 4; i++ {
					if !(self.Result[i][0]/10 == self.Result[i][1]/10 || self.Result[i][0]/10 == self.Result[i][2]/10 || self.Result[i][1]/10 == self.Result[i][2]/10) {
						b1 = 4
					}
				}
			}
			total += self.GetMoneyPos(i, robot) * b1
		}

	} else if self.room.Type%230000 == 1 { //! 高倍场
		for i := 0; i < 4; i++ {
			total += self.GetMoneyPos(i, robot) * 11
		}
	}
	return total
}

//! 下注
func (self *Game_GoldBrNN) GameBets(uid int64, index int, gold int) {
	if uid == 0 {
		return
	}

	if index < 0 || index >= len(self.Bets) {
		return
	}

	if gold <= 0 {
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() >= int64(self.BetTime-2) {
		self.room.SendErr(uid, "正在开奖,请稍后下注")
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
			self.room.SendErr(uid, fmt.Sprintf("%d金币以上才能下注。", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("%d金币以上才能下注。", lib.GetManyMgr().GetProperty(self.room.Type).MinBet))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("%d万金币以上才能下注。", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/10000))
		}
		return
	}

	if person.Total < gold {
		self.room.SendErr(uid, "您的金币不足，请前往充值。")
		return
	}

	b1 := 11
	if self.room.Type%230000 == 0 { //! 低倍场
		for i := 0; i < len(self.Result[0]); i++ {
			if self.Result[0][i] == 0 {
				break
			}
			if self.Result[0][i]/10 <= 10 {
				b1 = 5
			}
		}
		if b1 == 5 {
			if !(self.Result[0][0]/10 == self.Result[0][1]/10 || self.Result[0][0]/10 == self.Result[0][2]/10 || self.Result[0][1]/10 == self.Result[0][2]/10) {
				b1 = 4
			}
		}
	}
	if (person.Bets+gold)*(b1-1) > (person.Total - gold) {
		self.room.SendErr(uid, "您的金币可能不够赔率,请前往充值.")
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

	{
		//!　模拟庄家是否够赔
		self.Total += gold
		self.Bets[index][person] += gold
		//! 得到庄家的钱
		dealmoney := 0
		if self.Dealer != nil {
			dealmoney = self.Dealer.Total
		} else if self.RobotDealer != nil {
			dealmoney = self.RobotDealer.GetMoney()
		} else if lib.GetManyMgr().GetProperty(self.room.Type).DealChange == 1 {
			dealmoney = self.Money
		}
		if dealmoney > 0 {
			dealwin := self.GetMaxLost(true)
			if dealmoney+self.Total+self.Robot.RobotTotal-dealwin < dealmoney/5 {
				self.Total -= gold
				self.Bets[index][person] -= gold
				self.room.SendErr(uid, "庄家已到最大赔率")
				self.OnBegin()
				return
			}
		}
	}

	person.Bets += gold
	person.Total -= gold
	person.BetInfo[index] += gold
	person.Round = 0

	var msg Msg_GameGoldBrNN_Bets
	msg.Uid = uid
	msg.Index = index
	msg.Gold = gold
	msg.Total = person.Total
	self.room.broadCastMsg("gamegoldbrnnbets", &msg)

}

//!　庄家可赢多少
func (self *Game_GoldBrNN) GetResultWin() (int, bool) {
	win := 0
	num := 0
	ct, maxcard := GetBrNiuNiuScore(self.Result[0])
	bs := GetGoldBrNNBS(ct, self.room.Type%230000)
	for i := 1; i < len(self.Result); i++ {
		_ct, _maxcard := GetBrNiuNiuScore(self.Result[i])
		_bs := GetGoldBrNNBS(_ct, self.room.Type%230000)
		dealwin := true
		if ct > _ct {
			dealwin = true
		} else if ct < _ct {
			dealwin = false
		} else {
			if maxcard < _maxcard {
				dealwin = false
			} else {
				dealwin = true
			}
		}
		if dealwin {
			for _, value := range self.Bets[i-1] {
				win += value * bs
			}
			if self.Dealer != nil { //! 玩家庄要计算机器人
				for _, value := range self.Robot.RobotsBet[i-1] {
					win += value * bs
				}
			}
			num++
		} else {
			for _, value := range self.Bets[i-1] {
				win -= value * _bs
			}
			if self.Dealer != nil { //! 玩家庄要计算机器人
				for _, value := range self.Robot.RobotsBet[i-1] {
					win -= value * _bs
				}
			}
			num--
		}
	}

	return win, (num == 4 || num == -4)
}

//! 得到机器人输赢
func (self *Game_GoldBrNN) GetRobotWin() (int, bool) {
	num := 0
	win := 0
	ct, maxcard := GetBrNiuNiuScore(self.Result[0])
	bs := GetGoldBrNNBS(ct, self.room.Type%230000)
	for i := 1; i < len(self.Result); i++ {
		_ct, _maxcard := GetBrNiuNiuScore(self.Result[i])
		_bs := GetGoldBrNNBS(_ct, self.room.Type%230000)
		dealwin := true
		if ct > _ct {
			dealwin = true
		} else if ct < _ct {
			dealwin = false
		} else {
			if maxcard < _maxcard {
				dealwin = false
			} else {
				dealwin = true
			}
		}
		if dealwin {
			if self.Dealer != nil { //! 玩家庄要计算机器人
				for _, value := range self.Robot.RobotsBet[i-1] {
					win -= value * bs
				}
			}
			num++
		} else {
			if self.Dealer != nil { //! 玩家庄要计算机器人
				for _, value := range self.Robot.RobotsBet[i-1] {
					win += value * _bs
				}
			}
			num--
		}
	}
	return win, (num == 4 || num == -4)
}

func (self *Game_GoldBrNN) RobotNeedWin() {
	for i := 0; i < 100; i++ {
		card := NewCard_NiuNiu(false)
		for j := 0; j < 5; j++ { //!　删除已发的牌
			self.Result[j] = self.Result[j][0:3] //! 把盖着的两张牌置空
			for k := 0; k < 3; k++ {
				card.Del(self.Result[j][k])
			}
		}

		for j := 0; j < 5; j++ { //! 模拟发牌
			c := card.Deal(2)
			self.Result[j] = append(self.Result[j], c...)
		}

		win, _ := self.GetRobotWin()
		if win >= 0 {
			return
		}
	}
}

func (self *Game_GoldBrNN) WinOrLost(iswin bool) {
	if iswin {
		winLst := make([]BrNN_Lst, 0)
		killLst := make([]BrNN_Lst, 0)
		for i := 0; i < 100; i++ {
			card := NewCard_NiuNiu(false)
			for j := 0; j < 5; j++ { //!　删除已发的牌
				self.Result[j] = self.Result[j][0:3] //! 把盖着的两张牌置空
				for k := 0; k < 3; k++ {
					card.Del(self.Result[j][k])
				}
			}

			for j := 0; j < 5; j++ { //! 模拟发牌
				c := card.Deal(2)
				self.Result[j] = append(self.Result[j], c...)
			}

			win, kill := self.GetResultWin()
			if win > 0 {
				//				lib.GetLogMgr().Output(lib.LOG_DEBUG, "``````````````````````````````")
				//				lib.GetLogMgr().Output(lib.LOG_DEBUG, " ", self.Result)
				//				lib.GetLogMgr().Output(lib.LOG_DEBUG, " ", win)

				var msg BrNN_Lst
				lib.HF_DeepCopy(&msg.Reuslt, &self.Result)
				if kill {
					killLst = append(killLst, msg)
				} else {
					winLst = append(winLst, msg)
				}
			}
		}
		if len(winLst) > 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "--------------- 赢1111")
			//			for i := 0; i < len(winLst); i++ {
			//				lib.GetLogMgr().Output(lib.LOG_DEBUG, "```````````---------------")
			//				lib.GetLogMgr().Output(lib.LOG_DEBUG, "-----------winlst : ", winLst[i])
			//			}

			self.Result = winLst[lib.HF_GetRandom(len(winLst))].Reuslt
		} else if len(killLst) > 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------------- 赢2222")
			self.Result = killLst[lib.HF_GetRandom(len(killLst))].Reuslt
		} else {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "--------------- 随机111111")
			card := NewCard_NiuNiu(false)
			for j := 0; j < 5; j++ {
				self.Result[j] = self.Result[j][0:3]
				for k := 0; k < 3; k++ {
					card.Del(self.Result[j][k])
				}
			}
			for j := 0; j < 5; j++ {
				c := card.Deal(2)
				self.Result[j] = append(self.Result[j], c...)
			}
		}

	} else {
		lostLst := make([]BrNN_Lst, 0)
		killLst := make([]BrNN_Lst, 0)
		for i := 0; i < 100; i++ {
			card := NewCard_NiuNiu(false)
			for j := 0; j < 5; j++ { //!　删除已发的牌
				self.Result[j] = self.Result[j][0:3] //! 把盖着的两张牌置空
				for k := 0; k < 3; k++ {
					card.Del(self.Result[j][k])
				}
			}

			for j := 0; j < 5; j++ { //! 模拟发牌
				c := card.Deal(2)
				self.Result[j] = append(self.Result[j], c...)
			}

			win, kill := self.GetResultWin()
			if win < 0 {
				var msg BrNN_Lst
				lib.HF_DeepCopy(&msg.Reuslt, &self.Result)
				if kill {
					killLst = append(killLst, msg)
				} else {
					lostLst = append(lostLst, msg)
				}
			}
		}

		if len(lostLst) > 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "---------- 输11111")
			self.Result = lostLst[lib.HF_GetRandom(len(lostLst))].Reuslt
		} else if len(killLst) > 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "---------- 输222222")
			self.Result = killLst[lib.HF_GetRandom(len(killLst))].Reuslt
		} else {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------- 随机2222222")
			card := NewCard_NiuNiu(false)
			for j := 0; j < 5; j++ {
				self.Result[j] = self.Result[j][0:3]
				for k := 0; k < 3; k++ {
					card.Del(self.Result[j][k])
				}
			}
			for j := 0; j < 5; j++ {
				c := card.Deal(2)
				self.Result[j] = append(self.Result[j], c...)
			}
		}
	}
}

//! 换庄
func (self *Game_GoldBrNN) ChageDeal() {
	if self.Dealer != nil {
		self.Dealer.Seat = -1
	} else if self.RobotDealer != nil {
		self.RobotDealer.SetSeat(-1)
	}

	self.Dealer = nil
	self.RobotDealer = nil
	for len(self.LstDeal) > 0 {
		if self.LstDeal[0].Robot != nil {
			find := false
			for i := 0; i < len(self.Robot.Robots); i++ {
				if self.Robot.Robots[i] == self.LstDeal[0].Robot {
					find = true
					break
				}
			}
			if !find { //! 要上庄的机器人已经走了
				self.LstDeal = self.LstDeal[1:]
				continue
			}
		}
		if self.LstDeal[0].GetTotal() >= lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney {
			if self.LstDeal[0].Person != nil {
				self.Dealer = self.LstDeal[0].Person
				self.Dealer.Seat = 100
				for i := 0; i < len(self.Seat); i++ {
					if self.Seat[i].Person == self.Dealer {
						var msg Msg_GameGoldBrNN_UpdSeat
						msg.Index = i
						self.room.broadCastMsg("gamegoldbrnnseat", &msg)
						self.Seat[i].Person = nil
						break
					}
				}
			} else if self.LstDeal[0].Robot != nil {
				self.RobotDealer = self.LstDeal[0].Robot
				self.RobotDealer.SetSeat(100)
				for i := 0; i < len(self.Seat); i++ {
					if self.Seat[i].Robot == self.RobotDealer {
						var msg Msg_GameGoldBrNN_UpdSeat
						msg.Index = i
						self.room.broadCastMsg("gamegoldbrnnseat", &msg)
						self.Seat[i].Robot = nil
						break
					}
				}
			}
			self.LstDeal = self.LstDeal[1:]
			break
		} else {
			self.LstDeal = self.LstDeal[1:]
		}
	}
	self.DownUid = 0
	self.Round = 0

	var msg Msg_GameGoldBZW_Deal
	if self.Dealer != nil {
		msg.Uid = self.Dealer.Uid
		msg.Name = self.Dealer.Name
		msg.Head = self.Dealer.Head
		msg.Total = self.Dealer.Total
		msg.IP = self.Dealer.IP
		msg.Address = self.Dealer.Address
		msg.Sex = self.Dealer.Sex
	} else if self.RobotDealer != nil {
		msg.Uid = self.RobotDealer.Id
		msg.Name = self.RobotDealer.Name
		msg.Head = self.RobotDealer.Head
		msg.Total = self.RobotDealer.GetMoney()
		msg.IP = self.RobotDealer.IP
		msg.Address = self.RobotDealer.Address
		msg.Sex = self.RobotDealer.Sex
	} else {
		msg.Total = self.Money
	}

	self.room.broadCastMsg("gamerob", &msg)
}

func (self *Game_GoldBrNN) OnBegin() {
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")

	if self.room.IsBye() {
		return
	}
	self.room.Begin = true

	win, _ := self.GetResultWin()

	if self.Dealer != nil {
		if self.Robot.RobotTotal == 0 {
			if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost == 102 { //!　平衡模式
				if GetServer().BrNNUserMoney[self.room.Type%230000]+int64(win) > lib.GetManyMgr().GetProperty(self.room.Type).PlayerMax { //!玩家必须输
					self.WinOrLost(false)
				} else if GetServer().BrNNUserMoney[self.room.Type%230000]+int64(win) < lib.GetManyMgr().GetProperty(self.room.Type).PlayerMin { //!　玩家必须赢
					self.WinOrLost(true)
				}
			} else if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost == 101 { //! 随机模式

			} else if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost <= 100 { //! 概率模式
				if lib.HF_GetRandom(100) < lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost { //! 玩家庄胜利
					self.WinOrLost(true)
				} else { //! 玩家庄失败
					self.WinOrLost(false)
				}
			}
		} else {
			win, _ := self.GetRobotWin()
			if lib.GetRobotMgr().GetRobotWin(self.room.Type)+win < 0 || win >= 0 {
				self.RobotNeedWin()
			}
		}
	} else {
		if GetServer().BrNNSysMoney[self.room.Type%230000]+int64(win) < lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin { //!  庄家必赢
			self.WinOrLost(true)
		} else if GetServer().BrNNSysMoney[self.room.Type%230000]+int64(win) > lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax { //! 专家必输
			self.WinOrLost(false)
		}
	}

	//	self.Result[0] = []int{84, 72, 63, 51, 114}
	//	self.Result[1] = []int{64, 92, 43, 134, 61}
	//	self.Result[2] = []int{53, 12, 62, 62, 73}
	//	self.Result[3] = []int{133, 71, 41, 44, 33}
	//	self.Result[4] = []int{24, 43, 94, 93, 102}

	self.OnEnd()
}

func (self *Game_GoldBrNN) OnEnd() {
	self.room.Begin = false

	dealwin := 0
	robotwin := 0
	kill := 0
	dealCT, dealMaxCard := GetBrNiuNiuScore(self.Result[0])
	dealBS := GetGoldBrNNBS(dealCT, self.room.Type%230000)
	ct := make([]int, 0)
	ct = append(ct, dealCT)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "庄家 : ", dealCT, " , ", dealMaxCard, " , ", dealBS)
	for i := 1; i < len(self.Result); i++ {
		xianCT, xianMaxCard := GetBrNiuNiuScore(self.Result[i])
		xianBS := GetGoldBrNNBS(xianCT, self.room.Type%230000)
		ct = append(ct, xianCT)
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "闲家 : ", xianCT, " , ", xianMaxCard, " , ", xianBS)
		win := true
		if dealCT > xianCT {
			win = true
		} else if dealCT < xianCT {
			win = false
		} else {
			if dealMaxCard < xianMaxCard {
				win = false
			} else {
				win = true
			}
		}

		if win { //!　庄赢
			for key, value := range self.Bets[i-1] {
				dealwin += value * dealBS
				key.Win -= (value*dealBS - value)
			}

			for key, value := range self.Robot.RobotsBet[i-1] {
				if self.Dealer != nil {
					dealwin += value * dealBS
				}
				key.AddWin(-(value*dealBS - value))
				robotwin -= value * dealBS
			}

			tmp := []int{0}
			tmp = append(tmp, self.Trend[i-1]...)
			if len(tmp) > 8 {
				tmp = tmp[0:8]
			}
			self.Trend[i-1] = tmp
			kill++
		} else {
			for key, value := range self.Bets[i-1] {
				dealwin -= value * xianBS
				key.Win += (value*xianBS + value)
				key.Cost += int(math.Ceil(float64(value*xianBS) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
			}

			for key, value := range self.Robot.RobotsBet[i-1] {
				if self.Dealer != nil {
					dealwin -= value * xianBS
				}
				key.AddWin((value*xianBS + value))
				key.AddCost(int(math.Ceil(float64(value*xianBS) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
				robotwin += value * xianBS
			}

			tmp := []int{1}
			tmp = append(tmp, self.Trend[i-1]...)
			if len(tmp) > 8 {
				tmp = tmp[0:8]
			}
			self.Trend[i-1] = tmp
		}
	}

	var bigwin *GameGold_BigWin
	for _, value := range self.PersonMgr {
		if value.Win > 0 {
			value.Win -= value.Cost
			GetServer().SqlAgentGoldLog(value.Uid, value.Cost, self.room.Type)
			GetServer().SqlAgentBillsLog(value.Uid, value.Cost/2, self.room.Type)
		} else if value.Win-value.Bets < 0 {
			cost := int(math.Ceil(float64(value.Win-value.Bets) * float64(lib.GetManyMgr().GetProperty(self.room.Type).Cost) / 200.0))
			GetServer().SqlAgentBillsLog(value.Uid, cost, self.room.Type)
		}

		if value.Win != 0 {
			value.Total += value.Win
			var msg Msg_GameGoldBrNN_Balance
			msg.Uid = value.Uid
			msg.Total = value.Total
			msg.Win = value.Win
			find := false
			for i := 0; i < len(self.Seat); i++ {
				if self.Seat[i].Person == value {
					self.room.broadCastMsg("gamegoldbrnnbalance", &msg)
					find = true
					break
				}
			}
			if !find {
				self.room.SendMsg(value.Uid, "gamegoldbrnnbalance", &msg)
			}
		}

		if value.Win > 0 {
			if bigwin == nil {
				bigwin = &GameGold_BigWin{value.Uid, value.Name, value.Head, value.Win}
			} else if value.Win > bigwin.Win {
				bigwin = &GameGold_BigWin{value.Uid, value.Name, value.Head, value.Win}
			}
		}

		if value.Bets > 0 {
			var record Rec_BrNN_Info
			record.Time = time.Now().Unix()
			record.GameType = self.room.Type
			var rec Son_Rec_BrNN_Person
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
		}

		if self.Robot.Robots[i].GetWin() != 0 {
			self.Robot.Robots[i].AddMoney(self.Robot.Robots[i].GetWin())
		}

		for j := 0; j < len(self.Seat); j++ {
			if self.Seat[j].Robot == self.Robot.Robots[i] {
				var msg Msg_GameGoldBrNN_Balance
				msg.Uid = self.Robot.Robots[i].Id
				msg.Total = self.Robot.Robots[i].GetMoney()
				msg.Win = self.Robot.Robots[i].GetWin()
				self.room.broadCastMsg("gamegoldbrnnbalance", &msg)
				break
			}
		}

		if self.Robot.Robots[i].GetWin() > 0 {
			if bigwin == nil {
				bigwin = &GameGold_BigWin{self.Robot.Robots[i].Id, self.Robot.Robots[i].Name, self.Robot.Robots[i].Head, self.Robot.Robots[i].GetWin()}
			} else if self.Robot.Robots[i].GetWin() > bigwin.Win {
				bigwin = &GameGold_BigWin{self.Robot.Robots[i].Id, self.Robot.Robots[i].Name, self.Robot.Robots[i].Head, self.Robot.Robots[i].GetWin()}
			}
		}
	}

	a, b := self.GetResultWin()
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "-----dealwin : ", dealwin, " a ： ", a, " b : ", b)

	if self.Dealer == nil && dealwin != 0 { //!  系统庄
		GetServer().SqlBZWLog(&SQL_BZWLog{1, dealwin, time.Now().Unix(), self.room.Type})
	}
	if self.Dealer != nil && robotwin != 0 { //! 玩家庄
		lib.GetRobotMgr().AddRobotWin(self.room.Type, robotwin)
		GetServer().SqlBZWLog(&SQL_BZWLog{1, robotwin, time.Now().Unix(), self.room.Type + 10000000})
	}
	if self.Dealer != nil && lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost == 102 {
		GetServer().SetBrNNUserMoney(self.room.Type%230000, GetServer().BrNNUserMoney[self.room.Type%230000]+int64(dealwin))
	}

	if self.Dealer != nil { //! 玩家庄
		if dealwin > 0 {
			bl := lib.GetManyMgr().GetProperty(self.room.Type).Cost
			cost := int(math.Ceil(float64(dealwin) * bl / 100.0))
			dealwin -= cost
			GetServer().SqlAgentGoldLog(self.Dealer.Uid, cost, self.room.Type)
			GetServer().SqlAgentBillsLog(self.Dealer.Uid, cost/2, self.room.Type)
		} else if dealwin < 0 {
			cost := int(math.Ceil(float64(dealwin) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 200.0))
			GetServer().SqlAgentBillsLog(self.Dealer.Uid, cost, self.room.Type)
		}
	} else {
		_dealwin := dealwin
		if _dealwin > 0 {
			bl := lib.GetManyMgr().GetProperty(self.room.Type).DealCost
			cost := int(math.Ceil(float64(_dealwin) * bl / 100.0))
			_dealwin -= cost
		}
		GetServer().SetBrNNMoney(self.room.Type%230000, GetServer().BrNNSysMoney[self.room.Type%230000]+int64(_dealwin))
		dealwin -= robotwin
		if dealwin > 0 {
			bl := lib.GetManyMgr().GetProperty(self.room.Type).Cost
			cost := int(math.Ceil(float64(dealwin) * bl / 100.0))
			dealwin -= cost
		}
	}

	if self.Dealer != nil {
		var record Rec_BrNN_Info
		record.Time = time.Now().Unix()
		record.GameType = self.room.Type
		var rec Son_Rec_BrNN_Person
		rec.Uid = self.Dealer.Uid
		rec.Name = self.Dealer.Name
		rec.Head = self.Dealer.Head
		rec.Score = dealwin
		rec.Result = self.Result
		rec.Bets = self.Dealer.BetInfo
		record.Info = append(record.Info, rec)
		GetServer().InsertRecord(self.room.Type, self.Dealer.Uid, lib.HF_JtoA(&record), rec.Score)
	}

	{
		//!　庄家信息
		var msg Msg_GameGoldBrNN_Balance
		if self.Dealer != nil {
			self.Dealer.Total += dealwin
			msg.Uid = self.Dealer.Uid
			msg.Total = self.Dealer.Total
		} else if self.RobotDealer != nil { //! 机器人庄
			self.RobotDealer.AddMoney(dealwin)
			msg.Uid = self.RobotDealer.Id
			msg.Total = self.RobotDealer.GetMoney()
		} else {
			if lib.GetManyMgr().GetProperty(self.room.Type).DealChange == 1 {
				self.Money += dealwin
				if self.Money <= 0 {
					self.Money = lib.GetManyMgr().GetProperty(self.room.Type).DealInitMoney
				}
			} else {
				self.Money = lib.GetManyMgr().GetProperty(self.room.Type).DealInitMoney
			}
			msg.Total = self.Money
		}
		msg.Win = dealwin
		self.room.broadCastMsg("gamegoldbrnnbalance", &msg)
	}

	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 12
	self.SetTime(self.BetTime)

	{
		//!　总结算
		var msg Msg_GameGoldBrnn_End
		msg.Result = self.Result
		msg.CT = ct
		msg.Money = lib.GetManyMoneyMgr().GetProperty(self.room.Type).Money
		msg.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime
		if bigwin != nil {
			msg.Uid = bigwin.Uid
			msg.Name = bigwin.Name
			msg.Head = bigwin.Head
		}
		self.GameCard()

		//		self.Result[0] = []int{84, 72, 63, 51, 114}
		//		self.Result[1] = []int{64, 92, 43, 134, 61}
		//		self.Result[2] = []int{53, 12, 62, 62, 73}
		//		self.Result[3] = []int{133, 71, 41, 44, 33}
		//		self.Result[4] = []int{24, 43, 94, 93, 102}

		msg.Kill = 0
		if kill == 4 { //! 通杀
			msg.Kill = 1
		} else if kill == 0 { //! 通赔
			msg.Kill = 2
		}
		lib.HF_DeepCopy(&msg.Next, &self.Result)
		for i := 0; i < 5; i++ {
			msg.Next[i][3] = 0
			msg.Next[i][4] = 0
		}
		self.room.broadCastMsg("gamegoldbrnnend", &msg)
	}

	self.Total = 0

	//!清理玩家
	for key, value := range self.PersonMgr {
		if value.Online {
			find := false
			for i := 0; i < len(self.LstDeal); i++ {
				if self.LstDeal[i].Person == value {
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
				for j := 0; j < len(value.BeBetInfo); j++ {
					value.BetInfo[j] = 0
				}
				continue
			}
		}

		//!　走的人在上庄列表上
		for j := 0; j < len(self.LstDeal); j++ {
			if self.LstDeal[j].Person == value {
				copy(self.LstDeal[j:], self.LstDeal[j+1:])
				self.LstDeal = self.LstDeal[:len(self.LstDeal)-1]
				break
			}
		}

		//! 走的人是庄家
		if self.Dealer == value {
			self.ChageDeal()
		}

		//! 走的人是位置上的人
		for j := 0; j < len(self.Seat); j++ {
			if self.Seat[j].Person == value {
				self.Seat[j].Person = nil
				var msg Msg_GameGoldBrNN_UpdSeat
				msg.Index = j
				self.room.broadCastMsg("gamegoldbrnnseat", &msg)
				break
			}
		}
		delete(self.PersonMgr, key)
	}

	//! 返回机器人
	self.Robot.Init(4, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)

	for i := 0; i < len(self.Bets); i++ {
		self.Bets[i] = make(map[*Game_GoldBrNN_Person]int)
	}

	//! 判断庄家是否能继续连庄
	if self.Dealer != nil {
		if self.Dealer.Total < lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney || self.DownUid == self.Dealer.Uid || GetPersonMgr().GetPerson(self.Dealer.Uid) == nil {
			self.ChageDeal()
		} else {
			if self.Round >= 10 && len(self.LstDeal) > 0 {
				self.ChageDeal()
			} else {
				self.Round++
			}
		}
	} else if self.RobotDealer != nil {
		if self.RobotDealer.GetMoney() < lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney || !lib.GetRobotMgr().GetRobotSet(self.room.Type).NeedRobot {
			self.ChageDeal()
		} else {
			if self.Round >= lib.HF_GetRandom(6)+3 && len(self.LstDeal) > 0 {
				self.ChageDeal()
			} else {
				self.Round++
			}
		}
	} else if len(self.LstDeal) > 0 {
		self.ChageDeal()
	}

	//! 判断坐下的人是否还能坐下
	for i := 0; i < len(self.Seat); i++ {
		if self.Seat[i].Person == nil {
			continue
		}
		if self.Seat[i].Person.Total < lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney {
			self.Seat[i].Person.Seat = -1
			var msg Msg_GameGoldBrNN_UpdSeat
			msg.Index = i
			self.room.broadCastMsg("gamegoldbrnnseat", &msg)
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
			var msg Msg_GameGoldBrNN_UpdSeat
			msg.Index = i
			self.room.broadCastMsg("gamegoldbrnnseat", &msg)
			self.Seat[i].Robot = nil
		}
	}
}

func (self *Game_GoldBrNN) OnInit(room *Room) {
	self.room = room

	self.Money = lib.GetManyMgr().GetProperty(self.room.Type).DealInitMoney
	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 12

	//! 载入机器人
	self.Robot.Init(4, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)
}

func (self *Game_GoldBrNN) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldBrNN) OnSendInfo(person *Person) {
	if self.Time == 0 {
		self.GameCard()
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
		person.SendMsg("gamegoldbrnninfo", self.getInfo(person.Uid, value.Total))
		return
	}

	_person := new(Game_GoldBrNN_Person)
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
	person.SendMsg("gamegoldbrnninfo", self.getInfo(person.Uid, person.Gold))
}

func (self *Game_GoldBrNN) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold": //! 同步金币
		person := self.GetPerson(msg.V.(*staticfunc.Msg_SynchroGold).Uid)
		if person != nil {
			person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal(person.Uid, person.Total)
		}
	case "gamebrttzbets": //! 下注
		self.GameBets(msg.Uid, msg.V.(*Msg_GameGoldBZW_Bets).Index, msg.V.(*Msg_GameGoldBZW_Bets).Gold)
	case "gamerob": //!　上庄
		self.GameUpDeal(msg.Uid)
	case "gameredeal": //! 下庄
		self.GameReDeal(msg.Uid)
	case "gamebrttzseat": //! 坐下
		self.GameSeat(msg.Uid, msg.V.(*Msg_GameGoldBZW_Seat).Index)
	case "gameplayerlist": //申请无座玩家列表
		self.GamePlayerList(msg.Uid)
	}
}

func (self *Game_GoldBrNN) OnBye() {

}

func (self *Game_GoldBrNN) OnExit(uid int64) {
	value, ok := self.PersonMgr[uid]
	if ok {
		value.Online = false
		//! 退出房间同步金币
		gold := value.Total - value.Gold
		if gold > 0 {
			GetRoomMgr().AddCard(value.Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
		} else if gold < 0 {
			GetRoomMgr().CostCard(value.Uid, staticfunc.TYPE_GOLD, -gold, self.room)
		}
		value.Gold = value.Total
	}

}

func (self *Game_GoldBrNN) OnIsDealer(uid int64) bool {
	if self.Dealer != nil && self.Dealer == self.GetPerson(uid) {
		return true
	}
	return false
}

func (self *Game_GoldBrNN) OnIsBets(uid int64) bool {
	value, ok := self.PersonMgr[uid]
	if ok {
		return value.Bets > 0
	}
	return false
}

func (self *Game_GoldBrNN) OnBalance() {
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

func (self *Game_GoldBrNN) OnTime() {
	if lib.GetRobotMgr().GetRobotSet(self.room.Type).Dealer && len(self.Robot.Robots) > 0 && len(self.LstDeal) < 5 { //! 需要机器人上庄
		self.RobotUpDeal(self.Robot.Robots[lib.HF_GetRandom(len(self.Robot.Robots))])
	}

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
		self.RobotSeat(lib.HF_GetRandom(len(self.Seat)), self.Robot.Robots[i])
	}

	if time.Now().Unix() < self.Time {
		if self.Dealer == nil && self.RobotDealer == nil && lib.GetManyMgr().GetProperty(self.room.Type).SysNoBets == 1 {
			return
		}

		if self.Time-time.Now().Unix() >= int64(lib.GetManyMgr().GetProperty(self.room.Type).BetTime) {
			return
		}

		for i := 0; i < len(self.Robot.Robots); i++ {
			if self.Robot.Robots[i].GetSeat() == 100 { //! 庄家不能下注
				continue
			}

			if lib.HF_GetRandom(100) >= 100-lib.GetRobotMgr().GetRobotSet(self.room.Type).BetRate {
				continue
			}

			index, gold, bets := self.Robot.GameBets(self.Robot.Robots[i])
			if gold == 0 {
				continue
			}

			b1 := 11
			if self.room.Type%230000 == 0 { //! 低倍场
				for i := 0; i < len(self.Result[0]); i++ {
					if self.Result[0][i] == 0 {
						break
					}
					if self.Result[0][i]/10 <= 10 {
						b1 = 5
					}
				}
				if b1 == 5 {
					if !(self.Result[0][0]/10 == self.Result[0][1]/10 || self.Result[0][0]/10 == self.Result[0][2]/10 || self.Result[0][1]/10 == self.Result[0][2]/10) {
						b1 = 4
					}
				}
			}
			if bets*b1 > self.Robot.Robots[i].GetMoney() {
				self.Robot.GameBackBets(self.Robot.Robots[i], index, gold)
				continue
			}
			if self.Dealer != nil { //! 玩家庄判断是否能下
				dealwin := self.GetMaxLost(true)
				if self.Dealer.Total+self.Total+self.Robot.RobotTotal-dealwin < self.Dealer.Total/5 {
					self.Robot.GameBackBets(self.Robot.Robots[i], index, gold)
					continue
				}
			} else if self.RobotDealer != nil {
				dealwin := self.GetMaxLost(true)
				if self.RobotDealer.GetMoney()+self.Total+self.Robot.RobotTotal-dealwin < self.RobotDealer.GetMoney()/5 {
					self.Robot.GameBackBets(self.Robot.Robots[i], index, gold)
					continue
				}
			} else if lib.GetManyMgr().GetProperty(self.room.Type).DealChange == 1 {
				dealwin := self.GetMaxLost(true)
				if self.Money+self.Total+self.Robot.RobotTotal-dealwin < self.Money/5 {
					self.Robot.GameBackBets(self.Robot.Robots[i], index, gold)
					continue
				}
			}
			var msg Msg_GameGoldBrNN_Bets
			msg.Uid = self.Robot.Robots[i].Id
			msg.Index = index
			msg.Gold = gold
			msg.Total = self.Robot.Robots[i].GetMoney()
			self.room.broadCastMsg("gamegoldbrnnbets", &msg)
		}
		return
	}

	if !self.room.Begin {
		self.OnBegin()
	}
}
