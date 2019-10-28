package myrestful

import (
	"crypto/tls"
	"net/http"

	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/httprouter"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/myhttp"
	"gitlab.com/WebserverByGoGame7/serversetting"
)

var proxyData map[string]myhttp.RESTfulURL

func init() {
	proxyData = make(map[string]myhttp.RESTfulURL)
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
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},

			MaxIdleConns:        50,
			MaxIdleConnsPerHost: 50,
		}
		clientConnect.Client = &http.Client{
			Transport: httptr,
		}
	}
	return clientConnect.Client
}

// ListenProxy client -> Porxy -> processFun
func ListenProxy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	RESTfulInfo := proxyData[r.RequestURI[1:]]
	myhttp.AddHeader(&w)

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

			proxyData[RESTfulURLvalue.URL] = RESTfulURLvalue
			if RESTfulURLvalue.RequestType == "GET" {
				router.GET("/"+RESTfulURLvalue.URL, RESTfulURLvalue.Fun)
			} else if RESTfulURLvalue.RequestType == "POST" {
				router.POST("/"+RESTfulURLvalue.URL, ListenProxy)
			}
			router.OPTIONS("/"+RESTfulURLvalue.URL, myhttp.Option)

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

// HTTPSLisentRun ...
func HTTPSLisentRun(ListenIP, certFile, keyFile string, HandleURL ...[]myhttp.RESTfulURL) (err error) {
	router := httprouter.New()

	for _, RESTfulURLArray := range HandleURL {
		for _, RESTfulURLvalue := range RESTfulURLArray {
			messagehandle.LogPrintf("HTTPListen %v %s\n", RESTfulURLvalue.RequestType, RESTfulURLvalue.URL)

			proxyData[RESTfulURLvalue.URL] = RESTfulURLvalue
			if RESTfulURLvalue.RequestType == "GET" {
				router.GET("/"+RESTfulURLvalue.URL, RESTfulURLvalue.Fun)
			} else if RESTfulURLvalue.RequestType == "POST" {
				router.POST("/"+RESTfulURLvalue.URL, ListenProxy)
			}
			router.OPTIONS("/"+RESTfulURLvalue.URL, myhttp.Option)

		}
	}

	messagehandle.LogPrintln("Server run on", ListenIP)

	err = http.ListenAndServeTLS(ListenIP, certFile, keyFile, router)
	if err != nil {
		messagehandle.ErrorLogPrintln("ListenAndServe", err)
		return err
	}
	return nil
}

func maintain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := messagehandle.New()
	err.ErrorCode = code.Maintain
	err.Msg = "Maintain"
	myhttp.HTTPResponse(w, "", err)
}

// PostRawRequest connect pool
func PostRawRequest(url string, value []byte) []byte {
	return myhttp.HTTPPostRawRequest(connectPool(), url, value)
}
