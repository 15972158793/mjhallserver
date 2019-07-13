package backstage

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"lib"
	"log"
	"rjmgr"
	"sync"
	"time"
)

//! sql队列
type SQL_Queue struct {
	Sql   string
	Value []byte
	V     bool
}

//! 通用结构
type SQL_Uid struct {
	Uid int
}

//! 这里定义表结构
type Fa_Agent_User struct {
	Id             int64   `json:"id"`
	Agid           int     `json:"agid"`
	Open_Id        string  `json:"open_id"`
	Union_Id       string  `json:"union_id"`
	Top_Group      string  `json:"top_group"`
	Score          float64 `json:"score"`   //! 可提现推广额
	T_Score        int     `json:"t_score"` //! 已提现推广额
	Deepin         int     `json:"deepin"`
	Card           int     `json:"card"`
	Password       string  `json:"password"`
	NickName       string  `json:"nickname"`
	Head           string  `json:"head"`
	Rating         string  `json:"rating"`
	Level          int     `json:"level"`
	Add_Time       int64   `json:"add_time"`
	Todaygold      float64 `json:"todaygold"`
	Yestodaygold   float64 `json:"yestodaygold"`
	TodayTime      int64   `json:"todaytime"`
	Parent         int     `json:"parent"`
	Name           string  `json:"name"`
	Alipay         string  `json:"alipay"`
	AliName        string  `json:"aliname"`
	Bankcard       string  `json:"bankcard"`
	Bankname       string  `json:"bankname"`
	Phone          string  `json:"phone"`
	AllCost        int64   `json:"allcost"`        //! 总抽水
	AllBills       int64   `json:"allbills"`       //! 总流水
	DayBills       int64   `json:"daybills"`       //! 日流水
	WeekBills      int64   `json:"weekbills"`      //! 周流水
	MonthBills     int64   `json:"monthbills"`     //! 月流水
	TimeBills      int64   `json:"timebills"`      //! 统计流水时间
	Commission     int     `json:"commission"`     //! 总佣金
	T_Commission   int     `json:"t_commission"`   //! 已提佣金
	Bills1         int64   `json:"bills1"`         //! 直属流水
	Bills2         int64   `json:"bills2"`         //! 下级流水
	TimeCommission int64   `json:"timecommission"` //! 计算佣金时间
	AreaScale      int     `json:"areascale"`      //! 区域比例
	AreaScore      float64 `json:"areascore"`      //! 区域收益
	AreaTScore     int     `json:"areatscore"`     //! 已提现收益
}

//! 加自己的流水
func (self *Fa_Agent_User) AddBills(bills int64) {
	self.AllBills += bills
	if self.TimeBills != 0 {
		time1 := time.Unix(self.TimeBills, 0)
		if time1.Year() != time.Now().Year() || time1.Month() != time.Now().Month() || time1.Day() != time.Now().Day() {
			self.DayBills = 0
		}
		if time1.Year() != time.Now().Year() || time1.Month() != time.Now().Month() {
			self.MonthBills = 0
		}
		y1, w1 := time1.ISOWeek()
		y, w := time.Now().ISOWeek()
		if y1 != y || w1 != w {
			self.WeekBills = 0
		}
	}

	self.WeekBills += bills
	self.MonthBills += bills
	self.DayBills += bills

	self.TimeBills = time.Now().Unix()
}

//! 下级加流水
func (self *Fa_Agent_User) AddCommission(bills int64, direct bool /*是否是直属*/) bool {
	if direct {
		self.Bills1 += bills
	} else {
		self.Bills2 += bills
	}

	if self.TimeCommission == 0 {
		self.TimeCommission = time.Now().Unix()
	} else {
		time1 := time.Unix(self.TimeCommission, 0)
		if time1.Year() != time.Now().Year() || time1.Month() != time.Now().Month() || time1.Day() != time.Now().Day() {
			if GetServer().ModeBills == 2 { //! 日结
				self.TimeCommission = time.Now().Unix()
				self.CountCommission()
				return true
			}
		}
		if time1.Year() != time.Now().Year() || time1.Month() != time.Now().Month() {
			if GetServer().ModeBills == 1 { //! 月结
				self.TimeCommission = time.Now().Unix()
				self.CountCommission()
				return true
			}
		}
		y1, w1 := time1.ISOWeek()
		y, w := time.Now().ISOWeek()
		if y1 != y || w1 != w {
			if GetServer().ModeBills == 0 { //! 周结
				self.TimeCommission = time.Now().Unix()
				self.CountCommission()
				return true
			}
		}
	}
	return false
}

