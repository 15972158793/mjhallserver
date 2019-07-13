package gameserver

import (
	"fmt"
	"lib"
	"math"
	"sort"
	"staticfunc"
	"time"
)

//! param1
//! 00000000
const TYPE_GOLDZJH_BP = 0  //! 0比大小  1比花色   2全比
const TYPE_GOLDZJH_BZ = 1  //! 豹子是否额外奖励 0否 1是
const TYPE_GOLDZJH_SB = 2  //! 比牌双倍开  0否  1是
const TYPE_GOLDZJH_PDD = 3 //! 是否全压   0否  1是
const TYPE_GOLDZJH_FD = 4  //! 封顶   0,5轮  1,10轮   2,15轮
const TYPE_GOLDZJH_BL = 5  //! 比牌轮数 0,1轮   1,2轮   2,3轮
const TYPE_GOLDZJH_MP = 6  //! 闷牌轮数  0不闷  1,2轮   2,3轮   3,4轮
const TYPE_GOLDZJH_DF = 7  //! 0,50  1,100  2,200   3,500    4,1000
const TYPE_GOLDZJH_MAX = 8

type Game_GoldZJHRoom struct {
	Ready     []int64                `json:"ready"`    //! 已经准备的人
	Bets      []int64                `json:"bets"`     //! 已经下注的人
	CurStep   int                    `json:"curstep"`  //! 谁出牌
	Round     int                    `json:"round"`    //! 轮数
	Point     int                    `json:"point"`    //! 当前注
	Allpoint  int                    `json:"allpoint"` //! 当前总注
	PersonMgr []*Game_GoldZJH_Person `json:"personmgr"`
	FirstStep int                    `json:"firststep"`
	DF        int                    `json:"df"`
	Time      int64                  `json:"time"`    //! 自动倒计时
	DisList   []Game_GoldZJH_Discard `json:"dislist"` //! 弃了走的人
	AllIn     bool                   `json:"allin"`

	room *Room
}

func NewGame_GoldZJHRoom() *Game_GoldZJHRoom {
	game := new(Game_GoldZJHRoom)
	game.Ready = make([]int64, 0)
	game.Bets = make([]int64, 0)
	game.PersonMgr = make([]*Game_GoldZJH_Person, 0)

	return game
}

func (self *Game_GoldZJHRoom) GetParam(_type int) int {
	return self.room.Param1 % int(math.Pow(10.0, float64(TYPE_GOLDZJH_MAX-_type))) / int(math.Pow(10.0, float64(TYPE_GOLDZJH_MAX-_type-1)))
}

func (self *Game_GoldZJHRoom) OnInit(room *Room) {
	self.room = room

	df := self.GetParam(TYPE_GOLDZJH_DF)
	if df == 0 {
		self.DF = 50
	} else if df == 1 {
		self.DF = 100
	} else if df == 2 {
		self.DF = 200
	} else if df == 3 {
		self.DF = 500
	} else if df == 4 {
		self.DF = 1000
	}
}

func (self *Game_GoldZJHRoom) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldZJHRoom) OnSendInfo(person *Person) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			self.PersonMgr[i].SynchroGold(person.Gold)
			person.SendMsg("gamezjhinfo", self.getInfo(person.Uid))
			return
		}
	}

	person.SendMsg("gamezjhinfo", self.getInfo(person.Uid))

	if !self.room.Begin {
		if len(self.room.Uid)+len(self.room.Viewer) == lib.HF_Atoi(self.room.csv["minnum"]) { //! 进来的人满足最小开的人数
			self.SetTime(15)
		}

		if self.room.Seat(person.Uid) {
			_person := new(Game_GoldZJH_Person)
			_person.Uid = person.Uid
			_person.Name = person.Name
			_person.IP = person.ip
			_person.Head = person.Imgurl
			_person.Score = person.Gold
			_person.Gold = person.Gold
			self.PersonMgr = append(self.PersonMgr, _person)
		}
	}
}

func (self *Game_GoldZJHRoom) OnMsg(msg *RoomMsg) {
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
	case "gameshsuo":
		self.GameAllIn(msg.Uid)
	//case "gameallbets": //! 跟到底
	//	self.GameAllBets(msg.Uid, msg.V.(*Msg_GameBets).Bets)
	case "gamecompare": //! 比牌
		self.GameCompare(msg.Uid, msg.V.(*Msg_GameCompare).Destuid)
		//case "gameend": //! 解散
		//	self.GameEnd(msg.Uid)
	}
}

