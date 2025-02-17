package inmem

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	app "github.com/DustinHigginbotham/event-gen-user-example/gen"
)

type Store struct {
	sync.RWMutex
	events   map[string][]app.Event
	registry *app.Registry
}

// Get implements app.Store.
func (s *Store) Get(ctx context.Context, aggregateID string) ([]app.Event, error) {
	s.Lock()
	defer s.Unlock()

	var err error

	evs := make([]app.Event, len(s.events[aggregateID]))
	for i, ev := range s.events[aggregateID] {
		evs[i] = ev

		evs[i].Data, err = s.registry.Deserialize(ev.Type, ev.Data.([]byte))
		if err != nil {
			return nil, err
		}
	}

	return evs, nil
}

// Save implements app.Store.
func (s *Store) Save(ctx context.Context, aggregateID string, ev app.Eventer) error {
	s.Lock()
	defer s.Unlock()

	version := len(s.events[aggregateID])

	event := app.ExportEventer(aggregateID, ev)
	event.Version = version + +1
	b, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	event.Data = b

	s.events[aggregateID] = append(s.events[aggregateID], event)
	return nil
}

var _ app.Store = new(Store)

func NewStore(registry *app.Registry) *Store {
	return &Store{
		events:   make(map[string][]app.Event),
		registry: registry,
	}
}

func NewTestState(ctx context.Context, registry *app.Registry, id string, events []app.Eventer) *Store {
	store := NewStore(registry)
	for _, ev := range events {
		store.Save(ctx, id, ev)
	}
	return store
}

func (s *Store) MustAppend(ctx context.Context, aggregateID string, events []app.Eventer) *Store {
	for _, ev := range events {
		if err := s.Save(ctx, aggregateID, ev); err != nil {
			panic(err)
		}
	}
	return s
}
