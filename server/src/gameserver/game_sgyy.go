package gameserver

import (
	//"log"
	"lib"
	"math/rand"
	"staticfunc"
	"time"
)

type Game_SGYY struct {
	PersonMgr []*Game_SGYY_Person `json:"personmgr"`
	Ready     []int64             `json:"ready"` //! 已经准备的人
	Cardmgr   *CardMgr            //！扑克牌组
	TickTime  float64             `json:"ticktime"`  //! 计时器时间
	TimeState int                 `json:"timestate"` //! 计时器状态 1 抢庄 2下注 3明牌 4准备
	RecordMgr []Rec_GameNiuNiuJX  `json:"recordmgr"`
	room      *Room
}

//每个玩家的数据总览
type Game_SGYY_Person struct {
	Uid       int64 `json:"uid"`
	Deal      bool  `json:"deal"`      //!庄家
	Ready     bool  `json:"ready"`     //! 是否准备好
	Rub       bool  `json:"rub"`       //! 是否已经搓牌操作
	View      bool  `json:"view"`      //! 是否亮牌
	Card      []int `json:"card"`      //! 手里的牌
	Score     int   `json:"score"`     //! 单局输赢分数
	TuiBets   int   `json:"tuibets"`   //! 推注分数
	Bets      int   `json:"bets"`      //! 当局压点
	ModBets   int   `json:"modets"`    //! 模式3
	DealMul   int   `json:"dealmul"`   //! 抢庄倍率
	CardMul   int   `json:"cardmul"`   //! 牌型倍率
	CardValue int   `json:"cardvalue"` //! 牌值
	JoinStep  int   `json:"joinstep"`  //! 加入局数
	//CardPX    int   `json:"cardpx"`    //! 牌型
	//GongNum int `json:"gongnum"` //! 公牌数量
	Total int `json:"total"` //! 输赢分数总和
}

//每一局玩家的信息
type Son_GameSGYY_Info struct {
	Uid     int64 `json:"uid"`
	Dealer  bool  `json:"dealer"`
	Ready   bool  `json:"ready"`   //是否准备
	TuiBets int   `json:"tuibets"` //! 推注分数
	Card    []int `json:"card"`    //！牌
	Score   int   `json:"score"`   //! 当局得分
	Bets    int   `json:"bets"`    //! 当局压点
	Total   int   `json:"total"`   //玩家总分
	RobDeal int   `json:"robdeal"` //! 抢庄倍率
	View    bool  `json:"view"`    //是否亮牌
}
type Msg_GameSGYY_Info struct {
	Begin bool                `json:"begin"` //! 是否开始
	State int                 `json:"state"` //! 是否开始
	Host  int64               `json:"host"`  //! 房主UID
	Ready []int64             `json:"ready"` //! 已经准备的人
	Info  []Son_GameSGYY_Info `json:"info"`  //! 人的info
}

//! 单局结算
type Msg_GameSGYY_End struct {
	Info []Son_GameSGYY_End `json:"info"`
}
type Son_GameSGYY_End struct {
	Uid   int64 `json:"uid"`
	Score int   `json:"score"` //! 当局得分
	Total int   `json:"total"` //玩家总分
}
type Rec_GameSGYY struct {
	Info     []Son_Rec_GameSGYY `json:"info"`
	Roomid   int                `json:"roomid"`
	Time     int64              `json:"time"`
	RoomStep int                `json:"roomstep"`
}
type Son_Rec_GameSGYY struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Bets  int    `json:"bets"`
	Deal  bool   `json:"dealer"`
	Score int    `json:"score"`
}

