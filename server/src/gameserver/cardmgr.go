//! 扑克牌的逻辑

package gameserver

import (
	"lib"
	"math"
	"math/rand"
	"sort"
	"staticfunc"
	"time"
)

const TYPE_CARD_ONE = 1       //! 单张
const TYPE_CARD_TWO = 2       //! 对子
const TYPE_CARD_WANG = 3      //! 王炸
const TYPE_CARD_ZHA = 4       //! 炸弹
const TYPE_CARD_SAN = 5       //! 三张
const TYPE_CARD_SAN1 = 6      //! 三带一
const TYPE_CARD_SAN2 = 7      //! 三带二
const TYPE_CARD_SHUN = 8      //! 顺子
const TYPE_CARD_SHUNDUI = 9   //! 顺对
const TYPE_CARD_SHUNSAN = 10  //! 顺3
const TYPE_CARD_SHUNSAN1 = 12 //! 顺3带1
const TYPE_CARD_SHUNSAN2 = 13 //! 顺3带2
const TYPE_CARD_SI1 = 14      //! 4带单
const TYPE_CARD_SI2 = 15      //! 4带双
const TYPE_CARD_SI3 = 16      //! 4带三

const TYPE_TSDG_ONE = 1      //! 单张
const TYPE_TSDG_TWO = 2      //! 对子
const TYPE_TSDG_SAN = 3      //! 三张
const TYPE_TSDG_SHUNDUI = 4  //! 连队
const TYPE_TSDG_SHUNSAN = 5  //! 飞机
const TYPE_TSDG_BOOM = 6     //! 炸弹
const TYPE_TSDG_BOOM510K = 7 //! 510K炸弹
const TYPE_TSDG_BOOMKING = 8 //! 王炸

const TYPE_SSS_WULONG = 1      //! 乌龙
const TYPE_SSS_DUIZI = 2       //! 对子
const TYPE_SSS_LIANGDUI = 3    //! 两对
const TYPE_SSS_SANTIAO = 4     //! 三条
const TYPE_SSS_SHUNZI = 5      //! 顺子
const TYPE_SSS_TONGHUA = 6     //! 同花
const TYPE_SSS_HULU = 7        //! 葫芦
const TYPE_SSS_BOOM = 8        //! 炸弹
const TYPE_SSS_TONGHUASHUN = 9 //! 同花顺

const TYPE_CARD_BZ = 1         //! 一豹
const TYPE_CARD_SWANG = 2      //! 双王
const TYPE_CARD_SLIU = 3       //! 双六
const TYPE_CARD_SSHUN = 4      //! 四顺
const TYPE_CARD_SSHUNGUAI = 5  //! 三顺加拐
const TYPE_CARD_WUSHUN = 6     //! 五顺
const TYPE_CARD_SB = 7         //! 双豹
const TYPE_CARD_SISHUNGUAI = 8 //! 四顺加拐
const TYPE_CARD_QG = 9         //! 桥拐
const TYPE_CARD_SISHUNBAO = 10 //! 四顺加豹或三顺
const TYPE_CARD_LIUSHUN = 11   //! 六顺
const TYPE_CARD_WSG = 12       //! 五顺加拐
const TYPE_CARD_SANBAO = 13    //! 三豹
const TYPE_CARD_BGB = 14       //! 豹拐加豹
const TYPE_CARD_QGQ = 15       //! 桥拐加桥
const TYPE_CARD_SGQ = 16       //! 四顺加拐加桥
const TYPE_CARD_SHUNBAO = 17   //! 五顺加豹
const TYPE_CARD_SHUNGUAI = 18  //! 六顺加拐
const TYPE_CARD_QISHUN = 19    //! 七顺
const TYPE_CARD_ZHADAN = 20    //! 炸弹

type LstPoker []int

func (a LstPoker) Len() int { // 重写 Len() 方法
	return len(a)
}

func (a LstPoker) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}

func (a LstPoker) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return CardCompare(a[j], a[i]) > 0
}

//! 牌型=A-K为1-13
//! 颜色=桃-方4-1
//! 牌=牌型*10+颜色
//! 小王=1000,大王=2000
type CardMgr struct {
	Card   []int //! 剩余牌组
	random *rand.Rand
}

//! 鹤峰百胡
func NewCard_HFBH() *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 5; i++ {
		for j := 1; j <= 8; j++ {
			for k := 1; k <= 3; k++ {
				mgr.Card = append(mgr.Card, j*10+k)
			}
		}
	}

	return mgr
}

//! 百家乐
func NewCard_BJL156() *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	for j := 0; j < 8; j++ {
		for i := 1; i <= 13; i++ {
			mgr.Card = append(mgr.Card, i*10+4) //! 桃
			mgr.Card = append(mgr.Card, i*10+3) //! 心
			mgr.Card = append(mgr.Card, i*10+2) //! 梅
			mgr.Card = append(mgr.Card, i*10+1) //! 方
		}
	}
	return mgr
}

//! 翻牌机
func NewCard_FPJ(king bool) *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 13; i++ {
		mgr.Card = append(mgr.Card, i*10+4) //! 桃
		mgr.Card = append(mgr.Card, i*10+3) //! 心
		mgr.Card = append(mgr.Card, i*10+2) //! 梅
		mgr.Card = append(mgr.Card, i*10+1) //! 方
	}
	if king {
		mgr.Card = append(mgr.Card, 1000)
		mgr.Card = append(mgr.Card, 2000)
	}

	return mgr
}

//！捞腌菜牌组
func NewCard_LYC() *CardMgr {
	mgr := new(CardMgr)

	for i := 1; i <= 13; i++ {
		mgr.Card = append(mgr.Card, i*10+4) //! 桃
		mgr.Card = append(mgr.Card, i*10+3) //! 心
		mgr.Card = append(mgr.Card, i*10+2) //! 梅
		mgr.Card = append(mgr.Card, i*10+1) //! 方
	}

	return mgr
}

//! 牛牛牌组
func NewCard_NiuNiu(razz bool) *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 13; i++ {
		mgr.Card = append(mgr.Card, i*10+4) //! 桃
		mgr.Card = append(mgr.Card, i*10+3) //! 心
		mgr.Card = append(mgr.Card, i*10+2) //! 梅
		mgr.Card = append(mgr.Card, i*10+1) //! 方
	}

	if razz {
		//mgr.Card = append(mgr.Card, 1000)
		mgr.Card = append(mgr.Card, 2000)
	}

	return mgr
}

//! 斗地主牌组
func NewCard_ZGL() *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 10; i++ {
		mgr.Card = append(mgr.Card, i*10+4) //! 桃
		mgr.Card = append(mgr.Card, i*10+3) //! 心
		mgr.Card = append(mgr.Card, i*10+2) //! 梅
		mgr.Card = append(mgr.Card, i*10+1) //! 方
	}
	return mgr
}

//! 斗地主牌组
func NewCard_DDZ() *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 13; i++ {
		mgr.Card = append(mgr.Card, i*10+4) //! 桃
		mgr.Card = append(mgr.Card, i*10+3) //! 心
		mgr.Card = append(mgr.Card, i*10+2) //! 梅
		mgr.Card = append(mgr.Card, i*10+1) //! 方
	}
	mgr.Card = append(mgr.Card, 1000)
	mgr.Card = append(mgr.Card, 2000)

	return mgr
}

//! 跑得快金币场
func NewCard_GoldPDK() *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 13; i++ {
		if i == 1 {
			mgr.Card = append(mgr.Card, i*10+3) //! 心
			mgr.Card = append(mgr.Card, i*10+2) //! 梅
			mgr.Card = append(mgr.Card, i*10+1) //! 方
		} else if i == 2 {
			mgr.Card = append(mgr.Card, i*10+1) //! 方
		} else {
			mgr.Card = append(mgr.Card, i*10+4) //! 桃
			mgr.Card = append(mgr.Card, i*10+3) //! 心
			mgr.Card = append(mgr.Card, i*10+2) //! 梅
			mgr.Card = append(mgr.Card, i*10+1) //! 方
		}
	}

	return mgr
}

//! 斗板凳牌组
func NewCard_DBD() *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 13; i++ {
		if i == 2 {
			continue
		}
		mgr.Card = append(mgr.Card, i*10+4) //! 桃
		mgr.Card = append(mgr.Card, i*10+3) //! 心
		mgr.Card = append(mgr.Card, i*10+2) //! 梅
		mgr.Card = append(mgr.Card, i*10+1) //! 方
	}
	mgr.Card = append(mgr.Card, 1000)
	mgr.Card = append(mgr.Card, 2000)
	mgr.Card = append(mgr.Card, 3000) //！ 花牌
	mgr.Card = append(mgr.Card, 21)
	mgr.Card = append(mgr.Card, 22)

	return mgr

}

//! 梭哈28张牌
func NewCard_SuoHa28() *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 8; i <= 13; i++ {
		mgr.Card = append(mgr.Card, i*10+4) //! 桃
		mgr.Card = append(mgr.Card, i*10+3) //! 心
		mgr.Card = append(mgr.Card, i*10+2) //! 梅
		mgr.Card = append(mgr.Card, i*10+1) //! 方
	}
	mgr.Card = append(mgr.Card, 1*10+4)
	mgr.Card = append(mgr.Card, 1*10+3)
	mgr.Card = append(mgr.Card, 1*10+2)
	mgr.Card = append(mgr.Card, 1*10+1)
	return mgr
}

//! 十三道牌组
func NewCard_SSD(wang int, tmp int) *CardMgr {
	//0 不去牌 	 1 去2-4  	2 去2-6
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	mgr.Card = append(mgr.Card, 1*10+4) //! 桃
	mgr.Card = append(mgr.Card, 1*10+3) //! 心
	mgr.Card = append(mgr.Card, 1*10+2) //! 梅
	mgr.Card = append(mgr.Card, 1*10+1) //! 方

	if tmp == 0 {
		for i := 2; i <= 13; i++ {
			mgr.Card = append(mgr.Card, i*10+4) //! 桃
			mgr.Card = append(mgr.Card, i*10+3) //! 心
			mgr.Card = append(mgr.Card, i*10+2) //! 梅
			mgr.Card = append(mgr.Card, i*10+1) //! 方
		}
	} else if tmp == 1 {
		for i := 5; i <= 13; i++ {
			mgr.Card = append(mgr.Card, i*10+4) //! 桃
			mgr.Card = append(mgr.Card, i*10+3) //! 心
			mgr.Card = append(mgr.Card, i*10+2) //! 梅
			mgr.Card = append(mgr.Card, i*10+1) //! 方
		}
	} else if tmp == 2 {
		for i := 7; i <= 13; i++ {
			mgr.Card = append(mgr.Card, i*10+4) //! 桃
			mgr.Card = append(mgr.Card, i*10+3) //! 心
			mgr.Card = append(mgr.Card, i*10+2) //! 梅
			mgr.Card = append(mgr.Card, i*10+1) //! 方
		}
	}

	if wang == 1 { //带王
		mgr.Card = append(mgr.Card, 1000)
		mgr.Card = append(mgr.Card, 2000)
	}
	return mgr
}

//! 通山打拱牌组
func NewCard_TSDG(double bool) *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 13; i++ {
		mgr.Card = append(mgr.Card, i*10+4) //! 桃
		mgr.Card = append(mgr.Card, i*10+4) //! 桃
		mgr.Card = append(mgr.Card, i*10+3) //! 心
		mgr.Card = append(mgr.Card, i*10+3) //! 心
		mgr.Card = append(mgr.Card, i*10+2) //! 梅
		mgr.Card = append(mgr.Card, i*10+2) //! 梅
		mgr.Card = append(mgr.Card, i*10+1) //! 方
		mgr.Card = append(mgr.Card, i*10+1) //! 方
	}
	mgr.Card = append(mgr.Card, 1000)
	mgr.Card = append(mgr.Card, 2000)
	mgr.Card = append(mgr.Card, 1000)
	mgr.Card = append(mgr.Card, 2000)
	if double {
		mgr.Card = append(mgr.Card, 1000)
		mgr.Card = append(mgr.Card, 2000)
		mgr.Card = append(mgr.Card, 1000)
		mgr.Card = append(mgr.Card, 2000)
	}

	return mgr
}

//! 通山打拱牌组
func NewCard_BJL(amount int) *CardMgr { //！百家乐牌组
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 13; i++ {
		for j := 0; j < amount; j++ {
			mgr.Card = append(mgr.Card, i*10+4) //! 桃
			mgr.Card = append(mgr.Card, i*10+3) //! 心
			mgr.Card = append(mgr.Card, i*10+2) //! 梅
			mgr.Card = append(mgr.Card, i*10+1) //! 方
		}
	}
	return mgr
}

//! 跑的快经典无赖子牌组
func NewCard_NewPDK(tmp int) *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 3; i <= 12; i++ {
		mgr.Card = append(mgr.Card, i*10+4) //! 桃
		mgr.Card = append(mgr.Card, i*10+3) //! 心
		mgr.Card = append(mgr.Card, i*10+2) //! 梅
		mgr.Card = append(mgr.Card, i*10+1) //! 方
	}
	if tmp == 1 { //! 经典玩法
		mgr.Card = append(mgr.Card, 24)
		mgr.Card = append(mgr.Card, 14)
		mgr.Card = append(mgr.Card, 13)
		mgr.Card = append(mgr.Card, 12)
		mgr.Card = append(mgr.Card, 134)
		mgr.Card = append(mgr.Card, 133)
		mgr.Card = append(mgr.Card, 132)
		mgr.Card = append(mgr.Card, 131)
	} else if tmp == 2 { //! 15张玩法
		mgr.Card = append(mgr.Card, 24)
		mgr.Card = append(mgr.Card, 14)
		mgr.Card = append(mgr.Card, 134)
		mgr.Card = append(mgr.Card, 133)
		mgr.Card = append(mgr.Card, 132)
	}

	return mgr
}

