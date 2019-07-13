package gameserver

import (
	"lib"
	"math"
	"sort"
	"staticfunc"
	"time"
)

/*
个位 炸弹上限 0-3 1-4 2-5
十位 0-叫分 1-不叫分
百位 0-双王不可拆  1-双王可拆
*/

//!
type Msg_GameGoldDDZ_Info struct {
	Begin       bool                   `json:"begin"` //! 是否开始
	Deal        []Son_GameGoldDDZ_Deal `json:"deal"`  //! 抢地主的人
	Info        []Son_GameGoldDDZ_Info `json:"info"`
	CurStep     int64                  `json:"curstep"`     //! 当前谁出牌
	BefStep     int64                  `json:"befstep"`     //! 上局谁出
	LastCard    []int                  `json:"lastcard"`    //! 最后的牌
	LastAbsCard []int                  `json:"lastabscard"` //! 实际出牌
	DZCard      []int                  `json:"dzcard"`      //! 地主牌
	Bets        int                    `json:"bets"`        //! 底分
	Boom        int                    `json:"boom"`        //! 炸弹
	Razz        int                    `json:"razz"`
}

type Son_GameGoldDDZ_Info struct {
	Uid    int64 `json:"uid"`
	Card   []int `json:"card"`
	Dealer bool  `json:"dealer"`
	Score  int   `json:"score"`
	Total  int   `json:"total"`
	Ready  bool  `json:"ready"`
	Bets   int   `json:"bets"`
	Double int   `json:"isdouble"`
}
type Son_GameGoldDDZ_Deal struct {
	Uid int64 `json:"uid"`
	Ok  bool  `json:"ok"`
}

//! 出牌
type Msg_GameGoldDDZ_Step struct {
	Uid      int64 `json:"uid"`      //! 哪个uid
	Cards    []int `json:"cards"`    //! 出的啥牌
	AbsCards []int `json:"abscards"` //! 实际牌
	CurStep  int64 `json:"curstep"`  //! 下局谁出
}

//! 结算
type Msg_GameGoldDDZ_End struct {
	Info []Son_GameGoldDDZ_Info `json:"info"`
	Bets int                    `json:"bets"`
	Boom int                    `json:"boom"`
	CT   bool                   `json:"ct"`
}

//! 地主
type Msg_GameGoldDDZ_Dealer struct {
	Uid  int64 `json:"uid"`
	Card []int `json:"card"`
	Bets int   `json:"bets"`
	Razz int   `json:"razz"`
}

type Msg_GameGoldDDZ_Bets struct {
	Uid     int64 `json:"uid"`
	Bets    int   `json:"bets"`
	CurStep int64 `json:"curstep"`
}

type Msg_GameGoldDDZ_Double struct {
	Double int `json:"double"`
}

///////////////////////////////////////////////////////
type Game_GoldDDZ_Person struct {
	Uid      int64 `json:"uid"`
	Card     []int `json:"card"`     //! 手牌
	Bets     int   `json:"bets"`     //! 叫分
	Double   int   `json:"isdouble"` //! 是否加倍
	Gold     int   `json:"gold"`
	Total    int   `json:"total"`    //! 积分
	Dealer   bool  `json:"dealer"`   //! 是否是地主
	CurScore int   `json:"curscore"` //! 当前局的分数
	Ready    bool  `json:"ready"`    //! 是否准备好
	//	Boom     int   `json:"boom"`     //! 炸弹数量
	Trust bool `json:"trust"` //! 是否托管
}

func (self *Game_GoldDDZ_Person) Init() {
	self.Card = make([]int, 0)
	self.Dealer = false
	self.Double = -1
}

