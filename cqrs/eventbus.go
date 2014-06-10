package cqrs

import (
	//"fmt"
	"errors"
)

var (
	ErrInvalidEventBusState     = errors.New("Event bus not properly initialized")
	ErrInvalidNilTypeFilter     = errors.New("Cannot subscribe with nil type filter")
	ErrInvalidEmptyTypeFilter   = errors.New("Cannot subscribe with an empty type filter")
	ErrInvalidNilPublishedEvent = errors.New("Cannot publish a nil event")
	ErrInvalidNilUnSubscribe    = errors.New("Cannot unsubscribe a nil subscriber")
	ErrInvalidUnSubscribe       = errors.New("Cannot unsubscribe an invalid subscriber")
)

var EventBus EventRouter

/*
func init() {
	// Auto init the event bus as a global service
	EventBus = NewDefaultedEventBus()
	EventBus.Listen()
}
*/

type EventChanFactory func() chan Event
type CancelChanFactory func() chan struct{}

type EventFilterer interface {
	Predicate(event Event) bool
}

type Subscriber interface {
	EventChan() <-chan Event
	Publish(event Event)
	Filter() EventFilterer
	Cancel()
}

// EventRouter provides the abstraction over a bus for clients to connect against
type EventRouter interface {
	Listen()                                            // Iterates across the Step function in a goroutine loop
	Step()                                              // Grabs the next operation from the queue and processes it
	Publish(event Event) error                          // Pushes a copy of the event to all relevant subscribers
	Subscribe(filter EventFilterer) (Subscriber, error) // Registers a subscriber using it's filter
	UnSubscribe(subscriber Subscriber) error            // UnRegisters a subscriber from receiving events
}

type eventTypesFilter struct {
	eventTypes []uint32
}

func (filter *eventTypesFilter) Predicate(event Event) bool {
	for _, eventType := range filter.eventTypes {
		if eventType == event.GetEventType() {
			return true
		}
	}
	return false
}

func ByEventTypes(eventTypes ...uint32) EventFilterer {
	if len(eventTypes) == 0 {
		return nil
	}
	return &eventTypesFilter{
		eventTypes: eventTypes,
	}
}

type aggregateIdsFilter struct {
	aggregateIds []uint64
}

func (filter *aggregateIdsFilter) Predicate(event Event) bool {
	for _, aggregateId := range filter.aggregateIds {
		if aggregateId == event.GetId() {
			return true
		}
	}
	return false
}

func ByAggregateIds(aggregateIds ...uint64) EventFilterer {
	if len(aggregateIds) == 0 {
		return nil
	}
	return &aggregateIdsFilter{
		aggregateIds: aggregateIds,
	}
}

type subscription struct {
	eventBus  EventRouter
	filter    EventFilterer
	eventChan chan Event
}

func (s *subscription) EventChan() <-chan Event {
	return s.eventChan
}

func (s *subscription) Publish(event Event) {
	s.eventChan <- event
}

func (s *subscription) Filter() EventFilterer {
	return s.filter
}

func (s *subscription) Cancel() {
	s.eventBus.UnSubscribe(s)
}

// EventRouter implementation that uses Go chans to provide routing
type channelEventBus struct {
	subscribeChan    chan Subscriber
	unsubscribeChan  chan Subscriber
	publishChan      chan Event
	subscriptions    []Subscriber
	eventChanFactory EventChanFactory
}

func NewDefaultedEventBus() *channelEventBus {
	return NewChannelEventBus(
		make(chan Subscriber),
		make(chan Subscriber),
		make(chan Event),
		func() chan Event { return make(chan Event) },
	)
}

func NewChannelEventBus(
	subscriptionChan chan Subscriber,
	unsubscriptionChan chan Subscriber,
	publishChan chan Event,
	eventChanFactory EventChanFactory,
) *channelEventBus {
	bus := &channelEventBus{
		subscribeChan:    subscriptionChan,
		unsubscribeChan:  unsubscriptionChan,
		publishChan:      publishChan,
		subscriptions:    make([]Subscriber, 0, 10),
		eventChanFactory: eventChanFactory,
	}
	return bus
}

func (c *channelEventBus) Step() {
	select { // Synchronized select for event bus mutable actions
	case subscription := <-c.subscribeChan:
		{
			c.subscriptions = append(c.subscriptions, subscription)
		}
	case subscription := <-c.unsubscribeChan:
		{
			for index, s := range c.subscriptions {
				if subscription == s {
					c.subscriptions = append(c.subscriptions[:index], c.subscriptions[index+1:]...)
				}
			}
		}
	case event := <-c.publishChan:
		{
			for _, subscription := range c.subscriptions {
				if subscription.Filter().Predicate(event) {
					subscription.Publish(event)
				}
			}
		}
	}
}

func (c *channelEventBus) Listen() {
	go func() {
		for {
			c.Step()
		}
	}()
}

func (c *channelEventBus) Publish(event Event) error {
	if event == nil {
		return ErrInvalidNilPublishedEvent
	}
	select {
	case c.publishChan <- event:
	default:
		return ErrInvalidEventBusState
	}
	return nil
}

func (c *channelEventBus) Subscribe(filter EventFilterer) (Subscriber, error) {
	if filter == nil {
		return nil, ErrInvalidNilTypeFilter
	}
	handle := &subscription{
		eventBus:  c,
		filter:    filter,
		eventChan: c.eventChanFactory(),
	}
	c.subscribeChan <- handle
	return handle, nil
}

func (c *channelEventBus) UnSubscribe(subscription Subscriber) error {
	if subscription == nil {
		return ErrInvalidNilUnSubscribe
	}
	c.unsubscribeChan <- subscription
	return nil
}
