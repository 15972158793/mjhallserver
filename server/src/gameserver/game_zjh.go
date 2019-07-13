package gameserver

import (
	"lib"
	"sort"
	"staticfunc"
	"time"
)

/*
type = 80
param1
个位 封顶分数
十位 豹子奖励
百位 庄的模式
千位 1 顺子>金花
万位	1 杂色235打豹子
*/

////!
type Msg_GameZJH_Info struct {
	Begin    bool               `json:"begin"` //! 是否开始
	CurOp    int64              `json:"curop"` //! 当前操作的玩家
	Ready    []int64            `json:"ready"` //! 准备的人
	Deal     []Son_GameZJH_Deal `json:"deal"`  //! 抢庄的人
	Info     []Son_GameZJH_Info `json:"info"`
	Round    int                `json:"round"`    //! 轮数
	Point    int                `json:"point"`    //! 当前注
	Allpoint int                `json:"allpoint"` //! 当前总注
}
type Son_GameZJH_Info struct {
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
}
type Son_GameZJH_Deal struct {
	Uid int64 `json:"uid"`
	Ok  bool  `json:"ok"`
}

//! 结算
type Msg_GameZJH_End struct {
	Info     []Son_GameZJH_Info `json:"info"`
	Round    int                `json:"round"`    //! 轮数
	Point    int                `json:"point"`    //! 当前注
	Allpoint int                `json:"allpoint"` //! 当前总注
}

//! 房间结束
type Msg_GameZJH_Bye struct {
	Info []Son_GameZJH_Bye `json:"info"`
}

type Msg_GameZJH_View struct {
	Uid  int64 `json:"uid"`
	Card []int `json:"card"`
	Type int   `json:"type"`
}

type Msg_GameZJH_Discard struct {
	Uid   int64 `json:"uid"`
	Opuid int64 `json:"opuid"`
	Round int   `json:"round"`
}

//! 游戏下注
type Msg_GameZJH_Bets struct {
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
type Msg_GameZJH_Com struct {
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

type Son_GameZJH_Bye struct {
	Uid     int64 `json:"uid"`
	Win     int   `json:"win"`     //! 胜利次数
	Baozi   int   `json:"baozi"`   //! 豹子次数
	Shunjin int   `json:"shunjin"` //! 顺金次数
	Jinhua  int   `json:"jinhua"`  //! 金花次数
	Menpai  int   `json:"menpai"`
	Minpai  int   `json:"minpai"`
	Score   int   `json:"score"`
}

////! 得到最后一张牌
//type Msg_GameZJH_Card struct {
//	Card int   `json:"card"`
//	All  []int `json:"all"`
//}

/////////////////////////////////////////////////////////
type Game_ZJH_Person struct {
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

	_ct int //! 当前局牌型
	_cs int //! 当前局最大牌
}

type Game_ZJH struct {
	Ready     []int64                      `json:"ready"`    //! 已经准备的人
	Bets      []int64                      `json:"bets"`     //! 已经下注的人
	CurStep   int                          `json:"curstep"`  //! 谁出牌
	Round     int                          `json:"round"`    //! 轮数
	Point     int                          `json:"point"`    //! 当前注
	Allpoint  int                          `json:"allpoint"` //! 当前总注
	PersonMgr []*Game_ZJH_Person           `json:"personmgr"`
	Record    *staticfunc.Rec_GameZJH_Info `json:"record"` //! 记录
	FirstStep int                          `json:"firststep"`

	room *Room
}

func NewGame_ZJH() *Game_ZJH {
	game := new(Game_ZJH)
	game.Ready = make([]int64, 0)
	game.Bets = make([]int64, 0)
	game.PersonMgr = make([]*Game_ZJH_Person, 0)
	//game.Deal = make([]Son_GameNiuNiu_Deal, 0)

	return game
}

func (self *Game_ZJH) OnInit(room *Room) {
	self.room = room
}

func (self *Game_ZJH) OnRobot(robot *lib.Robot) {

}

func (self *Game_ZJH) OnSendInfo(person *Person) {
	person.SendMsg("gamezjhinfo", self.getInfo(person.Uid))
}

func (self *Game_ZJH) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "gameready": //! 游戏准备
		self.GameReady(msg.Uid)
	case "gamebets": //! 下注
		self.GameBets(msg.V.(*Msg_GameBets).Uid, msg.V.(*Msg_GameBets).Bets)
	//case "gameview": //! 亮牌
	//	self.GameView(msg.Uid)
	case "gamecompare": //! 比牌
		self.GameCompare(msg.Uid, msg.V.(*Msg_GameCompare).Destuid)
	}
}

