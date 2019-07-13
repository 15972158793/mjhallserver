package gameserver

import (
	"lib"
	"math"
	"sort"
	"staticfunc"
	"time"
)

//! 结算
type Msg_GameGoldPDK_End struct {
	Info []Son_GameGoldPDK_Info `json:"info"`
}

type Msg_GamePDK_Step struct {
	Uid     int64 `json:"uid"`     //! 哪个uid
	Cards   []int `json:"cards"`   //! 出的啥牌
	CurStep int64 `json:"curstep"` //! 下局谁出
}

//! 金币场跑得快
type Msg_GameGoldPDK_Info struct {
	Begin    bool                   `json:"begin"`    //! 是否开始
	Info     []Son_GameGoldPDK_Info `json:"info"`     //! 人的info
	CurStep  int64                  `json:"curstep"`  //! 这局谁出
	BefStep  int64                  `json:"befstep"`  //! 上局谁出
	LastCard []int                  `json:"lastcard"` //! 最后的牌
	Time     int64                  `json:"time"`
}

type Son_GameGoldPDK_Info struct {
	Uid   int64 `json:"uid"`
	Card  []int `json:"card"`
	Total int   `json:"total"`
	Score int   `json:"score"`
	Ready bool  `json:"ready"`
	Trust bool  `json:"trust"`
}

type Game_GoldPDK_Person struct {
	Uid      int64 `json:"uid"`
	Card     []int `json:"card"`  //! 手牌
	Ready    bool  `json:"ready"` //! 是否准备
	Trust    bool  `json:"trust"` //! 是否托管
	Total    int   `json:"total"`
	Gold     int   `json:"gold"`
	CurScore int   `json:"curscore"`
	Boom     int   `json:"boom"`
	IsRobot  bool  `json:"isrobot"`
}

func (self *Game_GoldPDK_Person) Init() {
	self.Card = make([]int, 0)
	self.Trust = false
	self.CurScore = 0
	self.Boom = 0
}

//! 同步金币
func (self *Game_GoldPDK_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

//! 设置托管
func (self *Game_GoldPDK_Person) SetTrust(trust bool) {
	if self.Trust == trust {
		return
	}

	self.Trust = trust

	var msg Msg_GameDeal
	msg.Uid = self.Uid
	msg.Ok = trust

	person := GetPersonMgr().GetPerson(self.Uid)
	if person != nil {
		person.SendMsg("gametrust", &msg)
	}
}

//! 是否托管
func (self *Game_GoldPDK_Person) IsTrush(_time int64) bool {
	if self.Trust {
		return true
	}

	if time.Now().Unix() < _time {
		return false
	}

	self.SetTrust(true)

	return true
}

type Game_GoldPDK struct {
	PersonMgr []*Game_GoldPDK_Person `json:"personmgr"`
	LastCard  []int                  `json:"lastcard"` //! 最后出的牌
	CurStep   int64                  `json:"curstep"`  //! 谁出牌
	BefStep   int64                  `json:"befstep"`  //! 上局谁出
	DF        int                    `json:"df"`       //! 底分
	Time      int64                  `json:"time"`     //! 自动选择时间
	BP        int64                  `json:"bp"`       //! 包赔id
	Card      int                    `json:"card"`

	//! 机器人相关
	RobotTime  int64           `json:"robottime"`
	RobotThink map[int64]int64 `json:"robotthink"`

	room *Room
}

func NewGame_GoldPDK() *Game_GoldPDK {
	game := new(Game_GoldPDK)
	game.PersonMgr = make([]*Game_GoldPDK_Person, 0)
	game.RobotThink = make(map[int64]int64)

	return game
}

func (self *Game_GoldPDK) GetPerson(uid int64) *Game_GoldPDK_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

//! 得到下一个uid
func (self *Game_GoldPDK) GetNextUid() int64 {
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

//! 得到上一个uid
func (self *Game_GoldPDK) GetBeforeUid() *Game_GoldPDK_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid != self.BefStep {
			continue
		}

		if i-1 >= 0 {
			return self.PersonMgr[i-1]
		} else {
			return self.PersonMgr[len(self.PersonMgr)-1]
		}
	}

	return nil
}

func (self *Game_GoldPDK) OnInit(room *Room) {
	self.room = room

	self.DF = staticfunc.GetCsvMgr().GetDF(self.room.Type)
}

func (self *Game_GoldPDK) OnRobot(robot *lib.Robot) {
	_person := new(Game_GoldPDK_Person)
	_person.Init()
	_person.Uid = robot.Id
	_person.Ready = false
	_person.Total = robot.GetMoney()
	_person.Gold = _person.Gold
	_person.Trust = false
	_person.IsRobot = true
	self.PersonMgr = append(self.PersonMgr, _person)

	//! 设置思考时间
	self.RobotThink[robot.Id] = time.Now().Unix() + int64(lib.HF_GetRandom(4))

	if len(self.room.Uid)+len(self.room.Viewer) == lib.HF_Atoi(self.room.csv["minnum"]) { //! 进来的人满足最小开的人数
		self.SetTime(15)
	}
}

