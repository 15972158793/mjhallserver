package gameserver

import (
	"fmt"
	"lib"
	"math"
	"staticfunc"
	"time"
)

var GOLDYSZ_BS []int = []int{2, 2, 10, 5, 5, 5, 5, 5, 5, 5, 5, 5}
var GOLDYSZ_DS []int = []int{0, 0, 0, 3, 4, 5, 6, 7, 8, 9, 10, 11}

//! 金币场记录
type Rec_YSZ_Info struct {
	GameType int                  `json:"gametype"`
	Time     int64                `json:"time"` //! 记录时间
	Info     []Son_Rec_YSZ_Person `json:"info"`
}
type Son_Rec_YSZ_Person struct {
	Uid    int64   `json:"uid"`
	Name   string  `json:"name"`
	Head   string  `json:"head"`
	Score  int     `json:"score"`
	Result [2]int  `json:"result"`
	Bets   [12]int `json:"bets"`
}

type Msg_GameGoldYSZ_Info struct {
	Begin   bool                    `json:"begin"`  //! 是否开始
	Time    int64                   `json:"time"`   //! 倒计时
	Seat    [8]Son_GameGoldBZW_Info `json:"info"`   //! 8个位置
	Bets    [12]int                 `json:"bets"`   //! 12个下注
	Dealer  Son_GameGoldBZW_Info    `json:"dealer"` //! 庄家
	Total   int                     `json:"total"`  //! 自己的钱
	Trend   [][2]int                `json:"trend"`  //! 走势
	IsDeal  bool                    `json:"isdeal"` //! 是否可下庄
	Money   []int                   `json:"money"`
	BetTime int                     `json:"bettime"`
}

type Msg_GameGoldYSZ_Begin struct {
	Result [2]int `json:"result"`
}

type Msg_GameGoldYSZ_End struct {
	Uid     int64  `json:"uid"` //! 大赢家
	Name    string `json:"name"`
	Head    string `json:"head"`
	Result  [2]int `json:"result"`
	Money   []int  `json:"money"`
	BetTime int    `json:"bettime"`
}

//! 摇色子续压
type Msg_GameGoldYSZ_Goon struct {
	Uid   int64   `json:"uid"`
	Gold  [12]int `json:"gold"`
	Total int     `json:"total"`
}

///////////////////////////////////////////////////////
type Game_GoldYSZ_Person struct {
	Uid       int64   `json:"uid"`
	Gold      int     `json:"gold"`      //! 进来时候的钱
	Total     int     `json:"total"`     //! 当前的钱
	Win       int     `json:"win"`       //! 本局赢的钱
	Cost      int     `json:"cost"`      //! 手续费
	Bets      int     `json:"bets"`      //! 本局下了多少钱
	BetInfo   [12]int `json:"bets"`      //! 本局下的注
	BeBets    int     `json:"bebets"`    //! 上把下了多少钱
	BeBetInfo [12]int `json:"bebetinfo"` //! 上把的下注
	Name      string  `json:"name"`      //! 名字
	Head      string  `json:"head"`      //! 头像
	Online    bool    `json:"online"`
	Round     int     `json:"round"` //! 不下注轮数
	Seat      int     `json:"seat"`  //! 0-7有座  -1无座  100庄家
	IP        string  `json:"ip"`
	Address   string  `json:"address"`
	Sex       int     `json:"sex"`
}

