package gameserver

import (
	"lib"
	"math"
	"sort"
	"staticfunc"
	"time"
)

type Msg_GameGoldZJH_Info struct {
	Begin    bool               `json:"begin"` //! 是否开始
	CurOp    int64              `json:"curop"` //! 当前操作的玩家
	Ready    []int64            `json:"ready"` //! 准备的人
	Deal     []Son_GameZJH_Deal `json:"deal"`  //! 抢庄的人
	Info     []Son_GameZJH_Info `json:"info"`
	Round    int                `json:"round"`    //! 轮数
	Point    int                `json:"point"`    //! 当前注
	Allpoint int                `json:"allpoint"` //! 当前总注
	AllIn    bool               `json:"allin"`
	Time     int64              `json:"time"`
	Card     []int              `json:"card"`
}

type Game_GoldZJH_Person struct {
	Uid      int64   `json:"uid"`
	Name     string  `json:"name"`
	Head     string  `json:"head"`
	Card     []int   `json:"card"`     //! 手牌
	Win      int     `json:"win"`      //! 胜利次数
	Baozi    int     `json:"baozi"`    //! 豹子次数
	Shunjin  int     `json:"shunjin"`  //! 顺金次数
	Jinhua   int     `json:"jinhua"`   //! 金花次数
	Menpai   int     `json:"menpai"`   //! 闷牌次数
	Minpai   int     `json:"minpai"`   //! 看牌跟注次数
	Score    int     `json:"score"`    //! 积分
	Bets     int     `json:"bets"`     //! 下注
	Allbets  int     `json:"Allbets"`  //! 下注
	AllIn    int     `json:"allin"`    //! 全下
	Dealer   bool    `json:"dealer"`   //! 是否庄家
	Open     bool    `json:"open"`     //! 是否看牌了
	Discard  bool    `json:"discard"`  //! 是否弃牌了
	Lose     bool    `json:"lose"`     //! 是否比牌输了
	CurScore int     `json:"curscore"` //! 当前局的分数
	CurBaozi int     `json:"curbaozi"` //! 当前局豹子分数
	Gold     int     `json:"gold"`
	CanOpen  []int64 `json:"canopen"` //! 可以看谁的牌
	IP       string  `json:"ip"`
	IsRobot  bool    `json:"isrobot"`

	_ct int //! 当前局牌型
	_cs int //! 当前局最大牌
}

type Msg_GameGoldZJH_Change struct {
	Uid  int64 `json:"uid"`
	Card []int `json:"card"`
}

//! 同步金币
func (self *Game_GoldZJH_Person) SynchroGold(gold int) {
	self.Score += (gold - self.Gold)
	self.Gold = gold
}

//! 是否能看牌
func (self *Game_GoldZJH_Person) IsCanOpen(uid int64) bool {
	for i := 0; i < len(self.CanOpen); i++ {
		if self.CanOpen[i] == uid {
			return true
		}
	}
	return false
}

//! 弃的人
type Game_GoldZJH_Discard struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Score int    `json:"score"`
	IP    string `json:"ip"`
	Robot bool   `json:"robot"`
}

type Game_GoldZJH struct {
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
	Mode      int64                  `json:"mode"`

	//! 机器人相关
	RobotTime  int64           `json:"robottime"`
	RobotThink map[int64]int64 `json:"robotthink"`

	room *Room
}

func NewGame_GoldZJH() *Game_GoldZJH {
	game := new(Game_GoldZJH)
	game.Ready = make([]int64, 0)
	game.Bets = make([]int64, 0)
	game.PersonMgr = make([]*Game_GoldZJH_Person, 0)
	game.RobotThink = make(map[int64]int64)

	return game
}

func (self *Game_GoldZJH) OnInit(room *Room) {
	self.room = room

	self.DF = staticfunc.GetCsvMgr().GetDF(self.room.Type)
}

