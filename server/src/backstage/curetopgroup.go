package backstage

import (
	"fmt"
	"log"
)

//! 上级代理链表
type TopPerson struct {
	Uid    int64      //! 自己的id
	Parent *TopPerson //! 上级
}

var mapTopPerson map[int64]*TopPerson = make(map[int64]*TopPerson)

type SQL_Agent struct {
	Id        int
	Agid      int64
	Top_group string
}

//! 修复一下topgroup
func CureTopGroup() {
	maxid := 0
	for {
		var agentnode SQL_Agent
		if !GetServer().HT_DB.GetOneData(fmt.Sprintf("select `id`, `agid`, `top_group` from `fa_agent_user` where `id` > %d order by id asc limit 1", maxid), &agentnode) {
			log.Println("出错id:", maxid)
			return
		}

		if agentnode.Agid == 0 {
			log.Println("终止id:", maxid)
			return
		}

		old := agentnode.Top_group

		my := GetTopPerson(agentnode.Agid)
		agentnode.Top_group = ""

		deep := 0
		parent := my.Parent
		if parent != nil {
			agentnode.Top_group = fmt.Sprintf("%d", parent.Uid)
			parent = parent.Parent
			deep = 1
			for parent != nil {
				agentnode.Top_group += ","
				agentnode.Top_group += fmt.Sprintf("%d", parent.Uid)
				parent = parent.Parent
				deep++
				if deep >= 10 {
					break
				}
			}
		}
		if agentnode.Top_group != old {
			GetServer().HT_DB.Exec("update fa_agent_user set `top_group` = '%s' where `agid` = %d", agentnode.Top_group, agentnode.Agid)
		}

		maxid = agentnode.Id
		log.Println(agentnode.Id, "修复完毕")
	}
}

func GetTopPerson(uid int64) *TopPerson {
	value, ok := mapTopPerson[uid]
	if ok {
		return value
	}

	my := &TopPerson{uid, nil}
	AddTopPerson(my)

	return my
}

type SQL_BindPlayers struct {
	Agid int64
}

func AddTopPerson(my *TopPerson) {
	mapTopPerson[my.Uid] = my
	var agentnode SQL_BindPlayers
	GetServer().HT_DB.GetOneData(fmt.Sprintf("select `agid` from `fa_bind_players` where `uid` = %d", my.Uid), &agentnode)
	if agentnode.Agid != 0 {
		value, ok := mapTopPerson[agentnode.Agid]
		if ok {
			my.Parent = value
			return
		}
		_my := &TopPerson{agentnode.Agid, nil}
		my.Parent = _my
		AddTopPerson(_my)
	}
}
