package gameserver

import (
	"lib"
	"math"
	"staticfunc"
	"time"
)

var GOLDJDQS_TIME = 3600
var GOLDJDQS_JACKBASE = 203509304
var JDQS_WILD int = 1   //! 万金油绝地求生
var JDQS_SCAT int = 2   //! 免费
var JDQS_GUNB1 int = 3  //! 步枪底
var JDQS_GUNJ1 int = 4  //! 机枪底
var JDQS_BOX int = 5    //! 箱子
var JDQS_POT int = 6    //! 锅
var JDQS_TXRA int = 7   //! 字母A
var JDQS_TXRK int = 8   //! 字母K
var JDQS_TXRQ int = 9   //! 字母Q
var JDQS_TXRJ int = 10  //! 字母J
var JDQS_GUNB2 int = 11 //! 步枪头
var JDQS_GUNJ2 int = 12 //! 机枪头

type Rec_JDQS_Info struct {
	GameType int                   `json:"gametype"`
	Time     int64                 `json:"time"` //! 记录时间
	Info     []Son_Rec_JDQS_Person `json:"info"`
}
type Son_Rec_JDQS_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Score int    `json:"score"`
	Bets  int    `json:"bets"`
}

type Game_JDQS struct {
	Person     *Game_JDQS_Person `json:"person"`
	Result     [5][3]int         `json:"result"`
	WinType    []int             `json:"wintype"` //！赢的图标数组
	Bet        int               `json:"bet"`     //! 下注
	Time       int64             `json:"time"`
	JackTime   int64             `json:"jacktime"`   //奖池刷新时间
	FreeTime   int               `json:"freetime"`   //免费次数
	FreeBet    int               `json:"freebet"`    //
	TotalMul   int               `json:"totalmul"`   //总倍率
	GameLevel  int               `json:"gamelevel"`  //游戏难度1~5  3为奖池
	PersonSet  int               `json:"personset"`  //个人设置1启用  0禁用
	StFreeRun  int               `json:"stfreerun"`  //设置免费
	PL         map[int][]int     `json:"pl"`         //! 赔率
	EnemyTotal int               `json:"enemytotal"` //! 敌人总数
	GetChicken bool              `json:"getchicken"` //! 是否吃鸡
	room       *Room
}
type Game_JDQS_Person struct {
	Uid     int64  `json:"uid"`
	Gold    int    `json:"gold"`  //! 进房时的金币
	Total   int    `json:"total"` //! 金币数
	Win     int    `json:"win"`   //! 赢了多少钱
	Kill    int    `json:"kill"`  //! 杀人数量
	Cost    int    `json:"cost"`  //!　抽水
	Name    string `json:"name"`  //! 名字
	Head    string `json:"head"`  //! 头像
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}
type Msg_GameJDQS_ChickenGold struct {
	Uid   int64 `json:"uid"`
	Gold  int   `json:"gold"`
	Total int   `json:"total"`
}

func NewGame_JDQS() *Game_JDQS {
	game := new(Game_JDQS)
	game.GameLevel = 3
	game.PersonSet = 0 //默认禁用
	game.StFreeRun = 0
	game.WinType = make([]int, 0)
	game.EnemyTotal = 100
	//! 初始化赔率345
	game.PL = make(map[int][]int, 0)
	game.PL[JDQS_GUNB1] = []int{10, 25, 100} //! 步枪
	game.PL[JDQS_GUNJ1] = []int{10, 25, 80}  //! 机枪
	game.PL[JDQS_BOX] = []int{5, 15, 60}     //! 箱子
	game.PL[JDQS_POT] = []int{5, 15, 50}     //! 锅
	game.PL[JDQS_TXRA] = []int{3, 10, 25}    //! A
	game.PL[JDQS_TXRK] = []int{3, 10, 25}    //! K
	game.PL[JDQS_TXRQ] = []int{2, 5, 15}     //! Q
	game.PL[JDQS_TXRJ] = []int{2, 5, 15}     //! J
	return game
}

type Msg_GameJDQS_Info struct {
	Begin      bool              `json:"begin"`
	Result     [5][3]int         `json:"result"`
	FreeTime   int               `json:"freetime"` //免费次数
	JackPot    int64             `json:"jackpot"`
	Money      int               `json:"money"`      //! 底分
	IsAdmin    bool              `json:"isadmin"`    //! 是否是超端
	EnemyTotal int               `json:"enemytotal"` //! 敌人总数
	Person     Son_GameJDQS_Info `json:"person"`
}

