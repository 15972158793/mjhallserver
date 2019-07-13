//! ip黑名单
package staticfunc

import (
	"lib"
	"sync"
)

type IpBlackMgr struct {
	Black map[string]int
	lock  *sync.RWMutex
}

var ipblackSingleton *IpBlackMgr = nil

//! 得到服务器指针
func GetIpBlackMgr() *IpBlackMgr {
	if ipblackSingleton == nil {
		ipblackSingleton = new(IpBlackMgr)
		ipblackSingleton.lock = new(sync.RWMutex)
		ipblackSingleton.Black = make(map[string]int)
	}

	return ipblackSingleton
}

//! 加入一个黑名单
func (self *IpBlackMgr) AddIp(ip string, info string) {
	self.lock.Lock()
	defer self.lock.Unlock()

	value, ok := self.Black[ip]
	if ok {
		self.Black[ip] = value + 1
	} else {
		self.Black[ip] = 1
	}

	lib.GetLogMgr().Output(lib.LOG_ERROR, ip, ":", info)
}

//! 是否在黑名单
func (self *IpBlackMgr) IsIp(ip string) bool {
	self.lock.RLock()
	defer self.lock.RUnlock()

	value, _ := self.Black[ip]

	//! 已经被记录了5次不合常规
	return value >= 10000
}
