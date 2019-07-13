package gameserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"lib"
	"runtime/debug"
	"staticfunc"
)

type S2C_Err struct {
	Info string `json:"info"`
}

type S2C_RecordDDZList struct { //! 得到战报
	Info []Son_RecordDDZCell `json:"info"`
}
type Son_RecordDDZCell struct {
	Roomid  int                   `json:"roomid"`
	Person  []Son_RecordDDZPerson `json:"person"`
	Time    int64                 `json:"time"`
	MaxStep int                   `json:"maxstep"`
	Type    int                   `json:"type"`
}
type Son_RecordDDZPerson struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Score int    `json:"score"`
	Head  string `json:"head"`
	Total int    `json:"total"`
}

type S2C_RecordQPDDZList struct { //! 得到战报
	Info []Son_RecordQPDDZCell `json:"info"`
}
type Son_RecordQPDDZCell struct {
	Roomid  int                     `json:"roomid"`
	Person  []Son_RecordQPDDZPerson `json:"person"`
	Time    int64                   `json:"time"`
	MaxStep int                     `json:"maxstep"`
	Type    int                     `json:"type"`
}
type Son_RecordQPDDZPerson struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Score int    `json:"score"`
	Head  string `json:"head"`
	Total int    `json:"total"`
}

type S2C_RecordKWX1 struct { //! 得到战报
	Info []Son_RecordKWX1 `json:"info"`
}
type Son_RecordKWX1 struct {
	Roomid  int                `json:"roomid"`
	Person  []Son_RecordPerson `json:"person"`
	Time    int64              `json:"time"`
	MaxStep int                `json:"maxstep"`
	Param1  int                `json:"param1"`
	Param2  int                `json:"param2"`
}
type Son_RecordPerson struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Score int    `json:"score"`
	Total int    `json:"total"`
}

type S2C_RecordKWX2 struct { //! 得到战报
	Info string `json:"info"`
}

type S2C_Record struct { //! 得到战报
	Type int      `json:"type"`
	Info []string `json:"info"`
}

func GetErr(info string) []byte {
	var msg S2C_Err
	msg.Info = info

	return lib.HF_EncodeMsg("err", &msg, false)
}

type RecordMethod int

func (self *RecordMethod) RecordMsg(args *staticfunc.Rpc_Args, reply *[]byte) error {
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
	case "record":
		*reply = Record(msg.Uid, msg.Type)
	case "recordkwx1":
		*reply = RecordKWX1(msg.Uid, msg.Type)
	case "recordkwx2":
		*reply = RecordKWX2(msg.Uid, msg.Type)
	case "recordgc1":
		*reply = RecordGC1(msg.Uid, msg.Type)
	case "recordgc2":
		*reply = RecordGC2(msg.Uid, msg.Type)
	case "recordszmj1":
		*reply = RecordSZMJ1(msg.Uid, msg.Type)
	case "recordszmj2":
		*reply = RecordSZMJ2(msg.Uid, msg.Type)
	case "recordxzdd1":
		*reply = RecordXZDD1(msg.Uid, msg.Type)
	case "recordxzdd2":
		*reply = RecordXZDD2(msg.Uid, msg.Type)
	case "recordddz1":
		*reply = RecordDDZ1(msg.Uid, msg.Type)
	case "recordddz2":
		*reply = RecordDDZ2(msg.Uid, msg.Type)
	case "recordzjh1":
		*reply = RecordZJH1(msg.Uid, msg.Type)
	case "recordzjh2":
		*reply = RecordZJH2(msg.Uid, msg.Type)
	case "recordcsmj1":
		*reply = RecordCSMJ1(msg.Uid, msg.Type)
	case "recordcsmj2":
		*reply = RecordCSMJ2(msg.Uid, msg.Type)
	case "recordaqmj1":
		*reply = RecordAQMJ1(msg.Uid, msg.Type)
	case "recordaqmj2":
		*reply = RecordAQMJ2(msg.Uid, msg.Type)
	case "recordgymj1":
		*reply = RecordGYMJ1(msg.Uid, msg.Type)
	case "recordgymj2":
		*reply = RecordGYMJ2(msg.Uid, msg.Type)
	case "recordsyhmj1":
		*reply = RecordSYHMJ1(msg.Uid, msg.Type)
	case "recordsyhmj2":
		*reply = RecordSYHMJ2(msg.Uid, msg.Type)
	case "recorddbd1":
		*reply = RecordDBD1(msg.Uid, msg.Type)
	case "recorddbd2":
		*reply = RecordDBD2(msg.Uid, msg.Type)
	case "recordqpddz1":
		*reply = RecordQPDDZ1(msg.Uid, msg.Type)
	case "recordqpddz2":
		*reply = RecordQPDDZ2(msg.Uid, msg.Type)
	case "recordtsmj1":
		*reply = RecordTSMJ1(msg.Uid, msg.Type)
	case "recordtsmj2":
		*reply = RecordTSMJ2(msg.Uid, msg.Type)
	case "recordncmj1":
		*reply = RecordNCMJ1(msg.Uid, msg.Type)
	case "recordncmj2":
		*reply = RecordNCMJ2(msg.Uid, msg.Type)
	case "recordgdcsmj1":
		*reply = RecordGDCSMJ1(msg.Uid, msg.Type)
	case "recordgdcsmj2":
		*reply = RecordGDCSMJ2(msg.Uid, msg.Type)
	case "recordxtrrhh1":
		*reply = RecordXTRRHH1(msg.Uid, msg.Type)
	case "recordxtrrhh2":
		*reply = RecordXTRRHH2(msg.Uid, msg.Type)
	case "recordbdddz1":
		*reply = RecordBDDDZ1(msg.Uid, msg.Type)
	case "recordbdddz2":
		*reply = RecordBDDDZ2(msg.Uid, msg.Type)
	case "recordhlgc1":
		*reply = RecordHLGC1(msg.Uid, msg.Type)
	case "recordhlgc2":
		*reply = RecordHLGC2(msg.Uid, msg.Type)
	}
	return nil
}

