package event

import (
	"time"

	"gitlab.com/WeberverByGo/db"
)

var eventList []EventInfo

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

func init() {
	eventList = append(eventList, EventInfo{
		IsLaunch:  false,
		EventType: LoopEvent,
		LoopTime:  5 * time.Second,
		Do: func(interface{}) interface{} {
			// slotgame.Update()
			return nil
		}})

	eventList = append(eventList, EventInfo{
		IsLaunch:  false,
		EventType: AtTimeEvent,
		DoTime:    time.Date(2019, 5, 15, 14, 58, 0, 0, time.Local).Unix(),
		Do: func(interface{}) interface{} {
			// slotgame.Update()
			return nil
		}})
}

// AddNewEvent add new event
func AddNewEvent(event EventInfo) {
	eventList = append(eventList, event)
}
