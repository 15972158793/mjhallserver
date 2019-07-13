package gameserver

import (
	//"fmt"
	"lib"
	"math"
	"staticfunc"
	"time"
)

//! 金币场记录
type Rec_TXLHD_Info struct {
	GameType int                    `json:"gametype"`
	Time     int64                  `json:"time"` //! 记录时间
	Deal     int64                  `json:"deal"`
	Info     []Son_Rec_TXLHD_Person `json:"info"`
}
type Son_Rec_TXLHD_Person struct {
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`
	Head   string `json:"head"`
	Score  int    `json:"score"`
	Result [2]int `json:"result"`
	Bets   [3]int `json:"bets"`
}

type Msg_GameTXLHD_Info struct {
	Match []Msg_GameTXLHD_Match `json:"match"`
	Total int                   `json:"total"` //! 当前金币
	Trend []*lib.QQOnline       `json:"trend"` //! 当前版本号
}

type Msg_GameTXLHD_Match struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`    //! 庄家的名字
	Bets    [3]int `json:"bets"`    //! 3个位置的下注
	AllBets [3]int `json:"allbets"` //! 3个位置的总下注
	Time    int64  `json:"time"`    //! 下注结束时间
	Create  int64  `json:"create"`
	Result  []int  `json:"result"` //! 结果
	Win     int    `json:"win"`
}

type Msg_GameTXLHD_Bets struct {
	Uid     int64  `json:"uid"`
	Total   int    `json:"total"`
	Id      int    `json:"id"`
	Bets    [3]int `json:"bets"`    //! 3个位置的下注
	AllBets [3]int `json:"allbets"` //! 3个位置的总下注
}

type Msg_GameTXLHD_AddMatch struct {
	Uid     int64  `json:"uid"`
	Total   int    `json:"total"`
	Id      int    `json:"id"`
	Name    string `json:"name"`    //! 庄家的名字
	Bets    [3]int `json:"bets"`    //! 3个位置的下注
	AllBets [3]int `json:"allbets"` //! 3个位置的总下注
	Time    int64  `json:"time"`    //! 下注结束时间
	Create  int64  `json:"create"`
	Result  []int  `json:"result"` //! 结果
}

type Msg_GameTXLHD_DelMatch struct {
	Id []int `json:"id"`
}

type Msg_GameTXLHD_ResultMatch struct {
	Id     int   `json:"id"`
	Result []int `json:"result"`
	Win    int   `json:"win"`
}

type Msg_GameTXLHD_Total struct {
	Total int `json:"total"`
}

type Msg_GameTXLHD_Trend struct {
	Trend *lib.QQOnline `json:"trend"`
}

type Clint_GameTXLHD_Bets struct {
	Id    int `json:"id"`
	Index int `json:"index"`
	Gold  int `json:"gold"`
}

///////////////////////////////////////////////////////
type Game_TXLHD_Person struct {
	Uid   int64  `json:"uid"`
	Gold  int    `json:"gold"`  //! 进来时候的钱
	Total int    `json:"total"` //! 当前的钱
	Name  string `json:"name"`  //! 名字
	Head  string `json:"head"`  //! 头像
}

