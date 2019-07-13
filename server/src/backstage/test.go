package backstage

import (
//	//"io/ioutil"
//	//"encoding/base64"
//	//"encoding/json"
//	"github.com/garyburd/redigo/redis"
//	//"io/ioutil"
//	"log"
//	"net/http"
//	_ "net/http/pprof"
//	//"net/url"
//	"fmt"
//	"io/ioutil"
//	"os"
//	"os/signal"
//	"runtime"
//	"staticfunc"
//	"strings"
//	"syscall"
//	"time"
)

////! 注册信号量
//func handleSignal(signalType os.Signal, handleFun func(*chan os.Signal)) {
//	ch := make(chan os.Signal)
//	signal.Notify(ch, signalType)
//	go handleFun(&ch)
//}

////! 管道破裂
//func handlePIPE(ch *chan os.Signal) {
//	for {
//		<-*ch
//		log.Println("get a SIGPIPE")
//	}
//}

////! ctrl+z
//func handleTSTP(ch *chan os.Signal) {
//	for {
//		<-*ch
//		log.Println("get a SIGTSTP")
//	}
//}

////! gdb trap
//func handleTRAP(ch *chan os.Signal) {
//	for {
//		<-*ch
//		log.Println("get a SIGTRAP")
//	}
//}

////! ctrl+c
//func handleINT(ch *chan os.Signal) {
//	for {
//		<-*ch
//		log.Println("get a SIGINT")
//		GetServer().Close()
//	}
//}

//func Run() {
//	ticker := time.NewTicker(time.Second * 300)
//	for {
//		<-ticker.C
//		Do()
//	}

//	//！ 关掉定时器
//	ticker.Stop()
//}

//type SQL_TopGroup struct {
//	Agid      int
//	Top_group string
//	Level     int
//}

//type SQL_AgentTop struct {
//	Deepin int
//	Rating string
//	Level  int
//}

//type SQL_AgentNum struct {
//	Num int
//}

//type SQL_BindUid struct {
//	Uid int
//}

////! 得到当前读取到第几条
//func GetCount() int64 {
//	var sql SQL_Score_Num
//	if !HT_DB.GetOneData("SELECT `score_num` FROM `fa_step`", &sql) {
//		log.Fatal("get score_num fail")
//		return 0
//	}
//	return sql.Score_num
//}

//func SetCount(id int64) {
//	a, b := HT_DB.Exec(fmt.Sprintf("UPDATE `fa_step` SET score_num = %d where id = 1", id))
//	if a == 0 && b == 0 {
//		log.Fatal("SetCount fail")
//	}
//}

//func Do() {
//	log.Println("开始统计")
//	maxid := GetCount()

//	agentgold := make(map[int]int)
//	var sql SQL_AgentGold
//	res := Game_DB.GetAllData(fmt.Sprintf("SELECT id,uid,gold FROM log_agentgold WHERE id > %d ORDER BY id ASC", maxid), &sql)
//	if len(res) == 0 {
//		return
//	}
//	log.Println("读取了:", len(res))
//	maxid = res[len(res)-1].(*SQL_AgentGold).Id
//	SetCount(maxid)
//	for i := 0; i < len(res); i++ {
//		agentgold[res[i].(*SQL_AgentGold).Uid] += res[i].(*SQL_AgentGold).Gold
//	}
//	log.Println("合并后有:", len(agentgold))

//	for key, value := range agentgold {
//		var data SQL_TopGroup
//		if !HT_DB.GetOneData(fmt.Sprintf("SELECT `agid`, `top_group`, `level` FROM fa_agent_user WHERE agid = %d", key), &data) {
//			continue
//		}

//		if data.Agid == 0 {
//			continue
//		}

//		top := make([]string, 0)
//		if data.Top_group != "" {
//			top = strings.Split(data.Top_group, ",")
//		}

//		if data.Level > 0 {
//			top = append([]string{fmt.Sprintf("%d", key)}, top...)
//		}

//		if len(top) == 0 {
//			continue
//		}

