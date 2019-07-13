package gameserver

import (
	"fmt"
	"lib"
	"math"
	"staticfunc"
	"time"
)

//! param1
//! 000000000
const TYPE_KWX = 0  //! 卡五星番数 0两番 1四番
const TYPE_PPH = 1  //! 碰碰胡/杠上炮番数 0两番 1四番
const TYPE_DL = 2   //! 0必须亮牌自摸 1自摸即可
const TYPE_ZDFS = 3 //! 最大番数 0八番 1十六番
const TYPE_PIAO = 4 //! 是否加飘 0不加飘 1自选加飘
const TYPE_SK = 5   //! 是否数坎 0不数坎 1数坎
const TYPE_MM = 6   //! 是否买马 0不买马 1独马 2六马
const TYPE_PH = 7   //! 是否荒庄赔胡 0不赔胡 1赔胡
const TYPE_PD = 8   //! 频道 0半频道 1全频道
const TYPE_MAX = 9

//! param2
const TYPE_SL = 0   //! 上楼   0不上楼  1上楼
const TYPE_DLDF = 1 //! 对亮对番   0不对番  1对番
const TYPE_PQMM = 2 //! 跑恰摸马   0不  1有
const TYPE_MAX2 = 3

//! 卡五星作弊
type GameKWX_GoldInfo struct {
	Info []int `json:"info"`
}

type GameKWX_GoldGet struct {
	Card int `json:"card"`
}

//! 记录结构
type Rec_GameKWX_Info struct {
	Roomid int                      `json:"roomid"` //! 房间号
	Time   int64                    `json:"time"`   //! 记录时间
	Param1 int                      `json:"param1"`
	Param2 int                      `json:"param2"`
	Person []Son_Rec_GameKWX_Person `json:"person"`
	Step   []Son_Rec_GameKWX_Step   `json:"step"`
}
type Son_Rec_GameKWX_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
	Total int    `json:"total"`
}
type Son_Rec_GameKWX_Step struct { //当前动作结构
	Uid  int64 `json:"uid"`
	Type int   `json:"type"` //! 0摸牌 1出牌 2碰牌 3亮牌 4胡牌 5暗杠 6明杠 7擦杠
	Card []int `json:"card"`
}

type Msg_GameKWX_Info struct {
	Begin    bool               `json:"begin"`    //! 是否开始
	Info     []Son_GameKWX_Info `json:"info"`     //! 人的info
	Num      int                `json:"num"`      //! 剩余数量
	CurStep  int64              `json:"curstep"`  //! 这局谁出
	BefStep  int64              `json:"befstep"`  //! 上局谁出
	Ma       []int              `json:"ma"`       //! 买马
	Hz       bool               `json:"hz"`       //! 荒庄
	Sl       [2]int             `json:"sl"`       //! 上楼
	LastCard int                `json:"lastcard"` //! 最后的牌
}

type Son_GameKWX_Info struct {
	Uid      int64           `json:"uid"`
	Deal     bool            `json:"deal"`
	Piao     int             `json:"piao"`
	Card     Mah_Card        `json:"card"`
	So_Card  int             `json:"so_card"`  //! 牌分
	So_Gang  int             `json:"so_gang"`  //! 杠分
	So_Other int             `json:"so_other"` //! 其他分
	Total    int             `json:"total"`    //! 总分
	Ready    bool            `json:"ready"`
	Peng     int             `json:"peng"` //! 是否能碰
	Gang     int             `json:"gang"` //! 是否能杠
	Hu       int             `json:"hu"`   //! 是否能胡
	State    []Son_KWX_State `json:"state"`
	Num      int             `json:"num"`
}

//! 出牌
type Msg_GameKWX_Step struct {
	Uid  int64 `json:"uid"`  //! 哪个uid
	Card int   `json:"card"` //! 出的啥牌
}

//! 摸牌
type Msg_GameKWX_Draw struct {
	Uid  int64 `json:"uid"`  //! 哪个uid
	Card int   `json:"card"` //! 摸的啥牌
	Hu   int   `json:"hu"`
	Gang int   `json:"gang"`
}

//! 卡操作
type Msg_GameKWX_Total struct {
	Info []Son_GameKWX_Total `json:"info"` //!
}
type Son_GameKWX_Total struct {
	Uid   int64 `json:"uid"` //!
	Total int   `json:"total"`
}

//! 卡操作
type Msg_GameKWX_Operator struct {
	Hu   int `json:"hu"`   //! 能胡
	Peng int `json:"peng"` //! 能碰
	Gang int `json:"gang"` //! 能杠
}

//! 碰牌
type Msg_GameKWX_Peng struct {
	Uid  int64 `json:"uid"`  //! 哪个uid
	Card int   `json:"card"` //! 碰的啥牌
}

type Msg_GameKWX_Gang struct {
	Uid  int64 `json:"uid"`  //! 哪个uid
	Card int   `json:"card"` //! 杠的啥牌
	View bool  `json:"view"` //! 是否明杠
}

type Msg_GameKWX_Kill struct {
	Uid  int64 `json:"uid"`  //! 哪个uid
	Card int   `json:"card"` //! 杠的啥牌
	View bool  `json:"view"` //! 是否明杠
	Want []int `json:"want"`
}

//! 结算
type Msg_GameKWX_End struct {
	Hz       bool               `json:"hz"`
	Ma       []int              `json:"ma"`
	LastCard int                `json:"lastcard"`
	Info     []Son_GameKWX_Info `json:"info"`
}

//! 亮牌
type Msg_GameKWX_View struct {
	Uid   int64 `json:"uid"`
	CardL []int `json:"cardl"`
	Want  []int `json:"want"`
	Card  int   `json:"card"`
}

//! 预定牌
type Msg_GameKWX_Need struct {
	Card []int `json:"card"`
}

//! 房间结束
type Msg_GameKWX_Bye struct {
	Info []Son_GameKWX_Bye `json:"info"`
}
type Son_GameKWX_Bye struct {
	Uid   int64 `json:"uid"`
	ZM    int   `json:"zm"` //! 自摸
	JP    int   `json:"jp"` //! 接炮
	DP    int   `json:"dp"` //! 点炮
	AG    int   `json:"ag"` //! 暗杠
	MG    int   `json:"mg"` //! 明杠
	PI    int   `json:"pi"` //! 漂
	Score int   `json:"score"`
}

///////////////////////////////////////////////////////
type Son_KWX_State struct {
	Id    int `json:"id"`
	Score int `json:"score"`
}

type Game_KWX_Person struct {
	Uid      int64           `json:"uid"`
	Deal     bool            `json:"deal"`
	Piao     int             `json:"piao"`     //! 几飘
	So_Card  int             `json:"so_card"`  //! 牌分
	So_Gang  int             `json:"so_gang"`  //! 杠分
	So_Other int             `json:"so_other"` //! 其他分
	Total    int             `json:"total"`    //! 总分
	Peng     int             `json:"peng"`     //! 是否能碰
	Gang     int             `json:"gang"`     //! 是否能杠
	Hu       int             `json:"hu"`       //! 是否能胡
	Ready    bool            `json:"ready"`    //! 是否准备好
	Card     Mah_Card        `json:"card"`
	Kan      int             `json:"kan"` //! 坎
	ZM       int             `json:"zm"`  //! 自摸数量
	JP       int             `json:"jp"`  //! 接炮数量
	DP       int             `json:"dp"`  //! 点炮数量
	AG       int             `json:"ag"`  //! 暗杠数量
	MG       int             `json:"mg"`  //! 明杠数量
	State    []Son_KWX_State `json:"state"`
	NextCard int             `json:"nextcard"`
	NoHu     int             `json:"nohu"`

	need []int //! 预定的牌
	hz   int   //! 荒庄规则

}

func (self *Game_KWX_Person) Init() {
	self.Deal = false
	self.Peng = 0
	self.Gang = 0
	self.Hu = 0
	self.So_Card = 0
	self.So_Gang = 0
	self.So_Other = 0
	self.State = make([]Son_KWX_State, 0)
	self.NextCard = 0
	self.Card.Init()
}

func (self *Game_KWX_Person) AddState(id, score int) {
	if id != 4 && id != 10 { //! 4归除外
		for i := 0; i < len(self.State); i++ {
			if self.State[i].Id == id {
				self.State[i].Score = score
				return
			}
		}
	}

	self.State = append(self.State, Son_KWX_State{id, score})
}

type Game_KWX struct {
	PersonMgr []*Game_KWX_Person `json:"personmgr"`
	Mah       *MahMgr            `json:"mah"`      //! 剩余
	CurStep   int64              `json:"curstep"`  //! 这局谁出
	BefStep   int64              `json:"befstep"`  //! 上局谁出
	Winer     int64              `json:"winer"`    //! 上局谁赢
	Gang      int                `json:"gang"`     //! 当前局第几杠
	LastCard  int                `json:"lastcard"` //! 最后一张牌
	ViewUid   int64              `json:"viewuid"`  //! 第一家亮牌
	GangUid   int64              `json:"ganguid"`  //! 最近杠牌
	Record    *Rec_GameKWX_Info  `json:"record"`   //! 卡五星记录
	IsRecord  bool               `json:"isrecord"` //! 是否记录了
	Ma        []int              `json:"ma"`       //! 买马
	Hz        bool               `json:"hz"`       //! 荒庄
	SL        [2]int             `json:"sl"`       //! 上楼

	room *Room
}

func NewGame_KWX() *Game_KWX {
	game := new(Game_KWX)
	game.PersonMgr = make([]*Game_KWX_Person, 0)

	return game
}

