package gameserver

import (
////"fmt"
//"lib"
////"math"
//"staticfunc"
//"time"
)

//! param1
//! 放炮*100 + 五八不碰*10 + 房费

//type Rec_GameHFBH_Info struct {
//	Roomid int                       `json:"roomid"` //! 房间号
//	Time   int64                     `json:"time"`   //! 记录时间
//	Param1 int                       `json:"param1"`
//	Param2 int                       `json:"param2"`
//	Jiang  int                       `json:"jiang"` //! 将牌
//	Person []Son_Rec_GameHFBH_Person `json:"person"`
//	Step   []Son_Rec_GameHFBH_Step   `json:"step"`
//	Back   []int                     `json:"back"`   //! 底牌
//	Dealer int64                     `json:"dealer"` //! 庄家
//}
//type Son_Rec_GameHFBH_Person struct {
//	Uid     int64        `json:"uid"`
//	Name    string       `json:"name"`
//	Head    string       `json:"head"`
//	Card    []int        `json:"card"`
//	Score   int          `json:"score"`
//	Total   int          `json:"total"`
//	EndCard lib.HFBHCard `json:"endcard"`
//}
//type Son_Rec_GameHFBH_Step struct { //当前动作结构
//	Uid  int64 `json:"uid"`
//	Type int   `json:"type"` //! 0摸牌 1出牌 13-15绍牌 23-25碰和招  32-33吃  99胡
//	Card []int `json:"card"`
//}

////! 房间招和绍
//type Msg_GameHFBH_ClientShao struct {
//	Num int `json:"num"`
//}
//type Msg_GameHFBH_ClientChi struct {
//	Type int `json:"type"`
//}
//type Msg_GameHFBH_ClientTuo struct {
//	Card  int `json:"card"`
//	Index int `json:"index"`
//}

////! 房间结束
//type Msg_GameHFBH_Bye struct {
//	Info []Son_GameHFBH_Bye `json:"info"`
//}
//type Son_GameHFBH_Bye struct {
//	Uid    int64 `json:"uid"`
//	Score  int   `json:"score"`
//	Hu30   int   `json:"hu30"`
//	Hu40   int   `json:"hu40"`
//	Hu50   int   `json:"hu50"`
//	HuMore int   `json:"humore"`
//	Pao    int   `json:"pao"`
//}

//type Msg_GameHFBH_Operator struct {
//	Step int `json:"step"` //! 1打牌  2下抓 3两者
//}

////! 招牌
//type Msg_GameHFBH_Zhao struct {
//	Uid      int64 `json:"uid"`
//	Card     int   `json:"card"`
//	Num      int   `json:"num"`
//	Zhao4    []int `json:"zhao4"`
//	Zhao5    []int `json:"zhao5"`
//	Shao3    []int `json:"shao3"`
//	Shao4    []int `json:"shao4"`
//	Dian     int   `json:"dian"`
//	Sequence []int `json:"sequence"`
//	Type     int   `json:"type"`
//	NoStep   []int `json:"nostep"`
//}

////! 绍牌
//type Msg_GameHFBH_Shao struct {
//	Uid      int64 `json:"uid"`
//	Num      int   `json:"num"`
//	Card     int   `json:"card"`
//	Shao3    []int `json:"shao3"`
//	Shao4    []int `json:"shao4"`
//	Shao5    []int `json:"shao5"`
//	Dian     int   `json:"dian"`
//	Sequence []int `json:"sequence"`
//	Type     int   `json:"type"` //! 0减手牌  1不减手牌
//	Find     bool  `json:"find"`
//}

////! 吃牌
//type Msg_GameHFBH_Chi struct {
//	Uid      int64   `json:"uid"`
//	Card     int     `json:"card"`
//	Chi2     [][]int `json:"chi2"`
//	Chi3     [][]int `json:"chi3"`
//	Type     int     `json:"type"`
//	Dian     int     `json:"dian"`
//	Sequence []int   `json:"sequence"`
//	NoStep   []int   `json:"nostep"`
//}

////! 拖牌
//type Msg_GameHFBH_Tuo struct {
//	Uid      int64   `json:"uid"`
//	Card     int     `json:"card"`
//	Chi2     [][]int `json:"chi2"`
//	Chi3     [][]int `json:"chi3"`
//	Dian     int     `json:"dian"`
//	Sequence []int   `json:"sequence"`
//	Index    int     `json:"index"`
//}

////! 抓牌
//type Msg_GameHFBH_Zhua struct {
//	Uid     int64 `json:"uid"`
//	ZhuaNum int   `json:"zhuanum"`
//}

////! 碰牌
//type Msg_GameHFBH_Peng struct {
//	Uid      int64 `json:"uid"`
//	Card     int   `json:"card"`
//	Peng     []int `json:"peng"`
//	Dian     int   `json:"dian"`
//	Sequence []int `json:"sequence"`
//	NoStep   []int `json:"nostep"`
//}

////! 翻牌
//type Msg_GameHFBH_Draw struct {
//	Uid      int64 `json:"uid"`  //! uid
//	Card     int   `json:"card"` //! 翻的牌
//	Hu       int   `json:"hu"`   //! 自己的操作
//	Peng     int   `json:"peng"`
//	Chi      int   `json:"chi"`
//	Shao     int   `json:"shao"`
//	Zhao     int   `json:"zhao"`
//	IsPlay   bool  `json:"isplay"` //! 是否是打出来的
//	State    int64 `json:"state"`  //! >0被uid跳了  0等待   -1没人用
//	Zhao4    []int `json:"zhao4"`
//	Tiao     []int `json:"tiao"`
//	Sequence []int `json:"sequence"`
//	Dian     int   `json:"dian"`
//}

//type Msg_GameHFBH_Deal struct {
//	Card [3]int `json:"card"` //! 发牌
//	Deal int    `json:"deal"` //! 庄家
//}

//type Msg_GameHFBH_End struct {
//	Card  int                 `json:"card"` //! 胡的牌
//	Jiang int                 `json:"jiang"`
//	Info  []Son_GameHFBH_Info `json:"info"`
//	Back  []int               `json:"back"`
//	Hu    int64               `json:"hu"` //! 胡的人uid
//}

//type Msg_GameHFBH_Info struct {
//	Begin   bool                `json:"begin"`   //! 是否开始
//	CurStep int64               `json:"curstep"` //! 当前局
//	Card    int                 `json:"card"`    //! 当前翻牌
//	Num     int                 `json:"num"`     //! 当前剩余牌数量
//	IsPlay  bool                `json:"isplay"`
//	Jiang   int                 `json:"jiang"`
//	Info    []Son_GameHFBH_Info `json:"info"`
//	Back    []int               `json:"back"`
//}
//type Son_GameHFBH_Info struct {
//	Uid     int64        `json:"uid"`
//	Deal    bool         `json:"deal"`    //! 是否庄家
//	Total   int          `json:"total"`   //! 当前总分
//	Dian    int          `json:"dian"`    //! 点数
//	ZhuaNum int          `json:"zhuanum"` //! 抓次数
//	Score   int          `json:"score"`   //! 当场积分
//	Hu      int          `json:"hu"`
//	Peng    int          `json:"peng"`
//	Chi     int          `json:"chi"`
//	Shao    int          `json:"shao"`
//	Zhao    int          `json:"zhao"`
//	Card    lib.HFBHCard `json:"card"` //! 牌结构
//	Step    int          `json:"step"`
//	Ready   bool         `json:"ready"`
//	NoStep  []int        `json:"nostep"`
//}

