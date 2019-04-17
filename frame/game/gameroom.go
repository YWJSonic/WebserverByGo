package game

import (
	"../../messagehandle/errorlog"
	"../code"
	"../player"
)

// error status
const (
	OK int8 = iota
	RoomFullError
	RoomLockError
	SelfInRoom
	SelfNotInRoom
)

// room status
const (
	Empty int8 = iota
	Someone
	RoomFull
	RoomLock
)

// RoomInfo base room info
type RoomInfo struct {
	MaxPlayer int32

	id        int
	gametype  string
	gamename  string
	locker    code.PlayerID
	status    int8 // status 0:empty, 1:haspeople ,2:full, 3:lock
	gameround int32
	isLock    bool
	players   []code.PlayerID
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
	return g.isLock
}

// RoomLocker ...
func (g RoomInfo) RoomLocker() code.PlayerID {
	return g.locker
}

// Players ...
func (g *RoomInfo) Players() []code.PlayerID {
	return g.players
}

// Join ...
func (g *RoomInfo) Join(player *player.PlayerInfo) (int8, errorlog.ErrorMsg) {
	// err = errormsg.New()
	var err errorlog.ErrorMsg

	// maybe meet self
	if g.IsRoomLock() && g.locker != player.ID {
		err.MsgNum = RoomLockError
		err.Msg = "RoomLockError"
		return g.status, err
	} else if g.isPlayerInRoom(player) {
		err.MsgNum = SelfInRoom
		err.Msg = "SelfInRoom"
		return g.status, err
	} else if g.status == RoomFull {
		err.MsgNum = RoomFullError
		err.Msg = "RoomFullError"
		return g.status, err
	}

	g.players = append(g.players, player.ID)
	g.status = Someone
	player.InGame = g.gametype
	player.InRoom = g.id

	if int32(len(g.Players())) >= g.MaxPlayer {
		g.status = RoomFull
	}
	return g.status, err
}

// Leave ...
func (g *RoomInfo) Leave(player *player.PlayerInfo) (int8, errorlog.ErrorMsg) {
	var err errorlog.ErrorMsg

	if !g.isPlayerInRoom(player) {
		err.MsgNum = SelfNotInRoom
		err.Msg = "SelfNotInRoom"

	} else {
		count := len(g.players)
		for index := 0; index < count; index++ {
			if g.players[index] == player.ID {
				g.players = append(g.players[:index], g.players[index+1:]...)
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
func (g *RoomInfo) isPlayerInRoom(player *player.PlayerInfo) bool {
	playerID := player.ID
	for _, comparetTarget := range g.players {
		if comparetTarget == playerID {
			return true
		}
	}
	return false
}
