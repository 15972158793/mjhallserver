package loginserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"lib"
	"log"
	"rjmgr"
	"staticfunc"
	"strings"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

func (s lstGoldRoom) Len() int           { return len(s) }
func (s lstGoldRoom) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s lstGoldRoom) Less(i, j int) bool { return s[i].Time < s[j].Time }

type JS_GoldRoom struct {
	RoomId int
	Time   int64
}

type lstGoldRoom []JS_GoldRoom

//! sql队列
type SQL_Queue struct {
	Sql   string
	Value []byte
	V     bool
}

//! 金币日志
type SQL_Gold struct {
	Id   int64
	Uid  int64
	Num  int
	Dec  int
	Time int64
}

//! 兑换日志
type SQL_Exchange struct {
	Id    int64
	Uid   int64
	Name  string
	Gold  int
	Time  int64
	State int
	Dec   string
}

type DB_Strc struct {
	Id       int64
	Value    []byte
	Time     int64
	Gold     int
	SaveGold int

	lib.DataUpdate
}

type DB_OpenId struct {
	Id     int64
	Openid string //! unionid
	Wyid   string //! openid

	lib.DataUpdate
}

type DB_OldUid struct {
	Openid string
	Uid    string
	Bit    int
}

type RoomInfo struct {
	GameType int
	Num      int
	Begin    bool
	Time     int64
	Player   map[int64]string
}

func (self *RoomInfo) IsHasIP(ip string) bool {
	if ip == "" {
		return false
	}

	for _, value := range self.Player {
		if value == ip {
			return true
		}
	}
	return false
}

type GameServerConfig struct {
	Id        int
	InIp      string //! 内网ip
	ExIp      string //! 外网ip(可能是高防域名)
	Type      int    //! 0任意房间 1扑克 2麻将 3匹配场
	Room      map[int]*RoomInfo
	RoomLock  *sync.RWMutex
	RpcClient *lib.ClientPool
}

func NewGameServerConfig() *GameServerConfig {
	p := new(GameServerConfig)
	p.RoomLock = new(sync.RWMutex)
	return p
}

func (self *GameServerConfig) AddRoom(room int, gametype int, num int, uid int64, ip string) int {
	self.RoomLock.Lock()
	defer self.RoomLock.Unlock()

	lib.GetLogMgr().Output(lib.LOG_INFO, "room加入一个人", room)

	_, ok := self.Room[room]
	if ok {
		self.Room[room].Num = num
		self.Room[room].Player[uid] = ip
		return self.Room[room].GameType
	} else {
		self.Room[room] = &RoomInfo{gametype, num, false, time.Now().Unix(), make(map[int64]string)}
		self.Room[room].Player[uid] = ip
		return gametype
	}
}

func (self *GameServerConfig) RemoveRoom(room int, uid int64) {
	self.RoomLock.Lock()
	defer self.RoomLock.Unlock()

	lib.GetLogMgr().Output(lib.LOG_INFO, "room减少一个人", room)

	_, ok := self.Room[room]
	if ok {
		self.Room[room].Num--
		delete(self.Room[room].Player, uid)
	}
}

func (self *GameServerConfig) DelRoom(room int) {
	self.RoomLock.Lock()
	defer self.RoomLock.Unlock()

	lib.GetLogMgr().Output(lib.LOG_INFO, "删除房间", room)

	delete(self.Room, room)
}

func (self *GameServerConfig) HasRoom(room int) bool {
	self.RoomLock.RLock()
	defer self.RoomLock.RUnlock()

	_, ok := self.Room[room]
	return ok
}

func (self *GameServerConfig) SetBegin(room int) {
	self.RoomLock.Lock()
	defer self.RoomLock.Unlock()

	_, ok := self.Room[room]
	if ok {
		self.Room[room].Begin = true
	}
}

