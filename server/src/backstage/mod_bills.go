package backstage

import (
	"fmt"
	//"github.com/garyburd/redigo/redis"
	"lib"
	//"log"
	"net/http"
	"rjmgr"
)

//! 得到代理模式
func GetAgentMode(w http.ResponseWriter, req *http.Request) {
	if !rjmgr.GetRJMgr().IsLicensing() { //! 未授权
		return
	}

	w.Write([]byte(fmt.Sprintf("%d", GetServer().Con.AgentMode)))
}

//! 修改无限代模式
func ModifyBills(w http.ResponseWriter, req *http.Request) {
	if !rjmgr.GetRJMgr().IsLicensing() { //! 未授权
		return
	}

	GetServer().InitBills()
}

//! 佣金兑金币
func CommissionToGold(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))
	money := lib.HF_Atoi(req.FormValue("money"))
	if uid == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	if agent.Commission-agent.T_Commission < money {
		w.Write(HF_Code(3, "佣金不足"))
		return
	}

	agent.T_Commission += money
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `t_commission` = %d where `agid` = %d", agent.T_Commission, agent.Agid), []byte(""), false)
	//GetServer().AddTodayMoney(uid)
	//GetServer().AddCostMoney(money)

	w.Write(HF_Code(0, "成功"))
}

//! 加佣金
func AddCommission(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))
	money := lib.HF_Atoi(req.FormValue("money"))

	if uid == 0 || money == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	agent.Commission += money
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `commission` = %d where `agid` = %d", agent.Commission, agent.Agid), []byte(""), false)
	//GetServer().AddTotalMoney(float64(money))

	w.Write(HF_Code(0, "成功"))
}