func (self *Game_KWX) GetParam(_type int) int {
	return self.room.Param1 % int(math.Pow(10.0, float64(TYPE_MAX-_type))) / int(math.Pow(10.0, float64(TYPE_MAX-_type-1)))
}

func (self *Game_KWX) SetParam(_type int, value int) {
	high := self.room.Param1 / int(math.Pow(10.0, float64(TYPE_MAX-_type)))
	next := self.room.Param1 % int(math.Pow(10.0, float64(TYPE_MAX-_type-1)))
	self.room.Param1 = high*int(math.Pow(10.0, float64(TYPE_MAX-_type))) + value*int(math.Pow(10.0, float64(TYPE_MAX-_type-1))) + next
}

func (self *Game_KWX) GetParam2(_type int) int {
	return self.room.Param2 % int(math.Pow(10.0, float64(TYPE_MAX2-_type))) / int(math.Pow(10.0, float64(TYPE_MAX2-_type-1)))
}

func (self *Game_KWX) SetParam2(_type int, value int) {
	high := self.room.Param2 / int(math.Pow(10.0, float64(TYPE_MAX2-_type)))
	next := self.room.Param2 % int(math.Pow(10.0, float64(TYPE_MAX2-_type-1)))
	self.room.Param2 = high*int(math.Pow(10.0, float64(TYPE_MAX2-_type))) + value*int(math.Pow(10.0, float64(TYPE_MAX2-_type-1))) + next
}

func (self *Game_KWX) GetPerson(uid int64) *Game_KWX_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

//! 得到下一个uid
func (self *Game_KWX) GetNextUid() int64 {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid != self.CurStep {
			continue
		}

		if i+1 < len(self.PersonMgr) {
			return self.PersonMgr[i+1].Uid
		} else {
			return self.PersonMgr[0].Uid
		}
	}

	return 0
}

func (self *Game_KWX) OnInit(room *Room) {
	self.room = room

	if self.room.Type == GAMETYPE_KWX_XG { //! 孝感
		self.SetParam(TYPE_PH, 1)
		self.SetParam(TYPE_PD, 0)
		self.SetParam2(TYPE_SL, 0)
		self.SetParam2(TYPE_PQMM, 0)
	} else if self.room.Type == GAMETYPE_KWX_XY { //! 襄阳
		self.SetParam(TYPE_SK, 0)
		self.SetParam2(TYPE_SL, 0)
		self.SetParam2(TYPE_DLDF, 0)
		self.SetParam2(TYPE_PQMM, 0)
	} else if self.room.Type == GAMETYPE_KWX_SY { //! 十堰
		self.SetParam(TYPE_SK, 0)
		self.SetParam2(TYPE_DLDF, 0)
		self.SetParam2(TYPE_PQMM, 0)
	} else if self.room.Type == GAMETYPE_KWX_SZ { //! 随州
		self.SetParam(TYPE_SK, 0)
		if self.GetParam(TYPE_MM) == 0 {
			self.SetParam(TYPE_MM, 1)
		}
		self.SetParam(TYPE_PH, 1)
		self.SetParam(TYPE_PD, 0)
		self.SetParam2(TYPE_SL, 0)
		self.SetParam2(TYPE_DLDF, 0)
		self.SetParam2(TYPE_PQMM, 0)
	} else if self.room.Type == GAMETYPE_KWX_YIC { //! 宜城
		self.SetParam(TYPE_SK, 0)
		self.SetParam(TYPE_PH, 1)
		self.SetParam(TYPE_PD, 0)
		self.SetParam2(TYPE_SL, 0)
		self.SetParam2(TYPE_DLDF, 0)
	} else { //! 应城
		//self.SetParam(TYPE_SK, 0)
		self.SetParam(TYPE_PH, 1)
		self.SetParam(TYPE_PD, 0)
		self.SetParam2(TYPE_SL, 0)
		//self.SetParam2(TYPE_DLDF, 0)
		self.SetParam2(TYPE_PQMM, 0)
	}
}

func (self *Game_KWX) OnRobot(robot *lib.Robot) {

}

func (self *Game_KWX) OnSendInfo(person *Person) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			person.SendMsg("gamekwxinfo", self.getInfo(person.Uid))
			return
		}
	}

	_person := new(Game_KWX_Person)
	_person.Init()
	_person.Uid = person.Uid
	_person.Ready = true
	_person.Piao = -1
	self.PersonMgr = append(self.PersonMgr, _person)

	if self.GetParam(TYPE_PIAO) == 0 { //! 不加飘
		if len(self.PersonMgr) >= lib.HF_Atoi(self.room.csv["minnum"]) {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
			self.OnBegin()
			self.room.flush()
			return
		}
	}

	person.SendMsg("gamekwxinfo", self.getInfo(person.Uid))
}

func (self *Game_KWX) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "gameready": //! 游戏准备
		self.GameReady(msg.Uid)
		self.room.flush()
	case "gamebets": //! 加飘
		self.GameBets(msg.Uid, msg.V.(*Msg_GameBets).Bets)
		self.room.flush()
	case "gamestep": //! 出牌
		self.GameStep(msg.Uid, msg.V.(*Msg_GameStep).Card)
		self.room.flush()
	case "gamecagang":
		self.GameCaGang(msg.Uid, msg.V.(*Msg_GameStep).Card)
		self.room.flush()
	case "gamepeng": //! 碰
		self.GamePeng(msg.Uid)
		self.room.flush()
	case "gamegang": //! 杠
		self.GameGang(msg.Uid, msg.V.(*Msg_GameStep).Card)
		self.room.flush()
	case "gamehu": //! 胡
		self.GameHu(msg.Uid)
		self.room.flush()
	case "gameguo": //! 过
		self.GameGuo(msg.Uid)
		self.room.flush()
	case "gamekwxview":
		self.GameView(msg.Uid, msg.V.(*Msg_GameKWX_View).CardL, msg.V.(*Msg_GameKWX_View).Want, msg.V.(*Msg_GameKWX_View).Card)
		self.room.flush()
	case "gamekwxneed":
		if GetServer().Con.IsNeed != 0 {
			person := self.GetPerson(msg.Uid)
			if person == nil {
				return
			}
			person.need = msg.V.(*Msg_GameKWX_Need).Card
		}
	case "gamekwxmygod":
		if GetServer().Con.IsNeed == 0 { //! 正式服只有配置的号才能用
			if !GetServer().IsAdmin(msg.Uid, staticfunc.ADMIN_KWX) {
				return
			}
		}
		person := GetPersonMgr().GetPerson(msg.Uid)
		if person == nil {
			return
		}
		var _msg GameKWX_GoldInfo
		if self.Mah != nil {
			_msg.Info = self.Mah.Card
		} else {
			_msg.Info = make([]int, 0)
		}
		person.SendMsg("gamekwxmygod", &_msg)
	case "kwxgetmygod":
		if GetServer().Con.IsNeed == 0 { //! 正式服只有配置的号才能用
			if !GetServer().IsAdmin(msg.Uid, staticfunc.ADMIN_KWX) {
				return
			}
		}
		person := self.GetPerson(msg.Uid)
		if person == nil {
			return
		}
		person.NextCard = msg.V.(*GameKWX_GoldGet).Card
	}
}

func (self *Game_KWX) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)
	if self.Winer == 0 {
		self.Winer = self.room.Uid[0]
	}

	self.Mah = NewMah_KWX()
	self.BefStep = 0
	self.ViewUid = 0
	self.Gang = 0
	self.GangUid = 0
	self.Hz = false
	self.Ma = make([]int, 0)
	self.Record = new(Rec_GameKWX_Info)
	self.IsRecord = false
	if self.GetParam2(TYPE_SL) == 1 { //! 上楼
		self.SL[0] = lib.HF_GetRandom(6) + 1
		self.SL[1] = lib.HF_GetRandom(6) + 1
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Init()
		if self.PersonMgr[i].Uid == self.Winer {
			self.PersonMgr[i].Deal = true
			self.PersonMgr[i].Card.Card1 = self.Mah.DealNeed(self.PersonMgr[i].need)
		} else {
			self.PersonMgr[i].Card.Card1 = self.Mah.DealNeed(self.PersonMgr[i].need)
		}
		//if i == 0 {
		//	self.PersonMgr[i].Card.Card1 = []int{18, 16, 3, 6, 4, 9, 32, 14, 17, 31, 3, 17, 2}
		//} else if i == 1 {
		//	self.PersonMgr[i].Card.Card1 = []int{1, 1, 15, 32, 33, 5, 33, 12, 7, 11, 14, 32, 8}
		//} else {
		//	self.PersonMgr[i].Card.Card1 = []int{18, 5, 7, 3, 32, 13, 9, 15, 31, 7, 14, 15, 5}
		//}
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "玩家", i, ":", self.PersonMgr[i].Card.Card1)

		//! 记录
		var rc_person Son_Rec_GameKWX_Person
		rc_person.Uid = self.PersonMgr[i].Uid
		rc_person.Name = self.room.GetName(rc_person.Uid)
		rc_person.Head = self.room.GetHead(rc_person.Uid)
		lib.HF_DeepCopy(&rc_person.Card, &self.PersonMgr[i].Card.Card1)
		self.Record.Person = append(self.Record.Person, rc_person)
	}
	//! 庄家进入摸牌阶段
	self.GameMoCard(self.Winer, false)
	self.Winer = 0

	for i := 0; i < len(self.PersonMgr); i++ {
		person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
		if person == nil {
			continue
		}
		person.SendMsg("gamekwxbegin", self.getInfo(person.Uid))
	}
}

