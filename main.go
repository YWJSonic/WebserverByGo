package main

import (
	_ "github.com/go-sql-driver/mysql"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/foundation/fileload"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/myhttp"
	"gitlab.com/ServerUtility/thirdparty/ulginfo"
	"gitlab.com/WeberverByGo/crontab"
	"gitlab.com/WeberverByGo/data"
	"gitlab.com/WeberverByGo/db"
	"gitlab.com/WeberverByGo/event"
	"gitlab.com/WeberverByGo/foundation/myrestful"
	"gitlab.com/WeberverByGo/game"
	"gitlab.com/WeberverByGo/lobby"
	"gitlab.com/WeberverByGo/login"
	"gitlab.com/WeberverByGo/service/api"
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
	messagehandle.IsPrintLog = foundation.InterfaceToBool(config["DebugLog"])
	data.EnableMaintain(foundation.InterfaceToBool(config["Maintain"]))

	ulginfo.LoginURL = foundation.InterfaceToString(config["ULGLoginURL"])
	ulginfo.GetuserURL = foundation.InterfaceToString(config["ULGGetuserURL"])
	ulginfo.AuthorizedURL = foundation.InterfaceToString(config["ULGAuthorizedURL"])
	ulginfo.ExchangeURL = foundation.InterfaceToString(config["ULGExchangeURL"])
	ulginfo.CheckoutURL = foundation.InterfaceToString(config["ULGCheckoutURL"])
	ulginfo.ULGMaintainCheckoutTime = foundation.InterfaceToString(config["ULGMaintainCheckoutTime"])

	var initArray [][]myhttp.RESTfulURL
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
	crontab.NewCron(ulginfo.ULGMaintainCheckoutTime, api.MaintainCheckout)

	go event.Update()
	myrestful.HTTPLisentRun(data.ServerURL(), initArray...)
}
