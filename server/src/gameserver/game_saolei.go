package gameserver

import (
	"lib"
	"math/rand"
	"staticfunc"
	"time"
)

var SaoLei_GetBoomTime int64 = 30

//! 记录结构
type Rec_GameSaoLei struct {
	Info     []Son_Rec_GameSaoLei `json:"info"`
	Roomid   int                  `json:"roomid"`
	Time     int64                `json:"time"`
	RoomStep int                  `json:"roomstep"`
}
type Son_Rec_GameSaoLei struct {
	Uid     int64   `json:"uid"`
	Name    string  `json:"name"`
	Head    string  `json:"head"`
	Dealer  bool    `json:"dealer"`
	Score   float64 `json:"score"`
	Total   float64 `json:"total"`
	BoomPos int     `json:"boompos"` //! 猜的埋雷地方
}

//!
type Msg_GameSaoLei_Info struct {
	Begin     bool                  `json:"begin"` //! 是否开始
	Gamestate int                   `json:"gamestate"`
	BoomNum   int                   `json:"boomnum"`
	Ready     []int64               `json:"ready"`   //! 准备的人
	BoomPos   int                   `json:"boompos"` //! 地雷的雷针数
	BoomUid   int64                 `json:"boomuid"` //! 谁埋得雷
	Bets      float64               `json:"bets"`    //! 中雷的时候的倍率
	HadTime   int64                 `json:"hadtime"` //!剩余时间
	Info      []Son_GameSaoLei_Info `json:"info"`
}
type Son_GameSaoLei_Info struct {
	Uid        int64   `json:"uid"`        //! uid
	Dealer     bool    `json:"dealer"`     //! 是否庄家
	Score      float64 `json:"score"`      //! 当局分数
	BoomPos    int     `json:"boompos"`    //! 猜的埋雷地方
	Total      float64 `json:"total"`      //! 总分
	SetBoomNum int     `json:"setboomnum"` //! 埋雷的次数
	GetBoomNum int     `json:"getboomnum"` //! 挖到地雷的次数
	BeBoomNum  int     `json:"beboomnum"`  //! 埋雷被挖到的次数
	BestNum    int     `json:"bestnum"`    //! 运气王次数
	BestScore  bool    `json:"bestscore"`  //! 是否是运气王
	Bets       float64 `json:"bets"`       //! 中雷的时候的倍率
	BoomScore  float64 `json:"boomscore"`  //! 中雷需要给与别人的分数
	BoomId     []int   `json:"boomid"`     //! 埋了哪几个雷
}

type Son_GameSaoLei_Bye struct {
	Uid        int64   `json:"uid"`        //! uid
	Dealer     bool    `json:"dealer"`     //! 是否庄家
	Score      float64 `json:"score"`      //! 当局分数
	BoomPos    int     `json:"boompos"`    //! 猜的埋雷地方
	Total      float64 `json:"total"`      //! 总分
	SetBoomNum int     `json:"setboomnum"` //! 埋雷的次数
	GetBoomNum int     `json:"getboomnum"` //! 挖到地雷的次数
	BeBoomNum  int     `json:"beboomnum"`  //! 埋雷被挖到的次数
	BestNum    int     `json:"bestnum"`    //! 运气王次数
	BestScore  bool    `json:"bestscore"`  //! 是否是运气王
	Bets       float64 `json:"bets"`       //! 中雷的时候的倍率
	BoomScore  float64 `json:"boomscore"`  //! 中雷需要给与别人的分数
	BoomId     []int   `json:"boomid"`     //! 埋了哪几个雷
}

//! 结算
type Msg_GameSaoLei_End struct {
	Info []Son_GameSaoLei_Info `json:"info"`
}

//! 房间结束
type Msg_GameSaoLei_Bye struct {
	Info []Son_GameSaoLei_Bye `json:"info"`
	Time int64                `json:"time"`
}

//! 放雷 和 扫雷
type Msg_GameSaoLei_Boom struct {
	Uid       int64   `json:"uid"`
	Pos       int     `json:"pos"`       //! 雷位置
	Score     float64 `json:"score"`     //! 当局分数
	BoomScore float64 `json:"boomscore"` //! 中雷需要给与别人的分数
	Bets      float64 `json:"bets"`      //! 中雷的时候的倍率
	Total     float64 `json:"total"`     //! 玩家总分
}

