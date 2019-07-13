package gameserver

import (
	"fmt"
	"lib"
	"math"
	"staticfunc"
	"time"
)

var GOLDBZW_BS []int = []int{2, 35, 2, 60, 30, 17, 12, 8, 7, 6, 6, 7, 8, 12, 17, 30, 60}
var GOLDBZW_DS []int = []int{0, 0, 0, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17}

//! 金币场记录
type Rec_BZW_Info struct {
	GameType int                  `json:"gametype"`
	Time     int64                `json:"time"` //! 记录时间
	Info     []Son_Rec_BZW_Person `json:"info"`
}
type Son_Rec_BZW_Person struct {
	Uid    int64   `json:"uid"`
	Name   string  `json:"name"`
	Head   string  `json:"head"`
	Score  int     `json:"score"`
	Result [3]int  `json:"result"`
	Bets   [17]int `json:"bets"`
}

type Msg_GameGoldBZW_Info struct {
	Begin   bool                    `json:"begin"`  //! 是否开始
	Time    int64                   `json:"time"`   //! 倒计时
	Seat    [8]Son_GameGoldBZW_Info `json:"info"`   //! 8个位置
	Bets    [17]int                 `json:"bets"`   //! 17个下注
	Dealer  Son_GameGoldBZW_Info    `json:"dealer"` //! 庄家
	Total   int                     `json:"total"`  //! 自己的钱
	Trend   [][3]int                `json:"trend"`  //! 走势
	IsDeal  bool                    `json:"isdeal"` //! 是否可下庄
	Change  bool                    `json:"change"` //! 是否可以打开超端
	Money   []int                   `json:"money"`  //! 筹码
	BetTime int                     `json:"bettime"`
}

type Son_GameGoldBZW_Info struct {
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Total   int    `json:"total"`
	Ip      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type Msg_GameGoldBZW_Begin struct {
	Result [3]int `json:"result"`
}

type Msg_GameGoldBZW_Balance struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"` //! 当前金币
	Win   int   `json:"win"`   //! 赢了多少金币
}

type Msg_GameGoldBZW_End struct {
	Uid     int64  `json:"uid"` //! 大赢家
	Name    string `json:"name"`
	Head    string `json:"head"`
	Result  [3]int `json:"result"`
	Money   []int  `json:"money"` //! 筹码
	BetTime int    `json:"bettime"`
}

type Msg_GameGoldBZW_Total struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"`
}

//! 豹子王下注
type Msg_GameGoldBZW_Bets struct {
	Uid   int64 `json:"uid"`
	Index int   `json:"index"`
	Gold  int   `json:"gold"`
	Total int   `json:"total"`
}

//! 豹子王续压
type Msg_GameGoldBZW_Goon struct {
	Uid   int64   `json:"uid"`
	Gold  [17]int `json:"gold"`
	Total int     `json:"total"`
}

//! 上庄
type Msg_GameGoldBZW_Deal struct {
	Uid     int64  `json:"uid"`
	Head    string `json:"head"`
	Name    string `json:"name"`
	Total   int    `json:"total"`
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

//! 换座位
type Msg_GameGoldBZW_Seat struct {
	Index int `json:"index"`
}

//! 刷新座位
type Msg_GameGoldBZW_UpdSeat struct {
	Index   int    `json:"index"`
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Total   int    `json:"total"`
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

//! 无座玩家列表
type Msg_GameGoldBZW_List struct {
	Info []Son_GameGoldBZW_Info `json:"info"`
}

//! 上庄和下庄
type Msg_GameGoldBZW_DealList struct {
	Type int                    `json:"type"` //! 0上庄  1下庄
	Info []Son_GameGoldBZW_Info `json:"info"`
}

//! 这局可能出现的点数
type GameGoldBZW_CanResult struct {
	DS int
	BZ bool
}

type GameGold_BigWin struct {
	Uid  int64
	Name string
	Head string
	Win  int
}

///////////////////////////////////////////////////////
type Game_GoldBZW_Person struct {
	Uid       int64   `json:"uid"`
	Gold      int     `json:"gold"`      //! 进来时候的钱
	Total     int     `json:"total"`     //! 当前的钱
	Win       int     `json:"win"`       //! 本局赢的钱
	Cost      int     `json:"cost"`      //! 手续费
	Bets      int     `json:"bets"`      //! 本局下了多少钱
	BetInfo   [17]int `json:"bets"`      //! 本局下的注
	BeBets    int     `json:"bebets"`    //! 上把下了多少钱
	BeBetInfo [17]int `json:"bebetinfo"` //! 上把的下注
	Name      string  `json:"name"`      //! 名字
	Head      string  `json:"head"`      //! 头像
	Online    bool    `json:"online"`
	Round     int     `json:"round"` //! 不下注轮数
	Seat      int     `json:"seat"`  //! 0-7有座  -1无座  100庄家
	IP        string  `json:"ip"`
	Address   string  `json:"address"`
	Sex       int     `json:"sex"`
}

//! 得到总下注
//func (self *Game_GoldBZW_Person) GetBets() int {
//	total := 0
//	for i := 0; i < len(self.Bets); i++ {
//		total += self.Bets[i]
//	}

//	return total
//}

//! 同步金币
func (self *Game_GoldBZW_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

type Game_GoldBZWSeat struct {
	Person *Game_GoldBZW_Person
	Robot  *lib.Robot
}

func (self *Game_GoldBZWSeat) GetTotal() int {
	if self.Person != nil {
		return self.Person.Total
	} else if self.Robot != nil {
		return self.Robot.GetMoney()
	}
	return 0
}

type Game_GoldBZW struct {
	PersonMgr   map[int64]*Game_GoldBZW_Person   `json:"personmgr"`
	Bets        [17]map[*Game_GoldBZW_Person]int `json:"bets"`
	Result      [3]int                           `json:"result"`
	Dealer      *Game_GoldBZW_Person             `json:"dealer"`      //! 庄家
	RobotDealer *lib.Robot                       `json:"robotdealer"` //! 机器人庄
	Round       int                              `json:"round"`       //! 连庄轮数
	DownUid     int64                            `json:"downuid"`     //! 下庄的人
	Time        int64                            `json:"time"`
	LstDeal     []Game_GoldBZWSeat               `json:"lstdeal"` //! 上庄列表
	Seat        [8]Game_GoldBZWSeat              `json:"seat"`    //! 8个位置
	Total       int                              `json:"total"`   //! 这局一共下了多少钱
	Money       int                              `json:"money"`   //! 系统庄的钱
	Trend       [][3]int                         `json:"trend"`   //! 走势
	Next        []int                            `json:"next"`
	BetTime     int                              `json:"bettime"`
	Robot       lib.ManyGameRobot                //! 机器人结构

	room *Room
}

func NewGame_GoldBZW() *Game_GoldBZW {
	game := new(Game_GoldBZW)
	game.PersonMgr = make(map[int64]*Game_GoldBZW_Person)
	for i := 0; i < len(game.Bets); i++ {
		game.Bets[i] = make(map[*Game_GoldBZW_Person]int)
	}
	for i := 0; i < 20; i++ {
		game.Trend = append(game.Trend, [3]int{lib.HF_GetRandom(6) + 1, lib.HF_GetRandom(6) + 1, lib.HF_GetRandom(6) + 1})
	}

	return game
}

func (self *Game_GoldBZW) OnInit(room *Room) {
	self.room = room
	self.Money = lib.GetManyMgr().GetProperty(self.room.Type).DealInitMoney
	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 11
	//! 载入机器人
	self.Robot.Init(17, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)
}

func (self *Game_GoldBZW) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldBZW) OnSendInfo(person *Person) {
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
		person.SendMsg("gamegoldbzwinfo", self.getInfo(person.Uid, value.Total))
		return
	}

	_person := new(Game_GoldBZW_Person)
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
	person.SendMsg("gamegoldbzwinfo", self.getInfo(person.Uid, person.Gold))
}

func (self *Game_GoldBZW) OnMsg(msg *RoomMsg) {
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
	case "gamechange":
		self.Change(msg.Uid, msg.V.(*Msg_GameChange).Card)
	}
}

