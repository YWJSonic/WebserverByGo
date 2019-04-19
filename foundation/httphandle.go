package foundation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"../frame/code"
	"../frame/transmission"
	"../messagehandle/errorlog"
	ErrorLog "../messagehandle/errorlog"
	"github.com/julienschmidt/httprouter"
)

// HTTPGet ...
func HTTPGet(ip string, values map[string][]string) []byte {
	res, err := http.Get(ip)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", result)
	return result
}

// HTTPPostRequest ...
func HTTPPostRequest(ip string, values map[string][]string) []byte {
	// res, err := http.Post(ip, "application/x-www-form-urlencoded", strings.NewReader("name=cjb"))
	res, err := http.PostForm(ip, values)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", result)
	return result

}

// HTTPResponse Respond to cliente
func HTTPResponse(httpconn http.ResponseWriter, data interface{}, err errorlog.ErrorMsg) {
	resoult := make(map[string]interface{})
	resoult["data"] = data
	resoult["error"] = err
	fmt.Fprint(httpconn, JSONToString(resoult))
}

// PostData get http post data
func PostData(r *http.Request) map[string]interface{} {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	v := r.Form

	postdata := v.Get("POST")
	date := map[string]interface{}{}
	if err := json.Unmarshal([]byte(postdata), &date); err != nil {
		panic(err)
	}
	return date
}

// HTTPLisentRun ...
func HTTPLisentRun(ListenIP string, HandleURL ...[]transmission.RESTfulURL) (err error) {
	router := httprouter.New()

	for _, RESTfulURLArray := range HandleURL {
		for _, RESTfulURLvalue := range RESTfulURLArray {
			fmt.Printf("HTTPListen %v %s\n", RESTfulURLvalue.RequestType, RESTfulURLvalue.URL)

			ProxyData[RESTfulURLvalue.URL] = RESTfulURLvalue
			if RESTfulURLvalue.RequestType == "GET" {
				router.GET("/"+RESTfulURLvalue.URL, RESTfulURLvalue.Fun)
			} else if RESTfulURLvalue.RequestType == "POST" {
				// router.POST("/"+RESTfulURLvalue.URL, RESTfulURLvalue.Fun)
				router.POST("/"+RESTfulURLvalue.URL, ListenProxy)
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

// ListenProxy client -> Porxy -> processFun
func ListenProxy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if ServerSetting.Maintain {
		maintain(w, r, ps)
	} else {
		ProxyData[r.URL.Path[1:]].Fun(w, r, ps)
	}
}

func maintain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := ErrorLog.New()
	err.ErrorCode = code.Maintain
	HTTPResponse(w, "", err)
}
