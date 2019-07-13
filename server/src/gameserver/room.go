package gameserver

import (
	"encoding/json"
	"lib"
	"runtime/debug"
	"staticfunc"
	"strings"
	"sync"
	"time"
)

//! 房间消息
//! 加入房间
type Msg_JoinRoom struct {
	Uid     int64  `json:"uid"` //! uid
	Roomid  int    `json:"roomid"`
	MInfo   string `json:"minfo"`
	Param   int    `json:"param"`
	UnionId string `json:"unionid"`
}

//! 玩家上线下线
type Msg_LinePerson struct {
	Uid  int64 `json:"uid"`
	Line bool  `json:"line"`
}

//! 聊天
type Msg_Chat struct {
	Uid  int64  `json:"uid"`
	Type int    `json:"type"`
	Chat string `json:"chat"`
}

/////////////////////////////////////////////////
//! 房间信息
type Msg_RoomInfo struct {
	RoomId  int              `json:"roomid"`
	Host    int64            `json:"host"`
	Name    string           `json:"name"`
	Head    string           `json:"head"`
	Person  []Son_PersonInfo `json:"person"`
	Time    int64            `json:"time"`
	Agree   []int64          `json:"agree"`
	ETime   []int            `json:"etime"`
	Step    int              `json:"step"`
	MaxStep int              `json:"maxstep"`
	Param1  int              `json:"param1"`
	Param2  int              `json:"param2"`
	Type    int              `json:"type"`
	Agent   bool             `json:"agent"`
}

type Son_PersonInfo struct {
	Uid       int64  `json:"uid"`
	Name      string `json:"name"`
	ImgUrl    string `json:"imgurl"`
	Sex       int    `json:"sex"`
	Ip        string `json:"ip"`
	Line      bool   `json:"line"`
	Address   string `json:"address"`
	Longitude string `json:"longitude"` //! 经度
	Latitude  string `json:"latitude"`  //! 纬度
	Param     int    `json:"param"`
}

//! 加入房间失败
type Msg_JoinRoomFail struct {
	Result int `json:"result"`
}

//! 退出房间
type Msg_ExitRoom struct {
	Result int `json:"result"`
}

//! 解散房间
type Msg_DissmissRoom struct {
	Agree []int64 `json:"agree"` //! 同意解散的人
	Time  int64   `json:"time"`  //! 剩余时间
	ETime []int   `json:"etime"`
}

/////////////////////////////////////////////////////////
type RoomMsg struct {
	Head string
	Data string
	Uid  int64
	V    interface{}
}

func NewRoomMsg(head string, data string, uid int64, v interface{}) *RoomMsg {
	roommsg := new(RoomMsg)
	roommsg.Head = head
	roommsg.Data = data
	roommsg.Uid = uid
	roommsg.V = v

	return roommsg
}

//! 房间逻辑
type Room struct {
	Id       int      `json:"id"`      //! 房间号
	Uid      []int64  `json:"uid"`     //! 房间里的人
	Name     []string `json:"name"`    //! 名字
	ImgUrl   []string `json:"imgurl"`  //! 头像
	Param    []int    `json:"param"`   //! 参数
	Sex      []int    `json:"sex"`     //! 性别
	Robot    []bool   `json:"robot"`   //! 是否是机器人
	Info     string   `json:"info"`    //! 游戏info
	Time     int64    `json:"time"`    //! 解散时间
	Agree    []int64  `json:"agree"`   //!
	ETime    []int    `json:"etime"`   //! 退出时间
	Type     int      `json:"type"`    //! 房间类型
	Begin    bool     `json:"begin"`   //! 游戏是否已经开始
	Step     int      `json:"step"`    //! 当前第几局
	MaxStep  int      `json:"maxstep"` //! 该房间最大局数
	Num      int      `json:"num"`     //! 消耗几张房卡
	ByeTime  int64    `json:"byetime"`
	LiveTime int64    `json:"livetime"` //! 活动时间
	Param1   int      `json:"param1"`   //! 房间附加参数
	Param2   int      `json:"param2"`   //! 房间附加参数2
	Viewer   []int64  `json:"viewer"`   //! 观众
	ClubId   int64    `json:"clubid"`   //! 是否为俱乐部房间
	Agent    bool     `json:"agent"`    //! 是否是代开房间
	Host     int64    `json:"host"`     //! 房主id
	HostName string   `json:"hostname"`
	HostHead string   `json:"hosthead"`
	HostCard int      `json:"hostcard"` //! 房主扣卡
	AACard   int      `json:"aacard"`   //! AA扣卡
	Many     int      `json:"many"`     //! 是否多人  0不多人  1多人不关闭 2多人关闭 3百人场可以解散

	chanlock *sync.RWMutex
	csv      staticfunc.CsvNode //! csv

	reciveChan chan *RoomMsg //! 操作队列

	game Game_Base //! 游戏
}

func NewRoom(id int, _type int, param1 int, param2 int, agent int64, clubid int64) *Room {
	room := new(Room)
	room.Id = id
	room.Type = _type
	room.Begin = false
	room.Uid = make([]int64, 0)
	room.Robot = make([]bool, 0)
	room.Agree = make([]int64, 0)
	room.ETime = make([]int, 0)
	room.Info = ""
	room.LiveTime = time.Now().Unix()
	room.Param1 = param1
	room.Param2 = param2
	room.Viewer = make([]int64, 0)
	room.ClubId = clubid
	if agent > 0 { //! 代理房间
		room.Agent = true
		room.Host = agent
		person := new(Person)
		value := GetServer().DB_GetData("user", agent)
		if string(value) != "" {
			json.Unmarshal(value, &person)
		} else { //! redis读不到，换服务器获取
			var _msg staticfunc.Msg_Uid
			_msg.Uid = agent
			result, _ := GetServer().CallLogin("ServerMethod.ServerMsg", "getperson", &_msg)
			json.Unmarshal(result, &person)
		}
		room.HostHead = person.Imgurl
		room.HostName = person.Name
	}

	return room
}

