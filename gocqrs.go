package gocqrs

import (
	"errors"
	"time"
)

var (
	// TODO: Bring these back?
	//ErrUsedTimestamp  = errors.New("timestamp used")
	//ErrUsedKey        = errors.New("datastore key used")

	// ErrInvalidApplication is used to inform a consumer when they've
	// provided an aggregate that doesn't have a valid application id
	// that the receiving service is able to process
	ErrInvalidApplication = errors.New("invalid application identifier")

	// ErrInvalidDomain is used to inform a consumer when they've
	// provided an aggregate that doesn't have a valid domain id that
	// the receiving service is able to process
	// * Domain is semantically equal to Aggregate Type
	ErrInvalidDomain = errors.New("invalid domain identifier")

	// ErrInvalidAggregateId is used to inform a consumer when they've
	// provided an aggregate id that is not available due to either
	// overlap with an existing aggregate or domain specific command
	// handler rules
	ErrInvalidAggregateId = errors.New("invalid aggregate identifier")

	// ErrInvalidVersion is used to inform a consumer when they've
	// provided an aggregate with a version that cannot be sync'd
	// with the current domain version
	ErrInvalidVersion = errors.New("invalid aggregate version")

	// ErrInvalidCommandType is used to inform a consumer when they've
	// provided a command type that isn't valid for the application and
	// domain partition
	ErrInvalidCommandType = errors.New("invalid command type identifier")

	// ErrInvalidEventType is used to inform a consumer when they've
	// provided an event type that isn't valid for the application and
	// domain partition
	ErrInvalidEventType = errors.New("invalid event type identifier")

	// ErrUnableToFindAggregate is used to inform a consumer when the
	// aggregate associate with a command wasn't found in the store
	ErrUnableToFindAggregate = errors.New("unable to locate specified aggregate")

	// ErrUnableToLoadAggregate is used to inform a consumer when the
	// aggregate loaded form the store failed to hydrate properly
	ErrUnableToLoadAggregate = errors.New("error occured loading aggregate")

	// ErrErrorApplyingCommand is used to inform a consumer when the
	// command handler returns an errory when applying the command
	// to the target aggregate
	ErrErrorApplyingCommand = errors.New("error occured applying command")

	// ErrErrorAppendingEvent is used to inform a consumer when there
	// is an error appending the event to the eventstore
	ErrErrorAppendingEvent = errors.New("error writing to the eventstore")

	// ErrErrorPublishingEvent is used to inform a consumer when there
	// is an error publishing the event produced by the command handler
	// This step occurs after the event has been stored
	ErrErrorPublishingEvent = errors.New("error publishing the event")
)

type TypeBuilder func(uint8, uint32) uint32

// MakeVersionedCommandType provides a utility to union a command's version and
// type identifiers and masks off the leftmost bit as 1 to indicate a command
func MakeVersionedCommandType(version uint8, typeId uint32) uint32 {
	return 0x80000000 | (uint32(version) << 24 & 0x7F000000) | (typeId & 0xFFFFFF)
}

// MakeVersionedEventType provides a utility to union an event's version and
// type identifiers and masks off the leftmost bit as 0 to indicate an event
func MakeVersionedEventType(version uint8, typeId uint32) uint32 {
	return 0x7FFFFFFF&(uint32(version)<<24&0x7F000000) | (typeId & 0xFFFFFF)
}

// EventStoreReaderWriter describes a type the can be used to either read
// or write events to an eventstore
type EventStoreReaderWriterGenerator interface {
	AggregateIdGenerater
	EventStoreWriter
	EventStoreReader
}

// AggregateIdGenerator is responsible for creating valid unique Ids for Aggregates
type AggregateIdGenerater interface {
	//GenerateAggregateId(application uint32, domain uint32) (uint64, error)
	GenerateAggregateId() (uint64, error)
}

// EventWriter is responsible for persisting Events to the EventStore
type EventStoreWriter interface {
	AppendEvent(Event) (time.Time, error)
}

