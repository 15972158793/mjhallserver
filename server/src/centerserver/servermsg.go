package centerserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"lib"
	"runtime/debug"
	"staticfunc"
	"time"
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
	case "getuserbase": //! 得到person信息
		var msg staticfunc.Msg_Uid
		json.Unmarshal(data, &msg)
		c := GetServer().Redis.Get()
		defer c.Close()
		value, _ := redis.Bytes(c.Do("GET", fmt.Sprintf("userbase_%d", msg.Uid)))
		*reply = value
	case "clubroomresult": //! 俱乐部战绩
		var msg staticfunc.Msg_ClubRoomResult
		json.Unmarshal(data, &msg)
		club := GetClubMgr().GetClub(msg.ClubId)
		if club == nil {
			return nil
		}

		var node JS_ClubRoomResult
		node.RoomId = msg.RoomId
		node.GameType = msg.GameType
		node.Param1 = msg.Param1
		node.Param2 = msg.Param2
		node.MaxStep = msg.MaxStep
		node.Info = msg.Info
		node.Num = msg.Num
		node.Time = time.Now().Unix()
		club.AddRoomResult(node)
	case "clubroomdel": //! 解散了俱乐部房间
		var msg staticfunc.Msg_ClubRoomDel
		json.Unmarshal(data, &msg)
		club := GetClubMgr().GetClub(msg.ClubId)
		if club == nil {
			return nil
		}
		club.DelRoomChat(msg.RoomId)
		club.BroadCastMsg(lib.HF_EncodeMsg("delroomchat", &msg, true), false)
	case "clubroombegin": //! 俱乐部房间开始游戏
		var msg staticfunc.Msg_ClubRoomDel
		json.Unmarshal(data, &msg)
		club := GetClubMgr().GetClub(msg.ClubId)
		if club == nil {
			return nil
		}
		club.StateRoomChat(msg.RoomId)
		club.BroadCastMsg(lib.HF_EncodeMsg("stateroomchat", &msg, true), false)
	case "synchro": //! 同步
		var msg staticfunc.Msg_Synchro
		json.Unmarshal(data, &msg)
		person := GetPersonMgr().GetPerson(msg.Uid, false)
		if person != nil {
			person.UpdCard(msg.Card, msg.Gold)
		}
	case "fishadddesk":
		var msg staticfunc.Msg_FishAddDesk1
		json.Unmarshal(data, &msg)
		var _msg staticfunc.Msg_FishAddDesk2
		_msg.Desk = msg.Desk
		data := lib.HF_EncodeMsg("fishadddesk", &_msg, true)
		for i := 0; i < len(msg.Lst); i++ {
			person := GetPersonMgr().GetPerson(msg.Lst[i], false)
			if person == nil {
				continue
			}
			person.SendMsg(data)
		}
	case "fishaddperson": //!
		var msg staticfunc.Msg_FishAddDesk1
		json.Unmarshal(data, &msg)
		var _msg staticfunc.Msg_FishAddDesk2
		_msg.Desk = msg.Desk
		data := lib.HF_EncodeMsg("fishaddperson", &_msg, true)
		for i := 0; i < len(msg.Lst); i++ {
			person := GetPersonMgr().GetPerson(msg.Lst[i], false)
			if person == nil {
				continue
			}
			person.SendMsg(data)
		}
	case "fishdelperson":
		var msg staticfunc.Msg_FishAddDesk1
		json.Unmarshal(data, &msg)
		var _msg staticfunc.Msg_FishAddDesk2
		_msg.Desk = msg.Desk
		data := lib.HF_EncodeMsg("fishaddperson", &_msg, true)
		for i := 0; i < len(msg.Lst); i++ {
			person := GetPersonMgr().GetPerson(msg.Lst[i], false)
			if person == nil {
				continue
			}
			person.SendMsg(data)
		}
	case "fishdeldesk":
		var msg staticfunc.Msg_FishDelDesk1
		json.Unmarshal(data, &msg)
		var _msg staticfunc.Msg_FishDelDesk2
		_msg.Index = msg.Index
		data := lib.HF_EncodeMsg("fishdeldesk", &_msg, true)
		for i := 0; i < len(msg.Lst); i++ {
			person := GetPersonMgr().GetPerson(msg.Lst[i], false)
			if person == nil {
				continue
			}
			person.SendMsg(data)
		}
	case "richnotice":
		var msg staticfunc.Msg_RichNotice
		json.Unmarshal(data, &msg)
		GetPersonMgr().BroadCastMsg2(lib.HF_EncodeMsg("richnotice", &msg, true))
	case "moneymode":
		var msg staticfunc.Msg_MoneyMode
		json.Unmarshal(data, &msg)
		GetServer().Con.MoneyMode = msg.MoneyMode
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "moneymode:", msg.MoneyMode)
	}

	return nil
}
