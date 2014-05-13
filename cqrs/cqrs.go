package cqrs

import (
	"errors"
	"fmt"
)

// TODO: Handle command & event versioning

type AggregateLoader interface {
	Load(events []Event)
}

type Aggregate interface {
	GetDomain() int32
	GetKind() int32
	GetId() int64
	GetVersion() int32
	MatchById(domain int32, kind int32, id int64) bool
}

type AggregateMemento struct {
	Domain int32 `json:"__domain"`		// Application
	Kind int32 `json:"__kind"`			// Aggregate Kind
	Id int64 `json:"__id"`				// Aggregate Id
	Version int32 `json:"__version"`	// Aggregate Version
}

func NewAggregate(domain int32, kind int32, id int64, version int32) AggregateMemento {
	return AggregateMemento {
		Domain: domain,
		Kind: kind,
		Id: id,
		Version: version,
	}
}

func (aggregate AggregateMemento) GetDomain() int32 {
	return aggregate.Domain
}

func (aggregate AggregateMemento) GetKind() int32 {
	return aggregate.Kind
}

func (aggregate AggregateMemento) GetId() int64 {
	return aggregate.Id
}

func (aggregate AggregateMemento) GetVersion() int32 {
	return aggregate.Version
}

func (aggregate AggregateMemento) MatchById(domain int32, kind int32, id int64) bool {
	return aggregate.Domain == domain && aggregate.Kind == kind && aggregate.Id == id
}

func (aggregate AggregateMemento) String() string {
	return fmt.Sprintf("%d.%d.%d @ %d", aggregate.Domain, aggregate.Kind, aggregate.Id, aggregate.Version)
}

type Command interface {
	GetCommandType() int32
}

type CommandMemento struct {
	AggregateMemento					// Aggregate
	CommandType int32 `json:"__ctype"`	// Command Type
}

func NewCommand(domain int32, kind int32, id int64, version int32, commandType int32) CommandMemento {
	return CommandMemento {
		AggregateMemento: NewAggregate(domain, kind, id, version),
		CommandType: commandType,
	}
}

func (command CommandMemento) GetCommandType() int32 {
	return command.CommandType
}

func (command CommandMemento) String() string {
	return fmt.Sprintf(" <C [ %s -> %d ] C\\> ", command.AggregateMemento.String(), command.CommandType)
}

type Event interface {
	GetEventType() int32
}

type EventMemento struct {
	AggregateMemento					// Aggregate
	EventType int32 `json:"__etype"`	// Event Type
}

func NewEvent(domain int32, kind int32, id int64, version int32, eventType int32) EventMemento {
	return EventMemento {
		AggregateMemento: NewAggregate(domain, kind, id, version),
		EventType: eventType,
	}
}

func (event EventMemento) GetEventType() int32 {
	return event.EventType
}

func (event EventMemento) String() string {
	return fmt.Sprintf(" <E [ %s -> %d ] E\\> ", event.AggregateMemento.String(), event.EventType)
}

type EventStorer interface {
	ReadAllEvents(aggregate Aggregate) ([]interface{}, error)
}

type MemoryEventStore struct {
	Data []interface{}
}

func (es *MemoryEventStore) ReadAllEvents(aggregate Aggregate) ([]interface{}, error) {
	matching := make([]interface{}, 0)
	for _, item := range es.Data {
		switch event := item.(type) {
			case Aggregate: {
				if (event.MatchById(aggregate.GetDomain(), aggregate.GetKind(), aggregate.GetId())) {
					matching = append(matching, item)
				}
				break
			}
			default: {
				return nil, errors.New(fmt.Sprintf("Item in MemoryEventStore isn't an event [ %s ]\n", item))
			}
		}
	}
	return matching, nil
}