//! 得到战报
func Record(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][_type]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	maxrc := lib.HF_Atoi(csv["maxrc"])
	if _type == 10000 {
		maxrc = 20
	}
	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, maxrc))

	var msg S2C_Record
	msg.Type = _type
	if err == nil {
		for _, v := range values {
			if _type == 7 {
				var record staticfunc.Rec_GameZJH_Info
				json.Unmarshal(v.([]byte), &record)
				if record.Time > 0 && record.Time < 1504019026 {
					continue
				}
			}
			msg.Info = append(msg.Info, string(v.([]byte)))
		}
	} else {
		msg.Info = make([]string, 0)
	}

	return lib.HF_EncodeMsg("record", &msg, true)
}

func RecordKWX1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][2]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordKWX1
	if err == nil {
		for _, v := range values {
			var info Son_RecordKWX1
			var record Rec_GameKWX_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordPerson
				person.Name = record.Person[i].Name
				person.Uid = record.Person[i].Uid
				person.Score = record.Person[i].Score
				person.Head = record.Person[i].Head
				person.Total = record.Person[i].Total
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	} else {
		msg.Info = make([]Son_RecordKWX1, 0)
	}

	return lib.HF_EncodeMsg("recordkwx1", &msg, true)
}

func RecordKWX2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][2]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordKWX2
	if err == nil {
		for _, v := range values {
			var record Rec_GameKWX_Info
			json.Unmarshal(v.([]byte), &record)
			if _type != record.Roomid {
				continue
			}
			msg.Info = string(v.([]byte))
		}
	}

	return lib.HF_EncodeMsg("recordkwx2", &msg, true)
}

func RecordGC1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][12]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordKWX1
	if err == nil {
		for _, v := range values {
			var info Son_RecordKWX1
			var record staticfunc.Rec_GameGC_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordPerson
				person.Name = record.Person[i].Name
				person.Uid = record.Person[i].Uid
				person.Score = record.Person[i].Score
				person.Head = record.Person[i].Head
				person.Total = record.Person[i].Total
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	} else {
		msg.Info = make([]Son_RecordKWX1, 0)
	}

	return lib.HF_EncodeMsg("recordgc1", &msg, true)
}

//!
func RecordGC2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][12]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordKWX2
	if err == nil {
		for _, v := range values {
			var record staticfunc.Rec_GameGC_Info
			json.Unmarshal(v.([]byte), &record)
			if _type != record.Roomid {
				continue
			}
			msg.Info = string(v.([]byte))
		}
	}

	return lib.HF_EncodeMsg("recordgc2", &msg, true)
}

