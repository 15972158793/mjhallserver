package gameserver

import (
	"lib"
	"math"
	"sort"
	"staticfunc"
	"time"
)

//! 结算
type Msg_GameGoldPDKRoom_End struct {
	Info []Son_GameGoldPDKRoom_Info `json:"info"`
}

type Msg_GamePDKRoom_Step struct {
	Uid     int64 `json:"uid"`     //! 哪个uid
	Cards   []int `json:"cards"`   //! 出的啥牌
	CurStep int64 `json:"curstep"` //! 下局谁出
}

//! 金币场跑得快
type Msg_GameGoldPDKRoom_Info struct {
	Begin    bool                       `json:"begin"`    //! 是否开始
	Info     []Son_GameGoldPDKRoom_Info `json:"info"`     //! 人的info
	CurStep  int64                      `json:"curstep"`  //! 这局谁出
	BefStep  int64                      `json:"befstep"`  //! 上局谁出
	LastCard []int                      `json:"lastcard"` //! 最后的牌
	Time     int64                      `json:"time"`
}

type Son_GameGoldPDKRoom_Info struct {
	Uid   int64 `json:"uid"`
	Card  []int `json:"card"`
	Total int   `json:"total"`
	Score int   `json:"score"`
	Ready bool  `json:"ready"`
	Trust bool  `json:"trust"`
}

type Game_GoldPDKRoom_Person struct {
	Uid      int64 `json:"uid"`
	Card     []int `json:"card"`  //! 手牌
	Ready    bool  `json:"ready"` //! 是否准备
	Trust    bool  `json:"trust"` //! 是否托管
	Total    int   `json:"total"`
	Gold     int   `json:"gold"`
	CurScore int   `json:"curscore"`
	Boom     int   `json:"boom"`
}

func (self *Game_GoldPDKRoom_Person) Init() {
	self.Card = make([]int, 0)
	self.Trust = false
	self.CurScore = 0
	self.Boom = 0
}

//! 同步金币
func (self *Game_GoldPDKRoom_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

//! 设置托管
func (self *Game_GoldPDKRoom_Person) SetTrust(trust bool) {
	if self.Trust == trust {
		return
	}

	self.Trust = trust

	var msg Msg_GameDeal
	msg.Uid = self.Uid
	msg.Ok = trust

	person := GetPersonMgr().GetPerson(self.Uid)
	if person != nil {
		person.SendMsg("gametrust", &msg)
	}
}

//! 是否托管
func (self *Game_GoldPDKRoom_Person) IsTrush(_time int64) bool {
	if self.Trust {
		return true
	}

	if time.Now().Unix() < _time {
		return false
	}

	self.SetTrust(true)

	return true
}

type Game_GoldPDKRoom struct {
	PersonMgr []*Game_GoldPDKRoom_Person `json:"personmgr"`
	LastCard  []int                      `json:"lastcard"` //! 最后出的牌
	CurStep   int64                      `json:"curstep"`  //! 谁出牌
	BefStep   int64                      `json:"befstep"`  //! 上局谁出
	DF        int                        `json:"df"`       //! 底分
	Time      int64                      `json:"time"`     //! 自动选择时间
	BP        int64                      `json:"bp"`       //! 包赔id
	Card      int                        `json:"card"`

	room *Room
}

func NewGame_GoldPDKRoom() *Game_GoldPDKRoom {
	game := new(Game_GoldPDKRoom)
	game.PersonMgr = make([]*Game_GoldPDKRoom_Person, 0)

	return game
}

func (self *Game_GoldPDKRoom) GetPerson(uid int64) *Game_GoldPDKRoom_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

//! 得到下一个uid
func (self *Game_GoldPDKRoom) GetNextUid() int64 {
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

//! 得到上一个uid
func (self *Game_GoldPDKRoom) GetBeforeUid() *Game_GoldPDKRoom_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid != self.BefStep {
			continue
		}

		if i-1 >= 0 {
			return self.PersonMgr[i-1]
		} else {
			return self.PersonMgr[len(self.PersonMgr)-1]
		}
	}

	return nil
}

