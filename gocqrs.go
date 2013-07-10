package gocqrs

type Aggregate struct {
	Id int64
}

type IAggregate interface{}

type Command struct {
	Id int64
}

type ICommand interface{}

type CommandHandler func(ICommand)

type Event struct {
	Id int64
}

type IEvent interface{}

type EventHandler func(IEvent)

type IEventStore interface {
}
