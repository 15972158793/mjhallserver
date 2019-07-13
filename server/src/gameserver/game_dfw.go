package gameserver

import (
	"lib"
	"math"
	"staticfunc"
	"time"
)

var DFW_MAP []int = []int{0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 2, 1, 0, 0, 0, 1, 0, 1, 0, 3, 1, 0, 0, 1, 5, 0, 0, 1, 0, 7, 0, 0, 1, 0, 15, 0, 0, 1, 1, 40, 0, 0, 0, 1000, 0, 10000, 25000}

type Rec_DFW_Info struct {
	GameType int                  `json:"gametype"`
	Time     int64                `json:"time"` //! 记录时间
	Info     []Son_Rec_DFW_Person `json:"info"`
}
type Son_Rec_DFW_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Score int    `json:"score"`
	Bets  int    `json:"bets"`
}

type Game_DFW struct {
	Person    *Game_DFW_Person `json:"person"`
	Bets      int              `json:"bets"`      //! 底分
	CurBet    int              `json:"curbet"`    //! 当前已经得了多少分
	Index     int              `json:"index"`     //! 当前位置
	Remaining int              `json:"remaining"` //! 还剩多少次机会
	Result    []int            `json:"result"`    //! 模拟结果
	Time      int64            `json:"time"`

	room *Room
}

func NewGame_DFW() *Game_DFW {
	game := new(Game_DFW)
	game.Bets = 100
	game.CurBet = 0
	game.Index = -1
	game.Remaining = 8
	return game
}

type Game_DFW_Person struct {
	Uid     int64  `json:"uid"`
	Gold    int    `json:"gold"`
	Total   int    `json:"total"`
	Win     int    `json:"win"`  //! 赢钱
	Cost    int    `json:"cost"` //!  抽水
	Address string `json:"address"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	IP      string `json:"ip"`
	Sex     int    `json:"sex"`
}

type Msg_GameDFW_Info struct {
	Begin     bool             `json:"begin"`
	Bets      int              `json:"bets"`      //! 底分
	CurBet    int              `json:"curbet"`    //! 当前已经得了多少分
	Index     int              `json:"index"`     //! 当前位置
	Remaining int              `json:"remaining"` //! 还剩多少次机会
	Person    Son_GameDFW_Info `json:"person"`
}

type Son_GameDFW_Info struct {
	Uid     int64  `json:"uid"`
	Win     int    `json:"win"`
	Total   int    `json:"total"`
	Address string `json:"address"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	IP      string `json:"ip"`
	Sex     int    `json:"sex"`
}

type Msg_GameDFW_Total struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"`
}

type Msg_GameDFW_DS struct {
	Ds        int `json:"ds"`
	Index     int `json:"index"`
	Remaining int `json:"remaining"` //! 还剩多少次机会
}

type Msg_GameDFW_End struct {
	Index     int              `json:"index"`
	Remaining int              `json:"remaining"` //! 还剩多少次机会
	EndType   int              `json:"endtype"`   //! 1-踩地雷 2-勾走 3-大赢家
	Person    Son_GameDFW_Info `json:"person"`
}

//! 同步金币
func (self *Game_DFW_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

//!　同步总分
func (self *Game_DFW) SendTotal(uid int64, total int) {
	var msg Msg_GameDFW_Total
	msg.Uid = uid
	msg.Total = total
	self.room.SendMsg(uid, "gametotal", &msg)
}

func (self *Game_DFW) OnIsBets(uid int64) bool {
	if self.room.Begin {
		return true
	} else {
		return false
	}
}

func (self *Game_DFW) getinfo(uid int64) *Msg_GameDFW_Info {
	var msg Msg_GameDFW_Info
	msg.Begin = self.room.Begin
	msg.Bets = self.Bets
	msg.CurBet = self.CurBet
	msg.Index = self.Index
	msg.Remaining = self.Remaining
	if self.Person != nil && uid == self.Person.Uid {
		msg.Person.Uid = self.Person.Uid
		msg.Person.Total = self.Person.Total
		msg.Person.Sex = self.Person.Sex
		msg.Person.Name = self.Person.Name
		msg.Person.IP = self.Person.IP
		msg.Person.Head = self.Person.Head
		msg.Person.Address = self.Person.Address
		msg.Person.Win = self.Person.Win
	}
	return &msg
}

func (self *Game_DFW) GameBets(uid int64, bet int) {

	if self.Person == nil || self.Person.Uid != uid {
		return
	}

	if bet < 100 || bet > 1000 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "bet<100 || bet>1000  bet : ", bet)
		return
	}

	if self.Person.Total < bet {
		self.room.SendErr(uid, "您的金币不足，请前往充值")
		return
	}

	self.Bets = bet
	self.Person.Total -= bet
	var msg Msg_GameDFW_Total
	msg.Uid = uid
	msg.Total = self.Person.Total
	self.room.SendMsg(uid, "gamedfwbet", &msg)

	self.OnBegin()
}

