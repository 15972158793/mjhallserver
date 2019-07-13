package centerserver

import (
//"encoding/json"
//"fmt"
////"log"
//"staticfunc"

//"github.com/garyburd/redigo/redis"
)

//type SQL_Friend struct {
//	Uid   int64
//	info  JS_Friend
//	Value []byte
//}

//type JS_Friend struct {
//	Friend []Son_Friend `json:"friend"`
//	Ask    []Son_Friend `json:"ask"`
//}

//type Son_Friend struct {
//	Uid    int64  `json:"uid"`
//	Name   string `json:"name"`
//	Imgurl string `json:"imgurl"` //! 头像
//	Online bool   `json:"online"` //! 是否在线
//}

//func (self *SQL_Friend) Decode() {
//	json.Unmarshal(self.Value, &self.info)
//}

//func (self *SQL_Friend) Encode() {
//	self.Value = staticfunc.HF_JtoB(self.info)
//}

//type Mod_Friend struct {
//	person *Person
//	friend SQL_Friend
//}

//func (self *Mod_Friend) OnGetData(person *Person) {
//	self.person = person

//	c := GetServer().Redis.Get()
//	defer c.Close()
//	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("friend_%d", self.person.Uid)))
//	if err == nil {
//		self.friend.Value = v
//		self.friend.Decode()
//	} else {
//		sql := fmt.Sprintf("select * from `friend` where uid = %d", self.person.Uid)
//		GetServer().DB.GetOneData(sql, &self.friend)
//		if self.friend.Uid <= 0 {
//			self.friend.Uid = self.person.Uid
//			self.friend.info.Friend = make([]Son_Friend, 0)
//			self.friend.info.Ask = make([]Son_Friend, 0)
//			self.friend.Encode()
//			sql := fmt.Sprintf("insert into `%s`(`uid`, `value`) values (%d, ?)", "friend", self.person.Uid)
//			GetServer().SqlQueue(sql, self.friend.Value, true)

//			c := GetServer().Redis.Get()
//			defer c.Close()
//			c.Do("SET", fmt.Sprintf("friend_%d", self.person.Uid), self.friend.Value)
//		} else {
//			self.friend.Decode()
//		}
//	}
//}

//func (self *Mod_Friend) OnGetOtherData() {
//}

//func (self *Mod_Friend) OnSave(sql bool) {
//	self.friend.Encode()
//	c := GetServer().Redis.Get()
//	defer c.Close()
//	c.Do("SET", fmt.Sprintf("friend_%d", self.person.Uid), self.friend.Value)

//	GetServer().SqlQueue(fmt.Sprintf("update `%s` set `value` = ? where `uid` = '%d'", "friend", self.person.Uid), self.friend.Value, true)
//}

//func (self *Mod_Friend) OnMsg(head string, body []byte) bool {
//	switch head {
//	case "friendGet": //！获取好友列表
//		self.SendInfo()
//		return true
//	case "friendDel": //！删除好友
//		var msg Msg_FriendUid
//		json.Unmarshal(body, &msg)
//		self.FriendDel(msg.Uid)
//		return true
//	case "friendAdd": //！确认添加好友
//		var msg Msg_FriendUid
//		json.Unmarshal(body, &msg)
//		self.FriendAdd(msg.Uid)
//		return true
//	case "friendSearch": //！请求添加好友
//		var msg Msg_FriendUid
//		json.Unmarshal(body, &msg)
//		self.FriendSearch(msg.Uid)
//		return true
//	case "friendAsk": //！请求添加好友
//		var msg Msg_FriendUid
//		json.Unmarshal(body, &msg)
//		self.FriendAsk(msg.Uid)
//		return true
//	case "friendRefuse": //！拒绝好友请求
//		var msg Msg_FriendUid
//		json.Unmarshal(body, &msg)
//		self.FriendRefuse(msg.Uid)
//		return true
//	}

