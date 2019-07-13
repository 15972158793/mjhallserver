package gameserver

import (
	//"fmt"
	"lib"
	"math"
	//"sort"
	"staticfunc"
	"time"
)

type Best_Result struct {
	Result int
	Score  int
	Card   []int
}

var FKFPJ_Score [3][lib.TYPE_FPJ_MAX]int = [3][lib.TYPE_FPJ_MAX]int{{0, 1, 2, 3, 4, 6, 8, 40, 100, 150, 300},
	{0, 1, 2, 3, 5, 7, 10, 50, 120, 200, 400},
	{0, 1, 2, 3, 6, 8, 12, 60, 140, 250, 500}}

type Msg_GameGoldFKFPJ_Info struct {
	Begin   bool                       `json:"begin"`   //! 是否开始
	Cur     int                        `json:"cur"`     //! 当前第几次翻牌
	Card    []int                      `json:"card"`    //! 当前牌
	Save    []int                      `json:"save"`    //! 当前保留
	Jackpot lib.Game_GoldFKFPJ_Jackpot `json:"jackpot"` //! 奖池
	BL      int                        `json:"bl"`      //! 当前比例
	Number  int                        `json:"number"`  //! 当前机号
	Index   int                        `json:"index"`   //! 中奖排数
	Score   int                        `json:"score"`   //! 当前分
	Select  int                        `json:"select"`  //! 选择分
	Total   int                        `json:"total"`   //! 当前金币
	Result  int                        `json:"result"`
	Mission int                        `json:"mission"`
}

type Msg_GameGoldFKFPJ_Step struct {
	Save []int `json:"save"`
}

type Msg_GameGoldFKFPJ_Result struct {
	Card    []int                      `json:"card"`   //! 结果
	Result  int                        `json:"result"` //! 结果
	Score   int                        `json:"score"`
	Jackpot lib.Game_GoldFKFPJ_Jackpot `json:"jackpot"` //! 奖池
	Used    []int                      `json:"used"`    //! 哪几张牌有用
}

type Msg_GameGoldFKFPJ_Select struct {
	Select int `json:"select"` //! 结果
}

type Msg_GameGoldFKFPJ_Get struct {
	Total int `json:"total"`
}

type Msg_GameGoldFKFPJ_Guess struct {
	GuessType int `json:"guesstype"`
	Guess     int `json:"guess"`
}

type Msg_GameGoldFKFPJ_GuessResult struct {
	Total   int                        `json:"total"`
	Score   int                        `json:"score"`
	Mission int                        `json:"mission"`
	Jackpot lib.Game_GoldFKFPJ_Jackpot `json:"jackpot"`
	Result  int                        `json:"result"`
	Card    int                        `json:"card"`
}

type Msg_Machine_List struct {
	Type int              `json:"type"`
	Info []*lib.FPJPerson `json:"info"`
}

type Msg_Machine_Del struct {
	Index int `json:"index"`
}

type Msg_Machine_Delete struct {
	Index []int `json:"index"`
}

type Msg_Machine_Keep struct {
	Index int `json:"index"`
}

type Msg_Client_HeadColor struct {
	TVer int `json:"tver"`
	YVer int `json:"yver"`
}

/////////////////////
type Game_GoldFKFPJ struct {
	Uid      int64                      `json:"uid"` //! 玩家uid
	Name     string                     `json:"name"`
	Head     string                     `json:"head"`
	Gold     int                        `json:"gold"`    //! 进场金币
	Total    int                        `json:"total"`   //! 当前金币
	BL       int                        `json:"bl"`      //! 当前场次比例
	Number   int                        `json:"number"`  //! 当前机号
	Index    int                        `json:"index"`   //! 第几排中奖
	Jackpot  lib.Game_GoldFKFPJ_Jackpot `json:"jackpot"` //! 奖池
	Mission  int                        `json:"mission"` //! 闯关奖池
	Cur      int                        `json:"cur"`     //! 当前第几次
	Card     []int                      `json:"card"`    //! 当前牌
	TmpCard  []int                      `json:"tmpcard"` //! 可能出现的牌
	Save     []int                      `json:"save"`    //! 当前保留
	Score    int                        `json:"score"`   //! 当前分
	Select   int                        `json:"select"`  //! 选择分
	Result   int                        `json:"result"`
	Mgr      *CardMgr                   `json:"mgr"`
	DealWin  int                        `json:"dealwin"` //! 庄家赢
	KickTime int64                      `json:"ticktime"`

	room *Room
}