type Game_GoldDDZ struct {
	PersonMgr   []*Game_GoldDDZ_Person `json:"personmgr"`   //! 玩家信息
	DZCard      []int                  `json:"dzcard"`      //! 地主牌
	LastCard    []int                  `json:"lastcard"`    //! 最后出的牌
	LastAbsCard []int                  `json:"lastabscard"` //! 实际出牌
	CurStep     int64                  `json:"curstep"`     //! 谁出牌
	BefStep     int64                  `json:"befstep"`     //! 上局谁出
	Winer       int64                  `json:"winer"`       //! 赢家
	Bets        int                    `json:"bets"`        //! 底分
	Boom        int                    `json:"boom"`        //! 炸弹数量
	Taxi        int                    `json:"taxi"`        //! 农民出牌次数
	DZNum       int                    `json:"dznum"`       //! 地主出牌次数
	Razz        int                    `json:"razz"`        //! 癞子
	State       int                    `json:"state"`       //!0－未开始 1-叫地主 2-加倍 3-出牌
	Time        int64                  `json:"time"`
	DF          int                    `json:"df"`

	room *Room
	hz   bool
}

func NewGame_GoldDDZ() *Game_GoldDDZ {
	game := new(Game_GoldDDZ)
	game.DZCard = make([]int, 0)
	game.LastCard = make([]int, 0)
	game.PersonMgr = make([]*Game_GoldDDZ_Person, 0)

	return game
}

func (self *Game_GoldDDZ) SetTime(t int) {
	if t == 0 {
		self.Time = 0
	} else {
		self.Time = time.Now().Unix() + int64(t)
	}

	var msg Msg_SetTime
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	self.room.broadCastMsg("gametime", &msg)
}

func (self *Game_GoldDDZ) OnInit(room *Room) {
	self.room = room

	self.DF = staticfunc.GetCsvMgr().GetDF(self.room.Type)
}

func (self *Game_GoldDDZ) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldDDZ) OnSendInfo(person *Person) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			person.SendMsg("gameddzinfo", self.getInfo(person.Uid))
			return
		}
	}

	_person := new(Game_GoldDDZ_Person)
	_person.Init()
	_person.Uid = person.Uid
	_person.Total = person.Gold
	_person.Gold = person.Gold
	_person.Ready = false

	self.PersonMgr = append(self.PersonMgr, _person)

	if len(self.PersonMgr) >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 人满了
		self.SetTime(5)
	}

	person.SendMsg("gameddzinfo", self.getInfo(person.Uid))
}

func (self *Game_GoldDDZ) OnMsg(msg *RoomMsg) {

	switch msg.Head {
	case "gamebets": //! 叫分
		self.GameBets(msg.Uid, msg.V.(*Msg_GameBets).Bets)
		//	case "gamedouble": //! 双倍
		//		self.GameDouble(msg.Uid, msg.V.(*Msg_GameGoldDDZ_Double).Double)
	case "gamesteps": //! 出牌
		self.GameStep(msg.Uid, msg.V.(*Msg_GameSteps).Cards, msg.V.(*Msg_GameSteps).AbsCards)
	case "gameready": //! 游戏准备
		self.GameReady(msg.Uid)
	case "gametrust": //! 托管
		self.GameTrust(msg.Uid, msg.V.(*Msg_GameDeal).Ok)
	}
}

//! 托管
func (self *Game_GoldDDZ) GameTrust(uid int64, ok bool) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	person.SetTrust(ok)
}

//! 同步总分
func (self *Game_GoldDDZ) SendTotal() {
	var msg Msg_GameKWX_Total
	for i := 0; i < len(self.PersonMgr); i++ {
		self.room.Param[i] = self.PersonMgr[i].Total
		msg.Info = append(msg.Info, Son_GameKWX_Total{self.PersonMgr[i].Uid, self.PersonMgr[i].Total})
	}
	self.room.broadCastMsg("gamegoldtotal", &msg)
}

func (self *Game_GoldDDZ) OnBegin() {
	if self.room.IsBye() {
		return
	}

	//! 扣除底分
	for i := 0; i < len(self.PersonMgr); i++ {
		cost := int(math.Ceil(float64(self.DF) * 50.0 / 100.0))
		self.PersonMgr[i].Total -= cost

	}
	self.SendTotal()

	self.SetTime(60)

	self.room.SetBegin(true)
	self.State = 1

	if self.Winer == 0 {
		self.Winer = self.PersonMgr[0].Uid
	}

	cardmgr := NewCard_DDZ()
	self.DZCard = cardmgr.Deal(3)
	self.LastCard = make([]int, 0)
	self.LastAbsCard = make([]int, 0)
	self.BefStep = 0
	self.CurStep = self.Winer
	self.hz = false
	self.Razz = 0
	self.Boom = 0
	self.Bets = 0
	self.Taxi = 0
	self.DZNum = 0
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Dealer = false
		self.PersonMgr[i].Bets = -1
		self.PersonMgr[i].Double = -1
		self.PersonMgr[i].Trust = false
		self.PersonMgr[i].Card = cardmgr.Deal(17)
		self.PersonMgr[i].CurScore = 0
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gameddzbegin", self.getInfo(person.Uid))
	}

	self.room.flush()
}

