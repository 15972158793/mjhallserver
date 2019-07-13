package gameserver

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"lib"
	"log"
	"rjmgr"
	//"net/http"
	"errors"
	"staticfunc"
	"strings"
	"sync"
	"time"
)

//! 分数日志
type SQL_Score struct {
	Id     int64
	Uid    int64
	Name   string
	Head   string
	Gameid int
	Room   int
	Time   int64
	Score  int
}

//! 代理日志
type SQL_Agent struct {
	Id       int64
	Uid      int64
	Agent    int
	GameType int
	Type     int
	Num      int
	Time     int64
}

//! 金币日志
type SQL_Gold struct {
	Id       int64
	Uid      int64
	Gold     int
	GameType int
	Time     int64
}

//! 流水日志
type SQL_Bills struct {
	Id       int64
	Uid      int64
	Num      int
	GameType int
	Time     int64
}

//! 消耗日志
type SQL_CostCard struct {
	Uid   int64
	Type  int
	Num   int
	Info  string
	Time  int64
	Table string
}

//! 房间日志
type SQL_RoomLog struct {
	Id   int64
	Uid  [6]int64
	IP   [6]string
	Win  [6]int
	Time int64
}

//! 五子棋日志
type SQL_WZQLog struct {
	Id   int64
	Uid1 int64
	Uid2 int64
	Gold int
	Time int64
}

//! 豹子王系统日志
type SQL_BZWLog struct {
	Id       int64
	Gold     int
	Time     int64
	GameType int
}

//! 超端日志
type SQL_SuperClientLog struct {
	Id       int64
	Uid      int64
	GameType int
	Next     string
	Time     int64
}

type Config struct {
	Id         int     `json:"id"`      //! 服务器id
	Host       string  `json:"host"`    //! 服务器外网ip
	InHost     string  `json:"inhost"`  //! 内网ip
	Type       int     `json:"type"`    //! 服务器类型
	Redis      string  `json:"redis"`   //! redis
	RcRedis    string  `json:"rcredis"` //! 记录redis
	RedisDB    int     `json:"redisdb"`
	RcRedisDB  int     `json:"rcredisdb"`
	Login      string  `json:"login"`  //! 登陆服ip
	Center     string  `json:"center"` //! 中心服务器ip
	DB         string  `json:"db"`     //! 数据库
	Log        []int64 `json:"log"`    //! 日志
	LogPro     int     `json:"logpro"`
	GM         []int64 `json:"gm"`
	IsNeed     int     `json:"need"`       //! 是否能要牌
	ScoreLog   int     `json:"scorelog"`   //! 分数日志
	PrintLevel int     `json:"printlevel"` //! 打印日志等级
	FileLevel  int     `json:"filelevel"`  //! 日志输出文件等级
	ManyGag    int     `json:"manygag"`    //! 百人场是否禁言
	MoneyMode  int     `json:"moneymode"`  //! 金币模式 0,1:100;1,1:1;2,1:10000
	WZQMoney   int     `json:"wzqmoney"`
	Flag       string  `json:"flag"` //! 唯一标识
	Sign       string  `json:"sign"` //! 签名
}

type Server struct {
	Con                *Config         //! 配置
	Wait               *sync.WaitGroup //! 同步阻塞
	ShutDown           bool            //! 是否正在执行关闭
	Redis              *redis.Pool
	RcRedis            *redis.Pool //! 对战记录redis
	DB                 *lib.DBServer
	Index              string
	Login              string
	SystemMoney        [3]int64 //! 系统奖池
	PlayerMoney        [3]int64 //! 玩家庄奖池
	YSZSYSMoney        [3]int64 //! 摇色子系统奖池
	YSZUSRMoney        [3]int64 //! 摇色子玩家奖池
	BrTTZMoney         [3]int64 //! 百人推筒子奖池
	BrTTZUsrMoney      [3]int64 //! 百人推筒子玩家奖池
	SxdbSysMoney       [3]int64 //! 神仙夺宝系统奖池
	SxdbUserMoney      [3]int64 //! 神仙夺宝玩家奖池
	MphSysMoney        [3]int64 //! 名品汇系统奖池
	MphUserMoney       [3]int64 //! 名品汇玩家奖池
	LHDMoney           [3]int64 //! 龙虎斗系统奖池
	TBUserMoney        [3]int64 //! 骰宝玩家奖池
	TBMoney            [3]int64 //! 骰宝系统奖池
	LzdbSysMoney       [3]int64 //! 龙珠夺宝系统奖池
	DfwSysMoney        [3]int64 //! 大富翁奖池
	FpjSysMoney        [3]int64 //! 翻牌机奖池
	FKFpjSysMoney      [3]int64 //! 疯狂翻牌机奖池
	SaiMaMoney         [3]int64 //! 赛马系统奖池
	HhdzSysMoney       [3]int64 //! 红黑大战奖池
	BrNNSysMoney       [3]int64 //! 百人牛牛系统奖池
	BrNNUserMoney      [3]int64 //! 百人牛牛玩家奖池
	YxxSysMoney        [3]int64 //! 鱼虾蟹系统奖池
	YxxUserMoney       [3]int64 //! 鱼虾蟹玩家奖池
	FishMoney          [3]int64 //!　捕鱼奖池
	LkpyMoney          [3]int64 //!　捕鱼奖池
	BJLSysMoney        [3]int64 //! 百家乐系统奖池
	BJLUserMoney       [3]int64 //! 百家乐玩家奖池
	SHZSysmoney        [3]int64 //! 水浒传奖池
	DwDyjlbSysMoney    [3]int64 //! 大赢家拉霸奖池
	DwJdqsSysMoney     [3]int64 //! 大赢家拉霸奖池
	RpcLogin           *lib.ClientPool
	LogChan            chan *SQL_Score
	AgentChan          chan *SQL_Agent
	AgentGoldChan      chan *SQL_Gold
	BillsChan          chan *SQL_Bills
	AgentBillsChan     chan *SQL_Gold
	CostCardChan       chan *SQL_CostCard
	RoomLogChan        chan *SQL_RoomLog
	WZQLogChan         chan *SQL_WZQLog
	BZWLogChan         chan *SQL_BZWLog
	SuperClientLogChan chan *SQL_SuperClientLog
}

