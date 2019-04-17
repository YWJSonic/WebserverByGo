package foundation

import (
	"../cache"
	"../frame/code"
	"../frame/player"
)

var CacheData cache.Cache
var CachePlayer map[code.PlayerID]*player.PlayerInfo

func init() {
	CachePlayer = make(map[code.PlayerID]*player.PlayerInfo)
}
