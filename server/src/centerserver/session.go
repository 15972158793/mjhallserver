package centerserver

import (
	"encoding/json"
	"fmt"
	"lib"
	"runtime/debug"
	"staticfunc"
)

//! 加入房间
type Msg_JoinRoom struct {
	Uid    int64 `json:"uid"` //! uid
	Roomid int   `json:"roomid"`
}

//! 加入房间失败
type Msg_JoinRoomFail struct {
	Result int `json:"result"`
}

type CodeResp struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
}

type MapInfo struct {
	Longitude string `json:"longitude"` //! 经度
	Latitude  string `json:"latitude"`  //! 纬度
	Country   string `json:"country"`   //! 国家
	Province  string `json:"province"`  //! 省
	City      string `json:"city"`      //! 市
	Citycode  string `json:"citycode"`  //! 城市编码
	District  string `json:"district"`  //! 区
	Adcode    string `json:"adcode"`    //! 区域码
	Address   string `json:"address"`   //! 地址
}

//! 得到省码
func (self *MapInfo) GetProvinceCode() string {
	if self.Adcode == "" {
		return ""
	}
	code := lib.HF_Atoi(self.Adcode)
	if code == 0 {
		return ""
	}

	code = code / 10000
	code = code * 10000
	return fmt.Sprintf("%d", code)
}

//! 得到市码
func (self *MapInfo) GetCityCode() string {
	if self.Adcode == "" {
		return ""
	}
	code := lib.HF_Atoi(self.Adcode)
	if code == 0 {
		return ""
	}

	code = code / 100
	code = code * 100
	return fmt.Sprintf("%d", code)
}

func GetErr(info string) []byte {
	var msg S2C_Err
	msg.Info = info

	return lib.HF_EncodeMsg("err", &msg, true)
}

