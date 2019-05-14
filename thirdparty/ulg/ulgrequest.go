package ulg

import (
	"encoding/json"
	"fmt"

	"../../code"
	"../../db"
	"../../foundation"
	"../../messagehandle/errorlog"
	"../../mycache"
)

// LoginURL ...
func LoginURL() string {
	return loginURL
}

// GetUserURL ...
func GetUserURL() string {
	return getuserURL
}

// AuthorizedURL ...
func AuthorizedURL() string {
	return authorizedURL
}

// ExchangeURL ...
func ExchangeURL() string {
	return exchangeURL
}

// CheckoutURL ...
func CheckoutURL() string {
	return checkoutURL
}

/////// API interface process

// Login ...
// func Login(username, password string) map[string]interface{} {
// 	var info map[string]interface{}
// 	postData := map[string][]string{
// 		"username": {username},
// 		"password": {password},
// 	}
// 	jsbyte := foundation.HTTPPostRequest(loginURL, postData)
// 	if err := json.Unmarshal(jsbyte, &info); err != nil {
// 		panic(err)
// 	}
// 	return info
// }

// NewULGInfo New ULGInfo
func NewULGInfo(playerid, coinamount int64, gametoken string) ULGInfo {
	info := ULGInfo{
		PlayerID:  playerid,
		GameToken: gametoken,
	}
	db.NewULGInfoRow(playerid, gametoken)
	mycache.SetULGInfo(fmt.Sprintf("ULG%d", playerid), info.ToJSONStr())
	return info
}

// GetULGInfo ...
func GetULGInfo(gametoken string) (*ULGInfo, errorlog.ErrorMsg) {
	ULGJsStr, err := mycache.GetULGInfoCache(gametoken)
	var ulginfo *ULGInfo

	// cache no data
	if err.ErrorCode != code.OK {
		ulginfomap, err := db.GetULGInfoRow(gametoken)

		// db no data
		if err.ErrorCode != code.OK {
			errorlog.ErrorLogPrintln("Cache", err)
			return nil, err
		}
		if len(ulginfomap) < 1 {
			errorlog.ErrorLogPrintln("DB", err)
			err.ErrorCode = code.NoThisPartyInfo
			err.Msg = "NoThisPartyInfo"
			return nil, err
		}
		ulginfo = MakeULGInfo(ulginfomap[0])

	} else {
		if errMsg := json.Unmarshal([]byte(ULGJsStr), ulginfo); errMsg != nil {
			errorlog.ErrorLogPrintln("Cache", errMsg)
			err.ErrorCode = code.NoThisPartyInfo
			err.Msg = "ULGInfoFormatError"
			return nil, err
		}
	}

	return ulginfo, err
}

// UpdateULGInfo ...
func UpdateULGInfo(ulginfo *ULGInfo, BetMoney, WinBet int64) {
	WinMoney := WinBet * BetMoney

	ulginfo.TotalBet += BetMoney
	ulginfo.TotalWin += (BetMoney * WinBet)
	if BetMoney > WinMoney {

		ulginfo.TotalLost += (WinMoney - BetMoney)
	}
	SaveULGInfo(ulginfo)
}

// SaveULGInfo ...
func SaveULGInfo(Info *ULGInfo) {
	mycache.SetULGInfo(fmt.Sprintf("ULG%d", Info.PlayerID), Info.ToJSONStr())
	db.UpdateULGInfoRow(Info.GameToken, Info.TotalWin, Info.ToJSONStr, Info.IsCheckOut)
}

// MakeULGInfo get ulg info form db
func MakeULGInfo(row map[string]interface{}) *ULGInfo {
	return &ULGInfo{
		PlayerID:   foundation.InterfaceToInt64(row["PlayerID"]),
		GameToken:  foundation.InterfaceToString(row["GameToken"]),
		TotalWin:   foundation.InterfaceToInt64(row["TotalWin"]),
		TotalLost:  foundation.InterfaceToInt64(row["TotalLost"]),
		IsCheckOut: foundation.InterfaceToBool(row["CheckOut"]),
	}
}

