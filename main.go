package main

import (
	"./account"
	_ "./crontab"
	"./data"
	"./db"
	"./event"
	"./foundation"
	"./game"
	"./lobby"
	"./service/api"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var initArray [][]foundation.RESTfulURL
	initArray = append(initArray, account.ServiceStart())
	initArray = append(initArray, game.ServiceStart())
	initArray = append(initArray, lobby.ServiceStart())
	initArray = append(initArray, api.ServiceStart())
	// foundation.NewToken()
	db.SetDBConn()

	// mycache.ClearAllCache()

	go event.Update()
	foundation.HTTPLisentRun(data.ServerURL(), initArray...)
}
