package centerserver

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"lib"
	"time"
)

//! 邀请记录
type Son_Invite struct {
	Name string `json:"name"`
	Time int64  `json:"time"`
}

type JS_Invite struct {
	IsInvite int          `json:"isinvite"`
	Invite   []Son_Invite `json:"invite"`
}

//! 邀请结构
type SQL_Invite struct {
	Uid   int64
	Value []byte

	info JS_Invite
}

func (self *SQL_Invite) Decode() {
	json.Unmarshal(self.Value, &self.info)
}

func (self *SQL_Invite) Encode() {
	self.Value = lib.HF_JtoB(self.info)
}

//! 邀请模块
type Mod_Invite struct {
	person *Person
	invite SQL_Invite
}

func (self *Mod_Invite) OnGetData(person *Person) {
	self.person = person

	c := GetServer().Redis.Get()
	defer c.Close()
	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("invite_%d", self.person.Uid)))
	if err == nil {
		self.invite.Value = v
		self.invite.Decode()
	} else {
		sql := fmt.Sprintf("select * from `invite` where uid = %d", self.person.Uid)
		GetServer().DB.GetOneData(sql, &self.invite)
		if self.invite.Uid <= 0 {
			self.invite.Uid = self.person.Uid
			self.invite.info.IsInvite = 0
			self.invite.info.Invite = make([]Son_Invite, 0)
			self.invite.Encode()
			sql := fmt.Sprintf("insert into `%s`(`uid`, `value`) values (%d, ?)", "invite", self.person.Uid)
			GetServer().SqlQueue(sql, self.invite.Value, true)

			c := GetServer().Redis.Get()
			defer c.Close()
			c.Do("SET", fmt.Sprintf("invite_%d", self.person.Uid), self.invite.Value)
			c.Do("EXPIRE", fmt.Sprintf("invite_%d", self.person.Uid), 86400*7)
		} else {
			self.invite.Decode()
		}
	}
}

func (self *Mod_Invite) OnGetOtherData() {
}

func (self *Mod_Invite) OnMsg(head string, body []byte) bool {
	switch head {
	case "invite": //! 使用邀请码
		var msg C2S_InviteCode
		json.Unmarshal(body, &msg)
		self.Invite(msg.Code)
	}
	return false
}

func (self *Mod_Invite) OnSave(sql bool) {
	self.invite.Encode()
	c := GetServer().Redis.Get()
	defer c.Close()
	c.Do("SET", fmt.Sprintf("invite_%d", self.person.Uid), self.invite.Value)
	c.Do("EXPIRE", fmt.Sprintf("invite_%d", self.person.Uid), 86400*7)

	GetServer().SqlQueue(fmt.Sprintf("update `%s` set `value` = ? where `uid` = '%d'", "invite", self.person.Uid), self.invite.Value, true)
}

//! 使用激活码
func (self *Mod_Invite) Invite(code int) {
	if self.person.Uid == int64(code) {
		self.person.SendErr("不能填写自己的激活码")
		return
	}

	if self.invite.info.IsInvite != 0 {
		self.person.SendErr("您已使用过激活码")
		return
	}

	//! 对方加卡
	host := GetPersonMgr().GetPerson(int64(code), true)
	if host == nil {
		self.person.SendErr("该激活码无效")
		return
	}

	ok, card, gold := GetServer().AddCard(int64(code), 5, "", -2)
	if !ok {
		self.person.SendErr("使用失败")
		return
	}

	host.GetModule("invite").(*Mod_Invite).invite.info.Invite = append(host.GetModule("invite").(*Mod_Invite).invite.info.Invite, Son_Invite{self.person.mod_info.Name, time.Now().Unix()})
	host.GetModule("invite").(*Mod_Invite).OnSave(false)

	if host.session != nil {
		host.UpdCard(card, gold)
		host.GetModule("invite").(*Mod_Invite).SendInfo()
	}

	//! 自己加卡
	ok, card, gold = GetServer().AddCard(self.person.Uid, 5, "", -2)
	if !ok {
		return
	}

	self.invite.info.IsInvite = 1
	self.OnSave(false)
	self.person.UpdCard(card, gold)
	self.SendInfo()
}

func (self *Mod_Invite) SendInfo() {
	self.person.SendMsg(lib.HF_EncodeMsg("invite", &self.invite.info, true))
}
