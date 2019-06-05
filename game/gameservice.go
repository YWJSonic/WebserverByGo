package game

import (
	"fmt"
	"net/http"
	"sync"

	"gitlab.com/WeberverByGo/code"
	"gitlab.com/WeberverByGo/data"
	"gitlab.com/WeberverByGo/foundation"
	"gitlab.com/WeberverByGo/messagehandle/errorlog"
	"gitlab.com/WeberverByGo/player"
	"gitlab.com/WeberverByGo/thirdparty/ulg"

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
	betIndex := foundation.InterfaceToInt64(postData["bet"])
	betMoney := GetBetMoney(betIndex)

	// gametype check
	err := errorlog.New()
	gametypeid := foundation.InterfaceToString(postData["gametypeid"])
	if gametypeid != data.GameTypeID {
		err.ErrorCode = code.GameTypeError
		err.Msg = "GameTypeError"
		foundation.HTTPResponse(w, "", err)
		return
	}

	// get player
	playerid := foundation.InterfaceToInt64(postData["playerid"])
	playerinfo, err := player.GetPlayerInfoByPlayerID(playerid)
	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return
	}

	if playerinfo.Money < betMoney {
		err.ErrorCode = code.NoMoneyToBet
		err.Msg = "NoMoneyToBet"
		foundation.HTTPResponse(w, "", err)
		return
	}

	// get thirdparty info data
	var ulginfo *ulg.ULGInfo
	ulginfo, err = ulg.GetULGInfo(playerinfo.ID, playerinfo.GameToken)
	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		fmt.Println(ulginfo)
		return
	}

	result, totalwinscore := gameRequest(playerinfo.ID, betMoney)
	playerinfo.Money = playerinfo.Money + totalwinscore - betMoney
	result["playermoney"] = playerinfo.Money

	ulginfo.TotalBet += betMoney
	ulginfo.TotalWin += totalwinscore
	player.SavePlayerInfo(playerinfo)
	ulg.SaveULGInfo(ulginfo)

	foundation.HTTPResponse(w, result, err)

}