func (self *Game_GoldPDK) OnSendInfo(person *Person) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			self.PersonMgr[i].IsRobot = false
			self.PersonMgr[i].SynchroGold(person.Gold)
			person.SendMsg("gamegoldpdkinfo", self.getInfo(person.Uid))
			delete(self.RobotThink, person.Uid)
			return
		}
	}

	_person := new(Game_GoldPDK_Person)
	_person.Init()
	_person.Uid = person.Uid
	_person.Ready = false
	_person.Total = person.Gold
	_person.Gold = person.Gold
	_person.Trust = false
	_person.IsRobot = false
	self.PersonMgr = append(self.PersonMgr, _person)

	self.RobotTime = time.Now().Unix() + int64(lib.HF_GetRandom(3)+2)

	if len(self.PersonMgr) >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 人满了
		self.SetTime(5)
	}

	person.SendMsg("gamegoldpdkinfo", self.getInfo(person.Uid))
}

func (self *Game_GoldPDK) OnMsg(msg *RoomMsg) {
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
	case "gamesteps": //! 出牌
		person := self.GetPerson(msg.Uid)
		if person != nil {
			person.SetTrust(false)
		}
		self.GameStep(msg.Uid, msg.V.(*Msg_GameSteps).Cards)
	case "gametrust": //! 托管
		self.GameTrust(msg.Uid, msg.V.(*Msg_GameDeal).Ok)
	}
}

func (self *Game_GoldPDK) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)

	//! 扣除底分
	for i := 0; i < len(self.PersonMgr); i++ {
		cost := int(math.Ceil(float64(self.DF) * 50.0 / 100.0))
		self.PersonMgr[i].Total -= cost
		if !self.PersonMgr[i].IsRobot { //! 机器人不能有抽水
			GetServer().SqlAgentGoldLog(self.PersonMgr[i].Uid, cost, self.room.Type)
			GetServer().SqlAgentBillsLog(self.PersonMgr[i].Uid, cost, self.room.Type)
		}
	}
	self.SendTotal()

	cardmgr := NewCard_GoldPDK()
	self.LastCard = make([]int, 0)
	self.BefStep = 0

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "当前奖池:", lib.GetRobotMgr().GetRobotWin(self.room.Type))
	if false { //lib.GetRobotMgr().GetRobotWin(self.room.Type) < 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "特殊模式")
		dealindex := lib.HF_GetRandom(len(self.PersonMgr))
		self.CurStep = self.PersonMgr[dealindex].Uid
		cardmgr.DealCard(11)
		cardmgr.DealCard(12)
		cardmgr.DealCard(13)
		cardmgr.DealCard(21)
		cardmgr.DealCard(34)
		first := true
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].IsRobot && first {
				self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, 11)
				self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, 12)
				self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, 13)
				self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, 21)
				first = false
			}
			if dealindex == i { //! 庄家先发一个34
				self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, 34)
			}
			self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, cardmgr.Deal(16-len(self.PersonMgr[i].Card))...)
		}
	} else {
		dealindex := lib.HF_GetRandom(len(self.PersonMgr))
		self.CurStep = self.PersonMgr[dealindex].Uid
		self.PersonMgr[dealindex].Card = append(self.PersonMgr[dealindex].Card, cardmgr.DealCard(34))
		self.PersonMgr[dealindex].Card = append(self.PersonMgr[dealindex].Card, cardmgr.Deal(15)...)
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Uid != self.CurStep {
				self.PersonMgr[i].Card = cardmgr.Deal(16)
			}
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gamegoldpdkbegin", self.getInfo(person.Uid))
	}

	self.SetTime(20)
	self.SetRobotThink(4, 1)
	self.room.flush()
}

//! 托管
func (self *Game_GoldPDK) GameTrust(uid int64, ok bool) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	person.SetTrust(ok)
}

