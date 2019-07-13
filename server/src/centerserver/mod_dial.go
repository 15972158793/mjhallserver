package centerserver

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"lib"
	"rjmgr"
	"time"
)

var DAIL_DS []int = []int{28900, 2, 8, 3888, 88, 288, 188, 8888, 6, 1000, 10000, 68}
var FREE int = 1

//!　概率 2~5% 6~5% 8~5% 68~15% 188~30% 288~20% 1000~9% 3888~5% 8888~3% 10000~2% 28900~1%

//! 大转盘
type Mod_Dial struct {
	person *Person
	dial   SQL_Dial
}

type SQL_Dial struct {
	Uid  int64 `json:"uid"`
	Num  int   `json:"num"`  //! 转了多少次
	Time int64 `json:"time"` //! 免费转的时间

	lib.DataUpdate
}

type Msg_DialInfo struct {
	Uid      int64 `json:"uid"`
	Surplus  int   `json:"surplus"`  //! 剩余免费转的次数
	Free     int   `json:"free"`     //! 一共免费多少次
	Expend1  int   `json:"expend1"`  //! 转一次消耗
	Expend5  int   `json:"expend5"`  //! 转五次消耗
	Expend10 int   `json:"expend10"` //! 转十次消耗
	DS       []int `json:"ds"`       //! 转盘上每格的数值
}

type Msg_DialResult struct {
	Uid     int64 `json:"uid"`
	Result  []int `json:"result"`  //! 每次的落点
	Get     int   `json:"get"`     //! 总收益
	Surplus int   `json:"surplus"` //! 剩余免费转的次数
}

func (self *Mod_Dial) OnGetData(person *Person) {
	self.person = person
}

func (self *Mod_Dial) OnGetOtherData() {
	c := GetServer().Redis.Get()
	defer c.Close()

	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("dial_%d", self.person.Uid)))
	if err == nil {
		json.Unmarshal(v, &self.dial)
	} else {
		sql := fmt.Sprintf("select * from `dial` where uid = %d ", self.person.Uid)
		GetServer().DB.GetOneData(sql, &self.dial)
		if self.dial.Uid <= 0 {
			self.dial.Uid = self.person.Uid
			self.dial.Num = 0
			self.dial.Time = 0
			sql := fmt.Sprintf("insert into `%s`(`uid`,`num`,`time`) value (%d, %d, %d)", "dial", self.person.Uid, 0, 0)
			GetServer().SqlQueue(sql, []byte(""), false)

			c := GetServer().Redis.Get()
			defer c.Close()
			c.Do("SET", fmt.Sprintf("dial_%d", self.person.Uid), lib.HF_JtoB(&self.dial))
			c.Do("EXPIRE", fmt.Sprintf("dial_%d", self.person.Uid), 86400*7)
		}
	}

	self.dial.Init("dial", &self.dial, GetServer().DB)
}

func (self *Mod_Dial) OnMsg(head string, body []byte) bool {
	switch head {
	case "dial":
		var msg Msg_ReadMail
		json.Unmarshal(body, &msg)
		self.Rotate(msg.Id)
	}
	return false
}

//! 旋转多少次
func (self *Mod_Dial) Rotate(num int) {
	if rjmgr.GetRJMgr().CloseDial == 1 {
		self.person.SendErr("该功能暂未开放")
		return
	}

	self.InitDial()

	cost := 0

	if num == 1 { //! 转一次
		if FREE-self.dial.Num <= 0 { //! 不是免费
			cost = 100
		} else { //! 免费 num++
			self.dial.Num++
			self.dial.Time = time.Now().Unix()
		}
	} else if num == 5 {
		cost = 470
	} else if num == 10 {
		cost = 900
	}
	if GetServer().Con.MoneyMode == 2 { // 1:10000
		cost *= 100
	}
	_, gold, _ := GetServer().GetCard(self.person.Uid)

	if gold < cost {
		self.person.SendErr("金币不足")
		return
	}

	self.OnSave(false)

	var msg Msg_DialResult
	msg.Uid = self.person.Uid
	for i := 0; i < num; i++ {
		result := 0
		rand := lib.HF_GetRandom(100)
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "----------- rand : ", rand)
		if rand < 50 {
			result = 2
		} else if rand < 65 {
			result = 6
		} else if rand < 85 {
			result = 8
		} else if rand < 95 {
			result = 68
		} else if rand < 100 {
			result = 188
		}

		find := false
		for i := 0; i < len(DAIL_DS); i++ {
			if DAIL_DS[i] == result {
				msg.Result = append(msg.Result, i)
				msg.Get += DAIL_DS[i]
				find = true
				break
			}
		}
		if !find { //!　没找到转盘上的值，分配最小的值
			msg.Result = append(msg.Result, 1)
			msg.Get += DAIL_DS[1]
		}
	}

	if GetServer().Con.MoneyMode == 2 { // 1:10000
		msg.Get *= 100
	}
	money := msg.Get - cost
	if money > 0 {
		ok, card, gold := GetServer().AddGold(self.person.Uid, money, "", -8)
		if !ok {
			return
		}
		self.person.UpdCard(card, gold)
	} else if money < 0 {
		ok, card, gold := GetServer().CostGold(self.person.Uid, -money, "", -9)
		if !ok {
			return
		}
		self.person.UpdCard(card, gold)
	}

	self.person.SendMsg(lib.HF_EncodeMsg("dail", &msg, true))

}

//! 重置免费转
func (self *Mod_Dial) InitDial() {
	if self.dial.Time == 0 {
		return
	}

	time1 := time.Unix(self.dial.Time, 0)
	if time1.Year() != time.Now().Year() || time1.Month() != time.Now().Month() || time1.Day() != time.Now().Day() {
		self.dial.Num = 0
	}
}

func (self *Mod_Dial) OnSave(sql bool) {
	c := GetServer().Redis.Get()
	defer c.Close()
	c.Do("SET", fmt.Sprintf("dial_%d", self.person.Uid), lib.HF_JtoB(&self.dial))
	c.Do("EXPIRE", fmt.Sprintf("dial_%d", self.person.Uid), 86400*7)

	GetServer().SqlQueue(fmt.Sprintf("update `%s` set `num` = %d ,time = %d where `uid` = %d ", "dial", self.dial.Num, self.dial.Time, self.person.Uid), []byte(""), false)
}

func (self *Mod_Dial) SendInfo() {
	self.InitDial()

	var msg Msg_DialInfo
	msg.Uid = self.dial.Uid
	lib.HF_DeepCopy(&msg.DS, &DAIL_DS)
	if GetServer().Con.MoneyMode == 2 { // 1:10000
		for i := 0; i < len(msg.DS); i++ {
			msg.DS[i] *= 100
		}
		msg.Expend1 = 10000
		msg.Expend5 = 47000
		msg.Expend10 = 90000
	} else {
		msg.Expend1 = 100
		msg.Expend5 = 470
		msg.Expend10 = 900
	}
	msg.Surplus = FREE - self.dial.Num
	msg.Free = FREE

	self.person.SendMsg(lib.HF_EncodeMsg("dialinfo", &msg, true))
}
