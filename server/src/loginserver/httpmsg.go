package loginserver

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io"
	"io/ioutil"
	"lib"
	"net/http"
	"rjmgr"
	//"net/url"
	"runtime/debug"
	"staticfunc"
	//"strings"
	"sync"
	"time"
)

//	"wxfd3530d2ed01875d", //! 颂游旗舰版4
//	"wx7276eb6f64caf74c", //! 榕迹6
//	"wx59691747a3e46ced", //! 至尊荣耀8
//	"wxabf47a88c67d225a", //! 颂游黑金10
//	"wxccfc348e0727bba5", //! 聚友演示版12
//	"wxf56a3742edf4d90c", //! 宝马棋牌 13
//	"wxb02bb19007803e40", //! 028棋牌　16
//	"wx15db0ef3f03a0288", //! 颂游贪玩17
//	"wxd06bc22daa56c7d8", //! 演示娱乐18
//	"wx5aa5f51e9526ca51", //! 028黑金版19
//	"wx73d30bdfbca35f1c", //! 至尊棋牌20
//	"wx511eaf329547acb5"} //! 演示娱乐2宝马版 21

//var APPSECRET = []string{"4aec89cc9201572b68b2dd0da7f9c277",
//	"8a3eb9a0b35ee7708f2602546aca5e16",
//	"92c3f0b01c478a29d3f0124222991bf7",
//	"53e980285f34e3278c2943a551df00d4",
//	"7bbdb625aa794591cd70d4fd6fc7f368",
//	"ec0678e936f7547258fc1b85d6b38d24",
//	"98cd1713df3f3f06d82aa9169de8cee8",
//	"29830ee28bd500c2e4c08ed51ae5214d",
//	"0539d44ff62326ed586978aac4adcff9",
//	"288493e6f1b26656fa208f48ed6b0e0f",
//	"d238aeb1392948a136a3f0899df3f3b1",
//	"4d860f94e3bb95437bbdc7b0b9b3d231",
//	"deae67f9c74c6ff0b70607dc000703a7",
//	"f90be6cfa8501b883987530e1627a15e",
//	"1ec2cf52d400a7ec2394db7e6345f62f",
//	"cd2c4b4eb76b4df965c6b5ca7590a247",
//	"2b8c5945120a7a5defdfc7cd78c1b9a7",
//	"b2dcd0f81435322f8581c9229e561b89",
//	"86953faa10eace25bb8d5cc7764f182a",
//	"4520c7f2b0f1cfbad2701be9a51d01cf", //!028黑金版19
//	"c3b3747e0df43e6d065affc0d17db3a0",
//	"dfb70c964a5750ad0033d3aed5342f3e"}

const BASEURL = "https://api.weixin.qq.com/sns"

//! 不要修改
const APP_KEY = "DEdeifnSDKRFIOJSNDDENGKsddssdsd"

//! 微信登陆结构
type Code2Token struct {
	Access_token string `json:"access_token"`
	Openid       string `json:"openid"`
}

type Token2UserInfo struct {
	Openid     string `json:"openid"`
	Unionid    string `json:"unionid"`
	Nickname   string `json:"nickname"`
	Headimgurl string `json:"headimgurl"`
	Sex        int    `json:"sex"`
}

//! 生成md5
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//! 生成Guid字串
func GetGuid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

type OldKWXCommand struct {
	CommId    int    `json:"commId"`
	Uid       string `json:"uid"`
	Num       int    `json:"num"`
	ProductId int    `json:"productId"`
	Text      string `json:"text"`
}

type OldKWXResult struct {
	Success int `json:"success"`
	Number  int `json:"number"`
}

/////////////////////////////////////////////////////////////////////
//! 玩家的数据结构
type Person struct {
	Uid      int64  `json:"uid"`
	Name     string `json:"name"`   //! 名字
	Imgurl   string `json:"imgurl"` //! 头像
	Card     int    `json:"card"`   //! 房卡数量
	Gold     int    `json:"gold"`   //! 金卡数量
	GameId   int    `json:"gameid"` //! 当前处于哪个game中
	RoomId   int    `json:"roomid"` //! 当前处于哪个room中
	GameType int    `json:"gametype"`
	Sex      int    `json:"sex"`  //! 性别
	Time     int64  `json:"time"` //! 登陆时间
	IP       string `json:"ip"`   //! 登陆ip
	Sign     string `json:"sign"` //! 签名
	UnionId  string `json:"unionid"`
	OpenId   string `json:"openid"`
	BindGold int    `json:"bindgold"` //! 绑定金币数量
	SaveGold int    `json:"savegold"` //! 多少存款
	Admin    int    `json:"admin"`
}

func (self *Person) Flush(db bool) {
	GetServer().DB_SetData("user", self.Uid, lib.HF_JtoB(self), self.Gold, self.SaveGold, db)
}

//! 同步金币到游戏
func (self *Person) SynchroGold() bool {
	if self.RoomId == 0 || self.GameId == 0 {
		return true
	}

	config := GetServer().GetGameServer(self.GameId)
	if config == nil {
		return true
	}

	var msg staticfunc.Msg_SynchroGold
	msg.Uid = self.Uid
	msg.Gold = self.Gold
	config.Call("ServerMethod.ServerMsg", "synchrogold", &msg)
	return false
}

//! 同步房卡和金币给客户端
func (self *Person) SynchroMoney() {
	var msg staticfunc.Msg_Synchro
	msg.Uid = self.Uid
	msg.Gold = self.Gold
	msg.Card = self.Card
	GetServer().CallCenter("ServerMethod.ServerMsg", "synchro", &msg)
}

//! 由于登陆服务器消息不多，故把消息全部写在这里
type C2S_Login struct { //! 登陆消息
	Code     string `json:"code"`
	Type     int    `json:"type"`
	Ver      int    `json:"ver"`
	AssetKey string `json:"assetkey"`
}

type C2S_LoginPwd struct { //! 登陆消息
	Account  string `json:"account"`
	Password string `json:"password"`
	Ver      int    `json:"ver"`
	Type     int    `json:"type"`
	AssetKey string `json:"assetkey"`
}

type C2S_Create struct { //! 创建房间
	Uid     int64  `json:"uid"`
	Type    int    `json:"type"`
	Num     int    `json:"num"`
	Param1  int    `json:"param1"`
	Param2  int    `json:"param2"`
	Agent   bool   `json:"agent"`
	UnionId string `json:"unionid"`
}

type C2S_Join struct { //! 加入房间
	Uid     int64  `json:"uid"`
	RoomId  int    `json:"roomid"`
	Group   int    `json:"group"`
	UnionId string `json:"unionid"`
}

type C2S_Agent struct { //! 绑定代理
	Uid   int64 `json:"uid"`
	Agent int   `json:"agent"`
}

type C2S_GetRoomList struct { //! 得到房间列表
	Uid     int64  `json:"uid"`
	UnionId string `json:"unionid"`
}

type C2S_SaveGold struct { //! 得到房间列表
	Uid     int64  `json:"uid"`
	Gold    int    `json:"gold"`
	UnionId string `json:"unionid"`
}

type C2S_AGiveGoldToB struct { //! 转增
	AUid    int64  `json:"auid"`
	Gold    int    `json:"gold"`
	UnionId string `json:"unionid"`
	BUid    int64  `json:"buid"`
}

type C2S_CardExchange struct { //! 得到房间列表
	Uid  int64 `json:"uid"`
	Type int   `json:"type"`
}
type C2S_Card2Gold struct { //!房卡换钻石
	Uid    int64 `json:"uid"`
	Type   int   `json:"type"`
	Amount int   `json:"amount"`
}
type C2S_Record struct { //! 得到战报
	Uid  int64 `json:"uid"`
	Type int   `json:"type"`
}
type C2S_GetTop struct { //! 得到排行榜
	Ver int `json:"ver"`
}

type C2S_Fish struct { //!
	Uid  int64 `json:"uid"`
	Type int   `json:"type"`
}

type GameEnter struct {
	DF int `json:"df"`
	ZR int `json:"zr"`
}

type S2C_Login struct { //! 登陆消息
	Uid       int64               `json:"uid"`
	Name      string              `json:"name"`
	Imgurl    string              `json:"imgurl"`
	Gold      int                 `json:"gold"`
	SaveGold  int                 `json:"savegold"`
	Card      int                 `json:"card"`
	Ip        string              `json:"ip"`
	Sign      string              `json:"sign"`
	Sex       int                 `json:"sex"`
	Room      int                 `json:"room"`
	Center    string              `json:"center"`
	Openid    string              `json:"openid"`
	Passwd    string              `json:"passwd"`
	Ver       bool                `json:"ver"`
	GameOver  int                 `json:"gameover"`
	Key       string              `json:"key"`
	GameMode  staticfunc.GameMode `json:"gamemode"`
	MoneyMode int                 `json:"moneymode"`
	AssetKey  string              `json:"assetkey"`
	Config    rjmgr.BaseConfig    `json:"config"`
	Enter     map[int][]GameEnter `json:"enter"`
}

type S2C_Ret struct { //! 登陆失败
	Ret int `json:"ret"`
}

type S2C_ReDownload struct { //! 重新下载
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	Download string `json:"download"`
}

type S2C_CreateRoom struct { //! 加入房间
	Ip    string `json:"ip"`
	Room  int    `json:"room"`
	Type  int    `json:"type"`
	Card  int    `json:"card"`
	Gold  int    `json:"gold"`
	Agent bool   `json:"agent"`
}

type S2C_JoinRoom struct { //! 加入房间
	Ip   string `json:"ip"`
	Room int    `json:"room"`
}

type S2C_JoinAG struct { //! 加入房间
	Url      string `json:"url"`
	GameType int    `json:"gametype"`
}

type S2C_Err struct {
	Info string `json:"info"`
}

//! 得到房间列表
type S2C_GetRoomList struct {
	Info []*CreateRoomInfo `json:"info"`
}

//! 得到排行榜
type S2C_GetTop struct {
	Info lstTop `json:"info"`
	Ver  int    `json:"ver"`
}

//! 设置额外信息
type Msg_SetExtraInfo struct {
	Uid     int64  `json:"uid"`
	Sex     int    `json:"sex"`
	Sign    string `json:"sign"`
	UnionId string `json:"unionid"`
}

//! 得到游戏人数
type Msg_GetGameNum struct {
	Info []JS_GameNum `json:"info"`
}

//! 得到当前金币和存金币
type Msg_SaveGold struct {
	Gold     int `json:"gold"`
	SaveGold int `json:"savegold"`
}

//! 转账
type Msg_AGiveGoldToB struct {
	Auid int64 `json:"auid"`
	Gold int   `json:"gold"`
	Buid int64 `json:"buid"`
}
type Msg_Rescord_GiveGold struct {
	AUid  int64  `json:"auid"`
	AName string `json:"aname"`
	BUid  int64  `json:"buid"`
	BName string `json:"bname"`
	Gold  int    `json:"gold"`
	Time  int64  `json:"time"`
}
type S2C_GiveGoldRecord struct {
	Uid  int64                  `json:"uid"`
	Info []Msg_Rescord_GiveGold `json:"info"`
}

//! 兑换金币
type Msg_ExchangeGold struct {
	Gold int `json:"gold"`
}

type Son_ExchangeInfo struct {
	Gold  int   `json:"gold"`
	State int   `json:"state"`
	Time  int64 `json:"time"`
}

//! 兑换记录
type Msg_ExchangeInfo struct {
	Info []Son_ExchangeInfo `json:"info"`
}

type S2C_Record struct { //! 得到战报
	Type int      `json:"type"`
	Info []string `json:"info"`
}

type S2C_RecordKWX1 struct { //! 得到战报
	Info []Son_RecordKWX1 `json:"info"`
}
type Son_RecordKWX1 struct {
	Roomid int                `json:"roomid"`
	Person []Son_RecordPerson `json:"person"`
	Time   int64              `json:"time"`
}
type Son_RecordPerson struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Score int    `json:"score"`
}

type S2C_RecordKWX2 struct { //! 得到战报
	Info string `json:"info"`
}

type S2C_Token struct { //!
	Token string `json:"token"`
}

//! 得到鱼池
type S2C_GetFishDesk struct {
	Type int              `json:"type"`
	Desk [12]lib.FishDesk `json:"desk"`
}

func GetErr(info string) []byte {
	var msg S2C_Err
	msg.Info = info

	return lib.HF_EncodeMsg("err", &msg, false)
}

func GetRet(head string, ret int) []byte {
	var msg S2C_Ret
	msg.Ret = ret

	return lib.HF_EncodeMsg(head, &msg, true)
}

