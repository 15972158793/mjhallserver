package centerserver

import (
//"encoding/json"
//"fmt"
//"staticfunc"

//"github.com/garyburd/redigo/redis"
)

//type Mod_Bank struct {
//	person *Person
//	bank   SQL_Bank
//}
//type SQL_Bank struct {
//	Uid   int64
//	info  JS_Bank
//	Value []byte
//}
//type JS_Bank struct {
//	money    int `json:"money"`
//	password int `json:"password"`
//}
//type Bank_Money struct {
//	Money int `json:"money"`
//}
//type Bank_MoneyGet struct {
//	Money int  `json:"money"`
//	First bool `json:"first"`
//}
//type Msg_BankSave struct {
//	Uid   int64 `json:"uid"`
//	Money int   `json:"money"`
//}
//type Msg_BankDraw struct {
//	Uid      int64 `json:"uid"`
//	Money    int   `json:"money"`
//	Password int   `json:"password"`
//}
//type Msg_BankPwd struct {
//	Uid      int64 `json:"uid"`
//	Password int   `json:"password"`
//	NewPwd   int   `json:"newpwd"`
//}

//const INIT_PASSWORD = 8888

//func (self *SQL_Bank) Decode() {
//	json.Unmarshal(self.Value, &self.info)
//}

//func (self *SQL_Bank) Encode() {
//	self.Value = staticfunc.HF_JtoB(self.info)
//}

//func (self *Mod_Bank) OnGetData(person *Person) {
//	self.person = person

//	c := GetServer().Redis.Get()
//	defer c.Close()
//	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("bank_%d", self.person.Uid)))
//	if err == nil {
//		self.bank.Value = v
//		self.bank.Decode()
//	} else {
//		sql := fmt.Sprintf("select * from `bank` where uid = %d", self.person.Uid)
//		GetServer().DB.GetOneData(sql, &self.bank)
//		if self.bank.Uid <= 0 {
//			self.bank.Uid = self.person.Uid
//			self.bank.info.money = 0
//			self.bank.info.password = INIT_PASSWORD
//			self.bank.Encode()
//			sql := fmt.Sprintf("insert into `%s`(`uid`, `value`) values (%d, ?)", "bank", self.person.Uid)
//			GetServer().SqlQueue(sql, self.bank.Value, true)

//			c := GetServer().Redis.Get()
//			defer c.Close()
//			c.Do("SET", fmt.Sprintf("bank_%d", self.person.Uid), self.bank.Value)
//		} else {
//			self.bank.Decode()
//		}
//	}
//}
//func (self *Mod_Bank) OnGetOtherData() {
//}
//func (self *Mod_Bank) OnMsg(head string, body []byte) bool {
//	switch head {
//	case "bankget": //！获取银行金币
//		self.SendInfo()
//		return true
//	case "banksave": //!存钱
//		var msg Msg_BankSave
//		json.Unmarshal(body, &msg)
//		self.BankSave(msg.Uid, msg.Money)
//		return true
//	case "bankdraw": //!取钱
//		var msg Msg_BankDraw
//		json.Unmarshal(body, &msg)
//		self.BankDraw(msg.Uid, msg.Money, msg.Password)
//		return true
//	case "bankpwd": //!修改密码
//		var msg Msg_BankPwd
//		json.Unmarshal(body, &msg)
//		self.BankPwd(msg.Uid, msg.Password, msg.NewPwd)
//		return true
//	}
//	return false
//}
//func (self *Mod_Bank) BankPwd(uid int64, password int, newpwd int) {
//	if self.person.Uid != uid {
//		self.person.SendErr("不能操作他人银行")
//		return
//	}
//	if password != self.bank.info.password && self.bank.info.password != INIT_PASSWORD {
//		self.person.SendErr("密码错误")
//		return
//	}
//	self.bank.info.password = newpwd
//	self.OnSave(false)
//	self.person.SendMsg(staticfunc.HF_EncodeMsg("bankpwd", nil, true))
//}

//func (self *Mod_Bank) BankDraw(uid int64, money int, password int) {
//	if self.person.Uid != uid {
//		self.person.SendErr("不能操作他人银行")
//		return
//	}
//	if password != self.bank.info.password {
//		self.person.SendErr("密码错误")
//		return
//	}
//	ok, card, gold := GetServer().AddGold(self.person.Uid, money, "", 0)
//	if !ok {
//		self.person.SendErr("使用失败")
//		return
//	}
//	if self.person.session != nil {
//		self.person.UpdCard(card, gold) //！更新房卡数量
//	}
//	self.bank.info.money -= money
//	self.OnSave(false)
//	var msg Bank_Money
//	msg.Money = self.bank.info.money
//	self.person.SendMsg(staticfunc.HF_EncodeMsg("bankdraw", &msg, true))
//}

//func (self *Mod_Bank) BankSave(uid int64, money int) { //存钱
//	if self.person.Uid != uid {
//		//////////////////////////////////////////
//		self.person.SendErr("不能操作他人银行")
//		return
//	}
//	//person := GetPersonMgr().GetPerson(uid, true)
//	//cost := money
//	ok, card, gold := GetServer().CostGold(self.person.Uid, money, "", 0)
//	if !ok {
//		self.person.SendErr("使用失败")
//		return
//	}
//	if self.person.session != nil {
//		self.person.UpdCard(card, gold) //！更新房卡数量
//	}
//	self.bank.info.money += money
//	self.OnSave(false)
//	var msg Bank_Money
//	msg.Money = self.bank.info.money
//	self.person.SendMsg(staticfunc.HF_EncodeMsg("banksave", &msg, true))
//}

//func (self *Mod_Bank) OnSave(sql bool) {
//	self.bank.Encode()
//	c := GetServer().Redis.Get()
//	defer c.Close()
//	c.Do("SET", fmt.Sprintf("bank_%d", self.person.Uid), self.bank.Value)

//	GetServer().SqlQueue(fmt.Sprintf("update `%s` set `value` = ? where `uid` = '%d'", "bank", self.person.Uid), self.bank.Value, true)
//}
//func (self *Mod_Bank) SendInfo() {
//	var msg Bank_MoneyGet
//	msg.Money = self.bank.info.money
//	if self.bank.info.password == INIT_PASSWORD {
//		msg.First = true
//	} else {
//		/////////////////////////////////////////////////////////////////
//		lib.GetLogMgr().Output(lib.LOG_DEBUG, "密码", self.bank.info.password)
//		////////////////////////////////////////////////////////////////
//		msg.First = false
//	}
//	self.person.SendMsg(staticfunc.HF_EncodeMsg("bank", &msg, true))
//}
