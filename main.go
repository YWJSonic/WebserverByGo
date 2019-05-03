package main

import (
	"fmt"
	"math/rand"
	"time"

	"./account"
	"./code"
	"./data"
	"./db"
	"./event"
	"./foundation"
	"./game/slotgame"
	"./lobby"
	"./player"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// var i int64
	// for i = 11; i < 13; i++ {
	// 	go getplayer(i)
	// }

	var initArray [][]foundation.RESTfulURL
	initArray = append(initArray, account.ServiceStart())
	initArray = append(initArray, slotgame.ServiceStart())
	initArray = append(initArray, lobby.ServiceStart())
	// foundation.NewToken()
	go db.SetDBCOnn()
	go event.Update()

	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		db.NewLogTable("20060102")
	}
	foundation.HTTPLisentRun(data.ServerURL(), initArray...)
}

func getplayer(i int64) {
	rmint := rand.Intn(5) + 1
	for {
		time.Sleep(time.Millisecond * time.Duration(rmint))
		playerInfo, err := player.GetPlayerInfoByPlayerID(i)
		fmt.Println("playerid", i, playerInfo, err)
		if err.ErrorCode != code.OK {
			return
		}
	}
}
