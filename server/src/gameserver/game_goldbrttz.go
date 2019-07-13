package gameserver

import (
	"fmt"
	"lib"
	"math"
	"sort"
	"staticfunc"
	"time"
)

func (s lstLost) Len() int           { return len(s) }
func (s lstLost) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s lstLost) Less(i, j int) bool { return s[i].Win > s[j].Win }

type JS_Lost struct {
	Result [4][]int `json:"result"`
	Win    int      `json:"win"`
}

type lstLost []JS_Lost

//! 金币场记录
type Rec_BrTTZ_Info struct {
	GameType int                    `json:"gametype"`
	Time     int64                  `json:"time"` //! 记录时间
	Info     []Son_Rec_BrTTZ_Person `json:"info"`
}
type Son_Rec_BrTTZ_Person struct {
	Uid    int64    `json:"uid"`
	Name   string   `json:"name"`
	Head   string   `json:"head"`
	Score  int      `json:"score"`
	Result [4][]int `json:"result"`
	Bets   [3]int   `json:"bets"`
}

type Msg_GameGoldBrTTZ_Info struct {
	Begin   bool                       `json:"begin"`  //! 是否开始
	Time    int64                      `json:"time"`   //! 倒计时
	Seat    [12]Son_GameGoldBrTTZ_Info `json:"info"`   //! 12个位置
	Bets    [3]int                     `json:"bets"`   //! 3个下注
	Dealer  Son_GameGoldBrTTZ_Info     `json:"dealer"` //! 庄家
	Total   int                        `json:"total"`  //! 自己的钱
	Trend   [3][]int                   `json:"trend"`  //! 走势
	IsDeal  bool                       `json:"isdeal"` //! 是否可下庄
	Result  [4][]int                   `json:"result"`
	Money   []int                      `json:"money"`
	BetTime int                        `json:"bettime"`
}

type Son_GameGoldBrTTZ_Info struct {
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Total   int    `json:"total"`
	Ip      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type Msg_GameGoldBrTTZ_Next struct {
	Next []int `json:"next"`
}

type Msg_GameGoldBrTTZ_End struct {
	Uid     int64    `json:"uid"` //! 大赢家
	Name    string   `json:"name"`
	Head    string   `json:"head"`
	Result  [4][]int `json:"result"`
	Next    [4][]int `json:"next"`
	Money   []int    `json:"money"`
	BetTime int      `json:"bettime"`
}

//! 推筒子续压
type Msg_GameGoldBrTTZ_Goon struct {
	Uid   int64  `json:"uid"`
	Gold  [3]int `json:"gold"`
	Total int    `json:"total"`
}

//! 推筒子抢庄
type Msg_GameGoldBrTTZ_Deal struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"`
}

