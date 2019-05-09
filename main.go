package main

import (
	"./account"
	"./data"
	"./db"
	"./event"
	"./foundation"
	"./game/slotgame"
	"./lobby"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var initArray [][]foundation.RESTfulURL
	initArray = append(initArray, account.ServiceStart())
	initArray = append(initArray, slotgame.ServiceStart())
	initArray = append(initArray, lobby.ServiceStart())
	// foundation.NewToken()
	db.SetDBCOnn()

	// mycache.ClearAllCache()

	go event.Update()
	foundation.HTTPLisentRun(data.ServerURL(), initArray...)
}
