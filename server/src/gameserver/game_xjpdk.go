package gameserver

import (
	"lib"
	"math/rand"
	"staticfunc"
	"time"
)

/*
param1
个位  黑三首出 	 	0-无 	  1-有
十位  显示剩余牌		0-无		  1-有
百位  红桃10抓鸟   	0-无      1-有
千位  必须管          0-无      1-有
万位  两人玩法        0-无     1-有
*/
/*
param2
个位  支付方式        0-房主    1-AA
*/

type Game_XJPDK struct {
	PersonMgr []*Game_XJPDK_Person         `json:"personmgr"`
	LastCard  []int                        `json:"lastcard"`
	CurStep   int64                        `json:"curstep"`
	BefStep   int64                        `json:"befstep"`
	State     int                          `json:"state"`
	Dealer    int64                        `json:"dealer"`
	Doubler   int64                        `json:"doubler"`
	Winner    int64                        `json:"winner"`
	SX        int64                        `json:"sx"`
	Must      bool                         `json:"must"`
	OtherCard []int                        `json:"othercard"`
	Record    *staticfunc.Rec_GamePDK_Info `json:"record"`

	room *Room
}

type Game_XJPDK_Person struct {
	Uid       int64 `json:"uid"`
	Ready     bool  `json:"ready"`
	Card      []int `json:"card"`
	Deal      bool  `json:"deal"`
	Doubler   bool  `json:"doubler"`
	Total     int   `json:"total"`
	CurScore  int   `json:"curscore"`
	Taxi      int   `json:"taxi"`
	WinNum    int   `json:"winnum"`
	Loser     int   `json:"loser"`
	AllBoom   int   `json:"allboom"`
	CTNum     int   `json:"ctnum"`
	CT        bool  `json:"ct"`
	Rate      int   `json:"rate"`
	Dan       int   `json:"dan"`
	Boom      int   `json:"boom"`      //! 炸弹数
	BoomScore int   `json:"boomscore"` //! 炸弹分
	MaxScore  int   `json:"maxscore"`
	CardNum   int   `json:"cardnum"`
}

type Msg_GameXJPDK_Step struct {
	Uid     int64 `json:"uid"`
	Cards   []int `json:"cards"`
	CardNum int   `json:"cardnum"`
	CurStep int64 `json:"curstep"`
}

func (self *Game_XJPDK_Person) Init() {
	self.Card = make([]int, 0)
	self.CurScore = 0
	self.Taxi = 0
	self.Rate = 1
	self.BoomScore = 0
	self.Boom = 0
	self.Deal = false
	self.CT = false
}

type Msg_GameXJPDK_Info struct {
	Doubler  int64 `json:"doubler"`
	Begin    bool  `json:"begin"`
	CurStep  int64 `json:"curstep"`
	BefStep  int64 `json:"befstep"`
	LastCard []int `json:"lastcard"`
	Dealer   int64 `json:"dealer"`
	State    int   `json:"state"`

	Must bool                 `json:"must"`
	Info []Son_GameXJPDK_Info `json:"info"`
}

type Son_GameXJPDK_Info struct {
	Uid      int64 `json:"uid"`
	Ready    bool  `json:"ready"`
	Card     []int `json:"card"`
	Deal     bool  `json:"deal"`
	Doubler  bool  `json:"doubler"`
	Total    int   `json:"total"`
	CurScore int   `json:"curscore"`
	CardNum  int   `json:"cardnum"`
}

type Msg_GameXJPDK_End struct {
	Winner  int64               `json:"winner"`
	Doubler int64               `json:"doubler"`
	Info    []Son_GameXJPDK_End `json:"info"`
}
type Son_GameXJPDK_End struct {
	Uid       int64 `json:"uid"`
	Card      []int `json:"card"`
	CardNum   int   `json:"cardnum"`
	CT        bool  `json:"ct"`
	Score     int   `json:"score"`
	Boom      int   `json:"boom"`
	CurScore  int   `json:"curscore"`
	BoomScore int   `json:"boomscore"`
}

