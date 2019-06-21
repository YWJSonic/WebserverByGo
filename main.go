package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "gitlab.com/WeberverByGo/crontab"
	"gitlab.com/WeberverByGo/data"
	"gitlab.com/WeberverByGo/db"
	"gitlab.com/WeberverByGo/event"
	"gitlab.com/WeberverByGo/foundation"
	"gitlab.com/WeberverByGo/foundation/fileload"
	"gitlab.com/WeberverByGo/game"
	"gitlab.com/WeberverByGo/lobby"
	"gitlab.com/WeberverByGo/login"
	"gitlab.com/WeberverByGo/service/api"
	"gitlab.com/WeberverByGo/thirdparty/ulg"
)

func main() {

	jsStr := fileload.Load("./file/config.json")
	config := foundation.StringToJSON(jsStr)

	data.GameTypeID = foundation.InterfaceToString(config["GameTypeID"])
	data.IP = foundation.InterfaceToString(config["IP"])
	data.PORT = foundation.InterfaceToString(config["PORT"])
	data.DBIP = foundation.InterfaceToString(config["DBIP"])
	data.DBPORT = foundation.InterfaceToString(config["DBPORT"])
	data.DBUser = foundation.InterfaceToString(config["DBUser"])
	data.DBPassword = foundation.InterfaceToString(config["DBPassword"])
	data.AccountEncodeStr = foundation.InterfaceToString(config["AccountEncodeStr"])
	data.RedisURL = foundation.InterfaceToString(config["RedisURL"])
	data.MaintainStartTime = foundation.InterfaceToString(config["MaintainStartTime"])
	data.MaintainFinishTime = foundation.InterfaceToString(config["MaintainFinishTime"])
	data.Maintain = foundation.InterfaceToBool(config["Maintain"])

	ulg.LoginURL = foundation.InterfaceToString(config["ULGLoginURL"])
	ulg.GetuserURL = foundation.InterfaceToString(config["ULGGetuserURL"])
	ulg.AuthorizedURL = foundation.InterfaceToString(config["ULGAuthorizedURL"])
	ulg.ExchangeURL = foundation.InterfaceToString(config["ULGExchangeURL"])
	ulg.CheckoutURL = foundation.InterfaceToString(config["ULGCheckoutURL"])

	var initArray [][]foundation.RESTfulURL
	initArray = append(initArray, login.ServiceStart())
	initArray = append(initArray, lobby.ServiceStart())
	initArray = append(initArray, game.ServiceStart())
	initArray = append(initArray, api.ServiceStart())
	db.SetDBConn()

	go event.Update()
	foundation.HTTPLisentRun(data.ServerURL(), initArray...)
}
