package events

// Fetcher is an interface for fetching a batch of events from an external source.
type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

// Processor is an interface for processing a single event.
type Processor interface {
	Process(e Event) error
}

// Type represents the type of an event.
type Type int

const (
	Unknown Type = iota
	Message
)

// Event represents a single event in the system, such as a user message.
type Event struct {
	Type Type
	Text string
	Meta any
}