//! 计算佣金
func (self *Fa_Agent_User) CountCommission() {
	//! 直属佣金
	total := self.Bills1 + self.Bills2
	tlevel := GetServer().GetBillsLevel(total)
	self.Commission += int((float64(self.Bills1) / 1000000.0) * float64(tlevel))

	//! 下级产生佣金
	xlevel := GetServer().GetBillsLevel(self.Bills2)
	self.Commission += int((float64(self.Bills2) / 1000000.0) * float64(tlevel-xlevel))

	self.Bills1 = 0
	self.Bills2 = 0
}

type Fa_Bind_Players struct {
	Id        int64   `json:"id"`
	Uid       int     `json:"uid"`
	Agid      int     `json:"agid"`
	Bind_time string  `json:"bind_time"`
	Score1    float64 `json:"score1"`
	Score2    float64 `json:"score2"`
}

type Fa_Children struct {
	Child []int `json:"child"`
}

func (self *Fa_Children) AddChild(uid int) {
	for i := 0; i < len(self.Child); i++ {
		if self.Child[i] == uid {
			return
		}
	}
	self.Child = append(self.Child, uid)
}

func (self *Fa_Children) DelChild(uid int) {
	for i := 0; i < len(self.Child); i++ {
		if self.Child[i] == uid {
			copy(self.Child[i:], self.Child[i+1:])
			self.Child = self.Child[:len(self.Child)-1]
			break
		}
	}
}

func (self *Fa_Children) HasChild(uid int) bool {
	for i := 0; i < len(self.Child); i++ {
		if self.Child[i] == uid {
			return true
		}
	}

	return false
}

type SQL_YJStatus struct {
	Id    int
	Quota int64
	YJ    int
	Level string
}

type SQL_YJType struct {
	Id   int
	Type int
}

type Config struct {
	Host       string   `json:"host"`
	Redis      string   `json:"redis"`
	GameDB     string   `json:"gamedb"`
	HTDB       string   `json:"htdb"`
	HostDB     string   `json:"hostdb"`
	GameIP     string   `json:"gameip"`
	LoginPort  int      `json:"loginport"`
	CenterPort int      `json:"centerport"`
	White      []string `json:"white"`     //! 白名单
	AgentMode  int      `json:"agentmode"` //! 代理模式 0抽水模式  1无限代模式
	Flag       string   `json:"flag"`      //! 唯一标识
	Sign       string   `json:"sign"`      //! 签名
}

type Server struct {
	Con               *Config //! 配置
	ShutDown          bool
	Wait              *sync.WaitGroup //! 同步阻塞
	Game_DB           *lib.DBServer
	HT_DB             *lib.DBServer
	Host_DB           *lib.DBServer
	Redis             *redis.Pool
	AgentUserChan     chan *SQL_Queue
	BindPlayerChan    chan *SQL_Queue
	ScoreLogChan      chan *SQL_Queue
	ParentLogChan     chan *SQL_Queue
	ParentCostLogChan chan *SQL_Queue
	Bills             []*SQL_YJStatus //! 佣金配置
	ModeBills         int             //! 0每周结算  1每月结算  2每日结算
	Tax               int             //! 新大区模式税点
}

var serverSingleton *Server = nil

//! 得到服务器指针
func GetServer() *Server {
	if serverSingleton == nil {
		serverSingleton = new(Server)
		serverSingleton.Con = new(Config)
		serverSingleton.Game_DB = new(lib.DBServer)
		serverSingleton.HT_DB = new(lib.DBServer)
		serverSingleton.Host_DB = new(lib.DBServer)
		serverSingleton.Wait = new(sync.WaitGroup)
	}

	return serverSingleton
}

func (self *Server) InitConfig() {
	//! 打开配置文件
	config, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("config err 1")
		return
	}
	err = json.Unmarshal(config, self.Con)
	if err != nil {
		log.Fatal("config err 2")
		return
	}
}

func (self *Server) Info() string {
	return lib.HF_JtoA(self.Con)
}