//! 跑得快48牌组
func NewCard_Run48() *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 3; i <= 13; i++ {
		mgr.Card = append(mgr.Card, i*10+4) //! 桃
		mgr.Card = append(mgr.Card, i*10+3) //! 心
		mgr.Card = append(mgr.Card, i*10+2) //! 梅
		mgr.Card = append(mgr.Card, i*10+1) //! 方
	}
	mgr.Card = append(mgr.Card, 24)
	mgr.Card = append(mgr.Card, 14)
	mgr.Card = append(mgr.Card, 13)
	mgr.Card = append(mgr.Card, 12)

	return mgr
}

//! 跑得快去掉大小王
func NewCard_Run52() *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 13; i++ {
		mgr.Card = append(mgr.Card, i*10+4) //! 桃
		mgr.Card = append(mgr.Card, i*10+3) //! 心
		mgr.Card = append(mgr.Card, i*10+2) //! 梅
		mgr.Card = append(mgr.Card, i*10+1) //! 方
	}

	return mgr
}

//!2副牌组
func NewCard_Run108() *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	for j := 0; j <= 1; j++ {
		for i := 1; i <= 13; i++ {
			mgr.Card = append(mgr.Card, i*10+4) //! 桃
			mgr.Card = append(mgr.Card, i*10+3) //! 心
			mgr.Card = append(mgr.Card, i*10+2) //! 梅
			mgr.Card = append(mgr.Card, i*10+1) //! 方
		}
		mgr.Card = append(mgr.Card, 1000)
		mgr.Card = append(mgr.Card, 2000)
	}
	return mgr
}

//! 跑得快45牌组
func NewCard_Run45() *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 3; i <= 12; i++ {
		mgr.Card = append(mgr.Card, i*10+4) //! 桃
		mgr.Card = append(mgr.Card, i*10+3) //! 心
		mgr.Card = append(mgr.Card, i*10+2) //! 梅
		mgr.Card = append(mgr.Card, i*10+1) //! 方
	}
	mgr.Card = append(mgr.Card, 14)
	mgr.Card = append(mgr.Card, 24)
	mgr.Card = append(mgr.Card, 134)
	mgr.Card = append(mgr.Card, 133)
	mgr.Card = append(mgr.Card, 132)

	return mgr
}

//! 炸金花牌组
func NewCard_ZJH() *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 13; i++ {
		mgr.Card = append(mgr.Card, i*10+4) //! 桃
		mgr.Card = append(mgr.Card, i*10+3) //! 心
		mgr.Card = append(mgr.Card, i*10+2) //! 梅
		mgr.Card = append(mgr.Card, i*10+1) //! 方
	}

	return mgr
}

//! 天九牌组
func NewCard_TJ() *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 21; i++ {
		mgr.Card = append(mgr.Card, i)
		if i <= 11 {
			mgr.Card = append(mgr.Card, i)
		}
	}

	return mgr
}

