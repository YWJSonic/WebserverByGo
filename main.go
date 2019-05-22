package main

import (
	_ "./crontab"
	"./data"
	"./db"
	"./event"
	"./foundation"
	"./game"
	"./lobby"
	"./login"
	"./service/api"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var initArray [][]foundation.RESTfulURL
	initArray = append(initArray, login.ServiceStart())
	initArray = append(initArray, lobby.ServiceStart())
	initArray = append(initArray, game.ServiceStart())
	initArray = append(initArray, api.ServiceStart())
	// foundation.NewToken()
	db.SetDBConn()

	// mycache.ClearAllCache()

	go event.Update()
	foundation.HTTPLisentRun(data.ServerURL(), initArray...)
}