func (self *Game_GoldBZW) Change(uid int64, card []int) {

	if !GetServer().IsAdmin(uid, staticfunc.ADMIN_GOLDBZW) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "非超端玩家")
		return
	}

	if len(card) == 3 {
		for i := 0; i < len(card); i++ {
			if card[i] < 1 && card[i] > 6 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "设置点数不符合规则 card : ", card)
				return
			}
		}
		lib.HF_DeepCopy(&self.Next, &card)
	} else if len(card) == 1 {
		//! 3-18 点数  21-大 22-小 23-豹子  31-庄赢 32-闲赢
		if !((card[0] >= 3 && card[0] <= 18) || (card[0] >= 21 && card[0] <= 23) || (card[0] == 31 || card[0] == 32)) {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "设置点数不符合规则 card : ", card)
			return
		}
		lib.HF_DeepCopy(&self.Next, &card)
	} else {
		return
	}

	GetServer().SqlSuperClientLog(&SQL_SuperClientLog{1, uid, self.room.Type, lib.HF_JtoA(&card), time.Now().Unix()})
	self.room.SendMsg(uid, "ok", nil)
}

func (self *Game_GoldBZW) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.Begin = true

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "next:", self.Next)
	if len(self.Next) == 3 {
		self.Result[0] = self.Next[0]
		self.Result[1] = self.Next[1]
		self.Result[2] = self.Next[2]
		self.Next = make([]int, 0)
		self.OnEnd()
		return
	} else if len(self.Next) == 1 {
		if self.Next[0] >= 3 && self.Next[0] <= 18 { //! 开点数
			self.GetResult(self.Next[0], false)
			self.Next = make([]int, 0)
			self.OnEnd()
			return
		}

		if self.Next[0] >= 21 && self.Next[0] <= 23 { //!21-大 22-小 23-豹子
			lst := make([]GameGoldBZW_CanResult, 0)
			if self.Next[0] == 21 {
				for i := 11; i <= 17; i++ { //! 模拟11到17点
					win := self.GetDealWin(i, false)
					if GetServer().SystemMoney[self.room.Type%40000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
						lst = append(lst, GameGoldBZW_CanResult{i, false})
					}
				}
			} else if self.Next[0] == 22 {
				for i := 4; i <= 10; i++ { //! 模拟4到10点
					win := self.GetDealWin(i, false)
					if GetServer().SystemMoney[self.room.Type%40000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
						lst = append(lst, GameGoldBZW_CanResult{i, false})
					}
				}
			} else {
				for i := 3; i <= 18; i++ {
					if i%3 == 0 { //! 可能是豹子
						win := self.GetDealWin(i, true)
						if GetServer().SystemMoney[self.room.Type%40000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
							lst = append(lst, GameGoldBZW_CanResult{i, true})
						}
					}
				}
			}
			if len(lst) == 0 {
				if self.Next[0] == 21 {
					self.GetResult(lib.HF_GetRandom(7)+11, false)
				} else if self.Next[0] == 22 {
					self.GetResult(lib.HF_GetRandom(7)+4, false)
				} else {
					self.GetResult((lib.HF_GetRandom(6)+1)*3, true)
				}
			} else {
				result := lst[lib.HF_GetRandom(len(lst))]
				self.GetResult(result.DS, result.BZ)
			}
			self.Next = make([]int, 0)
			self.OnEnd()
			return
		}

		if self.Next[0] == 31 || self.Next[0] == 32 { //! 31-庄赢 32-闲赢
			lst := make([]GameGoldBZW_CanResult, 0)
			if self.Next[0] == 31 {
				for i := 3; i <= 18; i++ { //! 模拟3到18点
					if i != 3 && i != 18 { //! 3和18必然是豹子
						win := self.GetDealWin(i, false)
						if GetServer().SystemMoney[self.room.Type%40000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && win > 0 {
							lst = append(lst, GameGoldBZW_CanResult{i, false})
						}
					}
				}
				if len(lst) == 0 || lib.HF_GetRandom(100) < 15 {
					for i := 3; i <= 18; i++ {
						if i%3 == 0 { //! 可能是豹子
							win := self.GetDealWin(i, true)
							if GetServer().SystemMoney[self.room.Type%40000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && win > 0 {
								lst = append(lst, GameGoldBZW_CanResult{i, true})
							}
						}
					}
				}
			} else {
				for i := 3; i <= 18; i++ { //! 模拟3到18点
					if i != 3 && i != 18 { //! 3和18必然是豹子
						win := self.GetDealWin(i, false)
						if GetServer().SystemMoney[self.room.Type%40000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && win < 0 {
							lst = append(lst, GameGoldBZW_CanResult{i, false})
						}
					}
				}
				if len(lst) == 0 || lib.HF_GetRandom(100) < 15 {
					for i := 3; i <= 18; i++ {
						if i%3 == 0 { //! 可能是豹子
							win := self.GetDealWin(i, true)
							if GetServer().SystemMoney[self.room.Type%40000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && win < 0 {
								lst = append(lst, GameGoldBZW_CanResult{i, true})
							}
						}
					}
				}
			}

			if len(lst) == 0 {
				lib.GetLogMgr().Output(lib.LOG_ERROR, "self.next 没有合适方案 ")
				for i := 0; i < len(self.Result); i++ {
					self.Result[i] = lib.HF_GetRandom(6) + 1
				}
			} else {
				result := lst[lib.HF_GetRandom(len(lst))]
				self.GetResult(result.DS, result.BZ)
			}

			self.Next = make([]int, 0)
			self.OnEnd()
			return
		}
	}

	if self.Dealer != nil {
		if self.Robot.RobotTotal == 0 { //! 没有机器人下注
			if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost == 102 { //! 系统算法
				lst := make([]GameGoldBZW_CanResult, 0)
				winlst := make([]GameGoldBZW_CanResult, 0)
				lostlst := make([]GameGoldBZW_CanResult, 0)
				for i := 3; i <= 18; i++ { //! 模拟3到18点
					if i != 3 && i != 18 { //! 3和18必然是豹子
						win := self.GetDealWin(i, false)
						if GetServer().PlayerMoney[self.room.Type%40000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMin && GetServer().PlayerMoney[self.room.Type%40000]+int64(win) <= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMax {
							lst = append(lst, GameGoldBZW_CanResult{i, false})
						}
						if win > 0 {
							winlst = append(winlst, GameGoldBZW_CanResult{i, false})
						} else {
							lostlst = append(lostlst, GameGoldBZW_CanResult{i, false})
						}
					}
				}
				if len(lst) == 0 || lib.HF_GetRandom(100) < 10 {
					for i := 3; i <= 18; i++ {
						if i%3 == 0 { //! 可能是豹子
							win := self.GetDealWin(i, true)
							if GetServer().PlayerMoney[self.room.Type%40000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMin && GetServer().PlayerMoney[self.room.Type%40000]+int64(win) <= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMax {
								lst = append(lst, GameGoldBZW_CanResult{i, true})
							}
							if win > 0 {
								winlst = append(winlst, GameGoldBZW_CanResult{i, true})
							} else {
								lostlst = append(lostlst, GameGoldBZW_CanResult{i, true})
							}
						}
					}
				}
				if len(lst) == 0 { //! 如果都输钱，就随机，理论上不会发生这种情况
					if GetServer().PlayerMoney[self.room.Type%40000] >= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMax && len(lostlst) > 0 { //! 一定输
						result := lostlst[lib.HF_GetRandom(len(lostlst))]
						self.GetResult(result.DS, result.BZ)
					} else if GetServer().PlayerMoney[self.room.Type%40000] <= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMin && len(winlst) > 0 {
						result := winlst[lib.HF_GetRandom(len(winlst))]
						self.GetResult(result.DS, result.BZ)
					} else {
						lib.GetLogMgr().Output(lib.LOG_ERROR, "豹子王亏了1")
						for i := 0; i < len(self.Result); i++ {
							self.Result[i] = lib.HF_GetRandom(6) + 1
						}
					}
				} else {
					result := lst[lib.HF_GetRandom(len(lst))]
					self.GetResult(result.DS, result.BZ)
				}
			} else if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost > 100 { //! 纯随机
				for i := 0; i < len(self.Result); i++ {
					self.Result[i] = lib.HF_GetRandom(6) + 1
				}
			} else { //! 设置概率
				iswin := lib.HF_GetRandom(100) < lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost
				lst := make([]GameGoldBZW_CanResult, 0)
				smalllst := make([]GameGoldBZW_CanResult, 0)
				biglst := make([]GameGoldBZW_CanResult, 0)
				for i := 3; i <= 18; i++ {
					if i%3 == 0 { //! 可能是豹子
						win := self.GetDealWin(i, true)
						if iswin && win >= 0 {
							smalllst = append(smalllst, GameGoldBZW_CanResult{i, true})
						} else if !iswin && win < 0 {
							biglst = append(biglst, GameGoldBZW_CanResult{i, true})
						}
					}
					if i != 3 && i != 18 { //! 3和18必然是豹子
						win := self.GetDealWin(i, false)
						if iswin && win >= 0 {
							lst = append(lst, GameGoldBZW_CanResult{i, false})
						} else if !iswin && win < 0 {
							if i == 4 || i == 17 {
								biglst = append(biglst, GameGoldBZW_CanResult{i, false})
							} else if i == 5 || i == 16 {
								smalllst = append(smalllst, GameGoldBZW_CanResult{i, false})
							} else {
								lst = append(lst, GameGoldBZW_CanResult{i, false})
							}
						}
					}
				}
				if len(lst) == 0 || lib.HF_GetRandom(100) < 10 {
					lst = append(lst, smalllst...)
				}
				if len(lst) == 0 || lib.HF_GetRandom(100) < 5 {
					lst = append(lst, biglst...)
				}
				if len(lst) == 0 { //! 如果都输钱，就随机，理论上不会发生这种情况
					lib.GetLogMgr().Output(lib.LOG_ERROR, "豹子王庄家算错了")
					for i := 0; i < len(self.Result); i++ {
						self.Result[i] = lib.HF_GetRandom(6) + 1
					}
				} else {
					result := lst[lib.HF_GetRandom(len(lst))]
					self.GetResult(result.DS, result.BZ)
				}
			}
		} else {
			lst := make([]GameGoldBZW_CanResult, 0)
			for i := 3; i <= 18; i++ { //! 模拟3到18点
				if i != 3 && i != 18 { //! 3和18必然是豹子
					win := self.GetRobotWin(i, false)
					if lib.GetRobotMgr().GetRobotWin(self.room.Type)+win >= 0 {
						lst = append(lst, GameGoldBZW_CanResult{i, false})
					}
				}
			}
			if len(lst) == 0 || lib.HF_GetRandom(100) < 10 {
				for i := 3; i <= 18; i++ {
					if i%3 == 0 { //! 可能是豹子
						win := self.GetRobotWin(i, true)
						if lib.GetRobotMgr().GetRobotWin(self.room.Type)+win >= 0 {
							lst = append(lst, GameGoldBZW_CanResult{i, true})
						}
					}
				}
			}
			if len(lst) == 0 { //! 如果都输钱，就随机，理论上不会发生这种情况
				for i := 0; i < len(self.Result); i++ {
					self.Result[i] = lib.HF_GetRandom(6) + 1
				}
			} else {
				result := lst[lib.HF_GetRandom(len(lst))]
				self.GetResult(result.DS, result.BZ)
			}
		}
	} else {
		//! 先随机试试
		//{
		//	ds := 0
		//	for i := 0; i < len(self.Result); i++ {
		//		self.Result[i] = lib.HF_GetRandom(6) + 1
		//		ds += self.Result[i]
		//	}
		//	win := self.GetDealWin(ds, self.Result[0] == self.Result[1] && self.Result[1] == self.Result[2])
		//	if GetServer().SystemMoney[self.room.Type%40000]+int64(win) >= GetServer().Con.BZWMinMoney[self.room.Type%40000] {
		//		self.OnEnd()
		//		return
		//	}
		//}

		lst := make([]GameGoldBZW_CanResult, 0)
		for i := 3; i <= 18; i++ { //! 模拟3到18点
			if i != 3 && i != 18 { //! 3和18必然是豹子
				win := self.GetDealWin(i, false)
				if GetServer().SystemMoney[self.room.Type%40000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
					lst = append(lst, GameGoldBZW_CanResult{i, false})
				} else if win >= 0 {
					lst = append(lst, GameGoldBZW_CanResult{i, false})
				}
			}
		}
		if len(lst) == 0 || lib.HF_GetRandom(100) < 15 {
			for i := 3; i <= 18; i++ {
				if i%3 == 0 { //! 可能是豹子
					win := self.GetDealWin(i, true)
					if GetServer().SystemMoney[self.room.Type%40000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
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
func (self *Game_GoldBZW) GetDealWinByPos(pos int, robot bool) int {
	lost := self.GetMoneyPos(pos, robot) * GOLDBZW_BS[pos] //! 压中这个点本身赢的钱
	maxlost := 0
	if pos == 0 {
		for i := 3; i <= 9; i++ {
			tmp := self.GetMoneyPos(i, robot) * GOLDBZW_BS[i]
			if tmp > maxlost {
				maxlost = tmp
			}
		}
	} else if pos == 1 {
		for i := 3; i <= 16; i++ {
			if (i+1)%3 == 0 {
				tmp := self.GetMoneyPos(i, robot) * GOLDBZW_BS[i]
				if tmp > maxlost {
					maxlost = tmp
				}
			}
		}
	} else if pos == 2 {
		for i := 10; i <= 16; i++ {
			tmp := self.GetMoneyPos(i, robot) * GOLDBZW_BS[i]
			if tmp > maxlost {
				maxlost = tmp
			}
		}
	} else {
		if (pos+1)%3 == 0 {
			maxlost = self.GetMoneyPos(1, robot) * GOLDBZW_BS[1]
		}
		if pos+1 >= 4 && pos+1 <= 10 {
			tmp := self.GetMoneyPos(0, robot) * GOLDBZW_BS[0]
			if tmp > maxlost {
				maxlost = tmp
			}
		} else {
			tmp := self.GetMoneyPos(2, robot) * GOLDBZW_BS[2]
			if tmp > maxlost {
				maxlost = tmp
			}
		}
	}

	return lost + maxlost
}

//! 得到庄家可以赢的钱
func (self *Game_GoldBZW) GetDealWin(ds int, bz bool) int {
	lost := 0
	if bz {
		lost += self.GetMoneyPos(1, false) * GOLDBZW_BS[1]
	} else {
		if ds >= 4 && ds <= 10 {
			lost += self.GetMoneyPos(0, false) * GOLDBZW_BS[0]
		}
		if ds >= 11 && ds <= 17 {
			lost += self.GetMoneyPos(2, false) * GOLDBZW_BS[2]
		}
	}
	for i := 3; i < 17; i++ {
		if ds == i+1 {
			lost += self.GetMoneyPos(i, false) * GOLDBZW_BS[i]
			break
		}
	}

	if self.Dealer == nil {
		return self.Total - lost
	} else {
		return self.Total + self.Robot.RobotTotal - lost
	}
}

//! 得到机器人可以赢的钱
func (self *Game_GoldBZW) GetRobotWin(ds int, bz bool) int {
	win := 0
	if bz {
		win += self.GetMoneyPosByRobot(1) * GOLDBZW_BS[1]
	} else {
		if ds >= 4 && ds <= 10 {
			win += self.GetMoneyPosByRobot(0) * GOLDBZW_BS[0]
		}
		if ds >= 11 && ds <= 17 {
			win += self.GetMoneyPosByRobot(2) * GOLDBZW_BS[2]
		}
	}
	for i := 3; i < 17; i++ {
		if ds == i+1 {
			win += self.GetMoneyPosByRobot(i) * GOLDBZW_BS[i]
			break
		}
	}

	return win - self.Robot.RobotTotal
}

//! 根据点数得到实际
func (self *Game_GoldBZW) GetResult(ds int, bz bool) {
	if bz {
		self.Result[0] = ds / 3
		self.Result[1] = ds / 3
		self.Result[2] = ds / 3
		return
	}

	//! 得到第一个数
	max := lib.HF_MinInt(6, ds-2)
	min := lib.HF_MaxInt(1, ds-12)
	self.Result[0] = lib.HF_GetRandom(max-min) + min
	//! 得到第二个数
	ds -= self.Result[0]
	max = lib.HF_MinInt(6, ds-1)
	min = lib.HF_MaxInt(1, ds-6)
	self.Result[1] = lib.HF_GetRandom(max-min) + min
	//! 得到第三个数
	self.Result[2] = ds - self.Result[1]

	if self.Result[0] == self.Result[1] && self.Result[1] == self.Result[2] { //! 不能是豹子
		a := self.Result[0] - 1
		b := 6 - self.Result[2]
		c := lib.HF_MinInt(a, b)
		self.Result[0] -= c
		self.Result[2] += c
	}
}

func (self *Game_GoldBZW) IsType() int { //! 0小 1豹子  2大
	if self.Result[0] == self.Result[1] && self.Result[0] == self.Result[2] {
		return 1
	}

	if self.Result[0]+self.Result[1]+self.Result[2] <= 10 {
		return 0
	}

	return 2
}

func (self *Game_GoldBZW) GameUpDeal(uid int64) {
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
			if self.LstDeal[i].Person == person {
				self.room.SendErr(uid, "您已经在上庄列表中，请等待上庄")
				return
			}
		}
		if len(self.LstDeal) == 0 {
			self.Round = 0
		}
		self.LstDeal = append(self.LstDeal, Game_GoldBZWSeat{person, nil})
	}
	person.Round = 0

	var msg Msg_GameGoldBZW_DealList
	msg.Type = 0
	msg.Info = make([]Son_GameGoldBZW_Info, 0)
	for i := 0; i < len(self.LstDeal); i++ {
		if self.LstDeal[i].Person != nil {
			msg.Info = append(msg.Info, Son_GameGoldBZW_Info{self.LstDeal[i].Person.Uid, self.LstDeal[i].Person.Name, self.LstDeal[i].Person.Head, self.LstDeal[i].Person.Total, self.LstDeal[i].Person.IP, self.LstDeal[i].Person.Address, self.LstDeal[i].Person.Sex})
		} else if self.LstDeal[i].Robot != nil {
			msg.Info = append(msg.Info, Son_GameGoldBZW_Info{self.LstDeal[i].Robot.Id, self.LstDeal[i].Robot.Name, self.LstDeal[i].Robot.Head, self.LstDeal[i].Robot.GetMoney(), self.LstDeal[i].Robot.IP, self.LstDeal[i].Robot.Address, self.LstDeal[i].Robot.Sex})
		}
	}
	self.room.SendMsg(uid, "gamebzwdeal", &msg)
}

//! 机器人上庄
func (self *Game_GoldBZW) RobotUpDeal(robot *lib.Robot) {
	if robot.GetMoney() < lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney {
		return
	}

	if robot.GetSeat() == 100 {
		return
	}

	for i := 0; i < len(self.LstDeal); i++ {
		if self.LstDeal[i].Robot == robot {
			return
		}
	}

	if len(self.LstDeal) == 0 {
		self.Round = 0
	}
	self.LstDeal = append(self.LstDeal, Game_GoldBZWSeat{nil, robot})
}

func (self *Game_GoldBZW) GameReDeal(uid int64) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if self.Dealer == person { //! 正在庄
		self.DownUid = uid
		self.room.SendErr(uid, "您已成功下庄，请等待本局结束")
	} else {
		for i := 0; i < len(self.LstDeal); i++ {
			if self.LstDeal[i].Person == person {
				copy(self.LstDeal[i:], self.LstDeal[i+1:])
				self.LstDeal = self.LstDeal[:len(self.LstDeal)-1]
				break
			}
		}
	}

	var msg Msg_GameGoldBZW_DealList
	msg.Type = 1
	msg.Info = make([]Son_GameGoldBZW_Info, 0)
	self.room.SendMsg(uid, "gamebzwdeal", &msg)
}

//! 坐下
func (self *Game_GoldBZW) GameSeat(uid int64, index int) {
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
		if self.Seat[i].Person == person {
			return
		}
	}

	if self.Seat[index].Person != nil {
		if person.Total <= self.Seat[index].Person.Total {
			self.room.SendErr(uid, "该位置已经有人坐了")
			return
		}
		//! 把原来在这个位置上的人挤下去
		self.Seat[index].Person.Seat = -1
	} else if self.Seat[index].Robot != nil {
		if person.Total <= self.Seat[index].Robot.GetMoney() {
			self.room.SendErr(uid, "该位置已经有人坐了")
			return
		}
		//! 把原来在这个位置上的人挤下去
		self.Seat[index].Robot.SetSeat(-1)
	}

	self.Seat[index].Person = person
	self.Seat[index].Robot = nil
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
	self.room.broadCastMsg("gamebzwseat", &msg)
}

//! 机器人坐下
func (self *Game_GoldBZW) RobotSeat(index int, robot *lib.Robot) {
	if index < 0 || index > 7 {
		return
	}

	if robot.GetMoney() < lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney {
		return
	}

	for i := 0; i < len(self.Seat); i++ {
		if self.Seat[i].Robot == robot {
			return
		}
	}

	if self.Seat[index].Person != nil || self.Seat[index].Robot != nil {
		return
	}

	self.Seat[index].Person = nil
	self.Seat[index].Robot = robot
	robot.SetSeat(index)

	var msg Msg_GameGoldBZW_UpdSeat
	msg.Uid = robot.Id
	msg.Index = index
	msg.Head = robot.Head
	msg.Name = robot.Name
	msg.Total = robot.GetMoney()
	msg.IP = robot.IP
	msg.Address = robot.Address
	msg.Sex = robot.Sex
	self.room.broadCastMsg("gamebzwseat", &msg)
}

func (self *Game_GoldBZW) GameBets(uid int64, index int, gold int) {
	if uid == 0 {
		return
	}

	if index < 0 || index > 16 {
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
			self.room.SendErr(uid, fmt.Sprintf("%d金币以上才能下注。", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("%d金币以上才能下注。", lib.GetManyMgr().GetProperty(self.room.Type).MinBet))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("%d万金币以上才能下注。", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/10000))
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
		if self.Dealer != nil {
			dealmoney = self.Dealer.Total
		} else if self.RobotDealer != nil {
			dealmoney = self.RobotDealer.GetMoney()
		} else if lib.GetManyMgr().GetProperty(self.room.Type).DealChange == 1 {
			dealmoney = self.Money
		}
		//! 模拟庄家赢钱
		if dealmoney > 0 {
			dealwin := self.GetDealWinByPos(index, true)
			if dealmoney+self.Total+self.Robot.RobotTotal-dealwin < dealmoney/5 {
				self.Total -= gold
				self.Bets[index][person] -= gold
				self.room.SendErr(uid, "庄家金币不足，该位置无法下注。")
				return
			}
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
	self.room.broadCastMsg("gamebzwbets", &msg)
}

//! 续压
func (self *Game_GoldBZW) GameGoOn(uid int64) {
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
			self.room.SendErr(uid, fmt.Sprintf("%d金币以上才能下注。", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/100))
		} else if GetServer().Con.MoneyMode == 0 {
			self.room.SendErr(uid, fmt.Sprintf("%d金币以上才能下注。", lib.GetManyMgr().GetProperty(self.room.Type).MinBet))
		} else {
			self.room.SendErr(uid, fmt.Sprintf("%d万金币以上才能下注。", lib.GetManyMgr().GetProperty(self.room.Type).MinBet/10000))
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
		if self.Dealer != nil {
			dealmoney = self.Dealer.Total
		} else if self.RobotDealer != nil {
			dealmoney = self.RobotDealer.GetMoney()
		} else if lib.GetManyMgr().GetProperty(self.room.Type).DealChange == 1 {
			dealmoney = self.Money
		}
		//! 模拟庄家赢钱
		if dealmoney > 0 {
			dealwin := 0
			for i := 0; i < len(person.BeBetInfo); i++ {
				if person.BeBetInfo[i] == 0 {
					continue
				}
				tmp := self.GetDealWinByPos(i, true)
				if tmp > dealwin {
					dealwin = tmp
				}
			}
			if dealmoney+self.Total+self.Robot.RobotTotal-dealwin < dealmoney/5 {
				for i := 0; i < len(person.BeBetInfo); i++ {
					self.Total -= person.BeBetInfo[i]
					self.Bets[i][person] -= person.BeBetInfo[i]
				}
				self.room.SendErr(uid, "庄家金币不足，无法续压。")
				return
			}
		}
	}

	person.Bets += person.BeBets
	person.Total -= person.BeBets
	for i := 0; i < len(person.BeBetInfo); i++ {
		person.BetInfo[i] += person.BeBetInfo[i]
	}
	person.Round = 0

	var msg Msg_GameGoldBZW_Goon
	msg.Uid = uid
	msg.Gold = person.BeBetInfo
	msg.Total = person.Total
	self.room.broadCastMsg("gamebzwgoon", &msg)
}

//! 结算
func (self *Game_GoldBZW) OnEnd() {
	self.room.Begin = false

	ds := 0
	for i := 0; i < len(self.Result); i++ {
		ds += self.Result[i]
	}

	trend := self.IsType()
	tmp := [][3]int{self.Result}
	tmp = append(tmp, self.Trend...)
	if len(tmp) > 20 {
		tmp = tmp[0:20]
	}
	self.Trend = tmp

	dealwin := 0
	robotwin := 0
	for i := 0; i < len(self.Bets); i++ {
		if i < 3 { //! 上面一排
			if trend == i { //!  下注赢了
				for key, value := range self.Bets[i] {
					winmoney := value * GOLDBZW_BS[i]
					dealwin -= winmoney
					key.Win += winmoney
					key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
				}

				for key, value := range self.Robot.RobotsBet[i] {
					winmoney := value * GOLDBZW_BS[i]
					key.AddWin(winmoney)
					key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
					robotwin += winmoney
					if self.Dealer != nil {
						dealwin -= winmoney
					}
				}
			}
		} else {
			if self.Result[0]+self.Result[1]+self.Result[2] == GOLDBZW_DS[i] { //! 下注赢了
				for key, value := range self.Bets[i] {
					winmoney := value * GOLDBZW_BS[i]
					dealwin -= winmoney
					key.Win += winmoney
					key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
				}

				for key, value := range self.Robot.RobotsBet[i] {
					winmoney := value * GOLDBZW_BS[i]
					key.AddWin(winmoney)
					key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
					robotwin += winmoney
					if self.Dealer != nil {
						dealwin -= winmoney
					}
				}
			}
		}
	}
	robotwin -= self.Robot.RobotTotal
	if self.Dealer != nil {
		dealwin += self.Robot.RobotTotal
	}
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "RobotTotal=", self.Robot.RobotTotal)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "robotwin=", robotwin)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "dealwin=", dealwin)

	var bigwin *GameGold_BigWin = nil
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
				if self.Seat[j].Person == value {
					self.room.broadCastMsg("gamegoldbzwbalance", &msg)
					find = true
					break
				}
			}
			if !find {
				self.room.SendMsg(value.Uid, "gamegoldbzwbalance", &msg)
			}

			if bigwin == nil {
				bigwin = &GameGold_BigWin{value.Uid, value.Name, value.Head, value.Win}
			} else if value.Win > bigwin.Win {
				bigwin = &GameGold_BigWin{value.Uid, value.Name, value.Head, value.Win}
			}
		} else if value.Win-value.Bets < 0 {
			cost := int(math.Ceil(float64(value.Win-value.Bets) * float64(lib.GetManyMgr().GetProperty(self.room.Type).Cost) / 200.0))
			GetServer().SqlAgentBillsLog(value.Uid, cost, self.room.Type)
		}

		//! 插入战绩
		if value.Bets > 0 {
			var record Rec_BZW_Info
			record.Time = time.Now().Unix()
			record.GameType = self.room.Type
			var rec Son_Rec_BZW_Person
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

	for i := 0; i < len(self.Robot.Robots); i++ {
		if self.Robot.Robots[i].GetWin() > 0 {
			self.Robot.Robots[i].AddWin(-self.Robot.Robots[i].GetCost())
			self.Robot.Robots[i].AddMoney(self.Robot.Robots[i].GetWin())

			for j := 0; j < len(self.Seat); j++ {
				if self.Seat[j].Robot == self.Robot.Robots[i] {
					var msg Msg_GameGoldBZW_Balance
					msg.Uid = self.Robot.Robots[i].Id
					msg.Total = self.Robot.Robots[i].GetMoney()
					msg.Win = self.Robot.Robots[i].GetWin()
					self.room.broadCastMsg("gamegoldbzwbalance", &msg)
					break
				}
			}

			if bigwin == nil {
				bigwin = &GameGold_BigWin{self.Robot.Robots[i].Id, self.Robot.Robots[i].Name, self.Robot.Robots[i].Head, self.Robot.Robots[i].GetWin()}
			} else if self.Robot.Robots[i].GetWin() > bigwin.Win {
				bigwin = &GameGold_BigWin{self.Robot.Robots[i].Id, self.Robot.Robots[i].Name, self.Robot.Robots[i].Head, self.Robot.Robots[i].GetWin()}
			}
		}
	}

	dealwin = self.Total + dealwin
	if self.Dealer == nil && dealwin != 0 { //! 系统庄
		GetServer().SqlBZWLog(&SQL_BZWLog{1, dealwin, time.Now().Unix(), self.room.Type})
	}
	if self.Dealer != nil && robotwin != 0 { //! 玩家庄
		lib.GetRobotMgr().AddRobotWin(self.room.Type, robotwin)
		GetServer().SqlBZWLog(&SQL_BZWLog{1, robotwin, time.Now().Unix(), self.room.Type + 10000000})
	}
	if self.Dealer != nil && lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost == 102 { //! 玩家庄并且是系统算法
		GetServer().SetPlayerMoney(self.room.Type%40000, GetServer().PlayerMoney[self.room.Type%40000]+int64(dealwin))
	}
	if self.Dealer != nil { //! 玩家庄
		if dealwin > 0 {
			bl := lib.GetManyMgr().GetProperty(self.room.Type).Cost
			cost := int(math.Ceil(float64(dealwin) * bl / 100.0))
			dealwin -= cost
			GetServer().SqlAgentGoldLog(self.Dealer.Uid, cost, self.room.Type)
			GetServer().SqlAgentBillsLog(self.Dealer.Uid, cost/2, self.room.Type)
		} else if dealwin < 0 {
			cost := int(math.Ceil(float64(dealwin) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 200.0))
			GetServer().SqlAgentBillsLog(self.Dealer.Uid, cost, self.room.Type)
		}
	} else {
		_dealwin := dealwin
		if _dealwin > 0 {
			bl := lib.GetManyMgr().GetProperty(self.room.Type).DealCost
			cost := int(math.Ceil(float64(_dealwin) * bl / 100.0))
			_dealwin -= cost
		}
		GetServer().SetSystemMoney(self.room.Type%40000, GetServer().SystemMoney[self.room.Type%40000]+int64(_dealwin))
		dealwin -= robotwin
		if dealwin > 0 {
			bl := lib.GetManyMgr().GetProperty(self.room.Type).Cost
			cost := int(math.Ceil(float64(dealwin) * bl / 100.0))
			dealwin -= cost
		}
	}

	if self.Dealer != nil {
		var record Rec_BZW_Info
		record.Time = time.Now().Unix()
		record.GameType = self.room.Type
		var rec Son_Rec_BZW_Person
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
		var msg Msg_GameGoldBZW_Balance
		if self.Dealer != nil { //! 玩家庄
			if self.Dealer.Total+dealwin > 0 {
				self.Dealer.Total += dealwin
			} else {
				self.Dealer.Total = 0
				dealwin = -self.Dealer.Total
			}
			msg.Uid = self.Dealer.Uid
			msg.Total = self.Dealer.Total
		} else if self.RobotDealer != nil { //! 机器人庄
			self.RobotDealer.AddMoney(dealwin)
			msg.Uid = self.RobotDealer.Id
			msg.Total = self.RobotDealer.GetMoney()
		} else {
			if lib.GetManyMgr().GetProperty(self.room.Type).DealChange == 1 {
				self.Money += dealwin
				if self.Money <= 0 {
					self.Money = lib.GetManyMgr().GetProperty(self.room.Type).DealInitMoney
				}
			} else {
				self.Money = lib.GetManyMgr().GetProperty(self.room.Type).DealInitMoney
			}
			msg.Total = self.Money
		}
		msg.Win = dealwin
		self.room.broadCastMsg("gamegoldbzwbalance", &msg)
	}

	//! 30秒的下注时间
	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 11
	self.SetTime(self.BetTime)

	//! 总结算
	{
		var msg Msg_GameGoldBZW_End
		msg.Result = self.Result
		if bigwin != nil {
			msg.Uid = bigwin.Uid
			msg.Name = bigwin.Name
			msg.Head = bigwin.Head
		}
		msg.Money = lib.GetManyMoneyMgr().GetProperty(self.room.Type).Money
		msg.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime
		self.room.broadCastMsg("gamegoldbzwend", &msg)
	}

	self.Total = 0

	//! 把不在room.uid里面的玩家清理出去
	for key, value := range self.PersonMgr {
		if value.Online {
			find := false
			for i := 0; i < len(self.LstDeal); i++ {
				if self.LstDeal[i].Person == value {
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
			if self.LstDeal[j].Person == value {
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
			if self.Seat[j].Person == value {
				self.Seat[j].Person = nil
				var msg Msg_GameGoldBZW_UpdSeat
				msg.Index = j
				self.room.broadCastMsg("gamebzwseat", &msg)
				break
			}
		}
		delete(self.PersonMgr, key)
	}

	//! 返回机器人
	self.Robot.Init(17, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)

	for i := 0; i < len(self.Bets); i++ {
		self.Bets[i] = make(map[*Game_GoldBZW_Person]int)
	}

	//! 判断庄家是否能继续连
	if self.Dealer != nil {
		if self.Dealer.Total < lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney || self.DownUid == self.Dealer.Uid || GetPersonMgr().GetPerson(self.Dealer.Uid) == nil {
			self.ChageDeal()
		} else {
			if self.Round >= 10 && len(self.LstDeal) > 0 {
				self.ChageDeal()
			} else {
				self.Round++
			}
		}
	} else if self.RobotDealer != nil {
		if self.RobotDealer.GetMoney() < lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney {
			self.ChageDeal()
		} else {
			if self.Round >= lib.HF_GetRandom(6)+3 && len(self.LstDeal) > 0 || !lib.GetRobotMgr().GetRobotSet(self.room.Type).NeedRobot {
				self.ChageDeal()
			} else {
				self.Round++
			}
		}
	} else if len(self.LstDeal) > 0 {
		self.ChageDeal()
	}

	//! 判断坐下的人是否能继续坐
	for i := 0; i < len(self.Seat); i++ {
		if self.Seat[i].Person == nil {
			continue
		}
		if self.Seat[i].Person.Total < lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney {
			self.Seat[i].Person.Seat = -1
			var msg Msg_GameGoldBZW_UpdSeat
			msg.Index = i
			self.room.broadCastMsg("gamebzwseat", &msg)
			self.Seat[i].Person = nil
		}
	}
	for i := 0; i < len(self.Seat); i++ {
		if self.Seat[i].Robot == nil {
			continue
		}
		find := false
		for j := 0; j < len(self.Robot.Robots); j++ {
			if self.Robot.Robots[j] == self.Seat[i].Robot {
				find = true
				break
			}
		}
		if !find || self.Seat[i].Robot.GetSeat() != i || self.Seat[i].Robot.GetMoney() < lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney {
			self.Seat[i].Robot.SetSeat(-1)
			var msg Msg_GameGoldBZW_UpdSeat
			msg.Index = i
			self.room.broadCastMsg("gamebzwseat", &msg)
			self.Seat[i].Robot = nil
		}
	}
}

func (self *Game_GoldBZW) OnBye() {
}

func (self *Game_GoldBZW) OnExit(uid int64) {
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
func (self *Game_GoldBZW) GetMoneyPos(index int, robot bool) int {
	total := 0
	for _, value := range self.Bets[index] {
		total += value
	}
	if robot || self.Dealer != nil { //! 是玩家庄,判断机器人下注
		for _, value := range self.Robot.RobotsBet[index] {
			total += value
		}
	}
	return total
}

//! 得到这个位置机器人下了多少钱
func (self *Game_GoldBZW) GetMoneyPosByRobot(index int) int {
	total := 0
	for _, value := range self.Robot.RobotsBet[index] {
		total += value
	}
	return total
}

func (self *Game_GoldBZW) getInfo(uid int64, total int) *Msg_GameGoldBZW_Info {
	var msg Msg_GameGoldBZW_Info
	msg.Begin = self.room.Begin
	msg.Time = self.Time - time.Now().Unix()
	msg.Total = total
	msg.Trend = self.Trend
	msg.IsDeal = false
	msg.Money = lib.GetManyMoneyMgr().GetProperty(self.room.Type).Money
	msg.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime
	if GetServer().IsAdmin(uid, staticfunc.ADMIN_GOLDBZW) {
		msg.Change = true
	} else {
		msg.Change = false
	}

	if self.Dealer != nil && self.Dealer.Uid == uid {
		msg.IsDeal = true
	} else {
		for i := 0; i < len(self.LstDeal); i++ {
			if self.LstDeal[i].Person != nil && self.LstDeal[i].Person.Uid == uid {
				msg.IsDeal = true
				break
			}
		}
	}
	for i := 0; i < len(self.Bets); i++ {
		msg.Bets[i] = self.GetMoneyPos(i, false)
	}
	if self.Dealer == nil { //! 系统庄的时候上面计算了机器人的下注
		for i := 0; i < len(self.Bets); i++ {
			msg.Bets[i] += self.GetMoneyPosByRobot(i)
		}
	}
	for i := 0; i < 8; i++ {
		if self.Seat[i].Person != nil {
			msg.Seat[i].Uid = self.Seat[i].Person.Uid
			msg.Seat[i].Name = self.Seat[i].Person.Name
			msg.Seat[i].Head = self.Seat[i].Person.Head
			msg.Seat[i].Total = self.Seat[i].Person.Total
			msg.Seat[i].Ip = self.Seat[i].Person.IP
			msg.Seat[i].Address = self.Seat[i].Person.Address
			msg.Seat[i].Sex = self.Seat[i].Person.Sex
		} else if self.Seat[i].Robot != nil {
			msg.Seat[i].Uid = self.Seat[i].Robot.Id
			msg.Seat[i].Name = self.Seat[i].Robot.Name
			msg.Seat[i].Head = self.Seat[i].Robot.Head
			msg.Seat[i].Total = self.Seat[i].Robot.GetMoney()
			msg.Seat[i].Ip = self.Seat[i].Robot.IP
			msg.Seat[i].Address = self.Seat[i].Robot.Address
			msg.Seat[i].Sex = self.Seat[i].Robot.Sex
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
	} else if self.RobotDealer != nil {
		msg.Dealer.Uid = self.RobotDealer.Id
		msg.Dealer.Name = self.RobotDealer.Name
		msg.Dealer.Head = self.RobotDealer.Head
		msg.Dealer.Total = self.RobotDealer.GetMoney()
		msg.Dealer.Ip = self.RobotDealer.IP
		msg.Dealer.Address = self.RobotDealer.Address
		msg.Dealer.Sex = self.RobotDealer.Sex
	} else {
		msg.Dealer.Total = self.Money
	}
	return &msg
}

func (self *Game_GoldBZW) GetPerson(uid int64) *Game_GoldBZW_Person {
	return self.PersonMgr[uid]
}

func (self *Game_GoldBZW) OnTime() {
	if lib.GetRobotMgr().GetRobotSet(self.room.Type).Dealer && len(self.Robot.Robots) > 0 && len(self.LstDeal) < 5 { //! 需要机器人上庄
		self.RobotUpDeal(self.Robot.Robots[lib.HF_GetRandom(len(self.Robot.Robots))])
	}

	if self.Time == 0 {
		return
	}

	for i := 0; i < len(self.Robot.Robots); i++ {
		if self.Robot.Robots[i].GetSeat() >= 0 {
			continue
		}
		if lib.HF_GetRandom(100) < 90 {
			continue
		}
		self.RobotSeat(lib.HF_GetRandom(8), self.Robot.Robots[i])
	}

	if time.Now().Unix() < self.Time {
		if self.Dealer == nil && self.RobotDealer == nil && lib.GetManyMgr().GetProperty(self.room.Type).SysNoBets == 1 {
			return
		}

		if self.Time-time.Now().Unix() >= int64(lib.GetManyMgr().GetProperty(self.room.Type).BetTime) {
			return
		}

		for i := 0; i < len(self.Robot.Robots); i++ {
			if self.Robot.Robots[i].GetSeat() == 100 { //! 庄家不能下注
				continue
			}

			if lib.HF_GetRandom(100) >= 100-lib.GetRobotMgr().GetRobotSet(self.room.Type).BetRate {
				continue
			}

			index, gold, _ := self.Robot.GameBets(self.Robot.Robots[i])
			if gold == 0 {
				continue
			}
			if self.Dealer != nil { //! 玩家庄判断是否能下
				dealwin := self.GetDealWinByPos(index, true)
				if self.Dealer.Total+self.Total+self.Robot.RobotTotal-dealwin < self.Dealer.Total/5 {
					self.Robot.GameBackBets(self.Robot.Robots[i], index, gold)
					continue
				}
			} else if self.RobotDealer != nil {
				dealwin := self.GetDealWinByPos(index, true)
				if self.RobotDealer.GetMoney()+self.Total+self.Robot.RobotTotal-dealwin < self.RobotDealer.GetMoney()/5 {
					self.Robot.GameBackBets(self.Robot.Robots[i], index, gold)
					continue
				}
			} else if lib.GetManyMgr().GetProperty(self.room.Type).DealChange == 1 {
				dealwin := self.GetDealWinByPos(index, true)
				if self.Money+self.Total+self.Robot.RobotTotal-dealwin < self.Money/5 {
					self.Robot.GameBackBets(self.Robot.Robots[i], index, gold)
					continue
				}
			}
			var msg Msg_GameGoldBZW_Bets
			msg.Uid = self.Robot.Robots[i].Id
			msg.Index = index
			msg.Gold = gold
			msg.Total = self.Robot.Robots[i].GetMoney()
			self.room.broadCastMsg("gamebzwbets", &msg)
		}
		return
	}

	if !self.room.Begin {
		self.OnBegin()
		return
	}
}

func (self *Game_GoldBZW) OnIsDealer(uid int64) bool {
	if self.Dealer != nil && self.Dealer == self.GetPerson(uid) {
		return true
	}
	return false
}

//! 申请无座玩家
func (self *Game_GoldBZW) GamePlayerList(uid int64) {
	var msg Msg_GameGoldBZW_List
	tmp := make(map[int64]Son_GameGoldBZW_Info)
	for _, value := range self.PersonMgr {
		if value.Seat >= 0 {
			continue
		}

		var node Son_GameGoldBZW_Info
		node.Uid = value.Uid
		node.Name = value.Name
		node.Total = value.Total
		node.Head = value.Head
		tmp[node.Uid] = node
	}
	for i := 0; i < len(self.Robot.Robots); i++ {
		if self.Robot.Robots[i].GetSeat() >= 0 {
			continue
		}

		var node Son_GameGoldBZW_Info
		node.Uid = self.Robot.Robots[i].Id
		node.Name = self.Robot.Robots[i].Name
		node.Total = self.Robot.Robots[i].GetMoney()
		node.Head = self.Robot.Robots[i].Head
		tmp[node.Uid] = node
	}
	for _, value := range tmp {
		msg.Info = append(msg.Info, value)
	}
	self.room.SendMsg(uid, "gameplayerlist", &msg)
}

//! 同步总分
func (self *Game_GoldBZW) SendTotal(uid int64, total int) {
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
func (self *Game_GoldBZW) SetTime(t int) {
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
func (self *Game_GoldBZW) ChageDeal() {
	if self.Dealer != nil {
		self.Dealer.Seat = -1
	} else if self.RobotDealer != nil {
		self.RobotDealer.SetSeat(-1)
	}

	self.Dealer = nil
	self.RobotDealer = nil
	for len(self.LstDeal) > 0 {
		if self.LstDeal[0].Robot != nil {
			find := false
			for i := 0; i < len(self.Robot.Robots); i++ {
				if self.Robot.Robots[i] == self.LstDeal[0].Robot {
					find = true
					break
				}
			}
			if !find { //! 要上庄的机器人已经走了
				self.LstDeal = self.LstDeal[1:]
				continue
			}
		}
		if self.LstDeal[0].GetTotal() >= lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney {
			if self.LstDeal[0].Person != nil {
				self.Dealer = self.LstDeal[0].Person
				self.Dealer.Seat = 100
				for i := 0; i < len(self.Seat); i++ {
					if self.Seat[i].Person == self.Dealer {
						var msg Msg_GameGoldBZW_UpdSeat
						msg.Index = i
						self.room.broadCastMsg("gamebzwseat", &msg)
						self.Seat[i].Person = nil
						break
					}
				}
			} else if self.LstDeal[0].Robot != nil {
				self.RobotDealer = self.LstDeal[0].Robot
				self.RobotDealer.SetSeat(100)
				for i := 0; i < len(self.Seat); i++ {
					if self.Seat[i].Robot == self.RobotDealer {
						var msg Msg_GameGoldBZW_UpdSeat
						msg.Index = i
						self.room.broadCastMsg("gamebzwseat", &msg)
						self.Seat[i].Robot = nil
						break
					}
				}
			}
			self.LstDeal = self.LstDeal[1:]
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
	} else if self.RobotDealer != nil {
		msg.Uid = self.RobotDealer.Id
		msg.Name = self.RobotDealer.Name
		msg.Head = self.RobotDealer.Head
		msg.Total = self.RobotDealer.GetMoney()
		msg.IP = self.RobotDealer.IP
		msg.Address = self.RobotDealer.Address
		msg.Sex = self.RobotDealer.Sex
	} else {
		msg.Total = self.Money
	}

	self.room.broadCastMsg("gamerob", &msg)
}

//! 是否下注了
func (self *Game_GoldBZW) OnIsBets(uid int64) bool {
	value, ok := self.PersonMgr[uid]
	if ok {
		return value.Bets > 0
	}
	return false
}

//! 结算所有人
func (self *Game_GoldBZW) OnBalance() {
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
