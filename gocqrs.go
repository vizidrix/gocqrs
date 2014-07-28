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

/*
TODO: Re-add or Re-move the following...
type CommandSerialConverter interface {
	CommandSerializer
	CommandDeserializer
}

type CommandSerializer interface {
	SerializeCommand(Command) ([]byte, error)
}

type CommandDeserializer interface {
	DeserializeCommand([]byte) (Command, error)
}

func GetCommandTypeFromJson(jsonCommand []byte) (commandType uint64, err error) {
	var command CommandMemento
	err = json.Unmarshal(jsonCommand, &command)
	return command.GetCommandType(), err
}

*/

// MakeVersionedCommandType provides a utility to union a command's version and
// type identifiers and masks off the leftmost bit as 1 to indicate a command
func MakeVersionedCommandType(version uint8, typeId uint32) uint32 {
	return 0x80000000 | (uint32(version) << 24  & 0x7F000000) | (typeId & 0xFFFFFF)
}

// MakeVersionedEventType provides a utility to union an event's version and
// type identifiers and masks off the leftmost bit as 0 to indicate an event
func MakeVersionedEventType(version uint8, typeId uint32) uint32 {
	return 0x7FFFFFFF & (uint32(version) << 24  & 0x7F000000) | (typeId & 0xFFFFFF)
}

// EventStoreReaderWriter describes a type the can be used to either read
// or write events to an eventstore
type EventStoreReaderWriter interface {
	AggregateIdGenerater
	EventStoreWriter
	EventStoreReader
}

// AggregateIdGenerator is esponsible for creating valid unique Ids for Aggregates
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
/*
type EventStoreReader interface {
	LoadEventsByAggregate(application uint32, domain uint32, aggregate uint64) ([]Event, error)
	LoadEventsByEventType(application uint32, domain uint32, eventType uint32) ([]Event, error)
	LoadEventsByMultipleEventTypes(application uint32, domain uint32, eventTypes ...uint32) ([]Event, error)
	LoadEventsByDomain(application uint32, domain uint32) ([]Event, error)
}
*/

// Aggregate provides a base interface for things that contain
// aggregate header information
type Aggregate interface {
	GetApplication() uint32
	GetDomain() uint32
	GetId() uint64
	GetVersion() uint32
}

// Command provides a base interface for all commands in the
// system which includes aggregate header information to identity
// the target of the command
type Command interface {
	Aggregate
	GetCommandType() uint32
}

// AggregateHydrator describes a method which processes a slice of events to produce
// a populated aggregate instance
type HydrateAggregate func([]Event)(Aggregate, error)

// ApplyCommand describes a method which, given an aggregate instance and a command,
// attempts to apply the command's intent then reports back success/fail event
type ApplyCommand func(aggregate Aggregate, command Command)(Event, error)

// PublishEvent describes a method which broadcasts an event to a pub/sub
type PublishEvent func(event Event)(error)

// CommandHandler describes a type that can be used to process commands
type CommandHandler interface {
	Handle(command Command) (error)
}

// Event provides a base interface for all events in the system
// which includes aggregate header information to identify the
// target of the event
type Event interface {
	Aggregate
	GetEventType() uint32
}

// EventHandler describes a type that can be used to process events
type EventHandler interface {
	Handle(event Event) (time.Time, error)
}

// EventSerializerDeSerializer  describes a type that can be used to
// either serialize or deserialize an Event to/from a byte slice
type EventSerializerDeSerializer interface {
	EventSerializer
	EventDeserializer
}

// EventSerializer describes a type that can be used to serialize
// Events to a raw byte slice
type EventSerializer interface {
	SerializeEvent(Event) ([]byte, error)
}

// EventDeserializer describes a type that can be used to deserialize
// Events from a raw byte slice
type EventDeserializer interface {
	DeserializeEvent([]byte) (Event, error)
}

type InformedSerialConverter interface {
	EventSerializer
	InformedDeserializer
}

type InformedDeserializer interface {
	InformedDeserializeEvent(uint32, []byte) (Event, error)
}

// aggregate is a structured header describing the UUId of an aggregate instance
type aggregate struct {
	// application is the target aggregate belongs to, provides multi-tenancy
	// at the application level partition for like domains within the same service
	application uint32 `json:"_app"`
	// domain is type of aggregate (type is semantically equivalent to doman)
	domain      uint32 `json:"_domain"`
	// id is an [application / domain] unique identifier for the aggregate instance
	// and should never be duplicated within that partition
	id          uint64 `json:"_id"`
	// version is derived from the number of events applied to the aggregate
	// and provides guaranteed event ordering within it's
	// [appliction / domain / id] partition
	version     uint32 `json:"_ver"`
}