func (self *Game_GoldPDKRoom) OnInit(room *Room) {
	self.room = room

	if self.room.Param1%10 == 0 {
		self.DF = 50
	} else if self.room.Param1%10 == 1 {
		self.DF = 100
	} else if self.room.Param1%10 == 2 {
		self.DF = 200
	} else if self.room.Param1%10 == 3 {
		self.DF = 300
	} else if self.room.Param1%10 == 4 {
		self.DF = 500
	} else {
		self.DF = 1000
	}
}

func (self *Game_GoldPDKRoom) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldPDKRoom) OnSendInfo(person *Person) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			self.PersonMgr[i].SynchroGold(person.Gold)
			person.SendMsg("gamegoldpdkinfo", self.getInfo(person.Uid))
			return
		}
	}

	_person := new(Game_GoldPDKRoom_Person)
	_person.Init()
	_person.Uid = person.Uid
	_person.Ready = false
	_person.Total = person.Gold
	_person.Gold = person.Gold
	_person.Trust = false
	self.PersonMgr = append(self.PersonMgr, _person)

	if len(self.PersonMgr) >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 人满了
		self.SetTime(5)
	}

	person.SendMsg("gamegoldpdkinfo", self.getInfo(person.Uid))
}

func (self *Game_GoldPDKRoom) OnMsg(msg *RoomMsg) {
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
	case "gamesteps": //! 出牌
		person := self.GetPerson(msg.Uid)
		if person != nil {
			person.SetTrust(false)
		}
		self.GameStep(msg.Uid, msg.V.(*Msg_GameSteps).Cards)
	case "gametrust": //! 托管
		self.GameTrust(msg.Uid, msg.V.(*Msg_GameDeal).Ok)
	}
}

func (self *Game_GoldPDKRoom) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)

	//! 扣除底分
	for i := 0; i < len(self.PersonMgr); i++ {
		cost := int(math.Ceil(float64(self.DF) * 50.0 / 100.0))
		self.PersonMgr[i].Total -= cost
		GetServer().SqlAgentGoldLog(self.PersonMgr[i].Uid, cost, self.room.Type)
		GetServer().SqlAgentBillsLog(self.PersonMgr[i].Uid, cost, self.room.Type)
	}
	self.SendTotal()

	cardmgr := NewCard_GoldPDK()
	self.LastCard = make([]int, 0)
	self.BefStep = 0

	dealindex := lib.HF_GetRandom(len(self.PersonMgr))
	self.CurStep = self.PersonMgr[dealindex].Uid
	self.PersonMgr[dealindex].Card = append(self.PersonMgr[dealindex].Card, cardmgr.DealCard(34))
	self.PersonMgr[dealindex].Card = append(self.PersonMgr[dealindex].Card, cardmgr.Deal(15)...)
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid != self.CurStep {
			self.PersonMgr[i].Card = cardmgr.Deal(16)
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gamegoldpdkbegin", self.getInfo(person.Uid))
	}

	self.SetTime(20)
	self.room.flush()
}

//! 托管
func (self *Game_GoldPDKRoom) GameTrust(uid int64, ok bool) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	person.SetTrust(ok)
}

