package gameserver

import (
	"lib"
	"math"
	"math/rand"
	"staticfunc"
	"time"
)

type Msg_GameGoldPTJ_Info struct {
	Begin bool                   `json:"begin"` //! 是否开始
	Info  []Son_GameGoldPTJ_Info `json:"info"`
	Time  int64                  `json:"time"`
	State int                    `json:"state"`
}
type Son_GameGoldPTJ_Info struct {
	Uid     int64  `json:"uid"`
	Card    []int  `json:"card"`
	Bets    int    `json:"bets"`
	Dealer  bool   `json:"dealer"`
	Score   int    `json:"score"`
	Total   int    `json:"total"`
	CT      [2]int `json:"ct"`
	View    bool   `json:"view"`
	Ready   bool   `json:"ready"`
	RobDeal int    `json:"robdeal"`
}

///////////////////////////////////////////////////////
type Game_GoldPTJ_Person struct {
	Uid      int64  `json:"uid"`
	Card     []int  `json:"card"`     //! 手牌
	Score    int    `json:"score"`    //! 积分
	Bets     int    `json:"bets"`     //! 下注
	Dealer   bool   `json:"dealer"`   //! 是否庄家
	CurScore int    `json:"curscore"` //! 当前局的分数
	View     bool   `json:"view"`     //! 是否亮牌
	CT       [2]int `json:"ct"`       //! 当前牌型
	Ready    bool   `json:"ready"`    //! 是否准备
	RobDeal  int    `json:"robdeal"`  //! 是否抢庄
	Gold     int    `json:"gold"`     //! 当前金币
}

func (self *Game_GoldPTJ_Person) Init() {
	self.Card = make([]int, 0)
	self.Bets = 0
	self.Dealer = false
	self.CurScore = 0
	self.View = false
	self.CT[0] = 0
	self.CT[1] = 0
	self.RobDeal = -1
}

//! 同步金币
func (self *Game_GoldPTJ_Person) SynchroGold(gold int) {
	self.Score += (gold - self.Gold)
	self.Gold = gold
}

type Game_GoldPTJ struct {
	PersonMgr []*Game_GoldPTJ_Person `json:"personmgr"`
	State     int                    `json:"state"` //! 0准备阶段  1等待抢庄   2等待下注   3等待亮牌
	PJ        *CardMgr               `json:"pj"`
	DF        int                    `json:"df"` //! 底分
	Time      int64                  `json:"time"`

	room *Room
}

func NewGame_GoldPTJ() *Game_GoldPTJ {
	game := new(Game_GoldPTJ)
	game.PersonMgr = make([]*Game_GoldPTJ_Person, 0)
	game.PJ = NewCard_TJ()
	//game.Card = make([]int, 0)

	return game
}

func (self *Game_GoldPTJ) OnInit(room *Room) {
	self.room = room

	self.DF = staticfunc.GetCsvMgr().GetDF(self.room.Type)
}

func (self *Game_GoldPTJ) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldPTJ) OnSendInfo(person *Person) {
	//! 观众模式游戏,观众进来只发送游戏信息
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			self.PersonMgr[i].SynchroGold(person.Gold)
			person.SendMsg("gameptjinfo", self.getInfo(person.Uid))
			return
		}
	}

	if !self.room.Begin {
		if len(self.room.Uid)+len(self.room.Viewer) == lib.HF_Atoi(self.room.csv["minnum"]) { //! 进来的人满足最小开的人数
			self.SetTime(15)
		}
	}

	person.SendMsg("gameptjinfo", self.getInfo(0))

	if !self.room.Begin {
		if self.room.Seat(person.Uid) {
			_person := new(Game_GoldPTJ_Person)
			_person.Init()
			_person.Uid = person.Uid
			_person.Score = person.Gold
			_person.Gold = person.Gold
			_person.Ready = false
			self.PersonMgr = append(self.PersonMgr, _person)
		}
	}
}

