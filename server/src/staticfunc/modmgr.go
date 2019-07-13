package staticfunc

import (
	"encoding/json"
	//"fmt"
	"lib"

	"github.com/garyburd/redigo/redis"
	//"sync"
)

var DefaulModProperty = ModProperty{true, true, 2, 500, 200}

//!
type ModProperty struct {
	OpenDial  bool `json:"opendial"`  //! 是否打开转盘
	OpenAlms  bool `json:"openalms"`  //! 是否打开救济金
	AlmsNum   int  `json:"almsnum"`   //! 每天领取救济金次数
	AlmsMoney int  `json:"almsmoney"` //! 救济金数量
	AlmsLimit int  `json:"almslimit"` //! 少于多少领取
}

//! 模块管理者
type ModMgr struct {
}

var modmgrSingleton *ModMgr = nil

//! 得到服务器指针
func GetModMgr() *ModMgr {
	if modmgrSingleton == nil {
		modmgrSingleton = new(ModMgr)
	}

	return modmgrSingleton
}

//! 得到属性
func (self *ModMgr) GetModProperty(_redis *redis.Pool) ModProperty {
	//! 先从redis里面找一下
	c := _redis.Get()
	defer c.Close()

	var value ModProperty

	v, err := redis.Bytes(c.Do("GET", "modproperty"))
	if err == nil {
		json.Unmarshal(v, &value)
		return value
	}

	return DefaulModProperty
}

//! 设置属性
func (self *ModMgr) SetModProperty(_redis *redis.Pool, value ModProperty) {
	//! 先从redis里面找一下
	c := _redis.Get()
	defer c.Close()

	c.Do("SET", "modproperty", lib.HF_JtoB(&value))
}
