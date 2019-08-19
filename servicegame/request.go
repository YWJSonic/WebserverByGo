package game

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"gitlab.com/WeberverByGoGame6/apithirdparty/ulg"
	db "gitlab.com/WeberverByGoGame6/handledb"
	"gitlab.com/WeberverByGoGame6/player"
	"gitlab.com/WeberverByGoGame6/serversetting"

	gameRule "gitlab.com/WeberverByGoGame6/gamerule"
	attach "gitlab.com/WeberverByGoGame6/handleattach"
	mycache "gitlab.com/WeberverByGoGame6/handlecache"
	log "gitlab.com/WeberverByGoGame6/handlelog"

	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/gamelimit"
	"gitlab.com/ServerUtility/httprouter"
	"gitlab.com/ServerUtility/loginfo"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/myhttp"
	"gitlab.com/ServerUtility/playerinfo"
	"gitlab.com/ServerUtility/thirdparty/ulginfo"
)

var mu *sync.RWMutex

// RoomCount Current room count
var RoomCount = 0

// HandleURL ...
var HandleURL []myhttp.RESTfulURL

func init() {
	mu = new(sync.RWMutex)
	HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "game/gameresult", Fun: gameresult, ConnType: myhttp.Client})
	HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "game/scottergameresult", Fun: scottergameresult, ConnType: myhttp.Client})
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
	serverTotalPayScore := serversetting.GetServerTotalPayScore()

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
	var otherdata map[string]interface{}
	att = attach.GetAttachByType(playerInfo.ID, gameRule.GameIndex, gameRule.DayScotterGameCountKey, gameRule.IsAttachSaveToDB)

	for index, max := 0, 2; index < max; index++ {
		result, newatt, otherdata = gameRule.GameRequest(playerID, betIndex, att)
		totalwinscore = otherdata["totalwinscore"].(int64)

		if gamelimit.IsInTotalMoneyWinLimit(gameRule.WinScoreLimit, betMoney, totalwinscore) && gamelimit.IsInTotalBetRateWinLimit(gameRule.WinBetRateLimit, betMoney, totalwinscore) {
			break
		}

	}

	playerInfo.Money = playerInfo.Money + totalwinscore - betMoney
	result["playermoney"] = playerInfo.Money

	attach.SaveAttach(playerInfo.ID, gameRule.GameIndex, []map[string]interface{}{newatt[0]}, false)
	attach.SaveAttachToDB(playerInfo.ID, gameRule.GameIndex, newatt)

	ulginfo.TotalBet += betMoney
	ulginfo.TotalWin += totalwinscore
	player.SavePlayerInfo(playerInfo)
	ulg.SaveULGInfo(ulginfo)

	msg := foundation.JSONToString(result)
	msg = strings.ReplaceAll(msg, "\"", "\\\"")
	loginfo := loginfo.New(loginfo.GameResult)
	loginfo.PlayerID = playerID
	loginfo.IValue1 = foundation.InterfaceToInt64(result["totalwinscore"])
	loginfo.IValue2 = otherdata["betmoney"].(int64)
	// loginfo.IValue3 = scotterid
	loginfo.Msg = msg
	log.SaveLog(loginfo)

	if totalwinscore > 0 {
		serverTotalPayScore += totalwinscore
		serversetting.SetServerTotalPayScore(serverTotalPayScore)
		db.UpdateSetting(foundation.ServerTotalPayScoreKey(gameRule.GameIndex), serverTotalPayScore, "")
	}
	myhttp.HTTPResponse(w, result, err)

}

func scottergameresult(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	postData := myhttp.PostData(r)
	token := foundation.InterfaceToString(postData["token"])
	scotterid := foundation.InterfaceToInt64(postData["scotterid"])
	luckydrawselect := foundation.InterfaceToInt64(postData["luckydrawselect"])
	err := messagehandle.New()

	if scotterid%10 != 0 {
		err.ErrorCode = code.DataLoss
		err.Msg = "Scotter Game Error"
		myhttp.HTTPResponse(w, "", err)
		return
	}
	if luckydrawselect > 6 {
		err.ErrorCode = code.DataLoss
		err.Msg = "Scotter Game Select Error"
		myhttp.HTTPResponse(w, "", err)
		return
	}

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

	// get thirdparty info data
	ulginfo, err := ulg.GetULGInfo(playerInfo.ID, playerInfo.GameToken)
	if err.ErrorCode != code.OK {
		myhttp.HTTPResponse(w, "", err)
		fmt.Println(ulginfo)
		return
	}

	att := attach.GetAttachByType(playerInfo.ID, gameRule.GameIndex, gameRule.DayScotterGameCountKey, gameRule.IsAttachSaveToDB)
	noCacheAtt := attach.GetAttachByTypeRange(playerInfo.ID, gameRule.GameIndex, scotterid, scotterid+2)
	att = append(att, noCacheAtt...)
	attInfo := gameRule.ConvertToGameAttach(playerID, att)

	scotterInfo, ok := attInfo.ScotterInfos[scotterid]
	if !ok {
		err.ErrorCode = code.FailedPrecondition
		err.Msg = "Scotter Game ID Error"
		myhttp.HTTPResponse(w, "", err)
		return
	} else if scotterInfo.DayScotterGameInfo == 1 {
		err.ErrorCode = code.FailedPrecondition
		err.Msg = "Scotter Isfinish Error"
		myhttp.HTTPResponse(w, "", err)
		return
	}

	betMoney := scotterInfo.FreeGameBetLockMoney
	result, _, _ := scotterGameProcess(playerInfo, ulginfo, att, betMoney, scotterid, luckydrawselect)
	myhttp.HTTPResponse(w, result, err)

}

