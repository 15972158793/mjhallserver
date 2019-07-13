package loginserver

import (
	//"crypto/md5"
	//"crypto/rand"
	//"encoding/base64"
	//"encoding/hex"
	"encoding/json"
	"fmt"
	//"io"
	//"io/ioutil"
	"lib"
	"net/http"

	"github.com/garyburd/redigo/redis"
	////"net/url"
	"rjmgr"
	"runtime/debug"
	"staticfunc"
)

//! 得到奖池
func GetJackpot(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetJackpot") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	gametype := lib.HF_Atoi(req.FormValue("gametype"))
	value := int64(0)

	c := GetServer().Redis.Get()
	defer c.Close()

	if gametype/10000 == 4 { //! 豹子王4000
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("systemmoney%d", gametype%10000)))
	} else if gametype/10000 == 12 { //! 摇色子120000
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("yszsysmoney%d", gametype%10000)))
	} else if gametype/10000 == 6 { //! 百人推筒子60000
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("brttzmoney%d", gametype%10000)))
	} else if gametype/10000 == 9 { //! 神仙夺宝90000
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("sxdbsysmoney%d", gametype%10000)))
	} else if gametype/10000 == 20 { //! 名品汇200000
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("mphsysmoney%d", gametype%10000)))
	} else if gametype/10000 == 10 { //! 龙虎斗100001
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("lhdmoney%d", gametype%10000)))
	} else if gametype/10000 == 14 { //! 单双140000
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("tbmoney%d", gametype%10000)))
	} else if gametype/10000 == 16 { //! 龙珠夺宝160000
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("lzdbsysmoney%d", gametype%10000)))
	} else if gametype/10000 == 17 { //! 地穴探宝170000
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("dfwsysmoney%d", gametype%10000)))
	} else if gametype/10000 == 13 { //! 赛马130001
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("saimamoney%d", gametype%10000)))
	} else if gametype/10000 == 21 { //! 红黑大战210000
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("hhdzsysmoney%d", gametype%10000)))
	} else if gametype/10000 == 19 { //! 翻牌机190000-190002
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("fkfpjsysmoney%d", gametype%10000)))
	} else if gametype/10000 == 23 { //! 百人牛牛230000
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("brnnsysmoney%d", gametype%10000)))
	} else if gametype/10000 == 24 { //! 鱼虾蟹240000
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("yxxsysmoney%d", gametype%10000)))
	} else if gametype/10000 == 25 { //! 3D捕鱼250000
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("fishmoney%d", gametype%10000)))
	} else if gametype/10000 == 26 { //! 百家乐260000
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("bjlsysmoney%d", gametype%10000)))
	} else if gametype/10000 == 27 { //! 水浒传270000
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("shzsysmoney%d", gametype%10000)))
	} else if gametype/10000 == 28 { //! 李逵劈鱼280000
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("lkpymoney%d", gametype%10000)))
	} else if gametype/10000 == 31 { //! 大赢家拉霸310000
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("dwdyjlbsysmoney%d", gametype%10000)))
	} else if gametype/10000 == 32 { //! 绝地求生320000
		value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("dwjdqssysmoney%d", gametype%10000)))
	}

	w.Write([]byte(fmt.Sprintf("%d", value)))
	return
}

//! 设置奖池
func SetJackpot(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "SetJackpot") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	gametype := lib.HF_Atoi(req.FormValue("gametype"))
	value := lib.HF_Atoi(req.FormValue("value"))

	csv, ok := staticfunc.GetCsvMgr().Data["game"][gametype]
	if !ok {
		w.Write([]byte("设置失败"))
		return
	}

	config := GetServer().GetGameServer(lib.HF_Atoi(csv["gametype"]))
	if config == nil {
		w.Write([]byte("设置失败"))
		return
	}

	var msg staticfunc.Msg_SetJackpot
	msg.GameType = gametype
	msg.Value = int64(value)
	config.Call("ServerMethod.ServerMsg", "setjackpot", &msg)

	w.Write([]byte("设置成功"))
}

///////////////////////////////////////////////////////////////////////
//! 机器人奖池
func GetRobotJackpot(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetRobotJackpot") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	gametype := lib.HF_Atoi(req.FormValue("gametype"))
	value := int64(0)

	c := GetServer().Redis.Get()
	defer c.Close()

	value, _ = redis.Int64(c.Do("GET", fmt.Sprintf("robotwin_%d", gametype)))

	w.Write([]byte(fmt.Sprintf("%d", value)))
	return
}

