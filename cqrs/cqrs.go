package cqrs

import (
	"errors"
)

// TODO: Handle command & event versioning

type AggregateLoader interface {
	Load(events []Event)
}

type Aggregate struct {
	Domain int32 `json:"__domain"`		// Application
	Kind int32 `json:"__kind"`			// Aggregate Kind
	Id int64 `json:"__id"`				// Aggregate Id
	Version int32 `json:"__version"`	// Aggregate Version
}

func NewAggregate(domain int32, kind int32, id int64, version int32) Aggregate {
	return Aggregate {
		Domain: domain,
		Kind: kind,
		Id: id,
		Version: version,
	}
}

/*
type Message struct {
	CorrelationId int64
}
*/

type Command struct {
	Aggregate							// Aggregate
	CommandType int32 `json:"__ctype"`	// Command Type
}

func NewCommand(domain int32, kind int32, id int64, version int32, commandType int32) Command {
	return Command {
		Aggregate: NewAggregate(domain, kind, id, version),
		CommandType: commandType,
	}
}

type Event struct {
	Aggregate							// Aggregate
	EventType int32 `json:"__etype"`	// Event Type
}

func NewEvent(domain int32, kind int32, id int64, version int32, eventType int32) Event {
	return Event {
		Aggregate: NewAggregate(domain, kind, id, version),
		EventType: eventType,
	}
}

type EventStorer interface {
	ReadAllEvents(domain int32, kind int32, id int64) ([]interface{}, error)
}

type MemoryEventStore struct {
	Data []interface{}
}

func (es *MemoryEventStore) ReadAllEvents(domain int32, kind int32, id int64) ([]interface{}, error) {
	return nil, errors.New("ES Not Configured")
}