func (self *Room) Init(num int) {
	self.reciveChan = make(chan *RoomMsg, 2000)
	self.chanlock = new(sync.RWMutex)
	self.csv = staticfunc.GetCsvMgr().Data["game"][self.Type]
	if self.MaxStep <= 0 {
		self.Num = num
		self.MaxStep = num * lib.HF_Atoi(self.csv["step"])
	}
	self.HostCard = 0
	self.AACard = 0
	self.Many = 3
	if self.Type/10000 == 1 { //! 卡五星金币场
		self.game = NewGame_GoldKWX()
	} else if self.Type/10000 == 2 { //! 扎金花金币场
		self.game = NewGame_GoldZJH()
	} else if self.Type/10000 == 3 { //! 牛牛金币场
		self.game = NewGame_GoldNN()
	} else if self.Type/10000 == 4 { //! 豹子王
		self.game = NewGame_GoldBZW()
		self.Many = 1
	} else if self.Type/10000 == 5 { //! 拼天九
		self.game = NewGame_GoldPTJ()
	} else if self.Type/10000 == 6 { //! 百人推筒子
		self.game = NewGame_GoldBrTTZ()
		self.Many = 1
	} else if self.Type/10000 == 7 { //! 推筒子
		self.game = NewGame_GoldTTZ()
	} else if self.Type/10000 == 8 { //! 跑得快
		self.game = NewGame_GoldPDK()
	} else if self.Type/10000 == 9 { //! 神仙夺宝
		self.game = NewGame_GoldSXDB()
		self.Many = 1
	} else if self.Type/10000 == 10 { //! 龙虎斗
		self.game = NewGame_GoldLHD()
		self.Many = 1
	} else if self.Type/10000 == 11 { //! 一夜暴富
		self.game = NewGame_GoldYYBF()
		self.Many = 1
	} else if self.Type/10000 == 12 { //! 摇色子
		self.game = NewGame_GoldYSZ()
		self.Many = 1
	} else if self.Type/10000 == 13 { //! 赛马
		self.game = NewGame_GoldSaiMa()
		self.Many = 1
	} else if self.Type/10000 == 14 { //! 单双
		self.game = NewGame_GoldTB()
		self.Many = 1
	} else if self.Type/10000 == 16 { //!  龙珠探宝
		self.game = NewGame_LZDB()
		self.Many = 2
	} else if self.Type == 1000 { //! 扎金花包房
		self.game = NewGame_GoldZJHRoom()
	} else if self.Type == 2000 { //! 骰宝包房
		self.game = NewGame_GoldTBRoom()
		self.Many = 2
	} else if self.Type/1000 == 3 { //!牛牛包房
		self.game = NewGame_GoldNNRoom()
	} else if self.Type == 4000 { //! 跑的快包房
		self.game = NewGame_GoldPDKRoom()
	} else if self.Type/10000 == 17 { //! 大富翁
		self.game = NewGame_DFW()
		self.Many = 2
	} else if self.Type/10000 == 18 { //! 翻牌机
		self.game = NewGame_FPJ()
		self.Many = 2
	} else if self.Type/10000 == 19 { //! 疯狂翻牌机
		self.game = NewGame_GoldFKFPJ()
		self.Many = 2
	} else if self.Type/10000 == 20 { //!  名品汇
		self.game = NewGame_GoldMPH()
		self.Many = 1
	} else if self.Type/10000 == 21 { //！ 红黑大战
		self.game = NewGame_GoldHHDZ()
		self.Many = 1
	} else if self.Type/10000 == 22 { //! 腾讯龙虎斗
		self.game = NewGame_TXLHD()
		self.Many = 1
	} else if self.Type/10000 == 23 { //! 百人牛牛
		self.game = NewGame_GoldBrNN()
		self.Many = 1
	} else if self.Type/10000 == 24 { //! 鱼虾蟹
		self.game = NewGame_GoldYXX()
		self.Many = 1
	} else if self.Type/10000 == 25 { //! 捕鱼
		self.game = NewGame_Fishing()
		self.Many = 3
	} else if self.Type/10000 == 26 { //! 百家乐
		self.game = NewGame_GoldBJL()
		self.Many = 1
	} else if self.Type/10000 == 27 { //! 水浒传
		self.game = NewGame_GoldSHZ()
		self.Many = 2
	} else if self.Type/10000 == 28 {
		self.game = NewGame_Fishing_LKPY()
		self.Many = 3
	} else if self.Type/10000 == 29 {
		self.game = NewGame_GoldDDZ()
	} else if self.Type/10000 == 30 { //! 红包扫雷
		self.game = NewGame_GoldHBSL()
		self.Many = 1
	} else if self.Type/10000 == 31 { //! 大赢家拉霸
		self.game = NewGame_DYJLB()
		self.Many = 2
	} else if self.Type/10000 == 32 { //! 绝地求生
		self.game = NewGame_JDQS()
		self.Many = 2
	} else if self.Type == 1 { //! 牛牛
		self.game = NewGame_NiuNiu()
	} else if self.Type >= 2 && self.Type <= 5 || self.Type == 13 || self.Type == 14 { //! 卡五星
		self.game = NewGame_KWX()
	} else if self.Type == 6 || self.Type == 8 { //! 斗地主
		self.game = NewGame_DDZ()
	} else if self.Type == 7 { //! 炸金花
		self.game = NewGame_ZJH()
	} else if self.Type == 10 { //! 十点半
		self.game = NewGame_TenHalf()
	} else if self.Type == 11 { //! 吃火锅
		self.game = NewGame_EatHot()
	} else if self.Type == 16 { //! 江西牛牛
		self.game = NewGame_NiuNiuJX()
		if self.MaxStep == 25 {
			self.AACard = 1
		} else {
			self.HostCard = self.MaxStep / lib.HF_Atoi(self.csv["step"])
		}
	} else if self.Type == 17 { //! 推筒子
		self.game = NewGame_T()
	} else if self.Type == 19 { //! 拼天九
		self.game = NewGame_PTJ()
	} else if self.Type == 24 { //！三公演义
		self.game = NewGame_SGYY()
		if self.Param2/10%10 == 0 { //房主
			if self.MaxStep == 10 {
				self.HostCard = 3
			} else {
				self.HostCard = 6
			}
		} else { //AA
			if self.MaxStep == 10 {
				self.AACard = 1
			} else {
				self.AACard = 2
			}
		}
	} else if self.Type == 25 { //逍遥炸金花
		self.game = NewGame_XYZJH()
	} else if self.Type == 36 { //！扫雷
		self.game = NewGame_SaoLei()
	} else if self.Type == 51 { //! 十三道
		if self.Param2%10 == 0 {
			self.AACard = (self.MaxStep / 10) * 2
		} else {
			self.HostCard = (self.MaxStep / 10) * (self.Param1 / 10000 % 10) * 2
		}
		self.game = NewGame_SSD()
	} else if self.Type == 65 || self.Type == 66 || self.Type == 67 || self.Type == 68 || self.Type == 69 { //! 明牌抢庄 八人明牌 自由抢庄
		self.game = NewGame_MPQZ()
	} else if self.Type == 77 {
		self.game = NewGame_WZQ()
	} else if self.Type == 78 {
		self.game = NewGame_XJPDK()
	} else if self.Type == 100 {
		//self.game = NewGame_HFBH()
	}
	if self.HostCard == 0 && self.AACard == 0 { //! 没有选择房卡模式，默认
		self.HostCard = self.MaxStep / lib.HF_Atoi(self.csv["step"])
		if self.Type == 28 || GetServer().Con.Type >= 3 || self.Type == 49 {
			self.HostCard = 0
			self.AACard = 0
		}
	}
	if self.Type == 77 { //! 五子棋为下分游戏，不消耗房卡
		self.HostCard = 0
		self.AACard = 0
	}
	if self.Info != "" {
		json.Unmarshal([]byte(self.Info), self.game)
	} else {
		self.flush()
	}
	self.game.OnInit(self)
	go self.run()
	go self.runLive()
	if self.Time != 0 { //! 解散时间不为0
		go self.dissmissThread()
	}
	if self.ByeTime != 0 { //! bye时间不为0
		go self.byeThread()
	}

	if self.Many > 0 {
		GetServer().Wait.Add(1)
	}
}

//! 发送操作
func (self *Room) Operator(op *RoomMsg) {
	self.chanlock.Lock()
	defer self.chanlock.Unlock()

	if self.reciveChan != nil {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "op:", op.Head)
		self.reciveChan <- op
	}
}

//! 设置开始结束
func (self *Room) SetBegin(begin bool) {
	if begin && self.Step == 0 && (self.Agent || self.ClubId != 0) { //! 第一局开始
		var msg staticfunc.Msg_BeginRoom
		msg.Id = GetServer().Con.Id
		msg.RoomId = self.Id
		msg.Host = self.Host
		msg.ClubId = self.ClubId
		GetServer().CallLogin("ServerMethod.ServerMsg", "beginroom", &msg)
	}

	self.Begin = begin

	if begin {
		self.Step++
	}

	if self.Agent { //! 代开的房间在创建房间时扣房卡
		return
	}

	if lib.HF_Atoi(self.csv["view"]) == 1 && self.AACard == 0 { //! 观众模式在创建房间时扣卡，故这里不扣
		return
	}

	if GetServer().Con.Type >= 3 { //! 匹配模式不扣房卡
		return
	}

	if self.Begin { //! 开始游戏
		if self.Step == 1 { //! 第一局开始
			if lib.HF_Atoi(self.csv["begincost"]) == 0 { //! 开始时扣房卡
				if self.HostCard > 0 {
					GetRoomMgr().CostCard(self.Uid[0], lib.HF_Atoi(self.csv["costtype"]), self.HostCard, self)
				} else if self.AACard > 0 {
					for i := 0; i < len(self.Uid); i++ {
						GetRoomMgr().CostCard(self.Uid[i], lib.HF_Atoi(self.csv["costtype"]), self.AACard, self)
					}
				}
			}
		}
	} else {
		if self.Step == 1 { //! 第一局结束
			if lib.HF_Atoi(self.csv["begincost"]) == 1 { //! 结束时扣房卡
				if self.HostCard > 0 {
					GetRoomMgr().CostCard(self.Uid[0], lib.HF_Atoi(self.csv["costtype"]), self.HostCard, self)
				} else if self.AACard > 0 {
					for i := 0; i < len(self.Uid); i++ {
						GetRoomMgr().CostCard(self.Uid[i], lib.HF_Atoi(self.csv["costtype"]), self.AACard, self)
					}
				}
			}
		}
	}
}