//! 同步金币
func (self *Game_GoldFKFPJ) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

func NewGame_GoldFKFPJ() *Game_GoldFKFPJ {
	game := new(Game_GoldFKFPJ)

	return game
}

func (self *Game_GoldFKFPJ) OnInit(room *Room) {
	self.room = room

	//! 不同场次的比例
	self.BL = staticfunc.GetCsvMgr().GetDF(self.room.Type)

	self.Card = make([]int, 0)
	self.TmpCard = make([]int, 0)
	self.Save = make([]int, 0)

	//! 选10分
	self.Select = 10 * self.BL

	self.KickTime = time.Now().Unix() + 300
}

func (self *Game_GoldFKFPJ) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldFKFPJ) OnSendInfo(person *Person) {
	if self.Uid == 0 {
		self.Uid = person.Uid
		self.Gold = person.Gold
		self.Total = person.Gold
		self.Head = person.Imgurl
		self.Name = person.Name
	} else {
		self.SynchroGold(person.Gold)
		if self.room.Begin { //! 上一局未结束
			self.Reset()
		}
	}

	//! 第一个场直接匹配
	if self.room.Type%10 == 0 && self.Number == 0 {
		self.Number = lib.GetFPJMGr().HasKeep(self.room.Type%10, person.Uid, 0)
		if self.Number == 0 {
			self.Number = lib.GetFPJMGr().Up(self.room.Type%10, self.room.Id, person.Uid, person.Name, person.Imgurl, false)
		}
		self.Jackpot = lib.GetFPJMGr().Load(self.room.Type%10, self.Number)
		if len(self.Jackpot.Card) == 0 {
			self.Jackpot.Card = NewCard_FPJ(false).Deal(5)
		}
	}

	if self.Number == 0 {
		var msg Msg_Machine_List
		msg.Type = self.room.Type % 10
		msg.Info = lib.GetFPJMGr().GetInfo(self.room.Type % 10)
		person.SendMsg("gamegoldfkfpjlist", &msg)
		lib.GetFPJMGr().Enter(self.room.Type%10, self.Uid)
	} else {
		person.SendMsg("gamegoldfkfpjinfo", self.getInfo(person.Uid))
	}
}

func (self *Game_GoldFKFPJ) OnMsg(msg *RoomMsg) {
	self.KickTime = time.Now().Unix() + 300
	switch msg.Head {
	case "synchrogold": //! 同步金币
		if self.Uid == msg.V.(*staticfunc.Msg_SynchroGold).Uid {
			self.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal(self.Uid, self.Total)
		}
	case "gamebegin": //! 开始
		self.OnBegin()
	case "gamefpjstep":
		self.GameStep(msg.V.(*Msg_GameGoldFKFPJ_Step).Save)
	case "gamefpjselect":
		self.GameSelect(msg.V.(*Msg_GameGoldFKFPJ_Select).Select)
	case "gamefpjget":
		self.GameGet()
	case "gamefpjguess":
		self.GameGuess(msg.V.(*Msg_GameGoldFKFPJ_Guess).GuessType, msg.V.(*Msg_GameGoldFKFPJ_Guess).Guess)
	case "gamefpjup": //! 上机
		self.GameUp(msg.V.(*Msg_GameGoldFKFPJ_Select).Select)
	case "gamefpjquick":
		self.GameQuick()
	case "gamefpjkeep": //! 留机
		self.GameKeep()
	case "gamefpjunkeep": //! 留机
		self.GameUnKeep()
	case "gamefpjheadcolor": //! 头彩
		self.GameHeadColor(msg.V.(*Msg_Client_HeadColor).TVer, msg.V.(*Msg_Client_HeadColor).YVer)
	}
}

