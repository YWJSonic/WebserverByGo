package gamesystem

// IsInTotalMoneyWinLimit if int limit return true
func IsInTotalMoneyWinLimit(betMoney, totalWin, limitScore int64) bool {
	if totalWin > limitScore {
		return false
	}
	return true
}

// IsInTotalBetRateWinLimit if int limit return true
func IsInTotalBetRateWinLimit(betMoney, totalWin, limitBetRate int64) bool {
	if (totalWin / betMoney) > limitBetRate {
		return false
	}
	return true
}
