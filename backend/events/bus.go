package events

import (
	"context"
	"log/slog"
)

type IObserver interface {
	OnEvent(ctx context.Context, event Event) error
}

type EventBus struct {
	observers []IObserver
	logger    *slog.Logger
}

func NewEventBus(logger *slog.Logger) *EventBus {
	return &EventBus{logger: logger}
}

func (b *EventBus) Subscribe(observer IObserver) {
	b.observers = append(b.observers, observer)
}

func (b *EventBus) Publish(ctx context.Context, event Event) {
	for _, o := range b.observers {
		if err := o.OnEvent(ctx, event); err != nil {
			b.logger.Error("observer failed",
				slog.String("event_type", event.Type),
				slog.Any("error", err),
			)
		}
	}
}
