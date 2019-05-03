package player

import (
	"encoding/json"
	"fmt"
	"time"

	"../code"
	"../data"
	"../db"
	"../foundation"
	"../messagehandle/errorlog"
	"../mycache"
)

// GetPlayerInfoByPlayerID ...
func GetPlayerInfoByPlayerID(playerid int64) (*PlayerInfo, errorlog.ErrorMsg) {

	var player PlayerInfo
	// err := errorlog.New()
	info, err := mycache.GetPlayerInfo(playerid) // get playerinfo form cache
	// info, errMsg := (*data.CacheRef).Get(fmt.Sprintf("ID%dJS", playerid))

	// if errMsg != nil {
	// 	fmt.Println(errMsg)
	// 	panic(errMsg)
	// }
	// if err.ErrorCode != code.OK {
	// 	fmt.Println("GetPlayerInfoByPlayerID ERROR:", err.Msg, "info", info)
	// }
	if info == nil {
		err.ErrorCode = code.NoThisPlayer
		err.Msg = "NoThisPlayer"
		return nil, err
	}

	if errMsg := json.Unmarshal(info.([]byte), &player); errMsg != nil {
		panic(errMsg)
	}

	return &player, err
}

// SavePlayerInfo ...
func SavePlayerInfo(playerInfo *PlayerInfo) {
	playerInfo.LastCheckTime = time.Now().Unix()

	(*data.CacheRef).Set(fmt.Sprintf("ID%dJS", playerInfo.ID), playerInfo.ToJsonStr(), time.Hour)
	mycache.SetPlayerInfo(playerInfo.ID, playerInfo.ToJsonStr())
	mycache.SetPlayerID(playerInfo.GameAccount, playerInfo.ID)

}

// New Create a new PlayerInfo
func New(GameAccount string) *PlayerInfo {
	playerID, err := db.NewGameAccount(GameAccount, 0)

	if err.ErrorCode != code.OK {
		fmt.Println(err)
	}

	info := PlayerInfo{
		GameAccount: GameAccount,
		ID:          playerID,
		Money:       0,
	}

	return &info
}

// MakePlayer Get player form db
func MakePlayer(row map[string]interface{}) *PlayerInfo {
	return &PlayerInfo{
		ID:          foundation.InterfaceToInt64(row["PlayerID"]),
		Money:       foundation.InterfaceToInt64(row["GameMoney"]),
		GameAccount: foundation.InterfaceToString(row["GameAccount"]),
	}
}

// JSONMakePlayer create playerinfo
func JSONMakePlayer(jsbyte []byte) PlayerInfo {
	var info PlayerInfo
	if err := json.Unmarshal(jsbyte, &info); err != nil {
		panic(err)
	}
	return info
}
