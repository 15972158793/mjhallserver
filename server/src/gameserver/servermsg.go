package gameserver

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"lib"
	"net/http"
	"runtime/debug"
	"staticfunc"
)

type ServerMethod int

func (self *ServerMethod) ServerMsg(args *staticfunc.Rpc_Args, reply *[]byte) error {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	if args == nil || reply == nil {
		return errors.New("nil paramters !")
	}

	head, _, data, ok := lib.HF_DecodeMsg(args.MsgData)
	if !ok {
		return errors.New("args err !")
	}

	var msg staticfunc.Msg_Record
	json.Unmarshal(data, &msg)

	switch head {
	case "loginserver": //! login服务器断开
		go GetServer().ConnectLogin()
	case "createroom": //! 创建一个房间
		var msg staticfunc.Msg_CreateRoom
		json.Unmarshal(data, &msg)
		if GetRoomMgr().CreateRoom(msg.Id, msg.Type, msg.Num, msg.Param1, msg.Param2, msg.Agent, msg.ClubId) {
			*reply = []byte("true")
		} else {
			*reply = []byte("false")
		}
	case "joinroom": //! 加入房间
		var msg staticfunc.Msg_JoinRoom
		json.Unmarshal(data, &msg)
		ok, room := GetRoomMgr().JoinRoom(msg.Id)
		if ok {
			*reply = []byte(fmt.Sprintf("%d", len(room.Uid)+len(room.Viewer)))
		} else {
			*reply = []byte("-1")
		}
	case "synchrogold":
		var msg staticfunc.Msg_SynchroGold
		json.Unmarshal(data, &msg)
		person := GetPersonMgr().GetPerson(msg.Uid)
		if person == nil {
			return nil
		}
		person.Gold = msg.Gold
		if person.room == nil {
			return nil
		}
		person.room.Operator(NewRoomMsg("synchrogold", "", msg.Uid, &msg))
	case "addblack": //! 加入黑名单
		var msg staticfunc.Msg_Uid
		json.Unmarshal(data, &msg)
		person := GetPersonMgr().GetPerson(msg.Uid)
		if person == nil {
			return nil
		}
		if person.room == nil {
			return nil
		}
		person.black = true
		if person.room.Type == 77 { //! 五子棋房间直接解散
			person.room.Operator(NewRoomMsg("byeroom", "", 0, nil))
		}
	case "dissmiss": //! 强制解散房间
		var msg staticfunc.Msg_JoinRoom
		json.Unmarshal(data, &msg)
		room := GetRoomMgr().GetRoom(msg.Id)
		if room == nil {
			return nil
		}
		room.Operator(NewRoomMsg("byeroom", "", 0, nil))
	case "setdealwin":
		var msg staticfunc.Msg_SetDealWin
		json.Unmarshal(data, &msg)
		lib.GetManyMgr().SetProperty(msg.GameType, msg.Property)
	case "setdealmoney":
		var msg staticfunc.Msg_SetDealMoney
		json.Unmarshal(data, &msg)
		lib.GetManyMoneyMgr().SetProperty(msg.GameType, msg.Property)
	case "setdealnext":
		var msg staticfunc.Msg_SetDealNext
		json.Unmarshal(data, &msg)
		room := GetRoomMgr().GetRoom(msg.RoomId)
		if room == nil {
			return nil
		}
		room.Operator(NewRoomMsg("gamesetnext", "", 0, &msg))
	case "setadmin":
		var msg staticfunc.Msg_SetAdmin
		json.Unmarshal(data, &msg)
		person := GetPersonMgr().GetPerson(msg.Uid)
		if person != nil {
			person.Admin = msg.Admin
		}
	case "loadrobot":
		var msg staticfunc.Msg_Uid
		json.Unmarshal(data, &msg)
		lib.GetRobotMgr().LoadFromID(int(msg.Uid))
	case "delrobot":
		var msg staticfunc.Msg_Uid
		json.Unmarshal(data, &msg)
		lib.GetRobotMgr().DelRobot(msg.Uid)
	case "setgamerobotset":
		var msg staticfunc.Msg_SetGameRobotSet
		json.Unmarshal(data, &msg)
		lib.GetRobotMgr().SetRobotSet(msg.GameType, msg.Set)
	case "setjackpot":
		var msg staticfunc.Msg_SetJackpot
		json.Unmarshal(data, &msg)
		if msg.GameType/10000 == 4 { //! 豹子王4000
			GetServer().SetSystemMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 12 { //! 摇色子120000
			GetServer().SetYSZSYSMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 6 { //! 百人推筒子60000
			GetServer().SetBrTTZMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 9 { //! 神仙夺宝90000
			GetServer().SetSxdbSysMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 20 { //! 名品汇200000
			GetServer().SetMphSysMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 10 { //! 龙虎斗100001
			GetServer().SetLHDMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 14 { //! 单双140000
			GetServer().SetTBMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 16 { //! 龙珠夺宝160000
			GetServer().SetLzdbSysMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 17 { //! 地穴探宝170000
			GetServer().SetDfwSysMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 13 { //! 赛马130001
			GetServer().SetSaiMaMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 21 { //! 红黑大战210000
			GetServer().SetHhdzMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 19 { //! 翻牌机190000-190002
			GetServer().SetFKFpjSysMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 23 { //! 百人牛牛230000
			GetServer().SetBrNNMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 24 { //! 鱼虾蟹240000
			GetServer().SetYxxSysMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 25 { //! 3D捕鱼250000
			GetServer().SetFishMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 26 { //! 百家乐260000
			GetServer().SetBJLSysMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 27 {
			GetServer().SetSHZSysMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 28 {
			GetServer().SetLkpyMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 31 {
			GetServer().SetDwDyjlbSysMoney(msg.GameType%10000, msg.Value)
		} else if msg.GameType/10000 == 32 {
			GetServer().SetDwJdqsSysMoney(msg.GameType%10000, msg.Value)
		}
	case "setrobotjackpot":
		var msg staticfunc.Msg_SetJackpot
		json.Unmarshal(data, &msg)
		lib.GetRobotMgr().SetRobotWin(msg.GameType, int(msg.Value))
	case "moneymode":
		var msg staticfunc.Msg_MoneyMode
		json.Unmarshal(data, &msg)
		GetServer().Con.MoneyMode = msg.MoneyMode
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "moneymode:", msg.MoneyMode)
	}

	return nil
}

func GMMsg(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	msghead := req.FormValue("msghead")
	if msghead == "initcard" {
		roomid := req.FormValue("roomid")
		room := GetRoomMgr().GetRoom(lib.HF_Atoi(roomid))
		if room == nil {
			w.Write([]byte("房间错误"))
			return
		}
		data := req.FormValue("data")
		text, err := base64.URLEncoding.DecodeString(data)
		if err != nil {
			w.Write([]byte("参数错误1"))
			return
		}
		var msg Msg_GameKWX_Need
		err = json.Unmarshal([]byte(text), &msg)
		if err != nil {
			w.Write([]byte("参数错误2"))
			return
		}
		room.Operator(NewRoomMsg("gameinitcard", "", 0, &msg))
		w.Write([]byte("ok"))
	}
}
