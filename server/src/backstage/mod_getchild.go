package backstage

import (
	"lib"
	"net/http"
	"sort"
	"strings"
)

type Son_Child struct {
	Uid          int     `json:"uid"`
	NickName     string  `json:"nickname"`
	Head         string  `json:"head"`
	Score1       float64 `json:"score1"`
	Score2       float64 `json:"score2"`
	Level        int     `json:"level"`
	Deepin       int     `json:"deepin"`
	Bind_time    string  `json:"bind_time"`
	Score        float64 `json:"score"`
	T_Score      int     `json:"t_score"`
	AgentCard    int     `json:"agentcard"`
	AllCost      int64   `json:"allcost"`      //! 总抽水
	AllBills     int64   `json:"allbills"`     //! 总流水
	DayBills     int64   `json:"daybills"`     //! 日流水
	WeekBills    int64   `json:"weekbills"`    //! 周流水
	MonthBills   int64   `json:"monthbills"`   //! 月流水
	TimeBills    int64   `json:"timebills"`    //! 统计流水时间
	Commission   int     `json:"commission"`   //! 总佣金
	T_Commission int     `json:"t_commission"` //! 已提佣金
	AreaScale    int     `json:"areascale"`    //! 区域比例
	AreaScore    float64 `json:"areascore"`    //! 区域收益
	AreaTScore   int     `json:"areatscore"`   //! 已提现收益
}

type LstSon_Child []Son_Child

func (a LstSon_Child) Len() int { // 重写 Len() 方法
	return len(a)
}

func (a LstSon_Child) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}

func (a LstSon_Child) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[i].Score1+a[i].Score2 > a[j].Score1+a[j].Score2
}

type HT_Child struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Score float64     `json:"score"`
	Num   int         `json:"num"`
	Add   int         `json:"add"`
	Deep  int         `json:"deep"`
	Info  []Son_Child `json:"info"`
}

//! 得到推广成员
func GetChild(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))
	page := lib.HF_Atoi(req.FormValue("page")) //! 第几页
	num := lib.HF_Atoi(req.FormValue("num"))   //! 每页几个
	if uid == 0 || page == 0 || num == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	all := GetServer().GetChildren(uid)
	sonlst := make(LstSon_Child, 0)
	for i := 0; i < len(all.Child); i++ {
		var son Son_Child
		agent := GetServer().GetAgentUser(all.Child[i])
		if agent.Agid == 0 {
			continue
		}
		bind := GetServer().GetBindPlayers(all.Child[i])
		if bind.Uid == 0 {
			continue
		}
		son.Uid = all.Child[i]
		son.NickName = agent.NickName
		son.Head = agent.Head
		son.Score1 = bind.Score1
		son.Score2 = bind.Score2
		son.Level = agent.Level
		son.Deepin = agent.Deepin
		son.Bind_time = bind.Bind_time
		son.Score = agent.Score
		son.T_Score = agent.T_Score
		son.AgentCard = agent.Card
		son.AllBills = agent.AllBills
		son.DayBills = agent.DayBills
		son.WeekBills = agent.WeekBills
		son.MonthBills = agent.MonthBills
		son.TimeBills = agent.TimeBills
		son.Commission = agent.Commission
		son.T_Commission = agent.T_Commission
		son.AreaScale = agent.AreaScale
		son.AreaScore = agent.AreaScore
		son.AreaTScore = agent.AreaTScore
		sonlst = append(sonlst, son)
	}
	sort.Sort(LstSon_Child(sonlst))

	var node HT_Child
	node.Score = GetServer().GetMyMoney(uid)
	node.Num = len(all.Child)
	node.Add = GetServer().GetTodayChild(uid)
	node.Deep = 1
	node.Info = make([]Son_Child, 0)
	for i := (page - 1) * num; i < lib.HF_MinInt((page-1)*num+num, len(sonlst)); i++ {
		node.Info = append(node.Info, sonlst[i])
	}
	w.Write(lib.HF_JtoB(&node))
}