//! 同步金币
func (self *Game_TXLHD_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

type Game_TXLHD_Match struct {
	Id      int              //! 比赛id
	Uid     int64            //! 庄家uid
	Name    string           //! 庄家名字
	Bets    [3]map[int64]int //! 下注
	AllBets [3]int           //! 总下注
	Create  int64            //! 创建时间
	Time    int64            //! 开奖的时间
	Result  []int            //! 结果
	Win     map[int64]int    //! 每个uid的赢钱
	Robot   bool             //! 庄家是否为机器人
}

func (self *Game_TXLHD_Match) Init() {
	self.Bets[0] = make(map[int64]int)
	self.Bets[1] = make(map[int64]int)
	self.Bets[2] = make(map[int64]int)
	self.Result = make([]int, 0)
	self.Win = make(map[int64]int)
}

func (self *Game_TXLHD_Match) GetBets(uid int64) [3]int {
	var bets [3]int
	for i := 0; i < 3; i++ {
		bets[i] = self.Bets[i][uid]
	}
	return bets
}

func (self *Game_TXLHD_Match) GetAllBets(uid int64) int {
	var bets int
	for i := 0; i < 3; i++ {
		bets += self.Bets[i][uid]
	}
	return bets
}

//! 得到状态 0下注中  1封注中   2结算中
func (self *Game_TXLHD_Match) GetState() int {
	if len(self.Result) > 0 {
		return 2
	}

	if time.Now().Unix() >= self.Time-5 {
		return 1
	}

	return 0
}

type Game_TXLHD struct {
	PersonMgr map[int64]*Game_TXLHD_Person
	Match     []*Game_TXLHD_Match
	MatchId   int
	DealMoney int
	Trend     []*lib.QQOnline
	MinDeal   int //! 是否有机器人庄
	Ver       int //! 当前版本号

	room *Room
}

func NewGame_TXLHD() *Game_TXLHD {
	game := new(Game_TXLHD)
	game.PersonMgr = make(map[int64]*Game_TXLHD_Person)
	game.Match = make([]*Game_TXLHD_Match, 0)
	game.MatchId = 1
	game.Trend = make([]*lib.QQOnline, 0)
	game.RefreshTrend()
	game.Ver = 1231
	game.MinDeal = -1

	return game
}

func (self *Game_TXLHD) OnInit(room *Room) {
	self.room = room

	if self.room.Type%10 == 0 {
		self.DealMoney = 100000
	} else if self.room.Type%10 == 1 {
		self.DealMoney = 500000
	} else if self.room.Type%10 == 2 {
		self.DealMoney = 1000000
	}

	//! 机器人当庄
	robot := lib.GetRobotMgr().GetRobot(220000, 1000000, 1000000, 1000000)
	if robot != nil {
		self.RobotDeal(robot)
	}
}

func (self *Game_TXLHD) OnRobot(robot *lib.Robot) {

}

func (self *Game_TXLHD) OnSendInfo(person *Person) {
	//! 观众模式游戏,观众进来只发送游戏信息
	value, ok := self.PersonMgr[person.Uid]
	if ok {
		value.SynchroGold(person.Gold)
		person.SendMsg("gametxlhdinfo", self.getInfo(person.Uid, value.Total))
		return
	}

	_person := new(Game_TXLHD_Person)
	_person.Uid = person.Uid
	_person.Gold = person.Gold
	_person.Total = person.Gold
	_person.Name = person.Name
	_person.Head = person.Imgurl
	self.PersonMgr[person.Uid] = _person
	person.SendMsg("gametxlhdinfo", self.getInfo(person.Uid, person.Gold))
}

func (self *Game_TXLHD) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold": //! 同步金币
		person := self.GetPerson(msg.V.(*staticfunc.Msg_SynchroGold).Uid)
		if person != nil {
			person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal(person.Uid, person.Total)
		}
	case "gametxlhdbets": //! 下注
		self.GameBets(msg.Uid, msg.V.(*Clint_GameTXLHD_Bets).Id, msg.V.(*Clint_GameTXLHD_Bets).Index, msg.V.(*Clint_GameTXLHD_Bets).Gold)
	case "gamerob": //! 上庄
		self.GameDeal(msg.Uid)
	}
}

func (self *Game_TXLHD) GetMatchId() int {
	self.MatchId++
	return self.MatchId
}

func (self *Game_TXLHD) GetPerson(uid int64) *Game_TXLHD_Person {
	return self.PersonMgr[uid]
}

