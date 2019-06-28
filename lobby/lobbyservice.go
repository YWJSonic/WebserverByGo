package lobby

import (
	"net/http"
	"sync"

	"gitlab.com/ServerUtility/myhttp"
	"gitlab.com/WeberverByGo/code"
	"gitlab.com/WeberverByGo/db"
	"gitlab.com/WeberverByGo/foundation"
	"gitlab.com/WeberverByGo/foundation/myrestful"
	"gitlab.com/WeberverByGo/messagehandle/errorlog"
	"gitlab.com/WeberverByGo/messagehandle/log"
	"gitlab.com/WeberverByGo/player"
	"gitlab.com/WeberverByGo/thirdparty/ulg"

	"github.com/julienschmidt/httprouter"
	gameRule "gitlab.com/game7"
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
	err := errorlog.New()
	postData := myhttp.PostData(r)
	token := foundation.InterfaceToString(postData["token"])
	gametypeid := foundation.InterfaceToString(postData["gametypeid"])

	if err = foundation.CheckGameType(gametypeid); err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "", err)
		return
	}

	GameAccount := foundation.InterfaceToString(postData["gameaccount"])
	if GameAccount == "" {
		err.ErrorCode = code.Unauthenticated
		err.Msg = "GameAccountError"
		myrestful.HTTPResponse(w, "", err)
		return
	}

	// check token
	if err = foundation.CheckToken(GameAccount, token); err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "", err)
		return
	}

	var dbresult []map[string]interface{}
	dbresult, err = db.GetPlayerInfoByGameAccount(GameAccount)
	if err.ErrorCode != code.OK {
		errorlog.ErrorLogPrintln("Lobby", err.Msg)
		myrestful.HTTPResponse(w, "", err)
		return
	}

	var playerInfo *player.PlayerInfo
	if len(dbresult) <= 0 {
		playerInfo, err = player.New(GameAccount)
	} else {
		playerInfo = player.MakePlayer(dbresult[0])
	}

	if playerInfo == nil {
		errorlog.ErrorLogPrintln("Lobby", err.Msg)
		myrestful.HTTPResponse(w, "", err)
		return
	}

	// ulg result info
	thirdpartyresult := make(map[string]interface{})
	thirdpartyresult["isexchange"] = 0
	if playerInfo.GameToken != "" {
		thirdpartyresult["isexchange"] = 1
	}
	result["thirdparty"] = thirdpartyresult

	gameRule.InitAttach(playerInfo.ID)
	player.SavePlayerInfo(playerInfo)
	result["player"] = playerInfo.ToJSONClient()
	result["reel"] = gameRule.GetInitScroll()
	result["betrate"] = gameRule.GetInitBetRate()
	// result["gameattach"] = game.GetAttach(playerInfo.ID)

	myrestful.HTTPResponse(w, result, err)
}

