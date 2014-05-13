package web

import (
	"github.com/vizidrix/gocqrs/cqrs"
	"fmt"
)

var DOMAIN int32 = 1

const ( /* Aggregates */ _ = 0 + iota
	Visitor
)

const ( /* Commands */ _ = 0 + iota
	//StartSession
	BanVisitorBySystem
)

const ( /* Events */ _ = 0 + iota
	VisitorRequestReceived
	//SessionStarted
)

type VisitorAggregate struct {
	cqrs.Aggregate
	IPAddress string `json:"ipaddress"`
	IsBanned bool `json:"isbanned"`
}

func NewVisitor(id int64) *VisitorAggregate {
	return &VisitorAggregate {
		Aggregate: cqrs.NewAggregate(DOMAIN, Visitor, id, 0),
		IPAddress: "",
		IsBanned: false,
	}
}

type BanVisitorBySystemCommand struct {
	cqrs.Command
	// TODO: Add requestor details
}

func NewBanVisitorBySystem(visitorId int64) *BanVisitorBySystemCommand {
	return &BanVisitorBySystemCommand {
		Command: cqrs.NewCommand(DOMAIN, Visitor, visitorId, -1, BanVisitorBySystem),
	}
}

type VisitorRequestReceivedEvent struct {
	cqrs.Event
	IPAddress int32 `json:"ipaddress"`
	Request []byte `json:"request"`
}

// TODO: Create id from IP hash
func NewVisitorRequestReceived(visitorId int64, ipAddress int32, request []byte) *VisitorRequestReceivedEvent {
	return &VisitorRequestReceivedEvent {
		Event: cqrs.NewEvent(DOMAIN, Visitor, visitorId, -1, VisitorRequestReceived),
		IPAddress: ipAddress,
		Request: request,
	}
}

func HandleBanVisitorBySystemCommand(eventBus chan<-interface{}, es cqrs.EventStorer, command *BanVisitorBySystemCommand) {
	fmt.Printf("Trying to ban visitor by the system!\t%v\n", command)
	// Load aggregate
	// TODO: Change cqrs to es
	events, err := es.ReadAllEvents(command.Domain, command.Kind, command.Id)
	
	fmt.Printf("E [ %v ] \tErr [ %s ]\n", events, err)
	// Check validation
	// Emit events
	// Poof done
}













