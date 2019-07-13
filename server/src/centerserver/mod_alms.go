package centerserver

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"lib"
	"rjmgr"
	"staticfunc"
	"time"
)

//! 救助金
type Mod_Alms struct {
	person *Person
	alms   SQL_Alms
}

type SQL_Alms struct {
	Uid  int64 `json:"uid"`
	Num  int   `json:"num"`  //! 领取次数
	Time int64 `json:"time"` //! 免费转的时间

	lib.DataUpdate
}

type Msg_AlmsInfo struct {
	AlmsNum   int `json:"almsnum"`   //! 每天领取救济金次数
	AlmsMoney int `json:"almsmoney"` //! 救济金数量
	AlmsLimit int `json:"almslimit"` //! 少于多少领取
}

func (self *Mod_Alms) OnGetData(person *Person) {
	self.person = person
}

func (self *Mod_Alms) OnGetOtherData() {
	c := GetServer().Redis.Get()
	defer c.Close()

	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("alms_%d", self.person.Uid)))
	if err == nil {
		json.Unmarshal(v, &self.alms)
	} else {
		sql := fmt.Sprintf("select * from `alms` where id = %d ", self.person.Uid)
		GetServer().DB.GetOneData(sql, &self.alms)
		if self.alms.Uid <= 0 {
			self.alms.Uid = self.person.Uid
			self.alms.Num = 0
			self.alms.Time = 0
			sql := fmt.Sprintf("insert into `%s`(`uid`,`num`,`time`) value (%d, %d, %d)", "alms", self.person.Uid, 0, 0)
			GetServer().SqlQueue(sql, []byte(""), false)

			c := GetServer().Redis.Get()
			defer c.Close()
			c.Do("SET", fmt.Sprintf("alms_%d", self.person.Uid), lib.HF_JtoB(&self.alms))
			c.Do("EXPIRE", fmt.Sprintf("alms_%d", self.person.Uid), 86400*7)
		}
	}

	self.alms.Init("alms", &self.alms, GetServer().DB)
}

func (self *Mod_Alms) OnMsg(head string, body []byte) bool {
	switch head {
	case "alms":
		self.GetAlms()
	}
	return false
}

func (self *Mod_Alms) GetAlms() {
	if rjmgr.GetRJMgr().CloseAlms == 1 {
		self.person.SendErr("该功能暂未开放")
		return
	}
	self.InitAlms()

	_, gold, _gold := GetServer().GetCard(self.person.Uid)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "------------- gold : ", gold, " staticfunc.GetModMgr().GetModProperty(GetServer().Redis).AlmsMoney : ", staticfunc.GetModMgr().GetModProperty(GetServer().Redis).AlmsLimit)
	if gold+_gold > staticfunc.GetModMgr().GetModProperty(GetServer().Redis).AlmsLimit {
		self.person.SendErr(fmt.Sprintf("金币低于%d可以领取", staticfunc.GetModMgr().GetModProperty(GetServer().Redis).AlmsLimit))
		return
	}

	if self.alms.Num >= staticfunc.GetModMgr().GetModProperty(GetServer().Redis).AlmsNum {
		self.person.SendErr("领取次数不足")
		return
	}

	self.alms.Time = time.Now().Unix()
	self.alms.Num++

	self.OnSave(false)

	ok, card, gold := GetServer().AddGold(self.person.Uid, staticfunc.GetModMgr().GetModProperty(GetServer().Redis).AlmsMoney, "", -12)
	if !ok {
		return
	}
	self.person.UpdCard(card, gold)

	var msg Msg_AlmsInfo
	msg.AlmsLimit = staticfunc.GetModMgr().GetModProperty(GetServer().Redis).AlmsLimit
	msg.AlmsMoney = staticfunc.GetModMgr().GetModProperty(GetServer().Redis).AlmsMoney
	msg.AlmsNum = staticfunc.GetModMgr().GetModProperty(GetServer().Redis).AlmsNum - self.alms.Num

	self.person.SendMsg(lib.HF_EncodeMsg("alms", &msg, true))
}

//! 重置免费转
func (self *Mod_Alms) InitAlms() {
	if self.alms.Time == 0 {
		return
	}

	time1 := time.Unix(self.alms.Time, 0)
	if time1.Year() != time.Now().Year() || time1.Month() != time.Now().Month() || time1.Day() != time.Now().Day() {
		self.alms.Num = 0
	}
}

func (self *Mod_Alms) OnSave(sql bool) {
	c := GetServer().Redis.Get()
	defer c.Close()
	c.Do("SET", fmt.Sprintf("alms_%d", self.person.Uid), lib.HF_JtoB(&self.alms))
	c.Do("EXPIRE", fmt.Sprintf("alms_%d", self.person.Uid), 86400*7)

	GetServer().SqlQueue(fmt.Sprintf("update `%s` set `num` = %d ,time = %d where `uid` = %d ", "alms", self.alms.Num, self.alms.Time, self.person.Uid), []byte(""), false)
}

func (self *Mod_Alms) SendInfo() {
	self.InitAlms()

	var msg Msg_AlmsInfo
	msg.AlmsLimit = staticfunc.GetModMgr().GetModProperty(GetServer().Redis).AlmsLimit
	msg.AlmsMoney = staticfunc.GetModMgr().GetModProperty(GetServer().Redis).AlmsMoney
	msg.AlmsNum = staticfunc.GetModMgr().GetModProperty(GetServer().Redis).AlmsNum - self.alms.Num

	self.person.SendMsg(lib.HF_EncodeMsg("almsinfo", &msg, true))
}