func (self *Game_KWX) OnEnd() {
	//! 结算时先终止擦杠
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Gang > 100 {
			self.StopCaGang(self.PersonMgr[i], self.PersonMgr[i].Gang-100)
			break
		}
	}

	self.room.SetBegin(false)

	//! 判断是否荒庄
	self.Hz = true
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].So_Card != 0 {
			self.Hz = false
			break
		}
	}
	//! 荒庄赔胡规则
	if self.Hz && self.GetParam(TYPE_PH) == 1 {
		num := 0
		hzstate := 0
		var no1person *Game_KWX_Person = nil //! 不是第一个亮牌的人
		for _, value := range self.PersonMgr {
			self.getTiHuFans(value)
			if len(value.Card.CardL) > 0 {
				if self.ViewUid == value.Uid {
					value.hz = 1
				} else {
					if self.room.Type == GAMETYPE_KWX_SZ { //! 随州玩法第一家亮牌的赔胡
						value.hz = 2
						no1person = value
					} else {
						value.hz = 1
					}
				}
			} else if value.So_Card > 0 {
				value.hz = 10
				num++
				hzstate = 1
			} else {
				value.hz = 0
				hzstate = 2
			}
		}
		if self.room.Type == GAMETYPE_KWX_SZ && no1person != nil { //! 纠正随州的赔胡
			if hzstate == 2 { //! 有人未听胡
				no1person.hz = 1
			}
		}
		if num == 3 { //! 3家都听胡不亮牌
			self.Winer = self.CurStep
		}
	}

	if self.Winer == 0 {
		winnum := 0
		gaofen := 0
		difen := 0
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].So_Card > 0 {
				winnum++
			}
			if self.PersonMgr[i].So_Card > self.PersonMgr[gaofen].So_Card {
				gaofen = i
			}
			if self.PersonMgr[i].So_Card < self.PersonMgr[difen].So_Card {
				difen = i
			}
		}
		if winnum == 2 { //! 一炮双响
			self.Winer = self.PersonMgr[difen].Uid
		} else {
			self.Winer = self.PersonMgr[gaofen].Uid
		}
	}

	score := make([]int, 0)
	for i := 0; i < len(self.PersonMgr); i++ {
		score = append(score, self.PersonMgr[i].So_Card)
		self.PersonMgr[i].So_Card = 0

		self.PersonMgr[i].AG += len(self.PersonMgr[i].Card.CardAG)
		self.PersonMgr[i].MG += (len(self.PersonMgr[i].Card.CardMG) + len(self.PersonMgr[i].Card.CardCG))
	}

	//! 统计牌分
	for key, value := range self.PersonMgr {
		if len(value.Card.CardL) > 0 {
			value.AddState(7, 2)
		}
		if score[key] != 0 {
			self.AddRecordStep(value.Uid, 4, []int{})
			if self.GetParam(TYPE_SK) == 1 {
				if value.Kan > 0 {
					value.AddState(32, value.Kan)
				}
			}
			if value.Piao >= 0 {
				value.AddState(31, value.Piao)
			}
			//! 计算跑恰摸八
			pqmb := 0
			if self.GetParam2(TYPE_PQMM) == 1 && !self.Hz {
				maxnum := 0
				pqmb += 1                      //! 跑
				pqmb += value.Kan              //! 恰
				if self.CurStep == value.Uid { //! 摸
					pqmb += 1
					maxnum = MahMaxNum(&value.Card, value.Card.CardM)
				} else {
					maxnum = MahMaxNum(&value.Card, self.LastCard)
				}
				if maxnum > 7 { //! 八
					pqmb += (maxnum - 7)
				}
				value.AddState(41, pqmb)
			}
			if self.CurStep == value.Uid || self.Hz { //! 自摸或荒庄
				if !self.Hz {
					value.ZM++
					value.AddState(17, 1)
					self.LastCard = value.Card.CardM
				}
				mascore := 0
				if !self.Hz && (self.GetParam(TYPE_DL) == 1 || len(value.Card.CardL) > 0) && len(self.Mah.Card) > 0 { //! 没有荒庄并且亮牌
					if self.GetParam(TYPE_MM) == 1 { //! 独马
						self.Ma = append(self.Ma, self.Mah.Draw())
						if self.Ma[0]/10 == 3 {
							mascore = 10
						} else {
							mascore = self.Ma[0] % 10
						}
					} else if self.GetParam(TYPE_MM) == 2 { //! 六马
						six := true
						for i := 0; i < 6; i++ {
							if len(self.Mah.Card) == 0 {
								six = false
								break
							}
							self.Ma = append(self.Ma, self.Mah.Draw())
							if self.Ma[i]/10 == 3 { //! 中发白+1
								mascore += 1
							} else {
								m1 := self.Ma[i] % 10
								m2 := value.Card.CardM % 10
								if (m1 == 1 || m1 == 4 || m1 == 7) && (m2 == 1 || m2 == 4 || m2 == 7) {
									mascore += 1
								} else if (m1 == 2 || m1 == 5 || m1 == 8) && (m2 == 2 || m2 == 5 || m2 == 8) {
									mascore += 1
								} else if (m1 == 3 || m1 == 6 || m1 == 9) && (m2 == 3 || m2 == 6 || m2 == 9) {
									mascore += 1
								}
							}
						}
						if mascore == 0 && six {
							mascore = 6
						}
					}
				}
				if mascore > 0 {
					value.AddState(20, mascore)
				}
				for i := 0; i < len(self.PersonMgr); i++ {
					if i == key {
						continue
					}
					if self.Hz && value.hz <= self.PersonMgr[i].hz {
						continue
					}
					fan := score[key]
					if len(value.Card.CardL) > 0 || len(self.PersonMgr[i].Card.CardL) > 0 || (self.Hz && (self.PersonMgr[i].hz == 1 || self.PersonMgr[i].hz == 2)) {
						fan *= 2
						if self.GetParam2(TYPE_DLDF) == 1 && ((len(value.Card.CardL) > 0 && len(self.PersonMgr[i].Card.CardL) > 0) || (self.Hz && self.PersonMgr[i].hz == 1 && self.PersonMgr[i].hz == 2)) {
							self.PersonMgr[i].AddState(19, 2)
							fan *= 2
						}
					}
					if self.GetParam2(TYPE_SL) == 1 && self.SL[0] == self.SL[1] {
						value.AddState(40, 2)
						fan *= 2
					}
					if self.GetParam(TYPE_ZDFS) == 0 {
						fan = lib.HF_MinInt(fan, 8)
					} else {
						fan = lib.HF_MinInt(fan, 16)
					}
					value.So_Card += fan
					self.PersonMgr[i].So_Card -= fan

					piao := lib.HF_MaxInt(0, value.Piao) + lib.HF_MaxInt(0, self.PersonMgr[i].Piao)
					value.So_Other += piao
					self.PersonMgr[i].So_Other -= piao

					if self.GetParam(TYPE_SK) == 1 {
						value.So_Other += value.Kan
						self.PersonMgr[i].So_Other -= value.Kan
					}

					value.So_Other += mascore
					self.PersonMgr[i].So_Other -= mascore

					value.So_Other += pqmb
					self.PersonMgr[i].So_Other -= pqmb
				}
			} else {
				value.JP++
				self.LastCard = self.LastCard
				for i := 0; i < len(self.PersonMgr); i++ {
					if self.PersonMgr[i].Uid != self.CurStep {
						continue
					}
					self.PersonMgr[i].DP++
					self.PersonMgr[i].AddState(18, 1)
					fan := score[key]
					bh := false
					if len(value.Card.CardL) > 0 { //! 亮牌了
						fan *= 2
						tmp := fan
						if self.GetParam2(TYPE_DLDF) == 1 && len(self.PersonMgr[i].Card.CardL) > 0 {
							self.PersonMgr[i].AddState(19, 2)
							fan *= 2
						}
						if (self.room.Type == GAMETYPE_KWX_XY || self.room.Type == GAMETYPE_KWX_SY) && len(self.PersonMgr[i].Card.CardL) == 0 { //! 襄阳和十堰进入包胡
							for j := 0; j < len(self.PersonMgr); j++ {
								if self.PersonMgr[j].Uid == self.CurStep || self.PersonMgr[j].Uid == value.Uid {
									continue
								}
								if self.GetParam2(TYPE_DLDF) == 1 && len(self.PersonMgr[j].Card.CardL) > 0 {
									self.PersonMgr[j].AddState(19, 2)
									tmp *= 2
								}
								fan += tmp
								break
							}
							bh = true
						}
					} else { //! 没有亮牌
						if len(self.PersonMgr[i].Card.CardL) > 0 {
							fan *= 2
						}
					}
					if self.GetParam2(TYPE_SL) == 1 && self.SL[0] == self.SL[1] {
						value.AddState(40, 2)
						fan *= 2
					}
					if self.GetParam(TYPE_ZDFS) == 0 {
						fan = lib.HF_MinInt(fan, 8)
					} else {
						fan = lib.HF_MinInt(fan, 16)
					}
					value.So_Card += fan
					self.PersonMgr[i].So_Card -= fan

					piao := lib.HF_MaxInt(0, value.Piao) + lib.HF_MaxInt(0, self.PersonMgr[i].Piao)
					if bh { //! 包胡还要加另外一家的飘
						for j := 0; j < len(self.PersonMgr); j++ {
							if self.PersonMgr[j].Uid == self.CurStep || self.PersonMgr[j].Uid == value.Uid {
								continue
							}
							piao += lib.HF_MaxInt(0, value.Piao) + lib.HF_MaxInt(0, self.PersonMgr[j].Piao)
							break
						}
					}
					value.So_Other += piao
					self.PersonMgr[i].So_Other -= piao

					if self.GetParam(TYPE_SK) == 1 {
						kan := value.Kan
						if bh {
							kan *= 2
						}
						value.So_Other += kan
						self.PersonMgr[i].So_Other -= kan
					}

					value.So_Other += pqmb
					self.PersonMgr[i].So_Other -= pqmb
				}
			}
		}
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Total += self.PersonMgr[i].So_Card
		//self.PersonMgr[i].Total += self.PersonMgr[i].So_Gang
		self.PersonMgr[i].Total += self.PersonMgr[i].So_Other

		for j := 0; j < len(self.Record.Person); j++ {
			if self.Record.Person[j].Uid == self.PersonMgr[i].Uid {
				self.Record.Person[j].Score = (self.PersonMgr[i].So_Card + self.PersonMgr[i].So_Gang + self.PersonMgr[i].So_Other)
				self.Record.Person[j].Total = self.PersonMgr[i].Total
				break
			}
		}
	}

	//! 发消息
	var msg Msg_GameKWX_End
	msg.Hz = self.Hz
	msg.Ma = self.Ma
	msg.LastCard = self.LastCard
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameKWX_Info
		son.Uid = self.PersonMgr[i].Uid
		son.Deal = self.PersonMgr[i].Deal
		son.Piao = self.PersonMgr[i].Piao
		son.Card = self.PersonMgr[i].Card
		son.Card.Want = make([]int, 0)
		son.So_Card = self.PersonMgr[i].So_Card
		son.So_Gang = self.PersonMgr[i].So_Gang
		son.So_Other = self.PersonMgr[i].So_Other
		son.Total = self.PersonMgr[i].Total
		son.State = self.PersonMgr[i].State
		self.PersonMgr[i].Hu = score[i]
		son.Hu = score[i]
		msg.Info = append(msg.Info, son)
	}
	self.room.broadCastMsg("gamekwxend", &msg)

	self.Record.Roomid = self.room.Id*100 + self.room.Step
	self.Record.Time = time.Now().Unix()
	self.Record.Param1 = self.room.Param1
	self.Record.Param2 = self.room.Param2
	self.room.AddRecord(lib.HF_JtoA(self.Record))
	self.IsRecord = true

	if self.room.IsBye() {
		self.OnBye()
		self.room.Bye()
		return
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Ready = false
		self.PersonMgr[i].hz = 0
	}
}

