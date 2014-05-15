package visitor

import (
	"fmt"
	"github.com/vizidrix/gocqrs/cqrs"
)

var DOMAIN_NAME string = "github.com/vizidrix/gocqrs/cqrs/web/visitor"
var DOMAIN uint32 = 0xD937B694

var ( /* Commands */
	C_Register         = cqrs.C(1, 1)
	C_HandleRequest    = cqrs.C(1, 2)
	C_Blacklist        = cqrs.C(1, 3)
	C_RescindBlacklist = cqrs.C(1, 4)
	C_Whitelist        = cqrs.C(1, 5)
	C_RescindWhitelist = cqrs.C(1, 6)
)

var ( /* Events */
	E_Registered         = cqrs.E(1, 1)
	E_RequestHandled     = cqrs.E(1, 2)
	E_Blacklisted        = cqrs.E(1, 3)
	E_BlacklistRescinded = cqrs.E(1, 4)
	E_Whitelisted        = cqrs.E(1, 5)
	E_WhitelistRescinded = cqrs.E(1, 6)
)

type VisitorMemento struct {
	cqrs.AggregateMemento
	IPV4Address uint32    `json:"ipv4"`
	IPV6Address [2]uint64 `json:"ipv6"`
	Observed    bool      `json:"observed"`
	Blacklisted bool      `json:"blacklisted"`
	Whitelisted bool      `json:"whitelisted"`
}

func NewVisitor(id uint64) VisitorMemento {
	return VisitorMemento{
		AggregateMemento: cqrs.NewAggregate(DOMAIN, id, 0),
		IPV4Address:      0,
		IPV6Address:      [2]uint64{0, 0},
		Blacklisted:      false,
		Whitelisted:      false,
	}
}

type Register struct {
	cqrs.CommandMemento
	RequesterId uint64    `json:"requestinguserid"`
	IPV4Address uint32    `json:"ipv4"`
	IPV6Address [2]uint64 `json:"ipv6"`
}

func NewRegister(visitorId uint64, requestingUserId uint64, ipv4Address uint32, ipv6Address [2]uint64) Register {
	return Register{
		CommandMemento: cqrs.NewCommand(DOMAIN, C_Register, visitorId, -1),
		RequesterId:    requestingUserId,
		IPV4Address:    ipv4Address,
		IPV6Address:    ipv6Address,
	}
}

func NewRegisterIPV4(visitorId uint64, requestingUserId uint64, ipv4Address uint32) Register {
	return NewRegister(visitorId, requestingUserId, ipv4Address, [2]uint64{0, 0})
}

type HandleRequest struct {
	cqrs.CommandMemento
	RequesterId uint64 `json:"requestinguserid"`
	Request     []byte `json:"request"`
}

func NewHandleRequest(visitorId uint64, request []byte) HandleRequest {
	return HandleRequest{
		CommandMemento: cqrs.NewCommand(DOMAIN, C_HandleRequest, visitorId, -1),
		Request:        request,
	}
}

type Blacklist struct {
	cqrs.CommandMemento
	RequesterId uint64 `json:"requestintUserId"`
}

func NewBlacklist(visitorId uint64, requestingUserId uint64) Blacklist {
	return Blacklist{
		CommandMemento: cqrs.NewCommand(DOMAIN, C_Blacklist, visitorId, -1),
		RequesterId:    requestingUserId,
	}
}

type RescindBlacklist struct {
	cqrs.CommandMemento
	RequesterId uint64 `json:"requestintUserId"`
}

func NewRescindBlacklist(visitorId uint64, requestingUserId uint64) RescindBlacklist {
	return RescindBlacklist{
		CommandMemento: cqrs.NewCommand(DOMAIN, C_RescindBlacklist, visitorId, -1),
		RequesterId:    requestingUserId,
	}
}

type Whitelist struct {
	cqrs.CommandMemento
	RequesterId uint64 `json:"requestintUserId"`
}

func NewWhitelist(visitorId uint64, requestingUserId uint64) Whitelist {
	return Whitelist{
		CommandMemento: cqrs.NewCommand(DOMAIN, C_Whitelist, visitorId, -1),
		RequesterId:    requestingUserId,
	}
}

type RescindWhitelist struct {
	cqrs.CommandMemento
	RequesterId uint64 `json:"requestintUserId"`
}

func NewRescindWhitelist(visitorId, requestingUserId uint64) RescindWhitelist {
	return RescindWhitelist{
		CommandMemento: cqrs.NewCommand(DOMAIN, C_RescindWhitelist, visitorId, -1),
		RequesterId:    requestingUserId,
	}
}

type Registered struct {
	cqrs.EventMemento
	RequesterId uint64    `json:"requestinguserid"`
	IPV4Address uint32    `json:"ipv4"`
	IPV6Address [2]uint64 `json:"ipv6"`
}