type Msg_GameSaoLeiDealer struct {
	BoomNum int     `json:"boomnum"` //!剩余雷数
	BoomId  int     `json:"boomid"`  // !雷的编号
	Score   float64 `json:"score"`
	Total   float64 `json:"total"` //! 玩家总分
}

type Msg_GameStateChange struct {
	GameState int     `json:"gamestate"` //游戏状态
	BoomPos   int     `json:"boompos"`   //雷针数字
	BoomUid   int64   `json:"boomuid"`   //! 谁埋得雷
	Bets      float64 `json:"bets"`      //! 中雷的时候的倍率
	BoomNum   int     `json:"boomnum"`   //!剩余雷数
}

//!玩家挖雷开始时间
type Msg_GameSaoLeiTimeToGetBoom struct {
	Time int64 `json:"time"`
}

///////////////////////////////////////////////////////
type Game_SaoLei_Person struct {
	Uid        int64   `json:"uid"`        //! uid
	Total      float64 `json:"total"`      //! 积分
	Dealer     bool    `json:"dealer"`     //! 是否庄家
	CurScore   float64 `json:"curscore"`   //! 当前局的分数
	BoomPos    int     `json:"boompos"`    //! 猜的埋雷地方
	SetBoomNum int     `json:"setboomnum"` //! 埋雷的次数
	GetBoomNum int     `json:"getboomnum"` //! 挖到地雷的次数
	BeBoomNum  int     `json:"beboomnum"`  //! 埋雷被挖到的次数
	BestScore  bool    `json:"bestscore"`  //! 是否是运气王
	BestNum    int     `json:"bestnum"`    //! 运气王次数
	Bets       float64 `json:"bets"`       //! 中雷的时候的倍率
	BoomScore  float64 `json:"boomscore"`  //! 中雷需要给与别人的分数
	BoomId     []int   `json:"boomid"`     //! 埋了哪几个雷
}

func (self *Game_SaoLei_Person) Init() {
	self.Dealer = false
	self.CurScore = 0
	self.BoomPos = -1
	self.Bets = 0
	self.BoomScore = 0
}

type Game_SaoLei struct {
	Ready          []int64               `json:"ready"` //! 已经准备的人
	PersonMgr      []*Game_SaoLei_Person `json:"personmgr"`
	CurStep        int64                 `json:"curstep"`        //! 当前操作人
	BOOM           [10]int               `json:"boom"`           //! 雷区    一个雷
	BoomPos        int                   `json:"boompos"`        //! 地雷的雷针数
	BoomUid        int64                 `json:"boomuid"`        //! 谁埋得雷
	Dealer_i       int                   `json:"dealer_i"`       //! 庄家在人物数组的位置
	GetBoomTime    bool                  `json:"getboomtime"`    //! 是否挖雷时间
	GetBoomEndTime int64                 `json:"getboomendtime"` //! 挖雷结束时间
	DealerAll      []int64               `json:"dealerall"`      //! 抢庄模式所有的庄
	BoomPosAll     []int                 `json:"boompos"`        //! 每局的雷针数
	ScoreAll       []float64             `json:"scoreall"`       //! 每个玩家分值
	geti           int                   `json:"geti"`           //! 第几个玩家挖
	Gamestate      int                   `json:"gamestate"`      //! 玩家处于的状态  0 准备  1 埋雷 2 扫雷
	BoomNum        int                   `json:"boomnum"`        //! 雷的数目
	Bets           float64               `json:"bets"`           //! 中雷的时候的倍率

	room *Room
}

func NewGame_SaoLei() *Game_SaoLei {
	game := new(Game_SaoLei)
	game.Ready = make([]int64, 0)
	game.PersonMgr = make([]*Game_SaoLei_Person, 0)
	game.GetBoomTime = false
	game.Gamestate = 0

	return game
}

func (self *Game_SaoLei) OnInit(room *Room) {
	self.room = room
	if self.room.Param1%10 == 0 {
		self.BoomNum = 5
	} else {
		self.BoomNum = 10
	}
}

