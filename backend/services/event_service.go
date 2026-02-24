package services

import (
	"log/slog"
	"sync"
)

// EventType identifies the kind of event being published.
type EventType string

const (
	EventAnswerSubmitted   EventType = "answer.submitted"
	EventStreakMilestone   EventType = "streak.milestone"
	EventPlanCompleted     EventType = "plan.completed"
	EventProficiencyChange EventType = "proficiency.changed"
)

// Event is a lightweight message published to subscribers.
type Event struct {
	Type    EventType
	UserID  int
	Payload map[string]interface{}
}

// EventHandler is a callback invoked when an event is published.
type EventHandler func(Event)

// EventService is a simple in-process pub/sub bus.
// All handlers are executed asynchronously in goroutines so publishing never blocks.
type EventService struct {
	mu          sync.RWMutex
	subscribers map[EventType][]EventHandler
}

// NewEventService creates a new event bus.
func NewEventService() *EventService {
	return &EventService{
		subscribers: make(map[EventType][]EventHandler),
	}
}

// Subscribe registers a handler for the given event type.
func (es *EventService) Subscribe(eventType EventType, handler EventHandler) {
	es.mu.Lock()
	defer es.mu.Unlock()
	es.subscribers[eventType] = append(es.subscribers[eventType], handler)
}

// Publish sends an event to all registered handlers asynchronously.
// Each handler runs in its own goroutine; panics are recovered and logged.
func (es *EventService) Publish(event Event) {
	es.mu.RLock()
	handlers := es.subscribers[event.Type]
	es.mu.RUnlock()

	for _, h := range handlers {
		go func(handler EventHandler) {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("event handler panicked",
						slog.String("event_type", string(event.Type)),
						slog.Int("user_id", event.UserID),
						slog.Any("panic", r),
					)
				}
			}()
			handler(event)
		}(h)
	}
}
