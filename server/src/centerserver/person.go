package centerserver

import (
	"encoding/json"
	"lib"
	"staticfunc"
	"sync"
	"time"
)

////////////////////////////////////////////////////////////////////
//! 玩家的数据结构
type BaseInfo struct {
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`   //! 名字
	Imgurl string `json:"imgurl"` //! 头像
	Sex    int    `json:"sex"`    //! 性别
}

//! 玩家的数据结构
type Person struct {
	Uid       int64
	session   *lib.Session
	group     int
	mapinfo   MapInfo
	Module    *Mod_All
	OtherData bool
	MsgTime   int64
	mod_info  BaseInfo //! 基础信息
}

//! 得到person
func NewPerson(uid int64) *Person {
	p := new(Person)
	p.Module = NewModAll(p)
	p.Uid = uid
	p.MsgTime = time.Now().Unix()
	return p
}

//! 保存
//func (self *Person) Save() {
//	self.Module.Save(true)
//}

//! 得到基础信息
func (self *Person) InitPlayerData() {
	self.MsgTime = time.Now().Unix()
	self.Module.GetData()

	go self.RunTime()
}

//! 得到其他信息
func (self *Person) OtherPlayerData() {
	if self.OtherData {
		return
	}
	self.Module.GetOtherData()
	self.OtherData = true
}

//! 得到模块
func (self *Person) GetModule(name string) Mod_Base {
	return self.Module.GetModule(name)
}

//! 定时器
func (self *Person) RunTime() {
	ticker := time.NewTicker(time.Second * 600)
	for {
		<-ticker.C
		if time.Now().Unix()-self.MsgTime >= 600 && self.session == nil {
			break
		}
	}

	//！ 关掉定时器
	ticker.Stop()
	GetPersonMgr().DelPerson(self.Uid)
}

func (self *Person) SendMsg(data []byte) {
	if self.session == nil {
		return
	}

	self.session.SendByteMsg(data)
}

func (self *Person) CloseSession() {
	if self.session == nil {
		return
	}

	self.session.SafeClose()
}

func (self *Person) GetBaseInfo() {
	value := GetServer().DB_GetData("user", self.Uid)
	if string(value) != "" {
		json.Unmarshal(value, &self.mod_info)
	} else { //! redis读不到，换服务器获取
		var _msg staticfunc.Msg_Uid
		_msg.Uid = self.Uid
		result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "getperson", &_msg)
		if err != nil || string(result) == "" {
			return
		}
		json.Unmarshal(result, &self.mod_info)
	}
}

func (self *Person) SendErr(err string) {
	var msg staticfunc.Msg_Err
	msg.Err = err
	self.SendMsg(lib.HF_EncodeMsg("err", &msg, true))
}

func (self *Person) SendWarning(err string) {
	var msg staticfunc.Msg_Err
	msg.Err = err
	self.SendMsg(lib.HF_EncodeMsg("warning", &msg, true))
}

func (self *Person) SendRet(head string) {
	var msg staticfunc.Msg_Null
	self.SendMsg(lib.HF_EncodeMsg(head, &msg, true))
}

//! 同步房卡
//! card 房卡+gold总和
func (self *Person) UpdCard(card int, gold int) {
	var msg staticfunc.S2C_UpdCard
	msg.Card = card + gold
	msg.Gold = gold
	self.SendMsg(lib.HF_EncodeMsg("updcard", &msg, true))
}

//////////////////////////////////////////////////////////////
//! 玩家管理者
type PersonMgr struct {
	MapPerson map[int64]*Person

	lock *sync.RWMutex
}

var personmgrSingleton *PersonMgr = nil

//! 得到服务器指针
func GetPersonMgr() *PersonMgr {
	if personmgrSingleton == nil {
		personmgrSingleton = new(PersonMgr)
		personmgrSingleton.MapPerson = make(map[int64]*Person)
		personmgrSingleton.lock = new(sync.RWMutex)
	}

	return personmgrSingleton
}

//! 删玩家
func (self *PersonMgr) DelPerson(uid int64) {
	self.lock.Lock()
	defer self.lock.Unlock()

	delete(self.MapPerson, uid)
	//GetServer().Wait.Done()
}

//! 该玩家是否存在
func (self *PersonMgr) GetPerson(uid int64, add bool) *Person {
	self.lock.RLock()
	defer self.lock.RUnlock()

	person, ok := self.MapPerson[uid]
	if ok {
		person.MsgTime = time.Now().Unix()
		return person
	}

	if add {
		person = NewPerson(uid)
		person.GetBaseInfo()
		person.InitPlayerData()
		self.MapPerson[person.Uid] = person
		//GetServer().Wait.Add(1)
		return person
	}

	return nil
}

func (self *PersonMgr) SaveAll() {
	//for _, _ = range self.MapPerson {
	//	GetServer().Wait.Done()
	//}
}

//! 广播消息
func (self *PersonMgr) BroadCastMsg(area string, body []byte, group int) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	for _, value := range self.MapPerson {
		if value.group != group {
			continue
		}

		if value.mapinfo.Adcode == area {
			value.SendMsg(body)
		} else if area == "" && GetServer().AreaNotice[group-1][value.mapinfo.Adcode] == "" {
			value.SendMsg(body)
		}
	}
}

//! 广播消息
func (self *PersonMgr) BroadCastMsg2(body []byte) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	for _, value := range self.MapPerson {
		value.SendMsg(body)
	}
}

func (self *PersonMgr) GetNum(group int) int {
	self.lock.RLock()
	defer self.lock.RUnlock()

	num := 0
	for _, value := range self.MapPerson {
		if value.group != group {
			continue
		}
		if value.session == nil {
			continue
		}

		num++
	}

	return num
}

//! 得到在线玩家
func (self *PersonMgr) GetOnline() []int64 {
	self.lock.RLock()
	defer self.lock.RUnlock()

	lst := make([]int64, 0)
	for key, value := range self.MapPerson {
		if value.session == nil {
			continue
		}
		lst = append(lst, key)
	}
	return lst
}
