package foundation

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"gitlab.com/WeberverByGo/code"
	"gitlab.com/WeberverByGo/messagehandle/errorlog"
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
func JSONToString(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", " ")
	STR := string(data)
	STR = strings.ReplaceAll(STR, string(10), ``)
	return STR
}

// StringToJSON ...
func StringToJSON(jsStr string) map[string]interface{} {
	return ByteToJSON([]byte(jsStr))
}

// ByteToJSON ...
func ByteToJSON(jsByte []byte) map[string]interface{} {
	var data map[string]interface{}
	if errMsg := json.Unmarshal(jsByte, &data); errMsg != nil {
		panic(errMsg)
	}

	return data
}

// ToJSONStr Convert to json string
func ToJSONStr(data interface{}) []byte {
	jsonString, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	return jsonString
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