//! 得到一个人数
func (self *GameServerConfig) GetRoomId(gametype int, noroomid int, ip string) (bool, int) {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][gametype]
	if !ok {
		lib.GetLogMgr().Output(lib.LOG_INFO, "获取gametype失败")
		return false, 0
	}

	self.RoomLock.RLock()
	defer self.RoomLock.RUnlock()

	var lst lstGoldRoom
	for key, value := range self.Room {
		lib.GetLogMgr().Output(lib.LOG_INFO, "匹配房间列表", self.Room, value.Begin, value.Num)
		if value.GameType != gametype {
			continue
		}
		if key == noroomid {
			continue
		}
		if value.Num >= lib.HF_Atoi(csv["maxnum"]) {
			continue
		}
		if (gametype/10000 == 1 || gametype/10000 == 2) && GetServer().Con.CheckIP == 1 && value.IsHasIP(ip) { //! 卡五星和拼三张
			continue
		}
		lst = append(lst, JS_GoldRoom{key, value.Time})
		if len(lst) >= 10 {
			break
		}
	}
	if len(lst) == 0 {
		return true, 0
	}

	return true, lst[lib.HF_GetRandom(len(lst))].RoomId
}

//! 关闭rpc
func (self *GameServerConfig) Close() {
	if self.RpcClient == nil {
		return
	}
	self.RpcClient.CloseAll()
}

//! 发送消息
func (self *GameServerConfig) Call(method string, msghead string, v interface{}) ([]byte, error) {
	data := lib.HF_EncodeMsg(msghead, v, false)

	err, client := self.RpcClient.RandomGetConn()
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

type Config struct {
	Host        string   `json:"host"`
	Redis       string   `json:"redis"`
	RcRedis     string   `json:"rcredis"` //! 记录redis
	DB          string   `json:"db"`
	NewCard     int      `json:"newcard"`
	NewGold     int      `json:"newgold"`
	Version     int      `json:"version"`
	MinVersion  int      `json:"minversion"`
	RedisDB     int      `json:"redisdb"`
	RcRedisDB   int      `json:"rcredisdb"`
	GateServer  []string `json:"gateserver"`
	GoldRoomNum int      `json:"goldroomnum"` //! 一个服务器可容纳多少金币场房间
	PrintLevel  int      `json:"printlevel"`  //! 打印日志等级
	FileLevel   int      `json:"filelevel"`   //! 日志输出文件等级
	CheckIP     int      `json:"checkip"`
	GameOver    int      `json:"gameover"`
	White       []string `json:"white"`
	AG          int      `json:"ag"`      //! 是否开启真人视讯
	AGWhite     []int    `json:"agwhite"` //! 百名单
	Flag        string   `json:"flag"`    //! 唯一标识
	Sign        string   `json:"sign"`    //! 签名
}

type Server struct {
	Con             *Config         //! 配置
	MoneyMode       int             //! 金币模式 0,1:100;1,1:1;2,1:10000
	Wait            *sync.WaitGroup //! 同步阻塞
	ShutDown        bool            //! 是否正在执行关闭
	Redis           *redis.Pool
	RcRedis         *redis.Pool
	DB              *lib.DBServer
	MapGameServer   map[int]*GameServerConfig
	InCenterServer  string          //! 中心服务器内网
	ExCenterServer  string          //! 中心服务器外网
	RpcCenterServer *lib.ClientPool //! 中心服务器rpc
	GameLock        *sync.RWMutex
	Black           map[int64]int //! 黑名单
	BlackLock       *sync.RWMutex
	MapCode         map[string]int
	CodeLock        *sync.RWMutex
	SqlChan         chan *SQL_Queue
	SqlLogChan      chan *SQL_Queue
	SqlGoldChan     chan *SQL_Gold
	SqlExchangeChan chan *SQL_Exchange
	ManyRoodId      map[int]int //! 百人场的roomid
	InitMoney       int         //! 初始金币数量
	WZQMode         staticfunc.WZQMode
	GameMode        staticfunc.GameMode
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
		serverSingleton.MapGameServer = make(map[int]*GameServerConfig)
		serverSingleton.GameLock = new(sync.RWMutex)
		serverSingleton.Black = make(map[int64]int)
		serverSingleton.BlackLock = new(sync.RWMutex)
		serverSingleton.MapCode = make(map[string]int)
		serverSingleton.CodeLock = new(sync.RWMutex)
		serverSingleton.ManyRoodId = make(map[int]int)
		serverSingleton.InitMoney = 0
	}

	return serverSingleton
}

