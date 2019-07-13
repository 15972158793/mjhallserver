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

type GoldRoom struct {
	RoomId int
	Lock   *sync.RWMutex
}

var goldroomSingleton *GoldRoom = nil

//! 得到服务器指针
func GetGoldRoom() *GoldRoom {
	if goldroomSingleton == nil {
		goldroomSingleton = new(GoldRoom)
		goldroomSingleton.Lock = new(sync.RWMutex)
		goldroomSingleton.get()
	}

	return goldroomSingleton
}

//! 得到一个id
func (self *GoldRoom) GetID() int {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	self.RoomId++
	if self.RoomId > 999999 {
		self.RoomId = 0
	}

	self.set()

	return (lib.HF_GetRandom(9)+1)*10000000 + self.RoomId*10 + lib.HF_GetRandom(10)
}

func (self *GoldRoom) get() {
	//! 保存到redis
	c := GetServer().Redis.Get()
	defer c.Close()
	v, err := redis.Int(c.Do("GET", "roomid_gold"))
	if err == nil {
		self.RoomId = v
		return
	}

	self.RoomId = 1231
}

func (self *GoldRoom) set() {
	c := GetServer().Redis.Get()
	defer c.Close()
	c.Do("SET", "roomid_gold", self.RoomId)
}