func (self *Game_GoldPTJ) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold": //! 同步金币
		person := self.GetPerson(msg.V.(*staticfunc.Msg_SynchroGold).Uid)
		if person != nil {
			person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal()
		}
		self.room.flush()
	case "gameready": //! 游戏准备
		self.GameReady(msg.Uid)
	case "gamebets": //! 下注
		self.GameBets(msg.Uid, msg.V.(*Msg_GameBets).Bets)
	case "gameptjview": //! 亮牌
		self.GameView(msg.Uid, msg.V.(*Msg_GamePTJ_View).View)
	case "gamedeal": //! 抢庄
		self.GameDeal(msg.Uid, msg.V.(*Msg_GameDeal).Ok)
	}
}

func (self *Game_GoldPTJ) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)
	self.State = 1

	bl := 50.0
	for i := 0; i < len(self.PersonMgr); i++ {
		cost := int(math.Ceil(float64(self.DF) * bl / 100.0))
		self.PersonMgr[i].Score -= cost
		GetServer().SqlAgentGoldLog(self.PersonMgr[i].Uid, cost, self.room.Type)
		GetServer().SqlAgentBillsLog(self.PersonMgr[i].Uid, cost, self.room.Type)
	}
	self.SendTotal()

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Init()
	}

	//! 发牌
	self.PJ = NewCard_TJ()

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gameptjbegin", self.getInfo(person.Uid))
	}
	self.room.broadCastMsgView("gameptjbegin", self.getInfo(0))

	self.SetTime(8)

	self.room.flush()
}

//! 抢庄
func (self *Game_GoldPTJ) GameDeal(uid int64, ok bool) {
	if !self.room.Begin { //! 未开始不能抢庄
		return
	}

	robnum := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].RobDeal >= 0 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能重复抢庄")
				return
			} else {
				if ok {
					self.PersonMgr[i].RobDeal = 1
				} else {
					self.PersonMgr[i].RobDeal = 0
				}
			}
		}
		if self.PersonMgr[i].RobDeal >= 0 {
			robnum++
		}
	}

	//! 广播
	var msg Msg_GameDeal
	msg.Uid = uid
	msg.Ok = ok
	self.room.broadCastMsg("gamedeal", &msg)

	if robnum == len(self.PersonMgr) { //! 全部发表了意见
		deal := make([]*Game_GoldPTJ_Person, 0)
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].RobDeal > 0 {
				deal = append(deal, self.PersonMgr[i])
			}
		}
		if len(deal) == 0 {
			for i := 0; i < len(self.PersonMgr); i++ {
				deal = append(deal, self.PersonMgr[i])
			}
		}

		dealer := deal[lib.HF_GetRandom(len(deal))]
		dealer.Dealer = true

		var msg staticfunc.Msg_Uid
		msg.Uid = dealer.Uid
		self.room.broadCastMsg("gamedealer", &msg)

		//! 下注
		self.State = 2
		if len(deal) > 1 {
			self.SetTime(8)
		} else {
			self.SetTime(6)
		}
	}

	self.room.flush()
}

//! 亮牌
func (self *Game_GoldPTJ) GameView(uid int64, view []int) {
	if !self.room.Begin {
		return
	}

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

	person := self.GetPerson(uid)
	if person == nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person")
		return
	}

	if len(person.Card) == 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "person无牌")
		return
	}

	if person.View {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经亮牌了")
		return
	}

	if len(view) != len(person.Card) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "亮牌不对")
		return
	}

	var _card []int
	lib.HF_DeepCopy(&_card, &person.Card)
	for i := 0; i < len(view); i++ {
		for j := 0; j < len(_card); {
			if _card[j] == view[i] {
				copy(_card[j:], _card[j+1:])
				_card = _card[:len(_card)-1]
			} else {
				j++
			}
		}
	}
	if len(_card) > 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "亮牌错误")
		return
	}

	person.View = true
	person.Card = view
	person.CT[0] = GetPTJType(view[0], view[1])
	if len(view) == 4 {
		person.CT[1] = GetPTJType(view[2], view[3])
	}

	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].View {
			num++
		}

		var msg Msg_GamePTJ_Send_View
		msg.Uid = uid
		if self.PersonMgr[i].Uid == uid || self.room.Type%10 == 1 {
			msg.CT = person.CT
		}
		self.room.SendMsg(self.PersonMgr[i].Uid, "gameview", &msg)
	}

	if num >= len(self.PersonMgr) {
		self.OnEnd()
		return
	}

	self.room.flush()
}