func (self *Server) Init() {
	self.Redis = lib.NewPool(self.Con.Redis, 2)

	if lib.HF_MD5(self.Con.Sign) != rjmgr.GetRJMgr().Flag {
		log.Fatal("sign err")
		return
	}

	self.Con.AgentMode = rjmgr.GetRJMgr().AgentMode

	//! 初始化数据库
	db := "root:"
	db += rjmgr.GetRJMgr().SQL
	db += "@tcp(127.0.0.1:3308)/"
	db += self.Con.GameDB
	db += "?charset=utf8&timeout=10s"
	if !self.Game_DB.Init(db) {
		log.Fatal("db1 fail")
		return
	}

	db = "root:"
	db += rjmgr.GetRJMgr().SQL
	db += "@tcp(127.0.0.1:3308)/"
	db += self.Con.HTDB
	db += "?charset=utf8&timeout=10s"
	if !self.HT_DB.Init(db) {
		log.Fatal("db2 fail")
		return
	}

	db = "root:"
	db += rjmgr.GetRJMgr().SQL
	db += "@tcp(127.0.0.1:3308)/"
	db += self.Con.HostDB
	db += "?charset=utf8&timeout=10s"
	if !self.Host_DB.Init(db) {
		log.Fatal("db3 fail")
		return
	}

	//! 载入无限代配置
	self.InitBills()
	self.Tax = 17

	rjmgr.GetRJMgr().InitIP(self.Con.Host)
}

//! 载入无限代配置
func (self *Server) InitBills() {
	{
		self.Bills = make([]*SQL_YJStatus, 0)
		var sql SQL_YJStatus
		res := self.Host_DB.GetAllData("select * from `yjstatus`", &sql)
		for i := 0; i < len(res); i++ {
			self.Bills = append(self.Bills, res[i].(*SQL_YJStatus))
		}
	}

	{
		var sql SQL_YJType
		self.Host_DB.GetOneData("select * from `yjtype`", &sql)
		self.ModeBills = sql.Type
	}
}

//! 确认档位
func (self *Server) GetBillsLevel(money int64) int {
	if money == 0 {
		return 0
	}

	if len(self.Bills) == 0 {
		return 0
	}

	for i := len(self.Bills) - 2; i >= 0; i-- {
		if money > self.Bills[i].Quota {
			return self.Bills[i+1].YJ
		}
	}
	return self.Bills[0].YJ
}

func (self *Server) Close() {
	self.ShutDown = true
	self.QueueAgentUser("", []byte(""), false)
	self.QueueBindPlayers("")
	self.QueueScoreLog("")
	self.QueueParentLog("")
	self.QueueParentCostLog("")

	//! 等所有存入数据库的指令执行完了
	self.Wait.Wait()

	self.HT_DB.Close()
	self.Host_DB.Close()
	self.Game_DB.Close()
	self.Redis.Close()

	log.Fatalln("server shutdown")
}

//! 得到agent_user
func (self *Server) GetAgentUser(agid int) *Fa_Agent_User {
	value := new(Fa_Agent_User)

	c := GetServer().Redis.Get()
	defer c.Close()

	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("fa_agent_user1_%d", agid)))
	if err == nil {
		json.Unmarshal(v, value)
		if value.Head == "" { //! 头像是空的
			name, head, _, _, ok := HF_GetPlayer(value.Agid, 1)
			if ok {
				value.NickName = name
				value.Head = head
				c.Do("SET", fmt.Sprintf("fa_agent_user1_%d", agid), lib.HF_JtoB(value))
				GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `nickname` = ?, `head` = '%s' where `agid` = %d", head, agid), []byte(name), true)
			}
		}
		c.Do("EXPIRE", fmt.Sprintf("fa_agent_user1_%d", agid), 86400*7)
		return value
	}

	if !self.HT_DB.GetOneData(fmt.Sprintf("select * from `fa_agent_user` where `agid` = %d", agid), value) {
		return value
	}

	if value.Agid == 0 {
		return value
	}

	if value.Head == "" { //! 头像是空的
		name, head, _, _, ok := HF_GetPlayer(value.Agid, 1)
		if ok {
			value.NickName = name
			value.Head = head
			GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `nickname` = ?, `head` = '%s' where `agid` = %d", head, agid), []byte(name), true)
		}
	}

	c.Do("SET", fmt.Sprintf("fa_agent_user1_%d", agid), lib.HF_JtoB(value))
	c.Do("EXPIRE", fmt.Sprintf("fa_agent_user1_%d", agid), 86400*7)

	return value
}