type Game_GoldBrTTZ_Robot struct {
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`   //! 名字
	Imgurl string `json:"imgurl"` //! 头像
	Sex    int    `json:"sex"`    //! 性别
	IP     string `json:"ip"`
}

///////////////////////////////////////////////////////
type Game_GoldBrTTZ_Deal struct {
	Person   *Game_GoldBrTTZ_Person `json:"dealer"`
	Robot    *lib.Robot             `json:"robot"`
	DealType int                    `json:"dealtype"`
}

func (self *Game_GoldBrTTZ_Deal) GetTotal() int {
	if self.Person != nil {
		return self.Person.Total
	} else if self.Robot != nil {
		return self.Robot.GetMoney()
	}
	return 0
}

///////////////////////////////////////////////////////
type Game_GoldBrTTZ_Person struct {
	Uid       int64  `json:"uid"`
	Gold      int    `json:"gold"`      //! 进来时候的钱
	Total     int    `json:"total"`     //! 当前的钱
	Win       int    `json:"win"`       //! 本局赢的钱
	Cost      int    `json:"cost"`      //! 手续费
	Bets      int    `json:"bets"`      //! 本局下了多少钱
	BetInfo   [3]int `json:"bets"`      //! 本局下的注
	BeBets    int    `json:"bebets"`    //! 上把下了多少钱
	BeBetInfo [3]int `json:"bebetinfo"` //! 上把的下注
	Name      string `json:"name"`      //! 名字
	Head      string `json:"head"`      //! 头像
	Online    bool   `json:"online"`
	Round     int    `json:"round"` //! 不下注轮数
	Seat      int    `json:"seat"`  //! 0-11有座  -1无座  100庄家
	IP        string `json:"ip"`
	Address   string `json:"address"`
	Sex       int    `json:"sex"`
}

//! 同步金币
func (self *Game_GoldBrTTZ_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

type Game_GoldBrTTZSeat struct {
	Person *Game_GoldBrTTZ_Person
	Robot  *lib.Robot
}

type Game_GoldBrTTZ struct {
	PersonMgr   map[int64]*Game_GoldBrTTZ_Person  `json:"personmgr"`
	Bets        [3]map[*Game_GoldBrTTZ_Person]int `json:"bets"`
	Result      [4][]int                          `json:"result"`
	Dealer      *Game_GoldBrTTZ_Person            `json:"dealer"`      //! 庄家
	RobotDealer *lib.Robot                        `json:"robotdealer"` //! 机器人庄
	DealType    int                               `json:"dealtype"`    //! 庄类型
	Round       int                               `json:"round"`       //! 连庄轮数
	DownUid     int64                             `json:"downuid"`     //! 下庄的人
	Time        int64                             `json:"time"`
	LstDeal     []Game_GoldBrTTZ_Deal             `json:"lstdeal"` //! 上庄列表
	Seat        [12]Game_GoldBrTTZSeat            `json:"seat"`    //! 12个位置
	Total       int                               `json:"total"`   //! 这局一共下了多少钱
	Money       int                               `json:"money"`   //! 系统庄的钱
	Trend       [3][]int                          `json:"trend"`   //! 走势
	Next        []int                             `json:"next"`
	BetTime     int                               `json:"bettime"`
	Robot       lib.ManyGameRobot                 //! 机器人结构

	room *Room
}

func NewGame_GoldBrTTZ() *Game_GoldBrTTZ {
	game := new(Game_GoldBrTTZ)
	game.PersonMgr = make(map[int64]*Game_GoldBrTTZ_Person)
	for i := 0; i < len(game.Bets); i++ {
		game.Bets[i] = make(map[*Game_GoldBrTTZ_Person]int)
	}
	for i := 0; i < 8; i++ {
		game.Trend[0] = append(game.Trend[0], lib.HF_GetRandom(2))
		game.Trend[1] = append(game.Trend[1], lib.HF_GetRandom(2))
		game.Trend[2] = append(game.Trend[2], lib.HF_GetRandom(2))
	}

	return game
}

func (self *Game_GoldBrTTZ) OnInit(room *Room) {
	self.room = room
	self.Money = lib.GetManyMgr().GetProperty(self.room.Type).DealInitMoney
	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 10
	//! 载入机器人
	self.Robot.Init(3, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)
}

func (self *Game_GoldBrTTZ) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldBrTTZ) OnSendInfo(person *Person) {
	if self.Time == 0 {
		self.GameCard()
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
		person.SendMsg("gamegoldbrttzinfo", self.getInfo(person.Uid, value.Total))
		return
	}

	_person := new(Game_GoldBrTTZ_Person)
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
	person.SendMsg("gamegoldbrttzinfo", self.getInfo(person.Uid, person.Gold))
}

func (self *Game_GoldBrTTZ) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold": //! 同步金币
		person := self.GetPerson(msg.V.(*staticfunc.Msg_SynchroGold).Uid)
		if person != nil {
			person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal(person.Uid, person.Total)
		}
	case "gamebrttzbets": //! 下注
		self.GameBets(msg.Uid, msg.V.(*Msg_GameGoldBZW_Bets).Index, msg.V.(*Msg_GameGoldBZW_Bets).Gold)
	//case "gamebrttzgoon":
	//	self.GameGoOn(msg.Uid)
	case "gamerob": //! 上庄
		self.GameUpDeal(msg.Uid, 0)
	case "gamebrttzrob":
		self.GameUpDeal(msg.Uid, msg.V.(*Msg_GameGoldBrTTZ_Deal).Type)
	case "gameredeal": //! 下庄
		self.GameReDeal(msg.Uid)
	case "gamebrttzseat":
		self.GameSeat(msg.Uid, msg.V.(*Msg_GameGoldBZW_Seat).Index)
	case "gameplayerlist":
		self.GamePlayerList(msg.Uid)
	case "gamesetnext":
		self.Next = msg.V.(*staticfunc.Msg_SetDealNext).Next
	}
}

//! 发牌
func (self *Game_GoldBrTTZ) GameCard() {
	card := NewMah_TTZ()
	for i := 0; i < len(self.Result); i++ {
		self.Result[i] = card.Deal(2)
	}
}

func (self *Game_GoldBrTTZ) RobotNeedWin() {
	lstWin := make(lstLost, 0)
	lstKill := make(lstLost, 0)
	for i := 10; i < 20; i++ {
		self.Result[0][1] = 0
		self.Result[1][1] = 0
		self.Result[2][1] = 0
		self.Result[3][1] = 0

		value := i
		if value == 10 {
			value = 37
		}
		self.Result[0][1] = value

		for j := 10; j < 20; j++ {
			value := j
			if value == 10 {
				value = 37
			}
			if !self.IsOk(value) {
				continue
			}
			self.Result[1][1] = value

			for k := 10; k < 20; k++ {
				value := k
				if value == 10 {
					value = 37
				}
				if !self.IsOk(value) {
					continue
				}
				self.Result[2][1] = value

				for l := 10; l < 20; l++ {
					value := l
					if value == 10 {
						value = 37
					}
					if !self.IsOk(value) {
						continue
					}
					self.Result[3][1] = value

					win, kill := self.GetRobotWin()
					if lib.GetRobotMgr().GetRobotWin(self.room.Type)+win < 0 {
						continue
					}
					var node JS_Lost
					lib.HF_DeepCopy(&node.Result, &self.Result)
					node.Win = win
					if kill {
						lstKill = append(lstKill, node)
					} else {
						lstWin = append(lstWin, node)
					}
				}
			}
		}
	}

	if len(lstWin) > 0 {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "赢1")
		self.Result = lstWin[lib.HF_GetRandom(len(lstWin))].Result
	} else if len(lstKill) > 0 {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "赢2")
		self.Result = lstKill[lib.HF_GetRandom(len(lstKill))].Result
	} else {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "哈哈哈啊1")
		card := NewMah_TTZ()
		for i := 0; i < 4; i++ {
			card.Del(self.Result[i][0])
		}
		for i := 0; i < 4; i++ {
			self.Result[i][1] = card.Deal(1)[0]
		}
	}
}

//!
func (self *Game_GoldBrTTZ) WinOrLost(iswin bool) {
	if lib.HF_GetRandom(100) < 30 {
		if iswin {
			lstWin := make(lstLost, 0)
			lstKill := make(lstLost, 0)
			for i := 10; i < 20; i++ {
				self.Result[0][1] = 0
				self.Result[1][1] = 0
				self.Result[2][1] = 0
				self.Result[3][1] = 0

				value := i
				if value == 10 {
					value = 37
				}
				self.Result[0][1] = value

				for j := 10; j < 20; j++ {
					value := j
					if value == 10 {
						value = 37
					}
					if !self.IsOk(value) {
						continue
					}
					self.Result[1][1] = value

					for k := 10; k < 20; k++ {
						value := k
						if value == 10 {
							value = 37
						}
						if !self.IsOk(value) {
							continue
						}
						self.Result[2][1] = value

						for l := 10; l < 20; l++ {
							value := l
							if value == 10 {
								value = 37
							}
							if !self.IsOk(value) {
								continue
							}
							self.Result[3][1] = value

							win, kill := self.GetResultWin()
							if win < 0 {
								continue
							}
							var node JS_Lost
							lib.HF_DeepCopy(&node.Result, &self.Result)
							node.Win = win
							if kill {
								lstKill = append(lstKill, node)
							} else {
								lstWin = append(lstWin, node)
							}
						}
					}
				}
			}

			if len(lstWin) > 0 {
				lib.GetLogMgr().Output(lib.LOG_ERROR, "赢1")
				self.Result = lstWin[lib.HF_GetRandom(len(lstWin))].Result
			} else if len(lstKill) > 0 {
				lib.GetLogMgr().Output(lib.LOG_ERROR, "赢2")
				self.Result = lstKill[lib.HF_GetRandom(len(lstKill))].Result
			} else {
				lib.GetLogMgr().Output(lib.LOG_ERROR, "哈哈哈啊1")
				card := NewMah_TTZ()
				for i := 0; i < 4; i++ {
					card.Del(self.Result[i][0])
				}
				for i := 0; i < 4; i++ {
					self.Result[i][1] = card.Deal(1)[0]
				}
			}
		} else {
			lstWin := make(lstLost, 0)
			for i := 10; i < 20; i++ {
				self.Result[0][1] = 0
				self.Result[1][1] = 0
				self.Result[2][1] = 0
				self.Result[3][1] = 0

				value := i
				if value == 10 {
					value = 37
				}
				self.Result[0][1] = value

				for j := 10; j < 20; j++ {
					self.Result[1][1] = 0
					self.Result[2][1] = 0
					self.Result[3][1] = 0

					value := j
					if value == 10 {
						value = 37
					}
					if !self.IsOk(value) {
						continue
					}
					self.Result[1][1] = value

					for k := 10; k < 20; k++ {
						self.Result[2][1] = 0
						self.Result[3][1] = 0

						value := k
						if value == 10 {
							value = 37
						}
						if !self.IsOk(value) {
							continue
						}
						self.Result[2][1] = value

						for l := 10; l < 20; l++ {
							self.Result[3][1] = 0

							value := l
							if value == 10 {
								value = 37
							}
							if !self.IsOk(value) {
								continue
							}
							self.Result[3][1] = value

							win, _ := self.GetResultWin()
							if win > 0 {
								continue
							}
							var node JS_Lost
							lib.HF_DeepCopy(&node.Result, &self.Result)
							node.Win = win
							lstWin = append(lstWin, node)
						}
					}
				}
			}
			sort.Sort(lstLost(lstWin))
			if len(lstWin) > 0 {
				lib.GetLogMgr().Output(lib.LOG_ERROR, "输1")
				self.Result = lstWin[lib.HF_GetRandom(len(lstWin)/2)].Result
			} else {
				lib.GetLogMgr().Output(lib.LOG_ERROR, "哈哈哈啊2")
				card := NewMah_TTZ()
				for i := 0; i < 4; i++ {
					card.Del(self.Result[i][0])
				}
				for i := 0; i < 4; i++ {
					self.Result[i][1] = card.Deal(1)[0]
				}
			}
		}
	} else {
		card := NewMah_TTZ()
		for i := 0; i < 4; i++ {
			card.Del(self.Result[i][0])
		}
		if iswin {
			//! 尽量不通杀
			find := false
			for i := 0; i < 2000; i++ {
				win, kill := self.GetResultWin()
				if kill {
					continue
				}
				if win >= 0 {
					find = true
					break
				}
				_card := NewMah_TTZ()
				lib.HF_DeepCopy(_card, card)
				for j := 0; j < 4; j++ {
					self.Result[j][1] = _card.Deal(1)[0]
				}
			}
			if !find {
				for i := 0; i < 100; i++ {
					win, _ := self.GetResultWin()
					if win >= 0 {
						break
					}
					_card := NewMah_TTZ()
					lib.HF_DeepCopy(_card, card)
					for j := 0; j < 4; j++ {
						self.Result[j][1] = _card.Deal(1)[0]
					}
				}
			}
		} else {
			//! 尽量不通赔
			find := false
			for i := 0; i < 2000; i++ {
				win, kill := self.GetResultWin()
				if kill {
					continue
				}
				if win <= 0 {
					find = true
					break
				}
				_card := NewMah_TTZ()
				lib.HF_DeepCopy(_card, card)
				for j := 0; j < 4; j++ {
					self.Result[j][1] = _card.Deal(1)[0]
				}
			}
			if !find {
				for i := 0; i < 100; i++ {
					win, _ := self.GetResultWin()
					if win <= 0 {
						break
					}
					_card := NewMah_TTZ()
					lib.HF_DeepCopy(_card, card)
					for j := 0; j < 4; j++ {
						self.Result[j][1] = _card.Deal(1)[0]
					}
				}
			}
		}
	}
}

func (self *Game_GoldBrTTZ) OnBegin() {
	if self.room.IsBye() {
		return
	}

	if len(self.Next) > 0 {
		card := NewMah_TTZ()
		for i := 0; i < 4; i++ {
			card.Del(self.Result[i][0])
		}
		for i := 0; i < lib.HF_MinInt(len(self.Next), 4); i++ {
			value := 0
			if self.Next[i] == 0 {
				value = 37
			} else if self.Next[i] >= 1 && self.Next[i] <= 9 {
				value = 10 + self.Next[i]
			}
			self.Result[i][1] = card.Draw4(value)
		}

		self.Next = make([]int, 0)
		self.OnEnd()
		return
	}

	self.room.Begin = true
	if self.Dealer != nil { //! 玩家庄
		if self.Robot.RobotTotal == 0 {
			if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost == 102 { //! 平衡模式
				win, _ := self.GetResultWin()
				if GetServer().BrTTZUsrMoney[self.room.Type%60000]+int64(win) > lib.GetManyMgr().GetProperty(self.room.Type).PlayerMax { //! 玩家需要输
					self.WinOrLost(false)
				} else if GetServer().BrTTZUsrMoney[self.room.Type%60000]+int64(win) < lib.GetManyMgr().GetProperty(self.room.Type).PlayerMin { //! 玩家需要赢
					self.WinOrLost(true)
				}
			} else if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost == 101 { //! 随机模式

			} else if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost <= 100 { //! 概率模式
				if lib.HF_GetRandom(100) < lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost { //! 玩家庄胜利
					self.WinOrLost(true)
				} else { //! 玩家庄失败
					self.WinOrLost(false)
				}
			}
		} else {
			win, _ := self.GetRobotWin()
			if lib.GetRobotMgr().GetRobotWin(self.room.Type)+win < 0 {
				self.RobotNeedWin()
			}
		}
	} else { //! 系统庄
		win, _ := self.GetResultWin()
		if GetServer().BrTTZMoney[self.room.Type%60000]+int64(win) < lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
			self.WinOrLost(true)
		} else if GetServer().BrTTZMoney[self.room.Type%60000]+int64(win) > lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
			self.WinOrLost(false)
		}
		//willwin := true
		//if GetServer().BrTTZMoney[self.room.Type%60000] < lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
		//	willwin = true
		//} else if GetServer().BrTTZMoney[self.room.Type%60000] > lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
		//	willwin = false
		//} else {
		//	sum := lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax - lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin
		//	s1 := GetServer().BrTTZMoney[self.room.Type%60000] - lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin
		//	pro := 1000000 * s1 / sum
		//	if lib.HF_GetRandom(1000000) < int(pro) {
		//		willwin = false
		//	} else {
		//		willwin = true
		//	}
		//}
		//self.WinOrLost(willwin)
	}

	self.OnEnd()
}

//! 是否可以换
func (self *Game_GoldBrTTZ) IsOk(card int) bool {
	num := 0
	for i := 0; i < len(self.Result); i++ {
		if self.Result[i][0] == card {
			num++
		}
		if self.Result[i][1] == card {
			num++
		}
		if num >= 4 {
			return false
		}
	}

	return true
}

func (self *Game_GoldBrTTZ) GameUpDeal(uid int64, _type int) {
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
		self.DealType = _type
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
		self.LstDeal = append(self.LstDeal, Game_GoldBrTTZ_Deal{person, nil, _type})
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
	self.room.SendMsg(uid, "gamebrttzdeal", &msg)
}

//! 机器人上庄
func (self *Game_GoldBrTTZ) RobotUpDeal(robot *lib.Robot) {
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
	self.LstDeal = append(self.LstDeal, Game_GoldBrTTZ_Deal{nil, robot, 0})
}

func (self *Game_GoldBrTTZ) GameReDeal(uid int64) {
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
	self.room.SendMsg(uid, "gamebrttzdeal", &msg)
}

//! 坐下
func (self *Game_GoldBrTTZ) GameSeat(uid int64, index int) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	if index < 0 || index >= len(self.Seat) {
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
	self.room.broadCastMsg("gamebrttzseat", &msg)
}

//! 机器人坐下
func (self *Game_GoldBrTTZ) RobotSeat(index int, robot *lib.Robot) {
	if index < 0 || index >= len(self.Seat) {
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
	self.room.broadCastMsg("gamebrttzseat", &msg)
}

func (self *Game_GoldBrTTZ) GameBets(uid int64, index int, gold int) {
	if uid == 0 {
		return
	}

	if index < 0 || index >= len(self.Bets) {
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

	bl := 20
	if self.Result[0][0] != 37 {
		bl = 15
	}
	if (person.Bets+gold)*(bl-1) > (person.Total - gold) {
		self.room.SendErr(uid, "您的金币可能不够赔率，请前往充值。")
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
		if dealmoney > 0 {
			dealwin := self.GetMaxLost(true)
			if dealmoney+self.Total+self.Robot.RobotTotal-dealwin < dealmoney/5 {
				self.Total -= gold
				self.Bets[index][person] -= gold
				self.room.SendErr(uid, "庄家已到最大赔率")
				self.OnBegin()
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
	self.room.broadCastMsg("gamebrttzbets", &msg)
}

//! 得到当前庄家输赢
func (self *Game_GoldBrTTZ) GetResultWin() (int, bool) {
	num := 0
	dealwin := 0
	dealds, dealbs := GetTTZResult(self.Result[0])
	for i := 1; i < len(self.Result); i++ {
		xiands, xianbs := GetTTZResult(self.Result[i])
		if dealds >= xiands { //! 庄赢
			for _, value := range self.Bets[i-1] {
				dealwin += value * dealbs
			}
			if self.Dealer != nil { //! 玩家庄要计算机器人
				for _, value := range self.Robot.RobotsBet[i-1] {
					dealwin += value * dealbs
				}
			}
			num++
		} else {
			for _, value := range self.Bets[i-1] {
				dealwin -= value * xianbs
			}
			if self.Dealer != nil { //! 玩家庄要计算机器人
				for _, value := range self.Robot.RobotsBet[i-1] {
					dealwin -= value * xianbs
				}
			}
			num--
		}
	}
	return dealwin, (num == 3 || num == -3)
}

//! 得到机器人输赢
func (self *Game_GoldBrTTZ) GetRobotWin() (int, bool) {
	num := 0
	robotwin := 0
	dealds, dealbs := GetTTZResult(self.Result[0])
	for i := 1; i < len(self.Result); i++ {
		xiands, xianbs := GetTTZResult(self.Result[i])
		if dealds >= xiands { //! 庄赢
			for _, value := range self.Robot.RobotsBet[i-1] {
				robotwin -= value * dealbs
			}
			num++
		} else {
			for _, value := range self.Robot.RobotsBet[i-1] {
				robotwin += value * xianbs
			}
			num--
		}
	}
	return robotwin, (num == 3 || num == -3)
}

//! 结算
func (self *Game_GoldBrTTZ) OnEnd() {
	self.room.Begin = false

	dealwin := 0
	robotwin := 0
	dealds, dealbs := GetTTZResult(self.Result[0])
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "庄家:", dealds, ",", dealbs)
	for i := 1; i < len(self.Result); i++ {
		xiands, xianbs := GetTTZResult(self.Result[i])
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "闲家:", xiands, ",", xianbs)
		if dealds >= xiands { //! 庄赢
			for key, value := range self.Bets[i-1] {
				dealwin += value * dealbs
				key.Win -= (value*dealbs - value)
			}

			for key, value := range self.Robot.RobotsBet[i-1] {
				if self.Dealer != nil {
					dealwin += value * dealbs
				}
				key.AddWin(-(value*dealbs - value))
				robotwin -= value * dealbs
			}

			tmp := []int{0}
			tmp = append(tmp, self.Trend[i-1]...)
			if len(tmp) > 8 {
				tmp = tmp[0:8]
			}
			self.Trend[i-1] = tmp
		} else {
			for key, value := range self.Bets[i-1] {
				dealwin -= value * xianbs
				key.Win += (value*xianbs + value)
				key.Cost += int(math.Ceil(float64(value*xianbs) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
			}

			for key, value := range self.Robot.RobotsBet[i-1] {
				if self.Dealer != nil {
					dealwin -= value * xianbs
				}
				key.AddWin((value*xianbs + value))
				key.AddCost(int(math.Ceil(float64(value*xianbs) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
				robotwin += value * xianbs
			}

			tmp := []int{1}
			tmp = append(tmp, self.Trend[i-1]...)
			if len(tmp) > 8 {
				tmp = tmp[0:8]
			}
			self.Trend[i-1] = tmp
		}
	}

	var bigwin *GameGold_BigWin = nil
	for _, value := range self.PersonMgr {
		if value.Win > 0 { //! 赢了要手续费
			value.Win -= value.Cost
			GetServer().SqlAgentGoldLog(value.Uid, value.Cost, self.room.Type)
			GetServer().SqlAgentBillsLog(value.Uid, value.Cost/2, self.room.Type)
		} else if value.Win-value.Bets < 0 {
			cost := int(math.Ceil(float64(value.Win-value.Bets) * float64(lib.GetManyMgr().GetProperty(self.room.Type).Cost) / 200.0))
			GetServer().SqlAgentBillsLog(value.Uid, cost, self.room.Type)
		}
		if value.Win != 0 {
			value.Total += value.Win
			var msg Msg_GameGoldBZW_Balance
			msg.Uid = value.Uid
			msg.Total = value.Total
			msg.Win = value.Win
			find := false
			for j := 0; j < len(self.Seat); j++ {
				if self.Seat[j].Person == value {
					self.room.broadCastMsg("gamegoldbrttzbalance", &msg)
					find = true
					break
				}
			}
			if !find {
				self.room.SendMsg(value.Uid, "gamegoldbrttzbalance", &msg)
			}
		}
		if value.Win > 0 {
			if bigwin == nil {
				bigwin = &GameGold_BigWin{value.Uid, value.Name, value.Head, value.Win}
			} else if value.Win > bigwin.Win {
				bigwin = &GameGold_BigWin{value.Uid, value.Name, value.Head, value.Win}
			}
		}

		//! 插入战绩
		if value.Bets > 0 {
			var record Rec_BrTTZ_Info
			record.Time = time.Now().Unix()
			record.GameType = self.room.Type
			var rec Son_Rec_BrTTZ_Person
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
		}

		if self.Robot.Robots[i].GetWin() != 0 {
			self.Robot.Robots[i].AddMoney(self.Robot.Robots[i].GetWin())
		}

		for j := 0; j < len(self.Seat); j++ {
			if self.Seat[j].Robot == self.Robot.Robots[i] {
				var msg Msg_GameGoldBZW_Balance
				msg.Uid = self.Robot.Robots[i].Id
				msg.Total = self.Robot.Robots[i].GetMoney()
				msg.Win = self.Robot.Robots[i].GetWin()
				self.room.broadCastMsg("gamegoldbrttzbalance", &msg)
				break
			}
		}

		if self.Robot.Robots[i].GetWin() > 0 {
			if bigwin == nil {
				bigwin = &GameGold_BigWin{self.Robot.Robots[i].Id, self.Robot.Robots[i].Name, self.Robot.Robots[i].Head, self.Robot.Robots[i].GetWin()}
			} else if self.Robot.Robots[i].GetWin() > bigwin.Win {
				bigwin = &GameGold_BigWin{self.Robot.Robots[i].Id, self.Robot.Robots[i].Name, self.Robot.Robots[i].Head, self.Robot.Robots[i].GetWin()}
			}
		}
	}

	if self.Dealer == nil && dealwin != 0 { //! 系统庄
		GetServer().SqlBZWLog(&SQL_BZWLog{1, dealwin, time.Now().Unix(), self.room.Type})
	}
	if self.Dealer != nil && robotwin != 0 { //! 玩家庄
		lib.GetRobotMgr().AddRobotWin(self.room.Type, robotwin)
		GetServer().SqlBZWLog(&SQL_BZWLog{1, robotwin, time.Now().Unix(), self.room.Type + 10000000})
	}
	if self.Dealer != nil && lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost == 102 { //! 玩家庄并且是平衡算法
		GetServer().SetBrTTZUsrMoney(self.room.Type%60000, GetServer().BrTTZUsrMoney[self.room.Type%60000]+int64(dealwin))
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
		GetServer().SetBrTTZMoney(self.room.Type%60000, GetServer().BrTTZMoney[self.room.Type%60000]+int64(_dealwin))
		dealwin -= robotwin
		if dealwin > 0 {
			bl := lib.GetManyMgr().GetProperty(self.room.Type).Cost
			cost := int(math.Ceil(float64(dealwin) * bl / 100.0))
			dealwin -= cost
		}
	}

	if self.Dealer != nil {
		var record Rec_BrTTZ_Info
		record.Time = time.Now().Unix()
		record.GameType = self.room.Type
		var rec Son_Rec_BrTTZ_Person
		rec.Uid = self.Dealer.Uid
		rec.Name = self.Dealer.Name
		rec.Head = self.Dealer.Head
		rec.Score = dealwin
		rec.Result = self.Result
		rec.Bets = self.Dealer.BetInfo
		record.Info = append(record.Info, rec)
		GetServer().InsertRecord(self.room.Type, self.Dealer.Uid, lib.HF_JtoA(&record), rec.Score)
	}

	//! 发送庄家结算
	{
		var msg Msg_GameGoldBZW_Balance
		if self.Dealer != nil {
			self.Dealer.Total += dealwin
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
		self.room.broadCastMsg("gamegoldbrttzbalance", &msg)
	}

	//! 30秒的下注时间
	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 10
	self.SetTime(self.BetTime)

	//! 总结算
	{
		var msg Msg_GameGoldBrTTZ_End
		msg.Result = self.Result
		if bigwin != nil {
			msg.Uid = bigwin.Uid
			msg.Name = bigwin.Name
			msg.Head = bigwin.Head
		}
		self.GameCard()
		lib.HF_DeepCopy(&msg.Next, &self.Result)
		if self.Dealer != nil && self.DealType == 1 {
			for i := 0; i < len(msg.Next); i++ {
				msg.Next[i][0] = 0
				msg.Next[i][1] = 0
			}
		} else {
			for i := 0; i < len(msg.Next); i++ {
				msg.Next[i][1] = 0
			}
		}
		msg.Money = lib.GetManyMoneyMgr().GetProperty(self.room.Type).Money
		msg.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime
		self.room.broadCastMsg("gamegoldbrttzend", &msg)
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
				self.room.broadCastMsg("gamebrttzseat", &msg)
				break
			}
		}
		delete(self.PersonMgr, key)
	}

	//! 返回机器人
	self.Robot.Init(3, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)

	for i := 0; i < len(self.Bets); i++ {
		self.Bets[i] = make(map[*Game_GoldBrTTZ_Person]int)
	}

	//! 判断庄家是否能继续连庄
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
		if self.RobotDealer.GetMoney() < lib.GetManyMgr().GetProperty(self.room.Type).UpDealMoney || !lib.GetRobotMgr().GetRobotSet(self.room.Type).NeedRobot {
			self.ChageDeal()
		} else {
			if self.Round >= lib.HF_GetRandom(6)+3 && len(self.LstDeal) > 0 {
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
			self.room.broadCastMsg("gamebrttzseat", &msg)
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
			self.room.broadCastMsg("gamebrttzseat", &msg)
			self.Seat[i].Robot = nil
		}
	}
}

func (self *Game_GoldBrTTZ) OnBye() {
}

func (self *Game_GoldBrTTZ) OnExit(uid int64) {
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
func (self *Game_GoldBrTTZ) GetMoneyPos(index int, robot bool) int {
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
func (self *Game_GoldBrTTZ) GetMoneyPosByRobot(index int) int {
	total := 0
	for _, value := range self.Robot.RobotsBet[index] {
		total += value
	}
	return total
}

//! 得到庄家最多可能赔多少
func (self *Game_GoldBrTTZ) GetMaxLost(robot bool) int {
	total := 0
	for i := 0; i < 3; i++ {
		bl := 20
		if self.Result[i+1][0] != 37 {
			bl = 15
		}
		total += self.GetMoneyPos(i, robot) * bl
	}
	return total
}

func (self *Game_GoldBrTTZ) getInfo(uid int64, total int) *Msg_GameGoldBrTTZ_Info {
	var msg Msg_GameGoldBrTTZ_Info
	msg.Begin = self.room.Begin
	msg.Time = self.Time - time.Now().Unix()
	msg.Total = total
	msg.Trend = self.Trend
	msg.IsDeal = false
	msg.Money = lib.GetManyMoneyMgr().GetProperty(self.room.Type).Money
	msg.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime
	lib.HF_DeepCopy(&msg.Result, &self.Result)
	if self.Dealer != nil && self.DealType == 1 { //!  暗牌
		for i := 0; i < len(msg.Result); i++ {
			msg.Result[i][0] = 0
			msg.Result[i][1] = 0
		}
	} else {
		for i := 0; i < len(msg.Result); i++ {
			msg.Result[i][1] = 0
		}
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
	for i := 0; i < len(self.Seat); i++ {
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

func (self *Game_GoldBrTTZ) GetPerson(uid int64) *Game_GoldBrTTZ_Person {
	return self.PersonMgr[uid]
}

func (self *Game_GoldBrTTZ) OnTime() {
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
		self.RobotSeat(lib.HF_GetRandom(len(self.Seat)), self.Robot.Robots[i])
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

			index, gold, bets := self.Robot.GameBets(self.Robot.Robots[i])
			if gold == 0 {
				continue
			}
			bl := 20
			if self.Result[0][0] != 37 {
				bl = 15
			}
			if bets*bl > self.Robot.Robots[i].GetMoney() {
				self.Robot.GameBackBets(self.Robot.Robots[i], index, gold)
				continue
			}
			//! 判断是否够赔
			if self.Dealer != nil { //! 玩家庄判断是否能下
				dealwin := self.GetMaxLost(true)
				if self.Dealer.Total+self.Total+self.Robot.RobotTotal-dealwin < self.Dealer.Total/5 {
					self.Robot.GameBackBets(self.Robot.Robots[i], index, gold)
					continue
				}
			} else if self.RobotDealer != nil {
				dealwin := self.GetMaxLost(true)
				if self.RobotDealer.GetMoney()+self.Total+self.Robot.RobotTotal-dealwin < self.RobotDealer.GetMoney()/5 {
					self.Robot.GameBackBets(self.Robot.Robots[i], index, gold)
					continue
				}
			} else if lib.GetManyMgr().GetProperty(self.room.Type).DealChange == 1 {
				dealwin := self.GetMaxLost(true)
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
			self.room.broadCastMsg("gamebrttzbets", &msg)
		}
		return
	}

	if !self.room.Begin {
		self.OnBegin()
		return
	}
}

func (self *Game_GoldBrTTZ) OnIsDealer(uid int64) bool {
	if self.Dealer != nil && self.Dealer == self.GetPerson(uid) {
		return true
	}
	return false
}

//! 申请无座玩家
func (self *Game_GoldBrTTZ) GamePlayerList(uid int64) {
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
func (self *Game_GoldBrTTZ) SendTotal(uid int64, total int) {
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
func (self *Game_GoldBrTTZ) SetTime(t int) {
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
func (self *Game_GoldBrTTZ) ChageDeal() {
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
						self.room.broadCastMsg("gamebrttzseat", &msg)
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
						self.room.broadCastMsg("gamebrttzseat", &msg)
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
func (self *Game_GoldBrTTZ) OnIsBets(uid int64) bool {
	value, ok := self.PersonMgr[uid]
	if ok {
		return value.Bets > 0
	}
	return false
}

//! 结算所有人
func (self *Game_GoldBrTTZ) OnBalance() {
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