func (self *Game_GoldPDKRoom) OnEnd() {
	self.room.SetBegin(false)
	self.Time = 0

	//! 记录
	var record staticfunc.Rec_Gold_Info
	record.Time = time.Now().Unix()
	record.GameType = self.room.Type

	self.SetTime(10)

	var bp *Game_GoldPDKRoom_Person = nil
	if self.BP != 0 {
		before := self.GetPerson(self.BP)
		if self.GetHasSomeCard(before.Card) {
			bp = before
		}
	}
	total := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if len(self.PersonMgr[i].Card) == 16 { //! 关门
			score := 0
			if bp != nil {
				score = lib.HF_MinInt(bp.Total, 16*self.DF*2)
				bp.CurScore -= score
			} else {
				score = lib.HF_MinInt(self.PersonMgr[i].Total, 16*self.DF*2)
				self.PersonMgr[i].CurScore -= score
			}
			total += score
		} else if len(self.PersonMgr[i].Card) > 1 { //! 输
			score := 0
			if bp != nil {
				score = lib.HF_MinInt(bp.Total, len(self.PersonMgr[i].Card)*self.DF)
				bp.CurScore -= score
			} else {
				score = lib.HF_MinInt(self.PersonMgr[i].Total, len(self.PersonMgr[i].Card)*self.DF)
				self.PersonMgr[i].CurScore -= score
			}
			total += score
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if len(self.PersonMgr[i].Card) == 0 { //! 赢了
			self.PersonMgr[i].CurScore += total
			break
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Total += self.PersonMgr[i].CurScore
	}

	//! 发消息
	var msg Msg_GameGoldPDKRoom_End
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Ready = false

		var son Son_GameGoldPDKRoom_Info
		son.Uid = self.PersonMgr[i].Uid
		son.Card = self.PersonMgr[i].Card
		son.Total = self.PersonMgr[i].Total
		son.Score = self.PersonMgr[i].CurScore
		msg.Info = append(msg.Info, son)

		var rec staticfunc.Son_Rec_Gold_Person
		rec.Uid = self.PersonMgr[i].Uid
		rec.Name = self.room.GetName(self.PersonMgr[i].Uid)
		rec.Head = self.room.GetHead(self.PersonMgr[i].Uid)
		rec.Score = self.PersonMgr[i].CurScore + self.PersonMgr[i].Boom
		record.Info = append(record.Info, rec)

		self.room.Param[i] = self.PersonMgr[i].Total

		self.PersonMgr[i].Init()
	}
	recordinfo := lib.HF_JtoA(&record)
	for i := 0; i < len(record.Info); i++ {
		GetServer().InsertRecord(self.room.Type, record.Info[i].Uid, recordinfo, -record.Info[i].Score)
	}
	self.room.broadCastMsg("gamegoldpdkend", &msg)

	//if self.room.IsBye() {
	//	self.OnBye()
	//	self.room.Bye()
	//	return
	//}

	self.room.flush()
}

func (self *Game_GoldPDKRoom) OnBye() {
}

func (self *Game_GoldPDKRoom) OnExit(uid int64) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			//! 退出房间同步金币
			gold := self.PersonMgr[i].Total - self.PersonMgr[i].Gold
			if gold > 0 {
				GetRoomMgr().AddCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
			} else if gold < 0 {
				GetRoomMgr().CostCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, -gold, self.room)
			}
			self.PersonMgr[i].Gold = self.PersonMgr[i].Total

			copy(self.PersonMgr[i:], self.PersonMgr[i+1:])
			self.PersonMgr = self.PersonMgr[:len(self.PersonMgr)-1]

			//! 有人退出之后取消自动操作
			self.SetTime(0)
			break
		}
	}
}

//! 准备,第一局自动准备
func (self *Game_GoldPDKRoom) GameReady(uid int64) {
	if self.room.IsBye() {
		return
	}

	if self.room.Begin {
		return
	}

	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Ready {
				return
			} else {
				if self.PersonMgr[i].Total < self.DF*20 { //! 携带的金币不足，踢出去
					self.room.KickPerson(uid, 1)
					return
				}

				self.PersonMgr[i].Ready = true
				num++
			}
		} else if self.PersonMgr[i].Ready {
			num++
		}
	}

	if num == len(self.room.Uid) && num >= lib.HF_Atoi(self.room.csv["minnum"]) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
		self.OnBegin()
		return
	}

	var msg staticfunc.Msg_Uid
	msg.Uid = uid
	self.room.broadCastMsg("gameready", &msg)

	self.room.flush()
}

