package web

import (
	"github.com/vizidrix/gocqrs/cqrs"
	"fmt"
)

var DOMAIN_NAME string = "github.com/vizidrix/gocqrs/cqrs/web/visitor"
var DOMAIN int32 = 1

const V1 int64 = 0x00010000

const ( /* Aggregates */ _ = 0 + iota
	A_Visitor
)

const ( /* Commands */ //_ = 0 + iota
	//StartSession
	C_RegisterVisitor 			= V1 & 0x0001
	C_ReportObservedVisitor 	= V1 & 0x0002
	C_BlacklistVisitor 			= V1 & 0x0003
	C_RescindVisitorBlacklist	= V1 & 0x0004
	C_WhitelistVisitor			= V1 & 0x0005
	C_RescindVisitorWhitelist	= V1 & 0x0006
)

const ( /* Events */ _ = 0 + iota
	E_VisitorRegistered
	E_RegisteredVisitorObserved
	E_UnknownVisitorObserved
	E_VisitorBlacklisted
	E_VisitorBlacklistRescinded
	E_VisitorWhitelisted
	E_VisitorWhitelistRescinded
	E_VisitorRequestReceived
	//SessionStarted
)

type Visitor struct {
	cqrs.AggregateMemento
	IPV4Address int32 `json:"ipv4"`
	IPV6Address [2]int64 `json:"ipv6"`
	Blacklisted bool `json:"blacklisted"`
	Whitelisted bool `json:"whitelisted"`
}

func NewVisitor(id int64) Visitor {
	return Visitor {
		AggregateMemento: cqrs.NewAggregate(DOMAIN, _Visitor, id, 0),
		IPV4Address: 0,
		IPV6Address: [2]int64 { 0, 0 },
		Blacklisted: false,
		Whitelisted: false,
	}
}

type BlacklistVisitor struct {
	cqrs.CommandMemento
	// TODO: Add requestor details (service acct for system)
}

func NewBlacklistVisitor(visitorId int64) BanVisitor {
	return BlacklistVisitor {
		CommandMemento: cqrs.NewCommand(DOMAIN, _Visitor, visitorId, -1, _BlacklistVisitor),
	}
}

type RescindVisitorBlacklist struct {
	cqrs.CommandMemento
}

func NewRescindVisitorBlacklist(visitorId int64) RescindVisitorBlacklist {
	return LiftVisitorBan {
		CommandMemento: cqrs.NewCommand(DOMAIN, _Visitor, visitorId, -1, _RescindVisitorBlacklist),
	}
}

type VisitorInitialized struct {
	cqrs.EventMemento
	IPAddress int32 `json:"ipaddress"`
}

func NewVisitorInitialized(visitorId int64, ipAddress int32) VisitorInitialized {
	return VisitorInitialized {
		EventMemento: cqrs.NewEvent(DOMAIN, _Visitor, visitorId, -1, _VisitorInitialized),
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

			if (len(events) == 0) { // New visitor
				eventBus <- NewVisitorInitialized(command.(cqrs.Aggregate).GetId())
			}
			
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