//! 版本消息
func VersionMsg(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if staticfunc.GetIpBlackMgr().IsIp(clientip) { //! ip在黑名单里
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	//w.Header().Set("content-type", "application/json")             //返回数据格式是json

	w.Write([]byte(fmt.Sprintf("%d", GetServer().Con.Version)))
}

//! 0-100为玩家庄胜率 101为玩家庄随机 102为玩家庄控制
//! 0-100为系统奖池抽分
func GetDealMoneyValue(gametype int) lib.ManyMoney {
	c := GetServer().Redis.Get()
	defer c.Close()

	var property lib.ManyMoney
	value, err := redis.Bytes(c.Do("GET", fmt.Sprintf("deal_money_%d", gametype)))
	if err == nil {
		json.Unmarshal(value, &property)
	} else {
		property = lib.DefaultManyMoney
	}

	return property
}

func SetDealMoneyValue(gametype int, property lib.ManyMoney) bool {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][gametype]
	if !ok {
		return false
	}

	config := GetServer().GetGameServer(lib.HF_Atoi(csv["gametype"]))
	if config == nil {
		return false
	}

	c := GetServer().Redis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("deal_money_%d", gametype), lib.HF_JtoB(&property))

	var msg staticfunc.Msg_SetDealMoney
	msg.GameType = gametype
	msg.Property = property
	config.Call("ServerMethod.ServerMsg", "setdealmoney", &msg)

	return true
}

//! 得到筹码
func GetDealMoney(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetDealMoney") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	gametype := lib.HF_Atoi(req.FormValue("gametype"))

	property := GetDealMoneyValue(gametype)

	w.Write(lib.HF_JtoB(&property))
}

//! 设置豹子王庄家胜率
func SetDealMoney(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "SetDealMoney") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	gametype := lib.HF_Atoi(req.FormValue("gametype"))
	value := req.FormValue("value")

	var property lib.ManyMoney
	err := json.Unmarshal([]byte(value), &property)
	if err != nil {
		w.Write([]byte("设置失败"))
		return
	}

	if SetDealMoneyValue(gametype, property) {
		w.Write([]byte("设置成功"))
	} else {
		w.Write([]byte("设置失败"))
	}
}

//! 0-100为玩家庄胜率 101为玩家庄随机 102为玩家庄控制
//! 0-100为系统奖池抽分
func GetDealValue(gametype int) lib.ManyProperty {
	c := GetServer().Redis.Get()
	defer c.Close()

	var property lib.ManyProperty
	value, err := redis.Bytes(c.Do("GET", fmt.Sprintf("deal_property_%d", gametype)))
	if err == nil {
		json.Unmarshal(value, &property)
	} else {
		property = lib.DefaultManyProperty
	}

	if property.BetTime == 0 {
		switch gametype / 10000 {
		case 9: //! 神仙夺宝
			property.BetTime = 15
		case 24: //! 鱼虾蟹
			property.BetTime = 17
		case 26: //! 百家乐
			property.BetTime = 16
		case 14: //! 单双
			property.BetTime = 14
		case 23: //! 百人牛牛
			property.BetTime = 18
		case 21: //! 红黑大战
			property.BetTime = 18
		case 6: //! 百人推筒子
			property.BetTime = 20
		case 10: //! 龙虎斗
			property.BetTime = 21
		case 4: //! 豹子王
			property.BetTime = 19
		case 20: //! 名品汇
			property.BetTime = 14
		case 12: //! 骰宝
			property.BetTime = 19
		}
	}

	return property
}

func SetDealValue(gametype int, property lib.ManyProperty) bool {
	csv, ok := staticfunc.GetCsvMgr().Data["game"][gametype]
	if !ok {
		return false
	}

	config := GetServer().GetGameServer(lib.HF_Atoi(csv["gametype"]))
	if config == nil {
		return false
	}

	if property.DealCost > 100 || property.DealCost < 0 {
		property.DealCost = 5
	}

	if property.PlayerCost > 102 || property.PlayerCost < 0 {
		property.PlayerCost = 102
	}

	if property.BetTime == 0 {
		switch gametype / 10000 {
		case 9: //! 神仙夺宝
			property.BetTime = 15
		case 24: //! 鱼虾蟹
			property.BetTime = 17
		case 26: //! 百家乐
			property.BetTime = 16
		case 14: //! 单双
			property.BetTime = 14
		case 23: //! 百人牛牛
			property.BetTime = 18
		case 21: //! 红黑大战
			property.BetTime = 18
		case 6: //! 百人推筒子
			property.BetTime = 20
		case 10: //! 龙虎斗
			property.BetTime = 21
		case 4: //! 豹子王
			property.BetTime = 19
		case 20: //! 名品汇
			property.BetTime = 14
		case 12: //! 骰宝
			property.BetTime = 19
		}
	}

	c := GetServer().Redis.Get()
	defer c.Close()

	c.Do("SET", fmt.Sprintf("deal_property_%d", gametype), lib.HF_JtoB(&property))

	var msg staticfunc.Msg_SetDealWin
	msg.GameType = gametype
	msg.Property = property
	config.Call("ServerMethod.ServerMsg", "setdealwin", &msg)

	return true
}

//! 得到豹子王庄家胜率
func GetDealWin(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetDealWin") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	gametype := lib.HF_Atoi(req.FormValue("gametype"))

	property := GetDealValue(gametype)

	w.Write(lib.HF_JtoB(&property))
}

//! 设置豹子王庄家胜率
func SetDealWin(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "SetDealWin") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	gametype := lib.HF_Atoi(req.FormValue("gametype"))
	value := req.FormValue("value")

	var property lib.ManyProperty
	err := json.Unmarshal([]byte(value), &property)
	if err != nil {
		w.Write([]byte("设置失败"))
		return
	}

	if SetDealValue(gametype, property) {
		w.Write([]byte("设置成功"))
	} else {
		w.Write([]byte("设置失败"))
	}
}

func SetDealNext(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "SetDealNext") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	gametype := lib.HF_Atoi(req.FormValue("gametype"))
	next := req.FormValue("next")

	csv, ok := staticfunc.GetCsvMgr().Data["game"][gametype]
	if !ok {
		w.Write([]byte("设置失败"))
		return
	}

	config := GetServer().GetGameServer(lib.HF_Atoi(csv["gametype"]))
	if config == nil {
		w.Write([]byte("设置失败"))
		return
	}

	roomid := GetServer().ManyRoodId[gametype]
	if roomid == 0 {
		w.Write([]byte("设置失败"))
		return
	}

	c := GetServer().Redis.Get()
	defer c.Close()

	var msg staticfunc.Msg_SetDealNext
	msg.RoomId = roomid
	json.Unmarshal([]byte(next), &msg.Next)
	config.Call("ServerMethod.ServerMsg", "setdealnext", &msg)

	lib.GetLogMgr().Output(lib.LOG_ERROR, msg)

	w.Write([]byte("设置成功"))
}

func GetInitMoney(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetInitMoney") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	w.Write([]byte(fmt.Sprintf("%d", GetServer().InitMoney)))
}

func SetInitMoney(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "SetInitMoney") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	money := lib.HF_Atoi(req.FormValue("money"))
	GetServer().InitMoney = money

	c := GetServer().Redis.Get()
	defer c.Close()

	c.Do("SET", "initmoney", money)

	w.Write([]byte("设置成功"))
}

func GetWZQMode(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetWZQMode") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	w.Write(lib.HF_JtoB(&GetServer().WZQMode))
}

func SetWZQMode(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "SetWZQMode") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	value := req.FormValue("value")

	var mode staticfunc.WZQMode
	err := json.Unmarshal([]byte(value), &mode)
	if err != nil {
		w.Write([]byte("设置失败"))
		return
	}

	GetServer().WZQMode = mode

	c := GetServer().Redis.Get()
	defer c.Close()

	c.Do("SET", "wzqmode", value)

	w.Write([]byte("设置成功"))
}

//! 游戏开关
type GameModeMsg struct {
	GoldGame    []int `json:"goldgame"`    //! 金币场开启的游戏
	RoomGame    []int `json:"roomgame"`    //! 房卡场开启的游戏
	CanGoldGame []int `json:"cangoldgame"` //! 金币场可以开启的游戏
	CanRoomGame []int `json:"canroomgame"` //! 房卡场可以开启的游戏
}

func GetGameMode(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetGameMode") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	var msg GameModeMsg
	msg.GoldGame = GetServer().GameMode.GoldGame
	msg.RoomGame = GetServer().GameMode.RoomGame
	if len(rjmgr.GetRJMgr().GoldGame) == 0 {
		msg.CanGoldGame = ALLGOLDGAME
	} else {
		msg.CanGoldGame = rjmgr.GetRJMgr().GoldGame
	}
	if len(rjmgr.GetRJMgr().RoomGame) == 0 {
		msg.CanRoomGame = ALLROOMGAME
	} else {
		msg.CanRoomGame = rjmgr.GetRJMgr().RoomGame
	}

	w.Write(lib.HF_JtoB(&msg))
}

func SetGameMode(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "SetGameMode") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	value := req.FormValue("value")

	var mode staticfunc.GameMode
	err := json.Unmarshal([]byte(value), &mode)
	if err != nil {
		w.Write([]byte("设置失败"))
		return
	}

	GetServer().GameMode = mode

	c := GetServer().Redis.Get()
	defer c.Close()

	c.Do("SET", "gamemode", value)

	w.Write([]byte("设置成功"))
}

func GetGameRobotSet(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetGameRobotSet") {
		return
	}

	gametype := lib.HF_Atoi(req.FormValue("gametype"))

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	c := GetServer().Redis.Get()
	defer c.Close()
	var value lib.GameRobotSet
	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("robotset_%d", gametype)))
	if err != nil { //! 找不到设置默认值
		value = lib.DefaultGameRobotSet
	} else {
		json.Unmarshal(v, &value)
	}
	w.Write(lib.HF_JtoB(&value))
}

//!-------------------------------------- 李逵劈鱼
func SetGameRobotSet(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "SetGameRobotSet") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	gametype := lib.HF_Atoi(req.FormValue("gametype"))
	value := req.FormValue("value")

	var property lib.GameRobotSet
	err := json.Unmarshal([]byte(value), &property)
	if err != nil {
		w.Write([]byte("设置失败"))
		return
	}

	csv, ok := staticfunc.GetCsvMgr().Data["game"][gametype]
	if !ok {
		w.Write([]byte("设置失败"))
		return
	}

	config := GetServer().GetGameServer(lib.HF_Atoi(csv["gametype"]))
	if config == nil {
		w.Write([]byte("设置失败"))
		return
	}

	c := GetServer().Redis.Get()
	defer c.Close()
	c.Do("SET", fmt.Sprintf("robotset_%d", gametype), lib.HF_JtoB(&property))

	var msg staticfunc.Msg_SetGameRobotSet
	msg.GameType = gametype
	msg.Set = property
	config.Call("ServerMethod.ServerMsg", "setgamerobotset", &msg)

	w.Write([]byte("设置成功"))
}

func SetLkpyFish(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "SetFish") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	fishid := lib.HF_Atoi(req.FormValue("fishid"))
	value := req.FormValue("value")

	var property staticfunc.FishProperty
	err := json.Unmarshal([]byte(value), &property)
	if err != nil {
		w.Write([]byte("设置失败"))
		return
	}

	staticfunc.GetFishMgr().SetLKPYFishProperty(fishid, GetServer().Redis, property)

	w.Write([]byte("设置成功"))
}

func GetLkpyFish(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetFish") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	w.Write(lib.HF_JtoB(staticfunc.GetFishMgr().GetAllLKPYFishProperty(GetServer().Redis)))
}

func SetLkpyGun(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "SetGun") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	gunid := lib.HF_Atoi(req.FormValue("gunid"))
	value := req.FormValue("value")
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "value : ", value)
	var property staticfunc.GunProperty
	err := json.Unmarshal([]byte(value), &property)
	if err != nil {
		w.Write([]byte("设置失败"))
		return
	}

	staticfunc.GetFishMgr().SetLKPYGunProperty(gunid, GetServer().Redis, property)

	w.Write([]byte("设置成功"))
}

func GetLkpyGun(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetGun") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	w.Write(lib.HF_JtoB(staticfunc.GetFishMgr().GetAllLKPYGunProperty(GetServer().Redis)))
}

func SetLkpyFishValue(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "SetFishValue") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	value := req.FormValue("fishvalue")
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "value : ", value)
	var property staticfunc.FishValue
	err := json.Unmarshal([]byte(value), &property)
	if err != nil {
		w.Write([]byte("设置失败"))
		return
	}

	staticfunc.GetFishMgr().SetLKPYFishValue(GetServer().Redis, property)

	w.Write([]byte("设置成功"))
}

