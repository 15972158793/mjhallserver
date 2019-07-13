package backstage

import (
	"fmt"
	"lib"
	"net/http"
)

//! 推广额兑换
func GetMoney(w http.ResponseWriter, req *http.Request) {
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

	//if money != 100 && money != 200 && money != 300 && money != 500 && money != 1000 && money != 2000 {
	//	w.Write(HF_Code(1, "参数错误"))
	//	return
	//}

	if money <= 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	if agent.Score-float64(agent.T_Score) < float64(money) {
		w.Write(HF_Code(3, "可兑换推广额不足"))
		return
	}

	if GetServer().GetTodayMoney(uid) >= 1 {
		w.Write(HF_Code(4, "今日提现次数已达上限"))
		return
	}

	agent.T_Score += money
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `t_score` = %d where `agid` = %d", agent.T_Score, agent.Agid), []byte(""), false)
	GetServer().AddTodayMoney(uid)
	GetServer().AddCostMoney(money)

	w.Write(HF_Code(0, "成功"))
}

//! 推广额兑换金币
func GetMoneyToGold(w http.ResponseWriter, req *http.Request) {
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

	if money <= 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	if agent.Score-float64(agent.T_Score) < float64(money) {
		w.Write(HF_Code(3, "可兑换推广额不足"))
		return
	}

	agent.T_Score += money
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `t_score` = %d where `agid` = %d", agent.T_Score, agent.Agid), []byte(""), false)
	//GetServer().AddTodayMoney(uid)
	GetServer().AddCostMoney(money)

	w.Write(HF_Code(0, "成功"))
}

//! 回退推广额
func ReBackMoney(w http.ResponseWriter, req *http.Request) {
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

	//if money != 100 && money != 200 && money != 300 && money != 500 && money != 1000 && money != 2000 {
	//	w.Write(HF_Code(1, "参数错误"))
	//	return
	//}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	if agent.T_Score < money {
		w.Write(HF_Code(3, "回退失败"))
		return
	}

	agent.T_Score -= money
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `t_score` = %d where `agid` = %d", agent.T_Score, agent.Agid), []byte(""), false)
	GetServer().DelTodayMoney(uid)
	GetServer().AddCostMoney(-money)

	w.Write(HF_Code(0, "成功"))
}

//! 加推广额
func AddMoney(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))
	money := lib.HF_Atof(req.FormValue("money"))

	if uid == 0 || money == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	agent.Score += float64(money)
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `score` = %f where `agid` = %d", agent.Score, agent.Agid), []byte(""), false)
	GetServer().AddTotalMoney(float64(money))

	w.Write(HF_Code(0, "成功"))
}

//! 清除
func ClearMoney(w http.ResponseWriter, req *http.Request) {
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

	if money == 0 {
		GetServer().AddTotalMoney(-agent.Score)
		GetServer().AddCostMoney(-agent.T_Score)
		agent.Score = 0
		agent.T_Score = 0
	} else {
		GetServer().AddTotalMoney(float64(-money))
		GetServer().AddCostMoney(-money)
		agent.Score -= float64(money)
		agent.T_Score -= money
	}
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `score` = %f, `t_score` = %d where `agid` = %d", agent.Score, agent.T_Score, agent.Agid), []byte(""), false)

	w.Write(HF_Code(0, "成功"))
}

//! 手动兑换
func HandMoveMoney(w http.ResponseWriter, req *http.Request) {
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

	if money <= 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	if agent.Score-float64(agent.T_Score) < float64(money) {
		w.Write(HF_Code(3, "可提现推广额不足"))
		return
	}

	agent.T_Score += money
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `t_score` = %d where `agid` = %d", agent.T_Score, agent.Agid), []byte(""), false)
	GetServer().AddTodayMoney(uid)
	GetServer().AddCostMoney(money)

	w.Write(HF_Code(0, "成功"))
}

//! 得到已提现和可提现
func GetCountMoney(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	w.Write([]byte(fmt.Sprintf("%f,%d", GetServer().GetTotalMoney(), GetServer().GetCostMoney())))
}

//! 区域收益兑换金币
func AreaScoreToGold(w http.ResponseWriter, req *http.Request) {
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

	if money <= 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	if int(agent.AreaScore-float64(agent.AreaTScore)) < money {
		w.Write(HF_Code(3, "收益不足"))
		return
	}

	agent.AreaTScore += money
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `areatscore` = %d where `agid` = %d", agent.AreaTScore, agent.Agid), []byte(""), false)

	w.Write(HF_Code(0, "成功"))
}

//! 加区域收益
func AddAreaScore(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))
	money := lib.HF_Atof(req.FormValue("money"))

	if uid == 0 || money == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	agent.AreaScore += float64(money)
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `areascore` = %d where `agid` = %d", agent.AreaScore, agent.Agid), []byte(""), false)

	w.Write(HF_Code(0, "成功"))
}
