package gameserver

import (
	"lib"
	"staticfunc"
	"time"
)

//! param1%10 = 0看牌抢庄 1世界大战 2房主当庄 3轮流当庄
//! param1/10%10 = 若是看牌抢庄则为0扣一张 1扣两张；若是其他模式为1,2,3,5倍
//! param1/100%10 = 0没有癞子 1有癞子
//! param2 人数2-5

//! 记录结构
type Rec_GameNiuNiuJX struct {
	Info    []Son_Rec_GameNiuNiuJX `json:"info"`
	Roomid  int                    `json:"roomid"`
	MaxStep int                    `json:"maxstep"`
	Param1  int                    `json:"param1"`
	Param2  int                    `json:"param2"`
	Time    int64                  `json:"time"`
}
type Son_Rec_GameNiuNiuJX struct {
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Card    []int  `json:"card"`
	Bets    int    `json:"bets"`
	Dealer  bool   `json:"dealer"`
	Score   int    `json:"score"`
	Type    int    `json:"type"`
	RobDeal int    `json:"robdeal"`
}

//! 牛牛亮牌
type Msg_GameNiuNiuJX_View struct {
	Uid  int64 `json:"uid"`
	Card []int `json:"card"`
	CT   int   `json:"ct"`
}

//!
type Msg_GameNiuNiuJX_Info struct {
	Begin bool                    `json:"begin"` //! 是否开始
	Info  []Son_GameNiuNiuJX_Info `json:"info"`
	Time  int64                   `json:"time"` //! 倒计时
}

type Son_GameNiuNiuJX_Info struct {
	Uid     int64 `json:"uid"`
	Card    []int `json:"card"`
	Ready   bool  `json:"ready"`
	Bets    int   `json:"bets"`
	Dealer  bool  `json:"dealer"`
	Score   int   `json:"score"`
	Total   int   `json:"total"`
	View    bool  `json:"view"`
	CT      int   `json:"ct"`
	RobDeal int   `json:"robdeal"`
}

//! 结算
type Msg_GameNiuNiuJX_End struct {
	Info []Son_GameNiuNiuJX_Info `json:"info"`
}

//! 得到最后一张牌
type Msg_GameNiuNiuJX_Card struct {
	Card []int `json:"card"`
}

///////////////////////////////////////////////////////
type Game_NiuNiuJX_Person struct {
	Uid      int64 `json:"uid"`
	Card     []int `json:"card"`     //! 手牌
	Win      int   `json:"win"`      //! 胜利次数
	Niu      int   `json:"niu"`      //! 牛牛次数
	Kill     int   `json:"kill"`     //! 通杀次数
	Dead     int   `json:"dead"`     //! 通赔次数
	Deal     int   `json:"deal"`     //! 坐庄次数
	Ready    bool  `json:"ready"`    //! 是否准备
	Score    int   `json:"score"`    //! 积分
	Dealer   bool  `json:"dealer"`   //! 是否庄家
	RobDeal  int   `json:"robdeal"`  //! 是否抢庄
	CurScore int   `json:"curscore"` //! 当前局的分数
	View     bool  `json:"view"`     //! 是否亮牌
	CT       int   `json:"ct"`       //! 当前牌型
	Bets     int   `json:"bets"`     //! 下注
	CS       int   `json:"cs"`       //! 当前局最大牌
}

func (self *Game_NiuNiuJX_Person) Init() {
	self.CT = 0
	self.CS = 0
	self.CurScore = 0
	self.Dealer = false
	self.Bets = 0
	self.RobDeal = -1
	self.View = false
}

type Game_NiuNiuJX struct {
	PersonMgr []*Game_NiuNiuJX_Person `json:"personmgr"`
	Card      *CardMgr                `json:"card"`
	Time      int64                   `json:"time"` //! 强制操作时间
	Mode      int                     `json:"mode"` //! 强制操作模式

	room *Room
}

func NewGame_NiuNiuJX() *Game_NiuNiuJX {
	game := new(Game_NiuNiuJX)
	game.PersonMgr = make([]*Game_NiuNiuJX_Person, 0)

	return game
}

func (self *Game_NiuNiuJX) GetBets() int {
	if self.room.Param1/10%10 == 0 {
		return 1
	} else if self.room.Param1/10%10 == 1 {
		return 3
	} else if self.room.Param1/10%10 == 2 {
		return 5
	}
	return 0
}

