package inmem

import (
	"context"
	"sync"

	app "github.com/DustinHigginbotham/event-gen-user-example/gen"
)

type ev struct {
	app.Eventer
	ctx context.Context
}

var _ app.EventQueuePull = new(InMemQueue)

type InMemQueue struct {
	mu       sync.Mutex
	events   chan *ev
	handlers []app.EventPullHandler
	registry *app.Registry
}

var _ app.EventQueue = new(InMemQueue)

func NewQueue(app *app.App, bufferSize int) *InMemQueue {
	return &InMemQueue{
		events:   make(chan *ev, bufferSize),
		registry: app.Registry(),
		handlers: []app.EventPullHandler{},
	}
}

func (i *InMemQueue) Publish(ctx context.Context, event app.Eventer) error {
	select {
	case i.events <- &ev{Eventer: event, ctx: ctx}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (i *InMemQueue) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case e := <-i.events:
			i.mu.Lock()
			for _, h := range i.handlers {
				h(e.ctx, e.Eventer)
			}
			i.mu.Unlock()
		}
	}
}

func (i *InMemQueue) Subscribe(ctx context.Context, h app.EventPullHandler) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.handlers = append(i.handlers, h)
	return nil
}