func (self *Game_ZJH) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)
	self.Point = 1
	self.Allpoint = 0
	self.Round = 0
	self.Record = new(staticfunc.Rec_GameZJH_Info)

	//! 庄家的位置
	DealerPos := -1
	WinPos := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].CurScore > self.PersonMgr[WinPos].CurScore {
			WinPos = i
		}
		if self.PersonMgr[i].Dealer {
			DealerPos = i
		}
	}

	for i := 0; i < len(self.room.Uid); i++ { //! 重新初始化人
		self.Allpoint++
		if i >= len(self.PersonMgr) {
			person := new(Game_ZJH_Person)
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
		}

		self.PersonMgr[i].Bets = 1
		self.PersonMgr[i].Allbets = 1
	}

	//! 确定庄家
	if self.room.Param1/100%10 == 0 {
		if DealerPos+1 >= len(self.PersonMgr) {
			DealerPos = -1
		}
		self.PersonMgr[DealerPos+1].Dealer = true
		self.CurStep = DealerPos + 1
		self.FirstStep = DealerPos + 1
	} else {
		self.PersonMgr[WinPos].Dealer = true
		self.CurStep = WinPos
		self.FirstStep = WinPos
	}
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
		var rc_person staticfunc.Son_Rec_GameZJH_Person
		rc_person.Uid = self.PersonMgr[i].Uid
		rc_person.Name = self.room.GetName(rc_person.Uid)
		rc_person.Head = self.room.GetHead(rc_person.Uid)
		rc_person.Dealer = self.PersonMgr[i].Dealer
		lib.HF_DeepCopy(&rc_person.Card, &self.PersonMgr[i].Card)
		self.Record.Person = append(self.Record.Person, rc_person)
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gamezjhbegin", self.getInfo(person.Uid))
	}

	self.room.flush()
}

////! 比牌
func (self *Game_ZJH) GameCompare(uid int64, destuid int64) {
	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
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
	if self.room.Type == 80 {
		win = ZjhCardCompare3(self.PersonMgr[player1].Card, self.PersonMgr[player2].Card, self.room.Param1/10000, self.room.Param1/1000%10)
	} else {
		win = ZjhCardCompare3(self.PersonMgr[player1].Card, self.PersonMgr[player2].Card, 0, 0)
	}

	if win == 0 {
		self.PersonMgr[player2].Lose = true
	} else {
		self.PersonMgr[player1].Lose = true
	}

	addpoint := self.Point * 2

	if self.PersonMgr[player1].Open {
		addpoint *= 2
	}

	self.PersonMgr[player1].Bets = addpoint
	self.PersonMgr[player1].Allbets += addpoint

	self.Allpoint += addpoint

	self.NextPlayer()

	var msg Msg_GameZJH_Com
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

	if self.Record != nil {
		self.Record.Step = append(self.Record.Step, staticfunc.Son_Rec_GameJZH_Step{uid, destuid, win})
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
	}
}

//! 看牌
func (self *Game_ZJH) GameView(uid int64) {
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
				for j := 0; j < 3; j++ {
					if self.PersonMgr[i].Uid == uid || GetServer().IsAdmin(self.PersonMgr[i].Uid, staticfunc.ADMIN_SZP) {
						msg.Card = append(msg.Card, card[j])
					} else {
						msg.Card = append(msg.Card, 0)
					}
				}
				self.room.SendMsg(self.PersonMgr[i].Uid, "gameview", &msg)
			}

			if self.Record != nil {
				self.Record.Step = append(self.Record.Step, staticfunc.Son_Rec_GameJZH_Step{uid, int64(0), self.PersonMgr[i].Allbets})
			}

			break
		}
	}

	self.room.flush()
}

func (self *Game_ZJH) NextPlayer() {
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
func (self *Game_ZJH) GameDiscard(uid int64) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Discard {
				return
			}

			self.PersonMgr[i].Discard = true

			if self.Record != nil {
				self.Record.Step = append(self.Record.Step, staticfunc.Son_Rec_GameJZH_Step{uid, int64(-1), self.PersonMgr[i].Allbets})
			}

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

	if count == 1 {
		self.OnEnd()
	}

	self.room.flush()

}

