package gameserver

import (
	"fmt"
	"lib"
	"math"
	"sort"
	"staticfunc"
	"time"
)

var GOLDBJL_BS []int = []int{12, 12, 2, 9, 2, 3, 33, 3} //! 闲对,庄对,闲,平,庄,闲天王,同点平,庄天王

type Rec_GoldBJL_Info struct {
	GameType int                      `json:"gametype"`
	Time     int64                    `json:"time"`
	Info     []Son_Rec_GoldBJL_Person `json:"info"`
}

type Son_Rec_GoldBJL_Person struct {
	Uid    int64    `json:"uid"`
	Name   string   `json:"name"`
	Head   string   `json:"head"`
	Scroe  int      `json:"score"`
	Result [2][]int `json:"result"`
	Bets   [8]int   `json:"bets"`
}

type Game_GoldBJLSeat struct {
	Person *Game_GoldBJL_Person
	Robot  *lib.Robot
}

func (self *Game_GoldBJLSeat) GetTotal() int {
	if self.Person != nil {
		return self.Person.Total
	} else if self.Robot != nil {
		return self.Robot.GetMoney()
	}
	return 0
}

type Game_GoldBJL struct {
	PersonMgr   map[int64]*Game_GoldBJL_Person  `json:"personmgr"`
	Dealer      *Game_GoldBJL_Person            `json:"dealer"`
	RobotDealer *lib.Robot                      `json:"robotdealer"` //! 机器人庄
	LstDeal     []Game_GoldBJLSeat              `json:"lstdeal"`
	Card        *CardMgr                        `json:"card"`
	Step        int                             `json:"Step"`    //! 第几局
	Round       int                             `json:"round"`   //! 连庄轮数
	DownUid     int64                           `json:"downuid"` //! 下庄uid
	Bets        [8]map[*Game_GoldBJL_Person]int `json:"bets"`    //! 闲对,庄对,闲,平,庄,闲天王,同点平,庄天王
	Result      [2][]int                        `json:"result"`  //! 0-闲 1-庄
	Time        int64                           `json:"time"`
	Seat        [6]Game_GoldBJLSeat             `json:"seat"`
	Total       int                             `json:"total"`
	Trend       []Msg_GameGoldBjl_Trend         `json:"trend"` //! 0-闲赢 1-庄赢 2-平局
	Money       int                             `json:"money"`
	WinLost     [2]int                          `json:"winlost"`
	Robot       lib.ManyGameRobot               //! 机器人结构
	BetTime     int                             `json:"bettime"`

	room *Room
}

func NewGame_GoldBJL() *Game_GoldBJL {
	game := new(Game_GoldBJL)
	game.PersonMgr = make(map[int64]*Game_GoldBJL_Person)
	card := NewCard_BJL156()
	for i := 0; i < 50; i++ {
		game.Result[0] = card.Deal(2)
		game.Result[1] = card.Deal(2)
		game.BoCard(game.Result[0], game.Result[1], card)
		ct, _ := game.IsType(game.Result[0], game.Result[1])
		xNum := 0
		zNum := 0
		for i := 0; i < len(game.Result[0]); i++ {
			if game.Result[0][i]/10 < 10 {
				xNum += game.Result[0][i] / 10
			}
		}
		xNum = xNum % 10

		for i := 0; i < len(game.Result[1]); i++ {
			if game.Result[1][i]/10 < 10 {
				zNum += game.Result[1][i] / 10
			}
		}
		zNum = zNum % 10

		var msg Msg_GameGoldBjl_Trend
		lib.HF_DeepCopy(&msg.Result, &game.Result)
		msg.Num = append(msg.Num, xNum)
		msg.Num = append(msg.Num, zNum)
		msg.Trend = ct

		game.Trend = append(game.Trend, msg)
	}
	for i := 0; i < len(game.Bets); i++ {
		game.Bets[i] = make(map[*Game_GoldBJL_Person]int)
	}
	game.WinLost[0] = 0
	game.WinLost[1] = 0

	return game
}

type Game_GoldBJL_Person struct {
	Uid       int64  `json:"uid"`
	Gold      int    `json:"gold"`
	Total     int    `json:"total"`
	Win       int    `json:"win"`       //! 赢钱
	Cost      int    `json:"cost"`      //! 抽水
	Bets      int    `json:"bets"`      //! 单局下注
	BetInfo   [8]int `json:"betinfo"`   //! 下注详情
	BeBets    int    `json:"bebets"`    //! 上局下注
	BeBetInfo [8]int `json:"bebetinfo"` //! 上局下注详情
	Name      string `json:"name"`
	Head      string `json:"head"`
	Online    bool   `json:"online"`
	Round     int    `json:"round"` //! 未下注轮数
	Seat      int    `json:"seat"`
	IP        string `json:"ip"`
	Address   string `json:"address"`
	Sex       int    `json:"sex"`
}

type Msg_GameGoldBJL_Total struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"`
}