func (self *Game_SaoLei) OnRobot(robot *lib.Robot) {

}

func (self *Game_SaoLei) OnSendInfo(person *Person) {
	person.SendMsg("gameSaoLeiinfo", self.getInfo(person.Uid))
}

func (self *Game_SaoLei) GetPerson(uid int64) *Game_SaoLei_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

func (self *Game_SaoLei) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "gameready": //! 游戏准备
		self.GameReady(msg.Uid)
		//	case "gamesldealer": //抢庄
		//		self.GameDealer(msg.Uid, msg.V.(*Msg_GameLzDealer).Type)
	case "gamesetboom": //埋雷
		self.GameSetBoom(msg.Uid, msg.V.(*Msg_GameBoom).Pos)
	case "gamegetboom": //挖雷
		self.GameGetBoom(msg.Uid, msg.V.(*Msg_GameBoom).Pos)
	}
}

func (self *Game_SaoLei) GameSetBoom(uid int64, Pos int) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏还没有开始，不能埋雷")
		return
	}

	if self.room.Param2 <= 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "param2参数错误")
		return
	}

	if Pos < 0 || Pos >= 10 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "埋雷地方错误")
		return
	}

	if self.Gamestate != 1 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不是埋雷的状态")
		return
	}

	if self.room.Param1/10%10 == 0 { //房主模式
		if self.PersonMgr[self.Dealer_i].Uid != uid {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "不是庄家，不能埋雷")
			return
		} else {
			self.BoomPos = Pos
			//self.PersonMgr[self.Dealer_i].CurScore -= float64(self.room.Param2)
			self.PersonMgr[self.Dealer_i].SetBoomNum++
			self.PersonMgr[self.Dealer_i].Total -= float64(self.room.Param2)
			self.PersonMgr[self.Dealer_i].Total = self.ToFloatNormal(self.PersonMgr[self.Dealer_i].Total)
		}
		//		self.BoomNum--

		self.Gamestate = 2
		var msg Msg_GameSaoLei_Boom
		msg.Pos = Pos
		msg.Uid = uid
		msg.Score = self.PersonMgr[self.Dealer_i].CurScore
		msg.Total = self.PersonMgr[self.Dealer_i].Total
		self.room.broadCastMsg("gamesetboom", &msg)

		self.GetBoomTime = true //可以开始挖雷
		//		for i := 0; i < len(self.PersonMgr); i++ {
		//			person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)

		//			var msg Msg_GameSaoLeiTimeToGetBoom
		//			msg.Time = SaoLei_GetBoomTime
		//			person.SendMsg("gamegetboomtime", &msg)
		//		}

		self.ScoreAll = self.GetScore()

		var msg1 Msg_GameStateChange
		msg1.GameState = self.Gamestate
		msg1.BoomPos = self.BoomPos
		msg1.BoomUid = self.BoomUid
		self.room.broadCastMsg("gamestatechange", &msg1)

		self.GetBoomEndTime = time.Now().Unix() + SaoLei_GetBoomTime

	} else { //抢庄模式
		var booo int
		if self.room.Param1%10 == 0 {
			booo = 5
		} else {
			booo = 10
		}
		if booo == len(self.DealerAll) {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "抢庄已经完成不能再抢了")
			return
		}

		var i int
		for i = 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Uid == uid {
				self.PersonMgr[i].SetBoomNum++
				self.DealerAll = append(self.DealerAll, uid)
				self.BoomPosAll = append(self.BoomPosAll, Pos)

				self.PersonMgr[i].BoomId = append(self.PersonMgr[i].BoomId, booo-self.BoomNum) //抢到雷的编号

				self.PersonMgr[i].Total -= float64(self.room.Param2)
				self.PersonMgr[i].Total = self.ToFloatNormal(self.PersonMgr[i].Total)
				break
			}
		}
		if i >= len(self.PersonMgr) {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "没有该玩家")
			return
		}

		self.BoomNum--
		person := self.GetPerson(uid)

		var msg Msg_GameSaoLeiDealer
		msg.BoomNum = self.BoomNum
		msg.BoomId = booo - self.BoomNum
		msg.Score = person.CurScore
		msg.Total = person.Total
		self.room.broadCastMsg("gamesldealer", &msg)

		if booo == len(self.DealerAll) {
			if self.room.Param1%10 == 0 { //! BoomNum此时相当于局数
				self.BoomNum = 5
			} else {
				self.BoomNum = 10
			}

			self.OnBegin()
		}
		self.room.flush()
	}
}

