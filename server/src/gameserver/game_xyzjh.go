package gameserver

import (
	"lib"
	"sort"
	"staticfunc"
	"time"
)

const GameXYZJH_WaitTime = 30

////!
type Msg_GameXYZJH_Info struct {
	Begin    bool                 `json:"begin"` //! 是否开始
	CurOp    int64                `json:"curop"` //! 当前操作的玩家
	Ready    []int64              `json:"ready"` //! 准备的人
	Deal     []Son_GameXYZJH_Deal `json:"deal"`  //! 抢庄的人
	Info     []Son_GameXYZJH_Info `json:"info"`
	Round    int                  `json:"round"`    //! 轮数
	Point    int                  `json:"point"`    //! 当前注
	Allpoint int                  `json:"allpoint"` //! 当前总注
	CurTime  int64                `json:"curtime"`  //! 当前时间
}
type Son_GameXYZJH_Info struct {
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Open    bool   `json:"open"`
	Discard bool   `json:"discard"`
	Lose    bool   `json:"lose"`
	Card    []int  `json:"card"`
	Bets    int    `json:"bets"`
	Allbets int    `json:"allbets"`
	Dealer  bool   `json:"dealer"`
	Score   int    `json:"score"`
	Baozi   int    `json:"baozi"`
	Total   int    `json:"total"`
	IsVip   bool   `json:"isvip"`
}
type Son_GameXYZJH_Deal struct {
	Uid int64 `json:"uid"`
	Ok  bool  `json:"ok"`
}

//! 结算
type Msg_GameXYZJH_End struct {
	Info     []Son_GameXYZJH_Info `json:"info"`
	Round    int                  `json:"round"`    //! 轮数
	Point    int                  `json:"point"`    //! 当前注
	Allpoint int                  `json:"allpoint"` //! 当前总注
}

//! 房间结束
type Msg_GameXYZJH_Bye struct {
	Info []Son_GameXYZJH_Bye `json:"info"`
}

type Msg_GameXYZJH_View struct {
	Uid  int64 `json:"uid"`
	Card []int `json:"card"`
}

type Msg_GameXYZJH_Discard struct {
	Uid   int64 `json:"uid"`
	Opuid int64 `json:"opuid"`
	Round int   `json:"round"`
}

//!换玩家发送时间
type Msg_GameXYZJH_WaitTime struct {
	Uid  int64 `json:"uid"`
	Time int64 `json:"time"`
}

//! 游戏下注
type Msg_GameXYZJH_Bets struct {
	Uid      int64 `json:"uid"`
	OpUid    int64 `json:"opuid"`
	Bets     int   `json:"bets"`
	Point    int   `json:"point"`
	Addpoint int   `json:"addpoint"`
	Allpoint int   `json:"allpoint"`
	Allbets  int   `json:"allbets"`
	Round    int   `json:"round"`
}

//! 玩家比牌
type Msg_GameXYZJH_Com struct {
	Uid      int64 `json:"uid"`
	Destuid  int64 `json:"destuid"`
	Win      bool  `json:"win"`
	Point    int   `json:"point"`
	Addpoint int   `json:"addpoint"`
	Allpoint int   `json:"allpoint"`
	Allbets  int   `json:"allbets"`
	Round    int   `json:"round"`
	Opuid    int64 `json:"opuid"`
	Card1    []int `json:"card1"`
	Card2    []int `json:"card2"`
}

type Son_GameXYZJH_Bye struct {
	Uid     int64 `json:"uid"`
	Win     int   `json:"win"`     //! 胜利次数
	Baozi   int   `json:"baozi"`   //! 豹子次数
	Shunjin int   `json:"shunjin"` //! 顺金次数
	Jinhua  int   `json:"jinhua"`  //! 金花次数
	Menpai  int   `json:"menpai"`
	Minpai  int   `json:"minpai"`
	Score   int   `json:"score"`
}

type Msg_GameXYZJH_ChangeCard struct {
	Uid  int64 `json:"uid"`
	Card []int `json:"card"`
}

////! 得到最后一张牌
//type Msg_GameXYZJH_Card struct {
//	Card int   `json:"card"`
//	All  []int `json:"all"`
//}