//	return false
//}

//func (self *Mod_Friend) SendInfo() {
//	self.person.SendMsg(staticfunc.HF_EncodeMsg("friend", &self.friend.info, true))
//}

//func (self *Mod_Friend) FriendDel(uid int64) {
//	if self.person.Uid == uid {
//		//////////////////////////////////////////
//		self.person.SendErr("不能删除自己")
//		return
//	}
//	delPerson := GetPersonMgr().GetPerson(uid, true)
//	if nil == delPerson {
//		//////////////////////////////////////////
//		self.person.SendErr("删除玩家不存在")
//		return
//	}
//	for i := 0; i < len(self.friend.info.Friend); i++ { //！删除自己列表
//		if self.friend.info.Friend[i].Uid == uid {
//			copy(self.friend.info.Friend[i:], self.friend.info.Friend[i+1:])
//			self.friend.info.Friend = self.friend.info.Friend[:len(self.friend.info.Friend)-1]
//			break
//		}
//	}
//	mod_friend := delPerson.GetModule("friend").(*Mod_Friend)
//	for i := 0; i < len(mod_friend.friend.info.Friend); i++ { //！删除对方好友列表
//		if mod_friend.friend.info.Friend[i].Uid == self.person.Uid {
//			copy(mod_friend.friend.info.Friend[i:], mod_friend.friend.info.Friend[i+1:])
//			mod_friend.friend.info.Friend = mod_friend.friend.info.Friend[:len(mod_friend.friend.info.Friend)-1]
//			break
//		}
//	}
//	var msg Msg_FriendUid
//	msg.Uid = uid
//	var _msg Msg_FriendUid
//	_msg.Uid = self.person.Uid
//	mod_friend.OnSave(false)
//	self.OnSave(false)
//	if delPerson.session != nil {
//		delPerson.SendMsg(staticfunc.HF_EncodeMsg("frienddel", &_msg, true))
//	}
//	self.person.SendMsg(staticfunc.HF_EncodeMsg("frienddel", &msg, true))
//}

//func (self *Mod_Friend) FriendAdd(uid int64) {
//	if self.person.Uid == uid {
//		//////////////////////////////////////////
//		self.person.SendErr("不能添加自己")
//		return
//	}
//	delPerson := GetPersonMgr().GetPerson(uid, true)
//	if nil == delPerson {
//		//////////////////////////////////////////
//		self.person.SendErr("添加玩家不存在")
//		return
//	}
//	for i := 0; i < len(self.friend.info.Friend); i++ {
//		if self.friend.info.Friend[i].Uid == uid {
//			////////////////////////////////////////////
//			self.person.SendErr("该玩家已经在好友列表中")
//			return
//		}
//	}

//	find := false
//	var son Son_Friend
//	for i := 0; i < len(self.friend.info.Ask); i++ {
//		if self.friend.info.Ask[i].Uid == uid {
//			find = true
//			son = self.friend.info.Ask[i]
//			self.friend.info.Friend = append(self.friend.info.Friend, son)
//			//delPerson.mod_friend.info.Friend = append(delPerson.mod_friend.info.Friend, son)
//			copy(self.friend.info.Ask[i:], self.friend.info.Ask[i+1:])
//			self.friend.info.Ask = self.friend.info.Ask[:len(self.friend.info.Ask)-1]
//			break
//		}
//	}
//	if !find {
//		//////////////////////////////////////////
//		self.person.SendErr("添加玩家不在请求列表中")
//		return
//	}
//	var son1 Son_Friend
//	son1.Uid = self.person.mod_info.Uid
//	son1.Imgurl = self.person.mod_info.Imgurl
//	son1.Name = self.person.mod_info.Name

//	mod_friend := delPerson.GetModule("friend").(*Mod_Friend)
//	mod_friend.friend.info.Friend = append(mod_friend.friend.info.Friend, son1)

