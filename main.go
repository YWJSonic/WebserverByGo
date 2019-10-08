package main

import (
	_ "github.com/go-sql-driver/mysql"
	"gitlab.com/WeberverByGoGame5/crontab"
	"gitlab.com/WeberverByGoGame5/data"
	"gitlab.com/WeberverByGoGame5/db"
	"gitlab.com/WeberverByGoGame5/event"
	"gitlab.com/WeberverByGoGame5/foundation"
	"gitlab.com/WeberverByGoGame5/foundation/fileload"
	"gitlab.com/WeberverByGoGame5/game"
	gameRules "gitlab.com/WeberverByGoGame5/game/game5"
	"gitlab.com/WeberverByGoGame5/lobby"
	"gitlab.com/WeberverByGoGame5/login"
	"gitlab.com/WeberverByGoGame5/messagehandle/errorlog"
	"gitlab.com/WeberverByGoGame5/service/api"
	"gitlab.com/WeberverByGoGame5/thirdparty/ulg"
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
	errorlog.IsPrintLog = foundation.InterfaceToBool(config["DebugLog"])

	gameRules.SetInfo(gameRules.GameIndex, config)
	data.EnableMaintain(foundation.InterfaceToBool(config["Maintain"]))

	ulg.LoginURL = foundation.InterfaceToString(config["ULGLoginURL"])
	ulg.GetuserURL = foundation.InterfaceToString(config["ULGGetuserURL"])
	ulg.AuthorizedURL = foundation.InterfaceToString(config["ULGAuthorizedURL"])
	ulg.ExchangeURL = foundation.InterfaceToString(config["ULGExchangeURL"])
	ulg.CheckoutURL = foundation.InterfaceToString(config["ULGCheckoutURL"])
	ulg.ULGMaintainCheckoutTime = foundation.InterfaceToString(config["ULGMaintainCheckoutTime"])

	var initArray [][]foundation.RESTfulURL
	initArray = append(initArray, login.ServiceStart())
	initArray = append(initArray, lobby.ServiceStart())
	initArray = append(initArray, game.ServiceStart())
	initArray = append(initArray, api.ServiceStart())
	db.SetDBConn()

	crontab.NewCron(data.MaintainStartTime, func() {
		data.EnableMaintain(true)
	})

	crontab.NewCron(data.MaintainFinishTime, func() {
		data.EnableMaintain(false)
	})
	crontab.NewCron(ulg.ULGMaintainCheckoutTime, api.MaintainCheckout)

	go event.Update()
	foundation.HTTPLisentRun(data.ServerURL(), initArray...)
}
