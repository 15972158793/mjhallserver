package gameserver

import (
	//	"fmt"
	"lib"
	"math"
	//	"sort"
	"staticfunc"
	"time"
)

var GOLDSHZ_TIME = 3600

var SHZ_SHZ int = 1  //! 水浒传
var SHZ_ZYT int = 2  //! 忠义堂
var SHZ_TTXD int = 3 //! 替天行道
var SHZ_SJ int = 4   //! 宋江
var SHZ_LC int = 5   //! 林冲
var SHZ_LZS int = 6  //! 鲁智深
var SHZ_D int = 7    //! 刀
var SHZ_Q int = 8    //! 枪
var SHZ_F int = 9    //! 斧

var RUN []int = []int{SHZ_ZYT, SHZ_F, SHZ_D, -1, SHZ_SJ, SHZ_LZS, SHZ_Q, SHZ_F, SHZ_TTXD, -1, SHZ_LC, SHZ_D, SHZ_F, SHZ_LZS, SHZ_SJ, -1, SHZ_LC, SHZ_Q, SHZ_F, SHZ_D, SHZ_TTXD, -1, SHZ_LZS, SHZ_Q}

type Rec_SHZ_Info struct {
	GameType int                  `json:"gametype"`
	Time     int64                `json:"time"` //! 记录时间
	Info     []Son_Rec_SHZ_Person `json:"info"`
}
type Son_Rec_SHZ_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Score int    `json:"score"`
	Bets  int    `json:"bets"`
}

type Game_GoldSHZ struct {
	Person      *Game_GoldSHZ_Person `json:"person"`
	Reslut      [5][3]int            `json:"result"` //! 左下角坐标为（0,0）右上角坐标为（4,2）
	Time        int64                `json:"time"`
	LineNum     int                  `json:"linenum"`   //! 选了几条线
	Bet         int                  `json:"bet"`       //! 每条线的花费
	PL          map[int][]int        `json:"pl"`        //! 赔率  []int{三个,四个，五个，全屏}
	RunPL       map[int]int          `json:"runpl"`     //! 跑灯赔率
	Line        [][]int              `json:"line"`      //! 线路
	RunNum      int                  `json:"runnum"`    //! 跑灯次数
	YSZTrend    []int                `json:"ysztrend"`  //! 摇色子战绩
	CurRunNum   int                  `json:"currunnum"` //! 当跑灯大于等于5的时候必然会转到exit
	State       int                  `json:"state"`     //! 0-游戏未开始 1-摇色子 2-摇色子选择大小 3-跑灯
	DoubleMoney int                  `json:"doublemoney"`

	room *Room
}

func NewGame_GoldSHZ() *Game_GoldSHZ {
	game := new(Game_GoldSHZ)
	game.Line = make([][]int, 0)

	/*
		(0,0) (1,0) (2,0) (3,0) (4,0)
		(0,1) (1,1) (2,1) (3,1) (4,1)
		(0,2) (1,2) (2,2) (3,2) (4,2)
	*/
	//!　初始化9条路线
	game.Line = append(game.Line, []int{1, 11, 21, 31, 41})
	game.Line = append(game.Line, []int{0, 10, 20, 30, 40})
	game.Line = append(game.Line, []int{2, 12, 22, 32, 42})
	game.Line = append(game.Line, []int{0, 11, 22, 31, 40})
	game.Line = append(game.Line, []int{2, 11, 20, 31, 42})
	game.Line = append(game.Line, []int{0, 10, 21, 30, 40})
	game.Line = append(game.Line, []int{2, 12, 21, 32, 42})
	game.Line = append(game.Line, []int{1, 12, 22, 32, 41})
	game.Line = append(game.Line, []int{1, 10, 20, 30, 41})
	//! 初始化每个图标的赔率
	game.PL = make(map[int][]int, 0)
	game.PL[SHZ_SHZ] = []int{1, 1, 1000, 5000}
	game.PL[SHZ_ZYT] = []int{50, 200, 1000, 2500}
	game.PL[SHZ_TTXD] = []int{20, 80, 400, 1000}
	game.PL[SHZ_SJ] = []int{15, 40, 200, 500}
	game.PL[SHZ_LC] = []int{10, 30, 160, 400}
	game.PL[SHZ_LZS] = []int{7, 20, 100, 250}
	game.PL[SHZ_D] = []int{5, 15, 60, 150}
	game.PL[SHZ_Q] = []int{3, 10, 40, 100}
	game.PL[SHZ_F] = []int{2, 5, 20, 50}

	//! 跑灯赔率
	game.RunPL = make(map[int]int, 0)
	game.RunPL[SHZ_ZYT] = 200
	game.RunPL[SHZ_TTXD] = 100
	game.RunPL[SHZ_SJ] = 70
	game.RunPL[SHZ_LC] = 50
	game.RunPL[SHZ_LZS] = 20
	game.RunPL[SHZ_D] = 10
	game.RunPL[SHZ_Q] = 5
	game.RunPL[SHZ_F] = 2

	game.CurRunNum = 0
	game.LineNum = 9
	game.Bet = 100
	game.RunNum = 0
	game.State = 0

	return game
}