//! 内蒙古帕斯牌组(36张牌)
func NewCard_PS() *CardMgr {
	mgr := new(CardMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 7; i <= 14; i++ {
		mgr.Card = append(mgr.Card, i*10+4) //! 桃
		mgr.Card = append(mgr.Card, i*10+3) //! 心
		mgr.Card = append(mgr.Card, i*10+2) //! 梅
		mgr.Card = append(mgr.Card, i*10+1) //! 方
	}
	mgr.Card = append(mgr.Card, 1000)
	mgr.Card = append(mgr.Card, 2000)
	mgr.Card = append(mgr.Card, 64)
	mgr.Card = append(mgr.Card, 61)

	return mgr
}

//! 发牌
func (self *CardMgr) Deal(num int) []int {
	if self.random == nil {
		self.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	lst := make([]int, 0)

	for i := 0; i < num; i++ {
		index := self.random.Intn(len(self.Card))

		lst = append(lst, self.Card[index])

		copy(self.Card[index:], self.Card[index+1:])
		self.Card = self.Card[:len(self.Card)-1]
	}

	return lst
}

//! 去掉一张牌
func (self *CardMgr) Del(card int) {
	for i := 0; i < len(self.Card); i++ {
		if self.Card[i] == card {
			copy(self.Card[i:], self.Card[i+1:])
			self.Card = self.Card[:len(self.Card)-1]
			break
		}
	}
}

//! 发指定的牌
func (self *CardMgr) DealCard(card int) int {
	for i := 0; i < len(self.Card); i++ {
		if self.Card[i] == card {
			copy(self.Card[i:], self.Card[i+1:])
			self.Card = self.Card[:len(self.Card)-1]
			return card
		}
	}

	return 0
}

////////////////////////////////////////////////////
func GetNiuNiuByRazz(card []int) (int, int) {
	var _card []int
	lib.HF_DeepCopy(&_card, &card)

	//! 去掉癞子
	razzcard := 0
	razznum := 0
	for i := 0; i < len(_card); {
		if _card[i] == 1000 || _card[i] == 2000 {
			if _card[i] > razzcard {
				razzcard = _card[i]
			}
			copy(_card[i:], _card[i+1:])
			_card = _card[:len(_card)-1]
			razznum++
		} else {
			i++
		}
	}

	if razznum == 0 {
		return GetNiuNiuJXScore(_card, 1, 1, 1)
	} else if razznum == 1 {
		maxct, maxcs := 0, 0
		for i := 13; i >= 1; i-- {
			var cardrazz []int
			lib.HF_DeepCopy(&cardrazz, &_card)
			cardrazz = append(cardrazz, i*10+5)
			ct, cs := GetNiuNiuJXScore(cardrazz, 1, 1, 1)
			if ct > maxct {
				maxct = ct
				maxcs = cs
			} else if ct == maxct && cs > maxcs {
				maxct = ct
				maxcs = cs
			}
		}
		//if razzcard > maxcs {
		//	maxcs = razzcard
		//}
		return maxct, maxcs
	} else {
		maxct, maxcs := 0, 0
		for i := 13; i >= 1; i-- {
			for j := 13; j >= 1; j-- {
				var cardrazz []int
				lib.HF_DeepCopy(&cardrazz, &_card)
				cardrazz = append(cardrazz, i*10+5)
				cardrazz = append(cardrazz, j*10+5)
				ct, cs := GetNiuNiuJXScore(cardrazz, 1, 1, 1)
				if ct > maxct {
					maxct = ct
					maxcs = cs
				} else if ct == maxct && cs > maxcs {
					maxct = ct
					maxcs = cs
				}
			}
		}
		//if razzcard > maxcs {
		//	maxcs = razzcard
		//}
		return maxct, maxcs
	}
}

func GetNiuNiuJXBS(ct int) int {
	if ct == 400 {
		return 8
	} else if ct == 300 {
		return 6
	} else if ct == 200 {
		return 5
	} else if ct == 100 {
		return 4
	} else if ct == 99 {
		return 3
	} else if ct == 98 {
		return 2
	}

	return 1
}

func GetMPQZBS(ct int, _type int) int {
	if ct == 800 { //! 顺金牛
		return 10
	} else if ct == 700 { //! 炸弹牛
		return 8
	} else if ct == 600 { //! 葫芦牛
		return 7
	} else if ct == 500 { //! 五小牛
		return 6
	} else if ct == 400 { //! 同花牛
		return 6
	} else if ct == 300 { //! 五花牛
		return 5
	} else if ct == 200 { //! 顺子牛
		return 5
	}

	if _type == 0 {
		if ct == 100 {
			return 3
		} else if ct >= 98 {
			return 2
		}
	} else {
		if ct == 100 {
			return 4
		} else if ct == 99 {
			return 3
		} else if ct >= 97 {
			return 2
		}
	}

	return 1
}

func GetGoldNiuNiuBS(ct int) int {
	if ct == 800 { //! 顺金牛
		return 10
	} else if ct == 700 { //! 炸弹牛
		return 8
	} else if ct == 600 { //! 葫芦牛
		return 7
	} else if ct == 500 { //! 五小牛
		return 6
	} else if ct == 400 { //! 同花牛
		return 6
	} else if ct == 300 { //! 五花牛
		return 5
	} else if ct == 200 { //! 顺子牛
		return 5
	} else if ct == 100 {
		return 4
	} else if ct == 99 {
		return 3
	} else if ct >= 97 {
		return 2
	}

	return 1
}

func GetGoldBrNNBS(ct int, tmp int) int { //0-低倍 1-高倍
	if tmp == 0 {
		if ct == 300 { //! 五花牛
			return 5
		} else if ct == 200 { //! 四炸
			return 4
		} else if ct == 100 { //! 牛牛
			return 3
		} else if ct >= 97 { //! 牛七-牛九
			return 2
		}
		return 1
	} else {
		if ct >= 100 { //! 牛牛,四炸,
			return 10
		} else if ct == 99 {
			return 9
		} else if ct == 98 {
			return 8
		} else if ct == 97 {
			return 7
		} else if ct == 96 {
			return 6
		} else if ct == 95 {
			return 5
		} else if ct == 94 {
			return 4
		} else if ct == 93 {
			return 3
		} else if ct == 92 {
			return 2
		}
		return 1
	}
}

func GetBrNiuNiuScore(card []int) (int, int) {
	//! 先取最大的牌
	maxcard := 0
	for i := 0; i < len(card); i++ {
		if card[i] > maxcard {
			maxcard = card[i]
		}
	}

	//! 判断五花牛
	{
		tmp := 0
		for i := 0; i < len(card); i++ {
			if card[i]/10 > 10 {
				tmp++
			}
		}
		if tmp >= 5 {
			return 300, maxcard
		}
	}

	//! 判断炸弹
	{
		tmp := make(map[int]int)
		for i := 0; i < len(card); i++ {
			tmp[card[i]/10] += 1
		}
		for _, value := range tmp {
			if value >= 4 {
				return 200, maxcard
			}
		}
	}

	isniu, niu := GetNiuNiuType(card)
	if isniu {
		//! 牛牛
		if niu >= 10 {
			return 100, maxcard
		} else {
			return 90 + niu, maxcard
		}
	}

	return 0, maxcard
}

//! 金币场牛元帅
func GetGoldNiuNiuScore(card []int, sjn bool, zdn bool, hln bool, wxn bool, thn bool, whn bool, szn bool) (int, int) {
	//! 先取最大的牌
	maxcard := 0
	for i := 0; i < len(card); i++ {
		if card[i] > maxcard {
			maxcard = card[i]
		}
	}

	//! 判断是否是顺金牛
	if sjn {
		istrue := true
		_card := make([]int, 0)
		color := card[0] % 10
		a := 0
		hasa := false
		for i := 0; i < len(card); i++ {
			if card[i]%10 != color {
				istrue = false
				break
			}
			if card[i]/10 == 1 { //! 有A
				hasa = true
			} else {
				if card[i]/10 == 2 {
					a = 1
				} else if card[i]/10 == 13 {
					a = 14
				}
				_card = append(_card, card[i]/10)
			}
		}
		if istrue {
			if hasa {
				_card = append(_card, a)
			}
			if len(_card) == 5 {
				sort.Ints(_card)
				for i := 1; i < len(_card); i++ {
					if math.Abs(float64(_card[i]-_card[i-1])) != 1 {
						istrue = false
						break
					}
				}
			} else {
				istrue = false
			}
		}
		if istrue {
			lib.GetLogMgr().Output(lib.LOG_ERROR, "顺金牛", card)
			return 800, maxcard
		}
	}

	//! 判断炸弹
	if zdn {
		tmp := make(map[int]int)
		for i := 0; i < len(card); i++ {
			tmp[card[i]/10] += 1
		}
		for _, value := range tmp {
			if value >= 4 {
				lib.GetLogMgr().Output(lib.LOG_ERROR, "炸弹牛", card)
				return 700, maxcard
			}
		}
	}

	//! 判断是否葫芦牛
	if hln {
		tmp := make(map[int]int)
		for i := 0; i < len(card); i++ {
			tmp[card[i]/10] += 1
		}
		istrue := true
		for _, value := range tmp {
			if value < 2 || value > 3 {
				istrue = false
				break
			}
		}
		if istrue {
			lib.GetLogMgr().Output(lib.LOG_ERROR, "葫芦牛", card)
			return 600, maxcard
		}
	}

	//! 判断是否是五小牛
	if wxn {
		tmp := 0
		for i := 0; i < len(card); i++ {
			if card[i]/10 >= 5 {
				tmp = -1
				break
			}
			tmp += card[i] / 10
		}
		if tmp != -1 && tmp <= 10 {
			lib.GetLogMgr().Output(lib.LOG_ERROR, "五小牛", card)
			return 500, maxcard
		}
	}

	//! 同花牛400
	if thn {
		color := card[0] % 10
		istrue := true
		for i := 0; i < len(card); i++ {
			if card[i]%10 != color {
				istrue = false
				break
			}
		}
		if istrue {
			lib.GetLogMgr().Output(lib.LOG_ERROR, "同花牛", card)
			return 400, maxcard
		}
	}

	//! 判断五花牛
	if whn {
		tmp := 0
		for i := 0; i < len(card); i++ {
			if card[i]/10 > 10 {
				tmp++
			}
		}
		if tmp >= 5 {
			lib.GetLogMgr().Output(lib.LOG_ERROR, "五花牛", card)
			return 300, maxcard
		}
	}

	//! 顺子牛200
	if szn {
		istrue := true
		_card := make([]int, 0)
		a := 0
		hasa := false
		for i := 0; i < len(card); i++ {
			if card[i]/10 == 1 { //! 有A
				hasa = true
			} else {
				if card[i]/10 == 2 {
					a = 1
				} else if card[i]/10 == 13 {
					a = 14
				}
				_card = append(_card, card[i]/10)
			}
		}
		if hasa {
			_card = append(_card, a)
		}
		if len(_card) == 5 {
			sort.Ints(_card)
			for i := 1; i < len(_card); i++ {
				if math.Abs(float64(_card[i]-_card[i-1])) != 1 {
					istrue = false
					break
				}
			}
			if istrue {
				lib.GetLogMgr().Output(lib.LOG_ERROR, "顺子牛", card)
				return 200, maxcard
			}
		}
	}

	isniu, niu := GetNiuNiuType(card)
	if isniu {
		//! 牛牛
		if niu >= 10 {
			return 100, maxcard
		} else {
			return 90 + niu, maxcard
		}
	}

	return 0, maxcard
}

func GetNiuNiuJXScore(card []int, wx int, zd int, wh int) (int, int) {
	//! 先取最大的牌
	maxcard := 0
	for i := 0; i < len(card); i++ {
		if card[i] > maxcard {
			maxcard = card[i]
		}
	}

	//! 判断是否是五小牛
	if wx != 0 {
		tmp := 0
		for i := 0; i < len(card); i++ {
			if card[i]/10 >= 5 {
				tmp = -1
				break
			}
			tmp += card[i] / 10
		}
		if tmp != -1 && tmp <= 10 {
			return 400, maxcard
		}
	}

	//! 判断炸弹
	if zd != 0 {
		tmp := make(map[int]int)
		for i := 0; i < len(card); i++ {
			tmp[card[i]/10] += 1
		}
		for _, value := range tmp {
			if value >= 4 {
				return 300, maxcard
			}
		}
	}

	//! 判断五花牛
	if wh != 0 {
		tmp := 0
		for i := 0; i < len(card); i++ {
			if card[i]/10 > 10 {
				tmp++
			}
		}
		if tmp >= 5 {
			return 200, maxcard
		}
	}

	isniu, niu := GetNiuNiuType(card)
	if isniu {
		//! 牛牛
		if niu >= 10 {
			return 100, maxcard
		} else {
			return 90 + niu, maxcard
		}
	}

	return 0, maxcard
}

func GetNiuNiuBS(ct int) int {
	if ct == 400 {
		return 6
	} else if ct == 300 {
		return 5
	} else if ct == 200 {
		return 4
	} else if ct == 100 {
		return 3
	} else if ct >= 97 {
		return 2
	}

	return 1
}

func GetNiuNiuScore(card []int) (int, int) {
	//! 先取最大的牌
	maxcard := 0
	for i := 0; i < len(card); i++ {
		if card[i] > maxcard {
			maxcard = card[i]
		}
	}

	//! 判断是否是五小牛
	{
		tmp := 0
		for i := 0; i < len(card); i++ {
			if card[i]/10 >= 5 {
				tmp = -1
				break
			}
			tmp += card[i] / 10
		}
		if tmp != -1 && tmp <= 10 {
			return 400, maxcard
		}
	}

	//! 判断五花牛
	{
		tmp := 0
		for i := 0; i < len(card); i++ {
			if card[i]/10 > 10 {
				tmp++
			}
		}
		if tmp >= 5 {
			return 300, maxcard
		}
	}

	//! 判断炸弹
	{
		tmp := make(map[int]int)
		for i := 0; i < len(card); i++ {
			tmp[card[i]/10] += 1
		}
		for _, value := range tmp {
			if value >= 4 {
				return 200, maxcard
			}
		}
	}

	isniu, niu := GetNiuNiuType(card)
	if isniu {
		//! 牛牛
		if niu >= 10 {
			return 100, maxcard
		} else {
			return 90 + niu, maxcard
		}
	}

	return 0, maxcard
}

func GetNiuNiuType(_card []int) (bool, int) {
	niu := 0
	card := make([]int, 0)
	for i := 0; i < len(_card); i++ {
		card = append(card, lib.HF_MinInt(_card[i]/10, 10))
	}

	tmp := card[0] + card[1] + card[2]
	if tmp%10 == 0 {
		_tmp := (card[3] + card[4]) % 10
		if _tmp == 0 { //! 已经是牛牛
			return true, 10
		} else if _tmp > niu {
			niu = _tmp
		}
	}

	tmp = card[0] + card[1] + card[3]
	if tmp%10 == 0 {
		_tmp := (card[2] + card[4]) % 10
		if _tmp == 0 { //! 已经是牛牛
			return true, 10
		} else if _tmp > niu {
			niu = _tmp
		}
	}

	tmp = card[0] + card[1] + card[4]
	if tmp%10 == 0 {
		_tmp := (card[2] + card[3]) % 10
		if _tmp == 0 { //! 已经是牛牛
			return true, 10
		} else if _tmp > niu {
			niu = _tmp
		}
	}

	tmp = card[0] + card[2] + card[3]
	if tmp%10 == 0 {
		_tmp := (card[1] + card[4]) % 10
		if _tmp == 0 { //! 已经是牛牛
			return true, 10
		} else if _tmp > niu {
			niu = _tmp
		}
	}

	tmp = card[0] + card[2] + card[4]
	if tmp%10 == 0 {
		_tmp := (card[1] + card[3]) % 10
		if _tmp == 0 { //! 已经是牛牛
			return true, 10
		} else if _tmp > niu {
			niu = _tmp
		}
	}

	tmp = card[0] + card[3] + card[4]
	if tmp%10 == 0 {
		_tmp := (card[1] + card[2]) % 10
		if _tmp == 0 { //! 已经是牛牛
			return true, 10
		} else if _tmp > niu {
			niu = _tmp
		}
	}

	tmp = card[1] + card[2] + card[3]
	if tmp%10 == 0 {
		_tmp := (card[0] + card[4]) % 10
		if _tmp == 0 { //! 已经是牛牛
			return true, 10
		} else if _tmp > niu {
			niu = _tmp
		}
	}

	tmp = card[1] + card[2] + card[4]
	if tmp%10 == 0 {
		_tmp := (card[0] + card[3]) % 10
		if _tmp == 0 { //! 已经是牛牛
			return true, 10
		} else if _tmp > niu {
			niu = _tmp
		}
	}

	tmp = card[1] + card[3] + card[4]
	if tmp%10 == 0 {
		_tmp := (card[0] + card[2]) % 10
		if _tmp == 0 { //! 已经是牛牛
			return true, 10
		} else if _tmp > niu {
			niu = _tmp
		}
	}

	tmp = card[2] + card[3] + card[4]
	if tmp%10 == 0 {
		_tmp := (card[0] + card[1]) % 10
		if _tmp == 0 { //! 已经是牛牛
			return true, 10
		} else if _tmp > niu {
			niu = _tmp
		}
	}

	if niu == 0 {
		return false, 0
	}

	return true, niu
}

func GetZjhType(_card []int) (int, []int) {

	sortcard := make([]int, 0)
	for i := 0; i < len(_card); i++ {
		card := _card[i]
		if card/10 == 1 {
			card = 140 + card%10
		}
		sortcard = append(sortcard, card)
	}

	if len(_card) != 3 {
		return 0, sortcard
	}

	sort.Sort((LstCard(sortcard)))

	/*for i := 0; i < 2; i++ {
		for j := i + 1; j < 3; j++ {
			if sortcard[i] < sortcard[j] {
				temp := sortcard[i]
				sortcard[i] = sortcard[j]
				sortcard[j] = temp
			}
		}
	}*/

	if sortcard[0]/10 == sortcard[1]/10 && sortcard[1]/10 == sortcard[2]/10 {
		return 600, sortcard
	}

	shunzi := false
	jinhua := false

	if sortcard[0]/10 == 14 && sortcard[1]/10 == 3 && sortcard[2]/10 == 2 {
		shunzi = true

		temp := make([]int, 0)
		temp = append(temp, sortcard[1])
		temp = append(temp, sortcard[2])
		temp = append(temp, sortcard[0]%10+10)

		sortcard = temp
	}

	if sortcard[0]/10 == sortcard[1]/10+1 && sortcard[1]/10 == sortcard[2]/10+1 {
		shunzi = true
	}

	if sortcard[0]%10 == sortcard[1]%10 && sortcard[1]%10 == sortcard[2]%10 {
		jinhua = true
	}

	if shunzi && jinhua {
		return 500, sortcard
	}

	if jinhua {
		return 400, sortcard
	}

	if shunzi {
		return 300, sortcard
	}

	if sortcard[0]/10 == sortcard[1]/10 || sortcard[1]/10 == sortcard[2]/10 || sortcard[0]/10 == sortcard[2]/10 {

		point := 0
		if sortcard[0]/10 == sortcard[1]/10 {
			point = sortcard[0] / 10
		} else if sortcard[1]/10 == sortcard[2]/10 {
			point = sortcard[1] / 10
		} else if sortcard[0]/10 == sortcard[2]/10 {
			point = sortcard[0] / 10
		}

		return 200 + point, sortcard
	}

	return 100, sortcard
}

func GetZjhType1(_card []int, tmp int) (int, []int) {

	sortcard := make([]int, 0)
	for i := 0; i < len(_card); i++ {
		card := _card[i]
		if card/10 == 1 {
			card = 140 + card%10
		}
		sortcard = append(sortcard, card)
	}

	if len(_card) != 3 {
		return 0, sortcard
	}

	sort.Sort((LstCard(sortcard)))

	/*for i := 0; i < 2; i++ {
		for j := i + 1; j < 3; j++ {
			if sortcard[i] < sortcard[j] {
				temp := sortcard[i]
				sortcard[i] = sortcard[j]
				sortcard[j] = temp
			}
		}
	}*/

	if sortcard[0]/10 == sortcard[1]/10 && sortcard[1]/10 == sortcard[2]/10 {
		return 600, sortcard
	}

	shunzi := false
	jinhua := false

	if sortcard[0]/10 == 14 && sortcard[1]/10 == 3 && sortcard[2]/10 == 2 {
		shunzi = true

		temp := make([]int, 0)
		temp = append(temp, sortcard[1])
		temp = append(temp, sortcard[2])
		temp = append(temp, sortcard[0]%10+10)

		sortcard = temp
	}

	if sortcard[0]/10 == sortcard[1]/10+1 && sortcard[1]/10 == sortcard[2]/10+1 {
		shunzi = true
	}

	if sortcard[0]%10 == sortcard[1]%10 && sortcard[1]%10 == sortcard[2]%10 {
		jinhua = true
	}

	if shunzi && jinhua {
		return 500, sortcard
	}

	if jinhua && tmp == 0 {
		return 400, sortcard
	} else if jinhua && tmp == 1 {
		return 300, sortcard
	}

	if shunzi && tmp == 0 {
		return 300, sortcard
	} else if shunzi && tmp == 1 {
		return 400, sortcard
	}

	if sortcard[0]/10 == sortcard[1]/10 || sortcard[1]/10 == sortcard[2]/10 || sortcard[0]/10 == sortcard[2]/10 {

		point := 0
		if sortcard[0]/10 == sortcard[1]/10 {
			point = sortcard[0] / 10
		} else if sortcard[1]/10 == sortcard[2]/10 {
			point = sortcard[1] / 10
		} else if sortcard[0]/10 == sortcard[2]/10 {
			point = sortcard[0] / 10
		}

		return 200 + point, sortcard
	}

	return 100, sortcard
}

func ZjhIs235(cardlsttype int, cardlst []int) bool {
	if cardlsttype == 100 && cardlst[0]/10 == 5 && cardlst[1]/10 == 3 && cardlst[2]/10 == 2 {
		return true
	}

	return false
}

//! 比牌 0胜 1负 2平  //比大小
func ZjhCardCompare(_card []int, _card2 []int) int {
	cardlsttype, cardlst := GetZjhType(_card)
	cardlsttype2, cardlst2 := GetZjhType(_card2)

	//log.Panicln(cardlsttype, cardlst, cardlsttype2, cardlst2)

	for i := 0; i < 3; i++ {
		cardpoint1 := cardlst[i] / 10
		cardpoint2 := cardlst2[i] / 10

		if cardpoint1 == 1 {
			cardpoint1 = 14
		}

		if cardpoint2 == 1 {
			cardpoint2 = 14
		}

		cardlst[i] = cardpoint1*10 + cardlst[i]%10
		cardlst2[i] = cardpoint2*10 + cardlst2[i]%10
	}

	sort.Sort(LstCard(cardlst))
	sort.Sort(LstCard(cardlst2))

	/*if cardlsttype == 600 && ZjhIs235(cardlsttype2, cardlst2) {
		return 1
	}

	if cardlsttype2 == 600 && ZjhIs235(cardlsttype, cardlst) {
		return 0
	}*/

	if cardlsttype > cardlsttype2 {
		return 0
	} else if cardlsttype < cardlsttype2 {
		return 1
	}

	for i := 0; i < 3; i++ {
		cardpoint1 := cardlst[i] / 10
		cardpoint2 := cardlst2[i] / 10

		if cardpoint1 == cardpoint2 {
			continue
		}
		if cardpoint1 > cardpoint2 {
			return 0
		} else {
			return 1
		}
	}

	return 2
}

func ZjhCardCompare3(_card []int, _card2 []int, tmp int, tmp2 int) int {
	cardlsttype, cardlst := GetZjhType1(_card, tmp2)
	cardlsttype2, cardlst2 := GetZjhType1(_card2, tmp2)

	//log.Panicln(cardlsttype, cardlst, cardlsttype2, cardlst2)

	for i := 0; i < 3; i++ {
		cardpoint1 := cardlst[i] / 10
		cardpoint2 := cardlst2[i] / 10

		if cardpoint1 == 1 {
			cardpoint1 = 14
		}

		if cardpoint2 == 1 {
			cardpoint2 = 14
		}

		cardlst[i] = cardpoint1*10 + cardlst[i]%10
		cardlst2[i] = cardpoint2*10 + cardlst2[i]%10
	}

	sort.Sort(LstCard(cardlst))
	sort.Sort(LstCard(cardlst2))

	if cardlsttype == 600 && ZjhIs235(cardlsttype2, cardlst2) {
		return 1
	}

	if cardlsttype2 == 600 && ZjhIs235(cardlsttype, cardlst) {
		return 0
	}

	if cardlsttype > cardlsttype2 {
		return 0
	} else if cardlsttype < cardlsttype2 {
		return 1
	}

	for i := 0; i < 3; i++ {
		cardpoint1 := cardlst[i] / 10
		cardpoint2 := cardlst2[i] / 10

		if cardpoint1 == cardpoint2 {
			continue
		}
		if cardpoint1 > cardpoint2 {
			return 0
		} else {
			return 1
		}
	}

	return 2
}

//! 比牌 0胜 1负  //比花色
func ZjhCardCompare1(_card []int, _card2 []int) int {
	cardlsttype, cardlst := GetZjhType(_card)
	cardlsttype2, cardlst2 := GetZjhType(_card2)

	//log.Panicln(cardlsttype, cardlst, cardlsttype2, cardlst2)

	for i := 0; i < 3; i++ {
		cardpoint1 := cardlst[i] / 10
		cardpoint2 := cardlst2[i] / 10

		if cardpoint1 == 1 {
			cardpoint1 = 14
		}

		if cardpoint2 == 1 {
			cardpoint2 = 14
		}

		cardlst[i] = cardpoint1*10 + cardlst[i]%10
		cardlst2[i] = cardpoint2*10 + cardlst2[i]%10
	}

	sort.Sort(LstCard(cardlst))
	sort.Sort(LstCard(cardlst2))

	/*if cardlsttype == 600 && ZjhIs235(cardlsttype2, cardlst2) {
		return false
	}

	if cardlsttype2 == 600 && ZjhIs235(cardlsttype, cardlst) {
		return true
	}*/

	if cardlsttype > cardlsttype2 {
		return 0
	} else if cardlsttype < cardlsttype2 {
		return 1
	}

	for i := 0; i < 3; i++ { //带花色比较
		cardpoint1 := cardlst[i]
		cardpoint2 := cardlst2[i]
		if cardpoint1 > cardpoint2 {
			return 0
		} else {
			return 1
		}
	}
	return 0
}

//! 比牌 0胜 1负 2平  //全比，先比大小，再比花色
func ZjhCardCompare2(_card []int, _card2 []int) int {
	cardlsttype, cardlst := GetZjhType(_card)
	cardlsttype2, cardlst2 := GetZjhType(_card2)

	//log.Panicln(cardlsttype, cardlst, cardlsttype2, cardlst2)

	for i := 0; i < 3; i++ {
		cardpoint1 := cardlst[i] / 10
		cardpoint2 := cardlst2[i] / 10

		if cardpoint1 == 1 {
			cardpoint1 = 14
		}

		if cardpoint2 == 1 {
			cardpoint2 = 14
		}

		cardlst[i] = cardpoint1*10 + cardlst[i]%10
		cardlst2[i] = cardpoint2*10 + cardlst2[i]%10
	}

	sort.Sort(LstCard(cardlst))
	sort.Sort(LstCard(cardlst2))

	/*if cardlsttype == 600 && ZjhIs235(cardlsttype2, cardlst2) {
		return false
	}

	if cardlsttype2 == 600 && ZjhIs235(cardlsttype, cardlst) {
		return true
	}*/

	if cardlsttype > cardlsttype2 {
		return 0
	} else if cardlsttype < cardlsttype2 {
		return 1
	}

	//比大小
	for i := 0; i < 3; i++ {
		cardpoint1 := cardlst[i] / 10
		cardpoint2 := cardlst2[i] / 10

		if cardpoint1 == cardpoint2 {
			continue
		}
		if cardpoint1 > cardpoint2 {
			return 0
		} else {
			return 1
		}
	}

	//比花色
	for i := 0; i < 3; i++ { //带花色比较
		cardpoint1 := cardlst[i]
		cardpoint2 := cardlst2[i]
		if cardpoint1 > cardpoint2 {
			return 0
		} else {
			return 1
		}
	}
	return 0
}

////////////////////////////////////////////////////
//！捞腌菜获取数组最大值下标
func GetLycMaxIndex(card []int) int {
	temp, index, value, max := 0, -1, 0, 0

	for i := 0; i < len(card); i++ {
		value = card[i] / 10
		if 1 == value {
			temp = card[i] + 13*10
		} else {
			temp = card[i]
		}
		if temp > max {
			index = i
			max = temp
		}
	}
	return index
}

//捞腌菜获取牌值和倍率 10为豹子 其他为正常牌值
func GetLycCardValueMul(card []int) (int, int) {
	total, num, num2, multiple := 0, 0, 0, 0
	for i := 0; i < len(card)-1; i++ {
		if card[i]/10 == card[i+1]/10 {
			num += 1
		}
		if 0 == (card[i]-card[i+1])%10 {
			num2 += 1
		}
	}
	if num == len(card)-1 && 3 == len(card) {
		multiple = 5 //！豹子倍率
		total = 10

	} else {
		for i := 0; i < len(card); i++ {
			if card[i] >= 100 {
				total += 0
			} else {
				total += card[i] / 10
			}
		}
		total = total % 10
		if num == len(card)-1 && 2 == len(card) { //！一对倍率
			multiple = 2
		} else if num2 == len(card)-1 && 3 == len(card) { //!三张同花色倍率
			multiple = 3
		} else if num2 == len(card)-1 && 2 == len(card) { //两张同花色
			multiple = 2
		} else {
			multiple = 1
		}
	}
	return total, multiple
}

////////////////////////////////////////////////////
//！三公演义获取牌值和公牌数量
func GetSgyyCardValue(card []int) (int, int) {
	total, num := 0, 0
	for i := 0; i < len(card); i++ {
		if card[i]/10 > 10 {
			total += 0
			num++
		} else {
			total += card[i] / 10
		}
	}
	return total % 10, num
}

//！三公演义获取牌型和最大牌值
func GetSgyyCardType(card []int) (int, int) {
	style := 0
	for i := 0; i < len(card); i++ {
		for j := i + 1; j < len(card); j++ {
			if card[i] > card[j] {
				temp := card[i]
				card[i] = card[j]
				card[j] = temp
			}
		}
	}
	if IsSgyyThree(card) {
		if card[0] >= 110 {
			style = 5 //大三公
		} else {
			style = 4 //小三公
		}

	} else if IsSgyyConnect(card) { //同花顺
		style = 3
	} else if IsSgyyOtherThree(card) { //混三公
		style = 2
	} else {
		style = 1
	}
	return style, card[len(card)-1]
}
func IsSgyyThree(card []int) bool { //！是否三公
	for i := 0; i < len(card)-1; i++ {
		if card[i]/10 != card[i+1]/10 {
			return false
		}
	}
	return true
}
func IsSgyyConnect(card []int) bool { //！是否同花顺
	for i := 0; i < len(card)-1; i++ {
		if card[i]%10 != card[i+1]%10 || card[i]/10+1 != card[i+1]/10 {
			return false
		}
	}
	return true
}
func IsSgyyOtherThree(card []int) bool { //是否混三公
	for i := 0; i < len(card); i++ {
		if card[i]/10 <= 10 {
			return false
		}
	}
	return true
}

////////////////////////////////////////////////////
//! 通山打拱判断牌型
func IsOkByTSCards(card []int, state int) (int, int) { //1有赖子 2没赖子
	if len(card) == 1 { //! 单张
		return card[0]/10*100 + TYPE_TSDG_ONE, 0
	}
	if len(card) == 2 { //! 对子
		if card[0]/10 == card[1]/10 {
			return card[0]/10*100 + TYPE_TSDG_TWO, 0
		}
	}
	if len(card) == 3 {
		if card[0]/10 == card[1]/10 && card[0]/10 == card[2]/10 { //! 3张
			return card[2]/10*100 + TYPE_TSDG_SAN, 0
		}
		sort.Ints(card)
		if card[0]/10 == 5 && card[1]/10 == 10 && card[2]/10 == 13 {
			return card[2]/10*100 + TYPE_TSDG_BOOM510K, 0
		}
	}
	handcard := make(map[int]int)
	for i := 0; i < len(card); i++ {
		handcard[card[i]/10]++
	}
	var Acard [16][]int
	for i := 0; i < 16; i++ {
		Acard[i] = make([]int, 0) //i张牌的数组
	}
	for key, value := range handcard {
		if value == 1 {
			Acard[0] = append(Acard[0], key*10)
			sort.Ints(Acard[0])
		} else if value == 2 {
			Acard[1] = append(Acard[1], key*10)
			sort.Ints(Acard[1])
		} else if value == 3 {
			Acard[2] = append(Acard[2], key*10)
			sort.Ints(Acard[2])
		} else if value == 4 {
			Acard[3] = append(Acard[3], key*10)
			sort.Ints(Acard[3])
		} else if value == 5 {
			Acard[4] = append(Acard[4], key*10)
			sort.Ints(Acard[4])
		} else if value == 6 {
			Acard[5] = append(Acard[5], key*10)
			sort.Ints(Acard[5])
		} else if value == 7 {
			Acard[6] = append(Acard[6], key*10)
			sort.Ints(Acard[6])
		} else if value == 8 {
			Acard[7] = append(Acard[7], key*10)
			sort.Ints(Acard[7])
		} else if value == 9 {
			Acard[8] = append(Acard[8], key*10)
			sort.Ints(Acard[8])
		} else if value == 10 {
			Acard[9] = append(Acard[9], key*10)
			sort.Ints(Acard[9])
		} else if value == 11 {
			Acard[10] = append(Acard[10], key*10)
			sort.Ints(Acard[10])
		} else if value == 12 {
			Acard[11] = append(Acard[11], key*10)
			sort.Ints(Acard[11])
		} else if value == 13 {
			Acard[12] = append(Acard[12], key*10)
			sort.Ints(Acard[12])
		} else if value == 14 {
			Acard[13] = append(Acard[13], key*10)
			sort.Ints(Acard[13])
		} else if value == 15 {
			Acard[14] = append(Acard[14], key*10)
			sort.Ints(Acard[14])
		} else if value == 16 {
			Acard[15] = append(Acard[15], key*10)
			sort.Ints(Acard[15])
		}
	}
	if len(card) >= 4 {
		for i := 1; i < 16; i++ {
			if len(Acard[i]) == len(handcard) {
				if 1 == len(Acard[i]) {
					if Acard[i][0] == 2000 && i >= 3 {
						long := (i - 1) * 2
						return TYPE_TSDG_BOOMKING, long
					}
					if i >= 3 && i <= 5 { //！没拢炸弹
						return Acard[i][len(Acard[i])-1]*10 + TYPE_TSDG_BOOM, 0
					} else {
						long := Square2Tsdg(i - 5)
						if 1 == state {
							long = long / 2
						}
						return Acard[i][len(Acard[i])-1]*10 + TYPE_TSDG_BOOM, long
					}
				}
				if IsArrayConnect(Acard[i]) {
					if 1 == i { //！连队
						return Acard[i][len(Acard[i])-1]*10 + TYPE_TSDG_SHUNDUI, 0
					} else if 2 == i { //！飞机
						return Acard[i][len(Acard[i])-1]*10 + TYPE_TSDG_SHUNSAN, 0
					} else {
						long := Square2Tsdg(len(Acard[i]) + i - 4)
						if 1 == state {
							long = long / 2
						}
						return Acard[i][len(Acard[i])-1]*10 + TYPE_TSDG_BOOM, long
					}
				} else if Acard[i][0]/10 == 5 && Acard[i][1]/10 == 10 && Acard[i][2]/10 == 13 {
					if i >= 1 && i <= 2 { //！2或者3对510K
						return Acard[i][len(Acard[i])-1]*10 + TYPE_TSDG_BOOM510K, 0
					} else { //！四对及以上
						long := Square2Tsdg(i - 1)
						if 1 == state {
							long = long / 2
						}
						return Acard[i][len(Acard[i])-1]*10 + TYPE_TSDG_BOOM510K, long
					}
				} else {
					return 0, 0
				}
			}
		}
	}
	return 0, 0
}
func IsArrayConnect(card []int) bool {
	for i := 1; i < len(card); i++ {
		if card[i]/10 == 2 {
			return false
		}
		res := CardCompare(card[i-1], card[i])
		if res != -1 {
			if 1 != card[0]/10 {
				return false
			} else {
				if res != len(card)-1 {
					return false
				}
			}
		}
	}
	return true
}
func Square2Tsdg(num int) int {
	if 0 == num {
		return 1
	} else {
		res := 1
		for i := 0; i < num; i++ {
			res = res * 2
		}
		return res
	}
}
func IsTongHua(card []int) int {
	for i := 1; i < len(card); i++ {
		if card[i]%10 != card[i-1]%10 {
			return 0
		}
	}
	return card[0] % 10
}

////////////////////////////////////////////////////
//! 十三水判断牌型
func IsOkBySSSCards(card []int) int {
	sort.Ints(card)
	handcard := make(map[int]int)
	for i := 0; i < len(card); i++ {
		handcard[card[i]/10]++
	}
	var Acard [4][]int
	for i := 0; i < 4; i++ {
		Acard[i] = make([]int, 0) //i张牌的数组
	}
	for key, value := range handcard {
		if value == 1 {
			Acard[0] = append(Acard[0], key*10)
			sort.Ints(Acard[0])
		} else if value == 2 {
			Acard[1] = append(Acard[1], key*10)
			sort.Ints(Acard[1])
		} else if value == 3 {
			Acard[2] = append(Acard[2], key*10)
			sort.Ints(Acard[2])
		} else if value == 4 {
			Acard[3] = append(Acard[3], key*10)
			sort.Ints(Acard[3])
		}
	}
	if 3 == len(card) {
		if 1 == len(Acard[2]) {
			return card[len(card)-1]*10 + TYPE_SSS_SANTIAO
		}
		if 1 == len(Acard[1]) && 1 == len(Acard[0]) {
			if Acard[1][0]/10 == card[len(card)-1]/10 {
				return card[len(card)-1]*10 + TYPE_SSS_DUIZI
			} else {
				return card[1]*10 + TYPE_SSS_DUIZI
			}
		}
		if 3 == len(Acard[0]) {
			if 1 == card[0]/10 {
				return card[0]*10 + TYPE_SSS_WULONG
			} else {
				return card[len(card)-1]*10 + TYPE_SSS_WULONG
			}
		}
	}

	if 5 == len(card) {
		isTonghua := false
		for i := 1; i < len(card); i++ {
			if CardTHCompare(card[i-1]%10, card[i]%10) != 0 {
				isTonghua = false
				break
			}
			isTonghua = true
		}
		if len(Acard[0]) == len(handcard) { //顺子 同花顺  乌龙
			if 1 == card[0]/10 {
				for i := 2; i < len(Acard[0]); i++ {
					if CardSSSCompare(Acard[0][i-1], Acard[0][i]) != -1 {
						if isTonghua {
							return card[0]*10 + TYPE_SSS_TONGHUA
						} else {
							return card[0]*10 + TYPE_SSS_WULONG
						}
					}
				}
				if 13 == card[len(card)-1]/10 || 5 == card[len(card)-1]/10 {
					maxcard := card[len(card)-1] / 10
					if isTonghua {
						if 13 == maxcard {
							return card[0]*10 + TYPE_SSS_TONGHUASHUN
						} else {
							return card[len(card)-1]*10 + TYPE_SSS_TONGHUASHUN
						}

					} else {
						if 13 == maxcard {
							return card[0]*10 + TYPE_SSS_SHUNZI
						} else {
							return card[len(card)-1]*10 + TYPE_SSS_SHUNZI
						}
					}

				} else {
					if isTonghua {
						return card[0]*10 + TYPE_SSS_TONGHUA
					} else {
						return card[0]*10 + TYPE_SSS_WULONG
					}
				}
			} else {
				for i := 1; i < len(Acard[0]); i++ {
					if CardSSSCompare(Acard[0][i-1], Acard[0][i]) != -1 {
						if isTonghua {
							return card[len(card)-1]*10 + TYPE_SSS_TONGHUA
						} else {
							return card[len(card)-1]*10 + TYPE_SSS_WULONG
						}
					}
				}
				if isTonghua {
					return card[len(card)-1]*10 + TYPE_SSS_TONGHUASHUN
				} else {
					return card[len(card)-1]*10 + TYPE_SSS_SHUNZI
				}
			}
		}
		if 1 == len(Acard[3]) && 1 == len(Acard[0]) { //！炸弹
			if Acard[3][0]/10 == card[len(card)-1]/10 {
				return card[len(card)-1]*10 + TYPE_SSS_BOOM
			} else {
				return card[3]*10 + TYPE_SSS_BOOM
			}
		}
		if 1 == len(Acard[2]) && 1 == len(Acard[1]) { //！葫芦
			if Acard[2][0]/10 == card[len(card)-1]/10 {
				return card[len(card)-1]*10 + TYPE_SSS_HULU
			} else {
				return card[2]*10 + TYPE_SSS_HULU
			}
		}
		if 1 == len(Acard[2]) && 2 == len(Acard[0]) { //！三条
			if Acard[2][0]/10 == card[len(card)-1]/10 {
				return card[len(card)-1]*10 + TYPE_SSS_SANTIAO
			}
			if Acard[2][0]/10 == card[0]/10 {
				return card[2]*10 + TYPE_SSS_SANTIAO
			}
			if Acard[2][0]/10 != card[0]/10 && Acard[2][0]/10 != card[len(card)-1]/10 {
				return card[3]*10 + TYPE_SSS_SANTIAO
			}
		}
		if 2 == len(Acard[1]) && 1 == len(Acard[0]) { //2对
			return Acard[1][1]/10*1000 + Acard[1][0]/10*10 + TYPE_SSS_LIANGDUI
			//			if Acard[1][1]/10 == card[len(card)-1]/10 {
			//				return card[len(card)-1]*10 + TYPE_SSS_LIANGDUI
			//			} else {
			//				return card[3]*10 + TYPE_SSS_LIANGDUI
			//			}
		}
		if 1 == len(Acard[1]) && 3 == len(Acard[0]) { //对子
			index := -1
			for i := 0; i < len(card); i++ {
				if card[i]/10 == Acard[1][0]/10 {
					index = i
					break
				}
			}
			return card[index+1]*10 + TYPE_SSS_DUIZI
		}
	}
	return 0
}
func CardSSSCompare(card1 int, card2 int) int {
	if card1/10 == card2/10 {
		return 0
	}
	tmp1 := card1 / 10
	tmp2 := card2 / 10
	return tmp1 - tmp2
}
func CardTHCompare(card1 int, card2 int) int {

	tmp1 := card1 % 10
	tmp2 := card2 % 10
	return tmp1 - tmp2
}

//////////////////////////////////////////////////////////
//！百家乐计算点数
func GetBJLCards(card []int) (int, int) { //点数，类型1一般 2对子
	if len(card) < 2 {
		return -1, 0
	}
	tmp, style := 0, 0
	for i := 0; i < len(card); i++ {
		if card[i]/10 < 10 {
			tmp += card[i] / 10
		}
	}
	tmp = tmp % 10
	if card[0]/10 == card[1]/10 {
		style = 2
	} else {
		style = 1
	}
	return tmp, style
}

//////////////////////////////////////////////////////////
////////////////////////////////////////////////////
//! 判断是否符合要求
func IsOkByCards(card []int) int {
	if len(card) == 1 { //! 单张
		return card[0]/10*100 + TYPE_CARD_ONE
	}

	if len(card) == 2 {
		if (card[0] == 1000 || card[0] == 2000) && (card[1] == 1000 || card[1] == 2000) { //! 王炸
			return TYPE_CARD_WANG
		}

		if card[0]/10 == card[1]/10 { //! 对子
			return card[0]/10*100 + TYPE_CARD_TWO
		}
	}

	if len(card) == 3 {
		if card[0]/10 == card[1]/10 && card[0]/10 == card[2]/10 { //! 3张
			return card[0]/10*100 + TYPE_CARD_SAN
		}
	}

	handcard := make(map[int]int)
	for i := 0; i < len(card); i++ {
		handcard[card[i]/10]++
	}

	if len(card) == 4 {
		if card[0]/10 == card[1]/10 && card[0]/10 == card[2]/10 && card[0]/10 == card[3]/10 { //! 炸弹
			return card[0]/10*100 + TYPE_CARD_ZHA
		}
		for key, value := range handcard {
			if value >= 3 { //! 3带1
				return key*100 + TYPE_CARD_SAN1
			}
		}
	}

	if len(card) == 5 {
		for key, value := range handcard {
			if value >= 3 && len(handcard) == 2 { //! 3带2
				return key*100 + TYPE_CARD_SAN2
			}
		}
	}

	if len(card) >= 5 {
		one := make(LstPoker, 0)
		two := make(LstPoker, 0)
		three := make(LstPoker, 0)
		four := make(LstPoker, 0)
		for key, value := range handcard {
			if value == 1 {
				one = append(one, key*10)
				sort.Sort(LstPoker(one))
			} else if value == 2 {
				two = append(two, key*10)
				sort.Sort(LstPoker(two))
			} else if value == 3 {
				three = append(three, key*10)
				sort.Sort(LstPoker(three))
			} else if value == 4 {
				four = append(four, key*10)
				sort.Sort(LstPoker(four))
			}
		}
		if len(one) == len(handcard) {
			for i := 1; i < len(one); i++ {
				if one[i]/10 == 2 {
					return 0
				}
				if CardCompare(one[i-1], one[i]) != -1 {
					return 0
				}
			}
			return one[0]*10 + TYPE_CARD_SHUN
		}
		if len(two) == len(handcard) { //! 是否顺对
			for i := 1; i < len(two); i++ {
				if two[i]/10 == 2 {
					return 0
				}
				if CardCompare(two[i-1], two[i]) != -1 {
					return 0
				}
			}
			return two[0]*10 + TYPE_CARD_SHUNDUI
		}
		if len(three) == len(handcard) {
			for i := 1; i < len(three); i++ {
				if three[i]/10 == 2 {
					return 0
				}
				if CardCompare(three[i-1], three[i]) != -1 {
					return 0
				}
			}
			return three[0]*10 + TYPE_CARD_SHUNSAN
		}
		if len(four) == 0 && len(three) == len(one) && len(two) == 0 {
			for i := 1; i < len(three); i++ {
				if three[i]/10 == 2 {
					return 0
				}
				if CardCompare(three[i-1], three[i]) != -1 {
					return 0
				}
			}
			return three[0]*10 + TYPE_CARD_SHUNSAN1
		}
		if len(four) == 0 && len(three) == len(two) && len(one) == 0 {
			for i := 1; i < len(three); i++ {
				if three[i]/10 == 2 {
					return 0
				}
				if CardCompare(three[i-1], three[i]) != -1 {
					return 0
				}
			}
			return three[0]*10 + TYPE_CARD_SHUNSAN2
		}
		if len(four) == 0 && len(three) == (len(two)*2+len(one)) {
			for i := 1; i < len(three); i++ {
				if three[i]/10 == 2 {
					return 0
				}
				if CardCompare(three[i-1], three[i]) != -1 {
					return 0
				}
			}
			return three[0]*10 + TYPE_CARD_SHUNSAN1
		}
		if len(four) == 1 && len(one) == 2 && len(three) == 0 && len(two) == 0 {
			return four[0]*10 + TYPE_CARD_SI1
		}
		if len(four) == 1 && len(two) == 1 && len(three) == 0 && len(one) == 0 {
			return four[0]*10 + TYPE_CARD_SI1
		}
		if len(four) == 1 && len(two) == 2 && len(three) == 0 && len(one) == 0 {
			return four[0]*10 + TYPE_CARD_SI2
		}
	}

	return 0
}

///////////////////////////////////////////////////
//! 判断四人斗地主是否符合要求
func IsOKByCardsDBD(card []int) int {
	if len(card) == 0 {
		return 0
	}

	if len(card) == 1 { //! 单张
		if card[0] != 1000 && card[0] != 2000 && card[0] != 3000 {
			return card[0]/10*100 + TYPE_CARD_ONE
		}
	}

	if len(card) == 2 {
		if card[0]/10 == card[1]/10 { //! 对子
			return card[0]/10*100 + TYPE_CARD_TWO
		}
	}

	if len(card) == 3 {
		if card[0]/10 == card[1]/10 && card[0]/10 == card[2]/10 { //! 3张
			return card[0]/10*100 + TYPE_CARD_SAN
		}

		num, index := 0, 0
		for i := 0; i < len(card); i++ {
			if card[i] == 1000 || card[i] == 2000 || card[i] == 3000 {
				num++
				index = i
			}
		}

		if num == 2 {
			for i := 0; i < len(card); i++ {
				if card[i] != 1000 && card[i] != 2000 && card[i] != 3000 {
					return card[i]/10*100 + TYPE_CARD_SAN
				}
			}
		}

		if num == 1 {
			card1, card2 := 0, 0
			for i := 0; i < len(card); i++ {
				if i == index {
					continue
				}
				if card1 == 0 {
					card1 = card[i]
					continue
				}
				if card2 == 0 {
					card2 = card[i]
				}
			}
			if card1/10 == card2/10 && card1 != 0 {
				return card1/10*100 + TYPE_CARD_SAN
			}
		}

	}

	handcard := make(map[int]int)
	for i := 0; i < len(card); i++ {
		handcard[card[i]/10]++
	}

	if len(card) == 4 {
		if card[0]/10 == card[1]/10 && card[0]/10 == card[2]/10 && card[0]/10 == card[3]/10 { //! 炸弹
			return card[0]/10*100 + TYPE_CARD_ZHA
		}

		num, sum, dian := 0, 0, 0
		for key, value := range handcard {
			if key == 100 || key == 200 || key == 300 {
				num++
			} else {
				sum = value
				dian = key
			}
		}
		if (num + sum) == 4 {
			return dian*100 + TYPE_CARD_ZHA
		}
	}

	if len(card) >= 4 {
		one := make(LstPoker, 0)
		two := make(LstPoker, 0)
		three := make(LstPoker, 0)
		four := make(LstPoker, 0)
		for key, value := range handcard {
			if value == 1 {
				one = append(one, key*10)
				sort.Sort(LstPoker(one))
			} else if value == 2 {
				two = append(two, key*10)
				sort.Sort(LstPoker(two))
			} else if value == 3 {
				three = append(three, key*10)
				sort.Sort(LstPoker(three))
			} else if value == 4 {
				four = append(four, key*10)
				sort.Sort(LstPoker(four))
			}
		}
		if len(one) == len(handcard) {
			for i := 1; i < len(one); i++ {
				if one[i]/10 == 2 {
					return 0
				}
				if CardCompare(one[i-1], one[i]) != -1 {
					return 0
				}
			}
			return one[0]*10 + TYPE_CARD_SHUN
		}
		if len(two) == len(handcard) { //! 是否顺对
			for i := 1; i < len(two); i++ {
				if two[i]/10 == 2 {
					return 0
				}
				if CardCompare(two[i-1], two[i]) != -1 {
					return 0
				}
			}
			return two[0]*10 + TYPE_CARD_SHUNDUI
		}
	}

	return 0
}

////////////////////////////////////////////////////
//! 判断跑得快是否符合要求
func IsOkByCardsPDK(card []int, _type int) (int, int) {
	if len(card) == 1 { //! 单张
		return card[0]/10*100 + TYPE_CARD_ONE, -1
	}

	if len(card) == 2 {
		if (card[0] == 1000 || card[0] == 2000) && (card[1] == 1000 || card[1] == 2000) { //! 王炸
			return TYPE_CARD_WANG, -1
		}

		if card[0]/10 == card[1]/10 { //! 对子
			return card[0]/10*100 + TYPE_CARD_TWO, -1
		}
	}

	if len(card) == 3 {
		if card[0]/10 == card[1]/10 && card[0]/10 == card[2]/10 { //! 3张
			return card[0]/10*100 + TYPE_CARD_SAN, -1
		}
	}

	handcard := make(map[int]int)
	for i := 0; i < len(card); i++ {
		handcard[card[i]/10]++
	}

	if len(card) == 4 {
		if card[0]/10 == card[1]/10 && card[0]/10 == card[2]/10 && card[0]/10 == card[3]/10 { //! 炸弹
			return card[0]/10*100 + TYPE_CARD_ZHA, -1
		}
		for key, value := range handcard {
			if value >= 3 { //! 3带1
				return key*100 + TYPE_CARD_SAN1, -1
			}
		}
		if card[0]/10 == card[1]/10 && card[2]/10 == card[3]/10 && card[0]/10+1 == card[2]/10 { //!双对
			return card[0]/10*100 + TYPE_CARD_SHUNDUI, -1
		}

	}

	if len(card) == 5 {
		for key, value := range handcard {
			if value >= 3 && (len(handcard) == 2 || len(handcard) == 3) { //! 3带2
				return key*10 + TYPE_CARD_SAN2, -1
			}
		}
	}

	if len(card) >= 5 {
		one := make(LstPoker, 0)
		two := make(LstPoker, 0)
		three := make(LstPoker, 0)
		four := make(LstPoker, 0)
		for key, value := range handcard {
			if value == 1 {
				one = append(one, key*10)
				sort.Sort(LstPoker(one))
			} else if value == 2 {
				two = append(two, key*10)
				sort.Sort(LstPoker(two))
			} else if value == 3 {
				three = append(three, key*10)
				sort.Sort(LstPoker(three))
			} else if value == 4 {
				four = append(four, key*10)
				sort.Sort(LstPoker(four))
			}
		}
		if len(one) == len(handcard) {
			for i := 1; i < len(one); i++ {
				if one[i]/10 == 2 {
					return 0, -1
				}
				if CardCompare(one[i-1], one[i]) != -1 {
					return 0, -1
				}
			}
			return one[0]*10 + TYPE_CARD_SHUN, -1
		}
		if len(two) == len(handcard) { //! 是否顺对
			for i := 1; i < len(two); i++ {
				if two[i]/10 == 2 {
					return 0, -1
				}
				if CardCompare(two[i-1], two[i]) != -1 {
					return 0, -1
				}
			}
			return two[0]*10 + TYPE_CARD_SHUNDUI, -1
		}

		//需要判斷 333 444 555 777類似這種情況
		if _type%10 == 3 { //三顺
			var i int
			for i = 1; i < len(three); i++ {
				if three[i]/10 == 2 {
					return 0, -1
				}
				if CardCompare(three[i-1], three[i]) != -1 {
					break
					//return 0, -1
				}
			}
			if i >= len(three) {
				return three[0]*10 + TYPE_CARD_SHUNSAN, len(three)
			}
		}
		if _type%10 == 4 { //顺3带1
			maxShunnum := 0
			for i := 1; i < len(three); i++ {
				if three[i]/10 == 2 {
					return 0, -1
				}
				if three[i] <= _type/10 {
					continue
				}
				if CardCompare(three[i-1], three[i]) != -1 {
					return 0, -1
				}
				maxShunnum++
			}
			return _type/10*10 + TYPE_CARD_SHUNSAN1, maxShunnum
		}
		if _type%10 == 5 { //顺3带2
			maxShunnum := 0
			for i := 1; i < len(three); i++ {
				if three[i]/10 == 2 {
					return 0, -1
				}
				if three[i] <= _type/10 {
					continue
				}
				if CardCompare(three[i-1], three[i]) != -1 {
					return 0, -1
				}
				maxShunnum++
			}
			return _type/10*10 + TYPE_CARD_SHUNSAN2, maxShunnum
		}

		if len(four) == 1 && len(one) == 2 && len(three) == 0 && len(two) == 0 {
			return four[0]*10 + TYPE_CARD_SI1, -1
		}
		if len(four) == 1 && len(two) == 1 && len(three) == 0 && len(one) == 0 {
			return four[0]*10 + TYPE_CARD_SI1, -1
		}
		if len(four) == 1 && len(two) == 2 && len(three) == 0 && len(one) == 0 {
			return four[0]*10 + TYPE_CARD_SI2, -1
		}
	}

	return 0, -1
}

func IsBoom(card []int) bool {
	return card[0]/10 == card[1]/10 && card[0]/10 == card[2]/10 && card[0]/10 == card[3]/10
}

func CardCompare(card1 int, card2 int) int {
	if card1/10 == card2/10 {
		return 0
	}

	tmp1 := card1 / 10
	if tmp1 == 1 {
		tmp1 = 14
	} else if tmp1 == 2 {
		tmp1 = 20
	}

	tmp2 := card2 / 10
	if tmp2 == 1 {
		tmp2 = 14
	} else if tmp2 == 2 {
		tmp2 = 20
	}

	return tmp1 - tmp2
}

//! 天王 1000
//! 花五小 500
//! 五小 200
//! 十点半 100
//! 高牌 1-20
//! 爆牌 -1
func GetTenHalfType(card []int) int {
	_card := make([]float32, 0)
	total := float32(0)
	flower := true
	for i := 0; i < len(card); i++ {
		if card[i]/10 > 10 {
			total += 0.5
			_card = append(_card, 0.5)
		} else {
			flower = false
			total += float32(card[i] / 10)
			_card = append(_card, float32(card[i]/10))
		}
	}

	if total > 10.5 {
		return -1
	}

	if len(card) == 5 {
		if total == 10.5 {
			return 1000
		}

		if flower {
			return 500
		}

		return 200
	}

	if total == 10.5 {
		return 100
	}

	return int(total * 2)
}

//! 得到倍数
func GetTenHalfBS(_type int, card int) int {
	if _type == 0 {
		if card == 1000 {
			return 3
		} else if card == 500 {
			return 3
		} else if card == 200 {
			return 2
		} else if card == 100 {
			return 2
		} else {
			return 1
		}
	} else {
		if card == 1000 {
			return 3
		} else if card == 500 {
			return 3
		} else if card == 200 {
			return 2
		} else if card == 100 {
			return 2
		} else {
			return 1
		}
	}
}

//! 得到牌类型
func GetPTJCardType(card int) int {
	if card == 1 { //! 天牌
		return 9
	} else if card == 2 { //! 地牌
		return 8
	} else if card == 3 { //! 人牌
		return 7
	} else if card == 4 { //! 娥牌
		return 6
	} else if card >= 5 && card <= 7 { //! 长牌
		return 5
	} else if card >= 8 && card <= 11 { //! 短牌
		return 4
	}
	return 2
}

//! 得到牌点数
func GetPTJCardNum(card int) int {
	arr := []int{12, 2, 8, 4, 10, 6, 4, 11, 10, 7, 6, 9, 9, 8, 8, 7, 7, 6, 5, 5, 3}

	return arr[card-1]
}

//! 返回牌九类型
func GetPTJType(card1, card2 int) int {
	if card1 > card2 { //! 把小的放前面
		tmp := card1
		card1 = card2
		card2 = tmp
	}

	value, ok := staticfunc.GetCsvMgr().GetPTJ(card1, card2)
	if !ok {
		return 221
	}

	return lib.HF_Atoi(value["id"])
}

////////////////////////////////////////////////
//! 小牌九
func GetXPTJTeShuCardType(card []int) int { //! 得到特殊牌型
	if card[0] > card[1] { //! 小的放前
		temp := card[1]
		card[1] = card[0]
		card[0] = temp
	}
	if card[0] == 18 && card[1] == 21 { //! 大小王
		return 10
	} else if card[0] == 1 && card[1] == 1 { //! 对天
		return 9
	} else if card[0] == 2 && card[1] == 2 { //! 对地
		return 8
	} else if card[0] == 3 && card[1] == 3 { //! 对人
		return 7
	} else if card[0] == 4 && card[1] == 4 { //! 对鹅
		return 6
	} else if card[0] == card[1] && (card[0] >= 5 || card[0] <= 7) { //! 对长
		return 5
	} else if card[0] == card[1] && (card[0] >= 8 || card[0] <= 11) { //! 对短
		return 4
	} else if (card[0] == 12 && card[1] == 13) || (card[0] == 14 && card[1] == 15) || (card[0] == 16 && card[1] == 17) || (card[0] == 19 && card[1] == 20) { //! 对杂
		return 3
	} else if card[0] == 1 && (card[1] == 3 || card[1] == 14 || card[1] == 15) { //! 天杠
		return 2
	} else if card[0] == 2 && (card[1] == 3 || card[1] == 14 || card[1] == 15) { //! 地杠
		return 1
	}
	return 0 //! 非特殊牌
}

func GetXPTJCardNum(no int) int {
	switch no {
	case 1:
		return 12
	case 2:
		return 2
	case 3:
		return 8
	case 4:
		return 4
	case 5:
		return 10
	case 6:
		return 6
	case 7:
		return 4
	case 8:
		return 11
	case 9:
		return 10
	case 10:
		return 7
	case 11:
		return 6
	case 12:
		return 9
	case 13:
		return 9
	case 14:
		return 8
	case 15:
		return 8
	case 16:
		return 7
	case 17:
		return 7
	case 18:
		return 6
	case 19:
		return 5
	case 20:
		return 5
	case 21:
		return 3
	}
	return 0
}

func GetXPTJCardTyoe(card []int) int { //! 非特殊牌型 前缀*10+点数 毕十:0
	/*
		2 杂  	4 短 	5 长 	6 鹅 	7 人 	8 地 	9 天
	*/

	//! 获取前缀
	qian1 := GetPTJCardType(card[0])
	qian2 := GetPTJCardType(card[1])
	var qian int
	if qian1 > qian2 {
		qian = qian1
	} else {
		qian = qian2
	}

	dian := (GetXPTJCardNum(card[0]) + GetXPTJCardNum(card[1])) % 10

	if dian == 0 { //! 毕十
		return 0
	}

	return qian*10 + dian
}

////////////////////////////////////////////////
//! 返回梭哈的类型
func GetSuoHaType(card []int) int {
	//返回三位数 个位 牌型 	十百位 最大点数
	//散牌 1	  对子 2	  两对 3	  三条 4	  满堂红 5  顺子 6  同花 7  四条 8  同花顺 9

	shun := lib.IsShun(card)
	if len(card) == 1 {
		shun = 0
	}
	tongHua := true
	for i := 0; i < len(card)-1; i++ {
		if card[i]%10 != card[i+1]%10 {
			tongHua = false
			break
		}
	}
	if len(card) <= 1 {
		tongHua = false
	}

	maxNum := card[0] / 10
	for i := 1; i < len(card); i++ {
		if card[i-1]/10 == 1 || card[i]/10 == 1 {
			maxNum = 14
			break
		}
		if maxNum < card[i]/10 {
			maxNum = card[i] / 10
		}
	}

	if shun != 0 && tongHua { //! 同花顺
		if card[0]/10 == 14 && card[1]/10 == 2 && card[2]/10 == 3 && card[3]/10 == 4 && card[4]/10 == 5 {
			return 9 + 5*10
		}
		return 9 + maxNum*10
	}
	if shun != 0 { //! 顺子
		if card[0]/10 == 14 && card[1]/10 == 2 && card[2]/10 == 3 && card[3]/10 == 4 && card[4]/10 == 5 {
			return 6 + 5*10
		}
		return 6 + maxNum*10
	}
	if tongHua { //! 同花
		return 7 + maxNum*10
	}

	var _card []int
	lib.HF_DeepCopy(&_card, &card)
	for i := 0; i < len(_card); i++ {
		if _card[i]/10 == 1 {
			_card[i] = 14*10 + _card[i]%10
		}
	}
	handcard := make(map[int]int)
	for i := 0; i < len(_card); i++ {
		handcard[_card[i]/10]++
	}
	one := make(LstPoker, 0)
	two := make(LstPoker, 0)
	three := make(LstPoker, 0)
	four := make(LstPoker, 0)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "handcard", handcard)
	for key, value := range handcard {
		if value == 1 {
			one = append(one, key*10)
			sort.Sort(LstPoker(one))
		} else if value == 2 {
			two = append(two, key*10)
			sort.Sort(LstPoker(two))
		} else if value == 3 {
			three = append(three, key*10)
			sort.Sort(LstPoker(three))
		} else if value == 4 {
			four = append(four, key*10)
			sort.Sort(LstPoker(four))
		}
	}
	if len(one) == len(card) { //! 散牌
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "One[0]", one)
		return 1 + one[len(one)-1]
	}
	if len(two) == 1 && len(three) == 0 { //! 对子
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "Two[0]", two)
		return 2 + two[len(two)-1]
	}
	if len(two) == 2 { //! 两对
		return 3 + two[0]
	}
	if len(three) == 1 && len(two) == 0 { //! 三条
		return 4 + three[0]
	}
	if len(three) == 1 && len(two) == 1 { //!
		return 5 + three[0]
	}
	if len(four) == 1 {
		return 8 + four[0]
	}
	return 0
}

