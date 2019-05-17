package api

import (
	"fmt"
	"net/http"
	"sync"

	"../../crontab"
	"../../data"
	"../../foundation"
	"../../messagehandle/errorlog"
	"../../mycache"
	"../../thirdparty/ulg"

	"github.com/julienschmidt/httprouter"
)

var isInit bool
var mu *sync.RWMutex

// ServiceStart ...
func ServiceStart() []foundation.RESTfulURL {
	var HandleURL []foundation.RESTfulURL

	if isInit {
		return HandleURL
	}

	mu = new(sync.RWMutex)
	isInit = true

	crontab.NewCron("0 20 15 * * *", MaintainCheckout)
	// crontab.NewCron("*/10 * * * * *", MaintainCheckout)

	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "api/cronstart", Fun: CronStart, ConnType: foundation.Backend})
	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "api/cronstop", Fun: CronStop, ConnType: foundation.Backend})
	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "api/cronadd", Fun: CronAdd, ConnType: foundation.Backend})
	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "api/maintainstart", Fun: MaintainStart, ConnType: foundation.Backend})
	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "api/maintainend", Fun: MaintainEnd, ConnType: foundation.Backend})
	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "api/clearcache", Fun: ClearAllCache, ConnType: foundation.Backend})

	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "api/maintaincheckout", Fun: maintaincheckout, ConnType: foundation.Backend})

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
	data.Maintain = true
}

// MaintainEnd Maintain API
func MaintainEnd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	data.Maintain = false
}

// MaintainCheckout auto checkout ulg
func MaintainCheckout() {
	if !data.Maintain {
		errorlog.LogPrintln("API Warning: MaintainCheckout not in maintain")
		return
	}

	infos, err := ulg.MaintainULGInfos()
	fmt.Println(infos, err)
}
func maintaincheckout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	MaintainCheckout()
}

// ClearAllCache clear all cache data
func ClearAllCache(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mycache.ClearAllCache()
}
