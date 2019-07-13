package gameserver

import (
	"lib"
	"math/rand"
	"staticfunc"
	"time"
)

//! param1%10 = 0大牌九  1小牌九
//! param1/10%10 = 0抢庄 1轮庄 2霸王庄
//! param1/100%10 = 0两道  1三道
//! param1/1000%10 0每次选分 >0固定分

//! param2%10 = 炸弹
//! param2/10%10 = 地九娘娘
//! param2/100%10 = 鬼子
//! param2/1000%10 = 天王九

const (
	TYPE_PTJ_DXJ = iota
	TYPE_PTJ_DEAL
	TYPE_PTJ_SD
	TYPE_PTJ_BET
	TYPE_PTJ_ZD
	TYPE_PTJ_DJNN
	TYPE_PTJ_GZ
	TYPE_PTJ_TWJ
)

//! 拼天九亮牌
type Msg_GamePTJ_View struct {
	View []int `json:"view"`
}

//! 拼天九亮牌
type Msg_GamePTJ_Send_View struct {
	Uid int64  `json:"uid"`
	CT  [2]int `json:"ct"`
}

//! 记录结构
type Rec_GamePTJ struct {
	Info    []Son_Rec_GamePTJ `json:"info"`
	Roomid  int               `json:"roomid"`
	MaxStep int               `json:"maxstep"`
	Param1  int               `json:"param1"`
	Param2  int               `json:"param2"`
	Time    int64             `json:"time"`
	View    []int             `json:"view"`
}
type Son_Rec_GamePTJ struct {
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Card    []int  `json:"card"`
	Bets    int    `json:"bets"`
	Dealer  bool   `json:"dealer"`
	Score   int    `json:"score"`
	Total   int    `json:"total"`
	RobDeal int    `json:"robdeal"`
	CT      [2]int `json:"ct"`
}

//!
type Msg_GamePTJ_Info struct {
	Begin bool               `json:"begin"` //! 是否开始
	Info  []Son_GamePTJ_Info `json:"info"`
	Card  []int              `json:"card"`
}
type Son_GamePTJ_Info struct {
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

//! 结算
type Msg_GamePTJ_End struct {
	Info []Son_GamePTJ_Info `json:"info"`
}

//! 房间结束
type Msg_GamePTJ_Bye struct {
	Info []Son_GamePTJ_Bye `json:"info"`
}
type Son_GamePTJ_Bye struct {
	Uid      int64 `json:"uid"`
	Win      int   `json:"win"`      //! 胜利次数
	MaxScore int   `json:"maxscore"` //! 最大分数
	Kill     int   `json:"kill"`     //! 通杀次数
	Dead     int   `json:"dead"`     //! 通赔次数
	MaxType  int   `json:"maxtype"`  //! 最大牌型
	Score    int   `json:"score"`
}

//! 得到最后一张牌
type Msg_GamePTJ_Card struct {
	SZ   [2]int `json:"sz"`
	Card []int  `json:"card"`
	View []int  `json:"view"`
}

///////////////////////////////////////////////////////
type Game_PTJ_Person struct {
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
	Win      int    `json:"win"`      //! 胜利次数
	MaxScore int    `json:"maxscore"` //! 最大分数
	Kill     int    `json:"kill"`     //! 通杀次数
	Dead     int    `json:"dead"`     //! 通赔次数
	MaxType  int    `json:"maxtype"`  //! 最大牌型
}

type Game_PTJ struct {
	PersonMgr []*Game_PTJ_Person `json:"personmgr"`
	Card      []int              `json:"card"`
	PJ        *CardMgr           `json:"pj"`
	Cache     []int              `json:"cache"`

	room *Room
}

func NewGame_PTJ() *Game_PTJ {
	game := new(Game_PTJ)
	game.PersonMgr = make([]*Game_PTJ_Person, 0)
	game.PJ = NewCard_TJ()
	game.Card = make([]int, 0)

	return game
}

func (self *Game_PTJ) GetParam(_type int) int {
	switch _type {
	case TYPE_PTJ_DXJ:
		return self.room.Param1 % 10
	case TYPE_PTJ_DEAL:
		return self.room.Param1 / 10 % 10
	case TYPE_PTJ_SD:
		return self.room.Param1 / 100 % 10
	case TYPE_PTJ_BET:
		bet := self.room.Param1 / 1000 % 10
		if bet == 0 {
			return 0
		} else if bet == 1 {
			return 1
		} else if bet == 2 {
			return 2
		} else if bet == 3 {
			return 5
		} else if bet == 4 {
			return 8
		} else if bet == 5 {
			return 10
		}
	case TYPE_PTJ_ZD:
		return self.room.Param2 % 10
	case TYPE_PTJ_DJNN:
		return self.room.Param2 / 10 % 10
	case TYPE_PTJ_GZ:
		return self.room.Param2 / 100 % 10
	case TYPE_PTJ_TWJ:
		return self.room.Param2 / 1000 % 10
	}
	return 0
}

func (self *Game_PTJ) OnInit(room *Room) {
	self.room = room
}

func (self *Game_PTJ) OnRobot(robot *lib.Robot) {

}

func (self *Game_PTJ) OnSendInfo(person *Person) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			person.SendMsg("gameptjinfo", self.getInfo(person.Uid))
			return
		}
	}

	_person := new(Game_PTJ_Person)
	_person.Uid = person.Uid
	self.PersonMgr = append(self.PersonMgr, _person)
	person.SendMsg("gameptjinfo", self.getInfo(person.Uid))
}

