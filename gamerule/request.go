package gamerule

// InitAttach init game attach
func InitAttach(playerID int64) {

}

// ConvertToGameAttach ...
func ConvertToGameAttach(playerID int64, attinfo []map[string]interface{}) AttachInfo {
	return attachDataToAttachInfo(playerID, attinfo)
}

// SetInfo ...
func SetInfo(gameIndex int, att map[string]interface{}) {
}

// GetInitScroll ...
func GetInitScroll() interface{} {
	var scrollmap map[string][][]int
	return scrollmap
}

// GetInitBetRate init info
func GetInitBetRate() interface{} {
	tmp := make(map[string]interface{})
	return tmp
}

// GetBetMoney ...
func GetBetMoney(betIndex int64) int64 {
	betrate := BetRate()
	return betrate[betIndex]
}

// BetRate ...
func BetRate() []int64 {
	return betRate
}

// GameRequest game server api return game result, game attach, totalwin
func GameRequest(playerID, betIndex int64, attach []map[string]interface{}) (map[string]interface{}, []map[string]interface{}, map[string]int64) {
	var result map[string]interface{}
	var otherdata map[string]int64
	return result, attach, otherdata
}
