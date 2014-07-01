package clients

import (
	"errors"
	"fmt"
	"github.com/vizidrix/gocqrs/cqrs"
)

var (
	ErrInvalidSession = errors.New("Cannot locate client by invalid session")
	ErrActiveSession  = errors.New("Cannot register client by active session")
	ErrInactiveClient = errors.New("No sessions to delete for an inactive client")
)

type ClientSessionsView struct {
	Clients map[string]uint64
}

func NewClientSessionsView() ClientSessionsView {
	return ClientSessionsView{
		Clients: make(map[string]uint64),
	}
}

func (view *ClientSessionsView) RegisterClientBySession(session string, client uint64) error {
	if _, active := view.Clients[session]; !active {
		view.Clients[session] = client
		return nil
	} else {
		view.Clients[session] = client
		return ErrActiveSession
	}
}

func (view *ClientSessionsView) GetBySession(session string) (uint64, error) {
	if client, valid := view.Clients[session]; !valid {
		return 0, ErrInvalidSession
	} else {
		return client, nil
	}
}

func (view *ClientSessionsView) DeleteByClient(clientid uint64) error {
	sessions := 0
	for session, client := range view.Clients {
		if client == clientid {
			delete(view.Clients, session)
			sessions++
		}
	}

	if sessions > 0 {
		return nil
	} else {
		return ErrInactiveClient
	}
}

func (view *ClientSessionsView) HandleEvent(newevent cqrs.Event) {
	switch event := newevent.(type) {
	case ClientAdded:
		//		fmt.Printf("\nClient View Handler: Adding client %d to View", event.GetId())
		view.RegisterClientBySession(event.Session, event.GetId())
	case ClientSessionUpdated:
		//		fmt.Printf("\nClient View Handler: Updating session for client %d in View", event.GetId())
		view.DeleteByClient(event.GetId())
		view.RegisterClientBySession(event.Session, event.GetId())
	case ClientRemoved:
		//		fmt.Printf("\nClient View Handler: Removing client %d from View", event.GetId())
		view.DeleteByClient(event.GetId())
	default:
		fmt.Println(errors.New("Invalid client view event"))
	}
}

func (view *ClientSessionsView) HandleStream(eventchan chan cqrs.Event) {
	for {
		select {
		case event := <-eventchan:
			view.HandleEvent(event)
		}
		fmt.Printf("\nClients ->\n\t %+v", view)
	}
}

func NewClientSessionsViewHandler(eventbus cqrs.EventRouter, view *ClientSessionsView) error {
	subscription, err := eventbus.Subscribe(cqrs.ByEventTypes(
		E_ClientAdded,
		E_ClientSessionUpdated,
		E_ClientRemoved,
	))
	if err != nil {
		return err
	}

	go view.HandleStream(subscription.EventChan())

	return nil
}
