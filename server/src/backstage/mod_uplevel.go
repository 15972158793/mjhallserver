package backstage

import (
	"fmt"
	"lib"
	"log"
	"net/http"
	"strings"
	"time"
)

//! 授权
func UpLevel(w http.ResponseWriter, req *http.Request) {
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
	if uid == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	if destuid == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	if uid != -1 {
		agent := GetServer().GetAgentUser(uid)
		if agent.Agid == 0 {
			w.Write(HF_Code(2, "无法获取用户信息"))
			return
		}

		//if agent.Level < 2 {
		//	w.Write(HF_Code(4, "你没有权利提升权限"))
		//	return
		//}
	}

	destagent := GetServer().GetAgentUser(destuid)
	if destagent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	if destagent.Level > 0 {
		w.Write(HF_Code(0, "成功"))
		return
	}

	if uid != -1 {
		if destagent.Top_Group != "" {
			find := false
			top := strings.Split(destagent.Top_Group, ",")
			for i := 0; i < len(top); i++ {
				if lib.HF_Atoi(top[i]) == uid {
					find = true
					break
				}
			}
			if !find {
				w.Write(HF_Code(3, "只能给自己的下级提升权限"))
				return
			}
		}
	}

	addgold := GetServer().GetTodayActive(destuid)
	destagent.Level = 1
	destagent.Score += addgold
	GetServer().SetAgentUser(destagent)
	GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET `score` = %f, `level` = %d WHERE agid = %d", destagent.Score, destagent.Level, destagent.Agid), []byte(""), false)
	if addgold > 0 {
		if destagent.Parent != 0 {
			GetServer().QueueParentCostLog(fmt.Sprintf("INSERT INTO `parent_cost_log` (`parent`, `score`) VALUES(%d, %f)", destagent.Parent, addgold))
		}
		GetServer().QueueScoreLog(fmt.Sprintf("INSERT INTO `score_log` (operator_id,agid,change_score,save_time,`action`) VALUES(%d, %d, %f, '%s', 1)", destuid, destuid, addgold, time.Now().Format(lib.TIMEFORMAT)))
		GetServer().AddMyMoney(destuid, addgold)
		GetServer().AddTotalMoney(addgold)
	}
	w.Write(HF_Code(0, "成功"))
}

//! 修改返利深度
func SetDeepin(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))
	deepin := lib.HF_Atoi(req.FormValue("deepin"))
	rating := req.FormValue("rating")
	if uid == 0 || rating == "" {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	if deepin <= 0 || deepin > 10 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	if len(strings.Split(rating, ",")) != 10 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	if deepin >= 5 {
		agent.Level = 2
	} else {
		agent.Level = 1
	}
	agent.Deepin = deepin
	agent.Rating = rating
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET `deepin` = %d, `rating` = '%s', `level` = %d WHERE agid = %d", agent.Deepin, agent.Rating, agent.Level, agent.Agid), []byte(""), false)

	if deepin >= 7 {
		ChangeParentFunc(uid, 0)
	}

	w.Write(HF_Code(0, "成功"))
}

//! 修改parent
func SetParent(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))
	parent := lib.HF_Atoi(req.FormValue("parent"))
	log.Println(uid)
	log.Println(parent)
	if uid == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	if parent == 0 { //! 把特殊账号去掉
		agent.Parent = 0
		GetServer().SetAgentUser(agent)
		GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET `parent` = %d  WHERE agid = %d", agent.Parent, agent.Agid), []byte(""), false)
	} else { //! 加入特殊账号
		if agent.Parent == parent {
			w.Write(HF_Code(0, "成功"))
			return
		}
		agent.Parent = parent
		GetServer().SetAgentUser(agent)
		GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET `parent` = %d  WHERE agid = %d", agent.Parent, agent.Agid), []byte(""), false)
		if uid == parent {
			ChangeParentFunc(uid, 0)
		}

		//! 改变所有的下级
		lst := make([]int, 0)
		HF_FindSon(uid, &lst, 0)
		for i := 0; i < len(lst); i++ {
			son := GetServer().GetAgentUser(lst[i])
			if son.Agid == 0 {
				continue
			}
			son.Parent = parent
			GetServer().SetAgentUser(son)
			GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET `parent` = %d WHERE agid = %d", son.Parent, son.Agid), []byte(""), false)
		}
	}

	w.Write(HF_Code(0, "成功"))
}