func (self *Game_KWX) OnBye() {
	if !self.IsRecord && self.Record != nil {
		self.Record.Roomid = self.room.Id*100 + self.room.Step
		self.Record.Time = time.Now().Unix()
		self.Record.Param1 = self.room.Param1
		self.Record.Param2 = self.room.Param2
		self.room.AddRecord(lib.HF_JtoA(self.Record))
		self.IsRecord = true
	}

	info := make([]staticfunc.JS_CreateRoomMem, 0)
	var msg Msg_GameKWX_Bye
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameKWX_Bye
		son.Uid = self.PersonMgr[i].Uid
		son.ZM = self.PersonMgr[i].ZM
		son.JP = self.PersonMgr[i].JP
		son.DP = self.PersonMgr[i].DP
		son.AG = self.PersonMgr[i].AG
		son.MG = self.PersonMgr[i].MG
		son.PI = lib.HF_MaxInt(0, self.PersonMgr[i].Piao)
		son.Score = self.PersonMgr[i].Total
		msg.Info = append(msg.Info, son)
		info = append(info, staticfunc.JS_CreateRoomMem{son.Uid, "", "", son.Score})
	}
	self.room.broadCastMsg("gamekwxbye", &msg)

	self.room.ClubResult(info)
}

func (self *Game_KWX) OnExit(uid int64) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			copy(self.PersonMgr[i:], self.PersonMgr[i+1:])
			self.PersonMgr = self.PersonMgr[:len(self.PersonMgr)-1]
			break
		}
	}
}

//! 准备,第一局自动准备
func (self *Game_KWX) GameReady(uid int64) {
	if self.room.IsBye() {
		return
	}

	if self.room.Begin {
		return
	}

	if self.room.Step == 0 { //! 一局之后才有这个消息
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

	if num == len(self.room.Uid) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
		self.OnBegin()
		return
	}

	var msg staticfunc.Msg_Uid
	msg.Uid = uid
	self.room.broadCastMsg("gameready", &msg)
}

//! 加飘
func (self *Game_KWX) GameBets(uid int64, bets int) {
	if self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经开始了:GameBets")
		return
	}

	if bets < 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "加漂小于0:GameBets")
		return
	}

	if self.GetParam(TYPE_PIAO) == 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能加漂:GameBets")
		return
	}

	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Piao != -1 {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经加漂过:GameBets")
				return
			} else {
				self.PersonMgr[i].Piao = bets
				num++
			}
		} else if self.PersonMgr[i].Piao != -1 {
			num++
		}
	}

	if num == len(self.PersonMgr) && len(self.PersonMgr) >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 所有人都加飘
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
		self.OnBegin()
		return
	}

	var msg Msg_GameBets
	msg.Uid = uid
	msg.Bets = bets
	self.room.broadCastMsg("gamebets", &msg)
}

//! 摸牌阶段(系统自动跳入)
func (self *Game_KWX) GameMoCard(uid int64, send bool) {
	if len(self.Mah.Card) == 0 { //! 荒庄
		self.OnEnd()
		return
	}

	self.CurStep = uid

	person := self.GetPerson(uid)
	if person.Card.CardM != 0 {
		person.Card.Card1 = append(person.Card.Card1, person.Card.CardM)
	}
	person.Card.CardM = self.Mah.Draw4(person.NextCard)
	person.NoHu = 0
	person.NextCard = 0
	self.AddRecordStep(person.Uid, 0, []int{person.Card.CardM})

	var msg Msg_GameKWX_Draw
	msg.Uid = uid

	if len(person.Card.Want) > 0 {
		for i := 0; i < len(person.Card.Want); i++ {
			if person.Card.Want[i] == person.Card.CardM {
				self.getHuFans(person, person.Card.CardM, true, false)
				break
			}
		}
	} else {
		self.getHuFans(person, person.Card.CardM, true, false)
	}

	if person.So_Card > 0 { //! 可以胡
		if len(person.Card.CardL) > 0 { //! 已经亮牌
			if send { //! 发消息
				for i := 0; i < len(self.PersonMgr); i++ {
					_person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
					if _person == nil {
						continue
					}
					if self.PersonMgr[i].Uid == uid {
						msg.Card = person.Card.CardM
					} else {
						msg.Card = 0
					}
					_person.SendMsg("gamekwx_draw", &msg)
				}
			}
			self.OnEnd()
			return
		}
		if self.room.Type == GAMETYPE_KWX_YC { //! 应城自摸必须2番胡
			if person.So_Card > 1 {
				person.Hu = 1
			}
		} else {
			person.Hu = 1
		}
	}

	if len(self.Mah.Card) > 0 && MahIsGang(&person.Card, person.Card.CardM, true) { //! 可以杠
		person.Gang = 1
	}
	if send { //! 发消息
		for i := 0; i < len(self.PersonMgr); i++ {
			_person := GetPersonMgr().GetPerson(self.PersonMgr[i].Uid)
			if _person == nil {
				continue
			}
			if self.PersonMgr[i].Uid == uid {
				msg.Card = person.Card.CardM
				msg.Hu = person.Hu
				msg.Gang = person.Gang
			} else {
				msg.Card = 0
				msg.Hu = 0
				msg.Gang = 0
			}
			_person.SendMsg("gamekwx_draw", &msg)
		}
	}
	//! 不能胡并且已经亮牌了，判断杀马
	if person.Hu == 0 && len(person.Card.CardL) > 0 && self.room.Type == GAMETYPE_KWX_SZ {
		num := 0
		for i := 0; i < len(person.Card.CardL); i++ {
			if person.Card.CardL[i] == person.Card.CardM {
				num++
			}
		}
		if num == 3 { //! 摸的牌亮出刻子
			person.Card.CardAG = append(person.Card.CardAG, person.Card.CardM)
			self.AddRecordStep(person.Uid, 5, []int{person.Card.CardM})
			person.So_Gang += 4 * int(math.Pow(float64(2), float64(self.Gang)))
			person.Total += 4 * int(math.Pow(float64(2), float64(self.Gang)))
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Uid == person.Uid {
					continue
				}
				self.PersonMgr[i].So_Gang -= 2 * int(math.Pow(float64(2), float64(self.Gang)))
				self.PersonMgr[i].Total -= 2 * int(math.Pow(float64(2), float64(self.Gang)))
			}
			self.Gang++
			for i := 0; i < len(person.Card.CardL); {
				if person.Card.CardL[i] == person.Card.CardM {
					copy(person.Card.CardL[i:], person.Card.CardL[i+1:])
					person.Card.CardL = person.Card.CardL[:len(person.Card.CardL)-1]
				} else {
					i++
				}
			}

			for i := 0; i < len(person.Card.Want); {
				self.getHuFans(person, person.Card.Want[i], true, true)
				if person.So_Card == 0 {
					copy(person.Card.Want[i:], person.Card.Want[i+1:])
					person.Card.Want = person.Card.Want[:len(person.Card.Want)-1]
				} else {
					i++
				}
				person.So_Card = 0
				person.State = make([]Son_KWX_State, 0)
			}

			var msg Msg_GameKWX_Kill
			msg.Uid = uid
			msg.Card = person.Card.CardM
			msg.View = false
			msg.Want = person.Card.Want
			self.room.broadCastMsg("gamekwx_kill", &msg)

			self.SendTotal()

			person.Card.CardM = 0

			self.GameMoCard(uid, true)
		}
	}
}

