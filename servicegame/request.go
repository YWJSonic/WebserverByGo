package game

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"gitlab.com/WeberverByGo/apithirdparty/ulg"
	"gitlab.com/WeberverByGo/player"
	"gitlab.com/WeberverByGo/serversetting"

	gameRule "gitlab.com/WeberverByGo/gamerule"
	attach "gitlab.com/WeberverByGo/handleattach"
	mycache "gitlab.com/WeberverByGo/handlecache"
	log "gitlab.com/WeberverByGo/handlelog"

	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/gamelimit"
	"gitlab.com/ServerUtility/loginfo"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/myhttp"

	"gitlab.com/ServerUtility/httprouter"
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
		myhttp.HTTPResponse(w, "", err)
		return
	}

	// get player
	playerID := foundation.InterfaceToInt64(postData["playerid"])
	playerInfo, err := player.GetPlayerInfoByPlayerID(playerID)
	if err.ErrorCode != code.OK {
		myhttp.HTTPResponse(w, "", err)
		return
	}

	// check token
	if err = foundation.CheckToken(mycache.GetToken(playerInfo.GameAccount), token); err.ErrorCode != code.OK {
		myhttp.HTTPResponse(w, "", err)
		return
	}

	// money check
	if playerInfo.Money < betMoney {
		err.ErrorCode = code.NoMoneyToBet
		err.Msg = "NoMoneyToBet"
		myhttp.HTTPResponse(w, "", err)
		return
	}

	// get thirdparty info data
	ulginfo, err := ulg.GetULGInfo(playerInfo.ID, playerInfo.GameToken)
	if err.ErrorCode != code.OK {
		myhttp.HTTPResponse(w, "", err)
		fmt.Println(ulginfo)
		return
	}

	var att []map[string]interface{}
	var totalwinscore int64
	var result map[string]interface{}
	var newatt []map[string]interface{}
	var otherdata map[string]int64
	att = attach.GetAttach(playerInfo.ID, gameRule.GameIndex, gameRule.IsAttachSaveToDB)

	for index, max := 0, 2; index < max; index++ {
		result, newatt, otherdata = gameRule.GameRequest(playerID, betIndex, att)
		totalwinscore = otherdata["totalwinscore"]
		JackPartBonusx2 := otherdata["JackPartBonusx2"]
		JackPartBonusx3 := otherdata["JackPartBonusx3"]
		JackPartBonusx5 := otherdata["JackPartBonusx5"]
		noJpWin := totalwinscore - JackPartBonusx2 - JackPartBonusx3 - JackPartBonusx5

		if gamelimit.IsInTotalMoneyWinLimit(gameRule.WinScoreLimit, betMoney, noJpWin) && gamelimit.IsInTotalBetRateWinLimit(gameRule.WinBetRateLimit, betMoney, noJpWin) {
			break
		}

	}

	playerInfo.Money = playerInfo.Money + totalwinscore - betMoney
	attach.SaveAttach(playerInfo.ID, gameRule.GameIndex, newatt, gameRule.IsAttachSaveToDB)
	result["playermoney"] = playerInfo.Money
	result["attach"] = gameRule.ConvertToGameAttach(playerInfo.ID, newatt)

	ulginfo.TotalBet += betMoney
	ulginfo.TotalWin += totalwinscore
	player.SavePlayerInfo(playerInfo)
	ulg.SaveULGInfo(ulginfo)

	msg := foundation.JSONToString(result)
	msg = strings.ReplaceAll(msg, "\"", "\\\"")
	loginfo := loginfo.New(loginfo.GameResult)
	loginfo.PlayerID = playerID
	loginfo.IValue1 = foundation.InterfaceToInt64(result["totalwinscore"])
	loginfo.IValue2 = otherdata["betmoney"]
	loginfo.Msg = msg
	log.SaveLog(loginfo)

	myhttp.HTTPResponse(w, result, err)

}
