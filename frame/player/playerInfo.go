package player

// CachePlayer memory cache player
var CachePlayer map[int64]PlayerInfo

// ThirdpartyInfo ...
type ThirdpartyInfo interface {
	GetInfo() interface{}
}

// PlayerInfo Player information
type PlayerInfo struct {
	ID         int64
	Money      int64
	Token      string
	GameToken  string
	InGame     string         // gametype
	InRoom     int            // room index
	Thirdparty ThirdpartyInfo // plant data
}

// IsInGameRoom is player in game room
func (p PlayerInfo) IsInGameRoom() bool {
	return p.InRoom != -1
}
