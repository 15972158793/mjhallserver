package centerserver

import (
	"encoding/json"
	"fmt"
	"lib"
	"net/http"
	"net/url"
	"staticfunc"
)

type HTMsg_AreaInfo struct {
	WChat string `json:"wchat"` //! 微信
	Phone string `json:"phone"` //! 电话
}

type HTMsg_GiveCard struct {
	Uid             int64 `json:"uid"`
	Productid       int   `json:"productid"`
	Repertory_count int   `json:"Repertor_count"`
}

//! 后台消息
func GiveCard(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	uid := lib.HF_Atoi(req.FormValue("uid"))
	productid := lib.HF_Atoi(req.FormValue("productid")) //! 1银卡card  2金卡gold
	product_count := lib.HF_Atoi(req.FormValue("product_count"))
	_type := lib.HF_Atoi(req.FormValue("type"))
	if _type == 0 {
		_type = -1
	}

	if uid <= 0 {
		w.Write([]byte("uid错误"))
		return
	}

	if productid != staticfunc.TYPE_CARD && productid != staticfunc.TYPE_GOLD {
		w.Write([]byte("productid错误"))
		return
	}

	if product_count <= 0 {
		w.Write([]byte("product_count错误"))
		return
	}

	ok, card, gold := false, 0, 0
	if productid == staticfunc.TYPE_CARD {
		ok, card, gold = GetServer().AddCard(int64(uid), product_count, lib.HF_GetHttpIP(req), _type)
	} else {
		ok, card, gold = GetServer().AddGold(int64(uid), product_count, lib.HF_GetHttpIP(req), _type)
	}

	if !ok {
		w.Write([]byte("发送失败"))
		return
	}

	person := GetPersonMgr().GetPerson(int64(uid), false)
	if person != nil {
		person.UpdCard(card, gold)
	}

	//! 回复给后台
	var _msg HTMsg_GiveCard
	_msg.Uid = int64(uid)
	_msg.Productid = productid
	if productid == staticfunc.TYPE_CARD {
		_msg.Repertory_count = card
	} else {
		_msg.Repertory_count = gold
	}
	w.Write(lib.HF_JtoB(&_msg))
}

func MoveCard(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	uid := lib.HF_Atoi(req.FormValue("uid"))
	productid := lib.HF_Atoi(req.FormValue("productid")) //! 1银卡  2金卡
	product_count := lib.HF_Atoi(req.FormValue("product_count"))
	_type := lib.HF_Atoi(req.FormValue("type"))
	if _type == 0 {
		_type = -1
	}

	if uid <= 0 {
		w.Write([]byte("uid错误"))
		return
	}

	if productid != staticfunc.TYPE_CARD && productid != staticfunc.TYPE_GOLD {
		w.Write([]byte("productid错误"))
		return
	}

	if product_count <= 0 {
		w.Write([]byte("product_count错误"))
		return
	}

	ok, card, gold := false, 0, 0
	if productid == staticfunc.TYPE_CARD {
		ok, card, gold = GetServer().CostCard(int64(uid), product_count, lib.HF_GetHttpIP(req), _type)
	} else {
		ok, card, gold = GetServer().CostGold(int64(uid), product_count, lib.HF_GetHttpIP(req), _type)
	}

	if !ok {
		w.Write([]byte("发送失败"))
		return
	}

	person := GetPersonMgr().GetPerson(int64(uid), false)
	if person != nil {
		person.UpdCard(card, gold)
	}

	//! 回复给后台
	var _msg HTMsg_GiveCard
	_msg.Uid = int64(uid)
	_msg.Productid = productid
	if productid == staticfunc.TYPE_CARD {
		_msg.Repertory_count = card
	} else {
		_msg.Repertory_count = gold
	}
	w.Write(lib.HF_JtoB(&_msg))
}