// NewAggregate creates an aggregate instance with UUId derived from the provided values
func NewAggregate(application uint32, domain uint32, id uint64, version uint32) Aggregate {
	return &aggregate{
		application: application,
		domain:      domain,
		id:          id,
		version:     version,
	}
}

// GetApplication returns the application id this aggregate
// was designed within
func (aggregate *aggregate) GetApplication() uint32 {
	return aggregate.application
}

// GetDomain returns the domain (or aggregate type) of this aggregate
func (aggregate *aggregate) GetDomain() uint32 {
	return aggregate.domain
}

// GetId returns the id of the aggregate which is unique within the
// partition provided by the combination of application and domain
func (aggregate *aggregate) GetId() uint64 {
	return aggregate.id
}

// GetVersion returns the version of the aggregate represented by
// this aggregate instance.  Not guaranteed to be the current version
// just the version state of the aggregate when this instance was
// loaded
func (aggregate *aggregate) GetVersion() uint32 {
	return aggregate.version
}

// command is a structured header describing the UUID of a Command instance
type command struct {
	// aggregate is the base structure that binds the command instance
	// to the target aggregate by capturing the aggregate's full UUId
	// partition information [ application / domain / id / version ]
	aggregate
	// commandType is an [ application / domain ] unique identifier for the type of 
	// command message which captures the semantic intent of the command
	commandType   uint32 `json:"_ctype"`
}

// NewCommand creates a command instance with UUID derived from the provided values
// including the header of the targeted aggregate instance
func NewCommand(application uint32, domain uint32, id uint64, version uint32, commandType uint32) Command {
	return &command {
		aggregate: aggregate {
			application: application,
			domain: domain,
			id: id,
			version: version,
			},
		commandType: commandType,
	}
}

// GetCommandType returns the command type of the event that is unique within
// the [ application / domain ] partition
func (command *command) GetCommandType() uint32 {
	return command.commandType
}

// commandHandler is a container which wraps all the dependencies needed for
// the general command handler case to perform all of it's related duties
type commandHandler struct {
	reader 		EventStoreReader
	writer 		EventStoreWriter
	hydrator 	HydrateAggregate
	applicator 	ApplyCommand
	publisher 	PublishEvent
	application uint32
	domain 		uint32
}

// NewCommandHandler provides a strongly typed list of dependencies needed to boot
// a generalized command handler
// Custom initializers could be built over this format to reduce the parameter list
func NewCommandHandler(reader EventStoreReader, writer EventStoreWriter, hydrator HydrateAggregate, applicator ApplyCommand, publisher PublishEvent, application uint32, domain uint32) (CommandHandler) {
	return &commandHandler {
		reader: reader,
		writer: writer,
		hydrator: hydrator,
		applicator: applicator,
		publisher: publisher,
		application: application,
		domain: domain,
	}
}

// Handle processes the provided command throgh the necessary steps to validate
// the command, load the events from the store, populate the aggregate, apply
// the command, append the resulting event and, finally, publish the event
// TODO: Create an async Handle option
func (c *commandHandler) Handle(command Command) (error) {
	// Validate the partition of the command
	if c.application != command.GetApplication() {
		return ErrInvalidApplication
	}
	if c.domain != command.GetDomain() {
		return ErrInvalidDomain
	}
	// Read the events from the store
	events, err := c.reader.LoadEventsByAggregate(command.GetId())
	if err != nil {
		return ErrUnableToFindAggregate
	}
	// Populate an aggregate using the retrieved events
	aggregate, err := c.hydrator(events)
	if err != nil {
		return ErrUnableToLoadAggregate
	}
	// Evaluate the command against the aggregate
	event, err := c.applicator(aggregate, command)
	if err != nil {
		return ErrErrorApplyingCommand
	}
	// Commit the event to the eventstore
	_, err = c.writer.AppendEvent(event)
	if err != nil {
		return ErrErrorAppendingEvent
	}
	// Broadcast the created event to all observers
	err = c.publisher(event)
	if err != nil {
		return ErrErrorPublishingEvent
	}
	return err
}

// event is a structured header describing the UUID of an Event instance
type event struct {
	// aggregate is the base structure that binds the event instance
	// to the target aggregate by capturing the aggregate's full UUId
	// partition information [ application / domain / id / version ]
	aggregate
	// eventType is an [ application / domain ] unique identifier for the type of 
	// event message which captures the semantic intent of the event
	eventType   uint32 `json:"_etype"`
}

// NewEvent creates an event instance with UUID derived from the provided values
// including the header of the targeted aggregate instance
func NewEvent(application uint32, domain uint32, id uint64, version uint32, eventType uint32) Event {
	return &event{
		aggregate: aggregate {
			application: application,
			domain: domain,
			id: id,
			version: version,
			},
		eventType: eventType,
	}
}

// GetEventType returns the event type of the event that is unique within
// the [ application / domain ] partition
func (event *event) GetEventType() uint32 {
	return event.eventType
}