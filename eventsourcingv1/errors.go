package eventsourcingv1

type EventError struct {
	error
	Events []Event
}
