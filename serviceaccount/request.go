package login

import (
	"net/http"
	"sync"

	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/httprouter"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/myhttp"
	"gitlab.com/WeberverByGoBase/apithirdparty"
	log "gitlab.com/WeberverByGoBase/handlelog"
	"gitlab.com/WeberverByGoBase/serversetting"
)

var isInit = false
var mu *sync.RWMutex

// ServiceStart ...
func ServiceStart() []myhttp.RESTfulURL {
	var HandleURL []myhttp.RESTfulURL

	if isInit {
		return HandleURL
	}

	mu = new(sync.RWMutex)
	isInit = true

	HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "account/login", Fun: login, ConnType: myhttp.Client})
	return HandleURL
}
func login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var result = make(map[string]interface{})

	postData := myhttp.PostData(r)
	logintype := foundation.InterfaceToInt(postData["logintype"])
	gametypeid := foundation.InterfaceToString(postData["gametypeid"])

	err := messagehandle.New()
	if err = foundation.CheckGameType(serversetting.GameTypeID, gametypeid); err.ErrorCode != code.OK {
		myhttp.HTTPResponse(w, "", err)
		return
	}

	result, accountInfo, err := apithirdparty.GetAccount(logintype, foundation.InterfaceToString(postData["accounttoken"]), gametypeid)
	if err.ErrorCode != code.OK {
		myhttp.HTTPResponse(w, "account login", err)
		return
	}

	result["gameaccount"] = accountInfo.GameAccount
	result["token"] = accountInfo.Token
	result["serversetting"] = serversetting.New()

	log.AcouuntLogin(accountInfo.GameAccount)
	myhttp.HTTPResponse(w, result, err)
}
