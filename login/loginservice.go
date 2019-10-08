package login

import (
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/WeberverByGoGame5/code"
	"gitlab.com/WeberverByGoGame5/db"
	"gitlab.com/WeberverByGoGame5/foundation"
	"gitlab.com/WeberverByGoGame5/messagehandle/errorlog"
	"gitlab.com/WeberverByGoGame5/messagehandle/log"
	"gitlab.com/WeberverByGoGame5/player"
	"gitlab.com/WeberverByGoGame5/setting"
	"gitlab.com/WeberverByGoGame5/thirdparty/guest"
	"gitlab.com/WeberverByGoGame5/thirdparty/ulg"
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
	logintype := foundation.InterfaceToInt(postData["logintype"])
	gametypeid := foundation.InterfaceToString(postData["gametypeid"])
	errorlog.LogPrintln("Login postData:", postData)

	err := errorlog.New()
	if err = foundation.CheckGameType(gametypeid); err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return
	}

	// mutility party switch
	var iPratyAccount player.IPratyAccount
	switch logintype {
	case player.Guest:
		account := foundation.InterfaceToString(postData["account"])
		guestinfo := guest.GuestInfo{
			Account: account,
		}

		iPratyAccount = &guestinfo
	case player.Self:
	case player.Ulg:
		accounttoken := foundation.InterfaceToString(postData["accounttoken"])
		errorlog.LogPrintln("ULG GetUser:", accounttoken, gametypeid)
		UserInfo, err := ulg.GetUser(accounttoken, gametypeid)
		if err.ErrorCode != code.OK {
			err.ErrorCode = code.FailedPrecondition
			err.Msg = "ULG User Error"
			foundation.HTTPResponse(w, "", err)
			return
		}

		iPratyAccount = &UserInfo
		result["userCoinQuota"] = UserInfo.UserCoinQuota
		result["gameInfo"] = UserInfo.GameInfo
	default:
		errorlog.ErrorLogPrintln("logintype Error", logintype, gametypeid, postData)
	}

	Info, err := db.GetAccountInfo(iPratyAccount.PartyAccount())
	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "LoginType", err)
		return
	}

	var accountInfo player.AccountInfo
	if len(Info) < 1 {
		accountInfo = player.NewAccountInfo(iPratyAccount.PartyAccount(), iPratyAccount.GameAccount(), iPratyAccount.AccountType())
		db.NewAccount(accountInfo.Account, accountInfo.GameAccount, accountInfo.AccountType)

	} else {
		result := Info[0]
		accountInfo = player.NewAccountInfo(result["Account"].(string), result["GameAccount"].(string), result["AccountType"].(int64))
		// db.UpdateAccount(accountInfo.Account, accountInfo.LoginTime)
	}

	player.SaveAccountInfo(&accountInfo)
	result["gameaccount"] = accountInfo.GameAccount
	result["token"] = accountInfo.Token
	result["serversetting"] = setting.New()

	loginfo := log.New(log.Login)
	loginfo.Account = accountInfo.GameAccount
	log.SaveLog(loginfo)
	foundation.HTTPResponse(w, result, err)
}