/////////////////////////////////////////////////////////
//type Game_HFBH_Person struct {
//	Uid     int64        `json:"uid"`
//	Deal    bool         `json:"deal"`    //! 是否庄家
//	Card    lib.HFBHCard `json:"card"`    //! 牌结构
//	Total   int          `json:"total"`   //! 总积分
//	Dian    int          `json:"dian"`    //! 点数
//	ZhuaNum int          `json:"zhuanum"` //! 开局抓次数
//	Score   int          `json:"score"`   //! 当场积分
//	Peng    int          `json:"peng"`    //! 0不能碰  1可以碰  2已选择了碰
//	Chi     int          `json:"chi"`     //! 0不能吃  1可以吃  2已经选择了吃
//	Shao    int          `json:"shao"`    //! 0不能绍  1可以绍  2已经选择了绍
//	Zhao    int          `json:"zhao"`    //! 0不能招  1可以招  2已经选择了招
//	Hu      int          `json:"hu"`      //! 0不能胡  1可以胡
//	Step    int          `json:"step"`    //! 出牌标志  0没有选择  1仅打  2仅抓  3可打可抓
//	Ready   bool         `json:"ready"`   //! 是否准备
//	Zhua    int          `json:"zhua"`    //! 有几次机会可以下抓
//	NoStep  []int        `json:"nostep"`  //! 不能打
//	NoZOP   []int        `json:"nozop"`   //! 不能碰或招
//	Steps   []int        `json:"steps"`   //! 打过的牌
//	Chis    []int        `json:"chis"`    //! 吃过的牌
//	HuCard  lib.HFBHCard `json:"hucard"`  //! 胡牌的结构
//	Hu30    int          `json:"hu30"`    //! 30胡次数
//	Hu40    int          `json:"hu40"`    //! 40胡次数
//	Hu50    int          `json:"hu50"`    //! 50胡次数
//	HuMore  int          `json:"humore"`  //! 超过50胡
//	Pao     int          `json:"pao"`     //! 点炮次数
//}

//func (self *Game_HFBH_Person) Init() {
//	self.Deal = false
//	self.Score = 0
//	self.Peng = 0
//	self.Chi = 0
//	self.Shao = 0
//	self.Zhao = 0
//	self.Ready = false
//	self.Step = 0
//	self.Dian = 0
//	self.ZhuaNum = 0
//	self.Hu = 0
//	self.NoStep = make([]int, 0)
//	self.NoZOP = make([]int, 0)
//	self.Steps = make([]int, 0)
//	self.Chis = make([]int, 0)
//	self.Card.Init()
//}

//func (self *Game_HFBH_Person) Reset() {
//	self.Peng = 0
//	self.Chi = 0
//	self.Shao = 0
//	self.Zhao = 0
//}

//func (self *Game_HFBH_Person) IsCanPlay() bool {
//	if len(self.Card.Card1) == 0 {
//		return false
//	}

//	for i := 0; i < len(self.Card.Card1); i++ {
//		find := false
//		for j := 0; j < len(self.NoStep); j++ {
//			if self.Card.Card1[i] == self.NoStep[j] {
//				find = true
//				break
//			}
//		}
//		if !find {
//			return true
//		}
//	}
//	return false
//}

//type Game_HFBH struct {
//	PersonMgr []*Game_HFBH_Person `json:"personmgr"`
//	Mgr       *CardMgr            `json:"mgr"`     //! 剩余
//	CurStep   int64               `json:"curstep"` //! 当前谁的局
//	BefStep   int64               `json:"befstep"` //! 上局谁出
//	Winer     int64               `json:"winer"`   //! 上局谁赢
//	IsPlay    bool                `json:"isplay"`
//	Card      int                 `json:"card"`   //! 当前翻牌
//	Jiang     int                 `json:"jiang"`  //! 将
//	Record    *Rec_GameHFBH_Info  `json:"record"` //! 记录

//	room *Room
//}

//func NewGame_HFBH() *Game_HFBH {
//	game := new(Game_HFBH)
//	game.PersonMgr = make([]*Game_HFBH_Person, 0)

//	return game
//}

//func (self *Game_HFBH) GetPerson(uid int64) *Game_HFBH_Person {
//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Uid == uid {
//			return self.PersonMgr[i]
//		}
//	}

//	return nil
//}

////! 得到上一个uid
//func (self *Game_HFBH) GetBeforePerson(uid int64) *Game_HFBH_Person {
//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Uid != uid {
//			continue
//		}

//		if i-1 < 0 {
//			return self.PersonMgr[len(self.PersonMgr)-1]
//		} else {
//			return self.PersonMgr[i-1]
//		}
//	}

//	return nil
//}

////! 得到下一个uid
//func (self *Game_HFBH) GetNextUid(uid int64) int64 {
//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Uid != uid {
//			continue
//		}

//		if i+1 < len(self.PersonMgr) {
//			return self.PersonMgr[i+1].Uid
//		} else {
//			return self.PersonMgr[0].Uid
//		}
//	}

//	return 0
//}

////! 得到下一个person
//func (self *Game_HFBH) GetNextPerson(uid int64) *Game_HFBH_Person {
//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Uid != uid {
//			continue
//		}

//		if i+1 < len(self.PersonMgr) {
//			return self.PersonMgr[i+1]
//		} else {
//			return self.PersonMgr[0]
//		}
//	}

//	return nil
//}

//func (self *Game_HFBH) OnInit(room *Room) {
//	self.room = room
//}

//func (self *Game_HFBH) OnRobot(robot *lib.Robot) {

//}

//func (self *Game_HFBH) OnSendInfo(person *Person) {
//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Uid == person.Uid {
//			person.SendMsg("gamehfbhinfo", self.getInfo(person.Uid))
//			return
//		}
//	}

//	_person := new(Game_HFBH_Person)
//	_person.Uid = person.Uid
//	_person.Init()
//	self.PersonMgr = append(self.PersonMgr, _person)

//	person.SendMsg("gamehfbhinfo", self.getInfo(person.Uid))
//}

//func (self *Game_HFBH) OnMsg(msg *RoomMsg) {
//	switch msg.Head {
//	case "gameready": //! 游戏准备
//		self.GameReady(msg.Uid)
//		self.room.flush()
//	case "gameunready": //! 取消准备
//		self.GameUnReady(msg.Uid)
//		self.room.flush()
//	case "gamestep": //! 出牌
//		self.GameStep(msg.Uid, msg.V.(*Msg_GameStep).Card)
//		self.room.flush()
//	case "gamepeng": //! 碰
//		self.GamePeng(msg.Uid)
//		self.room.flush()
//	case "gamehu": //! 胡
//		self.GameHu(msg.Uid)
//		self.room.flush()
//	case "gamezhao":
//		self.GameZhao(msg.Uid, msg.V.(*Msg_GameHFBH_ClientShao).Num)
//		self.room.flush()
//	case "gameshao":
//		self.GameShao(msg.Uid, msg.V.(*Msg_GameHFBH_ClientShao).Num)
//		self.room.flush()
//	case "gamehfbhchi":
//		self.GameChi(msg.Uid, msg.V.(*Msg_GameHFBH_ClientChi).Type)
//		self.room.flush()
//	case "gameguo": //! 过
//		self.GameGuo(msg.Uid)
//		self.room.flush()
//	case "gamezhua": //! 下抓
//		self.GameZhua(msg.Uid)
//		self.room.flush()
//	case "gametuo": //! 拖
//		self.GameTuo(msg.Uid, msg.V.(*Msg_GameHFBH_ClientTuo).Card, msg.V.(*Msg_GameHFBH_ClientTuo).Index)
//		self.room.flush()
//	}
//}

//func (self *Game_HFBH) OnBegin() {
//	if self.room.IsBye() {
//		return
//	}

//	self.room.SetBegin(true)
//	if self.Winer == 0 { //! 发牌决定谁是庄家
//		index := 0
//		var card [3]int
//		mgr := NewCard_HFBH()
//		card[0] = mgr.Deal(1)[0]
//		card[1] = mgr.Deal(1)[0]
//		card[2] = mgr.Deal(1)[0]
//		if card[1] < card[0] {
//			index = 1
//		}
//		if card[2] < card[index] {
//			index = 2
//		}
//		self.Winer = self.room.Uid[index]

//		var msg Msg_GameHFBH_Deal
//		msg.Card = card
//		msg.Deal = index
//		self.room.broadCastMsg("gamehfbhdeal", &msg)
//	}

//	self.Jiang = (lib.HF_GetRandom(8)+1)*10 + lib.HF_GetRandom(3) + 1
//	self.Mgr = NewCard_HFBH()
//	self.BefStep = 0
//	self.Card = 0
//	self.CurStep = self.Winer
//	for i := 0; i < len(self.PersonMgr); i++ {
//		self.PersonMgr[i].Init()
//		if self.PersonMgr[i].Uid == self.Winer {
//			self.PersonMgr[i].Deal = true
//			self.PersonMgr[i].Card.Card1 = self.Mgr.Deal(31)
//			self.PersonMgr[i].Card.Sort()
//			self.PersonMgr[i].Zhua = self.PersonMgr[i].Card.GetZhuaNum()
//			self.PersonMgr[i].Step = 1
//		} else {
//			self.PersonMgr[i].Card.Card1 = self.Mgr.Deal(30)
//			self.PersonMgr[i].Card.Sort()
//			self.PersonMgr[i].Zhua = self.PersonMgr[i].Card.GetZhuaNum()
//		}
//	}
//	self.Winer = 0

