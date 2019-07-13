package centerserver

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"lib"
	"log"
	"rjmgr"
	"staticfunc"
	"strings"
	"sync"
	"time"
)

//! sql队列
type SQL_Queue struct {
	Sql   string
	Value []byte
	V     bool
}

//! 举报队列
type SQL_Report struct {
	Id   int64
	Uid  int64
	Rid  int64
	Type int
	Dec  string
	Time int64
}

type Config struct {
	Host       string   `json:"host"`   //! 服务器ip
	InHost     string   `json:"inhost"` //! 服务器ip
	Login      string   `json:"login"`  //! 登陆服ip
	Redis      string   `json:"redis"`
	RedisDB    int      `json:"redisdb"`
	DB         string   `json:"db"`         //! 数据库
	White      []string `json:"white"`      //! 白名单
	PrintLevel int      `json:"printlevel"` //! 打印日志等级
	FileLevel  int      `json:"filelevel"`  //! 日志输出文件等级
	MoneyMode  int      `json:"moneymode"`  //! 金币模式 0,1:100;1,1:1;2,1:10000
	Flag       string   `json:"flag"`       //! 唯一标识
	Sign       string   `json:"sign"`       //! 签名
}

type Server struct {
	Con           *Config         //! 配置
	Wait          *sync.WaitGroup //! 同步阻塞
	Redis         *redis.Pool
	DB            *lib.DBServer
	ShutDown      bool                 //! 是否正在执行关闭
	Notice        [2]string            //! 公告内容
	AreaNotice    [2]map[string]string //! 地区公告
	AreaInfo      [2]map[string]string //! 地区信息
	RpcLogin      *lib.ClientPool
	SqlChan       chan *SQL_Queue
	SqlChanReport chan *SQL_Report
}

var serverSingleton *Server = nil

//! 得到服务器指针
func GetServer() *Server {
	if serverSingleton == nil {
		serverSingleton = new(Server)
		serverSingleton.ShutDown = false
		serverSingleton.Con = new(Config)
		serverSingleton.Wait = new(sync.WaitGroup)
		serverSingleton.AreaNotice[0] = make(map[string]string)
		serverSingleton.AreaNotice[1] = make(map[string]string)
		serverSingleton.AreaInfo[0] = make(map[string]string)
		serverSingleton.AreaInfo[1] = make(map[string]string)
		serverSingleton.DB = new(lib.DBServer)
	}

	return serverSingleton
}

//! 初始化
func (self *Server) Init() {
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

	c := GetServer().Redis.Get()
	self.Notice[0], _ = redis.String(c.Do("GET", "notice"))
	c.Close()

	c = GetServer().Redis.Get()
	self.Notice[1], _ = redis.String(c.Do("GET", "notice_mj"))
	c.Close()

	//! 连接数据库
	db := "root:"
	db += rjmgr.GetRJMgr().SQL
	db += "@tcp(127.0.0.1:3308)/"
	db += self.Con.DB
	db += "?charset=utf8&timeout=10s"
	self.DB.Init(db)

	GetNoticeMgr().GetData()
	GetClubMgr().GetData()
	GetSignMgr().GetData()

	rjmgr.GetRJMgr().InitIP(self.Con.Host)

	//! 告诉服务器开启
	go self.ConnectLogin()
}

func (self *Server) InitConfig() {
	//! 打开配置文件
	config, err := ioutil.ReadFile("./center_config.json")
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

	//self.Con.Login = fmt.Sprintf("http://%s/servermsg", self.Con.Login)
}

//! 关闭服务器
func (self *Server) Close() {
	//! 告诉服务器关闭
	var msg staticfunc.Msg_CenterServer
	msg.InIp = ""
	msg.ExIp = ""
	GetServer().CallLogin("ServerMethod.ServerMsg", "centerserver", &msg)

	self.ShutDown = true

	self.SqlQueue("", []byte(""), false)
	self.SqlReport(0, 0, 0, "")

	GetClubMgr().Save()

	GetPersonMgr().SaveAll()

	self.Wait.Wait()

	log.Fatalln("server shutdown")
}

func (self *Server) ConnectLogin() {
	//! 每1秒来一次心跳包
	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		var msg staticfunc.Msg_CenterServer
		msg.ExIp = self.Con.Host
		msg.InIp = self.Con.InHost
		result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "centerserver", &msg)
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

//! 得到公告
func (self *Server) GetAreaNotice(code string, group int) string {
	notice, ok := self.AreaNotice[group-1][code]
	if ok {
		return notice
	}

	c := self.Redis.Get()
	defer c.Close()
	if group == 1 {
		v, err := redis.String(c.Do("GET", "area_"+code))
		if err == nil {
			GetServer().AreaNotice[group-1][code] = v
			return v
		}
	} else {
		v, err := redis.String(c.Do("GET", "areamj_"+code))
		if err == nil {
			GetServer().AreaNotice[group-1][code] = v
			return v
		}
	}

	return ""
}

//! 得到联系方式
func (self *Server) GetAreaInfo(code string, group int) string {
	notice, ok := self.AreaInfo[group-1][code]
	if ok {
		return notice
	}

	c := self.Redis.Get()
	defer c.Close()
	if group == 1 {
		v, err := redis.String(c.Do("GET", "areainfo_"+code))
		if err == nil {
			GetServer().AreaInfo[group-1][code] = v
			return v
		}
	} else {
		v, err := redis.String(c.Do("GET", "areainfomj_"+code))
		if err == nil {
			GetServer().AreaInfo[group-1][code] = v
			return v
		}
	}

	return ""
}