//	mod_friend.OnSave(false)
//	self.OnSave(false)
//	if delPerson.session != nil {
//		delPerson.SendMsg(staticfunc.HF_EncodeMsg("friendadd", &son1, true))
//	}
//	self.person.SendMsg(staticfunc.HF_EncodeMsg("friendadd", &son, true))
//}

//func (self *Mod_Friend) FriendSearch(uid int64) {
//	delPerson := GetPersonMgr().GetPerson(uid, true)
//	if nil == delPerson {
//		//////////////////////////////////////////
//		self.person.SendErr("查找的玩家不存在")
//		return
//	}
//	var msg Son_Friend
//	msg.Uid = delPerson.mod_info.Uid
//	msg.Imgurl = delPerson.mod_info.Imgurl
//	msg.Name = delPerson.mod_info.Name
//	self.person.SendMsg(staticfunc.HF_EncodeMsg("friendsearch", &msg, true))
//}

//func (self *Mod_Friend) FriendAsk(uid int64) {
//	if self.person.Uid == uid {
//		//////////////////////////////////////////
//		self.person.SendErr("不能请求自己")
//		return
//	}
//	delPerson := GetPersonMgr().GetPerson(uid, true)
//	if nil == delPerson {
//		//////////////////////////////////////////
//		self.person.SendErr("添加玩家不存在")
//		return
//	}
//	for i := 0; i < len(self.friend.info.Friend); i++ {
//		if self.friend.info.Friend[i].Uid == uid {
//			//////////////////////////////////////////
//			self.person.SendErr("该玩家已在好友列表中")
//			return
//		}
//	}

//	mod_friend := delPerson.GetModule("friend").(*Mod_Friend)
//	find := false
//	for i := 0; i < len(mod_friend.friend.info.Ask); i++ {
//		if mod_friend.friend.info.Ask[i].Uid == self.person.Uid {
//			//////////////////////////////////////////
//			//			self.SendErr("请勿重复申请")
//			//			return
//			find = true
//			break
//		}
//	}
//	var son Son_Friend
//	son.Uid = self.person.mod_info.Uid
//	son.Name = self.person.mod_info.Name
//	son.Imgurl = self.person.mod_info.Imgurl

//	if !find {
//		mod_friend.friend.info.Ask = append(mod_friend.friend.info.Ask, son)
//		mod_friend.OnSave(false)
//	}

//	if delPerson.session != nil {
//		delPerson.SendMsg(staticfunc.HF_EncodeMsg("friendasked", &son, true))
//	}

//	self.person.SendMsg(staticfunc.HF_EncodeMsg("friendask", nil, true))
//}

//func (self *Mod_Friend) FriendRefuse(uid int64) {
//	if self.person.Uid == uid {
//		//////////////////////////////////////////
//		self.person.SendErr("不能拒绝自己")
//		return
//	}
//	delPerson := GetPersonMgr().GetPerson(uid, true)
//	if nil == delPerson {
//		//////////////////////////////////////////
//		self.person.SendErr("拒绝玩家不存在")
//		return
//	}
//	find := false
//	for i := 0; i < len(self.friend.info.Ask); i++ {
//		if self.friend.info.Ask[i].Uid == uid {
//			find = true
//			//son = self.mod_friend.info.Ask[i]
//			//self.mod_friend.info.Friend = append(self.mod_friend.info.Friend, son)
//			copy(self.friend.info.Ask[i:], self.friend.info.Ask[i+1:])
//			self.friend.info.Ask = self.friend.info.Ask[:len(self.friend.info.Ask)-1]
//			break
//		}
//	}
//	if !find {
//		self.person.SendErr("拒绝玩家不在请求列表中")
//		return
//	}
//	var msg Msg_FriendUid
//	msg.Uid = uid
//	self.OnSave(false)
//	self.person.SendMsg(staticfunc.HF_EncodeMsg("friendrefuse", &msg, true))
//}