//! 初始化
func (self *Server) Init() {
	if self.Con.AG == 1 {
		if !lib.GetAGVideoMgr().Init("https://api.a45.me/api/wallet/Gateway.php", "zdswapi", "4de29a7d43a036e7b7e5451d7216b188") {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "ag视讯初始化失败，可能导致ag视讯无法正常工作")
		}
	} else {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "ag视讯关闭")
	}

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
		log.Fatal("rcredis err")
		return
	}

	self.MoneyMode = rjmgr.GetRJMgr().MoneyMode

	//! 连接数据库
	db := "root:"
	db += rjmgr.GetRJMgr().SQL
	db += "@tcp(127.0.0.1:3308)/"
	db += self.Con.DB
	db += "?charset=utf8&timeout=10s"
	self.DB.Init(db)

	//! 载入黑名单
	self.GetBlack()

	rjmgr.GetRJMgr().InitIP(self.Con.Host)

	c := GetServer().Redis.Get()
	defer c.Close()
	v, err := redis.Int64(c.Do("GET", "initmoney"))
	if err == nil {
		self.InitMoney = int(v)
	} else {
		self.InitMoney = 0
	}

	value, err := redis.Bytes(c.Do("GET", "wzqmode"))
	if err == nil {
		json.Unmarshal(value, &self.WZQMode)
	}

	//! 初始化允许打开哪些游戏
	InitOpenGame()
}

func (self *Server) Run() {
	ticker := time.NewTicker(time.Second * 60)
	for {
		<-ticker.C
		self.ClearCode()
		for i := 0; i < 3; i++ {
			lib.GetFishMgr().Run(i)
		}
		if self.ShutDown {
			break
		}
	}
}

func (self *Server) InitConfig() {
	//! 打开配置文件
	config, err := ioutil.ReadFile("./login_config.json")
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

	lib.GetLogMgr().Output(lib.LOG_DEBUG, "当前版本号:", self.Con.Version)
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "最小版本号:", self.Con.MinVersion)

	if self.Con.GoldRoomNum == 0 {
		self.Con.GoldRoomNum = 2000
	}
}

//! 关闭服务器
func (self *Server) Close() {
	self.ShutDown = true

	self.SqlQueue("", []byte(""), false)
	self.SqlLog("", []byte(""), false)
	self.SqlGoldLog(0, 0, 0)
	self.SqlExchangeLog(0, "", 0)

	//! 告诉各个login服务器
	var msg staticfunc.Msg_Null
	for _, value := range self.MapGameServer {
		value.Call("ServerMethod.ServerMsg", "loginserver", &msg)
	}
	self.CallCenter("ServerMethod.ServerMsg", "loginserver", &msg)

	//! 等所有存入数据库的指令执行完了
	self.Wait.Wait()

	self.DB.Close()
	self.Redis.Close()

	log.Fatalln("server shutdown")
}

//! 发送moneymode
func (self *Server) SendMoneyMode() {
	lst := self.GetGameServerFromType(3)
	var msg staticfunc.Msg_MoneyMode
	msg.MoneyMode = self.MoneyMode
	for i := 0; i < len(lst); i++ {
		lst[i].Call("ServerMethod.ServerMsg", "moneymode", &msg)
	}
	self.CallCenter("ServerMethod.ServerMsg", "moneymode", &msg)
}

//! 通过openid得到uid
var UidLock *sync.RWMutex = new(sync.RWMutex)

func (self *Server) DB_GetUid(openid string, wyid string, insert bool, olduid int64) int64 {
	UidLock.Lock()
	defer UidLock.Unlock()

	//! 先从redis里面得到
	c := self.Redis.Get()
	defer c.Close()

	v, err := redis.Int64(c.Do("GET", "account_"+openid))
	if err == nil {
		c.Do("EXPIRE", "account_"+openid, 86400*7)
		return v
	}

	var dbstr DB_OpenId
	sql := fmt.Sprintf("select * from `account` where `openid` = '%s'", openid)
	GetServer().DB.GetOneData(sql, &dbstr)
	if dbstr.Id <= 0 {
		if insert {
			dbstr.Id = olduid
			dbstr.Openid = openid
			dbstr.Wyid = wyid
			uid := olduid
			if olduid == 0 {
				uid = lib.InsertTable("account", &dbstr, 1, self.DB)
			} else {
				lib.InsertTable("account", &dbstr, 0, self.DB)
			}
			c.Do("SET", "account_"+openid, uid)
			c.Do("EXPIRE", "account_"+openid, 86400*7)
			return uid
		} else {
			return 0
		}
	} else {
		c.Do("SET", "account_"+openid, dbstr.Id)
		c.Do("EXPIRE", "account_"+openid, 86400*7)
		return dbstr.Id
	}
}