///////////////////////////////////////////////
//! 返回十三道特殊牌型
func GetSSDSpecialType(card []int, special int) int {
	/*
		special
		1 三顺
		2 四条炸弹
		3 四对
		4 双条炸弹
		5 一条龙
	*/

	/*
		return
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

	handcard := make(map[int]int)
	for i := 0; i < len(card); i++ {
		handcard[card[i]/10]++
	}
	si := 0
	dui := 0
	for key := range handcard {
		if handcard[key] == 4 {
			si++
		}
		if handcard[key] == 2 {
			dui++
		}
	}

	head := card[0:2]
	body := card[2:5]
	tail := card[5:]
	sanShun := false
	num := 0
	if lib.IsShun(head) != 0 && lib.IsShun(body) != 0 && lib.IsShun(tail) != 0 {
		sanShun = true
		tongHua := true
		for i := 0; i < len(head)-1; i++ {
			if head[i]%10 != head[i+1]%10 {
				tongHua = false
				break
			}
		}
		if tongHua {
			num++
		}

		tongHua = true
		for i := 0; i < len(body)-1; i++ {
			if body[i]%10 != body[i+1]%10 {
				tongHua = false
				break
			}
		}
		if tongHua {
			num++
		}

		tongHua = true
		for i := 0; i < len(tail)-1; i++ {
			if tail[i]%10 != tail[i+1]%10 {
				tongHua = false
				break
			}
		}
		if tongHua {
			num++
		}
	}

	if si == 2 {
		return 8
	}
	if num == 3 && sanShun { //! 3同花顺
		return 7
	}
	if num == 2 && sanShun { //! 2同花顺
		return 6
	}
	if lib.IsShun(card) != 0 {
		return 5
	}
	if si == 1 {
		return 4
	}
	if dui == 4 {
		return 3
	}
	if num == 1 && sanShun { //! 1同花顺
		return 2
	}
	if num == 0 && sanShun { //! 0同花顺
		return 1
	}
	return 0
}

//! 普通牌型
func GetSSDCardType(card []int) int {
	/*
		 1 乌龙
		 2 对子
		 3 顺子
		 4 炸弹
		 5 同花顺
		A K Q 最大点数为 15
		A 2 3 最大点数为 14
		return 牌型+最大点数*10
	*/
	if len(card) == 2 { //! 头道
		if (card[0] == 1000 && card[1] == 2000) || (card[0] == 2000 && card[1] == 1000) { //对王
			return 2 + 15*10
		}
		if card[0]/10 == card[1]/10 { //! 对子
			if card[0]/10 == 1 {
				return 2 + 14*10
			}
			return 2 + card[0]/10*10
		} else { //! 乌龙
			if card[0]/10 == 1 || card[1]/10 == 1 {
				return 14*10 + 1
			}
			if card[0]/10 > card[1]/10 {
				return card[0]/10*10 + 1
			}
			return card[1]/10*10 + 1
		}
	} else { //! 中 尾道
		for i := 0; i < len(card); i++ { //! 降序排列
			for j := i + 1; j < len(card); j++ {
				if card[j]/10 > card[i]/10 {
					temp := card[i]
					card[i] = card[j]
					card[j] = temp
				}
			}
		}
		shun := true
		for i := 0; i < len(card)-1; i++ {
			if (card[i]/10 - card[i+1]/10) != 1 {
				shun = false
			}
		}

		te := false
		if card[0]/10 == 13 && card[1]/10 == 12 && card[2]/10 == 1 {
			te = true
			shun = true
		}
		if shun {
			tongHua := true
			for i := 0; i < len(card)-1; i++ {
				if card[i]%10 != card[i+1]%10 {
					tongHua = false
					break
				}
			}
			if tongHua {
				if card[len(card)-1]/10 == 1 {
					if te {
						return 5 + 15*10
					}
					return 5 + 14*10
				}
				return 5 + card[0]/10*10 //! 同花顺
			} else {
				if card[len(card)-1]/10 == 1 {
					if te {
						return 3 + 15*10
					}
					return 3 + 14*10
				}
				return 3 + card[0]/10*10 //! 顺子
			}
		}

		handcard := make(map[int]int)
		for i := 0; i < len(card); i++ {
			handcard[card[i]/10]++
		}

		for key := range handcard {
			if handcard[key] == 3 { //! 炸弹
				if key == 1 {
					return 4 + 14*10
				}
				return 4 + key*10
			}
			if handcard[key] == 2 { //! 对子
				if key == 1 {
					return 2 + 14*10
				}
				return 2 + key*10
			}
		}

		if card[2]/10 == 1 {
			return 14*10 + 1
		}
		return 1 + card[0]/10*10
	}
	return 1
}

