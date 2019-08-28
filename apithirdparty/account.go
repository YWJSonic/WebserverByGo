package apithirdparty

import (
	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/playerinfo"
	"gitlab.com/ServerUtility/thirdparty/ulginfo"
	"gitlab.com/WeberverByGoBase/apithirdparty/ulg"
	db "gitlab.com/WeberverByGoBase/handledb"
	"gitlab.com/WeberverByGoBase/player"
	"gitlab.com/WeberverByGoBase/serversetting"
)

// GetAccount ...
func GetAccount(logintype int, accounttoken, gametypeid string) (map[string]interface{}, *playerinfo.AccountInfo, messagehandle.ErrorMsg) {
	result := make(map[string]interface{})
	var iPratyAccount playerinfo.IPratyAccount

	switch logintype {
	case playerinfo.Ulg:
		UserInfo, err := ulg.GetUser(accounttoken, gametypeid)
		if err.ErrorCode != code.OK {
			err.ErrorCode = code.FailedPrecondition
			return result, nil, err
		}

		iPratyAccount = &UserInfo
		result["userCoinQuota"] = UserInfo.UserCoinQuota
		result["gameInfo"] = UserInfo.GameInfo
	default:
		messagehandle.ErrorLogPrintln("logintype Error", logintype, gametypeid, logintype, accounttoken)
		err := messagehandle.New()
		err.ErrorCode = code.LoginTypeError
		err.Msg = "logintype Error"
		return result, nil, err
	}

	Info, err := db.GetAccountInfo(iPratyAccount.PartyAccount())
	if err.ErrorCode != code.OK {
		return result, nil, err
	}

	var accountInfo playerinfo.AccountInfo
	if len(Info) < 1 {
		accountInfo = player.NewAccountInfo(iPratyAccount.PartyAccount(), iPratyAccount.GameAccount(serversetting.AccountEncodeStr), foundation.NewToken(iPratyAccount.PartyAccount()), iPratyAccount.AccountType())
		db.NewAccount(accountInfo.Account, accountInfo.GameAccount, accountInfo.AccountType)

	} else {
		result := Info[0]
		accountInfo = player.NewAccountInfo(result["Account"].(string), result["GameAccount"].(string), foundation.NewToken(result["Account"].(string)), result["AccountType"].(int64))
		// db.UpdateAccount(accountInfo.Account, accountInfo.LoginTime)
	}

	player.SaveAccountInfo(&accountInfo)
	return result, &accountInfo, messagehandle.New()
}

// Excahnge ...
func Excahnge(playerInfo *playerinfo.Info, accountToken, gametypeid string, cointype, coinamount int64) (*ulginfo.Result, messagehandle.ErrorMsg) {

	// new thirdparty token
	ulguser, err := ulg.Authorized(accountToken, gametypeid)
	if err.ErrorCode != code.OK {
		err.ErrorCode = code.FailedPrecondition
		return nil, err
	}

	// exchange
	messagehandle.LogPrintln("AccountUserInfo.GameToken", ulguser.GameToken)
	ulgResult, err := ulg.Exchange(ulguser.GameToken, gametypeid, accountToken, cointype, coinamount)
	if err.ErrorCode != code.OK {
		return nil, err
	}

	_, err = ulg.NewULGInfo(playerInfo.ID, cointype, coinamount, ulguser.GameToken, accountToken)
	if err.ErrorCode != code.OK {
		return nil, err
	}

	OldMoney := playerInfo.Money
	playerInfo.Money = ulgResult.GameCoin
	playerInfo.GameToken = ulguser.GameToken
	player.SavePlayerInfo(playerInfo)
	db.SetExchange(playerInfo.GameAccount, playerInfo.GameToken, cointype, coinamount, playerInfo.Money, OldMoney, foundation.ServerNowTime())
	return &ulgResult, messagehandle.New()
}

// CheckOut ...
func CheckOut(playerInfo *playerinfo.Info, gameTypeID string) (interface{}, messagehandle.ErrorMsg) {

	ulgInfo, err := ulg.GetULGInfo(playerInfo.ID, playerInfo.GameToken)
	if err.ErrorCode != code.OK {
		return nil, err
	}
	if ulgInfo.IsCheckOut {
		playerInfo.Money = 0
		playerInfo.GameToken = ""
		player.SavePlayerInfo(playerInfo)
		return nil, err
	}

	ulgCheckOutResult, err := ulg.Checkout(ulgInfo, gameTypeID)
	if err.ErrorCode != code.OK && err.ErrorCode != code.ExchangeError {
		return nil, err
	}

	ulg.UpdateUlgInfoCheckOut(playerInfo.GameToken)
	return ulgCheckOutResult.UserCoinQuota, messagehandle.New()
}

// Refresh ...
func Refresh(accountToken, gameTypeID string) (interface{}, messagehandle.ErrorMsg) {
	UserInfo, err := ulg.GetUser(accountToken, gameTypeID)
	if err.ErrorCode != code.OK {
		return nil, err
	}

	return UserInfo.UserCoinQuota, err
}