//! 同步金币
func (self *Game_GoldYSZ_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

type Game_GoldYSZ struct {
	PersonMgr map[int64]*Game_GoldYSZ_Person   `json:"personmgr"`
	Bets      [12]map[*Game_GoldYSZ_Person]int `json:"bets"`
	Result    [2]int                           `json:"result"`
	Dealer    *Game_GoldYSZ_Person             `json:"dealer"`  //! 庄家
	Round     int                              `json:"round"`   //! 连庄轮数
	DownUid   int64                            `json:"downuid"` //! 下庄的人
	Time      int64                            `json:"time"`
	LstDeal   []*Game_GoldYSZ_Person           `json:"lstdeal"` //! 上庄列表
	Seat      [8]*Game_GoldYSZ_Person          `json:"seat"`    //! 8个位置
	Total     int                              `json:"total"`   //! 这局一共下了多少钱
	Money     int                              `json:"money"`   //! 系统庄的钱
	Trend     [][2]int                         `json:"trend"`   //! 走势
	BetTime   int                              `json:"bettime"`

	room *Room
}

func NewGame_GoldYSZ() *Game_GoldYSZ {
	game := new(Game_GoldYSZ)
	game.PersonMgr = make(map[int64]*Game_GoldYSZ_Person)
	for i := 0; i < len(game.Bets); i++ {
		game.Bets[i] = make(map[*Game_GoldYSZ_Person]int)
	}
	for i := 0; i < 20; i++ {
		game.Trend = append(game.Trend, [2]int{lib.HF_GetRandom(6) + 1, lib.HF_GetRandom(6) + 1})
	}

	return game
}

func (self *Game_GoldYSZ) OnInit(room *Room) {
	self.room = room
	self.Money = lib.GetManyMgr().GetProperty(self.room.Type).DealInitMoney
	if lib.GetManyMgr().GetProperty(self.room.Type).DealChange == 1 {
		self.Money += int(GetServer().YSZSYSMoney[self.room.Type%120000])
	}
	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 11
}

func (self *Game_GoldYSZ) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldYSZ) OnSendInfo(person *Person) {
	if self.Time == 0 {
		self.SetTime(lib.GetManyMgr().GetProperty(self.room.Type).BetTime)
	}

	//! 观众模式游戏,观众进来只发送游戏信息
	value, ok := self.PersonMgr[person.Uid]
	if ok {
		value.Online = true
		value.Round = 0
		value.IP = person.ip
		value.Address = person.minfo.Address
		value.Sex = person.Sex
		value.SynchroGold(person.Gold)
		person.SendMsg("gamegoldyszinfo", self.getInfo(person.Uid, value.Total))
		return
	}

	_person := new(Game_GoldYSZ_Person)
	_person.Uid = person.Uid
	_person.Seat = -1
	_person.Gold = person.Gold
	_person.Total = person.Gold
	_person.Name = person.Name
	_person.Head = person.Imgurl
	_person.IP = person.ip
	_person.Address = person.minfo.Address
	_person.Sex = person.Sex
	_person.Online = true
	self.PersonMgr[person.Uid] = _person
	person.SendMsg("gamegoldyszinfo", self.getInfo(person.Uid, person.Gold))
}

func (self *Game_GoldYSZ) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold": //! 同步金币
		person := self.GetPerson(msg.V.(*staticfunc.Msg_SynchroGold).Uid)
		if person != nil {
			person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal(person.Uid, person.Total)
		}
	case "gamebzwbets": //! 下注
		self.GameBets(msg.Uid, msg.V.(*Msg_GameGoldBZW_Bets).Index, msg.V.(*Msg_GameGoldBZW_Bets).Gold)
	case "gamebzwgoon":
		self.GameGoOn(msg.Uid)
	case "gamerob": //! 上庄
		self.GameUpDeal(msg.Uid)
	case "gameredeal": //! 下庄
		self.GameReDeal(msg.Uid)
	case "gamebzwseat":
		self.GameSeat(msg.Uid, msg.V.(*Msg_GameGoldBZW_Seat).Index)
	case "gameplayerlist":
		self.GamePlayerList(msg.Uid)
	}
}