///////////////////////////////////////////////
func IsOkByCardsQPDDZ(card []int) int {
	if len(card) == 1 { //! 单张
		return card[0]/10*100 + TYPE_CARD_ONE
	}

	if len(card) == 2 { //! 对子
		if card[0]/10 == card[1]/10 {
			return card[0]/10*100 + TYPE_CARD_TWO
		}
	}

	if len(card) == 3 { //! 三张
		if card[0]/10 == card[1]/10 && card[0]/10 == card[2]/10 {
			return card[0]/10*100 + TYPE_CARD_SAN
		}
	}

	if len(card) == 4 { //! 王炸
		if (card[0] == 1000 || card[0] == 2000) && (card[1] == 1000 || card[1] == 2000) && (card[2] == 1000 || card[2] == 2000) && (card[3] == 1000 || card[3] == 2000) {
			return TYPE_CARD_WANG
		}
	}

	if len(card) >= 4 { //炸弹
		zha := true
		for i := 1; i < len(card); i++ {
			if card[0]/10 != card[i]/10 {
				zha = false
				break
			}
		}
		if zha {
			return card[0]/10*100 + TYPE_CARD_ZHA
		}
	}

	handcard := make(map[int]int)
	for i := 0; i < len(card); i++ {
		handcard[card[i]/10]++
	}

	if len(card) == 5 {
		for key, value := range handcard {
			if value >= 3 && len(handcard) == 2 { //! 三带二
				return key*100 + TYPE_CARD_SAN2
			}
		}
	}

	if len(card) >= 5 {
		one := make(LstPoker, 0)
		two := make(LstPoker, 0)
		three := make(LstPoker, 0)
		four := make(LstPoker, 0)
		for key, value := range handcard {
			if value == 1 {
				one = append(one, key*10)
				sort.Sort(LstPoker(one))
			} else if value == 2 {
				two = append(two, key*10)
				sort.Sort(LstPoker(two))
			} else if value == 3 {
				three = append(three, key*10)
				sort.Sort(LstPoker(three))
			} else if value == 4 {
				four = append(four, key*10)
				sort.Sort(LstPoker(four))
			}
		}
		if len(one) == len(handcard) { //! 顺子
			for i := 1; i < len(one); i++ {
				if one[i]/10 == 2 {
					return 0
				}
				if CardCompare(one[i-1], one[i]) != -1 {
					return 0
				}
			}
			return one[0]*10 + TYPE_CARD_SHUN
		}
		if len(two) == len(handcard) { //! 顺对
			for i := 1; i < len(two); i++ {
				if two[i]/10 == 2 {
					return 0
				}
				if CardCompare(three[i-1], three[i]) != -1 {
					return 0
				}
			}
			return two[0]*10 + TYPE_CARD_SHUNDUI
		}
		if len(three) == len(handcard) { //! 顺三
			for i := 1; i < len(three); i++ {
				if three[i]/10 == 2 {
					return 0
				}
				if CardCompare(three[i-1], three[i]) != -1 {
					return 0
				}
			}
			return three[0]*10 + TYPE_CARD_SHUNSAN
		}
		if len(four) == 0 && len(three) == len(two) && len(one) == 0 { //! 顺三带二
			for i := 1; i < len(three); i++ {
				if three[i]/10 == 2 {
					return 0
				}
				if CardCompare(three[i-1], three[i]) != -1 {
					return 0
				}
			}
			return three[0]*10 + TYPE_CARD_SHUNSAN2
		}
	}
	return 0
}