//! 下注
func (self *Game_GoldDDZ) GameBets(uid int64, bets int) {
	if !self.room.Begin { //! 未开始不能叫分
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "未开始不能抢庄")
		return
	}

	if uid != self.CurStep {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不归你叫")
		return
	}

	if bets == 3 || (self.room.Type%290000 == 1 && bets >= 2) { //!
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Uid == uid {
				self.PersonMgr[i].Dealer = true
				self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, self.DZCard...)
				break
			}
		}

		if self.room.Param1/10%10 == 1 {
			self.Bets = 1
		} else {
			if self.room.Type%290000 == 1 {
				self.Bets = 2
			} else {
				self.Bets = 3
			}
		}
		self.CurStep = uid
		if self.room.Type%290000 == 1 {
			self.Razz = lib.HF_GetRandom(13) + 1
		}

		self.State = 2
		var msg Msg_GameGoldDDZ_Dealer
		msg.Uid = uid
		msg.Card = self.DZCard
		msg.Bets = self.Bets
		msg.Razz = self.Razz
		self.room.broadCastMsg("gamedealer", &msg)
	} else {
		num, max, muid := 0, -1, int64(0)
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Bets >= 0 {
				num++
				if self.PersonMgr[i].Bets > max {
					max = self.PersonMgr[i].Bets
					muid = self.PersonMgr[i].Uid
				}
			}
		}

		if num >= len(self.PersonMgr)-1 { //! 最后一个人表态
			if bets > max {
				max = bets
				muid = uid
			}

			if max == 0 { //! 都不叫庄
				self.hz = true
				self.OnEnd()
				return
			}

			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Uid == muid {
					self.PersonMgr[i].Dealer = true
					self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, self.DZCard...)
					break
				}
			}

			self.Bets = lib.HF_MaxInt(max, 1)
			self.CurStep = muid
			if self.room.Type%290000 == 1 {
				self.Razz = lib.HF_GetRandom(13) + 1
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "癞子是:", self.Razz)
			}

			self.State = 2
			var msg Msg_GameGoldDDZ_Dealer
			msg.Uid = muid
			msg.Card = self.DZCard
			msg.Bets = self.Bets
			msg.Razz = self.Razz
			self.room.broadCastMsg("gamedealer", &msg)
		} else {
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Uid == uid {
					self.PersonMgr[i].Bets = bets
					break
				}
			}
			self.CurStep = self.GetNextUid()

			var msg Msg_GameGoldDDZ_Bets
			msg.Uid = uid
			msg.Bets = bets
			msg.CurStep = self.CurStep
			self.room.broadCastMsg("gamebets", &msg)
		}
	}
	self.SetTime(60)
	self.room.flush()
}

//! 加倍
func (self *Game_GoldDDZ) GameDouble(uid int64, double int) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "未开始不能加倍")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if double <= 0 || double > 2 {
		return
	}

	person.Double = double

	var msg Msg_GameDouble
	msg.Uid = uid
	msg.Double = double
	self.room.broadCastMsg("gamedouble", &msg)

	self.SetTime(60)
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Double == -1 {
			return
		}
	}
	self.State = 3
	self.room.flush()
}

