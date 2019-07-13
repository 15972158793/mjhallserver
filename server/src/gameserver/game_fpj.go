package gameserver

import (
	"lib"
	"math"
	"staticfunc"
	"time"
)

var FPJ_SYSTEM_MINBET []int64 = []int64{0, 0, 0}
var FPJ_SYSTEM_MAXBET []int64 = []int64{500000, 500000, 500000}

//!　散牌　10以上对子　2对　3条　顺子　同花　葫芦　四条　同花顺　同花大顺　五条
var FPJ_PL []int = []int{0, 1, 2, 3, 5, 5, 10, 20, 30, 50, 100}

type Rec_FPJ_Info struct {
	GameType int                  `json:"gametype"`
	Time     int64                `json:"time"` //! 记录时间
	Info     []Son_Rec_FPJ_Person `json:"info"`
}
type Son_Rec_FPJ_Person struct {
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`
	Head   string `json:"head"`
	Score  int    `json:"score"`
	Bets   int    `json:"bets"`
	Result []int  `json:"result"`
}

type Game_FPJ struct {
	Person    *Game_FPJ_Person `json:"person"`
	Result    []int            `json:"result"` //! 结果
	Bet       int              `json:"bet"`    //! 底分
	DealWin   int              `json:"dealwin"`
	CurBet    int              `json:"curbet"`    //! 当前可得分
	IsDouble  bool             `json:"isdouble"`  //! 是否是翻倍阶段
	IsHuan    bool             `json:"ishuan"`    //! 是否可以换牌
	DoubleNum int              `json:"doublenum"` //! 加倍次数
	Time      int64            `json:"time"`
	room      *Room
}

func NewGame_FPJ() *Game_FPJ {
	game := new(Game_FPJ)
	game.Result = make([]int, 0)
	game.Bet = 100
	game.CurBet = 0
	game.IsDouble = false
	game.IsHuan = false
	game.DoubleNum = 1
	return game
}

type Game_FPJ_Person struct {
	Uid     int64  `json:"uid"`
	Gold    int    `json:"gold"`
	Total   int    `json:"total"`
	Win     int    `json:"win"`  //! 赢钱
	Cost    int    `json:"cost"` //!  抽水
	Address string `json:"address"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	IP      string `json:"ip"`
	Sex     int    `json:"sex"`
}

type Msg_GameFPJ_Result struct {
	Result    []int `json:"result"`    //! 结果
	CardType  int   `json:"cardtype"`  //! 牌型 10-五条 9-同花大顺 8-同花顺 7-四条 6-葫芦 5-同花 4-顺子 3-3条 2-两对 1-1对10以上 0-散牌
	Bet       int   `json:"bet"`       //! 押分
	CurBet    int   `json:"curbet"`    //! 可赢多少分
	IsDouble  bool  `json:"isdouble"`  //! 是否可以加倍
	DoubleNum int   `json:"doublenum"` //! 加倍次数
	IsHuan    bool  `json:"ishuan"`    //! 是否可以换牌
}

type Msg_GameFPJ_Double struct {
	Card      int `json:"card"`      //!　开的什么牌
	Tmp       int `json:"tmp"`       //! 玩家猜的什么 1-小 2-大
	CurBet    int `json:"curbet"`    //! 当前可以赢多少
	DoubleNum int `json:"doublenum"` //! 翻倍次数
}

type Msg_GameFPJ_Info struct {
	Begin     bool             `json:"begin"`
	Result    []int            `json:"result"`
	Bet       int              `json:"bet"`
	CurBet    int              `json:"curbet"`
	IsDouble  bool             `json:"isdouble"`
	DoubleNum int              `json:"doublenum"` //! 加倍次数
	IsHuan    bool             `json:"ishuan"`
	Person    Son_GameFPJ_Info `json:"person"`
}