func (self *Game_GoldFKFPJ) OnBegin() {
	if self.room.IsBye() {
		return
	}

	if self.room.Begin {
		return
	}

	person := GetPersonMgr().GetPerson(self.Uid)
	if person == nil {
		return
	}

	if self.Total < self.Select {
		person.SendErr("金币不足,请前往充值")
		return
	}

	self.room.Begin = true

	//! 当前翻牌一次
	self.Cur = 1

	self.Mgr = NewCard_FPJ(true)
	self.Result = lib.GetFPJMGr().GetType(self.room.Type % 190000)
	self.TmpCard = lib.GetFPJMGr().GetCard(self.Result)
	if self.Result <= lib.TYPE_FPJ_TH {
		self.Card = self.TmpCard
		self.TmpCard = make([]int, 0)
	} else {
		mgr := NewCard_FPJ(true)
		for i := 0; i < len(self.TmpCard); i++ {
			mgr.DealCard(self.TmpCard[i])
		}
		index1 := lib.HF_GetRandom(5)
		index2 := lib.HF_GetRandom(5)
		for i := 0; i < len(self.TmpCard); i++ {
			if i == index1 || i == index2 {
				self.Card = append(self.Card, mgr.Deal(1)[0])
			} else {
				self.Card = append(self.Card, self.TmpCard[i])
			}
		}
	}
	self.Score = self.Select * FKFPJ_Score[self.Index][self.Result]
	for i := 0; i < len(self.Card); i++ {
		self.Mgr.DealCard(self.Card[i])
	}

	//! 哪些保留
	self.Save = lib.FPJSaveCard(self.Card)
	//! 当前中第几排
	if lib.HF_GetRandom(100) < 10 { //! 低概率
		self.Index = 2
	} else {
		self.Index = lib.HF_GetRandom(2)
	}
	//! 当前total
	self.Total -= self.Select
	self.DealWin += self.Select

	self.Jackpot.AddNum()

	person.SendMsg("gamegoldfkfpjbegin", self.getInfo(person.Uid))
}

