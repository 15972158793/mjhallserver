package backstage

import (
	"fmt"
	//"io/ioutil"
	//"log"
	"lib"
	"net/http"
	"strings"
	"sync"
)

//! 更换上级

var ChangeParentLock *sync.RWMutex = new(sync.RWMutex)

func ChangeParent(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	sonid := lib.HF_Atoi(req.FormValue("sonid"))
	parentid := lib.HF_Atoi(req.FormValue("parentid"))

	w.Write(HF_Code(ChangeParentFunc(sonid, parentid)))
}

func ChangeParentFunc(sonid int, parentid int) (int, string) {
	ChangeParentLock.Lock()
	defer ChangeParentLock.Unlock()

	if sonid == 0 {
		return 1, "参数错误"
	}

	if sonid == parentid {
		return 2, "不能移动到自己下级"
	}

	sondata := GetServer().GetAgentUser(sonid)
	if sondata.Agid == 0 {
		return 1, "参数错误"
	}

	parentdata := new(Fa_Agent_User)
	if parentid != 0 {
		parentdata = GetServer().GetAgentUser(parentid)
		if parentdata.Agid == 0 {
			return 1, "参数错误"
		}

		top := make([]string, 0)
		if parentdata.Top_Group != "" {
			top = strings.Split(parentdata.Top_Group, ",")
		}
		for i := 0; i < len(top); i++ {
			if lib.HF_Atoi(top[i]) == sonid {
				return 3, "移动失败，上级玩家在你的下级"
			}
		}
	}

	if parentid == 0 { //! 删除上级
		oldparentid := 0
		if sondata.Top_Group == "" {
			return 0, "移动成功"
		} else {
			top := strings.Split(sondata.Top_Group, ",")
			oldparentid = lib.HF_Atoi(top[0])
		}
		GetServer().DelBindPlayers(sonid)

		//! 更新自己
		sondata.Top_Group = ""
		GetServer().SetAgentUser(sondata)
		GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET `top_group` = '' WHERE agid = %d", sonid), []byte(""), false)

		//! 更新下级
		lst := make([]int, 0)
		HF_FindSon(sonid, &lst, 0)
		for i := 0; i < len(lst); i++ {
			son := GetServer().GetAgentUser(lst[i])
			if son.Agid == 0 {
				continue
			}
			top := make([]string, 0)
			if son.Top_Group != "" {
				top = strings.Split(son.Top_Group, ",")
			}
			for i := 0; i < len(top); i++ {
				if lib.HF_Atoi(top[i]) == oldparentid {
					top = top[0:i]
					break
				}
			}
			son.Top_Group = ""
			for i := 0; i < len(top); i++ {
				if i == 0 {
					son.Top_Group += top[i]
				} else {
					son.Top_Group += ","
					son.Top_Group += top[i]
				}
			}
			GetServer().SetAgentUser(son)
			GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET `top_group` = '%s' WHERE agid = %d", son.Top_Group, son.Agid), []byte(""), false)
		}
	} else { //! 移动上级
		oldparentid := 0
		if sondata.Top_Group == "" { //! 本来没有上级
			GetServer().InsertBindPlayers(sonid, parentid)
		} else {
			top := strings.Split(sondata.Top_Group, ",")
			if len(top) > 0 && lib.HF_Atoi(top[0]) == parentid { //! 上级没有发生改变
				return 0, "移动成功"
			}
			oldparentid = lib.HF_Atoi(top[0])
			GetServer().SetBindPlayers(sonid, parentid)
		}

		//! 更新下线
		lst := []int{sonid}
		HF_FindSon(sonid, &lst, 0)
		for i := 0; i < len(lst); i++ {
			son := GetServer().GetAgentUser(lst[i])
			if son.Agid == 0 {
				continue
			}
			son.Parent = parentdata.Parent
			top := make([]string, 0)
			if son.Top_Group != "" {
				top = strings.Split(son.Top_Group, ",")
			}
			for i := 0; i < len(top); i++ {
				if lib.HF_Atoi(top[i]) == oldparentid {
					top = top[0:i]
					break
				}
			}
			top = append(top, fmt.Sprintf("%d", parentid))
			if parentdata.Top_Group != "" {
				newtop := strings.Split(parentdata.Top_Group, ",")
				top = append(top, newtop...)
			}
			//if len(top) > 10 {
			//	top = top[0:10]
			//}
			son.Top_Group = ""
			for i := 0; i < len(top); i++ {
				if i == 0 {
					son.Top_Group += top[i]
				} else {
					son.Top_Group += ","
					son.Top_Group += top[i]
				}
			}
			GetServer().SetAgentUser(son)
			GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET `top_group` = '%s', `parent` = %d WHERE agid = %d", son.Top_Group, son.Parent, son.Agid), []byte(""), false)
		}
	}

	return 0, "移动成功"
}
