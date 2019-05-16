package log

import (
	"time"

	"../db"
)

const (
	Login      = 10
	GetPlayer  = 11
	Exchange   = 12
	JoinGame   = 13
	LeaveGame  = 14
	GameResult = 15
	CheckOut   = 16
)

type Log struct {
	Account       string
	PlayerID      int64
	Time          int64
	ActivityEvent int
	IValue1       int64
	IValue2       int64
	IValue3       int64
	SValue1       string // varchar(128)
	SValue2       string // varchar(128)
	SValue3       string // varchar(128)
	Msg           string // text
}

// New Default log
func New(ActivityEvent int) Log {
	return Log{
		Time:          time.Now().Unix(),
		ActivityEvent: ActivityEvent,
	}
}

// SaveLog Save to db
func SaveLog(log Log) {
	db.SetLog(
		log.Account,
		log.PlayerID,
		log.Time,
		log.ActivityEvent,
		log.IValue1,
		log.IValue2,
		log.IValue3,
		log.SValue1,
		log.SValue2,
		log.SValue3,
		log.Msg)
}