//! 插入agent_user
func (self *Server) InsertAgentUser(agid int, openid string, unionid string, top string, parent int) {
	c := GetServer().Redis.Get()
	defer c.Close()

	value := new(Fa_Agent_User)
	value.Agid = agid
	value.Open_Id = openid
	value.Union_Id = unionid
	value.Top_Group = top
	value.Deepin = 3
	if self.Con.AgentMode == 2 { //! 新抽水模式
		value.Rating = "40,15,7,0,0,0,0,0,0,0"
	} else {
		value.Rating = "35,10,5,5,5,5,5,5,5,5"
	}
	value.Parent = parent
	if self.Con.AgentMode == 1 { //! 无限代模式，进来都是代理
		value.Level = 1
	} else {
		value.Level = 0
	}
	c.Do("SET", fmt.Sprintf("fa_agent_user1_%d", agid), lib.HF_JtoB(value))
	c.Do("EXPIRE", fmt.Sprintf("fa_agent_user1_%d", agid), 86400*7)

	self.QueueAgentUser(fmt.Sprintf("insert into fa_agent_user(`agid`, `open_id`, `union_id`, `top_group`, `password`, `parent`, `level`, `nickname`, `rating`) values(%d, '%s', '%s', '%s', '', %d, %d, ?, '%s')", agid, openid, unionid, top, parent, value.Level, value.Rating), []byte(""), true)
}

//! 更新agent_user
//! 只更新redis,数据库单独更新
func (self *Server) SetAgentUser(value *Fa_Agent_User) {
	c := GetServer().Redis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("fa_agent_user1_%d", value.Agid), lib.HF_JtoB(value))
	c.Do("EXPIRE", fmt.Sprintf("fa_agent_user1_%d", value.Agid), 86400*7)
}

//! 得到bind_players
func (self *Server) GetBindPlayers(uid int) *Fa_Bind_Players {
	value := new(Fa_Bind_Players)

	c := GetServer().Redis.Get()
	defer c.Close()

	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("fa_bind_players1_%d", uid)))
	if err == nil {
		json.Unmarshal(v, value)
		c.Do("EXPIRE", fmt.Sprintf("fa_bind_players1_%d", uid), 86400*7)
		return value
	}

	if !self.HT_DB.GetOneData(fmt.Sprintf("select * from `fa_bind_players` where `uid` = %d", uid), value) {
		return value
	}

	if value.Agid == 0 {
		return value
	}

	c.Do("SET", fmt.Sprintf("fa_bind_players1_%d", uid), lib.HF_JtoB(value))
	c.Do("EXPIRE", fmt.Sprintf("fa_bind_players1_%d", uid), 86400*7)

	return value
}

//! 插入bind_players
func (self *Server) InsertBindPlayers(uid int, agid int) {
	c := GetServer().Redis.Get()
	defer c.Close()

	value := new(Fa_Bind_Players)
	value.Uid = uid
	value.Agid = agid
	value.Bind_time = time.Now().Format(lib.TIMEFORMAT)
	c.Do("SET", fmt.Sprintf("fa_bind_players1_%d", uid), lib.HF_JtoB(value))
	c.Do("EXPIRE", fmt.Sprintf("fa_bind_players1_%d", uid), 86400*7)

	self.QueueBindPlayers(fmt.Sprintf("INSERT into fa_bind_players(`uid`, `agid`, `bind_time`, `score1`, `score2`) values(%d, %d, '%s', %d, %d)", uid, agid, time.Now().Format(lib.TIMEFORMAT), 0, 0))

	all := self.GetChildren(agid)
	all.AddChild(uid)
	self.SetChildren(agid, all)
}

//! 更新bind_players
//! 只移动上级
func (self *Server) SetBindPlayers(uid int, agid int) {
	value := self.GetBindPlayers(uid)
	if value.Uid == 0 {
		return
	}

	c := GetServer().Redis.Get()
	defer c.Close()

	if value.Agid != agid {
		all := self.GetChildren(value.Agid)
		all.DelChild(uid)
		self.SetChildren(value.Agid, all)

		all = self.GetChildren(agid)
		all.AddChild(uid)
		self.SetChildren(agid, all)
	}

	value.Agid = agid
	c.Do("SET", fmt.Sprintf("fa_bind_players1_%d", uid), lib.HF_JtoB(value))
	c.Do("EXPIRE", fmt.Sprintf("fa_bind_players1_%d", uid), 86400*7)

	self.QueueBindPlayers(fmt.Sprintf("UPDATE `fa_bind_players` SET `agid` = %d WHERE uid = %d", agid, uid))
}

func (self *Server) AddBindPlayers(uid int, score1 float64, score2 float64) {
	value := self.GetBindPlayers(uid)
	if value.Uid == 0 {
		return
	}

	c := GetServer().Redis.Get()
	defer c.Close()

	value.Score1 += score1
	value.Score2 += score2
	c.Do("SET", fmt.Sprintf("fa_bind_players1_%d", uid), lib.HF_JtoB(value))
	c.Do("EXPIRE", fmt.Sprintf("fa_bind_players1_%d", uid), 86400*7)

	self.QueueBindPlayers(fmt.Sprintf("UPDATE `fa_bind_players` SET score1 = %f, score2 = %f WHERE uid = %d", value.Score1, value.Score2, uid))
}

