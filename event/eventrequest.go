package event

import (
	"time"

	"gitlab.com/ServerUtility/eventinfo"
	"gitlab.com/WeberverByGo/db"
)

var eventList []eventinfo.Info

// Update event check loop
func Update() {
	tick1 := time.Tick(1e9) //  1 second
	for {
		select {
		case <-tick1:
			// go slotgame.Update()
			for _, event := range eventList {
				if event.IsLaunch {
					continue
				}
			}
		}
		db.SQLSelect()
	}
}

// AddNewEvent add new event
func AddNewEvent(event eventinfo.Info) {
	eventList = append(eventList, event)
}