func GetLkpyFishValue(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetFishValue") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	w.Write(lib.HF_JtoB(staticfunc.GetFishMgr().GetLKPYFishValue(GetServer().Redis)))
}

//!------------------------------------ 捕鱼
func SetFish(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "SetFish") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	fishid := lib.HF_Atoi(req.FormValue("fishid"))
	value := req.FormValue("value")

	var property staticfunc.FishProperty
	err := json.Unmarshal([]byte(value), &property)
	if err != nil {
		w.Write([]byte("设置失败"))
		return
	}

	staticfunc.GetFishMgr().SetFishProperty(fishid, GetServer().Redis, property)

	w.Write([]byte("设置成功"))
}

func GetFish(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetFish") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	w.Write(lib.HF_JtoB(staticfunc.GetFishMgr().GetAllFishProperty(GetServer().Redis)))
}

func SetGun(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "SetGun") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	gunid := lib.HF_Atoi(req.FormValue("gunid"))
	value := req.FormValue("value")
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "value : ", value)
	var property staticfunc.GunProperty
	err := json.Unmarshal([]byte(value), &property)
	if err != nil {
		w.Write([]byte("设置失败"))
		return
	}

	staticfunc.GetFishMgr().SetGunProperty(gunid, GetServer().Redis, property)

	w.Write([]byte("设置成功"))
}

func GetGun(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetGun") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	w.Write(lib.HF_JtoB(staticfunc.GetFishMgr().GetAllGunProperty(GetServer().Redis)))
}

func SetFishValue(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "SetFishValue") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	value := req.FormValue("fishvalue")
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "value : ", value)
	var property staticfunc.FishValue
	err := json.Unmarshal([]byte(value), &property)
	if err != nil {
		w.Write([]byte("设置失败"))
		return
	}

	staticfunc.GetFishMgr().SetFishValue(GetServer().Redis, property)

	w.Write([]byte("设置成功"))
}

func GetFishValue(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetFishValue") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	w.Write(lib.HF_JtoB(staticfunc.GetFishMgr().GetFishValue(GetServer().Redis)))
}

func SetFPJValue(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "SetFPJValue") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	_type := lib.HF_Atoi(req.FormValue("type"))
	value := req.FormValue("value")
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "value : ", value)
	var property []int
	err := json.Unmarshal([]byte(value), &property)
	if err != nil {
		w.Write([]byte("设置失败"))
		return
	}

	c := GetServer().Redis.Get()
	defer c.Close()
	c.Do("SET", fmt.Sprintf("fkfpj_pro_%d", _type), value)

	w.Write([]byte("设置成功"))
}

func GetFPJValue(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetFPJValue") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	_type := lib.HF_Atoi(req.FormValue("type"))

	var value []int

	c := GetServer().Redis.Get()
	defer c.Close()
	v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("fkfpj_pro_%d", _type)))
	if err != nil { //! 找不到设置默认值
		value = lib.DefaultFKFPJPro
	} else {
		json.Unmarshal(v, &value)
	}

	w.Write(lib.HF_JtoB(value))
}

var AccountLock *sync.RWMutex = new(sync.RWMutex)

//! openid转uid
func WxToUid(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "WxToUid") {
		return
	}

	AccountLock.Lock()
	defer AccountLock.Unlock()

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	unionid := req.FormValue("unionid")
	openid := req.FormValue("openid")
	avatar := req.FormValue("avatar")
	nickname := req.FormValue("nickname")
	ip := req.FormValue("ip")
	buid := lib.HF_Atoi(req.FormValue("buid"))
	AddBindNode(ip, int64(buid), time.Now().Unix())

	if unionid == "" || openid == "" {
		w.Write([]byte(fmt.Sprintf("%d", 0)))
		return
	}
	uid := GetServer().DB_GetUid(unionid, openid, true, 0)
	value, _ := GetServer().DB_GetData("user", uid, 2)
	if string(value) == "" { //! 一个新用户
		var person Person
		person.Uid = uid
		person.Card = GetServer().Con.NewCard
		person.Gold = GetServer().InitMoney
		person.BindGold = GetServer().InitMoney
		person.Time = time.Now().Unix()
		person.IP = clientip
		person.Imgurl = avatar
		person.Name = nickname
		person.UnionId = unionid
		person.OpenId = openid
		person.Flush(true)
		if GetServer().InitMoney > 0 {
			GetServer().SqlGoldLog(person.Uid, GetServer().InitMoney, -5)
		}
	}

	w.Write([]byte(fmt.Sprintf("%d", uid)))
}

//! 注册
//func RegisterToUid(w http.ResponseWriter, req *http.Request) {
//	defer func() {
//		x := recover()
//		if x != nil {
//			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
//		}
//	}()

//	clientip := lib.HF_GetHttpIP(req)
//	if !GetServer().IsWhite(clientip, "RegisterToUid") {
//		return
//	}

//	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
//	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

//	account := req.FormValue("account")
//	passwd := req.FormValue("passwd")
//	avatar := req.FormValue("avatar")
//	nickname := req.FormValue("nickname")
//	uid := int64(lib.HF_Atoi(req.FormValue("uid")))
//	if uid == 0 {
//		w.Write([]byte(fmt.Sprintf("%d", 0)))
//		return
//	}

//	uid = GetServer().DB_GetUidFromPasswd(account, passwd, uid)
//	value, _ := GetServer().DB_GetData("user", uid, 2)
//	if string(value) == "" { //! 一个新用户
//		var person Person
//		person.Uid = uid
//		person.Card = GetServer().Con.NewCard
//		person.Gold = GetServer().InitMoney
//		person.BindGold = GetServer().InitMoney
//		person.Time = time.Now().Unix()
//		person.IP = clientip
//		person.Imgurl = avatar
//		person.Name = nickname
//		person.UnionId = account
//		person.OpenId = account
//		person.Flush(true)
//		if GetServer().InitMoney > 0 {
//			GetServer().SqlGoldLog(person.Uid, GetServer().InitMoney, -5)
//		}
//	}

//	w.Write([]byte(fmt.Sprintf("%d", uid)))
//}

//! 修改密码
//func ModifyPasswd(w http.ResponseWriter, req *http.Request) {
//	defer func() {
//		x := recover()
//		if x != nil {
//			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
//		}
//	}()

//	clientip := lib.HF_GetHttpIP(req)
//	if !GetServer().IsWhite(clientip, "RegisterToUid") {
//		return
//	}

//	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
//	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

//	account := req.FormValue("account")
//	passwd := req.FormValue("passwd")
//	newpasswd := req.FormValue("newpasswd")

//	uid := GetServer().DB_GetUidFromPasswd(account, passwd, 0)
//	if uid == 0 { //! 找不到这个账号密码
//		w.Write([]byte(fmt.Sprintf("%d", 0)))
//		return
//	} else {
//		GetServer().DB_ModifyPasswd(account, passwd, newpasswd, uid)
//		w.Write([]byte(fmt.Sprintf("%d", uid)))
//	}
//}

//! 得到在线人数
func GetOnlineNum(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetOnlineNum") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	w.Write([]byte(fmt.Sprintf("%d", GetNumMgr().GetNum())))
}

//! 得到不同游戏在线人数
func GetOnlineNumByGameType(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetOnlineNumByGameType") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	gametype := req.FormValue("gametype")

	w.Write([]byte(fmt.Sprintf("%d", GetNumMgr().GetNumByGameType(lib.HF_Atoi(gametype)))))
}

//! 强制解散房间
func ForceDissmiss(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "ForceDissmiss") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	roomid := lib.HF_Atoi(req.FormValue("roomid"))
	config := GetServer().GetGameRoom(roomid, -1)
	if config == nil {
		w.Write([]byte("true"))
		return
	}

	var msg staticfunc.Msg_JoinRoom
	msg.Id = roomid
	_, err := config.Call("ServerMethod.ServerMsg", "dissmiss", &msg)
	if err == nil {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}
}

//! 获取游戏比例
func GetMoneyMode(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "GetMoneyMode") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	w.Write([]byte(fmt.Sprintf("%d", GetServer().MoneyMode)))
}

//! 查找玩家
func FindPlayer(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "FindPlayer") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	uid := lib.HF_Atoi(req.FormValue("uid"))
	value, _ := GetServer().DB_GetData("user", int64(uid), 1)
	if string(value) != "" {
		var person Person
		json.Unmarshal(value, &person)
		var msg staticfunc.Msg_PersonInfo
		msg.Uid = person.Uid
		msg.Nickname = person.Name
		msg.Headurl = person.Imgurl
		msg.RoomId = person.RoomId
		msg.GameType = person.GameType
		var son staticfunc.Son_Card
		son.Productid = staticfunc.TYPE_CARD
		son.Repertory_count = person.Card
		msg.Repertory = append(msg.Repertory, son)
		son.Productid = staticfunc.TYPE_GOLD
		son.Repertory_count = person.Gold
		msg.Repertory = append(msg.Repertory, son)
		w.Write(lib.HF_JtoB(&msg))
		return
	}
	w.Write([]byte("false"))
}

type JS_TokenInfo struct {
	State int           `json:"state"`
	Info  Son_TokenInfo `json:"info"`
}
type Son_TokenInfo struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Openid   string `json:"openid"`
	Unionid  string `json:"unionid"`
	Gameid   int64  `json:"gameid"`
}

//! 得到信息
func GetInfoFromToken(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	token := req.FormValue("token")

	c := GetServer().Redis.Get()
	defer c.Close()

	var msg JS_TokenInfo
	v, err := redis.Bytes(c.Do("GET", token))
	if err != nil {
		msg.State = 1
	} else {
		var person Person
		err = json.Unmarshal(v, &person)
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "v err : ", err)
		if err != nil {
			msg.State = 1
		} else {
			msg.State = 0
			msg.Info.Avatar = person.Imgurl
			msg.Info.Gameid = person.Uid
			msg.Info.Nickname = person.Name
			msg.Info.Openid = person.OpenId
			msg.Info.Unionid = person.UnionId
		}
	}

	w.Write(lib.HF_JtoB(&msg))
}

