package main

import (
	"./foundation"
	"./frame/transmission"
	"./game/slotgame"
	"./lobby"
)

func main() {
	// var err error
	// foundation.CacheData, err = cache.New(cache.Config{
	// 	RedisURL:  "127.0.0.1:6379",
	// 	MustRedis: true,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	var initArray [][]transmission.RESTfulURL
	initArray = append(initArray, slotgame.ServiceStart())
	initArray = append(initArray, lobby.ServiceStart())

	foundation.HTTPLisentRun(foundation.ServerSetting.URL(), initArray...)
}
