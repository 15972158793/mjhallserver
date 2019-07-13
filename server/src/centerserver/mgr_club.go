package centerserver

import (
	"encoding/json"
	"fmt"
	"lib"
	"staticfunc"
	"sync"
	"time"
)

const CLUBMEM_HOST = 0 //! 主席
const CLUBMEM_MEM = 1  //! 成员

type ChgType struct {
	Chg bool
}

type SQL_ClubMgr struct {
	Id   int64
	Info []byte

	info       JS_ClubMgr
	chg        ChgType
	lock       *sync.RWMutex //! 俱乐部成员和申请操作锁
	roomlock   *sync.RWMutex //! 俱乐部房间锁
	chatlock   *sync.RWMutex //! 聊天锁
	resultlock *sync.RWMutex //! 战绩锁
	eventlock  *sync.RWMutex //! 事件锁
	costlock   *sync.RWMutex //! 消耗锁
}

func (self *SQL_ClubMgr) Decode() {
	json.Unmarshal(self.Info, &self.info)
}

func (self *SQL_ClubMgr) Encode() {
	self.Info = lib.HF_JtoB(self.info)
}

func (self *SQL_ClubMgr) Save() {
	if !self.chg.Chg {
		return
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	self.Encode()
	sql := fmt.Sprintf("update `%s` set `info` = ? where `id` = '%d'", "clubmgr", self.Id)
	GetServer().DB.DB.Exec(sql, self.Info)
	self.chg.Chg = false
}

//! 加入事件
func (self *SQL_ClubMgr) AddEvent(uid int64, name string, head string, event string) {
	self.eventlock.Lock()
	defer self.eventlock.Unlock()

	//! 先清理过期的事件
	for i := 0; i < len(self.info.Event); {
		if time.Now().Unix()-self.info.Event[i].Time >= 24*3600 {
			copy(self.info.Event[i:], self.info.Event[i+1:])
			self.info.Event = self.info.Event[:len(self.info.Event)-1]
		} else {
			i++
		}
	}

	self.info.Event = append(self.info.Event, JS_ClubEvent{uid, name, head, time.Now().Unix(), event})
}

//! 得到事件
func (self *SQL_ClubMgr) GetEvent() []JS_ClubEvent {
	self.eventlock.RLock()
	defer self.eventlock.RUnlock()

	return self.info.Event
}

//! 加入消耗
func (self *SQL_ClubMgr) AddCostCard(num int) {
	self.costlock.Lock()
	defer self.costlock.Unlock()

	find := false
	for i := 0; i < len(self.info.CostCard); i++ {
		if self.info.CostCard[i].Year == time.Now().Year() && self.info.CostCard[i].Month == int(time.Now().Month()) && self.info.CostCard[i].Day == time.Now().Day() {
			self.info.CostCard[i].Num += num
			find = true
			break
		}
	}

	if !find {
		self.info.CostCard = append(self.info.CostCard, JS_ClubCostCard{time.Now().Year(), int(time.Now().Month()), time.Now().Day(), num})
		if len(self.info.CostCard) > 7 {
			self.info.CostCard = self.info.CostCard[len(self.info.CostCard)-7 : len(self.info.CostCard)]
		}
	}
}

//! 得到消耗
func (self *SQL_ClubMgr) GetCostCard() []JS_ClubCostCard {
	self.costlock.RLock()
	defer self.costlock.RUnlock()

	return self.info.CostCard
}

func (self *SQL_ClubMgr) HasMem(uid int64) bool {
	self.lock.RLock()
	defer self.lock.RUnlock()

	for i := 0; i < len(self.info.Mem); i++ {
		if self.info.Mem[i].Uid == uid {
			return true
		}
	}

	return false
}

//! 更新成员信息
func (self *SQL_ClubMgr) UpdMem(uid int64, name string, head string, online bool) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	for i := 0; i < len(self.info.Mem); i++ {
		if self.info.Mem[i].Uid == uid {
			self.info.Mem[i].Name = name
			self.info.Mem[i].Head = head
			self.info.Mem[i].online = online
			break
		}
	}
}

func (self *SQL_ClubMgr) HasApply(uid int64) bool {
	self.lock.RLock()
	defer self.lock.RUnlock()

	for i := 0; i < len(self.info.Apply); i++ {
		if self.info.Apply[i].Uid == uid {
			return true
		}
	}

	return false
}

//! 加入申请
func (self *SQL_ClubMgr) AddApply(uid int64, name string, head string) bool {
	self.lock.Lock()
	defer self.lock.Unlock()

	if len(self.info.Apply) >= 100 {
		return false
	}

	self.info.Apply = append(self.info.Apply, JS_ClubApply{uid, name, head, time.Now().Unix()})

	self.chg.Chg = true

	return true
}