//! 得到所有下级
func GetAllChild(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))
	deep := lib.HF_Atoi(req.FormValue("deep"))
	if uid == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	all := make([]int, 0)
	HF_FindSonFromDeep(uid, &all, 0, deep)
	sonlst := make(LstSon_Child, 0)
	for i := 0; i < len(all); i++ {
		var son Son_Child
		agent := GetServer().GetAgentUser(all[i])
		if agent.Agid == 0 {
			continue
		}
		bind := GetServer().GetBindPlayers(all[i])
		if bind.Uid == 0 {
			continue
		}
		son.Uid = all[i]
		son.NickName = agent.NickName
		son.Head = agent.Head
		son.Score1 = bind.Score1
		son.Score2 = bind.Score2
		son.Level = agent.Level
		son.Deepin = agent.Deepin
		son.Bind_time = bind.Bind_time
		son.Score = agent.Score
		son.T_Score = agent.T_Score
		son.AgentCard = agent.Card
		son.AllCost = agent.AllCost
		son.AllBills = agent.AllBills
		son.DayBills = agent.DayBills
		son.WeekBills = agent.WeekBills
		son.MonthBills = agent.MonthBills
		son.TimeBills = agent.TimeBills
		son.Commission = agent.Commission
		son.T_Commission = agent.T_Commission
		son.AreaScale = agent.AreaScale
		son.AreaScore = agent.AreaScore
		son.AreaTScore = agent.AreaTScore
		sonlst = append(sonlst, son)
	}
	sort.Sort(LstSon_Child(sonlst))

	var node HT_Child
	node.Score = GetServer().GetMyMoney(uid)
	node.Num = len(all)
	node.Add = GetServer().GetTodayChild(uid)
	node.Deep = deep
	node.Info = make([]Son_Child, 0)
	for i := 0; i < len(sonlst); i++ {
		node.Info = append(node.Info, sonlst[i])
	}
	w.Write(lib.HF_JtoB(&node))
}

//! 搜索单个child
func GetOneChild(w http.ResponseWriter, req *http.Request) {
	//if !rjmgr.GetRJMgr().IsLicensing() { //! 未授权
	//	return
	//}

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

	destagent := GetServer().GetAgentUser(destuid)
	if destagent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	find := false
	deep := 0
	if destagent.Top_Group != "" {
		top := strings.Split(destagent.Top_Group, ",")
		for i := 0; i < len(top); i++ {
			deep++
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

	sonlst := make(LstSon_Child, 0)

	var son Son_Child
	agent = destagent
	bind := GetServer().GetBindPlayers(destuid)
	if bind.Uid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}
	son.Uid = destuid
	son.NickName = agent.NickName
	son.Head = agent.Head
	son.Score1 = bind.Score1
	son.Score2 = bind.Score2
	son.Level = agent.Level
	son.Deepin = agent.Deepin
	son.Bind_time = bind.Bind_time
	son.Score = agent.Score
	son.T_Score = agent.T_Score
	son.AgentCard = agent.Card
	son.AllCost = agent.AllCost
	son.AllBills = agent.AllBills
	son.DayBills = agent.DayBills
	son.WeekBills = agent.WeekBills
	son.MonthBills = agent.MonthBills
	son.TimeBills = agent.TimeBills
	son.Commission = agent.Commission
	son.T_Commission = agent.T_Commission
	son.AreaScale = agent.AreaScale
	son.AreaScore = agent.AreaScore
	son.AreaTScore = agent.AreaTScore
	sonlst = append(sonlst, son)

	sort.Sort(LstSon_Child(sonlst))

	var node HT_Child
	node.Score = GetServer().GetMyMoney(uid)
	node.Num = 1
	node.Add = GetServer().GetTodayChild(uid)
	node.Deep = deep
	node.Info = make([]Son_Child, 0)
	for i := 0; i < len(sonlst); i++ {
		node.Info = append(node.Info, sonlst[i])
	}
	w.Write(lib.HF_JtoB(&node))
}
