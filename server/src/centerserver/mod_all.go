package centerserver

type Mod_All struct {
	person *Person
	Module map[string]Mod_Base
}

func NewModAll(person *Person) *Mod_All {
	p := new(Mod_All)
	p.person = person
	p.Module = make(map[string]Mod_Base)

	//! 这里注册
	p.Module["share"] = new(Mod_Share)
	p.Module["invite"] = new(Mod_Invite)
	//p.Module["friend"] = new(Mod_Friend)
	//p.Module["gift"] = new(Mod_Gift)
	//p.Module["clothes"] = new(Mod_Clothes)
	//p.Module["record"] = new(Mod_Record)
	//p.Module["task"] = new(Mod_Task)
	//p.Module["packs"] = new(Mod_Packs)
	p.Module["club"] = new(Mod_Club)
	//p.Module["bank"] = new(Mod_Bank)
	p.Module["dial"] = new(Mod_Dial)
	p.Module["alms"] = new(Mod_Alms)
	p.Module["sign"] = new(Mod_Sign)

	return p
}

func (self *Mod_All) GetModule(name string) Mod_Base {
	return self.Module[name]
}

//! 得到数据
func (self *Mod_All) GetData() {
	for _, value := range self.Module {
		value.OnGetData(self.person)
	}
}

//! 得到数据
func (self *Mod_All) GetOtherData() {
	for _, value := range self.Module {
		value.OnGetOtherData()
	}
}

//! 保存数据
func (self *Mod_All) Save(sql bool) {
	for _, value := range self.Module {
		value.OnSave(sql)
	}
}

//! 得到消息
func (self *Mod_All) OnMsg(head string, body []byte) {
	for _, value := range self.Module {
		if value.OnMsg(head, body) {
			break
		}
	}
}
