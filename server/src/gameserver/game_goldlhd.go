package gameserver

import (
	"fmt"
	"lib"
	"math"
	"staticfunc"
	"time"
)

//! 金币场记录
type Rec_LHD_Info struct {
	GameType int                  `json:"gametype"`
	Time     int64                `json:"time"` //! 记录时间
	Info     []Son_Rec_LHD_Person `json:"info"`
}
type Son_Rec_LHD_Person struct {
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`
	Head   string `json:"head"`
	Score  int    `json:"score"`
	Result [2]int `json:"result"`
	Bets   [3]int `json:"bets"`
}

type Game_GoldLHDSeat struct {
	Person *Game_GoldLHD_Person
	Robot  *lib.Robot
}

func (self *Game_GoldLHDSeat) GetTotal() int {
	if self.Person != nil {
		return self.Person.Total
	} else if self.Robot != nil {
		return self.Robot.GetMoney()
	}
	return 0
}

type Game_GoldLHD struct {
	PersonMgr   map[int64]*Game_GoldLHD_Person  `json:"personmgr"`
	Bets        [3]map[*Game_GoldLHD_Person]int `json:"bets"`
	Dealer      *Game_GoldLHD_Person            `json:"dealer"`      //! 庄家
	RobotDealer *lib.Robot                      `json:"robotdealer"` //! 机器人庄
	Round       int                             `json:"round"`       //! 连庄轮数
	DownUid     int64                           `json:"downuid"`     //! 下庄的人
	Result      [2]int                          `json:"result"`      //! 0-龙 1-虎
	LstDeal     []Game_GoldLHDSeat              `json:"lstdeal"`     //! 上庄列表
	Time        int64                           `json:"time"`
	Seat        [8]Game_GoldLHDSeat             `json:"seat"`
	Total       int                             `json:"total"` //! 这局一共下了多少钱
	Money       int                             `json:"money"` //! 系统庄的钱
	Trend       []int                           `json:"trend"` //! 走势
	Tmp         int                             `json:"tmp"`   //! 1-下注时间 2-开奖时间
	Next        []int                           `json:"next"`
	BetTime     int                             `json:"bettime"`
	Robot       lib.ManyGameRobot               //! 机器人结构

	room *Room
}

func NewGame_GoldLHD() *Game_GoldLHD {
	game := new(Game_GoldLHD)
	game.PersonMgr = make(map[int64]*Game_GoldLHD_Person)
	game.Money = 10000000
	for i := 0; i < 48; i++ {
		cardmgr := NewCard_LYC()
		game.Result[0] = cardmgr.Deal(1)[0]
		game.Result[1] = cardmgr.Deal(1)[0]
		game.Trend = append(game.Trend, game.IsType())
	}
	for i := 0; i < len(game.Bets); i++ {
		game.Bets[i] = make(map[*Game_GoldLHD_Person]int)
	}
	game.Tmp = 1
	return game
}

type Game_GoldLHD_Person struct {
	Uid       int64  `json:"uid"`
	Gold      int    `json:"gold"`      //! 进房间的时候有多少钱
	Total     int    `json:"total"`     //! 当前有多少钱
	Win       int    `json:"win"`       //! 本局赢了多少
	Cost      int    `json:"cost"`      //! 抽水钱
	Bets      int    `json:"bets"`      //! 本局下了多少钱
	BetInfo   [3]int `json:"betinfo"`   //! 本局下的注	0-龙 1-虎 2-和
	BeBets    int    `json:"bebets"`    //! 上局下来多少注
	BeBetInfo [3]int `json:"bebetinfo"` //! 上局下的注
	Name      string `json:"name"`      //! 名字
	Head      string `json:"head"`      //! 头像
	Online    bool   `json:"online"`
	Round     int    `json:"round"` //! 不下注轮数
	Seat      int    `json:"seat"`  //! 0-7有座 -1无座 100庄家
	IP        string `json:"ip"`
	Address   string `json:"address"`
	Sex       int    `json:"sex"`
}

type Msg_GameGoldLHD_Info struct {
	Begin   bool                    `json:"begin"`  //! 是否开始游戏
	Time    int64                   `json:"time"`   //! 倒计时
	Seat    [8]Son_GameGoldLHD_Info `json:"info"`   //!  8个位置
	Bets    [3]int                  `json:"bets"`   //! 三个下注
	Total   int                     `json:"total"`  //! 自己的钱
	Trend   []int                   `json:"trend"`  //! 走势
	Tmp     int                     `json:"tmp"`    //! 1-下注时间 2-开奖时间
	Change  bool                    `json:"change"` //! 是否可以打开超端
	Money   []int                   `json:"money"`
	IsDeal  bool                    `json:"isdeal"` //! 是否可下庄
	Dealer  Son_GameGoldLHD_Info    `json:"dealer"` //! 庄家
	BetTime int                     `json:"bettime"`
}
type Son_GameGoldLHD_Info struct {
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Total   int    `json:"total"`
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type Msg_GameGoldLHD_Balance struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"` //! 当前金币
	Win   int   `json:"win"`   //! 赢了多少金币
}

