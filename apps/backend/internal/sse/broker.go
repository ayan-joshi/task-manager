package sse

import (
	"sync"

	"github.com/google/uuid"

	"taskmanager/internal/dto"
)

// Broker fans out task events to connected SSE clients. Events are scoped per
// user: a client only receives events for the user it is subscribed as, which
// preserves the same ownership isolation enforced by the REST API.
type Broker struct {
	mu          sync.RWMutex
	subscribers map[uuid.UUID]map[chan dto.TaskEvent]struct{}
}

// NewBroker constructs an empty Broker.
func NewBroker() *Broker {
	return &Broker{subscribers: make(map[uuid.UUID]map[chan dto.TaskEvent]struct{})}
}

// Subscribe registers a new client channel for the given user and returns it
// along with an unsubscribe function the caller must defer.
func (b *Broker) Subscribe(userID uuid.UUID) (chan dto.TaskEvent, func()) {
	ch := make(chan dto.TaskEvent, 16)

	b.mu.Lock()
	if b.subscribers[userID] == nil {
		b.subscribers[userID] = make(map[chan dto.TaskEvent]struct{})
	}
	b.subscribers[userID][ch] = struct{}{}
	b.mu.Unlock()

	unsubscribe := func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if set, ok := b.subscribers[userID]; ok {
			delete(set, ch)
			if len(set) == 0 {
				delete(b.subscribers, userID)
			}
		}
		close(ch)
	}
	return ch, unsubscribe
}

// Publish delivers an event to all of a user's subscribers. Delivery is
// non-blocking: a slow consumer is skipped rather than stalling the publisher.
func (b *Broker) Publish(userID uuid.UUID, event dto.TaskEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.subscribers[userID] {
		select {
		case ch <- event:
		default:
		}
	}
}