var serverSingleton *Server = nil

//! 得到服务器指针
func GetServer() *Server {
	if serverSingleton == nil {
		serverSingleton = new(Server)
		serverSingleton.ShutDown = false
		serverSingleton.Con = new(Config)
		serverSingleton.Wait = new(sync.WaitGroup)
		serverSingleton.DB = new(lib.DBServer)
	}

	return serverSingleton
}

//! 初始化
func (self *Server) Init(index string, ssc bool /*是否开启时时彩*/) {
	//! 连接rpc
	ip := strings.Split(self.Con.Login, ":")
	_ip := ip[0] + ":1" + ip[1]
	self.RpcLogin = lib.CreateClientPool([]string{_ip}, _ip)

	if lib.HF_MD5(self.Con.Sign) != rjmgr.GetRJMgr().Flag {
		log.Fatal("sign err")
		return
	}

	//! 连接redis
	self.Redis = lib.NewPool(self.Con.Redis, self.Con.RedisDB)
	if self.Redis == nil {
		log.Fatal("redis err")
		return
	}

	self.RcRedis = lib.NewPool(self.Con.RcRedis, self.Con.RcRedisDB)
	if self.RcRedis == nil {
		log.Fatal("RcRedis err")
		return
	}

	//! 连接数据库
	db := "root:"
	db += rjmgr.GetRJMgr().SQL
	db += "@tcp(127.0.0.1:3308)/"
	db += self.Con.DB
	db += "?charset=utf8&timeout=10s"
	self.DB.Init(db)

	//! 载入系统奖池
	self.LoadSystemMoney()
	self.LoadPlayerMoney()
	self.LoadYSZSYSMoney()
	self.LoadYSZUSRMoney()
	self.LoadBrTTZMoney()
	self.LoadBrTTZUsrMoney()
	self.LoadSxdbSysMoney()
	self.LoadSxdbUserMoney()
	self.LoadMphSysMoney()
	self.LoadMphUserMoney()
	self.LoadLHDMoney()
	self.LoadTBMoney()
	self.LoadTBUsrMoney()
	self.LoadLzdbSysMoney()
	self.LoadDfwSysMoney()
	//self.LoadFpjSysMoney()
	self.LoadSaiMaMoney()
	self.LoadFKFpjSysMoney()
	self.LoadHhdzMoney()
	self.LoadBrNNMoney()
	self.LoadBrNNUserMoney()
	self.LoadYxxSysMoney()
	self.LoadYxxUserMoney()
	self.LoadFishMoney()
	self.LoadLkpyMoney()
	self.LoadBJLSysMoney()
	self.LoadBJLUserMoney()
	self.LoadSHZSysMoney()
	self.LoadDwDyjlbSysMoney()
	self.LoadJdqsSysMoney()

	//! 载入配置
	lib.GetManyMgr().Init(self.Redis)
	lib.GetManyMoneyMgr().Init(self.Redis)
	lib.GetSingleMgr().Init(self.Redis)
	//! 载入机器人
	lib.GetRobotMgr().Init(self.DB, self.Redis)
	//! 疯狂翻牌机
	lib.GetFPJMGr().Init(self.Redis, self.SendNotice)

	rjmgr.GetRJMgr().InitIP(self.Con.Host)

	if self.Con.Type >= 3 { //! 金币场的game
		if ssc {
			lib.GetSSCMgr().Run()
		}
		go self.OnTime()
	}

	//! 告诉服务器开启
	go self.ConnectLogin()
}

func (self *Server) InitConfig() {
	//! 打开配置文件
	config, err := ioutil.ReadFile("./game_config" + self.Index + ".json")
	if err != nil {
		log.Fatal("config err 1")
		return
	}
	err = json.Unmarshal(config, self.Con)
	if err != nil {
		log.Fatal("config err 2")
		return
	}

	lib.GetLogMgr().SetLevel(self.Con.PrintLevel, self.Con.FileLevel)

	if self.Con.WZQMoney == 0 {
		self.Con.WZQMoney = 1000000
	}

	//self.Con.Login = fmt.Sprintf("http://%s/servermsg", self.Con.Login)
}

//! 关闭服务器
func (self *Server) Close() {
	//! 告诉服务器关闭
	var msg staticfunc.Msg_GameServer
	msg.Id = self.Con.Id
	msg.InIp = ""
	msg.ExIp = ""
	GetServer().CallLogin("ServerMethod.ServerMsg", "gameserver", &msg)

	self.ShutDown = true
	self.SqlScoreLog(0, "", "", 0, 0, 0)
	self.SqlAgentLog(0, 0, 0, 0, 0)
	self.SqlAgentGoldLog(0, 0, 0)
	self.SqlAgentBillsLog(0, 0, 0)
	self.SqlCostCardLog(0, 0, 0, "", "")
	self.SqlBillsLog(0, 0, 0)
	self.SqlRoomLog(new(SQL_RoomLog))
	self.SqlWZQLog(new(SQL_WZQLog))
	self.SqlBZWLog(new(SQL_BZWLog))
	self.SqlSuperClientLog(new(SQL_SuperClientLog))

	self.Wait.Wait()

	self.Redis.Close()
	self.DB.Close()

	log.Fatalln("server shutdown")
}