// Responsible for serving Streams as queries against the EventStore
type EventStoreReader interface {
	LoadEvents() ([]Event, error)
	LoadEventsByAggregate(aggregate uint64) ([]Event, error)
	LoadEventsByEventType(eventType uint32) ([]Event, error)
	LoadEventsByEventTypes(eventTypes ...uint32) ([]Event, error)
}

// Aggregate provides a base interface for things that contain
// aggregate header information
type Aggregate interface {
	GetApplication() uint32
	GetDomain() uint32
	GetId() uint64
	GetVersion() uint32
}

// AggregateHydrator describes a type which processes a slice of events to produce
// a populated aggregate instance
type AggregateHydrator interface {
	LoadAggregate([]event) (Aggregate, error)
}

// Command provides a base interface for all commands in the
// system which includes aggregate header information to identity
// the target of the command
type Command interface {
	Aggregate
	GetCommandType() uint32
}

// CommandHandler describes a type that can be used to process commands
type CommandHandler interface {
	Handle(command Command) error
}

// CommandSerializerDeSerializer  describes a type that can be used to
// either serialize or deserialize a Command to/from a byte slice
type CommandSerializerDeserializer interface {
	CommandSerializer
	CommandDeserializer
}

// CommandSerializer describes a type that can be used to serialize
// Commands to a raw byte slice
type CommandSerializer interface {
	Serialize(Command) ([]byte, error)
}

// CommandDeserializer describes a type that can be used to deserialize
// Commands from a raw byte slice
type CommandDeserializer interface {
	Deserialize([]byte) (Command, error)
}

// TypedCommandSerializerDeserializer describes a type that can be used to serialize
// or deserialize Cp,,amds from a raw byte slice given the commandType
type TypedCommandSerializerDeserializer interface {
	CommandSerializer
	TypedCommandDeserializer
}

// TypedCommandDeserializer describes a type that can be used to deserialize
// Command from a raw byte slice given the commandType
type TypedCommandDeserializer interface {
	Deserialize(uint32, []byte) (Command, error)
}

// Event provides a base interface for all events in the system
// which includes aggregate header information to identify the
// target of the event
type Event interface {
	Aggregate
	GetEventType() uint32
}

// EventPublisher describes a type that can be used to publish events to a bus
type EventPublisher interface {
	Publish(Event) error
}

// EventHandler describes a type that can be used to process events
type EventHandler interface {
	Handle(event Event) (time.Time, error)
}

// EventSerializerDeSerializer  describes a type that can be used to
// either serialize or deserialize an Event to/from a byte slice
type EventSerializerDeserializer interface {
	EventSerializer
	EventDeserializer
}

// EventSerializer describes a type that can be used to serialize
// Events to a raw byte slice
type EventSerializer interface {
	Serialize(Event) ([]byte, error)
}

// EventDeserializer describes a type that can be used to deserialize
// Events from a raw byte slice
type EventDeserializer interface {
	Deserialize([]byte) (Event, error)
}

// TypedEventSerializerDeserializer describes a type that can be used to serialize
// or deserialize Events from a raw byte slice given the eventType
type TypedEventSerializerDeserializer interface {
	EventSerializer
	TypedEventDeserializer
}

// TypedEventDeserializer describes a type that can be used to deserialize
// Events from a raw byte slice given the eventType
type TypedEventDeserializer interface {
	Deserialize(uint32, []byte) (Event, error)
}

// aggregate is a structured header describing the UUId of an aggregate instance
type aggregate struct {
	// application the target aggregate belongs to, provides multi-tenancy
	// at the application level partition for like domains within the same service
	Application uint32 `json:"_app"`
	// domain is the type of aggregate (type is semantically equivalent to doman)
	Domain uint32 `json:"_domain"`
	// id is an [application / domain] unique identifier for the aggregate instance
	// and should never be duplicated within that partition
	Id uint64 `json:"_id"`
	// version is derived from the number of events applied to the aggregate
	// and provides guaranteed event ordering within it's
	// [appliction / domain / id] partition
	Version uint32 `json:"_ver"`
}

// NewAggregate creates an aggregate instance with UUId derived from the provided values
func NewAggregate(application uint32, domain uint32, id uint64, version uint32) Aggregate {
	return aggregate{
		Application: application,
		Domain:      domain,
		Id:          id,
		Version:     version,
	}
}

