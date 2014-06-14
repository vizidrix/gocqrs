package clients

import (
	"github.com/vizidrix/gocqrs/cqrs"
)

var DOMAIN_NAME = "clients"
var DOMAIN uint32 = 0x00C82E74

var (
	C_AddClient           uint64 = cqrs.C(DOMAIN, 1, 1)
	C_UpdateClientSession uint64 = cqrs.C(DOMAIN, 1, 2)
	C_RemoveClient        uint64 = cqrs.C(DOMAIN, 1, 3)
)

var (
	E_ClientAdded          uint64 = cqrs.E(DOMAIN, 1, 1)
	E_ClientSessionUpdated uint64 = cqrs.E(DOMAIN, 1, 2)
	E_ClientRemoved        uint64 = cqrs.E(DOMAIN, 1, 3)
)

type ClientMemento struct {
	cqrs.AggregateMemento
	Session string
}

func NewClientMemento(id uint64) ClientMemento {
	return ClientMemento{
		AggregateMemento: cqrs.NewAggregate(DOMAIN, id, 0),
		Session:          "",
	}
}

/* Commands */

type AddClient struct {
	cqrs.CommandMemento
	Session string
}

func NewAddClient(clientId uint64, session string) AddClient {
	return AddClient{
		CommandMemento: cqrs.NewCommand(clientId, 0, C_AddClient),
		Session:        session,
	}
}

type UpdateClientSession struct {
	cqrs.CommandMemento
	Session string
}

func NewUpdateClientSession(clientId uint64, clientVersion uint32, clientid uint64, session string) UpdateClientSession {
	return UpdateClientSession{
		CommandMemento: cqrs.NewCommand(clientId, clientVersion, C_UpdateClientSession),
		Session:        session,
	}
}

type RemoveClient struct {
	cqrs.CommandMemento
}

func NewRemoveClient(clientId uint64, clientid uint64) RemoveClient {
	return RemoveClient{
		CommandMemento: cqrs.NewCommand(clientId, 0, C_RemoveClient),
	}
}

/* Events */

type ClientAdded struct {
	cqrs.EventMemento
	Session string
}

func (event ClientAdded) String() string {
	return fmt.Sprintf(" <E [ <A D[%d] ID[%d] V[%d] \\> -> E[%d] ] E\\> ",
		event.GetDomain(), event.GetId(), event.GetVersion(), "Client Added")
}

func NewClientAdded(clientId uint64, clientId uint64, session string) ClientAdded {
	return ClientAdded{
		EventMemento: cqrs.NewEvent(clientId, 0, E_ClientAdded),
		Session:      session,
	}
}

type ClientSessionUpdated struct {
	cqrs.EventMemento
	Session string
}

func (event ClientSessionUpdated) String() string {
	return fmt.Sprintf(" <E [ <A D[%d] ID[%d] V[%d] \\> -> E[%d] ] E\\> ",
		event.GetDomain(), event.GetId(), event.GetVersion(), "Client Session Updated")
}

func NewClientSessionUpdated(clientId uint64, clientId uint64, session string) ClientSessionUpdated {
	return ClientSessionUpdated{
		EventMemento: cqrs.NewEvent(clientId, 0, E_ClientSessionUpdated),
		Session:      session,
	}
}

type ClientRemoved struct {
	cqrs.EventMemento
}

func (event ClientRemoved) String() string {
	return fmt.Sprintf(" <E [ <A D[%d] ID[%d] V[%d] \\> -> E[%d] ] E\\> ",
		event.GetDomain(), event.GetId(), event.GetVersion(), "Client Removed")
}

func NewClientRemoved(clientId uint64, clientId uint64) ClientRemoved {
	return ClientRemoved{
		EventMemento: cqrs.NewEvent(clientId, 0, E_ClientRemoved),
	}
}
