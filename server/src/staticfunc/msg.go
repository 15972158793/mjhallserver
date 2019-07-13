//! 服务器之间的消息
package staticfunc

import (
	"lib"
)

//! 五子棋mode
type WZQMode struct {
	Mode  int     `json:"mode"`
	White []int64 `json:"white"`
}

//! 允许哪些游戏开启
type GameMode struct {
	GoldGame []int `json:"goldgame"`
	RoomGame []int `json:"roomgame"`
}

type JS_CreateRoomMem struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Score int    `json:"score"`
}

type Rpc_Args struct {
	MsgData []byte
}

type Msg_Err struct {
	Err string `json:"info"`
}

//! 战绩
type Msg_Record struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"`
}

//!
type Msg_Null struct {
}

//! 设置管理
type Msg_SetAdmin struct {
	Uid   int64 `json:"uid"`
	Admin int   `json:"admin"`
}

//! gameserver开关
type Msg_GameServer struct {
	Id   int    `json:"id"`
	InIp string `json:"inip"`
	ExIp string `json:"exip"`
	Type int    `json:"type"`
}

//! centerserver开关
type Msg_CenterServer struct {
	InIp string `json:"inip"`
	ExIp string `json:"exip"`
}

//! 改变游戏比例
type Msg_MoneyMode struct {
	MoneyMode int `json:"moneymode"` //! 金币模式 0,1:100;1,1:1;2,1:10000
}

//! 创建房间
type Msg_CreateRoom struct {
	Id     int   `json:"id"`
	Type   int   `json:"type"`
	Num    int   `json:"num"`
	Param1 int   `json:"param1"`
	Param2 int   `json:"param2"`
	Agent  int64 `json:"agent"`
	ClubId int64 `json:"clubid"`
}

//! 创建俱乐部
type Msg_ClubCreateRoom struct {
	ClubId   int64  `json:"clubid"`
	GameType int    `json:"gametype"` //! 游戏type
	Num      int    `json:"num"`      //! 消耗几张房卡
	Param1   int    `json:"param1"`   //! 玩法
	Param2   int    `json:"param2"`
	Host     int64  `json:"host"` //! 公会主席
	Uid      int64  `json:"uid"`  //! 创建者
	IP       string `json:"ip"`
}

//! 加入房间
type Msg_JoinRoom struct {
	Id int `json:"id"`
}

//! 删除房间
type Msg_DelRoom struct {
	Id       int     `json:"id"`
	RoomId   int     `json:"roomid"`
	Uid      []int64 `json:"uid"`
	Host     int64   `json:"host"`
	Agent    bool    `json:"agent"`
	ClubId   int64   `json:"clubid"`
	GameType int     `json:"gametype"`
}

//! 开始房间
type Msg_BeginRoom struct {
	Id     int   `json:"id"`
	RoomId int   `json:"roomid"`
	Host   int64 `json:"host"`
	ClubId int64 `json:"clubid"`
}

