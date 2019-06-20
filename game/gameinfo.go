package game

// Attach game att data
type Attach struct {
	FreeCount    int `json:"freecount"`
	IsLockBet    int `json:"islockbet"`
	LockBetIndex int `json:"lockbetindex"`
}
