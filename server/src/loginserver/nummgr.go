package loginserver

import (
	//"encoding/json"
	//"fmt"
	//"log"
	//"staticfunc"
	//"sort"
	"sync"
	//"github.com/garyburd/redigo/redis"
)

type JS_GameNum struct {
	GameType int `json:"gametype"`
	Num      int `json:"num"`
}

type NumMgr struct {
	GameNum  map[int]map[int64]int
	GameLock *sync.RWMutex
}

var nummgrSingleton *NumMgr = nil

//! 得到服务器指针
func GetNumMgr() *NumMgr {
	if nummgrSingleton == nil {
		nummgrSingleton = new(NumMgr)
		nummgrSingleton.GameNum = make(map[int]map[int64]int)
		nummgrSingleton.GameLock = new(sync.RWMutex)
	}

	return nummgrSingleton
}

func (self *NumMgr) AddGameOne(gametype int, uid int64) {
	if gametype == 0 {
		return
	}

	self.GameLock.Lock()
	defer self.GameLock.Unlock()

	_, ok := self.GameNum[gametype]
	if !ok {
		self.GameNum[gametype] = make(map[int64]int)
	}

	self.GameNum[gametype][uid] = 1
}

func (self *NumMgr) DoneGameOne(gametype int, uid int64) {
	if gametype == 0 {
		return
	}

	self.GameLock.Lock()
	defer self.GameLock.Unlock()

	_, ok := self.GameNum[gametype]
	if !ok {
		return
	}

	delete(self.GameNum[gametype], uid)
}

func (self *NumMgr) GetGameNum() []JS_GameNum {
	self.GameLock.RLock()
	defer self.GameLock.RUnlock()

	lst := make([]JS_GameNum, 0)
	for key, value := range self.GameNum {
		if key < 10000 { //! 只发送金币场的人数
			continue
		}

		if len(value) == 0 {
			continue
		}

		lst = append(lst, JS_GameNum{key, len(value)})
	}

	return lst
}

func (self *NumMgr) GetNum() int {
	self.GameLock.RLock()
	defer self.GameLock.RUnlock()

	num := 0
	for _, value := range self.GameNum {
		num += len(value)
	}

	return num
}

func (self *NumMgr) GetNumByGameType(gametype int) int {
	self.GameLock.RLock()
	defer self.GameLock.RUnlock()

	return len(self.GameNum[gametype])
}