//! 观战房间人数变化
type Msg_ViewRoomNum struct {
	Host   int64  `json:"host"`
	RoomId int    `json:"roomid"`
	Num    int    `json:"num"`
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`
	Head   string `json:"head"`
}

//! 更新分数
type Msg_ViewRoomScore struct {
	Host   int64              `json:"host"`
	RoomId int                `json:"roomid"`
	Node   []JS_CreateRoomMem `json:"node"`
}

//! joinfail
type Msg_JoinFail struct {
	Uid      int64 `json:"uid"`
	Id       int   `json:"id"`
	Room     int   `json:"room"`
	GameType int   `json:"gametype"`
}

//! 玩家
type Msg_Uid struct {
	Uid int64 `json:"uid"`
}

//! 玩家
type Msg_RoomId struct {
	roomid int `json:"roomid"`
}

//! 玩家
type Msg_Openid struct {
	Openid string `json:"openid"`
}

//! 消耗
type Msg_CostCard struct {
	Uid  int64 `json:"uid"`
	Num  int   `json:"num"`
	Type int   `json:"type"`
	Dec  int   `json:"dec"`
}

//! 发卡
type Msg_GiveCard struct {
	Uid  int64  `json:"uid"`
	Pid  int    `json:"pid"`
	Num  int    `json:"num"`
	Ip   string `json:"ip"`
	Dec  int    `json:"dec"`
	Sync bool   `json:"sync"` //! 是否需要同步
}

//! 发卡
type Msg_PersonInfo struct {
	Uid       int64      `json:"uid"`
	Nickname  string     `json:"nickname"`
	Headurl   string     `json:"headurl"`
	Repertory []Son_Card `json:"repertory"`
	RoomId    int        `json:"roomid"`
	GameType  int        `json:"gametype"`
}
type Son_Card struct {
	Productid       int `json:"productid"`
	Repertory_count int `json:"repertory_count"`
}

//! 俱乐部战绩
type Msg_ClubRoomResult struct {
	ClubId   int64              `json:"clubid"`
	GameType int                `json:"gametype"`
	Param1   int                `json:"param1"`
	Param2   int                `json:"param2"`
	RoomId   int                `json:"roomid"`
	MaxStep  int                `json:"maxstep"`
	Info     []JS_CreateRoomMem `json:"info"`
	Num      int                `json:"num"`
}

//! 删除俱乐部房间
type Msg_ClubRoomDel struct {
	ClubId int64 `json:"clubid"`
	RoomId int   `json:"roomid"`
}

//! 同步金币
type Msg_SynchroGold struct {
	Uid  int64 `json:"uid"`
	Gold int   `json:"gold"`
}

type Msg_Synchro struct {
	Uid  int64 `json:"uid"`
	Card int   `json:"card"`
	Gold int   `json:"gold"`
}

//! 设置胜率
type Msg_SetDealWin struct {
	GameType int              `json:"gametype"`
	Property lib.ManyProperty `json:"property"`
}

//! 设置筹码
type Msg_SetDealMoney struct {
	GameType int           `json:"gametype"`
	Property lib.ManyMoney `json:"property"`
}

//! 设置next
type Msg_SetDealNext struct {
	RoomId int   `json:"roomid"`
	Next   []int `json:"next"`
}

//! 设置游戏机器人
type Msg_SetGameRobotSet struct {
	GameType int              `json:"gametype"`
	Set      lib.GameRobotSet `json:"set"`
}

//! 发公告
type Msg_RichNotice struct {
	Content string `json:"content"`
}

//! 设置奖池
type Msg_SetJackpot struct {
	GameType int   `json:"gametype"`
	Value    int64 `json:"value"`
}

//!
type Msg_AddRobot struct {
	Id       int    `json:"id"`
	Room     int    `json:"room"`
	Uid      int64  `json:"uid"`
	GameType int    `json:"gametype"`
	Num      int    `json:"num"`
	IP       string `json:"ip"`
}

type Msg_DelRobot struct {
	Id       int   `json:"id"`
	Room     int   `json:"room"`
	Uid      int64 `json:"uid"`
	GameType int   `json:"gametype"`
}

//////////
//! 告诉客户端卡变化
type S2C_UpdCard struct {
	Card int `json:"card"`
	Gold int `json:"gold"`
}

//! 捕鱼桌子
type Msg_FishAddDesk1 struct {
	Desk lib.FishDesk `json:"desk"`
	Lst  []int64      `json:"lst"`
}
type Msg_FishAddDesk2 struct {
	Desk lib.FishDesk `json:"desk"`
}

//! 删除桌子
type Msg_FishDelDesk1 struct {
	Index int     `json:"index"`
	Lst   []int64 `json:"lst"`
}
type Msg_FishDelDesk2 struct {
	Index int `json:"index"`
}

////////////////////////////
//! 记录
//! 斗板凳
type Rec_GameDBD_Info struct {
	Roomid  int                      `json:"roomid"` //! 房间号
	Time    int64                    `json:"time"`   //! 记录时间
	Person  []Son_Rec_GameDBD_Person `json:"person"`
	Step    []Son_Rec_GameDBD_Step   `json:"step"`
	MaxStep int                      `json:"maxstep"`
	Type    int                      `json:"type"`
}
type Son_Rec_GameDBD_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
	Total int    `json:"total"`
}
type Son_Rec_GameDBD_Step struct {
	Uid     int64 `json:"uid"`
	Card    []int `json:"card"`
	AbsCard []int `json:"abscard"`
}

////////////////////////////
//! 记录
//! 斗地主
type Rec_GameDDZ_Info struct {
	Roomid  int                      `json:"roomid"` //! 房间号
	Time    int64                    `json:"time"`   //! 记录时间
	Person  []Son_Rec_GameDDZ_Person `json:"person"`
	Step    []Son_Rec_GameDDZ_Step   `json:"step"`
	Razz    int                      `json:"razz"`
	MaxStep int                      `json:"maxstep"`
	Type    int                      `json:"type"`
}
type Son_Rec_GameDDZ_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
	Total int    `json:"total"`
}
type Son_Rec_GameDDZ_Step struct {
	Uid     int64 `json:"uid"`
	Card    []int `json:"card"`
	AbsCard []int `json:"abscard"`
}

///////////////////////////////
//! 记录
//! 枪炮斗地主
type Rec_GameQPDDZ_Info struct {
	Roomid  int                        `json:"roomid"` //! 房间号
	Time    int64                      `json:"time"`   //! 记录时间
	Person  []Son_Rec_GameQPDDZ_Person `json:"person"`
	Step    []Son_Rec_GameQPDDZ_Step   `json:"step"`
	MaxStep int                        `json:"maxstep"`
	Type    int                        `json:"type"`
}
type Son_Rec_GameQPDDZ_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
	Total int    `json:"total"`
}
type Son_Rec_GameQPDDZ_Step struct {
	Uid     int64 `json:"uid"`
	Card    []int `json:"card"`
	AbsCard []int `json:"abscard"`
}

////////////////////////////
//! 记录
//! 众人斗地主
type Rec_GameZRDDZ_Info struct {
	Roomid  int                        `json:"roomid"` //! 房间号
	Time    int64                      `json:"time"`   //! 记录时间
	Person  []Son_Rec_GameZRDDZ_Person `json:"person"`
	Step    []Son_Rec_GameZRDDZ_Step   `json:"step"`
	Razz    int                        `json:"razz"`
	MaxStep int                        `json:"maxstep"`
	Type    int                        `json:"type"`
}
type Son_Rec_GameZRDDZ_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
	Total int    `json:"total"`
}
type Son_Rec_GameZRDDZ_Step struct {
	Uid     int64 `json:"uid"`
	Card    []int `json:"card"`
	AbsCard []int `json:"abscard"`
}

////////////////////////////
//! 记录
//! 众人扎金花
type Rec_GameZRZJH_Info struct {
	Roomid  int                        `json:"roomid"` //! 房间号
	Time    int64                      `json:"time"`   //! 记录时间
	Person  []Son_Rec_GameZRZJH_Person `json:"person"`
	Step    []Son_Rec_GameZRZJH_Step   `json:"step"`
	MaxStep int                        `json:"maxstep"`
}
type Son_Rec_GameZRZJH_Person struct {
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`
	Head   string `json:"head"`
	Card   []int  `json:"card"`
	Score  int    `json:"score"`
	Total  int    `json:"total"`
	Dealer bool   `json:"dealer"`
}
type Son_Rec_GameZRZJH_Step struct {
	Uid  int64 `json:"uid"`
	Type int64 `json:"type"` //! -1表示弃牌 0表示看牌 1表示跟住 2表示加注 >1000比牌
	Bets int   `json:"bets"`
}

