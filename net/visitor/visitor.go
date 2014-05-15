package visitor

import (
	"fmt"
	"github.com/vizidrix/gocqrs/cqrs"
)

// TODO: Create id from IP hash

var DOMAIN_NAME string = "github.com/vizidrix/gocqrs/cqrs/web/server"
var DOMAIN uint32 = 0xD937B694

var ( /* Commands */
	C_StartServer			= cqrs.C(1, 1)
	C_StopServer 			= cqrs.C(1, 2)
	C_HandleHttpRequest    	= cqrs.C(1, 3)
)

var ( /* Events */
	E_ServerStarted			= cqrs.E(1, 1)
	E_ServerStopped			= cqrs.E(1, 2)
	E_NewVisitorObserved	= cqrs.E(1, 3)
	E_HttpRequestHandled    = cqrs.E(1, 4)
)

type VisitorMemento struct {
	cqrs.AggregateMemento
	IPV4Address uint32    	`json:"ipv4"`
	IPV6Address [2]uint64 	`json:"ipv6"`
}

func NewVisitor(id uint64) VisitorMemento {
	return VisitorMemento{
		AggregateMemento: cqrs.NewAggregate(DOMAIN, id, 0),
		IPV4Address:      0,
		IPV6Address:      [2]uint64{0, 0},
	}
}

type Register struct {
	cqrs.CommandMemento
	RequestedBy uint64    `json:"requestedby"`
	IPV4Address uint32    `json:"ipv4"`
	IPV6Address [2]uint64 `json:"ipv6"`
}

func NewRegister(visitorId uint64, requestedBy uint64, ipv4Address uint32, ipv6Address [2]uint64) Register {
	return Register{
		CommandMemento: cqrs.NewCommand(DOMAIN, C_Register, visitorId, -1),
		RequestedBy:    requestedBy,
		IPV4Address:    ipv4Address,
		IPV6Address:    ipv6Address,
	}
}

func NewRegisterIPV4(visitorId uint64, requestedBy uint64, ipv4Address uint32) Register {
	return NewRegister(visitorId, requestedBy, ipv4Address, [2]uint64{0, 0})
}

func NewRegisterIPV6(visitorId, requestedBy uint64, ipv6Address [2]uint64) Register {
	return NewRegister(visitorId, requestedBy, 0, ipv6Address)
}

type HandleRequest struct {
	cqrs.CommandMemento
	RequestedBy uint64 		`json:"requestedby"`
	Request     []byte 		`json:"request"`
	IPV4Address uint32    	`json:"ipv4"`
	IPV6Address [2]uint64 	`json:"ipv6"`
}

func NewHandleRequest(visitorId uint64, request []byte, ipv4Address uint32, ipv6Address [2]uint64) HandleRequest {
	return HandleRequest{
		CommandMemento: cqrs.NewCommand(DOMAIN, C_HandleRequest, visitorId, -1),
		Request:        request,
	}
}

type Registered struct {
	cqrs.EventMemento
	RegisteredById uint64   `json:"registeredbyid"`
	IPV4Address uint32    	`json:"ipv4"`
	IPV6Address [2]uint64 	`json:"ipv6"`
}

func NewRegistered(visitorId, requestingUserId uint64, ipv4Address uint32, ipv6Address [2]uint64) Registered {
	return Registered{
		EventMemento: cqrs.NewEvent(DOMAIN, E_Registered, visitorId, -1),
		RequestingUserId:  requestingUserId,
		IPV4Address:  ipv4Address,
		IPV6Address:  ipv6Address,
	}
}

type RequestHandled struct {
	cqrs.EventMemento
	RequestingUserId uint64 	`json:"requestinguserid"`
	Request   []byte 			`json:"request"`
	IPV4Address uint32    		`json:"ipv4"`
	IPV6Address [2]uint64 		`json:"ipv6"`
}

func NewRequestHandled(visitorId, requestingUserId uint64, request []byte, ipv4Address uint32, ipv6Addrss [2]uint64) RequestHandled {
	return RequestHandled{
		EventMemento: cqrs.NewEvent(DOMAIN, E_RequestHandled, visitorId, -1),

		Request:      request,
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
	default:
		{
			fmt.Printf("Visitor was unable to handle command: [ %v ]\n", cmd)
		}
	}

}