func LoginMsg(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	if GetServer().RpcCenterServer == nil {
		return
	}

	clientip := lib.HF_GetHttpIP(req)
	if staticfunc.GetIpBlackMgr().IsIp(clientip) { //! ip在黑名单里
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	if GetServer().ShutDown { //! 已经关服了
		w.Write(GetErr("关服中"))
		return
	}

	msgdata := []byte(req.FormValue("msgdata"))
	if len(msgdata) >= 2048 {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "消息超过2048")
		return
	}

	head, _, data, ok := lib.HF_DecodeMsg(msgdata)
	if !ok {
		lib.GetLogMgr().Output(lib.LOG_INFO, "s2cerr:", string(msgdata))
		w.Write([]byte("消息解析错误"))
		staticfunc.GetIpBlackMgr().AddIp(clientip, "消息解析错误1")
		return
	}

	lib.GetLogMgr().Output(lib.LOG_INFO, "s2c:", head, "...", string(data))

	switch head {
	case "register":
		var msg C2S_Register
		json.Unmarshal(data, &msg)
		w.Write(Register(msg.Account, msg.Passwd, msg.NickName, clientip))
	case "chgpasswd":
		var msg C2S_ChgPasswd
		json.Unmarshal(data, &msg)
		w.Write(ChgPasswd(msg.Account, msg.Passwd, msg.NewPasswd, msg.NewPasswd1, clientip))
	case "modifyhead":
		var msg C2S_ModifyHead
		json.Unmarshal(data, &msg)
		w.Write(ModifyHead(msg.Uid, msg.Unionid, msg.Head))
	case "loginWX": //! 登陆
		var msg C2S_Login
		json.Unmarshal(data, &msg)
		w.Write(LoginFromCode(msg.Code, msg.Type, msg.Ver, msg.AssetKey, clientip))
	case "loginYK": //! 游客登陆
		if rjmgr.GetRJMgr().Guest == 0 {
			staticfunc.GetIpBlackMgr().AddIp(clientip, "不允许游客登陆")
			w.Write(GetErr("不允许游客登陆"))
			return
		}
		var msg C2S_Login
		json.Unmarshal(data, &msg)
		w.Write(LoginFromGuest(msg.Code, msg.Ver, msg.AssetKey))
	case "loginOP":
		var msg C2S_Login
		json.Unmarshal(data, &msg)
		w.Write(LoginFromOpenid(msg.Code, msg.Type, msg.Ver, msg.AssetKey, clientip))
	case "loginGuest": //! 游客登陆,非网页登陆
		var msg C2S_Login
		json.Unmarshal(data, &msg)
		w.Write(LoginFromYK(msg.Type, msg.Ver, msg.AssetKey, clientip))
	case "loginPwd":
		var msg C2S_LoginPwd
		json.Unmarshal(data, &msg)
		w.Write(LoginFromPasswd(msg.Account, msg.Password, msg.Type, msg.Ver, msg.AssetKey, clientip))
	case "create": //! 创建房间
		var msg C2S_Create
		json.Unmarshal(data, &msg)
		w.Write(Create(msg.Uid, msg.UnionId, msg.Type, msg.Num, msg.Param1, msg.Param2, msg.Agent, clientip))
	case "join": //! 加入房间
		var msg C2S_Join
		json.Unmarshal(data, &msg)
		w.Write(Join(msg.Uid, msg.UnionId, msg.RoomId, msg.Group, clientip))
	case "fastjoin": //! 快速加入游戏
		var msg C2S_Create
		json.Unmarshal(data, &msg)
		w.Write(FastJoin(msg.Uid, msg.UnionId, msg.Type, msg.Num, msg.Param1, msg.Param2, msg.Agent, clientip))
	case "exitag": //! 退出ag房间
		var msg staticfunc.Msg_Uid
		json.Unmarshal(data, &msg)
		w.Write(ExitAG(msg.Uid))
	case "getroomlist":
		var msg C2S_GetRoomList
		json.Unmarshal(data, &msg)
		w.Write(GetRoomList(msg.Uid))
	//case "card2gold":
	//	var msg C2S_Card2Gold
	//	json.Unmarshal(data, &msg)
	//	w.Write(GetCard2Gold(msg.Uid, msg.Type, msg.Amount))
	case "cardexchange":
		var msg C2S_CardExchange
		json.Unmarshal(data, &msg)
		w.Write(GetNewCard(msg.Uid, msg.Type))
	case "createtoken":
		var msg C2S_GetRoomList
		json.Unmarshal(data, &msg)
		w.Write(CreateToken(msg.Uid, msg.UnionId, clientip))
	case "gettop":
		var msg C2S_GetTop
		json.Unmarshal(data, &msg)
		w.Write(GetTop(msg.Ver))
	case "setextrainfo":
		var msg Msg_SetExtraInfo
		json.Unmarshal(data, &msg)
		w.Write(SetExtraInfo(msg.Uid, msg.UnionId, msg.Sex, msg.Sign))
	case "getextrainfo":
		var msg C2S_GetRoomList
		json.Unmarshal(data, &msg)
		w.Write(GetExtraInfo(msg.Uid))
	case "getgamenum":
		w.Write(GetGameNum())
	case "savegold":
		var msg C2S_SaveGold
		json.Unmarshal(data, &msg)
		w.Write(SaveGold(msg.Uid, msg.UnionId, msg.Gold))
	case "agivegoldtob":
		var msg C2S_AGiveGoldToB
		json.Unmarshal(data, &msg)
		w.Write(AGiveGoldToB(msg.AUid, msg.UnionId, msg.Gold, msg.BUid))
	case "getgiverecord":
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(GetGiveGoldRecord(msg.Uid, msg.Type))
	case "drawgold":
		var msg C2S_SaveGold
		json.Unmarshal(data, &msg)
		w.Write(DrawGold(msg.Uid, msg.UnionId, msg.Gold))
	case "exchangegold":
		var msg C2S_SaveGold
		json.Unmarshal(data, &msg)
		w.Write(ExchangeGold(msg.Uid, msg.Gold))
	case "exchangeinfo":
		var msg C2S_GetRoomList
		json.Unmarshal(data, &msg)
		w.Write(ExchangeInfo(msg.Uid))
	case "enterfish": //! 加入鱼池
		var msg C2S_Fish
		json.Unmarshal(data, &msg)
		w.Write(EnterFish(msg.Type, msg.Uid))
	case "quitfish": //! 退出鱼池
		var msg C2S_Fish
		json.Unmarshal(data, &msg)
		lib.GetFishMgr().Quit(msg.Type, msg.Uid)
	case "record": //! 得到战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(Record(msg.Uid, msg.Type, clientip))
	case "recordkwx1": //! 得到战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordKWX1(msg.Uid, msg.Type, clientip))
	case "recordkwx2": //! 得到战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordKWX2(msg.Uid, msg.Type, clientip))
	case "recordszmj1":
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordSZMJ1(msg.Uid, msg.Type, clientip))
	case "recordszmj2":
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordSZMJ2(msg.Uid, msg.Type, clientip))
	case "recordxzdd1":
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordXZDD1(msg.Uid, msg.Type, clientip))
	case "recordxzdd2":
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordXZDD2(msg.Uid, msg.Type, clientip))
	case "recordgc1": //! 得到战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordGC1(msg.Uid, msg.Type, clientip))
	case "recordgc2": //! 得到战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordGC2(msg.Uid, msg.Type, clientip))
	case "recordddz1": //! 得到战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordDDZ1(msg.Uid, msg.Type, clientip))
	case "recordddz2": //! 得到战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordDDZ2(msg.Uid, msg.Type, clientip))
	case "recordbdddz1": //! 得到战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordBDDDZ1(msg.Uid, msg.Type, clientip))
	case "recordbdddz2": //! 得到战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordBDDDZ2(msg.Uid, msg.Type, clientip))
	case "recorddbd1": //! 得到战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordDBD1(msg.Uid, msg.Type, clientip))
	case "recorddbd2": //! 得到战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordDBD2(msg.Uid, msg.Type, clientip))
	case "recordzjh1": //! 得到战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordZJH1(msg.Uid, msg.Type, clientip))
	case "recordzjh2": //! 得到战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordZJH2(msg.Uid, msg.Type, clientip))
	case "recordcsmj1": //!常熟战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordCSMJ1(msg.Uid, msg.Type, clientip))
	case "recordcsmj2":
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordCSMJ2(msg.Uid, msg.Type, clientip))
	case "recordaqmj1": //!安慶战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordAQMJ1(msg.Uid, msg.Type, clientip))
	case "recordaqmj2":
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordAQMJ2(msg.Uid, msg.Type, clientip))
	case "recordsyhmj1": //!上虞花战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordSYHMJ1(msg.Uid, msg.Type, clientip))
	case "recordsyhmj2":
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordSYHMJ2(msg.Uid, msg.Type, clientip))
	case "recordtsmj1": //!通山战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordTSMJ1(msg.Uid, msg.Type, clientip))
	case "recordtsmj2":
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordTSMJ2(msg.Uid, msg.Type, clientip))
	case "recordncmj1": //!南昌战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordNCMJ1(msg.Uid, msg.Type, clientip))
	case "recordncmj2":
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordNCMJ2(msg.Uid, msg.Type, clientip))
	case "recordqpddz1": //!枪炮斗地主
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordQPDDZ1(msg.Uid, msg.Type, clientip))
	case "recordqpddz2":
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordQPDDZ2(msg.Uid, msg.Type, clientip))
	case "recordxtrrhh1": //!仙桃晃晃
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordXTRRHH1(msg.Uid, msg.Type, clientip))
	case "recordxtrrhh2":
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordXTRRHH2(msg.Uid, msg.Type, clientip))
	case "recordgdcsmj1": //!潮汕麻将
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordGDCSMJ1(msg.Uid, msg.Type, clientip))
	case "recordgdcsmj2":
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordGDCSMJ2(msg.Uid, msg.Type, clientip))
	case "recordgymj1": //!涡阳麻将
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordGYMJ1(msg.Uid, msg.Type, clientip))
	case "recordgymj2":
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordGYMJ2(msg.Uid, msg.Type, clientip))
	case "recordhlgc1": //! 得到战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordHLGC1(msg.Uid, msg.Type, clientip))
	case "recordhlgc2": //! 得到战报
		var msg C2S_Record
		json.Unmarshal(data, &msg)
		w.Write(RecordHLGC2(msg.Uid, msg.Type, clientip))
	default:
		w.Write([]byte("消息解析错误"))
		staticfunc.GetIpBlackMgr().AddIp(clientip, "消息解析错误2")
	}
}

//! 生成一个token
func CreateToken(uid int64, unionid string, clientip string) []byte {
	if GetServer().Con.MinVersion == 99999 {
		return GetErr("亲爱的玩家，系统维护升级中，预计2分钟，感谢您的支持")
	}

	value, _ := GetServer().DB_GetData("user", uid, 1)
	if string(value) == "" {
		return GetErr("无法获取token")
	}

	var person Person
	json.Unmarshal(value, &person)

	if person.UnionId != "" && person.UnionId != unionid {
		return GetErr("授权失败,请重新打开游戏")
	}

	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return GetErr("无法获取token")
	}
	h := md5.New()
	h.Write([]byte(base64.URLEncoding.EncodeToString(b)))
	token := hex.EncodeToString(h.Sum(nil))

	c := GetServer().Redis.Get()
	defer c.Close()

	c.Do("SET", token, value)
	c.Do("EXPIRE", token, 60)

	var msg S2C_Token
	msg.Token = token
	return lib.HF_EncodeMsg("createtoken", &msg, true)
}

func LoginFromCode(code string, _type int, ver int, assetkey string, clientip string) []byte {
	if GetServer().Con.MinVersion == 99999 {
		return GetErr("亲爱的玩家，系统维护升级中，预计2分钟，感谢您的支持")
	}

	if ver > GetServer().Con.Version { //! 客户端版本过高
		return GetRet("loginfail", 1)
	} else if ver < GetServer().Con.MinVersion { //! 版本过低
		return ReDownload(0, "当前版本号过低，请下载最新版本进行游戏", _type)
	}

	if !rjmgr.GetRJMgr().AssetKeyTrue(assetkey) {
		return GetErr("您的版本过低,请下载最新版本")
	}

	if len(rjmgr.GetRJMgr().AppID) == 0 || len(rjmgr.GetRJMgr().AppSecret) == 0 {
		return GetErr("您的版本过低,请下载最新版本")
	}

	appid := ""
	appsecret := ""
	if _type >= len(rjmgr.GetRJMgr().AppID) {
		appid = rjmgr.GetRJMgr().AppID[0]
		appsecret = rjmgr.GetRJMgr().AppSecret[0]
	} else {
		appid = rjmgr.GetRJMgr().AppID[_type]
		appsecret = rjmgr.GetRJMgr().AppSecret[_type]
	}

	if GetServer().HasCode(code) {
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "重复的code")
		staticfunc.GetIpBlackMgr().AddIp(clientip, "重复code")
		return []byte("")
	}
	GetServer().AddCode(code)

	//! 先从微信读取用户信息
	str := fmt.Sprintf("%s/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", BASEURL, appid, appsecret, code)
	response, _ := http.Get(str)
	if response == nil {
		return GetErr("登陆失败,请稍后再试")
	}
	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	lib.GetLogMgr().Output(lib.LOG_INFO, "微信登陆:", string(body))
	var code2token Code2Token
	json.Unmarshal(body, &code2token)
	if code2token.Access_token == "" || code2token.Openid == "" {
		return GetErr("您的版本过低,请下载最新版本")
	}

	//! 从token获取用户信息
	str = fmt.Sprintf("%s/userinfo?access_token=%s&openid=%s", BASEURL, code2token.Access_token, code2token.Openid)
	response, _ = http.Get(str)
	if response == nil {
		return GetErr("您的版本过低,请下载最新版本response")
	}
	body, _ = ioutil.ReadAll(response.Body)
	response.Body.Close()
	lib.GetLogMgr().Output(lib.LOG_INFO, "微信登陆response:", string(body))
	var token2userinfo Token2UserInfo
	json.Unmarshal(body, &token2userinfo)
	if token2userinfo.Unionid == "" {
		if token2userinfo.Openid != "" {
			token2userinfo.Unionid = token2userinfo.Openid
		} else {
			return GetErr("微信系统繁忙，请稍后再试")
		}
	}

	var person Person
	uid := GetServer().DB_GetUid(token2userinfo.Unionid, token2userinfo.Openid, true, 0)
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	value, db := GetServer().DB_GetData("user", uid, 2)
	if string(value) == "" { //! 一个新用户
		//c := GetServer().Redis.Get()
		//defer c.Close()
		//result, err := c.Do("GET", clientip)

		person.Uid = uid
		person.Card = GetServer().Con.NewCard
		person.Gold = GetServer().InitMoney
		person.BindGold = GetServer().InitMoney
		//if err != nil || result == "" { //! 这个ip没有记录
		//	person.Gold = GetServer().Con.NewGold
		//	person.BindGold = GetServer().Con.NewGold
		//	c.Do("SET", clientip, "1")
		//	c.Do("EXPIRE", clientip, 86400)
		//} else { //! 有这个ip
		//	person.Gold = 0
		//	person.BindGold = 0
		//}
		person.Name = token2userinfo.Nickname
		person.Imgurl = token2userinfo.Headimgurl
		person.Sex = token2userinfo.Sex
		person.Time = time.Now().Unix()
		person.IP = clientip
		person.UnionId = token2userinfo.Unionid
		person.OpenId = token2userinfo.Openid
		person.Flush(true)
		if GetServer().InitMoney > 0 {
			GetServer().SqlGoldLog(person.Uid, GetServer().InitMoney, -5)
		}
	} else { //! 老用户
		json.Unmarshal(value, &person)
		if db && person.GameId != 0 { //! 若是从数据库读的数据，则说明redis的缓存没有了，房间信息被清空了
			person.GameId = 0
			person.RoomId = 0
			person.GameType = 0
		}
		person.Name = token2userinfo.Nickname
		person.Imgurl = token2userinfo.Headimgurl
		person.Time = time.Now().Unix()
		person.IP = clientip
		person.UnionId = token2userinfo.Unionid
		person.OpenId = token2userinfo.Openid
		person.Flush(true)
	}

	lib.GetAGVideoMgr().Register(fmt.Sprintf("rj_%d_%d", uid, uid), person.Name)
	bind := GetBindNode(clientip)
	if bind.UID != 0 { //! 尝试绑定
		res, err := lib.HF_Get(fmt.Sprintf("http://127.0.0.1:1231/abindb?auid=%d&aunionid=%s&aopenid=%s&buid=%d", person.Uid, person.UnionId, person.OpenId, bind.UID), 1)
		if err != nil { //! 调用出错
			AddBindNode(bind.IP, bind.UID, bind.Time)
		} else {
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil || string(body) == "" || string(body) == "false" { //! 没有绑定上
				AddBindNode(bind.IP, bind.UID, bind.Time)
			}
		}
	}

	Ip := ""
	if person.GameId != 0 && person.GameId != 999 { //! 当前处于某个房间中
		config := GetServer().GetGameServer(person.GameId)
		if config != nil {
			num := GetServer().JoinRoom(person.GameId, person.RoomId)
			if num >= 0 {
				Ip = config.ExIp
				GetServer().AddGameRoom(person.GameId, person.RoomId, person.GameType, num, person.Uid, clientip)
				GetNumMgr().AddGameOne(person.GameType, person.Uid)
			} else {
				GetNumMgr().DoneGameOne(person.GameType, person.Uid)
				person.GameId = 0
				person.RoomId = 0
				person.GameType = 0
				person.Flush(false)
			}
		} else { //! 服务器没开
			GetNumMgr().DoneGameOne(person.GameType, person.Uid)
			person.GameId = 0
			person.RoomId = 0
			person.GameType = 0
			person.Flush(false)
		}
	}

	var msg S2C_Login
	msg.Uid = person.Uid
	msg.Imgurl = person.Imgurl
	msg.Card = person.Card + person.Gold
	msg.Gold = person.Gold
	msg.Name = person.Name
	msg.Ip = Ip
	msg.Sign = person.Sign
	msg.Sex = person.Sex
	msg.Room = person.RoomId
	msg.Center = GetServer().ExCenterServer
	msg.Openid = person.UnionId
	msg.Passwd = person.OpenId
	msg.SaveGold = person.SaveGold
	msg.GameOver = person.GameType
	msg.Key = GetKey(person.Uid)
	msg.GameMode = GameModeToClient()
	msg.MoneyMode = GetServer().MoneyMode
	msg.AssetKey = lib.HF_MD5(assetkey + "2019")
	msg.Config = rjmgr.GetRJMgr().GetConfig(_type)
	msg.Enter = GetGameEnter()
	if ver == GetServer().Con.Version {
		msg.Ver = true
	}

	return lib.HF_EncodeMsg("login", &msg, true)
}

