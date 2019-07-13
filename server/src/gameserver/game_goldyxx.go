package gameserver

import (
	"fmt"
	"lib"
	"math"
	"staticfunc"
	"time"
)

var GOLDYXX_BS []int = []int{2, 35, 2, 60, 30, 17, 12, 8, 7, 6, 6, 7, 8, 12, 17, 30, 60, 2, 2, 2, 2, 2, 2}
var GOLDYXX_DS []int = []int{0, 0, 0, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 1, 2, 3, 4, 5, 6}

//! 金币场记录
type Rec_YXX_Info struct {
	GameType int                  `json:"gametype"`
	Time     int64                `json:"time"` //! 记录时间
	Info     []Son_Rec_YXX_Person `json:"info"`
}
type Son_Rec_YXX_Person struct {
	Uid    int64   `json:"uid"`
	Name   string  `json:"name"`
	Head   string  `json:"head"`
	Score  int     `json:"score"`
	Result []int   `json:"result"`
	Bets   [23]int `json:"bets"`
}

type Msg_GameGoldYXX_Info struct {
	Begin   bool                    `json:"begin"`  //! 是否开始
	Time    int64                   `json:"time"`   //! 倒计时
	Seat    [8]Son_GameGoldYXX_Info `json:"info"`   //! 8个位置
	Bets    [23]int                 `json:"bets"`   //! 17个下注
	Dealer  Son_GameGoldYXX_Info    `json:"dealer"` //! 庄家
	Total   int                     `json:"total"`  //! 自己的钱
	Trend   [][]int                 `json:"trend"`  //! 走势
	IsDeal  bool                    `json:"isdeal"` //! 是否可下庄
	Money   []int                   `json:"money"`
	BetTime int                     `json:"bettime"` //! 下注时间
}

type Son_GameGoldYXX_Info struct {
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Total   int    `json:"total"`
	Ip      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type Msg_GameGoldYXX_Balance struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"` //! 当前金币
	Win   int   `json:"win"`   //! 赢了多少金币
}

type Msg_GameGoldYXX_End struct {
	Uid     int64  `json:"uid"` //! 大赢家
	Name    string `json:"name"`
	Head    string `json:"head"`
	Result  []int  `json:"result"`
	Money   []int  `json:"money"`
	BetTime int    `json:"bettime"` //! 下注时间
}

type Msg_GameGoldYXX_Total struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"`
}

//! 豹子王下注
type Msg_GameGoldYXX_Bets struct {
	Uid   int64 `json:"uid"`
	Index int   `json:"index"`
	Gold  int   `json:"gold"`
	Total int   `json:"total"`
}

//! 豹子王续压
type Msg_GameGoldYXX_Goon struct {
	Uid   int64   `json:"uid"`
	Gold  [23]int `json:"gold"`
	Total int     `json:"total"`
}

