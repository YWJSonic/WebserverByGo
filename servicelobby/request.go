package lobby

import (
	"net/http"
	"sync"

	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/httprouter"
	"gitlab.com/ServerUtility/loginfo"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/myhttp"
	"gitlab.com/WebserverByGoGame6/apithirdparty"
	gameRule "gitlab.com/WebserverByGoGame6/gamerule"
	attach "gitlab.com/WebserverByGoGame6/handleattach"
	mycache "gitlab.com/WebserverByGoGame6/handlecache"
	log "gitlab.com/WebserverByGoGame6/handlelog"
	"gitlab.com/WebserverByGoGame6/player"
	"gitlab.com/WebserverByGoGame6/serversetting"
	"gitlab.com/WebserverByGoGame6/servicethirdparty/api"
)

var isInit = false
var mu *sync.RWMutex

// ServiceStart ...
func ServiceStart() []myhttp.RESTfulURL {
	var HandleURL []myhttp.RESTfulURL

	if isInit {
		return HandleURL
	}
	mu = new(sync.RWMutex)
	isInit = true

	HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "lobby/init", Fun: gameinit, ConnType: myhttp.Client})
	HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "lobby/refresh", Fun: refresh, ConnType: myhttp.Client})
	HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "lobby/exchange", Fun: exchange, ConnType: myhttp.Client})
	HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "lobby/checkout", Fun: checkout, ConnType: myhttp.Client})
	return HandleURL
}

func gameinit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()
	result := make(map[string]interface{})
	// loginfo := log.New(log.GetPlayer)
	err := messagehandle.New()
	postData := myhttp.PostData(r)
	token := foundation.InterfaceToString(postData["token"])
	gametypeid := foundation.InterfaceToString(postData["gametypeid"])

	if err = foundation.CheckGameType(serversetting.GameTypeID, gametypeid); err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("gameinit-1", err, token, gametypeid)
		myhttp.HTTPResponse(w, "", err)
		return
	}

	GameAccount := foundation.InterfaceToString(postData["gameaccount"])
	if GameAccount == "" {
		err.ErrorCode = code.Unauthenticated
		err.Msg = "GameAccountError"
		messagehandle.ErrorLogPrintln("gameinit-2", err, token, gametypeid)
		myhttp.HTTPResponse(w, "", err)
		return
	}

	// check token
	if err = foundation.CheckToken(mycache.GetToken(GameAccount), token); err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("gameinit-3", err, token, gametypeid)
		myhttp.HTTPResponse(w, "", err)
		return
	}

	playerInfo, err := player.GetPlayerInfoByGameAccount(GameAccount)
	if err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("gameinit-4", err, token, gametypeid)
		myhttp.HTTPResponse(w, "Lobby", err)
	}

	// ulg result info
	thirdpartyresult := make(map[string]interface{})
	thirdpartyresult["isexchange"] = 0
	if playerInfo.GameToken != "" {
		thirdpartyresult["isexchange"] = 1
	}

	result["thirdparty"] = thirdpartyresult
	result["player"] = playerInfo.ToJSONClient()
	result["reel"] = gameRule.GetInitScroll()
	result["betrate"] = gameRule.GetInitBetRate()

	if gameRule.IsAttachInit {
		result["attach"] = gameRule.ConvertToGameAttach(playerInfo.ID, attach.GetAttach(playerInfo.ID, gameRule.GameIndex, gameRule.IsAttachSaveToDB))
	}
	myhttp.HTTPResponse(w, result, err)
}

func refresh(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()
	err := messagehandle.New()
	postData := myhttp.PostData(r)

	gametypeid := foundation.InterfaceToString(postData["gametypeid"])
	if err = foundation.CheckGameType(serversetting.GameTypeID, gametypeid); err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("refresh-1", err, gametypeid)
		myhttp.HTTPResponse(w, "", err)
		return
	}

	accounttoken := foundation.InterfaceToString(postData["accounttoken"])
	userCoinQuota, err := apithirdparty.Refresh(accounttoken, gametypeid)
	if err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("refresh-2", err, gametypeid, accounttoken)
		myhttp.HTTPResponse(w, "", err)
		return
	}

	result := make(map[string]interface{})
	result["userCoinQuota"] = userCoinQuota
	myhttp.HTTPResponse(w, result, err)
}

