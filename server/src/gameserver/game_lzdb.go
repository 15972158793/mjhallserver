package gameserver

import (
	"lib"
	"math"
	"staticfunc"
	"time"
)

var LZDB_ICON []int = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11} //! 1-4 小赔率图标 5-9大赔率图标 10-万金油 11-特殊图标
//!1－9图标的赔率
var LZDB_PL [9][3]int = [9][3]int{{2, 15, 80}, {5, 20, 100}, {8, 25, 120}, {10, 30, 125}, {15, 60, 150}, {15, 75, 175}, {20, 80, 200}, {20, 100, 250}, {30, 125, 500}}
var LZDB_RATE int = 1     //! 1硬币等于多少金币
var LZDB_LZ int = 15      //! 这一行出万金油的概率
var LZDB_LZ_GRID int = 90 //!　这一格出万金油的概率
var LZDB_SPECIAL int = 3  //!　特殊玩法2出现概率

type Rec_LZDB_Info struct {
	GameType int                   `json:"gametype"`
	Time     int64                 `json:"time"` //! 记录时间
	Info     []Son_Rec_LZDB_Person `json:"info"`
}
type Son_Rec_LZDB_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Score int    `json:"score"`
	Bets  int    `json:"bets"`
}

type Game_LZDB struct {
	Person   *Game_LZDB_Person `json:"person"`
	Money    int               `json:"money"` //!　庄家钱
	Result   [5][3]int         `json:"result"`
	Special1 int               `json:"special1"` //!　特殊玩法1,可以免费转多少轮
	Special2 int               `json:"special2"` //!　特殊玩法2,消哪一行
	LZ       [3]int            `json:"lz"`       //! 3-5行是否开万金油
	Time     int64             `json:"time"`

	//!-------------测试
	C bool

	room *Room
}

func NewGame_LZDB() *Game_LZDB {
	game := new(Game_LZDB)
	game.Money = 1000000000
	game.Special1 = -1
	game.Special2 = -1
	game.C = true
	for i := 0; i < len(game.Result); i++ {
		for j := 0; j < len(game.Result[i]); j++ {
			if i == 1 {
				for true {
					tmp := lib.HF_GetRandom(9) + 1
					for z := 0; z < 3; z++ {
						if tmp == game.Result[0][z] {
							tmp = -1
							break
						}
					}
					if tmp == -1 {
						continue
					}
					game.Result[i][j] = tmp
					break
				}
			} else {
				game.Result[i][j] = lib.HF_GetRandom(9) + 1
			}
		}
	}
	return game
}