//! 出牌(玩家选择)
func (self *Game_KWX) GameStep(uid int64, card int) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始")
		return
	}

	if card == 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "GameStep(card=0)")
		return
	}

	if self.CurStep != uid {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "当前不是你的局:GameStep")
		return
	}

	if self.BefStep == uid {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经出过牌:GameStep")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person:GameStep")
		return
	}

	if person.Gang == 1 || person.Peng == 1 || person.Hu == 1 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "有操作不让出牌:GameStep")
		return
	}

	//if len(person.Card.CardL) == 0 {
	//	for i := 0; i < len(self.PersonMgr); i++ {
	//		if self.PersonMgr[i].Uid == uid {
	//			continue
	//		}
	//		for j := 0; j < len(self.PersonMgr[i].Card.Want); j++ {
	//			if self.PersonMgr[i].Card.Want[j] == card {
	//				lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能打亮牌需要的牌")
	//				return
	//			}
	//		}
	//	}
	//}

	if person.Card.CardM != card {
		if len(person.Card.CardL) > 0 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "亮牌后只能打摸的牌:GameStep")
			return
		}

		find := false
		for i := 0; i < len(person.Card.Card1); i++ {
			if person.Card.Card1[i] == card {
				copy(person.Card.Card1[i:], person.Card.Card1[i+1:])
				person.Card.Card1 = person.Card.Card1[:len(person.Card.Card1)-1]
				find = true
				break
			}
		}
		if !find {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "出牌错误:GameStep")
			return
		}
		if person.Card.CardM != 0 {
			person.Card.Card1 = append(person.Card.Card1, person.Card.CardM)
		}
	}

	//! 出牌广播
	var msg Msg_GameKWX_Step
	msg.Uid = uid
	msg.Card = card
	self.room.broadCastMsg("gamekwxstep", &msg)

	cur := time.Now().UnixNano()

	person.Card.CardM = 0
	person.Card.Card2 = append(person.Card.Card2, card)
	person.So_Card = 0
	person.State = make([]Son_KWX_State, 0)
	self.LastCard = card
	self.AddRecordStep(person.Uid, 1, []int{card})

	self.BefStep = uid

	//! 轮流判断一下其他人是否决定
	finduid := make([]int, 0)
	isend := false
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			continue
		}
		find := false
		if len(self.PersonMgr[i].Card.Want) > 0 {
			for j := 0; j < len(self.PersonMgr[i].Card.Want); j++ {
				if card == self.PersonMgr[i].Card.Want[j] {
					self.getHuFans(self.PersonMgr[i], card, false, false)
					break
				}
			}
		} else {
			self.getHuFans(self.PersonMgr[i], card, false, false)
		}

		if self.PersonMgr[i].So_Card > 1 || (self.PersonMgr[i].So_Card > 0 && (len(self.PersonMgr[i].Card.CardL) > 0 || len(person.Card.CardL) > 0)) { //! 可以胡
			self.PersonMgr[i].Hu = 1
			self.PersonMgr[i].NoHu = card
			find = true
			if len(self.PersonMgr[i].Card.CardL) > 0 {
				isend = true
			}
		} else {
			self.PersonMgr[i].So_Card = 0
			self.PersonMgr[i].State = make([]Son_KWX_State, 0)
		}

		if len(self.Mah.Card) > 0 && MahIsGang(&self.PersonMgr[i].Card, card, false) { //! 可以杠
			self.PersonMgr[i].Gang = 1
			find = true
		}

		if MahIsPeng(&self.PersonMgr[i].Card, card) { //! 可以碰
			self.PersonMgr[i].Peng = 1
			find = true
		}

		if find {
			finduid = append(finduid, i)
		}
	}

	delay := time.Now().UnixNano() - cur

	lib.GetLogMgr().Output(lib.LOG_DEBUG, fmt.Sprintf("消息处理毫秒(%d)纳秒(%d)", (delay/1000000), delay))

	if isend {
		self.OnEnd()
		return
	}

	if len(finduid) == 0 {
		if self.room.Type == GAMETYPE_KWX_SZ { //! 判断杀马
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Uid == uid {
					continue
				}

				if len(self.PersonMgr[i].Card.CardL) > 0 {
					num := 0
					for j := 0; j < len(self.PersonMgr[i].Card.CardL); j++ {
						if self.PersonMgr[i].Card.CardL[j] == card {
							num++
						}
					}
					if num != 3 {
						continue
					}

					person.So_Gang -= 2 * int(math.Pow(float64(2), float64(self.Gang)))
					person.Total -= 2 * int(math.Pow(float64(2), float64(self.Gang)))
					person.Card.Card2 = person.Card.Card2[:len(person.Card.Card2)-1]
					self.PersonMgr[i].Card.CardMG = append(self.PersonMgr[i].Card.CardMG, card)
					self.AddRecordStep(self.PersonMgr[i].Uid, 6, []int{card})
					self.PersonMgr[i].So_Gang += 2 * int(math.Pow(float64(2), float64(self.Gang)))
					self.PersonMgr[i].Total += 2 * int(math.Pow(float64(2), float64(self.Gang)))
					self.LastCard = 0
					self.Gang++
					self.GangUid = self.PersonMgr[i].Uid
					for j := 0; j < len(self.PersonMgr[i].Card.CardL); {
						if self.PersonMgr[i].Card.CardL[j] == card {
							copy(self.PersonMgr[i].Card.CardL[j:], self.PersonMgr[i].Card.CardL[j+1:])
							self.PersonMgr[i].Card.CardL = self.PersonMgr[i].Card.CardL[:len(self.PersonMgr[i].Card.CardL)-1]
						} else {
							j++
						}
					}

					for j := 0; j < len(self.PersonMgr[i].Card.Want); {
						self.getHuFans(self.PersonMgr[i], self.PersonMgr[i].Card.Want[j], true, true)
						if self.PersonMgr[i].So_Card == 0 {
							copy(self.PersonMgr[i].Card.Want[j:], self.PersonMgr[i].Card.Want[j+1:])
							self.PersonMgr[i].Card.Want = self.PersonMgr[i].Card.Want[:len(self.PersonMgr[i].Card.Want)-1]
						} else {
							j++
						}
						self.PersonMgr[i].So_Card = 0
						self.PersonMgr[i].State = make([]Son_KWX_State, 0)
					}

					var msg Msg_GameKWX_Kill
					msg.Uid = self.PersonMgr[i].Uid
					msg.Card = card
					msg.View = true
					msg.Want = self.PersonMgr[i].Card.Want
					self.room.broadCastMsg("gamekwx_kill", &msg)

					self.SendTotal()

					self.GameMoCard(self.PersonMgr[i].Uid, true)

					return
				}
			}
		}

		self.Gang = 0
		self.GangUid = 0
		self.GameMoCard(self.GetNextUid(), true)
	} else {
		for i := 0; i < len(finduid); i++ {
			self.SendOperator(self.PersonMgr[finduid[i]].Uid, self.PersonMgr[finduid[i]].Hu, self.PersonMgr[finduid[i]].Gang, self.PersonMgr[finduid[i]].Peng)
		}
	}
}

//! 擦杠(玩家选择)
func (self *Game_KWX) GameCaGang(uid int64, card int) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始:GameCaGang")
		return
	}

	if self.CurStep != uid {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "当前不是你的局:GameCaGang")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person:GameCaGang")
		return
	}

	if person.Card.CardM != card {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "擦杠必须是自摸的牌:GameCaGang")
		return
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			continue
		}
		for j := 0; j < len(self.PersonMgr[i].Card.Want); j++ {
			if self.PersonMgr[i].Card.Want[j] == card {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能打亮牌需要的牌:GameCaGang")
				return
			}
		}
	}

	index := -1
	for i := 0; i < len(person.Card.CardP); i++ {
		if person.Card.CardP[i] == card {
			index = i
			break
		}
	}

	if index == -1 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不可擦杠:GameCaGang")
		return
	}
	person.Card.CardM = 0
	person.Card.Card2 = append(person.Card.Card2, card)
	person.So_Card = 0
	person.Peng = 0
	person.Hu = 0
	person.State = make([]Son_KWX_State, 0)
	self.LastCard = card

	self.BefStep = uid

	//! 判断是否能抢杠
	findindex := make([]int, 0)
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			continue
		}

		self.getHuFans(self.PersonMgr[i], card, false, false)
		if self.PersonMgr[i].So_Card > 0 {
			self.PersonMgr[i].AddState(37, 2)
			self.PersonMgr[i].So_Card *= 2
			self.PersonMgr[i].Hu = 1
			findindex = append(findindex, i)
		}
	}

	if len(findindex) == 0 {
		self.CaGangCard(person, card)
		self.GameMoCard(uid, true)
	} else {
		for i := 0; i < len(findindex); i++ {
			self.SendOperator(self.PersonMgr[findindex[i]].Uid, 1, 0, 0)
		}
		self.WillCaGang(person, card)
		person.Gang = 100 + card
	}
}

//! 碰
func (self *Game_KWX) GamePeng(uid int64) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始:GamePeng")
		return
	}

	if self.CurStep == uid {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "自己局不能碰:GamePeng")
		return
	}

	if self.LastCard == 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "最后牌错误:GamePeng")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person:GamePeng")
		return
	}

	if person.Peng != 1 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能碰:GamePeng")
		return
	}

	person.Peng = 2
	person.Gang = 0
	person.Hu = 0
	person.So_Card = 0
	person.State = make([]Son_KWX_State, 0)

	find := false
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			continue
		}
		if self.PersonMgr[i].Hu == 1 { //! 胡的人没表态
			find = true
			break
		}
	}

	if !find { //! 没有可以胡的人
		self.PengCard(person)
		self.CurStep = uid
		if len(self.Mah.Card) > 0 && MahIsGang(&person.Card, 0, true) { //! 可以杠
			person.Gang = 1
			self.SendOperator(person.Uid, 0, 1, 0)
		}
	}
}