func (self *Game_GoldYSZ) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.Begin = true
	if self.Dealer != nil { //! 当前是玩家庄
		lst := make([]GameGoldBZW_CanResult, 0)
		winlst := make([]GameGoldBZW_CanResult, 0)
		lostlst := make([]GameGoldBZW_CanResult, 0)
		for i := 2; i <= 12; i++ {
			if i != 2 && i != 12 { //! 2和12必然是豹子
				win := self.GetDealWin(i, false)
				if GetServer().YSZUSRMoney[self.room.Type%120000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMin && GetServer().YSZUSRMoney[self.room.Type%120000]+int64(win) <= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMax {
					lst = append(lst, GameGoldBZW_CanResult{i, false})
				}
				if win > 0 {
					winlst = append(winlst, GameGoldBZW_CanResult{i, false})
				} else {
					lostlst = append(lostlst, GameGoldBZW_CanResult{i, false})
				}
			}
		}
		if len(lst) == 0 || lib.HF_GetRandom(100) < 40 {
			for i := 2; i <= 12; i++ {
				if i%2 == 0 { //! 可能是豹子
					win := self.GetDealWin(i, true)
					if GetServer().YSZUSRMoney[self.room.Type%120000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMin && GetServer().YSZUSRMoney[self.room.Type%120000]+int64(win) <= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMax {
						lst = append(lst, GameGoldBZW_CanResult{i, false})
					}
					if win > 0 {
						winlst = append(winlst, GameGoldBZW_CanResult{i, true})
					} else {
						lostlst = append(lostlst, GameGoldBZW_CanResult{i, true})
					}
				}
			}
		}
		if len(lst) == 0 {
			if GetServer().YSZUSRMoney[self.room.Type%120000] >= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMax && len(lostlst) > 0 { //! 一定输
				result := lostlst[lib.HF_GetRandom(len(lostlst))]
				self.GetResult(result.DS, result.BZ)
			} else if GetServer().YSZUSRMoney[self.room.Type%120000] <= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMin && len(winlst) > 0 {
				result := winlst[lib.HF_GetRandom(len(winlst))]
				self.GetResult(result.DS, result.BZ)
			} else {
				lib.GetLogMgr().Output(lib.LOG_ERROR, "摇色子亏了")
				for i := 0; i < len(self.Result); i++ {
					self.Result[i] = lib.HF_GetRandom(6) + 1
				}
			}
		} else {
			result := lst[lib.HF_GetRandom(len(lst))]
			self.GetResult(result.DS, result.BZ)
		}
	} else {
		lst := make([]GameGoldBZW_CanResult, 0)
		for i := 2; i <= 12; i++ {
			if i != 2 && i != 12 { //! 2和12必然是豹子
				win := self.GetDealWin(i, false)
				if GetServer().YSZSYSMoney[self.room.Type%120000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
					lst = append(lst, GameGoldBZW_CanResult{i, false})
				} else if win >= 0 {
					lst = append(lst, GameGoldBZW_CanResult{i, false})
				}
			}
		}
		if len(lst) == 0 || lib.HF_GetRandom(100) < 40 {
			for i := 2; i <= 12; i++ {
				if i%2 == 0 {
					win := self.GetDealWin(i, true)
					if GetServer().YSZSYSMoney[self.room.Type%120000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
						lst = append(lst, GameGoldBZW_CanResult{i, true})
					} else if win >= 0 {
						lst = append(lst, GameGoldBZW_CanResult{i, true})
					}
				}
			}
		}
		if len(lst) == 0 { //! 如果都输钱，就随机，理论上不会发生这种情况
			lib.GetLogMgr().Output(lib.LOG_ERROR, "豹子王亏了")
			for i := 0; i < len(self.Result); i++ {
				self.Result[i] = lib.HF_GetRandom(6) + 1
			}
		} else {
			result := lst[lib.HF_GetRandom(len(lst))]
			self.GetResult(result.DS, result.BZ)
		}
	}

	self.OnEnd()
}

//! 得到压中了这个点可以赢的钱
func (self *Game_GoldYSZ) GetDealWinByPos(pos int) int {
	lost := self.GetMoneyPos(pos) * GOLDYSZ_BS[pos] //! 压中这个点本身赢的钱
	maxlost := 0
	if pos == 0 {
		for i := 3; i <= 11; i++ {
			if i%2 == 0 {
				continue
			}
			tmp := self.GetMoneyPos(i) * GOLDYSZ_BS[i]
			if tmp > maxlost {
				maxlost = tmp
			}
		}
	} else if pos == 1 {
		for i := 3; i <= 11; i++ {
			if i%2 == 1 {
				continue
			}
			tmp := self.GetMoneyPos(i) * GOLDYSZ_BS[i]
			if tmp > maxlost {
				maxlost = tmp
			}
		}
	} else if pos == 2 {
		for i := 3; i <= 11; i++ {
			if i%2 == 1 {
				continue
			}
			tmp := self.GetMoneyPos(i) * GOLDYSZ_BS[i]
			if tmp > maxlost {
				maxlost = tmp
			}
		}
	} else {
		if pos%2 == 0 {
			maxlost = self.GetMoneyPos(2) * GOLDYSZ_BS[2]
		} else {
			maxlost = self.GetMoneyPos(0) * GOLDYSZ_BS[0]
		}
	}

	return lost + maxlost
}

//! 得到庄家可以赢的钱
func (self *Game_GoldYSZ) GetDealWin(ds int, bz bool) int {
	lost := 0
	if bz {
		lost += self.GetMoneyPos(2) * GOLDYSZ_BS[2]
	} else {
		if ds%2 == 1 {
			lost += self.GetMoneyPos(0) * GOLDYSZ_BS[0]
		} else if ds%2 == 0 {
			lost += self.GetMoneyPos(1) * GOLDYSZ_BS[1]
		}
	}

	for i := 3; i < 11; i++ {
		if ds == i {
			lost += self.GetMoneyPos(i) * GOLDYSZ_BS[i]
			break
		}
	}

	return self.Total - lost
}

//! 根据点数得到实际
func (self *Game_GoldYSZ) GetResult(ds int, bz bool) {
	if bz {
		self.Result[0] = ds / 2
		self.Result[1] = ds / 2
		return
	}

	//! 得到第一个数
	max := lib.HF_MinInt(6, ds-1)
	min := lib.HF_MaxInt(1, ds-6)
	self.Result[0] = lib.HF_GetRandom(max-min) + min
	//! 得到第二个数
	self.Result[1] = ds - self.Result[0]

	if self.Result[0] == self.Result[1] { //! 不能是豹子
		a := self.Result[0] - 1
		b := 6 - self.Result[1]
		c := lib.HF_MinInt(a, b)
		self.Result[0] -= c
		self.Result[1] += c
	}
}

func (self *Game_GoldYSZ) IsType() int { //! 0单 1双 2豹子
	if self.Result[0] == self.Result[1] {
		return 2
	}

	if (self.Result[0]+self.Result[1])%2 == 1 {
		return 0
	}

	return 1
}

func (self *Game_GoldYSZ) GameUpDeal(uid int64) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if person.Total < lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("金币必须大于%d才能上庄", lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("金币必须大于%d才能上庄", lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("金币必须大于%d万才能上庄", lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney/10000))
		}
		return
	}

	if self.Dealer == person {
		self.DownUid = 0
	} else {
		for i := 0; i < len(self.LstDeal); i++ {
			if self.LstDeal[i] == person {
				self.room.SendErr(uid, "您已经在上庄列表中，请等待上庄")
				return
			}
		}
		if len(self.LstDeal) == 0 {
			self.Round = 0
		}
		self.LstDeal = append(self.LstDeal, person)
	}
	person.Round = 0

	var msg Msg_GameGoldBZW_DealList
	msg.Type = 0
	msg.Info = make([]Son_GameGoldBZW_Info, 0)
	for i := 0; i < len(self.LstDeal); i++ {
		msg.Info = append(msg.Info, Son_GameGoldBZW_Info{self.LstDeal[i].Uid, self.LstDeal[i].Name, self.LstDeal[i].Head, self.LstDeal[i].Total, self.LstDeal[i].IP, self.LstDeal[i].Address, self.LstDeal[i].Sex})
	}
	self.room.SendMsg(uid, "gameyszdeal", &msg)
}