//! 结束
func (self *Game_GoldZJHRoom) GameEnd(uid int64) {
	if self.room.Host != uid {
		return
	}

	self.room.Bye()
}

func (self *Game_GoldZJHRoom) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)
	self.Point = 1
	self.Allpoint = 0
	self.Round = 0
	self.AllIn = false
	//self.Record = new(staticfunc.Rec_Gold_Info)

	//! 扣除底分
	for i := 0; i < len(self.PersonMgr); i++ {
		cost := int(math.Ceil(float64(self.DF) * 50.0 / 100.0))
		self.PersonMgr[i].Score -= cost
		GetServer().SqlAgentGoldLog(self.PersonMgr[i].Uid, cost, self.room.Type)
		GetServer().SqlAgentBillsLog(self.PersonMgr[i].Uid, cost, self.room.Type)
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i]._ct = 0
		self.PersonMgr[i]._cs = 0
		self.PersonMgr[i].CurScore = 0
		self.PersonMgr[i].CurBaozi = 0
		self.PersonMgr[i].Lose = false
		self.PersonMgr[i].Open = false
		self.PersonMgr[i].Discard = false
		self.PersonMgr[i].Dealer = false
		self.PersonMgr[i].Score -= self.DF
		self.PersonMgr[i].Allbets = self.DF
		self.PersonMgr[i].CanOpen = make([]int64, 0)
		self.PersonMgr[i].AllIn = 0
		self.Allpoint += self.DF
	}
	self.SendTotal()

	//! 确定庄家
	DealerPos := lib.HF_GetRandom(len(self.PersonMgr)) - 1
	self.PersonMgr[DealerPos+1].Dealer = true
	self.CurStep = DealerPos + 1
	self.FirstStep = DealerPos + 1
	//self.NextPlayer()

	//! 发牌
	cardmgr := NewCard_ZJH()
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Card = cardmgr.Deal(3)
		sort.Ints(self.PersonMgr[i].Card)

		//! 记录
		var rc_person staticfunc.Son_Rec_Gold_Person
		rc_person.Uid = self.PersonMgr[i].Uid
		rc_person.Name = self.room.GetName(rc_person.Uid)
		rc_person.Head = self.room.GetHead(rc_person.Uid)
		//self.Record.Info = append(self.Record.Info, rc_person)
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gamezjhbegin", self.getInfo(person.Uid))
	}

	self.SetTime(60)

	self.room.flush()
}

