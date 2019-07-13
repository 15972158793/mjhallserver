package gameserver

import (
	"encoding/json"
	"lib"
	"runtime/debug"
	"staticfunc"
)

type Msg_ClientLog struct {
	Log string `json:"log"`
}

func OnReceive(self *lib.Session, msg []byte) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	head, _, data, ok := lib.HF_DecodeMsg(msg)
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(self.IP, "消息解析错误1")
		return
	}

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "client:", head, "...", string(data))

	uid := int64(0)
	if self.Person != nil {
		uid = self.Person.(*Person).Uid
	}

	ip := self.IP

	switch head {
	case "joinroom": //! 加入房间
		var msg Msg_JoinRoom
		json.Unmarshal(data, &msg)

		//! 先验证合法性
		person := new(Person)
		self.Person = person
		value := GetServer().DB_GetData("user", msg.Uid)
		if string(value) != "" {
			json.Unmarshal(value, &person)
		} else { //! redis读不到，换服务器获取
			var _msg staticfunc.Msg_Uid
			_msg.Uid = msg.Uid
			result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "getperson", &_msg)
			if err != nil || string(result) == "" {
				self.SafeClose()
				return
			}
			json.Unmarshal(result, &person)
		}
		if person.UnionId != "" && person.UnionId != msg.UnionId {
			var msg Msg_JoinRoomFail
			msg.Result = 3
			self.SendMsg("joinroomfail", &msg)
			self.SafeClose()
			return
		}
		if person.RoomId != msg.Roomid {
			var msg Msg_JoinRoomFail
			msg.Result = 3
			self.SendMsg("joinroomfail", &msg)
			self.SafeClose()
			return
		}
		if !JoinRoom(self, person, msg.MInfo, ip, msg.Param) { //! 加入失败，告诉登陆服务器
			var _msg staticfunc.Msg_JoinFail
			_msg.Uid = msg.Uid
			_msg.Id = GetServer().Con.Id
			_msg.Room = person.RoomId
			GetServer().CallLogin("ServerMethod.ServerMsg", "joinfail", &_msg)
		}
	case "dissmissroom": //! 解散房间
		room := GetRoom(self)
		if room != nil {
			var msg Msg_DissRoom
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("dissmissroom", "", uid, &msg))
		}
	case "nodissmissroom":
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("nodissmissroom", "", uid, nil))
		}
	case "gameready": //! 准备
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gameready", "", uid, nil))
		}
	case "gameunready": //! 取消准备
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gameunready", "", uid, nil))
		}
	case "gameseat": //! 坐下
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gameseat", "", uid, nil))
		}
	case "gamethree": //! 换三张
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameThree
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamethree", "", uid, &msg))
		}
	//!**********************
	case "ChoiceScore": //! 推筒子选分
		room := GetRoom(self)
		if room != nil {
			var msg ChoiceScore
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("ChoiceScore", "", uid, &msg))
		}
	case "gameTnext": //! 推筒子选择下一局
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gameTnext", "", uid, nil))
		}
	//!***********************
	case "gamestart": //！游戏开始
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamestart", "", uid, nil))
		}
	case "lzdbstart": //！游戏开始
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameBets
			json.Unmarshal(data, &msg)
			msg.Uid = uid
			room.Operator(NewRoomMsg("lzdbstart", "", uid, &msg))
		}
	case "gamebets": //! 下注
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameBets
			json.Unmarshal(data, &msg)
			msg.Uid = uid
			room.Operator(NewRoomMsg("gamebets", "", uid, &msg))
		}
	case "gameallbets": //! 跟到底
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameBets
			json.Unmarshal(data, &msg)
			msg.Uid = uid
			room.Operator(NewRoomMsg("gameallbets", "", uid, &msg))
		}
	case "gamebjlbets": //! 下注
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameBJLBets
			json.Unmarshal(data, &msg)
			msg.Uid = uid
			room.Operator(NewRoomMsg("gamebjlbets", "", uid, &msg))
		}
	case "gametdzbets": //! 下注
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameTDZBets
			json.Unmarshal(data, &msg)
			msg.Uid = uid
			room.Operator(NewRoomMsg("gametdzbets", "", uid, &msg))
		}
	case "gamelyccard": //! 捞腌菜游戏操作
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GamePlay
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamelyccard", "", uid, &msg))
		}
	case "gamelycrub": //！捞腌菜搓牌
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamelycrub", "", uid, nil))
		}
	case "gameview": //! 亮牌
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gameview", "", uid, nil))
		}
	case "gameniuniuview": //! 牛牛亮牌
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameNiuNiu_View
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gameniuniuview", "", uid, &msg))
		}
	case "gamefpjhuan":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameFPJ_Huan
			json.Unmarshal(data, &msg)
			msg.Uid = uid
			room.Operator(NewRoomMsg("gamefpjhuan", "", uid, &msg))
		}
	case "gameptjview": //! 拼天九亮牌
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GamePTJ_View
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gameptjview", "", uid, &msg))
		}
	case "gamecompare": //! 比牌
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameCompare
			json.Unmarshal(data, &msg)
			msg.Uid = uid
			room.Operator(NewRoomMsg("gamecompare", "", uid, &msg))
		}
	case "gamerob": //！上庄
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamerob", "", uid, nil))
		}
	case "gamebrttzrob":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameGoldBrTTZ_Deal
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamebrttzrob", "", uid, &msg))
		}
	case "gameredeal": //！下庄
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gameredeal", "", uid, nil))
		}
	case "gamedeal": //! 抢庄
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameDeal
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamedeal", "", uid, &msg))
		}
	case "gametrust": //! 抢庄
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameDeal
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gametrust", "", uid, &msg))
		}
	case "gamedealer": //! 抢庄
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameDealer
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamedealer", "", uid, &msg))
		}
	case "gamerating":
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamerating", "", uid, nil))
		}
	case "gamechangecard": //! 换牌
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameChangeCard
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamechangecard", "", uid, &msg))
		}
	case "gamechange": //! 换牌
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameChange
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamechange", "", uid, &msg))
		}
	case "gamegetcard": //！看发的底牌
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamegetcard", "", uid, nil))
		}
	case "gametranslate": //！换发的底牌
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gametranslate", "", uid, nil))
		}
	case "gamelzdealer": //! 抢庄
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameLzDealer
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamelzdealer", "", uid, &msg))
		}
	case "gamesetboom": //! 扫雷埋雷
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameBoom
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamesetboom", "", uid, &msg))
		}
	case "gamegetboom": //! 扫雷挖雷
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameBoom
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamegetboom", "", uid, &msg))
		}
	case "gamestep": //! 出牌
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameStep
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamestep", "", uid, &msg))
		}
	case "gamesteps": //! 出牌
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameSteps
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamesteps", "", uid, &msg))
		}
	case "gamestepspdk": //! 出牌跑得快
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameStepsPDK
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamestepspdk", "", self.Person.(*Person).Uid, &msg))
		}
	case "gameshgen": //!梭哈跟注
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameBets
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gameshgen", "", self.Person.(*Person).Uid, &msg))
		}
	case "gameshjia": //!梭哈加注
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameBets
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gameshjia", "", self.Person.(*Person).Uid, &msg))
		}
	case "gameshsuo": //!梭哈all in
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameBets
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gameshsuo", "", self.Person.(*Person).Uid, &msg))
		}
	case "gameshqi": //! 梭哈弃牌
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gameshqi", "", self.Person.(*Person).Uid, nil))
		}
	case "gamewzqstep": //! 五子棋
		room := GetRoom(self)
		if room != nil {
			var msg C2S_GameWZQ_Step
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamewzqstep", "", self.Person.(*Person).Uid, &msg))
		}
	case "gamelose": //! 认输
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamelose", "", self.Person.(*Person).Uid, nil))
		}
	case "gamebegin": //! 开始
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamebegin", "", self.Person.(*Person).Uid, nil))
		}
	case "gameend": //! 结束
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gameend", "", self.Person.(*Person).Uid, nil))
		}
	case "gameplayerlist":
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gameplayerlist", "", self.Person.(*Person).Uid, nil))
		}
	case "gameopen": //! 开牌
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gameopen", "", self.Person.(*Person).Uid, nil))
		}
	case "gamepiao": //! 飘分
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GamePiao
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamepiao", "", uid, &msg))
		}
	case "gameplay": //! 游戏操作
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GamePlay
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gameplay", "", uid, &msg))
		}
	case "gamematch": //! 配牌
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameMatch
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamematch", "", uid, &msg))
		}
	case "gamematchs": //! 十三道配牌
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameSSDMatch
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamematchs", "", uid, &msg))
		}
	case "gamecontinue": //！继续游戏
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GamePlay
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamecontinue", "", uid, &msg))
		}
	case "gametbkill":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameGoldBZW_Seat
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gametbkill", "", uid, &msg))
		}
	case "gamedouble": //! 加倍
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameDDZ_Double
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamedouble", "", uid, &msg))
		}
	case "gamecagang": //! 擦杠
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameStep
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamecagang", "", uid, &msg))
		}
	case "gamepeng": //! 碰
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamepeng", "", uid, nil))
		}
	case "gamegang": //! 杠
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameStep
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamegang", "", uid, &msg))
		}
	case "gamechi": //! 吃
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameChi
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamechi", "", uid, &msg))
		}
	case "gamehu": //! 胡
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamehu", "", uid, nil))
		}
	//case "gamezhao": //! 招
	//	room := GetRoom(self)
	//	if room != nil {
	//		var msg Msg_GameHFBH_ClientShao
	//		json.Unmarshal(data, &msg)
	//		room.Operator(NewRoomMsg("gamezhao", "", uid, &msg))
	//	}
	//case "gameshao": //! 绍
	//	room := GetRoom(self)
	//	if room != nil {
	//		var msg Msg_GameHFBH_ClientShao
	//		json.Unmarshal(data, &msg)
	//		room.Operator(NewRoomMsg("gameshao", "", uid, &msg))
	//	}
	//case "gamehfbhchi": //! 吃
	//	room := GetRoom(self)
	//	if room != nil {
	//		var msg Msg_GameHFBH_ClientChi
	//		json.Unmarshal(data, &msg)
	//		room.Operator(NewRoomMsg("gamehfbhchi", "", uid, &msg))
	//	}
	//case "gamezhua": //! 抓
	//	room := GetRoom(self)
	//	if room != nil {
	//		room.Operator(NewRoomMsg("gamezhua", "", uid, nil))
	//	}
	//case "gametuo": //! 拖
	//	room := GetRoom(self)
	//	if room != nil {
	//		var msg Msg_GameHFBH_ClientTuo
	//		json.Unmarshal(data, &msg)
	//		room.Operator(NewRoomMsg("gametuo", "", uid, &msg))
	//	}
	case "gameting": //!
		room := GetRoom(self)

		if room != nil {
			var msg Msg_GameStep
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamehu", "", uid, &msg))
		}
	case "gamezhi": //! 掷筛子
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamezhi", "", uid, nil))
		}
	case "gamezhifinish": //! 掷筛子
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamezhifinish", "", uid, nil))
		}
	case "gameci":
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gameci", "", uid, nil))
		}
	case "gameguo": //! 过
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gameguo", "", uid, nil))
		}
	case "gamekwxview":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameKWX_View
			json.Unmarshal(data, &msg)
			msg.Uid = uid
			room.Operator(NewRoomMsg("gamekwxview", "", uid, &msg))
		}
	case "gamekwxneed":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameKWX_Need
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamekwxneed", "", uid, &msg))
		}
	case "gamekwxmygod":
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamekwxmygod", "", uid, nil))
		}
	case "kwxgetmygod":
		room := GetRoom(self)
		if room != nil {
			var msg GameKWX_GoldGet
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("kwxgetmygod", "", uid, &msg))
		}
	case "gamebzwbets":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameGoldBZW_Bets
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamebzwbets", "", uid, &msg))
		}
	case "gamebzwgoon":
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamebzwgoon", "", uid, nil))
		}
	case "gameopenfire": //! 开炮
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameFishing_OpenFire
			json.Unmarshal(data, &msg)
			msg.Uid = uid
			room.Operator(NewRoomMsg("gameopenfire", "", uid, &msg))
		}
	case "gamesetcannon": //! 设置火炮等级
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameFishing_SetCannon
			json.Unmarshal(data, &msg)
			msg.Uid = uid
			room.Operator(NewRoomMsg("gamesetcannon", "", uid, &msg))
		}
	case "gamehitfish": //! 捕鱼击中
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameFishing_HitFish
			json.Unmarshal(data, &msg)
			msg.Uid = uid
			room.Operator(NewRoomMsg("gamehitfish", "", uid, &msg))
		}
	case "gamebrttzbets":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameGoldBZW_Bets
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamebrttzbets", "", uid, &msg))
		}
	case "gamebrttzgoon":
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamebrttzgoon", "", uid, nil))
		}
	case "gametxlhdbets":
		room := GetRoom(self)
		if room != nil {
			var msg Clint_GameTXLHD_Bets
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gametxlhdbets", "", uid, &msg))
		}
	case "chatroom":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_Chat
			json.Unmarshal(data, &msg)
			if GetServer().Con.ManyGag == 1 && room.Many > 0 && msg.Type == 1 { //! 百人场禁言
				var msg staticfunc.Msg_Err
				msg.Err = "该场次暂时禁止自由发言"
				self.SendMsg("err", &msg)
			} else {
				msg.Uid = uid
				room.Operator(NewRoomMsg("broadcast", "chatroom", uid, &msg))
			}
		}
	case "clientlog":
		var msg Msg_ClientLog
		json.Unmarshal(data, &msg)
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "客户端日志:", msg.Log)
	case "gameping":
		var msg staticfunc.Msg_Null
		self.SendMsg("gameping", &msg)
	case "gameline":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_LinePerson
			json.Unmarshal(data, &msg)
			msg.Uid = uid
			room.Operator(NewRoomMsg("broadcast", "gameline", uid, &msg))

			if self.Person.(*Person) != nil {
				self.Person.(*Person).line = msg.Line
			}
		}
	case "gamebzwseat":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameGoldBZW_Seat
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamebzwseat", "", uid, &msg))
		}
	case "gamebrttzseat":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameGoldBZW_Seat
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamebrttzseat", "", uid, &msg))
		}
	case "gametuoguan":
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gametuoguan", "", uid, nil))
		}
	case "gamenotuoguan":
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamenotuoguan", "", uid, nil))
		}
	case "hidcard": //! 自闷,倒闷
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameBets
			json.Unmarshal(data, &msg)
			msg.Uid = uid
			room.Operator(NewRoomMsg("hidcard", "", uid, &msg))
		}
	case "clearbets": //! 清除下注
		room := GetRoom(self)
		if room != nil {
			var msg Msg_CancleOrder
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("clearbets", "", uid, &msg))
		}
	case "gamehistory": //! 下注记录
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamehistory", "", uid, nil))
		}
	case "gamefpjstep":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameGoldFKFPJ_Step
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamefpjstep", "", uid, &msg))
		}
	case "gamefpjselect":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameGoldFKFPJ_Select
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamefpjselect", "", uid, &msg))
		}
	case "gamefpjget":
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamefpjget", "", uid, nil))
		}
	case "gamefpjguess":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameGoldFKFPJ_Guess
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamefpjguess", "", uid, &msg))
		}
	case "gamefpjup":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameGoldFKFPJ_Select
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamefpjup", "", uid, &msg))
		}
	case "gamefpjquick":
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamefpjquick", "", uid, nil))
		}
	case "gamefpjkeep":
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamefpjkeep", "", uid, nil))
		}
	case "gamefpjunkeep":
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamefpjunkeep", "", uid, nil))
		}
	case "gamefpjheadcolor":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_Client_HeadColor
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamefpjheadcolor", "", uid, &msg))
		}
	case "gamehelp":
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamehelp", "", uid, nil))
		}
	case "gameshzstart":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameStart
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gameshzstart", "", uid, &msg))
		}
	case "gamedeallist":
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamedeallist", "", uid, nil))
		}
	case "gameland":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_GameLand
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gameland", "", uid, &msg))
		}
	case "gamedwsetroom":
		room := GetRoom(self)
		if room != nil {
			var msg Msg_DwSetRoom
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamedwsetroom", "", uid, &msg))
		}
	case "gamedwroomlist":
		room := GetRoom(self)
		if room != nil {
			room.Operator(NewRoomMsg("gamedwroomlist", "", uid, nil))
		}
	case "gamedwsetjack": //!设置电玩奖池
		room := GetRoom(self)
		if room != nil {
			var msg Msg_DwSetJackPot
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamedwsetjack", "", uid, &msg))
		}
	case "gamedwsetpro": //!设置电玩属性
		room := GetRoom(self)
		if room != nil {
			var msg Msg_DwSetPro
			json.Unmarshal(data, &msg)
			room.Operator(NewRoomMsg("gamedwsetpro", "", uid, &msg))
		}
	case "closesession":
		OnClose(self)
	default:
		staticfunc.GetIpBlackMgr().AddIp(ip, "消息解析错误2")
	}
}