///////////////////////////////////////////////
func IsOkByCardsZYPDK(card []int, can bool) int {
	if len(card) == 1 { //! 单张
		return card[0]/10*100 + TYPE_CARD_ONE
	}

	if len(card) == 2 { //! 对子
		if card[0]/10 == card[1]/10 {
			return card[0]/10*100 + TYPE_CARD_TWO
		}
	}

	if len(card) == 3 { //! 三张
		if card[0]/10 == card[1]/10 && card[0]/10 == card[2]/10 {
			return card[0]/10*100 + TYPE_CARD_SAN
		}
	}

	if len(card) == 4 { //! 炸弹
		zha := true
		for i := 1; i < len(card); i++ {
			if card[0]/10 != card[i]/10 {
				zha = false
				break
			}
		}
		if zha {
			return card[0]/10*100 + TYPE_CARD_ZHA
		}
	}

	handcard := make(map[int]int)
	for i := 0; i < len(card); i++ {
		handcard[card[i]/10]++
	}
	one := make(LstPoker, 0)
	two := make(LstPoker, 0)
	three := make(LstPoker, 0)
	four := make(LstPoker, 0)
	for key, value := range handcard {
		if value == 1 {
			one = append(one, key*10)
			sort.Sort(LstPoker(one))
		} else if value == 2 {
			two = append(two, key*10)
			sort.Sort(LstPoker(two))
		} else if value == 3 {
			three = append(three, key*10)
			sort.Sort(LstPoker(three))
		} else if value == 4 {
			four = append(four, key*10)
			sort.Sort(LstPoker(four))
		}
	}

	if len(card) >= 6 && len(three) >= 2 {
		san := 1
		sanNum := 1
		sum := make([]int, 0)
		sum = append(sum, three[0])
		for i := 1; i < len(three); i++ {
			if CardCompare(three[i-1], three[i]) != -1 {
				san = 1
				sum = append(sum, -1)
				sum = append(sum, three[i])
				continue
			}
			san++
			sum = append(sum, three[i])
			if san >= sanNum {
				sanNum = san
			}
		}
		if sanNum == (len(card)-3*sanNum)/2 && !can {
			san = 0
			s := sum[0]
			for i := 0; i < len(sum); i++ {
				if sum[i] == -1 {
					s = sum[i+1]
				}
				san++
				if san == sanNum {
					break
				}
			}
			return s*10 + TYPE_CARD_SHUNSAN2
		} else if sanNum >= (len(card)-3*sanNum)/2 && can {
			san = 0
			s := sum[0]
			for i := 0; i < len(sum); i++ {
				if sum[i] == -1 {
					s = sum[i+1]
				}
				san++
				if san == sanNum {
					break
				}
			}
			return s*10 + TYPE_CARD_SHUNSAN2
		}
	}

	if len(card) == 5 {
		for key, value := range handcard {
			if value == 3 && (len(handcard) == 2 || len(handcard) == 3) { //! 三带二
				return key*100 + TYPE_CARD_SAN2
			}
		}
	}

	if len(card) == 7 {
		if len(four) == 1 {
			return four[0]*10 + TYPE_CARD_SI3
		}
	}

	if len(card) == 4 {
		if len(two) == len(handcard) { //! 顺对
			for i := 1; i < len(two); i++ {
				if two[i]/10 == 2 {
					return 0
				}
				if CardCompare(two[i-1], two[i]) != -1 {
					return 0
				}
			}
			return two[0]*10 + TYPE_CARD_SHUNDUI
		}
		if len(three) == 1 && len(one) == 1 && len(two) == 0 && len(four) == 0 && can {
			return three[0]*10 + TYPE_CARD_SAN2
		}
	}

	if len(card) >= 5 {
		if len(one) == len(handcard) { //! 顺子
			for i := 1; i < len(one); i++ {
				if one[i]/10 == 2 {
					return 0
				}
				if CardCompare(one[i-1], one[i]) != -1 {
					return 0
				}
			}
			return one[0]*10 + TYPE_CARD_SHUN
		}
		if len(two) == len(handcard) { //! 顺对
			for i := 1; i < len(two); i++ {
				if two[i]/10 == 2 {
					return 0
				}
				if CardCompare(two[i-1], two[i]) != -1 {
					return 0
				}
			}
			return two[0]*10 + TYPE_CARD_SHUNDUI
		}
		/*
			if can {
				if len(four) == 0 && len(three)*2 >= (len(two)*2+len(one)) {
					for i := 1; i < len(three); i++ {
						if three[i]/10 == 2 {
							return 0
						}
						if CardCompare(three[i-1], three[i]) != -1 {
							return 0
						}
					}
					return three[0]*10 + TYPE_CARD_SHUNSAN2
				}
			} else {
				if len(four) == 0 && len(three) == (len(two)*2+len(one))/2 {
					for i := 1; i < len(three); i++ {
						if three[i]/10 == 2 {
							return 0
						}
						if CardCompare(three[i-1], three[i]) != -1 {
							return 0
						}
					}
					return three[0]*10 + TYPE_CARD_SHUNSAN2
				}
			}
		*/
		if len(four) == 1 && (len(two) == 1 || len(one) == 2) {
			return four[0]*10 + TYPE_CARD_SI2
		}
	}
	return 0
}

