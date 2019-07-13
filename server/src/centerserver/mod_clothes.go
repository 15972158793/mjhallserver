package centerserver

import (
//"encoding/json"
//"fmt"
////	"log"
////	"log"
//"staticfunc"

//"github.com/garyburd/redigo/redis"
)

//type Mod_Clothes struct {
//	person  *Person
//	clothes SQL_Clothes
//}
//type SQL_Clothes struct {
//	Uid   int64
//	Value []byte
//	info  JS_Clothes
//}

//type JS_Clothes struct {
//	ClothesInfo []Clothes_Info `json:"clothes"`
//}

//type Clothes_Info struct {
//	Id    int    `json:"id"`
//	Name  string `json:"name"`  //! 名字
//	State int    `json:"state"` //!状态
//	Price int    `json:"price"`
//}

//func (self *SQL_Clothes) Decode() {
//	json.Unmarshal(self.Value, &self.info)
//}

//func (self *SQL_Clothes) Encode() {
//	self.Value = staticfunc.HF_JtoB(self.info)
//}
//func (self *Mod_Clothes) OnGetData(person *Person) {
//	self.person = person

//	c := GetServer().Redis.Get()
//	defer c.Close()
//	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("clothes_%d", self.person.Uid)))
//	if err == nil {
//		self.clothes.Value = v
//		self.clothes.Decode()
//	} else {
//		sql := fmt.Sprintf("select * from `clothes` where uid = %d", self.person.Uid)
//		GetServer().DB.GetOneData(sql, &self.clothes)
//		if self.clothes.Uid <= 0 {
//			self.clothes.Uid = self.person.Uid
//			self.clothes.info.ClothesInfo = make([]Clothes_Info, 5)
//			for i := 0; i < 5; i++ {
//				self.clothes.info.ClothesInfo[i].Id = i
//				self.clothes.info.ClothesInfo[i].Name = "adds"
//				self.clothes.info.ClothesInfo[i].Price = i * 10
//				if i == 0 {
//					self.clothes.info.ClothesInfo[i].State = 2
//				} else {
//					self.clothes.info.ClothesInfo[i].State = 0
//				}
//			}
//			self.clothes.Encode()
//			sql := fmt.Sprintf("insert into `%s`(`uid`, `value`) values (%d, ?)", "clothes", self.person.Uid)
//			GetServer().SqlQueue(sql, self.clothes.Value, true)

//			c := GetServer().Redis.Get()
//			defer c.Close()
//			c.Do("SET", fmt.Sprintf("clothes_%d", self.person.Uid), self.clothes.Value)
//		} else {
//			self.clothes.Decode()
//		}
//	}
//}

//func (self *Mod_Clothes) OnGetOtherData() {

//}

//func (self *Mod_Clothes) SendInfo() {
//	self.person.SendMsg(staticfunc.HF_EncodeMsg("clothes", &self.clothes.info, true))
//}
//func (self *Mod_Clothes) OnMsg(head string, body []byte) bool {
//	switch head {
//	case "clothesGet": //！获取服装列表
//		self.SendInfo()
//		return true
//	case "clothesBuy":
//		var msg C2S_Clothes
//		json.Unmarshal(body, &msg)
//		self.ClothesBuy(msg.Id)
//		return true
//	case "clothesUse":
//		var msg C2S_ClothesUse
//		json.Unmarshal(body, &msg)
//		self.ClothesUse(msg.CurId, msg.UseId)
//		return true
//	case "clothesFind":
//		var msg C2S_ClothesFind
//		json.Unmarshal(body, &msg)
//		self.ClothesFind(msg.Id)
//		return true
//	}
//	return false
//}

//func (self *Mod_Clothes) ClothesBuy(Id int) {
//	var i int
//	for i = 0; i < len(self.clothes.info.ClothesInfo); i++ {
//		if Id == self.clothes.info.ClothesInfo[i].Id {
//			if self.clothes.info.ClothesInfo[i].State != 0 {
//				self.person.SendErr("服装已经购买")
//				return
//			}
//			_, card1, _ := GetServer().CostCard(self.person.Uid, 0, "", 0)
//			if card1 < self.clothes.info.ClothesInfo[i].Price {
//				self.person.SendErr("钻石不足")
//				return
//			}

//			ok, card, gold := GetServer().CostCard(self.person.Uid, self.clothes.info.ClothesInfo[i].Price, "", 0)

//			if !ok {
//				self.person.SendErr("购买失败")
//				return
//			}
//			if self.person.session != nil {
//				self.person.UpdCard(card, gold) //！更新房卡数量
//				self.clothes.info.ClothesInfo[i].State = 1
//			}

//			var msg S2C_ClothesState
//			msg.Id = Id
//			msg.State = self.clothes.info.ClothesInfo[i].State
//			self.OnSave(false)
//			self.person.SendMsg(staticfunc.HF_EncodeMsg("clothesupdate", &msg, true))
//			break
//		}
//	}
//	if i >= len(self.clothes.info.ClothesInfo) {
//		self.person.SendErr("未找到该服装")
//	}
//}

//func (self *Mod_Clothes) ClothesUse(CurId, UseId int) {
//	var i int
//	j := -1
//	for i = 0; i < len(self.clothes.info.ClothesInfo); i++ {
//		if CurId == self.clothes.info.ClothesInfo[i].Id {
//			self.clothes.info.ClothesInfo[i].State = 1
//			j = i
//		}
//	}
//	if j == -1 {
//		self.person.SendErr("未找穿的服装")
//		return
//	}
//	for i = 0; i < len(self.clothes.info.ClothesInfo); i++ {
//		if UseId == self.clothes.info.ClothesInfo[i].Id {
//			if self.clothes.info.ClothesInfo[i].State == 0 {
//				self.person.SendErr("服装未购买")
//				return
//			}

//			self.clothes.info.ClothesInfo[i].State = 2

//			var msg S2C_ClothesState
//			msg.Id = UseId
//			msg.State = self.clothes.info.ClothesInfo[i].State
//			self.OnSave(false)
//			self.person.SendMsg(staticfunc.HF_EncodeMsg("clothesupdate", &msg, true))
//			break
//		}
//	}
//	if i >= len(self.clothes.info.ClothesInfo) {
//		self.person.SendErr("未找到该服装")
//		self.clothes.info.ClothesInfo[j].State = 2
//		return
//	}
//}

//func (self *Mod_Clothes) ClothesFind(uid int64) {
//	desPerson := GetPersonMgr().GetPerson(uid, true)
//	if nil == desPerson {
//		//////////////////////////////////////////
//		self.person.SendErr("玩家不存在")
//		return
//	}
//	mod_clothes := desPerson.GetModule("clothes").(*Mod_Clothes)
//	var msg JS_Clothes
//	msg.ClothesInfo = mod_clothes.clothes.info.ClothesInfo
//	self.person.SendMsg(staticfunc.HF_EncodeMsg("clothesinfo", &msg, true))
//}

//func (self *Mod_Clothes) OnSave(sql bool) {
//	self.clothes.Encode()
//	c := GetServer().Redis.Get()
//	defer c.Close()
//	c.Do("SET", fmt.Sprintf("clothes_%d", self.person.Uid), self.clothes.Value)

//	GetServer().SqlQueue(fmt.Sprintf("update `%s` set `value` = ? where `uid` = '%d'", "clothes", self.person.Uid), self.clothes.Value, true)
//}
