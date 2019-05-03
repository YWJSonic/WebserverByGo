package api

import (
	"net/http"
	"sync"

	"../../foundation"
	"../../messagehandle/errorlog"

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

	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "api/changeRTP", Fun: changeRTP})
	// HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "api/addroom", Fun: addroom})
	// HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "api/getRoom", Fun: getRoom})
	return HandleURL
}

func changeRTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()

	err := errorlog.New()
	// postData := foundation.PostData(r)
	// newRTP := postData["RTP"]

	foundation.HTTPResponse(w, "", err)
}

// GetRTP dynamic RTP
func GetRTP() int {
	mu.RLock()
	defer mu.RUnlock()

	return 97
}
