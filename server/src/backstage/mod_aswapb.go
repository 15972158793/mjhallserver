package backstage

import (
	"fmt"
	"io/ioutil"
	"lib"
	"log"
	"net/http"
	"rjmgr"
	"strings"
	"sync"
	"time"
)

//! a扫码b模块

var SwapLock *sync.RWMutex = new(sync.RWMutex)

func ABindB(w http.ResponseWriter, req *http.Request) {
	if !rjmgr.GetRJMgr().IsLicensing() { //! 未授权
		w.Write([]byte("false"))
		return
	}

	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		w.Write([]byte("false"))
		return
	}

	if GetServer().ShutDown {
		w.Write([]byte("false"))
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	SwapLock.Lock()
	defer SwapLock.Unlock()

	auid := lib.HF_Atoi(req.FormValue("auid"))
	aunionid := req.FormValue("aunionid")
	aopenid := req.FormValue("aopenid")
	buid := lib.HF_Atoi(req.FormValue("buid"))
	lib.GetLogMgr().Output(lib.LOG_DEBUG, auid, ".", aunionid, ".", aopenid, ".", buid)
	if auid == 0 || buid == 0 {
		w.Write([]byte("false"))
		log.Println("参数错误")
		return
	}

	if auid == buid {
		w.Write([]byte("false"))
		log.Println("无法绑定自己")
		return
	}

	bagent := GetServer().GetAgentUser(buid)
	if bagent.Agid == 0 {
		w.Write([]byte("true"))
		log.Println("b不存在")
		return
	}

	btop := make([]string, 0)
	if bagent.Top_Group != "" {
		btop = strings.Split(bagent.Top_Group, ",")
		for i := 0; i < len(btop); i++ {
			if lib.HF_Atoi(btop[i]) == auid {
				w.Write([]byte("true"))
				log.Println("扫码者在被扫码者的上级")
				return
			}
		}
	}

	abind := GetServer().GetBindPlayers(auid)
	if abind.Uid != 0 {
		w.Write([]byte("true"))
		log.Println("该用户已经绑定了上级")
		return
	}

	//! 更新用户表
	aagent := GetServer().GetAgentUser(auid)
	//if aagent.Deepin >= 7 { //! 这个深度的人不能绑定上级
	//	w.Write(HF_Code(-7, "绑定失败，高级代理无法绑定上级"))
	//	log.Println("-7")
	//	return
	//}

	if aagent.AreaScale != 0 {
		w.Write([]byte("true"))
		log.Println("高级代理无法绑定上级")
		return
	}

	if aagent.Parent == auid {
		w.Write([]byte("true"))
		log.Println("高级代理无法绑定上级")
		return
	}

	//! 插入绑定表
	GetServer().InsertBindPlayers(auid, buid)

	//!
	aagent.Parent = bagent.Parent

	//! 修改数据表
	if len(btop) == 0 {
		aagent.Top_Group = fmt.Sprintf("%d", buid)
	} else {
		aagent.Top_Group = fmt.Sprintf("%d", buid)
		for i := 0; i < len(btop); i++ {
			aagent.Top_Group += ","
			aagent.Top_Group += btop[i]
		}
	}
	if aagent.Agid == 0 { //! 用户表里没有
		aagent.Agid = auid
		GetServer().InsertAgentUser(aagent.Agid, aopenid, aunionid, aagent.Top_Group, aagent.Parent)
	} else {
		GetServer().SetAgentUser(aagent)
		GetServer().QueueAgentUser(fmt.Sprintf("update fa_agent_user set `top_group` = '%s', `parent` = %d where `agid` = %d", aagent.Top_Group, aagent.Parent, aagent.Agid), []byte(""), false)
	}

	atop := strings.Split(aagent.Top_Group, ",")

	//! 修改A的所有下级
	lst := make([]int, 0)
	HF_FindSon(aagent.Agid, &lst, 0)
	for i := 0; i < len(lst); i++ {
		nagent := GetServer().GetAgentUser(lst[i])
		if nagent.Agid == 0 {
			continue
		}

		if nagent.Top_Group == "" {
			continue
		}

		for j := 0; j < len(atop); j++ {
			nagent.Top_Group += ","
			nagent.Top_Group += atop[j]
		}
		GetServer().SetAgentUser(nagent)
		GetServer().QueueAgentUser(fmt.Sprintf("update fa_agent_user set `top_group` = '%s' where `agid` = %d", nagent.Top_Group, nagent.Agid), []byte(""), false)
	}

	//! 新增
	GetServer().AddTodayChild(buid)
	GetServer().AddMonthChild(buid)
	GetServer().AddWeekChild(buid)

	if bagent.Level == 0 { //! b是0级,判断是否升级
		addgold := GetServer().GetTodayActive(buid)
		if addgold >= 15 {
			if len(GetServer().GetChildren(buid).Child) >= 5 {
				bagent.Level = 1
				bagent.Score += addgold
				GetServer().SetAgentUser(bagent)
				GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET score = %f, `level` = %d WHERE agid = %d", bagent.Score, bagent.Level, buid), []byte(""), false)
				GetServer().QueueScoreLog(fmt.Sprintf("INSERT INTO `score_log` (operator_id,agid,change_score,save_time,`action`) VALUES(%d, %d, %f, '%s', 1)", buid, buid, addgold, time.Now().Format(lib.TIMEFORMAT)))
				GetServer().AddMyMoney(buid, addgold)
				GetServer().AddTotalMoney(addgold)
			}
		}
	}

	log.Println("绑定成功")
	w.Write([]byte("true"))
}

