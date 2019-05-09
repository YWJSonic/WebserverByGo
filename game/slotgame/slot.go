package slotgame

import (
	"fmt"
	"net/http"
	"sync"

	"../../code"
	"../../db"
	"../../foundation"
	"../../game"
	"../../log"
	"../../player"

	"../../messagehandle/errorlog"
	"github.com/julienschmidt/httprouter"
)

// ErrorNum
const (
	SelfNoInRoom int8 = iota
)

// SlotRoomInfo ...
type SlotRoomInfo struct {
	BaseInfo *game.RoomInfo
}

// SimpleGameRoomInfo ...
type SimpleGameRoomInfo struct {
	Lock   bool
	Locker int64
}

var gameRoomArray map[int]SlotRoomInfo // gameroom with array
var isInit = false
var mu *sync.RWMutex

// RoomCount Current room count
var RoomCount = 0

// ServiceStart ...
func ServiceStart() []foundation.RESTfulURL {
	var HandleURL []foundation.RESTfulURL

	if isInit {
		return HandleURL
	}

	mu = new(sync.RWMutex)
	isInit = true

	gameRoomArray = addGameRoomInArray(2)

	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "slotgame/join", Fun: join})
	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "slotgame/leave", Fun: leave})
	HandleURL = append(HandleURL, foundation.RESTfulURL{RequestType: "POST", URL: "slotgame/gameresult", Fun: gameresult})
	return HandleURL
}

// Update ...
func Update() {
	CleanGameRoom()
}

// CleanGameRoom Remove dead player
func CleanGameRoom() {
	mu.Lock()
	defer mu.Unlock()

	for _, GameRoom := range gameRoomArray {
		GameRoom.BaseInfo.ClearRoom()
	}
}
func addGameRoomInArray(RoomLimit int) map[int]SlotRoomInfo {
	var tmpRoom = make(map[int]SlotRoomInfo)
	gametypeID := "slot"

	for i := 0; i < RoomLimit; i++ {
		RoomCount++
		tmp := game.CreatedGameRoom(RoomCount, 1, gametypeID, gametypeID)
		tmpRoom[RoomCount] = SlotRoomInfo{&tmp}
	}
	return tmpRoom
}

func join(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()

	result := make(map[string]interface{})
	POSTData := foundation.PostData(r)
	// token := foundation.InterfaceToString(POSTData["token"])
	playerID := foundation.InterfaceToInt64(POSTData["playerid"])
	playerInfo, err := player.GetPlayerInfoByPlayerID(playerID)

	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return
	}

	if playerInfo.IsInGameRoom() {
		err.ErrorCode = code.AlreadyInGame
		err.Msg = "AlreadyInGame"
		foundation.HTTPResponse(w, "", err)
		return
	}

	loginfo := log.New(log.JoinGame)
	loginfo.PlayerID = playerInfo.ID
	RoomCount := len(gameRoomArray)
	for index := 1; index <= RoomCount; index++ {
		_, err = gameRoomArray[index].BaseInfo.Join(playerInfo)
		if err.ErrorCode != code.OK {
			continue
		}

		result["roomState"] = gameRoomArray[index].BaseInfo.Status()
		result["roomid"] = gameRoomArray[index].BaseInfo.ID()
		loginfo.IValue1 = int64(gameRoomArray[index].BaseInfo.ID())
		break
	}

	if !playerInfo.IsInGameRoom() {
		err.ErrorCode = code.RoomFull
		err.Msg = "RoomFull"
		foundation.HTTPResponse(w, "", err)
		return
	}

	player.SavePlayerInfo(playerInfo)
	log.SaveLog(loginfo)
	foundation.HTTPResponse(w, result, err)
}
func leave(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()
	POSTData := foundation.PostData(r)
	// token := foundation.InterfaceToString(POSTData["token"])
	playerID := foundation.InterfaceToInt64(POSTData["playerid"])
	playerid := playerID
	playerInfo, err := player.GetPlayerInfoByPlayerID(playerid)
	loginfo := log.New(log.LeaveGame)
	loginfo.PlayerID = playerid

	result := make(map[string]interface{})
	roomInfo, err := getRoomInfo(playerInfo.InRoom)
	if roomInfo == nil {
		foundation.HTTPResponse(w, "", err)
		return
	}

	_, err = roomInfo.BaseInfo.Leave(playerInfo.ID)
	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return
	}

	result["roomState"] = roomInfo.BaseInfo.Status()
	result["roomid"] = roomInfo.BaseInfo.ID()

	loginfo.IValue1 = int64(roomInfo.BaseInfo.ID())
	player.SavePlayerInfo(playerInfo)
	log.SaveLog(loginfo)
	foundation.HTTPResponse(w, result, err)
}

