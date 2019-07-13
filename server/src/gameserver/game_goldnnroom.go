package gameserver

import (
	//"log"
	////"sort"
	"lib"
	"staticfunc"
	//"strings"
	"fmt"
	"math"
	"time"
)

/*
Param1
十位 低分 0-50 1-100 2-200 3-300 4-500 其他-1000
个位 0-看牌抢庄 1-自由抢庄  type 3000
八人抢庄 3001
*/

type Msg_GameGoldNNRoom_Info struct {
	Begin bool                      `json:"begin"` //! 是否开始
	Info  []Son_GameGoldNNRoom_Info `json:"info"`
	State int                       `json:"state"`
	Time  int64                     `json:"time"` //! 倒计时
	Bets  int                       `json:"bets"`
}

type Msg_GameGoldNNRoom_Dealer struct {
	Uid  int64 `json:"uid"`
	Bets int   `json:"bets"`
}

type Msg_GameGoldNNRoom_Open struct {
	Card []int `json:"card"`
	CT   int   `json:"ct"`
}

type Son_GameGoldNNRoom_Info struct {
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
	TZ      int   `json:"tz"`
	Open    bool  `json:"open"`
}

///////////////////////////////////////////////////////
type Game_GoldNNRoom_Person struct {
	Uid      int64 `json:"uid"`
	Card     []int `json:"card"`     //! 手牌
	CardZB   []int `json:"cardzb"`   //! 缓存手牌
	HPL      int   `json:"hpl"`      //! 好牌率
	Ready    bool  `json:"ready"`    //! 是否准备
	Score    int   `json:"score"`    //! 积分
	Dealer   bool  `json:"dealer"`   //! 是否庄家
	RobDeal  int   `json:"robdeal"`  //! 是否抢庄
	CurScore int   `json:"curscore"` //! 当前局的分数
	View     bool  `json:"view"`     //! 是否亮牌
	CT       int   `json:"ct"`       //! 当前牌型
	Bets     int   `json:"bets"`     //! 下注
	CS       int   `json:"cs"`       //! 当前局最大牌
	TZ       int   `json:"tz"`       //! 当前可推注
	Gold     int   `json:"gold"`     //! 当前金币
	Open     bool  `json:"open"`     //! 是否看牌
	MaxBets  int   `json:"maxbets"`  //! 最大注
}

func (self *Game_GoldNNRoom_Person) Init() {
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
func (self *Game_GoldNNRoom_Person) SynchroGold(gold int) {
	self.Score += (gold - self.Gold)
	self.Gold = gold
}

type Game_GoldNNRoom struct {
	PersonMgr []*Game_GoldNNRoom_Person `json:"personmgr"`
	Card      *CardMgr                  `json:"card"`
	State     int                       `json:"state"` //! 0准备阶段  1等待抢庄   2等待下注   3等待亮牌
	Time      int64                     `json:"time"`
	DF        int                       `json:"df"` //! 底分
	CardZB    [][]int                   `json:"cardzb"`
	//Bets      int                   `json:"bets"`

	room *Room
}

func NewGame_GoldNNRoom() *Game_GoldNNRoom {
	game := new(Game_GoldNNRoom)
	game.PersonMgr = make([]*Game_GoldNNRoom_Person, 0)

	return game
}

func (self *Game_GoldNNRoom) OnInit(room *Room) {
	self.room = room

	if self.room.Param1/10%10 == 0 {
		self.DF = 50
	} else if self.room.Param1/10%10 == 1 {
		self.DF = 100
	} else if self.room.Param1/10%10 == 2 {
		self.DF = 200
	} else if self.room.Param1/10%10 == 3 {
		self.DF = 300
	} else if self.room.Param1/10%10 == 4 {
		self.DF = 500
	} else {
		self.DF = 1000
	}
}

func (self *Game_GoldNNRoom) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldNNRoom) OnSendInfo(person *Person) {
	//! 观众模式游戏,观众进来只发送游戏信息
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			self.PersonMgr[i].SynchroGold(person.Gold)
			person.SendMsg("gamegoldnninfo", self.getInfo(person.Uid))
			return
		}
	}

	if !self.room.Begin {
		if len(self.room.Uid)+len(self.room.Viewer) == lib.HF_Atoi(self.room.csv["minnum"]) { //! 进来的人满足最小开的人数
			self.SetTime(15)
		}
	}

	person.SendMsg("gamegoldnninfo", self.getInfo(0))

	if !self.room.Begin {
		if self.room.Seat(person.Uid) {
			_person := new(Game_GoldNNRoom_Person)
			_person.Init()
			_person.Uid = person.Uid
			_person.Score = person.Gold
			_person.Gold = person.Gold
			_person.Ready = false
			self.PersonMgr = append(self.PersonMgr, _person)
		}
	}
}

