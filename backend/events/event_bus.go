package events

import "context"

type IObserver interface {
	OnEvent(ctx context.Context, event Event) error
}

type EventBus struct {
	observers []IObserver
}

func NewEventBus() *EventBus {
	return &EventBus{}
}

func (b *EventBus) Subscribe(observer IObserver) {
	b.observers = append(b.observers, observer)
}

func (b *EventBus) Publish(ctx context.Context, event Event) {
	for _, o := range b.observers {
		if err := o.OnEvent(ctx, event); err != nil {
			// log but never block the main flow
		}
	}
}
