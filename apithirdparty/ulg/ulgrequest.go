package ulg

import (
	"encoding/json"
	"fmt"

	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/thirdparty/ulginfo"
	"gitlab.com/WeberverByGoBase/foundation/myrestful"
	mycache "gitlab.com/WeberverByGoBase/handlecache"
	db "gitlab.com/WeberverByGoBase/handledb"
)

/////// API interface process

// Login ...
// func Login(username, password string) map[string]interface{} {
// 	var info map[string]interface{}
// 	postData := map[string][]string{
// 		"username": {username},
// 		"password": {password},
// 	}
// 	jsbyte := myrestful.PostRawRequest(loginURL, postData)
// 	if err := json.Unmarshal(jsbyte, &info); err != nil {
// 		panic(err)
// 	}
// 	return info
// }

// NewULGInfo New ULGInfo
func NewULGInfo(playerid, cointype, exchangAmount int64, gameToken, accountToken string) (*ulginfo.Info, messagehandle.ErrorMsg) {
	info := ulginfo.Info{
		PlayerID:       playerid,
		GameToken:      gameToken,
		AccountToken:   accountToken,
		ExchangeType:   cointype,
		ExchangeAmount: exchangAmount,
	}
	err := db.NewULGInfoRow(playerid, gameToken, accountToken, exchangAmount, cointype)
	if err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("NewULGInfo-1", err)
		return nil, err
	}
	mycache.SetULGInfo(info.PlayerID, foundation.JSONToString(info))
	return &info, err
}

// GetULGInfo ...
func GetULGInfo(playerid int64, gametoken string) (*ulginfo.Info, messagehandle.ErrorMsg) {
	var ulginfo *ulginfo.Info
	var err messagehandle.ErrorMsg

	ULGJsStr := mycache.GetULGInfoCache(playerid)
	// cache no data
	if ULGJsStr == nil {
		var ulginfomap []map[string]interface{}
		ulginfomap, err = db.GetULGInfoRow(gametoken)

		// db no data
		if err.ErrorCode != code.OK {
			messagehandle.ErrorLogPrintln("GetULGInfo-1", err)
			return nil, err
		}
		if len(ulginfomap) < 1 {
			err.ErrorCode = code.NoExchange
			err.Msg = "NoExchange"
			messagehandle.ErrorLogPrintln("GetULGInfo-2", err, playerid, gametoken)
			return nil, err
		}
		ulginfo = MakeULGInfo(ulginfomap[0])

	} else {
		if errMsg := json.Unmarshal(ULGJsStr.([]byte), &ulginfo); errMsg != nil {
			err.ErrorCode = code.ULGInfoFormatError
			err.Msg = "ULGInfoFormatError"
			messagehandle.ErrorLogPrintln("GetULGInfo-3", errMsg, playerid, gametoken)
			return nil, err
		}
	}

	return ulginfo, err
}

// MaintainULGInfos ...
func MaintainULGInfos() ([]ulginfo.Info, messagehandle.ErrorMsg) {
	var Infos []ulginfo.Info
	result, err := db.ULGMaintainCheckoutRow()

	// db no data
	if err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("MaintainULGInfos-1", err)
		return nil, err
	}

	Infos = make([]ulginfo.Info, len(result))
	for i, row := range result {
		Infos[i] = *MakeULGInfo(row)
	}

	return Infos, err
}

// UpdateULGInfo ...
func UpdateULGInfo(ulginfo *ulginfo.Info, BetMoney, WinBet int64) {
	WinMoney := WinBet * BetMoney

	ulginfo.TotalBet += BetMoney
	ulginfo.TotalWin += (BetMoney * WinBet)
	if BetMoney > WinMoney {
		ulginfo.TotalLost += (WinMoney - BetMoney)
	}
	SaveULGInfo(ulginfo)
}

// UpdateUlgInfoCheckOut ...
func UpdateUlgInfoCheckOut(gametoken string) messagehandle.ErrorMsg {
	err := db.UpdateCheckUlgRow(gametoken)
	return err
}

// SaveULGInfo ...
func SaveULGInfo(info *ulginfo.Info) {
	mycache.SetULGInfo(info.PlayerID, foundation.JSONToString(info))
	db.UpdateULGInfoRow(info.GameToken, info.TotalBet, info.TotalWin, info.TotalLost, info.IsCheckOut)
}

// MakeULGInfo get ulg info form db
func MakeULGInfo(row map[string]interface{}) *ulginfo.Info {
	info := &ulginfo.Info{
		PlayerID:       foundation.InterfaceToInt64(row["PlayerID"]),
		GameToken:      foundation.InterfaceToString(row["GameToken"]),
		AccountToken:   foundation.InterfaceToString(row["AccountToken"]),
		ExchangeType:   foundation.InterfaceToInt64(row["ExchangeType"]),
		ExchangeAmount: foundation.InterfaceToInt64(row["ExchangeAmount"]),
		TotalBet:       foundation.InterfaceToInt64(row["TotalBet"]),
		TotalWin:       foundation.InterfaceToInt64(row["TotalWin"]),
		TotalLost:      foundation.InterfaceToInt64(row["TotalLost"]),
		IsCheckOut:     foundation.InterfaceToBool(row["CheckOut"]),
	}

	if _, ok := row["AccountToken"]; ok {
		info.AccountToken = foundation.InterfaceToString(row["AccountToken"])
	}

	return info
}