//! 杠
func (self *Game_KWX) GameGang(uid int64, card int) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始:GameGang")
		return
	}

	if self.CurStep != uid {
		if self.LastCard != card {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "杠的不是最后一张牌:GameGang")
			return
		}
	}

	person := self.GetPerson(uid)
	if person == nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person:GameGang")
		return
	}

	if person.Gang != 1 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能杠:GameGang")
		return
	}

	if self.CurStep == uid {
		num := 0
		for i := 0; i < len(person.Card.Card1); i++ {
			if person.Card.Card1[i] == card {
				num++
			}
		}
		for i := 0; i < len(person.Card.CardL); i++ {
			if person.Card.CardL[i] == card {
				num++
			}
		}
		if person.Card.CardM == card {
			num++
		}
		if num < 4 {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "card错误:GameGang")
			return
		}
	}

	person.Gang = 2
	person.Peng = 0
	person.Hu = 0
	person.So_Card = 0
	person.State = make([]Son_KWX_State, 0)

	find := false
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			continue
		}
		if self.PersonMgr[i].Hu == 1 { //! 胡的人没表态
			find = true
			break
		}
	}

	if !find { //! 没有可以胡的人
		self.GangCard(person, card)
		self.GameMoCard(uid, true)
	}
}

//! 胡
func (self *Game_KWX) GameHu(uid int64) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始:GameHu")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person:GameHu")
		return
	}

	if person.Hu != 1 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能胡:GameHu")
		return
	}

	person.Hu = 2

	find := false
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			continue
		}
		if self.PersonMgr[i].Hu == 1 { //! 胡的人没表态
			find = true
			break
		}
	}

	if !find { //! 没有可以胡的人
		self.OnEnd()
		return
	}
}

//! 过
func (self *Game_KWX) GameGuo(uid int64) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未开始:GameGuo")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person:GameGuo")
		return
	}

	if person.Hu == 0 && person.Peng == 0 && person.Gang == 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "不能过:GameGuo")
		return
	}

	person.Hu = 0
	person.Peng = 0
	person.Gang = 0
	person.So_Card = 0
	person.State = make([]Son_KWX_State, 0)

	//! 若还有人有操作
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			continue
		}
		if self.PersonMgr[i].Hu == 1 || self.PersonMgr[i].Gang == 1 || self.PersonMgr[i].Peng == 1 {
			//person.SendNullMsg()
			return
		}
	}

	//! 判断有人胡
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			continue
		}
		if self.PersonMgr[i].Hu == 2 {
			self.OnEnd()
			return
		}
	}

	//! 判断有人杠
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			continue
		}
		if self.PersonMgr[i].Gang > 100 { //! 擦杠
			self.FinishCaGang(person, self.PersonMgr[i].Gang-100)
			self.GameMoCard(self.PersonMgr[i].Uid, true)
			return
		}
		if self.PersonMgr[i].Gang == 2 {
			self.GangCard(self.PersonMgr[i], self.LastCard)
			self.GameMoCard(self.PersonMgr[i].Uid, true)
			return
		}
	}

	//! 判断有人碰
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			continue
		}
		if self.PersonMgr[i].Peng == 2 {
			self.PengCard(self.PersonMgr[i])
			self.CurStep = self.PersonMgr[i].Uid
			return
		}
	}

	if self.CurStep != uid {
		self.Gang = 0
		self.GangUid = 0
		self.GameMoCard(self.GetNextUid(), true)
	}
}

func (self *Game_KWX) GameView(uid int64, cardl []int, want []int, card int) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "GameView游戏未开始")
		return
	}

	if card == 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "GameView(card=0)")
		return
	}

	if self.CurStep != uid {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "GameView当前不是你的局")
		return
	}

	if self.BefStep == uid {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "GameView已经出过牌")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {

		lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到person")
		return
	}

	//if len(self.Mah.Card) <= 12 && self.GetParam(TYPE_12VIEW) == 1 {
	//	lib.GetLogMgr().Output("GameView小于12张牌不能亮牌")
	//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "小于12张牌不能亮牌")
	//	return
	//}

	if person.Gang == 1 || person.Peng == 1 || person.Hu == 1 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "GameView有操作不让出牌:", person.Gang, person.Peng, person.Hu)
		return
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			continue
		}
		for j := 0; j < len(self.PersonMgr[i].Card.Want); j++ {
			if self.PersonMgr[i].Card.Want[j] == card {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "GameView不能打亮牌需要的牌:", card)
				return
			}
		}
	}

	if person.Card.CardM != card {
		find := false
		for i := 0; i < len(person.Card.Card1); i++ {
			if person.Card.Card1[i] == card {
				copy(person.Card.Card1[i:], person.Card.Card1[i+1:])
				person.Card.Card1 = person.Card.Card1[:len(person.Card.Card1)-1]
				find = true
				break
			}
		}
		if !find {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "GameView出牌错误")
			return
		}
		if person.Card.CardM != 0 {
			person.Card.Card1 = append(person.Card.Card1, person.Card.CardM)
		}
	}
	person.Card.CardM = 0
	person.Card.Card2 = append(person.Card.Card2, card)
	person.So_Card = 0
	person.State = make([]Son_KWX_State, 0)
	self.LastCard = card
	self.AddRecordStep(person.Uid, 1, []int{card})

	self.BefStep = uid

	//! 出牌广播
	{
		var msg Msg_GameKWX_Step
		msg.Uid = uid
		msg.Card = card
		self.room.broadCastMsg("gamekwxstep", &msg)
	}

	if len(cardl) == 0 {
		return
	}

	if len(person.Card.CardL) != 0 {
		return
	}

	var lst LstCard
	lib.HF_DeepCopy(&lst, &person.Card.Card1)

	for _, value := range cardl {
		find := false
		for i := 0; i < len(lst); i++ {
			if lst[i] == value {
				copy(lst[i:], lst[i+1:])
				lst = lst[:len(lst)-1]
				find = true
				break
			}
		}
		if !find {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到要亮的牌:", value, lst)
			return
		}
	}

	for i := 0; i < len(want); i++ {
		hu := MahIsHu(&person.Card, want[i])
		if !hu {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "want错误:", want[i], person.Card)
			return
		}
	}

	person.Card.CardL = cardl
	person.Card.Want = want
	person.Card.Card1 = lst
	if self.ViewUid == 0 {
		self.ViewUid = uid
	}
	self.AddRecordStep(person.Uid, 3, cardl)

	//! 广播亮牌
	{
		var msg Msg_GameKWX_View
		msg.Uid = uid
		msg.CardL = cardl
		msg.Want = want
		self.room.broadCastMsg("gamekwxview", &msg)
	}

	//! 轮流判断一下其他人是否决定
	finduid := make([]int, 0)
	isend := false
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			continue
		}
		find := false
		if len(self.PersonMgr[i].Card.Want) > 0 {
			for j := 0; j < len(self.PersonMgr[i].Card.Want); j++ {
				if card == self.PersonMgr[i].Card.Want[j] {
					self.getHuFans(self.PersonMgr[i], card, false, false)
					break
				}
			}
		} else {
			self.getHuFans(self.PersonMgr[i], card, false, false)
		}

		if self.PersonMgr[i].So_Card > 1 || (self.PersonMgr[i].So_Card > 0 && (len(self.PersonMgr[i].Card.CardL) > 0 || len(person.Card.CardL) > 0)) { //! 可以胡
			self.PersonMgr[i].Hu = 1
			self.PersonMgr[i].NoHu = card
			find = true
			if len(self.PersonMgr[i].Card.CardL) > 0 {
				isend = true
			}
		} else {
			self.PersonMgr[i].So_Card = 0
			self.PersonMgr[i].State = make([]Son_KWX_State, 0)
		}

		if len(self.Mah.Card) > 0 && MahIsGang(&self.PersonMgr[i].Card, card, false) { //! 可以杠
			self.PersonMgr[i].Gang = 1
			find = true
		}

		if MahIsPeng(&self.PersonMgr[i].Card, card) { //! 可以碰
			self.PersonMgr[i].Peng = 1
			find = true
		}

		if find {
			finduid = append(finduid, i)
		}
	}

	if isend {
		self.OnEnd()
		return
	}

	if len(finduid) == 0 {
		if self.room.Type == GAMETYPE_KWX_SZ { //! 判断杀马
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Uid == uid {
					continue
				}

				if len(self.PersonMgr[i].Card.CardL) > 0 {
					num := 0
					for j := 0; j < len(self.PersonMgr[i].Card.CardL); j++ {
						if self.PersonMgr[i].Card.CardL[j] == card {
							num++
						}
					}
					if num != 3 {
						continue
					}

					person.So_Gang -= 2 * int(math.Pow(float64(2), float64(self.Gang)))
					person.Total -= 2 * int(math.Pow(float64(2), float64(self.Gang)))
					person.Card.Card2 = person.Card.Card2[:len(person.Card.Card2)-1]
					self.PersonMgr[i].Card.CardMG = append(self.PersonMgr[i].Card.CardMG, card)
					self.AddRecordStep(self.PersonMgr[i].Uid, 6, []int{card})
					self.PersonMgr[i].So_Gang += 2 * int(math.Pow(float64(2), float64(self.Gang)))
					self.PersonMgr[i].Total += 2 * int(math.Pow(float64(2), float64(self.Gang)))
					self.LastCard = 0
					self.Gang++
					self.GangUid = self.PersonMgr[i].Uid
					for j := 0; j < len(self.PersonMgr[i].Card.CardL); {
						if self.PersonMgr[i].Card.CardL[j] == card {
							copy(self.PersonMgr[i].Card.CardL[j:], self.PersonMgr[i].Card.CardL[j+1:])
							self.PersonMgr[i].Card.CardL = self.PersonMgr[i].Card.CardL[:len(self.PersonMgr[i].Card.CardL)-1]
						} else {
							j++
						}
					}

					for j := 0; j < len(self.PersonMgr[i].Card.Want); {
						self.getHuFans(self.PersonMgr[i], self.PersonMgr[i].Card.Want[j], true, true)
						if self.PersonMgr[i].So_Card == 0 {
							copy(self.PersonMgr[i].Card.Want[j:], self.PersonMgr[i].Card.Want[j+1:])
							self.PersonMgr[i].Card.Want = self.PersonMgr[i].Card.Want[:len(self.PersonMgr[i].Card.Want)-1]
						} else {
							j++
						}
						self.PersonMgr[i].So_Card = 0
						self.PersonMgr[i].State = make([]Son_KWX_State, 0)
					}

					var msg Msg_GameKWX_Kill
					msg.Uid = self.PersonMgr[i].Uid
					msg.Card = card
					msg.View = true
					msg.Want = self.PersonMgr[i].Card.Want
					self.room.broadCastMsg("gamekwx_kill", &msg)

					self.SendTotal()

					self.GameMoCard(self.PersonMgr[i].Uid, true)

					return
				}
			}
		}

		self.Gang = 0
		self.GangUid = 0
		self.GameMoCard(self.GetNextUid(), true)
	} else {
		for i := 0; i < len(finduid); i++ {
			self.SendOperator(self.PersonMgr[finduid[i]].Uid, self.PersonMgr[finduid[i]].Hu, self.PersonMgr[finduid[i]].Gang, self.PersonMgr[finduid[i]].Peng)
		}
	}
}