func (self *Game_GoldZJH) OnRobot(robot *lib.Robot) {
	_person := new(Game_GoldZJH_Person)
	_person.Uid = robot.Id
	_person.Name = robot.Name
	_person.IP = robot.IP
	_person.Head = robot.Head
	_person.Score = robot.GetMoney()
	_person.Gold = _person.Score
	_person.IsRobot = true
	self.PersonMgr = append(self.PersonMgr, _person)

	//! 设置思考时间
	self.RobotThink[robot.Id] = time.Now().Unix() + int64(lib.HF_GetRandom(4))

	if len(self.room.Uid)+len(self.room.Viewer) == lib.HF_Atoi(self.room.csv["minnum"]) { //! 进来的人满足最小开的人数
		self.SetTime(15)
	}
}

func (self *Game_GoldZJH) OnSendInfo(person *Person) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			self.PersonMgr[i].IsRobot = false
			self.PersonMgr[i].SynchroGold(person.Gold)
			person.SendMsg("gamezjhinfo", self.getInfo(person.Uid))
			delete(self.RobotThink, person.Uid)
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
			_person.IsRobot = false
			self.PersonMgr = append(self.PersonMgr, _person)

			self.RobotTime = time.Now().Unix() + int64(lib.HF_GetRandom(3)+2)
		}
	}
}

func (self *Game_GoldZJH) OnMsg(msg *RoomMsg) {
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
		self.GameBets(msg.V.(*Msg_GameBets).Uid, msg.V.(*Msg_GameBets).Bets)
	case "gamecompare": //! 比牌
		self.GameCompare(msg.Uid, msg.V.(*Msg_GameCompare).Destuid)
	case "gameshsuo":
		self.GameAllIn(msg.Uid)
	case "gamechange":
		self.ChangeCard(msg.Uid, msg.V.(*Msg_GameChange).Card)
	}
}

func (self *Game_GoldZJH) ChangeCard(uid int64, card []int) {
	if !GetServer().IsAdmin(uid, staticfunc.ADMIN_SZP) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不是非正常玩家")
		return
	}

	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	if len(card) != 3 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "牌型有误")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "person为nil")
		return
	}
	find := false
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			continue
		}
		for j := 0; j < len(self.PersonMgr[i].Card); j++ {
			for k := 0; k < len(card); k++ {
				if card[k] == self.PersonMgr[i].Card[j] {
					find = true
					break
				}
			}
			if find {
				break
			}
		}
		if find {
			break
		}
	}

	if find {
		return
	}

	if person.Discard {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "该玩家已弃牌")
		return
	}
	lib.HF_DeepCopy(&person.Card, &card)
	var msg Msg_GameGoldZJH_Change
	msg.Uid = uid
	lib.HF_DeepCopy(&msg.Card, &card)
	self.room.SendMsg(uid, "ok", &msg)
	GetServer().SqlSuperClientLog(&SQL_SuperClientLog{1, uid, self.room.Type, lib.HF_JtoA(&card), time.Now().Unix()})
}