//		//deep := 0
//		theno := 0
//		if len(top) <= 4 && len(top) > 1 {
//			top = append(top, "338888")
//			top = append(top, "284177")
//		}
//		for i := 0; i < len(top); i++ {
//			topid := lib.HF_Atoi(top[i])
//			if topid != key {
//				theno++
//			}
//			if topid == 0 {
//				continue
//			}
//			var agenttop SQL_AgentTop
//			if !HT_DB.GetOneData(fmt.Sprintf("SELECT `deepin`, `rating`, `level` FROM fa_agent_user WHERE agid = %d", topid), &agenttop) {
//				continue
//			}
//			if agenttop.Deepin <= i && topid != 338888 && topid != 284177 {
//				continue
//			}
//			if agenttop.Level <= 0 {
//				if i == 0 || i == 1 { //! 模拟返利深度
//					willgold := float64(value*35) / 10000.0
//					if i == 1 {
//						willgold = float64(value*10) / 10000.0
//					}
//					c := Redis.Get()
//					goldkey := fmt.Sprintf("%s_%d_%d_%d_%d", "willgold", topid, time.Now().Year(), time.Now().Month(), time.Now().Day())
//					addgold, err := redis.Float64(c.Do("GET", goldkey))
//					if err != nil {
//						addgold = 0
//					}
//					addgold += willgold
//					c.Do("SET", goldkey, addgold)
//					c.Do("EXPIRE", goldkey, 86400*3)
//					c.Close()

//					if addgold >= 30 { //! 活跃度达标
//						var agentnum SQL_AgentNum
//						if HT_DB.GetOneData(fmt.Sprintf("select count(*) as `num` from fa_bind_players where agid = %d", topid), &agentnum) {
//							if agentnum.Num >= 10 {
//								agenttop.Level = 1
//								HT_DB.Exec(fmt.Sprintf("UPDATE `fa_agent_user` SET score = score + %f, `level` = %d WHERE agid = %d", addgold, agenttop.Level, topid))
//								HT_DB.Exec(fmt.Sprintf("INSERT INTO `score_log` (operator_id,agid,change_score,save_time,`action`) VALUES(%d, %d, %f, '%s', 1)", topid, topid, addgold, time.Now().Format(staticfunc.TIMEFORMAT)))
//							}
//						}
//					}
//				}

//				continue
//			}
//			fanli := 0
//			rating := strings.Split(agenttop.Rating, ",")
//			if i >= len(rating) {
//				continue //! 深度太大就放弃
//			} else {
//				fanli = lib.HF_Atoi(rating[i])
//				//if i >= 3 {
//				//	fanli--
//				//	if fanli <= 0 { //! 最少返利1%
//				//		fanli = 1
//				//	}
//				//}
//			}
//			if topid == 338888 || topid == 284177 {
//				fanli = 1
//			}
//			addgold := float64(value*fanli) / 10000.0
//			//deep++

//			HT_DB.Exec(fmt.Sprintf("UPDATE `fa_agent_user` SET score = score + %f WHERE agid = %d", addgold, topid))
//			HT_DB.Exec(fmt.Sprintf("INSERT INTO `score_log` (operator_id,agid,change_score,save_time,`action`) VALUES(%d, %d, %f, '%s', 1)", key, topid, addgold, time.Now().Format(staticfunc.TIMEFORMAT)))
//			if topid == 338888 || topid == 284177 {
//				continue
//			}
//			if theno == 1 { //! 直属上级
//				HT_DB.Exec(fmt.Sprintf("UPDATE `fa_bind_players` SET score1 = score1 + %f WHERE uid = %d", addgold, key))

//				c := Redis.Get()
//				goldkey := fmt.Sprintf("%s_%d_%d_%d_%d", "score1", key, time.Now().Year(), time.Now().Month(), time.Now().Day())
//				thegold, err := redis.Float64(c.Do("GET", goldkey))
//				if err != nil {
//					thegold = 0
//				}
//				thegold += addgold
//				c.Do("SET", goldkey, thegold)
//				c.Do("EXPIRE", goldkey, 86400*3)
//				c.Close()
//			} else if theno > 1 {
//				_topid := lib.HF_Atoi(top[i-1])
//				HT_DB.Exec(fmt.Sprintf("UPDATE `fa_bind_players` SET score2 = score2 + %f WHERE uid = %d", addgold, _topid))

