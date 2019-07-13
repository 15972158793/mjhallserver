package loginserver

import (
	"encoding/json"
	//"fmt"
	//"log"
	//"staticfunc"
	"sort"
	"sync"
	//"github.com/garyburd/redigo/redis"
)

func (s lstTop) Len() int           { return len(s) }
func (s lstTop) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s lstTop) Less(i, j int) bool { return s[i].Gold > s[j].Gold }

type JS_Top struct {
	Uid  int64  `json:"uid"`
	Name string `json:"name"`
	Head string `json:"head"`
	Gold int    `json:"gold"`
}

type lstTop []JS_Top

type TopMgr struct {
	Top  lstTop
	Ver  int
	Lock *sync.RWMutex
}

var topmgrSingleton *TopMgr = nil

//! 得到服务器指针
func GetTopMgr() *TopMgr {
	if topmgrSingleton == nil {
		topmgrSingleton = new(TopMgr)
		topmgrSingleton.Lock = new(sync.RWMutex)
	}

	return topmgrSingleton
}

//! 载入数据
func (self *TopMgr) GetData() {
	var db DB_Strc
	res := GetServer().DB.GetAllData("select * from `user` order by `gold` desc limit 50", &db)
	for i := 0; i < len(res); i++ {
		var p Person
		err := json.Unmarshal(res[i].(*DB_Strc).Value, &p)
		if err != nil {
			continue
		}
		var node JS_Top
		node.Uid = p.Uid
		node.Name = p.Name
		node.Head = p.Imgurl
		node.Gold = p.Gold
		self.Top = append(self.Top, node)
	}
	self.Ver = 1231
}

//!
func (self *TopMgr) UpdData(person *Person) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	needsort := false
	for i := 0; i < len(self.Top); i++ {
		if self.Top[i].Uid == person.Uid {
			self.Top[i].Name = person.Name
			self.Top[i].Head = person.Imgurl
			self.Top[i].Gold = person.Gold
			if i > 0 && self.Top[i].Gold > self.Top[i-1].Gold { //! 比上一名要多
				needsort = true
				break
			} else if i < len(self.Top)-1 && self.Top[i].Gold < self.Top[i+1].Gold { //! 比下一名要少
				needsort = true
				break
			} else {
				return
			}
		}
	}

	if needsort {
		sort.Sort(lstTop(self.Top))
		self.Ver++
		return
	}

	if len(self.Top) == 0 {
		self.Top = append(self.Top, JS_Top{person.Uid, person.Name, person.Imgurl, person.Gold})
		self.Ver++
		return
	}

	if len(self.Top) < 50 {
		self.Top = append(self.Top, JS_Top{person.Uid, person.Name, person.Imgurl, person.Gold})
		if self.Top[len(self.Top)-1].Gold > self.Top[len(self.Top)-2].Gold {
			sort.Sort(lstTop(self.Top))
			self.Ver++
			return
		}
	}

	if person.Gold > self.Top[len(self.Top)-1].Gold {
		self.Top = append(self.Top, JS_Top{person.Uid, person.Name, person.Imgurl, person.Gold})
		sort.Sort(lstTop(self.Top))
		self.Top = self.Top[0:50]
		self.Ver++
	}
}

//!
func (self *TopMgr) UpdBaseInfo(person *Person) {
	self.Lock.RLock()
	defer self.Lock.RUnlock()

	for i := 0; i < len(self.Top); i++ {
		if self.Top[i].Uid == person.Uid {
			self.Top[i].Name = person.Name
			self.Top[i].Head = person.Imgurl
			self.Ver++
			break
		}
	}
}

func (self *TopMgr) GetTop(ver int) lstTop {
	self.Lock.RLock()
	defer self.Lock.RUnlock()

	if self.Ver == ver {
		return make(lstTop, 0)
	}

	return self.Top
}