type Game_GoldSHZ_Person struct {
	Uid     int64  `json:"uid"`
	Gold    int    `json:"gold"`
	Total   int    `json:"total"`
	Win     int    `json:"win"`    //! 单局总得分
	RunWin  int    `json:"runwin"` //! 跑灯得分
	Cost    int    `json:"cost"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type Msg_GameGoldSHZ_Total struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"`
}

type Msg_GameGoldSHZ_CanResult struct {
	Result [5][3]int `json:"result"`
}

type Msg_GameGoldSHZ_RunCanResult struct {
	Icon  []int `json:"icon"`
	Index int   `json:"index"`
}

type Msg_GameGoldSHZ_Result struct {
	Uid     int64     `json:"uid"`
	Win     int       `json:"win"`
	Result  [5][3]int `json:"result"`
	RunNum  int       `json:"runnum"`  //! 跑灯次数
	DelAll  int       `json:"delall"`  //! 是否是全屏消  0~否 1-9~水浒传-斧 10~全人物 11~全武器
	WinLine []int     `json:"winline"` //! 哪几条线得分
}

type Msg_GameGoldSHZ_YSZResult struct {
	Uid    int64 `json:"uid"`
	Win    int   `json:"win"` //! 当前赢分
	Trend  []int `json:"trend"`
	Result []int `json:"result"`
	Choose int   `json:"choose"` //! 选择压大小还是退出
}

type Msg_GameGoldSHZ_YSZChoose struct {
	Uid    int64 `json:"uid"`
	Choose int   `json:"choose"`
	CurWin int   `json:"curwin"`
}

type Msg_GameGoldSHZ_RunResult struct {
	Uid  int64                       `json:"uid"`
	Info []Son_GameGoldSHZ_RunResult `json:"info"`
}
type Son_GameGoldSHZ_RunResult struct {
	CurRunWin int   `json:"currunwin"` //! 这一转赢了多少
	RunNum    int   `json:"runnum"`    //! 跑灯次数
	Index     int   `json:"index"`     //! 落点下标
	Icon      []int `json:"icon"`      //! 中间四个图案图标
}

type Msg_GameGoldSHZ_End struct {
	Uid   int64 `json:"uid"`
	Win   int   `json:"win"`
	Total int   `json:"total"`
}

type Msg_GameGoldSHZ_Info struct {
	Begin   bool                    `json:"begin"`
	Result  [5][3]int               `json:"result"`
	LineNum int                     `json:"linenum"` //! 选了几条线
	Bet     int                     `json:"bet"`     //! 每条线的花费
	RunNum  int                     `json:"runnum"`  //! 跑灯次数
	State   int                     `json:"state"`   //! 1-选择怎么加倍加倍 2-选择压大小  3-跑灯
	Person  Son_GameGoldSHZ_Info    `json:"person"`
	YSZ     Son_GameGoldSHZ_YSZInfo `json:"ysz"`
}

type Son_GameGoldSHZ_YSZInfo struct {
	CurScore int   `json:"curscore"`
	Trend    []int `json:"trend"`
}

