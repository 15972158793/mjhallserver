package backstage

import (
	"fmt"
	"io/ioutil"
	"lib"
	"log"
	"net/http"
	"rjmgr"
	"strings"
)

//! 得到信息
func GetInfo(w http.ResponseWriter, req *http.Request) {
	if !rjmgr.GetRJMgr().IsLicensing() { //! 未授权
		return
	}

	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))
	unionid := req.FormValue("unionid")
	openid := req.FormValue("openid")
	log.Println(uid, "...", unionid, "...", openid)
	if uid == 0 || unionid == "" || openid == "" {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		if unionid == "jyqp" {
			w.Write(HF_Code(2, "无法获取用户信息"))
			return
		}
		GetServer().InsertAgentUser(uid, openid, unionid, "", 0)
		agent = GetServer().GetAgentUser(uid)
		if agent.Agid == 0 {
			w.Write(HF_Code(2, "无法获取用户信息"))
			return
		}
	}

	name, head, card, gold, ok := HF_GetPlayer(uid, 5)
	if !ok {
		w.Write(HF_Code(2, "游戏服务器无响应"))
		return
	}

	if agent.NickName != name || agent.Head != head {
		agent.NickName = name
		agent.Head = head
		GetServer().SetAgentUser(agent)
		GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `nickname` = ?, `head` = '%s' where `agid` = %d", head, agent.Agid), []byte(name), true)
	}

	bind := GetServer().GetBindPlayers(uid)

	var node HT_Player
	node.Uid = uid
	node.NickName = agent.NickName
	node.Head = agent.Head
	node.Level = agent.Level
	node.Deepin = agent.Deepin
	node.Card = card
	node.Gold = gold
	node.Score = agent.Score
	node.T_Score = agent.T_Score
	node.AgentCard = agent.Card
	node.Rating = agent.Rating
	node.Top_Group = agent.Top_Group
	node.PassWord = agent.Password
	node.Union_id = agent.Union_Id
	node.ChildNum = len(GetServer().GetChildren(uid).Child)
	node.Score1 = bind.Score1
	node.Score2 = bind.Score2
	node.CreateTime = bind.Bind_time
	node.Parent = agent.Parent
	node.Name = agent.Name
	node.Alipay = agent.Alipay
	node.AliName = agent.AliName
	node.Bankcard = agent.Bankcard
	node.Bankname = agent.Bankname
	node.Phone = agent.Phone
	node.AllCost = agent.AllCost
	node.AllBills = agent.AllBills
	node.DayBills = agent.DayBills
	node.WeekBills = agent.WeekBills
	node.MonthBills = agent.MonthBills
	if agent.AddCommission(0, true) {
		GetServer().SetAgentUser(agent)
		GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET `commission` = %d, `bills1` = %d, `bills2` = %d, `timecommission` = %d WHERE agid = %d", agent.Commission, agent.Bills1, agent.Bills2, agent.TimeCommission, agent.Agid), []byte(""), false)
	}
	node.Bills1 = agent.Bills1             //! 直属业绩
	node.Bills2 = agent.Bills2             //! 推广员业绩
	node.Commission = agent.Commission     //! 总佣金
	node.T_Commission = agent.T_Commission //! 已提佣金
	//! 直属佣金
	total := node.Bills1 + node.Bills2
	tlevel := GetServer().GetBillsLevel(total)
	node.Commission1 = int((float64(node.Bills1) / 1000000.0) * float64(tlevel))
	//! 下级产生佣金
	xlevel := GetServer().GetBillsLevel(node.Bills2)
	node.Commission2 = int((float64(node.Bills2) / 1000000.0) * float64(tlevel-xlevel))
	node.TodayAdd = GetServer().GetTodayChild(uid)
	node.MonthAdd = GetServer().GetMonthChild(uid)
	node.WeekAdd = GetServer().GetWeekChild(uid)
	node.AgentMode = GetServer().Con.AgentMode
	node.AreaScale = agent.AreaScale
	node.AreaScore = agent.AreaScore
	node.AreaTScore = agent.AreaTScore

	str := fmt.Sprintf("http://%s:%d/getmoneymode", GetServer().Con.GameIP, GetServer().Con.LoginPort)
	res, err := lib.HF_Get(str, 3)
	if err == nil {
		body, _ := ioutil.ReadAll(res.Body)
		node.MoneyMode = lib.HF_Atoi(string(body))
	}
	res.Body.Close()

	w.Write(lib.HF_JtoB(&node))
}

//! 得到基本信息
//! 只能搜索自己的下级
func GetBaseInfo(w http.ResponseWriter, req *http.Request) {
	if !rjmgr.GetRJMgr().IsLicensing() { //! 未授权
		return
	}

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
	if uid == 0 || destuid == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	if agent.Level < 2 { //! 总代理以下，只能看自己的直属下级
		all := GetServer().GetChildren(agent.Agid)
		if !all.HasChild(destuid) {
			w.Write(HF_Code(3, "该用户不在你的下级"))
			return
		}
	}

	destagent := GetServer().GetAgentUser(destuid)
	if destagent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	if agent.Level >= 2 {
		find := false
		if destagent.Top_Group != "" {
			top := strings.Split(destagent.Top_Group, ",")
			for i := 0; i < len(top); i++ {
				if lib.HF_Atoi(top[i]) == uid {
					find = true
					break
				}
			}
		}
		if !find {
			w.Write(HF_Code(3, "该用户不在你的下级"))
			return
		}
	}

	name, head, _, _, ok := HF_GetPlayer(destuid, 5)
	if !ok {
		w.Write(HF_Code(2, "游戏服务器无响应"))
		return
	}

	if destagent.NickName != name || destagent.Head != head {
		destagent.NickName = name
		destagent.Head = head
		GetServer().SetAgentUser(destagent)
		GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `nickname` = ?, `head` = '%s' where `agid` = %d", head, destagent.Agid), []byte(name), true)
	}

	bind := GetServer().GetBindPlayers(destuid)

	var node HT_BasePlayer
	node.Uid = destuid
	node.NickName = destagent.NickName
	node.Head = destagent.Head
	node.Level = destagent.Level
	node.Deepin = destagent.Deepin
	node.AgentCard = destagent.Card
	node.Score1 = bind.Score1
	node.Score2 = bind.Score2
	node.CreateTime = bind.Bind_time
	w.Write(lib.HF_JtoB(&node))
}