func (self *Game_SaoLei) GameGetBoom(uid int64, Pos int) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏还没有开始，不能埋雷")
		return
	}

	if Pos == -1 {
		Pos = self.GetNoBoom()
	}

	//	if Pos < 0 || Pos >= 10 {
	//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "挖雷地方错误")
	//		return
	//	}
	//	if self.BOOM[Pos] != 0 {
	//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "此地有人挖了")
	//		return
	//	}

	var BoomScore float64
	var i int
	for i = 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].BoomPos != -1 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "该玩家已经扫雷了")
				return
			}
			self.PersonMgr[i].BoomPos = Pos
			self.PersonMgr[i].CurScore = self.ScoreAll[self.geti]
			self.PersonMgr[i].CurScore = self.ToFloatNormal(self.PersonMgr[i].CurScore)

			lib.GetLogMgr().Output(lib.LOG_DEBUG, "庄：", self.Dealer_i, "该玩家：", i)
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "雷针：", self.BoomPos)
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "玩家尾数：", int(self.ScoreAll[self.geti]*float64(100))%10)

			s_float := self.PersonMgr[i].CurScore * float64(100)
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "浮点型*100 = ", s_float)
			s_float = self.ToFloatNormal(s_float)
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "戴上安全套的浮点型：", s_float)

			s_int := int(s_float)
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "整型*100 = ", s_int)

			w_num := s_int % 10
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "w_num = ", w_num)

			//中雷
			if w_num == self.BoomPos && self.Dealer_i != i {
				BoomScore = float64(self.room.Param2) * (float64(2) - 0.2*float64(len(self.PersonMgr)-5))
				self.PersonMgr[i].Bets = (float64(2) - 0.2*float64(len(self.PersonMgr)-5))

				//				self.PersonMgr[i].CurScore -= BoomScore
				//				self.PersonMgr[i].CurScore = self.ToFloatNormal(self.PersonMgr[i].CurScore)

				self.PersonMgr[self.Dealer_i].Total += BoomScore
				self.PersonMgr[self.Dealer_i].Total = self.ToFloatNormal(self.PersonMgr[self.Dealer_i].Total)

				self.PersonMgr[i].GetBoomNum++
				self.PersonMgr[self.Dealer_i].BeBoomNum++
				self.PersonMgr[i].BoomScore -= BoomScore
				self.PersonMgr[self.Dealer_i].BoomScore += BoomScore

				self.PersonMgr[i].Total -= BoomScore
				self.PersonMgr[i].Total = self.ToFloatNormal(self.PersonMgr[i].Total)
			}
			self.PersonMgr[i].Total += self.PersonMgr[i].CurScore
			self.PersonMgr[i].Total = self.ToFloatNormal(self.PersonMgr[i].Total)
			break
		}
	}
	if i >= len(self.PersonMgr) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "没找到玩家")
		return
	}

	self.BOOM[Pos] = 1
	self.geti++

	//person := self.GetPerson(uid)
	//	var msg Msg_GameSaoLei_Boom
	//	msg.Pos = Pos
	//	msg.Uid = uid
	//	msg.Score = person.CurScore
	//	msg.BoomScore = BoomScore
	//	msg.Total = person.Total
	//	msg.Bets = float64(2) - 0.2*float64(len(self.PersonMgr)-5)
	self.room.broadCastMsg("gamegetboom", self.GetOneInfo(uid))

	if self.geti >= len(self.PersonMgr) {
		self.BoomNum--
		self.OnEnd()
	}
}

func (self *Game_SaoLei) ToFloatNormal(value float64) float64 {
	value *= float64(1000)
	value_int := int(value)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "保护前：", value_int)
	zf := 0
	if value_int < 0 {
		value_int *= -1
		zf = 1
	}
	if value_int%10 == 9 {
		value_int = int(value_int)/10 + 1
	} else {
		value_int = int(value_int) / 10
	}
	if zf == 1 {
		value_int *= -1
	}
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "保护后：", value_int)
	a := float64(value_int) / float64(100)
	return a
}

