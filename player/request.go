package player

import (
	"encoding/json"
	"time"

	"../code"
	"../db"
	"../foundation"
	"../messagehandle/errorlog"
	"../mycache"
)

// GetAccountInfoByGameAccount Get accountinfo struct
func GetAccountInfoByGameAccount(GameAccount string) (*AccountInfo, errorlog.ErrorMsg) {
	info, err := mycache.GetAccountInfo(GameAccount)
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
func New(GameAccount string) (*PlayerInfo, errorlog.ErrorMsg) {
	playerID, err := db.NewGameAccount(GameAccount, 0, "")

	if err.ErrorCode != code.OK {
		return nil, err
	}

	info := PlayerInfo{
		GameAccount: GameAccount,
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
