package backend

import (
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/WeberverByGoGame5/db"
	"gitlab.com/WeberverByGoGame5/foundation"
	"gitlab.com/WeberverByGoGame5/messagehandle/errorlog"
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

	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "backend/ChangeCheckoutFlag", Fun: ChangeCheckoutFlag, ConnType: foundation.Backend})

	return HandleURL
}

// ChangeCheckoutFlag ...
func ChangeCheckoutFlag(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := errorlog.New()
	postData := foundation.PostData(r)
	gametoken := foundation.InterfaceToString(postData["gametoken"])
	playerid := foundation.InterfaceToInt64(postData["playerid"])

	err = db.ULGSetFlagFalse(gametoken, playerid)
	foundation.HTTPResponse(w, "", err)
}
