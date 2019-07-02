package myrestful

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/myhttp"
	"gitlab.com/WeberverByGo/serversetting"
)

var ProxyData map[string]myhttp.RESTfulURL

func init() {
	ProxyData = make(map[string]myhttp.RESTfulURL)
}

type httpClient struct {
	Client *http.Client
}

var clientConnect *httpClient

// HttpClient http get http request connect pool
func connectPool() *http.Client {
	if clientConnect == nil {
		clientConnect = new(httpClient)
		httptr := &http.Transport{

			MaxIdleConns:        50,
			MaxIdleConnsPerHost: 50,
		}
		clientConnect.Client = &http.Client{
			Transport: httptr,
		}
	}
	return clientConnect.Client
}

// PostRawRequest connect pool
func PostRawRequest(url string, value []byte) []byte {
	return myhttp.HTTPPostRawRequest(connectPool(), url, value)
}

// ListenProxy client -> Porxy -> processFun
func ListenProxy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	RESTfulInfo := ProxyData[r.RequestURI[1:]]
	addHeader(&w)

	switch RESTfulInfo.ConnType {
	case myhttp.Client:
		if serversetting.IsMaintain() {
			maintain(w, r, ps)
		} else {
			RESTfulInfo.Fun(w, r, ps)
		}
	case myhttp.Backend:
		RESTfulInfo.Fun(w, r, ps)
	}
}

// HTTPLisentRun ...
func HTTPLisentRun(ListenIP string, HandleURL ...[]myhttp.RESTfulURL) (err error) {
	router := httprouter.New()

	for _, RESTfulURLArray := range HandleURL {
		for _, RESTfulURLvalue := range RESTfulURLArray {
			messagehandle.LogPrintf("HTTPListen %v %s\n", RESTfulURLvalue.RequestType, RESTfulURLvalue.URL)

			ProxyData[RESTfulURLvalue.URL] = RESTfulURLvalue
			if RESTfulURLvalue.RequestType == "GET" {
				router.GET("/"+RESTfulURLvalue.URL, RESTfulURLvalue.Fun)
			} else if RESTfulURLvalue.RequestType == "POST" {
				router.POST("/"+RESTfulURLvalue.URL, ListenProxy)
			}
			router.OPTIONS("/"+RESTfulURLvalue.URL, option)

		}
	}

	messagehandle.LogPrintln("Server run on", ListenIP)

	err = http.ListenAndServe(ListenIP, router)
	if err != nil {
		messagehandle.ErrorLogPrintln("ListenAndServe", err)
		return err
	}
	return nil
}

// HTTPResponse Respond to cliente
func HTTPResponse(httpconn http.ResponseWriter, data interface{}, err messagehandle.ErrorMsg) {
	result := make(map[string]interface{})
	result["data"] = data
	result["error"] = err
	fmt.Fprint(httpconn, foundation.JSONToString(result))
}

func maintain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := messagehandle.New()
	err.ErrorCode = code.Maintain
	err.Msg = "Maintain"
	HTTPResponse(w, "", err)
}

func addHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	// (*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*w).Header().Set("Content-Type", "application/json")
	// (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func option(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers := w.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")
	headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
	headers.Add("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
	w.WriteHeader(http.StatusOK)
}