//! 删除
func (self *Server) DelBindPlayers(uid int) {
	value := self.GetBindPlayers(uid)
	if value.Uid == 0 {
		return
	}

	c := GetServer().Redis.Get()
	defer c.Close()

	c.Do("DEL", fmt.Sprintf("fa_bind_players1_%d", uid))

	self.QueueBindPlayers(fmt.Sprintf("DELETE from fa_bind_players where `uid` = %d", uid))

	all := self.GetChildren(value.Agid)
	all.DelChild(uid)
	self.SetChildren(value.Agid, all)
}

//! 得到agid的所有下级
func (self *Server) GetChildren(agid int) *Fa_Children {
	value := new(Fa_Children)

	c := GetServer().Redis.Get()
	defer c.Close()

	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("fa_child1_%d", agid)))
	if err == nil {
		json.Unmarshal(v, value)
		c.Do("EXPIRE", fmt.Sprintf("fa_child1_%d", agid), 86400*7)
		return value
	}

	var sql_uid SQL_Uid
	res := self.HT_DB.GetAllData(fmt.Sprintf("select `uid` from `fa_bind_players` where `agid` = %d", agid), &sql_uid)
	for i := 0; i < len(res); i++ {
		value.Child = append(value.Child, res[i].(*SQL_Uid).Uid)
	}

	c.Do("SET", fmt.Sprintf("fa_child1_%d", agid), lib.HF_JtoB(value))
	c.Do("EXPIRE", fmt.Sprintf("fa_child1_%d", agid), 86400*7)

	return value
}

//! 设置agid的所有下级
func (self *Server) SetChildren(agid int, value *Fa_Children) {
	c := GetServer().Redis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("fa_child1_%d", agid), lib.HF_JtoB(value))
	c.Do("EXPIRE", fmt.Sprintf("fa_child1_%d", agid), 86400*7)
}

//! 今日新增绑定
func (self *Server) AddTodayChild(agid int) {
	c := GetServer().Redis.Get()
	defer c.Close()
	childkey := fmt.Sprintf("%s_%d_%d_%d_%d", "addchild", agid, time.Now().Year(), time.Now().Month(), time.Now().Day())
	thenum, err := redis.Int(c.Do("GET", childkey))
	if err != nil {
		thenum = 0
	}
	thenum += 1
	c.Do("SET", childkey, thenum)
	c.Do("EXPIRE", childkey, 86400)
}

//! 得到今日新增
func (self *Server) GetTodayChild(agid int) int {
	c := GetServer().Redis.Get()
	defer c.Close()
	childkey := fmt.Sprintf("%s_%d_%d_%d_%d", "addchild", agid, time.Now().Year(), time.Now().Month(), time.Now().Day())
	thenum, err := redis.Int(c.Do("GET", childkey))
	if err != nil {
		thenum = 0
	}
	return thenum
}

//! 本月新增绑定
func (self *Server) AddMonthChild(agid int) {
	c := GetServer().Redis.Get()
	defer c.Close()
	childkey := fmt.Sprintf("%s_%d_%d_%d", "addchildmonth", agid, time.Now().Year(), time.Now().Month())
	thenum, err := redis.Int(c.Do("GET", childkey))
	if err != nil {
		thenum = 0
	}
	thenum += 1
	c.Do("SET", childkey, thenum)
	c.Do("EXPIRE", childkey, 86400*31)
}

//! 得到本月新增
func (self *Server) GetMonthChild(agid int) int {
	c := GetServer().Redis.Get()
	defer c.Close()
	childkey := fmt.Sprintf("%s_%d_%d_%d", "addchildmonth", agid, time.Now().Year(), time.Now().Month())
	thenum, err := redis.Int(c.Do("GET", childkey))
	if err != nil {
		thenum = 0
	}
	return thenum
}

//! 本周新增绑定
func (self *Server) AddWeekChild(agid int) {
	c := GetServer().Redis.Get()
	defer c.Close()
	year, week := time.Now().ISOWeek()
	childkey := fmt.Sprintf("%s_%d_%d_%d", "addchildweek", agid, year, week)
	thenum, err := redis.Int(c.Do("GET", childkey))
	if err != nil {
		thenum = 0
	}
	thenum += 1
	c.Do("SET", childkey, thenum)
	c.Do("EXPIRE", childkey, 86400*31)
}

//! 得到本周新增
func (self *Server) GetWeekChild(agid int) int {
	c := GetServer().Redis.Get()
	defer c.Close()
	year, week := time.Now().ISOWeek()
	childkey := fmt.Sprintf("%s_%d_%d_%d", "addchildweek", agid, year, week)
	thenum, err := redis.Int(c.Do("GET", childkey))
	if err != nil {
		thenum = 0
	}
	return thenum
}