func (self *Game_SaoLei) GetOneInfo(uid int64) *Son_GameSaoLei_Info {
	value := self.GetPerson(uid)
	if value == nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不存在该玩家")
		return nil
	}

	var son Son_GameSaoLei_Info
	son.Uid = value.Uid
	son.Dealer = value.Dealer
	son.Score = value.CurScore
	son.Total = value.Total
	son.BoomPos = value.BoomPos
	son.Bets = value.Bets
	son.BoomScore = value.BoomScore
	son.SetBoomNum = value.SetBoomNum
	son.GetBoomNum = value.GetBoomNum
	son.BeBoomNum = value.BeBoomNum
	son.BestNum = value.BestNum
	son.BoomId = value.BoomId

	return &son
}

func (self *Game_SaoLei) FindQDealer() int {
	Uid := self.DealerAll[self.room.Step-1]
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == Uid {
			return i
		}
	}

	return -1
}

func (self *Game_SaoLei) FindBoomPos() int {
	return self.BoomPosAll[self.room.Step-1]
}

func (self *Game_SaoLei) OnBegin() {
	if self.room.Param2 <= 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "param2参数错误")
		return
	}
	if self.room.IsBye() {
		return
	}

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "param1: ", self.room.Param1, "step", self.room.Step)
	if self.room.Param1/10 == 0 || (self.Gamestate != 1) {
		self.room.SetBegin(true)

		for i := 0; i < len(self.room.Uid); i++ { //! 重新初始化人
			if i >= len(self.PersonMgr) {
				person := new(Game_SaoLei_Person)
				person.Uid = self.room.Uid[i]
				self.PersonMgr = append(self.PersonMgr, person)
			}
		}

		for i := 0; i < len(self.PersonMgr); i++ {
			self.PersonMgr[i].Init()
		}
		self.Bets = (float64(2) - 0.2*float64(len(self.PersonMgr)-5))
	}

	if self.room.Param1/10 == 0 {
		self.Gamestate = 1
	} else {
		self.Gamestate = 2
		//	self.BoomNum--

		self.GetBoomTime = true //可以开始挖雷

		//		for i := 0; i < len(self.PersonMgr); i++ {
		//			person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)

		//			var msg Msg_GameSaoLeiTimeToGetBoom
		//			msg.Time = SaoLei_GetBoomTime
		//			person.SendMsg("gamegetboomtime", &msg)
		//		}
		self.GetBoomEndTime = time.Now().Unix() + SaoLei_GetBoomTime
	}

	//初始化雷
	for i := 0; i < 10; i++ {
		self.BOOM[i] = 0
	}
	self.BoomPos = -1

	if self.room.Param1/10 == 0 { //! 房主庄
		self.PersonMgr[0].Dealer = true
		self.CurStep = self.PersonMgr[0].Uid
		self.Dealer_i = 0
		self.BoomUid = self.PersonMgr[0].Uid
	} else { //! 抢庄
		DealerPos := self.FindQDealer()
		if DealerPos == -1 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "当前庄家不存在")
			return
		}
		self.PersonMgr[DealerPos].Dealer = true
		self.Dealer_i = DealerPos
		self.BoomPos = self.FindBoomPos()
		self.BoomUid = self.PersonMgr[self.Dealer_i].Uid
		self.ScoreAll = self.GetScore()
	}

	self.geti = 0

	var msg1 Msg_GameStateChange
	msg1.GameState = self.Gamestate
	msg1.BoomPos = self.BoomPos
	msg1.BoomUid = self.BoomUid
	if self.Gamestate == 1 {
		msg1.Bets = (float64(2) - 0.2*float64(len(self.PersonMgr)-5))
		msg1.BoomNum = self.BoomNum
	}
	self.room.broadCastMsg("gamestatechange", &msg1)

	//	for i := 0; i < len(self.PersonMgr); i++ {
	//		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
	//		if person == nil {
	//			continue
	//		}
	//		person.SendMsg("gameSaoLeibegin", self.getInfo(person.Uid))
	//	}

	self.room.flush()
}

