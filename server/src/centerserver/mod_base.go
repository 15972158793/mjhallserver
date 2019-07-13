package centerserver

import ()

type Mod_Base interface {
	OnGetData(person *Person)            //! 得到数据,缓存时会读取
	OnGetOtherData()                     //! 得到数据,缓存时不会读取
	OnMsg(head string, body []byte) bool //! 模块收到消息
	OnSave(sql bool)                     //! 保存数据
}
