package gameserver

import (
	"lib"
	"staticfunc"
	"time"
)

/*
palam1
	个位 	0 不去牌 	 1 去2-4  	2 去2-6
	十位 	0 不翻车 	 1 翻车
	百位		0 不要特殊牌	 1 要特殊牌
	千位 	0 不带王		 1 带王
	万位 	2 2人 		 3 3人 		4 4人		 5 5人	 	6 6人
*/
/*
palam2 == 0 AA扣卡 	==1 房主扣卡
*/
/*
特殊牌型
0 非特殊牌型
1 三顺(0同花顺)	3
2 三顺(一同花顺)	5
3 四对子			8
4 四炸			8
5 一条龙			不去牌14 去牌13-14
6 三顺(双同花顺)	15
7 三顺(三同花顺)	30
8 双四条炸弹		50
*/
//! 记录
type Rec_GameSSD struct {
	Info    []Son_Rec_GameSSD `json:"info"`
	Roomid  int               `json:"roomid"`
	MaxStep int               `json:"maxstep"`
	Param1  int               `json:"param1"`
	Param2  int               `json:"param2"`
	Time    int64             `json:"time"`
}

type Son_Rec_GameSSD struct {
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Head    string `json:"head"`
	Card    []int  `json:"card"`
	Bets    int    `json:"bets"`
	Dealer  bool   `json:"dealer"`
	Score   int    `json:"score"`
	Total   int    `json:"total"`
	Type    int    `json:"type"`
	RobDeal int    `json:"robdeal"`
	Special int    `json:"special"`
}

type Game_SSD struct {
	PersonMgr []*Game_SSD_Person `json:"personmgr"` //! 玩家信息
	ShowCard  bool               `json:"showcard"`  //! 是否可见其他玩家的牌
	MaxPerson int                `json:"maxperson"` //! 房间最大人数
	Special   bool               `json:"special"`   //! 是否要特殊牌
	FanChe    bool               `json:"fanche"`    //! 是否翻车
	State     int                `json:"state"`     //! 状态 0初始 1未开始 2开始

	room *Room
}

type Game_SSD_Person struct {
	Uid int64 `json:"uid"`

	Ready       bool  `json:"ready"`       //! 是否准备
	Match       bool  `json:"match"`       //! 配牌是否完成
	FanChe      bool  `json:"fanche"`      //! 是否翻车
	Score       int   `json:"score"`       //! 全局分数
	CurScore    int   `json:"curscore"`    //! 单局分数
	HeadScore   int   `json:"headscore"`   //! 头道得分
	BodyScore   int   `json:"bodyscore"`   //! 中道得分
	TailScore   int   `json:"tailscore"`   //! 尾道得分
	Cards       []int `json:"cards"`       //! 手牌
	HeadCards   []int `json:"headcards"`   //! 头道牌
	BodyCards   []int `json:"bodycards"`   //! 中道牌
	TailCards   []int `json:"tailcards"`   //! 尾道牌
	HeadType    int   `json:"headtype"`    //! 头道类型
	BodyType    int   `json:"bodytype"`    //! 中道类型
	TailType    int   `json:"tailtype"`    //! 尾道类型
	SpecialCard int   `json:"specialcard"` //! 特殊牌型
	FanNum      int   `json:"fannum"`
	Winner      int   `json:"winner"`
	Loser       int   `json:"loser"`
	DogFall     int   `json:"dogfall"`

	//////////////////////////////////////
	SpecialCards []int `json:"specialcards"` //定制牌型
}

type Msg_GameSSD_Info struct {
	Begin bool               `json:"begin"` //! 是否开始
	State int                `json:"state"` //! 状态 0初始 1未开始 2开始
	Info  []Son_GameSSD_Info `json:"info"`
}
type Son_GameSSD_Info struct {
	Uid         int64 `json:"uid"`
	Ready       bool  `json:"ready"`       //! 是否准备
	Match       bool  `json:"match"`       //! 配牌是否完成
	FanChe      bool  `json:"fanche"`      //! 是否翻车
	Score       int   `json:"score"`       //! 全局分数
	CurScore    int   `json:"curscore"`    //! 单局分数
	HeadScore   int   `json:"headscore"`   //! 头道得分
	BodyScore   int   `json:"bodyscore"`   //! 中道得分
	TailScore   int   `json:"tailscore"`   //! 尾道得分
	Card        []int `json:"card"`        //! 手牌
	HeadCard    []int `json:"headcard"`    //! 头道牌
	BodyCard    []int `json:"bodycard"`    //! 中道牌
	TailCard    []int `json:"tailcard"`    //! 尾道牌
	HeadType    int   `json:"headtype"`    //! 头道类型
	BodyType    int   `json:"bodytype"`    //! 中道类型
	TailType    int   `json:"tailtype"`    //! 尾道类型
	SpecialCard int   `json:"specialcard"` //! 特殊牌型

}

type Msg_GameSSD_Match struct {
	Uid int64 `json:"uid"`
}

type Msg_GameSSD_Bye struct {
	Info []Son_GameSSD_Bye `json:"info"`
}
type Son_GameSSD_Bye struct {
	Uid     int64 `json:"uid"`
	Score   int   `json:"score"`
	Winner  int   `json:"winner"`
	Loser   int   `json:"loser"`
	DogFall int   `json:"dogfall"`
}

