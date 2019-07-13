package gameserver

import (
	"lib"
	"math"
	"staticfunc"
	"time"
)

//! param1
const TYPE_MPQZ_DF = 0   //! 0,1/2  1,2/4  2,3/6   3,4/8   4,5/10
const TYPE_MPQZ_FB = 1   //! 翻倍规则 0  1
const TYPE_MPQZ_ZDQZ = 2 //! 最大抢庄 1-4
const TYPE_MPQZ_ZDKS = 3 //! 自动开始 0手动开始 1满4人开 2满5人开  3满6人开
const TYPE_MPQZ_XJTZ = 4 //! 闲家推注 0无  1,5倍   2,10倍   3，15倍
const TYPE_MPQZ_JZCP = 5 //! 禁止搓牌 0不禁止  1禁止
const TYPE_MPQZ_XZXZ = 6 //! 下注限制 0不限制  1限制
const TYPE_MPQZ_MAX = 7

//! param2
const TYPE_MPQZ_SZN = 0   //! 顺子牛 0没有  1有  下同
const TYPE_MPQZ_WHN = 1   //! 五花牛
const TYPE_MPQZ_THN = 2   //! 同花牛
const TYPE_MPQZ_WXN = 3   //! 五小牛
const TYPE_MPQZ_HLN = 4   //! 葫芦牛
const TYPE_MPQZ_ZDN = 5   //! 炸弹牛
const TYPE_MPQZ_SJN = 6   //! 顺金牛
const TYPE_MPQZ_ENTER = 7 //! 开始是否能进入 0可以  1不能
const TYPE_MPQZ_MAX2 = 8

//! 记录结构
type Rec_GameNiuNiuXY struct {
	Info    []Son_Rec_GameNiuNiuXY `json:"info"`
	Roomid  int                    `json:"roomid"`
	MaxStep int                    `json:"maxstep"`
	Param1  int                    `json:"param1"`
	Param2  int                    `json:"param2"`
	Time    int64                  `json:"time"`
	Host    int64                  `json:"host"`
}
type Son_Rec_GameNiuNiuXY struct {
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Card    []int  `json:"card"`
	Bets    int    `json:"bets"`
	Dealer  bool   `json:"dealer"`
	Score   int    `json:"score"`
	Type    int    `json:"type"`
	RobDeal int    `json:"robdeal"`
}

//!
type Msg_GameMPQZ_Info struct {
	Begin bool                `json:"begin"` //! 是否开始
	Info  []Son_GameMPQZ_Info `json:"info"`
	State int                 `json:"state"`
	Time  int64               `json:"time"` //! 倒计时
}

type Son_GameMPQZ_Info struct {
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
	Trust   bool  `json:"trust"` //! 是否托管
	Open    bool  `json:"open"`
}

//! 结算
type Msg_GameMPQZ_End struct {
	Info []Son_GameMPQZ_Info `json:"info"`
}

//! 房间结束
type Msg_GameMPQZ_Bye struct {
	Host int64              `json:"host"`
	Info []Son_GameMPQZ_Bye `json:"info"`
}
type Son_GameMPQZ_Bye struct {
	Uid   int64 `json:"uid"`
	Score int   `json:"score"`
	QZ    int   `json:"qz"`
	ZZ    int   `json:"zz"`
	TZ    int   `json:"tz"`
}

/////////////////////////////////////////////////////
type Game_MPQZ_Person struct {
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
	Trust    bool  `json:"trust"`    //! 是否托管
	Open     bool  `json:"open"`     //! 是否看牌
	QZNum    int   `json:"qz"`
	ZZNum    int   `json:"zz"`
	TZNum    int   `json:"tz"`
}

func (self *Game_MPQZ_Person) Init() {
	self.CT = 0
	self.CS = 0
	self.CurScore = 0
	self.Dealer = false
	self.Bets = 0
	self.Ready = false
	self.RobDeal = -1
	self.View = false
	self.Card = make([]int, 0)
	self.Open = false
}

type Game_MPQZ struct {
	PersonMgr []*Game_MPQZ_Person `json:"personmgr"`
	Card      *CardMgr            `json:"card"`
	State     int                 `json:"state"` //! 0准备阶段  1等待抢庄   2等待下注   3等待亮牌
	Time      int64               `json:"time"`
	Dealer    int64               `json:"dealer"`
	CardZB    [][]int             `json:"cardzb"`

	room *Room
}