//! 准备
func (self *Game_SaoLei) GameReady(uid int64) {
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
		//判断房主是否点击了开始
		for i := 0; i < len(self.Ready); i++ {
			if self.Ready[i] == self.room.Uid[0] {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
				if self.room.Param1/10 == 0 { //房主庄
					self.OnBegin()
					return
				} else { //抢庄
					if self.room.Step == 0 { //第一局开始抢庄
						self.room.SetBegin(true)

						for i := 0; i < len(self.room.Uid); i++ { //! 重新初始化人
							if i >= len(self.PersonMgr) {
								person := new(Game_SaoLei_Person)
								person.Uid = self.room.Uid[i]
								self.PersonMgr = append(self.PersonMgr, person)
							}
						}
						for i := 0; i < len(self.PersonMgr); i++ {
							self.PersonMgr[i].Init()
						}
						self.Bets = (float64(2) - 0.2*float64(len(self.PersonMgr)-5))
						self.Gamestate = 1
						var msg1 Msg_GameStateChange
						msg1.GameState = self.Gamestate
						msg1.BoomNum = self.BoomNum
						if self.Gamestate == 1 {
							msg1.Bets = (float64(2) - 0.2*float64(len(self.PersonMgr)-5))
						}
						self.room.broadCastMsg("gamestatechange", &msg1)

						var msg Msg_GameSaoLeiDealer
						msg.BoomNum = self.BoomNum
						self.room.broadCastMsg("gamestartdealer", &msg)
					} else {
						self.GetBoomTime = true //第二局可以直接挖雷
						self.OnBegin()
						return
					}
				}
			}
		}
	}

	var msg staticfunc.Msg_Uid
	msg.Uid = uid
	self.room.broadCastMsg("gameready", &msg)

	self.room.flush()
}

func (self *Game_SaoLei) GetScore() []float64 {
	var Score []int
	Score = make([]int, len(self.PersonMgr))
	allgetscore := (self.room.Param2 * 100)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "当前中雷概率：30%")
	DesScore := make([]int, len(self.PersonMgr))
	for i := 0; i < len(self.PersonMgr); i++ {
		if random.Intn(100) < 30 {
			DesScore[i] = self.BoomPos
		} else {
			DesScore[i] = random.Intn(10)
			if DesScore[i] == self.BoomPos {
				DesScore[i] = self.BoomPos - 1
			}
		}
	}
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "小的分数 : ", DesScore)

	A1 := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		A1 += DesScore[i]
	}

	allgetscore -= A1

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "分数 : ", allgetscore)
	Score[0] = random.Intn(allgetscore/2-10-50)/10*10 + 50
	var HadScore int    //当前玩家的分数
	HadScore = Score[0] //之前所有玩家的总和
	var BefScore int    //当前玩家的分数
	for i := 1; i < len(self.PersonMgr)-1; i++ {
		BefScore = allgetscore - HadScore - (len(self.PersonMgr)-i-1)*50
		Score[i] = random.Intn(BefScore-50)/10*10 + 50
		HadScore += Score[i]
	}
	Score[len(self.PersonMgr)-1] = allgetscore - HadScore

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "大的分数 : ", Score)

	for i := 0; i < len(self.PersonMgr); i++ {
		Score[i] = Score[i] + DesScore[i]
	}
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "总的分数 : ", Score)

	var fscore []float64
	fscore = make([]float64, len(self.PersonMgr))
	for i := 0; i < len(self.PersonMgr); i++ {
		fscore[i] = float64(Score[i]) / float64(100)
	}

	//打乱数组
	var Score1 []float64

	for i := 0; i < len(self.PersonMgr); i++ {
		r_i := random.Intn(len(fscore))
		Score1 = append(Score1, fscore[r_i])
		copy(fscore[r_i:], fscore[r_i+1:])
		fscore = fscore[:len(fscore)-1]
	}

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "所有分数: ", Score1)

	return Score1
}