type Son_GameFPJ_Info struct {
	Uid     int64  `json:"uid"`
	Total   int    `json:"total"`
	Win     int    `json:"win"` //! 赢钱
	Address string `json:"address"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	IP      string `json:"ip"`
	Sex     int    `json:"sex"`
}

type Msg_GameFPJ_Huan struct {
	Uid   int64 `json:"uid"`
	Bet   int   `json:"bet"`
	HCard []int `json:"hcard"`
}

func (self *Game_FPJ) getinfo(uid int64) *Msg_GameFPJ_Info {
	var msg Msg_GameFPJ_Info
	msg.Begin = self.room.Begin
	msg.Result = self.Result
	msg.Bet = self.Bet
	msg.CurBet = self.CurBet
	msg.IsDouble = self.IsDouble
	msg.DoubleNum = self.DoubleNum
	msg.IsHuan = self.IsHuan
	if self.Person != nil && self.Person.Uid == uid {
		msg.Person.Uid = uid
		msg.Person.Total = self.Person.Total
		msg.Person.Win = self.Person.Win
		msg.Person.Address = self.Person.Address
		msg.Person.Name = self.Person.Name
		msg.Person.Head = self.Person.Head
		msg.Person.IP = self.Person.IP
		msg.Person.Sex = self.Person.Sex
	}
	return &msg
}

//! 同步金币
func (self *Game_FPJ_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

//!　同步总分
func (self *Game_FPJ) SendTotal(uid int64, total int) {
	var msg Msg_GameDFW_Total
	msg.Uid = uid
	msg.Total = total
	self.room.SendMsg(uid, "gametotal", &msg)
}

func (self *Game_FPJ) OnIsBets(uid int64) bool {
	if self.room.Begin {
		return true
	} else {
		return false
	}
}

func (self *Game_FPJ) GameBet(uid int64, bet int) {
	if self.Person == nil || self.Person.Uid != uid {
		return
	}

	if bet < 100 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "bet少于100 bet : ", bet)
		return
	}

	if self.Person.Total < bet {
		self.room.SendErr(uid, "您的金币不足，请前往充值")
		return
	}
	if self.room.Begin {
		return
	}

	self.Bet = bet
	self.DealWin += bet
	self.Person.Total -= bet
	var msg Msg_GameDFW_Total
	msg.Uid = uid
	msg.Total = self.Person.Total
	self.room.SendMsg(uid, "gamefpjbet", &msg)

	self.OnBegin()
}

func (self *Game_FPJ) GameHuan(uid int64, bet int, hcard []int) {

	if len(hcard) == 0 {
		self.OnEnd()
		return
	}

	if !self.IsHuan {
		return
	}

	if self.Person == nil || self.Person.Uid != uid {
		return
	}

	if bet != self.Bet {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "bet != self.bet bet : ", bet, " self.bet : ", self.Bet)
		return
	}

	if self.Person.Total < bet {
		self.room.SendErr(uid, "您的金币不足，请前往充值")
		return
	}

	if len(hcard) == 0 || len(self.Result) == 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------- len(hcard) : ", len(hcard), "  len(self.result) : ", len(self.Result))
		return
	}

	self.IsHuan = false

	for i := 0; i < len(hcard); i++ {
		find := false
		for j := 0; j < len(self.Result); j++ {
			if hcard[i] == self.Result[j] {
				find = true
				break
			}
		}
		if !find {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, " hcard[i] : ", hcard[i], " self.result : ", self.Result)
			return
		}
	}

	self.Bet = bet
	self.DealWin += bet
	self.Person.Total -= bet
	var msg Msg_GameDFW_Total
	msg.Uid = uid
	msg.Total = self.Person.Total
	self.room.SendMsg(uid, "gamefpjbet", &msg)

	lst := make([][]int, 0)
	winLst := make([][]int, 0)
	lostLst := make([][]int, 0)

	for k := 0; k < 100; k++ {
		cardmgr := NewCard_DDZ()
		for i := 0; i < len(self.Result); i++ {
			cardmgr.DealCard(self.Result[i])
		}
		card := make([]int, 0)
		lib.HF_DeepCopy(&card, &self.Result)

		for i := 0; i < len(hcard); i++ {
			for j := 0; j < len(card); j++ {
				if hcard[i] == card[j] {
					_card := cardmgr.Deal(1)
					card[j] = _card[0]
				}
			}
		}

		cardtype := GetFPJType(card)
		win := self.DealWin - self.Bet*FPJ_PL[cardtype]
		if GetServer().FpjSysMoney[self.room.Type%180000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && GetServer().FpjSysMoney[self.room.Type%180000]+int64(win) <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
			lst = append(lst, card)
		}
		if win >= 0 {
			winLst = append(winLst, card)
		} else {
			lostLst = append(lostLst, card)
		}

		if len(winLst) >= 10 && len(lostLst) >= 10 {
			break
		}
	}

	if len(lst) == 0 {
		if GetServer().FpjSysMoney[self.room.Type%180000]+int64(self.DealWin) < lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && len(winLst) != 0 { //! 一定赢
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "----------换牌  庄家必赢")
			self.Result = winLst[lib.HF_GetRandom(len(winLst))]
		} else if GetServer().FpjSysMoney[self.room.Type%180000]+int64(self.DealWin) > lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax && len(lostLst) != 0 { //! 一定输
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "----------换牌  庄家必输")
			self.Result = lostLst[lib.HF_GetRandom(len(lostLst))]
		} else {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "----------换牌  没有合适选项2 纯随机")

			cardmgr := NewCard_DDZ()
			for i := 0; i < len(hcard); i++ {
				cardmgr.DealCard(hcard[i])
			}
			for i := 0; i < len(hcard); i++ {
				for j := 0; j < len(self.Result); j++ {
					if hcard[i] == self.Result[j] {
						_card := cardmgr.Deal(1)
						self.Result[j] = _card[0]
					}
				}
			}
		}
	} else {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "----------换牌  len(lst)")
		self.Result = lst[lib.HF_GetRandom(len(lst))]
	}

	var result Msg_GameFPJ_Result
	result.Result = self.Result
	result.Bet = self.Bet
	result.CardType = GetFPJType(self.Result)
	self.CurBet = self.Bet * FPJ_PL[result.CardType]
	result.CurBet = self.Bet * FPJ_PL[result.CardType]
	if result.CurBet > 0 {
		self.IsDouble = true
	}
	result.IsDouble = self.IsDouble
	result.IsHuan = false
	result.DoubleNum = self.DoubleNum
	self.room.broadCastMsg("gamefpjhuan", &result)

	if !self.IsDouble {
		self.OnEnd()
	}
	self.Time = time.Now().Unix() + 3600
}

func (self *Game_FPJ) OnBegin() {
	if self.room.Begin {
		return
	}

	if self.Person == nil {
		return
	}
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
	self.room.Begin = true

	self.Result = make([]int, 0)

	lst := make([][]int, 0)
	winLst := make([][]int, 0)
	lostLst := make([][]int, 0)
	bjLst := make([][]int, 0)

	for i := 0; i < 100; i++ {
		cardmgr := NewCard_DDZ()
		card := cardmgr.Deal(5)
		cardtype := GetFPJType(card)
		win := self.DealWin - self.Bet*FPJ_PL[cardtype]
		if GetServer().FpjSysMoney[self.room.Type%180000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && GetServer().FpjSysMoney[self.room.Type%180000]+int64(win) <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
			lst = append(lst, card)
			if win == 0 {
				bjLst = append(bjLst, card)
			}
		}
		if win >= 0 {
			winLst = append(winLst, card)
		} else {
			lostLst = append(lostLst, card)
		}

		if len(winLst) >= 10 && len(lostLst) >= 10 {
			break
		}
	}

	if len(lst) == 0 {
		if GetServer().FpjSysMoney[self.room.Type%180000]+int64(self.DealWin) < lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && len(winLst) != 0 { //! 一定赢
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "----------发牌 庄必赢")
			self.Result = winLst[lib.HF_GetRandom(len(winLst))]
		} else if GetServer().FpjSysMoney[self.room.Type%180000]+int64(self.DealWin) > lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax && len(lostLst) != 0 { //! 一定输
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "----------发牌 庄必输")
			self.Result = lostLst[lib.HF_GetRandom(len(lostLst))]
		} else {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "----------发牌 没有合适选项2 纯随机")
			cardmgr := NewCard_DDZ()
			self.Result = cardmgr.Deal(5)
		}
	} else {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "----------发牌 len(lst)  len(lst)", len(lst), " len(bjlst) : ", len(bjLst))
		self.Result = lst[lib.HF_GetRandom(len(lst))]
		cardtype := GetFPJType(self.Result)
		win := self.Bet * FPJ_PL[cardtype]
		if win < 100 {
			if lib.HF_GetRandom(100) < 30 && len(bjLst) != 0 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "----------发牌 赢本金")
				self.Result = bjLst[lib.HF_GetRandom(len(bjLst))]
			}
		}
	}

	//self.Result = []int{91, 81, 82, 83, 84}

	var msg Msg_GameFPJ_Result
	msg.Result = self.Result
	msg.Bet = self.Bet
	msg.CardType = GetFPJType(self.Result)
	self.CurBet = self.Bet * FPJ_PL[msg.CardType]
	msg.CurBet = self.CurBet
	if msg.CurBet > 0 {
		self.IsDouble = true
		self.IsHuan = false
	} else {
		self.IsDouble = false
		self.IsHuan = true
	}
	msg.IsDouble = self.IsDouble
	msg.IsHuan = self.IsHuan
	msg.DoubleNum = self.DoubleNum
	self.room.broadCastMsg("gamefpjbegin", &msg)

	self.Time = time.Now().Unix() + 3600
}

func (self *Game_FPJ) GameDouble(uid int64, tmp int) { //! 1-小 2-大 其他-放弃加倍
	if self.Person == nil {
		return
	}

	if !self.IsDouble {
		return
	}

	if tmp != 1 && tmp != 2 {
		self.OnEnd()
		return
	}

	if self.DoubleNum <= 0 {
		self.room.SendErr(uid, "翻倍次数不足，无法翻倍")
		return
	}
	self.IsDouble = false

	card := make([]int, 0)
	if GetServer().FpjSysMoney[self.room.Type%180000]+int64(self.CurBet*2) <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax && GetServer().FpjSysMoney[self.room.Type%180000]+int64(self.CurBet*2) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
		_tmp := 2
		cardmgr := NewCard_DDZ()
		card = cardmgr.Deal(1)
		if card[0] == 1000 || card[0] == 2000 {
			self.CurBet *= 2
		} else if card[0]/10 <= 8 && card[0]/10 != 1 {
			_tmp = 1
			if _tmp == tmp {
				self.CurBet *= 2
			} else {
				self.CurBet = 0
			}
		} else {
			if _tmp == tmp {
				self.CurBet *= 2
			} else {
				self.CurBet = 0
			}
		}
	} else if GetServer().FpjSysMoney[self.room.Type%180000]+int64(self.CurBet*2) > lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax { //! 奖池超过最大值，必猜中
		if tmp == 1 { //! 开小牌 2-8
			card = append(card, (lib.HF_GetRandom(7)+2)*10+lib.HF_GetRandom(4)+1)
		} else { //! 开大牌 9-13
			card = append(card, (lib.HF_GetRandom(5)+9)*10+lib.HF_GetRandom(4)+1)
		}
		self.CurBet *= 2
	} else if GetServer().FpjSysMoney[self.room.Type%180000]+int64(self.CurBet*2) < lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin { //!　奖池低与最小值，必猜错
		if tmp == 2 { //! 开小牌 2-8
			card = append(card, (lib.HF_GetRandom(7)+2)*10+lib.HF_GetRandom(4)+1)
		} else { //! 开大牌 9-13
			card = append(card, (lib.HF_GetRandom(5)+9)*10+lib.HF_GetRandom(4)+1)
		}
		self.CurBet = 0
	}

	self.DoubleNum--
	var msg Msg_GameFPJ_Double
	msg.Card = card[0]
	msg.Tmp = tmp
	msg.CurBet = self.CurBet
	msg.DoubleNum = self.DoubleNum
	self.room.SendMsg(uid, "gamefpjdouble", &msg)

	if self.DoubleNum == 0 || self.CurBet == 0 {
		self.OnEnd()
	}
	self.Time = time.Now().Unix() + 60
}

func (self *Game_FPJ) OnEnd() {
	if !self.room.Begin {
		return
	}

	self.room.Begin = false
	var msg Son_GameFPJ_Info
	msg.Uid = self.Person.Uid
	msg.Sex = self.Person.Sex
	msg.Name = self.Person.Name
	msg.IP = self.Person.IP
	msg.Head = self.Person.Head
	msg.Address = self.Person.Address
	if self.CurBet > 0 {
		self.Person.Win = self.CurBet
		if self.Person.Win-self.DealWin > 0 {
			self.Person.Cost = int(math.Ceil(float64(self.Person.Win-self.DealWin) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
			self.Person.Win -= self.Person.Cost
			GetServer().SqlAgentGoldLog(self.Person.Uid, self.Person.Cost, self.room.Type)
		}

		self.Person.Total += self.Person.Win
	}

	dealwin := self.DealWin - self.CurBet
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "------- dealwin : ", dealwin, " self.dealwin : ", self.DealWin, " self.curbet : ", self.CurBet)
	if dealwin != 0 {
		GetServer().SqlBZWLog(&SQL_BZWLog{1, dealwin, time.Now().Unix(), self.room.Type})

		if dealwin > 0 {
			cost := 0
			cost = int(math.Ceil(float64(dealwin) * float64(lib.GetManyMgr().GetProperty(self.room.Type).DealCost) / 100.0))
			dealwin -= cost
		}
		GetServer().SetFpjSysMoney(self.room.Type%180000, GetServer().FpjSysMoney[self.room.Type%180000]+int64(dealwin))
	}

	msg.Win = self.Person.Win
	msg.Total = self.Person.Total
	self.room.SendMsg(self.Person.Uid, "gamefpjend", &msg)
	var record Rec_FPJ_Info
	record.GameType = self.room.Type
	record.Time = time.Now().Unix()
	var rec Son_Rec_FPJ_Person
	rec.Uid = self.Person.Uid
	rec.Head = self.Person.Head
	rec.Name = self.Person.Name
	rec.Result = self.Result
	rec.Score = self.Person.Win - self.DealWin
	record.Info = append(record.Info, rec)
	GetServer().InsertRecord(self.room.Type, self.Person.Uid, lib.HF_JtoA(&record), rec.Score)

	//! 初始化
	self.IsHuan = true
	if self.Person.Win > 0 {
		self.IsHuan = false
	}
	self.CurBet = 0
	self.DoubleNum = 1
	self.IsDouble = false
	self.Person.Cost = 0
	self.Person.Win = 0
	self.DealWin = 0
	self.Result = make([]int, 0)

	self.Time = time.Now().Unix() + 60
}

func (self *Game_FPJ) OnInit(room *Room) {
	self.room = room
}

func (self *Game_FPJ) OnRobot(robot *lib.Robot) {

}

func (self *Game_FPJ) OnSendInfo(person *Person) {
	if self.Person != nil && self.Person.Uid == person.Uid {
		self.Person.SynchroGold(person.Gold)
		self.room.broadCastMsg("gamefpjinfo", self.getinfo(person.Uid))
		return
	}

	_person := new(Game_FPJ_Person)
	_person.Uid = person.Uid
	_person.Total = person.Gold
	_person.Sex = person.Sex
	_person.Name = person.Name
	_person.IP = person.ip
	_person.Head = person.Imgurl
	_person.Gold = person.Gold
	_person.Address = person.minfo.Address
	_person.Cost = 0
	_person.Win = 0
	self.Person = _person
	self.room.broadCastMsg("gamefpjinfo", self.getinfo(person.Uid))
	self.Time = time.Now().Unix() + 60
}

func (self *Game_FPJ) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold": //! 同步金币
		if self.Person.Uid == msg.V.(*staticfunc.Msg_SynchroGold).Uid {
			self.Person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal(self.Person.Uid, self.Person.Total)
		}
	case "gamefpjhuan":
		self.GameHuan(msg.Uid, msg.V.(*Msg_GameFPJ_Huan).Bet, msg.V.(*Msg_GameFPJ_Huan).HCard)
	case "gamebets":
		self.GameBet(msg.Uid, msg.V.(*Msg_GameBets).Bets)
	case "gameplay":
		self.GameDouble(msg.Uid, msg.V.(*Msg_GamePlay).Type)

	}
}

func (self *Game_FPJ) OnBye() {

}

func (self *Game_FPJ) OnExit(uid int64) {
	if self.Person == nil {
		return
	}
	if uid == self.Person.Uid {
		gold := self.Person.Total - self.Person.Gold
		if gold > 0 {
			GetRoomMgr().AddCard(self.Person.Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
		} else if gold < 0 {
			GetRoomMgr().CostCard(self.Person.Uid, staticfunc.TYPE_GOLD, -gold, self.room)
		}
		self.Person.Gold = self.Person.Total
	}
}

func (self *Game_FPJ) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_FPJ) OnBalance() {
	if self.Person == nil {
		return
	}

	gold := self.Person.Total - self.Person.Gold
	if gold > 0 {
		GetRoomMgr().AddCard(self.Person.Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
	} else if gold < 0 {
		GetRoomMgr().CostCard(self.Person.Uid, staticfunc.TYPE_GOLD, -gold, self.room)
	}
	self.Person.Gold = self.Person.Total
}

func (self *Game_FPJ) OnTime() {
	if self.Person == nil {
		return
	}
	if self.Time == 0 {
		return
	}
	if time.Now().Unix() >= self.Time {
		if self.room.Begin {
			self.OnEnd()
		}
		self.room.KickViewByUid(self.Person.Uid, 96)
	}

}
