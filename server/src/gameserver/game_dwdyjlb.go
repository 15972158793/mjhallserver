package gameserver

import (
	"lib"
	"math"
	"staticfunc"
	"time"
)

var GOLDDYJLB_TIME = 3600
var GOLDDYJLB_JACKBASE = 203509304
var DYJLB_TOTALLINE = 9
var DYJLB_WILD int = 1  //! 万金油
var DYJLB_SCAT int = 2  //! 免费
var DYJLB_SEVE int = 3  //! 777
var DYJLB_DEMO int = 4  //! 钻石
var DYJLB_BAR int = 5   //! BAR
var DYJLB_LING int = 6  //! 铃铛
var DYJLB_CHER int = 7  //! 樱桃
var DYJLB_MELO int = 8  //! 西瓜
var DYJLB_GRAP int = 9  //! 葡萄
var DYJLB_LEMO int = 10 //! 柠檬
var DYJLB_APPL int = 11 //! 苹果

type Rec_DYJLB_Info struct {
	GameType int                    `json:"gametype"`
	Time     int64                  `json:"time"` //! 记录时间
	Info     []Son_Rec_DYJLB_Person `json:"info"`
}
type Son_Rec_DYJLB_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Score int    `json:"score"`
	Bets  int    `json:"bets"`
}