func FindPlayer(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	uid := lib.HF_Atoi(req.FormValue("uid"))

	if uid <= 0 {
		w.Write([]byte("uid错误"))
		return
	}

	var msg staticfunc.Msg_Uid
	msg.Uid = int64(uid)
	result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "findplayer", &msg)

	if err != nil || string(result) == "false" {
		w.Write([]byte("找不到玩家"))
		return
	}

	w.Write(result)
}

func FindPlayerFromOpenid(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	openid := req.FormValue("openid")

	if openid == "" {
		w.Write([]byte("openid错误"))
		return
	}

	var msg staticfunc.Msg_Openid
	msg.Openid = openid
	result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "findplayerfromopenid", &msg)

	if err != nil || string(result) == "false" {
		w.Write([]byte("找不到玩家"))
		return
	}

	w.Write(result)
}

func SendWarning(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	uid := lib.HF_Atoi(req.FormValue("uid"))
	context := req.FormValue("context")

	person := GetPersonMgr().GetPerson(int64(uid), false)
	if person == nil {
		w.Write([]byte("该玩家不在线"))
		return
	}

	person.SendWarning(context)
	w.Write([]byte("警告已发送"))
}

func IssueNotice(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	context := req.FormValue("context")
	group := lib.HF_Atoi(req.FormValue("group"))

	if group != 1 {
		group = 2
	}

	if GetServer().Notice[group-1] == context {
		w.Write([]byte("ok"))
		return
	}

	GetServer().Notice[group-1] = context
	c := GetServer().Redis.Get()
	defer c.Close()
	if group == 1 { //! 扑克
		c.Do("SET", "notice", context)
	} else {
		c.Do("SET", "notice_mj", context)
	}

	if context != "" {
		var msg S2C_Notice
		msg.Context = context
		GetPersonMgr().BroadCastMsg("", lib.HF_EncodeMsg("notice", &msg, true), group)
		w.Write([]byte("ok"))
	} else {
		w.Write([]byte("fail"))
	}
}

func IssueAreaNotice(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	area := req.FormValue("area")
	info := req.FormValue("info")
	group := lib.HF_Atoi(req.FormValue("group"))

	if group != 1 {
		group = 2
	}

	if info == "" {
		w.Write([]byte("fail"))
		return
	}

	GetServer().AreaNotice[group-1][area] = info

	c := GetServer().Redis.Get()
	defer c.Close()

	if group == 1 {
		c.Do("SET", "area_"+area, info)
	} else {
		c.Do("SET", "areamj_"+area, info)
	}

	var msg S2C_Notice
	msg.Context = info
	GetPersonMgr().BroadCastMsg(area, lib.HF_EncodeMsg("notice", &msg, true), group)
	w.Write([]byte("ok"))
}

func IssueAreaInfo(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	area := req.FormValue("area")
	info := req.FormValue("info")
	group := lib.HF_Atoi(req.FormValue("group"))

	if group != 1 {
		group = 2
	}

	if info == "" {
		w.Write([]byte("fail"))
		return
	}

	var ainfo HTMsg_AreaInfo
	err := json.Unmarshal([]byte(info), &ainfo)
	if err != nil {
		w.Write([]byte("fail"))
		return
	}
	if ainfo.WChat == "" {
		w.Write([]byte("wchat err"))
		return
	}
	GetServer().AreaInfo[group-1][area] = info

	c := GetServer().Redis.Get()
	defer c.Close()
	if group == 1 {
		c.Do("SET", "areainfo_"+area, info)
	} else {
		c.Do("SET", "areainfomj_"+area, info)
	}

	var msg S2C_Notice
	msg.Context = info
	GetPersonMgr().BroadCastMsg(area, lib.HF_EncodeMsg("areainfo", &msg, true), group)
	w.Write([]byte("ok"))
}

