package events

type Event struct {
	Name string
	Data any
}

type subscriber struct {
	channel chan Event
}
