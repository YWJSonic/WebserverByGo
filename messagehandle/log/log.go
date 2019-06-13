package log

// Login event
const (
	Login      = 10
	GetPlayer  = 11
	Exchange   = 12
	JoinGame   = 13
	LeaveGame  = 14
	GameResult = 15
	CheckOut   = 16
)

// LogInfo struect
type LogInfo struct {
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
