package ulg

import (
	"encoding/json"
	"strconv"
	"strings"

	"../../foundation"
	"../../player"
)

// ULGInfo plant ULG game play data
type ULGInfo struct {
	PlayerID   int64  `json:"PlayerID"`
	GameToken  string `json:"GameToken"` // platform GameToken
	TotalBet   int64  `json:"TotalBet"`
	TotalWin   int64  `json:"TotalWin"`
	TotalLost  int64  `json:"TotalLost"`
	IsCheckOut bool   `json:"CheckOut"`

	AccountToken string `json:"AccountToken"` // Maintan checkout use
}

// ToJSONStr ...
func (ulg ULGInfo) ToJSONStr() string {
	data, _ := json.MarshalIndent(ulg, "", " ")
	STR := string(data)
	STR = strings.ReplaceAll(STR, string(10), ``)
	return STR
}

// UlgResult plant ULG API Result
type UlgResult struct {
	Result        int           `json:"result,omitempty"`
	AccountID     int64         `json:"userID,omitempty"`
	Status        int           `json:"status,omitempty"` // 0: empty 1:exchange 2:checkout
	AccountName   string        `json:"accountName,omitempty"`
	ErrorMsg      string        `json:"errorMsg,omitempty"`
	UserName      string        `json:"userName,omitempty"`   // not use, give default value
	AccountToken  string        `json:"token,omitempty"`      // for plant token
	GameToken     string        `json:"game_token,omitempty"` // for game token
	UserPhone     string        `json:"userPhone,omitempty"`
	GameCoin      int64         `json:"gameCoin,omitempty"`
	UserCoinQuota []CoinQuota   `json:"userCoinQuota,CoinQuota,omitempty"`
	Coinsetting   []CoinSetting `json:"coinsetting,CoinSetting,omitempty"`
	GameInfo      []CoinInfo    `json:"gameInfo,CoinInfo,omitempty"`
	// CheckOutCoin  AmountCoin    `json:"amountCoin,AmountCoin,omitempty"`
}

// UlgCheckOutResult Ulg check result
type UlgCheckOutResult struct {
	Result        int       `json:"result"`
	ErrorMsg      string    `json:"errorMsg"`
	UserCoinQuota CoinQuota `json:"userCoinQuota,CoinQuota"`
}

// PartyAccount ...
func (ulg *UlgResult) PartyAccount() string {
	return foundation.NewAccount("ulg", strconv.FormatInt(ulg.AccountID, 10))
}

// GameAccount ...
func (ulg *UlgResult) GameAccount() string {
	return foundation.NewGameAccount(string(ulg.AccountID))
}

// PartyToken ...
func (ulg *UlgResult) PartyToken() string {
	return ulg.AccountToken
}

// AccountType ...
func (ulg *UlgResult) AccountType() int64 {
	return player.Ulg
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
	CoinType string `json:"type,omitempty"`
	Amount   int64  `json:"amount"`

	Coin1Out     int64 `json:"coin1_out,omitempty"`
	Coin2Out     int64 `json:"coin2_out,omitempty"`
	Coin3Out     int64 `json:"coin3_out,omitempty"`
	Coin4Out     int64 `json:"coin4_out,omitempty"`
	Betting      int64 `json:"betting,omitempty"`
	Win          int64 `json:"win,omitempty"`
	Lost         int64 `json:"lost,omitempty"`
	OutboundTime int64 `json:"outbound_time,omitempty"`
	Status       int   `json:"status,omitempty"`
}

// ToJSONClient ...
func (c CoinQuota) ToJSONClient() map[string]interface{} {
	result := make(map[string]interface{})
	result["cointype"] = c.CoinType
	result["amount"] = c.Amount
	return result
}

// CoinSetting ulg CoinSetting
type CoinSetting struct {
	Cointype string  `json:"cointype"` // money type
	Status   int     `json:"status"`   // enable status
	Rate     float32 `json:"rate"`     // exchange rate
	Sort     int     `json:"sort"`     // sort index

}

// AmountCoin coin check out
type AmountCoin struct {
	Coin1 int64 `json:"coin1"`
	Coin2 int64 `json:"coin2"`
	Coin3 int64 `json:"coin3"`
	Coin4 int64 `json:"coin4"`
}

// platform api url
const (
	loginURL      string = "http://54.65.188.126/api/v1/game/login"
	getuserURL    string = "http://54.65.188.126/api/v1/game/get_user"
	authorizedURL string = "http://54.65.188.126/api/v1/game/authorized"
	exchangeURL   string = "http://54.65.188.126/api/v1/game/exchange"
	checkoutURL   string = "http://54.65.188.126/api/v1/game/checkout"
)