/////////////////////////////////////////////////////////
type Game_XYZJH_Person struct {
	Uid      int64  `json:"uid"`
	Name     string `json:"name"`
	Card     []int  `json:"card"`     //! 手牌
	Win      int    `json:"win"`      //! 胜利次数
	Baozi    int    `json:"baozi"`    //! 豹子次数
	Shunjin  int    `json:"shunjin"`  //! 顺金次数
	Jinhua   int    `json:"jinhua"`   //! 金花次数
	Menpai   int    `json:"menpai"`   //! 闷牌次数
	Minpai   int    `json:"minpai"`   //! 看牌跟注次数
	Score    int    `json:"score"`    //! 积分
	Bets     int    `json:"bets"`     //! 下注
	Allbets  int    `json:"Allbets"`  //! 下注
	Dealer   bool   `json:"dealer"`   //! 是否庄家
	Open     bool   `json:"open"`     //! 是否看牌了
	Discard  bool   `json:"discard"`  //! 是否弃牌了
	Lose     bool   `json:"lose"`     //! 是否比牌输了
	CurScore int    `json:"curscore"` //! 当前局的分数
	CurBaozi int    `json:"curbaozi"` //! 当前局豹子分数
	IsVip    bool   `json:"isvip"`    //! 	是否是作弊玩家

	_ct int //! 当前局牌型
	_cs int //! 当前局最大牌
}

type Game_XYZJH struct {
	Ready       []int64                        `json:"ready"`       //! 已经准备的人
	Bets        []int64                        `json:"bets"`        //! 已经下注的人
	CurStep     int                            `json:"curstep"`     //! 谁出牌
	Round       int                            `json:"round"`       //! 轮数
	Point       int                            `json:"point"`       //! 当前注
	Allpoint    int                            `json:"allpoint"`    //! 当前总注
	WaitEndTime int64                          `json:"waitendtime"` //!当前等待结束时间
	PersonMgr   []*Game_XYZJH_Person           `json:"personmgr"`
	record      *staticfunc.Rec_GameXYZJH_Info `json:"record"` //! 记录

	Timestart bool `json:"timestart"`

	room *Room
}

func NewGame_XYZJH() *Game_XYZJH {
	game := new(Game_XYZJH)
	game.Ready = make([]int64, 0)
	game.Bets = make([]int64, 0)
	game.PersonMgr = make([]*Game_XYZJH_Person, 0)
	//game.Deal = make([]Son_GameNiuNiu_Deal, 0)

	return game
}

func (self *Game_XYZJH) OnInit(room *Room) {
	self.room = room
}

func (self *Game_XYZJH) OnRobot(robot *lib.Robot) {

}

func (self *Game_XYZJH) OnSendInfo(person *Person) {
	person.SendMsg("gamexyzjhinfo", self.getInfo(person.Uid))
}

func (self *Game_XYZJH) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "gameready": //! 游戏准备
		self.GameReady(msg.Uid)
	case "gamebets": //! 下注
		self.GameBets(msg.V.(*Msg_GameBets).Uid, msg.V.(*Msg_GameBets).Bets)
	//case "gameview": //! 亮牌
	//	self.GameView(msg.Uid)
	case "gamecompare": //! 比牌
		self.GameCompare(msg.Uid, msg.V.(*Msg_GameCompare).Destuid)
	case "gamechangecard": //换牌
		if self.room.Type == 94 {
			self.ChangeCard(msg.Uid, msg.V.(*Msg_GameChangeCard).Card, msg.V.(*Msg_GameChangeCard).ChCard)
		}
	}
}

func (self *Game_XYZJH) GetPerson(uid int64) *Game_XYZJH_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}
	return nil
}

func (self *Game_XYZJH) ChangeCard(uid int64, card int, _card int) {
	person := GetPersonMgr().GetPerson(uid)
	_person := self.GetPerson(uid)
	for i := 0; i < len(self.PersonMgr); i++ {
		for j := 0; j < len(self.PersonMgr[i].Card); j++ {
			if self.PersonMgr[i].Card[j] == _card {
				person.SendErr("换牌重复！")
				return
			}
		}
	}
	find := false
	for i := 0; i < len(_person.Card); i++ {
		if _person.Card[i] == card {
			_person.Card[i] = _card
			find = true
		}
	}

	if !find {
		lib.GetLogMgr().Output(staticfunc.BUY_GOLD, "没找到要换的牌")
		return
	}

	var msg Msg_GameXYZJH_ChangeCard
	msg.Uid = _person.Uid
	msg.Card = _person.Card
	person.SendMsg("changecard", &msg)

}