type Game_LZDB_Person struct {
	Uid     int64  `json:"uid"`
	Gold    int    `json:"gold"`  //! 进房时的金币
	Total   int    `json:"total"` //! 金币数
	Win     int    `json:"win"`   //! 赢了多少钱
	Coin    int    `json:"coin"`  //! 下注硬币数量
	Cost    int    `json:"cost"`  //!　抽水
	Name    string `json:"name"`  //! 名字
	Head    string `json:"head"`  //! 头像
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type Msg_GameLZDB_Info struct {
	Begin  bool              `json:"begin"`
	Money  int               `json:"money"` //!　庄家钱
	Result [5][3]int         `json:"result"`
	Person Son_GameLZDB_Info `json:"person"`
}

type Son_GameLZDB_Info struct {
	Uid     int64  `json:"uid"`
	Total   int    `json:"total"` //! 金币数
	Coin    int    `json:"coin"`  //! 下注硬币数量
	Name    string `json:"name"`  //! 名字
	Head    string `json:"head"`  //! 头像
	IP      string `json:"ip"`
	Address string `json:"address"`
	Sex     int    `json:"sex"`
}

type Msg_GameLZDB_Total struct {
	Uid   int64 `json:"uid"`
	Total int   `json:"total"`
}

type Msg_GameLZDB_End struct {
	Special1  int                `json:"special1"` //!　特殊玩法1,可以免费转多少轮
	Special2  int                `json:"special2"` //!　特殊玩法2,消哪一行
	Info      []Son_GameLZDB_End `json:"info"`
	Uid       int64              `json:"uid"`
	Total     int                `json:"total"`     //! 金币数
	Playerwin int                `json:"playerwin"` //! 赢了多少钱
	Name      string             `json:"name"`      //! 名字
	Head      string             `json:"head"`      //! 头像
	IP        string             `json:"ip"`
	Address   string             `json:"address"`
	Sex       int                `json:"sex"`
}

type Son_GameLZDB_End struct {
	Result    [5][3]int `json:"result"`
	PersonWin int       `json:"personwin"` //!这一轮赢了多少
	DelIcon   []int     `json:"delicon"`
}

func (self *Game_LZDB) getinfo(uid int64) *Msg_GameLZDB_Info {
	var msg Msg_GameLZDB_Info
	msg.Begin = self.room.Begin
	msg.Money = self.Money
	msg.Result = self.Result
	if self.Person != nil && uid == self.Person.Uid {
		msg.Person.Uid = self.Person.Uid
		msg.Person.Total = self.Person.Total
		msg.Person.Coin = self.Person.Coin
		msg.Person.Address = self.Person.Address
		msg.Person.Head = self.Person.Head
		msg.Person.IP = self.Person.IP
		msg.Person.Name = self.Person.Name
		msg.Person.Sex = self.Person.Sex
	}
	return &msg
}

//! 同步金币
func (self *Game_LZDB_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

//! 同步总分
func (self *Game_LZDB) SendTotal(uid int64, total int) {
	var msg Msg_GameLZDB_Total
	msg.Uid = uid
	msg.Total = total
	self.room.SendMsg(uid, "gametotal", &msg)
}

//!　这样开能赢多少 return 这局赢的分数 ， 特殊图标的个数
func (self *Game_LZDB) GetResultWin(result [5][3]int) (int, int) {
	playerwin := make([]int, 3)
	icon := make([]int, 3)

	for i := 0; i < 3; i++ {
		length := 1 //! 长度
		for j := 1; j < 5; j++ {
			if result[0][i] == 11 {
				break
			}
			find := false
			for z := 0; z < 3; z++ {
				if result[j][z] == result[0][i] || result[j][z] == 10 {
					length++
					find = true
					break
				}
			}
			if !find {
				break
			}
		}

		if length >= 3 { //! 可以获得分数
			num := 1 //! 条数
			icon[i] = result[0][i]
			for j := 1; j < length; j++ {
				tmp := 0
				for z := 0; z < 3; z++ {
					if result[j][z] == result[0][i] || result[j][z] == 10 {
						tmp++
					}
				}
				num *= tmp
			}

			playerwin[i] = (num * LZDB_PL[result[0][i]-1][length-3] * self.Person.Coin)
		} else {
			icon[i] = 0
		}
	}

	//! 得分是否翻倍
	small := 0
	big := 0

	for i := 0; i < len(icon); i++ {
		find := false
		for j := 0; j < i; j++ {
			if icon[i] == icon[j] {
				find = true
				break
			}
		}
		if find {
			continue
		}
		if icon[i] >= 1 && icon[i] <= 4 {
			small++
		} else if icon[i] >= 5 && icon[i] <= 9 {
			big++
		}
	}

	if small >= 2 {
		for i := 0; i < 3; i++ {
			if icon[i] >= 1 && icon[i] <= 4 {
				playerwin[i] *= small
			}
		}
	} else if big >= 2 {
		for i := 0; i < 3; i++ {
			if icon[i] >= 5 && icon[i] <= 9 {
				playerwin[i] *= big
			}
		}
	}

	//!　判断是否又特殊得分
	specialScore := 0
	specialNum := 0
	for i := 0; i < 5; i++ {
		for j := 0; j < 3; j++ {
			if result[i][j] == 11 {
				specialNum++
			}
		}
	}
	if specialNum == 3 {
		specialScore = 2500
		if GetServer().Con.MoneyMode == 2 {
			specialScore = 250000
		}
	} else if specialNum == 4 {
		specialScore = 25000
		if GetServer().Con.MoneyMode == 2 {
			specialScore = 2500000
		}
	} else if specialNum == 5 {
		specialScore = 125000
		if GetServer().Con.MoneyMode == 2 {
			specialScore = 12500000
		}
	}

	return playerwin[0] + playerwin[1] + playerwin[2] + specialScore, specialNum
}

//!　消掉得分图标 return 是否下落
func (self *Game_LZDB) DelIcon(del bool) (bool, []int) {
	length := 1

	indexIcon := make([]int, 0)
	for i := 0; i < 3; i++ {
		if self.Result[0][i] == 0 {
			continue
		}
		tmp := 1
		for j := 1; j < 5; j++ {
			find := false
			for z := 0; z < 3; z++ {
				if self.Result[0][i] == self.Result[j][z] || self.Result[j][z] == 10 {
					tmp++
					find = true
					break
				}
			}
			if !find {
				break
			}
		}

		if tmp >= 3 {
			if tmp > length {
				length = tmp
			}
			for j := 1; j < tmp; j++ {
				for z := 0; z < 3; z++ {
					if self.Result[j][z] == self.Result[0][i] {
						if del {
							self.Result[j][z] = 0
						} else {
							indexIcon = append(indexIcon, j*10+z)
						}

					}
				}
			}
			for j := i + 1; j < 3; j++ {
				if self.Result[0][i] == self.Result[0][j] {
					if del {
						self.Result[0][j] = 0
					} else {
						indexIcon = append(indexIcon, 0*10+j)
					}

				}
			}
			if del {
				self.Result[0][i] = 0
			} else {
				indexIcon = append(indexIcon, 0*10+i)
			}

		}
	}

	for i := 1; i < length; i++ {
		for j := 0; j < 3; j++ {
			if self.Result[i][j] == 10 {
				self.Result[i][j] = 0
			}
		}
	}

	goOn := false

	delIcon := make([]int, 0)
	for i := 0; i < 5; i++ {
		for j := 0; j < 3; j++ {
			if self.Result[i][j] == 0 {
				delIcon = append(delIcon, i*10+j)
			}
		}
	}

	for i := 0; i < 5; i++ { //! 消除后图标下落
		for j := 0; j < 3; j++ {
			if self.Result[i][j] == 0 {
				goOn = true
				for z := j + 1; z <= 2; z++ {
					if self.Result[i][z] != 0 {
						self.Result[i][j] = self.Result[i][z]
						self.Result[i][z] = 0
						break
					}
				}
			}
		}
	}
	if del {
		return goOn, delIcon
	} else {
		return false, indexIcon
	}

}

//! 填充消除后的图标
func (self *Game_LZDB) GetGoOnIcon(playWin int) {
	icon := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i := 2; i < 5; i++ { //!先填充后三行
		lz := false
		if self.LZ[i-2] == 1 { //! 这一行之前是万金油图标
			lz = true
		}
		for j := 0; j < 3; j++ {
			if self.Result[i][j] == 0 {
				if lz {
					if lib.HF_GetRandom(100) < LZDB_LZ_GRID {
						self.Result[i][j] = 10
					} else {

						self.Result[i][j] = icon[lib.HF_GetRandom(len(icon))]

					}
				} else {

					self.Result[i][j] = icon[lib.HF_GetRandom(len(icon))]

				}
			}
		}
	}

	//!　先填充第一行空的图标
	for i := 0; i < 3; i++ {
		if self.Result[0][i] == 0 {
			self.Result[0][i] = icon[lib.HF_GetRandom(len(icon))]
		}
	}

	//! 填充第二行空的图标
	num := 0 //!第二行需要填充的个数
	for i := 0; i < 3; i++ {
		if self.Result[1][i] == 0 {
			num++
		}
	}
	if num > 0 {
		icon = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		for i := 0; i < 3; i++ { //!把第一行出现过得图标去掉
			for j := 0; j < len(icon); {
				if self.Result[0][i] == icon[j] {
					copy(icon[j:], icon[j+1:])
					icon = icon[:len(icon)-1]
				} else {
					j++
				}
			}
		}
		sysMoney := 0 - playWin
		if self.Special1 == -1 {
			sysMoney += self.Person.Coin * LZDB_RATE
		}

		if GetServer().LzdbSysMoney[self.room.Type%160000]+int64(sysMoney) < lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {

			if num == 1 {
				self.Result[1][2] = icon[lib.HF_GetRandom(len(icon))]
			} else if num == 2 {
				self.Result[1][2] = icon[lib.HF_GetRandom(len(icon))]
				self.Result[1][1] = icon[lib.HF_GetRandom(len(icon))]
			} else if num == 3 {
				self.Result[1][2] = icon[lib.HF_GetRandom(len(icon))]
				self.Result[1][1] = icon[lib.HF_GetRandom(len(icon))]
				self.Result[1][0] = icon[lib.HF_GetRandom(len(icon))]
			}
		} else {

			if num == 1 {
				loop := 0
				for true && loop < 10 {
					result := self.Result
					result[1][2] = lib.HF_GetRandom(9) + 1
					playerwin, _ := self.GetResultWin(result)
					if GetServer().LzdbSysMoney[self.room.Type%160000]+int64(sysMoney)+int64(0-playerwin) > lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && GetServer().LzdbSysMoney[self.room.Type%160000]+int64(sysMoney)+int64(0-playerwin) < lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {

						self.Result = result
						break
					}
					loop++
				}
				if loop >= 10 {
					self.Result[1][2] = icon[lib.HF_GetRandom(len(icon))]
				}
			} else if num == 2 {
				loop := 0
				for true && loop < 10 {
					result := self.Result
					result[1][2] = lib.HF_GetRandom(9) + 1
					result[1][1] = lib.HF_GetRandom(9) + 1
					playerwin, _ := self.GetResultWin(result)
					if GetServer().LzdbSysMoney[self.room.Type%160000]+int64(sysMoney)+int64(0-playerwin) > lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && GetServer().LzdbSysMoney[self.room.Type%160000]+int64(sysMoney)+int64(0-playerwin) < lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
						self.Result = result
						break
					}
					loop++
				}
				if loop >= 10 {
					self.Result[1][2] = icon[lib.HF_GetRandom(len(icon))]
					self.Result[1][1] = icon[lib.HF_GetRandom(len(icon))]
				}

			} else if num == 3 {
				loop := 0
				for true && loop < 10 {
					result := self.Result
					result[1][2] = lib.HF_GetRandom(9) + 1
					result[1][1] = lib.HF_GetRandom(9) + 1
					result[1][0] = lib.HF_GetRandom(9) + 1
					playerwin, _ := self.GetResultWin(result)
					if GetServer().LzdbSysMoney[self.room.Type%160000]+int64(sysMoney)+int64(0-playerwin) > lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && GetServer().LzdbSysMoney[self.room.Type%160000]+int64(sysMoney)+int64(0-playerwin) < lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
						self.Result = result
						break
					}
					loop++
				}
				if loop >= 10 {
					self.Result[1][2] = icon[lib.HF_GetRandom(len(icon))]
					self.Result[1][1] = icon[lib.HF_GetRandom(len(icon))]
					self.Result[1][0] = icon[lib.HF_GetRandom(len(icon))]
				}
			}
		}
	}

}