type Son_GameGoldSHZ_Info struct {
	Uid     int64  `json:"uid"`
	Total   int    `json:"total"`
	Win     int    `json:"win"`    //! 单局总得分
	RunWin  int    `json:"runwin"` //! 跑灯得分
	Name    string `json:"name"`   //! 名字
	Head    string `json:"head"`   //! 头像
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

func (self *Game_GoldSHZ) getInfo(uid int64) *Msg_GameGoldSHZ_Info {
	var msg Msg_GameGoldSHZ_Info
	msg.Begin = self.room.Begin
	lib.HF_DeepCopy(&msg.Result, &self.Reslut)
	msg.RunNum = self.RunNum
	msg.LineNum = self.LineNum
	msg.Bet = self.Bet
	msg.State = self.State
	msg.YSZ.Trend = make([]int, 0)
	lib.HF_DeepCopy(&msg.YSZ.Trend, &self.YSZTrend)
	if self.State == 2 || self.State == 1 {
		msg.YSZ.CurScore = self.Person.Win
	}
	if self.Person != nil && self.Person.Uid == uid {
		msg.Person.Uid = uid
		msg.Person.Total = self.Person.Total
		msg.Person.Win = self.Person.Win
		msg.Person.RunWin = self.Person.RunWin
		msg.Person.Name = self.Person.Name
		msg.Person.Head = self.Person.Head
		msg.Person.IP = self.Person.IP
		msg.Person.Address = self.Person.Address
		msg.Person.Sex = self.Person.Sex
	}
	return &msg
}

func (self *Game_GoldSHZ_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}
func (self *Game_GoldSHZ) SendTotal(uid int64, total int) {
	var msg Msg_GameGoldSHZ_Total
	msg.Uid = uid
	msg.Total = total
	self.room.SendMsg(uid, "gamegoldtotal", &msg)
}

//! 设置时间
func (self *Game_GoldSHZ) SetTime(t int) {
	if t == 0 {
		self.Time = 0
	} else {
		self.Time = time.Now().Unix() + int64(t)
	}

	var msg Msg_SetTime
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	self.room.broadCastMsg("gametime", &msg)
}

//! 是否消全屏 return 分数 全屏消的是什么(10~全屏人物 11~全屏武器)
func (self *Game_GoldSHZ) IsDelAll() (int, int) {
	del := true
	for i := 0; i < len(self.Reslut); i++ {
		for j := 0; j < len(self.Reslut[i]); j++ {
			if self.Reslut[i][j] != self.Reslut[0][0] {
				del = false
				break
			}
		}
		if !del {
			break
		}
	}

	if del { //! 全屏同一图标
		return self.LineNum * self.PL[self.Reslut[0][0]][3] * self.Bet, self.Reslut[0][0]
	}

	del = true //! 判断是否全屏人物或者全屏武器
	for i := 0; i < len(self.Reslut); i++ {
		for j := 0; j < len(self.Reslut[i]); j++ {
			if self.Reslut[0][0] == SHZ_LC || self.Reslut[0][0] == SHZ_SJ || self.Reslut[0][0] == SHZ_LZS {
				if self.Reslut[i][j] != SHZ_LC && self.Reslut[i][j] != SHZ_SJ && self.Reslut[i][j] != SHZ_LZS {
					del = false
					break
				}
			} else if self.Reslut[0][0] == SHZ_D || self.Reslut[0][0] == SHZ_Q || self.Reslut[0][0] == SHZ_F {
				if self.Reslut[i][j] != SHZ_D && self.Reslut[i][j] != SHZ_Q && self.Reslut[i][j] != SHZ_F {
					del = false
					break
				}
			} else {
				del = false
				break
			}
		}
		if !del {
			break
		}
	}

	if del {
		if self.Reslut[0][0] == SHZ_LC || self.Reslut[0][0] == SHZ_SJ || self.Reslut[0][0] == SHZ_LZS {
			return self.LineNum * 50 * self.Bet, 10
		} else {
			return self.LineNum * 50 * self.Bet, 11
		}
	}

	return 0, 0
}

//! 获取一条线路的得分  int~得分 int~哪几条线得分 int~送跑灯游戏的次数
func (self *Game_GoldSHZ) GetScoreByLine(line []int, _num int) (int, int, int) {
	score := 0
	del := false
	run := 0

	card := make([]int, 0) //！ 取这条线路的数据
	for i := 0; i < len(line); i++ {
		card = append(card, self.Reslut[line[i]/10][line[i]%10])
	}

	//! 获取这条线是否有能送跑灯游戏
	run = self.GetRunNum(card)

	//! 这一条线全部是水浒传
	if card[0] == SHZ_SHZ && card[1] == SHZ_SHZ && card[2] == SHZ_SHZ && card[3] == SHZ_SHZ && card[4] == SHZ_SHZ {
		score = self.PL[SHZ_SHZ][2] * self.Bet
		return score, _num, run
	}

	//! 先从右到左找相同图标
	right := 0
	for i := 0; i < len(card); i++ {
		if card[i] == 1 {
			continue
		}
		right = card[i]
		break
	}
	num := 0
	for i := 0; i < len(card); i++ {
		if card[i] == 1 || card[i] == right {
			num++
		} else {
			break
		}
	}
	if num >= 3 {
		score += self.PL[right][num-3] * self.Bet
		lst := make([]int, 0)
		for i := 0; i < num; i++ {
			lst = append(lst, line[i])
		}
		del = true
	}

	if num < 5 {
		//! 从左到右找相同图标
		left := 0
		for i := len(card) - 1; i >= 0; i-- {
			if card[i] == 1 {
				continue
			}
			left = card[i]
			break
		}
		num = 0
		for i := len(card) - 1; i >= 0; i-- {
			if card[i] == 1 || card[i] == left {
				num++
			} else {
				break
			}
		}

		if num >= 3 {
			score += self.PL[left][num-3] * self.Bet
			lst := make([]int, 0)
			for i := len(card) - 1; i >= len(card)-num; i-- {
				lst = append(lst, i)
			}
			del = true
		}
	}

	if del {
		return score, _num, run
	} else {
		return score, -1, run
	}

}

//! 返回这条线免费跑灯数量
func (self *Game_GoldSHZ) GetRunNum(card []int) int {
	num := 0
	for i := 0; i < 3; i++ {
		if card[i] == SHZ_SHZ {
			_num := 1
			for j := i + 1; j < len(card); j++ {
				if card[j] == SHZ_SHZ {
					_num++
				} else {
					break
				}
			}
			if _num >= 3 && _num > num {
				num = _num
			}
		}
	}
	if num == 3 {
		return 1
	} else if num == 4 {
		return 2
	} else if num == 5 {
		return 3
	}
	return 0
}

//! 开始游戏
func (self *Game_GoldSHZ) GameStart(uid int64, lineNum int, gold int) {
	if self.Person == nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "GameStart person == nil")
		return
	}

	if lineNum > 9 || lineNum < 1 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "线路有误 linenum : ", lineNum)
		return
	}

	if uid != self.Person.Uid {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "GameStart person.uid != uid")
		return
	}

	if self.State != 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏已经开始 state111 : ", self.State)
		return
	}

	if lineNum*gold > self.Person.Total {
		person := GetPersonMgr().GetPerson(uid)
		person.SendErr("金币不足！")
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "金币不足!")
		return
	}
	self.SetTime(GOLDSHZ_TIME)

	self.LineNum = lineNum
	self.Bet = gold
	self.Person.Total -= self.LineNum * self.Bet
	self.SendTotal(self.Person.Uid, self.Person.Total)

	self.OnBegin()
}

