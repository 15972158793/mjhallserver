package centerserver

import (
	"lib"
	"staticfunc"
	"strings"
)

//! 买金卡
func BuyGold(uid int64, num int, receipt string) bool {
	var msg staticfunc.Msg_GiveCard
	msg.Uid = uid
	msg.Pid = staticfunc.TYPE_CARD
	msg.Num = num
	msg.Ip = receipt
	msg.Dec = 0
	result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "buycard", &msg)

	if err != nil || string(result) == "false" {
		return false
	}

	tmp := strings.Split(string(result), "_")
	if len(tmp) != 2 {
		return false
	}

	now1 := lib.HF_Atoi(tmp[0])
	now2 := lib.HF_Atoi(tmp[1])

	person := GetPersonMgr().GetPerson(int64(uid), false)
	if person != nil {
		var msg staticfunc.S2C_UpdCard
		msg.Card = now1 + now2
		msg.Gold = now2
		person.SendMsg(lib.HF_EncodeMsg("updcard", &msg, true))
	}

	return true
}
