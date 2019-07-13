package centerserver

/////////////////////////////
type Msg_SetUid struct {
	Uid     int64  `json:"uid"`
	Openid  string `json:"openid"`
	MapInfo string `json:"mapinfo"`
	Group   int    `json:"group"`
}

type Msg_SetMapInfo struct {
	MapInfo string `json:"mapinfo"`
}

type Msg_ReadMail struct {
	Id int `json:"id"`
}

type Msg_FriendUid struct {
	Uid int64 `json:"uid"`
}

//! 举报
type Msg_Report struct {
	Uid  int64  `json:"uid"`
	Type int    `json:"type"`
	Dec  string `json:"dec"`
}

type S2C_Notice struct {
	Context string `json:"context"`
}
type S2C_GiftAmount struct {
	Id     int64 `json:"id"`
	Amount int   `json:"amount"`
}
type C2S_GiftGive struct {
	Uid    int64 `json:"uid"`
	Type   int   `json:"type"`
	Id     int64 `json:"id"`
	Amount int   `json:"amount"`
}
type C2S_GiftBuy struct {
	Type   int   `json:"type"`
	Id     int64 `json:"id"`
	Amount int   `json:"amount"`
}
type C2S_Clothes struct {
	Id int `json:"id"`
}
type C2S_ClothesFind struct {
	Id int64 `json:"id"`
}
type C2S_ClothesUse struct {
	CurId int `json:"curid"`
	UseId int `json:"useid"`
}
type S2C_ClothesState struct {
	Id    int `json:"id"`
	State int `json:"state"`
}

type C2S_TaskId struct {
	Id int `json:"id"`
}

type C2S_TaskTimes struct {
	Id    int `json:"id"`
	Times int `json:"times"`
}

type S2C_TaskTimes struct {
	Uid   int64 `json:"uid"`
	Times int   `json:"times"`
}

type S2C_TaskState struct {
	Uid   int64 `json:"uid"`
	State bool  `json:"state"`
}

type S2C_PacksState struct {
	Uid   int64 `json:"uid"`
	State int   `json:"state"`
}

type C2S_Charge struct {
	Receipt string `json:"receipt"`
	Sandbox bool   `json:"sandbox"`
}

type C2S_Code struct {
	Code string `json:"code"`
}

type C2S_InviteCode struct {
	Code int `json:"code"`
}

type S2C_Err struct {
	Info string `json:"info"`
}

/////////////////////////////////
//! c_s
type Msg_ClubID struct {
	ClubId int64 `json:"clubid"`
}

//! 创建俱乐部
type C2S_ClubCreate struct {
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Notice string `json:"notice"`
	Mode   int    `json:"mode"`
}

//! 处理俱乐部
type C2S_ClubOrder struct {
	ClubId int64 `json:"clubid"`
	Uid    int64 `json:"uid"`
	Agree  bool  `json:"agree"`
}

//! 俱乐部创建房间
type C2S_ClubCreateRoom struct {
	ClubId   int64 `json:"clubid"`
	GameType int   `json:"gametype"`
	Param1   int   `json:"param1"`
	Param2   int   `json:"param2"`
	Num      int   `json:"num"`
}

//! 俱乐部开房模式
type C2S_ClubMode struct {
	ClubId int64 `json:"clubid"`
	Mode   int   `json:"mode"`
}

//! 设置俱乐部游戏
type C2S_ClubGame struct {
	ClubId int64         `json:"clubid"`
	Game   []JS_ClubGame `json:"game"`
}

//! 设置俱乐部名字
type C2S_ClubName struct {
	ClubId int64  `json:"clubid"`
	Name   string `json:"name"`
}

//! 离开俱乐部
type C2S_ClubLeave struct {
	ClubId int64 `json:"clubid"`
	Uid    int64 `json:"uid"`
}

//! 修改公告
type C2S_ClubNotice struct {
	ClubId int64  `json:"clubid"`
	Notice string `json:"notice"`
}

//! 修改头像
type C2S_ClubIcon struct {
	ClubId int64  `json:"clubid"`
	Icon   string `json:"icon"`
}

//! s_c
//! 自己俱乐部情况
type Msg_MyClubInfo struct {
	Info  []Son_ClubList `json:"info"`
	Apply []int64        `json:"apply"`
}

//! 俱乐部列表
type Son_ClubList struct {
	Id     int64         `json:"id"`     //! 俱乐部id
	Name   string        `json:"name"`   //! 俱乐部名字
	Icon   string        `json:"icon"`   //! 俱乐部icon
	Notice string        `json:"notice"` //! 俱乐部公告
	Num    int           `json:"num"`    //! 俱乐部人数
	Game   []JS_ClubGame `json:"game"`   //! 俱乐部游戏
}

type Msg_ClubList struct {
	Info []Son_ClubList `json:"info"` //! 俱乐部列表
}

//! 俱乐部信息
type Msg_ClubInfo struct {
	Id       int64         `json:"id"`       //! 俱乐部id
	Name     string        `json:"name"`     //! 俱乐部名字
	Icon     string        `json:"icon"`     //! 俱乐部icon
	Host     int64         `json:"host"`     //! 俱乐部主席
	Mode     int           `json:"mode"`     //! 开房模式
	InNotice string        `json:"innotice"` //! 俱乐部对内公告
	ExNotice string        `json:"exnotice"` //! 俱乐部对外公告
	Member   []JS_ClubMem  `json:"member"`   //! 俱乐部成员
	Game     []JS_ClubGame `json:"game"`     //! 游戏
	Red      int           `json:"red"`
	MsgRed   int           `json:"msgred"`
}

//! 刷新成员信息
type Msg_ClubMem struct {
	Id     int64          `json:"id"`     //! 俱乐部id
	Member []JS_ClubMem   `json:"member"` //! 俱乐部成员
	Apply  []JS_ClubApply `json:"apply"`  //! 俱乐部申请列表
}

//! 发送俱乐部聊天
type Msg_ClubAddRoomChat struct {
	Info JS_ClubRoomChat `json:"info"`
}

//! 发送俱乐部战绩
type Msg_ClubAddRoomResult struct {
	Info JS_ClubRoomResult `json:"info"`
}

//! 俱乐部房间列表
type Msg_ClubRoomChat struct {
	Info []JS_ClubRoomChat `json:"info"`
}

//! 俱乐部房间列表
type Msg_ClubRoomResult struct {
	Info []JS_ClubRoomResult `json:"info"`
}

//! 俱乐部开房模式
type S2C_ClubMode struct {
	Id   int64 `json:"id"`
	Mode int   `json:"mode"`
}

//! 设置俱乐部游戏
type S2C_ClubGame struct {
	Id   int64         `json:"id"`
	Game []JS_ClubGame `json:"game"`
}

//! 设置俱乐部名字
type S2C_ClubName struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

//! 离开俱乐部
type S2C_ClubLeave struct {
	Id  int64 `json:"id"`
	Uid int64 `json:"uid"`
}

//! 得到申请列表
type S2C_ClubApplyList struct {
	Id   int64          `json:"id"`
	Info []JS_ClubApply `json:"info"`
}

//! 得到事件列表
type S2C_ClubEventList struct {
	Id   int64          `json:"id"`
	Info []JS_ClubEvent `json:"info"`
}

//! 得到房卡统计
type S2C_ClubCostCard struct {
	ClubId int64             `json:"clubid"`
	Info   []JS_ClubCostCard `json:"info"`
}