//! 增加记录
func (self *Room) AddRecord(info string) {
	for i := 0; i < len(self.Uid); i++ {
		if self.Robot[i] {
			continue
		}
		if self.Type == 25 {
			GetServer().InsertRecord(7, self.Uid[i], info, 0)
		} else {
			GetServer().InsertRecord(self.Type, self.Uid[i], info, 0)
		}
	}
}

//! 是否局数达到
func (self *Room) IsBye() bool {
	return self.Step >= self.MaxStep
}

//! 是否有uid
func (self *Room) IsHasUid(uid int64) bool {
	for i := 0; i < len(self.Uid); i++ {
		if self.Uid[i] == uid {
			return true
		}
	}
	return false
}

//! 是否允许加入
func (self *Room) isJoin(uid int64) bool {
	if self.ByeTime != 0 {
		return false
	}

	if !self.IsLive() {
		return false
	}

	for i := 0; i < len(self.Uid); i++ {
		if self.Uid[i] == uid { //! 在房间里，则允许加入
			return true
		}
	}

	//! 若不在房间
	if self.Time != 0 { //! 该房间已经申请解散了
		return false
	}

	if GetServer().Con.Type < 3 { //! 金币场开始了也能加入房间
		if lib.HF_Atoi(self.csv["beginenter"]) == 0 {
			if self.Step > 0 { //! 已经开始了
				return false
			}

			if self.Begin { //! 该房间游戏已经开始
				return false
			}
		}
	}

	if (len(self.Uid) + len(self.Viewer)) >= lib.HF_Atoi(self.csv["maxnum"]) { //! 人数已上限
		return false
	}

	return true
}

func (self *Room) isView(uid int64) bool { //! 是否能观战
	if self.ByeTime != 0 {
		return false
	}

	if !self.IsLive() {
		return false
	}

	for i := 0; i < len(self.Viewer); i++ {
		if self.Viewer[i] == uid { //! 正在观战
			return true
		}
	}

	//! 若不在房间
	if self.Time != 0 { //! 该房间已经申请解散了
		return false
	}

	if self.Many == 0 && self.Param2%10 == 1 { //! 不是百人场，观战房间param % 10默认为是否开始了还能观战
		if self.Step > 0 { //! 已经开始了
			return false
		}

		if self.Begin { //! 该房间游戏已经开始
			return false
		}
	}

	return true
}

func (self *Room) run() {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
			self.clear(true)
		}
	}()

	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case <-ticker.C:
			if self.reciveChan == nil {
				return
			}

			//! 判断没有被退出的人
			exit := false
			for i := 0; i < len(self.ETime); {
				self.ETime[i]--
				if self.ETime[i] <= 0 {
					if len(self.Uid) <= 2 { //! 只剩下两个人
						self.game.OnBye()
						self.clear(true)
						ticker.Stop()
						return
					}
					exit = true
					for j := 0; j < len(self.Uid); j++ {
						if self.Uid[j] == self.Agree[i] {
							copy(self.Uid[j:], self.Uid[j+1:])
							self.Uid = self.Uid[:len(self.Uid)-1]
							copy(self.Name[j:], self.Name[j+1:])
							self.Name = self.Name[:len(self.Name)-1]
							copy(self.ImgUrl[j:], self.ImgUrl[j+1:])
							self.ImgUrl = self.ImgUrl[:len(self.ImgUrl)-1]
							copy(self.Sex[j:], self.Sex[j+1:])
							self.Sex = self.Sex[:len(self.Sex)-1]
							copy(self.Param[j:], self.Param[j+1:])
							self.Param = self.Param[:len(self.Param)-1]
							copy(self.Robot[j:], self.Robot[j+1:])
							self.Robot = self.Robot[:len(self.Robot)-1]

							self.broadCastMsg("roominfo", self.getRoomMsg())
							self.game.OnExit(self.Agree[i])

							var _msg staticfunc.Msg_JoinFail
							_msg.Uid = self.Agree[i]
							_msg.Id = GetServer().Con.Id
							_msg.Room = self.Id
							_msg.GameType = self.Type
							GetServer().CallLogin("ServerMethod.ServerMsg", "joinfail", &_msg)

							person := GetPersonMgr().GetPerson(self.Agree[i])
							if person != nil {
								var msg Msg_ExitRoom
								person.SendMsg("exitroom", &msg)
								person.CloseSession()
								GetPersonMgr().DelPerson(self.Agree[i])
							}
							break
						}
					}

					copy(self.Agree[i:], self.Agree[i+1:])
					self.Agree = self.Agree[:len(self.Agree)-1]
					copy(self.ETime[i:], self.ETime[i+1:])
					self.ETime = self.ETime[:len(self.ETime)-1]
					self.flush()
				} else {
					i++
				}
			}
			if exit { //! 有人被系统判断退出，则发消息，庄家可以更新
				var msg Msg_DissmissRoom
				msg.Agree = self.Agree
				msg.Time = self.Time
				msg.ETime = self.ETime
				self.broadCastMsg("dissmissroom", &msg)
			}

			if self.Many > 0 && GetServer().ShutDown { //! 百人场在关服的时候必须关闭场次
				self.clear(true)
				return
			}

			self.game.OnTime()
		case op := <-self.reciveChan:
			switch op.Head {
			case "joinroom": //! 加入房间
				self.JoinRoom(op.V.(*Person))
			case "dissmissroom": //! 申请解散
				self.LiveTime = time.Now().Unix()
				self.dismiss(op.Uid, op.V.(*Msg_DissRoom).Type)
			case "nodissmissroom": //! 不解散
				self.LiveTime = time.Now().Unix()
				self.nodismiss(op.Uid)
			case "broadcast": //! 广播
				self.LiveTime = time.Now().Unix()
				self.broadCastMsg(op.Data, op.V)
			case "byeroom": //! 结算房间
				if self.Step > 0 {
					self.game.OnBye()
				}
				self.clear(true)
				ticker.Stop()
				return
			case "delroom": //! 解散房间
				self.clear(false)
				ticker.Stop()
				return
			default:
				self.LiveTime = time.Now().Unix()
				self.game.OnMsg(op)
			}
		}
	}
}

//! 加入机器人
//! 0成功  1机器人不足   2该房间已满
func (self *Room) AddRobot(gametype int, money int, limitmin int, limitmax int) int {
	robots := lib.GetRobotMgr().InitRobotFromGame(gametype, 1, money, limitmin, limitmax)
	if len(robots) == 0 {
		return 1
	}
	robot := robots[0]
	if !self.isJoin(robot.Id) {
		for i := 0; i < len(robots); i++ {
			lib.GetRobotMgr().BackRobot(robots[i].Id, false)
		}
		return 2
	}

	self.LiveTime = time.Now().Unix()

	for i := 0; i < len(self.Uid); i++ {
		if self.Uid[i] == robot.Id {
			if self.Robot[i] {
				return 0
			} else {
				for j := 0; j < len(robots); j++ {
					lib.GetRobotMgr().BackRobot(robots[j].Id, false)
				}
				return 2
			}
		}
	}

	self.Uid = append(self.Uid, robot.Id)
	self.Name = append(self.Name, robot.Name)
	self.ImgUrl = append(self.ImgUrl, robot.Head)
	self.Sex = append(self.Sex, robot.Sex)
	self.Robot = append(self.Robot, true)
	self.Param = append(self.Param, robot.GetMoney())

	self.flush()

	self.broadCastMsg("roominfo", self.getRoomMsg()) //! 广播房间信息
	self.game.OnRobot(robot)                         //! 发送游戏信息

	var msg staticfunc.Msg_AddRobot
	msg.Id = GetServer().Con.Id
	msg.Room = self.Id
	msg.Uid = robot.Id
	msg.GameType = self.Type
	msg.Num = len(self.Uid)
	msg.IP = robot.IP
	GetServer().CallLogin("ServerMethod.ServerMsg", "addrobot", &msg)

	return 0
}

