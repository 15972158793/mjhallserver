package gameserver

import (
	//"log"
	////"sort"
	"lib"
	"staticfunc"
	//"strings"
	"fmt"
	"math"
	"time"
)

var GOLDMPH_BS []int = []int{2, 2, 12, 8, 8, 12, 4, 24, 24, 4, 100, 50}

var MPH_LV int = 1      //! lv
var MPH_DIOR int = 2    //! 迪奥
var MPH_CHANEL int = 3  //! 香奈儿
var MPH_ROLEX int = 4   //! 劳力士
var MPH_BENZ int = 5    //! 奔驰
var MPH_BENZ100 int = 6 //! 奔驰*100
var MPH_TP int = 7      //! 通赔
var MPH_TS int = 8      //! 通杀
var MPH_GUCCI int = 9   //! 古驰
var MPH_HERMES int = 10 //! 爱马仕
var MPH_ARMANI int = 11 //! 阿玛尼
var MPH_YSL int = 12    //! 圣罗兰

var GOLDMPH_DS []int = []int{MPH_BENZ100,
	MPH_YSL,
	MPH_YSL,
	MPH_YSL,
	MPH_HERMES,
	MPH_HERMES,
	MPH_HERMES,
	MPH_TS,
	MPH_GUCCI,
	MPH_GUCCI,
	MPH_GUCCI,
	MPH_ARMANI,
	MPH_ARMANI,
	MPH_ARMANI,
	MPH_BENZ,
	MPH_ROLEX,
	MPH_ROLEX,
	MPH_ROLEX,
	MPH_DIOR,
	MPH_DIOR,
	MPH_DIOR,
	MPH_TP,
	MPH_LV,
	MPH_LV,
	MPH_LV,
	MPH_CHANEL,
	MPH_CHANEL,
	MPH_CHANEL}

///////////////////////////////////////////////////////
type Game_GoldMPH_Person struct {
	Uid       int64   `json:"uid"`
	Gold      int     `json:"gold"`  //! 进来时候的钱
	Total     int     `json:"total"` //! 当前的钱
	CurBets   int     `json:"curbets"`
	Win       int     `json:"win"`       //! 本局赢的钱
	Cost      int     `json:"cost"`      //! 手续费
	Bets      int     `json:"bets"`      //! 本局下了多少钱
	BetInfo   [11]int `json:"bets"`      //! 本局下的注
	BeBets    int     `json:"bebets"`    //! 上把下了多少钱
	BeBetInfo [11]int `json:"bebetinfo"` //! 上把的下注
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
func (self *Game_GoldMPH_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

type Game_GoldMPHSeat struct {
	Person *Game_GoldMPH_Person
	Robot  *lib.Robot
}

func (self *Game_GoldMPHSeat) GetTotal() int {
	if self.Person != nil {
		return self.Person.Total
	} else if self.Robot != nil {
		return self.Robot.GetMoney()
	}
	return 0
}

type Game_GoldMPH struct {
	PersonMgr   map[int64]*Game_GoldMPH_Person   `json:"personmgr"`
	Bets        [11]map[*Game_GoldMPH_Person]int `json:"bets"`
	Result      int                              `json:"result"`
	Dealer      *Game_GoldMPH_Person             `json:"dealer"`      //! 庄家
	RobotDealer *lib.Robot                       `json:"robotdealer"` //! 机器人庄
	Round       int                              `json:"round"`       //! 连庄轮数
	DownUid     int64                            `json:"downuid"`     //! 下庄的人
	Time        int64                            `json:"time"`
	LstDeal     []Game_GoldMPHSeat               `json:"lstdeal"` //! 上庄列表
	Seat        [8]*Game_GoldMPH_Person          `json:"seat"`    //! 8个位置
	Total       int                              `json:"total"`   //! 这局一共下了多少钱
	Money       int                              `json:"money"`   //! 系统庄的钱
	Trend       []int                            `json:"trend"`   //! 走势
	BetTime     int                              `json:"bettime"`
	Robot       lib.ManyGameRobot                //! 机器人结构

	room *Room
}

func NewGame_GoldMPH() *Game_GoldMPH {
	game := new(Game_GoldMPH)
	game.PersonMgr = make(map[int64]*Game_GoldMPH_Person)
	for i := 0; i < len(game.Bets); i++ {
		game.Bets[i] = make(map[*Game_GoldMPH_Person]int)
	}
	for i := 0; i < 20; i++ {
		game.Trend = append(game.Trend, lib.HF_GetRandom(len(GOLDMPH_DS)))
	}

	return game
}

func (self *Game_GoldMPH) OnInit(room *Room) {
	self.room = room
	self.Money = lib.GetManyMgr().GetProperty(self.room.Type).DealInitMoney
	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 16
	//! 载入机器人
	self.Robot.Init(11, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)
}

func (self *Game_GoldMPH) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldMPH) OnSendInfo(person *Person) {
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
		value.CurBets = 0
		value.SynchroGold(person.Gold)
		person.SendMsg("gamegoldsxdbinfo", self.getInfo(person.Uid, value.Total))
		return
	}

	_person := new(Game_GoldMPH_Person)
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
	person.SendMsg("gamegoldsxdbinfo", self.getInfo(person.Uid, person.Gold))
}