type Msg_GameSSD_End struct {
	Info []Son_GameSSD_End `json:"info"`
}
type Son_GameSSD_End struct {
	Uid        int64   `json:"uid"`
	HeadPX     int     `json:"headpx"`     //! 头道牌类型
	BodyPX     int     `json:"bodypx"`     //! 中道类型
	TailPX     int     `json:"tailpx"`     //! 尾道类型
	HeadScore  int     `json:"headscore"`  //! 头道总得分
	BodyScore  int     `json:"bodyscore"`  //! 中道总得分
	TailScore  int     `json:"tailscore"`  //! 尾道总得分
	GunScore   int     `json:"gunscore"`   //! 翻车总得分
	BeGunScore int     `json:"begunscore"` //! 被翻车总分
	GunNum     int     `json:"gunnum"`     //! 翻车次数
	Card       []int   `json:"card"`       //! 手牌
	Score      int     `json:"score"`      //! 单局分
	Total      int     `json:"total"`      //! 全局分
	GunUid     []int64 `json:"gunuid"`     //! 翻车的人
	SwatScore  int     `json:"swatscore"`  //! 红波浪得分
	Special    int     `json:"special"`    //! 特殊牌型
}

func NewGame_SSD() *Game_SSD {
	game := new(Game_SSD)
	game.PersonMgr = make([]*Game_SSD_Person, 0)
	return game
}

func (self *Game_SSD) Init() {
	self.MaxPerson = self.room.Param1 / 10000
	if self.room.Param1/100%10 == 0 {
		self.Special = false
	} else {
		self.Special = true
	}
	if self.room.Param1/10%10 == 0 {
		self.FanChe = false
	} else {
		self.FanChe = true
	}
}

func (self *Game_SSD) GetSpecialScore(tmp int, card []int) int {
	switch tmp {
	case 1:
		return 3
	case 2:
		return 5
	case 3:
		return 8
	case 4:
		return 8
	case 5:
		if self.room.Param1%10 == 0 {
			return 14
		} else {
			score := 0
			for i := 0; i < len(card); i++ {
				if card[i]/10 == 1 {
					card[i] = 141
				}
				if card[i]/10 > score {
					score = card[i] / 10
				}
			}
			return score
		}
	case 6:
		return 15
	case 7:
		return 30
	case 8:
		return 50
	}
	return 0
}