func (self *Room) JoinRoom(person *Person) {
	if lib.HF_Atoi(self.csv["view"]) == 1 { //! 观战模式
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "观战模式")
		finduid := false
		for i := 0; i < len(self.Uid); i++ {
			if self.Uid[i] == person.Uid {
				finduid = true
				break
			}
		}
		if !finduid { //! 没有坐下，默认为观众
			if !self.isView(person.Uid) { //! 不能加入观众
				var _msg staticfunc.Msg_JoinFail
				_msg.Uid = person.Uid
				_msg.Id = GetServer().Con.Id
				_msg.Room = self.Id
				_msg.GameType = self.Type
				GetServer().CallLogin("ServerMethod.ServerMsg", "joinfail", &_msg)

				var msg Msg_JoinRoomFail
				msg.Result = 2
				person.SendMsg("joinroomfail", &msg)
				person.CloseSession()
			} else {
				if 28 == self.Type {
					if (person.Gold < self.Param1*20 && self.Param2/100%10 != 3) || (person.Gold < self.Param1+20 && self.Param2/100%10 == 3) {
						var _msg staticfunc.Msg_JoinFail
						_msg.Uid = person.Uid
						_msg.Id = GetServer().Con.Id
						_msg.Room = self.Id
						_msg.GameType = self.Type
						GetServer().CallLogin("ServerMethod.ServerMsg", "joinfail", &_msg)

						var msg Msg_JoinRoomFail
						msg.Result = 4
						person.SendMsg("joinroomfail", &msg)
						person.CloseSession()
						return
					}
				}
				//! 发送游戏信息
				self.LiveTime = time.Now().Unix()
				GetPersonMgr().AddPerson(person)
				if self.Host == 0 {
					self.Host = person.Uid
					self.HostName = person.Name
					self.HostHead = person.Imgurl
					//! 观战模式扣卡
					if self.HostCard > 0 {
						GetRoomMgr().CostCard(person.Uid, lib.HF_Atoi(self.csv["costtype"]), self.HostCard, self)
					}
				}
				self.addView(person)
				person.SendMsg("roominfo", self.getRoomMsg())
				self.game.OnSendInfo(person)
			}
			return
		}
	}
	join := false
	for i := 0; i < len(self.Uid); i++ {
		if self.Uid[i] == person.Uid { //! 在房间里，则允许加入
			join = true
			break
		}
	}
	if !join { //! 不在房间，判断房卡
		join = true
		if !self.Agent { //! 不是代开房，判断房卡
			if self.HostCard > 0 && len(self.Uid) == 0 {
				if person.Card+person.Gold < self.HostCard {
					join = false
				}
			} else if self.AACard > 0 {
				if person.Card+person.Gold < self.AACard { //! 房卡不足
					join = false
				}
			}
		}
	}
	if !join {
		//! 通知服务器
		var _msg staticfunc.Msg_JoinFail
		_msg.Uid = person.Uid
		_msg.Id = GetServer().Con.Id
		_msg.Room = self.Id
		_msg.GameType = self.Type
		GetServer().CallLogin("ServerMethod.ServerMsg", "joinfail", &_msg)

		var msg Msg_JoinRoomFail
		msg.Result = 4
		person.SendMsg("joinroomfail", &msg)
		person.CloseSession()
	} else {
		if !self.isJoin(person.Uid) { //! 无法加入
			//! 通知服务器
			var _msg staticfunc.Msg_JoinFail
			_msg.Uid = person.Uid
			_msg.Id = GetServer().Con.Id
			_msg.Room = self.Id
			_msg.GameType = self.Type
			GetServer().CallLogin("ServerMethod.ServerMsg", "joinfail", &_msg)

			var msg Msg_JoinRoomFail
			msg.Result = 2
			person.SendMsg("joinroomfail", &msg)
			person.CloseSession()
		} else {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "joinsucceed")
			self.LiveTime = time.Now().Unix()
			GetPersonMgr().AddPerson(person)
			if self.Agent && self.ClubId == 0 { //! 代理模式加人
				var msg staticfunc.Msg_ViewRoomNum
				msg.Host = self.Host
				msg.RoomId = self.Id
				msg.Num = 1
				msg.Uid = person.Uid
				msg.Name = person.Name
				msg.Head = person.Imgurl
				GetServer().CallLogin("ServerMethod.ServerMsg", "viewroomnum", &msg)
			}
			if self.Host == 0 {
				self.Host = person.Uid
				self.HostName = person.Name
				self.HostHead = person.Imgurl
			}
			self.addUid(person)
			self.broadCastMsg("roominfo", self.getRoomMsg()) //! 广播房间信息
			self.game.OnSendInfo(person)                     //! 发送游戏信息
		}
	}
}

func (self *Room) runLive() {
	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		if self.LiveTime == 0 {
			break
		}

		if self.reciveChan == nil {
			break
		}

		if !self.IsLive() {
			self.Operator(NewRoomMsg("byeroom", "now", 0, nil))
			break
		}
	}

	//！ 关掉定时器
	ticker.Stop()
}

//! 清理房间
func (self *Room) clear(send bool) {
	if self.Step == 0 && self.Host > 0 && self.HostCard > 0 {
		if lib.HF_Atoi(self.csv["view"]) == 1 || self.Agent { //! 观战或代开模式一局未开始,归还卡
			GetRoomMgr().AddCard(self.Host, lib.HF_Atoi(self.csv["costtype"]), self.HostCard, self.Type)
		}
	}

	//! 结算所有剩余玩家
	self.game.OnBalance()

	GetRoomMgr().DelRoom(self.Id)
	self.Time = 0

	self.chanlock.Lock()
	if self.reciveChan != nil {
		close(self.reciveChan)
		self.reciveChan = nil
	}
	self.chanlock.Unlock()

	self.ByeTime = 0
	self.LiveTime = 0
	self.ETime = make([]int, 0)
	self.Agree = make([]int64, 0)

	//! 清理玩家
	for i := 0; i < len(self.Uid); i++ {
		if self.Robot[i] { //! 不算机器人
			lib.GetRobotMgr().BackRobot(self.Uid[i], false)
			continue
		}
		person := GetPersonMgr().GetPerson(self.Uid[i])
		if person == nil {
			continue
		}
		if person.room == nil || person.room.Id != self.Id {
			continue
		}

		//! 告诉玩家退出
		var msg Msg_ExitRoom
		person.SendMsg("exitroom", &msg)

		//! 断开连接
		person.CloseSession()
		GetPersonMgr().DelPerson(self.Uid[i])
	}

	if send {
		var msg staticfunc.Msg_DelRoom
		msg.Id = GetServer().Con.Id
		msg.RoomId = self.Id
		msg.Uid = self.Uid
		msg.Host = self.Host
		msg.Agent = self.Agent
		msg.ClubId = self.ClubId
		msg.GameType = self.Type
		GetServer().CallLogin("ServerMethod.ServerMsg", "delroom", &msg)
	}

	//! 清理观战
	if lib.HF_Atoi(self.csv["view"]) == 1 {
		for i := 0; i < len(self.Viewer); i++ {
			person := GetPersonMgr().GetPerson(self.Viewer[i])
			if person == nil {
				continue
			}
			if person.room == nil || person.room.Id != self.Id {
				continue
			}

			//! 告诉玩家退出
			var msg Msg_ExitRoom
			person.SendMsg("exitroom", &msg)

			//! 断开连接
			person.CloseSession()
			GetPersonMgr().DelPerson(self.Viewer[i])
		}

		if send && len(self.Viewer) > 0 {
			var msg staticfunc.Msg_DelRoom
			msg.Id = GetServer().Con.Id
			msg.RoomId = self.Id
			msg.Uid = self.Viewer
			msg.Host = self.Host
			msg.Agent = self.Agent
			msg.ClubId = self.ClubId
			msg.GameType = self.Type
			GetServer().CallLogin("ServerMethod.ServerMsg", "delroom", &msg)
		}
	}

	if self.Many > 0 { //! 豹子王关闭服务器必须关闭场次
		GetServer().Wait.Done()
	}
}