func (self *Game_XYZJH) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)
	self.Point = 1
	self.Allpoint = 0
	self.Round = 0
	self.record = new(staticfunc.Rec_GameXYZJH_Info)

	if self.room.Param2/1000 == 1 {
		self.Timestart = true
	}

	//! 庄家的位置
	WinPos := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].CurScore > self.PersonMgr[WinPos].CurScore {
			WinPos = i
		}

	}

	for i := 0; i < len(self.room.Uid); i++ { //! 重新初始化人
		self.Allpoint++
		if i >= len(self.PersonMgr) {
			person := new(Game_XYZJH_Person)
			person.Uid = self.room.Uid[i]
			person.Name = self.room.Name[i]
			self.PersonMgr = append(self.PersonMgr, person)
		} else {
			self.PersonMgr[i]._ct = 0
			self.PersonMgr[i]._cs = 0
			self.PersonMgr[i].CurScore = 0
			self.PersonMgr[i].CurBaozi = 0
			self.PersonMgr[i].Lose = false
			self.PersonMgr[i].Open = false
			self.PersonMgr[i].Discard = false
			self.PersonMgr[i].Dealer = false
			if GetServer().IsAdmin(self.PersonMgr[i].Uid, staticfunc.ADMIN_SZP) {
				self.PersonMgr[i].IsVip = true
			} else {
				self.PersonMgr[i].IsVip = false
			}
		}

		self.PersonMgr[i].Bets = 1
		self.PersonMgr[i].Allbets = 1
	}

	//! 确定庄家   赢家当庄

	self.PersonMgr[WinPos].Dealer = true
	self.CurStep = WinPos

	self.NextPlayer()

	//! 发牌
	cardmgr := NewCard_ZJH()
	for i := 0; i < len(self.PersonMgr); i++ {
		//if i == 0 {
		//	self.PersonMgr[i].Card = []int{23, 33, 43}
		//} else if i == 1 {
		//	self.PersonMgr[i].Card = []int{71, 102, 113}
		//} else if i == 2 {
		//	self.PersonMgr[i].Card = []int{72, 101, 114}
		//} else if i == 3 {
		//	self.PersonMgr[i].Card = []int{51, 103, 111}
		//} else {
		//	self.PersonMgr[i].Card = []int{52, 104, 112}
		//}
		self.PersonMgr[i].Card = cardmgr.Deal(3)
		sort.Ints(self.PersonMgr[i].Card)

		//! 记录
		var rc_person staticfunc.Son_Rec_GameXYZJH_Person
		rc_person.Uid = self.PersonMgr[i].Uid
		rc_person.Name = self.room.GetName(rc_person.Uid)
		rc_person.Head = self.room.GetHead(rc_person.Uid)
		lib.HF_DeepCopy(&rc_person.Card, &self.PersonMgr[i].Card)
		self.record.Person = append(self.record.Person, rc_person)
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gamexyzjhbegin", self.getInfo(person.Uid))
	}

	self.room.flush()

	if self.room.Param2/1000 == 1 {
		self.WaitEndTime = time.Now().Unix() + GameXYZJH_WaitTime
		var msg Msg_GameXYZJH_WaitTime
		msg.Uid = self.PersonMgr[self.CurStep].Uid
		msg.Time = GameXYZJH_WaitTime
		for j := 0; j < len(self.PersonMgr); j++ {
			person := GetPersonMgr().GetPerson(self.PersonMgr[j].Uid)
			if person == nil {
				continue
			}
			person.SendMsg("gamewaittime", &msg)
		}
	}
}

