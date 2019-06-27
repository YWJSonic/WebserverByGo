package game

// Attach game att data
type Attach struct {
	FreeCount    int64 `json:"freecount"`
	IsLockBet    int64 `json:"islockbet"`
	LockBetIndex int64 `json:"lockbetindex"`
}
