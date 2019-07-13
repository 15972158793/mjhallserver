//! 日志
package staticfunc

import ()

const TYPE_CARD = 1 //! 银卡
const TYPE_GOLD = 2 //! 金卡

const ADD_CARD = 11  //! 后台加银卡
const ADD_GOLD = 12  //! 后台加金卡
const COST_CARD = 21 //! 游戏消耗银卡
const COST_GOLD = 22 //! 游戏消耗金卡
const BUY_CARD = 31  //! 购买银卡
const BUY_GOLD = 32  //! 购买金卡
const MOVE_CARD = 41 //! 后台减银卡
const MOVE_GOLD = 42 //! 后台减金卡

type Log_Base struct {
	Id            int64
	Uid           int64
	Type          int
	Num           int
	Info          string
	Creation_time int64
}

type Log_Room struct {
	Roomid int     `json:"roomid"`
	Uid    []int64 `json:"uid"`
}
