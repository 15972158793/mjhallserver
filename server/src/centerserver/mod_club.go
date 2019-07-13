package centerserver

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"lib"
	"staticfunc"
	"strings"
	"sync"
	"time"
)

type JS_Club struct {
	ClubId []int64 `json:"id"`    //! 所在俱乐部的id
	Apply  []int64 `json:"apply"` //! 申请的俱乐部id
}

//! 俱乐部结构
type SQL_Club struct {
	Uid  int64
	Info []byte

	info JS_Club
}

func (self *SQL_Club) Decode() {
	json.Unmarshal(self.Info, &self.info)
}

func (self *SQL_Club) Encode() {
	self.Info = lib.HF_JtoB(self.info)
}

//! 俱乐部模块
type Mod_Club struct {
	person *Person
	club   SQL_Club
	lock   *sync.RWMutex
}

func (self *Mod_Club) OnGetData(person *Person) {
	self.person = person
	self.lock = new(sync.RWMutex)
}

func (self *Mod_Club) OnGetOtherData() {
	c := GetServer().Redis.Get()
	defer c.Close()
	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("myclub_%d", self.person.Uid)))
	if err == nil {
		self.club.Info = v
		self.club.Decode()
	} else {
		sql := fmt.Sprintf("select * from `club` where uid = %d", self.person.Uid)
		GetServer().DB.GetOneData(sql, &self.club)
		if self.club.Uid <= 0 {
			self.club.Uid = self.person.Uid
			self.club.info.ClubId = make([]int64, 0)
			self.club.info.Apply = make([]int64, 0)
			self.club.Encode()
			sql := fmt.Sprintf("insert into `%s`(`uid`, `info`) values (%d, ?)", "club", self.person.Uid)
			GetServer().SqlQueue(sql, self.club.Info, true)

			c := GetServer().Redis.Get()
			defer c.Close()
			c.Do("SET", fmt.Sprintf("myclub_%d", self.person.Uid), self.club.Info)
			c.Do("EXPIRE", fmt.Sprintf("myclub_%d", self.person.Uid), 86400*7)
		} else {
			self.club.Decode()
		}
	}
}

func (self *Mod_Club) OnMsg(head string, body []byte) bool {
	switch head {
	case "getclublist": //! 得到俱乐部列表
		GetClubMgr().GetClubList(self.person)
	case "clubapply": //! 申请俱乐部
		var msg Msg_ClubID
		json.Unmarshal(body, &msg)
		self.Apply(msg.ClubId)
	case "clubunapply":
		var msg Msg_ClubID
		json.Unmarshal(body, &msg)
		self.UnApply(msg.ClubId)
	case "clubcreate": //!  创建俱乐部
		var msg C2S_ClubCreate
		json.Unmarshal(body, &msg)
		self.Create(msg.Name, msg.Icon, msg.Notice, msg.Mode)
	case "clubenter": //! 进入俱乐部
		var msg Msg_ClubID
		json.Unmarshal(body, &msg)
		self.Get(msg.ClubId)
	case "clubleave": //! 离开俱乐部
		var msg Msg_ClubID
		json.Unmarshal(body, &msg)
		self.Exit(msg.ClubId)
	case "clubexit": //! 俱乐部踢人
		var msg C2S_ClubLeave
		json.Unmarshal(body, &msg)
		self.Leave(msg.ClubId, msg.Uid)
	case "cluborder": //! 俱乐部处理请求
		var msg C2S_ClubOrder
		json.Unmarshal(body, &msg)
		self.Order(msg.ClubId, msg.Uid, msg.Agree)
	case "clubroom": //! 俱乐部开房
		var msg C2S_ClubCreateRoom
		json.Unmarshal(body, &msg)
		self.CreateRoom(msg.ClubId, msg.GameType, msg.Num, msg.Param1, msg.Param2)
	case "clubroomchat": //! 俱乐部开房列表
		var msg Msg_ClubID
		json.Unmarshal(body, &msg)
		self.GetRoomChat(msg.ClubId)
	case "clubroomresult": //! 俱乐部战绩
		var msg Msg_ClubID
		json.Unmarshal(body, &msg)
		self.GetRoomResult(msg.ClubId)
	case "clubmode": //! 设置俱乐部开房模式
		var msg C2S_ClubMode
		json.Unmarshal(body, &msg)
		self.SetMode(msg.ClubId, msg.Mode)
	case "clubgame": //! 设置俱乐部游戏
		var msg C2S_ClubGame
		json.Unmarshal(body, &msg)
		self.SetGame(msg.ClubId, msg.Game)
	case "clubname": //! 俱乐部名字
		var msg C2S_ClubName
		json.Unmarshal(body, &msg)
		self.SetName(msg.ClubId, msg.Name)
	case "clubapplylist": //! 得到申请列表
		var msg Msg_ClubID
		json.Unmarshal(body, &msg)
		self.GetApplyList(msg.ClubId)
	case "clubeventlist": //! 得到事件列表
		var msg Msg_ClubID
		json.Unmarshal(body, &msg)
		self.GetEventList(msg.ClubId)
	case "clubnotice": //! 修改公告
		var msg C2S_ClubNotice
		json.Unmarshal(body, &msg)
		self.SetNotice(msg.ClubId, msg.Notice)
	case "clubicon": //! 修改头像
		var msg C2S_ClubIcon
		json.Unmarshal(body, &msg)
		self.SetIcon(msg.ClubId, msg.Icon)
	case "clubcostcard": //! 俱乐部消耗房卡统计
		var msg Msg_ClubID
		json.Unmarshal(body, &msg)
		self.GetCostCard(msg.ClubId)
	}
	return false
}

