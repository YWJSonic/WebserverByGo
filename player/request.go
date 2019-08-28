package player

import (
	"encoding/json"
	"time"

	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/playerinfo"
	mycache "gitlab.com/WeberverByGoGame7/handlecache"
	db "gitlab.com/WeberverByGoGame7/handledb"
)

// CachePlayer memory cache player
var CachePlayer map[int64]playerinfo.Info

// GetAccountInfoByGameAccount Get accountinfo struct
func GetAccountInfoByGameAccount(gameAccount, token string) (*playerinfo.Info, messagehandle.ErrorMsg) {

	// check token
	if err := foundation.CheckToken(mycache.GetToken(gameAccount), token); err.ErrorCode != code.OK {
		return nil, err
	}

	var dbresult []map[string]interface{}
	dbresult, err := db.GetPlayerInfoByGameAccount(gameAccount)
	if err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("Lobby", err.Msg)
		return nil, err
	}

	var playerInfo *playerinfo.Info
	if len(dbresult) <= 0 {
		playerInfo, err = New(gameAccount)
	} else {
		playerInfo = MakePlayer(dbresult[0])
	}

	if playerInfo == nil {
		messagehandle.ErrorLogPrintln("Lobby", err.Msg)
		return nil, err
	}
	SavePlayerInfo(playerInfo)
	return playerInfo, err

}

// SaveAccountInfo ...
func SaveAccountInfo(accInfo *playerinfo.AccountInfo) {
	mycache.SetAccountInfo(accInfo.GameAccount, foundation.JSONToString(accInfo))
	mycache.SetToken(accInfo.GameAccount, accInfo.Token)
}

// GetPlayerInfoByGameAccount Get PlayerInfo by gameAccount with new token
func GetPlayerInfoByGameAccount(gameAccount string) (*playerinfo.Info, messagehandle.ErrorMsg) {
	dbresult, err := db.GetPlayerInfoByGameAccount(gameAccount)
	if err.ErrorCode != code.OK {
		messagehandle.ErrorLogPrintln("Lobby", err.Msg)
		return nil, err
	}

	var playerInfo *playerinfo.Info
	if len(dbresult) <= 0 {
		playerInfo, err = New(gameAccount)
	} else {
		playerInfo = MakePlayer(dbresult[0])
	}

	if playerInfo == nil {
		messagehandle.ErrorLogPrintln("Lobby", err.Msg)
		return nil, err
	}

	SavePlayerInfo(playerInfo)
	return playerInfo, messagehandle.New()
}

// GetPlayerInfoByPlayerID ...
func GetPlayerInfoByPlayerID(playerid int64) (*playerinfo.Info, messagehandle.ErrorMsg) {
	info, err := mycache.GetPlayerInfo(playerid)
	if info == nil {
		err.ErrorCode = code.NoThisPlayer
		err.Msg = "NoThisPlayer"
		return nil, err
	}

	var player playerinfo.Info
	if errMsg := json.Unmarshal(info.([]byte), &player); errMsg != nil {
		messagehandle.ErrorLogPrintln("Player", errMsg)
		err.ErrorCode = code.NoThisPlayer
		err.Msg = "PlayerFormatError"
		return nil, err
	}

	return &player, err
}

// SavePlayerInfo ...
func SavePlayerInfo(playerInfo *playerinfo.Info) {
	playerInfo.LastCheckTime = time.Now().Unix()

	mycache.SetPlayerInfo(playerInfo.ID, foundation.JSONToString(playerInfo))
	db.UpdatePlayerInfo(playerInfo.ID, playerInfo.Money, playerInfo.GameToken)

}

// New Create a new PlayerInfo
func New(gameAccount string) (*playerinfo.Info, messagehandle.ErrorMsg) {
	playerID, err := db.NewGameAccount(gameAccount, 0, "")

	if err.ErrorCode != code.OK {
		return nil, err
	}

	info := playerinfo.Info{
		GameAccount: gameAccount,
		ID:          playerID,
		Money:       0,
	}

	return &info, err
}

// MakePlayer Get player form db
func MakePlayer(row map[string]interface{}) *playerinfo.Info {
	return &playerinfo.Info{
		ID:          foundation.InterfaceToInt64(row["PlayerID"]),
		Money:       foundation.InterfaceToInt64(row["GameMoney"]),
		GameAccount: foundation.InterfaceToString(row["GameAccount"]),
		GameToken:   foundation.InterfaceToString(row["GameToken"]),
	}
}

// NewAccountInfo ...
func NewAccountInfo(account, gameAccount, token string, accountType int64) playerinfo.AccountInfo {
	return playerinfo.AccountInfo{
		Account:     account,
		GameAccount: gameAccount,
		Token:       token,
		LoginTime:   time.Now().Unix(),
		AccountType: accountType,
	}
}
