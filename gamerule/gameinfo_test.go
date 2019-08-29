package gamerule

func RTP(totalwin, TotalBet int64) float64 {
	return float64(totalwin) / float64(TotalBet) * 100
}