func (self *Mod_Club) OnSave(sql bool) {
	self.club.Encode()
	c := GetServer().Redis.Get()
	defer c.Close()
	c.Do("SET", fmt.Sprintf("myclub_%d", self.person.Uid), self.club.Info)
	c.Do("EXPIRE", fmt.Sprintf("myclub_%d", self.person.Uid), 86400*7)

	GetServer().SqlQueue(fmt.Sprintf("update `%s` set `info` = ? where `uid` = '%d'", "club", self.person.Uid), self.club.Info, true)
}

//! 发送消息
func (self *Mod_Club) SendInfo() {
	chg := false
	var msg Msg_MyClubInfo
	for i := 0; i < len(self.club.info.ClubId); {
		club := GetClubMgr().GetClub(self.club.info.ClubId[i])
		if club == nil {
			copy(self.club.info.ClubId[i:], self.club.info.ClubId[i+1:])
			self.club.info.ClubId = self.club.info.ClubId[:len(self.club.info.ClubId)-1]
			chg = true
		} else {
			club.UpdMem(self.person.Uid, self.person.mod_info.Name, self.person.mod_info.Imgurl, false)
			msg.Info = append(msg.Info, Son_ClubList{club.Id, club.info.Name, club.info.Icon, club.info.ExNotice, len(club.info.Mem), club.info.Game})
			i++
		}
	}
	msg.Apply = self.club.info.Apply
	self.person.SendMsg(lib.HF_EncodeMsg("myclubinfo", &msg, true))

	if chg {
		self.OnSave(false)
	}
}

//! 得到拥有俱乐部和申请俱乐部的数量
func (self *Mod_Club) GetInfoNum() (int, int) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	return len(self.club.info.ClubId), len(self.club.info.Apply)
}

//! 是否加入俱乐部
func (self *Mod_Club) IsJoin(clubid int64) bool {
	self.lock.RLock()
	defer self.lock.RUnlock()

	for i := 0; i < len(self.club.info.ClubId); i++ {
		if self.club.info.ClubId[i] == clubid {
			return true
		}
	}

	return false
}

//! 加入
func (self *Mod_Club) AddJoin(clubid int64) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.club.info.ClubId = append(self.club.info.ClubId, clubid)
	self.OnSave(false)
}

