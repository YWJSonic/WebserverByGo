package mycache

import (
	"fmt"
	"time"

	"../code"
	"../data"
	"../messagehandle/errorlog"
)

// SetToken ...
func SetToken(GameAccount, Token string) {
	now := time.Now()
	lastHour := 23 - now.Hour()
	lastMinute := 59 - now.Minute()
	lastsecod := 60 - now.Second()
	lasttime := time.Duration(lastHour*60*60+lastMinute*60+lastsecod) * time.Second
	runSet(fmt.Sprintf("TOK%s", GameAccount), Token, lasttime)
}

// GetToken ...
func GetToken(GameAccount string) string {
	value, err := getString(fmt.Sprintf("TOK%s", GameAccount))

	if err != nil {
		return ""
	}
	return value
}

// GetAccountInfo Get Account Struct
func GetAccountInfo(GameAccount string) (interface{}, errorlog.ErrorMsg) {
	err := errorlog.New()
	info, errMsg := get(fmt.Sprintf("ACC%s", GameAccount))

	if errMsg != nil {
		errorlog.ErrorLogPrintln("Cache GetAccountInfo", GameAccount, errMsg)
		err.ErrorCode = code.FailedPrecondition
		err.Msg = fmt.Sprintln(errMsg)
		return nil, err
	}

	return info, err
}

// SetAccountInfo Set Account Struct
func SetAccountInfo(GameAccount string, Value interface{}) {
	runSet(fmt.Sprintf("ACC%s", GameAccount), Value, data.CacheDeleteTime)
}

// GetPlayerInfo Get PlayerInfo Struct
func GetPlayerInfo(playerid int64) (interface{}, errorlog.ErrorMsg) {
	err := errorlog.New()
	info, errMsg := get(fmt.Sprintf("ID%dJS", playerid))

	if errMsg != nil {
		errorlog.ErrorLogPrintln("Cache GetPlayerInfo", playerid, errMsg)
		err.ErrorCode = code.FailedPrecondition
		err.Msg = fmt.Sprintln(errMsg)
		return nil, err
	}

	return info, err
}

// SetPlayerInfo Set PlayerInfo Struct
func SetPlayerInfo(playerid int64, Value interface{}) {
	runSet(fmt.Sprintf("ID%dJS", playerid), Value, data.CacheDeleteTime)
}

// ClearPlayerCache ...
func ClearPlayerCache(playerid int64, GameAccount string) {
	del(fmt.Sprintf("ID%dJS", playerid))
	del(fmt.Sprintf("ACC%s", GameAccount))
	del(fmt.Sprintf("TOK%s", GameAccount))
}

// ClearAllCache ...
func ClearAllCache() {
	runFlush()
}

// third party request

// SetULGInfo Set ULG info
func SetULGInfo(key string, value interface{}) {
	runSet(key, value, data.CacheDeleteTime)
}

// GetULGInfoCache Get ULG info
func GetULGInfoCache(key string) interface{} {
	err := errorlog.New()
	info, errMsg := get(key)

	if errMsg != nil {
		errorlog.ErrorLogPrintln("Cache GetULGInfoCache", key)
		err.ErrorCode = code.FailedPrecondition
		err.Msg = fmt.Sprintln(errMsg)
		return nil
	}

	return info
}