func (self *Game_GoldPDK) OnEnd() {
	self.room.SetBegin(false)
	self.Time = 0

	//! 记录
	var record staticfunc.Rec_Gold_Info
	record.Time = time.Now().Unix()
	record.GameType = self.room.Type

	self.SetTime(10)

	var bp *Game_GoldPDK_Person = nil
	if self.BP != 0 {
		before := self.GetPerson(self.BP)
		if self.GetHasSomeCard(before.Card) {
			bp = before
		}
	}
	total := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if len(self.PersonMgr[i].Card) == 16 { //! 关门
			score := 0
			if bp != nil {
				score = lib.HF_MinInt(bp.Total, 16*self.DF*2)
				bp.CurScore -= score
			} else {
				score = lib.HF_MinInt(self.PersonMgr[i].Total, 16*self.DF*2)
				self.PersonMgr[i].CurScore -= score
			}
			total += score
		} else if len(self.PersonMgr[i].Card) > 1 { //! 输
			score := 0
			if bp != nil {
				score = lib.HF_MinInt(bp.Total, len(self.PersonMgr[i].Card)*self.DF)
				bp.CurScore -= score
			} else {
				score = lib.HF_MinInt(self.PersonMgr[i].Total, len(self.PersonMgr[i].Card)*self.DF)
				self.PersonMgr[i].CurScore -= score
			}
			total += score
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if len(self.PersonMgr[i].Card) == 0 { //! 赢了
			self.PersonMgr[i].CurScore += total
			break
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Total += self.PersonMgr[i].CurScore
		if self.PersonMgr[i].IsRobot {
			lib.GetRobotMgr().AddRobotWin(self.room.Type, self.PersonMgr[i].CurScore)
			GetServer().SqlBZWLog(&SQL_BZWLog{1, self.PersonMgr[i].CurScore, time.Now().Unix(), 10080000})
		}
	}

	//! 发消息
	var msg Msg_GameGoldPDK_End
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Ready = false

		var son Son_GameGoldPDK_Info
		son.Uid = self.PersonMgr[i].Uid
		son.Card = self.PersonMgr[i].Card
		son.Total = self.PersonMgr[i].Total
		son.Score = self.PersonMgr[i].CurScore
		msg.Info = append(msg.Info, son)

		var rec staticfunc.Son_Rec_Gold_Person
		rec.Uid = self.PersonMgr[i].Uid
		rec.Name = self.room.GetName(self.PersonMgr[i].Uid)
		rec.Head = self.room.GetHead(self.PersonMgr[i].Uid)
		rec.Score = self.PersonMgr[i].CurScore + self.PersonMgr[i].Boom
		rec.Robot = self.PersonMgr[i].IsRobot
		record.Info = append(record.Info, rec)

		self.room.Param[i] = self.PersonMgr[i].Total

		self.PersonMgr[i].Init()
	}
	recordinfo := lib.HF_JtoA(&record)
	for i := 0; i < len(record.Info); i++ {
		if record.Info[i].Robot {
			continue
		}
		GetServer().InsertRecord(self.room.Type, record.Info[i].Uid, recordinfo, -record.Info[i].Score)
	}
	self.room.broadCastMsg("gamegoldpdkend", &msg)

	//if self.room.IsBye() {
	//	self.OnBye()
	//	self.room.Bye()
	//	return
	//}

	self.room.flush()
}

func (self *Game_GoldPDK) OnBye() {
}

func (self *Game_GoldPDK) OnExit(uid int64) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if !self.PersonMgr[i].IsRobot {
				//! 退出房间同步金币
				gold := self.PersonMgr[i].Total - self.PersonMgr[i].Gold
				if gold > 0 {
					GetRoomMgr().AddCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
				} else if gold < 0 {
					GetRoomMgr().CostCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, -gold, self.room)
				}
				self.PersonMgr[i].Gold = self.PersonMgr[i].Total
			} else {
				gold := self.PersonMgr[i].Total - self.PersonMgr[i].Gold
				robot := lib.GetRobotMgr().GetRobotFromId(uid)
				if robot != nil {
					robot.AddMoney(gold)
				}
				delete(self.RobotThink, uid)
			}

			copy(self.PersonMgr[i:], self.PersonMgr[i+1:])
			self.PersonMgr = self.PersonMgr[:len(self.PersonMgr)-1]

			//! 有人退出之后取消自动操作
			self.SetTime(0)
			break
		}
	}
}

