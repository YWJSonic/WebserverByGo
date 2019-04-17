package main

import (
	"fmt"
	"net/http"

	"./foundation"
	"./frame/transmission"
	"./game/slotgame"
	"./lobby"
	"./thirdparty/ulg"

	ErrorLog "./messagehandle/errorlog"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// var err error
	// foundation.CacheData, err = cache.New(cache.Config{
	// 	RedisURL:  "127.0.0.1:6379",
	// 	MustRedis: true,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	var initArray [][]transmission.RESTfulURL
	initArray = append(initArray, slotgame.ServiceStart())
	initArray = append(initArray, lobby.ServiceStart())
	initArray = append(initArray, []transmission.RESTfulURL{transmission.RESTfulURL{RequestType: "POST", URL: "ulg/testapi", Fun: testapi}})

	serverRun("192.168.1.10:8000", initArray...)

}

// ServerRun ...
func serverRun(ListenIP string, HandleURL ...[]transmission.RESTfulURL) (err error) {
	router := httprouter.New()

	for _, RESTfulURLArray := range HandleURL {
		for _, RESTfulURLvalue := range RESTfulURLArray {
			fmt.Printf("HTTPListen %v %s\n", RESTfulURLvalue.RequestType, RESTfulURLvalue.URL)

			if RESTfulURLvalue.RequestType == "GET" {
				router.GET("/"+RESTfulURLvalue.URL, RESTfulURLvalue.Fun)
			} else if RESTfulURLvalue.RequestType == "POST" {
				router.POST("/"+RESTfulURLvalue.URL, RESTfulURLvalue.Fun)
			}
		}
	}

	fmt.Println("Server run on", ListenIP)
	err = http.ListenAndServe(ListenIP, router)
	if err != nil {
		ErrorLog.ErrorLogPrintln("ListenAndServe", err)
		return err
	}
	return err
}

func testapi(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var Data = foundation.PostData(r)
	var err ErrorLog.ErrorMsg

	var request map[string]interface{}
	switch Data["api"] {
	case "login":
		request = ulg.Login("game_test", "123456")
	case "getuser":
		request = ulg.Getuser("game_test", "123456")
	case "authorized":
		request = ulg.Authorized("fwe", "qwer1")
	case "exchange":
		request = ulg.Exchange("fwe", "qwer1", "qwef", 123, 123)
	case "checkout":
		request = ulg.Checkout("fwe", "qwer1", "fwe", "qwer1", 123, 123)
	}

	foundation.HTTPResponse(w, request, err.New())

}
