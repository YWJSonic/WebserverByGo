package main

import (
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/foundation/fileload"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/myhttp"
	_ "gitlab.com/ServerUtility/mysql"
	"gitlab.com/ServerUtility/thirdparty/ulginfo"
	"gitlab.com/WeberverByGoGame8/foundation/myrestful"
	"gitlab.com/WeberverByGoGame8/gamerule"
	crontab "gitlab.com/WeberverByGoGame8/handlecrontab"
	db "gitlab.com/WeberverByGoGame8/handledb"
	event "gitlab.com/WeberverByGoGame8/handleevent"
	"gitlab.com/WeberverByGoGame8/serversetting"
	login "gitlab.com/WeberverByGoGame8/serviceaccount"
	game "gitlab.com/WeberverByGoGame8/servicegame"
	lobby "gitlab.com/WeberverByGoGame8/servicelobby"
	"gitlab.com/WeberverByGoGame8/servicethirdparty/api"
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

	var initArray [][]myhttp.RESTfulURL
	initArray = append(initArray, login.ServiceStart())
	initArray = append(initArray, lobby.ServiceStart())
	initArray = append(initArray, game.ServiceStart())
	initArray = append(initArray, api.ServiceStart())
	db.SetDBConn()

	crontab.NewCron(serversetting.MaintainStartTime, func() {
		serversetting.EnableMaintain(true)
	})

	crontab.NewCron(serversetting.MaintainFinishTime, func() {
		serversetting.EnableMaintain(false)
	})
	crontab.NewCron(ulginfo.ULGMaintainCheckoutTime, api.MaintainCheckout)

	go event.Update()
	myrestful.HTTPLisentRun(serversetting.ServerURL(), initArray...)
}