func (self *Game_NiuNiuJX) OnInit(room *Room) {
	self.room = room

	if self.room.Param2 < 2 {
		self.room.Param2 = 2
	} else if self.room.Param2 > 5 {
		self.room.Param2 = 5
	}
}

func (self *Game_NiuNiuJX) OnRobot(robot *lib.Robot) {

}

func (self *Game_NiuNiuJX) OnSendInfo(person *Person) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			person.SendMsg("gameniuniujxinfo", self.getInfo(person.Uid))
			return
		}
	}

	_person := new(Game_NiuNiuJX_Person)
	_person.Init()
	_person.Uid = person.Uid
	_person.Ready = true
	self.PersonMgr = append(self.PersonMgr, _person)

	minnum := self.room.Param2
	if len(self.PersonMgr) >= minnum {
		self.OnBegin()
		return
	}

	person.SendMsg("gameniuniujxinfo", self.getInfo(person.Uid))
}

func (self *Game_NiuNiuJX) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "gameready": //! 游戏准备
		self.GameReady(msg.Uid)
	case "gamebets": //! 下注
		self.GameBets(msg.Uid, msg.V.(*Msg_GameBets).Bets)
	case "gameview": //! 亮牌
		self.GameView(msg.Uid, true)
	case "gamedealer": //! 抢庄
		self.GameDeal(msg.Uid, msg.V.(*Msg_GameDealer).Score)
	}
}

func (self *Game_NiuNiuJX) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)

	DealerPos := -1
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Dealer {
			DealerPos = i
			break
		}
	}

	//! 初始化游戏人
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Init()
		if self.room.Param1%10 != 0 {
			self.PersonMgr[i].Bets = self.GetBets()
		}
	}

	//! 确定庄家
	if self.room.Param1%10 == 2 { //! 房主当庄
		self.PersonMgr[0].Dealer = true
		self.PersonMgr[0].Deal++
		if self.GetBets() != 0 {
			self.PersonMgr[0].RobDeal = 1
		}
	} else if self.room.Param1%10 == 3 { //! 轮流当庄
		if DealerPos+1 >= len(self.PersonMgr) {
			DealerPos = -1
		}
		self.PersonMgr[DealerPos+1].Dealer = true
		self.PersonMgr[DealerPos+1].Deal++
		if self.GetBets() != 0 {
			self.PersonMgr[DealerPos+1].RobDeal = 1
		}
	}

	//! 发牌
	self.Card = NewCard_NiuNiu(self.room.Param1/100%10 == 1)
	if self.room.Param1%10 == 0 { //! 看牌抢庄
		for i := 0; i < len(self.PersonMgr); i++ {
			self.PersonMgr[i].Card = self.Card.Deal(lib.HF_MaxInt(3, 5-self.room.Param1/10%10-1))
		}
	} else {
		if self.GetBets() == 0 {
			for i := 0; i < len(self.PersonMgr); i++ {
				self.PersonMgr[i].Card = self.Card.Deal(4)
			}
		} else {
			for i := 0; i < len(self.PersonMgr); i++ {
				self.PersonMgr[i].Card = self.Card.Deal(5)
			}
		}
	}

	if self.room.Param1%10 == 0 || self.GetBets() == 0 { //! 要抢庄
		if self.room.Param1%10 == 1 {
			self.Mode = 2
			self.Time = time.Now().Unix() + 11
		} else {
			self.Mode = 1
			self.Time = time.Now().Unix() + 11
		}
	} else { //! 亮牌
		self.Mode = 3
		self.Time = time.Now().Unix() + 5
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gameniuniujxbegin", self.getInfo(person.Uid))
	}

	self.room.flush()
}

