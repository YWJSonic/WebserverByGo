package lobby

import (
	"fmt"
	"net/http"
	"sync"

	"../code"
	"../db"
	"../foundation"
	"../log"
	"../messagehandle/errorlog"
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

	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "lobby/getplayer", Fun: getplayer})
	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "lobby/exchange", Fun: exchange})
	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "lobby/checkout", Fun: checkout})
	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "lobby/cachetest", Fun: cachetest})
	return HandleURL
}

func getplayer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	result := make(map[string]interface{})
	loginfo := log.New(log.GetPlayer)
	err := errorlog.New()
	postData := foundation.PostData(r)
	GameAccount := foundation.InterfaceToString(postData["gameaccount"])
	// token := foundation.InterfaceToString(postData["token"])
	// ServerToken := data.GetToken(GameAccount)

	// if ServerToken != token {
	// 	err.ErrorCode = code.Unauthenticated
	// 	foundation.HTTPResponse(w, playerInfo, err)
	// 	return
	// }

	if GameAccount == "" {
		err.ErrorCode = code.Unauthenticated
		err.Msg = "GameAccountError"
		foundation.HTTPResponse(w, result, err)
		return
	}

	dbresult, err := db.GetPlayerInfoByGameAccount(GameAccount)
	if err.ErrorCode != code.OK {
		fmt.Print(err.Msg)
	}

	var playerInfo *player.PlayerInfo
	if len(dbresult) <= 0 {
		playerInfo = player.New(GameAccount)
	} else {
		playerInfo = player.MakePlayer(dbresult[0])
	}

	loginfo.PlayerID = playerInfo.ID
	player.SavePlayerInfo(playerInfo)
	log.SaveLog(loginfo)
	foundation.HTTPResponse(w, playerInfo, err)
}

func cachetest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	postData := foundation.PostData(r)
	playerid := foundation.InterfaceToInt64(postData["playerid"])

	// fmt.Println("playerInfo", playerid)
	playerInfo, err := player.GetPlayerInfoByPlayerID(playerid)
	foundation.HTTPResponse(w, playerInfo, err)

}

func exchange(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	loginfo := log.New(log.Exchange)
	var result = make(map[string]interface{})
	postData := foundation.PostData(r)
	// token := foundation.InterfaceToString(postData["token"])
	playerID := foundation.InterfaceToInt64(postData["playerid"])
	accounttoken := postData["accounttoken"].(string)
	gametoken := postData["gametoken"].(string)
	gametypeid := postData["gametypeid"].(string)
	cointype := foundation.InterfaceToInt(postData["cointype"])
	coinamount := foundation.InterfaceToInt(postData["coinamount"])

	ulgResult, err := ulg.Exchange(gametoken, gametypeid, accounttoken, cointype, coinamount)

	// if !exchangeInfo["result"].(bool) {
	// 	err.ErrorCode = code.FailedPrecondition
	// 	err.Msg = exchangeInfo["errorMsg"].(string)
	// 	foundation.HTTPResponse(w, exchangeInfo, err)
	// }

	var playerInfo *player.PlayerInfo
	playerInfo, err = player.GetPlayerInfoByPlayerID(playerID)
	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, result, err)
		return
	}

	playerInfo.Money += ulgResult.GameCoin
	player.SavePlayerInfo(playerInfo)
	loginfo.PlayerID = playerInfo.ID
	loginfo.IValue1 = int64(cointype)
	loginfo.IValue2 = int64(coinamount)
	loginfo.IValue3 = playerInfo.Money
	log.SaveLog(loginfo)
	foundation.HTTPResponse(w, ulgResult, err)
}
func checkout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}
