package lobby

import (
	"net/http"
	"sync"

	"../frame/transmission"

	"../foundation"
	"../messagehandle/errorlog"
	"../thirdparty/ulg"

	"github.com/julienschmidt/httprouter"
)

var isInit = false
var mu *sync.RWMutex

// ServiceStart ...
func ServiceStart() []transmission.RESTfulURL {
	var HandleURL []transmission.RESTfulURL

	if isInit {
		return HandleURL
	}
	mu = new(sync.RWMutex)
	isInit = true

	HandleURL = append(HandleURL, transmission.RESTfulURL{RequestType: "POST", URL: "lobby/getplayer", Fun: getplayer})
	HandleURL = append(HandleURL, transmission.RESTfulURL{RequestType: "POST", URL: "lobby/exchange", Fun: exchange})
	return HandleURL
}

func getplayer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := errorlog.ErrorMsg{}
	postData := foundation.PostData(r)
	token := (postData["token"]).(string)
	gameid := (postData["gameid"]).(string)

	AccountUserInfo := ulg.Authorized(token, gameid)

	foundation.HTTPResponse(w, AccountUserInfo, err)
}

func exchange(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := errorlog.New()
	postData := foundation.PostData(r)
	playerID := foundation.InterfaceToInt(postData["playerid"])
	token := postData["token"].(string)
	gametoken := postData["gametoken"].(string)
	gameid := postData["gameid"].(string)
	cointype := foundation.InterfaceToInt(postData["cointype"])
	coinamount := foundation.InterfaceToInt(postData["coinamount"])

	ulgResult := ulg.Exchange(gametoken, gameid, token, cointype, coinamount)

	// if !exchangeInfo["result"].(bool) {
	// 	err.ErrorCode = code.FailedPrecondition
	// 	err.Msg = exchangeInfo["errorMsg"].(string)
	// 	foundation.HTTPResponse(w, exchangeInfo, err)
	// }

	player := foundation.GetPlayerInfo(playerID)
	player.Money += ulgResult.GameCoin
	foundation.SavePlayerInfo(player)

	foundation.HTTPResponse(w, ulgResult, err)
}