func NewGame_MPQZ() *Game_MPQZ {
	game := new(Game_MPQZ)
	game.PersonMgr = make([]*Game_MPQZ_Person, 0)

	return game
}

func (self *Game_MPQZ) GetParam(_type int) int {
	return self.room.Param1 % int(math.Pow(10.0, float64(TYPE_MPQZ_MAX-_type))) / int(math.Pow(10.0, float64(TYPE_MPQZ_MAX-_type-1)))
}

func (self *Game_MPQZ) GetParam2(_type int) int {
	return self.room.Param2 % int(math.Pow(10.0, float64(TYPE_MPQZ_MAX2-_type))) / int(math.Pow(10.0, float64(TYPE_MPQZ_MAX2-_type-1)))
}

func (self *Game_MPQZ) OnInit(room *Room) {
	self.room = room
}

func (self *Game_MPQZ) OnRobot(robot *lib.Robot) {

}

func (self *Game_MPQZ) OnSendInfo(person *Person) {
	//! 观众模式游戏,观众进来只发送游戏信息
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			person.SendMsg("gamempqzinfo", self.getInfo(person.Uid))
			return
		}
	}

	if self.room.Type == 68 && self.room.Host == person.Uid { //! 固定庄家
		if self.room.Seat(person.Uid) {
			_person := new(Game_MPQZ_Person)
			_person.Init()
			if self.GetParam(TYPE_MPQZ_ZDQZ) == 2 {
				_person.Score = 100
			} else if self.GetParam(TYPE_MPQZ_ZDQZ) == 3 {
				_person.Score = 150
			} else if self.GetParam(TYPE_MPQZ_ZDQZ) == 4 {
				_person.Score = 200
			}
			_person.Uid = person.Uid
			self.PersonMgr = append(self.PersonMgr, _person)
		}
	}

	person.SendMsg("gamempqzinfo", self.getInfo(0))
}

func (self *Game_MPQZ) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "gameseat": //! 游戏坐下
		self.GameSeat(msg.Uid)
	case "gameready": //! 游戏准备
		self.GameReady(msg.Uid)
	case "gamebegin": //! 开始游戏
		self.GameBegin(msg.Uid)
	case "gamebets": //! 下注
		self.GameBets(msg.Uid, msg.V.(*Msg_GameBets).Bets)
	case "gameview": //! 亮牌
		self.GameView(msg.Uid, true)
	case "gamedealer": //! 抢庄
		self.GameDeal(msg.Uid, msg.V.(*Msg_GameDealer).Score)
	case "gameopen":
		self.GameOpen(msg.Uid)
		//case "gametrust": //! 托管
		//	self.GameTrust(msg.Uid, msg.V.(*Msg_GameDeal).Ok)
	}
}

