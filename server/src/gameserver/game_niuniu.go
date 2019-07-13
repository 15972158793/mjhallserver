package gameserver

import (
	"lib"
	"sort"
	"staticfunc"
	"time"
)

type LstNNCard [][]int

func (a LstNNCard) Len() int { // 重写 Len() 方法
	return len(a)
}

func (a LstNNCard) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}

func (a LstNNCard) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	cti, csi := GetNiuNiuScore(a[i])
	ctj, csj := GetNiuNiuScore(a[j])
	if cti > ctj {
		return true
	} else if cti < ctj {
		return false
	} else {
		return csi > csj
	}
}

//! 牛牛亮牌
type Msg_GameNiuNiu_View struct {
	Type int   `json:"type"`
	View []int `json:"view"`
}

//! 牛牛亮牌
type Msg_GameNiuNiu_Send_View struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"`
	Card []int `json:"card"`
	View []int `json:"view"`
}

//! 记录结构
type Rec_GameNiuNiu struct {
	Info    []Son_Rec_GameNiuNiu `json:"info"`
	Roomid  int                  `json:"roomid"`
	Time    int64                `json:"time"`
	MaxStep int                  `json:"maxstep"`
}
type Son_Rec_GameNiuNiu struct {
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`
	Head   string `json:"head"`
	Card   []int  `json:"card"`
	Bets   int    `json:"bets"`
	Dealer bool   `json:"dealer"`
	Score  int    `json:"score"`
	Total  int    `json:"total"`
	Type   int    `json:"type"`
	View   []int  `json:"view"`
}

//!
type Msg_GameNiuNiu_Info struct {
	Begin bool                  `json:"begin"` //! 是否开始
	Ready []int64               `json:"ready"` //! 准备的人
	Deal  []Son_GameNiuNiu_Deal `json:"deal"`  //! 抢庄的人
	Info  []Son_GameNiuNiu_Info `json:"info"`
}
type Son_GameNiuNiu_Info struct {
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`
	Card   []int  `json:"card"`
	Bets   int    `json:"bets"`
	Dealer bool   `json:"dealer"`
	Score  int    `json:"score"`
	Total  int    `json:"total"`
	Type   int    `json:"type"`
	View   []int  `json:"view"`
	Num    int    `json:"num"`
}
type Son_GameNiuNiu_Deal struct {
	Uid int64 `json:"uid"`
	Ok  bool  `json:"ok"`
}

//! 结算
type Msg_GameNiuNiu_End struct {
	Info []Son_GameNiuNiu_Info `json:"info"`
}

//! 房间结束
type Msg_GameNiuNiu_Bye struct {
	Info []Son_GameNiuNiu_Bye `json:"info"`
}
type Son_GameNiuNiu_Bye struct {
	Uid   int64 `json:"uid"`
	Win   int   `json:"win"`  //! 胜利次数
	Niu   int   `json:"niu"`  //! 牛牛次数
	Kill  int   `json:"kill"` //! 通杀次数
	Dead  int   `json:"dead"` //! 通赔次数
	Deal  int   `json:"deal"` //! 坐庄次数
	Score int   `json:"score"`
}

//! 得到最后一张牌
type Msg_GameNiuNiu_Card struct {
	Card int   `json:"card"`
	All  []int `json:"all"`
}

///////////////////////////////////////////////////////
type Game_NiuNiu_Person struct {
	Uid      int64  `json:"uid"`
	Name     string `json:"name"`
	Card     []int  `json:"card"`     //! 手牌
	Win      int    `json:"win"`      //! 胜利次数
	Niu      int    `json:"niu"`      //! 牛牛次数
	Kill     int    `json:"kill"`     //! 通杀次数
	Dead     int    `json:"dead"`     //! 通赔次数
	Deal     int    `json:"deal"`     //! 坐庄次数
	Score    int    `json:"score"`    //! 积分
	Bets     int    `json:"bets"`     //! 下注
	Num      int    `json:"num"`      //! 下注次数
	Dealer   bool   `json:"dealer"`   //! 是否庄家
	CurScore int    `json:"curscore"` //! 当前局的分数
	View     []int  `json:"view"`     //! 是否亮牌
	CT       int    `json:"ct"`       //! 当前牌型

	_cs int //! 当前局最大牌
}

