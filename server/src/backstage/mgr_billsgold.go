package backstage

import (
	"fmt"
	//"github.com/garyburd/redigo/redis"
	"lib"
	"log"
	"runtime/debug"
	"strings"
	"time"
)

type BillsGoldMgr struct {
}

var billsgoldSingleton *BillsGoldMgr = nil

//! 得到服务器指针
func GetBillsGoldMgr() *BillsGoldMgr {
	if billsgoldSingleton == nil {
		billsgoldSingleton = new(BillsGoldMgr)
	}

	return billsgoldSingleton
}

type SQL_Bill_Num struct {
	Topup_num int64
}

//! 得到当前读取到第几条
func (self *BillsGoldMgr) GetCount() int64 {
	var sql SQL_Bill_Num
	if !GetServer().HT_DB.GetOneData("SELECT `topup_num` FROM `fa_step`", &sql) {
		log.Fatal("get score_num fail")
		return 0
	}
	return sql.Topup_num
}

func (self *BillsGoldMgr) SetCount(id int64) {
	a, b := GetServer().HT_DB.Exec(fmt.Sprintf("UPDATE `fa_step` SET topup_num = %d where id = 1", id))
	if a == 0 && b == 0 {
		log.Fatal("SetCount fail")
	}
}

type SQL_BillGold struct {
	Id       int64
	Uid      int
	Num      int
	GameType int
	Time     int64
}

//! 开始计算代理返利
func (self *BillsGoldMgr) Do() {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	//if GetServer().Con.AgentMode == 0 { //! 抽水模式不计算
	//	return
	//}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	lib.GetLogMgr().Output(lib.LOG_ERROR, "开始统计")
	maxid := self.GetCount()

	billgold := make(map[int]int) //! 总流水的
	wingold := make(map[int]int)  //! 赢的流水
	lostgold := make(map[int]int) //! 输的流水
	var sql SQL_BillGold
	res := GetServer().Game_DB.GetAllData(fmt.Sprintf("select * from `log_bills` WHERE id > %d ORDER BY id ASC", maxid), &sql)
	if len(res) == 0 {
		lib.GetLogMgr().Output(lib.LOG_ERROR, "结束统计")
		return
	}
	log.Println("读取了:", len(res))
	maxid = res[len(res)-1].(*SQL_BillGold).Id
	self.SetCount(maxid)
	for i := 0; i < len(res); i++ {
		if res[i].(*SQL_BillGold).Num >= 0 {
			billgold[res[i].(*SQL_BillGold).Uid] += res[i].(*SQL_BillGold).Num
			wingold[res[i].(*SQL_BillGold).Uid] += res[i].(*SQL_BillGold).Num
		} else {
			billgold[res[i].(*SQL_BillGold).Uid] -= res[i].(*SQL_BillGold).Num
			lostgold[res[i].(*SQL_BillGold).Uid] -= res[i].(*SQL_BillGold).Num
		}
	}
	log.Println("合并后有:", len(billgold))

	for key, value := range wingold {
		data := GetServer().GetAgentUser(key)
		if data.Agid == 0 {
			continue
		}
		if data.Parent != 0 { //! 一个特殊账号的下级
			parent := GetServer().GetAgentUser(data.Parent)
			if parent.Parent != data.Parent { //! 资格被取消
				data.Parent = 0
				GetServer().SetAgentUser(data)
				GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET `parent` = %d WHERE agid = %d", data.Parent, data.Agid), []byte(""), false)
			} else {
				GetServer().QueueParentLog(fmt.Sprintf("INSERT INTO `parent_win_bills` (`uid`, `gold`, `parent`, `time`) VALUES(%d, %d, %d, %d)", key, value, data.Parent, time.Now().Unix()))
			}
		}
	}
	for key, value := range lostgold {
		data := GetServer().GetAgentUser(key)
		if data.Agid == 0 {
			continue
		}
		if data.Parent != 0 { //! 一个特殊账号的下级
			parent := GetServer().GetAgentUser(data.Parent)
			if parent.Parent != data.Parent { //! 资格被取消
				data.Parent = 0
				GetServer().SetAgentUser(data)
				GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET `parent` = %d WHERE agid = %d", data.Parent, data.Agid), []byte(""), false)
			} else {
				GetServer().QueueParentLog(fmt.Sprintf("INSERT INTO `parent_lost_bills` (`uid`, `gold`, `parent`, `time`) VALUES(%d, %d, %d, %d)", key, value, data.Parent, time.Now().Unix()))
			}
		}
	}

	for key, value := range billgold {
		data := GetServer().GetAgentUser(key)
		if data.Agid == 0 {
			continue
		}

		//! 先给自己加流水
		data.AddBills(int64(value))
		GetServer().SetAgentUser(data)
		GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET `allbills` = %d, `daybills` = %d, `weekbills` = %d, `monthbills` = %d, `timebills` = %d WHERE agid = %d", data.AllBills, data.DayBills, data.WeekBills, data.MonthBills, data.TimeBills, key), []byte(""), false)

		//! 给上级加流水
		top := make([]string, 0)
		if data.Top_Group != "" {
			top = strings.Split(data.Top_Group, ",")
		}

		if len(top) == 0 {
			continue
		}

		for i := 0; i < len(top); i++ {
			topid := lib.HF_Atoi(top[i])
			if topid == 0 {
				continue
			}
			agenttop := GetServer().GetAgentUser(topid)
			if agenttop.Agid == 0 {
				continue
			}
			if i == 0 { //! 直属上级
				agenttop.AddCommission(int64(value), true)
			} else { //! 下级
				agenttop.AddCommission(int64(value), false)
			}
			GetServer().SetAgentUser(agenttop)
			GetServer().QueueAgentUser(fmt.Sprintf("UPDATE `fa_agent_user` SET `commission` = %d, `bills1` = %d, `bills2` = %d, `timecommission` = %d WHERE agid = %d", agenttop.Commission, agenttop.Bills1, agenttop.Bills2, agenttop.TimeCommission, topid), []byte(""), false)
		}
	}

	lib.GetLogMgr().Output(lib.LOG_ERROR, "结束统计")
}