type Game_DYJLB struct {
	Person    *Game_DYJLB_Person `json:"person"`
	Result    [5][3]int          `json:"result"`
	Bet       int                `json:"bet"` //! 每条线的花费
	Time      int64              `json:"time"`
	JackTime  int64              `json:"jacktime"`
	FreeTime  int                `json:"freetime"`  //免费次数
	FreeBet   int                `json:"freebet"`   //
	TotalMul  int                `json:"totalmul"`  //总倍率
	GameLevel int                `json:"gamelevel"` //游戏难度1~5  3为奖池
	PersonSet int                `json:"personset"` //个人设置1启用  0禁用
	StFreeRun int                `json:"stfreerun"` //
	Line      [][]int            `json:"line"`      //! 线路
	PL        map[int][]int      `json:"pl"`        //! 赔率
	room      *Room
}
type Game_DYJLB_Person struct {
	Uid     int64  `json:"uid"`
	Gold    int    `json:"gold"`  //! 进房时的金币
	Total   int    `json:"total"` //! 金币数
	Win     int    `json:"win"`   //! 赢了多少钱
	Cost    int    `json:"cost"`  //!　抽水
	Name    string `json:"name"`  //! 名字
	Head    string `json:"head"`  //! 头像
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

func NewGame_DYJLB() *Game_DYJLB {
	game := new(Game_DYJLB)
	game.GameLevel = 3
	game.PersonSet = 0 //默认禁用
	game.StFreeRun = 0
	//! 初始化线路
	game.Line = make([][]int, 0)
	game.Line = append(game.Line, []int{1, 11, 21, 31, 41}) //1
	game.Line = append(game.Line, []int{0, 10, 20, 30, 40}) //2
	game.Line = append(game.Line, []int{2, 12, 22, 32, 42}) //3
	game.Line = append(game.Line, []int{0, 11, 22, 31, 40}) //4
	game.Line = append(game.Line, []int{2, 11, 20, 31, 42}) //5
	game.Line = append(game.Line, []int{1, 10, 20, 30, 41}) //6
	game.Line = append(game.Line, []int{1, 12, 22, 32, 41}) //7
	game.Line = append(game.Line, []int{0, 10, 21, 32, 42}) //8
	game.Line = append(game.Line, []int{2, 12, 21, 30, 40}) //9
	//! 初始化赔率345
	game.PL = make(map[int][]int, 0)
	game.PL[DYJLB_SEVE] = []int{20, 100, 500} //! 777
	game.PL[DYJLB_DEMO] = []int{15, 75, 150}  //! 钻石
	game.PL[DYJLB_BAR] = []int{10, 50, 100}   //! bar
	game.PL[DYJLB_LING] = []int{8, 35, 70}    //! 铃铛
	game.PL[DYJLB_CHER] = []int{7, 30, 60}    //! 樱桃
	game.PL[DYJLB_MELO] = []int{6, 25, 50}    //! 西瓜
	game.PL[DYJLB_GRAP] = []int{5, 20, 40}    //! 葡萄
	game.PL[DYJLB_LEMO] = []int{4, 15, 30}    //! 柠檬
	game.PL[DYJLB_APPL] = []int{3, 10, 20}    //! 苹果
	return game
}

type Msg_GameDYJLB_Info struct {
	Begin    bool               `json:"begin"`
	Result   [5][3]int          `json:"result"`
	FreeTime int                `json:"freetime"` //免费次数
	JackPot  int64              `json:"jackpot"`
	Money    []int              `json:"money"`   //! 筹码
	IsAdmin  bool               `json:"isadmin"` //! 是否是超端
	Person   Son_GameDYJLB_Info `json:"person"`
}

type Son_GameDYJLB_Info struct {
	Uid     int64  `json:"uid"`
	Total   int    `json:"total"` //! 金币数
	Coin    int    `json:"coin"`  //! 下注硬币数量
	Name    string `json:"name"`  //! 名字
	Head    string `json:"head"`  //! 头像
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type Msg_GameDYJLB_Total struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"`
}
type Msg_GameDYJLB_JackPot struct {
	Uid     int64 `json:"uid"`
	JackPot int64 `json:"jackpot"`
}
type Msg_GameDYJLB_ResultArr struct {
	Result [5][3]int `json:"result"`
}
type Msg_GameDYJLB_End struct {
	Uid      int64     `json:"uid"`
	Win      int       `json:"win"` //总赢钱
	Result   [5][3]int `json:"result"`
	WinLine  []int     `json:"winline"`  //! 哪几条线得分
	FreeNum  int       `json:"freenum"`  //! 免费次数
	Total    int       `json:"total"`    //总金币
	TotalMul int       `json:"totalmul"` //总倍率
}
type Msg_GameDYJLB_STinfo struct {
	Uid         int64 `json:"uid"`
	JackPot     int64 `json:"jackpot"`     //！奖池
	JackPotMax  int64 `json:"jackpotmax"`  //！最大奖池
	JackPotMin  int64 `json:"jackpotmin"`  //! 最小奖池
	GameLevel   int   `json:"gamelevel"`   //! 难度级别
	PersonalSet int   `json:"personalset"` //是否禁用个人配置
}

type Msg_GameDYJLB_RoomInfo struct {
	Info []Son_GameDYJLB_RoomInfo `json:"info"`
}
type Son_GameDYJLB_RoomInfo struct {
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
func (self *Game_DYJLB_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}
func (self *Game_DYJLB) OnMsg(msg *RoomMsg) {
	//lib.GetLogMgr().Output(lib.LOG_DEBUG, "OnMsg", msg.Head)
	switch msg.Head {
	case "gamebets": //
		self.GameStart(msg.Uid, msg.V.(*Msg_GameBets).Bets)
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
func (self *Game_DYJLB) GameSetRoomPro(roomid int, gamelevel int, personset int, freerun int) {
	//lib.GetLogMgr().Output(lib.LOG_DEBUG, "非超端玩家GameSetRoomPro", self.room.Id)
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
func (self *Game_DYJLB) GameGetRoomDetail(uid int64, roomid int, gamelevel int, personset int, freerun int) {
	if !GetServer().IsAdmin(uid, staticfunc.ADMIN_DWDYJLB) {
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
func (self *Game_DYJLB) GameGetRoomList(uid int64) {
	if !GetServer().IsAdmin(uid, staticfunc.ADMIN_DWDYJLB) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "非超端玩家")
		return
	}
	var msg Msg_GameDYJLB_RoomInfo
	msg.Info = make([]Son_GameDYJLB_RoomInfo, 0)
	GetRoomMgr().lock.RLock()
	for id, value := range GetRoomMgr().MapRoom {
		if value.Type == self.room.Type {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "GetRoomMgr", id, value.Viewer, value.game.(*Game_DYJLB).GameLevel)
			var son Son_GameDYJLB_RoomInfo
			son.Id = value.Id
			son.Uid = value.Viewer[0]
			son.Name = value.HostName
			son.LiveTime = value.LiveTime
			son.Total = value.game.(*Game_DYJLB).Person.Total
			son.Win = value.game.(*Game_DYJLB).Person.Total - value.game.(*Game_DYJLB).Person.Gold
			son.GameLevel = value.game.(*Game_DYJLB).GameLevel
			son.PersonSet = value.game.(*Game_DYJLB).PersonSet
			msg.Info = append(msg.Info, son)
		}
	}
	GetRoomMgr().lock.RUnlock()
	self.room.SendMsg(uid, "gamedyjlbstroom", &msg)
	//lib.GetLogMgr().Output(lib.LOG_ERROR, "GameGetRoomList!", GetRoomMgr().MapRoom, len(GetRoomMgr().MapRoom))
}
func (self *Game_DYJLB) GameSetJackPot(uid int64, jackPot int64, jackPotMax int64, jackPotMin int64) {
	if !GetServer().IsAdmin(uid, staticfunc.ADMIN_DWDYJLB) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "非超端玩家")
		return
	}
	lib.GetLogMgr().Output(lib.LOG_ERROR, "超端设置奖池!", uid, jackPot, jackPotMax, jackPotMin)
	GetServer().SetDwDyjlbSysMoney(self.room.Type%10000, jackPot)
	lib.GetManyMgr().SetDWProperty(self.room.Type, jackPotMax, jackPotMin)

	var msg Msg_DwSetJackPot
	msg.JackPot = jackPot
	msg.JackPotMax = jackPotMax
	msg.JackPotMin = jackPotMin
	GetServer().SqlSuperClientLog(&SQL_SuperClientLog{1, uid, self.room.Type, lib.HF_JtoA(&msg), time.Now().Unix()})

	self.room.SendMsg(uid, "ok", nil)
}
func (self *Game_DYJLB) GameSetProperty(uid int64, gameLevel int, personalset int) {
	if !GetServer().IsAdmin(uid, staticfunc.ADMIN_DWDYJLB) {
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
func (self *Game_DYJLB) GameGetST(uid int64) {
	if !GetServer().IsAdmin(uid, staticfunc.ADMIN_DWDYJLB) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "非超端玩家")
		return
	}
	var msg Msg_GameDYJLB_STinfo
	msg.Uid = uid
	msg.JackPot = GetServer().DwDyjlbSysMoney[self.room.Type%10000]
	msg.JackPotMax = lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax
	msg.JackPotMin = lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin
	msg.GameLevel = lib.GetSingleMgr().GetProperty(self.room.Type).GameLevel
	msg.PersonalSet = lib.GetSingleMgr().GetProperty(self.room.Type).PersonalSet
	self.room.SendMsg(uid, "gamedyjlbstinfo", &msg)
}
func (self *Game_DYJLB) GameStart(uid int64, bets int) {
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
	if DYJLB_TOTALLINE*bets > self.Person.Total && self.FreeTime <= 0 {
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
		self.Person.Total -= DYJLB_TOTALLINE * bets
	}
	self.SendTotal(self.Person.Uid, self.Person.Total)

	self.OnBegin()
}
func (self *Game_DYJLB) SendTotal(uid int64, total int) {
	var msg Msg_GameDYJLB_Total
	msg.Uid = uid
	msg.Total = total
	self.room.SendMsg(uid, "gamegoldtotal", &msg)
}
func (self *Game_DYJLB) OnBalance() {
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
func (self *Game_DYJLB) OnBegin() {
	self.room.Begin = true

	winLst := make([]Msg_GameDYJLB_ResultArr, 0)
	lostLst := make([]Msg_GameDYJLB_ResultArr, 0)
	levelLst1 := make([]Msg_GameDYJLB_ResultArr, 0) //送分
	levelLst2 := make([]Msg_GameDYJLB_ResultArr, 0) //少量送
	levelLst4 := make([]Msg_GameDYJLB_ResultArr, 0) //少量杀
	levelLst5 := make([]Msg_GameDYJLB_ResultArr, 0) //杀分
	for i := 0; i < 150; i++ {
		for j := 0; j < len(self.Result); j++ {
			for k := 0; k < len(self.Result[j]); k++ {
				self.Result[j][k] = lib.HF_GetRandom(11) + 1
			}
		}
		perfree := 0 //每行免费的数量
		findfree := false
		for j := 0; j < len(self.Result); j++ {
			for k := 0; k < len(self.Result[j]); k++ {
				if self.Result[j][k] == DYJLB_SCAT {
					perfree++
				}
			}
			if perfree >= 2 {
				findfree = true
				break
			}
			perfree = 0
		}
		if findfree {
			continue
		}
		score := 0 //得分
		for j := 0; j < DYJLB_TOTALLINE; j++ {
			_score := self.GetPerLineScore(self.Line[j])
			score += _score
		}
		mywin := 0
		if self.FreeTime > 0 {
			mywin = score
		} else {
			mywin = score - self.Bet*DYJLB_TOTALLINE
		}
		if score <= 0 {
			var msg Msg_GameDYJLB_ResultArr
			lib.HF_DeepCopy(&msg.Result, &self.Result)
			levelLst5 = append(levelLst5, msg) //杀分
		} else if score > 0 && mywin <= 0 {
			var msg Msg_GameDYJLB_ResultArr
			lib.HF_DeepCopy(&msg.Result, &self.Result)
			levelLst4 = append(levelLst4, msg) //少量杀分
		} else if mywin > 0 && mywin < self.Bet*DYJLB_TOTALLINE {
			var msg Msg_GameDYJLB_ResultArr
			lib.HF_DeepCopy(&msg.Result, &self.Result)
			levelLst2 = append(levelLst2, msg) // 少量送分
		} else if mywin >= self.Bet*DYJLB_TOTALLINE {
			var msg Msg_GameDYJLB_ResultArr
			lib.HF_DeepCopy(&msg.Result, &self.Result)
			levelLst1 = append(levelLst1, msg) // 送分
		}
		//! 输分列表
		if score <= 0 || mywin < 0 {
			var msg Msg_GameDYJLB_ResultArr
			lib.HF_DeepCopy(&msg.Result, &self.Result)
			lostLst = append(lostLst, msg)
			continue
		}
		//奖池允许列表
		if GetServer().DwDyjlbSysMoney[self.room.Type%10000]-int64(mywin) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && GetServer().DwDyjlbSysMoney[self.room.Type%10000]-int64(mywin) <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
			var msg Msg_GameDYJLB_ResultArr
			lib.HF_DeepCopy(&msg.Result, &self.Result)
			winLst = append(winLst, msg)
		}
	}
	//lib.GetLogMgr().Output(lib.LOG_ERROR, "允许奖池1", levelLst1)
	//lib.GetLogMgr().Output(lib.LOG_ERROR, "允许奖池2", levelLst2)
	//lib.GetLogMgr().Output(lib.LOG_ERROR, "允许奖池4", levelLst4)
	//lib.GetLogMgr().Output(lib.LOG_ERROR, "允许奖池5", levelLst5)
	adminState := self.GetScoreState()
	if 1 == adminState { //送分
		if len(levelLst1) > 0 {
			lib.HF_DeepCopy(&self.Result, &levelLst1[lib.HF_GetRandom(len(levelLst1))].Result)
		} else { //纯随机
			for i := 0; i < len(self.Result); i++ {
				for j := 0; j < len(self.Result[i]); j++ {
					self.Result[i][j] = lib.HF_GetRandom(11) + 1
				}
			}
		}
	} else if 2 == adminState {
		if len(levelLst2) > 0 {
			lib.HF_DeepCopy(&self.Result, &levelLst2[lib.HF_GetRandom(len(levelLst2))].Result)
		} else { //纯随机
			for i := 0; i < len(self.Result); i++ {
				for j := 0; j < len(self.Result[i]); j++ {
					self.Result[i][j] = lib.HF_GetRandom(11) + 1
				}
			}
		}
	} else if 4 == adminState {
		if len(levelLst4) > 0 {
			lib.HF_DeepCopy(&self.Result, &levelLst4[lib.HF_GetRandom(len(levelLst4))].Result)
		} else { //纯随机
			lib.HF_DeepCopy(&self.Result, &lostLst[lib.HF_GetRandom(len(lostLst))].Result)
		}
	} else if 5 == adminState {
		if len(levelLst5) > 0 {
			lib.HF_DeepCopy(&self.Result, &levelLst5[lib.HF_GetRandom(len(levelLst5))].Result)
		} else { //纯随机
			lib.HF_DeepCopy(&self.Result, &lostLst[lib.HF_GetRandom(len(lostLst))].Result)
		}
	} else { //奖池模式
		if GetServer().DwDyjlbSysMoney[self.room.Type%10000] <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && len(lostLst) > 0 {
			//!奖池不够必须输
			lib.HF_DeepCopy(&self.Result, &lostLst[lib.HF_GetRandom(len(lostLst))].Result)
		} else if len(winLst) > 0 {
			lib.HF_DeepCopy(&self.Result, &winLst[lib.HF_GetRandom(len(winLst))].Result)
		} else if len(lostLst) > 0 { //! 都找不到就纯随机输1
			lib.HF_DeepCopy(&self.Result, &lostLst[lib.HF_GetRandom(len(lostLst))].Result)
		} else { //! 都找不到就纯随机
			for i := 0; i < len(self.Result); i++ {
				for j := 0; j < len(self.Result[i]); j++ {
					self.Result[i][j] = lib.HF_GetRandom(11) + 1
				}
			}
		}
	}

	scatArr := make([]int, 0) //
	if self.StFreeRun > 0 {
		totalWild := 0
		for i := 0; i < len(self.Result); i++ {
			count := 0
			for j := 0; j < len(self.Result[i]); j++ { //DYJLB_SCAT
				if DYJLB_SCAT == self.Result[i][j] {
					count++
					totalWild++
				}
			}
			if count <= 0 {
				scatArr = append(scatArr, i)
			}
		}
		lib.GetLogMgr().Output(lib.LOG_ERROR, "aaaaaGetScoreState", scatArr, totalWild, self.StFreeRun)
		if totalWild < self.StFreeRun {
			for i := 0; i < self.StFreeRun-totalWild; i++ {
				index1 := scatArr[i]
				index2 := lib.HF_GetRandom(3)
				lib.GetLogMgr().Output(lib.LOG_ERROR, "sssssGetScoreState", index1, index2)
				self.Result[index1][index2] = DYJLB_SCAT
			}
		}
		self.StFreeRun = 0
	}
	lib.GetLogMgr().Output(lib.LOG_ERROR, "结果", self.Result)
	self.OnEnd()
}
func (self *Game_DYJLB) GetScoreState() int {
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
func (self *Game_DYJLB) GetPerLineScore(line []int) int {
	card := make([]int, 0)
	score := 0 //这条线的得分

	for i := 0; i < len(line); i++ {
		card = append(card, self.Result[line[i]/10][line[i]%10])
	}
	var temppic = 0
	for i := 0; i < len(card); i++ {
		if card[i] == DYJLB_WILD {
			continue
		}
		temppic = card[i]
		break
	}
	num := 0
	for i := 0; i < len(card); i++ {
		if card[i] == DYJLB_WILD || card[i] == temppic {
			num++
		} else {
			break
		}
	}
	if num >= 3 && temppic > DYJLB_SCAT {
		score += self.PL[temppic][num-3] * self.Bet
	}
	return score
}
func (self *Game_DYJLB) OnBye() {
}
func (self *Game_DYJLB) OnExit(uid int64) {
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
func (self *Game_DYJLB) OnEnd() {
	totalWild := 0
	for i := 0; i < len(self.Result); i++ {
		for j := 0; j < len(self.Result[j]); j++ {
			if self.Result[i][j] == DYJLB_SCAT {
				totalWild++
			}
		}
	}

	var msg Msg_GameDYJLB_End
	msg.WinLine = make([]int, 0)

	totalwin := 0 //总赢分
	for j := 0; j < DYJLB_TOTALLINE; j++ {
		_perwin := self.GetPerLineScore(self.Line[j])
		totalwin += _perwin
		if _perwin > 0 {
			msg.WinLine = append(msg.WinLine, j)
		}
	}
	self.Person.Win = totalwin

	playerCost := lib.GetManyMgr().GetProperty(self.room.Type).Cost  //玩家抽水
	sysCost := lib.GetManyMgr().GetProperty(self.room.Type).DealCost //系统抽水
	//系统抽水
	sysWin := 0
	if self.FreeTime > 0 {
		sysWin = 0 - self.Person.Win
	} else {
		sysWin = DYJLB_TOTALLINE*self.Bet - self.Person.Win
	}
	if sysWin != 0 {
		GetServer().SqlBZWLog(&SQL_BZWLog{1, sysWin, time.Now().Unix(), self.room.Type})
	}
	if sysWin > 0 {
		cost := int(math.Ceil(float64(sysWin) * sysCost / 100.0))
		sysWin -= cost
	}
	//!更新奖池
	GetServer().SetDwDyjlbSysMoney(self.room.Type%10000, GetServer().DwDyjlbSysMoney[self.room.Type%10000]+int64(sysWin))
	//！玩家抽水
	playerWin := 0
	self.TotalMul = self.Person.Win / self.Bet
	if self.FreeTime > 0 {
		playerWin = self.Person.Win
	} else {
		playerWin = self.Person.Win - DYJLB_TOTALLINE*self.Bet
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
		if 3 == totalWild {
			self.FreeTime = 3
		} else if 4 == totalWild {
			self.FreeTime = 5
		} else {
			self.FreeTime = 10
		}
	} else if self.FreeTime > 0 {
		isFree = true
		self.FreeTime--
		if self.FreeTime <= 0 {
			self.FreeTime = 0
			self.FreeBet = 0
		}
	}
	msg.TotalMul = self.TotalMul
	msg.Uid = self.Person.Uid
	msg.Win = self.Person.Win
	lib.HF_DeepCopy(&msg.Result, &self.Result)
	msg.FreeNum = self.FreeTime
	msg.Total = self.Person.Total
	self.room.SendMsg(self.Person.Uid, "gamedyjlbend", &msg)
	//lib.GetLogMgr().Output(lib.LOG_ERROR, "结果分数"ua, totalwin, playerCost, msg.Win)

	var record Rec_DYJLB_Info
	record.GameType = self.room.Type
	record.Time = time.Now().Unix()
	var rec Son_Rec_DYJLB_Person
	rec.Uid = self.Person.Uid
	rec.Name = self.Person.Name
	rec.Head = self.Person.Head
	if isFree {
		rec.Score = self.Person.Win
		rec.Bets = 0
	} else {
		rec.Score = self.Person.Win - DYJLB_TOTALLINE*self.Bet
		rec.Bets = DYJLB_TOTALLINE * self.Bet
	}
	record.Info = append(record.Info, rec)
	GetServer().InsertRecord(self.room.Type, self.Person.Uid, lib.HF_JtoA(&record), rec.Score)
	////////
	self.room.Begin = false
	self.TotalMul = 0
	self.Person.Win = 0
	self.Person.Cost = 0

	self.SetTime(GOLDDYJLB_TIME)
}
func (self *Game_DYJLB) OnInit(room *Room) {
	self.room = room
}
func (self *Game_DYJLB) OnIsBets(uid int64) bool {
	return false
}
func (self *Game_DYJLB) OnIsDealer(uid int64) bool {
	return false
}
func (self *Game_DYJLB) OnRobot(robot *lib.Robot) {

}
func (self *Game_DYJLB) getInfo(uid int64) *Msg_GameDYJLB_Info {
	var msg Msg_GameDYJLB_Info
	msg.FreeTime = self.FreeTime
	df := staticfunc.GetCsvMgr().GetDF(self.room.Type)
	msg.Money = append(msg.Money, df)
	msg.JackPot = int64(GOLDDYJLB_JACKBASE) + GetServer().DwDyjlbSysMoney[self.room.Type%10000]
	//lib.GetLogMgr().Output(lib.LOG_ERROR, "结果分数getInfo", df, msg.Money)
	if GetServer().IsAdmin(uid, staticfunc.ADMIN_DWDYJLB) {
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
		msg.Person.Address = self.Person.Address
		msg.Person.Sex = self.Person.Sex
	}
	return &msg
}
func (self *Game_DYJLB) OnSendInfo(person *Person) {
	if self.Person != nil && self.Person.Uid == person.Uid {
		self.Person.SynchroGold(person.Gold)
		person.SendMsg("gamedyjlbinfo", self.getInfo(person.Uid))
		return
	}

	//lib.GetLogMgr().Output(lib.LOG_ERROR, "OnSendInfo:", person.Uid)

	_person := new(Game_DYJLB_Person)
	_person.Uid = person.Uid
	_person.Gold = person.Gold
	_person.Total = person.Gold
	_person.Name = person.Name
	_person.Sex = person.Sex
	_person.IP = person.ip
	_person.Head = person.Imgurl
	_person.Address = person.minfo.Address
	self.Person = _person
	//self.Init()
	person.SendMsg("gamedyjlbinfo", self.getInfo(person.Uid))

	self.SetTime(GOLDDYJLB_TIME)

}

func (self *Game_DYJLB) FlushJackState() {
	var msg Msg_GameDYJLB_JackPot
	msg.Uid = self.Person.Uid
	msg.JackPot = int64(GOLDDYJLB_JACKBASE) + GetServer().DwDyjlbSysMoney[self.room.Type%10000]
	self.room.SendMsg(self.Person.Uid, "gamedyjlbjackpot", &msg)
	self.JackTime = time.Now().Unix() + 3
}

//! 设置时间
func (self *Game_DYJLB) SetTime(t int) {
	if t == 0 {
		self.Time = 0
	} else {
		self.Time = time.Now().Unix() + int64(t)
	}

	var msg Msg_SetTime
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	self.room.broadCastMsg("gametime", &msg)
}
func (self *Game_DYJLB) OnTime() {
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
