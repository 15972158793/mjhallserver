package gameserver

import (
	"fmt"
	"lib"
	"math"
	"staticfunc"
	"time"
)

var YYBF_MIN_GOLD int = 100
var YYBF_MIN_BET int = 100
var YYBF_BET_TIME int = 30

type Rec_YYBF_Info struct {
	GameType int                   `json:"gametype"`
	Time     int64                 `json:"time"` //! 记录时间
	Info     []Son_Rec_YYBF_Person `json:"info"`
}
type Son_Rec_YYBF_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Bets  int    `json:"bets"`
	Score int    `json:"score"`
}

type Game_GoldYYBF struct {
	PersonMgr map[int64]*Game_GoldYYBF_Person `json:"personmgr"`
	Bets      map[*Game_GoldYYBF_Person]int   `json:"bets"`
	Time      int64                           `json:"time"`
	Uids      []int64                         `json:"uids"`
	Total     int                             `json:"total"` //! 这一局一共下了多少钱
	Trend     []Son_GameGoldYYBF_Trend        `json:"trend"`
	Rating    []Son_GameGoldYYBF_Rating       `json:"rating"`

	room *Room
}
type Son_GameGoldYYBF_Trend struct {
	Uid  int64   `json:"uid"`
	Name string  `json:"name"`
	Head string  `json:"head"`
	Prob float64 `json:"prob"`
	Win  int     `json:"win"`
}

func NewGame_GoldYYBF() *Game_GoldYYBF {
	game := new(Game_GoldYYBF)
	game.PersonMgr = make(map[int64]*Game_GoldYYBF_Person)
	game.Bets = make(map[*Game_GoldYYBF_Person]int)
	game.Uids = make([]int64, 0)
	return game
}

type Msg_GameGoldYYBF_List struct {
	Info []Son_GameGoldYYBF_Info `json:"info"`
}
type Son_GameGoldYYBF_Info struct {
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Total   int    `json:"total"`
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
	Bets    int    `json:"bets"`
}

type Game_GoldYYBF_Person struct {
	Uid     int64                      `json:"uid"`
	Gold    int                        `json:"gold"`
	Total   int                        `json:"total"`
	Win     int                        `json:"win"`    //! 本局赢了多少
	Cost    int                        `json:"cost"`   //! 抽水钱
	Bets    int                        `json:"bets"`   //! 本局下了多少钱
	BeBets  int                        `json:"bebets"` //! 上一局下了多少
	Name    string                     `json:"name"`   //! 名字
	Head    string                     `json:"head"`   //! 头像
	Online  bool                       `json:"Online"`
	Round   int                        `json:"round"` //! 不下注轮数
	IP      string                     `json:"ip"`
	Address string                     `json:"address"`
	Sex     int                        `json:"sex"`
	History []Msg_GameGoldYYBF_History `json:"history"` //! 记录
}

type Msg_GameGoldYYBF_HistoryList struct {
	Info []Msg_GameGoldYYBF_History `json:"info"`
}

type Msg_GameGoldYYBF_History struct {
	Uid       int64   `json:"uid"`
	Time      int64   `json:"time"`
	Bets      int     `json:"bets"`
	Prob      float64 `json:"prob"`
	GameTotal int     `json:"gametotal"`
	Win       int     `json:"win"`
}

type Msg_GameGoldYYBF_Info struct {
	Begin     bool    `json:"begin"`     //! 是否开始游戏
	Time      int64   `json:"time"`      //! 倒计时
	Bets      int     `json:"bets"`      //! 自己下了多少注
	Total     int     `json:"total"`     //! 自己的钱
	GameTotal int     `json:"gametotal"` //! 奖池
	PerNum    int     `json:"pernum"`
	WinName   string  `json:"winname"`
	WinProb   float64 `json:"winprob"`
	Win       int     `json:"win"`
}

type Msg_GameGoldYYBF_Total struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"`
}

type Msg_GameGoldYYBF_Bets struct {
	Uid       int64  `json:"uid"`
	Name      string `json:"name"`
	Gold      int    `json:"gold"`      //! 本次下了多少注
	PerBets   int    `json:"perbets"`   //! 一共下了多少注
	Total     int    `json:"total"`     //! 还有多少钱
	GameTotal int    `json:"gametotal"` //! 奖池有多少钱
	PerNum    int    `json:"pernum"`    //! 有多少人下注
}

type Msg_GameGoldYYBF_End struct {
	Uid   int64   `json:"uid"`
	Win   int     `json:"win"`
	Bets  int     `json:"bets"`
	Prob  float64 `json:"prob"`
	Total int     `json:"total"`
	Head  string  `json:"head"`
	Name  string  `json:"name"`
}

