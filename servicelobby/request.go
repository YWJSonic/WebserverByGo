package lobby

import (
	"net/http"
	"sync"

	"gitlab.com/WeberverByGo/apithirdparty"
	attach "gitlab.com/WeberverByGo/handleattach"

	"gitlab.com/WeberverByGo/serversetting"

	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/loginfo"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/myhttp"
	mycache "gitlab.com/WeberverByGo/handlecache"
	log "gitlab.com/WeberverByGo/handlelog"
	"gitlab.com/WeberverByGo/player"

	"github.com/julienschmidt/httprouter"
	gameRule "gitlab.com/gamerule"
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
		myhttp.HTTPResponse(w, "", err)
		return
	}

	GameAccount := foundation.InterfaceToString(postData["gameaccount"])
	if GameAccount == "" {
		err.ErrorCode = code.Unauthenticated
		err.Msg = "GameAccountError"
		myhttp.HTTPResponse(w, "", err)
		return
	}

	// check token
	if err = foundation.CheckToken(mycache.GetToken(GameAccount), token); err.ErrorCode != code.OK {
		myhttp.HTTPResponse(w, "", err)
		return
	}

	playerInfo, err := player.GetPlayerInfoByGameAccount(GameAccount)
	if err.ErrorCode != code.OK {
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
	result["attach"] = gameRule.ConvertToGameAttach(playerInfo.ID, attach.GetAttach(playerInfo.ID, gameRule.GameIndex, gameRule.IsAttachSaveToDB))

	myhttp.HTTPResponse(w, result, err)
}

func refresh(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()
	err := messagehandle.New()
	postData := myhttp.PostData(r)

	gametypeid := foundation.InterfaceToString(postData["gametypeid"])
	if err = foundation.CheckGameType(serversetting.GameTypeID, gametypeid); err.ErrorCode != code.OK {
		myhttp.HTTPResponse(w, "", err)
		return
	}

	accounttoken := foundation.InterfaceToString(postData["accounttoken"])
	userCoinQuota, err := apithirdparty.Refresh(accounttoken, gametypeid)
	if err.ErrorCode != code.OK {
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
		myhttp.HTTPResponse(w, "", err)
		return
	}

	// get player
	playerInfo, err := player.GetPlayerInfoByPlayerID(playerID)
	if err.ErrorCode != code.OK {
		myhttp.HTTPResponse(w, "", err)
		return
	}

	// check token
	if err = foundation.CheckToken(mycache.GetToken(playerInfo.GameAccount), token); err.ErrorCode != code.OK {
		myhttp.HTTPResponse(w, "", err)
		return
	}

	ulgResult, err := apithirdparty.Excahnge(playerInfo, accountToken, gametypeid, cointype, coinamount)
	if err.ErrorCode != code.OK {
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
		myhttp.HTTPResponse(w, "", err)
		return
	}

	playerInfo, err := player.GetPlayerInfoByPlayerID(playerID)
	if err.ErrorCode != code.OK {
		myhttp.HTTPResponse(w, "", err)
		return
	}

	// check token
	if err = foundation.CheckToken(mycache.GetToken(playerInfo.GameAccount), token); err.ErrorCode != code.OK {
		myhttp.HTTPResponse(w, "", err)
		return
	}

	if playerInfo.GameToken == "" {
		err.ErrorCode = code.NoExchange
		err.Msg = "NoExchange"
		myhttp.HTTPResponse(w, "", err)
		return
	}

	ulgCheckOutResult, err := apithirdparty.CheckOut(playerInfo, serversetting.GameTypeID)
	if err.ErrorCode != code.OK {
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
