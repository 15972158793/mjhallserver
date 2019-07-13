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

type Msg_Sign struct {
	Sign   bool                    `json:"sign"`  //! 今日是否能签到
	Index  int                     `json:"index"` //! 当前已签到第几个
	Config map[int]*Sql_SignConfig `json:"config"`
}

type Msg_SignOK struct {
	Index int `json:"index"`
	Money int `json:"money"`
}

//! 签到模块
type Mod_Sign struct {
	person *Person
	sign   SQL_Sign
}

type SQL_Sign struct {
	Uid   int64 `json:"uid"`
	Index int   `json:"index"` //! 第几次签到
	Time  int64 `json:"time"`  //! 最后一次签到时间

	lib.DataUpdate
}

func (self *Mod_Sign) OnGetData(person *Person) {
	self.person = person
}

func (self *Mod_Sign) OnGetOtherData() {
	c := GetServer().Redis.Get()
	defer c.Close()
	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("sign_%d", self.person.Uid)))
	if err == nil {
		json.Unmarshal(v, &self.sign)
	} else {
		sql := fmt.Sprintf("select * from `sign` where uid = %d", self.person.Uid)
		GetServer().DB.GetOneData(sql, &self.sign)
		if self.sign.Uid <= 0 {
			self.sign.Uid = self.person.Uid
			self.sign.Index = 0
			self.sign.Time = 0
			sql := fmt.Sprintf("insert into `%s`(`uid`, `index`, `time`) values (%d, %d, %d)", "sign", self.person.Uid, 0, 0)
			GetServer().SqlQueue(sql, []byte(""), false)

			c := GetServer().Redis.Get()
			defer c.Close()
			c.Do("SET", fmt.Sprintf("sign_%d", self.person.Uid), lib.HF_JtoB(&self.sign))
			c.Do("EXPIRE", fmt.Sprintf("sign_%d", self.person.Uid), 86400*7)
		}
	}
	self.sign.Init("sign", &self.sign, GetServer().DB)
}

func (self *Mod_Sign) OnMsg(head string, body []byte) bool {
	switch head {
	case "sign": //! 签到
		self.Sign()
		return true
	}
	return false
}

func (self *Mod_Sign) OnSave(sql bool) {
	c := GetServer().Redis.Get()
	defer c.Close()
	c.Do("SET", fmt.Sprintf("sign_%d", self.person.Uid), lib.HF_JtoB(&self.sign))
	c.Do("EXPIRE", fmt.Sprintf("sign_%d", self.person.Uid), 86400*7)

	GetServer().SqlQueue(fmt.Sprintf("update `%s` set `index` = %d, time = %d where `uid` = '%d'", "sign", self.sign.Index, self.sign.Time, self.person.Uid), []byte(""), false)
}

//! 重置签到
func (self *Mod_Sign) InitSign() bool {
	if self.sign.Time == 0 {
		return true
	}

	time1 := time.Unix(self.sign.Time, 0)
	if time1.Year() != time.Now().Year() || time1.Month() != time.Now().Month() || time1.Day() != time.Now().Day() {
		y1, w1 := time1.ISOWeek()
		y, w := time.Now().ISOWeek()
		if y1 != y || w1 != w {
			self.sign.Index = 0
		}
		return true
	}

	return false
}

//! 签到
func (self *Mod_Sign) Sign() {
	if rjmgr.GetRJMgr().CloseSign == 1 {
		self.person.SendErr("该功能暂未开放")
		return
	}

	if !self.InitSign() { //! 今日不能签到
		return
	}

	config := GetSignMgr().Sign[self.sign.Index+1]
	if config == nil {
		return
	}

	self.sign.Index += 1
	self.sign.Time = time.Now().Unix()
	self.OnSave(false)

	ok, card, gold := GetServer().AddGold(self.person.Uid, config.Money, "", -13)
	if !ok {
		return
	}
	self.person.UpdCard(card, gold)

	var msg Msg_SignOK
	msg.Index = self.sign.Index
	msg.Money = config.Money
	self.person.SendMsg(lib.HF_EncodeMsg("signok", &msg, true))
}

func (self *Mod_Sign) SendInfo() {
	var msg Msg_Sign
	msg.Sign = self.InitSign()
	msg.Index = self.sign.Index
	msg.Config = GetSignMgr().Sign
	self.person.SendMsg(lib.HF_EncodeMsg("sign", &msg, true))
}

////////////////////////////////////////////////////////////////////
type Sql_SignConfig struct {
	Id    int `json:"id"`
	Icon  int `json:"icon"`
	Money int `json:"money"`
}

type SignMgr struct {
	Sign map[int]*Sql_SignConfig
}

var signmgrsingleton *SignMgr = nil

//! public
func GetSignMgr() *SignMgr {
	if signmgrsingleton == nil {
		signmgrsingleton = new(SignMgr)
	}

	return signmgrsingleton
}

func (self *SignMgr) GetData() {
	self.Sign = make(map[int]*Sql_SignConfig)
	for key, value := range staticfunc.GetCsvMgr().Data["sign"] {
		node := new(Sql_SignConfig)
		node.Id = lib.HF_Atoi(value["id"])
		node.Icon = lib.HF_Atoi(value["icon"])
		node.Money = lib.HF_Atoi(value["money"])
		self.Sign[key] = node
	}
}