//! 结算
func (self *Game_SaoLei) OnEnd() {
	self.room.SetBegin(false)

	self.Gamestate = 0

	//运气王
	best := 0
	for i := 1; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].CurScore > self.PersonMgr[best].CurScore {
			best = i
		}
	}
	self.PersonMgr[best].BestScore = true
	self.PersonMgr[best].BestNum++

	//! 记录
	var record Rec_GameSaoLei
	record.Time = time.Now().Unix()
	record.Roomid = self.room.Id*100 + self.room.Step
	record.RoomStep = self.room.Step

	var msg Msg_GameSaoLei_End
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameSaoLei_Info
		son.Dealer = self.PersonMgr[i].Dealer
		son.Score = self.PersonMgr[i].CurScore
		son.Total = self.PersonMgr[i].Total
		son.Uid = self.PersonMgr[i].Uid
		son.BeBoomNum = self.PersonMgr[i].BeBoomNum
		son.GetBoomNum = self.PersonMgr[i].GetBoomNum
		son.SetBoomNum = self.PersonMgr[i].SetBoomNum
		son.BoomPos = self.PersonMgr[i].BoomPos
		son.BestNum = self.PersonMgr[i].BestNum
		son.BestScore = self.PersonMgr[i].BestScore
		son.Bets = self.PersonMgr[i].Bets
		son.BoomScore = self.PersonMgr[i].BoomScore
		son.BoomId = self.PersonMgr[i].BoomId
		msg.Info = append(msg.Info, son)

		var _son Son_Rec_GameSaoLei
		_son.Uid = self.PersonMgr[i].Uid
		_son.Name = self.room.GetName(self.PersonMgr[i].Uid)
		_son.Head = self.room.GetHead(self.PersonMgr[i].Uid)
		_son.Dealer = self.PersonMgr[i].Dealer
		_son.Score = self.PersonMgr[i].CurScore
		_son.BoomPos = self.PersonMgr[i].BoomPos
		_son.Total = self.PersonMgr[i].Total
		record.Info = append(record.Info, _son)
	}
	self.room.AddRecord(lib.HF_JtoA(&record))
	self.room.broadCastMsg("gameSaoLeiend", &msg)

	if self.room.IsBye() {
		self.OnBye()
		self.room.Bye()
		return
	}

	var msg1 Msg_GameStateChange
	msg1.GameState = self.Gamestate
	self.room.broadCastMsg("gamestatechange", &msg1)

	self.Ready = make([]int64, 0)

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].BestScore = false
	}
	self.GetBoomTime = false

	self.room.flush()
}

func (self *Game_SaoLei) OnBye() {
	var score float64

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].CurScore > 0 {
			score += self.PersonMgr[i].CurScore
			score = self.ToFloatNormal(score)
		}
	}
	if self.room.Param2 != int(score) {
		self.PersonMgr[self.Dealer_i].Total += (float64(self.room.Param2) - score)
		self.PersonMgr[self.Dealer_i].Total = self.ToFloatNormal(self.PersonMgr[self.Dealer_i].Total)
	}

	if self.room.Param1/10%10 == 1 {
		step := self.room.Step
		for i := step; i < len(self.DealerAll); i++ {
			for j := 0; j < len(self.PersonMgr); j++ {
				if self.PersonMgr[j].Uid == self.DealerAll[i] {
					self.PersonMgr[j].Total += float64(self.room.Param2)
					self.PersonMgr[j].Total = self.ToFloatNormal(self.PersonMgr[j].Total)
					break
				}
			}
		}
	}

	var index []int //手气王总分下标
	best := 0
	for i := 1; i < len(self.PersonMgr); i++ {
		a := self.PersonMgr[best].Total
		a = a * float64(100)
		a = self.ToFloatNormal(a)
		b := self.PersonMgr[i].Total
		b = b * float64(100)
		b = self.ToFloatNormal(b)
		if int(a) < int(b) {
			best = i
		} else if int(a) == int(b) {
			index = append(index, i)
		}
		self.PersonMgr[i].BestScore = false
	}
	index = append(index, best)
	for i := 0; i < len(index); i++ {
		self.PersonMgr[index[i]].BestScore = true
	}

	info := make([]staticfunc.JS_CreateRoomMem, 0)
	var msg Msg_GameSaoLei_Bye
	msg.Time = time.Now().Unix()
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameSaoLei_Bye
		son.Dealer = self.PersonMgr[i].Dealer
		son.Score = self.PersonMgr[i].CurScore
		son.Total = self.PersonMgr[i].Total
		son.Uid = self.PersonMgr[i].Uid
		son.BeBoomNum = self.PersonMgr[i].BeBoomNum
		son.GetBoomNum = self.PersonMgr[i].GetBoomNum
		son.SetBoomNum = self.PersonMgr[i].SetBoomNum
		son.BoomPos = self.PersonMgr[i].BoomPos
		son.BestNum = self.PersonMgr[i].BestNum
		son.BestScore = self.PersonMgr[i].BestScore
		son.Bets = self.PersonMgr[i].Bets
		son.BoomScore = self.PersonMgr[i].BoomScore
		son.BoomId = self.PersonMgr[i].BoomId
		msg.Info = append(msg.Info, son)
		info = append(info, staticfunc.JS_CreateRoomMem{son.Uid, "", "", int(son.Total)})
	}
	self.room.broadCastMsg("gameSaoLeibye", &msg)

	self.room.ClubResult(info)
}