func NewGame_SGYY() *Game_SGYY {
	game := new(Game_SGYY)
	game.PersonMgr = make([]*Game_SGYY_Person, 0)
	game.RecordMgr = make([]Rec_GameNiuNiuJX, 0)
	game.Ready = make([]int64, 0)
	game.Cardmgr = new(CardMgr)
	game.TimeState = 0
	return game
}
func (self *Game_SGYY) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "gamestart": //！游戏开始
		self.GameStart(msg.Uid)
	case "gameseat": //! 游戏坐下
		self.GameSeat(msg.Uid)
	case "gameready": //! 游戏准备
		self.GameReady(msg.Uid)
	case "gamedealer": //! 抢庄
		self.GameDeal(msg.Uid, msg.V.(*Msg_GameDealer).Score)
	case "gamebets": //! 下注
		self.GameBets(msg.Uid, msg.V.(*Msg_GameBets).Bets)
	case "gameview": //！亮牌
		self.GameView(msg.Uid)
	case "gameredeal": //！下庄 仅限于模式3
		self.GameReDeal(msg.Uid)
	}
}

func (self *Game_SGYY) GameStart(uid int64) {
	if self.room.Host != uid {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不是房主")
		return
	}
	if self.room.Step != 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不是第一局")
		return
	}
	if len(self.Ready) < lib.HF_Atoi(self.room.csv["minnum"]) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "人数不够")
		return
	}
	if len(self.Ready) != len(self.PersonMgr) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "有人未准备，不能开始")
		return
	}
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
	self.TimeState = 0
	self.OnBegin()
	self.room.flush()
}

//! 坐下
func (self *Game_SGYY) GameSeat(uid int64) {
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

	_person := new(Game_SGYY_Person)
	_person.Uid = uid
	_person.DealMul = -1
	_person.JoinStep = self.room.Step
	self.PersonMgr = append(self.PersonMgr, _person)
	self.room.broadCastMsg("gamesgyyseat", self.getInfo(0))
	self.GameReady(_person.Uid)
	self.room.flush()
}
func (self *Game_SGYY) GameReDeal(uid int64) {
	if 3 != self.room.Param1/10%10 || 4 == self.room.Param1/100%10 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不是模式3或者下庄分数为0")
		return
	}
	person := self.GetPerson(uid)
	if nil == person {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "玩家不存在")
		return
	}
	if !person.Deal {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "只有庄家能下庄")
		return
	}
	reScore := (self.room.Param1 / 100 % 10) * 100 //！下庄分数
	max := reScore
	maxArr := make([]int, 0)
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			continue
		}
		if self.PersonMgr[i].Total > max {
			max = self.PersonMgr[i].Total
			maxArr = make([]int, 0)
			maxArr = append(maxArr, i)
		} else if self.PersonMgr[i].Total == max {
			maxArr = append(maxArr, i)
		}
	}
	if len(maxArr) <= 0 { //!总结算
		self.OnBye()
		self.room.Bye()
	} else {
		self.PersonMgr[maxArr[0]].Deal = true
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Uid != uid {
				continue
			}
			self.PersonMgr[i].Deal = false
		}
	}
	self.GameReady(uid)
	self.room.flush()
}
func (self *Game_SGYY) GameView(uid int64) {
	person := self.GetPerson(uid)
	if nil == person {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "玩家不存在")
		return
	}
	if person.View {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "该玩家已经亮牌")
		return
	}
	person.View = true

	var msg Msg_GameView
	msg.Uid = person.Uid
	msg.Card = person.Card
	self.room.broadCastMsg("gameview", &msg)

	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].View {
			num++
		}
	}
	if num >= len(self.PersonMgr) {
		//self.ChooseView = false
		self.TimeState = 0
		self.OnEnd()
	}

	self.room.flush()
}