////////////////////////////
//! 记录
//! 慈善大地主
type Rec_GameCSDDZ_Info struct {
	Roomid  int                        `json:"roomid"` //! 房间号
	Time    int64                      `json:"time"`   //! 记录时间
	Person  []Son_Rec_GameCSDDZ_Person `json:"person"`
	Step    []Son_Rec_GameCSDDZ_Step   `json:"step"`
	Razz    int                        `json:"razz"`
	MaxStep int                        `json:"maxstep"`
	Type    int                        `json:"type"`
}
type Son_Rec_GameCSDDZ_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
	Total int    `json:"total"`
}
type Son_Rec_GameCSDDZ_Step struct {
	Uid     int64 `json:"uid"`
	Card    []int `json:"card"`
	AbsCard []int `json:"abscard"`
}

////////////////////////////
//! 记录
//! 跑的快
type Rec_GamePDK_Info struct {
	Roomid  int                      `json:"roomid"` //! 房间号
	Time    int64                    `json:"time"`   //! 记录时间
	MaxStep int                      `json:"maxstep"`
	Type    int                      `json:"type"`
	Person  []Son_Rec_GamePDK_Person `json:"person"`
	Step    []Son_Rec_GamePDK_Step   `json:"step"`
	Razz    int                      `json:"razz"`
}
type Son_Rec_GamePDK_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
	Total int    `json:"total"`
}
type Son_Rec_GamePDK_Step struct {
	Uid     int64 `json:"uid"`
	Card    []int `json:"card"`
	AbsCard []int `json:"abscard"`
}

