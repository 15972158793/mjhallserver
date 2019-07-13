package backstage

import (
	"rjmgr"
)

func BackStageInit(flag string, sign string) bool {
	return rjmgr.GetRJMgr().DoLicensing(flag, sign)
}
