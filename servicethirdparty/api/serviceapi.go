package api

import (
	"fmt"
	"net/http"
	"sync"

	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/myhttp"
	"gitlab.com/ServerUtility/playerinfo"
	"gitlab.com/WeberverByGoGame6/apithirdparty/ulg"
	gameRules "gitlab.com/WeberverByGoGame6/gamerule"
	attach "gitlab.com/WeberverByGoGame6/handleattach"
	mycache "gitlab.com/WeberverByGoGame6/handlecache"
	crontab "gitlab.com/WeberverByGoGame6/handlecrontab"
	db "gitlab.com/WeberverByGoGame6/handledb"
	"gitlab.com/WeberverByGoGame6/player"
	"gitlab.com/WeberverByGoGame6/serversetting"
	game "gitlab.com/WeberverByGoGame6/servicegame"

	"gitlab.com/ServerUtility/httprouter"
)

var isInit bool
var mu *sync.RWMutex

// ServiceStart ...
func ServiceStart() []myhttp.RESTfulURL {
	var HandleURL []myhttp.RESTfulURL

	if isInit {
		return HandleURL
	}

	mu = new(sync.RWMutex)
	isInit = true

	// HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "api/cronstart", Fun: CronStart, ConnType: myhttp.Backend})
	// HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "api/cronstop", Fun: CronStop, ConnType: myhttp.Backend})
	// HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "api/cronadd", Fun: CronAdd, ConnType: myhttp.Backend})
	// HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "api/maintainstart", Fun: MaintainStart, ConnType: myhttp.Backend})
	// HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "api/maintainend", Fun: MaintainEnd, ConnType: myhttp.Backend})
	// HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "api/clearcache", Fun: ClearAllCache, ConnType: myhttp.Backend})
	// HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "api/gamerulesset", Fun: GameRulesSet, ConnType: myhttp.Backend})

	return HandleURL
}

// CronStart cron API
func CronStart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	crontab.CronStart()
	fmt.Println("CronStart")
}

// CronStop cron API
func CronStop(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	crontab.CronStop()
	fmt.Println("CronStop")
}

// CronAdd cron API
func CronAdd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("AddCron")
}

// MaintainStart Maintain API
func MaintainStart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	serversetting.EnableMaintain(true)
}

// MaintainEnd Maintain API
func MaintainEnd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	serversetting.EnableMaintain(false)
}

// MaintainCheckout auto checkout ulg
func MaintainCheckout() {
	if !serversetting.IsMaintain() {
		messagehandle.LogPrintln("API Warning: MaintainCheckout not in maintain")
		return
	}

	infos, err := ulg.MaintainULGInfos()
	if err.ErrorCode != code.OK {
		fmt.Println(infos, err)
	}

	for _, ulginfo := range infos {
		RunNotFinishSoctter(ulginfo.PlayerID)
		_, err = ulg.Checkout(&ulginfo, serversetting.GameTypeID) //(ulginfo.AccountToken, ulginfo.GameToken, serversetting.GameTypeID, fmt.Sprint(ulginfo.TotalBet), fmt.Sprint(ulginfo.TotalWin), fmt.Sprint(ulginfo.TotalLost))
		if err.ErrorCode != code.OK {
			messagehandle.ErrorLogPrintln("Crontab MaintainCheckout", err, ulginfo)
		}
	}

	mycache.ClearAllCache()
	db.Game6ClearDBScotterCount()
	db.ULGMaintainCheckOutUpdate()
}

// ClearAllCache clear all cache data
func ClearAllCache(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mycache.ClearAllCache()
}

// GameRulesSet ...
func GameRulesSet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	postData := myhttp.PostData(r)
	configstr := foundation.InterfaceToString(postData["configstr"])
	gameindex := foundation.InterfaceToInt(postData["gameindex"])

	config := foundation.StringToJSON(configstr)
	gameRules.SetInfo(gameindex, config)
}

// RunNotFinishSoctter auto finish scotter
func RunNotFinishSoctter(playerID int64) {

	var ScotterAutoInfo [][]map[string]interface{}
	var err messagehandle.ErrorMsg
	var playerInfo *playerinfo.Info

	ScotterAutoInfo, err = db.Game6AttachGameScotterAutoFinish(playerID)
	if err.ErrorCode != code.OK {
		fmt.Println(ScotterAutoInfo)
		return
	}
	playerAttach := ScotterAutoInfo[0]
	playerGameAccount := ScotterAutoInfo[1]

	playerInfo, err = player.GetPlayerInfoByPlayerID(playerID)
	if err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("RunNotFinishSoctter Cache Get Player Error")

		if len(playerGameAccount) != 1 {
			messagehandle.ErrorLogPrintln("RunNotFinishSoctter GameAccount Count Error")
			return
		}

		playerInfo, err = player.GetPlayerInfoByGameAccount(foundation.InterfaceToString(playerGameAccount[0]["GameAccount"]))
		if err.ErrorCode != code.OK {
			messagehandle.ErrorLogPrintln("RunNotFinishSoctter GameAccount Get Player Error")
			return
		}

	}

	ulginfo, err := ulg.GetULGInfo(playerInfo.ID, playerInfo.GameToken)
	if err.ErrorCode != code.OK {
		fmt.Println(ulginfo)
		return
	}

	for ii, imax := 0, len(playerAttach); ii < imax; ii += 2 {
		scotterAtt := attach.GetAttachByType(playerInfo.ID, gameRules.GameIndex, gameRules.DayScotterGameCountKey, gameRules.IsAttachSaveToDB)
		scotterAtt = append(scotterAtt, playerAttach[ii])
		scotterAtt = append(scotterAtt, playerAttach[ii+1])

		if ivalue, ok := playerAttach[ii]["IValue"]; !ok && ivalue != 0 {
			break
		}
		// scotterCheckID := foundation.InterfaceToInt64(playerAttach[ii]["Type"])
		// scotterID := gameRules.GameInfoKeyToScotterID(scotterCheckID)
		game.AutoRunScotterGameResulttest(playerInfo, ulginfo, scotterAtt, 6)

	}
	// fmt.Println(playerAttach)
}
