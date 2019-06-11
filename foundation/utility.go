package foundation

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"gitlab.com/WeberverByGo/code"
	"gitlab.com/WeberverByGo/data"
	"gitlab.com/WeberverByGo/messagehandle/errorlog"
	"gitlab.com/WeberverByGo/mycache"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// DeleteArrayElement ...
func DeleteArrayElement(elementIndex interface{}, array []interface{}) []interface{} {
	count := len(array)
	for index := 0; index < count; index++ {
		if elementIndex == array[index] {
			return append(array[:index], array[index+1:]...)
		}
	}
	return array
}

// JSONToString conver JsonStruct to JsonString
func JSONToString(v interface{}) (out string) {
	str, err := json.Marshal(v)
	if err != nil {
		return
	}

	out = string(str)
	return
}

// InterfaceTofloat64 ...
func InterfaceTofloat64(v interface{}) float64 {
	return v.(float64)
}

// InterfaceToInt ...
func InterfaceToInt(v interface{}) int {
	switch v.(type) {
	case float64:
		return int(InterfaceTofloat64(v))
	case int:
		return v.(int)
	case int64:
		return int(v.(int64))
	default:
		panic("Conver Error")
	}
}

// InterfaceToInt64 ...
func InterfaceToInt64(v interface{}) int64 {
	switch v.(type) {
	case float64:
		return int64(v.(float64))
	case int:
		return int64(v.(int))
	case int64:
		return v.(int64)
	default:
		errorlog.ErrorLogPrintln("Conver", v)
		panic("Conver Error")
	}
}

// InterfaceToBool ...
func InterfaceToBool(v interface{}) bool {
	switch v.(type) {
	case int:
		return v.(bool)
	case bool:
		return v.(bool)
	default:
		panic("Conver Error")
	}
}

// InterfaceToDynamicInt ...
func InterfaceToDynamicInt(v interface{}) code.Code {
	return code.Code(InterfaceTofloat64(v))
}

// InterfaceToString ...
func InterfaceToString(v interface{}) string {
	return v.(string)
}

// NewAccount convert all plant account to server account
func NewAccount(plant, account string) string {
	return fmt.Sprintf("%s:%s", plant, account)
}

// NewGameAccount new game account
func NewGameAccount(account string) string {
	return MD5Code(data.AccountEncodeStr + account)
}

// NewToken ...
func NewToken(gameAccount string) string {
	return MD5Code(fmt.Sprintf("%s%d", gameAccount, ServerNowTime()))
}

// CheckToken Check Token func
func CheckToken(gameAccount, token string) errorlog.ErrorMsg {
	err := errorlog.New()
	ServerToken := mycache.GetToken(gameAccount)
	if ServerToken != token {
		err.ErrorCode = code.Unauthenticated
		err.Msg = "TokenError"
	}
	return err
}

// CheckGameType Check Game Type
func CheckGameType(gameTypeID string) errorlog.ErrorMsg {
	err := errorlog.New()
	if gameTypeID != data.GameTypeID {
		err.ErrorCode = code.GameTypeError
		err.Msg = "GameTypeError"
	}
	return err
}

// MD5Code encode MD5
func MD5Code(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

// RangeRandom array random index
func RangeRandom(rangeInt []int) int {
	Sum := 0

	for _, value := range rangeInt {
		Sum += value
	}

	random := rand.Intn(Sum)

	Sum = 0
	for i, value := range rangeInt {
		Sum += value
		if Sum > random {
			return i
		}
	}
	return -1

}

// ServerNowTime Get now Unix time
func ServerNowTime() int64 {
	return time.Now().Unix()
}

// ServerNow Get now time
func ServerNow() time.Time {
	return time.Now()
}

// ConevrToTimeInt64 Get time point
func ConevrToTimeInt64(year int, month time.Month, day, hour, min, sec, nsec int) int64 {
	return time.Date(year, month, day, hour, min, sec, nsec, time.Local).Unix()
}

// AppendMap map append map
func AppendMap(Target map[string]interface{}, Source map[string]interface{}) map[string]interface{} {
	for Key, Value := range Source {
		Target[Key] = Value
	}
	return Target
}
