package gameserver

import (
	"lib"
	"sort"
	"staticfunc"
)

//! 记录结构
//type Rec_GameTenHalf struct {
//	Info   []Son_GameTenHalf_Info `json:"info"`
//	Roomid int                    `json:"roomid"`
//	Time   int64                  `json:"time"`
//}

//! 游戏操作
type Msg_GameEatHotPlay struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"`
}

//! 游戏操作
type Msg_GameEatHotResult struct {
	Uid      int64 `json:"uid"`
	Bets     int   `json:"bets"`
	Card     []int `json:"card"`
	NextCard []int `json:"nextcard"`
	CurScore int   `json:"curscore"`
	CurStep  int64 `json:"curstep"`
	Step     int   `json:"step"`
	HotBets  int   `json:"hotbets"`
	One      int64 `json:"one"`
}

//!
type Msg_GameEatHot_Info struct {
	Begin   bool                  `json:"begin"` //! 是否开始
	Ready   []int64               `json:"ready"` //! 准备的人
	Agree   []int64               `json:"agree"` //! 同意平分锅底的人
	Info    []Son_GameEatHot_Info `json:"info"`
	CurStep int64                 `json:"curstep"` //! 该谁操作
	Bets    int                   `json:"bets"`    //! 锅底
	Step    int                   `json:"step"`    //! 第几圈
	One     int64                 `json:"one"`     //! 第一个人
}
type Son_GameEatHot_Info struct {
	Uid   int64 `json:"uid"`   //! uid
	Card  []int `json:"card"`  //! 手牌
	Bets  int   `json:"bets"`  //! 下注
	Score int   `json:"score"` //! 当局分数
	Total int   `json:"total"` //! 总分
}

//! 结算
type Msg_GameEatHot_End struct {
	Info []Son_GameEatHot_Info `json:"info"`
}

//! 房间结束
type Msg_GameEatHot_Bye struct {
	Info []Son_GameEatHot_Bye `json:"info"`
}
type Son_GameEatHot_Bye struct {
	Uid   int64 `json:"uid"`
	Win   int   `json:"win"`   //! 胜利次数
	BZ    int   `json:"bz"`    //! 豹子次数
	SZ    int   `json:"sz"`    //! 顺子次数
	KF    int   `json:"kf"`    //! 卡飞次数
	ZZ    int   `json:"zz"`    //! 撞柱次数
	Score int   `json:"score"` //! 总分
}

///////////////////////////////////////////////////////
type Game_EatHot_Person struct {
	Uid      int64 `json:"uid"`      //! uid
	Card     []int `json:"card"`     //! 手牌
	Score    int   `json:"score"`    //! 积分
	Bets     int   `json:"bets"`     //! 下注
	CurScore int   `json:"curscore"` //! 当前局的分数
	Win      int   `json:"win"`      //! 胜利次数
	BZ       int   `json:"bz"`       //! 豹子次数
	SZ       int   `json:"sz"`       //! 顺子次数
	KF       int   `json:"kf"`       //! 卡飞次数
	ZZ       int   `json:"zz"`       //! 撞柱次数
}

type Game_EatHot struct {
	Ready     []int64               `json:"ready"` //! 已经准备的人
	Agree     []int64               `json:"agree"` //! 同意平分锅底的人
	PersonMgr []*Game_EatHot_Person `json:"personmgr"`
	Bets      int                   `json:"bets"`    //! 总注
	Step      int                   `json:"step"`    //! 第几圈
	CurStep   int64                 `json:"curstep"` //! 当前操作人
	Card      *CardMgr              `json:"card"`    //! 剩余牌
	First     int64                 `json:"firset"`  //! 第一个赢了钱的人
	One       int64                 `json:"one"`     //! 每一圈第一的人

	room *Room
}

func NewGame_EatHot() *Game_EatHot {
	game := new(Game_EatHot)
	game.Ready = make([]int64, 0)
	game.PersonMgr = make([]*Game_EatHot_Person, 0)

	return game
}

func (self *Game_EatHot) OnInit(room *Room) {
	self.room = room
}

func (self *Game_EatHot) OnRobot(robot *lib.Robot) {

}

func (self *Game_EatHot) OnSendInfo(person *Person) {
	person.SendMsg("gameeathotinfo", self.getInfo(person.Uid))
}

func (self *Game_EatHot) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "gameready": //! 游戏准备
		self.GameReady(msg.Uid)
	case "gamebets": //! 下注
		self.GameBets(msg.Uid, msg.V.(*Msg_GameBets).Bets)
	case "gameplay": //! 平分锅底
		self.GamePlay(msg.Uid, msg.V.(*Msg_GamePlay).Type)
	}
}