//! 得到自己的推广额
func (self *Server) GetMyMoney(uid int) float64 {
	c := GetServer().Redis.Get()
	defer c.Close()
	moneykey := fmt.Sprintf("mymoney_%d", uid)
	thenum, err := redis.Float64(c.Do("GET", moneykey))
	if err != nil {
		var money SQL_MoneyF
		self.HT_DB.GetOneData(fmt.Sprintf("select sum(change_score) as `money` from `score_log` where `operator_id` = %d and `agid` = %d and `action` = 1", uid, uid), &money)
		thenum = money.Money
	}
	c.Do("SET", moneykey, thenum)
	c.Do("EXPIRE", moneykey, 86400*7)
	return thenum
}

//! 增加自己推广额
func (self *Server) AddMyMoney(uid int, money float64) {
	c := GetServer().Redis.Get()
	defer c.Close()
	moneykey := fmt.Sprintf("mymoney_%d", uid)
	thenum, err := redis.Float64(c.Do("GET", moneykey))
	if err != nil {
		thenum = self.GetMyMoney(uid)
	}
	thenum += money
	c.Do("SET", moneykey, thenum)
	c.Do("EXPIRE", moneykey, 86400*7)
}

//! 今日新增活跃度
func (self *Server) AddTodayActive(uid int, active float64) float64 {
	c := GetServer().Redis.Get()
	defer c.Close()
	moneykey := fmt.Sprintf("%s_%d_%d_%d_%d", "willgold", uid, time.Now().Year(), time.Now().Month(), time.Now().Day())
	thenum, err := redis.Float64(c.Do("GET", moneykey))
	if err != nil {
		thenum = 0
	}
	thenum += active
	c.Do("SET", moneykey, thenum)
	c.Do("EXPIRE", moneykey, 86400*3)
	return thenum
}

//! 得到今日提现
func (self *Server) GetTodayActive(uid int) float64 {
	c := GetServer().Redis.Get()
	defer c.Close()
	moneykey := fmt.Sprintf("%s_%d_%d_%d_%d", "willgold", uid, time.Now().Year(), time.Now().Month(), time.Now().Day())
	thenum, err := redis.Float64(c.Do("GET", moneykey))
	if err != nil {
		thenum = 0
	}
	return thenum
}

//! 清除新增活跃度
func (self *Server) DelTodayActive(uid int) {
	c := GetServer().Redis.Get()
	defer c.Close()
	moneykey := fmt.Sprintf("%s_%d_%d_%d_%d", "willgold", uid, time.Now().Year(), time.Now().Month(), time.Now().Day())
	c.Do("DEL", moneykey)
	return
}

//! 今日新增提现
func (self *Server) AddTodayMoney(uid int) {
	c := GetServer().Redis.Get()
	defer c.Close()
	moneykey := fmt.Sprintf("%s_%d_%d_%d_%d", "todaymoney", uid, time.Now().Year(), time.Now().Month(), time.Now().Day())
	thenum, err := redis.Int(c.Do("GET", moneykey))
	if err != nil {
		thenum = 0
	}
	thenum += 1
	c.Do("SET", moneykey, thenum)
	c.Do("EXPIRE", moneykey, 86400)
}

//! 得到今日提现
func (self *Server) GetTodayMoney(uid int) int {
	c := GetServer().Redis.Get()
	defer c.Close()
	moneykey := fmt.Sprintf("%s_%d_%d_%d_%d", "todaymoney", uid, time.Now().Year(), time.Now().Month(), time.Now().Day())
	thenum, err := redis.Int(c.Do("GET", moneykey))
	if err != nil {
		thenum = 0
	}
	return thenum
}

//! 今日新增提现
func (self *Server) AddTodayUidMoney(uid int, gold float64) {
	c := GetServer().Redis.Get()
	defer c.Close()
	moneykey := fmt.Sprintf("%s_%d_%d_%d_%d", "todayuidmoney", uid, time.Now().Year(), time.Now().Month(), time.Now().Day())
	thenum, err := redis.Float64(c.Do("GET", moneykey))
	if err != nil {
		thenum = 0
	}
	thenum += gold
	c.Do("SET", moneykey, thenum)
	c.Do("EXPIRE", moneykey, 86400)
}

