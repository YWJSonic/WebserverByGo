package mycache

import (
	"fmt"

	"../code"
	"../data"
	"../messagehandle/errorlog"
)

// SetULGInfo Set ULG info
func SetULGInfo(key, value string) {
	runSet(key, value, data.CacheDeleteTime)
}

// GetULGInfoCache Get ULG info
func GetULGInfoCache(key string) (string, errorlog.ErrorMsg) {
	err := errorlog.New()

	info, errMsg := getString(key)

	if errMsg != nil {
		errorlog.ErrorLogPrintln("Cache", key, errMsg)
		err.ErrorCode = code.FailedPrecondition
		err.Msg = fmt.Sprintln(errMsg)
		return "", err
	}

	return info, err
}