//! 掷骰子
func (self *Game_DFW) GameZhi(uid int64) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	ds := self.Result[8-self.Remaining]
	/*
		for i := 1; i <= ds; i++ {
			if DFW_MAP[self.Index+i] != 0 && DFW_MAP[self.Index+i] != 1 {
				self.CurBet = self.Bets * DFW_MAP[self.Index+i]
			}
		}
	*/
	self.Index += ds
	self.Remaining--
	var msg Msg_GameDFW_DS
	msg.Ds = ds
	msg.Index = self.Index
	msg.Remaining = self.Remaining
	self.room.SendMsg(uid, "gamedfwsz", &msg)
	if self.Remaining == 0 || DFW_MAP[self.Index] == 1 || self.CurBet == self.Bets*25000 {
		self.OnEnd()
	}

	self.Time = time.Now().Unix() + 3600
}

func (self *Game_DFW) OnBegin() {
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "开始游戏")

	if self.Person == nil {
		return
	}

	self.room.Begin = true

	//!　初始化
	self.CurBet = 0
	self.Index = -1
	self.Remaining = 8
	self.Result = make([]int, 0)
	self.Person.Win = 0
	self.Person.Cost = 0

	lib.GetLogMgr().Output(lib.LOG_ERROR, "sysmoney :", GetServer().DfwSysMoney[self.room.Type%170000], " jackpotmin : ", lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin, " max ", lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax)

	//!　模拟游戏结果
	{
		lst := make([][]int, 0)
		winlst := make([][]int, 0)
		lostlst := make([][]int, 0)
		for i := 0; i < 10; i++ {
			index := -1              //! 开始位置
			result := make([]int, 0) //! 每一轮掷的点数
			win := 0                 //! 赢了多少
			for j := 0; j < 8; j++ {
				sz := lib.HF_GetRandom(6) + 1
				for k := 1; k <= sz; k++ {
					if DFW_MAP[index+k] != 0 && DFW_MAP[index+k] != 1 {
						win = self.Bets * DFW_MAP[index+k]
					}
				}
				index += sz
				result = append(result, sz)
				if DFW_MAP[index] == 1 {
					break
				}
			}
			win -= self.Bets
			lib.GetLogMgr().Output(lib.LOG_DEBUG, " win : ", win)
			if GetServer().DfwSysMoney[self.room.Type%170000]-int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && GetServer().DfwSysMoney[self.room.Type%170000]-int64(win) <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, " 添加lst")
				lst = append(lst, result)
			}
			if win > 0 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, " 添加win")
				winlst = append(winlst, result)
			} else {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, " 添加lost")
				lostlst = append(lostlst, result)
			}
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "  ")
		}

		lib.GetLogMgr().Output(lib.LOG_DEBUG, "len(list) : ", len(lst), " len(winlst) : ", len(winlst), " len(lost) : ", len(lostlst))

		if len(lst) == 0 {
			if GetServer().DfwSysMoney[self.room.Type%170000] >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax && len(winlst) > 0 { //! 一定输
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "---- 庄家一定输")
				self.Result = winlst[lib.HF_GetRandom(len(winlst))]
			} else if GetServer().DfwSysMoney[self.room.Type%170000] <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && len(lostlst) > 0 { //! 一定赢
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "---- 庄家一定赢")
				self.Result = lostlst[lib.HF_GetRandom(len(lostlst))]
			} else { //! 纯随机
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "---- 纯随机")
				index := -1 //! 开始位置
				for i := 0; i < 8; i++ {
					sz := lib.HF_GetRandom(6) + 1
					index += sz
					self.Result = append(self.Result, sz)
					if DFW_MAP[index] == 1 {
						break
					}
				}
			}
		} else {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "---- 随机lst")
			self.Result = lst[lib.HF_GetRandom(len(lst))]
		}
	}

	//self.Result = []int{6, 6, 6, 6, 6, 6, 6, 6}
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "----- self.result : ", self.Result)

	self.room.broadCastMsg("gamedfwbegin", self.getinfo(self.Person.Uid))

	index := -1
	for i := 0; i < len(self.Result); i++ {
		for k := 1; k <= self.Result[i]; k++ {
			if DFW_MAP[index+k] != 0 && DFW_MAP[index+k] != 1 {
				self.CurBet = self.Bets * DFW_MAP[index+k]
			}
		}
		index += self.Result[i]
		if DFW_MAP[index] == 1 {
			break
		}
	}

	if self.CurBet > 0 {
		self.Person.Win = self.CurBet
		if self.Person.Win-self.Bets > 0 {
			self.Person.Cost = int(math.Ceil(float64(self.Person.Win-self.Bets) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
			self.Person.Win -= self.Person.Cost
			GetServer().SqlAgentGoldLog(self.Person.Uid, self.Person.Cost, self.room.Type)
			GetServer().SqlAgentBillsLog(self.Person.Uid, self.Person.Cost/2, self.room.Type)
		} else if self.Person.Win-self.Bets < 0 {
			cost := int(math.Ceil(float64(self.Person.Win-self.Bets) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 200.0))
			GetServer().SqlAgentBillsLog(self.Person.Uid, cost, self.room.Type)
		}

		self.Person.Total += self.Person.Win
	}
	dealwin := self.Bets - self.CurBet
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "--------　dealwin : ", dealwin)
	if dealwin != 0 {
		GetServer().SqlBZWLog(&SQL_BZWLog{1, dealwin, time.Now().Unix(), self.room.Type})

		if dealwin > 0 {
			cost := 0
			cost = int(math.Ceil(float64(dealwin) * float64(lib.GetManyMgr().GetProperty(self.room.Type).DealCost) / 100.0))
			dealwin -= cost
		}
		GetServer().SetDfwSysMoney(self.room.Type%170000, GetServer().DfwSysMoney[self.room.Type%170000]+int64(dealwin))
	}

	self.Time = time.Now().Unix() + 3600
}

