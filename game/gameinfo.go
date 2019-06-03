package game

import (
	"encoding/json"
	"strings"
)

type gameInfo struct {
	FreeCount int64 `json:"freecount"`
}

// ToJSONStr ...
func (p gameInfo) ToJSONStr() string {
	data, _ := json.MarshalIndent(p, "", " ")
	STR := string(data)
	STR = strings.ReplaceAll(STR, string(10), ``)
	return STR
}
