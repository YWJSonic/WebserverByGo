package event

import (
	"time"

	"../db"
	"../game/slotgame"
)

var EventList []Event

var count int

// Update event check loop
func Update() {
	tick := time.Tick(1e9) // 1,000,000,000 = 1 second
	// tick2 := time.Tick(5e8) // 5,000,000,00 = 0.5 second
	for {
		select {
		case <-tick:
			slotgame.Update()
			// case dbLogQuery := <-db.QueryLogChan:
			// 	db.ExecLog(dbLogQuery)
			// case dbgameQuery := <-db.WriteGameChan:
			// 	db.CallWrite(dbgameQuery.Quary, dbgameQuery.Args...)
			// case dbpayQuary := <-db.WritePayChan:
			// 	db.CallWrite(dbpayQuary.Quary, dbpayQuary.Args...)
		}
		db.SQLSelect()
	}
}