//!  刷新牌路
func (self *Game_TXLHD) RefreshTrend() {
	_time := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(),
		time.Now().Hour(), time.Now().Minute(), 0, 0, time.Now().Location()).Unix()
	if len(self.Trend) == 0 {
		for i := _time; ; i -= 60 {
			value := lib.GetSSCMgr().GetQQOnline(i)
			if value == nil {
				if i == _time {
					continue
				}
				break
			}
			self.Trend = append(self.Trend, value)
			if len(self.Trend) >= 20 {
				break
			}
		}
	} else if self.Trend[0].Time != _time {
		value := lib.GetSSCMgr().GetQQOnline(_time)
		if value == nil {
			return
		}
		self.Trend = append([]*lib.QQOnline{value}, self.Trend...)
		if len(self.Trend) >= 20 {
			self.Trend = self.Trend[0:20]
		}
		var msg Msg_GameTXLHD_Trend
		msg.Trend = value
		self.room.broadCastMsg("gametxlhdtrend", &msg)
	}
}

//! 当庄
func (self *Game_TXLHD) GameDeal(uid int64) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if person.Total < self.DealMoney {
		self.room.SendErr(uid, "金币不足,无法当庄")
		return
	}

	curtime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(),
		time.Now().Hour(), time.Now().Minute()+1, 0, 0, time.Now().Location()).Unix()
	if time.Now().Second() >= 30 { //! 已经过了30秒
		curtime += 60
	}
	num := 0
	for i := 0; i < len(self.Match); i++ {
		if self.Match[i].Time == curtime {
			num++
		}
	}
	if num >= 3 {
		self.room.SendErr(uid, "本期当庄人数已达上限,请等待下一期。")
		return
	}

	match := new(Game_TXLHD_Match)
	match.Id = self.GetMatchId()
	match.Uid = person.Uid
	match.Name = person.Name
	match.Time = curtime
	match.Create = time.Now().Unix()
	match.Robot = false
	match.Init()
	self.Match = append(self.Match, match)

	person.Total -= self.DealMoney

	var msg Msg_GameTXLHD_AddMatch
	msg.Uid = uid
	msg.Total = person.Total
	msg.Id = match.Id
	msg.Name = match.Name
	msg.Time = match.Time
	msg.Create = match.Create
	msg.Result = match.Result
	self.room.broadCastMsg("gametxlhddeal", &msg)
}

//! 机器人当庄
func (self *Game_TXLHD) RobotDeal(robot *lib.Robot) {
	curtime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(),
		time.Now().Hour(), time.Now().Minute()+1, 0, 0, time.Now().Location()).Unix()
	if time.Now().Second() >= 30 { //! 已经过了30秒
		curtime += 60
	}
	num := 0
	for i := 0; i < len(self.Match); i++ {
		if self.Match[i].Time == curtime {
			num++
		}
	}

	match := new(Game_TXLHD_Match)
	match.Id = self.GetMatchId()
	match.Uid = robot.Id
	match.Name = robot.Name
	match.Time = curtime
	match.Create = time.Now().Unix()
	match.Robot = true
	match.Init()
	self.Match = append(self.Match, match)

	var msg Msg_GameTXLHD_AddMatch
	msg.Uid = robot.Id
	msg.Total = 1000000
	msg.Id = match.Id
	msg.Name = match.Name
	msg.Time = match.Time
	msg.Create = match.Create
	msg.Result = match.Result
	self.room.broadCastMsg("gametxlhddeal", &msg)
}

