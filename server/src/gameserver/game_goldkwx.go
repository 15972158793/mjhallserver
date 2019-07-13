package gameserver

import (
	"fmt"
	"lib"
	"math"
	"staticfunc"
	"time"
)

//! 金币场卡五星
type Msg_GameGoldKWX_Info struct {
	Begin    bool                   `json:"begin"`    //! 是否开始
	Info     []Son_GameGoldKWX_Info `json:"info"`     //! 人的info
	Num      int                    `json:"num"`      //! 剩余数量
	CurStep  int64                  `json:"curstep"`  //! 这局谁出
	BefStep  int64                  `json:"befstep"`  //! 上局谁出
	Ma       []int                  `json:"ma"`       //! 买马
	Hz       bool                   `json:"hz"`       //! 荒庄
	Sl       [2]int                 `json:"sl"`       //! 上楼
	LastCard int                    `json:"lastcard"` //! 最后的牌
}

type Son_GameGoldKWX_Info struct {
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
	Trust    bool            `json:"trust"`
}

type Game_GoldKWX_Person struct {
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
	Gold     int             `json:"gold"`
	NoHu     int             `json:"nohu"`
	Trust    bool            `json:"trust"` //! 是否托管

	need []int //! 预定的牌
	hz   int   //! 荒庄规则
}

func (self *Game_GoldKWX_Person) Init() {
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

//! 同步金币
func (self *Game_GoldKWX_Person) SynchroGold(gold int) {
	self.Total += (gold - self.Gold)
	self.Gold = gold
}

//! 设置托管
func (self *Game_GoldKWX_Person) SetTrust(trust bool) {
	if self.Trust == trust {
		return
	}

	self.Trust = trust

	var msg Msg_GameDeal
	msg.Uid = self.Uid
	msg.Ok = trust

	person := GetPersonMgr().GetPerson(self.Uid)
	if person != nil {
		person.SendMsg("gametrust", &msg)
	}
}

//! 是否托管
func (self *Game_GoldKWX_Person) IsTrush(_time int64) bool {
	if self.Trust {
		return true
	}

	if time.Now().Unix() < _time {
		return false
	}

	self.SetTrust(true)

	return true
}

func (self *Game_GoldKWX_Person) AddState(id, score int) {
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

type Game_GoldKWX struct {
	PersonMgr []*Game_GoldKWX_Person `json:"personmgr"`
	Mah       *MahMgr                `json:"mah"`      //! 剩余
	CurStep   int64                  `json:"curstep"`  //! 这局谁出
	BefStep   int64                  `json:"befstep"`  //! 上局谁出
	Winer     int64                  `json:"winer"`    //! 上局谁赢
	Gang      int                    `json:"gang"`     //! 当前局第几杠
	LastCard  int                    `json:"lastcard"` //! 最后一张牌
	ViewUid   int64                  `json:"viewuid"`  //! 第一家亮牌
	GangUid   int64                  `json:"ganguid"`  //! 最近杠牌
	//Record    *staticfunc.Rec_Gold_Info `json:"record"`   //! 卡五星记录
	//IsRecord bool   `json:"isrecord"` //! 是否记录了
	Ma    []int  `json:"ma"` //! 买马
	Hz    bool   `json:"hz"` //! 荒庄
	SL    [2]int `json:"sl"` //! 上楼
	Type  int    `json:"type"`
	DF    int    `json:"df"`   //! 底分
	Time  int64  `json:"time"` //! 自动选择时间
	Error int    `json:"error"`
	//Piao      bool                   `json:"piao"` //! 是否选了漂

	room *Room
}

func NewGame_GoldKWX() *Game_GoldKWX {
	game := new(Game_GoldKWX)
	game.PersonMgr = make([]*Game_GoldKWX_Person, 0)

	return game
}

func (self *Game_GoldKWX) GetParam(_type int) int {
	return self.room.Param1 % int(math.Pow(10.0, float64(TYPE_MAX-_type))) / int(math.Pow(10.0, float64(TYPE_MAX-_type-1)))
}

func (self *Game_GoldKWX) SetParam(_type int, value int) {
	high := self.room.Param1 / int(math.Pow(10.0, float64(TYPE_MAX-_type)))
	next := self.room.Param1 % int(math.Pow(10.0, float64(TYPE_MAX-_type-1)))
	self.room.Param1 = high*int(math.Pow(10.0, float64(TYPE_MAX-_type))) + value*int(math.Pow(10.0, float64(TYPE_MAX-_type-1))) + next
}

func (self *Game_GoldKWX) GetParam2(_type int) int {
	return self.room.Param2 % int(math.Pow(10.0, float64(TYPE_MAX2-_type))) / int(math.Pow(10.0, float64(TYPE_MAX2-_type-1)))
}

func (self *Game_GoldKWX) SetParam2(_type int, value int) {
	high := self.room.Param2 / int(math.Pow(10.0, float64(TYPE_MAX2-_type)))
	next := self.room.Param2 % int(math.Pow(10.0, float64(TYPE_MAX2-_type-1)))
	self.room.Param2 = high*int(math.Pow(10.0, float64(TYPE_MAX2-_type))) + value*int(math.Pow(10.0, float64(TYPE_MAX2-_type-1))) + next
}

func (self *Game_GoldKWX) GetPerson(uid int64) *Game_GoldKWX_Person {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			return self.PersonMgr[i]
		}
	}

	return nil
}