// GetUser client request getplayer info
func GetUser(token, gameid string) (UlgResult, errorlog.ErrorMsg) {
	var info UlgResult
	err := errorlog.New()
	postData := map[string][]string{
		"token":   {token},
		"game_id": {gameid},
	}
	jsbyte := foundation.HTTPPostRequest(getuserURL, postData)
	if jserr := json.Unmarshal(jsbyte, &info); jserr != nil {
		err.ErrorCode = code.GetUserError
		err.Msg = "UserFormatError"
	}

	if info.Result == 1 {
		err.ErrorCode = code.OK
	} else {
		err.ErrorCode = code.GetUserError
		err.Msg = info.ErrorMsg
	}
	return info, err
}

// Authorized ...
func Authorized(token, gametypeid string) (UlgResult, errorlog.ErrorMsg) {
	var info UlgResult
	err := errorlog.New()
	postData := map[string][]string{
		"token":   {token},
		"game_id": {gametypeid},
	}
	jsbyte := foundation.HTTPPostRequest(authorizedURL, postData)
	if jserr := json.Unmarshal(jsbyte, &info); jserr != nil {
		err.ErrorCode = code.AuthorizedError
		err.Msg = "AuthorizedFormatError"
	}

	if info.Result == 1 {
		err.ErrorCode = code.OK
	} else {
		err.ErrorCode = code.AuthorizedError
		err.Msg = info.ErrorMsg
	}
	return info, err
}

// Exchange ...
func Exchange(gametoken, gametypeid, accounttoken string, cointype, coinamount int) (UlgResult, errorlog.ErrorMsg) { // map[string]interface{} {
	var info UlgResult
	err := errorlog.New()
	postData := map[string][]string{
		"game_token":  {gametoken},
		"game_id":     {gametypeid},
		"token":       {accounttoken},
		"coin_type":   {fmt.Sprint(cointype)},
		"coin_amount": {fmt.Sprint(coinamount)},
	}
	jsbyte := foundation.HTTPPostRequest(exchangeURL, postData)
	if jserr := json.Unmarshal(jsbyte, &info); jserr != nil {
		err.ErrorCode = code.ExchangeError
		err.Msg = "ExchangeFormatError"
	}

	if info.Result == 1 {
		err.ErrorCode = code.OK
	} else {
		err.ErrorCode = code.ExchangeError
		err.Msg = info.ErrorMsg
	}
	return info, err
}

// Checkout ...
func Checkout(accounttoken, gametoken, gameid, amount string, totalwin, totalost int64) (UlgResult, errorlog.ErrorMsg) {
	var info UlgResult
	err := errorlog.New()
	postData := map[string][]string{
		"game_token": {gametoken},
		"game_id":    {gameid},
		"token":      {accounttoken},
		"amount":     {amount},
		"win":        {fmt.Sprint(int(600))},
		"lost":       {fmt.Sprint(int(300))},
	}
	jsbyte := foundation.HTTPPostRequest(checkoutURL, postData)
	if jserr := json.Unmarshal(jsbyte, &info); jserr != nil {
		err.ErrorCode = code.ExchangeError
		err.Msg = "ExchangeError"
	}
	//{"data":{"result":0,"userID":0,"status":0,"accountName":"","errorMsg":"\u0008checkout - 無遊戲紀錄","userName":"","token":"","game_token":"","gameCoin":0,"userCoinQuota":null,"coinsetting":null,"gameInfo":null},"error":{"ErrorCode":20,"Msg":"ExchangeError"}}
	if info.Result == 1 {
		err.ErrorCode = code.OK
	} else {
		err.ErrorCode = code.ExchangeError
		err.Msg = info.ErrorMsg
	}
	return info, err
}
