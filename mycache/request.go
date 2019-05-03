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
	runSet(GameAccount, Token, lasttime)
}

func GetPlayerInfo(playerid int64) (interface{}, errorlog.ErrorMsg) {
	err := errorlog.New()
	info, errMsg := get(fmt.Sprintf("ID%dJS", playerid))

	if errMsg != nil {
		fmt.Println(errMsg)
		err.ErrorCode = code.FailedPrecondition
		err.Msg = fmt.Sprintln(errMsg)
		return nil, err
	}

	return info, err
}
func SetPlayerInfo(playerid int64, Value interface{}) {
	runSet(fmt.Sprintf("ID%dJS", playerid), Value, data.CacheDeleteTime)
}

func GetPlayerID(GameAccount string) int64 {
	id, err := getInt64(GameAccount)
	if err != nil {
		panic("GetPlayerID Error")
	}
	return id
}
func SetPlayerID(GameAccount string, playerID int64) {
	runSet(GameAccount, playerID, data.CacheDeleteTime)
}