//! 抢庄
func (self *Game_NiuNiuJX) GameDeal(uid int64, score int) {
	if !self.room.Begin { //! 未开始不能抢庄
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if self.room.Param1%10 != 0 { //! 不是看牌抢庄模式不能抢庄
		if !person.Dealer {
			return
		}
		if score <= 0 {
			score = 1
		}
	}

	robnum := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].RobDeal >= 0 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能重复抢庄")
				return
			} else {
				self.PersonMgr[i].RobDeal = score
			}
		}
		if self.PersonMgr[i].RobDeal >= 0 {
			robnum++
		}
	}

	//! 广播
	var msg Msg_GameDealer
	msg.Uid = uid
	msg.Score = score
	self.room.broadCastMsg("gamedeal", &msg)

	if self.room.Param1%10 != 0 {
		var msg staticfunc.Msg_Uid
		msg.Uid = person.Uid
		self.room.broadCastMsg("gamedealer", &msg)

		//! 下注
		self.Mode = 2
		self.Time = time.Now().Unix() + 8
	} else {
		if robnum == len(self.PersonMgr) { //! 全部发表了意见
			deal := make([]*Game_NiuNiuJX_Person, 0)
			for i := 0; i < len(self.PersonMgr); i++ {
				if len(deal) == 0 {
					deal = append(deal, self.PersonMgr[i])
				} else {
					if self.PersonMgr[i].RobDeal > deal[0].RobDeal {
						deal = make([]*Game_NiuNiuJX_Person, 0)
						deal = append(deal, self.PersonMgr[i])
					} else if self.PersonMgr[i].RobDeal == deal[0].RobDeal {
						deal = append(deal, self.PersonMgr[i])
					}
				}
			}

			dealer := deal[lib.HF_GetRandom(len(deal))]
			dealer.Dealer = true
			dealer.Deal++
			if dealer.RobDeal <= 0 {
				dealer.RobDeal = 1
			}

			var msg staticfunc.Msg_Uid
			msg.Uid = dealer.Uid
			self.room.broadCastMsg("gamedealer", &msg)

			//! 下注
			self.Mode = 2
			self.Time = time.Now().Unix() + 8
		}
	}

	self.room.flush()
}

//! 亮牌
func (self *Game_NiuNiuJX) GameView(uid int64, send bool) {
	if !self.room.Begin {
		return
	}

	if self.room.Param1%10 != 1 { //! 不是世界大战判断庄家
		find := false
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Dealer {
				find = true
				break
			}
		}
		if !find {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "没有庄家无法亮牌")
			return
		}
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if len(person.Card) < 5 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "还不能亮牌")
		return
	}

	if person.View {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经亮牌了")
		return
	}

	person.CT, person.CS = GetNiuNiuByRazz(person.Card)

	if self.room.Param1%10 == 0 && send {
		var msg Msg_GameNiuNiuJX_View
		msg.Uid = uid
		msg.Card = person.Card
		msg.CT = person.CT
		self.room.broadCastMsg("gameview", &msg)
	}

	person.View = true
	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].View {
			num++
		}
	}

	if num >= len(self.PersonMgr) {
		self.OnEnd()
		return
	}

	self.room.flush()
}

//! 准备
func (self *Game_NiuNiuJX) GameReady(uid int64) {
	if self.room.IsBye() {
		return
	}

	if self.room.Begin {
		return
	}

	if self.room.Step == 0 { //! 一局之后才有这个消息
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

	if num == len(self.room.Uid) {
		self.OnBegin()
		return
	}

	var msg staticfunc.Msg_Uid
	msg.Uid = uid
	self.room.broadCastMsg("gameready", &msg)

	self.room.flush()
}

//! 下注
func (self *Game_NiuNiuJX) GameBets(uid int64, bets int) {
	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	if bets <= 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "下注无效")
		return
	}

	//if self.room.Param1%10 != 0 { //! 不是看牌抢庄模式
	//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "该模式无法下注")
	//	return
	//}

	var dealer *Game_NiuNiuJX_Person = nil
	if self.room.Param1%10 != 1 { //! 不是世界大战
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Dealer {
				dealer = self.PersonMgr[i]
				break
			}
		}
		if dealer == nil {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "没有庄家无法下注")
			return
		}
	}

	if bets >= 5 {
		person := GetPersonMgr().GetPerson(uid)
		if person == nil {
			bets = 4
		} else if person.Card+person.Gold < 1 {
			bets = 4
		}
	}

	betnum := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Dealer { //! 是庄家
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "庄家不用下注")
				return
			}

			if self.PersonMgr[i].Bets > 0 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能重复下注")
				return
			}

			self.PersonMgr[i].Bets = bets

			//! 广播这个人下注
			var msg Msg_GameBets
			msg.Uid = self.PersonMgr[i].Uid
			msg.Bets = self.PersonMgr[i].Bets
			self.room.broadCastMsg("gamebets", &msg)
		}

		if self.PersonMgr[i].Bets > 0 {
			betnum++
		}
	}

	if dealer == nil {
		if betnum == len(self.PersonMgr) {
			for i := 0; i < len(self.PersonMgr); i++ {
				card := self.Card.Deal(5 - len(self.PersonMgr[i].Card))
				self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, card...)

				var msg Msg_GameNiuNiuJX_Card
				msg.Card = card
				person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
				if person != nil {
					person.SendMsg("gameniuniucard", &msg)
				}
			}
			//! 亮牌
			self.Mode = 3
			self.Time = time.Now().Unix() + 3
		}
	} else {
		if betnum == len(self.PersonMgr)-1 && dealer.RobDeal > 0 {
			for i := 0; i < len(self.PersonMgr); i++ {
				card := self.Card.Deal(5 - len(self.PersonMgr[i].Card))
				self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, card...)

				var msg Msg_GameNiuNiuJX_Card
				msg.Card = card
				person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
				if person != nil {
					person.SendMsg("gameniuniucard", &msg)
				}
			}
			//! 亮牌
			self.Mode = 3
			self.Time = time.Now().Unix() + 3
		}
	}

	self.room.flush()
}