//! 下注
func (self *Game_TXLHD) GameBets(uid int64, id int, index int, gold int) {
	if index < 0 || index >= 3 {
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if person.Total < gold {
		self.room.SendErr(uid, "金币不足,无法下注")
		return
	}

	var match *Game_TXLHD_Match = nil
	for i := 0; i < len(self.Match); i++ {
		if self.Match[i].Id == id {
			match = self.Match[i]
			break
		}
	}

	if match == nil {
		return
	}

	state := match.GetState()
	if state == 2 {
		self.room.SendErr(uid, "结算中,无法下注")
		return
	} else if state == 1 {
		self.room.SendErr(uid, "已封注,无法下注")
		return
	}

	if index == 0 {
		if self.DealMoney+match.AllBets[1]+match.AllBets[2]-(match.AllBets[0]+gold) < 0 {
			self.room.SendErr(uid, "该位置暂达上限,请下其他位置")
			return
		}
	} else if index == 1 {
		if self.DealMoney+match.AllBets[0]+match.AllBets[2]-(match.AllBets[1]+gold) < 0 {
			self.room.SendErr(uid, "该位置暂达上限,请下其他位置")
			return
		}
	} else {
		if self.DealMoney+match.AllBets[0]+match.AllBets[1]-4*(match.AllBets[2]+gold) < 0 {
			self.room.SendErr(uid, "该位置暂达上限,请下其他位置")
			return
		}
	}

	match.Bets[index][uid] += gold
	match.AllBets[index] += gold
	person.Total -= gold

	var msg Msg_GameTXLHD_Bets
	msg.Uid = uid
	msg.Total = person.Total
	msg.Bets = match.GetBets(uid)
	msg.Id = match.Id
	msg.AllBets = match.AllBets
	self.room.broadCastMsg("gametxlhdbets", &msg)
}

func (self *Game_TXLHD) AddGold(uid int64, gold int) {
	person := self.GetPerson(uid)
	if person != nil {
		person.Total += gold
		self.SendTotal(uid, person.Total)
		return
	}

	GetRoomMgr().AddCard2(uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
}

func (self *Game_TXLHD) SendTotal(uid int64, total int) {
	var msg Msg_GameTXLHD_Total
	msg.Total = total
	self.room.SendMsg(uid, "gametxlhdtotal", &msg)
}

func (self *Game_TXLHD) OnBegin() {
	if self.room.IsBye() {
		return
	}
}

//! 结算
func (self *Game_TXLHD) OnEnd() {
	self.room.Begin = false
}

func (self *Game_TXLHD) OnBye() {
}

func (self *Game_TXLHD) OnExit(uid int64) {
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

func (self *Game_TXLHD) getInfo(uid int64, total int) *Msg_GameTXLHD_Info {
	var msg Msg_GameTXLHD_Info
	msg.Total = total
	msg.Trend = self.Trend
	msg.Match = make([]Msg_GameTXLHD_Match, 0)
	for i := 0; i < len(self.Match); i++ {
		var info Msg_GameTXLHD_Match
		info.Id = self.Match[i].Id
		info.Name = self.Match[i].Name
		info.Bets = self.Match[i].GetBets(uid)
		info.AllBets = self.Match[i].AllBets
		info.Time = self.Match[i].Time
		info.Create = self.Match[i].Create
		info.Result = self.Match[i].Result
		info.Win = self.Match[i].Win[uid]
		msg.Match = append(msg.Match, info)
	}
	return &msg
}

func (self *Game_TXLHD) OnTime() {
	if time.Now().Second() == 30 && time.Now().Minute() != self.MinDeal {
		robot := lib.GetRobotMgr().GetRobot(220000, 1000000, 1000000, 1000000)
		if robot != nil {
			self.RobotDeal(robot)
		}
		self.MinDeal = time.Now().Minute()
	}

	//! 刷新战绩
	self.RefreshTrend()

	lst := make([]int, 0)
	//! 结算所有比赛
	for i := 0; i < len(self.Match); i++ {
		value := self.Match[i]
		if len(value.Result) > 0 {
			if time.Now().Unix()-value.Time >= 300 {
				lst = append(lst, value.Id)
				copy(self.Match[i:], self.Match[i+1:])
				self.Match = self.Match[:len(self.Match)-1]
				i--
			}
			continue
		}
		result := lib.GetSSCMgr().GetQQOnline(value.Time)
		if result == nil {
			continue
		}
		value.Result = append(value.Result, result.Onlinenumber/10%10)
		value.Result = append(value.Result, result.Onlinenumber%10)

		dealmoney := self.DealMoney
		dealmoney += value.AllBets[0]
		dealmoney += value.AllBets[1]
		dealmoney += value.AllBets[2]
		for j := 0; j < 3; j++ {
			for uid, gold := range value.Bets[j] {
				if j == 0 {
					if value.Result[0] > value.Result[1] {
						dealmoney -= 2 * gold
						value.Win[uid] += 2 * gold
					} else {
						value.Win[uid] += 0
					}
				} else if j == 1 {
					if value.Result[0] < value.Result[1] {
						dealmoney -= 2 * gold
						value.Win[uid] += 2 * gold
					} else {
						value.Win[uid] += 0
					}
				} else {
					if value.Result[0] == value.Result[1] {
						dealmoney -= 5 * gold
						value.Win[uid] += 5 * gold
					} else {
						value.Win[uid] += 0
					}
				}
			}
		}

		if !value.Robot {
			value.Win[value.Uid] += dealmoney
		} else {
			lib.GetRobotMgr().BackRobot(value.Uid, false)
			if dealmoney != 0 {
				GetServer().SqlBZWLog(&SQL_BZWLog{1, dealmoney - self.DealMoney, time.Now().Unix(), 220000 + 10000000})
			}
		}

		for uid, gold := range value.Win {
			bets := value.GetAllBets(uid)
			if uid == value.Uid {
				bets += self.DealMoney
			}
			addgold := gold - bets
			cost := 0
			if addgold > 0 {
				cost = int(math.Ceil(float64(addgold) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
				GetServer().SqlAgentGoldLog(uid, cost, self.room.Type)
				GetServer().SqlAgentBillsLog(uid, cost/2, self.room.Type)
				addgold -= cost
			} else if addgold < 0 {
				cost = int(math.Ceil(float64(addgold) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 200.0))
				GetServer().SqlAgentBillsLog(uid, cost, self.room.Type)
			}
			self.AddGold(uid, gold-cost)

			var record Rec_TXLHD_Info
			record.Time = value.Time
			record.GameType = self.room.Type
			record.Deal = value.Uid
			var rec Son_Rec_TXLHD_Person
			rec.Uid = uid
			p := self.GetPerson(uid)
			if p != nil {
				rec.Name = p.Name
				rec.Head = p.Head
			} else {
				p := GetPersonMgr().ForcePerson(uid)
				rec.Name = p.Name
				rec.Head = p.Imgurl
			}
			rec.Score = addgold
			rec.Result = [2]int{value.Result[0], value.Result[1]}
			rec.Bets = value.GetBets(uid)
			record.Info = append(record.Info, rec)
			GetServer().InsertRecord(self.room.Type, uid, lib.HF_JtoA(&record), rec.Score)
		}

		for _, p := range self.PersonMgr {
			var msg Msg_GameTXLHD_ResultMatch
			msg.Id = value.Id
			msg.Result = value.Result
			msg.Win = value.Win[p.Uid]
			self.room.SendMsg(p.Uid, "gametxlhdresult", &msg)
		}
	}

	if len(lst) > 0 {
		var msg Msg_GameTXLHD_DelMatch
		msg.Id = lst
		self.room.broadCastMsg("gametxlhddel", &msg)
	}
}

func (self *Game_TXLHD) OnIsDealer(uid int64) bool {
	return false
}

//! 是否下注了
func (self *Game_TXLHD) OnIsBets(uid int64) bool {
	return false
}

//! 结算所有人
func (self *Game_TXLHD) OnBalance() {
	//! 被clear时先返还本场下注
	for i := 0; i < len(self.Match); i++ {
		if len(self.Match[i].Result) > 0 {
			continue
		}
		for j := 0; j < 3; j++ {
			for key, value := range self.Match[i].Bets[j] {
				person := self.GetPerson(key)
				if person != nil {
					person.Total += value
				} else {
					GetRoomMgr().AddCard2(key, staticfunc.TYPE_GOLD, value, self.room.Type)
				}
			}
		}

		if !self.Match[i].Robot {
			person := self.GetPerson(self.Match[i].Uid)
			if person != nil {
				person.Total += self.DealMoney
			} else {
				GetRoomMgr().AddCard2(self.Match[i].Uid, staticfunc.TYPE_GOLD, self.DealMoney, self.room.Type)
			}
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