func LoginFromGuest(openid string, ver int, assetkey string) []byte {
	if ver > GetServer().Con.Version { //! 客户端版本过高
		return GetRet("loginfail", 1)
	} else if ver < GetServer().Con.MinVersion { //! 版本过低
		return GetErr("当前版本号过低，请下载最新版本进行游戏")
	}

	if openid == "" {
		openid = GetGuid()
	}

	var person Person
	uid := GetServer().DB_GetUid(openid, "", true, 0)
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	value, db := GetServer().DB_GetData("user", uid, 2)
	if string(value) == "" { //! 一个新用户
		person.Uid = uid
		person.Card = GetServer().Con.NewCard
		person.Gold = GetServer().Con.NewGold
		person.BindGold = 0
		person.Name = "游客"
		person.Imgurl = ""
		person.Time = time.Now().Unix()
		person.Flush(true)
		if GetServer().Con.NewGold > 0 {
			GetServer().SqlGoldLog(person.Uid, GetServer().Con.NewGold, -5)
		}
	} else { //! 老用户
		json.Unmarshal(value, &person)
		if db && person.GameId != 0 { //! 若是从数据库读的数据，则说明redis的缓存没有了，房间信息被清空了
			person.GameId = 0
			person.RoomId = 0
			person.GameType = 0
			person.Flush(false)
		}
		person.Time = time.Now().Unix()
		person.Flush(true)
	}

	lib.GetAGVideoMgr().Register(fmt.Sprintf("rj_%d_%d", uid, uid), person.Name)

	Ip := ""
	if person.GameId != 0 && person.GameId != 999 { //! 当前处于某个房间中
		config := GetServer().GetGameServer(person.GameId)
		if config != nil {
			num := GetServer().JoinRoom(person.GameId, person.RoomId)
			if num >= 0 {
				Ip = config.ExIp
				GetServer().AddGameRoom(person.GameId, person.RoomId, person.GameType, num, person.Uid, "")
				GetNumMgr().AddGameOne(person.GameType, person.Uid)
			} else {
				GetNumMgr().DoneGameOne(person.GameType, person.Uid)
				person.GameId = 0
				person.RoomId = 0
				person.GameType = 0
				person.Flush(false)
			}
		} else { //! 这个房间的服务器没开，则销毁房间
			GetNumMgr().DoneGameOne(person.GameType, person.Uid)
			person.GameId = 0
			person.RoomId = 0
			person.GameType = 0
			person.Flush(false)
		}
	}

	var msg S2C_Login
	msg.Uid = person.Uid
	msg.Imgurl = person.Imgurl
	msg.Card = person.Card + person.Gold
	msg.Gold = person.Gold
	msg.Name = person.Name
	msg.Ip = Ip
	msg.Sign = person.Sign
	msg.Sex = person.Sex
	msg.Room = person.RoomId
	msg.Center = GetServer().ExCenterServer
	msg.Openid = person.UnionId
	msg.Passwd = person.OpenId
	msg.SaveGold = person.SaveGold
	msg.GameOver = person.GameType
	msg.Key = GetKey(person.Uid)
	msg.GameMode = GameModeToClient()
	msg.MoneyMode = GetServer().MoneyMode
	msg.Config = rjmgr.GetRJMgr().GetConfig(0)
	msg.Enter = GetGameEnter()
	if ver == GetServer().Con.Version {
		msg.Ver = true
	}

	return lib.HF_EncodeMsg("login", &msg, true)
}

func LoginFromOpenid(openid string, _type int, ver int, assetkey string, clientip string) []byte {
	if GetServer().Con.MinVersion == 99999 {
		return GetErr("亲爱的玩家，系统维护升级中，预计2分钟，感谢您的支持")
	}

	if ver > GetServer().Con.Version { //! 客户端版本过高
		return GetRet("loginfail", 1)
	} else if ver < GetServer().Con.MinVersion { //! 版本过低
		return ReDownload(0, "当前版本号过低，请下载最新版本进行游戏", _type)
	}

	if !rjmgr.GetRJMgr().AssetKeyTrue(assetkey) {
		return GetErr("您的版本过低,请下载最新版本")
	}

	if openid == "" {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "openid错误")
		return []byte("")
	}

	var person Person
	uid := GetServer().DB_GetUid(openid, "", false, 0)
	if uid == 0 {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "openid登陆错误")
		return ReDownload(1, "当前版本号过低，请下载最新版本进行游戏", _type)
	}
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}
	value, db := GetServer().DB_GetData("user", uid, 1)
	if string(value) == "" { //! 一个新用户
		staticfunc.GetIpBlackMgr().AddIp(clientip, "openid登陆错误")
		return []byte("")
	} else { //! 老用户
		json.Unmarshal(value, &person)
		if db && person.GameId != 0 { //! 若是从数据库读的数据，则说明redis的缓存没有了，房间信息被清空了
			person.GameId = 0
			person.RoomId = 0
			person.GameType = 0
			person.Flush(false)
		}
		person.Time = time.Now().Unix()
		person.IP = clientip
		person.Flush(true)
	}

	lib.GetAGVideoMgr().Register(fmt.Sprintf("rj_%d_%d", uid, uid), person.Name)
	bind := GetBindNode(clientip)
	if bind.UID != 0 { //! 尝试绑定
		res, err := lib.HF_Get(fmt.Sprintf("http://127.0.0.1:1231/abindb?auid=%d&aunionid=%s&aopenid=%s&buid=%d", person.Uid, person.UnionId, person.OpenId, bind.UID), 1)
		if err != nil { //! 调用出错
			AddBindNode(bind.IP, bind.UID, bind.Time)
		} else {
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil || string(body) == "" || string(body) == "false" { //! 没有绑定上
				AddBindNode(bind.IP, bind.UID, bind.Time)
			}
		}
	}

	Ip := ""
	if person.GameId != 0 && person.GameId != 999 { //! 当前处于某个房间中
		config := GetServer().GetGameServer(person.GameId)
		if config != nil {
			num := GetServer().JoinRoom(person.GameId, person.RoomId)
			if num >= 0 {
				Ip = config.ExIp
				GetServer().AddGameRoom(person.GameId, person.RoomId, person.GameType, num, person.Uid, clientip)
				GetNumMgr().AddGameOne(person.GameType, person.Uid)
			} else {
				GetNumMgr().DoneGameOne(person.GameType, person.Uid)
				person.GameId = 0
				person.RoomId = 0
				person.GameType = 0
				person.Flush(false)
			}
		} else { //! 这个房间的服务器没开，则销毁房间
			GetNumMgr().DoneGameOne(person.GameType, person.Uid)
			person.GameId = 0
			person.RoomId = 0
			person.GameType = 0
			person.Flush(false)
		}
	}

	var msg S2C_Login
	msg.Uid = person.Uid
	msg.Imgurl = person.Imgurl
	msg.Card = person.Card + person.Gold
	msg.Gold = person.Gold
	msg.Name = person.Name
	msg.Ip = Ip
	msg.Sign = person.Sign
	msg.Sex = person.Sex
	msg.Room = person.RoomId
	msg.Center = GetServer().ExCenterServer
	msg.Openid = person.UnionId
	msg.Passwd = person.OpenId
	msg.SaveGold = person.SaveGold
	msg.GameOver = person.GameType
	msg.Key = GetKey(person.Uid)
	msg.GameMode = GameModeToClient()
	msg.MoneyMode = GetServer().MoneyMode
	msg.AssetKey = lib.HF_MD5(assetkey + "2019")
	msg.Config = rjmgr.GetRJMgr().GetConfig(_type)
	msg.Enter = GetGameEnter()
	if ver == GetServer().Con.Version {
		msg.Ver = true
	}

	return lib.HF_EncodeMsg("login", &msg, true)
}