////////////////////////////
//! 记录
//! 扎金花
type Rec_GameZJH_Info struct {
	Roomid  int                      `json:"roomid"` //! 房间号
	Time    int64                    `json:"time"`   //! 记录时间
	Person  []Son_Rec_GameZJH_Person `json:"person"`
	Step    []Son_Rec_GameJZH_Step   `json:"step"`
	MaxStep int                      `json:"maxstep"`
}
type Son_Rec_GameZJH_Person struct {
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`
	Head   string `json:"head"`
	Card   []int  `json:"card"`
	Score  int    `json:"score"`
	Total  int    `json:"total"`
	Dealer bool   `json:"dealer"`
}
type Son_Rec_GameJZH_Step struct {
	Uid  int64 `json:"uid"`
	Type int64 `json:"type"` //! -1表示弃牌 0表示看牌 1表示跟住 2表示加注 >1000比牌
	Bets int   `json:"bets"`
}

////////////////////////////
//! 记录
//! 逍遥扎金花
type Rec_GameXYZJH_Info struct {
	Roomid  int                        `json:"roomid"` //! 房间号
	Time    int64                      `json:"time"`   //! 记录时间
	MaxStep int                        `json:"maxstep"`
	Person  []Son_Rec_GameXYZJH_Person `json:"person"`
	Step    []Son_Rec_GameXYZJH_Step   `json:"step"`
}
type Son_Rec_GameXYZJH_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Total int    `json:"total"`
	Score int    `json:"score"`
}
type Son_Rec_GameXYZJH_Step struct {
	Uid  int64 `json:"uid"`
	Type int64 `json:"type"` //! -1表示弃牌 0表示看牌 1表示跟住 2表示加注 >1000比牌
}

////////////////////////////
//! 记录
//! 咸宁扎金花
type Rec_GameXNZJH_Info struct {
	Roomid  int                        `json:"roomid"` //! 房间号
	Time    int64                      `json:"time"`   //! 记录时间
	Person  []Son_Rec_GameXNZJH_Person `json:"person"`
	Step    []Son_Rec_GameXNJZH_Step   `json:"step"`
	MaxStep int                        `json:"maxstep"`
}
type Son_Rec_GameXNZJH_Person struct {
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`
	Head   string `json:"head"`
	Card   []int  `json:"card"`
	Score  int    `json:"score"`
	Total  int    `json:"total"`
	Dealer bool   `json:"dealer"`
}
type Son_Rec_GameXNJZH_Step struct {
	Uid  int64 `json:"uid"`
	Type int64 `json:"type"` //! -1表示弃牌 0表示看牌 1表示跟住 2表示加注 >1000比牌
	Bets int   `json:"bets"`
}

