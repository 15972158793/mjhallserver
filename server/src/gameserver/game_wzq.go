package gameserver

import (
	"lib"
	"staticfunc"
	"time"
)

type C2S_GameWZQ_Step struct {
	W int `json:"w"`
	H int `json:"h"`
}

type S2C_GameWZQ_Step struct {
	W  int `json:"w"`
	H  int `json:"h"`
	Zi int `json:"zi"`
}

//!
type Msg_GameWZQ_Info struct {
	Begin   bool                                 `json:"begin"` //! 是否开始
	Info    []Son_GameWZQ_Info                   `json:"info"`
	Board   [GAME_WZQ_WIDTH][GAME_WZQ_HEIGHT]int `json:"board"`
	CurStep int64                                `json:"curstep"` //! 当前谁出牌
	Gold    int                                  `json:"gold"`    //! 学费
}

type Son_GameWZQ_Info struct {
	Uid   int64 `json:"uid"`
	Black bool  `json:"black"`
	Total int   `json:"total"`
}

type Msg_GameWZQ_Bye struct {
	Winer int64 `json:"winer"`
	Gold  int   `json:"gold"`
}

///////////////////////////////////////////////////////
const GAME_WZQ_WIDTH = 15
const GAME_WZQ_HEIGHT = 15

type Game_WZQ_Person struct {
	Uid   int64 `json:"uid"`
	Black bool  `json:"black"`
	Total int   `json:"total"`
}

//func (self *Game_WZQ_Person) Init() {
//	self.Black = false
//}

type Game_WZQ struct {
	PersonMgr []*Game_WZQ_Person                   `json:"personmgr"` //! 玩家信息
	Info      [GAME_WZQ_WIDTH][GAME_WZQ_HEIGHT]int `json:"info"`      //! 当前棋盘
	CurStep   int64                                `json:"curstep"`
	Gold      int                                  `json:"gold"` //! 学费
	Win       int64                                `json:"win"`  //! 赢家

	room *Room
}

func NewGame_WZQ() *Game_WZQ {
	game := new(Game_WZQ)
	game.PersonMgr = make([]*Game_WZQ_Person, 0)

	return game
}

func (self *Game_WZQ) OnInit(room *Room) {
	self.room = room
}

func (self *Game_WZQ) OnRobot(robot *lib.Robot) {

}

func (self *Game_WZQ) OnSendInfo(person *Person) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			person.SendMsg("gamewzqinfo", self.getInfo(person.Uid))
			return
		}
	}

	_person := new(Game_WZQ_Person)
	_person.Uid = person.Uid
	_person.Total = lib.HF_MaxInt(person.Gold-person.BindGold, 0)
	_person.Black = (len(self.PersonMgr) == 0)
	self.PersonMgr = append(self.PersonMgr, _person)

	person.SendMsg("gamewzqinfo", self.getInfo(person.Uid))
}

func (self *Game_WZQ) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "gamewzqstep": //! 叫分
		self.GameStep(msg.Uid, msg.V.(*C2S_GameWZQ_Step).W, msg.V.(*C2S_GameWZQ_Step).H)
	case "gamelose": //! 认输
		self.Lose(msg.Uid)
	case "gamebets": //! 设置学费
		self.SetGold(msg.Uid, msg.V.(*Msg_GameBets).Bets)
	case "gamebegin":
		self.GameBegin(msg.Uid)
	case "gameend":
		self.GameEnd(msg.Uid)
	}
}

func (self *Game_WZQ) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)

	self.CurStep = self.PersonMgr[0].Uid

	self.PersonMgr[0].Black = true
	self.PersonMgr[1].Black = false

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gamewzqbegin", self.getInfo(person.Uid))
	}

	self.room.flush()
}