func (self *Game_GoldSHZ) OnBegin() {
	if self.State != 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏已经开始 state : ", self.State)
		return
	}

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始！")
	self.room.Begin = true

	lostLst := make([]Msg_GameGoldSHZ_CanResult, 0)
	//winLessLst := make([]Msg_GameGoldSHZ_CanResult, 0)
	//winMoreLst := make([]Msg_GameGoldSHZ_CanResult, 0)
	//specialLst := make([]Msg_GameGoldSHZ_CanResult, 0)
	canLst := make([]Msg_GameGoldSHZ_CanResult, 0)

	for k := 0; k < 100; k++ {
		for i := 0; i < len(self.Reslut); i++ {
			for j := 0; j < len(self.Reslut[i]); j++ {
				self.Reslut[i][j] = lib.HF_GetRandom(9) + 1
			}
		}

		score := 0                 //! 玩家赢多少
		runNum := 0                //! 跑灯次数
		score, _ = self.IsDelAll() //! 是否全屏同一图标
		if score <= 0 {            //! 非全屏同一图标，计算每条路线收益
			for i := 0; i < self.LineNum; i++ {
				_score, _, _runNum := self.GetScoreByLine(self.Line[i], i)
				score += _score
				if _runNum > runNum {
					runNum = _runNum
				}
			}
		}

		//! 输了的列表
		if score == 0 || score < self.Bet*self.LineNum {
			var msg Msg_GameGoldSHZ_CanResult
			lib.HF_DeepCopy(&msg.Result, &self.Reslut)
			lostLst = append(lostLst, msg)
		}

		//! 在奖池之类的列表
		if GetServer().SHZSysmoney[self.room.Type%270000]+int64(self.Bet*self.LineNum)-int64(score) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && GetServer().SHZSysmoney[self.room.Type%270000]+int64(self.Bet*self.LineNum)-int64(score) <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
			var msg Msg_GameGoldSHZ_CanResult
			lib.HF_DeepCopy(&msg.Result, &self.Reslut)
			canLst = append(canLst, msg)
			//if runNum > 0 {
			//	specialLst = append(specialLst, msg)
			//} else {
			//	winMoreLst = append(winMoreLst, msg)
			//}
		}

		if len(canLst) >= 50 {
			break
		}

		//if len(winLessLst) >= 10 && len(winMoreLst) >= 10 && len(lostLst) >= 10 ||  {
		//	break
		//}
	}

	//lib.GetLogMgr().Output(lib.LOG_DEBUG, "------------len(lostlst) : ", len(lostLst))
	//lib.GetLogMgr().Output(lib.LOG_DEBUG, "------------len(winLessLst) : ", len(winLessLst))
	//lib.GetLogMgr().Output(lib.LOG_DEBUG, "------------len(winMoreLst) : ", len(winMoreLst))
	//lib.GetLogMgr().Output(lib.LOG_DEBUG, "------------len(specialLst) : ", len(specialLst))

	if GetServer().SHZSysmoney[self.room.Type%270000] <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && len(lostLst) > 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "------------ 必输 随机lostlst")
		index := lib.HF_GetRandom(len(lostLst))
		lib.HF_DeepCopy(&self.Reslut, &lostLst[index].Result)
	} else if len(canLst) > 0 {
		index := lib.HF_GetRandom(len(canLst))
		lib.HF_DeepCopy(&self.Reslut, &canLst[index].Result)

		//num := lib.HF_GetRandom(100)
		//if num < 70 && len(winLessLst) > 0 {
		//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "------------- 赢少许 随机winlesslst")
		//	index := lib.HF_GetRandom(len(winLessLst))
		//	lib.HF_DeepCopy(&self.Reslut, &winLessLst[index].Result)
		//} else if num < 78 && len(winMoreLst) > 0 {
		//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "--------- 赢许多 随机winmorelst")
		//	index := lib.HF_GetRandom(len(winMoreLst))
		//	lib.HF_DeepCopy(&self.Reslut, &winMoreLst[index].Result)
		//} else if num < 83 && len(specialLst) > 0 {
		//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "---------- 特殊玩法 随机speciallst")
		//	index := lib.HF_GetRandom(len(specialLst))
		//	lib.HF_DeepCopy(&self.Reslut, &specialLst[index].Result)
		//} else {
		//	if len(lostLst) > 0 {
		//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------- 输钱 随机lostlst")
		//		index := lib.HF_GetRandom(len(lostLst))
		//		lib.HF_DeepCopy(&self.Reslut, &lostLst[index].Result)
		//	} else {
		//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "--------- 纯随机 ")
		//		for i := 0; i < len(self.Reslut); i++ {
		//			for j := 0; j < len(self.Reslut[i]); j++ {
		//				self.Reslut[i][j] = lib.HF_GetRandom(9) + 1
		//			}
		//		}
		//	}
		//}
	} else { //! 都找不到就纯随机
		for i := 0; i < len(self.Reslut); i++ {
			for j := 0; j < len(self.Reslut[i]); j++ {
				self.Reslut[i][j] = lib.HF_GetRandom(9) + 1
			}
		}
	}
	//self.Reslut[0][0] = SHZ_SHZ
	//self.Reslut[1][0] = SHZ_SHZ
	//self.Reslut[2][0] = SHZ_SHZ

	self.SendResult()
}

