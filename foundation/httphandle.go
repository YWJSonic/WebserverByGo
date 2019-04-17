package foundation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"../messagehandle/errorlog"
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
