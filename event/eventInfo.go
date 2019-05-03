package event

type Event struct {
	StartTime int64

	Do func(interface{}) interface{}
}
