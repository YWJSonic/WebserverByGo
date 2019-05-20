package account

import (
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"

	"../code"
	"../data"
	"../db"
	"../foundation"
	"../log"
	"../messagehandle/errorlog"
	"../mycache"
	"../player"
	"../thirdparty"
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

	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "account/login", Fun: login, ConnType: foundation.Client})
	return HandleURL
}
func login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var result = make(map[string]interface{})

	postData := foundation.PostData(r)
	thirdpartyid := foundation.InterfaceToInt(postData["thirdpartyid"])
	accounttoken := foundation.InterfaceToString(postData["accounttoken"])
	gametypeid := foundation.InterfaceToString(postData["gametypeid"])

	err := errorlog.New()
	if gametypeid != data.GameTypeID {
		err.ErrorCode = code.GameTypeError
		err.Msg = "GameTypeError"
		foundation.HTTPResponse(w, "", err)
		return
	}

	// mutility party switch
	var iPratyAccount thirdparty.IPratyAccount
	switch thirdpartyid {
	case thirdparty.Self:
	case thirdparty.Ulg:
	default:
		UserInfo, err := ulg.GetUser(accounttoken, gametypeid)
		if err.ErrorCode != code.OK {
			err.ErrorCode = code.FailedPrecondition
			foundation.HTTPResponse(w, "", err)
			return
		}

		iPratyAccount = &UserInfo
		result["partyinfo"] = UserInfo.UserCoinQuota
	}

	Info, err := db.GetAccountInfo(iPratyAccount.PartyAccount())
	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return
	}

	var accountInfo player.AccountInfo
	if len(Info) < 1 {
		accountInfo = player.NewAccountInfo(iPratyAccount.PartyAccount(), iPratyAccount.GameAccount())
		db.NewAccount(accountInfo.Account, accountInfo.GameAccount, accountInfo.LoginTime)

	} else {
		result := Info[0]
		accountInfo = player.NewAccountInfo(result["Account"].(string), result["GameAccount"].(string))
		db.UpdateAccount(accountInfo.Account, accountInfo.LoginTime)
	}

	mycache.SetAccountInfo(accountInfo.GameAccount, accountInfo.ToJSONStr())
	mycache.SetToken(accountInfo.GameAccount, accountInfo.Token)
	result["gameaccount"] = accountInfo.GameAccount
	result["token"] = accountInfo.Token

	loginfo := log.New(log.Login)
	loginfo.Account = accountInfo.GameAccount
	log.SaveLog(loginfo)
	foundation.HTTPResponse(w, result, err)
}
