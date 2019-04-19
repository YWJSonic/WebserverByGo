package api

import (
	"net/http"
	"sync"

	"../../frame/transmission"

	"../../foundation"
	"../../messagehandle/errorlog"

	"github.com/julienschmidt/httprouter"
)

var isInit bool
var mu *sync.RWMutex

// ServiceStart ...
func ServiceStart() []transmission.RESTfulURL {
	var HandleURL []transmission.RESTfulURL

	if isInit {
		return HandleURL
	}

	mu = new(sync.RWMutex)
	isInit = true

	HandleURL = append(HandleURL, transmission.RESTfulURL{RequestType: "POST", URL: "webservice/changeRTP", Fun: changeRTP})
	return HandleURL
}

func changeRTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := errorlog.New()
	// postData := foundation.PostData(r)
	// newRTP := postData["RTP"]

	foundation.HTTPResponse(w, "", err)
}

// GetRTP dynamic RTP
func GetRTP() int {
	return 97
}
