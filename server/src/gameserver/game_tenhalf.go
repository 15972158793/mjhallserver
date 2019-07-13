package gameserver

import (
	"lib"
	"staticfunc"
	"time"
)

//type LstNNCard [][]int

//func (a LstNNCard) Len() int { // 重写 Len() 方法
//	return len(a)
//}

//func (a LstNNCard) Swap(i, j int) { // 重写 Swap() 方法
//	a[i], a[j] = a[j], a[i]
//}

//func (a LstNNCard) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
//	cti, csi := GetNiuNiuScore(a[i])
//	ctj, csj := GetNiuNiuScore(a[j])
//	if cti > ctj {
//		return true
//	} else if cti < ctj {
//		return false
//	} else {
//		return csi > csj
//	}
//}

//! 记录结构
type Rec_GameTenHalf struct {
	Info    []Son_Rec_GameTenHalf `json:"info"`
	Roomid  int                   `json:"roomid"`
	Time    int64                 `json:"time"`
	MaxStep int                   `json:"maxstep"`
}
type Son_Rec_GameTenHalf struct {
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`
	Head   string `json:"head"`
	Card   []int  `json:"card"`
	Bets   int    `json:"bets"`
	Dealer bool   `json:"dealer"`
	Score  int    `json:"score"`
	Total  int    `json:"total"`
}

//! 游戏操作
type Msg_GameTenHalfPlay struct {
	Uid     int64 `json:"uid"`
	Type    int   `json:"type"`
	Card    []int `json:"card"`
	CurStep int64 `json:"curstep"`
}

//!
type Msg_GameTenHalf_Info struct {
	Begin   bool                   `json:"begin"` //! 是否开始
	Ready   []int64                `json:"ready"` //! 准备的人
	Info    []Son_GameTenHalf_Info `json:"info"`
	CurStep int64                  `json:"curstep"` //! 该谁操作
}
type Son_GameTenHalf_Info struct {
	Uid    int64 `json:"uid"`    //! uid
	Card   []int `json:"card"`   //! 手牌
	Bets   int   `json:"bets"`   //! 底分
	Dealer bool  `json:"dealer"` //! 是否庄家
	Stop   bool  `json:"stop"`   //! 是否停牌
	Score  int   `json:"score"`  //! 当局分数
	Total  int   `json:"total"`  //! 总分
}

//! 下注
type Msg_GameTenHalfBets struct {
	Uid  int64 `json:"uid"`
	Bets int   `json:"bets"`
	Card int   `json:"card"`
}

//! 结算
type Msg_GameTenHalf_End struct {
	Info []Son_GameTenHalf_Info `json:"info"`
}

//! 房间结束
type Msg_GameTenHalf_Bye struct {
	Info []Son_GameTenHalf_Bye `json:"info"`
}
type Son_GameTenHalf_Bye struct {
	Uid   int64 `json:"uid"`
	Win   int   `json:"win"`   //! 胜利次数
	Lose  int   `json:"lose"`  //! 失败次数
	Deal  int   `json:"deal"`  //! 坐庄次数
	TW    int   `json:"tw"`    //! 天王次数
	HWX   int   `json:"hwx"`   //! 花五小次数
	WX    int   `json:"wx"`    //! 五小次数
	SDB   int   `json:"sdb"`   //! 十点半次数
	GP    int   `json:"gp"`    //! 高牌次数
	Score int   `json:"score"` //! 总分
}

///////////////////////////////////////////////////////
type Game_TenHalf_Person struct {
	Uid      int64 `json:"uid"`      //! uid
	Card     []int `json:"card"`     //! 手牌
	Score    int   `json:"score"`    //! 积分
	Bets     int   `json:"bets"`     //! 下注
	Dealer   bool  `json:"dealer"`   //! 是否庄家
	CurScore int   `json:"curscore"` //! 当前局的分数
	Boom     bool  `json:"boom"`     //! 是否爆牌
	Stop     bool  `json:"stop"`     //! 是否停牌
	View     bool  `json:"view"`     //! 是否亮牌
	Win      int   `json:"win"`      //! 胜利次数
	Lose     int   `json:"lose"`     //! 失败次数
	Deal     int   `json:"deal"`     //! 坐庄次数
	TW       int   `json:"tw"`       //! 天王次数
	HWX      int   `json:"hwx"`      //! 花五小次数
	WX       int   `json:"wx"`       //! 五小次数
	SDB      int   `json:"sdb"`      //! 十点半次数
	GP       int   `json:"gp"`       //! 高牌次数
}

func (self *Game_TenHalf_Person) Init() {
	self.Card = make([]int, 0)
	self.Bets = 0
	self.Dealer = false
	self.CurScore = 0
	self.Boom = false
	self.Stop = false
	self.View = false
}

type Game_TenHalf struct {
	Ready     []int64                `json:"ready"` //! 已经准备的人
	PersonMgr []*Game_TenHalf_Person `json:"personmgr"`
	Bets      int                    `json:"bets"`    //! 底分
	CurStep   int64                  `json:"curstep"` //! 当前操作人
	Card      *CardMgr               `json:"card"`    //! 剩余牌

	room *Room
}

func NewGame_TenHalf() *Game_TenHalf {
	game := new(Game_TenHalf)
	game.Ready = make([]int64, 0)
	game.PersonMgr = make([]*Game_TenHalf_Person, 0)

	return game
}

func (self *Game_TenHalf) OnInit(room *Room) {
	self.room = room
}

func (self *Game_TenHalf) OnRobot(robot *lib.Robot) {

}

func (self *Game_TenHalf) OnSendInfo(person *Person) {
	person.SendMsg("gametenhalfinfo", self.getInfo(person.Uid))
}

func (self *Game_TenHalf) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "gameready": //! 游戏准备
		self.GameReady(msg.Uid)
	case "gamebets": //! 下注
		self.GameBets(msg.Uid, msg.V.(*Msg_GameBets).Bets)
	case "gameplay": //! 操作
		self.GamePlay(msg.Uid, msg.V.(*Msg_GamePlay).Type)
	}
}

func (self *Game_TenHalf) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)

	//! 庄家的位置
	WinPos := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].CurScore > self.PersonMgr[WinPos].CurScore {
			WinPos = i
		}
	}

	for i := 0; i < len(self.room.Uid); i++ { //! 重新初始化人
		if i >= len(self.PersonMgr) {
			person := new(Game_TenHalf_Person)
			person.Uid = self.room.Uid[i]
			self.PersonMgr = append(self.PersonMgr, person)
		} else {
			self.PersonMgr[i].Init()
		}
	}
	if self.room.Param1/100 == 0 { //! 房主庄
		self.PersonMgr[0].Dealer = true
		self.PersonMgr[0].Deal++
		self.CurStep = self.PersonMgr[0].Uid
	} else { //! 赢家庄
		self.PersonMgr[WinPos].Dealer = true
		self.PersonMgr[WinPos].Deal++
		self.CurStep = self.PersonMgr[WinPos].Uid
	}
	self.CurStep = self.GetNextUid()

	//! 用牛牛的牌组
	self.Card = NewCard_NiuNiu(false)
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Card = self.Card.Deal(1)
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gametenhalfbegin", self.getInfo(person.Uid))
	}

	self.room.flush()
}

//! 准备
func (self *Game_TenHalf) GameReady(uid int64) {
	if self.room.IsBye() {
		return
	}

	if self.room.Begin { //! 已经开始了不允许准备
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经开始了，不能准备")
		return
	}

	for i := 0; i < len(self.Ready); i++ {
		if self.Ready[i] == uid {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "同一个玩家准备")
			return
		}
	}

	self.Ready = append(self.Ready, uid)

	if len(self.Ready) == len(self.room.Uid) && len(self.Ready) >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 准备的人数达到游戏最小人数
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
		self.OnBegin()
		return
	}

	var msg staticfunc.Msg_Uid
	msg.Uid = uid
	self.room.broadCastMsg("gameready", &msg)

	self.room.flush()
}

//! 下注
func (self *Game_TenHalf) GameBets(uid int64, bets int) {
	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	if bets <= 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "下注无效")
		return
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Dealer { //! 是庄家
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "庄家不用下注")
				return
			}

			if self.PersonMgr[i].Bets > 0 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "重复下注")
				return
			}

			self.PersonMgr[i].Bets = bets
			break
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}

		var msg Msg_GameTenHalfBets
		msg.Uid = uid
		msg.Bets = bets
		if person.Uid == uid || GetServer().IsAdmin(uid, staticfunc.ADMIN_TENHALF) {
			msg.Card = self.PersonMgr[i].Card[0]
		} else {
			msg.Card = 0
		}

		person.SendMsg("gamebets", &msg)
	}

	self.room.flush()
}

//! 操作
func (self *Game_TenHalf) GamePlay(uid int64, _type int) {
	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	if self.CurStep != uid {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不归该玩家操作")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if person.Bets <= 0 && !person.Dealer {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "没有下注前不能操作")
		return
	}

	if len(person.Card) == 5 || person.Stop || person.View {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能操作")
		return
	}

	if _type == 0 { //! 要牌
		if person.Boom {
			return
		}
		person.Card = append(person.Card, self.Card.Deal(1)...)
		ct := GetTenHalfType(person.Card)
		if ct == -1 {
			person.Boom = true
			person.View = true
		} else if ct >= 100 {
			person.View = true
		}
	} else if _type == 1 { //! 停牌
		person.Stop = true
	}

	if person.Dealer && (person.Boom || len(person.Card) == 5 || person.Stop || person.View) {
		self.OnEnd()
		return
	}

	if len(person.Card) == 5 || person.Stop || person.View {
		self.CurStep = self.GetNextUid()
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		var msg Msg_GameTenHalfPlay
		msg.Uid = uid
		msg.Type = _type
		if person.View || self.PersonMgr[i].Uid == person.Uid {
			msg.Card = person.Card
		} else {
			msg.Card = append(msg.Card, 0)
			msg.Card = append(msg.Card, person.Card[1:]...)
		}
		msg.CurStep = self.CurStep
		self.room.SendMsg(self.PersonMgr[i].Uid, "gametenhalfplay", &msg)
	}

	self.room.flush()
}

//! 结算
func (self *Game_TenHalf) OnEnd() {
	self.room.SetBegin(false)

	var person *Game_TenHalf_Person = nil
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Dealer {
			person = self.PersonMgr[i]
			break
		}
	}

	dt := GetTenHalfType(person.Card)
	if dt == 1000 {
		person.TW++
	} else if dt == 500 {
		person.HWX++
	} else if dt == 200 {
		person.WX++
	} else if dt == 100 {
		person.SDB++
	} else if !person.Boom {
		person.GP++
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Dealer {
			continue
		}

		score := self.PersonMgr[i].Bets
		if self.PersonMgr[i].Boom { //! 爆牌了
			self.PersonMgr[i].CurScore -= score
			person.CurScore += score
			self.PersonMgr[i].Lose++
		} else {
			xt := GetTenHalfType(self.PersonMgr[i].Card)
			if xt == 1000 {
				self.PersonMgr[i].TW++
			} else if xt == 500 {
				self.PersonMgr[i].HWX++
			} else if xt == 200 {
				self.PersonMgr[i].WX++
			} else if xt == 100 {
				self.PersonMgr[i].SDB++
			} else if !self.PersonMgr[i].Boom {
				self.PersonMgr[i].GP++
			} else if self.PersonMgr[i].Boom {
				self.PersonMgr[i].Lose++
			}
			if dt >= xt { //! 庄家赢
				self.PersonMgr[i].CurScore -= score
				person.CurScore += score
			} else { //! 闲家赢
				score *= GetTenHalfBS(self.room.Param1%10, xt)
				self.PersonMgr[i].CurScore += score
				person.CurScore -= score
			}
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].CurScore > 0 {
			self.PersonMgr[i].Win++
		}
		self.PersonMgr[i].Score += self.PersonMgr[i].CurScore
	}

	self.Ready = make([]int64, 0)

	//! 记录
	var record Rec_GameTenHalf
	record.Time = time.Now().Unix()
	record.Roomid = self.room.Id*100 + self.room.Step
	record.MaxStep = self.room.MaxStep

	//! 发消息
	var msg Msg_GameTenHalf_End
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameTenHalf_Info
		son.Uid = self.PersonMgr[i].Uid
		son.Bets = self.PersonMgr[i].Bets
		son.Card = self.PersonMgr[i].Card
		son.Dealer = self.PersonMgr[i].Dealer
		son.Stop = self.PersonMgr[i].Stop
		son.Score = self.PersonMgr[i].CurScore
		son.Total = self.PersonMgr[i].Score
		msg.Info = append(msg.Info, son)

		var _son Son_Rec_GameTenHalf
		_son.Uid = self.PersonMgr[i].Uid
		_son.Name = self.room.GetName(self.PersonMgr[i].Uid)
		_son.Head = self.room.GetHead(self.PersonMgr[i].Uid)
		_son.Card = self.PersonMgr[i].Card
		_son.Bets = self.PersonMgr[i].Bets
		_son.Dealer = self.PersonMgr[i].Dealer
		_son.Score = self.PersonMgr[i].CurScore
		_son.Total = self.PersonMgr[i].Score
		record.Info = append(record.Info, _son)
	}
	self.room.AddRecord(lib.HF_JtoA(&record))
	self.room.broadCastMsg("gametenhalfend", &msg)

	if self.room.IsBye() {
		self.OnBye()
		self.room.Bye()
		return
	}

	self.room.flush()
}

func (self *Game_TenHalf) OnBye() {
	info := make([]staticfunc.JS_CreateRoomMem, 0)
	var msg Msg_GameTenHalf_Bye
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameTenHalf_Bye
		son.Uid = self.PersonMgr[i].Uid
		son.Win = self.PersonMgr[i].Win
		son.Lose = self.PersonMgr[i].Lose
		son.Deal = self.PersonMgr[i].Deal
		son.TW = self.PersonMgr[i].TW
		son.HWX = self.PersonMgr[i].HWX
		son.WX = self.PersonMgr[i].WX
		son.SDB = self.PersonMgr[i].SDB
		son.GP = self.PersonMgr[i].GP
		son.Score = self.PersonMgr[i].Score
		msg.Info = append(msg.Info, son)
		info = append(info, staticfunc.JS_CreateRoomMem{son.Uid, "", "", son.Score})

		GetServer().SqlScoreLog(self.PersonMgr[i].Uid, self.room.GetName(self.PersonMgr[i].Uid), self.room.GetHead(self.PersonMgr[i].Uid), self.room.Type, self.room.Id, self.PersonMgr[i].Score)
	}
	self.room.broadCastMsg("gametenhalfbye", &msg)

	self.room.ClubResult(info)
}

func (self *Game_TenHalf) OnExit(uid int64) {
	find := false
	for i := 0; i < len(self.Ready); i++ {
		if self.Ready[i] == uid {
			copy(self.Ready[i:], self.Ready[i+1:])
			self.Ready = self.Ready[:len(self.Ready)-1]
			find = true
			break
		}
	}

	if !find {
		if len(self.Ready) == len(self.room.Uid) && len(self.Ready) >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 准备的人数达到游戏最小人数
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
			self.OnBegin()
			return
		}
	}
}

func (self *Game_TenHalf) getInfo(uid int64) *Msg_GameTenHalf_Info {
	var msg Msg_GameTenHalf_Info
	msg.Begin = self.room.Begin
	msg.Ready = make([]int64, 0)
	msg.CurStep = self.CurStep
	msg.Info = make([]Son_GameTenHalf_Info, 0)
	if !msg.Begin { //! 没有开始,看哪些人已准备
		msg.Ready = self.Ready
	}
	for _, value := range self.PersonMgr {
		var son Son_GameTenHalf_Info
		son.Uid = value.Uid
		son.Bets = value.Bets
		son.Dealer = value.Dealer
		son.Stop = value.Stop
		son.Score = value.CurScore
		son.Total = value.Score
		if self.room.Begin {
			if (value.Uid == uid && (value.Bets > 0 || value.Dealer)) || value.View || GetServer().IsAdmin(uid, staticfunc.ADMIN_TENHALF) {
				son.Card = value.Card
			} else {
				son.Card = append(son.Card, 0)
				son.Card = append(son.Card, value.Card[1:]...)
			}
		} else {
			son.Card = value.Card
		}
		msg.Info = append(msg.Info, son)
	}
	return &msg
}

func (self *Game_TenHalf) GetPerson(uid int64) *Game_TenHalf_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

//! 得到下一个uid
func (self *Game_TenHalf) GetNextUid() int64 {
	dealer := int64(0)
	find := false
	i := 0
	for {
		if find {
			if !self.PersonMgr[i].Stop && !self.PersonMgr[i].View && !self.PersonMgr[i].Dealer && len(self.PersonMgr[i].Card) < 5 {
				return self.PersonMgr[i].Uid
			}
		}

		if self.PersonMgr[i].Dealer {
			dealer = self.PersonMgr[i].Uid
		}

		if self.PersonMgr[i].Uid == self.CurStep {
			if find { //! 已经找过这个人
				break
			}
			find = true
		}

		i++
		if i == len(self.PersonMgr) {
			i = 0
		}
	}

	return dealer
}

func (self *Game_TenHalf) OnTime() {

}

func (self *Game_TenHalf) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_TenHalf) OnIsBets(uid int64) bool {
	return false
}

//! 结算所有人
func (self *Game_TenHalf) OnBalance() {
}