func (self *Game_GoldMPH) OnMsg(msg *RoomMsg) {
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

func (self *Game_GoldMPH) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.Begin = true
	lst := make([]int, 0)
	if self.Dealer != nil { //! 玩家庄
		if self.Robot.RobotTotal == 0 { //! 没有机器人下注
			if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost == 102 { //! 平衡模式
				winlst := make([]int, 0)
				lostlst := make([]int, 0)
				for i := 0; i < len(GOLDMPH_DS); i++ {
					lost := self.GetMoneyByPos(i)
					win := GetServer().MphUserMoney[self.room.Type%200000] + int64(self.Total-lost)
					if win >= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMin && win <= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMax {
						lst = append(lst, i)
					} else {
						if self.Total > lost {
							winlst = append(winlst, i)
						} else {
							lostlst = append(lostlst, i)
						}
					}
				}
				if len(lst) == 0 {
					if GetServer().MphUserMoney[self.room.Type%200000] >= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMax {
						lst = append(lst, lostlst...)
					} else if GetServer().MphUserMoney[self.room.Type%200000] <= lib.GetManyMgr().GetProperty(self.room.Type).PlayerMin {
						lst = append(lst, winlst...)
					}
				}
			} else if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost == 101 { //! 随机模式

			} else if lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost < 100 {
				if lib.HF_GetRandom(100) < lib.GetManyMgr().GetProperty(self.room.Type).PlayerCost { //! 玩家庄胜利
					for i := 0; i < len(GOLDMPH_DS); i++ {
						if self.Total-self.GetMoneyByPos(i) >= 0 {
							lst = append(lst, i)
						}
					}
				} else { //! 玩家庄失败
					for i := 0; i < len(GOLDMPH_DS); i++ {
						if self.Total-self.GetMoneyByPos(i) <= 0 {
							lst = append(lst, i)
						}
					}
				}
			}
		} else {
			for i := 0; i < len(GOLDMPH_DS); i++ {
				win := self.GetRobotByPos(i)
				win -= self.Robot.RobotTotal
				if lib.GetRobotMgr().GetRobotWin(self.room.Type)+win >= 0 || win >= 0 {
					lst = append(lst, i)
				}
			}
		}
	} else {
		for i := 0; i < len(GOLDMPH_DS); i++ {
			lost := self.GetMoneyByPos(i)
			win := GetServer().MphSysMoney[self.room.Type%200000] + int64(self.Total-lost)
			if win >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
				lst = append(lst, i)
			}
		}
		if len(lst) == 0 {
			for i := 0; i < len(GOLDMPH_DS); i++ {
				lost := self.GetMoneyByPos(i)
				if lost <= self.Total {
					lst = append(lst, i)
				}
			}
		}
	}

	if len(lst) == 0 {
		self.Result = lib.HF_GetRandom(len(GOLDMPH_DS))
		lib.GetLogMgr().Output(lib.LOG_ERROR, "名品汇随机了")
	} else {
		self.Result = lst[lib.HF_GetRandom(len(lst))]
	}
	//	self.Result = 14

	self.OnEnd()
}