//! 设置机器人奖池
func SetRobotJackpot(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "SetRobotJackpot") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	gametype := lib.HF_Atoi(req.FormValue("gametype"))
	value := lib.HF_Atoi(req.FormValue("value"))

	csv, ok := staticfunc.GetCsvMgr().Data["game"][gametype]
	if !ok {
		w.Write([]byte("设置失败"))
		return
	}

	config := GetServer().GetGameServer(lib.HF_Atoi(csv["gametype"]))
	if config == nil {
		w.Write([]byte("设置失败"))
		return
	}

	var msg staticfunc.Msg_SetJackpot
	msg.GameType = gametype
	msg.Value = int64(value)
	config.Call("ServerMethod.ServerMsg", "setrobotjackpot", &msg)

	w.Write([]byte("设置成功"))
}

///////////////////////////////////////////////////////////////////////
var ALLGOLDGAME = []int{1000000, 320000, 300000, 310000, 30000, 20000, 270000, 250000, 220000, 260000, 240000, 290000, 210000, 190000, 230000, 140000, 280000, 80000, 90000, 60000, 70000, 77, 10000, 50000, 100001, 40000, 200000, 110001, 160000, 120000, 170000, 130001}
var ALLROOMGAME = []int{19, 20000, 1, 6, 78, 10, 17, 36, 65, 10000, 24, 51}

//! 游戏类型转热更类型
func GameTypeToAssetType(id int) int {
	if id >= 10000 && id < 20000 { //! 金币卡五星
		return 10000
	} else if id >= 20000 && id < 30000 { //! 金币金花
		return 20000
	} else if id >= 30000 && id < 40000 { //! 金币牛牛
		return 30000
	} else if id >= 40000 && id < 50000 { //! 豹子王
		return 40000
	} else if id >= 50000 && id < 60000 { //! 金币牌九
		return 50000
	} else if id >= 60000 && id < 70000 { //! 百人推筒子
		return 60000
	} else if id >= 70000 && id < 80000 { //! 金币推筒子
		return 70000
	} else if id >= 80000 && id < 90000 { //! 金币跑得快
		return 80000
	} else if id >= 90000 && id < 100000 { //! 神仙夺宝
		return 90000
	} else if id >= 100000 && id < 110000 { //! 龙虎斗
		return 100001
	} else if id >= 110000 && id < 120000 { //! 一夜暴富
		return 110001
	} else if id >= 120000 && id < 130000 { //! 摇色资
		return 120000
	} else if id >= 130000 && id < 140000 { //! 赛马
		return 130001
	} else if id >= 140000 && id < 150000 { //! 单双
		return 140000
	} else if id >= 160000 && id < 170000 { //! 龙珠
		return 160000
	} else if id >= 170000 && id < 180000 { //! 地穴
		return 170000
	} else if id >= 190000 && id < 200000 { //! 翻牌机
		return 190000
	} else if id >= 200000 && id < 210000 { //! 名品会
		return 200000
	} else if id >= 210000 && id < 220000 { //! 红黑
		return 210000
	} else if id >= 220000 && id < 230000 { //! 腾讯龙虎斗
		return 220000
	} else if id >= 230000 && id < 240000 { //! 百人牛牛
		return 230000
	} else if id >= 240000 && id < 250000 { //! 鱼虾蟹
		return 240000
	} else if id >= 250000 && id < 260000 { //! 捕鱼
		return 250000
	} else if id >= 260000 && id < 270000 { //! 百家乐
		return 260000
	} else if id >= 270000 && id < 280000 { //! 水浒传
		return 270000
	} else if id >= 280000 && id < 290000 { //! 李逵劈鱼
		return 280000
	} else if id >= 290000 && id < 300000 { //! 斗地主
		return 290000
	} else if id >= 300000 && id < 310000 { //! 红包扫雷
		return 300000
	} else if id >= 310000 && id < 320000 { //! 大赢家拉霸
		return 310000
	} else if id >= 320000 && id < 330000 { //! 绝地求生
		return 320000
	} else if id >= 1000 && id < 2000 { //! 金花包厢
		return 20000
	} else if id >= 2000 && id < 3000 { //! 单双包厢
		return 140000
	} else if id >= 3000 && id < 4000 { //! 牛牛包厢
		return 30000
	} else if id >= 4000 && id < 5000 { //! 跑得快包厢
		return 4000
	} else if id == 1 || id == 16 { //! 牛牛房卡
		return 1
	} else if id >= 2 && id <= 5 || id == 13 || id == 14 { //! 卡五星房卡
		return 10000
	} else if id == 6 || id == 8 { //! 斗地主房卡
		return 6
	} else if id == 7 || id == 25 { //! 金花房卡
		return 20000
	} else if id == 10 { //! 十点半
		return 10
	} else if id == 17 { //! 推筒子房卡
		return 17
	} else if id == 19 { //! 牌九房卡
		return 19
	} else if id == 24 { //! 三公房卡
		return 24
	} else if id == 36 { //! 扫雷房卡
		return 36
	} else if id >= 65 && id <= 69 { //! 牛元帅
		return 65
	} else if id == 77 { //! 五子棋
		return 77
	} else if id == 78 { //! 跑得快房卡
		return 78
	} else if id == 51 { //! 八张清
		return 51
	} else if id >= 1000000 && id <= 1000100 {
		return 1000000
	}
	return id
}