//	for i := 0; i < len(self.PersonMgr); i++ {
//		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
//		if person == nil {
//			continue
//		}
//		person.SendMsg("gamehfbhbegin", self.getInfo(person.Uid))
//	}

//	self.Record = new(Rec_GameHFBH_Info)
//	self.Record.Roomid = self.room.Id*100 + self.room.Step
//	self.Record.Time = time.Now().Unix()
//	self.Record.Param1 = self.room.Param1
//	self.Record.Param2 = self.room.Param2
//	self.Record.Jiang = self.Jiang
//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Deal {
//			self.Record.Dealer = self.PersonMgr[i].Uid
//		}
//		var node Son_Rec_GameHFBH_Person
//		node.Uid = self.PersonMgr[i].Uid
//		node.Name = self.room.GetName(self.PersonMgr[i].Uid)
//		node.Head = self.room.GetHead(self.PersonMgr[i].Uid)
//		node.Card = self.PersonMgr[i].Card.Card1
//		self.Record.Person = append(self.Record.Person, node)
//	}
//}

//func (self *Game_HFBH) OnEnd() {
//	self.room.SetBegin(false)

//	if self.Winer != 0 { //! 没有荒庄
//		allscore := 0
//		person := self.GetPerson(self.Winer)
//		hu := (person.Card.Score + 5) / 10 * 10
//		for i := 0; i < len(self.PersonMgr); i++ {
//			if self.PersonMgr[i].Uid == self.Winer {
//				continue
//			}

//			if self.BefStep != 0 {
//				if self.PersonMgr[i].Uid == self.BefStep {
//					self.PersonMgr[i].Pao++
//					score := hu*2 + self.room.Param1/100*10
//					allscore += score
//					self.PersonMgr[i].Score = -score
//					self.PersonMgr[i].Total += self.PersonMgr[i].Score
//					break
//				}
//			} else {
//				score := hu
//				allscore += score
//				self.PersonMgr[i].Score = -score
//				self.PersonMgr[i].Total += self.PersonMgr[i].Score
//			}
//		}
//		if hu > 50 {
//			person.HuMore++
//		} else if hu == 50 {
//			person.Hu50++
//		} else if hu == 40 {
//			person.Hu40++
//		} else if hu == 30 {
//			person.Hu30++
//		}
//		person.Score = allscore
//		person.Total += person.Score
//	}

//	//! 发消息
//	var msg Msg_GameHFBH_End
//	msg.Hu = self.Winer
//	msg.Card = self.Card
//	msg.Jiang = self.Jiang
//	msg.Back = self.Mgr.Card
//	for i := 0; i < len(self.PersonMgr); i++ {
//		var son Son_GameHFBH_Info
//		son.Uid = self.PersonMgr[i].Uid
//		son.Deal = self.PersonMgr[i].Deal
//		son.Card = self.PersonMgr[i].Card
//		son.Score = self.PersonMgr[i].Score
//		son.Total = self.PersonMgr[i].Total
//		msg.Info = append(msg.Info, son)

//		self.Record.Person[i].EndCard = self.PersonMgr[i].Card
//		self.Record.Person[i].Score = self.PersonMgr[i].Score
//		self.Record.Person[i].Total = self.PersonMgr[i].Total
//	}
//	self.room.broadCastMsg("gamehfbhend", &msg)
//	self.room.AddRecord(lib.HF_JtoA(self.Record))

//	if self.room.IsBye() {
//		self.OnBye()
//		self.room.Bye()
//		return
//	}
//}

//func (self *Game_HFBH) OnBye() {
//	info := make([]staticfunc.JS_CreateRoomMem, 0)
//	var msg Msg_GameHFBH_Bye
//	for i := 0; i < len(self.PersonMgr); i++ {
//		var son Son_GameHFBH_Bye
//		son.Hu30 = self.PersonMgr[i].Hu30
//		son.Hu40 = self.PersonMgr[i].Hu40
//		son.Hu50 = self.PersonMgr[i].Hu50
//		son.HuMore = self.PersonMgr[i].HuMore
//		son.Pao = self.PersonMgr[i].Pao
//		son.Uid = self.PersonMgr[i].Uid
//		son.Score = self.PersonMgr[i].Total
//		msg.Info = append(msg.Info, son)
//		info = append(info, staticfunc.JS_CreateRoomMem{son.Uid, "", "", son.Score})
//	}
//	self.room.broadCastMsg("gamehfbhbye", &msg)

//	self.room.ClubResult(info)
//}

//func (self *Game_HFBH) OnExit(uid int64) {
//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Uid == uid {
//			copy(self.PersonMgr[i:], self.PersonMgr[i+1:])
//			self.PersonMgr = self.PersonMgr[:len(self.PersonMgr)-1]
//			break
//		}
//	}
//}

////! 准备
//func (self *Game_HFBH) GameReady(uid int64) {
//	if self.room.IsBye() {
//		return
//	}

//	if self.room.Begin {
//		return
//	}

//	num := 0
//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Uid == uid {
//			if self.PersonMgr[i].Ready {
//				return
//			} else {
//				self.PersonMgr[i].Ready = true
//				num++
//			}
//		} else if self.PersonMgr[i].Ready {
//			num++
//		}
//	}

//	if num == len(self.room.Uid) && len(self.PersonMgr) >= lib.HF_Atoi(self.room.csv["minnum"]) {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
//		self.OnBegin()
//		return
//	}

//	var msg staticfunc.Msg_Uid
//	msg.Uid = uid
//	self.room.broadCastMsg("gameready", &msg)
//}

////! 取消准备
//func (self *Game_HFBH) GameUnReady(uid int64) {
//	if self.room.IsBye() {
//		return
//	}

//	if self.room.Begin {
//		return
//	}

//	person := self.GetPerson(uid)
//	if person == nil {
//		return
//	}

//	person.Ready = false

//	var msg staticfunc.Msg_Uid
//	msg.Uid = uid
//	self.room.broadCastMsg("gameunready", &msg)
//}

////! 摸牌阶段
//func (self *Game_HFBH) GameMoCard(uid int64) {
//	if len(self.Mgr.Card) == 0 { //! 荒庄
//		self.OnEnd()
//		return
//	}

//	//! 每一轮先还原所有人的操作
//	self.Reset()

//	self.CurStep = uid
//	self.BefStep = 0
//	self.Card = self.Mgr.Deal(1)[0]
//	self.IsPlay = false
//	self.AddStep(uid, 0, []int{self.Card})

//	//! 判断当前的人是否能胡
//	person := self.GetPerson(uid)
//	hu, peng, zhao, chi := self.IsHu(person)
//	if hu {
//		person.Card.Mo = self.Card
//		lst := lib.GetHFBHMgr().IsHu(person.Card, true, peng, zhao, chi, self.Jiang)
//		if len(lst) > 0 { //! 能胡
//			person.HuCard = lst[0]
//			person.Hu = 1
//			for i := 0; i < len(self.PersonMgr); i++ {
//				var msg Msg_GameHFBH_Draw
//				msg.IsPlay = false
//				msg.Uid = uid
//				msg.State = 0
//				if self.PersonMgr[i].Uid == uid { //! 发给自己
//					msg.Card = self.Card
//					msg.Hu = self.PersonMgr[i].Hu
//				} else {
//					msg.Card = -1
//					msg.Hu = 0
//				}
//				self.room.SendMsg(self.PersonMgr[i].Uid, "gamehfbhdraw", &msg)
//			}
//			return
//		}
//	}