//! 这个位置能赢多少钱
func (self *Game_GoldMPH) GetMoneyByPos(pos int) int {
	if pos == 0 { //! 发100
		return self.GetMoneyPos(10, false) * GOLDMPH_BS[10]
	} else if pos == 7 { //! 通杀
		return 0
	} else if pos == 14 { //! 发
		return self.GetMoneyPos(10, false) * GOLDMPH_BS[11]
	} else if pos == 21 { //! 通赔
		all := self.Total
		if self.Dealer != nil {
			all += self.Robot.RobotTotal
		}
		return 2 * all
	} else if pos == 1 || pos == 2 || pos == 3 { //! ysl
		return self.GetMoneyPos(1, false)*GOLDMPH_BS[1] + self.GetMoneyPos(9, false)*GOLDMPH_BS[9]
	} else if pos == 4 || pos == 5 || pos == 6 { //! 爱马仕
		return self.GetMoneyPos(1, false)*GOLDMPH_BS[1] + self.GetMoneyPos(5, false)*GOLDMPH_BS[5]
	} else if pos == 8 || pos == 9 || pos == 10 { //! gucci
		return self.GetMoneyPos(1, false)*GOLDMPH_BS[1] + self.GetMoneyPos(4, false)*GOLDMPH_BS[4]
	} else if pos == 11 || pos == 12 || pos == 13 { //! 阿玛尼
		return self.GetMoneyPos(1, false)*GOLDMPH_BS[1] + self.GetMoneyPos(8, false)*GOLDMPH_BS[8]
	} else if pos == 15 || pos == 16 || pos == 17 { //! rolex
		return self.GetMoneyPos(0, false)*GOLDMPH_BS[0] + self.GetMoneyPos(7, false)*GOLDMPH_BS[7]
	} else if pos == 18 || pos == 19 || pos == 20 { //! dior
		return self.GetMoneyPos(0, false)*GOLDMPH_BS[0] + self.GetMoneyPos(3, false)*GOLDMPH_BS[3]
	} else if pos == 22 || pos == 23 || pos == 24 { //! lv
		return self.GetMoneyPos(0, false)*GOLDMPH_BS[0] + self.GetMoneyPos(2, false)*GOLDMPH_BS[2]
	} else if pos == 25 || pos == 26 || pos == 27 { //! chanel
		return self.GetMoneyPos(0, false)*GOLDMPH_BS[0] + self.GetMoneyPos(6, false)*GOLDMPH_BS[6]
	}

	return 0
}