////////////////////////////
//! 记录
//! 华众扎金花
type Rec_GameHZZJH_Info struct {
	Roomid  int                        `json:"roomid"` //! 房间号
	Time    int64                      `json:"time"`   //! 记录时间
	Person  []Son_Rec_GameHZZJH_Person `json:"person"`
	Step    []Son_Rec_GameHZZJH_Step   `json:"step"`
	MaxStep int                        `json:"maxstep"`
}
type Son_Rec_GameHZZJH_Person struct {
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`
	Head   string `json:"head"`
	Card   []int  `json:"card"`
	Score  int    `json:"score"`
	Total  int    `json:"total"`
	Dealer bool   `json:"dealer"`
}
type Son_Rec_GameHZZJH_Step struct {
	Uid  int64 `json:"uid"`
	Type int64 `json:"type"` //! -1表示弃牌 0表示看牌 1表示跟住 2表示加注 >1000比牌
	Bets int   `json:"bets"`
}

////////////////////////////
//! 记录
//! 杠次
type Rec_GameGC_Info struct {
	Roomid int                     `json:"roomid"` //! 房间号
	Time   int64                   `json:"time"`   //! 记录时间
	Person []Son_Rec_GameGC_Person `json:"person"`
	Step   []Son_Rec_GameGC_Step   `json:"step"`
}
type Son_Rec_GameGC_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
	Total int    `json:"total"`
}
type Son_Rec_GameGC_Step struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"` //! 0摸牌 1出牌 2碰牌 3次牌 4胡牌 5暗杠 6明杠 7补杠
	Card []int `json:"card"`
}

////////////////////////////
//! 记录
//! 一脚赖油
type Rec_GameYJLY_Info struct {
	Roomid int                       `json:"roomid"` //! 房间号
	Time   int64                     `json:"time"`   //! 记录时间
	Person []Son_Rec_GameYJLY_Person `json:"person"`
	Step   []Son_Rec_GameYJLY_Step   `json:"step"`
}
type Son_Rec_GameYJLY_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
	Total int    `json:"total"`
}
type Son_Rec_GameYJLY_Step struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"` //! 0摸牌 1出牌 2碰牌  3胡牌 4暗杠 5明杠 6补杠 7赖子杠
	Card []int `json:"card"`
}

////////////////////////////
//! 记录
//! 仙桃人人晃晃
type Rec_GameXTRRHH_Info struct {
	Roomid  int                         `json:"roomid"` //! 房间号
	Time    int64                       `json:"time"`   //! 记录时间
	Person  []Son_Rec_GameXTRRHH_Person `json:"person"`
	Step    []Son_Rec_GameXTRRHH_Step   `json:"step"`
	Param1  int                         `json:"param1"`
	Param2  int                         `json:"param2"`
	MaxStep int                         `json:"maxstep"`
}
type Son_Rec_GameXTRRHH_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
	Total int    `json:"total"`
}
type Son_Rec_GameXTRRHH_Step struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"` //! 0摸牌 1出牌 2碰牌  3胡牌 4暗杠 5明杠 6补杠 7赖子杠
	Card []int `json:"card"`
}

////////////////////////////
//! 记录
//! 宿州麻将
type Rec_GameSZMJ_Info struct {
	Roomid int                       `json:"roomid"` //! 房间号
	Time   int64                     `json:"time"`   //! 记录时间
	Person []Son_Rec_GameSZMJ_Person `json:"person"`
	Step   []Son_Rec_GameSZMJ_Step   `json:"step"`
}
type Son_Rec_GameSZMJ_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
}
type Son_Rec_GameSZMJ_Step struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"` //! 0摸牌 1出牌 2碰牌 3胡牌 4暗杠 5明杠 6补杠 7摊花 8补牌
	Card []int `json:"card"`
}

////////////////////////////
//! 记录
//! 安庆麻将
type Rec_GameAQMJ_Info struct {
	Roomid  int                       `json:"roomid"` //! 房间号
	Time    int64                     `json:"time"`   //! 记录时间
	Person  []Son_Rec_GameAQMJ_Person `json:"person"`
	Step    []Son_Rec_GameAQMJ_Step   `json:"step"`
	Param1  int                       `json:"param1"`
	Param2  int                       `json:"param2"`
	MaxStep int                       `json:"maxstep"`
}
type Son_Rec_GameAQMJ_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
	Total int    `json:"total"`
}
type Son_Rec_GameAQMJ_Step struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"` //! 0摸牌 1出牌 2碰牌 3胡牌 4暗杠 5明杠 6补杠 7摊花 8补牌 9吃牌
	Card []int `json:"card"`
}