type Msg_GameXJPDK_Bye struct {
	Info []Son_GameXJPDK_Bye `json:"info"`
}
type Son_GameXJPDK_Bye struct {
	Uid      int64 `json:"uid"`
	Score    int   `json:"score"`
	WinNum   int   `json:"winnum"`
	AllBoom  int   `json:"allboom"`
	Loser    int   `json:"loser"`
	MaxScore int   `json:"maxscore"`
}

func NewGame_XJPDK() *Game_XJPDK {
	game := new(Game_XJPDK)
	game.PersonMgr = make([]*Game_XJPDK_Person, 0)
	return game
}

func (self *Game_XJPDK) OnIsBets(uid int64) bool {
	return false
}

func (self *Game_XJPDK) GetPerson(uid int64) *Game_XJPDK_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}
	return nil
}

func (self *Game_XJPDK) GetNextUid() int64 {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid != self.CurStep {
			continue
		}
		if i+1 < len(self.PersonMgr) {
			return self.PersonMgr[i+1].Uid
		} else {
			return self.PersonMgr[0].Uid
		}
	}
	return 0
}

func (self *Game_XJPDK) getInfo(uid int64) *Msg_GameXJPDK_Info {
	var msg Msg_GameXJPDK_Info
	msg.Begin = self.room.Begin
	msg.BefStep = self.BefStep
	msg.CurStep = self.CurStep
	msg.Dealer = self.Dealer
	msg.LastCard = self.LastCard
	msg.State = self.State
	msg.Must = self.Must
	msg.Doubler = self.Doubler
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameXJPDK_Info
		son.Uid = self.PersonMgr[i].Uid
		son.CurScore = self.PersonMgr[i].CurScore
		son.Total = self.PersonMgr[i].Total
		son.Deal = self.PersonMgr[i].Deal
		son.Doubler = self.PersonMgr[i].Doubler
		son.Ready = self.PersonMgr[i].Ready
		if self.room.Param1/10%10 == 1 {
			son.CardNum = self.PersonMgr[i].CardNum
		} else {
			son.CardNum = -1
		}
		if self.PersonMgr[i].Uid == uid || self.State != 1 {
			son.Card = self.PersonMgr[i].Card
		} else {
			son.Card = make([]int, len(self.PersonMgr[i].Card))
		}
		msg.Info = append(msg.Info, son)
	}
	return &msg
}

func (self *Game_XJPDK) Ready(uid int64) {
	if self.room.IsBye() {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "----------Ready() room is bye")
		return
	}

	if self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "----------Ready() room is begin")
		return
	}

	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Ready {
				return
			} else {
				self.PersonMgr[i].Ready = true
				num++
			}
		} else if self.PersonMgr[i].Ready {
			num++
		}
	}
	if self.room.Param1/10000%10 == 0 {
		if num == 3 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
			self.OnBegin()
			return
		}
	} else {
		if num == 2 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
			self.OnBegin()
			return
		}
	}

	self.room.flush()

	var msg staticfunc.Msg_Uid
	msg.Uid = uid
	self.room.broadCastMsg("gameXJPDKready", &msg)

}