////! 准备
func (self *Game_ZJH) GameReady(uid int64) {
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

////! 下注 -1表示弃牌 0表示看牌 1表示跟住 2表示加注
func (self *Game_ZJH) GameBets(uid int64, bets int) {
	//	//if self.room.IsBye() {
	//	//	return
	//	//}

	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	if bets == 0 {
		if self.Round < 1 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "还不能看牌")
			return
		}
		self.GameView(uid)
		return
	} else if bets == -1 {
		if self.Round < 1 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "还不能弃牌")
			return
		}
		self.GameDiscard(uid)
		return
	}

	if self.room.Uid[self.CurStep] != uid {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "当前玩家不能操作")
		return
	}

	if bets == 2 {
		self.Point = 2
	}

	addpoint := self.Point

	self.NextPlayer()

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

			if self.Record != nil {
				self.Record.Step = append(self.Record.Step, staticfunc.Son_Rec_GameJZH_Step{uid, int64(bets), self.PersonMgr[i].Allbets})
			}

			break
		}
	}

	self.Allpoint += addpoint

	var msg Msg_GameZJH_Bets
	msg.Uid = uid
	msg.OpUid = self.room.Uid[self.CurStep]
	msg.Bets = bets
	msg.Point = self.Point
	msg.Addpoint = addpoint
	msg.Allpoint = self.Allpoint
	msg.Allbets = allbets
	msg.Round = self.Round
	self.room.broadCastMsg("gamebets", &msg)

	if self.Allpoint >= 100 && self.room.Param1%10 == 1 {
		self.OnEnd()
		return
	} else if self.Allpoint >= 50 && self.room.Param1%10 == 0 {
		self.OnEnd()
		return
	}

	self.room.flush()
}

////! 结算
func (self *Game_ZJH) OnEnd() {
	allpoint := 0
	emplst := make([]*Game_ZJH_Person, 0)
	for i := 0; i < len(self.PersonMgr); i++ {
		if !self.PersonMgr[i].Discard && !self.PersonMgr[i].Lose {
			emplst = append(emplst, self.PersonMgr[i])
		} else {
			allpoint += self.PersonMgr[i].Allbets
		}
	}

	oldwin := make([]int, 0)
	oldwin = append(oldwin, 0)
	for i := 1; i < len(emplst); i++ {
		var win int
		if self.room.Type == 80 {
			win = ZjhCardCompare3(emplst[oldwin[0]].Card, emplst[i].Card, self.room.Param1/10000, self.room.Param1/1000%10)
		} else {
			win = ZjhCardCompare3(emplst[oldwin[0]].Card, emplst[i].Card, 0, 0)
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

	self.room.SetBegin(false)

	self.Ready = make([]int64, 0)
	self.Bets = make([]int64, 0)

	sy := allpoint % len(oldwin)
	if sy > 0 {
		emplst[oldwin[lib.HF_GetRandom(len(oldwin))]].CurScore += sy
	}
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

			if self.room.Param1/10%10 == 1 {
				baozijiangli = 5
				self.PersonMgr[i].CurBaozi += 5
			} else if self.room.Param1/10%10 == 2 {
				baozijiangli = 10
				self.PersonMgr[i].CurBaozi += 10
			}

			lib.GetLogMgr().Output(lib.LOG_DEBUG, "当前豹子模式", self.room.Param1/10%10, baozijiangli)

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

	//! 记录
	if self.Record != nil {
		self.Record.Roomid = self.room.Id*100 + self.room.Step
		self.Record.Time = time.Now().Unix()
		self.Record.MaxStep = self.room.MaxStep
		self.room.AddRecord(lib.HF_JtoA(self.Record))
	}

	//! 发消息
	agentinfo := make([]staticfunc.JS_CreateRoomMem, 0)
	var msg Msg_GameZJH_End
	msg.Point = self.Point
	msg.Round = self.Round
	msg.Allpoint = self.Allpoint
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameZJH_Info

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
		agentinfo = append(agentinfo, staticfunc.JS_CreateRoomMem{son.Uid, "", "", son.Total})
	}
	self.room.broadCastMsg("gameniuniuend", &msg)
	self.room.AgentResult(agentinfo)

	if self.room.IsBye() {
		self.OnBye()
		self.room.Bye()
		return
	}

	self.room.flush()
}

func (self *Game_ZJH) OnBye() {
	info := make([]staticfunc.JS_CreateRoomMem, 0)
	var msg Msg_GameZJH_Bye
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameZJH_Bye
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

		GetServer().SqlScoreLog(self.PersonMgr[i].Uid, self.room.GetName(self.PersonMgr[i].Uid), self.room.GetHead(self.PersonMgr[i].Uid), self.room.Type, self.room.Id, self.PersonMgr[i].Score)
	}
	self.room.broadCastMsg("gameniuniubye", &msg)

	self.room.ClubResult(info)
}

func (self *Game_ZJH) OnExit(uid int64) {
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

func (self *Game_ZJH) getInfo(uid int64) *Msg_GameZJH_Info {
	var msg Msg_GameZJH_Info
	msg.Begin = self.room.Begin
	msg.CurOp = self.room.Uid[self.CurStep]
	msg.Ready = make([]int64, 0)
	msg.Round = self.Round
	msg.Point = self.Point
	msg.Allpoint = self.Allpoint
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

func (self *Game_ZJH) OnTime() {

}

func (self *Game_ZJH) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_ZJH) OnIsBets(uid int64) bool {
	return false
}

//! 结算所有人
func (self *Game_ZJH) OnBalance() {
}
