package gameserver

import (
	"fmt"
	"lib"
	"math"
	"staticfunc"
	"time"
)

type Msg_GameGoldTTZ_Info struct {
	Begin bool                   `json:"begin"` //! 是否开始
	Info  []Son_GameGoldTTZ_Info `json:"info"`
	State int                    `json:"state"`
	Time  int64                  `json:"time"` //! 倒计时
	Card  []int                  `json:"card"`
	SZ    []int                  `json:"sz"`
	Bets  int                    `json:"bets"`
}

type Son_GameGoldTTZ_Info struct {
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
	Open    bool  `json:"open"`
}

type Son_GameTTZ_Info struct {
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
	Open    bool  `json:"open"`
}

//! 结算
type Msg_GameTTZ_End struct {
	Cache []int              `json:"cache"`
	Info  []Son_GameTTZ_Info `json:"info"`
}

///////////////////////////////////////////////////////
type Game_GoldTTZ_Person struct {
	Uid      int64 `json:"uid"`
	Card     []int `json:"card"`     //! 手牌
	Ready    bool  `json:"ready"`    //! 是否准备
	Score    int   `json:"score"`    //! 积分
	Dealer   bool  `json:"dealer"`   //! 是否庄家
	RobDeal  int   `json:"robdeal"`  //! 是否抢庄
	CurScore int   `json:"curscore"` //! 当前局的分数
	View     bool  `json:"view"`     //! 是否亮牌
	CT       int   `json:"ct"`       //! 当前牌型
	Bets     int   `json:"bets"`     //! 下注
	CS       int   `json:"cs"`       //! 当前局最大牌
	Gold     int   `json:"gold"`     //! 当前金币
	Open     bool  `json:"open"`     //! 是否看牌
	MaxBets  int   `json:"maxbets"`  //! 最大注
}

func (self *Game_GoldTTZ_Person) Init() {
	self.CT = 0
	self.CS = 0
	self.CurScore = 0
	self.Dealer = false
	self.Bets = 0
	self.RobDeal = -1
	self.View = false
	self.Card = make([]int, 0)
	self.Open = false
}

//! 同步金币
func (self *Game_GoldTTZ_Person) SynchroGold(gold int) {
	self.Score += (gold - self.Gold)
	self.Gold = gold
}

type Game_GoldTTZ struct {
	PersonMgr []*Game_GoldTTZ_Person `json:"personmgr"`
	Card      *MahMgr                `json:"card"`
	State     int                    `json:"state"` //! 0准备阶段  1等待抢庄   2等待下注   3等待亮牌
	Time      int64                  `json:"time"`
	DF        int                    `json:"df"` //! 底分
	Cache     []int                  `json:"cache"`
	Stack     []int                  `json:"stack"`
	SZ        []int                  `json:"sz"`

	room *Room
}

func NewGame_GoldTTZ() *Game_GoldTTZ {
	game := new(Game_GoldTTZ)
	game.PersonMgr = make([]*Game_GoldTTZ_Person, 0)

	return game
}

func (self *Game_GoldTTZ) OnInit(room *Room) {
	self.room = room

	self.DF = staticfunc.GetCsvMgr().GetDF(self.room.Type)
}

func (self *Game_GoldTTZ) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldTTZ) OnSendInfo(person *Person) {
	//! 观众模式游戏,观众进来只发送游戏信息
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			self.PersonMgr[i].SynchroGold(person.Gold)
			person.SendMsg("gamegoldttzinfo", self.getInfo(person.Uid))
			return
		}
	}

	person.SendMsg("gamegoldttzinfo", self.getInfo(0))

	if !self.room.Begin {
		if len(self.room.Uid)+len(self.room.Viewer) == lib.HF_Atoi(self.room.csv["minnum"]) { //! 进来的人满足最小开的人数
			self.SetTime(15)
		}
	}

	if !self.room.Begin {
		if self.room.Seat(person.Uid) {
			_person := new(Game_GoldTTZ_Person)
			_person.Init()
			_person.Uid = person.Uid
			_person.Score = person.Gold
			_person.Gold = person.Gold
			_person.Ready = false
			self.PersonMgr = append(self.PersonMgr, _person)
		}
	}
}