func RecordSZMJ1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][15]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))
	var msg S2C_RecordKWX1
	if err == nil {
		for _, v := range values {
			var info Son_RecordKWX1
			var record staticfunc.Rec_GameSZMJ_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordPerson
				person.Name = record.Person[i].Name
				person.Uid = record.Person[i].Uid
				person.Score = record.Person[i].Score
				person.Head = record.Person[i].Head
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	} else {
		msg.Info = make([]Son_RecordKWX1, 0)
	}

	return lib.HF_EncodeMsg("recordgc1", &msg, true)
}

//!
func RecordSZMJ2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][15]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordKWX2
	if err == nil {
		for _, v := range values {
			var record staticfunc.Rec_GameSZMJ_Info
			json.Unmarshal(v.([]byte), &record)
			if _type != record.Roomid {
				continue
			}
			msg.Info = string(v.([]byte))
		}
	}

	return lib.HF_EncodeMsg("recordgc2", &msg, true)
}

//! 血战到底
func RecordXZDD1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][18]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))
	var msg S2C_RecordKWX1
	if err == nil {
		for _, v := range values {
			var info Son_RecordKWX1
			var record staticfunc.Rec_GameXZDD_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			info.Param1 = record.Param1
			info.Param2 = record.Param2
			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordPerson
				person.Name = record.Person[i].Name
				person.Uid = record.Person[i].Uid
				person.Score = record.Person[i].Score
				person.Head = record.Person[i].Head
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	} else {
		msg.Info = make([]Son_RecordKWX1, 0)
	}

	return lib.HF_EncodeMsg("recordxzdd1", &msg, true)
}

//!
func RecordXZDD2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][18]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordKWX2
	if err == nil {
		for _, v := range values {
			var record staticfunc.Rec_GameXZDD_Info
			json.Unmarshal(v.([]byte), &record)
			if _type != record.Roomid {
				continue
			}
			msg.Info = string(v.([]byte))
		}
	}

	return lib.HF_EncodeMsg("recordxzdd2", &msg, true)
}

//! 斗地主列表
func RecordDDZ1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][6]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordDDZList
	msg.Info = make([]Son_RecordDDZCell, 0)
	if err == nil {
		for _, v := range values {
			var info Son_RecordDDZCell
			var record staticfunc.Rec_GameDDZ_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			info.MaxStep = record.MaxStep
			info.Type = record.Type
			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordDDZPerson
				person.Uid = record.Person[i].Uid
				person.Name = record.Person[i].Name
				person.Score = record.Person[i].Score
				person.Head = record.Person[i].Head
				person.Total = record.Person[i].Total
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	}

	return lib.HF_EncodeMsg("recordddz1", &msg, true)
}

//! 斗地主详细战报
func RecordDDZ2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][6]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	if err == nil {
		for _, v := range values {
			var msg staticfunc.Rec_GameDDZ_Info
			json.Unmarshal(v.([]byte), &msg)
			if _type != msg.Roomid {
				continue
			}
			return lib.HF_EncodeMsg("recordddz2", &msg, true)
		}
	}

	return GetErr("未找到战报")
}

//! 斗地主列表
func RecordBDDDZ1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][75]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordDDZList
	msg.Info = make([]Son_RecordDDZCell, 0)
	if err == nil {
		for _, v := range values {
			var info Son_RecordDDZCell
			var record staticfunc.Rec_GameDDZ_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			info.MaxStep = record.MaxStep
			info.Type = record.Type
			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordDDZPerson
				person.Uid = record.Person[i].Uid
				person.Name = record.Person[i].Name
				person.Score = record.Person[i].Score
				person.Head = record.Person[i].Head
				person.Total = record.Person[i].Total
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	}

	return lib.HF_EncodeMsg("recordbdddz1", &msg, true)
}

//! 斗地主详细战报
func RecordBDDDZ2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][75]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	if err == nil {
		for _, v := range values {
			var msg staticfunc.Rec_GameDDZ_Info
			json.Unmarshal(v.([]byte), &msg)
			if _type != msg.Roomid {
				continue
			}
			return lib.HF_EncodeMsg("recordbdddz2", &msg, true)
		}
	}

	return GetErr("未找到战报")
}

//! 枪炮斗地主列表
func RecordQPDDZ1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][6]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordQPDDZList
	msg.Info = make([]Son_RecordQPDDZCell, 0)
	if err == nil {
		for _, v := range values {
			var info Son_RecordQPDDZCell
			var record staticfunc.Rec_GameQPDDZ_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			info.MaxStep = record.MaxStep
			info.Type = record.Type
			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordQPDDZPerson
				person.Uid = record.Person[i].Uid
				person.Name = record.Person[i].Name
				person.Score = record.Person[i].Score
				person.Head = record.Person[i].Head
				person.Total = record.Person[i].Total
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	}

	return lib.HF_EncodeMsg("recordqpddz1", &msg, true)
}