func (self *Game_GoldYSZ) GameReDeal(uid int64) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if self.Dealer == person { //! 正在庄
		self.DownUid = uid
		self.room.SendErr(uid, "您已成功下庄，请等待本局结束")
	} else {
		for i := 0; i < len(self.LstDeal); i++ {
			if self.LstDeal[i] == person {
				copy(self.LstDeal[i:], self.LstDeal[i+1:])
				self.LstDeal = self.LstDeal[:len(self.LstDeal)-1]
				break
			}
		}
	}

	var msg Msg_GameGoldBZW_DealList
	msg.Type = 1
	msg.Info = make([]Son_GameGoldBZW_Info, 0)
	self.room.SendMsg(uid, "gameyszdeal", &msg)
}

//! 坐下
func (self *Game_GoldYSZ) GameSeat(uid int64, index int) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if index < 0 || index > 7 {
		return
	}

	if self.Dealer == person {
		self.room.SendErr(uid, "庄家无法坐下")
		return
	}

	if person.Total < lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("金币必须大于%d才能坐下", lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("金币必须大于%d才能坐下", lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("金币必须大于%d万才能坐下", lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney/10000))
		}
		return
	}

	for i := 0; i < len(self.Seat); i++ {
		if self.Seat[i] == person {
			return
		}
	}

	if self.Seat[index] != nil {
		if person.Total <= self.Seat[index].Total {
			self.room.SendErr(uid, "该位置已经有人坐了")
			return
		}
		//! 把原来在这个位置上的人挤下去
		self.Seat[index].Seat = -1
	}

	self.Seat[index] = person
	person.Seat = index

	var msg Msg_GameGoldBZW_UpdSeat
	msg.Uid = uid
	msg.Index = index
	msg.Head = person.Head
	msg.Name = person.Name
	msg.Total = person.Total
	msg.IP = person.IP
	msg.Address = person.Address
	msg.Sex = person.Sex
	self.room.broadCastMsg("gameyszseat", &msg)
}