//! 出牌(玩家选择)
func (self *Game_GoldDDZ) GameStep(uid int64, card []int, abscard []int) {
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

	if self.room.Type%290000 != 1 { //! 普通模式
		abscard = card
	} else {
		if len(abscard) != len(card) {
			return
		}
		tmpcard := make([]int, 0)
		lib.HF_DeepCopy(&tmpcard, &card)
		for _, value := range abscard {
			find := false
			for i := 0; i < len(tmpcard); i++ {
				if tmpcard[i] == value {
					find = true
					copy(tmpcard[i:], tmpcard[i+1:])
					tmpcard = tmpcard[:len(tmpcard)-1]
					break
				}
			}

			if find {
				continue
			}

			for i := 0; i < len(tmpcard); i++ {
				if tmpcard[i]/10 == self.Razz {
					find = true
					copy(tmpcard[i:], tmpcard[i+1:])
					tmpcard = tmpcard[:len(tmpcard)-1]
					break
				}
			}

			if !find {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "癞子错误:", card, ",", abscard)
				return
			}
		}
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

		tmp := IsOkByCards(abscard)

		if tmp == 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌错误:", card)
			return
		}
		if len(self.LastCard) != 0 && self.BefStep != uid {
			_tmp := IsOkByCards(self.LastAbsCard)
			if tmp == TYPE_CARD_WANG { //! 王炸
				self.Boom++
				//				person.Boom++
			} else if tmp%100 == TYPE_CARD_ZHA { //! 炸弹
				if _tmp%100 == TYPE_CARD_ZHA {
					compareboom := true //! 是否比炸弹
					if self.room.Type%290000 == 1 {
						if IsBoom(card) && !IsBoom(self.LastCard) { //! 硬炸一定压软炸
							compareboom = false
						}
					}

					if compareboom {
						if CardCompare(tmp/10, _tmp/10) <= 0 {
							lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌更小:", card, ",", self.LastCard)
							return
						}
					}
				}
				self.Boom++
				//				person.Boom++
			} else {
				if len(card) != len(self.LastCard) {
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌张数不匹配:", card, ",", self.LastCard)
					return
				}
				if tmp%100 != _tmp%100 {
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌不匹配:", card, ",", self.LastCard)
					return
				}
				if CardCompare(tmp/10, _tmp/10) <= 0 {
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌更小:", card, ",", self.LastCard)
					return
				}
			}
		} else {
			if tmp == TYPE_CARD_WANG || tmp%100 == TYPE_CARD_ZHA {
				self.Boom++
				//				person.Boom++
			}
		}
		self.LastCard = card
		self.LastAbsCard = abscard
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

	} else {
		if len(self.LastCard) == 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "第一局不能跳过")
			return
		}
	}

	self.CurStep = self.GetNextUid()
	if len(card) != 0 {
		if person.Dealer {
			self.DZNum++
		} else {
			self.Taxi++
		}
	}

	self.SetTime(20)

	var msg Msg_GameGoldDDZ_Step
	msg.Uid = uid
	msg.Cards = card
	msg.AbsCards = abscard
	msg.CurStep = self.CurStep
	self.room.broadCastMsg("gameddzstep", &msg)

	if len(person.Card) == 0 { //! 牌出完了
		self.OnEnd()
		return
	}

	self.room.flush()
}

