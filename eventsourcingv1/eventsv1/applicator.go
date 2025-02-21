package eventsv1

import (
	"context"

	"github.com/ooqls/go-eventsource/eventsourcingv1"
)

type Applicator interface {
	Apply(ctx context.Context, ent eventsourcingv1.Entity) error
}

type SQLApplicator struct {
	source eventsourcingv1.EventSource
	r      Reader
}

func NewSQLApplicator(r Reader, source eventsourcingv1.EventSource) *SQLApplicator {
	return &SQLApplicator{
		source: source,
		r:      r,
	}
}

func (a *SQLApplicator) Apply(ctx context.Context, ent eventsourcingv1.Entity) error {
	next, err := a.r.Get(ctx, ent.GetId())
	if err != nil {
		return err
	}
	ev := next()
	for ev != nil {
		ent.Apply(*ev)
		ev = next()
	}
	return nil
}