func refresh(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()
	err := errorlog.New()
	postData := myhttp.PostData(r)
	//
	// logintype := foundation.InterfaceToInt(postData["logintype"])
	// if logintype != account.Ulg {
	// 	err.ErrorCode = code.AccountTypeError
	// 	err.Msg = "AccountTypeError"
	// 	foundation.HTTPResponse(w, "", err)
	// 	return
	// }

	gametypeid := foundation.InterfaceToString(postData["gametypeid"])
	if err = foundation.CheckGameType(gametypeid); err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "", err)
		return
	}

	accounttoken := foundation.InterfaceToString(postData["accounttoken"])
	UserInfo, err := ulg.GetUser(accounttoken, gametypeid)
	if err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "", err)
		return
	}

	result := make(map[string]interface{})
	result["userCoinQuota"] = UserInfo.UserCoinQuota

	myrestful.HTTPResponse(w, result, err)
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

	err := errorlog.New()
	if err = foundation.CheckGameType(gametypeid); err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "", err)
		return
	}

	// get player
	var playerInfo *player.PlayerInfo
	playerInfo, err = player.GetPlayerInfoByPlayerID(playerID)
	if err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "", err)
		return
	}

	// check token
	if err = foundation.CheckToken(playerInfo.GameAccount, token); err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "", err)
		return
	}

	// if playerInfo.GameToken != "" {
	// 	err.ErrorCode = code.NoCheckoutError
	// 	err.Msg = "NoCheckoutError"
	// 	foundation.HTTPResponse(w, "", err)
	// 	return
	// }

	// new thirdparty token
	ulguser, err := ulg.Authorized(accountToken, gametypeid)
	if err.ErrorCode != code.OK {
		err.ErrorCode = code.FailedPrecondition
		myrestful.HTTPResponse(w, "", err)
		return
	}

	// exchange
	errorlog.LogPrintln("AccountUserInfo.GameToken", ulguser.GameToken)
	ulgResult, err := ulg.Exchange(ulguser.GameToken, gametypeid, accountToken, cointype, coinamount)
	if err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "", err)
		return
	}

	_, err = ulg.NewULGInfo(playerInfo.ID, cointype, coinamount, ulguser.GameToken, accountToken)
	if err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "", err)
		return
	}

	OldMoney := playerInfo.Money
	playerInfo.Money = ulgResult.GameCoin
	playerInfo.GameToken = ulguser.GameToken
	player.SavePlayerInfo(playerInfo)
	db.SetExchange(playerInfo.GameAccount, playerInfo.GameToken, cointype, coinamount, playerInfo.Money, OldMoney, foundation.ServerNowTime())

	loginfo := log.New(log.Exchange)
	loginfo.PlayerID = playerInfo.ID
	loginfo.IValue1 = int64(cointype)
	loginfo.IValue2 = int64(coinamount)
	loginfo.IValue3 = ulgResult.GameCoin
	log.SaveLog(loginfo)

	myrestful.HTTPResponse(w, ulgResult, err)
}

func checkout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()
	postData := myhttp.PostData(r)
	token := foundation.InterfaceToString(postData["token"])
	playerID := foundation.InterfaceToInt64(postData["playerid"])
	gametypeid := foundation.InterfaceToString(postData["gametypeid"])
	// accountToken := foundation.InterfaceToString(postData["accounttoken"])

	err := errorlog.New()
	if err = foundation.CheckGameType(gametypeid); err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "", err)
		return
	}

	playerInfo, err := player.GetPlayerInfoByPlayerID(playerID)
	if err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "", err)
		return
	}

	// check token
	if err = foundation.CheckToken(playerInfo.GameAccount, token); err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "", err)
		return
	}

	if playerInfo.GameToken == "" {
		err.ErrorCode = code.NoExchange
		err.Msg = "NoExchange"
		myrestful.HTTPResponse(w, "", err)
		return
	}

	var ulginfo *ulg.ULGInfo
	ulginfo, err = ulg.GetULGInfo(playerInfo.ID, playerInfo.GameToken)
	if err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "", err)
		return
	}
	if ulginfo.IsCheckOut {
		playerInfo.Money = 0
		playerInfo.GameToken = ""
		player.SavePlayerInfo(playerInfo)
		myrestful.HTTPResponse(w, "", err)
		return
	}

	var ulgCheckOutResult ulg.UlgCheckOutResult
	ulgCheckOutResult, err = ulg.Checkout(ulginfo, gametypeid)
	if err.ErrorCode != code.OK && err.ErrorCode != code.ExchangeError {
		myrestful.HTTPResponse(w, "", err)
		return
	}

	loginfo := log.New(log.CheckOut)
	loginfo.PlayerID = playerInfo.ID
	loginfo.IValue1 = playerInfo.Money
	loginfo.SValue1 = playerInfo.GameToken
	log.SaveLog(loginfo)

	ulg.UpdateUlgInfoCheckOut(playerInfo.GameToken)

	playerInfo.Money = 0
	playerInfo.GameToken = ""
	player.SavePlayerInfo(playerInfo)

	result := make(map[string]interface{})
	result["userCoinQuota"] = ulgCheckOutResult.UserCoinQuota

	myrestful.HTTPResponse(w, result, err)
}