type Msg_GameGoldLHD_End struct {
	Uid     int64  `json:"uid"` //! 大赢家
	Name    string `json:"name"`
	Head    string `json:"head"`
	Result  [2]int `json:"result"`
	Money   []int  `json:"money"`
	BetTime int    `json:"bettime"`
}

//! 刷新座位
type Msg_GameGoldLHD_UpdSeat struct {
	Index   int    `json:"index"`
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Total   int    `json:"total"`
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type Msg_GameGoldLHD_List struct {
	Info []Son_GameGoldLHD_Info `json:"info"`
}

type Msg_GameGoldLHD_Total struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"`
}

type Msg_GameGoldLHD_Bets struct {
	Uid   int64 `json:"uid"`
	Index int   `json:"index"`
	Gold  int   `json:"gold"`
	Total int   `json:"total"`
}

type Msg_GameGoldLHD_GoOn struct {
	Uid   int64  `json:"uid"`
	Gold  [3]int `json:"gold"`
	Total int    `json:"total"`
}

func (self *Game_GoldLHD) getinfo(uid int64, total int) *Msg_GameGoldLHD_Info {
	var msg Msg_GameGoldLHD_Info
	msg.Begin = self.room.Begin
	msg.Time = self.Time - time.Now().Unix()
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
	if GetServer().IsAdmin(uid, staticfunc.ADMIN_GOLDLHD) {
		msg.Change = true
	} else {
		msg.Change = false
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

func (self *Game_GoldLHD) GetPerson(uid int64) *Game_GoldLHD_Person {
	return self.PersonMgr[uid]
}

//! 得到这个位置下了多少钱
func (self *Game_GoldLHD) GetMoneyPos(index int, robot bool) int {
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
func (self *Game_GoldLHD) GetMoneyPosByRobot(index int) int {
	total := 0
	for _, value := range self.Robot.RobotsBet[index] {
		total += value
	}
	return total
}

//! 同步金币
func (self *Game_GoldLHD_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

//! 同步总分
func (self *Game_GoldLHD) SendTotal(uid int64, total int) {
	var msg Msg_GameGoldLHD_Total
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

//! 是否下注
func (self *Game_GoldLHD) OnIsBets(uid int64) bool {
	value, ok := self.PersonMgr[uid]
	if ok {
		return value.Bets > 0
	}
	return false
}

//! 设置时间
func (self *Game_GoldLHD) SetTime(t int) {
	if t == 0 {
		self.Time = 0
	} else {
		self.Time = time.Now().Unix() + int64(t)
	}

	//self.Tmp = tmp

	var msg Msg_SetTime
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	self.room.broadCastMsg("gametime", &msg)
}

//! 申请无座玩家列表
func (self *Game_GoldLHD) GamePlayerList(uid int64) {
	var msg Msg_GameGoldLHD_List
	tmp := make(map[int64]Son_GameGoldLHD_Info)
	for _, value := range self.PersonMgr {
		if value.Seat >= 0 {
			continue
		}

		var node Son_GameGoldLHD_Info
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

		var node Son_GameGoldLHD_Info
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

func (self *Game_GoldLHD) GameUpDeal(uid int64) {
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
		self.LstDeal = append(self.LstDeal, Game_GoldLHDSeat{person, nil})
	}
	person.Round = 0

	var msg Msg_GameGoldBZW_DealList
	msg.Type = 0
	msg.Info = make([]Son_GameGoldBZW_Info, 0)
	for i := 0; i < len(self.LstDeal); i++ {
		if self.LstDeal[i].Person != nil {
			msg.Info = append(msg.Info, Son_GameGoldBZW_Info{self.LstDeal[i].Person.Uid, self.LstDeal[i].Person.Name, self.LstDeal[i].Person.Head, self.LstDeal[i].Person.Total, self.LstDeal[i].Person.IP, self.LstDeal[i].Person.Address, self.LstDeal[i].Person.Sex})
		} else if self.LstDeal[i].Robot != nil {
			msg.Info = append(msg.Info, Son_GameGoldBZW_Info{self.LstDeal[i].Robot.Id, self.LstDeal[i].Robot.Name, self.LstDeal[i].Robot.Head, self.LstDeal[i].Robot.GetMoney(), self.LstDeal[i].Robot.IP, self.LstDeal[i].Robot.Address, self.LstDeal[i].Robot.Sex})
		}
	}
	self.room.SendMsg(uid, "gamelhddeal", &msg)
}

//! 机器人上庄
func (self *Game_GoldLHD) RobotUpDeal(robot *lib.Robot) {
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
	self.LstDeal = append(self.LstDeal, Game_GoldLHDSeat{nil, robot})
}

func (self *Game_GoldLHD) GameReDeal(uid int64) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if self.Dealer == person { //! 正在庄
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

	var msg Msg_GameGoldBZW_DealList
	msg.Type = 1
	msg.Info = make([]Son_GameGoldBZW_Info, 0)
	self.room.SendMsg(uid, "gamelhddeal", &msg)
}

//! 坐下
func (self *Game_GoldLHD) GameSeat(uid int64, index int) {
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

	var msg Msg_GameGoldLHD_UpdSeat
	msg.Uid = uid
	msg.Index = index
	msg.Head = person.Head
	msg.Name = person.Name
	msg.Total = person.Total
	msg.IP = person.IP
	msg.Address = person.Address
	msg.Sex = person.Sex
	self.room.broadCastMsg("gamelhdseat", &msg)
}

//! 机器人坐下
func (self *Game_GoldLHD) RobotSeat(index int, robot *lib.Robot) {
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

	var msg Msg_GameGoldLHD_UpdSeat
	msg.Uid = robot.Id
	msg.Index = index
	msg.Head = robot.Head
	msg.Name = robot.Name
	msg.Total = robot.GetMoney()
	msg.IP = robot.IP
	msg.Address = robot.Address
	msg.Sex = robot.Sex
	self.room.broadCastMsg("gamelhdseat", &msg)
}

//! 下注
func (self *Game_GoldLHD) GameBets(uid int64, index int, gold int) {
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
		self.room.SendErr(uid, "正在开奖，请稍后下注。")
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

	//! 判断庄家是否够赔
	{
		//! 模拟总下注
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
		//! 模拟庄家赢钱
		if dealmoney > 0 {
			dealwin := self.GetDealWinByPos(index, true)
			if dealmoney+self.Total+self.Robot.RobotTotal-dealwin < dealmoney/5 {
				self.Total -= gold
				self.Bets[index][person] -= gold
				self.room.SendErr(uid, "庄家金币不足，该位置无法下注。")
				return
			}
		}
	}

	person.Bets += gold
	person.Total -= gold
	person.BetInfo[index] += gold
	person.Round = 0

	var msg Msg_GameGoldLHD_Bets
	msg.Uid = uid
	msg.Index = index
	msg.Gold = gold
	msg.Total = person.Total
	self.room.broadCastMsg("gamelhdbets", &msg)
}

//! 续压
func (self *Game_GoldLHD) GameGoOn(uid int64) {
	if uid == 0 {
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() >= int64(self.BetTime-2) {
		self.room.SendErr(uid, "正在开奖，请稍后下注。")
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

	//! 判断庄家是否够赔
	{
		//! 模拟总下注
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
		//! 模拟庄家赢钱
		if dealmoney > 0 {
			dealwin := 0
			for i := 0; i < len(person.BeBetInfo); i++ {
				if person.BeBetInfo[i] == 0 {
					continue
				}
				tmp := self.GetDealWinByPos(i, true)
				if tmp > dealwin {
					dealwin = tmp
				}
			}
			if dealmoney+self.Total+self.Robot.RobotTotal-dealwin < dealmoney/5 {
				for i := 0; i < len(person.BeBetInfo); i++ {
					self.Total -= person.BeBetInfo[i]
					self.Bets[i][person] -= person.BeBetInfo[i]
				}
				self.room.SendErr(uid, "庄家金币不足，无法续压。")
				return
			}
		}
	}

	person.Bets += person.BeBets
	person.Total -= person.BeBets
	for i := 0; i < len(person.BeBetInfo); i++ {
		person.BetInfo[i] += person.BeBetInfo[i]
	}
	person.Round = 0

	var msg Msg_GameGoldLHD_GoOn
	msg.Uid = uid
	msg.Gold = person.BeBetInfo
	msg.Total = person.Total
	self.room.broadCastMsg("gamelhdgoon", &msg)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------------uid : ", uid)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------------Gold : ", msg.Gold)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------------Total : ", msg.Total)

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "self.bet", self.Bets)
}

//! 换庄
func (self *Game_GoldLHD) ChageDeal() {
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
						var msg Msg_GameGoldBZW_UpdSeat
						msg.Index = i
						self.room.broadCastMsg("gamelhdseat", &msg)
						self.Seat[i].Person = nil
						break
					}
				}
			} else if self.LstDeal[0].Robot != nil {
				self.RobotDealer = self.LstDeal[0].Robot
				self.RobotDealer.SetSeat(100)
				for i := 0; i < len(self.Seat); i++ {
					if self.Seat[i].Robot == self.RobotDealer {
						var msg Msg_GameGoldBZW_UpdSeat
						msg.Index = i
						self.room.broadCastMsg("gamelhdseat", &msg)
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

func (self *Game_GoldLHD) IsType() int { //! 0-龙 1-虎 2-和
	if self.Result[0]/10 == self.Result[1]/10 {
		return 2
	}
	if self.Result[0]/10 > self.Result[1]/10 {
		return 0
	}
	return 1
}

func (self *Game_GoldLHD) GetDealWinByPos(index int, robot bool) int {
	result := self.GetMoneyPos(index, robot)
	if index == 0 || index == 1 {
		return 2 * result
	} else {
		return 9 * result
	}
}

//! 庄家必赢或必输
func (self *Game_GoldLHD) WinOrLost(win bool) {
	result := make([]int, 3)
	for i := 0; i < 3; i++ {
		result[i] = self.GetMoneyPos(i, false)
	}
	longwin := result[1] + result[2] - result[0]
	huwin := result[0] + result[2] - result[1]

	trend := self.IsType()

	if trend == 2 { //! 和局可以开就直接开
		if !win { //! 必输可以开
			return
		} else {
			for {
				cardmgr := NewCard_LYC()
				self.Result[0] = cardmgr.Deal(1)[0]
				self.Result[1] = cardmgr.Deal(1)[0]
				trend = self.IsType()
				if trend != 2 {
					break
				}
			}
		}
	}

	canret := make([]int, 0)
	if win { //! 必赢
		if longwin >= 0 {
			canret = append(canret, 0)
		}
		if huwin >= 0 {
			canret = append(canret, 1)
		}
		if len(canret) == 0 {
			if longwin > huwin {
				canret = append(canret, 0)
			} else if longwin < huwin {
				canret = append(canret, 1)
			} else {
				return
			}
		}
	} else { //! 必输
		if longwin <= 0 {
			canret = append(canret, 0)
		}
		if huwin <= 0 {
			canret = append(canret, 1)
		}
		if len(canret) == 0 {
			if longwin > huwin {
				canret = append(canret, 1)
			} else if longwin < huwin {
				canret = append(canret, 0)
			} else {
				return
			}
		}
	}

	ret := canret[lib.HF_GetRandom(len(canret))]
	if ret == trend {
		return
	}
	lib.GetLogMgr().Output(lib.LOG_ERROR, "龙虎斗改变")
	tmp := self.Result[0]
	self.Result[0] = self.Result[1]
	self.Result[1] = tmp
}

//! 改变通过机器人下注
func (self *Game_GoldLHD) RobotNeedWin() {
	result := make([]int, 3)
	for i := 0; i < 3; i++ {
		result[i] = self.GetMoneyPosByRobot(i)
	}
	longwin := result[0]*2 - self.Robot.RobotTotal
	huwin := result[1]*2 - self.Robot.RobotTotal

	trend := self.IsType()

	if trend == 2 { //! 和局可以开就直接开
		return
	}

	canret := make([]int, 0)
	if longwin >= 0 || lib.GetRobotMgr().GetRobotWin(self.room.Type)+longwin >= 0 {
		canret = append(canret, 0)
	}
	if huwin >= 0 || lib.GetRobotMgr().GetRobotWin(self.room.Type)+huwin >= 0 {
		canret = append(canret, 1)
	}
	if len(canret) == 0 {
		return
	}

	ret := canret[lib.HF_GetRandom(len(canret))]
	if ret == trend {
		return
	}
	lib.GetLogMgr().Output(lib.LOG_ERROR, "龙虎斗改变")
	tmp := self.Result[0]
	self.Result[0] = self.Result[1]
	self.Result[1] = tmp
}

func (self *Game_GoldLHD) Change() {
	if self.Dealer == nil { //! 系统庄,只计算玩家的下注
		result := make([]int, 3)
		for i := 0; i < 3; i++ {
			result[i] = self.GetMoneyPos(i, false)
		}
		longwin := result[1] + result[2] - result[0]
		huwin := result[0] + result[2] - result[1]
		hewin := -result[2] * 8

		trend := self.IsType()

		if trend == 2 { //! 和局可以开就直接开
			if GetServer().LHDMoney[self.room.Type%100000]+int64(hewin) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
				return
			} else {
				for {
					cardmgr := NewCard_LYC()
					self.Result[0] = cardmgr.Deal(1)[0]
					self.Result[1] = cardmgr.Deal(1)[0]
					trend = self.IsType()
					if trend != 2 {
						break
					}
				}
			}
		}

		canret := make([]int, 0)
		if GetServer().LHDMoney[self.room.Type%100000] <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin { //! 必赢
			if longwin >= 0 {
				canret = append(canret, 0)
			}
			if huwin >= 0 {
				canret = append(canret, 1)
			}
			if len(canret) == 0 {
				if longwin > huwin {
					canret = append(canret, 0)
				} else if longwin < huwin {
					canret = append(canret, 1)
				} else {
					return
				}
			}
		} else if GetServer().LHDMoney[self.room.Type%100000] >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax { //! 必输
			if longwin <= 0 {
				canret = append(canret, 0)
			}
			if huwin <= 0 {
				canret = append(canret, 1)
			}
			if len(canret) == 0 {
				if longwin > huwin {
					canret = append(canret, 1)
				} else if longwin < huwin {
					canret = append(canret, 0)
				} else {
					return
				}
			}
		} else {
			return
		}

		ret := canret[lib.HF_GetRandom(len(canret))]
		if ret == trend {
			return
		}
		lib.GetLogMgr().Output(lib.LOG_ERROR, "龙虎斗改变")
		tmp := self.Result[0]
		self.Result[0] = self.Result[1]
		self.Result[1] = tmp
	} else { //! 玩家庄
		if self.Robot.RobotTotal == 0 { //! 机器人没有下注
			if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost <= 100 {
				self.WinOrLost(lib.HF_GetRandom(100) < lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost)
			}
		} else { //! 机器人下注了
			self.RobotNeedWin()
		}
	}
}

func (self *Game_GoldLHD) OnBegin() {
	if self.room.IsBye() {
		return
	}
	self.room.Begin = true

	cardmgr := NewCard_LYC()
	tmp1 := cardmgr.Deal(1) //! 龙
	tmp2 := cardmgr.Deal(1) //! 虎
	self.Result[0] = tmp1[0]
	self.Result[1] = tmp2[0]

	self.OnEnd()
}

func (self *Game_GoldLHD) OnEnd() {
	self.room.Begin = false

	if len(self.Next) == 1 {
		loop := 0
		for self.Next[0] != self.IsType() && loop < 100000 {
			cardmgr := NewCard_LYC()
			self.Result[0] = cardmgr.Deal(1)[0]
			self.Result[1] = cardmgr.Deal(1)[0]
			loop++
		}
		self.Next = make([]int, 0)
	} else {
		self.Change()
	}

	trend := self.IsType()
	tmp := make([]int, 0)
	tmp = append(tmp, trend)
	tmp = append(tmp, self.Trend...)
	if len(tmp) > 48 {
		tmp = tmp[0:48]
	}
	self.Trend = tmp

	dealwin := 0
	robotwin := 0

	for i := 0; i < len(self.Bets); i++ {
		if trend == 2 && i != 2 { //! 和局
			for key, value := range self.Bets[i] {
				key.Win += value
				dealwin -= value
			}

			for key, value := range self.Robot.RobotsBet[i] {
				key.AddWin(value)
				robotwin += value
				if self.Dealer != nil {
					dealwin -= value
				}
			}
		}

		if trend == i {
			for key, value := range self.Bets[i] {
				if i == 2 {
					winmoney := value * 9
					dealwin -= winmoney
					key.Win += winmoney
					key.Cost = int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
				} else {
					winmoney := value * 2
					dealwin -= winmoney
					key.Win += winmoney
					key.Cost = int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
				}
			}

			for key, value := range self.Robot.RobotsBet[i] {
				if i == 2 {
					winmoney := value * 9
					key.AddWin(winmoney)
					key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
					robotwin += winmoney
					if self.Dealer != nil {
						dealwin -= winmoney
					}
				} else {
					winmoney := value * 2
					key.AddWin(winmoney)
					key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
					robotwin += winmoney
					if self.Dealer != nil {
						dealwin -= winmoney
					}
				}
			}
		}
	}
	robotwin -= self.Robot.RobotTotal
	if self.Dealer != nil {
		dealwin += self.Robot.RobotTotal
	}

	var bigwin *GameGold_BigWin = nil //! 大赢家
	for _, value := range self.PersonMgr {
		if value.Win > 0 {
			if value.Win > 0 && value.Cost != 0 {
				value.Win -= value.Cost
				GetServer().SqlAgentGoldLog(value.Uid, value.Cost, self.room.Type)
				GetServer().SqlAgentBillsLog(value.Uid, value.Cost/2, self.room.Type)
			}
			value.Total += value.Win

			var msg Msg_GameGoldLHD_Balance
			msg.Uid = value.Uid
			msg.Win = value.Win
			msg.Total = value.Total
			find := false
			for j := 0; j < len(self.Seat); j++ {
				if self.Seat[j].Person == value {
					self.room.broadCastMsg("gamegoldlhdbalance", &msg)
					find = true
					break
				}
			}
			if !find {
				self.room.SendMsg(value.Uid, "gamegoldlhdbalance", &msg)
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
		//! 插入战绩
		if value.Bets > 0 {
			var record Rec_LHD_Info
			record.Time = time.Now().Unix()
			record.GameType = self.room.Type
			var rec Son_Rec_LHD_Person
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
					var msg Msg_GameGoldLHD_Balance
					msg.Uid = self.Robot.Robots[i].Id
					msg.Total = self.Robot.Robots[i].GetMoney()
					msg.Win = self.Robot.Robots[i].GetWin()
					self.room.broadCastMsg("gamegoldlhdbalance", &msg)
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

	dealwin = self.Total + dealwin
	if self.Dealer == nil && dealwin != 0 { //! 系统庄
		GetServer().SqlBZWLog(&SQL_BZWLog{1, dealwin, time.Now().Unix(), self.room.Type})
	}
	if self.Dealer != nil && robotwin != 0 { //! 玩家庄
		lib.GetRobotMgr().AddRobotWin(self.room.Type, robotwin)
		GetServer().SqlBZWLog(&SQL_BZWLog{1, robotwin, time.Now().Unix(), self.room.Type + 10000000})
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
		GetServer().SetLHDMoney(self.room.Type%100000, GetServer().LHDMoney[self.room.Type%100000]+int64(_dealwin))
		dealwin -= robotwin
		if dealwin > 0 {
			bl := lib.GetManyMgr().GetProperty(self.room.Type).Cost
			cost := int(math.Ceil(float64(dealwin) * bl / 100.0))
			dealwin -= cost
		}
	}

	if self.Dealer != nil {
		var record Rec_LHD_Info
		record.Time = time.Now().Unix()
		record.GameType = self.room.Type
		var rec Son_Rec_LHD_Person
		rec.Uid = self.Dealer.Uid
		rec.Name = self.Dealer.Name
		rec.Head = self.Dealer.Head
		rec.Score = dealwin
		rec.Result = self.Result
		record.Info = append(record.Info, rec)
		GetServer().InsertRecord(self.room.Type, self.Dealer.Uid, lib.HF_JtoA(&record), rec.Score)
	}

	//! 发送庄家结算
	{
		var msg Msg_GameGoldBZW_Balance
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
		self.room.broadCastMsg("gamegoldlhdbalance", &msg)
	}

	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 9
	self.SetTime(self.BetTime)

	//! 大赢家
	{
		var msg Msg_GameGoldLHD_End
		msg.Result = self.Result
		msg.Money = lib.GetManyMoneyMgr().GetProperty(self.room.Type).Money
		msg.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime
		if bigwin != nil {
			msg.Uid = bigwin.Uid
			msg.Name = bigwin.Name
			msg.Head = bigwin.Head
		}
		self.room.broadCastMsg("gamegoldlhdend", &msg)
	}
	self.Total = 0

	//! 把不在room.uid里面的玩家清理出去
	for key, value := range self.PersonMgr {
		if value.Online {
			find := false
			for i := 0; i < len(self.LstDeal); i++ {
				if self.LstDeal[i].Person == value {
					find = true
					break
				}
			}
			if !find && value.Seat < 0 { //! 无座玩家不下注轮数++
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

		//! 走的人正在上庄列表
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

		//! 走的人是位置上面的人
		for j := 0; j < len(self.Seat); j++ {
			if self.Seat[j].Person == value {
				self.Seat[j].Person = nil
				var msg Msg_GameGoldBZW_UpdSeat
				msg.Index = j
				self.room.broadCastMsg("gamelhdseat", &msg)
				break
			}
		}
		delete(self.PersonMgr, key)
	}

	//! 载入机器人
	self.Robot.Init(3, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)

	for i := 0; i < len(self.Bets); i++ {
		self.Bets[i] = make(map[*Game_GoldLHD_Person]int)
	}

	//! 判断庄家是否能继续连
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
		if self.RobotDealer.GetMoney() < lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney {
			self.ChageDeal()
		} else {
			if self.Round >= lib.HF_GetRandom(6)+3 && len(self.LstDeal) > 0 || !lib.GetRobotMgr().GetRobotSet(self.room.Type).NeedRobot {
				self.ChageDeal()
			} else {
				self.Round++
			}
		}
	} else if len(self.LstDeal) > 0 {
		self.ChageDeal()
	}

	//! 判断坐下的玩家是否还能继续坐下
	for i := 0; i < len(self.Seat); i++ {
		if self.Seat[i].Person == nil {
			continue
		}
		if self.Seat[i].Person.Total < lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney {
			self.Seat[i].Person.Seat = -1
			var msg Msg_GameGoldLHD_UpdSeat
			msg.Index = i
			self.room.broadCastMsg("gamelhdseat", &msg)
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
			var msg Msg_GameGoldLHD_UpdSeat
			msg.Index = i
			self.room.broadCastMsg("gamelhdseat", &msg)
			self.Seat[i].Robot = nil
		}
	}
}

func (self *Game_GoldLHD) OnInit(room *Room) {
	self.room = room
	self.Money = lib.GetManyMgr().GetProperty(self.room.Type).DealInitMoney
	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 9

	//! 载入机器人
	self.Robot.Init(3, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)
}

func (self *Game_GoldLHD) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldLHD) OnSendInfo(person *Person) {
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
		person.SendMsg("gamegoldlhdinfo", self.getinfo(person.Uid, value.Total))
		return
	}

	_person := new(Game_GoldLHD_Person)
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
	person.SendMsg("gamegoldlhdinfo", self.getinfo(person.Uid, person.Gold))
}

func (self *Game_GoldLHD) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold": //! 同步金币
		person := self.GetPerson(msg.V.(*staticfunc.Msg_SynchroGold).Uid)
		if person != nil {
			person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal(person.Uid, person.Total)
		}
	case "gamerob": //! 上庄
		self.GameUpDeal(msg.Uid)
	case "gameredeal": //! 下庄
		self.GameReDeal(msg.Uid)
	case "gamebzwbets": //! 下注
		self.GameBets(msg.Uid, msg.V.(*Msg_GameGoldBZW_Bets).Index, msg.V.(*Msg_GameGoldBZW_Bets).Gold)
	case "gamebzwgoon":
		self.GameGoOn(msg.Uid)
	case "gamebzwseat":
		self.GameSeat(msg.Uid, msg.V.(*Msg_GameGoldBZW_Seat).Index)
	case "gameplayerlist":
		self.GamePlayerList(msg.Uid)
	case "gamesetnext":
		self.Next = msg.V.(*staticfunc.Msg_SetDealNext).Next
	case "gamechange":
		self.ChangeCard(msg.Uid, msg.V.(*Msg_GameChange).Card)
	}
}

func (self *Game_GoldLHD) ChangeCard(uid int64, card []int) {
	if !GetServer().IsAdmin(uid, staticfunc.ADMIN_GOLDLHD) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "非超端玩家")
		return
	}

	if len(card) != 1 {
		return
	}

	if card[0] < 0 && card[0] > 2 {
		return
	}
	lib.HF_DeepCopy(&self.Next, card)
	GetServer().SqlSuperClientLog(&SQL_SuperClientLog{1, uid, self.room.Type, lib.HF_JtoA(&card), time.Now().Unix()})
	self.room.SendMsg(uid, "ok", nil)
}

func (self *Game_GoldLHD) OnBye() {

}

func (self *Game_GoldLHD) OnExit(uid int64) {
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

func (self *Game_GoldLHD) OnIsDealer(uid int64) bool {
	if self.Dealer != nil && self.Dealer == self.GetPerson(uid) {
		return true
	}
	return false
}

func (self *Game_GoldLHD) OnBalance() { //! 结算
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

func (self *Game_GoldLHD) OnTime() {
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
		self.RobotSeat(lib.HF_GetRandom(8), self.Robot.Robots[i])
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
				dealwin := self.GetDealWinByPos(index, true)
				if self.Dealer.Total+self.Total+self.Robot.RobotTotal-dealwin < self.Dealer.Total/5 {
					self.Robot.GameBackBets(self.Robot.Robots[i], index, gold)
					continue
				}
			} else if self.RobotDealer != nil {
				dealwin := self.GetDealWinByPos(index, true)
				if self.RobotDealer.GetMoney()+self.Total+self.Robot.RobotTotal-dealwin < self.RobotDealer.GetMoney()/5 {
					self.Robot.GameBackBets(self.Robot.Robots[i], index, gold)
					continue
				}
			} else if lib.GetManyMgr().GetProperty(self.room.Type).DealChange == 1 {
				dealwin := self.GetDealWinByPos(index, true)
				if self.Money+self.Total+self.Robot.RobotTotal-dealwin < self.Money/5 {
					self.Robot.GameBackBets(self.Robot.Robots[i], index, gold)
					continue
				}
			}
			var msg Msg_GameGoldBZW_Bets
			msg.Uid = self.Robot.Robots[i].Id
			msg.Index = index
			msg.Gold = gold
			msg.Total = self.Robot.Robots[i].GetMoney()
			self.room.broadCastMsg("gamelhdbets", &msg)
		}
		return
	}

	if !self.room.Begin {
		self.OnBegin()
	}
}