//! 出牌(玩家选择)
func (self *Game_GoldPDKRoom) GameStep(uid int64, card []int) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	if self.CurStep != uid {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "当前不是你的局")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person")
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
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到出牌:", card)
				return
			}
		}

		tmp := IsOkByGoldPDKCards(card, person.Card)
		if tmp == 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌错误:", card)
			return
		}
		if len(self.LastCard) != 0 && self.BefStep != uid {
			_tmp := IsOkByGoldPDKCards(self.LastCard, make([]int, 0))
			if tmp%100 == TYPE_CARD_ZHA { //! 炸弹
				if _tmp%100 == TYPE_CARD_ZHA {
					if CardCompare(tmp/10, _tmp/10) <= 0 {
						lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌更小:", card, ",", self.LastCard)
						return
					}
				}
			} else {
				//if len(card) != len(self.LastCard) {
				//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌张数不匹配:", card, ",", self.LastCard)
				//	return
				//}
				if tmp%100 != _tmp%100 { //! 类型不同
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌不匹配:", card, ",", self.LastCard)
					return
				}
				if CardCompare(tmp/10, _tmp/10) <= 0 {
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌更小:", card, ",", self.LastCard)
					return
				}
			}
		} else if len(card) == 1 {
			self.BP = uid
			self.Card = card[0]
		}
		self.LastCard = card
		self.BefStep = uid

		for _, value := range card {
			for i := 0; i < len(person.Card); i++ {
				if person.Card[i] == value {
					copy(person.Card[i:], person.Card[i+1:])
					person.Card = person.Card[:len(person.Card)-1]
					break
				}
			}
		}

		if tmp%100 == TYPE_CARD_ZHA { //! 炸弹
			person.Total += self.DF * 10
			person.Boom += self.DF * 10
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Uid != person.Uid {
					self.PersonMgr[i].Total -= self.DF * 5
					self.PersonMgr[i].Boom -= self.DF * 5
				}
			}
			self.SendTotal()
		}
	} else {
		if len(self.LastCard) == 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "第一局不能跳过")
			return
		} else {
			if len(self.GetStepCard(person)) > 0 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "有牌必须管", self.LastCard, "...", person.Card)
				return
			}
		}
	}

	self.CurStep = self.GetNextUid()

	self.SetTime(20)

	var msg Msg_GamePDKRoom_Step
	msg.Uid = uid
	msg.Cards = card
	msg.CurStep = self.CurStep
	self.room.broadCastMsg("gamegoldpdkstep", &msg)

	if len(person.Card) == 0 { //! 牌出完了
		if self.BP == uid {
			self.BP = 0
		}
		self.OnEnd()
		return
	} else if self.BP != uid {
		self.BP = 0
	}

	self.room.flush()
}

func (self *Game_GoldPDKRoom) getInfo(uid int64) *Msg_GameGoldPDKRoom_Info {
	var msg Msg_GameGoldPDKRoom_Info
	if self.Time != 0 {
		msg.Time = self.Time - time.Now().Unix()
	}
	msg.Begin = self.room.Begin
	msg.CurStep = self.CurStep
	msg.BefStep = self.BefStep
	msg.LastCard = self.LastCard
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameGoldPDKRoom_Info
		son.Uid = self.PersonMgr[i].Uid
		if son.Uid == uid || !msg.Begin {
			son.Card = self.PersonMgr[i].Card
		} else {
			son.Card = make([]int, len(self.PersonMgr[i].Card))
		}
		son.Ready = self.PersonMgr[i].Ready
		son.Total = self.PersonMgr[i].Total
		son.Trust = self.PersonMgr[i].Trust
		son.Score = self.PersonMgr[i].CurScore
		msg.Info = append(msg.Info, son)
	}
	return &msg
}