//! 广播消息
func (self *Room) broadCastMsg(head string, v interface{}) {
	if head == "lineperson" {
		if !v.(*Msg_LinePerson).Line { //! 下线判断
			person := GetPersonMgr().GetPerson(v.(*Msg_LinePerson).Uid)
			if person != nil && person.session != nil {
				return
			}
		}
	}

	msg := lib.HF_EncodeMsg(head, v, true)

	for i := 0; i < len(self.Uid); i++ {
		if self.Robot[i] { //! 不算机器人
			continue
		}
		person := GetPersonMgr().GetPerson(self.Uid[i])
		if person == nil {
			continue
		}
		if person.room == nil || person.room.Id != self.Id {
			continue
		}

		person.SendByteMsg(msg)
	}

	if lib.HF_Atoi(self.csv["view"]) == 1 { //! 观战模式下发给观战的人
		for i := 0; i < len(self.Viewer); i++ {
			person := GetPersonMgr().GetPerson(self.Viewer[i])
			if person == nil {
				continue
			}
			if person.room == nil || person.room.Id != self.Id {
				continue
			}

			person.SendByteMsg(msg)
		}
	}
}

//! 广播消息
func (self *Room) broadCastMsgView(head string, v interface{}) {
	if head == "lineperson" {
		if !v.(*Msg_LinePerson).Line { //! 下线判断
			person := GetPersonMgr().GetPerson(v.(*Msg_LinePerson).Uid)
			if person != nil && person.session != nil {
				return
			}
		}
	}

	if lib.HF_Atoi(self.csv["view"]) == 1 { //! 观战模式下发给观战的人
		for i := 0; i < len(self.Viewer); i++ {
			person := GetPersonMgr().GetPerson(self.Viewer[i])
			if person == nil {
				continue
			}
			if person.room == nil || person.room.Id != self.Id {
				continue
			}

			person.SendMsg(head, v)
		}
	}
}

//! 发送消息
func (self *Room) SendMsg(uid int64, head string, v interface{}) {
	for i := 0; i < len(self.Uid); i++ {
		if self.Uid[i] == uid {
			if self.Robot[i] {
				return
			}
			break
		}
	}

	person := GetPersonMgr().GetPerson(uid)
	if person == nil {
		return
	}
	if person.room == nil || person.room.Id != self.Id {
		return
	}

	person.SendMsg(head, v)
}

//! 发送错误提示消息
func (self *Room) SendErr(uid int64, info string) {
	for i := 0; i < len(self.Uid); i++ {
		if self.Uid[i] == uid {
			if self.Robot[i] {
				return
			}
			break
		}
	}

	person := GetPersonMgr().GetPerson(uid)
	if person == nil {
		return
	}
	if person.room == nil || person.room.Id != self.Id {
		return
	}

	person.SendErr(info)
}

//! 加入人
func (self *Room) addUid(person *Person) {
	self.LiveTime = time.Now().Unix()

	for i := 0; i < len(self.Uid); i++ {
		if self.Uid[i] == person.Uid {
			if self.Robot[i] {
				self.Robot[i] = false
			}
			return
		}
	}

	self.Uid = append(self.Uid, person.Uid)
	self.Name = append(self.Name, person.Name)
	self.ImgUrl = append(self.ImgUrl, person.Imgurl)
	self.Sex = append(self.Sex, person.Sex)
	self.Robot = append(self.Robot, false)
	if person.Param == 0 {
		if self.Type == 77 {
			self.Param = append(self.Param, lib.HF_MaxInt(person.Gold-person.BindGold, 0))
		} else {
			self.Param = append(self.Param, lib.HF_MaxInt(person.Gold, 0))
		}
	} else {
		self.Param = append(self.Param, person.Param)
	}

	self.flush()
}

//! 加入观战
func (self *Room) addView(person *Person) {
	self.LiveTime = time.Now().Unix()

	for i := 0; i < len(self.Viewer); i++ {
		if self.Viewer[i] == person.Uid {
			return
		}
	}

	self.Viewer = append(self.Viewer, person.Uid)

	self.flush()
}

//! 观战的人坐下
func (self *Room) Seat(uid int64) bool {
	if lib.HF_Atoi(self.csv["view"]) == 0 {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "view为0错误")
		return false
	}

	if len(self.Uid) >= lib.HF_Atoi(self.csv["maxnum"]) { //! 人数已上限
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "人数上限错误")
		return false
	}

	for i := 0; i < len(self.Viewer); i++ {
		if self.Viewer[i] == uid {
			person := GetPersonMgr().GetPerson(uid)
			if person == nil {
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "Viewer为空错误")
				return false
			}

			if self.AACard > 0 { //! AA模式坐下需要判断房卡
				if person.Card+person.Gold < self.AACard {
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "AA房卡不足")
					person.SendMsg("notenoughcard", nil)
					return false
				}

				if self.Step > 0 { //! 如果已经开始了，则直接扣卡
					GetRoomMgr().CostCard(uid, lib.HF_Atoi(self.csv["costtype"]), self.AACard, self)
				}
			}

			copy(self.Viewer[i:], self.Viewer[i+1:])
			self.Viewer = self.Viewer[:len(self.Viewer)-1]

			self.addUid(person)
			self.broadCastMsg("roomseat", self.getRoomMsg())

			if self.ClubId == 0 && GetServer().Con.Type < 3 {
				var msg staticfunc.Msg_ViewRoomNum
				msg.Host = self.Host
				msg.RoomId = self.Id
				msg.Num = 1
				msg.Uid = person.Uid
				msg.Name = person.Name
				msg.Head = person.Imgurl
				GetServer().CallLogin("ServerMethod.ServerMsg", "viewroomnum", &msg)
			}

			return true
		}
	}

	return false
}

//! 得到房间消息
func (self *Room) getRoomMsg() interface{} {
	var msg Msg_RoomInfo
	msg.RoomId = self.Id
	msg.Host = self.Host
	msg.Name = self.HostName
	msg.Head = self.HostHead
	msg.Agent = self.Agent
	for i := 0; i < len(self.Uid); i++ {
		var son Son_PersonInfo
		son.Uid = self.Uid[i]
		son.Name = self.Name[i]
		son.ImgUrl = self.ImgUrl[i]
		if i >= len(self.Param) {
			son.Param = 0
		} else {
			son.Param = self.Param[i]
		}
		if i >= len(self.Sex) {
			son.Sex = 0
		} else {
			son.Sex = self.Sex[i]
		}
		if self.Robot[i] {
			robot := lib.GetRobotMgr().GetRobotFromId(self.Uid[i])
			if robot == nil {
				son.Ip = ""
				son.Line = false
				son.Address = ""
				son.Latitude = ""
				son.Longitude = ""
			} else {
				son.Ip = robot.IP
				son.Line = true
				son.Address = robot.Address
			}
		} else {
			person := GetPersonMgr().GetPerson(self.Uid[i])
			if person == nil || person.session == nil {
				son.Ip = ""
				son.Line = false
				son.Address = ""
				son.Latitude = ""
				son.Longitude = ""
			} else {
				if GetServer().IsGM(person.Uid) {
					son.Ip = ""
				} else {
					son.Ip = person.ip
				}
				son.Line = person.line
				son.Address = person.minfo.Address
				son.Latitude = person.minfo.Latitude
				son.Longitude = person.minfo.Longitude
			}
		}
		msg.Person = append(msg.Person, son)
	}
	msg.Time = self.Time
	msg.Agree = self.Agree
	msg.ETime = self.ETime
	msg.Step = self.Step
	msg.MaxStep = self.MaxStep
	msg.Param1 = self.Param1
	msg.Param2 = self.Param2
	msg.Type = self.Type

	return &msg
}

func (self *Room) flush() {
	self.Info = lib.HF_JtoA(self.game)
	GetServer().DB_SetData("roominfo", int64(self.Id), lib.HF_JtoB(self))
}