//! 是否是白名单
func (self *Server) IsWhite(ip string) bool {
	lib.GetLogMgr().Output(lib.LOG_INFO, ip)

	if ip == "127.0.0.1" {
		return true
	}

	for i := 0; i < len(self.Con.White); i++ {
		if self.Con.White[i] == "0.0.0.0" || self.Con.White[i] == ip {
			return true
		}
	}

	return false
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

//! 加卡
func (self *Server) AddCard(uid int64, num int, ip string, dec int) (bool, int, int) {
	var msg staticfunc.Msg_GiveCard
	msg.Uid = uid
	msg.Pid = staticfunc.TYPE_CARD
	msg.Num = num
	msg.Ip = ip
	msg.Dec = dec
	result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "givecard", &msg)

	if err != nil || string(result) == "false" {
		return false, 0, 0
	}

	tmp := strings.Split(string(result), "_")
	if len(tmp) != 2 {
		return false, 0, 0
	}

	now1 := lib.HF_Atoi(tmp[0])
	now2 := lib.HF_Atoi(tmp[1])

	return true, now1, now2
}

//! 加金币
func (self *Server) AddGold(uid int64, num int, ip string, dec int) (bool, int, int) {
	var msg staticfunc.Msg_GiveCard
	msg.Uid = uid
	msg.Pid = staticfunc.TYPE_GOLD
	msg.Num = num
	msg.Ip = ip
	msg.Dec = dec
	result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "givecard", &msg)

	if err != nil || string(result) == "false" {
		return false, 0, 0
	}

	tmp := strings.Split(string(result), "_")
	if len(tmp) != 2 {
		return false, 0, 0
	}

	now1 := lib.HF_Atoi(tmp[0])
	now2 := lib.HF_Atoi(tmp[1])

	return true, now1, now2
}

//! 减卡
func (self *Server) CostCard(uid int64, num int, ip string, dec int) (bool, int, int) {
	var msg staticfunc.Msg_GiveCard
	msg.Uid = uid
	msg.Pid = staticfunc.TYPE_CARD
	msg.Num = num
	msg.Ip = ip
	msg.Dec = dec
	result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "movecard", &msg)

	if err != nil || string(result) == "false" {
		return false, 0, 0
	}

	tmp := strings.Split(string(result), "_")
	if len(tmp) != 2 {
		return false, 0, 0
	}

	now1 := lib.HF_Atoi(tmp[0])
	now2 := lib.HF_Atoi(tmp[1])

	return true, now1, now2
}

//! 得到卡
func (self *Server) GetCard(uid int64) (int, int, int) {
	var msg staticfunc.Msg_Uid
	msg.Uid = uid
	result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "getcard", &msg)

	if err != nil || string(result) == "false" {
		return 0, 0, 0
	}

	tmp := strings.Split(string(result), "_")
	if len(tmp) != 3 {
		return 0, 0, 0
	}

	now1 := lib.HF_Atoi(tmp[0])
	now2 := lib.HF_Atoi(tmp[1])
	now3 := lib.HF_Atoi(tmp[2])

	return now1, now2, now3
}

//! 加金币
func (self *Server) CostGold(uid int64, num int, ip string, dec int) (bool, int, int) {
	var msg staticfunc.Msg_GiveCard
	msg.Uid = uid
	msg.Pid = staticfunc.TYPE_GOLD
	msg.Num = num
	msg.Ip = ip
	msg.Dec = dec
	result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "movecard", &msg)

	if err != nil || string(result) == "false" {
		return false, 0, 0
	}

	tmp := strings.Split(string(result), "_")
	if len(tmp) != 2 {
		return false, 0, 0
	}

	now1 := lib.HF_Atoi(tmp[0])
	now2 := lib.HF_Atoi(tmp[1])

	return true, now1, now2
}

//! 数据库队列
func (self *Server) SqlQueue(sql string, value []byte, v bool) {
	if self.SqlChan == nil {
		return
	}

	self.SqlChan <- &SQL_Queue{sql, value, v}
}

//! 数据上报
func (self *Server) RunSqlQueue() {
	self.SqlChan = make(chan *SQL_Queue, 50000)
	self.Wait.Add(1)
	for msg := range self.SqlChan {
		if msg.Sql == "" {
			break
		}
		for i := 0; i < 5; i++ {
			if msg.V {
				_, err := GetServer().DB.DB.Exec(msg.Sql, msg.Value)
				if err == nil {
					break
				}
			} else {
				_, err := GetServer().DB.DB.Exec(msg.Sql)
				if err == nil {
					break
				}
			}
		}
	}
	self.Wait.Done()
	close(self.SqlChan)
	self.SqlChan = nil
}

//! 举报队列
func (self *Server) SqlReport(uid int64, rid int64, _type int, dec string) {
	return             //! 暂时关闭举报功能
	if rid == 100000 { //! 系统不被举报
		return
	}

	if self.SqlChanReport == nil {
		return
	}

	if len(self.SqlChanReport) >= 10000 {
		return
	}

	self.SqlChanReport <- &SQL_Report{0, uid, rid, _type, dec, time.Now().Unix()}
}

//! 数据上报
func (self *Server) RunSqlReport() {
	return //! 暂时关闭举报功能
	self.SqlChanReport = make(chan *SQL_Report, 10000)
	self.Wait.Add(1)
	for msg := range self.SqlChanReport {
		if msg.Uid == 0 {
			break
		}
		for i := 0; i < 5; i++ {
			if lib.InsertTable("report", msg, 1, GetServer().DB) > 0 {
				break
			}
		}
	}
	self.Wait.Done()
	close(self.SqlChanReport)
	self.SqlChanReport = nil
}