//! 计算玩家与另玩家比牌后的得分
func (self *Game_SSD) Compare(person *Game_SSD_Person, _person *Game_SSD_Person) (int, int, int) {
	win := true //! true person赢  false _person赢
	//! 计算头道分
	headScore := 0
	headCard := person.HeadCards
	_headCard := _person.HeadCards

	if person.HeadType%10 > _person.HeadType%10 { //! person 类型大
		win = true
	} else if person.HeadType%10 < _person.HeadType%10 { //! person 类型小
		win = false
	} else {
		if person.HeadType%10 == 1 { //! 同为乌龙
			if person.HeadType/10 > _person.HeadType/10 {
				win = true
			} else if person.HeadType/10 < _person.HeadType/10 {
				win = false
			} else { //! 最大点数相同 ，比另一张牌
				for i := 0; i < len(headCard); i++ {
					if headCard[i]/10 == 1 {
						headCard[i] = 14*10 + headCard[i]%10
					}
					if _headCard[i]/10 == 1 {
						_headCard[i] = 14*10 + _headCard[i]%10
					}
				}
				if headCard[0] < headCard[1] {
					tmp := headCard[1]
					headCard[1] = headCard[0]
					headCard[0] = tmp
				}
				if _headCard[0] < _headCard[1] {
					_tmp := _headCard[1]
					_headCard[1] = _headCard[0]
					_headCard[0] = _tmp
				}
				if headCard[1]/10 > _headCard[1]/10 { //! person第二张牌大
					win = true
				} else if headCard[1]/10 < _headCard[1]/10 { //! person第二张牌小
					win = false
				} else {
					if headCard[0]%10 > _headCard[0]%10 { //! person最大牌花色大
						win = true
					} else { //! person最大牌花色小
						win = false
					}
				}
			}
		} else if person.HeadType%10 == 2 { //! 同为对子
			if person.HeadType/10 > _person.HeadType/10 {
				win = true
			} else if person.HeadType/10 < _person.HeadType/10 {
				win = false
			} else if person.HeadType/10 == _person.HeadType/10 { //! 对子一样大比花色
				if headCard[0] < headCard[1] {
					tmp := headCard[0]
					headCard[0] = headCard[1]
					headCard[1] = tmp
				}
				if _headCard[0] < _headCard[1] {
					tmp := _headCard[0]
					_headCard[0] = _headCard[1]
					_headCard[1] = tmp
				}
				for i := 0; i < len(headCard); i++ {
					if headCard[i] > _headCard[i] {
						win = true
						break
					} else if headCard[i] < _headCard[i] {
						win = false
						break
					}
				}
			}
		}
	}

	if win {
		if person.HeadType%10 == 1 { //! 乌龙
			headScore += 1
		} else if person.HeadType%10 == 2 { //! 对子
			headScore += person.HeadType / 10
		}
	} else {
		if _person.HeadType%10 == 1 { //! 乌龙
			headScore -= 1
		} else if _person.HeadType%10 == 2 { //! 对子
			headScore -= _person.HeadType / 10
		}
	}

	//! 计算中道分
	win = true
	bodyScore := 0
	bodyCard := person.BodyCards
	_bodyCard := _person.BodyCards
	if person.BodyType%10 > _person.BodyType%10 { //! person 类型大
		win = true
	} else if person.BodyType%10 < _person.BodyType%10 { //! person 类型小
		win = false
	} else {
		if person.BodyType%10 == 3 || person.BodyType%10 == 5 { //! 同为 顺子 || 同花顺
			if person.BodyType/10 > _person.BodyType/10 { //! person的顺子大
				win = true
			} else if person.BodyType/10 < _person.BodyType/10 { //! person的顺子小
				win = false
			} else { //! 顺子一样大
				/*
					for i := 0; i < len(bodyCard); i++ { //! 降序排列
						for j := i + 1; j < len(bodyCard); j++ {
							if bodyCard[j]/10 > bodyCard[i]/10 {
								temp := bodyCard[i]
								bodyCard[i] = bodyCard[j]
								bodyCard[j] = temp
							}
						}
					}
					for i := 0; i < len(_bodyCard); i++ { //! 降序排列
						for j := i + 1; j < len(_bodyCard); j++ {
							if _bodyCard[j]/10 > _bodyCard[i]/10 {
								temp := _bodyCard[i]
								_bodyCard[i] = _bodyCard[j]
								_bodyCard[j] = temp
							}
						}
					}
					for i := 0; i < len(bodyCard); i++ {
						if bodyCard[i]%10 > _bodyCard[i]%10 {
							win = true
						} else if bodyCard[i]%10 < _bodyCard[i]%10 {
							win = false
						}
					}
				*/

				num, _num := 0, 0
				for i := 0; i < len(bodyCard); i++ {
					if bodyCard[i]/10 == person.BodyType/10 || (bodyCard[i]/10 == 1 && (person.BodyType/10 == 14 || person.BodyType/10 == 15)) {
						num = bodyCard[i]
					}
					if _bodyCard[i]/10 == _person.BodyType/10 || (_bodyCard[i]/10 == 1 && (_person.BodyType/10 == 14 || _person.BodyType/10 == 15)) {
						_num = _bodyCard[i]
					}
				}
				if num%10 > _num%10 { //! person的最大牌的花色大
					win = true
				} else if num%10 < _num%10 {
					win = false
				} else {
					find := false
					cards := person.Cards[2:5]
					for i := 0; i < len(cards); i++ {
						if num == cards[i] {
							win = true
							find = true
							break
						}
					}
					if !find {
						win = false
					}
				}

			}

		} else if person.BodyType%10 == 4 { //! 同为炸弹
			if person.BodyType/10 > _person.BodyType/10 { //! person 的炸弹大
				win = true
			} else if person.BodyType/10 < _person.BodyType/10 {
				win = false
			} else {
				for i := 0; i < len(bodyCard); i++ { //! 降序排列
					for j := i + 1; j < len(bodyCard); j++ {
						if bodyCard[j]/10 > bodyCard[i]/10 {
							temp := bodyCard[i]
							bodyCard[i] = bodyCard[j]
							bodyCard[j] = temp
						}
					}
				}
				for i := 0; i < len(_bodyCard); i++ { //! 降序排列
					for j := i + 1; j < len(_bodyCard); j++ {
						if _bodyCard[j]/10 > _bodyCard[i]/10 {
							temp := _bodyCard[i]
							_bodyCard[i] = _bodyCard[j]
							_bodyCard[j] = temp
						}
					}
				}
				for i := 0; i < len(bodyCard); i++ {
					if bodyCard[i] > _bodyCard[i] {
						win = true
					} else if bodyCard[i] < _bodyCard[i] {
						win = false
					}
				}
			}
		} else if person.BodyType%10 == 1 { //! 同为乌龙
			for i := 0; i < len(bodyCard); i++ {
				if bodyCard[i]/10 == 1 {
					bodyCard[i] = 14*10 + bodyCard[i]%10
				}
				if _bodyCard[i]/10 == 1 {
					_bodyCard[i] = 14*10 + _bodyCard[i]%10
				}
			}

			for i := 0; i < len(bodyCard); i++ { //! 降序排列
				for j := i + 1; j < len(bodyCard); j++ {
					if bodyCard[j]/10 > bodyCard[i]/10 {
						temp := bodyCard[i]
						bodyCard[i] = bodyCard[j]
						bodyCard[j] = temp
					}
				}
			}
			for i := 0; i < len(_bodyCard); i++ { //! 降序排列
				for j := i + 1; j < len(_bodyCard); j++ {
					if _bodyCard[j]/10 > _bodyCard[i]/10 {
						temp := _bodyCard[i]
						_bodyCard[i] = _bodyCard[j]
						_bodyCard[j] = temp
					}
				}
			}
			tong := true
			for i := 0; i < len(bodyCard); i++ { //! 比牌的大小
				if bodyCard[i]/10 == _bodyCard[i]/10 { //! 一样大比下一张
					continue
				} else if bodyCard[i]/10 > _bodyCard[i]/10 { //! person的牌大
					win = true
					tong = false
					break
				} else if bodyCard[i]/10 < _bodyCard[i]/10 { //! person的牌小
					win = false
					tong = false
					break
				}
			}
			if tong { //! 三张牌都一样 ，比花色
				if bodyCard[0]%10 > _bodyCard[0]%10 { //! person的最大牌花色较大
					win = true
				} else { //! person 的最大牌花色较小
					win = false
				}
			}

		} else if person.BodyType%10 == 2 { //! 同为对子
			if person.BodyType/10 > _person.BodyType/10 { //! person的对子大
				win = true
			} else if person.BodyType/10 < _person.BodyType/10 { //! person的对子小
				win = false
			} else { //! 对子一样大

				num, _num := 0, 0
				for i := 0; i < len(bodyCard); i++ {
					if bodyCard[i]/10 != person.BodyType/10 {
						if person.BodyType/10 == 14 && bodyCard[i]/10 == 1 {

						} else {
							num = bodyCard[i]
						}
					}
					if _bodyCard[i]/10 != _person.BodyType/10 {
						if _person.BodyType/10 == 14 && _bodyCard[i]/10 == 1 {

						} else {
							_num = _bodyCard[i]
						}
					}
				}
				if num/10 == 1 {
					num = 140
				}
				if _num/10 == 1 {
					_num = 140
				}
				if num/10 > _num/10 { //! person的单张大
					win = true
				} else if num/10 < _num/10 { //! person的单张小
					win = false
				} else {
					num = 0
					_num = 0
					for i := 0; i < len(bodyCard); i++ {
						if bodyCard[i]/10 == 1 && person.BodyType/10 == 14 || (bodyCard[i]/10 == person.BodyType/10) {
							if num < bodyCard[i]%10 {
								num = bodyCard[i] % 10
							}
						}
						if _bodyCard[i]/10 == 1 && _person.BodyType/10 == 14 || (_bodyCard[i]/10 == _person.BodyType/10) {
							if _num < _bodyCard[i]%10 {
								_num = _bodyCard[i] % 10
							}
						}
					}
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "num :", num, " card :", bodyCard, " _num: ", _num, " _card :", _bodyCard)
					if num < _num {
						win = false
					} else {
						win = true
					}
				}
			}
		}
	}

	if win {
		if person.BodyType%10 == 1 || person.BodyType%10 == 2 || person.BodyType%10 == 3 { //! 乌龙 || 对子 || 顺子
			bodyScore += 1
		} else if person.BodyType%10 == 4 { //! 炸弹
			bodyScore += 8
		} else if person.BodyType%10 == 5 { //! 同花顺
			bodyScore += 10
		}
	} else {
		if _person.BodyType%10 == 1 || _person.BodyType%10 == 2 || _person.BodyType%10 == 3 { //! 乌龙 || 对子 || 顺子
			bodyScore -= 1
		} else if _person.BodyType%10 == 4 { //!  炸弹
			bodyScore -= 8
		} else if _person.BodyType%10 == 5 { //! 同花顺
			bodyScore -= 10
		}
	}

	//! 计算尾道分
	win = true
	tailScore := 0
	tailCard := person.TailCards
	_tailCard := _person.TailCards
	if person.TailType%10 > _person.TailType%10 { //! person 类型大
		win = true
	} else if person.TailType%10 < _person.TailType%10 { //! person 类型小
		win = false
	} else {
		if person.TailType%10 == 3 || person.TailType%10 == 5 { //! 同为 顺子 || 同花顺

			if person.TailType/10 > _person.TailType/10 { //! person的顺子大
				win = true
			} else if person.TailType/10 < _person.TailType/10 { //! person的顺子小
				win = false
			} else { //! 顺子一样大
				num, _num := 0, 0
				for i := 0; i < len(tailCard); i++ {
					if tailCard[i]/10 == person.TailType/10 || (tailCard[i]/10 == 1 && (person.TailType/10 == 14 || person.TailType/10 == 15)) {
						num = tailCard[i]
					}
					if _tailCard[i]/10 == _person.TailType/10 || (_tailCard[i]/10 == 1 && (_person.TailType/10 == 14 || _person.TailType/10 == 15)) {
						_num = _tailCard[i]
					}
				}
				if num%10 > _num%10 { //! person的最大牌的花色大
					win = true
				} else if num%10 < _num%10 {
					win = false
				} else {
					find := false
					cards := person.Cards[5:]
					for i := 0; i < len(cards); i++ {
						if num == cards[i] {
							win = true
							find = true
							break
						}
					}
					if !find {
						win = false
					}
				}
			}

		} else if person.TailType%10 == 4 { //! 同为炸弹
			if person.TailType/10 > _person.TailType/10 { //! person 的炸弹大
				win = true
			} else if person.TailType/10 < _person.TailType/10 {
				win = false
			} else {
				for i := 0; i < len(tailCard); i++ { //! 降序排列
					for j := i + 1; j < len(tailCard); j++ {
						if tailCard[j]/10 > tailCard[i]/10 {
							temp := tailCard[i]
							tailCard[i] = tailCard[j]
							tailCard[j] = temp
						}
					}
				}
				for i := 0; i < len(_tailCard); i++ { //! 降序排列
					for j := i + 1; j < len(_tailCard); j++ {
						if _tailCard[j]/10 > _tailCard[i]/10 {
							temp := _tailCard[i]
							_tailCard[i] = _tailCard[j]
							_tailCard[j] = temp
						}
					}
				}
				for i := 0; i < len(tailCard); i++ {
					if tailCard[i] > _tailCard[i] {
						win = true
					} else if tailCard[i] < _tailCard[i] {
						win = false
					}
				}
			}
		} else if person.TailType%10 == 1 { //! 同为乌龙
			for i := 0; i < len(tailCard); i++ {
				if tailCard[i]/10 == 1 {
					tailCard[i] = 14*10 + tailCard[i]%10
				}
				if _tailCard[i]/10 == 1 {
					_tailCard[i] = 14*10 + _tailCard[i]%10
				}
			}

			for i := 0; i < len(tailCard); i++ { //! 降序排列
				for j := i + 1; j < len(tailCard); j++ {
					if tailCard[j]/10 > tailCard[i]/10 {
						temp := tailCard[i]
						tailCard[i] = tailCard[j]
						tailCard[j] = temp
					}
				}
			}
			for i := 0; i < len(_tailCard); i++ { //! 降序排列
				for j := i + 1; j < len(_tailCard); j++ {
					if _tailCard[j]/10 > _tailCard[i]/10 {
						temp := _tailCard[i]
						_tailCard[i] = _tailCard[j]
						_tailCard[j] = temp
					}
				}
			}
			tong := true
			for i := 0; i < len(tailCard); i++ { //! 比牌的大小
				if tailCard[i]/10 == _tailCard[i]/10 { //! 一样大比下一张
					continue
				} else if tailCard[i]/10 > _tailCard[i]/10 { //! person的牌大
					win = true
					tong = false
					break
				} else if tailCard[i]/10 < _tailCard[i]/10 { //! person的牌小
					win = false
					tong = false
					break
				}
			}
			if tong { //! 三张牌都一样 ，比花色
				if tailCard[0]%10 > _tailCard[0]%10 { //! person的最大牌花色较大
					win = true
				} else { //! person 的最大牌花色较小
					win = false
				}
			}

		} else if person.TailType%10 == 2 { //! 同为对子
			if person.TailType/10 > _person.TailType/10 { //! person的对子大
				win = true
			} else if person.TailType/10 < _person.TailType/10 { //! person的对子小
				win = false
			} else { //! 对子一样大
				num, _num := 0, 0
				for i := 0; i < len(tailCard); i++ {
					if tailCard[i]/10 != person.TailType/10 {
						if person.TailType/10 == 14 && tailCard[i]/10 == 1 {

						} else {
							num = tailCard[i]
						}
					}
					if _tailCard[i]/10 != _person.TailType/10 {
						if _person.TailType/10 == 14 && _tailCard[i]/10 == 1 {

						} else {
							_num = _tailCard[i]
						}
					}
				}
				if num/10 == 1 {
					num = 140
				}
				if _num/10 == 1 {
					_num = 140
				}
				if num/10 > _num/10 { //! person的单张大
					win = true
				} else if num/10 < _num/10 { //! person的单张小
					win = false
				} else {
					num = 0
					_num = 0
					for i := 0; i < len(tailCard); i++ {
						if tailCard[i]/10 == 1 && person.TailType/10 == 14 || (tailCard[i]/10 == person.TailType/10) {
							if num < tailCard[i]%10 {
								num = tailCard[i] % 10
							}
						}
						if _tailCard[i]/10 == 1 && _person.TailType/10 == 14 || (_tailCard[i]/10 == _person.TailType/10) {
							if _num < _tailCard[i]%10 {
								_num = _tailCard[i] % 10
							}
						}
					}
					if num < _num {
						win = false
					} else {
						win = true
					}
				}
			}
		}
	}

	if win {
		if person.TailType%10 == 1 || person.TailType%10 == 2 || person.TailType%10 == 3 {
			tailScore += 1
		} else if person.TailType%10 == 4 {
			tailScore += 4
		} else if person.TailType%10 == 5 {
			tailScore += 5
		}
	} else {
		if _person.TailType%10 == 1 || _person.TailType%10 == 2 || _person.TailType%10 == 3 {
			tailScore -= 1
		} else if _person.TailType%10 == 4 {
			tailScore -= 4
		} else if _person.TailType%10 == 5 {
			tailScore -= 5
		}
	}

	return headScore, bodyScore, tailScore
}

