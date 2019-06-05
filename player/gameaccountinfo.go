package player

import (
	"encoding/json"
	"strings"
	"time"

	"gitlab.com/WeberverByGo/foundation"
)

// Account Type
const (
	None = iota
	Guest
	Self
	Ulg
)

// IPratyAccount thirdparty api interface
type IPratyAccount interface {
	PartyAccount() string
	GameAccount() string
	AccountType() int64
}

// AccountInfo ...
type AccountInfo struct {
	Account     string `json:"Account"`
	GameAccount string `json:"GameAccount"`
	AccountType int64  `json:"AccountType"`
	LoginTime   int64  `json:"LoginTime"`

	AccountToken string `json:"AccountToken"` // platform AccountToken
	Token        string `json:"Token"`        // Server Token
}

// PartyInfo ThirdPartyInfo
type PartyInfo interface{}

// NewAccountInfo ...
func NewAccountInfo(account, gameAccount string, accountType int64) AccountInfo {
	return AccountInfo{
		Account:     account,
		GameAccount: gameAccount,
		Token:       foundation.NewToken(account),
		LoginTime:   time.Now().Unix(),
		AccountType: accountType,
	}
}

// ToJSONStr account info to json string
func (a AccountInfo) ToJSONStr() string {
	data, _ := json.MarshalIndent(a, "", " ")
	STR := string(data)
	STR = strings.ReplaceAll(STR, string(10), ``)
	return STR
}
