package foundation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"../code"
	"../data"
	"../messagehandle/errorlog"
	"github.com/julienschmidt/httprouter"
)

var ProxyData map[string]RESTfulURL

func init() {
	ProxyData = make(map[string]RESTfulURL)
}

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
	errorlog.LogPrintf("%s", result)
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
	return result

}

// HTTPResponse Respond to cliente
func HTTPResponse(httpconn http.ResponseWriter, data interface{}, err errorlog.ErrorMsg) {
	result := make(map[string]interface{})
	result["data"] = data
	result["error"] = err
	fmt.Fprint(httpconn, JSONToString(result))
}

// PostData get http post data
func PostData(r *http.Request) map[string]interface{} {
	data := map[string]interface{}{}
	contentType := r.Header.Get("Content-type")

	if contentType == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		v := r.Form
		postdata := v.Get("POST")
		if err := json.Unmarshal([]byte(postdata), &data); err != nil {
			panic(err)
		}

	} else {
		d := json.NewDecoder(r.Body)
		err := d.Decode(&data)
		if err != nil {
			panic(err)
		}
	}

	return data
}

// HTTPLisentRun ...
func HTTPLisentRun(ListenIP string, HandleURL ...[]RESTfulURL) (err error) {
	router := httprouter.New()

	for _, RESTfulURLArray := range HandleURL {
		for _, RESTfulURLvalue := range RESTfulURLArray {
			errorlog.LogPrintf("HTTPListen %v %s\n", RESTfulURLvalue.RequestType, RESTfulURLvalue.URL)

			ProxyData[RESTfulURLvalue.URL] = RESTfulURLvalue
			if RESTfulURLvalue.RequestType == "GET" {
				router.GET("/"+RESTfulURLvalue.URL, RESTfulURLvalue.Fun)
			} else if RESTfulURLvalue.RequestType == "POST" {
				// router.POST("/"+RESTfulURLvalue.URL, RESTfulURLvalue.Fun)
				router.POST("/"+RESTfulURLvalue.URL, ListenProxy)
			}
			router.OPTIONS("/"+RESTfulURLvalue.URL, option)

		}
	}

	errorlog.LogPrintln("Server run on", ListenIP)

	// HTTPS Server
	// ListenAndServeTLS

	err = http.ListenAndServe(ListenIP, router)
	if err != nil {
		errorlog.ErrorLogPrintln("ListenAndServe", err)
		return err
	}
	return err
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

// ListenProxy client -> Porxy -> processFun
func ListenProxy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	addHeader(&w)

	if data.Maintain {
		maintain(w, r, ps)
	} else {
		ProxyData[r.URL.Path[1:]].Fun(w, r, ps)
	}
}

func addHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	// (*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*w).Header().Set("Content-Type", "application/json")
	// (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func maintain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := errorlog.New()
	err.ErrorCode = code.Maintain
	HTTPResponse(w, "", err)
}
