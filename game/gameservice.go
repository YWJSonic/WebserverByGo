package game

import (
	"net/http"
	"sync"

	"gitlab.com/WeberverByGo/foundation"
	"gitlab.com/WeberverByGo/messagehandle/errorlog"

	"github.com/julienschmidt/httprouter"
)

var mu *sync.RWMutex

// RoomCount Current room count
var RoomCount = 0

// HandleURL ...
var HandleURL []foundation.RESTfulURL

func init() {
	mu = new(sync.RWMutex)
	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "game/gameresult", Fun: gameresult, ConnType: foundation.Client})
}

// ServiceStart ...
func ServiceStart() []foundation.RESTfulURL {
	return HandleURL
}

func gameresult(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()

	// var result = make(map[string]interface{})
	postData := foundation.PostData(r)
	// gametoken := foundation.InterfaceToString(postData["token"])
	BetIndex := foundation.InterfaceToInt64(postData["bet"])

	// gametype check
	err := errorlog.New()
	// gametypeid := foundation.InterfaceToString(postData["gametypeid"])
	// if gametypeid != data.GameTypeID {
	// 	err.ErrorCode = code.GameTypeError
	// 	err.Msg = "GameTypeError"
	// 	foundation.HTTPResponse(w, "", err)
	// 	return
	// }

	// get player
	// playerid := foundation.InterfaceToInt64(postData["playerid"])
	// playerinfo, err := player.GetPlayerInfoByPlayerID(playerid)
	// if err.ErrorCode != code.OK {
	// 	foundation.HTTPResponse(w, "", err)
	// 	return
	// }

	// get thirdparty info data
	// var ulginfo *ulg.ULGInfo
	// ulginfo, err = ulg.GetULGInfo(playerinfo.GameToken)
	// if err.ErrorCode != code.OK {
	// 	foundation.HTTPResponse(w, "", err)
	// 	fmt.Println(ulginfo)
	// 	return
	// }

	result, totalwinscore := gameRequest(1, BetIndex)
	result["playermoney"] = totalwinscore
	// loginfo := log.New(log.GameResult)
	// loginfo.PlayerID = playerid
	// loginfo.IValue1 = int64(WinBet * BetMoney)
	// loginfo.SValue1 = fmt.Sprint(newplate)
	// loginfo.SValue2 = fmt.Sprint(ScatterIndex)
	// log.SaveLog(loginfo)

	// result["ScatterGame"] = ScatterIndex
	// result["NormalGameIndex"] = newplateIndex
	// result["NormalGame"] = newplate
	// result["WinMoney"] = WinMoney
	foundation.HTTPResponse(w, result, err)

}
