package eventsv1

import "github.com/ooqls/go-eventsource/eventsourcingv1"

type ApplicatorError struct {
	error
	events []eventsourcingv1.Event
}