//! 同步总分
func (self *Game_GoldPDKRoom) SendTotal() {
	var msg Msg_GameKWX_Total
	for i := 0; i < len(self.PersonMgr); i++ {
		self.room.Param[i] = self.PersonMgr[i].Total
		msg.Info = append(msg.Info, Son_GameKWX_Total{self.PersonMgr[i].Uid, self.PersonMgr[i].Total})
	}
	self.room.broadCastMsg("gamegoldtotal", &msg)
}

func (self *Game_GoldPDKRoom) OnTime() {
	if self.Time == 0 {
		return
	}

	//! 60秒之后自动选择
	if !self.room.Begin {
		if time.Now().Unix() < self.Time {
			return
		}
		for i := 0; i < len(self.PersonMgr); {
			if !self.PersonMgr[i].Ready {
				self.room.KickPerson(self.PersonMgr[i].Uid, 98)
			} else {
				i++
			}
		}

		return
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == self.CurStep {
			if self.PersonMgr[i].IsTrush(self.Time) {
				self.GameStep(self.PersonMgr[i].Uid, self.GetStepCard(self.PersonMgr[i]))
			}
			return
		}
	}
}

func (self *Game_GoldPDKRoom) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_GoldPDKRoom) OnIsBets(uid int64) bool {
	return false
}

//! 设置时间
func (self *Game_GoldPDKRoom) SetTime(t int) {
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
func (self *Game_GoldPDKRoom) OnBalance() {
	for i := 0; i < len(self.PersonMgr); i++ {
		//! 退出房间同步金币
		gold := self.PersonMgr[i].Total - self.PersonMgr[i].Gold
		if gold > 0 {
			GetRoomMgr().AddCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
		} else if gold < 0 {
			GetRoomMgr().CostCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, -gold, self.room)
		}
		self.PersonMgr[i].Gold = self.PersonMgr[i].Total
	}
}

