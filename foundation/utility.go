package foundation

import (
	"encoding/json"

	"../frame/code"
	"../frame/player"
)

// DeleteArrayElement ...
func DeleteArrayElement(ElementIndex interface{}, array []interface{}) []interface{} {
	count := len(array)
	for index := 0; index < count; index++ {
		return append(array[:index], array[index+1:]...)
	}
	return array
}

// NewPlayerInfo Create a new PlayerInfo
func NewPlayerInfo(id code.PlayerID) player.PlayerInfo {
	return player.PlayerInfo{
		ID:    id,
		Money: 10000,
		Token: "1234567890",
	}
}

// GetPlayerInfo getplayerinfo
func GetPlayerInfo(id code.PlayerID) player.PlayerInfo {
	if player, ok := CachePlayer[id]; ok {
		return *player
	}
	player := NewPlayerInfo(id)
	CachePlayer[id] = &player
	return player
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
func InterfaceTofloat64(v interface{}) int {
	return int(v.(float64))
}

// InterfaceToInt ...
func InterfaceToInt(v interface{}) int {
	return int(InterfaceTofloat64(v))
}

// InterfaceToString ...
func InterfaceToString(v interface{}) string {
	return v.(string)
}