func (self *Game_SSD) IsMatch(person *Game_SSD_Person) bool { //! 配牌是否正确
	//! 头道 vs 中道
	if person.HeadType%10 > person.BodyType%10 { //! 头道牌型大于中道
		return false
	} else if person.HeadType%10 == person.BodyType%10 { //! 头道类型等于中道
		if person.HeadType%10 == 2 {
			return true
		}
		if person.HeadType/10 > person.BodyType/10 { //! 头道点数大于中道
			return false
		} else if person.HeadType/10 == person.BodyType/10 { //! 点数相同
			if person.HeadType%10 == 1 { //! 乌龙
				card := make([]int, 0)
				_card := make([]int, 0)
				card = person.HeadCards
				_card = person.BodyCards

				for i := 0; i < len(card); i++ { //! 降序排列
					for j := i + 1; j < len(card); j++ {
						if card[j]/10 > card[i]/10 {
							temp := card[i]
							card[i] = card[j]
							card[j] = temp
						}
					}
				}

				for i := 0; i < len(_card); i++ { //! 降序排列
					for j := i + 1; j < len(_card); j++ {
						if _card[j]/10 > _card[i]/10 {
							temp := _card[i]
							_card[i] = _card[j]
							_card[j] = temp
						}
					}
				}

				if card[1]/10 > _card[1]/10 { //! 第二大的牌头道大
					return false
				} else if card[1]/10 == _card[1]/10 { //! 第二大的牌相等
					if card[0]%10 > _card[0]%10 { //! 头道第一大牌花色大
						return false
					}
				}

			} else { //! 对子
				tao := false
				for i := 0; i < len(person.HeadCards); i++ {
					if person.HeadCards[i]%10 == 4 {
						tao = true
						break
					}
				}
				if !tao {
					return false
				}
			}
		}
	}

	//! 中道 vs 尾道
	if person.BodyType%10 > person.TailType%10 { //! 中道牌型大于尾道
		return false
	} else if person.BodyType%10 == person.TailType%10 { //! 中道类型等于尾道
		if person.BodyType/10 > person.TailType/10 { //! 中道最大点数大于尾道
			return false
		} else if person.BodyType/10 == person.TailType/10 { //! 中道最大点数等于尾道
			if person.BodyType%10 == 3 || person.BodyType%10 == 5 { //! 同花顺||顺子
				cardNum, _cardNum := 0, 0
				for i := 0; i < len(person.BodyCards); i++ {
					if person.BodyCards[i]/10 == person.BodyType/10 { //! 中道最大牌花色点数
						cardNum = person.BodyCards[i]
					}
					if person.TailCards[i]/10 == person.TailType/10 { //! 尾道最大牌花色点数
						_cardNum = person.TailCards[i]
					}
				}
				if cardNum%10 > _cardNum%10 { //! 中道花色大于尾道
					return false
				}
			}

			if person.BodyType%10 == 2 { //! 对子
				dan, _dan := 0, 0
				for i := 0; i < len(person.BodyCards); i++ {
					if person.BodyCards[i]/10 != person.BodyType/10 { //! 中道单张点数花色
						dan = person.BodyCards[i]
					}

					if person.TailCards[i]/10 != person.TailType/10 { //! 尾道单张点数花色
						_dan = person.TailCards[i]
					}

				}
				if dan/10 > _dan/10 { //! 中道单张 大于 尾道单张
					return false
				} else if dan/10 == _dan/10 { //! 单张相同
					if dan%10 > _dan%10 { //! 中道单张花色大
						return false
					}
				}
			}

			if person.BodyType%10 == 1 { //! 乌龙
				card := person.BodyCards
				_card := person.TailCards

				for i := 0; i < len(card); i++ {
					if card[i]/10 == 1 {
						card[i] = 14*10 + card[i]%10
					}
					if _card[i]/10 == 1 {
						_card[i] = 14*10 + _card[i]%10
					}
				}

				for i := 0; i < len(card); i++ { //! 降序排列
					for j := i + 1; j < len(card); j++ {
						if card[j]/10 > card[i]/10 {
							temp := card[i]
							card[i] = card[j]
							card[j] = temp
						}
					}
				}

				for i := 0; i < len(_card); i++ { //! 降序排列
					for j := i + 1; j < len(_card); j++ {
						if _card[j]/10 > _card[i]/10 {
							temp := _card[i]
							_card[i] = _card[j]
							_card[j] = temp
						}
					}
				}

				equ := true
				for i := 0; i < len(card); i++ {
					if card[i]/10 != _card[i]/10 {
						equ = false
					}
					if card[i]/10 > _card[i]/10 { //! 中道比尾道点数大
						return false
					} else if card[i]/10 < _card[i]/10 {
						break
					}
				}
				if equ { //! 三张全相等
					if card[0]%10 > _card[0]%10 { //! 最大单张中道的花色大
						return false
					}
				}
			}
		}
	}

	return true
}