func (self *Game_GoldFKFPJ) GameStep(save []int) {
	person := GetPersonMgr().GetPerson(self.Uid)
	if person == nil {
		return
	}

	if self.Cur >= 2 {
		return
	}

	if len(save) != 5 {
		return
	}

	isfind := true
	if len(self.TmpCard) != 0 { //! 有预先选好的牌
		for i := 0; i < 5; i++ {
			if save[i] == 1 && self.Card[i] != self.TmpCard[i] {
				isfind = false
				break
			}
		}
	} else {
		isfind = false
	}

	if isfind { //! 用预先选好的牌
		self.Score = self.Select * FKFPJ_Score[self.Index][self.Result]
		self.Result = self.Result
		self.Card = self.TmpCard
		self.TmpCard = make([]int, 0)
	} else {
		if save[0] == 0 && save[1] == 0 && save[2] == 0 && save[3] == 0 && save[4] == 0 { //! 没有保留任何一个
			self.Result = lib.GetFPJMGr().GetType(self.room.Type % 190000)
			self.Card = lib.GetFPJMGr().GetCard(self.Result)
			self.Score = self.Select * FKFPJ_Score[self.Index][self.Result]
		} else if save[0] == 1 && save[1] == 1 && save[2] == 1 && save[3] == 1 && save[4] == 1 { //! 都保留
			if len(self.TmpCard) != 0 {
				self.Result = lib.FPJCardResult(self.Card)
				self.Score = self.Select * FKFPJ_Score[self.Index][self.Result]
			}
		} else {
			result := lib.GetFPJMGr().GetType(self.room.Type % 190000)
			find := false
			var best = Best_Result{0, 1000000000, []int{}}
			lst := make([]Best_Result, 0)
			for loop := 0; loop < 200; loop++ {
				mgr := new(CardMgr)
				lib.HF_DeepCopy(mgr, self.Mgr)
				for i := 0; i < len(self.Card); i++ {
					if save[i] == 0 {
						self.Card[i] = mgr.Deal(1)[0]
					}
				}
				self.Result = lib.FPJCardResult(self.Card)
				self.Score = self.Select * FKFPJ_Score[self.Index][self.Result]
				if self.Score < best.Score {
					best.Score = self.Score
					best.Result = self.Result
					lib.HF_DeepCopy(&best.Card, &self.Card)
				}
				if self.Result == result {
					find = true
					break
				}
				if self.Result > result && (self.Select >= self.Score || GetServer().FKFpjSysMoney[self.room.Type%190000]-int64(self.Score) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin) {
					find = true
					break
				}
				if self.Select >= self.Score || GetServer().FKFpjSysMoney[self.room.Type%190000]-int64(self.Score) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
					var node Best_Result
					node.Result = self.Result
					node.Score = self.Score
					lib.HF_DeepCopy(&node.Card, &self.Card)
					lst = append(lst, node)
				}
			}
			if !find {
				if len(lst) > 0 {
					var node = lst[lib.HF_GetRandom(len(lst))]
					self.Score = node.Score
					self.Result = node.Result
					self.Card = node.Card
				} else {
					self.Score = best.Score
					self.Result = best.Result
					self.Card = best.Card
				}
			}
		}
	}

	self.Notice(self.Result)

	if self.Result == lib.TYPE_FPJ_SM {
		self.Score += self.Jackpot.Jackpot[0] * self.BL
		self.Jackpot.Jackpot[0] = 200
		lib.GetFPJMGr().UpdHeadColor(person.Uid, person.Name, person.Imgurl, self.room.Type%190000, 3)
	} else if self.Result == lib.TYPE_FPJ_THS {
		self.Score += self.Jackpot.Jackpot[1] * self.BL
		self.Jackpot.Jackpot[1] = 500
		lib.GetFPJMGr().UpdHeadColor(person.Uid, person.Name, person.Imgurl, self.room.Type%190000, 2)
	} else if self.Result == lib.TYPE_FPJ_THDS {
		self.Score += self.Jackpot.Jackpot[2] * self.BL
		self.Jackpot.Jackpot[2] = 2000
		lib.GetFPJMGr().UpdHeadColor(person.Uid, person.Name, person.Imgurl, self.room.Type%190000, 1)
	} else if self.Result == lib.TYPE_FPJ_WM {
		self.Score += self.Jackpot.Jackpot[3] * self.BL
		self.Jackpot.Jackpot[3] = 5000
		lib.GetFPJMGr().UpdHeadColor(person.Uid, person.Name, person.Imgurl, self.room.Type%190000, 0)
	}
	self.Cur++
	self.Jackpot.Record[self.Result]++

	var msg Msg_GameGoldFKFPJ_Result
	msg.Card = self.Card
	msg.Result = self.Result
	msg.Score = self.Score
	msg.Jackpot = self.Jackpot
	if self.Result != 0 {
		msg.Used = lib.FPJUseCard(self.Result, self.Card)
	}
	self.room.SendMsg(self.Uid, "gamegoldfkfpjresult", &msg)

	if self.Result == 0 { //! 没中到end流程
		self.OnEnd()
		return
	}
}

func (self *Game_GoldFKFPJ) GameSelect(_select int) {
	if self.Uid == 0 {
		return
	}

	if _select <= 0 {
		return
	}

	self.Select = _select * self.BL
}

func (self *Game_GoldFKFPJ) GameGet() {
	if !self.room.Begin {
		return
	}

	if self.Uid == 0 {
		return
	}

	if self.Score == 0 {
		return
	}

	self.Total += self.Score
	self.DealWin -= self.Score
	self.OnEnd()

	var msg Msg_GameGoldFKFPJ_Get
	msg.Total = self.Total
	self.room.SendMsg(self.Uid, "gamegoldfkfpjget", &msg)
}