func (self *Server) ConnectLogin() {
	//! 每1秒来一次心跳包
	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		var msg staticfunc.Msg_GameServer
		msg.Id = self.Con.Id
		msg.InIp = self.Con.InHost
		msg.ExIp = self.Con.Host
		msg.Type = self.Con.Type
		result, err := self.CallLogin("ServerMethod.ServerMsg", "gameserver", &msg)
		if err == nil {
			self.Con.MoneyMode = lib.HF_Atoi(string(result))

			if staticfunc.LIBDEBUG != 1 {
				ip := strings.Split(GetServer().Con.Host, ":")

				str := "http://"
				str += ip[0]
				str += ":8031/getinfo"

				res, err := lib.HF_Get(str, 5)
				if err != nil {
					log.Fatal("")
					return
				}
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				if err != nil || string(body) != "hello asset" {
					log.Fatal("")
					return
				}
			}
			break
		}

		if self.ShutDown {
			break
		}
	}

	//！ 关掉定时器
	ticker.Stop()
}

func (self *Server) OnTime() {
	//! 每1秒来一次心跳包
	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		if self.ShutDown {
			break
		}

		for i := 0; i < 3; i++ {
			lst := lib.GetFPJMGr().Run(i)
			if len(lst) > 0 {
				lstp := lib.GetFPJMGr().GetPerson(i)
				var msg Msg_Machine_Delete
				msg.Index = lst
				_msg := lib.HF_EncodeMsg("gamegoldfkfpjdelete", &msg, true)
				for k := 0; k < len(lstp); k++ {
					person := GetPersonMgr().GetPerson(lstp[k])
					if person == nil {
						continue
					}
					person.SendByteMsg(_msg)
				}
			}
		}
	}

	//！ 关掉定时器
	ticker.Stop()
}

//! 得到一个websocket处理句柄
func (self *Server) GetConnectHandler() websocket.Handler {
	connectHandler := func(ws *websocket.Conn) {
		if self.ShutDown { //! 关服了
			ws.Close()
			return
		}
		if staticfunc.GetIpBlackMgr().IsIp(lib.HF_GetHttpIP(ws.Request())) { //! ip在黑名单里
			ws.Close()
			return
		}
		session := lib.GetSessionMgr().GetNewSession(ws, OnReceive, OnClose)
		session.Run()
	}
	return websocket.Handler(connectHandler)
}

//! 获取一个数据
func (self *Server) DB_GetData(table string, uid int64) []byte {
	//! 先从redis里面得到该数据
	c := self.Redis.Get()
	defer c.Close()
	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("%s_%d", table, uid)))
	if err == nil { //! redis找得到，则直接返回redis里的数据
		return v
	}

	return []byte("")
}

//! 改变数据
func (self *Server) DB_SetData(table string, uid int64, value []byte) {
	//! 保存到redis
	c := self.Redis.Get()
	defer c.Close()
	key := fmt.Sprintf("%s_%d", table, uid)
	c.Do("SET", key, value)
	c.Do("EXPIRE", key, 86400*7)
}

//! 插入log数据
func (self *Server) InsertLog(gametype int, uid int64, _type int, num int, info string) {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][gametype]
	if !ok {
		return
	}
	self.SqlCostCardLog(uid, _type, num, info, csv["logtable"])
}

//! 插入record数据
func (self *Server) InsertRecord(gametype int, uid int64, info string, gold int /*这个记录实际的输赢*/) {
	c := self.RcRedis.Get()
	defer c.Close()

	csv, ok := staticfunc.GetCsvMgr().Data["game"][gametype]
	if !ok {
		return
	}

	table := fmt.Sprintf("%s_%d", csv["rctable"], uid)

	v, err := redis.Int(c.Do("LLEN", table))
	if err != nil {
		v = 0
	}

	for v >= lib.HF_Atoi(csv["maxrc"]) { //! 超过上限，删除
		c.Do("RPOP", table)
		v--
	}

	c.Do("LPUSH", table, info)
	c.Do("EXPIRE", table, 86400*2)

	//! 插入流水
	if gold != 0 {
		GetServer().SqlBillsLog(uid, gold, gametype)
	}
}

//func (self *Server) IsLog(uid int64) bool {
//	for i := 0; i < len(self.Con.Log); i++ {
//		if self.Con.Log[i] == uid {
//			return true
//		}
//	}

//	return false
//}

func (self *Server) IsGM(uid int64) bool {
	for i := 0; i < len(self.Con.GM); i++ {
		if self.Con.GM[i] == uid {
			return true
		}
	}

	return false
}

func (self *Server) IsAdmin(uid int64, admin int) bool {
	person := GetPersonMgr().GetPerson(uid)
	if person == nil {
		return false
	}

	return person.Admin&admin != 0
}

//! 触发好牌
func (self *Server) IsSuper() bool {
	if self.Con.LogPro == 0 {
		return false
	}

	if lib.HF_GetRandom(100) < self.Con.LogPro {
		return true
	}

	return false
}

//! 载入骰宝奖池
func (self *Server) LoadDSMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.TBMoney[0], _ = redis.Int64(c.Do("GET", "dsmoney0"))
	self.TBMoney[1], _ = redis.Int64(c.Do("GET", "dsmoney1"))
	self.TBMoney[2], _ = redis.Int64(c.Do("GET", "dsmoney2"))
}

//! 设置系统奖池
func (self *Server) SetDSMoney(index int, money int64) {
	self.TBMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("dsmoney%d", index), self.TBMoney[index])
}

//! 载入骰宝玩家奖池
func (self *Server) LoadDSUsrMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.TBUserMoney[0], _ = redis.Int64(c.Do("GET", "dsusermoney0"))
	self.TBUserMoney[1], _ = redis.Int64(c.Do("GET", "dsusermoney1"))
	self.TBUserMoney[2], _ = redis.Int64(c.Do("GET", "dsusermoney2"))
}