//! 准备,第一局自动准备
func (self *Game_GoldPDK) GameReady(uid int64) {
	if self.room.IsBye() {
		return
	}

	if self.room.Begin {
		return
	}

	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Ready {
				return
			} else {
				if self.PersonMgr[i].Total < staticfunc.GetCsvMgr().GetZR(self.room.Type) { //! 携带的金币不足，踢出去
					self.room.KickPerson(uid, 1)
					return
				}

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

	var msg staticfunc.Msg_Uid
	msg.Uid = uid
	self.room.broadCastMsg("gameready", &msg)

	self.room.flush()
}

//! 出牌(玩家选择)
func (self *Game_GoldPDK) GameStep(uid int64, card []int) {
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

		tmp := IsOkByGoldPDKCards(card, person.Card)
		if tmp == 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌错误:", card)
			return
		}
		if len(self.LastCard) != 0 && self.BefStep != uid {
			_tmp := IsOkByGoldPDKCards(self.LastCard, make([]int, 0))
			if tmp%100 == TYPE_CARD_ZHA { //! 炸弹
				if _tmp%100 == TYPE_CARD_ZHA {
					if CardCompare(tmp/10, _tmp/10) <= 0 {
						lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌更小:", card, ",", self.LastCard)
						return
					}
				}
			} else {
				//if len(card) != len(self.LastCard) {
				//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌张数不匹配:", card, ",", self.LastCard)
				//	return
				//}
				if tmp%100 != _tmp%100 { //! 类型不同
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌不匹配:", card, ",", self.LastCard)
					return
				}
				if CardCompare(tmp/10, _tmp/10) <= 0 {
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌更小:", card, ",", self.LastCard)
					return
				}
			}
		} else if len(card) == 1 {
			self.BP = uid
			self.Card = card[0]
		}
		self.LastCard = card
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

		if tmp%100 == TYPE_CARD_ZHA { //! 炸弹
			person.Total += self.DF * 10
			person.Boom += self.DF * 10
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Uid != person.Uid {
					self.PersonMgr[i].Total -= self.DF * 5
					self.PersonMgr[i].Boom -= self.DF * 5
				}
			}
			self.SendTotal()
		}
	} else {
		if len(self.LastCard) == 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "第一局不能跳过")
			return
		} else {
			if len(self.GetStepCard(person)) > 0 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "有牌必须管", self.LastCard, "...", person.Card)
				return
			}
		}
	}

	self.CurStep = self.GetNextUid()

	self.SetTime(20)
	self.SetRobotThink(4, 1)

	var msg Msg_GamePDK_Step
	msg.Uid = uid
	msg.Cards = card
	msg.CurStep = self.CurStep
	self.room.broadCastMsg("gamegoldpdkstep", &msg)

	if len(person.Card) == 0 { //! 牌出完了
		if self.BP == uid {
			self.BP = 0
		}
		self.OnEnd()
		return
	} else if self.BP != uid {
		self.BP = 0
	}

	self.room.flush()
}

func (self *Game_GoldPDK) getInfo(uid int64) *Msg_GameGoldPDK_Info {
	var msg Msg_GameGoldPDK_Info
	if self.Time != 0 {
		msg.Time = self.Time - time.Now().Unix()
	}
	msg.Begin = self.room.Begin
	msg.CurStep = self.CurStep
	msg.BefStep = self.BefStep
	msg.LastCard = self.LastCard
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameGoldPDK_Info
		son.Uid = self.PersonMgr[i].Uid
		if son.Uid == uid || !msg.Begin || GetServer().IsAdmin(uid, staticfunc.ADMIN_PDK) {
			//		if son.Uid == uid || !msg.Begin || self.PersonMgr[0].Uid == uid {
			son.Card = self.PersonMgr[i].Card
		} else {
			son.Card = make([]int, len(self.PersonMgr[i].Card))
		}
		son.Ready = self.PersonMgr[i].Ready
		son.Total = self.PersonMgr[i].Total
		son.Trust = self.PersonMgr[i].Trust
		son.Score = self.PersonMgr[i].CurScore
		msg.Info = append(msg.Info, son)
	}
	return &msg
}

//! 同步总分
func (self *Game_GoldPDK) SendTotal() {
	var msg Msg_GameKWX_Total
	for i := 0; i < len(self.PersonMgr); i++ {
		self.room.Param[i] = self.PersonMgr[i].Total
		msg.Info = append(msg.Info, Son_GameKWX_Total{self.PersonMgr[i].Uid, self.PersonMgr[i].Total})
	}
	self.room.broadCastMsg("gamegoldtotal", &msg)
}

