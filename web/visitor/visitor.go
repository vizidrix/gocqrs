package visitor

import (
	"github.com/vizidrix/gocqrs/cqrs"
	"fmt"
)

const DOMAIN_NAME string = "github.com/vizidrix/gocqrs/cqrs/web/visitor"
const DOMAIN uint32 = 0xD937B694

const C_V1 uint32 = 0x80010000
const E_V1 uint32 = 0x00010000

const ( /* Commands */
	C_Register 				= C_V1 | 0x0001
	C_HandleRequest			= C_V1 | 0x0002
	C_Blacklist 			= C_V1 | 0x0003
	C_RescindBlacklist		= C_V1 | 0x0004
	C_Whitelist				= C_V1 | 0x0005
	C_RescindWhitelist		= C_V1 | 0x0006
)

const ( /* Events */
	E_Registered			= E_V1 | 0x0001
	E_RequestHandled		= E_V1 | 0x0002
	E_Blacklisted 			= E_V1 | 0x0003
	E_BlacklistRescinded	= E_V1 | 0x0004
	E_Whitelisted 			= E_V1 | 0x0005
	E_WhitelistRescinded	= E_V1 | 0x0006
)

type VisitorMemento struct {
	cqrs.AggregateMemento
	IPV4Address uint32 `json:"ipv4"`
	IPV6Address [2]uint64 `json:"ipv6"`
	Observed bool `json:"observed"`
	Blacklisted bool `json:"blacklisted"`
	Whitelisted bool `json:"whitelisted"`
}

func NewVisitor(id uint64) VisitorMemento {
	return VisitorMemento {
		AggregateMemento: cqrs.NewAggregate(DOMAIN, id, 0),
		IPV4Address: 0,
		IPV6Address: [2]uint64 { 0, 0 },
		Blacklisted: false,
		Whitelisted: false,
	}
}

type Register struct {
	cqrs.CommandMemento
	RequestingUserId uint64 `json:"requestinguserid"`
	IPV4Address uint32 `json:"ipv4"`
	IPV6Address [2]uint64 `json:"ipv6"`
}

func NewRegister(visitorId uint64, requestingUserId uint64, ipv4Address uint32, ipv6Address [2]uint64) Register {
	return Register {
		CommandMemento: cqrs.NewCommand(DOMAIN, C_Register, visitorId, -1),
		RequestingUserId: requestingUserId,
		IPV4Address: ipv4Address,
		IPV6Address: ipv6Address,
	}
}

func NewRegisterIPV4(visitorId uint64, requestingUserId uint64, ipv4Address uint32) Register {
	return NewRegister(visitorId, requestingUserId, ipv4Address, [2]uint64 { 0, 0 })
}

type HandleRequest struct {
	cqrs.CommandMemento
	RequestingUserId uint64 `json:"requestinguserid"`
	Request []byte `json:"request"`
}

func NewHandleRequest(visitorId uint64, request []byte) HandleRequest {
	return HandleRequest {
		CommandMemento: cqrs.NewCommand(DOMAIN, C_HandleRequest, visitorId, -1),
		Request: request,
	}
}

type Blacklist struct {
	cqrs.CommandMemento
	RequestingUserId uint64 `json:"requestintUserId"`
}	

func NewBlacklist(visitorId uint64, requestingUserId uint64) Blacklist {
	return Blacklist {
		CommandMemento: cqrs.NewCommand(DOMAIN, C_Blacklist, visitorId, -1),
		RequestingUserId: requestingUserId,
	}
}

type RescindBlacklist struct {
	cqrs.CommandMemento
	RequestingUserId uint64 `json:"requestintUserId"`
}

func NewRescindBlacklist(visitorId uint64, requestingUserId uint64) RescindBlacklist {
	return RescindBlacklist {
		CommandMemento: cqrs.NewCommand(DOMAIN, C_RescindBlacklist, visitorId, -1),
		RequestingUserId: requestingUserId,
	}
}

type Whitelist struct {
	cqrs.CommandMemento
	RequestingUserId uint64 `json:"requestintUserId"`
}	

func NewWhitelist(visitorId uint64, requestingUserId uint64) Whitelist {
	return Whitelist {
		CommandMemento: cqrs.NewCommand(DOMAIN, C_Whitelist, visitorId, -1),
		RequestingUserId: requestingUserId,
	}
}

type RescindWhitelist struct {
	cqrs.CommandMemento
	RequestingUserId uint64 `json:"requestintUserId"`
}

func NewRescindWhitelist(visitorId, requestingUserId uint64) RescindWhitelist {
	return RescindWhitelist {
		CommandMemento: cqrs.NewCommand(DOMAIN, C_RescindWhitelist, visitorId, -1),
		RequestingUserId: requestingUserId,
	}
}

type Registered struct {
	cqrs.EventMemento
	RequestingUserId uint64 `json:"requestinguserid"`
	IPV4Address uint32 `json:"ipv4"`
	IPV6Address [2]uint64 `json:"ipv6"`
}