//! 枪炮斗地主详细战报
func RecordQPDDZ2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][6]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	if err == nil {
		for _, v := range values {
			var msg staticfunc.Rec_GameQPDDZ_Info
			json.Unmarshal(v.([]byte), &msg)
			if _type != msg.Roomid {
				continue
			}
			return lib.HF_EncodeMsg("recordqpddz2", &msg, true)
		}
	}

	return GetErr("未找到战报")
}

//! 斗板凳列表
func RecordDBD1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][6]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordDDZList
	msg.Info = make([]Son_RecordDDZCell, 0)
	if err == nil {
		for _, v := range values {
			var info Son_RecordDDZCell
			var record staticfunc.Rec_GameDDZ_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			info.MaxStep = record.MaxStep
			info.Type = record.Type
			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordDDZPerson
				person.Uid = record.Person[i].Uid
				person.Name = record.Person[i].Name
				person.Score = record.Person[i].Score
				person.Head = record.Person[i].Head
				person.Total = record.Person[i].Total
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	}

	return lib.HF_EncodeMsg("recorddbd1", &msg, true)
}

func RecordDBD2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][6]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	if err == nil {
		for _, v := range values {
			var msg staticfunc.Rec_GameDDZ_Info
			json.Unmarshal(v.([]byte), &msg)
			if _type != msg.Roomid {
				continue
			}
			return lib.HF_EncodeMsg("recorddbd2", &msg, true)
		}
	}

	return GetErr("未找到战报")
}

func RecordZJH1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][7]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordDDZList
	msg.Info = make([]Son_RecordDDZCell, 0)
	if err == nil {
		for _, v := range values {
			var info Son_RecordDDZCell
			var record staticfunc.Rec_GameZJH_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			info.MaxStep = record.MaxStep
			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordDDZPerson
				person.Uid = record.Person[i].Uid
				person.Name = record.Person[i].Name
				person.Score = record.Person[i].Score
				person.Head = record.Person[i].Head
				person.Total = record.Person[i].Total
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	}

	return lib.HF_EncodeMsg("recordzjh1", &msg, true)
}

func RecordZJH2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][7]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	if err == nil {
		for _, v := range values {
			var msg staticfunc.Rec_GameZJH_Info
			json.Unmarshal(v.([]byte), &msg)
			if _type != msg.Roomid {
				continue
			}
			return lib.HF_EncodeMsg("recordzjh2", &msg, true)
		}
	}

	return GetErr("未找到战报")
}

//! 常熟麻将
func RecordCSMJ1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][32]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))
	var msg S2C_RecordKWX1
	if err == nil {
		for _, v := range values {
			var info Son_RecordKWX1
			var record staticfunc.Rec_GameCSMJ_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			info.Param1 = record.Param1
			info.Param2 = record.Param2
			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordPerson
				person.Name = record.Person[i].Name
				person.Uid = record.Person[i].Uid
				person.Score = record.Person[i].Score
				person.Head = record.Person[i].Head
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	} else {
		msg.Info = make([]Son_RecordKWX1, 0)
	}

	return lib.HF_EncodeMsg("recordcsmj1", &msg, true)
}

//!
func RecordCSMJ2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][32]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordKWX2
	if err == nil {
		for _, v := range values {
			var record staticfunc.Rec_GameCSMJ_Info
			json.Unmarshal(v.([]byte), &record)
			if _type != record.Roomid {
				continue
			}
			msg.Info = string(v.([]byte))
		}
	}

	return lib.HF_EncodeMsg("recordcsmj2", &msg, true)
}

//! 安慶麻将
func RecordAQMJ1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][27]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))
	var msg S2C_RecordKWX1
	if err == nil {
		for _, v := range values {
			var info Son_RecordKWX1
			var record staticfunc.Rec_GameAQMJ_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			info.Param1 = record.Param1
			info.Param2 = record.Param2
			info.MaxStep = record.MaxStep
			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordPerson
				person.Name = record.Person[i].Name
				person.Uid = record.Person[i].Uid
				person.Score = record.Person[i].Score
				person.Total = record.Person[i].Total
				person.Head = record.Person[i].Head
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	} else {
		msg.Info = make([]Son_RecordKWX1, 0)
	}

	return lib.HF_EncodeMsg("recordaqmj1", &msg, true)
}