//! 结算
func (self *Game_NiuNiuJX) OnEnd() {
	self.room.SetBegin(false)

	var dealer *Game_NiuNiuJX_Person = nil
	if self.room.Param1%10 == 1 { //! 世界大战
		for i := 0; i < len(self.PersonMgr); i++ {
			if dealer == nil {
				dealer = self.PersonMgr[i]
			} else {
				if self.PersonMgr[i].CT > dealer.CT {
					dealer = self.PersonMgr[i]
				} else if self.PersonMgr[i].CT == dealer.CT && self.PersonMgr[i].CS > dealer.CS {
					dealer = self.PersonMgr[i]
				}
			}
		}
	} else {
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Dealer {
				dealer = self.PersonMgr[i]
				break
			}
		}
	}
	if dealer.RobDeal <= 0 {
		dealer.RobDeal = 1
	}

	lst := make([]*Game_NiuNiuJX_Person, 0)
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Ready = false
		if self.PersonMgr[i].CT >= 100 {
			self.PersonMgr[i].Niu++
		}

		if self.PersonMgr[i].Uid != dealer.Uid {
			lst = append(lst, self.PersonMgr[i])
		}
	}

	win := 0
	for i := 0; i < len(lst); i++ {
		dealerwin := false
		if dealer.CT > lst[i].CT { //! 庄家赢
			dealerwin = true
		} else if dealer.CT < lst[i].CT { //! 闲家赢
			dealerwin = false
		} else {
			if dealer.CS > lst[i].CS { //! 庄家赢
				dealerwin = true
			} else { //! 闲家赢
				dealerwin = false
			}
		}

		if dealerwin { //! 庄家赢
			bs := GetNiuNiuJXBS(dealer.CT)
			score := lst[i].Bets * dealer.RobDeal * bs
			if self.room.Param1%10 == 1 {
				score *= dealer.Bets
			}
			dealer.CurScore += score
			lst[i].CurScore += -score
			win++
		} else { //! 闲家赢
			bs := GetNiuNiuJXBS(lst[i].CT)
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "闲家赢:", bs)
			score := lst[i].Bets * dealer.RobDeal * bs
			if self.room.Param1%10 == 1 {
				score *= dealer.Bets
			}
			lst[i].CurScore += score
			dealer.CurScore += -score
			lst[i].Win++
			win--
		}
		lst[i].Score += lst[i].CurScore
	}
	dealer.Score += dealer.CurScore
	if dealer.CurScore > 0 {
		dealer.Win++
	}
	if win == len(lst) {
		dealer.Kill++
	} else if -win == len(lst) {
		dealer.Dead++
	}

	//! 记录
	var record Rec_GameNiuNiuJX
	record.Time = time.Now().Unix()
	record.Roomid = self.room.Id*100 + self.room.Step
	record.MaxStep = self.room.MaxStep
	record.Param1 = self.room.Param1
	record.Param2 = self.room.Param2

	//! 发消息
	var msg Msg_GameNiuNiuJX_End
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameNiuNiuJX_Info
		son.Uid = self.PersonMgr[i].Uid
		son.Bets = self.PersonMgr[i].Bets
		son.Card = self.PersonMgr[i].Card
		son.Dealer = self.PersonMgr[i].Dealer
		son.Score = self.PersonMgr[i].CurScore
		son.Total = self.PersonMgr[i].Score
		son.CT = self.PersonMgr[i].CT
		son.View = self.PersonMgr[i].View
		msg.Info = append(msg.Info, son)

		var rec Son_Rec_GameNiuNiuJX
		rec.Uid = self.PersonMgr[i].Uid
		rec.Name = self.room.GetName(self.PersonMgr[i].Uid)
		rec.Head = self.room.GetHead(self.PersonMgr[i].Uid)
		rec.Card = self.PersonMgr[i].Card
		rec.Dealer = self.PersonMgr[i].Dealer
		rec.Bets = self.PersonMgr[i].Bets
		rec.Score = self.PersonMgr[i].CurScore
		rec.Type = self.PersonMgr[i].CT
		rec.RobDeal = self.PersonMgr[i].RobDeal
		record.Info = append(record.Info, rec)
	}
	self.room.AddRecord(lib.HF_JtoA(&record))
	self.room.broadCastMsg("gameniuniuend", &msg)

	if self.room.IsBye() {
		self.OnBye()
		self.room.Bye()
		return
	}

	self.Mode = 4
	self.Time = time.Now().Unix() + 8

	self.room.flush()
}