//! 快速上机
func (self *Game_GoldFKFPJ) GameQuick() {
	if self.Uid == 0 {
		return
	}

	person := GetPersonMgr().GetPerson(self.Uid)
	if person == nil {
		return
	}

	self.Number = lib.GetFPJMGr().HasKeep(self.room.Type%10, person.Uid, 0)
	if self.Number == 0 {
		self.Number = lib.GetFPJMGr().Up(self.room.Type%10, self.room.Id, person.Uid, person.Name, person.Imgurl, false)
	} else {
		if !lib.GetFPJMGr().UpIndex(self.room.Type%10, self.Number, self.room.Id, person.Uid, person.Name, person.Imgurl, false) {
			person.SendErr("该位置已被其他玩家上机,请选择其他机器")
			return
		}
	}

	self.Jackpot = lib.GetFPJMGr().Load(self.room.Type%10, self.Number)
	if len(self.Jackpot.Card) == 0 {
		self.Jackpot.Card = NewCard_FPJ(false).Deal(5)
	}

	person.SendMsg("gamegoldfkfpjinfo", self.getInfo(person.Uid))
	lib.GetFPJMGr().Quit(self.room.Type%10, person.Uid)

	//! 广播给其他人上机
	var msg lib.Msg_Machine_Add
	msg.Info = &lib.FPJPerson{self.Number, self.room.Id, person.Uid, person.Name, person.Imgurl, 0, false}
	self.BroadCastMsg("gamegoldfkfpjadd", &msg)
}

//! 留机
func (self *Game_GoldFKFPJ) GameKeep() {
	if self.Uid == 0 {
		return
	}

	person := GetPersonMgr().GetPerson(self.Uid)
	if person == nil {
		return
	}

	lib.GetFPJMGr().Keep(self.room.Type%10, self.Number)

	//! 广播给其他人留机
	var msg Msg_Machine_Keep
	msg.Index = self.Number
	self.BroadCastMsg("gamegoldfkfpjkeep", &msg)
}

func (self *Game_GoldFKFPJ) GameUnKeep() {
	if self.Uid == 0 {
		return
	}

	person := GetPersonMgr().GetPerson(self.Uid)
	if person == nil {
		return
	}

	//! 广播给其他人留机
	var msg Msg_Machine_Del
	msg.Index = lib.GetFPJMGr().UnKeep(self.room.Type%10, person.Uid)
	self.BroadCastMsg("gamegoldfkfpjdel", &msg)
}

func (self *Game_GoldFKFPJ) GameHeadColor(tver int, yver int) {
	if self.Uid == 0 {
		return
	}

	person := GetPersonMgr().GetPerson(self.Uid)
	if person == nil {
		return
	}

	self.room.SendMsg(self.Uid, "gamegoldfkfpjheadcolor", lib.GetFPJMGr().GetHeadColor(self.Uid, tver, yver))
}

func (self *Game_GoldFKFPJ) GameUp(index int) {
	if self.room.Begin {
		return
	}

	if self.Number != 0 {
		return
	}

	if self.Uid == 0 {
		return
	}

	if index < 1 || index > 64 {
		return
	}

	person := GetPersonMgr().GetPerson(self.Uid)
	if person == nil {
		return
	}

	if lib.GetFPJMGr().HasKeep(self.room.Type%10, person.Uid, index) > 0 {
		person.SendErr("您已在其他机器留机")
		return
	}

	if !lib.GetFPJMGr().UpIndex(self.room.Type%10, index, self.room.Id, person.Uid, person.Name, person.Imgurl, false) {
		person.SendErr("该位置已被其他玩家上机,请选择其他机器")
		return
	}

	self.Number = index
	self.Jackpot = lib.GetFPJMGr().Load(self.room.Type%10, self.Number)
	if len(self.Jackpot.Card) == 0 {
		self.Jackpot.Card = NewCard_FPJ(false).Deal(5)
	}

	person.SendMsg("gamegoldfkfpjinfo", self.getInfo(person.Uid))
	lib.GetFPJMGr().Quit(self.room.Type%10, person.Uid)

	//! 广播给其他人上机
	var msg lib.Msg_Machine_Add
	msg.Info = &lib.FPJPerson{index, self.room.Id, person.Uid, person.Name, person.Imgurl, 0, false}
	self.BroadCastMsg("gamegoldfkfpjadd", &msg)
}

