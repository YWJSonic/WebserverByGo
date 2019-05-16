package foundation

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"../code"
)

// DeleteArrayElement ...
func DeleteArrayElement(ElementIndex interface{}, array []interface{}) []interface{} {
	count := len(array)
	for index := 0; index < count; index++ {
		if ElementIndex == array[index] {
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
	return int(InterfaceTofloat64(v))
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

// NewToken ...
func NewToken(GameAccount string) string {
	return MD5Code(fmt.Sprintf("%s%d", GameAccount, ServerNowTime))
}

// MD5Code encode MD5
func MD5Code(Data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(Data)))
}

// RangeRandom array random index
func RangeRandom(Range []int) int {
	Sum := 0
	rand.Seed(time.Now().UnixNano())

	for _, value := range Range {
		Sum += value
	}

	random := rand.Intn(Sum)

	Sum = 0
	for i, value := range Range {
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
