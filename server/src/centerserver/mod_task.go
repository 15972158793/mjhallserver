package centerserver

import (
//"encoding/json"
//"fmt"
////	"log"
////	"log"
//"staticfunc"
//"time"

//"github.com/garyburd/redigo/redis"
)

//type Mod_Task struct {
//	person *Person
//	task   SQL_Task
//}
//type SQL_Task struct {
//	Uid   int64
//	Value []byte
//	info  JS_Task
//}

//type JS_Task struct {
//	TaskInfo []Task_Info `json:"taskinfo"`
//}

//type Task_Info struct {
//	Id       int   `json:"id"`       //! 任务id
//	GameId   int   `json:"gameid"`   //! 对应游戏id
//	Times    int   `json:"times"`    //! 任务完成的进度
//	Complete int   `json:"complete"` //! 任务完成的需要的次数
//	State    bool  `json:"state"`    //! 任务是否已经领取
//	Daytime  int64 `json:"daytime"`  //! 第一次做任务的时间
//}

//func (self *SQL_Task) Decode() {
//	json.Unmarshal(self.Value, &self.info)
//}

//func (self *SQL_Task) Encode() {
//	self.Value = staticfunc.HF_JtoB(self.info)
//}
//func (self *Mod_Task) OnGetData(person *Person) {
//	self.person = person

//	c := GetServer().Redis.Get()
//	defer c.Close()
//	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("task_%d", self.person.Uid)))
//	if err == nil {
//		self.task.Value = v
//		self.task.Decode()

//		//判断时间		//先紧针对南昌
//		curtime := time.Now().Unix()
//		if curtime-self.task.info.TaskInfo[1].Daytime >= 86400 {
//			self.task.info.TaskInfo[1].Id = 1
//			self.task.info.TaskInfo[1].GameId = 33
//			self.task.info.TaskInfo[1].Times = 0
//			self.task.info.TaskInfo[1].Complete = 20
//			self.task.info.TaskInfo[1].State = false
//			self.task.info.TaskInfo[1].Daytime = time.Now().Unix()

//			self.task.info.TaskInfo[2].Id = 2
//			self.task.info.TaskInfo[2].GameId = 33
//			self.task.info.TaskInfo[2].Times = 0
//			self.task.info.TaskInfo[2].Complete = 8
//			self.task.info.TaskInfo[2].State = false
//			self.task.info.TaskInfo[2].Daytime = time.Now().Unix()

//			self.task.info.TaskInfo[3].Id = 3
//			self.task.info.TaskInfo[3].GameId = 33
//			self.task.info.TaskInfo[3].Times = 0
//			self.task.info.TaskInfo[3].Complete = 1
//			self.task.info.TaskInfo[3].State = false
//			self.task.info.TaskInfo[3].Daytime = time.Now().Unix()
//		}

//	} else {
//		sql := fmt.Sprintf("select * from `task` where uid = %d", self.person.Uid)
//		GetServer().DB.GetOneData(sql, &self.task)
//		if self.task.Uid <= 0 {
//			self.task.Uid = self.person.Uid
//			self.task.info.TaskInfo = make([]Task_Info, 4)

//			self.task.info.TaskInfo[1].Id = 1
//			self.task.info.TaskInfo[1].GameId = 33
//			self.task.info.TaskInfo[1].Times = 0
//			self.task.info.TaskInfo[1].Complete = 20
//			self.task.info.TaskInfo[1].State = false
//			self.task.info.TaskInfo[1].Daytime = time.Now().Unix()

//			self.task.info.TaskInfo[2].Id = 2
//			self.task.info.TaskInfo[2].GameId = 33
//			self.task.info.TaskInfo[2].Times = 0
//			self.task.info.TaskInfo[2].Complete = 8
//			self.task.info.TaskInfo[2].State = false
//			self.task.info.TaskInfo[2].Daytime = time.Now().Unix()

//			self.task.info.TaskInfo[3].Id = 3
//			self.task.info.TaskInfo[3].GameId = 33
//			self.task.info.TaskInfo[3].Times = 0
//			self.task.info.TaskInfo[3].Complete = 1
//			self.task.info.TaskInfo[3].State = false
//			self.task.info.TaskInfo[3].Daytime = time.Now().Unix()