func SetEveryGold(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	EveryGold = lib.HF_Atoi(req.FormValue("gold"))
	w.Write(HF_Code(0, "成功"))
}

//! 特殊账号的授权
func SetLevel(w http.ResponseWriter, req *http.Request) {
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
	level := lib.HF_Atoi(req.FormValue("level")) //! 0普通玩家  1普通代理  2高级代理   3总代理
	if uid == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	if destuid == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	if level < 0 || level > 3 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	if uid != -1 {
		agent := GetServer().GetAgentUser(uid)
		if agent.Agid == 0 {
			w.Write(HF_Code(2, "无法获取用户信息"))
			return
		}

		if agent.Parent != uid {
			w.Write(HF_Code(4, "你没有权利修改权限"))
			return
		}
	}

	destagent := GetServer().GetAgentUser(destuid)
	if destagent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	if uid != -1 && destagent.Parent != uid {
		w.Write(HF_Code(3, "只能给自己的下级修改权限"))
		return
	}

	setlevel := 0
	setdeepin := 0
	switch level {
	case 0:
		setlevel = 0
		setdeepin = 3
	case 1:
		setlevel = 1
		setdeepin = 3
	case 2:
		setlevel = 1
		setdeepin = 4
	case 3:
		setlevel = 2
		setdeepin = 5
	}

	if destagent.Level == setlevel && destagent.Deepin == setdeepin {
		w.Write(HF_Code(0, "成功"))
		return
	}

	if level == 0 { //! 修改为普通玩家
		if destagent.Level > 0 {
			GetServer().DelTodayActive(destuid)
		}

		destagent.Level = 0
		destagent.Deepin = 3
	} else {
		if destagent.Level == 0 {
			addgold := GetServer().GetTodayActive(destuid)
			destagent.Score += addgold
			if addgold > 0 {
				GetServer().QueueParentCostLog(fmt.Sprintf("INSERT INTO `parent_cost_log` (`parent`, `score`) VALUES(%d, %f)", destagent.Parent, addgold))
				GetServer().QueueScoreLog(fmt.Sprintf("INSERT INTO `score_log` (operator_id,agid,change_score,save_time,`action`) VALUES(%d, %d, %f, '%s', 1)", destuid, destuid, addgold, time.Now().Format(lib.TIMEFORMAT)))
				GetServer().AddMyMoney(destuid, addgold)
				GetServer().AddTotalMoney(addgold)
			}
		}

		if level == 1 {
			destagent.Level = 1
			destagent.Deepin = 3
		} else if level == 2 {
			destagent.Level = 1
			destagent.Deepin = 4
		} else if level == 3 {
			destagent.Level = 2
			destagent.Deepin = 5
		}

		GetServer().SetAgentUser(destagent)
		GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET `score` = %f, `level` = %d, `deepin` = %d WHERE agid = %d", destagent.Score, destagent.Level, destagent.Deepin, destagent.Agid), []byte(""), false)
	}

	w.Write(HF_Code(0, "成功"))
}

//! 授权区域代理
func SetAreaAgent(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))     //! 被授权的人
	pid := lib.HF_Atoi(req.FormValue("pid"))     //! 操作的人
	scale := lib.HF_Atoi(req.FormValue("scale")) //! 比例
	if uid == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	if pid > 0 { //! 操作人的
		if agent.Top_Group == "" {
			w.Write(HF_Code(3, "只能给自己的直属下级授权"))
			return
		} else {
			top := strings.Split(agent.Top_Group, ",")
			if len(top) == 0 || lib.HF_Atoi(top[0]) != pid {
				w.Write(HF_Code(3, "只能给自己的直属下级授权"))
				return
			}
		}
	} else { //! 平台直接授权则要脱离上级
		if scale > 0 {
			ChangeParentFunc(uid, 0)
		}
	}

	agent.AreaScale = scale
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET `areascale` = %d  WHERE agid = %d", agent.AreaScale, agent.Agid), []byte(""), false)

	w.Write(HF_Code(0, "成功"))
}