func OnClose(self *lib.Session) {
	if self.Person == nil {
		return
	}

	if self.Person.(*Person).session != self {
		return
	}

	self.Person.(*Person).session = nil
	GetPersonMgr().DelPerson(self.Person.(*Person).Uid)

	//! 通知其他玩家该玩家掉线
	if self.Person.(*Person).room != nil {
		var msg Msg_LinePerson
		msg.Uid = self.Person.(*Person).Uid
		msg.Line = false

		go self.Person.(*Person).room.Operator(NewRoomMsg("broadcast", "lineperson", 0, &msg))
	}
}

//! 加入房间
func JoinRoom(self *lib.Session, person *Person, minfo string, ip string, param int) bool {
	oldperson := GetPersonMgr().GetPerson(person.Uid)
	if oldperson != nil { //! 踢掉原来的人
		if self.Person.(*Person) != nil {
			if oldperson.session != nil && oldperson.session != self {
				var msg Msg_ExitRoom
				oldperson.SendMsg("tickroom", &msg)
				oldperson.CloseSession()
			}
		}
		GetPersonMgr().DelPerson(person.Uid)
	}

	if person.GameId == 0 { //! 进错了房间
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "进错房间")
		var msg Msg_JoinRoomFail
		msg.Result = 3
		self.SendMsg("joinroomfail", &msg)
		self.SafeClose()
		return false
	}

	room := GetRoomMgr().GetRoom(person.RoomId)
	if room == nil { //! 没有房间
		_, room = GetRoomMgr().JoinRoom(person.RoomId)
		if room == nil {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "找不到房间3")
			var msg Msg_JoinRoomFail
			msg.Result = 1
			self.SendMsg("joinroomfail", &msg)
			self.SafeClose()
			return false
		}
	}

	person.room = room
	person.session = self
	person.ip = ip
	person.line = true
	person.Param = param
	if minfo != "" {
		json.Unmarshal([]byte(minfo), &person.minfo)

		//! 转换一下地址
		if person.minfo.Province != "" && person.minfo.City != "" {
			province := person.minfo.Province
			sprovince := []byte(province)
			city := person.minfo.City
			scity := []byte(city)

			if len(sprovince) > 3 && len(scity) > 3 {
				province = string(sprovince[0 : len(sprovince)-3])
				city = string(scity[0 : len(scity)-3])

				person.minfo.Address = province + city
				lib.GetLogMgr().Output(lib.LOG_DEBUG, minfo)
			}
		}
	}

	room.Operator(NewRoomMsg("joinroom", "", person.Uid, person))

	return true
}

//! 得到room
func GetRoom(self *lib.Session) *Room {
	if self.Person == nil {
		return nil
	}

	return self.Person.(*Person).room
}