//! 得到下一个uid
func (self *Game_GoldKWX) GetNextUid() int64 {
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

func (self *Game_GoldKWX) OnInit(room *Room) {
	self.room = room

	if self.room.Type%10 == 0 { //! 买马玩法
		self.Type = GAMETYPE_KWX_XG
		self.SetParam(TYPE_KWX, 1)
		self.SetParam(TYPE_PPH, 1)
		self.SetParam(TYPE_DL, 0)
		self.SetParam(TYPE_ZDFS, 1)
		self.SetParam(TYPE_PIAO, 1)
		self.SetParam(TYPE_SK, 1)
		self.SetParam(TYPE_MM, 1)
		self.SetParam(TYPE_PH, 1)
		self.SetParam(TYPE_PD, 0)
		self.SetParam2(TYPE_SL, 0)
		self.SetParam2(TYPE_DLDF, 0)
		self.SetParam2(TYPE_PQMM, 0)
	} else if self.room.Type%10 == 1 { //! 不买马玩法
		self.Type = GAMETYPE_KWX_XG
		self.SetParam(TYPE_KWX, 1)
		self.SetParam(TYPE_PPH, 1)
		self.SetParam(TYPE_DL, 0)
		self.SetParam(TYPE_ZDFS, 1)
		self.SetParam(TYPE_PIAO, 1)
		self.SetParam(TYPE_SK, 1)
		self.SetParam(TYPE_MM, 0)
		self.SetParam(TYPE_PH, 1)
		self.SetParam(TYPE_PD, 0)
		self.SetParam2(TYPE_SL, 0)
		self.SetParam2(TYPE_DLDF, 0)
		self.SetParam2(TYPE_PQMM, 0)
	} else if self.room.Type%10 == 2 { //! 全频道买马
		self.Type = GAMETYPE_KWX_XY
		self.SetParam(TYPE_KWX, 0)
		self.SetParam(TYPE_PPH, 0)
		self.SetParam(TYPE_DL, 0)
		self.SetParam(TYPE_ZDFS, 0)
		self.SetParam(TYPE_PIAO, 1)
		self.SetParam(TYPE_SK, 0)
		self.SetParam(TYPE_MM, 1)
		self.SetParam(TYPE_PH, 1)
		self.SetParam(TYPE_PD, 1)
		self.SetParam2(TYPE_SL, 0)
		self.SetParam2(TYPE_DLDF, 0)
		self.SetParam2(TYPE_PQMM, 0)
	} else if self.room.Type%10 == 3 { //! 全频道不买马
		self.Type = GAMETYPE_KWX_XY
		self.SetParam(TYPE_KWX, 0)
		self.SetParam(TYPE_PPH, 0)
		self.SetParam(TYPE_DL, 0)
		self.SetParam(TYPE_ZDFS, 0)
		self.SetParam(TYPE_PIAO, 1)
		self.SetParam(TYPE_SK, 0)
		self.SetParam(TYPE_MM, 0)
		self.SetParam(TYPE_PH, 1)
		self.SetParam(TYPE_PD, 1)
		self.SetParam2(TYPE_SL, 0)
		self.SetParam2(TYPE_DLDF, 0)
		self.SetParam2(TYPE_PQMM, 0)
	} else { //! 上楼玩法
		self.Type = GAMETYPE_KWX_SY
		self.SetParam(TYPE_KWX, 1)
		self.SetParam(TYPE_PPH, 0)
		self.SetParam(TYPE_DL, 0)
		self.SetParam(TYPE_ZDFS, 0)
		self.SetParam(TYPE_PIAO, 1)
		self.SetParam(TYPE_SK, 0)
		self.SetParam(TYPE_MM, 0)
		self.SetParam(TYPE_PH, 1)
		self.SetParam(TYPE_PD, 1)
		self.SetParam2(TYPE_SL, 1)
		self.SetParam2(TYPE_DLDF, 0)
		self.SetParam2(TYPE_PQMM, 0)
	}

	self.DF = staticfunc.GetCsvMgr().GetDF(self.room.Type)
}

func (self *Game_GoldKWX) OnRobot(robot *lib.Robot) {

}

func (self *Game_GoldKWX) OnSendInfo(person *Person) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			self.PersonMgr[i].SynchroGold(person.Gold)
			person.SendMsg("gamekwxinfo", self.getInfo(person.Uid))
			return
		}
	}

	_person := new(Game_GoldKWX_Person)
	_person.Init()
	_person.Uid = person.Uid
	_person.Ready = false
	_person.Piao = -1
	_person.Total = person.Gold
	_person.Gold = person.Gold
	_person.Trust = false
	self.PersonMgr = append(self.PersonMgr, _person)

	if len(self.PersonMgr) >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 人满了进入一分钟倒计时
		self.SetTime(5)
	}

	person.SendMsg("gamekwxinfo", self.getInfo(person.Uid))
}