func NewRegistered(visitorId uint64, requestingUserId uint64, ipv4Address uint32, ipv6Address [2]uint64) Registered {
	return Registered {
		EventMemento: cqrs.NewEvent(DOMAIN, E_Registered, visitorId, -1),
		RequestingUserId: requestingUserId,
		IPV4Address: ipv4Address,
		IPV6Address: ipv6Address,
	}
}

type RequestHandled struct {
	cqrs.EventMemento
	IPAddress int32 `json:"ipaddress"`
	Request []byte `json:"request"`
}

// TODO: Create id from IP hash
func NewRequestHandled(visitorId uint64, request []byte) RequestHandled {
	return RequestHandled {
		EventMemento: cqrs.NewEvent(DOMAIN, E_RequestHandled, visitorId, -1),
		Request: request,
	}
}

type Blacklisted struct {
	cqrs.EventMemento
	RequestingUserId uint64 `json:"requestinguserid"`
}

func NewBlacklisted(visitorId uint64, requestingUserId uint64) Blacklisted {
	return Blacklisted {
		EventMemento: cqrs.NewEvent(DOMAIN, E_Blacklisted, visitorId, -1),
		RequestingUserId: requestingUserId,
	}
}

type BlacklistRescinded struct {
	cqrs.EventMemento
	RequestingUserId uint64 `json:"requestinguserid"`
}

func NewBlacklistRescinded(visitorId uint64, requestingUserId uint64) BlacklistRescinded {
	return BlacklistRescinded {
		EventMemento: cqrs.NewEvent(DOMAIN, E_BlacklistRescinded, visitorId, -1),
		RequestingUserId: requestingUserId,
	}
}

type Whitelisted struct {
	cqrs.EventMemento
	RequestingUserId uint64 `json:"requestinguserid"`
}

func NewWhitelisted(visitorId uint64, requestingUserId uint64) Whitelisted {
	return Whitelisted {
		EventMemento: cqrs.NewEvent(DOMAIN, E_Whitelisted, visitorId, -1),
		RequestingUserId: requestingUserId,
	}
}

type WhitelistRescinded struct {
	cqrs.EventMemento
	RequestingUserId uint64 `json:"requestinguserid"`
}

func NewWhitelistRescinded(visitorId uint64, requestingUserId uint64) WhitelistRescinded {
	return WhitelistRescinded {
		EventMemento: cqrs.NewEvent(DOMAIN, E_WhitelistRescinded, visitorId, -1),
		RequestingUserId: requestingUserId,
	}
}


func Handle(eventBus chan<-cqrs.Event, es cqrs.EventStorer, command cqrs.Command) {
	switch cmd := command.(type) {
		case Register: {
			fmt.Printf("Trying to register user\n->\t%v\n", command)
			eventBus <- NewRegistered(cmd.GetId(), cmd.RequestingUserId, cmd.IPV4Address, cmd.IPV6Address)
		}
		case HandleRequest: {
			fmt.Printf("Trying to handle request\n->\t%v\n", command)
			eventBus <- NewRequestHandled(cmd.GetId(), cmd.Request)
		}
		case Blacklist: {
			fmt.Printf("Trying to ban visitor!\n->\t%v\n", command)
			// Load aggregate
			// TODO: Change cqrs to es
			events, err := es.ReadAllEvents(command.(cqrs.Aggregate))

			if (len(events) == 0) { // New visitor
				//eventBus <- NewRegistered(command.(cqrs.Aggregate).GetId(), cmd.RequestingUserId)
			}
			
			fmt.Printf("E [ \n\t%v\n ] \nErr [ %s ]\n\n", events, err)
			// Check validation
			// Emit events
			eventBus <- NewBlacklisted(command.(cqrs.Aggregate).GetId(), cmd.RequestingUserId)
			// Poof done
		}
		case RescindBlacklist: {
			fmt.Printf("Trying to lift a visitor ban...\n->\t%v\n", command)
			// Load aggregate
			// TODO: Change cqrs to es
			events, err := es.ReadAllEvents(command.(cqrs.Aggregate))
			
			fmt.Printf("E [ \n\t%v\n ] \tErr [ %s ]\n\n", events, err)
			// Check validation
			// Emit events
			eventBus <- NewBlacklistRescinded(command.(cqrs.Aggregate).GetId(), cmd.RequestingUserId)
			// Poof done
		}
		case Whitelist: {
			fmt.Printf("Trying to whitelist visitor\n->\t%v\n", command)
			eventBus <- NewWhitelisted(cmd.GetId(), cmd.RequestingUserId)
		}
		case RescindWhitelist: {
			fmt.Printf("Trying to rescind whitelist for visitor\n->\t%v\n", command)
			eventBus <- NewWhitelistRescinded(cmd.GetId(), cmd.RequestingUserId)
		}
		default: {
			fmt.Printf("Visitor was unable to handle command: [ %v ]\n", cmd)
		}
	}
	
}