func (self *Game_GoldPDK) OnTime() {
	if !self.room.Begin {
		if self.RobotTime > 0 && time.Now().Unix() >= self.RobotTime { //! 加入机器人
			if self.room.AddRobot(80000, lib.HF_GetRandom(staticfunc.GetCsvMgr().GetZR(self.room.Type)*3)+staticfunc.GetCsvMgr().GetZR(self.room.Type), staticfunc.GetCsvMgr().GetZR(self.room.Type), staticfunc.GetCsvMgr().GetZR(self.room.Type)*4) != 2 {
				self.RobotTime = time.Now().Unix() + int64(lib.HF_GetRandom(4)+2)
			} else {
				self.RobotTime = 0
			}
		}

		for key, value := range self.RobotThink {
			if value == 0 || time.Now().Unix() < value {
				continue
			}
			if lib.HF_GetRandom(100) < 10 || !self.IsLivePerson() || !lib.GetRobotMgr().GetRobotSet(80000).NeedRobot { //! 10%概率不玩了
				self.room.KickPerson(key, 0)
				self.RobotTime = time.Now().Unix() + int64(lib.HF_GetRandom(4)+2)
				continue
			}
			robot := lib.GetRobotMgr().GetRobotFromId(key)
			if robot == nil {
				self.room.KickPerson(key, 0)
				self.RobotTime = time.Now().Unix() + int64(lib.HF_GetRandom(4)+2)
				continue
			}

			p := self.GetPerson(key)
			if p == nil || p.Ready {
				continue
			}
			self.RobotThink[key] = 0
			self.GameReady(p.Uid)
		}
	} else {
		for key, value := range self.RobotThink {
			if value == 0 || time.Now().Unix() < value {
				continue
			}
			p := self.GetPerson(key)
			if p == nil {
				continue
			}
			self.RobotThink[key] = 0
			if key != self.CurStep {
				continue
			}
			self.GameStep(key, self.GetStepCard(self.GetPerson(key)))
		}
	}

	if self.Time == 0 {
		return
	}

	//! 60秒之后自动选择
	if !self.room.Begin {
		if time.Now().Unix() < self.Time {
			return
		}
		for i := 0; i < len(self.PersonMgr); {
			if !self.PersonMgr[i].Ready {
				self.room.KickPerson(self.PersonMgr[i].Uid, 98)
			} else {
				i++
			}
		}

		return
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == self.CurStep {
			if self.PersonMgr[i].IsTrush(self.Time) {
				self.GameStep(self.PersonMgr[i].Uid, self.GetStepCard(self.PersonMgr[i]))
			}
			return
		}
	}
}

func (self *Game_GoldPDK) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_GoldPDK) OnIsBets(uid int64) bool {
	return false
}

//! 设置时间
func (self *Game_GoldPDK) SetTime(t int) {
	if t == 0 {
		self.Time = 0
	} else {
		self.Time = time.Now().Unix() + int64(t)
	}

	var msg Msg_SetTime
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	self.room.broadCastMsg("gametime", &msg)
}

//! 设置所有机器人的思考时间
func (self *Game_GoldPDK) SetRobotThink(t int, init int) {
	for key, _ := range self.RobotThink {
		if t == 0 {
			self.RobotThink[key] = 0
		} else {
			self.RobotThink[key] = time.Now().Unix() + int64(lib.HF_GetRandom(t)+init)
		}
	}
}

//! 是否还有活人
func (self *Game_GoldPDK) IsLivePerson() bool {
	for i := 0; i < len(self.PersonMgr); i++ {
		if !self.PersonMgr[i].IsRobot {
			return true
		}
	}
	return false
}

//! 结算所有人
func (self *Game_GoldPDK) OnBalance() {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].IsRobot {
			continue
		}
		//! 退出房间同步金币
		gold := self.PersonMgr[i].Total - self.PersonMgr[i].Gold
		if gold > 0 {
			GetRoomMgr().AddCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
		} else if gold < 0 {
			GetRoomMgr().CostCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, -gold, self.room)
		}
		self.PersonMgr[i].Gold = self.PersonMgr[i].Total
	}
}

