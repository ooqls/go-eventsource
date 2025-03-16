package eventsourcingv1

type Applicator[T any] interface {
	Apply(event Event, target *T) *EventError
}
