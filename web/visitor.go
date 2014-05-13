package web

import (
	"github.com/vizidrix/gocqrs/cqrs"
	"fmt"
)

var DOMAIN int32 = 1

const ( /* Aggregates */ _ = 0 + iota
	_Visitor
)

const ( /* Commands */ _ = 0 + iota
	//StartSession
	_BanVisitor
	_LiftVisitorBan
)

const ( /* Events */ _ = 0 + iota
	_VisitorRequestReceived
	_VisitorBanned
	_VisitorBanLifted
	//SessionStarted
)

type Visitor struct {
	cqrs.AggregateMemento
	IPAddress int32 `json:"ipaddress"`
	IsBanned bool `json:"isbanned"`
}

func (aggregate Visitor) GetIPAddress() int32 {
	return aggregate.IPAddress
}

func (aggregate Visitor) GetIsBanned() bool {
	return aggregate.IsBanned
}

func NewVisitor(id int64) Visitor {
	return Visitor {
		AggregateMemento: cqrs.NewAggregate(DOMAIN, _Visitor, id, 0),
		IPAddress: 0,
		IsBanned: false,
	}
}

type BanVisitor struct {
	cqrs.CommandMemento
	// TODO: Add requestor details (service acct for system)
}

func NewBanVisitor(visitorId int64) BanVisitor {
	return BanVisitor {
		CommandMemento: cqrs.NewCommand(DOMAIN, _Visitor, visitorId, -1, _BanVisitor),
	}
}

type LiftVisitorBan struct {
	cqrs.CommandMemento
}

func NewLiftVisitorBan(visitorId int64) LiftVisitorBan {
	return LiftVisitorBan {
		CommandMemento: cqrs.NewCommand(DOMAIN, _Visitor, visitorId, -1, _LiftVisitorBan),
	}
}

type VisitorRequestReceived struct {
	cqrs.EventMemento
	IPAddress int32 `json:"ipaddress"`
	Request []byte `json:"request"`
}

// TODO: Create id from IP hash
func NewVisitorRequestReceived(visitorId int64, ipAddress int32, request []byte) VisitorRequestReceived {
	return VisitorRequestReceived {
		EventMemento: cqrs.NewEvent(DOMAIN, _Visitor, visitorId, -1, _VisitorRequestReceived),
		IPAddress: ipAddress,
		Request: request,
	}
}

type VisitorBanned struct {
	cqrs.EventMemento
}

func NewVisitorBanned(visitorId int64) VisitorBanned {
	return VisitorBanned {
		EventMemento: cqrs.NewEvent(DOMAIN, _Visitor, visitorId, -1, _VisitorBanned),
	}
}

type VisitorBanLifted struct {
	cqrs.EventMemento
}

func NewVisitorBanLifted(visitorId int64) VisitorBanLifted {
	return VisitorBanLifted {
		EventMemento: cqrs.NewEvent(DOMAIN, _Visitor, visitorId, -1, _VisitorBanLifted),
	}
}

func Handle(eventBus chan<-interface{}, es cqrs.EventStorer, command cqrs.Command) {
	switch cmd := command.(type) {
		case BanVisitor: {
			fmt.Printf("Trying to ban visitor!\n->\t%v\n", command)
			// Load aggregate
			// TODO: Change cqrs to es
			events, err := es.ReadAllEvents(command.(cqrs.Aggregate))
			
			fmt.Printf("E [ \n\t%v\n ] \nErr [ %s ]\n\n", events, err)
			// Check validation
			// Emit events
			eventBus <- NewVisitorBanned(command.(cqrs.Aggregate).GetId())
			// Poof done
		}
		case LiftVisitorBan: {
			fmt.Printf("Trying to lift a visitor ban...\n->\t%v\n", command)
			// Load aggregate
			// TODO: Change cqrs to es
			events, err := es.ReadAllEvents(command.(cqrs.Aggregate))
			
			fmt.Printf("E [ \n\t%v\n ] \tErr [ %s ]\n\n", events, err)
			// Check validation
			// Emit events
			eventBus <- NewVisitorBanned(command.(cqrs.Aggregate).GetId())
			// Poof done
		}
		default: {
			fmt.Printf("Visitor was unable to handle command: [ %v ]\n", cmd)
		}
	}
	
}