func (self *Game_XJPDK) GameStep(uid int64, card []int, abscard []int) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}
	if self.CurStep != uid {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "---------GameStep() 该局不归你操作")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "未找到person")
		return
	}

	if len(card) != 0 {
		for _, value := range card {
			find := false
			for i := 0; i < len(person.Card); i++ {
				if person.Card[i] == value {
					find = true
					break
				}
			}
			if !find {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, card, ":出牌找不到")
				return
			}
		}
		if self.room.Step == 1 && self.room.Param1%10 == 1 && person.Taxi == 0 && self.Dealer == uid && self.room.Param2/10000 != 1 {
			find := false
			for i := 0; i < len(card); i++ {
				if card[i] == 34 {
					find = true
					break
				}
			}
			if !find {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "黑三首出")
				return
			}
		}
		tmp := 0
		if len(card) == len(person.Card) {
			tmp = IsOkByCardsZYPDK(card, true)
		} else {
			tmp = IsOkByCardsZYPDK(card, false)
		}
		if tmp == 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌错误")
			return
		}
		if len(self.LastCard) != 0 && self.BefStep != uid {
			_tmp := IsOkByCardsZYPDK(self.LastCard, false)
			if tmp%100 == TYPE_CARD_ZHA {
				if _tmp%100 == TYPE_CARD_ZHA {
					if CardCompare(tmp/10, _tmp/10) <= 0 {
						lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌更小")
						return
					}
				}
				person.Boom++
				person.AllBoom++
			} else {
				if len(card) != len(self.LastCard) {
					lib.GetLogMgr().Output(lib.LOG_DEBUG, card, ":出牌张数不匹配")
					return
				}
				if tmp%100 != _tmp%100 {
					lib.GetLogMgr().Output(lib.LOG_DEBUG, card, ":出牌不匹配")
					return
				}
				if CardCompare(tmp/10, _tmp/10) <= 0 {
					lib.GetLogMgr().Output(lib.LOG_DEBUG, card, ":出牌更小")
					return
				}
			}
		} else {
			if tmp%100 == TYPE_CARD_ZHA {
				person.Boom++
				person.AllBoom++
			}
			_tmp := IsOkByCardsZYPDK(self.LastCard, false)
			if self.BefStep == uid && _tmp%100 == TYPE_CARD_ZHA {
				score := 0
				for i := 0; i < len(self.PersonMgr); i++ {
					if self.PersonMgr[i].Uid == uid {
						continue
					}
					self.PersonMgr[i].BoomScore -= 20
					score += 20
				}
				person.BoomScore += score
			}
		}

		self.LastCard = card
		self.BefStep = uid
		person.Taxi++

		for _, value := range card {
			for i := 0; i < len(person.Card); i++ {
				if person.Card[i] == value {
					copy(person.Card[i:], person.Card[i+1:])
					person.Card = person.Card[:len(person.Card)-1]
					break
				}
			}
		}
		if len(person.Card) == 1 {
			person.Dan++
		}
	} else {
		if len(self.LastCard) == 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "第一局不可以过")
			return
		}
	}
	self.CurStep = self.GetNextUid()
	person.CardNum = len(person.Card)

	if self.Record != nil {
		self.Record.Step = append(self.Record.Step, staticfunc.Son_Rec_GamePDK_Step{uid, card, abscard})
	}

	var msg Msg_GameXJPDK_Step
	msg.Uid = uid
	msg.Cards = card
	if self.room.Param1/10%10 == 1 {
		msg.CardNum = person.CardNum
	} else {
		msg.CardNum = -1
	}

	msg.CurStep = self.CurStep
	self.room.broadCastMsg("gameXJPDKstep", &msg)

	if len(abscard) != 0 {
		self.SX = int64(abscard[0])
	}

	if person.CardNum == 0 {
		_tmp := IsOkByCardsZYPDK(self.LastCard, false)
		if _tmp%100 == TYPE_CARD_ZHA {
			score := 0
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Uid == uid {
					continue
				}
				self.PersonMgr[i].BoomScore -= 20
				score += 20
			}
			person.BoomScore += score
		}
		self.Winner = uid
		self.OnEnd()
		return
	}

	self.room.flush()
}

