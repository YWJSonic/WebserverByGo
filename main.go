package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "gitlab.com/WeberverByGo/crontab"
	"gitlab.com/WeberverByGo/data"
	"gitlab.com/WeberverByGo/db"
	"gitlab.com/WeberverByGo/event"
	"gitlab.com/WeberverByGo/foundation"
	"gitlab.com/WeberverByGo/game"
	"gitlab.com/WeberverByGo/lobby"
	"gitlab.com/WeberverByGo/login"
	"gitlab.com/WeberverByGo/service/api"
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