func gameresult(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()

	var result = make(map[string]interface{})

	postData := foundation.PostData(r)
	playerid := foundation.InterfaceToInt64(postData["playerid"])
	playerinfo, err := player.GetPlayerInfoByPlayerID(playerid)
	// gametoken := postData["token"].(string)
	// gametypeid := postData["gametypeid"].(string)
	bet := foundation.InterfaceToInt(postData["bet"])

	if err.ErrorCode != code.OK {
		foundation.HTTPResponse(w, "", err)
		return
	}
	loginfo := log.New(log.GameResult)
	loginfo.PlayerID = playerid

	WinBet := 0
	ScatterIndex := -1
	NormalGameResult := 0
	var newIndex []int
	var newplate []int
	var GameResult [][]int

	for i := 0; i < 2; i++ {
		newIndex, newplate = newPlate([]int{1, 1, 1}, [][]int{Sroll1, Sroll2, Sroll3})
		GameResult = gameResult(newplate)

		if len(GameResult) > 0 {
			NormalGameResult = GameResult[0][3]
			if NormalGameResult > 0 {
				WinBet = NormalGameResult
			} else if NormalGameResult == -100 {
				ScatterIndex, WinBet = scatter1()
			} else if NormalGameResult == -101 {
				ScatterIndex, WinBet = scatter2()
			}
		}
		if (WinBet * bet) < 50000 {
			break
		}
	}
	// GameResult = append(GameResult, newplate)
	WinMoney := int64(WinBet * bet)
	LostMoney := int64(bet)
	playerinfo.TotalWin += WinMoney
	playerinfo.TotalLost += LostMoney
	playerinfo.Money = playerinfo.Money + WinMoney - LostMoney
	result["ScatterGame"] = ScatterIndex
	result["NormalGameIndex"] = newIndex
	result["NormalGame"] = newplate
	result["WinMoney"] = WinMoney
	// result["Player"] = playerInfo.ToJson()
	loginfo.IValue1 = int64(WinBet * bet)
	loginfo.SValue1 = fmt.Sprint(newplate)
	loginfo.SValue2 = fmt.Sprint(ScatterIndex)

	db.UpdatePlayerInfo(playerinfo.ID, playerinfo.Money, playerinfo.TotalWin, playerinfo.TotalLost)
	log.SaveLog(loginfo)
	foundation.HTTPResponse(w, result, err)
}
func addroom(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()

	NewRoomInfo := addGameRoomInArray(5)

	for roomid, roominfo := range NewRoomInfo {
		gameRoomArray[roomid] = roominfo
	}

	foundation.HTTPResponse(w, "OK", errorlog.New())
}
func getRoomInfo(roomID int) (*SlotRoomInfo, errorlog.ErrorMsg) {
	err := errorlog.New()

	if roomInfo, ok := gameRoomArray[roomID]; ok {
		return &roomInfo, err
	}

	err.ErrorCode = code.RoomNotExistence
	err.Msg = "RoomNotExistence"
	return nil, err
}
func getRoom(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := errorlog.New()
	var result map[string]interface{}
	result = make(map[string]interface{})
	for i, x := range gameRoomArray {
		result[fmt.Sprintf("%d", i)] = x.BaseInfo.JoinTime()
	}

	foundation.HTTPResponse(w, result, err)
}