func (self *Game_XJPDK) OnBegin() {
	if self.room.IsBye() {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "----------onbegin()  room is bye")
		return
	}
	self.room.SetBegin(true)
	self.State = 1
	self.LastCard = make([]int, 0)
	self.BefStep = 0
	self.Dealer = 0
	self.SX = 0
	self.Record = new(staticfunc.Rec_GamePDK_Info)

	cardmgr := NewCard_Run48()

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Init()

		self.PersonMgr[i].Card = cardmgr.Deal(16)
		var rc_person staticfunc.Son_Rec_GamePDK_Person
		rc_person.Uid = self.PersonMgr[i].Uid
		rc_person.Name = self.room.GetName(rc_person.Uid)
		rc_person.Head = self.room.GetHead(rc_person.Uid)
		lib.HF_DeepCopy(&rc_person.Card, &self.PersonMgr[i].Card)
		self.Record.Person = append(self.Record.Person, rc_person)
	}
	if self.Winner == 0 {
		if len(self.PersonMgr) == 2 {
			rd := rand.New(rand.NewSource(time.Now().UnixNano()))
			self.Winner = self.PersonMgr[rd.Intn(2)].Uid
			self.Dealer = self.Winner
		} else {
			for i := 0; i < len(self.PersonMgr); i++ {
				for j := 0; j < len(self.PersonMgr[i].Card); j++ {
					if self.PersonMgr[i].Card[j] == 34 {
						self.Dealer = self.PersonMgr[i].Uid
						break
					}
				}
				if self.Dealer != 0 {
					break
				}
			}
		}
	} else {
		self.Dealer = self.Winner
	}
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == self.Dealer {
			self.PersonMgr[i].Deal = true
		}
	}
	self.CurStep = self.Dealer
	self.Doubler = 0
	if self.room.Param1/100%10 == 1 {
		for i := 0; i < len(self.PersonMgr); i++ {
			for j := 0; j < len(self.PersonMgr[i].Card); j++ {
				if self.PersonMgr[i].Card[j] == 103 {
					self.Doubler = self.PersonMgr[i].Uid
					self.PersonMgr[i].Rate = 2
					self.PersonMgr[i].Doubler = true
					break
				}
			}
			if self.Doubler != 0 {
				break
			}
		}
	}
	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gameXJPDKbegin", self.getInfo(person.Uid))
	}

}

func (self *Game_XJPDK) OnEnd() {
	self.room.SetBegin(false)
	self.State = 2
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Ready = false
	}

	if self.Winner == self.Doubler {
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Uid == self.Winner {
				continue
			}
			self.PersonMgr[i].Rate *= 2
		}
	}

	winner := self.GetPerson(self.Winner)
	if winner == nil {
		return
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------- uid : ", self.PersonMgr[i].Uid, " boomscore : ", self.PersonMgr[i].BoomScore)
	}

	score := 0
	if self.SX == 0 { //! 不包赔
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Uid == self.Winner {
				continue
			}
			if self.PersonMgr[i].Taxi == 0 { //! 春天
				self.PersonMgr[i].Rate *= 2
				self.PersonMgr[i].CT = true
			}
			if len(self.PersonMgr[i].Card) == 1 {
				self.PersonMgr[i].CurScore += self.PersonMgr[i].BoomScore
				self.PersonMgr[i].Total += self.PersonMgr[i].CurScore
				continue
			}
			self.PersonMgr[i].CurScore += (self.PersonMgr[i].BoomScore - len(self.PersonMgr[i].Card)*self.PersonMgr[i].Rate)
			self.PersonMgr[i].Total += self.PersonMgr[i].CurScore
			score += len(self.PersonMgr[i].Card) * self.PersonMgr[i].Rate
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------- uid : ", self.PersonMgr[i].Uid, " CurScore : ", self.PersonMgr[i].CurScore, " Total : ", self.PersonMgr[i].Total)
		}
	} else { //! 包赔
		sx := self.GetPerson(self.SX)
		if sx == nil {
			return
		}
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Uid == self.Winner {
				continue
			}
			if self.PersonMgr[i].Taxi == 0 { //! 春天
				self.PersonMgr[i].Rate *= 2
				self.PersonMgr[i].CT = true
			}
			if len(self.PersonMgr[i].Card) == 1 {
				score -= self.PersonMgr[i].BoomScore
				continue
			}
			if self.PersonMgr[i].BoomScore < 0 {
				score -= self.PersonMgr[i].BoomScore
			}

			score += len(self.PersonMgr[i].Card) * self.PersonMgr[i].Rate
		}
		sx.CurScore -= score
		sx.Total -= score
	}
	winner.CurScore = score + winner.BoomScore
	winner.Total += winner.CurScore
	winner.WinNum++
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "---------- winner.uid : ", winner.Uid, " winner.cursocre : ", winner.CurScore, " winner.total : ", winner.Total, " winnner.boomscore : ", winner.BoomScore)

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].MaxScore < self.PersonMgr[i].CurScore {
			self.PersonMgr[i].MaxScore = self.PersonMgr[i].CurScore
		}
		if len(self.PersonMgr[i].Card) > 0 {
			self.PersonMgr[i].Loser++
		}

	}

	var msg Msg_GameXJPDK_End
	msg.Doubler = self.Doubler
	msg.Winner = self.Winner
	msg.Info = make([]Son_GameXJPDK_End, 0)

	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameXJPDK_End
		son.Uid = self.PersonMgr[i].Uid
		son.Card = self.PersonMgr[i].Card
		son.CardNum = len(self.PersonMgr[i].Card)
		son.CT = self.PersonMgr[i].CT
		son.CurScore = self.PersonMgr[i].CurScore
		son.Boom = self.PersonMgr[i].Boom
		son.Score = self.PersonMgr[i].Total
		son.BoomScore = self.PersonMgr[i].BoomScore
		msg.Info = append(msg.Info, son)
		if self.Record != nil {
			for j := 0; j < len(self.Record.Person); j++ {
				if self.Record.Person[j].Uid == self.PersonMgr[i].Uid {
					self.Record.Person[j].Total = self.PersonMgr[i].Total
					self.Record.Person[j].Score = self.PersonMgr[i].CurScore
					break
				}
			}
		}
	}

	if self.Record != nil {
		self.Record.Roomid = self.room.Id*100 + self.room.Step
		self.Record.MaxStep = self.room.MaxStep
		self.Record.Type = self.room.Type
		self.Record.Time = time.Now().Unix()
		self.Record.Razz = 0
		self.room.AddRecord(lib.HF_JtoA(self.Record))
	}

	self.room.broadCastMsg("gameXJPDKend", &msg)

	if self.room.IsBye() {
		self.OnBye()
		self.room.Bye()
		return
	}
	self.room.flush()
}