//! 发送结果
func (self *Game_GoldSHZ) SendResult() {
	win, delall := self.IsDelAll()
	if win > 0 { //! 消全屏
		var msg Msg_GameGoldSHZ_Result
		msg.Uid = self.Person.Uid
		msg.Win = win
		msg.DelAll = delall
		msg.RunNum = 0
		lineLst := make([]int, 0)
		for i := 0; i < len(self.Reslut); i++ {
			for j := 0; j < len(self.Reslut[i]); j++ {
				lineLst = append(lineLst, i*10+j)
			}
		}
		msg.WinLine = []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
		lib.HF_DeepCopy(&msg.Result, &self.Reslut)
		self.State = 1 //! 摇色子阶段
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "------ SendResult() 选择是否摇色子 state ", self.State)
		self.room.SendMsg(self.Person.Uid, "gamegoldshzresult", &msg)

		self.Person.Win = msg.Win
		return
	}

	var msg Msg_GameGoldSHZ_Result
	msg.Uid = self.Person.Uid
	msg.DelAll = 0
	msg.WinLine = make([]int, 0)
	lib.HF_DeepCopy(&msg.Result, &self.Reslut)
	for i := 0; i < self.LineNum; i++ { //!计算每条线路的分数
		score, del, run := self.GetScoreByLine(self.Line[i], i)
		msg.Win += score
		if del != -1 {
			msg.WinLine = append(msg.WinLine, del)
		}
		if msg.RunNum < run {
			msg.RunNum = run
		}
	}
	if msg.RunNum <= 0 {
		self.State = 1
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "------ SendResult()1111 选择是否摇色子 state ", self.State)
	} else {
		self.State = 3
		self.RunNum = msg.RunNum
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "------ SendResult() 跑灯阶段 state ", self.State)
	}
	self.Person.Win = msg.Win
	self.room.SendMsg(self.Person.Uid, "gamegoldshzresult", &msg)

	if self.State == 3 {
		self.GameRun()
		return
	}

	if self.Person.Win <= 0 {
		self.OnEnd()
		return
	}
}

//! 摇色子玩法 0-不摇 1-普通加倍 2-双倍加倍
func (self *Game_GoldSHZ) GameYSZ(uid int64, choose int) {
	if self.Person != nil && self.Person.Uid != uid {
		return
	}

	if self.State != 1 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "GameYSZ 状态错误 state : ", self.State)
		return
	}
	self.State = 2
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "------ GameYSZ() 摇色子选择大小 state ", self.State)

	if choose == 2 {
		self.Person.Total -= self.Person.Win
		self.DoubleMoney += self.Person.Win
		self.Person.Win *= 2
		self.SendTotal(self.Person.Uid, self.Person.Total)
	}

	var msg Msg_GameGoldSHZ_YSZChoose
	msg.Uid = uid
	msg.Choose = choose
	msg.CurWin = self.Person.Win
	self.room.SendMsg(uid, "gamegoldshzyszchoose", &msg)

	if choose == 0 {
		self.OnEnd()
		return
	}

	self.SetTime(GOLDSHZ_TIME)
}