//				c := Redis.Get()
//				goldkey := fmt.Sprintf("%s_%d_%d_%d_%d", "score2", _topid, time.Now().Year(), time.Now().Month(), time.Now().Day())
//				thegold, err := redis.Float64(c.Do("GET", goldkey))
//				if err != nil {
//					thegold = 0
//				}
//				thegold += addgold
//				c.Do("SET", goldkey, thegold)
//				c.Do("EXPIRE", goldkey, 86400*3)
//				c.Close()
//			}
//		}
//	}

//	log.Println("统计完毕")
//}

////! 更换上级
//func ChageParent(w http.ResponseWriter, req *http.Request) { //sonid int, parentid int) []byte {
//	sonid := lib.HF_Atoi(req.FormValue("sonid"))
//	parentid := lib.HF_Atoi(req.FormValue("parentid"))

//	if sonid == 0 {
//		w.Write([]byte("参数错误"))
//		return
//	}

//	if sonid == parentid {
//		w.Write([]byte("不能移动到自己下级"))
//		return
//	}

//	var sondata SQL_TopGroup
//	if !HT_DB.GetOneData(fmt.Sprintf("SELECT `agid`, `top_group`, `level` FROM fa_agent_user WHERE agid = %d", sonid), &sondata) {
//		if sondata.Agid == 0 {
//			w.Write([]byte("sonid错误"))
//			return
//		}
//	}

//	var parentdata SQL_TopGroup
//	if parentid != 0 { //! 移动上级
//		if !HT_DB.GetOneData(fmt.Sprintf("SELECT `agid`, `top_group`, `level` FROM fa_agent_user WHERE agid = %d", parentid), &parentdata) {
//			w.Write([]byte("上级id错误"))
//			return
//		}

//		if parentdata.Agid == 0 {
//			w.Write([]byte("上级id错误"))
//			return
//		}

//		top := make([]string, 0)
//		if parentdata.Top_group != "" {
//			top = strings.Split(parentdata.Top_group, ",")
//		}
//		for i := 0; i < len(top); i++ {
//			if lib.HF_Atoi(top[i]) == sonid {
//				w.Write([]byte("移动失败，上级玩家在你的下级"))
//				return
//			}
//		}
//	}

//	if parentid == 0 { //! 删除上级
//		oldparentid := 0
//		if sondata.Top_group == "" {
//			w.Write([]byte("移动成功"))
//			return
//		} else {
//			top := strings.Split(sondata.Top_group, ",")
//			oldparentid = lib.HF_Atoi(top[0])
//		}
//		HT_DB.Exec(fmt.Sprintf("DELETE from fa_bind_players where `uid` = %d", sonid))

//		//! 更新自己
//		HT_DB.Exec(fmt.Sprintf("UPDATE `fa_agent_user` SET `top_group` = '' WHERE agid = %d", sonid))

