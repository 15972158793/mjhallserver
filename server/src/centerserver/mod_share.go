package centerserver

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"lib"
	"staticfunc"
	"time"
)

//! 分享模块
type Mod_Share struct {
	person *Person
	share  SQL_Share
}

type SQL_Share struct {
	Uid  int64 `json:"uid"`
	Num  int   `json:"num"`
	Get  int   `json:"get"`
	Time int64 `json:"time"`

	lib.DataUpdate
}

func (self *Mod_Share) OnGetData(person *Person) {
	self.person = person
}

func (self *Mod_Share) OnGetOtherData() {
	c := GetServer().Redis.Get()
	defer c.Close()
	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("share_%d", self.person.Uid)))
	if err == nil {
		json.Unmarshal(v, &self.share)
	} else {
		sql := fmt.Sprintf("select * from `share` where uid = %d", self.person.Uid)
		GetServer().DB.GetOneData(sql, &self.share)
		if self.share.Uid <= 0 {
			self.share.Uid = self.person.Uid
			self.share.Num = 0
			self.share.Time = 0
			sql := fmt.Sprintf("insert into `%s`(`uid`, `num`, `get`, `time`) values (%d, %d, %d, %d)", "share", self.person.Uid, 0, 0, 0)
			GetServer().SqlQueue(sql, []byte(""), false)

			c := GetServer().Redis.Get()
			defer c.Close()
			c.Do("SET", fmt.Sprintf("share_%d", self.person.Uid), lib.HF_JtoB(&self.share))
			c.Do("EXPIRE", fmt.Sprintf("share_%d", self.person.Uid), 86400*7)
		}
	}
	self.share.Init("share", &self.share, GetServer().DB)
}

func (self *Mod_Share) OnMsg(head string, body []byte) bool {
	switch head {
	case "shareOK":
		var msg Msg_ReadMail
		json.Unmarshal(body, &msg)
		self.Share(msg.Id)
		return true
	case "shareGet":
		self.ShareGet()
		return true
	}
	return false
}

func (self *Mod_Share) OnSave(sql bool) {
	c := GetServer().Redis.Get()
	defer c.Close()
	c.Do("SET", fmt.Sprintf("share_%d", self.person.Uid), lib.HF_JtoB(&self.share))
	c.Do("EXPIRE", fmt.Sprintf("share_%d", self.person.Uid), 86400*7)

	GetServer().SqlQueue(fmt.Sprintf("update `%s` set `num` = %d, time = %d, `get` = %d where `uid` = '%d'", "share", self.share.Num, self.share.Time, self.share.Get, self.person.Uid), []byte(""), false)
}

//! 重置分享
func (self *Mod_Share) InitShare() {
	if self.share.Time == 0 {
		return
	}

	time1 := time.Unix(self.share.Time, 0)
	if time1.Year() != time.Now().Year() || time1.Month() != time.Now().Month() || time1.Day() != time.Now().Day() {
		self.share.Num = 0
	}
}

//! 分享
func (self *Mod_Share) Share(id int) {
	self.InitShare()

	csv, _ := staticfunc.GetCsvMgr().Data["game"][id]
	sharenum := lib.HF_Atoi(csv["sharenum"])
	if sharenum == 0 {
		sharenum = 1
	}
	sharecard := lib.HF_Atoi(csv["sharecard"])
	if sharecard == 0 {
		sharecard = 3
	}

	if self.share.Num >= sharenum {
		return
	}

	self.share.Num++
	self.share.Get += sharecard
	self.share.Time = time.Now().Unix()
	self.OnSave(false)

	if lib.HF_Atoi(csv["sharetype"]) == 0 { //! 立即送卡
		ok, card, gold := GetServer().AddCard(self.person.Uid, self.share.Get, "", -3)
		if !ok {
			return
		}
		self.person.UpdCard(card, gold)

		self.share.Get = 0
		self.OnSave(false)
	}

	self.person.SendMsg(lib.HF_EncodeMsg("share", &self.share, true))
}

//! 分享领取
func (self *Mod_Share) ShareGet() {
	if self.share.Get == 0 {
		return
	}

	ok, card, gold := GetServer().AddCard(self.person.Uid, self.share.Get, "", -3)
	if !ok {
		return
	}
	self.person.UpdCard(card, gold)

	self.share.Get = 0
	self.OnSave(false)

	self.person.SendMsg(lib.HF_EncodeMsg("share", &self.share, true))
}

func (self *Mod_Share) SendInfo() {
	self.InitShare()
	self.person.SendMsg(lib.HF_EncodeMsg("share", &self.share, true))
}