func (self *Game_GoldKWX) OnMsg(msg *RoomMsg) {
	switch msg.Head {
	case "synchrogold": //! 同步金币
		person := self.GetPerson(msg.V.(*staticfunc.Msg_SynchroGold).Uid)
		if person != nil {
			person.SynchroGold(msg.V.(*staticfunc.Msg_SynchroGold).Gold)
			self.SendTotal()
		}
		self.room.flush()
	//case "gameready": //! 游戏准备
	//	self.GameReady(msg.Uid)
	//	self.room.flush()
	case "gamebets": //! 加飘
		person := self.GetPerson(msg.Uid)
		if person != nil {
			person.SetTrust(false)
		}
		self.GameBets(msg.Uid, msg.V.(*Msg_GameBets).Bets)
		self.room.flush()
	case "gamestep": //! 出牌
		person := self.GetPerson(msg.Uid)
		if person != nil {
			person.SetTrust(false)
		}
		self.GameStep(msg.Uid, msg.V.(*Msg_GameStep).Card)
		self.room.flush()
	case "gamecagang":
		person := self.GetPerson(msg.Uid)
		if person != nil {
			person.SetTrust(false)
		}
		self.GameCaGang(msg.Uid, msg.V.(*Msg_GameStep).Card)
		self.room.flush()
	case "gamepeng": //! 碰
		person := self.GetPerson(msg.Uid)
		if person != nil {
			person.SetTrust(false)
		}
		self.GamePeng(msg.Uid)
		self.room.flush()
	case "gamegang": //! 杠
		person := self.GetPerson(msg.Uid)
		if person != nil {
			person.SetTrust(false)
		}
		self.GameGang(msg.Uid, msg.V.(*Msg_GameStep).Card)
		self.room.flush()
	case "gamehu": //! 胡
		person := self.GetPerson(msg.Uid)
		if person != nil {
			person.SetTrust(false)
		}
		self.GameHu(msg.Uid)
		self.room.flush()
	case "gameguo": //! 过
		person := self.GetPerson(msg.Uid)
		if person != nil {
			person.SetTrust(false)
		}
		self.GameGuo(msg.Uid)
		self.room.flush()
	case "gamekwxview":
		person := self.GetPerson(msg.Uid)
		if person != nil {
			person.SetTrust(false)
		}
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
	case "gametrust": //! 托管
		self.GameTrust(msg.Uid, msg.V.(*Msg_GameDeal).Ok)
	}
}

func (self *Game_GoldKWX) OnBegin() {
	if self.room.IsBye() {
		return
	}

	self.room.SetBegin(true)
	self.Winer = self.room.Uid[lib.HF_GetRandom(len(self.room.Uid))]

	//! 扣除底分
	for i := 0; i < len(self.PersonMgr); i++ {
		cost := int(math.Ceil(float64(self.DF) * 35.0 / 100.0))
		self.PersonMgr[i].Total -= cost
		GetServer().SqlAgentGoldLog(self.PersonMgr[i].Uid, cost, self.room.Type)
		GetServer().SqlAgentBillsLog(self.PersonMgr[i].Uid, cost, self.room.Type)
	}
	self.SendTotal()

	self.Mah = NewMah_KWX()
	self.BefStep = 0
	self.ViewUid = 0
	self.Gang = 0
	self.GangUid = 0
	self.Hz = false
	self.Ma = make([]int, 0)
	//self.Record = new(staticfunc.Rec_Gold_Info)
	//self.IsRecord = false
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

		//! 记录
		//var rc_person staticfunc.Son_Rec_Gold_Person
		//rc_person.Uid = self.PersonMgr[i].Uid
		//rc_person.Name = self.room.GetName(rc_person.Uid)
		//rc_person.Head = self.room.GetHead(rc_person.Uid)
		//self.Record.Info = append(self.Record.Info, rc_person)
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

	self.SetTime(60)
}

