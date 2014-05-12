
/*
type IHasKind interface {
	GetKind() string
}

type IHasId interface {
	GetId() int64
}

type IHasVersion interface {
	GetVersion() int64
}

type Aggregate struct {
	kind    string `datastore:",noindex" json:"kind"`
	id      int64  `datastore:",noindex" json:"id"`
	version int64  `json:"version"`
}

type IAggregate interface {
	IncrementVersion()
}

func NewAggregate(kind string, id int64) Aggregate {
	return Aggregate{
		kind:    kind,
		id:      id,
		version: 0,
	}
}

func (aggregate *Aggregate) GetKind() string {
	return aggregate.kind
}

func (aggregate *Aggregate) GetId() int64 {
	return aggregate.id
}

func (aggregate *Aggregate) GetVersion() int64 {
	return aggregate.version
}

func (aggregate *Aggregate) IncrementVersion() {
	aggregate.version++
}

type Command struct {
	id      int64 `datastore:",noindex" json:"id"`
	version int64 `json:"version"`
}

type ICommand interface{}

func NewCommand(id int64, version int64) Command {
	return Command{
		id:      id,
		version: version,
	}
}

func (command *Command) GetId() int64 {
	return command.id
}

func (command *Command) GetVersion() int64 {
	return command.version
}

type Event struct {
	id      int64 `datastore:",noindex" json:"id"`
	version int64 `json:"version"`
}

type IEvent interface{}

func NewEvent(id int64, version int64) Event {
	return Event{
		id:      id,
		version: version,
	}
}

func (event *Event) GetId() int64 {
	return event.id
}

func (event *Event) GetVersion() int64 {
	return event.version
}

type AggregateLoader func(interface{}, <-chan IEvent) (interface{}, error)
type EventHandler func(IEvent)
type CommandHandler func(commandSub <-chan ICommand, eventPub chan<- IEvent) chan error

*/

/*
type IEventStore interface {
	func GetEvents(kind string, id int64) <-chan IEvent
}
*/
/*
// Given an aggregate kind this func returns a writer for all related events
type EventWriteStore func(kind string)

// Given an event this func puts the event in the appropriate store
type EventWriter func(event *Event) error

// Given an aggregate kind this func returns a retreiver for all related events
type EventReadStore func(kind string) EventSource

// Given an id and optional min version this func produces a read chan of events
type EventSource func(id int64, version int) <-chan *Event

type AggregateKind struct {
	Kind string
}

type AggregateStore struct {
	aggregateKeyMap map[int64]*EventSet
}

type EventSet struct {
	events []*Event
}

type AggregateKindMapper interface {
	Of(kind string) AggregateIdMapper
}

type AggregateIdMapper interface {
	WithId(id int64) EventSet
}

//func (eventStore *MemoryEventStore) Of(kind string) *AggregateStore {
func (eventStore *MemoryEventStore) Of(kind string) AggregateKindMapper {
	// Do a quick check to ensure that the inner map exsits
	if _, ok := eventStore.aggregateKindMap[kind]; !ok {
		eventStore.aggregateKindMap[kind] = &AggregateStore{
			aggregateKeyMap: make(map[int64]*EventSet),
		}
	}
	return eventStore.aggregateKindMap[kind]
	/
		return func(id int64) EventWriter {
			if _, ok :=
			return func(event *Event) error {
				if aggregateMap, ok := eventStore[kind]; ok {

				} else {
					return errors.New("text")
				}
			}
			return nil
		}
	/
	//return func(event *Event) error {
}

func (aggregateStore *AggregateStore) WithId(id int64) *EventSet {
	if _, ok := aggregateStore.aggregateKeyMap[id]; !ok {
		aggregateStore.aggregateKeyMap[id] = &EventSet{
			events: make([]*Event, 1),
		}
	}
	return aggregateStore.aggregateKeyMap[id]
}
*/