//! 上庄
type Msg_GameGoldYXX_Deal struct {
	Uid     int64  `json:"uid"`
	Head    string `json:"head"`
	Name    string `json:"name"`
	Total   int    `json:"total"`
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

//! 换座位
type Msg_GameGoldYXX_Seat struct {
	Index int `json:"index"`
}

//! 刷新座位
type Msg_GameGoldYXX_UpdSeat struct {
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
type Msg_GameGoldYXX_List struct {
	Info []Son_GameGoldYXX_Info `json:"info"`
}

//! 上庄和下庄
type Msg_GameGoldYXX_DealList struct {
	Type int                    `json:"type"` //! 0上庄  1下庄
	Info []Son_GameGoldYXX_Info `json:"info"`
}

//! 这局可能出现的点数
type GameGoldYXX_CanResult struct {
	DS int
	BZ bool
}

///////////////////////////////////////////////////////
type Game_GoldYXX_Person struct {
	Uid       int64   `json:"uid"`
	Gold      int     `json:"gold"`      //! 进来时候的钱
	Total     int     `json:"total"`     //! 当前的钱
	Win       int     `json:"win"`       //! 本局赢的钱
	Cost      int     `json:"cost"`      //! 手续费
	Bets      int     `json:"bets"`      //! 本局下了多少钱
	BetInfo   [23]int `json:"bets"`      //! 本局下的注
	BeBets    int     `json:"bebets"`    //! 上把下了多少钱
	BeBetInfo [23]int `json:"bebetinfo"` //! 上把的下注
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
//func (self *Game_GoldYXX_Person) GetBets() int {
//	total := 0
//	for i := 0; i < len(self.Bets); i++ {
//		total += self.Bets[i]
//	}

//	return total
//}

//! 同步金币
func (self *Game_GoldYXX_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

type Game_GoldYXXSeat struct {
	Person *Game_GoldYXX_Person
	Robot  *lib.Robot
}

func (self *Game_GoldYXXSeat) GetTotal() int {
	if self.Person != nil {
		return self.Person.Total
	} else if self.Robot != nil {
		return self.Robot.GetMoney()
	}
	return 0
}

type Game_GoldYXX struct {
	PersonMgr   map[int64]*Game_GoldYXX_Person   `json:"personmgr"`
	Bets        [23]map[*Game_GoldYXX_Person]int `json:"bets"`
	Result      []int                            `json:"result"`
	Dealer      *Game_GoldYXX_Person             `json:"dealer"`      //! 庄家
	RobotDealer *lib.Robot                       `json:"robotdealer"` //! 机器人庄
	Round       int                              `json:"round"`       //! 连庄轮数
	DownUid     int64                            `json:"downuid"`     //! 下庄的人
	Time        int64                            `json:"time"`
	LstDeal     []Game_GoldYXXSeat               `json:"lstdeal"` //! 上庄列表
	Seat        [8]Game_GoldYXXSeat              `json:"seat"`    //! 8个位置
	Total       int                              `json:"total"`   //! 这局一共下了多少钱
	Money       int                              `json:"money"`   //! 系统庄的钱
	Trend       [][]int                          `json:"trend"`   //! 走势
	Robot       lib.ManyGameRobot                //! 机器人结构
	Next        []int                            `json:"next"`
	BetTime     int                              `json:"bettime"`

	room *Room
}

func NewGame_GoldYXX() *Game_GoldYXX {
	game := new(Game_GoldYXX)
	game.PersonMgr = make(map[int64]*Game_GoldYXX_Person)
	for i := 0; i < len(game.Bets); i++ {
		game.Bets[i] = make(map[*Game_GoldYXX_Person]int)
	}
	for i := 0; i < 20; i++ {
		game.Trend = append(game.Trend, []int{lib.HF_GetRandom(6) + 1, lib.HF_GetRandom(6) + 1, lib.HF_GetRandom(6) + 1})
	}

	return game
}

func (self *Game_GoldYXX) OnInit(room *Room) {
	self.room = room
	self.Money = lib.GetManyMgr().GetProperty(self.room.Type).DealInitMoney
	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 13
	//! 载入机器人
	self.Robot.Init(23, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)
}

func (self *Game_GoldYXX) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldYXX) OnSendInfo(person *Person) {
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

	_person := new(Game_GoldYXX_Person)
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

func (self *Game_GoldYXX) OnMsg(msg *RoomMsg) {
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
	case "gamesetnext":
		self.Next = msg.V.(*staticfunc.Msg_SetDealNext).Next
		if len(self.Next) != 3 {
			self.Next = make([]int, 0)
		} else {
			for i := 0; i < len(self.Next); i++ {
				if self.Next[i] > 6 || self.Next[i] < 1 {
					self.Next = make([]int, 0)
					break
				}
			}
		}
	}
}

func (self *Game_GoldYXX) OnBegin() {
	if self.room.IsBye() {
		return
	}
	self.room.Begin = true

	if len(self.Next) == 3 {
		lib.HF_DeepCopy(&self.Result, &self.Next)
		self.Next = make([]int, 0)
		self.OnEnd()
		return
	}

	self.Result = make([]int, 0)
	if self.Dealer != nil { //! 玩家庄
		if self.Robot.RobotTotal == 0 { //! 没有机器人下注
			if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost == 102 { //! 系统算法
				lst := make([][]int, 0)
				winlst := make([][]int, 0)
				lostlst := make([][]int, 0)
				for i := 1; i <= 6; i++ {
					for j := i; j <= 6; j++ {
						for k := j; k <= 6; k++ {
							if i == j && i == k {
								continue
							}
							win := self.GetDealWin(i, j, k)
							if GetServer().YxxUserMoney[self.room.Type%240000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMin && GetServer().YxxUserMoney[self.room.Type%240000]+int64(win) <= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMax {
								lst = append(lst, []int{i, j, k})
							}
							if win > 0 {
								winlst = append(lst, []int{i, j, k})
							} else {
								lostlst = append(lst, []int{i, j, k})
							}
						}
					}
				}
				if len(lst) == 0 || lib.HF_GetRandom(100) < 10 {
					for i := 1; i <= 6; i++ {
						win := self.GetDealWin(i, i, i)
						if GetServer().YxxUserMoney[self.room.Type%240000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMin && GetServer().YxxUserMoney[self.room.Type%240000]+int64(win) <= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMax {
							lst = append(lst, []int{i, i, i})
						}
						if win > 0 {
							winlst = append(lst, []int{i, i, i})
						} else {
							lostlst = append(lst, []int{i, i, i})
						}
					}
				}
				if len(lst) == 0 { //! 如果都输钱，就随机，理论上不会发生这种情况
					if GetServer().YxxUserMoney[self.room.Type%240000] >= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMax && len(lostlst) > 0 { //! 一定输
						lib.GetLogMgr().Output(lib.LOG_DEBUG, "---------------- 玩家庄 无机器人 一定输 ")
						result := lostlst[lib.HF_GetRandom(len(lostlst))]
						self.GetResult(result)
					} else if GetServer().YxxUserMoney[self.room.Type%240000] <= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMin && len(winlst) > 0 { //! 一定赢
						lib.GetLogMgr().Output(lib.LOG_DEBUG, "---------------- 玩家庄 无机器人 一定赢 ")
						result := winlst[lib.HF_GetRandom(len(winlst))]
						self.GetResult(result)
					} else {
						lib.GetLogMgr().Output(lib.LOG_DEBUG, "---------------- 玩家庄 无机器人 亏了11111 ")
						for i := 0; i < 3; i++ {
							self.Result = append(self.Result, lib.HF_GetRandom(6)+1)
						}
					}
				} else {
					result := lst[lib.HF_GetRandom(len(lst))]
					self.GetResult(result)
				}

			} else if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost > 100 { //! 纯随机
				for i := 0; i < 3; i++ {
					self.Result = append(self.Result, lib.HF_GetRandom(6)+1)
				}
			} else { //! 设置概率
				iswin := lib.HF_GetRandom(100) < lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost
				lst := make([][]int, 0)
				smalllst := make([][]int, 0)
				biglst := make([][]int, 0)
				for i := 1; i <= 6; i++ {
					for j := i; j <= 6; j++ {
						for k := j; k <= 6; k++ {
							if i == j && i == k {
								win := self.GetDealWin(i, j, k)
								if iswin && win > 0 {
									smalllst = append(smalllst, []int{i, j, k})
								} else if !iswin && win < 0 {
									biglst = append(biglst, []int{i, j, k})
								}
							} else {
								win := self.GetDealWin(i, j, k)
								if iswin && win >= 0 {
									lst = append(lst, []int{i, j, k})
								} else if !iswin && win < 0 {
									if i+j+k == 4 || i+k+j == 17 {
										biglst = append(biglst, []int{i, j, k})
									} else if i+j+k == 5 || i+j+k == 16 {
										smalllst = append(smalllst, []int{i, j, k})
									} else {
										lst = append(lst, []int{i, j, k})
									}
								}
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
				if len(lst) == 0 {
					lib.GetLogMgr().Output(lib.LOG_ERROR, "豹子王庄家算错了")
					for i := 0; i < 3; i++ {
						self.Result = append(self.Result, lib.HF_GetRandom(6)+1)
					}
				} else {
					result := lst[lib.HF_GetRandom(len(lst))]
					self.GetResult(result)
				}
			}
		} else { //!  有机器人下注
			lst := make([][]int, 0)
			for i := 1; i <= 6; i++ {
				for j := i; j <= 6; j++ {
					for k := j; k <= 6; k++ {
						if i == j && i == k {
							continue
						}
						win := self.GetRobotWin(i, j, k)
						if lib.GetRobotMgr().GetRobotWin(self.room.Type)+win >= 0 {
							lst = append(lst, []int{i, j, k})
						}
					}
				}
			}
			if len(lst) == 0 || lib.HF_GetRandom(100) < 10 {
				for i := 1; i <= 6; i++ {
					win := self.GetRobotWin(i, i, i)
					if lib.GetRobotMgr().GetRobotWin(self.room.Type)+win >= 0 {
						lst = append(lst, []int{i, i, i})
					}
				}
			}
			if len(lst) == 0 { //! 如果都输钱，就随机，理论上不会发生这种情况
				for i := 0; i < 3; i++ {
					self.Result = append(self.Result, lib.HF_GetRandom(6)+1)
				}
			} else {
				result := lst[lib.HF_GetRandom(len(lst))]
				self.GetResult(result)
			}
		}
	} else { //! 系统庄
		lst := make([][]int, 0)
		for i := 1; i <= 6; i++ {
			for j := i; j <= 6; j++ {
				for k := j; k <= 6; k++ {
					if i == j && i == k {
						continue
					}
					win := self.GetDealWin(i, j, k)
					if GetServer().YxxSysMoney[self.room.Type%240000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
						lst = append(lst, []int{i, j, k})
					} else if win >= 0 {
						lst = append(lst, []int{i, j, k})
					}
				}
			}
		}
		for i := 1; i <= 6; i++ {
			win := self.GetDealWin(i, i, i)
			if GetServer().YxxSysMoney[self.room.Type%240000]+int64(win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
				lst = append(lst, []int{i, i, i})
			} else if win >= 0 {
				lst = append(lst, []int{i, i, i})
			}
		}
		if len(lst) == 0 {
			lib.GetLogMgr().Output(lib.LOG_ERROR, "豹子王亏了")
			for i := 0; i < 3; i++ {
				self.Result = append(self.Result, lib.HF_GetRandom(6)+1)
			}
		} else {
			result := lst[lib.HF_GetRandom(len(lst))]
			self.GetResult(result)
		}

	}

	self.OnEnd()
}

//! 开这三个点数17到22位置可以赢的钱
func (self *Game_GoldYXX) GetDealWinByPos17to22(num1 int, num2 int, num3 int, robot bool) int {
	lost := 0
	tmp := make(map[int]int)
	tmp[num1] += 1
	tmp[num2] += 2
	tmp[num3] += 3
	for key, value := range tmp {
		lost += self.GetMoneyPos(key+16, robot) * (GOLDYXX_BS[key+16] + value - 1)
	}

	return lost
}

//! 得到压中了这个点可以赢的钱
func (self *Game_GoldYXX) GetDealWinByPos(pos int, robot bool) int {
	lost := 0
	if pos < 17 { //!压龙虾蟹，这个点本身赢得钱单独来算
		lost = self.GetMoneyPos(pos, robot) * GOLDYXX_BS[pos] //! 压中这个点本身赢的钱
	}

	maxlost := 0
	if pos == 0 { //!　压小
		for i := 3; i <= 9; i++ {
			tmp := self.GetMoneyPos(i, robot) * GOLDYXX_BS[i]
			if tmp > maxlost {
				maxlost = tmp
			}
		}
		//!  模拟具体开的点数
		tmp := 0
		for i := 1; i <= 6; i++ {
			for j := i; j <= 6; j++ {
				for k := j; k <= 6; k++ {
					if i+j+k >= 11 {
						continue
					}
					if i == j && i == k {
						continue
					}
					_tmp := self.GetDealWinByPos17to22(i, j, k, robot)
					if tmp < _tmp {
						tmp = _tmp
					}
				}
			}
		}
		maxlost += tmp
	} else if pos == 1 { //! 压豹子
		for i := 3; i <= 16; i++ {
			if (i+1)%3 == 0 {
				tmp := self.GetMoneyPos(i, robot) * GOLDYXX_BS[i]
				if tmp > maxlost {
					maxlost = tmp
				}
			}
		}

		tmp := 0
		for i := 1; i <= 6; i++ {
			_tmp := self.GetDealWinByPos17to22(i, i, i, robot)
			if tmp < _tmp {
				tmp = _tmp
			}
		}
		maxlost += tmp
	} else if pos == 2 { //!  压大
		for i := 10; i <= 16; i++ {
			tmp := self.GetMoneyPos(i, robot) * GOLDYXX_BS[i]
			if tmp > maxlost {
				maxlost = tmp
			}
		}

		tmp := 0
		for i := 1; i <= 6; i++ {
			for j := i; j <= 6; j++ {
				for k := j; k <= 6; k++ {
					if i+j+k < 11 {
						continue
					}
					if i == j && i == k {
						continue
					}
					_tmp := self.GetDealWinByPos17to22(i, j, k, robot)
					if tmp < _tmp {
						tmp = _tmp
					}
				}
			}
		}
		maxlost += tmp
	} else if pos > 2 && pos < 17 { //!　压点数
		if (pos+1)%3 == 0 {
			maxlost = self.GetMoneyPos(1, robot) * GOLDYXX_BS[1]

			tmp := 0
			for i := 1; i <= 6; i++ {
				_tmp := self.GetDealWinByPos17to22(i, i, i, robot)
				if tmp < _tmp {
					tmp = _tmp
				}
			}
			maxlost += tmp
		} else {
			if pos+1 >= 4 && pos+1 <= 10 {
				maxlost = self.GetMoneyPos(0, robot) * GOLDYXX_BS[0]

				tmp := 0
				for i := 1; i <= 6; i++ {
					for j := i; j <= 6; j++ {
						for k := j; k <= 6; k++ {
							if i+j+k >= 11 {
								continue
							}
							if i == j && i == k {
								continue
							}
							_tmp := self.GetDealWinByPos17to22(i, j, k, robot)
							if tmp < _tmp {
								tmp = _tmp
							}
						}
					}
				}
				maxlost += tmp
			} else {
				maxlost = self.GetMoneyPos(2, robot) * GOLDYXX_BS[2]

				tmp := 0
				for i := 1; i <= 6; i++ {
					for j := i; j <= 6; j++ {
						for k := j; k <= 6; k++ {
							if i+j+k < 11 {
								continue
							}
							if i == j && i == k {
								continue
							}
							_tmp := self.GetDealWinByPos17to22(i, j, k, robot)
							if tmp < _tmp {
								tmp = _tmp
							}
						}
					}
				}
				maxlost += tmp
			}
		}
	} else { //!  压龙虾蟹
		tmp := 0
		for i := 1; i <= 3; i++ { //! 可能开出的个数 1-3
			if i == 1 {
				for j := 1; j <= 6; j++ {
					for k := j; k <= 6; k++ {
						_tmp := self.GetDealWinByPos17to22(pos-16, j, k, robot)
						if pos-16+k+j >= 11 {
							_tmp += self.GetMoneyPos(2, robot) * GOLDYXX_BS[2]
						} else {
							_tmp += self.GetMoneyPos(0, robot) * GOLDYXX_BS[0]
						}
						_tmp += self.GetMoneyPos(pos-16+j+k-1, robot) * GOLDYXX_BS[pos-16+j+k-1]
						if tmp < _tmp {
							tmp = _tmp
						}
					}
				}
			} else if i == 2 {
				for j := 1; j <= 6; j++ {
					_tmp := self.GetDealWinByPos17to22(pos-16, pos-16, j, robot)
					if (pos-16)*2+j >= 11 {
						_tmp += self.GetMoneyPos(2, robot) * GOLDYXX_BS[2]
					} else {
						_tmp += self.GetMoneyPos(0, robot) * GOLDYXX_BS[0]
					}
					_tmp += self.GetMoneyPos((pos-16)*2+j-1, robot) * GOLDYXX_BS[(pos-16)*2+j-1]
					if tmp < _tmp {
						tmp = _tmp
					}
				}
			} else if i == 3 {
				_tmp := self.GetDealWinByPos17to22(pos-16, pos-16, pos-16, robot)
				_tmp += self.GetMoneyPos(1, robot) * GOLDYXX_BS[1]
				_tmp += self.GetMoneyPos((pos-16)*3-1, robot) * GOLDYXX_BS[(pos-16)*3-1]
				if tmp < _tmp {
					tmp = _tmp
				}
			}
		}
		maxlost += tmp
	}

	return lost + maxlost
}

//! 得到庄家可以赢的钱
func (self *Game_GoldYXX) GetDealWin(num1 int, num2 int, num3 int) int {
	lost := 0
	if num1 == num2 && num1 == num3 {
		lost += self.GetMoneyPos(1, false) * GOLDYXX_BS[1]
	} else {
		if num1+num2+num3 >= 4 && num1+num2+num3 <= 10 {
			lost += self.GetMoneyPos(0, false) * GOLDYXX_BS[0]
		}
		if num1+num2+num3 >= 11 && num1+num2+num3 <= 17 {
			lost += self.GetMoneyPos(2, false) * GOLDYXX_BS[2]
		}
	}

	for i := 3; i < 17; i++ {
		if num1+num2+num3 == i+1 {
			lost += self.GetMoneyPos(i, false) * GOLDYXX_BS[i]
			break
		}
	}

	for i := 17; i < 23; i++ {
		tmp := 0
		if i-16 == num1 {
			tmp++
		}
		if i-16 == num2 {
			tmp++
		}
		if i-16 == num3 {
			tmp++
		}
		if tmp <= 0 {
			continue
		}
		lost += self.GetMoneyPos(i, false) * (GOLDYXX_BS[i] + tmp - 1)
	}

	if self.Dealer == nil {
		return self.Total - lost
	} else {
		return self.Total + self.Robot.RobotTotal - lost
	}
}

/*
func (self *Game_GoldYXX) GetDealWin(ds int, bz bool) int {
	lost := 0
	if bz {
		lost += self.GetMoneyPos(1) * GOLDYXX_BS[1]
	} else {
		if ds >= 4 && ds <= 10 {
			lost += self.GetMoneyPos(0) * GOLDYXX_BS[0]
		}
		if ds >= 11 && ds <= 17 {
			lost += self.GetMoneyPos(2) * GOLDYXX_BS[2]
		}
	}
	for i := 3; i < 17; i++ {
		if ds == i+1 {
			lost += self.GetMoneyPos(i) * GOLDYXX_BS[i]
			break
		}
	}

	if self.Dealer == nil {
		return self.Total - lost
	} else {
		return self.Total + self.Robot.RobotTotal - lost
	}
}
*/

//! 得到机器人可以赢的钱
func (self *Game_GoldYXX) GetRobotWin(num1 int, num2 int, num3 int) int {
	win := 0
	if num1 == num2 && num1 == num3 {
		win += self.GetMoneyPosByRobot(1) * GOLDYXX_BS[1]
	} else {
		if num1+num2+num3 >= 4 && num1+num2+num3 <= 10 {
			win += self.GetMoneyPosByRobot(0) * GOLDYXX_BS[0]
		}
		if num1+num2+num3 >= 11 && num1+num2+num3 <= 17 {
			win += self.GetMoneyPosByRobot(2) * GOLDYXX_BS[2]
		}
	}
	for i := 3; i < 17; i++ {
		if num1+num2+num3 == i+1 {
			win += self.GetMoneyPosByRobot(i) * GOLDYXX_BS[i]
			break
		}
	}
	for i := 17; i < 23; i++ {
		tmp := 0
		if i-16 == num1 {
			tmp++
		}
		if i-16 == num2 {
			tmp++
		}
		if i-16 == num3 {
			tmp++
		}
		if tmp <= 0 {
			continue
		}
		win += self.GetMoneyPosByRobot(i) * (GOLDYXX_BS[i] + tmp - 1)
	}
	return win - self.Robot.RobotTotal
}

/*
func (self *Game_GoldYXX) GetRobotWin(ds int, bz bool) int {
	win := 0
	if bz {
		win += self.GetMoneyPosByRobot(1) * GOLDYXX_BS[1]
	} else {
		if ds >= 4 && ds <= 10 {
			win += self.GetMoneyPosByRobot(0) * GOLDYXX_BS[0]
		}
		if ds >= 11 && ds <= 17 {
			win += self.GetMoneyPosByRobot(2) * GOLDYXX_BS[2]
		}
	}
	for i := 3; i < 17; i++ {
		if ds == i+1 {
			win += self.GetMoneyPosByRobot(i) * GOLDYXX_BS[i]
			break
		}
	}

	return win - self.Robot.RobotTotal
}
*/

//!根据三个点数得到具体排序
func (self *Game_GoldYXX) GetResult(result []int) {
	lib.GetLogMgr().Output(lib.LOG_DEBUG, " ------------- result : ", result)
	for i := 0; i < 3; i++ {
		index := lib.HF_GetRandom(len(result))
		self.Result = append(self.Result, result[index])
		copy(result[index:], result[index+1:])
		result = result[:len(result)-1]
	}
	//	self.Result[0] = 1
	//	self.Result[1] = 1
	//	self.Result[2] = 1
}

/*
//! 根据点数得到实际
func (self *Game_GoldYXX) GetResult(ds int, bz bool) {
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
*/

func (self *Game_GoldYXX) IsType() int { //! 0小 1豹子  2大
	if self.Result[0] == self.Result[1] && self.Result[0] == self.Result[2] {
		return 1
	}

	if self.Result[0]+self.Result[1]+self.Result[2] <= 10 {
		return 0
	}

	return 2
}

func (self *Game_GoldYXX) GameUpDeal(uid int64) {
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
		self.LstDeal = append(self.LstDeal, Game_GoldYXXSeat{person, nil})
	}
	person.Round = 0

	var msg Msg_GameGoldYXX_DealList
	msg.Type = 0
	msg.Info = make([]Son_GameGoldYXX_Info, 0)
	for i := 0; i < len(self.LstDeal); i++ {
		if self.LstDeal[i].Person != nil {
			msg.Info = append(msg.Info, Son_GameGoldYXX_Info{self.LstDeal[i].Person.Uid, self.LstDeal[i].Person.Name, self.LstDeal[i].Person.Head, self.LstDeal[i].Person.Total, self.LstDeal[i].Person.IP, self.LstDeal[i].Person.Address, self.LstDeal[i].Person.Sex})
		} else if self.LstDeal[i].Robot != nil {
			msg.Info = append(msg.Info, Son_GameGoldYXX_Info{self.LstDeal[i].Robot.Id, self.LstDeal[i].Robot.Name, self.LstDeal[i].Robot.Head, self.LstDeal[i].Robot.GetMoney(), self.LstDeal[i].Robot.IP, self.LstDeal[i].Robot.Address, self.LstDeal[i].Robot.Sex})
		}
	}
	self.room.SendMsg(uid, "gamebzwdeal", &msg)
}

//! 机器人上庄
func (self *Game_GoldYXX) RobotUpDeal(robot *lib.Robot) {
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
	self.LstDeal = append(self.LstDeal, Game_GoldYXXSeat{nil, robot})
}

func (self *Game_GoldYXX) GameReDeal(uid int64) {
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

	var msg Msg_GameGoldYXX_DealList
	msg.Type = 1
	msg.Info = make([]Son_GameGoldYXX_Info, 0)
	self.room.SendMsg(uid, "gamebzwdeal", &msg)
}

//! 坐下
func (self *Game_GoldYXX) GameSeat(uid int64, index int) {
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

	var msg Msg_GameGoldYXX_UpdSeat
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
func (self *Game_GoldYXX) RobotSeat(index int, robot *lib.Robot) {
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

	var msg Msg_GameGoldYXX_UpdSeat
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

func (self *Game_GoldYXX) GameBets(uid int64, index int, gold int) {
	if uid == 0 {
		return
	}

	if index < 0 || index > 22 {
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

	var msg Msg_GameGoldYXX_Bets
	msg.Uid = uid
	msg.Index = index
	msg.Gold = gold
	msg.Total = person.Total
	self.room.broadCastMsg("gamebzwbets", &msg)
}

//! 续压
func (self *Game_GoldYXX) GameGoOn(uid int64) {
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

	var msg Msg_GameGoldYXX_Goon
	msg.Uid = uid
	msg.Gold = person.BeBetInfo
	msg.Total = person.Total
	self.room.broadCastMsg("gamebzwgoon", &msg)
}

//! 结算
func (self *Game_GoldYXX) OnEnd() {
	self.room.Begin = false

	ds := 0
	for i := 0; i < 3; i++ {
		ds += self.Result[i]
	}

	trend := self.IsType()
	tmp := [][]int{self.Result}
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
					winmoney := value * GOLDYXX_BS[i]
					dealwin -= winmoney
					key.Win += winmoney
					key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
				}

				for key, value := range self.Robot.RobotsBet[i] {
					winmoney := value * GOLDYXX_BS[i]
					key.AddWin(winmoney)
					key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
					robotwin += winmoney
					if self.Dealer != nil {
						dealwin -= winmoney
					}
				}
			}
		} else if i >= 3 && i <= 16 {
			if self.Result[0]+self.Result[1]+self.Result[2] == GOLDYXX_DS[i] { //! 下注赢了
				for key, value := range self.Bets[i] {
					winmoney := value * GOLDYXX_BS[i]
					dealwin -= winmoney
					key.Win += winmoney
					key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
				}

				for key, value := range self.Robot.RobotsBet[i] {
					winmoney := value * GOLDYXX_BS[i]
					key.AddWin(winmoney)
					key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
					robotwin += winmoney
					if self.Dealer != nil {
						dealwin -= winmoney
					}
				}
			}
		} else { //! 计算鱼虾蟹
			num := 0
			for j := 0; j < 3; j++ {
				if self.Result[j] == i-16 {
					num++
				}
			}
			if num > 0 {
				for key, value := range self.Bets[i] {
					winmoney := value * (GOLDYXX_BS[i] + num - 1)
					dealwin -= winmoney
					key.Win += winmoney
					key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
				}

				for key, value := range self.Robot.RobotsBet[i] {
					winmoney := value * (GOLDYXX_BS[i] + num - 1)
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

			var msg Msg_GameGoldYXX_Balance
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
			var record Rec_YXX_Info
			record.Time = time.Now().Unix()
			record.GameType = self.room.Type
			var rec Son_Rec_YXX_Person
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
					var msg Msg_GameGoldYXX_Balance
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
		GetServer().SetYxxUserMoney(self.room.Type%240000, GetServer().YxxUserMoney[self.room.Type%240000]+int64(dealwin))
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
		GetServer().SetYxxSysMoney(self.room.Type%240000, GetServer().YxxSysMoney[self.room.Type%240000]+int64(_dealwin))
		dealwin -= robotwin
		if dealwin > 0 {
			bl := lib.GetManyMgr().GetProperty(self.room.Type).Cost
			cost := int(math.Ceil(float64(dealwin) * bl / 100.0))
			dealwin -= cost
		}
	}

	if self.Dealer != nil {
		var record Rec_YXX_Info
		record.Time = time.Now().Unix()
		record.GameType = self.room.Type
		var rec Son_Rec_YXX_Person
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
		var msg Msg_GameGoldYXX_Balance
		if self.Dealer != nil {
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

	//! 下注时间
	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 13
	self.SetTime(self.BetTime)

	//! 总结算
	{
		var msg Msg_GameGoldYXX_End
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
				var msg Msg_GameGoldYXX_UpdSeat
				msg.Index = j
				self.room.broadCastMsg("gamebzwseat", &msg)
				break
			}
		}
		delete(self.PersonMgr, key)
	}

	//! 返回机器人
	self.Robot.Init(23, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)

	for i := 0; i < len(self.Bets); i++ {
		self.Bets[i] = make(map[*Game_GoldYXX_Person]int)
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
			var msg Msg_GameGoldYXX_UpdSeat
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

func (self *Game_GoldYXX) OnBye() {
}

func (self *Game_GoldYXX) OnExit(uid int64) {
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
func (self *Game_GoldYXX) GetMoneyPos(index int, robot bool) int {
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
func (self *Game_GoldYXX) GetMoneyPosByRobot(index int) int {
	total := 0
	for _, value := range self.Robot.RobotsBet[index] {
		total += value
	}
	return total
}

func (self *Game_GoldYXX) getInfo(uid int64, total int) *Msg_GameGoldYXX_Info {
	var msg Msg_GameGoldYXX_Info
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

func (self *Game_GoldYXX) GetPerson(uid int64) *Game_GoldYXX_Person {
	return self.PersonMgr[uid]
}

func (self *Game_GoldYXX) OnTime() {
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
			var msg Msg_GameGoldYXX_Bets
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

func (self *Game_GoldYXX) OnIsDealer(uid int64) bool {
	if self.Dealer != nil && self.Dealer == self.GetPerson(uid) {
		return true
	}
	return false
}

//! 申请无座玩家
func (self *Game_GoldYXX) GamePlayerList(uid int64) {
	var msg Msg_GameGoldYXX_List
	for _, value := range self.PersonMgr {
		if value.Seat >= 0 {
			continue
		}

		var node Son_GameGoldYXX_Info
		node.Uid = value.Uid
		node.Name = value.Name
		node.Total = value.Total
		node.Head = value.Head
		msg.Info = append(msg.Info, node)
	}
	for i := 0; i < len(self.Robot.Robots); i++ {
		if self.Robot.Robots[i].GetSeat() >= 0 {
			continue
		}

		var node Son_GameGoldYXX_Info
		node.Uid = self.Robot.Robots[i].Id
		node.Name = self.Robot.Robots[i].Name
		node.Total = self.Robot.Robots[i].GetMoney()
		node.Head = self.Robot.Robots[i].Head
		msg.Info = append(msg.Info, node)
	}
	self.room.SendMsg(uid, "gameplayerlist", &msg)
}

//! 同步总分
func (self *Game_GoldYXX) SendTotal(uid int64, total int) {
	var msg Msg_GameGoldYXX_Total
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
func (self *Game_GoldYXX) SetTime(t int) {
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
func (self *Game_GoldYXX) ChageDeal() {
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
						var msg Msg_GameGoldYXX_UpdSeat
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
						var msg Msg_GameGoldYXX_UpdSeat
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
func (self *Game_GoldYXX) OnIsBets(uid int64) bool {
	value, ok := self.PersonMgr[uid]
	if ok {
		return value.Bets > 0
	}
	return false
}

//! 结算所有人
func (self *Game_GoldYXX) OnBalance() {
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