//	//! 不能胡就判断是否能绍
//	if person.Card.IsCanShao(self.Card) {
//		person.Shao = 1
//		if person.Card.IsCanChi(self.Card, self.GetNoChi(person)) {
//			person.Chi = 1
//		}
//		for i := 0; i < len(self.PersonMgr); i++ {
//			var msg Msg_GameHFBH_Draw
//			msg.IsPlay = false
//			msg.Uid = uid
//			msg.State = 0
//			if self.PersonMgr[i].Uid == uid { //! 发给自己
//				msg.Card = self.Card
//				msg.Shao = self.PersonMgr[i].Shao
//				msg.Chi = self.PersonMgr[i].Chi
//			} else {
//				msg.Card = -1
//				msg.Shao = 0
//				msg.Chi = 0
//			}
//			self.room.SendMsg(self.PersonMgr[i].Uid, "gamehfbhdraw", &msg)
//		}
//		return
//	}

//	//! 先判断胡
//	find := false
//	nextperson := self.GetNextPerson(uid)
//	for nextperson != person {
//		hu, peng, zhao, chi := self.IsHu(nextperson)
//		if hu {
//			nextperson.Card.Mo = self.Card
//			lst := lib.GetHFBHMgr().IsHu(nextperson.Card, false, peng, zhao, chi, self.Jiang)
//			if len(lst) > 0 {
//				nextperson.HuCard = lst[0]
//				nextperson.Hu = 1
//				find = true
//				break
//			}
//		}
//		nextperson = self.GetNextPerson(nextperson.Uid)
//	}
//	if find {
//		for i := 0; i < len(self.PersonMgr); i++ {
//			var msg Msg_GameHFBH_Draw
//			msg.IsPlay = false
//			msg.Uid = uid
//			msg.State = 0
//			msg.Card = self.Card
//			msg.Hu = self.PersonMgr[i].Hu
//			msg.Shao = 0
//			msg.Peng = 0
//			msg.Chi = 0
//			msg.Zhao = 0
//			self.room.SendMsg(self.PersonMgr[i].Uid, "gamehfbhdraw", &msg)
//		}
//		return
//	}

//	//! 能跳就直接跳
//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Card.IsCanTiao(self.Card) {
//			self.PersonMgr[i].Card.TiaoCard(self.Card)
//			self.CurStep = self.PersonMgr[i].Uid
//			var msg Msg_GameHFBH_Draw
//			msg.IsPlay = false
//			msg.Uid = self.PersonMgr[i].Uid
//			msg.State = self.PersonMgr[i].Uid
//			msg.Card = self.Card
//			msg.Hu = 0
//			msg.Shao = 0
//			msg.Peng = 0
//			msg.Chi = 0
//			msg.Zhao = 0
//			msg.Zhao4 = self.PersonMgr[i].Card.Zhao4
//			msg.Tiao = self.PersonMgr[i].Card.Tiao
//			msg.Sequence = self.PersonMgr[i].Card.Sequence
//			msg.Dian = self.PersonMgr[i].Card.CountDian(self.Jiang)
//			self.room.broadCastMsg("gamehfbhdraw", &msg)

//			self.GameMoCard(self.GetNextUid(self.CurStep))
//			return
//		}
//	}

//	//! 不能胡不能绍,轮流判断
//	//! 再判断其他情况
//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Uid == uid { //! 是自己判断是否能吃
//			if self.PersonMgr[i].Card.IsCanChi(self.Card, self.GetNoChi(self.PersonMgr[i])) {
//				self.PersonMgr[i].Chi = 1
//				find = true
//			}
//		} else { //!
//			if self.PersonMgr[i].Card.IsCanPeng(self.Card, self.room.Param1/10%10 == 1, self.GetNoPengOrZhao(self.PersonMgr[i])) == 2 {
//				self.PersonMgr[i].Peng = 1
//				find = true
//			}
//			if self.PersonMgr[i].Card.IsCanZhao(self.Card, self.GetNoPengOrZhao(self.PersonMgr[i])) {
//				self.PersonMgr[i].Zhao = 1
//				find = true
//			}
//			if self.PersonMgr[i].Uid == self.GetNextUid(self.CurStep) && self.PersonMgr[i].Card.IsCanChi(self.Card, self.GetNoChi(self.PersonMgr[i])) {
//				self.PersonMgr[i].Chi = 1
//				find = true
//			}
//		}
//	}
//	if find {
//		for i := 0; i < len(self.PersonMgr); i++ {
//			var msg Msg_GameHFBH_Draw
//			msg.IsPlay = false
//			msg.Uid = uid
//			msg.State = 0
//			msg.Card = self.Card
//			msg.Hu = 0
//			msg.Shao = 0
//			msg.Peng = self.PersonMgr[i].Peng
//			msg.Chi = self.PersonMgr[i].Chi
//			msg.Zhao = self.PersonMgr[i].Zhao
//			self.room.SendMsg(self.PersonMgr[i].Uid, "gamehfbhdraw", &msg)
//		}
//		return
//	}

//	person.Card.Card2 = append(person.Card.Card2, self.Card)
//	var msg Msg_GameHFBH_Draw
//	msg.IsPlay = false
//	msg.Uid = uid
//	msg.State = -1
//	msg.Card = self.Card
//	msg.Hu = 0
//	msg.Shao = 0
//	msg.Peng = 0
//	msg.Chi = 0
//	msg.Zhao = 0
//	self.room.broadCastMsg("gamehfbhdraw", &msg)

//	self.GameMoCard(self.GetNextUid(self.CurStep))
//}

////! 出牌
//func (self *Game_HFBH) GameStep(uid int64, card int) {
//	if !self.room.Begin {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
//		return
//	}

//	if card == 0 {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "GameStep(card=0)")
//		return
//	}

//	person := self.GetPerson(uid)
//	if person == nil {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person:GameStep")
//		return
//	}

//	if person.Step == 0 || person.Step == 2 {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能出牌:GameStep")
//		return
//	}

//	for i := 0; i < len(person.NoStep); i++ {
//		if person.NoStep[i] == card {
//			lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能出吃的牌:GameStep")
//			return
//		}
//	}

//	if !person.Card.StepCard(card) {
//		return
//	}

//	person.Steps = append(person.Steps, card)
//	self.Card = card
//	person.Step = 0
//	if len(person.Card.Card1) < 30 {
//		person.Zhua = 0
//	}
//	self.CurStep = person.Uid
//	self.BefStep = person.Uid
//	self.IsPlay = true
//	self.AddStep(uid, 1, []int{self.Card})

//	//! 每一轮先还原所有人的操作
//	self.Reset()

//	//! 先判断胡
//	find := false
//	nextperson := self.GetNextPerson(uid)
//	for nextperson != person {
//		hu, peng, zhao, chi := self.IsHu(nextperson)
//		if hu {
//			nextperson.Card.Mo = self.Card
//			lst := lib.GetHFBHMgr().IsHu(nextperson.Card, false, peng, zhao, chi, self.Jiang)
//			if len(lst) > 0 {
//				nextperson.HuCard = lst[0]
//				nextperson.Hu = 1
//				find = true
//				break
//			}
//		}
//		nextperson = self.GetNextPerson(nextperson.Uid)
//	}
//	if find {
//		for i := 0; i < len(self.PersonMgr); i++ {
//			var msg Msg_GameHFBH_Draw
//			msg.IsPlay = true
//			msg.Uid = uid
//			msg.State = 0
//			msg.Card = self.Card
//			msg.Hu = self.PersonMgr[i].Hu
//			msg.Shao = 0
//			msg.Peng = 0
//			msg.Chi = 0
//			msg.Zhao = 0
//			self.room.SendMsg(self.PersonMgr[i].Uid, "gamehfbhstep", &msg)
//		}
//		return
//	}

