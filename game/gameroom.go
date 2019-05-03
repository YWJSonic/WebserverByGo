package game

import (
	"sync"
	"time"

	"../code"
	"../messagehandle/errorlog"
	"../player"
)

// room status
const (
	Empty int8 = iota
	Someone
	RoomFull
	RoomLock
)

// IRoom RoomInfo interface
type IRoom interface {
	ID() int
	GameName() string
	GameRound() int32
	IsRoomLock() bool
	RoomLocker() int64
	Players() []int64
	Join(player *player.PlayerInfo) (int8, errorlog.ErrorMsg)
	Leave(player *player.PlayerInfo) (int8, errorlog.ErrorMsg)
	Lock(player *player.PlayerInfo)
}

// RoomInfo base room info
type RoomInfo struct {
	MaxPlayer int32

	id        int
	gametype  string
	gamename  string
	locker    int64
	status    int8 // status 0:empty, 1:haspeople ,2:full, 3:lock
	gameround int32
	players   []int64
	jointime  []int64
	mu        *sync.RWMutex
}

// CreatedGameRoom  create room at init
func CreatedGameRoom(roomid int, maxPlayer int32, gameName, gametype string) RoomInfo {
	return RoomInfo{
		id:        roomid,
		MaxPlayer: maxPlayer,
		gametype:  gametype,
		gamename:  gameName,
		status:    0,
		gameround: 0,
		mu:        new(sync.RWMutex),
	}
}

// ID room ID
func (g RoomInfo) ID() int {
	return g.id
}

// GameName ...
func (g RoomInfo) GameName() string {
	return g.gamename
}

// Status status 0:empty, 1:full, 2:onlock
func (g RoomInfo) Status() int8 {
	return g.status
}

// GameRound ...
func (g RoomInfo) GameRound() int32 {
	return g.gameround
}

// IsRoomLock ...
func (g RoomInfo) IsRoomLock() bool {
	return g.locker != 0
}

// RoomLocker ...
func (g RoomInfo) RoomLocker() int64 {
	return g.locker
}

func (g RoomInfo) JoinTime() []int64 {
	return g.jointime
}

// Players ...
func (g *RoomInfo) Players() []int64 {
	return g.players
}

// Join ...
func (g *RoomInfo) Join(player *player.PlayerInfo) (int8, errorlog.ErrorMsg) {
	err := errorlog.New()

	// maybe meet self
	if g.IsRoomLock() && g.locker != player.ID {
		err.ErrorCode = code.RoomLock
		err.Msg = "RoomLockError"
		return g.status, err
	} else if g.isPlayerInRoom(player.ID) {
		err.ErrorCode = code.SelfInRoom
		err.Msg = "SelfInRoom"
		return g.status, err
	} else if g.status == RoomFull {
		err.ErrorCode = code.RoomFull
		err.Msg = "RoomFullError"
		return g.status, err
	}

	g.players = append(g.players, player.ID)
	g.jointime = append(g.jointime, time.Now().Unix())
	g.status = Someone
	player.InGame = g.gametype
	player.InRoom = g.id

	if int32(len(g.Players())) >= g.MaxPlayer {
		g.status = RoomFull
	}
	return g.status, err
}

// Leave ...
func (g *RoomInfo) Leave(player int64) (int8, errorlog.ErrorMsg) {
	err := errorlog.New()

	if !g.isPlayerInRoom(player) {
		err.ErrorCode = code.SelfNotInRoom
		err.Msg = "SelfNotInRoom"

	} else {
		count := len(g.players)
		for index := 0; index < count; index++ {
			if g.players[index] == player {
				g.players = append(g.players[:index], g.players[index+1:]...)
				g.jointime = append(g.jointime[:index], g.jointime[index+1:]...)
				break
			}
		}

		if int32(len(g.Players())) <= 0 {
			g.status = 0
		} else {
			g.status = 1
		}
	}

	return g.status, err
}

// Lock ...
func (g *RoomInfo) Lock(player *player.PlayerInfo) {
	g.locker = player.ID
}

func (g *RoomInfo) releace(player *player.PlayerInfo) {
	g.locker = -1
}
func (g *RoomInfo) isPlayerInRoom(playerID int64) bool {
	players := g.players
	for _, comparetTarget := range players {
		if comparetTarget == playerID {
			return true
		}
	}
	return false
}

// ClearRoom remove dead player
func (g *RoomInfo) ClearRoom() {
	if len(g.jointime) <= 0 {
		return
	}

	for _, playerID := range g.players {
		player, err := player.GetPlayerInfoByPlayerID(playerID)
		if err.ErrorCode != code.OK {
			g.Leave(playerID)
		}

		if player.IsPlayerConnect() {
			continue
		}

		g.Leave(playerID)
	}
}