type Son_GameJDQS_Info struct {
	Uid     int64  `json:"uid"`
	Total   int    `json:"total"` //! 金币数
	Coin    int    `json:"coin"`  //! 下注硬币数量
	Kill    int    `json:"kill"`  //! 杀人数量
	Name    string `json:"name"`  //! 名字
	Head    string `json:"head"`  //! 头像
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type Msg_GameJDQS_Total struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"`
}
type Msg_GameJDQS_JackPot struct {
	Uid     int64 `json:"uid"`
	JackPot int64 `json:"jackpot"`
}
type Msg_GameJDQS_ResultArr struct {
	Result [5][3]int `json:"result"`
}
type Msg_GameJDQS_End struct {
	Uid        int64     `json:"uid"`
	Win        int       `json:"win"`     //总赢钱
	WinType    []int     `json:"wintype"` //！赢的图标数组
	Result     [5][3]int `json:"result"`
	FreeNum    int       `json:"freenum"`    //! 免费次数
	Total      int       `json:"total"`      //总金币
	TotalMul   int       `json:"totalmul"`   //总倍率
	Kill       int       `json:"kill"`       //! 杀人数量
	CurKill    int       `json:"curkill"`    //! 本局杀人
	CurDead    int       `json:"curdead"`    //! 本局死亡人数
	GetChicken bool      `json:"getchicken"` //! 是否吃鸡
}
type Msg_GameJDQS_STinfo struct {
	Uid         int64 `json:"uid"`
	JackPot     int64 `json:"jackpot"`     //！奖池
	JackPotMax  int64 `json:"jackpotmax"`  //！最大奖池
	JackPotMin  int64 `json:"jackpotmin"`  //! 最小奖池
	GameLevel   int   `json:"gamelevel"`   //! 难度级别
	PersonalSet int   `json:"personalset"` //是否禁用个人配置
}

type Msg_GameJDQS_RoomInfo struct {
	Info []Son_GameJDQS_RoomInfo `json:"info"`
}
type Son_GameJDQS_RoomInfo struct {
	Id        int    `json:"id"`        //! 房间号
	Uid       int64  `json:"uid"`       //! 房间里的人
	Name      string `json:"name"`      //! 名字
	LiveTime  int64  `json:"livetime"`  //! 活动时间
	Total     int    `json:"total"`     //！当前总金币
	Win       int    `json:"win"`       //！当前总输赢
	GameLevel int    `json:"gamelevel"` //！游戏难度
	PersonSet int    `json:"personset"` //！个人设置
}