//	//! 不能胡不能绍,轮流判断
//	//! 再判断其他情况
//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Uid == uid { //! 是自己判断是否能吃
//			continue
//		} else { //!
//			if self.PersonMgr[i].Card.IsCanPeng(self.Card, self.room.Param1/10%10 == 1, self.GetNoPengOrZhao(self.PersonMgr[i])) == 2 {
//				self.PersonMgr[i].Peng = 1
//				find = true
//			}
//			if self.PersonMgr[i].Card.IsCanZhao(self.Card, self.GetNoPengOrZhao(self.PersonMgr[i])) {
//				self.PersonMgr[i].Zhao = 1
//				find = true
//			}
//			if self.PersonMgr[i].Uid == self.GetNextUid(self.CurStep) && self.PersonMgr[i].Card.IsCanChi(self.Card, self.GetNoChi(self.PersonMgr[i])) {
//				self.PersonMgr[i].Chi = 1
//				find = true
//			}
//		}
//	}
//	if find {
//		for i := 0; i < len(self.PersonMgr); i++ {
//			var msg Msg_GameHFBH_Draw
//			msg.IsPlay = true
//			msg.Uid = uid
//			msg.State = 0
//			msg.Card = self.Card
//			msg.Hu = 0
//			msg.Shao = 0
//			msg.Peng = self.PersonMgr[i].Peng
//			msg.Chi = self.PersonMgr[i].Chi
//			msg.Zhao = self.PersonMgr[i].Zhao
//			self.room.SendMsg(self.PersonMgr[i].Uid, "gamehfbhstep", &msg)
//		}
//		return
//	}

//	person.Card.Card2 = append(person.Card.Card2, self.Card)
//	var msg Msg_GameHFBH_Draw
//	msg.IsPlay = true
//	msg.Uid = uid
//	msg.State = -1
//	msg.Card = self.Card
//	msg.Hu = 0
//	msg.Shao = 0
//	msg.Peng = 0
//	msg.Chi = 0
//	msg.Zhao = 0
//	self.room.broadCastMsg("gamehfbhstep", &msg)

//	self.GameMoCard(self.GetNextUid(self.CurStep))
//}

////! 碰
//func (self *Game_HFBH) GamePeng(uid int64) {
//	if !self.room.Begin {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始:GamePeng")
//		return
//	}

//	if self.CurStep == uid {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "自己局不能碰:GamePeng")
//		return
//	}

//	if self.Card == 0 {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "错误的碰牌:GamePeng")
//		return
//	}

//	person := self.GetPerson(uid)
//	if person == nil {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person:GamePeng")
//		return
//	}

//	if person.Peng != 1 {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能碰:GamePeng")
//		return
//	}

//	//! 判断是否可以碰
//	nextuid := self.GetNextUid(self.CurStep)
//	nextperson := self.GetPerson(nextuid)
//	if nextperson == person || nextperson.Peng == 0 { //! 自己是下家或者下家不能碰
//		self.PengCard(person)
//		if person.Zhua > 0 { //! 可选择下抓
//			person.Step = 3
//		} else { //! 只能打牌
//			person.Step = 1
//		}
//		self.SendOperator(person)
//	} else { //! 只能标记了碰,但是要等待
//		person.Reset()
//		person.Peng = 2
//	}
//}

////! 绍
//func (self *Game_HFBH) GameShao(uid int64, num int) {
//	if !self.room.Begin {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始:GameShao")
//		return
//	}

//	if self.Card == 0 {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "杠的不是最后一张牌:GameShao")
//		return
//	}

//	if self.CurStep != uid {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不是自己的局:GameShao")
//		return
//	}

//	person := self.GetPerson(uid)
//	if person == nil {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person:GameShao")
//		return
//	}

//	if person.Shao != 1 {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能绍:GameShao")
//		return
//	}

//	result := false
//	_type := 0
//	if num == 3 {
//		result, _type = person.Card.Shao3Card(self.Card)
//	} else if num == 4 {
//		result, _type = person.Card.Shao4Card(self.Card)
//	} else if num == 5 {
//		result, _type = person.Card.Shao5Card(self.Card)
//	} else {
//		return
//	}

//	if !result {
//		return
//	}

//	find := false
//	for i := 0; i < len(person.NoZOP); i++ {
//		if person.NoZOP[i] == self.Card {
//			find = true //! 后绍
//			copy(person.NoZOP[i:], person.NoZOP[i+1:])
//			person.NoZOP = person.NoZOP[:len(person.NoZOP)-1]
//			break
//		}
//	}

//	self.Reset()
//	person.Dian = person.Card.CountDian(self.Jiang)
//	self.CurStep = person.Uid
//	self.AddStep(uid, 10+num, []int{self.Card})

//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Uid == person.Uid {
//			var msg Msg_GameHFBH_Shao
//			msg.Uid = person.Uid
//			msg.Num = num
//			msg.Card = self.Card
//			msg.Shao3 = person.Card.Shao3
//			msg.Shao4 = person.Card.Shao4
//			msg.Shao5 = person.Card.Shao5
//			msg.Dian = person.Dian
//			msg.Sequence = person.Card.Sequence
//			msg.Type = _type
//			msg.Find = find
//			self.room.SendMsg(self.PersonMgr[i].Uid, "gamehfbhshao", &msg)
//		} else {
//			var msg Msg_GameHFBH_Shao
//			msg.Uid = person.Uid
//			msg.Num = num
//			msg.Card = 0
//			msg.Shao3 = make([]int, len(person.Card.Shao3))
//			msg.Shao4 = make([]int, len(person.Card.Shao4))
//			msg.Shao5 = make([]int, len(person.Card.Shao5))
//			msg.Dian = person.Dian
//			msg.Sequence = person.Card.Sequence
//			msg.Type = 0
//			msg.Find = find
//			self.room.SendMsg(self.PersonMgr[i].Uid, "gamehfbhshao", &msg)
//		}
//	}

//	self.Card = 0

//	if person.Card.GetZhuaNum() > 0 {
//		person.Step = 3
//	} else {
//		person.Step = 1
//	}
//	self.SendOperator(person)
//}

////! 招
//func (self *Game_HFBH) GameZhao(uid int64, num int) {
//	if !self.room.Begin {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始:GameZhao")
//		return
//	}

//	if self.Card == 0 {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "招的不是最后一张牌:GameZhao")
//		return
//	}

//	if self.CurStep == uid {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "是自己的局:GameZhao")
//		return
//	}

//	person := self.GetPerson(uid)
//	if person == nil {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person:GameZhao")
//		return
//	}

//	if person.Zhao != 1 {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能招:GameZhao")
//		return
//	}

//	result := false
//	_type := 0
//	if num == 4 {
//		result, _type = person.Card.Zhao4Card(self.Card)
//	} else if num == 5 {
//		result, _type = person.Card.Zhao5Card(self.Card)
//	} else {
//		return
//	}

//	if !result {
//		return
//	}

//	person.NoStep = append(person.NoStep, self.Card)
//	self.Reset()
//	person.Dian = person.Card.CountDian(self.Jiang)
//	self.CurStep = person.Uid
//	self.AddStep(uid, 20+num, []int{self.Card})

//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Uid == person.Uid {
//			var msg Msg_GameHFBH_Zhao
//			msg.Uid = person.Uid
//			msg.Card = self.Card
//			msg.Num = num
//			msg.Zhao4 = person.Card.Zhao4
//			msg.Zhao5 = person.Card.Zhao5
//			msg.Shao3 = person.Card.Shao3
//			msg.Shao4 = person.Card.Shao4
//			msg.Dian = person.Dian
//			msg.Sequence = person.Card.Sequence
//			msg.Type = _type
//			msg.NoStep = person.NoStep
//			self.room.SendMsg(self.PersonMgr[i].Uid, "gamehfbhzhao", &msg)
//		} else {
//			var msg Msg_GameHFBH_Zhao
//			msg.Uid = person.Uid
//			msg.Card = self.Card
//			msg.Num = num
//			msg.Zhao4 = person.Card.Zhao4
//			msg.Zhao5 = person.Card.Zhao5
//			msg.Shao3 = make([]int, len(person.Card.Shao3))
//			msg.Shao4 = make([]int, len(person.Card.Shao4))
//			msg.Dian = person.Dian
//			msg.Sequence = person.Card.Sequence
//			msg.Type = _type
//			msg.NoStep = person.NoStep
//			self.room.SendMsg(self.PersonMgr[i].Uid, "gamehfbhzhao", &msg)
//		}
//	}

//	self.Card = 0

//	c4 := len(person.Card.Zhao4) + len(person.Card.Shao4)
//	c5 := len(person.Card.Zhao5) + len(person.Card.Shao5) + len(person.Card.Tiao)
//	if c4+c5 >= 2 {
//		person.Step = 2
//	} else if person.Card.GetZhuaNum() > 0 {
//		person.Step = 3
//	} else {
//		person.Step = 1
//	}

//	self.SendOperator(person)
//}