func (self *Game_EatHot) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)

	index := -1
	for i := 0; i < len(self.room.Uid); i++ { //! 重新初始化人
		if i >= len(self.PersonMgr) {
			person := new(Game_EatHot_Person)
			person.Uid = self.room.Uid[i]
			person.Card = make([]int, 0)
			self.PersonMgr = append(self.PersonMgr, person)
		} else {
			self.PersonMgr[i].Bets = 0
			self.PersonMgr[i].Card = make([]int, 0)
			self.PersonMgr[i].CurScore = 0
		}
		if self.PersonMgr[i].Uid == self.First {
			self.First = 0
			index = i
		}
	}

	//! 随机第一个人
	if index == -1 {
		index = lib.HF_GetRandom(len(self.PersonMgr))
	}
	self.CurStep = self.PersonMgr[index].Uid
	self.One = self.CurStep

	//! 用牛牛的牌组
	self.Agree = make([]int64, 0)
	self.Step = 0
	self.Bets = ((self.room.Param1 / 10 % 10) + 1) * len(self.PersonMgr)
	self.Card = NewCard_NiuNiu(false)
	self.PersonMgr[index].Card = self.Card.Deal(2)
	sort.Ints(self.PersonMgr[index].Card)

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].CurScore -= (self.room.Param1/10%10 + 1)
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gameeathotbegin", self.getInfo(person.Uid))
	}

	self.room.flush()
}

//! 准备
func (self *Game_EatHot) GameReady(uid int64) {
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
func (self *Game_EatHot) GameBets(uid int64, bets int) {
	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	person := self.GetPerson(uid)

	if bets > 0 {
		minbet := lib.HF_MinInt(self.Bets, (self.room.Param1/10%10)+1)
		if bets < minbet {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "下注无效1")
			return
		}

		maxbet := (self.room.Param1%10 + 1) * 50
		if person.Card[0]/10 == person.Card[1]/10 { //! 可能是豹子
			maxbet = 5
		} else if lib.HF_Abs(person.Card[0]/10-person.Card[1]/10) == 1 || lib.HF_Abs(person.Card[0]/10-person.Card[1]/10) == 2 { //! 可能是顺子
			maxbet = 5
		}

		maxbet = lib.HF_MinInt(self.Bets, maxbet)

		if bets > maxbet {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "下注无效2")
			return
		}

		person.Bets = bets
		person.Card = append(person.Card, self.Card.Deal(1)...)
		_card := make([]int, 0)
		lib.HF_DeepCopy(&_card, &person.Card)
		sort.Ints(_card)
		if _card[0]/10 == _card[1]/10 && _card[1]/10 == _card[2]/10 { //! 豹子
			person.BZ++
			bs := 1
			if self.room.Param1/100%10 == 1 {
				bs = 10
			} else if self.room.Param1/100%10 == 2 {
				bs = 20
			}
			score := lib.HF_MinInt(person.Bets*bs, self.Bets)
			person.CurScore += score
			self.Bets -= score
			if self.First == 0 {
				self.First = person.Uid
			}
		} else if _card[1]/10-_card[0]/10 == 1 && _card[2]/10-_card[0]/10 == 2 { //! 顺子
			person.SZ++
			bs := 1
			if self.room.Param1/1000%10 == 1 {
				bs = 5
			} else if self.room.Param1/1000%10 == 2 {
				bs = 10
			}
			score := lib.HF_MinInt(person.Bets*bs, self.Bets)
			person.CurScore += score
			self.Bets -= score
			if self.First == 0 {
				self.First = person.Uid
			}
		} else if person.Card[2]/10 > person.Card[0]/10 && person.Card[2]/10 < person.Card[1]/10 {
			score := lib.HF_MinInt(person.Bets, self.Bets)
			person.CurScore += score
			self.Bets -= score
			if self.First == 0 {
				self.First = person.Uid
			}
		} else { //! 负
			if person.Card[1]/10 == person.Card[0]/10 || person.Card[1]/10 == person.Card[2]/10 {
				person.ZZ++
			} else {
				person.KF++
			}
			score := lib.HF_MinInt(person.Bets, self.Bets)
			person.CurScore -= score
			self.Bets += score
		}

		if self.Bets <= 0 {
			self.OnEnd()
			return
		}
	}

	var msg Msg_GameEatHotResult
	msg.Uid = uid
	msg.Bets = bets
	msg.Card = person.Card
	msg.CurScore = person.CurScore

	self.Step++
	if self.Step%len(self.PersonMgr) == 0 {
		self.Card = NewCard_NiuNiu(false)
	}
	nextperson := self.GetNextPerson()
	nextperson.Card = self.Card.Deal(2)
	sort.Ints(nextperson.Card)
	self.CurStep = nextperson.Uid

	msg.NextCard = nextperson.Card
	msg.CurStep = self.CurStep
	msg.Step = self.Step / len(self.PersonMgr)
	msg.HotBets = self.Bets
	msg.One = self.One
	self.room.broadCastMsg("gameeathotresult", &msg)

	self.room.flush()
}

