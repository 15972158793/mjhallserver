package gameserver

import (
	//"fmt"
	"lib"
	"math"
	"staticfunc"
	"time"
)

//! 红包扫雷时间
var HBSL_TIME int = 10

//! 金币场记录
type Rec_HBSL_Info struct {
	GameType int                   `json:"gametype"`
	Time     int64                 `json:"time"` //! 记录时间
	Info     []Son_Rec_HBSL_Person `json:"info"`
	Land     int                   `json:"land"` //! 雷针
}
type Son_Rec_HBSL_Person struct {
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`
	Head   string `json:"head"`
	Score  int    `json:"score"`
	Result int    `json:"result"`
}

type Msg_GameGoldHBSL_Info struct {
	Begin  bool                 `json:"begin"`  //! 是否开始
	Time   int64                `json:"time"`   //! 倒计时
	Dealer Game_GoldHBSL_Dealer `json:"dealer"` //! 庄家
	Total  int                  `json:"total"`  //! 自己的钱
	Play   []Game_GoldHBSL_Play `json:"play"`
}

//! 申请埋雷
type Msg_GameLand struct {
	Money int `json:"money"` //! 埋雷金额
	Land  int `json:"land"`  //! 雷针
}

//! 上庄列表
type Msg_DealList struct {
	Info []*Game_GoldHBSL_Dealer `json:"info"`
}

//! 上庄
type Msg_GameHBSLDeal struct {
	Info *Game_GoldHBSL_Dealer `json:"info"`
}

//! 抢红包
type Msg_GameRob struct {
	Uid   int64                `json:"uid"`
	Info  []Game_GoldHBSL_Play `json:"info"`
	Total int                  `json:"total"`
	Num   int                  `json:"num"`
}

///////////////////////////////////////////////////////
type Game_GoldHBSL_Person struct {
	Uid     int64  `json:"uid"`
	Gold    int    `json:"gold"`  //! 进来时候的钱
	Total   int    `json:"total"` //! 当前的钱
	Name    string `json:"name"`  //! 名字
	Head    string `json:"head"`  //! 头像
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
	Round   int    `json:"round"` //! 不下注轮数
}

//! 庄家
type Game_GoldHBSL_Dealer struct {
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Money   int    `json:"money"` //! 红包金额
	Num     int    `json:"num"`   //! 红包数量
	Land    int    `json:"land"`  //! 雷针
	Win     int    `json:"win"`
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
	isrobot bool
}

//! 结算
type Msg_GameHBSLBalance struct {
	Win int `json:"win"` //!　赢了多少
}

//! 扫雷记录
type Game_GoldHBSL_Play struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Money int    `json:"money"` //! 扫到的钱
	Win   int    `json:"win"`
}

//! 同步金币
func (self *Game_GoldHBSL_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

type Game_GoldHBSL struct {
	PersonMgr map[int64]*Game_GoldHBSL_Person `json:"personmgr"`
	Dealer    *Game_GoldHBSL_Dealer           `json:"dealer"` //! 庄家
	Time      int64                           `json:"time"`
	State     int                             `json:"state"`
	LstDeal   []*Game_GoldHBSL_Dealer         `json:"lstdeal"` //! 上庄列表
	Play      []Game_GoldHBSL_Play            `json:"play"`
	Robot     lib.ManyGameRobot               //! 机器人结构
	NoDeal    int64                           `json:"nodeal"` //! 没有庄家的时间

	RedNum int     `json:"rednum"` //! 红包个数
	RedBS  float32 `json:"redbs"`  //! 红包倍数

	room *Room
}

func NewGame_GoldHBSL() *Game_GoldHBSL {
	game := new(Game_GoldHBSL)
	game.PersonMgr = make(map[int64]*Game_GoldHBSL_Person)

	return game
}

func (self *Game_GoldHBSL) OnInit(room *Room) {
	self.room = room
	//! 载入机器人
	self.Robot.Init(17, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)
	self.NoDeal = time.Now().Unix()

	if self.room.Type%10 == 0 { //! 10包1倍
		self.RedNum = 10
		self.RedBS = 1.0
	} else if self.room.Type%10 == 1 { //! 7包1.5倍
		self.RedNum = 7
		self.RedBS = 1.5
	}
}

func (self *Game_GoldHBSL) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldHBSL) OnSendInfo(person *Person) {
	//! 观众模式游戏,观众进来只发送游戏信息
	value, ok := self.PersonMgr[person.Uid]
	if ok {
		value.IP = person.ip
		value.Address = person.minfo.Address
		value.Sex = person.Sex
		value.SynchroGold(person.Gold)
		person.SendMsg("gamegoldhbslinfo", self.getInfo(person.Uid, value.Total))
		return
	}

	_person := new(Game_GoldHBSL_Person)
	_person.Uid = person.Uid
	_person.Gold = person.Gold
	_person.Total = person.Gold
	_person.Name = person.Name
	_person.Head = person.Imgurl
	_person.IP = person.ip
	_person.Address = person.minfo.Address
	_person.Sex = person.Sex
	self.PersonMgr[person.Uid] = _person
	person.SendMsg("gamegoldhbslinfo", self.getInfo(person.Uid, person.Gold))
}

func (self *Game_GoldHBSL) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold": //! 同步金币
		person := self.GetPerson(msg.V.(*staticfunc.Msg_SynchroGold).Uid)
		if person != nil {
			person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal(person.Uid, person.Total)
		}
	case "gamedeallist": //! 得到上庄列表
		self.GetDealList(msg.Uid)
	case "gameland": //! 埋雷
		self.GameLand(msg.Uid, msg.V.(*Msg_GameLand).Money, msg.V.(*Msg_GameLand).Land)
	case "gamerob": //! 抢红包
		self.GameRob(msg.Uid)
	case "gameplayerlist":
		self.GamePlayerList(msg.Uid)
	}
}

func (self *Game_GoldHBSL) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.Begin = true

	self.NoDeal = 0
	self.Play = make([]Game_GoldHBSL_Play, 0)
	self.Dealer = self.LstDeal[0]
	self.LstDeal = self.LstDeal[1:]

	var msg Msg_GameHBSLDeal
	msg.Info = self.Dealer
	self.room.broadCastMsg("gamerob", &msg)

	self.Time = time.Now().Unix() + 10
	self.State = 0
}

//! 结算
func (self *Game_GoldHBSL) OnEnd() {
	self.room.Begin = false

	self.Time = time.Now().Unix() + 2
	self.State = 1
	bl := 100
	if GetServer().Con.MoneyMode == 1 {
		bl = 1
	}

	self.Dealer.Win = self.Dealer.Money
	for i := 0; i < len(self.Play); i++ {
		self.Dealer.Win -= self.Play[i].Money
		if self.Play[i].Money/bl%10 == self.Dealer.Land { //! 踩雷了
			self.Dealer.Win += int(float32(self.Dealer.Money) * self.RedBS)
		}
	}
	if self.Dealer.Win > self.Dealer.Money {
		addgold := self.Dealer.Win - self.Dealer.Money
		cost := int(math.Ceil(float64(addgold) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		if !self.Dealer.isrobot {
			GetServer().SqlAgentGoldLog(self.Dealer.Uid, cost, self.room.Type)
			GetServer().SqlAgentBillsLog(self.Dealer.Uid, cost/2, self.room.Type)
		}
		self.Dealer.Win -= cost
	} else if self.Dealer.Win-self.Dealer.Money < 0 && !self.Dealer.isrobot {
		addgold := self.Dealer.Win - self.Dealer.Money
		cost := int(math.Ceil(float64(addgold) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 200.0))
		GetServer().SqlAgentBillsLog(self.Dealer.Uid, cost, self.room.Type)
	}
	if self.Dealer.Win != 0 {
		if !self.Dealer.isrobot { //! 如果不是机器人
			self.AddGold(self.Dealer.Uid, self.Dealer.Win)
		} else {
			robot := lib.GetRobotMgr().GetRobotFromId(self.Dealer.Uid)
			if robot != nil {
				robot.AddMoney(self.Dealer.Win)
			}
		}
	}

	var msg Msg_GameHBSLBalance
	msg.Win = self.Dealer.Win
	self.room.broadCastMsg("gamehbslbalance", &msg)

	if !self.Dealer.isrobot {
		var record Rec_HBSL_Info
		record.Time = time.Now().Unix()
		record.GameType = self.room.Type
		record.Land = self.Dealer.Land
		var rec Son_Rec_HBSL_Person
		rec.Uid = self.Dealer.Uid
		rec.Name = self.Dealer.Name
		rec.Head = self.Dealer.Head
		rec.Score = self.Dealer.Win - self.Dealer.Money
		rec.Result = 0
		record.Info = append(record.Info, rec)
		GetServer().InsertRecord(self.room.Type, self.Dealer.Uid, lib.HF_JtoA(&record), rec.Score)
	}

	//! 返回机器人
	for i := 0; i < len(self.Robot.Robots); i++ {
		if self.Robot.Robots[i].GetMoney() < self.GetMinMoney() {
			self.Robot.Robots[i].Dead()
		}
	}
	self.Robot.Init(17, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)

	for _, value := range self.PersonMgr {
		if self.IsWaitDeal(value.Uid) {
			continue
		}
		value.Round++
		if value.Round >= 5 && GetPersonMgr().GetPerson(value.Uid) == nil {
			self.room.KickViewByUid(value.Uid, 96)
		}
	}
}

func (self *Game_GoldHBSL) OnBye() {
}

func (self *Game_GoldHBSL) OnExit(uid int64) {
	value, ok := self.PersonMgr[uid]
	if ok {
		//! 退出房间同步金币
		gold := value.Total - value.Gold
		if gold > 0 {
			GetRoomMgr().AddCard(value.Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
		} else if gold < 0 {
			GetRoomMgr().CostCard(value.Uid, staticfunc.TYPE_GOLD, -gold, self.room)
		}
		value.Gold = value.Total
		delete(self.PersonMgr, uid)
	}
}

func (self *Game_GoldHBSL) getInfo(uid int64, total int) *Msg_GameGoldHBSL_Info {
	var msg Msg_GameGoldHBSL_Info
	msg.Begin = self.room.Begin
	if self.Time == 0 {
		msg.Time = 0
	} else {
		msg.Time = self.Time - time.Now().Unix()
	}
	msg.Total = total
	if self.Dealer != nil {
		lib.HF_DeepCopy(&msg.Dealer, self.Dealer)
	}
	msg.Play = self.Play
	return &msg
}

func (self *Game_GoldHBSL) GetPerson(uid int64) *Game_GoldHBSL_Person {
	return self.PersonMgr[uid]
}

func (self *Game_GoldHBSL) AddGold(uid int64, gold int) {
	person := self.GetPerson(uid)
	if person != nil {
		person.Total += gold
		self.SendTotal(uid, person.Total)
		return
	}

	GetRoomMgr().AddCard2(uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
}

//! 得到上庄列表
func (self *Game_GoldHBSL) GetDealList(uid int64) {
	var msg Msg_DealList
	msg.Info = self.LstDeal
	self.room.SendMsg(uid, "gamedeallist", &msg)
}

//! 是否在上庄列表里
func (self *Game_GoldHBSL) IsWaitDeal(uid int64) bool {
	for i := 0; i < len(self.LstDeal); i++ {
		if self.LstDeal[i].isrobot {
			continue
		}
		if self.LstDeal[i].Uid == uid {
			return true
		}
	}
	return false
}

//! 申请埋雷
func (self *Game_GoldHBSL) GameLand(uid int64, money int, land int) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if money < self.GetMinMoney() || money%100 != 0 { //! 埋雷金额错误
		self.room.SendErr(uid, "埋雷金额错误")
		return
	}

	if person.Total < money {
		self.room.SendErr(uid, "金币不足,无法埋雷")
		return
	}

	if land < 0 || land > 9 {
		self.room.SendErr(uid, "雷号错误")
		return
	}

	person.Round = 0

	person.Total -= money
	self.LstDeal = append(self.LstDeal, &Game_GoldHBSL_Dealer{person.Uid, person.Name, person.Head, money, self.RedNum, land, 0, person.IP, person.Address, person.Sex, false})

	var msg Msg_GameLand
	msg.Money = person.Total
	msg.Land = land
	self.room.SendMsg(uid, "gameland", &msg)

	if !self.room.Begin && self.State == 0 { //! 没有开始则开始
		self.OnBegin()
	}
}

//! 机器人埋雷
func (self *Game_GoldHBSL) RobotLand(robot *lib.Robot) {
	minmoney := staticfunc.GetCsvMgr().GetDF(self.room.Type)
	maxmoney := staticfunc.GetCsvMgr().GetZR(self.room.Type)
	lst := make([]int, 0)
	for i := minmoney; i <= maxmoney; i += 1000 {
		lst = append(lst, i)
	}
	money := lst[lib.HF_GetRandom(len(lst))]
	if robot.GetMoney() < money { //! 埋雷金额错误
		return
	}

	for i := 0; i < len(self.LstDeal); i++ {
		if !self.LstDeal[i].isrobot {
			continue
		}
		if self.LstDeal[i].Uid == robot.Id { //! 重复的不要
			return
		}
	}

	land := lib.HF_GetRandom(10)
	robot.AddMoney(-money)
	self.LstDeal = append(self.LstDeal, &Game_GoldHBSL_Dealer{robot.Id, robot.Name, robot.Head, money, self.RedNum, land, 0, robot.IP, robot.Address, robot.Sex, true})

	if !self.room.Begin && self.State == 0 { //! 没有开始则开始
		self.OnBegin()
	}
}

func (self *Game_GoldHBSL) GameRob(uid int64) {
	if self.Dealer == nil {
		self.room.SendErr(uid, "请等待玩家发红包")
		return
	}

	if !self.room.Begin {
		self.room.SendErr(uid, "红包已被抢空,请等待下一局")
		return
	}

	if self.Dealer.Uid == uid {
		self.room.SendErr(uid, "不能抢自己发的红包")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if person.Total < int(math.Ceil(float64(float32(self.Dealer.Money)*self.RedBS))) {
		self.room.SendErr(uid, "金币不足,无法抢红包")
		return
	}

	if self.Dealer.Num <= 0 {
		self.room.SendErr(uid, "红包已被抢空,请等待下一局")
		return
	}

	for i := 0; i < len(self.Play); i++ {
		if self.Play[i].Uid == uid {
			self.room.SendErr(uid, "您已抢过该红包了")
			return
		}
	}

	person.Round = 0

	bl := 100
	if GetServer().Con.MoneyMode == 1 {
		bl = 1
	}
	money := 0
	total := 0
	for i := 0; i < len(self.Play); i++ {
		total += self.Play[i].Money / bl
	}
	if len(self.Play) == self.RedNum-1 {
		money = self.Dealer.Money/bl - total
	} else {
		max := lib.HF_MinInt(self.Dealer.Money/bl-total-(self.RedNum-len(self.Play)), self.Dealer.Money/(bl*10)*2)
		money = lib.HF_GetRandom(lib.HF_MaxInt(1, max)-1) + 1
		if self.Dealer.isrobot && money%10 != self.Dealer.Land && lib.GetRobotMgr().GetRobotWin(self.room.Type)-money*bl < 0 { //! 机器人庄,奖池会负数,要踩雷
			money += (self.Dealer.Land - money%10)
			if money <= 0 {
				money = 10 + self.Dealer.Land
			}
		}
	}
	self.Dealer.Num--

	win := money * bl
	if money%10 == self.Dealer.Land {
		win -= int(float32(self.Dealer.Money) * self.RedBS)
	}
	if self.Dealer.isrobot { //! 机器人
		lib.GetRobotMgr().AddRobotWin(self.room.Type, -win)
		GetServer().SqlBZWLog(&SQL_BZWLog{1, -win, time.Now().Unix(), 300000 + 10000000})
	}
	if win > 0 { //! 要抽水
		cost := int(math.Ceil(float64(win) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		GetServer().SqlAgentGoldLog(uid, cost, self.room.Type)
		GetServer().SqlAgentBillsLog(uid, cost/2, self.room.Type)
		win -= cost
	} else if win < 0 {
		cost := int(math.Ceil(float64(win) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 200.0))
		GetServer().SqlAgentBillsLog(uid, cost, self.room.Type)
	}
	person.Total += win
	self.Play = append(self.Play, Game_GoldHBSL_Play{person.Uid, person.Name, person.Head, money * bl, win})

	var msg Msg_GameRob
	msg.Uid = uid
	msg.Info = self.Play
	msg.Total = person.Total
	msg.Num = self.Dealer.Num
	self.room.broadCastMsg("gamehbslrob", &msg)

	var record Rec_HBSL_Info
	record.Time = time.Now().Unix()
	record.GameType = self.room.Type
	record.Land = self.Dealer.Land
	var rec Son_Rec_HBSL_Person
	rec.Uid = uid
	rec.Name = person.Name
	rec.Head = person.Head
	rec.Score = win
	rec.Result = money * bl
	record.Info = append(record.Info, rec)
	GetServer().InsertRecord(self.room.Type, uid, lib.HF_JtoA(&record), rec.Score)

	if self.Dealer.Num <= 0 {
		self.OnEnd()
	}
}

//! 机器人抢红包
func (self *Game_GoldHBSL) RobotRob(robot *lib.Robot) {
	if self.Dealer == nil {
		return
	}

	if !self.room.Begin {
		return
	}

	if self.Dealer.Uid == robot.Id {
		return
	}

	if robot.GetMoney() < int(math.Ceil(float64(float32(self.Dealer.Money)*self.RedBS))) {
		return
	}

	if self.Dealer.Num <= 0 {
		return
	}

	for i := 0; i < len(self.Play); i++ {
		if self.Play[i].Uid == robot.Id {
			return
		}
	}

	bl := 100
	if GetServer().Con.MoneyMode == 1 {
		bl = 1
	}
	money := 0
	total := 0
	for i := 0; i < len(self.Play); i++ {
		total += self.Play[i].Money / bl
	}
	if len(self.Play) == self.RedNum-1 {
		money = self.Dealer.Money/bl - total
	} else {
		max := lib.HF_MinInt(self.Dealer.Money/bl-total-(self.RedNum-len(self.Play)), self.Dealer.Money/(bl*10)*2)
		money = lib.HF_GetRandom(lib.HF_MaxInt(max, 1)-1) + 1
		if !self.Dealer.isrobot && money%10 == self.Dealer.Land && lib.GetRobotMgr().GetRobotWin(self.room.Type)+(money*bl-self.Dealer.Money) <= 0 { //! 会踩雷
			money += 1
		}
	}
	self.Dealer.Num--

	win := money * bl
	if money%10 == self.Dealer.Land {
		win -= int(float32(self.Dealer.Money) * self.RedBS)
	}
	if !self.Dealer.isrobot { //! 玩家庄
		lib.GetRobotMgr().AddRobotWin(self.room.Type, win)
		GetServer().SqlBZWLog(&SQL_BZWLog{1, win, time.Now().Unix(), 300000 + 10000000})
	}
	if win > 0 { //! 要抽水
		cost := int(math.Ceil(float64(win) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		win -= cost
	}
	robot.AddMoney(win)
	self.Play = append(self.Play, Game_GoldHBSL_Play{robot.Id, robot.Name, robot.Head, money * bl, win})

	var msg Msg_GameRob
	msg.Uid = robot.Id
	msg.Info = self.Play
	msg.Total = robot.GetMoney()
	msg.Num = self.Dealer.Num
	self.room.broadCastMsg("gamehbslrob", &msg)

	if self.Dealer.Num <= 0 {
		self.OnEnd()
	}
}

//! 申请无座玩家
func (self *Game_GoldHBSL) GamePlayerList(uid int64) {
	var msg Msg_GameGoldBZW_List
	tmp := make(map[int64]Son_GameGoldBZW_Info)
	for _, value := range self.PersonMgr {
		var node Son_GameGoldBZW_Info
		node.Uid = value.Uid
		node.Name = value.Name
		node.Total = value.Total
		node.Head = value.Head
		tmp[node.Uid] = node
	}
	for i := 0; i < len(self.Robot.Robots); i++ {
		var node Son_GameGoldBZW_Info
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

//! 得到基础埋雷金额
func (self *Game_GoldHBSL) GetMinMoney() int {
	return staticfunc.GetCsvMgr().GetDF(self.room.Type)
}

func (self *Game_GoldHBSL) OnTime() {
	if lib.GetRobotMgr().GetRobotSet(self.room.Type).Dealer && len(self.Robot.Robots) > 0 && len(self.LstDeal) < 5 { //! 需要机器人上庄
		self.RobotLand(self.Robot.Robots[lib.HF_GetRandom(len(self.Robot.Robots))])
	}

	if self.Time == 0 {
		return
	}

	if time.Now().Unix() >= self.Time {
		self.Time = 0
		if self.State == 0 { //! 结算
			self.OnEnd()
		} else {
			self.State = 0
			if len(self.LstDeal) > 0 {
				self.OnBegin()
			} else {
				self.NoDeal = time.Now().Unix()
				self.Dealer = nil
				self.Play = make([]Game_GoldHBSL_Play, 0)
				var msg Msg_GameHBSLDeal
				msg.Info = new(Game_GoldHBSL_Dealer)
				self.room.broadCastMsg("gamerob", &msg)
			}
		}
	} else {
		if self.State == 0 { //! 抢红阶段
			for i := 0; i < len(self.Robot.Robots); i++ {
				if lib.HF_GetRandom(100) >= 20 {
					continue
				}
				if lib.HF_GetRandom(100) >= 100-lib.GetRobotMgr().GetRobotSet(self.room.Type).BetRate {
					continue
				}
				self.RobotRob(self.Robot.Robots[i])
			}
		}
	}
}

func (self *Game_GoldHBSL) OnIsDealer(uid int64) bool {
	return false
}

//! 同步总分
func (self *Game_GoldHBSL) SendTotal(uid int64, total int) {
	var msg Msg_GameGoldBZW_Total
	msg.Uid = uid
	msg.Total = total

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	self.room.SendMsg(uid, "gamegoldtotal", &msg)
}

//! 是否下注了
func (self *Game_GoldHBSL) OnIsBets(uid int64) bool {
	return false
}

//! 结算所有人
func (self *Game_GoldHBSL) OnBalance() {
	if self.Dealer != nil { //! 返回庄家的金币
		person := self.GetPerson(self.Dealer.Uid)
		if person != nil {
			person.Total += self.Dealer.Money
		} else {
			GetRoomMgr().AddCard2(self.Dealer.Uid, staticfunc.TYPE_GOLD, self.Dealer.Money, self.room.Type)
		}
	}

	for i := 0; i < len(self.LstDeal); i++ {
		person := self.GetPerson(self.LstDeal[i].Uid)
		if person != nil {
			person.Total += self.LstDeal[i].Money
		} else {
			GetRoomMgr().AddCard2(self.LstDeal[i].Uid, staticfunc.TYPE_GOLD, self.LstDeal[i].Money, self.room.Type)
		}
	}

	for _, value := range self.PersonMgr {
		gold := value.Total - value.Gold
		if gold > 0 {
			GetRoomMgr().AddCard(value.Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
		} else if gold < 0 {
			GetRoomMgr().CostCard(value.Uid, staticfunc.TYPE_GOLD, -gold, self.room)
		}
		value.Gold = value.Total
	}
}
