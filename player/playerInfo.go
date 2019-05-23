package player

import (
	"encoding/json"
	"strings"
	"time"

	"../data"
)

// CachePlayer memory cache player
var CachePlayer map[int64]PlayerInfo

// PlayerInfo Player information
type PlayerInfo struct {
	ID          int64  `json:"ID"`
	Money       int64  `json:"Money"`
	GameAccount string `json:"GameAccount"`

	///////// for Server value
	GameToken     string `json:"GameToken,omitempty"`
	InRoom        int    `json:"InRoom,omitempty"`        // room index
	LastCheckTime int64  `json:"LastCheckTime,omitempty"` // connect check time
	InGame        string `json:"InGame,omitempty"`        // gametype
}

// ToJSONClient ...
func (p PlayerInfo) ToJSONClient() map[string]interface{} {
	clientdata := make(map[string]interface{})
	clientdata["ID"] = p.ID
	clientdata["Money"] = p.Money
	clientdata["GameAccount"] = p.GameAccount
	return clientdata
}

// ToJSONStr ...
func (p PlayerInfo) ToJSONStr() string {
	data, _ := json.MarshalIndent(p, "", " ")
	STR := string(data)
	STR = strings.ReplaceAll(STR, string(10), ``)
	return STR
}

// ToJson conver to json string
func (p PlayerInfo) ToJson() []byte {
	data, _ := json.MarshalIndent(p, "", " ")
	return data
}

// ResultMap player conver to map, client data
func (p PlayerInfo) ResultMap() map[string]interface{} {
	return map[string]interface{}{
		"ID":    p.ID,
		"Money": p.Money,
		// "Token": p.Token,
	}
}

// IsPlayerConnect check player connect time
func (p PlayerInfo) IsPlayerConnect() bool {
	return (time.Now().Unix() - p.LastCheckTime) <= data.ConnectTimeOut
}

// IsInGameRoom is player in game room
func (p PlayerInfo) IsInGameRoom() bool {
	return p.InRoom != 0
}

// IsCheckOut is ulg checkout
func (p PlayerInfo) IsCheckOut() bool {
	return p.GameToken == ""
}