//!
func RecordAQMJ2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][27]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordKWX2
	if err == nil {
		for _, v := range values {
			var record staticfunc.Rec_GameAQMJ_Info
			json.Unmarshal(v.([]byte), &record)
			if _type != record.Roomid {
				continue
			}
			msg.Info = string(v.([]byte))
		}
	}

	return lib.HF_EncodeMsg("recordaqmj2", &msg, true)
}

//! 涡阳麻将
func RecordGYMJ1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][27]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))
	var msg S2C_RecordKWX1
	if err == nil {
		for _, v := range values {
			var info Son_RecordKWX1
			var record staticfunc.Rec_GameGYMJ_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			info.Param1 = record.Param1
			info.Param2 = record.Param2
			info.MaxStep = record.MaxStep
			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordPerson
				person.Name = record.Person[i].Name
				person.Uid = record.Person[i].Uid
				person.Score = record.Person[i].Score
				person.Total = record.Person[i].Total
				person.Head = record.Person[i].Head
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	} else {
		msg.Info = make([]Son_RecordKWX1, 0)
	}

	return lib.HF_EncodeMsg("recordgymj1", &msg, true)
}

//!
func RecordGYMJ2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][27]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordKWX2
	if err == nil {
		for _, v := range values {
			var record staticfunc.Rec_GameGYMJ_Info
			json.Unmarshal(v.([]byte), &record)
			if _type != record.Roomid {
				continue
			}
			msg.Info = string(v.([]byte))
		}
	}

	return lib.HF_EncodeMsg("recordgymj2", &msg, true)
}

//! 上虞花麻将
func RecordSYHMJ1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][41]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))
	var msg S2C_RecordKWX1
	if err == nil {
		for _, v := range values {
			var info Son_RecordKWX1
			var record staticfunc.Rec_GameSYHMJ_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			info.Param1 = record.Param1
			info.Param2 = record.Param2
			info.MaxStep = record.MaxStep

			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordPerson
				person.Name = record.Person[i].Name
				person.Uid = record.Person[i].Uid
				person.Score = record.Person[i].Score
				person.Head = record.Person[i].Head
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	} else {
		msg.Info = make([]Son_RecordKWX1, 0)
	}

	return lib.HF_EncodeMsg("recordsyhmj1", &msg, true)
}

//!
func RecordSYHMJ2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][41]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordKWX2
	if err == nil {
		for _, v := range values {
			var record staticfunc.Rec_GameSYHMJ_Info
			json.Unmarshal(v.([]byte), &record)
			if _type != record.Roomid {
				continue
			}
			msg.Info = string(v.([]byte))
		}
	}

	return lib.HF_EncodeMsg("recordsyhmj2", &msg, true)
}

//!
func RecordTSMJ1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][66]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))
	var msg S2C_RecordKWX1
	if err == nil {
		for _, v := range values {
			var info Son_RecordKWX1
			var record staticfunc.Rec_GameTSMJ_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			info.Param1 = record.Param1
			info.Param2 = record.Param2
			info.MaxStep = record.MaxStep

			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordPerson
				person.Name = record.Person[i].Name
				person.Uid = record.Person[i].Uid
				person.Score = record.Person[i].Score
				person.Head = record.Person[i].Head
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	} else {
		msg.Info = make([]Son_RecordKWX1, 0)
	}

	return lib.HF_EncodeMsg("recordtsmj1", &msg, true)
}

//!
func RecordTSMJ2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][66]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordKWX2
	if err == nil {
		for _, v := range values {
			var record staticfunc.Rec_GameTSMJ_Info
			json.Unmarshal(v.([]byte), &record)
			if _type != record.Roomid {
				continue
			}
			msg.Info = string(v.([]byte))
		}
	}

	return lib.HF_EncodeMsg("recordtsmj2", &msg, true)
}

//!
func RecordNCMJ1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][33]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))
	var msg S2C_RecordKWX1
	if err == nil {
		for _, v := range values {
			var info Son_RecordKWX1
			var record staticfunc.Rec_GameTSMJ_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			info.Param1 = record.Param1
			info.Param2 = record.Param2
			info.MaxStep = record.MaxStep

			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordPerson
				person.Name = record.Person[i].Name
				person.Uid = record.Person[i].Uid
				person.Score = record.Person[i].Score
				person.Head = record.Person[i].Head
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	} else {
		msg.Info = make([]Son_RecordKWX1, 0)
	}

	return lib.HF_EncodeMsg("recordncmj1", &msg, true)
}