//! 平分锅底
func (self *Game_EatHot) GamePlay(uid int64, _type int) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	for i := 0; i < len(self.Agree); i++ {
		if self.Agree[i] == uid {
			return
		}
	}

	var msg Msg_GameEatHotPlay
	msg.Uid = uid
	msg.Type = _type
	self.room.broadCastMsg("gameeathotplay", &msg)

	if _type == 0 { //! 同意
		self.Agree = append(self.Agree, uid)
		if len(self.Agree) == len(self.PersonMgr) {
			s1 := self.Bets / len(self.PersonMgr)
			s2 := self.Bets % len(self.PersonMgr)
			for i := 0; i < len(self.PersonMgr); i++ {
				self.PersonMgr[i].CurScore += s1
			}
			if s2 > 0 {
				index := 0
				for i := 1; i < len(self.PersonMgr); i++ {
					if self.PersonMgr[i].CurScore < self.PersonMgr[index].CurScore {
						index = i
					}
				}
				self.PersonMgr[index].CurScore += s2
			}
			self.OnEnd()
			return
		}
	} else { //! 不同意
		self.Agree = make([]int64, 0)
	}

	self.room.flush()
}

//! 结算
func (self *Game_EatHot) OnEnd() {
	self.room.SetBegin(false)

	self.Ready = make([]int64, 0)

	//! 发消息
	var msg Msg_GameEatHot_End
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Score += self.PersonMgr[i].CurScore

		var son Son_GameEatHot_Info
		son.Uid = self.PersonMgr[i].Uid
		son.Score = self.PersonMgr[i].CurScore
		son.Total = self.PersonMgr[i].Score
		son.Card = self.PersonMgr[i].Card
		son.Bets = self.PersonMgr[i].Bets
		msg.Info = append(msg.Info, son)
	}
	self.room.broadCastMsg("gameeathotend", &msg)

	if self.room.IsBye() {
		self.OnBye()
		self.room.Bye()
		return
	}

	self.room.flush()
}

func (self *Game_EatHot) OnBye() {
	var msg Msg_GameEatHot_Bye
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameEatHot_Bye
		son.Uid = self.PersonMgr[i].Uid
		son.Win = self.PersonMgr[i].Win
		son.BZ = self.PersonMgr[i].BZ
		son.SZ = self.PersonMgr[i].SZ
		son.KF = self.PersonMgr[i].KF
		son.ZZ = self.PersonMgr[i].ZZ
		son.Score = self.PersonMgr[i].Score
		msg.Info = append(msg.Info, son)
	}
	self.room.broadCastMsg("gameeathotbye", &msg)
}

func (self *Game_EatHot) OnExit(uid int64) {
	for i := 0; i < len(self.Ready); i++ {
		if self.Ready[i] == uid {
			copy(self.Ready[i:], self.Ready[i+1:])
			self.Ready = self.Ready[:len(self.Ready)-1]
			break
		}
	}
}

func (self *Game_EatHot) getInfo(uid int64) *Msg_GameEatHot_Info {
	var msg Msg_GameEatHot_Info
	msg.Begin = self.room.Begin
	msg.One = self.One
	msg.Ready = make([]int64, 0)
	msg.CurStep = self.CurStep
	msg.Bets = self.Bets
	if len(self.PersonMgr) == 0 {
		msg.Step = 1
	} else {
		msg.Step = self.Step/len(self.PersonMgr) + 1
	}
	msg.Agree = self.Agree
	msg.Info = make([]Son_GameEatHot_Info, 0)
	if !msg.Begin { //! 没有开始,看哪些人已准备
		msg.Ready = self.Ready
	}
	for _, value := range self.PersonMgr {
		var son Son_GameEatHot_Info
		son.Uid = value.Uid
		son.Bets = value.Bets
		son.Card = value.Card
		son.Score = value.CurScore
		son.Total = value.Score
		msg.Info = append(msg.Info, son)
	}
	return &msg
}

func (self *Game_EatHot) GetPerson(uid int64) *Game_EatHot_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

//! 得到下一个uid
func (self *Game_EatHot) GetNextPerson() *Game_EatHot_Person {
	if self.Step%len(self.PersonMgr) == 0 && self.First != 0 {
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Uid != self.First {
				continue
			}

			self.One = self.First
			self.First = 0
			return self.PersonMgr[i]
		}
	} else {
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Uid != self.CurStep {
				continue
			}

			if i+1 < len(self.PersonMgr) {
				return self.PersonMgr[i+1]
			} else {
				return self.PersonMgr[0]
			}
		}
	}

	return nil
}

func (self *Game_EatHot) OnTime() {

}

func (self *Game_EatHot) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_EatHot) OnIsBets(uid int64) bool {
	return false
}

//! 结算所有人
func (self *Game_EatHot) OnBalance() {
}