func (self *Game_MPQZ) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Init()
	}

	self.SetTime(11)
	//! 确定庄家
	if self.IsBeginDealer() {
		if self.Dealer == 0 {
			if self.room.Type == 68 { //! 固定庄家
				self.Dealer = self.room.Uid[0]
			} else { //! 牛牛上庄
				self.Dealer = self.room.Uid[lib.HF_GetRandom(len(self.room.Uid))]
			}
		}
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Uid == self.Dealer {
				self.PersonMgr[i].Dealer = true
				self.PersonMgr[i].RobDeal = 1
				break
			}
		}
		self.State = 2
	} else {
		self.State = 1
	}

	//! 发牌
	//	self.Card = NewCard_NiuNiu(false)

	//	for i := 0; i < len(self.PersonMgr); i++ {
	//		if self.IsBeginCard() {
	//			self.PersonMgr[i].Card = self.Card.Deal(4)
	//		} else {
	//			self.PersonMgr[i].Card = make([]int, 0)
	//		}
	//	}

	//! 发牌
	self.Card = NewCard_NiuNiu(false)
	self.CardZB = make([][]int, 0)
	maxIndex := -1 //!　最大牌的数组下标
	for i := 0; i < len(self.PersonMgr); i++ {
		_card := self.Card.Deal(5)
		if maxIndex == -1 {
			maxIndex = i
		} else {
			ct, cs := GetGoldNiuNiuScore(_card, true, true, true, true, true, true, true)
			_ct, _cs := GetGoldNiuNiuScore(self.CardZB[maxIndex], true, true, true, true, true, true, true)
			if _ct < ct {
				maxIndex = i
			} else if _ct == ct && cs > _cs {
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
		if self.IsBeginCard() {
			self.PersonMgr[i].Card = self.PersonMgr[i].CardZB[0:4]
		} else {
			self.PersonMgr[i].Card = make([]int, 0)
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

//! 坐下
func (self *Game_MPQZ) GameSeat(uid int64) {
	if self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏正在进行，无法坐下")
		return
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid { //! 已经坐下，无法坐下
			return
		}
	}

	if !self.room.Seat(uid) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "无法坐下")
		return
	}

	_person := new(Game_MPQZ_Person)
	_person.Init()
	_person.Uid = uid
	self.PersonMgr = append(self.PersonMgr, _person)
	self.room.broadCastMsg("gamempqzseat", self.getInfo(0))
}

//! 抢庄
func (self *Game_MPQZ) GameDeal(uid int64, score int) {
	if !self.room.Begin { //! 未开始不能抢庄
		return
	}

	if self.IsBeginDealer() {
		return
	}

	if self.room.Type == 65 || self.room.Type == 66 {
		if score > self.GetParam(TYPE_MPQZ_ZDQZ) {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "下注超过上限")
			return
		}
	}

	robnum := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].RobDeal >= 0 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能重复抢庄")
				return
			} else {
				self.PersonMgr[i].RobDeal = score
				if score > 0 {
					self.PersonMgr[i].QZNum++
				}
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
		deal := make([]*Game_MPQZ_Person, 0)
		for i := 0; i < len(self.PersonMgr); i++ {
			if len(deal) == 0 {
				deal = append(deal, self.PersonMgr[i])
			} else {
				if self.PersonMgr[i].RobDeal > deal[0].RobDeal {
					deal = make([]*Game_MPQZ_Person, 0)
					deal = append(deal, self.PersonMgr[i])
				} else if self.PersonMgr[i].RobDeal == deal[0].RobDeal {
					deal = append(deal, self.PersonMgr[i])
				}
			}
		}

		dealer := deal[lib.HF_GetRandom(len(deal))]
		dealer.Dealer = true
		dealer.TZ = 0
		dealer.ZZNum++
		if dealer.RobDeal <= 0 {
			dealer.RobDeal = 1
		}

		var msg staticfunc.Msg_Uid
		msg.Uid = dealer.Uid
		self.room.broadCastMsg("gamedealer", &msg)

		//! 下注
		self.State = 2
		self.SetTime(8)
	}

	self.room.flush()
}

//! 看牌
func (self *Game_MPQZ) GameOpen(uid int64) {
	if !self.room.Begin {
		return
	}

	if self.IsBeginCard() { //! 看牌抢庄没有这个步骤
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

	person.CT, person.CS = GetGoldNiuNiuScore(person.Card, self.GetParam2(TYPE_MPQZ_SJN) == 1, self.GetParam2(TYPE_MPQZ_ZDN) == 1, self.GetParam2(TYPE_MPQZ_HLN) == 1, self.GetParam2(TYPE_MPQZ_WXN) == 1, self.GetParam2(TYPE_MPQZ_THN) == 1, self.GetParam2(TYPE_MPQZ_WHN) == 1, self.GetParam2(TYPE_MPQZ_SZN) == 1)

	var msg Msg_GameGoldNN_Open
	msg.Card = person.Card
	msg.CT = person.CT
	self.room.SendMsg(person.Uid, "gamempqzopen", &msg)

	person.Open = true

	self.room.flush()
}

func (self *Game_MPQZ) GameTrust(uid int64, ok bool) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	person.Trust = ok

	self.room.flush()

	var msg Msg_GameDeal
	msg.Uid = uid
	msg.Ok = ok
	self.room.broadCastMsg("gametrust", &msg)
}