//! 判断账号和密码是否一致
func (self *Server) CheckAccountAndPasswd(account string, passwd string) int64 {
	c := self.Redis.Get()
	defer c.Close()

	v, err := redis.String(c.Do("GET", "accountplus_"+account))
	if err == nil {
		if v != passwd {
			return -1 //! 账号密码不一致
		} else {
			c.Do("EXPIRE", "accountplus_"+account, 86400*7)
			return self.DB_GetUid(account, "", false, 0)
		}
	} else {
		var dbstr DB_OpenId
		sql := fmt.Sprintf("select * from `account` where `openid` = '%s'", account)
		GetServer().DB.GetOneData(sql, &dbstr)
		if dbstr.Id <= 0 {
			return 0
		} else {
			if dbstr.Wyid != passwd {
				return -1
			}
			c.Do("SET", "accountplus_"+account, passwd)
			c.Do("EXPIRE", "accountplus_"+account, 86400*7)

			c.Do("SET", "account_"+account, dbstr.Id)
			c.Do("EXPIRE", "account_"+account, 86400*7)
			return dbstr.Id
		}
	}
}

func (self *Server) DB_ModifyPasswd(account string, newpasswd string) {
	//! 先将passwd做md5
	//passwd = lib.HF_MD5(passwd)

	//! 先从redis里面得到
	c := self.Redis.Get()
	defer c.Close()

	c.Do("SET", "accountplus_"+account, newpasswd)
	c.Do("EXPIRE", "accountplus_"+account, 86400*7)

	sql := fmt.Sprintf("update `account` set `wyid` = '%s' where `openid` = '%s'", newpasswd, account)
	self.SqlQueue(sql, []byte(""), false)
}

//! 获取一个数据
//! db 0仅从redis取 1从mysql里取但不插入 2mysql没有新建
func (self *Server) DB_GetData(table string, uid int64, db int) ([]byte, bool) {
	//! 先从redis里面得到该数据
	c := self.Redis.Get()
	defer c.Close()
	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("%s_%d", table, uid)))
	if err == nil { //! redis找得到，则直接返回redis里的数据
		c.Do("EXPIRE", fmt.Sprintf("%s_%d", table, uid), 86400*7)
		return v, false
	}

	if db == 0 {
		return []byte(""), false
	}

	var dbstr DB_Strc
	sql := fmt.Sprintf("select * from `%s` where `id` = '%d'", table, uid)
	GetServer().DB.GetOneData(sql, &dbstr)
	if dbstr.Id <= 0 {
		if db == 1 {
			return []byte(""), false
		}
		dbstr.Id = uid
		dbstr.Value = []byte("")
		sql := fmt.Sprintf("insert into `%s`(`id`, `value`, `createtime`, `gold`, `savegold`) values (%d, ?, %d, %d, %d)", table, uid, time.Now().Unix(), 0, 0)
		self.SqlQueue(sql, dbstr.Value, true)
	}

	c.Do("SET", fmt.Sprintf("%s_%d", table, uid), dbstr.Value)
	c.Do("EXPIRE", fmt.Sprintf("%s_%d", table, uid), 86400*7)

	return dbstr.Value, true
}

//! 改变数据
func (self *Server) DB_SetData(table string, uid int64, value []byte, gold int, savegold int, db bool) {
	//! 保存到redis
	c := self.Redis.Get()
	defer c.Close()
	key := fmt.Sprintf("%s_%d", table, uid)
	c.Do("SET", key, value)
	c.Do("EXPIRE", key, 86400*7)

	//! 保存数据库
	if db {
		sql := fmt.Sprintf("update `%s` set `value` = ?, `gold` = %d, `savegold` = %d where `id` = %d", table, gold, savegold, uid)
		self.SqlQueue(sql, value, true)
	}
}

//! 插入数据
func (self *Server) InsertLog(uid int64, _type int, num int, info string) {
	if info == "0.0.0.0" {
		return
	}
	sql := fmt.Sprintf("insert into `log_base`(`uid`, `type`, `num`, `info`, `creation_time`) values(%d, %d, %d, '%s', %d)", uid, _type, num, info, time.Now().Unix())
	self.SqlLog(sql, []byte(""), false)
}

