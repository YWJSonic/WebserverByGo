package main

import (
	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/foundation/fileload"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/myhttp"
	_ "gitlab.com/ServerUtility/mysql"
	"gitlab.com/ServerUtility/thirdparty/ulginfo"
	"gitlab.com/WeberverByGoGame9/foundation/myrestful"
	"gitlab.com/WeberverByGoGame9/gamerule"
	crontab "gitlab.com/WeberverByGoGame9/handlecrontab"
	db "gitlab.com/WeberverByGoGame9/handledb"
	event "gitlab.com/WeberverByGoGame9/handleevent"
	"gitlab.com/WeberverByGoGame9/serversetting"
	login "gitlab.com/WeberverByGoGame9/serviceaccount"
	game "gitlab.com/WeberverByGoGame9/servicegame"
	lobby "gitlab.com/WeberverByGoGame9/servicelobby"
	"gitlab.com/WeberverByGoGame9/servicethirdparty/api"
)

func main() {

	jsStr := fileload.Load("./file/config.json")
	config := foundation.StringToJSON(jsStr)

	serversetting.GameTypeID = foundation.InterfaceToString(config["GameTypeID"])
	serversetting.IP = foundation.InterfaceToString(config["IP"])
	serversetting.PORT = foundation.InterfaceToString(config["PORT"])
	serversetting.DBIP = foundation.InterfaceToString(config["DBIP"])
	serversetting.DBPORT = foundation.InterfaceToString(config["DBPORT"])
	serversetting.DBUser = foundation.InterfaceToString(config["DBUser"])
	serversetting.DBPassword = foundation.InterfaceToString(config["DBPassword"])
	serversetting.AccountEncodeStr = foundation.InterfaceToString(config["AccountEncodeStr"])
	serversetting.RedisURL = foundation.InterfaceToString(config["RedisURL"])
	serversetting.MaintainStartTime = foundation.InterfaceToString(config["MaintainStartTime"])
	serversetting.MaintainFinishTime = foundation.InterfaceToString(config["MaintainFinishTime"])
	messagehandle.IsPrintLog = foundation.InterfaceToBool(config["DebugLog"])
	serversetting.EnableMaintain(foundation.InterfaceToBool(config["Maintain"]))
	gamerule.SetInfo(gamerule.GameIndex, config)

	ulginfo.LoginURL = foundation.InterfaceToString(config["ULGLoginURL"])
	ulginfo.GetuserURL = foundation.InterfaceToString(config["ULGGetuserURL"])
	ulginfo.AuthorizedURL = foundation.InterfaceToString(config["ULGAuthorizedURL"])
	ulginfo.ExchangeURL = foundation.InterfaceToString(config["ULGExchangeURL"])
	ulginfo.CheckoutURL = foundation.InterfaceToString(config["ULGCheckoutURL"])
	ulginfo.ULGMaintainCheckoutTime = foundation.InterfaceToString(config["ULGMaintainCheckoutTime"])

	db.SetDBConn()
	go event.Update()

	result, err := db.GetSetting()
	if err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("Main", err)
		panic("DB GetSetting Error")
	}
	serversetting.InsertDBSetting(result, gamerule.GameIndex)

	crontab.NewCron(serversetting.MaintainStartTime, func() {
		serversetting.EnableMaintain(true)
	})
	crontab.NewCron(serversetting.MaintainFinishTime, func() {
		serversetting.EnableMaintain(false)
	})
	crontab.NewCron(ulginfo.ULGMaintainCheckoutTime, api.MaintainCheckout)

	var initArray [][]myhttp.RESTfulURL
	initArray = append(initArray, login.ServiceStart())
	initArray = append(initArray, lobby.ServiceStart())
	initArray = append(initArray, game.ServiceStart())
	initArray = append(initArray, api.ServiceStart())
	myrestful.HTTPLisentRun(serversetting.ServerURL(), initArray...)
}