//! 亮牌
func (self *Game_MPQZ) GameView(uid int64, send bool) {
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
		person.CT, person.CS = GetGoldNiuNiuScore(person.Card, self.GetParam2(TYPE_MPQZ_SJN) == 1, self.GetParam2(TYPE_MPQZ_ZDN) == 1, self.GetParam2(TYPE_MPQZ_HLN) == 1, self.GetParam2(TYPE_MPQZ_WXN) == 1, self.GetParam2(TYPE_MPQZ_THN) == 1, self.GetParam2(TYPE_MPQZ_WHN) == 1, self.GetParam2(TYPE_MPQZ_SZN) == 1)
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
func (self *Game_MPQZ) GameReady(uid int64) {
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
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能重复准备")
				return
			}
			self.PersonMgr[i].Ready = true
		}

		if self.PersonMgr[i].Ready {
			num++
		}
	}

	beginnum := 0
	if self.room.Step == 0 {
		if self.GetParam(TYPE_MPQZ_ZDKS) == 1 {
			beginnum = 4
		} else if self.GetParam(TYPE_MPQZ_ZDKS) == 2 {
			beginnum = 5
		} else if self.GetParam(TYPE_MPQZ_ZDKS) == 3 {
			beginnum = 6
		}
	} else {
		beginnum = len(self.room.Uid)
	}
	if beginnum > 0 {
		if num == len(self.room.Uid) && num >= beginnum { //! 准备的人数达到游戏最小人数
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
			self.OnBegin()
			return
		}
	}

	var msg staticfunc.Msg_Uid
	msg.Uid = uid
	self.room.broadCastMsg("gameready", &msg)

	self.room.flush()
}

//! 开始游戏
func (self *Game_MPQZ) GameBegin(uid int64) {
	if self.room.IsBye() {
		return
	}

	if self.room.Begin { //! 已经开始了不允许准备
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经开始了，不能准备")
		return
	}

	if self.room.Host != uid {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "非房主不能开始游戏")
		return
	}

	if len(self.room.Uid) >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 准备的人数达到游戏最小人数
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
		self.OnBegin()
		return
	}
}

