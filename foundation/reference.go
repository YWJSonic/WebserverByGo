package foundation

import (
	"../cache"
	"../frame/player"
	ServiceSetting "../frame/setting"
	"../frame/transmission"
)

var CacheData cache.Cache
var CachePlayer map[int64]*player.PlayerInfo
var ServerSetting ServiceSetting.ServerOption
var ProxyData map[string]transmission.RESTfulURL

func init() {
	ProxyData = make(map[string]transmission.RESTfulURL)
	CachePlayer = make(map[int64]*player.PlayerInfo)
	ServerSetting = ServiceSetting.ServerOption{
		IP:       "192.168.1.10",
		PORT:     "8000",
		Maintain: false,
	}
}