type Son_GameGoldBJL_Info struct {
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Total   int    `json:"total"`
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type Msg_GameGoldBJL_List struct {
	Info []Son_GameGoldBJL_Info `json:"info"`
}

type Msg_GameGoldBjl_UpdSeat struct {
	Index   int    `json:"index"`
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Total   int    `json:"total"` //! 金币数量
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type Msg_GameGoldBJL_Bets struct {
	Uid   int64 `json:"uid"`
	Index int   `json:"index"` //! 下注位置
	Gold  int   `json:"gold"`  //! 下了多少金币
	Total int   `json:"total"` //! 总金币
}

type Msg_GameGoldBJL_GoOn struct {
	Uid   int64  `json:"uid"`
	Gold  [8]int `json:"gold"`  //! 每个区域下了多少
	Total int    `json:"total"` //! 总金币
}

type Msg_GameGoldBJL_Info struct {
	Begin   bool                    `json:"begin"`  //! 是否开始游戏
	Time    int64                   `json:"time"`   //! 倒计时
	Seat    [6]Son_GameGoldBJL_Info `json:"seat"`   //! 座位信息
	Bets    [8]int                  `json:"bets"`   //! 每个区域的下注
	Dealer  Son_GameGoldBJL_Info    `json:"dealer"` //! 庄
	IsDeal  bool                    `json:"isdeal"` //! 是否可以下庄
	Total   int                     `json:"total"`  //! 自己的金币
	Card    []int                   `json:"card"`   //! 牌组
	Step    int                     `json:"step"`   //! 当前局数
	Trend   []Msg_GameGoldBjl_Trend `json:"trend"`  //! 走势
	WinLost [2]int                  `json:"winlost"`
	Money   []int                   `json:"money"`
	BetTime int                     `json:"bettime"`
}

type Msg_GameGoldBJL_Balance struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"` //! 总金币
	Win   int   `json:"win"`   //! 赢了多少
}

type Msg_GameGoldBJL_End struct {
	Uid     int64                   `json:"uid"`
	Name    string                  `json:"name"`
	Head    string                  `json:"head"`
	CT      []int                   `json:"ct"`
	Card    []int                   `json:"card"` //! 牌组
	Step    int                     `json:"step"` //! 当前局数
	Result  [2][]int                `json:"result"`
	Trend   []Msg_GameGoldBjl_Trend `json:"trend"`
	WinLost [2]int                  `json:"winlost"`
	Money   []int                   `json:"money"`
	BetTime int                     `json:"bettime"`
}

type Msg_GameGoldBjl_Trend struct {
	Result [2][]int `json:"result"` //! 0-闲 1-庄 	牌
	Num    []int    `json:"num"`    //! 0-闲 1-庄 	点数
	Trend  int      `json:"trend"`  //! 0-闲赢 1-庄赢 2-平局
}

type Msg_GameGoldBJL_DealList struct {
	Type int                    `json:"type"` //! 0-上庄 1-下庄
	Info []Son_GameGoldBJL_Info `json:"info"`
}

type Msg_GameGoldBJL_Deal struct {
	Uid     int64  `json:"uid"`
	Head    string `json:"head"`
	Name    string `json:"name"`
	Total   int    `json:"total"`
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type GameGoldBJL_CanResult struct {
	XCard []int `json:"xcard"`
	ZCard []int `json:"zcard"`
}

func (self *Game_GoldBJL) getInfo(uid int64, total int) *Msg_GameGoldBJL_Info {
	var msg Msg_GameGoldBJL_Info
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
		msg.Bets[i] = self.GetMoneyPos(i, false)
	}
	if self.Dealer == nil { //! 系统庄的时候上面计算了机器人的下注
		for i := 0; i < len(self.Bets); i++ {
			msg.Bets[i] += self.GetMoneyPosByRobot(i)
		}
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

	msg.Step = self.Step
	if self.Card != nil {
		lib.HF_DeepCopy(&msg.Card, &self.Card.Card)
	} else {
		self.Card = NewCard_BJL156()
		lib.HF_DeepCopy(&msg.Card, &self.Card.Card)
	}
	msg.WinLost = self.WinLost

	return &msg
}

func (self *Game_GoldBJL) GetPerson(uid int64) *Game_GoldBJL_Person {
	return self.PersonMgr[uid]
}

func (self *Game_GoldBJL) GetMoneyPos(index int, robot bool) int { //! 获取该位置一共下了多少
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

func (self *Game_GoldBJL_Person) SynchroGold(gold int) { //!　同步金币
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

func (self *Game_GoldBJL) SendTotal(uid int64, total int) { //! 同步总分
	var msg Msg_GameGoldBJL_Total
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

func (self *Game_GoldBJL) SetTime(t int) { //! 设置时间
	if t == 0 {
		self.Time = 0
	} else {
		self.Time = time.Now().Unix() + int64(t)
	}
	var msg Msg_SetTime
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	self.room.broadCastMsg("gametime", &msg)
}

func (self *Game_GoldBJL) GamePlayerList(uid int64) { //! 获取无座玩家列表
	var msg Msg_GameGoldBJL_List
	tmp := make(map[int64]Son_GameGoldBJL_Info)
	for _, value := range self.PersonMgr {
		if value.Seat >= 0 {
			continue
		}

		var node Son_GameGoldBJL_Info
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

		var node Son_GameGoldBJL_Info
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

func (self *Game_GoldBJL) GameSeat(uid int64, index int) { //! 坐下
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if self.Dealer == person {
		self.room.SendErr(uid, "庄家无法坐下")
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

	var msg Msg_GameGoldBjl_UpdSeat
	msg.Uid = uid
	msg.Index = index
	msg.Head = person.Head
	msg.Name = person.Name
	msg.Total = person.Total
	msg.IP = person.IP
	msg.Address = person.Address
	msg.Sex = person.Sex
	self.room.broadCastMsg("gamegoldbjlseat", &msg)
}

//! 机器人坐下
func (self *Game_GoldBJL) RobotSeat(index int, robot *lib.Robot) {
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

	var msg Msg_GameGoldBjl_UpdSeat
	msg.Uid = robot.Id
	msg.Index = index
	msg.Head = robot.Head
	msg.Name = robot.Name
	msg.Total = robot.GetMoney()
	msg.IP = robot.IP
	msg.Address = robot.Address
	msg.Sex = robot.Sex
	self.room.broadCastMsg("gamegoldbjlseat", &msg)
}

func (self *Game_GoldBJL) GameUpDeal(uid int64) { //! 上庄
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
				self.room.SendErr(uid, "您已经在上庄列表中,请等待上庄")
				return
			}
		}
		if len(self.LstDeal) == 0 {
			self.Round = 0
		}
		self.LstDeal = append(self.LstDeal, Game_GoldBJLSeat{person, nil})
	}

	person.Round = 0

	var msg Msg_GameGoldBJL_DealList
	msg.Type = 0
	msg.Info = make([]Son_GameGoldBJL_Info, 0)
	for i := 0; i < len(self.LstDeal); i++ {
		if self.LstDeal[i].Person != nil {
			msg.Info = append(msg.Info, Son_GameGoldBJL_Info{self.LstDeal[i].Person.Uid, self.LstDeal[i].Person.Name, self.LstDeal[i].Person.Head, self.LstDeal[i].Person.Total, self.LstDeal[i].Person.IP, self.LstDeal[i].Person.Address, self.LstDeal[i].Person.Sex})
		} else if self.LstDeal[i].Robot != nil {
			msg.Info = append(msg.Info, Son_GameGoldBJL_Info{self.LstDeal[i].Robot.Id, self.LstDeal[i].Robot.Name, self.LstDeal[i].Robot.Head, self.LstDeal[i].Robot.GetMoney(), self.LstDeal[i].Robot.IP, self.LstDeal[i].Robot.Address, self.LstDeal[i].Robot.Sex})
		}
	}
	self.room.SendMsg(uid, "gamegoldbjldeal", &msg)
}

//! 机器人上庄
func (self *Game_GoldBJL) RobotUpDeal(robot *lib.Robot) {
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
	self.LstDeal = append(self.LstDeal, Game_GoldBJLSeat{nil, robot})
}

func (self *Game_GoldBJL) GameReDeal(uid int64) { //! 下庄
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if self.Dealer == person {
		self.DownUid = uid
		self.room.SendErr(uid, "您已成功下庄,请等待本局结束")
	} else {
		for i := 0; i < len(self.LstDeal); i++ {
			if self.LstDeal[i].Person == person {
				copy(self.LstDeal[i:], self.LstDeal[i+1:])
				self.LstDeal = self.LstDeal[:len(self.LstDeal)-1]
				break
			}
		}
	}
	var msg Msg_GameGoldBJL_DealList
	msg.Type = 1
	msg.Info = make([]Son_GameGoldBJL_Info, 0)
	self.room.SendMsg(uid, "gamegoldbjldeal", &msg)
}

func (self *Game_GoldBJL) ChageDeal() { //! 换庄
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
						var msg Msg_GameGoldBjl_UpdSeat
						msg.Index = i
						self.room.broadCastMsg("gamegoldbjlseat", &msg)
						self.Seat[i].Person = nil
						break
					}
				}
			} else if self.LstDeal[0].Robot != nil {
				self.RobotDealer = self.LstDeal[0].Robot
				self.RobotDealer.SetSeat(100)
				for i := 0; i < len(self.Seat); i++ {
					if self.Seat[i].Robot == self.RobotDealer {
						var msg Msg_GameGoldBjl_UpdSeat
						msg.Index = i
						self.room.broadCastMsg("gamegoldbjlseat", &msg)
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

func (self *Game_GoldBJL) GetMaxLost(robot bool, he bool) int { //! 获取庄家最多可能输多少
	total := 0
	//! 闲赢最多可能输多少
	total += self.GetMoneyPos(0, robot) * GOLDBJL_BS[0] //! 闲对
	total += self.GetMoneyPos(2, robot) * GOLDBJL_BS[2] //! 闲赢
	total += self.GetMoneyPos(5, robot) * GOLDBJL_BS[5] //! 闲天王

	//! 庄赢最多可能输多少
	zWin := self.GetMoneyPos(1, robot) * GOLDBJL_BS[1] //! 庄对
	zWin += self.GetMoneyPos(4, robot) * GOLDBJL_BS[4] //! 庄赢
	zWin += self.GetMoneyPos(7, robot) * GOLDBJL_BS[7] //! 庄天王

	if zWin > total {
		total = zWin
	}

	//! 和最多可能输多少  闲-99 庄-99 和，同点和，庄天王，闲天王，庄对，闲对
	if he {
		hWin := self.GetMoneyPos(3, robot) * GOLDBJL_BS[3] //! 和
		hWin += self.GetMoneyPos(6, robot) * GOLDBJL_BS[6] //! 同点和
		hWin += self.GetMoneyPos(0, robot) * GOLDBJL_BS[0] //! 闲对
		hWin += self.GetMoneyPos(5, robot) * GOLDBJL_BS[5] //! 闲天王
		hWin += self.GetMoneyPos(1, robot) * GOLDBJL_BS[1] //! 庄对
		hWin += self.GetMoneyPos(7, robot) * GOLDBJL_BS[7] //! 庄天王
		hWin += self.GetMoneyPos(2, robot)                 //! 退闲钱
		hWin += self.GetMoneyPos(4, robot)                 //! 退庄钱

		if hWin > total {
			total = hWin
		}
	}

	return total
}

func (self *Game_GoldBJL) GameBets(uid int64, index int, gold int) { //! 下注
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
			dealwin := self.GetMaxLost(true, index != 2 && index != 4)
			if dealmoney+self.Total+self.Robot.RobotTotal-dealwin < dealmoney/5 {
				self.Total -= gold
				self.Bets[index][person] -= gold
				self.room.SendErr(uid, "该区域庄家已到最大赔率")
				return
			}
		}
	}

	person.Bets += gold
	person.Total -= gold
	person.BetInfo[index] += gold
	person.Round = 0

	var msg Msg_GameGoldBJL_Bets
	msg.Uid = uid
	msg.Index = index
	msg.Gold = gold
	msg.Total = person.Total
	self.room.broadCastMsg("gamegoldbjlbets", &msg)
}

func (self *Game_GoldBJL) GameGoOn(uid int64) { //! 续压
	if uid == 0 {
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

	if person.BeBets == 0 {
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

	if person.Total < person.BeBets {
		self.room.SendErr(uid, "您的金币不足，请前往充值。")
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

	{
		//!　模拟庄家是否够赔
		for i := 0; i < len(person.BeBetInfo); i++ {
			self.Total += person.BeBetInfo[i]
			self.Bets[i][person] += person.BeBetInfo[i]
		}
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
			dealwin := self.GetMaxLost(true, true)
			if dealmoney+self.Total+self.Robot.RobotTotal-dealwin < dealmoney/5 {
				for i := 0; i < len(person.BeBetInfo); i++ {
					self.Total -= person.BeBetInfo[i]
					self.Bets[i][person] -= person.BeBetInfo[i]
				}
				self.room.SendErr(uid, "庄家已到最大赔率")
				return
			}
		}
	}

	person.Bets += person.BeBets
	person.Total -= person.BeBets
	for i := 0; i < len(person.BeBetInfo); i++ {
		person.BetInfo[i] = person.BeBetInfo[i]
	}
	person.Round = 0

	var msg Msg_GameGoldBJL_GoOn
	msg.Uid = uid
	msg.Gold = person.BeBetInfo
	msg.Total = person.Total
	self.room.broadCastMsg("gamegoldbjlgoon", &msg)

}

func (self *Game_GoldBJL) GetDealWin(card []int, _card []int) int { //!　庄家可以赢多少
	_, trend := self.IsType(card, _card)
	win := 0
	for i := 0; i < len(trend); i++ {
		if trend[i] == 1 { //! 庄家在这个区域输了
			win += self.GetMoneyPos(i, false) * GOLDBJL_BS[i]
		}
		if trend[i] == 2 {
			win += self.GetMoneyPos(i, false)
		}
	}

	if self.Dealer == nil {
		return self.Total - win
	} else {
		return self.Total + self.Robot.RobotTotal - win
	}
}

//! 博牌
func (self *Game_GoldBJL) BoCard(card []int, _card []int, cardmgr *CardMgr) ([]int, []int) {
	xNum := 0 //! 闲点数
	zNum := 0 //! 庄点数
	for i := 0; i < len(card); i++ {
		if card[i]/10 < 10 {
			xNum += card[i] / 10
		}
		if _card[i]/10 < 10 {
			zNum += _card[i] / 10
		}
	}

	xNum = xNum % 10
	zNum = zNum % 10

	if xNum >= 8 || zNum >= 8 { //! 庄,闲点数为8,9,既定胜负,无法博牌
		return card, _card
	}

	xCard := -1    //! 闲家博的第三张牌的点数
	if xNum <= 5 { //! 闲家博牌
		xCard = cardmgr.Deal(1)[0]
		card = append(card, xCard)
		xCard = xCard / 10
	}

	zCard := -1    //!　庄家博的第三张牌的点数
	if zNum <= 2 { //! 庄家博牌
		zCard = -100
	} else if zNum == 3 {
		if xCard != 8 {
			zCard = -100
		}
	} else if zNum == 4 {
		if xCard != 1 && xCard != 8 && xCard != 9 && xCard != 0 {
			zCard = -100
		}
	} else if zNum == 5 {
		if xCard != 1 && xCard != 2 && xCard != 3 && xCard != 8 && xCard != 9 && xCard != 0 {
			zCard = -100
		}
	} else if zNum == 6 {
		if xCard != 1 && xCard != 2 && xCard != 3 && xCard != 4 && xCard != 5 && xCard != 8 && xCard != 9 && xCard != 0 {
			zCard = -100
		}
	}
	if zCard == -100 {
		zCard = cardmgr.Deal(1)[0]
		_card = append(_card, zCard)
		zCard = zCard / 10
	}
	return card, _card
}

//! 判断输赢
func (self *Game_GoldBJL) IsType(cards []int, _cards []int) (int, []int) { //! int(0-闲赢 1-庄赢 2-平局) []int(8区域的输赢 0-输 1-赢)
	card := make([]int, 0)
	_card := make([]int, 0)
	lib.HF_DeepCopy(&card, &cards)
	lib.HF_DeepCopy(&_card, &_cards)
	trend := make([]int, 0)
	win := -1
	xNum := 0
	zNum := 0
	for i := 0; i < len(card); i++ {
		if card[i]/10 < 10 {
			xNum += card[i] / 10
		}
	}
	xNum = xNum % 10

	for i := 0; i < len(_card); i++ {
		if _card[i]/10 < 10 {
			zNum += _card[i] / 10
		}
	}
	zNum = zNum % 10

	if card[0]/10 == card[1]/10 { //! 闲对
		trend = append(trend, 1)
	} else {
		trend = append(trend, 0)
	}

	if _card[0]/10 == _card[1]/10 { //! 庄对
		trend = append(trend, 1)
	} else {
		trend = append(trend, 0)
	}

	if xNum > zNum { //!　闲,平,庄
		trend = append(trend, []int{1, 0, 0}...)
		win = 0
	} else if xNum < zNum {
		trend = append(trend, []int{0, 0, 1}...)
		win = 1
	} else {
		win = 2
		trend = append(trend, []int{2, 1, 2}...)
	}

	if xNum >= 8 { //! 闲天王
		trend = append(trend, 1)
	} else {
		trend = append(trend, 0)
	}

	if len(card) != len(_card) || trend[3] != 1 { //! 同点和
		trend = append(trend, 0)
	} else {
		sort.Ints(card)
		sort.Ints(_card)
		for i := 0; i < len(card); i++ {
			if card[i]/10 != _card[i]/10 {
				trend = append(trend, 0)
				break
			}
		}

		if len(trend) == 6 {
			trend = append(trend, 1)
		}
	}

	if zNum >= 8 { //! 庄天王
		trend = append(trend, 1)
	} else {
		trend = append(trend, 0)
	}

	return win, trend
}

func (self *Game_GoldBJL) OnBegin() {
	if self.room.IsBye() {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "房间已解散")
		return
	}
	self.room.Begin = true
	if self.Card == nil || self.Step == 50 {
		self.Card = NewCard_BJL156()
		self.Step = 0
		self.WinLost[0] = 0
		self.WinLost[1] = 0
	}

	var card *CardMgr
	lib.HF_DeepCopy(&card, &self.Card)
	self.Result[0] = card.Deal(2)
	self.Result[1] = card.Deal(2)
	self.Result[0], self.Result[1] = self.BoCard(self.Result[0], self.Result[1], card)

	if self.Dealer != nil {
		if self.Robot.RobotTotal == 0 { //! 没有机器人下注
			if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost == 102 { //! 平衡模式
				lst := make([]GameGoldBJL_CanResult, 0) //! 0-闲 1-庄
				winLst := make([]GameGoldBJL_CanResult, 0)
				lostLst := make([]GameGoldBJL_CanResult, 0)

				for i := 0; i < 50; i++ {
					var card *CardMgr
					lib.HF_DeepCopy(&card, &self.Card)
					xCard := card.Deal(2)
					zCard := card.Deal(2)

					xCard, zCard = self.BoCard(xCard, zCard, card)
					dealwin := self.GetDealWin(xCard, zCard)
					if GetServer().BJLUserMoney[self.room.Type%260000]+int64(dealwin) >= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMin && GetServer().BJLUserMoney[self.room.Type%260000]+int64(dealwin) <= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMax {
						lst = append(lst, GameGoldBJL_CanResult{xCard, zCard})
					}
					if dealwin >= 0 {
						winLst = append(winLst, GameGoldBJL_CanResult{xCard, zCard})
					}
					if dealwin <= 0 {
						lostLst = append(lostLst, GameGoldBJL_CanResult{xCard, zCard})
					}
				}

				if len(lst) == 0 {
					if GetServer().BJLUserMoney[self.room.Type%260000] >= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMax && len(lostLst) > 0 {
						lib.GetLogMgr().Output(lib.LOG_DEBUG, "玩家庄102 必输 随机lostlst")
						index := lib.HF_GetRandom(len(lostLst))
						self.Result[0] = lostLst[index].XCard
						self.Result[1] = lostLst[index].ZCard
					} else if GetServer().BJLUserMoney[self.room.Type%260000] <= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMin && len(winLst) > 0 {
						lib.GetLogMgr().Output(lib.LOG_DEBUG, "玩家庄102 必赢 随机winlst")
						index := lib.HF_GetRandom(len(winLst))
						self.Result[0] = winLst[index].XCard
						self.Result[1] = winLst[index].ZCard
					} else {
						lib.GetLogMgr().Output(lib.LOG_DEBUG, "玩家庄102 纯随机")
						var card *CardMgr
						lib.HF_DeepCopy(&card, &self.Card)
						self.Result[0] = card.Deal(2)
						self.Result[1] = card.Deal(2)
						self.Result[0], self.Result[1] = self.BoCard(self.Result[0], self.Result[1], card)
					}

				} else {
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "玩家庄102  随机lst")
					index := lib.HF_GetRandom(len(lst))
					self.Result[0] = lst[index].XCard
					self.Result[1] = lst[index].ZCard
				}
			} else if lib.GetManyMgr().GetProperty(self.room.Type%260000).PlayerCost > 100 { //! 纯随机

			} else { //! 设置概率
				lst := make([]GameGoldBJL_CanResult, 0) //! 0-闲 1-庄
				winLst := make([]GameGoldBJL_CanResult, 0)
				lostLst := make([]GameGoldBJL_CanResult, 0)

				for i := 0; i < 50; i++ {
					var card *CardMgr
					lib.HF_DeepCopy(&card, &self.Card)
					xCard := card.Deal(2)
					zCard := card.Deal(2)

					xCard, zCard = self.BoCard(xCard, zCard, card)
					dealwin := self.GetDealWin(xCard, zCard)
					if GetServer().BJLUserMoney[self.room.Type%260000]+int64(dealwin) >= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMin && GetServer().BJLUserMoney[self.room.Type%260000]+int64(dealwin) <= lib.GetManyMgr().GetProperty(self.room.Type%260000).PlayerMax {
						lst = append(lst, GameGoldBJL_CanResult{xCard, zCard})
					}
					if dealwin >= 0 {
						winLst = append(winLst, GameGoldBJL_CanResult{xCard, zCard})
					}
					if dealwin <= 0 {
						lostLst = append(lostLst, GameGoldBJL_CanResult{xCard, zCard})
					}
				}

				iswin := lib.HF_GetRandom(100) < lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost

				if iswin && len(winLst) > 0 {
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "玩家庄 概率模式 赢")
					index := lib.HF_GetRandom(len(winLst))
					self.Result[0] = winLst[index].XCard
					self.Result[1] = winLst[index].ZCard
				} else if !iswin && len(lostLst) > 0 {
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "玩家庄 概率模式 输")
					index := lib.HF_GetRandom(len(lostLst))
					self.Result[0] = lostLst[index].XCard
					self.Result[1] = lostLst[index].ZCard
				} else {
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "玩家庄 概率模式 纯随机")
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "系统庄  随机lst")
					index := lib.HF_GetRandom(len(lst))
					self.Result[0] = lst[index].XCard
					self.Result[1] = lst[index].ZCard
				}
			}
		} else {
			lst := make([]GameGoldBJL_CanResult, 0)
			for i := 0; i < 50; i++ {
				var card *CardMgr
				lib.HF_DeepCopy(&card, &self.Card)
				xCard := card.Deal(2)
				zCard := card.Deal(2)

				xCard, zCard = self.BoCard(xCard, zCard, card)
				robotwin := self.GetRobotWin(xCard, zCard)
				if lib.GetRobotMgr().GetRobotWin(self.room.Type)+robotwin >= 0 || robotwin >= 0 {
					lst = append(lst, GameGoldBJL_CanResult{xCard, zCard})
				}
			}
			if len(lst) != 0 {
				index := lib.HF_GetRandom(len(lst))
				self.Result[0] = lst[index].XCard
				self.Result[1] = lst[index].ZCard
			} else {
				var card *CardMgr
				lib.HF_DeepCopy(&card, &self.Card)
				self.Result[0] = card.Deal(2)
				self.Result[1] = card.Deal(2)
				self.Result[0], self.Result[1] = self.BoCard(self.Result[0], self.Result[1], card)
			}
		}
	} else {
		lst := make([]GameGoldBJL_CanResult, 0) //! 0-闲 1-庄
		winLst := make([]GameGoldBJL_CanResult, 0)
		lostLst := make([]GameGoldBJL_CanResult, 0)

		for i := 0; i < 25; i++ {
			var card *CardMgr
			lib.HF_DeepCopy(&card, &self.Card)
			xCard := card.Deal(2)
			zCard := card.Deal(2)

			xCard, zCard = self.BoCard(xCard, zCard, card)
			//			lib.GetLogMgr().Output(lib.LOG_DEBUG, "\\\\\\\\\\\\\\         xcard : ", xCard, " zcard : ", zCard)
			dealwin := self.GetDealWin(xCard, zCard)
			//			lib.GetLogMgr().Output(lib.LOG_DEBUG, "\\\\\\\\\\\\\\-------- xcard : ", xCard, " zcard : ", zCard)
			//			lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------------------------------------------")
			if GetServer().BJLSysMoney[self.room.Type%260000]+int64(dealwin) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && GetServer().BJLSysMoney[self.room.Type%260000]+int64(dealwin) <= lib.GetManyMgr().GetProperty(self.room.Type%260000).JackPotMax {
				lst = append(lst, GameGoldBJL_CanResult{xCard, zCard})
			}
			if dealwin >= 0 {
				winLst = append(winLst, GameGoldBJL_CanResult{xCard, zCard})
			} else {
				lostLst = append(lostLst, GameGoldBJL_CanResult{xCard, zCard})
			}
		}

		if len(lst) == 0 {
			if GetServer().BJLSysMoney[self.room.Type%260000] >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax && len(lostLst) > 0 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "系统庄 必输 随机lostlst")
				index := lib.HF_GetRandom(len(lostLst))
				self.Result[0] = lostLst[index].XCard
				self.Result[1] = lostLst[index].ZCard
			} else if GetServer().BJLSysMoney[self.room.Type%260000] <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && len(winLst) > 0 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "系统庄 必赢 随机winlst")
				index := lib.HF_GetRandom(len(winLst))
				self.Result[0] = winLst[index].XCard
				self.Result[1] = winLst[index].ZCard
			} else {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "系统庄 纯随机")
				var card *CardMgr
				lib.HF_DeepCopy(&card, &self.Card)
				self.Result[0] = card.Deal(2)
				self.Result[1] = card.Deal(2)
				self.Result[0], self.Result[1] = self.BoCard(self.Result[0], self.Result[1], card)
			}

		} else {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "系统庄  随机lst")
			index := lib.HF_GetRandom(len(lst))
			self.Result[0] = lst[index].XCard
			self.Result[1] = lst[index].ZCard
		}
	}

	for i := 0; i < len(self.Result); i++ {
		for j := 0; j < len(self.Result[i]); j++ {
			self.Card.DealCard(self.Result[i][j])
		}
	}

	self.OnEnd()
}