//! 下注
func (self *Game_MPQZ) GameBets(uid int64, bets int) {
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
		bettype := self.GetParam(TYPE_MPQZ_DF)
		if bets != bettype+1 && bets != (bettype+1)*2 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "下注错误")
			return
		}
		person.TZ = 0
	} else {
		person.TZNum++
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
		//! 亮牌
		self.State = 3
		self.SetTime(8)

		for i := 0; i < len(self.PersonMgr); i++ {
			if self.IsBeginCard() { //! 看牌抢庄
				//				card := self.Card.Deal(5 - len(self.PersonMgr[i].Card))
				//				self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, card...)

				//				var msg Msg_GameNiuNiuJX_Card
				//				msg.Card = card
				//				person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
				//				if person != nil {
				//					person.SendMsg("gamempqzcard", &msg)
				//				}

				//				card := self.Card.Deal(5 - len(self.PersonMgr[i].Card))
				//				self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, card...)
				card := self.PersonMgr[i].CardZB[4:]
				//card = append(card, self.PersonMgr[i].CardZB[4])
				self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, card...)

				var msg Msg_GameNiuNiuJX_Card
				msg.Card = card
				person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
				if person != nil {
					person.SendMsg("gamempqzcard", &msg)
				}

			} else { //! 自由抢庄
				//self.PersonMgr[i].Card = self.Card.Deal(5)
				self.PersonMgr[i].Card = self.PersonMgr[i].CardZB

				var msg Msg_GameNiuNiuJX_Card
				msg.Card = make([]int, 5)
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
func (self *Game_MPQZ) OnEnd() {
	self.room.SetBegin(false)
	self.State = 0
	self.SetTime(0)

	dealcs := 0
	var dealer *Game_MPQZ_Person = nil
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Dealer {
			dealer = self.PersonMgr[i]
			if self.room.Type == 69 { //! 牛牛上庄
				if dealer.CT == 100 {
					self.Dealer = dealer.Uid
					dealcs = dealer.CS
				}
			}
			break
		}
	}

	lst := make([]*Game_MPQZ_Person, 0)
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Ready = false
		if self.PersonMgr[i].Uid != dealer.Uid {
			lst = append(lst, self.PersonMgr[i])
		}
	}

	for i := 0; i < len(lst); i++ {
		if self.room.Type == 69 { //! 牛牛上庄
			if lst[i].CT == 100 && lst[i].CS > dealcs {
				self.Dealer = lst[i].Uid
				dealcs = lst[i].CS
			}
		}
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
			bs := GetMPQZBS(dealer.CT, self.GetParam(TYPE_MPQZ_FB))
			score := lst[i].Bets * dealer.RobDeal * bs
			dealer.CurScore += score
			lst[i].CurScore += -score
			lst[i].TZ = 0
		} else { //! 闲家赢
			bs := GetMPQZBS(lst[i].CT, self.GetParam(TYPE_MPQZ_FB))
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "闲家赢:", bs)
			score := lst[i].Bets * dealer.RobDeal * bs
			lst[i].CurScore += score
			dealer.CurScore += -score
			if lst[i].TZ == 0 {
				if self.GetParam(TYPE_MPQZ_XJTZ) != 0 {
					lst[i].TZ = lib.HF_MinInt(score+lst[i].Bets, lst[i].Bets*self.GetParam(TYPE_MPQZ_XJTZ)*5)
				}
			} else {
				lst[i].TZ = 0
			}
		}
		lst[i].Score += lst[i].CurScore
	}
	dealer.Score += dealer.CurScore

	//! 记录
	var record Rec_GameNiuNiuXY
	record.Time = time.Now().Unix()
	record.Roomid = self.room.Id*100 + self.room.Step
	record.MaxStep = self.room.MaxStep
	record.Param1 = self.room.Param1
	record.Param2 = self.room.Param2
	record.Host = self.room.Host

	//! 发消息
	var msg Msg_GameMPQZ_End
	for i := 0; i < len(self.PersonMgr); i++ {
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

		var rec Son_Rec_GameNiuNiuXY
		rec.Uid = self.PersonMgr[i].Uid
		rec.Name = self.room.GetName(self.PersonMgr[i].Uid)
		rec.Head = self.room.GetHead(self.PersonMgr[i].Uid)
		rec.Card = self.PersonMgr[i].Card
		rec.Dealer = self.PersonMgr[i].Dealer
		rec.Bets = self.PersonMgr[i].Bets
		rec.Score = self.PersonMgr[i].CurScore
		rec.Type = self.PersonMgr[i].CT
		rec.RobDeal = self.PersonMgr[i].RobDeal
		record.Info = append(record.Info, rec)

		self.PersonMgr[i].Ready = false
	}
	self.room.AddRecord(lib.HF_JtoA(&record))
	self.room.broadCastMsg("gamempqzend", &msg)

	if self.room.IsBye() || (self.room.Type == 68 && self.GetParam(TYPE_MPQZ_ZDQZ) > 1 && dealer.Score < 0) {
		self.OnBye()
		self.room.Bye()
		return
	}

	self.State = 0
	self.SetTime(6)

	self.room.flush()
}

func (self *Game_MPQZ) OnBye() {
	self.State = 0
	self.SetTime(0)

	info := make([]staticfunc.JS_CreateRoomMem, 0)
	var msg Msg_GameMPQZ_Bye
	msg.Host = self.room.Host
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.room.Type == 68 && self.room.Host == self.PersonMgr[i].Uid { //! 固定庄家
			if self.GetParam(TYPE_MPQZ_ZDQZ) == 2 {
				self.PersonMgr[i].Score -= 100
			} else if self.GetParam(TYPE_MPQZ_ZDQZ) == 3 {
				self.PersonMgr[i].Score -= 150
			} else if self.GetParam(TYPE_MPQZ_ZDQZ) == 4 {
				self.PersonMgr[i].Score -= 200
			}
		}

		var son Son_GameMPQZ_Bye
		son.Uid = self.PersonMgr[i].Uid
		son.Score = self.PersonMgr[i].Score
		son.QZ = self.PersonMgr[i].QZNum
		son.ZZ = self.PersonMgr[i].ZZNum
		son.TZ = self.PersonMgr[i].TZNum
		msg.Info = append(msg.Info, son)
		info = append(info, staticfunc.JS_CreateRoomMem{son.Uid, "", "", son.Score})
	}
	self.room.broadCastMsg("gamempqzbye", &msg)

	self.room.ClubResult(info)
}

