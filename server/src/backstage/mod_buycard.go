package backstage

import (
	"fmt"
	"lib"
	"net/http"
)

//! 购买房卡
func BuyCard(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))
	card := lib.HF_Atoi(req.FormValue("card"))
	if uid == 0 || card == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	agent.Card += card
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `card` = %d where `agid` = %d", agent.Card, agent.Agid), []byte(""), false)

	w.Write(HF_Code(0, "成功"))
}

//! 发卡给代理
func SendCardToAgent(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))
	destuid := lib.HF_Atoi(req.FormValue("destuid"))
	card := lib.HF_Atoi(req.FormValue("card"))
	if uid == 0 || destuid == 0 || card == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	if uid == destuid {
		w.Write(HF_Code(4, "无法给自己发卡"))
		return
	}

	var agent *Fa_Agent_User = nil
	if uid != -1 {
		if card < 0 {
			w.Write(HF_Code(1, "参数错误"))
			return
		}

		agent = GetServer().GetAgentUser(uid)
		if agent.Agid == 0 {
			w.Write(HF_Code(2, "无法获取用户信息"))
			return
		}

		if agent.Card < card {
			w.Write(HF_Code(3, "卡不足"))
			return
		}
	}

	destagent := GetServer().GetAgentUser(destuid)
	if destagent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	if agent != nil {
		agent.Card -= card
		GetServer().SetAgentUser(agent)
		GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `card` = %d where `agid` = %d", agent.Card, agent.Agid), []byte(""), false)
	}
	destagent.Card = lib.HF_MaxInt(0, destagent.Card+card)
	GetServer().SetAgentUser(destagent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `card` = %d where `agid` = %d", destagent.Card, destagent.Agid), []byte(""), false)

	w.Write(HF_Code(0, "成功"))
}

//! 发卡给玩家
func SendCardToPlayer(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))
	destuid := lib.HF_Atoi(req.FormValue("destuid"))
	card := lib.HF_Atoi(req.FormValue("card"))
	if uid == 0 || destuid == 0 || card == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	if agent.Card < card {
		w.Write(HF_Code(3, "卡不足"))
		return
	}

	//! 发卡
	str := fmt.Sprintf("http://%s:%d/givecard?uid=%d&productid=%d&product_count=%d", GetServer().Con.GameIP, GetServer().Con.CenterPort, destuid, 1, card)
	res, err := lib.HF_Get(str, 3)
	if err != nil {
		w.Write(HF_Code(4, "游戏服务器无响应"))
		return
	}
	res.Body.Close()

	agent.Card -= card
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `card` = %d where `agid` = %d", agent.Card, agent.Agid), []byte(""), false)

	w.Write(HF_Code(0, "成功"))
}