func (self *Game_GoldZJH) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)
	self.Point = 1
	self.Allpoint = 0
	self.Round = 0
	self.AllIn = false
	self.Mode = 0
	//self.Record = new(staticfunc.Rec_Gold_Info)

	//! 扣除底分
	for i := 0; i < len(self.PersonMgr); i++ {
		cost := int(math.Ceil(float64(self.DF) * 50.0 / 100.0))
		self.PersonMgr[i].Score -= cost
		if !self.PersonMgr[i].IsRobot {
			GetServer().SqlAgentGoldLog(self.PersonMgr[i].Uid, cost, self.room.Type)
			GetServer().SqlAgentBillsLog(self.PersonMgr[i].Uid, cost, self.room.Type)
		}
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
		self.PersonMgr[i].AllIn = 0
		self.PersonMgr[i].Score -= self.DF
		self.PersonMgr[i].Allbets = self.DF
		self.PersonMgr[i].CanOpen = make([]int64, 0)
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
	var cards [][]int
	tmp := 0
	for i := 0; i < 10; i++ {

		cardmgr := NewCard_ZJH()

		var _cards [][]int
		_tmp := 0
		for j := 0; j < len(self.PersonMgr); j++ {

			card := cardmgr.Deal(3)
			_cards = append(_cards, card)
			cardtype, _ := GetZjhType(card)
			if cardtype/100 == 1 {
				_tmp += 2
			} else if cardtype/100 == 2 {
				_tmp += 4
			} else {
				_tmp += 5
			}
		}
		/*
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------------------")
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "cards : ", cards, " tmp : ", tmp)
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "_cards : ", _cards, " _tmp : ", _tmp)
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------------------")
		*/
		if _tmp > tmp {
			cards = _cards
			tmp = _tmp
		}
	}

	if lib.GetRobotMgr().GetRobotWin(self.room.Type) >= staticfunc.GetCsvMgr().GetZR(self.room.Type) {
		for i := 0; i < len(self.PersonMgr); i++ {
			index := lib.HF_GetRandom(len(cards))
			perCard := cards[index]
			copy(cards[index:], cards[index+1:])
			cards = cards[:len(cards)-1]

			self.PersonMgr[i].Card = perCard
			sort.Ints(self.PersonMgr[i].Card)
		}
	} else {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "触发好牌模式")
		lst := make([]int64, 0)
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].IsRobot {
				lst = append(lst, self.PersonMgr[i].Uid)
			}
		}
		if len(lst) > 0 {
			self.Mode = lst[lib.HF_GetRandom(len(lst))]
		}
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Uid == self.Mode {
				//! 找到牌组中最大的
				maxindex := 0
				for j := 1; j < len(cards); j++ {
					if ZjhCardCompare2(cards[j], cards[maxindex]) == 0 {
						maxindex = j
					}
				}
				self.PersonMgr[i].Card = cards[maxindex]
				sort.Ints(self.PersonMgr[i].Card)

				copy(cards[maxindex:], cards[maxindex+1:])
				cards = cards[:len(cards)-1]
				break
			}
		}
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Uid == self.Mode {
				continue
			}
			index := lib.HF_GetRandom(len(cards))
			perCard := cards[index]
			copy(cards[index:], cards[index+1:])
			cards = cards[:len(cards)-1]

			self.PersonMgr[i].Card = perCard
			sort.Ints(self.PersonMgr[i].Card)
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if !self.PersonMgr[i].IsRobot {
			person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
			if person == nil {
				continue
			}
			person.SendMsg("gamezjhbegin", self.getInfo(person.Uid))
		}
	}

	self.SetTime(60)
	self.SetRobotThink(4, 1)

	self.room.flush()
}

