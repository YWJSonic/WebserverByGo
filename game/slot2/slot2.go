package slot2

import (
	"fmt"
	"net/http"
	"sync"

	"../../code"
	"../../foundation"
	"../../log"
	"../../player"
	"../../thirdparty/ulg"
	"../gamelogic"

	"github.com/julienschmidt/httprouter"
)

var mu *sync.RWMutex

// RoomCount Current room count
var RoomCount = 0
var TotalServerWin int64 = 0
var TotalServerLost int64 = 0

// HandleURL ...
var HandleURL []foundation.RESTfulURL

func init() {
	mu = new(sync.RWMutex)
	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "slot2/gameresult", Fun: gameresult, ConnType: foundation.Client})
}

// ServiceStart ...
func ServiceStart() []foundation.RESTfulURL {
	return HandleURL
}

func gameresult(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()

	// POSTvalues := map[string][]string{"POST": []string{fmt.Sprintf(`{"bet": %d}`, 100)}}
	// gameresult := foundation.HTTPPostRequest("http://192.168.1.15:8100/slot2/gameresult", POSTvalues)

	// foundation.HTTPResponse(w, gameresult, errorlog.New())
	// fmt.Println(gameresult)
	// return

	var result = make(map[string]interface{})

	postData := foundation.PostData(r)
	// gametoken := foundation.InterfaceToString(postData["token"])
	// gametypeid := foundation.InterfaceToString(postData["gametypeid"])
	BetMoney := foundation.InterfaceToInt64(postData["bet"])

	playerid := foundation.InterfaceToInt64(postData["playerid"])
	playerinfo, err := player.GetPlayerInfoByPlayerID(playerid)
	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return
	}

	var ulginfo *ulg.ULGInfo
	ulginfo, err = ulg.GetULGInfo(playerinfo.GameToken)
	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return
	}

	var gameresult map[string]interface{}
	gameresult, err = gamelogic.GetGameResult("ulg", playerid, BetMoney)
	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return
	}

	var LostMoney int64
	var WinBet int64
	var ScatterIndex int
	var newplate []int
	var newplateIndex []int
	WinBet = foundation.InterfaceToInt64(gameresult["WinBet"])
	WinMoney := int64(WinBet) * BetMoney
	if BetMoney > WinMoney {
		LostMoney := foundation.Abs(WinMoney - BetMoney)
		TotalServerWin += LostMoney
	}

	TotalServerLost += WinMoney

	// ulginfo, err := ulg.GetULGInfo(playerinfo.GameToken)
	// ulginfo.TotalBet += BetMoney
	// ulginfo.TotalWin += WinMoney
	// ulginfo.TotalLost += LostMoney
	// ulg.SaveULGInfo(ulginfo)

	ulg.UpdateULGInfo(ulginfo, BetMoney, WinBet)

	playerinfo.Money = playerinfo.Money + WinMoney - LostMoney
	player.SavePlayerInfo(playerinfo)

	loginfo := log.New(log.GameResult)
	loginfo.PlayerID = playerid
	loginfo.IValue1 = int64(WinBet * BetMoney)
	loginfo.SValue1 = fmt.Sprint(newplate)
	loginfo.SValue2 = fmt.Sprint(ScatterIndex)
	log.SaveLog(loginfo)

	result["ScatterGame"] = ScatterIndex
	result["NormalGameIndex"] = newplateIndex
	result["NormalGame"] = newplate
	result["WinMoney"] = WinMoney
	foundation.HTTPResponse(w, result, err)

}
