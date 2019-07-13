package gameserver

import (
	"lib"
	"math"
	"sort"
	"staticfunc"
	"time"
)

//!
type Msg_GameDDZ_Info struct {
	Begin       bool               `json:"begin"` //! 是否开始
	Deal        []Son_GameDDZ_Deal `json:"deal"`  //! 抢地主的人
	Info        []Son_GameDDZ_Info `json:"info"`
	CurStep     int64              `json:"curstep"`     //! 当前谁出牌
	BefStep     int64              `json:"befstep"`     //! 上局谁出
	LastCard    []int              `json:"lastcard"`    //! 最后的牌
	LastAbsCard []int              `json:"lastabscard"` //! 实际出牌
	DZCard      []int              `json:"dzcard"`      //! 地主牌
	Bets        int                `json:"bets"`        //! 底分
	Boom        int                `json:"boom"`        //! 炸弹
	Razz        int                `json:"razz"`
}

type Son_GameDDZ_Info struct {
	Uid    int64 `json:"uid"`
	Card   []int `json:"card"`
	Dealer bool  `json:"dealer"`
	Score  int   `json:"score"`
	Total  int   `json:"total"`
	Ready  bool  `json:"ready"`
	Bets   int   `json:"bets"`
	Double int   `json:"isdouble"`
}
type Son_GameDDZ_Deal struct {
	Uid int64 `json:"uid"`
	Ok  bool  `json:"ok"`
}

//! 出牌
type Msg_GameDDZ_Step struct {
	Uid      int64 `json:"uid"`      //! 哪个uid
	Cards    []int `json:"cards"`    //! 出的啥牌
	AbsCards []int `json:"abscards"` //! 实际牌
	CurStep  int64 `json:"curstep"`  //! 下局谁出
}

//! 结算
type Msg_GameDDZ_End struct {
	Info []Son_GameDDZ_Info `json:"info"`
	Bets int                `json:"bets"`
	Boom int                `json:"boom"`
	CT   bool               `json:"ct"`
}

//! 房间结束
type Msg_GameDDZ_Bye struct {
	Info []Son_GameDDZ_Bye `json:"info"`
}
type Son_GameDDZ_Bye struct {
	Uid   int64 `json:"uid"`
	Win   int   `json:"win"`   //! 胜利次数
	Deal  int   `json:"deal"`  //! 失败次数
	High  int   `json:"high"`  //! 单局最高
	Boom  int   `json:"boom"`  //! 炸弹数量
	Score int   `json:"score"` //! 总分
}

//! 地主
type Msg_GameDDZ_Dealer struct {
	Uid  int64 `json:"uid"`
	Card []int `json:"card"`
	Bets int   `json:"bets"`
	Razz int   `json:"razz"`
}

type Msg_GameDDZ_Bets struct {
	Uid     int64 `json:"uid"`
	Bets    int   `json:"bets"`
	CurStep int64 `json:"curstep"`
}

type Msg_GameDDZ_Double struct {
	Double int `json:"double"`
}

///////////////////////////////////////////////////////
type Game_DDZ_Person struct {
	Uid      int64 `json:"uid"`
	Card     []int `json:"card"`     //! 手牌
	Bets     int   `json:"bets"`     //! 叫分
	Win      int   `json:"win"`      //! 胜利次数
	Deal     int   `json:"deal"`     //! 地主次数
	Double   int   `json:"isdouble"` //! 是否加倍
	Score    int   `json:"score"`    //! 积分
	Dealer   bool  `json:"dealer"`   //! 是否是地主
	CurScore int   `json:"curscore"` //! 当前局的分数
	Ready    bool  `json:"ready"`    //! 是否准备好
	Boom     int   `json:"boom"`     //! 炸弹数量
	High     int   `json:"high"`     //! 单局最高
}

func (self *Game_DDZ_Person) Init() {
	self.Card = make([]int, 0)
	self.Dealer = false
	self.Double = 0
}