func (self *Game_PTJ) OnMsg(msg *RoomMsg) {
	switch msg.Head {
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

func (self *Game_PTJ) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)

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

	for i := 0; i < len(self.PersonMgr); i++ { //! 重新初始化人
		self.PersonMgr[i].CurScore = 0
		self.PersonMgr[i].Bets = 0
		self.PersonMgr[i].Dealer = false
		self.PersonMgr[i].View = false
		self.PersonMgr[i].CT[0] = 0
		self.PersonMgr[i].CT[1] = 0
		self.PersonMgr[i].RobDeal = -1
		bets := self.GetParam(TYPE_PTJ_BET)
		if bets == 0 { //! 每次选分
			self.PersonMgr[i].Bets = 0
		} else {
			if self.GetParam(TYPE_PTJ_SD) == 0 {
				self.PersonMgr[i].Bets = bets*100 + bets
			} else {
				self.PersonMgr[i].Bets = bets*10000 + bets*100 + bets
			}
		}
		self.PersonMgr[i].Card = make([]int, 0)
		self.PersonMgr[i].Ready = false
	}

	if self.GetParam(TYPE_PTJ_DEAL) == 1 { //! 轮庄模式
		//! 确定庄家
		if DealerPos+1 >= len(self.PersonMgr) {
			DealerPos = -1
		}
		self.PersonMgr[DealerPos+1].Dealer = true
	} else if self.GetParam(TYPE_PTJ_DEAL) == 2 { //! 赢家庄
		self.PersonMgr[WinPos].Dealer = true
	}

	//! 发牌
	if len(self.PJ.Card) == 0 {
		self.PJ = NewCard_TJ()
		self.Card = make([]int, 0)
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gameptjbegin", self.getInfo(person.Uid))
	}

	if self.GetParam(TYPE_PTJ_DEAL) != 0 && self.GetParam(TYPE_PTJ_BET) != 0 { //! 不抢庄又不下注直接发牌
		self.GameCard()
	}

	self.room.flush()
}

//! 抢庄
func (self *Game_PTJ) GameDeal(uid int64, ok bool) {
	if !self.room.Begin { //! 未开始不能抢庄
		return
	}
	if self.GetParam(TYPE_PTJ_DEAL) != 0 { //! 不是抢庄模式
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
		deal := make([]*Game_PTJ_Person, 0)
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

		if self.GetParam(TYPE_PTJ_BET) != 0 { //! 不用下注就发牌
			self.GameCard()
		}
	}

	self.room.flush()
}

