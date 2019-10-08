package api

import (
	"net/http"
	"sync"

	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/gamelimit"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/myhttp"
	"gitlab.com/WeberverByGoGame7/apithirdparty/ulg"
	gameRules "gitlab.com/WeberverByGoGame7/gamerule"
	mycache "gitlab.com/WeberverByGoGame7/handlecache"
	db "gitlab.com/WeberverByGoGame7/handledb"
	"gitlab.com/WeberverByGoGame7/serversetting"

	"github.com/julienschmidt/httprouter"
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

// MaintainCheckout auto checkout ulg
func MaintainCheckout() {
	if !serversetting.IsMaintain() {
		messagehandle.LogPrintln("API Warning: MaintainCheckout not in maintain")
		return
	}

	infos, err := ulg.MaintainULGInfos()

	CheckoutErrorPlayerIDs := make([]int64, 0)
	for _, ulginfo := range infos {
		_, err = ulg.Checkout(&ulginfo, serversetting.GameTypeID) //(ulginfo.AccountToken, ulginfo.GameToken, serversetting.GameTypeID, fmt.Sprint(ulginfo.TotalBet), fmt.Sprint(ulginfo.TotalWin), fmt.Sprint(ulginfo.TotalLost))
		if err.ErrorCode != code.OK {
			messagehandle.ErrorLogPrintln("Crontab MaintainCheckout", err, ulginfo)
			CheckoutErrorPlayerIDs = append(CheckoutErrorPlayerIDs, ulginfo.PlayerID)
		}
	}

	if len(CheckoutErrorPlayerIDs) <= 0 {
		db.ULGMaintainCheckOutUpdate()
	} else {
		db.ULGMaintainCheckOutUpdateByPlayerID(CheckoutErrorPlayerIDs)
	}

	serversetting.RefreshDBSetting(gameRules.GameIndex, gamelimit.ServerDayPayDefault)
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