func (self *Game_GoldFKFPJ) GameGuess(guesstype int, guess int) {
	if !self.room.Begin {
		return
	}

	if self.Uid == 0 {
		return
	}

	if guesstype < 1 || guesstype > 3 {
		return
	}

	if guesstype == 1 && self.Total < self.Score { //! 比双倍但是不够
		return
	}

	if guesstype == 3 && self.Score/self.BL < 20 { //! 半比倍但是不能
		return
	}

	if guesstype == 1 {
		self.Total -= self.Score
		self.DealWin += self.Score
		self.Score *= 2
	} else if guesstype == 3 {
		self.Total += self.Score / 2
		self.DealWin -= self.Score / 2
		self.Score /= 2
	}

	mgr := NewCard_FPJ(false)
	for i := 0; i < len(self.Jackpot.Card); i++ {
		mgr.DealCard(self.Jackpot.Card[i])
	}

	lost := false
	if self.Mission >= 40 { //! 40关必输
		lost = true
	}
	if GetServer().FKFpjSysMoney[self.room.Type%190000]-int64(2*self.Score) < lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
		lost = true
	}

	card := mgr.Deal(1)[0]

	result := 1 //! 0输  1赢  2平
	for result > 0 {
		if card/10 == 7 { //! 开7
			result = 2
		} else if card/10 < 7 { //! 开小
			if guess == 0 {
				result = 0
			} else {
				result = 1
			}
		} else {
			if guess == 0 {
				result = 1
			} else {
				result = 0
			}
		}
		if result == 0 {
			break
		}
		if lost {
			card = mgr.Deal(1)[0]
		} else {
			break
		}
	}

	self.Jackpot.Card = append(self.Jackpot.Card, card)
	if len(self.Jackpot.Card) > 8 {
		self.Jackpot.Card = self.Jackpot.Card[1:]
	}

	if result == 0 { //! 输了
		self.Jackpot.BBJackpot += self.Score / self.BL / 10
		if self.Jackpot.BBJackpot >= 500 {
			self.Total += self.Jackpot.BBJackpot * self.BL
			self.DealWin -= self.Jackpot.BBJackpot * self.BL
			self.Jackpot.BBJackpot = 0
		}
		self.OnEnd()
	} else { //! 赢了
		self.Mission++
		if result == 1 {
			self.Score *= 2
		}
	}

	var msg Msg_GameGoldFKFPJ_GuessResult
	msg.Total = self.Total
	msg.Score = self.Score
	msg.Mission = self.Mission
	msg.Jackpot = self.Jackpot
	msg.Result = result
	msg.Card = card
	self.room.SendMsg(self.Uid, "gamegoldfkfpjguess", &msg)
}