//! 得到比这个牌大的牌
func (self *Game_GoldPDK) GetBetterCard(person *Game_GoldPDK_Person, lastcard []int, special bool) []int {
	stepcard := make([]int, 0)

	tmp := IsOkByGoldPDKCards(lastcard, make([]int, 0))

	handcard := make(map[int]int)
	for i := 0; i < len(person.Card); i++ {
		handcard[person.Card[i]/10]++
	}

	one := make(LstPoker, 0)
	two := make(LstPoker, 0)
	three := make(LstPoker, 0)
	four := make(LstPoker, 0)
	for key, value := range handcard {
		if value == 1 {
			one = append(one, key*10)
			sort.Sort(LstPoker(one))
		} else if value == 2 {
			two = append(two, key*10)
			sort.Sort(LstPoker(two))
		} else if value == 3 {
			three = append(three, key*10)
			sort.Sort(LstPoker(three))
		} else if value == 4 {
			four = append(four, key*10)
			sort.Sort(LstPoker(four))
		}
	}

	if tmp%100 == TYPE_CARD_ONE { //! 单张
		for i := 0; i < len(one); i++ {
			if CardCompare(one[i], tmp/10) > 0 || special {
				stepcard = append(stepcard, self.GetHasCard(person.Card, 1, one[i])...)
				return stepcard
			}
		}
		sort.Sort(LstPoker(person.Card))
		for i := 0; i < len(person.Card); i++ {
			if CardCompare(person.Card[i], tmp/10) > 0 || special {
				stepcard = append(stepcard, person.Card[i])
				return stepcard
			}
		}
	}

	if tmp%100 == TYPE_CARD_TWO { //! 对子
		for i := 0; i < len(two); i++ {
			if CardCompare(two[i], tmp/10) > 0 || special {
				stepcard = append(stepcard, self.GetHasCard(person.Card, 2, two[i])...)
				return stepcard
			}
		}

		for i := 0; i < len(three); i++ {
			if CardCompare(three[i], tmp/10) > 0 || special {
				stepcard = append(stepcard, self.GetHasCard(person.Card, 2, three[i])...)
				return stepcard
			}
		}
	}

	if tmp%100 == TYPE_CARD_ZHA { //! 炸弹
		for i := 0; i < len(four); i++ {
			if CardCompare(four[i], tmp/10) > 0 || special {
				stepcard = append(stepcard, self.GetHasCard(person.Card, 4, four[i])...)
				return stepcard
			}
		}
	}

	if tmp%100 == TYPE_CARD_SAN { //! 3张
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找到牌型:3张")
		for i := 0; i < len(three); i++ {
			if CardCompare(three[i], tmp/10) > 0 || special {
				stepcard = append(stepcard, self.GetHasCard(person.Card, 3, three[i])...)
				lib.GetLogMgr().Output(lib.LOG_DEBUG, stepcard)
				if len(two) > 0 {
					stepcard = append(stepcard, self.GetHasCard(person.Card, 2, two[0])...)
					lib.GetLogMgr().Output(lib.LOG_DEBUG, stepcard)
					return stepcard
				} else if len(three) > 1 {
					for j := 0; j < len(three); j++ {
						if i == j {
							continue
						}
						stepcard = append(stepcard, self.GetHasCard(person.Card, 2, three[j])...)
						lib.GetLogMgr().Output(lib.LOG_DEBUG, stepcard)
						return stepcard
					}
				} else if len(four) > 0 {
					stepcard = append(stepcard, self.GetHasCard(person.Card, 2, four[0])...)
					lib.GetLogMgr().Output(lib.LOG_DEBUG, stepcard)
					return stepcard
				} else if len(one) > 0 {
					for j := 0; j < lib.HF_MinInt(2, len(one)); j++ {
						stepcard = append(stepcard, self.GetHasCard(person.Card, 1, one[j])...)
						lib.GetLogMgr().Output(lib.LOG_DEBUG, stepcard)
					}
					return stepcard
				}
				lib.GetLogMgr().Output(lib.LOG_DEBUG, stepcard)
				return stepcard
			}
		}
	}

	if tmp%100 == TYPE_CARD_SHUN { //! 顺子
		_card := make([]int, 0)
		_card = append(_card, one...)
		_card = append(_card, two...)
		_card = append(_card, three...)
		_card = append(_card, four...)
		sort.Sort(LstPoker(_card))
		sort.Sort(LstPoker(lastcard))
		if len(_card) >= len(lastcard) {
			index := -1
			for i := 0; i < len(_card); i++ {
				if CardCompare(_card[i], lastcard[0]) > 0 || special {
					index = i
					break
				}
			}
			if index >= 0 {
				for {
					if index+len(lastcard)-1 >= len(_card) {
						break
					}
					find := true
					for i := 1; i <= len(lastcard)-1; i++ {
						if _card[index+i]/10 == 2 || _card[index+i]/10 == 14 {
							find = false
							break
						}
						if _card[index+i]/10 == 1 {
							_card[index+i] = 140
						}
						if _card[index+i]/10-_card[index+i-1]/10 != 1 {
							find = false
							break
						}
					}
					if find {
						for i := index; i < index+len(lastcard); i++ {
							stepcard = append(stepcard, self.GetHasCard(person.Card, 1, _card[i])...)
						}
						return stepcard
					}
					index++
				}
			}
		}
	}

	if tmp%100 == TYPE_CARD_SHUNDUI { //! 顺对
		_card := make([]int, 0)
		_card = append(_card, two...)
		_card = append(_card, three...)
		_card = append(_card, four...)
		sort.Sort(LstPoker(_card))
		sort.Sort(LstPoker(lastcard))
		if len(_card) >= len(lastcard)/2 {
			index := -1
			for i := 0; i < len(_card); i++ {
				if CardCompare(_card[i], lastcard[0]) > 0 || special {
					index = i
					break
				}
			}
			if index >= 0 {
				for {
					if index+len(lastcard)/2-1 >= len(_card) {
						break
					}
					find := true
					for i := 1; i <= len(lastcard)/2-1; i++ {
						if _card[index+i]/10 == 2 || _card[index+i]/10 == 14 {
							find = false
							break
						}
						if _card[index+i]/10 == 1 {
							_card[index+i] = 140
						}
						if _card[index+i]/10-_card[index+i-1]/10 != 1 {
							find = false
							break
						}
					}
					if find {
						for i := index; i < index+len(lastcard)/2; i++ {
							stepcard = append(stepcard, self.GetHasCard(person.Card, 2, _card[i])...)
						}
						return stepcard
					}
					index++
				}
			}
		}
	}

	if tmp%100 == TYPE_CARD_SHUNSAN { //! 顺3
		_card := make([]int, 0)
		_card = append(_card, three...)
		_card = append(_card, four...)
		sort.Sort(LstPoker(_card))
		last := IsOkByGoldPDKCards(lastcard, make([]int, 0))
		if len(_card) >= len(lastcard)/3 {
			index := -1
			for i := 0; i < len(_card); i++ {
				if CardCompare(_card[i], last/10) > 0 || special {
					index = i
					break
				}
			}
			if index >= 0 {
				for {
					if index+len(lastcard)/3-1 >= len(_card) {
						break
					}
					find := true
					for i := 1; i <= len(lastcard)/3-1; i++ {
						if _card[index+i]/10 == 2 || _card[index+i]/10 == 14 {
							find = false
							break
						}
						if _card[index+i]/10 == 1 {
							_card[index+i] = 140
						}
						if _card[index+i]/10-_card[index+i-1]/10 != 1 {
							find = false
							break
						}
					}
					if find {
						for i := index; i < index+len(lastcard)/3; i++ {
							stepcard = append(stepcard, self.GetHasCard(person.Card, 3, _card[i])...)
						}
						return stepcard
					}
					index++
				}
			}
		}
	}

	if tmp%100 == TYPE_CARD_SI1 { //! 4带
		for i := 0; i < len(four); i++ {
			if CardCompare(four[i], tmp/10) > 0 || special {
				stepcard = append(stepcard, self.GetHasCard(person.Card, 4, four[i])...)
				if len(two) > 0 {
					stepcard = append(stepcard, self.GetHasCard(person.Card, 2, two[0])...)
					return stepcard
				} else if len(three) > 0 {
					stepcard = append(stepcard, self.GetHasCard(person.Card, 2, three[0])...)
					return stepcard
				} else if len(four) > 1 {
					for j := 0; j < len(four); j++ {
						if i == j {
							continue
						}
						stepcard = append(stepcard, self.GetHasCard(person.Card, 2, four[j])...)
						return stepcard
					}
				} else if len(one) > 1 {
					for j := 0; j < 2; j++ {
						stepcard = append(stepcard, self.GetHasCard(person.Card, 1, one[j])...)
					}
					return stepcard
				}
				return stepcard
			}
		}
	}

	if tmp%100 != TYPE_CARD_ZHA { //! 不是炸弹
		for i := 0; i < len(four); i++ {
			stepcard = append(stepcard, self.GetHasCard(person.Card, 4, four[i])...)
			return stepcard
		}
	}

	return stepcard
}

