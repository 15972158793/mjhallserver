package backstage

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"lib"
	"log"
	"runtime/debug"
	"strings"
	"time"
)

var GOLDBL float64 = 100.0

type AgentGoldMgr struct {
}

var agentgoldSingleton *AgentGoldMgr = nil

//! 得到服务器指针
func GetAgentGoldMgr() *AgentGoldMgr {
	if agentgoldSingleton == nil {
		agentgoldSingleton = new(AgentGoldMgr)
	}

	return agentgoldSingleton
}

type SQL_Score_Num struct {
	Score_num int64
}

//! 得到当前读取到第几条
func (self *AgentGoldMgr) GetCount() int64 {
	var sql SQL_Score_Num
	if !GetServer().HT_DB.GetOneData("SELECT `score_num` FROM `fa_step`", &sql) {
		log.Fatal("get score_num fail")
		return 0
	}
	return sql.Score_num
}

func (self *AgentGoldMgr) SetCount(id int64) {
	a, b := GetServer().HT_DB.Exec(fmt.Sprintf("UPDATE `fa_step` SET score_num = %d where id = 1", id))
	if a == 0 && b == 0 {
		log.Fatal("SetCount fail")
	}
}

type SQL_AgentGold struct {
	Id   int64
	Uid  int
	Gold int
}

var EveryGold int = 0

//! 开始计算代理返利
func (self *AgentGoldMgr) Do() {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	if GetServer().Con.AgentMode == 2 { //! 新大区模式不读这个表
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	lib.GetLogMgr().Output(lib.LOG_ERROR, "开始统计")
	maxid := self.GetCount()

	agentgold := make(map[int]int)
	var sql SQL_AgentGold
	res := GetServer().Game_DB.GetAllData(fmt.Sprintf("SELECT id,uid,gold FROM log_agentgold WHERE id > %d ORDER BY id ASC", maxid), &sql)
	if len(res) == 0 {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "结束统计")
		return
	}
	log.Println("读取了:", len(res))
	maxid = res[len(res)-1].(*SQL_AgentGold).Id
	self.SetCount(maxid)
	for i := 0; i < len(res); i++ {
		if res[i].(*SQL_AgentGold).Gold <= 0 { //! 只统计正的
			continue
		}
		agentgold[res[i].(*SQL_AgentGold).Uid] += res[i].(*SQL_AgentGold).Gold
	}
	log.Println("合并后有:", len(agentgold))

	for key, value := range agentgold {
		data := GetServer().GetAgentUser(key)
		if data.Agid == 0 {
			continue
		}

		data.AllCost += int64(value)
		GetServer().SetAgentUser(data)
		GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET `allcost` = %d WHERE agid = %d", data.AllCost, key), []byte(""), false)

		if data.Parent != 0 { //! 一个特殊账号的下级
			parent := GetServer().GetAgentUser(data.Parent)
			if parent.Parent != data.Parent { //! 资格被取消
				data.Parent = 0
				GetServer().SetAgentUser(data)
				GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET `parent` = %d WHERE agid = %d", data.Parent, data.Agid), []byte(""), false)
			} else {
				GetServer().QueueParentLog(fmt.Sprintf("INSERT INTO `parent_log` (`uid`, `gold`, `parent`, `time`) VALUES(%d, %d, %d, %d)", key, value, data.Parent, time.Now().Unix()))
			}
		}

		top := make([]string, 0)
		if data.Top_Group != "" {
			top = strings.Split(data.Top_Group, ",")
		}

		if data.Level > 0 {
			top = append([]string{fmt.Sprintf("%d", key)}, top...)
		}

		if len(top) == 0 {
			continue
		}

		score_total := float64(0)
		theno := 0
		for i := 0; i < len(top); i++ {
			topid := lib.HF_Atoi(top[i])
			if topid != key {
				theno++
			}
			if topid == 0 {
				continue
			}
			agenttop := GetServer().GetAgentUser(topid)
			if agenttop.Agid == 0 {
				continue
			}
			if agenttop.Deepin <= i {
				continue
			}
			if agenttop.Level <= 0 {
				if i == 0 || i == 1 { //! 模拟返利深度
					willgold := float64(value*35) / 10000.0
					if i == 1 {
						willgold = float64(value*15) / 10000.0
					}
					addgold := GetServer().AddTodayActive(topid, willgold)
					if addgold >= 15 { //! 活跃度达标
						if len(GetServer().GetChildren(topid).Child) >= 5 {
							score_total += addgold
							agenttop.Level = 1
							agenttop.Score += addgold
							GetServer().SetAgentUser(agenttop)
							GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET score = %f, `level` = %d WHERE agid = %d", agenttop.Score, agenttop.Level, topid), []byte(""), false)
							GetServer().QueueScoreLog(fmt.Sprintf("INSERT INTO `score_log` (operator_id,agid,change_score,save_time,`action`) VALUES(%d, %d, %f, '%s', 1)", topid, topid, addgold, time.Now().Format(lib.TIMEFORMAT)))
							GetServer().AddMyMoney(topid, addgold)
							GetServer().AddTotalMoney(addgold)
						}
					}
				}

				continue
			}
			fanli := 0
			rating := strings.Split(agenttop.Rating, ",")
			if i >= len(rating) {
				continue //! 深度太大就放弃
			} else {
				fanli = lib.HF_Atoi(rating[i])
				if fanli <= 0 {
					continue
				}
			}
			addgold := float64(value*fanli) / GOLDBL / 100.0

			score_total += addgold
			agenttop.Score += addgold
			GetServer().SetAgentUser(agenttop)
			GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET score = %f, `level` = %d WHERE agid = %d", agenttop.Score, agenttop.Level, topid), []byte(""), false)
			GetServer().QueueScoreLog(fmt.Sprintf("INSERT INTO `score_log` (operator_id,agid,change_score,save_time,`action`) VALUES(%d, %d, %f, '%s', 1)", key, topid, addgold, time.Now().Format(lib.TIMEFORMAT)))
			GetServer().AddTotalMoney(addgold)
			if key == topid {
				GetServer().AddMyMoney(topid, addgold)
			}
			if theno == 1 { //! 直属上级
				GetServer().AddBindPlayers(key, addgold, 0)

				c := GetServer().Redis.Get()
				goldkey := fmt.Sprintf("%s_%d_%d_%d_%d", "score1", key, time.Now().Year(), time.Now().Month(), time.Now().Day())
				thegold, err := redis.Float64(c.Do("GET", goldkey))
				if err != nil {
					thegold = 0
				}
				thegold += addgold
				c.Do("SET", goldkey, thegold)
				c.Do("EXPIRE", goldkey, 86400*3)
				c.Close()
			} else if theno > 1 {
				_topid := lib.HF_Atoi(top[i-1])
				GetServer().AddBindPlayers(_topid, 0, addgold)

				c := GetServer().Redis.Get()
				goldkey := fmt.Sprintf("%s_%d_%d_%d_%d", "score2", _topid, time.Now().Year(), time.Now().Month(), time.Now().Day())
				thegold, err := redis.Float64(c.Do("GET", goldkey))
				if err != nil {
					thegold = 0
				}
				thegold += addgold
				c.Do("SET", goldkey, thegold)
				c.Do("EXPIRE", goldkey, 86400*3)
				c.Close()
			}
		}
		if data.Parent != 0 && score_total != 0 {
			GetServer().QueueParentCostLog(fmt.Sprintf("INSERT INTO `parent_cost_log` (`parent`, `score`) VALUES(%d, %f)", data.Parent, score_total))
		}
	}

	lib.GetLogMgr().Output(lib.LOG_ERROR, "结束统计")
}