//! 准备
func (self *Game_GoldPTJ) GameReady(uid int64) {
	if self.room.IsBye() {
		return
	}

	if self.room.Begin { //! 已经开始了不允许准备
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经开始了，不能准备")
		return
	}

	person := GetPersonMgr().GetPerson(uid)
	if person == nil {
		return
	}
	if person.black {
		self.room.KickPerson(uid, 95)
		return
	}

	find := false
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Ready {
				return
			}

			if self.PersonMgr[i].Score < staticfunc.GetCsvMgr().GetZR(self.room.Type) { //! 携带的金币不足，踢出去
				self.room.KickPerson(uid, 99)
				return
			}
			self.PersonMgr[i].Ready = true
			find = true
			break
		}
	}

	if !find { //! 坐下
		if !self.room.Seat(uid) {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "无法坐下")
			return
		}

		lib.GetLogMgr().Output(lib.LOG_DEBUG, "坐下后:", self.room.Viewer)

		person := GetPersonMgr().GetPerson(uid)
		if person == nil {
			return
		}

		_person := new(Game_GoldPTJ_Person)
		_person.Init()
		_person.Uid = uid
		_person.Score = person.Gold
		_person.Gold = person.Gold
		_person.Ready = true
		self.PersonMgr = append(self.PersonMgr, _person)
	}

	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Ready {
			num++
		}
	}

	if num == len(self.room.Uid)+len(self.room.Viewer) && num >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 准备的人数达到游戏最小人数
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
		self.OnBegin()
		return
	}

	var msg staticfunc.Msg_Uid
	msg.Uid = uid
	self.room.broadCastMsg("gameready", &msg)

	if num == lib.HF_Atoi(self.room.csv["minnum"]) {
		self.SetTime(10)
	}

	self.room.flush()
}

//! 下注
func (self *Game_GoldPTJ) GameBets(uid int64, bets int) {
	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	if bets <= 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "下注无效")
		return
	}

	find := false
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Dealer {
			find = true
			break
		}
	}
	if !find {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "没有庄家无法下注")
		return
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

	if betnum == len(self.PersonMgr)-1 {
		self.GameCard()
	}

	self.room.flush()
}