func (self *Game_GoldTTZ) OnMsg(msg *RoomMsg) {
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
	case "gameview": //! 亮牌
		self.GameView(msg.Uid, true)
	case "gamedealer": //! 抢庄
		self.GameDeal(msg.Uid, msg.V.(*Msg_GameDealer).Score)
	case "gameopen":
		self.GameOpen(msg.Uid)
	}
}

func (self *Game_GoldTTZ) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)
	self.State = 1

	//! 扣除底分
	bl := 35.0
	if self.room.Type%10 == 1 {
		bl = 50.0
	}
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

	self.SZ = make([]int, 2)
	self.SZ[0] = lib.HF_GetRandom(6) + 1
	self.SZ[1] = lib.HF_GetRandom(6) + 1

	//! 发牌
	self.Cache = make([]int, 0)
	if self.Card == nil || len(self.Card.Card) == 0 {
		self.Card = NewMah_TTZGold()
		self.Stack = make([]int, 0)
	}
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Card = self.Card.Deal(2)
	}
	for i := 0; i < 5-len(self.PersonMgr); i++ {
		self.Cache = append(self.Cache, self.Card.Deal(2)...)
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gamegoldttzbegin", self.getInfo(person.Uid))
	}

	self.room.broadCastMsgView("gamegoldttzbegin", self.getInfo(0))

	self.SetTime(15)

	self.room.flush()
}

//! 抢庄
func (self *Game_GoldTTZ) GameDeal(uid int64, score int) {
	if !self.room.Begin { //! 未开始不能抢庄
		return
	}

	if score > 4 || score < 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "下注超过上限")
		return
	}

	robnum := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].RobDeal >= 0 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能重复抢庄")
				return
			} else {
				if score > 0 {
					basebs := 5 * (len(self.room.Uid) - 1) * self.DF
					maxscore := self.PersonMgr[i].Score / basebs
					if score > maxscore {
						person := GetPersonMgr().GetPerson(uid)
						if person == nil {
							return
						}
						if maxscore == 0 {
							person.SendErr("金币不足，无法抢庄")
						} else {
							person.SendErr(fmt.Sprintf("金币不足，您最大只能抢%d倍", maxscore))
						}
						return
					}
				}
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

	if robnum == len(self.PersonMgr) { //! 全部发表了意见
		deal := make([]*Game_GoldTTZ_Person, 0)
		for i := 0; i < len(self.PersonMgr); i++ {
			if len(deal) == 0 {
				deal = append(deal, self.PersonMgr[i])
			} else {
				if self.PersonMgr[i].RobDeal > deal[0].RobDeal {
					deal = make([]*Game_GoldTTZ_Person, 0)
					deal = append(deal, self.PersonMgr[i])
				} else if self.PersonMgr[i].RobDeal == deal[0].RobDeal {
					deal = append(deal, self.PersonMgr[i])
				}
			}
		}

		dealer := deal[lib.HF_GetRandom(len(deal))]
		dealer.Dealer = true
		if dealer.RobDeal <= 0 {
			dealer.RobDeal = 1
		}

		basebs := dealer.RobDeal * 5 * self.DF
		for i := 0; i < len(self.PersonMgr); i++ {
			usescore := lib.HF_MinInt(self.PersonMgr[i].Score, dealer.Score/(len(self.PersonMgr)-1))
			self.PersonMgr[i].MaxBets = lib.HF_MinInt(lib.HF_MaxInt(usescore/basebs, 1), 5)
			var msg Msg_GameGoldNN_Dealer
			msg.Uid = dealer.Uid
			msg.Bets = self.PersonMgr[i].MaxBets
			self.room.SendMsg(self.PersonMgr[i].Uid, "gamedealer", &msg)
		}

		var msg Msg_GameGoldNN_Dealer
		msg.Uid = dealer.Uid
		msg.Bets = 1
		self.room.broadCastMsgView("gamedealer", &msg)

		//! 下注
		self.State = 2
		if len(deal) > 1 {
			self.SetTime(12)
		} else {
			self.SetTime(10)
		}
	}

	self.room.flush()
}

