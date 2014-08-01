package webclients

import (
	"fmt"
	"github.com/vizidrix/gocqrs/cqrs"
)

var DOMAIN_NAME = "webclients"
var DOMAIN uint32 = 0x00C82E74

var (
	C_AddWebClient           uint64 = cqrs.C(DOMAIN, 1, 1)
	C_UpdateWebClientSession uint64 = cqrs.C(DOMAIN, 1, 2)
	C_RemoveWebClient        uint64 = cqrs.C(DOMAIN, 1, 3)
)

var (
	E_WebClientAdded          uint64 = cqrs.E(DOMAIN, 1, 1)
	E_WebClientSessionUpdated uint64 = cqrs.E(DOMAIN, 1, 2)
	E_WebClientRemoved        uint64 = cqrs.E(DOMAIN, 1, 3)
)

type WebClientMemento struct {
	cqrs.AggregateMemento
	Session string
}

func NewWebClientMemento(id uint64) WebClientMemento {
	return WebClientMemento{
		AggregateMemento: cqrs.NewAggregate(DOMAIN, id, 0),
		Session:          "",
	}
}

/* Commands */

type AddWebClient struct {
	cqrs.CommandMemento
	Session string
}

func NewAddWebClient(webclientId uint64, session string) AddWebClient {
	return AddWebClient{
		CommandMemento: cqrs.NewCommand(DOMAIN, webclientId, 0, C_AddWebClient),
		Session:        session,
	}
}

type UpdateWebClientSession struct {
	cqrs.CommandMemento
	Session string
}

func NewUpdateWebClientSession(webclientId uint64, webclientVersion uint32, session string) UpdateWebClientSession {
	return UpdateWebClientSession{
		CommandMemento: cqrs.NewCommand(DOMAIN, webclientId, webclientVersion, C_UpdateWebClientSession),
		Session:        session,
	}
}

type RemoveWebClient struct {
	cqrs.CommandMemento
	Session string
}

func NewRemoveWebClient(webclientId uint64, webclientVersion uint32, session string) RemoveWebClient {
	return RemoveWebClient{
		CommandMemento: cqrs.NewCommand(DOMAIN, webclientId, webclientVersion, C_RemoveWebClient),
		Session:        session,
	}
}

/* Events */

type WebClientAdded struct {
	cqrs.EventMemento
	Session string
}

func (event WebClientAdded) String() string {
	return fmt.Sprintf(" <E [ <A D[%d] ID[%d] V[%d] \\> -> E[%d] ] E\\> ",
		event.GetDomain(), event.GetId(), event.GetVersion(), "WebClient Added")
}

func NewWebClientAdded(webclientId uint64, session string) WebClientAdded {
	return WebClientAdded{
		EventMemento: cqrs.NewEvent(DOMAIN, webclientId, 0, E_WebClientAdded),
		Session:      session,
	}
}

type WebClientSessionUpdated struct {
	cqrs.EventMemento
	Session string
}

func (event WebClientSessionUpdated) String() string {
	return fmt.Sprintf(" <E [ <A D[%d] ID[%d] V[%d] \\> -> E[%d] ] E\\> ",
		event.GetDomain(), event.GetId(), event.GetVersion(), "WebClient Session Updated")
}

func NewWebClientSessionUpdated(webclientId uint64, webclientVersion uint32, session string) WebClientSessionUpdated {
	return WebClientSessionUpdated{
		EventMemento: cqrs.NewEvent(DOMAIN, webclientId, webclientVersion, E_WebClientSessionUpdated),
		Session:      session,
	}
}

type WebClientRemoved struct {
	cqrs.EventMemento
	Session string
}

func (event WebClientRemoved) String() string {
	return fmt.Sprintf(" <E [ <A D[%d] ID[%d] V[%d] \\> -> E[%d] ] E\\> ",
		event.GetDomain(), event.GetId(), event.GetVersion(), "WebClient Removed")
}

func NewWebClientRemoved(webclientId uint64, webclientVersion uint32, session string) WebClientRemoved {
	return WebClientRemoved{
		EventMemento: cqrs.NewEvent(DOMAIN, webclientId, webclientVersion, E_WebClientRemoved),
		Session:      session,
	}
}