////! 比牌
func (self *Game_XYZJH) GameCompare(uid int64, destuid int64) {
	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	if self.room.Param1/100%10 == 1 && self.Round < 1 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "未达到比牌轮数")
		return
	} else if self.room.Param1/100%10 == 2 && self.Round < 2 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "未达到比牌轮数")
		return
	} else if self.room.Param1/100%10 == 3 && self.Round < 3 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "未达到比牌轮数")
		return
	}

	if uid == destuid {
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

	/*if !self.PersonMgr[player2].Open {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能和没有看牌的比牌")
		return
	}*/

	var win int
	if self.room.Param1%10 == 2 { //花色比较
		win = ZjhCardCompare1(self.PersonMgr[player1].Card, self.PersonMgr[player2].Card)
	} else if self.room.Param1%10 == 3 { //全比较
		win = ZjhCardCompare2(self.PersonMgr[player1].Card, self.PersonMgr[player2].Card)
	} else { //大小比较
		win = ZjhCardCompare(self.PersonMgr[player1].Card, self.PersonMgr[player2].Card)
	}

	if win == 0 {
		self.PersonMgr[player2].Lose = true
	} else {
		self.PersonMgr[player1].Lose = true
	}

	addpoint := self.Point
	if self.room.Param2/10%10 == 1 { //比牌双倍
		addpoint = self.Point * 2
	}
	if self.PersonMgr[player1].Open {
		addpoint *= 2
	}
	self.PersonMgr[player1].Bets = addpoint
	self.PersonMgr[player1].Allbets += addpoint

	self.Allpoint += addpoint

	self.NextPlayer()
	var msg Msg_GameXYZJH_Com
	msg.Uid = uid
	msg.Destuid = destuid
	msg.Win = (win == 0)
	msg.Point = self.Point
	msg.Addpoint = addpoint
	msg.Allpoint = self.Allpoint
	msg.Allbets = self.PersonMgr[player1].Allbets
	msg.Opuid = self.room.Uid[self.CurStep]
	msg.Round = self.Round

	for i := 0; i < len(self.PersonMgr); i++ {
		if uid == self.PersonMgr[i].Uid {
			msg.Card1 = self.PersonMgr[player1].Card

			if self.PersonMgr[player2].Open || msg.Win {
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

	if self.record != nil {
		self.record.Step = append(self.record.Step, staticfunc.Son_Rec_GameXYZJH_Step{uid, destuid})
	}

	self.room.flush()

	count := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if !self.PersonMgr[i].Discard && !self.PersonMgr[i].Lose {
			count++
		}
	}

	if count == 1 {
		self.OnEnd()
	} else {
		self.GameView(uid)

		if self.room.Param2/1000 == 1 {
			self.WaitEndTime = time.Now().Unix() + GameXYZJH_WaitTime
			var msg Msg_GameXYZJH_WaitTime
			msg.Uid = self.PersonMgr[self.CurStep].Uid
			msg.Time = GameXYZJH_WaitTime
			for j := 0; j < len(self.PersonMgr); j++ {
				person := GetPersonMgr().GetPerson(self.PersonMgr[j].Uid)
				if person == nil {
					continue
				}
				person.SendMsg("gamewaittime", &msg)
			}
		}
	}

}

//! 看牌
func (self *Game_XYZJH) GameView(uid int64) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Open {
				return
			}

			self.PersonMgr[i].Open = true
			card := self.PersonMgr[i].Card

			for i := 0; i < len(self.PersonMgr); i++ {
				var msg Msg_GameXYZJH_View
				msg.Uid = uid
				for j := 0; j < 3; j++ {
					if self.PersonMgr[i].Uid == uid || GetServer().IsAdmin(self.PersonMgr[i].Uid, staticfunc.ADMIN_SZP) {
						msg.Card = append(msg.Card, card[j])
					} else {
						msg.Card = append(msg.Card, 0)
					}
				}
				self.room.SendMsg(self.PersonMgr[i].Uid, "gameview", &msg)
			}

			break
		}
	}

	self.room.flush()
}

func (self *Game_XYZJH) NextPlayer() {
	for i := 0; i < len(self.PersonMgr); i++ {
		self.CurStep++
		self.CurStep %= len(self.PersonMgr)
		if self.PersonMgr[self.CurStep].Dealer {
			self.Round++
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "Round:", self.Round)
			//break
		}
		if !self.PersonMgr[self.CurStep].Lose && !self.PersonMgr[self.CurStep].Discard {
			break
		}
	}

	//	for i := 0; i < len(self.PersonMgr); i++ { //! 重新初始化人
	//		if self.PersonMgr[i].Dealer && self.CurStep == i {
	//			self.Round++
	//			lib.GetLogMgr().Output(lib.LOG_DEBUG, "Round:", self.Round)
	//			break
	//		}
	//	}
}

