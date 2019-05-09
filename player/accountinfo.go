package player

import (
	"encoding/json"
	"strings"
	"time"

	"../foundation"
)

// AccountInfo ...
type AccountInfo struct {
	Account     string `json:"Account"`
	GameAccount string `json:"GameAccount"`
	LoginTime   int64  `json:"LoginTime"`

	AccountToken string `json:"AccountToken"` // platform AccountToken
	GameToken    string `json:"GameToken"`    // platform GameToken
	Token        string `json:"Token"`        // Server Token

	IsGameTokenCheckout bool `json:"IsGameTokenCheckout"` // 
}

// PartyInfo ThirdPartyInfo
type PartyInfo interface{}

// NewAccount ...
func NewAccount(account, gameAccount string) AccountInfo {
	return AccountInfo{
		Account:     account,
		GameAccount: gameAccount,
		Token:       foundation.NewToken(gameAccount),
		LoginTime:   time.Now().Unix(),
	}
}

// ToJSONStr account info to json string
func (a AccountInfo) ToJSONStr() string {
	data, _ := json.MarshalIndent(a, "", " ")
	STR := string(data)
	STR = strings.ReplaceAll(STR, string(10), ``)
	return STR
}