//! 发牌
func (self *Game_GoldPTJ) GameCard() {
	if self.PJ.random == nil {
		self.PJ.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	var sz [2]int
	sz[0] = self.PJ.random.Intn(6) + 1
	sz[1] = self.PJ.random.Intn(6) + 1

	cardnum := 0
	if self.room.Type%10 == 0 {
		cardnum = 4
		for i := 0; i < len(self.PersonMgr); i++ {
			self.PersonMgr[i].Card = self.PJ.Deal(cardnum)
		}
	} else {
		cardnum = 2
		for i := 0; i < len(self.PersonMgr); i++ {
			self.PersonMgr[i].Card = self.PJ.Deal(cardnum)
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		var msg Msg_GamePTJ_Card
		msg.SZ = sz
		msg.Card = self.PersonMgr[i].Card
		self.room.SendMsg(self.PersonMgr[i].Uid, "gameptjcard", &msg)
	}

	var msg Msg_GamePTJ_Card
	msg.SZ = sz
	msg.Card = make([]int, cardnum)
	self.room.broadCastMsgView("gameptjcard", &msg)

	//! 亮牌
	self.State = 3
	if self.room.Type%10 == 0 {
		self.SetTime(19)
	} else {
		self.SetTime(14)
	}
}

//! 结算
func (self *Game_GoldPTJ) OnEnd() {
	self.room.SetBegin(false)
	self.State = 0
	self.Time = 0

	var dealer *Game_GoldPTJ_Person
	lst := make([]*Game_GoldPTJ_Person, 0)
	for _, value := range self.PersonMgr {
		value.Ready = false
		if value.Dealer {
			dealer = value
		} else {
			lst = append(lst, value)
		}
	}

	win := 0
	for i := 0; i < len(lst); i++ {
		bet1 := lst[i].Bets % 100
		bet2 := lst[i].Bets / 100 % 100
		bet3 := lst[i].Bets / 10000

		//! 计算第一道
		dealerscore := 0
		xianscore := 0
		dealerwin := self.GetWin(dealer, lst[i])
		if dealerwin > 0 {
			//dealer.CurScore += bet1 * self.DF
			//lst[i].CurScore -= bet1 * self.DF
			dealerscore += bet1 * self.DF
			xianscore -= bet1 * self.DF
			win++
		} else if dealerwin < 0 {
			//dealer.CurScore -= bet1 * self.DF
			//lst[i].CurScore += bet1 * self.DF
			dealerscore -= bet1 * self.DF
			xianscore += bet1 * self.DF
			win--
		}

		//! 计算第二道
		if dealerwin > 0 {
			if self.GetPTJValue(dealer.CT[0]) >= 83 && (self.room.Type%10 == 1 || self.GetPTJValue(dealer.CT[1]) >= 83) {
				//dealer.CurScore += bet2 * self.DF
				//lst[i].CurScore -= bet2 * self.DF
				dealerscore += bet2 * self.DF
				xianscore -= bet2 * self.DF
			}
		} else if dealerwin < 0 {
			if self.GetPTJValue(lst[i].CT[0]) >= 83 && (self.room.Type%10 == 1 || self.GetPTJValue(lst[i].CT[1]) >= 83) {
				//dealer.CurScore -= bet2 * self.DF
				//lst[i].CurScore += bet2 * self.DF
				dealerscore -= bet2 * self.DF
				xianscore += bet2 * self.DF
			}
		}

		//! 计算第三道
		if dealerwin > 0 {
			if self.GetPTJValue(dealer.CT[0]) >= 94 && (self.room.Type%10 == 1 || self.GetPTJValue(dealer.CT[1]) >= 94) {
				//dealer.CurScore += bet3 * self.DF
				//lst[i].CurScore -= bet3 * self.DF
				dealerscore += bet3 * self.DF
				xianscore -= bet3 * self.DF
			}
		} else if dealerwin < 0 {
			if self.GetPTJValue(lst[i].CT[0]) >= 94 && (self.room.Type%10 == 1 || self.GetPTJValue(lst[i].CT[1]) >= 94) {
				//dealer.CurScore -= bet3 * self.DF
				//lst[i].CurScore += bet3 * self.DF
				dealerscore -= bet3 * self.DF
				xianscore += bet3 * self.DF
			}
		}

		if dealerscore > 0 {
			if lst[i].CurScore+lst[i].Score >= dealerscore {
				lst[i].CurScore -= dealerscore
				dealer.CurScore += dealerscore
			} else {
				abs := lst[i].CurScore + lst[i].Score
				lst[i].CurScore -= abs
				dealer.CurScore += abs
			}
		} else if xianscore > 0 {
			if dealer.CurScore+dealer.Score >= xianscore {
				dealer.CurScore -= xianscore
				lst[i].CurScore += xianscore
			} else {
				abs := dealer.CurScore + dealer.Score
				dealer.CurScore -= abs
				lst[i].CurScore += abs
			}
		}

		lst[i].Score += lst[i].CurScore
	}
	dealer.Score += dealer.CurScore

	//! 记录
	var record staticfunc.Rec_Gold_Info
	record.Time = time.Now().Unix()
	record.GameType = self.room.Type

	self.State = 0

	//! 发消息
	var msg Msg_GamePTJ_End
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Ready = false

		var son Son_GamePTJ_Info
		son.Uid = self.PersonMgr[i].Uid
		son.Bets = self.PersonMgr[i].Bets
		son.Card = self.PersonMgr[i].Card
		son.Dealer = self.PersonMgr[i].Dealer
		son.Score = self.PersonMgr[i].CurScore
		son.Total = self.PersonMgr[i].Score
		son.CT = self.PersonMgr[i].CT
		son.View = self.PersonMgr[i].View
		msg.Info = append(msg.Info, son)

		var rec staticfunc.Son_Rec_Gold_Person
		rec.Uid = self.PersonMgr[i].Uid
		rec.Name = self.room.GetName(self.PersonMgr[i].Uid)
		rec.Head = self.room.GetHead(self.PersonMgr[i].Uid)
		rec.Score = self.PersonMgr[i].CurScore
		record.Info = append(record.Info, rec)

		self.room.Param[i] = self.PersonMgr[i].Score

		self.PersonMgr[i].Init()
		self.PersonMgr[i].View = true
	}
	recordinfo := lib.HF_JtoA(&record)
	for i := 0; i < len(record.Info); i++ {
		GetServer().InsertRecord(self.room.Type, record.Info[i].Uid, recordinfo, -record.Info[i].Score)
	}
	self.room.broadCastMsg("gameptjend", &msg)

	self.SetTime(30)

	//if self.room.IsBye() {
	//	self.OnBye()
	//	self.room.Bye()
	//	return
	//}

	for i := 0; i < len(self.room.Viewer); {
		person := GetPersonMgr().GetPerson(self.room.Viewer[i])
		if person == nil {
			i++
			continue
		}
		if self.room.Seat(self.room.Viewer[i]) {
			_person := new(Game_GoldPTJ_Person)
			_person.Init()
			_person.Uid = person.Uid
			_person.Score = person.Gold
			_person.Gold = person.Gold
			_person.Ready = false
			_person.View = true
			self.PersonMgr = append(self.PersonMgr, _person)
		} else {
			i++
		}
	}

	self.room.flush()
}

func (self *Game_GoldPTJ) OnBye() {
}

func (self *Game_GoldPTJ) OnExit(uid int64) {
	if self.room.Begin {
		return
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			//! 退出房间同步金币
			gold := self.PersonMgr[i].Score - self.PersonMgr[i].Gold
			if gold > 0 {
				GetRoomMgr().AddCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
			} else if gold < 0 {
				GetRoomMgr().CostCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, -gold, self.room)
			}
			self.PersonMgr[i].Gold = self.PersonMgr[i].Score

			copy(self.PersonMgr[i:], self.PersonMgr[i+1:])
			self.PersonMgr = self.PersonMgr[:len(self.PersonMgr)-1]
			break
		}
	}

	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Ready {
			num++
		}
	}

	if num == len(self.room.Uid)+len(self.room.Viewer) && num >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 准备的人数达到游戏最小人数
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
		self.OnBegin()
		return
	}

	if len(self.room.Uid) < lib.HF_Atoi(self.room.csv["minnum"]) {
		self.SetTime(0)
	}
}