func (self *Game_SaoLei) OnExit(uid int64) {
	//find := false
	for i := 0; i < len(self.Ready); i++ {
		if self.Ready[i] == uid {
			copy(self.Ready[i:], self.Ready[i+1:])
			self.Ready = self.Ready[:len(self.Ready)-1]
			//find = true
			break
		}
	}

	//	if !find {
	//		if len(self.Ready) == len(self.room.Uid) && len(self.Ready) >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 准备的人数达到游戏最小人数
	//			//判断房主是否点击了开始
	//			for i := 0; i < len(self.Ready); i++ {
	//				if self.Ready[i] == self.room.Uid[0] {
	//					if self.room.Param1/10 == 0 { //房主庄
	//						self.OnBegin()
	//						return
	//					} else { //抢庄
	//					}
	//				}
	//			}
	//		}
	//	}
}

func (self *Game_SaoLei) getInfo(uid int64) *Msg_GameSaoLei_Info {
	var msg Msg_GameSaoLei_Info
	msg.Begin = self.room.Begin
	msg.BoomNum = self.BoomNum
	msg.Ready = make([]int64, 0)
	msg.Gamestate = self.Gamestate
	msg.BoomPos = self.BoomPos
	msg.BoomUid = self.BoomUid
	msg.HadTime = self.GetBoomEndTime - time.Now().Unix()
	msg.Bets = self.Bets

	msg.Info = make([]Son_GameSaoLei_Info, 0)

	if !msg.Begin { //! 没有开始,看哪些人已准备
		msg.Ready = self.Ready
	}
	for _, value := range self.PersonMgr {
		var son Son_GameSaoLei_Info
		son.Uid = value.Uid
		son.Dealer = value.Dealer
		son.Score = value.CurScore
		son.Total = value.Total
		son.BoomPos = value.BoomPos
		son.Bets = value.Bets
		son.BoomScore = son.BoomScore
		son.SetBoomNum = value.SetBoomNum
		son.GetBoomNum = value.GetBoomNum
		son.BeBoomNum = value.BeBoomNum
		son.BestNum = value.BestNum
		son.BoomId = value.BoomId
		msg.Info = append(msg.Info, son)
	}

	return &msg
}

func (self *Game_SaoLei) OnTime() {

	//挖雷时间
	if self.GetBoomTime {
		if time.Now().Unix() > self.GetBoomEndTime {
			self.GetBoomTime = false

			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].BoomPos == -1 { //没猜雷的玩家
					pos := self.GetNoBoom()
					self.GameGetBoom(self.PersonMgr[i].Uid, pos)
				}
			}
		}
	}

}

//返回一个没有被人挖的地方
func (self *Game_SaoLei) GetNoBoom() int {
	var NoBoomPos []int
	NoBoomPos = make([]int, 0)
	for i := 0; i < 10; i++ {
		if self.BOOM[i] == 0 {
			NoBoomPos = append(NoBoomPos, i)
		}
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	pos := random.Intn(len(NoBoomPos))

	return NoBoomPos[pos]
}

func (self *Game_SaoLei) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_SaoLei) OnIsBets(uid int64) bool {
	return false
}

//! 结算所有人
func (self *Game_SaoLei) OnBalance() {
}
