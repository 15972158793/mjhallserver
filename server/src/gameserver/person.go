package gameserver

import (
	"encoding/json"
	"lib"
	"staticfunc"
	"sync"
)

type MapInfo struct {
	Longitude string `json:"longitude"` //! 经度
	Latitude  string `json:"latitude"`  //! 纬度
	Country   string `json:"country"`   //! 国家
	Province  string `json:"province"`  //! 省
	City      string `json:"city"`      //! 市
	Citycode  string `json:"citycode"`  //! 城市编码
	District  string `json:"district"`  //! 区
	Adcode    string `json:"adcode"`    //! 区域码
	Address   string `json:"address"`   //! 地址
}

type UserBase struct {
	Uid   int64 `json:"uid"`   //! uid
	Money int   `json:"money"` //! 金币
	Gem   int   `json:"gem"`   //! 钻石
	Charm int   `json:"charm"` //! 魅力
}

//! 玩家的数据结构
type Person struct {
	Uid      int64  `json:"uid"`
	Name     string `json:"name"`     //! 名字
	Imgurl   string `json:"imgurl"`   //! 头像
	Card     int    `json:"card"`     //! 房卡数量
	Gold     int    `json:"gold"`     //! 金卡数量
	BindGold int    `json:"bindgold"` //! 绑定金币数量
	GameId   int    `json:"gameid"`   //! 当前处于哪个game中
	RoomId   int    `json:"roomid"`   //! 当前处于哪个room中
	Sex      int    `json:"sex"`      //! 性别
	Param    int    `json:"param"`
	UnionId  string `json:"unionid"`
	Admin    int    `json:"admin"`

	session *lib.Session //! session
	ip      string
	minfo   MapInfo //! 地图
	line    bool
	room    *Room
	black   bool //! 是否黑名单了
}

func (self *Person) SendNullMsg() {
	var msg staticfunc.Msg_Null
	self.SendMsg("nothing", &msg)
}

func (self *Person) SendMsg(head string, v interface{}) {
	if self.session == nil {
		return
	}

	//s := GetSessionMgr().GetSession(self.session)
	//if s == nil {
	//	return
	//}

	if v == nil {
		v = &staticfunc.Msg_Null{}
	}

	self.session.SendMsg(head, v)
}

func (self *Person) SendByteMsg(msg []byte) {
	if self.session == nil {
		return
	}

	self.session.SendByteMsg(msg)
}

func (self *Person) CloseSession() {
	if self.session == nil {
		return
	}

	//s := GetSessionMgr().GetSession(self.session)
	//if s == nil {
	//	return
	//}

	self.session.SafeClose()
}

func (self *Person) SendErr(err string) {
	var msg staticfunc.Msg_Err
	msg.Err = err
	self.SendMsg("err", &msg)
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

//! 加入玩家
func (self *PersonMgr) AddPerson(person *Person) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.MapPerson[person.Uid] = person
}

//! 删玩家
func (self *PersonMgr) DelPerson(uid int64) {
	self.lock.Lock()
	defer self.lock.Unlock()

	delete(self.MapPerson, uid)
}

//! 该玩家是否存在
func (self *PersonMgr) GetPerson(uid int64) *Person {
	self.lock.RLock()
	defer self.lock.RUnlock()

	person, ok := self.MapPerson[uid]
	if ok {
		return person
	}

	return nil
}

//! 得到玩家指针
func (self *PersonMgr) ForcePerson(uid int64) *Person {
	self.lock.RLock()
	person, ok := self.MapPerson[uid]
	self.lock.RUnlock()
	if ok {
		return person
	}

	person = new(Person)
	value := GetServer().DB_GetData("user", uid)
	if string(value) != "" {
		json.Unmarshal(value, &person)
	} else { //! redis读不到，换服务器获取
		var _msg staticfunc.Msg_Uid
		_msg.Uid = uid
		result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "getperson", &_msg)
		if err != nil || string(result) == "" {
			return nil
		}
		json.Unmarshal(result, &person)
	}
	return person
}