func (self *Room) dismiss(uid int64, _type int) {
	if _type == 1 && self.Host == uid && self.Step == 0 && GetServer().Con.Type < 3 { //! 不是金币场，房主可以在未开始的时候强制解散
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "解散房间")
		self.clear(true)
		return
	}

	if lib.HF_Atoi(self.csv["view"]) == 1 { //! 可观战
		if self.Many > 0 { //! 百人场
			if self.game.OnIsDealer(uid) {
				self.SendErr(uid, "请下庄后再退出游戏")
				return
			} else if self.game.OnIsBets(uid) {
				self.SendErr(uid, "本局已下注，请等待本局结束。")
				return
			}
		}

		for i := 0; i < len(self.Viewer); i++ {
			if self.Viewer[i] == uid { //! 是观众随意进出
				copy(self.Viewer[i:], self.Viewer[i+1:])
				self.Viewer = self.Viewer[:len(self.Viewer)-1]

				self.game.OnExit(uid)

				var _msg staticfunc.Msg_JoinFail
				_msg.Uid = uid
				_msg.Id = GetServer().Con.Id
				_msg.Room = self.Id
				_msg.GameType = self.Type
				GetServer().CallLogin("ServerMethod.ServerMsg", "joinfail", &_msg)

				person := GetPersonMgr().GetPerson(uid)
				if person != nil {
					var msg Msg_ExitRoom
					person.SendMsg("exitroom", &msg)
					person.CloseSession()
					GetPersonMgr().DelPerson(uid)
				}

				self.flush()

				if len(self.Viewer) == 0 && len(self.Uid) == 0 && self.Many != 1 { //! 没有观众也没有玩家
					lib.GetLogMgr().Output(lib.LOG_DEBUG, "清理房间")
					self.clear(true)
				}

				return
			}
		}
	}
	if self.Many > 0 && self.Many != 3 { //! 百人场不能解散
		return
	}
	if GetServer().Con.Type >= 3 || self.Type == 77 { //金币场或五子棋
		if !self.Begin {
			self.KickPerson(uid, 0)
		} else {
			if self.Type >= 20000 && self.Type < 30000 { //! 炸金花金币场
				if self.game.(*Game_GoldZJH).IsDiscard(uid) {
					self.KickPerson(uid, 0)
					return
				}
			}
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏过程中禁止退出")
			person := GetPersonMgr().GetPerson(uid)
			if person != nil {
				person.SendErr("游戏过程中无法退出")
			}
		}
		return
	}

	if self.Step == 0 { //! 一局都未开始
		view := lib.HF_Atoi(self.csv["view"])
		if view == 0 && uid == self.Uid[0] && !self.Agent { //! 非观战模式下房主退出，直接解散
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "解散房间")
			self.clear(true)
		} else {
			self.KickPerson(uid, 0)
			if (view == 1 || self.Agent) && self.ClubId == 0 {
				var msg staticfunc.Msg_ViewRoomNum
				msg.Host = self.Host
				msg.RoomId = self.Id
				msg.Num = -1
				msg.Uid = uid
				msg.Name = ""
				msg.Head = ""
				GetServer().CallLogin("ServerMethod.ServerMsg", "viewroomnum", &msg)
			}
		}
		return
	}

	if lib.HF_Atoi(self.csv["hostagree"]) == 1 { //! 房主同意模式
		if self.Begin { //! 游戏开始时不能解散
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "游戏开始时不能解散hostagree")
			return
		}
		if self.game.OnIsDealer(uid) { //! 庄家
			if len(self.Agree) == 0 { //! 庄家不能申请退出
				return
			}
			if len(self.Uid) > 2 {
				for i := 0; i < len(self.Uid); i++ {
					if self.Uid[i] == self.Agree[0] {
						copy(self.Uid[i:], self.Uid[i+1:])
						self.Uid = self.Uid[:len(self.Uid)-1]
						copy(self.Name[i:], self.Name[i+1:])
						self.Name = self.Name[:len(self.Name)-1]
						copy(self.ImgUrl[i:], self.ImgUrl[i+1:])
						self.ImgUrl = self.ImgUrl[:len(self.ImgUrl)-1]
						copy(self.Sex[i:], self.Sex[i+1:])
						self.Sex = self.Sex[:len(self.Sex)-1]
						copy(self.Param[i:], self.Param[i+1:])
						self.Param = self.Param[:len(self.Param)-1]
						copy(self.Robot[i:], self.Robot[i+1:])
						self.Robot = self.Robot[:len(self.Robot)-1]

						self.broadCastMsg("roominfo", self.getRoomMsg())
						self.game.OnExit(self.Agree[0])

						var _msg staticfunc.Msg_JoinFail
						_msg.Uid = self.Agree[0]
						_msg.Id = GetServer().Con.Id
						_msg.Room = self.Id
						_msg.GameType = self.Type
						GetServer().CallLogin("ServerMethod.ServerMsg", "joinfail", &_msg)

						person := GetPersonMgr().GetPerson(self.Agree[0])
						if person != nil {
							var msg Msg_ExitRoom
							person.SendMsg("exitroom", &msg)
							person.CloseSession()
							GetPersonMgr().DelPerson(self.Agree[0])
						}
						break
					}
				}
				self.Agree = self.Agree[1:]
				self.ETime = self.ETime[1:]
				self.flush()
			} else {
				self.game.OnBye()
				self.clear(true)
			}
		} else {
			for i := 0; i < len(self.Agree); i++ {
				if self.Agree[i] == uid {
					return
				}
			}
			self.Agree = append(self.Agree, uid)
			disstime := lib.HF_Atoi(self.csv["disstime"])
			if disstime == 0 {
				disstime = 180
			}
			self.ETime = append(self.ETime, disstime)
			self.flush()

			var msg Msg_DissmissRoom
			msg.Agree = self.Agree
			msg.Time = self.Time
			msg.ETime = self.ETime
			self.broadCastMsg("dissmissroom", &msg)
		}
	} else { //! 自由解散模式
		for i := 0; i < len(self.Agree); i++ {
			if self.Agree[i] == uid {
				return
			}
		}

		self.Agree = append(self.Agree, uid)
		if len(self.Agree) == len(self.Uid) || (self.Type == 82 && len(self.Agree) >= 2) { //! 房间内所有人都同意解散
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "解散房间")
			self.game.OnBye()
			self.clear(true)
		} else {
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "申请解散", lib.HF_Atoi(self.csv["disstime"]))
			if self.Time == 0 {
				disstime := lib.HF_Atoi(self.csv["disstime"])
				if disstime == 0 {
					disstime = 180
				}
				self.Time = time.Now().Unix() + int64(disstime) //! 5分钟后解散
				go self.dissmissThread()
			}
			self.flush()

			var msg Msg_DissmissRoom
			msg.Agree = self.Agree
			msg.Time = self.Time
			msg.ETime = self.ETime
			self.broadCastMsg("dissmissroom", &msg)
		}
	}
}

func (self *Room) nodismiss(uid int64) {
	if lib.HF_Atoi(self.csv["hostagree"]) == 1 {
		if !self.game.OnIsDealer(uid) { //! 不是房主，不能反对
			return
		}

		if len(self.Agree) == 0 {
			return
		}

		person := GetPersonMgr().GetPerson(self.Agree[0])
		if person != nil {
			var msg staticfunc.Msg_Uid
			msg.Uid = uid
			person.SendMsg("nodissmissroom", &msg)
		}

		self.Agree = self.Agree[1:]
		self.ETime = self.ETime[1:]
	} else {
		find := false
		for i := 0; i < len(self.Uid); i++ {
			if self.Uid[i] == uid {
				find = true
				break
			}
		}
		if !find { //! 不是可以反对的人
			return
		}

		for i := 0; i < len(self.Agree); i++ {
			if self.Agree[i] == uid {
				return
			}
		}

		self.Agree = make([]int64, 0)
		self.Time = 0
		self.flush()

		var msg staticfunc.Msg_Uid
		msg.Uid = uid
		self.broadCastMsg("nodissmissroom", &msg)
	}
}

func (self *Room) Bye() {
	self.clear(true)
}

