package app

import (
	"context"
	"encoding/json"
	"time"
)

// Eventer is a base interface for all events.
type Eventer interface {
	GetType() string
}

// Event is a struct that represents an event.
type Event struct {
	ID       string
	Type     string
	Data     any
	Created  time.Time
	Version  int
	User     AuthUser
	Metadata map[string]any
}

// GetType returns the type of the event.
func (e Event) GetType() string {
	return e.Type
}

// AuthUser is a struct that represents a user from a session.
type AuthUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Registry is a struct that holds a map of event types to their constructors.
type Registry struct {
	registry map[string]func() Eventer
}

// NewRegistry creates a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		registry: make(map[string]func() Eventer),
	}
}

// Add adds a list of event constructors to the registry.
func (r *Registry) Add(fn ...Eventer) *Registry {

	for _, f := range fn {
		r.registry[f.GetType()] = func() Eventer { return f }
	}

	return r
}

// Serialize serializes an event to a byte slice.
func (r *Registry) Serialize(event Eventer) (string, []byte, error) {

	if _, ok := r.registry[event.GetType()]; ok {

		b, err := json.Marshal(event)
		if err != nil {
			return "", nil, err
		}

		return event.GetType(), b, nil
	}
	return "", nil, nil
}

// Deserialize deserializes an event from a byte slice.
func (r *Registry) Deserialize(eventType string, data []byte) (Eventer, error) {
	if f, ok := r.registry[eventType]; ok {
		e := f()
		err := json.Unmarshal(data, e)
		if err != nil {
			return nil, err
		}
		return e, nil
	}
	return nil, nil
}

// EventQueue is an abstraction for an event queue (channel or distributed system).
type EventQueue interface {
	Publish(ctx context.Context, event Eventer) error
}

// EventPullHandler is a function that handles events from a queue.
type EventPullHandler func(ctx context.Context, event Eventer) error

// EventQueuePull is an abstraction for an event queue that can pull events.
type EventQueuePull interface {
	EventQueue
	Subscribe(ctx context.Context, handler EventPullHandler) error
	Start(ctx context.Context) error
}

// Store is an interface for storing events.
type Store interface {
	Save(ctx context.Context, aggregateID string, events Eventer) error
	Get(ctx context.Context, aggregateID string) ([]Event, error)
}

// ExportEventer exports an event to an Event struct.
func ExportEventer(aggID string, event Eventer) Event {
	return Event{
		ID:      aggID,
		Created: time.Now().UTC(),
		Type:    event.GetType(),
		Data:    event,
	}
}