//！抢庄
func (self *Game_SGYY) GameDeal(uid int64, score int) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "未开始不能抢庄")
		return
	}
	if score < 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "错误的抢庄分数")
		return
	}
	if 3 == (self.room.Param1/10%10) || 4 == (self.room.Param1/10%10) { //！12才有抢庄
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "模式34没有抢庄")
		return
	}
	person := self.GetPerson(uid)
	if nil == person {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "玩家不存在")
		return
	}
	if person.DealMul >= 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "重复抢庄")
		return
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid != uid {
			continue
		}
		self.PersonMgr[i].DealMul = score //模式1发倍率  模式2为0不抢 1抢
	}
	//! 广播
	var msg Msg_GameDealer
	msg.Uid = uid
	msg.Score = score
	self.room.broadCastMsg("gamedealmul", &msg)
	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].DealMul >= 0 {
			num++
		}
	}
	if num >= len(self.PersonMgr) {
		self.TimeState = 0
		self.GameGetDealer()
	}
	self.room.flush()
}
func (self *Game_SGYY) GameReady(uid int64) {
	if self.room.IsBye() {
		return
	}

	//	if !self.GameSeat(uid) {
	//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "坐下失败")
	//		return
	//	}

	//设置玩家准备状态为true
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid != uid {
			continue
		}
		self.PersonMgr[i].Ready = true
		break
	}
	if self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经开始了，不能准备")
		return
	}
	for i := 0; i < len(self.Ready); i++ {
		if self.Ready[i] == uid {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "同一个玩家准备")
			return
		}
	}
	self.Ready = append(self.Ready, uid) //准备的玩家添加到数组

	var msg staticfunc.Msg_Uid
	msg.Uid = uid
	self.room.broadCastMsg("gameready", &msg)

	if len(self.Ready) == len(self.room.Uid) && len(self.Ready) >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 准备的人数达到游戏最小人数
		if self.room.Step == 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "第一局不开始")
			return
		}
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
		self.TimeState = 0
		self.OnBegin()
		return
	}
	self.room.flush()
}
func (self *Game_SGYY) GameBets(uid int64, bets int) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}
	if bets <= 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "下注无效")
		return
	}
	person := self.GetPerson(uid)
	if nil == person {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "玩家不存在")
		return
	}
	if person.Deal {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "庄家不能下注")
		return
	}
	if person.Bets > 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "重复下注")
		return
	}
	person.Bets = bets
	person.ModBets = bets
	var msg Msg_GameBets
	msg.Uid = person.Uid
	msg.Bets = bets
	self.room.broadCastMsg("gamebets", &msg)

	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Bets > 0 {
			num++
		}
	}
	if num >= len(self.PersonMgr)-1 && self.room.Param1/10%10 != 4 {
		self.TimeState = 0
		self.GameSendCard()
	}
	if num >= len(self.PersonMgr) && self.room.Param1/10%10 == 4 {
		self.TimeState = 0
		self.GameSendCard()
	}

	self.room.flush()
}
func (self *Game_SGYY) OnInit(room *Room) {
	self.room = room
}
func (self *Game_SGYY) OnBegin() {
	if self.room.IsBye() {
		return
	}
	self.room.SetBegin(true)

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "------ room.uid : ", self.room.Uid)
	for i := 0; i < len(self.PersonMgr); i++ {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------personmgr : ", self.PersonMgr[i].Uid)
	}

	for i := 0; i < len(self.room.Uid); i++ { //! 重新初始化人
		self.PersonMgr[i].Bets = 0
		self.PersonMgr[i].ModBets = 0
		self.PersonMgr[i].Score = 0
		self.PersonMgr[i].View = false
		self.PersonMgr[i].DealMul = -1
		self.PersonMgr[i].CardMul = 0
		self.PersonMgr[i].Ready = false
		if 3 != (self.room.Param1 / 10 % 10) { //！模式3庄家提前判断
			self.PersonMgr[i].Deal = false
		}
		self.PersonMgr[i].Card = make([]int, 0)
	}

	if 1 == self.room.Step && 3 == (self.room.Param1/10%10) { //！模式3第一局庄家为房主
		person := self.GetPerson(self.room.Uid[0])
		if nil == person {
			return
		}
		person.Deal = true
		if 4 != (self.room.Param1 / 100 % 10) {
			person.Total = (self.room.Param1 / 100 % 10) * 100
		}

	}

	self.Cardmgr = NewCard_LYC()
	//！明牌抢庄模式先发牌
	if 1 == (self.room.Param1 / 10 % 10) {
		for i := 0; i < len(self.PersonMgr); i++ { //每人先发两张牌
			self.PersonMgr[i].Card = self.Cardmgr.Deal(2)
		}
	}

	if 1 == (self.room.Param1/10%10) || 2 == (self.room.Param1/10%10) {
		self.TimeState = 1 //！抢庄计时开始
	} else if 3 == (self.room.Param1/10%10) || 4 == (self.room.Param1/10%10) {
		self.TimeState = 2 //！下注计时开始
	}
	self.TickTime = 0
	//！把消息发给玩家
	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gamebegin", self.getInfo(person.Uid))
	}
	self.room.broadCastMsgView("gamebegin", self.getInfo(0))

	self.room.flush()
}
func (self *Game_SGYY) OnBye() {
	var msg Msg_GameSGYY_End
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameSGYY_End
		son.Uid = self.PersonMgr[i].Uid
		son.Total = self.PersonMgr[i].Total
		msg.Info = append(msg.Info, son)
	}
	self.room.broadCastMsg("gamesgyybye", &msg)
}