//! 设置系统奖池
func (self *Server) SetDSUserMoney(index int, money int64) {
	self.TBUserMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("dsusermoney%d", index), self.TBUserMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadSaiMaMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.SaiMaMoney[0], _ = redis.Int64(c.Do("GET", "saimamoney0"))
	self.SaiMaMoney[1], _ = redis.Int64(c.Do("GET", "saimamoney1"))
	self.SaiMaMoney[2], _ = redis.Int64(c.Do("GET", "saimamoney2"))
}

//! 设置系统奖池
func (self *Server) SetSaiMaMoney(index int, money int64) {
	self.SaiMaMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("saimamoney%d", index), self.SaiMaMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadHhdzMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.HhdzSysMoney[0], _ = redis.Int64(c.Do("GET", "hhdzsysmoney0"))
	self.HhdzSysMoney[1], _ = redis.Int64(c.Do("GET", "hhdzsysmoney1"))
	self.HhdzSysMoney[2], _ = redis.Int64(c.Do("GET", "hhdzsysmoney2"))
}

//! 设置系统奖池
func (self *Server) SetHhdzMoney(index int, money int64) {
	self.HhdzSysMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("hhdzsysmoney%d", index), self.HhdzSysMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadBrNNMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.BrNNSysMoney[0], _ = redis.Int64(c.Do("GET", "brnnsysmoney0"))
	self.BrNNSysMoney[1], _ = redis.Int64(c.Do("GET", "brnnsysmoney1"))
	self.BrNNSysMoney[2], _ = redis.Int64(c.Do("GET", "brnnsysmoney2"))
}

//! 设置系统奖池
func (self *Server) SetBrNNMoney(index int, money int64) {
	self.BrNNSysMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("brnnsysmoney%d", index), self.BrNNSysMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadBrNNUserMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.BrNNUserMoney[0], _ = redis.Int64(c.Do("GET", "brnnusermoney0"))
	self.BrNNUserMoney[1], _ = redis.Int64(c.Do("GET", "brnnusermoney1"))
	self.BrNNUserMoney[2], _ = redis.Int64(c.Do("GET", "brnnusermoney2"))
}

//! 设置系统奖池
func (self *Server) SetBrNNUserMoney(index int, money int64) {
	self.BrNNUserMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("brnnusermoney%d", index), self.BrNNUserMoney[index])
}

func (self *Server) LoadYxxSysMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.YxxSysMoney[0], _ = redis.Int64(c.Do("GET", "yxxsysmoney0"))
	self.YxxSysMoney[1], _ = redis.Int64(c.Do("GET", "yxxsysmoney1"))
	self.YxxSysMoney[2], _ = redis.Int64(c.Do("GET", "yxxsysmoney2"))
}

//! 设置系统奖池
func (self *Server) SetYxxSysMoney(index int, money int64) {
	self.YxxSysMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("yxxsysmoney%d", index), self.YxxSysMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadFishMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.FishMoney[0], _ = redis.Int64(c.Do("GET", "fishmoney0"))
	self.FishMoney[1], _ = redis.Int64(c.Do("GET", "fishmoney1"))
	self.FishMoney[2], _ = redis.Int64(c.Do("GET", "fishmoney2"))
}

//! 设置系统奖池
func (self *Server) SetFishMoney(index int, money int64) {
	self.FishMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("fishmoney%d", index), self.FishMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadLkpyMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.LkpyMoney[0], _ = redis.Int64(c.Do("GET", "lkpymoney0"))
	self.LkpyMoney[1], _ = redis.Int64(c.Do("GET", "lkpymoney1"))
	self.LkpyMoney[2], _ = redis.Int64(c.Do("GET", "lkpymoney2"))
}

//! 设置系统奖池
func (self *Server) SetLkpyMoney(index int, money int64) {
	self.LkpyMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("lkpymoney%d", index), self.LkpyMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadSHZSysMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.SHZSysmoney[0], _ = redis.Int64(c.Do("GET", "shzsysmoney0"))
	self.SHZSysmoney[1], _ = redis.Int64(c.Do("GET", "shzsysmoney1"))
	self.SHZSysmoney[2], _ = redis.Int64(c.Do("GET", "shzsysmoney2"))
}

//! 设置系统奖池
func (self *Server) SetSHZSysMoney(index int, money int64) {
	self.SHZSysmoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("shzsysmoney%d", index), self.SHZSysmoney[index])
}

//! 载入系统奖池
func (self *Server) LoadBJLSysMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.BJLSysMoney[0], _ = redis.Int64(c.Do("GET", "bjlsysmoney0"))
	self.BJLSysMoney[1], _ = redis.Int64(c.Do("GET", "bjlsysmoney1"))
	self.BJLSysMoney[2], _ = redis.Int64(c.Do("GET", "bjlsysmoney2"))
}

//! 设置系统奖池
func (self *Server) SetBJLSysMoney(index int, money int64) {
	self.BJLSysMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("bjlsysmoney%d", index), self.BJLSysMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadBJLUserMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.BJLUserMoney[0], _ = redis.Int64(c.Do("GET", "bjlusermoney0"))
	self.BJLUserMoney[1], _ = redis.Int64(c.Do("GET", "bjlusermoney1"))
	self.BJLUserMoney[2], _ = redis.Int64(c.Do("GET", "bjlusermoney2"))
}

//! 设置系统奖池
func (self *Server) SetBJLUserMoney(index int, money int64) {
	self.BJLUserMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("bjlusermoney%d", index), self.BJLUserMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadYxxUserMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.YxxUserMoney[0], _ = redis.Int64(c.Do("GET", "yxxusermoney0"))
	self.YxxUserMoney[1], _ = redis.Int64(c.Do("GET", "yxxusermoney1"))
	self.YxxUserMoney[2], _ = redis.Int64(c.Do("GET", "yxxusermoney2"))
}