//! 创建房间
func Create(uid int64, unionid string, _type int, num int, param1 int, param2 int, agent bool, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	if GetServer().Con.MinVersion == 99999 {
		return GetErr("亲爱的玩家，系统维护升级中，预计2分钟，感谢您的支持")
	}

	value, _ := GetServer().DB_GetData("user", uid, 1)
	if string(value) == "" { //! 新用户应该先登录
		staticfunc.GetIpBlackMgr().AddIp(clientip, "未创建角色却创建房间")
		return GetErr("新用户请先登录")
	}

	csv, ok := staticfunc.GetCsvMgr().Data["game"][_type]
	if !ok || !IsOpenGame(_type) {
		//staticfunc.GetIpBlackMgr().AddIp(clientip, "创建游戏类型错误")
		return GetErr("即将开放，敬请期待")
	}

	if !IsGameOk(_type) {
		return GetErr("即将开放，敬请期待")
	}

	if _type == 77 && GetServer().WZQMode.Mode == 1 { //! 五子棋
		if !GetServer().IsWZQWhite(uid) {
			return GetErr("您没有权限创建此游戏")
		}
	}

	var person Person
	json.Unmarshal(value, &person)

	if person.UnionId != "" && person.UnionId != unionid {
		return GetErr("授权失败,请重新打开游戏")
	}

	if person.GameId != 0 { //! 已经在一个房间，不能重新创建一个房间
		return GetErr("正在房间中")
	}

	if num <= 0 || num > lib.HF_Atoi(csv["maxstep"]) {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "创建错误局数")
		return GetErr("局数错误")
	}

	cost := 0
	if _type == 15 { //! 宿州特殊
		if param2 == 0 {
			if num == 1 {
				cost = 2
			} else if num == 2 {
				cost = 3
			} else {
				cost = 6
			}
		} else {
			if num == 1 {
				cost = 1
			} else if num == 2 {
				cost = 1
			} else {
				cost = 2
			}
		}
		if person.Card < cost {
			return GetErr("钻石不足")
		}
	} else if _type == 23 { //! 兰州十点半
		if param2 == 0 {
			if num == 1 {
				cost = 38
			} else if num == 2 {
				cost = 57
			} else {
				cost = 76
			}
		} else {
			if num == 1 {
				cost = 10
			} else if num == 2 {
				cost = 15
			} else {
				cost = 20
			}
		}
		if person.Card < cost {
			return GetErr("钻石不足")
		}
	} else if _type == 16 && num == 5 {
		cost = 1
		if person.Card < cost {
			return GetErr("钻石不足")
		}
	} else if _type == 24 { //！三公演义
		if param2/10%10 == 0 {
			if num == 1 {
				cost = 3
			} else {
				cost = 6
			}
		} else {
			if num == 1 {
				cost = 1
			} else {
				cost = 2
			}
		}
		if person.Card < cost {
			return GetErr("钻石不足")
		}
	} else if _type == 26 { //！扎旮旯
		if 0 == param2%10 {
			if num == 1 {
				cost = 3
			} else {
				cost = 6
			}
		} else {
			if num == 1 {
				cost = 1
			} else {
				cost = 2
			}
		}
		if person.Card < cost {
			return GetErr("钻石不足")
		}
	} else if _type == 51 {
		cost = num * (param1 / 10000 % 10) * 2
		lib.GetLogMgr().Output(lib.LOG_DEBUG, "癞子错误:", param1, num, cost, person.Card)
		if person.Card < cost {
			return GetErr("房卡不足")
		}
	} else if _type == 28 {
		if param2/100%10 != 3 {
			cost = param1 * 20
		} else {
			cost = param1 + 20
		}
		if person.Gold < cost {
			return GetErr("钻石不足")
		}
	} else if _type == 49 {
		cost = 0
	} else if _type == 77 {
		cost = 0
	} else if _type == 1000 {
		cost = 0
	} else {
		cost = num
		if person.Card < cost {
			return GetErr("房卡不足")
		}
	}

	//! 得到人数最少的房间
	var config *GameServerConfig = nil
	if csv["gametype"] == "2" { //! 创建麻将房
		config = GetServer().GetGameServerFromRoom(2)
	} else if csv["gametype"] == "1" { //! 创建扑克房
		config = GetServer().GetGameServerFromRoom(1)
	} else if csv["gametype"] == "3" { //！创建金币场
		config = GetServer().GetGameServerFromRoom(3)
	}

	//! 告诉游戏服务器有人玩家加入进来
	if config != nil {
		roomid := GetRoom().GetID()
		if _type == 23 { //! 兰州十点半
			roomid = GetRoom().GetLzTenhalfID()
		}
		if roomid == 0 {
			return GetErr("创建房间失败")
		}

		agentuid := uid
		if !agent {
			agentuid = 0
		}
		if GetServer().CreateRoom(config.Id, roomid, _type, num, param1, param2, agentuid, 0) {
			GetServer().AddGameRoom(config.Id, roomid, _type, 1, person.Uid, clientip)
			if agentuid == 0 { //! 不是代开房间，加入房间
				person.GameId = config.Id
				person.RoomId = roomid
				person.GameType = _type
				person.Flush(false)
				GetNumMgr().AddGameOne(person.GameType, person.Uid)
			} else {
				if lib.HF_Atoi(csv["costtype"]) == staticfunc.TYPE_GOLD {
					person.Gold = lib.HF_MaxInt(0, person.Gold-cost)
					person.Flush(true)
					GetServer().InsertLog(person.Uid, staticfunc.MOVE_GOLD, cost, clientip)
				} else {
					person.Card = lib.HF_MaxInt(0, person.Card-cost)
					person.Flush(true)
					GetServer().InsertLog(person.Uid, staticfunc.MOVE_CARD, cost, clientip)
				}
			}

			if lib.HF_Atoi(csv["view"]) == 1 || agentuid > 0 { //! 观战或代开模式
				GetCreateRoomMgr().AddRoom(uid, &CreateRoomInfo{roomid,
					_type,
					cost,
					param1,
					param2,
					0,
					lib.HF_Atoi(csv["maxnum"]),
					make([]staticfunc.JS_CreateRoomMem, 0),
					0,
					time.Now().Unix()})
			}
		} else {
			return GetErr("创建房间失败")
		}
	} else {
		return GetErr("找不到房间1")
	}

	var msg S2C_CreateRoom
	msg.Ip = config.ExIp
	msg.Room = person.RoomId
	msg.Type = _type
	msg.Card = person.Card + person.Gold
	msg.Gold = person.Gold
	msg.Agent = agent

	return lib.HF_EncodeMsg("create", &msg, true)
}

//! 加入房间
func Join(uid int64, unionid string, roomid int, group int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	if GetServer().Con.MinVersion == 99999 {
		return GetErr("亲爱的玩家，系统维护升级中，预计2分钟，感谢您的支持")
	}

	value, _ := GetServer().DB_GetData("user", uid, 1)
	if string(value) == "" { //! 新用户应该先登录
		staticfunc.GetIpBlackMgr().AddIp(clientip, "未创建角色却加入房间")
		return GetErr("新用户请先登录")
	}

	var person Person
	json.Unmarshal(value, &person)

	if person.UnionId != "" && person.UnionId != unionid {
		return GetErr("授权失败,请重新打开游戏")
	}

	if person.GameId != 0 { //! 已经在一个房间，不能重新加入
		return GetErr("正在房间中")
	}

	//! 得到该room房间
	config := GetServer().GetGameRoom(roomid, group)

	//! 房间存在
	num := 0
	if config != nil {
		num = GetServer().JoinRoom(config.Id, roomid)
		if num >= 0 {
			person.GameId = config.Id
			person.RoomId = roomid
			person.GameType, _ = GetServer().AddGameRoom(config.Id, roomid, 0, num+1, person.Uid, clientip)
			person.Flush(false)
			GetNumMgr().AddGameOne(person.GameType, person.Uid)
			if person.GameType/10000 == 25 {
				desk := lib.GetFishMgr().AddPerson(person.GameType%250000, person.RoomId, person.Uid, person.Name, person.Imgurl)
				if desk != nil {
					var msg staticfunc.Msg_FishAddDesk1
					msg.Desk = *desk
					msg.Lst = lib.GetFishMgr().Get(person.GameType % 250000)
					GetServer().CallCenter("ServerMethod.ServerMsg", "fishaddperson", &msg)
				}
			}
		} else {
			return GetErr("无法加入房间")
		}
	} else {
		if roomid/1000000 > 0 { //! 俱乐部房间
			person.GameId, num = GetServer().ActiveRoom(roomid)
			if person.GameId > 0 {
				person.RoomId = roomid
				person.GameType, config = GetServer().AddGameRoom(person.GameId, roomid, 0, num+1, person.Uid, clientip)
				person.Flush(false)
				GetNumMgr().AddGameOne(person.GameType, person.Uid)
			} else {
				return GetErr(fmt.Sprintf("找不到房间:%d", roomid))
			}
		} else {
			return GetErr(fmt.Sprintf("找不到房间:%d", roomid))
		}
	}

	var msg S2C_JoinRoom
	msg.Ip = config.ExIp
	msg.Room = roomid

	return lib.HF_EncodeMsg("join", &msg, true)
}

//! 快速加入游戏
func FastJoin(uid int64, unionid string, _type int, num int, param1 int, param2 int, agent bool, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	if GetServer().Con.MinVersion == 99999 {
		return GetErr("亲爱的玩家，系统维护升级中，预计2分钟，感谢您的支持")
	}

	value, _ := GetServer().DB_GetData("user", uid, 1)
	if string(value) == "" { //! 新用户应该先登录
		staticfunc.GetIpBlackMgr().AddIp(clientip, "未创建角色却创建房间")
		return GetErr("新用户请先登录")
	}

	csv, ok := staticfunc.GetCsvMgr().Data["game"][_type]
	if !ok || !IsOpenGame(_type) {
		//staticfunc.GetIpBlackMgr().AddIp(clientip, "创建游戏类型错误")
		return GetErr("即将开放，敬请期待!")
	}

	if !IsGameOk(_type) {
		return GetErr("即将开放，敬请期待。")
	}

	var person Person
	json.Unmarshal(value, &person)

	if person.UnionId != "" && person.UnionId != unionid {
		return GetErr("授权失败,请重新打开游戏")
	}

	if person.GameId == 999 && person.GameType == _type { //! ag不判断

	} else {
		if person.GameId != 0 { //! 已经在一个房间，不能重新创建一个房间
			return GetErr("正在房间中")
		}
	}

	noroomid := num
	num = 1

	if num <= 0 || num > lib.HF_Atoi(csv["maxstep"]) {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "创建错误局数")
		return GetErr("局数错误")
	}

	many := true
	if _type/10000 == 1 || _type/10000 == 2 || _type/10000 == 3 || _type/10000 == 5 || _type/10000 == 7 || _type/10000 == 8 || _type/10000 == 29 { //! 金币卡五星,牛牛金币场
		many = false
		if _type/10%10 < 0 || _type/10%10 > 5 {
			return GetErr("参数错误")
		}
		df := staticfunc.GetCsvMgr().GetZR(_type)
		if person.Gold < df {
			return GetErr("您的金币不足，请前往充值。")
		}
	} else if _type/10000 == 19 { //! 翻牌机
		if _type%10 < 0 || _type%10 > 2 {
			return GetErr("参数错误")
		}

		df := staticfunc.GetCsvMgr().GetZR(_type)
		if person.Gold < df {
			return GetErr("您的金币不足，请前往充值。")
		}
	} else if 31 == _type/10000 || 32 == _type/10000 { //! 大赢家拉霸，绝地求生
		if _type%10 < 0 || _type%10 > 2 {
			return GetErr("参数错误")
		}

		zr := staticfunc.GetCsvMgr().GetZR(_type)
		if person.Gold < zr {
			return GetErr("您的金币不足，请前往充值。")
		}
	} else {
		if person.Card+person.Gold < lib.HF_Atoi(csv["needgold"]) {
			return GetErr("您的金币不足，请前往充值。")
		}
	}

	if _type >= 1000000 && _type <= 1000050 { //! 这个范围是ag视讯
		isTest := "0"
		if person.OpenId == "" {
			isTest = "1"
		}
		mode := ""
		switch _type {
		case 1000000: //! 百家乐
			mode = "onlybac"
		case 1000001: //! 龙虎斗
			mode = "onlydgtg"
		case 1000002: //! 轮盘
			mode = "onlyrou"
		case 1000003: //! 骰宝
			mode = "onlysicbo"
		case 1000004: //! 牛牛
			mode = "onlyniuniu"
		case 1000005: //! 三公
			mode = "onlysamgong"
		case 1000006: //! 番摊
			mode = "onlyfantan"
		case 1000007: //! 多台
			mode = "onlymultiple"
		case 1000008: //! 鱼虾蟹
			mode = "onlyfishshrimpcrab"
		}
		url := lib.GetAGVideoMgr().LoginGame(fmt.Sprintf("rj_%d_%d", person.Uid, person.Uid), mode, isTest)
		if url == "" {
			return GetErr("暂时无法加入,请稍后再试")
		}
		person.GameId = 999      //! ag的gameid
		person.RoomId = 99999999 //! ag的room
		person.GameType = _type
		person.Flush(false)
		var msg S2C_JoinAG
		msg.Url = url
		msg.GameType = _type
		return lib.HF_EncodeMsg("joinag", &msg, true)
	}

	//! 得到人数最少的房间
	var config *GameServerConfig = GetServer().GetGameServerFromGoldRoom(lib.HF_Atoi(csv["gametype"]))
	if config == nil {
		return GetErr("找不到房间")
	}

	_, roomid := config.GetRoomId(_type, noroomid, clientip)
	lib.GetLogMgr().Output(lib.LOG_INFO, "快速加入的房间号：", roomid)
	if roomid == 0 || (_type/10000 == 25 && param1 > 0) { //! 需要创建房间
		lib.GetLogMgr().Output(lib.LOG_INFO, "创建房间")
		if agent {
			roomid = GetRoom().GetID()
		} else {
			roomid = GetGoldRoom().GetID()
		}
		if many {
			GetServer().ManyRoodId[_type] = roomid
		}
		if _type/10000 == 25 {
			desk := lib.GetFishMgr().CreateDesk(_type%250000, param1, roomid, person.Uid, person.Name, person.Imgurl)
			if desk != nil {
				var msg staticfunc.Msg_FishAddDesk1
				msg.Desk = *desk
				msg.Lst = lib.GetFishMgr().Get(_type % 250000)
				GetServer().CallCenter("ServerMethod.ServerMsg", "fishadddesk", &msg)
			} else {
				return GetErr("渔场正在维护中")
			}
			param1 = desk.Index
		}
		if GetServer().CreateRoom(config.Id, roomid, _type, num, param1, param2, 0, 0) {
			GetServer().AddGameRoom(config.Id, roomid, _type, 1, person.Uid, clientip)
			person.GameId = config.Id
			person.RoomId = roomid
			person.GameType = _type
			person.Flush(false)
		} else {
			return GetErr("创建房间失败")
		}

		GetNumMgr().AddGameOne(_type, uid)

		var msg S2C_CreateRoom
		msg.Ip = config.ExIp
		msg.Room = roomid
		msg.Type = _type
		msg.Card = person.Card + person.Gold
		msg.Gold = person.Gold
		msg.Agent = false

		return lib.HF_EncodeMsg("create", &msg, true)
	} else { //! 加入
		num := GetServer().JoinRoom(config.Id, roomid)
		if num >= 0 {
			GetServer().AddGameRoom(config.Id, roomid, _type, num+1, person.Uid, clientip)
			person.GameId = config.Id
			person.RoomId = roomid
			person.GameType = _type
			person.Flush(false)
			if person.GameType/10000 == 25 {
				desk := lib.GetFishMgr().AddPerson(person.GameType%250000, person.RoomId, person.Uid, person.Name, person.Imgurl)
				if desk != nil {
					var msg staticfunc.Msg_FishAddDesk1
					msg.Desk = *desk
					msg.Lst = lib.GetFishMgr().Get(person.GameType % 250000)
					GetServer().CallCenter("ServerMethod.ServerMsg", "fishaddperson", &msg)
				}
			}
		} else {
			return GetErr("无法加入房间")
		}

		GetNumMgr().AddGameOne(_type, uid)

		var msg S2C_JoinRoom
		msg.Ip = config.ExIp
		msg.Room = roomid

		return lib.HF_EncodeMsg("join", &msg, true)
	}
}