//! 加一个游戏服务器
func (self *Server) AddGameServer(id int, inip string, exip string, _type int) {
	self.GameLock.Lock()
	defer self.GameLock.Unlock()

	_, ok := self.MapGameServer[id]
	if ok {
		return
	}

	config := NewGameServerConfig()
	config.Id = id
	config.InIp = inip
	config.ExIp = exip
	config.Type = _type
	config.Room = make(map[int]*RoomInfo)
	ip := strings.Split(inip, ":")
	_ip := ip[0] + ":1" + ip[1]
	config.RpcClient = lib.CreateClientPool([]string{_ip}, _ip)
	self.MapGameServer[id] = config

	lib.GetLogMgr().Output(lib.LOG_DEBUG, id, "连接:", inip)
}

//! 删除一个游戏服务器
func (self *Server) DelGameServer(id int) {
	self.GameLock.Lock()
	defer self.GameLock.Unlock()

	value, ok := self.MapGameServer[id]
	if ok {
		value.Close()
	}
	delete(self.MapGameServer, id)

	lib.GetLogMgr().Output(lib.LOG_DEBUG, id, "断开")
}

//! 得到游戏服务器
func (self *Server) GetGameServer(id int) *GameServerConfig {
	self.GameLock.RLock()
	defer self.GameLock.RUnlock()

	config, ok := self.MapGameServer[id]
	if ok {
		return config
	}

	return nil
}

//! 得到room最少的gameserver
func (self *Server) GetGameServerFromRoom(_type int) *GameServerConfig {
	self.GameLock.RLock()
	defer self.GameLock.RUnlock()

	id := 0
	num := 1000000000
	for key, value := range self.MapGameServer {
		if value.Type != 0 && value.Type != _type {
			continue
		}
		if len(value.Room) < num {
			num = len(value.Room)
			id = key
		}
	}

	if id == 0 {
		return nil
	}

	return self.MapGameServer[id]
}

//! 得到所有的gameserver
func (self *Server) GetGameServerFromType(_type int) []*GameServerConfig {
	self.GameLock.RLock()
	defer self.GameLock.RUnlock()

	lst := make([]*GameServerConfig, 0)
	for _, value := range self.MapGameServer {
		if value.Type != 0 && value.Type != _type {
			continue
		}
		lst = append(lst, value)
	}

	return lst
}

//! 得到匹配房间
func (self *Server) GetGameServerFromGoldRoom(_type int) *GameServerConfig {
	if _type < 3 {
		_type = 3
	}
	for _, value := range self.MapGameServer {
		if value.Type != 0 && value.Type != _type {
			continue
		}
		if len(value.Room) >= self.Con.GoldRoomNum {
			continue
		}
		return value
	}

	//return self.GetGameServerFromRoom(3)
	return nil
}

//! 加一个房间
func (self *Server) AddGameRoom(id int, room int, gametype int, num int, uid int64, ip string) (int, *GameServerConfig) {
	self.GameLock.RLock()
	defer self.GameLock.RUnlock()

	config, ok := self.MapGameServer[id]
	if ok {
		return config.AddRoom(room, gametype, num, uid, ip), config
	}

	return 0, config
}

//! 减少一个人
func (self *Server) RemoveGameRoom(id int, room int, uid int64) {
	self.GameLock.RLock()
	defer self.GameLock.RUnlock()

	config, ok := self.MapGameServer[id]
	if ok {
		config.RemoveRoom(room, uid)
	}
}

//! 删一个房间
func (self *Server) DelGameRoom(id int, room int) {
	self.GameLock.RLock()
	defer self.GameLock.RUnlock()

	config, ok := self.MapGameServer[id]
	if ok {
		config.DelRoom(room)
	}
}

//! 开始一个房间
func (self *Server) BeginGameRoom(id int, room int) {
	self.GameLock.RLock()
	defer self.GameLock.RUnlock()

	config, ok := self.MapGameServer[id]
	if ok {
		config.SetBegin(room)
	}
}

//! 得到房间
func (self *Server) GetGameRoom(room int, group int) *GameServerConfig {
	self.GameLock.RLock()
	defer self.GameLock.RUnlock()

	for _, value := range self.MapGameServer {
		if group != -1 && value.Type-1 != group {
			continue
		}
		if value.HasRoom(room) {
			return value
		}
	}

	return nil
}