//! 压大小 0-小 1-和 2-大 3-得分退出
func (self *Game_GoldSHZ) GameYSZResult(uid int64, index int) {
	if self.Person != nil && self.Person.Uid != uid {
		return
	}

	if self.State != 2 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "GameYSZResult 状态错误 state : ", self.State)
		return
	}

	self.SetTime(GOLDSHZ_TIME)

	if index == 3 {
		var msg Msg_GameGoldSHZ_YSZResult
		msg.Uid = self.Person.Uid
		msg.Choose = index
		self.room.SendMsg(self.Person.Uid, "gamegoldshzyszresult", &msg)
		self.OnEnd()
		return
	}

	lst := make([][]int, 0)
	lostLst := make([][]int, 0)
	for i := 0; i < 50; i++ {
		result := make([]int, 0)
		result = append(result, lib.HF_GetRandom(6)+1)
		result = append(result, lib.HF_GetRandom(6)+1)

		playwin := self.GetYSZWin(index, result)
		if playwin <= 0 { //! 必输的
			lostLst = append(lostLst, result)
		}

		//! 在奖池内的
		if GetServer().SHZSysmoney[self.room.Type%270000]+int64(self.LineNum*self.Bet)-int64(playwin) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && GetServer().SHZSysmoney[self.room.Type%270000]+int64(self.LineNum*self.Bet)-int64(playwin) <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
			lst = append(lst, result)
		}
	}
	var msg Msg_GameGoldSHZ_YSZResult
	msg.Uid = self.Person.Uid
	msg.Choose = index
	if GetServer().SHZSysmoney[self.room.Type%270000] < lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && len(lostLst) > 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "------- 摇色子玩法 必输  随机lostLst")
		msg.Result = lostLst[lib.HF_GetRandom(len(lostLst))]
		msg.Win = self.GetYSZWin(index, msg.Result)
	} else {
		if len(lst) > 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "------- 摇色子玩法 随机lst ")
			msg.Result = lst[lib.HF_GetRandom(len(lst))]
			msg.Win = self.GetYSZWin(index, msg.Result)
		} else {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "------- 摇色子玩法 纯随机 ")
			msg.Result = append(msg.Result, lib.HF_GetRandom(6)+1)
			msg.Result = append(msg.Result, lib.HF_GetRandom(6)+1)
			msg.Win = self.GetYSZWin(index, msg.Result)
		}
	}

	tmp := make([]int, 0)
	if msg.Result[0]+msg.Result[1] == 7 {
		tmp = append(tmp, 1)
	} else if msg.Result[0]+msg.Result[1] <= 6 {
		tmp = append(tmp, 0)
	} else {
		tmp = append(tmp, 2)
	}
	tmp = append(tmp, self.YSZTrend...)
	if len(tmp) > 20 {
		tmp = tmp[0:20]
	}
	self.YSZTrend = tmp

	self.Person.Win = msg.Win
	lib.HF_DeepCopy(&msg.Trend, &self.YSZTrend)
	self.room.SendMsg(self.Person.Uid, "gamegoldshzyszresult", &msg)
	if self.Person.Win == 0 {
		self.OnEnd()
	}
}

//! 摇色子这样开玩家赢多少
func (self *Game_GoldSHZ) GetYSZWin(index int, result []int) int {
	if index == 0 {
		if result[0]+result[1] <= 6 {
			if result[0] == result[1] {
				return self.Person.Win * 4
			} else {
				return self.Person.Win * 2
			}
		}
	} else if index == 1 {
		if result[0]+result[1] == 7 {
			return self.Person.Win * 6
		}
	} else {
		if result[0]+result[1] >= 8 {
			if result[0] == result[1] {
				return self.Person.Win * 4
			} else {
				return self.Person.Win * 2
			}
		}
	}
	return 0
}