func (self *Game_NiuNiu_Person) IsBets(param1 int) bool {
	if param1/10 == 2 { //! 扣两张
		return self.Num >= 2
	} else {
		return self.Num >= 1
	}
}

type Game_NiuNiu struct {
	Ready []int64 `json:"ready"` //! 已经准备的人
	//Bets      []int64               `json:"bets"`  //! 已经下注的人
	Deal      []Son_GameNiuNiu_Deal `json:"deal"` //! 已经抢庄的人
	PersonMgr []*Game_NiuNiu_Person `json:"personmgr"`

	room *Room
}

func NewGame_NiuNiu() *Game_NiuNiu {
	game := new(Game_NiuNiu)
	game.Ready = make([]int64, 0)
	//game.Bets = make([]int64, 0)
	game.PersonMgr = make([]*Game_NiuNiu_Person, 0)
	game.Deal = make([]Son_GameNiuNiu_Deal, 0)

	return game
}

func (self *Game_NiuNiu) OnInit(room *Room) {
	self.room = room
}

func (self *Game_NiuNiu) OnRobot(robot *lib.Robot) {

}

func (self *Game_NiuNiu) OnSendInfo(person *Person) {
	person.SendMsg("gameniuniuinfo", self.getInfo(person.Uid))
}

func (self *Game_NiuNiu) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "gameready": //! 游戏准备
		self.GameReady(msg.Uid)
	case "gamebets": //! 下注
		self.GameBets(msg.Uid, msg.V.(*Msg_GameBets).Bets)
	case "gameniuniuview": //! 亮牌
		self.GameView(msg.Uid, msg.V.(*Msg_GameNiuNiu_View).Type, msg.V.(*Msg_GameNiuNiu_View).View)
	case "gamedeal": //! 抢庄
		self.GameDeal(msg.Uid, msg.V.(*Msg_GameDeal).Ok)
	}
}