func (self *Game_SSD) getInfo(uid int64) *Msg_GameSSD_Info {
	var msg Msg_GameSSD_Info
	msg.Begin = self.room.Begin
	msg.State = self.State
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameSSD_Info
		son.Uid = self.PersonMgr[i].Uid
		son.Ready = self.PersonMgr[i].Ready

		son.Match = self.PersonMgr[i].Match
		son.FanChe = self.PersonMgr[i].FanChe
		son.Score = self.PersonMgr[i].Score
		son.CurScore = self.PersonMgr[i].CurScore
		son.HeadScore = self.PersonMgr[i].HeadScore
		son.BodyScore = self.PersonMgr[i].BodyScore
		son.TailScore = self.PersonMgr[i].TailScore
		if self.PersonMgr[i].Uid == uid || self.ShowCard { //! 自己的手牌 || 已经开始比牌
			son.Card = self.PersonMgr[i].Cards

			if len(self.PersonMgr[i].SpecialCards) > 0 {
				son.Card = self.PersonMgr[i].SpecialCards
			}

			son.HeadCard = self.PersonMgr[i].HeadCards
			son.BodyCard = self.PersonMgr[i].BodyCards
			son.TailCard = self.PersonMgr[i].TailCards
			son.SpecialCard = self.PersonMgr[i].SpecialCard
			son.HeadType = self.PersonMgr[i].HeadType
			son.BodyType = self.PersonMgr[i].BodyType
			son.TailType = self.PersonMgr[i].TailType
		} else {
			son.Card = make([]int, 8)
			son.HeadCard = make([]int, 2)
			son.BodyCard = make([]int, 3)
			son.TailCard = make([]int, 3)
			son.SpecialCard = 0
			son.HeadType = 0
			son.BodyType = 0
			son.TailType = 0
		}

		msg.Info = append(msg.Info, son)
	}
	return &msg
}