//! 跑灯玩法
func (self *Game_GoldSHZ) GameRun() {
	if self.State != 3 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "state != 2 : ", self.State)
		return
	}

	if self.RunNum <= 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "self.runnum <= 0  : ", self.RunNum)
		return
	}

	var msg Msg_GameGoldSHZ_RunResult
	msg.Uid = self.Person.Uid

	for j := 0; j < 100; j++ {
		winLessLst := make([]Msg_GameGoldSHZ_RunCanResult, 0)
		winMoreLst := make([]Msg_GameGoldSHZ_RunCanResult, 0)
		lostLst := make([]Msg_GameGoldSHZ_RunCanResult, 0)
		exitLst := make([]Msg_GameGoldSHZ_RunCanResult, 0)

		for i := 0; i < 50; i++ {
			icon := make([]int, 0)
			for i := 0; i < 4; i++ {
				icon = append(icon, lib.HF_GetRandom(8)+2)
			}
			index := lib.HF_GetRandom(len(RUN))
			win, exit := self.GetRunWin(index, icon)

			if exit {
				var msg Msg_GameGoldSHZ_RunCanResult
				lib.HF_DeepCopy(&msg.Icon, &icon)
				msg.Index = index
				exitLst = append(exitLst, msg)
			} else {
				if GetServer().SHZSysmoney[self.room.Type%270000]+int64(self.Bet*self.LineNum)-int64(win) > lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && GetServer().SHZSysmoney[self.room.Type%270000]+int64(self.Bet*self.LineNum)-int64(win) < lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
					if win <= 0 {
						var msg Msg_GameGoldSHZ_RunCanResult
						lib.HF_DeepCopy(&msg.Icon, &icon)
						msg.Index = index
						lostLst = append(lostLst, msg)
					} else if win > 0 && win <= self.Bet*self.RunNum*10 {
						var msg Msg_GameGoldSHZ_RunCanResult
						lib.HF_DeepCopy(&msg.Icon, &icon)
						msg.Index = index
						winLessLst = append(winLessLst, msg)
					} else {
						var msg Msg_GameGoldSHZ_RunCanResult
						lib.HF_DeepCopy(&msg.Icon, &icon)
						msg.Index = index
						winMoreLst = append(winMoreLst, msg)
					}
				}
			}
		}

		index := 0
		icon := make([]int, 0)
		if (lib.HF_GetRandom(100) < 20 || self.CurRunNum >= 5) && len(exitLst) > 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "----- 跑灯玩法 退出 随机exitLst ")
			num := lib.HF_GetRandom(len(exitLst))
			lib.HF_DeepCopy(&icon, &exitLst[num].Icon)
			index = exitLst[num].Index
		} else {
			rand := lib.HF_GetRandom(100)
			if len(winLessLst) > 0 && rand < 30 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "----- 跑灯玩法 赢10倍一下 随机winLessLst ")
				num := lib.HF_GetRandom(len(winLessLst))
				lib.HF_DeepCopy(&icon, &winLessLst[num].Icon)
				index = winLessLst[num].Index
			} else if len(winMoreLst) > 0 && rand < 40 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "----- 跑灯玩法 赢10倍以上 随机winMoreLst ")
				num := lib.HF_GetRandom(len(winMoreLst))
				lib.HF_DeepCopy(&icon, &winMoreLst[num].Icon)
				index = winMoreLst[num].Index
			} else {
				if len(lostLst) > 0 {
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "----- 跑灯玩法 没有赢 随机lostLst ")
					num := lib.HF_GetRandom(len(lostLst))
					lib.HF_DeepCopy(&icon, &lostLst[num].Icon)
					index = lostLst[num].Index
				} else {
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "----- 跑灯玩法 纯随机 ")
					for i := 0; i < 4; i++ {
						icon = append(icon, lib.HF_GetRandom(8)+2)
					}
					index = lib.HF_GetRandom(len(RUN))
				}
			}

		}

		win, exit := self.GetRunWin(index, icon)
		self.Person.RunWin += win
		self.Person.Win += win

		if exit {
			self.RunNum--
			self.CurRunNum = 0
		} else {
			self.CurRunNum++
		}

		var son Son_GameGoldSHZ_RunResult
		son.CurRunWin = win
		son.Icon = make([]int, 0)
		lib.HF_DeepCopy(&son.Icon, &icon)
		son.Index = index
		son.RunNum = self.RunNum
		msg.Info = append(msg.Info, son)

		if self.RunNum <= 0 {
			break
		}
	}
	self.room.SendMsg(self.Person.Uid, "gamegoldshzrunresult", &msg)

	self.OnEnd()

}

//! 跑灯这样开玩家可以赢多少
func (self *Game_GoldSHZ) GetRunWin(index int, icon []int) (int, bool) {
	score := 0
	num := 0
	for i := 0; i < len(icon); i++ {
		if RUN[index] == icon[i] {
			num++
		}
	}

	if num > 0 {
		score += self.RunPL[RUN[index]] * num * self.Bet * self.LineNum
	}

	if score > 0 {
		num = 0
		for i := 0; i < len(icon); i++ {
			if RUN[index] == icon[i] {
				num++
			}
		}
		if num == 3 {
			score *= 20
		}
		if num == 4 {
			score *= 500
		}
	}

	return score, RUN[index] == -1
}