//! 亮牌
func (self *Game_PTJ) GameView(uid int64, view []int) {
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
	if self.GetPTJValue(person.CT[0]) > self.GetPTJValue(person.MaxType) {
		person.MaxType = person.CT[0]
	}
	if len(view) == 4 {
		person.CT[1] = GetPTJType(view[2], view[3])
		if self.GetPTJValue(person.CT[1]) > self.GetPTJValue(person.MaxType) {
			person.MaxType = person.CT[1]
		}
	}

	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].View {
			num++
		}

		var msg Msg_GamePTJ_Send_View
		msg.Uid = uid
		if self.PersonMgr[i].Uid == uid || self.GetParam(TYPE_PTJ_DXJ) == 1 {
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
func (self *Game_PTJ) GameReady(uid int64) {
	if self.room.IsBye() {
		return
	}

	if self.room.Begin { //! 已经开始了不允许准备
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经开始了，不能准备")
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

	if num == len(self.room.Uid) && num >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 准备的人数达到游戏最小人数
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
		self.OnBegin()
		return
	}

	var msg staticfunc.Msg_Uid
	msg.Uid = uid
	self.room.broadCastMsg("gameready", &msg)

	self.room.flush()
}

//! 下注
func (self *Game_PTJ) GameBets(uid int64, bets int) {
	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	if bets <= 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "下注无效")
		return
	}

	if self.GetParam(TYPE_PTJ_BET) > 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "该模式无法下注")
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
func (self *Game_PTJ) GameCard() {
	if self.PJ.random == nil {
		self.PJ.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	var sz [2]int
	sz[0] = self.PJ.random.Intn(6) + 1
	sz[1] = self.PJ.random.Intn(6) + 1

	if self.GetParam(TYPE_PTJ_DXJ) == 0 {
		for i := 0; i < len(self.PersonMgr); i++ {
			self.PersonMgr[i].Card = self.PJ.Deal(4)
		}
		self.Cache = append(self.Cache, self.PJ.Deal(4*(4-len(self.PersonMgr)))...)
	} else {
		for i := 0; i < len(self.PersonMgr); i++ {
			self.PersonMgr[i].Card = self.PJ.Deal(2)
		}
		self.Cache = append(self.Cache, self.PJ.Deal(2*(4-len(self.PersonMgr)))...)
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		var msg Msg_GamePTJ_Card
		msg.SZ = sz
		msg.Card = self.PersonMgr[i].Card
		msg.View = self.Card
		self.room.SendMsg(self.PersonMgr[i].Uid, "gameptjcard", &msg)
	}
}

//! 结算
func (self *Game_PTJ) OnEnd() {
	self.room.SetBegin(false)

	var dealer *Game_PTJ_Person
	lst := make([]*Game_PTJ_Person, 0)
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
		dealerwin := self.GetWin(dealer, lst[i])
		if dealerwin > 0 {
			dealer.CurScore += bet1
			lst[i].CurScore -= bet1
			win++
		} else if dealerwin < 0 {
			dealer.CurScore -= bet1
			lst[i].CurScore += bet1
			win--
		}

		//! 计算第二道
		if dealerwin > 0 {
			if self.GetPTJValue(dealer.CT[0]) >= 83 && (self.GetParam(TYPE_PTJ_DXJ) == 1 || self.GetPTJValue(dealer.CT[1]) >= 83) {
				dealer.CurScore += bet2
				lst[i].CurScore -= bet2
			}
		} else if dealerwin < 0 {
			if self.GetPTJValue(lst[i].CT[0]) >= 83 && (self.GetParam(TYPE_PTJ_DXJ) == 1 || self.GetPTJValue(lst[i].CT[1]) >= 83) {
				dealer.CurScore -= bet2
				lst[i].CurScore += bet2
			}
		}

		//! 计算第三道
		if dealerwin > 0 {
			if self.GetPTJValue(dealer.CT[0]) >= 94 && (self.GetParam(TYPE_PTJ_DXJ) == 1 || self.GetPTJValue(dealer.CT[1]) >= 94) {
				dealer.CurScore += bet3
				lst[i].CurScore -= bet3
			}
		} else if dealerwin < 0 {
			if self.GetPTJValue(lst[i].CT[0]) >= 94 && (self.GetParam(TYPE_PTJ_DXJ) == 1 || self.GetPTJValue(lst[i].CT[1]) >= 94) {
				dealer.CurScore -= bet3
				lst[i].CurScore += bet3
			}
		}

		lst[i].Score += lst[i].CurScore
		if lst[i].CurScore > 0 {
			lst[i].Win++
			if lst[i].CurScore > lst[i].MaxScore {
				lst[i].MaxScore = lst[i].CurScore
			}
		}
	}
	dealer.Score += dealer.CurScore
	if dealer.CurScore > 0 {
		dealer.Win++
		if dealer.CurScore > dealer.MaxScore {
			dealer.MaxScore = dealer.CurScore
		}
	}
	if win == len(lst) {
		dealer.Kill++
	} else if -win == len(lst) {
		dealer.Dead++
	}

	//! 记录
	var record Rec_GamePTJ
	record.Time = time.Now().Unix()
	record.Roomid = self.room.Id*100 + self.room.Step
	record.MaxStep = self.room.MaxStep
	record.Param1 = self.room.Param1
	record.Param2 = self.room.Param2
	record.View = self.Card

	//! 发消息
	agentinfo := make([]staticfunc.JS_CreateRoomMem, 0)
	var msg Msg_GamePTJ_End
	for i := 0; i < len(self.PersonMgr); i++ {
		self.Card = append(self.Card, self.PersonMgr[i].Card...)

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
		agentinfo = append(agentinfo, staticfunc.JS_CreateRoomMem{son.Uid, "", "", son.Total})

		var rec Son_Rec_GamePTJ
		rec.Uid = self.PersonMgr[i].Uid
		rec.Name = self.room.GetName(self.PersonMgr[i].Uid)
		rec.Head = self.room.GetHead(self.PersonMgr[i].Uid)
		rec.Card = self.PersonMgr[i].Card
		rec.Bets = self.PersonMgr[i].Bets
		rec.Dealer = self.PersonMgr[i].Dealer
		rec.Score = self.PersonMgr[i].CurScore
		rec.Total = self.PersonMgr[i].Score
		rec.RobDeal = self.PersonMgr[i].RobDeal
		rec.CT = self.PersonMgr[i].CT
		record.Info = append(record.Info, rec)
	}
	self.Card = append(self.Card, self.Cache...)
	self.Cache = make([]int, 0)
	self.room.AddRecord(lib.HF_JtoA(&record))
	self.room.broadCastMsg("gameptjend", &msg)
	self.room.AgentResult(agentinfo)

	if self.room.IsBye() {
		self.OnBye()
		self.room.Bye()
		return
	}

	self.room.flush()
}