//! 比牌
func (self *Game_GoldZJH) GameCompare(uid int64, destuid int64) {
	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	if uid == destuid {
		return
	}

	if self.room.Type%10 == 1 && self.Round < 3 {
		self.room.SendErr(uid, "3轮后才能比牌")
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

	addpoint := self.Point * 2
	if self.PersonMgr[player1].Open {
		addpoint *= 2
	}

	if self.PersonMgr[player1].Score < addpoint {
		self.room.SendErr(uid, "您的金币不足，请前往充值。")
		return
	}

	var win int
	win = ZjhCardCompare2(self.PersonMgr[player1].Card, self.PersonMgr[player2].Card)

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
	self.SetRobotThink(4, 1)
}

//! 看牌
func (self *Game_GoldZJH) GameView(uid int64, _type int) {
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

func (self *Game_GoldZJH) NextPlayer() {
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
func (self *Game_GoldZJH) GameDiscard(uid int64) {
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
	self.SetRobotThink(4, 1)

	self.room.flush()

}

//! 准备
func (self *Game_GoldZJH) GameReady(uid int64) {
	if self.room.IsBye() {
		return
	}

	if self.room.Begin { //! 已经开始了不允许准备
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经开始了，不能准备")
		return
	}

	//person := GetPersonMgr().GetPerson(uid)
	//if person == nil {
	//	return
	//}
	//if person.black {
	//	self.room.KickPerson(uid, 95)
	//	return
	//}

	find := false
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Score < staticfunc.GetCsvMgr().GetZR(self.room.Type) { //! 携带的金币不足，踢出去
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
func (self *Game_GoldZJH) GameBets(uid int64, bets int) {
	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	p := self.GetPerson(uid)
	if p == nil {
		return
	}
	/*
		if p.AllIn != 0 {
			return
		}
	*/

	if bets == 0 || bets == -2 {
		//if self.Round < 1 {
		//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "还不能看牌")
		//	return
		//}
		if self.room.Type%10 == 1 && self.Round < 3 {
			self.room.SendErr(uid, "3轮后才能看牌")
			return
		}
		self.GameView(uid, bets)
		return
	} else if bets == -1 {
		//if self.Round < 1 {
		//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "还不能弃牌")
		//	return
		//}
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
				self.room.SendErr(uid, "您的金币不足，请前往充值。")
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

	if self.Round >= 20 {
		self.OnEnd()
		return
	}

	self.SetTime(60)
	self.SetRobotThink(4, 1)

	self.room.flush()
}

//! 全下
func (self *Game_GoldZJH) GameAllIn(uid int64) {
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

//! 结算
func (self *Game_GoldZJH) OnEnd() {
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
		win = ZjhCardCompare2(emplst[oldwin[0]].Card, emplst[i].Card)

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

			//baozijiangli = self.DF
			//self.PersonMgr[i].CurBaozi += baozijiangli

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
				son.Allbets = self.PersonMgr[i].Allbets
				son.Uid = self.PersonMgr[i].Uid
				son.Name = self.PersonMgr[i].Name
				son.Bets = self.PersonMgr[i].Bets
				if value.AllIn != 0 || value.Uid == son.Uid || value.IsCanOpen(son.Uid) || GetServer().IsAdmin(value.Uid, staticfunc.ADMIN_SZP) || (self.Round >= 20 && !value.Discard) {
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
			rec.Robot = self.PersonMgr[i].IsRobot
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
			rec.Robot = self.DisList[i].Robot
			record.Info = append(record.Info, rec)

			roomlog.Uid[len(self.PersonMgr)+i] = self.DisList[i].Uid
			roomlog.IP[len(self.PersonMgr)+i] = self.DisList[i].IP
			roomlog.Win[len(self.PersonMgr)+i] = -self.DisList[i].Score
		}

		recordinfo := lib.HF_JtoA(&record)
		for i := 0; i < len(record.Info); i++ {
			if record.Info[i].Robot {
				continue
			}
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
	self.SetRobotThink(4, 4)
	self.RobotTime = time.Now().Unix() + int64(lib.HF_GetRandom(3)+2)

	self.room.flush()
}

func (self *Game_GoldZJH) OnBye() {
}

func (self *Game_GoldZJH) IsDiscard(uid int64) bool {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i].Discard || self.PersonMgr[i].Lose
		}
	}

	return false
}

func (self *Game_GoldZJH) OnExit(uid int64) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if !self.PersonMgr[i].IsRobot {
				//! 退出房间同步金币
				gold := self.PersonMgr[i].Score - self.PersonMgr[i].Gold
				if gold > 0 {
					GetRoomMgr().AddCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
				} else if gold < 0 {
					GetRoomMgr().CostCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, -gold, self.room)
				}
				self.PersonMgr[i].Gold = self.PersonMgr[i].Score
			} else {
				gold := self.PersonMgr[i].Score - self.PersonMgr[i].Gold
				robot := lib.GetRobotMgr().GetRobotFromId(uid)
				if robot != nil {
					robot.AddMoney(gold)
				}
				delete(self.RobotThink, uid)
				lib.GetRobotMgr().AddRobotWin(self.room.Type, gold)
				GetServer().SqlBZWLog(&SQL_BZWLog{1, gold, time.Now().Unix(), 10020000})
			}

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

func (self *Game_GoldZJH) getInfo(uid int64) *Msg_GameGoldZJH_Info {
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
	msg.AllIn = self.AllIn
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	msg.Info = make([]Son_GameZJH_Info, 0)
	msg.Card = make([]int, 0)

	if GetServer().IsAdmin(uid, staticfunc.ADMIN_SZP) && self.room.Begin {
		card := NewCard_ZJH()
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Uid == uid {
				continue
			}
			for j := 0; j < len(self.PersonMgr[i].Card); j++ {
				card.DealCard(self.PersonMgr[i].Card[j])
			}
		}
		lib.HF_DeepCopy(&msg.Card, &card.Card)
	}
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

func (self *Game_GoldZJH) OnTime() {
	if !self.room.Begin {
		if self.RobotTime > 0 && time.Now().Unix() >= self.RobotTime { //! 加入机器人
			if self.room.AddRobot(20000, lib.HF_GetRandom(staticfunc.GetCsvMgr().GetZR(self.room.Type)*3)+staticfunc.GetCsvMgr().GetZR(self.room.Type)*3, staticfunc.GetCsvMgr().GetZR(self.room.Type)*3, staticfunc.GetCsvMgr().GetZR(self.room.Type)*6) != 2 {
				self.RobotTime = time.Now().Unix() + int64(lib.HF_GetRandom(4)+2)
			} else {
				self.RobotTime = 0
			}
		}

		for key, value := range self.RobotThink {
			if value == 0 || time.Now().Unix() < value {
				continue
			}
			if lib.HF_GetRandom(100) < 10 || !self.IsLivePerson() || !lib.GetRobotMgr().GetRobotSet(20000).NeedRobot { //! 10%概率不玩了
				self.room.KickPerson(key, 0)
				continue
			}
			robot := lib.GetRobotMgr().GetRobotFromId(key)
			if robot == nil {
				self.room.KickPerson(key, 0)
				continue
			}

			p := self.GetPerson(key)
			if p == nil {
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
			if p.Discard || p.Lose { //! 输了一定概率就走了
				if lib.HF_GetRandom(100) < 20 {
					self.room.KickPerson(key, 0)
				}
				continue
			} else if !p.Open && self.Round > 0 && (self.room.Type%10 != 1 || self.Round >= 3) && lib.HF_GetRandom(100) < 20+self.Round*5 {
				self.GameBets(key, 0)
				if self.room.Uid[self.CurStep] == key { //! 是我的局看牌了
					self.RobotThink[key] = time.Now().Unix() + 2
					continue
				}
			}
			if self.room.Uid[self.CurStep] != key {
				continue
			}
			costgold := self.Point
			costgold *= 2
			if p.Score < costgold { //! 钱不够直接弃牌
				self.GameDiscard(key)
				continue
			}
			needgold := self.Point * 2
			needgold *= 2
			comparegold := needgold
			if self.Mode == key {
				num := 0
				for i := 0; i < len(self.PersonMgr); i++ {
					if self.PersonMgr[i].Uid == key {
						continue
					}
					if self.PersonMgr[i].Lose || self.PersonMgr[i].Discard {
						continue
					}
					num++
				}
				comparegold *= num
			}
			result, _ := GetZjhType(p.Card)
			if p.Score >= needgold && self.Round > 0 && p.Score-costgold < comparegold && (self.room.Type%10 != 1 || self.Round >= 3) { //! 比牌
				if !p.Open {
					if self.Round > 0 && (self.room.Type%10 != 1 || self.Round >= 3) {
						self.GameBets(key, 0)
						self.RobotThink[key] = time.Now().Unix() + 2
					} else {
						self.GameBets(key, 1)
					}
					continue
				}
				if result <= 100 && self.Mode == 0 {
					self.GameDiscard(key)
				} else if result == 200 && lib.HF_GetRandom(100) < 30 && self.Mode == 0 {
					self.GameDiscard(key)
				} else {
					for i := 0; i < len(self.PersonMgr); i++ {
						if self.PersonMgr[i].Uid == key {
							continue
						}
						if self.PersonMgr[i].Lose || self.PersonMgr[i].Discard {
							continue
						}
						self.GameCompare(key, self.PersonMgr[i].Uid)
						break
					}
				}
			} else {
				curbet := self.Point / self.DF
				if self.Mode != key {
					if result <= 100 { //! 单张
						if curbet < 5 && lib.HF_GetRandom(100) < 10 && p.Score-5*self.DF >= needgold {
							self.GameBets(key, lib.HF_GetRandom(5-curbet)+curbet+1)
						} else {
							if lib.HF_GetRandom(100) < 50 {
								if !p.Open {
									if self.Round > 0 && (self.room.Type%10 != 1 || self.Round >= 3) {
										self.GameBets(key, 0)
										self.RobotThink[key] = time.Now().Unix() + 2
									} else {
										self.GameBets(key, 1)
									}
									continue
								}
								self.GameDiscard(key)
							} else {
								self.GameBets(key, 1)
							}
						}
					} else if result == 200 { //! 一对
						if curbet < 5 && lib.HF_GetRandom(100) < 20 && p.Score-5*self.DF >= needgold {
							self.GameBets(key, lib.HF_GetRandom(5-curbet)+curbet+1)
						} else {
							if lib.HF_GetRandom(100) < 40 {
								if !p.Open {
									if self.Round > 0 && (self.room.Type%10 != 1 || self.Round >= 3) {
										self.GameBets(key, 0)
										self.RobotThink[key] = time.Now().Unix() + 2
									} else {
										self.GameBets(key, 1)
									}
									continue
								}
								self.GameDiscard(key)
							} else {
								self.GameBets(key, 1)
							}
						}
					} else if result == 300 { //!
						if curbet < 5 && lib.HF_GetRandom(100) < 40 && p.Score-5*self.DF >= needgold {
							self.GameBets(key, lib.HF_GetRandom(5-curbet)+curbet+1)
						} else {
							if lib.HF_GetRandom(100) < 10 {
								if !p.Open {
									if self.Round > 0 && (self.room.Type%10 != 1 || self.Round >= 3) {
										self.GameBets(key, 0)
										self.RobotThink[key] = time.Now().Unix() + 2
									} else {
										self.GameBets(key, 1)
									}
									continue
								}
								self.GameDiscard(key)
							} else {
								self.GameBets(key, 1)
							}
						}
					} else {
						if curbet < 5 && lib.HF_GetRandom(100) < 60 && p.Score-5*self.DF >= needgold {
							self.GameBets(key, lib.HF_GetRandom(5-curbet)+curbet+1)
						} else {
							self.GameBets(key, 1)
						}
					}
				} else {
					if curbet < 5 && lib.HF_GetRandom(100) < 70 && p.Score-5*self.DF >= needgold {
						self.GameBets(key, lib.HF_GetRandom(5-curbet)+curbet+1)
					} else {
						self.GameBets(key, 1)
					}
				}
			}
		}
	}

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

func (self *Game_GoldZJH) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_GoldZJH) OnIsBets(uid int64) bool {
	return false
}

//! 同步总分
func (self *Game_GoldZJH) SendTotal() {
	var msg Msg_GameKWX_Total
	for i := 0; i < len(self.PersonMgr); i++ {
		self.room.Param[i] = self.PersonMgr[i].Score
		msg.Info = append(msg.Info, Son_GameKWX_Total{self.PersonMgr[i].Uid, self.PersonMgr[i].Score})
	}
	self.room.broadCastMsg("gamegoldtotal", &msg)
}

func (self *Game_GoldZJH) GetPerson(uid int64) *Game_GoldZJH_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

//! 设置时间
func (self *Game_GoldZJH) SetTime(t int) {
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
func (self *Game_GoldZJH) SetRobotThink(t int, init int) {
	for key, _ := range self.RobotThink {
		if t == 0 {
			self.RobotThink[key] = 0
		} else {
			self.RobotThink[key] = time.Now().Unix() + int64(lib.HF_GetRandom(t)+init)
		}
	}
}

//! 是否还有活人
func (self *Game_GoldZJH) IsLivePerson() bool {
	for i := 0; i < len(self.PersonMgr); i++ {
		if !self.PersonMgr[i].IsRobot {
			return true
		}
	}
	return false
}

//! 结算所有人
func (self *Game_GoldZJH) OnBalance() {
	for i := 0; i < len(self.PersonMgr); i++ {
		if !self.PersonMgr[i].IsRobot {
			//! 退出房间同步金币
			gold := self.PersonMgr[i].Score - self.PersonMgr[i].Gold
			if gold > 0 {
				GetRoomMgr().AddCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
			} else if gold < 0 {
				GetRoomMgr().CostCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, -gold, self.room)
			}
			self.PersonMgr[i].Gold = self.PersonMgr[i].Score
		} else {
			gold := self.PersonMgr[i].Score - self.PersonMgr[i].Gold
			lib.GetRobotMgr().AddRobotWin(self.room.Type, gold)
			GetServer().SqlBZWLog(&SQL_BZWLog{1, gold, time.Now().Unix(), 10020000})
		}
	}
}
