package cqrs

import (
	"fmt"
)

// http://crc32-checksum.waraxe.us/

const MESSAGE_TYPE_MASK = 0x80000000

func C(domain uint32, version uint64, typeId uint64) uint64 {
	return (uint64(domain) << 32) | uint64(MESSAGE_TYPE_MASK) | (version & 0x7FFF << 16) | (typeId & 0xFFFF)
}

func E(domain uint32, version uint64, typeId uint64) uint64 {
	return (uint64(domain) << 32) | (version & 0x7FFF << 16) | (typeId & 0xFFFF)
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
	GetDomain() uint32
	GetId() uint64
	GetVersion() uint32
	GetCommandType() uint64
}

type CommandMemento struct {
	Domain      uint32 `json:__domain`    // Aggregate Domain
	Id          uint64 `json:"__id"`      // Aggregate Id
	Version     uint32 `json:"__version"` // Aggregate Version
	CommandType uint32 `json:"__ctype"`   // Command Type
}

func NewCommand(domain uint32, id uint64, version uint32, commandType uint64) CommandMemento {
	return CommandMemento{
		Domain:      domain,
		Id:          id,
		Version:     version,
		CommandType: uint32((commandType << 32) >> 32),
	}
}

func (command CommandMemento) GetDomain() uint32 {
	return command.Domain
}

func (command CommandMemento) GetId() uint64 {
	return command.Id
}

func (command CommandMemento) GetVersion() uint32 {
	return command.Version
}

func (command CommandMemento) GetCommandType() uint64 {
	return (uint64(command.Domain) << 32) | uint64(command.CommandType)
}

func (command CommandMemento) String() string {
	return fmt.Sprintf(" <C [ <A D[%d] ID[%d] V[%d] \\> -> C[%d] ] C\\> ",
		command.GetDomain(), command.GetId(), command.GetVersion(), command.CommandType)
}

type Event interface {
	GetDomain() uint32
	GetId() uint64
	GetVersion() uint32
	GetEventType() uint64
}

type EventMemento struct {
	Domain    uint32 `json:__domain`    // Aggregate Domain
	Id        uint64 `json:"__id"`      // Aggregate Id
	Version   uint32 `json:"__version"` // Aggregate Version
	EventType uint32 `json:"__etype"`   // Event Type
}

func NewEvent(domain uint32, id uint64, version uint32, eventType uint64) EventMemento {
	return EventMemento{
		Domain:    domain,
		Id:        id,
		Version:   version,
		EventType: uint32((eventType << 32) >> 32),
	}
}

func (event EventMemento) GetDomain() uint32 {
	return event.Domain
}

func (event EventMemento) GetId() uint64 {
	return event.Id
}

func (event EventMemento) GetVersion() uint32 {
	return event.Version
}

func (event EventMemento) GetEventType() uint64 {
	return (uint64(event.Domain) << 32) | uint64(event.EventType)
}

func (event EventMemento) String() string {
	return fmt.Sprintf(" <E [ <A D[%d] ID[%d] V[%d] \\> -> E[%d] ] E\\> ",
		event.GetDomain(), event.GetId(), event.GetVersion(), event.EventType)
}
