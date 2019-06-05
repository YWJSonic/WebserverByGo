package game

import (
	"encoding/json"
	"strings"
)

// Attach game att data
type Attach struct {
	FreeCount int `json:"freecount"`
}

// ToJSONStr ...
func (gi Attach) ToJSONStr() string {
	data, _ := json.MarshalIndent(gi, "", " ")
	STR := string(data)
	STR = strings.ReplaceAll(STR, string(10), ``)
	return STR
}