func (self *Game_GoldNNRoom) OnMsg(msg *RoomMsg) {
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

func (self *Game_GoldNNRoom) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)
	self.State = 1
	self.SetTime(12)

	//! 扣除底分
	bl := 35.0
	if self.room.Param1%10 == 1 {
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

	//! 发牌
	self.Card = NewCard_NiuNiu(false)
	self.CardZB = make([][]int, 0)
	maxIndex := -1 //!　最大牌的数组下标
	for i := 0; i < len(self.PersonMgr); i++ {
		_card := self.Card.Deal(5)
		if maxIndex == -1 {
			maxIndex = i
		} else {
			ct, _ := GetGoldNiuNiuScore(_card, true, true, true, true, true, true, true)
			_ct, _ := GetGoldNiuNiuScore(self.CardZB[maxIndex], true, true, true, true, true, true, true)
			if _ct < ct {
				maxIndex = i
			}
		}
		self.CardZB = append(self.CardZB, _card)
	}

	isHPL := -1
	for i := 0; i < len(self.PersonMgr); i++ {
		if GetServer().IsAdmin(self.PersonMgr[i].Uid, staticfunc.ADMIN_NN30) {
			self.PersonMgr[i].HPL = 30
			if isHPL == -1 {
				isHPL = i
			} else if lib.HF_GetRandom(100) > 50 {
				isHPL = i
			}

		} else if GetServer().IsAdmin(self.PersonMgr[i].Uid, staticfunc.ADMIN_NN50) {
			self.PersonMgr[i].HPL = 50
			if isHPL == -1 {
				isHPL = i
			} else if lib.HF_GetRandom(100) > 50 {
				isHPL = i
			}
		} else if GetServer().IsAdmin(self.PersonMgr[i].Uid, staticfunc.ADMIN_NN80) {
			self.PersonMgr[i].HPL = 80
			if isHPL == -1 {
				isHPL = i
			} else if lib.HF_GetRandom(100) > 50 {
				isHPL = i
			}
		} else if GetServer().IsAdmin(self.PersonMgr[i].Uid, staticfunc.ADMIN_NN100) {
			self.PersonMgr[i].HPL = 100
			if isHPL == -1 {
				isHPL = i
			} else if lib.HF_GetRandom(100) > 50 {
				isHPL = i
			}
		} else {
			self.PersonMgr[i].HPL = 0
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].CardZB = self.CardZB[i]
	}

	if isHPL != -1 && lib.HF_GetRandom(100) < self.PersonMgr[isHPL].HPL {
		maxCard := self.PersonMgr[isHPL].CardZB
		self.PersonMgr[isHPL].CardZB = self.PersonMgr[maxIndex].CardZB
		self.PersonMgr[maxIndex].CardZB = maxCard
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.room.Type%10 == 1 {
			self.PersonMgr[i].Card = make([]int, 0)
		} else {
			self.PersonMgr[i].Card = self.PersonMgr[i].CardZB[0:4]
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gamempqzbegin", self.getInfo(person.Uid))
	}

	self.room.broadCastMsgView("gamempqzbegin", self.getInfo(0))

	self.room.flush()
}

//! 抢庄
func (self *Game_GoldNNRoom) GameDeal(uid int64, score int) {
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
					basebs := 6 * (len(self.room.Uid) - 1) * self.DF
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
		deal := make([]*Game_GoldNNRoom_Person, 0)
		for i := 0; i < len(self.PersonMgr); i++ {
			if len(deal) == 0 {
				deal = append(deal, self.PersonMgr[i])
			} else {
				if self.PersonMgr[i].RobDeal > deal[0].RobDeal {
					deal = make([]*Game_GoldNNRoom_Person, 0)
					deal = append(deal, self.PersonMgr[i])
				} else if self.PersonMgr[i].RobDeal == deal[0].RobDeal {
					deal = append(deal, self.PersonMgr[i])
				}
			}
		}

		dealer := deal[lib.HF_GetRandom(len(deal))]
		dealer.Dealer = true
		dealer.TZ = 0
		if dealer.RobDeal <= 0 {
			dealer.RobDeal = 1
		}

		//! 下注
		self.State = 2
		if len(deal) > 1 {
			self.SetTime(10)
		} else {
			self.SetTime(8)
		}

		basebs := dealer.RobDeal * 6 * self.DF
		for i := 0; i < len(self.PersonMgr); i++ {
			usescore := lib.HF_MinInt(self.PersonMgr[i].Score, dealer.Score/(len(self.PersonMgr)-1))
			self.PersonMgr[i].MaxBets = lib.HF_MinInt(lib.HF_MaxInt(usescore/basebs, 1), 5)
			var msg Msg_GameGoldNNRoom_Dealer
			msg.Uid = dealer.Uid
			msg.Bets = self.PersonMgr[i].MaxBets
			self.room.SendMsg(self.PersonMgr[i].Uid, "gamedealer", &msg)
		}

		var msg Msg_GameGoldNNRoom_Dealer
		msg.Uid = dealer.Uid
		msg.Bets = 1
		self.room.broadCastMsgView("gamedealer", &msg)
	}

	self.room.flush()
}

