package game

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"gitlab.com/WeberverByGo/apithirdparty/ulg"
	"gitlab.com/WeberverByGo/foundation/myrestful"
	attach "gitlab.com/WeberverByGo/handleattach"
	db "gitlab.com/WeberverByGo/handledb"
	log "gitlab.com/WeberverByGo/handlelog"
	"gitlab.com/WeberverByGo/player"
	"gitlab.com/WeberverByGo/serversetting"

	mycache "gitlab.com/WeberverByGo/handlecache"
	gameRule "gitlab.com/game7"

	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/loginfo"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/myhttp"
	"gitlab.com/ServerUtility/thirdparty/ulginfo"

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
	postData := myhttp.PostData(r)
	token := foundation.InterfaceToString(postData["token"])
	betIndex := foundation.InterfaceToInt64(postData["bet"])
	betMoney := gameRule.GetBetMoney(betIndex)

	// gametype check
	err := messagehandle.New()
	gametypeid := foundation.InterfaceToString(postData["gametypeid"])
	if gametypeid != serversetting.GameTypeID {
		err.ErrorCode = code.GameTypeError
		err.Msg = "GameTypeError"
		myrestful.HTTPResponse(w, "", err)
		return
	}

	// get player
	playerID := foundation.InterfaceToInt64(postData["playerid"])
	playerInfo, err := player.GetPlayerInfoByPlayerID(playerID)
	if err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "", err)
		return
	}

	// check token
	if err = foundation.CheckToken(mycache.GetToken(playerInfo.GameAccount), token); err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "", err)
		return
	}

	// money check
	if playerInfo.Money < betMoney {
		err.ErrorCode = code.NoMoneyToBet
		err.Msg = "NoMoneyToBet"
		myrestful.HTTPResponse(w, "", err)
		return
	}

	// get thirdparty info data
	var ulginfo *ulginfo.Info
	ulginfo, err = ulg.GetULGInfo(playerInfo.ID, playerInfo.GameToken)
	if err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "", err)
		fmt.Println(ulginfo)
		return
	}

	st := time.Now()

	var att []map[string]interface{}
	att, err = db.GetAttachKind(playerInfo.ID, gameRule.GameIndex)
	result, newatt, otherdata := gameRule.GameRequest(playerID, betIndex, att)
	totalwinscore := otherdata["totalwinscore"]
	playerInfo.Money = playerInfo.Money + totalwinscore - betMoney
	result["playermoney"] = playerInfo.Money
	result["attach"] = gameRule.GetAttach(newatt)

	fmt.Println(time.Since(st))
	st = time.Now()
	attach.SaveAttach(playerInfo.ID, gameRule.GameIndex, newatt, false)

	fmt.Println(time.Since(st))
	st = time.Now()

	ulginfo.TotalBet += betMoney
	ulginfo.TotalWin += totalwinscore
	player.SavePlayerInfo(playerInfo)
	ulg.SaveULGInfo(ulginfo)

	fmt.Println(time.Since(st))
	st = time.Now()

	msg := foundation.JSONToString(result)
	msg = strings.ReplaceAll(msg, "\"", "\\\"")
	loginfo := loginfo.New(loginfo.GameResult)
	loginfo.PlayerID = playerID
	loginfo.IValue1 = foundation.InterfaceToInt64(result["totalwinscore"])
	loginfo.IValue2 = otherdata["betmoney"]
	loginfo.Msg = msg
	log.SaveLog(loginfo)

	fmt.Println(time.Since(st))
	st = time.Now()
	myrestful.HTTPResponse(w, result, err)

	fmt.Println(time.Since(st))
	st = time.Now()

}