func (self *Room) KickPerson(uid int64, result int) {
	for i := 0; i < len(self.Uid); i++ {
		if self.Uid[i] == uid {
			robot := self.Robot[i]
			copy(self.Uid[i:], self.Uid[i+1:])
			self.Uid = self.Uid[:len(self.Uid)-1]
			copy(self.Name[i:], self.Name[i+1:])
			self.Name = self.Name[:len(self.Name)-1]
			copy(self.ImgUrl[i:], self.ImgUrl[i+1:])
			self.ImgUrl = self.ImgUrl[:len(self.ImgUrl)-1]
			copy(self.Sex[i:], self.Sex[i+1:])
			self.Sex = self.Sex[:len(self.Sex)-1]
			copy(self.Param[i:], self.Param[i+1:])
			self.Param = self.Param[:len(self.Param)-1]
			copy(self.Robot[i:], self.Robot[i+1:])
			self.Robot = self.Robot[:len(self.Robot)-1]

			self.broadCastMsg("roominfo", self.getRoomMsg())
			self.game.OnExit(uid)

			if robot { //! 是机器人直接跳出
				lib.GetRobotMgr().BackRobot(uid, false)
				var msg staticfunc.Msg_DelRobot
				msg.Id = GetServer().Con.Id
				msg.Room = self.Id
				msg.Uid = uid
				msg.GameType = self.Type
				GetServer().CallLogin("ServerMethod.ServerMsg", "delrobot", &msg)
				break
			}

			var _msg staticfunc.Msg_JoinFail
			_msg.Uid = uid
			_msg.Id = GetServer().Con.Id
			_msg.Room = self.Id
			_msg.GameType = self.Type
			GetServer().CallLogin("ServerMethod.ServerMsg", "joinfail", &_msg)

			person := GetPersonMgr().GetPerson(uid)
			if person != nil {
				if result == 99 {
					person.SendErr("金币不足，请选择其他场次")
				} else if result == 98 {
					person.SendErr("由于您长时间未准备，故被请出金币场，请重新选择")
				} else if result == 97 {
					person.SendErr("房间人数已满")
				} else if result == 96 {
					person.SendErr("由于您长时间未下注，故被请出金币场，请重新选择")
				} else if result == 95 {
					person.SendErr("您已经被封号")
				} else if result == 94 {
					person.SendErr("由于您长时间未操作，故被请出金币场，请重新选择")
				}

				var msg Msg_ExitRoom
				msg.Result = result
				person.SendMsg("exitroom", &msg)
				person.CloseSession()
				GetPersonMgr().DelPerson(uid)
			}
			break
		}
	}

	if GetServer().Con.Type >= 3 { //! 金币场
		if len(self.Uid) == 0 {
			self.clear(true)
			return
		}
	}

	self.flush()
}

//! 把观战的全部清理
func (self *Room) KickView() {
	for i := 0; i < len(self.Viewer); {
		uid := self.Viewer[i]

		copy(self.Viewer[i:], self.Viewer[i+1:])
		self.Viewer = self.Viewer[:len(self.Viewer)-1]

		var _msg staticfunc.Msg_JoinFail
		_msg.Uid = uid
		_msg.Id = GetServer().Con.Id
		_msg.Room = self.Id
		_msg.GameType = self.Type
		GetServer().CallLogin("ServerMethod.ServerMsg", "joinfail", &_msg)

		person := GetPersonMgr().GetPerson(uid)
		if person != nil {
			var msg Msg_ExitRoom
			person.SendMsg("exitroom", &msg)
			person.CloseSession()
			GetPersonMgr().DelPerson(uid)
		}

		self.game.OnExit(uid)
	}

	self.flush()
}

//! 把观战的踢掉
func (self *Room) KickViewByUid(uid int64, result int) {
	for i := 0; i < len(self.Viewer); i++ {
		if self.Viewer[i] == uid {
			copy(self.Viewer[i:], self.Viewer[i+1:])
			self.Viewer = self.Viewer[:len(self.Viewer)-1]

			self.game.OnExit(uid)

			var _msg staticfunc.Msg_JoinFail
			_msg.Uid = uid
			_msg.Id = GetServer().Con.Id
			_msg.Room = self.Id
			_msg.GameType = self.Type
			GetServer().CallLogin("ServerMethod.ServerMsg", "joinfail", &_msg)

			person := GetPersonMgr().GetPerson(uid)
			if person != nil {
				if result == 96 {
					person.SendErr("由于您长时间未下注，故被请出金币场，请重新选择")
				} else if result == 97 {
					person.SendErr("由于您长时间未操作，故被请出金币场，请重新选择")
				}

				var msg Msg_ExitRoom
				person.SendMsg("exitroom", &msg)
				person.CloseSession()
				GetPersonMgr().DelPerson(uid)
			}

			if len(self.Viewer) == 0 && len(self.Uid) == 0 && self.Many != 1 { //! 没有观众也没有玩家
				lib.GetLogMgr().Output(lib.LOG_DEBUG, "清理房间")
				self.clear(true)
				return
			}

			break
		}
	}

	self.flush()
}

//! 30分钟没有人操作，这个房间为不活动房间
func (self *Room) IsLive() bool {
	if GetServer().Con.Type >= 3 { //! 金币场给2天活跃时间
		return time.Now().Unix()-self.LiveTime < 86400*2
	}

	if self.Step == 0 && !self.Begin { //! 一局都未开始的房间，10分钟无人操作就解散
		return time.Now().Unix()-self.LiveTime < 600
	}
	return time.Now().Unix()-self.LiveTime < 1800
}

func (self *Room) dissmissThread() {
	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		if self.Time == 0 {
			break
		}

		if time.Now().Unix() >= self.Time {
			self.Operator(NewRoomMsg("byeroom", "now", 0, nil))
			break
		}
	}

	//！ 关掉定时器
	ticker.Stop()
}

func (self *Room) byeThread() {
	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		if self.ByeTime == 0 {
			break
		}

		if time.Now().Unix() >= self.ByeTime {
			self.Operator(NewRoomMsg("delroom", "", 0, nil))
			break
		}
	}

	//！ 关掉定时器
	ticker.Stop()
}

func (self *Room) GetLogRoom() string {
	var _log staticfunc.Log_Room
	_log.Roomid = self.Id
	_log.Uid = self.Uid
	return lib.HF_JtoA(&_log)
}

func (self *Room) GetName(uid int64) string {
	for i := 0; i < len(self.Uid); i++ {
		if self.Uid[i] == uid {
			return self.Name[i]
		}
	}

	return ""
}

func (self *Room) GetHead(uid int64) string {
	for i := 0; i < len(self.Uid); i++ {
		if self.Uid[i] == uid {
			return self.ImgUrl[i]
		}
	}

	return ""
}

//! 发送俱乐部战绩
func (self *Room) ClubResult(info []staticfunc.JS_CreateRoomMem) {
	if self.ClubId == 0 {
		return
	}

	var msg staticfunc.Msg_ClubRoomResult
	msg.ClubId = self.ClubId
	msg.GameType = self.Type
	msg.Param1 = self.Param1
	msg.Param2 = self.Param2
	msg.RoomId = self.Id
	msg.MaxStep = self.MaxStep
	msg.Num = self.Num
	for i := 0; i < len(info); i++ {
		info[i].Name = self.GetName(info[i].Uid)
		info[i].Head = self.GetHead(info[i].Uid)
	}
	msg.Info = info
	GetServer().CallLogin("ServerMethod.ServerMsg", "clubroomresult", &msg)
}

//! 同步代开分数
func (self *Room) AgentResult(info []staticfunc.JS_CreateRoomMem) {
	if !self.Agent { //! 不是代开不用同步
		return
	}

	if self.ClubId != 0 { //! 俱乐部房间不用同步
		return
	}

	var msg staticfunc.Msg_ViewRoomScore
	msg.Host = self.Host
	msg.RoomId = self.Id
	msg.Node = info
	GetServer().CallLogin("ServerMethod.ServerMsg", "viewroomscore", &msg)
}

///////////////////////////////////////////////////////////////////////////////////
//! 房间管理者
type RoomMgr struct {
	MapRoom map[int]*Room

	lock *sync.RWMutex
}

var roommgrSingleton *RoomMgr = nil

//! 得到服务器指针
func GetRoomMgr() *RoomMgr {
	if roommgrSingleton == nil {
		roommgrSingleton = new(RoomMgr)
		roommgrSingleton.MapRoom = make(map[int]*Room)
		roommgrSingleton.lock = new(sync.RWMutex)
	}

	return roommgrSingleton
}

