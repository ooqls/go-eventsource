package eventsv1

import (
	"context"

	"github.com/ooqls/go-eventsource/eventsourcingv1"
	"go.uber.org/zap"
)

type EventApplicator[T any] struct {
	adapter eventsourcingv1.Adapter[T]
	source  eventsourcingv1.EventSource
	r       Reader
}

func NewApplicator[T any](r Reader, source eventsourcingv1.EventSource) *EventApplicator[T] {
	return &EventApplicator[T]{
		source: source,
		r:      r,
	}
}

func (a *EventApplicator[T]) Apply(ctx context.Context, ent *T) *ApplicatorError {
	next, err := a.r.Get(ctx, a.adapter.GetEntityId(*ent))
	if err != nil {
		return &ApplicatorError{err, nil}
	}

	ev, err := next()
	for ev != nil {
		err = a.adapter.Apply(*ev, ent)
		if err != nil {
			l.Error("failed to apply event", zap.Error(err), zap.String("event", string(ev.Key)))
		}
		ev, err = next()
		if err != nil {
			return &ApplicatorError{err, nil}
		}
	}
	return nil
}