//! 准备,第一局自动准备
func (self *Game_GoldDDZ) GameReady(uid int64) {
	if self.room.IsBye() {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "roomisbye")
		return
	}

	if self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "roomisbegin")
		return
	}

	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Ready {
				return
			} else {
				if self.PersonMgr[i].Total < staticfunc.GetCsvMgr().GetZR(self.room.Type) { //! 携带的金币不足，踢出去
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

	self.room.flush()

	var msg staticfunc.Msg_Uid
	msg.Uid = uid
	self.room.broadCastMsg("gameready", &msg)
}

//! 结算
func (self *Game_GoldDDZ) OnEnd() {
	self.room.SetBegin(false)

	self.State = 0

	var record staticfunc.Rec_Gold_Info
	record.Time = time.Now().Unix()
	record.GameType = self.room.Type

	if self.hz {
		var msg Msg_GameGoldDDZ_End
		msg.Bets = 0
		msg.Boom = 0
		msg.CT = false

		for i := 0; i < len(self.PersonMgr); i++ {
			var son Son_GameGoldDDZ_Info
			son.Uid = self.PersonMgr[i].Uid
			son.Card = self.PersonMgr[i].Card
			son.Dealer = self.PersonMgr[i].Dealer
			son.Total = self.PersonMgr[i].Total
			son.Score = self.PersonMgr[i].CurScore
			son.Double = self.PersonMgr[i].Double
			msg.Info = append(msg.Info, son)
		}

		self.room.broadCastMsg("gameddzend", &msg)

		if self.room.IsBye() {
			self.OnBye()
			self.room.Bye()
			return
		}

		for i := 0; i < len(self.PersonMgr); i++ {
			self.PersonMgr[i].Ready = false
		}
		self.room.flush()

		return
	}

	DZWin := false
	for i := 0; i < len(self.PersonMgr); i++ {
		if len(self.PersonMgr[i].Card) == 0 {
			if self.PersonMgr[i].Dealer {
				DZWin = true
			}
			self.Winer = self.PersonMgr[i].Uid
			break
		}
	}

	maxboom := 5
	addbet := false
	if self.room.Type%290000 == 1 { //! 癞子
		maxboom = 9999 //!　炸弹不封顶
		addbet = true
	} else { //! 经典
		maxboom = 5
		addbet = true
	}

	_boom := self.Boom - maxboom
	if _boom > 0 {
		self.Boom = maxboom
	}
	score := self.Bets * int(math.Pow(2.0, float64(self.Boom))) * self.DF
	if DZWin && self.Taxi == 0 {
		score *= 2
	} else if !DZWin && self.DZNum <= 1 {
		score *= 2
	}
	if _boom > 0 && addbet {
		score += self.Bets * _boom
	}

	//! 先找到地主
	var dealer *Game_GoldDDZ_Person
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Dealer {
			dealer = self.PersonMgr[i]
			break
		}
	}

	if dealer.Double == 2 {
		score *= 2
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Dealer {
			continue
		}

		bs := 1
		if self.PersonMgr[i].Double == 2 {
			bs = 2
		}

		if DZWin {
			dealer.CurScore += score * bs
			self.PersonMgr[i].CurScore -= score * bs
		} else {
			dealer.CurScore -= score * bs
			self.PersonMgr[i].CurScore += score * bs
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Total += self.PersonMgr[i].CurScore
	}

	var msg Msg_GameGoldDDZ_End
	msg.Bets = self.Bets
	msg.Boom = self.Boom
	msg.CT = false
	if DZWin && self.Taxi == 0 {
		msg.CT = true
	} else if !DZWin && self.DZNum <= 1 {
		msg.CT = true
	}
	agentinfo := make([]staticfunc.JS_CreateRoomMem, 0)
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameGoldDDZ_Info
		son.Uid = self.PersonMgr[i].Uid
		son.Card = self.PersonMgr[i].Card
		son.Dealer = self.PersonMgr[i].Dealer
		son.Total = self.PersonMgr[i].Total
		son.Score = self.PersonMgr[i].CurScore
		son.Double = self.PersonMgr[i].Double
		msg.Info = append(msg.Info, son)
		agentinfo = append(agentinfo, staticfunc.JS_CreateRoomMem{son.Uid, "", "", son.Total})

		var rec staticfunc.Son_Rec_Gold_Person
		rec.Uid = self.PersonMgr[i].Uid
		rec.Name = self.room.GetName(self.PersonMgr[i].Uid)
		rec.Head = self.room.GetHead(self.PersonMgr[i].Uid)
		rec.Score = self.PersonMgr[i].CurScore
		record.Info = append(record.Info, rec)
	}

	self.room.AddRecord(lib.HF_JtoA(&record))
	self.room.broadCastMsg("gameddzend", &msg)

	self.SetTime(60)

	if self.room.IsBye() {
		self.OnBye()
		self.room.Bye()
		return
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Ready = false
	}
	self.room.flush()
}

func (self *Game_GoldDDZ) OnBye() {

}

