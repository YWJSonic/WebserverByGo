package game

import (
	"net/http"
	"sync"

	"gitlab.com/WeberverByGo/db"
	"gitlab.com/WeberverByGo/foundation/myrestful"

	gameRule "gitlab.com/game7"

	"gitlab.com/ServerUtility/myhttp"
	"gitlab.com/WeberverByGo/foundation"
	"gitlab.com/WeberverByGo/messagehandle/errorlog"

	"github.com/julienschmidt/httprouter"
)

var mu *sync.RWMutex

// RoomCount Current room count
var RoomCount = 0

// HandleURL ...
var HandleURL []myhttp.RESTfulURL

func init() {
	mu = new(sync.RWMutex)
	HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "game/gameresult", Fun: gameresult, ConnType: myhttp.Client})
}

// ServiceStart ...
func ServiceStart() []myhttp.RESTfulURL {
	return HandleURL
}

func gameresult(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()

	postData := myhttp.PostData(r)
	// token := foundation.InterfaceToString(postData["token"])
	betIndex := foundation.InterfaceToInt64(postData["bet"])
	// betMoney := gameRule.GetBetMoney(betIndex)

	// // gametype check
	err := errorlog.New()
	// gametypeid := foundation.InterfaceToString(postData["gametypeid"])
	// if gametypeid != data.GameTypeID {
	// 	err.ErrorCode = code.GameTypeError
	// 	err.Msg = "GameTypeError"
	// 	foundation.HTTPResponse(w, "", err)
	// 	return
	// }

	// // get player
	playerID := foundation.InterfaceToInt64(postData["playerid"])
	// playerinfo, err := player.GetPlayerInfoByPlayerID(playerID)
	// if err.ErrorCode != code.OK {
	// 	foundation.HTTPResponse(w, "", err)
	// 	return
	// }

	// // check token
	// if err = foundation.CheckToken(playerinfo.GameAccount, token); err.ErrorCode != code.OK {
	// 	foundation.HTTPResponse(w, "", err)
	// 	return
	// }

	// // money check
	// if playerinfo.Money < betMoney {
	// 	err.ErrorCode = code.NoMoneyToBet
	// 	err.Msg = "NoMoneyToBet"
	// 	foundation.HTTPResponse(w, "", err)
	// 	return
	// }

	// // get thirdparty info data
	// var ulginfo *ulg.ULGInfo
	// ulginfo, err = ulg.GetULGInfo(playerinfo.ID, playerinfo.GameToken)
	// if err.ErrorCode != code.OK {
	// 	foundation.HTTPResponse(w, "", err)
	// 	fmt.Println(ulginfo)
	// 	return
	// }

	att, _ := db.GetAttachKind(4, 7)
	// att := make([]map[string]interface{}, 0)
	result, att, _ := gameRule.GameRequest(playerID, betIndex, att)
	// playerinfo.Money = playerinfo.Money + totalwinscore - betMoney
	// result["playermoney"] = playerinfo.Money

	// ulginfo.TotalBet += betMoney
	// ulginfo.TotalWin += totalwinscore
	// player.SavePlayerInfo(playerinfo)
	// ulg.SaveULGInfo(ulginfo)

	myrestful.HTTPResponse(w, result, err)

}