func OnReceive(self *lib.Session, msg []byte) {
	defer func() {
		x := recover()
		if x != nil {
			lib.GetLogMgr().Output(lib.LOG_ERROR, x, string(debug.Stack()))
		}
	}()

	head, _, data, ok := lib.HF_DecodeMsg(msg)
	if !ok {
		staticfunc.GetIpBlackMgr().AddIp(self.IP, "消息解析错误1")
		return
	}

	lib.GetLogMgr().Output(lib.LOG_INFO, "client:", head, "...", string(data))

	if head != "setuid" {
		if self.Person == nil {
			return
		}
	}

	switch head {
	case "setuid":
		var msg Msg_SetUid
		json.Unmarshal(data, &msg)
		//! 先验证有效性
		//if !self.Check(msg.Uid, msg.Openid) {
		//	staticfunc.GetIpBlackMgr().AddIp(self.IP, "uid验证错误")
		//	self.SafeClose()
		//	return
		//}
		person := GetPersonMgr().GetPerson(msg.Uid, false)
		if person != nil && self != person.session && person.session != nil {
			person.CloseSession()
		}
		person = GetPersonMgr().GetPerson(msg.Uid, true)
		self.Person = person
		person.OtherPlayerData()
		person.session = self
		person.group = msg.Group
		if person.group != 1 { //! 不是扑克就是麻将
			person.group = 2
		}
		if msg.MapInfo != "" {
			json.Unmarshal([]byte(msg.MapInfo), &person.mapinfo)
		}

		//! 发送邀请
		person.GetModule("invite").(*Mod_Invite).SendInfo()
		//! 发送分享
		person.GetModule("share").(*Mod_Share).SendInfo()
		//! 发送俱乐部
		person.GetModule("club").(*Mod_Club).SendInfo()
		//! 发送转盘消息
		person.GetModule("dial").(*Mod_Dial).SendInfo()
		//! 发送救济金消息
		person.GetModule("alms").(*Mod_Alms).SendInfo()
		//! 发送签到消息
		person.GetModule("sign").(*Mod_Sign).SendInfo()
		//! 发送公告
		notice := ""
		code := person.mapinfo.Adcode
		if code != "" {
			notice = GetServer().GetAreaNotice(code, person.group)
		}
		if notice == "" {
			code = person.mapinfo.GetCityCode()
			if code != "" {
				notice = GetServer().GetAreaNotice(code, person.group)
			}
		}
		if notice == "" {
			code = person.mapinfo.GetProvinceCode()
			if code != "" {
				notice = GetServer().GetAreaNotice(code, person.group)
			}
		}

		if notice == "" {
			notice = GetServer().Notice[person.group-1]
		}
		if notice != "" {

		}
		var msg1 S2C_Notice
		msg1.Context = notice
		self.SendByteMsg(lib.HF_EncodeMsg("notice", &msg1, true))

		//! 发送地区信息
		notice = ""
		code = person.mapinfo.Adcode
		if code != "" {
			notice = GetServer().GetAreaInfo(code, person.group)
		}
		if notice == "" {
			code = person.mapinfo.GetCityCode()
			if code != "" {
				notice = GetServer().GetAreaInfo(code, person.group)
			}
		}
		if notice == "" {
			code = person.mapinfo.GetProvinceCode()
			if code != "" {
				notice = GetServer().GetAreaInfo(code, person.group)
			}
		}
		if notice == "" {
			code = "100000"
			if code != "" {
				notice = GetServer().GetAreaInfo(code, person.group)
			}
		}

		if notice == "" {
			var info HTMsg_AreaInfo
			info.WChat = "asnn999"
			notice = lib.HF_JtoA(&info)
		}

		if notice != "" {
			var msg S2C_Notice
			msg.Context = notice
			self.SendByteMsg(lib.HF_EncodeMsg("areainfo", &msg, true))
		}

		//! 发送邮件
		GetNoticeMgr().SendInfo(person)
	case "charge":
		var msg C2S_Charge
		json.Unmarshal(data, &msg)
		lib.GetChargeMgr().AddCharge(&lib.ChargeData{self.Person.(*Person).Uid, msg.Receipt, msg.Sandbox, self.IP})
	case "setmapinfo":
		var msg Msg_SetMapInfo
		json.Unmarshal(data, &msg)
		if msg.MapInfo != "" {
			json.Unmarshal([]byte(msg.MapInfo), &self.Person.(*Person).mapinfo)
		}

		//! 发送公告
		notice := ""
		code := self.Person.(*Person).mapinfo.Adcode
		if code != "" {
			notice = GetServer().GetAreaNotice(code, self.Person.(*Person).group)
		}
		if notice == "" {
			code = self.Person.(*Person).mapinfo.GetCityCode()
			if code != "" {
				notice = GetServer().GetAreaNotice(code, self.Person.(*Person).group)
			}
		}
		if notice == "" {
			code = self.Person.(*Person).mapinfo.GetProvinceCode()
			if code != "" {
				notice = GetServer().GetAreaNotice(code, self.Person.(*Person).group)
			}
		}

		if notice == "" {
			notice = GetServer().Notice[self.Person.(*Person).group-1]
		}

		if notice != "" {
			var msg S2C_Notice
			msg.Context = notice
			self.SendByteMsg(lib.HF_EncodeMsg("notice", &msg, true))
		}

		//! 发送地区信息
		notice = ""
		code = self.Person.(*Person).mapinfo.Adcode
		if code != "" {
			notice = GetServer().GetAreaInfo(code, self.Person.(*Person).group)
		}
		if notice == "" {
			code = self.Person.(*Person).mapinfo.GetCityCode()
			if code != "" {
				notice = GetServer().GetAreaInfo(code, self.Person.(*Person).group)
			}
		}
		if notice == "" {
			code = self.Person.(*Person).mapinfo.GetProvinceCode()
			if code != "" {
				notice = GetServer().GetAreaInfo(code, self.Person.(*Person).group)
			}
		}
		if notice == "" {
			code = "100000"
			if code != "" {
				notice = GetServer().GetAreaInfo(code, self.Person.(*Person).group)
			}
		}

		if notice == "" {
			var info HTMsg_AreaInfo
			info.WChat = "asnn999"
			notice = lib.HF_JtoA(&info)
		}

		if notice != "" {
			var msg S2C_Notice
			msg.Context = notice
			self.SendByteMsg(lib.HF_EncodeMsg("areainfo", &msg, true))
		}
	case "readmail":
		var msg Msg_ReadMail
		json.Unmarshal(data, &msg)
		if self.Person.(*Person).group == 1 {
			GetNoticeMgr().ReadMail(self.Person.(*Person), msg.Id)
		} else if self.Person.(*Person).group == 2 {
			GetNoticeMgr().ReadKWXMail(self.Person.(*Person), msg.Id)
		}
	case "report": //! 举报
		var msg Msg_Report
		json.Unmarshal(data, &msg)
		GetServer().SqlReport(self.Person.(*Person).Uid, msg.Uid, msg.Type, msg.Dec)
	default:
		self.Person.(*Person).Module.OnMsg(head, data)
	}
}

func OnClose(self *lib.Session) {
	if self.Person == nil {
		return
	}

	self.Person.(*Person).session = nil
	self.Person = nil
}