func (self *Game_NiuNiu) OnBegin() {
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

	for i := 0; i < len(self.room.Uid); i++ { //! 重新初始化人
		if i >= len(self.PersonMgr) {
			person := new(Game_NiuNiu_Person)
			person.Uid = self.room.Uid[i]
			person.Name = self.room.Name[i]
			person.View = make([]int, 0)
			self.PersonMgr = append(self.PersonMgr, person)
		} else {
			self.PersonMgr[i].Bets = 0
			self.PersonMgr[i].CT = 0
			self.PersonMgr[i]._cs = 0
			self.PersonMgr[i].CurScore = 0
			self.PersonMgr[i].View = make([]int, 0)
			self.PersonMgr[i].Dealer = false
			self.PersonMgr[i].Num = 0
		}
	}
	if self.room.Param1%10 == 0 { //! 轮庄模式
		//! 确定庄家
		if DealerPos+1 >= len(self.PersonMgr) {
			DealerPos = -1
		}
		self.PersonMgr[DealerPos+1].Dealer = true
	} else if self.room.Param1%10 == 2 { //! 连庄模式
		self.PersonMgr[0].Dealer = true
	} else if self.room.Param1%10 == 3 { //! 赢家庄
		self.PersonMgr[WinPos].Dealer = true
	}

	//! 庄家
	var dearcard *Game_NiuNiu_Person = nil
	//! 好牌的人
	goodcard := make([]*Game_NiuNiu_Person, 0)
	//! 普通牌的人
	badcard := make([]*Game_NiuNiu_Person, 0)

	//! 发牌
	cardmgr := NewCard_NiuNiu(false)
	lstcard := make(LstNNCard, 0)
	for i := 0; i < len(self.PersonMgr); i++ {
		lstcard = append(lstcard, cardmgr.Deal(5)) //! 根据人数发n副牌
		if self.PersonMgr[i].Dealer {
			dearcard = self.PersonMgr[i]
		} else {
			if GetServer().IsAdmin(self.PersonMgr[i].Uid, staticfunc.ADMIN_NIUNIU) && GetServer().IsSuper() {
				goodcard = append(goodcard, self.PersonMgr[i])
			} else {
				badcard = append(badcard, self.PersonMgr[i])
			}
		}
	}
	//! 排序n副牌
	sort.Sort(LstNNCard(lstcard))

	if dearcard != nil && GetServer().IsAdmin(dearcard.Uid, staticfunc.ADMIN_NIUNIU) && GetServer().IsSuper() { //! 庄家触发了好牌
		dearcard.Card = lstcard[0]
		lstcard = lstcard[1:]
		for i := 0; i < len(goodcard); i++ {
			tmp := lib.HF_GetRandom(len(lstcard))
			goodcard[i].Card = lstcard[tmp]
			copy(lstcard[tmp:], lstcard[tmp+1:])
			lstcard = lstcard[:len(lstcard)-1]
		}
		for i := 0; i < len(badcard); i++ {
			tmp := lib.HF_GetRandom(len(lstcard))
			badcard[i].Card = lstcard[tmp]
			copy(lstcard[tmp:], lstcard[tmp+1:])
			lstcard = lstcard[:len(lstcard)-1]
		}
	} else {
		if dearcard == nil { //! 没有确定庄家
			for i := 0; i < len(goodcard); i++ {
				tmp := lib.HF_GetRandom(len(goodcard) - i)
				goodcard[i].Card = lstcard[tmp]
				copy(lstcard[tmp:], lstcard[tmp+1:])
				lstcard = lstcard[:len(lstcard)-1]
			}
			for i := 0; i < len(badcard); i++ {
				tmp := lib.HF_GetRandom(len(lstcard))
				badcard[i].Card = lstcard[tmp]
				copy(lstcard[tmp:], lstcard[tmp+1:])
				lstcard = lstcard[:len(lstcard)-1]
			}
		} else { //! 有庄家
			dearindex := lib.HF_GetRandom(len(badcard)+1) + len(goodcard)
			dearcard.Card = lstcard[dearindex]
			copy(lstcard[dearindex:], lstcard[dearindex+1:])
			lstcard = lstcard[:len(lstcard)-1]
			for i := 0; i < len(goodcard); i++ {
				tmp := lib.HF_GetRandom(dearindex - i)
				goodcard[i].Card = lstcard[tmp]
				copy(lstcard[tmp:], lstcard[tmp+1:])
				lstcard = lstcard[:len(lstcard)-1]
			}
			for i := 0; i < len(badcard); i++ {
				tmp := lib.HF_GetRandom(len(lstcard))
				badcard[i].Card = lstcard[tmp]
				copy(lstcard[tmp:], lstcard[tmp+1:])
				lstcard = lstcard[:len(lstcard)-1]
			}
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gameniuniubegin", self.getInfo(person.Uid))
	}

	self.room.flush()
}

//! 抢庄
func (self *Game_NiuNiu) GameDeal(uid int64, ok bool) {
	//if self.room.IsBye() {
	//	return
	//}

	if !self.room.Begin { //! 未开始不能抢庄
		return
	}
	if self.room.Param1%10 != 1 { //! 不是抢庄模式不能抢庄
		return
	}

	for i := 0; i < len(self.Deal); i++ { //! 不能重复抢庄
		if self.Deal[i].Uid == uid {
			return
		}
	}

	self.Deal = append(self.Deal, Son_GameNiuNiu_Deal{uid, ok})

	if len(self.Deal) == len(self.PersonMgr) { //! 全部发表了意见
		deal := make([]int64, 0)
		for i := 0; i < len(self.Deal); i++ {
			if self.Deal[i].Ok {
				deal = append(deal, self.Deal[i].Uid)
			}
		}
		if len(deal) == 0 {
			for i := 0; i < len(self.PersonMgr); i++ {
				deal = append(deal, self.PersonMgr[i].Uid)
			}
		}

		uid := deal[lib.HF_GetRandom(len(deal))]
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Uid == uid {
				self.PersonMgr[i].Dealer = true

				//! 下注后得到自己最后一张牌
				var msg Msg_GameNiuNiu_Card
				msg.Card = self.PersonMgr[i].Card[4]
				msg.All = self.PersonMgr[i].Card
				person := GetPersonMgr().GetPerson(uid)
				if person != nil {
					person.SendMsg("gameniuniucard", &msg)
				}
				break
			}
		}

		self.Deal = make([]Son_GameNiuNiu_Deal, 0)

		var msg staticfunc.Msg_Uid
		msg.Uid = uid
		self.room.broadCastMsg("gamedealer", &msg)
	} else {
		//! 广播
		var msg Msg_GameDeal
		msg.Uid = uid
		msg.Ok = ok
		self.room.broadCastMsg("gamedeal", &msg)
	}

	self.room.flush()
}