//! 比牌
func (self *Game_GoldZJHRoom) GameCompare(uid int64, destuid int64) {
	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	if uid == destuid {
		return
	}

	if self.Round < self.GetParam(TYPE_GOLDZJH_BL)+1 {
		person := GetPersonMgr().GetPerson(uid)
		if person != nil {
			person.SendErr(fmt.Sprintf("%d轮后才能比牌", self.GetParam(TYPE_GOLDZJH_BL)+1))
		}
		return
	}

	find := false
	var player1 int
	var player2 int
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			player1 = i
			//player1 = self.PersonMgr[i]
		} else if self.PersonMgr[i].Uid == destuid {
			find = true
			player2 = i
			//player2 = &self.PersonMgr[i]

			if self.PersonMgr[i].Discard {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "比牌玩家已经弃牌")
				return
			}

			if self.PersonMgr[i].Lose {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "比牌玩家已经输牌")
				return
			}
		}
	}

	if !find {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "比牌没找到玩家")
		return
	}

	addpoint := self.Point
	if self.GetParam(TYPE_GOLDZJH_SB) == 1 {
		addpoint *= 2
	}
	if self.PersonMgr[player1].Open {
		addpoint *= 2
	}

	if self.PersonMgr[player1].Score < addpoint {
		person := GetPersonMgr().GetPerson(uid)
		if person != nil {
			person.SendErr("您的金币不足，请前往充值。")
		}
		return
	}

	var win int
	if self.GetParam(TYPE_GOLDZJH_BP) == 0 {
		win = ZjhCardCompare(self.PersonMgr[player1].Card, self.PersonMgr[player2].Card)
	} else if self.GetParam(TYPE_GOLDZJH_BP) == 1 {
		win = ZjhCardCompare1(self.PersonMgr[player1].Card, self.PersonMgr[player2].Card)
	} else {
		win = ZjhCardCompare2(self.PersonMgr[player1].Card, self.PersonMgr[player2].Card)
	}

	if win == 0 {
		self.PersonMgr[player2].Lose = true
	} else {
		self.PersonMgr[player1].Lose = true
	}
	self.PersonMgr[player1].CanOpen = append(self.PersonMgr[player1].CanOpen, self.PersonMgr[player2].Uid)
	self.PersonMgr[player2].CanOpen = append(self.PersonMgr[player2].CanOpen, self.PersonMgr[player1].Uid)

	self.PersonMgr[player1].Bets = addpoint
	self.PersonMgr[player1].Allbets += addpoint
	self.PersonMgr[player1].Score -= addpoint

	self.Allpoint += addpoint

	self.NextPlayer()

	var msg Msg_GameZJH_Com
	msg.Uid = uid
	msg.Destuid = destuid
	msg.Win = (win == 0)
	msg.Point = self.Point
	msg.Addpoint = addpoint
	msg.Allpoint = self.Allpoint
	msg.Allbets = self.PersonMgr[player1].Score
	msg.Opuid = self.room.Uid[self.CurStep]
	msg.Round = self.Round

	for i := 0; i < len(self.PersonMgr); i++ {
		if uid == self.PersonMgr[i].Uid {
			msg.Card1 = self.PersonMgr[player1].Card

			if msg.Win {
				msg.Card2 = self.PersonMgr[player2].Card
			} else {
				msg.Card2 = make([]int, 0)
				for j := 0; j < 3; j++ {
					msg.Card2 = append(msg.Card2, 0)
				}
			}
		} else {
			msg.Card1 = make([]int, 0)
			msg.Card2 = make([]int, 0)
			for j := 0; j < 3; j++ {
				msg.Card1 = append(msg.Card1, 0)
				msg.Card2 = append(msg.Card2, 0)
			}
		}

		self.room.SendMsg(self.PersonMgr[i].Uid, "gamecompare", &msg)
	}

	//if self.Record != nil {
	//	self.Record.Step = append(self.Record.Step, staticfunc.Son_Rec_GameJZH_Step{uid, destuid, win})
	//}

	self.room.flush()

	count := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if !self.PersonMgr[i].Discard && !self.PersonMgr[i].Lose {
			count++
		}
	}

	if count == 1 {
		self.OnEnd()
		return
	} else {
		self.GameView(uid, 0)
	}

	self.SetTime(60)
}

//! 看牌
func (self *Game_GoldZJHRoom) GameView(uid int64, _type int) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Open {
				return
			}

			self.PersonMgr[i].Open = true
			card := self.PersonMgr[i].Card

			for i := 0; i < len(self.PersonMgr); i++ {
				var msg Msg_GameZJH_View
				msg.Uid = uid
				msg.Type = _type
				for j := 0; j < 3; j++ {
					if self.PersonMgr[i].Uid == uid || GetServer().IsAdmin(self.PersonMgr[i].Uid, staticfunc.ADMIN_SZP) {
						msg.Card = append(msg.Card, card[j])
					} else {
						msg.Card = append(msg.Card, 0)
					}
				}
				self.room.SendMsg(self.PersonMgr[i].Uid, "gameview", &msg)
			}

			//if self.Record != nil {
			//	self.Record.Step = append(self.Record.Step, staticfunc.Son_Rec_GameJZH_Step{uid, int64(0), self.PersonMgr[i].Allbets})
			//}

			break
		}
	}

	self.room.flush()
}

func (self *Game_GoldZJHRoom) NextPlayer() {
	for i := 0; i < len(self.PersonMgr); i++ {
		self.CurStep++
		self.CurStep %= len(self.PersonMgr)
		if self.CurStep == self.FirstStep {
			self.Round++
		}
		if !self.PersonMgr[self.CurStep].Lose && !self.PersonMgr[self.CurStep].Discard {
			break
		}
	}
}

//! 弃牌
func (self *Game_GoldZJHRoom) GameDiscard(uid int64) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Discard {
				return
			}

			self.PersonMgr[i].Discard = true
			break
		}
	}

	if self.room.Uid[self.CurStep] == uid {
		self.NextPlayer()
	}

	var msg Msg_GameZJH_Discard
	msg.Uid = uid
	msg.Opuid = self.room.Uid[self.CurStep]
	msg.Round = self.Round
	self.room.broadCastMsg("gamediscard", &msg)

	count := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if !self.PersonMgr[i].Discard && !self.PersonMgr[i].Lose {
			count++
		}
	}

	end := true
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Discard || self.PersonMgr[i].Lose {
			continue
		}
		if self.PersonMgr[i].AllIn == 0 {
			end = false
			break
		}
	}

	if count == 1 || end {
		self.OnEnd()
		return
	}

	self.SetTime(60)

	self.room.flush()

}