func ExitAG(uid int64) []byte {
	value, _ := GetServer().DB_GetData("user", uid, 1)

	if string(value) == "" { //! 新用户应该先登录
		return GetErr("新用户请先登录")
	}

	var person Person
	json.Unmarshal(value, &person)
	if person.GameId != 999 {
		var msg staticfunc.Msg_Null
		return lib.HF_EncodeMsg("exitag", &msg, true)
	}

	person.GameId = 0 //! ag的gameid
	person.RoomId = 0 //! ag的room
	person.GameType = 0
	person.Flush(false)
	var msg staticfunc.Msg_Null
	return lib.HF_EncodeMsg("exitag", &msg, true)
}

func GetCard2Gold(uid int64, _type int, amount int) []byte {
	value, _ := GetServer().DB_GetData("user", uid, 1)

	if string(value) == "" { //! 新用户应该先登录
		return GetErr("新用户请先登录")
	}

	var person Person
	json.Unmarshal(value, &person)

	if 61 == _type {
		if 1 == amount {
			if person.Card < 30 {
				return GetErr("钻石不足")
			}
			person.Card -= 30
			person.Gold += 36000
		} else if 2 == amount {
			if person.Card < 60 {
				return GetErr("钻石不足")
			}
			person.Card -= 60
			person.Gold += 75600
		} else if 3 == amount {
			if person.Card < 100 {
				return GetErr("钻石不足")
			}
			person.Card -= 100
			person.Gold += 132000
		} else if 4 == amount {
			if person.Card < 200 {
				return GetErr("钻石不足")
			}
			person.Card -= 200
			person.Gold += 276000
		} else if 5 == amount {
			if person.Card < 500 {
				return GetErr("钻石不足")
			}
			person.Card -= 500
			person.Gold += 720000
		} else if 6 == amount {
			if person.Card < 1000 {
				return GetErr("钻石不足")
			}
			person.Card -= 1000
			person.Gold += 1500000
		}
	} else if _type == 1 {
		if rjmgr.GetRJMgr().Guest == 0 { //! 不允许游客登陆的服务器不能用房卡换金币
			return GetErr("非法操作")
		}
		if amount <= 0 {
			return GetErr("错误")
		}
		if person.Card < amount {
			return GetErr("房卡不足")
		}
		person.Card -= amount
		person.Gold += amount * 100
		GetTopMgr().UpdData(&person)
		GetServer().SqlGoldLog(person.Uid, amount*100, -3)
		person.SynchroGold()
	} else {
		return GetErr("错误的游戏类型")
	}
	person.Flush(true)
	var msg staticfunc.S2C_UpdCard
	msg.Gold = person.Gold
	msg.Card = person.Card + person.Gold
	return lib.HF_EncodeMsg("updcard", &msg, true)
}

func GetNewCard(uid int64, _type int) []byte {
	value, _ := GetServer().DB_GetData("user", uid, 1)

	if string(value) == "" { //! 新用户应该先登录
		return GetErr("新用户请先登录")
	}

	var person Person
	json.Unmarshal(value, &person)

	if 1 == _type {
		if person.Gold < 20 {
			return GetErr("钻石不足")
		}
		person.Gold -= 20
		person.Card += 1
	} else if 2 == _type {
		if person.Gold < 100 {
			return GetErr("钻石不足")
		}
		person.Gold -= 100
		person.Card += 6
	} else if 3 == _type {
		if person.Gold < 200 {
			return GetErr("钻石不足")
		}
		person.Gold -= 200
		person.Card += 13
	} else if 4 == _type {
		if person.Gold < 400 {
			return GetErr("钻石不足")
		}
		person.Gold -= 400
		person.Card += 28
	} else {
		return GetErr("错误的购买类型")
	}
	person.Flush(true)
	var msg staticfunc.S2C_UpdCard
	msg.Gold = person.Gold
	msg.Card = person.Card + person.Gold
	return lib.HF_EncodeMsg("updcard", &msg, true)
}

func GetTop(ver int) []byte {
	if rjmgr.GetRJMgr().CloseTop == 1 {
		return GetErr("该功能暂未开放")
	}

	var msg S2C_GetTop
	msg.Info = GetTopMgr().GetTop(ver)
	msg.Ver = GetTopMgr().Ver
	return lib.HF_EncodeMsg("gettop", &msg, true)
}

func SetExtraInfo(uid int64, unionid string, sex int, sign string) []byte {
	value, _ := GetServer().DB_GetData("user", uid, 1)
	if string(value) == "" {
		return []byte("")
	}

	var person Person
	json.Unmarshal(value, &person)

	if person.UnionId != "" && person.UnionId != unionid {
		return GetErr("授权失败,请重新打开游戏")
	}

	person.Sex = sex
	person.Sign = sign
	person.Flush(true)

	var msg staticfunc.Msg_Null
	return lib.HF_EncodeMsg("setextrainfo", &msg, true)
}

func GetExtraInfo(uid int64) []byte {
	var msg Msg_SetExtraInfo
	msg.Uid = uid

	value, _ := GetServer().DB_GetData("user", uid, 1)
	if string(value) != "" {
		var person Person
		json.Unmarshal(value, &person)
		msg.Sex = person.Sex
		msg.Sign = person.Sign
	} else {
		robot := lib.GetRobotMgr().GetRobotFromId(uid)
		if robot != nil {
			msg.Sex = robot.Sex
			msg.Sign = robot.Sign
		}
	}

	return lib.HF_EncodeMsg("getextrainfo", &msg, true)
}

func GetGameNum() []byte {
	var msg Msg_GetGameNum
	msg.Info = GetNumMgr().GetGameNum()
	return lib.HF_EncodeMsg("getgamenum", &msg, true)
}

var SaveLock *sync.RWMutex = new(sync.RWMutex)

func SaveGold(uid int64, unionid string, gold int) []byte {
	SaveLock.Lock()
	defer SaveLock.Unlock()

	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	if gold < 0 {
		GetServer().AddBlack(uid)
		return GetErr("存入金币失败")
	}

	value, _ := GetServer().DB_GetData("user", uid, 1)
	if string(value) == "" {
		return GetErr("存入金币失败")
	}

	var person Person
	json.Unmarshal(value, &person)

	if person.UnionId != "" && person.UnionId != unionid {
		return GetErr("授权失败,请重新打开游戏")
	}

	if person.Gold < gold {
		return GetErr("存入金币失败")
	}

	if person.GameType != 0 { //! 任何游戏中都不允许存钱
		//GetServer().AddBlack(uid)
		return GetErr("存入金币失败")
	}

	if gold > 0 {
		person.Gold -= gold
		person.SaveGold += gold
		GetTopMgr().UpdData(&person)
		GetServer().SqlGoldLog(person.Uid, -gold, -10)
		person.SynchroGold()
		person.Flush(true)
	}

	var msg Msg_SaveGold
	msg.Gold = person.Gold
	msg.SaveGold = person.SaveGold
	return lib.HF_EncodeMsg("savegold", &msg, true)
}

var DrawLock *sync.RWMutex = new(sync.RWMutex)

func DrawGold(uid int64, unionid string, gold int) []byte {
	DrawLock.Lock()
	defer DrawLock.Unlock()

	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	if gold < 0 {
		GetServer().AddBlack(uid)
		return GetErr("取出金币失败")
	}

	value, _ := GetServer().DB_GetData("user", uid, 1)
	if string(value) == "" {
		return GetErr("取出金币失败")
	}

	var person Person
	json.Unmarshal(value, &person)

	if person.UnionId != "" && person.UnionId != unionid {
		return GetErr("授权失败,请重新打开游戏")
	}

	if person.SaveGold < gold {
		return GetErr("取出金币失败")
	}

	if gold > 0 {
		person.Gold += gold
		person.SaveGold -= gold
		GetTopMgr().UpdData(&person)
		GetServer().SqlGoldLog(person.Uid, gold, -11)
		person.SynchroGold()
		person.Flush(true)
	}

	var msg Msg_SaveGold
	msg.Gold = person.Gold
	msg.SaveGold = person.SaveGold
	return lib.HF_EncodeMsg("drawgold", &msg, true)
}

//! A转账B
func AGiveGoldToB(uid int64, unionid string, gold int, _uid int64) []byte {
	DrawLock.Lock()
	defer DrawLock.Unlock()

	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	if gold < 0 {
		GetServer().AddBlack(uid)
		return GetErr("转赠失败")
	}

	if uid == _uid {
		return GetErr("转赠失败")
	}

	value, _ := GetServer().DB_GetData("user", uid, 1)
	if string(value) == "" {
		return GetErr("转赠失败")
	}

	_value, _ := GetServer().DB_GetData("user", _uid, 1)
	if string(_value) == "" {
		return GetErr("转赠失败")
	}

	var person Person
	json.Unmarshal(value, &person)
	if person.UnionId != "" && person.UnionId != unionid {
		return GetErr("授权失败,请重新打开游戏")
	}

	if person.Gold < gold {
		return GetErr("金币不足")
	}

	if person.GameType != 0 { //! 任何游戏中都不允许转增
		return GetErr("转赠失败")
	}

	var _person Person
	json.Unmarshal(_value, &_person)

	if gold > 0 {
		//! A扣钱
		person.Gold -= gold
		GetTopMgr().UpdData(&person)
		GetServer().SqlGoldLog(person.Uid, -gold, -int(_person.Uid)*100)
		person.SynchroGold()
		person.Flush(true)
		//! B加钱
		_person.Gold += gold
		GetTopMgr().UpdData(&_person)
		GetServer().SqlGoldLog(_person.Uid, gold, int(person.Uid)*100)
		_person.SynchroGold()
		_person.SynchroMoney()
		_person.Flush(true)
	}
	var record Msg_Rescord_GiveGold
	record.AName = person.Name
	record.AUid = person.Uid
	record.BName = _person.Name
	record.BUid = _person.Uid
	record.Time = time.Now().Unix()
	record.Gold = gold
	GetServer().InsertGiveGoldRecord(1, uid, _uid, lib.HF_JtoA(&record))
	GetServer().InsertGiveGoldRecord(2, uid, _uid, lib.HF_JtoA(&record))

	var msg Msg_AGiveGoldToB
	msg.Auid = uid
	msg.Buid = _uid
	msg.Gold = gold
	return lib.HF_EncodeMsg("givegold", &msg, true)
}
func GetGiveGoldRecord(uid int64, _type int) []byte {
	c := GetServer().RcRedis.Get()
	defer c.Close()
	var table string
	if _type == 1 {
		table = fmt.Sprintf("give_%d", uid)
	} else {
		table = fmt.Sprintf("get_%d", uid)
	}
	values, err := redis.Values(c.Do("LRANGE", table, 0, 20))

	var msg S2C_GiveGoldRecord
	if err == nil {
		for _, v := range values {
			var record Msg_Rescord_GiveGold
			json.Unmarshal(v.([]byte), &record)
			msg.Info = append(msg.Info, record)
		}
	} else {
		msg.Info = make([]Msg_Rescord_GiveGold, 0)
	}
	return lib.HF_EncodeMsg("givegoldrecord", &msg, true)
}