// GetApplication returns the application id this aggregate
// was designed within
func (aggregate aggregate) GetApplication() uint32 {
	return aggregate.Application
}

// GetDomain returns the domain (or aggregate type) of this aggregate
func (aggregate aggregate) GetDomain() uint32 {
	return aggregate.Domain
}

// GetId returns the id of the aggregate which is unique within the
// partition provided by the combination of application and domain
func (aggregate aggregate) GetId() uint64 {
	return aggregate.Id
}

// GetVersion returns the version of the aggregate represented by
// this aggregate instance.  Not guaranteed to be the current version
// just the version state of the aggregate when this instance was
// loaded
func (aggregate aggregate) GetVersion() uint32 {
	return aggregate.Version
}

// command is a structured header describing the UUID of a Command instance
type command struct {
	// aggregate is the base structure that binds the command instance
	// to the target aggregate by capturing the aggregate's full UUId
	// partition information [ application / domain / id / version ]
	aggregate
	// commandType is an [ application / domain ] unique identifier for the type of
	// command message which captures the semantic intent of the command
	CommandType uint32 `json:"_ctype"`
}

// NewCommand creates a command instance with UUID derived from the provided values
// including the header of the targeted aggregate instance
func NewCommand(application uint32, domain uint32, id uint64, version uint32, commandType uint32) Command {
	return command{
		aggregate: aggregate{
			Application: application,
			Domain:      domain,
			Id:          id,
			Version:     version,
		},
		CommandType: commandType,
	}
}

// GetCommandType returns the command type of the event that is unique within
// the [ application / domain ] partition
func (command command) GetCommandType() uint32 {
	return command.CommandType
}

// event is a structured header describing the UUID of an Event instance
type event struct {
	// aggregate is the base structure that binds the event instance
	// to the target aggregate by capturing the aggregate's full UUId
	// partition information [ application / domain / id / version ]
	aggregate
	// eventType is an [ application / domain ] unique identifier for the type of
	// event message which captures the semantic intent of the event
	EventType uint32 `json:"_etype"`
}

// NewEvent creates an event instance with UUID derived from the provided values
// including the header of the targeted aggregate instance
func NewEvent(application uint32, domain uint32, id uint64, version uint32, eventType uint32) Event {
	return event{
		aggregate: aggregate{
			Application: application,
			Domain:      domain,
			Id:          id,
			Version:     version,
		},
		EventType: eventType,
	}
}

// GetEventType returns the event type of the event that is unique within
// the [ application / domain ] partition
func (event event) GetEventType() uint32 {
	return event.EventType
}

// AggregateLoader describes a function which takes a slice of events and
// produces either a valid aggregate or an error
type AggregateLoader func([]Event) (Aggregate, error)

// CommandEvaluator describes a function which evaluates a
type CommandEvaluator func(AggregateIdGenerater, Aggregate, Command) (Event, error)

// DefaultCommandHandler provides a base implementation for domain specific command
// handlers to use if they follow a standard execution path
func DefaultCommandHandler(eventStore EventStoreReaderWriterGenerator, publisher EventPublisher, loader AggregateLoader, evaluator CommandEvaluator, command Command) (err error) {
	// Read the events from the store
	events, err := eventStore.LoadEventsByAggregate(command.GetId())
	if err != nil {
		return ErrUnableToFindAggregate
	}
	// Populate an aggregate using the retrieved events
	aggregate, err := loader(events)
	if err != nil {
		return ErrUnableToLoadAggregate
	}
	// Evaluate the command against the aggregate
	event, err := evaluator(eventStore, aggregate, command)
	if err != nil {
		return ErrErrorApplyingCommand
	}
	// Commit the event to the eventstore
	_, err = eventStore.AppendEvent(event)
	if err != nil {
		return ErrErrorAppendingEvent
	}
	// Broadcast the created event to all observers
	err = publisher.Publish(event)
	if err != nil {
		return ErrErrorPublishingEvent
	}
	return err
}