func (self *Game_GoldZJHRoom) GameAllIn(uid int64) {
	if self.GetParam(TYPE_GOLDZJH_FD) == 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "当前模式不能全压")
		return
	}

	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}
	if self.room.Uid[self.CurStep] != uid {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "当前玩家不能操作")
		return
	}
	person := self.GetPerson(uid)
	if person == nil {
		return
	}
	self.AllIn = true

	person.Bets = person.Score
	person.AllIn = person.Score
	person.Allbets += person.Score
	person.Score = 0

	self.NextPlayer()

	self.Point = person.Bets
	self.Allpoint += person.Bets

	var msg Msg_GameZJH_Bets
	msg.Uid = uid
	msg.OpUid = self.room.Uid[self.CurStep]
	msg.Bets = 0
	msg.Point = self.Point
	msg.Addpoint = person.Bets
	msg.Allpoint = self.Allpoint
	msg.Allbets = 0
	msg.Round = self.Round
	self.room.broadCastMsg("gameshsuo", &msg)

	end := true
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Discard || self.PersonMgr[i].Lose {
			continue
		}
		if self.PersonMgr[i].AllIn == 0 {
			end = false
			break
		}
	}

	if self.Round >= 20 || end {
		self.OnEnd()
		return
	}

	self.SetTime(60)

	self.room.flush()
}

//! 准备
func (self *Game_GoldZJHRoom) GameReady(uid int64) {
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
			if self.PersonMgr[i].Score < self.DF*20 { //! 携带的金币不足，踢出去
				self.room.KickPerson(uid, 99)
				return
			}
			find = true
			break
		}
	}

	if !find { //! 坐下
		if !self.room.Seat(uid) {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "无法坐下")
			return
		}

		person := GetPersonMgr().GetPerson(uid)
		if person == nil {
			return
		}

		_person := new(Game_GoldZJH_Person)
		_person.Uid = uid
		_person.Name = person.Name
		_person.Head = person.Imgurl
		_person.Score = person.Gold
		_person.Gold = person.Gold
		_person.IP = person.ip
		self.PersonMgr = append(self.PersonMgr, _person)
	}

	for i := 0; i < len(self.Ready); i++ {
		if self.Ready[i] == uid {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "同一个玩家准备")
			return
		}
	}

	self.Ready = append(self.Ready, uid)

	if len(self.Ready) == len(self.room.Uid)+len(self.room.Viewer) && len(self.Ready) >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 准备的人数达到游戏最小人数
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
		self.OnBegin()
		return
	}

	var msg staticfunc.Msg_Uid
	msg.Uid = uid
	self.room.broadCastMsg("gameready", &msg)

	if len(self.Ready) == lib.HF_Atoi(self.room.csv["minnum"]) {
		self.SetTime(10)
	}

	self.room.flush()
}

//! 下注 -1表示弃牌 0表示看牌 1表示跟住 2-5表示加注
func (self *Game_GoldZJHRoom) GameBets(uid int64, bets int) {
	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	if bets == 0 || bets == -2 {
		if self.GetParam(TYPE_GOLDZJH_MP) > 0 {
			if self.Round < self.GetParam(TYPE_GOLDZJH_MP)+1 {
				person := GetPersonMgr().GetPerson(uid)
				if person != nil {
					person.SendErr(fmt.Sprintf("%d轮后才能看牌", self.GetParam(TYPE_GOLDZJH_MP)+1))
				}
				return
			}
		}
		self.GameView(uid, bets)
		return
	} else if bets == -1 {
		self.GameDiscard(uid)
		return
	} else if bets < 1 || bets > 5 {
		return
	}

	if self.room.Uid[self.CurStep] != uid {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "当前玩家不能操作")
		return
	}

	_bets := bets
	if self.Point/self.DF > bets {
		bets = self.Point / self.DF
	}

	addpoint := bets * self.DF

	allbets := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Open {
				addpoint *= 2
				self.PersonMgr[i].Minpai++
			} else {
				self.PersonMgr[i].Menpai++
			}
			if self.PersonMgr[i].Score < addpoint {
				person := GetPersonMgr().GetPerson(uid)
				if person != nil {
					person.SendErr("您的金币不足，请前往充值。")
				}
				return
			}
			self.PersonMgr[i].Bets = addpoint
			self.PersonMgr[i].Allbets += addpoint
			self.PersonMgr[i].Score -= addpoint
			allbets = self.PersonMgr[i].Score

			//if self.Record != nil {
			//	self.Record.Step = append(self.Record.Step, staticfunc.Son_Rec_GameJZH_Step{uid, int64(bets), self.PersonMgr[i].Allbets})
			//}

			break
		}
	}

	self.NextPlayer()

	self.Point = bets * self.DF
	self.Allpoint += addpoint

	var msg Msg_GameZJH_Bets
	msg.Uid = uid
	msg.OpUid = self.room.Uid[self.CurStep]
	msg.Bets = _bets
	msg.Point = self.Point
	msg.Addpoint = addpoint
	msg.Allpoint = self.Allpoint
	msg.Allbets = allbets
	msg.Round = self.Round
	self.room.broadCastMsg("gamebets", &msg)

	if self.Round >= (self.GetParam(TYPE_GOLDZJH_FD)+1)*5 {
		self.OnEnd()
		return
	}

	self.SetTime(60)

	self.room.flush()
}

