package gamerule

// Result att 0: freecount
func logicResult(betMoney int64, attinfo *AttachInfo) map[string]interface{} {
	var result map[string]interface{}
	return result
}

// outputGame out put normal game result, mini game status, totalwin
func outputGame(betMoney int64, attinfo *AttachInfo) (map[string]interface{}, map[string]interface{}, int64) {
	var totalScores int64
	var result map[string]interface{}
	otherdata := make(map[string]interface{})

	return result, otherdata, totalScores
}

// outRespin out put respin result and totalwin
func outRespin(betMoney int64, attinfo *AttachInfo) ([]interface{}, int64) {
	var totalScores int64
	var result []interface{}
	return result, totalScores
}