func (self *Game_PTJ) OnBye() {
	info := make([]staticfunc.JS_CreateRoomMem, 0)
	var msg Msg_GamePTJ_Bye
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GamePTJ_Bye
		son.Uid = self.PersonMgr[i].Uid
		son.Win = self.PersonMgr[i].Win
		son.Kill = self.PersonMgr[i].Kill
		son.Dead = self.PersonMgr[i].Dead
		son.MaxScore = self.PersonMgr[i].MaxScore
		son.MaxType = self.PersonMgr[i].MaxType
		son.Score = self.PersonMgr[i].Score
		msg.Info = append(msg.Info, son)
		info = append(info, staticfunc.JS_CreateRoomMem{son.Uid, "", "", son.Score})

		GetServer().SqlScoreLog(self.PersonMgr[i].Uid, self.room.GetName(self.PersonMgr[i].Uid), self.room.GetHead(self.PersonMgr[i].Uid), self.room.Type, self.room.Id, self.PersonMgr[i].Score)
	}
	self.room.broadCastMsg("gameptjbye", &msg)

	self.room.ClubResult(info)
}

func (self *Game_PTJ) OnExit(uid int64) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
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

	if num == len(self.room.Uid) && num >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 准备的人数达到游戏最小人数
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
		self.OnBegin()
		return
	}
}

func (self *Game_PTJ) getInfo(uid int64) *Msg_GamePTJ_Info {
	var msg Msg_GamePTJ_Info
	msg.Begin = self.room.Begin
	msg.Card = make([]int, len(self.PJ.Card))
	msg.Card = append(msg.Card, self.Card...)
	msg.Info = make([]Son_GamePTJ_Info, 0)
	for _, value := range self.PersonMgr {
		var son Son_GamePTJ_Info
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

func (self *Game_PTJ) GetPerson(uid int64) *Game_PTJ_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

func (self *Game_PTJ) GetWin(deal *Game_PTJ_Person, xian *Game_PTJ_Person) int {
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

//! 是否上道
//func (self *Game_PTJ) IsSD(person *Game_PTJ_Person) bool {
//	if person.CT[0] <= person.CT[1] {
//		return person.CT[0] >= 82
//	} else {
//		return person.CT[1] >= 82
//	}
//}

func (self *Game_PTJ) OnTime() {

}

func (self *Game_PTJ) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_PTJ) OnIsBets(uid int64) bool {
	return false
}

//! 根据牌型得到牌值
func (self *Game_PTJ) GetPTJValue(ct int) int {
	if ct == 0 {
		return 0
	}

	csv, ok := staticfunc.GetCsvMgr().Data["ptj"][ct]
	if !ok {
		return 0
	}

	if lib.HF_Atoi(csv["type"]) == 1 && self.GetParam(TYPE_PTJ_ZD) == 1 {
		return lib.HF_Atoi(csv["value2"])
	}

	if lib.HF_Atoi(csv["type"]) == 2 && self.GetParam(TYPE_PTJ_DJNN) == 1 {
		return lib.HF_Atoi(csv["value2"])
	}

	if lib.HF_Atoi(csv["type"]) == 3 && self.GetParam(TYPE_PTJ_GZ) == 1 {
		return lib.HF_Atoi(csv["value2"])
	}

	if lib.HF_Atoi(csv["type"]) == 4 && self.GetParam(TYPE_PTJ_TWJ) == 1 {
		return lib.HF_Atoi(csv["value2"])
	}

	return lib.HF_Atoi(csv["value1"])
}

//! 结算所有人
func (self *Game_PTJ) OnBalance() {
}