//! 跟到底
//func (self *Game_GoldZJHRoom) GameAllBets(uid int64, bets int) {
//	if !self.room.Begin { //! 没有开始不能下注
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
//		return
//	}

//	if self.GetParam(TYPE_GOLDZJH_PDD) == 0 {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "非拼到底模式不发这个消息")
//		return
//	}

//	if self.room.Uid[self.CurStep] != uid {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "当前玩家不能操作")
//		return
//	}

//	if bets == 0 || bets == -2 { //! 看牌
//		if self.GetParam(TYPE_GOLDZJH_MP) > 0 {
//			if self.Round < self.GetParam(TYPE_GOLDZJH_MP)+1 {
//				person := GetPersonMgr().GetPerson(uid)
//				if person != nil {
//					person.SendErr(fmt.Sprintf("%d轮后才能看牌", self.GetParam(TYPE_GOLDZJH_MP)+1))
//				}
//				return
//			}
//		}
//		self.GameView(uid, bets)
//		return
//	} else if bets == -1 { //! 弃牌
//		self.GameDiscard(uid)
//		return
//	}

//	if bets < self.Point {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "下注比当前注小")
//		return
//	}

//	allbets := 0
//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Uid == uid {
//			if self.PersonMgr[i].Score < bets {
//				person := GetPersonMgr().GetPerson(uid)
//				if person != nil {
//					person.SendErr("您的金币不足，请前往充值。")
//				}
//				return
//			}

//			self.PersonMgr[i].Bets = bets
//			self.PersonMgr[i].Allbets += bets
//			self.PersonMgr[i].Score -= bets
//			allbets = self.PersonMgr[i].Score
//			break
//		}
//	}

//	self.NextPlayer()
//	self.Point = bets
//	self.Allpoint += self.Point

//	var msg Msg_GameZJH_Bets
//	msg.Uid = uid
//	msg.OpUid = self.room.Uid[self.CurStep]
//	msg.Bets = bets
//	msg.Point = self.Point
//	msg.Addpoint = bets
//	msg.Allpoint = self.Allpoint
//	msg.Allbets = allbets
//	msg.Round = self.Round
//	self.room.broadCastMsg("gamebets", &msg)

//	if self.Round >= (self.GetParam(TYPE_GOLDZJH_FD)+1)*5 {
//		self.OnEnd()
//		return
//	}

//	self.SetTime(60)

//	self.room.flush()
//}