// AutoRunScotterGameResult at checkout and maintain process not finish scotter
func AutoRunScotterGameResult(playerInfo *playerinfo.Info, ulginfo *ulginfo.Info, att []map[string]interface{}, luckydrawselect int64) {

	if luckydrawselect > 6 {
		messagehandle.ErrorLogPrintln("Scotter Game Error", att, luckydrawselect)
		return
	}

	var extraScotter []map[string]interface{}

	for {

		tmpTargetAtt := gameRule.ConvertToGameAttach(playerInfo.ID, att)

		var scotterid int64
		var Info gameRule.ScotterInfo
		for scotterid, Info = range tmpTargetAtt.ScotterInfos {

			if scotterid%10 != 0 {
				messagehandle.ErrorLogPrintln("Scotter Game Error", Info)
				return
			}

			if Info.DayScotterGameInfo == 1 {
				messagehandle.ErrorLogPrintln("Scotter Isfinish Error", Info)
				return
			}
		}

		betMoney := Info.FreeGameBetLockMoney
		_, newatt, _ := scotterGameProcess(playerInfo, ulginfo, att, betMoney, scotterid, luckydrawselect)

		if len(newatt) > 3 {
			att[0] = newatt[0]
			extraScotter = append(extraScotter, newatt[3:]...)
			att[1], extraScotter = foundation.ArrayShift(extraScotter)
			att[2], extraScotter = foundation.ArrayShift(extraScotter)

		} else if len(extraScotter) > 0 {
			att[1], extraScotter = foundation.ArrayShift(extraScotter)
			att[2], extraScotter = foundation.ArrayShift(extraScotter)

		} else {
			break

		}
	}
}

func scotterGameProcess(playerInfo *playerinfo.Info, ulginfo *ulginfo.Info, att []map[string]interface{}, betMoney, scotterid, luckydrawselect int64) (map[string]interface{}, []map[string]interface{}, map[string]interface{}) {

	serverTotalPayScore := serversetting.GetServerTotalPayScore()
	var totalwinscore int64
	var result map[string]interface{}
	var newatt []map[string]interface{}
	var otherdata map[string]interface{}

	for index, max := 0, 2; index < max; index++ {
		result, newatt, otherdata = gameRule.ScotterGameRequest(playerInfo.ID, betMoney, luckydrawselect, scotterid, att)
		totalwinscore = otherdata["totalwinscore"].(int64)

		if gamelimit.IsInTotalMoneyWinLimit(gameRule.WinScoreLimit, betMoney, totalwinscore) && gamelimit.IsInTotalBetRateWinLimit(gameRule.WinBetRateLimit, betMoney, totalwinscore) {
			break
		}

	}

	playerInfo.Money = playerInfo.Money + totalwinscore - betMoney

	attach.SaveAttach(playerInfo.ID, gameRule.GameIndex, []map[string]interface{}{newatt[0]}, false)
	attach.SaveAttachToDB(playerInfo.ID, gameRule.GameIndex, newatt[1:])
	result["playermoney"] = playerInfo.Money

	ulginfo.TotalWin += totalwinscore
	player.SavePlayerInfo(playerInfo)
	ulg.SaveULGInfo(ulginfo)

	msg := foundation.JSONToString(result)
	msg = strings.ReplaceAll(msg, "\"", "\\\"")
	loginfo := loginfo.New(loginfo.GameResult)
	loginfo.PlayerID = playerInfo.ID
	loginfo.IValue1 = foundation.InterfaceToInt64(result["totalwinscore"])
	loginfo.IValue2 = otherdata["betmoney"].(int64)
	loginfo.IValue3 = scotterid
	loginfo.Msg = msg
	log.SaveLog(loginfo)

	if totalwinscore > 0 {
		serverTotalPayScore += totalwinscore
		serversetting.SetServerTotalPayScore(serverTotalPayScore)
		db.UpdateSetting(foundation.ServerTotalPayScoreKey(gameRule.GameIndex), serverTotalPayScore, "")
	}

	return result, newatt, otherdata
}