func NewRegistered(visitorId uint64, requestingUserId uint64, ipv4Address uint32, ipv6Address [2]uint64) Registered {
	return Registered{
		EventMemento: cqrs.NewEvent(DOMAIN, E_Registered, visitorId, -1),
		RequesterId:  requestingUserId,
		IPV4Address:  ipv4Address,
		IPV6Address:  ipv6Address,
	}
}

type RequestHandled struct {
	cqrs.EventMemento
	IPAddress int32  `json:"ipaddress"`
	Request   []byte `json:"request"`
}

// TODO: Create id from IP hash
func NewRequestHandled(visitorId uint64, request []byte) RequestHandled {
	return RequestHandled{
		EventMemento: cqrs.NewEvent(DOMAIN, E_RequestHandled, visitorId, -1),
		Request:      request,
	}
}

type Blacklisted struct {
	cqrs.EventMemento
	RequesterId uint64 `json:"requestinguserid"`
}

func NewBlacklisted(visitorId uint64, requestingUserId uint64) Blacklisted {
	return Blacklisted{
		EventMemento: cqrs.NewEvent(DOMAIN, E_Blacklisted, visitorId, -1),
		RequesterId:  requestingUserId,
	}
}

type BlacklistRescinded struct {
	cqrs.EventMemento
	RequesterId uint64 `json:"requestinguserid"`
}

func NewBlacklistRescinded(visitorId uint64, requestingUserId uint64) BlacklistRescinded {
	return BlacklistRescinded{
		EventMemento: cqrs.NewEvent(DOMAIN, E_BlacklistRescinded, visitorId, -1),
		RequesterId:  requestingUserId,
	}
}

type Whitelisted struct {
	cqrs.EventMemento
	RequesterId uint64 `json:"requestinguserid"`
}

func NewWhitelisted(visitorId uint64, requestingUserId uint64) Whitelisted {
	return Whitelisted{
		EventMemento: cqrs.NewEvent(DOMAIN, E_Whitelisted, visitorId, -1),
		RequesterId:  requestingUserId,
	}
}

type WhitelistRescinded struct {
	cqrs.EventMemento
	RequesterId uint64 `json:"requestinguserid"`
}

func NewWhitelistRescinded(visitorId uint64, requestingUserId uint64) WhitelistRescinded {
	return WhitelistRescinded{
		EventMemento: cqrs.NewEvent(DOMAIN, E_WhitelistRescinded, visitorId, -1),
		RequesterId:  requestingUserId,
	}
}

func Handle(eventBus chan<- cqrs.Event, es cqrs.EventStorer, command cqrs.Command) {
	switch cmd := command.(type) {
	case Register:
		{
			fmt.Printf("Trying to register user\n->\t%v\n", command)
			eventBus <- NewRegistered(cmd.GetId(), cmd.RequesterId, cmd.IPV4Address, cmd.IPV6Address)
		}
	case HandleRequest:
		{
			fmt.Printf("Trying to handle request\n->\t%v\n", command)
			eventBus <- NewRequestHandled(cmd.GetId(), cmd.Request)
		}
	case Blacklist:
		{
			fmt.Printf("Trying to ban visitor!\n->\t%v\n", command)
			// Load aggregate
			// TODO: Change cqrs to es
			events, err := es.ReadAllEvents(command.(cqrs.Aggregate))

			if len(events) == 0 { // New visitor
				//eventBus <- NewRegistered(command.(cqrs.Aggregate).GetId(), cmd.RequesterId)
			}

			fmt.Printf("E [ \n\t%v\n ] \nErr [ %s ]\n\n", events, err)
			// Check validation
			// Emit events
			eventBus <- NewBlacklisted(command.(cqrs.Aggregate).GetId(), cmd.RequesterId)
			// Poof done
		}
	case RescindBlacklist:
		{
			fmt.Printf("Trying to lift a visitor ban...\n->\t%v\n", command)
			// Load aggregate
			// TODO: Change cqrs to es
			events, err := es.ReadAllEvents(command.(cqrs.Aggregate))

			fmt.Printf("E [ \n\t%v\n ] \tErr [ %s ]\n\n", events, err)
			// Check validation
			// Emit events
			eventBus <- NewBlacklistRescinded(command.(cqrs.Aggregate).GetId(), cmd.RequesterId)
			// Poof done
		}
	case Whitelist:
		{
			fmt.Printf("Trying to whitelist visitor\n->\t%v\n", command)
			eventBus <- NewWhitelisted(cmd.GetId(), cmd.RequesterId)
		}
	case RescindWhitelist:
		{
			fmt.Printf("Trying to rescind whitelist for visitor\n->\t%v\n", command)
			eventBus <- NewWhitelistRescinded(cmd.GetId(), cmd.RequesterId)
		}
	default:
		{
			fmt.Printf("Visitor was unable to handle command: [ %v ]\n", cmd)
		}
	}

}
