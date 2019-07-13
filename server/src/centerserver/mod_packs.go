package centerserver

import (
//"log"
//"encoding/json"
//"fmt"
//	"log"
//"staticfunc"

//"github.com/garyburd/redigo/redis"
)

//type Mod_Packs struct {
//	person *Person
//	packs  SQL_Packs
//}
//type SQL_Packs struct {
//	Uid   int64
//	State int
//}

//func (self *Mod_Packs) OnGetData(person *Person) {
//	self.person = person

//	sql := fmt.Sprintf("select * from `packs` where uid = %d", self.person.Uid)
//	GetServer().DB.GetOneData(sql, &self.packs)
//	if self.packs.Uid <= 0 {
//		self.packs.Uid = self.person.Uid
//		self.packs.State = 1

//		sql := fmt.Sprintf("insert into `%s`(`uid`, `value`) values (%d, %d)", "packs", self.person.Uid, self.packs.State)
//		GetServer().SqlQueue(sql, []byte(""), false)
//	}
//}

//func (self *Mod_Packs) OnGetOtherData() {

//}

//func (self *Mod_Packs) OnMsg(head string, body []byte) bool {
//	switch head {
//	case "packsGet": //！获取新手礼包
//		self.PackGets()
//		return true
//	}
//	return false
//}

//func (self *Mod_Packs) PackGets() {
//	if self.packs.State == 1 {
//		okgold, _, gold := GetServer().AddGold(self.person.Uid, 50000, "", 0)
//		if !okgold {
//			self.person.SendErr("加金币失败")
//			return
//		}
//		okcard, card, _ := GetServer().AddCard(self.person.Uid, 100, "", 0)
//		if !okcard {
//			self.person.SendErr("加钻石失败")
//			return
//		}
//		if self.person.session != nil {
//			self.person.UpdCard(card, gold) //！更新房卡数量
//		}

//		self.packs.State = 0
//		var msg S2C_PacksState
//		msg.State = self.packs.State
//		msg.Uid = self.person.Uid
//		//self.person.SendErr("领取成功")
//		self.OnSave(false)
//		self.person.SendMsg(staticfunc.HF_EncodeMsg("packsget", &msg, true))
//	} else {
//		self.person.SendErr("已经领取过了")
//	}
//}

//func (self *Mod_Packs) OnSave(sql bool) {
//	c := GetServer().Redis.Get()
//	defer c.Close()
//	c.Do("SET", fmt.Sprintf("packs_%d", self.person.Uid), self.packs.State)

//	GetServer().SqlQueue(fmt.Sprintf("update `%s` set `value` = %d where `uid` = '%d'", "packs", self.packs.State, self.person.Uid), []byte(""), false)
//}