//! 看牌
func (self *Game_GoldTTZ) GameOpen(uid int64) {
	if !self.room.Begin {
		return
	}

	if self.room.Type%10 != 1 { //! 看牌抢庄没有这个步骤
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
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "没有庄家无法看牌")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if len(person.Card) < 2 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "还不能看牌")
		return
	}

	if person.Open {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经看牌")
		return
	}

	person.CT, person.CS = GetTTZResult1(person.Card)

	var msg Msg_GameGoldNN_Open
	msg.Card = person.Card
	msg.CT = person.CT
	self.room.SendMsg(person.Uid, "gamegoldttzopen", &msg)

	person.Open = true

	self.room.flush()
}

//! 亮牌
func (self *Game_GoldTTZ) GameView(uid int64, send bool) {
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
		return
	}

	if len(person.Card) < 2 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "还不能亮牌")
		return
	}

	if person.View {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经亮牌了")
		return
	}

	if !person.Open {
		person.CT, person.CS = GetTTZResult1(person.Card)
	}

	if send || !person.Open {
		var msg Msg_GameNiuNiuJX_View
		msg.Uid = uid
		msg.Card = person.Card
		msg.CT = person.CT
		self.room.broadCastMsg("gameview", &msg)
	}
	person.Open = true
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
func (self *Game_GoldTTZ) GameReady(uid int64) {
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
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能重复准备")
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

		_person := new(Game_GoldTTZ_Person)
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
func (self *Game_GoldTTZ) GameBets(uid int64, bets int) {
	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	if bets <= 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "下注无效")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
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

			if bets > 1 {
				basebs := 1
				usescore := self.PersonMgr[i].Score
				for _, value := range self.PersonMgr {
					if value.Dealer {
						basebs = value.RobDeal
						usescore = lib.HF_MinInt(usescore, value.Score/(len(self.PersonMgr)-1))
						break
					}
				}
				basebs *= (5 * self.DF)
				maxscore := lib.HF_MaxInt(usescore/basebs, 1)
				if bets > maxscore {
					person := GetPersonMgr().GetPerson(uid)
					if person == nil {
						return
					}
					person.SendErr(fmt.Sprintf("您最大只能下%d倍", maxscore))
					return
				}
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
		for i := 0; i < len(self.PersonMgr); i++ {
			var msg Msg_GameNiuNiuJX_Card
			msg.Card = self.PersonMgr[i].Card
			person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
			if person != nil {
				person.SendMsg("gamegoldttzcard", &msg)
			}
		}

		//! 亮牌
		self.State = 3
		self.SetTime(8)

	}

	self.room.flush()
}

//! 结算
func (self *Game_GoldTTZ) OnEnd() {
	self.room.SetBegin(false)
	self.State = 0
	self.Time = 0

	var dealer *Game_GoldTTZ_Person = nil
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Dealer {
			dealer = self.PersonMgr[i]
			break
		}
	}

	lst := make([]*Game_GoldTTZ_Person, 0)
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Ready = false
		if self.PersonMgr[i].Uid != dealer.Uid {
			lst = append(lst, self.PersonMgr[i])
		}
		self.Stack = append(self.Stack, self.PersonMgr[i].Card...)
	}
	self.Stack = append(self.Stack, self.Cache...)

	for i := 0; i < len(lst); i++ {
		dealerwin := false
		if dealer.CT >= lst[i].CT { //! 庄家赢
			dealerwin = true
		} else { //! 闲家赢
			dealerwin = false
		}

		if dealerwin { //! 庄家赢
			score := lst[i].Bets * dealer.RobDeal * dealer.CS
			if lst[i].CurScore+lst[i].Score >= score*self.DF {
				lst[i].CurScore -= score * self.DF
				dealer.CurScore += score * self.DF
			} else {
				abs := lst[i].CurScore + lst[i].Score
				lst[i].CurScore -= abs
				dealer.CurScore += abs
			}
		} else { //! 闲家赢
			score := lst[i].Bets * dealer.RobDeal * lst[i].CS
			if dealer.CurScore+dealer.Score >= score*self.DF {
				dealer.CurScore -= score * self.DF
				lst[i].CurScore += score * self.DF
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

	//! 发消息
	var msg Msg_GameTTZ_End
	msg.Cache = self.Cache
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Ready = false

		var son Son_GameTTZ_Info
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
	}
	recordinfo := lib.HF_JtoA(&record)
	for i := 0; i < len(record.Info); i++ {
		GetServer().InsertRecord(self.room.Type, record.Info[i].Uid, recordinfo, -record.Info[i].Score)
	}
	self.room.broadCastMsg("gamegoldttzend", &msg)

	self.State = 0
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
			_person := new(Game_GoldTTZ_Person)
			_person.Init()
			_person.Uid = person.Uid
			_person.Score = person.Gold
			_person.Gold = person.Gold
			_person.Ready = false
			self.PersonMgr = append(self.PersonMgr, _person)
		} else {
			i++
		}
	}

	self.room.flush()
}