type Game_DDZ struct {
	PersonMgr   []*Game_DDZ_Person           `json:"personmgr"`   //! 玩家信息
	DZCard      []int                        `json:"dzcard"`      //! 地主牌
	LastCard    []int                        `json:"lastcard"`    //! 最后出的牌
	LastAbsCard []int                        `json:"lastabscard"` //! 实际出牌
	CurStep     int64                        `json:"curstep"`     //! 谁出牌
	BefStep     int64                        `json:"befstep"`     //! 上局谁出
	Winer       int64                        `json:"winer"`       //! 赢家
	Bets        int                          `json:"bets"`        //! 底分
	Boom        int                          `json:"boom"`        //! 炸弹数量
	Taxi        int                          `json:"taxi"`        //! 农民出牌次数
	DZNum       int                          `json:"dznum"`       //! 地主出牌次数
	Razz        int                          `json:"razz"`        //! 癞子
	Record      *staticfunc.Rec_GameDDZ_Info `json:"record"`      //! 记录

	room *Room
	hz   bool
}

func NewGame_DDZ() *Game_DDZ {
	game := new(Game_DDZ)
	game.DZCard = make([]int, 0)
	game.LastCard = make([]int, 0)
	game.PersonMgr = make([]*Game_DDZ_Person, 0)

	return game
}

func (self *Game_DDZ) OnInit(room *Room) {
	self.room = room
}

func (self *Game_DDZ) OnRobot(robot *lib.Robot) {

}

func (self *Game_DDZ) OnSendInfo(person *Person) {
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "maxstep : ", self.room.MaxStep)
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			person.SendMsg("gameddzinfo", self.getInfo(person.Uid))
			return
		}
	}

	_person := new(Game_DDZ_Person)
	_person.Init()
	_person.Uid = person.Uid
	if self.room.Type == 75 || self.room.Type == 67 { //! 可以准备的斗地主
		_person.Ready = false
	} else {
		_person.Ready = true
	}
	self.PersonMgr = append(self.PersonMgr, _person)

	if len(self.PersonMgr) >= lib.HF_Atoi(self.room.csv["minnum"]) && self.room.Type != 75 && self.room.Type != 67 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
		self.OnBegin()
		return
	}

	person.SendMsg("gameddzinfo", self.getInfo(person.Uid))
}

func (self *Game_DDZ) OnMsg(msg *RoomMsg) {
	if self.room.Type == 67 {
		switch msg.Head {
		case "gamebets": //! 叫分
			self.GameBets(msg.Uid, msg.V.(*Msg_GameBets).Bets)
		case "gamesteps": //! 出牌
			self.GameStep(msg.Uid, msg.V.(*Msg_GameSteps).Cards, msg.V.(*Msg_GameSteps).AbsCards)
		case "gameready": //! 游戏准备
			self.GameReady(msg.Uid)
		}
	} else {
		switch msg.Head {
		case "gamebets": //! 叫分
			self.GameBets(msg.Uid, msg.V.(*Msg_GameBets).Bets)
		case "gamedouble": //! 双倍
			self.GameDouble(msg.Uid, msg.V.(*Msg_GameDDZ_Double).Double)
		case "gamesteps": //! 出牌
			self.GameStep(msg.Uid, msg.V.(*Msg_GameSteps).Cards, msg.V.(*Msg_GameSteps).AbsCards)
		case "gameready": //! 游戏准备
			self.GameReady(msg.Uid)
		}
	}

}