func (self *Game_NiuNiuJX) OnBye() {
	self.Mode = 0
	self.Time = 0

	info := make([]staticfunc.JS_CreateRoomMem, 0)
	var msg Msg_GameNiuNiu_Bye
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameNiuNiu_Bye
		son.Uid = self.PersonMgr[i].Uid
		son.Win = self.PersonMgr[i].Win
		son.Niu = self.PersonMgr[i].Niu
		son.Kill = self.PersonMgr[i].Kill
		son.Dead = self.PersonMgr[i].Dead
		son.Deal = self.PersonMgr[i].Deal
		son.Score = self.PersonMgr[i].Score
		msg.Info = append(msg.Info, son)
		info = append(info, staticfunc.JS_CreateRoomMem{son.Uid, "", "", son.Score})
	}
	self.room.broadCastMsg("gameniuniubye", &msg)

	self.room.ClubResult(info)
}

func (self *Game_NiuNiuJX) OnExit(uid int64) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			copy(self.PersonMgr[i:], self.PersonMgr[i+1:])
			self.PersonMgr = self.PersonMgr[:len(self.PersonMgr)-1]
			break
		}
	}
}

func (self *Game_NiuNiuJX) getInfo(uid int64) *Msg_GameNiuNiuJX_Info {
	var msg Msg_GameNiuNiuJX_Info
	msg.Begin = self.room.Begin
	if self.Time != 0 {
		msg.Time = self.Time - time.Now().Unix()
	}
	msg.Info = make([]Son_GameNiuNiuJX_Info, 0)
	for _, value := range self.PersonMgr {
		var son Son_GameNiuNiuJX_Info
		son.Uid = value.Uid
		son.Ready = value.Ready
		son.Bets = value.Bets
		son.Dealer = value.Dealer
		son.Total = value.Score
		son.Score = value.CurScore
		son.RobDeal = value.RobDeal
		if value.Uid == uid || value.View || !msg.Begin { //! 是自己或者亮牌了或者已经结束了
			son.Card = value.Card
			son.CT = value.CT
		} else {
			son.Card = make([]int, 0)
			for i := 0; i < len(value.Card); i++ {
				son.Card = append(son.Card, 0)
			}
			son.CT = 0
		}

		msg.Info = append(msg.Info, son)
	}
	return &msg
}

func (self *Game_NiuNiuJX) GetPerson(uid int64) *Game_NiuNiuJX_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

func (self *Game_NiuNiuJX) OnTime() {
	if self.Mode == 0 || self.Time == 0 {
		return
	}

	if time.Now().Unix() >= self.Time {
		if self.Mode == 1 { //! 抢庄
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].RobDeal < 0 {
					self.GameDeal(self.PersonMgr[i].Uid, 0)
				}
			}
		} else if self.Mode == 2 { //! 下注
			for i := 0; i < len(self.PersonMgr); i++ {
				if !self.PersonMgr[i].Dealer && self.PersonMgr[i].Bets <= 0 {
					self.GameBets(self.PersonMgr[i].Uid, 1)
				}
			}
		} else if self.Mode == 3 { //! 亮牌
			for i := 0; i < len(self.PersonMgr); i++ {
				if !self.PersonMgr[i].View {
					self.GameView(self.PersonMgr[i].Uid, false)
				}
			}
		} else if self.Mode == 4 {
			for i := 0; i < len(self.PersonMgr); i++ {
				if !self.PersonMgr[i].Ready {
					self.GameReady(self.PersonMgr[i].Uid)
				}
			}
		}
	}
}

func (self *Game_NiuNiuJX) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_NiuNiuJX) OnIsBets(uid int64) bool {
	return false
}

//! 结算所有人
func (self *Game_NiuNiuJX) OnBalance() {
}