//			self.task.Encode()
//			sql := fmt.Sprintf("insert into `%s`(`uid`, `value`) values (%d, ?)", "task", self.person.Uid)
//			GetServer().SqlQueue(sql, self.task.Value, true)

//			c := GetServer().Redis.Get()
//			defer c.Close()
//			c.Do("SET", fmt.Sprintf("task_%d", self.person.Uid), self.task.Value)
//		} else {
//			self.task.Decode()
//		}
//	}
//}

//func (self *Mod_Task) OnGetOtherData() {

//}

//func (self *Mod_Task) SendInfo() {
//	self.person.SendMsg(staticfunc.HF_EncodeMsg("task", &self.task.info, true))
//}
//func (self *Mod_Task) OnMsg(head string, body []byte) bool {
//	switch head {
//	case "taskGet": //！获取所有任务列表
//		self.SendInfo()
//		return true
//	case "taskFinish": //!任务完成
//		var msg C2S_TaskId
//		json.Unmarshal(body, &msg)
//		self.TaskFinish(msg.Id)
//		return true
//	case "taskTimes": //!任务完成进度
//		var msg C2S_TaskTimes
//		json.Unmarshal(body, &msg)
//		self.TaskTimes(msg.Id, msg.Times)
//		return true
//	}
//	return false
//}

//func (self *Mod_Task) TaskFinish(id int) {
//	var i int
//	for i = 0; i < len(self.task.info.TaskInfo); i++ {
//		if id == self.task.info.TaskInfo[i].Id {
//			if self.task.info.TaskInfo[i].Times != self.task.info.TaskInfo[i].Complete {
//				self.person.SendErr("任务还没有完成")
//				return
//			}
//			if self.task.info.TaskInfo[i].State {
//				self.person.SendErr("奖励已经领取")
//				return
//			}
//			okcard, card, gold := GetServer().AddCard(self.person.Uid, 1, "", 0)
//			if !okcard {
//				self.person.SendErr("加卡失败")
//				return
//			}
//			if self.person.session != nil {
//				self.person.UpdCard(card, gold) //！更新房卡数量
//			}

//			self.task.info.TaskInfo[i].State = true

//			var msg S2C_TaskState
//			msg.State = self.task.info.TaskInfo[i].State
//			msg.Uid = self.person.Uid

//			self.OnSave(false)
//			self.person.SendMsg(staticfunc.HF_EncodeMsg("taskupdate", &msg, true))
//			break

//		}
//	}
//	if i >= len(self.task.info.TaskInfo) {
//		self.person.SendErr("没有找到任务")
//		return
//	}
//}

//func (self *Mod_Task) TaskTimes(id int, times int) {
//	var i int
//	for i = 0; i < len(self.task.info.TaskInfo); i++ {
//		if id == self.task.info.TaskInfo[i].Id {
//			if self.task.info.TaskInfo[i].State {
//				return
//			}
//			times += self.task.info.TaskInfo[i].Times

//			if times > self.task.info.TaskInfo[i].Complete {
//				self.task.info.TaskInfo[i].Times = self.task.info.TaskInfo[i].Complete
//			} else {
//				self.task.info.TaskInfo[i].Times = times
//			}

//			var msg S2C_TaskTimes
//			msg.Times = self.task.info.TaskInfo[i].Times
//			msg.Uid = self.person.Uid

//			self.OnSave(false)
//			self.person.SendMsg(staticfunc.HF_EncodeMsg("tasktimes", &msg, true))
//			break

//		}
//	}
//	if i >= len(self.task.info.TaskInfo) {
//		self.person.SendErr("没有找到任务")
//		return
//	}
//}

//func (self *Mod_Task) OnSave(sql bool) {
//	self.task.Encode()
//	c := GetServer().Redis.Get()
//	defer c.Close()
//	c.Do("SET", fmt.Sprintf("task_%d", self.person.Uid), self.task.Value)

//	GetServer().SqlQueue(fmt.Sprintf("update `%s` set `value` = ? where `uid` = '%d'", "task", self.person.Uid), self.task.Value, true)
//}
