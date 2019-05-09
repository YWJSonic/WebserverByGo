package account

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/julienschmidt/httprouter"

	"../code"
	"../data"
	"../db"
	"../foundation"
	"../log"
	"../mycache"
	"../player"
	"../thirdparty/ulg"
)

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
	accounttoken := foundation.InterfaceToString(postData["accounttoken"])
	gametypeid := foundation.InterfaceToString(postData["gametypeid"])

	UserInfo, err := ulg.GetUser(accounttoken, gametypeid)

	if err.ErrorCode != code.OK {
		err.ErrorCode = code.FailedPrecondition
		foundation.HTTPResponse(w, "", err)
		return
	}
	AccountStr := foundation.NewAccount("ulg", strconv.FormatInt(UserInfo.AccountID, 10))
	Info, err := db.GetAccountInfo(AccountStr)

	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return
	}

	var accountInfo player.AccountInfo
	if len(Info) < 1 {
		accountInfo = player.NewAccount(AccountStr, foundation.MD5Code(data.AccountEncodeStr+string(UserInfo.AccountID)))
		db.NewAccount(accountInfo.Account, accountInfo.GameAccount, accountInfo.LoginTime)
	} else {

		result := Info[0]
		accountInfo = player.AccountInfo{
			Account:      result["Account"].(string),
			GameAccount:  result["GameAccount"].(string),
			AccountToken: UserInfo.AccountToken,
			Token:        foundation.NewToken(Info[0]["Account"].(string)),
			LoginTime:    foundation.ServerNowTime(),
			// ThirdPartyInfo: AccountUserInfo,
		}
		db.UpdateAccount(accountInfo.Account, accountInfo.LoginTime)
	}

	mycache.SetAccountInfo(accountInfo.GameAccount, accountInfo.ToJSONStr())
	mycache.SetToken(accountInfo.GameAccount, accountInfo.Token)
	result["usercoinquota"] = UserInfo.UserCoinQuota
	result["gameaccount"] = accountInfo.GameAccount
	result["token"] = accountInfo.Token

	loginfo := log.New(log.Login)
	loginfo.Account = accountInfo.GameAccount
	log.SaveLog(loginfo)
	foundation.HTTPResponse(w, result, err)
}