//! 得到有没有牌大于上家
func (self *Game_GoldPDK) GetStepCard(person *Game_GoldPDK_Person) []int {
	if len(self.LastCard) == 0 || self.BefStep == person.Uid { //! 上家没有出牌
		if len(self.GetPerson(self.GetNextUid()).Card) == 1 { //! 下家保单
			//! 试试顺子能不能出
			shuncard := self.GetBetterCard(person, []int{31, 41, 51, 61, 71, 81, 91, 101}, true)
			if len(shuncard) > 0 && !self.IsBoom(shuncard) {
				return shuncard
			}
			shuncard = self.GetBetterCard(person, []int{31, 41, 51, 61, 71, 81, 91}, true)
			if len(shuncard) > 0 && !self.IsBoom(shuncard) {
				return shuncard
			}
			shuncard = self.GetBetterCard(person, []int{31, 41, 51, 61, 71, 81}, true)
			if len(shuncard) > 0 && !self.IsBoom(shuncard) {
				return shuncard
			}
			shuncard = self.GetBetterCard(person, []int{31, 41, 51, 61, 71}, true)
			if len(shuncard) > 0 && !self.IsBoom(shuncard) {
				return shuncard
			}
			//! 试试连对
			lianduicard := self.GetBetterCard(person, []int{31, 31, 41, 41}, true)
			if len(lianduicard) > 0 && !self.IsBoom(lianduicard) {
				return lianduicard
			}
			//! 试试三带二
			sandaiercard := self.GetBetterCard(person, []int{31, 31, 31, 41, 41}, true)
			if len(sandaiercard) > 0 && !self.IsBoom(sandaiercard) {
				return sandaiercard
			}
			//! 试试三带一
			sandaiyicard := self.GetBetterCard(person, []int{31, 31, 31, 41}, true)
			if len(sandaiyicard) > 0 && !self.IsBoom(sandaiyicard) {
				return sandaiyicard
			}
			//! 试试对子
			duicard := self.GetBetterCard(person, []int{31, 31}, true)
			if len(duicard) > 0 && !self.IsBoom(duicard) {
				return duicard
			}
			//! 炸弹
			boomcard := self.GetBetterCard(person, []int{21, 21}, true)
			if len(boomcard) > 0 {
				return boomcard
			}
			//! 得到最大的牌
			return []int{self.GetMaxCard(person.Card)}
		} else {
			canstep := make([][]int, 0)
			//! 试试顺子能不能出
			shuncard := self.GetBetterCard(person, []int{31, 41, 51, 61, 71, 81, 91, 101}, true)
			if len(shuncard) > 0 && !self.IsBoom(shuncard) && !self.GetBigCard(shuncard) {
				canstep = append(canstep, shuncard)
			}
			shuncard = self.GetBetterCard(person, []int{31, 41, 51, 61, 71, 81, 91}, true)
			if len(shuncard) > 0 && !self.IsBoom(shuncard) && !self.GetBigCard(shuncard) {
				canstep = append(canstep, shuncard)
			}
			shuncard = self.GetBetterCard(person, []int{31, 41, 51, 61, 71, 81}, true)
			if len(shuncard) > 0 && !self.IsBoom(shuncard) && !self.GetBigCard(shuncard) {
				canstep = append(canstep, shuncard)
			}
			shuncard = self.GetBetterCard(person, []int{31, 41, 51, 61, 71}, true)
			if len(shuncard) > 0 && !self.IsBoom(shuncard) && !self.GetBigCard(shuncard) {
				canstep = append(canstep, shuncard)
			}
			//! 试试连对
			lianduicard := self.GetBetterCard(person, []int{31, 31, 41, 41}, true)
			if len(lianduicard) > 0 && !self.IsBoom(lianduicard) && !self.GetBigCard(lianduicard) {
				canstep = append(canstep, lianduicard)
			}
			//! 试试三带二
			sandaiercard := self.GetBetterCard(person, []int{31, 31, 31, 41, 41}, true)
			if len(sandaiercard) > 0 && !self.IsBoom(sandaiercard) && !self.GetBigCard(sandaiercard) {
				canstep = append(canstep, sandaiercard)
			}
			//! 试试三带一
			sandaiyicard := self.GetBetterCard(person, []int{31, 31, 31, 41}, true)
			if len(sandaiyicard) > 0 && !self.IsBoom(sandaiyicard) && !self.GetBigCard(sandaiyicard) {
				canstep = append(canstep, sandaiyicard)
			}
			//! 试试对子
			duicard := self.GetBetterCard(person, []int{31, 31}, true)
			if len(duicard) > 0 && !self.IsBoom(duicard) && !self.GetBigCard(duicard) {
				canstep = append(canstep, duicard)
			}
			onecard := self.GetBetterCard(person, []int{31}, true)
			if len(onecard) > 0 && !self.IsBoom(onecard) && !self.GetBigCard(onecard) {
				canstep = append(canstep, onecard)
			}
			if len(canstep) > 0 {
				return canstep[lib.HF_GetRandom(len(canstep))]
			} else {
				return []int{person.Card[0]}
			}
		}
	}

	return self.GetBetterCard(person, self.LastCard, false)
}