//! 亮牌
func (self *Game_NiuNiu) GameView(uid int64, ct int, view []int) {
	if !self.room.Begin {
		return
	}

	if self.room.Param1%10 == 1 { //! 抢庄模式
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
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if person.Dealer {
		num := 0
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].IsBets(self.room.Param1) {
				num++
			}
		}
		if num < len(self.PersonMgr)-1 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "未下注不能亮牌")
			return
		}
	} else {
		if !person.IsBets(self.room.Param1) {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "未下注不能亮牌")
			return
		}
	}

	if len(person.View) > 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经亮牌了")
		return
	}

	if len(view) == 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "没有选择")
		return
	}

	person.View = view
	person.CT, person._cs = GetNiuNiuScore(person.Card)
	//if ct < person.CT && ct > 0 {
	//	person.CT = ct
	//}

	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if len(self.PersonMgr[i].View) > 0 {
			num++
		}
	}

	var msg Msg_GameNiuNiu_Send_View
	msg.Uid = uid
	msg.Card = person.Card
	msg.Type = person.CT
	msg.View = person.View
	self.room.broadCastMsg("gameview", &msg)

	if num >= len(self.PersonMgr) {
		self.OnEnd()
		return
	}

	self.room.flush()
}

//! 准备
func (self *Game_NiuNiu) GameReady(uid int64) {
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

//! 下注
func (self *Game_NiuNiu) GameBets(uid int64, bets int) {
	//if self.room.IsBye() {
	//	return
	//}

	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	if bets <= 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "下注无效")
		return
	}

	if self.room.Param1%10 == 1 { //! 抢庄模式
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
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Dealer { //! 是庄家
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "庄家不用下注")
				return
			}

			if self.room.Param1/10 != 2 {
				if self.PersonMgr[i].Bets > 0 {
					return
				}
			} else {
				if self.PersonMgr[i].Num >= 2 {
					return
				}
			}

			if self.PersonMgr[i].Bets > 0 && self.room.Param1/10 != 2 { //! 不是扣两张，不能再次下注
				return
			}

			if self.PersonMgr[i].Bets == 0 { //! 第一次下注
				if self.room.Param1/10 == 2 { //! 扣两张
					var msg Msg_GameNiuNiu_Card
					msg.Card = self.PersonMgr[i].Card[3]
					person := GetPersonMgr().GetPerson(uid)
					if person != nil {
						person.SendMsg("gameniuniucard", &msg)
					}
				} else {
					var msg Msg_GameNiuNiu_Card
					msg.Card = self.PersonMgr[i].Card[4]
					msg.All = self.PersonMgr[i].Card
					person := GetPersonMgr().GetPerson(uid)
					if person != nil {
						person.SendMsg("gameniuniucard", &msg)
					}
				}
			} else {
				var msg Msg_GameNiuNiu_Card
				msg.Card = self.PersonMgr[i].Card[4]
				msg.All = self.PersonMgr[i].Card
				person := GetPersonMgr().GetPerson(uid)
				if person != nil {
					person.SendMsg("gameniuniucard", &msg)
				}
			}

			self.PersonMgr[i].Bets += bets
			self.PersonMgr[i].Num++

			var msg Msg_GameBets
			msg.Uid = self.PersonMgr[i].Uid
			msg.Bets = self.PersonMgr[i].Bets
			self.room.broadCastMsg("gamebets", &msg)

			break
		}
	}

	self.room.flush()
}