func (self *Game_SSD) GetPerson(uid int64) *Game_SSD_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}
	return nil
}

func (self *Game_SSD) Ready(uid int64) {
	if self.room.IsBye() {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "房间已结束")
		return
	}

	if self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏已经开始，不能加入!")
		return
	}

	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Ready {
				return
			} else {
				self.PersonMgr[i].Ready = true
				num++
			}
		} else if self.PersonMgr[i].Ready {
			num++
		}
	}

	if num == self.MaxPerson { //! 准备的人数达到游戏人数
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
		self.OnBegin()
		return
	}

	var msg staticfunc.Msg_Uid
	msg.Uid = uid
	self.room.broadCastMsg("gamessdready", &msg)

	self.room.flush()
}

func (self *Game_SSD) GameMatch(uid int64, card []int, abscard []int, special int) { //! 组牌
	/*
		special
		1 三顺
		2 四条炸弹
		3 四对
		4 双条炸弹
		5 一条龙
	*/
	var _card []int

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "进入gamematch")
	person := self.GetPerson(uid)
	if person == nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "没找到玩家")
		return
	}
	if self.room.Param1/1000%10 == 0 {
		abscard = card
	}

	if len(card) != len(abscard) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "card 和 abscard长度不一样")
		return
	}

	if person.Match {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能重复配牌")
		return
	}
	lib.HF_DeepCopy(&_card, &abscard)
	tmpcard := make([]int, 0)
	lib.HF_DeepCopy(&tmpcard, &abscard)
	for _, value := range card {
		find := false
		for i := 0; i < len(tmpcard); i++ {
			if tmpcard[i] == value {
				find = true
				copy(tmpcard[i:], tmpcard[i+1:])
				tmpcard = tmpcard[:len(tmpcard)-1]
				break
			}
		}

		if find {
			continue
		}

		for i := 0; i < len(tmpcard); i++ {
			if tmpcard[i] == 1000 || tmpcard[i] == 2000 {
				find = true
				copy(tmpcard[i:], tmpcard[i+1:])
				tmpcard = tmpcard[:len(tmpcard)-1]
				break
			}
		}

		if !find {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "癞子错误:", card, ",", abscard)
			return
		}
	}

	for _, value := range abscard {
		find := false
		for i := 0; i < len(person.Cards); i++ {
			if person.Cards[i] == value {
				find = true
				break
			}
		}
		if !find {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, card, ":出牌找不到")
			return
		}
	}

	if special != 0 && self.room.Param1/100%10 == 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "对局不要特殊牌")
		return
	}

	if special != 0 { //! 有特殊牌型
		person.SpecialCard = GetSSDSpecialType(card, special)
		if person.SpecialCard == 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "该玩家没有特殊牌型")
			return
		}
		if special == 1 {
			person.HeadCards = card[0:2]
			person.BodyCards = card[2:5]
			person.TailCards = card[5:]
		}
	} else { //! 没有特殊牌型
		person.HeadCards = card[0:2]
		person.BodyCards = card[2:5]
		person.TailCards = card[5:]

		person.HeadType = GetSSDCardType(person.HeadCards)
		person.BodyType = GetSSDCardType(person.BodyCards)
		person.TailType = GetSSDCardType(person.TailCards)

		if !self.IsMatch(person) {
			_person := GetPersonMgr().GetPerson(person.Uid)
			if _person == nil {
				return
			}
			_person.SendErr("配牌出错")
			return
		}
	}

	lib.HF_DeepCopy(&person.Cards, &_card)

	person.Match = true
	var msg Msg_GameSSD_Match
	msg.Uid = person.Uid
	self.room.broadCastMsg("gamessdmatch", &msg)

	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Match {
			num++
		}
	}
	if num == self.MaxPerson { //! 配牌完成
		self.ShowCard = true
		self.OnEnd()
	}

	self.room.flush()
}