//! 是否申请该俱乐部
func (self *Mod_Club) IsApply(clubid int64) bool {
	self.lock.RLock()
	defer self.lock.RUnlock()

	for i := 0; i < len(self.club.info.Apply); i++ {
		if self.club.info.Apply[i] == clubid {
			return true
		}
	}

	return false
}

//! 加入申请
func (self *Mod_Club) AddApply(clubid int64) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.club.info.Apply = append(self.club.info.Apply, clubid)
	self.OnSave(false)
}

//! 撤销申请
func (self *Mod_Club) DelApply(clubid int64) {
	self.lock.Lock()
	defer self.lock.Unlock()

	for i := 0; i < len(self.club.info.Apply); i++ {
		if self.club.info.Apply[i] == clubid {
			copy(self.club.info.Apply[i:], self.club.info.Apply[i+1:])
			self.club.info.Apply = self.club.info.Apply[:len(self.club.info.Apply)-1]
			break
		}
	}
	self.OnSave(false)
}

//! 加入俱乐部
func (self *Mod_Club) AddClub(clubid int64) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.club.info.ClubId = append(self.club.info.ClubId, clubid)
	self.OnSave(false)
}

//! 退出俱乐部
func (self *Mod_Club) DelClub(clubid int64) {
	self.lock.Lock()
	defer self.lock.Unlock()

	for i := 0; i < len(self.club.info.ClubId); i++ {
		if self.club.info.ClubId[i] == clubid {
			copy(self.club.info.ClubId[i:], self.club.info.ClubId[i+1:])
			self.club.info.ClubId = self.club.info.ClubId[:len(self.club.info.ClubId)-1]
			break
		}
	}
	self.OnSave(false)
}

//! 申请俱乐部
func (self *Mod_Club) Apply(clubid int64) {
	if self.IsJoin(clubid) {
		self.person.SendErr("您已加入了该俱乐部")
		return
	}

	if self.IsApply(clubid) {
		self.person.SendErr("您已申请了该俱乐部")
		return
	}

	club, apply := self.GetInfoNum()
	if club+apply >= 5 {
		self.person.SendErr("您所申请或者加入的俱乐部超过上限")
		return
	}

	result := GetClubMgr().ApplyClub(self.person, clubid)
	if result == 1 {
		self.person.SendErr("您申请的俱乐部不存在")
		return
	}

	if result == 2 {
		self.person.SendErr("该俱乐部申请列表已满")
		return
	}

	self.AddApply(clubid)

	var msg Msg_ClubID
	msg.ClubId = clubid
	self.person.SendMsg(lib.HF_EncodeMsg("clubapply", &msg, true))
}

//! 取消申请
func (self *Mod_Club) UnApply(clubid int64) {
	if !self.IsApply(clubid) {
		self.person.SendErr("您没有申请该俱乐部")
		return
	}

	if !GetClubMgr().UnApplyClub(self.person, clubid) {
		self.person.SendErr("您申请的俱乐部不存在")
		return
	}

	self.DelApply(clubid)
	var msg Msg_ClubID
	msg.ClubId = clubid
	self.person.SendMsg(lib.HF_EncodeMsg("clubunapply", &msg, true))
}

//! 进入俱乐部
func (self *Mod_Club) Get(clubid int64) {
	club := GetClubMgr().GetClub(clubid)
	if club == nil {
		self.person.SendErr("俱乐部不存在")
		return
	}

	club.UpdMem(self.person.Uid, self.person.mod_info.Name, self.person.mod_info.Imgurl, true)

	var msg Msg_ClubInfo
	msg.Id = club.Id
	msg.Name = club.info.Name
	msg.Icon = club.info.Icon
	msg.Mode = club.info.Mode
	msg.InNotice = club.info.InNotice
	msg.ExNotice = club.info.ExNotice
	msg.Host = club.info.Host
	msg.Member = club.info.Mem
	msg.Game = club.info.Game
	if len(club.info.Apply) > 0 {
		msg.Red = 1
	} else {
		msg.Red = 0
	}
	if len(club.info.Event) > 0 {
		msg.MsgRed = 1
	} else {
		msg.MsgRed = 0
	}

	self.person.SendMsg(lib.HF_EncodeMsg("clubinfo", &msg, true))
}