// GetUser client request getplayer info
func GetUser(token, gameid string) (ulginfo.Result, messagehandle.ErrorMsg) {
	var info ulginfo.Result
	err := messagehandle.New()
	postData := map[string]string{
		"token":   token,
		"game_id": gameid,
	}
	messagehandle.LogPrintln("Ulg", postData)
	jsbyte := myrestful.PostRawRequest(ulginfo.GetuserURL, foundation.ToJSONStr(postData))
	if jserr := json.Unmarshal(jsbyte, &info); jserr != nil {
		err.ErrorCode = code.GetUserError
		err.Msg = "UserFormatError"
		messagehandle.ErrorLogPrintln("GetUser-1", token, gameid, err)
	}

	if info.Result == 1 {
		err.ErrorCode = code.OK
	} else {
		err.ErrorCode = code.GetUserError
		err.Msg = info.ErrorMsg
		messagehandle.ErrorLogPrintln("GetUser-2", token, gameid, err)
	}
	return info, err
}

// Authorized ...
func Authorized(token, gametypeid string) (ulginfo.Result, messagehandle.ErrorMsg) {
	var info ulginfo.Result
	err := messagehandle.New()
	postData := map[string]string{
		"token":   token,
		"game_id": gametypeid,
	}
	messagehandle.LogPrintln("Ulg", postData)
	jsbyte := myrestful.PostRawRequest(ulginfo.AuthorizedURL, foundation.ToJSONStr(postData))
	if jserr := json.Unmarshal(jsbyte, &info); jserr != nil {
		err.ErrorCode = code.AuthorizedError
		err.Msg = "AuthorizedFormatError"
		messagehandle.ErrorLogPrintln("Authorized-1", token, gametypeid, err)
	}

	if info.Result == 1 {
		err.ErrorCode = code.OK
	} else {
		err.ErrorCode = code.AuthorizedError
		err.Msg = info.ErrorMsg
		messagehandle.ErrorLogPrintln("Authorized-2", token, gametypeid, err)
	}
	return info, err
}

// Exchange ...
func Exchange(gametoken, gametypeid, accounttoken string, cointype, coinamount int64) (ulginfo.Result, messagehandle.ErrorMsg) { // map[string]interface{} {
	var info ulginfo.Result
	err := messagehandle.New()
	postData := map[string]string{
		"game_token":  gametoken,
		"game_id":     gametypeid,
		"token":       accounttoken,
		"coin_type":   fmt.Sprint(cointype),
		"coin_amount": fmt.Sprint(coinamount),
	}
	messagehandle.LogPrintln("Ulg", postData)
	jsbyte := myrestful.PostRawRequest(ulginfo.ExchangeURL, foundation.ToJSONStr(postData))
	if jserr := json.Unmarshal(jsbyte, &info); jserr != nil {
		err.ErrorCode = code.ExchangeError
		err.Msg = "ExchangeFormatError"
		messagehandle.ErrorLogPrintln("Exchange-1", err, gametoken, gametypeid, accounttoken, cointype, coinamount)
	}

	if info.Result == 1 {
		err.ErrorCode = code.OK
	} else {
		err.ErrorCode = code.ExchangeError
		err.Msg = info.ErrorMsg
		messagehandle.ErrorLogPrintln("Exchange-2", err, gametoken, gametypeid, accounttoken, cointype, coinamount)
	}
	return info, err
}

// Checkout ...
func Checkout(ulgInfo *ulginfo.Info, gameid string) (ulginfo.CheckOutResult, messagehandle.ErrorMsg) { //accounttoken, gametoken, gameid, amount, totalwin, totalost string) (UlgCheckOutResult, messagehandle.ErrorMsg) {
	var info ulginfo.CheckOutResult
	err := messagehandle.New()
	postData := map[string]string{
		"game_token": ulgInfo.GameToken,
		"game_id":    gameid,
		"token":      ulgInfo.AccountToken,
		"amount":     fmt.Sprint(ulgInfo.TotalBet),
		"win":        fmt.Sprint(ulgInfo.TotalWin),
		"lost":       fmt.Sprint(ulgInfo.TotalLost),
	}
	jsbyte := myrestful.PostRawRequest(ulginfo.CheckoutURL, foundation.ToJSONStr(postData))
	if jserr := json.Unmarshal(jsbyte, &info); jserr != nil {
		err.ErrorCode = code.CheckoutError
		err.Msg = "CheckoutError"
		messagehandle.ErrorLogPrintln("Checkout-1", jserr, ulgInfo, gameid)
	}

	if info.Result == 1 {
		err.ErrorCode = code.OK
		switch ulgInfo.ExchangeType {
		case 1:
			info.UserCoinQuota.Coin1Out = info.UserCoinQuota.Coin1Out - ulgInfo.ExchangeAmount
		case 2:
			info.UserCoinQuota.Coin2Out = info.UserCoinQuota.Coin2Out - ulgInfo.ExchangeAmount
		case 3:
			info.UserCoinQuota.Coin3Out = info.UserCoinQuota.Coin3Out - ulgInfo.ExchangeAmount
		case 4:
			info.UserCoinQuota.Coin4Out = info.UserCoinQuota.Coin4Out - ulgInfo.ExchangeAmount
		}
	} else {
		err.ErrorCode = code.ExchangeError
		err.Msg = info.ErrorMsg
		messagehandle.ErrorLogPrintln("Checkout-2", err, ulgInfo, gameid)
	}
	return info, err
}