//! 设置系统奖池
func (self *Server) SetYxxUserMoney(index int, money int64) {
	self.YxxUserMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("yxxusermoney%d", index), self.YxxUserMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadLzdbSysMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.LzdbSysMoney[0], _ = redis.Int64(c.Do("GET", "lzdbsysmoney0"))
	self.LzdbSysMoney[1], _ = redis.Int64(c.Do("GET", "lzdbsysmoney1"))
	self.LzdbSysMoney[2], _ = redis.Int64(c.Do("GET", "lzdbsysmoney2"))
}

//! 设置系统奖池
func (self *Server) SetLzdbSysMoney(index int, money int64) {
	self.LzdbSysMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("lzdbsysmoney%d", index), self.LzdbSysMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadSystemMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.SystemMoney[0], _ = redis.Int64(c.Do("GET", "systemmoney0"))
	self.SystemMoney[1], _ = redis.Int64(c.Do("GET", "systemmoney1"))
	self.SystemMoney[2], _ = redis.Int64(c.Do("GET", "systemmoney2"))
}

//! 设置系统奖池
func (self *Server) SetSystemMoney(index int, money int64) {
	self.SystemMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("systemmoney%d", index), self.SystemMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadDfwSysMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.DfwSysMoney[0], _ = redis.Int64(c.Do("GET", "dfwsysmoney0"))
	self.DfwSysMoney[1], _ = redis.Int64(c.Do("GET", "dfwsysmoney1"))
	self.DfwSysMoney[2], _ = redis.Int64(c.Do("GET", "dfwsysmoney2"))
}

//! 设置系统奖池
func (self *Server) SetDfwSysMoney(index int, money int64) {
	self.DfwSysMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("dfwsysmoney%d", index), self.DfwSysMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadFpjSysMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.FpjSysMoney[0], _ = redis.Int64(c.Do("GET", "fpjsysmoney0"))
	self.FpjSysMoney[1], _ = redis.Int64(c.Do("GET", "fpjsysmoney1"))
	self.FpjSysMoney[2], _ = redis.Int64(c.Do("GET", "fpjsysmoney2"))
}

//! 设置系统奖池
func (self *Server) SetFpjSysMoney(index int, money int64) {
	self.FpjSysMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("fpjsysmoney%d", index), self.FpjSysMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadFKFpjSysMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.FKFpjSysMoney[0], _ = redis.Int64(c.Do("GET", "fkfpjsysmoney0"))
	self.FKFpjSysMoney[1], _ = redis.Int64(c.Do("GET", "fkfpjsysmoney1"))
	self.FKFpjSysMoney[2], _ = redis.Int64(c.Do("GET", "fkfpjsysmoney2"))
}

//! 设置系统奖池
func (self *Server) SetFKFpjSysMoney(index int, money int64) {
	self.FKFpjSysMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("fkfpjsysmoney%d", index), self.FKFpjSysMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadLHDMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.LHDMoney[0], _ = redis.Int64(c.Do("GET", "lhdmoney0"))
	self.LHDMoney[1], _ = redis.Int64(c.Do("GET", "lhdmoney1"))
	self.LHDMoney[2], _ = redis.Int64(c.Do("GET", "lhdmoney2"))
}

//! 设置系统奖池
func (self *Server) SetLHDMoney(index int, money int64) {
	self.LHDMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("lhdmoney%d", index), self.LHDMoney[index])
}

//! 载入玩家奖池
func (self *Server) LoadPlayerMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.PlayerMoney[0], _ = redis.Int64(c.Do("GET", "playermoney0"))
	self.PlayerMoney[1], _ = redis.Int64(c.Do("GET", "playermoney1"))
	self.PlayerMoney[2], _ = redis.Int64(c.Do("GET", "playermoney2"))
}

//! 载入骰宝奖池
func (self *Server) LoadTBMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.TBMoney[0], _ = redis.Int64(c.Do("GET", "tbmoney0"))
	self.TBMoney[1], _ = redis.Int64(c.Do("GET", "tbmoney0"))
	self.TBMoney[2], _ = redis.Int64(c.Do("GET", "tbmoney0"))
}

//! 设置系统奖池
func (self *Server) SetTBMoney(index int, money int64) {
	self.TBMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("tbmoney%d", index), self.TBMoney[index])
}

//! 设置玩家奖池
func (self *Server) SetPlayerMoney(index int, money int64) {
	self.PlayerMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("playermoney%d", index), self.PlayerMoney[index])
}

//! 载入骰宝玩家奖池
func (self *Server) LoadTBUsrMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.TBUserMoney[0], _ = redis.Int64(c.Do("GET", "tbusermoney0"))
	self.TBUserMoney[1], _ = redis.Int64(c.Do("GET", "tbusermoney1"))
	self.TBUserMoney[2], _ = redis.Int64(c.Do("GET", "tbusermoney2"))
}

//! 设置系统奖池
func (self *Server) SetTBUserMoney(index int, money int64) {
	self.TBUserMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("tbusermoney%d", index), self.TBUserMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadYSZSYSMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.YSZSYSMoney[0], _ = redis.Int64(c.Do("GET", "yszsysmoney0"))
	self.YSZSYSMoney[1], _ = redis.Int64(c.Do("GET", "yszsysmoney1"))
	self.YSZSYSMoney[2], _ = redis.Int64(c.Do("GET", "yszsysmoney2"))
}

//! 设置系统奖池
func (self *Server) SetYSZSYSMoney(index int, money int64) {
	self.YSZSYSMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("yszsysmoney%d", index), self.YSZSYSMoney[index])
}

