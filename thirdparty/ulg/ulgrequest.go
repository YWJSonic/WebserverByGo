package ulg

import (
	"encoding/json"

	"../../foundation"
	"../../frame/code"
)

type ulgUser struct {
	userID        int
	status        int // 0: empty 1:exchange 2:checkout
	accountName   string
	errorMsg      string
	userName      string // not use, give default value
	token         string // for plant token
	gametoken     string // for game token
	userCoinQuota map[int]code.MoneyBuffer
	coinsetting   ulgCoinSetting
}

type ulgCoinSetting struct {
	cointype string  // money type
	status   int     // enable status
	rate     float32 // exchange rate
	sort     int     // sort index

}

// platform api url
const (
	loginURL      string = "http://54.65.188.126/api/v1/game/login"
	getuserURL    string = "http://54.65.188.126/api/v1/game/get_user"
	authorizedURL string = "http://54.65.188.126/api/v1/game/authorized"
	exchangeURL   string = "http://54.65.188.126/api/v1/game/exchange"
	checkoutURL   string = "http://54.65.188.126/api/v1/game/checkout"
)

// LoginURL ...
func LoginURL() string {
	return loginURL
}

// GetUserURL ...
func GetUserURL() string {
	return getuserURL
}

// AuthorizedURL ...
func AuthorizedURL() string {
	return authorizedURL
}

// ExchangeURL ...
func ExchangeURL() string {
	return exchangeURL
}

// CheckoutURL ...
func CheckoutURL() string {
	return checkoutURL
}

/////// API interface process

// Login ...
func Login(username, password string) map[string]interface{} {
	var info map[string]interface{}
	postData := map[string][]string{
		"username": {username},
		"password": {password},
	}
	jsbyte := foundation.HTTPPostRequest(loginURL, postData)
	if err := json.Unmarshal(jsbyte, &info); err != nil {
		panic(err)
	}
	return info
}

// Getuser client request getplayer info
func Getuser(token, gameid string) map[string]interface{} {
	var info map[string]interface{}
	postData := map[string][]string{
		"token":  {token},
		"gameid": {gameid},
	}
	jsbyte := foundation.HTTPPostRequest(getuserURL, postData)
	if err := json.Unmarshal(jsbyte, &info); err != nil {
		panic(err)
	}
	return info
}

// Authorized ...
func Authorized(token, gameid string) map[string]interface{} {
	var info map[string]interface{}
	// postData := map[string][]string{
	// 	"token":  {token},
	// 	"gameid": {gameid},
	// }
	// jsbyte := foundation.HTTPPostRequest(authorizedURL, postData)
	jsbyte := []byte(`{"result":1,"status":0,"errorMsg":"","userID":123456,"accountName":"qwer","userName":"develop","userCoinQuota":[{"type":"1","amount":4944630},{"type":"2","amount":5000000},{"type":"3","amount":4974000},{"type":"4","amount":0},{"type":"5","amount":5024000}],"gameInfo":[{"type":"1","status":1,"rate":"1.000","sort":1},{"type":"2","status":1,"rate":"0.500","sort":4},{"type":"3","status":1,"rate":"1.000","sort":2},{"type":"4","status":1,"rate":"0.500","sort":3}],"game_token":"1534058582D49E1ECC5D040BBAE11BC07ED6DDD42012"}`)
	if err := json.Unmarshal(jsbyte, &info); err != nil {
		panic(err)
	}
	return info
}

// Exchange ...
func Exchange(gametoken, gameid, token string, cointype, coinamount int) map[string]interface{} {
	var info map[string]interface{}
	// postData := map[string][]string{
	// 	"gametoken":  {gametoken},
	// 	"gameid":     {gameid},
	// 	"token":      {token},
	// 	"cointype":   {string(cointype)},
	// 	"coinamount": {string(coinamount)},
	// }
	// jsbyte := foundation.HTTPPostRequest(exchangeURL, postData)
	jsbyte := []byte(`{"result":1,"errorMsg":"","userCoinQuota":[{"type":"1","amount":1000},{"type":"2","amount":0},{"type":"3","amount":0},{"type":"4","amount":0}],"gameCoin":10000,"gameInfo":[{"type":"1","status":1,"rate":"1.000","sort":1},{"type":"2","status":1,"rate":"0.500","sort":4},{"type":"3","status":1,"rate":"1.000","sort":2},{"type":"4","status":1,"rate":"0.500","sort":3}]}`)
	if err := json.Unmarshal(jsbyte, &info); err != nil {
		panic(err)
	}
	return info
}

// Checkout ...
func Checkout(gametoken, gameid, token, amount string, win, lost int) map[string]interface{} {

	var info map[string]interface{}
	postData := map[string][]string{
		"gametoken": {gametoken},
		"gameid":    {gameid},
		"token":     {token},
		"amount":    {amount},
		"win":       {string(win)},
		"lost":      {string(lost)},
	}
	jsbyte := foundation.HTTPPostRequest(checkoutURL, postData)
	if err := json.Unmarshal(jsbyte, &info); err != nil {
		panic(err)
	}
	return info
}