//		//! 更新下级
//		lst := make([]int, 0)
//		FindSon(sonid, &lst, 0)
//		for i := 0; i < len(lst); i++ {
//			var son SQL_TopGroup
//			if !HT_DB.GetOneData(fmt.Sprintf("SELECT `agid`, `top_group`, `level` FROM fa_agent_user WHERE agid = %d", lst[i]), &son) {
//				continue
//			}
//			if son.Agid == 0 {
//				continue
//			}
//			top := make([]string, 0)
//			if son.Top_group != "" {
//				top = strings.Split(son.Top_group, ",")
//			}
//			for i := 0; i < len(top); i++ {
//				if lib.HF_Atoi(top[i]) == oldparentid {
//					top = top[0:i]
//					break
//				}
//			}
//			son.Top_group = ""
//			for i := 0; i < len(top); i++ {
//				if i == 0 {
//					son.Top_group += top[i]
//				} else {
//					son.Top_group += ","
//					son.Top_group += top[i]
//				}
//			}
//			HT_DB.Exec(fmt.Sprintf("UPDATE `fa_agent_user` SET `top_group` = '%s' WHERE agid = %d", son.Top_group, son.Agid))
//		}
//	} else { //! 移动上级
//		oldparentid := 0
//		if sondata.Top_group == "" { //! 本来没有上级
//			HT_DB.Exec(fmt.Sprintf("INSERT into fa_bind_players(`uid`, `agid`, `bind_time`, `score1`, `score2`) values(%d, %d, '%s', %d, %d)", sonid, parentid, time.Now().Format(staticfunc.TIMEFORMAT), 0, 0))
//		} else {
//			top := strings.Split(sondata.Top_group, ",")
//			if len(top) > 0 && lib.HF_Atoi(top[0]) == parentid { //! 上级没有发生改变
//				w.Write([]byte("移动成功"))
//				return
//			}
//			oldparentid = lib.HF_Atoi(top[0])
//			HT_DB.Exec(fmt.Sprintf("UPDATE `fa_bind_players` SET `agid` = %d WHERE uid = %d", parentid, sonid))
//		}

//		//! 更新自己
//		//if parentdata.Top_group
//		//HT_DB.Exec(fmt.Sprintf("UPDATE `fa_agent_user` SET `top_group` = '%d,%s' WHERE agid = %d", parentid, parentdata.Top_group, sonid))

//		//! 更新下线
//		lst := []int{sonid}
//		FindSon(sonid, &lst, 0)
//		for i := 0; i < len(lst); i++ {
//			var son SQL_TopGroup
//			if !HT_DB.GetOneData(fmt.Sprintf("SELECT `agid`, `top_group`, `level` FROM fa_agent_user WHERE agid = %d", lst[i]), &son) {
//				continue
//			}
//			if son.Agid == 0 {
//				continue
//			}
//			top := make([]string, 0)
//			if son.Top_group != "" {
//				top = strings.Split(son.Top_group, ",")
//			}
//			for i := 0; i < len(top); i++ {
//				if lib.HF_Atoi(top[i]) == oldparentid {
//					top = top[0:i]
//					break
//				}
//			}
//			top = append(top, fmt.Sprintf("%d", parentid))
//			if parentdata.Top_group != "" {
//				newtop := strings.Split(parentdata.Top_group, ",")
//				top = append(top, newtop...)
//			}
//			if len(top) > 10 {
//				top = top[0:10]
//			}
//			son.Top_group = ""
//			for i := 0; i < len(top); i++ {
//				if i == 0 {
//					son.Top_group += top[i]
//				} else {
//					son.Top_group += ","
//					son.Top_group += top[i]
//				}
//			}
//			HT_DB.Exec(fmt.Sprintf("UPDATE `fa_agent_user` SET `top_group` = '%s' WHERE agid = %d", son.Top_group, son.Agid))
//		}
//	}

//	w.Write([]byte("移动成功"))
//	return
//}

////! 查找一个玩家所有的子代理
//func FindSon(uid int, lst *[]int, deep int) {
//	deep++
//	if deep > 10 { //! 最多只能10层
//		return
//	}
//	var bind SQL_BindUid
//	res := HT_DB.GetAllData(fmt.Sprintf("select `uid` from fa_bind_players where `agid` = %d", uid), &bind)
//	for i := 0; i < len(res); i++ {
//		_uid := res[i].(*SQL_BindUid).Uid
//		*lst = append(*lst, _uid)
//		FindSon(_uid, lst, deep)
//	}
//}

//var Game_DB *DBServer = new(DBServer)
//var HT_DB *DBServer = new(DBServer)
//var Host_DB *DBServer = new(DBServer)
//var Redis *redis.Pool = staticfunc.NewPool("127.0.0.1:6379", 2)

//func main() {
//	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

//	//! 注册信号量
//	handleSignal(syscall.SIGPIPE, handlePIPE)
//	handleSignal(syscall.SIGTSTP, handleTSTP)
//	handleSignal(syscall.SIGTRAP, handleTRAP)
//	handleSignal(syscall.SIGINT, handleINT)