//////////////////////////////////////////////////////////////////
//! 擦杠
func (self *Game_KWX) CaGangCard(person *Game_KWX_Person, card int) bool {
	for i := 0; i < len(person.Card.CardP); i++ {
		if person.Card.CardP[i] == card { //! 擦杠
			person.Gang = 0
			person.Card.CardCG = append(person.Card.CardCG, card)
			self.AddRecordStep(person.Uid, 7, []int{card})
			copy(person.Card.CardP[i:], person.Card.CardP[i+1:])
			person.Card.CardP = person.Card.CardP[:len(person.Card.CardP)-1]
			person.Card.Card2 = person.Card.Card2[:len(person.Card.Card2)-1]
			person.So_Gang += 2 * int(math.Pow(float64(2), float64(self.Gang)))
			person.Total += 2 * int(math.Pow(float64(2), float64(self.Gang)))
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Uid == person.Uid {
					continue
				}
				self.PersonMgr[i].So_Gang -= 1 * int(math.Pow(float64(2), float64(self.Gang)))
				self.PersonMgr[i].Total -= 1 * int(math.Pow(float64(2), float64(self.Gang)))
			}
			self.Gang++
			self.LastCard = 0
			self.BefStep = 0 //! 擦杠之后,清理掉最后出牌的人
			self.GangUid = person.Uid

			var msg Msg_GameKWX_Gang
			msg.Uid = person.Uid
			msg.Card = card
			msg.View = true
			self.room.broadCastMsg("gamekwx_cagang", &msg)

			self.SendTotal()

			return true
		}
	}
	return false
}

//! 将要擦杠
func (self *Game_KWX) WillCaGang(person *Game_KWX_Person, card int) {
	for i := 0; i < len(person.Card.CardP); i++ {
		if person.Card.CardP[i] == card { //! 擦杠
			person.Card.CardCG = append(person.Card.CardCG, card)
			copy(person.Card.CardP[i:], person.Card.CardP[i+1:])
			person.Card.CardP = person.Card.CardP[:len(person.Card.CardP)-1]
			person.Card.Card2 = person.Card.Card2[:len(person.Card.Card2)-1]
			self.LastCard = 0

			var msg Msg_GameKWX_Gang
			msg.Uid = person.Uid
			msg.Card = card
			msg.View = true
			self.room.broadCastMsg("gamekwx_cagang", &msg)
		}
	}
}

//! 完成擦杠
func (self *Game_KWX) FinishCaGang(person *Game_KWX_Person, card int) {
	person.Gang = 0
	self.BefStep = 0 //! 擦杠之后,清理掉最后出牌的人
	self.AddRecordStep(person.Uid, 7, []int{card})
	person.So_Gang += 2 * int(math.Pow(float64(2), float64(self.Gang)))
	person.Total += 2 * int(math.Pow(float64(2), float64(self.Gang)))
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			continue
		}
		self.PersonMgr[i].So_Gang -= 1 * int(math.Pow(float64(2), float64(self.Gang)))
		self.PersonMgr[i].Total -= 1 * int(math.Pow(float64(2), float64(self.Gang)))
	}
	self.Gang++
	self.GangUid = person.Uid
	self.SendTotal()
}

//! 终止擦杠
func (self *Game_KWX) StopCaGang(person *Game_KWX_Person, card int) {
	person.Card.CardCG = person.Card.CardCG[:len(person.Card.CardCG)-1]
	person.Card.CardP = append(person.Card.CardP, card)
	person.Card.Card2 = append(person.Card.Card2, card)
	self.LastCard = card
}

//! 杠牌
func (self *Game_KWX) GangCard(person *Game_KWX_Person, card int) {
	person.Gang = 0
	num := 0
	if self.CurStep == person.Uid { //! 暗杠
		person.Card.CardAG = append(person.Card.CardAG, card)
		self.AddRecordStep(person.Uid, 5, []int{card})
		bs := 1
		if person.Card.CardM == card { //! 杠的牌是摸的牌
			bs = int(math.Pow(float64(2), float64(self.Gang)))
		}
		person.So_Gang += 4 * bs
		person.Total += 4 * bs
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Uid == person.Uid {
				continue
			}
			self.PersonMgr[i].So_Gang -= 2 * bs
			self.PersonMgr[i].Total -= 2 * bs
		}
		num = 4

		var msg Msg_GameKWX_Gang
		msg.Uid = person.Uid
		msg.Card = card
		msg.View = false
		self.room.broadCastMsg("gamekwx_gang", &msg)

		if card == person.Card.CardM || self.Gang == 0 {
			self.Gang++
		}
	} else { //! 明杠
		_person := self.GetPerson(self.BefStep)
		_person.So_Gang -= 2 * int(math.Pow(float64(2), float64(self.Gang)))
		_person.Total -= 2 * int(math.Pow(float64(2), float64(self.Gang)))
		person.So_Gang += 2 * int(math.Pow(float64(2), float64(self.Gang)))
		person.Total += 2 * int(math.Pow(float64(2), float64(self.Gang)))

		_person.Card.Card2 = _person.Card.Card2[:len(_person.Card.Card2)-1]

		person.Card.CardMG = append(person.Card.CardMG, card)
		self.AddRecordStep(person.Uid, 6, []int{card})
		num = 3
		self.LastCard = 0

		var msg Msg_GameKWX_Gang
		msg.Uid = person.Uid
		msg.Card = card
		msg.View = true
		self.room.broadCastMsg("gamekwx_gang", &msg)

		self.Gang++
	}
	if person.Card.CardM == card {
		person.Card.CardM = 0
		num--
	}
	//! 扣牌
	if num > 0 {
		for i := 0; i < len(person.Card.Card1); {
			if person.Card.Card1[i] == card {
				copy(person.Card.Card1[i:], person.Card.Card1[i+1:])
				person.Card.Card1 = person.Card.Card1[:len(person.Card.Card1)-1]
				num--
				if num <= 0 {
					break
				}
			} else {
				i++
			}
		}
	}
	if num > 0 {
		for i := 0; i < len(person.Card.CardL); {
			if person.Card.CardL[i] == card {
				copy(person.Card.CardL[i:], person.Card.CardL[i+1:])
				person.Card.CardL = person.Card.CardL[:len(person.Card.CardL)-1]
				num--
				if num <= 0 {
					break
				}
			} else {
				i++
			}
		}
	}
	self.GangUid = person.Uid

	self.SendTotal()
}

//! 碰牌
func (self *Game_KWX) PengCard(person *Game_KWX_Person) {
	person.Peng = 0
	person.Card.CardP = append(person.Card.CardP, self.LastCard)
	self.AddRecordStep(person.Uid, 2, []int{self.LastCard})
	//! 扣牌
	num := 2
	for i := 0; i < len(person.Card.Card1); {
		if person.Card.Card1[i] == self.LastCard {
			copy(person.Card.Card1[i:], person.Card.Card1[i+1:])
			person.Card.Card1 = person.Card.Card1[:len(person.Card.Card1)-1]
			num--
			if num <= 0 {
				break
			}
		} else {
			i++
		}
	}

	_person := self.GetPerson(self.BefStep)
	_person.Card.Card2 = _person.Card.Card2[:len(_person.Card.Card2)-1]

	var msg Msg_GameKWX_Peng
	msg.Uid = person.Uid
	msg.Card = self.LastCard
	self.room.broadCastMsg("gamekwx_peng", &msg)

	self.LastCard = 0
	self.Gang = 0
	self.GangUid = 0
}