//! 这个位置能赢多少钱
func (self *Game_GoldMPH) GetRobotByPos(pos int) int {
	if pos == 0 { //! 发100
		return self.GetMoneyPosByRobot(10) * GOLDMPH_BS[10]
	} else if pos == 7 { //! 通杀
		return 0
	} else if pos == 14 { //! 发
		return self.GetMoneyPosByRobot(10) * GOLDMPH_BS[11]
	} else if pos == 21 { //! 通赔
		return 2 * self.Robot.RobotTotal
	} else if pos == 1 || pos == 2 || pos == 3 { //! ysl
		return self.GetMoneyPosByRobot(1)*GOLDMPH_BS[1] + self.GetMoneyPosByRobot(9)*GOLDMPH_BS[9]
	} else if pos == 4 || pos == 5 || pos == 6 { //! 爱马仕
		return self.GetMoneyPosByRobot(1)*GOLDMPH_BS[1] + self.GetMoneyPosByRobot(5)*GOLDMPH_BS[5]
	} else if pos == 8 || pos == 9 || pos == 10 { //! gucci
		return self.GetMoneyPosByRobot(1)*GOLDMPH_BS[1] + self.GetMoneyPosByRobot(4)*GOLDMPH_BS[4]
	} else if pos == 11 || pos == 12 || pos == 13 { //! 阿玛尼
		return self.GetMoneyPosByRobot(1)*GOLDMPH_BS[1] + self.GetMoneyPosByRobot(8)*GOLDMPH_BS[8]
	} else if pos == 15 || pos == 16 || pos == 17 { //! rolex
		return self.GetMoneyPosByRobot(0)*GOLDMPH_BS[0] + self.GetMoneyPosByRobot(7)*GOLDMPH_BS[7]
	} else if pos == 18 || pos == 19 || pos == 20 { //! dior
		return self.GetMoneyPosByRobot(0)*GOLDMPH_BS[0] + self.GetMoneyPosByRobot(3)*GOLDMPH_BS[3]
	} else if pos == 22 || pos == 23 || pos == 24 { //! lv
		return self.GetMoneyPosByRobot(0)*GOLDMPH_BS[0] + self.GetMoneyPosByRobot(2)*GOLDMPH_BS[2]
	} else if pos == 25 || pos == 26 || pos == 27 { //! chanel
		return self.GetMoneyPosByRobot(0)*GOLDMPH_BS[0] + self.GetMoneyPosByRobot(6)*GOLDMPH_BS[6]
	}

	return 0
}

//! 得到压中了这个点可以赢的钱
func (self *Game_GoldMPH) GetDealWinByPos(pos int, robot bool) int {
	//! 先得到通赔的钱
	tp := self.Total

	lost := self.GetMoneyPos(pos, robot) * GOLDMPH_BS[pos] //! 压中这个点本身赢的钱
	maxlost := 0
	if pos == 0 {
		for i := 2; i <= 3; i++ {
			tmp := self.GetMoneyPos(i, robot) * GOLDMPH_BS[i]
			if tmp > maxlost {
				maxlost = tmp
			}
		}
		for i := 6; i <= 7; i++ {
			tmp := self.GetMoneyPos(i, robot) * GOLDMPH_BS[i]
			if tmp > maxlost {
				maxlost = tmp
			}
		}
	} else if pos == 1 {
		for i := 4; i <= 5; i++ {
			tmp := self.GetMoneyPos(i, robot) * GOLDMPH_BS[i]
			if tmp > maxlost {
				maxlost = tmp
			}
		}
		for i := 8; i <= 9; i++ {
			tmp := self.GetMoneyPos(i, robot) * GOLDMPH_BS[i]
			if tmp > maxlost {
				maxlost = tmp
			}
		}
	} else if pos == 2 || pos == 3 || pos == 6 || pos == 7 {
		maxlost = self.GetMoneyPos(0, robot) * GOLDMPH_BS[0]
	} else if pos == 4 || pos == 5 || pos == 8 || pos == 9 {
		maxlost = self.GetMoneyPos(1, robot) * GOLDMPH_BS[1]
	}

	if lost+maxlost > tp {
		tp = lost + maxlost
	}

	return tp
}

func (self *Game_GoldMPH) GameUpDeal(uid int64) {
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
		self.LstDeal = append(self.LstDeal, Game_GoldMPHSeat{person, nil})
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
	self.room.SendMsg(uid, "gamesxdbdeal", &msg)
}

//! 机器人上庄
func (self *Game_GoldMPH) RobotUpDeal(robot *lib.Robot) {
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
	self.LstDeal = append(self.LstDeal, Game_GoldMPHSeat{nil, robot})
}

func (self *Game_GoldMPH) GameReDeal(uid int64) {
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
	self.room.SendMsg(uid, "gamesxdbdeal", &msg)
}

//! 坐下
func (self *Game_GoldMPH) GameSeat(uid int64, index int) {
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
	self.room.broadCastMsg("gamesxdbseat", &msg)
}

