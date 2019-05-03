package ulg

import (
	"encoding/json"
	"fmt"

	"../../code"
	"../../foundation"
	"../../messagehandle/errorlog"
)

// UlgResult ...
type UlgResult struct {
	Result        int           `json:"result"`
	AccountID     int64         `json:"userID"`
	Status        int           `json:"status"` // 0: empty 1:exchange 2:checkout
	AccountName   string        `json:"accountName"`
	ErrorMsg      string        `json:"errorMsg"`
	UserName      string        `json:"userName"`   // not use, give default value
	AccountToken  string        `json:"token"`      // for plant token
	GameToken     string        `json:"game_token"` // for game token
	GameCoin      int64         `json:"gameCoin"`
	UserCoinQuota []CoinQuota   `json:"userCoinQuota,CoinQuota"`
	Coinsetting   []CoinSetting `json:"coinsetting,CoinSetting"`
	GameInfo      []CoinInfo    `json:"gameInfo,CoinInfo"`
}

// CoinInfo Coin rate info
type CoinInfo struct {
	CoinType string `json:"type"`
	Status   int    `json:"status"`
	Rate     string `json:"rate"`
	Sort     int    `json:"sort"`
}

// CoinQuota ulg user CoinQuota
type CoinQuota struct {
	CoinType string `json:"type"`
	Amount   int64  `json:"amount"`

	Coin1Out     int64  `json:"coin1_out"`
	Coin2Out     int64  `json:"coin2_out"`
	Coin3Out     int64  `json:"coin3_out"`
	Coin4Out     int64  `json:"coin4_out"`
	Betting      string `json:"betting"`
	Win          string `json:"win"`
	Lost         string `json:"lost"`
	OutboundTime int64  `json:"outbound_time"`
}

// CoinSetting ulg CoinSetting
type CoinSetting struct {
	Cointype string  `json:"cointype"` // money type
	Status   int     `json:"status"`   // enable status
	Rate     float32 `json:"rate"`     // exchange rate
	Sort     int     `json:"sort"`     // sort index

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

var count int

// Authorized ...
func Authorized(token, gametypeid string) (UlgResult, errorlog.ErrorMsg) {
	var info UlgResult
	var err errorlog.ErrorMsg
	// postData := map[string][]string{
	// 	"token":  {token},
	// 	"game_id": {gametypeid},
	// }
	// jsbyte := foundation.HTTPPostRequest(authorizedURL, postData)
	count++
	jsbyte := []byte(fmt.Sprintf(`
	{"result":1,
	"status":0,
	"errorMsg":"",
	"userID":%d,
	"accountName":"qwer",
	"userName":"develop",
	"userCoinQuota":[{"type":"1","amount":4944630},{"type":"2","amount":5000000},{"type":"3","amount":4974000},{"type":"4","amount":0},{"type":"5","amount":5024000}],"gameInfo":[{"type":"1","status":1,"rate":"1.000","sort":1},{"type":"2","status":1,"rate":"0.500","sort":4},{"type":"3","status":1,"rate":"1.000","sort":2},{"type":"4","status":1,"rate":"0.500","sort":3}],
	"game_token":"1534058582D49E1ECC5D040BBAE11BC07ED6DDD42012"}`, count))
	if jserr := json.Unmarshal(jsbyte, &info); jserr != nil {
		panic("AuthorizedFormatError")
	}

	if info.Result == 1 {
		err.ErrorCode = code.OK
	} else {
		err.ErrorCode = code.AuthorizedError
	}
	err.Msg = "AuthorizedError"
	return info, err
}

// Exchange ...
func Exchange(gametoken, gametypeid, accounttoken string, cointype, coinamount int) (UlgResult, errorlog.ErrorMsg) { // map[string]interface{} {
	var info UlgResult
	var err errorlog.ErrorMsg
	// postData := map[string][]string{
	// 	"gametoken":  {gametoken},
	// 	"gameid":     {gametypeid},
	// 	"token":      {accounttoken},
	// 	"cointype":   {string(cointype)},
	// 	"coinamount": {string(coinamount)},
	// }
	// jsbyte := foundation.HTTPPostRequest(exchangeURL, postData)
	jsbyte := []byte(`{"result":1,"errorMsg":"","userCoinQuota":[{"type":"1","amount":1000},{"type":"2","amount":0},{"type":"3","amount":0},{"type":"4","amount":0}],"gameCoin":10000,"gameInfo":[{"type":"1","status":1,"rate":"1.000","sort":1},{"type":"2","status":1,"rate":"0.500","sort":4},{"type":"3","status":1,"rate":"1.000","sort":2},{"type":"4","status":1,"rate":"0.500","sort":3}]}`)
	if jserr := json.Unmarshal(jsbyte, &info); jserr != nil {
		panic("ExchangeFormatError")
	}

	if info.Result == 1 {
		err.ErrorCode = code.OK
	} else {
		err.ErrorCode = code.ExchangeError
	}
	err.Msg = "ExchangeError"
	return info, err
}

// Checkout ...
func Checkout(gametoken, gameid, token, amount string, win, lost int) UlgResult {

	var info UlgResult
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
