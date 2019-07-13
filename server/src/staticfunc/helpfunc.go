package staticfunc

import (
	"lib"
	"net"
)

//! 超级游戏
const ADMIN_NIUNIU = 1
const ADMIN_SZP = 2
const ADMIN_DDZ = 4
const ADMIN_KWX = 8
const ADMIN_PTJ = 16
const ADMIN_NN30 = 32
const ADMIN_NN50 = 64
const ADMIN_NN80 = 128
const ADMIN_NN100 = 256
const ADMIN_PDK = 512
const ADMIN_GOLDBZW = 1024
const ADMIN_GOLDLHD = 2048
const ADMIN_TENHALF = 4096
const ADMIN_DWDYJLB = 8192
const ADMIN_DWJDQS = 16384

//! 本机调试改为1
var LIBDEBUG = 0

func HF_GetMacAddress() bool {
	macAddress := ""
	netInterfaces, _ := net.Interfaces()
	for i := range netInterfaces {
		netInterface := netInterfaces[i]
		hardwareAddr := netInterface.HardwareAddr.String()
		if hardwareAddr != "" {
			macAddress = hardwareAddr
			break
		}
	}
	lib.GetLogMgr().Output(lib.LOG_INFO, macAddress)
	return macAddress == "00:16:3e:00:73:2d"
}
