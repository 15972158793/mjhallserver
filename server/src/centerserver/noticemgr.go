package centerserver

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"lib"
	"sync"
)

type ReadMail struct {
	Read []int `json:"read"`
}

type Msg_NoticeList struct {
	Info []*Sql_Notice `json:"info"`
	Read []int         `json:"read"`
}

type Sql_Notice struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Date    string `json:"date"`
	Context string `json:"context"`

	lib.DataUpdate
}

/////////////////////////////////////////////////////////////
type NoticeMgr struct {
	Notice    []*Sql_Notice
	NoticeKWX []*Sql_Notice
	Lock      *sync.RWMutex
	MaxId     int
}

var noticemgrsingleton *NoticeMgr = nil

//! public
func GetNoticeMgr() *NoticeMgr {
	if noticemgrsingleton == nil {
		noticemgrsingleton = new(NoticeMgr)
		noticemgrsingleton.Notice = make([]*Sql_Notice, 0)
	}

	return noticemgrsingleton
}

func (self *NoticeMgr) GetData() {
	self.Notice = make([]*Sql_Notice, 0)
	self.NoticeKWX = make([]*Sql_Notice, 0)

	var notice Sql_Notice
	sql := fmt.Sprintf("select * from `notice`")
	res := GetServer().DB.GetAllData(sql, &notice)
	for i := 0; i < len(res); i++ {
		self.Notice = append(self.Notice, res[i].(*Sql_Notice))
		if res[i].(*Sql_Notice).Id > self.MaxId {
			self.MaxId = res[i].(*Sql_Notice).Id
		}
	}

	sql = fmt.Sprintf("select * from `noticekwx`")
	res = GetServer().DB.GetAllData(sql, &notice)
	for i := 0; i < len(res); i++ {
		self.NoticeKWX = append(self.NoticeKWX, res[i].(*Sql_Notice))
	}
}

func (self *NoticeMgr) InsertData(info string) bool {
	self.MaxId = self.MaxId + 1
	var notice Sql_Notice
	notice.Id = self.MaxId
	notice.Context = info

	lib.InsertTable("notice", &notice, 0, GetServer().DB)
	self.Notice = append(self.Notice, &notice)
	return true
}

func (self *NoticeMgr) SendInfo(person *Person) {
	var read ReadMail
	read.Read = make([]int, 0)

	c := GetServer().Redis.Get()
	defer c.Close()
	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("mail_%d", person.Uid)))
	if err == nil { //! redis找得到，则直接返回redis里的数据
		json.Unmarshal(v, &read)
	}

	var msg Msg_NoticeList
	msg.Info = self.Notice
	msg.Read = read.Read
	person.SendMsg(lib.HF_EncodeMsg("noticeinfo", &msg, true))
}

func (self *NoticeMgr) SendKWXInfo(person *Person) {
	var read ReadMail
	read.Read = make([]int, 0)

	c := GetServer().Redis.Get()
	defer c.Close()
	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("mailkwx_%d", person.Uid)))
	if err == nil { //! redis找得到，则直接返回redis里的数据
		json.Unmarshal(v, &read)
	}

	var msg Msg_NoticeList
	msg.Info = self.NoticeKWX
	msg.Read = read.Read
	person.SendMsg(lib.HF_EncodeMsg("noticeinfo", &msg, true))
}

func (self *NoticeMgr) ReadMail(person *Person, id int) {
	var read ReadMail
	read.Read = make([]int, 0)

	c := GetServer().Redis.Get()
	defer c.Close()

	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("mail_%d", person.Uid)))
	if err == nil { //! redis找得到，则直接返回redis里的数据
		json.Unmarshal(v, &read)
	}

	read.Read = append(read.Read, id)

	c = GetServer().Redis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("mail_%d", person.Uid), lib.HF_JtoB(&read))
}

func (self *NoticeMgr) ReadKWXMail(person *Person, id int) {
	var read ReadMail
	read.Read = make([]int, 0)

	c := GetServer().Redis.Get()
	defer c.Close()

	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("mailkwx_%d", person.Uid)))
	if err == nil { //! redis找得到，则直接返回redis里的数据
		json.Unmarshal(v, &read)
	}

	read.Read = append(read.Read, id)

	c = GetServer().Redis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("mailkwx_%d", person.Uid), lib.HF_JtoB(&read))
}