//!
func RecordNCMJ2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][33]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordKWX2
	if err == nil {
		for _, v := range values {
			var record staticfunc.Rec_GameTSMJ_Info
			json.Unmarshal(v.([]byte), &record)
			if _type != record.Roomid {
				continue
			}
			msg.Info = string(v.([]byte))
		}
	}

	return lib.HF_EncodeMsg("recordncmj2", &msg, true)
}

//!潮汕麻将
func RecordGDCSMJ1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][81]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))
	var msg S2C_RecordKWX1
	if err == nil {
		for _, v := range values {
			var info Son_RecordKWX1
			var record staticfunc.Rec_GameGDCSMJ_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			info.Param1 = record.Param1
			info.Param2 = record.Param2
			info.MaxStep = record.MaxStep

			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordPerson
				person.Name = record.Person[i].Name
				person.Uid = record.Person[i].Uid
				person.Score = record.Person[i].Score
				person.Head = record.Person[i].Head
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	} else {
		msg.Info = make([]Son_RecordKWX1, 0)
	}

	return lib.HF_EncodeMsg("recordgdcsmj1", &msg, true)
}

//!
func RecordGDCSMJ2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][81]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordKWX2
	if err == nil {
		for _, v := range values {
			var record staticfunc.Rec_GameGDCSMJ_Info
			json.Unmarshal(v.([]byte), &record)
			if _type != record.Roomid {
				continue
			}
			msg.Info = string(v.([]byte))
		}
	}

	return lib.HF_EncodeMsg("recordgdcsmj2", &msg, true)
}

//!仙桃人人晃晃
func RecordXTRRHH1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][82]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))
	var msg S2C_RecordKWX1
	if err == nil {
		for _, v := range values {
			var info Son_RecordKWX1
			var record staticfunc.Rec_GameXTRRHH_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			info.Param1 = record.Param1
			info.Param2 = record.Param2
			info.MaxStep = record.MaxStep

			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordPerson
				person.Name = record.Person[i].Name
				person.Uid = record.Person[i].Uid
				person.Score = record.Person[i].Score
				person.Total = record.Person[i].Total
				person.Head = record.Person[i].Head
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	} else {
		msg.Info = make([]Son_RecordKWX1, 0)
	}

	return lib.HF_EncodeMsg("recordxtrrhh1", &msg, true)
}

//!
func RecordXTRRHH2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][82]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordKWX2
	if err == nil {
		for _, v := range values {
			var record staticfunc.Rec_GameXTRRHH_Info
			json.Unmarshal(v.([]byte), &record)
			if _type != record.Roomid {
				continue
			}
			msg.Info = string(v.([]byte))
		}
	}

	return lib.HF_EncodeMsg("recordxtrrhh2", &msg, true)
}

func RecordHLGC1(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][83]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordKWX1
	if err == nil {
		for _, v := range values {
			var info Son_RecordKWX1
			var record staticfunc.Rec_GameGC_Info
			json.Unmarshal(v.([]byte), &record)
			info.Roomid = record.Roomid
			info.Time = record.Time
			for i := 0; i < len(record.Person); i++ {
				var person Son_RecordPerson
				person.Name = record.Person[i].Name
				person.Uid = record.Person[i].Uid
				person.Score = record.Person[i].Score
				person.Head = record.Person[i].Head
				person.Total = record.Person[i].Total
				info.Person = append(info.Person, person)
			}
			msg.Info = append(msg.Info, info)
		}
	} else {
		msg.Info = make([]Son_RecordKWX1, 0)
	}

	return lib.HF_EncodeMsg("recordhlgc1", &msg, true)
}

//!
func RecordHLGC2(uid int64, _type int) []byte {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][83]
	if !ok {
		return []byte("")
	}

	c := GetServer().RcRedis.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("LRANGE", fmt.Sprintf("%s_%d", csv["rctable"], uid), 0, lib.HF_Atoi(csv["maxrc"])))

	var msg S2C_RecordKWX2
	if err == nil {
		for _, v := range values {
			var record staticfunc.Rec_GameGC_Info
			json.Unmarshal(v.([]byte), &record)
			if _type != record.Roomid {
				continue
			}
			msg.Info = string(v.([]byte))
		}
	}

	return lib.HF_EncodeMsg("recordhlgc2", &msg, true)
}