func (self *Game_DFW) OnEnd() {
	self.room.Begin = false
	var msg Msg_GameDFW_End
	msg.Index = self.Index
	msg.Remaining = self.Remaining
	if DFW_MAP[self.Index] == 1 {
		msg.EndType = 1
	} else if msg.Remaining == 0 {
		msg.EndType = 2
	} else {
		msg.EndType = 3
	}
	msg.Person.Uid = self.Person.Uid
	msg.Person.Sex = self.Person.Sex
	msg.Person.Name = self.Person.Name
	msg.Person.IP = self.Person.IP
	msg.Person.Head = self.Person.Head
	msg.Person.Address = self.Person.Address

	msg.Person.Win = self.Person.Win
	msg.Person.Total = self.Person.Total
	self.room.SendMsg(self.Person.Uid, "gamedfwend", &msg)

	var record Rec_DFW_Info
	record.GameType = self.room.Type
	record.Time = time.Now().Unix()
	var rec Son_Rec_DFW_Person
	rec.Uid = self.Person.Uid
	rec.Name = self.Person.Name
	rec.Head = self.Person.Head
	rec.Bets = self.Bets
	rec.Score = self.Person.Win - self.Bets
	record.Info = append(record.Info, rec)
	GetServer().InsertRecord(self.room.Type, self.Person.Uid, lib.HF_JtoA(&record), rec.Score)

	self.Time = time.Now().Unix() + 60
}

func (self *Game_DFW) OnInit(room *Room) {
	self.room = room
}

func (self *Game_DFW) OnRobot(robot *lib.Robot) {

}

func (self *Game_DFW) OnSendInfo(person *Person) {
	if self.Person != nil && self.Person.Uid == person.Uid {
		self.Person.SynchroGold(person.Gold)
		person.SendMsg("gamedfwinfo", self.getinfo(person.Uid))
		return
	}

	_person := new(Game_DFW_Person)
	_person.Uid = person.Uid
	_person.Gold = person.Gold
	_person.Total = person.Gold
	_person.Name = person.Name
	_person.Sex = person.Sex
	_person.IP = person.ip
	_person.Head = person.Imgurl
	_person.Address = person.minfo.Address
	self.Person = _person
	person.SendMsg("gamedfwinfo", self.getinfo(person.Uid))

	self.Time = time.Now().Unix() + 60
}

func (self *Game_DFW) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold": //! 同步金币
		if self.Person.Uid == msg.V.(*staticfunc.Msg_SynchroGold).Uid {
			self.Person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal(self.Person.Uid, self.Person.Total)
		}
	case "gamebets":
		self.GameBets(msg.Uid, msg.V.(*Msg_GameBets).Bets)
	case "gamezhi":
		self.GameZhi(msg.Uid)
	}
}

func (self *Game_DFW) OnBye() {

}

func (self *Game_DFW) OnExit(uid int64) {
	if uid == self.Person.Uid {
		gold := self.Person.Total - self.Person.Gold
		if gold > 0 {
			GetRoomMgr().AddCard(self.Person.Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
		} else if gold < 0 {
			GetRoomMgr().CostCard(self.Person.Uid, staticfunc.TYPE_GOLD, -gold, self.room)
		}
		self.Person.Gold = self.Person.Total
	}
}

func (self *Game_DFW) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_DFW) OnBalance() {
	if self.Person == nil {
		return
	}

	gold := self.Person.Total - self.Person.Gold
	if gold > 0 {
		GetRoomMgr().AddCard(self.Person.Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
	} else if gold < 0 {
		GetRoomMgr().CostCard(self.Person.Uid, staticfunc.TYPE_GOLD, -gold, self.room)
	}
	self.Person.Gold = self.Person.Total
}

func (self *Game_DFW) OnTime() {
	if self.Person == nil {
		return
	}
	if self.Time == 0 {
		return
	}
	if time.Now().Unix() >= self.Time {
		self.room.KickViewByUid(self.Person.Uid, 96)
	}
}