//! 得到有没有牌大于上家
func (self *Game_GoldPDKRoom) GetStepCard(person *Game_GoldPDKRoom_Person) []int {
	stepcard := make([]int, 0)

	if len(self.LastCard) == 0 || self.BefStep == person.Uid { //! 上家没有出牌
		stepcard = append(stepcard, person.Card[0])
		return stepcard
	}

	tmp := IsOkByGoldPDKCards(self.LastCard, make([]int, 0))

	if tmp%100 == TYPE_CARD_ONE { //! 单张
		for i := 0; i < len(person.Card); i++ {
			if CardCompare(person.Card[i], tmp/10) > 0 {
				stepcard = append(stepcard, person.Card[i])
				return stepcard
			}
		}
	}

	handcard := make(map[int]int)
	for i := 0; i < len(person.Card); i++ {
		handcard[person.Card[i]/10]++
	}

	one := make(LstPoker, 0)
	two := make(LstPoker, 0)
	three := make(LstPoker, 0)
	four := make(LstPoker, 0)
	for key, value := range handcard {
		if value == 1 {
			one = append(one, key*10)
			sort.Sort(LstPoker(one))
		} else if value == 2 {
			two = append(two, key*10)
			sort.Sort(LstPoker(two))
		} else if value == 3 {
			three = append(three, key*10)
			sort.Sort(LstPoker(three))
		} else if value == 4 {
			four = append(four, key*10)
			sort.Sort(LstPoker(four))
		}
	}

	if tmp%100 == TYPE_CARD_TWO { //! 对子
		for i := 0; i < len(two); i++ {
			if CardCompare(two[i], tmp/10) > 0 {
				stepcard = append(stepcard, self.GetHasCard(person.Card, 2, two[i])...)
				return stepcard
			}
		}

		for i := 0; i < len(three); i++ {
			if CardCompare(three[i], tmp/10) > 0 {
				stepcard = append(stepcard, self.GetHasCard(person.Card, 2, three[i])...)
				return stepcard
			}
		}
	}

	if tmp%100 == TYPE_CARD_ZHA { //! 炸弹
		for i := 0; i < len(four); i++ {
			if CardCompare(four[i], tmp/10) > 0 {
				stepcard = append(stepcard, self.GetHasCard(person.Card, 4, four[i])...)
				return stepcard
			}
		}
	}

	if tmp%100 == TYPE_CARD_SAN { //! 3张
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找到牌型:3张")
		for i := 0; i < len(three); i++ {
			if CardCompare(three[i], tmp/10) > 0 {
				stepcard = append(stepcard, self.GetHasCard(person.Card, 3, three[i])...)
				lib.GetLogMgr().Output(lib.LOG_DEBUG, stepcard)
				if len(two) > 0 {
					stepcard = append(stepcard, self.GetHasCard(person.Card, 2, two[0])...)
					lib.GetLogMgr().Output(lib.LOG_DEBUG, stepcard)
					return stepcard
				} else if len(three) > 1 {
					for j := 0; j < len(three); j++ {
						if i == j {
							continue
						}
						stepcard = append(stepcard, self.GetHasCard(person.Card, 2, three[j])...)
						lib.GetLogMgr().Output(lib.LOG_DEBUG, stepcard)
						return stepcard
					}
				} else if len(four) > 0 {
					stepcard = append(stepcard, self.GetHasCard(person.Card, 2, four[0])...)
					lib.GetLogMgr().Output(lib.LOG_DEBUG, stepcard)
					return stepcard
				} else if len(one) > 0 {
					for j := 0; j < lib.HF_MinInt(2, len(one)); j++ {
						stepcard = append(stepcard, self.GetHasCard(person.Card, 1, one[j])...)
						lib.GetLogMgr().Output(lib.LOG_DEBUG, stepcard)
					}
					return stepcard
				}
				lib.GetLogMgr().Output(lib.LOG_DEBUG, stepcard)
				return stepcard
			}
		}
	}

	if tmp%100 == TYPE_CARD_SHUN { //! 顺子
		_card := make([]int, 0)
		_card = append(_card, one...)
		_card = append(_card, two...)
		_card = append(_card, three...)
		_card = append(_card, four...)
		sort.Sort(LstPoker(_card))
		sort.Sort(LstPoker(self.LastCard))
		if len(_card) >= len(self.LastCard) {
			index := -1
			for i := 0; i < len(_card); i++ {
				if CardCompare(_card[i], self.LastCard[0]) > 0 {
					index = i
					break
				}
			}
			if index >= 0 {
				for {
					if index+len(self.LastCard)-1 >= len(_card) {
						break
					}
					find := true
					for i := 1; i <= len(self.LastCard)-1; i++ {
						if _card[index+i]/10 == 2 || _card[index+i]/10 == 14 {
							find = false
							break
						}
						if _card[index+i]/10 == 1 {
							_card[index+i] = 140
						}
						if _card[index+i]/10-_card[index+i-1]/10 != 1 {
							find = false
							break
						}
					}
					if find {
						for i := index; i < index+len(self.LastCard); i++ {
							stepcard = append(stepcard, self.GetHasCard(person.Card, 1, _card[i])...)
						}
						return stepcard
					}
					index++
				}
			}
		}
	}

	if tmp%100 == TYPE_CARD_SHUNDUI { //! 顺对
		_card := make([]int, 0)
		_card = append(_card, two...)
		_card = append(_card, three...)
		_card = append(_card, four...)
		sort.Sort(LstPoker(_card))
		sort.Sort(LstPoker(self.LastCard))
		if len(_card) >= len(self.LastCard)/2 {
			index := -1
			for i := 0; i < len(_card); i++ {
				if CardCompare(_card[i], self.LastCard[0]) > 0 {
					index = i
					break
				}
			}
			if index >= 0 {
				for {
					if index+len(self.LastCard)/2-1 >= len(_card) {
						break
					}
					find := true
					for i := 1; i <= len(self.LastCard)/2-1; i++ {
						if _card[index+i]/10 == 2 || _card[index+i]/10 == 14 {
							find = false
							break
						}
						if _card[index+i]/10 == 1 {
							_card[index+i] = 140
						}
						if _card[index+i]/10-_card[index+i-1]/10 != 1 {
							find = false
							break
						}
					}
					if find {
						for i := index; i < index+len(self.LastCard)/2; i++ {
							stepcard = append(stepcard, self.GetHasCard(person.Card, 2, _card[i])...)
						}
						return stepcard
					}
					index++
				}
			}
		}
	}

	if tmp%100 == TYPE_CARD_SHUNSAN { //! 顺3
		_card := make([]int, 0)
		_card = append(_card, three...)
		_card = append(_card, four...)
		sort.Sort(LstPoker(_card))
		last := IsOkByGoldPDKCards(self.LastCard, make([]int, 0))
		if len(_card) >= len(self.LastCard)/3 {
			index := -1
			for i := 0; i < len(_card); i++ {
				if CardCompare(_card[i], last/10) > 0 {
					index = i
					break
				}
			}
			if index >= 0 {
				for {
					if index+len(self.LastCard)/3-1 >= len(_card) {
						break
					}
					find := true
					for i := 1; i <= len(self.LastCard)/3-1; i++ {
						if _card[index+i]/10 == 2 || _card[index+i]/10 == 14 {
							find = false
							break
						}
						if _card[index+i]/10 == 1 {
							_card[index+i] = 140
						}
						if _card[index+i]/10-_card[index+i-1]/10 != 1 {
							find = false
							break
						}
					}
					if find {
						for i := index; i < index+len(self.LastCard)/3; i++ {
							stepcard = append(stepcard, self.GetHasCard(person.Card, 3, _card[i])...)
						}
						return stepcard
					}
					index++
				}
			}
		}
	}

	if tmp%100 == TYPE_CARD_SI1 { //! 4带
		for i := 0; i < len(four); i++ {
			if CardCompare(four[i], tmp/10) > 0 {
				stepcard = append(stepcard, self.GetHasCard(person.Card, 4, four[i])...)
				if len(two) > 0 {
					stepcard = append(stepcard, self.GetHasCard(person.Card, 2, two[0])...)
					return stepcard
				} else if len(three) > 0 {
					stepcard = append(stepcard, self.GetHasCard(person.Card, 2, three[0])...)
					return stepcard
				} else if len(four) > 1 {
					for j := 0; j < len(four); j++ {
						if i == j {
							continue
						}
						stepcard = append(stepcard, self.GetHasCard(person.Card, 2, four[j])...)
						return stepcard
					}
				} else if len(one) > 1 {
					for j := 0; j < 2; j++ {
						stepcard = append(stepcard, self.GetHasCard(person.Card, 1, one[j])...)
					}
					return stepcard
				}
				return stepcard
			}
		}
	}

	if tmp%100 != TYPE_CARD_ZHA { //! 不是炸弹
		for i := 0; i < len(four); i++ {
			stepcard = append(stepcard, self.GetHasCard(person.Card, 4, four[i])...)
			return stepcard
		}
	}

	return stepcard
}

func (self *Game_GoldPDKRoom) GetHasCard(card []int, num int, key int) []int {
	tmp := make([]int, 0)
	for i := 0; i < len(card); i++ {
		if CardCompare(card[i], key) == 0 {
			tmp = append(tmp, card[i])
			if len(tmp) >= num {
				break
			}
		}
	}

	return tmp
}

func (self *Game_GoldPDKRoom) GetHasSomeCard(card []int) bool {
	handcard := make(map[int]int)
	for i := 0; i < len(card); i++ {
		handcard[card[i]/10]++
	}

	for _, value := range handcard {
		if value >= 2 {
			return true
		}
	}

	for i := 0; i < len(card); i++ {
		if CardCompare(card[i], self.Card) > 0 {
			return true
		}
	}

	return false
}