//	//! 初始化数据库
//	if !Game_DB.Init("root:clannad!@#$@tcp(jyqp.hbyouyou.com:3306)/qp?charset=utf8&timeout=10s") {
//		log.Fatal("db1 fail")
//		return
//	}

//	if !HT_DB.Init("jyqp:Clannad!@#$@tcp(rm-bp17873057h4utelbo.mysql.rds.aliyuncs.com:3306)/jyqp_ht?charset=utf8&timeout=10s") {
//		log.Fatal("db2 fail")
//		return
//	}

//	if !Host_DB.Init("jyqp:Clannad!@#$@tcp(rm-bp17873057h4utelbo.mysql.rds.aliyuncs.com:3306)/jyqp_host?charset=utf8&timeout=10s") {
//		log.Fatal("db3 fail")
//		return
//	}

//	go Do()
//	go Run()

//	//if !HT_DB.Init("root:yang1231231@tcp(127.0.0.1:3306)/jyqp_ht?charset=utf8&timeout=10s") {
//	//	log.Fatal("db2 fail")
//	//	return
//	//}

//	go GetServer().RunQueueAgentUser()
//	go GetServer().RunBindPlayers()
//	go GetServer().RunScoreLog()

//	http.HandleFunc("/chageparent", ChageParent)
//	http.HandleFunc("/aloveb", ALoveB)

//	//! 后台接口
//	http.HandleFunc("/gettodaywillgold", GetTodayWillGold)
//	http.HandleFunc("/getyestodaywillgold", GetYesTodayWillGold)
//	http.HandleFunc("/gettodayaddgold", GetTodayAddGold)
//	http.HandleFunc("/countdec", CountDec)
//	http.HandleFunc("/jyaswapb", ASwapB)             //! A扫码B
//	http.HandleFunc("/jychangeparent", ChangeParent) //! 更换上级
//	http.HandleFunc("/jygetinfo", GetInfo)           //! 得到信息

//	//! 绑定一个http服务
//	log.Fatal(http.ListenAndServe(":1231", nil))
//}

//type SQL_Money struct {
//	Money int
//}

//type SQL_MoneyF struct {
//	Money float32
//}

////! 统计
//func CountDec(w http.ResponseWriter, req *http.Request) {
//	//! 向游戏服务器获取在线人数
//	str := fmt.Sprintf("http://jyqp.hbyouyou.com:8031/getonline")
//	res, err := http.Get(str)
//	if err != nil {
//		w.Write([]byte("-1")) //! 服务器无响应
//		return
//	}

//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		w.Write([]byte("-2")) //! 游戏服务器返回错误
//		return
//	}

//	//! 在线人数
//	online := string(body)

//	now := time.Now()
//	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
//	yestoday := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local)

//	//! 获得今日金币充值
//	var money SQL_Money
//	Host_DB.GetOneData(fmt.Sprintf("select sum(money) as money from `order` where `status` = 1 and ptype = 1 and Date(`time`) >= '%s'", today.Format(staticfunc.TIMEFORMAT)), &money)
//	todaygold := money.Money

//	//! 获得昨日金币充值
//	Host_DB.GetOneData(fmt.Sprintf("select sum(money) as money from `order` where `status` = 1 and ptype = 1 and Date(`time`) >= '%s' and Date(`time`) < '%s'", yestoday.Format(staticfunc.TIMEFORMAT), today.Format(staticfunc.TIMEFORMAT)), &money)
//	yestodaygold := money.Money

//	//! 获得今日房卡充值
//	Host_DB.GetOneData(fmt.Sprintf("select sum(money) as money from `order` where `status` = 1 and ptype = 2 and Date(`time`) >= '%s'", today.Format(staticfunc.TIMEFORMAT)), &money)
//	todaycard := money.Money

//	//! 获得昨日房卡充值
//	Host_DB.GetOneData(fmt.Sprintf("select sum(money) as money from `order` where `status` = 1 and ptype = 2 and Date(`time`) >= '%s' and Date(`time`) < '%s'", yestoday.Format(staticfunc.TIMEFORMAT), today.Format(staticfunc.TIMEFORMAT)), &money)
//	yestodaycard := money.Money