////////////////////////////
//! 记录
//! 涡阳麻将
type Rec_GameGYMJ_Info struct {
	Roomid  int                       `json:"roomid"` //! 房间号
	Time    int64                     `json:"time"`   //! 记录时间
	Person  []Son_Rec_GameGYMJ_Person `json:"person"`
	Step    []Son_Rec_GameGYMJ_Step   `json:"step"`
	Param1  int                       `json:"param1"`
	Param2  int                       `json:"param2"`
	MaxStep int                       `json:"maxstep"`
}
type Son_Rec_GameGYMJ_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
	Total int    `json:"total"`
}
type Son_Rec_GameGYMJ_Step struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"` //! 0摸牌 1出牌 2碰牌 3胡牌 4暗杠 5明杠 6补杠 7摊花 8补牌 9吃牌
	Card []int `json:"card"`
}

////////////////////////////
//! 记录
//! 常熟麻将
type Rec_GameCSMJ_Info struct {
	Roomid int                       `json:"roomid"` //! 房间号
	Time   int64                     `json:"time"`   //! 记录时间
	Person []Son_Rec_GameCSMJ_Person `json:"person"`
	Step   []Son_Rec_GameCSMJ_Step   `json:"step"`
	Param1 int                       `json:"param1"`
	Param2 int                       `json:"param2"`
}
type Son_Rec_GameCSMJ_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Piao  int    `json:"piao"`
	Score int    `json:"score"`
}
type Son_Rec_GameCSMJ_Step struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"` //! 0摸牌 1出牌 2碰牌 3胡牌 4暗杠 5明杠 6补杠 7摊花 8补牌 9吃牌
	Card []int `json:"card"`
}

////////////////////////////
//! 记录
//! 通山麻将
type Rec_GameTSMJ_Info struct {
	Roomid  int                       `json:"roomid"` //! 房间号
	Time    int64                     `json:"time"`   //! 记录时间
	Person  []Son_Rec_GameTSMJ_Person `json:"person"`
	Step    []Son_Rec_GameTSMJ_Step   `json:"step"`
	Param1  int                       `json:"param1"`
	Param2  int                       `json:"param2"`
	MaxStep int                       `json:"maxstep"`
}
type Son_Rec_GameTSMJ_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Piao  int    `json:"piao"`
	Score int    `json:"score"`
}
type Son_Rec_GameTSMJ_Step struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"` //! 0摸牌 1出牌 2碰牌 3胡牌 4暗杠 5明杠 6补杠 7摊花 8补牌 9吃牌 10赖子
	Card []int `json:"card"`
}

//! 记录
//! 推倒胡
type Rec_GameTDH_Info struct {
	Roomid int                      `json:"roomid"` //! 房间号
	Time   int64                    `json:"time"`   //! 记录时间
	Person []Son_Rec_GameTDH_Person `json:"person"`
	Step   []Son_Rec_GameTDH_Step   `json:"step"`
	Param1 int                      `json:"param1"`
	Param2 int                      `json:"param2"`
}
type Son_Rec_GameTDH_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Piao  int    `json:"piao"`
	Score int    `json:"score"`
}
type Son_Rec_GameTDH_Step struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"` //! 0摸牌 1出牌 2碰牌 3胡牌 4暗杠 5明杠 6补杠 7摊花 8补牌 9吃牌
	Card []int `json:"card"`
}

////////////////////////////
//! 记录
//! 潮汕麻将
type Rec_GameGDCSMJ_Info struct {
	Roomid  int                         `json:"roomid"` //! 房间号
	Time    int64                       `json:"time"`   //! 记录时间
	Person  []Son_Rec_GameGDCSMJ_Person `json:"person"`
	Step    []Son_Rec_GameGDCSMJ_Step   `json:"step"`
	Param1  int                         `json:"param1"`
	Param2  int                         `json:"param2"`
	MaxStep int                         `json:"maxstep"`
}
type Son_Rec_GameGDCSMJ_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Piao  int    `json:"piao"`
	Score int    `json:"score"`
}
type Son_Rec_GameGDCSMJ_Step struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"` //! 0摸牌 1出牌 2碰牌 3胡牌 4暗杠 5明杠 6补杠 7摊花 8补牌 9吃牌
	Card []int `json:"card"`
}