func (self *Game_SSD) OnBegin() { //! 一局开始
	if self.room.IsBye() {
		return
	}
	self.room.SetBegin(true)
	self.ShowCard = false

	self.State = 2
	for i := 0; i < len(self.PersonMgr); i++ { //! 重新初始化人
		self.PersonMgr[i].Match = false
		self.PersonMgr[i].FanChe = false
		self.PersonMgr[i].CurScore = 0
		self.PersonMgr[i].HeadScore = 0
		self.PersonMgr[i].BodyScore = 0
		self.PersonMgr[i].TailScore = 0
		self.PersonMgr[i].Cards = make([]int, 0)
		self.PersonMgr[i].HeadCards = make([]int, 0)
		self.PersonMgr[i].BodyCards = make([]int, 0)
		self.PersonMgr[i].TailCards = make([]int, 0)
		self.PersonMgr[i].HeadType = 0
		self.PersonMgr[i].BodyType = 0
		self.PersonMgr[i].TailType = 0
		self.PersonMgr[i].FanNum = 0
		self.PersonMgr[i].SpecialCard = 0

	}

	cardMsg := NewCard_SSD(self.room.Param1/1000%10, self.room.Param1%10)

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Cards = cardMsg.Deal(8)

		////////////////
		if len(self.PersonMgr[i].SpecialCards) > 0 {
			self.PersonMgr[i].Cards = self.PersonMgr[i].SpecialCards
		}
		//////////////////
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}

		person.SendMsg("gamessdbegin", self.getInfo(person.Uid))
	}

	self.room.flush()
}

func (self *Game_SSD) OnEnd() { //! 一局结束
	self.room.SetBegin(false)
	self.State = 1

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Ready = false
		self.PersonMgr[i].SpecialCards = make([]int, 0)
	}

	fan := false
	if self.room.Param1/10%10 == 1 {
		fan = true
	}
	red := false
	special := false
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].SpecialCard != 0 {
			special = true
			break
		}
	}
	if !special && self.MaxPerson > 2 && self.room.Param1/10%10 == 1 {
		red = true
	}
	var record Rec_GameSSD
	record.Time = time.Now().Unix()
	record.Roomid = self.room.Id*100 + self.room.Step
	record.MaxStep = self.room.MaxStep
	record.Param1 = self.room.Param1
	record.Param2 = self.room.Param2

	var msg Msg_GameSSD_End
	msg.Info = make([]Son_GameSSD_End, 0)
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameSSD_End
		son.Uid = self.PersonMgr[i].Uid
		son.HeadPX = self.PersonMgr[i].HeadType
		son.BodyPX = self.PersonMgr[i].BodyType
		son.TailPX = self.PersonMgr[i].TailType
		son.Card = self.PersonMgr[i].Cards
		son.Special = self.PersonMgr[i].SpecialCard
		son.GunUid = make([]int64, 0)
		for j := 0; j < len(self.PersonMgr); j++ {
			if self.PersonMgr[i].Uid == self.PersonMgr[j].Uid {
				continue
			}
			if self.PersonMgr[i].SpecialCard == 0 { //! 非特殊牌型
				if self.PersonMgr[j].SpecialCard != 0 { //! 对方是特殊牌型
					score := self.GetSpecialScore(self.PersonMgr[j].SpecialCard, self.PersonMgr[j].Cards)
					self.PersonMgr[i].CurScore -= score
				} else { //! 双方都不是特殊牌型
					headscore, bodyscore, tailscore := self.Compare(self.PersonMgr[i], self.PersonMgr[j])
					self.PersonMgr[i].HeadScore += headscore
					self.PersonMgr[i].BodyScore += bodyscore
					self.PersonMgr[i].TailScore += tailscore
					self.PersonMgr[i].CurScore += (headscore + bodyscore + tailscore)
					if fan { //! 会出现翻车
						if headscore > 0 && bodyscore > 0 && tailscore > 0 {
							self.PersonMgr[i].FanNum++
							son.GunScore += (headscore + bodyscore + tailscore)
							son.GunUid = append(son.GunUid, self.PersonMgr[j].Uid)
						} else if headscore < 0 && bodyscore < 0 && tailscore < 0 {
							son.BeGunScore += (headscore + bodyscore + tailscore)
						}
					}
				}
			} else { //! 自己是特殊牌型
				if self.PersonMgr[j].SpecialCard == 0 { //! 对方不是特殊牌型
					score := self.GetSpecialScore(self.PersonMgr[i].SpecialCard, self.PersonMgr[i].Cards)
					self.PersonMgr[i].CurScore += score
				} else { //! 都是特殊牌型
					if (self.PersonMgr[i].SpecialCard == 3 || self.PersonMgr[i].SpecialCard == 4) && (self.PersonMgr[j].SpecialCard == 3 || self.PersonMgr[j].SpecialCard == 4) {
						self.PersonMgr[i].CurScore += 0
					} else if self.PersonMgr[i].SpecialCard > self.PersonMgr[j].SpecialCard {
						score := self.GetSpecialScore(self.PersonMgr[i].SpecialCard, self.PersonMgr[i].Cards)
						self.PersonMgr[i].CurScore += score
					} else if self.PersonMgr[i].SpecialCard < self.PersonMgr[j].SpecialCard {
						score := self.GetSpecialScore(self.PersonMgr[j].SpecialCard, self.PersonMgr[j].Cards)
						self.PersonMgr[i].CurScore -= score
					}
				}
			}
		}
		son.HeadScore = self.PersonMgr[i].HeadScore
		son.BodyScore = self.PersonMgr[i].BodyScore
		son.TailScore = self.PersonMgr[i].TailScore
		son.GunNum = self.PersonMgr[i].FanNum
		son.SwatScore = 0

		if red && len(son.GunUid) == len(self.PersonMgr)-1 {
			son.SwatScore = son.HeadScore + son.BodyScore + son.TailScore + son.GunScore

			for z := 0; z < len(self.PersonMgr); z++ {
				if self.PersonMgr[z].Uid == self.PersonMgr[i].Uid {
					continue
				}
				headScore, bodyScore, tailScore := self.Compare(self.PersonMgr[z], self.PersonMgr[i])
				score := (headScore + bodyScore + tailScore) * 2
				self.PersonMgr[z].CurScore += score
			}
		}
		msg.Info = append(msg.Info, son)
	}
	for i := 0; i < len(msg.Info); i++ {
		for j := 0; j < len(self.PersonMgr); j++ {
			if msg.Info[i].Uid == self.PersonMgr[j].Uid {
				self.PersonMgr[j].CurScore += (msg.Info[i].GunScore + msg.Info[i].SwatScore + msg.Info[i].BeGunScore)
				self.PersonMgr[j].Score += self.PersonMgr[j].CurScore
				msg.Info[i].Score = self.PersonMgr[j].CurScore
				msg.Info[i].Total = self.PersonMgr[j].Score

				var rec Son_Rec_GameSSD
				rec.Uid = self.PersonMgr[j].Uid
				rec.Name = self.room.GetName(self.PersonMgr[j].Uid)
				rec.Head = self.room.GetHead(self.PersonMgr[j].Uid)
				rec.Card = self.PersonMgr[j].Cards
				rec.Score = self.PersonMgr[j].CurScore
				rec.Total = self.PersonMgr[j].Score
				rec.Special = self.PersonMgr[j].SpecialCard
				record.Info = append(record.Info, rec)
			}
		}
	}

	self.room.AddRecord(lib.HF_JtoA(&record))
	self.room.broadCastMsg("gamessdend", &msg)

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].CurScore > 0 {
			self.PersonMgr[i].Winner++
		} else if self.PersonMgr[i].CurScore < 0 {
			self.PersonMgr[i].Loser++
		} else {
			self.PersonMgr[i].DogFall++
		}
	}

	if self.room.IsBye() {
		self.OnBye()
		self.room.Bye()
		return
	}

	self.room.flush()
}