////! 吃
////! 100chi2去吃  102手牌两张吃 <100用哪张牌吃
//func (self *Game_HFBH) GameChi(uid int64, _type int) {
//	if !self.room.Begin {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始:GameChi")
//		return
//	}

//	if self.Card == 0 {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "没牌吃:GameChi")
//		return
//	}

//	person := self.GetPerson(uid)
//	if person == nil {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person:GameChi")
//		return
//	}

//	if person.Chi != 1 {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能吃:GameChi")
//		return
//	}

//	if _type == 0 && self.CurStep != uid {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不是自己的局:GameChi")
//		return
//	}

//	if _type < 100 && len(person.Card.Chi2) >= 2 {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能吃2口以上:GameChi")
//		return
//	}

//	if person.Shao == 1 { //! 能绍的情况下点了吃
//		person.Reset()
//		person.Chi = _type
//		//! 先判断胡
//		find := false
//		nextperson := self.GetNextPerson(uid)
//		for nextperson != person {
//			hu, peng, zhao, chi := self.IsHu(nextperson)
//			if hu {
//				nextperson.Card.Mo = self.Card
//				lst := lib.GetHFBHMgr().IsHu(nextperson.Card, false, peng, zhao, chi, self.Jiang)
//				if len(lst) > 0 {
//					nextperson.HuCard = lst[0]
//					nextperson.Hu = 1
//					find = true
//					break
//				}
//			}
//			nextperson = self.GetNextPerson(nextperson.Uid)
//		}
//		if find {
//			for i := 0; i < len(self.PersonMgr); i++ {
//				var msg Msg_GameHFBH_Draw
//				msg.IsPlay = false
//				msg.Uid = uid
//				msg.State = 0
//				msg.Card = self.Card
//				msg.Hu = self.PersonMgr[i].Hu
//				msg.Shao = 0
//				msg.Peng = 0
//				msg.Chi = 0
//				msg.Zhao = 0
//				self.room.SendMsg(self.PersonMgr[i].Uid, "gamehfbhfan", &msg)
//			}
//			return
//		}

//		//! 能跳就直接跳
//		for i := 0; i < len(self.PersonMgr); i++ {
//			if self.PersonMgr[i].Uid == uid {
//				continue
//			}
//			if self.PersonMgr[i].Card.IsCanTiao(self.Card) {
//				self.PersonMgr[i].Card.TiaoCard(self.Card)
//				self.CurStep = self.PersonMgr[i].Uid
//				var msg Msg_GameHFBH_Draw
//				msg.IsPlay = false
//				msg.Uid = self.PersonMgr[i].Uid
//				msg.State = self.PersonMgr[i].Uid
//				msg.Card = self.Card
//				msg.Hu = 0
//				msg.Shao = 0
//				msg.Peng = 0
//				msg.Chi = 0
//				msg.Zhao = 0
//				msg.Zhao4 = self.PersonMgr[i].Card.Zhao4
//				msg.Tiao = self.PersonMgr[i].Card.Tiao
//				msg.Sequence = self.PersonMgr[i].Card.Sequence
//				msg.Dian = self.PersonMgr[i].Card.CountDian(self.Jiang)
//				self.room.broadCastMsg("gamehfbhfan", &msg)

//				self.GameMoCard(self.GetNextUid(self.CurStep))
//				return
//			}
//		}

//		//! 不能胡不能绍,轮流判断
//		//! 再判断其他情况
//		for i := 0; i < len(self.PersonMgr); i++ {
//			if self.PersonMgr[i].Uid == uid { //! 是自己判断是否能吃
//				continue
//			} else { //!
//				if self.PersonMgr[i].Card.IsCanPeng(self.Card, self.room.Param1/10%10 == 1, self.GetNoPengOrZhao(self.PersonMgr[i])) == 2 {
//					self.PersonMgr[i].Peng = 1
//					find = true
//				}
//				if self.PersonMgr[i].Card.IsCanZhao(self.Card, self.GetNoPengOrZhao(self.PersonMgr[i])) {
//					self.PersonMgr[i].Zhao = 1
//					find = true
//				}
//			}
//		}
//		if find {
//			for i := 0; i < len(self.PersonMgr); i++ {
//				var msg Msg_GameHFBH_Draw
//				msg.IsPlay = false
//				msg.Uid = uid
//				msg.State = 0
//				msg.Card = self.Card
//				msg.Hu = 0
//				msg.Shao = 0
//				msg.Peng = self.PersonMgr[i].Peng
//				msg.Chi = self.PersonMgr[i].Chi
//				msg.Zhao = self.PersonMgr[i].Zhao
//				self.room.SendMsg(self.PersonMgr[i].Uid, "gamehfbhfan", &msg)
//			}
//			return
//		}

//		if self.ChiCard(person, _type) {
//			if person.Zhua > 0 { //! 可选择下抓
//				person.Step = 3
//			} else { //! 只能打牌
//				person.Step = 1
//			}
//			self.SendOperator(person)
//		}

//		return
//	}

//	if person.Peng == 1 || person.Zhao == 1 {
//		person.NoZOP = append(person.NoZOP, self.Card)
//	}

//	//! 其他人有碰或者招
//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Uid == uid {
//			continue
//		}
//		if self.PersonMgr[i].Peng == 1 || self.PersonMgr[i].Zhao == 1 {
//			person.Reset()
//			person.Chi = _type
//			return
//		}
//	}

//	//! 其他人已经点了碰
//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Uid == uid {
//			continue
//		}
//		if self.PersonMgr[i].Peng > 1 {
//			self.PengCard(self.PersonMgr[i])
//			if self.PersonMgr[i].Zhua > 0 { //! 可选择下抓
//				self.PersonMgr[i].Step = 3
//			} else { //! 只能打牌
//				self.PersonMgr[i].Step = 1
//			}
//			self.SendOperator(self.PersonMgr[i])
//			return
//		}
//	}

//	if self.CurStep != uid { //! 不是我的局判断一下其他人是否能吃
//		for i := 0; i < len(self.PersonMgr); i++ {
//			if self.PersonMgr[i].Uid == uid {
//				continue
//			}
//			if self.PersonMgr[i].Chi == 1 {
//				person.Reset()
//				person.Chi = _type
//				return
//			} else if self.PersonMgr[i].Chi > 1 {
//				if self.ChiCard(self.PersonMgr[i], self.PersonMgr[i].Chi) {
//					if self.PersonMgr[i].Zhua > 0 { //! 可选择下抓
//						self.PersonMgr[i].Step = 3
//					} else { //! 只能打牌
//						self.PersonMgr[i].Step = 1
//					}
//					self.SendOperator(self.PersonMgr[i])
//				}
//				return
//			}
//		}
//	}

//	if self.ChiCard(person, _type) {
//		if person.Zhua > 0 { //! 可选择下抓
//			person.Step = 3
//		} else { //! 只能打牌
//			person.Step = 1
//		}
//		self.SendOperator(person)
//	}
//}

////! 拖
//func (self *Game_HFBH) GameTuo(uid int64, card int, index int) {
//	if !self.room.Begin {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始:GameTuo")
//		return
//	}

//	person := self.GetPerson(uid)
//	if person == nil {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person:GameTuo")
//		return
//	}

//	if index < 0 || index >= len(person.Card.Chi2) {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "index错误:GameTuo")
//		return
//	}

//	if person.Card.TuoCard(card, index) {
//		if person.Peng == 1 || person.Zhao == 1 || person.Chi == 1 || person.Shao == 1 {
//			if person.Peng == 1 && person.Card.IsCanPeng(self.Card, self.room.Param1/10%10 == 1, self.GetNoPengOrZhao(person)) != 2 {
//				person.Peng = 0
//			}
//			if person.Zhao == 1 && !person.Card.IsCanZhao(self.Card, self.GetNoPengOrZhao(person)) {
//				person.Zhao = 0
//			}
//			if person.Chi == 1 {
//				if !person.Card.IsCanChi(self.Card, self.GetNoChi(person)) {
//					person.Chi = 0
//				}
//			} else if len(person.Card.Chi2) < 2 {
//				if (self.CurStep == person.Uid || self.GetNextUid(self.CurStep) == person.Uid) && person.Card.IsCanChi(self.Card, self.GetNoChi(person)) {
//					person.Chi = 1
//				}
//			}
//			if person.Shao == 1 && !person.Card.IsCanShao(self.Card) {
//				person.Shao = 0
//			}
//		}

