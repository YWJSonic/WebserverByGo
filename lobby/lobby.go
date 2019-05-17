package lobby

import (
	"fmt"
	"net/http"
	"sync"

	"../code"
	"../data"
	"../db"
	"../foundation"
	"../log"
	"../messagehandle/errorlog"
	"../mycache"
	"../player"
	"../thirdparty/ulg"

	"github.com/julienschmidt/httprouter"
)

var isInit = false
var mu *sync.RWMutex

// ServiceStart ...
func ServiceStart() []foundation.RESTfulURL {
	var HandleURL []foundation.RESTfulURL

	if isInit {
		return HandleURL
	}
	mu = new(sync.RWMutex)
	isInit = true

	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "lobby/getplayer", Fun: getplayer, ConnType: foundation.Client})
	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "lobby/exchange", Fun: exchange, ConnType: foundation.Client})
	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "lobby/checkout", Fun: checkout, ConnType: foundation.Client})
	return HandleURL
}

func getplayer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()
	loginfo := log.New(log.GetPlayer)
	err := errorlog.New()
	postData := foundation.PostData(r)
	token := foundation.InterfaceToString(postData["token"])

	GameAccount := foundation.InterfaceToString(postData["gameaccount"])
	if GameAccount == "" {
		err.ErrorCode = code.Unauthenticated
		err.Msg = "GameAccountError"
		foundation.HTTPResponse(w, "", err)
		return
	}

	ServerToken := mycache.GetToken(GameAccount)
	if ServerToken != token {
		err.ErrorCode = code.Unauthenticated
		err.Msg = "TokenError"
		foundation.HTTPResponse(w, "", err)
		return
	}

	var dbresult []map[string]interface{}
	dbresult, err = db.GetPlayerInfoByGameAccount(GameAccount)
	if err.ErrorCode != code.OK {
		errorlog.ErrorLogPrintln("Lobby", err.Msg)
		foundation.HTTPResponse(w, "", err)
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
		foundation.HTTPResponse(w, "", err)
		return
	}

	loginfo.PlayerID = playerInfo.ID
	player.SavePlayerInfo(playerInfo)
	log.SaveLog(loginfo)
	foundation.HTTPResponse(w, playerInfo, err)
}

func exchange(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()
	postData := foundation.PostData(r)
	token := foundation.InterfaceToString(postData["token"])
	playerID := foundation.InterfaceToInt64(postData["playerid"])
	accounttoken := postData["accounttoken"].(string)
	gametypeid := postData["gametypeid"].(string)
	cointype := foundation.InterfaceToInt(postData["cointype"])
	coinamount := foundation.InterfaceToInt(postData["coinamount"])

	// get player
	var err errorlog.ErrorMsg
	if gametypeid != data.GameTypeID {
		err.ErrorCode = code.GameTypeError
		err.Msg = "GameTypeError"
		foundation.HTTPResponse(w, "", err)
		return
	}

	var playerInfo *player.PlayerInfo
	playerInfo, err = player.GetPlayerInfoByPlayerID(playerID)
	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return
	}

	// check token
	ServerToken := mycache.GetToken(playerInfo.GameAccount)
	if ServerToken != token {
		err.ErrorCode = code.Unauthenticated
		err.Msg = "TokenError"
		foundation.HTTPResponse(w, "", err)
		return
	}

	// get account
	var AccountInfo *player.AccountInfo
	AccountInfo, err = player.GetAccountInfoByGameAccount(playerInfo.GameAccount)
	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return
	}
	if playerInfo.GameToken != "" {
		err.ErrorCode = code.NoCheckoutError
		err.Msg = "NoCheckoutError"
		foundation.HTTPResponse(w, "", err)
		return
	}

	// new thirdparty token
	ulguser, err := ulg.Authorized(accounttoken, gametypeid)
	if err.ErrorCode != code.OK {
		err.ErrorCode = code.FailedPrecondition
		foundation.HTTPResponse(w, "", err)
		return
	}

	// exchange
	errorlog.LogPrintln("AccountUserInfo.GameToken", ulguser.GameToken)
	ulgResult, err := ulg.Exchange(ulguser.GameToken, gametypeid, accounttoken, cointype, coinamount)
	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return
	}

	_, err = ulg.NewULGInfo(playerInfo.ID, ulguser.GameToken, accounttoken)
	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return
	}

	OldMoney := playerInfo.Money
	playerInfo.Money = playerInfo.Money + ulgResult.GameCoin
	playerInfo.GameToken = ulguser.GameToken
	player.SavePlayerInfo(playerInfo)
	db.SetExchange(playerInfo.GameAccount, playerInfo.GameToken, cointype, coinamount, playerInfo.Money, OldMoney, foundation.ServerNowTime())
	mycache.SetAccountInfo(AccountInfo.GameAccount, AccountInfo.ToJSONStr())

	loginfo := log.New(log.Exchange)
	loginfo.PlayerID = playerInfo.ID
	loginfo.IValue1 = int64(cointype)
	loginfo.IValue2 = int64(coinamount)
	loginfo.IValue3 = ulgResult.GameCoin
	log.SaveLog(loginfo)

	foundation.HTTPResponse(w, ulgResult, err)
}

func checkout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()
	postData := foundation.PostData(r)
	token := foundation.InterfaceToString(postData["token"])
	playerID := foundation.InterfaceToInt64(postData["playerid"])
	gametypeid := foundation.InterfaceToString(postData["gametypeid"])
	accountToken := foundation.InterfaceToString(postData["accounttoken"])

	err := errorlog.New()
	if gametypeid != data.GameTypeID {
		err.ErrorCode = code.GameTypeError
		err.Msg = "GameTypeError"
		foundation.HTTPResponse(w, "", err)
		return
	}

	playerInfo, err := player.GetPlayerInfoByPlayerID(playerID)
	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return
	}

	var AccountInfo *player.AccountInfo
	AccountInfo, err = player.GetAccountInfoByGameAccount(playerInfo.GameAccount)
	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return
	}

	ServerToken := AccountInfo.Token
	if ServerToken != token {
		err.ErrorCode = code.Unauthenticated
		err.Msg = "TokenError"
		foundation.HTTPResponse(w, "", err)
		return
	}

	if playerInfo.GameToken == "" {
		err.ErrorCode = code.NoExchange
		err.Msg = "NoExchange"
		foundation.HTTPResponse(w, "", err)
		return
	}

	var ulginfo *ulg.ULGInfo
	ulginfo, err = ulg.GetULGInfo(playerInfo.GameToken)
	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return

	}

	var ulgResult ulg.UlgResult
	ulgResult, err = ulg.Checkout(accountToken, playerInfo.GameToken, gametypeid, fmt.Sprint(ulginfo.TotalBet), ulginfo.TotalWin, ulginfo.TotalLost)
	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return
	}

	// err =
	// if err.ErrorCode != code.OK {
	// 	foundation.HTTPResponse(w, "", err)
	// 	return
	// }

	loginfo := log.New(log.CheckOut)
	loginfo.PlayerID = playerInfo.ID
	loginfo.IValue1 = playerInfo.Money
	loginfo.SValue1 = playerInfo.GameToken
	log.SaveLog(loginfo)

	playerInfo.GameToken = ""
	player.SavePlayerInfo(playerInfo)

	foundation.HTTPResponse(w, ulgResult, err)
}
