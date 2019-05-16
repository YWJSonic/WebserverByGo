package event

// AddNewEvent add new event
func AddNewEvent(event EventInfo) {
	eventList = append(eventList, event)
}