//! 得到听胡番数
func (self *Game_KWX) getTiHuFans(person *Game_KWX_Person) {
	fan := 0
	kan := 0
	if len(person.Card.CardL) > 0 { //! 得到亮牌的胡牌番数
		for i := 0; i < len(person.Card.Want); i++ {
			self.getHuFans(person, person.Card.Want[i], true, true)
			if person.So_Card > fan {
				fan = person.So_Card
				kan = person.Kan
			}
		}
		person.So_Card = fan
		person.Kan = kan
		return
	}

	//! 得到听胡的胡牌番数
	fan = 0
	kan = 0
	for i := 0; i < len(person.Card.Card1); i++ {
		self.getHuFans(person, person.Card.Card1[i], true, true)
		if person.So_Card > fan {
			fan = person.So_Card
			kan = person.Kan
		}
	}
	for i := 0; i < len(person.Card.Card1); i++ {
		if person.Card.Card1[i]/10 >= 2 { //! 去掉万和字的顺子
			continue
		}
		if (person.Card.Card1[i]-1)%10 != 0 { //! 去掉1以下
			self.getHuFans(person, person.Card.Card1[i]-1, true, true)
			if person.So_Card > fan {
				fan = person.So_Card
				kan = person.Kan
			}
		}
		if (person.Card.Card1[i]+1)%10 != 0 { //! 去掉9以上
			self.getHuFans(person, person.Card.Card1[i]+1, true, true)
			if person.So_Card > fan {
				fan = person.So_Card
				kan = person.Kan
			}
		}
	}

	person.So_Card = fan
	person.Kan = kan
}

//! 得到番数
func (self *Game_KWX) getHuFans(person *Game_KWX_Person, card int, selfget bool, ti bool) {
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "///////////////////////////////////////")
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "判断胡牌:", person.Uid)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, lib.HF_JtoA(&person.Card))
	lib.GetLogMgr().Output(lib.LOG_DEBUG, card)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "///////////////////////////////////////")
	if !selfget && person.NoHu == card {
		person.So_Card = 0
		person.State = make([]Son_KWX_State, 0)
		person.Kan = 0
		return
	}

	ptype, jiang, kan := MahIsHuByKWX(&person.Card, card)
	if ptype == 0 {
		person.So_Card = 0
		person.State = make([]Son_KWX_State, 0)
		person.Kan = kan
		return
	}

	score := 0
	switch ptype {
	case TYPE_HU_7DUI:
		person.AddState(9, 4)
		score = 4
	case TYPE_HU_7DUI1:
		person.AddState(13, 8)
		score = 8
	case TYPE_HU_7DUI2:
		person.AddState(15, 16)
		score = 16
	case TYPE_HU_7DUI3:
		score = 32
	case TYPE_HU_KWX:
		if self.GetParam(TYPE_KWX) == 0 {
			person.AddState(2, 2)
			score = 2
		} else {
			person.AddState(2, 4)
			score = 4
		}
	case TYPE_HU_PPH:
		if self.GetParam(TYPE_PPH) == 0 {
			person.AddState(3, 2)
			score = 2
		} else {
			person.AddState(3, 4)
			score = 4
		}
	case TYPE_HU_PH:
		person.AddState(1, 1)
		score = 1
	}

	//! 清一色
	allcolor := MahIsAllColor(&person.Card, card)
	if allcolor {
		person.AddState(8, 4)
		score *= 4
	}

	//! 4归
	fourg := 0
	if ptype != TYPE_HU_7DUI && ptype != TYPE_HU_7DUI1 && ptype != TYPE_HU_7DUI2 && ptype != TYPE_HU_7DUI3 {
		if self.GetParam(TYPE_PD) == 0 { //! 半频道
			if MahIsMing4(&person.Card, card) {
				person.AddState(4, 2)
				score *= 2
				fourg++
			} else if MahIsAn4(&person.Card, card) {
				person.AddState(10, 4)
				score *= 4
				fourg++
			}
		} else { //! 全频道
			ming, an := MahHas4Num(&person.Card, card)
			for i := 0; i < ming; i++ {
				person.AddState(4, 2)
				score *= 2
				fourg++
			}
			for i := 0; i < an; i++ {
				person.AddState(10, 4)
				score *= 4
				fourg++
			}
		}
	}

	//! 杠翻
	if allcolor && ptype != TYPE_HU_7DUI && ptype != TYPE_HU_7DUI1 && ptype != TYPE_HU_7DUI2 && ptype != TYPE_HU_7DUI3 { //! 清一色&&不是7对
		num1, num2 := MahHas4Num(&person.Card, card)
		num := num1 + num2 - fourg
		if num > 0 {
			bs := int(math.Pow(float64(2), float64(num)))
			person.AddState(8, 4*bs)
			score *= bs
		}
	}

	//! 手抓一
	if MahIsZOne(&person.Card, card) {
		person.AddState(11, 2)
		score *= 2
	}

	//! 亮牌
	//if len(person.Card.CardL) > 0 {
	//	score *= 2
	//}

	yuannum := MahYuanNum(&person.Card, card)
	if yuannum == 3 { //! 大三元
		person.AddState(14, 8)
		score *= 8
	} else if yuannum == 2 && jiang/10 == 3 { //! 小三元
		person.AddState(12, 4)
		score *= 4
	}

	if !ti {
		//! 杠开和杠上炮
		if self.GangUid != 0 {
			if selfget && self.GangUid == person.Uid { //! 自摸
				if self.GetParam(TYPE_PPH) == 0 {
					person.AddState(5, 2)
					score *= 2
				} else {
					person.AddState(5, 4)
					score *= 4
				}
			} else if !selfget && self.GangUid == self.BefStep {
				if self.GetParam(TYPE_PPH) == 0 {
					person.AddState(6, 2)
					score *= 2
				} else {
					person.AddState(6, 4)
					score *= 4
				}
			}
		}

		//! 海底捞
		if self.GetParam(TYPE_PD) == 1 {
			if len(self.Mah.Card) == 0 {
				score *= 2
				if selfget {
					person.AddState(38, 2)
				} else {
					person.AddState(39, 2)
				}
			}
		}
	}

	person.So_Card = score
	person.Kan = kan
	person.Kan += len(person.Card.CardAG)
	person.Kan += len(person.Card.CardCG)
	person.Kan += len(person.Card.CardMG)

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "胡牌分数:", score)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, lib.HF_JtoA(&person.Card))
	//lib.GetLogMgr().Output("胡牌分数:", score, ",", person.Card, ",", card)
}

func (self *Game_KWX) getInfo(uid int64) *Msg_GameKWX_Info {
	var msg Msg_GameKWX_Info
	msg.Begin = self.room.Begin
	msg.CurStep = self.CurStep
	msg.BefStep = self.BefStep
	msg.Hz = self.Hz
	msg.Ma = self.Ma
	msg.Sl = self.SL
	msg.LastCard = self.LastCard
	if self.Mah == nil {
		msg.Num = 0
	} else {
		msg.Num = len(self.Mah.Card)
	}
	for i := 0; i < len(self.PersonMgr); i++ {
		var son Son_GameKWX_Info
		son.Uid = self.PersonMgr[i].Uid
		son.Deal = self.PersonMgr[i].Deal
		son.Piao = self.PersonMgr[i].Piao
		if son.Uid == uid || !msg.Begin {
			son.Card = self.PersonMgr[i].Card
		} else {
			for _, _ = range self.PersonMgr[i].Card.Card1 {
				son.Card.Card1 = append(son.Card.Card1, 0)
			}
			son.Card.Card2 = self.PersonMgr[i].Card.Card2
			son.Card.CardAG = self.PersonMgr[i].Card.CardAG
			son.Card.CardCG = self.PersonMgr[i].Card.CardCG
			son.Card.CardL = self.PersonMgr[i].Card.CardL
			son.Card.CardM = 0
			son.Card.CardP = self.PersonMgr[i].Card.CardP
			son.Card.CardMG = self.PersonMgr[i].Card.CardMG
			son.Card.Want = self.PersonMgr[i].Card.Want
		}
		son.So_Card = self.PersonMgr[i].So_Card
		son.So_Gang = self.PersonMgr[i].So_Gang
		son.So_Other = self.PersonMgr[i].So_Other
		son.Total = self.PersonMgr[i].Total
		son.Ready = self.PersonMgr[i].Ready
		son.Peng = self.PersonMgr[i].Peng
		son.Gang = self.PersonMgr[i].Gang
		son.Hu = self.PersonMgr[i].Hu
		son.State = self.PersonMgr[i].State
		son.Num = len(self.PersonMgr[i].Card.Card1) + len(self.PersonMgr[i].Card.CardL)
		if self.PersonMgr[i].Card.CardM != 0 {
			son.Num++
		}
		msg.Info = append(msg.Info, son)
	}
	return &msg
}

func (self *Game_KWX) SendOperator(uid int64, hu int, gang int, peng int) {
	var msg Msg_GameKWX_Operator
	msg.Hu = hu
	msg.Gang = gang
	msg.Peng = peng
	self.room.SendMsg(uid, "gamekwx_operator", &msg)
}

//! 同步总分
func (self *Game_KWX) SendTotal() {
	var msg Msg_GameKWX_Total
	for i := 0; i < len(self.PersonMgr); i++ {
		msg.Info = append(msg.Info, Son_GameKWX_Total{self.PersonMgr[i].Uid, self.PersonMgr[i].Total})
	}
	self.room.broadCastMsg("gamekwx_total", &msg)
}

func (self *Game_KWX) AddRecordStep(uid int64, _type int, card []int) {
	self.Record.Step = append(self.Record.Step, Son_Rec_GameKWX_Step{uid, _type, card})
}

func (self *Game_KWX) OnTime() {

}

func (self *Game_KWX) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_KWX) OnIsBets(uid int64) bool {
	return false
}

//! 结算所有人
func (self *Game_KWX) OnBalance() {
}
