package gamesystem

// IsTotalWinLimit if int limit return true
func IsTotalWinLimit(betMoney, totalWin int64) bool {
	if (totalWin / betMoney) > 100 {
		return true
	}
	return false
}