//! 撤销申请
func (self *SQL_ClubMgr) DelApply(uid int64) {
	self.lock.Lock()
	defer self.lock.Unlock()

	for i := 0; i < len(self.info.Apply); i++ {
		if self.info.Apply[i].Uid == uid {
			copy(self.info.Apply[i:], self.info.Apply[i+1:])
			self.info.Apply = self.info.Apply[:len(self.info.Apply)-1]
			break
		}
	}

	self.chg.Chg = true
}

//! 处理申请
func (self *SQL_ClubMgr) OrderApply(uid int64, agree bool) bool {
	self.lock.Lock()
	defer self.lock.Unlock()

	find := false
	for i := 0; i < len(self.info.Apply); i++ {
		if self.info.Apply[i].Uid == uid {
			if agree { //! 同意
				self.info.Mem = append(self.info.Mem, JS_ClubMem{self.info.Apply[i].Uid, self.info.Apply[i].Name, self.info.Apply[i].Head, CLUBMEM_MEM, false})
				self.AddEvent(self.info.Apply[i].Uid, self.info.Apply[i].Name, self.info.Apply[i].Head, fmt.Sprintf("%s加入俱乐部", self.info.Apply[i].Name))
			}
			copy(self.info.Apply[i:], self.info.Apply[i+1:])
			self.info.Apply = self.info.Apply[:len(self.info.Apply)-1]
			find = true
			break
		}
	}

	if !find {
		return false
	}

	self.chg.Chg = true
	return true
}

//! 离开
func (self *SQL_ClubMgr) Leave(uid int64) bool {
	self.lock.Lock()
	defer self.lock.Unlock()

	for i := 0; i < len(self.info.Mem); i++ {
		if self.info.Mem[i].Uid == uid {
			self.AddEvent(self.info.Mem[i].Uid, self.info.Mem[i].Name, self.info.Mem[i].Head, fmt.Sprintf("%s离开俱乐部", self.info.Mem[i].Name))
			copy(self.info.Mem[i:], self.info.Mem[i+1:])
			self.info.Mem = self.info.Mem[:len(self.info.Mem)-1]
			self.chg.Chg = true
			return true
		}
	}

	return false
}

//! 加入聊天记录
func (self *SQL_ClubMgr) AddRoomChat(chat JS_ClubRoomChat) {
	self.chatlock.Lock()
	defer self.chatlock.Unlock()

	self.info.RoomChat = append(self.info.RoomChat, chat)
	if len(self.info.RoomChat) > 200 {
		self.info.RoomChat = self.info.RoomChat[len(self.info.RoomChat)-200:]
	}
}

//! 删除聊天记录
func (self *SQL_ClubMgr) DelRoomChat(roomid int) {
	self.chatlock.Lock()
	defer self.chatlock.Unlock()

	for i := 0; i < len(self.info.RoomChat); i++ {
		if self.info.RoomChat[i].RoomId == roomid {
			copy(self.info.RoomChat[i:], self.info.RoomChat[i+1:])
			self.info.RoomChat = self.info.RoomChat[:len(self.info.RoomChat)-1]
			break
		}
	}
}

//! 设置记录
func (self *SQL_ClubMgr) StateRoomChat(roomid int) {
	self.chatlock.Lock()
	defer self.chatlock.Unlock()

	for i := 0; i < len(self.info.RoomChat); i++ {
		if self.info.RoomChat[i].RoomId == roomid {
			self.info.RoomChat[i].State = 1
			self.AddCostCard(self.info.RoomChat[i].Num)
			break
		}
	}
}

//! 战绩
func (self *SQL_ClubMgr) AddRoomResult(result JS_ClubRoomResult) {
	self.resultlock.Lock()
	defer self.resultlock.Unlock()

	//! 先清理过期的战绩
	for i := 0; i < len(self.info.RoomResult); {
		if time.Now().Unix()-self.info.RoomResult[i].Time >= 24*3600 {
			copy(self.info.RoomResult[i:], self.info.RoomResult[i+1:])
			self.info.RoomResult = self.info.RoomResult[:len(self.info.RoomResult)-1]
		} else {
			i++
		}
	}

	self.info.RoomResult = append(self.info.RoomResult, result)

	self.chg.Chg = true
}

//! 修改名字
func (self *SQL_ClubMgr) SetName(name string) bool {
	if GetClubMgr().GetClubFromName(name) != nil {
		return false
	}

	self.info.Name = name
	self.chg.Chg = true

	return true
}