func (self *Game_DDZ) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)

	if self.Winer == 0 {
		self.Winer = self.PersonMgr[0].Uid
	}

	cardmgr := NewCard_DDZ()
	self.Record = new(staticfunc.Rec_GameDDZ_Info)
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
		if self.room.Type == 75 { //! 保定(自选庄)
			if self.PersonMgr[i].Dealer {
				self.CurStep = self.PersonMgr[i].Uid
			}
		}
		self.PersonMgr[i].Dealer = false
		self.PersonMgr[i].Bets = -1
		self.PersonMgr[i].Double = 0
		self.PersonMgr[i].Card = cardmgr.Deal(17)
		self.PersonMgr[i].CurScore = 0

		//! 记录
		var rc_person staticfunc.Son_Rec_GameDDZ_Person
		rc_person.Uid = self.PersonMgr[i].Uid
		rc_person.Name = self.room.GetName(rc_person.Uid)
		rc_person.Head = self.room.GetHead(rc_person.Uid)
		lib.HF_DeepCopy(&rc_person.Card, &self.PersonMgr[i].Card)
		self.Record.Person = append(self.Record.Person, rc_person)
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gameddzbegin", self.getInfo(person.Uid))
	}

	if self.room.Type == 75 { //! 保定经典地主
		if self.room.Param1/10%10 == 1 { //! 双王或四个2,拿地主
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Uid == self.CurStep {
					num := 0
					num2 := 0
					for _, card := range self.PersonMgr[i].Card {
						if card == 2000 || card == 1000 {
							num++
						}
						if card/10 == 2 {
							num2++
						}
					}
					if num >= 2 || num2 >= 4 {
						lib.GetLogMgr().Output(lib.LOG_DEBUG, "双王/四个2", num, num2)
						self.PersonMgr[i].Dealer = true
						self.PersonMgr[i].Deal++
						self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, self.DZCard...)

						if self.Record != nil {
							self.Record.Person[i].Card = append(self.Record.Person[i].Card, self.DZCard...)
						}

						self.Bets = 1

						var msg Msg_GameDDZ_Dealer
						msg.Uid = self.CurStep
						msg.Card = self.DZCard
						msg.Bets = self.Bets
						self.room.broadCastMsg("gamedealer", &msg)
					}
				}
			}
		}
	}
	self.room.flush()
}