func (self *Game_GoldYSZ) GameBets(uid int64, index int, gold int) {
	if uid == 0 {
		return
	}

	if index < 0 || index > 11 {
		return
	}

	if gold <= 0 {
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() >= int64(self.BetTime-2) {
		self.room.SendErr(uid, "正在开奖,请稍后下注")
		return
	}

	if self.Dealer == nil && lib.GetManyMgr().GetProperty(self.room.Type).SysNoBets == 1 {
		self.room.SendErr(uid, "请等待玩家上庄")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if self.Dealer == person {
		self.room.SendErr(uid, "庄家不用下注")
		return
	}

	if person.Total < lib.GetManyMgr().GetProperty(self.room.Type).MinBet {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("%d金币才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("%d金币才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("%d万金币才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/10000))
		}
		return
	}

	if person.Total < gold {
		self.room.SendErr(uid, "您的金币不足，请前往充值。")
		return
	}

	if person.Bets+gold > lib.GetManyMgr().GetProperty(self.room.Type).MaxBet {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d。", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d。", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d万。", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet/10000))
		}
		return
	}

	//! 判断庄家是否够赔
	{
		//! 模拟总下注
		self.Total += gold
		self.Bets[index][person] += gold
		//! 得到庄家的钱
		dealmoney := 0
		if self.Dealer == nil {
			dealmoney = self.Money
		} else {
			dealmoney = self.Dealer.Total
		}
		//! 模拟庄家赢钱
		dealwin := self.GetDealWinByPos(index)
		if dealmoney+self.Total-dealwin < dealmoney/5 {
			self.Total -= gold
			self.Bets[index][person] -= gold
			self.room.SendErr(uid, "庄家金币不足，该位置无法下注。")
			return
		}
	}

	person.Bets += gold
	person.Total -= gold
	person.BetInfo[index] += gold
	person.Round = 0

	var msg Msg_GameGoldBZW_Bets
	msg.Uid = uid
	msg.Index = index
	msg.Gold = gold
	msg.Total = person.Total
	self.room.broadCastMsg("gameyszbets", &msg)
}

//! 续压
func (self *Game_GoldYSZ) GameGoOn(uid int64) {
	if uid == 0 {
		return
	}

	if self.Time != 0 && self.Time-time.Now().Unix() >= int64(self.BetTime-2) {
		self.room.SendErr(uid, "正在开奖,请稍后下注")
		return
	}

	if self.Dealer == nil && lib.GetManyMgr().GetProperty(self.room.Type).SysNoBets == 1 {
		self.room.SendErr(uid, "请等待玩家上庄")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if self.Dealer == person {
		self.room.SendErr(uid, "庄家不用下注")
		return
	}

	if person.BeBets == 0 {
		return
	}

	if person.Total < lib.GetManyMgr().GetProperty(self.room.Type).MinBet {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("%d金币才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("%d金币才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("%d万金币才能下注", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/10000))
		}
		return
	}

	if person.Total < person.BeBets {
		self.room.SendErr(uid, "您的金币不足，请前往充值。")
		return
	}

	if person.Bets+person.BeBets > lib.GetManyMgr().GetProperty(self.room.Type).MaxBet {
		if GetServer().Con.MoneyMode == 1 {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d。", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d。", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("单局下注不能超过%d万。", lib.GetManyMgr().GetProperty(self.room.Type).MaxBet/10000))
		}
		return
	}

	//! 判断庄家是否够赔
	{
		//! 模拟总下注
		for i := 0; i < len(person.BeBetInfo); i++ {
			self.Total += person.BeBetInfo[i]
			self.Bets[i][person] += person.BeBetInfo[i]
		}
		//! 得到庄家的钱
		dealmoney := 0
		if self.Dealer == nil {
			dealmoney = self.Money
		} else {
			dealmoney = self.Dealer.Total
		}
		//! 模拟庄家赢钱
		dealwin := 0
		for i := 0; i < len(person.BeBetInfo); i++ {
			if person.BeBetInfo[i] == 0 {
				continue
			}
			tmp := self.GetDealWinByPos(i)
			if tmp > dealwin {
				dealwin = tmp
			}
		}
		if dealmoney+self.Total-dealwin < dealmoney/5 {
			for i := 0; i < len(person.BeBetInfo); i++ {
				self.Total -= person.BeBetInfo[i]
				self.Bets[i][person] -= person.BeBetInfo[i]
			}
			self.room.SendErr(uid, "庄家金币不足，无法续压。")
			return
		}
	}

	person.Bets += person.BeBets
	person.Total -= person.BeBets
	for i := 0; i < len(person.BeBetInfo); i++ {
		person.BetInfo[i] += person.BeBetInfo[i]
	}
	person.Round = 0

	var msg Msg_GameGoldYSZ_Goon
	msg.Uid = uid
	msg.Gold = person.BeBetInfo
	msg.Total = person.Total
	self.room.broadCastMsg("gameyszgoon", &msg)
}

//! 结算
func (self *Game_GoldYSZ) OnEnd() {
	self.room.Begin = false

	ds := 0
	for i := 0; i < len(self.Result); i++ {
		ds += self.Result[i]
	}

	trend := self.IsType()
	tmp := [][2]int{self.Result}
	tmp = append(tmp, self.Trend...)
	if len(tmp) > 20 {
		tmp = tmp[0:20]
	}
	self.Trend = tmp

	dealwin := 0
	for i := 0; i < len(self.Bets); i++ {
		if i < 3 { //! 上面一排
			if trend == i { //!  下注赢了
				for key, value := range self.Bets[i] {
					winmoney := value * GOLDYSZ_BS[i]
					dealwin -= winmoney
					key.Win += winmoney
					key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
				}
			}
		} else {
			if self.Result[0]+self.Result[1] == GOLDYSZ_DS[i] { //! 下注赢了
				for key, value := range self.Bets[i] {
					winmoney := value * GOLDYSZ_BS[i]
					dealwin -= winmoney
					key.Win += winmoney
					key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
				}
			}
		}
	}

	var bigwin *Game_GoldYSZ_Person = nil
	for _, value := range self.PersonMgr {
		if value.Win > 0 {
			value.Win -= value.Cost
			GetServer().SqlAgentGoldLog(value.Uid, value.Cost, self.room.Type)
			GetServer().SqlAgentBillsLog(value.Uid, value.Cost/2, self.room.Type)
			value.Total += value.Win

			var msg Msg_GameGoldBZW_Balance
			msg.Uid = value.Uid
			msg.Total = value.Total
			msg.Win = value.Win
			find := false
			for j := 0; j < len(self.Seat); j++ {
				if self.Seat[j] == value {
					self.room.broadCastMsg("gamegoldyszbalance", &msg)
					find = true
					break
				}
			}
			if !find {
				self.room.SendMsg(value.Uid, "gamegoldyszbalance", &msg)
			}

			if bigwin == nil {
				bigwin = value
			} else if value.Win > bigwin.Win {
				bigwin = value
			}
		} else if value.Win-value.Bets < 0 {
			cost := int(math.Ceil(float64(value.Win-value.Bets) * float64(lib.GetManyMgr().GetProperty(self.room.Type).Cost) / 200.0))
			GetServer().SqlAgentBillsLog(value.Uid, cost, self.room.Type)
		}

		//! 插入战绩
		if value.Bets > 0 {
			var record Rec_YSZ_Info
			record.Time = time.Now().Unix()
			record.GameType = self.room.Type
			var rec Son_Rec_YSZ_Person
			rec.Uid = value.Uid
			rec.Name = value.Name
			rec.Head = value.Head
			rec.Score = value.Win - value.Bets
			rec.Result = self.Result
			rec.Bets = value.BetInfo
			record.Info = append(record.Info, rec)
			GetServer().InsertRecord(self.room.Type, value.Uid, lib.HF_JtoA(&record), rec.Score)
		}
	}

	dealwin = self.Total + dealwin
	if self.Dealer == nil && dealwin != 0 { //! 系统庄
		GetServer().SqlBZWLog(&SQL_BZWLog{1, dealwin, time.Now().Unix(), self.room.Type})
	}
	if self.Dealer != nil { //! 玩家庄并且是系统算法
		GetServer().SetYSZUSRMoney(self.room.Type%120000, GetServer().YSZUSRMoney[self.room.Type%120000]+int64(dealwin))
	}
	if dealwin > 0 {
		bl := lib.GetManyMgr().GetProperty(self.room.Type).Cost
		if self.Dealer == nil { //! 系统庄
			bl = float64(lib.GetManyMgr().GetProperty(self.room.Type).DealCost)
		}
		cost := int(math.Ceil(float64(dealwin) * bl / 100.0))
		dealwin -= cost
		if self.Dealer != nil { //! 赢了抽水
			GetServer().SqlAgentGoldLog(self.Dealer.Uid, cost, self.room.Type)
			GetServer().SqlAgentBillsLog(self.Dealer.Uid, cost/2, self.room.Type)
		}
	} else if dealwin < 0 && self.Dealer != nil {
		cost := int(math.Ceil(float64(dealwin) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 200.0))
		GetServer().SqlAgentBillsLog(self.Dealer.Uid, cost, self.room.Type)
	}
	if self.Dealer != nil {
		var record Rec_YSZ_Info
		record.Time = time.Now().Unix()
		record.GameType = self.room.Type
		var rec Son_Rec_YSZ_Person
		rec.Uid = self.Dealer.Uid
		rec.Name = self.Dealer.Name
		rec.Head = self.Dealer.Head
		rec.Score = dealwin
		rec.Result = self.Result
		record.Info = append(record.Info, rec)
		GetServer().InsertRecord(self.room.Type, self.Dealer.Uid, lib.HF_JtoA(&record), rec.Score)
	}

	//! 发送庄家结算
	{
		self.Money = lib.GetManyMgr().GetProperty(self.room.Type).DealInitMoney
		var msg Msg_GameGoldBZW_Balance
		if self.Dealer != nil {
			if self.Dealer.Total+dealwin > 0 {
				self.Dealer.Total += dealwin
			} else {
				self.Dealer.Total = 0
				dealwin = -self.Dealer.Total
			}
			msg.Uid = self.Dealer.Uid
			msg.Total = self.Dealer.Total
		} else {
			GetServer().SetYSZSYSMoney(self.room.Type%120000, GetServer().YSZSYSMoney[self.room.Type%120000]+int64(dealwin))
			if lib.GetManyMgr().GetProperty(self.room.Type).DealChange == 1 {
				self.Money += int(GetServer().YSZSYSMoney[self.room.Type%120000])
			}
			msg.Total = self.Money
		}
		msg.Win = dealwin
		self.room.broadCastMsg("gamegoldyszbalance", &msg)
	}

	//! 30秒的下注时间
	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 11
	self.SetTime(self.BetTime)

	//! 总结算
	{
		var msg Msg_GameGoldYSZ_End
		msg.Money = lib.GetManyMoneyMgr().GetProperty(self.room.Type).Money
		msg.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime
		msg.Result = self.Result
		if bigwin != nil {
			msg.Uid = bigwin.Uid
			msg.Name = bigwin.Name
			msg.Head = bigwin.Head
		}
		self.room.broadCastMsg("gamegoldyszend", &msg)
	}

	self.Total = 0

	//! 把不在room.uid里面的玩家清理出去
	for key, value := range self.PersonMgr {
		if value.Online {
			find := false
			for i := 0; i < len(self.LstDeal); i++ {
				if self.LstDeal[i] == value {
					find = true
					break
				}
			}
			if !find && value.Seat < 0 { //! 无座玩家不下注轮数++
				value.Round++
			}
			if value.Round >= 5 && GetPersonMgr().GetPerson(value.Uid) == nil {
				self.room.KickViewByUid(value.Uid, 96)
			} else {
				value.BeBets = value.Bets
				value.BeBetInfo = value.BetInfo
				value.Win = 0
				value.Cost = 0
				value.Bets = 0
				for j := 0; j < len(value.BetInfo); j++ {
					value.BetInfo[j] = 0
				}
				continue
			}
		}

		//! 走的人正在上庄列表
		for j := 0; j < len(self.LstDeal); j++ {
			if self.LstDeal[j] == value {
				copy(self.LstDeal[j:], self.LstDeal[j+1:])
				self.LstDeal = self.LstDeal[:len(self.LstDeal)-1]
				break
			}
		}

		//! 走的人是庄家
		if self.Dealer == value {
			self.ChageDeal()
		}

		//! 走的人是位置上面的人
		for j := 0; j < len(self.Seat); j++ {
			if self.Seat[j] == value {
				self.Seat[j] = nil
				var msg Msg_GameGoldBZW_UpdSeat
				msg.Index = j
				self.room.broadCastMsg("gameyszseat", &msg)
				break
			}
		}
		delete(self.PersonMgr, key)
	}

	for i := 0; i < len(self.Bets); i++ {
		self.Bets[i] = make(map[*Game_GoldYSZ_Person]int)
	}

	//! 判断庄家是否能继续连
	if self.Dealer != nil && (self.Dealer.Total < lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney || self.DownUid == self.Dealer.Uid || GetPersonMgr().GetPerson(self.Dealer.Uid) == nil) {
		self.ChageDeal()
	} else if self.Dealer == nil && len(self.LstDeal) > 0 {
		self.ChageDeal()
	} else if self.Dealer != nil {
		if self.Round >= 10 && len(self.LstDeal) > 0 {
			self.ChageDeal()
		} else {
			self.Round++
		}
	}

	//! 判断坐下的人是否能继续坐
	for i := 0; i < len(self.Seat); i++ {
		if self.Seat[i] == nil {
			continue
		}
		if self.Seat[i].Total < lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney {
			self.Seat[i].Seat = -1
			var msg Msg_GameGoldBZW_UpdSeat
			msg.Index = i
			self.room.broadCastMsg("gameyszseat", &msg)
			self.Seat[i] = nil
		}
	}
}

func (self *Game_GoldYSZ) OnBye() {
}

func (self *Game_GoldYSZ) OnExit(uid int64) {
	value, ok := self.PersonMgr[uid]
	if ok {
		value.Online = false
		//! 退出房间同步金币
		gold := value.Total - value.Gold
		if gold > 0 {
			GetRoomMgr().AddCard(value.Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
		} else if gold < 0 {
			GetRoomMgr().CostCard(value.Uid, staticfunc.TYPE_GOLD, -gold, self.room)
		}
		value.Gold = value.Total
	}
}

//! 得到这个位置下了多少钱
func (self *Game_GoldYSZ) GetMoneyPos(index int) int {
	total := 0
	for _, value := range self.Bets[index] {
		total += value
	}
	return total
}

func (self *Game_GoldYSZ) getInfo(uid int64, total int) *Msg_GameGoldYSZ_Info {
	var msg Msg_GameGoldYSZ_Info
	msg.Begin = self.room.Begin
	msg.Time = self.Time - time.Now().Unix()
	msg.Total = total
	msg.Trend = self.Trend
	msg.IsDeal = false
	msg.Money = lib.GetManyMoneyMgr().GetProperty(self.room.Type).Money
	msg.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime
	if self.Dealer != nil && self.Dealer.Uid == uid {
		msg.IsDeal = true
	} else {
		for i := 0; i < len(self.LstDeal); i++ {
			if self.LstDeal[i].Uid == uid {
				msg.IsDeal = true
				break
			}
		}
	}
	for i := 0; i < len(self.Bets); i++ {
		msg.Bets[i] = self.GetMoneyPos(i)
	}
	for i := 0; i < 8; i++ {
		if self.Seat[i] != nil {
			msg.Seat[i].Uid = self.Seat[i].Uid
			msg.Seat[i].Name = self.Seat[i].Name
			msg.Seat[i].Head = self.Seat[i].Head
			msg.Seat[i].Total = self.Seat[i].Total
			msg.Seat[i].Ip = self.Seat[i].IP
			msg.Seat[i].Address = self.Seat[i].Address
			msg.Seat[i].Sex = self.Seat[i].Sex
		}
	}
	if self.Dealer != nil {
		msg.Dealer.Uid = self.Dealer.Uid
		msg.Dealer.Name = self.Dealer.Name
		msg.Dealer.Head = self.Dealer.Head
		msg.Dealer.Total = self.Dealer.Total
		msg.Dealer.Ip = self.Dealer.IP
		msg.Dealer.Address = self.Dealer.Address
		msg.Dealer.Sex = self.Dealer.Sex
	} else {
		msg.Dealer.Total = self.Money
	}
	return &msg
}

func (self *Game_GoldYSZ) GetPerson(uid int64) *Game_GoldYSZ_Person {
	return self.PersonMgr[uid]
}

func (self *Game_GoldYSZ) OnTime() {
	if self.Time == 0 {
		return
	}

	if time.Now().Unix() < self.Time {
		return
	}

	if !self.room.Begin {
		self.OnBegin()
		return
	}
}

func (self *Game_GoldYSZ) OnIsDealer(uid int64) bool {
	if self.Dealer != nil && self.Dealer == self.GetPerson(uid) {
		return true
	}
	return false
}

//! 申请无座玩家
func (self *Game_GoldYSZ) GamePlayerList(uid int64) {
	has := false
	var msg Msg_GameGoldBZW_List
	for _, value := range self.PersonMgr {
		if value.Seat >= 0 {
			continue
		}

		var node Son_GameGoldBZW_Info
		node.Uid = value.Uid
		node.Name = value.Name
		node.Total = value.Total
		node.Head = value.Head
		msg.Info = append(msg.Info, node)
		if uid == value.Uid {
			has = true
		}
		if len(msg.Info) >= 30 {
			break
		}
	}
	if !has {
		person := self.GetPerson(uid)
		if person == nil {
			return
		}
		if person.Seat < 0 {
			msg.Info = append(msg.Info, Son_GameGoldBZW_Info{person.Uid, person.Name, person.Head, person.Total, person.IP, person.Address, person.Sex})
		}
	}
	self.room.SendMsg(uid, "gameplayerlist", &msg)
}

//! 同步总分
func (self *Game_GoldYSZ) SendTotal(uid int64, total int) {
	var msg Msg_GameGoldBZW_Total
	msg.Uid = uid
	msg.Total = total

	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if person.Seat < 0 { //! 不在座位上
		self.room.SendMsg(uid, "gamegoldtotal", &msg)
	} else {
		self.room.broadCastMsg("gamegoldtotal", &msg)
	}
}

//! 设置时间
func (self *Game_GoldYSZ) SetTime(t int) {
	if t == 0 {
		self.Time = 0
	} else {
		self.Time = time.Now().Unix() + int64(t)
	}

	var msg Msg_SetTime
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	self.room.broadCastMsg("gametime", &msg)
}

//! 换庄
func (self *Game_GoldYSZ) ChageDeal() {
	if self.Dealer != nil {
		self.Dealer.Seat = -1
	}

	self.Dealer = nil
	for len(self.LstDeal) > 0 {
		if self.LstDeal[0].Total >= lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney {
			self.Dealer = self.LstDeal[0]
			self.LstDeal = self.LstDeal[1:]
			self.Dealer.Seat = 100
			for i := 0; i < len(self.Seat); i++ {
				if self.Seat[i] == self.Dealer {
					var msg Msg_GameGoldBZW_UpdSeat
					msg.Index = i
					self.room.broadCastMsg("gameyszseat", &msg)
					self.Seat[i] = nil
					break
				}
			}
			break
		} else {
			self.LstDeal = self.LstDeal[1:]
		}
	}
	self.DownUid = 0
	self.Round = 0

	var msg Msg_GameGoldBZW_Deal
	if self.Dealer != nil {
		msg.Uid = self.Dealer.Uid
		msg.Name = self.Dealer.Name
		msg.Head = self.Dealer.Head
		msg.Total = self.Dealer.Total
		msg.IP = self.Dealer.IP
		msg.Address = self.Dealer.Address
		msg.Sex = self.Dealer.Sex
	} else {
		msg.Total = self.Money
	}

	self.room.broadCastMsg("gamerob", &msg)
}

//! 是否下注了
func (self *Game_GoldYSZ) OnIsBets(uid int64) bool {
	value, ok := self.PersonMgr[uid]
	if ok {
		return value.Bets > 0
	}
	return false
}

//! 结算所有人
func (self *Game_GoldYSZ) OnBalance() {
	for _, value := range self.PersonMgr {
		//! 被clear时先返还本场下注
		value.Total += value.Bets

		gold := value.Total - value.Gold
		if gold > 0 {
			GetRoomMgr().AddCard(value.Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
		} else if gold < 0 {
			GetRoomMgr().CostCard(value.Uid, staticfunc.TYPE_GOLD, -gold, self.room)
		}
		value.Gold = value.Total
	}
}
