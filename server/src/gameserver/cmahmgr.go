//! 麻将的逻辑

package gameserver

import (
	"lib"
	"math/rand"
	"time"
)

const TYPE_MAH_PENG = 1 //! 可碰
const TYPE_MAH_GANG = 2 //! 可杠，必然可碰

const TYPE_HU_PH = 1   //! 屁胡
const TYPE_HU_7DUI = 2 //! 7对
const TYPE_HU_PPH = 3  //! 砰砰胡
const TYPE_HU_KWX = 4  //! 卡五星
const TYPE_HU_7DUI1 = 5
const TYPE_HU_7DUI2 = 6
const TYPE_HU_7DUI3 = 7

type LstCard []int

func (a LstCard) Len() int { // 重写 Len() 方法
	return len(a)
}

func (a LstCard) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}

func (a LstCard) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j] < a[i]
}

//!
type Mah_Card struct { //玩家状态
	Card1  LstCard `json:"card1"`  //! 手牌
	Card2  []int   `json:"card2"`  //! 已经出过的牌
	CardL  []int   `json:"cardl"`  //! 亮出的牌
	CardM  int     `json:"cardm"`  //! 当前摸的牌
	CardP  []int   `json:"cardp"`  //! 当前碰的牌
	CardC  [][]int `json:"cardc"`  //! 当前吃的牌
	CardMG []int   `json:"cardmg"` //! 明杠
	CardAG []int   `json:"cardag"` //! 暗杠
	CardCG []int   `json:"cardcg"` //! 擦杠
	Want   []int   `json:"want"`   //! 要胡的牌，亮牌之后有用
}

func (self *Mah_Card) Init() { //开始游戏时，创建玩家状态对象，并初始化为空
	self.Card1 = make([]int, 0)
	self.Card2 = make([]int, 0)
	self.CardL = make([]int, 0)
	self.CardM = 0
	self.CardP = make([]int, 0)
	self.CardC = make([][]int, 0)
	self.CardMG = make([]int, 0)
	self.CardAG = make([]int, 0)
	self.CardCG = make([]int, 0)
	self.Want = make([]int, 0)
}

type Chi struct {
	C1 int `json:"c1"`
	C2 int `json:"c1"`
}

type MahMgr struct {
	Card   []int //! 剩余牌组
	Temp   []int
	random *rand.Rand

	//Want []int
}

//! 牌=牌型*10+点数
//! 牌型:0条,1筒,2万,3字,4花
//! 点数:1中,2发,3白,4东,5南,6西,7北
//! new:5中,6发,7白,1东,2南,3西,4北
//! 花: 1-4:春夏秋冬 5-8:梅兰竹菊
func NewMah_KWX() *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		mgr.Card = append(mgr.Card, 31)
		mgr.Card = append(mgr.Card, 32)
		mgr.Card = append(mgr.Card, 33)
		for j := 1; j <= 9; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 11; j <= 19; j++ {
			mgr.Card = append(mgr.Card, j)
		}
	}

	return mgr
}

//! 涡阳麻将
func NewMah_GYMJ() *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		for j := 1; j <= 9; j++ { //! 1-9条
			mgr.Card = append(mgr.Card, j)
		}
		for j := 11; j <= 19; j++ { //! 1-9筒
			mgr.Card = append(mgr.Card, j)
		}
		for j := 21; j <= 29; j++ { //! 1-9万
			mgr.Card = append(mgr.Card, j)
		}
		for j := 31; j <= 37; j++ { //! 东南西北中发白
			mgr.Card = append(mgr.Card, j)
		}
	}

	for i := 41; i <= 48; i++ { //! 春夏秋冬梅兰竹菊
		mgr.Card = append(mgr.Card, i)
	}
	return mgr
}

//! 推筒子牌组
func NewMah_TTZ() *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		for j := 11; j <= 19; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		mgr.Card = append(mgr.Card, 37)
	}

	return mgr
}

//! 推筒子牌组
func NewMah_TTZGold() *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 2; i++ {
		for j := 11; j <= 19; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		mgr.Card = append(mgr.Card, 37)
	}

	return mgr
}

//! 推对子牌组
func NewMah_TDZ() *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		for j := 11; j <= 19; j++ {
			mgr.Card = append(mgr.Card, j)
		}
	}

	return mgr
}

//! 杠次
func NewMah_GC(zi bool) *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		for j := 1; j <= 9; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 11; j <= 19; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 21; j <= 29; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		if zi {
			for j := 31; j <= 37; j++ {
				mgr.Card = append(mgr.Card, j)
			}
		}
	}

	return mgr
}

//! 一腳賴油
func NewMah_YJLY() *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		for j := 1; j <= 9; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 11; j <= 19; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 21; j <= 29; j++ {
			mgr.Card = append(mgr.Card, j)
		}
	}

	return mgr
}

//! 仙桃人人晃晃
func NewMah_XTRRHH1() *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		for j := 1; j <= 9; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 11; j <= 19; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 21; j <= 29; j++ {
			mgr.Card = append(mgr.Card, j)
		}
	}

	return mgr
}

func NewMah_XTRRHH2() *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		for j := 1; j <= 9; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 11; j <= 19; j++ {
			mgr.Card = append(mgr.Card, j)
		}
	}

	return mgr
}

//! 血战
func NewMah_XZDD() *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		for j := 1; j <= 9; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 11; j <= 19; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 21; j <= 29; j++ {
			mgr.Card = append(mgr.Card, j)
		}
	}

	return mgr
}

//! 宿州
func NewMah_SZMJ() *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		for j := 1; j <= 9; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 11; j <= 19; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 21; j <= 29; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 31; j <= 37; j++ {
			mgr.Card = append(mgr.Card, j)
		}
	}
	for j := 41; j <= 48; j++ {
		mgr.Card = append(mgr.Card, j)
	}

	return mgr
}

//! 安庆
func NewMah_AQMJ() *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		for j := 1; j <= 9; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 11; j <= 19; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 21; j <= 29; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 31; j <= 33; j++ { //南西北
			mgr.Card = append(mgr.Card, j)
		}
		for j := 41; j <= 44; j++ { //中发白东
			mgr.Card = append(mgr.Card, j)
		}
	}

	return mgr
}

//! 焦作麻将
func NewMah_JZMJ() *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		for j := 1; j <= 9; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 11; j <= 19; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 21; j <= 29; j++ {
			mgr.Card = append(mgr.Card, j)
		}
	}

	return mgr
}

//! 推倒胡
func NewMah_TDH() *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		for j := 1; j <= 9; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 11; j <= 19; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 21; j <= 29; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 31; j <= 37; j++ { //东南西北中发白
			mgr.Card = append(mgr.Card, j)
		}
	}

	for j := 41; j <= 48; j++ { //花牌
		mgr.Card = append(mgr.Card, j)
	}

	return mgr
}

//! 杭州
func NewMah_HZMJ() *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		for j := 1; j <= 9; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 11; j <= 19; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 21; j <= 29; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 31; j <= 37; j++ { //东南西北中发白
			mgr.Card = append(mgr.Card, j)
		}
	}
	return mgr
}

//! 常熟
func NewMah_CSMJ() *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		for j := 1; j <= 9; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 11; j <= 19; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 21; j <= 29; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 31; j <= 36; j++ { //东南西北中发
			mgr.Card = append(mgr.Card, j)
		}
	}
	for j := 41; j <= 48; j++ { //春夏秋冬梅兰竹菊白
		mgr.Card = append(mgr.Card, j)
	}
	for j := 0; j < 4; j++ { //白
		mgr.Card = append(mgr.Card, 49)
	}
	return mgr
}

