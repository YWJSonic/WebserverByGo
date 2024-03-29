package main

import (
	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/foundation/fileload"
	"gitlab.com/ServerUtility/gamelimit"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/myhttp"
	_ "gitlab.com/ServerUtility/mysql"
	"gitlab.com/ServerUtility/thirdparty/ulginfo"
	"gitlab.com/WebserverByGoBase/foundation/myrestful"
	"gitlab.com/WebserverByGoBase/gamerule"
	crontab "gitlab.com/WebserverByGoBase/handlecrontab"
	db "gitlab.com/WebserverByGoBase/handledb"
	event "gitlab.com/WebserverByGoBase/handleevent"
	"gitlab.com/WebserverByGoBase/serversetting"
	login "gitlab.com/WebserverByGoBase/serviceaccount"
	game "gitlab.com/WebserverByGoBase/servicegame"
	lobby "gitlab.com/WebserverByGoBase/servicelobby"
	"gitlab.com/WebserverByGoBase/servicethirdparty/api"
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

	gamelimit.ServerDayPayLimit = foundation.InterfaceToInt64(config["ServerDayPayLimit"])
	gamelimit.ServerDayPayDefault = foundation.InterfaceToInt64(config["ServerDayPayDefault"])

	setting := struct{ DBUser, DBPassword, DBIP, DBPORT string }{
		serversetting.DBUser,
		serversetting.DBPassword,
		serversetting.DBIP,
		serversetting.DBPORT,
	}
	db.SetDBConn(&setting)
	go event.Update()

	result, err := db.GetSetting()
	if err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("Main", err)
		panic("DB GetSetting Error")
	}
	serversetting.InsertDBSetting(result, gamerule.GameIndex)
	serversetting.RefreshDBSetting(gamerule.GameIndex, gamelimit.ServerDayPayDefault)

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
