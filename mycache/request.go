package mycache

import (
	"fmt"
	"time"

	"gitlab.com/WeberverByGo/code"
	"gitlab.com/WeberverByGo/data"
	"gitlab.com/WeberverByGo/messagehandle/errorlog"
)

// SetToken ...
func SetToken(gameAccount, Token string) {
	now := time.Now()
	lastHour := 23 - now.Hour()
	lastMinute := 59 - now.Minute()
	lastsecod := 60 - now.Second()
	lasttime := time.Duration(lastHour*60*60+lastMinute*60+lastsecod) * time.Second
	runSet(fmt.Sprintf("TOK%s", gameAccount), Token, lasttime)
}

// GetToken ...
func GetToken(gameAccount string) string {
	value, err := getString(fmt.Sprintf("TOK%s", gameAccount))

	if err != nil {
		return ""
	}
	return value
}

// GetAccountInfo Get Account Struct
func GetAccountInfo(gameAccount string) (interface{}, errorlog.ErrorMsg) {
	err := errorlog.New()
	info, errMsg := get(fmt.Sprintf("ACC%s", gameAccount))

	if errMsg != nil {
		errorlog.ErrorLogPrintln("Cache GetAccountInfo", gameAccount, errMsg)
		err.ErrorCode = code.FailedPrecondition
		err.Msg = fmt.Sprintln(errMsg)
		return nil, err
	}

	return info, err
}

// SetAccountInfo Set Account Struct
func SetAccountInfo(gameAccount string, Value interface{}) {
	runSet(fmt.Sprintf("ACC%s", gameAccount), Value, data.CacheDeleteTime)
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
func ClearPlayerCache(playerid int64, gameAccount string) {
	del(fmt.Sprintf("ID%dJS", playerid))
	del(fmt.Sprintf("ACC%s", gameAccount))
	del(fmt.Sprintf("TOK%s", gameAccount))
}

// ClearAllCache ...
func ClearAllCache() {
	runFlush()
}

// third party request

// SetULGInfo Set ULG info
func SetULGInfo(playerid int64, value interface{}) {
	key := fmt.Sprintf("ULG%d", playerid)
	runSet(key, value, data.CacheDeleteTime)
}

// GetULGInfoCache Get ULG info
func GetULGInfoCache(playerid int64) interface{} {
	err := errorlog.New()
	key := fmt.Sprintf("ULG%d", playerid)
	info, errMsg := get(key)

	if errMsg != nil {
		errorlog.ErrorLogPrintln("Cache GetULGInfoCache", key)
		err.ErrorCode = code.FailedPrecondition
		err.Msg = fmt.Sprintln(errMsg)
		return nil
	}

	return info
}

// game info per each player

// SetAttach ...
func SetAttach(playerid int64, value interface{}) {
	key := fmt.Sprintf("attach%d", playerid)
	runSet(key, value, data.CacheDeleteTime)
}

// GetAttach game data request
func GetAttach(playerid int64) interface{} {
	err := errorlog.New()
	key := fmt.Sprintf("attach%d", playerid)
	info, errMsg := get(key)

	if errMsg != nil {
		errorlog.ErrorLogPrintln("Cache GetULGInfoCache", key)
		err.ErrorCode = code.FailedPrecondition
		err.Msg = fmt.Sprintln(errMsg)
		return nil
	}

	return info
}