//! 进入离开
func (self *Mod_Club) Exit(clubid int64) {
	club := GetClubMgr().GetClub(clubid)
	if club == nil {
		self.person.SendErr("俱乐部不存在")
		return
	}

	club.UpdMem(self.person.Uid, self.person.mod_info.Name, self.person.mod_info.Imgurl, false)
}

//! 创建俱乐部
func (self *Mod_Club) Create(name string, icon string, notice string, mode int) {
	clubnum, applynum := self.GetInfoNum()
	if clubnum+applynum >= 5 {
		self.person.SendErr("您所申请或者加入的俱乐部超过上限")
		return
	}

	card, _, _ := GetServer().GetCard(self.person.Uid)
	if card < 100 {
		self.person.SendErr("您必须拥有100张房卡或更多才能创建俱乐部")
		return
	}

	club := GetClubMgr().Create(self.person, name, icon, notice, mode)
	if club == nil {
		return
	}

	self.AddJoin(club.Id)

	var msg Msg_ClubInfo
	msg.Id = club.Id
	msg.Name = club.info.Name
	msg.Icon = club.info.Icon
	msg.Mode = club.info.Mode
	msg.InNotice = club.info.InNotice
	msg.ExNotice = club.info.ExNotice
	msg.Host = club.info.Host
	msg.Member = club.info.Mem
	msg.Game = club.info.Game
	msg.Red = 0
	msg.MsgRed = 0

	self.person.SendMsg(lib.HF_EncodeMsg("clubcreate", &msg, true))
}

//! 处理
func (self *Mod_Club) Order(clubid int64, uid int64, agree bool) {
	club := GetClubMgr().GetClub(clubid)
	if club == nil {
		self.person.SendErr("俱乐部不存在")
		return
	}

	if club.info.Host != self.person.Uid {
		self.person.SendErr("非群主无法操作")
		return
	}

	if !club.OrderApply(uid, agree) {
		self.person.SendErr("处理失败,请刷新")
		return
	}

	person := GetPersonMgr().GetPerson(uid, true)
	if person != nil {
		person.GetModule("club").(*Mod_Club).DelApply(clubid)
		if agree {
			person.GetModule("club").(*Mod_Club).AddClub(clubid)
		}
		person.GetModule("club").(*Mod_Club).SendInfo()
	}

	var msg Msg_ClubMem
	msg.Id = club.Id
	msg.Member = club.info.Mem
	msg.Apply = club.info.Apply

	self.person.SendMsg(lib.HF_EncodeMsg("clubmember", &msg, true))
}

//! 设置模式
func (self *Mod_Club) SetMode(clubid int64, mode int) {
	club := GetClubMgr().GetClub(clubid)
	if club == nil {
		self.person.SendErr("俱乐部不存在")
		return
	}

	if club.info.Host != self.person.Uid {
		self.person.SendErr("非群主无法更改模式")
		return
	}

	club.info.Mode = mode

	club.chg.Chg = true

	var msg S2C_ClubMode
	msg.Id = clubid
	msg.Mode = mode
	self.person.SendMsg(lib.HF_EncodeMsg("clubmode", &msg, true))
}

//! 设置俱乐部游戏
func (self *Mod_Club) SetGame(clubid int64, game []JS_ClubGame) {
	club := GetClubMgr().GetClub(clubid)
	if club == nil {
		self.person.SendErr("俱乐部不存在")
		return
	}

	if club.info.Host != self.person.Uid {
		self.person.SendErr("非群主无法设置")
		return
	}

	club.info.Game = game

	club.chg.Chg = true

	var msg S2C_ClubGame
	msg.Id = clubid
	msg.Game = game
	self.person.SendMsg(lib.HF_EncodeMsg("clubgame", &msg, true))
}

