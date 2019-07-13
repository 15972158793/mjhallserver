package loginserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"lib"
	//"rjmgr"
	"runtime/debug"
	"staticfunc"
	"strings"
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
	case "getperson": //! 得到person信息
		var msg staticfunc.Msg_Uid
		json.Unmarshal(data, &msg)
		value, _ := GetServer().DB_GetData("user", msg.Uid, 1)
		*reply = value
	case "costcard": //! 扣除房卡
		var msg staticfunc.Msg_CostCard
		json.Unmarshal(data, &msg)
		value, _ := GetServer().DB_GetData("user", msg.Uid, 1)
		if string(value) != "" {
			var person Person
			json.Unmarshal(value, &person)
			card := person.Card
			gold := person.Gold
			if msg.Type == staticfunc.TYPE_GOLD { //!  优先扣除金卡
				person.Gold = lib.HF_MaxInt(0, person.Gold-msg.Num)
				if msg.Dec != 77 {
					person.BindGold = lib.HF_MaxInt(0, person.BindGold-msg.Num)
				}
				GetTopMgr().UpdData(&person)
				GetServer().SqlGoldLog(person.Uid, -msg.Num, msg.Dec)
				if msg.Dec <= 0 { //! 是其他事件带来的金币改变,需要同步游戏服务器金币
					person.SynchroGold()
				}
			} else {
				person.Card = lib.HF_MaxInt(0, person.Card-msg.Num)
			}
			person.Flush(true)
			*reply = []byte(fmt.Sprintf("%d_%d_%d_%d", card-person.Card, gold-person.Gold, person.Card, person.Gold))
		}
	case "givecard":
		var msg staticfunc.Msg_GiveCard
		json.Unmarshal(data, &msg)
		value, _ := GetServer().DB_GetData("user", msg.Uid, 1)
		if string(value) != "" {
			var person Person
			json.Unmarshal(value, &person)
			if msg.Pid == staticfunc.TYPE_GOLD {
				person.Gold += msg.Num
				if msg.Dec == -5 { //! 后台送的绑定金币
					person.BindGold += msg.Num
				}
				GetTopMgr().UpdData(&person)
				GetServer().SqlGoldLog(person.Uid, msg.Num, msg.Dec)
				if msg.Dec <= 0 || msg.Sync { //! 是其他事件带来的金币改变,需要同步游戏服务器金币
					msg.Sync = person.SynchroGold()
				}
				person.Flush(true)
				*reply = []byte(fmt.Sprintf("%d_%d", person.Card, person.Gold))
				if msg.Sync && msg.Dec > 0 { //! 需要同步,把消息给centerserver去同步
					var _msg staticfunc.Msg_Synchro
					_msg.Uid = person.Uid
					_msg.Gold = person.Gold
					_msg.Card = person.Card
					GetServer().CallCenter("ServerMethod.ServerMsg", "synchro", &_msg)
				}
				GetServer().InsertLog(person.Uid, staticfunc.ADD_GOLD, msg.Num, msg.Ip)
				return nil
			} else if msg.Pid == staticfunc.TYPE_CARD {
				person.Card += msg.Num
				person.Flush(true)
				GetServer().InsertLog(person.Uid, staticfunc.ADD_CARD, msg.Num, msg.Ip)
				*reply = []byte(fmt.Sprintf("%d_%d", person.Card, person.Gold))
				return nil
			}
		}
		*reply = []byte("false")
	case "movecard":
		var msg staticfunc.Msg_GiveCard
		json.Unmarshal(data, &msg)
		value, _ := GetServer().DB_GetData("user", msg.Uid, 1)
		if string(value) != "" {
			var person Person
			json.Unmarshal(value, &person)
			if msg.Pid == staticfunc.TYPE_GOLD { //! 提现的时候不能在游戏内
				if msg.Dec == -4 && person.GameType != 0 {
					*reply = []byte("false")
					return nil
				}

				if msg.Num > person.Gold && (msg.Dec == -1 || msg.Dec == -4) && person.SaveGold > 0 { //! 后台扣金币大于身上的金币，就扣银行的钱
					gold := lib.HF_MinInt(msg.Num-person.Gold, person.SaveGold)
					person.Gold += gold
					person.SaveGold -= gold
					GetServer().SqlGoldLog(person.Uid, gold, -11)
				}
				person.Gold = lib.HF_MaxInt(0, person.Gold-msg.Num)
				if msg.Dec != 77 {
					person.BindGold = lib.HF_MaxInt(0, person.BindGold-msg.Num)
				}
				GetTopMgr().UpdData(&person)
				GetServer().SqlGoldLog(person.Uid, -msg.Num, msg.Dec)
				if msg.Dec <= 0 { //! 是其他事件带来的金币改变,需要同步游戏服务器金币
					person.SynchroGold()
				}
				person.Flush(true)
				GetServer().InsertLog(person.Uid, staticfunc.MOVE_GOLD, msg.Num, msg.Ip)
				*reply = []byte(fmt.Sprintf("%d_%d", person.Card, person.Gold))
				return nil
			} else if msg.Pid == staticfunc.TYPE_CARD {
				person.Card = lib.HF_MaxInt(0, person.Card-msg.Num)
				person.Flush(true)
				GetServer().InsertLog(person.Uid, staticfunc.MOVE_CARD, msg.Num, msg.Ip)
				*reply = []byte(fmt.Sprintf("%d_%d", person.Card, person.Gold))
				return nil
			}
		}
		*reply = []byte("false")
	case "buycard": //! 这是用来接苹果内购的
		var msg staticfunc.Msg_GiveCard
		json.Unmarshal(data, &msg)
		value, _ := GetServer().DB_GetData("user", msg.Uid, 1)
		if string(value) != "" {
			var person Person
			json.Unmarshal(value, &person)
			person.Card += msg.Num
			person.Flush(true)
			GetServer().InsertLog(person.Uid, staticfunc.BUY_CARD, msg.Num, msg.Ip)
			*reply = []byte(fmt.Sprintf("%d_%d", person.Card, person.Gold))
			return nil
		}
		*reply = []byte("false")
	case "getcard":
		var msg staticfunc.Msg_GiveCard
		json.Unmarshal(data, &msg)
		value, _ := GetServer().DB_GetData("user", msg.Uid, 1)
		if string(value) != "" {
			var person Person
			json.Unmarshal(value, &person)
			*reply = []byte(fmt.Sprintf("%d_%d_%d", person.Card, person.Gold, person.SaveGold))
			return nil
		}
		*reply = []byte("false")
	case "findplayer":
		var msg staticfunc.Msg_Uid
		json.Unmarshal(data, &msg)
		value, _ := GetServer().DB_GetData("user", msg.Uid, 1)
		if string(value) != "" {
			var person Person
			json.Unmarshal(value, &person)
			var msg staticfunc.Msg_PersonInfo
			msg.Uid = person.Uid
			msg.Nickname = person.Name
			msg.Headurl = person.Imgurl
			msg.RoomId = person.RoomId
			var son staticfunc.Son_Card
			son.Productid = staticfunc.TYPE_CARD
			son.Repertory_count = person.Card
			msg.Repertory = append(msg.Repertory, son)
			son.Productid = staticfunc.TYPE_GOLD
			son.Repertory_count = person.Gold + person.SaveGold
			msg.Repertory = append(msg.Repertory, son)
			*reply = lib.HF_JtoB(&msg)
			return nil
		}
		*reply = []byte("false")
	case "findplayerfromopenid":
		var msg staticfunc.Msg_Openid
		json.Unmarshal(data, &msg)
		uid := GetServer().DB_GetUid(msg.Openid, "", false, 0)
		if uid != 0 {
			value, _ := GetServer().DB_GetData("user", uid, 1)
			if string(value) != "" {
				var person Person
				json.Unmarshal(value, &person)
				var msg staticfunc.Msg_PersonInfo
				msg.Uid = person.Uid
				msg.Nickname = person.Name
				msg.Headurl = person.Imgurl
				msg.RoomId = person.RoomId
				var son staticfunc.Son_Card
				son.Productid = staticfunc.TYPE_CARD
				son.Repertory_count = person.Card
				msg.Repertory = append(msg.Repertory, son)
				son.Productid = staticfunc.TYPE_GOLD
				son.Repertory_count = person.Gold
				msg.Repertory = append(msg.Repertory, son)
				*reply = lib.HF_JtoB(&msg)
				return nil
			}
		}
		*reply = []byte("false")
	case "closedownplayer1":
		var msg staticfunc.Msg_Uid
		json.Unmarshal(data, &msg)
		value, _ := GetServer().DB_GetData("user", msg.Uid, 1)
		if string(value) != "" {
			GetServer().AddBlack(msg.Uid)
			*reply = []byte("ok")
			return nil
		}
		*reply = []byte("false")
	case "closedownplayer2":
		var msg staticfunc.Msg_Openid
		json.Unmarshal(data, &msg)
		uid := GetServer().DB_GetUid(msg.Openid, "", false, 0)
		if uid != 0 {
			value, _ := GetServer().DB_GetData("user", uid, 1)
			if string(value) != "" {
				GetServer().AddBlack(uid)
				*reply = []byte("ok")
				return nil
			}
		}
		*reply = []byte("false")
	case "opendownplayer1":
		var msg staticfunc.Msg_Uid
		json.Unmarshal(data, &msg)
		value, _ := GetServer().DB_GetData("user", msg.Uid, 1)
		if string(value) != "" {
			GetServer().DelBlack(msg.Uid)
			*reply = []byte("ok")
			return nil
		}
		*reply = []byte("false")
	case "opendownplayer2":
		var msg staticfunc.Msg_Openid
		json.Unmarshal(data, &msg)
		uid := GetServer().DB_GetUid(msg.Openid, "", false, 0)
		if uid != 0 {
			value, _ := GetServer().DB_GetData("user", uid, 1)
			if string(value) != "" {
				GetServer().DelBlack(uid)
				*reply = []byte("ok")
				return nil
			}
		}
		*reply = []byte("false")
	case "gameserver": //! gameserver连接断开
		var msg staticfunc.Msg_GameServer
		json.Unmarshal(data, &msg)
		if msg.InIp != "" {
			GetServer().AddGameServer(msg.Id, msg.InIp, msg.ExIp, msg.Type)
		} else {
			GetServer().DelGameServer(msg.Id)
		}
		*reply = []byte(fmt.Sprintf("%d", GetServer().MoneyMode))
	case "centerserver": //! centerserver连接断开
		var msg staticfunc.Msg_CenterServer
		json.Unmarshal(data, &msg)
		GetServer().InCenterServer = msg.InIp
		GetServer().ExCenterServer = msg.ExIp
		if GetServer().RpcCenterServer == nil && msg.InIp != "" {
			ip := strings.Split(msg.InIp, ":")
			_ip := ip[0] + ":1" + ip[1]
			GetServer().RpcCenterServer = lib.CreateClientPool([]string{_ip}, _ip)
		}
		*reply = []byte(fmt.Sprintf("%d", GetServer().MoneyMode))
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "center连接:", msg.InIp)
	case "delroom": //! 删除房间
		var msg staticfunc.Msg_DelRoom
		json.Unmarshal(data, &msg)
		GetServer().DelGameRoom(msg.Id, msg.RoomId)
		GetCreateRoomMgr().DelRoom(msg.Host, msg.RoomId, msg.Agent)
		if msg.GameType/10000 == 25 {
			index := lib.GetFishMgr().DelDesk(msg.GameType%250000, msg.RoomId)
			var _msg staticfunc.Msg_FishDelDesk1
			_msg.Index = index
			_msg.Lst = lib.GetFishMgr().Get(msg.GameType % 250000)
			GetServer().CallCenter("ServerMethod.ServerMsg", "fishdeldesk", &_msg)
		}
		for i := 0; i < len(msg.Uid); i++ {
			GetNumMgr().DoneGameOne(msg.GameType, msg.Uid[i])
			value, _ := GetServer().DB_GetData("user", msg.Uid[i], 0)
			if string(value) == "" {
				continue
			}

			var person Person
			json.Unmarshal(value, &person)
			if person.RoomId != msg.RoomId { //! 不在这个房间内
				continue
			}
			person.GameId = 0
			person.RoomId = 0
			person.GameType = 0
			person.Flush(false)
		}
		if msg.ClubId != 0 { //! 解散了俱乐部房间
			var _msg staticfunc.Msg_ClubRoomDel
			_msg.ClubId = msg.ClubId
			_msg.RoomId = msg.RoomId
			GetServer().CallCenter("ServerMethod.ServerMsg", "clubroomdel", &_msg)
		}
	case "beginroom": //! 开始房间
		var msg staticfunc.Msg_BeginRoom
		json.Unmarshal(data, &msg)
		GetServer().BeginGameRoom(msg.Id, msg.RoomId)
		GetCreateRoomMgr().SetBegin(msg.Host, msg.RoomId)
		if msg.ClubId != 0 {
			var _msg staticfunc.Msg_ClubRoomDel
			_msg.ClubId = msg.ClubId
			_msg.RoomId = msg.RoomId
			GetServer().CallCenter("ServerMethod.ServerMsg", "clubroombegin", &_msg)
		}
	case "joinfail": //! 加入房间失败
		var msg staticfunc.Msg_JoinFail
		json.Unmarshal(data, &msg)
		value, _ := GetServer().DB_GetData("user", msg.Uid, 0)
		if string(value) == "" {
			return nil
		}
		GetServer().RemoveGameRoom(msg.Id, msg.Room, msg.Uid)
		var person Person
		json.Unmarshal(value, &person)
		GetNumMgr().DoneGameOne(person.GameType, person.Uid)
		person.GameId = 0
		person.RoomId = 0
		person.GameType = 0
		person.Flush(false)
		if msg.GameType/10000 == 25 {
			desk := lib.GetFishMgr().DelPerson(msg.GameType%250000, msg.Room, person.Uid)
			if desk != nil {
				var _msg staticfunc.Msg_FishAddDesk1
				_msg.Desk = *desk
				_msg.Lst = lib.GetFishMgr().Get(msg.GameType % 250000)
				GetServer().CallCenter("ServerMethod.ServerMsg", "fishdelperson", &_msg)
			}
		}
	case "viewroomnum": //
		var msg staticfunc.Msg_ViewRoomNum
		json.Unmarshal(data, &msg)
		GetCreateRoomMgr().Add(msg.Host, msg.RoomId, msg.Num, staticfunc.JS_CreateRoomMem{msg.Uid, msg.Name, msg.Head, 0})
	case "viewroomscore":
		var msg staticfunc.Msg_ViewRoomScore
		json.Unmarshal(data, &msg)
		GetCreateRoomMgr().UpdScore(msg.Host, msg.RoomId, msg.Node)
	case "createroom": //! 俱乐部创建房间
		var msg staticfunc.Msg_ClubCreateRoom
		json.Unmarshal(data, &msg)
		*reply = CreateClubRoom(msg.ClubId, msg.Host, msg.Uid, msg.GameType, msg.Num, msg.Param1, msg.Param2, msg.IP)
	case "clubroomresult": //! 发送战绩
		var msg staticfunc.Msg_ClubRoomResult
		json.Unmarshal(data, &msg)
		GetServer().CallCenter("ServerMethod.ServerMsg", "clubroomresult", &msg)
	case "richnotice":
		var msg staticfunc.Msg_RichNotice
		json.Unmarshal(data, &msg)
		GetServer().CallCenter("ServerMethod.ServerMsg", "richnotice", &msg)
	case "addrobot":
		var msg staticfunc.Msg_AddRobot
		json.Unmarshal(data, &msg)
		GetServer().AddGameRoom(msg.Id, msg.Room, msg.GameType, msg.Num, msg.Uid, msg.IP)
		GetNumMgr().AddGameOne(msg.GameType, msg.Uid)
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "addrobot")
	case "delrobot":
		var msg staticfunc.Msg_DelRobot
		json.Unmarshal(data, &msg)
		GetServer().RemoveGameRoom(msg.Id, msg.Room, msg.Uid)
		GetNumMgr().DoneGameOne(msg.GameType, msg.Uid)
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "delrobot")
	}

	return nil
}
