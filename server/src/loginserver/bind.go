package loginserver

import (
	"encoding/json"
	"fmt"
	"lib"
	//"net/http"
	//"runtime/debug"
	//"strings"
	"io/ioutil"
	"rjmgr"
	"sync"
	"time"
)

type BindParentNode struct {
	IP   string
	UID  int64
	Time int64
}

var BindParent []BindParentNode

var BindLock *sync.RWMutex = new(sync.RWMutex)

func AddBindNode(ip string, uid int64, t int64) {
	BindLock.Lock()
	defer BindLock.Unlock()

	BindParent = append(BindParent, BindParentNode{ip, uid, t})
	lib.GetLogMgr().Output(lib.LOG_DEBUG, "加一个ip:", ip, ",", uid)
}

func GetBindNode(ip string) BindParentNode {
	BindLock.Lock()
	defer BindLock.Unlock()

	var node BindParentNode
	for i := 0; i < len(BindParent); {
		if time.Now().Unix()-BindParent[i].Time >= 86400 { //! 这个绑定信息已经超时
			copy(BindParent[i:], BindParent[i+1:])
			BindParent = BindParent[:len(BindParent)-1]
			continue
		}
		if BindParent[i].IP == ip {
			node = BindParent[i]
			copy(BindParent[i:], BindParent[i+1:])
			BindParent = BindParent[:len(BindParent)-1]
			lib.GetLogMgr().Output(lib.LOG_DEBUG, "得到一个ip:", node.IP, ",", node.UID)
			break
		}
		i++
	}
	return node
}

///////////////////////////////////////////////////////////
type C2S_Register struct {
	Account  string `json:"account"`  //！账号
	Passwd   string `json:"passwd"`   //！密码
	NickName string `json:"nickname"` //！昵称
}

type S2C_Register struct {
	Account string `json:"account"`
	Passwd  string `json:"passwd"`
}

type C2S_ChgPasswd struct {
	Account    string `json:"account"`    //！账号
	Passwd     string `json:"passwd"`     //！密码
	NewPasswd  string `json:"newpasswd"`  //！新密码
	NewPasswd1 string `json:"newpasswd1"` //！新密码
}

type C2S_ModifyHead struct {
	Uid     int64  `json:"uid"`
	Head    string `json:"head"`
	Unionid string `json:"unionid"`
}

type S2C_ModifyHead struct {
	Head string `json:"head"`
}