//! 弃牌
func (self *Game_XYZJH) GameDiscard(uid int64) {

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Discard {
				return
			}

			self.PersonMgr[i].Discard = true
		}
	}

	if self.room.Uid[self.CurStep] == uid {
		self.NextPlayer()
	}

	var msg Msg_GameXYZJH_Discard
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

	if count == 1 {
		self.OnEnd()
	} else {
		if self.room.Param2/1000 == 1 {
			self.WaitEndTime = time.Now().Unix() + GameXYZJH_WaitTime
			var msg Msg_GameXYZJH_WaitTime
			msg.Uid = self.PersonMgr[self.CurStep].Uid
			msg.Time = GameXYZJH_WaitTime
			for j := 0; j < len(self.PersonMgr); j++ {
				person := GetPersonMgr().GetPerson(self.PersonMgr[j].Uid)
				if person == nil {
					continue
				}
				person.SendMsg("gamewaittime", &msg)
			}
		}
	}

	self.room.flush()

}

////! 准备
func (self *Game_XYZJH) GameReady(uid int64) {
	if self.room.IsBye() {
		return
	}

	if self.room.Begin { //! 已经开始了不允许准备
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经开始了，不能准备")
		return
	}

	for i := 0; i < len(self.Ready); i++ {
		if self.Ready[i] == uid {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "同一个玩家准备")
			return
		}
	}

	self.Ready = append(self.Ready, uid)

	if len(self.Ready) == len(self.room.Uid) && len(self.Ready) >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 准备的人数达到游戏最小人数
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
		self.OnBegin()
		return
	}

	var msg staticfunc.Msg_Uid
	msg.Uid = uid
	self.room.broadCastMsg("gameready", &msg)

	self.room.flush()
}

////! 下注 -1表示弃牌 0表示看牌 1表示跟住 2-5表示加注
func (self *Game_XYZJH) GameBets(uid int64, bets int) {
	//	//if self.room.IsBye() {
	//	//	return
	//	//}
	if bets < -1 || bets > 5 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "参数不对")
		return
	}

	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	if bets == 0 {
		if self.room.Param1/1000 == 1 && self.Round < 2 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "还不能看牌")
			return
		} else if self.room.Param1/1000 == 2 && self.Round < 3 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "还不能看牌")
			return
		} else if self.room.Param1/1000 == 3 && self.Round < 5 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "还不能看牌")
			return
		}

		self.GameView(uid)
		if self.record != nil {
			self.record.Step = append(self.record.Step, staticfunc.Son_Rec_GameXYZJH_Step{uid, int64(bets)})
		}
		return
	} else if bets == -1 {
		//		if self.Round < 1 {
		//			lib.GetLogMgr().Output(lib.LOG_DEBUG, "还不能弃牌")
		//			return
		//		}
		self.GameDiscard(uid)
		if self.record != nil {
			self.record.Step = append(self.record.Step, staticfunc.Son_Rec_GameXYZJH_Step{uid, int64(bets)})
		}
		return
	}

	if self.room.Uid[self.CurStep] != uid {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "当前玩家不能操作")
		return
	}

	//	if self.Point >= 5 {
	//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "加注无效，已经是最大值")
	//		return
	//	}

	if bets >= 2 && bets <= 5 {
		if self.Point < bets {
			self.Point = bets
		} else {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "当前注大于或等于加注数，加注无效")
			return
		}
	}

	addpoint := self.Point
	//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "圈数bet:", self.Round)
	self.NextPlayer()
	//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "圈数:", self.Round)

	allbets := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Open {
				addpoint *= 2
				self.PersonMgr[i].Minpai++
			} else {
				self.PersonMgr[i].Menpai++
			}
			self.PersonMgr[i].Bets = addpoint
			self.PersonMgr[i].Allbets += addpoint
			allbets = self.PersonMgr[i].Allbets
		}
	}

	self.Allpoint += addpoint

	var msg Msg_GameXYZJH_Bets
	msg.Uid = uid
	msg.OpUid = self.room.Uid[self.CurStep]
	msg.Bets = bets
	msg.Point = self.Point
	msg.Addpoint = addpoint
	msg.Allpoint = self.Allpoint
	msg.Allbets = allbets
	msg.Round = self.Round
	self.room.broadCastMsg("gamebets", &msg)

	if self.record != nil {
		self.record.Step = append(self.record.Step, staticfunc.Son_Rec_GameXYZJH_Step{uid, int64(bets)})
	}

	if self.Round >= 5 && self.room.Param1/10%10 == 1 {
		self.OnEnd()
		return
	} else if self.Round >= 10 && self.room.Param1/10%10 == 2 {
		self.OnEnd()
		return
	} else if self.Round >= 15 && self.room.Param1/10%10 == 3 {
		self.OnEnd()
		return
	}

	self.room.flush()

	if self.room.Param2/1000 == 1 {
		self.WaitEndTime = time.Now().Unix() + GameXYZJH_WaitTime
		var msg Msg_GameXYZJH_WaitTime
		msg.Uid = self.PersonMgr[self.CurStep].Uid
		msg.Time = GameXYZJH_WaitTime
		for j := 0; j < len(self.PersonMgr); j++ {
			person := GetPersonMgr().GetPerson(self.PersonMgr[j].Uid)
			if person == nil {
				continue
			}
			person.SendMsg("gamewaittime", &msg)
		}
	}
}

