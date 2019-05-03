package db

// AccountGetRead stored procedure format
type AccountGetRead struct {
	Account     string
	GameAccount string
	LoginTime   int64
}

// GameAccountGetRead stored procedure format
type GameAccountGetRead struct {
	GameAccount string
	GameID      int64
	GameMoney   int64
}

// Setting stored procedure format
type Setting struct {
	Key    string
	IValue int64
	SValue string
}
