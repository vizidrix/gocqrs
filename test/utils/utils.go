package test_utils

import (
	cqrs "github.com/vizidrix/gocqrs"
	"log"
	"testing"
)

func ignore() { log.Println("") }

// Creates a chainable wrapper around the default test harness
type AggregateTestContext struct {
	T *testing.T
}

// Appends the aggregate loader to the context
type AggregateTestGiven struct {
	AggregateTestContext
	Loader cqrs.AggregateLoader
}

// Appends the pre-condition events to the context
type AggregateTestWith struct {
	AggregateTestGiven
	StoredEvents []cqrs.IEvent
	Aggregate    interface{}
}

// Appends the post-condition events to the context
type AggregateTestWhen struct {
	AggregateTestWith
	NewEvents []cqrs.IEvent
}

// Returns a context which wrapps the test harness
func AggregateTest(t *testing.T) *AggregateTestContext {
	return &AggregateTestContext{
		T: t,
	}
}

// Returns the context with appended aggregate loader
func (context *AggregateTestContext) Given(loader cqrs.AggregateLoader) *AggregateTestGiven {
	return &AggregateTestGiven{
		AggregateTestContext: *context,
		Loader:               loader,
	}
}

// Returns the context with appended pre-condition events
func (given *AggregateTestGiven) With(storedEvents ...cqrs.IEvent) *AggregateTestWith {
	eventChannel := make(chan cqrs.IEvent)
	go func() {
		defer close(eventChannel)
		for _, event := range storedEvents {
			eventChannel <- event
		}
	}()
	aggregate, err := given.Loader(nil, eventChannel)
	if err != nil {
		log.Println("Loader failed on With")
		given.T.Fail()
	}
	return &AggregateTestWith{
		AggregateTestGiven: *given,
		StoredEvents:       storedEvents,
		Aggregate:          aggregate,
	}
}

// Returns the context with appended post-condition events
func (with *AggregateTestWith) When(newEvents ...cqrs.IEvent) *AggregateTestWhen {
	return &AggregateTestWhen{
		AggregateTestWith: *with,
		NewEvents:         newEvents,
	}
}

// Builds two channels, one to listen for the result and the other to listen for errors
func (when *AggregateTestWhen) Then() (<-chan interface{}, <-chan error) {
	resultChan := make(chan interface{})
	errorChan := make(chan error)
	eventChannel := make(chan cqrs.IEvent)
	go func() {
		defer close(eventChannel)
		for _, event := range when.NewEvents {
			eventChannel <- event
		}
	}()

	go func() {
		aggregate, err := when.Loader(when.Aggregate, eventChannel)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- aggregate
		}
	}()

	return resultChan, errorChan
}

// Builds a single channel to listen for the result, errors will trigger a test fail
func (when *AggregateTestWhen) ThenFailOnError() <-chan interface{} {
	resultChan := make(chan interface{})

	eventChannel := make(chan cqrs.IEvent)
	go func() {
		defer close(eventChannel)
		for _, event := range when.NewEvents {
			eventChannel <- event
		}
	}()

	go func() {
		aggregate, err := when.Loader(when.Aggregate, eventChannel)
		if err != nil {

			log.Printf("Error loading aggregate: %s", err)
			when.T.Fail()
		} else {
			resultChan <- aggregate
		}
	}()

	return resultChan
}