//! 看牌
func (self *Game_GoldNNRoom) GameOpen(uid int64) {
	if !self.room.Begin {
		return
	}

	if self.room.Param1%10 != 1 { //! 看牌抢庄没有这个步骤
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

	if len(person.Card) < 5 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "还不能看牌")
		return
	}

	if person.Open {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经看牌")
		return
	}

	//person.Card = self.Card.Deal(5)
	person.CT, person.CS = GetGoldNiuNiuScore(person.Card, true, true, true, true, true, true, true)

	var msg Msg_GameGoldNNRoom_Open
	msg.Card = person.Card
	msg.CT = person.CT
	self.room.SendMsg(person.Uid, "gamempqzopen", &msg)

	person.Open = true

	self.room.flush()
}

//! 亮牌
func (self *Game_GoldNNRoom) GameView(uid int64, send bool) {
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

	if len(person.Card) < 5 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "还不能亮牌")
		return
	}

	if person.View {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经亮牌了")
		return
	}

	if !person.Open {
		person.CT, person.CS = GetGoldNiuNiuScore(person.Card, true, true, true, true, true, true, true)
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
func (self *Game_GoldNNRoom) GameReady(uid int64) {
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

			if self.PersonMgr[i].Score < self.DF*20 { //! 携带的金币不足，踢出去
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

		_person := new(Game_GoldNNRoom_Person)
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

	if num == lib.HF_Atoi(self.room.csv["minnum"]) {
		self.SetTime(10)
	}

	var msg staticfunc.Msg_Uid
	msg.Uid = uid
	self.room.broadCastMsg("gameready", &msg)

	self.room.flush()
}

//! 下注
func (self *Game_GoldNNRoom) GameBets(uid int64, bets int) {
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

	if bets != person.TZ { //! 没有推注,下一轮可以推注
		if bets < 1 || bets > 5 {
			return
		}
		person.TZ = 0
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
				basebs *= (6 * self.DF)
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
		//! 亮牌
		self.State = 3
		self.SetTime(8)

		for i := 0; i < len(self.PersonMgr); i++ {
			if self.room.Param1%10 == 1 { //! 自由抢庄
				self.PersonMgr[i].Card = self.Card.Deal(5)

				var msg Msg_GameNiuNiuJX_Card
				msg.Card = make([]int, 5)
				person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
				if person != nil {
					person.SendMsg("gamempqzcard", &msg)
				}
			} else { //! 看牌抢庄
				//	card := self.Card.Deal(5 - len(self.PersonMgr[i].Card))
				//	self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, card...)
				card := self.PersonMgr[i].CardZB[4:]
				//card = append(card, self.PersonMgr[i].CardZB[4])
				self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, card...)

				var msg Msg_GameNiuNiuJX_Card
				msg.Card = card
				person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
				if person != nil {
					person.SendMsg("gamempqzcard", &msg)
				}
			}
		}
	}

	self.room.flush()
}

//! 结算
func (self *Game_GoldNNRoom) OnEnd() {
	self.room.SetBegin(false)
	self.State = 0
	self.Time = 0

	var dealer *Game_GoldNNRoom_Person = nil
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Dealer {
			dealer = self.PersonMgr[i]
			break
		}
	}

	lst := make([]*Game_GoldNNRoom_Person, 0)
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Ready = false
		if self.PersonMgr[i].Uid != dealer.Uid {
			lst = append(lst, self.PersonMgr[i])
		}
	}

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
			bs := GetGoldNiuNiuBS(dealer.CT)
			score := lst[i].Bets * dealer.RobDeal * bs
			if lst[i].CurScore+lst[i].Score >= score*self.DF {
				lst[i].CurScore -= score * self.DF
				dealer.CurScore += score * self.DF
			} else {
				abs := lst[i].CurScore + lst[i].Score
				lst[i].CurScore -= abs
				dealer.CurScore += abs
			}
			lst[i].TZ = 0
		} else { //! 闲家赢
			bs := GetGoldNiuNiuBS(lst[i].CT)
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "闲家赢:", bs)
			score := lst[i].Bets * dealer.RobDeal * bs
			if dealer.CurScore+dealer.Score >= score*self.DF {
				dealer.CurScore -= score * self.DF
				lst[i].CurScore += score * self.DF
			} else {
				abs := dealer.CurScore + dealer.Score
				dealer.CurScore -= abs
				lst[i].CurScore += abs
			}
			lst[i].TZ = 0
		}
		lst[i].Score += lst[i].CurScore
	}
	dealer.Score += dealer.CurScore

	//! 记录
	var record staticfunc.Rec_Gold_Info
	record.Time = time.Now().Unix()
	record.GameType = self.room.Type

	self.State = 0
	self.SetTime(30)

	//! 发消息
	var msg Msg_GameMPQZ_End
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Ready = false

		var son Son_GameMPQZ_Info
		son.Uid = self.PersonMgr[i].Uid
		son.Bets = self.PersonMgr[i].Bets
		son.Card = self.PersonMgr[i].Card
		son.Dealer = self.PersonMgr[i].Dealer
		son.Score = self.PersonMgr[i].CurScore
		son.Total = self.PersonMgr[i].Score
		son.CT = self.PersonMgr[i].CT
		son.View = self.PersonMgr[i].View
		son.TZ = self.PersonMgr[i].TZ
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
	self.room.broadCastMsg("gamempqzend", &msg)

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
			_person := new(Game_GoldNNRoom_Person)
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