func (self *Game_GoldPTJ) getInfo(uid int64) *Msg_GameGoldPTJ_Info {

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "```````````````````````````````````````````````````````````````` uid : ", uid, "  GetServer().IsAdmin(uid, staticfunc.ADMIN_PTJ) :  ", GetServer().IsAdmin(uid, staticfunc.ADMIN_PTJ))

	var msg Msg_GameGoldPTJ_Info
	msg.Begin = self.room.Begin
	msg.State = self.State
	if self.Time != 0 {
		msg.Time = self.Time - time.Now().Unix()
	}
	msg.Info = make([]Son_GameGoldPTJ_Info, 0)
	for _, value := range self.PersonMgr {
		var son Son_GameGoldPTJ_Info
		son.Uid = value.Uid
		son.Bets = value.Bets
		son.Dealer = value.Dealer
		son.Score = value.CurScore
		son.Total = value.Score
		son.View = value.View
		son.Ready = value.Ready
		son.RobDeal = value.RobDeal
		if value.Uid == uid || !msg.Begin || GetServer().IsAdmin(uid, staticfunc.ADMIN_PTJ) {
			son.Card = value.Card
			son.CT = value.CT
		} else {
			son.Card = make([]int, len(value.Card))
		}
		msg.Info = append(msg.Info, son)
	}
	return &msg
}

func (self *Game_GoldPTJ) GetPerson(uid int64) *Game_GoldPTJ_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

