package correlations

import (
	"fmt"
	"github.com/vizidrix/gocqrs/cqrs"
)

var DOMAIN_NAME = "correlations"
var DOMAIN uint32 = 0xE7ECE41A

var (
	C_AddClientUser    uint64 = cqrs.C(DOMAIN, 1, 1)
	C_ExpireClientUser uint64 = cqrs.C(DOMAIN, 1, 2)
)

var (
	E_ClientUserAdded   uint64 = cqrs.E(DOMAIN, 1, 1)
	E_ClientUserExpired uint64 = cqrs.E(DOMAIN, 1, 2)
)

/* Domain */

type ClientUserMemento struct {
	cqrs.AggregateMemento
	SessionId int `json:"session"` // 128 bit hex
}

func NewClientUser(id uint64, session int, clientuser uint64) ClientUserMemento {
	return ClientUserMemento{
		AggregateMemento: cqrs.NewAggregate(DOMAIN, id, 0),
		SessionId:        session,
	}
}

type CorrelatedCommand interface {
	cqrs.Command
	GetClientUserId() uint64
	GetCorrelation() uint64
}

type CorrelatedEvent interface {
	cqrs.Event
	GetClientUserId() uint64
	GetCorrelation() uint64
}

type CorrelationMemento struct {
	ClientUserId uint64 `json:"__userid"`
	Correlation  uint64 `json:"__correlation"`
}

func NewCorrelation(clientuserid uint64, correlation uint64) CorrelationMemento {
	return CorrelationMemento{
		ClientUserId: clientuserid,
		Correlation:  correlation,
	}
}

func (correlation CorrelationMemento) GetClientUserId() uint64 {
	return correlation.ClientUserId
}

func (correlation CorrelationMemento) GetCorrelation() uint64 {
	return correlation.Correlation
}

/* Commands */

type AddClientUser struct {
	cqrs.CommandMemento
	SessionId int `json:"session"`
}

func NewAddClientUser(id uint64, session int) AddClientUser {
	return AddClientUser{
		CommandMemento: cqrs.NewCommand(DOMAIN, id, 0, C_AddClientUser),
		SessionId:      session,
	}
}

type ExpireClientUser struct {
	cqrs.CommandMemento
}

func NewExpireClientUser(id uint64, version uint32) ExpireClientUser {
	return ExpireClientUser{
		CommandMemento: cqrs.NewCommand(DOMAIN, id, version, C_ExpireClientUser),
	}
}

/* Events */

type ClientUserAdded struct {
	cqrs.EventMemento
	SessionId int `json:"session"`
}

func (event ClientUserAdded) String() string {
	return fmt.Sprintf(" <E [ <A D[%d] ID[%d] V[%d] \\> -> E[%d] ] E\\> ",
		event.GetDomain(), event.GetId(), event.GetVersion(), "Client User Added")
}

func NewClientUserAdded(id uint64, session int) ClientUserAdded {
	return ClientUserAdded{
		EventMemento: cqrs.NewEvent(DOMAIN, id, 0, E_ClientUserAdded),
		SessionId:    session,
	}
}

type ClientUserExpired struct {
	cqrs.EventMemento
}

func (event ClientUserExpired) String() string {
	return fmt.Sprintf(" <E [ <A D[%d] ID[%d] V[%d] \\> -> E[%d] ] E\\> ",
		event.GetDomain(), event.GetId(), event.GetVersion(), "Client User Expired")
}

func NewClientUserExpired(id uint64, version uint32) ClientUserExpired {
	return ClientUserExpired{
		EventMemento: cqrs.NewEvent(DOMAIN, id, version, E_ClientUserExpired),
	}
}