//! 结算
func (self *Game_GoldZJHRoom) OnEnd() {
	allpoint := 0
	emplst := make([]*Game_GoldZJH_Person, 0)
	for i := 0; i < len(self.PersonMgr); i++ {
		if !self.PersonMgr[i].Discard && !self.PersonMgr[i].Lose {
			emplst = append(emplst, self.PersonMgr[i])
		} else {
			allpoint += self.PersonMgr[i].Allbets
		}
	}
	for i := 0; i < len(self.DisList); i++ {
		allpoint += self.DisList[i].Score
	}

	oldwin := make([]int, 0)
	oldwin = append(oldwin, 0)
	for i := 1; i < len(emplst); i++ {
		var win int
		if self.GetParam(TYPE_GOLDZJH_BP) == 0 {
			win = ZjhCardCompare(emplst[oldwin[0]].Card, emplst[i].Card)
		} else if self.GetParam(TYPE_GOLDZJH_BP) == 1 {
			win = ZjhCardCompare1(emplst[oldwin[0]].Card, emplst[i].Card)
		} else {
			win = ZjhCardCompare2(emplst[oldwin[0]].Card, emplst[i].Card)
		}

		if win == 0 { //! 胜利
			emplst[i].Lose = true
			allpoint += emplst[i].Allbets
		} else if win == 1 { //! 负
			for _, value := range oldwin {
				emplst[value].Lose = true
				allpoint += emplst[value].Allbets
			}
			oldwin = make([]int, 0)
			oldwin = append(oldwin, i)
		} else { //! 平
			oldwin = append(oldwin, i)
		}
	}

	minallin := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].AllIn == 0 {
			continue
		}
		if minallin == 0 || self.PersonMgr[i].AllIn < minallin {
			minallin = self.PersonMgr[i].AllIn
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].AllIn == 0 {
			continue
		}
		if !self.PersonMgr[i].Discard && !self.PersonMgr[i].Lose { //! 赢的人不算
			continue
		}
		if self.PersonMgr[i].AllIn > minallin {
			self.PersonMgr[i].Score += (self.PersonMgr[i].AllIn - minallin)
			allpoint -= (self.PersonMgr[i].AllIn - minallin)
			self.PersonMgr[i].Allbets -= (self.PersonMgr[i].AllIn - minallin)
		}
	}

	self.room.SetBegin(false)

	self.Ready = make([]int64, 0)
	self.Bets = make([]int64, 0)

	sy := allpoint % len(oldwin)
	if sy > 0 {
		emplst[oldwin[lib.HF_GetRandom(len(oldwin))]].CurScore += sy
	}
	for i := 0; i < len(self.PersonMgr); i++ {
		//baozijiangli := 0
		if !self.PersonMgr[i].Lose && !self.PersonMgr[i].Discard {
			self.PersonMgr[i].CurScore += allpoint / len(oldwin)
			self.PersonMgr[i].Win++
		} else {
			self.PersonMgr[i].CurScore -= self.PersonMgr[i].Allbets
		}

		cardtype, _ := GetZjhType(self.PersonMgr[i].Card)

		switch cardtype {
		case 600:
			self.PersonMgr[i].Baozi++

			//if self.GetParam(TYPE_GOLDZJH_BZ) == 1 {
			//	baozijiangli = self.DF
			//	self.PersonMgr[i].CurBaozi += baozijiangli
			//}

			//if baozijiangli > 0 {
			//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "22222")
			//	for j := 0; j < len(self.PersonMgr); j++ {
			//		if self.PersonMgr[j].Uid == self.PersonMgr[i].Uid {
			//			self.PersonMgr[j].CurScore += baozijiangli * (len(self.PersonMgr) - 1)
			//			//self.PersonMgr[j].Score += baozijiangli * (len(self.PersonMgr) - 1)
			//		} else {
			//			self.PersonMgr[j].CurScore -= baozijiangli
			//			//elf.PersonMgr[j].Score -= baozijiangli
			//		}
			//		lib.GetLogMgr().Output(lib.LOG_DEBUG, self.PersonMgr[j].CurScore, baozijiangli)
			//	}
			//}

		case 500:
			self.PersonMgr[i].Shunjin++
		case 400:
			self.PersonMgr[i].Jinhua++
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Score += self.PersonMgr[i].Allbets
		self.PersonMgr[i].Score += self.PersonMgr[i].CurScore
	}

	var record staticfunc.Rec_Gold_Info
	record.Time = time.Now().Unix()
	record.GameType = self.room.Type

	//! 发消息
	//! 发给玩家
	{
		for _, value := range self.PersonMgr {
			var msg Msg_GameZJH_End
			msg.Point = self.Point
			msg.Round = self.Round
			msg.Allpoint = self.Allpoint
			for i := 0; i < len(self.PersonMgr); i++ {
				var son Son_GameZJH_Info

				son.Uid = self.PersonMgr[i].Uid
				son.Name = self.PersonMgr[i].Name
				son.Bets = self.PersonMgr[i].Bets
				if value.Uid == son.Uid || value.IsCanOpen(son.Uid) || GetServer().IsAdmin(value.Uid, staticfunc.ADMIN_SZP) || (self.Round >= 20 && !value.Discard) {
					son.Card = self.PersonMgr[i].Card
				} else {
					son.Card = make([]int, len(self.PersonMgr[i].Card))
				}
				son.Discard = self.PersonMgr[i].Discard
				son.Dealer = self.PersonMgr[i].Dealer
				son.Score = self.PersonMgr[i].CurScore
				son.Total = self.PersonMgr[i].Score
				son.Baozi = self.PersonMgr[i].CurBaozi
				msg.Info = append(msg.Info, son)
			}
			self.room.SendMsg(value.Uid, "gameniuniuend", &msg)
		}
	}
	//! 发给观众
	{
		var roomlog SQL_RoomLog
		roomlog.Id = 1
		roomlog.Time = time.Now().Unix()
		var msg Msg_GameZJH_End
		msg.Point = self.Point
		msg.Round = self.Round
		msg.Allpoint = self.Allpoint
		for i := 0; i < len(self.PersonMgr); i++ {
			var son Son_GameZJH_Info

			son.Uid = self.PersonMgr[i].Uid
			son.Name = self.PersonMgr[i].Name
			son.Bets = self.PersonMgr[i].Bets
			son.Card = make([]int, len(self.PersonMgr[i].Card))
			son.Discard = self.PersonMgr[i].Discard
			son.Dealer = self.PersonMgr[i].Dealer
			son.Score = self.PersonMgr[i].CurScore
			son.Total = self.PersonMgr[i].Score
			son.Baozi = self.PersonMgr[i].CurBaozi
			msg.Info = append(msg.Info, son)

			var rec staticfunc.Son_Rec_Gold_Person
			rec.Uid = self.PersonMgr[i].Uid
			rec.Name = self.room.GetName(self.PersonMgr[i].Uid)
			rec.Head = self.room.GetHead(self.PersonMgr[i].Uid)
			rec.Score = self.PersonMgr[i].CurScore
			record.Info = append(record.Info, rec)

			roomlog.Uid[i] = self.PersonMgr[i].Uid
			roomlog.IP[i] = self.PersonMgr[i].IP
			roomlog.Win[i] = self.PersonMgr[i].CurScore

			self.room.Param[i] = self.PersonMgr[i].Score

			self.PersonMgr[i].Card = make([]int, 0)
			self.PersonMgr[i]._ct = 0
			self.PersonMgr[i]._cs = 0
			self.PersonMgr[i].CurScore = 0
			self.PersonMgr[i].CurBaozi = 0
			self.PersonMgr[i].Lose = false
			self.PersonMgr[i].Open = false
			self.PersonMgr[i].Discard = false
			self.PersonMgr[i].Dealer = false
			self.PersonMgr[i].Allbets = self.DF
		}

		for i := 0; i < len(self.DisList); i++ {
			var rec staticfunc.Son_Rec_Gold_Person
			rec.Uid = self.DisList[i].Uid
			rec.Name = self.DisList[i].Name
			rec.Head = self.DisList[i].Head
			rec.Score = -self.DisList[i].Score
			record.Info = append(record.Info, rec)

			roomlog.Uid[len(self.PersonMgr)+i] = self.DisList[i].Uid
			roomlog.IP[len(self.PersonMgr)+i] = self.DisList[i].IP
			roomlog.Win[len(self.PersonMgr)+i] = -self.DisList[i].Score
		}

		recordinfo := lib.HF_JtoA(&record)
		for i := 0; i < len(record.Info); i++ {
			GetServer().InsertRecord(self.room.Type, record.Info[i].Uid, recordinfo, -record.Info[i].Score)
		}
		self.room.broadCastMsgView("gameniuniuend", &msg)
		GetServer().SqlRoomLog(&roomlog)
	}

	//if self.room.IsBye() {
	//	self.OnBye()
	//	self.room.Bye()
	//	return
	//}

	self.DisList = make([]Game_GoldZJH_Discard, 0)

	for i := 0; i < len(self.room.Viewer); {
		person := GetPersonMgr().GetPerson(self.room.Viewer[i])
		if person == nil {
			i++
			continue
		}
		if self.room.Seat(self.room.Viewer[i]) {
			_person := new(Game_GoldZJH_Person)
			_person.Uid = person.Uid
			_person.Name = person.Name
			_person.Head = person.Imgurl
			_person.Score = person.Gold
			_person.Gold = person.Gold
			_person.IP = person.ip
			self.PersonMgr = append(self.PersonMgr, _person)
		} else {
			i++
		}
	}

	self.SetTime(30)

	self.room.flush()
}

