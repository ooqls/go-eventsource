package eventsv1

import "github.com/ooqls/go-eventsource/eventsourcingv1"

type ApplicatorError struct {
	error
	Events []eventsourcingv1.Event
}