type Msg_GameGoldYYBF_Rating struct {
	Info []Son_GameGoldYYBF_Rating `json:"info"`
}
type Son_GameGoldYYBF_Rating struct {
	Uid  int64  `json:"uid"`
	Name string `json:"name"`
	Head string `json:"head"`
	Win  int    `json:"win"`
}

func (self *Game_GoldYYBF) getinfo(uid int64, total int, bets int) *Msg_GameGoldYYBF_Info {
	var msg Msg_GameGoldYYBF_Info
	msg.Begin = self.room.Begin
	msg.Time = self.Time - time.Now().Unix()
	msg.Total = total
	msg.Bets = bets
	msg.GameTotal = self.Total
	msg.PerNum = len(self.Bets)
	if len(self.Trend) != 0 {
		msg.WinName = self.Trend[0].Name
		msg.WinProb = self.Trend[0].Prob
		msg.Win = self.Trend[0].Win
	}

	return &msg
}

//! 设置时间
func (self *Game_GoldYYBF) SetTime(t int) {
	if t == 0 {
		self.Time = 0
	} else {
		self.Time = time.Now().Unix() + int64(t)
	}
	var msg Msg_SetTime
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	self.room.broadCastMsg("gametime", &msg)
}

//! 同步金币
func (self *Game_GoldYYBF_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

func (self *Game_GoldYYBF) GetPerson(uid int64) *Game_GoldYYBF_Person {
	return self.PersonMgr[uid]
}

func (self *Game_GoldYYBF) GamePlayerList(uid int64) {
	var msg Msg_GameGoldYYBF_List
	msg.Info = make([]Son_GameGoldYYBF_Info, 0)
	for _, value := range self.PersonMgr {
		var son Son_GameGoldYYBF_Info
		son.Uid = value.Uid
		son.Head = value.Head
		son.IP = value.IP
		son.Name = value.Name
		son.Sex = value.Sex
		son.Total = value.Total
		son.Address = value.Address
		son.Bets = value.Bets
		msg.Info = append(msg.Info, son)
	}
	self.room.SendMsg(uid, "gameplayerlist", &msg)
}

//! 是否下注
func (self *Game_GoldYYBF) OnIsBets(uid int64) bool {
	value, ok := self.PersonMgr[uid]
	if ok {
		return value.Bets > 0
	}
	return false
}

func (self *Game_GoldYYBF) GameHistory(uid int64) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}
	var msg Msg_GameGoldYYBF_HistoryList
	msg.Info = person.History
	self.room.SendMsg(uid, "gamehistorylist", &msg)
}

func (self *Game_GoldYYBF) GameRating(uid int64) {
	for j := 0; j < len(self.Rating)-1; j++ {
		for i := 0; i < len(self.Rating)-1-j; i++ {
			if self.Rating[i].Win < self.Rating[i+1].Win {
				tmp := self.Rating[i]
				self.Rating[i] = self.Rating[i+1]
				self.Rating[i+1] = tmp
			}
		}
	}
	var msg Msg_GameGoldYYBF_Rating
	msg.Info = self.Rating
	if len(msg.Info) > 20 {
		msg.Info = msg.Info[0:20]
	}
	self.room.SendMsg(uid, "gamegoldyybfrating", &msg)
}

//! 同步总分
func (self *Game_GoldYYBF) SendTotal(uid int64, total int) {
	var msg Msg_GameGoldYYBF_Total
	msg.Uid = uid
	msg.Total = total

	person := self.GetPerson(uid)
	if person == nil {
		return
	}
	self.room.SendMsg(uid, "gamegoldtotal", &msg)
}

func (self *Game_GoldYYBF) GameGoOn(uid int64) {
	if uid == 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "------- uid == 0")
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() >= int64(YYBF_BET_TIME-2) {
		self.room.SendErr(uid, "正在开奖，请稍后下注。")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}
	if person.Total < YYBF_MIN_GOLD {
		self.room.SendErr(uid, fmt.Sprintf("%d金币以上才能下注。", YYBF_MIN_GOLD/100))
		return
	}

	if person.Total < person.BeBets {
		self.room.SendErr(uid, fmt.Sprintf("余额不足，不能续押"))
		return
	}

	person.Bets += person.BeBets
	person.Total -= person.BeBets
	person.Round = 0
	for i := 0; i < person.BeBets/YYBF_MIN_BET; i++ {
		self.Uids = append(self.Uids, uid)
	}
	self.Total += person.BeBets
	self.Bets[person] += person.BeBets

	var msg Msg_GameGoldYYBF_Bets
	msg.Uid = uid
	msg.Gold = person.BeBets
	msg.Name = person.Name
	msg.PerBets = person.Bets
	msg.Total = person.Total
	msg.GameTotal = self.Total
	msg.PerNum = len(self.Bets)
	self.room.broadCastMsg("gameyybfbets", &msg)

}

