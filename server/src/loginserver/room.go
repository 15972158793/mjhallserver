//! 房间号生成器

package loginserver

import (
	//"encoding/json"
	//"fmt"
	//"log"
	"lib"
	"sync"

	"github.com/garyburd/redigo/redis"
)

type Room struct {
	RoomId int
	Lock   *sync.RWMutex
}

var roomSingleton *Room = nil

//! 得到服务器指针
func GetRoom() *Room {
	if roomSingleton == nil {
		roomSingleton = new(Room)
		roomSingleton.Lock = new(sync.RWMutex)
		roomSingleton.get()
	}

	return roomSingleton
}

//! 得到一个id
func (self *Room) GetID() int {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	self.RoomId++
	if self.RoomId > 9999 {
		self.RoomId = 0
	}

	self.set()

	return (lib.HF_GetRandom(9)+1)*100000 + self.RoomId*10 + lib.HF_GetRandom(10)
}

//! 得到一个俱乐部id
func (self *Room) GetClubID() int {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	self.RoomId++
	if self.RoomId > 9999 {
		self.RoomId = 0
	}

	self.set()

	return (lib.HF_GetRandom(90)+10)*100000 + self.RoomId*10 + lib.HF_GetRandom(10)
}

func (self *Room) GetLzTenhalfID() int {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	self.RoomId++
	bai := self.RoomId / 100
	shi := self.RoomId / 10 % 10
	ge := self.RoomId % 10
	if bai == 4 || bai == 7 {
		self.RoomId += 100
	}
	if shi == 4 || shi == 7 {
		self.RoomId += 10
	}
	if ge == 4 || ge == 7 {
		self.RoomId += 1
	}
	if self.RoomId > 999 {
		self.RoomId = 0
	}

	self.set()

	return self.RoomId
}

func (self *Room) get() {
	//! 保存到redis
	c := GetServer().Redis.Get()
	defer c.Close()
	v, err := redis.Int(c.Do("GET", "roomid_new"))
	if err == nil {
		self.RoomId = v
		return
	}

	self.RoomId = 1231
}

func (self *Room) set() {
	c := GetServer().Redis.Get()
	defer c.Close()
	c.Do("SET", "roomid_new", self.RoomId)
}
