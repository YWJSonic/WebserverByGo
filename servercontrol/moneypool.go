package moneypool

import "fmt"

var (
	// JPCount ...
	JPCount = 3
)

// PoolGroup ...
// type PoolGroup struct {
// 	Pool []MoneyPool
// }

// Pool ...
type Pool struct {
	PlayerPool int64
	ServerPool int64
	JpPool     []int64
}

// RTP Get thie pool rtp
func (p Pool) RTP() string {
	return fmt.Sprintf("%.2f", float64(p.PlayerPool)/float64(p.ServerPool)*100)
}

var jpfee = []float64{0.004, 0.006, 0.007}
var poolcache Pool

// PoolInit ...
func init() {
	poolcache = Pool{JpPool: make([]int64, JPCount)}
}

// RTPControl ...
func RTPControl(Bet, Win int64) {
	var lastbet = Bet

	for i := range poolcache.JpPool {
		lastbet -= int64(float64(Bet) * jpfee[i])
		poolcache.JpPool[i] += int64(float64(Bet) * jpfee[i])
	}

	poolcache.ServerPool += lastbet
	poolcache.PlayerPool += Win
	fmt.Println("MoneyPool", poolcache.RTP())
}
