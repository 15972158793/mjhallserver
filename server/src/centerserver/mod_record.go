package centerserver

import (
//"encoding/json"
//"fmt"
////	"log"
//"staticfunc"

//"github.com/garyburd/redigo/redis"
)

//type Mod_Record struct {
//	person *Person
//	record SQL_Record
//}
//type SQL_Record struct {
//	Uid   int64
//	Value []byte
//	info  JS_Record
//}

//type JS_Record struct {
//	Record []Record_Info `json:"record"`
//}

//type Record_Info struct {
//	Id      int    `json:"id"`
//	Title   string `json:"title"`
//	Date    string `json:"date"`
//	Context string `json:"context"`
//}

//func (self *SQL_Record) Decode() {
//	json.Unmarshal(self.Value, &self.info)
//}

//func (self *SQL_Record) Encode() {
//	self.Value = staticfunc.HF_JtoB(self.info)
//}

//func (self *Mod_Record) GetInitDate() []Record_Info {
//	sss := []Record_Info{
//		{11, "今天的活动", "2017-11-1", "呵呵1"},
//		{12, "明天的活动", "2017-11-2", "呵呵2"},
//		{13, "后天的活动", "2017-11-3", "呵呵3"},
//		{14, "大后天的活动", "2017-11-4", "呵呵3"},
//		{15, "大大后台的活动", "2017-11-5", "呵呵5"}}
//	return sss
//}

//func (self *Mod_Record) OnGetData(person *Person) {
//	self.person = person

//	c := GetServer().Redis.Get()
//	defer c.Close()
//	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("record_%d", self.person.Uid)))
//	if err == nil {
//		self.record.Value = v
//		self.record.Decode()
//	} else {
//		sql := fmt.Sprintf("select * from `record` where uid = %d", self.person.Uid)
//		GetServer().DB.GetOneData(sql, &self.record)
//		if self.record.Uid <= 0 {
//			self.record.Uid = self.person.Uid
//			self.record.info.Record = make([]Record_Info, 5)
//			self.record.info.Record = self.GetInitDate()
//			self.record.Encode()
//			sql := fmt.Sprintf("insert into `%s`(`uid`, `value`) values (%d, ?)", "record", self.person.Uid)
//			GetServer().SqlQueue(sql, self.record.Value, true)

//			c := GetServer().Redis.Get()
//			defer c.Close()
//			c.Do("SET", fmt.Sprintf("record_%d", self.person.Uid), self.record.Value)
//		} else {
//			self.record.Decode()
//		}
//	}
//}

//func (self *Mod_Record) OnGetOtherData() {

//}

//func (self *Mod_Record) SendInfo() {
//	self.person.SendMsg(staticfunc.HF_EncodeMsg("record", &self.record.info, true))
//}
//func (self *Mod_Record) OnMsg(head string, body []byte) bool {
//	switch head {
//	case "recordGet": //！获取公告列表
//		self.SendInfo()
//		return true
//	}
//	return false
//}

//func (self *Mod_Record) OnSave(sql bool) {
//	self.record.Encode()
//	c := GetServer().Redis.Get()
//	defer c.Close()
//	c.Do("SET", fmt.Sprintf("record_%d", self.person.Uid), self.record.Value)

//	GetServer().SqlQueue(fmt.Sprintf("update `%s` set `value` = ? where `uid` = '%d'", "record", self.person.Uid), self.record.Value, true)
//}