//////////////////////////////////////////////////////////
//!  内蒙古帕斯牌型
func GetPsType(card []int) (int, int, int) { //! 待测试补充
	//////////! 重复量
	pair := make(map[int]int)
	shoes := make(map[int]int)
	total := 0 //! 总点数
	for i := 0; i < len(card); i++ {
		pair[card[i]/10]++  //! 牌/个数
		shoes[card[i]%10]++ //! 花色/个数
		if card[i]/10 == 200 {
			total += 16
		} else if card[i]/10 == 100 {
			total += 15
		} else {
			total += card[i] / 10
		}
	}

	//////////! 同花
	_pair := make(map[int]int)
	for key, val := range shoes {
		if val >= 3 {
			for i := 0; i < len(card); i++ {
				if key == card[i]%10 {
					_pair[card[i]+(1000*key)] = key
				}
			}
		}
	}

	//////////! 顺子(同花顺)
	shunzi := false                        //! 是否有顺子
	one := make(LstPoker, 0)               //! 排序用
	shunzi_length_max := make(map[int]int) //! 顺子长度/最大值/顺子数量

	for key, _ := range _pair {
		one = append(one, key)
		sort.Sort(LstPoker(one))
	}

	num := 0
	for i := 1; i < len(one); i++ {
		if (one[i-1] - one[i]) == -10 {
			num += 1
		} else if num >= 2 { //! 最后一次没进
			k, v := shunzi_length_max[num+1]
			if v {
				shunzi_length_max[2] = k
			}
			shunzi_length_max[num+1] = one[i-1] % 1000 / 10
			shunzi = true
			num = 0
		} else {
			num = 0
		}

		if i == len(one)-1 && num >= 2 {
			k, v := shunzi_length_max[num+1]
			if v {
				shunzi_length_max[2] = k
			}
			shunzi_length_max[num+1] = one[i] % 1000 / 10
			shunzi = true

		}
	}
	//////////! 炸,拐,王,豹
	max := 0
	zhadan := false    //! 是否有炸弹
	shuangliu := false //! 是否有双六
	king := false      //! 是否有大王
	joker := false     //! 是否有小王
	baozi := false     //! 是否有豹子
	baozi_count := 0   //! 豹子数量
	for key, val := range pair {
		if val == 4 {
			zhadan = true
			max = key
			break
		}
		if key == 6 && val == 2 {
			shuangliu = true
		}
		if key == 200 {
			king = true
		}
		if key == 100 {
			joker = true
		}
		if val == 3 {
			baozi = true
			baozi_count += 1
			if key > max {
				max = key
			}
		}
	}

	if zhadan {
		return TYPE_CARD_ZHADAN, max, 0 //! 20 炸弹
	}

	if shunzi {
		if len(shunzi_length_max) > 1 {
			for length, _max := range shunzi_length_max {
				if length == 4 && shuangliu {
					return TYPE_CARD_SGQ, 0, 0 //! 16 四顺加拐加桥
				} else if length == 4 {
					return TYPE_CARD_SISHUNBAO, _max, shunzi_length_max[3] //! 10 四顺加豹或三顺
				}
				if length == 3 && shuangliu && baozi {
					return TYPE_CARD_QGQ, 0, 0 //! 15 桥拐加桥
				} else if length == 3 && baozi {
					return TYPE_CARD_QG, _max, max //! 9 桥拐
				} else if length == 3 && shuangliu {
					return TYPE_CARD_QG, 0, 0 //! 9 桥拐
				}
			}
		}
		for length, _max := range shunzi_length_max {
			if length == 7 {
				return TYPE_CARD_QISHUN, _max, 0 //! 19 七顺
			}
			if length == 6 && shuangliu {
				return TYPE_CARD_SHUNGUAI, 0, 0 //! 18 六顺加拐
			} else if length == 6 {
				return TYPE_CARD_LIUSHUN, _max, total //! 11 六顺
			}
			if length == 5 && shuangliu && _max > 10 {
				return TYPE_CARD_SHUNBAO, 1, 0 //! 17 五顺加豹
			} else if length == 5 && king && joker {
				return TYPE_CARD_SHUNBAO, 0, 0 //! 17 五顺加豹
			} else if length == 5 && shuangliu {
				return TYPE_CARD_WSG, 15, 0 //! 12 五顺加拐
			} else if length == 5 && baozi {
				return TYPE_CARD_WSG, max, 0 //! 12 五顺加拐
			} else if length == 5 {
				return TYPE_CARD_WUSHUN, _max, total //! 6 五顺
			}
			if length == 4 && baozi && (max > _max || _max-3 > max) {
				return TYPE_CARD_SISHUNBAO, _max, max //! 10 四顺加豹或三顺
			} else if length == 4 && shuangliu {
				return TYPE_CARD_SISHUNGUAI, 0, 0 //! 8 四顺加拐
			} else if length == 4 && baozi {
				return TYPE_CARD_SISHUNGUAI, _max, max //! 8 四顺加拐
			} else if length == 4 {
				return TYPE_CARD_SSHUN, _max, total //! 4 四顺
			}
			if length == 3 && shuangliu && baozi {
				return TYPE_CARD_BGB, 1, 0 //! 14 豹拐加豹
			} else if length == 3 && baozi && king && joker {
				return TYPE_CARD_BGB, 1, 0 //! 14 豹拐加豹
			} else if length == 3 && shuangliu {
				return TYPE_CARD_SSHUNGUAI, 15, 0 //! 5 三顺加拐
			} else if length == 3 && baozi {
				return TYPE_CARD_SSHUNGUAI, _max, max //! 5 三顺加拐
			}
		}
	} else {
		if baozi_count == 3 {
			return TYPE_CARD_SANBAO, max, 0 //! 13 三豹
		}
		if shuangliu && king && joker {
			return TYPE_CARD_SB, 16, 0 //! 7 双豹
		}
		if shuangliu && baozi_count == 1 {
			return TYPE_CARD_SB, 15, 0 //! 7 双豹
		}
		if baozi_count == 2 {
			return TYPE_CARD_SB, max, 0 //! 7 双豹
		}
		if shuangliu {
			return TYPE_CARD_SLIU, 0, 0 //! 3 双六
		}
		if king && joker {
			return TYPE_CARD_SWANG, 0, 0 //! 2 双王
		}
		if baozi_count == 1 {
			return TYPE_CARD_BZ, max, 0 //! 1 一豹
		}
	}
	return 0, 0, total
}