//! 同步金币
func (self *Game_JDQS_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}
func (self *Game_JDQS) OnMsg(msg *RoomMsg) {
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "OnMsg", msg.Head)
	switch msg.Head {
	case "gamebets": //
		self.GameStart(msg.Uid, msg.V.(*Msg_GameBets).Bets)
	case "gamestart": //吃鸡
		self.GameChicken(msg.Uid)
	case "gameview": //超端获取信息
		self.GameGetST(msg.Uid)
	case "gamedwsetjack": //设置电玩奖池
		self.GameSetJackPot(msg.Uid, msg.V.(*Msg_DwSetJackPot).JackPot, msg.V.(*Msg_DwSetJackPot).JackPotMax, msg.V.(*Msg_DwSetJackPot).JackPotMin)
	case "gamedwsetpro": //设置电玩属性
		self.GameSetProperty(msg.Uid, msg.V.(*Msg_DwSetPro).GameLevel, msg.V.(*Msg_DwSetPro).PersonalSet)
	case "gamedwroomlist": //超端获取房间列表
		self.GameGetRoomList(msg.Uid)
	case "gamedwsetroom": //设置房间信息
		self.GameGetRoomDetail(msg.Uid, msg.V.(*Msg_DwSetRoom).RoomId, msg.V.(*Msg_DwSetRoom).GameLevel, msg.V.(*Msg_DwSetRoom).PersonalSet, msg.V.(*Msg_DwSetRoom).FreeRun)
	case "gamegetroomset": //设置房间信息
		self.GameSetRoomPro(msg.V.(*Msg_DwSetRoom).RoomId, msg.V.(*Msg_DwSetRoom).GameLevel, msg.V.(*Msg_DwSetRoom).PersonalSet, msg.V.(*Msg_DwSetRoom).FreeRun)
	}
}
func (self *Game_JDQS) GameSetRoomPro(roomid int, gamelevel int, personset int, freerun int) {
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "非超端玩家GameSetRoomPro", self.room.Id)
	self.GameLevel = gamelevel
	self.PersonSet = personset
	if 0 == personset {
		return
	}
	if 3 == freerun {
		self.StFreeRun = 3
	} else if 5 == freerun {
		self.StFreeRun = 4
	} else if 10 == freerun {
		self.StFreeRun = 5
	}
}
func (self *Game_JDQS) GameGetRoomDetail(uid int64, roomid int, gamelevel int, personset int, freerun int) {
	if !GetServer().IsAdmin(uid, staticfunc.ADMIN_DWJDQS) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "非超端玩家")
		return
	}
	GetRoomMgr().lock.RLock()
	for id, value := range GetRoomMgr().MapRoom {
		//lib.GetLogMgr().Output(lib.LOG_DEBUG, "GetRoomMgr", id, value.Uid, value.Type, self.room.Type)
		if id == roomid {
			var msg Msg_DwSetRoom
			msg.RoomId = roomid
			msg.GameLevel = gamelevel
			msg.PersonalSet = personset
			msg.FreeRun = freerun
			value.Operator(NewRoomMsg("gamegetroomset", "", uid, &msg))

			GetServer().SqlSuperClientLog(&SQL_SuperClientLog{1, uid, self.room.Type, lib.HF_JtoA(&msg), time.Now().Unix()})

			self.room.SendMsg(uid, "ok", nil)
			break
		}
	}
	GetRoomMgr().lock.RUnlock()
}
func (self *Game_JDQS) GameGetRoomList(uid int64) {
	if !GetServer().IsAdmin(uid, staticfunc.ADMIN_DWJDQS) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "非超端玩家")
		return
	}
	var msg Msg_GameJDQS_RoomInfo
	msg.Info = make([]Son_GameJDQS_RoomInfo, 0)
	GetRoomMgr().lock.RLock()
	for id, value := range GetRoomMgr().MapRoom {
		if value.Type == self.room.Type {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "GetRoomMgr", id, value.Viewer, value.game.(*Game_JDQS).GameLevel)
			var son Son_GameJDQS_RoomInfo
			son.Id = value.Id
			son.Uid = value.Viewer[0]
			son.Name = value.HostName
			son.LiveTime = value.LiveTime
			son.Total = value.game.(*Game_JDQS).Person.Total
			son.Win = value.game.(*Game_JDQS).Person.Total - value.game.(*Game_JDQS).Person.Gold
			son.GameLevel = value.game.(*Game_JDQS).GameLevel
			son.PersonSet = value.game.(*Game_JDQS).PersonSet
			msg.Info = append(msg.Info, son)
		}
	}
	GetRoomMgr().lock.RUnlock()
	self.room.SendMsg(uid, "gamejdqsstroom", &msg)
}
func (self *Game_JDQS) GameSetJackPot(uid int64, jackPot int64, jackPotMax int64, jackPotMin int64) {
	if !GetServer().IsAdmin(uid, staticfunc.ADMIN_DWJDQS) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "非超端玩家")
		return
	}
	lib.GetLogMgr().Output(lib.LOG_ERROR, "超端设置奖池!", uid, jackPot, jackPotMax, jackPotMin)
	GetServer().SetDwJdqsSysMoney(self.room.Type%10000, jackPot)
	lib.GetManyMgr().SetDWProperty(self.room.Type, jackPotMax, jackPotMin)

	var msg Msg_DwSetJackPot
	msg.JackPot = jackPot
	msg.JackPotMax = jackPotMax
	msg.JackPotMin = jackPotMin

	GetServer().SqlSuperClientLog(&SQL_SuperClientLog{1, uid, self.room.Type, lib.HF_JtoA(&msg), time.Now().Unix()})
	self.room.SendMsg(uid, "ok", nil)
}
func (self *Game_JDQS) GameSetProperty(uid int64, gameLevel int, personalset int) {
	if !GetServer().IsAdmin(uid, staticfunc.ADMIN_DWJDQS) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "非超端玩家")
		return
	}
	lib.GetLogMgr().Output(lib.LOG_ERROR, "超端设置游戏!", uid, gameLevel, personalset)
	value := lib.GetSingleMgr().GetProperty(self.room.Type)
	value.GameLevel = gameLevel
	value.PersonalSet = personalset
	lib.GetSingleMgr().SetProperty(self.room.Type, value)

	var msg Msg_DwSetPro
	msg.GameLevel = gameLevel
	msg.PersonalSet = personalset
	GetServer().SqlSuperClientLog(&SQL_SuperClientLog{1, uid, self.room.Type, lib.HF_JtoA(&msg), time.Now().Unix()})

	self.room.SendMsg(uid, "ok", nil)
}
func (self *Game_JDQS) GameGetST(uid int64) {
	if !GetServer().IsAdmin(uid, staticfunc.ADMIN_DWJDQS) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "非超端玩家")
		return
	}
	var msg Msg_GameJDQS_STinfo
	msg.Uid = uid
	msg.JackPot = GetServer().DwJdqsSysMoney[self.room.Type%10000]
	msg.JackPotMax = lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax
	msg.JackPotMin = lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin
	msg.GameLevel = lib.GetSingleMgr().GetProperty(self.room.Type).GameLevel
	msg.PersonalSet = lib.GetSingleMgr().GetProperty(self.room.Type).PersonalSet
	self.room.SendMsg(uid, "gamejdqsstinfo", &msg)
}
func (self *Game_JDQS) OnBalance() {
	if self.Person == nil {
		return
	}

	gold := self.Person.Total - self.Person.Gold
	if gold > 0 {
		GetRoomMgr().AddCard(self.Person.Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
	} else if gold < 0 {
		GetRoomMgr().CostCard(self.Person.Uid, staticfunc.TYPE_GOLD, -gold, self.room)
	}
	self.Person.Gold = self.Person.Total
}
func (self *Game_JDQS) GameChicken(uid int64) {
	if !self.GetChicken {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "当前未吃鸡")
		return
	}
	//df := staticfunc.GetCsvMgr().GetDF(self.room.Type)
	//temp := df * (lib.HF_GetRandom(9) + 1)
	self.GetChicken = false
	self.EnemyTotal = 100
	var msg Msg_GameJDQS_ChickenGold
	msg.Uid = self.Person.Uid
	//msg.Total
	self.room.SendMsg(uid, "gamechickengold", &msg)
}
func (self *Game_JDQS) GameStart(uid int64, bets int) {
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "GameStart!", uid, bets)
	if self.room.Begin {
		person := GetPersonMgr().GetPerson(uid)
		person.SendErr("游戏已开始")
		lib.GetLogMgr().Output(lib.LOG_ERROR, "游戏已开始!")
		return
	}
	if uid != self.Person.Uid {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "GameStart person.uid != uid")
		return
	}
	if bets > self.Person.Total && self.FreeTime <= 0 {
		person := GetPersonMgr().GetPerson(uid)
		person.SendErr("金币不足！")
		lib.GetLogMgr().Output(lib.LOG_ERROR, "金币不足!")
		return
	}
	if bets <= 0 && self.FreeTime <= 0 {
		person := GetPersonMgr().GetPerson(uid)
		person.SendErr("下注金额错误！")
		lib.GetLogMgr().Output(lib.LOG_ERROR, "下注金额错误！!")
		return
	}

	if self.FreeTime > 0 {
		self.Bet = self.FreeBet
	} else {
		self.Bet = bets
		self.Person.Total -= bets
	}
	self.SendTotal(self.Person.Uid, self.Person.Total)

	self.OnBegin()
}
func (self *Game_JDQS) SendTotal(uid int64, total int) {
	var msg Msg_GameJDQS_Total
	msg.Uid = uid
	msg.Total = total
	self.room.SendMsg(uid, "gamegoldtotal", &msg)
}
func (self *Game_JDQS) OnBegin() {
	self.room.Begin = true

	winLst := make([]Msg_GameJDQS_ResultArr, 0)
	lostLst := make([]Msg_GameJDQS_ResultArr, 0)
	levelLst1 := make([]Msg_GameJDQS_ResultArr, 0) //送分
	levelLst2 := make([]Msg_GameJDQS_ResultArr, 0) //少量送
	levelLst4 := make([]Msg_GameJDQS_ResultArr, 0) //少量杀
	levelLst5 := make([]Msg_GameJDQS_ResultArr, 0) //杀分
	for i := 0; i < 150; i++ {
		self.GetRandomResult()
		score, _ := self.GetAllScore(false)
		mywin := 0
		if self.FreeTime > 0 {
			mywin = score
		} else {
			mywin = score - self.Bet
		}
		if score <= 0 {
			var msg Msg_GameJDQS_ResultArr
			lib.HF_DeepCopy(&msg.Result, &self.Result)
			levelLst5 = append(levelLst5, msg) //杀分
		} else if score > 0 && mywin <= 0 {
			var msg Msg_GameJDQS_ResultArr
			lib.HF_DeepCopy(&msg.Result, &self.Result)
			levelLst4 = append(levelLst4, msg) //少量杀分
		} else if mywin > 0 && mywin < self.Bet {
			var msg Msg_GameJDQS_ResultArr
			lib.HF_DeepCopy(&msg.Result, &self.Result)
			levelLst2 = append(levelLst2, msg) // 少量送分
		} else if mywin >= self.Bet {
			var msg Msg_GameJDQS_ResultArr
			lib.HF_DeepCopy(&msg.Result, &self.Result)
			levelLst1 = append(levelLst1, msg) // 送分
		}
		//! 输分列表
		if score <= 0 || mywin < 0 {
			var msg Msg_GameJDQS_ResultArr
			lib.HF_DeepCopy(&msg.Result, &self.Result)
			lostLst = append(lostLst, msg)
			continue
		}
		//奖池允许列表
		if GetServer().DwJdqsSysMoney[self.room.Type%10000]-int64(mywin) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && GetServer().DwJdqsSysMoney[self.room.Type%10000]-int64(mywin) <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
			var msg Msg_GameJDQS_ResultArr
			lib.HF_DeepCopy(&msg.Result, &self.Result)
			winLst = append(winLst, msg)
		}
	}
	adminState := self.GetScoreState()
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "获取当前模式", adminState)
	if 1 == adminState { //送分
		if len(levelLst1) > 0 {
			lib.HF_DeepCopy(&self.Result, &levelLst1[lib.HF_GetRandom(len(levelLst1))].Result)
		} else { //纯随机
			self.GetRandomResult()
		}
	} else if 2 == adminState {
		if len(levelLst2) > 0 {
			lib.HF_DeepCopy(&self.Result, &levelLst2[lib.HF_GetRandom(len(levelLst2))].Result)
		} else { //纯随机
			self.GetRandomResult()
		}
	} else if 4 == adminState {
		if len(levelLst4) > 0 {
			lib.HF_DeepCopy(&self.Result, &levelLst4[lib.HF_GetRandom(len(levelLst4))].Result)
		} else {
			lib.HF_DeepCopy(&self.Result, &lostLst[lib.HF_GetRandom(len(lostLst))].Result)
		}
	} else if 5 == adminState {
		if len(levelLst5) > 0 {
			lib.HF_DeepCopy(&self.Result, &levelLst5[lib.HF_GetRandom(len(levelLst5))].Result)
		} else {
			lib.HF_DeepCopy(&self.Result, &lostLst[lib.HF_GetRandom(len(lostLst))].Result)
		}
	} else { //奖池模式
		if GetServer().DwJdqsSysMoney[self.room.Type%10000] <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && len(lostLst) > 0 {
			//!奖池不够必须输
			lib.HF_DeepCopy(&self.Result, &lostLst[lib.HF_GetRandom(len(lostLst))].Result)
		} else if len(winLst) > 0 {
			lib.HF_DeepCopy(&self.Result, &winLst[lib.HF_GetRandom(len(winLst))].Result)
		} else if len(lostLst) > 0 { //! 都找不到就纯随机输1
			lib.HF_DeepCopy(&self.Result, &lostLst[lib.HF_GetRandom(len(lostLst))].Result)
		} else { //! 都找不到就纯随机
			self.GetRandomResult()
		}
	}

	//！设置免费旋转
	scatArr := make([]int, 0)
	if self.StFreeRun > 0 {
		totalWild := 0
		for i := 0; i < len(self.Result); i++ {
			count := 0
			for j := 0; j < len(self.Result[i]); j++ { //JDQS_SCAT
				if JDQS_SCAT == self.Result[i][j] {
					count++
					totalWild++
				}
			}
			if count <= 0 {
				scatArr = append(scatArr, i)
			}
		}
		//lib.GetLogMgr().Output(lib.LOG_ERROR, "aaaaaGetScoreState", scatArr, totalWild, self.StFreeRun)
		if totalWild < self.StFreeRun {
			for i := 0; i < self.StFreeRun-totalWild; i++ {
				index1 := scatArr[i]
				index2 := lib.HF_GetRandom(3)
				//lib.GetLogMgr().Output(lib.LOG_ERROR, "sssssGetScoreState", index1, index2)
				self.Result[index1][index2] = JDQS_SCAT
			}
		}
		self.StFreeRun = 0
	}
	lib.GetLogMgr().Output(lib.LOG_ERROR, "结果", self.Result)
	self.OnEnd()
}
func (self *Game_JDQS) GetScoreState() int {
	var ret = 3
	allctl := lib.GetSingleMgr().GetProperty(self.room.Type).GameLevel
	personSet := lib.GetSingleMgr().GetProperty(self.room.Type).PersonalSet
	if 1 == personSet && 1 == self.PersonSet {
		ret = self.GameLevel
	} else { //禁用个人
		ret = allctl
	}
	//lib.GetLogMgr().Output(lib.LOG_ERROR, "获取gamelvel", allctl, personSet, ret)
	return ret
}
func (self *Game_JDQS) GetAllScore(final bool) (int, int) {
	//lib.GetLogMgr().Output(lib.LOG_DEBUG, "GetAllScore", self.Result)
	score := 0
	for j := 0; j < len(self.Result[0]); j++ {
		_score := self.GetPerScore(self.Result[0][j])
		score += _score
		if final && _score > 0 {
			var tempType = self.Result[0][j]
			if tempType > 10 {
				tempType -= 8
			}
			self.WinType = append(self.WinType, tempType)
		}
	}
	totalWild := 0
	for i := 0; i < len(self.Result); i++ {
		for j := 0; j < len(self.Result[i]); j++ {
			if JDQS_SCAT == self.Result[i][j] {
				totalWild++
			}
		}
	}
	if totalWild >= 3 {
		score += (self.Bet * 2)
	}
	return score, totalWild
}
func (self *Game_JDQS) GetPerScore(style int) int {
	var _type = style
	if _type > 10 {
		_type -= 8
	}
	if _type == JDQS_SCAT {
		return 0
	}
	total := 1
	line := 1
	score := 0
	for i := 1; i < len(self.Result); i++ {
		find := false
		for j := 0; j < len(self.Result[i]); j++ {
			var temp = self.Result[i][j]
			if temp > 10 {
				temp -= 8
			}
			if _type == temp || JDQS_WILD == temp {
				find = true
				total++
			}
		}
		if find {
			line++
		} else {
			if i <= 2 {
				return score
			} else {
				break
			}
		}
	}
	mul := total - line + 1
	//lib.GetLogMgr().Output(lib.LOG_DEBUG, "GetPerScore", style, line, mul)
	score = mul * self.PL[_type][line-3] * self.Bet / 10
	return score
}
func (self *Game_JDQS) GetRandomResult() {
	for j := 0; j < len(self.Result); j++ {
		findfree := false
		for k := 0; k < len(self.Result[j]); k++ {
			if findfree || 0 == j { //每列只有一个万金油，第一列没有万金油
				if 2 == k { //枪头单独只出现在第三列
					self.Result[j][k] = lib.HF_GetRandom(11) + 2
				} else {
					self.Result[j][k] = lib.HF_GetRandom(9) + 2
				}
			} else {
				if 2 == k {
					self.Result[j][k] = lib.HF_GetRandom(12) + 1
				} else {
					self.Result[j][k] = lib.HF_GetRandom(10) + 1
				}
			}
			if self.Result[j][k] == JDQS_WILD {
				findfree = true
			}
			if k > 0 && (self.Result[j][k] == JDQS_GUNB1 || self.Result[j][k] == JDQS_GUNJ1) {
				//!大于下标1的为枪尾，上一下标必为枪头
				if self.Result[j][k] == JDQS_GUNB1 {
					self.Result[j][k-1] = JDQS_GUNB2
				} else if self.Result[j][k] == JDQS_GUNJ1 {
					self.Result[j][k-1] = JDQS_GUNJ2
				}
			}
		}
		if self.Result[j][0] > 10 && ((self.Result[j][0] - 8) != self.Result[j][1]) {
			self.Result[j][0] = lib.HF_GetRandom(9) + 2
		}
	}

}
func (self *Game_JDQS) OnBye() {
}
func (self *Game_JDQS) OnExit(uid int64) {
	if uid == self.Person.Uid {
		gold := self.Person.Total - self.Person.Gold
		if gold > 0 {
			GetRoomMgr().AddCard(self.Person.Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
		} else if gold < 0 {
			GetRoomMgr().CostCard(self.Person.Uid, staticfunc.TYPE_GOLD, -gold, self.room)
		}
		self.Person.Gold = self.Person.Total
	}
}
func (self *Game_JDQS) OnEnd() {
	totalwin, totalWild := self.GetAllScore(true)
	self.Person.Win = totalwin
	playerCost := lib.GetManyMgr().GetProperty(self.room.Type).Cost  //玩家抽水
	sysCost := lib.GetManyMgr().GetProperty(self.room.Type).DealCost //系统抽水
	curDead := 0
	curKill := 0
	if totalwin > 0 {
		temp := lib.HF_GetRandom(1)
		curDead = lib.HF_GetRandom(1) + 1 //本局死亡人数
		if (self.EnemyTotal - curDead) <= 1 {
			curDead = self.EnemyTotal - 1
		}
		self.EnemyTotal -= curDead
		if 1 == temp { //自己杀的
			curKill = curDead
		}
		self.Person.Kill += curKill
		if self.EnemyTotal <= 1 {
			self.GetChicken = true
		}
	}
	//系统抽水
	sysWin := 0
	if self.FreeTime > 0 {
		sysWin = 0 - self.Person.Win
	} else {
		sysWin = self.Bet - self.Person.Win
	}
	if sysWin != 0 {
		GetServer().SqlBZWLog(&SQL_BZWLog{1, sysWin, time.Now().Unix(), self.room.Type})
	}
	if sysWin > 0 {
		cost := int(math.Ceil(float64(sysWin) * sysCost / 100.0))
		sysWin -= cost
	}
	//!更新奖池
	GetServer().SetDwJdqsSysMoney(self.room.Type%10000, GetServer().DwJdqsSysMoney[self.room.Type%10000]+int64(sysWin))
	//！玩家抽水
	playerWin := 0
	self.TotalMul = self.Person.Win / self.Bet
	if self.FreeTime > 0 {
		playerWin = self.Person.Win
	} else {
		playerWin = self.Person.Win - self.Bet
	}
	if playerWin > 0 {
		self.Person.Cost = int(math.Ceil(float64(playerWin) * playerCost / 100.0))
		GetServer().SqlAgentGoldLog(self.Person.Uid, self.Person.Cost, self.room.Type)
		GetServer().SqlAgentBillsLog(self.Person.Uid, self.Person.Cost/2, self.room.Type)
		self.Person.Win -= self.Person.Cost
	} else if playerWin < 0 {
		cost := int(math.Ceil(float64(playerWin) * playerCost / 200.0))
		GetServer().SqlAgentBillsLog(self.Person.Uid, cost, self.room.Type)
	}
	self.Person.Total += self.Person.Win

	var isFree = false
	if totalWild >= 3 && self.FreeTime <= 0 { //免费旋转中不触发免费
		//self.FreeTime = totalWild
		self.FreeBet = self.Bet //记录免费旋转的值
		self.FreeTime = 3

	} else if self.FreeTime > 0 {
		isFree = true
		self.FreeTime--
		if self.FreeTime <= 0 {
			self.FreeTime = 0
			self.FreeBet = 0
		}
	}
	var msg Msg_GameJDQS_End
	msg.WinType = make([]int, 0)
	msg.WinType = self.WinType
	msg.TotalMul = self.TotalMul
	msg.Uid = self.Person.Uid
	msg.Win = self.Person.Win
	lib.HF_DeepCopy(&msg.Result, &self.Result)
	msg.FreeNum = self.FreeTime
	msg.Total = self.Person.Total
	msg.CurKill = curKill
	msg.CurDead = curDead
	msg.Kill = self.Person.Kill
	msg.GetChicken = self.GetChicken
	self.room.SendMsg(self.Person.Uid, "gamejdqsend", &msg)

	var record Rec_JDQS_Info
	record.GameType = self.room.Type
	record.Time = time.Now().Unix()
	var rec Son_Rec_JDQS_Person
	rec.Uid = self.Person.Uid
	rec.Name = self.Person.Name
	rec.Head = self.Person.Head
	if isFree {
		rec.Score = self.Person.Win
		rec.Bets = 0
	} else {
		rec.Score = self.Person.Win - self.Bet
		rec.Bets = self.Bet
	}
	record.Info = append(record.Info, rec)
	GetServer().InsertRecord(self.room.Type, self.Person.Uid, lib.HF_JtoA(&record), rec.Score)
	////////
	self.room.Begin = false
	self.WinType = make([]int, 0)
	self.TotalMul = 0
	self.Person.Win = 0
	self.Person.Cost = 0

	self.SetTime(GOLDJDQS_TIME)
}
func (self *Game_JDQS) OnInit(room *Room) {
	self.room = room
}
func (self *Game_JDQS) OnIsBets(uid int64) bool {
	return false
}
func (self *Game_JDQS) OnIsDealer(uid int64) bool {
	return false
}
func (self *Game_JDQS) OnRobot(robot *lib.Robot) {

}
func (self *Game_JDQS) getInfo(uid int64) *Msg_GameJDQS_Info {
	var msg Msg_GameJDQS_Info
	msg.FreeTime = self.FreeTime
	msg.EnemyTotal = self.EnemyTotal
	//df := staticfunc.GetCsvMgr().GetDF(self.room.Type)
	msg.Money = staticfunc.GetCsvMgr().GetDF(self.room.Type)
	msg.JackPot = int64(GOLDJDQS_JACKBASE) + GetServer().DwJdqsSysMoney[self.room.Type%10000]
	if GetServer().IsAdmin(uid, staticfunc.ADMIN_DWJDQS) {
		msg.IsAdmin = true
	} else {
		msg.IsAdmin = false
	}

	if self.Person != nil && self.Person.Uid == uid {
		msg.Person.Uid = uid
		msg.Person.Total = self.Person.Total
		msg.Person.Name = self.Person.Name
		msg.Person.Head = self.Person.Head
		msg.Person.IP = self.Person.IP
		msg.Person.Kill = self.Person.Kill
		msg.Person.Address = self.Person.Address
		msg.Person.Sex = self.Person.Sex
	}
	return &msg
}
func (self *Game_JDQS) OnSendInfo(person *Person) {
	if self.Person != nil && self.Person.Uid == person.Uid {
		self.Person.SynchroGold(person.Gold)
		person.SendMsg("gamejdqsinfo", self.getInfo(person.Uid))
		return
	}

	_person := new(Game_JDQS_Person)
	_person.Uid = person.Uid
	_person.Gold = person.Gold
	_person.Total = person.Gold
	_person.Name = person.Name
	_person.Sex = person.Sex
	_person.IP = person.ip
	_person.Head = person.Imgurl
	_person.Address = person.minfo.Address
	self.Person = _person
	person.SendMsg("gamejdqsinfo", self.getInfo(person.Uid))

	self.SetTime(GOLDJDQS_TIME)
}
func (self *Game_JDQS) FlushJackState() {
	var msg Msg_GameJDQS_JackPot
	msg.Uid = self.Person.Uid
	msg.JackPot = int64(GOLDJDQS_JACKBASE) + GetServer().DwJdqsSysMoney[self.room.Type%10000]
	//self.room.SendMsg(self.Person.Uid, "gamejdqsjackpot", &msg)
	self.JackTime = time.Now().Unix() + 3
}

//! 设置时间
func (self *Game_JDQS) SetTime(t int) {
	if t == 0 {
		self.Time = 0
	} else {
		self.Time = time.Now().Unix() + int64(t)
	}

	var msg Msg_SetTime
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	self.room.broadCastMsg("gametime", &msg)
}
func (self *Game_JDQS) OnTime() {
	if self.Person == nil || (0 == self.Time && 0 == self.JackTime) {
		return
	}
	if time.Now().Unix() >= self.Time {
		self.room.KickViewByUid(self.Person.Uid, 96)
	}
	if time.Now().Unix() >= self.JackTime {
		self.FlushJackState()
	}
}