//! 得到今日提现
func (self *Server) GetTodayUidMoney(uid int) float64 {
	c := GetServer().Redis.Get()
	defer c.Close()
	moneykey := fmt.Sprintf("%s_%d_%d_%d_%d", "todayuidmoney", uid, time.Now().Year(), time.Now().Month(), time.Now().Day())
	thenum, err := redis.Float64(c.Do("GET", moneykey))
	if err != nil {
		thenum = 0
	}
	return thenum
}

//! 删除今日提现
func (self *Server) DelTodayMoney(uid int) {
	c := GetServer().Redis.Get()
	defer c.Close()
	moneykey := fmt.Sprintf("%s_%d_%d_%d_%d", "todaymoney", uid, time.Now().Year(), time.Now().Month(), time.Now().Day())
	thenum, err := redis.Int(c.Do("GET", moneykey))
	if err != nil {
		thenum = 0
	}
	if thenum > 0 {
		thenum--
		c.Do("SET", moneykey, thenum)
		c.Do("EXPIRE", moneykey, 86400)
	}
}

//!  得到总提现额
func (self *Server) GetTotalMoney() float64 {
	c := GetServer().Redis.Get()
	defer c.Close()
	moneykey := fmt.Sprintf("totalmoney")
	thenum, err := redis.Float64(c.Do("GET", moneykey))
	if err != nil {
		var money SQL_MoneyF
		self.HT_DB.GetOneData(fmt.Sprintf("select sum(score) as `money` from `fa_agent_user`"), &money)
		thenum = money.Money
	}
	c.Do("SET", moneykey, thenum)
	return thenum
}

//! 加总提现额
func (self *Server) AddTotalMoney(money float64) {
	c := GetServer().Redis.Get()
	defer c.Close()
	moneykey := fmt.Sprintf("totalmoney")
	thenum, err := redis.Float64(c.Do("GET", moneykey))
	if err != nil {
		thenum = self.GetTotalMoney()
	}
	thenum += money
	c.Do("SET", moneykey, thenum)
}

//!  得到总已提现额
func (self *Server) GetCostMoney() int {
	c := GetServer().Redis.Get()
	defer c.Close()
	moneykey := fmt.Sprintf("costmoney")
	thenum, err := redis.Int(c.Do("GET", moneykey))
	if err != nil {
		var money SQL_Money
		self.HT_DB.GetOneData(fmt.Sprintf("select sum(t_score) as `money` from `fa_agent_user`"), &money)
		thenum = money.Money
	}
	c.Do("SET", moneykey, thenum)
	return thenum
}

//! 加总已提现额
func (self *Server) AddCostMoney(money int) {
	c := GetServer().Redis.Get()
	defer c.Close()
	moneykey := fmt.Sprintf("costmoney")
	thenum, err := redis.Int(c.Do("GET", moneykey))
	if err != nil {
		thenum = self.GetCostMoney()
	}
	thenum += money
	c.Do("SET", moneykey, thenum)
}

//! 数据库队列
func (self *Server) QueueAgentUser(sql string, value []byte, v bool) {
	if self.AgentUserChan == nil {
		return
	}

	self.AgentUserChan <- &SQL_Queue{sql, value, v}
}

//! 数据上报
func (self *Server) RunQueueAgentUser() {
	self.AgentUserChan = make(chan *SQL_Queue, 5000000)
	self.Wait.Add(1)
	for msg := range self.AgentUserChan {
		if msg.Sql == "" {
			break
		}
		for i := 0; i < 5; i++ {
			if msg.V {
				_, err := self.HT_DB.DB.Exec(msg.Sql, msg.Value)
				if err == nil {
					break
				} else {
					lib.GetLogMgr().Output(lib.LOG_ERROR, "RunQueueAgentUser:", msg.Sql, ",err:", err)
				}
			} else {
				_, err := self.HT_DB.DB.Exec(msg.Sql)
				if err == nil {
					break
				} else {
					lib.GetLogMgr().Output(lib.LOG_ERROR, "RunQueueAgentUser:", msg.Sql, ",err:", err)
				}
			}
		}
	}
	self.Wait.Done()
	close(self.AgentUserChan)
	self.AgentUserChan = nil
}

//! 数据库队列
func (self *Server) QueueBindPlayers(sql string) {
	if self.BindPlayerChan == nil {
		return
	}

	self.BindPlayerChan <- &SQL_Queue{sql, []byte(""), false}
}

