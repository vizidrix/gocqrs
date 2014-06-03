package cqrs_test

import (
	//"fmt"
	"testing"
	. "github.com/vizidrix/gocqrs/cqrs"
)

var (
	DOMAIN uint32 = 0x11111111
	E_TestEvent uint32 = E(1, 1)
)

type TestEvent struct {
	EventMemento
	Value string
}

func NewTestEvent(id uint64, version uint32, value string) TestEvent {
	return TestEvent {
		EventMemento: NewEvent(DOMAIN, id, version, E_TestEvent),
		//Event: cqrs.NewEvent(DOMAIN, id, version, E_TestEvent),
		Value: value,
	}
}

func Test_Should_filter_non_matching_events_ByEventType(t *testing.T) {
	event := NewTestEvent(1, 1, "test")

	filter := ByEventTypes(10)

	if filter.Predicate(event) {
		t.Errorf("Expected filter ByEventType to return false")
	}
}

func Test_Should_filter_matching_events_ByEventType(t *testing.T) {
	event := NewTestEvent(1, 1, "test")

	filter := ByEventTypes(20)

	if filter.Predicate(event) {
		t.Errorf("Expected filter ByEventType to return true")
	}
}

func Test_Should_return_an_error_on_nil_filter(t *testing.T) {
	_, err := EventBus.Subscribe(nil)

	if err != ErrInvalidNilTypeFilter {
		t.Errorf("Should have returned an error for nil type filter but was [ %v ]\n", err)
	}
}
/*
func Test_Should_return_an_error_on_empty_filter(t *testing.T) {
	_, err := cqrs.EventBus.Subscribe(...[]uint32{})

	if err != cqrs.ErrInvalidEmptyTypeFilter {
		t.Errorf("Should have returned an error for empty type filter but was [ %v ]\n", err)
	}
}
*/

func Test_Should_return_subscription_token_for_valid_filter_set(t *testing.T) {
	//handle, err := cqrs.EventBus.Subscribe(
	//	cqrs.ByEventTypes(10, 20, 40),
	//	cqrs.ByAggregateIds(20, 21, 23),
	//	)
	handle, err := EventBus.Subscribe(ByEventTypes(10))

	if err != nil {
		t.Errorf("Should not have err but was [ %s ]\n", err)
		return
	}
	if handle == nil {
		t.Errorf("Should have returned a non nil handle\n")
		return
	}
}

func Test_Should_return_error_when_publishing_nil_event(t *testing.T) {
	err := EventBus.Publish(nil)

	if err != ErrInvalidNilPublishedEvent {
		t.Errorf("Should have returned an error for nil event in publish but was [ %v ]\n", err)
		return
	}
}

func Test_Should_receive_matching_event_when_published(t *testing.T) {
	subscriptionChan := make(chan *Subscription, 1)
	publishChan := make(chan Event, 1)
	eventChan := make(chan Event, 1)
	cancelChan := make(chan struct{})
	eventbus := NewChannelEventBus(
		subscriptionChan, 
		publishChan,
		func() chan Event { return eventChan },
		func() chan struct{} { return cancelChan },
	)
	handle, err := eventbus.Subscribe(ByEventTypes(
		E_TestEvent,
		))
	if err != nil {
		t.Errorf("Should not return error from subscribe [ %s ]\n")
		return
	}
	eventbus.Step()
	expected := NewTestEvent(1, 1, "publish test")
	if err := eventbus.Publish(expected); err != nil {
		t.Errorf("Should not return error from publish [ %s ]\n", err)
		return
	}
	eventbus.Step()
	select {
		case actual := <-handle.EventChan: {
			if expected != actual {
				t.Errorf("\nExpected\t[ %v ]\n\tbut was\t[ %v ]\n", expected, actual)
				return
			}
		}
		default: {
			t.Errorf("Should not have hit default case\n")
		}
	}
}

