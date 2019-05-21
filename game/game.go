package game

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"../data"
	"../foundation"

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

func apitest() {
	POSTvalues := map[string]interface{}{
		"d1": 0,
		"d2": 1,
		"d3": 2,
		"d4": 3,
		"d5": "4",
		"d6": "5",
		"d7": "6",
	}
	// gameresult := foundation.HTTPPostRequest(GameResultURL, POSTvalues)

	value, _ := json.MarshalIndent(POSTvalues, "", " ")
	// STR := string(data)
	// STR = strings.ReplaceAll(STR, string(10), ``)

	// fmt.Println(STR)
	// str := []byte(
	// 	`{"d1":0,"d2":1,"d3":2,"d4":3,"d5":"4","d6":"5","d7":"6"}`)
	// 	`{"d1":0,"d2":1,"d3":2,"d4":3,"d5":"4","d6":"5"}`
	gameresult := foundation.HTTPPostRawRequest(data.GameResultURL, value)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	fmt.Println("----", string(gameresult))
}

func gameresult(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()

	apitest()

	return

	// var result = make(map[string]interface{})
	// postData := foundation.PostData(r)
	// gametoken := foundation.InterfaceToString(postData["token"])
	// BetMoney := foundation.InterfaceToInt64(postData["bet"])

	// gametype check
	// err := errorlog.New()
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
	// 	return
	// }

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
	// foundation.HTTPResponse(w, result, err)

}
