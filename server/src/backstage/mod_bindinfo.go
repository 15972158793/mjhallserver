package backstage

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"lib"
	"log"
	"net/http"
	"rjmgr"
)

type SQL_Money struct {
	Money int
}

type SQL_MoneyF struct {
	Money float64
}

//! 绑定手机号
func BindPhone(w http.ResponseWriter, req *http.Request) {
	if !rjmgr.GetRJMgr().IsLicensing() { //! 未授权
		return
	}

	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))
	phone := req.FormValue("phone")
	if uid == 0 || phone == "" {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	if HasPhone(phone) {
		w.Write(HF_Code(3, "手机号已被占用"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	agent.Phone = phone
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `phone` = '%s' where `agid` = %d", agent.Phone, agent.Agid), []byte(""), false)

	phonekey := fmt.Sprintf("phone_%s", phone)
	c := GetServer().Redis.Get()
	defer c.Close()
	c.Do("SET", phonekey, 1)

	w.Write(HF_Code(0, "成功"))
}

//! 绑定支付宝
func BindAlipay(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))
	alipay := req.FormValue("alipay")
	aliname := req.FormValue("aliname")
	if uid == 0 || alipay == "" || aliname == "" {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	agent.Alipay = alipay
	agent.AliName = aliname
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `alipay` = '%s', `aliname` = '%s' where `agid` = %d", agent.Alipay, agent.AliName, agent.Agid), []byte(""), false)

	w.Write(HF_Code(0, "成功"))
}

//! 清空支付宝
func ClearAlipay(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))
	if uid == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	agent.Alipay = ""
	agent.AliName = ""
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `alipay` = '%s', `aliname` = '%s' where `agid` = %d", agent.Alipay, agent.AliName, agent.Agid), []byte(""), false)

	w.Write(HF_Code(0, "成功"))
}

//! 绑定银行卡
func BindBank(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))
	name := req.FormValue("name")
	bankcard := req.FormValue("bankcard")
	bankname := req.FormValue("bankname")
	if uid == 0 || name == "" || bankcard == "" || bankname == "" {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	agent.Name = name
	agent.Bankcard = bankcard
	agent.Bankname = bankname
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `name` = '%s', `bankname` = '%s', `bankcard` = '%s' where `agid` = %d", agent.Name, agent.Bankname, agent.Bankcard, agent.Agid), []byte(""), false)

	w.Write(HF_Code(0, "成功"))
}

//! 清空银行卡
func ClearBank(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	uid := lib.HF_Atoi(req.FormValue("uid"))
	if uid == 0 {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	agent := GetServer().GetAgentUser(uid)
	if agent.Agid == 0 {
		w.Write(HF_Code(2, "无法获取用户信息"))
		return
	}

	agent.Name = ""
	agent.Bankcard = ""
	agent.Bankname = ""
	GetServer().SetAgentUser(agent)
	GetServer().QueueAgentUser(fmt.Sprintf("update `fa_agent_user` set `name` = '%s', `bankname` = '%s', `bankcard` = '%s' where `agid` = %d", agent.Name, agent.Bankname, agent.Bankcard, agent.Agid), []byte(""), false)

	w.Write(HF_Code(0, "成功"))
}

func OnlyPhone(w http.ResponseWriter, req *http.Request) {
	if !GetServer().IsWhite(lib.HF_GetHttpIP(req)) {
		return
	}

	if GetServer().ShutDown {
		return
	}

	GetServer().Wait.Add(1)
	defer GetServer().Wait.Done()

	phone := req.FormValue("phone")
	log.Println(phone)
	if phone == "" {
		w.Write(HF_Code(1, "参数错误"))
		return
	}

	if HasPhone(phone) {
		w.Write(HF_Code(0, "true"))
	} else {
		w.Write(HF_Code(0, "false"))
	}
}

//! 判断手机号是否重复
func HasPhone(phone string) bool {
	c := GetServer().Redis.Get()
	defer c.Close()
	phonekey := fmt.Sprintf("phone_%s", phone)
	value, err := redis.Int(c.Do("GET", phonekey))
	if err != nil || value == 0 {
		var money SQL_Money
		GetServer().HT_DB.GetOneData(fmt.Sprintf("select `agid` from fa_agent_user where `phone` = '%s'", phone), &money)
		if money.Money == 0 {
			return false
		}

		c.Do("SET", phonekey, 1)
	}

	return true
}