//! 加一个房间
func (self *RoomMgr) CreateRoom(id int, _type int, num int, param1 int, param2 int, agent int64, clubid int64) bool {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][_type]
	if !ok {
		return false
	}

	if num <= 0 || num > lib.HF_Atoi(csv["maxstep"]) {
		return false
	}

	self.lock.RLock()
	room, ok := self.MapRoom[id]
	self.lock.RUnlock()
	if ok {
		return false
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	//! 从redis获取数据，或者初始化
	data := GetServer().DB_GetData("roominfo", int64(id))
	if string(data) == "" { //! 没有房间就创建
		room = NewRoom(id, _type, param1, param2, agent, clubid)
		room.Init(num)
		self.MapRoom[id] = room
		return true
	}

	room = new(Room)
	json.Unmarshal(data, room)
	if room.Time != 0 && time.Now().Unix() > room.Time { //! 从redis获取的房间已经解散了
		var msg staticfunc.Msg_DelRoom
		msg.Id = GetServer().Con.Id
		msg.RoomId = id
		msg.Uid = room.Uid
		msg.Host = room.Host
		msg.Agent = room.Agent
		msg.ClubId = room.ClubId
		msg.GameType = room.Type
		GetServer().CallLogin("ServerMethod.ServerMsg", "delroom", &msg)

		room = NewRoom(id, _type, param1, param2, agent, clubid)
		room.Init(num)
		self.MapRoom[id] = room
		return true
	}

	if room.ByeTime != 0 && time.Now().Unix() > room.ByeTime { //! 从redis获取的房间已经解散了
		var msg staticfunc.Msg_DelRoom
		msg.Id = GetServer().Con.Id
		msg.RoomId = id
		msg.Uid = room.Uid
		msg.Host = room.Host
		msg.Agent = room.Agent
		msg.ClubId = room.ClubId
		msg.GameType = room.Type
		GetServer().CallLogin("ServerMethod.ServerMsg", "delroom", &msg)

		room = NewRoom(id, _type, param1, param2, agent, clubid)
		room.Init(num)
		self.MapRoom[id] = room
		return true
	}

	if !room.IsLive() { //! 从redis获取的房间已经解散了
		var msg staticfunc.Msg_DelRoom
		msg.Id = GetServer().Con.Id
		msg.RoomId = id
		msg.Uid = room.Uid
		msg.Host = room.Host
		msg.Agent = room.Agent
		msg.ClubId = room.ClubId
		msg.GameType = room.Type
		GetServer().CallLogin("ServerMethod.ServerMsg", "delroom", &msg)

		room = NewRoom(id, _type, param1, param2, agent, clubid)
		room.Init(num)
		self.MapRoom[id] = room
		return true
	}

	room.Init(num)
	self.MapRoom[id] = room

	return false
}

func (self *RoomMgr) JoinRoom(id int) (bool, *Room) {
	self.lock.RLock()
	room, ok := self.MapRoom[id]
	self.lock.RUnlock()
	if ok { //! 房间存在
		if room.ByeTime != 0 { //! 将要解散的房间
			return false, nil
		}
		return true, room
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	data := GetServer().DB_GetData("roominfo", int64(id))
	if string(data) == "" { //! redis里面没有房间
		var msg staticfunc.Msg_DelRoom
		msg.Id = GetServer().Con.Id
		msg.RoomId = id
		GetServer().CallLogin("ServerMethod.ServerMsg", "delroom", &msg)
		return false, nil
	}

	room = new(Room)
	json.Unmarshal(data, room)

	if room.Time != 0 && time.Now().Unix() > room.Time { //! 从redis获取的房间已经解散了
		var msg staticfunc.Msg_DelRoom
		msg.Id = GetServer().Con.Id
		msg.RoomId = id
		msg.Uid = room.Uid
		msg.Host = room.Host
		msg.Agent = room.Agent
		msg.ClubId = room.ClubId
		msg.GameType = room.Type
		GetServer().CallLogin("ServerMethod.ServerMsg", "delroom", &msg)

		GetServer().DB_SetData("roominfo", int64(id), []byte(""))
		return false, nil
	}

	if room.ByeTime != 0 && time.Now().Unix() > room.ByeTime { //! 从redis获取的房间已经解散了
		var msg staticfunc.Msg_DelRoom
		msg.Id = GetServer().Con.Id
		msg.RoomId = id
		msg.Uid = room.Uid
		msg.Host = room.Host
		msg.Agent = room.Agent
		msg.ClubId = room.ClubId
		msg.GameType = room.Type
		GetServer().CallLogin("ServerMethod.ServerMsg", "delroom", &msg)

		GetServer().DB_SetData("roominfo", int64(id), []byte(""))
		return false, nil
	}

	if !room.IsLive() { //! 从redis获取的房间已经解散了
		var msg staticfunc.Msg_DelRoom
		msg.Id = GetServer().Con.Id
		msg.RoomId = id
		msg.Uid = room.Uid
		msg.Host = room.Host
		msg.Agent = room.Agent
		msg.ClubId = room.ClubId
		msg.GameType = room.Type
		GetServer().CallLogin("ServerMethod.ServerMsg", "delroom", &msg)

		GetServer().DB_SetData("roominfo", int64(id), []byte(""))
		return false, nil
	}

	//! 成功创建房间
	room.Init(0)

	self.MapRoom[id] = room

	return true, room
}

//! 删一个房间
func (self *RoomMgr) DelRoom(id int) {
	self.lock.Lock()
	defer self.lock.Unlock()

	GetServer().DB_SetData("roominfo", int64(id), []byte(""))
	delete(self.MapRoom, id)
}

//! 得到房间
func (self *RoomMgr) GetRoom(id int) *Room {
	self.lock.RLock()
	defer self.lock.RUnlock()

	room, ok := self.MapRoom[id]
	if !ok {
		return nil
	}

	return room
}

//! 扣房卡
func (self *RoomMgr) CostCard(uid int64, _type, num int, room *Room) bool {
	var msg staticfunc.Msg_CostCard
	msg.Uid = uid
	msg.Num = num
	msg.Type = _type
	if room != nil {
		msg.Dec = room.Type
	}
	result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "costcard", &msg)

	if err == nil && string(result) != "" { //! 插入消耗房卡的日志
		tmp := strings.Split(string(result), "_")
		if len(tmp) == 4 {
			card := lib.HF_Atoi(tmp[0])
			gold := lib.HF_Atoi(tmp[1])
			if room != nil {
				if card > 0 {
					GetServer().InsertLog(room.Type, uid, staticfunc.COST_CARD, card, room.GetLogRoom())
				}
				if gold > 0 {
					//GetServer().InsertLog(room.Type, uid, staticfunc.COST_GOLD, gold, room.GetLogRoom())
				}
			}
			nowcard := lib.HF_Atoi(tmp[2])
			nowgold := lib.HF_Atoi(tmp[3])
			person := GetPersonMgr().GetPerson(uid)
			if person != nil {
				person.Card = nowcard
				person.Gold = nowgold

				var msg staticfunc.S2C_UpdCard
				msg.Card = nowcard + nowgold
				msg.Gold = nowgold
				person.SendMsg("updcard", &msg)
			}
		}
		return true
	}

	return false
}

//! 直接同步到loginserver去
func (self *RoomMgr) AddCard2(uid int64, _type, num int, dec int) bool {
	var msg staticfunc.Msg_GiveCard
	msg.Uid = uid
	msg.Num = num
	msg.Pid = _type
	msg.Ip = ""
	msg.Dec = dec
	msg.Sync = true
	GetServer().CallLogin("ServerMethod.ServerMsg", "givecard", &msg)

	return true
}

//! 加房卡
func (self *RoomMgr) AddCard(uid int64, _type, num int, dec int) bool {
	var msg staticfunc.Msg_GiveCard
	msg.Uid = uid
	msg.Num = num
	msg.Pid = _type
	msg.Ip = ""
	msg.Dec = dec
	result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "givecard", &msg)

	if err != nil || string(result) == "false" {
		return false
	}

	tmp := strings.Split(string(result), "_")
	if len(tmp) != 2 {
		return false
	}

	nowcard := lib.HF_Atoi(tmp[0])
	nowgold := lib.HF_Atoi(tmp[1])

	person := GetPersonMgr().GetPerson(uid)
	if person != nil {
		person.Card = nowcard
		person.Gold = nowgold
		var msg staticfunc.S2C_UpdCard
		msg.Card = nowcard + nowgold
		msg.Gold = nowgold
		person.SendMsg("updcard", &msg)
	}

	return true
}