//! 申请兑换
func ExchangeGold(uid int64, gold int) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	if gold < 0 {
		GetServer().AddBlack(uid)
		return GetErr("申请失败")
	}

	value, _ := GetServer().DB_GetData("user", uid, 1)
	if string(value) == "" {
		return GetErr("申请失败")
	}

	gold *= 100

	var person Person
	json.Unmarshal(value, &person)
	if person.Gold < gold {
		return GetErr("兑换金币失败")
	}

	if person.GameType != 0 { //! 任何游戏中都不允许存钱
		return GetErr("申请失败")
	}

	if gold > 0 {
		person.Gold -= gold
		GetTopMgr().UpdData(&person)
		GetServer().SqlGoldLog(person.Uid, -gold, -100)
		person.SynchroGold()
		person.Flush(true)
		GetServer().SqlExchangeLog(uid, person.Name, gold)
	}

	var msg Msg_ExchangeGold
	msg.Gold = person.Gold
	return lib.HF_EncodeMsg("exchangegold", &msg, true)
}

//! 兑换记录
func ExchangeInfo(uid int64) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	var node SQL_Exchange
	sql := fmt.Sprintf("select * from `exchange` where uid = %d order by id desc", uid)
	res := GetServer().DB.GetAllData(sql, &node)

	var msg Msg_ExchangeInfo
	for i := 0; i < len(res); i++ {
		var info Son_ExchangeInfo
		info.Gold = res[i].(*SQL_Exchange).Gold
		info.State = res[i].(*SQL_Exchange).State
		info.Time = res[i].(*SQL_Exchange).Time
		msg.Info = append(msg.Info, info)
	}
	return lib.HF_EncodeMsg("exchangeinfo", &msg, true)
}

//! 进入鱼池
func EnterFish(_type int, uid int64) []byte {
	lib.GetFishMgr().Enter(_type, uid)

	var msg S2C_GetFishDesk
	msg.Type = _type
	msg.Desk = lib.GetFishMgr().Desk[_type]
	return lib.HF_EncodeMsg("enterfish", &msg, true)
}

//! 得到房间列表
func GetRoomList(uid int64) []byte {
	var msg S2C_GetRoomList
	msg.Info = GetCreateRoomMgr().Get(uid)
	return lib.HF_EncodeMsg("getroomlist", &msg, true)
}

//! 得到战报
func Record(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][_type]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "record", &msg)
	lib.GetLogMgr().Output(lib.LOG_INFO, "返回战报:", string(result))
	return result
}

//! 记录结构
type Rec_GameKWX_Info struct {
	Roomid int                      `json:"roomid"` //! 房间号
	Time   int64                    `json:"time"`   //! 记录时间
	Person []Son_Rec_GameKWX_Person `json:"person"`
}
type Son_Rec_GameKWX_Person struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Head  string `json:"head"`
	Card  []int  `json:"card"`
	Score int    `json:"score"`
}

func RecordKWX1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][2]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordkwx1", &msg)
	return result
}

func RecordKWX2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][2]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordkwx2", &msg)
	return result
}

func RecordGC1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][12]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordgc1", &msg)
	return result
}

func RecordGC2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][12]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordgc2", &msg)
	return result
}

func RecordSZMJ1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][15]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordszmj1", &msg)
	return result
}

func RecordSZMJ2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][15]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordszmj2", &msg)
	return result
}

func RecordXZDD1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][18]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordxzdd1", &msg)
	return result
}

func RecordXZDD2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][18]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordxzdd2", &msg)
	return result
}

//常熟战报
func RecordCSMJ1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][32]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordcsmj1", &msg)
	return result
}

func RecordCSMJ2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][32]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordcsmj2", &msg)
	return result
}

//安慶战报
func RecordAQMJ1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][27]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordaqmj1", &msg)
	return result
}

func RecordAQMJ2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][27]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordaqmj2", &msg)
	return result
}

//上虞花战报
func RecordSYHMJ1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][41]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordsyhmj1", &msg)
	return result
}

func RecordSYHMJ2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][41]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordsyhmj2", &msg)
	return result
}

//通山战报
func RecordTSMJ1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][66]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordtsmj1", &msg)
	return result
}

func RecordTSMJ2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][66]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordtsmj2", &msg)
	return result
}

//南昌战报
func RecordNCMJ1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][33]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordncmj1", &msg)
	return result
}

func RecordNCMJ2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][33]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordncmj2", &msg)
	return result
}

//仙桃人人晃晃战报
func RecordXTRRHH1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][82]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordxtrrhh1", &msg)
	return result
}

func RecordXTRRHH2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][82]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordxtrrhh2", &msg)
	return result
}

//潮汕
func RecordGDCSMJ1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][81]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordgdcsmj1", &msg)
	return result
}

func RecordGDCSMJ2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][81]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordgdcsmj2", &msg)
	return result
}

//涡阳
func RecordGYMJ1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][76]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordgymj1", &msg)
	return result
}

func RecordGYMJ2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][76]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordgymj2", &msg)
	return result
}

//枪炮斗地主
func RecordQPDDZ1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][87]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordqpddz1", &msg)
	return result
}

func RecordQPDDZ2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][87]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordqpddz2", &msg)
	return result
}

//! 斗地主列表
func RecordDDZ1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][6]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordddz1", &msg)
	return result
}

//! 斗地主详细战报
func RecordDDZ2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][6]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordddz2", &msg)
	return result
}

//! 斗地主列表
func RecordBDDDZ1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][75]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordbdddz1", &msg)
	return result
}

//! 斗地主详细战报
func RecordBDDDZ2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][75]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordbdddz2", &msg)
	return result
}

//! 斗板凳列表
func RecordDBD1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][39]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recorddbd1", &msg)
	return result
}

//! 斗板凳详细战报
func RecordDBD2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][39]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recorddbd2", &msg)
	return result
}

func RecordHLGC1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][83]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordhlgc1", &msg)
	return result
}

func RecordHLGC2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][83]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordhlgc2", &msg)
	return result
}

func RecordZJH1(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][7]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordzjh1", &msg)
	return result
}

//! 扎金花详细战报
func RecordZJH2(uid int64, _type int, clientip string) []byte {
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	_, ok := staticfunc.GetCsvMgr().Data["game"][7]
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(clientip, "获取错误类型的战报")
		return []byte("")
	}

	var msg staticfunc.Msg_Record
	msg.Uid = uid
	msg.Type = _type

	config := GetServer().GetGameServerFromRoom(1)
	result, _ := config.Call("RecordMethod.RecordMsg", "recordzjh2", &msg)
	return result
}

/////////////////////////////////////////////////////////////////////////////////////////
//! 创建俱乐部房间
func CreateClubRoom(clubid int64, host int64, uid int64, gametype int, num int, param1 int, param2 int, clientip string) []byte {
	value, _ := GetServer().DB_GetData("user", host, 1)
	if string(value) == "" { //! 新用户应该先登录
		return []byte("false")
	}

	csv, ok := staticfunc.GetCsvMgr().Data["game"][gametype]
	if !ok {
		return []byte("false")
	}

	var person Person
	json.Unmarshal(value, &person)

	if num <= 0 || num > lib.HF_Atoi(csv["maxstep"]) {
		return []byte("false")
	}

	cost := num

	if person.Card < cost {
		return []byte("nocard")
	}

	//! 得到人数最少的房间
	var config *GameServerConfig = nil
	if csv["gametype"] == "2" { //! 创建麻将房
		config = GetServer().GetGameServerFromRoom(2)
	} else if csv["gametype"] == "1" { //! 创建扑克房
		config = GetServer().GetGameServerFromRoom(1)
	} else if csv["gametype"] == "3" { //！创建金币场
		config = GetServer().GetGameServerFromRoom(3)
	}

	//! 告诉游戏服务器有人玩家加入进来
	if config == nil {
		return []byte("false")
	}

	roomid := GetRoom().GetClubID()
	if roomid == 0 {
		return []byte("false")
	}

	if GetServer().CreateRoom(config.Id, roomid, gametype, num, param1, param2, host, clubid) {
		GetServer().AddGameRoom(config.Id, roomid, gametype, 1, host, clientip)
	} else {
		return []byte("false")
	}

	if lib.HF_Atoi(csv["costtype"]) == staticfunc.TYPE_GOLD {
		person.Gold = lib.HF_MaxInt(0, person.Gold-cost)
		person.Flush(true)
		GetServer().InsertLog(person.Uid, staticfunc.MOVE_GOLD, cost, clientip)
	} else {
		person.Card = lib.HF_MaxInt(0, person.Card-cost)
		person.Flush(true)
		GetServer().InsertLog(person.Uid, staticfunc.MOVE_CARD, cost, clientip)
	}

	return []byte(fmt.Sprintf("%d_%d_%d_%d_%d", roomid, cost, num*lib.HF_Atoi(csv["step"]), person.Card, person.Gold))
}

//! 设置
func SetAdmin(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "SetAdmin") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	uid := lib.HF_Atoi(req.FormValue("uid"))
	admin := lib.HF_Atoi(req.FormValue("admin"))

	value, _ := GetServer().DB_GetData("user", int64(uid), 1)
	if string(value) != "" {
		var person Person
		json.Unmarshal(value, &person)
		person.Admin = admin
		person.Flush(true)

		if person.GameId != 0 {
			config := GetServer().GetGameServer(person.GameId)
			if config != nil {
				var msg staticfunc.Msg_SetAdmin
				msg.Uid = int64(uid)
				msg.Admin = admin
				config.Call("ServerMethod.ServerMsg", "setadmin", &msg)
			}
		}

		w.Write([]byte("true"))
		return
	}
	w.Write([]byte("false"))
}

//! 载入机器人
func LoadRobot(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "LoadRobot") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	id := lib.HF_Atoi(req.FormValue("id"))

	for i := 3; ; i++ {
		config := GetServer().GetGameServer(i)
		if config != nil {
			var msg staticfunc.Msg_Uid
			msg.Uid = int64(id)
			config.Call("ServerMethod.ServerMsg", "loadrobot", &msg)
		} else {
			break
		}
	}
	w.Write([]byte("true"))
}

//! 删除机器人
func DelRobot(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	clientip := lib.HF_GetHttpIP(req)
	if !GetServer().IsWhite(clientip, "DelRobot") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	id := lib.HF_Atoi(req.FormValue("id"))

	for i := 3; ; i++ {
		config := GetServer().GetGameServer(i)
		if config != nil {
			var msg staticfunc.Msg_Uid
			msg.Uid = int64(id)
			config.Call("ServerMethod.ServerMsg", "delrobot", &msg)
		} else {
			break
		}
	}
	w.Write([]byte("true"))
}

//! 得到key
func GetKey(uid int64) string {
	return lib.HF_MD5(fmt.Sprintf("%d%s", uid, APP_KEY))
}

//! 发送重新下载
func ReDownload(_code int, _msg string, _type int) []byte {
	var msg S2C_ReDownload
	msg.Code = _code
	msg.Msg = _msg
	msg.Download = rjmgr.GetRJMgr().GetConfig(_type).DownLoad

	return lib.HF_EncodeMsg("redownload", &msg, true)
}

//! 得到底分和准入
func GetGameEnter() map[int][]GameEnter {
	result := make(map[int][]GameEnter)

	csv, ok := staticfunc.GetCsvMgr().Data["enter"]
	if !ok {
		return result
	}

	for key, value := range csv {
		if !IsGameOk(key) {
			continue
		}

		tmp := make([]GameEnter, 0)
		for i := 0; i < 6; i++ {
			tmp = append(tmp, GameEnter{lib.HF_Atoi(value[fmt.Sprintf("df%d", i)]), lib.HF_Atoi(value[fmt.Sprintf("zr%d", i)])})
		}
		result[key] = tmp
	}

	return result
}