//	//! 获得今日抽水
//	Game_DB.GetOneData(fmt.Sprintf("select sum(gold) as money from log_agentgold where `time` >= %d", today.Unix()), &money)
//	todaylr := money.Money

//	//! 获得昨日抽水
//	Game_DB.GetOneData(fmt.Sprintf("select sum(gold) as money from log_agentgold where `time` >= %d and `time` < %d", yestoday.Unix(), today.Unix()), &money)
//	yestodaylr := money.Money

//	//! 总金币充值
//	Host_DB.GetOneData(fmt.Sprintf("select sum(money) as money from `order` where `status` = 1 and ptype = 2"), &money)
//	totalgold := money.Money

//	//! 总抽水
//	Game_DB.GetOneData(fmt.Sprintf("select sum(gold) as money from log_agentgold"), &money)
//	totallr := money.Money

//	//! 代理可体现
//	var moneyf SQL_MoneyF
//	HT_DB.GetOneData(fmt.Sprintf("select sum(score) as money from `fa_agent_user`"), &moneyf)
//	totalscore := moneyf.Money

//	result := fmt.Sprintf("在线人数%s人,充值房卡%d/%d,充值金币%d/%d,抽水%d/%d,总充值%d,总抽水%d,代理可提现%f", online, yestodaygold, todaygold, yestodaycard, todaycard, yestodaylr, todaylr, totalgold, totallr, totalscore)
//	w.Write([]byte(result))
//}

////! 得到今天的推广额
//func GetTodayWillGold(w http.ResponseWriter, req *http.Request) {
//	uid := lib.HF_Atoi(req.FormValue("uid"))
//	if uid == 0 {
//		w.Write([]byte("0"))
//		return
//	}

//	c := Redis.Get()
//	defer c.Close()

//	goldkey := fmt.Sprintf("%s_%d_%d_%d_%d", "willgold", uid, time.Now().Year(), time.Now().Month(), time.Now().Day())
//	value, err := redis.Float64(c.Do("GET", goldkey))
//	if err != nil {
//		value = 0
//	}

//	w.Write([]byte(fmt.Sprintf("%f", value)))
//}

////! 得到昨天的推广额
//func GetYesTodayWillGold(w http.ResponseWriter, req *http.Request) {
//	uid := lib.HF_Atoi(req.FormValue("uid"))
//	if uid == 0 {
//		w.Write([]byte("0"))
//		return
//	}

//	c := Redis.Get()
//	defer c.Close()

//	t := time.Unix(time.Now().Unix()-86400, 0)
//	goldkey := fmt.Sprintf("%s_%d_%d_%d_%d", "willgold", uid, t.Year(), t.Month(), t.Day())
//	value, err := redis.Float64(c.Do("GET", goldkey))
//	if err != nil {
//		value = 0
//	}

//	w.Write([]byte(fmt.Sprintf("%f", value)))
//}

////! 得到今天的推广额
//func GetTodayAddGold(w http.ResponseWriter, req *http.Request) {
//	uid := lib.HF_Atoi(req.FormValue("uid"))
//	if uid == 0 {
//		w.Write([]byte("0,0"))
//		return
//	}

//	c := Redis.Get()
//	defer c.Close()

//	goldkey := fmt.Sprintf("%s_%d_%d_%d_%d", "score1", uid, time.Now().Year(), time.Now().Month(), time.Now().Day())
//	value1, err := redis.Float64(c.Do("GET", goldkey))
//	if err != nil {
//		value1 = 0
//	}

//	goldkey = fmt.Sprintf("%s_%d_%d_%d_%d", "score2", uid, time.Now().Year(), time.Now().Month(), time.Now().Day())
//	value2, err := redis.Float64(c.Do("GET", goldkey))
//	if err != nil {
//		value2 = 0
//	}

//	w.Write([]byte(fmt.Sprintf("%f,%f", value1, value2)))
//}