//! 是否允许打开这个游戏
func IsGameOk(gametype int) bool {
	if len(rjmgr.GetRJMgr().GoldGame) == 0 && len(rjmgr.GetRJMgr().RoomGame) == 0 {
		return true
	}

	assettype := GameTypeToAssetType(gametype)
	isgold := true //! 是金币场
	if gametype < 1000 && gametype != 77 {
		isgold = false
	}

	if isgold {
		if len(rjmgr.GetRJMgr().GoldGame) == 0 {
			return true
		}
		for i := 0; i < len(rjmgr.GetRJMgr().GoldGame); i++ {
			if rjmgr.GetRJMgr().GoldGame[i] == assettype {
				return true
			}
		}
	} else {
		if len(rjmgr.GetRJMgr().RoomGame) == 0 {
			return true
		}
		for i := 0; i < len(rjmgr.GetRJMgr().RoomGame); i++ {
			if rjmgr.GetRJMgr().RoomGame[i] == assettype {
				return true
			}
		}
	}

	return false
}

//! 得到打开的游戏
func InitOpenGame() {
	c := GetServer().Redis.Get()
	defer c.Close()

	value, err := redis.Bytes(c.Do("GET", "gamemode"))
	if err == nil {
		json.Unmarshal(value, &GetServer().GameMode)

		for i := 0; i < len(GetServer().GameMode.GoldGame); i++ {
			if !IsGameOk(GetServer().GameMode.GoldGame[i]) {
				copy(GetServer().GameMode.GoldGame[i:], GetServer().GameMode.GoldGame[i+1:])
				GetServer().GameMode.GoldGame = GetServer().GameMode.GoldGame[:len(GetServer().GameMode.GoldGame)-1]
			} else {
				i++
			}
		}
		for i := 0; i < len(GetServer().GameMode.RoomGame); i++ {
			if !IsGameOk(GetServer().GameMode.RoomGame[i]) {
				copy(GetServer().GameMode.RoomGame[i:], GetServer().GameMode.RoomGame[i+1:])
				GetServer().GameMode.RoomGame = GetServer().GameMode.RoomGame[:len(GetServer().GameMode.RoomGame)-1]
			} else {
				i++
			}
		}
	} else {
		if len(rjmgr.GetRJMgr().GoldGame) == 0 {
			lib.HF_DeepCopy(&GetServer().GameMode.GoldGame, &ALLGOLDGAME)
		} else {
			lib.HF_DeepCopy(&GetServer().GameMode.GoldGame, &rjmgr.GetRJMgr().GoldGame)
		}

		if len(rjmgr.GetRJMgr().RoomGame) == 0 {
			lib.HF_DeepCopy(&GetServer().GameMode.RoomGame, &ALLROOMGAME)
		} else {
			lib.HF_DeepCopy(&GetServer().GameMode.RoomGame, &rjmgr.GetRJMgr().RoomGame)
		}
	}
}

//! 是否开放
func IsOpenGame(gametype int) bool {
	assettype := GameTypeToAssetType(gametype)

	lib.GetLogMgr().Output(lib.LOG_ERROR, "gametype:", gametype)
	lib.GetLogMgr().Output(lib.LOG_ERROR, "assettype:", assettype)
	lib.GetLogMgr().Output(lib.LOG_ERROR, "goldgame:", GetServer().GameMode.GoldGame)

	isgold := true //! 是金币场
	if gametype < 1000 && gametype != 77 {
		isgold = false
	}

	if isgold {
		for i := 0; i < len(GetServer().GameMode.GoldGame); i++ {
			if GetServer().GameMode.GoldGame[i] == assettype {
				return true
			}
		}
	} else {
		for i := 0; i < len(GetServer().GameMode.RoomGame); i++ {
			if GetServer().GameMode.RoomGame[i] == assettype {
				return true
			}
		}
	}
	return false
}

//! 转换成客户端数组
func GameModeToClient() staticfunc.GameMode {
	var gamemode staticfunc.GameMode
	lib.HF_DeepCopy(&gamemode.GoldGame, &GetServer().GameMode.GoldGame)
	lib.HF_DeepCopy(&gamemode.RoomGame, &GetServer().GameMode.RoomGame)
	for i := 0; i < len(gamemode.GoldGame); i++ {
		if gamemode.GoldGame[i] == 20000 {
			gamemode.GoldGame[i] = 7
		} else if gamemode.GoldGame[i] == 10000 {
			gamemode.GoldGame[i] = 2
		}
	}
	for i := 0; i < len(gamemode.RoomGame); i++ {
		if gamemode.RoomGame[i] == 20000 {
			gamemode.RoomGame[i] = 7
		} else if gamemode.RoomGame[i] == 10000 {
			gamemode.RoomGame[i] = 2
		}
	}
	return gamemode
}
