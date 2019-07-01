package login

import (
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/loginfo"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/myhttp"
	"gitlab.com/ServerUtility/playerinfo"
	"gitlab.com/WeberverByGo/db"
	"gitlab.com/WeberverByGo/foundation/myrestful"
	"gitlab.com/WeberverByGo/messagehandle/log"
	"gitlab.com/WeberverByGo/player"
	"gitlab.com/WeberverByGo/setting"
	"gitlab.com/WeberverByGo/thirdparty/guest"
	"gitlab.com/WeberverByGo/thirdparty/ulg"
)

var isInit = false
var mu *sync.RWMutex

// ServiceStart ...
func ServiceStart() []myhttp.RESTfulURL {
	var HandleURL []myhttp.RESTfulURL

	if isInit {
		return HandleURL
	}

	mu = new(sync.RWMutex)
	isInit = true

	HandleURL = append(HandleURL, myhttp.RESTfulURL{RequestType: "POST", URL: "account/login", Fun: login, ConnType: myhttp.Client})
	return HandleURL
}
func login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var result = make(map[string]interface{})

	postData := myhttp.PostData(r)
	logintype := foundation.InterfaceToInt(postData["logintype"])
	gametypeid := foundation.InterfaceToString(postData["gametypeid"])

	err := messagehandle.New()
	if err = foundation.CheckGameType(gametypeid); err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "", err)
		return
	}

	// mutility party switch
	var iPratyAccount playerinfo.IPratyAccount
	switch logintype {
	case playerinfo.Guest:
		account := foundation.InterfaceToString(postData["account"])
		guestinfo := guest.GuestInfo{
			Account: account,
		}

		iPratyAccount = &guestinfo
	case playerinfo.Self:
	case playerinfo.Ulg:
		accounttoken := foundation.InterfaceToString(postData["accounttoken"])
		UserInfo, err := ulg.GetUser(accounttoken, gametypeid)
		if err.ErrorCode != code.OK {
			err.ErrorCode = code.FailedPrecondition
			myrestful.HTTPResponse(w, "", err)
			return
		}

		iPratyAccount = &UserInfo
		result["userCoinQuota"] = UserInfo.UserCoinQuota
		result["gameInfo"] = UserInfo.GameInfo
	default:
		messagehandle.ErrorLogPrintln("logintype Error", logintype, gametypeid, postData)
	}

	Info, err := db.GetAccountInfo(iPratyAccount.PartyAccount())
	if err.ErrorCode != code.OK {
		myrestful.HTTPResponse(w, "LoginType", err)
		return
	}

	var accountInfo playerinfo.AccountInfo
	if len(Info) < 1 {
		accountInfo = playerinfo.NewAccountInfo(iPratyAccount.PartyAccount(), iPratyAccount.GameAccount(), foundation.NewToken(iPratyAccount.PartyAccount()), iPratyAccount.AccountType())
		db.NewAccount(accountInfo.Account, accountInfo.GameAccount, accountInfo.AccountType)

	} else {
		result := Info[0]
		accountInfo = playerinfo.NewAccountInfo(result["Account"].(string), result["GameAccount"].(string), foundation.NewToken(result["Account"].(string)), result["AccountType"].(int64))
		// db.UpdateAccount(accountInfo.Account, accountInfo.LoginTime)
	}

	player.SaveAccountInfo(&accountInfo)
	result["gameaccount"] = accountInfo.GameAccount
	result["token"] = accountInfo.Token
	result["serversetting"] = setting.New()

	loginfo := loginfo.New(loginfo.Login)
	loginfo.Account = accountInfo.GameAccount
	log.SaveLog(loginfo)
	myrestful.HTTPResponse(w, result, err)
}
