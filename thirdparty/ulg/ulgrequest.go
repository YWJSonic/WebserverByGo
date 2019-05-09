package ulg

import (
	"encoding/json"
	"fmt"
	"sync"

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
	UserPhone     string        `json:"userPhone"`
	GameCoin      int64         `json:"gameCoin"`
	UserCoinQuota []CoinQuota   `json:"userCoinQuota,CoinQuota"`
	Coinsetting   []CoinSetting `json:"coinsetting,CoinSetting"`
	GameInfo      []CoinInfo    `json:"gameInfo,CoinInfo"`
}

// CoinInfo Coin rate info
type CoinInfo struct {
	CoinType string  `json:"type"`
	Status   int     `json:"status"`
	Rate     float32 `json:"rate"`
	Sort     int     `json:"sort"`
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
// func Login(username, password string) map[string]interface{} {
// 	var info map[string]interface{}
// 	postData := map[string][]string{
// 		"username": {username},
// 		"password": {password},
// 	}
// 	jsbyte := foundation.HTTPPostRequest(loginURL, postData)
// 	if err := json.Unmarshal(jsbyte, &info); err != nil {
// 		panic(err)
// 	}
// 	return info
// }

// GetUser client request getplayer info
func GetUser(token, gameid string) (UlgResult, errorlog.ErrorMsg) {
	var info UlgResult
	err := errorlog.New()
	postData := map[string][]string{
		"token":   {token},
		"game_id": {gameid},
	}
	jsbyte := foundation.HTTPPostRequest(getuserURL, postData)
	if jserr := json.Unmarshal(jsbyte, &info); jserr != nil {
		err.ErrorCode = code.GetUserError
		err.Msg = "UserFormatError"
	}

	if info.Result == 1 {
		err.ErrorCode = code.OK
	} else {
		err.ErrorCode = code.GetUserError
		err.Msg = info.ErrorMsg
	}
	return info, err
}

var count int
var mu *sync.RWMutex

func init() {

	mu = new(sync.RWMutex)
}

// Authorized ...
func Authorized(token, gametypeid string) (UlgResult, errorlog.ErrorMsg) {
	mu.Lock()
	defer mu.Unlock()
	var info UlgResult
	err := errorlog.New()
	postData := map[string][]string{
		"token":   {token},
		"game_id": {gametypeid},
	}
	jsbyte := foundation.HTTPPostRequest(authorizedURL, postData)
	if jserr := json.Unmarshal(jsbyte, &info); jserr != nil {
		err.ErrorCode = code.AuthorizedError
		err.Msg = "AuthorizedFormatError"
	}

	if info.Result == 1 {
		err.ErrorCode = code.OK
	} else {
		err.ErrorCode = code.AuthorizedError
		err.Msg = info.ErrorMsg
	}
	return info, err
}

// Exchange ...
func Exchange(gametoken, gametypeid, accounttoken string, cointype, coinamount int) (UlgResult, errorlog.ErrorMsg) { // map[string]interface{} {
	var info UlgResult
	err := errorlog.New()
	postData := map[string][]string{
		"game_token":  {gametoken},
		"game_id":     {gametypeid},
		"token":       {accounttoken},
		"coin_type":   {fmt.Sprint(cointype)},
		"coin_amount": {fmt.Sprint(coinamount)},
	}
	jsbyte := foundation.HTTPPostRequest(exchangeURL, postData)
	if jserr := json.Unmarshal(jsbyte, &info); jserr != nil {
		err.ErrorCode = code.ExchangeError
		err.Msg = "ExchangeFormatError"
	}

	if info.Result == 1 {
		err.ErrorCode = code.OK
	} else {
		err.ErrorCode = code.ExchangeError
		err.Msg = info.ErrorMsg
	}
	return info, err
}

// Checkout ...
func Checkout(accounttoken, gametoken, gameid, amount string, totalwin, totalost int64) (UlgResult, errorlog.ErrorMsg) {
	var info UlgResult
	err := errorlog.New()
	postData := map[string][]string{
		"game_token": {gametoken},
		"game_id":    {gameid},
		"token":      {accounttoken},
		"amount":     {amount},
		"win":        {fmt.Sprint(int(600))},
		"lost":       {fmt.Sprint(int(300))},
	}
	jsbyte := foundation.HTTPPostRequest(checkoutURL, postData)
	if jserr := json.Unmarshal(jsbyte, &info); jserr != nil {
		err.ErrorCode = code.ExchangeError
		err.Msg = "ExchangeError"
	}
	//{"data":{"result":0,"userID":0,"status":0,"accountName":"","errorMsg":"\u0008checkout - 無遊戲紀錄","userName":"","token":"","game_token":"","gameCoin":0,"userCoinQuota":null,"coinsetting":null,"gameInfo":null},"error":{"ErrorCode":20,"Msg":"ExchangeError"}}
	if info.Result == 1 {
		err.ErrorCode = code.OK
	} else {
		err.ErrorCode = code.ExchangeError
		err.Msg = info.ErrorMsg
	}
	return info, err
}