//! 广播消息
func (self *SQL_ClubMgr) BroadCastMsg(msg []byte, all bool) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	for i := 0; i < len(self.info.Mem); i++ {
		if !self.info.Mem[i].online && !all {
			continue
		}

		person := GetPersonMgr().GetPerson(self.info.Mem[i].Uid, false)
		if person == nil {
			continue
		}
		person.SendMsg(msg)
	}
}

//! 俱乐部成员
type JS_ClubMem struct {
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`
	Head   string `json:"head"`
	Job    int    `json:"job"`
	online bool
}

//! 俱乐部申请
type JS_ClubApply struct {
	Uid  int64  `json:"uid"`
	Name string `json:"name"`
	Head string `json:"head"`
	Time int64  `json:"time"`
}

//! 房间聊天记录
type JS_ClubRoomChat struct {
	Uid      int64  `json:"uid"` //! 开房的人
	Name     string `json:"name"`
	Head     string `json:"head"`
	RoomId   int    `json:"roomid"`   //! 房间id
	GameType int    `json:"gametype"` //! 游戏type
	Param1   int    `json:"param1"`   //! 玩法
	Param2   int    `json:"param2"`
	MaxStep  int    `json:"maxstep"` //! 最大局数
	Num      int    `json:"num"`     //! 消耗房卡
	Time     int64  `json:"time"`    //! 开房时间
	State    int    `json:"state"`   //! 0未开始  1已开始
}

//! 房间战绩
type JS_ClubRoomResult struct {
	RoomId   int                           `json:"roomid"`
	GameType int                           `json:"gametype"` //! 游戏type
	Param1   int                           `json:"param1"`   //! 玩法
	Param2   int                           `json:"param2"`
	MaxStep  int                           `json:"maxstep"` //! 最大局数
	Num      int                           `json:"num"`     //! 消耗几张房卡
	Info     []staticfunc.JS_CreateRoomMem `json:"info"`
	Time     int64                         `json:"time"` //!
}

//! 俱乐部事件
type JS_ClubEvent struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Time  int64  `json:"time"`
	Event string `json:"event"`
}

//! 俱乐部消耗房卡
type JS_ClubCostCard struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
	Num   int `json:"num"`
}

//! 俱乐部游戏
type JS_ClubGame struct {
	GameType int `json:"gametype"`
	Param1   int `json:"param1"`
	Param2   int `json:"param2"`
	Num      int `json:"num"`
}

type JS_ClubMgr struct {
	Name       string              `json:"name"`       //! 俱乐部名字
	Icon       string              `json:"icon"`       //! 俱乐部图标
	Host       int64               `json:"host"`       //! 俱乐部主席
	InNotice   string              `json:"innotice"`   //! 俱乐部内部公告
	ExNotice   string              `json:"exnotice"`   //! 俱乐部外部公告
	Mode       int                 `json:"mode"`       //! 俱乐部开房模式 0只有主席能开 1任意成员能开
	Mem        []JS_ClubMem        `json:"mem"`        //! 俱乐部成员
	Apply      []JS_ClubApply      `json:"apply"`      //! 俱乐部申请
	RoomChat   []JS_ClubRoomChat   `jsom:"roomchat"`   //! 开房聊天记录
	RoomResult []JS_ClubRoomResult `json:"roomresult"` //! 房间战绩
	Event      []JS_ClubEvent      `json:"event"`      //! 俱乐部事件
	CostCard   []JS_ClubCostCard   `json:"costcard"`   //! 俱乐部消耗记录
	Game       []JS_ClubGame       `json:"game"`       //! 俱乐部游戏
}

//! 俱乐部管理者
type ClubMgr struct {
	Club map[int64]*SQL_ClubMgr
	Lock *sync.RWMutex
}

var clubsingleton *ClubMgr = nil

//! public
func GetClubMgr() *ClubMgr {
	if clubsingleton == nil {
		clubsingleton = new(ClubMgr)
		clubsingleton.Lock = new(sync.RWMutex)
		clubsingleton.Club = make(map[int64]*SQL_ClubMgr)
	}

	return clubsingleton
}

func (self *ClubMgr) GetData() {
	var sql SQL_ClubMgr
	res := GetServer().DB.GetAllData("select * from `clubmgr`", &sql)

	for i := 0; i < len(res); i++ {
		data := res[i].(*SQL_ClubMgr)
		data.lock = new(sync.RWMutex)
		data.roomlock = new(sync.RWMutex)
		data.chatlock = new(sync.RWMutex)
		data.resultlock = new(sync.RWMutex)
		data.eventlock = new(sync.RWMutex)
		data.costlock = new(sync.RWMutex)
		data.Decode()
		self.Club[data.Id] = data
	}
}

func (self *ClubMgr) Save() {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	for _, value := range self.Club {
		value.Save()
	}
}

//! 加入俱乐部
func (self *ClubMgr) AddClub(club *SQL_ClubMgr) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	self.Club[club.Id] = club
}

//! 查找俱乐部
func (self *ClubMgr) GetClub(id int64) *SQL_ClubMgr {
	self.Lock.RLock()
	defer self.Lock.RUnlock()

	return self.Club[id]
}

//! 删除俱乐部
func (self *ClubMgr) DelClub(id int64) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	delete(self.Club, id)
}

//! 查找俱乐部
func (self *ClubMgr) GetClubFromName(name string) *SQL_ClubMgr {
	self.Lock.RLock()
	defer self.Lock.RUnlock()

	for _, value := range self.Club {
		if value.info.Name == name {
			return value
		}
	}

	return nil
}

//! 得到俱乐部列表
func (self *ClubMgr) GetClubList(person *Person) {
	var msg Msg_ClubList

	self.Lock.RLock()
	defer self.Lock.RUnlock()

	for _, value := range self.Club {
		if value.Id != 100001 && !person.GetModule("club").(*Mod_Club).IsMyClub(value.Id) {
			continue
		}
		var son Son_ClubList
		son.Id = value.Id
		son.Name = value.info.Name
		son.Icon = value.info.Icon
		son.Notice = value.info.ExNotice
		son.Num = len(value.info.Mem)
		son.Game = value.info.Game
		msg.Info = append(msg.Info, son)
	}

	person.SendMsg(lib.HF_EncodeMsg("getclublist", &msg, true))
}

//! 申请俱乐部
func (self *ClubMgr) ApplyClub(person *Person, clubid int64) int {
	self.Lock.RLock()
	defer self.Lock.RUnlock()

	club, ok := self.Club[clubid]
	if !ok {
		return 1
	}

	if !club.AddApply(person.Uid, person.mod_info.Name, person.mod_info.Imgurl) {
		return 2
	}

	return 0
}

//! 取消申请
func (self *ClubMgr) UnApplyClub(person *Person, clubid int64) bool {
	self.Lock.RLock()
	defer self.Lock.RUnlock()

	club, ok := self.Club[clubid]
	if !ok {
		return false
	}

	club.DelApply(person.Uid)

	return true
}

//! 创建俱乐部
func (self *ClubMgr) Create(person *Person, name string, icon string, notice string, mode int) *SQL_ClubMgr {
	if !lib.HF_IsLicitName([]byte(name)) {
		person.SendErr("该名字不合法")
		return nil
	}

	club := self.GetClubFromName(name)
	if club != nil {
		person.SendErr("该名字已被占用")
		return nil
	}

	club = new(SQL_ClubMgr)
	club.lock = new(sync.RWMutex)
	club.roomlock = new(sync.RWMutex)
	club.chatlock = new(sync.RWMutex)
	club.resultlock = new(sync.RWMutex)
	club.eventlock = new(sync.RWMutex)
	club.costlock = new(sync.RWMutex)
	club.info.Name = name
	club.info.ExNotice = notice
	club.info.Mode = mode
	if icon == "" {
		club.info.Icon = person.mod_info.Imgurl
	} else {
		club.info.Icon = icon
	}
	club.info.Host = person.Uid
	club.info.Mem = append(club.info.Mem, JS_ClubMem{person.Uid, person.mod_info.Name, person.mod_info.Imgurl, CLUBMEM_HOST, true})
	club.info.Game = make([]JS_ClubGame, 0)
	club.AddEvent(person.mod_info.Uid, person.mod_info.Name, person.mod_info.Imgurl, fmt.Sprintf("%s加入俱乐部", person.mod_info.Name))
	club.info.Apply = make([]JS_ClubApply, 0)
	club.info.RoomChat = make([]JS_ClubRoomChat, 0)
	club.info.RoomResult = make([]JS_ClubRoomResult, 0)
	club.info.Event = make([]JS_ClubEvent, 0)
	club.info.CostCard = make([]JS_ClubCostCard, 0)
	club.Encode()
	sql := fmt.Sprintf("insert into `%s`(`info`) values (?)", "clubmgr")
	result, err := GetServer().DB.DB.Exec(sql, club.Info)
	if err != nil {
		person.SendErr("创建俱乐部失败")
		return nil
	}
	club.Id, _ = result.LastInsertId()
	self.AddClub(club)

	return club
}

//! 解散俱乐部
func (self *ClubMgr) Dissmiss(id int64) {
	sql := fmt.Sprintf("delete from `clubmgr` where `id` = %d", id)
	GetServer().DB.DB.Exec(sql)
	self.DelClub(id)
}