//! 下注
func (self *Game_DDZ) GameBets(uid int64, bets int) {
	if !self.room.Begin { //! 未开始不能叫分
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "未开始不能抢庄")
		return
	}

	if self.room.Type == 67 && self.CurStep == 0 { //第一局开始房主叫地主
		self.CurStep = self.PersonMgr[0].Uid
		self.Winer = self.CurStep
	}

	if uid != self.CurStep {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不归你叫")
		return
	}
	if self.room.Type == 67 {
		if bets == 3 { //! 叫到最大分
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Uid == uid {
					self.PersonMgr[i].Dealer = true
					self.PersonMgr[i].Deal++
					self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, self.DZCard...)

					if self.Record != nil {
						self.Record.Person[i].Card = append(self.Record.Person[i].Card, self.DZCard...)
					}
					break
				}
			}
			self.Bets = 3
			self.CurStep = uid
			var msg Msg_GameDDZ_Dealer
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

			if num >= len(self.PersonMgr)-1 {
				if bets > max {
					max = bets
					muid = uid
				}

				if max == 0 { //! 都不叫庄
					for i := 0; i < len(self.PersonMgr); i++ {
						if self.PersonMgr[i].Uid == self.Winer {
							self.PersonMgr[i].Dealer = true
							self.PersonMgr[i].Deal++
							self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, self.DZCard...)

							if self.Record != nil {
								self.Record.Person[i].Card = append(self.Record.Person[i].Card, self.DZCard...)
							}
							break
						}
					}
					self.Bets = 1
					self.CurStep = self.Winer

					var msg Msg_GameDDZ_Dealer
					msg.Uid = muid
					msg.Card = self.DZCard
					msg.Bets = self.Bets
					msg.Razz = self.Razz
					self.room.broadCastMsg("gamedealer", &msg)
					return
				}

				for i := 0; i < len(self.PersonMgr); i++ {
					if self.PersonMgr[i].Uid == muid {
						self.PersonMgr[i].Dealer = true
						self.PersonMgr[i].Deal++
						self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, self.DZCard...)

						if self.Record != nil {
							self.Record.Person[i].Card = append(self.Record.Person[i].Card, self.DZCard...)
						}
						break
					}
				}

				self.Bets = max
				self.CurStep = muid

				var msg Msg_GameDDZ_Dealer
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

				var msg Msg_GameDDZ_Bets
				msg.Uid = uid
				msg.Bets = bets
				msg.CurStep = self.CurStep
				self.room.broadCastMsg("gamebets", &msg)
			}

		}

	} else if self.room.Type == 75 { //! 保定经典地主
		if bets == 0 { //! 不要地主

			self.CurStep = self.GetNextUid()

			num := 0
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Uid == uid {
					if self.PersonMgr[i].Bets == 0 {
						return
					} else {
						self.PersonMgr[i].Bets = 0
						num++
					}
				} else if self.PersonMgr[i].Bets == 0 {
					num++
				}
			}

			if num == len(self.PersonMgr) && num >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 都不叫庄
				self.hz = true
				self.OnEnd()
				return
			}

			var msg Msg_GameDDZ_Bets
			msg.Uid = uid
			msg.Bets = bets
			msg.CurStep = self.CurStep
			self.room.broadCastMsg("gamebets", &msg)

			if self.room.Param1/10%10 == 1 { //! 双王或四个2,必须拿地主
				p := self.GetPerson(self.CurStep)
				if nil == p {
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "玩家不在")
					return
				}
				num := 0
				num2 := 0
				for _, card := range p.Card {
					if card == 2000 || card == 1000 {
						num++
					}
					if card/10 == 2 {
						num2++
					}
				}
				if num >= 2 || num2 >= 4 {
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "双王/四个2")
					for i := 0; i < len(self.PersonMgr); i++ {
						if self.PersonMgr[i].Uid == self.CurStep {
							self.PersonMgr[i].Dealer = true
							self.PersonMgr[i].Deal++
							self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, self.DZCard...)

							if self.Record != nil {
								self.Record.Person[i].Card = append(self.Record.Person[i].Card, self.DZCard...)
							}
							break
						}
					}

					self.Bets = 1

					var msg Msg_GameDDZ_Dealer
					msg.Uid = self.CurStep
					msg.Card = self.DZCard
					msg.Bets = self.Bets
					self.room.broadCastMsg("gamedealer", &msg)
					return
				}
			}
		} else if bets == 1 { //! 要地主
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Uid == uid {
					self.PersonMgr[i].Dealer = true
					self.PersonMgr[i].Deal++
					self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, self.DZCard...)

					if self.Record != nil {
						self.Record.Person[i].Card = append(self.Record.Person[i].Card, self.DZCard...)
					}
					break
				}
			}

			self.Bets = 1
			self.CurStep = uid

			var msg Msg_GameDDZ_Dealer
			msg.Uid = uid
			msg.Card = self.DZCard
			msg.Bets = self.Bets
			self.room.broadCastMsg("gamedealer", &msg)
		} else {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, bets, ":参数错误")
			return
		}
	} else {
		if bets == 3 || (self.room.Param1/10%10 == 1 && bets > 0) || (self.room.Type == 8 && bets >= 2) { //!
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Uid == uid {
					self.PersonMgr[i].Dealer = true
					self.PersonMgr[i].Deal++
					self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, self.DZCard...)

					if self.Record != nil {
						self.Record.Person[i].Card = append(self.Record.Person[i].Card, self.DZCard...)
					}
					break
				}
			}

			if self.room.Param1/10%10 == 1 {
				self.Bets = 1
			} else {
				if self.room.Type == 8 {
					self.Bets = 2
				} else {
					self.Bets = 3
				}
			}
			self.CurStep = uid
			if self.room.Type == 8 {
				self.Razz = lib.HF_GetRandom(13) + 1

				if self.Record != nil {
					self.Record.Razz = self.Razz
				}
			}

			var msg Msg_GameDDZ_Dealer
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
						self.PersonMgr[i].Deal++
						self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, self.DZCard...)

						if self.Record != nil {
							self.Record.Person[i].Card = append(self.Record.Person[i].Card, self.DZCard...)
						}
						break
					}
				}

				self.Bets = lib.HF_MaxInt(max, 1)
				self.CurStep = muid
				if self.room.Type == 8 {
					self.Razz = lib.HF_GetRandom(13) + 1
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "癞子是:", self.Razz)
				}

				var msg Msg_GameDDZ_Dealer
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

				var msg Msg_GameDDZ_Bets
				msg.Uid = uid
				msg.Bets = bets
				msg.CurStep = self.CurStep
				self.room.broadCastMsg("gamebets", &msg)
			}
		}

	}

	self.room.flush()
}