func (self *Game_MPQZ) OnExit(uid int64) {
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

	beginnum := 0
	if self.GetParam(TYPE_MPQZ_ZDKS) == 1 {
		beginnum = 4
	} else if self.GetParam(TYPE_MPQZ_ZDKS) == 2 {
		beginnum = 5
	} else if self.GetParam(TYPE_MPQZ_ZDKS) == 3 {
		beginnum = 6
	}
	if beginnum > 0 {
		if num == len(self.room.Uid) && num >= beginnum { //! 准备的人数达到游戏最小人数
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
			self.OnBegin()
			return
		}
	}
}

func (self *Game_MPQZ) getInfo(uid int64) *Msg_GameMPQZ_Info {
	var msg Msg_GameMPQZ_Info
	msg.Begin = self.room.Begin
	msg.State = self.State
	if self.Time != 0 {
		msg.Time = self.Time - time.Now().Unix()
	}
	msg.Info = make([]Son_GameMPQZ_Info, 0)
	for _, value := range self.PersonMgr {
		var son Son_GameMPQZ_Info
		son.Uid = value.Uid
		son.Ready = value.Ready
		son.Bets = value.Bets
		son.Dealer = value.Dealer
		son.Total = value.Score
		son.Score = value.CurScore
		son.RobDeal = value.RobDeal
		son.TZ = value.TZ
		son.Trust = value.Trust
		son.View = value.View
		son.Open = value.Open
		if self.IsBeginCard() { //! 看牌抢庄
			if value.Uid == uid || value.View || !msg.Begin { //! 是自己或者亮牌了或者已经结束了
				son.Card = value.Card
				son.CT = value.CT
			} else {
				if GetServer().IsAdmin(uid, staticfunc.ADMIN_NIUNIU) {
					son.Card = value.Card
					son.CT = value.CT
				} else {
					son.Card = make([]int, len(value.Card))
					son.CT = 0
				}

			}
		} else { //! 自由抢庄
			if (value.Uid == uid && value.Open) || !msg.Begin {
				son.Card = value.Card
				son.CT = value.CT
			} else {
				if GetServer().IsAdmin(uid, staticfunc.ADMIN_NIUNIU) {
					son.Card = value.Card
					son.CT = value.CT
				} else {
					son.Card = make([]int, len(value.Card))
					son.CT = 0
				}

			}
		}

		msg.Info = append(msg.Info, son)
	}
	return &msg
}

func (self *Game_MPQZ) GetPerson(uid int64) *Game_MPQZ_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

func (self *Game_MPQZ) OnTime() {
	if self.Time == 0 {
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
				if self.PersonMgr[i].Bets <= 0 {
					self.GameBets(self.PersonMgr[i].Uid, self.GetParam(TYPE_MPQZ_DF)+1)
				}
			}
		} else if self.State == 3 { //! 亮牌
			for i := 0; i < len(self.PersonMgr); i++ {
				if !self.PersonMgr[i].View {
					self.GameView(self.PersonMgr[i].Uid, false)
				}
			}
		} else if self.State == 0 {
			for i := 0; i < len(self.PersonMgr); i++ {
				if !self.PersonMgr[i].Ready {
					self.GameReady(self.PersonMgr[i].Uid)
				}
			}
		}
	}
}

func (self *Game_MPQZ) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_MPQZ) OnIsBets(uid int64) bool {
	return false
}

//! 结算所有人
func (self *Game_MPQZ) OnBalance() {
}

//! 设置时间
func (self *Game_MPQZ) SetTime(t int) {
	if t == 0 {
		self.Time = 0
	} else {
		self.Time = time.Now().Unix() + int64(t)
	}

	var msg Msg_SetTime
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	self.room.broadCastMsg("gametime", &msg)
}

//! 是否开始发牌
func (self *Game_MPQZ) IsBeginCard() bool {
	return self.room.Type == 65 || self.room.Type == 66
}

//! 是否开始有庄
func (self *Game_MPQZ) IsBeginDealer() bool {
	return self.room.Type == 68 || self.room.Type == 69
}