func (self *Game_GoldPTJ) GetWin(deal *Game_GoldPTJ_Person, xian *Game_GoldPTJ_Person) int {
	num := 0

	var card1 [2]int
	card1[0] = self.GetPTJValue(deal.CT[0])
	card1[1] = self.GetPTJValue(deal.CT[1])

	var card2 [2]int
	card2[0] = self.GetPTJValue(xian.CT[0])
	card2[1] = self.GetPTJValue(xian.CT[1])

	if card1[0] >= card2[0] {
		num++
	} else {
		num--
	}

	if card1[1] != 0 && card2[1] != 0 {
		if card1[1] >= card2[1] {
			num++
		} else {
			num--
		}
	}

	return num
}

func (self *Game_GoldPTJ) OnTime() {
	if self.Time == 0 {
		return
	}

	if time.Now().Unix() < self.Time {
		return
	}

	if !self.room.Begin {
		for i := 0; i < len(self.PersonMgr); {
			if !self.PersonMgr[i].Ready {
				self.room.KickPerson(self.PersonMgr[i].Uid, 98)
			} else {
				i++
			}
		}

		self.room.KickView()
		return
	}

	if time.Now().Unix() >= self.Time {
		if self.State == 1 { //! 抢庄
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].RobDeal < 0 {
					self.GameDeal(self.PersonMgr[i].Uid, false)
				}
			}
		} else if self.State == 2 { //! 下注
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Bets <= 0 && !self.PersonMgr[i].Dealer {
					self.GameBets(self.PersonMgr[i].Uid, 10101)
				}
			}
		} else if self.State == 3 { //! 亮牌
			size := len(self.PersonMgr)
			for i := 0; i < size; i++ {
				if !self.PersonMgr[i].View {
					if self.room.Type%10 == 0 {
						lib.GetLogMgr().Output(lib.LOG_DEBUG, self.PersonMgr[i].Card)
						ct0 := GetPTJType(self.PersonMgr[i].Card[0], self.PersonMgr[i].Card[1])
						ct1 := GetPTJType(self.PersonMgr[i].Card[2], self.PersonMgr[i].Card[3])
						if self.GetPTJValue(ct0) < self.GetPTJValue(ct1) {
							self.PersonMgr[i].Card[0], self.PersonMgr[i].Card[2] = self.PersonMgr[i].Card[2], self.PersonMgr[i].Card[0]
							self.PersonMgr[i].Card[1], self.PersonMgr[i].Card[3] = self.PersonMgr[i].Card[3], self.PersonMgr[i].Card[1]
						}
					}

					self.GameView(self.PersonMgr[i].Uid, self.PersonMgr[i].Card)
				}
			}
		}
	}
}

func (self *Game_GoldPTJ) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_GoldPTJ) OnIsBets(uid int64) bool {
	return false
}

//! 根据牌型得到牌值
func (self *Game_GoldPTJ) GetPTJValue(ct int) int {
	if ct == 0 {
		return 0
	}

	csv, ok := staticfunc.GetCsvMgr().Data["ptj"][ct]
	if !ok {
		return 0
	}

	if lib.HF_Atoi(csv["type"]) == 3 || lib.HF_Atoi(csv["type"]) == 4 {
		value2 := lib.HF_Atoi(csv["value2"])
		if value2 > 0 {
			return value2
		}
	}

	return lib.HF_Atoi(csv["value1"])
}

//! 同步总分
func (self *Game_GoldPTJ) SendTotal() {
	var msg Msg_GameKWX_Total
	for i := 0; i < len(self.PersonMgr); i++ {
		self.room.Param[i] = self.PersonMgr[i].Score
		msg.Info = append(msg.Info, Son_GameKWX_Total{self.PersonMgr[i].Uid, self.PersonMgr[i].Score})
	}
	self.room.broadCastMsg("gamegoldtotal", &msg)
}

//! 设置时间
func (self *Game_GoldPTJ) SetTime(t int) {
	if t == 0 {
		self.Time = 0
	} else {
		self.Time = time.Now().Unix() + int64(t)
	}

	var msg Msg_SetTime
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	self.room.broadCastMsg("gametime", &msg)
}

//! 结算所有人
func (self *Game_GoldPTJ) OnBalance() {
}
