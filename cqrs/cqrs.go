package cqrs

type AggregateRef struct {
	Domain int32 `json:"__domain"`				// Application
	Kind int32 `json:"__kind"`					// Aggregate Kind
	Id int64 `json:"__id"`						// Aggregate Id
	Version int32 `json:"__version"`			// Aggregate Version
}

type Command struct {
	Domain int32 `json:"__domain"`				// Application
	Kind int32 `json:"__kind"`					// Aggregate Kind
	Id int64 `json:"__id"`						// Aggregate Id
	Version int32 `json:"__version"`			// Aggregate Version
	CommandType int32 `json:"__commandtype"`	// Command Type
}

func NewCommand(domain int32, kind int32, id int64, version int32, commandType int32) Command {
	return Command {
		Domain: domain,
		Kind: kind,
		Id: id,
		Version: version,
		CommandType: commandType,
	}
}

type Event struct {
	Domain int32 `json:"__domain"`			// Application
	Kind int32 `json:"__kind"`				// Aggregate Kind
	Id int64 `json:"__id"`					// Aggregate Id
	Version int32 `json:"__version"`		// Aggregate Version
	EventType int32 `json:"__eventtype"`	// Event Type
}

func NewEvent(domain int32, kind int32, id int64, version int32, eventType int32) Event {
	return Event {
		Domain: domain,
		Kind: kind,
		Id: id,
		Version: version,
		EventType: eventType,
	}
}