func (self *Game_GoldSHZ) OnEnd() {
	self.room.Begin = false

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "------ self.doublemoney : ", self.DoubleMoney)
	dealwin := self.LineNum*self.Bet + self.DoubleMoney
	dealwin -= self.Person.Win
	if dealwin != 0 {
		GetServer().SqlBZWLog(&SQL_BZWLog{1, dealwin, time.Now().Unix(), self.room.Type})
	}

	if dealwin > 0 {
		b1 := float64(lib.GetManyMgr().GetProperty(self.room.Type).DealCost)
		cost := int(math.Ceil(float64(dealwin) * b1 / 100.0))
		dealwin -= cost
	}
	GetServer().SetSHZSysMoney(self.room.Type%270000, GetServer().SHZSysmoney[self.room.Type%270000]+int64(dealwin))

	if self.Person.Win-self.LineNum*self.Bet-self.DoubleMoney > 0 {
		self.Person.Cost = int(math.Ceil(float64(self.Person.Win-self.LineNum*self.Bet-self.DoubleMoney) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		GetServer().SqlAgentGoldLog(self.Person.Uid, self.Person.Cost, self.room.Type)
		GetServer().SqlAgentBillsLog(self.Person.Uid, self.Person.Cost/2, self.room.Type)
		self.Person.Win -= self.Person.Cost
	} else if self.Person.Win-self.LineNum*self.Bet-self.DoubleMoney < 0 {
		cost := int(math.Ceil(float64(self.Person.Win-self.LineNum*self.Bet-self.DoubleMoney) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 200.0))
		GetServer().SqlAgentBillsLog(self.Person.Uid, cost, self.room.Type)
	}

	self.Person.Total += self.Person.Win

	var msg Msg_GameGoldSHZ_End
	msg.Uid = self.Person.Uid
	msg.Total = self.Person.Total
	msg.Win = self.Person.Win
	self.room.SendMsg(self.Person.Uid, "gamegoldshzend", &msg)

	var record Rec_SHZ_Info
	record.GameType = self.room.Type
	record.Time = time.Now().Unix()
	var rec Son_Rec_SHZ_Person
	rec.Uid = self.Person.Uid
	rec.Name = self.Person.Name
	rec.Head = self.Person.Head
	rec.Score = self.Person.Win - self.LineNum*self.Bet - self.DoubleMoney
	rec.Bets = self.LineNum * self.Bet
	record.Info = append(record.Info, rec)
	GetServer().InsertRecord(self.room.Type, self.Person.Uid, lib.HF_JtoA(&record), rec.Score)

	self.State = 0
	self.Person.Win = 0
	self.Person.Cost = 0
	self.Person.RunWin = 0
	self.RunNum = 0
	self.CurRunNum = 0
	self.DoubleMoney = 0

	self.SetTime(GOLDSHZ_TIME)
}

func (self *Game_GoldSHZ) Init() {
	if self.Person.Win > 0 {
		self.Person.Total += self.Person.Win
	}
	self.State = 0
	self.Person.Win = 0
	self.Person.Cost = 0
	self.Person.RunWin = 0
	self.Bet = 0
	self.LineNum = 0
	self.RunNum = 0
	self.CurRunNum = 0
}

func (self *Game_GoldSHZ) OnInit(room *Room) {
	self.room = room
}

func (self *Game_GoldSHZ) OnSendInfo(person *Person) {
	if self.Person != nil && self.Person.Uid == person.Uid {
		//		self.Init()
		self.Person.SynchroGold(person.Gold)
		person.SendMsg("gamegoldshzinfo", self.getInfo(person.Uid))
		return
	}

	_person := new(Game_GoldSHZ_Person)
	_person.Uid = person.Uid
	_person.Gold = person.Gold
	_person.Total = person.Gold
	_person.Name = person.Name
	_person.Sex = person.Sex
	_person.IP = person.ip
	_person.Head = person.Imgurl
	_person.Address = person.minfo.Address
	self.Person = _person
	self.Init()
	person.SendMsg("gamegoldshzinfo", self.getInfo(person.Uid))

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "加入房间 state : ", self.State)

	self.SetTime(GOLDSHZ_TIME)

}

func (self *Game_GoldSHZ) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldSHZ) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold": //! 同步金币
		if self.Person.Uid == msg.V.(*staticfunc.Msg_SynchroGold).Uid {
			self.Person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal(self.Person.Uid, self.Person.Total)
		}
	case "gameshzstart": //! 开始游戏
		self.GameStart(msg.Uid, msg.V.(*Msg_GameStart).LineNum, msg.V.(*Msg_GameStart).Gold)
	case "gamestep": //! 摇色子玩法 0-不摇 1-普通加倍 2-双倍加倍
		self.GameYSZ(msg.Uid, msg.V.(*Msg_GameStep).Card)
	case "gamebets": //摇色子压大小 0-小 1-和 2-大  3-退出得分
		self.GameYSZResult(msg.Uid, msg.V.(*Msg_GameBets).Bets)
	}
}

func (self *Game_GoldSHZ) OnBye() {

}

func (self *Game_GoldSHZ) OnExit(uid int64) {
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

func (self *Game_GoldSHZ) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_GoldSHZ) OnIsBets(uid int64) bool {
	if self.State != 0 {
		return true
	}
	return false
}

func (self *Game_GoldSHZ) OnBalance() {
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

func (self *Game_GoldSHZ) OnTime() {
	if self.Person == nil {
		return
	}
	if self.Time == 0 {
		return
	}
	if time.Now().Unix() >= self.Time {
		if self.State == 0 {
			self.room.KickViewByUid(self.Person.Uid, 96)
		} else if self.State == 1 {
			self.GameYSZ(self.Person.Uid, 0)
		} else if self.State == 2 {
			self.GameYSZResult(self.Person.Uid, 3)
		}
	}
}
