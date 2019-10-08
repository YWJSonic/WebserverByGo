package setting

import (
	"gitlab.com/WeberverByGoGame5/crontab"
	"gitlab.com/WeberverByGoGame5/data"
	"gitlab.com/WeberverByGoGame5/foundation"
)

// ServerTime ...
func serverTime() int64 {
	return foundation.ServerNowTime()
}

// MaintainTime ...
func maintainTime() int64 {
	return crontab.SpecToTime(data.MaintainStartTime).Unix()
}

// New ...
func New() map[string]interface{} {

	setting := map[string]interface{}{
		"servertime":   serverTime(),
		"maintaintime": maintainTime(),
	}
	return setting
}
