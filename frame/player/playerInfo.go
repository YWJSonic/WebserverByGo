package player

import (
	"../code"
)

// CachePlayer memory cache player
var CachePlayer map[code.PlayerID]PlayerInfo

// ThirdpartyInfo ...
type ThirdpartyInfo interface {
	GetInfo() interface{}
}

// PlayerInfo Player information
type PlayerInfo struct {
	ID         code.PlayerID
	Money      int64
	Token      string
	GameToken  string
	InGame     string         // gametype
	InRoom     int            // room index
	Thirdparty ThirdpartyInfo // plant data
}