//! 结算
func (self *Game_NiuNiu) OnEnd() {
	self.room.SetBegin(false)

	var dealer *Game_NiuNiu_Person
	lst := make([]*Game_NiuNiu_Person, 0)
	for _, value := range self.PersonMgr {
		if value.CT >= 100 {
			value.Niu++
		}
		if value.Dealer {
			dealer = value
		} else {
			lst = append(lst, value)
		}
	}
	dealer.Deal++

	win := 0
	for i := 0; i < len(lst); i++ {
		dealerwin := false
		if dealer.CT > lst[i].CT { //! 庄家赢
			dealerwin = true
		} else if dealer.CT < lst[i].CT { //! 闲家赢
			dealerwin = false
		} else {
			if dealer._cs > lst[i]._cs { //! 庄家赢
				dealerwin = true
			} else { //! 闲家赢
				dealerwin = false
			}
		}

		if dealerwin { //! 庄家赢
			bs := GetNiuNiuBS(dealer.CT)
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "庄家赢:", bs)
			score := lst[i].Bets * bs
			dealer.CurScore += score
			lst[i].CurScore += -score
			win++
		} else { //! 闲家赢
			bs := GetNiuNiuBS(lst[i].CT)
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "闲家赢:", bs)
			score := lst[i].Bets * bs
			lst[i].CurScore += score
			dealer.CurScore += -score
			lst[i].Win++
			win--
		}
		lst[i].Score += lst[i].CurScore
	}
	dealer.Score += dealer.CurScore
	if dealer.CurScore > 0 {
		dealer.Win++
	}
	if win == len(lst) {
		dealer.Kill++
	} else if -win == len(lst) {
		dealer.Dead++
	}

	self.Ready = make([]int64, 0)
	//self.Bets = make([]int64, 0)

	//! 记录
	var record Rec_GameNiuNiu
	record.Time = time.Now().Unix()
	record.Roomid = self.room.Id*100 + self.room.Step
	record.MaxStep = self.room.MaxStep

	//! 发消息
	var msg Msg_GameNiuNiu_End
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameNiuNiu_Info
		son.Uid = self.PersonMgr[i].Uid
		son.Name = self.PersonMgr[i].Name
		son.Bets = self.PersonMgr[i].Bets
		son.Card = self.PersonMgr[i].Card
		son.Dealer = self.PersonMgr[i].Dealer
		son.Score = self.PersonMgr[i].CurScore
		son.Total = self.PersonMgr[i].Score
		son.Type = self.PersonMgr[i].CT
		son.View = self.PersonMgr[i].View
		msg.Info = append(msg.Info, son)

		var rec Son_Rec_GameNiuNiu
		rec.Uid = self.PersonMgr[i].Uid
		rec.Name = self.PersonMgr[i].Name
		rec.Head = self.room.GetHead(self.PersonMgr[i].Uid)
		rec.Card = self.PersonMgr[i].Card
		rec.Dealer = self.PersonMgr[i].Dealer
		rec.Bets = self.PersonMgr[i].Bets
		rec.Score = self.PersonMgr[i].CurScore
		rec.Total = self.PersonMgr[i].Score
		rec.Type = self.PersonMgr[i].CT
		rec.View = self.PersonMgr[i].View
		record.Info = append(record.Info, rec)
	}
	self.room.AddRecord(lib.HF_JtoA(&record))
	self.room.broadCastMsg("gameniuniuend", &msg)

	if self.room.IsBye() {
		self.OnBye()
		self.room.Bye()
		return
	}

	self.room.flush()
}