////! 结算
func (self *Game_XYZJH) OnEnd() {
	if self.room.Param2/1000 == 1 {
		self.Timestart = false
	}
	allpoint := 0
	var dealerpos int //庄家的位置
	emplst := make([]*Game_XYZJH_Person, 0)
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Dealer {
			dealerpos = i
		}

		if !self.PersonMgr[i].Discard && !self.PersonMgr[i].Lose {
			emplst = append(emplst, self.PersonMgr[i])
		} else {
			allpoint += self.PersonMgr[i].Allbets
		}
	}

	oldwin := make([]int, 0) //只能是一个人
	oldwin = append(oldwin, 0)
	for i := 1; i < len(emplst); i++ {
		var win int
		if self.room.Param1%10 == 2 { //比较花色
			win = ZjhCardCompare1(emplst[oldwin[0]].Card, emplst[i].Card)
		} else if self.room.Param1%10 == 3 { //全比较
			win = ZjhCardCompare2(emplst[oldwin[0]].Card, emplst[i].Card)
		} else { //比较大小
			win = ZjhCardCompare(emplst[oldwin[0]].Card, emplst[i].Card)
		}
		if win == 2 { //平局情况，从庄家开始，自小到大，末家赢
			if (oldwin[0] < dealerpos && i < dealerpos) || (oldwin[0] > dealerpos && i > dealerpos) {
				if oldwin[0] < i {
					win = 1
				} else {
					win = 0
				}
			} else if oldwin[0] < dealerpos && i > dealerpos {
				win = 0
			} else if oldwin[0] > dealerpos && i < dealerpos {
				win = 1
			}
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
		} /*else { //! 平
			oldwin = append(oldwin, i)
		}*/
	}

	self.room.SetBegin(false)

	self.Ready = make([]int64, 0)
	self.Bets = make([]int64, 0)

	//	sy := allpoint % len(oldwin)
	//	//平牌的情况随机一个人多一点
	//	if sy > 0 {
	//		emplst[oldwin[lib.HF_GetRandom(len(oldwin))]].CurScore += sy
	//	}
	for i := 0; i < len(self.PersonMgr); i++ {
		baozijiangli := 0
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

			if self.room.Param2%10 == 1 { //豹子奖励
				baozijiangli = 5
				self.PersonMgr[i].CurBaozi += 5
			}

			//lib.GetLogMgr().Output(lib.LOG_DEBUG, "当前豹子模式", self.room.Param2%10, baozijiangli)

			if baozijiangli > 0 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "22222")
				for j := 0; j < len(self.PersonMgr); j++ {
					if self.PersonMgr[j].Uid == self.PersonMgr[i].Uid {
						self.PersonMgr[j].CurScore += baozijiangli * (len(self.PersonMgr) - 1)
						//self.PersonMgr[j].Score += baozijiangli * (len(self.PersonMgr) - 1)
					} else {
						self.PersonMgr[j].CurScore -= baozijiangli
						//elf.PersonMgr[j].Score -= baozijiangli
					}
					lib.GetLogMgr().Output(lib.LOG_DEBUG, self.PersonMgr[j].CurScore, baozijiangli)
				}
			}

		case 500:
			self.PersonMgr[i].Shunjin++
		case 400:
			self.PersonMgr[i].Jinhua++
		}

	}

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Score += self.PersonMgr[i].CurScore

		if self.record != nil {
			for j := 0; j < len(self.record.Person); j++ {
				if self.record.Person[j].Uid == self.PersonMgr[i].Uid {
					self.record.Person[j].Score = self.PersonMgr[i].CurScore
					self.record.Person[j].Total = self.PersonMgr[i].Score
					break
				}
			}
		}
	}

	//! 记录
	if self.record != nil {
		self.record.Roomid = self.room.Id*100 + self.room.Step
		self.record.Time = time.Now().Unix()
		self.record.MaxStep = self.room.MaxStep
		self.room.AddRecord(lib.HF_JtoA(self.record))
	}

	//! 发消息
	var msg Msg_GameXYZJH_End
	msg.Point = self.Point
	msg.Round = self.Round
	msg.Allpoint = self.Allpoint
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameXYZJH_Info

		son.Uid = self.PersonMgr[i].Uid
		son.Name = self.PersonMgr[i].Name
		son.Bets = self.PersonMgr[i].Bets
		son.Card = self.PersonMgr[i].Card
		son.Discard = self.PersonMgr[i].Discard
		son.Dealer = self.PersonMgr[i].Dealer
		son.Score = self.PersonMgr[i].CurScore
		son.Total = self.PersonMgr[i].Score
		son.Baozi = self.PersonMgr[i].CurBaozi
		msg.Info = append(msg.Info, son)
	}
	self.room.broadCastMsg("gamexyzjhend", &msg)

	if self.room.IsBye() {
		self.OnBye()
		self.room.Bye()
		return
	}

	self.room.flush()
}