//! 注册
func Register(account string, passwd string, nickname string, clientip string) []byte {
	if account == "" || passwd == "" {
		return GetErr("账号密码不能为空")
	}

	if nickname == "" {
		return GetErr("请输入昵称")
	}

	if len([]byte(account)) < 5 {
		return GetErr("账号至少5个字符")
	}

	if len([]byte(passwd)) < 5 {
		return GetErr("密码至少5个字符")
	}

	if !lib.HF_IsLicitAccount([]byte(account)) {
		return GetErr("账号只能输入字母、数字、下划线和@")
	}

	if !lib.HF_IsLicitPasswd([]byte(passwd)) {
		return GetErr("密码包含非法字符,请重新输入")
	}

	if !lib.HF_IsLicitName([]byte(nickname)) {
		return GetErr("昵称输入不合法,请重新输入")
	}

	md5_passwd := lib.HF_MD5(passwd)

	uid := GetServer().CheckAccountAndPasswd(account, md5_passwd)
	if uid != 0 {
		return GetErr("该账号已存在")
	}

	uid = GetServer().DB_GetUid(account, md5_passwd, true, 0)
	if uid == 0 {
		return GetErr("当前无法注册，请稍后再试")
	}

	var person Person
	value, db := GetServer().DB_GetData("user", uid, 2)
	if string(value) == "" { //! 一个新用户
		person.Uid = uid
		person.Card = GetServer().Con.NewCard
		person.Gold = GetServer().InitMoney
		person.BindGold = GetServer().InitMoney
		person.Name = nickname
		person.Imgurl = fmt.Sprintf("%d", lib.HF_GetRandom(10)+1)
		person.Sex = 1
		person.Time = time.Now().Unix()
		person.IP = clientip
		person.UnionId = account
		person.OpenId = md5_passwd
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
		person.Time = time.Now().Unix()
		person.IP = clientip
		person.UnionId = account
		person.OpenId = md5_passwd
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

	var msg S2C_Register
	msg.Account = account
	msg.Passwd = passwd
	return lib.HF_EncodeMsg("register", &msg, true)
}

//! 修改密码
func ChgPasswd(account string, passwd string, newpasswd string, newpasswd1 string, clientip string) []byte {
	if account == "" || passwd == "" {
		return GetErr("账号密码错误")
	}

	if newpasswd == "" || newpasswd1 == "" {
		return GetErr("请输入新密码")
	}

	if newpasswd != newpasswd1 {
		return GetErr("两次密码输入不一致")
	}

	if passwd == newpasswd {
		return GetErr("新密码不能与旧密码相同")
	}

	if len([]byte(newpasswd)) < 5 {
		return GetErr("密码至少5个字符")
	}

	if len([]byte(passwd)) < 5 {
		return GetErr("密码错误")
	}

	if !lib.HF_IsLicitAccount([]byte(account)) {
		return GetErr("密码错误")
	}

	if !lib.HF_IsLicitPasswd([]byte(passwd)) {
		return GetErr("密码错误")
	}

	if !lib.HF_IsLicitName([]byte(newpasswd)) {
		return GetErr("密码错误")
	}

	md5_passwd := lib.HF_MD5(passwd)

	uid := GetServer().CheckAccountAndPasswd(account, md5_passwd)
	if uid == 0 {
		return GetErr("该账号不存在")
	} else if uid == -1 {
		return GetErr("密码错误")
	}

	md5_passwd = lib.HF_MD5(newpasswd)

	GetServer().DB_ModifyPasswd(account, md5_passwd)

	var msg S2C_Register
	msg.Account = account
	msg.Passwd = newpasswd
	return lib.HF_EncodeMsg("chgpasswd", &msg, true)
}

//! 账号密码登陆
func LoginFromPasswd(account string, passwd string, _type int, ver int, assetkey string, clientip string) []byte {
	if ver > GetServer().Con.Version { //! 客户端版本过高
		return GetRet("loginfail", 1)
	} else if ver < GetServer().Con.MinVersion { //! 版本过低
		return ReDownload(0, "当前版本号过低，请下载最新版本进行游戏", 0)
	}

	if account == "" || passwd == "" {
		return GetErr("账号密码不能为空")
	}

	if !lib.HF_IsLicitAccount([]byte(account)) || !lib.HF_IsLicitAccount([]byte(passwd)) {
		return GetErr("账号密码输入不合法")
	}

	//! 不允许游客登陆,就判断assetkey
	if rjmgr.GetRJMgr().Guest == 0 && !rjmgr.GetRJMgr().AssetKeyTrue(assetkey) {
		return GetErr("您的版本过低,请下载最新版本")
	}

	passwd = lib.HF_MD5(passwd)

	uid := GetServer().CheckAccountAndPasswd(account, passwd)
	if uid == 0 {
		return GetErr("账号不存在")
	} else if uid == -1 {
		return GetErr("密码错误")
	}
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}
	var person Person
	value, db := GetServer().DB_GetData("user", uid, 2)
	if string(value) == "" { //! 一个新用户
		person.Uid = uid
		person.Card = GetServer().Con.NewCard
		person.Gold = GetServer().InitMoney
		person.BindGold = GetServer().InitMoney
		person.Name = fmt.Sprintf("玩家%d", uid)
		person.Imgurl = fmt.Sprintf("%d", lib.HF_GetRandom(10)+1)
		person.Sex = 1
		person.Time = time.Now().Unix()
		person.IP = clientip
		person.UnionId = account
		person.OpenId = passwd
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
		person.Time = time.Now().Unix()
		person.IP = clientip
		person.UnionId = account
		person.OpenId = passwd
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

//! 游客登陆
func LoginFromYK(_type int, ver int, assetkey string, clientip string) []byte {
	if ver > GetServer().Con.Version { //! 客户端版本过高
		return GetRet("loginfail", 1)
	} else if ver < GetServer().Con.MinVersion { //! 版本过低
		return ReDownload(0, "当前版本号过低，请下载最新版本进行游戏", 0)
	}

	//! 不允许游客登陆,就判断assetkey
	if rjmgr.GetRJMgr().Guest == 0 && !rjmgr.GetRJMgr().AssetKeyTrue(assetkey) {
		return GetErr("您的版本过低,请下载最新版本")
	}

	openid := GetGuid()
	uid := GetServer().DB_GetUid(openid, "rjyklogin", true, 0)
	if GetServer().IsBlack(uid) {
		return GetErr("由于您被玩家多次举报，经核实，已封停该账号。如需申述，请联系微信客服")
	}

	var person Person
	value, db := GetServer().DB_GetData("user", uid, 2)
	if string(value) == "" { //! 一个新用户
		person.Uid = uid
		person.Card = GetServer().Con.NewCard
		person.Gold = GetServer().InitMoney
		person.BindGold = GetServer().InitMoney
		person.Name = fmt.Sprintf("玩家%d", uid)
		person.Imgurl = fmt.Sprintf("%d", lib.HF_GetRandom(10)+1)
		person.Sex = 1
		person.Time = time.Now().Unix()
		person.IP = clientip
		person.UnionId = openid
		person.OpenId = "rjyklogin"
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
		person.Time = time.Now().Unix()
		person.IP = clientip
		person.UnionId = openid
		person.OpenId = "rjyklogin"
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

//! 修改密码
func ModifyHead(uid int64, unionid string, head string) []byte {
	value, _ := GetServer().DB_GetData("user", uid, 1)
	if string(value) == "" {
		return GetErr("找不到用户")
	}

	var person Person
	json.Unmarshal(value, &person)

	if person.UnionId != "" && person.UnionId != unionid {
		return GetErr("授权失败,请重新打开游戏")
	}

	find := false
	for i := 1; i <= 10; i++ {
		if head == fmt.Sprintf("%d", i) {
			find = true
			break
		}
	}
	if !find {
		return GetErr("头像错误")
	}

	person.Imgurl = head
	person.Flush(true)
	GetTopMgr().UpdBaseInfo(&person)

	var msg S2C_ModifyHead
	msg.Head = head
	return lib.HF_EncodeMsg("modifyhead", &msg, true)
}