//! 上虞花
func NewMah_SYHMJ() *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		for j := 1; j <= 9; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 11; j <= 19; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 21; j <= 29; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 31; j <= 37; j++ { //东南西北中发白
			mgr.Card = append(mgr.Card, j)
		}
	}
	for j := 41; j <= 48; j++ { //春夏秋冬梅兰竹菊
		mgr.Card = append(mgr.Card, j)
	}
	return mgr
}

//! 南昌
func NewMah_NCMJ() *MahMgr {
	mgr := new(MahMgr)
	mgr.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		for j := 1; j <= 9; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 11; j <= 19; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 21; j <= 29; j++ {
			mgr.Card = append(mgr.Card, j)
		}
		for j := 31; j <= 37; j++ { //东南西北中发白
			mgr.Card = append(mgr.Card, j)
		}
	}
	return mgr
}

//! 发牌
func (self *MahMgr) Deal(num int) []int {
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

//! 发牌
func (self *MahMgr) DealNeed(card []int) []int {
	lst := make([]int, 0)

	for _, value := range card {
		for i := 0; i < len(self.Card); i++ {
			if self.Card[i] == value {
				copy(self.Card[i:], self.Card[i+1:])
				self.Card = self.Card[:len(self.Card)-1]
				if len(lst) < 13 {
					lst = append(lst, value)
				} else {
					self.Temp = append(self.Temp, value)
				}
				break
			}
		}
	}

	num := 13 - len(lst)
	for i := 0; i < num; i++ {
		lst = append(lst, self.Draw2())
	}

	return lst
}

//! 去掉一张牌
func (self *MahMgr) Del(card int) {
	for i := 0; i < len(self.Card); i++ {
		if self.Card[i] == card {
			copy(self.Card[i:], self.Card[i+1:])
			self.Card = self.Card[:len(self.Card)-1]
			break
		}
	}
}

//增加一张牌
func (self *MahMgr) AddCard(card int) {
	self.Card = append(self.Card, card)
}

//! 摸牌
func (self *MahMgr) Draw2() int {
	if self.random == nil {
		self.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	index := self.random.Intn(len(self.Card))
	card := self.Card[index]
	copy(self.Card[index:], self.Card[index+1:])
	self.Card = self.Card[:len(self.Card)-1]

	return card
}

//! 摸牌
func (self *MahMgr) Draw() int {
	if self.random == nil {
		self.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	if len(self.Card) == 0 {
		return -1
	}

	ts := 0
	//if len(self.Want) > 0 {
	//	ts = self.Want[0]
	//	self.Want = self.Want[1:]
	//}
	if ts != 0 {
		for i := 0; i < len(self.Card); i++ {
			if self.Card[i] == ts {
				copy(self.Card[i:], self.Card[i+1:])
				self.Card = self.Card[:len(self.Card)-1]
				return ts
			}
		}
	}

	if len(self.Temp) > 0 {
		card := self.Temp[0]
		self.Temp = self.Temp[1:]
		return card
	}

	index := self.random.Intn(len(self.Card))
	card := self.Card[index]
	copy(self.Card[index:], self.Card[index+1:])
	self.Card = self.Card[:len(self.Card)-1]

	return card
}

//! 摸牌不丢弃情况
func (self *MahMgr) Draw3() int {
	if self.random == nil {
		self.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	ts := 0
	//if len(self.Want) > 0 {
	//	ts = self.Want[0]
	//	self.Want = self.Want[1:]
	//}
	if ts != 0 {
		for i := 0; i < len(self.Card); i++ {
			if self.Card[i] == ts {
				copy(self.Card[i:], self.Card[i+1:])
				self.Card = self.Card[:len(self.Card)-1]
				return ts
			}
		}
	}

	if len(self.Temp) > 0 {
		card := self.Temp[0]
		self.Temp = self.Temp[1:]
		return card
	}

	index := self.random.Intn(len(self.Card))
	card := self.Card[index]

	return card
}

//! 摸指定牌
func (self *MahMgr) Draw4(card int) int {
	if self.random == nil {
		self.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	if card != 0 {
		for i := 0; i < len(self.Card); i++ {
			if self.Card[i] == card {
				copy(self.Card[i:], self.Card[i+1:])
				self.Card = self.Card[:len(self.Card)-1]
				return card
			}
		}
	}

	index := self.random.Intn(len(self.Card))
	card = self.Card[index]
	copy(self.Card[index:], self.Card[index+1:])
	self.Card = self.Card[:len(self.Card)-1]

	return card
}

//! 摸指定牌
func (self *MahMgr) Draw5(card int) int {
	if self.random == nil {
		self.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	if card != 0 {
		for i := 0; i < len(self.Card); i++ {
			if self.Card[i] == card {
				copy(self.Card[i:], self.Card[i+1:])
				self.Card = self.Card[:len(self.Card)-1]
				return card
			}
		}
	}
	return 0
}

//! 是否能吃  //安庆
func MahIsChi(card *Mah_Card, _card int) (bool, []Chi) {
	var ch []Chi
	if len(card.Want) != 0 {
		return false, ch
	}
	var IsChi bool
	var t1, t2, t3 Chi
	t1.C1 = _card + 1
	t1.C2 = _card + 2

	t2.C1 = _card + 1
	t2.C2 = _card - 1

	t3.C1 = _card - 1
	t3.C2 = _card - 2

	tmp := make(map[int]int)
	//	if _card != 0 {
	//		tmp[_card]++
	//	}
	for i := 0; i < len(card.Card1); i++ {
		if card.Card1[i]/10 < 3 {
			tmp[card.Card1[i]]++
		} else {
			tmp[card.Card1[i]] = 0
		}
	}
	if tmp[t1.C1] > 0 && tmp[t1.C2] > 0 {
		ch = append(ch, t1)
		IsChi = true
	}
	if tmp[t2.C1] > 0 && tmp[t2.C2] > 0 {
		ch = append(ch, t2)
		IsChi = true
	}
	if tmp[t3.C1] > 0 && tmp[t3.C2] > 0 {
		ch = append(ch, t3)
		IsChi = true
	}

	return IsChi, ch
}

//! 是否能吃  //南昌
func MahIsChiNC(card *Mah_Card, _card int) (bool, []Chi) {
	var ch []Chi
	if len(card.Want) != 0 {
		return false, ch
	}
	var IsChi bool
	var t1, t2, t3 Chi

	FengChi := map[int][][]int{
		31: {{32, 33}, {33, 34}, {34, 32}},
		32: {{33, 34}, {31, 33}, {31, 34}},
		33: {{31, 34}, {31, 32}, {32, 34}},
		34: {{31, 32}, {31, 33}, {32, 33}},
	}

	if _card/10 < 3 {
		t1.C1 = _card + 1
		t1.C2 = _card + 2

		t2.C1 = _card + 1
		t2.C2 = _card - 1

		t3.C1 = _card - 1
		t3.C2 = _card - 2
	} else if _card <= 34 {
		for i := 31; i <= 34; i++ {
			if i == _card {
				t1.C1 = FengChi[i][0][0]
				t1.C2 = FengChi[i][0][1]

				t2.C1 = FengChi[i][1][0]
				t2.C2 = FengChi[i][1][1]

				t3.C1 = FengChi[i][2][0]
				t3.C2 = FengChi[i][2][1]
			}
		}
	}

	tmp := make(map[int]int)
	for i := 0; i < len(card.Card1); i++ {
		tmp[card.Card1[i]]++
	}

	if _card <= 34 {
		if tmp[t1.C1] > 0 && tmp[t1.C2] > 0 {
			ch = append(ch, t1)
			IsChi = true
		}
		if tmp[t2.C1] > 0 && tmp[t2.C2] > 0 {
			ch = append(ch, t2)
			IsChi = true
		}
		if tmp[t3.C1] > 0 && tmp[t3.C2] > 0 {
			ch = append(ch, t3)
			IsChi = true
		}
	} else {
		for i := 35; i <= 37; i++ {
			if i != _card {
				if t1.C1 == 0 {
					t1.C1 = i
				} else {
					t1.C2 = i
				}
			}
		}
		if tmp[t1.C1] > 0 && tmp[t1.C2] > 0 {
			ch = append(ch, t1)
			IsChi = true
		}
	}

	return IsChi, ch
}

//! 是否能碰
func MahIsPeng(card *Mah_Card, _card int) bool {
	if len(card.Want) != 0 {
		return false
	}

	num := 0
	for i := 0; i < len(card.Card1); i++ {
		if card.Card1[i] == _card {
			num++
			if num >= 2 {
				return true
			}
		}
	}

	return false
}

//! 是否能碰   安庆麻将
func MahIsPengAQMJ(card *Mah_Card, _card int) bool {
	if _card/10 == 4 {
		return false
	}
	num := 0
	for i := 0; i < len(card.Card1); i++ {
		if card.Card1[i] == _card {
			num++
			if num >= 2 {
				return true
			}
		}
	}

	return false
}

//! 是否能杠
func MahIsGang(card *Mah_Card, _card int, self bool) bool {
	if self {
		tmp := make(map[int]int)
		if _card != 0 {
			tmp[_card]++
		}
		for i := 0; i < len(card.Card1); i++ {
			tmp[card.Card1[i]]++
		}

		for key, value := range tmp {
			if value >= 4 {
				if len(card.Want) > 0 { //! 已经亮牌，则要判断杠了之后是否能胡
					newcard := new(Mah_Card)
					lib.HF_DeepCopy(newcard, card)
					for i := 0; i < len(newcard.Card1); { //! 去掉杠出的牌
						if newcard.Card1[i] == key {
							copy(newcard.Card1[i:], newcard.Card1[i+1:])
							newcard.Card1 = newcard.Card1[:len(newcard.Card1)-1]
						} else {
							i++
						}
					}
					score := false
					for i := 0; i < len(newcard.Want); i++ {
						score = false
						score = MahIsHu(newcard, newcard.Want[i])
						if !score {
							break
						}
					}
					if score {
						return true
					}
				} else {
					return true
				}
			}
		}

		return false
	} else {
		num := 0
		for _, value := range card.Card1 {
			if value == _card {
				num++
				if num >= 3 {
					if len(card.Want) > 0 {
						newcard := new(Mah_Card)
						lib.HF_DeepCopy(newcard, card)
						for i := 0; i < len(newcard.Card1); { //! 去掉杠出的牌
							if newcard.Card1[i] == _card {
								copy(newcard.Card1[i:], newcard.Card1[i+1:])
								newcard.Card1 = newcard.Card1[:len(newcard.Card1)-1]
							} else {
								i++
							}
						}
						score := false
						for i := 0; i < len(newcard.Want); i++ {
							score = false
							score = MahIsHu(newcard, newcard.Want[i])
							if !score {
								break
							}
						}
						if score {
							return true
						}
					} else {
						return true
					}
				}
			}
		}

		return false
	}
}

//! 是否能杠(癞子)
//func MahIsGangByRazz(card *Mah_Card, _card int, self bool, razz int) bool {
//	if self {
//		var lst LstCard
//		lib.HF_DeepCopy(&lst, &card.Card1)
//		if _card != 0 {
//			lst = append(lst, _card)
//		}
//		razznum := 0
//		//! 将癞子拿出
//		for i := 0; i < len(lst); {
//			if lst[i] == razz {
//				razznum++
//				copy(lst[i:], lst[i+1:])
//				lst = lst[:len(lst)-1]
//			} else {
//				i++
//			}
//		}
//		tmp := make(map[int]int)
//		for i := 0; i < len(lst); i++ {
//			tmp[lst[i]]++
//		}
//		for _, value := range tmp {
//			if value+razznum >= 4 {
//				return true
//			}
//		}
//		return false
//	} else {
//		num := 0
//		for _, value := range card.Card1 {
//			if value == _card {
//				num++
//				if num >= 3 {
//					return true
//				}
//			}
//		}
//		return false
//	}
//}

//仙桃
func MahIsGangLZ(card *Mah_Card, _card int, self bool, Lz int) bool {
	if self {
		tmp := make(map[int]int)
		if _card != 0 && _card != Lz {
			tmp[_card]++
		}
		for i := 0; i < len(card.Card1); i++ {
			if card.Card1[i]/10 != 4 && card.Card1[i] != Lz {
				tmp[card.Card1[i]]++
			} else {
				tmp[card.Card1[i]] = 0
			}
		}

		for _, value := range tmp {
			if value >= 4 {
				return true
			}
		}

		return false
	} else {
		num := 0
		for _, value := range card.Card1 {
			if value == _card && _card != Lz {
				num++
				if num >= 3 {
					return true
				}
			}
		}

		return false
	}
}

//! 是否能杠   安庆麻将  //常熟 //南昌
func MahIsGangAQMJ(card *Mah_Card, _card int, self bool) bool {
	if self {
		tmp := make(map[int]int)
		if _card != 0 {
			tmp[_card]++
		}
		for i := 0; i < len(card.Card1); i++ {
			if card.Card1[i]/10 != 4 {
				tmp[card.Card1[i]]++
			} else {
				tmp[card.Card1[i]] = 0
			}
		}

		for _, value := range tmp {
			if value >= 4 {
				return true
			}
		}

		return false
	} else {
		num := 0
		for _, value := range card.Card1 {
			if value == _card {
				num++
				if num >= 3 {
					return true
				}
			}
		}

		return false
	}
}

//! 是否能杠
func MahIsGangYJLY(card *Mah_Card, _card int, self bool, lzp int) bool {
	if self {
		tmp := make(map[int]int)
		if _card != 0 {
			if _card == lzp {
				tmp[_card] += 2
			} else {
				tmp[_card]++
			}
		}
		for i := 0; i < len(card.Card1); i++ {
			tmp[card.Card1[i]]++
		}

		for _, value := range tmp {
			if value >= 4 {
				return true
			}
		}

		return false
	} else {
		num := 0
		for _, value := range card.Card1 {
			if value == _card {
				num++
				if _card == lzp {
					if num >= 2 {
						return true
					}
				} else {
					if num >= 3 {
						return true
					}
				}
			}
		}

		return false
	}
}

//! 是否能补杠
func MahIsBuGang(card *Mah_Card, _card int) bool {
	var lst LstCard
	lib.HF_DeepCopy(&lst, &card.Card1)
	if _card != 0 {
		lst = append(lst, _card)
	}

	for i := 0; i < len(lst); i++ {
		for j := 0; j < len(card.CardP); j++ {
			if card.CardP[j] == lst[i] {
				return true
			}
		}
	}

	return false
}

//! 是否能补杠
//func MahIsBuGangByQue(card *Mah_Card, _card int, que int) bool {
//	var lst LstCard
//	lib.HF_DeepCopy(&lst, &card.Card1)
//	lst = append(lst, _card)

//	for i := 0; i < len(lst); i++ {
//		if lst[i]/10 == que {
//			continue
//		}
//		for j := 0; j < len(card.CardP); j++ {
//			if card.CardP[j] == lst[i] {
//				return true
//			}
//		}
//	}

//	return false
//}

//! 是否能补杠(癞子)
//func MahIsBuGangByRazz(card *Mah_Card, _card int, razz int) bool {
//	if _card == razz {
//		return len(card.CardP) > 0
//	}

//	for i := 0; i < len(card.CardP); i++ {
//		if card.CardP[i] == _card {
//			return true
//		}
//	}

//	return false
//}

//! 是否能次
func MahIsCi(card *Mah_Card, _card int, self bool, ci int) int {
	if self {
		num := 0
		for i := 0; i < len(card.Card1); i++ {
			if card.Card1[i] == ci {
				num++
			}
		}
		if _card == ci {
			num++
		}
		if num >= 3 { //! 流氓次
			return 1000
		}

		tmp := make(map[int]int)
		if _card != 0 {
			tmp[_card]++
		}
		for i := 0; i < len(card.Card1); i++ {
			tmp[card.Card1[i]]++
		}

		//! 先判断一次暗次
		for key, value := range tmp {
			if value < 4 {
				continue
			}

			clonecard := new(Mah_Card)
			lib.HF_DeepCopy(clonecard, card)
			clonecard.Card1 = append(clonecard.Card1, _card)

			for i := 0; i < len(clonecard.Card1); {
				if clonecard.Card1[i] == key {
					copy(clonecard.Card1[i:], clonecard.Card1[i+1:])
					clonecard.Card1 = clonecard.Card1[:len(clonecard.Card1)-1]
				} else {
					i++
				}
			}

			ptype := MahIsHu(clonecard, ci)

			if ptype {
				return key
			}

			tmp[key] -= value
		}

		for i := 0; i < len(card.CardP); i++ {
			tmp[card.CardP[i]] += 3
		}

		for key, value := range tmp {
			if value < 4 {
				continue
			}

			clonecard := new(Mah_Card)
			lib.HF_DeepCopy(clonecard, card)
			clonecard.Card1 = append(clonecard.Card1, _card)

			for i := 0; i < len(clonecard.Card1); {
				if clonecard.Card1[i] == key {
					copy(clonecard.Card1[i:], clonecard.Card1[i+1:])
					clonecard.Card1 = clonecard.Card1[:len(clonecard.Card1)-1]
				} else {
					i++
				}
			}

			ptype := MahIsHu(clonecard, ci)

			if ptype {
				return key
			}
		}

		return 0
	} else {
		clonecard := new(Mah_Card)
		lib.HF_DeepCopy(clonecard, card)

		num := 3
		for i := 0; i < len(clonecard.Card1); {
			if clonecard.Card1[i] == _card {
				copy(clonecard.Card1[i:], clonecard.Card1[i+1:])
				clonecard.Card1 = clonecard.Card1[:len(clonecard.Card1)-1]
				num--
			} else {
				i++
			}
			if num <= 0 {
				break
			}
		}

		if num > 0 {
			return 0
		}

		ptype := MahIsHu(clonecard, ci)

		if ptype {
			return _card
		}
	}

	return 0
}

///////////////////////////////////////////////////////////
//! 是否7对
func MahIsDui(_card LstCard) (bool, int) {
	lst := _card

	if len(lst) != 14 {
		return false, 0
	}

	tmp := make(map[int]int)
	for i := 0; i < len(lst); i++ {
		tmp[lst[i]]++
	}

	num := 0
	for _, value := range tmp {
		if value%2 != 0 {
			return false, 0
		}
		if value/2 == 2 { //记录有几个杠
			num++
		}
	}

	return true, num
}

//!是否十三乱
func MashIsluan(_card LstCard) bool {
	if len(_card) != 14 {
		return false
	}
	tmp := make(map[int]int)
	for _, all_card := range _card {
		tmp[all_card]++
	}

	for _, all_card := range _card {
		if tmp[all_card] >= 2 {
			return false
		}

		if all_card/10 >= 0 && all_card/10 <= 2 {
			if all_card%10 > 2 { //左右两边不允许有牌，不能成顺子
				if tmp[all_card-1] > 0 || tmp[all_card-2] > 0 {
					return false
				}
			}
			if all_card%10 < 9 {
				if tmp[all_card+1] > 0 || tmp[all_card+2] > 0 {
					return false
				}
			}
		}
	}

	return true
}

//是否七星归位
func MashIsQXGW(_card LstCard) bool {
	if MashIsluan(_card) {
		var qi int
		for _, all_card := range _card {
			if all_card/10 >= 3 {
				qi++
			}
		}
		if qi == 7 {
			return true
		}
	}
	return false
}

//! 是否是碰碰胡
func MahIsPPH(card *Mah_Card, _card LstCard) (bool, int, int) {
	lst := _card

	num := 0
	num += len(card.CardP)
	num += len(card.CardMG)
	num += len(card.CardAG)
	num += len(card.CardCG)

	tmp := make(map[int]int)
	for i := 0; i < len(lst); i++ {
		tmp[lst[i]]++
	}

	d_num := 0
	jiang := 0
	kan := 0
	for key, value := range tmp {
		if value >= 3 {
			num++
			kan++
		} else if value == 2 {
			d_num++
			jiang = key
		}
	}

	return num == 4 && d_num == 1, jiang, kan
}

//! 是否胡(卡五星)
func MahIsHuByKWX(card *Mah_Card, _card int) (int, int, int) {
	var lst LstCard
	lib.HF_DeepCopy(&lst, &card.Card1)
	lst = append(lst, _card)
	lst = append(lst, card.CardL...)

	d7, num := MahIsDui(lst)
	if d7 {
		switch num {
		case 1:
			return TYPE_HU_7DUI1, 0, 0
		case 2:
			return TYPE_HU_7DUI2, 0, 0
		case 3:
			return TYPE_HU_7DUI3, 0, 0
		default:
			return TYPE_HU_7DUI, 0, 0
		}
	}

	ok, jiang, kan := MahIsPPH(card, lst)
	if ok {
		return TYPE_HU_PPH, jiang, kan
	}

	//! 将手牌转换成2维数组
	var handCards [4][10]int
	for i := 0; i < len(lst); i++ {
		handCards[lst[i]/10][lst[i]%10]++
		handCards[lst[i]/10][0]++
	}

	if _card%10 == 5 { //! 判断是否卡五星
		if handCards[_card/10][4] != 0 && handCards[_card/10][6] != 0 { //! 前后牌都有
			handCards[_card/10][4]--
			handCards[_card/10][6]--
			handCards[_card/10][5]--
			handCards[_card/10][0] -= 3
			ok, jiang, kan = lib.MahIsHuCard(handCards)
			if ok {
				return TYPE_HU_KWX, jiang, kan
			}
			handCards[_card/10][4]++
			handCards[_card/10][6]++
			handCards[_card/10][5]++
			handCards[_card/10][0] += 3
		}
	}
	ok, jiang, kan = lib.MahIsHuCard(handCards)
	if ok {
		return TYPE_HU_PH, jiang, kan
	}

	return 0, 0, 0
}

//! 是否胡(血战到底)
func MahIsHuByXZDD(card *Mah_Card, _card int) (int, int) {
	var lst LstCard
	lib.HF_DeepCopy(&lst, &card.Card1)
	lst = append(lst, _card)
	lst = append(lst, card.CardL...)

	d7, num := MahIsDui(lst)
	if d7 {
		if num > 0 {
			return TYPE_HU_7DUI1, 0
		} else {
			return TYPE_HU_7DUI, 0
		}
	}

	ok, jiang, _ := MahIsPPH(card, lst)
	if ok {
		return TYPE_HU_PPH, jiang
	}

	//! 将手牌转换成2维数组
	var handCards [4][10]int
	for i := 0; i < len(lst); i++ {
		handCards[lst[i]/10][lst[i]%10]++
		handCards[lst[i]/10][0]++
	}

	ok, jiang, _ = lib.MahIsHuCard(handCards)
	if ok {
		return TYPE_HU_PH, jiang
	}

	return 0, 0
}

//! 是否胡 (上虞花)
func MahIsHuByRazzSYH(card *Mah_Card, _card int, razz int, _type int) (bool, []int, int, int) {
	var lst LstCard
	lib.HF_DeepCopy(&lst, &card.Card1)
	lst = append(lst, _card)

	abs := make([]int, 0)
	if _type == 1 {
		ok, jiang, kan := MahIsHuLstCardSYH(lst)
		if ok {
			return true, abs, jiang, kan
		}
	}

	razznum := 0
	//! 将癞子拿出
	for i := 0; i < len(lst); {
		if lst[i] == razz {
			razznum++
			copy(lst[i:], lst[i+1:])
			lst = lst[:len(lst)-1]
		} else {
			i++
		}
	}

	if razznum == 0 { //!
		return false, abs, 0, 0
	}

	return MahIsHuToRazzSYH(card, lst, razznum, 0, abs, _type)
}

func MahIsHuToRazzSYH(card *Mah_Card, _lst LstCard, razznum int, step int, abs []int, _type int) (bool, []int, int, int) {
	if step == razznum {
		if _type == 1 {
			ok, jiang, kan := MahIsHuLstCardSYH(_lst)
			return ok, abs, jiang, kan
		}
		if _type == 6 {
			ok, jiang, _ := MahIsPPH(card, _lst)
			return ok, abs, jiang, 0
		}
	} else {
		step++
		for i := 1; i <= 9; i++ {
			var lst LstCard
			lib.HF_DeepCopy(&lst, &_lst)
			lst = append(lst, i)

			var _abs []int
			lib.HF_DeepCopy(&_abs, &abs)
			_abs = append(_abs, i)
			ok, a, jiang, kan := MahIsHuToRazzSYH(card, lst, razznum, step, _abs, _type)
			if ok {
				return true, a, jiang, kan
			}
		}
		for i := 11; i <= 19; i++ {
			var lst LstCard
			lib.HF_DeepCopy(&lst, &_lst)
			lst = append(lst, i)

			var _abs []int
			lib.HF_DeepCopy(&_abs, &abs)
			_abs = append(_abs, i)
			ok, a, jiang, kan := MahIsHuToRazzSYH(card, lst, razznum, step, _abs, _type)
			if ok {
				return true, a, jiang, kan
			}
		}
		for i := 21; i <= 29; i++ {
			var lst LstCard
			lib.HF_DeepCopy(&lst, &_lst)
			lst = append(lst, i)

			var _abs []int
			lib.HF_DeepCopy(&_abs, &abs)
			_abs = append(_abs, i)
			ok, a, jiang, kan := MahIsHuToRazzSYH(card, lst, razznum, step, _abs, _type)
			if ok {
				return true, a, jiang, kan
			}
		}
		for i := 31; i <= 34; i++ {
			var lst LstCard
			lib.HF_DeepCopy(&lst, &_lst)
			lst = append(lst, i)

			var _abs []int
			lib.HF_DeepCopy(&_abs, &abs)
			_abs = append(_abs, i)
			ok, a, jiang, kan := MahIsHuToRazzSYH(card, lst, razznum, step, _abs, _type)
			if ok {
				return true, a, jiang, kan
			}
		}
	}

	return false, abs, 0, 0
}

func MahIsHuLstCardSYH(lst LstCard) (bool, int, int) {
	//! 将手牌转换成2维数组
	var handCards [4][10]int
	for i := 0; i < len(lst); i++ {
		handCards[lst[i]/10][lst[i]%10]++
		handCards[lst[i]/10][0]++
	}

	ok, jiang, kan := lib.MahIsHuCard(handCards) //胡   将   看张数目
	if ok {
		return true, jiang, kan
	}

	return false, jiang, kan
}

////////////////////////////上虞花

//! 是否胡(癞子)
func MahIsHuByRazz(card *Mah_Card, _card int, razz int) (bool, []int) {
	var lst LstCard
	lib.HF_DeepCopy(&lst, &card.Card1)
	lst = append(lst, _card)
	//lst = append(lst, card.CardL...)

	abs := make([]int, 0)

	if MahIsHuLstCardBIAO(lst) {
		return true, abs
	}

	razznum := 0
	//! 将癞子拿出
	for i := 0; i < len(lst); {
		if lst[i] == razz {
			razznum++
			copy(lst[i:], lst[i+1:])
			lst = lst[:len(lst)-1]
		} else {
			i++
		}
	}

	if razznum == 0 { //!
		return false, abs
	}

	//	for i := 1; i <= razznum; i++ {
	//		ok, abs := MahIsHuToRazz(lst, i, 0, abs)
	//		if ok {
	//			return ok, abs
	//		}
	//		if i > 1 {
	//			return false, nil
	//		}
	//	}

	return MahIsHuToRazz(lst, razznum, 0, abs)
}

func MahIsHuToRazz(_lst LstCard, razznum int, step int, abs []int) (bool, []int) {
	if step == razznum {
		return MahIsHuLstCardBIAO(_lst), abs
	} else {
		step++
		for i := 1; i <= 9; i++ {
			var lst LstCard
			lib.HF_DeepCopy(&lst, &_lst)
			lst = append(lst, i)

			var _abs []int
			lib.HF_DeepCopy(&_abs, &abs)
			_abs = append(_abs, i)
			ok, a := MahIsHuToRazz(lst, razznum, step, _abs)
			if ok {
				return true, a
			}
		}
		for i := 11; i <= 19; i++ {
			var lst LstCard
			lib.HF_DeepCopy(&lst, &_lst)
			lst = append(lst, i)

			var _abs []int
			lib.HF_DeepCopy(&_abs, &abs)
			_abs = append(_abs, i)
			ok, a := MahIsHuToRazz(lst, razznum, step, _abs)
			if ok {
				return true, a
			}
		}
		for i := 21; i <= 29; i++ {
			var lst LstCard
			lib.HF_DeepCopy(&lst, &_lst)
			lst = append(lst, i)

			var _abs []int
			lib.HF_DeepCopy(&_abs, &abs)
			_abs = append(_abs, i)
			ok, a := MahIsHuToRazz(lst, razznum, step, _abs)
			if ok {
				return true, a
			}
		}
		for i := 31; i <= 34; i++ {
			var lst LstCard
			lib.HF_DeepCopy(&lst, &_lst)
			lst = append(lst, i)

			var _abs []int
			lib.HF_DeepCopy(&_abs, &abs)
			_abs = append(_abs, i)
			ok, a := MahIsHuToRazz(lst, razznum, step, _abs)
			if ok {
				return true, a
			}
		}
	}

	return false, abs
}

//! 是否胡
func MahIsHu(card *Mah_Card, _card int) bool {
	var lst LstCard
	lib.HF_DeepCopy(&lst, &card.Card1)
	lst = append(lst, _card)
	if len(card.Want) > 0 {
		lst = append(lst, card.CardL...)
	}

	return MahIsHuLstCard(lst)
}

func MahIsHuLstCard(lst LstCard) bool {
	d7, _ := MahIsDui(lst)
	if d7 {
		return true
	}

	//! 将手牌转换成2维数组
	var handCards [4][10]int
	for i := 0; i < len(lst); i++ {
		handCards[lst[i]/10][lst[i]%10]++
		handCards[lst[i]/10][0]++
	}

	ok, _, _ := lib.MahIsHuCard(handCards)
	if ok {
		return true
	}

	return false
}

///////////////////////////////////////////////////////////
//! 是否清一色
func MahIsAllColor(card *Mah_Card, _card int) bool {
	tmp := _card / 10
	for i := 0; i < len(card.Card1); i++ {
		if card.Card1[i]/10 != tmp {
			return false
		}
	}
	for i := 0; i < len(card.CardL); i++ {
		if card.CardL[i]/10 != tmp {
			return false
		}
	}
	for i := 0; i < len(card.CardP); i++ {
		if card.CardP[i]/10 != tmp {
			return false
		}
	}
	for i := 0; i < len(card.CardMG); i++ {
		if card.CardMG[i]/10 != tmp {
			return false
		}
	}
	for i := 0; i < len(card.CardAG); i++ {
		if card.CardAG[i]/10 != tmp {
			return false
		}
	}
	for i := 0; i < len(card.CardCG); i++ {
		if card.CardCG[i]/10 != tmp {
			return false
		}
	}

	return true
}

///////////////////////////////////////////////////////////
//! 是否清一色   //上虞花什么的
func MahIsAllColorGSQ(card *Mah_Card, _card int) bool {
	tmp := _card / 10
	for i := 0; i < len(card.Card1); i++ {
		if card.Card1[i]/10 != tmp {
			return false
		}
	}
	for i := 0; i < len(card.CardC); i++ {
		for j := 0; j < len(card.CardC[i]); j++ {
			if card.CardC[i][j]/10 != tmp {
				return false
			}
		}
	}
	for i := 0; i < len(card.CardP); i++ {
		if card.CardP[i]/10 != tmp {
			return false
		}
	}
	for i := 0; i < len(card.CardMG); i++ {
		if card.CardMG[i]/10 != tmp {
			return false
		}
	}
	for i := 0; i < len(card.CardAG); i++ {
		if card.CardAG[i]/10 != tmp {
			return false
		}
	}
	for i := 0; i < len(card.CardCG); i++ {
		if card.CardCG[i]/10 != tmp {
			return false
		}
	}

	return true
}

//! 是否混一色
func MahIsAll2Color(card *Mah_Card, _card int) bool {
	tmp := _card / 10
	for i := 0; i < len(card.Card1); i++ {
		if card.Card1[i]/10 != tmp && card.Card1[i]/10 != 3 {
			return false
		}
	}
	for i := 0; i < len(card.CardL); i++ {
		if card.CardL[i]/10 != tmp && card.Card1[i]/10 != 3 {
			return false
		}
	}
	for i := 0; i < len(card.CardP); i++ {
		if card.CardP[i]/10 != tmp && card.Card1[i]/10 != 3 {
			return false
		}
	}
	for i := 0; i < len(card.CardMG); i++ {
		if card.CardMG[i]/10 != tmp && card.Card1[i]/10 != 3 {
			return false
		}
	}
	for i := 0; i < len(card.CardAG); i++ {
		if card.CardAG[i]/10 != tmp && card.Card1[i]/10 != 3 {
			return false
		}
	}
	for i := 0; i < len(card.CardCG); i++ {
		if card.CardCG[i]/10 != tmp && card.Card1[i]/10 != 3 {
			return false
		}
	}

	return true
}

//! 是否混一色  、、上虞花
func MahIsAll2ColorGSQ(card *Mah_Card, _card int) bool {
	tmp := _card / 10
	for i := 0; i < len(card.Card1); i++ {
		if card.Card1[i]/10 != tmp && card.Card1[i]/10 != 3 {
			return false
		}
	}

	for i := 0; i < len(card.CardC); i++ {
		for j := 0; j < len(card.CardC[i]); j++ {
			if card.CardC[i][j]/10 != tmp && card.CardC[i][j] != 3 {
				return false
			}
		}
	}
	for i := 0; i < len(card.CardP); i++ {
		if card.CardP[i]/10 != tmp && card.Card1[i]/10 != 3 {
			return false
		}
	}
	for i := 0; i < len(card.CardMG); i++ {
		if card.CardMG[i]/10 != tmp && card.Card1[i]/10 != 3 {
			return false
		}
	}
	for i := 0; i < len(card.CardAG); i++ {
		if card.CardAG[i]/10 != tmp && card.Card1[i]/10 != 3 {
			return false
		}
	}
	for i := 0; i < len(card.CardCG); i++ {
		if card.CardCG[i]/10 != tmp && card.Card1[i]/10 != 3 {
			return false
		}
	}

	return true
}

//! 是否将对
func MahIsJong(card *Mah_Card, _card int) bool {
	for i := 0; i < len(card.CardP); i++ {
		if card.CardP[i]%10 != 2 && card.CardP[i]%10 != 5 && card.CardP[i]%10 != 8 {
			return false
		}
	}

	for i := 0; i < len(card.CardAG); i++ {
		if card.CardAG[i]%10 != 2 && card.CardAG[i]%10 != 5 && card.CardAG[i]%10 != 8 {
			return false
		}
	}

	for i := 0; i < len(card.CardMG); i++ {
		if card.CardMG[i]%10 != 2 && card.CardMG[i]%10 != 5 && card.CardMG[i]%10 != 8 {
			return false
		}
	}

	for i := 0; i < len(card.CardCG); i++ {
		if card.CardCG[i]%10 != 2 && card.CardCG[i]%10 != 5 && card.CardCG[i]%10 != 8 {
			return false
		}
	}

	for i := 0; i < len(card.Card1); i++ {
		if card.Card1[i]%10 != 2 && card.Card1[i]%10 != 5 && card.Card1[i]%10 != 8 {
			return false
		}
	}

	return true
}

//! 是否门清
func MahIsClear(card *Mah_Card) bool {
	if len(card.CardP) > 0 {
		return false
	}

	if len(card.CardAG) > 0 {
		return false
	}

	if len(card.CardCG) > 0 {
		return false
	}

	if len(card.CardMG) > 0 {
		return false
	}

	return true
}

//! 是否门清
func MahIsClearGSQ(card *Mah_Card) bool {
	if len(card.CardP) > 0 {
		return false
	}

	if len(card.CardC) > 0 {
		return false
	}

	if len(card.CardAG) > 0 {
		return false
	}

	if len(card.CardCG) > 0 {
		return false
	}

	if len(card.CardMG) > 0 {
		return false
	}

	return true
}

//! 是否中张
func MahIsMidle(card *Mah_Card) bool {
	for i := 0; i < len(card.CardP); i++ {
		if card.CardP[i]%10 == 1 || card.CardP[i]%10 == 9 {
			return false
		}
	}

	for i := 0; i < len(card.CardAG); i++ {
		if card.CardAG[i]%10 == 1 || card.CardAG[i]%10 == 9 {
			return false
		}
	}

	for i := 0; i < len(card.CardCG); i++ {
		if card.CardCG[i]%10 == 1 || card.CardCG[i]%10 == 9 {
			return false
		}
	}

	for i := 0; i < len(card.CardMG); i++ {
		if card.CardMG[i]%10 == 1 || card.CardMG[i]%10 == 9 {
			return false
		}
	}

	for i := 0; i < len(card.Card1); i++ {
		if card.Card1[i]%10 == 1 || card.Card1[i]%10 == 9 {
			return false
		}
	}

	return true
}

//! 有几个4个
func MahHas4Num(card *Mah_Card, _card int) (int, int) {
	var lst LstCard
	lib.HF_DeepCopy(&lst, &card.Card1)
	lst = append(lst, _card)
	lst = append(lst, card.CardL...)

	tmp := make(map[int]int)
	for i := 0; i < len(lst); i++ {
		tmp[lst[i]]++
	}
	for i := 0; i < len(card.CardP); i++ {
		tmp[card.CardP[i]] += 3
	}

	ming := 0
	an := 0
	for key, value := range tmp {
		if value >= 4 {
			find := false
			for i := 0; i < len(card.CardP); i++ {
				if card.CardP[i] == key {
					find = true
					break
				}
			}
			if !find {
				an++
			} else {
				ming++
			}
		}
	}

	return ming, an
}

//! 是否明四归
func MahIsMing4(card *Mah_Card, _card int) bool {
	for i := 0; i < len(card.CardP); i++ {
		if card.CardP[i] == _card {
			return true
		}
	}

	return false
}

//! 是否暗四归
func MahIsAn4(card *Mah_Card, _card int) bool {
	num := 0
	for i := 0; i < len(card.Card1); i++ {
		if card.Card1[i] == _card {
			num++
		}
	}
	for i := 0; i < len(card.CardL); i++ {
		if card.CardL[i] == _card {
			num++
		}
	}

	return num >= 3
}

//! 是否手抓一
func MahIsZOne(card *Mah_Card, _card int) bool {
	return len(card.Card1)+len(card.CardL) == 1
}

//! 元数量
func MahYuanNum(card *Mah_Card, _card int) int {
	num := 0
	for i := 0; i < len(card.CardP); i++ {
		if card.CardP[i]/10 == 3 {
			num++
		}
	}
	for i := 0; i < len(card.CardMG); i++ {
		if card.CardMG[i]/10 == 3 {
			num++
		}
	}
	for i := 0; i < len(card.CardAG); i++ {
		if card.CardAG[i]/10 == 3 {
			num++
		}
	}
	for i := 0; i < len(card.CardCG); i++ {
		if card.CardCG[i]/10 == 3 {
			num++
		}
	}

	tmp := make(map[int]int)
	for i := 0; i < len(card.Card1); i++ {
		tmp[card.Card1[i]]++
	}
	for i := 0; i < len(card.CardL); i++ {
		tmp[card.CardL[i]]++
	}
	tmp[_card]++

	for key, value := range tmp {
		if key/10 == 3 && value >= 3 {
			num++
		}
	}

	return num
}

//! 坎数量
func MahKanNum(card *Mah_Card, _card int) int {
	var lst LstCard
	lib.HF_DeepCopy(&lst, &card.Card1)
	if _card != 0 {
		lst = append(lst, _card)
	}
	lst = append(lst, card.CardL...)

	num := len(card.CardAG)
	num += len(card.CardCG)
	num += len(card.CardMG)
	tmp := make(map[int]int)
	for i := 0; i < len(lst); i++ {
		tmp[lst[i]]++
	}
	for _, value := range tmp {
		if value >= 3 {
			num++
		}
	}

	return num
}

//! 某一花色最多张数
func MahMaxNum(card *Mah_Card, _card int) int {
	var lst LstCard
	lib.HF_DeepCopy(&lst, &card.Card1)
	if _card != 0 {
		lst = append(lst, _card)
	}
	lst = append(lst, card.CardL...)

	tmp := make(map[int]int)

	for i := 0; i < len(lst); i++ {
		tmp[lst[i]/10]++
	}
	for i := 0; i < len(card.CardP); i++ {
		tmp[card.CardP[i]/10] += 3
	}
	for i := 0; i < len(card.CardAG); i++ {
		tmp[card.CardAG[i]/10] += 4
	}
	for i := 0; i < len(card.CardCG); i++ {
		tmp[card.CardCG[i]/10] += 4
	}
	for i := 0; i < len(card.CardMG); i++ {
		tmp[card.CardMG[i]/10] += 4
	}

	maxnum := 0
	for _, value := range tmp {
		if value > maxnum {
			maxnum = value
		}
	}

	return maxnum
}

//判断是否单吊 安庆麻将
func MahIsDanDiaoAQMJ(card *Mah_Card, _card int) bool {
	ptype := MahIsHuBIAO(card, _card)

	if ptype {
		for _, all_card := range card.Card1 {
			if all_card/10 == 3 {
				continue
			}
			if all_card != _card {
				Ishu := MahIsHuBIAO(card, all_card)
				if Ishu {
					return false
				}
			}

			if all_card%10 != 1 && all_card-1 != _card { //不是一 的情况
				Ishu := MahIsHuBIAO(card, all_card-1)
				if Ishu {
					return false
				}
			}
			if all_card%10 != 9 && all_card+1 != _card { //不是9 的情况
				Ishu := MahIsHuBIAO(card, all_card+1)
				if Ishu {
					return false
				}
			}
		}
	} else {
		return false
	}

	return true
}

//以下胡牌针对南昌麻将的

//! 是否胡(癞子)  //_type  胡牌类型    0 标准胡牌  1 七对 2 十三乱 3 七星 4碰碰胡
func MahIsHuByRazzNC(card *Mah_Card, _card int, razz1, razz2 int, _type int, selfMo bool) (bool, []int, int, int) {
	var lst LstCard
	lib.HF_DeepCopy(&lst, &card.Card1)
	lst = append(lst, _card)
	if _type == 4 {
		lst = append(lst, card.CardL...)
	}
	abs := make([]int, 0)

	if _type == 0 {
		ok, jiang, kan := MahIsHuLstCardNC(lst)
		if ok {
			return ok, abs, jiang, kan
		}
	}

	del_razz := false
	if !selfMo {
		del_razz = true
	}

	razznum := 0
	//! 将癞子拿出
	for i := 0; i < len(lst); {
		if del_razz && lst[i] == _card && (lst[i] == razz1 || lst[i] == razz2) {
			del_razz = false
			i++
			continue
		}
		if lst[i] == razz1 || lst[i] == razz2 {
			razznum++
			copy(lst[i:], lst[i+1:])
			lst = lst[:len(lst)-1]
		} else {
			i++
		}
	}

	if razznum == 0 { //!
		return false, abs, 0, 0
	}

	return MahIsHuToRazzNC(card, lst, razznum, 0, abs, _type)
}

func MahIsHuToRazzNC(card *Mah_Card, _lst LstCard, razznum int, step int, abs []int, _type int) (bool, []int, int, int) {
	if step == razznum {
		if _type == 0 { //平胡
			ishu, jiang, kan := MahIsHuLstCardNC(_lst)
			return ishu, abs, jiang, kan
		}
		if _type == 1 { //七对
			ishu, _ := MahIsDui(_lst)
			return ishu, abs, 0, 0
		}
		if _type == 2 { //十三乱
			ishu := MashIsluan(_lst)
			return ishu, abs, 0, 0
		}
		if _type == 3 { //七星十三乱
			ishu := MashIsQXGW(_lst)
			return ishu, abs, 0, 0
		}
		if _type == 4 { //碰碰胡
			ok, _, _ := MahIsPPH(card, _lst)
			return ok, abs, 0, 0
		}
	} else {
		step++
		for i := 1; i <= 9; i++ {
			var lst LstCard
			lib.HF_DeepCopy(&lst, &_lst)
			lst = append(lst, i)

			var _abs []int
			lib.HF_DeepCopy(&_abs, &abs)
			_abs = append(_abs, i)
			ok, a, jiang, kan := MahIsHuToRazzNC(card, lst, razznum, step, _abs, _type)
			if ok {
				return true, a, jiang, kan
			}
		}
		for i := 11; i <= 19; i++ {
			var lst LstCard
			lib.HF_DeepCopy(&lst, &_lst)
			lst = append(lst, i)

			var _abs []int
			lib.HF_DeepCopy(&_abs, &abs)
			_abs = append(_abs, i)
			ok, a, jiang, kan := MahIsHuToRazzNC(card, lst, razznum, step, _abs, _type)
			if ok {
				return true, a, jiang, kan
			}
		}
		for i := 21; i <= 29; i++ {
			var lst LstCard
			lib.HF_DeepCopy(&lst, &_lst)
			lst = append(lst, i)

			var _abs []int
			lib.HF_DeepCopy(&_abs, &abs)
			_abs = append(_abs, i)
			ok, a, jiang, kan := MahIsHuToRazzNC(card, lst, razznum, step, _abs, _type)
			if ok {
				return true, a, jiang, kan
			}
		}
		if _type != 3 {
			for i := 31; i <= 37; i++ {
				var lst LstCard
				lib.HF_DeepCopy(&lst, &_lst)
				lst = append(lst, i)

				var _abs []int
				lib.HF_DeepCopy(&_abs, &abs)
				_abs = append(_abs, i)
				ok, a, jiang, kan := MahIsHuToRazzNC(card, lst, razznum, step, _abs, _type)
				if ok {
					return true, a, jiang, kan
				}
			}
		}
	}

	return false, abs, 0, 0
}

//! 是否胡(针对南昌麻将类型，东西南北算顺子，中发白算顺子)
func MahIsHuNC(card *Mah_Card, _card int) (bool, int, int) {
	var lst LstCard
	lib.HF_DeepCopy(&lst, &card.Card1)
	lst = append(lst, _card)
	//	if len(card.Want) > 0 {
	//		lst = append(lst, card.CardL...)
	//	}

	return MahIsHuLstCardNC(lst)
}

func MahIsHuLstCardNC(lst LstCard) (bool, int, int) {
	//! 将手牌转换成2维数组
	var handCards [4][10]int
	for i := 0; i < len(lst); i++ {
		handCards[lst[i]/10][lst[i]%10]++
		handCards[lst[i]/10][0]++
	}

	ok, jiang, kan := MahIsHuCardNC(handCards)
	if ok {
		lib.GetLogMgr().Output(lib.LOG_INFO, "胡的牌型：", lst)
		return true, jiang, kan
	}

	return false, 0, 0
}

func MahIsHuCardNC(handCards [4][10]int) (bool, int, int) {
	douleLeft := 0
	singleLeft := 0
	for i := 0; i < 4; i++ {
		if handCards[i][0]%3 == 2 {
			douleLeft++
		}
		if handCards[i][0]%3 == 1 {
			singleLeft++
		}
	}

	if douleLeft > 1 || singleLeft > 1 {
		return false, 0, 0
	}

	//! 获取真正有可能做将牌的牌
	jiangVec := make([]int, 0)
	for i := 0; i < 4; i++ {
		for j := 1; j < 10; j++ {
			if handCards[i][j] >= 2 {
				jiangVec = append(jiangVec, i*10+j)
			}
		}
	}

	//! 轮询去掉将牌
	for _, card := range jiangVec {
		var tmpCards [4][10]int
		lib.HF_DeepCopy(&tmpCards, &handCards)
		tmpCards[card/10][card%10] -= 2
		tmpCards[card/10][0] -= 2

		//! 拆解3n牌结构
		tmp, kan := Remove3NCardsNC(tmpCards)

		if tmp[0][0]+tmp[1][0]+tmp[2][0]+tmp[3][0] == 0 {
			return true, card, kan //胡   将   看张数目
		}
	}

	return false, 0, 0
}

func Remove3NCardsNC(cards [4][10]int) ([4][10]int, int) {
	num := 0

	//! 去掉风牌的刻子
	for i := 1; i <= 7; i++ {
		if cards[3][i] <= 0 {
			continue
		}
		if cards[3][i] == 3 {
			cards[3][i] -= 3
			cards[3][0] -= 3
			num++
		}
	}
	//去掉风牌的顺子
	for i := 1; i <= 4; i++ {
		tmpCount := cards[3][i]
		f_shun := 1
		for k := 0; k < tmpCount; k++ {
			for j := 0; j <= 4; j++ {
				if j == i {
					continue
				}

				if cards[3][j] > 0 {
					cards[3][j]--
					f_shun++
				}

				if f_shun == 3 {
					cards[3][0] -= 3
					cards[3][i]--
					break
				}
			}
			if f_shun < 3 { //不符合
				return cards, num
			}
		}
	}
	if cards[3][5] > 0 { //中发白顺子
		tmpCount := cards[3][5]
		for k := 0; k < tmpCount; k++ {
			if cards[3][6] > 0 && cards[3][7] > 0 {
				cards[3][0] -= 3
				cards[3][5]--
				cards[3][6]--
				cards[3][7]--
			} else {
				return cards, num
			}
		}
	}

	//! 去掉1~7的刻子和顺子
	for i := 0; i < 3; i++ {
		for j := 1; j <= 7; j++ {
			if cards[i][j] <= 0 {
				continue
			}

			//! 去掉刻子
			if cards[i][j] >= 3 {
				cards[i][j] -= 3
				cards[i][0] -= 3
				num++
			}

			// 先去掉顺子
			tmpCount := cards[i][j]
			for k := 0; k < tmpCount; k++ {
				if cards[i][j+1] > 0 && cards[i][j+2] > 0 {
					cards[i][0] -= 3
					cards[i][j]--
					cards[i][j+1]--
					cards[i][j+2]--
				} else {
					break
				}
			}

			if cards[i][j] > 0 {
				return cards, num
			}
		}
	}

	//! 去掉8~9的刻子
	for i := 0; i < 3; i++ {
		for j := 8; j <= 9; j++ {
			if cards[i][j] <= 0 {
				continue
			}
			if cards[i][j] == 3 {
				cards[i][0] -= 3
				cards[i][j] -= 3
				num++
			} else {
				return cards, num
			}
		}
	}
	return cards, num
}

/////////////////////////////////////////////////////////////////////////////////

/////////////标准胡牌
//! 是否胡
func MahIsHuBIAO(card *Mah_Card, _card int) bool {
	var lst LstCard
	lib.HF_DeepCopy(&lst, &card.Card1)
	lst = append(lst, _card)
	//	if len(card.Want) > 0 {
	//		lst = append(lst, card.CardL...)
	//	}

	for i := 0; i < len(lst); i++ {
		if lst[i]/10 >= 4 {
			return false
		}
	}
	return MahIsHuLstCardBIAO(lst)
}

func MahIsHuLstCardBIAO(lst LstCard) bool {
	//! 将手牌转换成2维数组
	var handCards [4][10]int
	for i := 0; i < len(lst); i++ {
		handCards[lst[i]/10][lst[i]%10]++
		handCards[lst[i]/10][0]++
	}

	ok, _, _ := lib.MahIsHuCard(handCards)
	if ok {
		return true
	}

	return false
}

///////////////////////////////////////////////
//！推对子结果
func GetTDZResult(card []int) (int, int) { //牌值牌型
	if len(card) != 2 {
		return -1, -1
	}
	value, px := -1, -1
	for i := 0; i < len(card); i++ {
		value += card[i] % 10
	}
	value = value % 10
	if card[0] == card[1] {
		px = 2
	} else {
		px = 1
	}
	return value, px
}

//! 得到推筒子的点数和倍数
func GetTTZResult(card []int) (int, int) {
	if card[0] == 37 && card[1] == 37 { //! 至尊
		return 10000, 20
	}

	if card[0] == card[1] { //! 豹子
		return 9990 + card[0]%10, 15
	}

	if card[0] == 12 && card[1] == 18 || card[0] == 18 && card[1] == 12 {
		return 9980, 12
	}

	value := float32(0)
	for i := 0; i < len(card); i++ {
		if card[i] == 37 {
			value += 0.5
		} else {
			value += float32(card[i] % 10)
		}
	}

	if value >= 10 {
		value = float32(int(value) % 10)
	}

	if value >= 9 {
		return int(value * 10), 9
	} else if value >= 8 {
		return int(value * 10), 8
	} else if value >= 7 {
		return int(value * 10), 7
	} else if value >= 6 {
		return int(value * 10), 6
	} else if value >= 5 {
		return int(value * 10), 5
	} else if value >= 4 {
		return int(value * 10), 4
	} else if value >= 3 {
		return int(value * 10), 3
	} else if value >= 2 {
		return int(value * 10), 2
	} else {
		return int(value * 10), 1
	}
}

//! 得到推筒子的点数和倍数
func GetTTZResult1(card []int) (int, int) {
	if card[0] == 37 && card[1] == 37 { //! 至尊
		return 10000, 5
	}

	if card[0] == card[1] { //! 豹子
		return 9990 + card[0]%10, 4
	}

	if card[0] == 12 && card[1] == 18 || card[0] == 18 && card[1] == 12 {
		return 9980, 3
	}

	value := float32(0)
	for i := 0; i < len(card); i++ {
		if card[i] == 37 {
			value += 0.5
		} else {
			value += float32(card[i] % 10)
		}
	}

	if value >= 10 {
		value = float32(int(value) % 10)
	}

	if value >= 8 {
		return int(value * 10), 2
	}

	return int(value * 10), 1
}