func (self *Game_LZDB) SetCoin(uid int64, bet int) {

	self.Person.Coin = bet / LZDB_RATE
	self.OnBegin()
	self.Time = time.Now().Unix() + 300
}

func (self *Game_LZDB) OnBegin() {
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
	self.room.Begin = true
	if self.Person.Coin == 0 {
		self.Person.Coin = 1
	}

	if self.Person == nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "self.Person == nil")
		return
	}
	if self.Person.Total < self.Person.Coin*LZDB_RATE {
		self.room.SendErr(self.Person.Uid, "您的金币不足，请前往充值")
		return
	}

	self.OnEnd()
}

func (self *Game_LZDB) OnEnd() {
	self.room.Begin = false

	//! ------------------ 设置开奖结果 ------------------
	//! 扣金币
	if self.Special1 == -1 {
		self.Person.Total -= self.Person.Coin * LZDB_RATE
		var bet Msg_GameLZDB_Total
		bet.Uid = self.Person.Uid
		bet.Total = self.Person.Total
		self.room.broadCastMsg("gamelzdbbets", &bet)
	}

	icon := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	lostlst := make([]Msg_GameLZDB_End, 0)
	lostlst99 := make([]Msg_GameLZDB_End, 0)
	winlst := make([]Msg_GameLZDB_End, 0)
	for k := 0; k < 50; k++ {
		goOn := true //! 是否继续下落
		special := 5 //!  特殊图标最大数量
		spNum := 0
		for i := 0; i < 5; i++ { //! 确定1,3,4,5行图标

			if i == 1 { //! 第二行空出来,通过计算决定具体图标
				continue
			}

			lz := false //! 这一行是否开万金油
			if i > 1 {  //! 后三行可以开出万金油图标
				if LZDB_LZ > lib.HF_GetRandom(100) {
					lz = true
					self.LZ[i-2] = 1 //! 这一行开万金油
				}
			}

			for j := 0; j < 3; j++ {
				if lz {
					if LZDB_LZ_GRID > lib.HF_GetRandom(100) {
						self.Result[i][j] = 10
					} else {

						self.Result[i][j] = icon[lib.HF_GetRandom(len(icon))]

					}
				} else {
					if special > lib.HF_GetRandom(100) { //! 开特殊图标
						self.Result[i][j] = 11
						spNum++
						if spNum == 3 {
							special = 3
						} else if spNum == 4 {
							special = 1
						} else if spNum >= 5 {
							special = 0
						}

					} else {

						self.Result[i][j] = icon[lib.HF_GetRandom(len(icon))]

					}
				}
			}
		}
		loop := 0
		for true && loop < 100 {
			self.Result[1] = [3]int{icon[lib.HF_GetRandom(len(icon))], icon[lib.HF_GetRandom(len(icon))], icon[lib.HF_GetRandom(len(icon))]}
			playerwin, _ := self.GetResultWin(self.Result)
			sysMoney := 0 - playerwin
			if self.Special1 == -1 {
				sysMoney += self.Person.Coin * LZDB_RATE
			}
			if GetServer().LzdbSysMoney[self.room.Type%160000]+int64(sysMoney) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && GetServer().LzdbSysMoney[self.room.Type%160000]+int64(sysMoney) <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
				break
			}
			loop++
		}
		playerwin, specialNum := self.GetResultWin(self.Result) //!　玩家赢钱

		special1 := -1
		if specialNum == 3 {
			goOn = false
			if self.Special1 == -1 {
				special1 = 15
			}
		} else if specialNum == 4 {
			goOn = false
			if self.Special1 == -1 {
				special1 = 20
			}
		} else if specialNum == 5 {
			goOn = false
			if self.Special1 == -1 {
				special1 = 25
			}
		}

		win := 0
		var msg Msg_GameLZDB_End
		msg.Uid = self.Person.Uid
		msg.Sex = self.Person.Sex
		msg.Name = self.Person.Name
		msg.IP = self.Person.IP
		msg.Head = self.Person.Head
		msg.Address = self.Person.Address
		msg.Special1 = special1
		msg.Special2 = self.Special2
		msg.Info = make([]Son_GameLZDB_End, 0)
		{
			if playerwin > 0 { //! 玩家赢钱
				win += playerwin
			}

			var son Son_GameLZDB_End
			son.Result = self.Result
			son.PersonWin = playerwin
			son.DelIcon = make([]int, 0)
			if goOn {
				goOn, son.DelIcon = self.DelIcon(true) //! 消掉得分图标
			} else {
				_, son.DelIcon = self.DelIcon(false)
			}
			msg.Info = append(msg.Info, son)
		}

		if goOn { //! 会继续下落
			loop := 0
			for true && loop < 10 {

				self.GetGoOnIcon(win) //! 填充空位

				playerwin, _ := self.GetResultWin(self.Result)

				if playerwin > 0 { //! 玩家赢钱
					win += playerwin
				}
				var son Son_GameLZDB_End
				son.Result = self.Result
				son.PersonWin = playerwin
				goOn, son.DelIcon = self.DelIcon(true) //! 消掉得分图标
				msg.Info = append(msg.Info, son)
				if !goOn { //! 没有消除的，不用继续下落了
					break
				}
				loop++
			}
		}

		if win == 0 {
			lostlst = append(lostlst, msg)
		} else if win > 0 && win < self.Person.Coin*LZDB_RATE {
			if GetServer().LzdbSysMoney[self.room.Type%160000]+int64(0-win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin && GetServer().LzdbSysMoney[self.room.Type%160000]+int64(0-win) <= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax {
				lostlst99 = append(lostlst99, msg)
			}

		} else {
			if GetServer().LzdbSysMoney[self.room.Type%160000]+int64(0-win) >= lib.GetManyMgr().GetProperty(self.room.Type).JackPotMin {
				winlst = append(winlst, msg)
			}
		}

		if len(lostlst) >= 10 && len(lostlst99) >= 10 && len(winlst) >= 10 {
			break
		}
	}

	var msg Msg_GameLZDB_End
	num := lib.HF_GetRandom(100)
	if GetServer().LzdbSysMoney[self.room.Type%160000] > lib.GetManyMgr().GetProperty(self.room.Type).JackPotMax && len(winlst) != 0 {
		msg = winlst[lib.HF_GetRandom(len(winlst))]
	} else {
		if num < 70 && len(lostlst99) != 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------玩家赢 0 ~ 本金")
			msg = lostlst99[lib.HF_GetRandom(len(lostlst99))]
		} else if num > 90 && len(winlst) != 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------玩家赢 > 本金")
			msg = winlst[lib.HF_GetRandom(len(winlst))]
		} else {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "-------玩家赢 0")
			msg = lostlst[lib.HF_GetRandom(len(lostlst))]
		}
	}

	if self.Special1 == -1 {
		self.Special1 = msg.Special1
	} else {
		msg.Special1 = self.Special1
	}
	dealwin := 0

	if self.Special1 == -1 {
		dealwin += self.Person.Coin * LZDB_RATE
	}

	for i := 0; i < len(msg.Info); i++ {
		if msg.Info[i].PersonWin > 0 {
			dealwin += (0 - msg.Info[i].PersonWin)
			self.Person.Win += msg.Info[i].PersonWin
		}
	}

	bets := 0
	if self.Special1 == -1 {
		bets = self.Person.Coin * LZDB_RATE
	}
	if (self.Person.Win-self.Person.Coin*LZDB_RATE > 0 && self.Special1 == -1) || (self.Special1 != -1 && self.Person.Win > 0) {
		if self.Special1 == -1 {
			self.Person.Cost += int(math.Ceil(float64(self.Person.Win-self.Person.Coin*LZDB_RATE) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		} else {
			self.Person.Cost += int(math.Ceil(float64(self.Person.Win) * lib.GetManyMgr().GetProperty(self.room.Type).Cost / 100.0))
		}

		self.Person.Win -= self.Person.Cost
		GetServer().SqlAgentGoldLog(self.Person.Uid, self.Person.Cost, self.room.Type)
		GetServer().SqlAgentBillsLog(self.Person.Uid, self.Person.Cost/2, self.room.Type)
	} else if self.Person.Win-bets < 0 {
		cost := int(math.Ceil(float64(self.Person.Win-bets) * float64(lib.GetManyMgr().GetProperty(self.room.Type).Cost) / 200.0))
		GetServer().SqlAgentBillsLog(self.Person.Uid, cost, self.room.Type)
	}
	self.Person.Total += self.Person.Win

	var record Rec_LZDB_Info
	record.GameType = self.room.Type
	record.Time = time.Now().Unix()
	var rec Son_Rec_LZDB_Person
	if self.Special1 == -1 {
		rec.Bets = self.Person.Coin * LZDB_RATE
	}
	rec.Uid = self.Person.Uid
	rec.Score = self.Person.Win - rec.Bets
	rec.Name = self.Person.Name
	rec.Head = self.Person.Head
	record.Info = append(record.Info, rec)
	GetServer().InsertRecord(self.room.Type, self.Person.Uid, lib.HF_JtoA(&record), rec.Score)

	if dealwin != 0 {
		GetServer().SqlBZWLog(&SQL_BZWLog{1, dealwin, time.Now().Unix(), self.room.Type})
	}

	if dealwin != 0 {
		cost := 0
		if dealwin > 0 { //! 玩家输钱，系统赢钱
			cost = int(math.Ceil(float64(dealwin) * float64(lib.GetManyMgr().GetProperty(self.room.Type).DealCost) / 100.0))
			dealwin -= cost
		}
		GetServer().SetLzdbSysMoney(self.room.Type%160000, GetServer().LzdbSysMoney[self.room.Type%160000]+int64(dealwin))
	}

	msg.Playerwin = self.Person.Win
	msg.Total = self.Person.Total
	self.room.broadCastMsg("gamelzdbend", &msg)

	if self.Special1 == 0 {
		self.Special1 = -1
	}
	self.Special2 = -1
	self.LZ = [3]int{0, 0, 0}
	self.Person.Win = 0
	self.Person.Cost = 0

	if self.Special1 > 0 {
		self.Special1--
		self.OnBegin()
	}
}

func (self *Game_LZDB) OnInit(room *Room) {
	self.room = room
}

func (self *Game_LZDB) OnRobot(robot *lib.Robot) {

}

func (self *Game_LZDB) OnSendInfo(person *Person) {
	if self.Person != nil && self.Person.Uid == person.Uid {
		self.Person.SynchroGold(person.Gold)
		person.SendMsg("gamelzdbinfo", self.getinfo(person.Uid))
		return
	}

	_person := new(Game_LZDB_Person)
	_person.Uid = person.Uid
	_person.Gold = person.Gold
	_person.Total = person.Gold
	_person.Name = person.Name
	_person.Head = person.Imgurl
	_person.IP = person.ip
	_person.Sex = person.Sex
	_person.Address = person.minfo.Address
	self.Person = _person
	person.SendMsg("gamelzdbinfo", self.getinfo(person.Uid))

	self.Time = time.Now().Unix() + 300
}

func (self *Game_LZDB) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold": //! 同步金币
		if self.Person.Uid == msg.V.(*staticfunc.Msg_SynchroGold).Uid {
			self.Person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal(self.Person.Uid, self.Person.Total)
		}
	case "lzdbstart":
		self.SetCoin(msg.Uid, msg.V.(*Msg_GameBets).Bets)
	}
}

func (self *Game_LZDB) OnBye() {

}

func (self *Game_LZDB) OnExit(uid int64) {
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

func (self *Game_LZDB) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_LZDB) OnIsBets(uid int64) bool {
	return false
}

func (self *Game_LZDB) OnBalance() {
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

func (self *Game_LZDB) OnTime() {
	if self.Person == nil {
		return
	}

	if self.Time == 0 {
		return
	}

	if time.Now().Unix() >= self.Time {
		self.room.KickViewByUid(self.Person.Uid, 96)
	}
}