//! 下注
func (self *Game_GoldYYBF) GameBets(uid int64, gold int) {
	if uid == 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "------- uid == 0")
		return
	}
	if gold <= 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "------- gold == 0")
		return
	}
	if self.Time != 0 && self.Time-time.Now().Unix() >= int64(YYBF_BET_TIME-2) {
		self.room.SendErr(uid, "正在开奖，请稍后下注。")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if gold < YYBF_MIN_BET {
		self.room.SendErr(uid, fmt.Sprintf("下注不能少于%d金币。", YYBF_MIN_BET/100))
		return
	}

	if person.Total < YYBF_MIN_GOLD {
		self.room.SendErr(uid, fmt.Sprintf("%d金币以上才能下注。", YYBF_MIN_GOLD/100))
		return
	}

	if person.Total < gold {
		self.room.SendErr(uid, "您的金币不足，请及时充值。")
		return
	}

	person.Bets += gold
	person.Total -= gold
	person.Round = 0
	for i := 0; i < gold/YYBF_MIN_BET; i++ {
		self.Uids = append(self.Uids, uid)
	}
	self.Total += gold
	self.Bets[person] += gold
	var msg Msg_GameGoldYYBF_Bets
	msg.Uid = uid
	msg.Gold = gold
	msg.Name = person.Name
	msg.PerBets = person.Bets
	msg.Total = person.Total
	msg.GameTotal = self.Total
	msg.PerNum = len(self.Bets)
	self.room.broadCastMsg("gameyybfbets", &msg)

}

//! 一局开始
func (self *Game_GoldYYBF) OnBegin() {
	self.Total = 0
	self.room.Begin = true
	self.Bets = make(map[*Game_GoldYYBF_Person]int)
	self.Uids = make([]int64, 0)
}

