package cqrs

import (
	//"fmt"
	"errors"
)

var (
	ErrInvalidBusState = errors.New("Bus not properly initialized")
	ErrInvalidNilTypeFilter = errors.New("Cannot subscribe with nil type filter")
	ErrInvalidEmptyTypeFilter = errors.New("Cannot subscribe with an empty type filter")
	ErrInvalidNilPublishedEvent = errors.New("Cannot publish a nil event")
	ErrInvalidNilUnSubscribe = errors.New("Cannot unsubscribe a nil subscriber")
	ErrInvalidUnSubscribe = errors.New("Cannot unsubscribe an invalid subscriber")
)

var EventBus EventRouter

func init() {
	// Auto init the event bus as a global service
	EventBus = NewDefaultedEventBus()
	EventBus.Listen()
}

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
	Listen() // Iterates across the Step function in a goroutine loop
	Step()	// Grabs the next operation from the queue and processs it
	Publish(event Event) (error)
	Subscribe(filter EventFilterer) (Subscriber, error)
	UnSubscribe(subscriber Subscriber) (error)
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
	return &eventTypesFilter {
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
	return &aggregateIdsFilter {
		aggregateIds: aggregateIds,
	}
}


type subscription struct {
	//TypeFilter []uint32
	eventBus EventRouter
	filter EventFilterer
	eventChan  chan Event
	//cancelChan chan struct{}
	//EventTypes []uint32
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
	subscribeChan chan Subscriber
	unsubscribeChan chan Subscriber
	publishChan chan Event
	subscriptions []Subscriber
	eventChanFactory EventChanFactory
	cancelChanFactory CancelChanFactory
}

func NewDefaultedEventBus() *channelEventBus {
	return NewChannelEventBus(
		make(chan Subscriber),
		make(chan Subscriber),
		make(chan Event),
		func() chan Event { return make(chan Event)},
		func() chan struct{} { return make(chan struct{})},
		)
}

type EventChanFactory func() chan Event
type CancelChanFactory func() chan struct{}

func NewChannelEventBus(
	subscriptionChan chan Subscriber,
	unsubscriptionChan chan Subscriber,
	publishChan chan Event,
	eventChanFactory EventChanFactory,
	cancelChanFactory CancelChanFactory) *channelEventBus {
	bus := &channelEventBus {
		subscribeChan: subscriptionChan,
		unsubscribeChan: unsubscriptionChan,
		publishChan: publishChan,
		subscriptions: make([]Subscriber, 0, 10),
		eventChanFactory: eventChanFactory,
		cancelChanFactory: cancelChanFactory,
	}
	return bus
}

func (c *channelEventBus) Listen() {
	go func() {
		for {
			c.Step()
		}
	}()
}

func (c *channelEventBus) Step() {
	select { // Synchronized select for event bus mutable actions
	case subscription := <-c.subscribeChan: {
		//fmt.Printf("Subscribe [ %s ]\n", subscription)
		c.subscriptions = append(c.subscriptions, subscription)
	}
	case subscription := <-c.unsubscribeChan: {
		//fmt.Printf("Unsubscribe [ %s ]\n", subscription)
		for index, s := range c.subscriptions {
			if subscription == s {
				c.subscriptions = append(c.subscriptions[:index], c.subscriptions[index+1:]...)
			}
		}
	}
	case event := <- c.publishChan: {
		//fmt.Printf("Publish [ %s ]\n", event)
		for _, subscription := range c.subscriptions {
			if subscription.Filter().Predicate(event) {
				subscription.Publish(event)
				//subscription.EventChan <- event
			}
			/*
			for _, eventType := range subscription.TypeFilter {
				if eventType == event.GetEventType() {
					subscription.EventChan <- event
				}
			}
			*/
		}
	}
	}
}

func (c *channelEventBus) Publish(event Event) (error) {
	if (event == nil) {
		return ErrInvalidNilPublishedEvent
	}
	select {
	case c.publishChan <- event:
	default: 
		return ErrInvalidBusState
	}
	return nil
}

/*
func (c *channelEventBus) Subscribe(eventFilters ...EventPredicate) (*Subscription, error) {
	var compositeFilter := NewCompositeFilter(eventFilters...)

	return 
}
*/

func (c *channelEventBus) Subscribe(filter EventFilterer) (Subscriber, error) {
	//fmt.Printf("Subscribed to ChannelEventBus with\t%x\n", filter)
	if filter == nil {
		return nil, ErrInvalidNilTypeFilter	
	}
	handle := &subscription {
		eventBus: c,
		filter: filter,
		eventChan: c.eventChanFactory(),
		//cancelChan: c.cancelChanFactory(),
	}
	c.subscribeChan <- handle
	return handle, nil
}

func (c *channelEventBus) UnSubscribe(subscription Subscriber) (error) {
	if subscription == nil {
		return ErrInvalidNilUnSubscribe
	}
	//ErrInvalidUnSubscribe
	c.unsubscribeChan <- subscription
	return nil
}

/*
func (c *channelEventBus) Subscribe(typeFilter []uint32) (*Subscription, error) {
	fmt.Printf("Subscribed to ChannelEventBus with\n%#x\n", typeFilter)
	if typeFilter == nil {
		return nil, ErrInvalidNilTypeFilter
	}
	if len(typeFilter) == 0 {
		return nil, ErrInvalidEmptyTypeFilter
	}
	handle := &Subscription {
		TypeFilter: typeFilter,
		EventChan: make(chan Event),
		CancelChan: make(chan struct{}),
	}
	c.subscribeChan <- handle
	return handle, nil
}
*/

/*

type EventSubscriptionService struct {
	Subscriptions map[uint32][]chan Event
}

func NewEventSubscriptionService() EventSubscriptionService {
	return EventSubscriptionService{
		Subscriptions: make(map[uint32][]chan Event),
	}
}



func NewEventSubscription(eventchan chan Event, eventtypes []uint32) EventSubscription {
	return EventSubscription{
		EventChan:  eventchan,
		EventTypes: eventtypes,
	}
}

func (eventsubscriptions *EventSubscriptionService) BusEvents(eventbus <-chan Event, subscriberbus <-chan EventSubscription) {
	for {
		select {
		case event := <-eventbus:
			fmt.Printf("\nEvent Bus: %v", event)
			for _, listener := range eventsubscriptions.Subscriptions[event.GetEventType()] {
				listener <- event
			}
		case subscription := <-subscriberbus:
			for _, eventtype := range subscription.EventTypes {
				eventsubscriptions.Subscriptions[eventtype] = append(eventsubscriptions.Subscriptions[eventtype], subscription.EventChan)
			}
		}
	}
}

*/