func (self *Game_GoldBJL) OnEnd() {
	self.room.Begin = false
	self.Step++

	//	{
	//		aa := []int{32, 51}
	//		bb := []int{42, 63}
	//		card := NewCard_BJL156()
	//		aa, bb = self.BoCard(aa, bb, card)
	//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "------------- aa : ", aa, "  bb : ", bb)

	//	}

	dealwin := 0
	robotwin := 0
	trend, ct := self.IsType(self.Result[0], self.Result[1])
	if trend == 0 {
		self.WinLost[0]++
	} else if trend == 1 {
		self.WinLost[1]++
	}
	xNum := 0
	zNum := 0
	for i := 0; i < len(self.Result[0]); i++ {
		if self.Result[0][i]/10 < 10 {
			xNum += self.Result[0][i] / 10
		}
	}
	xNum = xNum % 10

	for i := 0; i < len(self.Result[1]); i++ {
		if self.Result[1][i]/10 < 10 {
			zNum += self.Result[1][i] / 10
		}
	}
	zNum = zNum % 10

	var tre Msg_GameGoldBjl_Trend
	lib.HF_DeepCopy(&tre.Result, &self.Result)
	tre.Trend = trend
	tre.Num = append(tre.Num, xNum)
	tre.Num = append(tre.Num, zNum)
	if len(self.Trend) > 49 {
		self.Trend = self.Trend[len(self.Trend)-49 : len(self.Trend)]
	}
	self.Trend = append(self.Trend, tre)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "ct : ", ct)

	for i := 0; i < len(self.Bets); i++ {
		if ct[i] == 1 {
			for key, value := range self.Bets[i] {
				winmoney := value * GOLDBJL_BS[i]
				dealwin -= winmoney
				key.Win += winmoney
				key.Cost += int(math.Ceil(float64(winmoney-value) * float64(lib.GetManyMgr().GetProperty(self.room.Type).Cost) / 100.0))
			}

			for key, value := range self.Robot.RobotsBet[i] {
				winmoney := value * GOLDBJL_BS[i]
				key.AddWin(winmoney)
				key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
				robotwin += winmoney
				if self.Dealer != nil {
					dealwin -= winmoney
				}
			}
		}
		if ct[i] == 2 { //!　和，退庄闲钱
			for key, value := range self.Bets[i] {
				winmoney := value
				dealwin -= winmoney
				key.Win += winmoney
			}
			for key, value := range self.Robot.RobotsBet[i] {
				winmoney := value
				key.AddWin(winmoney)
				robotwin += winmoney
				if self.Dealer != nil {
					dealwin -= winmoney
				}
			}
		}
	}
	robotwin -= self.Robot.RobotTotal
	if self.Dealer != nil {
		dealwin += self.Robot.RobotTotal
	}

	dealwin += self.Total

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "RobotTotal=", self.Robot.RobotTotal)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "robotwin=", robotwin)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "dealwin=", dealwin)

	var bigwin *GameGold_BigWin = nil
	for _, value := range self.PersonMgr {
		if value.Win > 0 {
			value.Win -= value.Cost
			GetServer().SqlAgentGoldLog(value.Uid, value.Cost, self.room.Type)
			GetServer().SqlAgentBillsLog(value.Uid, value.Cost/2, self.room.Type)
			value.Total += value.Win

			var msg Msg_GameGoldBJL_Balance
			msg.Uid = value.Uid
			msg.Total = value.Total
			msg.Win = value.Win
			find := false
			for j := 0; j < len(self.Seat); j++ {
				if self.Seat[j].Person == value {
					self.room.broadCastMsg("gamegoldbjlbalance", &msg)
					find = true
					break
				}
			}
			if !find {
				self.room.SendMsg(value.Uid, "gamegoldbjlbalance", &msg)
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
			var record Rec_GoldBJL_Info
			record.Time = time.Now().Unix()
			record.GameType = self.room.Type
			var rec Son_Rec_GoldBJL_Person
			rec.Uid = value.Uid
			rec.Name = value.Name
			rec.Head = value.Head
			rec.Scroe = value.Win - value.Bets
			rec.Result = self.Result
			rec.Bets = value.BetInfo
			record.Info = append(record.Info, rec)
			GetServer().InsertRecord(self.room.Type, value.Uid, lib.HF_JtoA(&record), rec.Scroe)
		}
	}

	for i := 0; i < len(self.Robot.Robots); i++ {
		if self.Robot.Robots[i].GetWin() > 0 {
			self.Robot.Robots[i].AddWin(-self.Robot.Robots[i].GetCost())
			self.Robot.Robots[i].AddMoney(self.Robot.Robots[i].GetWin())

			for j := 0; j < len(self.Seat); j++ {
				if self.Seat[j].Robot == self.Robot.Robots[i] {
					var msg Msg_GameGoldBJL_Balance
					msg.Uid = self.Robot.Robots[i].Id
					msg.Total = self.Robot.Robots[i].GetMoney()
					msg.Win = self.Robot.Robots[i].GetWin()
					self.room.broadCastMsg("gamegoldbjlbalance", &msg)
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

	if self.Dealer == nil && dealwin != 0 {
		GetServer().SqlBZWLog(&SQL_BZWLog{1, dealwin, time.Now().Unix(), self.room.Type})
	}
	if self.Dealer != nil && robotwin != 0 { //! 玩家庄
		lib.GetRobotMgr().AddRobotWin(self.room.Type, robotwin)
		GetServer().SqlBZWLog(&SQL_BZWLog{1, robotwin, time.Now().Unix(), self.room.Type + 10000000})
	}
	if self.Dealer != nil && lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost == 102 {
		GetServer().SetBJLUserMoney(self.room.Type%260000, GetServer().BJLUserMoney[self.room.Type%260000]+int64(dealwin))
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
		GetServer().SetBJLSysMoney(self.room.Type%260000, GetServer().BJLSysMoney[self.room.Type%260000]+int64(_dealwin))
		dealwin -= robotwin
		if dealwin > 0 {
			bl := lib.GetManyMgr().GetProperty(self.room.Type).Cost
			cost := int(math.Ceil(float64(dealwin) * bl / 100.0))
			dealwin -= cost
		}
	}

	if self.Dealer != nil {
		var record Rec_GoldBJL_Info
		record.Time = time.Now().Unix()
		record.GameType = self.room.Type
		var rec Son_Rec_GoldBJL_Person
		rec.Uid = self.Dealer.Uid
		rec.Name = self.Dealer.Name
		rec.Head = self.Dealer.Head
		rec.Scroe = dealwin
		rec.Result = self.Result
		rec.Bets = self.Dealer.BetInfo
		record.Info = append(record.Info, rec)
		GetServer().InsertRecord(self.room.Type, self.Dealer.Uid, lib.HF_JtoA(&record), rec.Scroe)
	}

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "------------------ dealwin : ", dealwin)

	{
		//! 庄家信息
		var msg Msg_GameGoldBJL_Balance
		if self.Dealer != nil {
			if self.Dealer.Total+dealwin > 0 {
				self.Dealer.Total += dealwin
			} else {
				self.Dealer.Total = 0
				dealwin = -self.Dealer.Total
			}
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
		self.room.broadCastMsg("gamegoldbjlbalance", &msg)
	}

	//! 返回机器人
	self.Robot.Init(8, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)

	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 14
	self.SetTime(self.BetTime)

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

	{
		//! 总结算
		var msg Msg_GameGoldBJL_End
		msg.Result = self.Result
		msg.CT = ct
		if bigwin != nil {
			msg.Uid = bigwin.Uid
			msg.Name = bigwin.Name
			msg.Head = bigwin.Head
		}
		msg.Trend = self.Trend
		lib.HF_DeepCopy(&msg.Card, &self.Card.Card)
		msg.Step = self.Step
		msg.WinLost = self.WinLost
		msg.Money = lib.GetManyMoneyMgr().GetProperty(self.room.Type).Money
		msg.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime
		self.room.broadCastMsg("gamegoldbjlend", &msg)
	}
	self.Total = 0

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

		//! 走的人在上庄列表上
		for j := 0; j < len(self.LstDeal); j++ {
			if self.LstDeal[j].Person == value {
				copy(self.LstDeal[j:], self.LstDeal[j+1:])
				self.LstDeal = self.LstDeal[:len(self.LstDeal)-1]
				break
			}
		}

		//!　走的人是庄家
		if self.Dealer == value {
			self.ChageDeal()
		}

		//! 走的人是位置上的人
		for j := 0; j < len(self.Seat); j++ {
			if self.Seat[j].Person == value {
				self.Seat[j].Person = nil
				var msg Msg_GameGoldBjl_UpdSeat
				msg.Index = j
				self.room.broadCastMsg("gamegoldbjlseat", &msg)
				break
			}
		}
		delete(self.PersonMgr, key)
	}

	for i := 0; i < len(self.Bets); i++ {
		self.Bets[i] = make(map[*Game_GoldBJL_Person]int)
	}

	//! 判断坐下的人是否还能继续坐下
	for i := 0; i < len(self.Seat); i++ {
		if self.Seat[i].Person == nil {
			continue
		}

		if self.Seat[i].Person.Total < lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney {
			self.Seat[i].Person.Seat = -1
			var msg Msg_GameGoldBjl_UpdSeat
			msg.Index = i
			self.room.broadCastMsg("gamegoldbjlseat", &msg)
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
			var msg Msg_GameGoldBZW_UpdSeat
			msg.Index = i
			self.room.broadCastMsg("gamegoldbjlseat", &msg)
			self.Seat[i].Robot = nil
		}
	}
}

func (self *Game_GoldBJL) OnInit(room *Room) {
	self.room = room
	self.Money = lib.GetManyMgr().GetProperty(self.room.Type).DealInitMoney
	self.Card = NewCard_BJL156()
	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 14
	//! 载入机器人
	self.Robot.Init(8, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)
}

func (self *Game_GoldBJL) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldBJL) OnSendInfo(person *Person) {
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
		person.SendMsg("gamegoldbjlinfo", self.getInfo(person.Uid, value.Total))
		return
	}

	_person := new(Game_GoldBJL_Person)
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
	person.SendMsg("gamegoldbjlinfo", self.getInfo(person.Uid, person.Gold))
}

func (self *Game_GoldBJL) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold":
		person := self.GetPerson(msg.V.(*staticfunc.Msg_SynchroGold).Uid)
		if person != nil {
			person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal(person.Uid, person.Total)
		}
	case "gamebzwgoon":
		self.GameGoOn(msg.Uid)
	case "gamebzwbets":
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

func (self *Game_GoldBJL) OnBye() {

}

func (self *Game_GoldBJL) OnExit(uid int64) {
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

func (self *Game_GoldBJL) OnIsDealer(uid int64) bool {
	if self.Dealer != nil && self.Dealer == self.GetPerson(uid) {
		return true
	}
	return false
}

func (self *Game_GoldBJL) OnIsBets(uid int64) bool {
	value, ok := self.PersonMgr[uid]
	if ok {
		return value.Bets > 0
	}
	return false
}

func (self *Game_GoldBJL) OnBalance() {
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

func (self *Game_GoldBJL) OnTime() {
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
		self.RobotSeat(lib.HF_GetRandom(6), self.Robot.Robots[i])
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

			index, gold, _ := self.Robot.GameBets(self.Robot.Robots[i])
			if gold == 0 {
				continue
			}
			if self.Dealer != nil { //! 玩家庄判断是否能下
				dealwin := self.GetMaxLost(true, index != 2 && index != 4)
				if self.Dealer.Total+self.Total+self.Robot.RobotTotal-dealwin < self.Dealer.Total/5 {
					self.Robot.GameBackBets(self.Robot.Robots[i], index, gold)
					continue
				}
			} else if self.RobotDealer != nil {
				dealwin := self.GetMaxLost(true, index != 2 && index != 4)
				if self.RobotDealer.GetMoney()+self.Total+self.Robot.RobotTotal-dealwin < self.RobotDealer.GetMoney()/5 {
					self.Robot.GameBackBets(self.Robot.Robots[i], index, gold)
					continue
				}
			} else if lib.GetManyMgr().GetProperty(self.room.Type).DealChange == 1 {
				dealwin := self.GetMaxLost(true, index != 2 && index != 4)
				if self.Money+self.Total+self.Robot.RobotTotal-dealwin < self.Money/5 {
					self.Robot.GameBackBets(self.Robot.Robots[i], index, gold)
					continue
				}
			}
			var msg Msg_GameGoldBJL_Bets
			msg.Uid = self.Robot.Robots[i].Id
			msg.Index = index
			msg.Gold = gold
			msg.Total = self.Robot.Robots[i].GetMoney()
			self.room.broadCastMsg("gamegoldbjlbets", &msg)
		}
		return
	}

	if !self.room.Begin {
		self.OnBegin()
	}
}

//! 得到这个位置机器人下了多少钱
func (self *Game_GoldBJL) GetMoneyPosByRobot(index int) int {
	total := 0
	for _, value := range self.Robot.RobotsBet[index] {
		total += value
	}
	return total
}

//! 得到机器人可以赢的钱
func (self *Game_GoldBJL) GetRobotWin(xcard []int, zcard []int) int {
	win := 0
	_, ct := self.IsType(xcard, zcard)
	for i := 0; i < len(self.Bets); i++ {
		if ct[i] == 1 {
			for _, value := range self.Robot.RobotsBet[i] {
				winmoney := value * GOLDBJL_BS[i]
				win += winmoney
			}
		}
		if ct[i] == 2 {
			for _, value := range self.Robot.RobotsBet[i] {
				winmoney := value
				win += winmoney
			}
		}
	}

	return win - self.Robot.RobotTotal
}
