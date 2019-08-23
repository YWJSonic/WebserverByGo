package mycache

import (
	"fmt"
	"time"

	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/redigo/redis"
	"gitlab.com/WeberverByGoGame9/serversetting"
)

var CachePool *redis.Pool

func init() {
	newCachePool()
}

const ConnectTimeout time.Duration = 20 * time.Second
const ReadTimeout time.Duration = 5 * time.Second
const WriteTimeout time.Duration = 10 * time.Second

func newCachePool() {
	CachePool = &redis.Pool{
		MaxIdle:     50,
		IdleTimeout: 240 * time.Second,
		MaxActive:   50,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", serversetting.RedisURL,
				redis.DialConnectTimeout(ConnectTimeout),
				redis.DialReadTimeout(ReadTimeout),
				redis.DialWriteTimeout(WriteTimeout))
			if err != nil {
				return nil, fmt.Errorf("redis connection error: %s", err)
			}
			//验证redis密码
			// if _, authErr := c.Do("AUTH", RedisPassword); authErr != nil {
			// 	return nil, fmt.Errorf("redis auth password error: %s", authErr)
			// }
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}
}

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
func GetAccountInfo(gameAccount string) (interface{}, messagehandle.ErrorMsg) {
	err := messagehandle.New()
	info, errMsg := get(fmt.Sprintf("ACC%s", gameAccount))

	if errMsg != nil {
		messagehandle.ErrorLogPrintln("Cache GetAccountInfo", gameAccount, errMsg)
		err.ErrorCode = code.FailedPrecondition
		err.Msg = fmt.Sprintln(errMsg)
		return nil, err
	}

	return info, err
}

// SetAccountInfo Set Account Struct
func SetAccountInfo(gameAccount string, Value interface{}) {
	runSet(fmt.Sprintf("ACC%s", gameAccount), Value, serversetting.CacheDeleteTime)
}

// GetPlayerInfo Get PlayerInfo Struct
func GetPlayerInfo(playerid int64) (interface{}, messagehandle.ErrorMsg) {
	err := messagehandle.New()
	info, errMsg := get(fmt.Sprintf("ID%dJS", playerid))

	if errMsg != nil {
		messagehandle.ErrorLogPrintln("Cache GetPlayerInfo", playerid, errMsg)
		err.ErrorCode = code.FailedPrecondition
		err.Msg = fmt.Sprintln(errMsg)
		return nil, err
	}

	return info, err
}

// SetPlayerInfo Set PlayerInfo Struct
func SetPlayerInfo(playerid int64, Value interface{}) {
	runSet(fmt.Sprintf("ID%dJS", playerid), Value, serversetting.CacheDeleteTime)
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
	runSet(key, value, serversetting.CacheDeleteTime)
}

// GetULGInfoCache Get ULG info
func GetULGInfoCache(playerid int64) interface{} {
	err := messagehandle.New()
	key := fmt.Sprintf("ULG%d", playerid)
	info, errMsg := get(key)

	if errMsg != nil {
		messagehandle.ErrorLogPrintln("Cache GetULGInfoCache", key)
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
	runSet(key, value, serversetting.CacheDeleteTime)
}

// GetAttach game data request
func GetAttach(playerid int64) interface{} {
	err := messagehandle.New()
	key := fmt.Sprintf("attach%d", playerid)
	info, errMsg := get(key)

	if errMsg != nil {
		messagehandle.ErrorLogPrintln("Cache GetULGInfoCache", key)
		err.ErrorCode = code.FailedPrecondition
		err.Msg = fmt.Sprintln(errMsg)
		return nil
	}

	return info
}
