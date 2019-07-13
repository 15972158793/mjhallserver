package loginserver

import (
	"encoding/json"
	"fmt"
	"lib"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

type JS_AGVideo struct {
	ErrorCode    int               `json:"errorCode"`
	ErrorMessage string            `json:"errorMessage"`
	Result       map[string]string `json:"result"`
}

var AGVideoDealid map[string]int = make(map[string]int)

func AGVideoMsg(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	//clientip := lib.HF_GetHttpIP(req)
	//if !GetServer().IsWhite(clientip, "AGVideoMsg") {
	//	return
	//}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	cmd := req.FormValue("cmd")
	//signature := req.FormValue("signature")
	user := req.FormValue("user")
	requestDate := req.FormValue("requestDate")

	lib.GetLogMgr().Output(lib.LOG_DEBUG, cmd)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, user)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, requestDate)

	if cmd == "CallBalance" {
		var msg JS_AGVideo
		msg.ErrorCode = 0
		msg.Result = make(map[string]string)

		arr := strings.Split(user, "_")
		if len(arr) != 3 {
			msg.ErrorCode = 101
			msg.ErrorMessage = "参数错误"
			w.Write(lib.HF_JtoB(&msg))
			return
		}
		uid := lib.HF_Atoi(arr[1])
		value, _ := GetServer().DB_GetData("user", int64(uid), 1)
		if string(value) == "" {
			msg.ErrorCode = 102
			msg.ErrorMessage = "找不到玩家"
			w.Write(lib.HF_JtoB(&msg))
			return
		}

		var person Person
		json.Unmarshal(value, &person)
		msg.Result["user"] = user
		if GetServer().MoneyMode == 1 {
			msg.Result["money"] = fmt.Sprintf("%f", float32(person.Gold)/100.0)
		} else {
			msg.Result["money"] = fmt.Sprintf("%d", person.Gold)
		}
		msg.Result["responseDate"] = time.Now().Format(lib.TIMEFORMAT)
		w.Write(lib.HF_JtoB(&msg))
	} else if cmd == "PointInout" {
		money := req.FormValue("money")
		dealid := req.FormValue("dealid")
		_type := req.FormValue("type")
		lib.GetLogMgr().Output(lib.LOG_DEBUG, _type)

		var msg JS_AGVideo
		msg.ErrorCode = 0
		msg.Result = make(map[string]string)

		arrtype := strings.Split(_type, "_")
		if len(arrtype) != 4 {
			msg.ErrorCode = 101
			msg.ErrorMessage = "参数错误"
			w.Write(lib.HF_JtoB(&msg))
			return
		}

		gametype := 0
		switch lib.HF_Atoi(arrtype[0]) {
		case 101:
			gametype = 1000000
		case 102:
			gametype = 1000001
		case 103:
			gametype = 1000002
		case 104:
			gametype = 1000003
		case 105:
			gametype = 1000004
		case 106:
			gametype = 1000005
		case 107:
			gametype = 1000006
		case 108:
			gametype = 1000008
		}

		i_money := lib.HF_Atof(money)

		arr := strings.Split(user, "_")
		if len(arr) != 3 {
			msg.ErrorCode = 101
			msg.ErrorMessage = "参数错误"
			w.Write(lib.HF_JtoB(&msg))
			return
		}

		uid := lib.HF_Atoi(arr[1])
		value, _ := GetServer().DB_GetData("user", int64(uid), 1)
		if string(value) == "" {
			msg.ErrorCode = 102
			msg.ErrorMessage = "找不到玩家"
			w.Write(lib.HF_JtoB(&msg))
			return
		}

		var person Person
		json.Unmarshal(value, &person)

		if GetServer().MoneyMode == 1 {
			i_money *= 100
		}
		if person.Gold+int(i_money) < 0 {
			msg.ErrorCode = 103
			msg.ErrorMessage = "金币不足"
			w.Write(lib.HF_JtoB(&msg))
			return
		}
		person.Gold += int(i_money)
		GetServer().SqlGoldLog(person.Uid, int(i_money), gametype)
		person.SynchroGold()
		person.SynchroMoney()
		person.Flush(true)
		AGVideoDealid[dealid] = gametype
		msg.Result["user"] = user
		msg.Result["money"] = money
		msg.Result["responseDate"] = time.Now().Format(lib.TIMEFORMAT)
		msg.Result["dealid"] = dealid
		if GetServer().MoneyMode == 1 {
			msg.Result["cash"] = fmt.Sprintf("%f", float32(person.Gold)/100.0)
		} else {
			msg.Result["cash"] = fmt.Sprintf("%d", person.Gold)
		}
		w.Write(lib.HF_JtoB(&msg))
	} else if cmd == "TimeoutBetReturn" {
		money := req.FormValue("money")
		dealid := req.FormValue("dealid")
		lib.GetLogMgr().Output(lib.LOG_DEBUG, money)

		i_money := lib.HF_Atof(money)

		var msg JS_AGVideo
		msg.ErrorCode = 0
		msg.Result = make(map[string]string)

		if AGVideoDealid[dealid] == 0 {
			msg.ErrorCode = 103
			msg.ErrorMessage = "找不到该订单"
			w.Write(lib.HF_JtoB(&msg))
			return
		}

		arr := strings.Split(user, "_")
		if len(arr) != 3 {
			msg.ErrorCode = 101
			msg.ErrorMessage = "参数错误"
			w.Write(lib.HF_JtoB(&msg))
			return
		}

		uid := lib.HF_Atoi(arr[1])
		value, _ := GetServer().DB_GetData("user", int64(uid), 1)
		if string(value) == "" {
			msg.ErrorCode = 102
			msg.ErrorMessage = "找不到玩家"
			w.Write(lib.HF_JtoB(&msg))
			return
		}

		var person Person
		json.Unmarshal(value, &person)

		if GetServer().MoneyMode == 1 {
			i_money *= 100
		}
		if i_money < 0 {
			i_money = -i_money
		}
		person.Gold += int(i_money)
		GetServer().SqlGoldLog(person.Uid, int(i_money), AGVideoDealid[dealid])
		person.SynchroGold()
		person.SynchroMoney()
		person.Flush(true)
		delete(AGVideoDealid, dealid)
		msg.Result["user"] = user
		msg.Result["money"] = money
		msg.Result["responseDate"] = time.Now().Format(lib.TIMEFORMAT)
		msg.Result["dealid"] = dealid
		if GetServer().MoneyMode == 1 {
			msg.Result["cash"] = fmt.Sprintf("%f", float32(person.Gold)/100.0)
		} else {
			msg.Result["cash"] = fmt.Sprintf("%d", person.Gold)
		}
		w.Write(lib.HF_JtoB(&msg))
	}
}
