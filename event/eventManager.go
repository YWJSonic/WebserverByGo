package event

import (
	"time"

	"../game/slotgame"
)

// Update event check loop
func Update() {
	tick := time.Tick(1e9) // 1,000,000,000 = 1 second
	for {
		select {
		case <-tick:
			slotgame.Update()
		}
	}
}