func (self *Game_XJPDK) OnInit(room *Room) {
	self.room = room
}

func (self *Game_XJPDK) OnRobot(robot *lib.Robot) {

}

func (self *Game_XJPDK) OnSendInfo(person *Person) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			person.SendMsg("gameXJPDKinfo", self.getInfo(person.Uid))
			return
		}
	}
	_person := new(Game_XJPDK_Person)
	if _person == nil {
		return
	}
	_person.Uid = person.Uid
	_person.Ready = true
	self.PersonMgr = append(self.PersonMgr, _person)
	person.SendMsg("gameXJPDKinfo", self.getInfo(person.Uid))

	if self.room.Param1/10000%10 == 0 {
		if len(self.PersonMgr) == 3 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
			self.OnBegin()
			return
		}
	} else {
		if len(self.PersonMgr) == 2 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
			self.OnBegin()
			return
		}
	}

}

func (self *Game_XJPDK) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "gameready":
		self.Ready(msg.Uid)
	case "gamesteps":
		self.GameStep(msg.Uid, msg.V.(*Msg_GameSteps).Cards, msg.V.(*Msg_GameSteps).AbsCards)
	}
}

func (self *Game_XJPDK) OnBye() {
	var msg Msg_GameXJPDK_Bye
	msg.Info = make([]Son_GameXJPDK_Bye, 0)
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameXJPDK_Bye
		son.Uid = self.PersonMgr[i].Uid
		son.Score = self.PersonMgr[i].Total
		son.WinNum = self.PersonMgr[i].WinNum
		son.AllBoom = self.PersonMgr[i].AllBoom
		son.Loser = self.PersonMgr[i].Loser
		son.MaxScore = self.PersonMgr[i].MaxScore
		msg.Info = append(msg.Info, son)
	}
	self.room.broadCastMsg("gameXJPDKbye", &msg)
}

func (self *Game_XJPDK) OnExit(uid int64) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if uid == self.PersonMgr[i].Uid {
			copy(self.PersonMgr[i:], self.PersonMgr[i+1:])
			self.PersonMgr = self.PersonMgr[:len(self.PersonMgr)-1]
			return
		}
	}
}

func (self *Game_XJPDK) OnBalance() {

} //! 结算

func (self *Game_XJPDK) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_XJPDK) OnTime() {

}