func (self *Game_GoldTTZ) OnBye() {
}

func (self *Game_GoldTTZ) OnExit(uid int64) {
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

func (self *Game_GoldTTZ) getInfo(uid int64) *Msg_GameGoldTTZ_Info {
	var msg Msg_GameGoldTTZ_Info
	msg.SZ = self.SZ
	msg.Begin = self.room.Begin
	msg.State = self.State
	msg.Card = self.Stack
	if self.Time != 0 {
		msg.Time = self.Time - time.Now().Unix()
	}
	msg.Info = make([]Son_GameGoldTTZ_Info, 0)
	for _, value := range self.PersonMgr {
		if value.Uid == uid {
			msg.Bets = lib.HF_MinInt(lib.HF_MaxInt(value.MaxBets, 1), 5)
		}
		var son Son_GameGoldTTZ_Info
		son.Uid = value.Uid
		son.Ready = value.Ready
		son.Bets = value.Bets
		son.Dealer = value.Dealer
		son.Total = value.Score
		son.Score = value.CurScore
		son.RobDeal = value.RobDeal
		son.Open = value.Open
		son.View = value.View
		if self.room.Type%10 == 1 { //! 自由抢庄
			if (value.Uid == uid && self.State == 3) || !msg.Begin || value.View {
				son.Card = value.Card
				son.CT = value.CT
			} else {
				son.Card = make([]int, len(value.Card))
				son.CT = 0
			}
		} else { //! 看牌抢庄
			if value.Uid == uid { //! 是自己或者亮牌了或者已经结束了
				lib.HF_DeepCopy(&son.Card, &value.Card)
				son.CT = value.CT
				if !value.View && msg.Begin {
					son.Card[1] = 0
					son.CT = 0
				}
			} else {
				son.Card = make([]int, len(value.Card))
				son.CT = 0
			}
		}

		msg.Info = append(msg.Info, son)
	}
	return &msg
}

func (self *Game_GoldTTZ) GetPerson(uid int64) *Game_GoldTTZ_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

func (self *Game_GoldTTZ) OnTime() {
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
					self.GameDeal(self.PersonMgr[i].Uid, 0)
				}
			}
		} else if self.State == 2 { //! 下注
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Bets <= 0 && !self.PersonMgr[i].Dealer {
					self.GameBets(self.PersonMgr[i].Uid, 1)
				}
			}
		} else if self.State == 3 { //! 亮牌
			for i := 0; i < len(self.PersonMgr); i++ {
				if !self.PersonMgr[i].View {
					self.GameView(self.PersonMgr[i].Uid, false)
				}
			}
		}
	}
}

func (self *Game_GoldTTZ) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_GoldTTZ) OnIsBets(uid int64) bool {
	return false
}

//! 同步总分
func (self *Game_GoldTTZ) SendTotal() {
	var msg Msg_GameKWX_Total
	for i := 0; i < len(self.PersonMgr); i++ {
		self.room.Param[i] = self.PersonMgr[i].Score
		msg.Info = append(msg.Info, Son_GameKWX_Total{self.PersonMgr[i].Uid, self.PersonMgr[i].Score})
	}
	self.room.broadCastMsg("gamegoldtotal", &msg)
}

//! 设置时间
func (self *Game_GoldTTZ) SetTime(t int) {
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
func (self *Game_GoldTTZ) OnBalance() {
	for i := 0; i < len(self.PersonMgr); i++ {
		//! 退出房间同步金币
		gold := self.PersonMgr[i].Score - self.PersonMgr[i].Gold
		if gold > 0 {
			GetRoomMgr().AddCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
		} else if gold < 0 {
			GetRoomMgr().CostCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, -gold, self.room)
		}
		self.PersonMgr[i].Gold = self.PersonMgr[i].Score
	}
}
