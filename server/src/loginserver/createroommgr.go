package loginserver

import (
	"staticfunc"
	"sync"
	"time"
)

type CreateRoomInfo struct {
	RoomId int                           `json:"roomid"` //! 房间号
	Type   int                           `json:"type"`   //! 游戏id
	Num    int                           `json:"num"`    //! 房卡数量
	Param1 int                           `json:"param1"` //! 玩法
	Param2 int                           `json:"param2"`
	Cur    int                           `json:"cur"` //! 当前人数
	Max    int                           `json:"max"` //! 最大人数
	Mem    []staticfunc.JS_CreateRoomMem `json:"mem"`
	State  int                           `json:"state"` //! 状态 0未开始 1已开始 2已结束
	Time   int64                         `json:"time"`  //! 时间
}

//! 玩家创建房间列表
type CreateRoomMgr struct {
	Room map[int64][]*CreateRoomInfo //! 玩家创建列表

	lock *sync.RWMutex
}

var createroommgr *CreateRoomMgr = nil

//! 得到服务器指针
func GetCreateRoomMgr() *CreateRoomMgr {
	if createroommgr == nil {
		createroommgr = new(CreateRoomMgr)
		createroommgr.Room = make(map[int64][]*CreateRoomInfo)
		createroommgr.lock = new(sync.RWMutex)
	}

	return createroommgr
}

//! 创建房间
func (self *CreateRoomMgr) AddRoom(uid int64, info *CreateRoomInfo) {
	self.lock.Lock()
	defer self.lock.Unlock()

	value, ok := self.Room[uid]
	if ok {
		self.Room[uid] = append(self.Room[uid], info)
		return
	}

	value = append(value, info)
	self.Room[uid] = value
}

//! 销毁房间
func (self *CreateRoomMgr) DelRoom(uid int64, roomid int, agent bool) {
	self.lock.Lock()
	defer self.lock.Unlock()

	_, ok := self.Room[uid]
	if !ok {
		return
	}

	//! 先删除过期的房间
	for i := 0; i < len(self.Room[uid]); {
		if time.Now().Unix()-self.Room[uid][i].Time >= 24*3600 {
			copy(self.Room[uid][i:], self.Room[uid][i+1:])
			self.Room[uid] = self.Room[uid][:len(self.Room[uid])-1]
		} else {
			i++
		}
	}

	for i := 0; i < len(self.Room[uid]); i++ {
		if self.Room[uid][i].RoomId == roomid {
			if !agent { //! 不是代开房间，直接删除列表
				copy(self.Room[uid][i:], self.Room[uid][i+1:])
				self.Room[uid] = self.Room[uid][:len(self.Room[uid])-1]
			} else {
				self.Room[uid][i].State = 2
			}
			break
		}
	}
}

//! 加人
func (self *CreateRoomMgr) Add(uid int64, roomid int, num int, node staticfunc.JS_CreateRoomMem) {
	self.lock.Lock()
	defer self.lock.Unlock()

	_, ok := self.Room[uid]
	if !ok {
		return
	}

	for i := 0; i < len(self.Room[uid]); i++ {
		if self.Room[uid][i].RoomId == roomid {
			self.Room[uid][i].Cur += num
			if num > 0 {
				self.Room[uid][i].Mem = append(self.Room[uid][i].Mem, node)
			} else {
				for j := 0; j < len(self.Room[uid][i].Mem); j++ {
					if self.Room[uid][i].Mem[j].Uid == node.Uid {
						copy(self.Room[uid][i].Mem[j:], self.Room[uid][i].Mem[j+1:])
						self.Room[uid][i].Mem = self.Room[uid][i].Mem[:len(self.Room[uid][i].Mem)-1]
						break
					}
				}
			}
			break
		}
	}
}

//! 得到当前房间列表
func (self *CreateRoomMgr) Get(uid int64) []*CreateRoomInfo {
	self.lock.RLock()
	defer self.lock.RUnlock()

	value, ok := self.Room[uid]
	if !ok {
		return make([]*CreateRoomInfo, 0)
	}

	return value
}

//! 设置开始
func (self *CreateRoomMgr) SetBegin(uid int64, roomid int) {
	self.lock.Lock()
	defer self.lock.Unlock()

	_, ok := self.Room[uid]
	if !ok {
		return
	}

	for i := 0; i < len(self.Room[uid]); i++ {
		if self.Room[uid][i].RoomId == roomid {
			self.Room[uid][i].State = 1
			break
		}
	}
}

//! 更新分数
func (self *CreateRoomMgr) UpdScore(uid int64, roomid int, node []staticfunc.JS_CreateRoomMem) {
	self.lock.Lock()
	defer self.lock.Unlock()

	_, ok := self.Room[uid]
	if !ok {
		return
	}

	for i := 0; i < len(self.Room[uid]); i++ {
		if self.Room[uid][i].RoomId == roomid {
			for j := 0; j < len(self.Room[uid][i].Mem); j++ {
				for k := 0; k < len(node); k++ {
					if self.Room[uid][i].Mem[j].Uid == node[k].Uid {
						self.Room[uid][i].Mem[j].Score = node[k].Score
						break
					}
				}
			}
			break
		}
	}
}