//! 载入玩家奖池
func (self *Server) LoadYSZUSRMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.YSZUSRMoney[0], _ = redis.Int64(c.Do("GET", "yszusrmoney0"))
	self.YSZUSRMoney[1], _ = redis.Int64(c.Do("GET", "yszusrmoney1"))
	self.YSZUSRMoney[2], _ = redis.Int64(c.Do("GET", "yszusrmoney2"))
}

//! 设置玩家奖池
func (self *Server) SetYSZUSRMoney(index int, money int64) {
	self.YSZUSRMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("yszusrmoney%d", index), self.YSZUSRMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadSxdbSysMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.SxdbSysMoney[0], _ = redis.Int64(c.Do("GET", "sxdbsysmoney0"))
	self.SxdbSysMoney[1], _ = redis.Int64(c.Do("GET", "sxdbsysmoney1"))
	self.SxdbSysMoney[2], _ = redis.Int64(c.Do("GET", "sxdbsysmoney2"))
}

//! 设置系统奖池
func (self *Server) SetSxdbSysMoney(index int, money int64) {
	self.SxdbSysMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("sxdbsysmoney%d", index), self.SxdbSysMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadMphSysMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.MphSysMoney[0], _ = redis.Int64(c.Do("GET", "mphsysmoney0"))
	self.MphSysMoney[1], _ = redis.Int64(c.Do("GET", "mphsysmoney1"))
	self.MphSysMoney[2], _ = redis.Int64(c.Do("GET", "mphsysmoney2"))
}

//! 设置系统奖池
func (self *Server) SetMphSysMoney(index int, money int64) {
	self.MphSysMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("mphsysmoney%d", index), self.MphSysMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadJdqsSysMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.DwJdqsSysMoney[0], _ = redis.Int64(c.Do("GET", "dwjdqssysmoney"))
	self.DwJdqsSysMoney[1], _ = redis.Int64(c.Do("GET", "dwjdqssysmoney"))
	self.DwJdqsSysMoney[2], _ = redis.Int64(c.Do("GET", "dwjdqssysmoney"))
}

//! 设置系统奖池
func (self *Server) SetDwJdqsSysMoney(index int, money int64) {
	self.DwJdqsSysMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("dwjdqssysmoney%d", index), self.DwJdqsSysMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadDwDyjlbSysMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.DwDyjlbSysMoney[0], _ = redis.Int64(c.Do("GET", "dwdyjlbsysmoney0"))
	self.DwDyjlbSysMoney[1], _ = redis.Int64(c.Do("GET", "dwdyjlbsysmoney1"))
	self.DwDyjlbSysMoney[2], _ = redis.Int64(c.Do("GET", "dwdyjlbsysmoney2"))
}

//! 设置系统奖池
func (self *Server) SetDwDyjlbSysMoney(index int, money int64) {
	self.DwDyjlbSysMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("dwdyjlbsysmoney%d", index), self.DwDyjlbSysMoney[index])
}

//! 载入系统奖池
func (self *Server) LoadSxdbUserMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.SxdbUserMoney[0], _ = redis.Int64(c.Do("GET", "sxdbusermoney0"))
	self.SxdbUserMoney[1], _ = redis.Int64(c.Do("GET", "sxdbusermoney1"))
	self.SxdbUserMoney[2], _ = redis.Int64(c.Do("GET", "sxdbusermoney2"))
}

func (self *Server) LoadMphUserMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.MphUserMoney[0], _ = redis.Int64(c.Do("GET", "mphusermoney0"))
	self.MphUserMoney[1], _ = redis.Int64(c.Do("GET", "mphusermoney1"))
	self.MphUserMoney[2], _ = redis.Int64(c.Do("GET", "mphusermoney2"))
}

//! 设置系统奖池
func (self *Server) SetSxdbUserMoney(index int, money int64) {
	self.SxdbUserMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("sxdbusermoney%d", index), self.SxdbUserMoney[index])
}

//! 设置玩家奖池
func (self *Server) SetMphUserMoney(index int, money int64) {
	self.MphUserMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("mphusermoney%d", index), self.MphUserMoney[index])
}

//! 载入推筒子奖池
func (self *Server) LoadBrTTZMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.BrTTZMoney[0], _ = redis.Int64(c.Do("GET", "brttzmoney0"))
	self.BrTTZMoney[1], _ = redis.Int64(c.Do("GET", "brttzmoney1"))
	self.BrTTZMoney[2], _ = redis.Int64(c.Do("GET", "brttzmoney2"))
}

//! 设置系统奖池
func (self *Server) SetBrTTZMoney(index int, money int64) {
	self.BrTTZMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("brttzmoney%d", index), self.BrTTZMoney[index])
}

//! 载入推筒子玩家奖池
func (self *Server) LoadBrTTZUsrMoney() {
	c := self.RcRedis.Get()
	defer c.Close()

	self.BrTTZUsrMoney[0], _ = redis.Int64(c.Do("GET", "brttzusrmoney0"))
	self.BrTTZUsrMoney[1], _ = redis.Int64(c.Do("GET", "brttzusrmoney1"))
	self.BrTTZUsrMoney[2], _ = redis.Int64(c.Do("GET", "brttzusrmoney2"))
}

//! 设置系统奖池
func (self *Server) SetBrTTZUsrMoney(index int, money int64) {
	self.BrTTZUsrMoney[index] = money

	c := self.RcRedis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("brttzusrmoney%d", index), self.BrTTZUsrMoney[index])
}

//! 调用登陆服务器
func (self *Server) CallLogin(method string, msghead string, v interface{}) ([]byte, error) {
	if self.RpcLogin == nil {
		return []byte(""), errors.New("RpcLogin is nil")
	}

	data := lib.HF_EncodeMsg(msghead, v, false)

	err, client := self.RpcLogin.RandomGetConn()
	if err != nil {
		return []byte(""), err
	}

	var reply []byte
	err = client.Call(method, staticfunc.Rpc_Args{data}, &reply)
	if err != nil {
		return []byte(""), err
	}

	return reply, nil
}