//! 下子
func (self *Game_WZQ) GameStep(uid int64, w int, h int) {
	if !self.room.Begin {
		return
	}

	if self.CurStep != uid {
		return
	}

	if w < 0 || w >= GAME_WZQ_WIDTH {
		return
	}

	if h < 0 || h >= GAME_WZQ_HEIGHT {
		return
	}

	if self.Info[w][h] != 0 {
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	zi := 1
	if !person.Black {
		zi = 2
	}
	self.Info[w][h] = zi

	var msg S2C_GameWZQ_Step
	msg.W = w
	msg.H = h
	msg.Zi = zi
	self.room.broadCastMsg("gamewzqstep", &msg)

	//! 上下
	{
		num := 1
		for i := 1; ; i++ {
			if h-i < 0 {
				break
			}
			if self.Info[w][h-i] != zi {
				break
			}
			num++
		}
		for i := 1; ; i++ {
			if h+i >= GAME_WZQ_HEIGHT {
				break
			}
			if self.Info[w][h+i] != zi {
				break
			}
			num++
		}
		if num >= 5 {
			self.Win = uid
			self.OnEnd()
			return
		}
	}

	//! 左右
	{
		num := 1
		for i := 1; ; i++ {
			if w-i < 0 {
				break
			}
			if self.Info[w-i][h] != zi {
				break
			}
			num++
		}
		for i := 1; ; i++ {
			if w+i >= GAME_WZQ_WIDTH {
				break
			}
			if self.Info[w+i][h] != zi {
				break
			}
			num++
		}
		if num >= 5 {
			self.Win = uid
			self.OnEnd()
			return
		}
	}

	//! 右斜
	{
		num := 1
		for i := 1; ; i++ {
			if w-i < 0 {
				break
			}
			if h-i < 0 {
				break
			}
			if self.Info[w-i][h-i] != zi {
				break
			}
			num++
		}
		for i := 1; ; i++ {
			if w+i >= GAME_WZQ_WIDTH {
				break
			}
			if h+i >= GAME_WZQ_HEIGHT {
				break
			}
			if self.Info[w+i][h+i] != zi {
				break
			}
			num++
		}
		if num >= 5 {
			self.Win = uid
			self.OnEnd()
			return
		}
	}

	//! 左斜
	{
		num := 1
		for i := 1; ; i++ {
			if w-i < 0 {
				break
			}
			if h+i >= GAME_WZQ_HEIGHT {
				break
			}
			if self.Info[w-i][h+i] != zi {
				break
			}
			num++
		}
		for i := 1; ; i++ {
			if w+i >= GAME_WZQ_WIDTH {
				break
			}
			if h-i < 0 {
				break
			}
			if self.Info[w+i][h-i] != zi {
				break
			}
			num++
		}
		if num >= 5 {
			self.Win = uid
			self.OnEnd()
			return
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid != self.CurStep {
			self.CurStep = self.PersonMgr[i].Uid
			break
		}
	}

	self.room.flush()
}

//! 认输
func (self *Game_WZQ) Lose(uid int64) {
	if !self.room.Begin {
		return
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid != uid {
			self.Win = self.PersonMgr[i].Uid
			break
		}
	}

	self.OnEnd()
}

//! 开始
func (self *Game_WZQ) GameBegin(uid int64) {
	if self.room.Begin {
		return
	}

	if self.room.Host != uid {
		return
	}

	if self.Gold == 0 {
		return
	}

	if len(self.PersonMgr) < lib.HF_Atoi(self.room.csv["minnum"]) {
		return
	}

	self.OnBegin()
}

//! 结束
func (self *Game_WZQ) GameEnd(uid int64) {
	if self.room.Host != uid {
		return
	}

	self.room.Bye()
}

//! 设置学费
func (self *Game_WZQ) SetGold(uid int64, gold int) {
	if self.room.Begin {
		return
	}

	//if gold < 5000 || gold > 100000 {
	//	return
	//}

	if self.room.Host != uid {
		return
	}

	self.Gold = gold

	var msg Msg_GameBets
	msg.Uid = uid
	msg.Bets = gold
	self.room.broadCastMsg("gamebets", &msg)
}

//! 结算
func (self *Game_WZQ) OnEnd() {
	self.room.SetBegin(false)

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid != self.Win {
			if self.PersonMgr[i].Total < self.Gold { //! 钱不够
				self.Gold = 0
			} else if !GetRoomMgr().CostCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, self.Gold, self.room) {
				self.Gold = 0
			}
			break
		}
	}

	if self.Gold > 0 {
		var record staticfunc.Rec_Gold_Info
		record.Time = time.Now().Unix()
		record.GameType = self.room.Type

		wzqlog := new(SQL_WZQLog)
		wzqlog.Id = 1
		wzqlog.Gold = self.Gold
		wzqlog.Time = time.Now().Unix()
		for i := 0; i < len(self.PersonMgr); i++ {
			var rec staticfunc.Son_Rec_Gold_Person
			rec.Uid = self.PersonMgr[i].Uid
			rec.Name = self.room.GetName(self.PersonMgr[i].Uid)
			rec.Head = self.room.GetHead(self.PersonMgr[i].Uid)

			if self.PersonMgr[i].Uid == self.Win {
				wzqlog.Uid2 = self.PersonMgr[i].Uid
				self.PersonMgr[i].Total += self.Gold
				GetRoomMgr().AddCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, self.Gold, self.room.Type)
				rec.Score = self.Gold
			} else {
				wzqlog.Uid1 = self.PersonMgr[i].Uid
				self.PersonMgr[i].Total -= self.Gold
				rec.Score = -self.Gold
			}

			record.Info = append(record.Info, rec)
		}

		self.room.AddRecord(lib.HF_JtoA(&record))
		GetServer().SqlWZQLog(wzqlog)
	}

	var msg Msg_GameWZQ_Bye
	msg.Winer = self.Win
	msg.Gold = self.Gold
	self.room.broadCastMsg("gamewzqend", &msg)

	self.Gold = 0
	for i := 0; i < GAME_WZQ_WIDTH; i++ {
		for j := 0; j < GAME_WZQ_HEIGHT; j++ {
			self.Info[i][j] = 0
		}
	}

	return
}

func (self *Game_WZQ) OnBye() {
}

func (self *Game_WZQ) OnExit(uid int64) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			copy(self.PersonMgr[i:], self.PersonMgr[i+1:])
			self.PersonMgr = self.PersonMgr[:len(self.PersonMgr)-1]
			break
		}
	}
}

func (self *Game_WZQ) getInfo(uid int64) *Msg_GameWZQ_Info {
	var msg Msg_GameWZQ_Info
	msg.Begin = self.room.Begin
	msg.CurStep = self.CurStep
	msg.Gold = self.Gold
	msg.Board = self.Info
	msg.Info = make([]Son_GameWZQ_Info, 0)
	for _, value := range self.PersonMgr {
		var son Son_GameWZQ_Info
		son.Uid = value.Uid
		son.Black = value.Black
		son.Total = value.Total
		msg.Info = append(msg.Info, son)
	}
	return &msg
}

func (self *Game_WZQ) GetPerson(uid int64) *Game_WZQ_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

func (self *Game_WZQ) OnTime() {

}

func (self *Game_WZQ) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_WZQ) OnIsBets(uid int64) bool {
	return false
}

//! 结算所有人
func (self *Game_WZQ) OnBalance() {
}
