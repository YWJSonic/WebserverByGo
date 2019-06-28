package api

import (
	"fmt"
	"net/http"
	"sync"

	"gitlab.com/ServerUtility/myhttp"
	"gitlab.com/WeberverByGo/code"
	"gitlab.com/WeberverByGo/crontab"
	"gitlab.com/WeberverByGo/data"
	"gitlab.com/WeberverByGo/db"
	"gitlab.com/WeberverByGo/foundation"
	"gitlab.com/WeberverByGo/messagehandle/errorlog"
	"gitlab.com/WeberverByGo/mycache"
	"gitlab.com/WeberverByGo/thirdparty/ulg"
	gameRules "gitlab.com/game7"

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

	HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "api/cronstart", Fun: CronStart, ConnType: myhttp.Backend})
	HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "api/cronstop", Fun: CronStop, ConnType: myhttp.Backend})
	HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "api/cronadd", Fun: CronAdd, ConnType: myhttp.Backend})
	HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "api/maintainstart", Fun: MaintainStart, ConnType: myhttp.Backend})
	HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "api/maintainend", Fun: MaintainEnd, ConnType: myhttp.Backend})
	HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "api/clearcache", Fun: ClearAllCache, ConnType: myhttp.Backend})
	HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "api/gamerulesset", Fun: GameRulesSet, ConnType: myhttp.Backend})

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
	data.EnableMaintain(true)
}

// MaintainEnd Maintain API
func MaintainEnd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	data.EnableMaintain(false)
}

// MaintainCheckout auto checkout ulg
func MaintainCheckout() {
	if !data.IsMaintain() {
		errorlog.LogPrintln("API Warning: MaintainCheckout not in maintain")
		return
	}

	infos, err := ulg.MaintainULGInfos()
	fmt.Println(infos, err)

	for _, ulginfo := range infos {
		_, err = ulg.Checkout(&ulginfo, data.GameTypeID) //(ulginfo.AccountToken, ulginfo.GameToken, data.GameTypeID, fmt.Sprint(ulginfo.TotalBet), fmt.Sprint(ulginfo.TotalWin), fmt.Sprint(ulginfo.TotalLost))
		if err.ErrorCode != code.OK {
			errorlog.ErrorLogPrintln("Crontab MaintainCheckout", err, ulginfo)
		}

		mycache.ClearAllCache()
	}

	db.ULGMaintainCheckOutUpdate()
}

// ClearAllCache clear all cache data
func ClearAllCache(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mycache.ClearAllCache()
}

func GameRulesSet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	postData := myhttp.PostData(r)
	configstr := foundation.InterfaceToString(postData["configstr"])
	gameindex := foundation.InterfaceToInt(postData["gameindex"])

	config := foundation.StringToJSON(configstr)
	gameRules.SetInfo(gameindex, config)
}