//! 一局结束
func (self *Game_GoldYYBF) OnEnd() {
	if len(self.Bets) < 2 {
		if len(self.Bets) == 1 {
			for key, value := range self.Bets {
				key.Total += value
				key.Bets = 0
				self.SendTotal(key.Uid, key.Total)
			}
		}
		//! 清理玩家
		for key, value := range self.PersonMgr {
			if value.Online {
				value.Round++
				if value.Round >= 5 {
					self.room.KickViewByUid(value.Uid, 96)
				} else {
					value.Bets = 0
					value.Win = 0
					value.Cost = 0
					continue
				}
			}
			delete(self.PersonMgr, key)
		}
		self.SetTime(30)
		var msg Msg_GameGoldYYBF_End
		self.room.broadCastMsg("gamegoldyybfend", &msg)
		self.OnBegin()
		return
	}
	self.room.Begin = false

	winner := self.GetPerson(self.Uids[lib.HF_GetRandom(len(self.Uids))])
	if winner == nil {
		return
	}

	winner.Win = self.Total
	winner.Cost = int(math.Ceil(float64(winner.Win-winner.Bets) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
	GetServer().SqlAgentGoldLog(winner.Uid, winner.Cost, self.room.Type)
	winner.Total += winner.Win - winner.Cost
	find := false
	for i := 0; i < len(self.Rating); i++ {
		if self.Rating[i].Uid == winner.Uid {
			find = true
			self.Rating[i].Win += winner.Win - winner.Bets - winner.Cost
		}
	}
	if !find {
		self.Rating = append(self.Rating, Son_GameGoldYYBF_Rating{winner.Uid, winner.Name, winner.Head, winner.Win - winner.Bets - winner.Cost})
	}

	tmp := []Son_GameGoldYYBF_Trend{{winner.Uid, winner.Name, winner.Head, float64(winner.Bets) / float64(self.Total), winner.Win - winner.Cost}}
	tmp = append(tmp, self.Trend...)
	if len(tmp) > 20 {
		tmp = tmp[0:20]
	}
	self.Trend = tmp

	for _, value := range self.PersonMgr {
		if value.Bets > 0 {
			var record Rec_YYBF_Info
			record.Time = time.Now().Unix()
			record.GameType = self.room.Type
			var rec Son_Rec_YYBF_Person
			rec.Uid = value.Uid
			rec.Head = value.Head
			rec.Name = value.Name
			rec.Score = value.Win - value.Bets
			rec.Bets = value.Bets
			record.Info = append(record.Info, rec)
			GetServer().InsertRecord(self.room.Type, value.Uid, lib.HF_JtoA(&record), rec.Score)

			var history Msg_GameGoldYYBF_History
			history.Uid = value.Uid
			history.Bets = value.Bets
			history.Prob = float64(value.Bets) / float64(self.Total)
			history.Time = time.Now().Unix()
			history.Win = value.Win - value.Cost - value.Bets
			if history.Win < 0 {
				history.Win = 0
			}
			history.GameTotal = self.Total
			tmp := make([]Msg_GameGoldYYBF_History, 0)
			tmp = append(tmp, history)
			tmp = append(tmp, value.History...)
			if len(tmp) > 20 {
				tmp = tmp[0:20]
			}
			value.History = tmp
		}

		value.BeBets = value.Bets
	}

	var msg Msg_GameGoldYYBF_End
	msg.Uid = winner.Uid
	msg.Win = self.Total - winner.Cost
	msg.Bets = winner.Bets
	msg.Total = winner.Total
	msg.Head = winner.Head
	msg.Name = winner.Name
	msg.Prob = float64(winner.Bets) / float64(self.Total)
	self.room.broadCastMsg("gamegoldyybfend", &msg)

	//! 清理玩家
	for key, value := range self.PersonMgr {
		if value.Online {
			value.Round++
			if value.Round >= 5 {
				self.room.KickViewByUid(value.Uid, 96)
			} else {
				value.Bets = 0
				value.Win = 0
				value.Cost = 0
				continue
			}
		}
		delete(self.PersonMgr, key)
	}

	self.SetTime(30)
	self.OnBegin()
}

//! 初始化
func (self *Game_GoldYYBF) OnInit(room *Room) {
	self.room = room
}

func (self *Game_GoldYYBF) OnRobot(robot *lib.Robot) {

}

//! 告诉玩家数据
func (self *Game_GoldYYBF) OnSendInfo(person *Person) {
	value, ok := self.PersonMgr[person.Uid]
	if ok {
		value.Online = true
		value.Round = 0
		value.IP = person.ip
		value.Address = person.minfo.Address
		value.Sex = person.Sex
		value.SynchroGold(person.Gold)
		person.SendMsg("gamegoldyybfinfo", self.getinfo(person.Uid, value.Total, value.Bets))
		return
	}

	_person := new(Game_GoldYYBF_Person)
	_person.Uid = person.Uid
	_person.Gold = person.Gold
	_person.Total = person.Gold
	_person.Online = true
	_person.IP = person.ip
	_person.Address = person.minfo.Address
	_person.Sex = person.Sex
	_person.Name = person.Name
	_person.Head = person.Imgurl
	_person.History = make([]Msg_GameGoldYYBF_History, 0)
	self.PersonMgr[person.Uid] = _person
	person.SendMsg("gamegoldyybfinfo", self.getinfo(person.Uid, person.Gold, 0))

	if self.Time == 0 {
		self.SetTime(15)
		self.OnBegin()
	}
}

//! 消息转发
func (self *Game_GoldYYBF) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold":
		person := self.GetPerson(msg.V.(*staticfunc.Msg_SynchroGold).Uid)
		if person != nil {
			person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal(person.Uid, person.Total)
		}
	case "gamebets":
		self.GameBets(msg.Uid, msg.V.(*Msg_GameBets).Bets)
	case "gameplayerlist":
		self.GamePlayerList(msg.Uid)
	case "gamehistory": //!获取个人战绩
		self.GameHistory(msg.Uid)
	case "gamebzwgoon": //!续押
		self.GameGoOn(msg.Uid)
	case "gamerating": //! 获取排行榜
		self.GameRating(msg.Uid)
	}
}

//! 游戏结算
func (self *Game_GoldYYBF) OnBye() {

}

//! 玩家退出
func (self *Game_GoldYYBF) OnExit(uid int64) {
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

//! 是否是庄家
func (self *Game_GoldYYBF) OnIsDealer(uid int64) bool {
	return false
}

//! 结算
func (self *Game_GoldYYBF) OnBalance() {
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

//! 每秒调用一次
func (self *Game_GoldYYBF) OnTime() {
	if time.Now().Hour() == 0 && time.Now().Minute() == 0 && time.Now().Second() == 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "清空排行榜")
		self.Rating = make([]Son_GameGoldYYBF_Rating, 0)
	}

	if self.Time == 0 {
		return
	}

	if time.Now().Unix() < self.Time {
		return
	}

	if self.room.Begin {
		self.OnEnd()
		return
	}
}