func IssueSysNotice(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	info := req.FormValue("info")
	text, err := url.QueryUnescape(info)
	if err != nil {
		w.Write([]byte("参数错误"))
		return
	}

	if GetNoticeMgr().InsertData(text) {
		w.Write([]byte("ok"))
	} else {
		w.Write([]byte("设置失败"))
	}
}

func GetOnlineNum(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	group := lib.HF_Atoi(req.FormValue("group"))

	if group != 1 {
		group = 2
	}

	w.Write([]byte(fmt.Sprintf("%d", GetPersonMgr().GetNum(group))))
}

func CloseDownPlayer(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	uid := lib.HF_Atoi(req.FormValue("uid"))
	openid := req.FormValue("openid")

	if uid > 0 {
		var msg staticfunc.Msg_Uid
		msg.Uid = int64(uid)
		result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "closedownplayer1", &msg)
		if err != nil || string(result) == "false" {
			w.Write([]byte("找不到玩家"))
			return
		}
		w.Write(result)
	} else if openid != "" {
		var msg staticfunc.Msg_Openid
		msg.Openid = openid
		result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "closedownplayer2", &msg)
		if err != nil || string(result) == "false" {
			w.Write([]byte("找不到玩家"))
			return
		}
		w.Write(result)
	} else {
		w.Write([]byte("参数错误"))
	}
}

func OpenDownPlayer(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	uid := lib.HF_Atoi(req.FormValue("uid"))
	openid := req.FormValue("openid")

	if uid > 0 {
		var msg staticfunc.Msg_Uid
		msg.Uid = int64(uid)
		result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "opendownplayer1", &msg)
		if err != nil || string(result) == "false" {
			w.Write([]byte("找不到玩家"))
			return
		}
		w.Write(result)
	} else if openid != "" {
		var msg staticfunc.Msg_Openid
		msg.Openid = openid
		result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "opendownplayer2", &msg)
		if err != nil || string(result) == "false" {
			w.Write([]byte("找不到玩家"))
			return
		}
		w.Write(result)
	} else {
		w.Write([]byte("参数错误"))
	}
}

//! 得到在线玩家
//! 玩家的数据结构
type DB_Person struct {
	Uid      int64  `json:"uid"`
	Name     string `json:"name"`   //! 名字
	Imgurl   string `json:"imgurl"` //! 头像
	Card     int    `json:"card"`   //! 房卡数量
	Gold     int    `json:"gold"`   //! 金卡数量
	RoomId   int    `json:"roomid"` //! 当前处于哪个room中
	GameType int    `json:"gametype"`
	Sex      int    `json:"sex"`      //! 性别
	Time     int64  `json:"time"`     //! 登陆时间
	IP       string `json:"ip"`       //! 登陆ip
	Sign     string `json:"sign"`     //! 签名
	SaveGold int    `json:"savegold"` //! 多少存款
}

type Msg_Online struct {
	Num  int         `json:"num"`
	Info []DB_Person `json:"info"`
}

type Msg_OnlineNum struct {
	Num int `json:"num"`
}

func GetOnlinePerson(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	lst := GetPersonMgr().GetOnline()
	var msg Msg_Online
	msg.Num = len(lst)
	for i := 0; i < len(lst); i++ {
		var person DB_Person
		value := GetServer().DB_GetData("user", lst[i])
		if string(value) != "" {
			json.Unmarshal(value, &person)
		} else { //! redis读不到，换服务器获取
			var _msg staticfunc.Msg_Uid
			_msg.Uid = lst[i]
			result, err := GetServer().CallLogin("ServerMethod.ServerMsg", "getperson", &_msg)
			if err != nil || string(result) == "" {
				continue
			}
			json.Unmarshal(result, &person)
		}
		msg.Info = append(msg.Info, person)
	}
	w.Write(lib.HF_JtoB(&msg))
}

func GetOnlinePersonNum(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	lst := GetPersonMgr().GetOnline()
	var msg Msg_OnlineNum
	msg.Num = len(lst)
	w.Write(lib.HF_JtoB(&msg))
}
