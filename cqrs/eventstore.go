package cqrs

import (
	"errors"
	"fmt"
)

type EventStorer interface {
	StoreEvent(event Event)
	ReadAllEvents() (int, []Event, error)
	ReadAggregateEvents(domain uint32, id uint64) ([]Event, error)
	ReadAggregateEventsFromSnapshot(domain uint32, id uint64, version int32) ([]Event, error)
}

type MemoryEventStore struct {
	Snapshots []Aggregate
	Data      []Event
}

func NewMemoryEventStore() MemoryEventStore {
	return MemoryEventStore{
		Snapshots: make([]Aggregate, 0),
		Data:      make([]Event, 0),
	}
}

func (eventstore *MemoryEventStore) StoreEvent(event Event) {
	eventstore.Data = append(eventstore.Data, event)
}

func (eventstore *MemoryEventStore) ReadAllEvents() (int, []Event, error) {
	return len(eventstore.Data), eventstore.Data, nil
}

func (eventstore *MemoryEventStore) ReadAggregateEvents(domain uint32, id uint64) ([]Event, error) {
	matching := make([]Event, 0)
	for _, item := range eventstore.Data {
		switch event := item.(type) {
		case Aggregate:
			{
				if event.GetDomain() != domain || event.GetId() != id {
					break
				}
				matching = append(matching, item.(Event))
			}
		default:
			{
				return nil, errors.New(fmt.Sprintf("Item in MemoryEventStore isn't an event [ %s ]\n", item))
			}
		}
	}
	return matching, nil
}

func (eventstore *MemoryEventStore) ReadAggregateEventsFromSnapshot(domain uint32, id uint64, version uint32) ([]Event, error) {
	matching := make([]Event, 0)
	for _, item := range eventstore.Data {
		switch event := item.(type) {
		case Aggregate:
			{
				if event.GetDomain() != domain || event.GetId() != id || event.GetVersion() < version {
					break
				}
				matching = append(matching, item.(Event))
			}
		default:
			{
				return nil, errors.New(fmt.Sprintf("Item in MemoryEventStore isn't an event [ %s ]\n", item))
			}
		}
	}
	return matching, nil
}
/*
type EventStorer interface {
	StoreEvent(event Event)
	ReadAllEvents() (int, []Event, error)
	ReadAggregateEvents(domain uint32, id uint64) ([]Event, error)
	ReadAggregateEventsFromSnapshot(domain uint32, id uint64, version int32) ([]Event, error)
}

type MemoryEventStore struct {
	Snapshots []Aggregate
	Data      []Event
}

func NewMemoryEventStore() MemoryEventStore {
	return MemoryEventStore{
		Snapshots: make([]Aggregate, 0),
		Data:      make([]Event, 0),
	}
}

func (eventstore *MemoryEventStore) StoreEvent(event Event) {
	eventstore.Data = append(eventstore.Data, event)
}

func (eventstore *MemoryEventStore) ReadAllEvents() (int, []Event, error) {
	return len(eventstore.Data), eventstore.Data, nil
}

func (eventstore *MemoryEventStore) ReadAggregateEvents(domain uint32, id uint64) ([]Event, error) {
	matching := make([]Event, 0)
	for _, item := range eventstore.Data {
		switch event := item.(type) {
		case Aggregate:
			{
				if event.GetDomain() != domain || event.GetId() != id {
					break
				}
				matching = append(matching, item.(Event))
			}
		default:
			{
				return nil, errors.New(fmt.Sprintf("Item in MemoryEventStore isn't an event [ %s ]\n", item))
			}
		}
	}
	return matching, nil
}

func (eventstore *MemoryEventStore) ReadAggregateEventsFromSnapshot(domain uint32, id uint64, version int32) ([]Event, error) {
	matching := make([]Event, 0)
	for _, item := range eventstore.Data {
		switch event := item.(type) {
		case Aggregate:
			{
				if event.GetDomain() != domain || event.GetId() != id || event.GetVersion() < version {
					break
				}
				matching = append(matching, item.(Event))
			}
		default:
			{
				return nil, errors.New(fmt.Sprintf("Item in MemoryEventStore isn't an event [ %s ]\n", item))
			}
		}
	}
	return matching, nil
}
*/