func (self *Game_GoldNNRoom) OnBye() {
}

func (self *Game_GoldNNRoom) OnExit(uid int64) {
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

func (self *Game_GoldNNRoom) getInfo(uid int64) *Msg_GameGoldNNRoom_Info {
	var msg Msg_GameGoldNNRoom_Info
	msg.Begin = self.room.Begin
	msg.State = self.State
	if self.Time != 0 {
		msg.Time = self.Time - time.Now().Unix()
	}
	msg.Info = make([]Son_GameGoldNNRoom_Info, 0)
	for _, value := range self.PersonMgr {
		if value.Uid == uid {
			msg.Bets = lib.HF_MinInt(lib.HF_MaxInt(value.MaxBets, 1), 5)
		}
		var son Son_GameGoldNNRoom_Info
		son.Uid = value.Uid
		son.Ready = value.Ready
		son.Bets = value.Bets
		son.Dealer = value.Dealer
		son.Total = value.Score
		son.Score = value.CurScore
		son.RobDeal = value.RobDeal
		son.TZ = value.TZ
		son.Open = value.Open
		son.View = value.View
		if self.room.Param1%10 == 1 { //! 自由抢庄
			if (value.Uid == uid && value.Open) || !msg.Begin || value.View || GetServer().IsAdmin(value.Uid, staticfunc.ADMIN_NIUNIU) {
				son.Card = value.Card
				son.CT = value.CT
			} else {
				son.Card = make([]int, len(value.Card))
				son.CT = 0
			}
		} else { //! 看牌抢庄
			if value.Uid == uid || value.View || !msg.Begin || GetServer().IsAdmin(value.Uid, staticfunc.ADMIN_NIUNIU) { //! 是自己或者亮牌了或者已经结束了
				son.Card = value.Card
				son.CT = value.CT
			} else {
				son.Card = make([]int, len(value.Card))
				son.CT = 0
			}
		}

		msg.Info = append(msg.Info, son)
	}
	return &msg
}

func (self *Game_GoldNNRoom) GetPerson(uid int64) *Game_GoldNNRoom_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

func (self *Game_GoldNNRoom) OnTime() {
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

func (self *Game_GoldNNRoom) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_GoldNNRoom) OnIsBets(uid int64) bool {
	return false
}

//! 同步总分
func (self *Game_GoldNNRoom) SendTotal() {
	var msg Msg_GameKWX_Total
	for i := 0; i < len(self.PersonMgr); i++ {
		self.room.Param[i] = self.PersonMgr[i].Score
		msg.Info = append(msg.Info, Son_GameKWX_Total{self.PersonMgr[i].Uid, self.PersonMgr[i].Score})
	}
	self.room.broadCastMsg("gamegoldtotal", &msg)
}

//! 设置时间
func (self *Game_GoldNNRoom) SetTime(t int) {
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
func (self *Game_GoldNNRoom) OnBalance() {
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