func (self *Game_GoldMPH) GameBets(uid int64, index int, gold int) {
	if uid == 0 {
		return
	}

	if index < 0 || index > 10 {
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

	person.CurBets += gold
	person.Bets += gold
	person.Total -= gold
	person.BetInfo[index] += gold
	person.Round = 0

	var msg Msg_GameGoldBZW_Bets
	msg.Uid = uid
	msg.Index = index
	msg.Gold = gold
	msg.Total = person.Total
	self.room.broadCastMsg("gamesxdbbets", &msg)
}

//! 续压
func (self *Game_GoldMPH) GameGoOn(uid int64) {
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

	person.CurBets += person.BeBets
	person.Bets += person.BeBets
	person.Total -= person.BeBets
	for i := 0; i < len(person.BeBetInfo); i++ {
		person.BetInfo[i] += person.BeBetInfo[i]
	}
	person.Round = 0

	var msg Msg_GameGoldSXDB_Goon
	msg.Uid = uid
	msg.Gold = person.BeBetInfo
	msg.Total = person.Total
	self.room.broadCastMsg("gamesxdbgoon", &msg)
}

//! 结算
func (self *Game_GoldMPH) OnEnd() {
	self.room.Begin = false

	tmp := []int{self.Result}
	tmp = append(tmp, self.Trend...)
	if len(tmp) > 20 {
		tmp = tmp[0:20]
	}
	self.Trend = tmp

	dealwin := 0
	robotwin := 0
	if GOLDMPH_DS[self.Result] == MPH_BENZ100 { //! 结果是发100
		for key, value := range self.Bets[10] {
			winmoney := value * 100
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[10] {
			winmoney := value * 100
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
	} else if GOLDMPH_DS[self.Result] == MPH_BENZ { //! 结果是发
		for key, value := range self.Bets[10] {
			winmoney := value * 50
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[10] {
			winmoney := value * 50
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
	} else if GOLDMPH_DS[self.Result] == MPH_TP { //! 结果是通赔
		for i := 0; i < len(self.Bets); i++ {
			for key, value := range self.Bets[i] {
				winmoney := value * 2
				dealwin -= winmoney
				key.Win += winmoney
				key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
			}
			for key, value := range self.Robot.RobotsBet[i] {
				winmoney := value * 2
				key.AddWin(winmoney)
				key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
				robotwin += winmoney
				if self.Dealer != nil {
					dealwin -= winmoney
				}
			}
		}
	} else if GOLDMPH_DS[self.Result] == MPH_GUCCI {
		for key, value := range self.Bets[1] {
			winmoney := value * GOLDMPH_BS[1]
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[1] {
			winmoney := value * GOLDMPH_BS[1]
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
		for key, value := range self.Bets[4] {
			winmoney := value * GOLDMPH_BS[4]
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[4] {
			winmoney := value * GOLDMPH_BS[4]
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
	} else if GOLDMPH_DS[self.Result] == MPH_HERMES {
		for key, value := range self.Bets[1] {
			winmoney := value * GOLDMPH_BS[1]
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[1] {
			winmoney := value * GOLDMPH_BS[1]
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
		for key, value := range self.Bets[5] {
			winmoney := value * GOLDMPH_BS[5]
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[5] {
			winmoney := value * GOLDMPH_BS[5]
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
	} else if GOLDMPH_DS[self.Result] == MPH_ARMANI {
		for key, value := range self.Bets[1] {
			winmoney := value * GOLDMPH_BS[1]
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[1] {
			winmoney := value * GOLDMPH_BS[1]
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
		for key, value := range self.Bets[8] {
			winmoney := value * GOLDMPH_BS[8]
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[8] {
			winmoney := value * GOLDMPH_BS[8]
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
	} else if GOLDMPH_DS[self.Result] == MPH_YSL {
		for key, value := range self.Bets[1] {
			winmoney := value * GOLDMPH_BS[1]
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[1] {
			winmoney := value * GOLDMPH_BS[1]
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
		for key, value := range self.Bets[9] {
			winmoney := value * GOLDMPH_BS[9]
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[9] {
			winmoney := value * GOLDMPH_BS[9]
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
	} else if GOLDMPH_DS[self.Result] == MPH_LV {
		for key, value := range self.Bets[0] {
			winmoney := value * GOLDMPH_BS[0]
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[0] {
			winmoney := value * GOLDMPH_BS[0]
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
		for key, value := range self.Bets[2] {
			winmoney := value * GOLDMPH_BS[2]
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[2] {
			winmoney := value * GOLDMPH_BS[2]
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
	} else if GOLDMPH_DS[self.Result] == MPH_DIOR {
		for key, value := range self.Bets[0] {
			winmoney := value * GOLDMPH_BS[0]
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[0] {
			winmoney := value * GOLDMPH_BS[0]
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
		for key, value := range self.Bets[3] {
			winmoney := value * GOLDMPH_BS[3]
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[3] {
			winmoney := value * GOLDMPH_BS[3]
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
	} else if GOLDMPH_DS[self.Result] == MPH_CHANEL {
		for key, value := range self.Bets[0] {
			winmoney := value * GOLDMPH_BS[0]
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[0] {
			winmoney := value * GOLDMPH_BS[0]
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
		for key, value := range self.Bets[6] {
			winmoney := value * GOLDMPH_BS[6]
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[6] {
			winmoney := value * GOLDMPH_BS[6]
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
	} else if GOLDMPH_DS[self.Result] == MPH_ROLEX {
		for key, value := range self.Bets[0] {
			winmoney := value * GOLDMPH_BS[0]
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[0] {
			winmoney := value * GOLDMPH_BS[0]
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
		for key, value := range self.Bets[7] {
			winmoney := value * GOLDMPH_BS[7]
			dealwin -= winmoney
			key.Win += winmoney
			key.Cost += int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}
		for key, value := range self.Robot.RobotsBet[7] {
			winmoney := value * GOLDMPH_BS[7]
			key.AddWin(winmoney)
			key.AddCost(int(math.Ceil(float64(winmoney-value) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0)))
			robotwin += winmoney
			if self.Dealer != nil {
				dealwin -= winmoney
			}
		}
	}

	robotwin -= self.Robot.RobotTotal
	if self.Dealer != nil {
		dealwin += self.Robot.RobotTotal
	}
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
				if self.Seat[j] == value {
					self.room.broadCastMsg("gamegoldsxdbbalance", &msg)
					find = true
					break
				}
			}
			if !find {
				self.room.SendMsg(value.Uid, "gamegoldsxdbbalance", &msg)
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
			var record Rec_SXDB_Info
			record.Time = time.Now().Unix()
			record.GameType = self.room.Type
			var rec Son_Rec_SXDB_Person
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
		GetServer().SetMphUserMoney(self.room.Type%200000, GetServer().MphUserMoney[self.room.Type%200000]+int64(dealwin))
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
		GetServer().SetMphSysMoney(self.room.Type%200000, GetServer().MphSysMoney[self.room.Type%200000]+int64(_dealwin))
		dealwin -= robotwin
		if dealwin > 0 {
			bl := lib.GetManyMgr().GetProperty(self.room.Type).Cost
			cost := int(math.Ceil(float64(dealwin) * bl / 100.0))
			dealwin -= cost
		}
	}

	if self.Dealer != nil {
		var record Rec_SXDB_Info
		record.Time = time.Now().Unix()
		record.GameType = self.room.Type
		var rec Son_Rec_SXDB_Person
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
		self.room.broadCastMsg("gamegoldsxdbbalance", &msg)
	}

	//! 30秒的下注时间
	self.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime + 16
	self.SetTime(self.BetTime)

	//! 总结算
	{
		var msg Msg_GameGoldSXDB_End
		msg.Money = lib.GetManyMoneyMgr().GetProperty(self.room.Type).Money
		msg.BetTime = lib.GetManyMgr().GetProperty(self.room.Type).BetTime
		msg.Result = self.Result
		if bigwin != nil {
			msg.Uid = bigwin.Uid
			msg.Name = bigwin.Name
			msg.Head = bigwin.Head
		}
		self.room.broadCastMsg("gamegoldsxdbend", &msg)
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
			if self.Seat[j] == value {
				self.Seat[j] = nil
				var msg Msg_GameGoldBZW_UpdSeat
				msg.Index = j
				self.room.broadCastMsg("gamesxdbseat", &msg)
				break
			}
		}
		delete(self.PersonMgr, key)
	}

	//! 返回机器人
	self.Robot.Init(11, lib.GetManyMoneyMgr().GetProperty(self.room.Type).RobotMoney)
	self.Robot.Refresh(self.room.Type)

	for i := 0; i < len(self.Bets); i++ {
		self.Bets[i] = make(map[*Game_GoldMPH_Person]int)
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
		if self.Seat[i] == nil {
			continue
		}
		if self.Seat[i].Total < lib.GetManyMgr().GetProperty(self.room.Type).UpSeatMoney {
			self.Seat[i].Seat = -1
			var msg Msg_GameGoldBZW_UpdSeat
			msg.Index = i
			self.room.broadCastMsg("gamesxdbseat", &msg)
			self.Seat[i] = nil
		}
	}
}

func (self *Game_GoldMPH) OnBye() {
}

func (self *Game_GoldMPH) OnExit(uid int64) {
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
func (self *Game_GoldMPH) GetMoneyPos(index int, robot bool) int {
	if index > 10 {
		index = 10
	}

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
func (self *Game_GoldMPH) GetMoneyPosByRobot(index int) int {
	total := 0
	for _, value := range self.Robot.RobotsBet[index] {
		total += value
	}
	return total
}

func (self *Game_GoldMPH) getInfo(uid int64, total int) *Msg_GameGoldSXDB_Info {
	var msg Msg_GameGoldSXDB_Info
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

func (self *Game_GoldMPH) GetPerson(uid int64) *Game_GoldMPH_Person {
	return self.PersonMgr[uid]
}

func (self *Game_GoldMPH) OnTime() {
	if lib.GetRobotMgr().GetRobotSet(self.room.Type).Dealer && len(self.Robot.Robots) > 0 && len(self.LstDeal) < 5 { //! 需要机器人上庄
		self.RobotUpDeal(self.Robot.Robots[lib.HF_GetRandom(len(self.Robot.Robots))])
	}

	if self.Time == 0 {
		return
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
			var msg Msg_GameGoldBZW_Bets
			msg.Uid = self.Robot.Robots[i].Id
			msg.Index = index
			msg.Gold = gold
			msg.Total = self.Robot.Robots[i].GetMoney()
			self.room.broadCastMsg("gamesxdbbets", &msg)
		}
		return
	}

	if !self.room.Begin {
		self.OnBegin()
		return
	}
}

func (self *Game_GoldMPH) OnIsDealer(uid int64) bool {
	if self.Dealer != nil && self.Dealer == self.GetPerson(uid) {
		return true
	}
	return false
}

//! 申请无座玩家
func (self *Game_GoldMPH) GamePlayerList(uid int64) {
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
func (self *Game_GoldMPH) SendTotal(uid int64, total int) {
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
func (self *Game_GoldMPH) SetTime(t int) {
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
func (self *Game_GoldMPH) ChageDeal() {
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
			} else if self.LstDeal[0].Robot != nil {
				self.RobotDealer = self.LstDeal[0].Robot
				self.RobotDealer.SetSeat(100)
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
func (self *Game_GoldMPH) OnIsBets(uid int64) bool {
	value, ok := self.PersonMgr[uid]
	if ok {
		return value.Bets > 0
	}
	return false
}

//! 结算所有人
func (self *Game_GoldMPH) OnBalance() {
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