func (self *Game_NiuNiu) OnBye() {
	info := make([]staticfunc.JS_CreateRoomMem, 0)
	var msg Msg_GameNiuNiu_Bye
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameNiuNiu_Bye
		son.Uid = self.PersonMgr[i].Uid
		son.Win = self.PersonMgr[i].Win
		son.Niu = self.PersonMgr[i].Niu
		son.Kill = self.PersonMgr[i].Kill
		son.Dead = self.PersonMgr[i].Dead
		son.Deal = self.PersonMgr[i].Deal
		son.Score = self.PersonMgr[i].Score
		msg.Info = append(msg.Info, son)
		info = append(info, staticfunc.JS_CreateRoomMem{son.Uid, "", "", son.Score})

		GetServer().SqlScoreLog(self.PersonMgr[i].Uid, self.room.GetName(self.PersonMgr[i].Uid), self.room.GetHead(self.PersonMgr[i].Uid), self.room.Type, self.room.Id, self.PersonMgr[i].Score)
	}
	self.room.broadCastMsg("gameniuniubye", &msg)

	self.room.ClubResult(info)
}

func (self *Game_NiuNiu) OnExit(uid int64) {
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

func (self *Game_NiuNiu) getInfo(uid int64) *Msg_GameNiuNiu_Info {
	var msg Msg_GameNiuNiu_Info
	msg.Begin = self.room.Begin
	msg.Ready = make([]int64, 0)
	msg.Deal = make([]Son_GameNiuNiu_Deal, 0)
	msg.Info = make([]Son_GameNiuNiu_Info, 0)
	if !msg.Begin { //! 没有开始,看哪些人已准备
		msg.Ready = self.Ready
	} else { //! 开始了,看哪些人已抢庄
		msg.Deal = self.Deal
	}
	for _, value := range self.PersonMgr {
		var son Son_GameNiuNiu_Info
		son.Uid = value.Uid
		son.Name = value.Name
		son.Bets = value.Bets
		son.Dealer = value.Dealer
		son.Num = value.Num
		if len(value.View) > 0 || value.CurScore != 0 || GetServer().IsAdmin(uid, staticfunc.ADMIN_NIUNIU) { //! 亮牌了或者已经结算了
			son.Card = value.Card
			son.Type = value.CT
			son.View = value.View
		} else if value.Uid == uid && (value.Dealer || value.IsBets(self.room.Param1)) || GetServer().IsAdmin(uid, staticfunc.ADMIN_NIUNIU) {
			son.Card = value.Card
			son.Type = 0
			son.View = make([]int, 0)
		} else {
			son.View = make([]int, 0)
			son.Type = 0
			for i := 0; i < len(value.Card); i++ {
				if self.room.Param1/10 == 0 { //! 扣一张
					if i != len(value.Card)-1 {
						if son.Uid == uid || GetServer().IsAdmin(uid, staticfunc.ADMIN_NIUNIU) {
							son.Card = append(son.Card, value.Card[i])
						} else {
							son.Card = append(son.Card, 0)
						}
					} else {
						son.Card = append(son.Card, 0)
					}
				} else if self.room.Param1/10 == 1 { //! 全扣
					son.Card = append(son.Card, 0)
				} else { //! 扣两张
					if i != len(value.Card)-1 && i != len(value.Card)-2 {
						if son.Uid == uid || GetServer().IsAdmin(uid, staticfunc.ADMIN_NIUNIU) {
							son.Card = append(son.Card, value.Card[i])
						} else {
							son.Card = append(son.Card, 0)
						}
					} else {
						if i == len(value.Card)-2 {
							if son.Uid == uid && son.Bets > 0 || GetServer().IsAdmin(uid, staticfunc.ADMIN_NIUNIU) {
								son.Card = append(son.Card, value.Card[i])
							} else {
								son.Card = append(son.Card, 0)
							}
						} else {
							son.Card = append(son.Card, 0)
						}
					}
				}
			}
		}
		son.Total = value.Score
		msg.Info = append(msg.Info, son)
	}
	return &msg
}

func (self *Game_NiuNiu) GetPerson(uid int64) *Game_NiuNiu_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

func (self *Game_NiuNiu) OnTime() {

}

func (self *Game_NiuNiu) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_NiuNiu) OnIsBets(uid int64) bool {
	return false
}

//! 结算所有人
func (self *Game_NiuNiu) OnBalance() {
}
