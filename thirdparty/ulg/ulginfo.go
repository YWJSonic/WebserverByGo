package ulg

import (
	"encoding/json"
	"strconv"
	"strings"

	"../../foundation"
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
	UserCoinQuota CoinQuota     `json:"userCoinQuota,CoinQuota"`
	Coinsetting   []CoinSetting `json:"coinsetting,CoinSetting"`
	GameInfo      []CoinInfo    `json:"gameInfo,CoinInfo"`
	CheckOutCoin  AmountCoin    `json:"amountCoin,AmountCoin"`
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
	Status       int    `json:"status"`
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
