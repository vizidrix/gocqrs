package cqrs

import (
	//"fmt"
	"errors"
)

var (
	ErrInvalidEventBusState     = errors.New("Event bus not properly initialized")
	ErrInvalidNilTypeFilter     = errors.New("Cannot subscribe with nil type filter")
	ErrInvalidSubscribeDomain   = errors.New("Cannot subscribe with an invalid domain")
	ErrInvalidEmptyTypeFilter   = errors.New("Cannot subscribe with an empty type filter")
	ErrInvalidNilPublishedEvent = errors.New("Cannot publish a nil event")
	ErrInvalidPublishDomain     = errors.New("Cannot publish to an invalid domain")
	ErrInvalidNilUnSubscribe    = errors.New("Cannot unsubscribe a nil subscriber")
	ErrInvalidUnSubscribe       = errors.New("Cannot unsubscribe an invalid subscriber")
	ErrInvalidUnSubscribeDomain = errors.New("Cannot unsubscribe with an invalid domain")
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
	Domain() uint32
	Filter() EventFilterer
	Cancel()
}

// EventRouter provides the abstraction over a bus for clients to connect against
type EventRouter interface {
	Listen()                                                           // Iterates across the Step function in a goroutine loop
	Step()                                                             // Grabs the next operation from the queue and processes it
	Publish(event Event) error                                         // Pushes a copy of the event to all relevant subscribers
	ValidateDomain(domain uint32) bool                                 //Returns a boolean of whether a domain is active in the bus
	Subscribe(domain uint32, filter EventFilterer) (Subscriber, error) // Registers a subscriber using it's filter
	UnSubscribe(subscriber Subscriber) error                           // UnRegisters a subscriber from receiving events
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
	domain    uint32
	filter    EventFilterer
	eventChan chan Event
}

func (s *subscription) EventChan() <-chan Event {
	return s.eventChan
}

func (s *subscription) Publish(event Event) {
	s.eventChan <- event
}

func (s *subscription) Domain() uint32 {
	return s.domain
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
	subscriptions    map[uint32][]Subscriber
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
		subscriptions:    make(map[uint32][]Subscriber),
		eventChanFactory: eventChanFactory,
	}
	return bus
}

func (c *channelEventBus) Step() {
	select { // Synchronized select for event bus mutable actions
	case subscription := <-c.subscribeChan:
		{
			if c.subscriptions[subscription.Domain()] == nil {
				c.subscriptions[subscription.Domain()] = make([]Subscriber, 0, 10)
			}
			c.subscriptions[subscription.Domain()] = append(c.subscriptions[subscription.Domain()], subscription)
		}
	case subscription := <-c.unsubscribeChan:
		{
			if c.subscriptions[subscription.Domain()] == nil {
				panic(ErrInvalidUnSubscribeDomain)
			}
			for index, s := range c.subscriptions[subscription.Domain()] {
				if subscription == s {
					c.subscriptions[subscription.Domain()] = append(
						c.subscriptions[subscription.Domain()][:index],
						c.subscriptions[subscription.Domain()][index+1:]...,
					)
				}
			}
		}
	case event := <-c.publishChan:
		{
			if c.subscriptions[event.GetDomain()] == nil {
				/*
					panic(ErrInvalidPublishDomain)
				*/
				break
			}
			for _, subscription := range c.subscriptions[event.GetDomain()] {
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

	/*
		if !c.ValidateDomain(event.GetDomain()) {
			return ErrInvalidPublishDomain
		}
	*/
	select {
	case c.publishChan <- event:
	default:
		return ErrInvalidEventBusState
	}
	return nil
}

func (c *channelEventBus) ValidateDomain(domain uint32) bool {
	if _, active := c.subscriptions[domain]; !active {
		return false
	}
	return true
}

func (c *channelEventBus) Subscribe(domain uint32, filter EventFilterer) (Subscriber, error) {
	if domain == 0 {
		return nil, ErrInvalidSubscribeDomain
	}
	if filter == nil {
		return nil, ErrInvalidNilTypeFilter
	}
	handle := &subscription{
		eventBus:  c,
		domain:    domain,
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

	if !c.ValidateDomain(subscription.Domain()) {
		return ErrInvalidUnSubscribeDomain
	}
	c.unsubscribeChan <- subscription
	return nil
}
