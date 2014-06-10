package cqrs

import (
	"fmt"
)

// http://crc32-checksum.waraxe.us/

const MESSAGE_TYPE_MASK = 0x80000000

func C(version uint32, typeId uint32) uint32 {
	return MESSAGE_TYPE_MASK | (version & 0x7FFF << 16) | (typeId & 0xFFFF)

}

func E(version uint32, typeId uint32) uint32 {
	return (version & 0x7FFF << 16) | (typeId & 0xFFFF)
}

type AggregateLoader interface {
	Load(events []Event)
}

type Aggregate interface {
	GetDomain() uint32
	GetId() uint64
	GetVersion() uint32
}

type AggregateMemento struct {
	Domain  uint32 `json:"__domain"`  // Aggregate Domain
	Id      uint64 `json:"__id"`      // Aggregate Id
	Version uint32 `json:"__version"` // Aggregate Version
}

func NewAggregate(domain uint32, id uint64, version uint32) AggregateMemento {
	return AggregateMemento{
		Domain:  domain,
		Id:      id,
		Version: version,
	}
}

func (aggregate AggregateMemento) GetDomain() uint32 {
	return aggregate.Domain
}

func (aggregate AggregateMemento) GetId() uint64 {
	return aggregate.Id
}

func (aggregate AggregateMemento) GetVersion() uint32 {
	return aggregate.Version
}

func (aggregate AggregateMemento) String() string {
	return fmt.Sprintf("<A D[%d] ID[%d] V[%d] \\>", aggregate.Domain, aggregate.Id, aggregate.Version)
}

type Command interface {
	Aggregate
	GetCommandType() uint32
}

type CommandMemento struct {
	AggregateMemento        // Aggregate
	CommandType      uint32 `json:"__ctype"` // Command Type
}

func NewCommand(domain uint32, id uint64, version uint32, commandType uint32) CommandMemento {
	return CommandMemento{
		AggregateMemento: NewAggregate(domain, id, version),
		CommandType:      commandType,
	}
}

func (command CommandMemento) GetCommandType() uint32 {
	return command.CommandType
}

func (command CommandMemento) String() string {
	return fmt.Sprintf(" <C [ %s -> C[%d] ] C\\> ", command.AggregateMemento.String(), command.CommandType)
}

type Event interface {
	Aggregate
	GetEventType() uint32
}

type EventMemento struct {
	AggregateMemento        // Aggregate
	EventType        uint32 `json:"__etype"` // Event Type
}

func NewEvent(domain uint32, id uint64, version uint32, eventType uint32) EventMemento {
	return EventMemento{
		AggregateMemento: NewAggregate(domain, id, version),
		EventType:        eventType,
	}
}

func (event EventMemento) GetEventType() uint32 {
	return event.EventType
}

func (event EventMemento) String() string {
	return fmt.Sprintf(" <E [ %s -> E[%d] ] E\\> ", event.AggregateMemento.String(), event.EventType)
}