//////
func (self *Server) CreateRoom(id int, roomid int, _type int, num int, param1 int, param2 int, agent int64, clubid int64) bool {
	config, ok := self.MapGameServer[id]
	if ok {
		var msg staticfunc.Msg_CreateRoom
		msg.Id = roomid
		msg.Type = _type
		msg.Num = num
		msg.Param1 = param1
		msg.Param2 = param2
		msg.Agent = agent
		msg.ClubId = clubid

		result, err := config.Call("ServerMethod.ServerMsg", "createroom", &msg)
		if err != nil || string(result) == "" {
			return false
		}

		if string(result) == "true" {
			return true
		} else {
			return false
		}
	}

	return false
}

func (self *Server) JoinRoom(id int, roomid int) int {
	config, ok := self.MapGameServer[id]
	if ok {
		var msg staticfunc.Msg_JoinRoom
		msg.Id = roomid

		result, err := config.Call("ServerMethod.ServerMsg", "joinroom", &msg)
		if err != nil || string(result) == "" {
			return -1
		}
		return lib.HF_Atoi(string(result))
	}

	return -1
}

//! 激活房间
func (self *Server) ActiveRoom(roomid int) (int, int) {
	for _, config := range self.MapGameServer {
		if config.Type >= 3 { //! 匹配场不能激活俱乐部房间
			continue
		}
		var msg staticfunc.Msg_JoinRoom
		msg.Id = roomid

		result, err := config.Call("ServerMethod.ServerMsg", "joinroom", &msg)
		if err != nil || string(result) == "" {
			continue
		}

		num := lib.HF_Atoi(string(result))
		if num == -1 {
			continue
		}

		return config.Id, num
	}

	return 0, -1
}

///////////
type DB_Black struct {
	Id         int64
	Createtime int64
}

//! 黑名单
func (self *Server) GetBlack() {
	var black DB_Black
	res := GetServer().DB.GetAllData("select * from `black`", &black)
	for i := 0; i < len(res); i++ {
		data := res[i].(*DB_Black)
		self.Black[data.Id] = 1
	}
}

//! 加入黑名单
func (self *Server) AddBlack(uid int64) {
	ok := self.IsBlack(uid)
	if ok {
		return
	}

	self.BlackLock.Lock()
	self.Black[uid] = 1
	self.BlackLock.Unlock()

	go func() {
		self.Wait.Add(1)
		defer self.Wait.Done()

		var black DB_Black
		black.Id = uid
		black.Createtime = time.Now().Unix()
		lib.InsertTable("black", &black, 0, self.DB)
	}()

	value, _ := GetServer().DB_GetData("user", uid, 1)
	if string(value) == "" {
		return
	}

	var person Person
	json.Unmarshal(value, &person)
	if person.GameId == 0 {
		return
	}

	config := GetServer().GetGameServer(person.GameId)
	if config == nil {
		return
	}

	var msg staticfunc.Msg_Uid
	msg.Uid = person.Uid
	config.Call("ServerMethod.ServerMsg", "addblack", &msg)
}

//! 删除黑名单
func (self *Server) DelBlack(uid int64) {
	ok := self.IsBlack(uid)
	if !ok {
		return
	}

	self.BlackLock.Lock()
	delete(self.Black, uid)
	self.BlackLock.Unlock()

	go func() {
		self.Wait.Add(1)
		defer self.Wait.Done()

		sql := fmt.Sprintf("delete from `black` where `id` = %d", uid)
		GetServer().DB.DB.Exec(sql)
	}()
}

//! 是否在黑名单
func (self *Server) IsBlack(uid int64) bool {
	self.BlackLock.RLock()
	defer self.BlackLock.RUnlock()

	_, ok := self.Black[uid]
	return ok
}

func (self *Server) AddCode(code string) {
	self.CodeLock.Lock()
	defer self.CodeLock.Unlock()

	self.MapCode[code] = 1
}

func (self *Server) ClearCode() {
	self.CodeLock.Lock()
	defer self.CodeLock.Unlock()

	self.MapCode = make(map[string]int)
}

func (self *Server) HasCode(code string) bool {
	self.CodeLock.RLock()
	defer self.CodeLock.RUnlock()

	_, ok := self.MapCode[code]

	return ok
}

//! 是否五子棋白名单
func (self *Server) IsWZQWhite(uid int64) bool {
	for i := 0; i < len(self.WZQMode.White); i++ {
		if self.WZQMode.White[i] == uid {
			return true
		}
	}

	return false
}

