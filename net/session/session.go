package sessions

import (
	"fmt"
	"github.com/vizidrix/gocqrs/cqrs"
)

var DOMAIN_NAME = "sessions"
var DOMAIN uint32 = 0x00C82E74

var (
	C_AddSession        uint64 = cqrs.C(DOMAIN, 1, 1)
	C_UpdateSessionUser uint64 = cqrs.C(DOMAIN, 1, 2)
	C_InvalidateSession uint64 = cqrs.C(DOMAIN, 1, 3)
)

var (
	E_SessionAdded          uint64 = cqrs.E(DOMAIN, 1, 1)
	E_SessionWebUserUpdated uint64 = cqrs.E(DOMAIN, 1, 2)
	E_SessionInvalidated    uint64 = cqrs.E(DOMAIN, 1, 3)
)

type SessionMemento struct {
	cqrs.AggregateMemento
	User cqrs.AggregateMemento
}

func NewSessionMemento(id uint64, user cqrs.AggregateMemento) SessionMemento {
	return SessionMemento{
		AggregateMemento: cqrs.NewAggregate(DOMAIN, id, 0),
		User:             user,
	}
}

/* Commands */

type AddSession struct {
	cqrs.CommandMemento
	User cqrs.AggregateMemento
}

func NewAddSession(sessionId uint64, user cqrs.AggregateMemento) AddSession {
	return AddSession{
		CommandMemento: cqrs.NewCommand(DOMAIN, sessionId, 0, C_AddSession),
		User:           user,
	}
}

type UpdateSessionUser struct {
	cqrs.CommandMemento
	User cqrs.AggregateMemento
}

func NewUpdateSessionUser(sessionId uint64, sessionVersion uint32, user cqrs.AggregateMemento) UpdateSessionUser {
	return UpdateSessionUser{
		CommandMemento: cqrs.NewCommand(DOMAIN, sessionId, sessionVersion, C_UpdateSessionUser),
		User:           user,
	}
}

type InvalidateSession struct {
	cqrs.CommandMemento
}

func NewInvalidateSession(sessionId uint64, sessionVersion uint32) InvalidateSession {
	return InvalidateSession{
		CommandMemento: cqrs.NewCommand(DOMAIN, sessionId, sessionVersion, C_InvalidateSession),
	}
}

/* Events */

type SessionAdded struct {
	cqrs.EventMemento
	User cqrs.AggregateMemento
}

func (event SessionAdded) String() string {
	return fmt.Sprintf(" <E [ <A D[%d] ID[%d] V[%d] \\> -> E[%d] ] E\\> ",
		event.GetDomain(), event.GetId(), event.GetVersion(), "Session Added")
}

func NewSessionAdded(sessionId uint64, user cqrs.AggregateMemento) SessionAdded {
	return SessionAdded{
		EventMemento: cqrs.NewEvent(DOMAIN, sessionId, 0, E_SessionAdded),
		User:         user,
	}
}

type SessionUserUpdated struct {
	cqrs.EventMemento
	User cqrs.AggregateMemento
}

func (event SessionUserUpdated) String() string {
	return fmt.Sprintf(" <E [ <A D[%d] ID[%d] V[%d] \\> -> E[%d] ] E\\> ",
		event.GetDomain(), event.GetId(), event.GetVersion(), "Session Session Updated")
}

func NewSessionUserUpdated(sessionId uint64, sessionVersion uint32, user cqrs.AggregateMemento) SessionUserUpdated {
	return SessionUserUpdated{
		EventMemento: cqrs.NewEvent(DOMAIN, sessionId, sessionVersion, E_SessionUserUpdated),
		User:         user,
	}
}

type SessionInvalidated struct {
	cqrs.EventMemento
}

func (event SessionInvalidated) String() string {
	return fmt.Sprintf(" <E [ <A D[%d] ID[%d] V[%d] \\> -> E[%d] ] E\\> ",
		event.GetDomain(), event.GetId(), event.GetVersion(), "Session Invalidated")
}

func NewSessionInvalidated(sessionId uint64, sessionVersion uint32) SessionInvalidated {
	return SessionInvalidated{
		EventMemento: cqrs.NewEvent(DOMAIN, sessionId, sessionVersion, E_SessionInvalidated),
	}
}