func (self *Game_GoldPDK) GetHasCard(card []int, num int, key int) []int {
	tmp := make([]int, 0)
	for i := 0; i < len(card); i++ {
		if CardCompare(card[i], key) == 0 {
			tmp = append(tmp, card[i])
			if len(tmp) >= num {
				break
			}
		}
	}

	return tmp
}

func (self *Game_GoldPDK) GetHasSomeCard(card []int) bool {
	handcard := make(map[int]int)
	for i := 0; i < len(card); i++ {
		handcard[card[i]/10]++
	}

	for _, value := range handcard {
		if value >= 2 {
			return true
		}
	}

	for i := 0; i < len(card); i++ {
		if CardCompare(card[i], self.Card) > 0 {
			return true
		}
	}

	return false
}

//! 有没有大牌
func (self *Game_GoldPDK) GetBigCard(card []int) bool {
	for i := 0; i < len(card); i++ {
		if card[i]/10 == 12 || card[i]/10 == 13 || card[i]/10 == 1 || card[i]/10 == 2 || card[i] == 1000 || card[i] == 2000 || card[i]/10 == 14 || card[i]/10 == 20 {
			return true
		}
	}
	return false
}

//! 是否是炸弹
func (self *Game_GoldPDK) IsBoom(card []int) bool {
	if len(card) != 4 {
		return false
	}
	return card[0]/10 == card[1]/10 && card[0]/10 == card[2]/10 && card[0]/10 == card[3]/10
}

//! 得到最大的牌
func (self *Game_GoldPDK) GetMaxCard(card []int) int {
	maxcard := 0
	for i := 0; i < len(card); i++ {
		if maxcard == 0 {
			maxcard = card[i]
			continue
		}
		if CardCompare(card[i], maxcard) > 0 {
			maxcard = card[i]
		}
	}

	return maxcard
}
