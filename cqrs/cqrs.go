package cqrs

import (
	"errors"
	"fmt"
)

const MESSAGE_TYPE_MASK = 0x80000000
// http://crc32-checksum.waraxe.us/

func C(version uint32, commandId uint32) uint32 {
	return MESSAGE_TYPE_MASK | (version << 16) | (commandId & 0xFF)
}

func E(version uint32, commandId uint32) uint32 {
	return (MESSAGE_TYPE_MASK-1) | (version << 16) | (commandId & 0xFF)
}

type AggregateLoader interface {
	Load(events []Event)
}

type Aggregate interface {
	GetDomain() uint32
	GetId() uint64
	GetVersion() int32
	MatchById(domain uint32, id uint64) bool
}

type AggregateMemento struct {
	Domain uint32 `json:"__domain"`		// Aggregate Domain
	Id uint64 `json:"__id"`				// Aggregate Id
	Version int32 `json:"__version"`	// Aggregate Version
}

func NewAggregate(domain uint32, id uint64, version int32) AggregateMemento {
	return AggregateMemento {
		Domain: domain,
		Id: id,
		Version: version,
	}
}

func (aggregate AggregateMemento) GetDomain() uint32 {
	return aggregate.Domain
}

func (aggregate AggregateMemento) GetId() uint64 {
	return aggregate.Id
}

func (aggregate AggregateMemento) GetVersion() int32 {
	return aggregate.Version
}

func (aggregate AggregateMemento) String() string {
	return fmt.Sprintf("DM[%d] ID[%d] V[%d]", aggregate.Domain, aggregate.Id, aggregate.Version)
}

type Command interface {
	GetCommandType() uint32
}

type CommandMemento struct {
	AggregateMemento					// Aggregate
	CommandType uint32 `json:"__ctype"`	// Command Type
}

func NewCommand(domain uint32, commandType uint32, id uint64, version int32) CommandMemento {
	return CommandMemento {
		AggregateMemento: NewAggregate(domain, id, version),
		CommandType: commandType,
	}
}

func (command CommandMemento) GetCommandType() uint32 {
	return command.CommandType
}

func (command CommandMemento) String() string {
	return fmt.Sprintf(" <C [ %s -> C[%d] ] C\\> ", command.AggregateMemento.String(), command.CommandType)
}

type Event interface {
	GetEventType() uint32
}

type EventMemento struct {
	AggregateMemento					// Aggregate
	EventType uint32 `json:"__etype"`	// Event Type
}

func NewEvent(domain uint32, eventType uint32, id uint64, version int32) EventMemento {
	return EventMemento {
		AggregateMemento: NewAggregate(domain, id, version),
		EventType: eventType,
	}
}

func (event EventMemento) GetEventType() uint32 {
	return event.EventType
}

func (event EventMemento) String() string {
	return fmt.Sprintf(" <E [ %s -> E[%d] ] E\\> ", event.AggregateMemento.String(), event.EventType)
}

type EventStorer interface {
	PersistEvent(event Event)
	ReadAllEvents() (int, []Event, error)
//	ReadAllEventsFrom(index int) (int, []Event, error)
	ReadAggregateEvents(aggregate Aggregate) ([]Event, error)
	ReadAggregateEventsFromSnapshot(aggregate Aggregate) ([]Event, error)
}

type MemoryEventStore struct {
	Snapshot []Aggregate
	EventsChan chan Event
	Data []Event
}

func (eventstore *MemoryEventStore) PersistEvent(event Event) {
	eventstore.Data = append(eventstore.Data, event)
}

func (eventstore *MemoryEventStore) ReadAllEvents() (int, []Event, error) {
	return len(eventstore.Data), eventstore.Data, nil
}

func (eventstore *MemoryEventStore) ReadAllEventsFrom(index int) (int, []Event, error) {
	if index < len(eventstore.Data) {
		events := make([]Event, 0)
		for i := index; i < len(eventstore.Data); i++ {
			events = append(events, eventstore.Data[i])
		}
		return len(eventstore.Data), events, nil
	} else {
		return 0, nil, errors.New("Index position out of range")
	}
}

func (eventstore *MemoryEventStore) ReadAggregateEvents(aggregate Aggregate) ([]Event, error) {
	matching := make([]Event, 0)
	for _, item := range eventstore.Data {
		switch event := item.(type) {
			case Aggregate: {
				if (event.GetDomain() != aggregate.GetDomain() || event.GetId() != aggregate.GetId())) {
					break
				}
				matching = append(matching, item.(Event))
			}
			default: {
				return nil, errors.New(fmt.Sprintf("Item in MemoryEventStore isn't an event [ %s ]\n", item))
			}
		}
	}
	return matching, nil
}

func (eventstore *MemoryEventStore) ReadAggregateEventsFromSnapshot(aggregate Aggregate) ([]Event, error) {
	matching := make([]Event, 0)
	for _, item := range eventstore.Data {
		switch event := item.(type) {
			case Aggregate: {
				if (event.GetDomain() != aggregate.GetDomain() || event.GetId() != aggregate.GetId() || event.GetVersion() < aggregate.GetVersion()) {
					break
				}
				matching = append(matching, item.(Event))
			}
			default: {
				return nil, errors.New(fmt.Sprintf("Item in MemoryEventStore isn't an event [ %s ]\n", item))
			}
		}
	}
	return matching, nil
}

func (eventstore *cqrs.MemoryEventStore) PersistenceLayer(eventBus chan cqrs.Event) {
	for {
		select {
		case event := <- eventstore.EventsChan:
			eventstore.StoreEvent(event)
			eventBus <- event
		}
	}
}








