package cqrs_test

import (
	. "github.com/vizidrix/gocqrs/cqrs"
	"testing"
)

var (
	DOMAIN      uint32 = 0x11111111
	E_TestEvent uint32 = E(1, 1)
)

type MockEventBus struct {
	SubscriptionChan   chan Subscriber
	UnSubscriptionChan chan Subscriber
	PublishChan        chan Event
	EventChan          chan Event
	CancelChan         chan struct{}
}

func NewMockEventBus() *MockEventBus {
	return &MockEventBus{
		SubscriptionChan:   make(chan Subscriber, 1),
		UnSubscriptionChan: make(chan Subscriber, 1),
		PublishChan:        make(chan Event, 1),
		EventChan:          make(chan Event, 1),
		CancelChan:         make(chan struct{}),
	}
}

func (mock *MockEventBus) Create() EventRouter {
	return NewChannelEventBus(
		mock.SubscriptionChan,
		mock.UnSubscriptionChan,
		mock.PublishChan,
		func() chan Event { return mock.EventChan },
	)
}

type MockSubscriber uint32

func (mock *MockSubscriber) EventChan() <-chan Event {
	return nil
}

func (mock *MockSubscriber) Publish(event Event) {
	return
}

func (mock *MockSubscriber) Domain() uint32 {
	return uint32(*mock)
}

func (mock *MockSubscriber) Filter() EventFilterer {
	return nil
}

func (mock *MockSubscriber) Cancel() {
	return
}

type TestEvent struct {
	EventMemento
	Value string
}

func NewTestEvent(id uint64, version uint32, value string) TestEvent {
	return TestEvent{
		EventMemento: NewEvent(DOMAIN, id, version, E_TestEvent),
		Value:        value,
	}
}

func Test_Should_filter_non_matching_events_ByEventType(t *testing.T) {
	event := NewTestEvent(1, 1, "test")
	filter := ByEventTypes(10)

	if filter.Predicate(event) {
		t.Errorf("Expected filter ByEventType for [ %d ] with [ %d ] to return false",
			E_TestEvent, event.GetEventType())
	}
}

func Test_Should_return_nil_for_empty_eventtype_filter_set(t *testing.T) {
	if filter := ByEventTypes(); filter != nil {
		t.Errorf("Should have returned nil for empty event type filter")
	}
}

func Test_Should_return_nil_for_empty_aggregateid_filter_set(t *testing.T) {
	if filter := ByAggregateIds(); filter != nil {
		t.Errorf("Should have returned nil for empty aggregate id filter")
	}
}

func Test_Should_filter_matching_events_ByEventType(t *testing.T) {
	event := NewTestEvent(1, 1, "test")
	filter := ByEventTypes(E_TestEvent)

	if !filter.Predicate(event) {
		t.Errorf("Expected filter ByEventType for [ %d ] with [ %d ] to return true",
			E_TestEvent, event.GetEventType())
	}
}

func Test_Should_filter_matching_events_ByAggregateIds(t *testing.T) {
	event := NewTestEvent(1, 1, "Test")
	filter := ByAggregateIds(1)

	if !filter.Predicate(event) {
		t.Errorf("Expected filter ByAggregateIds for [ %d ] with [ %d ] to return true",
			12, event.GetId())
	}
}

func Test_Should_return_false_for_unmatched_events_ByEventType(t *testing.T) {
	event := NewTestEvent(1, 1, "Test")
	filter := ByEventTypes(E_TestEvent + 1)

	if filter.Predicate(event) {
		t.Errorf("Expected filter ByEventType for [ %d ] with [ %d ] to return false", 10, event.GetEventType())
	}
}

func Test_Should_return_false_for_unmatched_events_ByAggregateIds(t *testing.T) {
	event := NewTestEvent(1, 1, "Test")
	filter := ByAggregateIds(10)

	if filter.Predicate(event) {
		t.Errorf("Expected filter ByAggregateIds for [ %d ] with [ %d ] to return false", 10, event.GetId())
	}
}

func Test_Should_return_an_error_on_nil_filter(t *testing.T) {
	eventbus := NewMockEventBus().Create()
	_, err := eventbus.Subscribe(DOMAIN, nil)

	if err != ErrInvalidNilTypeFilter {
		t.Errorf("Should have returned an error for nil type filter but was [ %v ]\n", err)
	}
}

func Test_Should_return_subscription_token_for_valid_filter_set(t *testing.T) {
	eventbus := NewMockEventBus().Create()
	handle, err := eventbus.Subscribe(DOMAIN, ByEventTypes(10))

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
	eventbus := NewMockEventBus().Create()
	if err := eventbus.Publish(nil); err != ErrInvalidNilPublishedEvent {
		t.Errorf("Should have returned an error for nil event in publish but was [ %v ]\n", err)
		return
	}
}

func Test_Should_not_return_error_from_valid_publish(t *testing.T) {
	eventbus := NewMockEventBus().Create()
	event := NewTestEvent(1, 1, "publish test")
	if err := eventbus.Publish(event); err != nil {
		t.Errorf("Should not return error from valid publish")
	}
}

func Test_Should_receive_matching_event_when_published(t *testing.T) {
	eventbus := NewMockEventBus().Create()
	handle, _ := eventbus.Subscribe(DOMAIN, ByEventTypes(E_TestEvent))
	eventbus.Step()
	expected := NewTestEvent(1, 1, "publish test")
	eventbus.Publish(expected)
	eventbus.Step()
	select {
	case actual := <-handle.EventChan():
		{
			if expected != actual {
				t.Errorf("\nExpected\t[ %v ]\n\tbut was\t[ %v ]\n", expected, actual)
			}
		}
	default:
		{
			t.Errorf("Should not have hit default case\n")
		}
	}
}

func Test_Should_return_error_when_unsubscribing_nil_subscriber(t *testing.T) {
	eventbus := NewMockEventBus().Create()
	if err := eventbus.UnSubscribe(nil); err == nil {
		t.Errorf("Should have raised an error when unsubscribing nil subscriber")
	}
}

func Test_Should_stop_receiving_matching_event_when_canceled(t *testing.T) {
	eventbus := NewMockEventBus().Create()
	handle, _ := eventbus.Subscribe(DOMAIN, ByEventTypes(E_TestEvent))
	eventbus.Step()
	eventbus.UnSubscribe(handle)
	eventbus.Step()
	event := NewTestEvent(1, 1, "publish test")
	eventbus.Publish(event)
	eventbus.Step()
	select {
	case actual := <-handle.EventChan():
		{
			t.Errorf("\nExpected subscription cancel but got\t[ %v ]\n", actual)
		}
	default:
		{ // Should not have published to this subscriber
		}
	}
}