//! 是否是白名单
func (self *Server) IsWhite(ip string, operate string) bool {
	lib.GetLogMgr().Output(lib.LOG_INFO, ip, ":", operate)

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

//! 得到网关ip
func (self *Server) GetGateIP(ip string) string {
	//if self.Con.Mode == 0 {
	//	return ip
	//}
	if len(self.Con.GateServer) == 0 {
		return ""
	}

	return self.Con.GateServer[lib.HF_GetRandom(len(self.Con.GateServer))]
}

//! 调用中心服务器
func (self *Server) CallCenter(method string, msghead string, v interface{}) ([]byte, error) {
	if self.RpcCenterServer == nil {
		return []byte(""), errors.New("RpcCenterServer is nil")
	}

	data := lib.HF_EncodeMsg(msghead, v, false)

	err, client := self.RpcCenterServer.RandomGetConn()
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

//! 数据库队列
func (self *Server) SqlLog(sql string, value []byte, v bool) {
	if self.SqlLogChan == nil {
		return
	}

	if len(self.SqlLogChan) >= 10000 {
		return
	}

	self.SqlLogChan <- &SQL_Queue{sql, value, v}
}

//! 数据上报
func (self *Server) RunSqlLog() {
	self.SqlLogChan = make(chan *SQL_Queue, 10000)
	self.Wait.Add(1)
	for msg := range self.SqlLogChan {
		if msg.Sql == "" {
			break
		}
		for i := 0; i < 5; i++ {
			_, err := GetServer().DB.DB.Exec(msg.Sql)
			if err == nil {
				break
			}
		}
	}
	self.Wait.Done()
	close(self.SqlLogChan)
	self.SqlLogChan = nil
}

//! 金币日志
func (self *Server) SqlGoldLog(uid int64, num int, dec int) {
	if self.SqlGoldChan == nil {
		return
	}

	if len(self.SqlGoldChan) >= 10000 {
		return
	}

	self.SqlGoldChan <- &SQL_Gold{0, uid, num, dec, time.Now().Unix()}
}

//! 数据上报
func (self *Server) RunSqlGoldLog() {
	self.SqlGoldChan = make(chan *SQL_Gold, 10000)
	self.Wait.Add(1)
	for msg := range self.SqlGoldChan {
		if msg.Uid == 0 {
			break
		}
		for i := 0; i < 5; i++ {
			if lib.InsertTable("log_gold", msg, 1, self.DB) > 0 {
				break
			}
		}
	}
	self.Wait.Done()
	close(self.SqlGoldChan)
	self.SqlGoldChan = nil
}

//! 兑换日志
func (self *Server) SqlExchangeLog(uid int64, name string, gold int) {
	if self.SqlExchangeChan == nil {
		return
	}

	if len(self.SqlExchangeChan) >= 10000 {
		return
	}

	self.SqlExchangeChan <- &SQL_Exchange{0, uid, name, gold, time.Now().Unix(), 0, ""}
}

//! 数据上报
func (self *Server) RunSqlExchangeLog() {
	self.SqlExchangeChan = make(chan *SQL_Exchange, 10000)
	self.Wait.Add(1)
	for msg := range self.SqlExchangeChan {
		if msg.Uid == 0 {
			break
		}
		for i := 0; i < 5; i++ {
			sql := fmt.Sprintf("insert into `exchange`(`uid`, `name`, `gold`, `time`, `state`, `dec`) values (%d, ?, %d, %d, 0, '')", msg.Uid, msg.Gold, msg.Time)
			_, err := GetServer().DB.DB.Exec(sql, []byte(msg.Name))
			if err == nil {
				break
			}
		}
	}
	self.Wait.Done()
	close(self.SqlExchangeChan)
	self.SqlExchangeChan = nil
}

func (self *Server) InsertGiveGoldRecord(giveType int, uid int64, _uid int64, info string) {
	c := self.RcRedis.Get()
	defer c.Close()

	var table string
	if giveType == 1 {
		table = fmt.Sprintf("give_%d", uid)
	} else {
		table = fmt.Sprintf("get_%d", _uid)
	}

	v, err := redis.Int(c.Do("LLEN", table))
	if err != nil {
		v = 0
	}

	for v >= 20 { //! 超过上限，删除
		c.Do("RPOP", table)
		v--
	}

	c.Do("LPUSH", table, info)
	c.Do("EXPIRE", table, 86400*2)
}