func (self *Game_SGYY) OnEnd() {
	self.room.SetBegin(false)

	self.Ready = make([]int64, 0)

	gameModol := self.room.Param1 / 10 % 10 //！游戏模式
	speStyle := self.room.Param1 / 100000 % 10
	for i := 0; i < len(self.PersonMgr); i++ {

		self.PersonMgr[i].Ready = false
		self.PersonMgr[i].TuiBets = 0

		style, maxCard := GetSgyyCardType(self.PersonMgr[i].Card)
		total, num := GetSgyyCardValue(self.PersonMgr[i].Card)
		if 1 == gameModol { //！模式1
			if 0 == speStyle { //！都没选
				if style >= 4 || 2 == style {
					self.PersonMgr[i].CardMul = 4
				} else if 3 == style {
					style = 1
				}
			} else if 1 == speStyle { //！同花顺
				if style >= 4 || 2 == style {
					self.PersonMgr[i].CardMul = 4
				} else if 3 == style {
					style = 6
					self.PersonMgr[i].CardMul = 5
				}
			} else if 2 == speStyle { //！小三公
				if 5 == style || 2 == style {
					self.PersonMgr[i].CardMul = 4
				} else if 3 == style {
					style = 1
				} else if 4 == style {
					self.PersonMgr[i].CardMul = 6
					style = 6
				}
			} else if 3 == speStyle { //！同花顺+小三公
				if 3 == style {
					style = 6
					self.PersonMgr[i].CardMul = 5
				} else if 4 == style {
					style = 7
					self.PersonMgr[i].CardMul = 6
				} else if 5 == style || 2 == style {
					self.PersonMgr[i].CardMul = 4
				}
			} else if 4 == speStyle { //！大三公
				if 5 == style {
					self.PersonMgr[i].CardMul = 7
				} else if 4 == style || 2 == style {
					self.PersonMgr[i].CardMul = 4
				} else if 3 == style {
					style = 1
				}
			} else if 5 == speStyle { //！同花顺+大三公
				if 5 == style {
					self.PersonMgr[i].CardMul = 7
					style = 7
				} else if 3 == style {
					self.PersonMgr[i].CardMul = 6
					style = 6
				} else if 2 == style || 4 == style {
					self.PersonMgr[i].CardMul = 4
				}
			} else if 6 == speStyle { //！小三公+大三公
				if 5 == style {
					self.PersonMgr[i].CardMul = 7
				} else if 4 == style {
					self.PersonMgr[i].CardMul = 6
				} else if 3 == style {
					style = 1
				} else if 2 == style {
					self.PersonMgr[i].CardMul = 4
				}
			} else if 7 == speStyle { //！全选
				if 5 == style {
					self.PersonMgr[i].CardMul = 7
				} else if 4 == style {
					self.PersonMgr[i].CardMul = 6
				} else if 3 == style {
					self.PersonMgr[i].CardMul = 5
				} else if 2 == style {
					self.PersonMgr[i].CardMul = 4
				}
			}
			if 1 == style {
				if 1 == self.room.Param1/10000%10 {
					if total <= 6 && total >= 0 {
						self.PersonMgr[i].CardMul = 1
					} else {
						self.PersonMgr[i].CardMul = 2
					}
				} else {
					if total <= 6 && total >= 0 {
						self.PersonMgr[i].CardMul = 1
					} else if total <= 8 && total >= 7 {
						self.PersonMgr[i].CardMul = 2
					} else {
						self.PersonMgr[i].CardMul = 3
					}
				}
			}
		} else {
			self.PersonMgr[i].CardMul = 1
			self.PersonMgr[i].DealMul = 1
			if 3 == style || 4 == style {
				style = 1
				if IsSgyyOtherThree(self.PersonMgr[i].Card) { //混三公
					style = 2
				}

			}
		}

		if style != 1 {
			self.PersonMgr[i].CardValue = style*100000 + maxCard
		} else {
			self.PersonMgr[i].CardValue = style*100000 + total*10000 + num*1000 + maxCard
		}
	}

	if 4 != gameModol { //！有庄家的模式算分
		dealpos := -1
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Deal {
				dealpos = i
				break
			}
		}

		for i := 0; i < len(self.PersonMgr); i++ {
			if dealpos == i {
				continue
			}
			if self.PersonMgr[i].CardValue > self.PersonMgr[dealpos].CardValue { //！闲家赢
				winscore := self.PersonMgr[dealpos].DealMul * self.PersonMgr[i].Bets * self.PersonMgr[i].CardMul
				self.PersonMgr[i].Score += winscore
				self.PersonMgr[i].Total += winscore

				self.PersonMgr[dealpos].Score -= winscore
				self.PersonMgr[dealpos].Total -= winscore
				if 1 == gameModol && self.room.Param1%10 >= 4 {
					self.PersonMgr[i].TuiBets = winscore + self.PersonMgr[i].Bets
					if 1 == self.room.Param1/100%10 && self.PersonMgr[i].TuiBets >= 10 {
						self.PersonMgr[i].TuiBets = 10
					} else if 2 == self.room.Param1/100%10 && self.PersonMgr[i].TuiBets >= 20 {
						self.PersonMgr[i].TuiBets = 20
					} else if 3 == self.room.Param1/100%10 && self.PersonMgr[i].TuiBets >= 40 {
						self.PersonMgr[i].TuiBets = 40
					}
					//lib.GetLogMgr().Output(lib.LOG_DEBUG, "推注分数", self.PersonMgr[i].TuiBets)
					if 1 == self.room.Param1/100%10 && self.PersonMgr[i].TuiBets == 2 {
						self.PersonMgr[i].TuiBets = 0
					} else if 2 == self.room.Param1/100%10 && self.PersonMgr[i].TuiBets == 4 {
						self.PersonMgr[i].TuiBets = 0
					} else if 3 == self.room.Param1/100%10 && self.PersonMgr[i].TuiBets == 8 {
						self.PersonMgr[i].TuiBets = 0
					}
				}

			} else { //！庄家赢
				winscore := self.PersonMgr[dealpos].DealMul * self.PersonMgr[i].Bets * self.PersonMgr[dealpos].CardMul
				self.PersonMgr[i].Score -= winscore
				self.PersonMgr[i].Total -= winscore

				self.PersonMgr[dealpos].Score += winscore
				self.PersonMgr[dealpos].Total += winscore
			}
		}
		if 3 == gameModol && (self.PersonMgr[dealpos].Total <= 0) { //！下庄
			self.GameReDeal(self.PersonMgr[dealpos].Uid)
		}
	} else { //！无庄家模式算分
		valueArr := make([]int, 0)
		indexArr := make([]int, 0)
		for i := 0; i < len(self.PersonMgr); i++ {
			valueArr = append(valueArr, self.PersonMgr[i].CardValue)
			indexArr = append(indexArr, i)
		}

		for i := 0; i < len(valueArr); i++ { //！对下标进行排序
			for j := i + 1; j < len(valueArr); j++ {
				if valueArr[i] < valueArr[j] {
					temp := valueArr[i]
					valueArr[i] = valueArr[j]
					valueArr[j] = temp

					tempIndex := indexArr[i]
					indexArr[i] = indexArr[j]
					indexArr[j] = tempIndex
				}
			}
		}
		head := 0
		tail := len(self.PersonMgr) - 1
		for {
			if self.PersonMgr[indexArr[head]].Bets > self.PersonMgr[indexArr[tail]].Bets { //!最后一名分不够扣
				winscore := self.PersonMgr[indexArr[tail]].Bets
				//！赢家加分
				self.PersonMgr[indexArr[head]].Score += winscore
				self.PersonMgr[indexArr[head]].Total += winscore
				//！输家减分
				self.PersonMgr[indexArr[tail]].Score -= winscore
				self.PersonMgr[indexArr[tail]].Total -= winscore

				self.PersonMgr[indexArr[head]].Bets -= winscore
				self.PersonMgr[indexArr[tail]].Bets = 0
				tail--
			} else if self.PersonMgr[indexArr[head]].Bets < self.PersonMgr[indexArr[tail]].Bets { //！最后一名分口不完
				winscore := self.PersonMgr[indexArr[head]].Bets
				//！赢家加分
				self.PersonMgr[indexArr[head]].Score += winscore
				self.PersonMgr[indexArr[head]].Total += winscore
				//！输家减分
				self.PersonMgr[indexArr[tail]].Score -= winscore
				self.PersonMgr[indexArr[tail]].Total -= winscore

				self.PersonMgr[indexArr[tail]].Bets -= winscore
				self.PersonMgr[indexArr[head]].Bets = 0
				head++
			} else { //！最后一名分数刚好够扣完
				winscore := self.PersonMgr[indexArr[tail]].Bets
				self.PersonMgr[indexArr[head]].Score += winscore
				self.PersonMgr[indexArr[head]].Total += winscore

				self.PersonMgr[indexArr[tail]].Score -= winscore
				self.PersonMgr[indexArr[tail]].Total -= winscore

				self.PersonMgr[indexArr[tail]].Bets = 0
				self.PersonMgr[indexArr[head]].Bets = 0
				tail--
				head++
			}

			if head >= tail {
				break
			}
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].JoinStep > 0 {
			self.PersonMgr[i].JoinStep = 0
			for j := 0; j < len(self.RecordMgr); j++ {
				GetServer().InsertRecord(self.room.Type, self.PersonMgr[i].Uid, lib.HF_JtoA(&self.RecordMgr[j]), 0)
			}
		}
	}

	//! 记录
	var record Rec_GameNiuNiuJX
	record.Time = time.Now().Unix()
	record.Roomid = self.room.Id*100 + self.room.Step
	record.MaxStep = self.room.MaxStep
	record.Param1 = self.room.Param1
	record.Param2 = self.room.Param2

	var msg Msg_GameSGYY_End
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Ready = false
		var son Son_GameSGYY_End
		son.Uid = self.PersonMgr[i].Uid
		son.Score = self.PersonMgr[i].Score
		son.Total = self.PersonMgr[i].Total
		msg.Info = append(msg.Info, son)

		var rec Son_Rec_GameNiuNiuJX
		rec.Uid = self.PersonMgr[i].Uid
		rec.Name = self.room.GetName(self.PersonMgr[i].Uid)
		rec.Head = self.room.GetHead(self.PersonMgr[i].Uid)
		rec.Card = self.PersonMgr[i].Card
		rec.Dealer = self.PersonMgr[i].Deal
		rec.Bets = self.PersonMgr[i].ModBets
		rec.Score = self.PersonMgr[i].Score
		rec.RobDeal = self.PersonMgr[i].DealMul
		record.Info = append(record.Info, rec)
	}
	self.room.AddRecord(lib.HF_JtoA(&record))
	self.RecordMgr = append(self.RecordMgr, record)
	self.room.broadCastMsg("gamesgyyend", &msg)
	if self.room.IsBye() {
		self.OnBye()
		self.room.Bye()
		return
	}

	self.TickTime = 0
	self.TimeState = 4 //准备时间
	self.room.flush()
}
func (self *Game_SGYY) OnExit(uid int64) {
	//	find := false
	for i := 0; i < len(self.Ready); i++ {
		if self.Ready[i] == uid {
			copy(self.Ready[i:], self.Ready[i+1:])
			self.Ready = self.Ready[:len(self.Ready)-1]
			//	find = true
			break
		}
	}
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			copy(self.PersonMgr[i:], self.PersonMgr[i+1:])
			self.PersonMgr = self.PersonMgr[:len(self.PersonMgr)-1]
			break
		}
	}
	//	if !find {
	//		if len(self.Ready) == len(self.room.Uid) && len(self.Ready) >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 准备的人数达到游戏最小人数
	//			lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
	//			self.OnBegin()
	//			return
	//		}
	//	}
}
func (self *Game_SGYY) OnSendInfo(person *Person) {
	//person.SendMsg("gamesgyyinfo", self.getInfo(person.Uid))
	//	if len(self.PersonMgr) <= 0 {
	//		self.GameSeat(person.Uid)
	//	}
	//! 观众模式游戏,观众进来只发送游戏信息
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			person.SendMsg("gamesgyyinfo", self.getInfo(person.Uid))
			return
		}
	}
	person.SendMsg("gamesgyyinfo", self.getInfo(0))
}
func (self *Game_SGYY) getInfo(uid int64) *Msg_GameSGYY_Info {
	var msg Msg_GameSGYY_Info
	msg.Begin = self.room.Begin
	msg.State = self.TimeState
	msg.Host = self.room.Host
	msg.Ready = make([]int64, 0)
	msg.Info = make([]Son_GameSGYY_Info, 0)
	if !msg.Begin {
		msg.Ready = self.Ready
	}
	for _, value := range self.PersonMgr {
		var son Son_GameSGYY_Info
		son.Uid = value.Uid
		son.Dealer = value.Deal
		son.Score = value.Score
		son.Bets = value.Bets
		son.Total = value.Total
		son.Ready = value.Ready
		son.TuiBets = value.TuiBets
		son.RobDeal = value.DealMul
		son.View = value.View
		if msg.Begin {
			if son.Uid == uid || value.View { //！自己牌或者亮牌的可见
				son.Card = value.Card
			} else { //！其他人牌不可见
				for i := 0; i < len(value.Card); i++ {
					son.Card = append(son.Card, 0)
				}
			}
		} else { //！游戏结束所有人牌可见
			son.Card = value.Card
		}
		msg.Info = append(msg.Info, son)
	}
	return &msg
}
func (self *Game_SGYY) GameGetDealer() {
	max := self.PersonMgr[0].DealMul
	maxDeal := make([]int, 0)
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].DealMul > max {
			max = self.PersonMgr[i].DealMul
			maxDeal = make([]int, 0)
			maxDeal = append(maxDeal, i)
		} else if self.PersonMgr[i].DealMul == max {
			maxDeal = append(maxDeal, i)
		}
	}
	if 0 == max {
		for i := 0; i < len(self.PersonMgr); i++ {
			self.PersonMgr[i].DealMul = 1
		}
	}
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := (random.Intn(10000)) % len(maxDeal)
	self.PersonMgr[maxDeal[index]].Deal = true

	var msg staticfunc.Msg_Uid
	msg.Uid = self.PersonMgr[maxDeal[index]].Uid
	self.room.broadCastMsg("gamedealer", &msg)

	self.TimeState = 2 //下注时间开始
	self.TickTime = 0
}
func (self *Game_SGYY) TimeEndDeal() {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].DealMul < 0 {
			self.PersonMgr[i].DealMul = 0
			var msg Msg_GameDealer
			msg.Uid = self.PersonMgr[i].Uid
			msg.Score = self.PersonMgr[i].DealMul
			self.room.broadCastMsg("gamedealmul", &msg)
		}
	}
	self.GameGetDealer()
}
func (self *Game_SGYY) GameSendCard() {
	for i := 0; i < len(self.PersonMgr); i++ {
		if 1 == (self.room.Param1/10%10) && len(self.PersonMgr[i].Card) < 3 { //！模式1追加一张牌
			self.PersonMgr[i].Card = append(self.PersonMgr[i].Card, self.Cardmgr.Deal(1)...)
		} else if 1 != (self.room.Param1/10%10) && len(self.PersonMgr[i].Card) < 2 {
			self.PersonMgr[i].Card = self.Cardmgr.Deal(3)
		}
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		var msg Msg_GameView
		msg.Uid = self.PersonMgr[i].Uid
		msg.Card = self.PersonMgr[i].Card
		if nil != person {
			person.SendMsg("gamesgyycard", &msg)
		}
	}

	self.TimeState = 3
	self.TickTime = 0
	self.room.flush()
}
func (self *Game_SGYY) TimeEndBets() {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Bets <= 0 {
			if 1 == (self.room.Param1 / 10 % 10) {
				if 1 == (self.room.Param1 / 100 % 10) {
					self.GameBets(self.PersonMgr[i].Uid, 1)
				} else if 2 == (self.room.Param1 / 100 % 10) {
					self.GameBets(self.PersonMgr[i].Uid, 2)
				} else if 3 == (self.room.Param1 / 100 % 10) {
					self.GameBets(self.PersonMgr[i].Uid, 4)
				}
			} else {
				self.GameBets(self.PersonMgr[i].Uid, 20)
			}
		}
	}
	self.GameSendCard()
}
func (self *Game_SGYY) TimeEndView() {
	for i := 0; i < len(self.PersonMgr); i++ {
		if !self.PersonMgr[i].View {
			var msg Msg_GameView
			msg.Uid = self.PersonMgr[i].Uid
			msg.Card = self.PersonMgr[i].Card
			self.room.broadCastMsg("gameview", &msg)
			self.PersonMgr[i].View = true
		}
	}
	self.OnEnd()
	self.room.flush()
}
func (self *Game_SGYY) TimeEndReady() {
label:
	for i := 1; i < len(self.PersonMgr); i++ {
		for j := 0; j < len(self.Ready); j++ {
			if self.Ready[j] == self.PersonMgr[i].Uid {
				continue label
			}
		}
		self.Ready = append(self.Ready, self.PersonMgr[i].Uid)
	}
	self.OnBegin()
}
func (self *Game_SGYY) OnTime() {
	if 0 == self.TimeState {
		return
	}
	self.TickTime += 0.1                            //! 计时器状态 1 抢庄 2下注 3明牌
	if self.TickTime >= 12 && 1 == self.TimeState { //！抢庄时间到
		if 3 == (self.room.Param1/10%10) || 4 == (self.room.Param1/10%10) {
			return
		}
		self.TimeState = 0
		self.TimeEndDeal()
	}
	if self.TickTime >= 12 && 2 == self.TimeState { //！下注时间到
		self.TimeState = 0
		self.TimeEndBets()
	}
	if self.TickTime >= 15 && 3 == self.TimeState { //！明牌时间到
		self.TimeState = 0
		self.TimeEndView()
	}
	if self.TickTime >= 18 && 4 == self.TimeState { //！准备时间到
		self.TimeState = 0
		self.TimeEndReady()
	}
}
func (self *Game_SGYY) GetPerson(uid int64) *Game_SGYY_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}
	return nil
}
func (self *Game_SGYY) OnBalance() {

}

func (self *Game_SGYY) OnIsBets(uid int64) bool {
	return false
}

func (self *Game_SGYY) OnRobot(robot *lib.Robot) {

}

func (self *Game_SGYY) OnIsDealer(uid int64) bool {
	return false
}