//! 修改名字
func (self *Mod_Club) SetName(clubid int64, name string) {
	club := GetClubMgr().GetClub(clubid)
	if club == nil {
		self.person.SendErr("俱乐部不存在")
		return
	}

	if club.info.Host != self.person.Uid {
		self.person.SendErr("非群主无法修改名字")
		return
	}

	if !club.SetName(name) {
		self.person.SendErr("该名字已被占用")
		return
	}

	var msg S2C_ClubName
	msg.Id = clubid
	msg.Name = name
	self.person.SendMsg(lib.HF_EncodeMsg("clubname", &msg, true))
}

//! 踢人
func (self *Mod_Club) Leave(clubid int64, uid int64) {
	club := GetClubMgr().GetClub(clubid)
	if club == nil {
		self.person.SendErr("俱乐部不存在")
		return
	}

	if uid > 0 { //! 踢人
		if club.info.Host != self.person.Uid {
			self.person.SendErr("非群主无法踢人")
			return
		}
	} else {
		uid = self.person.Uid
	}

	if uid == club.info.Host { //! 俱乐部主席离开则解散
		var msg Msg_ClubID
		msg.ClubId = clubid
		club.BroadCastMsg(lib.HF_EncodeMsg("clubdissmiss", &msg, true), true)

		GetClubMgr().Dissmiss(clubid)
	} else {
		if !club.Leave(uid) {
			self.person.SendErr("该用户已离开俱乐部")
			return
		}

		var msg S2C_ClubLeave
		msg.Id = clubid
		msg.Uid = uid

		person := GetPersonMgr().GetPerson(uid, true)
		if person != nil {
			person.GetModule("club").(*Mod_Club).DelClub(clubid)
			person.SendMsg(lib.HF_EncodeMsg("clubexit", &msg, true))
		}

		if uid != self.person.Uid {
			self.person.SendMsg(lib.HF_EncodeMsg("clubexit", &msg, true))
		}
	}
}

//! 开房
func (self *Mod_Club) CreateRoom(clubid int64, gametype int, num int, param1 int, param2 int) {
	club := GetClubMgr().GetClub(clubid)
	if club == nil {
		self.person.SendErr("俱乐部不存在")
		return
	}

	if !club.HasMem(self.person.Uid) {
		self.person.SendErr("你不是该俱乐部成员")
		return
	}

	if club.info.Mode == 0 && club.info.Host != self.person.Uid {
		self.person.SendErr("该俱乐部设置为群主才能开设房间")
		return
	}

	roomid := 0
	maxstep := 0
	card := 0
	gold := 0
	//! 向loginserver请求创建房间
	{
		var msg staticfunc.Msg_ClubCreateRoom
		msg.ClubId = clubid
		msg.GameType = gametype
		msg.Num = num
		msg.Param1 = param1
		msg.Param2 = param2
		msg.Host = club.info.Host
		msg.Uid = self.person.Uid
		msg.IP = self.person.session.IP
		result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "createroom", &msg)
		if err != nil || string(result) == "" || string(result) == "false" {
			self.person.SendErr("开设房间错误")
			return
		}
		if string(result) == "nocard" {
			self.person.SendErr("群主房卡不足")
			return
		}
		tmp := strings.Split(string(result), "_")
		if len(tmp) < 5 {
			self.person.SendErr("开设房间错误")
			return
		}

		roomid = lib.HF_Atoi(tmp[0])
		maxstep = lib.HF_Atoi(tmp[2])
		card = lib.HF_Atoi(tmp[3])
		gold = lib.HF_Atoi(tmp[4])
	}

	//! 扣除主席房卡
	hostperson := GetPersonMgr().GetPerson(club.info.Host, false)
	if hostperson != nil {
		hostperson.UpdCard(card, gold)
	}

	//! 发送聊天
	{
		var node JS_ClubRoomChat
		node.RoomId = roomid
		node.Param1 = param1
		node.Param2 = param2
		node.GameType = gametype
		node.MaxStep = maxstep
		node.Time = time.Now().Unix()
		node.Uid = self.person.Uid
		node.Name = self.person.mod_info.Name
		node.Head = self.person.mod_info.Imgurl
		node.State = 0
		node.Num = num
		club.AddRoomChat(node)

		club.chg.Chg = true

		var msg Msg_ClubAddRoomChat
		msg.Info = node
		club.BroadCastMsg(lib.HF_EncodeMsg("addroomchat", &msg, true), false)
	}
}