func ASwapB(w http.ResponseWriter, req *http.Request) {
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

	SwapLock.Lock()
	defer SwapLock.Unlock()

	aunionid := req.FormValue("aunionid")
	aopenid := req.FormValue("aopenid")
	aip := req.FormValue("ip")
	buid := lib.HF_Atoi(req.FormValue("buid"))
	log.Println("b..", buid)
	if buid == 0 {
		w.Write(HF_Code(-98, "参数错误"))
		log.Println("-98")
		return
	}

	//! 向游戏服务器获取uid
	str := fmt.Sprintf("http://%s:%d/wxtouid?unionid=%s&openid=%s&ip=%s&buid=%d", GetServer().Con.GameIP, GetServer().Con.LoginPort, aunionid, aopenid, aip, buid)
	res, err := lib.HF_Get(str, 3)
	if err != nil {
		w.Write(HF_Code(-1, "服务器无响应"))
		log.Println("-1")
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		w.Write(HF_Code(-2, "游戏服务器返回错误"))
		log.Println("-2")
		return
	}

	if aunionid == "" || aopenid == "" {
		w.Write(HF_Code(0, "绑定成功"))
		log.Println("0")
		return
	}

	auid := lib.HF_Atoi(string(body))
	log.Println("a..", auid)

	if auid == buid {
		w.Write(HF_Code(-6, "无法绑定自己"))
		log.Println("-6")
		return
	}

	bagent := GetServer().GetAgentUser(buid)
	if bagent.Agid == 0 {
		w.Write(HF_Code(-3, "绑定失败，被扫者不存在"))
		log.Println("-3")
		return
	}

	btop := make([]string, 0)
	if bagent.Top_Group != "" {
		btop = strings.Split(bagent.Top_Group, ",")
		for i := 0; i < len(btop); i++ {
			if lib.HF_Atoi(btop[i]) == auid {
				w.Write(HF_Code(-5, "绑定失败，扫码者在被扫码者的上级"))
				log.Println("-5")
				return
			}
		}
	}

	abind := GetServer().GetBindPlayers(auid)
	if abind.Uid != 0 {
		w.Write(HF_Code(-4, "绑定失败，该用户已经绑定了上级"))
		log.Println("-4")
		return
	}

	//! 更新用户表
	aagent := GetServer().GetAgentUser(auid)
	//if aagent.Deepin >= 7 { //! 这个深度的人不能绑定上级
	//	w.Write(HF_Code(-7, "绑定失败，高级代理无法绑定上级"))
	//	log.Println("-7")
	//	return
	//}

	if aagent.AreaScale != 0 {
		w.Write(HF_Code(-7, "绑定失败，高级代理无法绑定上级"))
		log.Println("-7")
		return
	}

	log.Println(aagent.Parent, "...", aagent.Agid)
	if aagent.Parent == auid {
		w.Write(HF_Code(-7, "绑定失败，高级代理无法绑定上级"))
		log.Println("-7")
		return
	}

	//! 插入绑定表
	GetServer().InsertBindPlayers(auid, buid)

	//!
	aagent.Parent = bagent.Parent

	//! 修改数据表
	if len(btop) == 0 {
		aagent.Top_Group = fmt.Sprintf("%d", buid)
	} else {
		aagent.Top_Group = fmt.Sprintf("%d", buid)
		for i := 0; i < len(btop); i++ {
			aagent.Top_Group += ","
			aagent.Top_Group += btop[i]
		}
	}
	if aagent.Agid == 0 { //! 用户表里没有
		aagent.Agid = auid
		GetServer().InsertAgentUser(aagent.Agid, aopenid, aunionid, aagent.Top_Group, aagent.Parent)
	} else {
		GetServer().SetAgentUser(aagent)
		GetServer().QueueAgentUser(fmt.Sprintf("update fa_agent_user set `top_group` = '%s', `parent` = %d where `agid` = %d", aagent.Top_Group, aagent.Parent, aagent.Agid), []byte(""), false)
	}

	atop := strings.Split(aagent.Top_Group, ",")

	//! 修改A的所有下级
	lst := make([]int, 0)
	HF_FindSon(aagent.Agid, &lst, 0)
	for i := 0; i < len(lst); i++ {
		nagent := GetServer().GetAgentUser(lst[i])
		if nagent.Agid == 0 {
			continue
		}

		if nagent.Top_Group == "" {
			continue
		}

		for j := 0; j < len(atop); j++ {
			nagent.Top_Group += ","
			nagent.Top_Group += atop[j]
		}
		GetServer().SetAgentUser(nagent)
		GetServer().QueueAgentUser(fmt.Sprintf("update fa_agent_user set `top_group` = '%s' where `agid` = %d", nagent.Top_Group, nagent.Agid), []byte(""), false)
	}

	//! 新增
	GetServer().AddTodayChild(buid)
	GetServer().AddMonthChild(buid)
	GetServer().AddWeekChild(buid)

	if bagent.Level == 0 { //! b是0级,判断是否升级
		addgold := GetServer().GetTodayActive(buid)
		if addgold >= 15 {
			if len(GetServer().GetChildren(buid).Child) >= 5 {
				bagent.Level = 1
				bagent.Score += addgold
				GetServer().SetAgentUser(bagent)
				GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET score = %f, `level` = %d WHERE agid = %d", bagent.Score, bagent.Level, buid), []byte(""), false)
				GetServer().QueueScoreLog(fmt.Sprintf("INSERT INTO `score_log` (operator_id,agid,change_score,save_time,`action`) VALUES(%d, %d, %f, '%s', 1)", buid, buid, addgold, time.Now().Format(lib.TIMEFORMAT)))
				GetServer().AddMyMoney(buid, addgold)
				GetServer().AddTotalMoney(addgold)
			}
		}
	}

	w.Write(HF_Code(0, "绑定成功"))
	log.Println("0")
}