func exchange(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()
	postData := myhttp.PostData(r)
	token := foundation.InterfaceToString(postData["token"])
	playerID := foundation.InterfaceToInt64(postData["playerid"])
	accountToken := foundation.InterfaceToString(postData["accounttoken"])
	gametypeid := postData["gametypeid"].(string)
	cointype := foundation.InterfaceToInt64(postData["cointype"])
	coinamount := foundation.InterfaceToInt64(postData["coinamount"])

	err := messagehandle.New()
	if err = foundation.CheckGameType(serversetting.GameTypeID, gametypeid); err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("exchange-1", err, playerID, token, accountToken, gametypeid, cointype, coinamount)
		myhttp.HTTPResponse(w, "", err)
		return
	}

	// get player
	playerInfo, err := player.GetPlayerInfoByPlayerID(playerID)
	if err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("exchange-2", err, playerID, token, accountToken, gametypeid, cointype, coinamount)
		myhttp.HTTPResponse(w, "", err)
		return
	}

	// check token
	if err = foundation.CheckToken(mycache.GetToken(playerInfo.GameAccount), token); err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("exchange-3", err, playerID, token, accountToken, gametypeid, cointype, coinamount, playerInfo)
		myhttp.HTTPResponse(w, "", err)
		return
	}

	ulgResult, err := apithirdparty.Excahnge(playerInfo, accountToken, gametypeid, cointype, coinamount)
	if err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("exchange-4", err, playerID, token, accountToken, gametypeid, cointype, coinamount, playerInfo)
		myhttp.HTTPResponse(w, "", err)
		return

	}

	loginfo := loginfo.New(loginfo.Exchange)
	loginfo.PlayerID = playerInfo.ID
	loginfo.IValue1 = int64(cointype)
	loginfo.IValue2 = int64(coinamount)
	loginfo.IValue3 = ulgResult.GameCoin
	log.SaveLog(loginfo)

	myhttp.HTTPResponse(w, ulgResult, err)
}

func checkout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()
	postData := myhttp.PostData(r)
	token := foundation.InterfaceToString(postData["token"])
	playerID := foundation.InterfaceToInt64(postData["playerid"])
	gametypeid := foundation.InterfaceToString(postData["gametypeid"])
	// accountToken := foundation.InterfaceToString(postData["accounttoken"])

	err := messagehandle.New()
	if err = foundation.CheckGameType(serversetting.GameTypeID, gametypeid); err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("checkout-1", err, playerID, token, gametypeid)
		myhttp.HTTPResponse(w, "", err)
		return
	}

	playerInfo, err := player.GetPlayerInfoByPlayerID(playerID)
	if err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("checkout-2", err, playerID, token, gametypeid)
		myhttp.HTTPResponse(w, "", err)
		return
	}

	// check token
	if err = foundation.CheckToken(mycache.GetToken(playerInfo.GameAccount), token); err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("checkout-2", err, playerID, token, gametypeid, playerInfo)
		myhttp.HTTPResponse(w, "", err)
		return
	}

	if playerInfo.GameToken == "" {
		err.ErrorCode = code.NoExchange
		err.Msg = "NoExchange"
		messagehandle.ErrorLogPrintln("checkout-3", err, playerID, token, gametypeid, playerInfo)
		myhttp.HTTPResponse(w, "", err)
		return
	}

	api.RunNotFinishSoctter(playerInfo.ID)

	ulgCheckOutResult, err := apithirdparty.CheckOut(playerInfo, serversetting.GameTypeID)
	if err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("checkout-4", err, playerID, token, gametypeid, playerInfo)
		myhttp.HTTPResponse(w, "", err)
		return
	}

	loginfo := loginfo.New(loginfo.CheckOut)
	loginfo.PlayerID = playerInfo.ID
	loginfo.IValue1 = playerInfo.Money
	loginfo.SValue1 = playerInfo.GameToken
	log.SaveLog(loginfo)

	playerInfo.Money = 0
	playerInfo.GameToken = ""
	player.SavePlayerInfo(playerInfo)

	result := make(map[string]interface{})
	result["userCoinQuota"] = ulgCheckOutResult

	myhttp.HTTPResponse(w, result, err)
}
