package correlations

import (
	"github.com/vizidrix/gocqrs/cqrs"
)

var DOMAIN_NAME = "correlations"
var DOMAIN uint32 = 0xE7ECE41A

var (
	C_AddClient    uint32 = cqrs.C(DOMAIN, 1, 1)
	C_ExpireClient uint32 = cqrs.C(DOMAIN, 1, 2)
)

var (
	E_ClientAdded   uint32 = cqrs.E(DOMAIN, 1, 1)
	E_ClientExpired uint32 = cqrs.E(DOMAIN, 1, 2)
)

/* Domain */

type ClientMemento struct {
	cqrs.AggregateMemento
	SessionId int `json:"session"` // 128 bit hex
}

func NewClient(id uint64, session int, client uint64) ClientMemento {
	return ClientMemento{
		AggregateMemento: cqrs.NewAggregate(DOMAIN, id, 0),
		SessionId:        session,
	}
}

type CorrelatedCommand interface {
	cqrs.Command
	GetClientId() uint64
	GetCorrelation() uint64
}

type CorrelatedEvent interface {
	cqrs.Event
	GetClientId() uint64
	GetCorrelation() uint64
}

type CorrelationMemento struct {
	ClientId    uint64 `json:"__clientid"`
	Correlation uint64 `json:"__correlation"`
}

func NewCorrelation(clientid uint64, correlation uint64) CorrelationMemento {
	return CorrelationMemento{
		Client:      clientid,
		Correlation: correlation,
	}
}

func (correlation CorrelationMemento) GetClientId() uint64 {
	return correlation.ClientId
}

func (correlation CorrelationMemento) GetCorrelation() uint64 {
	return correlation.Correlation
}

/* Commands */

type AddClient struct {
	cqrs.CommandMemento
	SessionId int `json:"session"`
}

func NewAddClient(id uint64, session int) AddClient {
	return AddClient{
		CommandMemento: cqrs.NewCommand(id, 0, C_AddClient),
		SessionId:      session,
	}
}

type ExpireClient struct {
	cqrs.Command
}

func NewExpireClient(id uint64, version uint32) ExpireClient {
	return ExpireClient{
		CommandMemento: cqrs.NewCommand(id, version, C_ExpireClient),
	}
}

/* Events */

type ClientAdded struct {
	cqrs.Event
	SessionId int `json:"session"`
}

func (event ClientAdded) String() string {
	return fmt.Sprintf(" <E [ <A D[%d] ID[%d] V[%d] \\> -> E[%d] ] E\\> ",
		event.GetDomain(), event.GetId(), event.GetVersion(), "Client Added")
}

func NewClientAdded(id uint64, session int) ClientAdded {
	return ClientAdded{
		EventMemento: cqrs.NewEvent(id, 0, E_ClientAdded),
		SessionId:    session,
	}
}

type ClientExpired struct {
	cqrs.Event
}

func (event ClientExpired) String() string {
	return fmt.Sprintf(" <E [ <A D[%d] ID[%d] V[%d] \\> -> E[%d] ] E\\> ",
		event.GetDomain(), event.GetId(), event.GetVersion(), "Client Expired")
}

func NewClientExpired(id uint64, version uint32) ClientExpired {
	return ClientExpired{
		EventMemento: cqrs.NewEvent(id, version, E_ClientExpired),
	}
}
