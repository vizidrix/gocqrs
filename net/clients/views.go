package clients

import (
	"errors"
	"fmt"
	"github.com/vizidrix/gocqrs/cqrs"
)

type ClientSessionsView struct {
	Clients map[string]uint64
}

func NewClientSessionsView() ClientSessionsView {
	return ClientSessionsView{
		Clients: make(map[string]uint64),
	}
}

func (view *ClientSessionsView) GetBySession(session string) (uint64, error) {
	if client, valid := view.Clients[session]; !valid {
		return 0, errors.New("invalid_session_id")
	} else {
		return client, nil
	}
}

func ClientSessionsViewHandler(eventChan chan cqrs.Event, clientview *ClientSessionsView) {
	for {
		select {
		case newevent := <-eventChan:
			switch event := newevent.(type) {
			case ClientAdded:
				fmt.Printf("\nClient View Handler: Adding client %d to View", event.GetId())
				clientview.Clients[event.Session] = event.GetId()
			case ClientSessionUpdated:
				fmt.Printf("\nClient View Handler: Updating session for client %d in View", event.GetId())
				for session, client := range clientview.Clients {
					if client == event.GetId() {
						delete(clientview.Clients, session)
						clientview.Clients[event.Session] = client
					}
				}
			case ClientRemoved:
				fmt.Printf("\nClient View Handler: Removing client %d from View", event.GetId())
				for session, client := range clientview.Clients {
					if client == event.GetId() {
						delete(clientview.Clients, session)
					}
				}
			default:
				fmt.Println(errors.New("Invalid client view event"))
			}
		}
		fmt.Printf("\nClients ->\n\t %+v", clientview)
	}
}

func NewClientSessionsViewHandler(eventbus cqrs.EventRouter, clientview *ClientSessionsView) error {
	subscription, err := eventbus.Subscribe(cqrs.ByEventTypes(
		E_ClientAdded,
		E_ClientSessionUpdated,
		E_ClientRemoved,
	))
	if err != nil {
		return err
	}

	go ClientSessionsViewHandler(subscription.EventChan(), clientview)

	return nil
}