////////////////////////////
//! 记录
//! 贵溪麻将
type Rec_GameGXMJ_Info struct {
	Roomid int                       `json:"roomid"` //! 房间号
	Time   int64                     `json:"time"`   //! 记录时间
	Person []Son_Rec_GameGXMJ_Person `json:"person"`
	Step   []Son_Rec_GameGXMJ_Step   `json:"step"`
	Param1 int                       `json:"param1"`
	Param2 int                       `json:"param2"`
}
type Son_Rec_GameGXMJ_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Piao  int    `json:"piao"`
	Score int    `json:"score"`
}
type Son_Rec_GameGXMJ_Step struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"` //! 0摸牌 1出牌 2碰牌 3胡牌 4暗杠 5明杠 6补杠 7摊花 8补牌 9吃牌
	Card []int `json:"card"`
}

////////////////////////////
//! 记录
//! 南昌麻将
type Rec_GameNCMJ_Info struct {
	Roomid int                       `json:"roomid"` //! 房间号
	Time   int64                     `json:"time"`   //! 记录时间
	Person []Son_Rec_GameNCMJ_Person `json:"person"`
	Step   []Son_Rec_GameNCMJ_Step   `json:"step"`
	Param1 int                       `json:"param1"`
	Param2 int                       `json:"param2"`
}
type Son_Rec_GameNCMJ_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
}
type Son_Rec_GameNCMJ_Step struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"` //! 0摸牌 1出牌 2碰牌 3胡牌 4暗杠 5明杠 6补杠 7摊花 8补牌 9吃牌
	Card []int `json:"card"`
}

////////////////////////////
//! 记录
//! 上虞花麻将
type Rec_GameSYHMJ_Info struct {
	Roomid  int                        `json:"roomid"` //! 房间号
	Time    int64                      `json:"time"`   //! 记录时间
	Person  []Son_Rec_GameSYHMJ_Person `json:"person"`
	Step    []Son_Rec_GameSYHMJ_Step   `json:"step"`
	Param1  int                        `json:"param1"`
	Param2  int                        `json:"param2"`
	MaxStep int                        `json:"maxstep"`
}
type Son_Rec_GameSYHMJ_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
}
type Son_Rec_GameSYHMJ_Step struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"` //! 0摸牌 1出牌 2碰牌 3胡牌 4暗杠 5明杠 6补杠 7摊花 8补牌 9吃牌
	Card []int `json:"card"`
}

////////////////////////////
//! 记录
//! 杭州麻将
type Rec_GameHZMJ_Info struct {
	Roomid int                       `json:"roomid"` //! 房间号
	Time   int64                     `json:"time"`   //! 记录时间
	Person []Son_Rec_GameHZMJ_Person `json:"person"`
	Step   []Son_Rec_GameHZMJ_Step   `json:"step"`
	Param1 int                       `json:"param1"`
	Param2 int                       `json:"param2"`
}
type Son_Rec_GameHZMJ_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
}
type Son_Rec_GameHZMJ_Step struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"` //! 0摸牌 1出牌 2碰牌 3胡牌 4暗杠 5明杠 6补杠 7摊花 8补牌 9吃牌
	Card []int `json:"card"`
}

////////////////////////////
//! 记录
//! 血战麻将
type Rec_GameXZDD_Info struct {
	Roomid int                       `json:"roomid"` //! 房间号
	Time   int64                     `json:"time"`   //! 记录时间
	Param1 int                       `json:"param1"`
	Param2 int                       `json:"param2"`
	Person []Son_Rec_GameXZDD_Person `json:"person"`
	Step   []Son_Rec_GameXZDD_Step   `json:"step"`
}
type Son_Rec_GameXZDD_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Que   int    `json:"que"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
}
type Son_Rec_GameXZDD_Step struct {
	Uid  int64 `json:"uid"`
	Type int   `json:"type"` //! 0摸牌 1出牌 2碰牌 3暗杠 4明杠 5补杠 6自摸胡 7点炮胡
	Card []int `json:"card"`
}

////////////////////////////
//! 金币场记录
type Rec_Gold_Info struct {
	GameType int                   `json:"gametype"`
	Time     int64                 `json:"time"` //! 记录时间
	Info     []Son_Rec_Gold_Person `json:"info"`
}
type Son_Rec_Gold_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Score int    `json:"score"`
	Robot bool   `json:"robot"`
}