func (self *Game_GoldZJHRoom) OnBye() {
}

func (self *Game_GoldZJHRoom) IsDiscard(uid int64) bool {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i].Discard || self.PersonMgr[i].Lose
		}
	}

	return false
}

func (self *Game_GoldZJHRoom) OnExit(uid int64) {
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

			if self.CurStep > i {
				self.CurStep--
			}
			if self.FirstStep >= i {
				if self.FirstStep > 0 {
					self.FirstStep--
				} else {
					self.FirstStep = len(self.PersonMgr) - 2
				}
			}

			if self.room.Begin {
				self.DisList = append(self.DisList, Game_GoldZJH_Discard{self.PersonMgr[i].Uid, self.PersonMgr[i].Name, self.PersonMgr[i].Head, self.PersonMgr[i].Allbets, self.PersonMgr[i].IP, self.PersonMgr[i].IsRobot})
			}

			copy(self.PersonMgr[i:], self.PersonMgr[i+1:])
			self.PersonMgr = self.PersonMgr[:len(self.PersonMgr)-1]
			break
		}
	}

	if self.room.Begin {
		return
	}

	for i := 0; i < len(self.Ready); i++ {
		if self.Ready[i] == uid {
			copy(self.Ready[i:], self.Ready[i+1:])
			self.Ready = self.Ready[:len(self.Ready)-1]
			break
		}
	}

	if len(self.Ready) == len(self.room.Uid)+len(self.room.Viewer) && len(self.Ready) >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 准备的人数达到游戏最小人数
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
		self.OnBegin()
		return
	}
	if len(self.room.Uid) < lib.HF_Atoi(self.room.csv["minnum"]) {
		self.SetTime(0)
	}
}

