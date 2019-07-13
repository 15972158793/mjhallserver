package gameserver

import (
	"lib"
	"math/rand"
	"time"

	"staticfunc"
)

//色子的点数
type Msg_GameT_SZ struct { //发牌起点
	Uid int64 `json:"uid"` //!开始发牌的玩家uid
	Num []int `json:"num"` //! 色子点数
}

//!告诉客户端发牌已完成
type Msg_GameT_Begin struct {
	Info        Son_GameT_Info `json:"info"`
	ChoiceScore int            `json:"choicescore"`
}

//! 单局结算
type Msg_GameT_End struct {
	Info    []Son_GameT_Info `json:"info"`
	WinAll  bool             `json:"winall"`  //!是否通杀
	LoseAll bool             `json:"loseall"` //!是否通赔
}

//! 记录结构
type Rec_GameT struct {
	Info    []Son_Rec_GameT `json:"info"`
	Roomid  int             `json:"roomid"`
	Time    int64           `json:"time"`
	MaxStep int             `json:"maxstep"`
}
type Son_Rec_GameT struct {
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`
	Head   string `json:"head"`
	Card   []int  `json:"card"`
	Bets   int    `json:"bets"`
	Dealer bool   `json:"dealer"`
	Score  int    `json:"score"`
	Total  int    `json:"total"`
}

//! 房间结束
type Msg_GameT_Bye struct {
	Info []Son_GameT_Bye `json:"info"`
}
type Son_GameT_Bye struct {
	Uid        int64 `json:"uid"`
	DealNum    int64 `json:"dealnum"`    //!坐庄次数
	WinAllNum  int64 `json:"winallnum"`  //!通吃次数
	LoseAllNum int64 `json:"loseallnum"` //!通赔次数
	WinNum     int64 `json:"winnum"`     //! 赢的次数
	LostNum    int64 `json:"lostnum"`    //! 赢的次数
	Total      int   `json:"total"`      //! 本局输赢的分数
	Score      int   `json:"score"`      //! 输赢分数总和
	NumBaoZi   int   `json:"numbaozi"`   //!豹子次数
}

//每个玩家的数据总览
type Game_T_Person struct {
	Uid      int64 `json:"uid"`
	Deal     bool  `json:"deal"` //!是否是庄
	NextDeal bool  `json:"nextdeal"`
	So_Card  int   `json:"so_card"` //! 倍率
	So_Df    int   `json:"so_df"`   //! 底分
	Total    int   `json:"total"`   //! 本局输赢总分
	Ready    bool  `json:"ready"`   //! 是否准备好
	Card     []int `json:"card"`    //! 手里的牌
	PaiXing  int   `json:"paixing"` //牌型
	//!new
	DealNum    int64 `json:"dealnum"`    //!坐庄次数
	WinAllNum  int64 `json:"winallnum"`  //!通吃次数
	LoseAllNum int64 `json:"loseallnum"` //!通赔次数
	WinNum     int64 `json:"winnum"`     //! 赢的次数
	LostNum    int64 `json:"lostnum"`    //! 赢的次数
	Score      int   `json:"score"`      //! 输赢分数总和
	Uid_KS     int64 `json:"uid_ks"`     //! 本局开始发牌的玩家UID
	NumBaoZi   int   `json:"numbaozi"`   //! 玩家拿到豹子的次数
}
type MahMgrT struct {
	Card []int //! 剩余牌组
	Temp []int //! 已经用过的牌
}
type Game_T struct {
	PersonMgr   []*Game_T_Person `json:"personmgr"`
	Ready       []int64          `json:"ready"`   //! 已经准备的人
	Mah         MahMgrT          `json:"mah"`     //! 剩余
	Winer       int              `json:"winer"`   //! 上局谁赢下标
	BefDeal     int              `json:"befdeal"` //! 上局庄家下标
	ChoiceScore int              `json:"choicescore"`
	DealUid     int64            `json:"dealuid"` //当局庄家的uid
	mahmgr      *MahMgr
	room        *Room
}
type Msg_GameT_Info struct {
	Begin       bool             `json:"begin"`    //! 是否开始
	Info        []Son_GameT_Info `json:"info"`     //! 人的info
	Num         int64            `json:"num"`      //! 起始发牌的玩家UID
	LastCard    []int            `json:"lastcard"` //! 奇数次用过的牌
	ChoiceScore int              `json:"choicescore"`
}

//每一局玩家的信息
type Son_GameT_Info struct {
	Uid         int64 `json:"uid"`
	Deal        bool  `json:"deal"`
	NextDeal    bool  `json:"nextdeal"`
	Card        []int `json:"card"`
	Score       int   `json:"score"` //! 当局点数
	So_Df       int   `json:"so_df"` //! 当局底分
	ChoiceScore int   `json:"choicescore"`
	Total       int
	Ready       bool `json:"ready"`
}

//! 选择底分倍数 0,1,2
type ChoiceScore struct {
	Score int `json:"score"`
}

func (self *Game_T) OnMsg(msg *RoomMsg) { //!
	switch msg.Head {
	case "gameready": //! 游戏准备 ok!
		self.GameReady(msg.Uid)
	case "ChoiceScore": //! 选择下注分数  ok!
		self.GameStep(msg.Uid, msg.V.(*ChoiceScore).Score)
	}
}

func (self *Game_T) GameStep(uid int64, Score int) {
	//lib.GetLogMgr().Output(lib.LOG_DEBUG, "")
	if !self.room.Begin { //! 没有开始不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}
	if self.DealUid == uid { //庄家不能下注
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "庄家不能下注")
		return
	}
	person := self.GetPerson(uid)
	if nil == person {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "玩家不存在")
		return
	}
	if person.So_Df > 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "重复下注")
		return
	}

	var msg Msg_GameBets
	msg.Uid = uid

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid != uid {
			continue
		}
		if 0 == Score {
			self.PersonMgr[i].So_Df = 3
			msg.Bets = 3
		} else if 1 == Score {
			self.PersonMgr[i].So_Df = 5
			msg.Bets = 5
		} else if 2 == Score {
			self.PersonMgr[i].So_Df = 7
			msg.Bets = 7
		}
		break
	}

	self.room.broadCastMsg("gameTbets", &msg)

	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].So_Df > 0 {
			num++
		}
	}
	if num >= len(self.PersonMgr)-1 { //庄家不下注
		self.OnEnd()
		return
	}
	self.room.flush()
}
func (self *Game_T) GameReady(uid int64) {
	if self.room.IsBye() {
		return
	}
	for i := 0; i < len(self.Ready); i++ {
		if self.Ready[i] == uid {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "同一个玩家准备")
			return
		}
	}
	//设置玩家准备状态为true
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid != uid {
			continue
		}
		self.PersonMgr[i].Ready = true
		break
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

func NewGame_T() *Game_T {
	game := new(Game_T)
	game.PersonMgr = make([]*Game_T_Person, 0)
	game.Ready = make([]int64, 0)
	//game.Winer = -1
	//game.BefDeal = -1
	return game
}
func (self *Game_T) OnInit(room *Room) {
	self.room = room
}

//每局开始
func (self *Game_T) OnBegin() {
	if self.room.IsBye() {
		return
	}
	self.room.SetBegin(true)

	DealerPos := -1 //! 庄家的位置
	WinPos := 0     //! 赢家的位置
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Total > self.PersonMgr[WinPos].Total {
			WinPos = i
		}
		if self.PersonMgr[i].Deal {
			DealerPos = i
		}
	}
	//固定选分还是自由选分
	paramge := -1
	Score := 0
	if 0 == self.room.Param1%10 {
		paramge = 0
	} else {
		paramge = 1
		if 1 == self.room.Param1%10 {
			Score = 3
		} else if 2 == self.room.Param1%10 {
			Score = 5
		} else if 3 == self.room.Param1%10 {
			Score = 7
		}
	}
	for i := 0; i < len(self.room.Uid); i++ { //! 重新初始化人
		if i >= len(self.PersonMgr) {
			person := new(Game_T_Person)
			person.Uid = self.room.Uid[i]
			if 1 == paramge {
				person.So_Df = Score
			} else {
				person.So_Df = 0
			}
			self.PersonMgr = append(self.PersonMgr, person)
		} else {
			self.PersonMgr[i].Deal = false
			self.PersonMgr[i].NextDeal = false
			self.PersonMgr[i].So_Card = 0
			self.PersonMgr[i].PaiXing = 0
			if 1 == paramge {
				self.PersonMgr[i].So_Df = Score
			} else {
				self.PersonMgr[i].So_Df = 0
			}
			self.PersonMgr[i].Total = 0
			self.PersonMgr[i].Ready = false
			self.PersonMgr[i].Card = make([]int, 0)
		}
	}
	self.ChoiceScore = paramge
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "self.score", Score)

	//下一局庄家
	if 0 == self.room.Param1/10%10 { //! 轮庄模式
		if DealerPos+1 >= len(self.PersonMgr) {
			DealerPos = -1
		}
		self.PersonMgr[DealerPos+1].Deal = true
		self.PersonMgr[DealerPos+1].DealNum++
		self.BefDeal = DealerPos + 1
		self.DealUid = self.PersonMgr[DealerPos+1].Uid
	} else if 1 == self.room.Param1/10%10 { //! 连庄模式
		self.PersonMgr[0].Deal = true
		self.PersonMgr[0].DealNum++
		self.BefDeal = 0
		self.DealUid = self.PersonMgr[0].Uid
	} else if 2 == self.room.Param1/10%10 { //! 霸王庄
		self.PersonMgr[WinPos].Deal = true
		self.PersonMgr[WinPos].DealNum++
		self.BefDeal = WinPos
		self.DealUid = self.PersonMgr[WinPos].Uid
	}
	//下一局牌 局数为奇数重新生成牌
	if 1 == self.room.Step%2 {
		//		self.Mah.Temp = make([]int, 0)
		//		self.Mah.Card = make([]int, 0)
		self.mahmgr = NewMah_TTZ()
	}
	//发牌
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Card = self.mahmgr.Deal(2)
	}
	//	if 1 == self.room.Step%2 {
	//		for i := 0; i < len(self.PersonMgr); i++ {
	//			for j := 0; j < 2; j++ {
	//				self.Mah.Temp = append(self.Mah.Temp, self.PersonMgr[i].Card[j])
	//			}
	//		}
	//	}
	self.Mah.Card = self.mahmgr.Card
	self.throwSZ() //! 掷色子，将结果发送给所有客户端
	//把消息发给玩家
	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gamebegin", self.getInfo(person.Uid))
	}
	if 1 == paramge {
		self.OnEnd()
	}
	self.room.flush()
}

//!掷色子
func (self *Game_T) throwSZ() {
	var msg Msg_GameT_SZ
	msg.Num = make([]int, 2)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 2; i++ {
		j := random.Intn(6) + 1
		msg.Num[i] = j
	}
	n := (msg.Num[0] + msg.Num[1] - 2) % (len(self.PersonMgr))
	msg.Uid = self.PersonMgr[n].Uid
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Uid_KS = self.PersonMgr[n].Uid
	}
	self.room.broadCastMsg("gameTsz", msg)
}

//每局结束
func (self *Game_T) OnEnd() {
	self.room.SetBegin(false)
	//奇数局结束存进去
	if 1 == self.room.Step%2 {
		for i := 0; i < len(self.PersonMgr); i++ {
			for j := 0; j < 2; j++ {
				self.Mah.Temp = append(self.Mah.Temp, self.PersonMgr[i].Card[j])
			}
		}
	}
	//偶数局结束清空
	if 0 == self.room.Step%2 {
		self.Mah.Temp = make([]int, 0)
		self.Mah.Card = make([]int, 0)
		//		self.mahmgr = NewMah_TTZ()
	}
	for i := 0; i < len(self.PersonMgr); i++ {
		if 37 == self.PersonMgr[i].Card[0] && 37 == self.PersonMgr[i].Card[1] {
			//双天自尊 110
			self.PersonMgr[i].So_Card = 2
			self.PersonMgr[i].NumBaoZi += 1
			self.PersonMgr[i].PaiXing = 110
		} else if self.PersonMgr[i].Card[0] == self.PersonMgr[i].Card[1] {
			//豹子 105
			self.PersonMgr[i].So_Card = 2
			self.PersonMgr[i].NumBaoZi += 1
			self.PersonMgr[i].PaiXing = 105
		} else if self.PersonMgr[i].Card[0]+self.PersonMgr[i].Card[1] == 30 {
			//二八杠 100
			if self.PersonMgr[i].Card[0] == 12 || self.PersonMgr[i].Card[0] == 18 {
				self.PersonMgr[i].So_Card = 2
				self.PersonMgr[i].PaiXing = 100
			}
		} else {
			//其他牌型 4九点半  3九点 2x点  1没点
			mahVal := 0
			if self.PersonMgr[i].Card[0] == 37 || self.PersonMgr[i].Card[1] == 37 {
				//有白板的情况
				mahVal = (self.PersonMgr[i].Card[0]+self.PersonMgr[i].Card[1]-37)%10*10 + 5
			} else {
				//没有白板的情况
				mahVal = (self.PersonMgr[i].Card[0] + self.PersonMgr[i].Card[1]) % 10 * 10
			}
			if mahVal >= 90 && mahVal <= 95 {
				self.PersonMgr[i].So_Card = 1
				self.PersonMgr[i].PaiXing = mahVal
			} else {
				self.PersonMgr[i].So_Card = 1
				self.PersonMgr[i].PaiXing = mahVal
			}
		}
	}
	//! 比较大小并计算输赢分数
	win := 0          //!记录本局庄家赢的次数
	lose := 0         //!记录本局庄家输的次数
	t := self.BefDeal //此处比较时  庄家的位置
	for i := 0; i < len(self.PersonMgr); i++ {
		if i == t {
			continue
		}
		if self.PersonMgr[i].PaiXing > self.PersonMgr[t].PaiXing {
			//!闲家赢
			winscore := self.PersonMgr[i].So_Card * self.PersonMgr[i].So_Df
			self.PersonMgr[i].Total += winscore
			self.PersonMgr[i].Score += winscore
			self.PersonMgr[i].WinNum += 1

			self.PersonMgr[t].Total -= winscore
			self.PersonMgr[t].Score -= winscore
			self.PersonMgr[t].LostNum += 1
			lose += 1
		} else {
			if (105 == self.PersonMgr[i].PaiXing) && (105 == self.PersonMgr[t].PaiXing) {
				//！都是豹子比较筒子大小
				if self.PersonMgr[i].Card[0] > self.PersonMgr[t].Card[0] {
					//!闲家赢
					winscore := self.PersonMgr[i].So_Card * self.PersonMgr[i].So_Df
					self.PersonMgr[i].Total += winscore
					self.PersonMgr[i].Score += winscore
					self.PersonMgr[i].WinNum += 1

					self.PersonMgr[t].Total -= winscore
					self.PersonMgr[t].Score -= winscore
					self.PersonMgr[t].LostNum += 1
					lose += 1
				} else {
					//！庄家赢
					winscore := self.PersonMgr[t].So_Card * self.PersonMgr[i].So_Df
					self.PersonMgr[i].Total -= winscore
					self.PersonMgr[i].Score -= winscore
					self.PersonMgr[i].LostNum += 1

					self.PersonMgr[t].Total += winscore
					self.PersonMgr[t].Score += winscore
					self.PersonMgr[t].WinNum += 1
					win += 1
				}
				continue
			}
			//!庄家赢
			winscore := self.PersonMgr[t].So_Card * self.PersonMgr[i].So_Df
			self.PersonMgr[i].Total -= winscore
			self.PersonMgr[i].Score -= winscore
			self.PersonMgr[i].LostNum += 1

			self.PersonMgr[t].Total += winscore
			self.PersonMgr[t].Score += winscore
			self.PersonMgr[t].WinNum += 1
			win += 1
		}
	}
	self.Ready = make([]int64, 0)

	//! 记录
	var record Rec_GameT
	record.Time = time.Now().Unix()
	record.Roomid = self.room.Id*100 + self.room.Step
	record.MaxStep = self.room.MaxStep

	//! 发消息
	var msg Msg_GameT_End

	if win >= len(self.PersonMgr)-1 {
		//！全赢
		self.PersonMgr[t].WinAllNum += 1
		msg.WinAll = true
		//!赢家下标
		self.Winer = t
	}
	if lose >= len(self.PersonMgr)-1 {
		//通赔
		self.PersonMgr[t].LoseAllNum += 1
		msg.LoseAll = true
		//!赢家下标
		WinPos := 0
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Total > self.PersonMgr[WinPos].Total {
				WinPos = i
			}
		}
		self.Winer = WinPos
	}

	DealerPos := -1 //! 结算后庄家的位置
	WinPos := 0     //! 结算后赢家的位置
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Total > self.PersonMgr[WinPos].Total {
			WinPos = i
		}
		if self.PersonMgr[i].Deal {
			DealerPos = i
		}
	}
	//下一局庄家
	if 0 == self.room.Param1/10%10 { //! 轮庄模式
		if DealerPos+1 >= len(self.PersonMgr) {
			DealerPos = -1
		}
		self.PersonMgr[DealerPos+1].NextDeal = true
	} else if 1 == self.room.Param1/10%10 { //! 连庄模式
		self.PersonMgr[0].NextDeal = true
	} else if 2 == self.room.Param1/10%10 { //! 霸王庄
		self.PersonMgr[WinPos].NextDeal = true
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Ready = false
		var son Son_GameT_Info
		son.Uid = self.PersonMgr[i].Uid
		son.Deal = self.PersonMgr[i].Deal
		son.Card = self.PersonMgr[i].Card
		son.Score = self.PersonMgr[i].Score
		son.So_Df = self.PersonMgr[i].So_Df
		son.Total = self.PersonMgr[i].Total
		son.NextDeal = self.PersonMgr[i].NextDeal
		son.Ready = false
		msg.Info = append(msg.Info, son)

		var rec Son_Rec_GameT
		rec.Uid = self.PersonMgr[i].Uid
		rec.Name = self.room.GetName(self.PersonMgr[i].Uid)
		rec.Head = self.room.GetHead(self.PersonMgr[i].Uid)
		rec.Card = self.PersonMgr[i].Card
		rec.Dealer = self.PersonMgr[i].Deal
		rec.Bets = self.PersonMgr[i].So_Df
		rec.Score = self.PersonMgr[i].Score
		rec.Total = self.PersonMgr[i].Total
		record.Info = append(record.Info, rec)
	}
	self.room.AddRecord(lib.HF_JtoA(&record))
	self.room.broadCastMsg("gameTend", &msg)
	if self.room.IsBye() {
		self.OnBye()
		self.room.Bye()
		return
	}
	self.room.flush()
}

//总结算
func (self *Game_T) OnBye() {
	info := make([]staticfunc.JS_CreateRoomMem, 0)
	var msg Msg_GameT_Bye
	msg.Info = make([]Son_GameT_Bye, 0)
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameT_Bye
		son.Uid = self.PersonMgr[i].Uid
		son.Total = self.PersonMgr[i].Total
		son.Score = self.PersonMgr[i].Score
		son.WinNum = self.PersonMgr[i].WinNum
		son.LostNum = self.PersonMgr[i].LostNum
		son.DealNum = self.PersonMgr[i].DealNum
		son.WinAllNum = self.PersonMgr[i].WinAllNum
		son.LoseAllNum = self.PersonMgr[i].LoseAllNum
		son.NumBaoZi = self.PersonMgr[i].NumBaoZi
		msg.Info = append(msg.Info, son)
		info = append(info, staticfunc.JS_CreateRoomMem{son.Uid, "", "", son.Score})

		GetServer().SqlScoreLog(self.PersonMgr[i].Uid, self.room.GetName(self.PersonMgr[i].Uid), self.room.GetHead(self.PersonMgr[i].Uid), self.room.Type, self.room.Id, self.PersonMgr[i].Score)
	}
	self.room.broadCastMsg("gameTbye", &msg)

	self.room.ClubResult(info)
}

//退出游戏
func (self *Game_T) OnExit(uid int64) {
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
func (self *Game_T) OnRobot(robot *lib.Robot) {

}
func (self *Game_T) OnSendInfo(person *Person) {
	person.SendMsg("gameTinfo", self.getInfo(person.Uid))
}
func (self *Game_T) getInfo(uid int64) *Msg_GameT_Info {
	var msg Msg_GameT_Info
	msg.Begin = self.room.Begin
	//	msg.LastCard = make([]int, 0)
	msg.Info = make([]Son_GameT_Info, 0)
	msg.LastCard = make([]int, 0)
	msg.ChoiceScore = self.ChoiceScore
	for _, value := range self.PersonMgr {
		var son Son_GameT_Info
		son.Uid = value.Uid
		son.Deal = value.Deal
		son.NextDeal = value.NextDeal
		son.Score = value.Score
		son.Ready = value.Ready
		son.So_Df = value.So_Df
		son.Total = value.Total
		if msg.Begin { //!游戏已经开始，只显示自己的一张牌给客户端
			if son.Uid == uid {
				for i := 0; i < len(value.Card); i++ {
					if i != len(value.Card)-1 {
						son.Card = append(son.Card, value.Card[i])
					} else {
						son.Card = append(son.Card, 0)
					}
				}
			} else {
				for i := 0; i < len(value.Card); i++ {
					son.Card = append(son.Card, 0)
				}
			}
		} else { //！游戏结束可以把当局牌全部发给客户端
			for i := 0; i < len(value.Card); i++ {
				son.Card = append(son.Card, value.Card[i])
			}
		}
		msg.Info = append(msg.Info, son)
	}
	for i := 0; i < len(self.Mah.Temp); i++ {
		msg.LastCard = append(msg.LastCard, self.Mah.Temp[i])
	}
	return &msg
}

func (self *Game_T) OnTime() {

}
func (self *Game_T) GetPerson(uid int64) *Game_T_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}
	return nil
}

func (self *Game_T) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_T) OnIsBets(uid int64) bool {
	return false
}

//! 结算所有人
func (self *Game_T) OnBalance() {
}