func (self *Game_GoldKWX) OnEnd() {
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
		var no1person *Game_GoldKWX_Person = nil //! 不是第一个亮牌的人
		for _, value := range self.PersonMgr {
			self.getTiHuFans(value)
			if len(value.Card.CardL) > 0 {
				if self.ViewUid == value.Uid {
					value.hz = 1
				} else {
					if self.Type == GAMETYPE_KWX_SZ { //! 随州玩法第一家亮牌的赔胡
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
		if self.Type == GAMETYPE_KWX_SZ && no1person != nil { //! 纠正随州的赔胡
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
						if (self.Type == GAMETYPE_KWX_XY || self.Type == GAMETYPE_KWX_SY) && len(self.PersonMgr[i].Card.CardL) == 0 { //! 襄阳和十堰进入包胡
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

	//! 计算一共赔了多少
	lostnum := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		num := (self.PersonMgr[i].So_Card + self.PersonMgr[i].So_Other) * self.DF
		if num >= 0 {
			continue
		}
		if self.PersonMgr[i].Total+num < 0 {
			num = -self.PersonMgr[i].Total
		}
		self.PersonMgr[i].Total = self.PersonMgr[i].Total + num
		lostnum += -num
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		num := (self.PersonMgr[i].So_Card + self.PersonMgr[i].So_Other) * self.DF
		if num <= 0 {
			continue
		}
		if lostnum >= num { //! 足够赔
			self.PersonMgr[i].Total += num
			lostnum -= num
		} else { //! 不够赔
			self.PersonMgr[i].Total += lostnum
			break
		}
	}

	var record staticfunc.Rec_Gold_Info
	record.Time = time.Now().Unix()
	record.GameType = self.room.Type

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
		son.So_Card = self.PersonMgr[i].So_Card * self.DF
		son.So_Gang = self.PersonMgr[i].So_Gang * self.DF
		son.So_Other = self.PersonMgr[i].So_Other * self.DF
		son.Total = self.PersonMgr[i].Total
		son.State = self.PersonMgr[i].State
		self.PersonMgr[i].Hu = score[i]
		son.Hu = score[i]
		msg.Info = append(msg.Info, son)
		self.room.Param[i] = self.PersonMgr[i].Total

		var rec staticfunc.Son_Rec_Gold_Person
		rec.Uid = self.PersonMgr[i].Uid
		rec.Name = self.room.GetName(self.PersonMgr[i].Uid)
		rec.Head = self.room.GetHead(self.PersonMgr[i].Uid)
		rec.Score = (self.PersonMgr[i].So_Card + self.PersonMgr[i].So_Gang + self.PersonMgr[i].So_Other) * self.DF
		record.Info = append(record.Info, rec)
	}
	recordinfo := lib.HF_JtoA(&record)
	for i := 0; i < len(record.Info); i++ {
		GetServer().InsertRecord(self.room.Type, record.Info[i].Uid, recordinfo, -record.Info[i].Score)
	}
	self.room.broadCastMsg("gamekwxend", &msg)

	//self.Record.Time = time.Now().Unix()
	//self.Record.GameType = self.room.Type
	//self.room.AddRecord(lib.HF_JtoA(self.Record))
	//self.IsRecord = true

	if self.room.IsBye() {
		self.OnBye()
		self.room.Bye()
		return
	}

	//self.Piao = false
	for i := 0; i < len(self.PersonMgr); i++ {
		self.PersonMgr[i].Ready = false
		self.PersonMgr[i].Piao = -1
		self.PersonMgr[i].hz = 0
	}
}

//! 托管
func (self *Game_GoldKWX) GameTrust(uid int64, ok bool) {
	person := self.GetPerson(uid)
	if person == nil {
		return
	}

	person.SetTrust(ok)
}

func (self *Game_GoldKWX) OnBye() {
	//if !self.IsRecord && self.Record != nil {
	//	self.IsRecord = true
	//}

	//var msg Msg_GameKWX_Bye
	//for i := 0; i < len(self.PersonMgr); i++ {
	//	var son Son_GameKWX_Bye
	//	son.Uid = self.PersonMgr[i].Uid
	//	son.ZM = self.PersonMgr[i].ZM
	//	son.JP = self.PersonMgr[i].JP
	//	son.DP = self.PersonMgr[i].DP
	//	son.AG = self.PersonMgr[i].AG
	//	son.MG = self.PersonMgr[i].MG
	//	son.PI = lib.HF_MaxInt(0, self.PersonMgr[i].Piao)
	//	son.Score = self.PersonMgr[i].Total
	//	msg.Info = append(msg.Info, son)
	//}
	//self.room.broadCastMsg("gamekwxbye", &msg)
}

func (self *Game_GoldKWX) OnExit(uid int64) {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			//! 退出房间同步金币
			gold := self.PersonMgr[i].Total - self.PersonMgr[i].Gold
			if gold > 0 {
				GetRoomMgr().AddCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
			} else if gold < 0 {
				GetRoomMgr().CostCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, -gold, self.room)
			}
			self.PersonMgr[i].Gold = self.PersonMgr[i].Total

			copy(self.PersonMgr[i:], self.PersonMgr[i+1:])
			self.PersonMgr = self.PersonMgr[:len(self.PersonMgr)-1]

			//! 有人退出之后取消自动操作
			self.SetTime(0)
			break
		}
	}
}

//! 准备,第一局自动准备
//func (self *Game_GoldKWX) GameReady(uid int64) {
//	if self.room.IsBye() {
//		return
//	}

//	if self.room.Begin {
//		return
//	}

//	num := 0
//	for i := 0; i < len(self.PersonMgr); i++ {
//		if self.PersonMgr[i].Uid == uid {
//			if self.PersonMgr[i].Ready {
//				return
//			} else {
//				if self.PersonMgr[i].Total < self.DF*20 { //! 携带的金币不足，踢出去
//					self.room.KickPerson(uid, 1)
//					return
//				}

//				self.PersonMgr[i].Ready = true
//				num++
//			}
//		} else if self.PersonMgr[i].Ready {
//			num++
//		}
//	}

//	if num == len(self.room.Uid) && num >= lib.HF_Atoi(self.room.csv["minnum"]) {
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始")
//		self.OnBegin()
//		return
//	}

//	var msg staticfunc.Msg_Uid
//	msg.Uid = uid
//	self.room.broadCastMsg("gameready", &msg)
//}

//! 加飘
func (self *Game_GoldKWX) GameBets(uid int64, bets int) {
	if self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "没有开始不能选飘:GameBets")
		return
	}

	if bets < 0 {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "加漂小于0:GameBets")
		return
	}

	if self.GetParam(TYPE_PIAO) == 0 {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "不能加漂:GameBets")
		return
	}

	person := GetPersonMgr().GetPerson(uid)
	if person == nil {
		return
	}
	if person.black {
		self.room.KickPerson(uid, 95)
		return
	}

	num := 0
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			if self.PersonMgr[i].Piao != -1 {
				lib.GetLogMgr().Output(lib.LOG_ERROR, "已经加漂过:GameBets")
				return
			} else {
				if self.PersonMgr[i].Total < staticfunc.GetCsvMgr().GetZR(self.room.Type) { //! 携带的金币不足，踢出去
					self.room.KickPerson(uid, 99)
					return
				}
				self.PersonMgr[i].Piao = bets
				self.PersonMgr[i].Ready = true
				num++
			}
		} else if self.PersonMgr[i].Piao != -1 {
			num++
		}
	}

	if num == len(self.PersonMgr) && len(self.PersonMgr) >= lib.HF_Atoi(self.room.csv["minnum"]) { //! 所有人都加飘
		//self.Piao = true
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
func (self *Game_GoldKWX) GameMoCard(uid int64, send bool) {
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
		if self.Type == GAMETYPE_KWX_YC { //! 应城自摸必须2番胡
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
	if person.Hu == 0 && len(person.Card.CardL) > 0 && self.Type == GAMETYPE_KWX_SZ {
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
			person.Total += 4 * int(math.Pow(float64(2), float64(self.Gang))) * self.DF
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Uid == person.Uid {
					continue
				}
				self.PersonMgr[i].So_Gang -= 2 * int(math.Pow(float64(2), float64(self.Gang)))
				self.PersonMgr[i].Total -= 2 * int(math.Pow(float64(2), float64(self.Gang))) * self.DF
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
func (self *Game_GoldKWX) GameStep(uid int64, card int) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "游戏未开始")
		return
	}

	//if !self.Piao {
	//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未加飘")
	//	return
	//}

	if card == 0 {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "GameStep(card=0):", self.room.Id)
		return
	}

	if self.CurStep != uid {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "当前不是你的局:GameStep")
		return
	}

	if self.BefStep == uid {
		self.Error++
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "已经出过牌:GameStep:", self.room.Id)
		if self.Error >= 10 {
			self.room.Bye()
		}
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "找不到person:GameStep")
		return
	}

	if person.Gang == 1 || person.Peng == 1 || person.Hu == 1 {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "有操作不让出牌:GameStep")
		return
	}

	self.Error = 0

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
			lib.GetLogMgr().Output(lib.LOG_ERROR, "亮牌后只能打摸的牌:GameStep")
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
			lib.GetLogMgr().Output(lib.LOG_ERROR, "出牌错误:GameStep")
			return
		}
		if person.Card.CardM != 0 {
			person.Card.Card1 = append(person.Card.Card1, person.Card.CardM)
		}
	}

	self.SetTime(60)

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
		if self.Type == GAMETYPE_KWX_SZ { //! 判断杀马
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
					person.Total -= 2 * int(math.Pow(float64(2), float64(self.Gang))) * self.DF
					person.Card.Card2 = person.Card.Card2[:len(person.Card.Card2)-1]
					self.PersonMgr[i].Card.CardMG = append(self.PersonMgr[i].Card.CardMG, card)
					self.AddRecordStep(self.PersonMgr[i].Uid, 6, []int{card})
					self.PersonMgr[i].So_Gang += 2 * int(math.Pow(float64(2), float64(self.Gang)))
					self.PersonMgr[i].Total += 2 * int(math.Pow(float64(2), float64(self.Gang))) * self.DF
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
func (self *Game_GoldKWX) GameCaGang(uid int64, card int) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "游戏未开始:GameCaGang")
		return
	}

	if self.CurStep != uid {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "当前不是你的局:GameCaGang")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "找不到person:GameCaGang")
		return
	}

	if person.Card.CardM != card {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "擦杠必须是自摸的牌:GameCaGang")
		return
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			continue
		}
		for j := 0; j < len(self.PersonMgr[i].Card.Want); j++ {
			if self.PersonMgr[i].Card.Want[j] == card {
				lib.GetLogMgr().Output(lib.LOG_ERROR, "不能打亮牌需要的牌:GameCaGang")
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
		lib.GetLogMgr().Output(lib.LOG_ERROR, "不可擦杠:GameCaGang")
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
func (self *Game_GoldKWX) GamePeng(uid int64) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "游戏未开始:GamePeng")
		return
	}

	if self.CurStep == uid {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "自己局不能碰:GamePeng")
		return
	}

	if self.LastCard == 0 {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "最后牌错误:GamePeng")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "找不到person:GamePeng")
		return
	}

	if person.Peng != 1 {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "不能碰:GamePeng")
		return
	}

	self.SetTime(60)

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
func (self *Game_GoldKWX) GameGang(uid int64, card int) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "游戏未开始:GameGang")
		return
	}

	//if !self.Piao {
	//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未加飘")
	//	return
	//}

	if self.CurStep != uid {
		if self.LastCard != card {
			lib.GetLogMgr().Output(lib.LOG_ERROR, "杠的不是最后一张牌:GameGang")
			return
		}
	}

	person := self.GetPerson(uid)
	if person == nil {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "找不到person:GameGang")
		return
	}

	if person.Gang != 1 {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "不能杠:GameGang")
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
			lib.GetLogMgr().Output(lib.LOG_ERROR, "card错误:GameGang")
			return
		}
	}

	self.SetTime(60)

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
func (self *Game_GoldKWX) GameHu(uid int64) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "游戏未开始:GameHu")
		return
	}

	//if !self.Piao {
	//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未加飘")
	//	return
	//}

	person := self.GetPerson(uid)
	if person == nil {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "找不到person:GameHu")
		return
	}

	if person.Hu != 1 {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "不能胡:GameHu")
		return
	}

	self.SetTime(20)

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
func (self *Game_GoldKWX) GameGuo(uid int64) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "游戏未开始:GameGuo")
		return
	}

	//if !self.Piao {
	//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未加飘")
	//	return
	//}

	person := self.GetPerson(uid)
	if person == nil {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "找不到person:GameGuo")
		return
	}

	if person.Hu == 0 && person.Peng == 0 && person.Gang == 0 {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "不能过:GameGuo")
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

