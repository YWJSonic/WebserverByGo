package account

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"

	"../code"
	"../data"
	"../db"
	"../foundation"
	"../log"
	"../thirdparty/ulg"
)

type AccountInfo struct {
	Account     string
	GameAccount string
	LoginTime   int64

	AccountToken   string // platform AccountToken
	GameToken      string // platform GameToken
	Token          string // Server Token
	ThirdPartyInfo PartyInfo
}

type PartyInfo interface{}

var isInit = false
var mu *sync.RWMutex

// ServiceStart ...
func ServiceStart() []foundation.RESTfulURL {
	var HandleURL []foundation.RESTfulURL

	if isInit {
		return HandleURL
	}

	mu = new(sync.RWMutex)
	isInit = true

	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "account/login", Fun: login})
	return HandleURL
}

func login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var result = make(map[string]interface{})

	postData := foundation.PostData(r)
	token := foundation.InterfaceToString(postData["token"])
	gametypeid := foundation.InterfaceToString(postData["gametypeid"])

	AccountUserInfo, err := ulg.Authorized(token, gametypeid)

	if err.ErrorCode != code.OK {
		err.ErrorCode = code.FailedPrecondition
		foundation.HTTPResponse(w, result, err)
		return
	}

	Info, err := db.GetAccountInfo(strconv.FormatInt(AccountUserInfo.AccountID, 10))

	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, result, err)
		return
	}

	var accountInfo AccountInfo
	if len(Info) < 1 {
		// gameid := db.GetGameID()
		accountInfo = newAccount(strconv.FormatInt(AccountUserInfo.AccountID, 10), foundation.MD5Code(data.AccountEncodeStr+string(AccountUserInfo.AccountID)))
		accountInfo.GameToken = AccountUserInfo.GameToken
		db.NewAccount(accountInfo.Account, accountInfo.GameAccount, accountInfo.LoginTime)
	} else {

		result := Info[0]
		accountInfo = AccountInfo{
			Account:        result["Account"].(string),
			GameAccount:    result["GameAccount"].(string),
			AccountToken:   AccountUserInfo.AccountToken,
			GameToken:      AccountUserInfo.GameToken,
			Token:          foundation.NewToken(Info[0]["Account"].(string)),
			LoginTime:      time.Now().Unix(),
			ThirdPartyInfo: AccountUserInfo,
		}
		db.UpdateAccount(accountInfo.Account, accountInfo.LoginTime)
	}

	loginfo := log.New(log.Login)
	loginfo.Account = accountInfo.GameAccount
	result["usercoinquota"] = AccountUserInfo.UserCoinQuota
	result["gameaccount"] = accountInfo.GameAccount
	result["gametoken"] = accountInfo.GameToken
	result["token"] = accountInfo.Token
	log.SaveLog(loginfo)
	foundation.HTTPResponse(w, result, err)
}

// NewAccount ...
func newAccount(account, gameAccount string) AccountInfo {
	return AccountInfo{
		Account:     account,
		GameAccount: gameAccount,
		Token:       foundation.NewToken(gameAccount),
		LoginTime:   time.Now().Unix(),
	}
}
