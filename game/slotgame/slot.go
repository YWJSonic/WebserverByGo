package slotgame

import (
	"fmt"
	"net/http"
	"sync"

	"../../foundation"
	"../../frame/code"
	"../../frame/game"
	"../../frame/transmission"
	"../../gamelogic"

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

// ServiceStart ...
func ServiceStart() []transmission.RESTfulURL {
	var HandleURL []transmission.RESTfulURL

	if isInit {
		return HandleURL
	}

	mu = new(sync.RWMutex)
	isInit = true

	gameRoomArray = addGameRoomInArray(20)

	HandleURL = append(HandleURL, transmission.RESTfulURL{RequestType: "POST", URL: "slotgame/join", Fun: join})
	HandleURL = append(HandleURL, transmission.RESTfulURL{RequestType: "POST", URL: "slotgame/leave", Fun: leave})
	HandleURL = append(HandleURL, transmission.RESTfulURL{RequestType: "POST", URL: "slotgame/gameresult", Fun: gameresult})
	return HandleURL
}

func addGameRoomInArray(RoomLimit int) map[int]SlotRoomInfo {
	var tmpRoom = make(map[int]SlotRoomInfo)
	gameID := "slot"

	for i := 0; i < RoomLimit; i++ {
		tmp := game.CreatedGameRoom(i+1, 1, gameID, gameID)
		tmpRoom[i] = SlotRoomInfo{&tmp}
		fmt.Println("Create", gameID, "game room", tmp.ID())
	}
	return tmpRoom
}

func join(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()

	POSTData := foundation.PostData(r)
	playerID := foundation.InterfaceToInt(POSTData["playerid"])
	playerid := playerID
	player := foundation.GetPlayerInfo(playerid)
	resoult := make(map[string]interface{})
	// var roomStat int8
	var err errorlog.ErrorMsg

	if !player.IsInGameRoom() {
		err.MsgNum = 4
		err.Msg = "AlreadyInGame"
		foundation.HTTPResponse(w, resoult, err)
	}

	RoomCount := len(gameRoomArray)
	for index := 0; index < RoomCount; index++ {
		_, err = gameRoomArray[index].BaseInfo.Join(player)
		if err.MsgNum != game.OK {
			continue
		}

		resoult["roomState"] = gameRoomArray[index].BaseInfo.Status()
		resoult["roomid"] = gameRoomArray[index].BaseInfo.ID()
		break
	}

	if !player.IsInGameRoom() {
		foundation.HTTPResponse(w, resoult, err)
	}

	foundation.SavePlayerInfo(player)
	foundation.HTTPResponse(w, resoult, err)
}
func leave(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()
	POSTData := foundation.PostData(r)
	playerID := foundation.InterfaceToInt(POSTData["playerid"])
	playerid := playerID
	player := foundation.GetPlayerInfo(playerid)

	resoult := make(map[string]interface{})
	roomInfo, err := getRoomInfo(player.InRoom)
	if roomInfo == nil {
		foundation.HTTPResponse(w, resoult, err)
		return
	}

	_, err = roomInfo.BaseInfo.Leave(player)

	if err.MsgNum != game.OK {
		foundation.HTTPResponse(w, resoult, err)
		return
	}

	resoult["roomState"] = roomInfo.BaseInfo.Status()
	resoult["roomid"] = roomInfo.BaseInfo.ID()

	foundation.SavePlayerInfo(player)
	foundation.HTTPResponse(w, resoult, err)
}

func gameresult(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// cleint request data
	// betmoney int
	// gameid string
	// playerid int

	err := errorlog.New()
	// postData := foundation.PostData(r)
	// token := postData["playerid"].(string)
	// gametoken := postData["token"].(string)
	// gameid := postData["gameid"].(string)
	// bet := foundation.InterfaceToInt(postData["bet"])

	reault := gamelogic.GetGameResult("slot", 200)

	foundation.HTTPResponse(w, reault, err)

}

func getRoomInfo(roomID int) (*SlotRoomInfo, errorlog.ErrorMsg) {
	err := errorlog.New()

	if roomInfo, ok := gameRoomArray[roomID-1]; ok {
		return &roomInfo, err
	}

	err.ErrorCode = code.RoomNotExistence
	err.Msg = "RoomNotExistence"
	return nil, err
}