func (self *Game_GoldKWX) GameView(uid int64, cardl []int, want []int, card int) {
	if !self.room.Begin {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "GameView游戏未开始")
		return
	}

	//if !self.Piao {
	//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏未加飘")
	//	return
	//}

	if card == 0 {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "GameView(card=0)")
		return
	}

	if self.CurStep != uid {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "GameView当前不是你的局")
		return
	}

	if self.BefStep == uid {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "GameView已经出过牌")
		return
	}

	person := self.GetPerson(uid)
	if person == nil {

		lib.GetLogMgr().Output(lib.LOG_ERROR, "找不到person")
		return
	}

	//if len(self.Mah.Card) <= 12 && self.GetParam(TYPE_12VIEW) == 1 {
	//	lib.GetLogMgr().Output("GameView小于12张牌不能亮牌")
	//	lib.GetLogMgr().Output(lib.LOG_DEBUG, "小于12张牌不能亮牌")
	//	return
	//}

	if person.Gang == 1 || person.Peng == 1 || person.Hu == 1 {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "GameView有操作不让出牌:", person.Gang, person.Peng, person.Hu)
		return
	}

	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			continue
		}
		for j := 0; j < len(self.PersonMgr[i].Card.Want); j++ {
			if self.PersonMgr[i].Card.Want[j] == card {
				lib.GetLogMgr().Output(lib.LOG_ERROR, "GameView不能打亮牌需要的牌:", card)
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
			lib.GetLogMgr().Output(lib.LOG_ERROR, "GameView出牌错误")
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
			lib.GetLogMgr().Output(lib.LOG_ERROR, "找不到要亮的牌:", value, lst)
			return
		}
	}

	for i := 0; i < len(want); i++ {
		hu := MahIsHu(&person.Card, want[i])
		if !hu {
			lib.GetLogMgr().Output(lib.LOG_ERROR, "want错误:", want[i], person.Card)
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
		if self.Type == GAMETYPE_KWX_SZ { //! 判断杀马
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
					person.Total -= 2 * int(math.Pow(float64(2), float64(self.Gang))) * self.DF
					person.Card.Card2 = person.Card.Card2[:len(person.Card.Card2)-1]
					self.PersonMgr[i].Card.CardMG = append(self.PersonMgr[i].Card.CardMG, card)
					self.AddRecordStep(self.PersonMgr[i].Uid, 6, []int{card})
					self.PersonMgr[i].So_Gang += 2 * int(math.Pow(float64(2), float64(self.Gang)))
					self.PersonMgr[i].Total += 2 * int(math.Pow(float64(2), float64(self.Gang))) * self.DF
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
func (self *Game_GoldKWX) CaGangCard(person *Game_GoldKWX_Person, card int) bool {
	for i := 0; i < len(person.Card.CardP); i++ {
		if person.Card.CardP[i] == card { //! 擦杠
			person.Gang = 0
			person.Card.CardCG = append(person.Card.CardCG, card)
			self.AddRecordStep(person.Uid, 7, []int{card})
			copy(person.Card.CardP[i:], person.Card.CardP[i+1:])
			person.Card.CardP = person.Card.CardP[:len(person.Card.CardP)-1]
			person.Card.Card2 = person.Card.Card2[:len(person.Card.Card2)-1]
			person.So_Gang += 2 * int(math.Pow(float64(2), float64(self.Gang)))
			person.Total += 2 * int(math.Pow(float64(2), float64(self.Gang))) * self.DF
			for i := 0; i < len(self.PersonMgr); i++ {
				if self.PersonMgr[i].Uid == person.Uid {
					continue
				}
				self.PersonMgr[i].So_Gang -= 1 * int(math.Pow(float64(2), float64(self.Gang)))
				self.PersonMgr[i].Total -= 1 * int(math.Pow(float64(2), float64(self.Gang))) * self.DF
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
func (self *Game_GoldKWX) WillCaGang(person *Game_GoldKWX_Person, card int) {
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
func (self *Game_GoldKWX) FinishCaGang(person *Game_GoldKWX_Person, card int) {
	person.Gang = 0
	self.BefStep = 0 //! 擦杠之后,清理掉最后出牌的人
	self.AddRecordStep(person.Uid, 7, []int{card})
	person.So_Gang += 2 * int(math.Pow(float64(2), float64(self.Gang)))
	person.Total += 2 * int(math.Pow(float64(2), float64(self.Gang))) * self.DF
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == person.Uid {
			continue
		}
		self.PersonMgr[i].So_Gang -= 1 * int(math.Pow(float64(2), float64(self.Gang)))
		self.PersonMgr[i].Total -= 1 * int(math.Pow(float64(2), float64(self.Gang))) * self.DF
	}
	self.Gang++
	self.GangUid = person.Uid
	self.SendTotal()
}

//! 终止擦杠
func (self *Game_GoldKWX) StopCaGang(person *Game_GoldKWX_Person, card int) {
	person.Card.CardCG = person.Card.CardCG[:len(person.Card.CardCG)-1]
	person.Card.CardP = append(person.Card.CardP, card)
	person.Card.Card2 = append(person.Card.Card2, card)
	self.LastCard = card
}

//! 杠牌
func (self *Game_GoldKWX) GangCard(person *Game_GoldKWX_Person, card int) {
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
		person.Total += 4 * bs * self.DF
		for i := 0; i < len(self.PersonMgr); i++ {
			if self.PersonMgr[i].Uid == person.Uid {
				continue
			}
			self.PersonMgr[i].So_Gang -= 2 * bs
			self.PersonMgr[i].Total -= 2 * bs * self.DF
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
		_person.Total -= 2 * int(math.Pow(float64(2), float64(self.Gang))) * self.DF
		person.So_Gang += 2 * int(math.Pow(float64(2), float64(self.Gang)))
		person.Total += 2 * int(math.Pow(float64(2), float64(self.Gang))) * self.DF

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
func (self *Game_GoldKWX) PengCard(person *Game_GoldKWX_Person) {
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
func (self *Game_GoldKWX) getTiHuFans(person *Game_GoldKWX_Person) {
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
func (self *Game_GoldKWX) getHuFans(person *Game_GoldKWX_Person, card int, selfget bool, ti bool) {
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

func (self *Game_GoldKWX) getInfo(uid int64) *Msg_GameGoldKWX_Info {
	var msg Msg_GameGoldKWX_Info
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
		var son Son_GameGoldKWX_Info
		son.Uid = self.PersonMgr[i].Uid
		son.Deal = self.PersonMgr[i].Deal
		son.Piao = self.PersonMgr[i].Piao
		if son.Uid == uid || !msg.Begin || GetServer().IsAdmin(uid, staticfunc.ADMIN_KWX) {
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
		son.So_Card = self.PersonMgr[i].So_Card * self.DF
		son.So_Gang = self.PersonMgr[i].So_Gang * self.DF
		son.So_Other = self.PersonMgr[i].So_Other * self.DF
		son.Total = self.PersonMgr[i].Total
		son.Ready = self.PersonMgr[i].Ready
		son.Peng = self.PersonMgr[i].Peng
		son.Gang = self.PersonMgr[i].Gang
		son.Hu = self.PersonMgr[i].Hu
		son.State = self.PersonMgr[i].State
		son.Num = len(self.PersonMgr[i].Card.Card1) + len(self.PersonMgr[i].Card.CardL)
		son.Trust = self.PersonMgr[i].Trust
		if self.PersonMgr[i].Card.CardM != 0 {
			son.Num++
		}
		msg.Info = append(msg.Info, son)
	}
	return &msg
}

func (self *Game_GoldKWX) SendOperator(uid int64, hu int, gang int, peng int) {
	var msg Msg_GameKWX_Operator
	msg.Hu = hu
	msg.Gang = gang
	msg.Peng = peng
	self.room.SendMsg(uid, "gamekwx_operator", &msg)
}

//! 同步总分
func (self *Game_GoldKWX) SendTotal() {
	var msg Msg_GameKWX_Total
	for i := 0; i < len(self.PersonMgr); i++ {
		self.room.Param[i] = self.PersonMgr[i].Total
		msg.Info = append(msg.Info, Son_GameKWX_Total{self.PersonMgr[i].Uid, self.PersonMgr[i].Total})
	}
	self.room.broadCastMsg("gamekwx_total", &msg)
}

func (self *Game_GoldKWX) AddRecordStep(uid int64, _type int, card []int) {
	//self.Record.Step = append(self.Record.Step, Son_Rec_GameKWX_Step{uid, _type, card})
}

func (self *Game_GoldKWX) OnTime() {
	if self.Time == 0 {
		return
	}

	//! 60秒之后自动选择
	if !self.room.Begin {
		if time.Now().Unix() < self.Time {
			return
		}
		for i := 0; i < len(self.PersonMgr); {
			if self.PersonMgr[i].Piao == -1 {
				self.room.KickPerson(self.PersonMgr[i].Uid, 98)
			} else {
				i++
			}
		}

		return
	}

	//! 判断是否有人能胡
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Hu == 1 {
			if self.PersonMgr[i].IsTrush(self.Time) {
				self.GameHu(self.PersonMgr[i].Uid)
			}
			return
		}
	}

	//! 判断是否有人能杠
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Gang == 1 {
			if self.PersonMgr[i].IsTrush(self.Time) {
				if self.CurStep != self.PersonMgr[i].Uid {
					self.GameGang(self.PersonMgr[i].Uid, self.LastCard)
				} else {
					self.GameGang(self.PersonMgr[i].Uid, self.MahIsGang(&self.PersonMgr[i].Card, self.PersonMgr[i].Card.CardM, true))
				}
			}
			return
		}
	}

	//! 判断是否有人能碰
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Peng == 1 {
			if self.PersonMgr[i].IsTrush(self.Time) {
				self.GamePeng(self.PersonMgr[i].Uid)
			}
			return
		}
	}

	//! 判断是否能出牌
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == self.CurStep {
			if self.PersonMgr[i].IsTrush(self.Time) {
				if len(self.PersonMgr[i].Card.CardL) > 0 && self.PersonMgr[i].Card.CardM != 0 { //! 亮牌了
					self.GameStep(self.PersonMgr[i].Uid, self.PersonMgr[i].Card.CardM)
				} else {
					if self.PersonMgr[i].Card.CardM != 0 && self.IsStep(self.PersonMgr[i].Uid, self.PersonMgr[i].Card.CardM) {
						self.GameStep(self.PersonMgr[i].Uid, self.PersonMgr[i].Card.CardM)
					} else {
						find := false
						for j := 0; j < len(self.PersonMgr[i].Card.Card1); j++ {
							if self.IsStep(self.PersonMgr[i].Uid, self.PersonMgr[i].Card.Card1[j]) {
								self.GameStep(self.PersonMgr[i].Uid, self.PersonMgr[i].Card.Card1[j])
								find = true
								break
							}
						}
						if !find {
							self.GameStep(self.PersonMgr[i].Uid, self.PersonMgr[i].Card.Card1[0])
						}
					}
				}
			}

			return
		}
	}
}

func (self *Game_GoldKWX) OnIsDealer(uid int64) bool {
	return false
}

func (self *Game_GoldKWX) OnIsBets(uid int64) bool {
	return false
}

//! 该牌是否能出
func (self *Game_GoldKWX) IsStep(uid int64, card int) bool {
	for i := 0; i < len(self.PersonMgr); i++ {
		if self.PersonMgr[i].Uid == uid {
			continue
		}

		for j := 0; j < len(self.PersonMgr[i].Card.Want); j++ {
			if self.PersonMgr[i].Card.Want[j] == card {
				return false
			}
		}
	}

	return true
}

//! 是否能杠
func (self *Game_GoldKWX) MahIsGang(card *Mah_Card, _card int, selfget bool) int {
	if selfget {
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
						return key
					}
				} else {
					return key
				}
			}
		}

		return 0
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
							return value
						}
					} else {
						return value
					}
				}
			}
		}

		return 0
	}
}

//! 设置时间
func (self *Game_GoldKWX) SetTime(t int) {
	if t == 0 {
		self.Time = 0
	} else {
		self.Time = time.Now().Unix() + int64(t)
	}

	var msg Msg_SetTime
	msg.Time = lib.HF_MaxInt64(0, self.Time-time.Now().Unix())
	self.room.broadCastMsg("gametime", &msg)
}

//! 结算所有人
func (self *Game_GoldKWX) OnBalance() {
	for i := 0; i < len(self.PersonMgr); i++ {
		//! 退出房间同步金币
		gold := self.PersonMgr[i].Total - self.PersonMgr[i].Gold
		if gold > 0 {
			GetRoomMgr().AddCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, gold, self.room.Type)
		} else if gold < 0 {
			GetRoomMgr().CostCard(self.PersonMgr[i].Uid, staticfunc.TYPE_GOLD, -gold, self.room)
		}
		self.PersonMgr[i].Gold = self.PersonMgr[i].Total
	}
}
