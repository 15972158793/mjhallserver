package centerserver

import (
//"encoding/json"
//"fmt"
////"log"
//"staticfunc"

//"github.com/garyburd/redigo/redis"
)

//type Mod_Gift struct {
//	person *Person
//	gift   SQL_Gift
//}
//type SQL_Gift struct {
//	Uid   int64
//	info  JS_Gift
//	Value []byte
//}

//type JS_Gift struct {
//	Mine  []Mine_Gift  `json:"mine"`
//	Other []Other_Gift `json:"other"`
//}
//type Mine_Gift struct {
//	Id     int64 `json:"id"`
//	Amount int   `json:"amount"`
//}
//type Other_Gift struct {
//	Uid       int64        `json:"uid"`
//	GiftArray []Gift_Array `json:"giftarray"`
//}

//type Gift_Array struct {
//	Id     int64 `json:"id"`
//	Amount int   `json:"amount"`
//}

//func (self *SQL_Gift) Decode() {
//	json.Unmarshal(self.Value, &self.info)
//}

//func (self *SQL_Gift) Encode() {
//	self.Value = staticfunc.HF_JtoB(self.info)
//}
//func (self *Mod_Gift) OnGetData(person *Person) {
//	self.person = person

//	c := GetServer().Redis.Get()
//	defer c.Close()
//	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("gift_%d", self.person.Uid)))
//	if err == nil {
//		self.gift.Value = v
//		self.gift.Decode()
//	} else {
//		sql := fmt.Sprintf("select * from `gift` where uid = %d", self.person.Uid)
//		GetServer().DB.GetOneData(sql, &self.gift)
//		if self.gift.Uid <= 0 {
//			self.gift.Uid = self.person.Uid
//			self.gift.info.Mine = make([]Mine_Gift, 0)
//			self.gift.info.Other = make([]Other_Gift, 0)
//			self.gift.Encode()
//			sql := fmt.Sprintf("insert into `%s`(`uid`, `value`) values (%d, ?)", "gift", self.person.Uid)
//			GetServer().SqlQueue(sql, self.gift.Value, true)

//			c := GetServer().Redis.Get()
//			defer c.Close()
//			c.Do("SET", fmt.Sprintf("gift_%d", self.person.Uid), self.gift.Value)
//		} else {
//			self.gift.Decode()
//		}
//	}
//}

//func (self *Mod_Gift) OnGetOtherData() {
//}
//func (self *Mod_Gift) SendInfo() {
//	self.person.SendMsg(staticfunc.HF_EncodeMsg("gift", &self.gift.info, true))
//}
//func (self *Mod_Gift) OnMsg(head string, body []byte) bool {
//	switch head {
//	case "giftGet": //！获取礼物列表
//		self.SendInfo()
//		return true
//	case "giftAdd":
//		var msg C2S_GiftBuy
//		json.Unmarshal(body, &msg)
//		self.GiftAdd(msg.Type, msg.Id, msg.Amount)
//		return true
//	case "giftgive":
//		var msg C2S_GiftGive
//		json.Unmarshal(body, &msg)
//		self.GiftGive(msg.Uid, msg.Type, msg.Id, msg.Amount)
//		return true
//	}
//	return false
//}
//func (self *Mod_Gift) GiftGive(desuid int64, _type int, gift int64, amount int) {
//	desPerson := GetPersonMgr().GetPerson(desuid, true)
//	if nil == desPerson {
//		//////////////////////////////////////////
//		self.person.SendErr("赠送的玩家不存在")
//		return
//	}
//	var msg S2C_GiftAmount

//	for i := 0; i < len(self.gift.info.Mine); i++ {
//		if gift == self.gift.info.Mine[i].Id {
//			if self.gift.info.Mine[i].Amount < amount {
//				self.person.SendErr("礼物数量不足")
//				return
//			} else {
//				self.gift.info.Mine[i].Amount -= amount
//				msg.Amount = self.gift.info.Mine[i].Amount
//				msg.Id = self.gift.info.Mine[i].Id
//			}
//			break
//		}
//	}
//	mod_gift := desPerson.GetModule("gift").(*Mod_Gift)
//	//mod_gift.gift.info.Other = append(mod_gift.gift.info.Other, _msg)
//	findUsr := false //查找送礼物用户
//	for i := 0; i < len(mod_gift.gift.info.Other); i++ {
//		if self.person.Uid == mod_gift.gift.info.Other[i].Uid {
//			findGift := false
//			for j := 0; j < len(mod_gift.gift.info.Other[i].GiftArray); j++ {
//				if mod_gift.gift.info.Other[i].GiftArray[j].Id == gift {
//					findGift = true
//					mod_gift.gift.info.Other[i].GiftArray[j].Amount += amount //增加礼物
//				}
//				break
//			}
//			if !findGift { //没有礼物创建礼物
//				var msg Gift_Array
//				msg.Id = gift
//				msg.Amount = amount
//				mod_gift.gift.info.Other[i].GiftArray = append(mod_gift.gift.info.Other[i].GiftArray, msg)
//			}
//			findUsr = true
//			break
//		}
//	}
//	if !findUsr {
//		var msg Other_Gift
//		msg.Uid = self.person.Uid
//		msg.GiftArray = make([]Gift_Array, 0)
//		var _msg Gift_Array
//		_msg.Id = gift
//		_msg.Amount = amount
//		msg.GiftArray = append(msg.GiftArray, _msg)
//		mod_gift.gift.info.Other = append(mod_gift.gift.info.Other, msg)
//	}
//	self.OnSave(false)
//	mod_gift.OnSave(false)
//	self.person.SendMsg(staticfunc.HF_EncodeMsg("giftupdate", &msg, true))
//	if desPerson.session != nil {
//		var msg C2S_GiftGive
//		msg.Uid = self.person.Uid
//		msg.Type = _type
//		msg.Amount = amount
//		msg.Id = gift
//		desPerson.SendMsg(staticfunc.HF_EncodeMsg("giftaccept", &msg, true))
//	}
//}

//func (self *Mod_Gift) GiftAdd(_type int, gift int64, amount int) {
//	cost := 0
//	if 28 == _type {
//		cost = amount * 500
//	}
//	ok, card, gold := GetServer().CostGold(self.person.Uid, cost, "", 0)
//	if !ok {
//		self.person.SendErr("使用失败")
//		return
//	}
//	if self.person.session != nil {
//		self.person.UpdCard(card, gold) //！更新房卡数量
//	}

//	var msg S2C_GiftAmount
//	for i := 0; i < len(self.gift.info.Mine); i++ {
//		if gift == self.gift.info.Mine[i].Id {
//			self.gift.info.Mine[i].Amount += amount
//			msg.Amount = self.gift.info.Mine[i].Amount
//			msg.Id = self.gift.info.Mine[i].Id
//			break
//		}
//	}
//	self.OnSave(false)
//	self.person.SendMsg(staticfunc.HF_EncodeMsg("giftupdate", &msg, true))
//}

//func (self *Mod_Gift) OnSave(sql bool) {
//	self.gift.Encode()
//	c := GetServer().Redis.Get()
//	defer c.Close()
//	c.Do("SET", fmt.Sprintf("gift_%d", self.person.Uid), self.gift.Value)

//	GetServer().SqlQueue(fmt.Sprintf("update `%s` set `value` = ? where `uid` = '%d'", "gift", self.person.Uid), self.gift.Value, true)
//}