//! 加倍
func (self *Game_DDZ) GameDouble(uid int64, double int) {
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

	self.room.flush()
}

//! 出牌(玩家选择)
func (self *Game_DDZ) GameStep(uid int64, card []int, abscard []int) {
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

	if self.room.Type != 8 { //! 普通模式
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

	if self.room.Type != 67 {
		if self.room.Param1/100%10 == 1 { //! 不可拆王
			num := 0
			wang := 0
			for i := 0; i < len(card); i++ {
				if card[i] == 1000 || card[i] == 2000 {
					num++
					wang = card[i]
				}
			}

			if num == 1 { //! 有一个王
				for i := 0; i < len(person.Card); i++ {
					if person.Card[i] != wang && (person.Card[i] == 1000 || person.Card[i] == 2000) {
						lib.GetLogMgr().Output(lib.LOG_DEBUG, "不可拆王:", card)
						return
					}
				}
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
		if self.room.Type == 75 && tmp%100 == 7 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌错误:", card)
			return
		}
		if self.room.Type == 67 && tmp%10 == TYPE_CARD_SHUN && self.room.Param1/1000%10 == 1 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能出顺子")
			return
		}
		if tmp == 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌错误:", card)
			return
		}
		if len(self.LastCard) != 0 && self.BefStep != uid {
			_tmp := IsOkByCards(self.LastAbsCard)
			if tmp == TYPE_CARD_WANG { //! 王炸
				self.Boom++
				person.Boom++
			} else if tmp%100 == TYPE_CARD_ZHA { //! 炸弹
				if _tmp%100 == TYPE_CARD_ZHA {
					compareboom := true //! 是否比炸弹
					if self.room.Type == 8 {
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
				person.Boom++
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
				person.Boom++
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

	if self.Record != nil {
		self.Record.Step = append(self.Record.Step, staticfunc.Son_Rec_GameDDZ_Step{uid, card, abscard})
	}

	var msg Msg_GameDDZ_Step
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
func (self *Game_DDZ) GameReady(uid int64) {
	if self.room.IsBye() {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "roomisbye")
		return
	}

	if self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "roomisbegin")
		return
	}

	if self.room.Type != 75 && self.room.Type != 67 && self.room.Step == 0 { //! 一局之后才有这个消息
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "room.step == 0")
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
func (self *Game_DDZ) OnEnd() {
	self.room.SetBegin(false)

	if self.hz {
		var msg Msg_GameDDZ_End
		msg.Bets = 0
		msg.Boom = 0
		msg.CT = false
		if self.room.Type == 75 {
			if self.room.Param1/100%10 == 1 { //! 保定比分模式
				king := self.GetKing()
				for i := 0; i < len(self.PersonMgr); i++ {
					var son Son_GameDDZ_Info
					son.Uid = self.PersonMgr[i].Uid
					son.Card = self.PersonMgr[i].Card
					son.Dealer = self.PersonMgr[i].Dealer
					if king == son.Uid {
						self.PersonMgr[i].CurScore = -2
					} else {
						self.PersonMgr[i].CurScore = 1
					}
					son.Score = self.PersonMgr[i].CurScore
					self.PersonMgr[i].Score += self.PersonMgr[i].CurScore
					son.Total = self.PersonMgr[i].Score
					son.Double = self.PersonMgr[i].Double
					msg.Info = append(msg.Info, son)
				}
			} else { //! 不比分
				for i := 0; i < len(self.PersonMgr); i++ {
					var son Son_GameDDZ_Info
					son.Uid = self.PersonMgr[i].Uid
					son.Card = self.PersonMgr[i].Card
					son.Dealer = self.PersonMgr[i].Dealer
					if self.CurStep == self.PersonMgr[i].Uid {
						self.PersonMgr[i].CurScore = -2
					} else {
						self.PersonMgr[i].CurScore = 1
					}
					son.Score = self.PersonMgr[i].CurScore
					self.PersonMgr[i].Score += self.PersonMgr[i].CurScore
					son.Total = self.PersonMgr[i].Score
					son.Double = self.PersonMgr[i].Double
					msg.Info = append(msg.Info, son)
				}
			}
		} else {
			for i := 0; i < len(self.PersonMgr); i++ {
				var son Son_GameDDZ_Info
				son.Uid = self.PersonMgr[i].Uid
				son.Card = self.PersonMgr[i].Card
				son.Dealer = self.PersonMgr[i].Dealer
				son.Score = self.PersonMgr[i].CurScore
				son.Total = self.PersonMgr[i].Score
				son.Double = self.PersonMgr[i].Double
				msg.Info = append(msg.Info, son)
			}
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
	if self.room.Type == 8 { //! 癞子
		maxboom = 5
		if self.room.Param1%10 == 1 {
			addbet = true
		}
	} else if self.room.Type == 67 {
		maxboom = 4
	} else { //! 经典
		maxboom = self.room.Param1%10 + 3
		addbet = true
	}

	_boom := self.Boom - maxboom
	if _boom > 0 {
		self.Boom = maxboom
	}
	score := self.Bets * int(math.Pow(2.0, float64(self.Boom)))
	if DZWin && self.Taxi == 0 {
		score *= 2
	} else if !DZWin && self.DZNum <= 1 {
		score *= 2
	}
	if _boom > 0 && addbet {
		score += self.Bets * _boom
	}

	//! 先找到地主
	var dealer *Game_DDZ_Person
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

	if self.room.Type == 75 {
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Dealer {
				continue
			}
			if dealer.Double == 2 && self.PersonMgr[i].Double == 1 { //! 农民不踹
				self.PersonMgr[i].CurScore /= 2
				dealer.CurScore += self.PersonMgr[i].CurScore
			}

			//! 触顶，重置玩家分数
			if self.room.Param1%10 == 1 { //! 50封顶
				if self.PersonMgr[i].CurScore < (-50) {
					dealer.CurScore -= (-50) - self.PersonMgr[i].CurScore
					self.PersonMgr[i].CurScore = -50
				}
				if dealer.CurScore < (-100) {
					bs := float32(self.PersonMgr[i].CurScore) / float32(-dealer.CurScore)
					self.PersonMgr[i].CurScore = int(bs * float32(100))
				}

			} else if self.room.Param1%10 == 2 { //! 100封顶
				if self.PersonMgr[i].CurScore < (-100) {
					dealer.CurScore -= (-100) - self.PersonMgr[i].CurScore
					self.PersonMgr[i].CurScore = -100
				}
				if dealer.CurScore < (-200) {
					bs := float32(self.PersonMgr[i].CurScore) / float32(-dealer.CurScore)
					self.PersonMgr[i].CurScore = int(bs * float32(200))
				}
			}
		}

		//! 触顶，重置地主分数
		if self.room.Param1%10 == 1 && dealer.CurScore < (-100) {
			dealer.CurScore = -100
		}
		if self.room.Param1%10 == 2 && dealer.CurScore < (-200) {
			dealer.CurScore = -200
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Score += self.PersonMgr[i].CurScore
		if self.PersonMgr[i].CurScore > 0 {
			self.PersonMgr[i].Win++
			if self.PersonMgr[i].CurScore > self.PersonMgr[i].High {
				self.PersonMgr[i].High = self.PersonMgr[i].CurScore
			}
		}
	}

	var msg Msg_GameDDZ_End
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
		var son Son_GameDDZ_Info
		son.Uid = self.PersonMgr[i].Uid
		son.Card = self.PersonMgr[i].Card
		son.Dealer = self.PersonMgr[i].Dealer
		son.Score = self.PersonMgr[i].CurScore
		son.Total = self.PersonMgr[i].Score
		son.Double = self.PersonMgr[i].Double
		msg.Info = append(msg.Info, son)
		agentinfo = append(agentinfo, staticfunc.JS_CreateRoomMem{son.Uid, "", "", son.Total})

		if self.Record != nil {
			for j := 0; j < len(self.Record.Person); j++ {
				if self.Record.Person[j].Uid == self.PersonMgr[i].Uid {
					self.Record.Person[j].Score = self.PersonMgr[i].CurScore
					self.Record.Person[j].Total = self.PersonMgr[i].Score
					break
				}
			}
		}
	}
	self.room.broadCastMsg("gameddzend", &msg)
	self.room.AgentResult(agentinfo)

	if self.Record != nil {
		self.Record.Roomid = self.room.Id*100 + self.room.Step
		self.Record.Time = time.Now().Unix()
		self.Record.Razz = self.Razz
		self.Record.MaxStep = self.room.MaxStep
		self.Record.Type = self.room.Type
		self.room.AddRecord(lib.HF_JtoA(self.Record))
	}

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

func (self *Game_DDZ) OnBye() {
	info := make([]staticfunc.JS_CreateRoomMem, 0)
	var msg Msg_GameDDZ_Bye
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameDDZ_Bye
		son.Uid = self.PersonMgr[i].Uid
		son.Win = self.PersonMgr[i].Win
		son.Deal = self.PersonMgr[i].Deal
		son.High = self.PersonMgr[i].High
		son.Boom = self.PersonMgr[i].Boom
		son.Score = self.PersonMgr[i].Score
		msg.Info = append(msg.Info, son)
		info = append(info, staticfunc.JS_CreateRoomMem{son.Uid, "", "", son.Score})

		GetServer().SqlScoreLog(self.PersonMgr[i].Uid, self.room.GetName(self.PersonMgr[i].Uid), self.room.GetHead(self.PersonMgr[i].Uid), self.room.Type, self.room.Id, self.PersonMgr[i].Score)
	}
	self.room.broadCastMsg("gameddzbye", &msg)

	self.room.ClubResult(info)
}

func (self *Game_DDZ) OnExit(uid int64) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			copy(self.PersonMgr[i:], self.PersonMgr[i+1:])
			self.PersonMgr = self.PersonMgr[:len(self.PersonMgr)-1]
			break
		}
	}
}

func (self *Game_DDZ) getInfo(uid int64) *Msg_GameDDZ_Info {
	var msg Msg_GameDDZ_Info
	msg.Begin = self.room.Begin
	msg.Deal = make([]Son_GameDDZ_Deal, 0)
	msg.Info = make([]Son_GameDDZ_Info, 0)
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
		var son Son_GameDDZ_Info
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
		son.Score = value.CurScore
		son.Total = value.Score
		son.Ready = value.Ready
		son.Double = value.Double
		msg.Info = append(msg.Info, son)
	}
	return &msg
}

func (self *Game_DDZ) GetPerson(uid int64) *Game_DDZ_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

//! 得到下一个uid
func (self *Game_DDZ) GetNextUid() int64 {
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

func (self *Game_DDZ) OnTime() {

}

func (self *Game_DDZ) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_DDZ) OnIsBets(uid int64) bool {
	return false
}

func (self *Game_DDZ) GetKing() int64 {
	score := make(map[int64]int)
	sor := make([]int, 0)
	king := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		sor = append(sor, 0)
		for _, card := range self.PersonMgr[i].Card {
			if card == 2000 { //! 大王4分
				score[self.PersonMgr[i].Uid] += 4
				sor[i] += 4
				king = int(self.PersonMgr[i].Uid)
			} else if card == 1000 { //! 小王3分
				score[self.PersonMgr[i].Uid] += 3
				sor[i] += 3
			} else if card/10 == 2 { //! 2是2分
				score[self.PersonMgr[i].Uid] += 2
				sor[i] += 2
			}
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sor))) //! 排序
	if sor[0] == sor[1] {
		for key, val := range score {
			if val == sor[0] && key == int64(king) {
				return key
			}
		}
	} else {
		for key, val := range score {
			if val == sor[0] {
				return key
			}
		}
	}

	return 0
}

//! 结算所有人
func (self *Game_DDZ) OnBalance() {
}