//! 发送公告
func (self *Server) SendNotice(content string) {
	var msg staticfunc.Msg_RichNotice
	msg.Content = content
	self.CallLogin("ServerMethod.ServerMsg", "richnotice", &msg)
}

//! 数据库日志
func (self *Server) SqlScoreLog(uid int64, name string, head string, gameid int, room int, score int) {
	if self.LogChan == nil {
		return
	}

	if len(self.LogChan) >= 10000 {
		return
	}

	if uid == 0 {
		self.LogChan <- &SQL_Score{0, uid, name, head, gameid, room, time.Now().Unix(), score}
		return
	}

	if self.Con.ScoreLog == 0 { //! 没有开启分数统计并且不是作弊者
		return
	}

	self.LogChan <- &SQL_Score{0, uid, name, head, gameid, room, time.Now().Unix(), score}
}

//! 数据上报
func (self *Server) RunSqlScoreLog() {
	self.LogChan = make(chan *SQL_Score, 10000)
	self.Wait.Add(1)
	for msg := range self.LogChan {
		if msg.Uid == 0 {
			break
		}
		for i := 0; i < 5; i++ {
			sql := fmt.Sprintf("insert into `log_score`(`uid`, `name`, `head`, `gameid`, `room`, `time`, `score`) values (%d, ?, '%s', %d, %d, %d, %d)", msg.Uid, msg.Head, msg.Gameid, msg.Room, msg.Time, msg.Score)
			_, err := GetServer().DB.DB.Exec(sql, msg.Name)
			if err == nil {
				break
			}
		}
	}
	self.Wait.Done()
	close(self.LogChan)
	self.LogChan = nil
}

//! 代理日志
func (self *Server) SqlAgentLog(uid int64, agent int, gametype int, _type int, num int) {
	if self.AgentChan == nil {
		return
	}

	if len(self.AgentChan) >= 10000 {
		return
	}

	self.AgentChan <- &SQL_Agent{0, uid, agent, gametype, _type, num, time.Now().Unix()}
}

//! 数据上报
func (self *Server) RunSqlAgentLog() {
	self.AgentChan = make(chan *SQL_Agent, 10000)
	self.Wait.Add(1)
	for msg := range self.AgentChan {
		if msg.Uid == 0 {
			break
		}
		for i := 0; i < 5; i++ {
			sql := fmt.Sprintf("insert into `log_agent`(`uid`, `agent`, `gametype`, `type`, `num`, `time`) values (%d, %d, %d, %d, %d, %d)", msg.Uid, msg.Agent, msg.GameType, msg.Type, msg.Num, msg.Time)
			_, err := GetServer().DB.DB.Exec(sql)
			if err == nil {
				break
			}
		}
	}
	self.Wait.Done()
	close(self.AgentChan)
	self.AgentChan = nil
}

//! 代理金币日志
func (self *Server) SqlAgentGoldLog(uid int64, gold int, gametype int) {
	if self.AgentGoldChan == nil {
		return
	}

	if len(self.AgentGoldChan) >= 100000 {
		return
	}

	self.AgentGoldChan <- &SQL_Gold{0, uid, gold, gametype, time.Now().Unix()}
}

//! 数据上报
func (self *Server) RunSqlAgentGoldLog() {
	self.AgentGoldChan = make(chan *SQL_Gold, 100000)
	self.Wait.Add(1)
	for msg := range self.AgentGoldChan {
		if msg.Uid == 0 {
			break
		}
		for i := 0; i < 5; i++ {
			sql := fmt.Sprintf("insert into `log_agentgold`(`uid`, `gold`, `gametype`, `time`) values (%d, %d, %d, %d)", msg.Uid, msg.Gold, msg.GameType, msg.Time)
			_, err := GetServer().DB.DB.Exec(sql)
			if err == nil {
				break
			}
		}
	}
	self.Wait.Done()
	close(self.AgentGoldChan)
	self.AgentGoldChan = nil
}

//! 代理金币日志
func (self *Server) SqlAgentBillsLog(uid int64, gold int, gametype int) {
	if self.AgentBillsChan == nil {
		return
	}

	if len(self.AgentBillsChan) >= 100000 {
		return
	}

	self.AgentBillsChan <- &SQL_Gold{0, uid, gold, gametype, time.Now().Unix()}
}

//! 数据上报
func (self *Server) RunSqlAgentBillsLog() {
	self.AgentBillsChan = make(chan *SQL_Gold, 100000)
	self.Wait.Add(1)
	for msg := range self.AgentBillsChan {
		if msg.Uid == 0 {
			break
		}
		for i := 0; i < 5; i++ {
			sql := fmt.Sprintf("insert into `log_agentbills`(`uid`, `gold`, `gametype`, `time`) values (%d, %d, %d, %d)", msg.Uid, msg.Gold, msg.GameType, msg.Time)
			_, err := GetServer().DB.DB.Exec(sql)
			if err == nil {
				break
			}
		}
	}
	self.Wait.Done()
	close(self.AgentBillsChan)
	self.AgentBillsChan = nil
}

//! 流水日志
func (self *Server) SqlBillsLog(uid int64, num int, gametype int) {
	if self.BillsChan == nil {
		return
	}

	if len(self.BillsChan) >= 100000 {
		return
	}

	self.BillsChan <- &SQL_Bills{0, uid, num, gametype, time.Now().Unix()}
}

