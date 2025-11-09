package event

const (
	EventAddToGroup      = "group.added"
	EventRemoveFromGroup = "group.removed"
)

type Event struct {
	Type string
	Data any
}

type Bus struct {
	bus chan Event
}

func NewEventBus() *Bus {
	return &Bus{
		bus: make(chan Event),
	}
}

func (e *Bus) Publish(event Event) {
	e.bus <- event
}

func (e *Bus) Subscribe() <-chan Event {
	return e.bus
}
