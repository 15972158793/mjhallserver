package gameserver

import (
	"lib"
)

const GAMETYPE_NIUNIU = 1
const GAMETYPE_KWX_XG = 2   //! 孝感
const GAMETYPE_KWX_XY = 3   //! 襄阳
const GAMETYPE_KWX_SY = 4   //! 十堰
const GAMETYPE_KWX_SZ = 5   //! 随州
const GAMETYPE_KWX_YIC = 13 //! 宜城
const GAMETYPE_KWX_YC = 14  //! 应城

//! 解散类型
type Msg_DissRoom struct {
	Type int `json:"type"`
}

//! 游戏下注
type Msg_GameBets struct {
	Uid  int64 `json:"uid"`
	Bets int   `json:"bets"`
}

//! 游戏下注
type Msg_GameBJLBets struct {
	Uid    int64 `json:"uid"`
	Bets   int   `json:"bets"`
	Type   int   `json:"type"`
	Amount int   `json:"amount"`
}

//! 游戏下注
type Msg_GameTDZBets struct {
	Uid  int64 `json:"uid"`
	Bets int   `json:"bets"`
	Type int   `json:"type"`
}

//! 游戏换三张
type Msg_GameThree struct {
	Uid   int64 `json:"uid"`
	Three []int `json:"three"`
}

//! 游戏飘
type Msg_GamePiao struct {
	Gold int `json:"gold"`
}

//! 游戏出牌
type Msg_GameStep struct {
	Card int `json:"card"`
}

//! 游戏吃牌
type Msg_GameChi struct {
	ChiCard []int `json:"chicard"`
}

type Msg_GameChangeCard struct {
	Card   int `json:"card"`
	ChCard int `json:"chcard"`
}

type Msg_GameChange struct {
	Card []int `json:"card"`
}

type Msg_GameStart struct {
	LineNum int `json:"linenum"`
	Gold    int `josn:"gold"`
}

//! 游戏出牌
type Msg_GameSteps struct {
	Cards    []int `json:"cards"`
	AbsCards []int `json:"abscards"`
}

//! 游戏跑得快出牌
type Msg_GameStepsPDK struct {
	Id       int64 `json:"id"`
	Cards    []int `json:"cards"`
	AbsCards []int `json:"abscards"`
	TypeCard int   `json:"typecard"`
}

//! 十三道配牌
type Msg_GameSSDMatch struct {
	Cards    []int `json:"cards"`
	AbsCards []int `json:"abscards"`
	Special  int   `json:"special"`
}

//! 游戏抢庄
type Msg_GameDeal struct {
	Uid int64 `json:"uid"`
	Ok  bool  `json:"ok"`
}

//! 游戏抢庄
type Msg_GameDealer struct {
	Uid   int64 `json:"uid"`
	Score int   `json:"score"`
}

//! 游戏比牌
type Msg_GameCompare struct {
	Uid     int64 `json:"uid"`
	Destuid int64 `json:"destuid"`
}

//! 游戏操作
type Msg_GamePlay struct {
	Type int `json:"type"`
}

//! 配牌操作
type Msg_GameMatch struct {
	Cards   []int `json:"cards"`
	Special int   `json:"special"`
}

//! 游戏抢庄
type Msg_GameLzDealer struct {
	//	Uid  int64 `json:"uid"`
	Type int `json:"type"`
}

//! 游戏扫雷
type Msg_GameBoom struct {
	//	Uid int64 `json:"uid"`
	Pos int `json:"pos"`
}

//! 亮牌
type Msg_GameView struct {
	Uid  int64 `json:"uid"`
	Card []int `json:"card"`
}

//! 加倍
type Msg_GameDouble struct {
	Uid    int64 `json:"uid"`
	Double int   `json:"isdouble"`
}

//! 时间
type Msg_SetTime struct {
	Time int64 `json:"time"`
}

//! 撤销订单
type Msg_CancleOrder struct {
	OrderId string `json:"orderid"`
}

//！电玩奖池
type Msg_DwSetJackPot struct {
	JackPot    int64 `json:"jackpot"`
	JackPotMax int64 `json:"jackpotmax"`
	JackPotMin int64 `json:"jackpotmin"`
}

//！电玩属性
type Msg_DwSetPro struct {
	GameLevel   int `json:"gamelevel"`
	PersonalSet int `json:"personalset"`
}

//！电玩房间详细
type Msg_DwSetRoom struct {
	RoomId      int `json:"roomid"`
	GameLevel   int `json:"gamelevel"`
	PersonalSet int `json:"personalset"`
	FreeRun     int `json:"freerun"`
}

type Game_Base interface {
	OnBegin() //! 一局开始
	OnEnd()   //! 一局结束

	OnInit(room *Room)         //! 初始化
	OnSendInfo(person *Person) //! 告诉玩家数据
	OnRobot(robot *lib.Robot)  //! 机器人加入
	OnMsg(msg *RoomMsg)        //! 消息转发
	OnBye()                    //! 游戏结算
	OnExit(uid int64)          //! 玩家退出
	OnIsDealer(uid int64) bool //! 是否是庄家
	OnIsBets(uid int64) bool   //! 是否下注
	OnBalance()                //! 结算

	OnTime() //! 每秒调用一次
}