func (self *Game_SSD) OnInit(room *Room) { //! 初始化
	self.room = room

}

func (self *Game_SSD) OnSendInfo(person *Person) { //! 告诉玩家数据

	for i := 0; i < len(self.PersonMgr); i++ {
		if person.Uid == self.PersonMgr[i].Uid {
			person.SendMsg("gamessdinfo", self.getInfo(person.Uid))
			return
		}
	}

	self.Init()

	if len(self.PersonMgr)+1 > self.MaxPerson { //! 房间人数已满
		self.OnExit(person.Uid)
		self.room.KickPerson(person.Uid, 97)
		return
	}

	_person := new(Game_SSD_Person)
	_person.Uid = person.Uid
	_person.Ready = false
	self.PersonMgr = append(self.PersonMgr, _person)
	person.SendMsg("gamessdinfo", self.getInfo(person.Uid))

	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Ready {
			num++
		}
	}

	if num == self.MaxPerson { //! 准备的人数达到游戏人数
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
		self.OnBegin()
		return
	}

	self.room.flush()
}

func (self *Game_SSD) OnMsg(msg *RoomMsg) { //! 消息转发
	switch msg.Head {
	case "gameready":
		self.Ready(msg.Uid)
	case "gamematchs":
		self.GameMatch(msg.Uid, msg.V.(*Msg_GameSSDMatch).Cards, msg.V.(*Msg_GameSSDMatch).AbsCards, msg.V.(*Msg_GameSSDMatch).Special)
	/////////////////////////////
	case "gamesteps":
		self.GameSteps(msg.Uid, msg.V.(*Msg_GameSteps).Cards)
	}
}

func (self *Game_SSD) GameSteps(uid int64, card []int) {
	person := self.GetPerson(uid)
	person.SpecialCards = card
	_person := GetPersonMgr().GetPerson(uid)
	_person.SendMsg("successcard", nil)

}

func (self *Game_SSD) OnBye() { //! 游戏结算
	var msg Msg_GameSSD_Bye
	msg.Info = make([]Son_GameSSD_Bye, 0)
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameSSD_Bye
		son.Uid = self.PersonMgr[i].Uid
		son.Score = self.PersonMgr[i].Score
		son.Winner = self.PersonMgr[i].Winner
		son.Loser = self.PersonMgr[i].Loser
		son.DogFall = self.PersonMgr[i].DogFall
		msg.Info = append(msg.Info, son)
	}
	self.room.broadCastMsg("gamessdbye", &msg)
}

func (self *Game_SSD) OnExit(uid int64) { //! 玩家退出
	for i := 0; i < len(self.PersonMgr); i++ {
		if uid == self.PersonMgr[i].Uid {
			copy(self.PersonMgr[i:], self.PersonMgr[i+1:])
			self.PersonMgr = self.PersonMgr[:len(self.PersonMgr)-1]
			return
		}
	}
}

func (self *Game_SSD) OnIsDealer(uid int64) bool { //! 是否是庄家
	return false
}

func (self *Game_SSD) OnTime() { //! 每秒调用一次

}
func (self *Game_SSD) OnIsBets(uid int64) bool {
	return false
}

//! 结算所有人
func (self *Game_SSD) OnBalance() {
}
func (self *Game_SSD) OnRobot(robot *lib.Robot) {

}
