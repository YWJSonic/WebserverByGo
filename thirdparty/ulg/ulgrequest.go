package ulg

import (
	"encoding/json"
	"fmt"

	"gitlab.com/WeberverByGo/code"
	"gitlab.com/WeberverByGo/db"
	"gitlab.com/WeberverByGo/foundation"
	"gitlab.com/WeberverByGo/messagehandle/errorlog"
	"gitlab.com/WeberverByGo/mycache"
)

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
func NewULGInfo(playerid int64, gametoken, accounttoken string) (*ULGInfo, errorlog.ErrorMsg) {
	info := ULGInfo{
		PlayerID:     playerid,
		GameToken:    gametoken,
		AccountToken: accounttoken,
	}
	err := db.NewULGInfoRow(playerid, gametoken, accounttoken)
	if err.ErrorCode != code.OK {
		return nil, err
	}
	mycache.SetULGInfo(playerid, info.ToJSONStr())
	return &info, err
}

// GetULGInfo ...
func GetULGInfo(playerid int64, gametoken string) (*ULGInfo, errorlog.ErrorMsg) {
	var ulginfo *ULGInfo
	var err errorlog.ErrorMsg

	ULGJsStr := mycache.GetULGInfoCache(playerid)
	// cache no data
	if ULGJsStr == nil {
		var ulginfomap []map[string]interface{}
		ulginfomap, err = db.GetULGInfoRow(gametoken)

		// db no data
		if err.ErrorCode != code.OK {
			errorlog.ErrorLogPrintln("DB GetULGInfo", err)
			return nil, err
		}
		if len(ulginfomap) < 1 {
			errorlog.ErrorLogPrintln("DB GetULGInfo", err)
			err.ErrorCode = code.NoExchange
			err.Msg = "NoExchange"
			return nil, err
		}
		ulginfo = MakeULGInfo(ulginfomap[0])

	} else {
		if errMsg := json.Unmarshal(ULGJsStr.([]byte), &ulginfo); errMsg != nil {
			errorlog.ErrorLogPrintln("Cache ULGInfoFormatError", errMsg)
			err.ErrorCode = code.ULGInfoFormatError
			err.Msg = "ULGInfoFormatError"
			return nil, err
		}
	}

	return ulginfo, err
}

// MaintainULGInfos ...
func MaintainULGInfos() ([]ULGInfo, errorlog.ErrorMsg) {
	var Infos []ULGInfo
	result, err := db.ULGMaintainCheckoutRow()

	// db no data
	if err.ErrorCode != code.OK {
		errorlog.ErrorLogPrintln("Cache GetULGInfo", err)
		return nil, err
	}

	Infos = make([]ULGInfo, len(result))
	for i, row := range result {

		Infos[i] = *MakeULGInfo(row)
	}

	return Infos, err
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

// UpdateUlgInfoCheckOut ...
func UpdateUlgInfoCheckOut(gametoken string) errorlog.ErrorMsg {
	err := db.UpdateCheckUlgRow(gametoken)
	return err
}

// SaveULGInfo ...
func SaveULGInfo(Info *ULGInfo) {
	mycache.SetULGInfo(Info.PlayerID, Info.ToJSONStr())
	db.UpdateULGInfoRow(Info.GameToken, Info.TotalBet, Info.TotalWin, Info.TotalLost, Info.IsCheckOut)
}

// MakeULGInfo get ulg info form db
func MakeULGInfo(row map[string]interface{}) *ULGInfo {
	info := &ULGInfo{
		PlayerID:   foundation.InterfaceToInt64(row["PlayerID"]),
		GameToken:  foundation.InterfaceToString(row["GameToken"]),
		TotalBet:   foundation.InterfaceToInt64(row["TotalBet"]),
		TotalWin:   foundation.InterfaceToInt64(row["TotalWin"]),
		TotalLost:  foundation.InterfaceToInt64(row["TotalLost"]),
		IsCheckOut: foundation.InterfaceToBool(row["CheckOut"]),
	}

	if _, ok := row["AccountToken"]; ok {
		info.AccountToken = foundation.InterfaceToString(row["AccountToken"])
	}

	return info
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
func Checkout(accounttoken, gametoken, gameid, amount, totalwin, totalost string) (UlgCheckOutResult, errorlog.ErrorMsg) {
	var info UlgCheckOutResult
	err := errorlog.New()
	postData := map[string][]string{
		"game_token": {gametoken},
		"game_id":    {gameid},
		"token":      {accounttoken},
		"amount":     {amount},
		"win":        {totalwin},
		"lost":       {totalost},
	}
	jsbyte := foundation.HTTPPostRequest(checkoutURL, postData)
	if jserr := json.Unmarshal(jsbyte, &info); jserr != nil {
		err.ErrorCode = code.CheckoutError
		err.Msg = "CheckoutError"
	}

	if info.Result == 1 {
		err.ErrorCode = code.OK
	} else {
		err.ErrorCode = code.ExchangeError
		err.Msg = info.ErrorMsg
	}
	return info, err
}