//! 得到房间开设
func (self *Mod_Club) GetRoomChat(clubid int64) {
	club := GetClubMgr().GetClub(clubid)
	if club == nil {
		self.person.SendErr("俱乐部不存在")
		return
	}

	var msg Msg_ClubRoomChat
	msg.Info = club.info.RoomChat
	self.person.SendMsg(lib.HF_EncodeMsg("clubroomchat", &msg, true))
}

//! 得到房间战绩
func (self *Mod_Club) GetRoomResult(clubid int64) {
	club := GetClubMgr().GetClub(clubid)
	if club == nil {
		self.person.SendErr("俱乐部不存在")
		return
	}

	var msg Msg_ClubRoomResult
	msg.Info = club.info.RoomResult
	self.person.SendMsg(lib.HF_EncodeMsg("clubroomresult", &msg, true))
}

//! 得到申请列表
func (self *Mod_Club) GetApplyList(clubid int64) {
	club := GetClubMgr().GetClub(clubid)
	if club == nil {
		self.person.SendErr("俱乐部不存在")
		return
	}

	var msg S2C_ClubApplyList
	msg.Id = clubid
	msg.Info = club.info.Apply
	self.person.SendMsg(lib.HF_EncodeMsg("clubapplylist", &msg, true))
}

//! 得到事件列表
func (self *Mod_Club) GetEventList(clubid int64) {
	club := GetClubMgr().GetClub(clubid)
	if club == nil {
		self.person.SendErr("俱乐部不存在")
		return
	}

	var msg S2C_ClubEventList
	msg.Id = clubid
	msg.Info = club.GetEvent()
	self.person.SendMsg(lib.HF_EncodeMsg("clubeventlist", &msg, true))
}

//! 修改公告
func (self *Mod_Club) SetNotice(clubid int64, notice string) {
	club := GetClubMgr().GetClub(clubid)
	if club == nil {
		self.person.SendErr("俱乐部不存在")
		return
	}

	if club.info.Host != self.person.Uid {
		self.person.SendErr("只有群主才能修改公告")
		return
	}

	club.info.ExNotice = notice
	club.chg.Chg = true

	self.person.SendErr("修改成功")
}

//! 修改头像
func (self *Mod_Club) SetIcon(clubid int64, icon string) {
	club := GetClubMgr().GetClub(clubid)
	if club == nil {
		self.person.SendErr("俱乐部不存在")
		return
	}

	if club.info.Host != self.person.Uid {
		self.person.SendErr("只有群主才能修改头像")
		return
	}

	club.info.Icon = icon
	club.chg.Chg = true

	self.person.SendErr("修改成功")
}

//! 得到房卡消耗统计
func (self *Mod_Club) GetCostCard(clubid int64) {
	club := GetClubMgr().GetClub(clubid)
	if club == nil {
		self.person.SendErr("俱乐部不存在")
		return
	}

	if club.info.Host != self.person.Uid {
		self.person.SendErr("只有群主才能查看房卡统计")
		return
	}

	var msg S2C_ClubCostCard
	msg.ClubId = clubid
	msg.Info = club.GetCostCard()
	self.person.SendMsg(lib.HF_EncodeMsg("clubcostcard", &msg, true))
}

//! 是否是有关联的俱乐部
func (self *Mod_Club) IsMyClub(clubid int64) bool {
	for i := 0; i < len(self.club.info.ClubId); i++ {
		if self.club.info.ClubId[i] == clubid {
			return true
		}
	}

	for i := 0; i < len(self.club.info.Apply); i++ {
		if self.club.info.Apply[i] == clubid {
			return true
		}
	}

	return false
}