//! 数据上报
func (self *Server) RunSqlBillsLog() {
	self.BillsChan = make(chan *SQL_Bills, 100000)
	self.Wait.Add(1)
	for msg := range self.BillsChan {
		if msg.Uid == 0 {
			break
		}
		for i := 0; i < 5; i++ {
			sql := fmt.Sprintf("insert into `log_bills`(`uid`, `num`, `gametype`, `time`) values (%d, %d, %d, %d)", msg.Uid, msg.Num, msg.GameType, msg.Time)
			_, err := GetServer().DB.DB.Exec(sql)
			if err == nil {
				break
			}
		}
	}
	self.Wait.Done()
	close(self.BillsChan)
	self.BillsChan = nil
}

//! 消耗房卡日志
func (self *Server) SqlCostCardLog(uid int64, _type int, num int, info string, table string) {
	if self.CostCardChan == nil {
		return
	}

	if len(self.CostCardChan) >= 10000 {
		return
	}

	self.CostCardChan <- &SQL_CostCard{uid, _type, num, info, time.Now().Unix(), table}
}

//! 数据上报
func (self *Server) RunSqlCostCardLog() {
	self.CostCardChan = make(chan *SQL_CostCard, 10000)
	self.Wait.Add(1)
	for msg := range self.CostCardChan {
		if msg.Uid == 0 {
			break
		}
		var log staticfunc.Log_Base
		log.Uid = msg.Uid
		log.Type = msg.Type
		log.Num = msg.Num
		log.Info = msg.Info
		log.Creation_time = msg.Time
		for i := 0; i < 5; i++ {
			if lib.InsertTable(msg.Table, &log, 1, GetServer().DB) > 0 {
				break
			}
		}
	}
	self.Wait.Done()
	close(self.CostCardChan)
	self.CostCardChan = nil
}

//! 房间日志
func (self *Server) SqlRoomLog(node *SQL_RoomLog) {
	if self.RoomLogChan == nil {
		return
	}

	if len(self.RoomLogChan) >= 10000 {
		return
	}

	self.RoomLogChan <- node
}

//! 数据上报
func (self *Server) RunSqlRoomLog() {
	self.RoomLogChan = make(chan *SQL_RoomLog, 10000)
	self.Wait.Add(1)
	for msg := range self.RoomLogChan {
		if msg.Id == 0 {
			break
		}
		for i := 0; i < 5; i++ {
			sql := fmt.Sprintf("insert into `log_room`(`p1`, `p2`, `p3`, `p4`, `p5`, `p6`, `ip1`, `ip2`, `ip3`, `ip4`, `ip5`, `ip6`, `win1`, `win2`, `win3`, `win4`, `win5`, `win6`, `time`) values (%d, %d, %d, %d, %d, %d, '%s', '%s', '%s', '%s', '%s', '%s', %d, %d, %d, %d, %d, %d, %d)",
				msg.Uid[0], msg.Uid[1], msg.Uid[2], msg.Uid[3], msg.Uid[4], msg.Uid[5], msg.IP[0], msg.IP[1], msg.IP[2], msg.IP[3], msg.IP[4], msg.IP[5], msg.Win[0], msg.Win[1], msg.Win[2], msg.Win[3], msg.Win[4], msg.Win[5], msg.Time)
			_, err := GetServer().DB.DB.Exec(sql)
			if err == nil {
				break
			}
		}
	}
	self.Wait.Done()
	close(self.RoomLogChan)
	self.RoomLogChan = nil
}

//! 五子棋日志
func (self *Server) SqlWZQLog(node *SQL_WZQLog) {
	if self.WZQLogChan == nil {
		return
	}

	if len(self.WZQLogChan) >= 10000 {
		return
	}

	self.WZQLogChan <- node
}

//! 数据上报
func (self *Server) RunSqlWZQLog() {
	self.WZQLogChan = make(chan *SQL_WZQLog, 10000)
	self.Wait.Add(1)
	for msg := range self.WZQLogChan {
		if msg.Id == 0 {
			break
		}
		for i := 0; i < 5; i++ {
			if lib.InsertTable("log_wzq", msg, 1, GetServer().DB) > 0 {
				break
			}
		}
	}
	self.Wait.Done()
	close(self.WZQLogChan)
	self.WZQLogChan = nil
}

//! 豹子王日志
func (self *Server) SqlBZWLog(node *SQL_BZWLog) {
	if self.BZWLogChan == nil {
		return
	}

	if len(self.BZWLogChan) >= 10000 {
		return
	}

	if node.Gold == 0 && node.Id != 0 {
		return
	}

	self.BZWLogChan <- node
}

//! 数据上报
func (self *Server) RunSqlBZWLog() {
	self.BZWLogChan = make(chan *SQL_BZWLog, 10000)
	self.Wait.Add(1)
	for msg := range self.BZWLogChan {
		if msg.Id == 0 {
			break
		}
		for i := 0; i < 5; i++ {
			if lib.InsertTable("log_bzw", msg, 1, GetServer().DB) > 0 {
				break
			}
		}
	}
	self.Wait.Done()
	close(self.BZWLogChan)
	self.BZWLogChan = nil
}

//! 超端日志
func (self *Server) SqlSuperClientLog(node *SQL_SuperClientLog) {
	if self.SuperClientLogChan == nil {
		return
	}

	if len(self.SuperClientLogChan) >= 10000 {
		return
	}

	self.SuperClientLogChan <- node
}

//! 数据上报
func (self *Server) RunSqlSuperClientLog() {
	self.SuperClientLogChan = make(chan *SQL_SuperClientLog, 10000)
	self.Wait.Add(1)
	for msg := range self.SuperClientLogChan {
		if msg.Id == 0 {
			break
		}
		for i := 0; i < 5; i++ {
			if lib.InsertTable("log_client", msg, 1, GetServer().DB) > 0 {
				break
			}
		}
	}
	self.Wait.Done()
	close(self.SuperClientLogChan)
	self.SuperClientLogChan = nil
}