func (self *Game_GoldZJHRoom) getInfo(uid int64) *Msg_GameGoldZJH_Info {
	var msg Msg_GameGoldZJH_Info
	msg.Begin = self.room.Begin
	if self.CurStep >= len(self.room.Uid) {
		msg.CurOp = 0
	} else {
		msg.CurOp = self.room.Uid[self.CurStep]
	}
	msg.Ready = make([]int64, 0)
	msg.Round = self.Round
	msg.Point = self.Point
	msg.Allpoint = self.Allpoint
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	msg.Info = make([]Son_GameZJH_Info, 0)
	if !msg.Begin { //! 没有开始,看哪些人已准备
		msg.Ready = self.Ready
	}
	for _, value := range self.PersonMgr {
		var son Son_GameZJH_Info
		son.Uid = value.Uid
		son.Name = value.Name
		son.Bets = value.Bets
		son.Allbets = value.Allbets
		son.Dealer = value.Dealer
		son.Open = value.Open
		son.Discard = value.Discard
		son.Lose = value.Lose
		son.Score = value.CurScore
		son.Total = value.Score
		son.Baozi = value.CurBaozi
		if msg.Begin {
			for i := 0; i < len(value.Card); i++ {
				if son.Uid == uid && son.Open || GetServer().IsAdmin(uid, staticfunc.ADMIN_SZP) {
					son.Card = append(son.Card, value.Card[i])
				} else {
					son.Card = append(son.Card, 0)
				}
			}
		} else {
			son.Card = value.Card
		}

		son.Total = value.Score
		msg.Info = append(msg.Info, son)
	}
	return &msg
}

func (self *Game_GoldZJHRoom) OnTime() {
	if self.Time == 0 {
		return
	}

	if time.Now().Unix() < self.Time {
		return
	}

	if !self.room.Begin {
		for i := 0; i < len(self.PersonMgr); {
			find := false
			for j := 0; j < len(self.Ready); j++ {
				if self.PersonMgr[i].Uid == self.Ready[j] {
					find = true
					break
				}
			}
			if !find {
				self.room.KickPerson(self.PersonMgr[i].Uid, 98)
			} else {
				i++
			}
		}
		self.room.KickView()
		return
	}

	self.GameDiscard(self.room.Uid[self.CurStep])
}

func (self *Game_GoldZJHRoom) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_GoldZJHRoom) OnIsBets(uid int64) bool {
	return false
}

//! 同步总分
func (self *Game_GoldZJHRoom) SendTotal() {
	var msg Msg_GameKWX_Total
	for i := 0; i < len(self.PersonMgr); i++ {
		self.room.Param[i] = self.PersonMgr[i].Score
		msg.Info = append(msg.Info, Son_GameKWX_Total{self.PersonMgr[i].Uid, self.PersonMgr[i].Score})
	}
	self.room.broadCastMsg("gamegoldtotal", &msg)
}

func (self *Game_GoldZJHRoom) GetPerson(uid int64) *Game_GoldZJH_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

//! 设置时间
func (self *Game_GoldZJHRoom) SetTime(t int) {
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
func (self *Game_GoldZJHRoom) OnBalance() {
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