//! 比牌
func PsCardCompare(_card []int, _card2 []int) int {
	cardtype, first, second := GetPsType(_card)
	cardtype2, first2, second2 := GetPsType(_card2)

	if cardtype > cardtype2 { //! 比牌型
		return 0
	} else if cardtype < cardtype2 {
		return 1
	}

	if first > first2 { //! 比最大值
		return 0
	} else if first < first2 {
		return 1
	}

	if second > second2 { //! 比总点数
		return 0
	} else if second < second2 {
		return 1
	}

	return 2 //! 和局
}

//! 判断是否符合要求
func IsOkByGoldPDKCards(card []int, allcard []int) int {
	if len(card) == 1 { //! 单张
		return card[0]/10*100 + TYPE_CARD_ONE
	}

	if len(card) == 2 {
		if (card[0] == 1000 || card[0] == 2000) && (card[1] == 1000 || card[1] == 2000) { //! 王炸
			return TYPE_CARD_WANG
		}

		if card[0]/10 == card[1]/10 { //! 对子
			return card[0]/10*100 + TYPE_CARD_TWO
		}
	}

	handcard := make(map[int]int)
	for i := 0; i < len(card); i++ {
		handcard[card[i]/10]++
	}

	one := make(LstPoker, 0)
	two := make(LstPoker, 0)
	three := make(LstPoker, 0)
	four := make(LstPoker, 0)
	for key, value := range handcard {
		if value == 1 {
			one = append(one, key*10)
			sort.Sort(LstPoker(one))
		} else if value == 2 {
			two = append(two, key*10)
			sort.Sort(LstPoker(two))
		} else if value == 3 {
			three = append(three, key*10)
			sort.Sort(LstPoker(three))
		} else if value == 4 {
			four = append(four, key*10)
			sort.Sort(LstPoker(four))
		}
	}

	if len(three) == 1 {
		if len(three)*2 == len(one)+len(two)*2 {
			return three[0]*10 + TYPE_CARD_SAN
		} else if len(one)+len(two)*2 < len(three)*2 {
			if len(card) == len(allcard) || len(allcard) == 0 {
				return three[0]*10 + TYPE_CARD_SAN
			}
		}
	}

	if len(card) == 4 {
		if card[0]/10 == card[1]/10 && card[0]/10 == card[2]/10 && card[0]/10 == card[3]/10 { //! 炸弹
			return card[0]/10*100 + TYPE_CARD_ZHA
		}
		if len(two) == len(handcard) { //! 是否顺对
			for i := 1; i < len(two); i++ {
				if two[i]/10 == 2 {
					return 0
				}
				if CardCompare(two[i-1], two[i]) != -1 {
					return 0
				}
			}
			return two[0]*10 + TYPE_CARD_SHUNDUI
		}
	}

	if len(card) >= 5 {
		if len(one) == len(handcard) {
			for i := 1; i < len(one); i++ { //! 顺子
				if one[i]/10 == 2 {
					return 0
				}
				if CardCompare(one[i-1], one[i]) != -1 {
					return 0
				}
			}
			return one[0]*10 + TYPE_CARD_SHUN
		}
		if len(two) == len(handcard) { //! 是否顺对
			for i := 1; i < len(two); i++ {
				if two[i]/10 == 2 {
					return 0
				}
				if CardCompare(two[i-1], two[i]) != -1 {
					return 0
				}
			}
			return two[0]*10 + TYPE_CARD_SHUNDUI
		}
		if len(three) >= 2 {
			for i := 1; i < len(three); i++ {
				if three[i]/10 == 2 {
					return 0
				}
				if CardCompare(three[i-1], three[i]) != -1 {
					return 0
				}
			}
			if len(three)*2 == len(one)+len(two)*2 {
				return three[0]*10 + TYPE_CARD_SHUNSAN
			} else if len(one)+len(two)*2 < len(three)*2 {
				//if len(card) == len(allcard) || len(allcard) == 0 {
				return three[0]*10 + TYPE_CARD_SHUNSAN
				//}
			}
		}
		if len(four) == 1 && len(one) == 2 && len(three) == 0 && len(two) == 0 {
			return four[0]*10 + TYPE_CARD_SI1
		}
		if len(four) == 1 && len(two) == 1 && len(three) == 0 && len(one) == 0 {
			return four[0]*10 + TYPE_CARD_SI1
		}
	}

	return 0
}

//////////////////////////////////////////////
//! 翻牌机牌型
func GetFPJType(_card []int) int { //! 0-没有奖 1-10以上对子 2-两对 3-三条 4-顺子 5-同花 6-葫芦 7-四条 8-同花顺 9-同花大顺 10-五条
	card := make([]int, 0)
	lib.HF_DeepCopy(&card, &_card)

	for i := 0; i < len(card); i++ {
		if card[i]/10 == 1 {
			card[i] = 140 + card[i]%10
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

	for i := 0; i < len(card); {
		if card[i] == 1000 || card[i] == 2000 {
			copy(card[i:], card[i+1:])
			card = card[:len(card)-1]
		} else {
			i++
		}
	}

	tonghua := true

	for i := 1; i < len(card); i++ {
		if card[0]%10 != card[i]%10 {
			tonghua = false
			break
		}
	}

	shunzi := true
	lz := 5 - len(card)
	for i := 0; i < len(card)-1; i++ {
		if card[i]/10-card[i+1]/10 != 1 {
			if lz == 0 {
				shunzi = false
				break
			} else {
				if card[i]/10-card[i+1]/10 == 2 {
					lz--
				} else {
					shunzi = false
					break
				}
			}
		}
	}

	{
		c := make([]int, 0)
		lib.HF_DeepCopy(&c, &card)
		for i := 0; i < len(c); i++ {
			if c[i]/10 == 14 {
				c[i] = 10 + c[i]%10
			}
		}

		for i := 0; i < len(c); i++ { //! 降序排列
			for j := i + 1; j < len(c); j++ {
				if c[j]/10 > c[i]/10 {
					temp := c[i]
					c[i] = c[j]
					c[j] = temp
				}
			}
		}

		_shunzi := true
		lz := 5 - len(c)
		for i := 0; i < len(c)-1; i++ {
			if c[i]/10-c[i+1]/10 != 1 {
				if lz == 0 {
					_shunzi = false
					break
				} else {
					if c[i]/10-c[i+1]/10 == 2 {
						lz--
					} else {
						_shunzi = false
						break
					}
				}
			}
		}

		if _shunzi {
			shunzi = true
		}

	}

	handcard := make(map[int]int)
	for i := 0; i < len(card); i++ {
		handcard[card[i]/10]++
	}
	one := make(LstPoker, 0)
	two := make(LstPoker, 0)
	three := make(LstPoker, 0)
	four := make(LstPoker, 0)
	for key, value := range handcard {
		if value == 1 {
			one = append(one, key*10)
			sort.Sort(LstPoker(one))
		} else if value == 2 {
			two = append(two, key*10)
			sort.Sort(LstPoker(two))
		} else if value == 3 {
			three = append(three, key*10)
			sort.Sort(LstPoker(three))
		} else if value == 4 {
			four = append(four, key*10)
			sort.Sort(LstPoker(four))
		}
	}

	if len(card) == 5 { //! 没有癞子
		if shunzi && tonghua {
			if card[0]/10 == 14 { //! 同花大顺
				return 9
			} else { //! 同花顺
				return 8
			}
		}
		if len(four) == 1 { //! 四条
			return 7
		}
		if len(three) == 1 && len(two) == 1 { //! 葫芦
			return 6
		}
		if tonghua { //! 同花
			return 5
		}
		if shunzi { //! 顺子
			return 4
		}
		if len(three) == 1 { //! 三条
			return 3
		}
		if len(two) == 2 { //! 两对
			return 2
		}
		if len(two) == 1 && two[0] >= 80 { //! 一对8以上
			return 1
		}
	} else if len(card) == 4 { //!　一个癞子
		if len(four) == 1 { //!　5条
			return 10
		}
		if shunzi && tonghua {
			dashun := true
			for i := 0; i < len(card); i++ {
				if card[i]/10 < 10 {
					dashun = false
					break
				}
			}
			if dashun {
				return 9
			} else {
				return 8
			}
		}
		if len(three) == 1 { //! 四条
			return 7
		}
		if len(two) == 2 { //!　葫芦
			return 6
		}
		if tonghua { //! 同花
			return 5
		}
		if shunzi { //!顺子
			return 4
		}
		if len(two) == 1 { //! 三条
			return 3
		}
		if card[0]/10 >= 8 { //! 一对10以上
			return 1
		}
	} else if len(card) == 3 { //! 两个癞子
		if len(three) == 1 { //! 五条
			return 10
		}
		if tonghua && shunzi {
			dashun := true
			for i := 0; i < len(card); i++ {
				if card[i]/10 < 10 {
					dashun = false
					break
				}
			}
			if dashun { //! 同花大顺
				return 9
			} else { //! 同花顺
				return 8
			}
		}
		if len(two) == 1 { //! 四条
			return 7
		}
		if tonghua {
			return 5
		}
		if shunzi {
			return 4
		}
		return 3
	}

	return 0
}