func (self *Game_XYZJH) OnBye() {
	if self.room.Param2/100%10 == 1 && !self.room.IsBye() { //解散算分
		self.OnEnd()
	}

	info := make([]staticfunc.JS_CreateRoomMem, 0)
	var msg Msg_GameXYZJH_Bye
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameXYZJH_Bye
		son.Uid = self.PersonMgr[i].Uid
		son.Win = self.PersonMgr[i].Win
		son.Baozi = self.PersonMgr[i].Baozi
		son.Shunjin = self.PersonMgr[i].Shunjin
		son.Jinhua = self.PersonMgr[i].Jinhua
		son.Menpai = self.PersonMgr[i].Menpai
		son.Minpai = self.PersonMgr[i].Minpai
		son.Score = self.PersonMgr[i].Score
		msg.Info = append(msg.Info, son)
		info = append(info, staticfunc.JS_CreateRoomMem{son.Uid, "", "", son.Score})
	}
	self.room.broadCastMsg("gamexyzjhbye", &msg)

	self.room.ClubResult(info)
}

func (self *Game_XYZJH) OnExit(uid int64) {
	find := false
	for i := 0; i < len(self.Ready); i++ {
		if self.Ready[i] == uid {
			copy(self.Ready[i:], self.Ready[i+1:])
			self.Ready = self.Ready[:len(self.Ready)-1]
			find = true
			break
		}
	}

	if !find {
		if len(self.Ready) == len(self.room.Uid) && len(self.Ready) >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 准备的人数达到游戏最小人数
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
			self.OnBegin()
			return
		}
	}
}

func (self *Game_XYZJH) getInfo(uid int64) *Msg_GameXYZJH_Info {
	var msg Msg_GameXYZJH_Info
	msg.Begin = self.room.Begin
	msg.CurOp = self.room.Uid[self.CurStep]
	msg.Ready = make([]int64, 0)
	msg.Round = self.Round
	msg.Point = self.Point
	msg.CurTime = self.WaitEndTime - time.Now().Unix()
	msg.Allpoint = self.Allpoint
	msg.Info = make([]Son_GameXYZJH_Info, 0)
	if !msg.Begin { //! 没有开始,看哪些人已准备
		msg.Ready = self.Ready
	}
	for _, value := range self.PersonMgr {
		var son Son_GameXYZJH_Info
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
		if GetServer().IsAdmin(value.Uid, staticfunc.ADMIN_SZP) {
			son.IsVip = true
		} else {
			son.IsVip = false
		}

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

func (self *Game_XYZJH) OnTime() {
	if self.room.Param2/1000 == 1 {
		if time.Now().Unix() >= self.WaitEndTime && self.Timestart {
			self.GameBets(self.PersonMgr[self.CurStep].Uid, -1)
			//			self.WaitEndTime = time.Now().Unix() + GameXYZJH_WaitTime
		}
	}
}

func (self *Game_XYZJH) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_XYZJH) OnIsBets(uid int64) bool {
	return false
}

//! 结算所有人
func (self *Game_XYZJH) OnBalance() {
}