func (self *Game_GoldDDZ) OnExit(uid int64) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {

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

func (self *Game_GoldDDZ) getInfo(uid int64) *Msg_GameGoldDDZ_Info {
	var msg Msg_GameGoldDDZ_Info
	msg.Begin = self.room.Begin
	msg.Deal = make([]Son_GameGoldDDZ_Deal, 0)
	msg.Info = make([]Son_GameGoldDDZ_Info, 0)
	msg.CurStep = self.CurStep
	msg.BefStep = self.BefStep
	msg.Bets = self.Bets
	msg.Boom = self.Boom
	msg.Razz = self.Razz
	msg.LastAbsCard = self.LastAbsCard
	if self.CurStep == 0 {
		msg.DZCard = []int{0, 0, 0}
	} else {
		msg.DZCard = self.DZCard
	}
	msg.LastCard = self.LastCard
	for _, value := range self.PersonMgr {
		var son Son_GameGoldDDZ_Info
		son.Uid = value.Uid
		son.Dealer = value.Dealer
		if value.CurScore == 0 {
			if son.Uid == uid || GetServer().IsAdmin(uid, staticfunc.ADMIN_DDZ) {
				son.Card = value.Card
			} else {
				for i := 0; i < len(value.Card); i++ {
					son.Card = append(son.Card, 0)
				}
			}
		} else {
			son.Card = value.Card
		}
		son.Bets = value.Bets
		son.Total = value.Total
		son.Ready = value.Ready
		son.Double = value.Double
		msg.Info = append(msg.Info, son)
	}
	return &msg
}

func (self *Game_GoldDDZ) GetPerson(uid int64) *Game_GoldDDZ_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

//! 得到下一个uid
func (self *Game_GoldDDZ) GetNextUid() int64 {
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

func (self *Game_GoldDDZ) OnTime() {

	//! 60秒之后自动选择
	if !self.room.Begin {
		if self.Time == 0 {
			return
		}

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

	if self.State == 1 {
		if time.Now().Unix() < self.Time {
			return
		}
		self.GameBets(self.CurStep, 0)
	} else if self.State == 2 {
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Uid == self.CurStep {
				if self.PersonMgr[i].IsTrush(self.Time) {
					self.GameStep(self.PersonMgr[i].Uid, self.GetStepCard(self.PersonMgr[i]), self.GetStepCard(self.PersonMgr[i]))
				}
				return
			}
		}
	}
}

func (self *Game_GoldDDZ) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_GoldDDZ) OnIsBets(uid int64) bool {
	return false
}

//! 结算所有人
func (self *Game_GoldDDZ) OnBalance() {
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
func (self *Game_GoldDDZ) GetStepCard(person *Game_GoldDDZ_Person) []int {
	if len(self.LastCard) == 0 || self.BefStep == person.Uid { //! 上家没有出牌
		if len(self.GetPerson(self.GetNextUid()).Card) == 1 { //! 下家保单
			//! 试试顺子能不能出
			shuncard := self.GetBetterCard(person, []int{31, 41, 51, 61, 71, 81, 91, 101}, true)
			if len(shuncard) > 0 && !self.IsBoom(shuncard) {
				return shuncard
			}
			shuncard = self.GetBetterCard(person, []int{31, 41, 51, 61, 71, 81, 91}, true)
			if len(shuncard) > 0 && !self.IsBoom(shuncard) {
				return shuncard
			}
			shuncard = self.GetBetterCard(person, []int{31, 41, 51, 61, 71, 81}, true)
			if len(shuncard) > 0 && !self.IsBoom(shuncard) {
				return shuncard
			}
			shuncard = self.GetBetterCard(person, []int{31, 41, 51, 61, 71}, true)
			if len(shuncard) > 0 && !self.IsBoom(shuncard) {
				return shuncard
			}
			//! 试试连对
			lianduicard := self.GetBetterCard(person, []int{31, 31, 41, 41}, true)
			if len(lianduicard) > 0 && !self.IsBoom(lianduicard) {
				return lianduicard
			}
			//! 试试三带二
			sandaiercard := self.GetBetterCard(person, []int{31, 31, 31, 41, 41}, true)
			if len(sandaiercard) > 0 && !self.IsBoom(sandaiercard) {
				return sandaiercard
			}
			//! 试试三带一
			sandaiyicard := self.GetBetterCard(person, []int{31, 31, 31, 41}, true)
			if len(sandaiyicard) > 0 && !self.IsBoom(sandaiyicard) {
				return sandaiyicard
			}
			//! 试试对子
			duicard := self.GetBetterCard(person, []int{31, 31}, true)
			if len(duicard) > 0 && !self.IsBoom(duicard) {
				return duicard
			}
			//! 炸弹
			boomcard := self.GetBetterCard(person, []int{21, 21}, true)
			if len(boomcard) > 0 {
				return boomcard
			}
			//! 得到最大的牌
			return []int{self.GetMaxCard(person.Card)}
		} else {
			canstep := make([][]int, 0)
			//! 试试顺子能不能出
			shuncard := self.GetBetterCard(person, []int{31, 41, 51, 61, 71, 81, 91, 101}, true)
			if len(shuncard) > 0 && !self.IsBoom(shuncard) && !self.GetBigCard(shuncard) {
				canstep = append(canstep, shuncard)
			}
			shuncard = self.GetBetterCard(person, []int{31, 41, 51, 61, 71, 81, 91}, true)
			if len(shuncard) > 0 && !self.IsBoom(shuncard) && !self.GetBigCard(shuncard) {
				canstep = append(canstep, shuncard)
			}
			shuncard = self.GetBetterCard(person, []int{31, 41, 51, 61, 71, 81}, true)
			if len(shuncard) > 0 && !self.IsBoom(shuncard) && !self.GetBigCard(shuncard) {
				canstep = append(canstep, shuncard)
			}
			shuncard = self.GetBetterCard(person, []int{31, 41, 51, 61, 71}, true)
			if len(shuncard) > 0 && !self.IsBoom(shuncard) && !self.GetBigCard(shuncard) {
				canstep = append(canstep, shuncard)
			}
			//! 试试连对
			lianduicard := self.GetBetterCard(person, []int{31, 31, 41, 41}, true)
			if len(lianduicard) > 0 && !self.IsBoom(lianduicard) && !self.GetBigCard(lianduicard) {
				canstep = append(canstep, lianduicard)
			}
			//! 试试三带二
			sandaiercard := self.GetBetterCard(person, []int{31, 31, 31, 41, 41}, true)
			if len(sandaiercard) > 0 && !self.IsBoom(sandaiercard) && !self.GetBigCard(sandaiercard) {
				canstep = append(canstep, sandaiercard)
			}
			//! 试试三带一
			sandaiyicard := self.GetBetterCard(person, []int{31, 31, 31, 41}, true)
			if len(sandaiyicard) > 0 && !self.IsBoom(sandaiyicard) && !self.GetBigCard(sandaiyicard) {
				canstep = append(canstep, sandaiyicard)
			}
			//! 试试对子
			duicard := self.GetBetterCard(person, []int{31, 31}, true)
			if len(duicard) > 0 && !self.IsBoom(duicard) && !self.GetBigCard(duicard) {
				canstep = append(canstep, duicard)
			}
			onecard := self.GetBetterCard(person, []int{31}, true)
			if len(onecard) > 0 && !self.IsBoom(onecard) && !self.GetBigCard(onecard) {
				canstep = append(canstep, onecard)
			}
			if len(canstep) > 0 {
				return canstep[lib.HF_GetRandom(len(canstep))]
			} else {
				return []int{person.Card[0]}
			}
		}
	}

	return self.GetBetterCard(person, self.LastCard, false)
}

//! 得到比这个牌大的牌
func (self *Game_GoldDDZ) GetBetterCard(person *Game_GoldDDZ_Person, lastcard []int, special bool) []int {
	stepcard := make([]int, 0)

	tmp := IsOkByCards(lastcard)

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

	if tmp%100 == TYPE_CARD_ONE { //! 单张
		for i := 0; i < len(one); i++ {
			if CardCompare(one[i], tmp/10) > 0 || special {
				stepcard = append(stepcard, self.GetHasCard(person.Card, 1, one[i])...)
				return stepcard
			}
		}
		sort.Sort(LstPoker(person.Card))
		for i := 0; i < len(person.Card); i++ {
			if CardCompare(person.Card[i], tmp/10) > 0 || special {
				stepcard = append(stepcard, person.Card[i])
				return stepcard
			}
		}
	}

	if tmp%100 == TYPE_CARD_TWO { //! 对子
		for i := 0; i < len(two); i++ {
			if CardCompare(two[i], tmp/10) > 0 || special {
				stepcard = append(stepcard, self.GetHasCard(person.Card, 2, two[i])...)
				return stepcard
			}
		}

		for i := 0; i < len(three); i++ {
			if CardCompare(three[i], tmp/10) > 0 || special {
				stepcard = append(stepcard, self.GetHasCard(person.Card, 2, three[i])...)
				return stepcard
			}
		}
	}

	if tmp%100 == TYPE_CARD_ZHA { //! 炸弹
		for i := 0; i < len(four); i++ {
			if CardCompare(four[i], tmp/10) > 0 || special {
				stepcard = append(stepcard, self.GetHasCard(person.Card, 4, four[i])...)
				return stepcard
			}
		}
	}

	if tmp%100 == TYPE_CARD_SAN { //! 3张
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找到牌型:3张")
		for i := 0; i < len(three); i++ {
			if CardCompare(three[i], tmp/10) > 0 || special {
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
		sort.Sort(LstPoker(lastcard))
		if len(_card) >= len(lastcard) {
			index := -1
			for i := 0; i < len(_card); i++ {
				if CardCompare(_card[i], lastcard[0]) > 0 || special {
					index = i
					break
				}
			}
			if index >= 0 {
				for {
					if index+len(lastcard)-1 >= len(_card) {
						break
					}
					find := true
					for i := 1; i <= len(lastcard)-1; i++ {
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
						for i := index; i < index+len(lastcard); i++ {
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
		sort.Sort(LstPoker(lastcard))
		if len(_card) >= len(lastcard)/2 {
			index := -1
			for i := 0; i < len(_card); i++ {
				if CardCompare(_card[i], lastcard[0]) > 0 || special {
					index = i
					break
				}
			}
			if index >= 0 {
				for {
					if index+len(lastcard)/2-1 >= len(_card) {
						break
					}
					find := true
					for i := 1; i <= len(lastcard)/2-1; i++ {
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
						for i := index; i < index+len(lastcard)/2; i++ {
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
		last := IsOkByGoldPDKCards(lastcard, make([]int, 0))
		if len(_card) >= len(lastcard)/3 {
			index := -1
			for i := 0; i < len(_card); i++ {
				if CardCompare(_card[i], last/10) > 0 || special {
					index = i
					break
				}
			}
			if index >= 0 {
				for {
					if index+len(lastcard)/3-1 >= len(_card) {
						break
					}
					find := true
					for i := 1; i <= len(lastcard)/3-1; i++ {
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
						for i := index; i < index+len(lastcard)/3; i++ {
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
			if CardCompare(four[i], tmp/10) > 0 || special {
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

//! 是否是炸弹
func (self *Game_GoldDDZ) IsBoom(card []int) bool {
	if len(card) != 4 {
		return false
	}
	return card[0]/10 == card[1]/10 && card[0]/10 == card[2]/10 && card[0]/10 == card[3]/10
}

//! 得到最大的牌
func (self *Game_GoldDDZ) GetMaxCard(card []int) int {
	maxcard := 0
	for i := 0; i < len(card); i++ {
		if maxcard == 0 {
			maxcard = card[i]
			continue
		}
		if CardCompare(card[i], maxcard) > 0 {
			maxcard = card[i]
		}
	}

	return maxcard
}

//! 有没有大牌
func (self *Game_GoldDDZ) GetBigCard(card []int) bool {
	for i := 0; i < len(card); i++ {
		if card[i]/10 == 12 || card[i]/10 == 13 || card[i]/10 == 1 || card[i]/10 == 2 || card[i] == 1000 || card[i] == 2000 || card[i]/10 == 14 || card[i]/10 == 20 {
			return true
		}
	}
	return false
}

func (self *Game_GoldDDZ) GetHasCard(card []int, num int, key int) []int {
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

func (self *Game_GoldDDZ_Person) IsTrush(_time int64) bool {
	if self.Trust {
		return true
	}

	if time.Now().Unix() < _time {
		return false
	}

	self.SetTrust(true)

	return true
}

//! 设置托管
func (self *Game_GoldDDZ_Person) SetTrust(trust bool) {
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