//		person.Dian = person.Card.CountDian(self.Jiang)

//		var msg Msg_GameHFBH_Tuo
//		msg.Uid = uid
//		msg.Index = index
//		msg.Card = card
//		msg.Chi2 = person.Card.Chi2
//		msg.Chi3 = person.Card.Chi3
//		msg.Dian = person.Dian
//		msg.Sequence = person.Card.Sequence
//		self.room.broadCastMsg("gamehfbhtuo", &msg)
//	}
//}

////! 抓
//func (self *Game_HFBH) GameZhua(uid int64) {
//	if !self.room.Begin {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始:GameZhua")
//		return
//	}

//	person := self.GetPerson(uid)
//	if person == nil {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person:GameZhua")
//		return
//	}

//	if person.Step == 1 {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能抓:GameZhua")
//		return
//	}

//	if person.Uid != self.CurStep {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能抓:GameZhua")
//		return
//	}

//	person.Step = 0
//	if person.Zhua > 0 {
//		person.ZhuaNum++
//		person.Zhua--
//	}

//	var msg Msg_GameHFBH_Zhua
//	msg.Uid = uid
//	msg.ZhuaNum = person.ZhuaNum
//	self.room.broadCastMsg("gamehfbhzhua", &msg)

//	self.GameMoCard(self.GetNextUid(self.CurStep))
//}

////! 胡
//func (self *Game_HFBH) GameHu(uid int64) {
//	if !self.room.Begin {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始:GameHu")
//		return
//	}

//	person := self.GetPerson(uid)
//	if person == nil {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person:GameHu")
//		return
//	}

//	if person.Hu != 1 {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能胡:GameHu")
//		return
//	}

//	self.AddStep(uid, 99, []int{self.Card})
//	person.Card = person.HuCard
//	self.Winer = uid
//	self.OnEnd()
//}

////! 过
//func (self *Game_HFBH) GameGuo(uid int64) {
//	if !self.room.Begin {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始:GameGuo")
//		return
//	}

//	person := self.GetPerson(uid)
//	if person == nil {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person:GameGuo")
//		return
//	}

//	if person.Peng == 0 && person.Zhao == 0 && person.Chi == 0 && person.Shao == 0 {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能过:GameGuo")
//		return
//	}

//	if person.Shao == 1 { //! 可以绍点了过
//		person.Peng = 0
//		person.Zhao = 0
//		person.Chi = 0
//		person.Shao = 0
//		//! 先判断胡
//		find := false
//		nextperson := self.GetNextPerson(uid)
//		for nextperson != person {
//			hu, peng, zhao, chi := self.IsHu(nextperson)
//			if hu {
//				nextperson.Card.Mo = self.Card
//				lst := lib.GetHFBHMgr().IsHu(nextperson.Card, false, peng, zhao, chi, self.Jiang)
//				if len(lst) > 0 {
//					nextperson.HuCard = lst[0]
//					nextperson.Hu = 1
//					find = true
//					break
//				}
//			}
//			nextperson = self.GetNextPerson(nextperson.Uid)
//		}
//		if find {
//			for i := 0; i < len(self.PersonMgr); i++ {
//				var msg Msg_GameHFBH_Draw
//				msg.IsPlay = false
//				msg.Uid = uid
//				msg.State = 0
//				msg.Card = self.Card
//				msg.Hu = self.PersonMgr[i].Hu
//				msg.Shao = 0
//				msg.Peng = 0
//				msg.Chi = 0
//				msg.Zhao = 0
//				self.room.SendMsg(self.PersonMgr[i].Uid, "gamehfbhfan", &msg)
//			}
//			return
//		}

//		//! 能跳就直接跳
//		for i := 0; i < len(self.PersonMgr); i++ {
//			if self.PersonMgr[i].Uid == uid {
//				continue
//			}
//			if self.PersonMgr[i].Card.IsCanTiao(self.Card) {
//				self.PersonMgr[i].Card.TiaoCard(self.Card)
//				self.CurStep = self.PersonMgr[i].Uid
//				var msg Msg_GameHFBH_Draw
//				msg.IsPlay = false
//				msg.Uid = self.PersonMgr[i].Uid
//				msg.State = self.PersonMgr[i].Uid
//				msg.Card = self.Card
//				msg.Hu = 0
//				msg.Shao = 0
//				msg.Peng = 0
//				msg.Chi = 0
//				msg.Zhao = 0
//				msg.Zhao4 = self.PersonMgr[i].Card.Zhao4
//				msg.Tiao = self.PersonMgr[i].Card.Tiao
//				msg.Sequence = self.PersonMgr[i].Card.Sequence
//				msg.Dian = self.PersonMgr[i].Card.CountDian(self.Jiang)
//				self.room.broadCastMsg("gamehfbhfan", &msg)

//				self.GameMoCard(self.GetNextUid(self.CurStep))
//				return
//			}
//		}

//		//! 不能胡不能绍,轮流判断
//		//! 再判断其他情况
//		for i := 0; i < len(self.PersonMgr); i++ {
//			if self.PersonMgr[i].Uid == uid { //! 是自己判断是否能吃
//				continue
//			} else { //!
//				if self.PersonMgr[i].Card.IsCanPeng(self.Card, self.room.Param1/10%10 == 1, self.GetNoPengOrZhao(self.PersonMgr[i])) == 2 {
//					self.PersonMgr[i].Peng = 1
//					find = true
//				}
//				if self.PersonMgr[i].Card.IsCanZhao(self.Card, self.GetNoPengOrZhao(self.PersonMgr[i])) {
//					self.PersonMgr[i].Zhao = 1
//					find = true
//				}
//				if self.PersonMgr[i].Uid == self.GetNextUid(self.CurStep) && self.PersonMgr[i].Card.IsCanChi(self.Card, self.GetNoChi(self.PersonMgr[i])) {
//					self.PersonMgr[i].Chi = 1
//					find = true
//				}
//			}
//		}
//		if find {
//			for i := 0; i < len(self.PersonMgr); i++ {
//				var msg Msg_GameHFBH_Draw
//				msg.IsPlay = false
//				msg.Uid = uid
//				msg.State = 0
//				msg.Card = self.Card
//				msg.Hu = 0
//				msg.Shao = 0
//				msg.Peng = self.PersonMgr[i].Peng
//				msg.Chi = self.PersonMgr[i].Chi
//				msg.Zhao = self.PersonMgr[i].Zhao
//				self.room.SendMsg(self.PersonMgr[i].Uid, "gamehfbhfan", &msg)
//			}
//			return
//		}

//		person.Card.Card2 = append(person.Card.Card2, self.Card)
//		var msg Msg_GameHFBH_Draw
//		msg.IsPlay = false
//		msg.Uid = uid
//		msg.State = -1
//		msg.Card = self.Card
//		msg.Hu = 0
//		msg.Shao = 0
//		msg.Peng = 0
//		msg.Chi = 0
//		msg.Zhao = 0
//		self.room.broadCastMsg("gamehfbhfan", &msg)

//		self.GameMoCard(self.GetNextUid(self.CurStep))
//		return
//	}

//	if person.Peng == 1 || person.Zhao == 1 {
//		person.NoZOP = append(person.NoZOP, self.Card)
//	}

//	person.Peng = 0
//	person.Zhao = 0
//	person.Chi = 0
//	person.Shao = 0

//	//! 看其他人有没有碰
//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Peng == 1 {
//			return
//		} else if self.PersonMgr[i].Peng == 2 {
//			self.PengCard(self.PersonMgr[i])
//			if self.PersonMgr[i].Zhua > 0 { //! 可选择下抓
//				self.PersonMgr[i].Step = 3
//			} else { //! 只能打牌
//				self.PersonMgr[i].Step = 1
//			}
//			self.SendOperator(self.PersonMgr[i])
//			return
//		}
//	}