//! 数据上报
func (self *Server) RunBindPlayers() {
	self.BindPlayerChan = make(chan *SQL_Queue, 5000000)
	self.Wait.Add(1)
	for msg := range self.BindPlayerChan {
		if msg.Sql == "" {
			break
		}
		for i := 0; i < 5; i++ {
			if msg.V {
				_, err := self.HT_DB.DB.Exec(msg.Sql, msg.Value)
				if err == nil {
					break
				} else {
					lib.GetLogMgr().Output(lib.LOG_ERROR, "RunBindPlayers:", msg.Sql, ",err:", err)
				}
			} else {
				_, err := self.HT_DB.DB.Exec(msg.Sql)
				if err == nil {
					break
				} else {
					lib.GetLogMgr().Output(lib.LOG_ERROR, "RunBindPlayers:", msg.Sql, ",err:", err)
				}
			}
		}
	}
	self.Wait.Done()
	close(self.BindPlayerChan)
	self.BindPlayerChan = nil
}

//! 数据库队列
func (self *Server) QueueScoreLog(sql string) {
	if self.ScoreLogChan == nil {
		return
	}

	self.ScoreLogChan <- &SQL_Queue{sql, []byte(""), false}
}

//! 数据上报
func (self *Server) RunScoreLog() {
	self.ScoreLogChan = make(chan *SQL_Queue, 5000000)
	self.Wait.Add(1)
	for msg := range self.ScoreLogChan {
		if msg.Sql == "" {
			break
		}
		for i := 0; i < 5; i++ {
			if msg.V {
				_, err := self.HT_DB.DB.Exec(msg.Sql, msg.Value)
				if err == nil {
					break
				} else {
					lib.GetLogMgr().Output(lib.LOG_ERROR, "RunScoreLog:", msg.Sql, ",err:", err)
				}
			} else {
				_, err := self.HT_DB.DB.Exec(msg.Sql)
				if err == nil {
					break
				} else {
					lib.GetLogMgr().Output(lib.LOG_ERROR, "RunScoreLog:", msg.Sql, ",err:", err)
				}
			}
		}
	}
	self.Wait.Done()
	close(self.ScoreLogChan)
	self.ScoreLogChan = nil
}

//! 数据库队列
func (self *Server) QueueParentLog(sql string) {
	if self.ParentLogChan == nil {
		return
	}

	self.ParentLogChan <- &SQL_Queue{sql, []byte(""), false}
}

//! 数据上报
func (self *Server) RunParentLog() {
	self.ParentLogChan = make(chan *SQL_Queue, 5000000)
	self.Wait.Add(1)
	for msg := range self.ParentLogChan {
		if msg.Sql == "" {
			break
		}
		for i := 0; i < 5; i++ {
			if msg.V {
				_, err := self.HT_DB.DB.Exec(msg.Sql, msg.Value)
				if err == nil {
					break
				} else {
					lib.GetLogMgr().Output(lib.LOG_ERROR, "RunParentLog:", msg.Sql, ",err:", err)
				}
			} else {
				_, err := self.HT_DB.DB.Exec(msg.Sql)
				if err == nil {
					break
				} else {
					lib.GetLogMgr().Output(lib.LOG_ERROR, "RunParentLog:", msg.Sql, ",err:", err)
				}
			}
		}
	}
	self.Wait.Done()
	close(self.ParentLogChan)
	self.ParentLogChan = nil
}

//! 数据库队列
func (self *Server) QueueParentCostLog(sql string) {
	if self.ParentCostLogChan == nil {
		return
	}

	self.ParentCostLogChan <- &SQL_Queue{sql, []byte(""), false}
}

//! 数据上报
func (self *Server) RunParentCostLog() {
	self.ParentCostLogChan = make(chan *SQL_Queue, 5000000)
	self.Wait.Add(1)
	for msg := range self.ParentCostLogChan {
		if msg.Sql == "" {
			break
		}
		for i := 0; i < 5; i++ {
			if msg.V {
				_, err := self.HT_DB.DB.Exec(msg.Sql, msg.Value)
				if err == nil {
					break
				} else {
					lib.GetLogMgr().Output(lib.LOG_ERROR, "RunParentLog:", msg.Sql, ",err:", err)
				}
			} else {
				_, err := self.HT_DB.DB.Exec(msg.Sql)
				if err == nil {
					break
				} else {
					lib.GetLogMgr().Output(lib.LOG_ERROR, "RunParentLog:", msg.Sql, ",err:", err)
				}
			}
		}
	}
	self.Wait.Done()
	close(self.ParentCostLogChan)
	self.ParentCostLogChan = nil
}

func (self *Server) IsWhite(ip string) bool {
	if ip == "127.0.0.1" {
		return true
	}

	for i := 0; i < len(self.Con.White); i++ {
		if self.Con.White[i] == "0.0.0.0" {
			return true
		}

		if self.Con.White[i] == ip {
			return true
		}
	}

	return false
}