//! 结算
func (self *Game_GoldFKFPJ) OnEnd() {
	self.room.Begin = false

	self.Score = 0
	self.Index = 0
	self.Cur = 0
	self.Save = make([]int, 0)
	self.Card = make([]int, 0)
	self.Result = 0
	self.Mission = 0

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "战绩")
	var record staticfunc.Rec_Gold_Info
	record.Time = time.Now().Unix()
	record.GameType = self.room.Type
	record.Info = append(record.Info, staticfunc.Son_Rec_Gold_Person{self.Uid, self.Name, self.Head, -self.DealWin, false})
	GetServer().InsertRecord(self.room.Type, self.Uid, lib.HF_JtoA(&record), -self.DealWin)

	//!
	if self.DealWin > 0 {
		cost := int(math.Ceil(float64(self.DealWin) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		self.DealWin -= cost
		GetServer().SqlAgentGoldLog(self.Uid, cost, self.room.Type)
		GetServer().SqlAgentBillsLog(self.Uid, cost/2, self.room.Type)
		GetServer().SqlBZWLog(&SQL_BZWLog{1, self.DealWin, time.Now().Unix(), self.room.Type})
		cost = int(math.Ceil(float64(self.DealWin) * lib.GetManyMgr().GetProperty(self.room.Type).DealCost / 100.0))
		self.DealWin -= cost
	} else if self.DealWin < 0 {
		GetServer().SqlBZWLog(&SQL_BZWLog{1, self.DealWin, time.Now().Unix(), self.room.Type})
		cost := int(math.Ceil(float64(self.DealWin) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 200.0))
		GetServer().SqlAgentBillsLog(self.Uid, cost, self.room.Type)
	}

	GetServer().SetFKFpjSysMoney(self.room.Type%190000, GetServer().FKFpjSysMoney[self.room.Type%190000]+int64(self.DealWin))
	self.DealWin = 0
}

//! 重置
func (self *Game_GoldFKFPJ) Reset() {
	self.Total += self.Score
	self.DealWin -= self.Score
	self.OnEnd()
}

func (self *Game_GoldFKFPJ) OnBye() {
}

func (self *Game_GoldFKFPJ) OnExit(uid int64) {
	if self.Uid != 0 {
		if self.room.Begin { //! 上一局未结束
			self.Reset()
		}
		//! 退出房间同步金币
		gold := self.Total - self.Gold
		if gold > 0 {
			GetRoomMgr().AddCard(self.Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
		} else if gold < 0 {
			GetRoomMgr().CostCard(self.Uid, staticfunc.TYPE_GOLD, -gold, self.room)
		}
		self.Gold = self.Total
		self.Uid = 0
	}
}

func (self *Game_GoldFKFPJ) getInfo(uid int64) *Msg_GameGoldFKFPJ_Info {
	var msg Msg_GameGoldFKFPJ_Info
	msg.Begin = self.room.Begin
	msg.Cur = self.Cur
	msg.Card = self.Card
	msg.Save = self.Save
	msg.Jackpot = self.Jackpot
	msg.BL = self.BL
	msg.Number = self.Number
	msg.Index = self.Index
	msg.Score = self.Score
	msg.Select = self.Select
	msg.Total = self.Total
	msg.Result = self.Result
	msg.Mission = self.Mission

	return &msg
}

func (self *Game_GoldFKFPJ) OnTime() {
	if time.Now().Unix() >= self.KickTime {
		self.room.KickViewByUid(self.Uid, 97)
	}
}

func (self *Game_GoldFKFPJ) OnIsDealer(uid int64) bool {
	return false
}

//! 同步总分
func (self *Game_GoldFKFPJ) SendTotal(uid int64, total int) {
	var msg Msg_GameGoldBZW_Total
	msg.Uid = uid
	msg.Total = total
	self.room.SendMsg(uid, "gamegoldtotal", &msg)
}

//! 是否下注了
func (self *Game_GoldFKFPJ) OnIsBets(uid int64) bool {
	return false
}

//! 结算所有人
func (self *Game_GoldFKFPJ) OnBalance() {
	lib.GetFPJMGr().Quit(self.room.Type%10, self.Uid)
	lib.GetFPJMGr().Save(self.room.Type%10, self.Number, self.Jackpot)
	if lib.GetFPJMGr().DownFromIndex(self.room.Type%10, self.Number, false) {
		var msg Msg_Machine_Del
		msg.Index = self.Number
		self.BroadCastMsg("gamegoldfkfpjdel", &msg)
	}

	if self.Uid != 0 {
		if self.room.Begin { //! 上一局未结束
			self.Reset()
		}
		//! 退出房间同步金币
		gold := self.Total - self.Gold
		if gold > 0 {
			GetRoomMgr().AddCard(self.Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
		} else if gold < 0 {
			GetRoomMgr().CostCard(self.Uid, staticfunc.TYPE_GOLD, -gold, self.room)
		}
		self.Gold = self.Total
	}
}

func (self *Game_GoldFKFPJ) BroadCastMsg(head string, v interface{}) {
	msg := lib.HF_EncodeMsg(head, v, true)

	lst := lib.GetFPJMGr().GetPerson(self.room.Type % 10)
	for i := 0; i < len(lst); i++ {
		person := GetPersonMgr().GetPerson(lst[i])
		if person == nil {
			continue
		}
		person.SendByteMsg(msg)
	}
}

//! 公告
func (self *Game_GoldFKFPJ) Notice(result int) {
	notice := lib.GetFPJMGr().Notice(result, self.Name, self.Number, self.room.Type%10)
	if notice != "" {
		GetServer().SendNotice(notice)
	}
}