//	//! 看其他人有没有吃
//	lstchi := make([]int, 0)
//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Chi == 1 {
//			return
//		} else if self.PersonMgr[i].Chi > 1 {
//			lstchi = append(lstchi, i)
//		}
//	}
//	if len(lstchi) == 1 { //! 只有一个人决定吃
//		self.ChiCard(self.PersonMgr[lstchi[0]], self.PersonMgr[lstchi[0]].Chi)
//		if self.PersonMgr[lstchi[0]].Zhua > 0 { //! 可选择下抓
//			self.PersonMgr[lstchi[0]].Step = 3
//		} else { //! 只能打牌
//			self.PersonMgr[lstchi[0]].Step = 1
//		}
//		self.SendOperator(self.PersonMgr[lstchi[0]])
//		return
//	} else {
//		for i := 0; i < len(lstchi); i++ {
//			if self.PersonMgr[lstchi[i]].Uid == self.CurStep {
//				self.ChiCard(self.PersonMgr[lstchi[i]], self.PersonMgr[lstchi[i]].Chi)
//				if self.PersonMgr[lstchi[i]].Zhua > 0 { //! 可选择下抓
//					self.PersonMgr[lstchi[i]].Step = 3
//				} else { //! 只能打牌
//					self.PersonMgr[lstchi[i]].Step = 1
//				}
//				self.SendOperator(self.PersonMgr[lstchi[i]])
//				return
//			}
//		}
//	}

//	//! 都没有操作,则进入下一轮
//	curperson := self.GetPerson(self.CurStep)
//	curperson.Card.Card2 = append(curperson.Card.Card2, self.Card)
//	self.GameMoCard(self.GetNextUid(self.CurStep))
//}

////! 碰牌
//func (self *Game_HFBH) PengCard(person *Game_HFBH_Person) bool {
//	if !person.Card.PengCard(self.Card) {
//		return false
//	}

//	person.NoStep = append(person.NoStep, self.Card)
//	self.Reset()
//	person.Dian = person.Card.CountDian(self.Jiang)
//	self.CurStep = person.Uid
//	self.AddStep(person.Uid, 23, []int{self.Card})

//	var msg Msg_GameHFBH_Peng
//	msg.Uid = person.Uid
//	msg.Card = self.Card
//	msg.Peng = person.Card.Peng
//	msg.Dian = person.Dian
//	msg.Sequence = person.Card.Sequence
//	msg.NoStep = person.NoStep
//	self.room.broadCastMsg("gamehfbhpeng", &msg)

//	self.Card = 0

//	return true
//}

////! 吃牌
//func (self *Game_HFBH) ChiCard(person *Game_HFBH_Person, _type int) bool {
//	if !person.Card.ChiCard(self.Card, _type) {
//		return false
//	}

//	person.Chis = append(person.Chis, self.Card)
//	person.NoStep = append(person.NoStep, self.Card)
//	self.Reset()
//	person.Dian = person.Card.CountDian(self.Jiang)
//	self.CurStep = person.Uid
//	self.AddStep(person.Uid, -_type, []int{self.Card})

//	var msg Msg_GameHFBH_Chi
//	msg.Uid = person.Uid
//	msg.Card = self.Card
//	msg.Chi2 = person.Card.Chi2
//	msg.Chi3 = person.Card.Chi3
//	msg.Type = _type
//	msg.Dian = person.Dian
//	msg.Sequence = person.Card.Sequence
//	msg.NoStep = person.NoStep
//	self.room.broadCastMsg("gamehfbhchi", &msg)

//	self.Card = 0
//	return true
//}

//func (self *Game_HFBH) getInfo(uid int64) *Msg_GameHFBH_Info {
//	var msg Msg_GameHFBH_Info
//	msg.Begin = self.room.Begin
//	msg.CurStep = self.CurStep
//	person := self.GetPerson(self.CurStep)
//	if person != nil && (person.Hu == 1 || person.Shao == 1) && person.Uid != uid {
//		msg.Card = -1
//	} else {
//		msg.Card = self.Card
//	}
//	msg.Jiang = self.Jiang
//	msg.IsPlay = self.IsPlay
//	if self.Mgr == nil {
//		msg.Num = 0
//	} else {
//		msg.Num = len(self.Mgr.Card)
//	}
//	if self.room.Begin || self.Mgr == nil {
//		msg.Back = make([]int, 0)
//	} else {
//		msg.Back = self.Mgr.Card
//	}

//	for i := 0; i < len(self.PersonMgr); i++ {
//		var son Son_GameHFBH_Info
//		son.Uid = self.PersonMgr[i].Uid
//		son.Deal = self.PersonMgr[i].Deal
//		son.Card = self.PersonMgr[i].Card
//		son.Dian = self.PersonMgr[i].Dian
//		son.ZhuaNum = self.PersonMgr[i].ZhuaNum
//		son.Step = self.PersonMgr[i].Step
//		son.NoStep = self.PersonMgr[i].NoStep
//		son.Hu = self.PersonMgr[i].Hu
//		if son.Uid != uid && msg.Begin {
//			son.Hu = 0
//			son.Dian = 0
//			son.Card.Card1 = make([]int, 0)
//			son.Card.Shao3 = make([]int, len(self.PersonMgr[i].Card.Shao3))
//			son.Card.Shao4 = make([]int, len(self.PersonMgr[i].Card.Shao4))
//			son.Card.Shao5 = make([]int, len(self.PersonMgr[i].Card.Shao5))
//			son.NoStep = make([]int, 0)
//		}
//		son.Score = self.PersonMgr[i].Score
//		son.Total = self.PersonMgr[i].Total
//		son.Ready = self.PersonMgr[i].Ready
//		son.Peng = self.PersonMgr[i].Peng
//		son.Chi = self.PersonMgr[i].Chi
//		son.Shao = self.PersonMgr[i].Shao
//		son.Zhao = self.PersonMgr[i].Zhao
//		msg.Info = append(msg.Info, son)
//	}
//	return &msg
//}

//func (self *Game_HFBH) SendOperator(person *Game_HFBH_Person) {
//	if !person.IsCanPlay() { //! 不能打,直接下抓
//		person.Step = 2
//	}

//	var msg Msg_GameHFBH_Operator
//	msg.Step = person.Step
//	self.room.SendMsg(person.Uid, "gamehfbhoperator", &msg)
//}

//func (self *Game_HFBH) IsHu(person *Game_HFBH_Person) (bool, bool, bool, bool) {
//	if self.Card == 0 {
//		return false, false, false, false
//	}

//	peng := true
//	zhao := true

//	lst := self.GetNoPengOrZhao(person)
//	for i := 0; i < len(lst); i++ {
//		if lst[i] == self.Card {
//			peng = false
//			zhao = false
//			break
//		}
//	}

//	lst = self.GetNoChi(person)
//	for i := 0; i < len(lst); i++ {
//		if lst[i] == self.Card {
//			return false, false, false, false
//		}
//	}

//	return true, peng, zhao, true
//}

//func (self *Game_HFBH) OnTime() {

//}

//func (self *Game_HFBH) OnIsDealer(uid int64) bool {
//	return false
//}

//func (self *Game_HFBH) OnIsBets(uid int64) bool {
//	return false
//}

////! 结算所有人
//func (self *Game_HFBH) OnBalance() {
//}

////! 重置所有人操作
//func (self *Game_HFBH) Reset() {
//	for i := 0; i < len(self.PersonMgr); i++ {
//		self.PersonMgr[i].Peng = 0
//		self.PersonMgr[i].Zhao = 0
//		self.PersonMgr[i].Hu = 0
//		self.PersonMgr[i].Shao = 0
//		self.PersonMgr[i].Chi = 0
//		self.PersonMgr[i].Step = 0
//	}
//}

////! 得到一个人不能吃的牌
//func (self *Game_HFBH) GetNoChi(person *Game_HFBH_Person) []int {
//	lst := make([]int, 0)
//	//lst = append(lst, person.Steps...)
//	lst = append(lst, person.Card.Card2...)
//	_person := self.GetBeforePerson(person.Uid)
//	lst = append(lst, _person.Card.Card2...)

//	return lst
//}

////! 得到一个人不能碰的牌
//func (self *Game_HFBH) GetNoPengOrZhao(person *Game_HFBH_Person) []int {
//	return person.NoZOP
//}

////! 加入一个操作
////! 0摸牌 1出牌 13-15绍牌 23-25碰和招  负数是吃  99胡
//func (self *Game_HFBH) AddStep(uid int64, _type int, card []int) {
//	self.Record.Step = append(self.Record.Step, Son_Rec_GameHFBH_Step{uid, _type, card})
//}
