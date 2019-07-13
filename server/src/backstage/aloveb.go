////! a扫码b

package backstage

import (
//	"fmt"
//	"io/ioutil"
//	"log"
//	"net/http"
//	//"net/url"
//	"staticfunc"
//	"strings"
//	"sync"
//	"time"
)

//type ALoveB_Agent struct {
//	Agid      int
//	Top_Group string
//	Deepin    int
//}

//type ALoveB_Bind struct {
//	Uid  int
//	Agid int
//}

//var AccountLock *sync.RWMutex = new(sync.RWMutex)

//func ALoveB(w http.ResponseWriter, req *http.Request) {
//	AccountLock.Lock()
//	defer AccountLock.Unlock()

//	aunionid := req.FormValue("aunionid")
//	aopenid := req.FormValue("aopenid")
//	buid := lib.HF_Atoi(req.FormValue("buid"))
//	log.Println(aunionid)
//	log.Println(aopenid)
//	log.Println(buid)
//	if buid == 0 {
//		w.Write([]byte("-98")) //! 参数错误
//		log.Println("-98")
//		return
//	}

//	//! 向游戏服务器获取uid
//	str := fmt.Sprintf("http://jyqp.hbyouyou.com:8031/wxtouid?unionid=%s&openid=%s", aunionid, aopenid)
//	res, err := http.Get(str)
//	if err != nil {
//		w.Write([]byte("-1")) //! 服务器无响应
//		return
//	}

//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		w.Write([]byte("-2")) //! 游戏服务器返回错误
//		log.Println("-2")
//		return
//	}

//	auid := lib.HF_Atoi(string(body))

//	if auid == buid {
//		w.Write([]byte("-6")) //! 数据库错误
//		log.Println("-6")
//		return
//	}

//	var bagent ALoveB_Agent
//	if !HT_DB.GetOneData(fmt.Sprintf("select `agid`, `top_group`, `deepin` from `fa_agent_user` where `agid` = %d", buid), &bagent) {
//		w.Write([]byte("-99")) //! 数据库错误
//		log.Println("-99")
//		return
//	}

//	if bagent.Agid == 0 {
//		w.Write([]byte("-3")) //! 绑定失败,被扫者不存在
//		log.Println("-3")
//		return
//	}

//	btop := make([]string, 0)
//	if bagent.Top_Group != "" {
//		btop = strings.Split(bagent.Top_Group, ",")
//		for i := 0; i < len(btop); i++ {
//			if lib.HF_Atoi(btop[i]) == auid {
//				w.Write([]byte("-5")) //! 扫码者在被扫码者的上级
//				log.Println("-5")
//				return
//			}
//		}
//	}

//	var abind ALoveB_Bind
//	if !HT_DB.GetOneData(fmt.Sprintf("select `uid`, `agid` from `fa_bind_players` where `uid` = %d", auid), &abind) {
//		w.Write([]byte("-99")) //! 数据库错误
//		log.Println("-99")
//		return
//	}

//	if abind.Uid != 0 {
//		w.Write([]byte("-4")) //! 该用户已经绑定了
//		log.Println("-4")
//		return
//	}

//	//! 更新用户表
//	var aagent ALoveB_Agent
//	if !HT_DB.GetOneData(fmt.Sprintf("select `agid`, `top_group`, `deepin` from `fa_agent_user` where `agid` = %d", auid), &aagent) {
//		w.Write([]byte("-99")) //! 数据库错误
//		log.Println("-99")
//		return
//	}

//	if aagent.Deepin >= 7 { //! 这个深度的人不能绑定上级
//		w.Write([]byte("-7"))
//		log.Println("-7")
//		return
//	}

//	//! 插入绑定表
//	HT_DB.Exec(fmt.Sprintf("INSERT into fa_bind_players(`uid`, `agid`, `bind_time`, `score1`, `score2`) values(%d, %d, '%s', %d, %d)", auid, buid, time.Now().Format(staticfunc.TIMEFORMAT), 0, 0))

//	//! 修改数据表
//	if len(btop) == 0 {
//		aagent.Top_Group = fmt.Sprintf("%d", buid)
//	} else {
//		aagent.Top_Group = fmt.Sprintf("%d", buid)
//		for i := 0; i < staticfunc.HF_MinInt(len(btop), 9); i++ {
//			aagent.Top_Group += ","
//			aagent.Top_Group += btop[i]
//		}
//	}
//	if aagent.Agid == 0 { //! 用户表里没有
//		aagent.Agid = auid
//		HT_DB.Exec(fmt.Sprintf("insert into fa_agent_user(`agid`, `open_id`, `union_id`, `top_group`, `password`) values(%d, '%s', '%s', '%s', '')", aagent.Agid, aopenid, aunionid, aagent.Top_Group))
//	} else {
//		HT_DB.Exec(fmt.Sprintf("update fa_agent_user set `top_group` = '%s' where `agid` = %d", aagent.Top_Group, aagent.Agid))
//	}

//	atop := strings.Split(aagent.Top_Group, ",")

//	//! 修改A的所有下级
//	lst := make([]int, 0)
//	FindSon(aagent.Agid, &lst, 0)
//	for i := 0; i < len(lst); i++ {
//		var nagent ALoveB_Agent
//		if !HT_DB.GetOneData(fmt.Sprintf("select `agid`, `top_group`, `deepin` from `fa_agent_user` where `agid` = %d", lst[i]), &nagent) {
//			continue
//		}

//		if nagent.Top_Group == "" {
//			continue
//		}

//		top := strings.Split(nagent.Top_Group, ",")
//		for j := 0; j < staticfunc.HF_MinInt(len(atop), 10-len(top)); j++ {
//			nagent.Top_Group += ","
//			nagent.Top_Group += atop[j]
//		}
//		HT_DB.Exec(fmt.Sprintf("update fa_agent_user set `top_group` = '%s' where `agid` = %d", nagent.Top_Group, nagent.Agid))
//	}

//	w.Write([]byte("0"))
//	log.Println("0")
//}
