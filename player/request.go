package player

import (
	"encoding/json"
	"time"

	"gitlab.com/WeberverByGo/code"
	"gitlab.com/WeberverByGo/db"
	"gitlab.com/WeberverByGo/foundation"
	"gitlab.com/WeberverByGo/messagehandle/errorlog"
	"gitlab.com/WeberverByGo/mycache"
)

// GetAccountInfoByGameAccount Get accountinfo struct
func GetAccountInfoByGameAccount(gameAccount string) (*AccountInfo, errorlog.ErrorMsg) {
	info, err := mycache.GetAccountInfo(gameAccount)
	if info == nil {
		err.ErrorCode = code.NoThisPlayer
		err.Msg = "NoThisGameAccount"
		return nil, err
	}

	var account AccountInfo
	if errMsg := json.Unmarshal(info.([]byte), &account); errMsg != nil {
		errorlog.ErrorLogPrintln("Player", errMsg)
		err.ErrorCode = code.NoThisPlayer
		err.Msg = "AccountFormatError"
		return nil, err
	}
	return &account, err

}

// SaveAccountInfo ...
func SaveAccountInfo(accInfo *AccountInfo) {
	mycache.SetAccountInfo(accInfo.GameAccount, accInfo.ToJSONStr())
	mycache.SetToken(accInfo.GameAccount, accInfo.Token)
}

// GetPlayerInfoByPlayerID ...
func GetPlayerInfoByPlayerID(playerid int64) (*PlayerInfo, errorlog.ErrorMsg) {
	info, err := mycache.GetPlayerInfo(playerid)
	if info == nil {
		err.ErrorCode = code.NoThisPlayer
		err.Msg = "NoThisPlayer"
		return nil, err
	}

	var player PlayerInfo
	if errMsg := json.Unmarshal(info.([]byte), &player); errMsg != nil {
		errorlog.ErrorLogPrintln("Player", errMsg)
		err.ErrorCode = code.NoThisPlayer
		err.Msg = "PlayerFormatError"
		return nil, err
	}

	return &player, err
}

// SavePlayerInfo ...
func SavePlayerInfo(playerInfo *PlayerInfo) {
	playerInfo.LastCheckTime = time.Now().Unix()

	// (*data.CacheRef).Set(fmt.Sprintf("ID%dJS", playerInfo.ID), playerInfo.ToJSONStr(), time.Hour)
	mycache.SetPlayerInfo(playerInfo.ID, playerInfo.ToJSONStr())
	db.UpdatePlayerInfo(playerInfo.ID, playerInfo.Money, playerInfo.GameToken)

}

// New Create a new PlayerInfo
func New(gameAccount string) (*PlayerInfo, errorlog.ErrorMsg) {
	playerID, err := db.NewGameAccount(gameAccount, 0, "")

	if err.ErrorCode != code.OK {
		return nil, err
	}

	info := PlayerInfo{
		GameAccount: gameAccount,
		ID:          playerID,
		Money:       0,
	}

	return &info, err
}

// MakePlayer Get player form db
func MakePlayer(row map[string]interface{}) *PlayerInfo {
	return &PlayerInfo{
		ID:          foundation.InterfaceToInt64(row["PlayerID"]),
		Money:       foundation.InterfaceToInt64(row["GameMoney"]),
		GameAccount: foundation.InterfaceToString(row["GameAccount"]),
		GameToken:   foundation.InterfaceToString(row["GameToken"]),
	}
}
