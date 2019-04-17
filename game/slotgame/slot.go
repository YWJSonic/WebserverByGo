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

// HandleURL ...
// var HandleURL []frame.RESTfulUrl

// SlotRoomInfo ...
type SlotRoomInfo struct {
	RoomInfo *game.RoomInfo
}

// SimpleGameRoomInfo ...
type SimpleGameRoomInfo struct {
	Lock   bool
	Locker code.PlayerID
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

	gameRoomArray = createGameRoomInArray(5)

	// HandleURL = append(HandleURL, frame.RESTfulUrl{RequestType: "POST", URL: "slotgame/", Fun: runa})
	HandleURL = append(HandleURL, transmission.RESTfulURL{RequestType: "POST", URL: "slotgame/join", Fun: join})
	HandleURL = append(HandleURL, transmission.RESTfulURL{RequestType: "POST", URL: "slotgame/leave", Fun: leave})
	HandleURL = append(HandleURL, transmission.RESTfulURL{RequestType: "POST", URL: "slotgame/gameresult", Fun: gameresult})
	HandleURL = append(HandleURL, transmission.RESTfulURL{RequestType: "POST", URL: "slotgame/bet", Fun: bet})
	// HandleURL = append(HandleURL, transmission.RESTfulURL{RequestType: "POST", URL: "slotgame/GetGameRoomList", Fun: getGameRoomList})
	return HandleURL
}

func createGameRoomInArray(RoomLimit int) map[int]SlotRoomInfo {
	var tmpRoom = make(map[int]SlotRoomInfo)
	gameID := "slot"

	for i := int(0); i < RoomLimit; i++ {
		tmp := game.CreatedGameRoom(i, 1, gameID, gameID)
		tmpRoom[i] = SlotRoomInfo{&tmp}
		fmt.Println("Create", gameID, "game room", i)
	}
	return tmpRoom
}

// func createGameRoomInArray(RoomLimit int) []SlotRoomInfo {
// 	var tmpRoom []SlotRoomInfo
// 	fmt.Println(RoomLimit, len(tmpRoom))
// 	for i := 0; i < RoomLimit; i++ {
// 		tmpRoom = append(tmpRoom, SlotRoomInfo{RoomInfo: game.CreatedGameRoom(i, 1, "slot", "slot")})
// 	}
// 	return tmpRoom
// }
func join(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// mu.Lock()
	// defer mu.Unlock()

	playerid := code.PlayerID(1)
	player := foundation.GetPlayerInfo(playerid)

	resoult := make(map[string]interface{})
	// var roomStat int8
	var err errorlog.ErrorMsg
	for index := 0; index < 5; index++ {
		_, err = gameRoomArray[index].RoomInfo.Join(&player)
		if err.MsgNum != game.OK {
			continue
		}

		resoult["roomState"] = gameRoomArray[index].RoomInfo.Status()
		resoult["roomid"] = gameRoomArray[index].RoomInfo.ID()
		fmt.Printf("RoomID: %d Player Count: %d\n", gameRoomArray[index].RoomInfo.ID(), len(gameRoomArray[index].RoomInfo.Players()))
		break
	}

	foundation.HTTPResponse(w, resoult, err)
}
func leave(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// mu.Lock()
	// defer mu.Unlock()
	POSTData := foundation.PostData(r)
	playerID := foundation.InterfaceToInt(POSTData["playerid"])
	playerid := code.PlayerID(playerID)
	player := foundation.GetPlayerInfo(playerid)

	resoult := make(map[string]interface{})
	// var roomStat int8
	var err errorlog.ErrorMsg

	fmt.Println(gameRoomArray)
	_, err = gameRoomArray[player.InRoom].RoomInfo.Leave(&player)
	if err.MsgNum != game.OK {
		panic("Room leave error " + string(player.ID))
	}

	resoult["roomState"] = gameRoomArray[player.InRoom].RoomInfo.Status()
	resoult["roomid"] = gameRoomArray[player.InRoom].RoomInfo.ID()
	fmt.Printf("RoomID: %d Player Count: %d\n", gameRoomArray[player.InRoom].RoomInfo.ID(), len(gameRoomArray[player.InRoom].RoomInfo.Players()))

	foundation.HTTPResponse(w, resoult, err)
}

func bet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
func gameresult(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := errorlog.ErrorMsg{}
	// postData := foundation.PostData(r)
	// token := postData["playerid"].(string)
	// gametoken := postData["token"].(string)
	// gameid := postData["gameid"].(string)
	// bet := foundation.InterfaceToInt(postData["bet"])

	reault := gamelogic.GameOutput([]int{0, 1, 2, 3}, []int{3, 3, 3, 3, 3})

	foundation.HTTPResponse(w, reault, err)

}
