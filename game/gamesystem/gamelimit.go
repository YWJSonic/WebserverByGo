package gamesystem

// IsInTotalMoneyWinLimit if int limit return true
func IsInTotalMoneyWinLimit(betMoney, totalWin int64) bool {
	if totalWin > 2000000 {
		return false
	}
	return true
}

// IsInTotalBetRateWinLimit if int limit return true
func IsInTotalBetRateWinLimit(betMoney, totalWin int64) bool {
	if (totalWin / betMoney) > 100 {
		return false
	}
	return true
}
